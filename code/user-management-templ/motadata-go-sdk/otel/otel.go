// Package otel provides unified OpenTelemetry initialization for logging, metrics, and tracing.
// It serves as the main entry point for configuring all OTEL components in the motadata-go-sdk.
//
// The package provides:
//   - Unified initialization from config, file, or environment
//   - Coordinated lifecycle management for all OTEL components
//   - Graceful shutdown with proper ordering (tracer -> metrics -> logger)
//   - Must* variants that panic on error for use in main()
//
// Architecture:
//
//	                ┌─────────────┐
//	                │    OTEL     │
//	                │  (facade)   │
//	                └──────┬──────┘
//	       ┌───────────────┼───────────────┐
//	       ▼               ▼               ▼
//	┌─────────────┐ ┌─────────────┐ ┌─────────────┐
//	│   Logger    │ │   Tracer    │ │   Metrics   │
//	│ (otelzap)   │ │ (OTLP SDK)  │ │ (OTLP SDK)  │
//	└─────────────┘ └─────────────┘ └─────────────┘
//	       │               │               │
//	       └───────────────┴───────────────┘
//	                       │
//	                ┌──────▼──────┐
//	                │    OTLP     │
//	                │  Collector  │
//	                └─────────────┘
//
// Usage:
//
//	// From config struct
//	otelInstance, err := otel.Init(cfg)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer otelInstance.Shutdown(context.Background())
//
//	// From YAML file
//	otelInstance, err := otel.InitFromFile("config.yaml")
//
//	// From environment variables
//	otelInstance, err := otel.InitFromEnv()
//
//	// Must variant (panics on error)
//	otelInstance := otel.MustInit(cfg)
//	defer otelInstance.Shutdown(context.Background())
package otel

import (
	"context"
	"fmt"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/config"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/common"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/metrics"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer"
)

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// OTEL holds references to all initialized OTEL components.
// It manages the lifecycle of logger, tracer, and metrics subsystems.
//
// The OTEL struct provides:
//   - Access to the Logger for structured logging
//   - Access to the Tracer for distributed tracing
//   - Coordinated shutdown of all components
//
// Note: Metrics are managed via package-level functions (metrics.R())
// rather than being exposed directly on this struct.
type OTEL struct {
	// Logger provides structured logging with OTEL integration.
	// Use this for application logging with trace correlation.
	Logger *logger.Logger

	// Tracer provides distributed tracing capabilities.
	// Use this to create spans and propagate trace context.
	Tracer *tracer.Tracer

	// config holds the original configuration for reference
	config config.Config

	// metrics tracks whether metrics subsystem was initialized
	// Used during shutdown to determine if metrics cleanup is needed
	metrics bool
}

/* ---------------------------------------- Constructor Functions -------------------------------------------------- */

// Init initializes all OTEL components (logger, metrics, tracer) from the unified config.
// This is the primary initialization function that sets up the complete OTEL pipeline.
//
// Initialization order:
//  1. Logger - initialized first so it's available for error logging
//  2. Metrics - initialized second (independent of logger)
//  3. Tracer - initialized last (may use logger for errors)
//
// On partial failure, previously initialized components are cleaned up.
//
// Parameters:
//   - cfg: Unified configuration containing logger, metrics, and tracer settings
//
// Returns:
//   - *OTEL: The initialized OTEL instance with all components
//   - error: Non-nil if any component fails to initialize
//
// Usage:
//
//	cfg := config.Default()
//	cfg.Logger.Level = "debug"
//	cfg.Tracer.Enabled = true
//
//	otelInstance, err := otel.Init(cfg)
//	if err != nil {
//	    log.Fatal("OTEL init failed:", err)
//	}
//	defer otelInstance.Shutdown(context.Background())
func Init(cfg config.Config) (*OTEL, error) {

	otelInstance := &OTEL{
		config: cfg,
	}

	// Initialize logger first - it's used for all other error reporting
	loggerConfig := cfg.GetLoggerConfig()
	loggerInstance, loggererror := logger.Init(loggerConfig)

	if loggererror != nil {

		return nil, fmt.Errorf("initialize logger: %w", loggererror)
	}

	otelInstance.Logger = loggerInstance

	// Initialize metrics if enabled
	metricsConfig := cfg.GetMetricsConfig()

	if metricsConfig.Enabled {

		metricserror := metrics.Init(metricsConfig)

		if metricserror != nil {
			// Cleanup logger on metrics initialization failure
			_ = logger.Close()

			return nil, fmt.Errorf("initialize metrics: %w", metricserror)
		}

		otelInstance.metrics = true
	}

	// Initialize tracer if enabled
	tracerConfig := cfg.GetTracerConfig()

	if tracerConfig.Enabled {

		tracerInstance, tracererror := tracer.Init(tracerConfig)

		if tracererror != nil {
			// Cleanup previously initialized components on failure
			_ = logger.Close()

			if otelInstance.metrics {
				_ = metrics.Shutdown(nil)
			}

			return nil, fmt.Errorf("initialize tracer: %w", tracererror)
		}

		otelInstance.Tracer = tracerInstance
	}

	return otelInstance, nil
}

// InitFromFile initializes all OTEL components from a YAML config file.
// Loads the configuration from the specified path and delegates to Init().
//
// Parameters:
//   - configPath: Path to the YAML configuration file
//
// Returns:
//   - *OTEL: The initialized OTEL instance
//   - error: Non-nil if file loading or initialization fails
//
// Usage:
//
//	otelInstance, err := otel.InitFromFile("/etc/myapp/config.yaml")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer otelInstance.Shutdown(context.Background())
func InitFromFile(configPath string) (*OTEL, error) {

	cfg, loaderror := config.Load(configPath)

	if loaderror != nil {

		return nil, fmt.Errorf("load config: %w", loaderror)
	}

	return Init(cfg)
}

// InitFromEnv initializes all OTEL components from environment variables.
// Environment variables are loaded via config.LoadFromEnv() and passed to Init().
//
// Common environment variables:
//   - OTEL_SERVICE_NAME: Service name for identification
//   - OTEL_SERVICE_VERSION: Service version
//   - OTEL_ENVIRONMENT: Deployment environment
//   - OTEL_EXPORTER_OTLP_ENDPOINT: OTLP collector endpoint
//   - OTEL_EXPORTER_OTLP_PROTOCOL: Protocol (grpc or http)
//   - OTEL_LOG_LEVEL: Logging level (debug, info, warn, error)
//
// Returns:
//   - *OTEL: The initialized OTEL instance
//   - error: Non-nil if initialization fails
//
// Usage:
//
//	// Set environment variables before starting
//	// export OTEL_SERVICE_NAME=myapp
//	// export OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317
//
//	otelInstance, err := otel.InitFromEnv()
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer otelInstance.Shutdown(context.Background())
func InitFromEnv() (*OTEL, error) {

	cfg := config.LoadFromEnv()

	return Init(cfg)
}

// MustInit initializes all OTEL components and panics on error.
// Use this in main() or init() where initialization failures should be fatal.
//
// Parameters:
//   - cfg: Unified configuration
//
// Returns:
//   - *OTEL: The initialized OTEL instance (panics on error)
//
// Usage:
//
//	func main() {
//	    otelInstance := otel.MustInit(config.Default())
//	    defer otelInstance.Shutdown(context.Background())
//	    // ... application code ...
//	}
func MustInit(cfg config.Config) *OTEL {

	otelInstance, initerror := Init(cfg)

	if initerror != nil {
		panic(fmt.Sprintf("failed to initialize OTEL: %v", initerror))
	}

	return otelInstance
}

// MustInitFromFile initializes all OTEL components from a YAML file and panics on error.
// Use this in main() or init() where initialization failures should be fatal.
//
// Parameters:
//   - configPath: Path to the YAML configuration file
//
// Returns:
//   - *OTEL: The initialized OTEL instance (panics on error)
//
// Usage:
//
//	func main() {
//	    otelInstance := otel.MustInitFromFile("config.yaml")
//	    defer otelInstance.Shutdown(context.Background())
//	    // ... application code ...
//	}
func MustInitFromFile(configPath string) *OTEL {

	otelInstance, initerror := InitFromFile(configPath)

	if initerror != nil {
		panic(fmt.Sprintf("failed to initialize OTEL from file: %v", initerror))
	}

	return otelInstance
}

// MustInitFromEnv initializes all OTEL components from environment and panics on error.
// Use this in main() or init() where initialization failures should be fatal.
//
// Returns:
//   - *OTEL: The initialized OTEL instance (panics on error)
//
// Usage:
//
//	func main() {
//	    otelInstance := otel.MustInitFromEnv()
//	    defer otelInstance.Shutdown(context.Background())
//	    // ... application code ...
//	}
func MustInitFromEnv() *OTEL {

	otelInstance, initerror := InitFromEnv()

	if initerror != nil {
		panic(fmt.Sprintf("failed to initialize OTEL from env: %v", initerror))
	}

	return otelInstance
}

/* ---------------------------------------- OTEL Instance Methods -------------------------------------------------- */

// Shutdown gracefully shuts down all OTEL components.
// Should be called before application exit, typically with defer.
//
// Shutdown order (important for proper cleanup):
//  1. Tracer - finish in-flight traces first
//  2. Metrics - flush pending metrics
//  3. Logger - close last so shutdown errors can be logged
//
// The shutdown is coordinated to ensure all telemetry data is exported
// before the application exits. Errors from individual shutdowns are
// collected and returned as a combined error.
//
// Parameters:
//   - ctx: Context with optional timeout for graceful shutdown.
//     Use context.Background() for unlimited wait, or
//     context.WithTimeout() to set a deadline.
//
// Returns:
//   - error: Combined errors from all shutdown operations
//
// Usage:
//
//	// With unlimited timeout
//	defer otelInstance.Shutdown(context.Background())
//
//	// With 5 second timeout
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	if err := otelInstance.Shutdown(ctx); err != nil {
//	    log.Printf("OTEL shutdown error: %v", err)
//	}
func (otelInstance *OTEL) Shutdown(ctx context.Context) error {

	// Use ShutdownCollector to aggregate errors from multiple shutdowns
	collector := common.NewShutdownCollector()

	// 1. Shutdown tracer first to finish any in-flight traces
	// This ensures all trace data is exported before we shut down
	if otelInstance.Tracer != nil {
		collector.Collect(tracer.Shutdown(ctx))
	}

	// 2. Shutdown metrics to flush any pending metric data
	if otelInstance.metrics {
		collector.Collect(metrics.Shutdown(ctx))
	}

	// 3. Shutdown logger last so we can log any shutdown errors
	// The logger should remain available until the very end
	if otelInstance.Logger != nil {
		collector.Collect(logger.Close())
	}

	return collector.Error()
}

/* ---------------------------------------- Package-level Functions -------------------------------------------------- */

// Shutdown is a package-level function to shutdown all global OTEL components.
// Use this when you don't have access to the OTEL instance or when using
// package-level initialization functions directly.
//
// This function shuts down:
//   - Global tracer provider
//   - Global metrics provider
//   - Global logger
//
// Parameters:
//   - ctx: Context with optional timeout for graceful shutdown
//
// Returns:
//   - error: Combined errors from all shutdown operations
//
// Usage:
//
//	// When using package-level functions directly
//	logger.Init(loggerConfig)
//	tracer.Init(tracerConfig)
//	metrics.Init(metricsConfig)
//
//	// ... application code ...
//
//	// Shutdown all at once
//	if err := otel.Shutdown(ctx); err != nil {
//	    log.Printf("OTEL shutdown error: %v", err)
//	}
func Shutdown(ctx context.Context) error {

	collector := common.NewShutdownCollector()

	// 1. Shutdown tracer first (to finish any in-flight traces)
	collector.Collect(tracer.Shutdown(ctx))

	// 2. Shutdown metrics
	collector.Collect(metrics.Shutdown(ctx))

	// 3. Shutdown logger last (so we can log any shutdown errors)
	collector.Collect(logger.Close())

	return collector.Error()
}
