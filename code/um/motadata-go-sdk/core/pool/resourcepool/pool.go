// Package resourcepool provides a high-performance, thread-safe resource pool
// for managing stateful objects like database connections, file handles, or any
// expensive-to-create resources that can be reused.
package resourcepool

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	. "dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"
)

// Pool manages reusable resources in a thread-safe manner using generics.
// It provides lifecycle management with configurable creation, reset, and destruction callbacks.
// The pool automatically handles resource creation up to a maximum limit and provides
// both blocking and non-blocking resource acquisition methods.
type Pool[T any] struct {
	// Core resource storage and management
	resources chan T            // buffered channel storing available resources
	onCreate  func() (T, error) // factory function to create new resources
	onReset   func(T) error     // optional function to reset resources before reuse
	onDestroy func(T)           // optional function to destroy resources

	// Pool configuration and state
	maxSize int          // maximum number of resources the pool can create
	created atomic.Int32 // total number of resources created (atomic counter)
	used    atomic.Int32 // number of resources currently checked out (atomic counter)
	closed  atomic.Bool  // indicates if the pool is closed (atomic flag)
	once    sync.Once    // ensures Close() is called only once

	// Synchronization primitives
	lock      sync.RWMutex   // reader-writer lock for consistent stats and resource creation
	waitGroup sync.WaitGroup // tracks checked-out resources during shutdown
	returned  chan struct{}  // signals when all resources are returned during shutdown
}

// PoolConfig contains configuration parameters for creating a resource pool.
type PoolConfig[T any] struct {
	// Required configuration
	MaxSize  int               // maximum number of resources the pool can create
	OnCreate func() (T, error) // factory function to create new resources (required)

	// Optional lifecycle hooks
	OnReset   func(T) error // function to reset/clean resources before reuse (optional)
	OnDestroy func(T)       // function to destroy resources when removed from pool (optional)
}

// NewResourcePool creates a new object pool for managing stateful resources.
// It validates the configuration and initializes the pool with the specified parameters.
func NewResourcePool[T any](config PoolConfig[T]) (*Pool[T], error) {
	// Validate required configuration parameters
	if config.MaxSize <= 0 {
		return nil, ErrInvalidSize
	}

	if config.OnCreate == nil {
		return nil, ErrFactoryFunctionRequired
	}

	// Initialize the pool with validated configuration
	pool := &Pool[T]{
		resources: make(chan T, config.MaxSize), // buffered channel with maxSize capacity
		onCreate:  config.OnCreate,              // required factory function
		onReset:   config.OnReset,               // optional reset function
		onDestroy: config.OnDestroy,             // optional destroy function
		maxSize:   config.MaxSize,               // maximum pool capacity
		returned:  make(chan struct{}),          // signal channel for shutdown coordination
	}

	return pool, nil
}

// Get acquires a resource from the pool with context support.
// It will return an available resource immediately, create a new one if possible,
// or wait for one to become available (respecting the context timeout/cancellation).
func (resourcePool *Pool[T]) Get(context context.Context) (T, error) {
	// Zero value to return in error cases
	var zero T

	// Check if pool is closed
	if resourcePool.closed.Load() {
		return zero, ErrPoolClosed
	}

	// Try to get a resource with context cancellation support
	select {
	case resource := <-resourcePool.resources:
		// Got an available resource from the pool
		resourcePool.used.Add(1)      // increment in-use counter atomically
		resourcePool.waitGroup.Add(1) // track this resource for shutdown coordination
		return resource, nil

	case <-context.Done():
		// Context was cancelled or timed out
		return zero, context.Err()

	default:
		// No available resources, try to create a new one
		return resourcePool.create(context)
	}
}

// TryGet attempts to get a resource without blocking.
// Returns immediately with either an available resource, a newly created resource,
// or an error if no resource is available and none can be created.
func (resourcePool *Pool[T]) TryGet() (T, error) {
	// Zero value to return in error cases
	var zero T

	// Check if pool is closed
	if resourcePool.closed.Load() {
		return zero, ErrPoolClosed
	}

	// Non-blocking attempt to get a resource
	select {
	case resource := <-resourcePool.resources:
		// Got an available resource from the pool
		resourcePool.used.Add(1)      // increment in-use counter
		resourcePool.waitGroup.Add(1) // track this resource for shutdown
		return resource, nil

	default:
		// No available resources, try to create one (non-blocking)
		return resourcePool.create(nil)
	}
}

// Put returns a resource back to the pool for reuse.
// The resource will be reset (if OnReset is configured) before being made available.
// Returns an error if the resource is invalid or if the pool is closed.
func (resourcePool *Pool[T]) Put(resource T) error {
	// Validate the resource is not a zero value
	if IsZeroValue(resource) {
		return fmt.Errorf("cannot put zero/nil resource back to pool")
	}

	// Atomically decrement in-use counter and validate
	used := resourcePool.used.Add(-1)
	if used < 0 {
		// Revert the decrement - this resource wasn't obtained from this pool
		resourcePool.used.Add(1)
		return fmt.Errorf("resource not obtained from this pool or pool counter inconsistent")
	}

	// Mark resource as returned for shutdown coordination
	defer resourcePool.waitGroup.Done()

	// Check if pool is closed - if so, destroy the resource
	if resourcePool.closed.Load() {
		// Pool is closed, destroy the resource
		resourcePool.safeDestroy(resource)

		// Signal completion if all resources are returned
		if resourcePool.used.Load() == 0 {
			select {
			case resourcePool.returned <- struct{}{}:
			default: // channel might be closed or full
			}
		}

		return ErrPoolClosed
	}

	// Reset the resource for reuse (if OnReset callback is provided)
	if err := resourcePool.reset(resource); err != nil {
		// Reset failed, destroy the resource
		resourcePool.safeDestroy(resource)
		return err
	}

	// Double-check if pool was closed during reset operation
	if resourcePool.closed.Load() {
		resourcePool.safeDestroy(resource)
		return ErrPoolClosed
	}

	// Try to return the resource to the pool (non-blocking)
	select {
	case resourcePool.resources <- resource:
		// Resource successfully returned to pool
		return nil

	default:
		// Pool channel is full - destroy the resource
		resourcePool.safeDestroy(resource)

		// Check if pool was closed during the operation
		if resourcePool.closed.Load() {
			return ErrPoolClosed
		}

		return fmt.Errorf("pool overflow: channel full")
	}
}

// ═══════════════════════════════════════════════════════════════
// POOL LIFECYCLE MANAGEMENT - SHUTDOWN OPERATIONS
// ═══════════════════════════════════════════════════════════════

// Close gracefully shuts down the pool without waiting for checked-out resources.
// Marks the pool as closed and destroys all available resources.
// Uses sync.Once to ensure it can only be called once.
func (resourcePool *Pool[T]) Close() {
	// Ensure Close is called only once
	resourcePool.once.Do(func() {
		// Mark pool as closed to stop accepting new requests
		resourcePool.closed.Store(true)

		// Drain and destroy all remaining resources in the pool
		// This prevents race conditions with concurrent Put operations
		for {
			select {
			case resource := <-resourcePool.resources:
				// Destroy each available resource
				resourcePool.safeDestroy(resource)
			default:
				// Channel is empty, shutdown complete
				return
			}
		}
	})
}

// CloseWithTimeout gracefully shuts down the pool with a timeout.
// Waits for all checked-out resources to be returned before completing shutdown.
// Returns an error if the timeout is exceeded while waiting for resources.
func (resourcePool *Pool[T]) CloseWithTimeout(timeout time.Duration) error {
	// Mark pool as closed to stop accepting new requests
	resourcePool.closed.Store(true)

	// Create timeout timer
	timer := time.NewTimer(timeout)
	defer timer.Stop() // ensure timer cleanup

	// Create channel to signal when all resources are returned
	done := make(chan struct{})

	// Wait for all resources to return in a separate goroutine
	go func() {
		resourcePool.waitGroup.Wait() // wait for all checked-out resources
		close(done)                   // signal completion
	}()

	// Wait for either completion or timeout
	select {
	case <-done:
		// All resources returned, perform final cleanup
		resourcePool.once.Do(func() {
			// Drain and destroy any remaining resources in the channel
			for {
				select {
				case resource := <-resourcePool.resources:
					resourcePool.safeDestroy(resource)
				default:
					// Channel is empty, cleanup complete
					return
				}
			}
		})
		return nil

	case <-timer.C:
		// Timeout exceeded while waiting for resources to return
		return fmt.Errorf("timeout waiting for %d resources to be returned", resourcePool.used.Load())
	}
}

// UsedCount returns the number of resources currently checked out
func (resourcePool *Pool[T]) UsedCount() int {

	return int(resourcePool.used.Load())
}

// HasWorkingResource returns true if there are resources currently checked out
func (resourcePool *Pool[T]) HasWorkingResource() bool {

	return resourcePool.used.Load() > 0
}

// create creates a new object with panic recovery
func (resourcePool *Pool[T]) create(context context.Context) (resource T, err error) {

	var zero T

	var pending bool

	// Panic recovery to ensure consistent state
	defer func() {
		if r := recover(); r != nil {
			resource = zero

			// If we reserved a creation slot but creation panicked, release it
			if pending {
				resourcePool.created.Add(-1)
			}

			// Convert panic to error
			if err2, ok := r.(error); ok {
				err = err2
			} else {
				err = fmt.Errorf("panic in onCreate: %v", r)
			}
		}
	}()

	// Attempt to create a new resource if pool capacity allows
	// Use mutex to ensure atomic creation check and increment
	resourcePool.lock.Lock()

	// Double-check closed state while holding lock
	if resourcePool.closed.Load() {
		resourcePool.lock.Unlock()
		return zero, ErrPoolClosed
	}

	// Check if we can create a new resource (haven't reached capacity)
	created := resourcePool.created.Load()
	if int(created) < resourcePool.maxSize {
		// Reserve a creation slot by incrementing counter while holding lock
		resourcePool.created.Add(1)
		pending = true
		resourcePool.lock.Unlock()

		// Create the resource outside of lock to avoid blocking other operations
		resource, err = resourcePool.onCreate()
		if err != nil {
			// Creation failed, release the reserved slot
			resourcePool.created.Add(-1)
			pending = false
			return zero, err
		}

		// Successfully created, increment in-use counter and track for shutdown
		resourcePool.used.Add(1)
		resourcePool.waitGroup.Add(1)
		return resource, nil
	}

	// Pool is at capacity, release lock and handle based on context
	resourcePool.lock.Unlock()

	// If no context provided, fail immediately
	if context == nil {
		return zero, ErrPoolExhausted
	}

	// Pool at max capacity, wait for an available resource with context cancellation
	select {
	case resource = <-resourcePool.resources:
		// Got an available resource, track its usage
		resourcePool.used.Add(1)
		resourcePool.waitGroup.Add(1)
		return resource, nil

	case <-context.Done():
		// Context was cancelled or timed out
		return zero, context.Err()
	}
}

// safeDestroy safely destroys a resource with proper counter management.
// Handles the OnDestroy callback with panic recovery to prevent pool corruption.
func (resourcePool *Pool[T]) safeDestroy(resource T) {
	// Decrement created counter in a thread-safe manner
	resourcePool.lock.Lock()
	created := resourcePool.created.Load()
	if created > 0 {
		resourcePool.created.Add(-1)
	}
	resourcePool.lock.Unlock()

	// Call destroy callback if configured, with panic recovery
	if resourcePool.onDestroy != nil {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Log panic but don't propagate it to avoid corrupting pool state
					// In production, use proper logging instead of Printf
					fmt.Printf("WARNING: panic in onDestroy: %v\n", r)
				}
			}()
			resourcePool.onDestroy(resource)
		}()
	}
}

// reset safely resets a resource with panic recovery.
// Calls the OnReset callback if configured, returning any errors.
func (resourcePool *Pool[T]) reset(resource T) (err error) {
	// No reset function configured, consider reset successful
	if resourcePool.onReset == nil {
		return nil
	}

	// Panic recovery for reset operations
	defer func() {
		if r := recover(); r != nil {
			// Convert panic to error
			if err2, ok := r.(error); ok {
				err = err2
			} else {
				err = fmt.Errorf("panic in onReset: %v", r)
			}
		}
	}()

	return resourcePool.onReset(resource)
}

// ═══════════════════════════════════════════════════════════════
// POOL STATISTICS AND MONITORING
// ═══════════════════════════════════════════════════════════════

// PoolStats provides comprehensive statistics about the resource pool's state.
// These metrics are useful for monitoring, debugging, and capacity planning.
type PoolStats struct {
	// Resource counts
	Created   int // total number of resources ever created
	Used      int // number of resources currently checked out
	Available int // number of resources available in the pool channel
	MaxSize   int // maximum number of resources the pool can create

	// Pool state
	IsClosed bool // whether the pool has been closed
}

// GetPoolStats returns a consistent snapshot of pool statistics.
// All metrics are captured atomically to provide a coherent view of the pool state.
func (resourcePool *Pool[T]) GetPoolStats() PoolStats {
	// Take a read lock to ensure consistent snapshot across all metrics
	resourcePool.lock.RLock()
	defer resourcePool.lock.RUnlock()

	// Calculate available resources safely (only when pool is open)
	available := 0
	if !resourcePool.closed.Load() {
		// Safe to check channel length when pool is not closed
		available = len(resourcePool.resources)
	}

	// Return comprehensive statistics
	return PoolStats{
		Created:   int(resourcePool.created.Load()), // total resources created
		Used:      int(resourcePool.used.Load()),    // currently checked out
		Available: available,                        // available in channel
		MaxSize:   resourcePool.maxSize,             // pool capacity limit
		IsClosed:  resourcePool.closed.Load(),       // shutdown state
	}
}
