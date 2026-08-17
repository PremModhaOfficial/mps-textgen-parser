package microservice_test

import (
	"context"
	"fmt"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/microservice"
)

// Example_basicServiceRegistration demonstrates basic service registration.
func Example_basicServiceRegistration() {
	registry := microservice.NewRegistry()

	// Create a service
	svc := &microservice.Info{
		ID:          "user-service",
		Name:        "user-service",
		DisplayName: "User Service",
		Type:        microservice.ServiceTypeAPI,
		Version:     "1.0.0",
		Status:      microservice.StatusRunning,
	}

	// Register the service
	if err := registry.RegisterService(svc); err != nil {
		fmt.Printf("Failed to register: %v\n", err)
		return
	}

	// Retrieve the service
	retrieved, found := registry.GetService("user-service")
	if found {
		fmt.Printf("Service: %s (v%s)\n", retrieved.DisplayName, retrieved.Version)
	}

	// Output:
	// Service: User Service (v1.0.0)
}

// Example_serviceWithInstances demonstrates registering a service with instances.
func Example_serviceWithInstances() {
	registry := microservice.NewRegistry()

	// Register service first
	svc := &microservice.Info{
		ID:     "api-service",
		Name:   "api-service",
		Type:   microservice.ServiceTypeAPI,
		Status: microservice.StatusRunning,
	}
	registry.RegisterService(svc)

	// Register instances
	instances := []*microservice.Instance{
		{
			ID:        "api-instance-1",
			ServiceID: "api-service",
			Host:      "10.0.1.10",
			Port:      8080,
			Status:    microservice.StatusRunning,
			Weight:    100,
		},
		{
			ID:        "api-instance-2",
			ServiceID: "api-service",
			Host:      "10.0.1.11",
			Port:      8080,
			Status:    microservice.StatusRunning,
			Weight:    100,
		},
		{
			ID:        "api-instance-3",
			ServiceID: "api-service",
			Host:      "10.0.1.12",
			Port:      8080,
			Status:    microservice.StatusDegraded,
			Weight:    50,
		},
	}

	for _, inst := range instances {
		registry.RegisterInstance(inst)
	}

	// Get all instances
	allInstances := registry.GetInstances("api-service")
	fmt.Printf("Total instances: %d\n", len(allInstances))

	// Get only healthy instances
	healthyInstances := registry.GetHealthyInstances("api-service")
	fmt.Printf("Healthy instances: %d\n", len(healthyInstances))

	// Get available instances (healthy or degraded)
	availableInstances := registry.GetAvailableInstances("api-service")
	fmt.Printf("Available instances: %d\n", len(availableInstances))

	// Output:
	// Total instances: 3
	// Healthy instances: 2
	// Available instances: 3
}

// Example_instanceTTL demonstrates instance registration with TTL.
func Example_instanceTTL() {
	registry := microservice.NewRegistry()

	// Register service
	registry.RegisterService(&microservice.Info{
		ID:     "worker-service",
		Name:   "worker-service",
		Type:   microservice.ServiceTypeWorker,
		Status: microservice.StatusRunning,
	})

	// Register instance with TTL
	inst := &microservice.Instance{
		ID:        "worker-1",
		ServiceID: "worker-service",
		Host:      "worker-host",
		Port:      9000,
		Status:    microservice.StatusRunning,
	}

	registry.RegisterInstanceWithTTL(inst, 30*time.Second)

	// Check if instance exists
	retrieved, found := registry.GetInstance("worker-1")
	if found {
		fmt.Printf("Instance: %s, expires: %v\n", retrieved.ID, !retrieved.ExpiresAt.IsZero())
	}

	// Renew the TTL
	registry.RenewInstance("worker-1", 30*time.Second)
	fmt.Println("TTL renewed")

	// Output:
	// Instance: worker-1, expires: true
	// TTL renewed
}

// Example_queryServices demonstrates querying services with filters.
func Example_queryServices() {
	registry := microservice.NewRegistry()

	// Register various services
	services := []*microservice.Info{
		{
			ID:           "user-api",
			Name:         "user-api",
			Type:         microservice.ServiceTypeAPI,
			Version:      "1.0.0",
			Status:       microservice.StatusRunning,
			Capabilities: []string{"auth", "users"},
			Tags:         []string{"core"},
		},
		{
			ID:           "order-api",
			Name:         "order-api",
			Type:         microservice.ServiceTypeAPI,
			Version:      "2.0.0",
			Status:       microservice.StatusRunning,
			Capabilities: []string{"orders", "checkout"},
			Tags:         []string{"commerce"},
		},
		{
			ID:           "email-worker",
			Name:         "email-worker",
			Type:         microservice.ServiceTypeWorker,
			Version:      "1.0.0",
			Status:       microservice.StatusRunning,
			Capabilities: []string{"email", "notifications"},
			Tags:         []string{"background"},
		},
		{
			ID:           "cache-service",
			Name:         "cache-service",
			Type:         microservice.ServiceTypeCache,
			Version:      "1.0.0",
			Status:       microservice.StatusDegraded,
			Capabilities: []string{"caching"},
			Tags:         []string{"infrastructure"},
		},
	}

	for _, svc := range services {
		registry.RegisterService(svc)
	}

	// Query by type
	apiServices := registry.QueryServices(&microservice.Query{
		Types: []microservice.ServiceType{microservice.ServiceTypeAPI},
	})
	fmt.Printf("API services: %d\n", len(apiServices))

	// Query healthy services
	healthyServices := registry.QueryServices(&microservice.Query{
		OnlyHealthy: true,
	})
	fmt.Printf("Healthy services: %d\n", len(healthyServices))

	// Query by capability
	authServices := registry.QueryServices(&microservice.Query{
		Capabilities: []string{"auth"},
	})
	fmt.Printf("Auth-capable services: %d\n", len(authServices))

	// Query by tag
	coreServices := registry.QueryServices(&microservice.Query{
		Tags: []string{"core"},
	})
	fmt.Printf("Core services: %d\n", len(coreServices))

	// Output:
	// API services: 2
	// Healthy services: 3
	// Auth-capable services: 1
	// Core services: 1
}

// Example_manager demonstrates using the microservice manager.
func Example_manager() {
	manager, err := microservice.NewManager(microservice.ManagerConfig{
		HealthCheck: microservice.HealthCheckConfig{
			Enabled: false, // Disabled for example
		},
	})
	if err != nil {
		fmt.Printf("Failed to create manager: %v\n", err)
		return
	}

	ctx := context.Background()

	// Register a service
	svc := &microservice.Info{
		ID:           "payment-service",
		Name:         "payment-service",
		Type:         microservice.ServiceTypeAPI,
		Version:      "3.0.0",
		Status:       microservice.StatusRunning,
		Capabilities: []string{"payments", "refunds"},
	}
	manager.RegisterServiceDirect(ctx, svc)

	// Register instances
	for i := 1; i <= 3; i++ {
		inst := &microservice.Instance{
			ID:        fmt.Sprintf("payment-instance-%d", i),
			ServiceID: "payment-service",
			Host:      fmt.Sprintf("10.0.1.%d", 10+i),
			Port:      8080,
			Status:    microservice.StatusRunning,
		}
		manager.RegisterInstanceDirect(ctx, inst)
	}

	// Discover services
	services := manager.Discover("payment-service",
		microservice.WithCapabilities("payments"),
		microservice.OnlyHealthy(),
	)
	fmt.Printf("Found %d services\n", len(services))

	// Get instances for load balancing
	instances := manager.DiscoverInstances("payment-service", true)
	fmt.Printf("Found %d healthy instances\n", len(instances))

	// Shutdown
	manager.Shutdown(ctx)

	// Output:
	// Found 1 services
	// Found 3 healthy instances
}

// Example_callbacks demonstrates using registry callbacks.
// Note: Callbacks run in goroutines, so we use channels to synchronize.
func Example_callbacks() {
	registry := microservice.NewRegistry()

	svcRegistered := make(chan string, 1)
	instRegistered := make(chan string, 1)
	statusChanged := make(chan string, 1)

	// Set up callbacks
	registry.OnServiceRegistered(func(svc *microservice.Info) {
		svcRegistered <- svc.Name
	})

	registry.OnInstanceRegistered(func(inst *microservice.Instance) {
		instRegistered <- inst.ID
	})

	registry.OnInstanceStatusChange(func(inst *microservice.Instance, old, new microservice.Status) {
		statusChanged <- fmt.Sprintf("%s: %s->%s", inst.ID, old, new)
	})

	// Register service
	registry.RegisterService(&microservice.Info{
		ID:     "callback-service",
		Name:   "callback-service",
		Status: microservice.StatusRunning,
	})

	// Wait for service registration callback
	fmt.Printf("Service registered: %s\n", <-svcRegistered)

	// Register instance
	registry.RegisterInstance(&microservice.Instance{
		ID:        "callback-instance",
		ServiceID: "callback-service",
		Host:      "localhost",
		Port:      8080,
		Status:    microservice.StatusRunning,
	})

	// Wait for instance registration callback
	fmt.Printf("Instance registered: %s\n", <-instRegistered)

	// Update status
	registry.UpdateInstanceStatus("callback-instance", microservice.StatusDegraded)

	// Wait for status change callback
	fmt.Printf("Status change: %s\n", <-statusChanged)

	// Output:
	// Service registered: callback-service
	// Instance registered: callback-instance
	// Status change: callback-instance: running->degraded
}

// Example_dependencies demonstrates service dependencies.
func Example_dependencies() {
	registry := microservice.NewRegistry()

	// Register services with dependencies
	registry.RegisterService(&microservice.Info{
		ID:           "order-service",
		Name:         "order-service",
		Type:         microservice.ServiceTypeAPI,
		Status:       microservice.StatusRunning,
		Dependencies: []string{"user-service", "inventory-service", "payment-service"},
		Capabilities: []string{"orders", "checkout"},
	})

	registry.RegisterService(&microservice.Info{
		ID:           "user-service",
		Name:         "user-service",
		Type:         microservice.ServiceTypeAPI,
		Status:       microservice.StatusRunning,
		Capabilities: []string{"users", "auth"},
	})

	registry.RegisterService(&microservice.Info{
		ID:           "inventory-service",
		Name:         "inventory-service",
		Type:         microservice.ServiceTypeAPI,
		Status:       microservice.StatusRunning,
		Capabilities: []string{"inventory", "stock"},
	})

	// Get service with dependencies
	svc, _ := registry.GetService("order-service")
	fmt.Printf("Service: %s\n", svc.Name)
	fmt.Printf("Dependencies: %v\n", svc.Dependencies)
	fmt.Printf("Capabilities: %v\n", svc.Capabilities)

	// Output:
	// Service: order-service
	// Dependencies: [user-service inventory-service payment-service]
	// Capabilities: [orders checkout]
}

// Example_serviceInfo demonstrates service info utilities.
func Example_serviceInfo() {
	svc := &microservice.Info{
		ID:           "example-service",
		Name:         "example-service",
		Type:         microservice.ServiceTypeAPI,
		Version:      "1.0.0",
		Status:       microservice.StatusRunning,
		Capabilities: []string{"feature-a", "feature-b"},
		Tags:         []string{"production", "tier-1"},
		Metadata: map[string]string{
			"team":        "platform",
			"cost-center": "engineering",
		},
	}

	// Add some instances
	svc.Instances = []*microservice.Instance{
		{ID: "inst-1", Status: microservice.StatusRunning},
		{ID: "inst-2", Status: microservice.StatusRunning},
		{ID: "inst-3", Status: microservice.StatusDegraded},
	}

	// Check capabilities
	fmt.Printf("Has feature-a: %v\n", svc.HasCapability("feature-a"))
	fmt.Printf("Has feature-c: %v\n", svc.HasCapability("feature-c"))

	// Check tags
	fmt.Printf("Is production: %v\n", svc.HasTag("production"))

	// Get metadata
	fmt.Printf("Team: %s\n", svc.GetMetadata("team"))

	// Instance counts
	fmt.Printf("Healthy instances: %d\n", svc.HealthyInstanceCount())
	fmt.Printf("Available instances: %d\n", svc.AvailableInstanceCount())

	// Output:
	// Has feature-a: true
	// Has feature-c: false
	// Is production: true
	// Team: platform
	// Healthy instances: 2
	// Available instances: 3
}

// Example_instanceInfo demonstrates instance info utilities.
func Example_instanceInfo() {
	inst := &microservice.Instance{
		ID:        "example-instance",
		ServiceID: "example-service",
		Host:      "10.0.1.100",
		Port:      8080,
		Protocol:  "http",
		Status:    microservice.StatusRunning,
		Region:    "us-east-1",
		Zone:      "us-east-1a",
		PodName:   "example-service-7d4f9b8c-x2k9l",
		Metadata: map[string]string{
			"node": "worker-1",
		},
	}

	// Get address
	fmt.Printf("Address: %s\n", inst.Address())

	// Check health
	fmt.Printf("Is healthy: %v\n", inst.IsHealthy())

	// Get metadata
	fmt.Printf("Node: %s\n", inst.GetMetadata("node"))

	// Check expiration
	fmt.Printf("Is expired: %v\n", inst.IsExpired())

	// Output:
	// Address: 10.0.1.100:8080
	// Is healthy: true
	// Node: worker-1
	// Is expired: false
}
