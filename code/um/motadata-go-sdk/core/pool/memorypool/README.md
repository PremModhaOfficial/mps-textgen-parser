# Memory Pool Package

## Overview

This package provides efficient memory pooling solutions for high-performance Go applications. It uses **PoolManager** as the single public interface to manage type-safe memory pools with leak detection.

**Key Features:**
- **Type-safe pools** for string, int64, float64, and byte arrays
- **Global vs Local pools** - Choose centralized thread-safe or per-goroutine pools
- **Automatic leak detection** and debugging capabilities
- **Expandable pools** that can grow beyond initial size
- **Zero-value initialization** for security and predictability

## Quick Decision Guide

```
What type of pool setup do you need?

SHARED ACROSS GOROUTINES → Use PoolManager with Global: true (thread-safe)
ISOLATED PER GOROUTINE → Use PoolManager with Global: false (no locking overhead)
MULTIPLE POOL TYPES → Use PoolManager (manages all 4 types: string, int64, float64, byte)
```

## Internal Architecture

### MemoryPool (Private Implementation)
The underlying `MemoryPool[T]` struct provides the core functionality but is **not directly accessible**. It:
- Manages **typed array slices** ([]string, []int64, []float64, []byte)
- Provides **leak detection** and debugging capabilities  
- Supports **expandable pools** that can grow beyond initial size
- Uses **bitmask tracking** for efficient pool state management
- **Zero-value initialization** for security and predictability

**Note:** All pool operations must go through `PoolManager` - the `MemoryPool` methods are private and not exported.

## PoolManager - The Public Interface

### What It Does
- **Single public interface** to manage all four typed pools (string, int64, float64, byte)
- Supports **Global vs Local** pool strategies
- **Global pools**: Thread-safe with RWMutex, shared across goroutines
- **Local pools**: No locking overhead, isolated per goroutine/component
- **Unified configuration** for all pool types
- **Centralized statistics** and leak detection

### Global vs Local Pool Strategy

**Global Pools (`Global: true`):**
- Thread-safe with automatic RWMutex locking
- Shared across all goroutines
- Best for: Web servers, shared data processing, concurrent access patterns
- Slightly higher overhead due to locking

**Local Pools (`Global: false`):**
- No locking overhead
- Isolated per PoolManager instance (per goroutine/component)
- Best for: Worker goroutines, isolated processing, performance-critical paths
- Lower overhead, but requires separate manager per goroutine

### Configuration and Creation

```go
// For global shared pools (thread-safe)
globalConfig := &memorypool.PoolConfig{
    PoolSize:   10,    // Number of pools per type
    PoolLength: 1000,  // Elements per pool
    Expandable: true,  // Allow pools to grow beyond PoolLength
    Global:     true,  // Global singleton with thread-safety
}
globalManager := memorypool.NewPoolManager(globalConfig)

// For local per-goroutine pools (no locking)
localConfig := &memorypool.PoolConfig{
    PoolSize:   5,     // Smaller pools for local use
    PoolLength: 500,   // Elements per pool
    Expandable: false, // Fixed size for predictability
    Global:     false, // Local instance, no locking
}
localManager := memorypool.NewPoolManager(localConfig)
```

### Basic Pool Operations

```go
// String pool operations
index, strings := manager.AcquireStringPool(500)
if index != utils.NotAvailable {
    defer manager.ReleaseStringPool(index)
    strings[0] = "Processing data"
    strings[1] = "in memorypool pool"
}

// Int64 pool operations  
index, numbers := manager.AcquireInt64Pool(1000)
if index != utils.NotAvailable {
    defer manager.ReleaseInt64Pool(index)
    for i := 0; i < 1000; i++ {
        numbers[i] = int64(i * 2)
    }
}

// Float64 pool operations
index, values := manager.AcquireFloat64Pool(200)
if index != utils.NotAvailable {
    defer manager.ReleaseFloat64Pool(index)
    for i := 0; i < 200; i++ {
        values[i] = float64(i) * 3.14
    }
}

// Byte pool operations
index, buffer := manager.AcquireBytePool(4096)
if index != utils.NotAvailable {
    defer manager.ReleaseBytePool(index)
    buffer[0] = 255
    buffer[1] = 128
}
```

### Pool Monitoring and Maintenance

```go
// Get current pool statistics
stats := manager.GetPoolStats()
fmt.Printf("Used pools - Strings: %d, Int64: %d, Float64: %d, Bytes: %d\n", 
    stats.UsedStringPools, stats.UsedInt64Pools, 
    stats.UsedFloat64Pools, stats.UsedBytePools)

// Shrink all oversized pools back to original size
manager.ShrinkAllPools()

// Check for memorypool leaks across all pools (prints stack trace if leaked)
manager.TestPoolLeak()

// Direct pool access for inspection (read-only)
if index != utils.NotAvailable {
    currentStrings := manager.GetStringPool(index)
    fmt.Printf("Current string pool size: %d\n", len(currentStrings))
}
```

## Pool Usage Patterns

### Global Pool Pattern (Shared Across Goroutines)

```go
// Create global pool manager once (typically at app startup)
var globalPoolManager = memorypool.NewPoolManager(&memorypool.PoolConfig{
    PoolSize:   10,
    PoolLength: 1000,
    Expandable: true,
    Global:     true, // Thread-safe, shared across goroutines
})

// Use in multiple goroutines safely
func processInGoroutine(data []string) {
    // Acquire string pool (thread-safe)
    index, strings := globalPoolManager.AcquireStringPool(len(data))
    if index != utils.NotAvailable {
        defer globalPoolManager.ReleaseStringPool(index)
        
        // Process data using pooled strings
        for i, d := range data {
            if i < len(strings) {
                strings[i] = d + " processed"
            }
        }
    }
}
```

### Local Pool Pattern (Per-Goroutine)

```go
// Each goroutine creates its own pool manager
func workerGoroutine(workerID int, dataChannel <-chan []byte) {
    // Local pool manager for this workerpool (no locking overhead)
    localManager := memorypool.NewPoolManager(&memorypool.PoolConfig{
        PoolSize:   5,
        PoolLength: 4096,
        Expandable: false,
        Global:     false, // Local to this goroutine, no locking
    })
    
    for data := range dataChannel {
        // Acquire byte pool (no locking, very fast)
        index, buffer := localManager.AcquireBytePool(len(data))
        if index != utils.NotAvailable {
            defer localManager.ReleaseBytePool(index)
            
            // Process data using pooled buffer
            copy(buffer, data)
            processBuffer(workerID, buffer[:len(data)])
        }
    }
}
```

### Pool Monitoring and Metrics

```go
// Monitor pool usage and performance
func monitorPools(manager *memorypool.PoolManager) {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        // Get current statistics
        stats := manager.GetPoolStats()
        
        // Log usage metrics
        fmt.Printf("Pool Usage - Strings: %d, Int64: %d, Float64: %d, Bytes: %d\n", 
            stats.UsedStringPools, stats.UsedInt64Pools, 
            stats.UsedFloat64Pools, stats.UsedBytePools)
        
        // Perform maintenance
        manager.ShrinkAllPools()  // Shrink oversized pools
        manager.TestPoolLeak()    // Check for leaks (debug builds)
    }
}
```

## Simple Usage Examples

### Example 1: Basic PoolManager Usage

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/core/memorypool"

// Create pool manager
manager := memorypool.NewPoolManager(&memorypool.PoolConfig{
    PoolSize:   5,     // 5 pools per type
    PoolLength: 1000,  // 1000 elements per pool  
    Expandable: true,  // Allow pools to grow
    Global:     true,  // Thread-safe global pools
})

// Use string pool
strIndex, strings := manager.AcquireStringPool(100)
if strIndex != utils.NotAvailable {
    defer manager.ReleaseStringPool(strIndex)
    strings[0] = "Hello"
    strings[1] = "World"
    fmt.Printf("First string: %s\n", strings[0])
}
```

### Example 2: Multiple Pool Types

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/core/memorypool"

manager := memorypool.NewPoolManager(&memorypool.PoolConfig{
    PoolSize:   8,
    PoolLength: 1000,
    Expandable: true,
    Global:     false, // Local pools, no locking
})

// Use multiple pool types simultaneously
strIndex, strings := manager.AcquireStringPool(100)
intIndex, numbers := manager.AcquireInt64Pool(100)
byteIndex, buffer := manager.AcquireBytePool(4096)

if strIndex != utils.NotAvailable {
    defer manager.ReleaseStringPool(strIndex)
    strings[0] = "processed data"
}

if intIndex != utils.NotAvailable {
    defer manager.ReleaseInt64Pool(intIndex)
    numbers[0] = 12345
}

if byteIndex != utils.NotAvailable {
    defer manager.ReleaseBytePool(byteIndex)
    buffer[0] = 255
}

// Check statistics
stats := manager.GetPoolStats()
fmt.Printf("Used pools - Strings: %d, Int64: %d, Bytes: %d\n", 
    stats.UsedStringPools, stats.UsedInt64Pools, stats.UsedBytePools)
```

### Example 3: Pool Maintenance

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/core/memorypool"

manager := memorypool.NewPoolManager(&memorypool.PoolConfig{
    PoolSize:   8,
    PoolLength: 1024,
    Expandable: true,
    Global:     true,
})

// Request larger than poolLength to test expansion
index, buffer := manager.AcquireBytePool(2048)
if index != utils.NotAvailable {
    defer manager.ReleaseBytePool(index)
    
    buffer[0] = 255
    buffer[1] = 128
    
    fmt.Printf("Acquired expanded buffer with %d bytes\n", len(buffer))
}

// Periodic maintenance operations
manager.ShrinkAllPools()    // Shrink oversized pools back to original size
manager.TestPoolLeak()      // Check for memorypool leaks (debug mode)

// Monitor current usage
stats := manager.GetPoolStats()
fmt.Printf("Current usage - Strings: %d, Bytes: %d\n", 
    stats.UsedStringPools, stats.UsedBytePools)
```

## Key Features Comparison

| Feature | Global Pools | Local Pools |
|---------|-------------|-------------|
| **Thread Safety** | Yes (RWMutex locking) | No (isolated per goroutine) |
| **Performance** | Medium (locking overhead) | High (no locking) |
| **Sharing** | Shared across goroutines | Isolated per PoolManager instance |
| **Memory Usage** | Shared pool memory | Per-goroutine pool memory |
| **Use Case** | Web servers, shared processing | Worker goroutines, isolated processing |
| **Configuration** | `Global: true` | `Global: false` |
| **Best For** | Concurrent access patterns | Performance-critical single-threaded paths |

## Best Practices

### PoolManager Best Practices
1. **Always release pools**: Use defer to ensure pools are returned
2. **Handle pool exhaustion**: Always check for utils.NotAvailable return value
3. **Choose Global vs Local**: Use Global=true for shared access, Global=false for isolated components
4. **Configure appropriately**: Balance PoolSize and PoolLength based on usage patterns
5. **Use expandable wisely**: Enable for unpredictable workloads, disable for consistent sizes
6. **Monitor statistics**: Use GetPoolStats() to track usage patterns
7. **Periodic maintenance**: Call ShrinkAllPools() and TestPoolLeak() for maintenance
8. **Zero sensitive data**: Pools automatically zero values on allocation for security

### General Best Practices
1. **Profile your application**: Measure impact of pooling on performance
2. **Choose Global vs Local**: Global for shared access, Local for isolated high-performance paths
3. **Monitor memory usage**: Track pool effectiveness with built-in GetPoolStats()
4. **Avoid premature optimization**: Pool when you have proven allocation pressure
5. **Test thoroughly**: Especially with concurrent access patterns using Global pools
6. **Proper error handling**: Always check for utils.NotAvailable return values

## Performance Tips

1. **Global vs Local**: Choose Local pools (Global=false) for maximum performance in single-threaded scenarios
2. **Pool sizing**: Pre-calculate optimal PoolSize and PoolLength based on usage patterns
3. **Size estimation**: Start with realistic estimates and monitor with GetPoolStats()
4. **Expandable pools**: Use for variable workloads, disable for predictable patterns
5. **Pool acquisition**: Minimize time between Acquire and Release calls
6. **Concurrent access**: Global pools handle thread-safety automatically with RWMutex
7. **Memory monitoring**: Use TestPoolLeak() and ShrinkAllPools() for maintenance

## Testing

Run tests for the memory pool implementations:

```bash
# Run all tests
go test ./...

# Run with race detection
go test -race ./...

# Run benchmarks  
go test -bench=. ./...

# Run specific memorypool pool tests
go test -run TestMemoryPool ./memorypool/
go test -run TestPoolManager ./memorypool/
```

## Migration Guide

### From Manual Allocation to MemoryPool

Before:
```go
// Manual allocation - creates garbage
func processData(size int) {
    data := make([]string, size) // New allocation each time
    defer data = nil // Doesn't help with GC pressure
    // Process data...
}
```

After:
```go
// Memory pool - reuses allocations
func processData(size int) {
    index, data := pool.acquirePool(size)
    if index != utils.NotAvailable {
        defer pool.releasePool(index)
        // Process data...
    }
}
```

### From sync.Pool to PoolManager

Before:
```go
var bytePool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 0, 1024)
    },
}

func process() {
    buf := bytePool.Get().([]byte)
    defer bytePool.Put(buf[:0])
    // Use buf...
}
```

After:
```go
// Create once during initialization
var poolManager = memorypool.NewPoolManager(&memorypool.PoolConfig{
    PoolSize:   5,     // 5 pools per type
    PoolLength: 1024,  // 1024 elements per pool
    Expandable: true,  // Allow growth
    Global:     true,  // Thread-safe shared pools
})

func process() {
    index, buf := poolManager.AcquireBytePool(1024)
    if index != utils.NotAvailable {
        defer poolManager.ReleaseBytePool(index)
        // Use buf... (automatically zero-initialized)
    }
}
```

## Summary

The **PoolManager** is the single public interface for type-safe memory pooling with these key benefits:

**Core Features:**
- **Type-safe pools** for string, int64, float64, and byte arrays
- **Global vs Local** pool strategies for different use cases
- **Automatic leak detection** and debugging capabilities
- **Expandable pools** that can grow beyond initial configuration
- **Zero-value initialization** for security

**Choose Global Pools (`Global: true`) when you need:**
- Shared access across multiple goroutines
- Thread-safe operations with automatic locking
- Centralized pool management for web servers or concurrent processing

**Choose Local Pools (`Global: false`) when you need:**
- Maximum performance in single-threaded scenarios
- Isolated pools per goroutine or component
- No locking overhead for performance-critical paths

**All pool operations go through PoolManager** - the underlying MemoryPool implementation is private and not directly accessible. This ensures consistent API usage and proper resource management across all pool types.