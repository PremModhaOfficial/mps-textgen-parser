// Package events provides a generalized event-driven messaging SDK.
//
// This package provides a transport-agnostic interface for publishing and
// subscribing to messages with support for:
//   - Multiple transport backends (NATS, Kafka, etc.)
//   - Multi-tenant connection management
//   - JWT-based authentication
//   - Middleware for tracing, retry, metrics, and logging
//   - Multiple encoding formats (JSON, Protobuf, bytes)
//   - Endpoint registration and discovery
//   - Microservice registration and lifecycle management
//
// # Basic Usage
//
//	// Create configuration
//	cfg := config.DefaultConfig()
//	cfg.Servers = []string{"nats://localhost:4222"}
//
//	// Create tenant manager
//	manager, _ := tenant.NewManager(tenant.ManagerConfig{
//	    Config:  cfg,
//	    Factory: nats.NewFactory(),
//	})
//
//	// Connect tenant
//	creds := auth.JWT(jwtToken, seed)
//	conn, _ := manager.Connect(ctx, "tenant-1", creds)
//
//	// Publish message
//	publisher := conn.Publisher()
//	msg := core.NewMessage([]byte(`{"hello": "world"}`))
//	publisher.Publish(ctx, "my.subject", msg)
//
// # Endpoint Registration
//
//	// Create endpoint manager
//	epManager, _ := endpoint.NewManager(endpoint.ManagerConfig{})
//
//	// Register an endpoint
//	ep := &endpoint.Info{
//	    ID:        "user-api-1",
//	    ServiceID: "user-service",
//	    Host:      "localhost",
//	    Port:      8080,
//	    Path:      "/api/v1/users",
//	    Protocol:  endpoint.ProtocolHTTP,
//	}
//	epManager.RegisterEndpoint(ctx, ep)
//
// # Microservice Registration
//
//	// Create microservice manager
//	svcManager, _ := microservice.NewManager(microservice.ManagerConfig{})
//
//	// Register a service
//	svc := &microservice.Info{
//	    ID:      "user-service",
//	    Name:    "user-service",
//	    Version: "1.0.0",
//	    Type:    microservice.ServiceTypeAPI,
//	}
//	svcManager.RegisterServiceDirect(ctx, svc)
//
// # Package Structure
//
//   - core: Core interfaces and types (Message, Connection, Publisher, Subscriber)
//   - config: Configuration management with env loading
//   - auth: Authentication providers and JWT credential management
//   - encoding: Message serialization (JSON, Protobuf, bytes)
//   - middleware: Cross-cutting concerns (tracing, retry, metrics, logging)
//   - tenant: Multi-tenant connection management
//   - transport: Transport implementations (NATS, etc.)
//   - endpoint: Endpoint registration, discovery, and health checking
//   - microservice: Microservice registration and lifecycle management
package events

import (
	"context"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/config"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/encoding"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/endpoint"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/microservice"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/middleware"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/tenant"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/transport/nats"
)

// Re-export commonly used types for convenience

// Core types
type (
	Message         = core.Message
	Headers         = core.Headers
	Connection      = core.Connection
	Publisher       = core.Publisher
	Subscriber      = core.Subscriber
	Subscription    = core.Subscription
	MessageHandler  = core.MessageHandler
	PubAck          = core.PubAck
	PubAckFuture    = core.PubAckFuture
	ConnectionState = core.ConnectionState
	HealthStatus    = core.HealthStatus
	Component       = core.Component
)

// Config types
type (
	Config          = config.Config
	TLSConfig       = config.TLSConfig
	ReconnectConfig = config.ReconnectConfig
	StreamConfig    = config.StreamConfig
	PublishConfig   = config.PublishConfig
	RetryConfig     = config.RetryConfig
	SubscribeConfig = config.SubscribeConfig
)

// Auth types
type (
	Credentials         = auth.Credentials
	AuthType            = auth.Type
	Provider            = auth.Provider
	CredentialManager   = auth.CredentialManager
	RegistrationPayload = auth.RegistrationPayload
	RegistrationHandler = auth.RegistrationHandler
	JWTCredential       = auth.JWTCredential
)

// Tenant types
type (
	TenantManager    = tenant.Manager
	TenantConnection = tenant.Connection
	TenantRegistry   = tenant.Registry
	TenantInfo       = tenant.Info
)

// Middleware types
type (
	PublishMiddleware   = middleware.PublishMiddleware
	SubscribeMiddleware = middleware.SubscribeMiddleware
	MiddlewareStack     = middleware.Stack
)

// Connection states
const (
	StateDisconnected = core.StateDisconnected
	StateConnecting   = core.StateConnecting
	StateConnected    = core.StateConnected
	StateReconnecting = core.StateReconnecting
	StateDraining     = core.StateDraining
	StateClosed       = core.StateClosed
)

// Auth types
const (
	AuthTypeNone            = auth.TypeNone
	AuthTypeUserPass        = auth.TypeUserPass
	AuthTypeToken           = auth.TypeToken
	AuthTypeNKey            = auth.TypeNKey
	AuthTypeCredentialsFile = auth.TypeCredentialsFile
	AuthTypeJWT             = auth.TypeJWT
)

// Common errors
var (
	ErrNotConnected       = core.ErrNotConnected
	ErrAlreadyConnected   = core.ErrAlreadyConnected
	ErrConnectionClosed   = core.ErrConnectionClosed
	ErrConnectionTimeout  = core.ErrConnectionTimeout
	ErrPublishFailed      = core.ErrPublishFailed
	ErrPublishTimeout     = core.ErrPublishTimeout
	ErrInvalidSubject     = core.ErrInvalidSubject
	ErrInvalidMessage     = core.ErrInvalidMessage
	ErrTenantNotFound     = core.ErrTenantNotFound
	ErrShutdownInProgress = core.ErrShutdownInProgress
)

// NewMessage creates a new message
func NewMessage(data []byte) *Message {
	return core.NewMessage(data)
}

// Message pool functions for reduced allocations
var (
	AcquireMessage         = core.AcquireMessage
	AcquireMessageWithData = core.AcquireMessageWithData
	ReleaseMessage         = core.ReleaseMessage
	AcquireHeaders         = core.AcquireHeaders
	ReleaseHeaders         = core.ReleaseHeaders
)

// Pool types
type (
	MessagePool    = core.MessagePool
	HeadersPool    = core.HeadersPool
	ByteBufferPool = core.ByteBufferPool
)

// NewMessagePool creates a new message pool
func NewMessagePool() *MessagePool {
	return core.NewMessagePool()
}

// NewConfig creates a new configuration with defaults
func NewConfig() *Config {
	return config.DefaultConfig()
}

// LoadConfigFromEnv loads configuration from environment variables
func LoadConfigFromEnv() (*Config, error) {
	return config.LoadFromEnv()
}

// DefaultPublishConfig returns default publish configuration
func DefaultPublishConfig() PublishConfig {
	return config.DefaultPublishConfig()
}

// DefaultSubscribeConfig returns default subscribe configuration
func DefaultSubscribeConfig() SubscribeConfig {
	return config.DefaultSubscribeConfig()
}

// NewCredentialManager creates a new credential manager
func NewCredentialManager(cfg auth.CredentialManagerConfig) *CredentialManager {
	return auth.NewCredentialManager(cfg)
}

// NewRegistrationHandler creates a new registration handler
func NewRegistrationHandler(cm *CredentialManager) *RegistrationHandler {
	return auth.NewRegistrationHandler(cm)
}

// Auth helper functions
var (
	UserPassAuth = auth.UserPass
	TokenAuth    = auth.Token
	NKeyAuth     = auth.NKey
	JWTAuth      = auth.JWT
)

// Encoding functions
var (
	EncodeJSON     = encoding.EncodeJSON
	DecodeJSON     = encoding.DecodeJSON
	EncodeBytes    = encoding.EncodeBytes
	DecodeBytes    = encoding.DecodeBytes
	EncodeProtobuf = encoding.EncodeProtobuf
	DecodeProtobuf = encoding.DecodeProtobuf
)

// Encoders
var (
	JSONEncoder     = encoding.JSON
	BytesEncoder    = encoding.Bytes
	ProtobufEncoder = encoding.Protobuf
)

// Middleware constructors
var (
	TracingMiddleware = middleware.Tracing
	RetryMiddleware   = middleware.Retry
	MetricsMiddleware = middleware.MetricsMiddleware
	LoggingMiddleware = middleware.Logging
)

// Circuit breaker types and functions
type (
	CircuitBreaker       = middleware.CircuitBreaker
	CircuitBreakerConfig = middleware.CircuitBreakerConfig
	CircuitState         = middleware.CircuitState
	MultiCircuitBreaker  = middleware.MultiCircuitBreaker
)

// Circuit breaker constants
const (
	CircuitClosed   = middleware.CircuitClosed
	CircuitOpen     = middleware.CircuitOpen
	CircuitHalfOpen = middleware.CircuitHalfOpen
)

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(cfg CircuitBreakerConfig) *CircuitBreaker {
	return middleware.NewCircuitBreaker(cfg)
}

// DefaultCircuitBreakerConfig returns default circuit breaker configuration
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return middleware.DefaultCircuitBreakerConfig()
}

// CircuitBreakerMiddleware returns publish middleware with circuit breaker
var CircuitBreakerMiddleware = middleware.CircuitBreakerMiddleware

// Rate limiter types and functions
type (
	RateLimiter           = middleware.RateLimiter
	RateLimiterConfig     = middleware.RateLimiterConfig
	PerSubjectRateLimiter = middleware.PerSubjectRateLimiter
	SlidingWindowLimiter  = middleware.SlidingWindowLimiter
)

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(cfg RateLimiterConfig) *RateLimiter {
	return middleware.NewRateLimiter(cfg)
}

// DefaultRateLimiterConfig returns default rate limiter configuration
func DefaultRateLimiterConfig() RateLimiterConfig {
	return middleware.DefaultRateLimiterConfig()
}

// Rate limiter middleware functions
var (
	RateLimitMiddleware          = middleware.RateLimitMiddleware
	RateLimitWaitMiddleware      = middleware.RateLimitWaitMiddleware
	PerSubjectRateLimitMiddleware = middleware.PerSubjectRateLimitMiddleware
)

// Middleware errors
var (
	ErrCircuitOpen       = middleware.ErrCircuitOpen
	ErrRateLimitExceeded = middleware.ErrRateLimitExceeded
)

// NewMiddlewareStack creates a new middleware stack
func NewMiddlewareStack() *MiddlewareStack {
	return middleware.NewStack()
}

// NewTenantManager creates a new tenant manager with NATS transport
func NewTenantManager(cfg *Config, credManager *CredentialManager) (*TenantManager, error) {
	return tenant.NewManager(tenant.ManagerConfig{
		Config:      cfg,
		Factory:     nats.NewFactory(),
		CredManager: credManager,
	})
}

// NewTenantRegistry creates a new tenant registry
func NewTenantRegistry() *TenantRegistry {
	return tenant.NewRegistry()
}

// Context helpers
var (
	WithTenantID          = core.WithTenantID
	TenantIDFromContext   = core.TenantIDFromContext
	WithTraceContext      = core.WithTraceContext
	TraceContextFromContext = core.TraceContextFromContext
	WithMessageID         = core.WithMessageID
	MessageIDFromContext  = core.MessageIDFromContext
	WithCorrelationID     = core.WithCorrelationID
	CorrelationIDFromContext = core.CorrelationIDFromContext
)

// Client is a convenience wrapper for common operations
type Client struct {
	manager     *TenantManager
	credManager *CredentialManager
	config      *Config
}

// NewClient creates a new events client
func NewClient(opts ...config.Option) (*Client, error) {
	eventsConfiguration := config.DefaultConfig()
	config.Apply(eventsConfiguration, opts...)

	if validationError := config.Validate(eventsConfiguration); validationError != nil {
		return nil, validationError
	}

	credentialManager := auth.NewCredentialManager(auth.DefaultCredentialManagerConfig())

	tenantManager, createError := NewTenantManager(eventsConfiguration, credentialManager)
	if createError != nil {
		return nil, createError
	}

	return &Client{
		manager:     tenantManager,
		credManager: credentialManager,
		config:      eventsConfiguration,
	}, nil
}

// Connect connects a tenant
func (c *Client) Connect(ctx context.Context, tenantID string, creds *Credentials) (TenantConnection, error) {
	return c.manager.Connect(ctx, tenantID, creds)
}

// ConnectWithJWT connects a tenant using JWT credentials
func (c *Client) ConnectWithJWT(ctx context.Context, tenantID, jwt, seed string) (TenantConnection, error) {
	creds := auth.JWT(jwt, seed)
	return c.manager.Connect(ctx, tenantID, creds)
}

// ConnectWithPayload connects using a registration payload
func (c *Client) ConnectWithPayload(ctx context.Context, payload RegistrationPayload) (TenantConnection, error) {
	return c.manager.ConnectWithPayload(ctx, payload)
}

// RegisterCredentials registers credentials for a tenant
func (c *Client) RegisterCredentials(tenantID, jwt, seed string) error {
	return c.credManager.Register(tenantID, jwt, seed)
}

// Get returns an existing connection
func (c *Client) Get(tenantID string) (TenantConnection, bool) {
	return c.manager.Get(tenantID)
}

// Disconnect disconnects a tenant
func (c *Client) Disconnect(ctx context.Context, tenantID string) error {
	return c.manager.Disconnect(ctx, tenantID)
}

// Shutdown shuts down the client
func (c *Client) Shutdown(ctx context.Context) error {
	c.credManager.Stop()
	return c.manager.Shutdown(ctx)
}

// Manager returns the tenant manager
func (c *Client) Manager() *TenantManager {
	return c.manager
}

// CredentialManager returns the credential manager
func (c *Client) CredentialManager() *CredentialManager {
	return c.credManager
}

// Config returns the configuration
func (c *Client) Config() *Config {
	return c.config
}

// Endpoint types
type (
	EndpointInfo           = endpoint.Info
	EndpointStatus         = endpoint.Status
	EndpointProtocol       = endpoint.Protocol
	EndpointRegistry       = endpoint.Registry
	EndpointManager        = endpoint.Manager
	EndpointQuery          = endpoint.Query
	EndpointEvent          = endpoint.Event
	EndpointHandler        = endpoint.Handler
	EndpointClient         = endpoint.Client
	EndpointHealthUpdate   = endpoint.HealthUpdate
)

// Microservice types
type (
	MicroserviceInfo       = microservice.Info
	MicroserviceStatus     = microservice.Status
	MicroserviceType       = microservice.ServiceType
	MicroserviceInstance   = microservice.Instance
	MicroserviceRegistry   = microservice.Registry
	MicroserviceManager    = microservice.Manager
	MicroserviceQuery      = microservice.Query
	MicroserviceEvent      = microservice.Event
	MicroserviceHandler    = microservice.Handler
	MicroserviceClient     = microservice.Client
)

// Endpoint status constants
const (
	EndpointStatusUnknown     = endpoint.StatusUnknown
	EndpointStatusHealthy     = endpoint.StatusHealthy
	EndpointStatusUnhealthy   = endpoint.StatusUnhealthy
	EndpointStatusDegraded    = endpoint.StatusDegraded
	EndpointStatusMaintenance = endpoint.StatusMaintenance
	EndpointStatusOffline     = endpoint.StatusOffline
)

// Endpoint protocol constants
const (
	EndpointProtocolHTTP  = endpoint.ProtocolHTTP
	EndpointProtocolHTTPS = endpoint.ProtocolHTTPS
	EndpointProtocolGRPC  = endpoint.ProtocolGRPC
	EndpointProtocolWS    = endpoint.ProtocolWS
	EndpointProtocolWSS   = endpoint.ProtocolWSS
	EndpointProtocolTCP   = endpoint.ProtocolTCP
	EndpointProtocolUDP   = endpoint.ProtocolUDP
)

// Microservice status constants
const (
	ServiceStatusUnknown     = microservice.StatusUnknown
	ServiceStatusStarting    = microservice.StatusStarting
	ServiceStatusRunning     = microservice.StatusRunning
	ServiceStatusDegraded    = microservice.StatusDegraded
	ServiceStatusStopping    = microservice.StatusStopping
	ServiceStatusStopped     = microservice.StatusStopped
	ServiceStatusFailed      = microservice.StatusFailed
	ServiceStatusMaintenance = microservice.StatusMaintenance
)

// Microservice type constants
const (
	ServiceTypeAPI       = microservice.ServiceTypeAPI
	ServiceTypeWorker    = microservice.ServiceTypeWorker
	ServiceTypeGateway   = microservice.ServiceTypeGateway
	ServiceTypeBroker    = microservice.ServiceTypeBroker
	ServiceTypeDatabase  = microservice.ServiceTypeDatabase
	ServiceTypeCache     = microservice.ServiceTypeCache
	ServiceTypeQueue     = microservice.ServiceTypeQueue
	ServiceTypeScheduler = microservice.ServiceTypeScheduler
	ServiceTypeMonitor   = microservice.ServiceTypeMonitor
	ServiceTypeCustom    = microservice.ServiceTypeCustom
)

// Endpoint errors
var (
	ErrEndpointNotFound  = endpoint.ErrEndpointNotFound
	ErrEndpointExists    = endpoint.ErrEndpointExists
	ErrInvalidEndpoint   = endpoint.ErrInvalidEndpoint
	ErrMissingEndpointID = endpoint.ErrMissingEndpointID
)

// Microservice errors
var (
	ErrServiceNotFound   = microservice.ErrServiceNotFound
	ErrServiceExists     = microservice.ErrServiceExists
	ErrInvalidService    = microservice.ErrInvalidService
	ErrMissingServiceID  = microservice.ErrMissingServiceID
	ErrInstanceNotFound  = microservice.ErrInstanceNotFound
	ErrInstanceExists    = microservice.ErrInstanceExists
	ErrInvalidInstance   = microservice.ErrInvalidInstance
	ErrMissingInstanceID = microservice.ErrMissingInstanceID
)

// NewEndpointRegistry creates a new endpoint registry
func NewEndpointRegistry(cfg ...endpoint.RegistryConfig) *EndpointRegistry {
	return endpoint.NewRegistry(cfg...)
}

// NewEndpointManager creates a new endpoint manager
func NewEndpointManager(cfg endpoint.ManagerConfig) (*EndpointManager, error) {
	return endpoint.NewManager(cfg)
}

// NewEndpointHandler creates a new endpoint handler
func NewEndpointHandler(manager *EndpointManager, subjects ...endpoint.Subjects) *EndpointHandler {
	var opts []endpoint.HandlerOption
	if len(subjects) > 0 {
		opts = append(opts, endpoint.WithSubjects(subjects[0]))
	}
	return endpoint.NewHandler(manager, opts...)
}

// NewEndpointClient creates a new endpoint client
func NewEndpointClient(publisher core.Publisher, subjects ...endpoint.Subjects) *EndpointClient {
	return endpoint.NewClient(publisher, subjects...)
}

// NewMicroserviceRegistry creates a new microservice registry
func NewMicroserviceRegistry(cfg ...microservice.RegistryConfig) *MicroserviceRegistry {
	return microservice.NewRegistry(cfg...)
}

// NewMicroserviceManager creates a new microservice manager
func NewMicroserviceManager(cfg microservice.ManagerConfig) (*MicroserviceManager, error) {
	return microservice.NewManager(cfg)
}

// NewMicroserviceHandler creates a new microservice handler
func NewMicroserviceHandler(manager *MicroserviceManager, subjects ...microservice.Subjects) *MicroserviceHandler {
	var opts []microservice.HandlerOption
	if len(subjects) > 0 {
		opts = append(opts, microservice.WithSubjects(subjects[0]))
	}
	return microservice.NewHandler(manager, opts...)
}

// NewMicroserviceClient creates a new microservice client
func NewMicroserviceClient(publisher core.Publisher, subjects ...microservice.Subjects) *MicroserviceClient {
	return microservice.NewClient(publisher, subjects...)
}

// DefaultEndpointSubjects returns default endpoint subjects
func DefaultEndpointSubjects() endpoint.Subjects {
	return endpoint.DefaultSubjects()
}

// DefaultMicroserviceSubjects returns default microservice subjects
func DefaultMicroserviceSubjects() microservice.Subjects {
	return microservice.DefaultSubjects()
}
