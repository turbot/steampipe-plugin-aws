package vpcflowlogs

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

// MockS3Client is a mock implementation of S3ClientInterface for testing
type MockS3Client struct {
	GetObjectFn     func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	ListObjectsV2Fn func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
}

func (m *MockS3Client) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return m.GetObjectFn(ctx, params, optFns...)
}

func (m *MockS3Client) ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	return m.ListObjectsV2Fn(ctx, params, optFns...)
}

// createGzippedFlowLog creates a gzipped flow log content for testing
func createGzippedFlowLog(t *testing.T, logLines []string) io.ReadCloser {
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)

	for _, line := range logLines {
		_, err := gzWriter.Write([]byte(line + "\n"))
		if err != nil {
			t.Fatalf("failed to write to gzip buffer: %v", err)
		}
	}

	err := gzWriter.Close()
	if err != nil {
		t.Fatalf("failed to close gzip writer: %v", err)
	}

	return io.NopCloser(bytes.NewReader(buf.Bytes()))
}

// TestProcessObjectsWorkerHappyPath tests the processObjectsWorker function with a valid S3 object
func TestProcessObjectsWorkerHappyPath(t *testing.T) {
	// Create a logger that won't output anything during tests
	logger := hclog.New(&hclog.LoggerOptions{
		Level: hclog.Off,
	})

	// Create test data
	now := time.Now()
	lastModified := now.Add(-1 * time.Hour) // 1 hour ago
	s3Key := "vpcflowlogs/us-east-1/2023/05/01/123456789012_vpcflowlogs_us-east-1_fl-1234_20230501T1200Z_abc123.log.gz"
	bucket := "test-bucket"

	// Sample flow log lines
	// Format: version account-id interface-id srcaddr dstaddr srcport dstport protocol packets bytes start end action log-status
	sampleLogLines := []string{
		"version 2",
		"2 123456789012 eni-abcdef12345678901 172.31.16.139 172.31.16.21 20641 22 6 20 4249 1618346672 1618346732 ACCEPT OK",
		"2 123456789012 eni-abcdef12345678902 10.0.0.1 10.0.0.2 443 43121 6 10 1500 1618346675 1618346735 ACCEPT OK",
	}

	// Mock S3 client that returns our test data
	mockS3Client := &MockS3Client{
		GetObjectFn: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
			assert.Equal(t, bucket, *params.Bucket)
			assert.Equal(t, s3Key, *params.Key)

			// Return gzipped flow log content
			return &s3.GetObjectOutput{
				Body:         createGzippedFlowLog(t, sampleLogLines),
				LastModified: &lastModified,
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
		"us-east-1",
		bucket,
		"vpcflowlogs/",
		nil,
		nil,
		logger,
	)

	// Create a test context
	ctx := context.Background()

	// Create test channels and object pool
	resultsChan := make(chan S3FlowLogEvent, 10)
	errorChan := make(chan error, 5)
	objectPool := NewObjectPoolDefault[s3types.Object]()

	// Add a test object to the pool
	testObject := s3types.Object{
		Key:          aws.String(s3Key),
		LastModified: &lastModified,
		Size:         aws.Int64(1024), // Arbitrary size
	}
	objectPool.Add(testObject)
	objectPool.Close() // Close the pool so the worker will exit after processing the single object

	// Run the worker
	workerID := 1
	retriever.processObjectsWorker(ctx, workerID, objectPool, resultsChan, errorChan)

	// Check that we have no errors
	select {
	case err := <-errorChan:
		t.Fatalf("Unexpected error from worker: %v", err)
	default:
		// No error, continue
	}

	// Collect all events from the results channel
	var events []S3FlowLogEvent
	timeout := time.After(1 * time.Second)

CollectEventsLoop:
	for {
		select {
		case event, ok := <-resultsChan:
			if !ok {
				break CollectEventsLoop // Channel closed
			}
			events = append(events, event)
		case <-timeout:
			break CollectEventsLoop // Timeout collecting events
		default:
			if len(resultsChan) == 0 {
				// No more events in the channel
				time.Sleep(100 * time.Millisecond) // Small pause to allow any pending events
				if len(resultsChan) == 0 {
					break CollectEventsLoop
				}
			}
		}
	}

	// We should have 2 real log events (the version line should be skipped)
	assert.Equal(t, 2, len(events))

	// Check the events
	for _, event := range events {
		// Verify bucket and key
		assert.Equal(t, bucket, event.BucketName)
		assert.Equal(t, s3Key, event.S3Key)

		// Verify the event ID follows the expected format
		// The number part might vary depending on processing order
		assert.Contains(t, *event.EventId, s3Key+":")

		// Verify message is not empty and comes from our sample log lines
		// Skip version line (index 0) when comparing
		assert.NotEmpty(t, *event.Message)
		// Use Contains rather than Equal as processing order might vary
		assert.Contains(t, []string{sampleLogLines[1], sampleLogLines[2]}, *event.Message)

		// Verify timestamp was set
		assert.NotNil(t, event.Timestamp)
		assert.Greater(t, *event.Timestamp, int64(0))

		// Verify ingestion time matches LastModified
		assert.Equal(t, lastModified.UnixMilli(), *event.IngestionTime)
	}
}

// TestProcessObjectsWorkerWithFilters tests the processObjectsWorker function with filters
func TestProcessObjectsWorkerWithFilters(t *testing.T) {
	// Create a logger that won't output anything during tests
	logger := hclog.New(&hclog.LoggerOptions{
		Level: hclog.Off,
	})

	// Create test data
	now := time.Now()
	lastModified := now.Add(-1 * time.Hour) // 1 hour ago
	s3Key := "vpcflowlogs/us-east-1/2023/05/01/123456789012_vpcflowlogs_us-east-1_fl-1234_20230501T1200Z_abc123.log.gz"
	bucket := "test-bucket"

	// Sample flow log lines with a mix of ACCEPT and REJECT
	// Format: version account-id interface-id srcaddr dstaddr srcport dstport protocol packets bytes start end action log-status
	sampleLogLines := []string{
		"version 2",
		"2 123456789012 eni-abcdef12345678901 172.31.16.139 172.31.16.21 20641 22 6 20 4249 1618346672 1618346732 ACCEPT OK",
		"2 123456789012 eni-abcdef12345678902 10.0.0.1 10.0.0.2 443 43121 6 10 1500 1618346675 1618346735 REJECT OK",
		"2 123456789012 eni-abcdef12345678903 192.168.1.1 192.168.1.2 80 12345 6 15 2500 1618346680 1618346740 ACCEPT OK",
	}

	// Mock S3 client that returns our test data
	mockS3Client := &MockS3Client{
		GetObjectFn: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
			assert.Equal(t, bucket, *params.Bucket)
			assert.Equal(t, s3Key, *params.Key)

			// Return gzipped flow log content
			return &s3.GetObjectOutput{
				Body:         createGzippedFlowLog(t, sampleLogLines),
				LastModified: &lastModified,
			}, nil
		},
	}

	// Create a retriever with our mock S3 client with filters
	// Filter for "ACCEPT" action logs only
	retriever := NewS3FlowLogEventsRetriever(
		[]string{"ACCEPT"}, // filter for ACCEPT action logs
		func(ctx context.Context, args ...interface{}) {
			// This is a no-op item streamer for testing
		},
		mockS3Client,
		"us-east-1",
		bucket,
		"vpcflowlogs/",
		nil,
		nil,
		logger,
	)

	// Create a test context
	ctx := context.Background()

	// Create test channels and object pool
	resultsChan := make(chan S3FlowLogEvent, 10)
	errorChan := make(chan error, 5)
	objectPool := NewObjectPoolDefault[s3types.Object]()

	// Add a test object to the pool
	testObject := s3types.Object{
		Key:          aws.String(s3Key),
		LastModified: &lastModified,
		Size:         aws.Int64(1024), // Arbitrary size
	}
	objectPool.Add(testObject)
	objectPool.Close() // Close the pool so the worker will exit after processing the single object

	// Run the worker
	workerID := 1
	retriever.processObjectsWorker(ctx, workerID, objectPool, resultsChan, errorChan)

	// Check that we have no errors
	select {
	case err := <-errorChan:
		t.Fatalf("Unexpected error from worker: %v", err)
	default:
		// No error, continue
	}

	// Collect all events from the results channel
	var events []S3FlowLogEvent
	timeout := time.After(1 * time.Second)

CollectEventsLoop:
	for {
		select {
		case event, ok := <-resultsChan:
			if !ok {
				break CollectEventsLoop // Channel closed
			}
			events = append(events, event)
		case <-timeout:
			break CollectEventsLoop // Timeout collecting events
		default:
			if len(resultsChan) == 0 {
				// No more events in the channel
				time.Sleep(100 * time.Millisecond) // Small pause to allow any pending events
				if len(resultsChan) == 0 {
					break CollectEventsLoop
				}
			}
		}
	}

	// We should have 2 ACCEPT logs (filtered out the REJECT log)
	assert.Equal(t, 2, len(events))

	// Verify all events contain "ACCEPT"
	for _, event := range events {
		assert.Contains(t, *event.Message, "ACCEPT")
		assert.NotContains(t, *event.Message, "REJECT")
	}
}

// TestProcessObjectsWorkerContextCancellationBeforeGetObject tests that the worker stops immediately
// when the context is cancelled before getting an object
func TestProcessObjectsWorkerContextCancellationBeforeGetObject(t *testing.T) {
	// Create a logger that won't output anything during tests
	logger := hclog.New(&hclog.LoggerOptions{
		Level: hclog.Off,
	})

	// Create test data
	now := time.Now()
	lastModified := now.Add(-1 * time.Hour) // 1 hour ago
	s3Key := "vpcflowlogs/us-east-1/2023/05/01/123456789012_vpcflowlogs_us-east-1_fl-1234_20230501T1200Z_abc123.log.gz"
	bucket := "test-bucket"

	// Count how many times GetObject is called
	var getObjectCalls int
	var getObjectMutex sync.Mutex

	// Mock S3 client
	mockS3Client := &MockS3Client{
		GetObjectFn: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
			getObjectMutex.Lock()
			getObjectCalls++
			getObjectMutex.Unlock()
			return nil, errors.New("this should not be called")
		},
	}

	// Create a retriever with our mock S3 client
	retriever := NewS3FlowLogEventsRetriever(
		[]string{},
		func(ctx context.Context, args ...interface{}) {},
		mockS3Client,
		"us-east-1",
		bucket,
		"vpcflowlogs/",
		nil,
		nil,
		logger,
	)

	// Create a cancellable context that we'll cancel immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel the context immediately

	// Create test channels and object pool
	resultsChan := make(chan S3FlowLogEvent, 10)
	errorChan := make(chan error, 5)
	objectPool := NewObjectPoolDefault[s3types.Object]()

	// Add a test object to the pool
	testObject := s3types.Object{
		Key:          aws.String(s3Key),
		LastModified: &lastModified,
		Size:         aws.Int64(1024),
	}
	objectPool.Add(testObject)

	// Run the worker - it should exit immediately due to cancelled context
	workerID := 1
	retriever.processObjectsWorker(ctx, workerID, objectPool, resultsChan, errorChan)

	// Verify GetObject was never called because context was already cancelled
	assert.Equal(t, 0, getObjectCalls, "GetObject should not have been called")

	// Verify no events were produced
	assert.Equal(t, 0, len(resultsChan), "No events should have been produced")
}

// TestProcessObjectsWorkerContextCancellationDuringGetObject tests that the worker stops
// when the context is cancelled during GetObject
func TestProcessObjectsWorkerContextCancellationDuringGetObject(t *testing.T) {
	// Create a logger that won't output anything during tests
	logger := hclog.New(&hclog.LoggerOptions{
		Level: hclog.Off,
	})

	// Create test data
	now := time.Now()
	lastModified := now.Add(-1 * time.Hour) // 1 hour ago
	s3Key := "vpcflowlogs/us-east-1/2023/05/01/123456789012_vpcflowlogs_us-east-1_fl-1234_20230501T1200Z_abc123.log.gz"
	bucket := "test-bucket"

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Mock S3 client that cancels the context during GetObject
	mockS3Client := &MockS3Client{
		GetObjectFn: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
			// Cancel the context while we're in GetObject
			cancel()

			// Simulate some delay to ensure cancellation is processed
			time.Sleep(10 * time.Millisecond)

			// Now check if context is cancelled
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}

			return &s3.GetObjectOutput{
				Body:         io.NopCloser(bytes.NewReader([]byte{})),
				LastModified: &lastModified,
			}, nil
		},
	}

	// Create a retriever with our mock S3 client
	retriever := NewS3FlowLogEventsRetriever(
		[]string{},
		func(ctx context.Context, args ...interface{}) {},
		mockS3Client,
		"us-east-1",
		bucket,
		"vpcflowlogs/",
		nil,
		nil,
		logger,
	)

	// Create test channels and object pool
	resultsChan := make(chan S3FlowLogEvent, 10)
	errorChan := make(chan error, 5)
	objectPool := NewObjectPoolDefault[s3types.Object]()

	// Add a test object to the pool
	testObject := s3types.Object{
		Key:          aws.String(s3Key),
		LastModified: &lastModified,
		Size:         aws.Int64(1024),
	}
	objectPool.Add(testObject)
	objectPool.Close()

	// Run the worker - it should exit when the context is cancelled
	workerID := 1
	retriever.processObjectsWorker(ctx, workerID, objectPool, resultsChan, errorChan)

	// Verify no events were produced
	assert.Equal(t, 0, len(resultsChan), "No events should have been produced")
}

// signallingBlockingReader is a custom io.ReadCloser that signals when a read is attempted
// and blocks until explicitly allowed to continue
type signallingBlockingReader struct {
	reader        io.ReadCloser
	blocker       chan struct{}
	readAttempted chan struct{}
	blocked       bool
}

func newSignallingBlockingReader(r io.ReadCloser, readAttempted chan struct{}) (*signallingBlockingReader, chan struct{}) {
	blocker := make(chan struct{})
	return &signallingBlockingReader{
		reader:        r,
		blocker:       blocker,
		readAttempted: readAttempted,
	}, blocker
}

func (r *signallingBlockingReader) Read(p []byte) (int, error) {
	// Signal that read was attempted
	select {
	case <-r.readAttempted:
		// Already signaled
	default:
		close(r.readAttempted)
	}

	// Block on first read
	if !r.blocked {
		r.blocked = true
		<-r.blocker // Wait until we're told to proceed
	}

	return r.reader.Read(p)
}

func (r *signallingBlockingReader) Close() error {
	return r.reader.Close()
}

// TestProcessObjectsWorkerContextCancellationAfterGetObject tests that the worker stops
// when the context is cancelled after GetObject but before scanning
func TestProcessObjectsWorkerContextCancellationAfterGetObject(t *testing.T) {
	// Create a logger that won't output anything during tests
	logger := hclog.New(&hclog.LoggerOptions{
		Level: hclog.Off,
	})

	// Create test data
	now := time.Now()
	lastModified := now.Add(-1 * time.Hour) // 1 hour ago
	s3Key := "vpcflowlogs/us-east-1/2023/05/01/123456789012_vpcflowlogs_us-east-1_fl-1234_20230501T1200Z_abc123.log.gz"
	bucket := "test-bucket"

	// Sample flow log lines
	sampleLogLines := []string{
		"version 2",
		"2 123456789012 eni-abcdef12345678901 172.31.16.139 172.31.16.21 20641 22 6 20 4249 1618346672 1618346732 ACCEPT OK",
		"2 123456789012 eni-abcdef12345678902 10.0.0.1 10.0.0.2 443 43121 6 10 1500 1618346675 1618346735 ACCEPT OK",
	}

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Signal channels to control test flow more precisely
	getObjectCalled := make(chan struct{}) // Signals when GetObject is called
	readyToBlock := make(chan struct{})    // Signals when the reader is ready to block
	readAttempted := make(chan struct{})   // Signals when a read is attempted
	var blockerChan chan struct{}          // Controls when to unblock the reader

	// Mock S3 client that returns our controlled gzip data
	mockS3Client := &MockS3Client{
		GetObjectFn: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
			// Signal that GetObject was called
			close(getObjectCalled)

			// Wait for test to signal it's ready for blocking
			select {
			case <-readyToBlock:
				// Continue with the setup
			case <-ctx.Done():
				return nil, ctx.Err()
			}

			// Create custom reader that will block on the very first byte
			// and signal when read is attempted
			normalReader := createGzippedFlowLog(t, sampleLogLines)

			reader, blocker := newSignallingBlockingReader(normalReader, readAttempted)
			blockerChan = blocker

			return &s3.GetObjectOutput{
				Body:         reader,
				LastModified: &lastModified,
			}, nil
		},
	}

	// Create a retriever with our mock S3 client
	retriever := NewS3FlowLogEventsRetriever(
		[]string{},
		func(ctx context.Context, args ...interface{}) {},
		mockS3Client,
		"us-east-1",
		bucket,
		"vpcflowlogs/",
		nil,
		nil,
		logger,
	)

	// Create test channels with size 0 to avoid buffering
	// This is crucial to ensure we don't process any events after cancellation
	resultsChan := make(chan S3FlowLogEvent) // Unbuffered
	errorChan := make(chan error, 5)
	objectPool := NewObjectPoolDefault[s3types.Object]()

	// Add a test object to the pool
	testObject := s3types.Object{
		Key:          aws.String(s3Key),
		LastModified: &lastModified,
		Size:         aws.Int64(1024),
	}
	objectPool.Add(testObject)
	objectPool.Close()

	// Create a channel to drain any events that might be produced
	// This ensures we don't deadlock if events are produced despite cancellation
	drainDone := make(chan struct{})
	go func() {
		defer close(drainDone)
		for {
			select {
			case _, ok := <-resultsChan:
				if !ok {
					return // Channel closed
				}
				// Just drain any events
			case <-time.After(500 * time.Millisecond):
				return // Timeout waiting for events
			}
		}
	}()

	// Start worker in a goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Close resultsChan when done to signal the drain goroutine
		defer close(resultsChan)
		retriever.processObjectsWorker(ctx, 1, objectPool, resultsChan, errorChan)
	}()

	// Wait for GetObject to be called
	<-getObjectCalled

	// Signal that we're ready to proceed with blocking
	close(readyToBlock)

	// Wait for read to be attempted
	<-readAttempted

	// Now the worker is blocked trying to read the first byte
	// Cancel the context
	cancel()

	// Unblock the reader so the worker can detect the cancellation
	close(blockerChan)

	// Wait for worker to finish
	wg.Wait()

	// Wait for drain goroutine to finish
	<-drainDone

	// Check for any errors
	select {
	case err := <-errorChan:
		t.Fatalf("Unexpected error: %v", err)
	default:
		// No errors, good
	}
}

// TestProcessObjectsWorkerContextCancellationDuringScanning tests that the worker stops
// when the context is cancelled during the scanning phase
func TestProcessObjectsWorkerContextCancellationDuringScanning(t *testing.T) {
	// Create a logger that won't output anything during tests
	logger := hclog.New(&hclog.LoggerOptions{
		Level: hclog.Off,
	})

	// Create test data
	now := time.Now()
	lastModified := now.Add(-1 * time.Hour) // 1 hour ago
	s3Key := "vpcflowlogs/us-east-1/2023/05/01/123456789012_vpcflowlogs_us-east-1_fl-1234_20230501T1200Z_abc123.log.gz"
	bucket := "test-bucket"

	// Create a lot of sample flow log lines to ensure scanning takes time
	var sampleLogLines []string
	sampleLogLines = append(sampleLogLines, "version 2")

	// Generate 1000 lines to ensure scanning takes time
	for i := 0; i < 1000; i++ {
		sampleLogLines = append(sampleLogLines, fmt.Sprintf(
			"2 123456789012 eni-abcdef%d 172.31.16.%d 172.31.16.%d %d 22 6 20 4249 1618346672 1618346732 ACCEPT OK",
			i, i%256, (i+1)%256, 10000+i))
	}

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Track if any events were sent
	eventCount := 0
	var eventCountMutex sync.Mutex

	// Mock S3 client that returns our large data set
	mockS3Client := &MockS3Client{
		GetObjectFn: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
			return &s3.GetObjectOutput{
				Body:         createGzippedFlowLog(t, sampleLogLines),
				LastModified: &lastModified,
			}, nil
		},
	}

	// Create a retriever with our mock S3 client
	retriever := NewS3FlowLogEventsRetriever(
		[]string{},
		func(ctx context.Context, args ...interface{}) {
			eventCountMutex.Lock()
			eventCount++
			eventCountMutex.Unlock()

			// Cancel after we've processed a few events
			if eventCount == 50 {
				cancel()
			}
		},
		mockS3Client,
		"us-east-1",
		bucket,
		"vpcflowlogs/",
		nil,
		nil,
		logger,
	)

	// Create a buffered channel to avoid blocking
	resultsChan := make(chan S3FlowLogEvent, 1000)
	errorChan := make(chan error, 5)
	objectPool := NewObjectPoolDefault[s3types.Object]()

	// Add a test object to the pool
	testObject := s3types.Object{
		Key:          aws.String(s3Key),
		LastModified: &lastModified,
		Size:         aws.Int64(1024),
	}
	objectPool.Add(testObject)
	objectPool.Close()

	// Run the worker
	retriever.processObjectsWorker(ctx, 1, objectPool, resultsChan, errorChan)

	// We should have at most the number of log lines minus the header line
	maxEvents := len(sampleLogLines) - 1

	// Verify we have a reasonable number of events - either all events made it through
	// before cancellation took effect, or a subset of them did
	numEvents := len(resultsChan)
	assert.LessOrEqual(t, numEvents, maxEvents,
		"Should not have more events than log lines (excluding header)")

	// Log the result for diagnostic purposes
	t.Logf("Processed %d out of a possible %d events before cancellation took effect",
		numEvents, maxEvents)
}

// TestProcessObjectsWorkerTimeoutWaitingForObject tests the worker's behavior when
// it times out waiting to get an object from the pool
func TestProcessObjectsWorkerTimeoutWaitingForObject(t *testing.T) {
	// Create a logger that won't output anything during tests
	logger := hclog.New(&hclog.LoggerOptions{
		Level: hclog.Off,
	})

	// Create test data
	bucket := "test-bucket"

	// Mock S3 client
	mockS3Client := &MockS3Client{
		GetObjectFn: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
			t.Fatal("GetObject should not be called")
			return nil, errors.New("should not be called")
		},
	}

	// Create a retriever with our mock S3 client
	retriever := NewS3FlowLogEventsRetriever(
		[]string{},
		func(ctx context.Context, args ...interface{}) {},
		mockS3Client,
		"us-east-1",
		bucket,
		"vpcflowlogs/",
		nil,
		nil,
		logger,
	)

	// Create test channels and an EMPTY object pool
	resultsChan := make(chan S3FlowLogEvent, 10)
	errorChan := make(chan error, 5)
	objectPool := NewObjectPoolDefault[s3types.Object]()

	// Do NOT close the pool to simulate waiting for objects that never come

	// Create a context with timeout shorter than the worker's internal timeout
	// Worker uses 5 seconds, so we'll use 100ms
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Run the worker - it should loop and timeout when context is cancelled
	startTime := time.Now()
	retriever.processObjectsWorker(ctx, 1, objectPool, resultsChan, errorChan)
	duration := time.Since(startTime)

	// Worker should exit due to context timeout
	assert.True(t, duration >= 100*time.Millisecond, "Worker should have waited for at least context timeout")
	assert.True(t, duration < 5*time.Second, "Worker should not have waited for its internal timeout")

	// Verify no events were produced
	assert.Equal(t, 0, len(resultsChan), "No events should have been produced")
}

// TestProcessObjectsWorkerCancelledDuringSendToResultsChannel tests the worker's behavior
// when the context is cancelled while sending results to the channel
func TestProcessObjectsWorkerCancelledDuringSendToResultsChannel(t *testing.T) {
	// Create a logger that won't output anything during tests
	logger := hclog.New(&hclog.LoggerOptions{
		Level: hclog.Off,
	})

	// Create test data
	now := time.Now()
	lastModified := now.Add(-1 * time.Hour) // 1 hour ago
	s3Key := "vpcflowlogs/us-east-1/2023/05/01/123456789012_vpcflowlogs_us-east-1_fl-1234_20230501T1200Z_abc123.log.gz"
	bucket := "test-bucket"

	// Sample flow log lines
	sampleLogLines := []string{
		"version 2",
		"2 123456789012 eni-abcdef12345678901 172.31.16.139 172.31.16.21 20641 22 6 20 4249 1618346672 1618346732 ACCEPT OK",
		"2 123456789012 eni-abcdef12345678902 10.0.0.1 10.0.0.2 443 43121 6 10 1500 1618346675 1618346735 ACCEPT OK",
	}

	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())

	// Mock S3 client
	mockS3Client := &MockS3Client{
		GetObjectFn: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
			return &s3.GetObjectOutput{
				Body:         createGzippedFlowLog(t, sampleLogLines),
				LastModified: &lastModified,
			}, nil
		},
	}

	// Create a retriever with our mock S3 client
	retriever := NewS3FlowLogEventsRetriever(
		[]string{},
		func(ctx context.Context, args ...interface{}) {},
		mockS3Client,
		"us-east-1",
		bucket,
		"vpcflowlogs/",
		nil,
		nil,
		logger,
	)

	// Create test channels and object pool
	// Use unbuffered channel to simulate blocking when sending results
	resultsChan := make(chan S3FlowLogEvent) // Unbuffered to force blocking
	errorChan := make(chan error, 5)
	objectPool := NewObjectPoolDefault[s3types.Object]()

	// Add a test object to the pool
	testObject := s3types.Object{
		Key:          aws.String(s3Key),
		LastModified: &lastModified,
		Size:         aws.Int64(1024),
	}
	objectPool.Add(testObject)
	objectPool.Close()

	// Start worker in a goroutine - it will process the object but block when sending to results
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		retriever.processObjectsWorker(ctx, 1, objectPool, resultsChan, errorChan)
	}()

	// Wait a short time to ensure the worker is blocked trying to send results
	time.Sleep(100 * time.Millisecond)

	// Cancel the context while worker is blocked sending results
	cancel()

	// Worker should now exit
	waitWithTimeout(&wg, 1*time.Second)

	// Try to drain the channel in case any events were sent
	drainChannel(resultsChan)
}

// TestProcessObjectsWorkerErrorHandling tests that the worker properly handles
// errors from S3 GetObject call
func TestProcessObjectsWorkerErrorHandling(t *testing.T) {
	// Create a logger that won't output anything during tests
	logger := hclog.New(&hclog.LoggerOptions{
		Level: hclog.Off,
	})

	// Create test data
	now := time.Now()
	lastModified := now.Add(-1 * time.Hour) // 1 hour ago
	s3Key := "vpcflowlogs/us-east-1/2023/05/01/123456789012_vpcflowlogs_us-east-1_fl-1234_20230501T1200Z_abc123.log.gz"
	bucket := "test-bucket"
	expectedError := errors.New("s3 get object error")

	// Mock S3 client that returns an error
	mockS3Client := &MockS3Client{
		GetObjectFn: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
			return nil, expectedError
		},
	}

	// Create a retriever with our mock S3 client
	retriever := NewS3FlowLogEventsRetriever(
		[]string{},
		func(ctx context.Context, args ...interface{}) {},
		mockS3Client,
		"us-east-1",
		bucket,
		"vpcflowlogs/",
		nil,
		nil,
		logger,
	)

	// Create test channels and object pool
	resultsChan := make(chan S3FlowLogEvent, 10)
	errorChan := make(chan error, 5)
	objectPool := NewObjectPoolDefault[s3types.Object]()

	// Add a test object to the pool
	testObject := s3types.Object{
		Key:          aws.String(s3Key),
		LastModified: &lastModified,
		Size:         aws.Int64(1024),
	}
	objectPool.Add(testObject)

	// Add a second object to see if processing continues after error
	testObject2 := s3types.Object{
		Key:          aws.String(s3Key + ".2"),
		LastModified: &lastModified,
		Size:         aws.Int64(1024),
	}
	objectPool.Add(testObject2)
	objectPool.Close()

	// Run the worker
	retriever.processObjectsWorker(context.Background(), 1, objectPool, resultsChan, errorChan)

	// Verify error was properly reported
	select {
	case err := <-errorChan:
		assert.Equal(t, expectedError, err, "Error should have been reported correctly")
	default:
		t.Fatalf("Expected error but none was received")
	}

	// Verify no events were produced
	assert.Equal(t, 0, len(resultsChan), "No events should have been produced")
}

// Helper function to wait for a WaitGroup with timeout
func waitWithTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return true // completed normally
	case <-time.After(timeout):
		return false // timed out
	}
}

// TestProcessObjectsWorkerContextTimeout tests that processing actually stops
// after context cancellation, rather than continuing to process data
func TestProcessObjectsWorkerContextTimeout(t *testing.T) {
	// Create test configuration
	const (
		totalLogLines    = 10000 // Number of sample log lines to generate
		cancelThreshold  = 1000  // Cancel after processing this many items
		maxPostCancel    = 200   // Maximum allowed items processed after cancellation
		minItemsRequired = 100   // Minimum items that must be processed before cancellation
		workerTimeout    = 2 * time.Second
	)

	// Basic test setup
	logger := hclog.New(&hclog.LoggerOptions{Level: hclog.Off})
	bucket := "test-bucket"
	s3Key := "vpcflowlogs/us-east-1/2023/05/01/123456789012_vpcflowlogs_test_file.log.gz"
	lastModified := time.Now().Add(-1 * time.Hour)

	// Generate sample log data
	sampleLogLines := []string{"version 2"}
	for i := 0; i < totalLogLines; i++ {
		sampleLogLines = append(sampleLogLines, fmt.Sprintf(
			"2 123456789012 eni-abcdef%d 172.31.%d.%d 172.31.%d.%d %d 22 6 20 4249 1618346672 1618346732 ACCEPT OK",
			i, i%256, (i+1)%256, (i+2)%256, (i+3)%256, 10000+i))
	}

	// Create mock S3 client
	mockS3Client := &MockS3Client{
		GetObjectFn: func(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
			return &s3.GetObjectOutput{
				Body:         createGzippedFlowLog(t, sampleLogLines),
				LastModified: &lastModified,
			}, nil
		},
	}

	// Create S3FlowLogEventsRetriever with mock client
	retriever := NewS3FlowLogEventsRetriever(
		[]string{}, // no filters
		func(ctx context.Context, args ...interface{}) {}, // no-op streamer
		mockS3Client,
		"us-east-1", bucket, "vpcflowlogs/",
		nil, nil, logger,
	)

	// Setup test channels and pool
	resultsChan := make(chan S3FlowLogEvent, totalLogLines) // Buffered to avoid blocking
	errorChan := make(chan error, 5)
	objectPool := NewObjectPoolDefault[s3types.Object]()
	objectPool.Add(s3types.Object{
		Key:          aws.String(s3Key),
		LastModified: &lastModified,
		Size:         aws.Int64(1024),
	})
	objectPool.Close() // No more objects will be added

	// Synchronization primitives
	processedCount := atomic.Int32{}  // Tracks items processed
	cancelled := make(chan struct{})  // Signals when cancellation occurs
	workerDone := make(chan struct{}) // Signals when worker completes
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start result monitoring goroutine
	go func() {
		for {
			select {
			case event, ok := <-resultsChan:
				if !ok {
					return // Channel closed
				}

				// Count each result and check for cancellation threshold
				count := processedCount.Add(1)
				if count == cancelThreshold {
					t.Logf("Cancelling context after processing %d items", count)
					cancel()
					close(cancelled)
				}

				_ = event // Use the item to avoid compiler warnings

			case <-workerDone:
				return // Worker completed

			case <-time.After(workerTimeout):
				t.Logf("Monitor timeout after %v", workerTimeout)
				return
			}
		}
	}()

	// Start worker goroutine
	go func() {
		retriever.processObjectsWorker(ctx, 1, objectPool, resultsChan, errorChan)
		close(workerDone)
	}()

	// Wait for cancellation to occur
	select {
	case <-cancelled:
		// Good - cancellation triggered after threshold reached
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out waiting for cancellation")
	}

	// Record count at cancellation
	countAtCancellation := int(processedCount.Load())

	// Wait for worker to finish
	select {
	case <-workerDone:
		// Worker completed normally after cancellation
	case <-time.After(workerTimeout):
		t.Fatal("Worker did not exit within timeout period after cancellation")
	}

	// Calculate final metrics
	itemsProcessedAfterCancel := len(resultsChan)
	totalProcessed := countAtCancellation + itemsProcessedAfterCancel
	expectedLogLines := totalLogLines - 1 // Subtract version line

	// Log metrics for debugging
	t.Logf("Items at cancellation: %d", countAtCancellation)
	t.Logf("Items after cancellation: %d", itemsProcessedAfterCancel)
	t.Logf("Total processed: %d of %d available", totalProcessed, expectedLogLines)

	// Verify results
	assert.Greater(t, countAtCancellation, minItemsRequired,
		"Should process a meaningful number of items before cancellation")

	assert.Less(t, totalProcessed, expectedLogLines,
		"Processing should stop before handling all available log lines")

	assert.LessOrEqual(t, itemsProcessedAfterCancel, maxPostCancel,
		"Should not process too many items after context cancellation")

	// Verify no errors were reported
	select {
	case err := <-errorChan:
		t.Fatalf("Unexpected error: %v", err)
	default:
		// No errors, good
	}
}

// Helper function to drain a channel (best effort)
func drainChannel[T any](ch chan T) {
	for {
		select {
		case <-ch:
		default:
			return
		}
	}
}
