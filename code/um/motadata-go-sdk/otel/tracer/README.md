# Tracer Package

A production-ready distributed tracing package built on OpenTelemetry, providing comprehensive tracing capabilities for Go applications.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [API Reference](#api-reference)
- [Use Cases & Examples](#use-cases--examples)
- [Good Practices vs Bad Practices](#good-practices-vs-bad-practices)
- [Performance Pitfalls](#performance-pitfalls)
- [Benchmarks](#benchmarks)

## Features

- **OpenTelemetry Integration**: Full OTEL compatibility with gRPC and HTTP exporters
- **Multiple Span Kinds**: Internal, Server, Client, Producer, Consumer
- **Context Propagation**: W3C TraceContext, B3, Jaeger propagators
- **Configurable Sampling**: Always, Never, or Ratio-based sampling
- **Convenience Helpers**: `WithSpan`, `TimeSpan` for common patterns
- **Thread-Safe**: Safe for concurrent use across goroutines

## Installation

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer"
```

## Quick Start

```go
package main

import (
    "context"
    "log"

    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer"
)

func main() {
    // Initialize tracer
    cfg := tracer.DefaultConfig()
    cfg.Enabled = true
    cfg.ServiceName = "my-service"
    cfg.OTELEndpoint = "localhost:4317"

    t, err := tracer.Init(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer tracer.Shutdown(context.Background())

    // Create a span
    ctx, span := tracer.Start(context.Background(), "my-operation")
    defer span.End()

    // Your code here
    doWork(ctx)
}

func doWork(ctx context.Context) {
    ctx, span := tracer.Start(ctx, "do-work")
    defer span.End()

    span.SetAttributes(tracer.StringAttr("key", "value"))
    // ... work
}
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `OTEL_EXPORTER_OTLP_ENDPOINT` | OTLP collector endpoint | `localhost:4317` |
| `OTEL_EXPORTER_OTLP_PROTOCOL` | Protocol (grpc/http) | `grpc` |
| `OTEL_EXPORTER_OTLP_INSECURE` | Disable TLS | `false` |
| `OTEL_SERVICE_NAME` | Service name | `app` |
| `OTEL_SERVICE_VERSION` | Service version | `0.0.0` |
| `OTEL_ENVIRONMENT` | Deployment environment | `development` |
| `OTEL_TRACES_SAMPLER_ARG` | Sampling ratio (0.0-1.0) | `1.0` |
| `OTEL_PROPAGATORS` | Propagators (comma-separated) | `tracecontext,baggage` |

### Programmatic Configuration

```go
cfg := tracer.Config{
    Enabled:        true,
    ServiceName:    "my-service",
    ServiceVersion: "1.0.0",
    Environment:    "production",
    OTELEndpoint:   "otel-collector:4317",
    OTELProtocol:   "grpc",
    OTELInsecure:   false,
    SamplingRatio:  0.1, // Sample 10% of traces
    Propagators:    []string{"tracecontext", "baggage", "b3"},
    MaxQueueSize:   2048,
    MaxExportBatch: 512,
    BatchTimeout:   5 * time.Second,
}
```

### Supported Propagators

| Propagator | Description |
|------------|-------------|
| `tracecontext` | W3C Trace Context (default) |
| `baggage` | W3C Baggage (default) |
| `b3` | Zipkin B3 single header |
| `b3multi` | Zipkin B3 multi header |
| `jaeger` | Jaeger propagation format |

## API Reference

### Span Creation

```go
// Basic span (internal kind)
ctx, span := tracer.Start(ctx, "operation-name")
defer span.End()

// Server span (incoming request)
ctx, span := tracer.StartServer(ctx, "HTTP GET /api/users")
defer span.End()

// Client span (outgoing request)
ctx, span := tracer.StartClient(ctx, "HTTP POST /external-api")
defer span.End()

// Producer span (message publishing)
ctx, span := tracer.StartProducer(ctx, "kafka.publish")
defer span.End()

// Consumer span (message consuming)
ctx, span := tracer.StartConsumer(ctx, "kafka.consume")
defer span.End()
```

### Span Options

```go
ctx, span := tracer.Start(ctx, "operation",
    tracer.WithSpanKind(tracer.SpanKindClient),
    tracer.WithAttributes(
        tracer.StringAttr("http.method", "GET"),
        tracer.IntAttr("http.status_code", 200),
    ),
    tracer.WithLinks(link1, link2),
)
```

### Span Operations

```go
// Set attributes
span.SetAttributes(
    tracer.StringAttr("user.id", userID),
    tracer.IntAttr("items.count", len(items)),
)

// Add events (milestones within span)
span.AddEvent("cache-miss")
span.AddEvent("retry-attempt")

// Set status
span.SetOK()                           // Success
span.SetError(err)                     // Error with message and event
span.SetStatus(codes.Error, "failed")  // Manual status

// Record error as event (without changing status)
span.RecordError(err)

// Rename span
span.SetName("updated-operation-name")
```

### Context Utilities

```go
// Extract span from context
span := tracer.SpanFromContext(ctx)

// Create context with span
ctx = tracer.ContextWithSpan(ctx, span)

// Extract trace/span IDs for logging
traceID := tracer.TraceIDFromContext(ctx)
spanID := tracer.SpanIDFromContext(ctx)
```

### Convenience Helpers

```go
// WithSpan: automatic error handling
err := tracer.WithSpan(ctx, "operation", func(ctx context.Context) error {
    return doWork(ctx)
})

// TimeSpan: measure duration
duration := tracer.TimeSpan(ctx, "timed-operation", func(ctx context.Context) {
    doWork(ctx)
})
```

### Attribute Constructors

```go
tracer.StringAttr("key", "value")
tracer.IntAttr("key", 42)
tracer.Int64Attr("key", 9223372036854775807)
tracer.Float64Attr("key", 3.14159)
tracer.BoolAttr("key", true)
tracer.StringsAttr("key", []string{"a", "b"})
tracer.IntsAttr("key", []int{1, 2, 3})
tracer.Int64sAttr("key", []int64{1, 2})
tracer.Float64sAttr("key", []float64{1.1, 2.2})
tracer.BoolsAttr("key", []bool{true, false})
```

## Use Cases & Examples

### 1. HTTP Server Tracing

```go
func HTTPMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract trace context from incoming headers
        ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

        // Create server span
        ctx, span := tracer.StartServer(ctx, fmt.Sprintf("%s %s", r.Method, r.URL.Path))
        defer span.End()

        // Add HTTP attributes
        span.SetAttributes(
            tracer.StringAttr("http.method", r.Method),
            tracer.StringAttr("http.url", r.URL.String()),
            tracer.StringAttr("http.host", r.Host),
            tracer.StringAttr("http.user_agent", r.UserAgent()),
        )

        // Wrap response writer to capture status
        wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}

        next.ServeHTTP(wrapped, r.WithContext(ctx))

        // Record response
        span.SetAttributes(tracer.IntAttr("http.status_code", wrapped.statusCode))
        if wrapped.statusCode >= 400 {
            span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", wrapped.statusCode))
        } else {
            span.SetOK()
        }
    })
}
```

### 2. Database Query Tracing

```go
func (repo *UserRepository) FindByID(ctx context.Context, id string) (*User, error) {
    ctx, span := tracer.StartClient(ctx, "db.query")
    defer span.End()

    span.SetAttributes(
        tracer.StringAttr("db.system", "postgresql"),
        tracer.StringAttr("db.operation", "SELECT"),
        tracer.StringAttr("db.sql.table", "users"),
        tracer.StringAttr("db.statement", "SELECT * FROM users WHERE id = $1"),
    )

    user, err := repo.db.QueryRow(ctx, "SELECT * FROM users WHERE id = $1", id).Scan(...)
    if err != nil {
        span.SetError(err)
        return nil, err
    }

    span.SetOK()
    return user, nil
}
```

### 3. External HTTP Client Tracing

```go
func (c *APIClient) GetUser(ctx context.Context, userID string) (*User, error) {
    ctx, span := tracer.StartClient(ctx, "http.client.request")
    defer span.End()

    url := fmt.Sprintf("%s/users/%s", c.baseURL, userID)
    span.SetAttributes(
        tracer.StringAttr("http.method", "GET"),
        tracer.StringAttr("http.url", url),
        tracer.StringAttr("peer.service", "user-service"),
    )

    req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

    // Inject trace context into outgoing headers
    otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

    resp, err := c.httpClient.Do(req)
    if err != nil {
        span.SetError(err)
        return nil, err
    }
    defer resp.Body.Close()

    span.SetAttributes(tracer.IntAttr("http.status_code", resp.StatusCode))

    if resp.StatusCode >= 400 {
        span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", resp.StatusCode))
        return nil, fmt.Errorf("request failed: %d", resp.StatusCode)
    }

    span.SetOK()
    // ... parse response
}
```

### 4. Message Queue Producer/Consumer

```go
// Producer
func (p *Producer) Publish(ctx context.Context, topic string, message []byte) error {
    ctx, span := tracer.StartProducer(ctx, "kafka.publish")
    defer span.End()

    span.SetAttributes(
        tracer.StringAttr("messaging.system", "kafka"),
        tracer.StringAttr("messaging.destination", topic),
        tracer.StringAttr("messaging.destination_kind", "topic"),
        tracer.IntAttr("messaging.message.payload_size_bytes", len(message)),
    )

    // Inject trace context into message headers
    headers := make(map[string]string)
    otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(headers))

    err := p.kafka.Publish(topic, message, headers)
    if err != nil {
        span.SetError(err)
        return err
    }

    span.SetOK()
    return nil
}

// Consumer
func (c *Consumer) Handle(msg *kafka.Message) error {
    // Extract trace context from message headers
    ctx := otel.GetTextMapPropagator().Extract(
        context.Background(),
        propagation.MapCarrier(msg.Headers),
    )

    ctx, span := tracer.StartConsumer(ctx, "kafka.consume")
    defer span.End()

    span.SetAttributes(
        tracer.StringAttr("messaging.system", "kafka"),
        tracer.StringAttr("messaging.destination", msg.Topic),
        tracer.StringAttr("messaging.message_id", msg.ID),
    )

    err := c.processMessage(ctx, msg)
    if err != nil {
        span.SetError(err)
        return err
    }

    span.SetOK()
    return nil
}
```

### 5. Background Job Tracing

```go
func (w *Worker) ProcessJob(ctx context.Context, job Job) error {
    ctx, span := tracer.Start(ctx, fmt.Sprintf("job.%s", job.Type))
    defer span.End()

    span.SetAttributes(
        tracer.StringAttr("job.id", job.ID),
        tracer.StringAttr("job.type", job.Type),
        tracer.IntAttr("job.attempt", job.Attempt),
    )

    // Add event for job start
    span.AddEvent("job.started")

    result, err := w.execute(ctx, job)

    if err != nil {
        span.AddEvent("job.failed", trace.WithAttributes(
            attribute.String("error.message", err.Error()),
        ))
        span.SetError(err)
        return err
    }

    span.AddEvent("job.completed", trace.WithAttributes(
        attribute.Int("result.items_processed", result.ItemsProcessed),
    ))
    span.SetOK()
    return nil
}
```

### 6. Multi-Service Request Tracing

```go
// Service A - Orchestrator
func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderRequest) (*Order, error) {
    ctx, span := tracer.Start(ctx, "order.create")
    defer span.End()

    // Call User Service
    user, err := s.userClient.GetUser(ctx, req.UserID)
    if err != nil {
        span.SetError(err)
        return nil, err
    }
    span.AddEvent("user.validated")

    // Call Inventory Service
    err = s.inventoryClient.ReserveItems(ctx, req.Items)
    if err != nil {
        span.SetError(err)
        return nil, err
    }
    span.AddEvent("inventory.reserved")

    // Call Payment Service
    payment, err := s.paymentClient.ProcessPayment(ctx, req.Payment)
    if err != nil {
        span.SetError(err)
        // Rollback inventory
        s.inventoryClient.ReleaseItems(ctx, req.Items)
        return nil, err
    }
    span.AddEvent("payment.processed")

    // Create order
    order, err := s.repo.Create(ctx, user, req.Items, payment)
    if err != nil {
        span.SetError(err)
        return nil, err
    }

    span.SetAttributes(tracer.StringAttr("order.id", order.ID))
    span.SetOK()
    return order, nil
}
```

### 7. Retry with Tracing

```go
func (c *Client) RequestWithRetry(ctx context.Context, req *Request) (*Response, error) {
    ctx, span := tracer.Start(ctx, "request.with_retry")
    defer span.End()

    var lastErr error
    for attempt := 1; attempt <= 3; attempt++ {
        // Create child span for each attempt
        attemptCtx, attemptSpan := tracer.StartClient(ctx, "request.attempt")
        attemptSpan.SetAttributes(tracer.IntAttr("retry.attempt", attempt))

        resp, err := c.doRequest(attemptCtx, req)
        if err == nil {
            attemptSpan.SetOK()
            attemptSpan.End()
            span.SetOK()
            return resp, nil
        }

        attemptSpan.RecordError(err)
        attemptSpan.SetStatus(codes.Error, "attempt failed")
        attemptSpan.End()

        lastErr = err
        span.AddEvent("retry.scheduled", trace.WithAttributes(
            attribute.Int("attempt", attempt),
            attribute.String("error", err.Error()),
        ))

        time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
    }

    span.SetError(lastErr)
    return nil, lastErr
}
```

### 8. Correlation with Logs

```go
func ProcessRequest(ctx context.Context, req *Request) error {
    ctx, span := tracer.Start(ctx, "process.request")
    defer span.End()

    // Extract trace context for structured logging
    traceID := tracer.TraceIDFromContext(ctx)
    spanID := tracer.SpanIDFromContext(ctx)

    // Use with your logger
    logger.Info(ctx, "Processing request",
        logger.String("trace_id", traceID),
        logger.String("span_id", spanID),
        logger.String("request_id", req.ID),
    )

    // Process...

    return nil
}
```

## Good Practices vs Bad Practices

### Span Naming

```go
// GOOD: Descriptive, consistent naming
ctx, span := tracer.StartServer(ctx, "HTTP GET /api/users/{id}")
ctx, span := tracer.StartClient(ctx, "postgresql.query")
ctx, span := tracer.Start(ctx, "order.validate")

// BAD: Generic or dynamic names causing high cardinality
ctx, span := tracer.Start(ctx, "process")              // Too generic
ctx, span := tracer.Start(ctx, fmt.Sprintf("/users/%s", userID))  // Dynamic value in name
ctx, span := tracer.Start(ctx, req.URL.String())       // Full URL as name
```

### Context Propagation

```go
// GOOD: Always pass context through the call chain
func HandleRequest(ctx context.Context) error {
    ctx, span := tracer.Start(ctx, "handle.request")
    defer span.End()
    return processData(ctx)  // Pass ctx
}

func processData(ctx context.Context) error {
    ctx, span := tracer.Start(ctx, "process.data")
    defer span.End()
    return saveToDB(ctx)  // Pass ctx
}

// BAD: Breaking context chain
func HandleRequest(ctx context.Context) error {
    ctx, span := tracer.Start(ctx, "handle.request")
    defer span.End()
    return processData(context.Background())  // Lost trace context!
}
```

### Error Handling

```go
// GOOD: Record errors with context
if err != nil {
    span.SetError(err)  // Sets status AND records error event
    return err
}
span.SetOK()

// GOOD: Record error event without changing status
span.RecordError(tempErr)  // Log error but continue
// ... retry logic
span.SetOK()  // Final success

// BAD: Ignoring errors in traces
if err != nil {
    return err  // Span ends without error info
}

// BAD: Setting error status without recording
span.SetStatus(codes.Error, "failed")  // Missing error details
```

### Attribute Usage

```go
// GOOD: Semantic, bounded attributes
span.SetAttributes(
    tracer.StringAttr("http.method", "GET"),
    tracer.IntAttr("http.status_code", 200),
    tracer.StringAttr("user.id", userID),
    tracer.IntAttr("items.count", len(items)),
)

// BAD: High-cardinality or sensitive data
span.SetAttributes(
    tracer.StringAttr("request.body", string(requestBody)),  // PII/large data
    tracer.StringAttr("sql.query", fullQuery),               // May contain data
    tracer.Int64Attr("timestamp.nanos", time.Now().UnixNano()), // High cardinality
)
```

### Span Lifecycle

```go
// GOOD: Defer End() immediately after creation
ctx, span := tracer.Start(ctx, "operation")
defer span.End()

// GOOD: End child spans before parent
ctx, parent := tracer.Start(ctx, "parent")
defer parent.End()

ctx, child := tracer.Start(ctx, "child")
defer child.End()  // Ends before parent due to defer stack

// BAD: Forgetting to end spans
func process(ctx context.Context) {
    ctx, span := tracer.Start(ctx, "process")
    // Missing span.End() - LEAK!
    doWork(ctx)
}

// BAD: Conditional end
ctx, span := tracer.Start(ctx, "operation")
if condition {
    span.End()  // Only ends sometimes!
}
```

### Span Granularity

```go
// GOOD: Meaningful operation boundaries
func ProcessOrder(ctx context.Context, order *Order) error {
    ctx, span := tracer.Start(ctx, "order.process")
    defer span.End()

    if err := validateOrder(ctx, order); err != nil {  // Child span inside
        return err
    }
    if err := chargePayment(ctx, order); err != nil {  // Child span inside
        return err
    }
    return saveOrder(ctx, order)  // Child span inside
}

// BAD: Too granular - span per line
func ProcessOrder(ctx context.Context, order *Order) error {
    ctx, span1 := tracer.Start(ctx, "get.order.id")
    id := order.ID
    span1.End()

    ctx, span2 := tracer.Start(ctx, "check.id.not.empty")
    if id == "" { ... }
    span2.End()
    // ... excessive spans
}

// BAD: Too coarse - single span for everything
func HandleEverything(ctx context.Context) {
    ctx, span := tracer.Start(ctx, "handle.everything")
    defer span.End()
    // 1000 lines of code, no child spans
}
```

### Using WithSpan Helper

```go
// GOOD: Use WithSpan for simple error-returning functions
err := tracer.WithSpan(ctx, "validate.input", func(ctx context.Context) error {
    return validator.Validate(input)
})

// GOOD: Manual control when you need the span
ctx, span := tracer.Start(ctx, "complex.operation")
defer span.End()
span.SetAttributes(...)  // Need span reference
span.AddEvent(...)

// BAD: Nesting WithSpan unnecessarily
err := tracer.WithSpan(ctx, "outer", func(ctx context.Context) error {
    return tracer.WithSpan(ctx, "inner", func(ctx context.Context) error {
        return tracer.WithSpan(ctx, "innermost", func(ctx context.Context) error {
            return doWork(ctx)  // Hard to read
        })
    })
})
```

## Performance Pitfalls

### 1. Span Creation in Hot Paths

```go
// PROBLEM: Creating spans for every item in a loop
for _, item := range items {  // 10,000 items
    ctx, span := tracer.Start(ctx, "process.item")
    process(item)
    span.End()  // 10,000 spans!
}

// SOLUTION: Single span for batch, events for important items
ctx, span := tracer.Start(ctx, "process.items.batch")
defer span.End()
span.SetAttributes(tracer.IntAttr("items.count", len(items)))

for i, item := range items {
    if i%1000 == 0 {
        span.AddEvent("batch.progress", trace.WithAttributes(
            attribute.Int("processed", i),
        ))
    }
    process(item)
}
```

### 2. Excessive Attributes

```go
// PROBLEM: Adding too many attributes
for i := 0; i < 200; i++ {
    span.SetAttributes(tracer.IntAttr(fmt.Sprintf("item.%d", i), values[i]))
}
// SDK typically limits to 128 attributes, extras dropped

// SOLUTION: Use events or aggregate data
span.SetAttributes(
    tracer.IntAttr("items.count", len(values)),
    tracer.Float64Attr("items.avg", average(values)),
    tracer.Int64Attr("items.sum", sum(values)),
)
```

### 3. Large Attribute Values

```go
// PROBLEM: Storing large payloads in attributes
span.SetAttributes(
    tracer.StringAttr("request.body", string(largeJSON)),  // 1MB body!
    tracer.StringAttr("response.body", string(response)),
)

// SOLUTION: Store only necessary metadata
span.SetAttributes(
    tracer.IntAttr("request.body.size", len(largeJSON)),
    tracer.StringAttr("request.content_type", "application/json"),
)
```

### 4. Synchronous Span Export

```go
// PROBLEM: Using SimpleSpanProcessor in production
// This exports synchronously, blocking your code

// SOLUTION: Use BatchSpanProcessor (default in this package)
// Configured via MaxQueueSize, MaxExportBatch, BatchTimeout
```

### 5. Missing Sampling in High-Traffic Services

```go
// PROBLEM: Sampling 100% in production
cfg.SamplingRatio = 1.0  // Every request creates traces

// SOLUTION: Use appropriate sampling
cfg.SamplingRatio = 0.01  // 1% in high-traffic production
cfg.SamplingRatio = 0.1   // 10% in medium-traffic
cfg.SamplingRatio = 1.0   // 100% only for debugging/low-traffic
```

### 6. Not Checking IsRecording

```go
// PROBLEM: Expensive attribute computation even when not sampling
func expensiveAttributes() []tracer.Attribute {
    // Complex computation that takes 10ms
    data := computeExpensiveData()
    return []tracer.Attribute{
        tracer.StringAttr("computed.data", data),
    }
}

ctx, span := tracer.Start(ctx, "operation")
span.SetAttributes(expensiveAttributes()...)  // Always computed!

// SOLUTION: Check IsRecording before expensive operations
ctx, span := tracer.Start(ctx, "operation")
if span.IsRecording() {
    span.SetAttributes(expensiveAttributes()...)
}
```

### 7. Context Allocation in Loops

```go
// PROBLEM: Creating contexts in tight loops
for i := 0; i < 1000000; i++ {
    ctx, span := tracer.Start(ctx, "loop.iteration")
    span.End()  // Context allocations add up
}

// SOLUTION: Batch processing with single span
ctx, span := tracer.Start(ctx, "loop.batch")
defer span.End()
for i := 0; i < 1000000; i++ {
    // Process without creating span per iteration
}
```

### 8. Forgetting to Close Tracer

```go
// PROBLEM: Not shutting down tracer
func main() {
    tracer.Init(cfg)
    // ... application runs
    // Tracer never shut down, spans may be lost!
}

// SOLUTION: Always shutdown gracefully
func main() {
    t, _ := tracer.Init(cfg)
    defer func() {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        tracer.Shutdown(ctx)
    }()
    // ... application runs
}
```

### 9. Nil Context Panic

```go
// PROBLEM: Passing nil context
ctx, span := tracer.Start(nil, "operation")  // PANICS!

// SOLUTION: Always use valid context
ctx, span := tracer.Start(context.Background(), "operation")
// Or use context from request/caller
```

### 10. String Formatting in Span Names

```go
// PROBLEM: Dynamic span names with high cardinality
for _, userID := range userIDs {
    ctx, span := tracer.Start(ctx, fmt.Sprintf("user.%s.process", userID))
    // Creates unique span name per user - terrible for aggregation
}

// SOLUTION: Use static names with attributes
for _, userID := range userIDs {
    ctx, span := tracer.Start(ctx, "user.process")
    span.SetAttributes(tracer.StringAttr("user.id", userID))
}
```

## Benchmarks

Run benchmarks:

```bash
go test -bench=. -benchmem ./otel/tracer/
```

Typical results (noop tracer):

| Benchmark | ns/op | B/op | allocs/op |
|-----------|-------|------|-----------|
| SpanCreation | ~150 | 48 | 1 |
| SpanCreationWithAttributes | ~300 | 144 | 3 |
| SpanCreation10Attributes | ~600 | 400 | 5 |
| NestedSpanCreation | ~300 | 96 | 2 |
| SpanSetAttributes | ~50 | 32 | 1 |
| SpanAddEvent | ~40 | 24 | 1 |
| SpanFromContext | ~5 | 0 | 0 |
| TraceIDFromContext | ~20 | 0 | 0 |
| ConcurrentSpanCreation | ~200 | 48 | 1 |
| WithSpan | ~200 | 48 | 1 |
| AttributeCreation | ~15 | 24 | 1 |

## Thread Safety

The tracer package is fully thread-safe:

- Global tracer access is protected by RWMutex
- Spans can be safely used across goroutines
- Provider shutdown is idempotent

```go
// Safe concurrent usage
var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        ctx, span := tracer.Start(context.Background(), "concurrent.span")
        span.SetAttributes(tracer.IntAttr("goroutine.id", id))
        span.End()
    }(i)
}
wg.Wait()
```

## Integration Examples

### With gRPC

```go
import (
    "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
    "google.golang.org/grpc"
)

// Server
server := grpc.NewServer(
    grpc.StatsHandler(otelgrpc.NewServerHandler()),
)

// Client
conn, _ := grpc.Dial(target,
    grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
)
```

### With net/http

```go
import (
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Server
handler := otelhttp.NewHandler(yourHandler, "my-server")
http.ListenAndServe(":8080", handler)

// Client
client := &http.Client{
    Transport: otelhttp.NewTransport(http.DefaultTransport),
}
```

## License

MIT License
