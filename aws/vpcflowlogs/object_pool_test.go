package vpcflowlogs

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestObjectPoolInitialization tests the pool initialization
func TestObjectPoolInitialization(t *testing.T) {
	t.Run("Default Capacity", func(t *testing.T) {
		pool := NewObjectPoolDefault[string]()
		if pool == nil {
			t.Fatal("Expected pool to be initialized")
		}
		if pool.Len() != 0 {
			t.Errorf("Expected empty pool, got %d items", pool.Len())
		}
		if pool.IsClosed() {
			t.Error("New pool should not be closed")
		}
	})

	t.Run("Custom Capacity", func(t *testing.T) {
		pool := NewObjectPool[int](50)
		if pool == nil {
			t.Fatal("Expected pool to be initialized")
		}
		if pool.Len() != 0 {
			t.Errorf("Expected empty pool, got %d items", pool.Len())
		}
	})
}

// TestObjectPoolBasicOperations tests basic add and get operations
func TestObjectPoolBasicOperations(t *testing.T) {
	t.Run("Add and Get Single Item", func(t *testing.T) {
		pool := NewObjectPoolDefault[string]()
		pool.Add("test-item")

		if pool.Len() != 1 {
			t.Errorf("Expected 1 item, got %d items", pool.Len())
		}

		item, ok := pool.GetRandom(context.Background())
		if !ok {
			t.Fatal("Failed to get item from pool")
		}
		if item != "test-item" {
			t.Errorf("Expected 'test-item', got '%s'", item)
		}

		if pool.Len() != 0 {
			t.Errorf("Expected empty pool after get, got %d items", pool.Len())
		}
	})

	t.Run("Add and Get Multiple Items", func(t *testing.T) {
		pool := NewObjectPoolDefault[int]()
		for i := 0; i < 5; i++ {
			pool.Add(i)
		}

		if pool.Len() != 5 {
			t.Errorf("Expected 5 items, got %d items", pool.Len())
		}

		// Get all items
		for i := 0; i < 5; i++ {
			_, ok := pool.GetRandom(context.Background())
			if !ok {
				t.Fatalf("Failed to get item %d from pool", i)
			}
		}

		if pool.Len() != 0 {
			t.Errorf("Expected empty pool after getting all items, got %d items", pool.Len())
		}
	})

	t.Run("Empty Pool", func(t *testing.T) {
		pool := NewObjectPoolDefault[string]()
		if !pool.IsEmpty() {
			t.Error("New pool should be empty")
		}

		// Create a context with timeout to avoid blocking the test indefinitely
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		// Trying to get from empty pool should not succeed without items or closing
		_, ok := pool.GetRandom(ctx)
		if ok {
			t.Error("Should not be able to get item from empty pool")
		}
	})
}

// TestObjectPoolConcurrency tests concurrent access to the pool
func TestObjectPoolConcurrency(t *testing.T) {
	t.Run("Concurrent Add", func(t *testing.T) {
		pool := NewObjectPoolDefault[int]()
		var wg sync.WaitGroup

		// Spawn 10 goroutines, each adding 10 items
		numGoroutines := 10
		itemsPerGoroutine := 10
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(base int) {
				defer wg.Done()
				for j := 0; j < itemsPerGoroutine; j++ {
					pool.Add(base*itemsPerGoroutine + j)
				}
			}(i)
		}

		wg.Wait()

		if pool.Len() != numGoroutines*itemsPerGoroutine {
			t.Errorf("Expected %d items, got %d items", numGoroutines*itemsPerGoroutine, pool.Len())
		}
	})

	t.Run("Concurrent Get", func(t *testing.T) {
		pool := NewObjectPoolDefault[int]()

		// Add 100 items
		for i := 0; i < 100; i++ {
			pool.Add(i)
		}

		var wg sync.WaitGroup
		results := make(chan int, 100)

		// Spawn 10 goroutines, each getting 10 items
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 10; j++ {
					item, ok := pool.GetRandom(context.Background())
					if ok {
						results <- item
					}
				}
			}()
		}

		wg.Wait()
		close(results)

		// Count how many items were retrieved
		count := 0
		for range results {
			count++
		}

		if count != 100 {
			t.Errorf("Expected 100 items retrieved, got %d items", count)
		}

		if pool.Len() != 0 {
			t.Errorf("Pool should be empty, got %d items", pool.Len())
		}
	})

	t.Run("Producers and Consumers", func(t *testing.T) {
		pool := NewObjectPoolDefault[int]()
		var producerWg sync.WaitGroup
		var consumerWg sync.WaitGroup

		// Counter for retrieved items
		var retrievedItems int32
		expectedItems := 100 // 5 producers * 20 items each

		// Start 5 producers
		for i := 0; i < 5; i++ {
			producerWg.Add(1)
			go func(id int) {
				defer producerWg.Done()
				for j := 0; j < 20; j++ {
					pool.Add(id*1000 + j)
					time.Sleep(time.Millisecond) // Small delay to interleave operations
				}
			}(i)
		}

		// Use a context with timeout for the entire test
		testCtx, testCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer testCancel()

		// Start 5 consumers
		for i := 0; i < 5; i++ {
			consumerWg.Add(1)
			go func() {
				defer consumerWg.Done()
				for {
					// Check if test is done
					select {
					case <-testCtx.Done():
						return
					default:
						// Continue processing
					}

					// Try to get an item with a short timeout
					ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
					_, ok := pool.GetRandom(ctx)
					cancel()

					if ok {
						// Successfully got an item
						atomic.AddInt32(&retrievedItems, 1)
					} else if pool.IsClosed() && pool.IsEmpty() {
						// Pool is closed and empty, we're done
						return
					}
					// Small sleep to reduce contention
					time.Sleep(time.Millisecond)
				}
			}()
		}

		// Wait for producers to finish
		producerWg.Wait()

		// Close the pool to signal no more items
		pool.Close()

		// Wait for consumers to finish or timeout
		done := make(chan struct{})
		go func() {
			consumerWg.Wait()
			close(done)
		}()

		// Wait for either completion or timeout
		select {
		case <-done:
			// Success - all consumers finished
		case <-testCtx.Done():
			t.Fatal("Test timeout")
		}

		// Verify all items were retrieved
		if int(retrievedItems) != expectedItems {
			t.Errorf("Expected %d items retrieved, got %d", expectedItems, retrievedItems)
		}

		if !pool.IsEmpty() {
			t.Errorf("Pool should be empty, but has %d items", pool.Len())
		}
	})
}

// TestObjectPoolContextCancellation tests context cancellation behavior
func TestObjectPoolContextCancellation(t *testing.T) {
	t.Run("Cancelled Context", func(t *testing.T) {
		pool := NewObjectPoolDefault[string]()
		ctx, cancel := context.WithCancel(context.Background())

		// Cancel the context before calling GetRandom
		cancel()

		_, ok := pool.GetRandom(ctx)
		if ok {
			t.Error("GetRandom should return false with cancelled context")
		}
	})

	t.Run("Context Timeout", func(t *testing.T) {
		pool := NewObjectPoolDefault[string]()

		// Create a context with a short timeout
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()

		// GetRandom should return after the timeout
		start := time.Now()
		_, ok := pool.GetRandom(ctx)
		elapsed := time.Since(start)

		if ok {
			t.Error("GetRandom should return false after context timeout")
		}

		if elapsed < 50*time.Millisecond {
			t.Errorf("GetRandom returned too quickly: %v", elapsed)
		}
	})
}

// TestObjectPoolCloseOperations tests the behavior of closing the pool
func TestObjectPoolCloseOperations(t *testing.T) {
	t.Run("Close Empty Pool", func(t *testing.T) {
		pool := NewObjectPoolDefault[int]()
		pool.Close()

		if !pool.IsClosed() {
			t.Error("Pool should be marked as closed")
		}

		// Adding to closed pool should fail
		if pool.Add(1) {
			t.Error("Should not be able to add to closed pool")
		}

		// Getting from empty closed pool should fail
		_, ok := pool.GetRandom(context.Background())
		if ok {
			t.Error("Should not be able to get from empty closed pool")
		}
	})

	t.Run("Close Non-Empty Pool", func(t *testing.T) {
		pool := NewObjectPoolDefault[int]()
		pool.Add(1)
		pool.Add(2)

		pool.Close()

		if !pool.IsClosed() {
			t.Error("Pool should be marked as closed")
		}

		// Adding to closed pool should fail
		if pool.Add(3) {
			t.Error("Should not be able to add to closed pool")
		}

		// Getting existing items should succeed
		_, ok1 := pool.GetRandom(context.Background())
		if !ok1 {
			t.Error("Should be able to get first item from closed but non-empty pool")
		}

		_, ok2 := pool.GetRandom(context.Background())
		if !ok2 {
			t.Error("Should be able to get second item from closed but non-empty pool")
		}

		// Now pool is empty and closed
		_, ok3 := pool.GetRandom(context.Background())
		if ok3 {
			t.Error("Should not be able to get from empty closed pool")
		}
	})

	t.Run("Close Unblocks Waiting Get", func(t *testing.T) {
		pool := NewObjectPoolDefault[int]()

		// Start a goroutine that waits for an item
		done := make(chan bool)
		go func() {
			ctx := context.Background()
			_, ok := pool.GetRandom(ctx)
			// Should return false when pool is closed
			done <- ok
		}()

		// Give the goroutine time to start waiting
		time.Sleep(50 * time.Millisecond)

		// Close the pool, which should unblock the goroutine
		pool.Close()

		// Check the result
		select {
		case ok := <-done:
			if ok {
				t.Error("GetRandom should return false when pool is closed with no items")
			}
		case <-time.After(500 * time.Millisecond):
			t.Fatal("Timeout waiting for GetRandom to unblock after Close")
		}
	})
}

// TestObjectPoolEdgeCases tests various edge cases
func TestObjectPoolEdgeCases(t *testing.T) {
	t.Run("Different Object Types", func(t *testing.T) {
		// Test with struct type
		type TestStruct struct {
			ID   int
			Name string
		}

		structPool := NewObjectPoolDefault[TestStruct]()
		structPool.Add(TestStruct{ID: 1, Name: "Test1"})
		structPool.Add(TestStruct{ID: 2, Name: "Test2"})

		item, ok := structPool.GetRandom(context.Background())
		if !ok {
			t.Fatal("Failed to get struct from pool")
		}

		if item.ID != 1 && item.ID != 2 {
			t.Errorf("Got unexpected item: %+v", item)
		}

		// Test with pointer type
		ptrPool := NewObjectPoolDefault[*TestStruct]()
		ptrPool.Add(&TestStruct{ID: 1, Name: "Test1"})

		ptrItem, ok := ptrPool.GetRandom(context.Background())
		if !ok {
			t.Fatal("Failed to get pointer from pool")
		}

		if ptrItem.ID != 1 {
			t.Errorf("Got unexpected pointer item: %+v", ptrItem)
		}
	})

	t.Run("Random Selection Distribution", func(t *testing.T) {
		// This test verifies that the selection isn't simply taking the first or last item
		pool := NewObjectPoolDefault[int]()

		// Add 1000 items
		for i := 0; i < 1000; i++ {
			pool.Add(i)
		}

		// Get 100 items and count their positions
		firstQuarter := 0
		secondQuarter := 0
		thirdQuarter := 0
		fourthQuarter := 0

		for i := 0; i < 100; i++ {
			item, _ := pool.GetRandom(context.Background())

			if item < 250 {
				firstQuarter++
			} else if item < 500 {
				secondQuarter++
			} else if item < 750 {
				thirdQuarter++
			} else {
				fourthQuarter++
			}
		}

		// We expect some distribution across all quarters, not concentrated in one area
		// This is a simple check, not a rigorous statistical test
		if firstQuarter == 0 || secondQuarter == 0 || thirdQuarter == 0 || fourthQuarter == 0 {
			t.Errorf("Poor random distribution: %d, %d, %d, %d",
				firstQuarter, secondQuarter, thirdQuarter, fourthQuarter)
		}
	})

	t.Run("High Contention", func(t *testing.T) {
		pool := NewObjectPoolDefault[int]()
		var wg sync.WaitGroup

		// Start with some items
		for i := 0; i < 50; i++ {
			pool.Add(i)
		}

		// Start 20 goroutines that constantly add and get
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()

				for j := 0; j < 50; j++ {
					// Add an item
					pool.Add(id*1000 + j)

					// Get an item with a short-timeout context
					ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
					pool.GetRandom(ctx)
					cancel()
				}
			}(i)
		}

		wg.Wait()

		// The pool should have some items but not be empty or full
		// (exact count is non-deterministic due to concurrent access)
		t.Logf("After high contention, pool has %d items", pool.Len())
	})
}
