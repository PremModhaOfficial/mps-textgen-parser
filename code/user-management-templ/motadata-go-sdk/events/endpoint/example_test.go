package endpoint_test

import (
	"context"
	"fmt"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/endpoint"
)

const (
	userServiceID  = "user-service"
	callbackTestID = "callback-test"
	myServiceID    = "my-service"
)

// Example_basicRegistration demonstrates basic endpoint registration.
func Example_basicRegistration() {
	// Create a registry
	registry := endpoint.NewRegistry()

	// Create an endpoint
	ep := &endpoint.Info{
		ID:          "user-api-v1",
		Name:        "User API",
		ServiceID:   userServiceID,
		ServiceName: "User Service",
		Host:        "localhost",
		Port:        8080,
		Path:        "/api/v1/users",
		Protocol:    endpoint.ProtocolHTTP,
		Status:      endpoint.StatusHealthy,
		Version:     "1.0.0",
	}

	// Register the endpoint
	if err := registry.Register(ep); err != nil {
		fmt.Printf("Failed to register: %v\n", err)
		return
	}

	// Retrieve the endpoint
	retrieved, found := registry.Get("user-api-v1")
	if found {
		fmt.Printf("Endpoint: %s at %s\n", retrieved.Name, retrieved.URL())
	}

	// Output:
	// Endpoint: User API at http://localhost:8080/api/v1/users
}

// Example_registrationWithTTL demonstrates endpoint registration with automatic expiration.
func Example_registrationWithTTL() {
	registry := endpoint.NewRegistry()

	ep := &endpoint.Info{
		ID:        "temp-endpoint",
		ServiceID: "temp-service",
		Host:      "localhost",
		Port:      9090,
		Path:      "/",
		Protocol:  endpoint.ProtocolHTTP,
		Status:    endpoint.StatusHealthy,
	}

	// Register with 5-minute TTL
	if err := registry.RegisterWithTTL(ep, 5*time.Minute); err != nil {
		fmt.Printf("Failed to register: %v\n", err)
		return
	}

	retrieved, _ := registry.Get("temp-endpoint")
	fmt.Printf("Expires at: %v\n", !retrieved.ExpiresAt.IsZero())

	// Output:
	// Expires at: true
}

// Example_queryEndpoints demonstrates querying endpoints with filters.
func Example_queryEndpoints() {
	registry := endpoint.NewRegistry()

	// Register multiple endpoints
	endpoints := []*endpoint.Info{
		{
			ID:        "user-api-1",
			ServiceID: userServiceID,
			Host:      "host1",
			Port:      8080,
			Path:      "/api/users",
			Protocol:  endpoint.ProtocolHTTP,
			Status:    endpoint.StatusHealthy,
			Tags:      []string{"api", "users"},
			Version:   "1.0.0",
		},
		{
			ID:        "user-api-2",
			ServiceID: userServiceID,
			Host:      "host2",
			Port:      8080,
			Path:      "/api/users",
			Protocol:  endpoint.ProtocolHTTP,
			Status:    endpoint.StatusHealthy,
			Tags:      []string{"api", "users"},
			Version:   "2.0.0",
		},
		{
			ID:        "order-api-1",
			ServiceID: "order-service",
			Host:      "host3",
			Port:      8080,
			Path:      "/api/orders",
			Protocol:  endpoint.ProtocolHTTP,
			Status:    endpoint.StatusUnhealthy,
			Tags:      []string{"api", "orders"},
			Version:   "1.0.0",
		},
	}

	for _, ep := range endpoints {
		registry.Register(ep)
	}

	// Query healthy endpoints for user-service
	results := registry.Query(&endpoint.Query{
		ServiceIDs:  []string{userServiceID},
		OnlyHealthy: true,
	})
	fmt.Printf("Healthy user-service endpoints: %d\n", len(results))

	// Query by tag
	results = registry.Query(&endpoint.Query{
		Tags:          []string{"api"},
		OnlyAvailable: true,
	})
	fmt.Printf("Available API endpoints: %d\n", len(results))

	// Query by version
	results = registry.Query(&endpoint.Query{
		Version: "1.0.0",
	})
	fmt.Printf("v1.0.0 endpoints: %d\n", len(results))

	// Output:
	// Healthy user-service endpoints: 2
	// Available API endpoints: 2
	// v1.0.0 endpoints: 2
}

// Example_callbacks demonstrates using registry callbacks.
// Note: Callbacks run in goroutines, so output order may vary.
func Example_callbacks() {
	registry := endpoint.NewRegistry()

	registered := make(chan string, 1)
	deregistered := make(chan string, 1)
	statusChanged := make(chan string, 1)

	// Set up callbacks
	registry.OnRegistered(func(ep *endpoint.Info) {
		registered <- ep.ID
	})

	registry.OnDeregistered(func(ep *endpoint.Info, reason string) {
		deregistered <- ep.ID
	})

	registry.OnStatusChange(func(ep *endpoint.Info, oldStatus, newStatus endpoint.Status) {
		statusChanged <- fmt.Sprintf("%s->%s", oldStatus, newStatus)
	})

	// Register an endpoint
	ep := &endpoint.Info{
		ID:        callbackTestID,
		ServiceID: "test-service",
		Host:      "localhost",
		Port:      8080,
		Protocol:  endpoint.ProtocolHTTP,
		Status:    endpoint.StatusHealthy,
	}
	registry.Register(ep)

	// Wait for registration callback
	fmt.Printf("Registered: %s\n", <-registered)

	// Update status
	registry.UpdateStatus(callbackTestID, endpoint.StatusDegraded)

	// Wait for status change callback
	fmt.Printf("Status change: %s\n", <-statusChanged)

	// Deregister
	registry.DeregisterWithReason(callbackTestID, "test complete")

	// Wait for deregistration callback
	fmt.Printf("Deregistered: %s\n", <-deregistered)

	// Output:
	// Registered: callback-test
	// Status change: healthy->degraded
	// Deregistered: callback-test
}

// Example_manager demonstrates using the endpoint manager.
func Example_manager() {
	// Create manager with health checking disabled for this example
	manager, err := endpoint.NewManager(endpoint.ManagerConfig{
		HealthCheck: endpoint.HealthCheckConfig{
			Enabled: false,
		},
	})
	if err != nil {
		fmt.Printf("Failed to create manager: %v\n", err)
		return
	}

	ctx := context.Background()

	// Register endpoints
	ep1 := &endpoint.Info{
		ID:        "api-1",
		ServiceID: myServiceID,
		Host:      "host1",
		Port:      8080,
		Path:      "/api",
		Protocol:  endpoint.ProtocolHTTP,
		Status:    endpoint.StatusHealthy,
		Tags:      []string{"primary"},
	}
	ep2 := &endpoint.Info{
		ID:        "api-2",
		ServiceID: myServiceID,
		Host:      "host2",
		Port:      8080,
		Path:      "/api",
		Protocol:  endpoint.ProtocolHTTP,
		Status:    endpoint.StatusHealthy,
		Tags:      []string{"secondary"},
	}

	manager.RegisterEndpoint(ctx, ep1)
	manager.RegisterEndpoint(ctx, ep2)

	// Discover endpoints with options
	endpoints := manager.Discover(myServiceID,
		endpoint.OnlyHealthy(),
	)
	fmt.Printf("Discovered %d endpoints\n", len(endpoints))

	// Get all services
	services := manager.Services()
	fmt.Printf("Services: %v\n", services)

	// Shutdown
	manager.Shutdown(ctx)

	// Output:
	// Discovered 2 endpoints
	// Services: [my-service]
}

// Example_multiTenant demonstrates multi-tenant endpoint management.
func Example_multiTenant() {
	registry := endpoint.NewRegistry()

	// Register endpoints for different tenants
	registry.Register(&endpoint.Info{
		ID:        "tenant1-api",
		ServiceID: "api-service",
		Host:      "tenant1.example.com",
		Port:      443,
		Protocol:  endpoint.ProtocolHTTPS,
		Status:    endpoint.StatusHealthy,
		TenantID:  "tenant-1",
		Region:    "us-east-1",
	})

	registry.Register(&endpoint.Info{
		ID:        "tenant2-api",
		ServiceID: "api-service",
		Host:      "tenant2.example.com",
		Port:      443,
		Protocol:  endpoint.ProtocolHTTPS,
		Status:    endpoint.StatusHealthy,
		TenantID:  "tenant-2",
		Region:    "eu-west-1",
	})

	// Query by tenant
	tenant1Endpoints := registry.Query(&endpoint.Query{
		TenantID: "tenant-1",
	})
	fmt.Printf("Tenant 1 endpoints: %d\n", len(tenant1Endpoints))

	// Query by region
	usEndpoints := registry.Query(&endpoint.Query{
		Region: "us-east-1",
	})
	fmt.Printf("US East endpoints: %d\n", len(usEndpoints))

	// Output:
	// Tenant 1 endpoints: 1
	// US East endpoints: 1
}

// Example_endpointInfo demonstrates endpoint info utilities.
func Example_endpointInfo() {
	ep := &endpoint.Info{
		ID:        "example-api",
		ServiceID: "example-service",
		Host:      "api.example.com",
		Port:      443,
		Path:      "/v1",
		Protocol:  endpoint.ProtocolHTTPS,
		Status:    endpoint.StatusHealthy,
		Tags:      []string{"production", "api"},
		Metadata: map[string]string{
			"environment": "prod",
			"version":     "1.2.3",
		},
	}

	// URL and Address
	fmt.Printf("URL: %s\n", ep.URL())
	fmt.Printf("Address: %s\n", ep.Address())

	// Check tags
	fmt.Printf("Has 'api' tag: %v\n", ep.HasTag("api"))
	fmt.Printf("Has 'staging' tag: %v\n", ep.HasTag("staging"))

	// Get metadata
	fmt.Printf("Environment: %s\n", ep.GetMetadata("environment"))

	// Check health
	fmt.Printf("Is healthy: %v\n", ep.IsHealthy())

	// Output:
	// URL: https://api.example.com/v1
	// Address: api.example.com:443
	// Has 'api' tag: true
	// Has 'staging' tag: false
	// Environment: prod
	// Is healthy: true
}
