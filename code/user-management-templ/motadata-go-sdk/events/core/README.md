# Core Package

The `core` package provides fundamental interfaces, types, and context propagation utilities for the events messaging system.

## Overview

This package contains:

- **Interfaces**: Publisher, Subscriber, Subscription (NATS-typed)
- **Types**: TraceContext, Metadata, MessageHandler
- **Context Helpers**: Tenant ID, Message ID, Correlation ID, Trace Context propagation
- **Header Constants**: Standard messaging headers (W3C, B3, NATS)
- **Request-Reply**: Service, Endpoint, and Requester patterns (NATS micro)

## Package Structure

```
core/
├── messaging.go     # Publisher, Subscriber, Subscription interfaces and headers
├── context.go       # Context value helpers for propagation
├── request.go       # Request-reply pattern and NATS micro service
├── doc.go           # Package documentation
├── context_test.go  # Context tests
├── request_test.go  # Request-reply tests
└── README.md        # This file
```

## Core Interfaces

### Publisher

Publishes messages to NATS subjects:

```go
type Publisher interface {
    // Publish sends a message to the specified subject
    Publish(ctx context.Context, subject string, msg *nats.Msg) error

    // PublishAsync sends a message asynchronously (JetStream)
    PublishAsync(ctx context.Context, subject string, msg *nats.Msg) PubAckFuture

    // Request sends a message and waits for a reply
    Request(ctx context.Context, subject string, msg *nats.Msg) (*nats.Msg, error)

    // Close gracefully closes the publisher
    Close(ctx context.Context) error
}
```

### Subscriber

Subscribes to NATS subjects and receives messages:

```go
type Subscriber interface {
    // Subscribe creates a subscription to the specified subject
    Subscribe(ctx context.Context, subject string, handler MessageHandler) (Subscription, error)

    // QueueSubscribe creates a queue subscription for load balancing
    QueueSubscribe(ctx context.Context, subject, queue string, handler MessageHandler) (Subscription, error)

    // Close gracefully closes the subscriber
    Close(ctx context.Context) error
}
```

### Subscription

Represents an active subscription:

```go
type Subscription interface {
    // Subject returns the subscribed subject
    Subject() string

    // Unsubscribe removes the subscription
    Unsubscribe() error

    // Drain unsubscribes and waits for messages to be processed
    Drain() error

    // IsValid returns true if subscription is still active
    IsValid() bool
}
```

### MessageHandler

Function signature for handling incoming messages:

```go
type MessageHandler func(ctx context.Context, msg *nats.Msg) error
```

## Standard Headers

### W3C Trace Context

| Constant | Header Key | Description |
|----------|------------|-------------|
| `HeaderTraceParent` | `traceparent` | W3C trace parent header |
| `HeaderTraceState` | `tracestate` | W3C trace state |

### B3 Propagation Format

| Constant | Header Key | Description |
|----------|------------|-------------|
| `HeaderB3TraceID` | `X-B3-TraceId` | B3 trace identifier |
| `HeaderB3SpanID` | `X-B3-SpanId` | B3 span identifier |
| `HeaderB3ParentSpanID` | `X-B3-ParentSpanId` | B3 parent span |
| `HeaderB3Sampled` | `X-B3-Sampled` | B3 sampling flag |

### Application Headers

| Constant | Header Key | Description |
|----------|------------|-------------|
| `HeaderContentType` | `Content-Type` | Message content type |
| `HeaderMessageID` | `Nats-Msg-Id` | Message ID for deduplication |
| `HeaderCorrelationID` | `Correlation-ID` | Request correlation ID |
| `HeaderTenantID` | `X-Tenant-ID` | Multi-tenant identifier |
| `HeaderTraceID` | `X-Trace-ID` | Distributed trace ID |
| `HeaderSpanID` | `X-Span-ID` | Span identifier |
| `HeaderServiceError` | `Nats-Service-Error` | Service error message |
| `HeaderServiceErrorCode` | `Nats-Service-Error-Code` | Service error code |

## TraceContext

Holds distributed tracing information:

```go
type TraceContext struct {
    TraceID  string  // Distributed trace identifier
    SpanID   string  // Current span identifier
    ParentID string  // Parent span identifier
    Sampled  bool    // Whether trace is sampled
    State    string  // Additional trace state
}
```

## Metadata

Container for additional message context:

```go
type Metadata struct {
    Values map[string]string
}
```

## Context Helpers

### Adding Context Values

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"

ctx := context.Background()

// Add tenant ID
ctx = core.WithTenantID(ctx, "tenant-123")

// Add message ID
ctx = core.WithMessageID(ctx, "msg-456")

// Add correlation ID
ctx = core.WithCorrelationID(ctx, "corr-789")

// Add trace context
tc := &core.TraceContext{
    TraceID: "abc123",
    SpanID:  "def456",
    Sampled: true,
}
ctx = core.WithTraceContext(ctx, tc)

// Add metadata
meta := &core.Metadata{Values: map[string]string{"key": "value"}}
ctx = core.WithMetadata(ctx, meta)
```

### Extracting Context Values

```go
// Extract tenant ID
tenantID, ok := core.TenantIDFromContext(ctx)

// Extract message ID
messageID, ok := core.MessageIDFromContext(ctx)

// Extract correlation ID
correlationID, ok := core.CorrelationIDFromContext(ctx)

// Extract trace context
traceCtx, ok := core.TraceContextFromContext(ctx)

// Extract metadata
meta, ok := core.MetadataFromContext(ctx)
```

### Header Injection/Extraction

```go
// Extract context values into NATS message headers
headers := make(nats.Header)
core.ExtractHeaders(ctx, headers)
// Headers now contain X-Tenant-ID, X-Trace-ID, Correlation-ID, etc.

// Inject NATS headers into context
ctx = core.InjectContext(ctx, headers)
// Context now contains tenant ID, trace context, etc.

// Extract OTEL trace context
tc := core.ExtractOTELTraceContext(ctx)
```

## Request-Reply Pattern

### Service Configuration

For NATS micro services:

```go
type ServiceConfig struct {
    Name        string
    Version     string
    Description string
    QueueGroup  string
    Metadata    map[string]string
}

type EndpointConfig struct {
    Subject   string
    Handler   RequestHandler
    QueueGroup string
    Metadata  map[string]string
}
```

### Request Options

```go
// Configure request behavior
opt := core.WithRequestTimeout(5 * time.Second)
opt = core.WithRequestHeaders(headers)
opt = core.WithMaxReplies(1)
opt = core.WithExpectedRTT(100 * time.Millisecond)
```

### Endpoint Options

```go
// Configure endpoint behavior
opt := core.WithEndpointSubject("custom.subject")
opt = core.WithEndpointMetadata(map[string]string{"version": "v2"})
opt = core.WithEndpointQueueGroup("workers")
opt = core.WithEndpointQueueGroupDisabled()
```

### Response Options

```go
// Configure response
opt := core.WithResponseHeaders(headers)
opt = core.WithResponseHeader("X-Custom", "value")
```

### Control Verbs

```go
const (
    PingVerb  Verb = "PING"
    StatsVerb Verb = "STATS"
    InfoVerb  Verb = "INFO"
)

// Generate control subject
subject := core.ControlSubject(PingVerb, "my-service", "v1")
```

### Service Statistics

```go
// Record request on an endpoint
core.RecordRequest(endpoint, duration, dataLen)

// Get endpoint statistics
stats := endpoint.Stats()

// Reset statistics
endpoint.Reset()
```

## Step-by-Step Usage Guide

### Step 1: Import the Package

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
```

### Step 2: Add Context Values in Middleware or Handlers

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // Propagate tenant and trace info
    ctx = core.WithTenantID(ctx, r.Header.Get("X-Tenant-ID"))
    ctx = core.WithCorrelationID(ctx, r.Header.Get("X-Correlation-ID"))

    // Pass context to messaging operations
    publisher.Publish(ctx, "events.request", msg)
}
```

### Step 3: Extract Context in Message Handlers

```go
sub.Subscribe(ctx, "events.>", func(ctx context.Context, msg *nats.Msg) error {
    tenantID, _ := core.TenantIDFromContext(ctx)
    traceCtx, _ := core.TraceContextFromContext(ctx)

    log.Printf("Processing message for tenant %s, trace %s",
        tenantID, traceCtx.TraceID)

    return nil
})
```

### Step 4: Propagate Headers Across Service Boundaries

```go
// Service A: Extract from context to headers
headers := make(nats.Header)
core.ExtractHeaders(ctx, headers)
msg := &nats.Msg{Subject: "service-b.process", Data: data, Header: headers}
publisher.Publish(ctx, msg.Subject, msg)

// Service B: Inject from headers to context
sub.Subscribe(ctx, "service-b.process", func(ctx context.Context, msg *nats.Msg) error {
    ctx = core.InjectContext(ctx, msg.Header)
    // Now ctx has tenant ID, trace context, etc. from Service A
    return nil
})
```

## Error Types

The core package uses error types from the `utils` package. Common errors:

| Error | Description |
|-------|-------------|
| `ErrNotConnected` | Operation on disconnected client |
| `ErrPublishFailed` | Message publish failed |
| `ErrPublishTimeout` | Publish acknowledgment timeout |
| `ErrInvalidSubject` | Invalid subject format |
| `ErrInvalidMessage` | Invalid message format |

See [`../utils/README.md`](../utils/README.md) for the full error reference.

## PubAckFuture

For asynchronous publish operations (JetStream):

```go
type PubAckFuture interface {
    // Ok returns a channel that receives the ack on success
    Ok() <-chan *PubAck

    // Err returns a channel that receives error on failure
    Err() <-chan error

    // Wait blocks until the publish completes or context cancels
    Wait(ctx context.Context) (*PubAck, error)
}

type PubAck struct {
    Stream   string  // Stream name
    Sequence uint64  // Sequence number
    Domain   string  // JetStream domain
}
```

**Usage:**

```go
future := publisher.PublishAsync(ctx, subject, msg)

// Non-blocking check
select {
case ack := <-future.Ok():
    fmt.Printf("Published to stream %s, seq %d\n", ack.Stream, ack.Sequence)
case err := <-future.Err():
    fmt.Printf("Publish failed: %v\n", err)
}

// Or blocking wait
ack, err := future.Wait(ctx)
```

## Testing

```bash
# Run tests
go test ./events/core/...

# Run with coverage
go test -cover ./events/core/...

# Run benchmarks
go test -bench=. ./events/core/...
```

## Usage in Other Packages

This package is used by:
- `events` (main) - Re-exports context helpers and types
- `events/tenant` - Connection management
- `events/middleware` - Middleware pipeline context propagation
- `events/endpoint` - Endpoint management
- `events/microservice` - Service management
- Application code - Context propagation in handlers
