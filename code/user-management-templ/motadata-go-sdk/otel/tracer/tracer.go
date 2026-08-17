// Package tracer provides distributed tracing capabilities with OpenTelemetry integration.
//
// This package wraps the OpenTelemetry SDK to provide simplified distributed tracing with:
//   - Automatic trace context propagation
//   - Span creation with convenience methods
//   - Multiple propagator support (W3C, B3, Jaeger)
//   - OTEL collector integration with batching
//   - Configurable sampling strategies
//
// Basic Usage:
//
//	// Initialize the global tracer
//	tracer.Init(tracer.DefaultConfig())
//	defer tracer.Shutdown(context.Background())
//
//	// Create spans
//	ctx, span := tracer.Start(ctx, "operation-name")
//	defer span.End()
//
//	span.SetAttributes(tracer.StringAttr("key", "value"))
//
// Error Handling:
//
//	ctx, span := tracer.Start(ctx, "db.query")
//	defer span.End()
//
//	result, err := db.Query(ctx, query)
//	if err != nil {
//	    span.SetError(err)
//	    return err
//	}
//	span.SetOK()
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
// This allows existing code using tracer.Config to continue working.
type Config = config.TracerConfig

// Tracer wraps the OTEL tracer with convenience methods.
// It provides a simplified API for distributed tracing while managing
// the underlying OTEL provider lifecycle.
//
// Thread Safety: Tracer is safe for concurrent use. Multiple goroutines
// may call methods on the same Tracer simultaneously.
type Tracer struct {
	// tracer is the underlying OTEL tracer that creates spans
	tracer trace.Tracer

	// provider is the OTEL tracer provider for lifecycle management
	provider *sdktrace.TracerProvider

	// closed tracks whether the tracer has been closed
	// Using atomic bool for thread-safe idempotent shutdown
	closed atomic.Bool
}

/* ---------------------------------------- Global State -------------------------------------------------- */

// Global tracer state.
// The global tracer provides a convenient default for applications that don't need
// multiple tracer instances. It's protected by a RWMutex for thread-safe access.
var (
	// globalMu protects access to globalProvider and globalTracer
	globalMu sync.RWMutex

	// globalProvider is the OTEL tracer provider
	globalProvider *sdktrace.TracerProvider

	// globalTracer is the default tracer instance
	globalTracer trace.Tracer

	// globalClosed tracks whether the global tracer has been shut down
	globalClosed atomic.Bool
)

/* ---------------------------------------- Constructor Functions -------------------------------------------------- */

// DefaultConfig returns the default tracer configuration with service defaults.
// This provides sensible defaults for development and can be customized as needed.
//
// Default values:
//   - ServiceName: "app"
//   - ServiceVersion: "0.0.0"
//   - Environment: "development"
//   - Enabled: true
//   - SamplingRatio: 1.0 (sample all traces)
//   - Propagators: ["tracecontext", "baggage"]
//
// Returns:
//   - Config: A tracer configuration with defaults applied
func DefaultConfig() Config {

	// Start with the base default config from the config package
	defaultConfig := config.DefaultTracerConfig()

	// Apply service identification defaults
	defaultConfig.ServiceName = "app"
	defaultConfig.ServiceVersion = "0.0.0"
	defaultConfig.Environment = "development"

	return defaultConfig
}

/* ---------------------------------------- Global Initialization Functions -------------------------------------------------- */

// Init initializes the global tracer with the given configuration.
// This should be called once at application startup.
//
// When tracing is disabled (config.Enabled=false), returns nil without error.
// This allows conditional tracing without code changes.
//
// Parameters:
//   - config: The tracer configuration
//
// Returns:
//   - *Tracer: The initialized tracer (nil if disabled)
//   - error: Non-nil if initialization fails
func Init(config Config) (*Tracer, error) {

	// Early return if tracing is disabled
	// This allows conditional tracing via configuration
	if !config.Enabled {

		return nil, nil
	}

	// Reset the closed flag for re-initialization
	globalClosed.Store(false)

	// Create the OTEL tracer provider with configured exporters and samplers
	provider, err := newTracerProvider(config)

	if err != nil {

		return nil, err
	}

	// Set as the global OTEL provider
	// This enables context propagation across the application
	otel.SetTracerProvider(provider)

	// Configure text map propagator for trace context propagation
	// This determines how trace context is encoded/decoded in headers
	propagator := createPropagator(config.Propagators)
	otel.SetTextMapPropagator(propagator)

	// Create a named tracer for the service
	tracer := provider.Tracer(config.ServiceName)

	// Update global state with write lock
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
// This is convenient for containerized deployments where configuration
// is passed via environment.
//
// Returns:
//   - *Tracer: The initialized tracer
//   - error: Non-nil if initialization fails
func InitFromEnv() (*Tracer, error) {

	config := ConfigFromEnv()

	return Init(config)
}

// InitFromUnifiedConfig initializes the global tracer from the unified config.
// This integrates with the SDK's unified configuration system.
//
// Parameters:
//   - unifiedConfig: The unified application configuration
//
// Returns:
//   - *Tracer: The initialized tracer
//   - error: Non-nil if initialization fails
func InitFromUnifiedConfig(unifiedConfig config.Config) (*Tracer, error) {

	return Init(unifiedConfig.GetTracerConfig())
}

// ConfigFromEnv loads tracer configuration from environment variables.
//
// Returns:
//   - Config: Configuration from environment
func ConfigFromEnv() Config {

	return config.LoadTracerConfigFromEnv()
}

// T returns the global tracer instance.
// If no global tracer has been initialized, it returns a no-op tracer.
//
// Returns:
//   - trace.Tracer: The global tracer
func T() trace.Tracer {

	globalMu.RLock()
	tracer := globalTracer
	globalMu.RUnlock()

	// Return no-op tracer if not initialized
	// This ensures the application doesn't crash if tracing isn't configured
	if tracer == nil {

		return otel.Tracer("")
	}

	return tracer
}

// Shutdown gracefully shuts down the global tracer.
// This method is idempotent - multiple calls are safe.
//
// It flushes any buffered spans and releases resources.
// Always call this before application exit to ensure all traces are exported.
//
// Parameters:
//   - ctx: Context with timeout for shutdown
//
// Returns:
//   - error: Non-nil if shutdown encounters errors
func Shutdown(ctx context.Context) error {

	// Ensure Shutdown is only executed once using atomic compare-and-swap
	if !globalClosed.CompareAndSwap(false, true) {

		return nil
	}

	globalMu.Lock()
	defer globalMu.Unlock()

	// Shutdown the provider and clear global state
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
// The span must be ended by calling span.End().
//
// Parameters:
//   - ctx: Parent context (may contain parent span)
//   - name: The span name (should describe the operation)
//   - options: Optional span configuration
//
// Returns:
//   - context.Context: New context containing the span
//   - Span: The created span (must call End())
//
// Example:
//
//	ctx, span := tracer.Start(ctx, "db.query")
//	defer span.End()
func (tracer *Tracer) Start(ctx context.Context, name string, options ...StartSpanOption) (context.Context, Span) {

	// Apply configuration options
	config := &startSpanConfig{
		kind: SpanKindInternal, // Default to internal span
	}

	for _, option := range options {

		option(config)
	}

	// Build OTEL span options from our configuration
	spanOptions := []trace.SpanStartOption{
		trace.WithSpanKind(config.kind),
	}

	// Add attributes if provided
	if len(config.attributes) > 0 {

		spanOptions = append(spanOptions, trace.WithAttributes(config.attributes...))
	}

	// Add links if provided
	if len(config.links) > 0 {

		spanOptions = append(spanOptions, trace.WithLinks(config.links...))
	}

	// Create the span
	ctx, span := tracer.tracer.Start(ctx, name, spanOptions...)

	return ctx, wrapSpan(span)
}

// StartServer creates a new server span.
// Use this for incoming requests (HTTP handlers, gRPC servers, etc.)
//
// Parameters:
//   - ctx: Parent context
//   - name: The span name
//   - options: Optional span configuration
//
// Returns:
//   - context.Context: New context containing the span
//   - Span: The created span
func (tracer *Tracer) StartServer(ctx context.Context, name string, options ...StartSpanOption) (context.Context, Span) {

	// Prepend server kind option (can be overridden by user options)
	options = append([]StartSpanOption{WithSpanKind(SpanKindServer)}, options...)

	return tracer.Start(ctx, name, options...)
}

// StartClient creates a new client span.
// Use this for outgoing requests (HTTP clients, database queries, etc.)
//
// Parameters:
//   - ctx: Parent context
//   - name: The span name
//   - options: Optional span configuration
//
// Returns:
//   - context.Context: New context containing the span
//   - Span: The created span
func (tracer *Tracer) StartClient(ctx context.Context, name string, options ...StartSpanOption) (context.Context, Span) {

	options = append([]StartSpanOption{WithSpanKind(SpanKindClient)}, options...)

	return tracer.Start(ctx, name, options...)
}

// StartProducer creates a new producer span.
// Use this when sending messages to queues or event systems.
//
// Parameters:
//   - ctx: Parent context
//   - name: The span name
//   - options: Optional span configuration
//
// Returns:
//   - context.Context: New context containing the span
//   - Span: The created span
func (tracer *Tracer) StartProducer(ctx context.Context, name string, options ...StartSpanOption) (context.Context, Span) {

	options = append([]StartSpanOption{WithSpanKind(SpanKindProducer)}, options...)

	return tracer.Start(ctx, name, options...)
}

// StartConsumer creates a new consumer span.
// Use this when receiving messages from queues or event systems.
//
// Parameters:
//   - ctx: Parent context
//   - name: The span name
//   - options: Optional span configuration
//
// Returns:
//   - context.Context: New context containing the span
//   - Span: The created span
func (tracer *Tracer) StartConsumer(ctx context.Context, name string, options ...StartSpanOption) (context.Context, Span) {

	options = append([]StartSpanOption{WithSpanKind(SpanKindConsumer)}, options...)

	return tracer.Start(ctx, name, options...)
}

// Shutdown shuts down the tracer provider.
// This method is idempotent - multiple calls are safe.
//
// Parameters:
//   - ctx: Context with timeout for shutdown
//
// Returns:
//   - error: Non-nil if shutdown encounters errors
func (tracer *Tracer) Shutdown(ctx context.Context) error {

	// Ensure Shutdown is only executed once
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
// This is a convenience function for applications using the global tracer.
//
// Parameters:
//   - ctx: Parent context
//   - name: The span name
//   - options: Optional span configuration
//
// Returns:
//   - context.Context: New context containing the span
//   - Span: The created span
func Start(ctx context.Context, name string, options ...StartSpanOption) (context.Context, Span) {

	// Apply configuration options
	config := &startSpanConfig{
		kind: SpanKindInternal,
	}

	for _, option := range options {

		option(config)
	}

	// Build OTEL span options
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
// Propagators determine how trace context is encoded in headers for cross-service propagation.
//
// Supported propagators:
//   - "tracecontext", "traceparent": W3C Trace Context (recommended)
//   - "baggage": W3C Baggage for propagating arbitrary key-value pairs
//   - "b3": Zipkin B3 single-header format
//   - "b3multi": Zipkin B3 multi-header format
//   - "jaeger": Jaeger native propagation format
//
// Parameters:
//   - names: List of propagator names to enable
//
// Returns:
//   - propagation.TextMapPropagator: The composite propagator
func createPropagator(names []string) propagation.TextMapPropagator {

	// Default to W3C standard propagators if none specified
	if len(names) == 0 {

		return propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, // W3C Trace Context
			propagation.Baggage{},      // W3C Baggage
		)
	}

	var propagators []propagation.TextMapPropagator

	for _, name := range names {

		switch name {

		case "tracecontext", "traceparent":
			// W3C Trace Context - the standard for trace propagation
			propagators = append(propagators, propagation.TraceContext{})

		case "baggage":
			// W3C Baggage - for propagating arbitrary key-value pairs
			propagators = append(propagators, propagation.Baggage{})

		case "b3":
			// Zipkin B3 single header format
			// Used by Zipkin and some other tracing systems
			propagators = append(propagators, b3.New(b3.WithInjectEncoding(b3.B3SingleHeader)))

		case "b3multi":
			// Zipkin B3 multi header format
			propagators = append(propagators, b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader)))

		case "jaeger":
			// Jaeger native propagation format
			propagators = append(propagators, jaeger.Jaeger{})
		}
	}

	// Fallback to default if no valid propagators were found
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
// This is a convenience function for the common pattern of creating a span,
// executing work, and setting the status based on the result.
//
// Parameters:
//   - ctx: Parent context
//   - name: The span name
//   - fn: The function to execute within the span
//   - options: Optional span configuration
//
// Returns:
//   - error: The error returned by fn (if any)
//
// Example:
//
//	err := tracer.WithSpan(ctx, "process.order", func(ctx context.Context) error {
//	    return processOrder(ctx, order)
//	})
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
// This is useful for timing operations and recording the duration in traces.
//
// Parameters:
//   - ctx: Parent context
//   - name: The span name
//   - fn: The function to execute and time
//   - options: Optional span configuration
//
// Returns:
//   - time.Duration: The elapsed time
//
// Example:
//
//	duration := tracer.TimeSpan(ctx, "data.process", func(ctx context.Context) {
//	    processData(ctx, data)
//	})
//	log.Printf("Processing took %v", duration)
func TimeSpan(ctx context.Context, name string, fn func(ctx context.Context), options ...StartSpanOption) time.Duration {

	start := time.Now()
	ctx, span := Start(ctx, name, options...)
	defer span.End()

	fn(ctx)
	span.SetOK()

	return time.Since(start)
}
