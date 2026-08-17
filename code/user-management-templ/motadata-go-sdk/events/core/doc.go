// Package core provides context utilities, header constants, and minimal
// NATS-typed interfaces for the events system.
//
// This package defines:
//   - Context utilities: Tenant ID, trace context, correlation ID propagation
//   - Header constants: Standard NATS header keys
//   - Minimal interfaces: Publisher, Subscriber, Subscription using *nats.Msg
//
// Error types and structured error handling have been moved to the utils package.
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
// Propagate context values through NATS message headers:
//
//	// Add tenant ID to context
//	ctx := core.WithTenantID(ctx, "tenant-123")
//
//	// Extract context values into NATS headers
//	headers := core.ExtractHeaders(ctx, msg.Header)
//
//	// Inject headers back into context (on receive)
//	ctx = core.InjectContext(ctx, msg.Header)
//
// # Error Handling
//
// Error types are now in the utils package:
//
//	if errors.Is(err, utils.ErrNotConnected) {
//	    // Handle not connected
//	}
//
//	if utils.IsRetryable(err) {
//	    // Retry the operation
//	}
//
// # NATS-typed Interfaces
//
// Minimal interfaces to avoid circular imports between packages:
//
//	core.Publisher   - Publish and request using *nats.Msg
//	core.Subscriber  - Subscribe using *nats.Msg handlers
//	core.Subscription - Active subscription management
package core
