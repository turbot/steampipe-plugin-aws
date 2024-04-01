package aws

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	cloudwatchlogsTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"

	cloudwatchlogsv1 "github.com/aws/aws-sdk-go/service/cloudwatchlogs"

	"github.com/turbot/go-kit/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func tableAwsCloudwatchLogEventListKeyColumns() []*plugin.KeyColumn {
	return []*plugin.KeyColumn{
		{Name: "log_group_name"},
		{Name: "log_stream_name", Require: plugin.Optional},
		{Name: "filter", Require: plugin.Optional, CacheMatch: "exact"},
		{Name: "region", Require: plugin.Optional},
		{Name: "timestamp", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
	}
}

func tableAwsCloudwatchLogEvent(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudwatch_log_event",
		Description: "AWS CloudWatch Log Event",
		List: &plugin.ListConfig{
			Hydrate:    listCloudwatchLogEvents,
			Tags:       map[string]string{"service": "logs", "action": "FilterLogEvents"},
			KeyColumns: tableAwsCloudwatchLogEventListKeyColumns(),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(cloudwatchlogsv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			// Top columns
			{
				Name:        "log_group_name",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("log_group_name"),
				Description: "The name of the log group to which this event belongs.",
			},
			{
				Name:        "log_stream_name",
				Type:        proto.ColumnType_STRING,
				Description: "The name of the log stream to which this event belongs.",
			},
			{
				Name:        "event_id",
				Type:        proto.ColumnType_STRING,
				Description: "The ID of the event.",
			},
			{
				Name:        "timestamp",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Timestamp").Transform(transform.UnixMsToTimestamp),
				Description: "The time when the event occurred.",
			},
			{
				Name:        "ingestion_time",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("IngestionTime").Transform(transform.UnixMsToTimestamp),
				Description: "The time when the event was ingested.",
			},
			{
				Name:        "filter",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromQual("filter"),
				Description: "Filter pattern for the search.",
			},
			{
				Name:        "message",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Message").Transform(trim),
				Description: "The data contained in the log event.",
			},
			{
				Name:        "message_json",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Message").Transform(trim).Transform(cloudwatchLogsMesssageJson),
				Description: "The data contained in the log event in json format. Only if data is valid json string.",
			},
			// Other columns
		}),
	}
}

func trim(_ context.Context, d *transform.TransformData) (interface{}, error) {
	if d.Value == nil {
		return nil, nil
	}
	s := types.ToString(d.Value)
	return strings.TrimSpace(s), nil
}

func listCloudwatchLogEvents(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Get client
	svc, err := CloudWatchLogsClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_cloudwatch_log_event.listCloudwatchLogEvents", "get_client_error", err)
		return nil, err
	}

	equalQuals := d.EqualsQuals

	// Limiting the results
	maxLimit := int32(10000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			if limit < 1 {
				maxLimit = 1
			} else {
				maxLimit = limit
			}
		}
	}

	params := &cloudwatchlogs.FilterLogEventsInput{
		LogGroupName: aws.String(equalQuals["log_group_name"].GetStringValue()),
		Limit:        aws.Int32(maxLimit),
	}

	if equalQuals["log_stream_name"] != nil {
		params.LogStreamNames = []string{(equalQuals["log_stream_name"].GetStringValue())}
	}

	if equalQuals["filter"] != nil {
		params.FilterPattern = aws.String(equalQuals["filter"].GetStringValue())
	} else {
		if filterStrings := buildFilter(equalQuals); len(filterStrings) > 0 {
			params.FilterPattern = aws.String(strings.Join(filterStrings, " "))
		}
	}

	quals := d.Quals

	if quals["timestamp"] != nil {
		for _, q := range quals["timestamp"].Quals {
			tsSecs := q.Value.GetTimestampValue().GetSeconds()
			tsMs := tsSecs * 1000
			switch q.Operator {
			case "=":
				params.StartTime = aws.Int64(tsMs)
				params.EndTime = aws.Int64(tsMs)
			case ">=", ">":
				params.StartTime = aws.Int64(tsMs)
			case "<", "<=":
				params.EndTime = aws.Int64(tsMs)
			}
		}
	}

	paginator := cloudwatchlogs.NewFilterLogEventsPaginator(svc, params, func(o *cloudwatchlogs.FilterLogEventsPaginatorOptions) {
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_cloudwatch_alarm.listCloudWatchAlarms", "api_error", err)
			return nil, err
		}
		for _, logEvent := range output.Events {
			d.StreamListItem(ctx, logEvent)
			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, err
}

//// TRANSFORM FUNCTIONS

func cloudwatchLogsMesssageJson(_ context.Context, d *transform.TransformData) (interface{}, error) {
	event := d.HydrateItem.(cloudwatchlogsTypes.FilteredLogEvent)
	var eventMessage interface{}
	err := json.Unmarshal([]byte(*event.Message), &eventMessage)
	if err != nil {
		return nil, nil
	}
	return eventMessage, nil
}
