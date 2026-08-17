package metrics

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// timeNow is a variable for testing time-dependent functionality.
// Replace with a mock function in tests to control time.
var timeNow = time.Now

// Registry manages metric registration and retrieval with namespace support.
// It provides a simplified API for creating and accessing metric instruments,
// with automatic deduplication and caching.
//
// Registry handles:
//   - Metric name prefixing with namespace
//   - Caching of created instruments for retrieval by name
//   - Graceful fallback to no-op metrics on creation errors
//   - Timer creation with automatic histogram binding
//
// Usage:
//
//	registry := metrics.NewRegistry(provider, "myapp")
//	counter := registry.Counter("requests", "Total requests")
//	counter.Inc(ctx)
type Registry struct {
	// mu protects access to metric maps
	mu sync.RWMutex

	// provider is the underlying metric provider (OTEL or no-op)
	provider Provider

	// counters maps namespaced names to counter instruments
	counters map[string]Counter

	// gauges maps namespaced names to gauge instruments
	gauges map[string]Gauge

	// histograms maps namespaced names to histogram instruments
	histograms map[string]Histogram

	// namespace is the prefix for all metrics in this registry
	namespace string
}

// RegistryStats contains statistics about registered metrics.
// Useful for debugging and monitoring metric registration.
type RegistryStats struct {
	// Counters is the number of registered counter instruments
	Counters int

	// Gauges is the number of registered gauge instruments
	Gauges int

	// Histograms is the number of registered histogram instruments
	Histograms int
}

/* ---------------------------------------- Constructor Functions -------------------------------------------------- */

// NewRegistry creates a new metric registry with the given provider and namespace.
// All metrics created from this registry will be prefixed with the namespace.
//
// Parameters:
//   - provider: The underlying metric provider
//   - namespace: Prefix for all metrics (use empty string for no prefix)
//
// Returns:
//   - *Registry: A new registry instance
//
// Usage:
//
//	// Create registry with namespace
//	httpRegistry := metrics.NewRegistry(provider, "http")
//	// Creates metrics like: http_requests_total, http_duration_ms
//
//	// Create registry without namespace
//	globalRegistry := metrics.NewRegistry(provider, "")
func NewRegistry(provider Provider, namespace string) *Registry {

	return &Registry{
		provider:   provider,
		counters:   make(map[string]Counter),
		gauges:     make(map[string]Gauge),
		histograms: make(map[string]Histogram),
		namespace:  namespace,
	}
}

/* ---------------------------------------- Registry Methods -------------------------------------------------- */

// fullName returns the namespaced metric name.
// If namespace is empty, returns the name unchanged.
//
// Parameters:
//   - name: The base metric name
//
// Returns:
//   - string: The namespaced name (e.g., "namespace_name")
func (registry *Registry) fullName(name string) string {

	if registry.namespace == "" {

		return name
	}

	return registry.namespace + "_" + name
}

// Counter creates or retrieves a counter metric.
// If a counter with the same namespaced name exists, returns the cached instance.
// On creation error, returns a no-op counter that silently discards data.
//
// Parameters:
//   - name: Metric name (will be prefixed with namespace)
//   - description: Human-readable description
//   - labels: Optional default labels for all measurements
//
// Returns:
//   - Counter: The counter instrument (never nil)
//
// Usage:
//
//	counter := registry.Counter("requests", "Total requests received")
//	counter.Inc(ctx)
//	counter.Add(ctx, 5, Labels{"method": "GET"})
func (registry *Registry) Counter(name, description string, labels ...Labels) Counter {

	registry.mu.Lock()
	defer registry.mu.Unlock()

	fullName := registry.fullName(name)

	// Return cached counter if exists
	if existing, exists := registry.counters[fullName]; exists {

		return existing
	}

	// Extract default labels from variadic argument
	var defaultLabels Labels

	if len(labels) > 0 {
		defaultLabels = labels[0]
	}

	// Create counter via provider
	counter, err := registry.provider.Counter(MetricOpts{
		Name:        fullName,
		Description: description,
		Labels:      defaultLabels,
	})

	if err != nil {
		// Return no-op counter on error to prevent nil pointer issues
		// Errors are silently swallowed - consider logging in production
		return &noopCounter{name: fullName}
	}

	// Cache and return the counter
	registry.counters[fullName] = counter

	return counter
}

// CounterWithUnit creates a counter with a specific unit.
// Similar to Counter but includes unit metadata for better documentation.
//
// Parameters:
//   - name: Metric name (will be prefixed with namespace)
//   - description: Human-readable description
//   - unit: Measurement unit (e.g., "1", "bytes", "requests")
//   - labels: Optional default labels
//
// Returns:
//   - Counter: The counter instrument (never nil)
func (registry *Registry) CounterWithUnit(name, description, unit string, labels ...Labels) Counter {

	registry.mu.Lock()
	defer registry.mu.Unlock()

	fullName := registry.fullName(name)

	// Return cached counter if exists
	if existing, exists := registry.counters[fullName]; exists {

		return existing
	}

	var defaultLabels Labels

	if len(labels) > 0 {
		defaultLabels = labels[0]
	}

	counter, err := registry.provider.Counter(MetricOpts{
		Name:        fullName,
		Description: description,
		Unit:        unit,
		Labels:      defaultLabels,
	})

	if err != nil {

		return &noopCounter{name: fullName}
	}

	registry.counters[fullName] = counter

	return counter
}

// Gauge creates or retrieves a gauge metric.
// If a gauge with the same namespaced name exists, returns the cached instance.
// On creation error, returns a no-op gauge that silently discards data.
//
// Parameters:
//   - name: Metric name (will be prefixed with namespace)
//   - description: Human-readable description
//   - labels: Optional default labels
//
// Returns:
//   - Gauge: The gauge instrument (never nil)
//
// Usage:
//
//	gauge := registry.Gauge("connections", "Current active connections")
//	gauge.Set(ctx, 42)
//	gauge.Inc(ctx)
//	gauge.Dec(ctx)
func (registry *Registry) Gauge(name, description string, labels ...Labels) Gauge {

	registry.mu.Lock()
	defer registry.mu.Unlock()

	fullName := registry.fullName(name)

	// Return cached gauge if exists
	if existing, exists := registry.gauges[fullName]; exists {

		return existing
	}

	var defaultLabels Labels

	if len(labels) > 0 {
		defaultLabels = labels[0]
	}

	gauge, err := registry.provider.Gauge(MetricOpts{
		Name:        fullName,
		Description: description,
		Labels:      defaultLabels,
	})

	if err != nil {

		return &noopGauge{name: fullName}
	}

	registry.gauges[fullName] = gauge

	return gauge
}

// GaugeWithUnit creates a gauge with a specific unit.
// Similar to Gauge but includes unit metadata for better documentation.
//
// Parameters:
//   - name: Metric name (will be prefixed with namespace)
//   - description: Human-readable description
//   - unit: Measurement unit (e.g., "bytes", "percent", "celsius")
//   - labels: Optional default labels
//
// Returns:
//   - Gauge: The gauge instrument (never nil)
func (registry *Registry) GaugeWithUnit(name, description, unit string, labels ...Labels) Gauge {

	registry.mu.Lock()
	defer registry.mu.Unlock()

	fullName := registry.fullName(name)

	// Return cached gauge if exists
	if existing, exists := registry.gauges[fullName]; exists {

		return existing
	}

	var defaultLabels Labels

	if len(labels) > 0 {
		defaultLabels = labels[0]
	}

	gauge, err := registry.provider.Gauge(MetricOpts{
		Name:        fullName,
		Description: description,
		Unit:        unit,
		Labels:      defaultLabels,
	})

	if err != nil {

		return &noopGauge{name: fullName}
	}

	registry.gauges[fullName] = gauge

	return gauge
}

// Histogram creates or retrieves a histogram metric.
// If a histogram with the same namespaced name exists, returns the cached instance.
// Uses default latency buckets if no custom buckets are specified.
// On creation error, returns a no-op histogram that silently discards data.
//
// Parameters:
//   - name: Metric name (will be prefixed with namespace)
//   - description: Human-readable description
//   - labels: Optional default labels
//
// Returns:
//   - Histogram: The histogram instrument (never nil)
//
// Usage:
//
//	histogram := registry.Histogram("duration", "Request duration")
//	histogram.Observe(ctx, 42.5)
//	histogram.ObserveDuration(ctx, startTime)
func (registry *Registry) Histogram(name, description string, labels ...Labels) Histogram {

	registry.mu.Lock()
	defer registry.mu.Unlock()

	fullName := registry.fullName(name)

	// Return cached histogram if exists
	if existing, exists := registry.histograms[fullName]; exists {

		return existing
	}

	var defaultLabels Labels

	if len(labels) > 0 {
		defaultLabels = labels[0]
	}

	histogram, err := registry.provider.Histogram(HistogramOpts{
		MetricOpts: MetricOpts{
			Name:        fullName,
			Description: description,
			Labels:      defaultLabels,
		},
		// Default buckets will be used by provider
	})

	if err != nil {

		return &noopHistogram{name: fullName}
	}

	registry.histograms[fullName] = histogram

	return histogram
}

// HistogramWithBuckets creates a histogram with custom bucket boundaries.
// Use this when default buckets don't fit your use case (e.g., for size metrics).
//
// Parameters:
//   - name: Metric name (will be prefixed with namespace)
//   - description: Human-readable description
//   - unit: Measurement unit (e.g., "ms", "bytes")
//   - buckets: Explicit bucket boundaries
//   - labels: Optional default labels
//
// Returns:
//   - Histogram: The histogram instrument (never nil)
//
// Usage:
//
//	histogram := registry.HistogramWithBuckets(
//	    "response_size",
//	    "Response size in bytes",
//	    "bytes",
//	    []float64{100, 1000, 10000, 100000},
//	)
func (registry *Registry) HistogramWithBuckets(name, description, unit string, buckets []float64, labels ...Labels) Histogram {

	registry.mu.Lock()
	defer registry.mu.Unlock()

	fullName := registry.fullName(name)

	// Return cached histogram if exists
	if existing, exists := registry.histograms[fullName]; exists {

		return existing
	}

	var defaultLabels Labels

	if len(labels) > 0 {
		defaultLabels = labels[0]
	}

	histogram, err := registry.provider.Histogram(HistogramOpts{
		MetricOpts: MetricOpts{
			Name:        fullName,
			Description: description,
			Unit:        unit,
			Labels:      defaultLabels,
		},
		Buckets: buckets,
	})

	if err != nil {

		return &noopHistogram{name: fullName}
	}

	registry.histograms[fullName] = histogram

	return histogram
}

// Timer creates a new timer for measuring duration.
// The timer is bound to a histogram and records the duration when stopped.
// Uses object pooling for zero-allocation operation in hot paths.
//
// Parameters:
//   - ctx: Context for the operation
//   - name: Histogram metric name (will be prefixed with namespace)
//   - description: Human-readable description
//   - labels: Optional default labels
//
// Returns:
//   - Timer: A started timer (never nil)
//
// Usage:
//
//	timer := registry.Timer(ctx, "operation_duration", "Operation duration")
//	defer timer.Stop()
//	// ... perform operation ...
func (registry *Registry) Timer(ctx context.Context, name, description string, labels ...Labels) Timer {

	// Get or create the histogram for this timer
	histogram := registry.Histogram(name, description, labels...)

	// Get a timer from the pool
	timer := getTimer()
	timer.histogram = histogram
	timer.start = timeNow()
	timer.labels = mergeAllLabels(labels...)

	return timer
}

// Shutdown gracefully shuts down the registry.
// Delegates to the underlying provider's Shutdown method.
//
// Parameters:
//   - ctx: Context with optional timeout for graceful shutdown
//
// Returns:
//   - error: Non-nil if shutdown encounters errors
func (registry *Registry) Shutdown(ctx context.Context) error {

	return registry.provider.Shutdown(ctx)
}

// Stats returns statistics about registered metrics.
// Useful for debugging and monitoring metric registration patterns.
//
// Returns:
//   - RegistryStats: Statistics about registered instruments
func (registry *Registry) Stats() RegistryStats {

	registry.mu.RLock()
	defer registry.mu.RUnlock()

	return RegistryStats{
		Counters:   len(registry.counters),
		Gauges:     len(registry.gauges),
		Histograms: len(registry.histograms),
	}
}

// String returns a human-readable string representation of RegistryStats.
func (stats RegistryStats) String() string {

	return fmt.Sprintf("Registry{counters=%d, gauges=%d, histograms=%d}",
		stats.Counters, stats.Gauges, stats.Histograms)
}

/* ---------------------------------------- Helper Functions -------------------------------------------------- */

// mergeAllLabels merges multiple label sets into a single Labels map.
// Later labels take precedence on key conflicts.
//
// Parameters:
//   - labels: Variadic label sets to merge
//
// Returns:
//   - Labels: Merged labels or nil if no labels provided
func mergeAllLabels(labels ...Labels) Labels {

	if len(labels) == 0 {

		return nil
	}

	result := make(Labels)

	for _, labelSet := range labels {

		for key, value := range labelSet {
			result[key] = value
		}
	}

	return result
}

/* ---------------------------------------- No-op Implementations -------------------------------------------------- */

// noopCounter implements Counter with no-op behavior.
// Used as fallback when metric creation fails to prevent nil pointer issues.
type noopCounter struct {
	// name is preserved for Name() method
	name string
}

// Name returns the counter's metric name.
func (counter *noopCounter) Name() string {

	return counter.name
}

// Inc is a no-op that discards the increment.
// Intentionally empty: this is a fallback implementation used when metric creation
// fails, ensuring callers can safely call methods without nil checks.
func (counter *noopCounter) Inc(ctx context.Context, labels ...Labels) {
	// No operation - metric data is intentionally discarded
}

// Add is a no-op that discards the value.
// Intentionally empty: this is a fallback implementation used when metric creation
// fails, ensuring callers can safely call methods without nil checks.
func (counter *noopCounter) Add(ctx context.Context, value float64, labels ...Labels) {
	// No operation - metric data is intentionally discarded
}

// noopGauge implements Gauge with no-op behavior.
// Used as fallback when metric creation fails to prevent nil pointer issues.
type noopGauge struct {
	// name is preserved for Name() method
	name string
}

// Name returns the gauge's metric name.
func (gauge *noopGauge) Name() string {

	return gauge.name
}

// Set is a no-op that discards the value.
// Intentionally empty: this is a fallback implementation used when metric creation
// fails, ensuring callers can safely call methods without nil checks.
func (gauge *noopGauge) Set(ctx context.Context, value float64, labels ...Labels) {
	// No operation - metric data is intentionally discarded
}

// Inc is a no-op that discards the increment.
// Intentionally empty: this is a fallback implementation used when metric creation
// fails, ensuring callers can safely call methods without nil checks.
func (gauge *noopGauge) Inc(ctx context.Context, labels ...Labels) {
	// No operation - metric data is intentionally discarded
}

// Dec is a no-op that discards the decrement.
// Intentionally empty: this is a fallback implementation used when metric creation
// fails, ensuring callers can safely call methods without nil checks.
func (gauge *noopGauge) Dec(ctx context.Context, labels ...Labels) {
	// No operation - metric data is intentionally discarded
}

// Add is a no-op that discards the value.
// Intentionally empty: this is a fallback implementation used when metric creation
// fails, ensuring callers can safely call methods without nil checks.
func (gauge *noopGauge) Add(ctx context.Context, value float64, labels ...Labels) {
	// No operation - metric data is intentionally discarded
}

// noopHistogram implements Histogram with no-op behavior.
// Used as fallback when metric creation fails to prevent nil pointer issues.
type noopHistogram struct {
	// name is preserved for Name() method
	name string
}

// Name returns the histogram's metric name.
func (histogram *noopHistogram) Name() string {

	return histogram.name
}

// Observe is a no-op that discards the observation.
// Intentionally empty: this is a fallback implementation used when metric creation
// fails, ensuring callers can safely call methods without nil checks.
func (histogram *noopHistogram) Observe(ctx context.Context, value float64, labels ...Labels) {
	// No operation - metric data is intentionally discarded
}

// ObserveDuration is a no-op that discards the duration observation.
// Intentionally empty: this is a fallback implementation used when metric creation
// fails, ensuring callers can safely call methods without nil checks.
func (histogram *noopHistogram) ObserveDuration(ctx context.Context, start time.Time, labels ...Labels) {
	// No operation - metric data is intentionally discarded
}
