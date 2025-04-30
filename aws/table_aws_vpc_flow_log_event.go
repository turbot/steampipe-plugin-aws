package aws

import (
	"bufio"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
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
		region := d.EqualsQualString("region") // always set when you use SupportedRegionMatrix
		if region == "" {
			region = "us-west-2" // (safety fallback – should never hit)
		}

		// Initialize the S3 client
		s3Client, err := S3Client(ctx, d, region)
		if err != nil {
			return nil, err
		}

		// Get extraction time from query parameters or use default (10 minutes)
		const defaultExtractionTime = 600
		extractionTime := getExtractionTime(d, defaultExtractionTime)

		return nil, listS3FlowLogEvents(ctx, d, s3Client, extractionTime)
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
func listS3FlowLogEvents(ctx context.Context, d *plugin.QueryData, s3Client S3ClientInterface, extractionTime int) error {
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

	// Use the extraction time passed as a parameter
	logger.Debug("listS3FlowLogEvents", "extraction_time", extractionTime)

	// Create a context with timeout based on extraction time
	ctx, cancel := context.WithTimeout(ctx, time.Duration(extractionTime)*time.Second)
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

	// Start parallel goroutine for progressive processing with direct time slot targeting
	// This implementation starts processing files immediately while continuing to discover more
	logger.Debug("listS3FlowLogEvents", "message", "Starting progressive S3 object processing")
	go func() {
		objectCount := 0
		processedCount := 0
		defer func() {
			// Close channels to signal completion
			close(objChan)
			close(listingDoneChan)
			logger.Debug("listS3FlowLogEvents", "message", "S3 object listing complete",
				"eligible_objects_found", objectCount, "processed_objects", processedCount)
		}()

		// If we have no time bounds, use default listing
		if startTime == nil && endTime == nil {
			// Use standard listing for cases with no time filtering
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

					// Found an eligible object, send it to worker pool
					objectCount++
					processedCount++
					logger.Trace("listS3FlowLogEvents", "message", "Found eligible object",
						"key", key, "count", objectCount)
					select {
					case objChan <- obj:
					case <-ctx.Done():
						return
					}
				}
			}
			return
		}

		// For time-bounded searches, use progressive processing with direct time slot targeting
		// First, compute date range for our search
		// Default to a reasonable range if not fully specified
		now := time.Now()
		searchStartTime := startTime
		searchEndTime := endTime
		if searchStartTime == nil {
			// Default to 1 day before end time or current time
			if searchEndTime != nil {
				defaultStart := searchEndTime.Add(-24 * time.Hour)
				searchStartTime = &defaultStart
			} else {
				defaultStart := now.Add(-24 * time.Hour)
				searchStartTime = &defaultStart
			}
		}
		if searchEndTime == nil {
			// Default to now if no end time provided
			searchEndTime = &now
		}

		// Create a structure to track dates and hours we need to process
		type timeTarget struct {
			date time.Time
			hour int
		}

		// Build a queue of time targets to process, distributed across the time range
		var timeTargets []timeTarget

		// Get start and end date at day granularity
		startDate := time.Date(searchStartTime.Year(), searchStartTime.Month(), searchStartTime.Day(), 0, 0, 0, 0, searchStartTime.Location())
		endDate := time.Date(searchEndTime.Year(), searchEndTime.Month(), searchEndTime.Day(), 0, 0, 0, 0, searchEndTime.Location())

		// Calculate total days in range
		totalDays := int(endDate.Sub(startDate).Hours()/24) + 1
		logger.Debug("listS3FlowLogEvents", "message", "Time range calculated",
			"start_date", startDate.Format("2006-01-02"),
			"end_date", endDate.Format("2006-01-02"),
			"total_days", totalDays)

		// Create time targets that are well-distributed across the time range
		// For each day, we'll create targets for specific hours
		for day := 0; day < totalDays; day++ {
			currentDate := startDate.AddDate(0, 0, day)

			// For the first and last day, respect the actual start/end time's hour
			var startHour, endHour int
			if day == 0 {
				// First day - start at the actual start time's hour
				startHour = searchStartTime.Hour()
			} else {
				startHour = 0
			}

			if day == totalDays-1 && currentDate.Equal(endDate) {
				// Last day - end at the actual end time's hour
				endHour = searchEndTime.Hour() + 1 // +1 to include the hour of end time
			} else {
				endHour = 24
			}

			// Create targets for specific hours, balanced across the time range
			// Use a stride to ensure even distribution
			if endHour > startHour {
				// Determine how many hours to sample based on the range
				// For narrow ranges, sample every hour
				// For wider ranges, sample fewer hours per day
				stride := 1
				if totalDays > 3 {
					stride = 4 // Sample every 4 hours for ranges > 3 days
				} else if totalDays > 1 {
					stride = 2 // Sample every 2 hours for 2-3 day ranges
				}

				// Add time targets with appropriate stride
				for hour := startHour; hour < endHour; hour += stride {
					timeTargets = append(timeTargets, timeTarget{date: currentDate, hour: hour})
				}
			}
		}

		// Log the targeting plan
		logger.Debug("listS3FlowLogEvents", "message", "Created time targets for distributed sampling",
			"target_count", len(timeTargets))

		// Shuffle time targets to improve distribution in case of timeout
		// This helps ensure we don't process all early times first
		// Using Fisher-Yates shuffle
		for i := len(timeTargets) - 1; i > 0; i-- {
			j := i
			if i > 0 {
				// Use a simple deterministic shuffle based on the indices
				// since we don't have a real random source here
				j = (i*7 + 3) % (i + 1)
			}
			timeTargets[i], timeTargets[j] = timeTargets[j], timeTargets[i]
		}

		// Process each time target concurrently (but limit concurrent API calls)
		// This allows us to process files from across the time range immediately
		// rather than waiting for the entire listing to complete

		// Create a throttling channel to limit concurrent API calls
		// This prevents too many simultaneous S3 requests
		concurrencyLimit := 5
		throttle := make(chan struct{}, concurrencyLimit)

		// Create a wait group to track completion of all list operations
		var listWg sync.WaitGroup

		// Process time targets
		for _, target := range timeTargets {
			listWg.Add(1)

			// Throttle concurrent API calls
			throttle <- struct{}{}

			// Start a goroutine to process this time target
			go func(date time.Time, hour int) {
				defer listWg.Done()
				defer func() { <-throttle }() // Release throttle when done

				// Format the hour for the pattern-matching prefix
				hourStr := fmt.Sprintf("%02d", hour)

				// Build pattern for direct time slot targeting
				// Structure is typically: base_prefix/YYYY/MM/DD/AccountID_vpcflowlogs_region_*_YYYYMMDDTHHMM*
				dateStr := fmt.Sprintf("%d%02d%02d", date.Year(), date.Month(), date.Day())
				timePattern := dateStr + "T" + hourStr

				// Build the date component of the prefix
				datePrefix := prefix
				if !strings.HasSuffix(datePrefix, "/") {
					datePrefix += "/"
				}

				// If the prefix doesn't already include date components, add them
				if !strings.Contains(datePrefix, fmt.Sprintf("/%d/%02d/%02d/", date.Year(), date.Month(), date.Day())) {
					datePrefix += fmt.Sprintf("%d/%02d/%02d/", date.Year(), date.Month(), date.Day())
				}

				logger.Debug("listS3FlowLogEvents", "message", "Directly targeting time slot",
					"date", dateStr, "hour", hourStr, "prefix", datePrefix)

				// List objects for this specific time slot
				paginator := s3.NewListObjectsV2Paginator(s3Client, &s3.ListObjectsV2Input{
					Bucket: aws.String(bucket),
					Prefix: aws.String(datePrefix),
					// Use a small page size to get first results quickly
					MaxKeys: aws.Int32(50),
				})

				// Process each page for this time slot
				timeSlotCount := 0
				for paginator.HasMorePages() {
					// Check context before fetching next page
					if ctx.Err() != nil {
						return
					}

					page, err := paginator.NextPage(ctx)
					if err != nil {
						logger.Error("listS3FlowLogEvents", "message", "Error listing objects for time slot",
							"date", dateStr, "hour", hourStr, "error", err)
						if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
							return
						}
						select {
						case errorChan <- err:
						case <-ctx.Done():
						}
						return
					}

					// Filter objects that match our time pattern
					var slotObjects []s3types.Object
					for _, obj := range page.Contents {
						key := aws.ToString(obj.Key)
						if !strings.HasSuffix(key, ".log.gz") {
							continue
						}

						// Match our specific time slot (hour)
						// Use the original regex to extract the timestamp
						if matches := reTs.FindStringSubmatch(key); len(matches) > 1 {
							tsPart := matches[1]

							// Check if this file belongs to our target hour
							if strings.HasPrefix(tsPart, timePattern) {
								fileTs, err := time.Parse("20060102T1504", tsPart)
								if err != nil {
									logger.Warn("listS3FlowLogEvents", "message", "Failed to parse timestamp from key",
										"key", key, "error", err)
									continue
								}

								// Double-check that file is within our time bounds
								if (startTime != nil && fileTs.Before(*startTime)) ||
									(endTime != nil && fileTs.After(*endTime)) {
									continue
								}

								slotObjects = append(slotObjects, obj)
								timeSlotCount++
							}
						}
					}

					// Process matched objects for this time slot
					// Send objects directly to processing
					for i := 0; i < len(slotObjects) && ctx.Err() == nil; i++ {
						obj := slotObjects[i]
						key := aws.ToString(obj.Key)

						// Atomically increment counts
						atomic.AddInt32((*int32)(&objectCount), 1)
						atomic.AddInt32((*int32)(&processedCount), 1)

						logger.Trace("listS3FlowLogEvents", "message", "Processing object from time slot",
							"date", dateStr, "hour", hourStr, "key", key)

						// Send object to processing channel
						select {
						case objChan <- obj:
						case <-ctx.Done():
							return
						}
					}
				}

				logger.Debug("listS3FlowLogEvents", "message", "Completed processing time slot",
					"date", dateStr, "hour", hourStr, "objects", timeSlotCount)

			}(target.date, target.hour)
		}

		// Wait for all listing operations to complete or context to be canceled
		listingDone := make(chan struct{})
		go func() {
			listWg.Wait()
			close(listingDone)
		}()

		// Wait for either completion or cancellation
		select {
		case <-listingDone:
			logger.Debug("listS3FlowLogEvents", "message", "All time slots processed",
				"total_objects", objectCount)
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				logger.Info("listS3FlowLogEvents", "message", "Time-distributed processing stopped due to timeout",
					"extraction_time_seconds", extractionTime, "objects_processed", processedCount)
			} else {
				logger.Info("listS3FlowLogEvents", "message", "Time-distributed processing cancelled",
					"objects_processed", processedCount)
			}
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
					if errors.Is(ctx.Err(), context.DeadlineExceeded) {
						logger.Info("listS3FlowLogEvents", "message", "Operation timed out after reaching extraction_time limit, returning partial results",
							"extraction_time_seconds", extractionTime, "results_found", resultCount)
						return nil
					}
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
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					logger.Info("listS3FlowLogEvents", "message", "Operation timed out after reaching extraction_time limit, returning partial results",
						"extraction_time_seconds", extractionTime, "results_found", resultCount)
					return nil
				}
				return ctx.Err()
			}
		case <-ctx.Done():
			// Context cancelled externally or timeout
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				logger.Info("listS3FlowLogEvents", "message", "Operation timed out after reaching extraction_time limit, returning partial results",
					"extraction_time_seconds", extractionTime, "results_found", resultCount)
				return nil
			}
			logger.Info("listS3FlowLogEvents", "message", "Operation cancelled by context",
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
