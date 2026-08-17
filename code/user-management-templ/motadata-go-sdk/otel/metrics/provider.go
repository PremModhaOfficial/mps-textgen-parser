package metrics

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/core/pool/resourcepool"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/common"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

/* ---------------------------------------- Object Pool Types -------------------------------------------------- */

// LabelsWrapper wraps a Labels map for object pooling.
// This reduces allocations when merging labels in hot paths.
type LabelsWrapper struct {
	// Labels is the wrapped label map
	Labels Labels
}

// AttributesWrapper wraps an OTEL attribute slice for object pooling.
// This reduces allocations when converting labels to OTEL attributes.
type AttributesWrapper struct {
	// Attrs is the wrapped attribute slice
	Attrs []attribute.KeyValue
}

// pooledAttrs holds pooled resources for attribute conversion.
// Using a struct avoids closure allocation overhead by keeping
// all related resources together for easy cleanup.
type pooledAttrs struct {
	// attrs is the resulting attribute slice (points to attrsWrapper.Attrs)
	attrs []attribute.KeyValue

	// attrsWrapper is borrowed from the attributes pool
	attrsWrapper *AttributesWrapper

	// labelsWrapper is borrowed from the labels pool
	labelsWrapper *LabelsWrapper
}

// release returns the pooled resources to their respective pools.
// Must be called when done using the pooledAttrs to avoid leaks.
func (pooledAttrs *pooledAttrs) release() {

	if pooledAttrs.labelsWrapper != nil {
		putLabelsWrapper(pooledAttrs.labelsWrapper)
		pooledAttrs.labelsWrapper = nil
	}

	if pooledAttrs.attrsWrapper != nil {
		putAttributesWrapper(pooledAttrs.attrsWrapper)
		pooledAttrs.attrsWrapper = nil
	}
}

/* ---------------------------------------- Object Pools -------------------------------------------------- */

// Object pools for reducing allocations in hot paths.
// These pools are initialized at package startup.
var (
	// labelsPool holds reusable LabelsWrapper instances
	labelsPool *resourcepool.Pool[*LabelsWrapper]

	// attributesPool holds reusable AttributesWrapper instances
	attributesPool *resourcepool.Pool[*AttributesWrapper]

	// pooledAttrsPool holds reusable pooledAttrs instances
	pooledAttrsPool *resourcepool.Pool[*pooledAttrs]
)

// initMetricsPools initializes all metric pools during package initialization.
// Uses utils constants for pool configuration to maintain consistency across the SDK.
func init() {

	var err error

	// Initialize labels pool for label merging operations
	labelsPool, err = resourcepool.NewResourcePool(resourcepool.PoolConfig[*LabelsWrapper]{
		// Maximum number of LabelsWrapper instances to keep pooled
		MaxSize: int(utils.OTELLabelsPoolMaxSize),

		// OnCreate allocates a new LabelsWrapper with pre-sized map
		OnCreate: func() (*LabelsWrapper, error) {
			return &LabelsWrapper{
				Labels: make(Labels, int(utils.OTELMaxLabels)),
			}, nil
		},

		// OnReset clears the map without deallocating
		OnReset: func(wrapper *LabelsWrapper) error {
			for k := range wrapper.Labels {
				delete(wrapper.Labels, k)
			}
			return nil
		},
	})

	if err != nil {
		panic("failed to initialize labels pool: " + err.Error())
	}

	// Initialize attributes pool for OTEL attribute conversion
	attributesPool, err = resourcepool.NewResourcePool(resourcepool.PoolConfig[*AttributesWrapper]{
		// Maximum number of AttributesWrapper instances to keep pooled
		MaxSize: int(utils.OTELAttributesPoolMaxSize),

		// OnCreate allocates a new AttributesWrapper with pre-sized slice
		OnCreate: func() (*AttributesWrapper, error) {
			return &AttributesWrapper{
				Attrs: make([]attribute.KeyValue, 0, int(utils.OTELMaxLabels)),
			}, nil
		},

		// OnReset clears the slice while retaining capacity
		OnReset: func(wrapper *AttributesWrapper) error {
			wrapper.Attrs = wrapper.Attrs[:0]
			return nil
		},
	})

	if err != nil {
		panic("failed to initialize attributes pool: " + err.Error())
	}

	// Initialize pooledAttrs pool for coordinated resource management
	pooledAttrsPool, err = resourcepool.NewResourcePool(resourcepool.PoolConfig[*pooledAttrs]{
		// Maximum number of pooledAttrs instances to keep pooled
		MaxSize: int(utils.OTELAttributesPoolMaxSize),

		// OnCreate allocates a new pooledAttrs struct
		OnCreate: func() (*pooledAttrs, error) {
			return &pooledAttrs{}, nil
		},

		// OnReset clears all references
		OnReset: func(p *pooledAttrs) error {
			p.attrs = nil
			p.attrsWrapper = nil
			p.labelsWrapper = nil
			return nil
		},
	})

	if err != nil {
		panic("failed to initialize pooledAttrs pool: " + err.Error())
	}
}

/* ---------------------------------------- Pool Helper Functions -------------------------------------------------- */

// getPooledAttrs retrieves a pooledAttrs from the pool or creates a new one.
//
// Returns:
//   - *pooledAttrs: A pooled or new instance ready for use
func getPooledAttrs() *pooledAttrs {

	p, err := pooledAttrsPool.TryGet()

	if err != nil {
		// Pool exhausted - create new instance as fallback
		return &pooledAttrs{}
	}

	return p
}

// putPooledAttrs returns a pooledAttrs to the pool after releasing its resources.
// Safe to call with nil.
//
// Parameters:
//   - p: The pooledAttrs to return (can be nil)
func putPooledAttrs(p *pooledAttrs) {

	if p == nil {

		return
	}

	// Release borrowed resources first
	p.release()
	p.attrs = nil

	// Return to pool
	_ = pooledAttrsPool.Put(p)
}

// getLabelsWrapper retrieves a LabelsWrapper from the pool or creates a new one.
//
// Returns:
//   - *LabelsWrapper: A pooled or new instance with empty map
func getLabelsWrapper() *LabelsWrapper {

	wrapper, err := labelsPool.TryGet()

	if err != nil {
		// Pool exhausted - create new instance as fallback
		return &LabelsWrapper{
			Labels: make(Labels, int(utils.OTELMaxLabels)),
		}
	}

	return wrapper
}

// putLabelsWrapper returns a LabelsWrapper to the pool.
// Safe to call with nil.
//
// Parameters:
//   - wrapper: The LabelsWrapper to return (can be nil)
func putLabelsWrapper(wrapper *LabelsWrapper) {

	if wrapper == nil {

		return
	}

	_ = labelsPool.Put(wrapper)
}

// getAttributesWrapper retrieves an AttributesWrapper from the pool or creates a new one.
//
// Returns:
//   - *AttributesWrapper: A pooled or new instance with empty slice
func getAttributesWrapper() *AttributesWrapper {

	wrapper, err := attributesPool.TryGet()

	if err != nil {
		// Pool exhausted - create new instance as fallback
		return &AttributesWrapper{
			Attrs: make([]attribute.KeyValue, 0, int(utils.OTELMaxLabels)),
		}
	}

	return wrapper
}

// putAttributesWrapper returns an AttributesWrapper to the pool.
// Safe to call with nil.
//
// Parameters:
//   - wrapper: The AttributesWrapper to return (can be nil)
func putAttributesWrapper(wrapper *AttributesWrapper) {

	if wrapper == nil {

		return
	}

	_ = attributesPool.Put(wrapper)
}

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// OTELProvider implements Provider using OpenTelemetry.
// It manages metric instruments (counters, gauges, histograms) and their OTLP export.
//
// The provider maintains maps of created instruments to enable retrieval
// by name and prevent duplicate instrument creation.
type OTELProvider struct {
	// mu protects access to instrument maps
	mu sync.RWMutex

	// meter is the OTEL meter for creating instruments
	meter metric.Meter

	// meterProvider is the underlying OTEL meter provider
	meterProvider *sdkmetric.MeterProvider

	// config holds the provider configuration
	config Config

	// counters maps names to counter instruments
	counters map[string]*otelCounter

	// gauges maps names to gauge instruments
	gauges map[string]*otelGauge

	// histograms maps names to histogram instruments
	histograms map[string]*otelHistogram

	// closed tracks whether Shutdown has been called
	closed atomic.Bool
}

/* ---------------------------------------- Constructor Functions -------------------------------------------------- */

// NewOTELProvider creates a new OpenTelemetry metrics provider.
// The provider is configured based on the provided Config, including:
//   - OTLP exporter setup (gRPC or HTTP)
//   - Periodic metric export at configured interval
//   - Service identification via OTEL resource
//
// Parameters:
//   - cfg: Metrics configuration
//
// Returns:
//   - *OTELProvider: The configured provider
//   - error: Non-nil if provider creation fails
//
// Usage:
//
//	provider, err := metrics.NewOTELProvider(cfg)
//	if err != nil {
//	    return fmt.Errorf("create metrics provider: %w", err)
//	}
//	defer provider.Shutdown(ctx)
func NewOTELProvider(cfg Config) (*OTELProvider, error) {

	ctx := context.Background()

	// Create OTEL resource with service identification
	// This attaches service metadata to all metrics for filtering
	resource, err := common.NewOTELResourceFromConfig(
		cfg.ServiceName,
		cfg.ServiceVersion,
		cfg.Environment,
	)

	if err != nil {

		return nil, fmt.Errorf("create resource: %w", err)
	}

	var meterProvider *sdkmetric.MeterProvider

	if cfg.OTELEnabled {
		// Create OTLP exporter based on configured protocol
		exporter, err := createOTLPMetricExporter(ctx, cfg)

		if err != nil {

			return nil, fmt.Errorf("%w: %v", utils.ErrOTELExporterCreationFailed, err)
		}

		// Create meter provider with periodic reader for metric export
		// The periodic reader collects and exports metrics at the configured interval
		meterProvider = sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(resource),
			sdkmetric.WithReader(
				sdkmetric.NewPeriodicReader(exporter,
					// Export interval determines how often metrics are pushed
					// Default is typically 60s, but can be configured for faster feedback
					sdkmetric.WithInterval(cfg.ExportInterval),
				),
			),
		)
	} else {
		// Create a minimal provider for local development without export
		// Metrics are still recorded but not sent anywhere
		meterProvider = sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(resource),
		)
	}

	// Register as global meter provider
	// This allows other libraries to use the same provider
	otel.SetMeterProvider(meterProvider)

	// Create meter for this service
	meter := meterProvider.Meter(cfg.ServiceName)

	return &OTELProvider{
		meter:         meter,
		meterProvider: meterProvider,
		config:        cfg,
		counters:      make(map[string]*otelCounter),
		gauges:        make(map[string]*otelGauge),
		histograms:    make(map[string]*otelHistogram),
	}, nil
}

/* ---------------------------------------- Exporter Factory Functions -------------------------------------------------- */

// createOTLPMetricExporter creates an OTLP exporter based on the configured protocol.
// Supports both gRPC and HTTP transports.
//
// Parameters:
//   - ctx: Context for exporter creation
//   - cfg: Configuration containing endpoint and protocol settings
//
// Returns:
//   - sdkmetric.Exporter: The configured exporter
//   - error: Non-nil if exporter creation fails
func createOTLPMetricExporter(ctx context.Context, cfg Config) (sdkmetric.Exporter, error) {

	// Resolve protocol string to canonical form
	protocol, err := common.ResolveProtocol(cfg.OTELProtocol)

	if err != nil {

		return nil, err
	}

	// Create exporter based on resolved protocol
	switch protocol {

	case common.ProtocolGRPC:
		// gRPC transport - high performance, bidirectional streaming
		return createGRPCMetricExporter(ctx, cfg)

	case common.ProtocolHTTP:
		// HTTP transport - simpler, firewall-friendly
		return createHTTPMetricExporter(ctx, cfg)

	default:
		// Should not reach here if ResolveProtocol works correctly
		return nil, fmt.Errorf("%w: %s", utils.ErrOTELUnsupportedProtocol, protocol)
	}
}

// createGRPCMetricExporter creates a gRPC-based OTLP metric exporter.
// gRPC provides high performance through HTTP/2 multiplexing and binary encoding.
//
// Parameters:
//   - ctx: Context for exporter creation
//   - cfg: Configuration containing endpoint and security settings
//
// Returns:
//   - sdkmetric.Exporter: The configured gRPC exporter
//   - error: Non-nil if exporter creation fails
func createGRPCMetricExporter(ctx context.Context, cfg Config) (sdkmetric.Exporter, error) {

	// Build exporter options starting with endpoint
	options := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(cfg.OTELEndpoint),
	}

	// Add insecure transport if configured
	// WARNING: Only use insecure mode for local development
	if cfg.OTELInsecure {

		dialOption := common.GRPCDialOption(true)

		if dialOption != nil {
			options = append(options, otlpmetricgrpc.WithDialOption(dialOption))
		}
	}

	return otlpmetricgrpc.New(ctx, options...)
}

// createHTTPMetricExporter creates an HTTP-based OTLP metric exporter.
// HTTP transport uses REST conventions with protobuf encoding.
//
// Parameters:
//   - ctx: Context for exporter creation
//   - cfg: Configuration containing endpoint and security settings
//
// Returns:
//   - sdkmetric.Exporter: The configured HTTP exporter
//   - error: Non-nil if exporter creation fails
func createHTTPMetricExporter(ctx context.Context, cfg Config) (sdkmetric.Exporter, error) {

	// Build exporter options starting with endpoint
	options := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(cfg.OTELEndpoint),
	}

	// Add insecure option if configured
	// WARNING: Only use insecure mode for local development
	if cfg.OTELInsecure {
		options = append(options, otlpmetrichttp.WithInsecure())
	}

	return otlpmetrichttp.New(ctx, options...)
}

/* ---------------------------------------- Provider Methods -------------------------------------------------- */

// Counter creates or retrieves a counter metric instrument.
// If a counter with the same name exists, returns the existing one.
//
// Parameters:
//   - options: Metric options including name, description, and labels
//
// Returns:
//   - Counter: The counter instrument
//   - error: Non-nil if instrument creation fails
func (otelProvider *OTELProvider) Counter(options MetricOpts) (Counter, error) {

	otelProvider.mu.Lock()
	defer otelProvider.mu.Unlock()

	// Return existing counter if already created
	if existing, exists := otelProvider.counters[options.Name]; exists {

		return existing, nil
	}

	// Create new OTEL counter instrument
	instrument, err := otelProvider.meter.Float64Counter(options.Name,
		metric.WithDescription(options.Description),
		metric.WithUnit(options.Unit),
	)

	if err != nil {

		return nil, fmt.Errorf("%w: create counter %s: %v", utils.ErrOTELMetricCreationFailed, options.Name, err)
	}

	// Wrap and cache the counter
	counter := &otelCounter{
		name:          options.Name,
		counter:       instrument,
		defaultLabels: mergeLabels(otelProvider.config.DefaultLabels, options.Labels),
	}

	otelProvider.counters[options.Name] = counter

	return counter, nil
}

// Gauge creates or retrieves a gauge metric instrument.
// If a gauge with the same name exists, returns the existing one.
// Uses Float64UpDownCounter under the hood to support both positive and negative values.
//
// Parameters:
//   - options: Metric options including name, description, and labels
//
// Returns:
//   - Gauge: The gauge instrument
//   - error: Non-nil if instrument creation fails
func (otelProvider *OTELProvider) Gauge(options MetricOpts) (Gauge, error) {

	otelProvider.mu.Lock()
	defer otelProvider.mu.Unlock()

	// Return existing gauge if already created
	if existing, exists := otelProvider.gauges[options.Name]; exists {

		return existing, nil
	}

	// Create new OTEL UpDownCounter for gauge behavior
	// UpDownCounter allows both positive and negative changes
	instrument, err := otelProvider.meter.Float64UpDownCounter(options.Name,
		metric.WithDescription(options.Description),
		metric.WithUnit(options.Unit),
	)

	if err != nil {

		return nil, fmt.Errorf("%w: create gauge %s: %v", utils.ErrOTELMetricCreationFailed, options.Name, err)
	}

	// Wrap and cache the gauge with value tracking for Set() support
	gauge := &otelGauge{
		name:          options.Name,
		gauge:         instrument,
		defaultLabels: mergeLabels(otelProvider.config.DefaultLabels, options.Labels),
		values:        make(map[string]float64),
	}

	otelProvider.gauges[options.Name] = gauge

	return gauge, nil
}

// Histogram creates or retrieves a histogram metric instrument.
// If a histogram with the same name exists, returns the existing one.
//
// Parameters:
//   - options: Histogram options including buckets and standard metric options
//
// Returns:
//   - Histogram: The histogram instrument
//   - error: Non-nil if instrument creation fails
func (otelProvider *OTELProvider) Histogram(options HistogramOpts) (Histogram, error) {

	otelProvider.mu.Lock()
	defer otelProvider.mu.Unlock()

	// Return existing histogram if already created
	if existing, exists := otelProvider.histograms[options.Name]; exists {

		return existing, nil
	}

	// Use default buckets if none specified
	buckets := options.Buckets

	if len(buckets) == 0 {
		buckets = DefaultHistogramBuckets
	}

	// Create new OTEL histogram instrument with explicit bucket boundaries
	instrument, err := otelProvider.meter.Float64Histogram(options.Name,
		metric.WithDescription(options.Description),
		metric.WithUnit(options.Unit),
		metric.WithExplicitBucketBoundaries(buckets...),
	)

	if err != nil {

		return nil, fmt.Errorf("%w: create histogram %s: %v", utils.ErrOTELMetricCreationFailed, options.Name, err)
	}

	// Wrap and cache the histogram
	histogram := &otelHistogram{
		name:          options.Name,
		histogram:     instrument,
		defaultLabels: mergeLabels(otelProvider.config.DefaultLabels, options.Labels),
	}

	otelProvider.histograms[options.Name] = histogram

	return histogram, nil
}

// Shutdown gracefully shuts down the provider and flushes pending metrics.
// Safe to call multiple times - subsequent calls return nil.
//
// Parameters:
//   - ctx: Context with optional timeout for graceful shutdown
//
// Returns:
//   - error: Non-nil if shutdown encounters errors
func (otelProvider *OTELProvider) Shutdown(ctx context.Context) error {

	// Ensure shutdown is only executed once using atomic flag
	if !otelProvider.closed.CompareAndSwap(false, true) {

		return nil
	}

	if otelProvider.meterProvider != nil {

		return otelProvider.meterProvider.Shutdown(ctx)
	}

	return nil
}

/* ---------------------------------------- OTEL Metric Implementations -------------------------------------------------- */

// otelCounter wraps an OTEL Float64Counter with default labels support.
type otelCounter struct {
	// name is the metric name for identification
	name string

	// counter is the underlying OTEL counter instrument
	counter metric.Float64Counter

	// defaultLabels are applied to all measurements
	defaultLabels Labels
}

// Name returns the counter's metric name.
func (counter *otelCounter) Name() string {

	return counter.name
}

// Inc increments the counter by 1.
// Additional labels can be provided to add dimensions.
func (counter *otelCounter) Inc(ctx context.Context, labels ...Labels) {

	counter.Add(ctx, 1, labels...)
}

// Add adds a non-negative value to the counter.
// Negative values are silently ignored (counters can only increase).
func (counter *otelCounter) Add(ctx context.Context, value float64, labels ...Labels) {

	// Counters can only increase - ignore negative values
	if value < 0 {

		return
	}

	// Fast path: no labels to merge
	if len(counter.defaultLabels) == 0 && len(labels) == 0 {
		counter.counter.Add(ctx, value)

		return
	}

	// Convert labels to OTEL attributes using pooled allocation
	pooled := labelsToAttributesPooled(counter.defaultLabels, labels...)
	counter.counter.Add(ctx, value, metric.WithAttributes(pooled.attrs...))
	putPooledAttrs(pooled)
}

// otelGauge wraps an OTEL Float64UpDownCounter with Set() support.
// Tracks current values per label set to enable absolute value setting.
type otelGauge struct {
	// name is the metric name for identification
	name string

	// gauge is the underlying OTEL UpDownCounter
	gauge metric.Float64UpDownCounter

	// defaultLabels are applied to all measurements
	defaultLabels Labels

	// mu protects access to values map
	mu sync.RWMutex

	// values tracks current value per label set for Set() support
	values map[string]float64
}

// Name returns the gauge's metric name.
func (gauge *otelGauge) Name() string {

	return gauge.name
}

// Set sets the gauge to an absolute value.
// Calculates delta from previous value to work with UpDownCounter semantics.
func (gauge *otelGauge) Set(ctx context.Context, value float64, labels ...Labels) {

	// Create unique key for this label combination
	key := labelsToKey(gauge.defaultLabels, labels...)

	// Calculate delta from previous value
	gauge.mu.Lock()
	previous := gauge.values[key]
	delta := value - previous
	gauge.values[key] = value
	gauge.mu.Unlock()

	// Only record if there's a change
	if delta != 0 {

		pooled := labelsToAttributesPooled(gauge.defaultLabels, labels...)
		gauge.gauge.Add(ctx, delta, metric.WithAttributes(pooled.attrs...))
		putPooledAttrs(pooled)
	}
}

// Inc increments the gauge by 1.
func (gauge *otelGauge) Inc(ctx context.Context, labels ...Labels) {

	gauge.Add(ctx, 1, labels...)
}

// Dec decrements the gauge by 1.
func (gauge *otelGauge) Dec(ctx context.Context, labels ...Labels) {

	gauge.Add(ctx, -1, labels...)
}

// Add adds a value to the gauge (can be positive or negative).
func (gauge *otelGauge) Add(ctx context.Context, value float64, labels ...Labels) {

	pooled := labelsToAttributesPooled(gauge.defaultLabels, labels...)
	gauge.gauge.Add(ctx, value, metric.WithAttributes(pooled.attrs...))
	putPooledAttrs(pooled)
}

// otelHistogram wraps an OTEL Float64Histogram with default labels support.
type otelHistogram struct {
	// name is the metric name for identification
	name string

	// histogram is the underlying OTEL histogram instrument
	histogram metric.Float64Histogram

	// defaultLabels are applied to all observations
	defaultLabels Labels
}

// Name returns the histogram's metric name.
func (histogram *otelHistogram) Name() string {

	return histogram.name
}

// Observe records a value in the histogram.
// The value is placed in the appropriate bucket based on boundaries.
func (histogram *otelHistogram) Observe(ctx context.Context, value float64, labels ...Labels) {

	// Fast path: no labels to merge
	if len(histogram.defaultLabels) == 0 && len(labels) == 0 {
		histogram.histogram.Record(ctx, value)

		return
	}

	// Convert labels to OTEL attributes using pooled allocation
	pooled := labelsToAttributesPooled(histogram.defaultLabels, labels...)
	histogram.histogram.Record(ctx, value, metric.WithAttributes(pooled.attrs...))
	putPooledAttrs(pooled)
}

// ObserveDuration records the duration since start time in milliseconds.
// Convenience method for timing operations.
func (histogram *otelHistogram) ObserveDuration(ctx context.Context, start time.Time, labels ...Labels) {

	// Calculate duration in milliseconds
	duration := float64(time.Since(start).Milliseconds())
	histogram.Observe(ctx, duration, labels...)
}

/* ---------------------------------------- Helper Functions -------------------------------------------------- */

// mergeLabels merges base and overlay labels, with overlay taking precedence.
// Returns nil if both inputs are empty.
//
// Parameters:
//   - base: Base labels
//   - overlay: Labels to overlay (takes precedence)
//
// Returns:
//   - Labels: Merged labels or nil if both inputs are empty
func mergeLabels(base Labels, overlay Labels) Labels {

	if len(base) == 0 && len(overlay) == 0 {

		return nil
	}

	result := make(Labels, len(base)+len(overlay))

	// Copy base labels first
	for key, value := range base {
		result[key] = value
	}

	// Overlay takes precedence
	for key, value := range overlay {
		result[key] = value
	}

	return result
}

// labelsToAttributesPooled converts labels to OTEL attributes using pooled allocation.
// Returns a pooledAttrs struct that MUST be released via putPooledAttrs after use.
//
// This function avoids closure allocation by returning a struct instead of
// using defer. The caller is responsible for calling putPooledAttrs.
//
// Parameters:
//   - defaultLabels: Default labels to include
//   - additionalLabels: Additional labels to merge
//
// Returns:
//   - *pooledAttrs: Struct containing the converted attributes
func labelsToAttributesPooled(defaultLabels Labels, additionalLabels ...Labels) *pooledAttrs {

	result := getPooledAttrs()

	// Get pooled map for merging labels
	result.labelsWrapper = getLabelsWrapper()

	// Merge all labels into the pooled map
	for key, value := range defaultLabels {
		result.labelsWrapper.Labels[key] = value
	}

	for _, labelSet := range additionalLabels {

		for key, value := range labelSet {
			result.labelsWrapper.Labels[key] = value
		}
	}

	// Fast path: no labels to convert
	if len(result.labelsWrapper.Labels) == 0 {
		result.attrs = nil

		return result
	}

	// Get pooled attributes slice
	result.attrsWrapper = getAttributesWrapper()

	// Convert labels to OTEL attributes
	for key, value := range result.labelsWrapper.Labels {
		result.attrsWrapper.Attrs = append(result.attrsWrapper.Attrs, attribute.String(key, value))
	}

	result.attrs = result.attrsWrapper.Attrs

	return result
}

// labelsToKey creates a unique string key from labels for tracking gauge values.
// Labels are sorted to ensure consistent keys regardless of iteration order.
//
// Parameters:
//   - defaultLabels: Default labels
//   - additionalLabels: Additional labels to merge
//
// Returns:
//   - string: Unique key for this label combination
func labelsToKey(defaultLabels Labels, additionalLabels ...Labels) string {

	// Merge all labels
	merged := make(Labels)

	for key, value := range defaultLabels {
		merged[key] = value
	}

	for _, labelSet := range additionalLabels {

		for key, value := range labelSet {
			merged[key] = value
		}
	}

	// Return empty string for no labels
	if len(merged) == 0 {

		return ""
	}

	// Sort keys for consistent ordering
	keys := make([]string, 0, len(merged))

	for key := range merged {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	// Build key string: "key1=value1|key2=value2|..."
	var builder strings.Builder

	for i, key := range keys {

		if i > 0 {
			builder.WriteByte('|')
		}

		builder.WriteString(key)
		builder.WriteByte('=')
		builder.WriteString(merged[key])
	}

	return builder.String()
}
