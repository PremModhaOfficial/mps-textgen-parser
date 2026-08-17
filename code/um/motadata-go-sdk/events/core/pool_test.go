package core

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ============================================================================
// Mock Connection for Testing
// ============================================================================

type mockPoolConnection struct {
	connected atomic.Bool
	state     ConnectionState
	mu        sync.RWMutex
	closeErr  error
}

func newMockPoolConnection() *mockPoolConnection {
	return &mockPoolConnection{
		state: StateDisconnected,
	}
}

func (c *mockPoolConnection) Connect(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connected.Store(true)
	c.state = StateConnected
	return nil
}

func (c *mockPoolConnection) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.connected.Store(false)
	c.state = StateClosed
	return c.closeErr
}

func (c *mockPoolConnection) IsConnected() bool {
	return c.connected.Load()
}

func (c *mockPoolConnection) State() ConnectionState {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.state
}

func (c *mockPoolConnection) Health() HealthStatus {
	return HealthStatus{
		State:   c.State(),
		Healthy: c.IsConnected(),
	}
}

// ============================================================================
// Pool Tests
// ============================================================================

func TestNewPool(t *testing.T) {
	factory := func(ctx context.Context) (Connection, error) {
		return newMockPoolConnection(), nil
	}

	t.Run("valid config", func(t *testing.T) {
		config := DefaultPoolConfig()
		config.Factory = factory

		pool, err := NewPool(config)
		if err != nil {
			t.Fatalf("NewPool() error = %v", err)
		}
		defer pool.Close(context.Background())

		if pool.Cap() != 10 {
			t.Errorf("Cap() = %d, want 10", pool.Cap())
		}
	})

	t.Run("missing factory", func(t *testing.T) {
		config := DefaultPoolConfig()

		_, err := NewPool(config)
		if err == nil {
			t.Error("NewPool() should fail without factory")
		}
	})

	t.Run("prewarm", func(t *testing.T) {
		config := DefaultPoolConfig()
		config.Factory = factory
		config.MinSize = 3
		config.PreWarm = true

		pool, err := NewPool(config)
		if err != nil {
			t.Fatalf("NewPool() error = %v", err)
		}
		defer pool.Close(context.Background())

		if pool.Len() != 3 {
			t.Errorf("Len() = %d, want 3 after prewarm", pool.Len())
		}
	})
}

func TestPoolGetPut(t *testing.T) {
	factory := func(ctx context.Context) (Connection, error) {
		return newMockPoolConnection(), nil
	}

	config := DefaultPoolConfig()
	config.Factory = factory
	config.MaxSize = 5

	pool, err := NewPool(config)
	if err != nil {
		t.Fatalf("NewPool() error = %v", err)
	}
	defer pool.Close(context.Background())

	t.Run("get creates connection", func(t *testing.T) {
		ctx := context.Background()
		conn, err := pool.Get(ctx)
		if err != nil {
			t.Fatalf("Get() error = %v", err)
		}

		if conn == nil {
			t.Fatal("Get() returned nil connection")
		}

		if !conn.IsConnected() {
			t.Error("Connection should be connected")
		}

		if err := pool.Put(conn); err != nil {
			t.Errorf("Put() error = %v", err)
		}
	})

	t.Run("put returns to available", func(t *testing.T) {
		ctx := context.Background()

		conn1, _ := pool.Get(ctx)
		stats := pool.Stats()
		if stats.InUse != 1 {
			t.Errorf("InUse = %d, want 1", stats.InUse)
		}

		pool.Put(conn1)
		stats = pool.Stats()
		if stats.Available != 1 {
			t.Errorf("Available = %d, want 1", stats.Available)
		}
	})

	t.Run("reuses connections", func(t *testing.T) {
		ctx := context.Background()

		conn1, _ := pool.Get(ctx)
		id1 := conn1.ID()
		pool.Put(conn1)

		conn2, _ := pool.Get(ctx)
		id2 := conn2.ID()
		pool.Put(conn2)

		if id1 != id2 {
			t.Errorf("Connection not reused: id1=%s, id2=%s", id1, id2)
		}
	})
}

func TestPoolMaxSize(t *testing.T) {
	factory := func(ctx context.Context) (Connection, error) {
		return newMockPoolConnection(), nil
	}

	config := DefaultPoolConfig()
	config.Factory = factory
	config.MaxSize = 2
	config.WaitForConnection = false

	pool, err := NewPool(config)
	if err != nil {
		t.Fatalf("NewPool() error = %v", err)
	}
	defer pool.Close(context.Background())

	ctx := context.Background()

	conn1, _ := pool.Get(ctx)
	conn2, _ := pool.Get(ctx)

	_, err = pool.Get(ctx)
	if !errors.Is(err, ErrPoolExhausted) {
		t.Errorf("Get() error = %v, want ErrPoolExhausted", err)
	}

	pool.Put(conn1)
	pool.Put(conn2)
}

func TestPoolWaitForConnection(t *testing.T) {
	factory := func(ctx context.Context) (Connection, error) {
		return newMockPoolConnection(), nil
	}

	config := DefaultPoolConfig()
	config.Factory = factory
	config.MaxSize = 1
	config.WaitForConnection = true

	pool, err := NewPool(config)
	if err != nil {
		t.Fatalf("NewPool() error = %v", err)
	}
	defer pool.Close(context.Background())

	ctx := context.Background()

	conn1, _ := pool.Get(ctx)

	var conn2 PooledConnection
	var getErr error
	done := make(chan struct{})

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		conn2, getErr = pool.Get(ctx)
		close(done)
	}()

	time.Sleep(20 * time.Millisecond)
	pool.Put(conn1)

	<-done

	if getErr != nil {
		t.Errorf("Get() error = %v", getErr)
	}
	if conn2 == nil {
		t.Error("Get() returned nil after wait")
	}

	if conn2 != nil {
		pool.Put(conn2)
	}
}

func TestPoolClose(t *testing.T) {
	factory := func(ctx context.Context) (Connection, error) {
		return newMockPoolConnection(), nil
	}

	config := DefaultPoolConfig()
	config.Factory = factory

	pool, err := NewPool(config)
	if err != nil {
		t.Fatalf("NewPool() error = %v", err)
	}

	ctx := context.Background()

	conn, _ := pool.Get(ctx)
	pool.Put(conn)

	if err := pool.Close(ctx); err != nil {
		t.Errorf("Close() error = %v", err)
	}

	_, err = pool.Get(ctx)
	if !errors.Is(err, ErrPoolClosed) {
		t.Errorf("Get() after close error = %v, want ErrPoolClosed", err)
	}
}

func TestPoolStats(t *testing.T) {
	factory := func(ctx context.Context) (Connection, error) {
		return newMockPoolConnection(), nil
	}

	config := DefaultPoolConfig()
	config.Factory = factory
	config.MaxSize = 5

	pool, err := NewPool(config)
	if err != nil {
		t.Fatalf("NewPool() error = %v", err)
	}
	defer pool.Close(context.Background())

	ctx := context.Background()

	conn1, _ := pool.Get(ctx)
	conn2, _ := pool.Get(ctx)

	stats := pool.Stats()

	if stats.TotalCreated != 2 {
		t.Errorf("TotalCreated = %d, want 2", stats.TotalCreated)
	}
	if stats.TotalAcquired != 2 {
		t.Errorf("TotalAcquired = %d, want 2", stats.TotalAcquired)
	}
	if stats.InUse != 2 {
		t.Errorf("InUse = %d, want 2", stats.InUse)
	}
	if stats.MaxSize != 5 {
		t.Errorf("MaxSize = %d, want 5", stats.MaxSize)
	}

	pool.Put(conn1)
	pool.Put(conn2)

	stats = pool.Stats()
	if stats.TotalReleased != 2 {
		t.Errorf("TotalReleased = %d, want 2", stats.TotalReleased)
	}
}

func TestPoolConcurrency(t *testing.T) {
	factory := func(ctx context.Context) (Connection, error) {
		return newMockPoolConnection(), nil
	}

	config := DefaultPoolConfig()
	config.Factory = factory
	config.MaxSize = 10

	pool, err := NewPool(config)
	if err != nil {
		t.Fatalf("NewPool() error = %v", err)
	}
	defer pool.Close(context.Background())

	var wg sync.WaitGroup
	goroutines := 50
	iterations := 100

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				ctx := context.Background()
				conn, err := pool.Get(ctx)
				if err != nil {
					continue
				}
				time.Sleep(time.Microsecond)
				pool.Put(conn)
			}
		}()
	}

	wg.Wait()

	stats := pool.Stats()
	if stats.TotalAcquired == 0 {
		t.Error("TotalAcquired should be > 0")
	}
	if stats.TotalReleased == 0 {
		t.Error("TotalReleased should be > 0")
	}
}

func TestPoolHealthChecker(t *testing.T) {
	healthyConn := newMockPoolConnection()
	unhealthyConn := newMockPoolConnection()

	callCount := atomic.Int32{}
	factory := func(ctx context.Context) (Connection, error) {
		count := callCount.Add(1)
		if count == 1 {
			return unhealthyConn, nil
		}
		return healthyConn, nil
	}

	config := DefaultPoolConfig()
	config.Factory = factory
	config.HealthCheckInterval = 50 * time.Millisecond
	config.HealthCheckTimeout = 10 * time.Millisecond

	pool, err := NewPool(config)
	if err != nil {
		t.Fatalf("NewPool() error = %v", err)
	}
	defer pool.Close(context.Background())

	ctx := context.Background()

	conn, _ := pool.Get(ctx)
	pool.Put(conn)

	unhealthyConn.connected.Store(false)

	time.Sleep(100 * time.Millisecond)

	stats := pool.Stats()
	if stats.TotalHealthChecksFailed == 0 {
		t.Error("TotalHealthChecksFailed should be > 0")
	}
}

func TestPooledConnectionUseCount(t *testing.T) {
	factory := func(ctx context.Context) (Connection, error) {
		return newMockPoolConnection(), nil
	}

	config := DefaultPoolConfig()
	config.Factory = factory

	pool, err := NewPool(config)
	if err != nil {
		t.Fatalf("NewPool() error = %v", err)
	}
	defer pool.Close(context.Background())

	ctx := context.Background()

	conn, _ := pool.Get(ctx)
	id := conn.ID()
	if conn.UseCount() != 1 {
		t.Errorf("UseCount() = %d, want 1", conn.UseCount())
	}
	pool.Put(conn)

	conn2, _ := pool.Get(ctx)
	if conn2.ID() == id && conn2.UseCount() != 2 {
		t.Errorf("UseCount() = %d, want 2", conn2.UseCount())
	}
	pool.Put(conn2)
}
