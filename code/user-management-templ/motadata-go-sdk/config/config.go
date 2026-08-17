package config

import "time"

// Config is the root configuration structure containing all module configs.
type Config struct {
	// Service identification (shared across all modules)
	Service ServiceConfig `yaml:"service"`

	// Logger configuration
	Logger LoggerConfig `yaml:"logger"`

	// Metrics configuration
	Metrics MetricsConfig `yaml:"metrics"`

	// Tracer configuration
	Tracer TracerConfig `yaml:"tracer"`
}

// ServiceConfig holds common service identification.
type ServiceConfig struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Environment string `yaml:"environment"`
}

// LoggerConfig holds logging configuration.
// This flat structure is used directly by the logger package.
type LoggerConfig struct {
	// Level is minimum log level (debug, info, warn, error)
	Level string

	// Service identification
	ServiceName    string
	ServiceVersion string
	Environment    string

	// Console output
	ConsoleEnabled bool
	ConsoleFormat  string // "json" or "console"

	// OTEL output
	OTELEnabled  bool
	OTELEndpoint string // e.g., "localhost:4317"
	OTELInsecure bool
	OTELProtocol string // "grpc" or "http"
	OTELDebug    bool   // Enable stdout exporter for debugging OTEL pipeline

	// File output
	FileEnabled    bool
	FilePath       string
	FileMaxSizeMB  int
	FileMaxBackups int
	FileMaxAgeDays int
	FileCompress   bool

	// Caller info
	AddCaller  bool
	CallerSkip int

	// Module-specific log levels
	ModuleLevels map[string]string

	// Batch processor settings for OTEL export
	BatchTimeout    time.Duration // Timeout for batch export (default: 30s)
	ExportInterval  time.Duration // Interval between exports (default: 1s)
	MaxExportBatch  int           // Max logs per batch (default: 512)
	MaxQueueSize    int           // Max logs in queue (default: 2048)
	ShutdownTimeout time.Duration // Timeout for graceful shutdown (default: 10s)
}

// MetricsConfig holds metrics configuration.
// This flat structure is used directly by the metrics package.
type MetricsConfig struct {
	// Enabled controls whether metrics collection is active
	Enabled bool

	// Service identification
	ServiceName    string
	ServiceVersion string
	Environment    string

	// OTEL exporter settings
	OTELEnabled  bool
	OTELEndpoint string // e.g., "localhost:4317"
	OTELInsecure bool
	OTELProtocol string // "grpc" or "http"

	// Prometheus exporter settings
	PrometheusEnabled bool
	PrometheusPort    int    // Port for /metrics endpoint
	PrometheusPath    string // Path for metrics endpoint

	// Export settings
	ExportInterval time.Duration

	// Default labels applied to all metrics
	DefaultLabels map[string]string
}

// TracerConfig holds tracing configuration.
// This flat structure is used directly by the tracer package.
type TracerConfig struct {
	// Enabled controls whether tracing is active
	Enabled bool

	// Service identification
	ServiceName    string
	ServiceVersion string
	Environment    string

	// OTEL exporter settings
	OTELEndpoint string // e.g., "localhost:4317"
	OTELInsecure bool
	OTELProtocol string // "grpc" or "http"
	OTELDebug    bool   // Enable stdout exporter for debugging

	// Sampling configuration
	SamplingRatio float64 // 0.0 to 1.0 (1.0 = sample everything)

	// Propagation
	Propagators []string // e.g., ["tracecontext", "baggage"]

	// Batch settings
	BatchTimeout   time.Duration
	MaxExportBatch int
	MaxQueueSize   int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Service: ServiceConfig{
			Name:        "app",
			Version:     "0.0.0",
			Environment: "development",
		},
		Logger:  DefaultLoggerConfig(),
		Metrics: DefaultMetricsConfig(),
		Tracer:  DefaultTracerConfig(),
	}
}

// DefaultLoggerConfig returns sensible defaults for logger configuration.
// Note: ServiceName, ServiceVersion, Environment are left empty so they
// can be populated from the Service config when using the unified Config.
func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Level:           "info",
		ServiceName:     "",
		ServiceVersion:  "",
		Environment:     "",
		ConsoleEnabled:  true,
		ConsoleFormat:   "json",
		OTELEnabled:     false,
		OTELEndpoint:    "localhost:4317",
		OTELInsecure:    true,
		OTELProtocol:    "grpc",
		OTELDebug:       false,
		FileEnabled:     false,
		FilePath:        "logs/app.log",
		FileMaxSizeMB:   100,
		FileMaxBackups:  5,
		FileMaxAgeDays:  7,
		FileCompress:    true,
		AddCaller:       true,
		CallerSkip:      2,
		BatchTimeout:    30 * time.Second,
		ExportInterval:  time.Second,
		MaxExportBatch:  512,
		MaxQueueSize:    2048,
		ShutdownTimeout: 10 * time.Second,
	}
}

// DefaultMetricsConfig returns sensible defaults for metrics configuration.
// Note: ServiceName, ServiceVersion, Environment are left empty so they
// can be populated from the Service config when using the unified Config.
func DefaultMetricsConfig() MetricsConfig {
	return MetricsConfig{
		Enabled:           true,
		ServiceName:       "",
		ServiceVersion:    "",
		Environment:       "",
		OTELEnabled:       false,
		OTELEndpoint:      "localhost:4317",
		OTELInsecure:      true,
		OTELProtocol:      "grpc",
		PrometheusEnabled: false,
		PrometheusPort:    9090,
		PrometheusPath:    "/metrics",
		ExportInterval:    15 * time.Second,
		DefaultLabels:     make(map[string]string),
	}
}

// DefaultTracerConfig returns sensible defaults for tracer configuration.
// Note: ServiceName, ServiceVersion, Environment are left empty so they
// can be populated from the Service config when using the unified Config.
func DefaultTracerConfig() TracerConfig {
	return TracerConfig{
		Enabled:        false,
		ServiceName:    "",
		ServiceVersion: "",
		Environment:    "",
		OTELEndpoint:   "localhost:4317",
		OTELInsecure:   true,
		OTELProtocol:   "grpc",
		OTELDebug:      false,
		SamplingRatio:  1.0,
		Propagators:    []string{"tracecontext", "baggage"},
		BatchTimeout:   5 * time.Second,
		MaxExportBatch: 512,
		MaxQueueSize:   2048,
	}
}

// GetLoggerConfig returns the LoggerConfig with service fields populated.
func (c *Config) GetLoggerConfig() LoggerConfig {
	cfg := c.Logger
	// Override with service config if not set
	if cfg.ServiceName == "" {
		cfg.ServiceName = c.Service.Name
	}
	if cfg.ServiceVersion == "" {
		cfg.ServiceVersion = c.Service.Version
	}
	if cfg.Environment == "" {
		cfg.Environment = c.Service.Environment
	}
	return cfg
}

// GetMetricsConfig returns the MetricsConfig with service fields populated.
func (c *Config) GetMetricsConfig() MetricsConfig {
	cfg := c.Metrics
	// Override with service config if not set
	if cfg.ServiceName == "" {
		cfg.ServiceName = c.Service.Name
	}
	if cfg.ServiceVersion == "" {
		cfg.ServiceVersion = c.Service.Version
	}
	if cfg.Environment == "" {
		cfg.Environment = c.Service.Environment
	}
	return cfg
}

// GetTracerConfig returns the TracerConfig with service fields populated.
func (c *Config) GetTracerConfig() TracerConfig {
	tracerConfiguration := c.Tracer
	// Override with service config if not set
	if tracerConfiguration.ServiceName == "" {
		tracerConfiguration.ServiceName = c.Service.Name
	}
	if tracerConfiguration.ServiceVersion == "" {
		tracerConfiguration.ServiceVersion = c.Service.Version
	}
	if tracerConfiguration.Environment == "" {
		tracerConfiguration.Environment = c.Service.Environment
	}
	return tracerConfiguration
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	return nil
}
