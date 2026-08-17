// Package events provides a NATS-based event-driven messaging SDK.
//
// This package provides direct NATS integration for publishing and
// subscribing to messages with support for:
//   - NATS Core, JetStream, KV Store, and Object Store
//   - Multi-tenant connection management
//   - JWT-based authentication
//   - Middleware for tracing, retry, metrics, and logging
//   - Codec support via core/codec (custom binary, MsgPack)
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
//	    Builder: myConnectionBuilder,
//	})
//
//	// Connect tenant
//	creds := auth.JWT(jwtToken, seed)
//	tc, _ := manager.Connect(ctx, "tenant-1", creds)
//
//	// Use the NATS connection directly
//	tc.Conn().Publish("my.subject", []byte("hello"))
//
// # Package Structure
//
//   - core: Error types, context helpers, header constants, minimal NATS-typed interfaces
//   - config: Configuration management with env loading
//   - auth: Authentication providers and JWT credential management
//   - codec: Binary serialization via core/codec (custom binary, MsgPack)
//   - middleware: Cross-cutting concerns (tracing, retry, metrics, logging, circuit breaker, rate limiting)
//   - tenant: Multi-tenant connection management
//   - endpoint: Endpoint registration, discovery, and health checking
//   - microservice: Microservice registration and lifecycle management
package events

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/config"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/endpoint"
	jsmod "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/jetstream"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/microservice"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/middleware"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/tenant"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"
)

// Configuration type aliases re-exported from the config package.
// These allow consumers to use configuration types without importing config directly.
type (
	// Config is the top-level events/NATS configuration including server URLs,
	// TLS, reconnect behavior, JetStream settings, and timeouts.
	Config = config.EventsConfig

	// TLSConfig holds TLS/SSL settings for secure NATS connections.
	TLSConfig = config.TLSConfig

	// ReconnectConfig controls automatic reconnection behavior on connection loss.
	ReconnectConfig = config.ReconnectConfig

	// StreamConfig defines JetStream stream creation parameters such as
	// retention policy, storage type, and subject filters.
	StreamConfig = config.StreamConfig

	// PublishConfig holds settings for publish operations including
	// timeout, retry behavior, and expected stream targeting.
	PublishConfig = config.PublishConfig

	// RetryConfig defines retry parameters such as max attempts,
	// backoff strategy, and retryable error classification.
	RetryConfig = config.RetryConfig

	// SubscribeConfig holds settings for subscribe operations including
	// queue group, durable name, and delivery policy.
	SubscribeConfig = config.SubscribeConfig

	// KVConfig defines KeyValue store bucket settings such as
	// bucket name, TTL, history depth, and replicas.
	KVConfig = config.KVConfig

	// ObjectStoreConfig defines Object Store bucket settings such as
	// bucket name, chunk size, and max object size.
	ObjectStoreConfig = config.ObjectStoreConfig
)

// Authentication type aliases re-exported from the auth package.
// These provide access to credential types without importing auth directly.
type (
	// Credentials holds authentication details (type, username, password, token,
	// NKey seed, JWT, or credentials file path) used to connect to NATS.
	Credentials = auth.Credentials

	// AuthType identifies the authentication mechanism (none, user/pass, token,
	// NKey, credentials file, or JWT).
	AuthType = auth.Type

	// Provider supplies Credentials on demand, enabling dynamic credential
	// resolution at connection time.
	Provider = auth.Provider

	// CredentialManager manages JWT credentials for multiple tenants, handling
	// storage, expiry checking, automatic refresh, and optional file persistence.
	CredentialManager = auth.CredentialManager

	// RegistrationPayload carries the tenant ID, JWT, and seed needed to
	// register a new tenant's credentials in a single call.
	RegistrationPayload = auth.RegistrationPayload

	// RegistrationHandler processes incoming registration payloads and delegates
	// credential storage to a CredentialManager.
	RegistrationHandler = auth.RegistrationHandler

	// JWTCredential represents a parsed and validated JWT credential pair
	// (JWT token + NKey seed) for a specific tenant.
	JWTCredential = auth.JWTCredential
)

// Tenant type aliases re-exported from the tenant package.
// These support multi-tenant NATS connection management.
type (
	// TenantManager orchestrates tenant lifecycle: connecting, disconnecting,
	// health monitoring, and idle connection cleanup across all tenants.
	TenantManager = tenant.Manager

	// TenantConnection wraps a NATS connection and optional JetStream context
	// for a single tenant, tracking last-used time for idle detection.
	TenantConnection = tenant.TenantConnection

	// TenantRegistry maintains an in-memory index of registered tenant metadata,
	// enabling lookup and enumeration of known tenants.
	TenantRegistry = tenant.Registry

	// TenantInfo holds descriptive metadata about a tenant such as its ID,
	// display name, and connection status.
	TenantInfo = tenant.Info
)

// JetStream type aliases re-exported from the jetstream sub-package.
// These provide access to JetStream primitives without importing the internal module.
type (
	// StreamManager provides CRUD operations for JetStream streams including
	// creation, update, deletion, purging, and stream info retrieval.
	StreamManager = jsmod.StreamManager

	// KVStore wraps a JetStream KeyValue bucket, providing get/put/delete
	// operations, key watching, and history access.
	KVStore = jsmod.KVStore

	// ObjectStore wraps a JetStream Object Store bucket for storing and
	// retrieving large binary objects that exceed NATS message size limits.
	ObjectStore = jsmod.ObjectStore

	// JSConsumer wraps a JetStream consumer, providing message fetching,
	// acknowledgment, and consumer info retrieval.
	JSConsumer = jsmod.Consumer

	// OrderedConsumer wraps a JetStream ordered consumer that guarantees
	// in-order delivery without acknowledgments or redelivery.
	OrderedConsumer = jsmod.OrderedConsumer
)

// Middleware type aliases re-exported from the middleware package.
type (
	// MiddlewareStack is an ordered collection of publish/subscribe middleware
	// that wraps NATS operations with cross-cutting concerns such as tracing,
	// retry, metrics, and logging.
	MiddlewareStack = middleware.Stack
)

// Authentication type constants identifying the supported NATS authentication
// mechanisms. Use these when constructing Credentials to specify how the
// client should authenticate with the NATS server.
const (
	AuthTypeNone            = auth.TypeNone            // No authentication
	AuthTypeUserPass        = auth.TypeUserPass        // Username and password
	AuthTypeToken           = auth.TypeToken           // Bearer token
	AuthTypeNKey            = auth.TypeNKey            // NKey-based ed25519 signing
	AuthTypeCredentialsFile = auth.TypeCredentialsFile // Path to a .creds file
	AuthTypeJWT             = auth.TypeJWT             // JWT token with NKey seed
)

// Sentinel errors re-exported from utils for convenient error checking with
// errors.Is without needing to import the utils package directly.
var (
	ErrNotConnected        = utils.ErrNotConnected        // NATS connection has not been established
	ErrAlreadyConnected    = utils.ErrAlreadyConnected    // Tenant is already connected
	ErrConnectionClosed    = utils.ErrConnectionClosed    // Connection was closed or lost
	ErrConnectionTimeout   = utils.ErrConnectionTimeout   // Connection attempt exceeded the deadline
	ErrPublishFailed       = utils.ErrPublishFailed       // Message publish operation failed
	ErrPublishTimeout      = utils.ErrPublishTimeout      // Publish did not complete within the timeout
	ErrInvalidSubject      = utils.ErrInvalidSubject      // Subject string is empty or malformed
	ErrInvalidMessage      = utils.ErrInvalidMessage      // Message payload failed validation
	ErrTenantNotFound      = utils.ErrTenantNotFound      // No connection exists for the given tenant ID
	ErrShutdownInProgress  = utils.ErrShutdownInProgress  // Client is shutting down; no new operations accepted
	ErrJetStreamNotEnabled = utils.ErrJetStreamNotEnabled // JetStream is not available on the connected server
)

// NewConfig creates a new Config populated with sensible defaults for local
// development (e.g., nats://localhost:4222, no TLS). Callers should override
// fields like Servers and TLS before using in production.
func NewConfig() *Config {
	return config.DefaultEventsConfig()
}

// LoadConfigFromEnv loads a Config by reading environment variables such as
// NATS_URL, NATS_TLS_CERT, etc. This is the recommended approach for
// production deployments where configuration is injected via the environment.
func LoadConfigFromEnv() (*Config, error) {
	return config.LoadEventsFromEnv()
}

// DefaultPublishConfig returns a PublishConfig with production-safe defaults
// including a reasonable timeout and no retries enabled.
func DefaultPublishConfig() PublishConfig {
	return config.DefaultPublishConfig()
}

// DefaultSubscribeConfig returns a SubscribeConfig with production-safe
// defaults including no queue group and automatic acknowledgment.
func DefaultSubscribeConfig() SubscribeConfig {
	return config.DefaultSubscribeConfig()
}

// NewCredentialManager creates a new CredentialManager that stores, validates,
// and auto-refreshes JWT credentials for tenants. The manager runs a background
// goroutine (once Start is called) that checks for expiring credentials at the
// interval specified in cfg.
func NewCredentialManager(cfg auth.CredentialManagerConfig) *CredentialManager {
	return auth.NewCredentialManager(cfg)
}

// NewRegistrationHandler creates a new RegistrationHandler that processes
// incoming RegistrationPayload messages and delegates credential storage to
// the provided CredentialManager.
func NewRegistrationHandler(cm *CredentialManager) *RegistrationHandler {
	return auth.NewRegistrationHandler(cm)
}

// Authentication helper functions that construct Credentials for each supported
// auth mechanism. These are convenience aliases so callers can write
// events.JWTAuth(jwt, seed) instead of importing the auth package.
var (
	// UserPassAuth creates Credentials for username/password authentication.
	UserPassAuth = auth.UserPass

	// TokenAuth creates Credentials for bearer-token authentication.
	TokenAuth = auth.Token

	// NKeyAuth creates Credentials for NKey-based ed25519 authentication.
	NKeyAuth = auth.NKey

	// JWTAuth creates Credentials for JWT + NKey seed authentication.
	JWTAuth = auth.JWT
)

// Middleware constructor functions re-exported for convenience. Each returns a
// publish or subscribe middleware that can be added to a MiddlewareStack.
var (
	// TracingMiddleware adds OpenTelemetry span propagation to publish/subscribe operations.
	TracingMiddleware = middleware.Tracing

	// RetryMiddleware wraps publish operations with configurable retry logic
	// and exponential backoff.
	RetryMiddleware = middleware.Retry

	// MetricsMiddleware records publish/subscribe latency and error counters
	// via OpenTelemetry metrics.
	MetricsMiddleware = middleware.MetricsMiddleware

	// LoggingMiddleware emits structured log entries for every publish and
	// subscribe operation, including subject, duration, and outcome.
	LoggingMiddleware = middleware.Logging
)

// Circuit breaker type aliases re-exported from the middleware package.
// A circuit breaker prevents cascading failures by temporarily blocking
// operations to an unhealthy downstream after a configurable failure threshold.
type (
	// CircuitBreaker tracks failure counts and transitions between closed,
	// open, and half-open states for a single target.
	CircuitBreaker = middleware.CircuitBreaker

	// CircuitBreakerConfig holds thresholds and timeouts that control when
	// the circuit opens, how long it stays open, and how many trial requests
	// are allowed in the half-open state.
	CircuitBreakerConfig = middleware.CircuitBreakerConfig

	// CircuitState represents the current state of a circuit breaker
	// (closed, open, or half-open).
	CircuitState = middleware.CircuitState

	// MultiCircuitBreaker maintains independent circuit breakers keyed by
	// subject, enabling per-destination failure isolation.
	MultiCircuitBreaker = middleware.MultiCircuitBreaker
)

// Circuit breaker state constants representing the three possible states
// in the circuit breaker state machine.
const (
	// CircuitClosed means the circuit is healthy; all requests pass through.
	CircuitClosed = middleware.CircuitClosed

	// CircuitOpen means the failure threshold was exceeded; requests are
	// immediately rejected until the recovery timeout elapses.
	CircuitOpen = middleware.CircuitOpen

	// CircuitHalfOpen means the recovery timeout has elapsed; a limited
	// number of trial requests are allowed to probe for recovery.
	CircuitHalfOpen = middleware.CircuitHalfOpen
)

// NewCircuitBreaker creates a new CircuitBreaker initialized in the closed
// (healthy) state with the thresholds and timeouts specified in cfg.
func NewCircuitBreaker(cfg CircuitBreakerConfig) *CircuitBreaker {
	return middleware.NewCircuitBreaker(cfg)
}

// DefaultCircuitBreakerConfig returns a CircuitBreakerConfig with conservative
// defaults suitable for most workloads (e.g., 5 consecutive failures to open,
// 30-second recovery timeout).
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return middleware.DefaultCircuitBreakerConfig()
}

// CircuitBreakerMiddleware is a publish middleware constructor that wraps
// publish calls with circuit breaker protection, returning ErrCircuitOpen
// when the breaker is in the open state.
var CircuitBreakerMiddleware = middleware.CircuitBreakerMiddleware

// Rate limiter type aliases re-exported from the middleware package.
// Rate limiters control throughput by capping the number of operations
// allowed within a time window.
type (
	// RateLimiter implements a token-bucket rate limiter that caps the
	// overall publish/subscribe throughput.
	RateLimiter = middleware.RateLimiter

	// RateLimiterConfig holds rate limiter parameters such as requests
	// per second and burst size.
	RateLimiterConfig = middleware.RateLimiterConfig

	// PerSubjectRateLimiter maintains independent rate limiters keyed by
	// NATS subject, enabling per-topic throughput control.
	PerSubjectRateLimiter = middleware.PerSubjectRateLimiter

	// SlidingWindowLimiter tracks request counts over a sliding time window,
	// providing smoother rate limiting compared to fixed-window approaches.
	SlidingWindowLimiter = middleware.SlidingWindowLimiter
)

// NewRateLimiter creates a new RateLimiter with the throughput limits
// defined in cfg. The limiter is immediately active upon creation.
func NewRateLimiter(cfg RateLimiterConfig) *RateLimiter {
	return middleware.NewRateLimiter(cfg)
}

// DefaultRateLimiterConfig returns a RateLimiterConfig with permissive
// defaults suitable for development and testing.
func DefaultRateLimiterConfig() RateLimiterConfig {
	return middleware.DefaultRateLimiterConfig()
}

// Rate limiter middleware constructor functions. Each returns a publish
// middleware that enforces throughput limits on outgoing messages.
var (
	// RateLimitMiddleware rejects requests immediately with ErrRateLimitExceeded
	// when the rate limit is exceeded.
	RateLimitMiddleware = middleware.RateLimitMiddleware

	// RateLimitWaitMiddleware blocks the caller until a token becomes available
	// rather than rejecting the request, providing backpressure.
	RateLimitWaitMiddleware = middleware.RateLimitWaitMiddleware

	// PerSubjectRateLimitMiddleware applies independent rate limits per NATS
	// subject, preventing a high-throughput topic from starving others.
	PerSubjectRateLimitMiddleware = middleware.PerSubjectRateLimitMiddleware
)

// Middleware sentinel errors re-exported for convenient error checking.
var (
	// ErrCircuitOpen is returned when a publish is attempted while the
	// circuit breaker is in the open state.
	ErrCircuitOpen = middleware.ErrCircuitOpen

	// ErrRateLimitExceeded is returned when a publish is rejected because
	// the rate limit has been exceeded.
	ErrRateLimitExceeded = middleware.ErrRateLimitExceeded
)

// NewMiddlewareStack creates a new middleware stack
func NewMiddlewareStack() *MiddlewareStack {
	return middleware.NewStack()
}

// NewTenantManager creates a new tenant manager
func NewTenantManager(cfg *Config, builder tenant.ConnectionBuilder, credManager *CredentialManager) (*TenantManager, error) {
	return tenant.NewManager(tenant.ManagerConfig{
		Config:      cfg,
		Builder:     builder,
		CredManager: credManager,
	})
}

// NewTenantRegistry creates a new tenant registry
func NewTenantRegistry() *TenantRegistry {
	return tenant.NewRegistry()
}

// Context helpers
var (
	WithTenantID             = core.WithTenantID
	TenantIDFromContext      = core.TenantIDFromContext
	WithTraceContext         = core.WithTraceContext
	TraceContextFromContext  = core.TraceContextFromContext
	WithMessageID            = core.WithMessageID
	MessageIDFromContext     = core.MessageIDFromContext
	WithCorrelationID        = core.WithCorrelationID
	CorrelationIDFromContext = core.CorrelationIDFromContext
)

// Client is a convenience wrapper for common operations
type Client struct {
	manager     *TenantManager
	credManager *CredentialManager
	config      *Config
	builder     tenant.ConnectionBuilder
}

// ClientConfig holds client configuration options
type ClientConfig struct {
	// Config is the NATS/events configuration
	Config *Config

	// Builder creates NATS connections for tenants
	Builder tenant.ConnectionBuilder
}

// NewClient creates a new events client
func NewClient(cfg ClientConfig, opts ...config.EventsOption) (*Client, error) {
	eventsConfiguration := cfg.Config
	if eventsConfiguration == nil {
		eventsConfiguration = config.DefaultEventsConfig()
	}
	config.ApplyEventsOptions(eventsConfiguration, opts...)

	if validationError := config.ValidateEventsConfig(eventsConfiguration); validationError != nil {
		return nil, validationError
	}

	credentialManager := auth.NewCredentialManager(auth.DefaultCredentialManagerConfig())

	if cfg.Builder == nil {
		return nil, utils.NewError("client", "create", utils.ErrInvalidConfig)
	}

	tenantManager, createError := NewTenantManager(eventsConfiguration, cfg.Builder, credentialManager)
	if createError != nil {
		return nil, createError
	}

	return &Client{
		manager:     tenantManager,
		credManager: credentialManager,
		config:      eventsConfiguration,
		builder:     cfg.Builder,
	}, nil
}

// Connect connects a tenant
func (c *Client) Connect(ctx context.Context, tenantID string, creds *Credentials) (*TenantConnection, error) {
	return c.manager.Connect(ctx, tenantID, creds)
}

// ConnectWithJWT connects a tenant using JWT credentials
func (c *Client) ConnectWithJWT(ctx context.Context, tenantID, jwt, seed string) (*TenantConnection, error) {
	creds := auth.JWT(jwt, seed)
	return c.manager.Connect(ctx, tenantID, creds)
}

// ConnectWithPayload connects using a registration payload
func (c *Client) ConnectWithPayload(ctx context.Context, payload RegistrationPayload) (*TenantConnection, error) {
	return c.manager.ConnectWithPayload(ctx, payload)
}

// RegisterCredentials registers credentials for a tenant
func (c *Client) RegisterCredentials(tenantID, jwt, seed string) error {
	return c.credManager.Register(tenantID, jwt, seed)
}

// Get returns an existing connection
func (c *Client) Get(tenantID string) (*TenantConnection, bool) {
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

// CredManager returns the credential manager
func (c *Client) CredManager() *CredentialManager {
	return c.credManager
}

// ClientConfig returns the configuration
func (c *Client) ClientConfig() *Config {
	return c.config
}

// Endpoint types
type (
	EndpointInfo         = endpoint.Info
	EndpointStatus       = endpoint.Status
	EndpointProtocol     = endpoint.Protocol
	EndpointRegistry     = endpoint.Registry
	EndpointManager      = endpoint.Manager
	EndpointQuery        = endpoint.Query
	EndpointEvent        = endpoint.Event
	EndpointHandler      = endpoint.Handler
	EndpointClient       = endpoint.Client
	EndpointHealthUpdate = endpoint.HealthUpdate
)

// Microservice types
type (
	MicroserviceInfo     = microservice.Info
	MicroserviceStatus   = microservice.Status
	MicroserviceType     = microservice.ServiceType
	MicroserviceInstance = microservice.Instance
	MicroserviceRegistry = microservice.Registry
	MicroserviceManager  = microservice.Manager
	MicroserviceQuery    = microservice.Query
	MicroserviceEvent    = microservice.Event
	MicroserviceHandler  = microservice.Handler
	MicroserviceClient   = microservice.Client
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

// ============================================================================
// JetStream convenience functions
// ============================================================================

// NewStreamManager creates a stream manager from a JetStream instance.
func NewStreamManager(js jetstream.JetStream) (*StreamManager, error) {
	return jsmod.NewStreamManager(js)
}

// CreateKVStore creates a new KeyValue store bucket.
func CreateKVStore(ctx context.Context, js jetstream.JetStream, cfg jetstream.KeyValueConfig) (*KVStore, error) {
	return jsmod.CreateKVStore(ctx, js, cfg)
}

// GetKVStore retrieves an existing KeyValue store bucket.
func GetKVStore(ctx context.Context, js jetstream.JetStream, bucket string) (*KVStore, error) {
	return jsmod.GetKVStore(ctx, js, bucket)
}

// DeleteKVStore deletes a KeyValue store bucket.
func DeleteKVStore(ctx context.Context, js jetstream.JetStream, bucket string) error {
	return jsmod.DeleteKVStore(ctx, js, bucket)
}

// KVStoreNames returns the names of all KeyValue store buckets.
func KVStoreNames(ctx context.Context, js jetstream.JetStream) ([]string, error) {
	return jsmod.KVStoreNames(ctx, js)
}

// CreateObjectStore creates a new Object Store bucket.
func CreateObjectStore(ctx context.Context, js jetstream.JetStream, cfg jetstream.ObjectStoreConfig) (*ObjectStore, error) {
	return jsmod.CreateObjectStore(ctx, js, cfg)
}

// GetObjectStore retrieves an existing Object Store bucket.
func GetObjectStore(ctx context.Context, js jetstream.JetStream, bucket string) (*ObjectStore, error) {
	return jsmod.GetObjectStore(ctx, js, bucket)
}

// DeleteObjectStore deletes an Object Store bucket.
func DeleteObjectStore(ctx context.Context, js jetstream.JetStream, bucket string) error {
	return jsmod.DeleteObjectStore(ctx, js, bucket)
}

// ObjectStoreNames returns the names of all Object Store buckets.
func ObjectStoreNames(ctx context.Context, js jetstream.JetStream) ([]string, error) {
	return jsmod.ObjectStoreNames(ctx, js)
}

// CreateOrUpdateConsumer creates or updates a JetStream consumer.
func CreateOrUpdateConsumer(ctx context.Context, js jetstream.JetStream, stream string, cfg jetstream.ConsumerConfig) (*JSConsumer, error) {
	return jsmod.CreateOrUpdateConsumer(ctx, js, stream, cfg)
}

// GetConsumer retrieves an existing JetStream consumer.
func GetConsumer(ctx context.Context, js jetstream.JetStream, stream, name string) (*JSConsumer, error) {
	return jsmod.GetConsumer(ctx, js, stream, name)
}

// DeleteConsumer deletes a JetStream consumer.
func DeleteConsumer(ctx context.Context, js jetstream.JetStream, stream, name string) error {
	return jsmod.DeleteConsumer(ctx, js, stream, name)
}

// CreateOrderedConsumer creates an ordered JetStream consumer.
func CreateOrderedConsumer(ctx context.Context, js jetstream.JetStream, stream string, cfg jetstream.OrderedConsumerConfig) (*OrderedConsumer, error) {
	return jsmod.CreateOrderedConsumer(ctx, js, stream, cfg)
}
