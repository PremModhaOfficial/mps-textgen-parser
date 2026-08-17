# Endpoint Package

The `endpoint` package provides endpoint registration, discovery, and health checking for microservice architectures.

## Overview

This package enables services to:

- **Register Endpoints**: Publish API endpoint information to a central registry
- **Discover Endpoints**: Find other service endpoints for communication
- **Health Checking**: Automatic HTTP health checks with status propagation
- **TTL Management**: Automatic cleanup of expired endpoints
- **Event Publishing**: Lifecycle event notifications

## Package Structure

```
endpoint/
├── types.go        # EndpointInfo, Status, Protocol types
├── registry.go     # In-memory endpoint registry
├── manager.go      # Endpoint lifecycle manager
├── handlers.go     # NATS message handlers
├── errors.go       # Error definitions
├── doc.go          # Package documentation
├── example_test.go # Example tests
└── README.md       # This file
```

## Key Types

### Info

Comprehensive endpoint information:

```go
type Info struct {
    // Identification
    ID          string    // Unique endpoint identifier
    Name        string    // Human-readable name
    ServiceID   string    // Parent service identifier
    ServiceName string    // Parent service name

    // Network
    Host     string    // Hostname or IP
    Port     int       // Port number
    Path     string    // URL path
    Protocol Protocol  // HTTP, HTTPS, gRPC, etc.

    // HTTP specifics
    Methods []Method   // Supported HTTP methods

    // Health
    Status          Status  // Current status
    HealthCheckPath string  // Path for health checks
    HealthCheckPort int     // Port for health checks (optional)

    // Metadata
    Tags        []string          // Searchable tags
    Version     string            // API version
    Description string            // Endpoint description
    Metadata    map[string]string // Custom metadata

    // Multi-tenancy
    TenantID string  // Tenant identifier
    Region   string  // Deployment region
    Zone     string  // Availability zone

    // Timestamps
    RegisteredAt time.Time
    LastSeenAt   time.Time
    ExpiresAt    time.Time
}
```

### Status

```go
type Status int

const (
    StatusUnknown     Status = iota  // Status not determined
    StatusHealthy                    // Endpoint is healthy
    StatusUnhealthy                  // Endpoint is unhealthy
    StatusDegraded                   // Running with issues
    StatusMaintenance                // Under maintenance
    StatusOffline                    // Offline
)
```

### Protocol

```go
type Protocol string

const (
    ProtocolHTTP  Protocol = "http"
    ProtocolHTTPS Protocol = "https"
    ProtocolGRPC  Protocol = "grpc"
    ProtocolWS    Protocol = "ws"
    ProtocolWSS   Protocol = "wss"
    ProtocolTCP   Protocol = "tcp"
    ProtocolUDP   Protocol = "udp"
)
```

## Registry

Thread-safe in-memory endpoint storage:

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/endpoint"

// Create registry with defaults
registry := endpoint.NewRegistry()

// Or with custom config
registry := endpoint.NewRegistry(endpoint.RegistryConfig{
    CleanupInterval:   30 * time.Second,
    ExpirationEnabled: true,
})
```

### Registry Methods

| Method | Description |
|--------|-------------|
| `Register(ep)` | Register endpoint (no expiration) |
| `RegisterWithTTL(ep, ttl)` | Register with automatic expiration |
| `Deregister(id)` | Remove endpoint |
| `Get(id)` | Get endpoint by ID |
| `GetByService(serviceID)` | Get all endpoints for a service |
| `Query(query)` | Query endpoints with filters |
| `UpdateStatus(id, status)` | Update endpoint status |
| `Renew(id, ttl)` | Extend endpoint TTL |
| `List()` | List all endpoints |
| `Count()` | Count endpoints |
| `Close()` | Stop cleanup goroutine |

### Registration

```go
ep := &endpoint.Info{
    ID:              "user-api-v1",
    ServiceID:       "user-service",
    Host:            "localhost",
    Port:            8080,
    Path:            "/api/v1/users",
    Protocol:        endpoint.ProtocolHTTP,
    Methods:         []endpoint.Method{endpoint.MethodGET, endpoint.MethodPOST},
    Status:          endpoint.StatusHealthy,
    Tags:            []string{"api", "users", "v1"},
    HealthCheckPath: "/health",
}

// Register permanently
registry.Register(ep)

// Register with 5-minute TTL
registry.RegisterWithTTL(ep, 5*time.Minute)

// Renew TTL (heartbeat)
registry.Renew("user-api-v1", 5*time.Minute)
```

### Querying

```go
// Query with filters
results := registry.Query(&endpoint.Query{
    ServiceIDs:  []string{"user-service"},
    Tags:        []string{"api"},
    Protocols:   []endpoint.Protocol{endpoint.ProtocolHTTP},
    OnlyHealthy: true,
    TenantID:    "tenant-123",
    Region:      "us-east-1",
})

// Get by service
userEndpoints := registry.GetByService("user-service")

// Get by ID
ep, ok := registry.Get("user-api-v1")
```

### Callbacks

```go
registry.OnRegistered(func(ep *endpoint.Info) {
    log.Printf("Endpoint registered: %s", ep.ID)
})

registry.OnDeregistered(func(ep *endpoint.Info, reason string) {
    log.Printf("Endpoint deregistered: %s, reason: %s", ep.ID, reason)
})

registry.OnStatusChange(func(ep *endpoint.Info, old, new endpoint.Status) {
    log.Printf("Endpoint %s: %s -> %s", ep.ID, old, new)
})
```

## Manager

Handles endpoint lifecycle with health checking:

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/endpoint"

manager, err := endpoint.NewManager(endpoint.ManagerConfig{
    Registry:  registry,
    Publisher: publisher, // Optional: for event publishing
    HealthCheck: endpoint.HealthCheckConfig{
        Enabled:            true,
        Interval:           30 * time.Second,
        Timeout:            5 * time.Second,
        HealthyThreshold:   2,  // Consecutive successes to become healthy
        UnhealthyThreshold: 3,  // Consecutive failures to become unhealthy
    },
})
if err != nil {
    log.Fatal(err)
}
defer manager.Shutdown(ctx)
```

### Manager Methods

| Method | Description |
|--------|-------------|
| `RegisterEndpoint(ctx, ep)` | Register and start health checking |
| `DeregisterEndpoint(ctx, id)` | Deregister and stop health checking |
| `UpdateEndpoint(ctx, ep)` | Update endpoint information |
| `Discover(serviceID, opts...)` | Discover endpoints |
| `UpdateHealth(ctx, id, update)` | Manual health update |
| `Shutdown(ctx)` | Graceful shutdown |

### Discovery Options

```go
endpoints := manager.Discover("user-service",
    endpoint.WithTags("api", "v2"),
    endpoint.WithProtocol(endpoint.ProtocolHTTP),
    endpoint.WithVersion("2.0.0"),
    endpoint.OnlyHealthy(),
    endpoint.InRegion("us-east-1"),
)
```

## Handler

NATS message handlers for distributed operations:

```go
handler := endpoint.NewHandler(manager)

// Register all subscriptions
subscriptions, err := handler.RegisterSubscriptions(ctx, subscriber)

// Or use queue subscriptions for load balancing
subscriptions, err := handler.RegisterQueueSubscriptions(ctx, subscriber, "endpoint-workers")
```

### Message Subjects

| Subject | Description |
|---------|-------------|
| `endpoints.register` | Register endpoint request |
| `endpoints.deregister` | Deregister endpoint request |
| `endpoints.query` | Query endpoints request |
| `endpoints.renew` | Renew endpoint TTL |
| `endpoints.health` | Update health status |

### Published Events

| Subject | Description |
|---------|-------------|
| `endpoints.registered` | Endpoint was registered |
| `endpoints.deregistered` | Endpoint was deregistered |
| `endpoints.updated` | Endpoint was updated |
| `endpoints.health.changed` | Health status changed |

## Client

Remote endpoint operations over messaging:

```go
client := endpoint.NewClient(publisher)

// Register endpoint remotely
resp, err := client.Register(ctx, ep, 5*time.Minute)

// Discover endpoints remotely
resp, err := client.Discover(ctx, "user-service",
    endpoint.OnlyHealthy(),
    endpoint.WithTags("api"),
)

// Renew TTL (heartbeat)
resp, err := client.Renew(ctx, "user-api-v1", 5*time.Minute)

// Deregister
resp, err := client.Deregister(ctx, "user-api-v1")
```

## Health Checking

### Configuration

```go
healthConfig := endpoint.HealthCheckConfig{
    Enabled:            true,
    Interval:           30 * time.Second,  // Check frequency
    Timeout:            5 * time.Second,   // HTTP timeout
    HealthyThreshold:   2,                 // Successes to become healthy
    UnhealthyThreshold: 3,                 // Failures to become unhealthy
}
```

### Health Check Behavior

| HTTP Status | Result |
|-------------|--------|
| 2xx | Healthy |
| 503 | Maintenance |
| Other/Error | Unhealthy |

### Endpoint Health Path

```go
ep := &endpoint.Info{
    Host:            "localhost",
    Port:            8080,
    Protocol:        endpoint.ProtocolHTTP,
    HealthCheckPath: "/health",      // Health check path
    HealthCheckPort: 8081,           // Optional: different port
}

// Health check URL: http://localhost:8081/health
```

### Manual Health Update

```go
manager.UpdateHealth(ctx, "user-api-v1", endpoint.HealthUpdate{
    Status:  endpoint.StatusMaintenance,
    Message: "Scheduled maintenance",
})
```

## HTTP Handlers

REST API integration:

```go
http.HandleFunc("/endpoints/register", manager.HandleRegistration)
http.HandleFunc("/endpoints/deregister", manager.HandleDeregistration)
http.HandleFunc("/endpoints/query", manager.HandleQuery)
http.HandleFunc("/endpoints/renew", manager.HandleRenew)
```

### Example Requests

**Register:**
```bash
POST /endpoints/register
Content-Type: application/json

{
    "id": "user-api-v1",
    "service_id": "user-service",
    "host": "localhost",
    "port": 8080,
    "protocol": "http",
    "ttl": "5m"
}
```

**Query:**
```bash
POST /endpoints/query
Content-Type: application/json

{
    "service_ids": ["user-service"],
    "only_healthy": true,
    "tags": ["api"]
}
```

## Complete Example

```go
package main

import (
    "context"
    "time"

    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/endpoint"
)

func main() {
    ctx := context.Background()

    // Create events client
    client, _ := events.NewClient(events.WithServers("nats://localhost:4222"))
    defer client.Shutdown(ctx)

    // Connect tenant
    conn, _ := client.ConnectWithJWT(ctx, "tenant-1", jwt, seed)

    // Create endpoint manager
    manager, _ := endpoint.NewManager(endpoint.ManagerConfig{
        Publisher: conn.Publisher(),
        HealthCheck: endpoint.HealthCheckConfig{
            Enabled:  true,
            Interval: 30 * time.Second,
        },
    })
    defer manager.Shutdown(ctx)

    // Register handlers for distributed discovery
    handler := endpoint.NewHandler(manager)
    handler.RegisterQueueSubscriptions(ctx, conn.Subscriber(), "endpoint-workers")

    // Register this service's endpoint
    ep := &endpoint.Info{
        ID:              "order-api-1",
        ServiceID:       "order-service",
        Host:            "localhost",
        Port:            8080,
        Path:            "/api/v1/orders",
        Protocol:        endpoint.ProtocolHTTP,
        HealthCheckPath: "/health",
        Tags:            []string{"api", "orders"},
    }
    manager.RegisterEndpoint(ctx, ep)

    // Discover other services
    userEndpoints := manager.Discover("user-service", endpoint.OnlyHealthy())
    for _, ep := range userEndpoints {
        fmt.Printf("Found: %s://%s:%d%s\n", ep.Protocol, ep.Host, ep.Port, ep.Path)
    }

    // Heartbeat loop
    go func() {
        ticker := time.NewTicker(time.Minute)
        for range ticker.C {
            client := endpoint.NewClient(conn.Publisher())
            client.Renew(ctx, "order-api-1", 5*time.Minute)
        }
    }()

    // Keep running...
    select {}
}
```

## Error Handling

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/endpoint"

err := manager.RegisterEndpoint(ctx, ep)
if err != nil {
    switch {
    case errors.Is(err, endpoint.ErrEndpointExists):
        // Already registered
    case errors.Is(err, endpoint.ErrInvalidEndpoint):
        // Missing required fields
    case errors.Is(err, endpoint.ErrMissingEndpointID):
        // ID is required
    }
}

ep, ok := registry.Get("unknown-id")
if !ok {
    // Endpoint not found
}
```

## Best Practices

### 1. Use TTL with Heartbeats

```go
// Register with TTL
manager.RegisterEndpoint(ctx, ep)
registry.RegisterWithTTL(ep, 5*time.Minute)

// Heartbeat at half the TTL interval
ticker := time.NewTicker(2*time.Minute + 30*time.Second)
for range ticker.C {
    registry.Renew(ep.ID, 5*time.Minute)
}
```

### 2. Handle Discovery Failures

```go
endpoints := manager.Discover("user-service", endpoint.OnlyHealthy())
if len(endpoints) == 0 {
    // No healthy endpoints - use fallback or retry
    time.Sleep(time.Second)
    endpoints = manager.Discover("user-service") // Include unhealthy
}
```

### 3. Monitor Health Changes

```go
registry.OnStatusChange(func(ep *endpoint.Info, old, new endpoint.Status) {
    if new == endpoint.StatusUnhealthy {
        alerting.Warn("Endpoint unhealthy", ep.ID)
    }
})
```

### 4. Graceful Shutdown

```go
// Deregister before shutdown
manager.DeregisterEndpoint(ctx, "my-endpoint-id")

// Then shutdown manager
manager.Shutdown(ctx)
```

## Thread Safety

All components are thread-safe for concurrent use.

## Testing

```bash
go test ./events/endpoint/...
go test -cover ./events/endpoint/...
```
