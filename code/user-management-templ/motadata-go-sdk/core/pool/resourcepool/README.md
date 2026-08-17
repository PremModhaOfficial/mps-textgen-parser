# Resource Pool Package

## Overview

This package provides efficient resource pooling solutions for managing expensive, stateful resources like database connections, HTTP clients, file handles, and other reusable objects. The **ResourcePool** offers thread-safe resource management with lifecycle callbacks and graceful shutdown capabilities.

**Key Features:**
- **Type-safe resource pooling** for any resource type using Go generics
- **Lifecycle management** with onCreate, onReset, and onDestroy callbacks
- **Thread-safe operations** with atomic counters and mutex protection
- **Context-aware resource acquisition** with timeout and cancellation support
- **Graceful shutdown** with resource return waiting and cleanup
- **Resource leak detection** and usage monitoring

## Quick Decision Guide

```
What type of resource pooling do you need?

DATABASE CONNECTIONS → Use ResourcePool with DB connection lifecycle management
HTTP CLIENTS → Use ResourcePool with client reuse and reset capabilities
FILE HANDLES → Use ResourcePool with proper file lifecycle management  
EXPENSIVE OBJECTS → Use ResourcePool for any costly-to-create stateful resources
CUSTOM RESOURCES → Use ResourcePool with your own create/reset/destroy logic
```

## Resource Pool Architecture

### ResourcePool[T] - The Core Implementation
The `ResourcePool[T]` struct provides comprehensive resource management:
- **Channel-based resource storage** for efficient get/put operations
- **Atomic counters** for tracking created and used resources
- **Lifecycle callbacks** for resource creation, reset, and destruction
- **Thread-safe operations** with proper synchronization
- **Context support** for timeouts and cancellation
- **Graceful shutdown** with resource cleanup

### Supported Operations
- `Get(context)` - Acquire resource with context support (blocking)
- `TryGet()` - Non-blocking resource acquisition
- `Put(resource)` - Return resource to pool with reset
- `Close()` - Immediate pool shutdown
- `CloseWithTimeout(timeout)` - Graceful shutdown with timeout
- `GetPoolStats()` - Get pool statistics and monitoring data

## Configuration and Creation

### Basic Pool Configuration

```go
// Basic resourcepool pool configuration
config := ResourcePoolConfig[*MyResource]{
    MaxSize:  10,    // Maximum number of resources in pool
    OnCreate: func() (*MyResource, error) {
        // Create new resourcepool instance
        return &MyResource{
            ID:        generateID(),
            CreatedAt: time.Now(),
        }, nil
    },
    OnReset: func(resource *MyResource) error {
        // Reset resourcepool state for reuse (optional)
        resource.Reset()
        return nil
    },
    OnDestroy: func(resource *MyResource) {
        // Cleanup resourcepool on destruction (optional)
        resource.Close()
    },
}

pool, err := NewResourcePool(config)
if err != nil {
    log.Fatal("Failed to create resourcepool pool:", err)
}
defer pool.Close()
```

## Usage Patterns

### Database Connection Pool

```go
import "motadatagosdk/core/pool/resourcepool"

// Database connection wrapper
type DBConnection struct {
    Conn      *sql.DB
    ID        int
    CreatedAt time.Time
    InTx      bool
}

func (db *DBConnection) Reset() {
    db.InTx = false
    // Reset any connection state
}

func (db *DBConnection) Close() {
    if db.Conn != nil {
        db.Conn.Close()
    }
}

// Create database connection pool
func NewDBConnectionPool(dsn string, maxSize int) (*resource.ResourcePool[*DBConnection], error) {
    connID := int32(0)
    
    config := resource.ResourcePoolConfig[*DBConnection]{
        MaxSize: maxSize,
        OnCreate: func() (*DBConnection, error) {
            db, err := sql.Open("postgres", dsn)
            if err != nil {
                return nil, err
            }
            
            // Test connection
            if err := db.Ping(); err != nil {
                db.Close()
                return nil, err
            }
            
            id := atomic.AddInt32(&connID, 1)
            return &DBConnection{
                Conn:      db,
                ID:        int(id),
                CreatedAt: time.Now(),
                InTx:      false,
            }, nil
        },
        OnReset: func(conn *DBConnection) error {
            // Reset connection state
            conn.Reset()
            return nil
        },
        OnDestroy: func(conn *DBConnection) {
            // Close database connection
            conn.Close()
        },
    }
    
    return resource.NewResourcePool(config)
}

// Usage example
func ProcessDatabaseQuery(pool *resource.ResourcePool[*DBConnection]) error {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Get database connection from pool
    conn, err := pool.Get(ctx)
    if err != nil {
        return fmt.Errorf("failed to get DB connection: %w", err)
    }
    defer pool.Put(conn) // Always return connection to pool
    
    // Use the connection
    rows, err := conn.Conn.Query("SELECT id, name FROM users WHERE active = $1", true)
    if err != nil {
        return err
    }
    defer rows.Close()
    
    // Process query results...
    for rows.Next() {
        var id int
        var name string
        if err := rows.Scan(&id, &name); err != nil {
            return err
        }
        fmt.Printf("User: %d - %s\n", id, name)
    }
    
    return rows.Err()
}
```

### HTTP Client Pool

```go
// HTTP client wrapper
type HTTPClient struct {
    Client    *http.Client
    ID        int
    CreatedAt time.Time
    BaseURL   string
}

func (c *HTTPClient) Reset() {
    // Reset any client state if needed
}

// Create HTTP client pool
func NewHTTPClientPool(baseURL string, maxSize int) (*resource.ResourcePool[*HTTPClient], error) {
    clientID := int32(0)
    
    config := resource.ResourcePoolConfig[*HTTPClient]{
        MaxSize: maxSize,
        OnCreate: func() (*HTTPClient, error) {
            client := &http.Client{
                Timeout: 30 * time.Second,
                Transport: &http.Transport{
                    MaxIdleConns:        100,
                    MaxIdleConnsPerHost: 10,
                    IdleConnTimeout:     90 * time.Second,
                },
            }
            
            id := atomic.AddInt32(&clientID, 1)
            return &HTTPClient{
                Client:    client,
                ID:        int(id),
                CreatedAt: time.Now(),
                BaseURL:   baseURL,
            }, nil
        },
        OnReset: func(client *HTTPClient) error {
            client.Reset()
            return nil
        },
        OnDestroy: func(client *HTTPClient) {
            // HTTP client doesn't need explicit cleanup
            // Transport will be garbage collected
        },
    }
    
    return resource.NewResourcePool(config)
}

// Usage example
func MakeAPIRequest(pool *resource.ResourcePool[*HTTPClient], endpoint string) (*http.Response, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    
    // Get HTTP client from pool
    httpClient, err := pool.Get(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get HTTP client: %w", err)
    }
    defer pool.Put(httpClient) // Return client to pool
    
    // Make API request
    url := httpClient.BaseURL + endpoint
    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    resp, err := httpClient.Client.Do(req)
    if err != nil {
        return nil, err
    }
    
    return resp, nil
}
```

### Buffer Pool for Data Processing

```go
import "bytes"

// Buffer wrapper for pooling
type Buffer struct {
    Buffer    *bytes.Buffer
    ID        int
    CreatedAt time.Time
    Capacity  int
}

func (b *Buffer) Reset() {
    b.Buffer.Reset()
}

// Create buffer pool
func NewBufferPool(initialCapacity, maxSize int) (*resource.ResourcePool[*Buffer], error) {
    bufferID := int32(0)
    
    config := resource.ResourcePoolConfig[*Buffer]{
        MaxSize: maxSize,
        OnCreate: func() (*Buffer, error) {
            id := atomic.AddInt32(&bufferID, 1)
            buffer := &bytes.Buffer{}
            buffer.Grow(initialCapacity)
            
            return &Buffer{
                Buffer:    buffer,
                ID:        int(id),
                CreatedAt: time.Now(),
                Capacity:  initialCapacity,
            }, nil
        },
        OnReset: func(buffer *Buffer) error {
            buffer.Reset()
            return nil
        },
        OnDestroy: func(buffer *Buffer) {
            // No special cleanup needed for buffers
        },
    }
    
    return resource.NewResourcePool(config)
}

// Usage example
func ProcessData(pool *resource.ResourcePool[*Buffer], data []byte) ([]byte, error) {
    // Get buffer from pool (non-blocking)
    buffer, err := pool.TryGet()
    if err != nil {
        return nil, fmt.Errorf("failed to get buffer: %w", err)
    }
    defer pool.Put(buffer) // Return buffer to pool
    
    // Process data using buffer
    buffer.Buffer.Write(data)
    buffer.Buffer.WriteString(" - processed")
    
    // Return processed data
    result := make([]byte, buffer.Buffer.Len())
    copy(result, buffer.Buffer.Bytes())
    
    return result, nil
}
```

## Resource Lifecycle Management

### Resource Creation and Initialization

```go
config := ResourcePoolConfig[*MyResource]{
    MaxSize: 5,
    OnCreate: func() (*MyResource, error) {
        // Expensive resourcepool creation
        resource, err := createExpensiveResource()
        if err != nil {
            return nil, fmt.Errorf("resourcepool creation failed: %w", err)
        }
        
        // Initialize resourcepool
        if err := resource.Initialize(); err != nil {
            resource.Close() // Cleanup on failure
            return nil, fmt.Errorf("resourcepool initialization failed: %w", err)
        }
        
        return resource, nil
    },
}
```

### Resource Reset and State Management

```go
config := ResourcePoolConfig[*MyResource]{
    OnReset: func(resource *MyResource) error {
        // Reset resourcepool to clean state for reuse
        resource.ClearState()
        
        // Verify resourcepool is still valid
        if !resource.IsValid() {
            return errors.New("resourcepool is no longer valid")
        }
        
        // Reset counters, flags, etc.
        resource.ResetCounters()
        
        return nil
    },
}
```

### Resource Cleanup and Destruction

```go
config := ResourcePoolConfig[*MyResource]{
    OnDestroy: func(resource *MyResource) {
        // Graceful resourcepool cleanup
        resource.Shutdown()
        
        // Release external resources
        resource.ReleaseHandles()
        
        // Close connections, files, etc.
        resource.Close()
    },
}
```

## Context-Aware Resource Management

### Timeout-Based Resource Acquisition

```go
func GetResourceWithTimeout(pool *resource.ResourcePool[*MyResource], timeout time.Duration) (*MyResource, error) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    
    resource, err := pool.Get(ctx)
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            return nil, fmt.Errorf("timeout acquiring resourcepool after %v", timeout)
        }
        return nil, fmt.Errorf("failed to acquire resourcepool: %w", err)
    }
    
    return resource, nil
}
```

### Cancellation-Aware Operations

```go
func ProcessWithCancellation(pool *resource.ResourcePool[*MyResource], ctx context.Context) error {
    // Get resourcepool with cancellation support
    resource, err := pool.Get(ctx)
    if err != nil {
        if errors.Is(err, context.Canceled) {
            return fmt.Errorf("operation was cancelled")
        }
        return err
    }
    defer pool.Put(resource)
    
    // Process with context cancellation checks
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
        return resource.ProcessData()
    }
}
```

## Pool Monitoring and Statistics

### Basic Pool Statistics

```go
// Monitor pool usage
func MonitorPool(pool *resource.ResourcePool[*MyResource]) {
    stats := pool.GetPoolStats()
    
    fmt.Printf("Pool Statistics:\n")
    fmt.Printf("  Created:   %d\n", stats.Created)
    fmt.Printf("  Used:      %d\n", stats.Used)
    fmt.Printf("  Available: %d\n", stats.Available)
    fmt.Printf("  MaxSize:   %d\n", stats.MaxSize)
    fmt.Printf("  IsClosed:  %t\n", stats.IsClosed)
    
    // Calculate efficiency metrics
    efficiency := float64(stats.Available) / float64(stats.MaxSize) * 100
    utilization := float64(stats.Used) / float64(stats.Created) * 100
    
    fmt.Printf("  Pool Efficiency: %.1f%%\n", efficiency)
    fmt.Printf("  Utilization:     %.1f%%\n", utilization)
}
```

### Continuous Pool Monitoring

```go
func StartPoolMonitoring(pool *resource.ResourcePool[*MyResource], interval time.Duration) {
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-ticker.C:
            stats := pool.GetPoolStats()
            
            // Log pool metrics
            log.Printf("Pool stats - Created: %d, Used: %d, Available: %d, Closed: %t",
                stats.Created, stats.Used, stats.Available, stats.IsClosed)
            
            // Alert on high utilization
            if stats.Used > stats.MaxSize*8/10 { // 80% utilization
                log.Printf("WARNING: High pool utilization: %d/%d resources in use",
                    stats.Used, stats.MaxSize)
            }
            
            // Alert on pool closure
            if stats.IsClosed {
                log.Printf("Pool has been closed")
                return
            }
            
        case <-time.After(time.Minute * 5):
            log.Printf("Pool monitoring stopped due to inactivity")
            return
        }
    }
}
```

## Graceful Shutdown and Cleanup

### Immediate Pool Shutdown

```go
func ShutdownPool(pool *resource.ResourcePool[*MyResource]) {
    log.Printf("Shutting down resourcepool pool...")
    
    // Get final statistics
    stats := pool.GetPoolStats()
    log.Printf("Pool state before shutdown - Created: %d, Used: %d, Available: %d",
        stats.Created, stats.Used, stats.Available)
    
    // Immediate shutdown
    pool.Close()
    
    log.Printf("Resource pool shutdown complete")
}
```

### Graceful Shutdown with Timeout

```go
func GracefulShutdown(pool *resource.ResourcePool[*MyResource], timeout time.Duration) error {
    log.Printf("Starting graceful shutdown with %v timeout...", timeout)
    
    // Check current usage
    stats := pool.GetPoolStats()
    if stats.Used > 0 {
        log.Printf("Waiting for %d resources to be returned...", stats.Used)
    }
    
    // Graceful shutdown with timeout
    err := pool.CloseWithTimeout(timeout)
    if err != nil {
        log.Printf("Graceful shutdown failed: %v", err)
        
        // Force shutdown
        pool.Close()
        return fmt.Errorf("forced shutdown after timeout: %w", err)
    }
    
    log.Printf("Graceful shutdown completed successfully")
    return nil
}
```

## Error Handling and Recovery

### Resource Creation Error Handling

```go
func CreateRobustPool() (*resource.ResourcePool[*MyResource], error) {
    config := resource.ResourcePoolConfig[*MyResource]{
        MaxSize: 5,
        OnCreate: func() (*MyResource, error) {
            // Retry logic for resourcepool creation
            var lastErr error
            for attempt := 1; attempt <= 3; attempt++ {
                resource, err := tryCreateResource()
                if err == nil {
                    log.Printf("Resource created successfully on attempt %d", attempt)
                    return resource, nil
                }
                
                lastErr = err
                log.Printf("Resource creation attempt %d failed: %v", attempt, err)
                
                // Exponential backoff
                if attempt < 3 {
                    time.Sleep(time.Duration(attempt) * time.Second)
                }
            }
            
            return nil, fmt.Errorf("failed to create resourcepool after 3 attempts: %w", lastErr)
        },
    }
    
    return resource.NewResourcePool(config)
}
```

### Resource Reset Error Recovery

```go
func CreatePoolWithResetRecovery() (*resource.ResourcePool[*MyResource], error) {
    config := resource.ResourcePoolConfig[*MyResource]{
        OnReset: func(res *MyResource) error {
            // Try to reset resourcepool
            if err := res.Reset(); err != nil {
                log.Printf("Resource reset failed: %v", err)
                
                // Try recovery
                if err := res.Recover(); err != nil {
                    log.Printf("Resource recovery failed: %v", err)
                    return fmt.Errorf("resourcepool cannot be reused: %w", err)
                }
                
                log.Printf("Resource recovered successfully")
            }
            
            return nil
        },
    }
    
    return resource.NewResourcePool(config)
}
```

## Best Practices

### Pool Configuration
1. **Right-size your pool**: Balance memory usage with resource availability
2. **Handle creation failures**: Implement retry logic and proper error handling
3. **Reset resources properly**: Ensure clean state for resource reuse
4. **Monitor pool usage**: Track statistics and set up alerting for issues
5. **Graceful shutdown**: Allow resources to be returned before pool closure

### Resource Management
1. **Always return resources**: Use defer to ensure resources are returned to pool
2. **Handle pool exhaustion**: Check for errors and implement fallback strategies
3. **Use context timeouts**: Prevent indefinite blocking on resource acquisition
4. **Validate resources**: Check resource state before and after use
5. **Implement proper cleanup**: Ensure resources are properly destroyed

### Performance Optimization
1. **Pre-warm pools**: Create initial resources during application startup
2. **Monitor resource lifecycle**: Track creation, usage, and destruction patterns
3. **Optimize reset operations**: Minimize overhead in resource reset callbacks
4. **Use non-blocking operations**: Use TryGet() when immediate response is needed
5. **Profile memory usage**: Monitor for resource leaks and excessive creation

### Error Handling
1. **Robust creation logic**: Handle transient failures and implement retries
2. **Reset error recovery**: Destroy resources that cannot be reset properly
3. **Context cancellation**: Respect context timeouts and cancellation
4. **Pool state validation**: Check pool state before operations
5. **Graceful degradation**: Implement fallbacks for pool exhaustion scenarios

## Testing

Run tests for the resource pool implementations:

```bash
# Run all tests
go test ./...

# Run with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...

# Run specific resourcepool pool tests
go test -run TestResourcePool ./resourcepool/
go test -run TestResourcePoolConcurrent ./resourcepool/

# Run stress tests
go test -run TestResourcePoolStress ./resourcepool/
```

## Performance Comparison

| Feature | ResourcePool | sync.Pool | Manual Management |
|---------|-------------|-----------|-------------------|
| **Type Safety** | Full (generics) | Type assertions required | Full |
| **Lifecycle Management** | Complete callbacks | Limited | Manual implementation |
| **Resource Limits** | Enforced MaxSize | Unlimited | Manual enforcement |
| **Context Support** | Full support | None | Manual implementation |
| **Statistics** | Built-in detailed stats | None | Manual tracking |
| **Graceful Shutdown** | Built-in | None | Manual implementation |
| **Memory Overhead** | Moderate | Low | Variable |

## Summary

The **ResourcePool** provides comprehensive resource management with these key benefits:

**Core Features:**
- **Type-safe resource pooling** using Go generics
- **Complete lifecycle management** with onCreate, onReset, and onDestroy callbacks
- **Thread-safe operations** with atomic counters and proper synchronization
- **Context-aware operations** with timeout and cancellation support
- **Monitoring and statistics** for pool usage tracking
- **Graceful shutdown** capabilities with resource cleanup

**Use ResourcePool when you need:**
- Management of expensive, stateful resources (DB connections, HTTP clients)
- Resource lifecycle control with creation, reset, and destruction logic
- Pool size limits and resource usage monitoring
- Context-aware operations with timeout and cancellation
- Thread-safe resource sharing across goroutines
- Graceful application shutdown with resource cleanup

**The ResourcePool is ideal for managing any expensive or stateful resources** where you need complete control over the resource lifecycle, usage monitoring, and graceful shutdown capabilities.