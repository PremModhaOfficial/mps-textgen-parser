// Package nats provides the NATS transport implementation for the events system.
//
// This package implements the transport interfaces for NATS messaging:
//   - Connection management with automatic reconnection
//   - JetStream support for persistent messaging
//   - Publishing with acknowledgment handling
//   - Subscribing with queue groups
//   - Health monitoring
//
// # Basic Usage
//
// Create a NATS transport and connect:
//
//	import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/transport/nats"
//
//	// Create transport
//	transport := nats.NewTransport()
//
//	// Configure
//	cfg := config.DefaultConfig()
//	cfg.Servers = []string{"nats://localhost:4222"}
//
//	// Create credentials
//	creds := auth.JWT(jwtToken, nkeySeed)
//
//	// Connect
//	conn, err := transport.Connect(ctx, cfg, creds)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer conn.Close(ctx)
//
// # Publisher
//
// Create a publisher for sending messages:
//
//	publisher, err := transport.Publisher(conn)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Publish a message
//	msg := core.NewMessage([]byte(`{"event": "user.created"}`))
//	err = publisher.Publish(ctx, "events.users", msg)
//
//	// Publish asynchronously with future
//	future := publisher.PublishAsync(ctx, "events.users", msg)
//	ack, err := future.Wait(ctx)
//
//	// Request-Reply pattern
//	reply, err := publisher.Request(ctx, "api.users.get", msg)
//
// # Subscriber
//
// Create a subscriber for receiving messages:
//
//	subscriber, err := transport.Subscriber(conn)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Subscribe to a subject
//	sub, err := subscriber.Subscribe(ctx, "events.users", func(ctx context.Context, msg *core.Message) error {
//	    log.Printf("Received: %s", msg.Data)
//	    return nil
//	})
//
//	// Queue subscribe for load balancing
//	sub, err = subscriber.QueueSubscribe(ctx, "events.users", "worker-group", handler)
//
//	// Unsubscribe
//	sub.Unsubscribe()
//
// # Connection Factory
//
// Use the factory with tenant manager for multi-tenant support:
//
//	factory := nats.NewFactory()
//
//	manager, err := tenant.NewManager(tenant.ManagerConfig{
//	    Config:  cfg,
//	    Factory: factory,
//	})
//
// # Tenant Connection
//
// Create tenant-aware connections directly:
//
//	tc := nats.NewTenantConnection("tenant-123", cfg, creds)
//	if err := tc.Connect(ctx); err != nil {
//	    log.Fatal(err)
//	}
//	defer tc.Close(ctx)
//
//	// Get publisher and subscriber
//	pub := tc.Publisher()
//	sub := tc.Subscriber()
//
// # JetStream Support
//
// Configure JetStream for persistent messaging:
//
//	cfg := config.DefaultConfig()
//	cfg.Stream = config.StreamConfig{
//	    Enabled: true,
//	    Domain:  "hub",  // For leaf nodes
//	    Prefix:  "js",   // API prefix
//	}
//
// JetStream features:
//   - Persistent message storage
//   - At-least-once delivery guarantees
//   - Message replay from any point
//   - Consumer groups with various delivery modes
//
// # Authentication
//
// Supported authentication methods:
//
//	// JWT authentication (recommended for multi-tenant)
//	creds := auth.JWT(jwtToken, nkeySeed)
//
//	// Username/password
//	creds := auth.UserPass("user", "password")
//
//	// Token
//	creds := auth.Token("secret-token")
//
//	// NKey
//	creds := auth.NKey("/path/to/nkey.nk")
//
//	// Credentials file
//	creds := auth.CredentialsFileAuth("/path/to/file.creds")
//
// # TLS Configuration
//
// Configure TLS for secure connections:
//
//	cfg := config.DefaultConfig()
//	cfg.Servers = []string{"tls://nats.example.com:4222"}
//	cfg.TLS = &config.TLSConfig{
//	    Enabled:  true,
//	    CertFile: "/path/to/client-cert.pem",
//	    KeyFile:  "/path/to/client-key.pem",
//	    CAFile:   "/path/to/ca-cert.pem",
//	}
//
// # Reconnection
//
// Configure automatic reconnection:
//
//	cfg := config.DefaultConfig()
//	cfg.Reconnect = config.ReconnectConfig{
//	    MaxAttempts:     -1,  // Infinite
//	    InitialInterval: 100 * time.Millisecond,
//	    MaxInterval:     30 * time.Second,
//	    Multiplier:      2.0,
//	    Jitter:          0.1,
//	}
//
// # Health Monitoring
//
// Check connection health:
//
//	health := conn.Health()
//	if !health.Healthy {
//	    log.Printf("Connection unhealthy: %s", health.Message)
//	}
//
//	// Check connection state
//	state := conn.State()
//	switch state {
//	case core.StateConnected:
//	    log.Println("Connected")
//	case core.StateReconnecting:
//	    log.Println("Reconnecting...")
//	case core.StateDisconnected:
//	    log.Println("Disconnected")
//	}
//
// # Subject Patterns
//
// NATS supports hierarchical subjects with wildcards:
//
//	// Exact match
//	"events.users.created"
//
//	// Single-level wildcard (*)
//	"events.users.*"  // Matches events.users.created, events.users.deleted
//
//	// Multi-level wildcard (>)
//	"events.>"  // Matches all events.* subjects
//
// # Error Handling
//
// Common errors:
//
//	core.ErrNotConnected      - Connection not established
//	core.ErrConnectionClosed  - Connection was closed
//	core.ErrPublishFailed     - Publish operation failed
//	core.ErrAuthFailed        - Authentication failed
//	core.ErrPermissionDenied  - Insufficient permissions
//
// # Registration
//
// Register NATS transport with the global registry:
//
//	nats.Register()
//
//	// Now accessible via transport registry
//	t, _ := transport.Get("nats")
//
// # Thread Safety
//
// All components (Connection, Publisher, Subscriber) are safe for concurrent use.
package nats
