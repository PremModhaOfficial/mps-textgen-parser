// Package microservice provides microservice registration, discovery, and lifecycle
// management for distributed systems.
//
// This package enables services to register themselves with a central registry,
// manage multiple instances, perform health checking, and discover other services
// for inter-service communication.
//
// # Key Components
//
// ## Info
//
// Info holds comprehensive microservice information including metadata, capabilities,
// and instance information:
//
//	svc := &microservice.Info{
//	    ID:           "user-service",
//	    Name:         "user-service",
//	    DisplayName:  "User Service",
//	    Description:  "Handles user authentication and management",
//	    Type:         microservice.ServiceTypeAPI,
//	    Version:      "2.1.0",
//	    APIVersion:   "v2",
//	    Dependencies: []string{"database-service", "cache-service"},
//	    Capabilities: []string{"auth", "user-management", "oauth"},
//	    Tags:         []string{"core", "auth"},
//	    Status:       microservice.StatusRunning,
//	}
//
// ## Instance
//
// Instance represents a running instance of a microservice:
//
//	inst := &microservice.Instance{
//	    ID:              "user-service-instance-1",
//	    ServiceID:       "user-service",
//	    Host:            "10.0.1.50",
//	    Port:            8080,
//	    Status:          microservice.StatusRunning,
//	    HealthCheckPath: "/health",
//	    Weight:          100,
//	    Region:          "us-east-1",
//	    Zone:            "us-east-1a",
//	    PodName:         "user-service-7d4f9b8c-x2k9l",
//	}
//
// ## Registry
//
// Registry provides thread-safe storage for services and instances:
//
//	// Create registry
//	registry := microservice.NewRegistry()
//
//	// Register a service
//	registry.RegisterService(svc)
//
//	// Register an instance
//	registry.RegisterInstance(inst)
//
//	// Register instance with TTL
//	registry.RegisterInstanceWithTTL(inst, 30 * time.Second)
//
//	// Query services
//	services := registry.QueryServices(&microservice.Query{
//	    Types:        []microservice.ServiceType{microservice.ServiceTypeAPI},
//	    OnlyHealthy:  true,
//	    HasInstances: true,
//	})
//
//	// Get healthy instances
//	instances := registry.GetHealthyInstances("user-service")
//
// ## Manager
//
// Manager handles service lifecycle including health checking and event publishing:
//
//	manager, err := microservice.NewManager(microservice.ManagerConfig{
//	    Registry:  registry,
//	    Publisher: publisher,
//	    HealthCheck: microservice.HealthCheckConfig{
//	        Enabled:            true,
//	        Interval:           30 * time.Second,
//	        Timeout:            5 * time.Second,
//	        HealthyThreshold:   2,
//	        UnhealthyThreshold: 3,
//	    },
//	})
//
//	// Register service
//	manager.RegisterServiceDirect(ctx, svc)
//
//	// Register instance
//	manager.RegisterInstanceDirect(ctx, inst)
//
//	// Discover services
//	services := manager.Discover("user-service",
//	    microservice.WithType(microservice.ServiceTypeAPI),
//	    microservice.WithCapabilities("auth"),
//	    microservice.OnlyHealthy(),
//	)
//
//	// Get instances for load balancing
//	instances := manager.DiscoverInstances("user-service", true)
//
// ## Handler
//
// Handler provides NATS message handlers for distributed operations:
//
//	handler := microservice.NewHandler(manager)
//
//	// Register subscriptions
//	subs, err := handler.RegisterSubscriptions(ctx, subscriber)
//
//	// Or with queue groups for load balancing
//	subs, err := handler.RegisterQueueSubscriptions(ctx, subscriber, "service-registry")
//
// ## Client
//
// Client provides remote service operations over messaging:
//
//	client := microservice.NewClient(publisher)
//
//	// Register service remotely
//	resp, err := client.RegisterService(ctx, svc)
//
//	// Register instance with TTL
//	resp, err := client.RegisterInstance(ctx, inst, 30*time.Second)
//
//	// Discover services
//	resp, err := client.Discover(ctx, "user-service",
//	    microservice.WithCapabilities("auth"),
//	)
//
//	// Heartbeat/renew instance TTL
//	resp, err := client.RenewInstance(ctx, "instance-id", 30*time.Second)
//
// # Service Types
//
// The package defines common service types:
//
//	microservice.ServiceTypeAPI       // REST/gRPC API services
//	microservice.ServiceTypeWorker    // Background workers
//	microservice.ServiceTypeGateway   // API gateways
//	microservice.ServiceTypeBroker    // Message brokers
//	microservice.ServiceTypeDatabase  // Database services
//	microservice.ServiceTypeCache     // Cache services (Redis, Memcached)
//	microservice.ServiceTypeQueue     // Queue services
//	microservice.ServiceTypeScheduler // Job schedulers
//	microservice.ServiceTypeMonitor   // Monitoring services
//	microservice.ServiceTypeCustom    // Custom service types
//
// # Service Status
//
// Services and instances can have the following statuses:
//
//	microservice.StatusUnknown     // Status not determined
//	microservice.StatusStarting    // Service is starting
//	microservice.StatusRunning     // Running normally
//	microservice.StatusDegraded    // Running with reduced capacity
//	microservice.StatusStopping    // Shutting down
//	microservice.StatusStopped     // Stopped
//	microservice.StatusFailed      // Failed
//	microservice.StatusMaintenance // Under maintenance
//
// # Health Checking
//
// The Manager performs automatic HTTP health checks on instances:
//
//	inst := &microservice.Instance{
//	    // ... basic info ...
//	    HealthCheckPath: "/health",
//	    Protocol:        "http",
//	}
//
// Health check behavior:
//   - HTTP 2xx: Instance is healthy
//   - HTTP 503: Instance is in maintenance
//   - Other codes/errors: Instance is unhealthy
//   - Thresholds determine when status actually changes
//
// # Instance TTL and Heartbeats
//
// Instances can be registered with TTL for automatic cleanup:
//
//	// Register with 30-second TTL
//	registry.RegisterInstanceWithTTL(inst, 30 * time.Second)
//
//	// Heartbeat loop to keep instance alive
//	go func() {
//	    ticker := time.NewTicker(10 * time.Second)
//	    for range ticker.C {
//	        registry.RenewInstance(inst.ID, 30 * time.Second)
//	    }
//	}()
//
// # Event Publishing
//
// When a Publisher is configured, lifecycle events are published:
//
// Service Events:
//   - services.registered - Service registered
//   - services.deregistered - Service deregistered
//   - services.updated - Service updated
//   - services.health - Service health changed
//
// Instance Events:
//   - instances.registered - Instance registered
//   - instances.deregistered - Instance deregistered
//   - instances.health - Instance health changed
//
// # Capabilities and Dependencies
//
// Services can declare capabilities and dependencies:
//
//	svc := &microservice.Info{
//	    ID:           "order-service",
//	    Name:         "order-service",
//	    Capabilities: []string{"orders", "checkout", "payments"},
//	    Dependencies: []string{"user-service", "inventory-service"},
//	}
//
//	// Find services with specific capabilities
//	services := registry.QueryServices(&microservice.Query{
//	    Capabilities: []string{"payments"},
//	})
//
// # Load Balancing
//
// Instances support weight-based load balancing:
//
//	inst1 := &microservice.Instance{
//	    ID:        "inst-1",
//	    ServiceID: "my-service",
//	    Weight:    100, // High priority
//	}
//	inst2 := &microservice.Instance{
//	    ID:        "inst-2",
//	    ServiceID: "my-service",
//	    Weight:    50, // Lower priority
//	}
//
// # Multi-Tenancy
//
// Services support multi-tenant environments:
//
//	svc := &microservice.Info{
//	    ID:       "tenant-api",
//	    Name:     "tenant-api",
//	    TenantID: "tenant-123",
//	    Owner:    "team-platform",
//	}
//
//	// Query by tenant
//	services := registry.QueryServices(&microservice.Query{
//	    TenantID: "tenant-123",
//	})
//
// # Callbacks
//
// Registry supports lifecycle callbacks:
//
//	registry.OnServiceRegistered(func(svc *microservice.Info) {
//	    log.Printf("Service registered: %s", svc.Name)
//	})
//
//	registry.OnInstanceRegistered(func(inst *microservice.Instance) {
//	    log.Printf("Instance registered: %s@%s", inst.ID, inst.ServiceID)
//	})
//
//	registry.OnInstanceStatusChange(func(inst *microservice.Instance, old, new microservice.Status) {
//	    log.Printf("Instance %s: %s -> %s", inst.ID, old, new)
//	})
//
// # HTTP Handlers
//
// Manager provides HTTP handlers for REST API integration:
//
//	http.HandleFunc("/services/register", manager.HandleServiceRegistration)
//	http.HandleFunc("/services/deregister", manager.HandleServiceDeregistration)
//	http.HandleFunc("/services/query", manager.HandleServiceQuery)
//	http.HandleFunc("/instances/register", manager.HandleInstanceRegistration)
//	http.HandleFunc("/instances/renew", manager.HandleInstanceRenew)
//
// # Thread Safety
//
// All components are safe for concurrent use from multiple goroutines.
package microservice
