# Config Package

The `config` package provides configuration management for the events messaging system with support for programmatic, environment variable, and file-based configuration.

## Overview

This package contains:

- **Configuration Structs**: Config, TLSConfig, ReconnectConfig, StreamConfig
- **Default Configurations**: Sensible defaults for all settings
- **Environment Loading**: Load configuration from environment variables
- **Functional Options**: Type-safe configuration builder pattern
- **Validation**: Configuration validation

## Package Structure

```
config/
├── config.go     # Configuration structs and defaults
├── loader.go     # Environment and file loading
├── options.go    # Functional options pattern
└── README.md     # This file
```

## Configuration Structs

### Config

Main configuration structure:

```go
type Config struct {
    // Connection settings
    Servers        []string        // Server URLs
    Name           string          // Client name for identification
    ConnectTimeout time.Duration   // Connection timeout
    DrainTimeout   time.Duration   // Drain timeout for graceful shutdown

    // TLS configuration
    TLS *TLSConfig

    // Reconnection settings
    Reconnect ReconnectConfig

    // Stream/JetStream settings
    Stream StreamConfig

    // Idle connection cleanup
    IdleTimeout     time.Duration  // Close connections idle longer than this
    CleanupInterval time.Duration  // How often to check for idle connections
}
```

### TLSConfig

TLS/SSL configuration:

```go
type TLSConfig struct {
    Enabled    bool    // Enable TLS
    CertFile   string  // Client certificate file
    KeyFile    string  // Client key file
    CAFile     string  // CA certificate file
    SkipVerify bool    // Skip server certificate verification (not recommended)
}
```

### ReconnectConfig

Reconnection settings with exponential backoff:

```go
type ReconnectConfig struct {
    MaxAttempts     int            // -1 for infinite retries
    InitialInterval time.Duration  // Starting backoff interval
    MaxInterval     time.Duration  // Maximum backoff interval
    Multiplier      float64        // Backoff multiplier
    Jitter          float64        // Random jitter factor (0.0-1.0)
}
```

### StreamConfig

JetStream configuration:

```go
type StreamConfig struct {
    Enabled bool    // Enable JetStream
    Domain  string  // JetStream domain (for leaf nodes)
    Prefix  string  // API prefix
}
```

### PublishConfig

Publishing configuration:

```go
type PublishConfig struct {
    Retry               RetryConfig    // Retry settings
    AckTimeout          time.Duration  // How long to wait for acknowledgment
    EnableDeduplication bool           // Enable message deduplication
    DeduplicationWindow time.Duration  // Deduplication time window
}
```

### SubscribeConfig

Subscription configuration:

```go
type SubscribeConfig struct {
    QueueGroup    string         // Queue group for load balancing
    MaxConcurrent int            // Max concurrent message handlers
    AckWait       time.Duration  // How long before message is redelivered
    BatchSize     int            // Messages per batch
    BatchWait     time.Duration  // Max wait time for batch
}
```

### RetryConfig

Retry settings:

```go
type RetryConfig struct {
    MaxAttempts     int            // Maximum retry attempts
    InitialInterval time.Duration  // Initial backoff interval
    MaxInterval     time.Duration  // Maximum backoff interval
    Multiplier      float64        // Backoff multiplier
    Jitter          float64        // Random jitter factor
}
```

## Default Values

| Setting | Default | Description |
|---------|---------|-------------|
| `Servers` | `["nats://localhost:4222"]` | Server URLs |
| `ConnectTimeout` | `10s` | Connection timeout |
| `DrainTimeout` | `30s` | Graceful shutdown timeout |
| `IdleTimeout` | `30m` | Idle connection cleanup |
| `CleanupInterval` | `1m` | Cleanup check frequency |
| `Reconnect.MaxAttempts` | `-1` (infinite) | Reconnection attempts |
| `Reconnect.InitialInterval` | `100ms` | Initial backoff |
| `Reconnect.MaxInterval` | `30s` | Maximum backoff |
| `Reconnect.Multiplier` | `2.0` | Backoff multiplier |
| `Reconnect.Jitter` | `0.1` | Random jitter |
| `Stream.Enabled` | `true` | JetStream enabled |
| `Retry.MaxAttempts` | `3` | Publish retry attempts |
| `Retry.InitialInterval` | `100ms` | Initial retry interval |
| `Retry.MaxInterval` | `5s` | Maximum retry interval |
| `Publish.AckTimeout` | `5s` | Acknowledgment timeout |
| `Publish.EnableDeduplication` | `true` | Deduplication enabled |
| `Publish.DeduplicationWindow` | `2m` | Deduplication window |
| `Subscribe.MaxConcurrent` | `10` | Concurrent handlers |
| `Subscribe.AckWait` | `30s` | Redelivery timeout |
| `Subscribe.BatchSize` | `100` | Batch size |
| `Subscribe.BatchWait` | `100ms` | Batch wait time |

## Usage

### Programmatic Configuration

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/config"

// Start with defaults
cfg := config.DefaultConfig()

// Customize settings
cfg.Servers = []string{"nats://server1:4222", "nats://server2:4222"}
cfg.Name = "my-service"
cfg.ConnectTimeout = 15 * time.Second

// TLS configuration
cfg.TLS = &config.TLSConfig{
    Enabled:  true,
    CertFile: "/path/to/cert.pem",
    KeyFile:  "/path/to/key.pem",
    CAFile:   "/path/to/ca.pem",
}

// Custom reconnection
cfg.Reconnect = config.ReconnectConfig{
    MaxAttempts:     10,
    InitialInterval: 200 * time.Millisecond,
    MaxInterval:     60 * time.Second,
    Multiplier:      2.0,
    Jitter:          0.2,
}

// Validate configuration
if err := config.Validate(cfg); err != nil {
    log.Fatal(err)
}
```

### Environment Variables

```go
// Load from environment variables
cfg, err := config.LoadFromEnv()
if err != nil {
    log.Fatal(err)
}
```

**Supported Environment Variables:**

| Variable | Default | Description |
|----------|---------|-------------|
| `EVENTS_SERVERS` | `nats://localhost:4222` | Comma-separated server URLs |
| `EVENTS_NAME` | `app` | Client name |
| `EVENTS_CONNECT_TIMEOUT` | `10s` | Connection timeout |
| `EVENTS_DRAIN_TIMEOUT` | `30s` | Drain timeout |
| `EVENTS_IDLE_TIMEOUT` | `30m` | Idle connection timeout |
| `EVENTS_CLEANUP_INTERVAL` | `1m` | Cleanup interval |
| `EVENTS_TLS_ENABLED` | `false` | Enable TLS |
| `EVENTS_TLS_CERT` | | TLS certificate file |
| `EVENTS_TLS_KEY` | | TLS key file |
| `EVENTS_TLS_CA` | | TLS CA file |
| `EVENTS_TLS_SKIP_VERIFY` | `false` | Skip TLS verification |
| `EVENTS_RECONNECT_MAX_ATTEMPTS` | `-1` | Max reconnection attempts |
| `EVENTS_RECONNECT_INITIAL_INTERVAL` | `100ms` | Initial reconnect interval |
| `EVENTS_RECONNECT_MAX_INTERVAL` | `30s` | Max reconnect interval |
| `EVENTS_STREAM_ENABLED` | `true` | Enable JetStream |
| `EVENTS_STREAM_DOMAIN` | | JetStream domain |

### Functional Options

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/config"

// Create config using functional options
cfg := config.DefaultConfig()
config.Apply(cfg,
    config.WithServers("nats://server1:4222", "nats://server2:4222"),
    config.WithName("order-service"),
    config.WithConnectTimeout(15 * time.Second),
    config.WithDrainTimeout(60 * time.Second),
    config.WithReconnect(config.ReconnectConfig{
        MaxAttempts:     10,
        InitialInterval: 200 * time.Millisecond,
        MaxInterval:     60 * time.Second,
        Multiplier:      2.0,
        Jitter:          0.2,
    }),
    config.WithTLS(&config.TLSConfig{
        Enabled:  true,
        CertFile: "/path/to/cert.pem",
        KeyFile:  "/path/to/key.pem",
        CAFile:   "/path/to/ca.pem",
    }),
    config.WithStream(config.StreamConfig{
        Enabled: true,
        Domain:  "hub",
    }),
)
```

**Available Options:**

| Option | Description |
|--------|-------------|
| `WithServers(urls...)` | Set server URLs |
| `WithName(name)` | Set client name |
| `WithConnectTimeout(d)` | Set connection timeout |
| `WithDrainTimeout(d)` | Set drain timeout |
| `WithIdleTimeout(d)` | Set idle connection timeout |
| `WithCleanupInterval(d)` | Set cleanup interval |
| `WithTLS(cfg)` | Set TLS configuration |
| `WithReconnect(cfg)` | Set reconnection configuration |
| `WithStream(cfg)` | Set JetStream configuration |

## Configuration Cloning

```go
// Deep copy configuration
clonedCfg := cfg.Clone()

// Modify clone without affecting original
clonedCfg.Servers = []string{"nats://other-server:4222"}
```

## Validation

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/config"

cfg := config.DefaultConfig()
cfg.Servers = []string{} // Invalid - empty servers

err := config.Validate(cfg)
if err != nil {
    fmt.Println("Invalid config:", err)
    // Output: Invalid config: at least one server URL is required
}
```

**Validation Rules:**

- At least one server URL required
- Server URLs must be valid NATS URLs
- Timeouts must be positive
- Reconnect multiplier must be >= 1.0
- Jitter must be between 0.0 and 1.0
- TLS files must exist if TLS is enabled

## Default Configurations

### DefaultConfig

```go
cfg := config.DefaultConfig()
// Returns Config with sensible defaults for development
```

### DefaultReconnectConfig

```go
reconnect := config.DefaultReconnectConfig()
// Returns ReconnectConfig with exponential backoff defaults
```

### DefaultRetryConfig

```go
retry := config.DefaultRetryConfig()
// Returns RetryConfig for publish retries
```

### DefaultPublishConfig

```go
publish := config.DefaultPublishConfig()
// Returns PublishConfig with deduplication enabled
```

### DefaultSubscribeConfig

```go
subscribe := config.DefaultSubscribeConfig()
// Returns SubscribeConfig for queue subscriptions
```

## Example Configurations

### Development

```go
cfg := config.DefaultConfig()
cfg.Name = "my-service-dev"
// Uses localhost:4222, no TLS, infinite reconnects
```

### Production

```go
cfg := config.DefaultConfig()
cfg.Servers = []string{
    "nats://nats1.prod.internal:4222",
    "nats://nats2.prod.internal:4222",
    "nats://nats3.prod.internal:4222",
}
cfg.Name = "my-service-prod"
cfg.TLS = &config.TLSConfig{
    Enabled:  true,
    CertFile: "/etc/ssl/client.crt",
    KeyFile:  "/etc/ssl/client.key",
    CAFile:   "/etc/ssl/ca.crt",
}
cfg.Reconnect.MaxAttempts = 100  // Limit reconnects in prod
cfg.IdleTimeout = 10 * time.Minute  // Shorter idle timeout
```

### High-Throughput

```go
cfg := config.DefaultConfig()
cfg.Servers = []string{"nats://nats-cluster:4222"}

// Faster reconnection for high availability
cfg.Reconnect = config.ReconnectConfig{
    MaxAttempts:     -1,  // Infinite
    InitialInterval: 50 * time.Millisecond,
    MaxInterval:     5 * time.Second,
    Multiplier:      1.5,
    Jitter:          0.1,
}

// Shorter drain timeout
cfg.DrainTimeout = 10 * time.Second
```

## Testing

```bash
# Run tests
go test ./events/config/...

# Run with coverage
go test -cover ./events/config/...
```

## Usage in Other Packages

This package is used by:
- `events/tenant` - Tenant manager configuration
- `events/transport` - Transport configuration
- `events` (main) - Client configuration
