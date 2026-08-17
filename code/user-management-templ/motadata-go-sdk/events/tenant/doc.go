// Package tenant provides multi-tenant connection management for the events system.
//
// This package enables managing multiple isolated tenant connections:
//   - Connection lifecycle management per tenant
//   - Automatic reconnection handling
//   - Idle connection cleanup
//   - Health monitoring
//   - Credential management integration
//
// # Architecture
//
// The package provides a Manager that handles tenant connections:
//
//	┌──────────────────────────────────────────┐
//	│              Tenant Manager              │
//	├──────────────────────────────────────────┤
//	│  Registry (tenant connections)           │
//	│  Transport Factory (NATS, etc.)          │
//	│  Credential Manager (JWT handling)       │
//	│  Health Checker (connection monitoring)  │
//	│  Cleanup Loop (idle connection removal)  │
//	└──────────────────────────────────────────┘
//
// # Basic Usage
//
// Create and use a tenant manager:
//
//	// Create manager
//	manager, err := tenant.NewManager(tenant.ManagerConfig{
//	    Config:      config.DefaultConfig(),
//	    Factory:     nats.NewFactory(),
//	    CredManager: credentialManager,
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer manager.Shutdown(context.Background())
//
//	// Connect a tenant
//	creds := auth.JWT(jwtToken, nkeySeed)
//	conn, err := manager.Connect(ctx, "tenant-123", creds)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Use the connection
//	publisher := conn.Publisher()
//	subscriber := conn.Subscriber()
//
// # Connection Management
//
// Manage tenant connections:
//
//	// Get existing connection
//	conn, exists := manager.Get("tenant-123")
//
//	// Check if tenant is connected
//	if manager.IsConnected("tenant-123") {
//	    // Use connection
//	}
//
//	// Disconnect a tenant
//	err := manager.Disconnect(ctx, "tenant-123")
//
//	// List all connected tenants
//	tenantIDs := manager.List()
//
// # Registration Payload
//
// Connect using a registration payload:
//
//	payload := auth.RegistrationPayload{
//	    TenantID: "tenant-123",
//	    JWT:      jwtToken,
//	    Seed:     nkeySeed,
//	}
//	conn, err := manager.ConnectWithPayload(ctx, payload)
//
// # Tenant Registry
//
// The Registry stores tenant connection state:
//
//	registry := tenant.NewRegistry()
//
//	// Store tenant info
//	info := &tenant.Info{
//	    TenantID:  "tenant-123",
//	    Status:    tenant.StatusConnected,
//	    CreatedAt: time.Now(),
//	}
//	registry.Store(info)
//
//	// Get tenant info
//	info, exists := registry.Get("tenant-123")
//
//	// Set callbacks
//	registry.OnConnect(func(tenantID string) {
//	    log.Printf("Tenant %s connected", tenantID)
//	})
//	registry.OnDisconnect(func(tenantID string) {
//	    log.Printf("Tenant %s disconnected", tenantID)
//	})
//
// # Tenant Connection
//
// The Connection wrapper provides tenant-specific operations:
//
//	type Connection interface {
//	    TenantID() string
//	    Publisher() core.Publisher
//	    Subscriber() core.Subscriber
//	    IsConnected() bool
//	    Health() core.HealthStatus
//	    LastActivity() time.Time
//	    Reconnect(ctx context.Context) error
//	    Close(ctx context.Context) error
//	}
//
// # Health Monitoring
//
// Monitor tenant connection health:
//
//	// Get health status
//	health := conn.Health()
//	if !health.Healthy {
//	    log.Printf("Tenant unhealthy: %s", health.Message)
//	}
//
//	// Check all tenants
//	for _, tenantID := range manager.List() {
//	    conn, _ := manager.Get(tenantID)
//	    health := conn.Health()
//	    log.Printf("Tenant %s: %s", tenantID, health.State)
//	}
//
// # Idle Connection Cleanup
//
// The manager automatically cleans up idle connections:
//
//	cfg := tenant.ManagerConfig{
//	    Config: &config.Config{
//	        IdleTimeout:     30 * time.Minute,  // Close after 30min idle
//	        CleanupInterval: 1 * time.Minute,   // Check every minute
//	    },
//	}
//
// # Graceful Shutdown
//
// Properly shutdown all tenant connections:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
//	defer cancel()
//
//	if err := manager.Shutdown(ctx); err != nil {
//	    log.Printf("Shutdown error: %v", err)
//	}
//
// # Error Handling
//
// Common tenant errors (defined in the utils package):
//
//	utils.ErrTenantNotFound     - Tenant not registered
//	utils.ErrTenantExists       - Tenant already connected
//	utils.ErrTenantDisconnected - Tenant connection lost
//
// Example:
//
//	conn, err := manager.Get("tenant-123")
//	if errors.Is(err, utils.ErrTenantNotFound) {
//	    // Tenant needs to connect first
//	    conn, err = manager.Connect(ctx, "tenant-123", creds)
//	}
package tenant
