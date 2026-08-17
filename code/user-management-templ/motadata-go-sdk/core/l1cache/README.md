# Cache Package

## Overview

This package provides high-performance, multi-tenant caching solutions for Go applications. It uses **CacheManager** as the single public interface to manage cache operations with tenant isolation and checksum validation.

**Key Features:**
- **Multi-tenant isolation** with tenant/org key prefixes
- **Zero-allocation operations** with buffer reuse support
- **CRC32 checksum validation** for data integrity
- **Configurable TTL** for cache entries
- **High-performance** FreeCache backend (default)
- **Memory pool integration** for optimal buffer reuse

## Quick Decision Guide

```
What type of cache setup do you need?

MULTI-TENANT ISOLATION → Use CacheManager with tenantID/orgID
ZERO-ALLOCATION READS → Use GetWithBuf() with pooled buffers
CUSTOM CACHE BACKEND → Implement Cache interface and use CacheFactory
HIGH-PERFORMANCE WRITES → Use Set() with pre-allocated buffers
```

## Internal Architecture

### Cache Interface
The underlying `Cache` interface provides core functionality:
- Get/GetWithBuf for reading values
- Set for storing values with TTL
- Del for removing entries
- Clear for removing all entries
- EntryCount and HitRate for statistics

### FreeCache Implementation
The default implementation uses FreeCache, a high-performance cache library that:
- Provides zero GC overhead
- Supports concurrent access
- Manages memory efficiently
- Offers LRU eviction

## CacheManager - The Public Interface

### What It Does
- **Multi-tenant isolation** using tenant/org key prefixes
- **Checksum validation** to ensure data integrity
- **Zero-allocation operations** with buffer reuse
- **Configurable TTL** for cache expiration
- **Unified statistics** for monitoring

### Configuration and Creation

```go
// Default configuration (512MB cache, 5 min TTL)
manager := cache.NewCacheManager(nil)

// Custom configuration
config := &cache.Config{
    TotalSize:  1024 * 1024 * 1024, // 1GB
    DefaultTTL: 10 * time.Minute,   // 10 minutes
}
manager := cache.NewCacheManager(config)

// Custom cache implementation
config := &cache.Config{
    TotalSize:    256 * 1024 * 1024, // 256MB
    DefaultTTL:   time.Hour,
    CacheFactory: myCustomCacheFactory,
}
manager := cache.NewCacheManager(config)
```

### Basic Cache Operations

```go
// Simple set and get operations
tenantID := "tenant123"
orgID := "org456"
key := "user:profile:789"

// Prepare value with checksum space
value := []byte("user profile data")
buffer := make([]byte, cache.ChecksumSize+len(value))
copy(buffer[cache.ChecksumSize:], value)

// Set value
err := manager.Set(tenantID, orgID, key, buffer)

// Get value (allocates memory)
data, err := manager.Get(tenantID, orgID, key)
if err == nil {
    fmt.Printf("Retrieved: %s\n", string(data))
}

// Delete value
deleted := manager.Delete(tenantID, orgID, key)

// Check existence
exists := manager.Exists(tenantID, orgID, key)
```

### Zero-Allocation Operations with Memory Pool

```go
import (
    "motadatagosdk/core/cache"
    "motadatagosdk/core/pool/memorypool"
)

// Create memory pool manager for buffer reuse
poolManager := memorypool.NewPoolManager(&memorypool.PoolConfig{
    PoolSize:   10,    // 10 pools
    PoolLength: 8192,  // 8KB per pool
    Expandable: true,  // Allow growth
    Global:     true,  // Thread-safe
})

// Create cache manager
cacheManager := cache.NewCacheManager(&cache.Config{
    TotalSize:  512 * 1024 * 1024, // 512MB
    DefaultTTL: 5 * time.Minute,
})

// Zero-allocation write using memory pool
func writeToCache(tenantID, orgID, key string, value []byte) error {
    // Acquire buffer from pool
    poolIndex, buffer := poolManager.AcquireBytePool(cache.ChecksumSize + len(value))
    if poolIndex == utils.NotAvailable {
        // Fall back to regular allocation if pool is exhausted
        buffer = make([]byte, cache.ChecksumSize+len(value))
        defer func() {} // No pool to release
    } else {
        defer poolManager.ReleaseBytePool(poolIndex)
    }
    
    // Copy value into buffer (leaving space for checksum)
    copy(buffer[cache.ChecksumSize:], value)
    
    // Set in cache (checksum is written automatically)
    return cacheManager.Set(tenantID, orgID, key, buffer[:cache.ChecksumSize+len(value)])
}

// Zero-allocation read using memory pool
func readFromCache(tenantID, orgID, key string) ([]byte, error) {
    // Acquire buffer from pool for read
    poolIndex, buffer := poolManager.AcquireBytePool(8192) // Reasonable size
    if poolIndex == utils.NotAvailable {
        // Fall back to regular allocation
        buffer = make([]byte, 8192)
        defer func() {} // No pool to release
    } else {
        defer poolManager.ReleaseBytePool(poolIndex)
    }
    
    // Read into pooled buffer
    data, err := cacheManager.GetWithBuf(tenantID, orgID, key, buffer)
    if err != nil {
        return nil, err
    }
    
    // Return copy of data (caller owns this)
    result := make([]byte, len(data))
    copy(result, data)
    return result, nil
}
```

### Cache Monitoring and Maintenance

```go
// Get cache statistics
stats := manager.Stats()
fmt.Printf("Entries: %d, Hit Rate: %.2f%%\n", 
    stats.EntryCount, stats.HitRate*100)

// Clear all cache entries
manager.Clear()

// Clean up resources
manager.Close()
```

## Cache Usage Patterns

### High-Throughput Pattern with Memory Pool

```go
// Initialize shared resources
var (
    cacheManager = cache.NewCacheManager(&cache.Config{
        TotalSize:  1024 * 1024 * 1024, // 1GB
        DefaultTTL: 10 * time.Minute,
    })
    
    poolManager = memorypool.NewPoolManager(&memorypool.PoolConfig{
        PoolSize:   20,     // 20 pools for high concurrency
        PoolLength: 16384,  // 16KB per pool
        Expandable: true,   // Allow growth for large values
        Global:     true,   // Thread-safe shared pools
    })
)

// High-throughput handler using pooled buffers
func handleRequest(tenantID, orgID string, req *Request) (*Response, error) {
    cacheKey := fmt.Sprintf("req:%s", req.ID)
    
    // Try to get from cache first
    poolIndex, readBuffer := poolManager.AcquireBytePool(16384)
    if poolIndex != utils.NotAvailable {
        defer poolManager.ReleaseBytePool(poolIndex)
        
        if cached, err := cacheManager.GetWithBuf(tenantID, orgID, cacheKey, readBuffer); err == nil {
            return parseResponse(cached), nil
        }
    }
    
    // Process request
    response := processRequest(req)
    responseBytes := response.Serialize()
    
    // Cache the response using pooled buffer
    writeIndex, writeBuffer := poolManager.AcquireBytePool(cache.ChecksumSize + len(responseBytes))
    if writeIndex != utils.NotAvailable {
        defer poolManager.ReleaseBytePool(writeIndex)
        
        // Copy response into buffer with checksum space
        copy(writeBuffer[cache.ChecksumSize:], responseBytes)
        
        // Store in cache
        cacheManager.Set(tenantID, orgID, cacheKey, writeBuffer[:cache.ChecksumSize+len(responseBytes)])
    }
    
    return response, nil
}
```

### Batch Processing Pattern

```go
// Batch processor with memory pool for efficiency
type BatchProcessor struct {
    cache       *cache.CacheManager
    poolManager *memorypool.PoolManager
}

func NewBatchProcessor() *BatchProcessor {
    return &BatchProcessor{
        cache: cache.NewCacheManager(&cache.Config{
            TotalSize:  512 * 1024 * 1024,
            DefaultTTL: time.Hour,
        }),
        poolManager: memorypool.NewPoolManager(&memorypool.PoolConfig{
            PoolSize:   5,      // Smaller pool for batch processing
            PoolLength: 65536,  // 64KB buffers for larger batches
            Expandable: false,  // Fixed size for predictability
            Global:     false,  // Local to processor for performance
        }),
    }
}

func (bp *BatchProcessor) ProcessBatch(tenantID, orgID string, items []Item) error {
    // Acquire large buffer for batch operations
    poolIndex, buffer := bp.poolManager.AcquireBytePool(65536)
    if poolIndex == utils.NotAvailable {
        return fmt.Errorf("pool exhausted")
    }
    defer bp.poolManager.ReleaseBytePool(poolIndex)
    
    for _, item := range items {
        // Serialize item
        itemBytes := item.Serialize()
        
        // Ensure buffer is large enough
        requiredSize := cache.ChecksumSize + len(itemBytes)
        if requiredSize > len(buffer) {
            continue // Skip items that are too large
        }
        
        // Reuse buffer for each item
        copy(buffer[cache.ChecksumSize:], itemBytes)
        
        // Cache each item
        key := fmt.Sprintf("item:%s", item.ID)
        bp.cache.Set(tenantID, orgID, key, buffer[:requiredSize])
    }
    
    return nil
}
```

### Session Cache Pattern

```go
// Session manager using cache with memory pools
type SessionManager struct {
    cache       *cache.CacheManager
    poolManager *memorypool.PoolManager
    sessionTTL  time.Duration
}

func NewSessionManager() *SessionManager {
    return &SessionManager{
        cache: cache.NewCacheManager(&cache.Config{
            TotalSize:  256 * 1024 * 1024, // 256MB for sessions
            DefaultTTL: 30 * time.Minute,  // Default 30 min sessions
        }),
        poolManager: memorypool.NewPoolManager(&memorypool.PoolConfig{
            PoolSize:   15,
            PoolLength: 4096, // 4KB typical session size
            Expandable: true,
            Global:     true,  // Thread-safe for concurrent sessions
        }),
        sessionTTL: 30 * time.Minute,
    }
}

func (sm *SessionManager) SaveSession(tenantID, orgID, sessionID string, data *SessionData) error {
    // Serialize session data
    sessionBytes := data.Serialize()
    
    // Get buffer from pool
    poolIndex, buffer := sm.poolManager.AcquireBytePool(cache.ChecksumSize + len(sessionBytes))
    if poolIndex == utils.NotAvailable {
        // Fall back to allocation if pool exhausted
        buffer = make([]byte, cache.ChecksumSize+len(sessionBytes))
    } else {
        defer sm.poolManager.ReleaseBytePool(poolIndex)
    }
    
    // Copy session data
    copy(buffer[cache.ChecksumSize:], sessionBytes)
    
    // Save with custom TTL
    return sm.cache.SetTTL(tenantID, orgID, sessionID, buffer[:cache.ChecksumSize+len(sessionBytes)], sm.sessionTTL)
}

func (sm *SessionManager) GetSession(tenantID, orgID, sessionID string) (*SessionData, error) {
    // Get buffer from pool for reading
    poolIndex, buffer := sm.poolManager.AcquireBytePool(4096)
    if poolIndex == utils.NotAvailable {
        buffer = make([]byte, 4096)
    } else {
        defer sm.poolManager.ReleaseBytePool(poolIndex)
    }
    
    // Read session from cache
    data, err := sm.cache.GetWithBuf(tenantID, orgID, sessionID, buffer)
    if err != nil {
        return nil, err
    }
    
    return DeserializeSessionData(data), nil
}
```

## Simple Usage Examples

### Example 1: Basic Cache Operations

```go
import "motadatagosdk/core/cache"

// Create cache manager
manager := cache.NewCacheManager(nil) // Uses default config

// Store a value
tenantID := "tenant1"
orgID := "org1"
key := "config:app"
value := []byte(`{"theme":"dark","lang":"en"}`)

// Prepare buffer with checksum space
buffer := make([]byte, cache.ChecksumSize+len(value))
copy(buffer[cache.ChecksumSize:], value)

// Set in cache
err := manager.Set(tenantID, orgID, key, buffer)
if err != nil {
    log.Printf("Cache set error: %v", err)
}

// Retrieve value
data, err := manager.Get(tenantID, orgID, key)
if err == nil {
    fmt.Printf("Config: %s\n", string(data))
}
```

### Example 2: Zero-Allocation with Memory Pool

```go
import (
    "motadatagosdk/core/cache"
    "motadatagosdk/core/pool/memorypool"
)

// Setup
cacheManager := cache.NewCacheManager(nil)
poolManager := memorypool.NewPoolManager(&memorypool.PoolConfig{
    PoolSize:   5,
    PoolLength: 8192,
    Expandable: true,
    Global:     true,
})

// Write with pooled buffer
poolIndex, buffer := poolManager.AcquireBytePool(1024)
if poolIndex != utils.NotAvailable {
    defer poolManager.ReleaseBytePool(poolIndex)
    
    value := []byte("cached data")
    copy(buffer[cache.ChecksumSize:], value)
    
    cacheManager.Set("tenant1", "org1", "key1", buffer[:cache.ChecksumSize+len(value)])
}

// Read with pooled buffer
readIndex, readBuf := poolManager.AcquireBytePool(1024)
if readIndex != utils.NotAvailable {
    defer poolManager.ReleaseBytePool(readIndex)
    
    data, err := cacheManager.GetWithBuf("tenant1", "org1", "key1", readBuf)
    if err == nil {
        fmt.Printf("Data: %s\n", string(data))
    }
}
```

### Example 3: Custom TTL and Statistics

```go
import "motadatagosdk/core/cache"

manager := cache.NewCacheManager(&cache.Config{
    TotalSize:  256 * 1024 * 1024, // 256MB
    DefaultTTL: time.Hour,         // 1 hour default
})

// Set with custom TTL
value := []byte("temporary data")
buffer := make([]byte, cache.ChecksumSize+len(value))
copy(buffer[cache.ChecksumSize:], value)

// Store for only 1 minute
manager.SetTTL("tenant1", "org1", "temp-key", buffer, 1*time.Minute)

// Check statistics
stats := manager.Stats()
fmt.Printf("Cache entries: %d, Hit rate: %.2f%%\n", 
    stats.EntryCount, stats.HitRate*100)

// Check if key exists
if manager.Exists("tenant1", "org1", "temp-key") {
    fmt.Println("Key exists")
}

// Clean up
manager.Clear()
manager.Close()
```

## Key Features Comparison

| Feature | Standard Operations | Memory Pool Operations |
|---------|-------------------|------------------------|
| **Memory Allocation** | New allocation per operation | Reuses pooled buffers |
| **GC Pressure** | Higher | Minimal |
| **Performance** | Good | Excellent |
| **Complexity** | Simple | Moderate |
| **Use Case** | Low-frequency access | High-throughput scenarios |
| **Buffer Management** | Automatic (GC) | Manual (acquire/release) |
| **Best For** | Regular caching | Performance-critical paths |

## Best Practices

### CacheManager Best Practices
1. **Always include checksum space**: Reserve ChecksumSize bytes at buffer start
2. **Use GetWithBuf for reads**: Avoid allocations in hot paths
3. **Pool buffers for performance**: Use memory pools for high-throughput
4. **Handle pool exhaustion**: Fall back to allocation when pools are full
5. **Size buffers appropriately**: Match buffer sizes to typical value sizes
6. **Monitor statistics**: Track hit rates and entry counts
7. **Set appropriate TTL**: Balance freshness with cache efficiency
8. **Clean up resources**: Call Close() when done

### Memory Pool Integration
1. **Choose pool strategy**: Global for shared access, Local for isolated
2. **Size pools correctly**: Based on concurrent operations and value sizes
3. **Always release pools**: Use defer to ensure pools are returned
4. **Handle NotAvailable**: Fall back to regular allocation gracefully
5. **Monitor pool usage**: Use GetPoolStats() to track efficiency
6. **Shrink oversized pools**: Call ShrinkAllPools() periodically
7. **Test for leaks**: Use TestPoolLeak() in development

## Performance Tips

1. **Buffer reuse**: Use memory pools to eliminate allocation overhead
2. **Batch operations**: Process multiple items with single pool acquisition
3. **Right-size buffers**: Allocate based on actual data sizes
4. **Local pools**: Use non-global pools for single-threaded scenarios
5. **Pre-allocate buffers**: Size pools based on expected load
6. **Avoid string conversions**: Work with []byte throughout
7. **Monitor and tune**: Use statistics to optimize cache size and TTL

## Testing

Run tests for the cache implementation:

```bash
# Run all tests
go test ./...

# Run with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...

# Run specific cache tests
go test -run TestCacheManager ./l1cache/
go test -run TestFreeCache ./l1cache/
```

## Migration Guide

### From sync.Map to CacheManager

Before:
```go
var cache sync.Map

// Write
cache.Store("key", value)

// Read
if val, ok := cache.Load("key"); ok {
    // Use val
}
```

After:
```go
manager := cache.NewCacheManager(nil)

// Write with tenant isolation
buffer := make([]byte, cache.ChecksumSize+len(value))
copy(buffer[cache.ChecksumSize:], value)
manager.Set(tenantID, orgID, "key", buffer)

// Read with checksum validation
if val, err := manager.Get(tenantID, orgID, "key"); err == nil {
    // Use val
}
```

### From Simple Cache to Memory Pool Integration

Before:
```go
// Simple cache without pooling
func cacheData(key string, data []byte) {
    buffer := make([]byte, cache.ChecksumSize+len(data))
    copy(buffer[cache.ChecksumSize:], data)
    manager.Set("tenant", "org", key, buffer)
}
```

After:
```go
// Cache with memory pool for zero allocations
func cacheData(key string, data []byte) {
    poolIndex, buffer := poolManager.AcquireBytePool(cache.ChecksumSize + len(data))
    if poolIndex != utils.NotAvailable {
        defer poolManager.ReleaseBytePool(poolIndex)
        copy(buffer[cache.ChecksumSize:], data)
        manager.Set("tenant", "org", key, buffer[:cache.ChecksumSize+len(data)])
    } else {
        // Fall back to allocation
        buffer := make([]byte, cache.ChecksumSize+len(data))
        copy(buffer[cache.ChecksumSize:], data)
        manager.Set("tenant", "org", key, buffer)
    }
}
```

## Summary

The **CacheManager** provides high-performance caching with these key benefits:

**Core Features:**
- **Multi-tenant isolation** with automatic key prefixing
- **CRC32 checksums** for data integrity validation
- **Zero-allocation operations** with buffer reuse support
- **Memory pool integration** for optimal performance
- **Configurable TTL** for flexible expiration
- **FreeCache backend** for zero-GC overhead

**Performance Optimization:**
- Use **memory pools** to eliminate allocation overhead
- Leverage **GetWithBuf()** for zero-allocation reads
- Pre-allocate buffers with **ChecksumSize** space
- Choose **Global vs Local** pools based on access patterns
- Monitor with **Stats()** to optimize configuration

**Best Use Cases:**
- High-throughput web services
- Multi-tenant applications
- Session management
- Configuration caching
- Batch data processing
- Real-time data pipelines

The cache package integrates seamlessly with the memory pool package to provide maximum performance with minimal GC pressure, making it ideal for performance-critical applications.