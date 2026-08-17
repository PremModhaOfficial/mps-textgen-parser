package metrics

import (
	"context"
	"fmt"
	"sync"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/config"
)

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// Config is an alias for config.MetricsConfig for backward compatibility.
// This allows the metrics package to be used without importing config directly.
type Config = config.MetricsConfig

// ServiceMetrics provides a convenient way to create common metrics for a service.
// It automatically namespaces all metrics and provides helpers for common patterns.
//
// Usage:
//
//	sm := metrics.NewServiceMetrics("myservice")
//	sm.Requests("api.get").Inc(ctx)
//	sm.Duration("api.get").Observe(ctx, 42.5)
//	timer := sm.Timer(ctx, "api.get")
//	defer timer.Stop()
type ServiceMetrics struct {
	// namespace is the prefix for all metrics created by this instance
	namespace string

	// registry manages metric registration and caching
	registry *Registry
}

/* ---------------------------------------- Global State -------------------------------------------------- */

// Global state for the metrics system.
// Protected by globalMu for thread-safe access.
var (
	// globalMu protects access to global provider and registry
	globalMu sync.RWMutex

	// globalProvider is the active metrics provider (OTEL or no-op)
	globalProvider Provider

	// globalRegistry manages global metric registration
	globalRegistry *Registry
)

/* ---------------------------------------- Configuration Functions -------------------------------------------------- */

// DefaultConfig returns the default metrics configuration with sensible defaults.
// Use this as a starting point and override specific values as needed.
//
// Returns:
//   - Config: Configuration with defaults for local development
func DefaultConfig() Config {

	defaultConfig := config.DefaultMetricsConfig()

	// Set service identification defaults
	defaultConfig.ServiceName = "app"
	defaultConfig.ServiceVersion = "0.0.0"
	defaultConfig.Environment = "development"

	return defaultConfig
}

// ConfigFromEnv loads metrics configuration from environment variables.
// Environment variables take precedence over defaults.
//
// Supported environment variables:
//   - OTEL_SERVICE_NAME: Service name for identification
//   - OTEL_SERVICE_VERSION: Service version
//   - OTEL_ENVIRONMENT: Deployment environment (production, staging, etc.)
//   - OTEL_EXPORTER_OTLP_ENDPOINT: OTLP collector endpoint
//   - OTEL_EXPORTER_OTLP_PROTOCOL: Protocol (grpc or http)
//
// Returns:
//   - Config: Configuration loaded from environment
func ConfigFromEnv() Config {

	return config.LoadMetricsConfigFromEnv()
}

/* ---------------------------------------- Global Initialization Functions -------------------------------------------------- */

// Init initializes the global metrics system with the given configuration.
// This must be called before using any metric functions.
// If metrics are disabled in config, this is a no-op.
//
// Parameters:
//   - cfg: Metrics configuration
//
// Returns:
//   - error: Non-nil if initialization fails
//
// Usage:
//
//	if err := metrics.Init(cfg); err != nil {
//	    log.Fatal("metrics init failed:", err)
//	}
//	defer metrics.Shutdown(ctx)
func Init(cfg Config) error {

	// Skip initialization if metrics are disabled
	if !cfg.Enabled {

		return nil
	}

	// Create OTEL provider with the given configuration
	provider, err := NewOTELProvider(cfg)

	if err != nil {

		return fmt.Errorf("create metrics provider: %w", err)
	}

	// Set global state atomically
	globalMu.Lock()
	globalProvider = provider
	globalRegistry = NewRegistry(provider, cfg.ServiceName)
	globalMu.Unlock()

	return nil
}

// InitFromEnv initializes metrics using environment variables.
// Convenience wrapper around Init(ConfigFromEnv()).
//
// Returns:
//   - error: Non-nil if initialization fails
func InitFromEnv() error {

	cfg := ConfigFromEnv()

	return Init(cfg)
}

// MustInit initializes metrics and panics on error.
// Use this in main() or init() where errors should be fatal.
//
// Parameters:
//   - cfg: Metrics configuration
func MustInit(cfg Config) {

	if err := Init(cfg); err != nil {
		panic(fmt.Sprintf("failed to init metrics: %v", err))
	}
}

// InitFromUnifiedConfig initializes metrics from the unified SDK config.
// This is used by the main otel.Init() function.
//
// Parameters:
//   - unifiedConfig: Unified SDK configuration
//
// Returns:
//   - error: Non-nil if initialization fails
func InitFromUnifiedConfig(unifiedConfig config.Config) error {

	return Init(unifiedConfig.GetMetricsConfig())
}

// Shutdown gracefully shuts down the global metrics system.
// This flushes any pending metrics and releases resources.
// Safe to call multiple times or when metrics aren't initialized.
//
// Parameters:
//   - ctx: Context with optional timeout for graceful shutdown
//
// Returns:
//   - error: Non-nil if shutdown encounters errors
func Shutdown(ctx context.Context) error {

	globalMu.Lock()
	defer globalMu.Unlock()

	if globalProvider != nil {

		err := globalProvider.Shutdown(ctx)

		// Clear global state regardless of error
		globalProvider = nil
		globalRegistry = nil

		return err
	}

	return nil
}

/* ---------------------------------------- Global Registry Access -------------------------------------------------- */

// R returns the global registry for creating metrics.
// If metrics aren't initialized, returns a no-op registry.
// This is the primary entry point for creating metrics.
//
// Returns:
//   - *Registry: The global metric registry
//
// Usage:
//
//	counter := metrics.R().Counter("requests", "Total requests")
//	counter.Inc(ctx)
func R() *Registry {

	globalMu.RLock()
	registry := globalRegistry
	globalMu.RUnlock()

	// Return no-op registry if not initialized
	// This allows metrics code to run safely before Init()
	if registry == nil {

		return defaultRegistry()
	}

	return registry
}

// defaultRegistry creates a no-op registry for when metrics aren't initialized.
// This ensures metrics code can run without errors even before Init().
// The no-op registry is lazily created and cached.
//
// Returns:
//   - *Registry: A no-op registry that discards all metrics
func defaultRegistry() *Registry {

	globalMu.Lock()
	defer globalMu.Unlock()

	// Double-check in case another goroutine initialized
	if globalRegistry != nil {

		return globalRegistry
	}

	// Create no-op provider and registry
	globalProvider = &noopProvider{}
	globalRegistry = NewRegistry(globalProvider, "")

	return globalRegistry
}

/* ---------------------------------------- Convenience Functions -------------------------------------------------- */

// NewCounter creates or retrieves a counter from the global registry.
// Counters are monotonically increasing metrics.
//
// Parameters:
//   - name: Metric name (will be prefixed with service name)
//   - description: Human-readable description
//   - labels: Optional default labels
//
// Returns:
//   - Counter: The counter metric instrument
//
// Usage:
//
//	counter := metrics.NewCounter("http_requests_total", "Total HTTP requests")
//	counter.Inc(ctx, Labels{"method": "GET"})
func NewCounter(name, description string, labels ...Labels) Counter {

	return R().Counter(name, description, labels...)
}

// NewGauge creates or retrieves a gauge from the global registry.
// Gauges are metrics that can go up and down.
//
// Parameters:
//   - name: Metric name (will be prefixed with service name)
//   - description: Human-readable description
//   - labels: Optional default labels
//
// Returns:
//   - Gauge: The gauge metric instrument
//
// Usage:
//
//	gauge := metrics.NewGauge("active_connections", "Current active connections")
//	gauge.Set(ctx, 42)
func NewGauge(name, description string, labels ...Labels) Gauge {

	return R().Gauge(name, description, labels...)
}

// NewHistogram creates or retrieves a histogram from the global registry.
// Histograms record observations in configurable buckets.
//
// Parameters:
//   - name: Metric name (will be prefixed with service name)
//   - description: Human-readable description
//   - labels: Optional default labels
//
// Returns:
//   - Histogram: The histogram metric instrument
//
// Usage:
//
//	histogram := metrics.NewHistogram("request_duration_ms", "Request duration")
//	histogram.Observe(ctx, 42.5)
func NewHistogram(name, description string, labels ...Labels) Histogram {

	return R().Histogram(name, description, labels...)
}

// NewTimer creates a new timer for measuring duration.
// The timer automatically records to a histogram when stopped.
//
// Parameters:
//   - ctx: Context for the operation
//   - name: Metric name for the underlying histogram
//   - description: Human-readable description
//   - labels: Optional default labels
//
// Returns:
//   - Timer: A started timer
//
// Usage:
//
//	timer := metrics.NewTimer(ctx, "operation_duration_ms", "Operation duration")
//	defer timer.Stop()
//	// ... perform operation ...
func NewTimer(ctx context.Context, name, description string, labels ...Labels) Timer {

	return R().Timer(ctx, name, description, labels...)
}

/* ---------------------------------------- Namespace Registry -------------------------------------------------- */

// Namespace creates a new registry with a namespace prefix.
// All metrics created from this registry will be prefixed with the namespace.
// Useful for organizing metrics by component or subsystem.
//
// Parameters:
//   - namespace: Prefix for all metrics (e.g., "http", "db", "cache")
//
// Returns:
//   - *Registry: A new registry with the namespace
//
// Usage:
//
//	httpMetrics := metrics.Namespace("http")
//	httpMetrics.Counter("requests", "HTTP requests").Inc(ctx)
//	// Creates metric: http_requests
func Namespace(namespace string) *Registry {

	globalMu.RLock()
	provider := globalProvider
	globalMu.RUnlock()

	// Use no-op provider if not initialized
	if provider == nil {
		provider = &noopProvider{}
	}

	return NewRegistry(provider, namespace)
}

/* ---------------------------------------- Common Metric Helpers -------------------------------------------------- */

// RequestCounter creates a counter for tracking requests.
// The metric name will be suffixed with "_requests_total".
//
// Parameters:
//   - name: Base name (will have "_requests_total" appended)
//   - labels: Optional default labels
//
// Returns:
//   - Counter: The request counter
//
// Usage:
//
//	counter := metrics.RequestCounter("api")
//	// Creates: api_requests_total
func RequestCounter(name string, labels ...Labels) Counter {

	return NewCounter(name+"_requests_total", "Total number of requests", labels...)
}

// ErrorCounter creates a counter for tracking errors.
// The metric name will be suffixed with "_errors_total".
//
// Parameters:
//   - name: Base name (will have "_errors_total" appended)
//   - labels: Optional default labels
//
// Returns:
//   - Counter: The error counter
//
// Usage:
//
//	counter := metrics.ErrorCounter("api")
//	// Creates: api_errors_total
func ErrorCounter(name string, labels ...Labels) Counter {

	return NewCounter(name+"_errors_total", "Total number of errors", labels...)
}

// DurationHistogram creates a histogram for tracking request duration.
// Uses default latency buckets optimized for API response times.
// The metric name will be suffixed with "_duration_ms".
//
// Parameters:
//   - name: Base name (will have "_duration_ms" appended)
//   - labels: Optional default labels
//
// Returns:
//   - Histogram: The duration histogram
//
// Usage:
//
//	histogram := metrics.DurationHistogram("api")
//	// Creates: api_duration_ms with latency buckets
func DurationHistogram(name string, labels ...Labels) Histogram {

	return R().HistogramWithBuckets(
		name+"_duration_ms",
		"Request duration in milliseconds",
		"ms",
		DefaultHistogramBuckets,
		labels...,
	)
}

// SizeHistogram creates a histogram for tracking sizes.
// Uses default size buckets optimized for payload/file sizes.
// The metric name will be suffixed with "_bytes".
//
// Parameters:
//   - name: Base name (will have "_bytes" appended)
//   - labels: Optional default labels
//
// Returns:
//   - Histogram: The size histogram
//
// Usage:
//
//	histogram := metrics.SizeHistogram("response")
//	// Creates: response_bytes with size buckets
func SizeHistogram(name string, labels ...Labels) Histogram {

	return R().HistogramWithBuckets(
		name+"_bytes",
		"Size in bytes",
		"bytes",
		DefaultSizeBuckets,
		labels...,
	)
}

// ActiveGauge creates a gauge for tracking active items.
// The metric name will be suffixed with "_active".
//
// Parameters:
//   - name: Base name (will have "_active" appended)
//   - labels: Optional default labels
//
// Returns:
//   - Gauge: The active gauge
//
// Usage:
//
//	gauge := metrics.ActiveGauge("connections")
//	// Creates: connections_active
func ActiveGauge(name string, labels ...Labels) Gauge {

	return NewGauge(name+"_active", "Number of active items", labels...)
}

/* ---------------------------------------- Timing Helpers -------------------------------------------------- */

// TimeFunc measures the execution time of a function.
// Records the duration to a histogram after the function completes.
//
// Parameters:
//   - ctx: Context for the operation
//   - name: Histogram metric name
//   - description: Human-readable description
//   - fn: Function to time
//
// Usage:
//
//	metrics.TimeFunc(ctx, "db_query_duration_ms", "Database query duration", func() {
//	    db.Query(...)
//	})
func TimeFunc(ctx context.Context, name, description string, fn func()) {

	start := time.Now()
	fn()

	histogram := NewHistogram(name, description)
	histogram.ObserveDuration(ctx, start)
}

// TimeFuncWithLabels measures execution time with labels.
// Similar to TimeFunc but allows adding labels to the measurement.
//
// Parameters:
//   - ctx: Context for the operation
//   - name: Histogram metric name
//   - description: Human-readable description
//   - labels: Labels to attach to the measurement
//   - fn: Function to time
//
// Usage:
//
//	metrics.TimeFuncWithLabels(ctx, "db_query_duration_ms", "Query duration",
//	    Labels{"table": "users"}, func() {
//	        db.Query(...)
//	    })
func TimeFuncWithLabels(ctx context.Context, name, description string, labels Labels, fn func()) {

	start := time.Now()
	fn()

	histogram := NewHistogram(name, description)
	histogram.ObserveDuration(ctx, start, labels)
}

/* ---------------------------------------- No-op Provider -------------------------------------------------- */

// noopProvider implements Provider with no-op behavior.
// Used when metrics are disabled or before initialization.
// All metrics created by this provider silently discard data.
type noopProvider struct{}

// Counter returns a no-op counter that discards all increments.
func (noopProvider *noopProvider) Counter(options MetricOpts) (Counter, error) {

	return &noopCounter{name: options.Name}, nil
}

// Gauge returns a no-op gauge that discards all updates.
func (noopProvider *noopProvider) Gauge(options MetricOpts) (Gauge, error) {

	return &noopGauge{name: options.Name}, nil
}

// Histogram returns a no-op histogram that discards all observations.
func (noopProvider *noopProvider) Histogram(options HistogramOpts) (Histogram, error) {

	return &noopHistogram{name: options.Name}, nil
}

// Shutdown is a no-op for the no-op provider.
func (noopProvider *noopProvider) Shutdown(ctx context.Context) error {

	return nil
}

/* ---------------------------------------- ServiceMetrics -------------------------------------------------- */

// NewServiceMetrics creates a new ServiceMetrics instance.
// ServiceMetrics provides a convenient API for common metric patterns.
//
// Parameters:
//   - name: Service/component name used as namespace
//
// Returns:
//   - *ServiceMetrics: A new service metrics instance
//
// Usage:
//
//	sm := metrics.NewServiceMetrics("api")
//	sm.Requests("get_users").Inc(ctx)
//	sm.Duration("get_users").Observe(ctx, 42.5)
func NewServiceMetrics(name string) *ServiceMetrics {

	return &ServiceMetrics{
		namespace: name,
		registry:  Namespace(name),
	}
}

// Requests returns a counter for tracking requests for the given operation.
// Creates a metric named "{namespace}_requests_total" with operation label.
//
// Parameters:
//   - operation: Operation name (added as label)
//
// Returns:
//   - Counter: The request counter
func (serviceMetrics *ServiceMetrics) Requests(operation string) Counter {

	return serviceMetrics.registry.Counter(
		"requests_total",
		"Total number of requests",
		Labels{"operation": operation},
	)
}

// Errors returns a counter for tracking errors for the given operation.
// Creates a metric named "{namespace}_errors_total" with operation label.
//
// Parameters:
//   - operation: Operation name (added as label)
//
// Returns:
//   - Counter: The error counter
func (serviceMetrics *ServiceMetrics) Errors(operation string) Counter {

	return serviceMetrics.registry.Counter(
		"errors_total",
		"Total number of errors",
		Labels{"operation": operation},
	)
}

// Duration returns a histogram for tracking duration for the given operation.
// Creates a metric named "{namespace}_duration_ms" with operation label.
// Uses default latency buckets.
//
// Parameters:
//   - operation: Operation name (added as label)
//
// Returns:
//   - Histogram: The duration histogram
func (serviceMetrics *ServiceMetrics) Duration(operation string) Histogram {

	return serviceMetrics.registry.HistogramWithBuckets(
		"duration_ms",
		"Operation duration in milliseconds",
		"ms",
		DefaultHistogramBuckets,
		Labels{"operation": operation},
	)
}

// Active returns a gauge for tracking active operations.
// Creates a metric named "{namespace}_active" with operation label.
//
// Parameters:
//   - operation: Operation name (added as label)
//
// Returns:
//   - Gauge: The active gauge
func (serviceMetrics *ServiceMetrics) Active(operation string) Gauge {

	return serviceMetrics.registry.Gauge(
		"active",
		"Number of active operations",
		Labels{"operation": operation},
	)
}

// Timer starts a timer for the given operation.
// The timer records to the duration histogram when stopped.
//
// Parameters:
//   - ctx: Context for the operation
//   - operation: Operation name (added as label)
//
// Returns:
//   - Timer: A started timer
//
// Usage:
//
//	timer := sm.Timer(ctx, "get_users")
//	defer timer.Stop()
func (serviceMetrics *ServiceMetrics) Timer(ctx context.Context, operation string) Timer {

	return serviceMetrics.registry.Timer(ctx, "duration_ms", "Operation duration in milliseconds",
		Labels{"operation": operation})
}

// RecordRequest records a request with duration and optional error.
// This is a convenience method that records to multiple metrics at once.
//
// Parameters:
//   - ctx: Context for the operation
//   - operation: Operation name
//   - duration: How long the operation took
//   - err: Error if the operation failed (nil for success)
//
// Usage:
//
//	start := time.Now()
//	result, err := doOperation()
//	sm.RecordRequest(ctx, "get_users", time.Since(start), err)
func (serviceMetrics *ServiceMetrics) RecordRequest(ctx context.Context, operation string, duration time.Duration, err error) {
	// Always record request count and duration
	serviceMetrics.Requests(operation).Inc(ctx)
	serviceMetrics.Duration(operation).Observe(ctx, float64(duration.Milliseconds()))

	// Record error if operation failed
	if err != nil {
		serviceMetrics.Errors(operation).Inc(ctx)
	}
}
