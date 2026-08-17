package tenant

import (
	"context"
	"sync"
	"testing"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/config"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

// ============== Mock Implementations ==============

type mockConnection struct {
	tenantID   string
	connected  bool
	lastTouch  time.Time
	health     core.HealthStatus
	publisher  *mockPublisher
	subscriber *mockSubscriber
	state      core.ConnectionState
	mu         sync.Mutex
}

func newMockConnection(tenantID string) *mockConnection {
	return &mockConnection{
		tenantID:   tenantID,
		lastTouch:  time.Now(),
		publisher:  &mockPublisher{},
		subscriber: &mockSubscriber{},
		state:      core.StateDisconnected,
		health: core.HealthStatus{
			State:   core.StateConnected,
			Healthy: true,
			Message: "healthy",
		},
	}
}

func (m *mockConnection) TenantID() string {
	return m.tenantID
}

func (m *mockConnection) Connect(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connected = true
	m.state = core.StateConnected
	return nil
}

func (m *mockConnection) Close(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connected = false
	m.state = core.StateClosed
	return nil
}

func (m *mockConnection) IsConnected() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.connected
}

func (m *mockConnection) State() core.ConnectionState {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.state
}

func (m *mockConnection) Touch() {
	m.mu.Lock()
	m.lastTouch = time.Now()
	m.mu.Unlock()
}

func (m *mockConnection) IsIdle(timeout time.Duration) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return time.Since(m.lastTouch) > timeout
}

func (m *mockConnection) Health() core.HealthStatus {
	return m.health
}

func (m *mockConnection) Publisher() core.Publisher {
	return m.publisher
}

func (m *mockConnection) Subscriber() core.Subscriber {
	return m.subscriber
}

type mockPublisher struct{}

func (m *mockPublisher) Publish(ctx context.Context, subject string, msg *core.Message) error {
	return nil
}

func (m *mockPublisher) PublishAsync(ctx context.Context, subject string, msg *core.Message) core.PubAckFuture {
	return &mockPubAckFuture{}
}

func (m *mockPublisher) Request(ctx context.Context, subject string, msg *core.Message) (*core.Message, error) {
	return core.NewMessage([]byte("response")), nil
}

func (m *mockPublisher) Close(ctx context.Context) error {
	return nil
}

type mockPubAckFuture struct{}

func (m *mockPubAckFuture) Ok() <-chan *core.PubAck {
	ch := make(chan *core.PubAck, 1)
	ch <- &core.PubAck{}
	return ch
}

func (m *mockPubAckFuture) Err() <-chan error {
	return make(chan error)
}

func (m *mockPubAckFuture) Wait(ctx context.Context) (*core.PubAck, error) {
	return &core.PubAck{}, nil
}

type mockSubscriber struct{}

func (m *mockSubscriber) Subscribe(ctx context.Context, subject string, handler core.MessageHandler) (core.Subscription, error) {
	return &mockSubscription{subject: subject}, nil
}

func (m *mockSubscriber) QueueSubscribe(ctx context.Context, subject, queue string, handler core.MessageHandler) (core.Subscription, error) {
	return &mockSubscription{subject: subject}, nil
}

func (m *mockSubscriber) Close(ctx context.Context) error {
	return nil
}

type mockSubscription struct {
	subject string
}

func (m *mockSubscription) Subject() string   { return m.subject }
func (m *mockSubscription) Unsubscribe() error { return nil }
func (m *mockSubscription) Drain() error       { return nil }
func (m *mockSubscription) IsValid() bool      { return true }

type mockConnectionFactory struct {
	connections map[string]*mockConnection
	createError error
}

func newMockConnectionFactory() *mockConnectionFactory {
	return &mockConnectionFactory{
		connections: make(map[string]*mockConnection),
	}
}

func (f *mockConnectionFactory) Create(ctx context.Context, tenantID string, cfg *config.Config, creds *auth.Credentials) (Connection, error) {
	if f.createError != nil {
		return nil, f.createError
	}
	conn := newMockConnection(tenantID)
	f.connections[tenantID] = conn
	return conn, nil
}

// ============== Registry Tests ==============

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	if r == nil {
		t.Fatal("NewRegistry returned nil")
	}
	if r.tenants == nil {
		t.Error("tenants map not initialized")
	}
}

func TestRegistryRegister(t *testing.T) {
	r := NewRegistry()

	// Valid registration
	info := &Info{
		ID:     "tenant-1",
		Name:   "Test Tenant",
		Active: true,
	}
	err := r.Register(info)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify registration
	if r.Count() != 1 {
		t.Errorf("expected count 1, got %d", r.Count())
	}

	// Duplicate registration
	err = r.Register(info)
	if err != core.ErrTenantExists {
		t.Errorf("expected ErrTenantExists, got %v", err)
	}

	// Nil info
	err = r.Register(nil)
	if err != core.ErrTenantNotFound {
		t.Errorf("expected ErrTenantNotFound for nil, got %v", err)
	}

	// Empty ID
	err = r.Register(&Info{})
	if err != core.ErrTenantNotFound {
		t.Errorf("expected ErrTenantNotFound for empty ID, got %v", err)
	}
}

func TestRegistryRegisterCallback(t *testing.T) {
	r := NewRegistry()

	var registeredID string
	r.OnRegistered(func(tenantID string) {
		registeredID = tenantID
	})

	r.Register(&Info{ID: "tenant-1"})
	if registeredID != "tenant-1" {
		t.Errorf("expected callback with tenant-1, got %s", registeredID)
	}
}

func TestRegistryUnregister(t *testing.T) {
	r := NewRegistry()
	r.Register(&Info{ID: "tenant-1"})

	// Unregister existing
	err := r.Unregister("tenant-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify removal
	if r.Count() != 0 {
		t.Errorf("expected count 0, got %d", r.Count())
	}

	// Unregister non-existent
	err = r.Unregister("unknown")
	if err != core.ErrTenantNotFound {
		t.Errorf("expected ErrTenantNotFound, got %v", err)
	}
}

func TestRegistryUnregisterCallback(t *testing.T) {
	r := NewRegistry()
	r.Register(&Info{ID: "tenant-1"})

	var unregisteredID string
	r.OnUnregistered(func(tenantID string) {
		unregisteredID = tenantID
	})

	r.Unregister("tenant-1")
	if unregisteredID != "tenant-1" {
		t.Errorf("expected callback with tenant-1, got %s", unregisteredID)
	}
}

func TestRegistryGet(t *testing.T) {
	r := NewRegistry()
	r.Register(&Info{ID: "tenant-1", Name: "Test"})

	// Get existing
	info, ok := r.Get("tenant-1")
	if !ok {
		t.Error("expected tenant to exist")
	}
	if info.Name != "Test" {
		t.Errorf("unexpected name: %s", info.Name)
	}

	// Get non-existent
	_, ok = r.Get("unknown")
	if ok {
		t.Error("expected tenant to not exist")
	}
}

func TestRegistryExists(t *testing.T) {
	r := NewRegistry()
	r.Register(&Info{ID: "tenant-1"})

	if !r.Exists("tenant-1") {
		t.Error("expected tenant-1 to exist")
	}
	if r.Exists("unknown") {
		t.Error("expected unknown to not exist")
	}
}

func TestRegistryList(t *testing.T) {
	r := NewRegistry()
	r.Register(&Info{ID: "tenant-1"})
	r.Register(&Info{ID: "tenant-2"})

	list := r.List()
	if len(list) != 2 {
		t.Errorf("expected 2 tenants, got %d", len(list))
	}
}

func TestRegistryAll(t *testing.T) {
	r := NewRegistry()
	r.Register(&Info{ID: "tenant-1", Name: "Test 1"})
	r.Register(&Info{ID: "tenant-2", Name: "Test 2"})

	all := r.All()
	if len(all) != 2 {
		t.Errorf("expected 2 infos, got %d", len(all))
	}
}

func TestRegistryCount(t *testing.T) {
	r := NewRegistry()
	if r.Count() != 0 {
		t.Errorf("expected count 0, got %d", r.Count())
	}

	r.Register(&Info{ID: "tenant-1"})
	if r.Count() != 1 {
		t.Errorf("expected count 1, got %d", r.Count())
	}
}

func TestRegistryUpdate(t *testing.T) {
	r := NewRegistry()
	r.Register(&Info{ID: "tenant-1", Name: "Original"})

	// Update existing
	err := r.Update("tenant-1", func(info *Info) {
		info.Name = "Updated"
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, _ := r.Get("tenant-1")
	if info.Name != "Updated" {
		t.Errorf("expected Updated, got %s", info.Name)
	}

	// Update non-existent
	err = r.Update("unknown", func(info *Info) {})
	if err != core.ErrTenantNotFound {
		t.Errorf("expected ErrTenantNotFound, got %v", err)
	}
}

func TestRegistrySetActive(t *testing.T) {
	r := NewRegistry()
	r.Register(&Info{ID: "tenant-1", Active: false})

	err := r.SetActive("tenant-1", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, _ := r.Get("tenant-1")
	if !info.Active {
		t.Error("expected active to be true")
	}

	// Set active on non-existent
	err = r.SetActive("unknown", true)
	if err != core.ErrTenantNotFound {
		t.Errorf("expected ErrTenantNotFound, got %v", err)
	}
}

func TestRegistryActiveTenants(t *testing.T) {
	r := NewRegistry()
	r.Register(&Info{ID: "tenant-1", Active: true})
	r.Register(&Info{ID: "tenant-2", Active: false})
	r.Register(&Info{ID: "tenant-3", Active: true})

	active := r.ActiveTenants()
	if len(active) != 2 {
		t.Errorf("expected 2 active tenants, got %d", len(active))
	}
}

func TestRegistryConcurrentAccess(t *testing.T) {
	r := NewRegistry()
	var wg sync.WaitGroup

	// Concurrent registrations
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			tenantID := "tenant-" + string(rune('a'+id%26))
			r.Register(&Info{ID: tenantID})
		}(i)
	}

	// Concurrent reads
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			tenantID := "tenant-" + string(rune('a'+id%26))
			r.Get(tenantID)
			r.Exists(tenantID)
			r.List()
		}(i)
	}

	wg.Wait()
}

// ============== Manager Tests ==============

func TestNewManager(t *testing.T) {
	factory := newMockConnectionFactory()
	cfg := ManagerConfig{
		Factory: factory,
	}

	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m == nil {
		t.Fatal("NewManager returned nil")
	}

	// Cleanup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	m.Shutdown(ctx)
}

func TestNewManagerInvalidConfig(t *testing.T) {
	cfg := ManagerConfig{
		Config: &config.Config{
			// CleanupInterval 0 is invalid
			CleanupInterval: 0,
		},
	}

	_, err := NewManager(cfg)
	if err == nil {
		t.Error("expected error for invalid config")
	}
}

func TestManagerConnect(t *testing.T) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{Factory: factory})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	creds := auth.UserPass("user", "pass")

	conn, err := m.Connect(ctx, "tenant-1", creds)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn == nil {
		t.Error("expected connection")
	}
	if conn.TenantID() != "tenant-1" {
		t.Errorf("unexpected tenant ID: %s", conn.TenantID())
	}
}

func TestManagerConnectExisting(t *testing.T) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{Factory: factory})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	creds := auth.UserPass("user", "pass")

	// First connect
	conn1, _ := m.Connect(ctx, "tenant-1", creds)

	// Second connect should return existing
	conn2, _ := m.Connect(ctx, "tenant-1", creds)
	if conn1.TenantID() != conn2.TenantID() {
		t.Error("expected same connection to be returned")
	}
}

func TestManagerConnectNoCredentials(t *testing.T) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{Factory: factory})
	defer m.Shutdown(context.Background())

	ctx := context.Background()

	_, err := m.Connect(ctx, "tenant-1", nil)
	if err == nil {
		t.Error("expected error for nil credentials")
	}
}

func TestManagerConnectNoFactory(t *testing.T) {
	m, _ := NewManager(ManagerConfig{})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	creds := auth.UserPass("user", "pass")

	_, err := m.Connect(ctx, "tenant-1", creds)
	if err == nil {
		t.Error("expected error for nil factory")
	}
}

func TestManagerConnectWithCredentialManager(t *testing.T) {
	factory := newMockConnectionFactory()
	cm := auth.NewCredentialManager(auth.DefaultCredentialManagerConfig())
	cm.Register("tenant-1", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig", "")

	m, _ := NewManager(ManagerConfig{
		Factory:     factory,
		CredManager: cm,
	})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	conn, err := m.ConnectWithCredentialManager(ctx, "tenant-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn == nil {
		t.Error("expected connection")
	}
}

func TestManagerConnectWithCredentialManagerNoManager(t *testing.T) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{Factory: factory})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	_, err := m.ConnectWithCredentialManager(ctx, "tenant-1")
	if err == nil {
		t.Error("expected error when no credential manager")
	}
}

func TestManagerConnectWithPayload(t *testing.T) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{Factory: factory})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	payload := auth.RegistrationPayload{
		TenantID: "tenant-1",
		JWT:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig",
	}

	conn, err := m.ConnectWithPayload(ctx, payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn == nil {
		t.Error("expected connection")
	}
}

func TestManagerConnectWithPayloadInvalid(t *testing.T) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{Factory: factory})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	payload := auth.RegistrationPayload{} // Invalid

	_, err := m.ConnectWithPayload(ctx, payload)
	if err == nil {
		t.Error("expected error for invalid payload")
	}
}

func TestManagerGet(t *testing.T) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{Factory: factory})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	creds := auth.UserPass("user", "pass")
	m.Connect(ctx, "tenant-1", creds)

	// Get existing
	conn, ok := m.Get("tenant-1")
	if !ok {
		t.Error("expected connection to exist")
	}
	if conn == nil {
		t.Error("expected non-nil connection")
	}

	// Get non-existent
	_, ok = m.Get("unknown")
	if ok {
		t.Error("expected connection to not exist")
	}
}

func TestManagerDisconnect(t *testing.T) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{Factory: factory})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	creds := auth.UserPass("user", "pass")
	m.Connect(ctx, "tenant-1", creds)

	// Disconnect
	err := m.Disconnect(ctx, "tenant-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify disconnected
	_, ok := m.Get("tenant-1")
	if ok {
		t.Error("expected connection to be removed")
	}

	// Disconnect non-existent
	err = m.Disconnect(ctx, "unknown")
	if err != core.ErrTenantNotFound {
		t.Errorf("expected ErrTenantNotFound, got %v", err)
	}
}

func TestManagerShutdown(t *testing.T) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{Factory: factory})

	ctx := context.Background()
	creds := auth.UserPass("user", "pass")
	m.Connect(ctx, "tenant-1", creds)
	m.Connect(ctx, "tenant-2", creds)

	// Shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := m.Shutdown(shutdownCtx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify all connections closed
	if m.ActiveConnections() != 0 {
		t.Errorf("expected 0 connections, got %d", m.ActiveConnections())
	}
}

func TestManagerShutdownTimeout(t *testing.T) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{Factory: factory})

	ctx := context.Background()
	creds := auth.UserPass("user", "pass")
	m.Connect(ctx, "tenant-1", creds)

	// Shutdown with very short timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Should not panic
	m.Shutdown(shutdownCtx)
}

func TestManagerErrors(t *testing.T) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{Factory: factory})
	defer m.Shutdown(context.Background())

	errChan := m.Errors()
	if errChan == nil {
		t.Error("expected error channel")
	}
}

func TestManagerOnError(t *testing.T) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{Factory: factory})
	defer m.Shutdown(context.Background())

	var receivedErr core.TenantError
	m.OnError(func(err core.TenantError) {
		receivedErr = err
	})

	// Dispatch error
	testErr := core.TenantError{TenantID: "tenant-1", Err: core.ErrTenantNotFound}
	m.dispatchError(testErr)

	if receivedErr.TenantID != "tenant-1" {
		t.Errorf("expected tenant-1, got %s", receivedErr.TenantID)
	}
}

func TestManagerHealth(t *testing.T) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{Factory: factory})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	creds := auth.UserPass("user", "pass")
	m.Connect(ctx, "tenant-1", creds)
	m.Connect(ctx, "tenant-2", creds)

	health := m.Health()
	if len(health) != 2 {
		t.Errorf("expected 2 health entries, got %d", len(health))
	}
}

func TestManagerTenantHealth(t *testing.T) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{Factory: factory})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	creds := auth.UserPass("user", "pass")
	m.Connect(ctx, "tenant-1", creds)

	// Get existing
	health, ok := m.TenantHealth("tenant-1")
	if !ok {
		t.Error("expected health status")
	}
	if !health.Healthy {
		t.Error("expected healthy to be true")
	}

	// Get non-existent
	_, ok = m.TenantHealth("unknown")
	if ok {
		t.Error("expected no health status for unknown tenant")
	}
}

func TestManagerActiveConnections(t *testing.T) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{Factory: factory})
	defer m.Shutdown(context.Background())

	if m.ActiveConnections() != 0 {
		t.Errorf("expected 0 connections, got %d", m.ActiveConnections())
	}

	ctx := context.Background()
	creds := auth.UserPass("user", "pass")
	m.Connect(ctx, "tenant-1", creds)

	if m.ActiveConnections() != 1 {
		t.Errorf("expected 1 connection, got %d", m.ActiveConnections())
	}
}

func TestManagerTenantIDs(t *testing.T) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{Factory: factory})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	creds := auth.UserPass("user", "pass")
	m.Connect(ctx, "tenant-1", creds)
	m.Connect(ctx, "tenant-2", creds)

	ids := m.TenantIDs()
	if len(ids) != 2 {
		t.Errorf("expected 2 IDs, got %d", len(ids))
	}
}

func TestManagerSetCredentialManager(t *testing.T) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{Factory: factory})
	defer m.Shutdown(context.Background())

	cm := auth.NewCredentialManager(auth.DefaultCredentialManagerConfig())
	m.SetCredentialManager(cm)

	if m.CredentialManager() != cm {
		t.Error("credential manager not set correctly")
	}
}

func TestManagerConnectAfterShutdown(t *testing.T) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{Factory: factory})

	// Shutdown first
	m.Shutdown(context.Background())

	// Try to connect
	ctx := context.Background()
	creds := auth.UserPass("user", "pass")
	_, err := m.Connect(ctx, "tenant-1", creds)
	if err != core.ErrShutdownInProgress {
		t.Errorf("expected ErrShutdownInProgress, got %v", err)
	}
}

func TestManagerConcurrentConnect(t *testing.T) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{Factory: factory})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	creds := auth.UserPass("user", "pass")
	var wg sync.WaitGroup

	// Concurrent connects to same tenant
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Connect(ctx, "tenant-1", creds)
		}()
	}

	wg.Wait()

	// Should only have one connection
	if m.ActiveConnections() != 1 {
		t.Errorf("expected 1 connection, got %d", m.ActiveConnections())
	}
}

func TestManagerDispatchErrorAfterShutdown(t *testing.T) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{Factory: factory})

	// Shutdown first
	m.Shutdown(context.Background())

	// Dispatch error should not panic
	testErr := core.TenantError{TenantID: "tenant-1", Err: core.ErrTenantNotFound}
	m.dispatchError(testErr)
}

func TestManagerCleanupIdleConnections(t *testing.T) {
	factory := newMockConnectionFactory()
	testConfig := config.DefaultConfig()
	testConfig.IdleTimeout = 50 * time.Millisecond     // Very short for testing
	testConfig.CleanupInterval = 25 * time.Millisecond // Very short for testing

	cfg := ManagerConfig{
		Factory: factory,
		Config:  testConfig,
	}
	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("failed to create manager: %v", err)
	}
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	creds := auth.UserPass("user", "pass")
	m.Connect(ctx, "tenant-1", creds)

	// Wait for idle timeout and cleanup
	time.Sleep(100 * time.Millisecond)

	// Connection should be cleaned up due to idle timeout
	if m.ActiveConnections() != 0 {
		t.Errorf("expected 0 connections after idle cleanup, got %d", m.ActiveConnections())
	}
}

func TestManagerConnectWithCredentialManagerFromManager(t *testing.T) {
	factory := newMockConnectionFactory()
	cm := auth.NewCredentialManager(auth.DefaultCredentialManagerConfig())
	cm.Register("tenant-1", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig", "")

	m, _ := NewManager(ManagerConfig{
		Factory:     factory,
		CredManager: cm,
	})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	// Connect with nil creds but credential manager set
	conn, err := m.Connect(ctx, "tenant-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn == nil {
		t.Error("expected connection")
	}
}

func TestManagerConnectWithPayloadAndCredentialManager(t *testing.T) {
	factory := newMockConnectionFactory()
	cm := auth.NewCredentialManager(auth.DefaultCredentialManagerConfig())

	m, _ := NewManager(ManagerConfig{
		Factory:     factory,
		CredManager: cm,
	})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	payload := auth.RegistrationPayload{
		TenantID: "tenant-1",
		JWT:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig",
	}

	conn, err := m.ConnectWithPayload(ctx, payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn == nil {
		t.Error("expected connection")
	}

	// Verify credentials were registered
	if !cm.HasValid("tenant-1") {
		t.Error("expected credentials to be registered")
	}
}

func TestManagerConnectWithEmptyCredentials(t *testing.T) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{Factory: factory})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	// Empty credentials
	creds := &auth.Credentials{}

	_, err := m.Connect(ctx, "tenant-1", creds)
	if err == nil {
		t.Error("expected error for empty credentials")
	}
}

func TestManagerConnectFactoryError(t *testing.T) {
	factory := newMockConnectionFactory()
	factory.createError = core.ErrNotConnected

	m, _ := NewManager(ManagerConfig{Factory: factory})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	creds := auth.UserPass("user", "pass")

	_, err := m.Connect(ctx, "tenant-1", creds)
	if err == nil {
		t.Error("expected error from factory")
	}
}

func TestManagerConnectWithCredentialManagerGetError(t *testing.T) {
	factory := newMockConnectionFactory()
	cm := auth.NewCredentialManager(auth.DefaultCredentialManagerConfig())
	// Don't register any credentials

	m, _ := NewManager(ManagerConfig{
		Factory:     factory,
		CredManager: cm,
	})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	_, err := m.ConnectWithCredentialManager(ctx, "unknown-tenant")
	if err == nil {
		t.Error("expected error for unknown tenant")
	}
}

func TestManagerDispatchErrorChannelFull(t *testing.T) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{
		Factory:         factory,
		ErrorBufferSize: 1, // Small buffer
	})
	defer m.Shutdown(context.Background())

	// Fill the error channel
	for i := 0; i < 10; i++ {
		testErr := core.TenantError{TenantID: "tenant-1", Err: core.ErrTenantNotFound}
		m.dispatchError(testErr)
	}
	// Should not panic
}

// ============== Benchmarks ==============

func BenchmarkRegistryRegister(b *testing.B) {
	r := NewRegistry()
	info := &Info{ID: "tenant-1", Name: "Test"}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r.tenants = make(map[string]*Info) // Reset
		r.Register(info)
	}
}

func BenchmarkRegistryGet(b *testing.B) {
	r := NewRegistry()
	r.Register(&Info{ID: "tenant-1", Name: "Test"})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r.Get("tenant-1")
	}
}

func BenchmarkRegistryList(b *testing.B) {
	r := NewRegistry()
	for i := 0; i < 100; i++ {
		r.Register(&Info{ID: "tenant-" + string(rune('a'+i%26))})
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r.List()
	}
}

func BenchmarkManagerConnect(b *testing.B) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{Factory: factory})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	creds := auth.UserPass("user", "pass")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		m.Connect(ctx, "tenant-1", creds)
	}
}

func BenchmarkManagerGet(b *testing.B) {
	factory := newMockConnectionFactory()
	m, _ := NewManager(ManagerConfig{Factory: factory})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	creds := auth.UserPass("user", "pass")
	m.Connect(ctx, "tenant-1", creds)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		m.Get("tenant-1")
	}
}
