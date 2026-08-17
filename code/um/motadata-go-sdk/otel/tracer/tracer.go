package tracer

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/config"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// Config is an alias for config.TracerConfig for backward compatibility.
type Config = config.TracerConfig

// Tracer wraps the OTEL tracer with convenience methods
type Tracer struct {
	tracer trace.Tracer

	provider *sdktrace.TracerProvider

	closed atomic.Bool
}

/* ---------------------------------------- Global State -------------------------------------------------- */

var (
	globalMu sync.RWMutex

	globalProvider *sdktrace.TracerProvider

	globalTracer trace.Tracer

	globalClosed atomic.Bool
)

/* ---------------------------------------- Constructor Functions -------------------------------------------------- */

// DefaultConfig returns the default tracer configuration with service defaults.
func DefaultConfig() Config {
	defaultConfig := config.DefaultTracerConfig()

	defaultConfig.ServiceName = "app"
	defaultConfig.ServiceVersion = "0.0.0"
	defaultConfig.Environment = "development"

	return defaultConfig
}

/* ---------------------------------------- Global Initialization Functions -------------------------------------------------- */

// Init initializes the global tracer with the given configuration.
// This should be called once at application startup.
func Init(config Config) (*Tracer, error) {

	if !config.Enabled {

		return nil, nil
	}

	globalClosed.Store(false)

	provider, err := newTracerProvider(config)

	if err != nil {

		return nil, err
	}

	otel.SetTracerProvider(provider)

	propagator := createPropagator(config.Propagators)
	otel.SetTextMapPropagator(propagator)

	tracer := provider.Tracer(config.ServiceName)

	globalMu.Lock()
	globalProvider = provider
	globalTracer = tracer
	globalMu.Unlock()

	return &Tracer{
		tracer:   tracer,
		provider: provider,
	}, nil
}

// InitFromEnv initializes the global tracer using environment variables.
func InitFromEnv() (*Tracer, error) {

	config := ConfigFromEnv()

	return Init(config)
}

// InitFromUnifiedConfig initializes the global tracer from the unified config.
func InitFromUnifiedConfig(unifiedConfig config.Config) (*Tracer, error) {

	return Init(unifiedConfig.GetTracerConfig())
}

// ConfigFromEnv loads tracer configuration from environment variables.
func ConfigFromEnv() Config {
	return config.LoadTracerConfigFromEnv()
}

// T returns the global tracer instance.
// If no global tracer has been initialized, it returns a no-op tracer.
func T() trace.Tracer {

	globalMu.RLock()
	tracer := globalTracer
	globalMu.RUnlock()

	if tracer == nil {

		return otel.Tracer("")
	}

	return tracer
}

// Shutdown gracefully shuts down the global tracer. Safe to call multiple times.
func Shutdown(ctx context.Context) error {

	if !globalClosed.CompareAndSwap(false, true) {

		return nil
	}

	globalMu.Lock()
	defer globalMu.Unlock()

	if globalProvider != nil {

		err := globalProvider.Shutdown(ctx)
		globalProvider = nil
		globalTracer = nil

		return err
	}

	return nil
}

/* ---------------------------------------- Tracer Methods -------------------------------------------------- */

// Start creates a new span with the given name.
func (tracer *Tracer) Start(ctx context.Context, name string, options ...StartSpanOption) (context.Context, Span) {

	config := &startSpanConfig{
		kind: SpanKindInternal,
	}

	for _, option := range options {

		option(config)
	}

	spanOptions := []trace.SpanStartOption{
		trace.WithSpanKind(config.kind),
	}

	if len(config.attributes) > 0 {

		spanOptions = append(spanOptions, trace.WithAttributes(config.attributes...))
	}

	if len(config.links) > 0 {

		spanOptions = append(spanOptions, trace.WithLinks(config.links...))
	}

	ctx, span := tracer.tracer.Start(ctx, name, spanOptions...)

	return ctx, wrapSpan(span)
}

// StartServer creates a new server span.
func (tracer *Tracer) StartServer(ctx context.Context, name string, options ...StartSpanOption) (context.Context, Span) {

	options = append([]StartSpanOption{WithSpanKind(SpanKindServer)}, options...)

	return tracer.Start(ctx, name, options...)
}

// StartClient creates a new client span.
func (tracer *Tracer) StartClient(ctx context.Context, name string, options ...StartSpanOption) (context.Context, Span) {

	options = append([]StartSpanOption{WithSpanKind(SpanKindClient)}, options...)

	return tracer.Start(ctx, name, options...)
}

// StartProducer creates a new producer span.
func (tracer *Tracer) StartProducer(ctx context.Context, name string, options ...StartSpanOption) (context.Context, Span) {

	options = append([]StartSpanOption{WithSpanKind(SpanKindProducer)}, options...)

	return tracer.Start(ctx, name, options...)
}

// StartConsumer creates a new consumer span.
func (tracer *Tracer) StartConsumer(ctx context.Context, name string, options ...StartSpanOption) (context.Context, Span) {

	options = append([]StartSpanOption{WithSpanKind(SpanKindConsumer)}, options...)

	return tracer.Start(ctx, name, options...)
}

// Shutdown shuts down the tracer provider. Safe to call multiple times.
func (tracer *Tracer) Shutdown(ctx context.Context) error {

	if !tracer.closed.CompareAndSwap(false, true) {

		return nil
	}

	if tracer.provider != nil {

		return tracer.provider.Shutdown(ctx)
	}

	return nil
}

/* ---------------------------------------- Package-level Functions -------------------------------------------------- */

// Start creates a new span using the global tracer.
func Start(ctx context.Context, name string, options ...StartSpanOption) (context.Context, Span) {

	config := &startSpanConfig{
		kind: SpanKindInternal,
	}

	for _, option := range options {

		option(config)
	}

	spanOptions := []trace.SpanStartOption{
		trace.WithSpanKind(config.kind),
	}

	if len(config.attributes) > 0 {

		spanOptions = append(spanOptions, trace.WithAttributes(config.attributes...))
	}

	if len(config.links) > 0 {

		spanOptions = append(spanOptions, trace.WithLinks(config.links...))
	}

	ctx, span := T().Start(ctx, name, spanOptions...)

	return ctx, wrapSpan(span)
}

// StartServer creates a new server span using the global tracer.
func StartServer(ctx context.Context, name string, options ...StartSpanOption) (context.Context, Span) {

	options = append([]StartSpanOption{WithSpanKind(SpanKindServer)}, options...)

	return Start(ctx, name, options...)
}

// StartClient creates a new client span using the global tracer.
func StartClient(ctx context.Context, name string, options ...StartSpanOption) (context.Context, Span) {

	options = append([]StartSpanOption{WithSpanKind(SpanKindClient)}, options...)

	return Start(ctx, name, options...)
}

// StartProducer creates a new producer span using the global tracer.
func StartProducer(ctx context.Context, name string, options ...StartSpanOption) (context.Context, Span) {

	options = append([]StartSpanOption{WithSpanKind(SpanKindProducer)}, options...)

	return Start(ctx, name, options...)
}

// StartConsumer creates a new consumer span using the global tracer.
func StartConsumer(ctx context.Context, name string, options ...StartSpanOption) (context.Context, Span) {

	options = append([]StartSpanOption{WithSpanKind(SpanKindConsumer)}, options...)

	return Start(ctx, name, options...)
}

/* ---------------------------------------- Helper Functions -------------------------------------------------- */

// createPropagator creates a composite propagator from the list of propagator names.
// Supported propagators:
//   - tracecontext, traceparent: W3C Trace Context
//   - baggage: W3C Baggage
//   - b3: Zipkin B3 (single header)
//   - b3multi: Zipkin B3 (multi header)
//   - jaeger: Jaeger propagation format
func createPropagator(names []string) propagation.TextMapPropagator {
	if len(names) == 0 {
		// Default propagators
		return propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)
	}

	var propagators []propagation.TextMapPropagator
	for _, name := range names {
		switch name {
		case "tracecontext", "traceparent":
			propagators = append(propagators, propagation.TraceContext{})
		case "baggage":
			propagators = append(propagators, propagation.Baggage{})
		case "b3":
			// B3 single header format (used by Zipkin)
			propagators = append(propagators, b3.New(b3.WithInjectEncoding(b3.B3SingleHeader)))
		case "b3multi":
			// B3 multi header format
			propagators = append(propagators, b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader)))
		case "jaeger":
			// Jaeger propagation format
			propagators = append(propagators, jaeger.Jaeger{})
		}
	}

	if len(propagators) == 0 {
		return propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		)
	}

	return propagation.NewCompositeTextMapPropagator(propagators...)
}

/* ---------------------------------------- Convenience Helpers -------------------------------------------------- */

// WithSpan executes a function within a span and handles errors.
func WithSpan(ctx context.Context, name string, fn func(ctx context.Context) error, options ...StartSpanOption) error {

	ctx, span := Start(ctx, name, options...)
	defer span.End()

	err := fn(ctx)

	if err != nil {

		span.SetError(err)

	} else {

		span.SetOK()
	}

	return err
}

// TimeSpan measures the duration of a function within a span.
func TimeSpan(ctx context.Context, name string, fn func(ctx context.Context), options ...StartSpanOption) time.Duration {

	start := time.Now()
	ctx, span := Start(ctx, name, options...)
	defer span.End()

	fn(ctx)
	span.SetOK()

	return time.Since(start)
}
