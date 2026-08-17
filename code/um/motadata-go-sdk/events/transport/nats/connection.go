package nats

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/config"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Connection wraps a NATS connection
type Connection struct {
	tenantID string
	conn     *nats.Conn
	js       jetstream.JetStream
	config   *config.Config
	creds    *auth.Credentials

	mu        sync.RWMutex
	state     core.ConnectionState
	lastError error

	lastUsed atomic.Int64

	// Callbacks
	onError func(error)
}

// NewConnection creates a new NATS connection wrapper
func NewConnection(tenantID string, cfg *config.Config, creds *auth.Credentials) *Connection {
	c := &Connection{
		tenantID: tenantID,
		config:   cfg,
		creds:    creds,
		state:    core.StateDisconnected,
	}
	c.Touch()
	return c
}

// Connect establishes the NATS connection
func (c *Connection) Connect(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "nats.connect",
		tracer.WithAttributes(
			tracer.StringAttr("tenant.id", c.tenantID),
		),
	)
	defer span.End()

	c.mu.Lock()
	if c.state == core.StateConnected {
		c.mu.Unlock()
		span.SetError(core.ErrAlreadyConnected)
		return core.ErrAlreadyConnected
	}
	c.state = core.StateConnecting
	c.mu.Unlock()

	opts := c.buildOptions()

	// Build server URL
	serverURL := strings.Join(c.config.Servers, ",")
	span.SetAttributes(tracer.StringAttr("servers", serverURL))

	logger.Debug(ctx, "connecting to NATS",
		logger.String("tenant_id", c.tenantID),
		logger.String("servers", serverURL),
	)

	// Connect with context timeout
	connectCtx, cancel := context.WithTimeout(ctx, c.config.ConnectTimeout)
	defer cancel()

	conn, err := nats.Connect(serverURL, opts...)
	if err != nil {
		c.mu.Lock()
		c.state = core.StateDisconnected
		c.mu.Unlock()
		span.SetError(err)
		logger.Error(ctx, "failed to connect to NATS",
			logger.String("tenant_id", c.tenantID),
			logger.Err(err),
		)
		return err
	}

	// Wait for connection
	select {
	case <-connectCtx.Done():
		conn.Close()
		c.mu.Lock()
		c.state = core.StateDisconnected
		c.mu.Unlock()
		span.SetError(connectCtx.Err())
		logger.Error(ctx, "NATS connection timeout",
			logger.String("tenant_id", c.tenantID),
			logger.Err(connectCtx.Err()),
		)
		return connectCtx.Err()
	default:
		if !conn.IsConnected() {
			time.Sleep(10 * time.Millisecond)
		}
	}

	c.mu.Lock()
	c.conn = conn
	c.state = core.StateConnected
	c.mu.Unlock()

	// Initialize JetStream if enabled
	if c.config.Stream.Enabled {
		var js jetstream.JetStream
		var jsErr error

		if c.config.Stream.Domain != "" {
			js, jsErr = jetstream.NewWithDomain(conn, c.config.Stream.Domain)
		} else if c.config.Stream.Prefix != "" {
			js, jsErr = jetstream.NewWithAPIPrefix(conn, c.config.Stream.Prefix)
		} else {
			js, jsErr = jetstream.New(conn)
		}

		if jsErr == nil {
			c.mu.Lock()
			c.js = js
			c.mu.Unlock()
			span.SetAttributes(tracer.BoolAttr("jetstream.enabled", true))
		} else {
			logger.Warn(ctx, "JetStream not available",
				logger.String("tenant_id", c.tenantID),
				logger.Err(jsErr),
			)
		}
	}

	c.Touch()

	span.SetOK()
	logger.Info(ctx, "NATS connection established",
		logger.String("tenant_id", c.tenantID),
		logger.String("connected_url", conn.ConnectedUrl()),
	)
	return nil
}

// Close gracefully closes the connection
func (c *Connection) Close(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "nats.close",
		tracer.WithAttributes(
			tracer.StringAttr("tenant.id", c.tenantID),
		),
	)
	defer span.End()

	c.mu.Lock()
	if c.state == core.StateClosed || c.state == core.StateDraining {
		c.mu.Unlock()
		span.SetOK()
		return nil
	}
	c.state = core.StateDraining
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
		}

		select {
		case <-drainCtx.Done():
			conn.Close()
		default:
		}
	}

	c.mu.Lock()
	c.state = core.StateClosed
	c.conn = nil
	c.js = nil
	c.mu.Unlock()

	span.SetOK()
	logger.Info(ctx, "NATS connection closed",
		logger.String("tenant_id", c.tenantID),
	)
	return nil
}

// IsConnected returns true if connected
func (c *Connection) IsConnected() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.state == core.StateConnected && c.conn != nil && c.conn.IsConnected()
}

// State returns the connection state
func (c *Connection) State() core.ConnectionState {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.state
}

// TenantID returns the tenant identifier
func (c *Connection) TenantID() string {
	return c.tenantID
}

// Touch updates the last activity timestamp
func (c *Connection) Touch() {
	c.lastUsed.Store(time.Now().Unix())
}

// IsIdle returns true if connection has been idle
func (c *Connection) IsIdle(timeout time.Duration) bool {
	lastUsed := time.Unix(c.lastUsed.Load(), 0)
	return time.Since(lastUsed) > timeout
}

// Health returns the health status
func (connection *Connection) Health() core.HealthStatus {

	connection.mu.RLock()
	defer connection.mu.RUnlock()

	healthy := connection.state == core.StateConnected && connection.conn != nil && connection.conn.IsConnected()

	// Build status message based on state
	var message string

	switch connection.state {

	case core.StateConnected:

		message = "connected to NATS server"

	case core.StateConnecting:

		message = "connecting to NATS server"

	case core.StateReconnecting:

		message = "reconnecting to NATS server"

	case core.StateDraining:

		message = "draining connection"

	case core.StateClosed:

		message = "connection closed"

	case core.StateDisconnected:

		message = "disconnected from NATS server"

	default:

		message = "unknown state"
	}

	// Build details
	details := map[string]any{
		"tenant_id": connection.tenantID,
		"servers":   connection.config.Servers,
	}

	if connection.conn != nil {

		details["connected_url"] = connection.conn.ConnectedUrl()
		details["reconnects"] = connection.conn.Stats().Reconnects
	}

	return core.HealthStatus{
		State:     connection.state,
		Healthy:   healthy,
		Message:   message,
		LastError: connection.lastError,
		Details:   details,
	}
}

// Conn returns the raw NATS connection
func (c *Connection) Conn() *nats.Conn {
	c.Touch()
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.conn
}

// JetStream returns the JetStream context
func (c *Connection) JetStream() jetstream.JetStream {
	c.Touch()
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.js
}

// OnError sets the error callback
func (c *Connection) OnError(cb func(error)) {
	c.onError = cb
}

// buildOptions builds NATS connection options
func (c *Connection) buildOptions() []nats.Option {
	opts := []nats.Option{
		nats.Name(c.config.Name + "-" + c.tenantID),
		nats.ReconnectWait(c.config.Reconnect.InitialInterval),
		nats.MaxReconnects(c.config.Reconnect.MaxAttempts),
		nats.ReconnectBufSize(8 * 1024 * 1024),
	}

	// Authentication
	if c.creds != nil {
		switch c.creds.Type {
		case auth.TypeUserPass:
			opts = append(opts, nats.UserInfo(c.creds.Username, c.creds.Password))
		case auth.TypeToken:
			opts = append(opts, nats.Token(c.creds.Token))
		case auth.TypeCredentialsFile:
			opts = append(opts, nats.UserCredentials(c.creds.CredentialsFile))
		case auth.TypeNKey:
			if c.creds.NKeyFile != "" {
				opt, err := nats.NkeyOptionFromSeed(c.creds.NKeyFile)
				if err == nil {
					opts = append(opts, opt)
				}
			}
		case auth.TypeJWT:
			if c.creds.JWT != "" {
				opts = append(opts, nats.UserJWTAndSeed(c.creds.JWT, c.creds.Seed))
			}
		}
	}

	// TLS
	if c.config.TLS != nil && c.config.TLS.Enabled {
		tlsConfig, tlsErr := c.buildTLSConfig()
		if tlsErr != nil {
			// Store error for later retrieval but continue - NATS might still connect
			c.mu.Lock()
			c.lastError = tlsErr
			c.mu.Unlock()
		}
		if tlsConfig != nil {
			opts = append(opts, nats.Secure(tlsConfig))
		}
	}

	// Event handlers
	opts = append(opts,
		nats.DisconnectErrHandler(c.handleDisconnect),
		nats.ReconnectHandler(c.handleReconnect),
		nats.ClosedHandler(c.handleClosed),
		nats.ErrorHandler(c.handleError),
	)

	return opts
}

func (c *Connection) buildTLSConfig() (*tls.Config, error) {
	if c.config.TLS == nil || !c.config.TLS.Enabled {
		return nil, nil
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.config.TLS.SkipVerify,
	}

	var errs []error

	if c.config.TLS.CertFile != "" && c.config.TLS.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(c.config.TLS.CertFile, c.config.TLS.KeyFile)
		if err != nil {
			errs = append(errs, core.NewError("tls", "load_cert", err))
		} else {
			tlsConfig.Certificates = []tls.Certificate{cert}
		}
	}

	if c.config.TLS.CAFile != "" {
		caCert, err := os.ReadFile(c.config.TLS.CAFile)
		if err != nil {
			errs = append(errs, core.NewError("tls", "load_ca", err))
		} else {
			caCertPool := x509.NewCertPool()
			if !caCertPool.AppendCertsFromPEM(caCert) {
				errs = append(errs, core.NewError("tls", "parse_ca", core.ErrInvalidConfig))
			} else {
				tlsConfig.RootCAs = caCertPool
			}
		}
	}

	if len(errs) > 0 {
		return tlsConfig, core.NewMultiError(errs)
	}

	return tlsConfig, nil
}

func (connection *Connection) handleDisconnect(conn *nats.Conn, err error) {
	ctx := context.Background()

	connection.mu.Lock()
	connection.state = core.StateReconnecting

	if err != nil {
		connection.lastError = err
	}
	connection.mu.Unlock()

	logger.Warn(ctx, "NATS connection disconnected",
		logger.String("tenant_id", connection.tenantID),
		logger.Err(err),
	)

	if err != nil && connection.onError != nil {
		connection.onError(err)
	}
}

func (connection *Connection) handleReconnect(conn *nats.Conn) {
	ctx := context.Background()

	connection.mu.Lock()
	connection.state = core.StateConnected
	connection.lastError = nil // Clear error on successful reconnect
	connection.mu.Unlock()

	connection.Touch()

	logger.Info(ctx, "NATS connection reconnected",
		logger.String("tenant_id", connection.tenantID),
		logger.String("connected_url", conn.ConnectedUrl()),
	)
}

func (connection *Connection) handleClosed(conn *nats.Conn) {
	ctx := context.Background()

	connection.mu.Lock()
	connection.state = core.StateClosed
	connection.mu.Unlock()

	logger.Info(ctx, "NATS connection closed by server",
		logger.String("tenant_id", connection.tenantID),
	)
}

func (connection *Connection) handleError(conn *nats.Conn, sub *nats.Subscription, err error) {
	ctx := context.Background()

	if err != nil {
		connection.mu.Lock()
		connection.lastError = err
		connection.mu.Unlock()

		subject := ""
		if sub != nil {
			subject = sub.Subject
		}
		logger.Error(ctx, "NATS error occurred",
			logger.String("tenant_id", connection.tenantID),
			logger.String("subject", subject),
			logger.Err(err),
		)
	}

	if connection.onError != nil {
		connection.onError(err)
	}
}
