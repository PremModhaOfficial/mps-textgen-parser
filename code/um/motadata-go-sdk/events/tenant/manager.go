package tenant

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/config"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer"
)

// Connection represents a tenant's connection
type Connection interface {
	core.Connection
	core.HealthChecker

	// TenantID returns the tenant identifier
	TenantID() string

	// Touch updates the last activity timestamp
	Touch()

	// IsIdle returns true if connection has been idle
	IsIdle(timeout time.Duration) bool

	// Publisher returns a publisher for this connection
	Publisher() core.Publisher

	// Subscriber returns a subscriber for this connection
	Subscriber() core.Subscriber
}

// ConnectionFactory creates connections for tenants
type ConnectionFactory interface {
	// Create creates a new connection for a tenant
	Create(ctx context.Context, tenantID string, cfg *config.Config, creds *auth.Credentials) (Connection, error)
}

// Manager manages connections for multiple tenants
type Manager struct {
	config  *config.Config
	factory ConnectionFactory

	mu          sync.RWMutex
	connections map[string]Connection

	// Credential management
	credManager *auth.CredentialManager

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Error handling
	errorChan   chan core.TenantError
	errorCb     func(core.TenantError)
	errorMu     sync.RWMutex
	errorClosed atomic.Bool // guards against send on closed channel
}

// ManagerConfig holds manager configuration
type ManagerConfig struct {
	Config          *config.Config
	Factory         ConnectionFactory
	CredManager     *auth.CredentialManager
	ErrorBufferSize int
}

// NewManager creates a new tenant manager
func NewManager(cfg ManagerConfig) (*Manager, error) {
	if cfg.Config == nil {
		cfg.Config = config.DefaultConfig()
	}

	if err := config.Validate(cfg.Config); err != nil {
		return nil, err
	}

	if cfg.ErrorBufferSize == 0 {
		cfg.ErrorBufferSize = 100
	}

	ctx, cancel := context.WithCancel(context.Background())

	m := &Manager{
		config:      cfg.Config,
		factory:     cfg.Factory,
		connections: make(map[string]Connection),
		credManager: cfg.CredManager,
		ctx:         ctx,
		cancel:      cancel,
		errorChan:   make(chan core.TenantError, cfg.ErrorBufferSize),
	}

	// Start cleanup loop
	m.startCleanupLoop()

	return m, nil
}

// startCleanupLoop starts the background cleanup goroutine
func (m *Manager) startCleanupLoop() {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()

		ticker := time.NewTicker(m.config.CleanupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-m.ctx.Done():
				return
			case <-ticker.C:
				m.cleanupIdleConnections()
			}
		}
	}()
}

// cleanupIdleConnections closes idle connections
func (m *Manager) cleanupIdleConnections() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for tenantID, conn := range m.connections {
		if conn.IsIdle(m.config.IdleTimeout) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			conn.Close(ctx)
			cancel()
			delete(m.connections, tenantID)

			logger.Info(ctx, "closed idle tenant connection",
				logger.String("tenant_id", tenantID),
				logger.Duration("idle_timeout", m.config.IdleTimeout),
			)
		}
	}
}

// Connect establishes a connection for a tenant
func (m *Manager) Connect(ctx context.Context, tenantID string, creds *auth.Credentials) (Connection, error) {
	// Start tracing span
	ctx, span := tracer.Start(ctx, "tenant.connect",
		tracer.WithAttributes(
			tracer.StringAttr("tenant.id", tenantID),
		),
	)
	defer span.End()

	// Check for shutdown
	select {
	case <-m.ctx.Done():
		span.SetError(core.ErrShutdownInProgress)
		return nil, core.ErrShutdownInProgress
	default:
	}

	// Check if already connected
	m.mu.Lock()
	if existing, ok := m.connections[tenantID]; ok {
		existing.Touch()
		m.mu.Unlock()
		span.SetAttributes(tracer.BoolAttr("connection.reused", true))
		span.SetOK()
		logger.Debug(ctx, "reusing existing tenant connection",
			logger.String("tenant_id", tenantID),
		)
		return existing, nil
	}
	m.mu.Unlock()

	// Get credentials if not provided
	if creds == nil && m.credManager != nil {
		var err error
		creds, err = m.credManager.GetCredentials(tenantID)
		if err != nil {
			span.SetError(err)
			logger.Error(ctx, "failed to get credentials for tenant",
				logger.String("tenant_id", tenantID),
				logger.Err(err),
			)
			return nil, err
		}
	}

	if creds == nil || creds.IsEmpty() {
		span.SetError(core.ErrInvalidCredential)
		logger.Error(ctx, "invalid or missing credentials for tenant",
			logger.String("tenant_id", tenantID),
		)
		return nil, core.ErrInvalidCredential
	}

	// Create connection
	if m.factory == nil {
		err := core.NewError("connection", "create", core.ErrNotConnected)
		span.SetError(err)
		return nil, err
	}

	conn, err := m.factory.Create(ctx, tenantID, m.config, creds)
	if err != nil {
		span.SetError(err)
		logger.Error(ctx, "failed to create tenant connection",
			logger.String("tenant_id", tenantID),
			logger.Err(err),
		)
		return nil, err
	}

	// Connect
	if err := conn.Connect(ctx); err != nil {
		span.SetError(err)
		logger.Error(ctx, "failed to establish tenant connection",
			logger.String("tenant_id", tenantID),
			logger.Err(err),
		)
		return nil, err
	}

	// Store connection
	m.mu.Lock()
	// Double-check for race condition
	if existing, ok := m.connections[tenantID]; ok {
		m.mu.Unlock()
		conn.Close(ctx) // Close the new connection
		existing.Touch()
		span.SetAttributes(tracer.BoolAttr("connection.reused", true))
		span.SetOK()
		return existing, nil
	}
	m.connections[tenantID] = conn
	m.mu.Unlock()

	span.SetAttributes(tracer.BoolAttr("connection.new", true))
	span.SetOK()
	logger.Info(ctx, "tenant connection established",
		logger.String("tenant_id", tenantID),
	)

	return conn, nil
}

// ConnectWithCredentialManager connects using the credential manager
func (m *Manager) ConnectWithCredentialManager(ctx context.Context, tenantID string) (Connection, error) {
	if m.credManager == nil {
		return nil, core.NewError("connection", "connect", auth.ErrInvalidCredential)
	}

	creds, err := m.credManager.GetCredentials(tenantID)
	if err != nil {
		return nil, err
	}

	return m.Connect(ctx, tenantID, creds)
}

// ConnectWithPayload connects using a registration payload
func (m *Manager) ConnectWithPayload(ctx context.Context, payload auth.RegistrationPayload) (Connection, error) {
	if err := payload.Validate(); err != nil {
		return nil, err
	}

	// Register credentials if we have a credential manager
	if m.credManager != nil {
		if err := m.credManager.RegisterFromPayload(payload); err != nil {
			return nil, err
		}
	}

	creds := auth.JWT(payload.JWT, payload.Seed)
	return m.Connect(ctx, payload.TenantID, creds)
}

// Get returns an existing connection
func (m *Manager) Get(tenantID string) (Connection, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	conn, ok := m.connections[tenantID]
	if ok {
		conn.Touch()
	}
	return conn, ok
}

// Disconnect closes a tenant's connection
func (m *Manager) Disconnect(ctx context.Context, tenantID string) error {
	ctx, span := tracer.Start(ctx, "tenant.disconnect",
		tracer.WithAttributes(
			tracer.StringAttr("tenant.id", tenantID),
		),
	)
	defer span.End()

	m.mu.Lock()
	conn, ok := m.connections[tenantID]
	if !ok {
		m.mu.Unlock()
		span.SetError(core.ErrTenantNotFound)
		return core.ErrTenantNotFound
	}
	delete(m.connections, tenantID)
	m.mu.Unlock()

	err := conn.Close(ctx)
	if err != nil {
		span.SetError(err)
		logger.Error(ctx, "failed to close tenant connection",
			logger.String("tenant_id", tenantID),
			logger.Err(err),
		)
	} else {
		span.SetOK()
		logger.Info(ctx, "tenant connection closed",
			logger.String("tenant_id", tenantID),
		)
	}

	return err
}

// Shutdown gracefully shuts down all connections
func (m *Manager) Shutdown(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "tenant.manager.shutdown")
	defer span.End()

	logger.Info(ctx, "initiating tenant manager shutdown")

	m.cancel()

	// Get all connections
	m.mu.Lock()
	connections := make([]Connection, 0, len(m.connections))
	tenantIDs := make([]string, 0, len(m.connections))
	for tenantID, conn := range m.connections {
		connections = append(connections, conn)
		tenantIDs = append(tenantIDs, tenantID)
	}
	m.connections = make(map[string]Connection)
	m.mu.Unlock()

	span.SetAttributes(tracer.IntAttr("connections.count", len(connections)))

	// Close all connections concurrently
	var wg sync.WaitGroup
	for i, conn := range connections {
		wg.Add(1)
		go func(c Connection, tid string) {
			defer wg.Done()
			if err := c.Close(ctx); err != nil {
				logger.Warn(ctx, "error closing tenant connection during shutdown",
					logger.String("tenant_id", tid),
					logger.Err(err),
				)
			}
		}(conn, tenantIDs[i])
	}

	// Wait for connections to close
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		logger.Warn(ctx, "shutdown timeout reached, some connections may not be closed cleanly")
		span.SetAttributes(tracer.BoolAttr("shutdown.timeout", true))
	case <-done:
		// All closed
	}

	// Wait for cleanup loop
	m.wg.Wait()

	// Mark error channel as closed before closing to prevent race
	m.errorClosed.Store(true)
	close(m.errorChan)

	span.SetOK()
	logger.Info(ctx, "tenant manager shutdown complete",
		logger.Int("connections_closed", len(connections)),
	)

	return nil
}

// Errors returns the error channel
func (m *Manager) Errors() <-chan core.TenantError {
	return m.errorChan
}

// OnError sets an error callback
func (m *Manager) OnError(cb func(core.TenantError)) {
	m.errorMu.Lock()
	m.errorCb = cb
	m.errorMu.Unlock()
}

// dispatchError dispatches an error
func (m *Manager) dispatchError(err core.TenantError) {
	// Check if shutdown is in progress
	if m.errorClosed.Load() {
		return
	}

	m.errorMu.RLock()
	cb := m.errorCb
	m.errorMu.RUnlock()

	if cb != nil {
		cb(err)
	}

	// Double-check after callback (shutdown might have occurred)
	if m.errorClosed.Load() {
		return
	}

	select {
	case m.errorChan <- err:
	default:
		// Channel full, drop error
	}
}

// Health returns health status for all tenants
func (m *Manager) Health() map[string]core.HealthStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()

	health := make(map[string]core.HealthStatus, len(m.connections))
	for tenantID, conn := range m.connections {
		health[tenantID] = conn.Health()
	}
	return health
}

// TenantHealth returns health for a specific tenant
func (m *Manager) TenantHealth(tenantID string) (core.HealthStatus, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	conn, ok := m.connections[tenantID]
	if !ok {
		return core.HealthStatus{}, false
	}
	return conn.Health(), true
}

// ActiveConnections returns the count of active connections
func (m *Manager) ActiveConnections() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.connections)
}

// TenantIDs returns all connected tenant IDs
func (m *Manager) TenantIDs() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ids := make([]string, 0, len(m.connections))
	for id := range m.connections {
		ids = append(ids, id)
	}
	return ids
}

// SetCredentialManager sets the credential manager
func (m *Manager) SetCredentialManager(cm *auth.CredentialManager) {
	m.credManager = cm
}

// CredentialManager returns the credential manager
func (m *Manager) CredentialManager() *auth.CredentialManager {
	return m.credManager
}
