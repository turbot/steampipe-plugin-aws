package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
)

type logStreamInfo = struct {
	LogStream *cloudwatchlogs.LogStream
	LogGroup  *string
}

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
		},
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
				Transform:   transform.FromField("LogStream.CreationTime").Transform(convertTimestamp),
			},
			{
				Name:        "first_event_timestamp",
				Description: "The time of the first event.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LogStream.FirstEventTimestamp").Transform(convertTimestamp),
			},
			{
				Name:        "last_event_timestamp",
				Description: "The time of the most recent log event in the log stream in CloudWatch Logs.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LogStream.LastEventTimestamp").Transform(convertTimestamp),
			},
			{
				Name:        "last_ingestion_time",
				Description: "Specifies the last log ingestion time.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("LogStream.LastIngestionTime").Transform(convertTimestamp),
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
	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listCloudwatchLogStreams", "AWS_REGION", defaultRegion)

	// Create session
	svc, err := CloudWatchLogsService(ctx, d.ConnectionManager, defaultRegion)
	if err != nil {
		return nil, err
	}

	logGroup := h.Item.(*cloudwatchlogs.LogGroup)

	err = svc.DescribeLogStreamsPages(
		&cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName: logGroup.LogGroupName,
		},
		func(page *cloudwatchlogs.DescribeLogStreamsOutput, isLast bool) bool {
			for _, logStream := range page.LogStreams {
				d.StreamLeafListItem(ctx, logStreamInfo{logStream, logGroup.LogGroupName})
			}
			return true
		},
	)

	return nil, err
}

//// HYDRATE FUNCTIONS

func getCloudwatchLogStream(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getCloudwatchLogStream")

	defaultRegion := GetDefaultRegion()
	name := d.KeyColumnQuals["name"].GetStringValue()
	logGroupName := d.KeyColumnQuals["log_group_name"].GetStringValue()

	// Error: pq: rpc error: code = Unknown desc = InvalidParameter: 2 validation error(s) found.
	// - minimum field size of 1, DescribeLogStreamsInput.LogGroupName.
	// - minimum field size of 1, DescribeLogStreamsInput.LogStreamNamePrefix.
	if len(name) < 1 || len(logGroupName) < 1 {
		return nil, nil
	}

	// Create session
	svc, err := CloudWatchLogsService(ctx, d.ConnectionManager, defaultRegion)
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
		if *logStream.LogStreamName == *logStream.LogStreamName {
			return logStreamInfo{logStream, aws.String(logGroupName)}, nil
		}
	}

	return nil, nil
}
