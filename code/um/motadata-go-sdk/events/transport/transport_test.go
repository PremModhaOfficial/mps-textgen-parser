package transport

import (
	"context"
	"testing"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/config"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

// ============== Mock Transport ==============

type mockTransport struct {
	name string
}

func newMockTransport(name string) *mockTransport {
	return &mockTransport{name: name}
}

func (m *mockTransport) Name() string {
	return m.name
}

func (m *mockTransport) Connect(ctx context.Context, cfg *config.Config, creds *auth.Credentials) (core.Connection, error) {
	return nil, nil
}

func (m *mockTransport) Publisher(conn core.Connection) (core.Publisher, error) {
	return nil, nil
}

func (m *mockTransport) Subscriber(conn core.Connection) (core.Subscriber, error) {
	return nil, nil
}

// ============== Registry Tests ==============

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	if r == nil {
		t.Fatal("NewRegistry returned nil")
	}
	if r.transports == nil {
		t.Error("transports map not initialized")
	}
}

func TestRegistryRegister(t *testing.T) {
	r := NewRegistry()
	transport := newMockTransport("test")

	r.Register(transport)

	if len(r.transports) != 1 {
		t.Errorf("expected 1 transport, got %d", len(r.transports))
	}
}

func TestRegistryGet(t *testing.T) {
	r := NewRegistry()
	transport := newMockTransport("nats")
	r.Register(transport)

	// Get existing
	got, ok := r.Get("nats")
	if !ok {
		t.Error("expected transport to exist")
	}
	if got.Name() != "nats" {
		t.Errorf("unexpected name: %s", got.Name())
	}

	// Get non-existent
	_, ok = r.Get("unknown")
	if ok {
		t.Error("expected transport to not exist")
	}
}

func TestRegistryNames(t *testing.T) {
	r := NewRegistry()
	r.Register(newMockTransport("nats"))
	r.Register(newMockTransport("kafka"))

	names := r.Names()
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}

func TestDefaultRegistry(t *testing.T) {
	if DefaultRegistry == nil {
		t.Fatal("DefaultRegistry is nil")
	}
}

func TestGlobalRegister(t *testing.T) {
	// Create a fresh registry for this test
	originalRegistry := DefaultRegistry
	DefaultRegistry = NewRegistry()
	defer func() { DefaultRegistry = originalRegistry }()

	transport := newMockTransport("global-test")
	Register(transport)

	got, ok := Get("global-test")
	if !ok {
		t.Error("expected transport to exist in default registry")
	}
	if got.Name() != "global-test" {
		t.Errorf("unexpected name: %s", got.Name())
	}
}

// ============== Benchmarks ==============

func BenchmarkRegistryRegister(b *testing.B) {
	r := NewRegistry()
	transport := newMockTransport("test")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r.Register(transport)
	}
}

func BenchmarkRegistryGet(b *testing.B) {
	r := NewRegistry()
	r.Register(newMockTransport("test"))

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r.Get("test")
	}
}
