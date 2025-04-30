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

// Precompiled regex patterns for better performance
var reCommentPrefix = regexp.MustCompile(`^(#|version\s)`)
var reTs = regexp.MustCompile(`_(\d{8}T\d{4})Z_`)

// S3FlowLogEvent represents a flow log event from S3 storage
type S3FlowLogEvent struct {
	types.FilteredLogEvent
	BucketName string
	S3Key      string
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

	// Constant for controlling concurrency
	const maxConcurrentDownloads = 5
	const resultsChannelSize = 1000
	const objectChannelSize = 100 // Buffer for eligible objects

	// Setup channels for concurrent processing - create a pipeline
	objChan := make(chan s3types.Object, objectChannelSize)
	resultsChan := make(chan S3FlowLogEvent, resultsChannelSize)
	errorChan := make(chan error, maxConcurrentDownloads)
	listingDoneChan := make(chan struct{})    // Signal that S3 listing is complete
	processingDoneChan := make(chan struct{}) // Signal that processing is complete

	// Use the extraction time passed as a parameter
	r.logger.Debug("listS3FlowLogEvents", "extraction_time", extractionTime)

	// Create a context with timeout based on extraction time
	ctx, cancel := context.WithTimeout(ctx, time.Duration(extractionTime)*time.Second)
	defer cancel()

	// Start workers to process objects concurrently
	// Workers begin processing immediately as objects become available from objChan
	r.logger.Debug("listS3FlowLogEvents", "workers", maxConcurrentDownloads,
		"message", "Starting worker goroutines")
	var workersWg sync.WaitGroup
	for i := 0; i < maxConcurrentDownloads; i++ {
		workersWg.Add(1)
		go func() {
			defer workersWg.Done()
			r.processObjectsWorker(ctx, i, objChan, resultsChan, errorChan)
		}()
	}

	// Start a goroutine to close results channel when all workers finish
	go func() {
		workersWg.Wait()
		close(resultsChan)
		close(processingDoneChan)
	}()

	// Start parallel goroutine for progressive processing with direct time slot targeting
	// This implementation starts processing files immediately while continuing to discover more
	r.logger.Debug("listS3FlowLogEvents", "message", "Starting progressive S3 object processing")
	go r.processS3Objects(ctx, objChan, errorChan, listingDoneChan, extractionTime)

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
		// First priority: check if context is done at the start of each loop iteration
		select {
		case <-doneChan:
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				r.logger.Info("listS3FlowLogEvents", "message", "Operation timed out after reaching extraction_time limit, returning partial results",
					"extraction_time_seconds", extractionTime, "results_found", resultCount)
				return nil
			}
			r.logger.Info("listS3FlowLogEvents", "message", "Operation cancelled by context",
				"error", ctx.Err())
			return ctx.Err()
		default:
			// Continue with normal processing
		}

		// Second priority: process events and check other channels
		select {
		case <-doneChan:
			// Double-check for cancellation (prioritized)
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				r.logger.Info("listS3FlowLogEvents", "message", "Operation timed out after reaching extraction_time limit, returning partial results",
					"extraction_time_seconds", extractionTime, "results_found", resultCount)
				return nil
			}
			r.logger.Info("listS3FlowLogEvents", "message", "Operation cancelled by context",
				"error", ctx.Err())
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
			select {
			case <-doneChan:
				// Context was cancelled while trying to stream
				r.logger.Info("listS3FlowLogEvents", "message", "Context cancelled during streaming, returning partial results",
					"results_found", resultCount)
				return nil
			default:
				r.itemStreamer(ctx, result)
			}

		case err := <-errorChan:
			// Cancel all in-progress work
			r.logger.Error("listS3FlowLogEvents", "message", "Error received from worker",
				"error", err)
			cancel()
			return err

		case <-listingDoneChan:
			// S3 listing is complete, but we continue processing existing items
			r.logger.Debug("listS3FlowLogEvents", "message", "S3 listing complete, continuing to process results")

		case <-processingDoneChan:
			// All workers finished and no more results
			r.logger.Debug("listS3FlowLogEvents", "message", "All processing complete, returning results",
				"total_results", resultCount)
			return nil
		case <-ctx.Done():
			// Context cancelled externally or timeout
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				r.logger.Info("listS3FlowLogEvents", "message", "Operation timed out after reaching extraction_time limit, returning partial results",
					"extraction_time_seconds", extractionTime, "results_found", resultCount)
				return nil
			}
			r.logger.Info("listS3FlowLogEvents", "message", "Operation cancelled by context",
				"error", ctx.Err())
			return ctx.Err()
		}
	}
}

// processObjectsWorker processes S3 objects from the object channel and sends flow log events to the results channel
func (r *S3FlowLogEventsRetriever) processObjectsWorker(
	ctx context.Context,
	workerID int,
	objChan <-chan s3types.Object,
	resultsChan chan<- S3FlowLogEvent,
	errorChan chan<- error,
) {
	r.logger.Trace("listS3FlowLogEvents", "worker_id", workerID,
		"message", "Worker started")
	for obj := range objChan {
		key := aws.ToString(obj.Key)
		r.logger.Trace("listS3FlowLogEvents", "worker_id", workerID,
			"message", "Processing object", "key", key)

		// download object
		objOut, err := r.s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(r.bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			r.logger.Error("listS3FlowLogEvents", "worker_id", workerID,
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
			r.logger.Error("listS3FlowLogEvents", "worker_id", workerID,
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

		hasFilters := len(r.filters) > 0

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

			select {
			case <-ctx.Done():
				gr.Close()
				objOut.Body.Close()
				return
			case resultsChan <- ev:
			}
			lineNum++
		}

		if err := scanner.Err(); err != nil {
			r.logger.Error("listS3FlowLogEvents", "worker_id", workerID,
				"message", "Scanner error while processing file",
				"key", key, "error", err)
			select {
			case errorChan <- err:
			case <-ctx.Done():
			}
		}

		r.logger.Trace("listS3FlowLogEvents", "worker_id", workerID,
			"message", "Completed processing object",
			"key", key, "lines_processed", lineNum)
		gr.Close()
		objOut.Body.Close()
	}
	r.logger.Trace("listS3FlowLogEvents", "worker_id", workerID,
		"message", "Worker finished")
}

// processTimeTarget processes a specific time slot and sends matching objects to the object channel
func (r *S3FlowLogEventsRetriever) processTimeTarget(
	ctx context.Context,
	date time.Time,
	hour int,
	objChan chan<- s3types.Object,
	errorChan chan<- error,
	objectCount *int32,
	processedCount *int32,
) {
	// Format the hour for the pattern-matching prefix
	hourStr := fmt.Sprintf("%02d", hour)

	// Build pattern for direct time slot targeting
	// Structure is typically: base_prefix/YYYY/MM/DD/AccountID_vpcflowlogs_region_*_YYYYMMDDTHHMM*
	dateStr := fmt.Sprintf("%d%02d%02d", date.Year(), date.Month(), date.Day())
	timePattern := dateStr + "T" + hourStr

	// Build the date component of the prefix
	datePrefix := r.prefix
	if !strings.HasSuffix(datePrefix, "/") {
		datePrefix += "/"
	}

	// Add region between vpcflowlogs prefix and date components
	if strings.Contains(datePrefix, "vpcflowlogs/") && !strings.Contains(datePrefix, fmt.Sprintf("vpcflowlogs/%s/", r.region)) {
		// Replace "vpcflowlogs/" with "vpcflowlogs/region/"
		datePrefix = strings.Replace(datePrefix, "vpcflowlogs/", fmt.Sprintf("vpcflowlogs/%s/", r.region), 1)
	}

	// If the prefix doesn't already include date components, add them
	if !strings.Contains(datePrefix, fmt.Sprintf("/%d/%02d/%02d/", date.Year(), date.Month(), date.Day())) {
		datePrefix += fmt.Sprintf("%d/%02d/%02d/", date.Year(), date.Month(), date.Day())
	}

	r.logger.Debug("listS3FlowLogEvents", "message", "Directly targeting time slot",
		"date", dateStr, "hour", hourStr, "prefix", datePrefix)

	// List objects for this specific time slot
	paginator := s3.NewListObjectsV2Paginator(r.s3Client, &s3.ListObjectsV2Input{
		Bucket: aws.String(r.bucket),
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
			r.logger.Error("listS3FlowLogEvents", "message", "Error listing objects for time slot",
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
						r.logger.Warn("listS3FlowLogEvents", "message", "Failed to parse timestamp from key",
							"key", key, "error", err)
						continue
					}

					// Double-check that file is within our time bounds
					if (r.startTime != nil && fileTs.Before(*r.startTime)) ||
						(r.endTime != nil && fileTs.After(*r.endTime)) {
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
			atomic.AddInt32(objectCount, 1)
			atomic.AddInt32(processedCount, 1)

			r.logger.Trace("listS3FlowLogEvents", "message", "Processing object from time slot",
				"date", dateStr, "hour", hourStr, "key", key)

			// Send object to processing channel
			select {
			case objChan <- obj:
			case <-ctx.Done():
				return
			}
		}
	}

	r.logger.Debug("listS3FlowLogEvents", "message", "Completed processing time slot",
		"date", dateStr, "hour", hourStr, "objects", timeSlotCount)
}

// processS3Objects handles S3 object discovery and sends objects to the worker pool
func (r *S3FlowLogEventsRetriever) processS3Objects(
	ctx context.Context,
	objChan chan<- s3types.Object,
	errorChan chan<- error,
	listingDoneChan chan<- struct{},
	extractionTime int,
) {
	var objectCount int32 = 0
	var processedCount int32 = 0
	defer func() {
		// Close channels to signal completion
		close(objChan)
		close(listingDoneChan)
		r.logger.Debug("listS3FlowLogEvents", "message", "S3 object listing complete",
			"eligible_objects_found", atomic.LoadInt32(&objectCount), "processed_objects", atomic.LoadInt32(&processedCount))
	}()

	// If we have no time bounds, use default listing
	if r.startTime == nil && r.endTime == nil {
		// Prepare prefix with region component
		listPrefix := r.prefix
		if !strings.HasSuffix(listPrefix, "/") {
			listPrefix += "/"
		}

		// Add region between vpcflowlogs prefix and any trailing components
		if strings.Contains(listPrefix, "vpcflowlogs/") && !strings.Contains(listPrefix, fmt.Sprintf("vpcflowlogs/%s/", r.region)) {
			// Replace "vpcflowlogs/" with "vpcflowlogs/region/"
			listPrefix = strings.Replace(listPrefix, "vpcflowlogs/", fmt.Sprintf("vpcflowlogs/%s/", r.region), 1)
		}

		r.logger.Debug("listS3FlowLogEvents", "message", "Listing S3 objects", "prefix", listPrefix)

		// Use standard listing for cases with no time filtering
		paginator := s3.NewListObjectsV2Paginator(r.s3Client, &s3.ListObjectsV2Input{
			Bucket: aws.String(r.bucket),
			Prefix: aws.String(listPrefix),
		})

		pageCount := 0
		for paginator.HasMorePages() {
			pageCount++
			r.logger.Trace("listS3FlowLogEvents", "message", "Retrieving page of S3 objects",
				"page", pageCount)

			page, err := paginator.NextPage(ctx)
			if err != nil {
				r.logger.Error("listS3FlowLogEvents", "message", "Error listing S3 objects",
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
				atomic.AddInt32(&objectCount, 1)
				atomic.AddInt32(&processedCount, 1)
				r.logger.Trace("listS3FlowLogEvents", "message", "Found eligible object",
					"key", key, "count", atomic.LoadInt32(&objectCount))
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
	searchStartTime := r.startTime
	searchEndTime := r.endTime
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
	r.logger.Debug("listS3FlowLogEvents", "message", "Time range calculated",
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

			r.processTimeTarget(ctx, date, hour, objChan, errorChan, &objectCount, &processedCount)
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
		r.logger.Debug("listS3FlowLogEvents", "message", "All time slots processed",
			"total_objects", atomic.LoadInt32(&objectCount))
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			r.logger.Info("listS3FlowLogEvents", "message", "Time-distributed processing stopped due to timeout",
				"extraction_time_seconds", extractionTime, "objects_processed", atomic.LoadInt32(&processedCount))
		} else {
			r.logger.Info("listS3FlowLogEvents", "message", "Time-distributed processing cancelled",
				"objects_processed", atomic.LoadInt32(&processedCount))
		}
	}
}
