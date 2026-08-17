package metrics

import (
	"context"
	"sync/atomic"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/core/pool/resourcepool"
)

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// MetricType represents the type of metric
type MetricType int

// Labels represents metric labels/tags as key-value pairs
type Labels map[string]string

// MetricOpts holds common options for creating metrics
type MetricOpts struct {
	Name string

	Description string

	Unit string

	Labels Labels
}

// HistogramOpts extends MetricOpts with histogram-specific options
type HistogramOpts struct {
	MetricOpts

	Buckets []float64
}

/* ---------------------------------------- Constants -------------------------------------------------- */

const (
	CounterType MetricType = iota

	GaugeType

	HistogramType
)

// DefaultHistogramBuckets provides sensible defaults for latency measurements (in milliseconds)
var DefaultHistogramBuckets = []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000}

// DefaultSizeBuckets provides sensible defaults for size measurements (in bytes)
var DefaultSizeBuckets = []float64{100, 1000, 10000, 100000, 1000000, 10000000, 100000000}

/* ---------------------------------------- MetricType Methods -------------------------------------------------- */

// String returns the metric type name
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

// Merge merges two label sets, with other taking precedence
func (labels Labels) Merge(other Labels) Labels {

	result := make(Labels, len(labels)+len(other))

	for key, value := range labels {
		result[key] = value
	}

	for key, value := range other {
		result[key] = value
	}

	return result
}

/* ---------------------------------------- Metric Interfaces -------------------------------------------------- */

// Counter is a monotonically increasing metric
type Counter interface {
	Inc(ctx context.Context, labels ...Labels)

	Add(ctx context.Context, value float64, labels ...Labels)

	Name() string
}

// Gauge is a metric that can go up and down
type Gauge interface {
	Set(ctx context.Context, value float64, labels ...Labels)

	Inc(ctx context.Context, labels ...Labels)

	Dec(ctx context.Context, labels ...Labels)

	Add(ctx context.Context, value float64, labels ...Labels)

	Name() string
}

// Histogram records observations and counts them in configurable buckets
type Histogram interface {
	Observe(ctx context.Context, value float64, labels ...Labels)

	ObserveDuration(ctx context.Context, start time.Time, labels ...Labels)

	Name() string
}

// Timer is a helper for timing operations
type Timer interface {
	Stop()

	StopWithLabels(labels Labels)
}

// Provider is the interface for metric backends
type Provider interface {
	Counter(opts MetricOpts) (Counter, error)

	Gauge(opts MetricOpts) (Gauge, error)

	Histogram(opts HistogramOpts) (Histogram, error)

	Shutdown(ctx context.Context) error
}

/* ---------------------------------------- Timer Implementation -------------------------------------------------- */

// Maximum pool size for timers
const timerPoolMaxSize = 1024

// timerPool reduces allocations for timer creation using ResourcePool
var timerPool *resourcepool.Pool[*timerImpl]

// initTimerPool initializes the timer pool
func init() {

	pool, err := resourcepool.NewResourcePool(resourcepool.PoolConfig[*timerImpl]{
		MaxSize: timerPoolMaxSize,
		OnCreate: func() (*timerImpl, error) {
			return &timerImpl{}, nil
		},
		OnReset: func(timer *timerImpl) error {
			// Only clear references, don't reset stopped flag here
			// The stopped flag must remain true to prevent double-stop panics
			// It will be reset in getTimer() when the timer is acquired for new use
			timer.ctx = nil
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

// timerImpl implements Timer with atomic stop flag
type timerImpl struct {
	ctx context.Context

	histogram Histogram

	start time.Time

	labels Labels

	stopped atomic.Bool
}

/* ---------------------------------------- Timer Methods -------------------------------------------------- */

// Stop stops the timer and records the duration
func (timer *timerImpl) Stop() {

	if !timer.stopped.CompareAndSwap(false, true) {

		return
	}

	// Capture references before clearing
	histogram := timer.histogram
	ctx := timer.ctx
	start := timer.start
	labels := timer.labels

	// Clear and return to pool
	timer.ctx = nil
	timer.histogram = nil
	timer.labels = nil
	_ = timerPool.Put(timer)

	// Record the observation after returning to pool
	histogram.ObserveDuration(ctx, start, labels)
}

// StopWithLabels stops the timer and records the duration with additional labels
func (timer *timerImpl) StopWithLabels(labels Labels) {

	if !timer.stopped.CompareAndSwap(false, true) {

		return
	}

	// Capture references before clearing
	histogram := timer.histogram
	ctx := timer.ctx
	start := timer.start
	timerLabels := timer.labels

	// Clear and return to pool
	timer.ctx = nil
	timer.histogram = nil
	timer.labels = nil
	_ = timerPool.Put(timer)

	// Record the observation after returning to pool
	if timerLabels == nil {
		histogram.ObserveDuration(ctx, start, labels)
	} else {
		mergedLabels := timerLabels.Merge(labels)
		histogram.ObserveDuration(ctx, start, mergedLabels)
	}
}

/* ---------------------------------------- Timer Pool Helper Functions -------------------------------------------------- */

// getTimer gets a timer from the pool
func getTimer() *timerImpl {

	timer, err := timerPool.TryGet()

	if err != nil {
		// Fallback to creating a new timer if pool is exhausted
		return &timerImpl{}
	}

	// Reset stopped flag for new use (must be done here, not in OnReset,
	// to prevent double-stop panics on the same timer pointer)
	timer.stopped.Store(false)

	return timer
}
