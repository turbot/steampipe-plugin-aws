package aws

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

//// TABLE DEFINITION

var msgRegex = regexp.MustCompile(`(?m)^(.*)\t([0-9A-Fa-f]{8}-[0-9A-Fa-f]{4}-4[0-9A-Fa-f]{3}-[89AB][0-9A-Fa-f]{3}-[0-9A-Fa-f]{12})\t(INFO|WARN|ERROR)\t(.*)`)

func tableAwsLambdaLogEvent(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_lambda_log_event",
		Description: "CloudTrail events from cloudwatch service for lambda functions.",
		List: &plugin.ListConfig{
			Hydrate:       listLambdaLogEvents,
			ParentHydrate: listAwsLambdaFunctions,
			KeyColumns:    tableLambdaLogEventsListKeyColumns(),
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			// Top columns
			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "The cloudwatch filter pattern for the search."},
			{Name: "log_stream_name", Type: proto.ColumnType_STRING, Description: "The name of the log stream to which this event belongs."},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Timestamp").Transform(transform.UnixMsToTimestamp), Description: "The time when the event occurred."},
			{Name: "timestamp_ms", Type: proto.ColumnType_INT, Transform: transform.FromField("Timestamp"), Description: "The time when the event occurred."},

			{Name: "event", Type: proto.ColumnType_STRING, Hydrate: getLambdaLogEvent, Description: "The CloudTrail event in the json format."},
			{Name: "function_name", Type: proto.ColumnType_STRING, Hydrate: getLambdaLogEvent, Description: "The CloudTrail event in the json format."},
			{Name: "request_id", Type: proto.ColumnType_STRING, Hydrate: getLambdaLogEvent, Description: "The CloudTrail event in the json format."},
			{Name: "message", Type: proto.ColumnType_STRING, Hydrate: getLambdaLogEvent, Description: "The CloudTrail event in the json format."},
			{Name: "type", Type: proto.ColumnType_STRING, Hydrate: getLambdaLogEvent, Description: "The CloudTrail event in the json format."},
			{Name: "event_timestamp", Type: proto.ColumnType_STRING, Hydrate: getLambdaLogEvent, Description: "The CloudTrail event in the json format."},
		}),
	}
}

func listLambdaLogEvents(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create session
	svc, err := CloudWatchLogsService(ctx, d)
	if err != nil {
		return nil, err
	}

	lambdaFn := h.Item.(*lambda.FunctionConfiguration)

	equalQuals := d.KeyColumnQuals
	// quals := d.Quals

	input := cloudwatchlogs.FilterLogEventsInput{
		LogGroupName: aws.String(fmt.Sprintf("/aws/lambda/%s", *lambdaFn.FunctionName)),
		// Default to the maximum allowed
		Limit: aws.Int64(10000),
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if limit := d.QueryContext.Limit; limit != nil && *input.Limit < *limit {
		input.Limit = limit
	}

	if equalQuals["log_stream_name"] != nil {
		input.LogStreamNames = []*string{aws.String(equalQuals["log_stream_name"].GetStringValue())}
	}

	queryFilter := ""
	filter := buildQueryFilter(equalQuals)

	if equalQuals["filter"] != nil {
		queryFilter = equalQuals["filter"].GetStringValue()
	}

	if queryFilter != "" {
		input.FilterPattern = aws.String(queryFilter)
	} else if len(filter) > 0 {
		input.FilterPattern = aws.String(fmt.Sprintf("{ %s }", strings.Join(filter, " && ")))
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

	if input.FilterPattern != nil {
		plugin.Logger(ctx).Debug("aws_cloudtrail_trail_event.listCloudwatchLogTrailEvents", "region", d.KeyColumnQualString(matrixKeyRegion), "filter query", *input.FilterPattern)
	}

	err = svc.FilterLogEventsPages(
		&input,
		func(page *cloudwatchlogs.FilterLogEventsOutput, _ bool) bool {
			for _, logEvent := range page.Events {
				if isControlLogEvent(ctx, *(logEvent.Message)) {
					continue
				}
				d.StreamListItem(ctx, logEvent)

				// Context can be cancelled due to manual cancellation or the limit has been hit
				if d.QueryStatus.RowsRemaining(ctx) == 0 {
					return false
				}
			}
			return (ctx.Err() == nil)
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

func isControlLogEvent(ctx context.Context, msg string) bool {
	return (strings.HasPrefix(msg, "START RequestId") ||
		strings.HasPrefix(msg, "END RequestId") ||
		strings.HasPrefix(msg, "REPORT RequestId"))
}

// https://docs.aws.amazon.com/AWSJavaScriptSDK/latest/AWS/CloudTrail.html#lookupEvents-property
func tableLambdaLogEventsListKeyColumns() []*plugin.KeyColumn {
	return []*plugin.KeyColumn{
		// CloudWatch fields
		{Name: "log_stream_name", Require: plugin.Optional},
		{Name: "filter", Require: plugin.Optional, CacheMatch: "exact"},
		{Name: "region", Require: plugin.Optional},
		{Name: "timestamp", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
	}
}

func getLambdaLogEvent(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	e := h.Item.(*cloudwatchlogs.FilteredLogEvent)
	matches := msgRegex.FindAllStringSubmatch(*e.Message, -1)
	cte := lambdaLoGEvent{
		FunctionName:   *h.ParentItem.(*lambda.FunctionConfiguration).FunctionName,
		EventTimestamp: matches[0][1],
		RequestId:      matches[0][2],
		Type:           matches[0][3],
		Message:        matches[0][4],
		Event:          *e.Message,
	}

	return cte, nil
}

type lambdaLoGEvent struct {
	Event          string
	FunctionName   string
	RequestId      string
	Message        string
	Type           string
	EventTimestamp string
}
