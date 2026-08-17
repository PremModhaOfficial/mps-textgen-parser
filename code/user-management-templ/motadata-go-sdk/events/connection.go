package events

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/config"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// ConnectionState represents the current lifecycle phase of a NATS connection.
// The state machine follows a well-defined transition order:
//
//	Disconnected -> Connecting -> Connected -> Draining -> Closed
//	                                  |
//	                                  +---> Reconnecting -> Connected (on success)
//	                                                    \-> Disconnected (on failure)
//
// All state transitions are protected by the Connection's sync.RWMutex to ensure
// thread-safe reads and writes from concurrent goroutines and NATS event handlers.
type ConnectionState int

const (
	// StateDisconnected indicates no active connection to any NATS server.
	// This is the initial state of a newly created Connection and the fallback
	// state when a connect attempt or reconnection fails.
	StateDisconnected ConnectionState = iota

	// StateConnecting indicates a connection attempt is in progress.
	// The connection transitions from Disconnected to Connecting when Connect
	// is called, and will move to Connected on success or back to Disconnected
	// on failure.
	StateConnecting

	// StateConnected indicates a live, usable connection to a NATS server.
	// Publishing, subscribing, and JetStream operations are only valid in
	// this state.
	StateConnected

	// StateReconnecting indicates the connection was lost and the NATS client
	// is automatically attempting to re-establish it using the configured
	// reconnect policy. Messages may be buffered up to defaultReconnectBufSize
	// during this phase.
	StateReconnecting

	// StateDraining indicates the connection is performing a graceful shutdown.
	// In-flight messages are flushed and subscriptions are unsubscribed before
	// the underlying NATS connection is closed. No new operations should be
	// started in this state.
	StateDraining

	// StateClosed indicates the connection has been permanently closed and
	// cannot be reused. All resources (conn, JetStream context) have been
	// released. A new Connection must be created to reconnect.
	StateClosed
)

// Internal connection tuning constants.
const (
	// defaultReconnectBufSize controls how many bytes of outgoing messages the
	// NATS client will buffer while in the StateReconnecting phase. Messages
	// exceeding this limit are dropped. 8 MB provides a reasonable window for
	// short network interruptions without excessive memory growth.
	defaultReconnectBufSize = 8 * 1024 * 1024

	// drainPollInterval is the polling interval used by closedCh to check
	// whether the NATS connection has finished draining. A shorter interval
	// reduces Close latency at the cost of slightly higher CPU usage during
	// the drain phase.
	drainPollInterval = 50 * time.Millisecond
)

// String returns a human-readable label for the ConnectionState, suitable for
// logging and diagnostic output. Unknown or future state values return "unknown".
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

// HealthStatus represents a point-in-time snapshot of a Connection's health.
// It aggregates the current lifecycle state, the most recent error (if any),
// and diagnostic details such as the connected server URL and reconnection
// count. HealthStatus values are safe to read concurrently once obtained,
// but they become stale as the connection state evolves.
type HealthStatus struct {
	// State is the connection's lifecycle phase at the time of the snapshot.
	State ConnectionState

	// Healthy is true only when the connection is in StateConnected and the
	// underlying NATS connection confirms it is still alive.
	Healthy bool

	// Message is a human-readable description of the current state, useful
	// for health-check endpoints and operational dashboards.
	Message string

	// LastError holds the most recent error observed on the connection, or nil
	// if no error has occurred since the last successful (re)connect.
	LastError error

	// Details contains additional diagnostic metadata such as tenant_id,
	// configured servers, connected_url, and reconnect count.
	Details map[string]any
}

// IsHealthy returns true only when the snapshot indicates a fully operational
// connection -- both the Healthy flag is set and the state is StateConnected.
// Use this as a quick predicate in readiness probes or load-balancer checks.
func (h HealthStatus) IsHealthy() bool {
	return h.Healthy && h.State == StateConnected
}

// Connection is a tenant-scoped wrapper around a raw [nats.Conn] that adds
// lifecycle management, JetStream integration, TLS configuration, automatic
// reconnection, and health monitoring.
//
// Thread safety: Connection is safe for concurrent use. A [sync.RWMutex] (mu)
// guards mutable fields such as state, conn, js, lastError, onError, and
// healthMonitor. The lastUsed timestamp is managed via [atomic.Int64] for
// lock-free reads and writes on the hot path.
//
// Typical usage:
//
//	conn := events.NewConnection(tenantID, cfg, creds)
//	if err := conn.Connect(ctx); err != nil { ... }
//	defer conn.Close(ctx)
//	// use conn.Conn() for core NATS, conn.JetStream() for streaming
type Connection struct {
	tenantID string
	conn     *nats.Conn
	js       jetstream.JetStream
	config   *config.EventsConfig
	creds    *auth.Credentials

	// mu protects state, conn, js, lastError, onError, and healthMonitor.
	// Use RLock for reads (IsConnected, State, Health) and Lock for writes
	// (Connect, Close, event handlers).
	mu        sync.RWMutex
	state     ConnectionState
	lastError error

	// lastUsed stores the Unix timestamp of the most recent activity on this
	// connection. It is updated atomically by Touch and read by IsIdle to
	// support idle-connection eviction in connection pools.
	lastUsed atomic.Int64

	// onError is an optional callback invoked on disconnect or async errors.
	// Protected by mu.
	onError func(error)

	// healthMonitor, when set, receives lifecycle events (connect, disconnect,
	// reconnect, error) for aggregation and reporting. Protected by mu.
	healthMonitor *HealthMonitor
}

// NewConnection creates a new Connection in the StateDisconnected phase.
// The returned Connection is not yet connected to NATS; call [Connection.Connect]
// to establish the network connection. The initial last-used timestamp is set to
// the current time so the connection is not immediately considered idle.
//
// Parameters:
//   - tenantID: unique identifier for the tenant that owns this connection,
//     used for subject prefixing, logging, and connection naming.
//   - cfg: events configuration containing server addresses, TLS settings,
//     reconnect policy, JetStream options, and timeouts.
//   - creds: optional authentication credentials (user/pass, token, NKey, JWT,
//     or credentials file). Pass nil for unauthenticated connections.
func NewConnection(tenantID string, cfg *config.EventsConfig, creds *auth.Credentials) *Connection {
	c := &Connection{
		tenantID: tenantID,
		config:   cfg,
		creds:    creds,
		state:    StateDisconnected,
	}
	c.Touch()
	return c
}

// SetHealthMonitor attaches a [HealthMonitor] to the connection so that all
// lifecycle events (connect, disconnect, reconnect, error) are automatically
// recorded. Pass nil to detach a previously attached monitor. This method is
// safe to call concurrently; the monitor reference is protected by the
// connection's mutex.
func (c *Connection) SetHealthMonitor(hm *HealthMonitor) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.healthMonitor = hm
}

// Connect establishes the underlying NATS connection and, when enabled,
// initializes the JetStream context. The method is idempotent-safe: calling
// Connect on an already-connected Connection returns [utils.ErrAlreadyConnected],
// and calling it on a closed or draining Connection returns
// [utils.ErrConnectionClosed].
//
// State transitions performed by Connect:
//
//	Disconnected -> Connecting -> Connected  (success)
//	Disconnected -> Connecting -> Disconnected  (failure)
//
// Connect is protected against concurrent Close calls -- if the state is
// changed to Closed or Draining while the NATS dial is in progress, the
// newly opened connection is immediately closed and an error is returned.
//
// JetStream is initialized according to the Stream configuration: a domain-scoped
// context is created when Stream.Domain is set, an API-prefix-scoped context when
// Stream.Prefix is set, and a default context otherwise. JetStream initialization
// failure is non-fatal; a warning is logged and the connection remains usable for
// core NATS operations.
//
// The provided context is used for OpenTelemetry span propagation and structured
// logging but does not control the connection timeout itself -- that is governed
// by config.ConnectTimeout.
func (c *Connection) Connect(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "nats.connect",
		tracer.WithAttributes(
			tracer.StringAttr("tenant.id", c.tenantID),
		),
	)
	defer span.End()

	c.mu.Lock()
	if c.state == StateConnected {
		c.mu.Unlock()
		span.SetError(utils.ErrAlreadyConnected)
		return utils.ErrAlreadyConnected
	}
	if c.state == StateClosed || c.state == StateDraining {
		c.mu.Unlock()
		span.SetError(utils.ErrConnectionClosed)
		return utils.ErrConnectionClosed
	}
	c.state = StateConnecting
	c.mu.Unlock()

	opts, buildErr := c.buildOptions()
	if buildErr != nil {
		c.mu.Lock()
		c.state = StateDisconnected
		c.mu.Unlock()
		span.SetError(buildErr)
		return buildErr
	}

	// Add connection timeout as a NATS option
	opts = append(opts, nats.Timeout(c.config.ConnectTimeout))

	// Build server URL
	serverURL := strings.Join(c.config.Servers, ",")
	span.SetAttributes(tracer.StringAttr("servers", serverURL))

	logger.Debug(ctx, "connecting to NATS",
		logger.String("tenant_id", c.tenantID),
		logger.String("servers", serverURL),
	)

	conn, err := nats.Connect(serverURL, opts...)
	if err != nil {
		c.mu.Lock()
		c.state = StateDisconnected
		c.mu.Unlock()
		span.SetError(err)
		logger.Error(ctx, "failed to connect to NATS",
			logger.String("tenant_id", c.tenantID),
			logger.Err(err),
		)
		return err
	}

	// nats.Connect returns only when connected (or error), so no polling needed.
	c.mu.Lock()
	// Guard against concurrent Close() call during connect
	if c.state != StateConnecting {
		c.mu.Unlock()
		conn.Close()
		span.SetError(utils.ErrConnectionClosed)
		return utils.ErrConnectionClosed
	}
	c.conn = conn
	c.state = StateConnected
	c.mu.Unlock()

	// Initialize JetStream if enabled
	if c.config.Stream.Enabled {
		c.initJetStream(ctx, conn, span)
	}

	c.Touch()

	// Record connect event on health monitor
	c.mu.RLock()
	hm := c.healthMonitor
	c.mu.RUnlock()
	if hm != nil {
		hm.RecordConnect()
	}

	span.SetOK()
	logger.Info(ctx, "NATS connection established",
		logger.String("tenant_id", c.tenantID),
		logger.String("connected_url", conn.ConnectedUrl()),
	)
	return nil
}

// Close performs a graceful shutdown of the NATS connection. It first attempts
// to drain the connection -- flushing in-flight messages and unsubscribing all
// subscriptions -- within the configured DrainTimeout. If the drain does not
// complete in time or fails outright, the connection is forcibly closed.
//
// State transitions:
//
//	Connected / Reconnecting / Connecting -> Draining -> Closed
//
// Close is idempotent: calling it on an already-closed or draining connection
// returns nil immediately. After Close returns, the Connection's conn and js
// fields are set to nil and the state is StateClosed. The Connection cannot be
// reused; create a new one via [NewConnection] to reconnect.
//
// The provided context is used for OpenTelemetry tracing and logging.
func (c *Connection) Close(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "nats.close",
		tracer.WithAttributes(
			tracer.StringAttr("tenant.id", c.tenantID),
		),
	)
	defer span.End()

	c.mu.Lock()
	if c.state == StateClosed || c.state == StateDraining {
		c.mu.Unlock()
		span.SetOK()
		return nil
	}
	c.state = StateDraining
	conn := c.conn
	c.mu.Unlock()

	logger.Debug(ctx, "closing NATS connection",
		logger.String("tenant_id", c.tenantID),
	)

	if conn != nil {
		drainCtx, cancel := context.WithTimeout(ctx, c.config.DrainTimeout)
		defer cancel()

		if err := conn.Drain(); err != nil {
			logger.Warn(ctx, "NATS drain failed, forcing close",
				logger.String("tenant_id", c.tenantID),
				logger.Err(err),
			)
			conn.Close()
		} else {
			// Wait for drain to complete or timeout
			select {
			case <-drainCtx.Done():
				logger.Warn(ctx, "NATS drain timeout, forcing close",
					logger.String("tenant_id", c.tenantID),
				)
				conn.Close()
			case <-closedCh(conn):
				// Drain completed successfully
			}
		}
	}

	c.mu.Lock()
	c.state = StateClosed
	c.conn = nil
	c.js = nil
	c.mu.Unlock()

	span.SetOK()
	logger.Info(ctx, "NATS connection closed",
		logger.String("tenant_id", c.tenantID),
	)
	return nil
}

// closedCh returns a channel that closes when the NATS connection is fully closed.
func closedCh(conn *nats.Conn) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		for conn.IsDraining() {
			time.Sleep(drainPollInterval)
		}
		close(ch)
	}()
	return ch
}

// IsConnected returns true if connected.
func (c *Connection) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.state == StateConnected && c.conn != nil && c.conn.IsConnected()
}

// State returns the connection state.
func (c *Connection) State() ConnectionState {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.state
}

// TenantID returns the tenant identifier.
func (c *Connection) TenantID() string {
	return c.tenantID
}

// Touch updates the last activity timestamp.
func (c *Connection) Touch() {
	c.lastUsed.Store(time.Now().Unix())
}

// IsIdle returns true if the connection has been idle longer than timeout.
func (c *Connection) IsIdle(timeout time.Duration) bool {
	lastUsed := time.Unix(c.lastUsed.Load(), 0)
	return time.Since(lastUsed) > timeout
}

// Health returns the health status.
func (c *Connection) Health() HealthStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()

	healthy := c.state == StateConnected && c.conn != nil && c.conn.IsConnected()

	// Build status message based on state
	var message string
	switch c.state {
	case StateConnected:
		message = "connected to NATS server"
	case StateConnecting:
		message = "connecting to NATS server"
	case StateReconnecting:
		message = "reconnecting to NATS server"
	case StateDraining:
		message = "draining connection"
	case StateClosed:
		message = "connection closed"
	case StateDisconnected:
		message = "disconnected from NATS server"
	default:
		message = "unknown state"
	}

	// Build details
	details := map[string]any{
		"tenant_id": c.tenantID,
		"servers":   c.config.Servers,
	}

	if c.conn != nil {
		details["connected_url"] = c.conn.ConnectedUrl()
		details["reconnects"] = c.conn.Stats().Reconnects
	}

	return HealthStatus{
		State:     c.state,
		Healthy:   healthy,
		Message:   message,
		LastError: c.lastError,
		Details:   details,
	}
}

// initJetStream creates the JetStream context using the configured domain,
// prefix, or default settings. Failure is non-fatal and logged as a warning.
func (c *Connection) initJetStream(ctx context.Context, conn *nats.Conn, span tracer.Span) {
	js, jsErr := c.createJetStream(conn)
	if jsErr != nil {
		logger.Warn(ctx, "JetStream not available",
			logger.String("tenant_id", c.tenantID),
			logger.Err(jsErr),
		)
		return
	}

	c.mu.Lock()
	c.js = js
	c.mu.Unlock()
	span.SetAttributes(tracer.BoolAttr("jetstream.enabled", true))
}

// createJetStream creates the JetStream context based on configuration.
func (c *Connection) createJetStream(conn *nats.Conn) (jetstream.JetStream, error) {
	if c.config.Stream.Domain != "" {
		return jetstream.NewWithDomain(conn, c.config.Stream.Domain)
	}
	if c.config.Stream.Prefix != "" {
		return jetstream.NewWithAPIPrefix(conn, c.config.Stream.Prefix)
	}
	return jetstream.New(conn)
}

// Conn returns the raw NATS connection.
func (c *Connection) Conn() *nats.Conn {
	c.Touch()
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn
}

// JetStream returns the JetStream context.
func (c *Connection) JetStream() jetstream.JetStream {
	c.Touch()
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.js
}

// OnError sets the error callback.
func (c *Connection) OnError(cb func(error)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.onError = cb
}

// buildOptions builds NATS connection options. Returns an error if TLS is
// enabled but certificate loading fails (TLS failure is a hard error).
func (c *Connection) buildOptions() ([]nats.Option, error) {
	opts := []nats.Option{
		nats.Name(c.config.Name + "-" + c.tenantID),
		nats.ReconnectWait(c.config.Reconnect.InitialInterval),
		nats.MaxReconnects(c.config.Reconnect.MaxAttempts),
		nats.ReconnectBufSize(defaultReconnectBufSize),
	}

	// Authentication
	if c.creds != nil {
		opts = c.appendAuthOptions(opts)
	}

	// TLS — failure is a hard error when TLS is explicitly enabled
	if c.config.TLS != nil && c.config.TLS.Enabled {
		tlsOpt, tlsErr := c.buildTLSOption()
		if tlsErr != nil {
			return nil, tlsErr
		}
		if tlsOpt != nil {
			opts = append(opts, tlsOpt)
		}
	}

	// Event handlers
	opts = append(opts,
		nats.DisconnectErrHandler(c.handleDisconnect),
		nats.ReconnectHandler(c.handleReconnect),
		nats.ClosedHandler(c.handleClosed),
		nats.ErrorHandler(c.handleError),
	)

	return opts, nil
}

// appendAuthOptions appends the appropriate NATS auth option based on credential type.
func (c *Connection) appendAuthOptions(opts []nats.Option) []nats.Option {
	switch c.creds.Type {
	case auth.TypeUserPass:
		return append(opts, nats.UserInfo(c.creds.Username, c.creds.Password))
	case auth.TypeToken:
		return append(opts, nats.Token(c.creds.Token))
	case auth.TypeCredentialsFile:
		return append(opts, nats.UserCredentials(c.creds.CredentialsFile))
	case auth.TypeNKey:
		if c.creds.NKeyFile != "" {
			if opt, err := nats.NkeyOptionFromSeed(c.creds.NKeyFile); err == nil {
				return append(opts, opt)
			}
		}
	case auth.TypeJWT:
		if c.creds.JWT != "" {
			return append(opts, nats.UserJWTAndSeed(c.creds.JWT, c.creds.Seed))
		}
	}
	return opts
}

// buildTLSOption builds the NATS TLS option, logging a warning if SkipVerify is set.
func (c *Connection) buildTLSOption() (nats.Option, error) {
	if c.config.TLS.SkipVerify {
		logger.Warn(context.Background(), "TLS certificate verification is disabled (InsecureSkipVerify=true)",
			logger.String("tenant_id", c.tenantID),
		)
	}

	tlsConfig, tlsErr := c.buildTLSConfig()
	if tlsErr != nil {
		return nil, utils.WrapError(tlsErr, "TLS configuration failed")
	}
	if tlsConfig == nil {
		return nil, nil
	}
	return nats.Secure(tlsConfig), nil
}

func (c *Connection) buildTLSConfig() (*tls.Config, error) {
	if c.config.TLS == nil || !c.config.TLS.Enabled {
		return nil, nil
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.config.TLS.SkipVerify,
	}

	if c.config.TLS.CertFile != "" && c.config.TLS.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(c.config.TLS.CertFile, c.config.TLS.KeyFile)
		if err != nil {
			return nil, utils.NewError("tls", "load_cert", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	if c.config.TLS.CAFile != "" {
		caCert, err := os.ReadFile(c.config.TLS.CAFile)
		if err != nil {
			return nil, utils.NewError("tls", "load_ca", err)
		}
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, utils.NewError("tls", "parse_ca", utils.ErrInvalidConfig)
		}
		tlsConfig.RootCAs = caCertPool
	}

	return tlsConfig, nil
}

func (c *Connection) handleDisconnect(_ *nats.Conn, err error) {
	ctx := context.Background()

	c.mu.Lock()
	c.state = StateReconnecting
	if err != nil {
		c.lastError = err
	}
	hm := c.healthMonitor
	cb := c.onError
	c.mu.Unlock()

	logger.Warn(ctx, "NATS connection disconnected",
		logger.String("tenant_id", c.tenantID),
		logger.Err(err),
	)

	if hm != nil {
		hm.RecordDisconnect()
	}

	if err != nil && cb != nil {
		cb(err)
	}
}

func (c *Connection) handleReconnect(conn *nats.Conn) {
	ctx := context.Background()

	c.mu.Lock()
	c.state = StateConnected
	c.lastError = nil // Clear error on successful reconnect
	hm := c.healthMonitor
	c.mu.Unlock()

	c.Touch()

	if hm != nil {
		hm.RecordReconnect()
	}

	logger.Info(ctx, "NATS connection reconnected",
		logger.String("tenant_id", c.tenantID),
		logger.String("connected_url", conn.ConnectedUrl()),
	)
}

func (c *Connection) handleClosed(_ *nats.Conn) {
	ctx := context.Background()

	c.mu.Lock()
	c.state = StateClosed
	c.mu.Unlock()

	logger.Info(ctx, "NATS connection closed by server",
		logger.String("tenant_id", c.tenantID),
	)
}

func (c *Connection) handleError(_ *nats.Conn, sub *nats.Subscription, err error) {
	c.mu.Lock()
	if err != nil {
		c.lastError = err
	}
	hm := c.healthMonitor
	cb := c.onError
	c.mu.Unlock()

	if err != nil {
		c.logSubscriptionError(sub, err)
		if hm != nil {
			hm.RecordError("subscription")
		}
	}

	if cb != nil {
		cb(err)
	}
}

func (c *Connection) logSubscriptionError(sub *nats.Subscription, err error) {
	subject := ""
	if sub != nil {
		subject = sub.Subject
	}
	logger.Error(context.Background(), "NATS error occurred",
		logger.String("tenant_id", c.tenantID),
		logger.String("subject", subject),
		logger.Err(err),
	)
}
