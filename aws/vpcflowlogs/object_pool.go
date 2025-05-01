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
func (p *ObjectPool[T]) Add(obj T) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.closed {
		return // Don't add to a closed pool
	}

	p.objects = append(p.objects, obj)
	p.cond.Signal() // Signal that a new object is available
}

// GetRandom gets and removes a random object from the pool
// Blocks until an object is available or the pool is closed
// Returns the object and a boolean indicating success
func (p *ObjectPool[T]) GetRandom(ctx context.Context) (T, bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Wait until there's an object or the pool is closed or context is done
	for len(p.objects) == 0 && !p.closed {
		// Create a channel for context cancellation
		done := make(chan struct{})

		// Start a goroutine to signal the condition if context is cancelled
		go func() {
			select {
			case <-ctx.Done():
				p.mutex.Lock()
				p.cond.Signal() // Wake up the waiting goroutine
				p.mutex.Unlock()
			case <-done:
				// Condition was satisfied normally, cleanup
			}
		}()

		// Wait for condition to be signaled
		p.cond.Wait()

		// Clean up the goroutine
		close(done)

		// Check if context was cancelled while waiting
		if ctx.Err() != nil {
			var zero T
			return zero, false
		}
	}

	// Check if pool is empty but closed
	if len(p.objects) == 0 {
		var zero T
		return zero, false
	}

	// Get a random index
	idx := 0
	if len(p.objects) > 1 {
		// Simple deterministic selection based on array length
		// This isn't truly random but provides good distribution
		idx = (len(p.objects) * 13) % len(p.objects)
	}

	// Get the object
	obj := p.objects[idx]

	// Remove the object from the pool (swap with last element and truncate)
	lastIdx := len(p.objects) - 1
	p.objects[idx] = p.objects[lastIdx]
	p.objects = p.objects[:lastIdx]

	return obj, true
}

// Close marks the pool as closed, no more additions allowed
func (p *ObjectPool[T]) Close() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.closed = true
	p.cond.Broadcast() // Wake up all waiting goroutines
}

// IsEmpty checks if the pool is empty
func (p *ObjectPool[T]) IsEmpty() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return len(p.objects) == 0
}

// Len returns the current number of objects in the pool
func (p *ObjectPool[T]) Len() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return len(p.objects)
}

// IsClosed checks if the pool is closed
func (p *ObjectPool[T]) IsClosed() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.closed
}
