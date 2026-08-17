# Microservice Package

The `microservice` package provides microservice registration, discovery, and lifecycle management for distributed systems.

## Overview

This package enables services to:

- **Register Services**: Publish service metadata to a central registry
- **Manage Instances**: Track multiple running instances per service
- **Health Checking**: Automatic HTTP health checks for instances
- **Service Discovery**: Find services by capabilities, type, or metadata
- **TTL Management**: Automatic cleanup of expired instances
- **Event Publishing**: Lifecycle event notifications

## Package Structure

```
microservice/
├── types.go        # Info, Instance, Status, ServiceType
├── registry.go     # In-memory service registry
├── manager.go      # Service lifecycle manager
├── handlers.go     # NATS message handlers
├── errors.go       # Error definitions
├── doc.go          # Package documentation
├── example_test.go # Example tests
└── README.md       # This file
```

## Key Types

### Info (Service)

Service metadata and configuration:

```go
type Info struct {
    // Identification
    ID          string  // Unique service identifier
    Name        string  // Service name
    DisplayName string  // Human-readable name
    Description string  // Service description

    // Classification
    Type       ServiceType  // API, Worker, Gateway, etc.
    Version    string       // Service version
    APIVersion string       // API version (v1, v2, etc.)

    // Dependencies
    Dependencies []string  // Required services
    Capabilities []string  // Provided capabilities

    // Status
    Status Status  // Current status

    // Metadata
    Tags     []string          // Searchable tags
    Metadata map[string]string // Custom metadata

    // Multi-tenancy
    TenantID string  // Tenant identifier
    Owner    string  // Team/owner

    // Timestamps
    RegisteredAt time.Time
    UpdatedAt    time.Time
}
```

### Instance

Running service instance:

```go
type Instance struct {
    // Identification
    ID        string  // Unique instance identifier
    ServiceID string  // Parent service ID

    // Network
    Host     string  // Hostname or IP
    Port     int     // Port number
    Protocol string  // http, https, grpc

    // Status
    Status          Status  // Current status
    HealthCheckPath string  // Path for health checks

    // Load balancing
    Weight int  // Load balancing weight (higher = more traffic)

    // Location
    Region string  // Deployment region
    Zone   string  // Availability zone

    // Kubernetes metadata
    PodName   string  // Pod name
    NodeName  string  // Node name
    Namespace string  // Namespace

    // Metadata
    Labels   map[string]string  // Instance labels
    Metadata map[string]string  // Custom metadata

    // Timestamps
    RegisteredAt time.Time
    LastSeenAt   time.Time
    ExpiresAt    time.Time
}
```

### ServiceType

```go
type ServiceType string

const (
    ServiceTypeAPI       ServiceType = "api"        // REST/gRPC APIs
    ServiceTypeWorker    ServiceType = "worker"     // Background workers
    ServiceTypeGateway   ServiceType = "gateway"    // API gateways
    ServiceTypeBroker    ServiceType = "broker"     // Message brokers
    ServiceTypeDatabase  ServiceType = "database"   // Database services
    ServiceTypeCache     ServiceType = "cache"      // Cache services
    ServiceTypeQueue     ServiceType = "queue"      // Queue services
    ServiceTypeScheduler ServiceType = "scheduler"  // Job schedulers
    ServiceTypeMonitor   ServiceType = "monitor"    // Monitoring
    ServiceTypeCustom    ServiceType = "custom"     // Custom types
)
```

### Status

```go
type Status int

const (
    StatusUnknown     Status = iota  // Status not determined
    StatusStarting                   // Service is starting
    StatusRunning                    // Running normally
    StatusDegraded                   // Running with issues
    StatusStopping                   // Shutting down
    StatusStopped                    // Stopped
    StatusFailed                     // Failed
    StatusMaintenance                // Under maintenance
)
```

## Registry

Thread-safe storage for services and instances:

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/microservice"

// Create registry
registry := microservice.NewRegistry()

// Or with custom config
registry := microservice.NewRegistry(microservice.RegistryConfig{
    CleanupInterval:   30 * time.Second,
    ExpirationEnabled: true,
})
```

### Service Operations

```go
// Register service
svc := &microservice.Info{
    ID:           "user-service",
    Name:         "user-service",
    Type:         microservice.ServiceTypeAPI,
    Version:      "2.0.0",
    Capabilities: []string{"auth", "user-management"},
    Dependencies: []string{"database-service"},
}
registry.RegisterService(svc)

// Get service
svc, ok := registry.GetService("user-service")

// Update service
registry.UpdateService(svc)

// Query services
services := registry.QueryServices(&microservice.Query{
    Types:       []microservice.ServiceType{microservice.ServiceTypeAPI},
    Capabilities: []string{"auth"},
    OnlyHealthy: true,
})

// Deregister service
registry.DeregisterService("user-service")
```

### Instance Operations

```go
// Register instance
inst := &microservice.Instance{
    ID:              "user-service-1",
    ServiceID:       "user-service",
    Host:            "10.0.1.50",
    Port:            8080,
    Protocol:        "http",
    Status:          microservice.StatusRunning,
    HealthCheckPath: "/health",
    Weight:          100,
    Region:          "us-east-1",
    Zone:            "us-east-1a",
}
registry.RegisterInstance(inst)

// Register with TTL
registry.RegisterInstanceWithTTL(inst, 30*time.Second)

// Get instances
instances := registry.GetInstances("user-service")
healthyInstances := registry.GetHealthyInstances("user-service")

// Renew TTL (heartbeat)
registry.RenewInstance("user-service-1", 30*time.Second)

// Update status
registry.UpdateInstanceStatus("user-service-1", microservice.StatusDegraded)

// Deregister instance
registry.DeregisterInstance("user-service-1")
```

### Callbacks

```go
registry.OnServiceRegistered(func(svc *microservice.Info) {
    log.Printf("Service registered: %s", svc.Name)
})

registry.OnInstanceRegistered(func(inst *microservice.Instance) {
    log.Printf("Instance registered: %s@%s", inst.ID, inst.ServiceID)
})

registry.OnInstanceStatusChange(func(inst *microservice.Instance, old, new microservice.Status) {
    log.Printf("Instance %s: %s -> %s", inst.ID, old, new)
})

registry.OnInstanceExpired(func(inst *microservice.Instance) {
    log.Printf("Instance expired: %s", inst.ID)
})
```

## Manager

Handles service lifecycle with health checking:

```go
manager, err := microservice.NewManager(microservice.ManagerConfig{
    Registry:  registry,
    Publisher: publisher,
    HealthCheck: microservice.HealthCheckConfig{
        Enabled:            true,
        Interval:           30 * time.Second,
        Timeout:            5 * time.Second,
        HealthyThreshold:   2,
        UnhealthyThreshold: 3,
    },
})
defer manager.Shutdown(ctx)
```

### Manager Methods

| Method | Description |
|--------|-------------|
| `RegisterServiceDirect(ctx, svc)` | Register service |
| `RegisterInstanceDirect(ctx, inst)` | Register instance |
| `UpdateService(ctx, svc)` | Update service |
| `UpdateInstance(ctx, inst)` | Update instance |
| `Discover(serviceID, opts...)` | Discover services |
| `DiscoverInstances(serviceID, healthyOnly)` | Get instances |
| `Shutdown(ctx)` | Graceful shutdown |

### Discovery Options

```go
services := manager.Discover("user-service",
    microservice.WithType(microservice.ServiceTypeAPI),
    microservice.WithCapabilities("auth", "oauth"),
    microservice.WithVersion("2.0.0"),
    microservice.OnlyHealthy(),
    microservice.HasInstances(),
)

instances := manager.DiscoverInstances("user-service", true)
```

## Handler

NATS message handlers for distributed operations:

```go
handler := microservice.NewHandler(manager)

// Register subscriptions
subscriptions, err := handler.RegisterSubscriptions(ctx, subscriber)

// Or with queue groups
subscriptions, err := handler.RegisterQueueSubscriptions(ctx, subscriber, "service-registry")
```

### Message Subjects

| Subject | Description |
|---------|-------------|
| `services.register` | Register service |
| `services.deregister` | Deregister service |
| `services.query` | Query services |
| `instances.register` | Register instance |
| `instances.deregister` | Deregister instance |
| `instances.renew` | Renew instance TTL |
| `instances.health` | Update instance health |

### Published Events

| Subject | Description |
|---------|-------------|
| `services.registered` | Service was registered |
| `services.deregistered` | Service was deregistered |
| `services.updated` | Service was updated |
| `instances.registered` | Instance was registered |
| `instances.deregistered` | Instance was deregistered |
| `instances.health.changed` | Instance health changed |

## Client

Remote operations over messaging:

```go
client := microservice.NewClient(publisher)

// Register service
resp, err := client.RegisterService(ctx, svc)

// Register instance with TTL
resp, err := client.RegisterInstance(ctx, inst, 30*time.Second)

// Discover services
resp, err := client.Discover(ctx, "user-service",
    microservice.WithCapabilities("auth"),
)

// Heartbeat
resp, err := client.RenewInstance(ctx, "user-service-1", 30*time.Second)
```

## Health Checking

### Configuration

```go
healthConfig := microservice.HealthCheckConfig{
    Enabled:            true,
    Interval:           30 * time.Second,
    Timeout:            5 * time.Second,
    HealthyThreshold:   2,
    UnhealthyThreshold: 3,
}
```

### Health Check Behavior

| HTTP Status | Result |
|-------------|--------|
| 2xx | Running |
| 503 | Maintenance |
| Other/Error | Failed/Degraded |

### Instance Health Path

```go
inst := &microservice.Instance{
    Host:            "localhost",
    Port:            8080,
    Protocol:        "http",
    HealthCheckPath: "/health",
}

// Health check URL: http://localhost:8080/health
```

## HTTP Handlers

REST API integration:

```go
http.HandleFunc("/services/register", manager.HandleServiceRegistration)
http.HandleFunc("/services/query", manager.HandleServiceQuery)
http.HandleFunc("/instances/register", manager.HandleInstanceRegistration)
http.HandleFunc("/instances/renew", manager.HandleInstanceRenew)
```

## Complete Example

```go
package main

import (
    "context"
    "time"

    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/microservice"
)

func main() {
    ctx := context.Background()

    // Create events client
    client, _ := events.NewClient(events.WithServers("nats://localhost:4222"))
    defer client.Shutdown(ctx)

    conn, _ := client.ConnectWithJWT(ctx, "tenant-1", jwt, seed)

    // Create microservice manager
    manager, _ := microservice.NewManager(microservice.ManagerConfig{
        Publisher: conn.Publisher(),
        HealthCheck: microservice.HealthCheckConfig{
            Enabled:  true,
            Interval: 30 * time.Second,
        },
    })
    defer manager.Shutdown(ctx)

    // Register handlers
    handler := microservice.NewHandler(manager)
    handler.RegisterQueueSubscriptions(ctx, conn.Subscriber(), "service-registry")

    // Register this service
    svc := &microservice.Info{
        ID:           "order-service",
        Name:         "order-service",
        Type:         microservice.ServiceTypeAPI,
        Version:      "1.0.0",
        Capabilities: []string{"orders", "checkout"},
        Dependencies: []string{"user-service", "inventory-service"},
    }
    manager.RegisterServiceDirect(ctx, svc)

    // Register this instance
    inst := &microservice.Instance{
        ID:              "order-service-1",
        ServiceID:       "order-service",
        Host:            getHostIP(),
        Port:            8080,
        Protocol:        "http",
        HealthCheckPath: "/health",
        Weight:          100,
        PodName:         os.Getenv("POD_NAME"),
    }
    manager.RegisterInstanceDirect(ctx, inst)

    // Heartbeat loop
    go func() {
        ticker := time.NewTicker(10 * time.Second)
        for range ticker.C {
            client := microservice.NewClient(conn.Publisher())
            client.RenewInstance(ctx, inst.ID, 30*time.Second)
        }
    }()

    // Discover dependencies
    userService := manager.Discover("user-service", microservice.OnlyHealthy())
    if len(userService) > 0 {
        instances := manager.DiscoverInstances("user-service", true)
        for _, inst := range instances {
            fmt.Printf("User service: %s:%d\n", inst.Host, inst.Port)
        }
    }

    // Keep running...
    select {}
}
```

## Capabilities and Dependencies

### Declaring Capabilities

```go
svc := &microservice.Info{
    ID:           "payment-service",
    Capabilities: []string{"payments", "refunds", "subscriptions"},
}
```

### Finding by Capability

```go
// Find services that can handle payments
services := registry.QueryServices(&microservice.Query{
    Capabilities: []string{"payments"},
    OnlyHealthy:  true,
})
```

### Dependency Tracking

```go
svc := &microservice.Info{
    ID:           "order-service",
    Dependencies: []string{"user-service", "payment-service", "inventory-service"},
}

// Check if dependencies are available
for _, dep := range svc.Dependencies {
    if services := manager.Discover(dep, microservice.OnlyHealthy()); len(services) == 0 {
        log.Printf("Warning: dependency %s has no healthy instances", dep)
    }
}
```

## Load Balancing

### Weight-Based Selection

```go
inst1 := &microservice.Instance{
    ID:        "service-1",
    ServiceID: "my-service",
    Weight:    100,  // High priority
}

inst2 := &microservice.Instance{
    ID:        "service-2",
    ServiceID: "my-service",
    Weight:    50,   // Lower priority
}

// Implement weighted selection in your load balancer
instances := manager.DiscoverInstances("my-service", true)
selected := weightedRandom(instances)
```

### Region/Zone Affinity

```go
// Query instances in same region
instances := registry.QueryInstances(&microservice.InstanceQuery{
    ServiceID: "user-service",
    Region:    "us-east-1",
    Zone:      "us-east-1a",
    OnlyHealthy: true,
})
```

## Error Handling

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/microservice"

// Service errors
err := manager.RegisterServiceDirect(ctx, svc)
switch {
case errors.Is(err, microservice.ErrServiceExists):
    // Already registered
case errors.Is(err, microservice.ErrInvalidService):
    // Missing required fields
case errors.Is(err, microservice.ErrMissingServiceID):
    // ID is required
}

// Instance errors
err = manager.RegisterInstanceDirect(ctx, inst)
switch {
case errors.Is(err, microservice.ErrInstanceExists):
    // Already registered
case errors.Is(err, microservice.ErrInvalidInstance):
    // Missing required fields
case errors.Is(err, microservice.ErrServiceNotFound):
    // Parent service not registered
}
```

## Best Practices

### 1. Register Service Before Instances

```go
// First register service metadata
manager.RegisterServiceDirect(ctx, svc)

// Then register instances
manager.RegisterInstanceDirect(ctx, inst)
```

### 2. Use TTL with Heartbeats

```go
// Register with TTL
registry.RegisterInstanceWithTTL(inst, 30*time.Second)

// Heartbeat at 1/3 the TTL
ticker := time.NewTicker(10 * time.Second)
for range ticker.C {
    registry.RenewInstance(inst.ID, 30*time.Second)
}
```

### 3. Graceful Shutdown

```go
// Deregister instance first
manager.DeregisterInstance(ctx, inst.ID)

// Wait for in-flight requests
time.Sleep(5 * time.Second)

// Shutdown
manager.Shutdown(ctx)
```

### 4. Monitor Service Health

```go
registry.OnInstanceStatusChange(func(inst *microservice.Instance, old, new microservice.Status) {
    metrics.SetGauge("instance_status",
        float64(new),
        "service", inst.ServiceID,
        "instance", inst.ID,
    )

    if new == microservice.StatusFailed {
        alerting.Critical("Instance failed", inst.ID)
    }
})
```

### 5. Check Dependencies at Startup

```go
func checkDependencies(manager *microservice.Manager, deps []string) error {
    for _, dep := range deps {
        services := manager.Discover(dep, microservice.OnlyHealthy())
        if len(services) == 0 {
            return fmt.Errorf("dependency %s not available", dep)
        }

        instances := manager.DiscoverInstances(dep, true)
        if len(instances) == 0 {
            return fmt.Errorf("dependency %s has no healthy instances", dep)
        }
    }
    return nil
}
```

## Step-by-Step Usage Guide

### Step 1: Create Microservice Manager

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/microservice"

manager, err := microservice.NewManager(microservice.ManagerConfig{
    Publisher: pub,
    HealthCheck: microservice.HealthCheckConfig{
        Enabled:            true,
        Interval:           30 * time.Second,
        Timeout:            5 * time.Second,
        HealthyThreshold:   2,
        UnhealthyThreshold: 3,
    },
})
defer manager.Shutdown(ctx)
```

### Step 2: Register Your Service

```go
svc := &microservice.Info{
    ID:           "order-service",
    Name:         "order-service",
    DisplayName:  "Order Service",
    Type:         microservice.ServiceTypeAPI,
    Version:      "2.0.0",
    Capabilities: []string{"orders", "checkout"},
    Dependencies: []string{"user-service", "inventory-service"},
}
manager.RegisterServiceDirect(ctx, svc)
```

### Step 3: Register Instance

```go
inst := &microservice.Instance{
    ID:              "order-service-1",
    ServiceID:       "order-service",
    Host:            getHostIP(),
    Port:            8080,
    Protocol:        "http",
    HealthCheckPath: "/health",
    Weight:          100,
    Region:          "us-east-1",
    Zone:            "us-east-1a",
}
manager.RegisterInstanceDirect(ctx, inst)
```

### Step 4: Set Up Heartbeat

```go
go func() {
    ticker := time.NewTicker(10 * time.Second) // 1/3 of 30s TTL
    for range ticker.C {
        manager.RenewInstance(inst.ID, 30*time.Second)
    }
}()
```

### Step 5: Set Up Distributed Handlers

```go
handler := microservice.NewHandler(manager)
handler.RegisterQueueSubscriptions(ctx, sub, "service-registry")
```

### Step 6: Discover Other Services

```go
services := manager.Discover("user-service",
    microservice.OnlyHealthy(),
    microservice.WithCapabilities("auth"),
)

instances := manager.DiscoverInstances("user-service", true)
for _, inst := range instances {
    fmt.Printf("Found: %s:%d (weight: %d)\n", inst.Host, inst.Port, inst.Weight)
}
```

### Step 7: Check Dependencies at Startup

```go
for _, dep := range svc.Dependencies {
    if services := manager.Discover(dep, microservice.OnlyHealthy()); len(services) == 0 {
        log.Printf("WARNING: Dependency %s not available", dep)
    }
}
```

### Step 8: Monitor Health Changes

```go
registry := manager.Registry()
registry.OnInstanceStatusChange(func(inst *microservice.Instance, old, new microservice.Status) {
    if new == microservice.StatusFailed {
        log.Printf("ALERT: Instance %s failed!", inst.ID)
    }
})
```

### Step 9: Graceful Shutdown

```go
manager.DeregisterInstance(ctx, inst.ID)
time.Sleep(5 * time.Second) // Let in-flight requests finish
manager.Shutdown(ctx)
```

## Thread Safety

All components are thread-safe for concurrent use.

## Testing

```bash
go test ./events/microservice/...
go test -cover ./events/microservice/...
```
