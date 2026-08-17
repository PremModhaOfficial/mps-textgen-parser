package core

import (
	"context"
)

// Connection represents a connection to a messaging system
type Connection interface {
	// Connect establishes the connection
	Connect(ctx context.Context) error

	// Close gracefully closes the connection
	Close(ctx context.Context) error

	// IsConnected returns true if currently connected
	IsConnected() bool

	// State returns the current connection state
	State() ConnectionState
}

// Publisher publishes messages to subjects/topics
type Publisher interface {
	// Publish sends a message to the specified subject
	Publish(ctx context.Context, subject string, msg *Message) error

	// PublishAsync sends a message asynchronously and returns a future
	PublishAsync(ctx context.Context, subject string, msg *Message) PubAckFuture

	// Request sends a message and waits for a reply (request-reply pattern)
	Request(ctx context.Context, subject string, msg *Message) (*Message, error)

	// Close gracefully closes the publisher
	Close(ctx context.Context) error
}

// Subscriber subscribes to subjects/topics and receives messages
type Subscriber interface {
	// Subscribe creates a subscription to the specified subject
	Subscribe(ctx context.Context, subject string, handler MessageHandler) (Subscription, error)

	// QueueSubscribe creates a queue subscription for load balancing
	QueueSubscribe(ctx context.Context, subject string, queue string, handler MessageHandler) (Subscription, error)

	// Close gracefully closes the subscriber
	Close(ctx context.Context) error
}

// Subscription represents an active subscription
type Subscription interface {
	// Subject returns the subscribed subject
	Subject() string

	// Unsubscribe removes the subscription
	Unsubscribe() error

	// Drain unsubscribes and waits for messages to be processed
	Drain() error

	// IsValid returns true if subscription is still active
	IsValid() bool
}

// MessageHandler processes received messages
type MessageHandler func(ctx context.Context, msg *Message) error

// PubAckFuture represents a pending publish acknowledgment
type PubAckFuture interface {
	// Ok returns a channel that receives the ack on success
	Ok() <-chan *PubAck

	// Err returns a channel that receives error on failure
	Err() <-chan error

	// Wait blocks until the publish completes or context cancels
	Wait(ctx context.Context) (*PubAck, error)
}

// PubAck represents a publish acknowledgment
type PubAck struct {
	Stream   string // Stream name (for JetStream)
	Sequence uint64 // Sequence number
	Domain   string // JetStream domain
}

// Component represents a managed component for graceful shutdown
type Component interface {
	// Name returns the component name for identification
	Name() string

	// Shutdown gracefully shuts down the component
	Shutdown(ctx context.Context) error
}

// HealthChecker provides health status
type HealthChecker interface {
	// Health returns the current health status
	Health() HealthStatus
}

// ConnectionState represents the state of a connection
type ConnectionState int

const (
	StateDisconnected ConnectionState = iota
	StateConnecting
	StateConnected
	StateReconnecting
	StateDraining
	StateClosed
)

// String returns a human-readable state name
func (s ConnectionState) String() string {
	switch s {
	case StateDisconnected:
		return "disconnected"
	case StateConnecting:
		return "connecting"
	case StateConnected:
		return "connected"
	case StateReconnecting:
		return "reconnecting"
	case StateDraining:
		return "draining"
	case StateClosed:
		return "closed"
	default:
		return "unknown"
	}
}

// HealthStatus represents the health of a component
type HealthStatus struct {
	State     ConnectionState
	Healthy   bool
	Message   string
	LastError error
	Details   map[string]any
}

// IsHealthy returns true if the status indicates healthy state
func (h HealthStatus) IsHealthy() bool {
	return h.Healthy && h.State == StateConnected
}
