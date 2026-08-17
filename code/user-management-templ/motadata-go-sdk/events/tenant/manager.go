package tenant

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/config"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer"
)

// Default tenant manager constants used during connection lifecycle management.
const (
	// defaultCleanupTimeout limits how long the manager waits for a single idle
	// connection to drain during background cleanup. If a connection does not
	// finish draining within this window, it is force-closed to avoid blocking
	// the cleanup loop.
	defaultCleanupTimeout = 10 * time.Second

	// defaultErrorBufferSize sets the capacity of the buffered error channel.
	// When the channel is full, new errors are silently dropped to prevent
	// blocking the producer goroutine.
	defaultErrorBufferSize = 100
)

// ConnectionBuilder is a factory function that creates NATS connections for
// individual tenants. It abstracts the transport-level connection details so
// the Manager remains agnostic to NATS-specific configuration (TLS, cluster
// URLs, authentication options, etc.). Implementations should apply the
// provided credentials and configuration to establish a fully authenticated
// connection and return both the raw NATS connection and the JetStream context
// for stream-based operations.
type ConnectionBuilder func(ctx context.Context, tenantID string, cfg *config.EventsConfig, creds *auth.Credentials) (*nats.Conn, jetstream.JetStream, error)

// TenantConnection holds a NATS connection that is scoped to a single tenant.
// It tracks the last time the connection was actively used (via an atomic
// timestamp) so the Manager can identify and reclaim idle connections during
// periodic cleanup sweeps. All exported methods are safe for concurrent use.
type TenantConnection struct {
	// tenantID is the immutable identifier set at creation time.
	tenantID string
	// conn is the underlying NATS connection for this tenant.
	conn *nats.Conn
	// js is the JetStream context derived from conn, used for stream operations.
	js jetstream.JetStream
	// lastUsed stores a UnixNano timestamp updated on every Touch call.
	// It is accessed atomically to avoid locking on the hot read path.
	lastUsed atomic.Int64
}

// TenantID returns the immutable tenant identifier that was assigned when the
// connection was created.
func (tc *TenantConnection) TenantID() string {
	return tc.tenantID
}

// Conn returns the underlying NATS connection. Callers should check
// IsConnected before performing operations, as the connection may have been
// dropped by the server or drained during shutdown.
func (tc *TenantConnection) Conn() *nats.Conn {
	return tc.conn
}

// JetStream returns the JetStream context associated with this tenant's NATS
// connection, enabling stream, consumer, key-value, and object store operations.
func (tc *TenantConnection) JetStream() jetstream.JetStream {
	return tc.js
}

// Touch updates the last-activity timestamp to the current time. This is
// called automatically by the Manager when a connection is retrieved or
// reused, which prevents the idle-cleanup loop from reclaiming connections
// that are still in active use.
func (tc *TenantConnection) Touch() {
	tc.lastUsed.Store(time.Now().UnixNano())
}

// IsIdle returns true if the connection has been idle longer than the given
// timeout. A connection with a zero lastUsed timestamp (never touched) is
// treated as non-idle to avoid prematurely closing newly created connections
// that have not yet been used.
func (tc *TenantConnection) IsIdle(timeout time.Duration) bool {
	last := tc.lastUsed.Load()
	if last == 0 {
		return false
	}
	return time.Since(time.Unix(0, last)) > timeout
}

// IsConnected returns true if the underlying NATS connection is non-nil and
// reports itself as connected. This is a point-in-time check; the connection
// status may change immediately after the call returns.
func (tc *TenantConnection) IsConnected() bool {
	return tc.conn != nil && tc.conn.IsConnected()
}

// Close gracefully closes the tenant connection by first attempting to drain
// all pending messages. Draining ensures that in-flight publishes are flushed
// and active subscriptions process remaining messages before the connection is
// terminated. If the provided context expires before draining completes, the
// connection is force-closed and the context error is returned.
func (tc *TenantConnection) Close(ctx context.Context) error {
	if tc.conn == nil {
		return nil
	}

	// Drain with context timeout awareness
	done := make(chan struct{})
	go func() {
		tc.conn.Drain()
		close(done)
	}()

	select {
	case <-ctx.Done():
		tc.conn.Close()
		return ctx.Err()
	case <-done:
		return nil
	}
}

// HealthInfo represents a snapshot of the health status for a single tenant
// connection. It is returned by Manager.Health and Manager.TenantHealth to
// allow monitoring dashboards and health-check endpoints to report per-tenant
// connectivity state.
type HealthInfo struct {
	// Connected indicates whether the NATS connection is active at the time of the check.
	Connected bool
	// TenantID identifies the tenant this health record belongs to.
	TenantID string
	// URL is the NATS server URL that the tenant is currently connected to.
	// It may be empty if the connection has been closed.
	URL string
}

// Manager manages the full lifecycle of NATS connections for multiple tenants.
//
// Connection lifecycle: A connection is created on the first call to Connect
// for a given tenant ID. Subsequent Connect or Get calls for the same tenant
// reuse the existing connection and refresh its activity timestamp. Connections
// are removed either explicitly via Disconnect or automatically by the
// background idle-cleanup loop when they exceed the configured IdleTimeout.
//
// Idle cleanup: A background goroutine runs on a configurable interval
// (CleanupInterval) and closes any connection whose last-activity timestamp
// is older than IdleTimeout. This prevents resource leaks from tenants that
// stop communicating without explicitly disconnecting.
//
// Reconnection: The Manager itself does not perform automatic reconnection;
// NATS client-level reconnection is handled by the nats.Conn options set in
// the ConnectionBuilder. If a connection is detected as disconnected, the
// caller should invoke Connect again to establish a fresh connection.
//
// All exported methods are safe for concurrent use. The internal connection
// map is protected by a sync.RWMutex to allow parallel reads with exclusive
// writes.
type Manager struct {
	// config holds the shared event system configuration (timeouts, intervals, etc.).
	config *config.EventsConfig
	// builder is the factory function used to create new NATS connections.
	builder ConnectionBuilder

	// mu protects the connections map. Read-lock is acquired for lookups;
	// write-lock is acquired for insertions, deletions, and cleanup sweeps.
	mu          sync.RWMutex
	connections map[string]*TenantConnection

	// credManager provides credential resolution when explicit credentials are
	// not supplied to Connect.
	credManager *auth.CredentialManager

	// done is closed when cancel is called, signalling the cleanup goroutine
	// to stop and preventing new connections from being established.
	done <-chan struct{}
	// cancel stops the Manager's background goroutines.
	cancel context.CancelFunc
	// wg tracks background goroutines (currently the cleanup loop) so Shutdown
	// can wait for them to exit before returning.
	wg sync.WaitGroup

	// errorChan is a buffered channel that receives tenant-scoped errors for
	// consumer-driven processing.
	errorChan chan utils.TenantError
	// errorCb is an optional synchronous callback invoked before the error is
	// sent to errorChan. Protected by errorMu.
	errorCb func(utils.TenantError)
	// errorMu protects concurrent reads and writes to errorCb.
	errorMu sync.RWMutex
	// errorClosed is set to true during Shutdown, just before errorChan is
	// closed, to prevent a panic from sending on a closed channel.
	errorClosed atomic.Bool
}

// ManagerConfig holds the configuration required to initialize a Manager.
// At minimum, a ConnectionBuilder should be supplied; the remaining fields
// have sensible defaults.
type ManagerConfig struct {
	// Config provides shared event system settings such as IdleTimeout and
	// CleanupInterval. If nil, DefaultEventsConfig is used.
	Config *config.EventsConfig
	// Builder is the factory function that creates NATS connections per tenant.
	Builder ConnectionBuilder
	// CredManager is an optional credential store used to resolve credentials
	// when they are not explicitly passed to Connect.
	CredManager *auth.CredentialManager
	// ErrorBufferSize sets the capacity of the internal error channel.
	// Defaults to defaultErrorBufferSize (100) when zero.
	ErrorBufferSize int
}

// NewManager creates a new tenant Manager and starts its background
// idle-cleanup goroutine. The returned Manager is ready to accept connections
// via Connect. The caller must call Shutdown when the Manager is no longer
// needed to release all connections and stop background work.
// Returns an error if the provided configuration fails validation.
func NewManager(cfg ManagerConfig) (*Manager, error) {
	if cfg.Config == nil {
		cfg.Config = config.DefaultEventsConfig()
	}

	if err := config.ValidateEventsConfig(cfg.Config); err != nil {
		return nil, err
	}

	if cfg.ErrorBufferSize == 0 {
		cfg.ErrorBufferSize = defaultErrorBufferSize
	}

	ctx, cancel := context.WithCancel(context.Background())

	m := &Manager{
		config:      cfg.Config,
		builder:     cfg.Builder,
		connections: make(map[string]*TenantConnection),
		credManager: cfg.CredManager,
		done:        ctx.Done(),
		cancel:      cancel,
		errorChan:   make(chan utils.TenantError, cfg.ErrorBufferSize),
	}

	// Start cleanup loop
	m.startCleanupLoop()

	return m, nil
}

// startCleanupLoop launches a background goroutine that periodically scans
// the connection map and closes any connections that have been idle longer
// than the configured IdleTimeout. The goroutine exits when the Manager's
// context is cancelled (i.e., during Shutdown).
func (m *Manager) startCleanupLoop() {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()

		ticker := time.NewTicker(m.config.CleanupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-m.done:
				return
			case <-ticker.C:
				m.cleanupIdleConnections()
			}
		}
	}()
}

// cleanupIdleConnections iterates over all active connections under a write
// lock and closes any whose last-activity timestamp exceeds IdleTimeout.
// Each connection is drained with a bounded timeout (defaultCleanupTimeout)
// to avoid blocking the cleanup loop on a slow drain. Closed connections are
// immediately removed from the map.
func (m *Manager) cleanupIdleConnections() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for tenantID, tc := range m.connections {
		if tc.IsIdle(m.config.IdleTimeout) {
			ctx, cancel := context.WithTimeout(context.Background(), defaultCleanupTimeout)
			if err := tc.Close(ctx); err != nil {
				logger.Warn(ctx, "error closing idle tenant connection",
					logger.String("tenant_id", tenantID),
					logger.Err(err),
				)
			}
			cancel()
			delete(m.connections, tenantID)

			logger.Info(ctx, "closed idle tenant connection",
				logger.String("tenant_id", tenantID),
				logger.Duration("idle_timeout", m.config.IdleTimeout),
			)
		}
	}
}

// Connect establishes or retrieves a NATS connection for the specified tenant.
//
// The method follows a connect-or-reuse strategy:
//  1. If the Manager is shutting down, it returns ErrShutdownInProgress.
//  2. If a connection for the tenant already exists, it is reused (Touch is
//     called to reset the idle timer).
//  3. Otherwise, credentials are resolved (from the argument or the
//     CredentialManager) and a new connection is built via the ConnectionBuilder.
//  4. The new connection is stored in the map. If another goroutine raced and
//     stored a connection for the same tenant in the meantime, the duplicate
//     is closed and the winner is returned.
//
// The creds parameter may be nil if a CredentialManager has been configured.
func (m *Manager) Connect(ctx context.Context, tenantID string, creds *auth.Credentials) (*TenantConnection, error) {
	ctx, span := tracer.Start(ctx, "tenant.connect",
		tracer.WithAttributes(
			tracer.StringAttr("tenant.id", tenantID),
		),
	)
	defer span.End()

	// Check for shutdown
	select {
	case <-m.done:
		span.SetError(utils.ErrShutdownInProgress)
		return nil, utils.ErrShutdownInProgress
	default:
	}

	// Reuse existing connection if present
	if existing := m.tryReuseConnection(tenantID); existing != nil {
		span.SetAttributes(tracer.BoolAttr("connection.reused", true))
		span.SetOK()
		return existing, nil
	}

	// Resolve credentials
	creds, err := m.resolveCredentials(ctx, tenantID, creds, span)
	if err != nil {
		return nil, err
	}

	// Build the connection
	tc, err := m.buildConnection(ctx, tenantID, creds, span)
	if err != nil {
		return nil, err
	}

	// Store, handling the race where another goroutine connected first
	return m.storeConnection(ctx, tenantID, tc, span)
}

// tryReuseConnection performs a read-locked lookup for an existing connection.
// If found, it updates the activity timestamp and returns the connection.
// Returns nil when no connection exists for the given tenant, signalling that
// a new connection must be built.
func (m *Manager) tryReuseConnection(tenantID string) *TenantConnection {
	m.mu.RLock()
	existing, ok := m.connections[tenantID]
	m.mu.RUnlock()

	if ok {
		existing.Touch()
		return existing
	}
	return nil
}

// resolveCredentials determines the credentials to use for a tenant connection.
// If explicit credentials are provided they are used directly; otherwise the
// CredentialManager is consulted. Returns ErrInvalidCredential if no valid
// credentials can be resolved from either source.
func (m *Manager) resolveCredentials(ctx context.Context, tenantID string, creds *auth.Credentials, span tracer.Span) (*auth.Credentials, error) {
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
		span.SetError(utils.ErrInvalidCredential)
		return nil, utils.ErrInvalidCredential
	}

	return creds, nil
}

// buildConnection delegates to the ConnectionBuilder to create a new NATS
// connection and wraps the result in a TenantConnection. The activity
// timestamp is set immediately so the connection is not considered idle
// before the caller has a chance to use it. Returns an error if no builder
// is configured or the builder itself fails.
func (m *Manager) buildConnection(ctx context.Context, tenantID string, creds *auth.Credentials, span tracer.Span) (*TenantConnection, error) {
	if m.builder == nil {
		err := utils.NewError("connection", "create", utils.ErrNotConnected)
		span.SetError(err)
		return nil, err
	}

	nc, js, err := m.builder(ctx, tenantID, m.config, creds)
	if err != nil {
		span.SetError(err)
		logger.Error(ctx, "failed to create tenant connection",
			logger.String("tenant_id", tenantID),
			logger.Err(err),
		)
		return nil, err
	}

	tc := &TenantConnection{
		tenantID: tenantID,
		conn:     nc,
		js:       js,
	}
	tc.Touch()
	return tc, nil
}

// storeConnection inserts the newly built connection into the map under a
// write lock. Because Connect does not hold the write lock while building the
// connection (to avoid blocking other tenants), two goroutines may build
// connections for the same tenant concurrently. If a connection already exists
// when the lock is acquired, the duplicate is closed and the existing
// connection is returned to maintain a single-connection-per-tenant invariant.
func (m *Manager) storeConnection(ctx context.Context, tenantID string, tc *TenantConnection, span tracer.Span) (*TenantConnection, error) {
	m.mu.Lock()
	if existing, ok := m.connections[tenantID]; ok {
		m.mu.Unlock()
		if closeErr := tc.Close(ctx); closeErr != nil {
			logger.Warn(ctx, "error closing duplicate tenant connection",
				logger.String("tenant_id", tenantID),
				logger.Err(closeErr),
			)
		}
		existing.Touch()
		span.SetAttributes(tracer.BoolAttr("connection.reused", true))
		span.SetOK()
		return existing, nil
	}
	m.connections[tenantID] = tc
	m.mu.Unlock()

	span.SetAttributes(tracer.BoolAttr("connection.new", true))
	span.SetOK()
	logger.Info(ctx, "tenant connection established",
		logger.String("tenant_id", tenantID),
	)

	return tc, nil
}

// ConnectWithCredentialManager establishes a connection for the tenant by
// looking up credentials exclusively from the configured CredentialManager.
// Returns an error if no CredentialManager has been set or if the lookup fails.
// This is a convenience wrapper around Connect for callers that always rely
// on centrally managed credentials.
func (m *Manager) ConnectWithCredentialManager(ctx context.Context, tenantID string) (*TenantConnection, error) {
	if m.credManager == nil {
		return nil, utils.NewError("connection", "connect", auth.ErrInvalidCredential)
	}

	creds, err := m.credManager.GetCredentials(tenantID)
	if err != nil {
		return nil, err
	}

	return m.Connect(ctx, tenantID, creds)
}

// ConnectWithPayload validates and uses a RegistrationPayload to establish a
// tenant connection. If a CredentialManager is configured, the credentials
// from the payload are persisted there so that future reconnections can
// resolve them automatically. This is the preferred entry point when onboarding
// a new tenant for the first time.
func (m *Manager) ConnectWithPayload(ctx context.Context, payload auth.RegistrationPayload) (*TenantConnection, error) {
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

// Get returns an existing connection for the given tenant, or (nil, false) if
// the tenant is not connected. Retrieving a connection also refreshes its
// activity timestamp, preventing it from being reclaimed by the idle-cleanup
// loop while it is actively being used.
func (m *Manager) Get(tenantID string) (*TenantConnection, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tc, ok := m.connections[tenantID]
	if ok {
		tc.Touch()
	}
	return tc, ok
}

// Disconnect removes and gracefully closes the connection for the specified
// tenant. The connection is drained before closure, respecting the provided
// context's deadline. Returns ErrTenantNotFound if no connection exists for
// the tenant. After Disconnect returns successfully, a subsequent Connect call
// for the same tenant will create a fresh connection.
func (m *Manager) Disconnect(ctx context.Context, tenantID string) error {
	ctx, span := tracer.Start(ctx, "tenant.disconnect",
		tracer.WithAttributes(
			tracer.StringAttr("tenant.id", tenantID),
		),
	)
	defer span.End()

	m.mu.Lock()
	tc, ok := m.connections[tenantID]
	if !ok {
		m.mu.Unlock()
		span.SetError(utils.ErrTenantNotFound)
		return utils.ErrTenantNotFound
	}
	delete(m.connections, tenantID)
	m.mu.Unlock()

	err := tc.Close(ctx)
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

// Shutdown gracefully shuts down the Manager by performing the following steps
// in order:
//  1. Cancels the Manager's internal context, which stops the cleanup loop and
//     prevents new connections from being accepted.
//  2. Drains all active connections from the map (under a write lock).
//  3. Closes every connection concurrently, respecting the caller's context
//     deadline. If the deadline expires, remaining connections may not be
//     drained cleanly.
//  4. Waits for the cleanup goroutine to exit.
//  5. Marks the error channel as closed and closes it.
//
// Shutdown is idempotent in the sense that calling it multiple times will not
// panic, but only the first call performs meaningful work.
func (m *Manager) Shutdown(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "tenant.manager.shutdown")
	defer span.End()

	logger.Info(ctx, "initiating tenant manager shutdown")

	m.cancel()

	connections, tenantIDs := m.drainConnections()
	span.SetAttributes(tracer.IntAttr("connections.count", len(connections)))

	m.closeAllConnections(ctx, connections, tenantIDs, span)

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

// drainConnections atomically removes all connections from the map and returns
// them along with their tenant IDs. After this call the map is empty, so no
// new operations can find existing connections. The returned slices are
// ordered consistently (connection at index i belongs to tenantIDs[i]).
func (m *Manager) drainConnections() ([]*TenantConnection, []string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	connections := make([]*TenantConnection, 0, len(m.connections))
	tenantIDs := make([]string, 0, len(m.connections))
	for tenantID, tc := range m.connections {
		connections = append(connections, tc)
		tenantIDs = append(tenantIDs, tenantID)
	}
	m.connections = make(map[string]*TenantConnection)
	return connections, tenantIDs
}

// closeAllConnections closes every provided connection in its own goroutine
// and waits for all of them to finish. If the context deadline expires before
// all connections have drained, a warning is logged and the span is annotated
// with a timeout attribute. This fan-out approach minimizes total shutdown
// time when many tenants are connected simultaneously.
func (m *Manager) closeAllConnections(ctx context.Context, connections []*TenantConnection, tenantIDs []string, span tracer.Span) {
	var wg sync.WaitGroup
	for i, tc := range connections {
		wg.Add(1)
		go func(c *TenantConnection, tid string) {
			defer wg.Done()
			if err := c.Close(ctx); err != nil {
				logger.Warn(ctx, "error closing tenant connection during shutdown",
					logger.String("tenant_id", tid),
					logger.Err(err),
				)
			}
		}(tc, tenantIDs[i])
	}

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
}

// Errors returns a receive-only channel that delivers tenant-scoped errors.
// Consumers should read from this channel to observe connection failures,
// credential errors, and other asynchronous problems. The channel is closed
// during Shutdown; consumers should handle channel closure gracefully.
func (m *Manager) Errors() <-chan utils.TenantError {
	return m.errorChan
}

// OnError registers a synchronous callback that is invoked for every
// tenant error before the error is sent to the Errors channel. This allows
// immediate processing (e.g., metrics emission) without requiring a separate
// goroutine to drain the channel. The callback is protected by a read-write
// mutex and may be replaced at any time.
func (m *Manager) OnError(cb func(utils.TenantError)) {
	m.errorMu.Lock()
	m.errorCb = cb
	m.errorMu.Unlock()
}

// dispatchError delivers a tenant error through both the synchronous callback
// (if set) and the buffered error channel. It guards against sending on a
// closed channel by checking the errorClosed flag both before and after the
// callback invocation (shutdown may occur while the callback is running).
// If the channel is full, the error is silently dropped to avoid blocking.
func (m *Manager) dispatchError(err utils.TenantError) {
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
func (m *Manager) Health() map[string]HealthInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	health := make(map[string]HealthInfo, len(m.connections))
	for tenantID, tc := range m.connections {
		info := HealthInfo{
			Connected: tc.IsConnected(),
			TenantID:  tenantID,
		}
		if tc.conn != nil {
			info.URL = tc.conn.ConnectedUrl()
		}
		health[tenantID] = info
	}
	return health
}

// TenantHealth returns health for a specific tenant
func (m *Manager) TenantHealth(tenantID string) (HealthInfo, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tc, ok := m.connections[tenantID]
	if !ok {
		return HealthInfo{}, false
	}

	info := HealthInfo{
		Connected: tc.IsConnected(),
		TenantID:  tenantID,
	}
	if tc.conn != nil {
		info.URL = tc.conn.ConnectedUrl()
	}
	return info, true
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
