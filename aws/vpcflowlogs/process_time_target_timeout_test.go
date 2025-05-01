package vpcflowlogs

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

// TestProcessTimeTargetWithContextTimeout tests processTimeTarget behavior when context timeout occurs
func TestProcessTimeTargetWithContextTimeout(t *testing.T) {
	// Create a logger that will output debug info during tests
	logger := hclog.New(&hclog.LoggerOptions{
		Level: hclog.Debug,
		Name:  "timeout-test",
	})

	// Test data
	bucket := "test-bucket"
	prefix := "vpcflowlogs/"
	region := "us-east-1"

	// Create a specific date to test
	targetDate := time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)

	// Create many mock objects (enough to cause timeout during processing)
	var mockObjects []s3types.Object
	for i := 0; i < 500; i++ {
		mockObjects = append(mockObjects, s3types.Object{
			Key:          aws.String(fmt.Sprintf("vpcflowlogs/2023/05/01/account_vpcflowlogs_fl-%d_20230501T1200Z_%d.log.gz", i, i)),
			LastModified: aws.Time(targetDate.Add(time.Duration(i) * time.Minute)),
			Size:         aws.Int64(1024 * int64(i%10+1)), // Varying sizes
		})
	}

	// Track API calls and timing
	var listObjectCalls int
	var processStartTime time.Time
	listingDelay := 100 * time.Millisecond // Introduce artificial delay to ensure timeout
	handledTimeout := false

	// Mock S3 client that introduces delay to ensure timeout
	mockS3Client := &MockS3Client{
		ListObjectsV2Fn: func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
			// Record start time on first call
			if listObjectCalls == 0 {
				processStartTime = time.Now()
			}

			listObjectCalls++

			// Simulate API call delay
			select {
			case <-time.After(listingDelay):
				// Normal delay path
			case <-ctx.Done():
				handledTimeout = true
				t.Logf("Context cancelled during simulated API delay (call #%d)", listObjectCalls)
				return nil, ctx.Err()
			}

			// Check if context is already done after delay
			if ctx.Err() != nil {
				handledTimeout = true
				t.Logf("Context already done after delay (call #%d): %v", listObjectCalls, ctx.Err())

				// Return the context error directly
				if errors.Is(ctx.Err(), context.DeadlineExceeded) {
					return nil, context.DeadlineExceeded
				}
				return nil, ctx.Err()
			}

			// Return a chunk of objects for each page
			pageSize := 50
			startIdx := (listObjectCalls - 1) * pageSize
			endIdx := startIdx + pageSize

			if startIdx < len(mockObjects) {
				if endIdx > len(mockObjects) {
					endIdx = len(mockObjects)
				}

				t.Logf("Returning %d objects (from index %d to %d)", endIdx-startIdx, startIdx, endIdx-1)

				return &s3.ListObjectsV2Output{
					Contents:    mockObjects[startIdx:endIdx],
					IsTruncated: aws.Bool(endIdx < len(mockObjects)),
				}, nil
			}

			return &s3.ListObjectsV2Output{
				Contents:    []s3types.Object{},
				IsTruncated: aws.Bool(false),
			}, nil
		},
	}

	// Create a retriever with our mock S3 client
	retriever := NewS3FlowLogEventsRetriever(
		[]string{}, // no filters
		func(ctx context.Context, args ...interface{}) {
			// This is a no-op item streamer for testing
		},
		mockS3Client,
		region,
		bucket,
		prefix,
		nil, // no start time filter
		nil, // no end time filter
		logger,
	)

	// Create test channels and object pool
	errorChan := make(chan error, 5)
	objectPool := NewObjectPoolDefault[s3types.Object]()

	// Create a context with a very short timeout to ensure it times out
	ctx, cancel := context.WithTimeout(context.Background(), 110*time.Millisecond)
	defer cancel()

	// Call the function under test
	t.Log("Starting processTimeTarget with short timeout...")
	retriever.processTimeTarget(ctx, targetDate, objectPool, errorChan)
	t.Logf("processTimeTarget completed, duration: %v", time.Since(processStartTime))

	// The context should have timed out
	assert.Equal(t, context.DeadlineExceeded, ctx.Err(), "Expected context to timeout")
	assert.True(t, handledTimeout, "Expected the S3 client to have handled the timeout")

	// Verify that the function responded to the context timeout
	objectPool.Close() // Close the pool to retrieve all objects

	var pooledObjects []s3types.Object
	maxAttempts := 1000              // Prevent infinite loop
	drainCtx := context.Background() // Use a fresh context for draining the pool

	for i := 0; i < maxAttempts; i++ {
		obj, ok := objectPool.GetRandom(drainCtx)
		if !ok {
			break
		}
		pooledObjects = append(pooledObjects, obj)
	}

	// Some objects should have been processed before timeout
	t.Logf("Objects processed before timeout: %d", len(pooledObjects))

	assert.True(t, len(pooledObjects) > 0, "Expected some objects to be processed before timeout")
	assert.True(t, len(pooledObjects) < len(mockObjects), "Expected less than all objects to be processed due to timeout")

	// No errors should be reported when the context times out
	select {
	case err := <-errorChan:
		t.Fatalf("Unexpected error on timeout: %v", err)
	default:
		// No error, good
	}
}

// TestProcessTimeTargetWithCancellation tests processTimeTarget behavior when context is cancelled
func TestProcessTimeTargetWithCancellation(t *testing.T) {
	// Create a logger that will output debug info during tests
	logger := hclog.New(&hclog.LoggerOptions{
		Level: hclog.Debug,
		Name:  "cancel-test",
	})

	// Test data
	bucket := "test-bucket"
	prefix := "vpcflowlogs/"
	region := "us-east-1"

	// Create a specific date to test
	targetDate := time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)

	// Create many mock objects
	var mockObjects []s3types.Object
	for i := 0; i < 500; i++ {
		mockObjects = append(mockObjects, s3types.Object{
			Key:          aws.String(fmt.Sprintf("vpcflowlogs/2023/05/01/account_vpcflowlogs_fl-%d_20230501T1200Z_%d.log.gz", i, i)),
			LastModified: aws.Time(targetDate.Add(time.Duration(i) * time.Minute)),
			Size:         aws.Int64(1024 * int64(i%10+1)), // Varying sizes
		})
	}

	// Track API calls and timing
	var listObjectCalls int
	listingDelay := 50 * time.Millisecond
	cancelAfterCalls := 2 // Cancel after this many ListObjectsV2 calls
	processCancelled := false

	// Mock S3 client that signals for cancellation after a certain number of calls
	mockS3Client := &MockS3Client{
		ListObjectsV2Fn: func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
			listObjectCalls++
			t.Logf("ListObjectsV2 call #%d", listObjectCalls)

			// Introduce delay to simulate real API call
			time.Sleep(listingDelay)

			// Check if context is already done
			if ctx.Err() != nil {
				processCancelled = true
				t.Logf("Context already cancelled during call #%d: %v", listObjectCalls, ctx.Err())
				if errors.Is(ctx.Err(), context.Canceled) {
					return nil, context.Canceled
				}
				return nil, ctx.Err()
			}

			// Return a chunk of objects for each page
			pageSize := 50
			startIdx := (listObjectCalls - 1) * pageSize
			endIdx := startIdx + pageSize

			if startIdx < len(mockObjects) {
				if endIdx > len(mockObjects) {
					endIdx = len(mockObjects)
				}

				return &s3.ListObjectsV2Output{
					Contents:    mockObjects[startIdx:endIdx],
					IsTruncated: aws.Bool(endIdx < len(mockObjects)),
				}, nil
			}

			return &s3.ListObjectsV2Output{
				Contents:    []s3types.Object{},
				IsTruncated: aws.Bool(false),
			}, nil
		},
	}

	// Create a retriever with our mock S3 client
	retriever := NewS3FlowLogEventsRetriever(
		[]string{}, // no filters
		func(ctx context.Context, args ...interface{}) {
			// This is a no-op item streamer for testing
		},
		mockS3Client,
		region,
		bucket,
		prefix,
		nil, // no start time filter
		nil, // no end time filter
		logger,
	)

	// Create test channels and object pool
	errorChan := make(chan error, 5)
	objectPool := NewObjectPoolDefault[s3types.Object]()

	// Create a context and cancel it after a specific number of API calls
	ctx, cancel := context.WithCancel(context.Background())

	// Start a goroutine to cancel the context after specified number of calls
	go func() {
		for {
			if listObjectCalls >= cancelAfterCalls {
				t.Logf("Cancelling context after %d ListObjectsV2 calls", listObjectCalls)
				cancel()
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Call the function under test
	t.Log("Starting processTimeTarget with planned cancellation...")
	retriever.processTimeTarget(ctx, targetDate, objectPool, errorChan)

	// The context should have been cancelled
	assert.Equal(t, context.Canceled, ctx.Err(), "Expected context to be cancelled")
	assert.True(t, processCancelled, "Expected the S3 client to have handled the cancellation")

	// Verify that the function responded to the context cancellation
	objectPool.Close() // Close the pool to retrieve all objects

	var pooledObjects []s3types.Object
	maxAttempts := 1000              // Prevent infinite loop
	drainCtx := context.Background() // Use a fresh context for draining the pool

	for i := 0; i < maxAttempts; i++ {
		obj, ok := objectPool.GetRandom(drainCtx)
		if !ok {
			break
		}
		pooledObjects = append(pooledObjects, obj)
	}

	// Some objects should have been processed before cancellation
	t.Logf("Objects processed before cancellation: %d", len(pooledObjects))

	assert.True(t, len(pooledObjects) < len(mockObjects), "Expected less than all objects to be processed due to cancellation")

	// No errors should be reported when the context is cancelled
	select {
	case err := <-errorChan:
		t.Fatalf("Unexpected error on cancellation: %v", err)
	default:
		// No error, good
	}
}
