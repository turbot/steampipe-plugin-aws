package aws

import (
	"context"
	"fmt"
	"github.com/turbot/steampipe-plugin-aws/aws/vpcflowlogs"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	cloudwatchlogsv1 "github.com/aws/aws-sdk-go/service/cloudwatchlogs"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"github.com/turbot/steampipe-plugin-sdk/v5/query_cache"
)

// -----------------------------------------------------------------------------
// Key columns
// -----------------------------------------------------------------------------

func tableAwsVpcFlowLogEventListKeyColumns() []*plugin.KeyColumn {
	return []*plugin.KeyColumn{
		// selector
		{Name: "log_source", Require: plugin.Optional}, // "cloudwatch" (default) | "s3"
		// cloudwatch‑only
		{Name: "log_group_name", Require: plugin.Optional},
		// s3‑only
		{Name: "bucket_name", Require: plugin.Optional},
		{Name: "s3_prefix", Require: plugin.Optional},
		{Name: "extraction_time", Require: plugin.Optional},
		// shared optionals
		{Name: "log_stream_name", Require: plugin.Optional},
		{Name: "filter", Require: plugin.Optional, CacheMatch: query_cache.CacheMatchExact},
		{Name: "region", Require: plugin.Optional},
		{Name: "timestamp", Operators: []string{">", ">=", "=", "<", "<="}, Require: plugin.Optional},
		{Name: "event_id", Require: plugin.Optional},
		{Name: "interface_id", Require: plugin.Optional},
		{Name: "src_addr", Require: plugin.Optional},
		{Name: "dst_addr", Require: plugin.Optional},
		{Name: "src_port", Require: plugin.Optional},
		{Name: "dst_port", Require: plugin.Optional},
		{Name: "action", Require: plugin.Optional},
		{Name: "log_status", Require: plugin.Optional},
	}
}

// -----------------------------------------------------------------------------
// Table definition
// -----------------------------------------------------------------------------

func tableAwsVpcFlowLogEvent(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_vpc_flow_log_event",
		Description: "AWS VPC Flow Log events from CloudWatch Logs or S3 objects (compressed .log.gz)",
		List: &plugin.ListConfig{
			Hydrate:    listVpcFlowLogEvents,
			Tags:       map[string]string{"service": "logs", "action": "FilterLogEvents"},
			KeyColumns: tableAwsVpcFlowLogEventListKeyColumns(),
		},
		GetMatrixItemFunc: SupportedRegionMatrix(cloudwatchlogsv1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			// selector / provenance cols
			{Name: "log_source", Type: proto.ColumnType_STRING, Transform: transform.FromQual("log_source"), Description: "Source of the flow logs: cloudwatch (default) or s3."},
			{Name: "bucket_name", Type: proto.ColumnType_STRING, Transform: transform.FromField("BucketName"), Description: "S3 bucket containing the log file (when log_source = 's3')."},
			{Name: "s3_prefix", Type: proto.ColumnType_STRING, Transform: transform.FromQual("s3_prefix"), Description: "S3 prefix to search for flow logs (when log_source = 's3')."},
			{Name: "s3_key", Type: proto.ColumnType_STRING, Transform: transform.FromField("S3Key"), Description: "Full S3 object key for the record (when log_source = 's3')."},
			{Name: "extraction_time", Type: proto.ColumnType_INT, Transform: transform.FromQual("extraction_time"), Description: "Time limit in seconds to process flow log files (when log_source = 's3'). Defaults to 600 seconds (10 minutes) if not specified."},

			// Top columns (existing)
			{Name: "log_group_name", Type: proto.ColumnType_STRING, Transform: transform.FromQual("log_group_name"), Description: "The name of the log group to which this event belongs (CloudWatch source)."},
			{Name: "log_stream_name", Type: proto.ColumnType_STRING, Description: "The name of the log stream to which this event belongs."},
			{Name: "timestamp", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("Timestamp").Transform(transform.UnixMsToTimestamp), Description: "The time when the event occurred (maps to the record 'end' field)."},
			{Name: "version", Type: proto.ColumnType_INT, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 0), Description: "The VPC Flow Logs version."},
			{Name: "interface_account_id", Type: proto.ColumnType_STRING, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 1), Description: "The AWS account ID of the owner of the network interface."},
			{Name: "interface_id", Type: proto.ColumnType_STRING, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 2), Description: "The ID of the network interface."},
			{Name: "src_addr", Type: proto.ColumnType_IPADDR, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 3), Description: "Source IP address."},
			{Name: "dst_addr", Type: proto.ColumnType_IPADDR, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 4), Description: "Destination IP address."},
			{Name: "src_port", Type: proto.ColumnType_INT, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 5), Description: "Source port."},
			{Name: "dst_port", Type: proto.ColumnType_INT, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 6), Description: "Destination port."},
			{Name: "protocol", Type: proto.ColumnType_INT, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 7), Description: "IANA protocol number."},
			{Name: "packets", Type: proto.ColumnType_INT, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 8), Description: "Number of packets transferred."},
			{Name: "bytes", Type: proto.ColumnType_INT, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 9), Description: "Number of bytes transferred."},
			{Name: "start", Type: proto.ColumnType_TIMESTAMP, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 10).Transform(transform.UnixToTimestamp), Description: "Time when first packet of the flow was received."},
			{Name: "end", Type: proto.ColumnType_TIMESTAMP, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 11).Transform(transform.UnixToTimestamp), Description: "Time when last packet of the flow was received."},
			{Name: "action", Type: proto.ColumnType_STRING, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 12), Description: "ACCEPT | REJECT."},
			{Name: "log_status", Type: proto.ColumnType_STRING, Hydrate: getMessageField, Transform: transform.FromValue().TransformP(getField, 13), Description: "Logging status."},

			// Other columns
			{Name: "event_id", Type: proto.ColumnType_STRING, Transform: transform.FromField("EventId"), Description: "Event ID (CloudWatch) or synthetic ID (S3)."},
			{Name: "filter", Type: proto.ColumnType_STRING, Transform: transform.FromQual("filter"), Description: "Filter pattern used (CloudWatch source)."},
			{Name: "ingestion_time", Type: proto.ColumnType_TIMESTAMP, Transform: transform.FromField("IngestionTime").Transform(transform.UnixMsToTimestamp), Description: "Time when event was ingested (CloudWatch) or object last modified (S3)."},
		}),
	}
}

// -----------------------------------------------------------------------------
// Routing List Hydrate – decides cloudwatch vs s3
// -----------------------------------------------------------------------------

func listVpcFlowLogEvents(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Default is cloudwatch
	logSource := "cloudwatch"
	if q := d.EqualsQuals["log_source"]; q != nil {
		v := strings.ToLower(q.GetStringValue())
		if v == "s3" {
			logSource = "s3"
		}
	}

	switch logSource {
	case "s3":
		return listCloudwatchLogEventsFromS3(ctx, d)
	default:
		return listCloudwatchLogEvents(ctx, d, h)
	}
}

func listCloudwatchLogEventsFromS3(ctx context.Context, d *plugin.QueryData) (interface{}, error) {
	region := d.EqualsQualString("region") // always set when you use SupportedRegionMatrix
	if region == "" {
		region = "us-west-2" // (safety fallback – should never hit)
	}

	bucketQual := d.EqualsQuals["bucket_name"]
	if bucketQual == nil {
		return nil, fmt.Errorf("bucket_name must be provided when log_source = 's3'")
	}
	bucket := bucketQual.GetStringValue()

	prefix := ""
	if q := d.EqualsQuals["s3_prefix"]; q != nil {
		prefix = q.GetStringValue()
	}

	// compile optional time range
	startTime, endTime := getStartEndTime(d)

	// Initialize the S3 client
	s3Client, err := S3Client(ctx, d, region)
	if err != nil {
		return nil, err
	}

	// Get extraction time from query parameters or use default
	extractionTime := getExtractionTime(d, vpcflowlogs.DefaultTimeoutSeconds)

	retriever := vpcflowlogs.NewS3FlowLogEventsRetriever(buildFilter(d.EqualsQuals), d.StreamListItem, s3Client, region, bucket, prefix, startTime, endTime, plugin.Logger(ctx))

	return nil, retriever.ListS3FlowLogEvents(ctx, extractionTime)
}

func getStartEndTime(d *plugin.QueryData) (*time.Time, *time.Time) {
	var startTime, endTime *time.Time
	if quals, ok := d.Quals["timestamp"]; ok {
		for _, q := range quals.Quals {
			t := q.Value.GetTimestampValue().AsTime()
			switch q.Operator {
			case ">", ">=":
				if startTime == nil || t.After(*startTime) {
					startTime = &t
				}
			case "<", "<=":
				if endTime == nil || t.Before(*endTime) {
					endTime = &t
				}
			case "=":
				startTime = &t
				endTime = &t
			}
		}
	}
	return startTime, endTime
}

func getMessageField(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	switch v := h.Item.(type) {
	case types.FilteredLogEvent:
		// CloudWatch rows
		fields := strings.Fields(*v.Message)
		return fields, nil
	case vpcflowlogs.S3FlowLogEvent:
		// S3 rows - access the embedded FilteredLogEvent's Message field
		fields := strings.Fields(*v.Message)
		return fields, nil
	default:
		return nil, fmt.Errorf("unknown item type %T in getMessageField", v)
	}
}

func getField(_ context.Context, d *transform.TransformData) (interface{}, error) {
	fields := d.Value.([]string)
	idx := d.Param.(int)
	if fields[idx] == "-" {
		return nil, nil
	}
	return fields[idx], nil
}

// getExtractionTime extracts the extraction time parameter from query data
// or returns the default extraction time if not set
func getExtractionTime(d *plugin.QueryData, defaultTime int) int {
	if q := d.EqualsQuals["extraction_time"]; q != nil {
		extractionTime := int(q.GetInt64Value())
		if extractionTime <= 0 {
			return defaultTime
		}
		return extractionTime
	}
	return defaultTime
}

func buildFilter(equalQuals plugin.KeyColumnEqualsQualMap) []string {
	filters := []string{}

	filterQuals := []string{"action", "log_status", "interface_id", "event_id", "src_addr", "dst_addr", "src_port", "dst_port"}

	for _, qual := range filterQuals {
		switch qual {
		case "action", "log_status", "interface_id", "event_id":
			if equalQuals[qual] != nil {
				filters = append(filters, equalQuals[qual].GetStringValue())
			}
		case "src_addr", "dst_addr":
			if equalQuals[qual] != nil {
				filters = append(filters, equalQuals[qual].GetInetValue().Addr)
			}
		case "src_port", "dst_port":
			if equalQuals[qual] != nil {
				filters = append(filters, strconv.Itoa(int(equalQuals[qual].GetInt64Value())))
			}
		}
	}

	return filters
}
