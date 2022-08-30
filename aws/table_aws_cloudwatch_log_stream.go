package aws

import (
	"context"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
)

type logStreamInfo = struct {
	LogStream *cloudwatchlogs.LogStream
	LogGroup  *string
}

//// TABLE DEFINITION

func tableAwsCloudwatchLogStream(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudwatch_log_stream",
		Description: "AWS CloudWatch Log Stream",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"log_group_name", "name"}),
			Hydrate:    getCloudwatchLogStream,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listCloudwatchLogGroups,
			Hydrate:       listCloudwatchLogStreams,
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
			},
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "name",
				Description: "The name of the log stream.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LogStream.LogStreamName"),
			},
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) of the log stream.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LogStream.Arn"),
			},
			{
				Name:        "log_group_name",
				Description: "The name of the log group, in which the log stream belongs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LogGroup"),
			},
			{
				Name:        "creation_time",
				Description: "The creation time of the log stream.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LogStream.CreationTime").Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "first_event_timestamp",
				Description: "The time of the first event.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LogStream.FirstEventTimestamp").Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "last_event_timestamp",
				Description: "The time of the most recent log event in the log stream in CloudWatch Logs.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LogStream.LastEventTimestamp").Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "last_ingestion_time",
				Description: "Specifies the last log ingestion time.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LogStream.LastIngestionTime").Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "upload_sequence_token",
				Description: "Specifies the log upload sequence token.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LogStream.UploadSequenceToken"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LogStream.LogStreamName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("LogStream.Arn").Transform(arnToAkas),
			},
		}),
	}
}

//// LIST FUNCTION

func listCloudwatchLogStreams(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get logGroup details
	logGroup := h.Item.(*cloudwatchlogs.LogGroup)

	// Create session
	svc, err := CloudWatchLogsService(ctx, d)
	if err != nil {
		return nil, err
	}

	input := &cloudwatchlogs.DescribeLogStreamsInput{
		Limit: aws.Int64(50),
	}
	input.LogGroupName = logGroup.LogGroupName

	// Additonal Filter
	equalQuals := d.KeyColumnQuals
	if equalQuals["name"] != nil {
		input.LogStreamNamePrefix = types.String(equalQuals["name"].GetStringValue())
	}

	// If the requested number of items is less than the paging max limit
	// set the limit to that instead
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.Limit {
			if *limit < 1 {
				input.Limit = types.Int64(1)
			} else {
				input.Limit = limit
			}
		}
	}

	err = svc.DescribeLogStreamsPages(
		input,
		func(page *cloudwatchlogs.DescribeLogStreamsOutput, _ bool) bool {
			for _, logStream := range page.LogStreams {
				d.StreamLeafListItem(ctx, logStreamInfo{logStream, logGroup.LogGroupName})

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return true
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCloudwatchLogStream(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCloudwatchLogStream")

	name := d.KeyColumnQuals["name"].GetStringValue()
	logGroupName := d.KeyColumnQuals["log_group_name"].GetStringValue()

	// Error: pq: rpc error: code = Unknown desc = InvalidParameter: 2 validation error(s) found.
	// - minimum field size of 1, DescribeLogStreamsInput.LogGroupName.
	// - minimum field size of 1, DescribeLogStreamsInput.LogStreamNamePrefix.
	if len(name) < 1 || len(logGroupName) < 1 {
		return nil, nil
	}

	// Create session
	svc, err := CloudWatchLogsService(ctx, d)
	if err != nil {
		return nil, err
	}

	params := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        aws.String(logGroupName),
		LogStreamNamePrefix: aws.String(name),
	}

	// execute list call
	item, err := svc.DescribeLogStreams(params)
	if err != nil {
		return nil, err
	}

	for _, logStream := range item.LogStreams {
		if types.SafeString(logStream.LogStreamName) == name {
			return logStreamInfo{logStream, aws.String(logGroupName)}, nil
		}
	}

	return nil, nil
}
