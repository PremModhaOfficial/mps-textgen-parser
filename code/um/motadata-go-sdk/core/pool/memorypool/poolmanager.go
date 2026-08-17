package memorypool

import (
	"sync"
)

// ═══════════════════════════════════════════════════════════════
// GLOBAL SINGLETON MANAGEMENT
// ═══════════════════════════════════════════════════════════════

// Global variables for singleton pattern implementation.
// Ensures a single global PoolManager instance when needed.
var (
	once        = sync.Once{} // ensures single initialization
	poolManager *PoolManager  // global singleton instance
)

// ═══════════════════════════════════════════════════════════════
// TYPE DEFINITIONS
// ═══════════════════════════════════════════════════════════════

// PoolManager manages all memory pools in a centralized way.
// It provides type-safe access to pools for different data types
// with optional thread-safety for global usage.
type PoolManager struct {
	// Type-specific memory pools
	stringPool  *Pool[string]  // pool for string slices
	int64Pool   *Pool[int64]   // pool for int64 slices
	float64Pool *Pool[float64] // pool for float64 slices
	bytePool    *Pool[byte]    // pool for byte slices

	// Concurrency control
	lock sync.RWMutex // reader-writer lock for thread-safe operations

	// Configuration
	global bool // indicates if this is a global singleton instance
}

// PoolConfig defines configuration parameters for creating memory pools.
// All type-specific pools will be created with the same configuration.
type PoolConfig struct {
	// Pool dimensions
	PoolSize   int // number of buffers in each pool
	PoolLength int // default size of each buffer

	// Behavior settings
	Expandable bool // whether buffers can grow beyond PoolLength
	Global     bool // whether to use global singleton pattern
}

// PoolStats provides usage statistics for all memory pools.
// Each field represents a bitmask where each bit corresponds to a pool index.
type PoolStats struct {
	// Usage bitmasks for each pool type (bit set = pool in use)
	UsedStringPools  int // bitmask of used string pools
	UsedInt64Pools   int // bitmask of used int64 pools
	UsedFloat64Pools int // bitmask of used float64 pools
	UsedBytePools    int // bitmask of used byte pools
}

// ═══════════════════════════════════════════════════════════════
// INITIALIZATION AND CONFIGURATION
// ═══════════════════════════════════════════════════════════════

// NewPoolManager creates a new PoolManager instance with the specified configuration.
// If Global is true, returns a singleton instance; otherwise creates a new instance.
// The singleton pattern ensures consistent global memory pool access.
func NewPoolManager(config *PoolConfig) *PoolManager {
	// Handle global singleton pattern
	if config.Global {
		// Use sync.Once to ensure thread-safe singleton initialization
		once.Do(func() {
			poolManager = &PoolManager{
				// Initialize all type-specific pools with the same configuration
				stringPool:  newMemoryPool[string](config.PoolSize, config.PoolLength, config.Expandable),
				int64Pool:   newMemoryPool[int64](config.PoolSize, config.PoolLength, config.Expandable),
				float64Pool: newMemoryPool[float64](config.PoolSize, config.PoolLength, config.Expandable),
				bytePool:    newMemoryPool[byte](config.PoolSize, config.PoolLength, config.Expandable),
				lock:        sync.RWMutex{},
				global:      true,
			}
		})
		return poolManager
	}

	// Create a new local (non-global) instance
	return &PoolManager{
		stringPool:  newMemoryPool[string](config.PoolSize, config.PoolLength, config.Expandable),
		int64Pool:   newMemoryPool[int64](config.PoolSize, config.PoolLength, config.Expandable),
		float64Pool: newMemoryPool[float64](config.PoolSize, config.PoolLength, config.Expandable),
		bytePool:    newMemoryPool[byte](config.PoolSize, config.PoolLength, config.Expandable),
		global:      false, // no locking needed for local instances
	}
}

// ═══════════════════════════════════════════════════════════════
// STRING POOL LIFECYCLE METHODS
// ═══════════════════════════════════════════════════════════════

// AcquireStringPool acquires a string buffer from the pool.
// Returns the buffer index and the slice, or NotAvailable if no buffer is available.
// Thread-safe for global instances.
func (poolManager *PoolManager) AcquireStringPool(size int) (int, []string) {
	// Apply thread-safety lock for global instances
	if poolManager.global {
		poolManager.lock.Lock()
		defer poolManager.lock.Unlock()
	}

	return poolManager.stringPool.acquirePool(size)
}

// ReleaseStringPool returns a string buffer back to the pool.
// The buffer at the specified index will be marked as available for reuse.
// Thread-safe for global instances.
func (poolManager *PoolManager) ReleaseStringPool(index int) {
	// Apply thread-safety lock for global instances
	if poolManager.global {
		poolManager.lock.Lock()
		defer poolManager.lock.Unlock()
	}

	poolManager.stringPool.releasePool(index)
}

// GetStringPool retrieves the string buffer at the specified index.
// This provides direct access to the underlying slice for reading/writing.
// Thread-safe for global instances (uses read lock).
func (poolManager *PoolManager) GetStringPool(index int) []string {
	// Apply read lock for global instances (allows concurrent reads)
	if poolManager.global {
		poolManager.lock.RLock()
		defer poolManager.lock.RUnlock()
	}

	return poolManager.stringPool.getPool(index)
}

// ═══════════════════════════════════════════════════════════════
// INT64 POOL LIFECYCLE METHODS
// ═══════════════════════════════════════════════════════════════

// AcquireInt64Pool acquires an int64 buffer from the pool.
// Returns the buffer index and the slice, or NotAvailable if no buffer is available.
// Thread-safe for global instances.
func (poolManager *PoolManager) AcquireInt64Pool(size int) (int, []int64) {
	// Apply thread-safety lock for global instances
	if poolManager.global {
		poolManager.lock.Lock()
		defer poolManager.lock.Unlock()
	}

	return poolManager.int64Pool.acquirePool(size)
}

// ReleaseInt64Pool returns an int64 buffer back to the pool.
// The buffer at the specified index will be marked as available for reuse.
// Thread-safe for global instances.
func (poolManager *PoolManager) ReleaseInt64Pool(index int) {
	// Apply thread-safety lock for global instances
	if poolManager.global {
		poolManager.lock.Lock()
		defer poolManager.lock.Unlock()
	}

	poolManager.int64Pool.releasePool(index)
}

// GetInt64Pool retrieves the int64 buffer at the specified index.
// This provides direct access to the underlying slice for reading/writing.
// Thread-safe for global instances (uses read lock).
func (poolManager *PoolManager) GetInt64Pool(index int) []int64 {
	// Apply read lock for global instances (allows concurrent reads)
	if poolManager.global {
		poolManager.lock.RLock()
		defer poolManager.lock.RUnlock()
	}

	return poolManager.int64Pool.getPool(index)
}

// ═══════════════════════════════════════════════════════════════
// FLOAT64 POOL LIFECYCLE METHODS
// ═══════════════════════════════════════════════════════════════

// AcquireFloat64Pool acquires a float64 buffer from the pool.
// Returns the buffer index and the slice, or NotAvailable if no buffer is available.
// Thread-safe for global instances.
func (poolManager *PoolManager) AcquireFloat64Pool(size int) (int, []float64) {
	// Apply thread-safety lock for global instances
	if poolManager.global {
		poolManager.lock.Lock()
		defer poolManager.lock.Unlock()
	}

	return poolManager.float64Pool.acquirePool(size)
}

// ReleaseFloat64Pool returns a float64 buffer back to the pool.
// The buffer at the specified index will be marked as available for reuse.
// Thread-safe for global instances.
func (poolManager *PoolManager) ReleaseFloat64Pool(index int) {
	// Apply thread-safety lock for global instances
	if poolManager.global {
		poolManager.lock.Lock()
		defer poolManager.lock.Unlock()
	}

	poolManager.float64Pool.releasePool(index)
}

// GetFloat64Pool retrieves the float64 buffer at the specified index.
// This provides direct access to the underlying slice for reading/writing.
// Thread-safe for global instances (uses read lock).
func (poolManager *PoolManager) GetFloat64Pool(index int) []float64 {
	// Apply read lock for global instances (allows concurrent reads)
	if poolManager.global {
		poolManager.lock.RLock()
		defer poolManager.lock.RUnlock()
	}

	return poolManager.float64Pool.getPool(index)
}

// ═══════════════════════════════════════════════════════════════
// BYTE POOL LIFECYCLE METHODS
// ═══════════════════════════════════════════════════════════════

// AcquireBytePool acquires a byte buffer from the pool.
// Returns the buffer index and the slice, or NotAvailable if no buffer is available.
// Thread-safe for global instances.
func (poolManager *PoolManager) AcquireBytePool(size int) (int, []byte) {
	// Apply thread-safety lock for global instances
	if poolManager.global {
		poolManager.lock.Lock()
		defer poolManager.lock.Unlock()
	}

	return poolManager.bytePool.acquirePool(size)
}

// ReleaseBytePool returns a byte buffer back to the pool.
// The buffer at the specified index will be marked as available for reuse.
// Thread-safe for global instances.
func (poolManager *PoolManager) ReleaseBytePool(index int) {
	// Apply thread-safety lock for global instances
	if poolManager.global {
		poolManager.lock.Lock()
		defer poolManager.lock.Unlock()
	}

	poolManager.bytePool.releasePool(index)
}

// GetBytePool retrieves the byte buffer at the specified index.
// This provides direct access to the underlying slice for reading/writing.
// Thread-safe for global instances (uses read lock).
func (poolManager *PoolManager) GetBytePool(index int) []byte {
	// Apply read lock for global instances (allows concurrent reads)
	if poolManager.global {
		poolManager.lock.RLock()
		defer poolManager.lock.RUnlock()
	}

	return poolManager.bytePool.getPool(index)
}

// ═══════════════════════════════════════════════════════════════
// POOL MAINTENANCE AND OPTIMIZATION
// ═══════════════════════════════════════════════════════════════

// ShrinkAllPools shrinks all oversized buffers back to their default poolLength.
// This helps reclaim memory from buffers that were temporarily expanded.
// Only operates on available (unused) buffers.
// Thread-safe for global instances.
func (poolManager *PoolManager) ShrinkAllPools() {
	// Apply thread-safety lock for global instances
	if poolManager.global {
		poolManager.lock.Lock()
		defer poolManager.lock.Unlock()
	}

	// Shrink all type-specific pools
	poolManager.stringPool.shrinkPool()
	poolManager.int64Pool.shrinkPool()
	poolManager.float64Pool.shrinkPool()
	poolManager.bytePool.shrinkPool()
}

// TestPoolLeak detects memory leaks across all pools by checking for unreleased buffers.
// If any pools are still marked as used, it logs the leaks and forcibly resets them.
// This is typically called during cleanup or testing phases.
// Thread-safe for global instances (uses read lock).
func (poolManager *PoolManager) TestPoolLeak() {
	// Apply read lock for global instances (read-only operation)
	if poolManager.global {
		poolManager.lock.RLock()
		defer poolManager.lock.RUnlock()
	}

	// Test all type-specific pools for leaks
	poolManager.stringPool.testPoolLeak()
	poolManager.int64Pool.testPoolLeak()
	poolManager.float64Pool.testPoolLeak()
	poolManager.bytePool.testPoolLeak()
}

// ═══════════════════════════════════════════════════════════════
// MONITORING AND STATISTICS
// ═══════════════════════════════════════════════════════════════

// GetPoolStats returns current usage statistics for all memory pools.
// This provides a snapshot of which pools are currently in use.
// Thread-safe for global instances (uses read lock).
func (poolManager *PoolManager) GetPoolStats() PoolStats {
	// Apply read lock for global instances (read-only operation)
	if poolManager.global {
		poolManager.lock.RLock()
		defer poolManager.lock.RUnlock()
	}

	// Gather usage statistics from all pool types
	return PoolStats{
		UsedStringPools:  poolManager.stringPool.getUsedPools(),
		UsedInt64Pools:   poolManager.int64Pool.getUsedPools(),
		UsedFloat64Pools: poolManager.float64Pool.getUsedPools(),
		UsedBytePools:    poolManager.bytePool.getUsedPools(),
	}
}
