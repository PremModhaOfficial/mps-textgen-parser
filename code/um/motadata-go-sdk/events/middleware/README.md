# Middleware Package

The `middleware` package provides a middleware pipeline for the events messaging system, enabling cross-cutting concerns like tracing, retry, metrics, and logging.

## Overview

This package contains:

- **Middleware Stack**: Composable middleware chain
- **Tracing Middleware**: Distributed tracing with OpenTelemetry
- **Retry Middleware**: Automatic retry with exponential backoff
- **Metrics Middleware**: Request metrics collection
- **Logging Middleware**: Structured request logging

## Package Structure

```
middleware/
├── middleware.go    # Core middleware interfaces and stack
├── tracing.go       # Distributed tracing middleware
├── retry.go         # Retry middleware
├── metrics.go       # Metrics middleware
├── logging.go       # Logging middleware
├── tracing_test.go  # Tracing tests
├── metrics_test.go  # Metrics tests
├── retry_test.go    # Retry tests
└── README.md        # This file
```

## Middleware Types

### PublishMiddleware

Intercepts publish operations:

```go
type PublishMiddleware func(PublishHandler) PublishHandler
type PublishHandler func(ctx context.Context, subject string, msg *core.Message) error
```

### SubscribeMiddleware

Intercepts subscribe operations:

```go
type SubscribeMiddleware func(SubscribeHandler) SubscribeHandler
type SubscribeHandler func(ctx context.Context, msg *core.Message) error
```

### Interceptor

Combined publish and subscribe interception:

```go
type Interceptor interface {
    InterceptPublish() PublishMiddleware
    InterceptSubscribe() SubscribeMiddleware
}
```

## Middleware Stack

### Creating a Stack

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/middleware"

stack := middleware.NewStack()

// Add middleware (executed in order added)
stack.UsePublish(tracingMiddleware.InterceptPublish())
stack.UsePublish(metricsMiddleware.InterceptPublish())
stack.UsePublish(retryMiddleware)
stack.UsePublish(loggingMiddleware.InterceptPublish())

stack.UseSubscribe(tracingMiddleware.InterceptSubscribe())
stack.UseSubscribe(metricsMiddleware.InterceptSubscribe())
stack.UseSubscribe(loggingMiddleware.InterceptSubscribe())
```

### Using Interceptors

```go
// Add both publish and subscribe middleware at once
stack.UseInterceptor(middleware.Tracing())
stack.UseInterceptor(middleware.MetricsMiddleware("my-service"))
stack.UseInterceptor(middleware.Logging())
```

### Wrapping Handlers

```go
// Wrap a publish handler
basePublish := func(ctx context.Context, subject string, msg *core.Message) error {
    return publisher.Publish(ctx, subject, msg)
}
wrappedPublish := stack.WrapPublish(basePublish)

// Use wrapped handler
err := wrappedPublish(ctx, "orders.created", msg)

// Wrap a subscribe handler
baseSubscribe := func(ctx context.Context, msg *core.Message) error {
    return processMessage(msg)
}
wrappedSubscribe := stack.WrapSubscribe(baseSubscribe)
```

### Chaining Middleware

```go
// Chain multiple middleware manually
combined := middleware.Chain(
    tracingMiddleware,
    metricsMiddleware,
    loggingMiddleware,
)

handler := combined(baseHandler)

// For subscribe middleware
combinedSubscribe := middleware.ChainSubscribe(
    tracing.InterceptSubscribe(),
    metrics.InterceptSubscribe(),
    logging.InterceptSubscribe(),
)
```

## Tracing Middleware

Distributed tracing with OpenTelemetry integration:

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/middleware"

// Create tracing middleware
tracing := middleware.Tracing()

// Or with custom configuration
tracing := middleware.NewTracingMiddleware()
tracing.UseOTEL = true  // Enable OpenTelemetry (default)
tracing.Sampler = func() bool { return true }  // Custom sampler
```

### Features

- **OpenTelemetry Integration**: Automatic span creation and propagation
- **W3C Trace Context**: Standard `traceparent` header support
- **Context Propagation**: Trace context injected into message headers
- **Span Attributes**: Messaging system, destination, operation

### Usage

```go
stack := middleware.NewStack()
stack.UseInterceptor(middleware.Tracing())

// Publish - creates producer span
err := wrappedPublish(ctx, "orders.created", msg)
// Span: "events.publish" with messaging.destination="orders.created"

// Subscribe - creates consumer span
err := wrappedSubscribe(ctx, msg)
// Span: "events.receive" with messaging.destination from msg.Subject
```

### Trace Context Headers

| Header | Description |
|--------|-------------|
| `traceparent` | W3C trace context parent |
| `tracestate` | W3C trace state |
| `X-Trace-ID` | Trace identifier |
| `X-Span-ID` | Span identifier |

### W3C Trace Context Support

```go
// Extract trace context from W3C traceparent header
tc := middleware.ExtractW3CTraceParent("00-abc123...-def456...-01")

// Format trace context as W3C traceparent
traceparent := middleware.FormatW3CTraceParent(tc)
// "00-abc123...-def456...-01"
```

## Retry Middleware

Automatic retry with exponential backoff:

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/middleware"

// Create retry middleware with defaults
retry := middleware.Retry(middleware.RetryConfig{
    MaxAttempts:     3,
    InitialInterval: 100 * time.Millisecond,
    MaxInterval:     5 * time.Second,
    Multiplier:      2.0,
    Jitter:          0.1,
})

stack.UsePublish(retry)
```

### RetryConfig

```go
type RetryConfig struct {
    MaxAttempts     int           // Maximum retry attempts (default: 3)
    InitialInterval time.Duration // Initial backoff (default: 100ms)
    MaxInterval     time.Duration // Maximum backoff (default: 5s)
    Multiplier      float64       // Backoff multiplier (default: 2.0)
    Jitter          float64       // Random jitter 0.0-1.0 (default: 0.1)
    RetryableErrors []error       // Errors to retry (default: all)
}
```

### Backoff Calculation

```
interval = min(InitialInterval * (Multiplier ^ attempt), MaxInterval)
interval = interval * (1 + random(-Jitter, +Jitter))
```

**Example with defaults:**
- Attempt 1: ~100ms
- Attempt 2: ~200ms
- Attempt 3: ~400ms
- (Total: ~700ms before final failure)

### Selective Retry

```go
retry := middleware.Retry(middleware.RetryConfig{
    MaxAttempts: 3,
    RetryableErrors: []error{
        core.ErrPublishTimeout,
        core.ErrNotConnected,
    },
})

// Only retry on specified errors
// Other errors fail immediately
```

## Metrics Middleware

Request metrics collection:

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/middleware"

// Create metrics middleware
metrics := middleware.MetricsMiddleware("my-service")

stack.UseInterceptor(metrics)
```

### Collected Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `events_publish_total` | Counter | Total publish operations |
| `events_publish_errors_total` | Counter | Failed publish operations |
| `events_publish_duration_seconds` | Histogram | Publish latency |
| `events_receive_total` | Counter | Total received messages |
| `events_receive_errors_total` | Counter | Failed message processing |
| `events_receive_duration_seconds` | Histogram | Processing latency |

### Labels

- `service`: Service name
- `subject`: Message subject
- `status`: `success` or `error`

### Integration with OTEL Metrics

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/metrics"

// Initialize OTEL metrics first
metrics.Init(cfg)

// Metrics middleware automatically uses OTEL
metricsMiddleware := middleware.MetricsMiddleware("my-service")
```

## Logging Middleware

Structured request logging:

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/middleware"

// Create logging middleware
logging := middleware.Logging()

// Or with custom logger
logging := middleware.LoggingWithLogger(customLogger)

stack.UseInterceptor(logging)
```

### Logged Fields

**Publish:**
- `subject`: Target subject
- `message_id`: Message ID
- `size_bytes`: Message size
- `duration`: Operation duration
- `error`: Error if failed

**Subscribe:**
- `subject`: Source subject
- `message_id`: Message ID
- `size_bytes`: Message size
- `duration`: Processing duration
- `error`: Error if failed

### Log Levels

| Event | Level |
|-------|-------|
| Publish start | Debug |
| Publish success | Debug |
| Publish error | Error |
| Receive start | Debug |
| Receive success | Debug |
| Receive error | Error |

## Complete Example

```go
import (
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/middleware"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

func main() {
    // Create middleware stack
    stack := middleware.NewStack()

    // Add middleware in order (first added = outermost)
    stack.UseInterceptor(middleware.Tracing())           // Tracing first
    stack.UseInterceptor(middleware.MetricsMiddleware("order-service"))
    stack.UsePublish(middleware.Retry(middleware.RetryConfig{
        MaxAttempts:     3,
        InitialInterval: 100 * time.Millisecond,
    }))
    stack.UseInterceptor(middleware.Logging())

    // Wrap handlers
    publish := stack.WrapPublish(func(ctx context.Context, subject string, msg *core.Message) error {
        return natsPublisher.Publish(ctx, subject, msg)
    })

    subscribe := stack.WrapSubscribe(func(ctx context.Context, msg *core.Message) error {
        return processOrder(ctx, msg)
    })

    // Use wrapped handlers
    msg := core.NewMessage([]byte(`{"order_id": "123"}`))
    err := publish(ctx, "orders.created", msg)
    // Middleware chain: Tracing -> Metrics -> Retry -> Logging -> Publish
}
```

## Middleware Execution Order

```
Request Flow (Publish):
┌─────────┐  ┌─────────┐  ┌───────┐  ┌─────────┐  ┌──────────┐
│ Tracing │->│ Metrics │->│ Retry │->│ Logging │->│ Publish  │
└─────────┘  └─────────┘  └───────┘  └─────────┘  └──────────┘
     │            │           │           │             │
     │            │           │           │             │
     │<───────────│<──────────│<──────────│<────────────│
     │            │           │           │             │
Response Flow

Request Flow (Subscribe):
┌──────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐
│ Message  │->│ Tracing │->│ Metrics │->│ Logging │->│ Handler │
└──────────┘  └─────────┘  └─────────┘  └─────────┘  └─────────┘
```

## Custom Middleware

### Simple Middleware

```go
func TimingMiddleware() middleware.PublishMiddleware {
    return func(next middleware.PublishHandler) middleware.PublishHandler {
        return func(ctx context.Context, subject string, msg *core.Message) error {
            start := time.Now()
            err := next(ctx, subject, msg)
            duration := time.Since(start)

            if err != nil {
                log.Printf("Publish to %s failed after %v: %v", subject, duration, err)
            } else {
                log.Printf("Publish to %s succeeded in %v", subject, duration)
            }
            return err
        }
    }
}
```

### Interceptor Implementation

```go
type AuthMiddleware struct {
    tokenValidator func(ctx context.Context) bool
}

func (m *AuthMiddleware) InterceptPublish() middleware.PublishMiddleware {
    return func(next middleware.PublishHandler) middleware.PublishHandler {
        return func(ctx context.Context, subject string, msg *core.Message) error {
            if !m.tokenValidator(ctx) {
                return errors.New("unauthorized")
            }
            return next(ctx, subject, msg)
        }
    }
}

func (m *AuthMiddleware) InterceptSubscribe() middleware.SubscribeMiddleware {
    return func(next middleware.SubscribeHandler) middleware.SubscribeHandler {
        return func(ctx context.Context, msg *core.Message) error {
            // Extract and validate token from headers
            token := msg.Headers.Get("Authorization")
            if token == "" {
                return errors.New("missing authorization")
            }
            return next(ctx, msg)
        }
    }
}

// Usage
stack.UseInterceptor(&AuthMiddleware{tokenValidator: validateToken})
```

### Header Injection Middleware

```go
func TenantHeaderMiddleware(tenantID string) middleware.PublishMiddleware {
    return func(next middleware.PublishHandler) middleware.PublishHandler {
        return func(ctx context.Context, subject string, msg *core.Message) error {
            if msg.Headers == nil {
                msg.Headers = make(core.Headers)
            }
            msg.Headers.Set("X-Tenant-ID", tenantID)
            return next(ctx, subject, msg)
        }
    }
}
```

## Best Practices

### 1. Order Matters

```go
// Recommended order:
stack.UseInterceptor(middleware.Tracing())    // 1. Tracing (outermost)
stack.UseInterceptor(middleware.Metrics())    // 2. Metrics
stack.UsePublish(middleware.Retry(cfg))       // 3. Retry (publish only)
stack.UseInterceptor(middleware.Logging())    // 4. Logging (innermost)
```

### 2. Don't Panic in Middleware

```go
func SafeMiddleware() middleware.PublishMiddleware {
    return func(next middleware.PublishHandler) middleware.PublishHandler {
        return func(ctx context.Context, subject string, msg *core.Message) error {
            defer func() {
                if r := recover(); r != nil {
                    log.Printf("Panic in middleware: %v", r)
                }
            }()
            return next(ctx, subject, msg)
        }
    }
}
```

### 3. Respect Context Cancellation

```go
func TimeoutMiddleware(timeout time.Duration) middleware.PublishMiddleware {
    return func(next middleware.PublishHandler) middleware.PublishHandler {
        return func(ctx context.Context, subject string, msg *core.Message) error {
            ctx, cancel := context.WithTimeout(ctx, timeout)
            defer cancel()

            done := make(chan error, 1)
            go func() {
                done <- next(ctx, subject, msg)
            }()

            select {
            case err := <-done:
                return err
            case <-ctx.Done():
                return ctx.Err()
            }
        }
    }
}
```

## Testing

```bash
# Run tests
go test ./events/middleware/...

# Run with coverage
go test -cover ./events/middleware/...

# Run benchmarks
go test -bench=. ./events/middleware/...
```

## Usage in Other Packages

This package is used by:
- `events/transport` - Transport implementations with middleware
- `events` (main) - Re-exports for convenience
- Application code - Custom middleware pipelines
