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
			Tags:       map[string]string{"service": "logs", "action": "DescribeLogStreams"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listCloudwatchLogGroups,
			Hydrate:       listCloudwatchLogStreams,
			Tags:          map[string]string{"service": "logs", "action": "DescribeLogStreams"},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "name",
					Require: plugin.Optional,
				},
				{
					Name:    "log_stream_name_prefix",
					Require: plugin.Optional,
				},
				{
					Name:    "descending",
					Require: plugin.Optional,
				},
				{
					Name:    "order_by",
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
				Name:        "log_stream_name_prefix",
				Description: "The prefix to match the name of the log stream.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("log_stream_name_prefix"),
			},
			{
				Name:        "descending",
				Description: "If the value is true, results are returned in descending order. If the value is to false, results are returned in ascending order. The default value is false.",
				Type:        proto.ColumnType_BOOL,
				Transform:   transform.FromQual("descending"),
				Default:     false,
			},
			{
				Name:        "order_by",
				Description: "If the value is LogStreamName, the results are ordered by log stream name. If the value is LastEventTime, the results are ordered by the event time. The default value is LogStreamName.If you order the results by event time, you cannot specify the logStreamNamePrefix parameter. lastEventTimestamp represents the time of the most recent log event in the log stream in CloudWatch Logs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("order_by"),
				Default:     "LogStreamName",
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
	if equalQuals["descending"] != nil {
		input.Descending = aws.Bool(equalQuals["descending"].GetBoolValue())
	}

	// If the value is LogStreamName, the results are ordered by log stream name. If the value is LastEventTime, the results are ordered by the event time. The default value is LogStreamName.
	// If you order the results by event time, you cannot specify the logStreamNamePrefix parameter.
	if equalQuals["order_by"] != nil {
		input.OrderBy = types.OrderBy(equalQuals["order_by"].GetStringValue())

		if input.OrderBy != types.OrderByLastEventTime {
			if equalQuals["name"] != nil {
				input.LogStreamNamePrefix = aws.String(equalQuals["name"].GetStringValue())
			}
			if equalQuals["log_stream_name_prefix"] != nil {
				input.LogStreamNamePrefix = aws.String(equalQuals["log_stream_name_prefix"].GetStringValue())
			}
		}
		if (input.OrderBy == types.OrderByLastEventTime && equalQuals["log_stream_name_prefix"] != nil) || (input.OrderBy == types.OrderByLastEventTime && equalQuals["log_stream_name_prefix"] != nil && equalQuals["name"] != nil) {
			input.LogStreamNamePrefix = aws.String(equalQuals["log_stream_name_prefix"].GetStringValue())
		}
	} else {
		if equalQuals["name"] != nil {
			input.LogStreamNamePrefix = aws.String(equalQuals["name"].GetStringValue())
		}
		if equalQuals["log_stream_name_prefix"] != nil {
			input.LogStreamNamePrefix = aws.String(equalQuals["log_stream_name_prefix"].GetStringValue())
		}
	}

	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Info("aws_cloudwatch_log_stream.listCloudwatchLogGroups", "api_error", err)
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
		plugin.Logger(ctx).Trace("aws_cloudwatch_log_stream.getCloudwatchLogStream", "client_error", err)
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