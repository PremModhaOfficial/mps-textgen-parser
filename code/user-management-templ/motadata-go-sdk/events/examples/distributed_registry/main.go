// Package main demonstrates a distributed service registry pattern using
// NATS messaging for service registration and discovery.
//
// This example shows:
//   - Setting up a registry server that handles registration requests
//   - Using handlers for NATS-based registration
//   - Client-side service registration and discovery
//   - Heartbeat/TTL renewal pattern
//
// Note: This example requires a running NATS server. For demonstration
// purposes, it shows the code structure without actual NATS connection.
//
// Architecture:
//
//	+------------------+       NATS        +------------------+
//	|  Service A       | <===============> |  Registry Server |
//	|  (registers)     |                   |  (handlers)      |
//	+------------------+                   +------------------+
//	        ^                                      ^
//	        |              NATS                    |
//	        v                                      v
//	+------------------+       NATS        +------------------+
//	|  Service B       | <===============> |  Service C       |
//	|  (discovers)     |                   |  (registers)     |
//	+------------------+                   +------------------+
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/endpoint"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/microservice"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"
)

func main() {
	ctx := context.Background()
	log := logger.L()

	log.Info(ctx, "=== Distributed Service Registry Example ===")

	// This example demonstrates the patterns without requiring NATS
	demonstrateRegistryServer(ctx, log)
	demonstrateServiceRegistration(ctx, log)
	demonstrateServiceDiscovery(ctx, log)
	demonstrateHeartbeatPattern(ctx, log)

	log.Info(ctx, "=== Example Complete ===")
}

// demonstrateRegistryServer shows how to set up a registry server
func demonstrateRegistryServer(ctx context.Context, log *logger.Logger) {
	log.Info(ctx, "1. Registry Server Setup")
	log.Info(ctx, "-------------------------")

	// In production, you would create these with NATS connection
	// For this example, we show the setup pattern

	os.Stdout.WriteString(`
// Create managers
svcManager, _ := microservice.NewManager(microservice.ManagerConfig{
    HealthCheck: microservice.HealthCheckConfig{
        Enabled:  true,
        Interval: 30 * time.Second,
    },
})

epManager, _ := endpoint.NewManager(endpoint.ManagerConfig{
    HealthCheck: endpoint.HealthCheckConfig{
        Enabled:  true,
        Interval: 30 * time.Second,
    },
})

// Create handlers
svcHandler := microservice.NewHandler(svcManager)
epHandler := endpoint.NewHandler(epManager)

// Register subscriptions (requires NATS subscriber)
// svcSubs, _ := svcHandler.RegisterQueueSubscriptions(ctx, subscriber, "registry")
// epSubs, _ := epHandler.RegisterQueueSubscriptions(ctx, subscriber, "registry")

// Subjects being handled:
// - services.register
// - services.deregister
// - services.query
// - services.discover
// - instances.register
// - instances.deregister
// - instances.renew
// - endpoints.register
// - endpoints.deregister
// - endpoints.query
// - endpoints.discover
`)

	// Create actual managers to show they work
	svcManager, _ := microservice.NewManager(microservice.ManagerConfig{
		HealthCheck: microservice.HealthCheckConfig{Enabled: false},
	})
	epManager, _ := endpoint.NewManager(endpoint.ManagerConfig{
		HealthCheck: endpoint.HealthCheckConfig{Enabled: false},
	})

	// Create handlers (they work even without NATS for local operations)
	_ = microservice.NewHandler(svcManager)
	_ = endpoint.NewHandler(epManager)

	log.Info(ctx, "[OK] Registry server components created")
}

// demonstrateServiceRegistration shows how a service registers itself
func demonstrateServiceRegistration(ctx context.Context, log *logger.Logger) {
	log.Info(ctx, "2. Service Registration Pattern")
	log.Info(ctx, "--------------------------------")

	os.Stdout.WriteString(`
// Service startup registration pattern
func registerService(ctx context.Context, client *microservice.Client) error {
    // 1. Register the service
    svc := &microservice.Info{
        ID:           "my-service",
        Name:         "my-service",
        Version:      "1.0.0",
        Type:         microservice.ServiceTypeAPI,
        Capabilities: []string{"feature-a", "feature-b"},
    }

    resp, err := client.RegisterService(ctx, svc)
    if err != nil {
        return err
    }
    log.Printf("Service registered: %s", resp.ServiceID)

    // 2. Register this instance
    inst := &microservice.Instance{
        ID:              hostname + "-" + uuid.New().String()[:8],
        ServiceID:       "my-service",
        Host:            getLocalIP(),
        Port:            8080,
        HealthCheckPath: "/health",
    }

    // Register with 30-second TTL
    resp, err = client.RegisterInstance(ctx, inst, 30*time.Second)
    if err != nil {
        return err
    }
    log.Printf("Instance registered: %s", resp.InstanceID)

    // 3. Start heartbeat goroutine
    go heartbeatLoop(ctx, client, inst.ID)

    return nil
}
`)
	var serviceName = "demo-service"
	// Demonstrate with local registry
	registry := microservice.NewRegistry()

	svc := &microservice.Info{
		ID:           serviceName,
		Name:         serviceName,
		Version:      "1.0.0",
		Type:         microservice.ServiceTypeAPI,
		Status:       microservice.StatusRunning,
		Capabilities: []string{"demo"},
	}
	registry.RegisterService(svc)

	inst := &microservice.Instance{
		ID:              "demo-instance-1",
		ServiceID:       serviceName,
		Host:            "127.0.0.1",
		Port:            8080,
		Status:          microservice.StatusRunning,
		HealthCheckPath: "/health",
	}
	registry.RegisterInstanceWithTTL(inst, 30*time.Second)

	log.Info(ctx, "[OK] Service and instance registered")
	log.Info(ctx, "Service registered", logger.String("name", svc.Name), logger.String("version", svc.Version))
	log.Info(ctx, "Instance registered", logger.String("id", inst.ID), logger.String("address", inst.Address()))
}

// demonstrateServiceDiscovery shows how to discover other services
func demonstrateServiceDiscovery(ctx context.Context, log *logger.Logger) {
	log.Info(ctx, "3. Service Discovery Pattern")
	log.Info(ctx, "-----------------------------")

	os.Stdout.WriteString(`
// Discover and call another service
func callUserService(ctx context.Context, client *microservice.Client) error {
    // 1. Discover the service
    resp, err := client.Discover(ctx, "user-service",
        microservice.WithCapabilities("auth"),
        microservice.OnlyHealthy(),
    )
    if err != nil {
        return err
    }

    if resp.Count == 0 {
        return errors.New("no user-service instances available")
    }

    // 2. Get instances for the service
    instResp, err := client.QueryInstances(ctx, "user-service", true)
    if err != nil {
        return err
    }

    // 3. Select an instance (simple round-robin or weighted)
    instance := selectInstance(instResp.Instances)

    // 4. Make the HTTP call
    url := fmt.Sprintf("http://%s/api/users", instance.Address())
    resp, err := http.Get(url)
    // ...
}

// Simple weighted selection
func selectInstance(instances []*microservice.Instance) *microservice.Instance {
    // Weight-based selection
    totalWeight := 0
    for _, inst := range instances {
        totalWeight += inst.Weight
    }

    r := rand.Intn(totalWeight)
    for _, inst := range instances {
        r -= inst.Weight
        if r < 0 {
            return inst
        }
    }
    return instances[0]
}
`)

	// Demonstrate with local registry
	registry := microservice.NewRegistry()

	// Register some services
	for _, name := range []string{"user-service", "order-service"} {
		registry.RegisterService(&microservice.Info{
			ID:           name,
			Name:         name,
			Status:       microservice.StatusRunning,
			Capabilities: []string{"api"},
		})
		for i := 1; i <= 2; i++ {
			registry.RegisterInstance(&microservice.Instance{
				ID:        fmt.Sprintf("%s-inst-%d", name, i),
				ServiceID: name,
				Host:      fmt.Sprintf("10.0.%d.%d", i, i*10),
				Port:      8080,
				Status:    microservice.StatusRunning,
				Weight:    100,
			})
		}
	}

	// Discover
	services := registry.QueryServices(&microservice.Query{
		Capabilities:  []string{"api"},
		OnlyAvailable: true,
	})

	log.Info(ctx, "[OK] Discovered services with 'api' capability", logger.Int("count", len(services)))
	for _, svc := range services {
		instances := registry.GetHealthyInstances(svc.ID)
		log.Info(ctx, "Service instances", logger.String("name", svc.Name), logger.Int("healthy_instances", len(instances)))
	}
}

// demonstrateHeartbeatPattern shows the heartbeat/TTL renewal pattern
func demonstrateHeartbeatPattern(ctx context.Context, log *logger.Logger) {
	log.Info(ctx, "4. Heartbeat Pattern")
	log.Info(ctx, "--------------------")

	os.Stdout.WriteString(`
// Heartbeat loop to keep instance registration alive
func heartbeatLoop(ctx context.Context, client *microservice.Client, instanceID string) {
    ttl := 30 * time.Second
    interval := 10 * time.Second // Renew at 1/3 of TTL

    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            // Graceful shutdown - deregister
            client.DeregisterInstance(context.Background(), instanceID, "shutdown")
            return

        case <-ticker.C:
            resp, err := client.RenewInstance(ctx, instanceID, ttl)
            if err != nil {
                log.Printf("Heartbeat failed: %v", err)
                // Attempt to re-register if needed
                continue
            }
            log.Printf("Heartbeat OK, expires: %v", resp.ExpiresAt)
        }
    }
}

// Graceful shutdown pattern
func gracefulShutdown(ctx context.Context, client *microservice.Client, instanceID string) {
    // Create shutdown context with timeout
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // Deregister instance
    _, err := client.DeregisterInstance(shutdownCtx, instanceID, "graceful shutdown")
    if err != nil {
        log.Printf("Failed to deregister: %v", err)
    }

    // Deregister endpoints if any
    // ...
}
`)

	var serviceName = "heartbeat-demo"
	// Demonstrate heartbeat simulation
	registry := microservice.NewRegistry()
	registry.RegisterService(&microservice.Info{
		ID:     serviceName,
		Name:   serviceName,
		Status: microservice.StatusRunning,
	})

	inst := &microservice.Instance{
		ID:        "heartbeat-inst",
		ServiceID: serviceName,
		Host:      "localhost",
		Port:      8080,
		Status:    microservice.StatusRunning,
	}
	registry.RegisterInstanceWithTTL(inst, 30*time.Second)

	// Simulate heartbeat
	log.Info(ctx, "[OK] Instance registered with 30s TTL")

	// Renew
	registry.RenewInstance(inst.ID, 30*time.Second)
	log.Info(ctx, "[OK] TTL renewed")

	// Check instance
	retrieved, found := registry.GetInstance(inst.ID)
	if found {
		log.Info(ctx, "Instance expiration", logger.String("expires_at", retrieved.ExpiresAt.Format(time.RFC3339)))
	}
}

// Additional code patterns that would be useful

func init() {
	// This shows additional patterns as documentation

	_ = `
// HTTP Health Check Handler
func healthHandler(w http.ResponseWriter, r *http.Request) {
    health := struct {
        Status    string    ` + "`json:\"status\"`" + `
        Timestamp time.Time ` + "`json:\"timestamp\"`" + `
        Version   string    ` + "`json:\"version\"`" + `
    }{
        Status:    "healthy",
        Timestamp: time.Now(),
        Version:   version,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(health)
}

// Service Mesh Integration Pattern
type ServiceMesh struct {
    svcManager *microservice.Manager
    epManager  *endpoint.Manager
    svcClient  *microservice.Client
    epClient   *endpoint.Client
}

func NewServiceMesh(publisher core.Publisher) *ServiceMesh {
    return &ServiceMesh{
        svcClient: microservice.NewClient(publisher),
        epClient:  endpoint.NewClient(publisher),
    }
}

func (m *ServiceMesh) Register(ctx context.Context, svc *microservice.Info, inst *microservice.Instance, eps []*endpoint.Info) error {
    // Register service
    if _, err := m.svcClient.RegisterService(ctx, svc); err != nil {
        return fmt.Errorf("register service: %w", err)
    }

    // Register instance
    if _, err := m.svcClient.RegisterInstance(ctx, inst, 30*time.Second); err != nil {
        return fmt.Errorf("register instance: %w", err)
    }

    // Register endpoints
    for _, ep := range eps {
        if _, err := m.epClient.Register(ctx, ep, 30*time.Second); err != nil {
            return fmt.Errorf("register endpoint %s: %w", ep.ID, err)
        }
    }

    return nil
}

func (m *ServiceMesh) Discover(ctx context.Context, serviceName string) ([]*microservice.Instance, error) {
    resp, err := m.svcClient.QueryInstances(ctx, serviceName, true)
    if err != nil {
        return nil, err
    }
    return resp.Instances, nil
}
`
}
