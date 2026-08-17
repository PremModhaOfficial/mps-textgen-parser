package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"
)

// EventsLoader loads events configuration from various sources
type EventsLoader struct {
	prefix string
}

// NewEventsLoader creates a new events config loader with optional env prefix
func NewEventsLoader(prefix string) *EventsLoader {
	return &EventsLoader{prefix: prefix}
}

// DefaultEventsLoader returns a loader with "EVENTS_" prefix
func DefaultEventsLoader() *EventsLoader {
	return NewEventsLoader("EVENTS_")
}

// LoadFromEnv loads events configuration from environment variables
func (loader *EventsLoader) LoadFromEnv() (*EventsConfig, error) {
	cfg := DefaultEventsConfig()

	// Servers
	if serversValue := loader.getEnv("SERVERS"); serversValue != "" {
		cfg.Servers = loader.parseServers(serversValue)
	}

	// Client name
	if nameValue := loader.getEnv("NAME"); nameValue != "" {
		cfg.Name = nameValue
	}

	// Timeouts
	if d := loader.getDuration("CONNECT_TIMEOUT"); d > 0 {
		cfg.ConnectTimeout = d
	}

	if d := loader.getDuration("DRAIN_TIMEOUT"); d > 0 {
		cfg.DrainTimeout = d
	}

	// Reconnection
	if v := loader.getInt("RECONNECT_MAX_ATTEMPTS"); v != 0 {
		cfg.Reconnect.MaxAttempts = v
	}

	if d := loader.getDuration("RECONNECT_INITIAL"); d > 0 {
		cfg.Reconnect.InitialInterval = d
	}

	if d := loader.getDuration("RECONNECT_MAX"); d > 0 {
		cfg.Reconnect.MaxInterval = d
	}

	if v := loader.getFloat("RECONNECT_MULTIPLIER"); v > 0 {
		cfg.Reconnect.Multiplier = v
	}

	if v := loader.getFloat("RECONNECT_JITTER"); v >= 0 {
		cfg.Reconnect.Jitter = v
	}

	// TLS
	if loader.getBool("TLS_ENABLED") {
		cfg.TLS = &TLSConfig{
			Enabled:    true,
			CertFile:   loader.getEnv("TLS_CERT_FILE"),
			KeyFile:    loader.getEnv("TLS_KEY_FILE"),
			CAFile:     loader.getEnv("TLS_CA_FILE"),
			SkipVerify: loader.getBool("TLS_SKIP_VERIFY"),
		}
	}

	// Stream/JetStream
	cfg.Stream.Enabled = !loader.getBool("STREAM_DISABLED")
	if v := loader.getEnv("STREAM_DOMAIN"); v != "" {
		cfg.Stream.Domain = v
	}
	if v := loader.getEnv("STREAM_PREFIX"); v != "" {
		cfg.Stream.Prefix = v
	}

	// Idle cleanup
	if d := loader.getDuration("IDLE_TIMEOUT"); d > 0 {
		cfg.IdleTimeout = d
	}

	if d := loader.getDuration("CLEANUP_INTERVAL"); d > 0 {
		cfg.CleanupInterval = d
	}

	if err := ValidateEventsConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// LoadFromEnvWithOptions loads from env and applies additional options
func (loader *EventsLoader) LoadFromEnvWithOptions(options ...EventsOption) (*EventsConfig, error) {
	cfg, err := loader.LoadFromEnv()
	if err != nil {
		return nil, err
	}

	ApplyEventsOptions(cfg, options...)

	if err := ValidateEventsConfig(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// getEnv gets an environment variable with prefix
func (loader *EventsLoader) getEnv(key string) string {
	return os.Getenv(loader.prefix + key)
}

// getBool parses a boolean from environment
func (loader *EventsLoader) getBool(key string) bool {
	v := strings.ToLower(strings.TrimSpace(loader.getEnv(key)))
	return v == "true" || v == "1" || v == "yes"
}

// getInt parses an integer from environment
func (loader *EventsLoader) getInt(key string) int {
	if v := loader.getEnv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return 0
}

// getFloat parses a float from environment
func (loader *EventsLoader) getFloat(key string) float64 {
	if v := loader.getEnv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
	}
	return 0
}

// getDuration parses a duration from environment
func (loader *EventsLoader) getDuration(key string) time.Duration {
	if v := loader.getEnv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return 0
}

// parseServers splits a comma-separated server list
func (loader *EventsLoader) parseServers(serversString string) []string {
	parts := strings.Split(serversString, ",")
	servers := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			servers = append(servers, trimmed)
		}
	}
	return servers
}

// LoadEventsFromEnv loads events config using the default loader
func LoadEventsFromEnv() (*EventsConfig, error) {
	return DefaultEventsLoader().LoadFromEnv()
}

// LoadEventsFromEnvWithPrefix loads events config using a custom prefix
func LoadEventsFromEnvWithPrefix(prefix string) (*EventsConfig, error) {
	return NewEventsLoader(prefix).LoadFromEnv()
}

// ValidateEventsConfig validates the events configuration
func ValidateEventsConfig(cfg *EventsConfig) error {
	if len(cfg.Servers) == 0 {
		return utils.ConfigError{Field: "Servers", Message: "at least one server is required"}
	}

	for i, server := range cfg.Servers {
		if server == "" {
			return utils.ConfigError{Field: "Servers", Message: "server URL cannot be empty"}
		}
		if !strings.HasPrefix(server, "nats://") && !strings.HasPrefix(server, "tls://") && !strings.HasPrefix(server, "ws://") && !strings.HasPrefix(server, "wss://") {
			return utils.ConfigError{Field: "Servers", Message: "server " + strconv.Itoa(i) + " must have valid scheme (nats://, tls://, ws://, or wss://)"}
		}
	}

	if cfg.ConnectTimeout <= 0 {
		return utils.ConfigError{Field: "ConnectTimeout", Message: "must be positive"}
	}

	if cfg.DrainTimeout <= 0 {
		return utils.ConfigError{Field: "DrainTimeout", Message: "must be positive"}
	}

	if cfg.Reconnect.InitialInterval <= 0 {
		return utils.ConfigError{Field: "Reconnect.InitialInterval", Message: "must be positive"}
	}

	if cfg.Reconnect.MaxInterval < cfg.Reconnect.InitialInterval {
		return utils.ConfigError{Field: "Reconnect.MaxInterval", Message: "must be >= InitialInterval"}
	}

	if cfg.Reconnect.Multiplier < 1.0 {
		return utils.ConfigError{Field: "Reconnect.Multiplier", Message: "must be >= 1.0"}
	}

	if cfg.Reconnect.Jitter < 0 || cfg.Reconnect.Jitter > 1 {
		return utils.ConfigError{Field: "Reconnect.Jitter", Message: "must be between 0.0 and 1.0"}
	}

	if cfg.IdleTimeout <= 0 {
		return utils.ConfigError{Field: "IdleTimeout", Message: "must be positive"}
	}

	if cfg.CleanupInterval <= 0 {
		return utils.ConfigError{Field: "CleanupInterval", Message: "must be positive"}
	}

	if cfg.TLS != nil && cfg.TLS.Enabled {
		if err := ValidateEventsTLS(cfg.TLS); err != nil {
			return err
		}
	}

	return nil
}

// ValidateEventsTLS validates TLS configuration
func ValidateEventsTLS(tlsConfig *TLSConfig) error {
	if tlsConfig == nil {
		return nil
	}

	if (tlsConfig.CertFile != "" && tlsConfig.KeyFile == "") ||
		(tlsConfig.CertFile == "" && tlsConfig.KeyFile != "") {
		return utils.ConfigError{Field: "TLS", Message: "both CertFile and KeyFile must be provided together"}
	}

	if tlsConfig.CertFile != "" {
		if _, err := os.Stat(tlsConfig.CertFile); os.IsNotExist(err) {
			return utils.ConfigError{Field: "TLS.CertFile", Message: "certificate file does not exist: " + tlsConfig.CertFile}
		}
	}

	if tlsConfig.KeyFile != "" {
		if _, err := os.Stat(tlsConfig.KeyFile); os.IsNotExist(err) {
			return utils.ConfigError{Field: "TLS.KeyFile", Message: "key file does not exist: " + tlsConfig.KeyFile}
		}
	}

	if tlsConfig.CAFile != "" {
		if _, err := os.Stat(tlsConfig.CAFile); os.IsNotExist(err) {
			return utils.ConfigError{Field: "TLS.CAFile", Message: "CA file does not exist: " + tlsConfig.CAFile}
		}
	}

	return nil
}

// ValidateEventsPublish validates publish configuration
func ValidateEventsPublish(cfg *PublishConfig) error {
	if cfg.Retry.MaxAttempts < 0 {
		return utils.ConfigError{Field: "Retry.MaxAttempts", Message: "cannot be negative"}
	}

	if cfg.AckTimeout <= 0 {
		return utils.ConfigError{Field: "AckTimeout", Message: "must be positive"}
	}

	return nil
}

// ValidateEventsSubscribe validates subscribe configuration
func ValidateEventsSubscribe(cfg *SubscribeConfig) error {
	if cfg.MaxConcurrent <= 0 {
		return utils.ConfigError{Field: "MaxConcurrent", Message: "must be positive"}
	}

	if cfg.AckWait <= 0 {
		return utils.ConfigError{Field: "AckWait", Message: "must be positive"}
	}

	return nil
}
