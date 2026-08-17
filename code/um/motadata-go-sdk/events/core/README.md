# Core Package

The `core` package provides fundamental interfaces, types, and error definitions for the events messaging system.

## Overview

This package contains:

- **Interfaces**: Connection, Publisher, Subscriber, Subscription
- **Types**: Message, Headers, TraceContext, HealthStatus
- **Context Helpers**: Tenant ID, Request ID, Trace Context propagation
- **Error Definitions**: Standard errors for messaging operations

## Package Structure

```
core/
├── interfaces.go     # Core interfaces (Connection, Publisher, Subscriber)
├── message.go        # Message and Headers types
├── context.go        # Context value helpers
├── errors.go         # Error definitions
├── message_test.go   # Message tests
├── context_test.go   # Context tests
├── errors_test.go    # Error tests
└── README.md         # This file
```

## Core Interfaces

### Connection

Represents a connection to a messaging system:

```go
type Connection interface {
    // Connect establishes the connection
    Connect(ctx context.Context) error

    // Close gracefully closes the connection
    Close(ctx context.Context) error

    // IsConnected returns true if currently connected
    IsConnected() bool

    // State returns the current connection state
    State() ConnectionState
}
```

### Publisher

Publishes messages to subjects/topics:

```go
type Publisher interface {
    // Publish sends a message to the specified subject
    Publish(ctx context.Context, subject string, msg *Message) error

    // PublishAsync sends a message asynchronously and returns a future
    PublishAsync(ctx context.Context, subject string, msg *Message) PubAckFuture

    // Request sends a message and waits for a reply (request-reply pattern)
    Request(ctx context.Context, subject string, msg *Message) (*Message, error)

    // Close gracefully closes the publisher
    Close(ctx context.Context) error
}
```

### Subscriber

Subscribes to subjects/topics and receives messages:

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

## Message Types

### Message

The primary message structure for publishing and receiving:

```go
type Message struct {
    Subject   string        // Destination subject/topic
    Data      []byte        // Message payload
    Headers   Headers       // Message metadata
    Reply     string        // Reply subject for request-reply patterns
    Timestamp time.Time     // When the message was created
    ID        string        // Unique message identifier (for deduplication)
}
```

**Usage:**

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"

// Create a new message
msg := core.NewMessage([]byte(`{"user": "john"}`))

// Fluent API for building messages
msg = msg.
    WithSubject("users.created").
    WithHeader("Content-Type", "application/json").
    WithID("msg-12345")

// Add multiple headers
msg.WithHeaders(core.Headers{
    "X-Tenant-ID": []string{"tenant-1"},
    "X-Trace-ID":  []string{"trace-abc"},
})
```

### Headers

HTTP-like headers for message metadata:

```go
type Headers map[string][]string

// Methods
func (h Headers) Get(key string) string       // Get first value
func (h Headers) Set(key, value string)       // Set value (replaces existing)
func (h Headers) Add(key, value string)       // Add value (supports multiple)
func (h Headers) Del(key string)              // Delete header
func (h Headers) Has(key string) bool         // Check if exists
func (h Headers) Values(key string) []string  // Get all values
func (h Headers) Clone() Headers              // Deep copy
```

**Standard Header Constants:**

| Constant | Header Key | Description |
|----------|------------|-------------|
| `HeaderContentType` | `Content-Type` | Message content type |
| `HeaderMessageID` | `Nats-Msg-Id` | Message ID for deduplication |
| `HeaderCorrelationID` | `Correlation-ID` | Request correlation ID |
| `HeaderTenantID` | `X-Tenant-ID` | Multi-tenant identifier |
| `HeaderTraceID` | `X-Trace-ID` | Distributed trace ID |
| `HeaderSpanID` | `X-Span-ID` | Span identifier |
| `HeaderTraceParent` | `traceparent` | W3C trace context |
| `HeaderTraceState` | `tracestate` | W3C trace state |

### TraceContext

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

## Connection States

```go
const (
    StateDisconnected ConnectionState = iota  // Not connected
    StateConnecting                           // Connection in progress
    StateConnected                            // Successfully connected
    StateReconnecting                         // Attempting reconnection
    StateDraining                             // Draining pending messages
    StateClosed                               // Connection closed
)
```

**State Transitions:**

```
[Start] --> Disconnected --> Connecting --> Connected
                                 |              |
                                 v              v
                            [Error]       Reconnecting
                                               |
                                               v
                                          Connected
                                               |
                                               v
                                          Draining --> Closed
```

## Health Status

```go
type HealthStatus struct {
    State     ConnectionState    // Current connection state
    Healthy   bool               // Overall health indicator
    Message   string             // Human-readable status message
    LastError error              // Most recent error
    Details   map[string]any     // Additional details
}

// Check if healthy
func (h HealthStatus) IsHealthy() bool {
    return h.Healthy && h.State == StateConnected
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
```

### Header Injection/Extraction

```go
// Extract context values into headers
headers := make(core.Headers)
core.ExtractHeaders(ctx, headers)
// Headers now contain X-Tenant-ID, X-Trace-ID, etc.

// Inject headers into context
ctx = core.InjectContext(ctx, headers)
// Context now contains tenant ID, trace context, etc.
```

## Error Definitions

### Standard Errors

| Error | Description |
|-------|-------------|
| `ErrNotConnected` | Operation on disconnected client |
| `ErrAlreadyConnected` | Duplicate connection attempt |
| `ErrConnectionClosed` | Connection was closed |
| `ErrConnectionTimeout` | Connection timed out |
| `ErrPublishFailed` | Message publish failed |
| `ErrPublishTimeout` | Publish acknowledgment timeout |
| `ErrInvalidSubject` | Invalid subject format |
| `ErrInvalidMessage` | Invalid message format |
| `ErrTenantNotFound` | Tenant not registered |
| `ErrShutdownInProgress` | System is shutting down |
| `ErrInvalidCredential` | Invalid authentication credentials |

### Error Handling

```go
import (
    "errors"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

err := publisher.Publish(ctx, subject, msg)
if err != nil {
    switch {
    case errors.Is(err, core.ErrNotConnected):
        // Reconnect and retry
    case errors.Is(err, core.ErrPublishTimeout):
        // Retry with backoff
    case errors.Is(err, core.ErrInvalidSubject):
        // Fix subject format
    default:
        // Handle unknown error
    }
}
```

### TenantError

Wraps errors with tenant context:

```go
type TenantError struct {
    TenantID string
    Err      error
}

func (e TenantError) Error() string {
    return fmt.Sprintf("tenant %s: %v", e.TenantID, e.Err)
}

func (e TenantError) Unwrap() error {
    return e.Err
}
```

### Custom Error Creation

```go
// Create wrapped error with operation context
err := core.NewError("connection", "connect", core.ErrConnectionTimeout)
// Error: "connection connect: connection timeout"
```

## PubAckFuture

For asynchronous publish operations:

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
    Stream   string  // Stream name (for JetStream)
    Sequence uint64  // Sequence number
    Domain   string  // JetStream domain
}
```

**Usage:**

```go
// Async publish
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

## Component Interface

For managed shutdown:

```go
type Component interface {
    // Name returns the component name for identification
    Name() string

    // Shutdown gracefully shuts down the component
    Shutdown(ctx context.Context) error
}
```

## MessageHandler

Function signature for message handlers:

```go
type MessageHandler func(ctx context.Context, msg *Message) error
```

**Usage:**

```go
handler := func(ctx context.Context, msg *Message) error {
    fmt.Printf("Received: %s\n", string(msg.Data))

    // Extract tenant from context
    tenantID, _ := core.TenantIDFromContext(ctx)
    fmt.Printf("Tenant: %s\n", tenantID)

    return nil
}

subscription, err := subscriber.Subscribe(ctx, "orders.*", handler)
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
- `events/tenant` - Connection management
- `events/transport` - Transport implementations
- `events/middleware` - Middleware pipeline
- `events/endpoint` - Endpoint management
- `events/microservice` - Service management
- `events` (main) - Re-exports for convenience
