package vpcflowlogs

import (
	"context"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

// TestProcessTimeTargetHappyPath tests the happy path of processTimeTarget
func TestProcessTimeTargetHappyPath(t *testing.T) {
	// Create a logger that won't output anything during tests
	logger := hclog.New(&hclog.LoggerOptions{
		Level: hclog.Off,
	})

	// Test data
	bucket := "test-bucket"
	prefix := "vpcflowlogs/"
	region := "us-east-1"

	// Create a specific date to test
	targetDate := time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)

	// Expected S3 prefixes for this date
	expectedPrefixes := []string{
		"vpcflowlogs/2023/05/01/",           // YYYY/MM/DD format
		"vpcflowlogs/us-east-1/2023/05/01/", // region/YYYY/MM/DD format
	}

	// Create mock objects that should match the filters
	mockObjects := []s3types.Object{
		{
			Key:          aws.String("vpcflowlogs/2023/05/01/account_vpcflowlogs_fl-1234_20230501T1200Z_abc123.log.gz"),
			LastModified: aws.Time(targetDate.Add(12 * time.Hour)), // noon on the target date
			Size:         aws.Int64(1024),
		},
		{
			Key:          aws.String("vpcflowlogs/us-east-1/2023/05/01/account_vpcflowlogs_us-east-1_fl-5678_20230501T1800Z_def456.log.gz"),
			LastModified: aws.Time(targetDate.Add(18 * time.Hour)), // 6 PM on the target date
			Size:         aws.Int64(2048),
		},
	}

	// Create objects that should NOT match (wrong date or format)
	nonMatchingObjects := []s3types.Object{
		{
			Key:          aws.String("vpcflowlogs/2023/05/02/account_vpcflowlogs_fl-1234_20230502T1200Z_abc123.log.gz"), // Wrong date
			LastModified: aws.Time(targetDate.Add(36 * time.Hour)),
			Size:         aws.Int64(1024),
		},
		{
			Key:          aws.String("vpcflowlogs/2023/05/01/account_vpcflowlogs_fl-1234_20230501T1200Z_abc123.txt"), // Wrong extension
			LastModified: aws.Time(targetDate.Add(12 * time.Hour)),
			Size:         aws.Int64(1024),
		},
	}

	// Track which prefixes were requested
	requestedPrefixes := make(map[string]bool)
	var listObjectCalls int

	// Mock S3 client that returns our test objects
	mockS3Client := &MockS3Client{
		ListObjectsV2Fn: func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
			assert.Equal(t, bucket, *params.Bucket)

			// Track which prefix was requested
			requestedPrefixes[*params.Prefix] = true
			listObjectCalls++

			// Determine which objects to return based on the prefix
			var matchingObjects []s3types.Object
			prefix := *params.Prefix

			for _, obj := range mockObjects {
				key := *obj.Key

				// Return objects whose keys start with the requested prefix
				if strings.HasPrefix(key, prefix) {
					matchingObjects = append(matchingObjects, obj)
				}
			}

			// Add non-matching objects to ensure filtering works
			if len(matchingObjects) > 0 {
				matchingObjects = append(matchingObjects, nonMatchingObjects...)
			}

			return &s3.ListObjectsV2Output{
				Contents:    matchingObjects,
				IsTruncated: aws.Bool(false), // No pagination for simplicity
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

	// Counter variables for the function
	var objectCount int32
	var processedCount int32

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Call the function under test
	retriever.processTimeTarget(ctx, targetDate, objectPool, errorChan, &objectCount, &processedCount)

	// Verify that all expected prefixes were requested
	for _, expectedPrefix := range expectedPrefixes {
		assert.True(t, requestedPrefixes[expectedPrefix], "Expected prefix %s was not requested", expectedPrefix)
	}

	// Verify ListObjectsV2 was called the expected number of times (once for each prefix)
	assert.Equal(t, len(expectedPrefixes), listObjectCalls, "ListObjectsV2 should be called once per prefix")

	// Verify object counts
	assert.Equal(t, int32(2), atomic.LoadInt32(&objectCount), "Expected 2 objects to be found")
	assert.Equal(t, int32(2), atomic.LoadInt32(&processedCount), "Expected 2 objects to be processed")

	// Verify objects were added to the pool
	objectPool.Close() // Close the pool to retrieve all objects

	var pooledObjects []s3types.Object
	maxAttempts := 100 // Prevent infinite loop
	for i := 0; i < maxAttempts; i++ {
		obj, ok := objectPool.GetRandom(ctx)
		if !ok {
			break
		}
		pooledObjects = append(pooledObjects, obj)
	}

	// Verify we got exactly the matching objects in the pool
	assert.Equal(t, 2, len(pooledObjects), "Expected 2 objects in the pool")

	// Check that the objects in the pool match our expected objects
	keyMap := make(map[string]bool)
	for _, obj := range pooledObjects {
		keyMap[*obj.Key] = true
	}

	// Verify each expected object is in the pool
	for _, expectedObj := range mockObjects {
		assert.True(t, keyMap[*expectedObj.Key], "Expected object %s should be in the pool", *expectedObj.Key)
	}

	// Verify no errors were reported
	select {
	case err := <-errorChan:
		t.Fatalf("Unexpected error: %v", err)
	default:
		// No error, good
	}
}

// TestProcessTimeTargetWithPagination tests processTimeTarget with pagination
func TestProcessTimeTargetWithPagination(t *testing.T) {
	// Create a logger that won't output anything during tests
	logger := hclog.New(&hclog.LoggerOptions{
		Level: hclog.Off,
	})

	// Test data
	bucket := "test-bucket"
	prefix := "vpcflowlogs/"
	region := "us-east-1"

	// Create a specific date to test
	targetDate := time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)

	// Expected S3 prefixes for this date - both standard format and region-specific format
	expectedPrefixes := []string{
		"vpcflowlogs/2023/05/01/",           // Standard format
		"vpcflowlogs/us-east-1/2023/05/01/", // Region-specific format
	}

	// Create many mock objects (more than a single page)
	var mockObjects []s3types.Object
	for i := 0; i < 50; i++ {
		mockObjects = append(mockObjects, s3types.Object{
			Key:          aws.String("vpcflowlogs/2023/05/01/account_vpcflowlogs_fl-" + string(rune('A'+i)) + "_20230501T1200Z_" + string(rune('a'+i)) + ".log.gz"),
			LastModified: aws.Time(targetDate.Add(time.Duration(i) * time.Hour)),
			Size:         aws.Int64(1024),
		})
	}

	// Track the pagination state
	var continuationToken *string
	var listObjectCalls int
	var prefixesCalled = make(map[string]bool)

	// Mock S3 client that returns our test objects with pagination
	mockS3Client := &MockS3Client{
		ListObjectsV2Fn: func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
			assert.Equal(t, bucket, *params.Bucket)
			prefix := *params.Prefix

			prefixesCalled[prefix] = true
			listObjectCalls++

			// Only return objects for the standard prefix format, empty for region format
			if prefix == "vpcflowlogs/us-east-1/2023/05/01/" {
				return &s3.ListObjectsV2Output{
					Contents:    []s3types.Object{},
					IsTruncated: aws.Bool(false),
				}, nil
			}

			// Only handle pagination for the standard prefix
			if prefix != "vpcflowlogs/2023/05/01/" {
				return &s3.ListObjectsV2Output{
					Contents:    []s3types.Object{},
					IsTruncated: aws.Bool(false),
				}, nil
			}

			// Track the continuation token
			assert.Equal(t, continuationToken, params.ContinuationToken,
				"Continuation token should match what was provided in the previous response")

			// Determine which page of results to return
			startIdx := 0
			if continuationToken != nil {
				// For the second page, start from the middle
				startIdx = 25
			}

			// Return first or second half of the objects
			var pageObjects []s3types.Object
			if startIdx+25 <= len(mockObjects) {
				pageObjects = mockObjects[startIdx : startIdx+25]
			} else {
				pageObjects = mockObjects[startIdx:]
			}

			// For the first call, indicate there are more results and provide a continuation token
			var isTruncated bool
			var nextToken *string
			if continuationToken == nil {
				isTruncated = true
				token := "next-page-token"
				nextToken = &token
				continuationToken = nextToken // save for the next call
			} else {
				// For the second call, indicate no more results
				isTruncated = false
				nextToken = nil
				continuationToken = nil // reset for test clarity
			}

			return &s3.ListObjectsV2Output{
				Contents:              pageObjects,
				IsTruncated:           aws.Bool(isTruncated),
				NextContinuationToken: nextToken,
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

	// Counter variables for the function
	var objectCount int32
	var processedCount int32

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Call the function under test
	retriever.processTimeTarget(ctx, targetDate, objectPool, errorChan, &objectCount, &processedCount)

	// Check that all expected prefixes were called
	for _, prefix := range expectedPrefixes {
		assert.True(t, prefixesCalled[prefix], "Expected prefix %s should have been called", prefix)
	}

	// Verify ListObjectsV2 was called the expected number of times
	// We expect at least 3 calls total: both prefixes (2 calls) plus at least one pagination call
	assert.GreaterOrEqual(t, listObjectCalls, 3, "ListObjectsV2 should be called at least 3 times for pagination")

	// Verify object counts
	assert.Equal(t, int32(50), atomic.LoadInt32(&objectCount), "Expected 50 objects to be found")
	assert.Equal(t, int32(50), atomic.LoadInt32(&processedCount), "Expected 50 objects to be processed")

	// Verify objects were added to the pool
	objectPool.Close() // Close the pool to retrieve all objects

	var pooledObjects []s3types.Object
	maxAttempts := 100 // Prevent infinite loop
	for i := 0; i < maxAttempts; i++ {
		obj, ok := objectPool.GetRandom(ctx)
		if !ok {
			break
		}
		pooledObjects = append(pooledObjects, obj)
	}

	// Verify we got all the objects in the pool
	assert.Equal(t, 50, len(pooledObjects), "Expected 50 objects in the pool")

	// Check that all expected keys are in the pool
	keyMap := make(map[string]bool)
	for _, obj := range pooledObjects {
		keyMap[*obj.Key] = true
	}

	// Verify each expected object is in the pool
	for _, expectedObj := range mockObjects {
		assert.True(t, keyMap[*expectedObj.Key], "Expected object %s should be in the pool", *expectedObj.Key)
	}

	// Verify no errors were reported
	select {
	case err := <-errorChan:
		t.Fatalf("Unexpected error: %v", err)
	default:
		// No error, good
	}
}

// TestProcessTimeTargetWithTimeFilters tests processTimeTarget with time-based filtering
func TestProcessTimeTargetWithTimeFilters(t *testing.T) {
	// Create a logger that won't output anything during tests
	logger := hclog.New(&hclog.LoggerOptions{
		Level: hclog.Off,
	})

	// Test data
	bucket := "test-bucket"
	prefix := "vpcflowlogs/"
	region := "us-east-1"

	// Create a specific date to test
	targetDate := time.Date(2023, 5, 1, 0, 0, 0, 0, time.UTC)

	// Create time filters - only accept logs between 6 AM and 6 PM (exclusive)
	// Start time is inclusive, end time is exclusive
	startTime := targetDate.Add(6 * time.Hour) // 6 AM
	endTime := targetDate.Add(18 * time.Hour)  // 6 PM

	// Create test objects with timestamps throughout the day
	// Create file paths with proper timestamp notation in the key
	// The implementation uses regex to extract timestamps from the key name
	mockObjects := []s3types.Object{
		// Objects that should be accepted (within time range)
		{
			Key:          aws.String("vpcflowlogs/2023/05/01/account_vpcflowlogs_fl-1_20230501T0600Z_abc.log.gz"), // 6 AM - should match
			LastModified: aws.Time(targetDate.Add(6 * time.Hour)),
			Size:         aws.Int64(1024),
		},
		{
			Key:          aws.String("vpcflowlogs/2023/05/01/account_vpcflowlogs_fl-2_20230501T1200Z_def.log.gz"), // 12 PM - should match
			LastModified: aws.Time(targetDate.Add(12 * time.Hour)),
			Size:         aws.Int64(1024),
		},
		{
			Key:          aws.String("vpcflowlogs/2023/05/01/account_vpcflowlogs_fl-3_20230501T1800Z_ghi.log.gz"), // 6 PM - exactly at endTime, should NOT match (exclusive)
			LastModified: aws.Time(targetDate.Add(18 * time.Hour)),
			Size:         aws.Int64(1024),
		},

		// Objects that should be rejected (outside time range)
		{
			Key:          aws.String("vpcflowlogs/2023/05/01/account_vpcflowlogs_fl-4_20230501T0500Z_jkl.log.gz"), // 5 AM - too early
			LastModified: aws.Time(targetDate.Add(5 * time.Hour)),
			Size:         aws.Int64(1024),
		},
		{
			Key:          aws.String("vpcflowlogs/2023/05/01/account_vpcflowlogs_fl-5_20230501T1900Z_mno.log.gz"), // 7 PM - too late
			LastModified: aws.Time(targetDate.Add(19 * time.Hour)),
			Size:         aws.Int64(1024),
		},
	}

	// Separate the objects into those that should be filtered and those that should pass
	// With our updated logic, only objects before the endTime should match (exclusive)
	filteredObjects := mockObjects[:2] // First 2 objects match the time filter

	// Mock S3 client that returns our test objects
	mockS3Client := &MockS3Client{
		ListObjectsV2Fn: func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
			assert.Equal(t, bucket, *params.Bucket)

			prefix := *params.Prefix

			// Return all objects for the standard prefix, empty for region-specific prefix
			if prefix == "vpcflowlogs/2023/05/01/" {
				return &s3.ListObjectsV2Output{
					Contents:    mockObjects,
					IsTruncated: aws.Bool(false),
				}, nil
			}

			// Empty results for other prefixes
			return &s3.ListObjectsV2Output{
				Contents:    []s3types.Object{},
				IsTruncated: aws.Bool(false),
			}, nil
		},
	}

	// Create a retriever with our mock S3 client and time filters
	retriever := NewS3FlowLogEventsRetriever(
		[]string{}, // no content filters
		func(ctx context.Context, args ...interface{}) {
			// This is a no-op item streamer for testing
		},
		mockS3Client,
		region,
		bucket,
		prefix,
		&startTime, // start time filter
		&endTime,   // end time filter
		logger,
	)

	// Create test channels and object pool
	errorChan := make(chan error, 5)
	objectPool := NewObjectPoolDefault[s3types.Object]()

	// Counter variables for the function
	var objectCount int32
	var processedCount int32

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Call the function under test
	retriever.processTimeTarget(ctx, targetDate, objectPool, errorChan, &objectCount, &processedCount)

	// Verify object counts
	// Note: The implementation counts only objects within time range for both objectCount and processedCount
	assert.Equal(t, int32(len(filteredObjects)), atomic.LoadInt32(&objectCount), "Expected %d objects to be found", len(filteredObjects))
	assert.Equal(t, int32(len(filteredObjects)), atomic.LoadInt32(&processedCount), "Expected %d objects to be processed (within time range)", len(filteredObjects))

	// Verify objects were added to the pool
	objectPool.Close() // Close the pool to retrieve all objects

	var pooledObjects []s3types.Object
	maxAttempts := 100 // Prevent infinite loop
	for i := 0; i < maxAttempts; i++ {
		obj, ok := objectPool.GetRandom(ctx)
		if !ok {
			break
		}
		pooledObjects = append(pooledObjects, obj)
	}

	// Verify we got exactly the matching objects in the pool
	assert.Equal(t, len(filteredObjects), len(pooledObjects), "Expected %d objects in the pool (within time range)", len(filteredObjects))

	// Check that the objects in the pool match our expected time-filtered objects
	keyMap := make(map[string]bool)
	for _, obj := range pooledObjects {
		keyMap[*obj.Key] = true
	}

	// Verify correct objects were included/excluded
	for _, obj := range filteredObjects {
		assert.True(t, keyMap[*obj.Key], "Expected key %s should be in the pool", *obj.Key)
	}

	// Verify excluded objects
	excludedKeys := []string{
		"vpcflowlogs/2023/05/01/account_vpcflowlogs_fl-3_20230501T1800Z_ghi.log.gz", // 6 PM - exactly at endTime
		"vpcflowlogs/2023/05/01/account_vpcflowlogs_fl-4_20230501T0500Z_jkl.log.gz", // 5 AM - before startTime
		"vpcflowlogs/2023/05/01/account_vpcflowlogs_fl-5_20230501T1900Z_mno.log.gz", // 7 PM - after endTime
	}

	for _, key := range excludedKeys {
		assert.False(t, keyMap[key], "Excluded key %s should not be in the pool", key)
	}

	// Verify no errors were reported
	select {
	case err := <-errorChan:
		t.Fatalf("Unexpected error: %v", err)
	default:
		// No error, good
	}
}

// TestProcessTimeTargetWithRealWorldFileNames tests processTimeTarget with real-world S3 object names
func TestProcessTimeTargetWithRealWorldFileNames(t *testing.T) {
	// Create a logger that won't output anything during tests
	logger := hclog.New(&hclog.LoggerOptions{
		Level: hclog.Off,
	})

	// Test data
	bucket := "test-bucket"
	prefix := "vpcflowlogs/"
	region := "us-west-2"

	// Create a specific date to test (April 29, 2025)
	targetDate := time.Date(2025, 4, 29, 0, 0, 0, 0, time.UTC)

	// Define time boundaries for filtering (8:55 to 9:15 UTC exclusive)
	// Time boundaries: start time is inclusive, end time is exclusive
	// We'll expect the objects in time slots 8:55 (3 objects), 9:00 (3 objects), 9:05 (3 objects), and 9:10 (4 objects)
	// but not the objects from 8:50 or earlier (24 objects) or from 9:15 or later (16 objects)
	startTime := time.Date(2025, 4, 29, 8, 55, 0, 0, time.UTC)
	endTime := time.Date(2025, 4, 29, 9, 15, 0, 0, time.UTC)

	// Create real-world VPC flow log object names with comprehensive examples
	mockObjects := []s3types.Object{
		// Objects in the 8:15 time slot - should be excluded (before startTime)
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0815Z_b6ca9687.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 8, 21, 7, 0, time.UTC)),
			Size:         aws.Int64(728200),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0815Z_ce89f1a4.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 8, 16, 6, 0, time.UTC)),
			Size:         aws.Int64(97000),
		},

		// Objects in the 8:20 time slot - should be excluded (before startTime)
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0820Z_1cb9131d.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 8, 26, 7, 0, time.UTC)),
			Size:         aws.Int64(745600),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0820Z_a970085d.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 8, 21, 6, 0, time.UTC)),
			Size:         aws.Int64(88600),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0820Z_add42d36.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 8, 31, 6, 0, time.UTC)),
			Size:         aws.Int64(185),
		},

		// Objects in the 8:25 time slot - should be excluded (before startTime)
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0825Z_228e1b6b.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 8, 26, 7, 0, time.UTC)),
			Size:         aws.Int64(102200),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0825Z_7019cd3d.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 8, 36, 6, 0, time.UTC)),
			Size:         aws.Int64(185),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0825Z_e05d5fde.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 8, 31, 7, 0, time.UTC)),
			Size:         aws.Int64(714100),
		},

		// Objects in the 8:30 time slot - should be excluded (before startTime)
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0830Z_8a208661.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 8, 31, 7, 0, time.UTC)),
			Size:         aws.Int64(86300),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0830Z_a1598a90.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 8, 36, 7, 0, time.UTC)),
			Size:         aws.Int64(726200),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0830Z_b876325c.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 8, 41, 6, 0, time.UTC)),
			Size:         aws.Int64(201),
		},

		// Objects in the 8:35 time slot - should be excluded (before startTime)
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0835Z_42ab621a.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 8, 36, 7, 0, time.UTC)),
			Size:         aws.Int64(104400),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0835Z_6a1d00d5.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 8, 41, 7, 0, time.UTC)),
			Size:         aws.Int64(704800),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0835Z_6cdbb6fa.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 8, 46, 6, 0, time.UTC)),
			Size:         aws.Int64(184),
		},

		// Objects in the 8:40 time slot - should be excluded (before startTime)
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0840Z_1f206c97.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 8, 41, 6, 0, time.UTC)),
			Size:         aws.Int64(91100),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0840Z_61bdaf6d.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 1, 7, 0, time.UTC)),
			Size:         aws.Int64(320),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0840Z_6a93aaa0.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 8, 46, 7, 0, time.UTC)),
			Size:         aws.Int64(715200),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0840Z_85d8cff0.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 8, 51, 6, 0, time.UTC)),
			Size:         aws.Int64(185),
		},

		// Objects in the 8:45 time slot - should be excluded (before startTime)
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0845Z_22b878c2.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 8, 56, 6, 0, time.UTC)),
			Size:         aws.Int64(199),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0845Z_b18435d1.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 8, 51, 7, 0, time.UTC)),
			Size:         aws.Int64(706100),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0845Z_d124ba88.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 8, 46, 7, 0, time.UTC)),
			Size:         aws.Int64(94200),
		},

		// Objects in the 8:50 time slot - should be excluded (before startTime)
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0850Z_73486c16.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 8, 56, 7, 0, time.UTC)),
			Size:         aws.Int64(716100),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0850Z_9c5a562c.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 8, 51, 6, 0, time.UTC)),
			Size:         aws.Int64(88400),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0850Z_b373e928.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 1, 6, 0, time.UTC)),
			Size:         aws.Int64(550000),
		},

		// Objects in the 8:55 time slot - should be excluded (before startTime)
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0855Z_c96c623a.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 8, 56, 7, 0, time.UTC)),
			Size:         aws.Int64(89700),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0855Z_dcae4449.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 6, 6, 0, time.UTC)),
			Size:         aws.Int64(203),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0855Z_dcf8fd81.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 1, 7, 0, time.UTC)),
			Size:         aws.Int64(721500),
		},

		// Objects in the 9:00 time slot - should be included
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0900Z_422abc34.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 6, 7, 0, time.UTC)),
			Size:         aws.Int64(702400),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0900Z_89263aaa.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 1, 6, 0, time.UTC)),
			Size:         aws.Int64(103800),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0900Z_a7404713.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 11, 6, 0, time.UTC)),
			Size:         aws.Int64(203),
		},

		// Objects in the 9:05 time slot - should be included
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0905Z_1c2bf943.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 16, 6, 0, time.UTC)),
			Size:         aws.Int64(201),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0905Z_3a40e5ce.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 6, 7, 0, time.UTC)),
			Size:         aws.Int64(103200),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0905Z_e059c109.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 11, 7, 0, time.UTC)),
			Size:         aws.Int64(710100),
		},

		// Objects in the 9:10 time slot - should be included
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0910Z_30bd154e.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 21, 6, 0, time.UTC)),
			Size:         aws.Int64(185),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0910Z_4c373ca0.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 16, 7, 0, time.UTC)),
			Size:         aws.Int64(717800),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0910Z_54b24e06.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 11, 6, 0, time.UTC)),
			Size:         aws.Int64(88700),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0910Z_a214b6ec.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 31, 6, 0, time.UTC)),
			Size:         aws.Int64(479),
		},

		// Objects in the 9:15 time slot - should be excluded (at endTime)
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0915Z_2a72f189.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 16, 6, 0, time.UTC)),
			Size:         aws.Int64(98400),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0915Z_832c88e3.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 21, 7, 0, time.UTC)),
			Size:         aws.Int64(711600),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0915Z_a6d8842a.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 26, 7, 0, time.UTC)),
			Size:         aws.Int64(184),
		},

		// Objects in the 9:20 time slot - should be excluded (after endTime)
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0920Z_9d454c22.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 21, 6, 0, time.UTC)),
			Size:         aws.Int64(95100),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0920Z_e2b51af8.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 26, 8, 0, time.UTC)),
			Size:         aws.Int64(711300),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0920Z_ff6db31c.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 31, 6, 0, time.UTC)),
			Size:         aws.Int64(184),
		},

		// Objects in the 9:25 time slot - should be excluded (after endTime)
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0925Z_933d1a10.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 26, 7, 0, time.UTC)),
			Size:         aws.Int64(99900),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0925Z_b06e6c91.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 36, 6, 0, time.UTC)),
			Size:         aws.Int64(185),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0925Z_be9d3bdc.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 31, 7, 0, time.UTC)),
			Size:         aws.Int64(700100),
		},

		// Objects in the 9:30 time slot - should be excluded (after endTime)
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0930Z_8f0dbc3d.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 31, 6, 0, time.UTC)),
			Size:         aws.Int64(106200),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0930Z_e37d89fa.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 36, 7, 0, time.UTC)),
			Size:         aws.Int64(719300),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0930Z_e86f96da.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 41, 7, 0, time.UTC)),
			Size:         aws.Int64(201),
		},

		// Objects in the 9:35 time slot - should be excluded (after endTime)
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0935Z_190ed438.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 46, 6, 0, time.UTC)),
			Size:         aws.Int64(203),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0935Z_932711a0.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 41, 7, 0, time.UTC)),
			Size:         aws.Int64(708900),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0935Z_d4bfa834.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 36, 6, 0, time.UTC)),
			Size:         aws.Int64(101800),
		},

		// Objects in the 9:40 time slot - should be excluded (after endTime)
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0940Z_3c98df71.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 46, 7, 0, time.UTC)),
			Size:         aws.Int64(712500),
		},
		{
			Key:          aws.String("vpcflowlogs/2025/04/29/891377056770_vpcflowlogs_us-west-2_fl-0b57d2f8b1b7bffa1_20250429T0940Z_76dfe88a.log.gz"),
			LastModified: aws.Time(time.Date(2025, 4, 29, 9, 41, 6, 0, time.UTC)),
			Size:         aws.Int64(105500),
		},
	}

	// The objects that should match our filter (8:55 to 9:15 exclusive)
	// Time filtering uses the timestamp in the filename, not the LastModified time

	// Let's get the correct indices by counting the time slots in our array:
	// 2 objects in 8:15
	// 3 objects in 8:20
	// 3 objects in 8:25
	// 3 objects in 8:30
	// 3 objects in 8:35
	// 4 objects in 8:40
	// 3 objects in 8:45
	// 3 objects in 8:50
	// Total: 24 objects before 8:55

	// 3 objects in 8:55 - should be included (at startTime)
	// 3 objects in 9:00 - should be included
	// 3 objects in 9:05 - should be included
	// 4 objects in 9:10 - should be included
	// Total: 13 objects from 8:55-9:10

	// 3 objects in 9:15 - should be excluded (at endTime)
	// 13 objects in 9:20 and later - should be excluded

	// For our time range (8:55 to 9:15 exclusive), we expect 8:55, 9:00, 9:05, and 9:10
	// which is 3 + 3 + 3 + 4 = 13 objects
	expectedMatches := mockObjects[24:37] // Objects from 8:55, 9:00, 9:05, and 9:10 time slots

	// Mock S3 client that returns our test objects
	mockS3Client := &MockS3Client{
		ListObjectsV2Fn: func(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
			assert.Equal(t, bucket, *params.Bucket)
			prefix := *params.Prefix

			// Match expected prefix pattern for the day
			if prefix == "vpcflowlogs/2025/04/29/" || prefix == "vpcflowlogs/us-west-2/2025/04/29/" {
				var filteredObjects []s3types.Object
				for _, obj := range mockObjects {
					if strings.HasPrefix(*obj.Key, prefix) {
						filteredObjects = append(filteredObjects, obj)
					}
				}
				return &s3.ListObjectsV2Output{
					Contents:    filteredObjects,
					IsTruncated: aws.Bool(false),
				}, nil
			}

			// Return empty result for non-matching prefixes
			return &s3.ListObjectsV2Output{
				Contents:    []s3types.Object{},
				IsTruncated: aws.Bool(false),
			}, nil
		},
	}

	// Create a retriever with our mock S3 client and time filters
	retriever := NewS3FlowLogEventsRetriever(
		[]string{}, // no content filters
		func(ctx context.Context, args ...interface{}) {
			// This is a no-op item streamer for testing
		},
		mockS3Client,
		region,
		bucket,
		prefix,
		&startTime, // Time window start (8:40 UTC)
		&endTime,   // Time window end (8:50 UTC)
		logger,
	)

	// Create test channels and object pool
	errorChan := make(chan error, 5)
	objectPool := NewObjectPoolDefault[s3types.Object]()

	// Counter variables for the function
	var objectCount int32
	var processedCount int32

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Call the function under test
	retriever.processTimeTarget(ctx, targetDate, objectPool, errorChan, &objectCount, &processedCount)

	// Verify object counts
	assert.Equal(t, int32(len(expectedMatches)), atomic.LoadInt32(&objectCount), "Expected %d matching objects", len(expectedMatches))
	assert.Equal(t, int32(len(expectedMatches)), atomic.LoadInt32(&processedCount), "Expected %d objects processed", len(expectedMatches))

	// Verify objects were added to the pool
	objectPool.Close() // Close the pool to retrieve all objects

	var pooledObjects []s3types.Object
	maxAttempts := 100 // Prevent infinite loop
	for i := 0; i < maxAttempts; i++ {
		obj, ok := objectPool.GetRandom(ctx)
		if !ok {
			break
		}
		pooledObjects = append(pooledObjects, obj)
	}

	// Verify we got the right number of objects in the pool
	assert.Equal(t, len(expectedMatches), len(pooledObjects), "Expected %d objects in the pool", len(expectedMatches))

	// Build a map of keys to verify objects
	keyMap := make(map[string]bool)
	for _, obj := range pooledObjects {
		keyMap[*obj.Key] = true
	}

	// Verify each expected object is in the pool
	for _, expectedObj := range expectedMatches {
		assert.True(t, keyMap[*expectedObj.Key], "Expected object %s should be in the pool", *expectedObj.Key)
	}

	// Verify non-matching objects are excluded
	// This includes pre-startTime objects (0-23) and objects at or after endTime (37+)
	excludedObjects := append(mockObjects[:24], mockObjects[37:]...)
	for _, obj := range excludedObjects {
		assert.False(t, keyMap[*obj.Key], "Object %s should not be in the pool", *obj.Key)
	}

	// Verify no errors were reported
	select {
	case err := <-errorChan:
		t.Fatalf("Unexpected error: %v", err)
	default:
		// No error, good
	}
}
