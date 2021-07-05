package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCloudwatchLogEventJSON(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cloudwatch_log_event_json",
		Description: "AWS CloudWatch Log Events in JSON format",
		List: &plugin.ListConfig{
			Hydrate:    listCloudwatchLogEvents,
			KeyColumns: tableAwsCloudwatchLogEventListKeyColumns(),
		},
		GetMatrixItem: BuildRegionList,
		Columns: awsRegionalColumns([]*plugin.Column{
			// Top columns
			{Name: "log_group_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("log_group_name"), Description: "The name of the log group to which this event belongs."},
			{Name: "log_stream_name", Type: proto.ColumnType_STRING, Description: "The name of the log stream to which this event belongs."},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Timestamp").Transform(transform.UnixMsToTimestamp), Description: "The time when the event occurred."},
			{Name: "message", Type: proto.ColumnType_JSON, Description: "The data contained in the log event in JSON format."},
			// Other columns
			{Name: "event_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("EventId"), Description: "The ID of the event."},
			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Filter pattern for the search."},
			{Name: "ingestion_time", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("IngestionTime").Transform(transform.UnixMsToTimestamp), Description: "The time when the event was ingested."},
		}),
	}
}
