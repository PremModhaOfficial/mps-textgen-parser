package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

// Loader loads configuration from various sources
type Loader struct {
	prefix string
}

// NewLoader creates a new config loader with optional env prefix
func NewLoader(prefix string) *Loader {
	return &Loader{prefix: prefix}
}

// DefaultLoader returns a loader with "EVENTS_" prefix
func DefaultLoader() *Loader {
	return NewLoader("EVENTS_")
}

// LoadFromEnv loads configuration from environment variables
func (loader *Loader) LoadFromEnv() (*Config, error) {
	eventsConfiguration := DefaultConfig()

	// Servers
	if serversValue := loader.getEnv("SERVERS"); serversValue != "" {
		eventsConfiguration.Servers = loader.parseServers(serversValue)
	}

	// Client name
	if nameValue := loader.getEnv("NAME"); nameValue != "" {
		eventsConfiguration.Name = nameValue
	}

	// Timeouts
	if connectTimeoutDuration := loader.getDuration("CONNECT_TIMEOUT"); connectTimeoutDuration > 0 {
		eventsConfiguration.ConnectTimeout = connectTimeoutDuration
	}

	if drainTimeoutDuration := loader.getDuration("DRAIN_TIMEOUT"); drainTimeoutDuration > 0 {
		eventsConfiguration.DrainTimeout = drainTimeoutDuration
	}

	// Reconnection
	if maxAttemptsValue := loader.getInt("RECONNECT_MAX_ATTEMPTS"); maxAttemptsValue != 0 {
		eventsConfiguration.Reconnect.MaxAttempts = maxAttemptsValue
	}

	if initialIntervalDuration := loader.getDuration("RECONNECT_INITIAL"); initialIntervalDuration > 0 {
		eventsConfiguration.Reconnect.InitialInterval = initialIntervalDuration
	}

	if maxIntervalDuration := loader.getDuration("RECONNECT_MAX"); maxIntervalDuration > 0 {
		eventsConfiguration.Reconnect.MaxInterval = maxIntervalDuration
	}

	if multiplierValue := loader.getFloat("RECONNECT_MULTIPLIER"); multiplierValue > 0 {
		eventsConfiguration.Reconnect.Multiplier = multiplierValue
	}

	if jitterValue := loader.getFloat("RECONNECT_JITTER"); jitterValue >= 0 {
		eventsConfiguration.Reconnect.Jitter = jitterValue
	}

	// TLS
	if loader.getBool("TLS_ENABLED") {
		eventsConfiguration.TLS = &TLSConfig{
			Enabled:    true,
			CertFile:   loader.getEnv("TLS_CERT_FILE"),
			KeyFile:    loader.getEnv("TLS_KEY_FILE"),
			CAFile:     loader.getEnv("TLS_CA_FILE"),
			SkipVerify: loader.getBool("TLS_SKIP_VERIFY"),
		}
	}

	// Stream/JetStream
	eventsConfiguration.Stream.Enabled = !loader.getBool("STREAM_DISABLED")
	if domainValue := loader.getEnv("STREAM_DOMAIN"); domainValue != "" {
		eventsConfiguration.Stream.Domain = domainValue
	}
	if prefixValue := loader.getEnv("STREAM_PREFIX"); prefixValue != "" {
		eventsConfiguration.Stream.Prefix = prefixValue
	}

	// Idle cleanup
	if idleTimeoutDuration := loader.getDuration("IDLE_TIMEOUT"); idleTimeoutDuration > 0 {
		eventsConfiguration.IdleTimeout = idleTimeoutDuration
	}

	if cleanupIntervalDuration := loader.getDuration("CLEANUP_INTERVAL"); cleanupIntervalDuration > 0 {
		eventsConfiguration.CleanupInterval = cleanupIntervalDuration
	}

	if validationError := Validate(eventsConfiguration); validationError != nil {
		return nil, validationError
	}

	return eventsConfiguration, nil
}

// LoadFromEnvWithOptions loads from env and applies additional options
func (loader *Loader) LoadFromEnvWithOptions(options ...Option) (*Config, error) {
	eventsConfiguration, loadError := loader.LoadFromEnv()
	if loadError != nil {
		return nil, loadError
	}

	Apply(eventsConfiguration, options...)

	if validationError := Validate(eventsConfiguration); validationError != nil {
		return nil, validationError
	}

	return eventsConfiguration, nil
}

// getEnv gets an environment variable with prefix
func (loader *Loader) getEnv(key string) string {
	return os.Getenv(loader.prefix + key)
}

// getBool parses a boolean from environment
func (loader *Loader) getBool(key string) bool {
	envValue := strings.ToLower(strings.TrimSpace(loader.getEnv(key)))
	return envValue == "true" || envValue == "1" || envValue == "yes"
}

// getInt parses an integer from environment
func (loader *Loader) getInt(key string) int {
	if envValue := loader.getEnv(key); envValue != "" {
		if parsedInteger, parseError := strconv.Atoi(envValue); parseError == nil {
			return parsedInteger
		}
	}
	return 0
}

// getFloat parses a float from environment
func (loader *Loader) getFloat(key string) float64 {
	if envValue := loader.getEnv(key); envValue != "" {
		if parsedFloat, parseError := strconv.ParseFloat(envValue, 64); parseError == nil {
			return parsedFloat
		}
	}
	return 0
}

// getDuration parses a duration from environment
func (loader *Loader) getDuration(key string) time.Duration {
	if envValue := loader.getEnv(key); envValue != "" {
		if parsedDuration, parseError := time.ParseDuration(envValue); parseError == nil {
			return parsedDuration
		}
	}
	return 0
}

// parseServers splits a comma-separated server list
func (loader *Loader) parseServers(serversString string) []string {
	serverParts := strings.Split(serversString, ",")
	serversList := make([]string, 0, len(serverParts))
	for _, serverPart := range serverParts {
		if trimmedServer := strings.TrimSpace(serverPart); trimmedServer != "" {
			serversList = append(serversList, trimmedServer)
		}
	}
	return serversList
}

// LoadFromEnv loads using the default loader
func LoadFromEnv() (*Config, error) {
	return DefaultLoader().LoadFromEnv()
}

// LoadFromEnvWithPrefix loads using a custom prefix
func LoadFromEnvWithPrefix(prefix string) (*Config, error) {
	return NewLoader(prefix).LoadFromEnv()
}

// Validate validates the configuration
func Validate(configuration *Config) error {
	if len(configuration.Servers) == 0 {
		return core.ConfigError{Field: "Servers", Message: "at least one server is required"}
	}

	// Validate server URLs
	for i, server := range configuration.Servers {
		if server == "" {
			return core.ConfigError{Field: "Servers", Message: "server URL cannot be empty"}
		}
		// Check for valid NATS scheme
		if !strings.HasPrefix(server, "nats://") && !strings.HasPrefix(server, "tls://") && !strings.HasPrefix(server, "ws://") && !strings.HasPrefix(server, "wss://") {
			return core.ConfigError{Field: "Servers", Message: "server " + strconv.Itoa(i) + " must have valid scheme (nats://, tls://, ws://, or wss://)"}
		}
	}

	if configuration.ConnectTimeout <= 0 {
		return core.ConfigError{Field: "ConnectTimeout", Message: "must be positive"}
	}

	if configuration.DrainTimeout <= 0 {
		return core.ConfigError{Field: "DrainTimeout", Message: "must be positive"}
	}

	if configuration.Reconnect.InitialInterval <= 0 {
		return core.ConfigError{Field: "Reconnect.InitialInterval", Message: "must be positive"}
	}

	if configuration.Reconnect.MaxInterval < configuration.Reconnect.InitialInterval {
		return core.ConfigError{Field: "Reconnect.MaxInterval", Message: "must be >= InitialInterval"}
	}

	if configuration.Reconnect.Multiplier < 1.0 {
		return core.ConfigError{Field: "Reconnect.Multiplier", Message: "must be >= 1.0"}
	}

	if configuration.Reconnect.Jitter < 0 || configuration.Reconnect.Jitter > 1 {
		return core.ConfigError{Field: "Reconnect.Jitter", Message: "must be between 0.0 and 1.0"}
	}

	if configuration.IdleTimeout <= 0 {
		return core.ConfigError{Field: "IdleTimeout", Message: "must be positive"}
	}

	if configuration.CleanupInterval <= 0 {
		return core.ConfigError{Field: "CleanupInterval", Message: "must be positive"}
	}

	// Validate TLS configuration
	if configuration.TLS != nil && configuration.TLS.Enabled {
		if err := ValidateTLS(configuration.TLS); err != nil {
			return err
		}
	}

	return nil
}

// ValidateTLS validates TLS configuration
func ValidateTLS(tlsConfig *TLSConfig) error {
	if tlsConfig == nil {
		return nil
	}

	// If cert file is provided, key file must also be provided
	if (tlsConfig.CertFile != "" && tlsConfig.KeyFile == "") ||
		(tlsConfig.CertFile == "" && tlsConfig.KeyFile != "") {
		return core.ConfigError{Field: "TLS", Message: "both CertFile and KeyFile must be provided together"}
	}

	// If cert file is provided, verify it exists
	if tlsConfig.CertFile != "" {
		if _, err := os.Stat(tlsConfig.CertFile); os.IsNotExist(err) {
			return core.ConfigError{Field: "TLS.CertFile", Message: "certificate file does not exist: " + tlsConfig.CertFile}
		}
	}

	// If key file is provided, verify it exists
	if tlsConfig.KeyFile != "" {
		if _, err := os.Stat(tlsConfig.KeyFile); os.IsNotExist(err) {
			return core.ConfigError{Field: "TLS.KeyFile", Message: "key file does not exist: " + tlsConfig.KeyFile}
		}
	}

	// If CA file is provided, verify it exists
	if tlsConfig.CAFile != "" {
		if _, err := os.Stat(tlsConfig.CAFile); os.IsNotExist(err) {
			return core.ConfigError{Field: "TLS.CAFile", Message: "CA file does not exist: " + tlsConfig.CAFile}
		}
	}

	return nil
}

// ValidatePublish validates publish configuration
func ValidatePublish(publishConfiguration *PublishConfig) error {
	if publishConfiguration.Retry.MaxAttempts < 0 {
		return core.ConfigError{Field: "Retry.MaxAttempts", Message: "cannot be negative"}
	}

	if publishConfiguration.AckTimeout <= 0 {
		return core.ConfigError{Field: "AckTimeout", Message: "must be positive"}
	}

	return nil
}

// ValidateSubscribe validates subscribe configuration
func ValidateSubscribe(subscribeConfiguration *SubscribeConfig) error {
	if subscribeConfiguration.MaxConcurrent <= 0 {
		return core.ConfigError{Field: "MaxConcurrent", Message: "must be positive"}
	}

	if subscribeConfiguration.AckWait <= 0 {
		return core.ConfigError{Field: "AckWait", Message: "must be positive"}
	}

	return nil
}
