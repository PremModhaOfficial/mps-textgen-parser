package tracer

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// SpanKind represents the type of span in a distributed trace.
// The span kind determines how the span relates to its parent and children.
//
// Span kinds:
//   - Internal: Default, represents internal operations within the service
//   - Server: Represents the server-side of a synchronous RPC call
//   - Client: Represents the client-side of a synchronous RPC call
//   - Producer: Represents the producer of an async message (e.g., message queue)
//   - Consumer: Represents the consumer of an async message
type SpanKind = trace.SpanKind

// Attribute is a key-value pair attached to spans for context.
// Attributes provide additional metadata about the operation being traced,
// enabling filtering and analysis in observability tools.
//
// Use the provided constructors (StringAttr, IntAttr, etc.) to create attributes.
type Attribute = attribute.KeyValue

// SpanContext contains the identifying information about a span.
// This includes the trace ID, span ID, and trace flags.
// SpanContext is immutable and can be safely passed across API boundaries.
type SpanContext = trace.SpanContext

/* ---------------------------------------- Constants -------------------------------------------------- */

// Span kind constants for common operation types.
// These determine how traces are visualized and analyzed in tracing backends.
const (
	// SpanKindInternal represents internal operations within the service.
	// This is the default kind and doesn't indicate a remote call.
	SpanKindInternal = trace.SpanKindInternal

	// SpanKindServer represents the server-side of a synchronous RPC call.
	// Use for incoming HTTP requests, gRPC calls, etc.
	SpanKindServer = trace.SpanKindServer

	// SpanKindClient represents the client-side of a synchronous RPC call.
	// Use for outgoing HTTP requests, database queries, gRPC calls, etc.
	SpanKindClient = trace.SpanKindClient

	// SpanKindProducer represents the producer of an asynchronous message.
	// Use when sending messages to queues (Kafka, RabbitMQ, etc.)
	SpanKindProducer = trace.SpanKindProducer

	// SpanKindConsumer represents the consumer of an asynchronous message.
	// Use when receiving messages from queues.
	SpanKindConsumer = trace.SpanKindConsumer
)

/* ---------------------------------------- Attribute Constructors -------------------------------------------------- */

// Attribute constructors (re-exported from OTEL for convenience).
// These provide strongly-typed attribute creation for span metadata.
// Using these constructors ensures type safety and optimal performance.
//
// Usage:
//
//	span.SetAttributes(
//	    StringAttr("user.id", "123"),
//	    IntAttr("items.count", 42),
//	    BoolAttr("is_premium", true),
//	)
var (
	// StringAttr creates a string attribute
	StringAttr = attribute.String

	// IntAttr creates an int attribute
	IntAttr = attribute.Int

	// Int64Attr creates an int64 attribute
	Int64Attr = attribute.Int64

	// Float64Attr creates a float64 attribute
	Float64Attr = attribute.Float64

	// BoolAttr creates a bool attribute
	BoolAttr = attribute.Bool

	// StringsAttr creates a string slice attribute
	StringsAttr = attribute.StringSlice

	// IntsAttr creates an int slice attribute
	IntsAttr = attribute.IntSlice

	// Int64sAttr creates an int64 slice attribute
	Int64sAttr = attribute.Int64Slice

	// Float64sAttr creates a float64 slice attribute
	Float64sAttr = attribute.Float64Slice

	// BoolsAttr creates a bool slice attribute
	BoolsAttr = attribute.BoolSlice
)

/* ---------------------------------------- Span Interface -------------------------------------------------- */

// Span wraps trace.Span with additional convenience methods.
// This interface provides a simplified API for common span operations
// while maintaining full access to the underlying OTEL span.
//
// Usage:
//
//	ctx, span := tracer.Start(ctx, "operation-name")
//	defer span.End()
//
//	span.SetAttributes(StringAttr("key", "value"))
//
//	if err != nil {
//	    span.SetError(err)
//	    return err
//	}
//
//	span.SetOK()
type Span interface {
	// End completes the span and records its duration.
	// Always call End() to ensure the span is exported.
	End()

	// SetName changes the span name.
	// Use sparingly as the name should typically be set at creation.
	SetName(name string)

	// SetStatus sets the span status code and description.
	// Use SetError() or SetOK() for common cases.
	SetStatus(code codes.Code, description string)

	// SetError marks the span as error and records the error.
	// This sets the status to Error and records the error as an event.
	SetError(err error)

	// SetOK marks the span as successful.
	// Call this when the operation completes successfully.
	SetOK()

	// SetAttributes adds attributes to the span.
	// Attributes provide context about the operation being traced.
	SetAttributes(attrs ...Attribute)

	// AddEvent adds a timestamped event to the span.
	// Events represent significant moments during the span's lifetime.
	AddEvent(name string, opts ...trace.EventOption)

	// RecordError records an error as a span event.
	// Unlike SetError, this doesn't change the span status.
	RecordError(err error)

	// SpanContext returns the span's context.
	// This contains the trace ID, span ID, and trace flags.
	SpanContext() SpanContext

	// IsRecording returns true if the span is recording events.
	// Non-recording spans are created when sampling is disabled.
	IsRecording() bool

	// TracerProvider returns the provider that created this span.
	TracerProvider() trace.TracerProvider

	// Unwrap returns the underlying trace.Span.
	// Use this when you need direct access to the OTEL span.
	Unwrap() trace.Span
}

/* ---------------------------------------- Span Wrapper Implementation -------------------------------------------------- */

// spanWrapper wraps trace.Span to implement our Span interface.
// This provides convenience methods while delegating to the underlying span.
type spanWrapper struct {
	// span is the underlying OTEL span
	span trace.Span
}

// End completes the span and records its duration.
func (spanInstance *spanWrapper) End() {

	spanInstance.span.End()
}

// SetName changes the span name.
func (spanInstance *spanWrapper) SetName(name string) {

	spanInstance.span.SetName(name)
}

// SetStatus sets the span status code and description.
func (spanInstance *spanWrapper) SetStatus(code codes.Code, description string) {

	spanInstance.span.SetStatus(code, description)
}

// SetError marks the span as error and records the error.
// This is a convenience method that:
// 1. Sets the span status to Error with the error message
// 2. Records the error as a span event for detailed analysis
func (spanInstance *spanWrapper) SetError(err error) {

	// Guard against nil error to prevent panic
	if err != nil {

		spanInstance.span.SetStatus(codes.Error, err.Error())
		spanInstance.span.RecordError(err)
	}
}

// SetOK marks the span as successful.
func (spanInstance *spanWrapper) SetOK() {

	spanInstance.span.SetStatus(codes.Ok, "")
}

// SetAttributes adds attributes to the span.
func (spanInstance *spanWrapper) SetAttributes(attrs ...Attribute) {

	spanInstance.span.SetAttributes(attrs...)
}

// AddEvent adds a timestamped event to the span.
func (spanInstance *spanWrapper) AddEvent(name string, opts ...trace.EventOption) {

	spanInstance.span.AddEvent(name, opts...)
}

// RecordError records an error as a span event.
func (spanInstance *spanWrapper) RecordError(err error) {

	spanInstance.span.RecordError(err)
}

// SpanContext returns the span's context.
func (spanInstance *spanWrapper) SpanContext() SpanContext {

	return spanInstance.span.SpanContext()
}

// IsRecording returns true if the span is recording events.
func (spanInstance *spanWrapper) IsRecording() bool {

	return spanInstance.span.IsRecording()
}

// TracerProvider returns the provider that created this span.
func (spanInstance *spanWrapper) TracerProvider() trace.TracerProvider {

	return spanInstance.span.TracerProvider()
}

// Unwrap returns the underlying trace.Span.
func (spanInstance *spanWrapper) Unwrap() trace.Span {

	return spanInstance.span
}

// wrapSpan creates a Span wrapper around a trace.Span.
func wrapSpan(traceSpan trace.Span) Span {

	return &spanWrapper{span: traceSpan}
}

/* ---------------------------------------- Span Options -------------------------------------------------- */

// StartSpanOption configures span creation.
// Use the provided option functions to customize span properties.
type StartSpanOption func(*startSpanConfig)

// startSpanConfig holds span creation options.
type startSpanConfig struct {
	// kind specifies the span kind (internal, server, client, etc.)
	kind SpanKind

	// attributes are key-value pairs added to the span at creation
	attributes []Attribute

	// links connect this span to other spans across trace boundaries
	links []trace.Link
}

// WithSpanKind sets the span kind.
// This determines how the span relates to its parent and children.
//
// Example:
//
//	ctx, span := tracer.Start(ctx, "http.request", WithSpanKind(SpanKindServer))
func WithSpanKind(kind SpanKind) StartSpanOption {

	return func(cfg *startSpanConfig) {
		cfg.kind = kind
	}
}

// WithAttributes adds attributes to the span at creation time.
// These attributes provide context about the operation being traced.
//
// Example:
//
//	ctx, span := tracer.Start(ctx, "db.query",
//	    WithAttributes(StringAttr("db.system", "postgresql")))
func WithAttributes(attrs ...Attribute) StartSpanOption {

	return func(cfg *startSpanConfig) {
		cfg.attributes = append(cfg.attributes, attrs...)
	}
}

// WithLinks adds links to other spans.
// Links connect spans across trace boundaries (e.g., batch processing).
//
// Example:
//
//	link := trace.Link{SpanContext: prevSpan.SpanContext()}
//	ctx, span := tracer.Start(ctx, "batch.process", WithLinks(link))
func WithLinks(links ...trace.Link) StartSpanOption {

	return func(cfg *startSpanConfig) {
		cfg.links = append(cfg.links, links...)
	}
}

/* ---------------------------------------- Context Helper Functions -------------------------------------------------- */

// SpanFromContext returns the span from context.
// If no span exists in the context, returns a no-op span.
//
// Example:
//
//	span := tracer.SpanFromContext(ctx)
//	span.SetAttributes(StringAttr("key", "value"))
func SpanFromContext(ctx context.Context) Span {

	traceSpan := trace.SpanFromContext(ctx)

	return wrapSpan(traceSpan)
}

// ContextWithSpan returns a new context with the span attached.
// This enables propagating spans through function calls.
//
// Example:
//
//	newCtx := tracer.ContextWithSpan(ctx, span)
func ContextWithSpan(ctx context.Context, spanInstance Span) context.Context {

	// Extract the underlying OTEL span from our wrapper
	if wrapper, ok := spanInstance.(*spanWrapper); ok {

		return trace.ContextWithSpan(ctx, wrapper.span)
	}

	return ctx
}

// TraceIDFromContext extracts trace ID from context.
// Returns empty string if no valid trace ID exists.
//
// Example:
//
//	traceID := tracer.TraceIDFromContext(ctx)
//	if traceID != "" {
//	    log.Printf("trace: %s", traceID)
//	}
func TraceIDFromContext(ctx context.Context) string {

	spanContext := trace.SpanContextFromContext(ctx)

	if spanContext.HasTraceID() {

		return spanContext.TraceID().String()
	}

	return ""
}

// SpanIDFromContext extracts span ID from context.
// Returns empty string if no valid span ID exists.
//
// Example:
//
//	spanID := tracer.SpanIDFromContext(ctx)
func SpanIDFromContext(ctx context.Context) string {

	spanContext := trace.SpanContextFromContext(ctx)

	if spanContext.HasSpanID() {

		return spanContext.SpanID().String()
	}

	return ""
}
