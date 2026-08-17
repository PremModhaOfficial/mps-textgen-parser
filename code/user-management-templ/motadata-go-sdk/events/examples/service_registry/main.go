// Package main demonstrates a complete service registry example using the
// endpoint and microservice packages.
//
// This example shows how to:
//   - Set up a microservice registry
//   - Register services and instances
//   - Register API endpoints
//   - Perform service discovery
//   - Handle health checking
//   - Use callbacks for monitoring
//
// Run this example:
//
//	go run main.go
package main

import (
	"context"
	"log"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/endpoint"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/microservice"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"
)

const (
	svcIDUser         = "user-service"
	svcIDOrder        = "order-service"
	svcIDNotification = "notification-worker"
	svcIDGateway      = "api-gateway"

	svcNameUser    = "User Service"
	svcNameOrder   = "Order Service"
	svcNameGateway = "API Gateway"

	healthPath = "/health"
	regionEast = "us-east-1"
	apiHost    = "api.example.com"

	msgRegistryState = "Registry state"
)

func main() {
	ctx := context.Background()
	appLog := logger.L()

	appLog.Info(ctx, "=== Service Registry Example ===")

	// Step 1: Create registries
	appLog.Info(ctx, "1. Creating registries...")
	svcRegistry := microservice.NewRegistry()
	epRegistry := endpoint.NewRegistry()

	// Step 2: Set up callbacks for monitoring
	appLog.Info(ctx, "2. Setting up monitoring callbacks...")
	setupCallbacks(svcRegistry, epRegistry)

	// Step 3: Register microservices
	appLog.Info(ctx, "3. Registering microservices...")
	registerServices(svcRegistry)

	// Step 4: Register service instances
	appLog.Info(ctx, "4. Registering service instances...")
	registerInstances(svcRegistry)

	// Step 5: Register API endpoints
	appLog.Info(ctx, "5. Registering API endpoints...")
	registerEndpoints(epRegistry)

	// Step 6: Demonstrate service discovery
	appLog.Info(ctx, "6. Service Discovery Demo...")
	demonstrateDiscovery(ctx, appLog, svcRegistry, epRegistry)

	// Step 7: Create managers for health checking
	appLog.Info(ctx, "7. Creating managers with health checking...")
	svcManager, epManager := createManagers(svcRegistry, epRegistry)
	defer svcManager.Shutdown(ctx)
	defer epManager.Shutdown(ctx)

	// Step 8: Simulate status changes
	appLog.Info(ctx, "8. Simulating status changes...")
	simulateStatusChanges(ctx, appLog, svcRegistry, epRegistry)

	// Step 9: Final state
	appLog.Info(ctx, "9. Final Registry State:")
	printRegistryState(ctx, appLog, svcRegistry, epRegistry)

	appLog.Info(ctx, "=== Example Complete ===")
}

func setupCallbacks(svcReg *microservice.Registry, epReg *endpoint.Registry) {
	// Microservice callbacks
	svcReg.OnServiceRegistered(func(svc *microservice.Info) {
		log.Printf("[SVC] Registered: %s (type: %s)", svc.Name, svc.Type)
	})

	svcReg.OnInstanceRegistered(func(inst *microservice.Instance) {
		log.Printf("[INST] Registered: %s @ %s", inst.ID, inst.Address())
	})

	svcReg.OnInstanceStatusChange(func(inst *microservice.Instance, old, new microservice.Status) {
		log.Printf("[INST] Status change: %s: %s -> %s", inst.ID, old, new)
	})

	// Endpoint callbacks
	epReg.OnRegistered(func(ep *endpoint.Info) {
		log.Printf("[EP] Registered: %s -> %s", ep.ID, ep.URL())
	})

	epReg.OnStatusChange(func(ep *endpoint.Info, old, new endpoint.Status) {
		log.Printf("[EP] Status change: %s: %s -> %s", ep.ID, old, new)
	})

	// Allow callbacks to be set up
	time.Sleep(10 * time.Millisecond)
}

func registerServices(registry *microservice.Registry) {
	services := []*microservice.Info{
		{
			ID:           svcIDUser,
			Name:         svcIDUser,
			DisplayName:  svcNameUser,
			Description:  "Handles user authentication and profile management",
			Type:         microservice.ServiceTypeAPI,
			Version:      "2.1.0",
			APIVersion:   "v2",
			Status:       microservice.StatusRunning,
			Capabilities: []string{"auth", "users", "profiles", "oauth"},
			Tags:         []string{"core", "auth"},
			Dependencies: []string{"database-service", "cache-service"},
		},
		{
			ID:           svcIDOrder,
			Name:         svcIDOrder,
			DisplayName:  svcNameOrder,
			Description:  "Manages customer orders and checkout",
			Type:         microservice.ServiceTypeAPI,
			Version:      "3.0.0",
			APIVersion:   "v3",
			Status:       microservice.StatusRunning,
			Capabilities: []string{"orders", "checkout", "payments"},
			Tags:         []string{"commerce"},
			Dependencies: []string{svcIDUser, "inventory-service", "payment-gateway"},
		},
		{
			ID:           svcIDNotification,
			Name:         svcIDNotification,
			DisplayName:  "Notification Worker",
			Description:  "Sends emails, SMS, and push notifications",
			Type:         microservice.ServiceTypeWorker,
			Version:      "1.5.0",
			Status:       microservice.StatusRunning,
			Capabilities: []string{"email", "sms", "push"},
			Tags:         []string{"background", "notifications"},
		},
		{
			ID:           svcIDGateway,
			Name:         svcIDGateway,
			DisplayName:  svcNameGateway,
			Description:  "Main API gateway for all services",
			Type:         microservice.ServiceTypeGateway,
			Version:      "4.0.0",
			Status:       microservice.StatusRunning,
			Capabilities: []string{"routing", "rate-limiting", "auth"},
			Tags:         []string{"infrastructure"},
		},
	}

	for _, svc := range services {
		if err := registry.RegisterService(svc); err != nil {
			log.Printf("Failed to register service %s: %v", svc.Name, err)
		}
	}

	// Allow callbacks to execute
	time.Sleep(50 * time.Millisecond)
}

func registerInstances(registry *microservice.Registry) {
	instances := []*microservice.Instance{
		// User service instances
		{
			ID:              "user-svc-1",
			ServiceID:       svcIDUser,
			Host:            "10.0.1.10",
			Port:            8080,
			Protocol:        "http",
			Status:          microservice.StatusRunning,
			HealthCheckPath: healthPath,
			Weight:          100,
			Region:          regionEast,
			Zone:            "us-east-1a",
		},
		{
			ID:              "user-svc-2",
			ServiceID:       svcIDUser,
			Host:            "10.0.1.11",
			Port:            8080,
			Protocol:        "http",
			Status:          microservice.StatusRunning,
			HealthCheckPath: healthPath,
			Weight:          100,
			Region:          regionEast,
			Zone:            "us-east-1b",
		},
		// Order service instances
		{
			ID:              "order-svc-1",
			ServiceID:       svcIDOrder,
			Host:            "10.0.2.10",
			Port:            8080,
			Protocol:        "http",
			Status:          microservice.StatusRunning,
			HealthCheckPath: healthPath,
			Weight:          100,
			Region:          regionEast,
			Zone:            "us-east-1a",
		},
		{
			ID:              "order-svc-2",
			ServiceID:       svcIDOrder,
			Host:            "10.0.2.11",
			Port:            8080,
			Protocol:        "http",
			Status:          microservice.StatusDegraded,
			HealthCheckPath: healthPath,
			Weight:          50,
			Region:          regionEast,
			Zone:            "us-east-1b",
		},
		// Notification worker instances
		{
			ID:              "notif-worker-1",
			ServiceID:       svcIDNotification,
			Host:            "10.0.3.10",
			Port:            9000,
			Status:          microservice.StatusRunning,
			HealthCheckPath: healthPath,
			Weight:          100,
		},
		// API Gateway instances
		{
			ID:              "gateway-1",
			ServiceID:       svcIDGateway,
			Host:            "10.0.0.10",
			Port:            443,
			Protocol:        "https",
			Status:          microservice.StatusRunning,
			HealthCheckPath: healthPath,
			Weight:          100,
		},
		{
			ID:              "gateway-2",
			ServiceID:       svcIDGateway,
			Host:            "10.0.0.11",
			Port:            443,
			Protocol:        "https",
			Status:          microservice.StatusRunning,
			HealthCheckPath: healthPath,
			Weight:          100,
		},
	}

	for _, inst := range instances {
		// Register with 5-minute TTL
		if err := registry.RegisterInstanceWithTTL(inst, 5*time.Minute); err != nil {
			log.Printf("Failed to register instance %s: %v", inst.ID, err)
		}
	}

	// Allow callbacks to execute
	time.Sleep(50 * time.Millisecond)
}

func registerEndpoints(registry *endpoint.Registry) {
	endpoints := []*endpoint.Info{
		// User service endpoints
		{
			ID:              "user-api-auth",
			Name:            "User Authentication",
			ServiceID:       svcIDUser,
			ServiceName:     svcNameUser,
			Host:            apiHost,
			Port:            443,
			Path:            "/api/v2/auth",
			Protocol:        endpoint.ProtocolHTTPS,
			Methods:         []endpoint.Method{endpoint.MethodPOST},
			Status:          endpoint.StatusHealthy,
			Tags:            []string{"auth", "public"},
			Version:         "2.1.0",
			HealthCheckPath: healthPath,
		},
		{
			ID:              "user-api-users",
			Name:            "User Management",
			ServiceID:       svcIDUser,
			ServiceName:     svcNameUser,
			Host:            apiHost,
			Port:            443,
			Path:            "/api/v2/users",
			Protocol:        endpoint.ProtocolHTTPS,
			Methods:         []endpoint.Method{endpoint.MethodGET, endpoint.MethodPOST, endpoint.MethodPUT, endpoint.MethodDELETE},
			Status:          endpoint.StatusHealthy,
			Tags:            []string{"users", "protected"},
			Version:         "2.1.0",
			HealthCheckPath: healthPath,
		},
		// Order service endpoints
		{
			ID:              "order-api-orders",
			Name:            "Order Management",
			ServiceID:       svcIDOrder,
			ServiceName:     svcNameOrder,
			Host:            apiHost,
			Port:            443,
			Path:            "/api/v3/orders",
			Protocol:        endpoint.ProtocolHTTPS,
			Methods:         []endpoint.Method{endpoint.MethodGET, endpoint.MethodPOST},
			Status:          endpoint.StatusHealthy,
			Tags:            []string{"orders", "protected"},
			Version:         "3.0.0",
			HealthCheckPath: healthPath,
		},
		{
			ID:              "order-api-checkout",
			Name:            "Checkout",
			ServiceID:       svcIDOrder,
			ServiceName:     svcNameOrder,
			Host:            apiHost,
			Port:            443,
			Path:            "/api/v3/checkout",
			Protocol:        endpoint.ProtocolHTTPS,
			Methods:         []endpoint.Method{endpoint.MethodPOST},
			Status:          endpoint.StatusHealthy,
			Tags:            []string{"checkout", "protected"},
			Version:         "3.0.0",
			HealthCheckPath: healthPath,
		},
		// Gateway endpoint
		{
			ID:          "gateway-main",
			Name:        svcNameGateway,
			ServiceID:   svcIDGateway,
			ServiceName: svcNameGateway,
			Host:        apiHost,
			Port:        443,
			Path:        "/",
			Protocol:    endpoint.ProtocolHTTPS,
			Methods:     []endpoint.Method{endpoint.MethodAny},
			Status:      endpoint.StatusHealthy,
			Tags:        []string{"gateway", "public"},
			Version:     "4.0.0",
			Priority:    100,
		},
	}

	for _, ep := range endpoints {
		if err := registry.Register(ep); err != nil {
			log.Printf("Failed to register endpoint %s: %v", ep.ID, err)
		}
	}

	// Allow callbacks to execute
	time.Sleep(50 * time.Millisecond)
}

func demonstrateDiscovery(ctx context.Context, appLog *logger.Logger, svcReg *microservice.Registry, epReg *endpoint.Registry) {
	appLog.Info(ctx, "--- Service Discovery ---")

	// Find API services
	apiServices := svcReg.QueryServices(&microservice.Query{
		Types:       []microservice.ServiceType{microservice.ServiceTypeAPI},
		OnlyHealthy: true,
	})
	appLog.Info(ctx, "API Services", logger.Int("count", len(apiServices)))
	for _, svc := range apiServices {
		appLog.Info(ctx, "Service", logger.String("name", svc.Name), logger.String("version", svc.Version), logger.Int("instances", len(svc.Instances)))
	}

	// Find services with auth capability
	authServices := svcReg.QueryServices(&microservice.Query{
		Capabilities: []string{"auth"},
	})
	appLog.Info(ctx, "Services with 'auth' capability", logger.Int("count", len(authServices)))
	for _, svc := range authServices {
		appLog.Info(ctx, "Auth service", logger.String("name", svc.Name))
	}

	// Get healthy instances for user-service
	userInstances := svcReg.GetHealthyInstances(svcIDUser)
	appLog.Info(ctx, "Healthy user-service instances", logger.Int("count", len(userInstances)))
	for _, inst := range userInstances {
		appLog.Info(ctx, "Instance", logger.String("id", inst.ID), logger.String("address", inst.Address()), logger.Int("weight", inst.Weight))
	}

	appLog.Info(ctx, "--- Endpoint Discovery ---")

	// Find protected endpoints
	protectedEndpoints := epReg.Query(&endpoint.Query{
		Tags:        []string{"protected"},
		OnlyHealthy: true,
	})
	appLog.Info(ctx, "Protected endpoints", logger.Int("count", len(protectedEndpoints)))
	for _, ep := range protectedEndpoints {
		appLog.Info(ctx, "Endpoint", logger.String("name", ep.Name), logger.String("url", ep.URL()))
	}

	// Find endpoints by service
	orderEndpoints := epReg.GetByService(svcIDOrder)
	appLog.Info(ctx, "Order service endpoints", logger.Int("count", len(orderEndpoints)))
	for _, ep := range orderEndpoints {
		appLog.Info(ctx, "Endpoint", logger.String("path", ep.Path), logger.Any("methods", ep.Methods))
	}
}

func createManagers(svcReg *microservice.Registry, epReg *endpoint.Registry) (*microservice.Manager, *endpoint.Manager) {
	svcManager, _ := microservice.NewManager(microservice.ManagerConfig{
		Registry: svcReg,
		HealthCheck: microservice.HealthCheckConfig{
			Enabled: false, // Disabled for example (no real servers)
		},
	})

	epManager, _ := endpoint.NewManager(endpoint.ManagerConfig{
		Registry: epReg,
		HealthCheck: endpoint.HealthCheckConfig{
			Enabled: false, // Disabled for example
		},
	})

	return svcManager, epManager
}

func simulateStatusChanges(ctx context.Context, appLog *logger.Logger, svcReg *microservice.Registry, epReg *endpoint.Registry) {
	// Simulate instance becoming unhealthy
	appLog.Info(ctx, "Simulating instance failure...")
	svcReg.UpdateInstanceStatus("order-svc-2", microservice.StatusFailed)

	// Simulate endpoint going to maintenance
	appLog.Info(ctx, "Simulating endpoint maintenance...")
	epReg.UpdateStatus("order-api-checkout", endpoint.StatusMaintenance)

	// Allow callbacks to execute
	time.Sleep(50 * time.Millisecond)
}

func printRegistryState(ctx context.Context, appLog *logger.Logger, svcReg *microservice.Registry, epReg *endpoint.Registry) {
	appLog.Info(ctx, msgRegistryState, logger.Int("services", svcReg.ServiceCount()))
	appLog.Info(ctx, msgRegistryState, logger.Int("instances", svcReg.InstanceCount()))
	appLog.Info(ctx, msgRegistryState, logger.Int("endpoints", epReg.Count()))

	// Print healthy vs total
	allInstances := svcReg.AllInstances()
	healthyCount := 0
	for _, inst := range allInstances {
		if inst.IsHealthy() {
			healthyCount++
		}
	}
	appLog.Info(ctx, "Instance health", logger.Int("healthy", healthyCount), logger.Int("total", len(allInstances)))

	healthyEndpoints := epReg.GetHealthy()
	appLog.Info(ctx, "Endpoint health", logger.Int("healthy", len(healthyEndpoints)), logger.Int("total", epReg.Count()))
}
