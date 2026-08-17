package tracer

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// SpanKind represents the type of span.
type SpanKind = trace.SpanKind

// Attribute is an alias for attribute.KeyValue
type Attribute = attribute.KeyValue

// SpanContext wraps trace.SpanContext for convenience
type SpanContext = trace.SpanContext

/* ---------------------------------------- Constants -------------------------------------------------- */

const (
	SpanKindInternal = trace.SpanKindInternal

	SpanKindServer = trace.SpanKindServer

	SpanKindClient = trace.SpanKindClient

	SpanKindProducer = trace.SpanKindProducer

	SpanKindConsumer = trace.SpanKindConsumer
)

/* ---------------------------------------- Attribute Constructors -------------------------------------------------- */

var (
	StringAttr = attribute.String

	IntAttr = attribute.Int

	Int64Attr = attribute.Int64

	Float64Attr = attribute.Float64

	BoolAttr = attribute.Bool

	StringsAttr = attribute.StringSlice

	IntsAttr = attribute.IntSlice

	Int64sAttr = attribute.Int64Slice

	Float64sAttr = attribute.Float64Slice

	BoolsAttr = attribute.BoolSlice
)

/* ---------------------------------------- Span Interface -------------------------------------------------- */

// Span wraps trace.Span with additional convenience methods
type Span interface {
	// End completes the span
	End()

	// SetName sets the span name
	SetName(name string)

	// SetStatus sets the span status
	SetStatus(code codes.Code, description string)

	// SetError marks the span as error with the given error
	SetError(err error)

	// SetOK marks the span as successful
	SetOK()

	// SetAttributes adds attributes to the span
	SetAttributes(attrs ...Attribute)

	// AddEvent adds an event to the span
	AddEvent(name string, opts ...trace.EventOption)

	// RecordError records an error as a span event
	RecordError(err error)

	// SpanContext returns the span context
	SpanContext() SpanContext

	// IsRecording returns true if the span is recording
	IsRecording() bool

	// TracerProvider returns the tracer provider
	TracerProvider() trace.TracerProvider

	// Unwrap returns the underlying trace.Span
	Unwrap() trace.Span
}

/* ---------------------------------------- Span Wrapper Implementation -------------------------------------------------- */

// spanWrapper wraps trace.Span
type spanWrapper struct {
	span trace.Span
}

// End completes the span
func (spanInstance *spanWrapper) End() {

	spanInstance.span.End()
}

// SetName sets the span name
func (spanInstance *spanWrapper) SetName(name string) {

	spanInstance.span.SetName(name)
}

// SetStatus sets the span status
func (spanInstance *spanWrapper) SetStatus(code codes.Code, description string) {

	spanInstance.span.SetStatus(code, description)
}

// SetError marks the span as error
func (spanInstance *spanWrapper) SetError(err error) {

	if err != nil {

		spanInstance.span.SetStatus(codes.Error, err.Error())
		spanInstance.span.RecordError(err)
	}
}

// SetOK marks the span as successful
func (spanInstance *spanWrapper) SetOK() {

	spanInstance.span.SetStatus(codes.Ok, "")
}

// SetAttributes adds attributes to the span
func (spanInstance *spanWrapper) SetAttributes(attrs ...Attribute) {

	spanInstance.span.SetAttributes(attrs...)
}

// AddEvent adds an event to the span
func (spanInstance *spanWrapper) AddEvent(name string, opts ...trace.EventOption) {

	spanInstance.span.AddEvent(name, opts...)
}

// RecordError records an error as a span event
func (spanInstance *spanWrapper) RecordError(err error) {

	spanInstance.span.RecordError(err)
}

// SpanContext returns the span context
func (spanInstance *spanWrapper) SpanContext() SpanContext {

	return spanInstance.span.SpanContext()
}

// IsRecording returns true if the span is recording
func (spanInstance *spanWrapper) IsRecording() bool {

	return spanInstance.span.IsRecording()
}

// TracerProvider returns the tracer provider
func (spanInstance *spanWrapper) TracerProvider() trace.TracerProvider {

	return spanInstance.span.TracerProvider()
}

// Unwrap returns the underlying trace.Span
func (spanInstance *spanWrapper) Unwrap() trace.Span {

	return spanInstance.span
}

// wrapSpan wraps a trace.Span
func wrapSpan(traceSpan trace.Span) Span {
	return &spanWrapper{span: traceSpan}
}

/* ---------------------------------------- Span Options -------------------------------------------------- */

// StartSpanOption configures span creation
type StartSpanOption func(*startSpanConfig)

type startSpanConfig struct {
	kind SpanKind

	attributes []Attribute

	links []trace.Link
}

// WithSpanKind sets the span kind
func WithSpanKind(kind SpanKind) StartSpanOption {
	return func(cfg *startSpanConfig) {
		cfg.kind = kind
	}
}

// WithAttributes adds attributes to the span
func WithAttributes(attrs ...Attribute) StartSpanOption {
	return func(cfg *startSpanConfig) {
		cfg.attributes = append(cfg.attributes, attrs...)
	}
}

// WithLinks adds links to other spans
func WithLinks(links ...trace.Link) StartSpanOption {
	return func(cfg *startSpanConfig) {
		cfg.links = append(cfg.links, links...)
	}
}

/* ---------------------------------------- Context Helper Functions -------------------------------------------------- */

// SpanFromContext returns the span from context
func SpanFromContext(ctx context.Context) Span {

	traceSpan := trace.SpanFromContext(ctx)

	return wrapSpan(traceSpan)
}

// ContextWithSpan returns a new context with the span
func ContextWithSpan(ctx context.Context, spanInstance Span) context.Context {

	if wrapper, ok := spanInstance.(*spanWrapper); ok {

		return trace.ContextWithSpan(ctx, wrapper.span)
	}

	return ctx
}

// TraceIDFromContext extracts trace ID from context
func TraceIDFromContext(ctx context.Context) string {

	spanContext := trace.SpanContextFromContext(ctx)

	if spanContext.HasTraceID() {

		return spanContext.TraceID().String()
	}

	return ""
}

// SpanIDFromContext extracts span ID from context
func SpanIDFromContext(ctx context.Context) string {

	spanContext := trace.SpanContextFromContext(ctx)

	if spanContext.HasSpanID() {

		return spanContext.SpanID().String()
	}

	return ""
}
