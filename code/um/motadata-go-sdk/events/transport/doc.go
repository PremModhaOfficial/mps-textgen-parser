// Package transport provides the transport layer abstraction for the events system.
//
// This package defines the interfaces and registry for transport implementations:
//   - Transport interface for creating connections
//   - Factory pattern for transport creation
//   - Registry for managing multiple transport types
//
// # Transport Interface
//
// The Transport interface defines how to create connections:
//
//	type Transport interface {
//	    Name() string
//	    Connect(ctx context.Context, cfg *config.Config, creds *auth.Credentials) (core.Connection, error)
//	    Publisher(conn core.Connection) (core.Publisher, error)
//	    Subscriber(conn core.Connection) (core.Subscriber, error)
//	}
//
// # Built-in Transports
//
// The SDK includes the following transport implementations:
//
//	transport/nats - NATS messaging transport with JetStream support
//
// # Using the NATS Transport
//
//	import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/transport/nats"
//
//	// Create transport
//	transport := nats.NewTransport()
//
//	// Connect
//	conn, err := transport.Connect(ctx, cfg, creds)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Get publisher and subscriber
//	publisher, _ := transport.Publisher(conn)
//	subscriber, _ := transport.Subscriber(conn)
//
// # Transport Registry
//
// Register and retrieve transports by name:
//
//	// Get the default registry
//	registry := transport.DefaultRegistry
//
//	// Register a custom transport
//	registry.Register(myTransport)
//
//	// Get a transport by name
//	t, ok := registry.Get("nats")
//	if !ok {
//	    log.Fatal("transport not found")
//	}
//
//	// List registered transports
//	names := registry.Names()
//
// # Creating Custom Transports
//
// Implement the Transport interface for custom messaging systems:
//
//	type KafkaTransport struct {
//	    // ...
//	}
//
//	func (t *KafkaTransport) Name() string {
//	    return "kafka"
//	}
//
//	func (t *KafkaTransport) Connect(ctx context.Context, cfg *config.Config, creds *auth.Credentials) (core.Connection, error) {
//	    // Create Kafka connection
//	}
//
//	func (t *KafkaTransport) Publisher(conn core.Connection) (core.Publisher, error) {
//	    // Create Kafka publisher
//	}
//
//	func (t *KafkaTransport) Subscriber(conn core.Connection) (core.Subscriber, error) {
//	    // Create Kafka subscriber
//	}
//
//	// Register custom transport
//	transport.Register(&KafkaTransport{})
//
// # Factory Pattern
//
// Use factories for creating transports with tenant manager:
//
//	// Create NATS factory
//	factory := nats.NewFactory()
//
//	// Use with tenant manager
//	manager, err := tenant.NewManager(tenant.ManagerConfig{
//	    Config:  cfg,
//	    Factory: factory,
//	})
//
// # Global Functions
//
// Convenience functions using the default registry:
//
//	// Register a transport globally
//	transport.Register(myTransport)
//
//	// Get a transport globally
//	t, ok := transport.Get("nats")
//
// # Connection Lifecycle
//
// Transport connections follow this lifecycle:
//
//	1. Connect(ctx, cfg, creds) - Establish connection
//	2. Use Publisher/Subscriber - Send/receive messages
//	3. conn.Close(ctx) - Graceful shutdown
//
// The transport handles reconnection automatically based on configuration.
//
// # Error Handling
//
// Transport operations may return these errors:
//
//	core.ErrNotConnected      - Not connected to server
//	core.ErrConnectionTimeout - Connection attempt timed out
//	core.ErrAuthFailed        - Authentication failed
//	core.ErrPermissionDenied  - Insufficient permissions
package transport
