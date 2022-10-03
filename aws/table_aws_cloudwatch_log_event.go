package aws

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
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
			KeyColumns: tableAwsCloudwatchLogEventListKeyColumns(),
		},
		GetMatrixItemFunc: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			// Top columns
			{Name: "log_group_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("log_group_name"), Description: "The name of the log group to which this event belongs."},
			{Name: "log_stream_name", Type: proto.ColumnType_STRING, Description: "The name of the log stream to which this event belongs."},
			{Name: "event_id", Type: proto.ColumnType_STRING, Description: "The ID of the event."},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Timestamp").Transform(transform.UnixMsToTimestamp), Description: "The time when the event occurred."},
			{Name: "ingestion_time", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("IngestionTime").Transform(transform.UnixMsToTimestamp), Description: "The time when the event was ingested."},
			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Filter pattern for the search."},
			{Name: "message", Type: proto.ColumnType_STRING, Transform: transform.FromField("Message").Transform(trim), Description: "The data contained in the log event."},
			{Name: "message_json", Type: proto.ColumnType_JSON, Transform: transform.FromField("Message").Transform(trim).Transform(cloudwatchLogsMesssageJson), Description: "The data contained in the log event in json format. Only if data is valid json string."},
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

	// Create session
	svc, err := CloudWatchLogsService(ctx, d)
	if err != nil {
		return nil, err
	}

	equalQuals := d.KeyColumnQuals

	input := cloudwatchlogs.FilterLogEventsInput{
		LogGroupName: aws.String(equalQuals["log_group_name"].GetStringValue()),
		// Default to the maximum allowed
		Limit: aws.Int64(10000),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *input.Limit {
			input.Limit = limit
		}
	}

	if equalQuals["log_stream_name"] != nil {
		input.LogStreamNames = []*string{aws.String(equalQuals["log_stream_name"].GetStringValue())}
	}

	if equalQuals["filter"] != nil {
		input.FilterPattern = aws.String(equalQuals["filter"].GetStringValue())
	} else {
		if filterStrings := buildFilter(equalQuals); len(filterStrings) > 0 {
			input.FilterPattern = aws.String(strings.Join(filterStrings, " "))
		}
	}

	if input.FilterPattern != nil {
		plugin.Logger(ctx).Trace("listCloudwatchLogTrailEvents", "input.FilterPattern", *input.FilterPattern)
	}

	quals := d.Quals

	if quals["timestamp"] != nil {
		for _, q := range quals["timestamp"].Quals {
			tsSecs := q.Value.GetTimestampValue().GetSeconds()
			tsMs := tsSecs * 1000
			switch q.Operator {
			case "=":
				input.StartTime = aws.Int64(tsMs)
				input.EndTime = aws.Int64(tsMs)
			case ">=", ">":
				input.StartTime = aws.Int64(tsMs)
			case "<", "<=":
				input.EndTime = aws.Int64(tsMs)
			}
		}
	}

	err = svc.FilterLogEventsPages(
		&input,
		func(page *cloudwatchlogs.FilterLogEventsOutput, _ bool) bool {
			for _, logEvent := range page.Events {
				d.StreamListItem(ctx, logEvent)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			// Abort if we've been cancelled, which probably means we've reached the requested limit
			select {
			case <-ctx.Done():
				return false
			default:
				return true
			}
		},
	)

	// Handle log group not found errors gracefully
	if awsErr, ok := err.(awserr.Error); ok {
		if awsErr.Code() == "ResourceNotFoundException" {
			return nil, nil
		}
	}

	return nil, err
}

//// TRANSFORM FUNCTIONS

func cloudwatchLogsMesssageJson(_ context.Context, d *transform.TransformData) (interface{}, error) {
	event := d.HydrateItem.(*cloudwatchlogs.FilteredLogEvent)
	var eventMessage interface{}
	err := json.Unmarshal([]byte(*event.Message), &eventMessage)
	if err != nil {
		return nil, nil
	}
	return eventMessage, nil
}
