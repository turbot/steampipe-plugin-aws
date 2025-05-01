package vpcflowlogs

import (
	"bufio"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/go-hclog"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Configuration constants for processing S3 flow logs
const (
	// Concurrent processing settings
	DefaultMaxConcurrentDownloads = 5    // Number of concurrent worker goroutines
	DefaultResultsChannelSize     = 1000 // Size of the channel buffer for results
	DefaultConcurrencyLimit       = 5    // Limit for concurrent S3 list operations

	// Default time targeting settings
	DefaultTimeoutSeconds = 600 // Default timeout for extraction operations (10 minutes)
	DefaultLookbackHours  = 24  // Default lookback period when not specified

	// S3 pagination settings
	DefaultS3PageSize = 50 // Page size for S3 listing operations

	// Scanner settings
	MaxScanTokenSize = 1024 * 1024 // Maximum size for scanner buffer (1MB)
)

// Precompiled regex patterns for better performance
var reCommentPrefix = regexp.MustCompile(`^(#|version\s)`) // Pattern to match comment lines in logs
var reTs = regexp.MustCompile(`_(\d{8}T\d{4})Z_`)          // Pattern to extract timestamps from S3 keys

// S3FlowLogEvent represents a flow log event from S3 storage
type S3FlowLogEvent struct {
	types.FilteredLogEvent
	BucketName string
	S3Key      string
}

// Helper functions for context handling

// sendWithContext sends a value to a channel only if the context is not done
// returns true if the value was sent, false if the context was done
func sendWithContext[T any](ctx context.Context, ch chan<- T, value T) bool {
	select {
	case <-ctx.Done():
		return false
	case ch <- value:
		return true
	}
}

// logContextError logs the appropriate message when a context error occurs
func logContextError(ctx context.Context, logger hclog.Logger, extractionTime int, resultCount int) {
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		logger.Info("listS3FlowLogEvents", "message", "Operation timed out after reaching extraction_time limit, returning partial results",
			"extraction_time_seconds", extractionTime, "results_found", resultCount)
	} else {
		logger.Info("listS3FlowLogEvents", "message", "Operation cancelled by context",
			"error", ctx.Err())
	}
}

// S3ClientInterface defines the methods we use from the S3 client for testing
type S3ClientInterface interface {
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

type ItemStreamer func(context.Context, ...interface{})

type S3FlowLogEventsRetriever struct {
	filters      []string
	itemStreamer ItemStreamer
	s3Client     S3ClientInterface
	region       string
	bucket       string
	prefix       string
	startTime    *time.Time
	endTime      *time.Time
	logger       hclog.Logger
}

func NewS3FlowLogEventsRetriever(filters []string, itemStreamer ItemStreamer, s3Client S3ClientInterface, region string, bucket string, prefix string, startTime *time.Time, endTime *time.Time, logger hclog.Logger) *S3FlowLogEventsRetriever {
	return &S3FlowLogEventsRetriever{
		filters:      filters,
		itemStreamer: itemStreamer,
		s3Client:     s3Client,
		region:       region,
		bucket:       bucket,
		prefix:       prefix,
		startTime:    startTime,
		endTime:      endTime,
		logger:       logger,
	}
}

func (r *S3FlowLogEventsRetriever) ListS3FlowLogEvents(ctx context.Context, extractionTime int) error {
	r.logger.Debug("listS3FlowLogEvents", "message", "Starting S3 flow log retrieval operation")

	// Create object pool and setup channels for concurrent processing
	objectPool := NewObjectPoolDefault[s3types.Object]()
	resultsChan := make(chan S3FlowLogEvent, DefaultResultsChannelSize)
	errorChan := make(chan error, DefaultMaxConcurrentDownloads)
	processingDoneChan := make(chan struct{}) // Signal that processing is complete

	// Use the extraction time passed as a parameter
	r.logger.Debug("listS3FlowLogEvents", "extraction_time", extractionTime)

	// Create a context with timeout based on extraction time
	ctx, cancel := context.WithTimeout(ctx, time.Duration(extractionTime)*time.Second)
	defer cancel()

	// Start workers to process objects concurrently
	// Workers begin processing immediately as objects become available from the pool
	r.logger.Debug("listS3FlowLogEvents", "workers", DefaultMaxConcurrentDownloads,
		"message", "Starting worker goroutines")
	var workersWg sync.WaitGroup
	for i := 0; i < DefaultMaxConcurrentDownloads; i++ {
		workersWg.Add(1)
		go func(workerID int) {
			defer workersWg.Done()
			r.processObjectsWorker(ctx, workerID, objectPool, resultsChan, errorChan)
		}(i) // Pass worker ID to avoid closure issues
	}

	// Start a goroutine to close results channel when all workers finish
	go func() {
		workersWg.Wait()
		close(resultsChan)
		close(processingDoneChan)
	}()

	r.logger.Debug("listS3FlowLogEvents", "message", "Starting progressive S3 object processing")
	r.processS3Objects(ctx, objectPool, errorChan, extractionTime)

	// Stream results to output while checking for errors
	// Prioritize context check before starting the streaming loop
	if ctx.Err() != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			r.logger.Info("listS3FlowLogEvents", "message", "Operation timed out before streaming results",
				"extraction_time_seconds", extractionTime)
			return nil
		}
		return ctx.Err()
	}

	resultCount := 0
	r.logger.Debug("listS3FlowLogEvents", "message", "Starting result streaming phase")

	// Create a separate goroutine to monitor context cancellation
	// This ensures we can break out of any waiting state immediately
	doneChan := make(chan struct{})
	go func() {
		<-ctx.Done()
		close(doneChan)
	}()

	for {
		// Check if context is done at the start of each loop iteration
		if ctx.Err() != nil {
			logContextError(ctx, r.logger, extractionTime, resultCount)
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return nil // Return success on timeout, just with partial results
			}
			return ctx.Err()
		}

		// Process events and check channels
		select {
		case <-doneChan:
			// Context cancellation occurred
			logContextError(ctx, r.logger, extractionTime, resultCount)
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return nil // Return success on timeout, just with partial results
			}
			return ctx.Err()

		case result, ok := <-resultsChan:
			if !ok {
				// resultsChan closed, all processing complete
				r.logger.Debug("listS3FlowLogEvents", "message", "Result streaming complete",
					"total_results", resultCount)
				// No need to wait for additional signals - just return what we have
				return nil
			}
			resultCount++
			if resultCount%1000 == 0 {
				r.logger.Debug("listS3FlowLogEvents", "message", "Streaming results",
					"count", resultCount)
			}
			// Stream with context awareness
			if ctx.Err() != nil {
				r.logger.Info("listS3FlowLogEvents", "message", "Context cancelled during streaming, returning partial results",
					"results_found", resultCount)
				return nil
			}

			// Pass result to the streamer
			r.itemStreamer(ctx, result)

		case err := <-errorChan:
			// Cancel all in-progress work
			r.logger.Error("listS3FlowLogEvents", "message", "Error received from worker",
				"error", err)
			cancel()
			return err

		case <-processingDoneChan:
			// All workers finished and no more results
			r.logger.Debug("listS3FlowLogEvents", "message", "All processing complete, returning results",
				"total_results", resultCount)
			return nil
		case <-ctx.Done():
			// Context cancelled externally or timeout
			logContextError(ctx, r.logger, extractionTime, resultCount)
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return nil
			}
			return ctx.Err()
		}
	}
}

// processObjectsWorker processes S3 objects from the object pool and sends flow log events to the results channel
func (r *S3FlowLogEventsRetriever) processObjectsWorker(
	ctx context.Context,
	workerID int,
	objectPool *ObjectPool[s3types.Object],
	resultsChan chan<- S3FlowLogEvent,
	errorChan chan<- error,
) {
	r.logger.Trace("listS3FlowLogEvents", "worker_id", workerID,
		"message", "Worker started")

	// Continue processing objects until context is cancelled or pool is empty and closed
	for {
		// Check if context is done before getting the next object
		if ctx.Err() != nil {
			r.logger.Trace("listS3FlowLogEvents", "worker_id", workerID,
				"message", "Context done, worker exiting early", "reason", ctx.Err())
			return
		}

		// Set a timeout for getting the next object
		// This ensures workers don't block indefinitely even if the pool is not empty
		getCtx, cancelGet := context.WithTimeout(ctx, 5*time.Second)
		obj, ok := objectPool.GetRandom(getCtx)
		cancelGet()

		if !ok {
			// Check again with the main context to determine why we didn't get an object
			if ctx.Err() != nil {
				r.logger.Trace("listS3FlowLogEvents", "worker_id", workerID,
					"message", "Context cancelled while waiting for object", "reason", ctx.Err())
				return
			}

			// If the main context isn't done, check if the pool is empty and closed
			if objectPool.IsEmpty() && objectPool.IsClosed() {
				r.logger.Trace("listS3FlowLogEvents", "worker_id", workerID,
					"message", "Object pool empty and closed, worker finished")
				return
			}

			// If neither of those conditions are true, we timed out waiting for an object
			// Try again in the next loop iteration
			r.logger.Trace("listS3FlowLogEvents", "worker_id", workerID,
				"message", "Timed out waiting for object, checking context and trying again")
			continue
		}

		key := aws.ToString(obj.Key)
		r.logger.Trace("listS3FlowLogEvents", "worker_id", workerID,
			"message", "Processing object", "key", key)

		// Check context before heavy processing
		if ctx.Err() != nil {
			r.logger.Trace("listS3FlowLogEvents", "worker_id", workerID,
				"message", "Context done before processing object", "key", key, "reason", ctx.Err())
			return
		}

		// download object
		objOut, err := r.s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(r.bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			if ctx.Err() != nil {
				// Context cancelled during the download
				r.logger.Debug("listS3FlowLogEvents", "worker_id", workerID,
					"message", "Context cancelled during object download", "key", key, "reason", ctx.Err())
				return
			}

			r.logger.Error("listS3FlowLogEvents", "worker_id", workerID,
				"message", "Failed to get object from S3",
				"key", key, "error", err)
			if !sendWithContext(ctx, errorChan, err) {
				return
			}
			continue
		}

		// Check context again before proceeding with decompression
		if ctx.Err() != nil {
			r.logger.Trace("listS3FlowLogEvents", "worker_id", workerID,
				"message", "Context done after object download", "key", key, "reason", ctx.Err())
			objOut.Body.Close()
			return
		}

		gr, err := gzip.NewReader(objOut.Body)
		if err != nil {
			r.logger.Error("listS3FlowLogEvents", "worker_id", workerID,
				"message", "Failed to create gzip reader",
				"key", key, "error", err)
			objOut.Body.Close()
			if !sendWithContext(ctx, errorChan, err) {
				return
			}
			continue
		}

		// Use a scanner for better line reading performance
		scanner := bufio.NewScanner(gr)
		// Increase buffer size for large lines
		buf := make([]byte, MaxScanTokenSize)
		scanner.Buffer(buf, MaxScanTokenSize)

		hasFilters := len(r.filters) > 0

		lineNum := 0
		linesSinceContextCheck := 0
		for scanner.Scan() {
			// Check context periodically during scanning to ensure timely termination
			// Only check every 100 lines to balance responsiveness and performance
			linesSinceContextCheck++
			if linesSinceContextCheck >= 100 {
				linesSinceContextCheck = 0
				if ctx.Err() != nil {
					r.logger.Debug("listS3FlowLogEvents", "worker_id", workerID,
						"message", "Context done during scanning", "key", key,
						"lines_processed", lineNum, "reason", ctx.Err())
					gr.Close()
					objOut.Body.Close()
					return
				}
			}

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
				for _, f := range r.filters {
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

			ev := S3FlowLogEvent{
				FilteredLogEvent: types.FilteredLogEvent{
					Message:       aws.String(line),
					EventId:       aws.String(fmt.Sprintf("%s:%d", key, lineNum)),
					Timestamp:     aws.Int64(endMillis),
					IngestionTime: ingestion,
				},
				BucketName: r.bucket,
				S3Key:      key,
			}

			if !sendWithContext(ctx, resultsChan, ev) {
				r.logger.Debug("listS3FlowLogEvents", "worker_id", workerID,
					"message", "Failed to send event to channel, context done",
					"key", key, "lines_processed", lineNum)
				gr.Close()
				objOut.Body.Close()
				return
			}
			lineNum++
		}

		if err := scanner.Err(); err != nil {
			r.logger.Error("listS3FlowLogEvents", "worker_id", workerID,
				"message", "Scanner error while processing file",
				"key", key, "error", err)
			if ctx.Err() == nil {
				// Only send the error if the context isn't done
				sendWithContext(ctx, errorChan, err)
			}
		}

		r.logger.Trace("listS3FlowLogEvents", "worker_id", workerID,
			"message", "Completed processing object",
			"key", key, "lines_processed", lineNum)
		gr.Close()
		objOut.Body.Close()

		// Check context after finishing an object
		if ctx.Err() != nil {
			r.logger.Debug("listS3FlowLogEvents", "worker_id", workerID,
				"message", "Context done after processing object", "key", key, "reason", ctx.Err())
			return
		}
	}
}

// processTimeTarget processes a full day and adds matching objects to the object pool
func (r *S3FlowLogEventsRetriever) processTimeTarget(
	ctx context.Context,
	date time.Time,
	objectPool *ObjectPool[s3types.Object],
	errorChan chan<- error,
	objectCount *int32,
	processedCount *int32,
) {
	// Build pattern for direct day targeting
	// Structure is typically: base_prefix/YYYY/MM/DD/AccountID_vpcflowlogs_region_*_YYYYMMDD*
	dateStr := fmt.Sprintf("%d%02d%02d", date.Year(), date.Month(), date.Day())

	// Create a list of prefixes to check - we need to try different patterns since
	// VPC flow logs can be stored in different directory structures
	var prefixesToCheck []string

	// Standard format: prefix/YYYY/MM/DD/
	basePrefix := r.prefix
	if !strings.HasSuffix(basePrefix, "/") {
		basePrefix += "/"
	}
	standardPrefix := fmt.Sprintf("%s%d/%02d/%02d/", basePrefix, date.Year(), date.Month(), date.Day())
	prefixesToCheck = append(prefixesToCheck, standardPrefix)

	// Region-specific format: prefix/region/YYYY/MM/DD/
	if strings.Contains(basePrefix, "vpcflowlogs/") {
		regionPrefix := strings.Replace(basePrefix, "vpcflowlogs/", fmt.Sprintf("vpcflowlogs/%s/", r.region), 1)
		regionDatePrefix := fmt.Sprintf("%s%d/%02d/%02d/", regionPrefix, date.Year(), date.Month(), date.Day())
		prefixesToCheck = append(prefixesToCheck, regionDatePrefix)
	}

	dayObjectCount := 0

	// Process each prefix pattern
	for _, datePrefix := range prefixesToCheck {
		r.logger.Debug("listS3FlowLogEvents", "message", "Directly targeting full day",
			"date", dateStr, "prefix", datePrefix)

		// List objects for this prefix
		paginator := s3.NewListObjectsV2Paginator(r.s3Client, &s3.ListObjectsV2Input{
			Bucket: aws.String(r.bucket),
			Prefix: aws.String(datePrefix),
			// Use a small page size to get first results quickly
			MaxKeys: aws.Int32(DefaultS3PageSize),
		})

		// Process each page for this prefix
		for paginator.HasMorePages() {
			// Check context before fetching next page
			if ctx.Err() != nil {
				r.logger.Debug("listS3FlowLogEvents", "message", "Context done during day processing",
					"date", dateStr, "reason", ctx.Err())
				return
			}

			page, err := paginator.NextPage(ctx)
			if err != nil {
				r.logger.Error("listS3FlowLogEvents", "message", "Error listing objects for day",
					"date", dateStr, "error", err)
				if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
					return
				}
				sendWithContext(ctx, errorChan, err)
				return
			}

			// Filter objects for this day
			var dayObjects []s3types.Object
			for _, obj := range page.Contents {
				// Check context frequently to ensure we stop promptly when time is up
				if ctx.Err() != nil {
					r.logger.Debug("listS3FlowLogEvents", "message", "Context done during object filtering",
						"date", dateStr, "reason", ctx.Err())
					return
				}

				key := aws.ToString(obj.Key)
				if !strings.HasSuffix(key, ".log.gz") {
					continue
				}

				// Extract timestamp from the key
				if matches := reTs.FindStringSubmatch(key); len(matches) > 1 {
					tsPart := matches[1]

					// Check if this file belongs to our target day
					if strings.HasPrefix(tsPart, dateStr) {
						fileTs, err := time.Parse("20060102T1504", tsPart)
						if err != nil {
							r.logger.Warn("listS3FlowLogEvents", "message", "Failed to parse timestamp from key",
								"key", key, "error", err)
							continue
						}

						// Double-check that file is within our time bounds
						// Start time is inclusive, end time is exclusive
						if (r.startTime != nil && fileTs.Before(*r.startTime)) ||
							(r.endTime != nil && fileTs.Compare(*r.endTime) >= 0) {
							// Skip incrementing objectCount for objects that don't match time filters
							continue
						}

						dayObjects = append(dayObjects, obj)
						dayObjectCount++
					}
				}
			}

			// Process matched objects for this day
			// Send objects directly to processing
			for i := 0; i < len(dayObjects); i++ {
				// Check context before each object to ensure timely termination
				if ctx.Err() != nil {
					r.logger.Debug("listS3FlowLogEvents", "message", "Context done during object processing",
						"date", dateStr, "reason", ctx.Err(), "processed", i, "total", len(dayObjects))
					return
				}

				obj := dayObjects[i]
				key := aws.ToString(obj.Key)

				// Atomically increment counts
				atomic.AddInt32(objectCount, 1)
				atomic.AddInt32(processedCount, 1)

				r.logger.Trace("listS3FlowLogEvents", "message", "Processing object from day",
					"date", dateStr, "key", key)

				// Add object to processing pool with context awareness
				if !objectPool.AddWithContext(ctx, obj) {
					r.logger.Debug("listS3FlowLogEvents", "message", "Context done when adding to object pool",
						"date", dateStr, "key", key)
					return
				}
			}
		}
	}

	r.logger.Debug("listS3FlowLogEvents", "message", "Completed processing day",
		"date", dateStr, "objects", dayObjectCount)
}

// processS3Objects handles S3 object discovery and adds objects to the worker pool
func (r *S3FlowLogEventsRetriever) processS3Objects(
	ctx context.Context,
	objectPool *ObjectPool[s3types.Object],
	errorChan chan<- error,
	extractionTime int,
) {
	var objectCount int32 = 0
	var processedCount int32 = 0

	// Create a channel to signal early completion based on timeout
	extractionDone := make(chan struct{})

	// Set up a timer to track extraction time separately from context
	// This provides an early warning when approaching the extraction time limit
	extractionTimer := time.NewTimer(time.Duration(extractionTime) * time.Second * 9 / 10) // 90% of extraction time

	// Start a goroutine to monitor extraction time
	go func() {
		select {
		case <-extractionTimer.C:
			r.logger.Info("listS3FlowLogEvents", "message", "Approaching extraction time limit, signaling early completion",
				"extraction_time_seconds", extractionTime,
				"processed_count", atomic.LoadInt32(&processedCount),
				"found_count", atomic.LoadInt32(&objectCount))
			close(extractionDone)
		case <-ctx.Done():
			// Context already done, no need to do anything
			if !extractionTimer.Stop() {
				<-extractionTimer.C // Drain the channel if timer already fired
			}
		}
	}()

	defer func() {
		// Clean up the timer
		if !extractionTimer.Stop() {
			select {
			case <-extractionTimer.C:
				// Drain the channel if timer has fired
			default:
				// Timer was already stopped or drained
			}
		}

		// Close the object pool to signal completion
		objectPool.Close()
		r.logger.Debug("listS3FlowLogEvents", "message", "S3 object listing complete",
			"eligible_objects_found", atomic.LoadInt32(&objectCount), "processed_objects", atomic.LoadInt32(&processedCount))
	}()

	// For time-bounded searches, use progressive processing with direct time slot targeting
	// First, compute date range for our search
	// Default to a reasonable range if not fully specified
	now := time.Now()
	searchStartTime := r.startTime
	searchEndTime := r.endTime
	if searchStartTime == nil {
		// Default to lookback period before end time or current time
		if searchEndTime != nil {
			defaultStart := searchEndTime.Add(-time.Duration(DefaultLookbackHours) * time.Hour)
			searchStartTime = &defaultStart
		} else {
			defaultStart := now.Add(-time.Duration(DefaultLookbackHours) * time.Hour)
			searchStartTime = &defaultStart
		}
	}
	if searchEndTime == nil {
		// Default to now if no end time provided
		searchEndTime = &now
	}

	// Create a structure to track dates we need to process
	type timeTarget struct {
		date time.Time
	}

	// Build a queue of time targets to process, distributed across the time range
	var timeTargets []timeTarget

	// Get start and end date at day granularity
	startDate := time.Date(searchStartTime.Year(), searchStartTime.Month(), searchStartTime.Day(), 0, 0, 0, 0, searchStartTime.Location())
	endDate := time.Date(searchEndTime.Year(), searchEndTime.Month(), searchEndTime.Day(), 0, 0, 0, 0, searchEndTime.Location())

	// Calculate total days in range
	totalDays := int(endDate.Sub(startDate).Hours()/24) + 1
	r.logger.Debug("listS3FlowLogEvents", "message", "Time range calculated",
		"start_date", startDate.Format("2006-01-02"),
		"end_date", endDate.Format("2006-01-02"),
		"total_days", totalDays)

	// Create time targets for each day in the range
	for day := 0; day < totalDays; day++ {
		currentDate := startDate.AddDate(0, 0, day)
		timeTargets = append(timeTargets, timeTarget{date: currentDate})
	}

	// Log the targeting plan
	r.logger.Debug("listS3FlowLogEvents", "message", "Created time targets for distributed sampling",
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
	throttle := make(chan struct{}, DefaultConcurrencyLimit)

	// Create a wait group to track completion of all list operations
	var listWg sync.WaitGroup

	// Create a merged context that also respects our early completion signal
	mergedCtx, cancelMerged := context.WithCancel(ctx)
	defer cancelMerged()

	// Set up goroutine to cancel merged context when extractionDone is signaled
	go func() {
		select {
		case <-extractionDone:
			cancelMerged() // Cancel mergedCtx when extractionDone is signaled
		case <-ctx.Done():
			// Original context is already done, no need to cancel again
		}
	}()

	// Process time targets with regular checks for timeout
	for i, target := range timeTargets {
		// Check context before starting each target to avoid unnecessary work
		select {
		case <-mergedCtx.Done():
			r.logger.Info("listS3FlowLogEvents", "message", "Stopping time target processing early",
				"reason", mergedCtx.Err(), "completed", i, "total", len(timeTargets),
				"processed_count", atomic.LoadInt32(&processedCount))
			goto WaitForCompletion // Skip to wait for already started goroutines
		default:
			// Context still valid, continue processing
		}

		listWg.Add(1)

		// Throttle concurrent API calls
		select {
		case throttle <- struct{}{}:
			// Throttle acquired, proceed
		case <-mergedCtx.Done():
			// Context done while waiting for throttle
			listWg.Done() // Undo the Add(1) since we're not starting the goroutine
			r.logger.Info("listS3FlowLogEvents", "message", "Context cancelled while throttling",
				"reason", mergedCtx.Err(), "completed", i, "total", len(timeTargets))
			goto WaitForCompletion // Skip to wait for already started goroutines
		}

		// Start a goroutine to process this time target
		go func(date time.Time) {
			defer listWg.Done()
			defer func() { <-throttle }() // Release throttle when done

			r.processTimeTarget(mergedCtx, date, objectPool, errorChan, &objectCount, &processedCount)
		}(target.date)
	}

WaitForCompletion:
	// Wait for all listing operations to complete or context to be canceled
	listingDone := make(chan struct{})
	go func() {
		listWg.Wait()
		close(listingDone)
	}()

	// Wait for either completion or cancellation
	select {
	case <-listingDone:
		r.logger.Debug("listS3FlowLogEvents", "message", "All time slots processed",
			"total_objects", atomic.LoadInt32(&objectCount))
	case <-extractionDone:
		r.logger.Info("listS3FlowLogEvents", "message", "Extraction time approaching limit, returning partial results",
			"extraction_time_seconds", extractionTime, "processed_objects", atomic.LoadInt32(&processedCount))
	case <-ctx.Done():
		logContextError(ctx, r.logger, extractionTime, int(atomic.LoadInt32(&processedCount)))
	}
}
