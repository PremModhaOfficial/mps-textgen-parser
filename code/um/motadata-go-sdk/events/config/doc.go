// Package config provides configuration management for the events system.
//
// This package handles all configuration aspects including:
//   - Connection settings (servers, timeouts)
//   - TLS/SSL configuration
//   - Reconnection with exponential backoff
//   - JetStream/Stream settings
//   - Publishing and subscription options
//   - Environment variable loading
//
// # Basic Usage
//
// Create a configuration with defaults:
//
//	cfg := config.DefaultConfig()
//	cfg.Servers = []string{"nats://localhost:4222"}
//
// # Functional Options
//
// Use functional options for fluent configuration:
//
//	cfg := config.DefaultConfig()
//	config.Apply(cfg,
//	    config.WithServers("nats://server1:4222", "nats://server2:4222"),
//	    config.WithName("my-service"),
//	    config.WithConnectTimeout(5 * time.Second),
//	    config.WithMaxReconnects(-1), // Infinite reconnects
//	)
//
// # Environment Variables
//
// Load configuration from environment variables with the default EVENTS_ prefix:
//
//	cfg, err := config.LoadFromEnv()
//
// Or with a custom prefix:
//
//	cfg, err := config.LoadFromEnvWithPrefix("MYAPP_NATS_")
//
// Supported environment variables:
//
//	EVENTS_SERVERS              - Comma-separated server URLs
//	EVENTS_NAME                 - Client name
//	EVENTS_CONNECT_TIMEOUT      - Connection timeout (e.g., "10s")
//	EVENTS_DRAIN_TIMEOUT        - Drain timeout (e.g., "30s")
//	EVENTS_RECONNECT_MAX_ATTEMPTS - Max reconnection attempts (-1 for infinite)
//	EVENTS_RECONNECT_INITIAL    - Initial reconnect interval
//	EVENTS_RECONNECT_MAX        - Maximum reconnect interval
//	EVENTS_RECONNECT_MULTIPLIER - Backoff multiplier
//	EVENTS_RECONNECT_JITTER     - Random jitter (0.0-1.0)
//	EVENTS_TLS_ENABLED          - Enable TLS (true/false)
//	EVENTS_TLS_CERT_FILE        - Path to TLS certificate
//	EVENTS_TLS_KEY_FILE         - Path to TLS key
//	EVENTS_TLS_CA_FILE          - Path to CA certificate
//	EVENTS_TLS_SKIP_VERIFY      - Skip certificate verification
//	EVENTS_STREAM_DISABLED      - Disable JetStream
//	EVENTS_STREAM_DOMAIN        - JetStream domain
//	EVENTS_IDLE_TIMEOUT         - Idle connection timeout
//	EVENTS_CLEANUP_INTERVAL     - Cleanup check interval
//
// # TLS Configuration
//
// Configure TLS for secure connections:
//
//	cfg := config.DefaultConfig()
//	config.Apply(cfg,
//	    config.WithTLSFiles(
//	        "/path/to/client-cert.pem",
//	        "/path/to/client-key.pem",
//	        "/path/to/ca-cert.pem",
//	    ),
//	)
//
// # Reconnection Settings
//
// Configure exponential backoff for reconnection:
//
//	cfg := config.DefaultConfig()
//	cfg.Reconnect = config.ReconnectConfig{
//	    MaxAttempts:     -1,                      // Infinite retries
//	    InitialInterval: 100 * time.Millisecond,
//	    MaxInterval:     30 * time.Second,
//	    Multiplier:      2.0,                     // Double each retry
//	    Jitter:          0.1,                     // 10% random jitter
//	}
//
// # Publish Configuration
//
// Configure publishing behavior:
//
//	pubCfg := config.DefaultPublishConfig()
//	pubCfg.Retry.MaxAttempts = 5
//	pubCfg.AckTimeout = 10 * time.Second
//	pubCfg.EnableDeduplication = true
//
// # Subscribe Configuration
//
// Configure subscription behavior:
//
//	subCfg := config.DefaultSubscribeConfig()
//	subCfg.QueueGroup = "my-service-workers"
//	subCfg.MaxConcurrent = 20
//	subCfg.AckWait = 60 * time.Second
//
// # Validation
//
// Validate configuration before use:
//
//	if err := config.Validate(cfg); err != nil {
//	    log.Fatalf("Invalid configuration: %v", err)
//	}
//
// # Configuration Types
//
// The package provides the following configuration types:
//   - Config: Main connection configuration
//   - TLSConfig: TLS/SSL settings
//   - ReconnectConfig: Reconnection with exponential backoff
//   - StreamConfig: JetStream settings
//   - PublishConfig: Publishing settings
//   - SubscribeConfig: Subscription settings
//   - RetryConfig: Retry settings for operations
package config
