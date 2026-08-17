// Package endpoint provides endpoint registration, discovery, and lifecycle management
// for microservice architectures.
//
// This package enables services to register their API endpoints with a central registry,
// discover other service endpoints, and manage endpoint health through automatic health
// checking and status propagation.
//
// # Key Components
//
// ## Info
//
// Info holds comprehensive endpoint information including network details, health status,
// and metadata:
//
//	ep := &endpoint.Info{
//	    ID:          "user-api-v1",
//	    Name:        "User API",
//	    ServiceID:   "user-service",
//	    ServiceName: "User Service",
//	    Host:        "localhost",
//	    Port:        8080,
//	    Path:        "/api/v1/users",
//	    Protocol:    endpoint.ProtocolHTTP,
//	    Methods:     []endpoint.Method{endpoint.MethodGET, endpoint.MethodPOST},
//	    Status:      endpoint.StatusHealthy,
//	    Tags:        []string{"api", "users", "v1"},
//	    Version:     "1.0.0",
//	    HealthCheckPath: "/health",
//	}
//
// ## Registry
//
// Registry provides thread-safe in-memory storage for endpoints with automatic
// expiration and efficient indexing:
//
//	// Create registry with default config
//	registry := endpoint.NewRegistry()
//
//	// Or with custom config
//	registry := endpoint.NewRegistry(endpoint.RegistryConfig{
//	    CleanupInterval:   30 * time.Second,
//	    ExpirationEnabled: true,
//	})
//
//	// Register an endpoint
//	registry.Register(ep)
//
//	// Register with TTL for automatic expiration
//	registry.RegisterWithTTL(ep, 5 * time.Minute)
//
//	// Query endpoints
//	endpoints := registry.Query(&endpoint.Query{
//	    ServiceIDs:    []string{"user-service"},
//	    OnlyHealthy:   true,
//	    Tags:          []string{"api"},
//	})
//
//	// Get by service
//	userEndpoints := registry.GetByService("user-service")
//
//	// Deregister
//	registry.Deregister("user-api-v1")
//
// ## Manager
//
// Manager handles endpoint lifecycle including health checking and event publishing:
//
//	manager, err := endpoint.NewManager(endpoint.ManagerConfig{
//	    Registry:  registry,
//	    Publisher: publisher, // Optional: for publishing events
//	    HealthCheck: endpoint.HealthCheckConfig{
//	        Enabled:            true,
//	        Interval:           30 * time.Second,
//	        Timeout:            5 * time.Second,
//	        HealthyThreshold:   2,
//	        UnhealthyThreshold: 3,
//	    },
//	})
//
//	// Register endpoint
//	err = manager.RegisterEndpoint(ctx, ep)
//
//	// Discover endpoints for a service
//	endpoints := manager.Discover("user-service",
//	    endpoint.WithTags("api"),
//	    endpoint.WithVersion("1.0.0"),
//	    endpoint.OnlyHealthy(),
//	)
//
//	// Graceful shutdown
//	manager.Shutdown(ctx)
//
// ## Handler
//
// Handler provides NATS message handlers for distributed endpoint operations:
//
//	handler := endpoint.NewHandler(manager)
//
//	// Register all subscription handlers
//	subs, err := handler.RegisterSubscriptions(ctx, subscriber)
//
//	// Or use queue subscriptions for load balancing
//	subs, err := handler.RegisterQueueSubscriptions(ctx, subscriber, "endpoint-workers")
//
// ## Client
//
// Client provides a client interface for remote endpoint operations:
//
//	client := endpoint.NewClient(publisher)
//
//	// Register remotely
//	resp, err := client.Register(ctx, ep, 5*time.Minute)
//
//	// Discover remotely
//	resp, err := client.Discover(ctx, "user-service",
//	    endpoint.WithTags("api"),
//	    endpoint.OnlyHealthy(),
//	)
//
// # Health Checking
//
// The Manager automatically performs HTTP health checks on registered endpoints:
//
//	ep := &endpoint.Info{
//	    // ... basic info ...
//	    HealthCheckPath: "/health",     // Path for health check
//	    HealthCheckPort: 8081,          // Optional: different port for health
//	}
//
// Health checks follow these rules:
//   - HTTP 2xx responses mark endpoint as healthy
//   - HTTP 503 marks endpoint as maintenance
//   - Other responses or errors mark endpoint as unhealthy
//   - Consecutive successes/failures determine final status based on thresholds
//
// # Event Publishing
//
// When a Publisher is configured, the Manager publishes lifecycle events:
//
//   - endpoints.registered - When an endpoint is registered
//   - endpoints.deregistered - When an endpoint is deregistered
//   - endpoints.updated - When an endpoint is updated
//   - endpoints.health - When health status changes
//
// # Multi-Tenancy
//
// Endpoints support multi-tenant environments:
//
//	ep := &endpoint.Info{
//	    // ... basic info ...
//	    TenantID: "tenant-123",
//	    Region:   "us-east-1",
//	    Zone:     "us-east-1a",
//	}
//
//	// Query by tenant
//	endpoints := registry.Query(&endpoint.Query{
//	    TenantID: "tenant-123",
//	    Region:   "us-east-1",
//	})
//
// # Callbacks
//
// Registry supports callbacks for lifecycle events:
//
//	registry.OnRegistered(func(ep *endpoint.Info) {
//	    log.Printf("Endpoint registered: %s", ep.ID)
//	})
//
//	registry.OnDeregistered(func(ep *endpoint.Info, reason string) {
//	    log.Printf("Endpoint deregistered: %s, reason: %s", ep.ID, reason)
//	})
//
//	registry.OnStatusChange(func(ep *endpoint.Info, old, new endpoint.Status) {
//	    log.Printf("Endpoint %s status: %s -> %s", ep.ID, old, new)
//	})
//
// # HTTP Handlers
//
// Manager provides HTTP handlers for REST API integration:
//
//	http.HandleFunc("/endpoints/register", manager.HandleRegistration)
//	http.HandleFunc("/endpoints/deregister", manager.HandleDeregistration)
//	http.HandleFunc("/endpoints/query", manager.HandleQuery)
//	http.HandleFunc("/endpoints/renew", manager.HandleRenew)
//
// # Thread Safety
//
// All components (Registry, Manager, Handler, Client) are safe for concurrent use.
package endpoint
