package config

import "time"

// Option is a functional option for configuring the events system
type Option func(*Config)

// WithServers sets the server URLs
func WithServers(servers ...string) Option {
	return func(configuration *Config) {
		configuration.Servers = servers
	}
}

// WithName sets the client name
func WithName(name string) Option {
	return func(configuration *Config) {
		configuration.Name = name
	}
}

// WithConnectTimeout sets the connection timeout
func WithConnectTimeout(timeout time.Duration) Option {
	return func(configuration *Config) {
		configuration.ConnectTimeout = timeout
	}
}

// WithDrainTimeout sets the drain timeout
func WithDrainTimeout(timeout time.Duration) Option {
	return func(configuration *Config) {
		configuration.DrainTimeout = timeout
	}
}

// WithTLS sets TLS configuration
func WithTLS(tlsConfiguration *TLSConfig) Option {
	return func(configuration *Config) {
		configuration.TLS = tlsConfiguration
	}
}

// WithTLSFiles sets TLS using certificate files
func WithTLSFiles(certFile, keyFile, caFile string) Option {
	return func(configuration *Config) {
		configuration.TLS = &TLSConfig{
			Enabled:  true,
			CertFile: certFile,
			KeyFile:  keyFile,
			CAFile:   caFile,
		}
	}
}

// WithTLSSkipVerify enables TLS with certificate verification disabled
func WithTLSSkipVerify() Option {
	return func(configuration *Config) {
		if configuration.TLS == nil {
			configuration.TLS = &TLSConfig{Enabled: true}
		}
		configuration.TLS.SkipVerify = true
	}
}

// WithReconnect sets reconnection configuration
func WithReconnect(reconnectConfiguration ReconnectConfig) Option {
	return func(configuration *Config) {
		configuration.Reconnect = reconnectConfiguration
	}
}

// WithMaxReconnects sets the maximum reconnection attempts
func WithMaxReconnects(maxAttempts int) Option {
	return func(configuration *Config) {
		configuration.Reconnect.MaxAttempts = maxAttempts
	}
}

// WithStream sets stream configuration
func WithStream(streamConfiguration StreamConfig) Option {
	return func(configuration *Config) {
		configuration.Stream = streamConfiguration
	}
}

// WithStreamDisabled disables stream/JetStream
func WithStreamDisabled() Option {
	return func(configuration *Config) {
		configuration.Stream.Enabled = false
	}
}

// WithStreamDomain sets the JetStream domain
func WithStreamDomain(domain string) Option {
	return func(configuration *Config) {
		configuration.Stream.Domain = domain
	}
}

// WithIdleTimeout sets the idle connection timeout
func WithIdleTimeout(timeout time.Duration) Option {
	return func(configuration *Config) {
		configuration.IdleTimeout = timeout
	}
}

// WithCleanupInterval sets the cleanup interval
func WithCleanupInterval(interval time.Duration) Option {
	return func(configuration *Config) {
		configuration.CleanupInterval = interval
	}
}

// Apply applies options to a config
func Apply(configuration *Config, options ...Option) *Config {
	for _, option := range options {
		option(configuration)
	}
	return configuration
}

// PublishOption is a functional option for publish configuration
type PublishOption func(*PublishConfig)

// WithRetry sets retry configuration for publishing
func WithRetry(retryConfiguration RetryConfig) PublishOption {
	return func(publishConfiguration *PublishConfig) {
		publishConfiguration.Retry = retryConfiguration
	}
}

// WithMaxRetries sets max retry attempts
func WithMaxRetries(maxAttempts int) PublishOption {
	return func(publishConfiguration *PublishConfig) {
		publishConfiguration.Retry.MaxAttempts = maxAttempts
	}
}

// WithAckTimeout sets the acknowledgment timeout
func WithAckTimeout(timeout time.Duration) PublishOption {
	return func(publishConfiguration *PublishConfig) {
		publishConfiguration.AckTimeout = timeout
	}
}

// WithDeduplication enables/disables deduplication
func WithDeduplication(enabled bool) PublishOption {
	return func(publishConfiguration *PublishConfig) {
		publishConfiguration.EnableDeduplication = enabled
	}
}

// WithDeduplicationWindow sets the deduplication window
func WithDeduplicationWindow(windowDuration time.Duration) PublishOption {
	return func(publishConfiguration *PublishConfig) {
		publishConfiguration.DeduplicationWindow = windowDuration
	}
}

// SubscribeOption is a functional option for subscription configuration
type SubscribeOption func(*SubscribeConfig)

// WithQueueGroup sets the queue group for subscriptions
func WithQueueGroup(queueGroup string) SubscribeOption {
	return func(subscribeConfiguration *SubscribeConfig) {
		subscribeConfiguration.QueueGroup = queueGroup
	}
}

// WithMaxConcurrent sets max concurrent handlers
func WithMaxConcurrent(maxConcurrentHandlers int) SubscribeOption {
	return func(subscribeConfiguration *SubscribeConfig) {
		subscribeConfiguration.MaxConcurrent = maxConcurrentHandlers
	}
}

// WithAckWait sets how long to wait before redelivery
func WithAckWait(waitDuration time.Duration) SubscribeOption {
	return func(subscribeConfiguration *SubscribeConfig) {
		subscribeConfiguration.AckWait = waitDuration
	}
}

// WithBatch sets batch processing settings
func WithBatch(batchSize int, batchWaitDuration time.Duration) SubscribeOption {
	return func(subscribeConfiguration *SubscribeConfig) {
		subscribeConfiguration.BatchSize = batchSize
		subscribeConfiguration.BatchWait = batchWaitDuration
	}
}
