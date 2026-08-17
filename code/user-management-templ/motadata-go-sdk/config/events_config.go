package config

import (
	"time"
)

// EventsConfig holds the base configuration for the events system
type EventsConfig struct {
	// Connection settings
	Servers        []string      // Server URLs
	Name           string        // Client name for identification
	ConnectTimeout time.Duration // Connection timeout
	DrainTimeout   time.Duration // Drain timeout for graceful shutdown

	// TLS configuration
	TLS *TLSConfig

	// Reconnection settings
	Reconnect ReconnectConfig

	// Stream/JetStream settings
	Stream StreamConfig

	// Idle connection cleanup
	IdleTimeout     time.Duration // Close connections idle longer than this
	CleanupInterval time.Duration // How often to check for idle connections
}

// TLSConfig holds TLS configuration
type TLSConfig struct {
	Enabled    bool
	CertFile   string
	KeyFile    string
	CAFile     string
	SkipVerify bool
}

// ReconnectConfig holds reconnection settings with exponential backoff
type ReconnectConfig struct {
	MaxAttempts     int           // -1 for infinite retries
	InitialInterval time.Duration // Starting backoff interval
	MaxInterval     time.Duration // Maximum backoff interval
	Multiplier      float64       // Backoff multiplier
	Jitter          float64       // Random jitter factor (0.0-1.0)
}

// StreamConfig holds stream/JetStream configuration
type StreamConfig struct {
	Enabled bool
	Domain  string // JetStream domain (for leaf nodes)
	Prefix  string // API prefix
}

// PublishConfig holds publishing configuration
type PublishConfig struct {
	// Retry settings
	Retry RetryConfig

	// Timeouts
	AckTimeout time.Duration // How long to wait for acknowledgment

	// Deduplication
	EnableDeduplication bool
	DeduplicationWindow time.Duration
}

// RetryConfig holds retry settings
type RetryConfig struct {
	MaxAttempts     int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	Jitter          float64
}

// SubscribeConfig holds subscription configuration
type SubscribeConfig struct {
	// Queue group for load balancing
	QueueGroup string

	// Processing settings
	MaxConcurrent int           // Max concurrent message handlers
	AckWait       time.Duration // How long before message is redelivered

	// Batch settings
	BatchSize int           // Messages per batch
	BatchWait time.Duration // Max wait time for batch
}

// KVConfig holds KV store configuration
type KVConfig struct {
	BucketName   string
	TTL          time.Duration
	History      int
	MaxValueSize int32
	Storage      string // "file" or "memory"
	Replicas     int
	Description  string
}

// ObjectStoreConfig holds object store configuration
type ObjectStoreConfig struct {
	BucketName    string
	MaxChunkSize  int32
	MaxObjectSize int64
	Storage       string // "file" or "memory"
	Replicas      int
	Description   string
}

// DefaultEventsConfig returns sensible defaults for the events system
func DefaultEventsConfig() *EventsConfig {
	return &EventsConfig{
		Servers:         []string{"nats://localhost:4222"},
		ConnectTimeout:  10 * time.Second,
		DrainTimeout:    30 * time.Second,
		Reconnect:       DefaultReconnectConfig(),
		Stream:          StreamConfig{Enabled: true},
		IdleTimeout:     30 * time.Minute,
		CleanupInterval: 1 * time.Minute,
	}
}

// DefaultReconnectConfig returns default reconnection settings
func DefaultReconnectConfig() ReconnectConfig {
	return ReconnectConfig{
		MaxAttempts:     -1, // Infinite
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     30 * time.Second,
		Multiplier:      2.0,
		Jitter:          0.1,
	}
}

// DefaultRetryConfig returns default retry settings
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:     3,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     5 * time.Second,
		Multiplier:      2.0,
		Jitter:          0.1,
	}
}

// DefaultPublishConfig returns default publish settings
func DefaultPublishConfig() PublishConfig {
	return PublishConfig{
		Retry:               DefaultRetryConfig(),
		AckTimeout:          5 * time.Second,
		EnableDeduplication: true,
		DeduplicationWindow: 2 * time.Minute,
	}
}

// DefaultSubscribeConfig returns default subscription settings
func DefaultSubscribeConfig() SubscribeConfig {
	return SubscribeConfig{
		MaxConcurrent: 10,
		AckWait:       30 * time.Second,
		BatchSize:     100,
		BatchWait:     100 * time.Millisecond,
	}
}

// Clone creates a deep copy of the events config
func (c *EventsConfig) Clone() *EventsConfig {
	if c == nil {
		return nil
	}

	cloned := *c
	cloned.Servers = append([]string{}, c.Servers...)

	if c.TLS != nil {
		tlsCopy := *c.TLS
		cloned.TLS = &tlsCopy
	}

	return &cloned
}
