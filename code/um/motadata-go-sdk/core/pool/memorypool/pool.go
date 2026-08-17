// Package memorypool provides high-performance, type-safe memory pools for different data types.
// It uses generics to provide zero-allocation buffer management with automatic expansion,
// leak detection, and comprehensive memory optimization features.
package memorypool

import (
	"fmt"
	"runtime"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/core/types"
	. "dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"
)

// ═══════════════════════════════════════════════════════════════
// TYPE DEFINITIONS AND CONSTANTS
// ═══════════════════════════════════════════════════════════════

// PoolType defines the constraint for data types supported by the memory pool.
// All basic data types that implement the Interface constraint can be used.
// This includes String, Int64, Float64, Byte, and other primitive types.
type PoolType interface{}

// Pool represents a type-safe memory pool using generics for efficient buffer management.
// It maintains multiple pre-allocated buffers of type T and tracks their usage through bitmasks.
// The pool minimizes allocations by reusing buffers and supports dynamic expansion.
type Pool[T PoolType] struct {
	// Core pool data
	buffers [][]T // collection of pre-allocated buffers

	// Pool metadata
	dataType string // string representation of the data type for debugging

	// Usage tracking (bitmask for efficient pool tracking)
	usedPools int64 // bitmask indicating which pools are currently in use (bit n = pool n)

	// Configuration
	poolLength int  // default length of each buffer
	poolSize   int  // total number of buffers in the pool
	expandable bool // whether buffers can grow beyond poolLength
}

// ═══════════════════════════════════════════════════════════════
// POOL INITIALIZATION
// ═══════════════════════════════════════════════════════════════

// newMemoryPool creates a new generic memory pool for type T with specified configuration.
// It determines the data type at compile time for better debugging and monitoring.
// Parameters:
//   - poolSize: number of buffers to maintain in the pool
//   - poolLength: default size of each buffer
//   - expandable: whether buffers can grow beyond poolLength
func newMemoryPool[T PoolType](poolSize, poolLength int, expandable bool) *Pool[T] {
	// Type detection for debugging and monitoring purposes
	var poolType T
	var dataType string

	// Determine the actual type for logging and debugging
	switch any(poolType).(type) {
	case string:
		dataType = types.TypeSTRING
	case int64:
		dataType = types.TypeINT64
	case float64:
		dataType = types.TypeFLOAT64
	case byte:
		dataType = types.TypeByte
	default:
		dataType = types.TypeUnknown
	}

	// Initialize the pool with pre-allocated buffer slice
	return &Pool[T]{
		buffers:    make([][]T, poolSize), // pre-allocate slice for buffer pointers
		usedPools:  0,                     // no pools used initially
		poolLength: poolLength,            // default buffer size
		poolSize:   poolSize,              // total number of buffers
		expandable: expandable,            // expansion policy
		dataType:   dataType,              // type identifier for debugging
	}
}

// ═══════════════════════════════════════════════════════════════
// POOL ACQUISITION AND RELEASE
// ═══════════════════════════════════════════════════════════════

// acquirePool acquires a buffer of the specified size from the pool.
// Returns the buffer index and the slice, or NotAvailable if no buffer is available.
// The acquired buffer is automatically marked as in-use and initialized if needed.
func (pool *Pool[T]) acquirePool(size int) (int, []T) {
	// Validate and normalize the requested size
	size = pool.getSize(size)

	// Search for an available buffer in the pool
	for index := range pool.buffers[:pool.poolSize] {
		if pool.isAvailable(index) {
			// Mark this pool as used in the bitmask
			pool.usedPools |= 1 << index

			// Allocate or resize the buffer as needed
			pool.allocatePool(index, size)

			return index, pool.buffers[index]
		}
	}

	// No available buffers found
	return NotAvailable, nil
}

// releasePool releases the buffer at the specified index back to the pool.
// Includes leak detection to identify buffers that were already released.
// The buffer remains allocated but is marked as available for reuse.
func (pool *Pool[T]) releasePool(index int) {
	// Check if buffer is already available (potential double-release or leak)
	if pool.isAvailable(index) {
		// Capture stack trace for debugging
		stackTraceBytes := make([]byte, 1<<20)

		// Log potential memory leak
		fmt.Printf("%s pool %d not released, reason: %s pool leaked\n", pool.dataType, index, pool.dataType)
		fmt.Printf("Stack trace: %s\n", string(stackTraceBytes[:runtime.Stack(stackTraceBytes, false)]))
		return
	}

	// Mark the pool as available by clearing the corresponding bit
	pool.usedPools &^= 1 << index
}

// ═══════════════════════════════════════════════════════════════
// POOL ACCESS AND INFORMATION
// ═══════════════════════════════════════════════════════════════

// getPool returns the buffer at the specified index.
// This provides direct access to the underlying buffer slice.
// No validation is performed - caller must ensure index is valid.
func (pool *Pool[T]) getPool(index int) []T {
	return pool.buffers[index]
}

// getUsedPools returns the bitmask representing which pools are currently in use.
// Each bit corresponds to a pool index (bit 0 = pool 0, bit 1 = pool 1, etc.)
// This is useful for monitoring pool utilization and debugging.
func (pool *Pool[T]) getUsedPools() int {
	return int(pool.usedPools)
}

// getPoolLength returns the configured default length for each buffer in the pool.
// This is the size buffers will have when initially allocated or after shrinking.
func (pool *Pool[T]) getPoolLength() int {
	return pool.poolLength
}

// getPoolSize returns the total number of buffers managed by this pool.
// This represents the maximum number of concurrent buffer users.
func (pool *Pool[T]) getPoolSize() int {
	return pool.poolSize
}

// ═══════════════════════════════════════════════════════════════
// POOL EXPANSION AND SHRINKING
// ═══════════════════════════════════════════════════════════════

// expandPool expands the buffer at the specified index to the given size.
// This is called when a buffer needs to grow beyond its current capacity.
// Only works if the pool was created with expandable=true.
func (pool *Pool[T]) expandPool(index, size int) []T {
	// Validate and normalize the requested size
	size = pool.getSize(size)

	// Log expansion for monitoring and debugging
	fmt.Printf("%s Pool expanded from %d to %d\n", pool.dataType, len(pool.buffers[index]), size)

	// Create additional zero-value elements
	values := make([]T, size-len(pool.buffers[index]))

	// Extend the existing buffer with new elements
	pool.buffers[index] = append(pool.buffers[index], values...)

	return pool.buffers[index]
}

// shrinkPool shrinks oversized buffers back to their default poolLength.
// This helps reclaim memory from buffers that were temporarily expanded.
// Only operates on available (unused) buffers and only if the pool is expandable.
func (pool *Pool[T]) shrinkPool() {
	// Only shrink if the pool supports expansion
	if pool.expandable {
		// Check each buffer in the pool
		for index := range pool.buffers[:pool.poolSize] {
			// Only shrink buffers that are currently available
			if pool.isAvailable(index) {
				// Shrink buffers that exceed the default length
				if len(pool.buffers[index]) > pool.poolLength {
					fmt.Printf("%s pool shrink from %d to %d\n", pool.dataType, len(pool.buffers[index]), pool.poolLength)
					// Replace with a new buffer of default size
					pool.buffers[index] = make([]T, pool.poolLength)
				}
			}
		}
	}
}

// ═══════════════════════════════════════════════════════════════
// LEAK DETECTION AND DEBUGGING
// ═══════════════════════════════════════════════════════════════

// testPoolLeak detects memory leaks by checking for unreleased pools.
// If any pools are still marked as used, it logs the leak and forcibly resets them.
// This is typically called during cleanup or testing phases.
func (pool *Pool[T]) testPoolLeak() {
	// Check if any pools are still marked as used
	if pool.usedPools != 0 {
		// Force reset all used pools
		pool.usedPools = 0

		// Capture stack trace for debugging
		stackTraceBytes := make([]byte, 1<<20)

		// Log the leak with stack trace
		fmt.Printf("%s pool leaked\n", pool.dataType)
		fmt.Printf("Stack trace: %s\n", string(stackTraceBytes[:runtime.Stack(stackTraceBytes, false)]))
	}
}

// ═══════════════════════════════════════════════════════════════
// PRIVATE HELPER METHODS
// ═══════════════════════════════════════════════════════════════

// isAvailable checks if a pool at the given index is available for use.
// Uses bitmask operations for efficient availability checking.
// Returns true if the pool is available, false if it's in use.
func (pool *Pool[T]) isAvailable(index int) bool {
	// Special case: NotAvailable index is considered available
	if index == NotAvailable {
		return true
	}

	// Check if the corresponding bit is clear (0 = available, 1 = used)
	return pool.usedPools&(1<<index) == 0
}

// getSize validates and normalizes the requested buffer size.
// Applies pool policies regarding expansion and size limits.
// Returns the actual size that should be allocated.
func (pool *Pool[T]) getSize(size int) int {
	// Use default pool length if no specific size requested
	if size == NotAvailable {
		size = pool.poolLength
	} else if size > pool.poolLength && !pool.expandable {
		// Log size limitation and cap to pool length
		fmt.Printf("failed to acquire %s pool, reason: requested size %d > %d\n", pool.dataType, size, pool.poolLength)
		size = pool.poolLength
	}

	return size
}

// allocatePool initializes or resizes a buffer at the specified index.
// Handles both new buffer creation and existing buffer management.
// Ensures all elements are zero-initialized for safety.
func (pool *Pool[T]) allocatePool(index, size int) {
	// Zero value for the type T (used for clearing buffers)
	var value T

	// Handle new buffer allocation
	if pool.buffers[index] == nil {
		if size > pool.poolLength {
			// Log expansion during initial allocation
			fmt.Printf("%s Pool initialized with expanded size %d\n", pool.dataType, size)
			pool.buffers[index] = make([]T, size)
		} else {
			// Create buffer with size but reserve capacity up to poolLength
			pool.buffers[index] = make([]T, size, pool.poolLength)
		}
	} else {
		// Handle existing buffer resize
		if size > len(pool.buffers[index]) {
			// Buffer needs to grow
			if len(pool.buffers[index]) >= pool.poolLength {
				// Log expansion if buffer was already at default size
				fmt.Printf("%s Pool expanded from %d to %d\n", pool.dataType, len(pool.buffers[index]), size)
			}

			// Clear existing values to zero value of type T
			for i := range pool.buffers[index] {
				pool.buffers[index][i] = value
			}

			// Create additional elements and append them
			values := make([]T, size-len(pool.buffers[index]))
			pool.buffers[index] = append(pool.buffers[index], values...)
		} else {
			// Buffer can be shrunk or kept same size
			pool.buffers[index] = pool.buffers[index][:size]

			// Clear all elements to zero value
			for i := range pool.buffers[index][:size] {
				pool.buffers[index][i] = value
			}
		}
	}
}
