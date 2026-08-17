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
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

/* ---------------------------------------- Constants -------------------------------------------------- */

const maxLabels = 16

// Maximum pool sizes
const (
	labelsPoolMaxSize     = 1024
	attributesPoolMaxSize = 1024
	pooledAttrsMaxSize    = 1024
)

/* ---------------------------------------- Object Pool Types -------------------------------------------------- */

// LabelsWrapper wraps a Labels map for pooling
type LabelsWrapper struct {
	Labels Labels
}

// AttributesWrapper wraps an attribute slice for pooling
type AttributesWrapper struct {
	Attrs []attribute.KeyValue
}

// pooledAttrs holds pooled resources for attribute conversion
// Using a struct avoids closure allocation overhead
type pooledAttrs struct {
	attrs         []attribute.KeyValue
	attrsWrapper  *AttributesWrapper
	labelsWrapper *LabelsWrapper
}

// release returns the pooled resources to their pools
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

var (
	labelsPool      *resourcepool.Pool[*LabelsWrapper]
	attributesPool  *resourcepool.Pool[*AttributesWrapper]
	pooledAttrsPool *resourcepool.Pool[*pooledAttrs]
)

// initMetricsPools initializes all metric pools
func init() {

	var err error

	// Initialize labels pool
	labelsPool, err = resourcepool.NewResourcePool(resourcepool.PoolConfig[*LabelsWrapper]{
		MaxSize: labelsPoolMaxSize,
		OnCreate: func() (*LabelsWrapper, error) {
			return &LabelsWrapper{
				Labels: make(Labels, maxLabels),
			}, nil
		},
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

	// Initialize attributes pool
	attributesPool, err = resourcepool.NewResourcePool(resourcepool.PoolConfig[*AttributesWrapper]{
		MaxSize: attributesPoolMaxSize,
		OnCreate: func() (*AttributesWrapper, error) {
			return &AttributesWrapper{
				Attrs: make([]attribute.KeyValue, 0, maxLabels),
			}, nil
		},
		OnReset: func(wrapper *AttributesWrapper) error {
			wrapper.Attrs = wrapper.Attrs[:0]
			return nil
		},
	})

	if err != nil {
		panic("failed to initialize attributes pool: " + err.Error())
	}

	// Initialize pooledAttrs pool
	pooledAttrsPool, err = resourcepool.NewResourcePool(resourcepool.PoolConfig[*pooledAttrs]{
		MaxSize: pooledAttrsMaxSize,
		OnCreate: func() (*pooledAttrs, error) {
			return &pooledAttrs{}, nil
		},
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

// getPooledAttrs gets a pooledAttrs from the pool
func getPooledAttrs() *pooledAttrs {

	p, err := pooledAttrsPool.TryGet()

	if err != nil {

		return &pooledAttrs{}
	}

	return p
}

// putPooledAttrs returns a pooledAttrs to the pool
func putPooledAttrs(p *pooledAttrs) {

	if p == nil {

		return
	}

	p.release()
	p.attrs = nil
	_ = pooledAttrsPool.Put(p)
}

// getLabelsWrapper gets a labels wrapper from the pool
func getLabelsWrapper() *LabelsWrapper {

	wrapper, err := labelsPool.TryGet()

	if err != nil {

		return &LabelsWrapper{
			Labels: make(Labels, maxLabels),
		}
	}

	return wrapper
}

// putLabelsWrapper returns a labels wrapper to the pool
func putLabelsWrapper(wrapper *LabelsWrapper) {

	if wrapper == nil {

		return
	}

	_ = labelsPool.Put(wrapper)
}

// getAttributesWrapper gets an attributes wrapper from the pool
func getAttributesWrapper() *AttributesWrapper {

	wrapper, err := attributesPool.TryGet()

	if err != nil {

		return &AttributesWrapper{
			Attrs: make([]attribute.KeyValue, 0, maxLabels),
		}
	}

	return wrapper
}

// putAttributesWrapper returns an attributes wrapper to the pool
func putAttributesWrapper(wrapper *AttributesWrapper) {

	if wrapper == nil {

		return
	}

	_ = attributesPool.Put(wrapper)
}

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// OTELProvider implements Provider using OpenTelemetry
type OTELProvider struct {
	mu sync.RWMutex

	meter metric.Meter

	meterProvider *sdkmetric.MeterProvider

	config Config

	counters map[string]*otelCounter

	gauges map[string]*otelGauge

	histograms map[string]*otelHistogram

	closed atomic.Bool
}

/* ---------------------------------------- Constructor Functions -------------------------------------------------- */

// NewOTELProvider creates a new OpenTelemetry metrics provider
func NewOTELProvider(cfg Config) (*OTELProvider, error) {

	ctx := context.Background()

	// Create resource using common utility
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
		// Create OTLP exporter based on protocol
		exporter, err := createOTLPMetricExporter(ctx, cfg)

		if err != nil {

			return nil, fmt.Errorf("create exporter: %w", err)
		}

		meterProvider = sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(resource),
			sdkmetric.WithReader(
				sdkmetric.NewPeriodicReader(exporter,
					sdkmetric.WithInterval(cfg.ExportInterval),
				),
			),
		)
	} else {
		// Create a no-op provider for local development
		meterProvider = sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(resource),
		)
	}

	// Set global meter provider
	otel.SetMeterProvider(meterProvider)

	// Create meter
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
func createOTLPMetricExporter(ctx context.Context, cfg Config) (sdkmetric.Exporter, error) {

	protocol, err := common.ResolveProtocol(cfg.OTELProtocol)

	if err != nil {

		return nil, err
	}

	switch protocol {

	case common.ProtocolGRPC:

		return createGRPCMetricExporter(ctx, cfg)

	case common.ProtocolHTTP:

		return createHTTPMetricExporter(ctx, cfg)

	default:

		return nil, fmt.Errorf("unsupported OTEL protocol: %s", protocol)
	}
}

// createGRPCMetricExporter creates a gRPC-based OTLP metric exporter.
func createGRPCMetricExporter(ctx context.Context, cfg Config) (sdkmetric.Exporter, error) {

	options := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(cfg.OTELEndpoint),
	}

	if cfg.OTELInsecure {

		dialOption := common.GRPCDialOption(true)

		if dialOption != nil {
			options = append(options, otlpmetricgrpc.WithDialOption(dialOption))
		}
	}

	return otlpmetricgrpc.New(ctx, options...)
}

// createHTTPMetricExporter creates an HTTP-based OTLP metric exporter.
func createHTTPMetricExporter(ctx context.Context, cfg Config) (sdkmetric.Exporter, error) {

	options := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(cfg.OTELEndpoint),
	}

	if cfg.OTELInsecure {
		options = append(options, otlpmetrichttp.WithInsecure())
	}

	return otlpmetrichttp.New(ctx, options...)
}

/* ---------------------------------------- Provider Methods -------------------------------------------------- */

// Counter creates or retrieves a counter metric
func (otelProvider *OTELProvider) Counter(options MetricOpts) (Counter, error) {

	otelProvider.mu.Lock()
	defer otelProvider.mu.Unlock()

	if existing, exists := otelProvider.counters[options.Name]; exists {

		return existing, nil
	}

	instrument, err := otelProvider.meter.Float64Counter(options.Name,
		metric.WithDescription(options.Description),
		metric.WithUnit(options.Unit),
	)

	if err != nil {

		return nil, fmt.Errorf("create counter %s: %w", options.Name, err)
	}

	counter := &otelCounter{
		name:          options.Name,
		counter:       instrument,
		defaultLabels: mergeLabels(otelProvider.config.DefaultLabels, options.Labels),
	}

	otelProvider.counters[options.Name] = counter

	return counter, nil
}

// Gauge creates or retrieves a gauge metric
func (otelProvider *OTELProvider) Gauge(options MetricOpts) (Gauge, error) {

	otelProvider.mu.Lock()
	defer otelProvider.mu.Unlock()

	if existing, exists := otelProvider.gauges[options.Name]; exists {

		return existing, nil
	}

	instrument, err := otelProvider.meter.Float64UpDownCounter(options.Name,
		metric.WithDescription(options.Description),
		metric.WithUnit(options.Unit),
	)

	if err != nil {

		return nil, fmt.Errorf("create gauge %s: %w", options.Name, err)
	}

	gauge := &otelGauge{
		name:          options.Name,
		gauge:         instrument,
		defaultLabels: mergeLabels(otelProvider.config.DefaultLabels, options.Labels),
		values:        make(map[string]float64),
	}

	otelProvider.gauges[options.Name] = gauge

	return gauge, nil
}

// Histogram creates or retrieves a histogram metric
func (otelProvider *OTELProvider) Histogram(options HistogramOpts) (Histogram, error) {

	otelProvider.mu.Lock()
	defer otelProvider.mu.Unlock()

	if existing, exists := otelProvider.histograms[options.Name]; exists {

		return existing, nil
	}

	buckets := options.Buckets

	if len(buckets) == 0 {
		buckets = DefaultHistogramBuckets
	}

	instrument, err := otelProvider.meter.Float64Histogram(options.Name,
		metric.WithDescription(options.Description),
		metric.WithUnit(options.Unit),
		metric.WithExplicitBucketBoundaries(buckets...),
	)

	if err != nil {

		return nil, fmt.Errorf("create histogram %s: %w", options.Name, err)
	}

	histogram := &otelHistogram{
		name:          options.Name,
		histogram:     instrument,
		defaultLabels: mergeLabels(otelProvider.config.DefaultLabels, options.Labels),
	}

	otelProvider.histograms[options.Name] = histogram

	return histogram, nil
}

// Shutdown gracefully shuts down the provider. Safe to call multiple times.
func (otelProvider *OTELProvider) Shutdown(ctx context.Context) error {

	// Ensure Shutdown is only executed once
	if !otelProvider.closed.CompareAndSwap(false, true) {

		return nil
	}

	if otelProvider.meterProvider != nil {

		return otelProvider.meterProvider.Shutdown(ctx)
	}

	return nil
}

/* ---------------------------------------- OTEL Metric Implementations -------------------------------------------------- */

type otelCounter struct {
	name          string
	counter       metric.Float64Counter
	defaultLabels Labels
}

func (counter *otelCounter) Name() string {

	return counter.name
}

func (counter *otelCounter) Inc(ctx context.Context, labels ...Labels) {

	counter.Add(ctx, 1, labels...)
}

func (counter *otelCounter) Add(ctx context.Context, value float64, labels ...Labels) {

	if value < 0 {

		return // counters can only increase
	}

	// Fast path: no default labels and no additional labels
	if len(counter.defaultLabels) == 0 && len(labels) == 0 {
		counter.counter.Add(ctx, value)

		return
	}

	pooled := labelsToAttributesPooled(counter.defaultLabels, labels...)
	counter.counter.Add(ctx, value, metric.WithAttributes(pooled.attrs...))
	putPooledAttrs(pooled)
}

type otelGauge struct {
	name          string
	gauge         metric.Float64UpDownCounter
	defaultLabels Labels
	mu            sync.RWMutex
	values        map[string]float64 // track current values per label set
}

func (gauge *otelGauge) Name() string {

	return gauge.name
}

func (gauge *otelGauge) Set(ctx context.Context, value float64, labels ...Labels) {

	// Calculate delta from previous value to set absolute value
	key := labelsToKey(gauge.defaultLabels, labels...)

	gauge.mu.Lock()
	previous := gauge.values[key]
	delta := value - previous
	gauge.values[key] = value
	gauge.mu.Unlock()

	// Only add if there's a change
	if delta != 0 {

		pooled := labelsToAttributesPooled(gauge.defaultLabels, labels...)
		gauge.gauge.Add(ctx, delta, metric.WithAttributes(pooled.attrs...))
		putPooledAttrs(pooled)
	}
}

func (gauge *otelGauge) Inc(ctx context.Context, labels ...Labels) {

	gauge.Add(ctx, 1, labels...)
}

func (gauge *otelGauge) Dec(ctx context.Context, labels ...Labels) {

	gauge.Add(ctx, -1, labels...)
}

func (gauge *otelGauge) Add(ctx context.Context, value float64, labels ...Labels) {

	pooled := labelsToAttributesPooled(gauge.defaultLabels, labels...)
	gauge.gauge.Add(ctx, value, metric.WithAttributes(pooled.attrs...))
	putPooledAttrs(pooled)
}

type otelHistogram struct {
	name          string
	histogram     metric.Float64Histogram
	defaultLabels Labels
}

func (histogram *otelHistogram) Name() string {

	return histogram.name
}

func (histogram *otelHistogram) Observe(ctx context.Context, value float64, labels ...Labels) {

	// Fast path: no default labels and no additional labels
	if len(histogram.defaultLabels) == 0 && len(labels) == 0 {
		histogram.histogram.Record(ctx, value)

		return
	}

	pooled := labelsToAttributesPooled(histogram.defaultLabels, labels...)
	histogram.histogram.Record(ctx, value, metric.WithAttributes(pooled.attrs...))
	putPooledAttrs(pooled)
}

func (histogram *otelHistogram) ObserveDuration(ctx context.Context, start time.Time, labels ...Labels) {

	duration := float64(time.Since(start).Milliseconds())
	histogram.Observe(ctx, duration, labels...)
}

/* ---------------------------------------- Helper Functions -------------------------------------------------- */

func mergeLabels(base Labels, overlay Labels) Labels {

	if len(base) == 0 && len(overlay) == 0 {

		return nil
	}

	result := make(Labels, len(base)+len(overlay))

	for key, value := range base {
		result[key] = value
	}

	for key, value := range overlay {
		result[key] = value
	}

	return result
}

// labelsToAttributesPooled converts labels to OTEL attributes using pooled allocation.
// Returns a pooledAttrs struct that must be released after use via putPooledAttrs.
// This version avoids closure allocation overhead by returning a struct instead.
func labelsToAttributesPooled(defaultLabels Labels, additionalLabels ...Labels) *pooledAttrs {

	result := getPooledAttrs()

	// Get pooled map for merging
	result.labelsWrapper = getLabelsWrapper()

	// Merge all labels
	for key, value := range defaultLabels {
		result.labelsWrapper.Labels[key] = value
	}

	for _, labelSet := range additionalLabels {

		for key, value := range labelSet {
			result.labelsWrapper.Labels[key] = value
		}
	}

	// Fast path: no labels
	if len(result.labelsWrapper.Labels) == 0 {
		result.attrs = nil

		return result
	}

	// Get pooled attributes slice
	result.attrsWrapper = getAttributesWrapper()

	// Convert to attributes
	for key, value := range result.labelsWrapper.Labels {
		result.attrsWrapper.Attrs = append(result.attrsWrapper.Attrs, attribute.String(key, value))
	}

	result.attrs = result.attrsWrapper.Attrs

	return result
}

// labelsToKey creates a unique string key from labels for tracking gauge values.
// Labels are sorted to ensure consistent keys regardless of iteration order.
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

	if len(merged) == 0 {

		return ""
	}

	// Sort keys for consistent ordering
	keys := make([]string, 0, len(merged))

	for key := range merged {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	// Build key string
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
