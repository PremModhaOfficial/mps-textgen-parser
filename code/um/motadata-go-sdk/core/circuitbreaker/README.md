# Circuit Breaker Package - Fault Tolerance

## Overview

The circuit breaker package provides a fault tolerance mechanism for protecting your application from cascading failures when calling external services. It wraps the battle-tested [sony/gobreaker](https://github.com/sony/gobreaker) library with a simplified, idiomatic Go API.

## What is a Circuit Breaker?

A circuit breaker is a design pattern that prevents an application from repeatedly trying to execute an operation that's likely to fail. Like an electrical circuit breaker, it "trips" when failures exceed a threshold and "resets" after a timeout period.

```
┌─────────────────────────────────────────────────────────────────┐
│                    Circuit Breaker States                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│   ┌──────────┐    failures >= threshold    ┌──────────┐        │
│   │  CLOSED  │ ─────────────────────────── │   OPEN   │        │
│   │ (normal) │                             │  (fail   │        │
│   └──────────┘                             │   fast)  │        │
│        ▲                                   └──────────┘        │
│        │                                        │              │
│        │ success                     timeout expires           │
│        │                                        │              │
│        │            ┌───────────┐               │              │
│        └─────────── │ HALF-OPEN │ ◄─────────────┘              │
│                     │  (test)   │                              │
│                     └───────────┘                              │
│                          │                                     │
│                     failure → back to OPEN                     │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## Quick Start

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/core/circuitbreaker"

// Create with default configuration
cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.DefaultConfig("my-service"))

// Execute operations through the circuit breaker
result, err := cb.Execute(func() (any, error) {
    return callExternalService()
})

if err != nil {
    if errors.Is(err, circuitbreaker.ErrCircuitOpen) {
        // Circuit is open - service is unavailable
        return fallbackResponse()
    }
    // Handle other errors
}
```

## States

| State | Description | Behavior |
|-------|-------------|----------|
| **Closed** | Normal operation | Requests pass through, failures are counted |
| **Open** | Circuit tripped | Requests fail immediately with `ErrCircuitOpen` |
| **Half-Open** | Testing recovery | Limited requests allowed to test if service recovered |

## Configuration

### Default Configuration

```go
cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.DefaultConfig("service-name"))
```

Default values:
- `FailureThreshold`: 5 consecutive failures to trip
- `Timeout`: 60 seconds before trying recovery
- `MaxRequests`: 1 request allowed in half-open state
- `Interval`: 0 (failure counts never reset in closed state)

### Custom Configuration

```go
cfg := circuitbreaker.Config{
    // Name identifies this circuit breaker (required)
    Name: "payment-service",

    // Number of consecutive failures to trip the circuit
    FailureThreshold: 5,

    // Time to wait before attempting recovery (entering half-open state)
    Timeout: 30 * time.Second,

    // Maximum requests allowed through in half-open state
    MaxRequests: 3,

    // Interval to reset failure counts in closed state (0 = never reset)
    Interval: 1 * time.Minute,

    // Number of consecutive successes to close circuit from half-open
    SuccessThreshold: 2,

    // Callback when state changes
    OnStateChange: func(name string, from, to circuitbreaker.State) {
        log.Printf("[CircuitBreaker] %s: %s -> %s", name, from, to)
    },

    // Custom logic to determine if an error should count as failure
    IsSuccessful: func(err error) bool {
        // Don't count 4xx errors as failures (client errors)
        var httpErr *HTTPError
        if errors.As(err, &httpErr) && httpErr.StatusCode >= 400 && httpErr.StatusCode < 500 {
            return true // Not a failure
        }
        return false // Count as failure
    },
}

cb := circuitbreaker.NewCircuitBreaker(cfg)
```

## API Reference

### Creating a Circuit Breaker

```go
// With default configuration
cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.DefaultConfig("name"))

// With custom configuration
cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.Config{...})
```

### Executing Operations

```go
// Execute a function through the circuit breaker
result, err := cb.Execute(func() (any, error) {
    // Your operation here
    return doSomething()
})
```

### Checking State

```go
// Get current state
state := cb.State()  // StateClosed, StateHalfOpen, or StateOpen

// Convenience methods
if cb.IsOpen() {
    // Circuit is open
}

if cb.IsClosed() {
    // Circuit is closed (normal operation)
}
```

### Getting Statistics

```go
counts := cb.Counts()

fmt.Printf("Requests: %d\n", counts.Requests)
fmt.Printf("Successes: %d\n", counts.TotalSuccesses)
fmt.Printf("Failures: %d\n", counts.TotalFailures)
fmt.Printf("Consecutive Successes: %d\n", counts.ConsecutiveSuccesses)
fmt.Printf("Consecutive Failures: %d\n", counts.ConsecutiveFailures)
```

### Other Methods

```go
// Get the circuit breaker name
name := cb.Name()
```

## Errors

| Error | Description |
|-------|-------------|
| `ErrCircuitOpen` | Circuit is open, request was rejected |
| `ErrTooManyRequest` | Too many requests in half-open state |

## Usage Patterns

### Basic Error Handling

```go
result, err := cb.Execute(func() (any, error) {
    return client.Call()
})

switch {
case err == nil:
    // Success
    processResult(result)
case errors.Is(err, circuitbreaker.ErrCircuitOpen):
    // Circuit is open - use fallback
    return getFallbackData()
default:
    // Other error from the operation
    return nil, err
}
```

### With Fallback

```go
func callWithFallback() (any, error) {
    result, err := cb.Execute(func() (any, error) {
        return primaryService.Call()
    })

    if errors.Is(err, circuitbreaker.ErrCircuitOpen) {
        // Primary service unavailable, use fallback
        return fallbackService.Call()
    }

    return result, err
}
```

### With Retry (Outside Circuit Breaker)

```go
func callWithRetry(maxRetries int) (any, error) {
    var lastErr error

    for i := 0; i < maxRetries; i++ {
        result, err := cb.Execute(func() (any, error) {
            return service.Call()
        })

        if err == nil {
            return result, nil
        }

        // Don't retry if circuit is open
        if errors.Is(err, circuitbreaker.ErrCircuitOpen) {
            return nil, err
        }

        lastErr = err
        time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
    }

    return nil, lastErr
}
```

### Multiple Circuit Breakers

```go
// Create separate circuit breakers for different services
paymentCB := circuitbreaker.NewCircuitBreaker(circuitbreaker.Config{
    Name:             "payment-service",
    FailureThreshold: 3,
    Timeout:          10 * time.Second,
})

inventoryCB := circuitbreaker.NewCircuitBreaker(circuitbreaker.Config{
    Name:             "inventory-service",
    FailureThreshold: 5,
    Timeout:          30 * time.Second,
})

// Use them independently
paymentResult, _ := paymentCB.Execute(func() (any, error) {
    return paymentService.Process(order)
})

inventoryResult, _ := inventoryCB.Execute(func() (any, error) {
    return inventoryService.Reserve(items)
})
```

### Monitoring State Changes

```go
cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.Config{
    Name:             "monitored-service",
    FailureThreshold: 5,
    Timeout:          30 * time.Second,
    OnStateChange: func(name string, from, to circuitbreaker.State) {
        // Log state changes
        log.Printf("[CircuitBreaker] %s: %s -> %s", name, from, to)

        // Send metrics
        metrics.CircuitBreakerState(name, to.String())

        // Alert on open
        if to == circuitbreaker.StateOpen {
            alerting.Send(fmt.Sprintf("Circuit breaker %s opened!", name))
        }
    },
})
```

## Best Practices

### 1. Choose Appropriate Thresholds

```go
// For critical, fast services (e.g., auth)
cfg := circuitbreaker.Config{
    FailureThreshold: 3,   // Trip quickly
    Timeout:          5 * time.Second,  // Recover quickly
    MaxRequests:      1,   // Conservative recovery
}

// For less critical, slower services (e.g., reporting)
cfg := circuitbreaker.Config{
    FailureThreshold: 10,  // More tolerant
    Timeout:          60 * time.Second, // Longer recovery
    MaxRequests:      5,   // Allow more test requests
}
```

### 2. Use Meaningful Names

```go
// Good - descriptive names
cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.DefaultConfig("payment-gateway-v2"))
cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.DefaultConfig("user-service-auth"))

// Bad - generic names
cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.DefaultConfig("service1"))
```

### 3. Handle Circuit Open Gracefully

```go
result, err := cb.Execute(fn)
if errors.Is(err, circuitbreaker.ErrCircuitOpen) {
    // Don't log as error - this is expected behavior
    log.Debug("Circuit open, using cached data")
    return getCachedData()
}
```

### 4. Don't Wrap Everything

Only use circuit breakers for:
- External service calls (HTTP, gRPC, database)
- Operations that can fail due to external factors
- Operations where fast-fail is beneficial

Don't use for:
- In-memory operations
- Local file operations (usually)
- Operations that should always be attempted

### 5. Consider the Interval Setting

```go
// Use Interval to reset counts periodically
cfg := circuitbreaker.Config{
    FailureThreshold: 5,
    Interval:         1 * time.Minute,  // Reset counts every minute
}
// This prevents occasional failures from accumulating over time
// and eventually tripping the circuit
```

## Thread Safety

The circuit breaker is **thread-safe** and can be safely used from multiple goroutines concurrently.

```go
cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.DefaultConfig("shared-service"))

// Safe to use from multiple goroutines
for i := 0; i < 100; i++ {
    go func() {
        cb.Execute(func() (any, error) {
            return service.Call()
        })
    }()
}
```

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...

# Run specific benchmark
go test -bench=BenchmarkExecute_Success -benchmem ./...
```

## Performance

The circuit breaker adds minimal overhead to your operations:

| Scenario | Overhead |
|----------|----------|
| Closed state (success) | ~100-200ns |
| Closed state (failure) | ~100-200ns |
| Open state (rejection) | ~50-100ns |
| Parallel execution | Scales linearly |

## Summary

- **Use circuit breakers** for external service calls to prevent cascading failures
- **Configure thresholds** based on your service's characteristics
- **Handle `ErrCircuitOpen`** gracefully with fallbacks or cached data
- **Monitor state changes** for observability
- **Keep it simple** - start with defaults and tune as needed