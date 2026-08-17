package core

import (
	"container/heap"
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// ============================================================================
// Connection Pool - Production-Level Implementation
// ============================================================================

// Pool manages a pool of connections with health checking, load balancing,
// and automatic recovery.
type Pool interface {
	// Get returns a healthy connection from the pool
	Get(ctx context.Context) (PooledConnection, error)

	// Put returns a connection to the pool
	Put(conn PooledConnection) error

	// Close closes all connections in the pool
	Close(ctx context.Context) error

	// Stats returns pool statistics
	Stats() PoolStats

	// Len returns the number of connections in the pool
	Len() int

	// Cap returns the pool capacity
	Cap() int

	// SetHealthChecker sets a custom health checker
	SetHealthChecker(checker PoolHealthChecker)
}

// PooledConnection represents a connection from the pool
type PooledConnection interface {
	Connection
	HealthChecker

	// ID returns the connection identifier
	ID() string

	// CreatedAt returns when the connection was created
	CreatedAt() time.Time

	// LastUsedAt returns when the connection was last used
	LastUsedAt() time.Time

	// UseCount returns how many times the connection has been used
	UseCount() int64

	// MarkUsed marks the connection as used
	MarkUsed()

	// Unwrap returns the underlying connection
	Unwrap() Connection
}

// PoolHealthChecker checks connection health
type PoolHealthChecker func(ctx context.Context, conn PooledConnection) bool

// PoolFactory creates new connections for the pool
type PoolFactory func(ctx context.Context) (Connection, error)

// ============================================================================
// Pool Configuration
// ============================================================================

// PoolConfig configures the connection pool
type PoolConfig struct {
	// MinSize is the minimum number of connections (default: 1)
	MinSize int

	// MaxSize is the maximum number of connections (default: 10)
	MaxSize int

	// MaxIdleTime is how long a connection can be idle before eviction
	MaxIdleTime time.Duration

	// MaxLifetime is the maximum lifetime of a connection
	MaxLifetime time.Duration

	// AcquireTimeout is the timeout for acquiring a connection
	AcquireTimeout time.Duration

	// HealthCheckInterval is how often to check connection health
	HealthCheckInterval time.Duration

	// HealthCheckTimeout is the timeout for health checks
	HealthCheckTimeout time.Duration

	// WaitForConnection determines if Get() waits for a connection
	WaitForConnection bool

	// PreWarm creates MinSize connections on initialization
	PreWarm bool

	// Factory creates new connections
	Factory PoolFactory

	// OnAcquire is called when a connection is acquired
	OnAcquire func(conn PooledConnection)

	// OnRelease is called when a connection is released
	OnRelease func(conn PooledConnection)

	// OnClose is called when a connection is closed
	OnClose func(conn PooledConnection)
}

// DefaultPoolConfig returns default pool configuration
func DefaultPoolConfig() PoolConfig {
	return PoolConfig{
		MinSize:             1,
		MaxSize:             10,
		MaxIdleTime:         30 * time.Minute,
		MaxLifetime:         0,
		AcquireTimeout:      30 * time.Second,
		HealthCheckInterval: 30 * time.Second,
		HealthCheckTimeout:  5 * time.Second,
		WaitForConnection:   true,
		PreWarm:             false,
	}
}

// ============================================================================
// Pool Statistics
// ============================================================================

// PoolStats contains pool statistics
type PoolStats struct {
	Size                    int           `json:"size"`
	Available               int           `json:"available"`
	InUse                   int           `json:"in_use"`
	MaxSize                 int           `json:"max_size"`
	MinSize                 int           `json:"min_size"`
	TotalCreated            int64         `json:"total_created"`
	TotalClosed             int64         `json:"total_closed"`
	TotalAcquired           int64         `json:"total_acquired"`
	TotalReleased           int64         `json:"total_released"`
	TotalFailed             int64         `json:"total_failed"`
	TotalHealthChecksFailed int64         `json:"total_health_checks_failed"`
	AverageAcquireTime      time.Duration `json:"average_acquire_time"`
	AverageUseTime          time.Duration `json:"average_use_time"`
	WaitingRequests         int           `json:"waiting_requests"`
}

// ============================================================================
// Connection Pool Implementation
// ============================================================================

type connectionPool struct {
	mu sync.Mutex

	config        PoolConfig
	factory       PoolFactory
	healthChecker PoolHealthChecker

	connections []*pooledConn
	available   *connHeap
	inUse       map[string]*pooledConn

	waiters     []*waiter
	waitersLock sync.Mutex

	stats poolStatistics

	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	closed    atomic.Bool
	closeOnce sync.Once
}

type poolStatistics struct {
	totalCreated            atomic.Int64
	totalClosed             atomic.Int64
	totalAcquired           atomic.Int64
	totalReleased           atomic.Int64
	totalFailed             atomic.Int64
	totalHealthChecksFailed atomic.Int64
	totalAcquireTime        atomic.Int64
	totalUseTime            atomic.Int64
}

type waiter struct {
	ch       chan *pooledConn
	ctx      context.Context
	deadline time.Time
}

// NewPool creates a new connection pool
func NewPool(config PoolConfig) (Pool, error) {
	if config.Factory == nil {
		return nil, errors.New("pool factory is required")
	}

	if config.MinSize < 0 {
		config.MinSize = 0
	}
	if config.MaxSize <= 0 {
		config.MaxSize = 10
	}
	if config.MinSize > config.MaxSize {
		config.MinSize = config.MaxSize
	}
	if config.AcquireTimeout <= 0 {
		config.AcquireTimeout = 30 * time.Second
	}
	if config.HealthCheckInterval <= 0 {
		config.HealthCheckInterval = 30 * time.Second
	}
	if config.HealthCheckTimeout <= 0 {
		config.HealthCheckTimeout = 5 * time.Second
	}

	ctx, cancel := context.WithCancel(context.Background())

	p := &connectionPool{
		config:      config,
		factory:     config.Factory,
		connections: make([]*pooledConn, 0, config.MaxSize),
		available:   &connHeap{},
		inUse:       make(map[string]*pooledConn),
		waiters:     make([]*waiter, 0),
		ctx:         ctx,
		cancel:      cancel,
	}

	heap.Init(p.available)
	p.healthChecker = p.defaultHealthChecker

	if config.PreWarm {
		if err := p.preWarm(ctx); err != nil {
			cancel()
			return nil, err
		}
	}

	p.startMaintenance()

	return p, nil
}

func (p *connectionPool) preWarm(ctx context.Context) error {
	for i := 0; i < p.config.MinSize; i++ {
		conn, err := p.createConnection(ctx)
		if err != nil {
			return err
		}
		p.mu.Lock()
		p.connections = append(p.connections, conn)
		heap.Push(p.available, conn)
		p.mu.Unlock()
	}
	return nil
}

func (p *connectionPool) startMaintenance() {
	if p.config.HealthCheckInterval > 0 {
		p.wg.Add(1)
		go p.healthCheckLoop()
	}

	if p.config.MaxIdleTime > 0 {
		p.wg.Add(1)
		go p.evictionLoop()
	}

	if p.config.MinSize > 0 {
		p.wg.Add(1)
		go p.minSizeLoop()
	}
}

func (p *connectionPool) healthCheckLoop() {
	defer p.wg.Done()

	ticker := time.NewTicker(p.config.HealthCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			p.performHealthChecks()
		}
	}
}

func (p *connectionPool) performHealthChecks() {
	p.mu.Lock()
	toCheck := make([]*pooledConn, 0, p.available.Len())
	for _, conn := range *p.available {
		toCheck = append(toCheck, conn)
	}
	p.mu.Unlock()

	for _, conn := range toCheck {
		ctx, cancel := context.WithTimeout(p.ctx, p.config.HealthCheckTimeout)
		if !p.healthChecker(ctx, conn) {
			p.stats.totalHealthChecksFailed.Add(1)
			p.mu.Lock()
			p.removeConnection(conn)
			p.mu.Unlock()
			conn.Close(ctx)
		}
		cancel()
	}
}

func (p *connectionPool) evictionLoop() {
	defer p.wg.Done()

	ticker := time.NewTicker(p.config.MaxIdleTime / 2)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			p.evictIdle()
		}
	}
}

func (p *connectionPool) evictIdle() {
	now := time.Now()

	p.mu.Lock()
	var toClose []*pooledConn

	newAvailable := &connHeap{}
	heap.Init(newAvailable)

	for p.available.Len() > 0 {
		conn := heap.Pop(p.available).(*pooledConn)

		isIdle := now.Sub(conn.lastUsed) > p.config.MaxIdleTime
		isExpired := p.config.MaxLifetime > 0 && now.Sub(conn.created) > p.config.MaxLifetime
		overMin := len(p.connections)-len(toClose) > p.config.MinSize

		if (isIdle || isExpired) && overMin {
			toClose = append(toClose, conn)
		} else {
			heap.Push(newAvailable, conn)
		}
	}

	p.available = newAvailable

	for _, conn := range toClose {
		p.removeFromSlice(conn)
	}

	p.mu.Unlock()

	for _, conn := range toClose {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		conn.Close(ctx)
		cancel()
		p.stats.totalClosed.Add(1)
	}
}

func (p *connectionPool) minSizeLoop() {
	defer p.wg.Done()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			return
		case <-ticker.C:
			p.maintainMinSize()
		}
	}
}

func (p *connectionPool) maintainMinSize() {
	p.mu.Lock()
	needed := p.config.MinSize - len(p.connections)
	p.mu.Unlock()

	for i := 0; i < needed; i++ {
		ctx, cancel := context.WithTimeout(p.ctx, 10*time.Second)
		conn, err := p.createConnection(ctx)
		cancel()

		if err != nil {
			break
		}

		p.mu.Lock()
		p.connections = append(p.connections, conn)
		heap.Push(p.available, conn)
		p.mu.Unlock()
	}
}

func (p *connectionPool) Get(ctx context.Context) (PooledConnection, error) {
	if p.closed.Load() {
		return nil, ErrPoolClosed
	}

	start := time.Now()
	defer func() {
		p.stats.totalAcquireTime.Add(int64(time.Since(start)))
	}()

	p.mu.Lock()
	for p.available.Len() > 0 {
		conn := heap.Pop(p.available).(*pooledConn)

		if conn.IsConnected() {
			conn.MarkUsed()
			p.inUse[conn.id] = conn
			p.mu.Unlock()

			p.stats.totalAcquired.Add(1)
			if p.config.OnAcquire != nil {
				p.config.OnAcquire(conn)
			}
			return conn, nil
		}

		p.removeFromSlice(conn)
		go func(c *pooledConn) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			c.Close(ctx)
			cancel()
		}(conn)
	}

	if len(p.connections) < p.config.MaxSize {
		p.mu.Unlock()

		conn, err := p.createConnection(ctx)
		if err != nil {
			p.stats.totalFailed.Add(1)
			return nil, err
		}

		p.mu.Lock()
		p.connections = append(p.connections, conn)
		conn.MarkUsed()
		p.inUse[conn.id] = conn
		p.mu.Unlock()

		p.stats.totalAcquired.Add(1)
		if p.config.OnAcquire != nil {
			p.config.OnAcquire(conn)
		}
		return conn, nil
	}

	if !p.config.WaitForConnection {
		p.mu.Unlock()
		p.stats.totalFailed.Add(1)
		return nil, ErrPoolExhausted
	}

	w := &waiter{
		ch:  make(chan *pooledConn, 1),
		ctx: ctx,
	}

	if deadline, ok := ctx.Deadline(); ok {
		w.deadline = deadline
	} else {
		w.deadline = time.Now().Add(p.config.AcquireTimeout)
	}

	p.waitersLock.Lock()
	p.waiters = append(p.waiters, w)
	p.waitersLock.Unlock()

	p.mu.Unlock()

	select {
	case <-ctx.Done():
		p.removeWaiter(w)
		p.stats.totalFailed.Add(1)
		return nil, ctx.Err()
	case conn := <-w.ch:
		if conn == nil {
			p.stats.totalFailed.Add(1)
			return nil, ErrPoolExhausted
		}
		p.stats.totalAcquired.Add(1)
		if p.config.OnAcquire != nil {
			p.config.OnAcquire(conn)
		}
		return conn, nil
	}
}

func (p *connectionPool) Put(conn PooledConnection) error {
	if p.closed.Load() {
		return conn.Close(context.Background())
	}

	pc, ok := conn.(*pooledConn)
	if !ok {
		return errors.New("invalid connection type")
	}

	if p.config.OnRelease != nil {
		p.config.OnRelease(conn)
	}

	useTime := time.Since(pc.lastUsed)
	p.stats.totalUseTime.Add(int64(useTime))
	p.stats.totalReleased.Add(1)

	if !pc.IsConnected() {
		p.mu.Lock()
		delete(p.inUse, pc.id)
		p.removeFromSlice(pc)
		p.mu.Unlock()

		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			pc.Close(ctx)
			cancel()
		}()

		p.stats.totalClosed.Add(1)
		return nil
	}

	p.mu.Lock()
	delete(p.inUse, pc.id)

	p.waitersLock.Lock()
	for len(p.waiters) > 0 {
		w := p.waiters[0]
		p.waiters = p.waiters[1:]

		if time.Now().Before(w.deadline) {
			pc.MarkUsed()
			p.inUse[pc.id] = pc
			p.waitersLock.Unlock()
			p.mu.Unlock()

			select {
			case w.ch <- pc:
				return nil
			default:
				p.mu.Lock()
				delete(p.inUse, pc.id)
				p.waitersLock.Lock()
			}
		}
	}
	p.waitersLock.Unlock()

	heap.Push(p.available, pc)
	p.mu.Unlock()

	return nil
}

func (p *connectionPool) Close(ctx context.Context) error {
	var err error

	p.closeOnce.Do(func() {
		p.closed.Store(true)
		p.cancel()

		p.waitersLock.Lock()
		for _, w := range p.waiters {
			close(w.ch)
		}
		p.waiters = nil
		p.waitersLock.Unlock()

		p.mu.Lock()
		connections := make([]*pooledConn, len(p.connections))
		copy(connections, p.connections)
		p.connections = nil
		p.available = &connHeap{}
		p.inUse = make(map[string]*pooledConn)
		p.mu.Unlock()

		var wg sync.WaitGroup
		for _, conn := range connections {
			wg.Add(1)
			go func(c *pooledConn) {
				defer wg.Done()
				if p.config.OnClose != nil {
					p.config.OnClose(c)
				}
				c.Close(ctx)
			}(conn)
		}

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-ctx.Done():
			err = ctx.Err()
		case <-done:
		}

		p.wg.Wait()
	})

	return err
}

func (p *connectionPool) Stats() PoolStats {
	p.mu.Lock()
	defer p.mu.Unlock()

	acquired := p.stats.totalAcquired.Load()
	var avgAcquire, avgUse time.Duration
	if acquired > 0 {
		avgAcquire = time.Duration(p.stats.totalAcquireTime.Load() / acquired)
		avgUse = time.Duration(p.stats.totalUseTime.Load() / acquired)
	}

	p.waitersLock.Lock()
	waiters := len(p.waiters)
	p.waitersLock.Unlock()

	return PoolStats{
		Size:                    len(p.connections),
		Available:               p.available.Len(),
		InUse:                   len(p.inUse),
		MaxSize:                 p.config.MaxSize,
		MinSize:                 p.config.MinSize,
		TotalCreated:            p.stats.totalCreated.Load(),
		TotalClosed:             p.stats.totalClosed.Load(),
		TotalAcquired:           p.stats.totalAcquired.Load(),
		TotalReleased:           p.stats.totalReleased.Load(),
		TotalFailed:             p.stats.totalFailed.Load(),
		TotalHealthChecksFailed: p.stats.totalHealthChecksFailed.Load(),
		AverageAcquireTime:      avgAcquire,
		AverageUseTime:          avgUse,
		WaitingRequests:         waiters,
	}
}

func (p *connectionPool) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.connections)
}

func (p *connectionPool) Cap() int {
	return p.config.MaxSize
}

func (p *connectionPool) SetHealthChecker(checker PoolHealthChecker) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.healthChecker = checker
}

func (p *connectionPool) createConnection(ctx context.Context) (*pooledConn, error) {
	conn, err := p.factory(ctx)
	if err != nil {
		return nil, err
	}

	if err := conn.Connect(ctx); err != nil {
		return nil, err
	}

	pc := newPooledConn(conn)
	p.stats.totalCreated.Add(1)

	return pc, nil
}

func (p *connectionPool) defaultHealthChecker(ctx context.Context, conn PooledConnection) bool {
	return conn.IsConnected()
}

func (p *connectionPool) removeConnection(conn *pooledConn) {
	for i := 0; i < p.available.Len(); i++ {
		if (*p.available)[i] == conn {
			heap.Remove(p.available, i)
			break
		}
	}
	p.removeFromSlice(conn)
	delete(p.inUse, conn.id)
}

func (p *connectionPool) removeFromSlice(conn *pooledConn) {
	for i, c := range p.connections {
		if c == conn {
			p.connections = append(p.connections[:i], p.connections[i+1:]...)
			break
		}
	}
}

func (p *connectionPool) removeWaiter(w *waiter) {
	p.waitersLock.Lock()
	defer p.waitersLock.Unlock()

	for i, waiter := range p.waiters {
		if waiter == w {
			p.waiters = append(p.waiters[:i], p.waiters[i+1:]...)
			break
		}
	}
}

// ============================================================================
// Pooled Connection Wrapper
// ============================================================================

type pooledConn struct {
	conn Connection

	id       string
	created  time.Time
	lastUsed time.Time
	useCount atomic.Int64

	mu sync.RWMutex
}

var connCounter atomic.Int64

func newPooledConn(conn Connection) *pooledConn {
	id := connCounter.Add(1)
	now := time.Now()
	return &pooledConn{
		conn:     conn,
		id:       fmt.Sprintf("conn-%d", id),
		created:  now,
		lastUsed: now,
	}
}

func (c *pooledConn) Connect(ctx context.Context) error {
	return c.conn.Connect(ctx)
}

func (c *pooledConn) Close(ctx context.Context) error {
	return c.conn.Close(ctx)
}

func (c *pooledConn) IsConnected() bool {
	return c.conn.IsConnected()
}

func (c *pooledConn) State() ConnectionState {
	return c.conn.State()
}

func (c *pooledConn) Health() HealthStatus {
	if checker, ok := c.conn.(HealthChecker); ok {
		return checker.Health()
	}
	return HealthStatus{
		State:   c.conn.State(),
		Healthy: c.conn.IsConnected(),
	}
}

func (c *pooledConn) ID() string {
	return c.id
}

func (c *pooledConn) CreatedAt() time.Time {
	return c.created
}

func (c *pooledConn) LastUsedAt() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastUsed
}

func (c *pooledConn) UseCount() int64 {
	return c.useCount.Load()
}

func (c *pooledConn) MarkUsed() {
	c.mu.Lock()
	c.lastUsed = time.Now()
	c.mu.Unlock()
	c.useCount.Add(1)
}

func (c *pooledConn) Unwrap() Connection {
	return c.conn
}

// ============================================================================
// Connection Heap (LRU ordering)
// ============================================================================

type connHeap []*pooledConn

func (h connHeap) Len() int { return len(h) }

func (h connHeap) Less(i, j int) bool {
	return h[i].lastUsed.Before(h[j].lastUsed)
}

func (h connHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *connHeap) Push(x any) {
	*h = append(*h, x.(*pooledConn))
}

func (h *connHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// ============================================================================
// Pool Errors
// ============================================================================

var (
	ErrPoolExhausted = errors.New("connection pool exhausted")
	ErrPoolClosed    = errors.New("connection pool closed")
)
