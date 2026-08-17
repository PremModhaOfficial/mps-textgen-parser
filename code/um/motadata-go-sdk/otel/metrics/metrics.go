package metrics

import (
	"context"
	"fmt"
	"sync"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/config"
)

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// Config is an alias for config.MetricsConfig for backward compatibility
type Config = config.MetricsConfig

// ServiceMetrics provides a convenient way to create common metrics for a service
type ServiceMetrics struct {
	namespace string

	registry *Registry
}

/* ---------------------------------------- Global State -------------------------------------------------- */

var (
	globalMu sync.RWMutex

	globalProvider Provider

	globalRegistry *Registry
)

/* ---------------------------------------- Constructor Functions -------------------------------------------------- */

// DefaultConfig returns the default metrics configuration with service defaults
func DefaultConfig() Config {

	defaultConfig := config.DefaultMetricsConfig()

	defaultConfig.ServiceName = "app"
	defaultConfig.ServiceVersion = "0.0.0"
	defaultConfig.Environment = "development"

	return defaultConfig
}

// ConfigFromEnv loads metrics configuration from environment variables
func ConfigFromEnv() Config {

	return config.LoadMetricsConfigFromEnv()
}

/* ---------------------------------------- Global Initialization Functions -------------------------------------------------- */

// Init initializes the global metrics system with the given configuration
func Init(cfg Config) error {

	if !cfg.Enabled {

		return nil
	}

	provider, err := NewOTELProvider(cfg)

	if err != nil {

		return fmt.Errorf("create metrics provider: %w", err)
	}

	globalMu.Lock()
	globalProvider = provider
	globalRegistry = NewRegistry(provider, cfg.ServiceName)
	globalMu.Unlock()

	return nil
}

// InitFromEnv initializes metrics using environment variables
func InitFromEnv() error {

	cfg := ConfigFromEnv()

	return Init(cfg)
}

// MustInit initializes metrics and panics on error
func MustInit(cfg Config) {

	if err := Init(cfg); err != nil {
		panic(fmt.Sprintf("failed to init metrics: %v", err))
	}
}

// InitFromUnifiedConfig initializes metrics from the unified config
func InitFromUnifiedConfig(unifiedConfig config.Config) error {

	return Init(unifiedConfig.GetMetricsConfig())
}

// Shutdown gracefully shuts down the global metrics system
func Shutdown(ctx context.Context) error {

	globalMu.Lock()
	defer globalMu.Unlock()

	if globalProvider != nil {

		err := globalProvider.Shutdown(ctx)
		globalProvider = nil
		globalRegistry = nil

		return err
	}

	return nil
}

// R returns the global registry
func R() *Registry {

	globalMu.RLock()
	registry := globalRegistry
	globalMu.RUnlock()

	if registry == nil {

		return defaultRegistry()
	}

	return registry
}

// defaultRegistry creates a no-op registry for when metrics aren't initialized
func defaultRegistry() *Registry {

	globalMu.Lock()
	defer globalMu.Unlock()

	if globalRegistry != nil {

		return globalRegistry
	}

	globalProvider = &noopProvider{}
	globalRegistry = NewRegistry(globalProvider, "")

	return globalRegistry
}

/* ---------------------------------------- Convenience Functions -------------------------------------------------- */

// NewCounter creates or retrieves a counter from the global registry
func NewCounter(name, description string, labels ...Labels) Counter {

	return R().Counter(name, description, labels...)
}

// NewGauge creates or retrieves a gauge from the global registry
func NewGauge(name, description string, labels ...Labels) Gauge {

	return R().Gauge(name, description, labels...)
}

// NewHistogram creates or retrieves a histogram from the global registry
func NewHistogram(name, description string, labels ...Labels) Histogram {

	return R().Histogram(name, description, labels...)
}

// NewTimer creates a new timer for measuring duration
func NewTimer(ctx context.Context, name, description string, labels ...Labels) Timer {

	return R().Timer(ctx, name, description, labels...)
}

/* ---------------------------------------- Namespace Registry -------------------------------------------------- */

// Namespace creates a new registry with a namespace prefix
func Namespace(namespace string) *Registry {

	globalMu.RLock()
	provider := globalProvider
	globalMu.RUnlock()

	if provider == nil {
		provider = &noopProvider{}
	}

	return NewRegistry(provider, namespace)
}

/* ---------------------------------------- Common Metric Helpers -------------------------------------------------- */

// RequestCounter creates a counter for tracking requests
func RequestCounter(name string, labels ...Labels) Counter {

	return NewCounter(name+"_requests_total", "Total number of requests", labels...)
}

// errorCounter creates a counter for tracking errors
func ErrorCounter(name string, labels ...Labels) Counter {

	return NewCounter(name+"_errors_total", "Total number of errors", labels...)
}

// DurationHistogram creates a histogram for tracking request duration
func DurationHistogram(name string, labels ...Labels) Histogram {

	return R().HistogramWithBuckets(
		name+"_duration_ms",
		"Request duration in milliseconds",
		"ms",
		DefaultHistogramBuckets,
		labels...,
	)
}

// SizeHistogram creates a histogram for tracking sizes
func SizeHistogram(name string, labels ...Labels) Histogram {

	return R().HistogramWithBuckets(
		name+"_bytes",
		"Size in bytes",
		"bytes",
		DefaultSizeBuckets,
		labels...,
	)
}

// ActiveGauge creates a gauge for tracking active items
func ActiveGauge(name string, labels ...Labels) Gauge {

	return NewGauge(name+"_active", "Number of active items", labels...)
}

/* ---------------------------------------- Timing Helpers -------------------------------------------------- */

// TimeFunc measures the execution time of a function
func TimeFunc(ctx context.Context, name, description string, fn func()) {

	start := time.Now()
	fn()

	histogram := NewHistogram(name, description)
	histogram.ObserveDuration(ctx, start)
}

// TimeFuncWithLabels measures execution time with labels
func TimeFuncWithLabels(ctx context.Context, name, description string, labels Labels, fn func()) {

	start := time.Now()
	fn()

	histogram := NewHistogram(name, description)
	histogram.ObserveDuration(ctx, start, labels)
}

/* ---------------------------------------- No-op Provider -------------------------------------------------- */

type noopProvider struct{}

func (noopProvider *noopProvider) Counter(options MetricOpts) (Counter, error) {

	return &noopCounter{name: options.Name}, nil
}

func (noopProvider *noopProvider) Gauge(options MetricOpts) (Gauge, error) {

	return &noopGauge{name: options.Name}, nil
}

func (noopProvider *noopProvider) Histogram(options HistogramOpts) (Histogram, error) {

	return &noopHistogram{name: options.Name}, nil
}

func (noopProvider *noopProvider) Shutdown(ctx context.Context) error {

	return nil
}

/* ---------------------------------------- ServiceMetrics -------------------------------------------------- */

// NewServiceMetrics creates a new ServiceMetrics instance
func NewServiceMetrics(name string) *ServiceMetrics {

	return &ServiceMetrics{
		namespace: name,
		registry:  Namespace(name),
	}
}

// Requests returns a counter for tracking requests
func (serviceMetrics *ServiceMetrics) Requests(operation string) Counter {

	return serviceMetrics.registry.Counter(
		"requests_total",
		"Total number of requests",
		Labels{"operation": operation},
	)
}

// errors returns a counter for tracking errors
func (serviceMetrics *ServiceMetrics) errors(operation string) Counter {

	return serviceMetrics.registry.Counter(
		"errors_total",
		"Total number of errors",
		Labels{"operation": operation},
	)
}

// Duration returns a histogram for tracking duration
func (serviceMetrics *ServiceMetrics) Duration(operation string) Histogram {

	return serviceMetrics.registry.HistogramWithBuckets(
		"duration_ms",
		"Operation duration in milliseconds",
		"ms",
		DefaultHistogramBuckets,
		Labels{"operation": operation},
	)
}

// Active returns a gauge for tracking active operations
func (serviceMetrics *ServiceMetrics) Active(operation string) Gauge {

	return serviceMetrics.registry.Gauge(
		"active",
		"Number of active operations",
		Labels{"operation": operation},
	)
}

// Timer starts a timer for the given operation
func (serviceMetrics *ServiceMetrics) Timer(ctx context.Context, operation string) Timer {

	return serviceMetrics.registry.Timer(ctx, "duration_ms", "Operation duration in milliseconds",
		Labels{"operation": operation})
}

// RecordRequest records a request with duration
func (serviceMetrics *ServiceMetrics) RecordRequest(ctx context.Context, operation string, duration time.Duration, err error) {
	serviceMetrics.Requests(operation).Inc(ctx)
	serviceMetrics.Duration(operation).Observe(ctx, float64(duration.Milliseconds()))

	if err != nil {
		serviceMetrics.errors(operation).Inc(ctx)
	}
}
