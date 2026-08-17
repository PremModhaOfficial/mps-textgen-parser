package utils

import "time"

// ============================================================================
// Common String Constants
// ============================================================================

const (
	// Empty is an empty string constant
	Empty = ""

	// DefaultSeparator is the default separator for subject tokens
	DefaultSeparator = "."

	// WildcardSingle matches any single token in a subject
	WildcardSingle = "*"

	// WildcardMulti matches any number of tokens (must be last)
	WildcardMulti = ">"
)

// ============================================================================
// Default Values
// ============================================================================

const (
	// DefaultTimeout is the default timeout for operations
	DefaultTimeout = 30 * time.Second

	// DefaultConnectTimeout is the default connection timeout
	DefaultConnectTimeout = 10 * time.Second

	// DefaultRequestTimeout is the default request-reply timeout
	DefaultRequestTimeout = 5 * time.Second

	// DefaultReconnectWait is the default wait time between reconnection attempts
	DefaultReconnectWait = 2 * time.Second

	// DefaultMaxReconnects is the default maximum reconnection attempts (-1 for unlimited)
	DefaultMaxReconnects = -1

	// DefaultHealthCheckInterval is the default interval for health checks
	DefaultHealthCheckInterval = 30 * time.Second

	// DefaultHealthCheckTimeout is the default timeout for health checks
	DefaultHealthCheckTimeout = 5 * time.Second

	// DefaultBufferSize is the default buffer size for channels
	DefaultBufferSize = 256

	// DefaultMaxPending is the default maximum pending messages
	DefaultMaxPending = 65536

	// DefaultMaxPayload is the default maximum message payload size (1MB)
	DefaultMaxPayload = 1024 * 1024

	// DefaultBatchSize is the default batch size for batch operations
	DefaultBatchSize = 100
)

// ============================================================================
// Subject Prefixes
// ============================================================================

const (
	// SubjectPrefixService is the prefix for service-related subjects
	SubjectPrefixService = "service"

	// SubjectPrefixEndpoint is the prefix for endpoint-related subjects
	SubjectPrefixEndpoint = "endpoint"

	// SubjectPrefixEvent is the prefix for event subjects
	SubjectPrefixEvent = "events"

	// SubjectPrefixRequest is the prefix for request-reply subjects
	SubjectPrefixRequest = "request"

	// SubjectPrefixBroadcast is the prefix for broadcast subjects
	SubjectPrefixBroadcast = "broadcast"

	// SubjectPrefixInternal is the prefix for internal subjects
	SubjectPrefixInternal = "_internal"
)

// ============================================================================
// Event Types
// ============================================================================

const (
	// EventTypeCreated indicates a resource was created
	EventTypeCreated = "created"

	// EventTypeUpdated indicates a resource was updated
	EventTypeUpdated = "updated"

	// EventTypeDeleted indicates a resource was deleted
	EventTypeDeleted = "deleted"

	// EventTypeStarted indicates a process/service started
	EventTypeStarted = "started"

	// EventTypeStopped indicates a process/service stopped
	EventTypeStopped = "stopped"

	// EventTypeHealthy indicates a resource is healthy
	EventTypeHealthy = "healthy"

	// EventTypeUnhealthy indicates a resource is unhealthy
	EventTypeUnhealthy = "unhealthy"

	// EventTypeConnected indicates a connection was established
	EventTypeConnected = "connected"

	// EventTypeDisconnected indicates a connection was lost
	EventTypeDisconnected = "disconnected"

	// EventTypeError indicates an error occurred
	EventTypeError = "error"
)

// ============================================================================
// Content Types
// ============================================================================

const (
	// ContentTypeJSON is the MIME type for JSON content
	ContentTypeJSON = "application/json"

	// ContentTypeProtobuf is the MIME type for Protocol Buffers content
	ContentTypeProtobuf = "application/protobuf"

	// ContentTypeMsgPack is the MIME type for MessagePack content
	ContentTypeMsgPack = "application/msgpack"

	// ContentTypeText is the MIME type for plain text content
	ContentTypeText = "text/plain"

	// ContentTypeBinary is the MIME type for binary content
	ContentTypeBinary = "application/octet-stream"
)

// ============================================================================
// Retry Configuration
// ============================================================================

const (
	// DefaultMaxRetries is the default number of retry attempts
	DefaultMaxRetries = 3

	// DefaultRetryDelay is the default delay between retries
	DefaultRetryDelay = 100 * time.Millisecond

	// DefaultRetryMaxDelay is the default maximum delay between retries
	DefaultRetryMaxDelay = 10 * time.Second

	// DefaultRetryMultiplier is the default multiplier for exponential backoff
	DefaultRetryMultiplier = 2.0
)

// ============================================================================
// Circuit Breaker Configuration
// ============================================================================

const (
	// DefaultCircuitBreakerThreshold is the default failure threshold to open circuit
	DefaultCircuitBreakerThreshold = 5

	// DefaultCircuitBreakerTimeout is the default timeout before half-open state
	DefaultCircuitBreakerTimeout = 30 * time.Second

	// DefaultCircuitBreakerHalfOpenMax is the default max requests in half-open state
	DefaultCircuitBreakerHalfOpenMax = 3
)

// ============================================================================
// Limits
// ============================================================================

const (
	// MaxSubjectLength is the maximum length of a subject name
	MaxSubjectLength = 256

	// MaxHeaderValueLength is the maximum length of a header value
	MaxHeaderValueLength = 4096

	// MaxMetadataSize is the maximum size of metadata (16KB)
	MaxMetadataSize = 16 * 1024

	// MaxMessageIDLength is the maximum length of a message ID
	MaxMessageIDLength = 64

	// MaxTenantIDLength is the maximum length of a tenant ID
	MaxTenantIDLength = 128

	// MaxServiceNameLength is the maximum length of a service name
	MaxServiceNameLength = 256
)
