package config

import "time"

// EventsOption is a functional option for configuring the events system
type EventsOption func(*EventsConfig)

// WithServers sets the server URLs
func WithServers(servers ...string) EventsOption {
	return func(c *EventsConfig) {
		c.Servers = servers
	}
}

// WithName sets the client name
func WithName(name string) EventsOption {
	return func(c *EventsConfig) {
		c.Name = name
	}
}

// WithConnectTimeout sets the connection timeout
func WithConnectTimeout(timeout time.Duration) EventsOption {
	return func(c *EventsConfig) {
		c.ConnectTimeout = timeout
	}
}

// WithDrainTimeout sets the drain timeout
func WithDrainTimeout(timeout time.Duration) EventsOption {
	return func(c *EventsConfig) {
		c.DrainTimeout = timeout
	}
}

// WithTLS sets TLS configuration
func WithTLS(tlsConfig *TLSConfig) EventsOption {
	return func(c *EventsConfig) {
		c.TLS = tlsConfig
	}
}

// WithTLSFiles sets TLS using certificate files
func WithTLSFiles(certFile, keyFile, caFile string) EventsOption {
	return func(c *EventsConfig) {
		c.TLS = &TLSConfig{
			Enabled:  true,
			CertFile: certFile,
			KeyFile:  keyFile,
			CAFile:   caFile,
		}
	}
}

// WithTLSSkipVerify enables TLS with certificate verification disabled
func WithTLSSkipVerify() EventsOption {
	return func(c *EventsConfig) {
		if c.TLS == nil {
			c.TLS = &TLSConfig{Enabled: true}
		}
		c.TLS.SkipVerify = true
	}
}

// WithReconnect sets reconnection configuration
func WithReconnect(reconnectConfig ReconnectConfig) EventsOption {
	return func(c *EventsConfig) {
		c.Reconnect = reconnectConfig
	}
}

// WithMaxReconnects sets the maximum reconnection attempts
func WithMaxReconnects(maxAttempts int) EventsOption {
	return func(c *EventsConfig) {
		c.Reconnect.MaxAttempts = maxAttempts
	}
}

// WithStream sets stream configuration
func WithStream(streamConfig StreamConfig) EventsOption {
	return func(c *EventsConfig) {
		c.Stream = streamConfig
	}
}

// WithStreamDisabled disables stream/JetStream
func WithStreamDisabled() EventsOption {
	return func(c *EventsConfig) {
		c.Stream.Enabled = false
	}
}

// WithStreamDomain sets the JetStream domain
func WithStreamDomain(domain string) EventsOption {
	return func(c *EventsConfig) {
		c.Stream.Domain = domain
	}
}

// WithIdleTimeout sets the idle connection timeout
func WithIdleTimeout(timeout time.Duration) EventsOption {
	return func(c *EventsConfig) {
		c.IdleTimeout = timeout
	}
}

// WithCleanupInterval sets the cleanup interval
func WithCleanupInterval(interval time.Duration) EventsOption {
	return func(c *EventsConfig) {
		c.CleanupInterval = interval
	}
}

// ApplyEventsOptions applies options to an events config
func ApplyEventsOptions(c *EventsConfig, options ...EventsOption) *EventsConfig {
	for _, option := range options {
		option(c)
	}
	return c
}

// PublishOption is a functional option for publish configuration
type PublishOption func(*PublishConfig)

// WithRetry sets retry configuration for publishing
func WithRetry(retryConfig RetryConfig) PublishOption {
	return func(c *PublishConfig) {
		c.Retry = retryConfig
	}
}

// WithMaxRetries sets max retry attempts
func WithMaxRetries(maxAttempts int) PublishOption {
	return func(c *PublishConfig) {
		c.Retry.MaxAttempts = maxAttempts
	}
}

// WithAckTimeout sets the acknowledgment timeout
func WithAckTimeout(timeout time.Duration) PublishOption {
	return func(c *PublishConfig) {
		c.AckTimeout = timeout
	}
}

// WithDeduplication enables/disables deduplication
func WithDeduplication(enabled bool) PublishOption {
	return func(c *PublishConfig) {
		c.EnableDeduplication = enabled
	}
}

// WithDeduplicationWindow sets the deduplication window
func WithDeduplicationWindow(windowDuration time.Duration) PublishOption {
	return func(c *PublishConfig) {
		c.DeduplicationWindow = windowDuration
	}
}

// SubscribeOption is a functional option for subscription configuration
type SubscribeOption func(*SubscribeConfig)

// WithQueueGroup sets the queue group for subscriptions
func WithQueueGroup(queueGroup string) SubscribeOption {
	return func(c *SubscribeConfig) {
		c.QueueGroup = queueGroup
	}
}

// WithMaxConcurrent sets max concurrent handlers
func WithMaxConcurrent(maxConcurrentHandlers int) SubscribeOption {
	return func(c *SubscribeConfig) {
		c.MaxConcurrent = maxConcurrentHandlers
	}
}

// WithAckWait sets how long to wait before redelivery
func WithAckWait(waitDuration time.Duration) SubscribeOption {
	return func(c *SubscribeConfig) {
		c.AckWait = waitDuration
	}
}

// WithBatch sets batch processing settings
func WithBatch(batchSize int, batchWaitDuration time.Duration) SubscribeOption {
	return func(c *SubscribeConfig) {
		c.BatchSize = batchSize
		c.BatchWait = batchWaitDuration
	}
}
