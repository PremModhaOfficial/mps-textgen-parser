# Transport Package

The `transport` package defines the transport abstraction layer for the events messaging system and contains transport implementations.

## Overview

This package provides:

- **Transport Abstraction**: Pluggable transport interface
- **NATS Transport**: Full-featured NATS implementation (in `nats/` subdirectory)
- **Future Transports**: Designed for additional transports (Kafka, Redis, etc.)

## Package Structure

```
transport/
├── transport.go   # Transport interface definitions
├── nats/          # NATS transport implementation
│   ├── transport.go
│   ├── connection.go
│   ├── publisher.go
│   ├── subscriber.go
│   └── README.md
└── README.md      # This file
```

## Transport Interface

The transport layer implements the `tenant.ConnectionFactory` interface:

```go
type ConnectionFactory interface {
    Create(ctx context.Context, tenantID string,
           cfg *config.Config, creds *auth.Credentials) (Connection, error)
}
```

## Available Transports

### NATS (Recommended)

Full-featured messaging with NATS:

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/transport/nats"

factory := nats.NewFactory()
```

**Features:**
- Core NATS pub/sub
- JetStream persistence
- Request-reply pattern
- Queue groups
- Automatic reconnection
- TLS support
- Multiple auth methods

See [nats/README.md](nats/README.md) for detailed documentation.

## Usage with Tenant Manager

```go
import (
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/tenant"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/transport/nats"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/config"
)

// Create tenant manager with NATS transport
manager, err := tenant.NewManager(tenant.ManagerConfig{
    Config:  config.DefaultConfig(),
    Factory: nats.NewFactory(),
})
```

## Adding New Transports

To add a new transport (e.g., Kafka):

1. Create a new subdirectory: `transport/kafka/`

2. Implement the `ConnectionFactory` interface:

```go
package kafka

type Factory struct {
    // Kafka-specific configuration
}

func NewFactory() *Factory {
    return &Factory{}
}

func (f *Factory) Create(ctx context.Context, tenantID string,
    cfg *config.Config, creds *auth.Credentials) (tenant.Connection, error) {
    // Create Kafka connection
    return &KafkaConnection{tenantID: tenantID}, nil
}
```

3. Implement the `tenant.Connection` interface:

```go
type KafkaConnection struct {
    tenantID string
    // Kafka client
}

func (c *KafkaConnection) Connect(ctx context.Context) error { ... }
func (c *KafkaConnection) Close(ctx context.Context) error { ... }
func (c *KafkaConnection) IsConnected() bool { ... }
func (c *KafkaConnection) State() core.ConnectionState { ... }
func (c *KafkaConnection) TenantID() string { ... }
func (c *KafkaConnection) Touch() { ... }
func (c *KafkaConnection) IsIdle(timeout time.Duration) bool { ... }
func (c *KafkaConnection) Publisher() core.Publisher { ... }
func (c *KafkaConnection) Subscriber() core.Subscriber { ... }
func (c *KafkaConnection) Health() core.HealthStatus { ... }
```

4. Implement `Publisher` and `Subscriber` interfaces

## Transport Comparison

| Feature | NATS | Kafka (Future) | Redis (Future) |
|---------|------|----------------|----------------|
| Latency | Very Low | Low | Very Low |
| Throughput | High | Very High | High |
| Persistence | JetStream | Built-in | Optional |
| Ordering | Per-subject | Per-partition | Per-stream |
| Clustering | Built-in | Built-in | Cluster mode |
| Best For | Microservices | Event Sourcing | Caching + Pub/Sub |

## Configuration

Each transport uses `config.Config` with transport-specific fields:

```go
cfg := config.DefaultConfig()

// Common settings
cfg.Servers = []string{"nats://localhost:4222"}
cfg.ConnectTimeout = 10 * time.Second
cfg.DrainTimeout = 30 * time.Second

// Transport-specific (NATS)
cfg.Stream.Enabled = true
cfg.Stream.Domain = "hub"
```

## Testing

```bash
# Test all transports
go test ./events/transport/...

# Test specific transport
go test ./events/transport/nats/...
```
