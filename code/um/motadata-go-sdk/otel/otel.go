// Package otel provides unified OpenTelemetry initialization for logging, metrics, and tracing.
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
type OTEL struct {
	Logger *logger.Logger

	Tracer *tracer.Tracer

	config config.Config

	metrics bool
}

/* ---------------------------------------- Constructor Functions -------------------------------------------------- */

// Init initializes all OTEL components (logger, metrics, tracer) from the unified config.
// Returns an OTEL instance that can be used to shutdown all components.
func Init(cfg config.Config) (*OTEL, error) {

	otelInstance := &OTEL{
		config: cfg,
	}

	// Initialize logger
	loggerConfig := cfg.GetLoggerConfig()
	loggerInstance, loggererror := logger.Init(loggerConfig)

	if loggererror != nil {

		return nil, fmt.Errorf("initialize logger: %w", loggererror)
	}

	otelInstance.Logger = loggerInstance

	// Initialize metrics
	metricsConfig := cfg.GetMetricsConfig()

	if metricsConfig.Enabled {

		metricserror := metrics.Init(metricsConfig)

		if metricserror != nil {
			// Cleanup logger on error
			_ = logger.Close()

			return nil, fmt.Errorf("initialize metrics: %w", metricserror)
		}

		otelInstance.metrics = true
	}

	// Initialize tracer
	tracerConfig := cfg.GetTracerConfig()

	if tracerConfig.Enabled {

		tracerInstance, tracererror := tracer.Init(tracerConfig)

		if tracererror != nil {
			// Cleanup on error
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
func InitFromFile(configPath string) (*OTEL, error) {

	cfg, loaderror := config.Load(configPath)

	if loaderror != nil {

		return nil, fmt.Errorf("load config: %w", loaderror)
	}

	return Init(cfg)
}

// InitFromEnv initializes all OTEL components from environment variables.
func InitFromEnv() (*OTEL, error) {

	cfg := config.LoadFromEnv()

	return Init(cfg)
}

// MustInit initializes all OTEL components and panics on error.
func MustInit(cfg config.Config) *OTEL {

	otelInstance, initerror := Init(cfg)

	if initerror != nil {
		panic(fmt.Sprintf("failed to initialize OTEL: %v", initerror))
	}

	return otelInstance
}

// MustInitFromFile initializes all OTEL components from a YAML file and panics on error.
func MustInitFromFile(configPath string) *OTEL {

	otelInstance, initerror := InitFromFile(configPath)

	if initerror != nil {
		panic(fmt.Sprintf("failed to initialize OTEL from file: %v", initerror))
	}

	return otelInstance
}

// MustInitFromEnv initializes all OTEL components from environment and panics on error.
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
func (otelInstance *OTEL) Shutdown(ctx context.Context) error {

	collector := common.NewShutdownCollector()

	// Shutdown tracer first (to finish any in-flight traces)
	if otelInstance.Tracer != nil {
		collector.Collect(tracer.Shutdown(ctx))
	}

	// Shutdown metrics
	if otelInstance.metrics {
		collector.Collect(metrics.Shutdown(ctx))
	}

	// Shutdown logger last (so we can log any shutdown errors)
	if otelInstance.Logger != nil {
		collector.Collect(logger.Close())
	}

	return collector.Error()
}

/* ---------------------------------------- Package-level Functions -------------------------------------------------- */

// Shutdown is a package-level function to shutdown all global OTEL components.
func Shutdown(ctx context.Context) error {

	collector := common.NewShutdownCollector()

	// Shutdown tracer first (to finish any in-flight traces)
	collector.Collect(tracer.Shutdown(ctx))

	// Shutdown metrics
	collector.Collect(metrics.Shutdown(ctx))

	// Shutdown logger last (so we can log any shutdown errors)
	collector.Collect(logger.Close())

	return collector.Error()
}
