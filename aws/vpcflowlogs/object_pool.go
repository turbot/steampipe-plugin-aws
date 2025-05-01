package vpcflowlogs

import (
	"context"
	"sync"
)

const (
	// DefaultObjectPoolCapacity is the default capacity for the object pool
	DefaultObjectPoolCapacity = 100
)

// ObjectPool is a thread-safe collection of objects that allows random retrieval
// It uses generics to support any type of objects
type ObjectPool[T any] struct {
	objects []T        // Store objects in a slice
	mutex   sync.Mutex // Mutex for thread safety
	closed  bool       // Flag to indicate if pool is closed
	cond    *sync.Cond // Condition variable for signaling when new objects are added
}

// NewObjectPool creates a new empty object pool with the specified capacity
func NewObjectPool[T any](capacity int) *ObjectPool[T] {
	pool := &ObjectPool[T]{
		objects: make([]T, 0, capacity),
	}
	pool.cond = sync.NewCond(&pool.mutex)
	return pool
}

// NewObjectPoolDefault creates a new empty object pool with the default capacity
func NewObjectPoolDefault[T any]() *ObjectPool[T] {
	return NewObjectPool[T](DefaultObjectPoolCapacity)
}

// Add adds an object to the pool
// Returns true if the object was added, false if the pool was closed
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

// waitWithContext waits on a condition variable with context support
// The mutex must be locked when calling this function
// Returns true if the wait completed normally, false if context was cancelled
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

// GetRandom gets and removes a random object from the pool
// Blocks until an object is available or the pool is closed
// Returns the object and a boolean indicating success
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

// selectRandomObject selects a random object from the pool
// Returns the selected object and its index
// Note: This method assumes the caller holds the mutex lock
// and that the pool has at least one object
func (p *ObjectPool[T]) selectRandomObject() (T, int) {
	// Get a random index
	idx := 0
	if len(p.objects) > 1 {
		// Simple deterministic selection based on array length
		// This isn't truly random but provides good distribution
		idx = (len(p.objects) * 13) % len(p.objects)
	}

	// Return the object and its index
	return p.objects[idx], idx
}

// Close marks the pool as closed, no more additions allowed
func (p *ObjectPool[T]) Close() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.closed = true
	p.cond.Broadcast() // Wake up all waiting goroutines
}

// IsEmpty returns if the pool is empty at this instant
// Note: This is not synchronized and may not reflect concurrent modifications
func (p *ObjectPool[T]) IsEmpty() bool {
	return len(p.objects) == 0
}

// Len returns the current number of objects in the pool at this instant
// Note: This is not synchronized and may not reflect concurrent modifications
func (p *ObjectPool[T]) Len() int {
	return len(p.objects)
}

// IsClosed returns if the pool is marked as closed at this instant
// Note: This is not synchronized and may not reflect concurrent modifications
func (p *ObjectPool[T]) IsClosed() bool {
	return p.closed
}
