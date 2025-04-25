package aws

import (
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
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
		region := d.EqualsQualString("region") // always set when you use SupportedRegionMatrix
		if region == "" {
			region = "us-west-2" // (safety fallback – should never hit)
		}

		// Initialize the S3 client
		s3Client, err := S3Client(ctx, d, region)
		if err != nil {
			return nil, err
		}
		return nil, listS3FlowLogEvents(ctx, d, s3Client)
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

// S3ClientInterface defines the methods we use from the S3 client for testing
type S3ClientInterface interface {
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

// listS3FlowLogEvents queries and processes VPC Flow Logs stored in S3.
//
// This function retrieves VPC flow logs from AWS S3 storage, processes them,
// and streams the results back to the client. It supports various query filters
// to narrow down the results based on user requirements.
//
// Input parameters:
//   - ctx: The context for the query, which can be used for cancellation
//   - d: The QueryData object containing qualifiers and filters for the query
//
// Expected qualifiers in QueryData:
//   - bucket_name (required): The S3 bucket containing flow log files
//   - s3_prefix (optional): Prefix path within the bucket to limit the search scope
//   - region (optional): AWS region where the bucket is located
//   - timestamp (optional): Time range filters to limit results by event timestamp
//   - event_id, interface_id, src_addr, dst_addr, src_port, dst_port, action, log_status (optional):
//     Content filters for specific fields in the flow logs
//
// Processing behavior:
//   - Files are filtered by the ".log.gz" extension to identify flow log archives
//   - Timestamp filtering is applied to filenames (if they match the expected pattern)
//   - Downloaded files are decompressed and processed line by line
//   - Content filtering is applied to match specific criteria when provided
//
// Output behavior:
//   - Each matching log line is converted to an s3FlowLogEvent object
//   - Results are streamed incrementally via d.StreamListItem() as they're processed
//   - Processing continues until all matching files are processed or an error occurs
//
// Return value:
//   - Returns nil on successful completion or if no matching files are found
//   - Returns an error if bucket_name is missing or any processing error occurs
func listS3FlowLogEvents(ctx context.Context, d *plugin.QueryData, s3Client S3ClientInterface) error {
	logger := plugin.Logger(ctx)
	logger.Debug("listS3FlowLogEvents", "message", "Starting S3 flow log retrieval operation")

	// Constant for controlling concurrency
	const maxConcurrentDownloads = 5
	const resultsChannelSize = 1000
	const objectChannelSize = 100 // Buffer for eligible objects

	bucketQual := d.EqualsQuals["bucket_name"]
	if bucketQual == nil {
		logger.Error("listS3FlowLogEvents", "error", "bucket_name qualifier missing")
		return fmt.Errorf("bucket_name must be provided when log_source = 's3'")
	}

	bucket := bucketQual.GetStringValue()

	// prefix: default to account‑agnostic if not provided
	prefix := ""
	if q := d.EqualsQuals["s3_prefix"]; q != nil {
		prefix = q.GetStringValue()
	}
	logger.Debug("listS3FlowLogEvents", "bucket", bucket, "prefix", prefix)

	// Precompile patterns for better performance
	reTs := regexp.MustCompile(`_(\d{8}T\d{4})Z_`)
	reCommentPrefix := regexp.MustCompile(`^(#|version\s)`)

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
		logger.Debug("listS3FlowLogEvents", "time_filter", "applied",
			"start_time", startTime, "end_time", endTime)
	}

	// Setup channels for concurrent processing - create a pipeline
	objChan := make(chan s3types.Object, objectChannelSize)
	resultsChan := make(chan s3FlowLogEvent, resultsChannelSize)
	errorChan := make(chan error, maxConcurrentDownloads)
	listingDoneChan := make(chan struct{})    // Signal that S3 listing is complete
	processingDoneChan := make(chan struct{}) // Signal that processing is complete

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Start workers to process objects concurrently
	// Workers begin processing immediately as objects become available from objChan
	logger.Debug("listS3FlowLogEvents", "workers", maxConcurrentDownloads,
		"message", "Starting worker goroutines")
	var workersWg sync.WaitGroup
	for i := 0; i < maxConcurrentDownloads; i++ {
		workersWg.Add(1)
		go func(workerID int) {
			logger.Trace("listS3FlowLogEvents", "worker_id", workerID,
				"message", "Worker started")
			defer workersWg.Done()
			for obj := range objChan {
				key := aws.ToString(obj.Key)
				logger.Trace("listS3FlowLogEvents", "worker_id", workerID,
					"message", "Processing object", "key", key)

				// download object
				objOut, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
					Bucket: aws.String(bucket),
					Key:    aws.String(key),
				})
				if err != nil {
					logger.Error("listS3FlowLogEvents", "worker_id", workerID,
						"message", "Failed to get object from S3",
						"key", key, "error", err)
					select {
					case errorChan <- err:
					case <-ctx.Done():
						return
					}
					continue
				}

				gr, err := gzip.NewReader(objOut.Body)
				if err != nil {
					logger.Error("listS3FlowLogEvents", "worker_id", workerID,
						"message", "Failed to create gzip reader",
						"key", key, "error", err)
					objOut.Body.Close()
					select {
					case errorChan <- err:
					case <-ctx.Done():
						return
					}
					continue
				}

				// Use a scanner for better line reading performance
				scanner := bufio.NewScanner(gr)
				// Increase buffer size for large lines
				const maxScanTokenSize = 1024 * 1024
				buf := make([]byte, maxScanTokenSize)
				scanner.Buffer(buf, maxScanTokenSize)

				// Pre-compute filters once before the loop for efficiency
				// This avoids rebuilding the filter list for each line
				filters := buildFilter(d.EqualsQuals)
				hasFilters := len(filters) > 0

				lineNum := 0
				for scanner.Scan() {
					line := scanner.Text()

					// Skip empty lines and comments more efficiently
					if line == "" || reCommentPrefix.MatchString(line) {
						continue
					}

					fields := strings.Fields(line)
					if len(fields) < 14 {
						continue // malformed
					}

					// If filters exist, do a quick check before creating the full event
					// This avoids the overhead of creating the event object for records
					// that will be filtered out anyway
					if hasFilters {
						// Quick check for filter matches
						match := true
						for _, f := range filters {
							if !strings.Contains(line, f) {
								match = false
								break
							}
						}
						if !match {
							lineNum++
							continue
						}
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

					select {
					case resultsChan <- ev:
					case <-ctx.Done():
						gr.Close()
						objOut.Body.Close()
						return
					}
					lineNum++
				}

				if err := scanner.Err(); err != nil {
					logger.Error("listS3FlowLogEvents", "worker_id", workerID,
						"message", "Scanner error while processing file",
						"key", key, "error", err)
					select {
					case errorChan <- err:
					case <-ctx.Done():
					}
				}

				logger.Trace("listS3FlowLogEvents", "worker_id", workerID,
					"message", "Completed processing object",
					"key", key, "lines_processed", lineNum)
				gr.Close()
				objOut.Body.Close()
			}
			logger.Trace("listS3FlowLogEvents", "worker_id", workerID,
				"message", "Worker finished")
		}(i) // Pass worker ID to the goroutine
	}

	// Start a goroutine to close results channel when all workers finish
	go func() {
		workersWg.Wait()
		close(resultsChan)
		close(processingDoneChan)
	}()

	// Start parallel goroutine to list and filter S3 objects
	// This sends objects to workers as they're discovered, without collecting them all first
	logger.Debug("listS3FlowLogEvents", "message", "Starting S3 object listing")
	go func() {
		objectCount := 0
		defer func() {
			// If we found zero objects, we need to close the objChan
			// to signal workers they're done and can exit
			close(objChan)
			close(listingDoneChan)
			logger.Debug("listS3FlowLogEvents", "message", "S3 object listing complete",
				"eligible_objects_found", objectCount)
		}()

		paginator := s3.NewListObjectsV2Paginator(s3Client, &s3.ListObjectsV2Input{
			Bucket: aws.String(bucket),
			Prefix: aws.String(prefix),
		})

		pageCount := 0
		for paginator.HasMorePages() {
			pageCount++
			logger.Trace("listS3FlowLogEvents", "message", "Retrieving page of S3 objects",
				"page", pageCount)

			page, err := paginator.NextPage(ctx)
			if err != nil {
				logger.Error("listS3FlowLogEvents", "message", "Error listing S3 objects",
					"error", err, "page", pageCount)
				select {
				case errorChan <- err:
				case <-ctx.Done():
				}
				return
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

				// Found an eligible object, send it to worker pool
				objectCount++
				logger.Trace("listS3FlowLogEvents", "message", "Found eligible object",
					"key", aws.ToString(obj.Key), "count", objectCount)
				select {
				case objChan <- obj:
				case <-ctx.Done():
					return
				}
			}
		}

		// If we processed the entire listing with zero objects found, log an info message
		if objectCount == 0 {
			plugin.Logger(ctx).Info("listS3FlowLogEvents", "message", "No eligible log files found in S3 bucket",
				"bucket", bucket, "prefix", prefix)
		}
	}()

	// Stream results to output while checking for errors
	noObjectsFound := false
	resultCount := 0
	logger.Debug("listS3FlowLogEvents", "message", "Starting result streaming phase")
	for {
		select {
		case result, ok := <-resultsChan:
			if !ok {
				// resultsChan closed, all processing complete
				logger.Debug("listS3FlowLogEvents", "message", "Result streaming complete",
					"total_results", resultCount)
				if noObjectsFound {
					// We know there were no objects to process, so report success
					return nil
				}
				// Wait for listing to complete before returning
				select {
				case <-listingDoneChan:
					return nil
				case err := <-errorChan:
					cancel()
					return err
				case <-ctx.Done():
					return ctx.Err()
				}
			}
			resultCount++
			if resultCount%1000 == 0 {
				logger.Debug("listS3FlowLogEvents", "message", "Streaming results",
					"count", resultCount)
			}
			d.StreamListItem(ctx, result)
		case err := <-errorChan:
			// Cancel all in-progress work
			logger.Error("listS3FlowLogEvents", "message", "Error received from worker",
				"error", err)
			cancel()
			return err
		case <-listingDoneChan:
			// S3 listing is complete, but we still need to wait for processing to finish
			noObjectsFound = true
		case <-processingDoneChan:
			// All workers finished and no more results
			logger.Debug("listS3FlowLogEvents", "message", "All processing complete")
			// Wait for listing to complete before returning
			select {
			case <-listingDoneChan:
				logger.Debug("listS3FlowLogEvents", "message", "S3 flow log retrieval operation complete",
					"total_results", resultCount)
				return nil
			case err := <-errorChan:
				cancel()
				return err
			case <-ctx.Done():
				return ctx.Err()
			}
		case <-ctx.Done():
			// Context cancelled externally
			logger.Warn("listS3FlowLogEvents", "message", "Operation cancelled by context",
				"error", ctx.Err())
			return ctx.Err()
		}
	}
}

func getMessageField(_ context.Context, _ *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	switch v := h.Item.(type) {
	case types.FilteredLogEvent:
		// CloudWatch rows
		fields := strings.Fields(*v.Message)
		return fields, nil
	case s3FlowLogEvent:
		// S3‑sourced rows (our custom struct embeds FilteredLogEvent)
		fields := strings.Fields(*v.Message)
		return fields, nil
	default:
		return nil, fmt.Errorf("Unknown item type %T in getMessageField", v)
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
