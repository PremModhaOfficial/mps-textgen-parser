# Utils Package

The `utils` package provides common error types, type constraints, thread-safe collections, and helper utilities for the events messaging system.

## Overview

This package contains:

- **Error Types**: Categorized errors for messaging operations (connection, publish, subscribe, auth, tenant, config, JetStream)
- **Error Utilities**: Wrapping, chaining, collection, and classification helpers
- **Type Constraints**: Generic type constraints for Go generics (Signed, Unsigned, Integer, Float, Numeric, Ordered)
- **Optional/Result**: Functional types for safe value handling
- **Thread-Safe Collections**: SafeMap and SafeSlice for concurrent access
- **Panic Recovery**: SafeGo and RecoverError for goroutine safety

## Package Structure

```
utils/
├── errors.go           # Error types and sentinel errors
├── types.go            # Type constraints, Optional, Result, SafeMap, SafeSlice
├── const.go            # Package constants
├── utils.go            # Helper functions
├── errors_core_test.go # Error handling tests
├── utils_test.go       # Utility tests
└── README.md           # This file
```

## Error Types

### Sentinel Errors

All sentinel errors are predefined and can be checked with `errors.Is()`.

#### Connection Errors

| Error | Description |
|-------|-------------|
| `ErrNotConnected` | Operation attempted on disconnected client |
| `ErrAlreadyConnected` | Duplicate connection attempt |
| `ErrConnectionClosed` | Connection was closed |
| `ErrConnectionTimeout` | Connection timed out |
| `ErrReconnectFailed` | Reconnection attempts exhausted |

#### Publishing Errors

| Error | Description |
|-------|-------------|
| `ErrPublishFailed` | Message publish failed |
| `ErrPublishTimeout` | Publish acknowledgment timeout |
| `ErrNoAck` | No acknowledgment received |
| `ErrDuplicateMsg` | Duplicate message detected |
| `ErrStreamNotFound` | Target stream not found |
| `ErrInvalidSubject` | Invalid subject format |
| `ErrInvalidMessage` | Invalid message format |
| `ErrMessageTooLarge` | Message exceeds size limit |
| `ErrRequestTimeout` | Request-reply timeout |
| `ErrNoReply` | No reply received |

#### Subscription Errors

| Error | Description |
|-------|-------------|
| `ErrSubscriptionClosed` | Subscription was closed |
| `ErrSubscriptionInvalid` | Invalid subscription |
| `ErrMaxMsgsExceeded` | Maximum messages exceeded |

#### Authentication Errors

| Error | Description |
|-------|-------------|
| `ErrAuthFailed` | Authentication failed |
| `ErrAuthExpired` | Credentials expired |
| `ErrInvalidCredential` | Invalid credential format |
| `ErrPermissionDenied` | Insufficient permissions |

#### Tenant Errors

| Error | Description |
|-------|-------------|
| `ErrTenantNotFound` | Tenant not registered |
| `ErrTenantExists` | Tenant already exists |
| `ErrTenantDisconnected` | Tenant connection lost |

#### Configuration Errors

| Error | Description |
|-------|-------------|
| `ErrInvalidConfig` | Invalid configuration |
| `ErrMissingConfig` | Required configuration missing |

#### Serialization Errors

| Error | Description |
|-------|-------------|
| `ErrSerializationFailed` | Serialization failed |
| `ErrDeserializationFailed` | Deserialization failed |

#### Lifecycle Errors

| Error | Description |
|-------|-------------|
| `ErrShutdownInProgress` | System is shutting down |
| `ErrAlreadyStarted` | Component already started |
| `ErrNotStarted` | Component not started |
| `ErrJetStreamNotEnabled` | JetStream not enabled on server |

#### Utility Errors

| Error | Description |
|-------|-------------|
| `ErrInvalidArgument` | Invalid function argument |
| `ErrNilValue` | Nil value provided |
| `ErrEmptyValue` | Empty value provided |
| `ErrOutOfRange` | Value out of allowed range |
| `ErrNotFound` | Resource not found |
| `ErrAlreadyExists` | Resource already exists |
| `ErrOperationCanceled` | Operation was canceled |
| `ErrTimeout` | Operation timed out |
| `ErrClosed` | Resource is closed |

### Error Handling

```go
import (
    "errors"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"
)

err := publisher.Publish(ctx, subject, msg)
if err != nil {
    switch {
    case errors.Is(err, utils.ErrNotConnected):
        // Reconnect and retry
    case errors.Is(err, utils.ErrPublishTimeout):
        // Retry with backoff
    case errors.Is(err, utils.ErrInvalidSubject):
        // Fix subject format
    case errors.Is(err, utils.ErrShutdownInProgress):
        // System is shutting down, stop operations
    default:
        // Handle unknown error
    }
}
```

### Custom Error Types

#### Error (with operation context)

```go
// Create error with component and operation context
err := utils.NewError("publisher", "publish", utils.ErrPublishTimeout)
// Error: "publisher publish: publish timeout"

// Wrap an existing error with context
wrapped := utils.WrapError("connection", "connect", originalErr)
wrapped = utils.WrapErrorf("tenant %s", tenantID, originalErr)
```

#### TenantError

```go
// Create tenant-specific error
err := utils.NewTenantError("tenant-123", utils.ErrConnectionClosed)
fmt.Println(err) // "tenant tenant-123: connection closed"
fmt.Println(err.TenantID) // "tenant-123"
errors.Is(err, utils.ErrConnectionClosed) // true
```

#### ValidationError

```go
// Create validation error
err := utils.NewValidationError("port", "must be between 1 and 65535")
err = utils.NewValidationErrorWithValue("port", "must be positive", -1)
```

#### ConfigError

```go
// Configuration-specific error with field context
type ConfigError struct {
    Field   string
    Message string
    Err     error
}
```

#### SerializationError

```go
// Serialization/deserialization error with format context
type SerializationError struct {
    Format  string  // "json", "protobuf", etc.
    Message string
    Err     error
}
```

#### MultiError

```go
// Collect multiple errors
multi := utils.NewMultiError()
multi.Add(err1)
multi.Add(err2)

if multi.HasErrors() {
    fmt.Printf("%d errors: %v\n", multi.Count(), multi)
}
```

### Error Classification

```go
// Check if an error is retryable
if utils.IsRetryable(err) {
    // Safe to retry: timeouts, temporary failures
}

// Check if an error is temporary
if utils.IsTemporary(err) {
    // Error is temporary and may resolve
}
```

### Error Chain Utilities

```go
// Check error type with generics
if utils.Is[*utils.TenantError](err) {
    // Error is a TenantError
}

// Extract typed error from chain
tenantErr, ok := utils.As[*utils.TenantError](err)

// Check if error matches any of the given errors
if utils.IsAnyOf(err, utils.ErrPublishTimeout, utils.ErrNotConnected) {
    // Matches one of the listed errors
}

// Get full error chain
chain := utils.ErrorChain(err)

// Get root cause
root := utils.RootCause(err)

// Format error for logging
formatted := utils.FormatError(err)
```

### Error Collector

```go
collector := utils.NewErrorCollector()

// Collect errors from multiple operations
collector.Add(operation1())
collector.Add(operation2())
collector.AddAll(operation3(), operation4())

if collector.HasErrors() {
    fmt.Printf("%d errors occurred\n", collector.Count())
    for _, err := range collector.All() {
        fmt.Println(err)
    }
}
```

### Panic Recovery

```go
// Recover from panic as error
err := utils.RecoverError(func() error {
    // Code that might panic
    return riskyOperation()
})

// Run goroutine with panic recovery
utils.SafeGo(func() {
    // Panic-safe goroutine
    processMessages()
})
```

## Type Constraints

Generic type constraints for Go generics:

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"

// Signed integer types (int, int8, int16, int32, int64)
type Signed interface { ~int | ~int8 | ~int16 | ~int32 | ~int64 }

// Unsigned integer types (uint, uint8, uint16, uint32, uint64, uintptr)
type Unsigned interface { ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr }

// All integer types
type Integer interface { Signed | Unsigned }

// Floating-point types (float32, float64)
type Float interface { ~float32 | ~float64 }

// All numeric types
type Numeric interface { Integer | Float }

// Types that support ordering (numeric + string)
type Ordered interface { Numeric | ~string }
```

### Usage

```go
func Max[T utils.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}

func Sum[T utils.Numeric](values []T) T {
    var total T
    for _, v := range values {
        total += v
    }
    return total
}
```

## Optional Type

Safe value wrapper that eliminates nil checks:

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"

// Create Optional values
opt := utils.Some("hello")   // Has value
empty := utils.None[string]() // No value

// Check and extract
if opt.HasValue() {
    fmt.Println(opt.Value()) // "hello"
}

// Get with default
val := empty.ValueOr("default") // "default"
```

## Result Type

Operation result that may succeed or fail:

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"

// Create Result values
ok := utils.Ok(42)
fail := utils.Err[int](errors.New("failed"))

// Check and extract
if ok.IsOk() {
    fmt.Println(ok.Value()) // 42
}

if fail.IsErr() {
    fmt.Println(fail.Error()) // "failed"
}
```

## Thread-Safe Collections

### SafeMap

Thread-safe map for concurrent access:

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"

// Create SafeMap
m := utils.NewSafeMap[string, int]()
m = utils.NewSafeMapWithCapacity[string, int](100)

// Operations
m.Set("key", 42)
val, ok := m.Get("key") // 42, true
m.Delete("key")

// Check existence
exists := m.Has("key")

// Get size
count := m.Len()

// Iterate safely
m.Range(func(key string, value int) bool {
    fmt.Printf("%s: %d\n", key, value)
    return true // continue iteration
})

// Get all keys/values
keys := m.Keys()
values := m.Values()
```

### SafeSlice

Thread-safe slice for concurrent access:

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"

// Create SafeSlice
s := utils.NewSafeSlice[string]()
s = utils.NewSafeSliceWithCapacity[string](100)

// Operations
s.Append("item1")
s.Append("item2")

val := s.Get(0) // "item1"
count := s.Len()

// Get snapshot
items := s.Slice()

// Iterate safely
s.Range(func(index int, value string) bool {
    fmt.Printf("[%d]: %s\n", index, value)
    return true
})
```

## Pair Type

Simple key-value pair:

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"

pair := utils.Pair[string, int]{Key: "count", Value: 42}
fmt.Println(pair.Key, pair.Value) // "count" 42
```

## Step-by-Step Usage Guide

### Step 1: Import the Package

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"
```

### Step 2: Use Sentinel Errors for Error Handling

```go
conn, err := client.Connect(ctx, tenantID, creds)
if err != nil {
    if errors.Is(err, utils.ErrConnectionTimeout) {
        // Retry connection
    }
    if errors.Is(err, utils.ErrAuthFailed) {
        // Check credentials
    }
}
```

### Step 3: Wrap Errors with Context

```go
func connectTenant(tenantID string) error {
    conn, err := nats.Connect(url)
    if err != nil {
        return utils.WrapErrorf("tenant %s: connect failed", tenantID, err)
    }
    return nil
}
```

### Step 4: Use Error Collector for Multi-Step Operations

```go
func shutdown(ctx context.Context) error {
    collector := utils.NewErrorCollector()
    collector.Add(publisher.Close(ctx))
    collector.Add(subscriber.Close(ctx))
    collector.Add(connection.Close(ctx))

    if collector.HasErrors() {
        return collector.Error()
    }
    return nil
}
```

### Step 5: Use SafeMap for Concurrent State

```go
connections := utils.NewSafeMap[string, *Connection]()

// From multiple goroutines:
connections.Set(tenantID, conn)
conn, ok := connections.Get(tenantID)
connections.Delete(tenantID)
```

### Step 6: Use Optional for Nullable Values

```go
func findService(id string) utils.Optional[*ServiceInfo] {
    svc, ok := registry.Get(id)
    if !ok {
        return utils.None[*ServiceInfo]()
    }
    return utils.Some(svc)
}

result := findService("user-service")
if result.HasValue() {
    processService(result.Value())
}
```

## Best Practices

### 1. Use Sentinel Errors Over String Comparison

```go
// GOOD
if errors.Is(err, utils.ErrNotConnected) { ... }

// BAD
if err.Error() == "not connected" { ... }
```

### 2. Wrap Errors to Preserve Context

```go
// GOOD - preserves error chain
return utils.WrapError("tenant", "connect", err)

// BAD - loses original error
return fmt.Errorf("connect failed")
```

### 3. Use SafeMap Instead of mutex + map

```go
// GOOD
connections := utils.NewSafeMap[string, *Connection]()

// BAD
var mu sync.RWMutex
connections := make(map[string]*Connection)
```

### 4. Use SafeGo for Background Goroutines

```go
// GOOD - panic-safe
utils.SafeGo(func() {
    processMessages()
})

// BAD - panic crashes the application
go processMessages()
```

### 5. Check Retryability Before Retrying

```go
for attempt := 0; attempt < maxRetries; attempt++ {
    err := publish(msg)
    if err == nil {
        return nil
    }
    if !utils.IsRetryable(err) {
        return err // Non-retryable, fail fast
    }
    time.Sleep(backoff(attempt))
}
```

## Thread Safety

All collection types are thread-safe:
- `SafeMap` - Uses `sync.RWMutex` for concurrent reads, exclusive writes
- `SafeSlice` - Uses `sync.RWMutex` for concurrent reads, exclusive writes
- `ErrorCollector` - Uses `sync.Mutex` for concurrent error collection

## Testing

```bash
# Run tests
go test ./events/utils/...

# Run with coverage
go test -cover ./events/utils/...

# Run with race detection
go test -race ./events/utils/...
```

## Usage in Other Packages

This package is used by:
- `events/core` - Context helpers and interfaces
- `events/auth` - Authentication error types
- `events/tenant` - Tenant error types
- `events/middleware` - Error classification for retry
- `events/endpoint` - Endpoint error types
- `events/microservice` - Microservice error types
- `events` (main) - Re-exports common errors
