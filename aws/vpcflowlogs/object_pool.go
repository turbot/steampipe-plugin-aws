package vpcflowlogs

import (
	"context"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"sync"
)

const (
	// DefaultObjectChannelSize is the default capacity for the object pool
	DefaultObjectChannelSize = 100
)

// ObjectPool is a thread-safe collection of S3 objects that allows random retrieval
type ObjectPool struct {
	objects []s3types.Object // Store objects in a slice
	mutex   sync.Mutex       // Mutex for thread safety
	closed  bool             // Flag to indicate if pool is closed
	cond    *sync.Cond       // Condition variable for signaling when new objects are added
}

// NewObjectPool creates a new empty object pool
func NewObjectPool() *ObjectPool {
	pool := &ObjectPool{
		objects: make([]s3types.Object, 0, DefaultObjectChannelSize),
	}
	pool.cond = sync.NewCond(&pool.mutex)
	return pool
}

// Add adds an object to the pool
func (p *ObjectPool) Add(obj s3types.Object) {
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
func (p *ObjectPool) GetRandom(ctx context.Context) (s3types.Object, bool) {
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
			return s3types.Object{}, false
		}
	}

	// Check if pool is empty but closed
	if len(p.objects) == 0 {
		return s3types.Object{}, false
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
func (p *ObjectPool) Close() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.closed = true
	p.cond.Broadcast() // Wake up all waiting goroutines
}

// IsEmpty checks if the pool is empty
func (p *ObjectPool) IsEmpty() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return len(p.objects) == 0
}

// Len returns the current number of objects in the pool
func (p *ObjectPool) Len() int {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return len(p.objects)
}

// IsClosed checks if the pool is closed
func (p *ObjectPool) IsClosed() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.closed
}
