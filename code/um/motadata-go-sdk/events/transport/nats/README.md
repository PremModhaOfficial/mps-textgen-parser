# NATS Transport Package

The `nats` package provides the NATS transport implementation for the events messaging system, supporting both Core NATS and JetStream.

## Overview

This package contains:

- **Factory**: Connection factory for tenant management
- **Connection**: NATS connection implementation
- **Publisher**: Message publishing with Core NATS and JetStream
- **Subscriber**: Message subscription with queue groups

## Package Structure

```
transport/nats/
├── transport.go   # Factory implementation
├── connection.go  # NATS connection wrapper
├── publisher.go   # Publisher implementation
├── subscriber.go  # Subscriber implementation
└── README.md      # This file
```

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    NATS Transport                        │
│                                                          │
│  ┌────────────┐                                         │
│  │  Factory   │───> Creates connections per tenant      │
│  └────────────┘                                         │
│        │                                                 │
│        ▼                                                 │
│  ┌────────────────────────────────────────────────┐     │
│  │               NATSConnection                    │     │
│  │  ┌──────────────┐    ┌──────────────────┐      │     │
│  │  │ Core NATS    │    │    JetStream     │      │     │
│  │  │  Connection  │    │     Context      │      │     │
│  │  └──────────────┘    └──────────────────┘      │     │
│  │        │                      │                │     │
│  │        ▼                      ▼                │     │
│  │  ┌──────────┐          ┌──────────┐           │     │
│  │  │Publisher │          │Publisher │           │     │
│  │  │ (Core)   │          │  (JS)    │           │     │
│  │  └──────────┘          └──────────┘           │     │
│  │        │                      │                │     │
│  │        ▼                      ▼                │     │
│  │  ┌──────────┐          ┌──────────┐           │     │
│  │  │Subscriber│          │Subscriber│           │     │
│  │  │ (Core)   │          │  (JS)    │           │     │
│  │  └──────────┘          └──────────┘           │     │
│  └────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────┘
```

## Factory

Creates connections for the tenant manager:

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/transport/nats"

// Create factory
factory := nats.NewFactory()

// Use with tenant manager
manager, _ := tenant.NewManager(tenant.ManagerConfig{
    Config:  config.DefaultConfig(),
    Factory: factory,
})
```

### Factory Interface

```go
type ConnectionFactory interface {
    Create(ctx context.Context, tenantID string,
           cfg *config.Config, creds *auth.Credentials) (tenant.Connection, error)
}
```

## Connection

### Creating a Connection

```go
import (
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/transport/nats"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/config"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"
)

cfg := config.DefaultConfig()
cfg.Servers = []string{"nats://localhost:4222"}
cfg.Stream.Enabled = true

creds := auth.JWT(jwtToken, seed)

// Create via factory
factory := nats.NewFactory()
conn, err := factory.Create(ctx, "tenant-123", cfg, creds)
if err != nil {
    log.Fatal(err)
}

// Connect
err = conn.Connect(ctx)
if err != nil {
    log.Fatal(err)
}
defer conn.Close(ctx)
```

### Connection Methods

| Method | Description |
|--------|-------------|
| `Connect(ctx)` | Establish connection |
| `Close(ctx)` | Close connection gracefully |
| `IsConnected()` | Check connection status |
| `State()` | Get connection state |
| `TenantID()` | Get tenant identifier |
| `Touch()` | Update activity timestamp |
| `IsIdle(timeout)` | Check if idle |
| `Publisher()` | Get publisher |
| `Subscriber()` | Get subscriber |
| `Health()` | Get health status |

### Connection Options

The connection respects configuration from `config.Config`:

```go
cfg := config.DefaultConfig()

// Server configuration
cfg.Servers = []string{
    "nats://server1:4222",
    "nats://server2:4222",
}
cfg.Name = "my-service"
cfg.ConnectTimeout = 10 * time.Second
cfg.DrainTimeout = 30 * time.Second

// TLS
cfg.TLS = &config.TLSConfig{
    Enabled:  true,
    CertFile: "/path/to/cert.pem",
    KeyFile:  "/path/to/key.pem",
    CAFile:   "/path/to/ca.pem",
}

// Reconnection
cfg.Reconnect = config.ReconnectConfig{
    MaxAttempts:     -1,  // Infinite
    InitialInterval: 100 * time.Millisecond,
    MaxInterval:     30 * time.Second,
    Multiplier:      2.0,
    Jitter:          0.1,
}

// JetStream
cfg.Stream = config.StreamConfig{
    Enabled: true,
    Domain:  "hub",
}
```

## Publisher

### Basic Publishing

```go
pub := conn.Publisher()

// Create message
msg := core.NewMessage([]byte(`{"order_id": "123"}`)).
    WithHeader("Content-Type", "application/json")

// Publish (Core NATS)
err := pub.Publish(ctx, "orders.created", msg)
```

### Async Publishing

```go
// Publish asynchronously
future := pub.PublishAsync(ctx, "orders.created", msg)

// Wait for acknowledgment
ack, err := future.Wait(ctx)
if err != nil {
    log.Printf("Publish failed: %v", err)
} else {
    log.Printf("Published to stream %s, seq %d", ack.Stream, ack.Sequence)
}

// Or non-blocking check
select {
case ack := <-future.Ok():
    log.Printf("Success: seq %d", ack.Sequence)
case err := <-future.Err():
    log.Printf("Failed: %v", err)
}
```

### Request-Reply

```go
// Send request and wait for reply
request := core.NewMessage([]byte(`{"user_id": "123"}`))
response, err := pub.Request(ctx, "users.get", request)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Response: %s\n", string(response.Data))
```

### JetStream Publishing

When JetStream is enabled, publish uses JetStream for persistence:

```go
cfg.Stream.Enabled = true

// Messages are persisted to JetStream
err := pub.Publish(ctx, "orders.created", msg)

// Get acknowledgment with stream info
future := pub.PublishAsync(ctx, "orders.created", msg)
ack, _ := future.Wait(ctx)
fmt.Printf("Stream: %s, Sequence: %d\n", ack.Stream, ack.Sequence)
```

### Message Deduplication

```go
// Set message ID for deduplication
msg := core.NewMessage(data).
    WithID("unique-msg-id-123")

// NATS will deduplicate within the window
err := pub.Publish(ctx, subject, msg)
```

## Subscriber

### Basic Subscription

```go
sub := conn.Subscriber()

// Subscribe to a subject
subscription, err := sub.Subscribe(ctx, "orders.*", func(ctx context.Context, msg *core.Message) error {
    fmt.Printf("Received order event: %s\n", msg.Subject)
    fmt.Printf("Data: %s\n", string(msg.Data))
    return nil
})
if err != nil {
    log.Fatal(err)
}
defer subscription.Unsubscribe()
```

### Queue Subscription

Load balance across multiple subscribers:

```go
// Multiple instances share the queue
subscription, err := sub.QueueSubscribe(ctx, "orders.created", "order-processors", handler)

// Messages are distributed across queue members
// Only one subscriber in the queue receives each message
```

### Wildcard Subscriptions

```go
// Single token wildcard (*)
sub.Subscribe(ctx, "orders.*", handler)
// Matches: orders.created, orders.updated
// Not: orders.items.added

// Multi-token wildcard (>)
sub.Subscribe(ctx, "orders.>", handler)
// Matches: orders.created, orders.items.added, orders.a.b.c
```

### Subscription Management

```go
subscription, _ := sub.Subscribe(ctx, subject, handler)

// Check if active
if subscription.IsValid() {
    fmt.Println("Subscription is active")
}

// Get subject
fmt.Println("Subscribed to:", subscription.Subject())

// Unsubscribe immediately
subscription.Unsubscribe()

// Or drain (process pending, then unsubscribe)
subscription.Drain()
```

## Authentication

### JWT Authentication

```go
creds := auth.JWT(jwtToken, nkeySeed)
conn, _ := factory.Create(ctx, tenantID, cfg, creds)
```

### User/Password

```go
creds := auth.UserPass("username", "password")
conn, _ := factory.Create(ctx, tenantID, cfg, creds)
```

### Token

```go
creds := auth.Token("secret-token")
conn, _ := factory.Create(ctx, tenantID, cfg, creds)
```

### NKey

```go
creds := auth.NKey("/path/to/user.nkey")
// or
creds := auth.NKeySeed("SUAM...")
conn, _ := factory.Create(ctx, tenantID, cfg, creds)
```

### Credentials File

```go
creds := auth.CredentialsFileAuth("/path/to/user.creds")
conn, _ := factory.Create(ctx, tenantID, cfg, creds)
```

## TLS Configuration

```go
cfg.TLS = &config.TLSConfig{
    Enabled:    true,
    CertFile:   "/etc/ssl/client.crt",
    KeyFile:    "/etc/ssl/client.key",
    CAFile:     "/etc/ssl/ca.crt",
    SkipVerify: false, // Don't skip in production!
}
```

## Reconnection Handling

The NATS connection automatically handles reconnection:

```go
cfg.Reconnect = config.ReconnectConfig{
    MaxAttempts:     -1,                    // Infinite retries
    InitialInterval: 100 * time.Millisecond,
    MaxInterval:     30 * time.Second,
    Multiplier:      2.0,
    Jitter:          0.1,
}

// Connection state changes can be monitored
// via Health() or parent TenantManager error channel
```

## Health Monitoring

```go
health := conn.Health()

fmt.Printf("State: %s\n", health.State)
fmt.Printf("Healthy: %v\n", health.Healthy)
fmt.Printf("Message: %s\n", health.Message)

if health.LastError != nil {
    fmt.Printf("Last Error: %v\n", health.LastError)
}

if health.Details != nil {
    fmt.Printf("RTT: %v\n", health.Details["rtt"])
}
```

## Error Handling

```go
// Connection errors
err := conn.Connect(ctx)
if err != nil {
    switch {
    case errors.Is(err, core.ErrConnectionTimeout):
        // Server unreachable
    case errors.Is(err, core.ErrInvalidCredential):
        // Auth failed
    default:
        // Other error
    }
}

// Publish errors
err = pub.Publish(ctx, subject, msg)
if err != nil {
    switch {
    case errors.Is(err, core.ErrNotConnected):
        // Connection lost
    case errors.Is(err, core.ErrPublishTimeout):
        // Ack timeout
    case errors.Is(err, core.ErrInvalidSubject):
        // Bad subject format
    }
}

// Subscribe errors
_, err = sub.Subscribe(ctx, subject, handler)
if err != nil {
    switch {
    case errors.Is(err, core.ErrNotConnected):
        // Connection lost
    case errors.Is(err, core.ErrInvalidSubject):
        // Bad subject format
    }
}
```

## JetStream Features

### Stream Configuration

JetStream streams are typically configured server-side, but the client interacts with them:

```go
cfg.Stream = config.StreamConfig{
    Enabled: true,
    Domain:  "hub",    // For leaf nodes
    Prefix:  "",       // Custom API prefix
}
```

### Publish Acknowledgment

```go
future := pub.PublishAsync(ctx, subject, msg)
ack, err := future.Wait(ctx)
if err == nil {
    fmt.Printf("Stream: %s\n", ack.Stream)
    fmt.Printf("Sequence: %d\n", ack.Sequence)
    fmt.Printf("Domain: %s\n", ack.Domain)
}
```

## Best Practices

### 1. Reuse Connections

```go
// Create connection once
conn, _ := factory.Create(ctx, tenantID, cfg, creds)
conn.Connect(ctx)

// Reuse for all operations
pub := conn.Publisher()
sub := conn.Subscriber()

// Don't create new connections per request
```

### 2. Use Queue Groups for Scaling

```go
// Multiple instances can process messages
sub.QueueSubscribe(ctx, "orders.created", "order-workers", handler)

// Messages are load-balanced automatically
```

### 3. Handle Reconnection Gracefully

```go
// Trust the built-in reconnection
cfg.Reconnect.MaxAttempts = -1 // Infinite

// Monitor health for alerting
go func() {
    ticker := time.NewTicker(30 * time.Second)
    for range ticker.C {
        health := conn.Health()
        if !health.IsHealthy() {
            alerting.Warn("Connection unhealthy", health.Message)
        }
    }
}()
```

### 4. Use Message IDs for Deduplication

```go
msg := core.NewMessage(data).
    WithID(fmt.Sprintf("%s-%d", orderID, time.Now().UnixNano()))

// Prevents duplicate processing
pub.Publish(ctx, subject, msg)
```

### 5. Drain Before Shutdown

```go
// Graceful shutdown
ctx, cancel := context.WithTimeout(context.Background(), cfg.DrainTimeout)
defer cancel()

// Close drains pending messages
conn.Close(ctx)
```

## Testing

```bash
# Run tests (requires NATS server)
go test ./events/transport/nats/...

# Run with coverage
go test -cover ./events/transport/nats/...

# Run integration tests
NATS_URL=nats://localhost:4222 go test -tags=integration ./events/transport/nats/...
```

## Dependencies

| Dependency | Version | Purpose |
|------------|---------|---------|
| `github.com/nats-io/nats.go` | v1.30+ | NATS client |

## Usage in Other Packages

This package is used by:
- `events/tenant` - Via ConnectionFactory interface
- `events` (main) - Default transport factory
