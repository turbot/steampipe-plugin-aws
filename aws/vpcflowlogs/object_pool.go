package vpcflowlogs

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

const (
	// DefaultObjectPoolCapacity is the default capacity for the object pool
	DefaultObjectPoolCapacity = 100
)

// ObjectPool is a thread-safe collection that allows random access to stored objects.
// Objects are added to the pool by producers and randomly retrieved by consumers.
// The pool supports waiting with context for objects to become available.
//
// Example usage:
//
//	// Create a pool for string objects
//	pool := NewObjectPoolDefault[string]()
//
//	// Add objects to the pool
//	pool.Add("item1")
//	pool.Add("item2")
//
//	// Get a random object with context
//	ctx := context.Background()
//	item, ok := pool.GetRandom(ctx)
//	if ok {
//	    // Process the item
//	}
//
//	// When finished, close the pool to signal no more items will be added
//	pool.Close()
type ObjectPool[T any] struct {
	objects []T        // Store objects in a slice
	mutex   sync.Mutex // Mutex for thread safety
	closed  bool       // Flag to indicate if pool is closed
	cond    *sync.Cond // Condition variable for signaling when new objects are added
}

// NewObjectPool creates a new empty object pool with the specified initial capacity.
// The capacity parameter is only a hint for initial memory allocation.
func NewObjectPool[T any](capacity int) *ObjectPool[T] {
	pool := &ObjectPool[T]{
		objects: make([]T, 0, capacity),
	}
	pool.cond = sync.NewCond(&pool.mutex)
	return pool
}

// NewObjectPoolDefault creates a new empty object pool with a default capacity.
// This is the recommended way to create a pool when no specific capacity is needed.
func NewObjectPoolDefault[T any]() *ObjectPool[T] {
	return NewObjectPool[T](DefaultObjectPoolCapacity)
}

// Add puts an object into the pool for later retrieval.
// Returns true if the object was added, false if the pool is closed.
//
// This method is thread-safe and can be called from multiple goroutines.
func (p *ObjectPool[T]) Add(obj T) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.closed {
		return false // Don't add to a closed pool
	}

	p.objects = append(p.objects, obj)
	p.cond.Signal() // Signal that a new object is available
	return true
}

// AddWithContext puts an object into the pool for later retrieval, but respects
// the provided context for cancellation.
// Returns true if the object was added, false if the pool is closed or context is done.
//
// This method is thread-safe and can be called from multiple goroutines.
func (p *ObjectPool[T]) AddWithContext(ctx context.Context, obj T) bool {
	// First check if context is already done
	if ctx.Err() != nil {
		return false
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.closed {
		return false // Don't add to a closed pool
	}

	// Check context again with lock held
	select {
	case <-ctx.Done():
		return false
	default:
		p.objects = append(p.objects, obj)
		p.cond.Signal() // Signal that a new object is available
		return true
	}
}

// GetRandom retrieves and removes a random object from the pool.
// It blocks until an object is available, the pool is closed, or the context is cancelled.
//
// Returns:
//   - The retrieved object and true if successful
//   - A zero value and false if the pool is closed or the context is cancelled
//
// This method is thread-safe and can be called from multiple goroutines.
func (p *ObjectPool[T]) GetRandom(ctx context.Context) (T, bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Wait until there's an object or the pool is closed
	for len(p.objects) == 0 && !p.closed {
		if !waitWithContext(ctx, p.cond) {
			var zero T
			return zero, false
		}
	}

	// Check if pool is empty but closed
	if len(p.objects) == 0 {
		var zero T
		return zero, false
	}

	// Select and remove a random object from the pool
	obj, idx := p.selectRandomObject()

	// Remove the object from the pool (swap with last element and truncate)
	lastIdx := len(p.objects) - 1
	p.objects[idx] = p.objects[lastIdx]
	p.objects = p.objects[:lastIdx]

	return obj, true
}

// Close marks the pool as closed, preventing new objects from being added.
// All goroutines waiting in GetRandom will be woken up.
// This should be called when no more objects will be added to the pool.
//
// This method is thread-safe and can be called from multiple goroutines.
func (p *ObjectPool[T]) Close() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.closed = true
	p.cond.Broadcast() // Wake up all waiting goroutines
}

// IsEmpty checks if the pool currently has no objects.
//
// Note: This is a non-blocking method that returns a point-in-time snapshot.
// In concurrent environments, the pool's state may change immediately after this call.
func (p *ObjectPool[T]) IsEmpty() bool {
	return len(p.objects) == 0
}

// Len returns the current number of objects in the pool.
//
// Note: This is a non-blocking method that returns a point-in-time snapshot.
// In concurrent environments, the pool's state may change immediately after this call.
func (p *ObjectPool[T]) Len() int {
	return len(p.objects)
}

// IsClosed checks if the pool has been marked as closed.
//
// Note: This is a non-blocking method that returns a point-in-time snapshot.
// In concurrent environments, the pool's state may change immediately after this call.
func (p *ObjectPool[T]) IsClosed() bool {
	return p.closed
}

// ----------------------------------------------------------------------------
// Private helper methods
// ----------------------------------------------------------------------------

// waitWithContext waits on a condition variable with context support.
// Returns true if the wait completed normally, false if context was cancelled.
//
// Note: The mutex must be locked when calling this function.
func waitWithContext(ctx context.Context, cond *sync.Cond) bool {
	// Check context first
	if ctx.Err() != nil {
		return false
	}

	// Set up context monitoring
	done := make(chan struct{})
	waiting := true

	go func() {
		select {
		case <-ctx.Done():
			// Context cancelled, acquire lock and signal condition
			cond.L.Lock()
			// Only signal if the waiter is still waiting
			if waiting {
				cond.Signal()
			}
			cond.L.Unlock()
		case <-done:
			// Wait completed normally
		}
	}()

	// Wait for the condition
	cond.Wait()

	// Mark that we're done waiting (to avoid unnecessary signal)
	waiting = false
	close(done)

	// Check if context was cancelled during wait
	return ctx.Err() == nil
}

// Initialize the random number generator with a time-based seed
var randSource = rand.New(rand.NewSource(time.Now().UnixNano()))
var randMutex sync.Mutex // Protect the non-thread-safe rand

// selectRandomObject selects a random object from the pool.
// Returns the selected object and its index.
//
// Note: This method assumes the caller holds the mutex lock
// and that the pool has at least one object.
func (p *ObjectPool[T]) selectRandomObject() (T, int) {
	// Get a random index
	idx := 0
	if len(p.objects) > 1 {
		// Use thread-safe access to the random number generator
		randMutex.Lock()
		idx = randSource.Intn(len(p.objects))
		randMutex.Unlock()
	}

	// Return the object and its index
	return p.objects[idx], idx
}
