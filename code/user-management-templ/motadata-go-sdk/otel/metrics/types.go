// Package metrics provides OpenTelemetry-based metrics collection for the motadata-go-sdk.
// It supports counters, gauges, and histograms with OTLP export capabilities.
//
// The package provides:
//   - Multiple metric types (Counter, Gauge, Histogram, Timer)
//   - Zero-allocation pooling for high-performance metric operations
//   - Label/attribute management with efficient merging
//   - Thread-safe timer implementation with automatic pooling
//
// Usage:
//
//	// Initialize metrics
//	metrics.Init(config)
//	defer metrics.Shutdown(ctx)
//
//	// Create and use metrics
//	counter := metrics.NewCounter("requests_total", "Total requests")
//	counter.Inc(ctx)
//
//	histogram := metrics.NewHistogram("latency_ms", "Request latency")
//	timer := histogram.Timer(ctx)
//	defer timer.Stop()
package metrics

import (
	"context"
	"sync/atomic"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/core/pool/resourcepool"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"
)

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// MetricType represents the type of metric instrument.
// Each type has different characteristics for how values are recorded and aggregated.
type MetricType int

// Labels represents metric labels/tags as key-value pairs.
// Labels are used to add dimensions to metrics for filtering and grouping.
// Example: Labels{"service": "api", "method": "GET", "status": "200"}
type Labels map[string]string

// MetricOpts holds common options for creating metric instruments.
// These options are shared across all metric types (counter, gauge, histogram).
type MetricOpts struct {
	// Name is the unique identifier for the metric.
	// Should follow naming conventions: snake_case with meaningful prefixes.
	// Example: "http_requests_total", "db_connections_active"
	Name string

	// Description provides human-readable documentation for the metric.
	// This is displayed in metric backends and documentation.
	Description string

	// Unit specifies the measurement unit (e.g., "ms", "bytes", "1").
	// Use UCUM conventions where possible.
	Unit string

	// Labels are default labels attached to all measurements of this metric.
	// Additional labels can be provided at recording time.
	Labels Labels
}

// HistogramOpts extends MetricOpts with histogram-specific options.
// Histograms record distributions of values in configurable buckets.
type HistogramOpts struct {
	// Embedded common metric options
	MetricOpts

	// Buckets defines the histogram bucket boundaries.
	// Values are upper bounds (exclusive) for each bucket.
	// Example: []float64{10, 50, 100, 500, 1000} for latency in ms
	Buckets []float64
}

/* ---------------------------------------- Constants -------------------------------------------------- */

const (
	// CounterType represents a monotonically increasing metric.
	// Counters can only increase or be reset to zero.
	// Use for: request counts, error counts, total bytes processed.
	CounterType MetricType = iota

	// GaugeType represents a metric that can go up and down.
	// Use for: current temperature, memory usage, active connections.
	GaugeType

	// HistogramType represents a distribution of values.
	// Records observations in pre-defined buckets.
	// Use for: request latency, response size, queue wait time.
	HistogramType
)

// DefaultHistogramBuckets provides sensible defaults for latency measurements (in milliseconds).
// These buckets cover typical API response times from 1ms to 10 seconds.
// Bucket boundaries: 1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000 ms
var DefaultHistogramBuckets = []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000}

// DefaultSizeBuckets provides sensible defaults for size measurements (in bytes).
// These buckets cover typical payload sizes from 100 bytes to 100MB.
// Bucket boundaries: 100B, 1KB, 10KB, 100KB, 1MB, 10MB, 100MB
var DefaultSizeBuckets = []float64{100, 1000, 10000, 100000, 1000000, 10000000, 100000000}

/* ---------------------------------------- MetricType Methods -------------------------------------------------- */

// String returns the human-readable name of the metric type.
// Used for logging and debugging purposes.
func (metricType MetricType) String() string {

	switch metricType {

	case CounterType:

		return "counter"

	case GaugeType:

		return "gauge"

	case HistogramType:

		return "histogram"

	default:

		return "unknown"
	}
}

/* ---------------------------------------- Labels Methods -------------------------------------------------- */

// Merge combines two label sets, with the other set taking precedence on conflicts.
// This creates a new Labels map without modifying either input.
//
// Parameters:
//   - other: Labels to merge with (takes precedence on key conflicts)
//
// Returns:
//   - Labels: New merged label set
//
// Example:
//
//	base := Labels{"service": "api", "env": "prod"}
//	extra := Labels{"method": "GET", "env": "staging"}
//	merged := base.Merge(extra)
//	// Result: {"service": "api", "env": "staging", "method": "GET"}
func (labels Labels) Merge(other Labels) Labels {

	result := make(Labels, len(labels)+len(other))

	// Copy base labels first
	for key, value := range labels {
		result[key] = value
	}

	// Overlay with other labels (takes precedence)
	for key, value := range other {
		result[key] = value
	}

	return result
}

/* ---------------------------------------- Metric Interfaces -------------------------------------------------- */

// Counter is a monotonically increasing metric.
// Counter values can only increase or be reset to zero on restart.
// This is the most common metric type, used for counting events.
//
// Usage:
//
//	counter := metrics.NewCounter("http_requests_total", "Total HTTP requests")
//	counter.Inc(ctx)                           // Increment by 1
//	counter.Add(ctx, 5)                        // Add arbitrary value
//	counter.Inc(ctx, Labels{"method": "GET"})  // With labels
type Counter interface {
	// Inc increments the counter by 1.
	// Additional labels can be provided to add dimensions to this measurement.
	Inc(ctx context.Context, labels ...Labels)

	// Add adds an arbitrary non-negative value to the counter.
	// Negative values are silently ignored (counters can only increase).
	Add(ctx context.Context, value float64, labels ...Labels)

	// Name returns the metric name.
	Name() string
}

// Gauge is a metric that can go up and down.
// Gauges represent instantaneous values that can increase or decrease.
// Unlike counters, gauges can be set to arbitrary values.
//
// Usage:
//
//	gauge := metrics.NewGauge("active_connections", "Current active connections")
//	gauge.Set(ctx, 42)                         // Set absolute value
//	gauge.Inc(ctx)                             // Increment by 1
//	gauge.Dec(ctx)                             // Decrement by 1
//	gauge.Add(ctx, -5)                         // Add/subtract value
type Gauge interface {
	// Set sets the gauge to an absolute value.
	Set(ctx context.Context, value float64, labels ...Labels)

	// Inc increments the gauge by 1.
	Inc(ctx context.Context, labels ...Labels)

	// Dec decrements the gauge by 1.
	Dec(ctx context.Context, labels ...Labels)

	// Add adds a value to the gauge (can be negative).
	Add(ctx context.Context, value float64, labels ...Labels)

	// Name returns the metric name.
	Name() string
}

// Histogram records observations and counts them in configurable buckets.
// Histograms are used to measure distributions of values like latency or size.
// Each bucket counts observations less than or equal to its upper bound.
//
// Usage:
//
//	histogram := metrics.NewHistogram("request_duration_ms", "Request duration")
//	histogram.Observe(ctx, 42.5)                          // Record observation
//	histogram.ObserveDuration(ctx, startTime)             // Record duration since start
//	histogram.Observe(ctx, 100, Labels{"status": "200"})  // With labels
type Histogram interface {
	// Observe records a value in the histogram.
	// The value is placed in the appropriate bucket based on bucket boundaries.
	Observe(ctx context.Context, value float64, labels ...Labels)

	// ObserveDuration records the duration since start time in milliseconds.
	// This is a convenience method for timing operations.
	ObserveDuration(ctx context.Context, start time.Time, labels ...Labels)

	// Name returns the metric name.
	Name() string
}

// Timer is a helper for timing operations with automatic recording.
// Timers are created from histograms and automatically record duration when stopped.
// Uses object pooling for zero-allocation operation in hot paths.
//
// Usage:
//
//	timer := metrics.NewTimer(ctx, "operation_duration_ms", "Operation duration")
//	defer timer.Stop()                                    // Records duration on stop
//	// ... do work ...
//
//	// Or with labels:
//	timer := metrics.NewTimer(ctx, "operation_duration_ms", "Operation duration")
//	defer timer.StopWithLabels(Labels{"status": "success"})
type Timer interface {
	// Stop stops the timer and records the duration to the histogram.
	// Safe to call multiple times (subsequent calls are no-ops).
	Stop()

	// StopWithLabels stops the timer and records duration with additional labels.
	// Labels are merged with any labels set at timer creation.
	StopWithLabels(labels Labels)
}

// Provider is the interface for metric backends.
// Providers are responsible for creating metric instruments and managing their lifecycle.
// The default provider uses OpenTelemetry OTLP export.
type Provider interface {
	// Counter creates or retrieves a counter metric instrument.
	Counter(opts MetricOpts) (Counter, error)

	// Gauge creates or retrieves a gauge metric instrument.
	Gauge(opts MetricOpts) (Gauge, error)

	// Histogram creates or retrieves a histogram metric instrument.
	Histogram(opts HistogramOpts) (Histogram, error)

	// Shutdown gracefully shuts down the provider and flushes pending metrics.
	Shutdown(ctx context.Context) error
}

/* ---------------------------------------- Timer Pool Configuration -------------------------------------------------- */

// timerPool reduces allocations for timer creation using ResourcePool.
// Timers are frequently created and destroyed, making pooling essential for performance.
var timerPool *resourcepool.Pool[*timerImpl]

// initTimerPool initializes the timer pool during package initialization.
// Uses utils constants for pool configuration to maintain consistency.
func init() {

	pool, err := resourcepool.NewResourcePool(resourcepool.PoolConfig[*timerImpl]{
		// Maximum number of timers to keep in pool
		// Higher values use more memory but reduce allocation overhead
		MaxSize: int(utils.OTELTimerPoolMaxSize),

		// OnCreate allocates a new timer when pool is empty
		OnCreate: func() (*timerImpl, error) {
			return &timerImpl{}, nil
		},

		// OnReset clears timer state when returned to pool
		// Note: stopped flag is NOT reset here - it must remain true
		// to prevent double-stop issues. It's reset in getTimer() instead.
		OnReset: func(timer *timerImpl) error {
			timer.histogram = nil
			timer.labels = nil
			return nil
		},
	})

	if err != nil {
		panic("failed to initialize timer pool: " + err.Error())
	}

	timerPool = pool
}

/* ---------------------------------------- Timer Implementation -------------------------------------------------- */

// timerImpl implements Timer with atomic stop flag for thread safety.
// The timer is automatically returned to the pool when stopped.
type timerImpl struct {
	// histogram is the target histogram for duration recording
	histogram Histogram

	// start is when the timer was created
	start time.Time

	// labels are default labels to include with the measurement
	labels Labels

	// stopped tracks whether Stop() has been called (prevents double-stop)
	stopped atomic.Bool
}

/* ---------------------------------------- Timer Methods -------------------------------------------------- */

// Stop stops the timer and records the duration to the associated histogram.
// This method is safe to call multiple times - only the first call records.
// The timer is automatically returned to the pool after stopping.
func (timer *timerImpl) Stop() {

	// Atomically check and set stopped flag to prevent double recording
	// CompareAndSwap returns true only if we successfully changed false->true
	if !timer.stopped.CompareAndSwap(false, true) {

		return
	}

	// Capture references before clearing (required for pool safety)
	// These values must be captured before returning timer to pool
	histogram := timer.histogram
	start := timer.start
	labels := timer.labels

	// Clear references and return timer to pool immediately
	// This allows the timer to be reused while we record the observation
	timer.histogram = nil
	timer.labels = nil
	_ = timerPool.Put(timer)

	// Record the observation after returning to pool
	// This ordering is safe because we captured all needed values above
	histogram.ObserveDuration(context.Background(), start, labels)
}

// StopWithLabels stops the timer and records the duration with additional labels.
// Additional labels are merged with any labels set at timer creation.
// This method is safe to call multiple times - only the first call records.
func (timer *timerImpl) StopWithLabels(labels Labels) {

	// Atomically check and set stopped flag to prevent double recording
	if !timer.stopped.CompareAndSwap(false, true) {

		return
	}

	// Capture references before clearing
	histogram := timer.histogram
	start := timer.start
	timerLabels := timer.labels

	// Clear references and return timer to pool
	timer.histogram = nil
	timer.labels = nil
	_ = timerPool.Put(timer)

	// Record the observation with merged labels
	// Timer labels serve as base, additional labels override on conflict
	ctx := context.Background()
	if timerLabels == nil {
		histogram.ObserveDuration(ctx, start, labels)
	} else {
		mergedLabels := timerLabels.Merge(labels)
		histogram.ObserveDuration(ctx, start, mergedLabels)
	}
}

/* ---------------------------------------- Timer Pool Helper Functions -------------------------------------------------- */

// getTimer retrieves a timer from the pool or creates a new one if pool is exhausted.
// The stopped flag is reset here to ensure the timer is ready for use.
//
// Returns:
//   - *timerImpl: A timer instance ready for use
func getTimer() *timerImpl {

	timer, err := timerPool.TryGet()

	if err != nil {
		// Pool exhausted - create new timer as fallback
		// This ensures we never block on timer creation
		return &timerImpl{}
	}

	// Reset stopped flag for new use
	// This MUST be done here (not in OnReset) to prevent race conditions
	// where a timer pointer could be double-stopped
	timer.stopped.Store(false)

	return timer
}
