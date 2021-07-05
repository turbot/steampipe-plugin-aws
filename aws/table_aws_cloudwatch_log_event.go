package aws

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

func tableAwsCloudwatchLogEventListKeyColumns() []*plugin.KeyColumn {
	return []*plugin.KeyColumn{
		{Name: "log_group_name"},
		{Name: "log_stream_name", Require: plugin.Optional},
		{Name: "filter", Require: plugin.Optional},
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
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			// Top columns
			{Name: "log_group_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("log_group_name"), Description: "The name of the log group to which this event belongs."},
			{Name: "log_stream_name", Type: proto.ColumnType_STRING, Description: "The name of the log stream to which this event belongs."},
			{Name: "event_id", Type: proto.ColumnType_STRING, Description: "The ID of the event."},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Timestamp").Transform(transform.UnixMsToTimestamp), Description: "The time when the event occurred."},
			{Name: "message", Type: proto.ColumnType_STRING, Transform: transform.FromField("Message").Transform(trim), Description: "The data contained in the log event."},
			// Other columns
			{Name: "ingestion_time", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("IngestionTime").Transform(transform.UnixMsToTimestamp), Description: "The time when the event was ingested."},
			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Filter pattern for the search."},
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

func listCloudwatchLogEvents(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	var region string
	matrixRegion := plugin.GetMatrixItem(ctx)[matrixKeyRegion]
	if matrixRegion != nil {
		region = matrixRegion.(string)
	}

	// Create session
	svc, err := CloudWatchLogsService(ctx, d, region)
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
				break
			case ">=", ">":
				input.StartTime = aws.Int64(tsMs)
			case "<", "<=":
				input.EndTime = aws.Int64(tsMs)
			}
		}
	}

	err = svc.FilterLogEventsPages(
		&input,
		func(page *cloudwatchlogs.FilterLogEventsOutput, isLast bool) bool {
			for _, logEvent := range page.Events {
				d.StreamListItem(ctx, logEvent)
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
