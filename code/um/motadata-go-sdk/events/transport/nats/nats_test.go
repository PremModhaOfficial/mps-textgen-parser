package nats

import (
	"context"
	"sync"
	"testing"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/config"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

// ============== Transport Tests ==============

func TestNewTransport(t *testing.T) {
	tr := NewTransport()
	if tr == nil {
		t.Fatal("NewTransport returned nil")
	}
}

func TestTransportName(t *testing.T) {
	tr := NewTransport()
	if tr.Name() != TransportName {
		t.Errorf("expected %s, got %s", TransportName, tr.Name())
	}
	if tr.Name() != "nats" {
		t.Errorf("expected 'nats', got %s", tr.Name())
	}
}

func TestTransportPublisherWrongType(t *testing.T) {
	tr := NewTransport()
	// Pass wrong connection type
	_, err := tr.Publisher(&mockConnection{})
	if err == nil {
		t.Error("expected error for wrong connection type")
	}
}

func TestTransportSubscriberWrongType(t *testing.T) {
	tr := NewTransport()
	// Pass wrong connection type
	_, err := tr.Subscriber(&mockConnection{})
	if err == nil {
		t.Error("expected error for wrong connection type")
	}
}

// Mock connection that doesn't implement *Connection
type mockConnection struct{}

func (m *mockConnection) Connect(ctx context.Context) error { return nil }
func (m *mockConnection) Close(ctx context.Context) error   { return nil }
func (m *mockConnection) IsConnected() bool                 { return true }
func (m *mockConnection) State() core.ConnectionState       { return core.StateConnected }

// ============== TenantConnection Tests ==============

func TestNewTenantConnection(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")

	tc := NewTenantConnection("tenant-1", cfg, creds)
	if tc == nil {
		t.Fatal("NewTenantConnection returned nil")
	}
	if tc.Connection == nil {
		t.Error("Connection not set")
	}
	if tc.TenantID() != "tenant-1" {
		t.Errorf("unexpected tenant ID: %s", tc.TenantID())
	}
}

func TestTenantConnectionPublisher(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")

	tc := NewTenantConnection("tenant-1", cfg, creds)

	// First call creates publisher
	pub1 := tc.Publisher()
	if pub1 == nil {
		t.Error("Publisher returned nil")
	}

	// Second call returns same publisher
	pub2 := tc.Publisher()
	if pub1 != pub2 {
		t.Error("expected same publisher instance")
	}
}

func TestTenantConnectionSubscriber(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")

	tc := NewTenantConnection("tenant-1", cfg, creds)

	// First call creates subscriber
	sub1 := tc.Subscriber()
	if sub1 == nil {
		t.Error("Subscriber returned nil")
	}

	// Second call returns same subscriber
	sub2 := tc.Subscriber()
	if sub1 != sub2 {
		t.Error("expected same subscriber instance")
	}
}

// ============== Factory Tests ==============

func TestNewFactory(t *testing.T) {
	f := NewFactory()
	if f == nil {
		t.Fatal("NewFactory returned nil")
	}
}

func TestFactoryCreate(t *testing.T) {
	f := NewFactory()
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	ctx := context.Background()

	conn, err := f.Create(ctx, "tenant-1", cfg, creds)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if conn == nil {
		t.Error("Create returned nil connection")
	}
	if conn.TenantID() != "tenant-1" {
		t.Errorf("unexpected tenant ID: %s", conn.TenantID())
	}
}

// ============== Health Monitor Config Tests ==============

func TestDefaultHealthMonitorConfig(t *testing.T) {
	cfg := DefaultHealthMonitorConfig()
	if cfg.CheckInterval != 10*time.Second {
		t.Errorf("expected CheckInterval 10s, got %v", cfg.CheckInterval)
	}
	if cfg.RTTSamples != 10 {
		t.Errorf("expected RTTSamples 10, got %d", cfg.RTTSamples)
	}
}

func TestNewHealthMonitor(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	conn := NewConnection("tenant-1", cfg, creds)

	hm := NewHealthMonitor(conn, DefaultHealthMonitorConfig())
	if hm == nil {
		t.Fatal("NewHealthMonitor returned nil")
	}
	if hm.conn != conn {
		t.Error("connection not set correctly")
	}
	if hm.checkInterval != 10*time.Second {
		t.Errorf("unexpected checkInterval: %v", hm.checkInterval)
	}
}

func TestNewHealthMonitorDefaults(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	conn := NewConnection("tenant-1", cfg, creds)

	// Zero config should use defaults
	hm := NewHealthMonitor(conn, HealthMonitorConfig{})
	if hm.checkInterval != 10*time.Second {
		t.Errorf("expected default checkInterval, got %v", hm.checkInterval)
	}
	if hm.rttSamples != 10 {
		t.Errorf("expected default rttSamples, got %d", hm.rttSamples)
	}
}

func TestHealthMonitorStartStop(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	conn := NewConnection("tenant-1", cfg, creds)

	hmConfig := HealthMonitorConfig{
		CheckInterval: 10 * time.Millisecond, // Fast for testing
		RTTSamples:    5,
	}
	hm := NewHealthMonitor(conn, hmConfig)

	hm.Start()
	time.Sleep(25 * time.Millisecond)
	hm.Stop()
	// Should not panic
}

func TestHealthMonitorStats(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	conn := NewConnection("tenant-1", cfg, creds)

	hm := NewHealthMonitor(conn, DefaultHealthMonitorConfig())
	stats := hm.Stats()

	if stats.MessagesPublished != 0 {
		t.Error("expected MessagesPublished to be 0")
	}
}

func TestHealthMonitorIsHealthy(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	conn := NewConnection("tenant-1", cfg, creds)

	hm := NewHealthMonitor(conn, DefaultHealthMonitorConfig())
	// Initial state should be unhealthy (not connected)
	if hm.IsHealthy() {
		t.Error("expected unhealthy state initially")
	}
}

func TestHealthMonitorOnHealthChange(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	conn := NewConnection("tenant-1", cfg, creds)

	hm := NewHealthMonitor(conn, DefaultHealthMonitorConfig())

	var callbackCalled bool
	hm.OnHealthChange(func(healthy bool, status core.HealthStatus) {
		callbackCalled = true
	})

	// Verify callback is set (we can't easily test the callback being called
	// without a real connection)
	if hm.onHealthChange == nil {
		t.Error("callback not set")
	}
	_ = callbackCalled // Suppress unused warning
}

func TestHealthMonitorRecordPublish(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	conn := NewConnection("tenant-1", cfg, creds)

	hm := NewHealthMonitor(conn, DefaultHealthMonitorConfig())

	hm.RecordPublish(100)
	stats := hm.Stats()

	if stats.MessagesPublished != 1 {
		t.Errorf("expected MessagesPublished 1, got %d", stats.MessagesPublished)
	}
	if stats.BytesSent != 100 {
		t.Errorf("expected BytesSent 100, got %d", stats.BytesSent)
	}
}

func TestHealthMonitorRecordReceive(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	conn := NewConnection("tenant-1", cfg, creds)

	hm := NewHealthMonitor(conn, DefaultHealthMonitorConfig())

	hm.RecordReceive(200)
	stats := hm.Stats()

	if stats.MessagesReceived != 1 {
		t.Errorf("expected MessagesReceived 1, got %d", stats.MessagesReceived)
	}
	if stats.BytesReceived != 200 {
		t.Errorf("expected BytesReceived 200, got %d", stats.BytesReceived)
	}
}

func TestHealthMonitorRecordError(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	conn := NewConnection("tenant-1", cfg, creds)

	hm := NewHealthMonitor(conn, DefaultHealthMonitorConfig())

	hm.RecordError("publish")
	hm.RecordError("subscription")
	hm.RecordError("other")

	stats := hm.Stats()
	if stats.Errors != 3 {
		t.Errorf("expected Errors 3, got %d", stats.Errors)
	}
	if stats.PublishErrors != 1 {
		t.Errorf("expected PublishErrors 1, got %d", stats.PublishErrors)
	}
	if stats.SubscriptionErrors != 1 {
		t.Errorf("expected SubscriptionErrors 1, got %d", stats.SubscriptionErrors)
	}
}

func TestHealthMonitorRecordConnect(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	conn := NewConnection("tenant-1", cfg, creds)

	hm := NewHealthMonitor(conn, DefaultHealthMonitorConfig())

	hm.RecordConnect()
	stats := hm.Stats()

	if stats.ConnectedAt.IsZero() {
		t.Error("expected ConnectedAt to be set")
	}
}

func TestHealthMonitorRecordDisconnect(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	conn := NewConnection("tenant-1", cfg, creds)

	hm := NewHealthMonitor(conn, DefaultHealthMonitorConfig())

	hm.RecordDisconnect()
	stats := hm.Stats()

	if stats.LastDisconnect.IsZero() {
		t.Error("expected LastDisconnect to be set")
	}
}

func TestHealthMonitorRecordReconnect(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	conn := NewConnection("tenant-1", cfg, creds)

	hm := NewHealthMonitor(conn, DefaultHealthMonitorConfig())

	hm.RecordReconnect()
	stats := hm.Stats()

	if stats.LastReconnect.IsZero() {
		t.Error("expected LastReconnect to be set")
	}
	if stats.Reconnects != 1 {
		t.Errorf("expected Reconnects 1, got %d", stats.Reconnects)
	}
}

func TestHealthMonitorConcurrentRecords(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	conn := NewConnection("tenant-1", cfg, creds)

	hm := NewHealthMonitor(conn, DefaultHealthMonitorConfig())

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			hm.RecordPublish(10)
			hm.RecordReceive(20)
			hm.RecordError("publish")
			hm.Stats()
		}()
	}
	wg.Wait()

	stats := hm.Stats()
	if stats.MessagesPublished != 100 {
		t.Errorf("expected MessagesPublished 100, got %d", stats.MessagesPublished)
	}
}

// ============== SimpleHealthChecker Tests ==============

func TestNewHealthChecker(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	conn := NewConnection("tenant-1", cfg, creds)

	hc := NewHealthChecker(conn)
	if hc == nil {
		t.Fatal("NewHealthChecker returned nil")
	}
	if hc.conn != conn {
		t.Error("connection not set correctly")
	}
}

func TestHealthCheckerCheckNotConnected(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	conn := NewConnection("tenant-1", cfg, creds)

	hc := NewHealthChecker(conn)
	err := hc.Check(context.Background())
	if err != core.ErrNotConnected {
		t.Errorf("expected ErrNotConnected, got %v", err)
	}
}

// ============== Connection Tests (basic) ==============

func TestNewConnection(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")

	conn := NewConnection("tenant-1", cfg, creds)
	if conn == nil {
		t.Fatal("NewConnection returned nil")
	}
	if conn.TenantID() != "tenant-1" {
		t.Errorf("unexpected tenant ID: %s", conn.TenantID())
	}
}

func TestConnectionTenantID(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")

	conn := NewConnection("my-tenant", cfg, creds)
	if conn.TenantID() != "my-tenant" {
		t.Errorf("expected 'my-tenant', got %s", conn.TenantID())
	}
}

func TestConnectionState(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")

	conn := NewConnection("tenant-1", cfg, creds)
	state := conn.State()
	if state != core.StateDisconnected {
		t.Errorf("expected StateDisconnected, got %v", state)
	}
}

func TestConnectionIsConnected(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")

	conn := NewConnection("tenant-1", cfg, creds)
	if conn.IsConnected() {
		t.Error("expected IsConnected to return false")
	}
}

func TestConnectionHealth(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")

	conn := NewConnection("tenant-1", cfg, creds)
	health := conn.Health()

	if health.Healthy {
		t.Error("expected Healthy to be false")
	}
	if health.State != core.StateDisconnected {
		t.Errorf("expected StateDisconnected, got %v", health.State)
	}
}

func TestConnectionTouch(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")

	conn := NewConnection("tenant-1", cfg, creds)
	time.Sleep(10 * time.Millisecond)
	conn.Touch()

	// IsIdle should return false after Touch
	if conn.IsIdle(1 * time.Hour) {
		t.Error("expected not idle after Touch")
	}
}

func TestConnectionIsIdle(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")

	conn := NewConnection("tenant-1", cfg, creds)

	// With very long timeout, should not be idle
	if conn.IsIdle(1 * time.Hour) {
		t.Error("expected not idle with long timeout")
	}

	// Wait a bit and check with short timeout
	time.Sleep(20 * time.Millisecond)
	if !conn.IsIdle(10 * time.Millisecond) {
		t.Error("expected idle with short timeout")
	}
}

func TestConnectionConn(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")

	conn := NewConnection("tenant-1", cfg, creds)

	// Should be nil when not connected
	if conn.Conn() != nil {
		t.Error("expected Conn to return nil when not connected")
	}
}

// ============== Publisher Tests (basic) ==============

func TestNewPublisher(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	conn := NewConnection("tenant-1", cfg, creds)

	pub := NewPublisher(conn)
	if pub == nil {
		t.Fatal("NewPublisher returned nil")
	}
}

func TestPublisherPublishNotConnected(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	conn := NewConnection("tenant-1", cfg, creds)

	pub := NewPublisher(conn)
	msg := core.NewMessage([]byte("test"))
	err := pub.Publish(context.Background(), "test.subject", msg)
	if err == nil {
		t.Error("expected error when not connected")
	}
}

func TestPublisherRequestNotConnected(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	conn := NewConnection("tenant-1", cfg, creds)

	pub := NewPublisher(conn)
	msg := core.NewMessage([]byte("test"))
	_, err := pub.Request(context.Background(), "test.subject", msg)
	if err == nil {
		t.Error("expected error when not connected")
	}
}

// ============== Subscriber Tests (basic) ==============

func TestNewSubscriber(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	conn := NewConnection("tenant-1", cfg, creds)

	sub := NewSubscriber(conn)
	if sub == nil {
		t.Fatal("NewSubscriber returned nil")
	}
}

func TestSubscriberSubscribeNotConnected(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	conn := NewConnection("tenant-1", cfg, creds)

	sub := NewSubscriber(conn)
	handler := func(ctx context.Context, msg *core.Message) error { return nil }
	_, err := sub.Subscribe(context.Background(), "test.subject", handler)
	if err == nil {
		t.Error("expected error when not connected")
	}
}

func TestSubscriberQueueSubscribeNotConnected(t *testing.T) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	conn := NewConnection("tenant-1", cfg, creds)

	sub := NewSubscriber(conn)
	handler := func(ctx context.Context, msg *core.Message) error { return nil }
	_, err := sub.QueueSubscribe(context.Background(), "test.subject", "queue", handler)
	if err == nil {
		t.Error("expected error when not connected")
	}
}

// ============== Benchmarks ==============

func BenchmarkHealthMonitorRecordPublish(b *testing.B) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	conn := NewConnection("tenant-1", cfg, creds)
	hm := NewHealthMonitor(conn, DefaultHealthMonitorConfig())

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		hm.RecordPublish(100)
	}
}

func BenchmarkHealthMonitorStats(b *testing.B) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")
	conn := NewConnection("tenant-1", cfg, creds)
	hm := NewHealthMonitor(conn, DefaultHealthMonitorConfig())

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		hm.Stats()
	}
}

func BenchmarkNewConnection(b *testing.B) {
	cfg := config.DefaultConfig()
	creds := auth.UserPass("user", "pass")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		NewConnection("tenant-1", cfg, creds)
	}
}
