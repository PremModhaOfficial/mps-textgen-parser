# Tenant Package

The `tenant` package provides multi-tenant connection management for the events messaging system, enabling isolated connections for each tenant with automatic lifecycle management.

## Overview

This package contains:

- **Manager**: Central manager for tenant connections
- **Connection**: Tenant-specific connection interface
- **Registry**: Tenant information storage
- **Idle Cleanup**: Automatic cleanup of idle connections
- **Health Monitoring**: Per-tenant health status

## Package Structure

```
tenant/
├── manager.go   # Tenant manager implementation
├── registry.go  # Tenant registry
└── README.md    # This file
```

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Tenant Manager                        │
│  ┌─────────────────────────────────────────────────┐    │
│  │              Connection Pool                     │    │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐         │    │
│  │  │Tenant A │  │Tenant B │  │Tenant C │  ...    │    │
│  │  │  Conn   │  │  Conn   │  │  Conn   │         │    │
│  │  └─────────┘  └─────────┘  └─────────┘         │    │
│  └─────────────────────────────────────────────────┘    │
│                                                          │
│  ┌──────────────┐  ┌───────────────┐  ┌────────────┐   │
│  │ Credential   │  │  Connection   │  │  Cleanup   │   │
│  │   Manager    │  │    Factory    │  │   Loop     │   │
│  └──────────────┘  └───────────────┘  └────────────┘   │
└─────────────────────────────────────────────────────────┘
```

## Manager

### ManagerConfig

```go
type ManagerConfig struct {
    Config          *config.Config        // Base configuration
    Factory         ConnectionFactory     // Connection factory
    CredManager     *auth.CredentialManager // Credential manager
    ErrorBufferSize int                   // Error channel buffer size
}
```

### Creating a Manager

```go
import (
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/tenant"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/config"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/transport/nats"
)

// Create manager with NATS transport
manager, err := tenant.NewManager(tenant.ManagerConfig{
    Config:      config.DefaultConfig(),
    Factory:     nats.NewFactory(),
    CredManager: credManager, // Optional
})
if err != nil {
    log.Fatal(err)
}
defer manager.Shutdown(context.Background())
```

### Manager Methods

| Method | Description |
|--------|-------------|
| `Connect(ctx, tenantID, creds)` | Connect a tenant |
| `ConnectWithCredentialManager(ctx, tenantID)` | Connect using stored credentials |
| `ConnectWithPayload(ctx, payload)` | Connect using registration payload |
| `Get(tenantID)` | Get existing connection |
| `Disconnect(ctx, tenantID)` | Disconnect a tenant |
| `Shutdown(ctx)` | Shutdown all connections |
| `Health()` | Get all tenant health statuses |
| `TenantHealth(tenantID)` | Get specific tenant health |
| `ActiveConnections()` | Count active connections |
| `TenantIDs()` | List connected tenant IDs |
| `Errors()` | Error notification channel |
| `OnError(callback)` | Set error callback |

## Connection Interface

```go
type Connection interface {
    core.Connection      // Connect, Close, IsConnected, State
    core.HealthChecker   // Health

    // TenantID returns the tenant identifier
    TenantID() string

    // Touch updates the last activity timestamp
    Touch()

    // IsIdle returns true if connection has been idle
    IsIdle(timeout time.Duration) bool

    // Publisher returns a publisher for this connection
    Publisher() core.Publisher

    // Subscriber returns a subscriber for this connection
    Subscriber() core.Subscriber
}
```

## Usage Examples

### Basic Connection

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/tenant"

// Connect with explicit credentials
creds := auth.JWT(jwtToken, seed)
conn, err := manager.Connect(ctx, "tenant-123", creds)
if err != nil {
    log.Fatal(err)
}

// Use the connection
pub := conn.Publisher()
sub := conn.Subscriber()
```

### Using Credential Manager

```go
// Pre-register credentials
manager.CredentialManager().Register("tenant-123", jwt, seed)

// Connect without passing credentials
conn, err := manager.ConnectWithCredentialManager(ctx, "tenant-123")
```

### Connection with Registration Payload

```go
payload := auth.RegistrationPayload{
    TenantID: "tenant-123",
    JWT:      jwtToken,
    Seed:     nkeySeed,
}

conn, err := manager.ConnectWithPayload(ctx, payload)
```

### Getting Existing Connections

```go
// Get existing connection (doesn't create new)
conn, ok := manager.Get("tenant-123")
if !ok {
    // Tenant not connected
    conn, err = manager.Connect(ctx, "tenant-123", creds)
}

// Connection activity is updated (Touch)
pub := conn.Publisher()
```

### Disconnecting Tenants

```go
// Disconnect specific tenant
err := manager.Disconnect(ctx, "tenant-123")
if errors.Is(err, core.ErrTenantNotFound) {
    // Tenant was not connected
}

// Shutdown all connections
err := manager.Shutdown(ctx)
```

## Health Monitoring

### Check All Tenants

```go
health := manager.Health()
for tenantID, status := range health {
    if !status.IsHealthy() {
        log.Printf("Tenant %s unhealthy: %s", tenantID, status.Message)
    }
}
```

### Check Specific Tenant

```go
status, ok := manager.TenantHealth("tenant-123")
if !ok {
    // Tenant not connected
}

if status.IsHealthy() {
    fmt.Println("Tenant is healthy")
} else {
    fmt.Printf("Unhealthy: %s (last error: %v)\n", status.Message, status.LastError)
}
```

### Health Status Structure

```go
type HealthStatus struct {
    State     ConnectionState    // Current state
    Healthy   bool               // Health indicator
    Message   string             // Status message
    LastError error              // Most recent error
    Details   map[string]any     // Additional details
}
```

## Error Handling

### Error Channel

```go
// Listen for tenant errors
go func() {
    for err := range manager.Errors() {
        log.Printf("Tenant %s error: %v", err.TenantID, err.Err)

        // Optionally reconnect
        if errors.Is(err.Err, core.ErrConnectionClosed) {
            go reconnectTenant(err.TenantID)
        }
    }
}()
```

### Error Callback

```go
// Set error callback
manager.OnError(func(err core.TenantError) {
    metrics.IncrementCounter("tenant_errors", err.TenantID)
    alerting.Notify(fmt.Sprintf("Tenant %s: %v", err.TenantID, err.Err))
})
```

## Idle Connection Cleanup

The manager automatically cleans up idle connections:

```go
// Configuration
cfg := config.DefaultConfig()
cfg.IdleTimeout = 30 * time.Minute    // Close after 30 min idle
cfg.CleanupInterval = 1 * time.Minute // Check every minute

manager, _ := tenant.NewManager(tenant.ManagerConfig{
    Config: cfg,
    // ...
})
```

**Cleanup Behavior:**

1. Background goroutine runs every `CleanupInterval`
2. Connections idle longer than `IdleTimeout` are closed
3. Activity is tracked via `Touch()` on each operation
4. Cleanup respects shutdown context

## Connection Factory

Interface for creating connections:

```go
type ConnectionFactory interface {
    Create(ctx context.Context, tenantID string,
           cfg *config.Config, creds *auth.Credentials) (Connection, error)
}
```

### Using NATS Factory

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/transport/nats"

factory := nats.NewFactory()
manager, _ := tenant.NewManager(tenant.ManagerConfig{
    Factory: factory,
    // ...
})
```

### Custom Factory

```go
type CustomFactory struct {
    // Custom configuration
}

func (f *CustomFactory) Create(ctx context.Context, tenantID string,
    cfg *config.Config, creds *auth.Credentials) (tenant.Connection, error) {
    // Create custom connection
    return &MyConnection{tenantID: tenantID}, nil
}
```

## Registry

Optional tenant information storage:

```go
type Registry struct {
    // Thread-safe tenant storage
}

type Info struct {
    ID        string
    Name      string
    Metadata  map[string]string
    CreatedAt time.Time
}
```

### Using Registry

```go
registry := tenant.NewRegistry()

// Register tenant info
registry.Register(&tenant.Info{
    ID:   "tenant-123",
    Name: "Acme Corp",
    Metadata: map[string]string{
        "plan": "enterprise",
        "region": "us-east-1",
    },
})

// Get tenant info
info, ok := registry.Get("tenant-123")

// List all tenants
tenants := registry.List()

// Deregister
registry.Deregister("tenant-123")
```

## Connection Lifecycle

```
┌──────────┐   Connect()   ┌────────────┐   Close()   ┌────────┐
│ Initial  │──────────────>│ Connected  │────────────>│ Closed │
└──────────┘               └────────────┘             └────────┘
                                 │
                                 │ Connection Lost
                                 ▼
                          ┌──────────────┐
                          │ Reconnecting │
                          └──────────────┘
                                 │
                    ┌────────────┴────────────┐
                    │                         │
              Reconnected                   Failed
                    │                         │
                    ▼                         ▼
             ┌────────────┐            ┌────────────┐
             │ Connected  │            │  Draining  │──> Closed
             └────────────┘            └────────────┘
```

## Best Practices

### 1. Initialize Once

```go
func main() {
    manager, err := tenant.NewManager(cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer manager.Shutdown(context.Background())

    // Pass manager to handlers
    server := NewServer(manager)
    server.Run()
}
```

### 2. Reuse Connections

```go
func handleRequest(tenantID string) {
    // Get existing or create new
    conn, ok := manager.Get(tenantID)
    if !ok {
        conn, _ = manager.Connect(ctx, tenantID, nil)
    }
    // Use connection...
}
```

### 3. Handle Reconnection

```go
manager.OnError(func(err core.TenantError) {
    if errors.Is(err.Err, core.ErrConnectionClosed) {
        // Schedule reconnection with backoff
        go func() {
            time.Sleep(time.Second)
            manager.Connect(ctx, err.TenantID, nil)
        }()
    }
})
```

### 4. Graceful Shutdown

```go
// Handle shutdown signals
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

<-sigChan

// Graceful shutdown with timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

if err := manager.Shutdown(ctx); err != nil {
    log.Printf("Shutdown error: %v", err)
}
```

### 5. Monitor Active Connections

```go
// Periodic monitoring
go func() {
    ticker := time.NewTicker(time.Minute)
    for range ticker.C {
        count := manager.ActiveConnections()
        metrics.SetGauge("active_tenant_connections", float64(count))

        health := manager.Health()
        unhealthy := 0
        for _, status := range health {
            if !status.IsHealthy() {
                unhealthy++
            }
        }
        metrics.SetGauge("unhealthy_tenants", float64(unhealthy))
    }
}()
```

## Thread Safety

All Manager and Registry methods are thread-safe:

- Connection map protected by `sync.RWMutex`
- Safe concurrent `Connect`, `Get`, `Disconnect`
- Thread-safe error channel and callbacks
- Safe for concurrent use from multiple goroutines

## Testing

```bash
# Run tests
go test ./events/tenant/...

# Run with coverage
go test -cover ./events/tenant/...

# Run with race detection
go test -race ./events/tenant/...
```

## Usage in Other Packages

This package is used by:
- `events` (main) - Client and TenantManager exports
- Application code - Multi-tenant connection management
