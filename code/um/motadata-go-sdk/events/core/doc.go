// Package core provides the fundamental interfaces, types, and utilities for the events system.
//
// This package defines the core abstractions used throughout the events SDK:
//   - Message: The fundamental unit of communication
//   - Headers: Message metadata and context propagation
//   - Connection: Interface for managing transport connections
//   - Publisher: Interface for sending messages
//   - Subscriber: Interface for receiving messages
//   - Error types: Structured error handling
//   - Context utilities: Tenant ID, trace context, correlation ID propagation
//   - Object pools: Memory-efficient message and header allocation
//
// # Message Structure
//
// Messages are the fundamental unit of communication:
//
//	msg := core.NewMessage([]byte(`{"event": "user.created"}`))
//	msg.WithSubject("events.users")
//	msg.WithHeader("Content-Type", "application/json")
//	msg.WithID("unique-id-for-deduplication")
//
// # Headers
//
// Headers provide message metadata and support context propagation:
//
//	headers := make(core.Headers)
//	headers.Set("Content-Type", "application/json")
//	headers.Set("X-Tenant-ID", "tenant-123")
//	headers.Add("X-Tags", "important")
//	headers.Add("X-Tags", "urgent")
//
//	// Get single value
//	contentType := headers.Get("Content-Type")
//
//	// Get all values
//	tags := headers.Values("X-Tags") // ["important", "urgent"]
//
// # Standard Headers
//
// The package defines standard header keys:
//
//	core.HeaderContentType   - MIME content type
//	core.HeaderMessageID     - Unique message ID (NATS deduplication)
//	core.HeaderCorrelationID - Request correlation ID
//	core.HeaderTenantID      - Multi-tenant identifier
//	core.HeaderTraceID       - Distributed trace ID
//	core.HeaderSpanID        - Span ID
//	core.HeaderTraceParent   - W3C Trace Context format
//	core.HeaderTraceState    - W3C Trace State
//
// # Context Propagation
//
// Propagate context values through message headers:
//
//	// Add tenant ID to context
//	ctx := core.WithTenantID(ctx, "tenant-123")
//
//	// Add trace context
//	ctx = core.WithTraceContext(ctx, &core.TraceContext{
//	    TraceID: "abc123",
//	    SpanID:  "def456",
//	    Sampled: true,
//	})
//
//	// Extract context values into headers
//	headers := core.ExtractHeaders(ctx, msg.Headers)
//
//	// Inject headers back into context (on receive)
//	ctx = core.InjectContext(ctx, msg.Headers)
//
// # Connection Interface
//
// The Connection interface manages transport connections:
//
//	type Connection interface {
//	    Connect(ctx context.Context) error
//	    Close(ctx context.Context) error
//	    IsConnected() bool
//	    State() ConnectionState
//	}
//
// Connection states:
//
//	core.StateDisconnected  - Not connected
//	core.StateConnecting    - Connection in progress
//	core.StateConnected     - Successfully connected
//	core.StateReconnecting  - Reconnection in progress
//	core.StateDraining      - Draining before close
//	core.StateClosed        - Connection closed
//
// # Publisher Interface
//
// The Publisher interface handles message publishing:
//
//	type Publisher interface {
//	    Publish(ctx context.Context, subject string, msg *Message) error
//	    PublishAsync(ctx context.Context, subject string, msg *Message) PubAckFuture
//	    Request(ctx context.Context, subject string, msg *Message) (*Message, error)
//	    Close(ctx context.Context) error
//	}
//
// # Subscriber Interface
//
// The Subscriber interface handles message subscriptions:
//
//	type Subscriber interface {
//	    Subscribe(ctx context.Context, subject string, handler MessageHandler) (Subscription, error)
//	    QueueSubscribe(ctx context.Context, subject string, queue string, handler MessageHandler) (Subscription, error)
//	    Close(ctx context.Context) error
//	}
//
// # Error Handling
//
// The package provides structured error types:
//
//	// Check for specific errors
//	if errors.Is(err, core.ErrNotConnected) {
//	    // Handle not connected
//	}
//
//	// Check if error is retryable
//	if core.IsRetryable(err) {
//	    // Retry the operation
//	}
//
//	// Check if error is temporary
//	if core.IsTemporary(err) {
//	    // Wait and retry
//	}
//
//	// Handle tenant-specific errors
//	var tenantErr core.TenantError
//	if errors.As(err, &tenantErr) {
//	    log.Printf("Tenant %s error: %v", tenantErr.TenantID, tenantErr.Err)
//	}
//
// # Object Pools
//
// Use object pools for memory-efficient operations in high-throughput scenarios:
//
//	// Acquire message from pool
//	msg := core.AcquireMessage()
//	msg.Data = []byte("data")
//
//	// Release back to pool when done
//	defer core.ReleaseMessage(msg)
//
//	// Acquire headers
//	headers := core.AcquireHeaders()
//	defer core.ReleaseHeaders(headers)
//
// # Error Categories
//
// Common error categories:
//   - Connection errors: ErrNotConnected, ErrConnectionClosed, ErrConnectionTimeout
//   - Publishing errors: ErrPublishFailed, ErrPublishTimeout, ErrInvalidSubject
//   - Subscription errors: ErrSubscriptionClosed, ErrSubscriptionInvalid
//   - Authentication errors: ErrAuthFailed, ErrPermissionDenied
//   - Tenant errors: ErrTenantNotFound, ErrTenantExists
//   - Configuration errors: ErrInvalidConfig, ErrMissingConfig
package core
