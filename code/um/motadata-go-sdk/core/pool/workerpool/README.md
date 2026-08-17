# WorkerPool Package

A high-performance, thread-safe worker pool implementation with auto-scaling capabilities, panic recovery, and comprehensive concurrency management.

## Features

- ✅ **Auto-scaling worker pool** with configurable min/max workers
- ✅ **Synchronous and asynchronous task execution**
- ✅ **Context-based cancellation and timeout support**
- ✅ **Panic recovery** with graceful error handling
- ✅ **Batch processing** for multiple concurrent tasks
- ✅ **Scheduling support** (After, Every) for delayed/recurring tasks
- ✅ **Comprehensive metrics** and monitoring
- ✅ **Thread-safe singleton manager** for global usage
- ✅ **Resource cleanup** and graceful shutdown
- ✅ **Retry mechanisms** with exponential backoff
- ✅ **Type-safe operations** using Go generics

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "time"

    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/core/pool/workerpool"
)

func main() {
    // Initialize the global worker pool
    err := workerpool.Init(workerpool.PoolConfig{
        MinWorkers: 2,
        MaxWorkers: 10,
        Timeout:    30 * time.Second,
    })
    if err != nil {
        panic(err)
    }
    defer workerpool.Shutdown()

    // Fire-and-forget task
    workerpool.Async(func() {
        fmt.Println("Hello from worker pool!")
    })

    // Wait for completion
    time.Sleep(100 * time.Millisecond)
}
```

### Context-based Operations

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// Synchronous execution with context
err := workerpool.SyncWithCtx(ctx, func() error {
    // Your task logic here
    time.Sleep(2 * time.Second)
    return nil
})

// Asynchronous execution with context
workerpool.AsyncWithCtx(ctx, func(c context.Context) {
    select {
    case <-c.Done():
        fmt.Println("Task cancelled")
        return
    default:
        // Your task logic
        fmt.Println("Task completed")
    }
})
```

## Core API

### Pool Management

```go
// Initialize with custom configuration
workerpool.Init(workerpool.PoolConfig{
    MinWorkers: 5,
    MaxWorkers: 50,
    Timeout:    time.Minute,
})

// Or use MustInit for panic on failure
workerpool.MustInit()

// Graceful shutdown
workerpool.Shutdown()

// Shutdown with timeout
err := workerpool.ShutdownWithTimeout(10 * time.Second)
```

### Task Execution

#### Fire-and-Forget
```go
// Simple async task
workerpool.Async(func() {
    fmt.Println("Background task")
})

// With context
ctx := context.Background()
workerpool.AsyncWithCtx(ctx, func(c context.Context) {
    // Task with context
})
```

#### Synchronous Execution
```go
ctx := context.Background()

// Basic synchronous execution
err := workerpool.SyncWithCtx(ctx, func() error {
    // Your logic here
    return nil
})

// With timeout
err = workerpool.SyncWithTimeout(5*time.Second, func() error {
    // Your logic here
    return nil
})
```

### Result Channels

```go
// Get result from async task
resultChan := workerpool.AsyncWithError(func() error {
    return doSomeWork()
})
err := <-resultChan

// Get typed result
valueChan := workerpool.AsyncResult(func() (string, error) {
    return "result", nil
})
result := <-valueChan
if result.Error != nil {
    // Handle error
}
fmt.Println("Value:", result.Value)

// With context
ctx := context.Background()
resultChan = workerpool.AsyncWithCtxError(ctx, func(c context.Context) error {
    return doWorkWithContext(c)
})
```

### Batch Processing

```go
ctx := context.Background()

// Run multiple tasks concurrently
tasks := []func() error{
    func() error { return task1() },
    func() error { return task2() },
    func() error { return task3() },
}

// All tasks run in parallel, wait for all to complete
errors := workerpool.Batch(ctx, tasks...)

// Stop on first error
err := workerpool.Parallel(ctx, tasks...)

// Run sequentially, stop on first error
err = workerpool.Sequential(ctx, tasks...)
```

### Advanced Task Options

```go
ctx := context.Background()

task := workerpool.Task{
    Name:       "complex-task",
    Timeout:    10 * time.Second,
    Retries:    3,
    RetryDelay: time.Second,
    OnSuccess: func() {
        fmt.Println("Task succeeded!")
    },
    OnError: func(err error) {
        fmt.Printf("Task failed: %v\n", err)
    },
    OnRetry: func(attempt int, err error) {
        fmt.Printf("Retry %d: %v\n", attempt, err)
    },
    Fn: func(c context.Context) error {
        // Your complex task logic
        return doComplexWork(c)
    },
}

err := workerpool.Submit(ctx, task)
```

### Scheduling

```go
// Execute once after delay
task := workerpool.After(5*time.Second, func() {
    fmt.Println("Delayed task executed")
})

// Can cancel if needed
cancelled := task.Cancel()

// Execute repeatedly
ctx := context.Background()
recurring := workerpool.Every(ctx, time.Minute, func() {
    fmt.Println("Recurring task")
})

// Cancel recurring task
recurring.Cancel()
```

## Pool Statistics

```go
stats := workerpool.GetStats()
fmt.Printf("Running: %d, Waiting: %d, Capacity: %d\n", 
    stats.Running, stats.Waiting, stats.Capacity)
fmt.Printf("Submitted: %d, Completed: %d, Failed: %d\n",
    stats.Submitted, stats.Completed, stats.Failed)
```

## Direct Pool Usage

For more control, you can create and manage pools directly:

```go
pool, err := workerpool.NewWorkerPool(workerpool.PoolConfig{
    MinWorkers: 3,
    MaxWorkers: 15,
    Timeout:    30 * time.Second,
})
if err != nil {
    return err
}
defer pool.close()

// Use pool directly
err = pool.async(func() {
    // Task logic
})

// With options
ctx := context.Background()
err = pool.syncWithOptions(ctx, func(c context.Context) error {
    return nil
}, 
    workerpool.WithTimeout(5*time.Second),
    workerpool.WithRetry(3, time.Second),
)
```

## Configuration Options

### PoolConfig
```go
type PoolConfig struct {
    MinWorkers int           // Minimum number of workers (default: 2)
    MaxWorkers int           // Maximum number of workers (default: 6)
    Timeout    time.Duration // Default task timeout (default: 30s)
}
```

### Task Options
```go
// Available task options
workerpool.WithName("task-name")
workerpool.WithTimeout(10 * time.Second)
workerpool.WithRetry(3, time.Second)
workerpool.WithOnSuccess(func() { /* success callback */ })
workerpool.WithOnError(func(err error) { /* error callback */ })
workerpool.WithOnRetry(func(attempt int, err error) { /* retry callback */ })
```

## Error Handling

The workerpool provides comprehensive error handling:

```go
// Panic recovery
workerpool.Async(func() {
    panic("This won't crash the pool")
    // Pool continues to function normally
})

// Context cancellation
ctx, cancel := context.WithTimeout(context.Background(), time.Second)
defer cancel()

err := workerpool.SyncWithCtx(ctx, func() error {
    time.Sleep(2 * time.Second) // Will be cancelled
    return nil
})
// err will be context.DeadlineExceeded
```

## Best Practices

### 1. Pool Lifecycle
```go
// Initialize once at application startup
func init() {
    workerpool.MustInit(workerpool.PoolConfig{
        MinWorkers: 5,
        MaxWorkers: 100,
        Timeout:    30 * time.Second,
    })
}

// Graceful shutdown
func shutdown() {
    workerpool.ShutdownWithTimeout(10 * time.Second)
}
```

### 2. Context Usage
```go
// Always use context for long-running tasks
func processRequest(ctx context.Context) error {
    return workerpool.SyncWithCtx(ctx, func() error {
        // Use context in your task
        return doWork(ctx)
    })
}
```

### 3. Batch Processing
```go
// Process large datasets efficiently
func processBatch(items []Item) error {
    ctx := context.Background()
    
    // Create tasks for each item
    var tasks []func() error
    for _, item := range items {
        item := item // Capture loop variable
        tasks = append(tasks, func() error {
            return processItem(item)
        })
    }
    
    // Process all in parallel
    errors := workerpool.Batch(ctx, tasks...)
    
    // Handle errors
    for i, err := range errors {
        if err != nil {
            log.Printf("Item %d failed: %v", i, err)
        }
    }
    
    return nil
}
```

### 4. Resource Management
```go
// Use defer for cleanup
func handleRequest() {
    // Schedule cleanup
    cleanup := workerpool.After(time.Hour, func() {
        cleanupResources()
    })
    
    // Cancel if not needed
    defer cleanup.Cancel()
    
    // Process request...
}
```

## Performance Considerations

- **Auto-scaling**: The pool automatically adjusts worker count based on load
- **Panic Recovery**: Failed tasks don't affect other workers
- **Memory Efficiency**: Workers are reused, minimizing allocation overhead
- **Context Cancellation**: Proper cancellation prevents resource leaks
- **Batch Operations**: Use batch processing for multiple related tasks

## Thread Safety

All operations are thread-safe and can be called from multiple goroutines:

```go
// Safe to call from multiple goroutines
for i := 0; i < 100; i++ {
    go func(id int) {
        workerpool.Async(func() {
            fmt.Printf("Task %d\n", id)
        })
    }(i)
}
```

## Testing

The package includes comprehensive test coverage (88.3%) with tests for:

- ✅ Concurrent operations and race conditions
- ✅ Panic recovery and error handling
- ✅ Context cancellation and timeouts
- ✅ Auto-scaling behavior
- ✅ Resource cleanup and memory leaks
- ✅ Edge cases and error scenarios

## License

This package is part of the Motadata Go SDK. See the main repository for license information.