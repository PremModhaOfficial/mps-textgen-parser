package metrics

import (
	"context"
	"fmt"
	"sync"
	"time"
)

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// timeNow is a variable for testing
var timeNow = time.Now

// Registry manages metric registration and retrieval
type Registry struct {
	mu sync.RWMutex

	provider Provider

	counters map[string]Counter

	gauges map[string]Gauge

	histograms map[string]Histogram

	namespace string
}

// RegistryStats contains registry statistics
type RegistryStats struct {
	Counters int

	Gauges int

	Histograms int
}

/* ---------------------------------------- Constructor Functions -------------------------------------------------- */

// NewRegistry creates a new metric registry
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

// fullName returns the namespaced metric name
func (registry *Registry) fullName(name string) string {

	if registry.namespace == "" {

		return name
	}

	return registry.namespace + "_" + name
}

// Counter creates or retrieves a counter metric
func (registry *Registry) Counter(name, description string, labels ...Labels) Counter {

	registry.mu.Lock()
	defer registry.mu.Unlock()

	fullName := registry.fullName(name)

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
		Labels:      defaultLabels,
	})

	if err != nil {
		// Return a no-op counter on error
		return &noopCounter{name: fullName}
	}

	registry.counters[fullName] = counter

	return counter
}

// CounterWithUnit creates a counter with a specific unit
func (registry *Registry) CounterWithUnit(name, description, unit string, labels ...Labels) Counter {

	registry.mu.Lock()
	defer registry.mu.Unlock()

	fullName := registry.fullName(name)

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

// Gauge creates or retrieves a gauge metric
func (registry *Registry) Gauge(name, description string, labels ...Labels) Gauge {

	registry.mu.Lock()
	defer registry.mu.Unlock()

	fullName := registry.fullName(name)

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

// GaugeWithUnit creates a gauge with a specific unit
func (registry *Registry) GaugeWithUnit(name, description, unit string, labels ...Labels) Gauge {

	registry.mu.Lock()
	defer registry.mu.Unlock()

	fullName := registry.fullName(name)

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

// Histogram creates or retrieves a histogram metric
func (registry *Registry) Histogram(name, description string, labels ...Labels) Histogram {

	registry.mu.Lock()
	defer registry.mu.Unlock()

	fullName := registry.fullName(name)

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
	})

	if err != nil {

		return &noopHistogram{name: fullName}
	}

	registry.histograms[fullName] = histogram

	return histogram
}

// HistogramWithBuckets creates a histogram with custom buckets
func (registry *Registry) HistogramWithBuckets(name, description, unit string, buckets []float64, labels ...Labels) Histogram {

	registry.mu.Lock()
	defer registry.mu.Unlock()

	fullName := registry.fullName(name)

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

// Timer creates a new timer for measuring duration
func (registry *Registry) Timer(ctx context.Context, name, description string, labels ...Labels) Timer {

	histogram := registry.Histogram(name, description, labels...)

	timer := getTimer()
	timer.ctx = ctx
	timer.histogram = histogram
	timer.start = timeNow()
	timer.labels = mergeAllLabels(labels...)

	return timer
}

// Shutdown gracefully shuts down the registry
func (registry *Registry) Shutdown(ctx context.Context) error {

	return registry.provider.Shutdown(ctx)
}

// Stats returns statistics about registered metrics
func (registry *Registry) Stats() RegistryStats {

	registry.mu.RLock()
	defer registry.mu.RUnlock()

	return RegistryStats{
		Counters:   len(registry.counters),
		Gauges:     len(registry.gauges),
		Histograms: len(registry.histograms),
	}
}

// string returns a string representation of RegistryStats
func (stats RegistryStats) string() string {

	return fmt.Sprintf("Registry{counters=%d, gauges=%d, histograms=%d}",
		stats.Counters, stats.Gauges, stats.Histograms)
}

/* ---------------------------------------- Helper Functions -------------------------------------------------- */

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

type noopCounter struct {
	name string
}

func (counter *noopCounter) Name() string {

	return counter.name
}

func (counter *noopCounter) Inc(ctx context.Context, labels ...Labels) {}

func (counter *noopCounter) Add(ctx context.Context, value float64, labels ...Labels) {}

type noopGauge struct {
	name string
}

func (gauge *noopGauge) Name() string {

	return gauge.name
}

func (gauge *noopGauge) Set(ctx context.Context, value float64, labels ...Labels) {}

func (gauge *noopGauge) Inc(ctx context.Context, labels ...Labels) {}

func (gauge *noopGauge) Dec(ctx context.Context, labels ...Labels) {}

func (gauge *noopGauge) Add(ctx context.Context, value float64, labels ...Labels) {}

type noopHistogram struct {
	name string
}

func (histogram *noopHistogram) Name() string {

	return histogram.name
}

func (histogram *noopHistogram) Observe(ctx context.Context, value float64, labels ...Labels) {}

func (histogram *noopHistogram) ObserveDuration(ctx context.Context, start time.Time, labels ...Labels) {
}
