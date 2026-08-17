package transport

import (
	"context"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/config"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

// Transport represents a messaging transport backend
type Transport interface {
	// Name returns the transport name (e.g., "nats", "kafka")
	Name() string

	// Connect establishes a connection
	Connect(ctx context.Context, cfg *config.Config, creds *auth.Credentials) (core.Connection, error)

	// Publisher creates a publisher for the connection
	Publisher(conn core.Connection) (core.Publisher, error)

	// Subscriber creates a subscriber for the connection
	Subscriber(conn core.Connection) (core.Subscriber, error)
}

// Factory creates transport instances
type Factory interface {
	// Create creates a transport by name
	Create(name string) (Transport, error)
}

// Registry holds registered transports
type Registry struct {
	transports map[string]Transport
}

// NewRegistry creates a new transport registry
func NewRegistry() *Registry {
	return &Registry{
		transports: make(map[string]Transport),
	}
}

// Register adds a transport to the registry
func (r *Registry) Register(t Transport) {
	r.transports[t.Name()] = t
}

// Get retrieves a transport by name
func (r *Registry) Get(name string) (Transport, bool) {
	t, ok := r.transports[name]
	return t, ok
}

// Names returns all registered transport names
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.transports))
	for name := range r.transports {
		names = append(names, name)
	}
	return names
}

// DefaultRegistry is the global transport registry
var DefaultRegistry = NewRegistry()

// Register adds a transport to the default registry
func Register(t Transport) {
	DefaultRegistry.Register(t)
}

// Get retrieves a transport from the default registry
func Get(name string) (Transport, bool) {
	return DefaultRegistry.Get(name)
}
