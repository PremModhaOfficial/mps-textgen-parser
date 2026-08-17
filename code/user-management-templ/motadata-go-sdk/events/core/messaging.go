package core

import (
	"context"

	"github.com/nats-io/nats.go"
)

// MessageHandler is a callback function that processes received NATS messages.
// It accepts a context carrying tenant, trace, and correlation metadata, along with
// the raw NATS message. Returning a non-nil error signals the caller (typically a
// subscriber or middleware chain) that processing failed, enabling retry or
// dead-letter handling upstream.
type MessageHandler func(ctx context.Context, msg *nats.Msg) error

// Publisher defines the contract for publishing messages to NATS subjects.
// Implementations are responsible for connection management, header injection
// (tenant ID, trace context), and message serialization. The interface is
// intentionally thin so that middleware (retry, tracing, metrics) can wrap it
// without knowledge of the underlying NATS connection details.
type Publisher interface {
	// Publish sends a message to the given subject in a fire-and-forget fashion.
	// Context values (tenant ID, trace context) should be propagated into
	// message headers by the implementation.
	Publish(ctx context.Context, subject string, msg *nats.Msg) error

	// Request sends a message and waits for a single reply, implementing the
	// request-reply pattern. The context deadline controls the maximum wait time.
	Request(ctx context.Context, subject string, msg *nats.Msg) (*nats.Msg, error)

	// Close gracefully shuts down the publisher, draining any in-flight messages
	// before releasing the underlying NATS connection resources.
	Close(ctx context.Context) error
}

// Subscriber defines the contract for consuming messages from NATS subjects.
// It supports both regular subscriptions (every subscriber gets every message)
// and queue subscriptions (messages are load-balanced across subscribers in the
// same queue group). Implementations should inject incoming NATS headers into
// the handler's context so that downstream code can access tenant and trace data.
type Subscriber interface {
	// Subscribe creates a subscription where every instance receives all messages
	// published to the subject. Suitable for fan-out scenarios such as cache
	// invalidation or event sourcing projections.
	Subscribe(ctx context.Context, subject string, handler MessageHandler) (Subscription, error)

	// QueueSubscribe creates a queue-group subscription where messages are
	// distributed across subscribers sharing the same queue name. This provides
	// automatic load balancing for horizontally scaled services.
	QueueSubscribe(ctx context.Context, subject string, queue string, handler MessageHandler) (Subscription, error)

	// Close gracefully shuts down the subscriber, draining active subscriptions
	// so that in-flight messages finish processing before the connection is released.
	Close(ctx context.Context) error
}

// Subscription represents an active NATS subscription and provides lifecycle
// control. Callers should prefer Drain over Unsubscribe for graceful shutdown,
// as Drain allows in-flight message handlers to complete before removing the
// subscription.
type Subscription interface {
	// Subject returns the NATS subject this subscription is listening on.
	Subject() string

	// Unsubscribe immediately removes the subscription, dropping any pending
	// messages. Use Drain instead when graceful shutdown is preferred.
	Unsubscribe() error

	// Drain signals NATS to stop delivering new messages while allowing
	// in-flight handlers to finish. Once drained, the subscription becomes
	// invalid and cannot be reused.
	Drain() error

	// IsValid reports whether the subscription is still active and capable
	// of receiving messages. Returns false after Unsubscribe or Drain completes.
	IsValid() bool
}
