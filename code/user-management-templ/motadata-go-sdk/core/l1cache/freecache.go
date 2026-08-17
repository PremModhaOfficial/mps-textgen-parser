// Package l1cache - FreeCache Adapter Implementation
//
// ARCHITECTURE OVERVIEW:
//
// This module provides the FreeCache adapter that implements the Cache interface,
// wrapping the high-performance github.com/coocood/freecache library. FreeCache
// is chosen as the default backend due to its zero-allocation design, predictable
// memory usage, and excellent concurrent performance characteristics.
//
// FREECACHE DESIGN PRINCIPLES:
//
// FreeCache is a cache library that:
// 1. Uses fixed-size memory allocation (no GC pressure)
// 2. Implements segmented locking for concurrency
// 3. Provides LRU eviction with approximation
// 4. Achieves zero-allocation on Gets
// 5. Maintains consistent memory bounds
//
// INTERNAL ARCHITECTURE:
//
// FreeCache Structure:
//
//	┌────────────────────────────────────┐
//	│         FreeCache Instance         │
//	├────────────────────────────────────┤
//	│  Segment 0  │ RingBuffer + Index   │
//	│  Segment 1  │ RingBuffer + Index   │
//	│  ...        │ ...                  │
//	│  Segment 255│ RingBuffer + Index   │
//	└────────────────────────────────────┘
//
// Each segment contains:
// - Ring buffer for data storage
// - Hash index for fast lookups
// - Independent lock for concurrency
//
// MEMORY MANAGEMENT:
//
// Fixed Memory Allocation:
// FreeCache pre-allocates all memory upfront:
//
//	cache := freecache.NewCache(512 * 1024 * 1024) // 512MB
//
// Memory layout per segment:
//
//	Total Size / 256 segments = Per-segment size
//	512MB / 256 = 2MB per segment
//
// This ensures:
// - No runtime memory allocations
// - Predictable memory footprint
// - No GC pressure from cache operations
// - Protection against memory leaks
//
// SEGMENTATION STRATEGY:
//
// Key Distribution:
// Keys are distributed across 256 segments using:
//
//	segment_id = hash(key) % 256
//
// Benefits:
// - Reduced lock contention (1/256th)
// - Better CPU cache utilization
// - Parallel operations on different segments
// - Scalable to multi-core systems
//
// EVICTION ALGORITHM:
//
// Approximated LRU (Least Recently Used):
// FreeCache uses an approximated LRU that:
// 1. Tracks access time with 1-second granularity
// 2. Evicts oldest entries when space needed
// 3. Scans limited entries for efficiency
// 4. Balances accuracy vs performance
//
// Eviction triggers:
// - Segment full → evict oldest in segment
// - TTL expired → lazy deletion on access
// - Explicit deletion → immediate removal
//
// CONCURRENCY MODEL:
//
// Locking Strategy:
// - 256 independent locks (one per segment)
// - Read operations: Shared lock per segment
// - Write operations: Exclusive lock per segment
// - No global locks
//
// Concurrency characteristics:
// - Theoretical max parallelism: 256 operations
// - Practical parallelism: ~8-16 (CPU cores)
// - Lock hold time: ~100ns typical
// - Contention probability: 1/256 per operation
//
// ZERO-ALLOCATION OPERATIONS:
//
// GetWithBuf Implementation:
// The zero-allocation get is achieved through:
// 1. Caller provides buffer
// 2. Direct copy into buffer
// 3. No intermediate allocations
// 4. Return slice of buffer
//
// Usage pattern:
//
//	buf := make([]byte, 1024) // Reusable buffer
//	value, err := cache.GetWithBuf(key, buf)
//	// value points into buf, no allocation
//
// PERFORMANCE CHARACTERISTICS:
//
// Operation Benchmarks (typical):
// - Get: 50-100ns
// - GetWithBuf: 40-80ns (zero-alloc)
// - Set: 100-200ns
// - Delete: 50-100ns
//
// Scalability:
// - Single-core: 10M+ ops/sec
// - Multi-core: Near-linear scaling to 8 cores
// - 16+ cores: Diminishing returns
//
// Memory efficiency:
// - Overhead: ~10% of configured size
// - No heap allocations in steady state
// - Bounded memory usage
// - No memory fragmentation
//
// TTL HANDLING:
//
// TTL Implementation:
// - Second-precision timestamps
// - Lazy expiration on access
// - No background expiration threads
// - Space reclaimed on eviction
//
// TTL storage:
//
//	Entry: [TTL:4bytes][Key:variable][Value:variable]
//
// ADAPTER PATTERN:
//
// The FreeCache struct adapts the FreeCache API to our Cache interface:
//
//	Cache Interface          FreeCache Implementation
//	├─ Get()          →     cache.Get()
//	├─ GetWithBuf()   →     cache.GetWithBuf()
//	├─ Set()          →     cache.Set()
//	├─ Del()          →     cache.Del()
//	├─ Clear()        →     cache.Clear()
//	├─ EntryCount()   →     cache.EntryCount()
//	└─ HitRate()      →     cache.HitRate()
//
// This abstraction enables:
// - Swappable cache backends
// - Consistent API across implementations
// - Testability with mock caches
// - Future migration paths
//
// ERROR HANDLING:
//
// FreeCache error semantics:
// - Key not found: Returns error
// - Set failure: Returns error (memory full)
// - Delete miss: Returns false
// - Corrupted data: Panic (should never happen)
//
// The adapter passes errors through transparently.
//
// MONITORING CAPABILITIES:
//
// Built-in metrics:
// - EntryCount(): Current number of entries
// - HitRate(): Cumulative hit rate (0.0-1.0)
// - AverageAccessTime(): Not exposed, available internally
// - EvacuateCount(): Not exposed, eviction counter
//
// These metrics enable:
// - Cache efficiency monitoring
// - Capacity planning
// - Performance debugging
// - Eviction rate tracking
//
// LIMITATIONS AND TRADE-OFFS:
//
// FreeCache limitations:
// 1. Fixed memory size (no dynamic growth)
// 2. Approximate LRU (not exact)
// 3. Second-precision TTL
// 4. No persistence
// 5. Local-only (not distributed)
//
// Trade-offs made for performance:
// - Accuracy vs speed (approximate LRU)
// - Flexibility vs predictability (fixed size)
// - Features vs simplicity (no persistence)
//
// BEST PRACTICES:
//
//  1. Size configuration:
//     // Set to 70% of available memory
//     size := int(0.7 * availableMemory)
//
//  2. Buffer pooling:
//     var bufferPool = sync.Pool{
//     New: func() any { return make([]byte, 4096) },
//     }
//
//  3. Error handling:
//     if err != nil && err != freecache.ErrNotFound {
//     // Handle unexpected error
//     }
//
// FUTURE OPTIMIZATIONS:
//
// - SIMD-accelerated hash computation
// - NUMA-aware segment placement
// - Adaptive eviction policies
// - Compression for large values
// - Tiered storage (memory + SSD)
package l1cache

import (
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"

	"github.com/coocood/freecache"
)

// FreeCache implements Cache interface using FreeCache.
// It wraps the FreeCache library to provide a consistent caching interface
// with zero-allocation operations and predictable memory usage.
type FreeCache struct {
	cache *freecache.Cache
}

// initFreeCache creates a FreeCache instance as the default CacheFactory.
// The sizeBytes parameter determines the fixed memory allocation for the cache.
func initFreeCache(sizeBytes int) Cache {
	return &FreeCache{
		cache: freecache.NewCache(sizeBytes),
	}
}

// Get retrieves a value by key from the cache
func (freeCache *FreeCache) Get(key []byte) ([]byte, error) {
	return freeCache.cache.Get(key)
}

// GetWithBuf retrieves a value by key using provided buffer to avoid allocation
func (freeCache *FreeCache) GetWithBuf(key, buf []byte) ([]byte, error) {
	return freeCache.cache.GetWithBuf(key, buf)
}

// Set stores a key-value pair with TTL in seconds
func (freeCache *FreeCache) Set(key, value []byte, ttlSeconds int) error {
	return freeCache.cache.Set(key, value, ttlSeconds)
}

// Del deletes a key from the cache
func (freeCache *FreeCache) Del(key []byte) bool {
	return freeCache.cache.Del(key)
}

// Clear removes all entries from the cache
func (freeCache *FreeCache) Clear() {
	if !utils.IsZeroValue(freeCache.cache) {
		freeCache.cache.Clear()
	}
}

// EntryCount returns the number of entries in the cache
func (freeCache *FreeCache) EntryCount() int64 {
	return freeCache.cache.EntryCount()
}

// HitRate returns the cache hit rate
func (freeCache *FreeCache) HitRate() float64 {
	return freeCache.cache.HitRate()
}
