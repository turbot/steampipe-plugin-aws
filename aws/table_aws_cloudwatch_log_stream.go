package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"

	cloudwatchlogsv1 "github.com/aws/aws-sdk-go/service/cloudwatchlogs"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type logStreamInfo = struct {
	LogStream types.LogStream
	LogGroup  string
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
		GetMatrixItemFunc: SupportedRegionMatrix(cloudwatchlogsv1.EndpointsID),
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
	logGroup := h.Item.(types.LogGroup)

	// Get client
	svc, err := CloudWatchLogsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_stream.listCloudwatchLogStreams", "client_error", err)
		return nil, err
	}

	maxItems := int32(50)

	// Reduce the basic request limit down if the user has only requested a small number
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}

	input := &cloudwatchlogs.DescribeLogStreamsInput{
		Limit:        &maxItems,
		LogGroupName: logGroup.LogGroupName,
	}

	paginator := cloudwatchlogs.NewDescribeLogStreamsPaginator(svc, input, func(o *cloudwatchlogs.DescribeLogStreamsPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	// Additonal Filter
	equalQuals := d.EqualsQuals
	if equalQuals["name"] != nil {
		input.LogStreamNamePrefix = aws.String(equalQuals["name"].GetStringValue())
	}

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Info("aws_cloudwatch_log_group.listCloudwatchLogGroups", "api_error", err)
			return nil, err
		}

		for _, logStream := range output.LogStreams {
			d.StreamListItem(ctx, logStreamInfo{logStream, *logGroup.LogGroupName})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getCloudwatchLogStream(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	name := d.EqualsQuals["name"].GetStringValue()
	logGroupName := d.EqualsQuals["log_group_name"].GetStringValue()

	// Error: pq: rpc error: code = Unknown desc = InvalidParameter: 2 validation error(s) found.
	// - minimum field size of 1, DescribeLogStreamsInput.LogGroupName.
	// - minimum field size of 1, DescribeLogStreamsInput.LogStreamNamePrefix.
	if len(name) < 1 || len(logGroupName) < 1 {
		return nil, nil
	}

	// Create session
	svc, err := CloudWatchLogsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Trace("aws_cloudwatch_log_group.getCloudwatchLogStream", "client_error", err)
		return nil, err
	}

	params := &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        aws.String(logGroupName),
		LogStreamNamePrefix: aws.String(name),
	}

	// execute list call
	item, err := svc.DescribeLogStreams(ctx, params)
	if err != nil {
		return nil, err
	}

	for _, logStream := range item.LogStreams {
		if *logStream.LogStreamName == name {
			return logStreamInfo{logStream, logGroupName}, nil
		}
	}

	return nil, nil
}
