# Pool Package - Complete Pooling Solutions

The Pool package provides three specialized pooling implementations designed for different use cases in high-performance Go applications. Each pool type is optimized for specific scenarios and resource management patterns.

## 📦 Pool Types Overview

| Pool Type | Purpose | Best For | Key Features |
|-----------|---------|----------|--------------|
| **WorkerPool** | Concurrent task execution | Background jobs, API processing | Auto-scaling, context support, panic recovery |
| **MemoryPool** | Memory allocation reuse | Data processing, buffer management | Type-safe arrays, leak detection, zero-init |
| **ResourcePool** | Expensive resource management | DB connections, HTTP clients | Lifecycle callbacks, graceful shutdown |

---

## 🚀 WorkerPool - Concurrent Task Execution

**Purpose**: Execute tasks concurrently with automatic worker management and resource control.

### Quick Start
```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/core/pool/workerpool"

// Initialize global pool
workerpool.MustInit(workerpool.PoolConfig{
    MinWorkers: 5,
    MaxWorkers: 50,
    Timeout:    30 * time.Second,
})
defer workerpool.Shutdown()

// Fire-and-forget task
workerpool.Async(func() {
    fmt.Println("Background task")
})

// Synchronous execution with context
ctx := context.Background()
err := workerpool.SyncWithCtx(ctx, func() error {
    return processData()
})

// Batch processing
tasks := []func() error{task1, task2, task3}
errors := workerpool.Batch(ctx, tasks...)
```

### Key Features
- ✅ **Auto-scaling**: 2-50 workers based on load
- ✅ **Context support**: Cancellation and timeouts
- ✅ **Panic recovery**: Failed tasks don't crash pool
- ✅ **Batch operations**: Process multiple tasks efficiently
- ✅ **Scheduling**: Delayed and recurring tasks
- ✅ **Metrics**: Real-time pool statistics

### Use Cases
- **Web API backends**: Handle HTTP requests concurrently
- **Data processing**: Process files, messages, or streams
- **Background jobs**: Execute async tasks without blocking
- **Microservices**: Scale worker count based on demand

### Performance
- **Throughput**: 1M+ operations/sec for simple tasks
- **Latency**: Sub-microsecond task submission
- **Memory**: 24B/operation base allocation
- **Scaling**: Linear up to CPU core count

---

## 💾 MemoryPool - Efficient Memory Reuse

**Purpose**: Reuse pre-allocated memory arrays to reduce garbage collection pressure.

### Quick Start
```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/core/pool/memorypool"

// Global shared pools (thread-safe)
manager := memorypool.NewPoolManager(&memorypool.PoolConfig{
    PoolSize:   10,    // 10 pools per type
    PoolLength: 1000,  // 1000 elements per pool
    Expandable: true,  // Allow growth beyond initial size
    Global:     true,  // Thread-safe shared access
})

// Acquire string pool
index, strings := manager.AcquireStringPool(500)
if index != utils.NotAvailable {
    defer manager.ReleaseStringPool(index)
    
    strings[0] = "Processing data"
    strings[1] = "in memory pool"
    // Use strings array...
}

// Multiple pool types
intIndex, numbers := manager.AcquireInt64Pool(1000)
byteIndex, buffer := manager.AcquireBytePool(4096)
floatIndex, values := manager.AcquireFloat64Pool(200)
```

### Key Features
- ✅ **Type-safe pools**: string, int64, float64, byte arrays
- ✅ **Global vs Local**: Choose thread-safe shared or isolated pools
- ✅ **Leak detection**: Automatic memory leak tracking
- ✅ **Zero initialization**: Security through clean slate
- ✅ **Expandable pools**: Grow beyond initial configuration
- ✅ **Usage statistics**: Monitor pool efficiency

### Pool Strategies

**Global Pools (`Global: true`)**:
```go
// Shared across goroutines, thread-safe with mutex locking
globalManager := memorypool.NewPoolManager(&memorypool.PoolConfig{
    PoolSize:   10,
    PoolLength: 1000,
    Global:     true,  // Thread-safe, shared access
})
```

**Local Pools (`Global: false`)**:
```go
// Per-goroutine pools, no locking overhead
localManager := memorypool.NewPoolManager(&memorypool.PoolConfig{
    PoolSize:   5,
    PoolLength: 500,
    Global:     false, // High-performance, isolated access
})
```

### Use Cases
- **Data processing**: Reuse arrays for parsing, calculations
- **Web servers**: Buffer pools for request/response handling
- **Stream processing**: Reduce allocations in hot paths
- **Batch operations**: Process large datasets efficiently

### Performance
- **Global pools**: Medium overhead due to locking
- **Local pools**: High performance, zero locking
- **Memory savings**: 60-90% reduction in allocations
- **GC pressure**: Significant reduction in garbage collection

---

## 🔗 ResourcePool - Expensive Resource Management

**Purpose**: Manage costly, stateful resources like database connections, HTTP clients, and file handles.

### Quick Start
```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/core/pool/resourcepool"

// Database connection pool
config := resourcepool.ResourcePoolConfig[*sql.DB]{
    MaxSize: 10,
    OnCreate: func() (*sql.DB, error) {
        return sql.Open("postgres", dsn)
    },
    OnReset: func(db *sql.DB) error {
        // Reset connection state
        return nil
    },
    OnDestroy: func(db *sql.DB) {
        db.Close()
    },
}

pool, err := resourcepool.NewResourcePool(config)
if err != nil {
    log.Fatal(err)
}
defer pool.Shutdown()

// Use resource with context
ctx := context.Background()
db, err := pool.Get(ctx)
if err != nil {
    return err
}
defer pool.Put(db) // Always return to pool

// Use the database connection
rows, err := db.Query("SELECT * FROM users")
// Process results...
```

### Key Features
- ✅ **Type-safe generics**: Any resource type supported
- ✅ **Lifecycle management**: onCreate, onReset, onDestroy callbacks
- ✅ **Context-aware**: Timeout and cancellation support
- ✅ **Thread-safe**: Concurrent resource acquisition/release
- ✅ **Graceful shutdown**: Wait for resources or force cleanup
- ✅ **Resource monitoring**: Usage statistics and leak detection

### Advanced Examples

**HTTP Client Pool**:
```go
type HTTPClient struct {
    Client    *http.Client
    BaseURL   string
}

config := resourcepool.ResourcePoolConfig[*HTTPClient]{
    MaxSize: 5,
    OnCreate: func() (*HTTPClient, error) {
        return &HTTPClient{
            Client: &http.Client{Timeout: 30 * time.Second},
            BaseURL: "https://api.example.com",
        }, nil
    },
    OnReset: func(client *HTTPClient) error {
        // Reset client state if needed
        return nil
    },
}
```

**File Handle Pool**:
```go
config := resourcepool.ResourcePoolConfig[*os.File]{
    MaxSize: 10,
    OnCreate: func() (*os.File, error) {
        return os.OpenFile("temp.txt", os.O_RDWR, 0644)
    },
    OnDestroy: func(file *os.File) {
        file.Close()
    },
}
```

### Use Cases
- **Database connections**: Connection pooling with lifecycle management
- **HTTP clients**: Reuse configured clients for external APIs
- **File handles**: Manage limited file descriptors
- **Network connections**: TCP/UDP connection reuse
- **Expensive objects**: Any costly-to-create stateful resources

### Performance
- **Resource limits**: Enforced MaxSize prevents exhaustion
- **Lifecycle efficiency**: Proper resource reuse and cleanup
- **Context support**: Timeout-based resource acquisition
- **Statistics**: Monitor usage patterns and efficiency

---

## 🎯 Choosing the Right Pool

### Decision Matrix

| Scenario | Recommended Pool | Reason |
|----------|------------------|---------|
| **Execute background tasks** | WorkerPool | Concurrent execution with worker management |
| **Process large datasets** | MemoryPool + WorkerPool | Combine memory reuse with parallel processing |
| **Web API backend** | WorkerPool | Handle requests concurrently with auto-scaling |
| **Database operations** | ResourcePool | Manage expensive DB connections |
| **Stream processing** | MemoryPool (Local) | High-performance buffer reuse |
| **File processing** | ResourcePool | Manage file handles safely |
| **Batch data processing** | All three | Workers for concurrency, memory for buffers, resources for DB/files |


---

## 🔧 Integration Examples

### Web Server with All Pools
```go
// Initialize all pools
workerpool.MustInit(workerpool.PoolConfig{
    MinWorkers: 10,
    MaxWorkers: 100,
})

memoryManager := memorypool.NewPoolManager(&memorypool.PoolConfig{
    PoolSize:   20,
    PoolLength: 4096,
    Global:     true, // Shared across requests
})

dbPool, _ := resourcepool.NewResourcePool(dbConfig)
defer dbPool.Close()

// HTTP handler
func handleRequest(w http.ResponseWriter, r *http.Request) {
    // Use WorkerPool for async processing
    workerpool.Async(func() {
        // Get memory buffer
        bufIndex, buffer := memoryManager.AcquireBytePool(4096)
        if bufIndex != utils.NotAvailable {
            defer memoryManager.ReleaseBytePool(bufIndex)
            
            // Get database connection
            ctx := context.Background()
            db, err := dbPool.Get(ctx)
            if err == nil {
                defer dbPool.Put(db)
                
                // Process request with all resources
                processRequest(db, buffer, r)
            }
        }
    })
}
```

### Data Processing Pipeline
```go
func ProcessDataPipeline(data [][]byte) error {
    ctx := context.Background()
    
    // Create processing tasks
    var tasks []func() error
    for _, chunk := range data {
        chunk := chunk // Capture loop variable
        tasks = append(tasks, func() error {
            // Get memory buffer
            bufIndex, buffer := memoryManager.AcquireBytePool(len(chunk))
            if bufIndex == utils.NotAvailable {
                return errors.New("no buffer available")
            }
            defer memoryManager.ReleaseBytePool(bufIndex)
            
            // Get database connection for results
            db, err := dbPool.Get(ctx)
            if err != nil {
                return err
            }
            defer dbPool.Put(db)
            
            // Process chunk
            copy(buffer, chunk)
            return processChunk(db, buffer[:len(chunk)])
        })
    }
    
    // Execute all tasks in parallel
    errors := workerpool.Batch(ctx, tasks...)
    
    // Check for errors
    for _, err := range errors {
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

---

## 📊 Monitoring and Maintenance

### Pool Statistics
```go
// WorkerPool stats
workerStats := workerpool.GetPoolStats()
fmt.Printf("Workers - Running: %d, Capacity: %d, Tasks: %d completed\n",
    workerStats.Running, workerStats.Capacity, workerStats.Completed)

// MemoryPool stats
memStats := memoryManager.GetPoolStats()
fmt.Printf("Memory - String pools used: %d, Byte pools used: %d\n",
    memStats.UsedStringPools, memStats.UsedBytePools)

// ResourcePool stats
resStats := dbPool.GetPoolStats()
fmt.Printf("Resources - Created: %d, Used: %d, Available: %d\n",
    resStats.Created, resStats.Used, resStats.Available)
```

### Maintenance Operations
```go
// Memory pool maintenance
memoryManager.ShrinkAllPools()    // Shrink oversized pools
memoryManager.TestPoolLeak()      // Check for memory leaks

// Resource pool graceful shutdown
err := dbPool.CloseWithTimeout(30 * time.Second)
if err != nil {
    log.Printf("Forced resource pool shutdown: %v", err)
}

// Worker pool graceful shutdown
err = workerpool.ShutdownWithTimeout(10 * time.Second)
if err != nil {
    log.Printf("Worker pool shutdown timeout: %v", err)
}
```

---

## 🏗️ Best Practices

### 1. **Pool Lifecycle Management**
```go
// Initialize pools at application startup
func initPools() {
    workerpool.MustInit(defaultWorkerConfig)
    memoryManager = memorypool.NewPoolManager(defaultMemoryConfig)
    resourcePool, _ = resourcepool.NewResourcePool(defaultResourceConfig)
}

// Graceful shutdown in reverse order
func shutdownPools() {
    resourcePool.CloseWithTimeout(10 * time.Second)
    memoryManager.ShrinkAllPools()
    workerpool.ShutdownWithTimeout(10 * time.Second)
}
```

### 2. **Error Handling**
```go
// Always check for pool exhaustion
index, buffer := memoryManager.AcquireBytePool(size)
if index == utils.NotAvailable {
    // Handle pool exhaustion - use fallback or retry
    return errors.New("memory pool exhausted")
}
defer memoryManager.ReleaseBytePool(index)

// Context timeouts for resource pools
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

resource, err := resourcePool.Get(ctx)
if err != nil {
    // Handle timeout or other errors
    return fmt.Errorf("resource acquisition failed: %w", err)
}
defer resourcePool.Put(resource)
```

### 3. **Performance Optimization**
```go
// Use local memory pools for high-performance paths
localMemory := memorypool.NewPoolManager(&memorypool.PoolConfig{
    Global: false, // No locking overhead
})

// Pre-warm resource pools
for i := 0; i < maxSize; i++ {
    resource, err := resourcePool.Get(context.Background())
    if err == nil {
        resourcePool.Put(resource) // Return immediately
    }
}

// Batch operations for efficiency
var tasks []func() error
for _, item := range items {
    // Create tasks...
}
errors := workerpool.Batch(ctx, tasks...)
```

### 4. **Monitoring Integration**
```go
// Periodic health checks
func monitorPools() {
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()
    
    for range ticker.C {
        workerStats := workerpool.GetPoolStats()
        memStats := memoryManager.GetPoolStats()
        resStats := resourcePool.GetPoolStats()
        
        // Log metrics
        log.Printf("Pool health - Workers: %d/%d, Memory: %d pools, Resources: %d/%d",
            workerStats.Running, workerStats.Capacity,
            memStats.UsedStringPools + memStats.UsedBytePools,
            resStats.Used, resStats.MaxSize)
        
        // Alert on high utilization
        if workerStats.Running > workerStats.Capacity*8/10 {
            log.Printf("HIGH WORKER UTILIZATION: %d/%d", 
                workerStats.Running, workerStats.Capacity)
        }
    }
}
```

---

## 🧪 Testing

Each pool type includes comprehensive test coverage:

```bash
# Test all pools
go test ./pool/...

# Test with race detection
go test -race ./pool/...

# Run benchmarks
go test -bench=. ./pool/...

# Test specific pool types
go test ./pool/workerpool/
go test ./pool/memorypool/
go test ./pool/resourcepool/

# Run stress tests
go test -run Stress ./pool/...
```

---

## 🎯 Summary

The Pool package provides three complementary pooling solutions:

- **WorkerPool**: For concurrent task execution with auto-scaling and panic recovery
- **MemoryPool**: For efficient memory reuse with type safety and leak detection  
- **ResourcePool**: For managing expensive, stateful resources with lifecycle control

Choose the appropriate pool type based on your specific needs, or combine multiple pools for comprehensive resource management in complex applications. Each pool is designed for production use with extensive testing, monitoring capabilities, and graceful shutdown support.