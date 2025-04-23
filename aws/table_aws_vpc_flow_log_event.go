package aws

import (
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"

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
		return nil, listS3FlowLogEvents(ctx, d)
	default:
		return listCloudwatchLogEvents(ctx, d, h)
	}
}

// -----------------------------------------------------------------------------
// S3 implementation
// -----------------------------------------------------------------------------

type s3FlowLogEvent struct {
	types.FilteredLogEvent
	BucketName string
	S3Key      string
}

func listS3FlowLogEvents(ctx context.Context, d *plugin.QueryData) error {
	bucketQual := d.EqualsQuals["bucket_name"]
	if bucketQual == nil {
		return fmt.Errorf("bucket_name must be provided when log_source = 's3'")
	}

	bucket := bucketQual.GetStringValue()

	// prefix: default to account‑agnostic if not provided
	prefix := ""
	if q := d.EqualsQuals["s3_prefix"]; q != nil {
		prefix = q.GetStringValue()
	}

	region := d.EqualsQualString("region") // always set when you use SupportedRegionMatrix
	if region == "" {
		region = "us-west-2" // (safety fallback – should never hit)
	}

	svc, err := S3Client(ctx, d, region)
	if err != nil {
		return err
	}

	paginator := s3.NewListObjectsV2Paginator(svc, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	})

	// compile optional time range
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

	// object key timestamp regex: _YYYYMMDDTHHMMZ_
	reTs := regexp.MustCompile(`_(\d{8}T\d{4})Z_`)

	// stream each event
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		for _, obj := range page.Contents {
			key := aws.ToString(obj.Key)
			if !strings.HasSuffix(key, ".log.gz") {
				continue
			}

			if (startTime != nil || endTime != nil) && reTs.MatchString(key) {
				tsPart := reTs.FindStringSubmatch(key)[1]
				fileTs, _ := time.Parse("20060102T1504", tsPart)
				if startTime != nil && fileTs.Before(*startTime) {
					continue
				}
				if endTime != nil && fileTs.After(*endTime) {
					continue
				}
			}

			// download object
			objOut, err := svc.GetObject(ctx, &s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(key),
			})
			if err != nil {
				return err
			}
			defer objOut.Body.Close()

			gr, err := gzip.NewReader(objOut.Body)
			if err != nil {
				return err
			}
			defer gr.Close()

			reader := bufio.NewReader(gr)
			lineNum := 0
			for {
				line, err := reader.ReadString('\n')
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}
				line = strings.TrimSpace(line)
				if line == "" ||
					strings.HasPrefix(line, "#") ||
					strings.HasPrefix(line, "version ") { // <─ NEW
					continue
				}

				fields := strings.Fields(line)
				if len(fields) < 14 {
					continue // malformed
				}

				// build event
				endField := fields[11]
				endUnix, _ := strconv.ParseInt(endField, 10, 64)
				endMillis := endUnix * 1000

				var ingestion *int64
				if obj.LastModified != nil {
					ingestion = aws.Int64(obj.LastModified.UnixMilli())
				}

				ev := s3FlowLogEvent{
					FilteredLogEvent: types.FilteredLogEvent{
						Message:       aws.String(line),
						EventId:       aws.String(fmt.Sprintf("%s:%d", key, lineNum)),
						Timestamp:     aws.Int64(endMillis),
						IngestionTime: ingestion,
					},
					BucketName: bucket,
					S3Key:      key,
				}

				// in‑memory row‑level filtering (reuse helper)
				if !matchesQuals(ev, d) {
					lineNum++
					continue
				}

				d.StreamListItem(ctx, ev)
				lineNum++
			}
		}
	}

	return nil
}

// matchesQuals evaluates optional quals for S3‑sourced rows
func matchesQuals(ev s3FlowLogEvent, d *plugin.QueryData) bool {
	equalQuals := d.EqualsQuals

	// leverage buildFilter helper to construct list of strings required.
	filters := buildFilter(equalQuals)
	if len(filters) == 0 {
		return true
	}

	msg := aws.ToString(ev.Message)
	for _, f := range filters {
		if !strings.Contains(msg, f) {
			return false
		}
	}
	return true
}

func getMessageField(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	e := h.Item.(types.FilteredLogEvent)
	fields := strings.Fields(*e.Message)
	return fields, nil
}

func getField(_ context.Context, d *transform.TransformData) (interface{}, error) {
	fields := d.Value.([]string)
	idx := d.Param.(int)
	if fields[idx] == "-" {
		return nil, nil
	}
	return fields[idx], nil
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
