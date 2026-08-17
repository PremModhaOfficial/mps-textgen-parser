package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer"
)

// idCounter provides a monotonically increasing fallback counter used by
// fallbackIDGenerator when crypto/rand is unavailable. Combined with
// a nanosecond timestamp, it ensures uniqueness even under degraded entropy.
var idCounter uint64

// TracingMiddleware adds distributed tracing to NATS publish and subscribe operations.
// It supports two modes of operation:
//
//  1. OpenTelemetry mode (UseOTEL=true): Creates OTEL spans via the tracer package,
//     extracts W3C-compatible trace/span IDs from the OTEL span context, and propagates
//     them through NATS message headers. This integrates seamlessly with OTEL collectors,
//     Jaeger, Zipkin, and other W3C Trace Context-compatible backends.
//
//  2. Manual/fallback mode (UseOTEL=false): Generates trace and span IDs using
//     configurable generators (defaulting to crypto/rand), manages parent-child span
//     relationships manually, and propagates context via NATS headers. This mode is
//     useful when OTEL infrastructure is not available.
//
// The middleware implements the Interceptor interface, providing both publish-side
// (producer) and subscribe-side (consumer) instrumentation.
type TracingMiddleware struct {
	// TraceIDGenerator produces unique trace identifiers. In manual mode, this is
	// called to create a new trace when no existing trace context is found. Defaults
	// to a 128-bit (32 hex character) cryptographically random ID, compatible with
	// the W3C Trace Context trace-id format.
	TraceIDGenerator func() string

	// SpanIDGenerator produces unique span identifiers for each operation within a
	// trace. Defaults to a 128-bit cryptographically random ID. In W3C Trace Context,
	// span IDs are typically 64-bit, but the extra length does not break compatibility.
	SpanIDGenerator func() string

	// Sampler is a function that determines whether a new trace should be sampled
	// (i.e., recorded and exported). Returns true to sample. Defaults to always-sample.
	// This is only used in manual mode; in OTEL mode, the OTEL SDK's sampler is used.
	Sampler func() bool

	// UseOTEL controls whether OpenTelemetry integration is active. When true,
	// the middleware creates OTEL spans and relies on the OTEL SDK for trace/span ID
	// generation, sampling decisions, and context propagation.
	UseOTEL bool
}

// NewTracingMiddleware creates a new TracingMiddleware with sensible defaults:
// cryptographically random ID generators, always-sample policy, and OTEL enabled.
// The OTEL integration creates producer/consumer spans that appear in distributed
// traces alongside HTTP, gRPC, and other instrumented operations.
func NewTracingMiddleware() *TracingMiddleware {

	return &TracingMiddleware{
		TraceIDGenerator: defaultIDGenerator,
		SpanIDGenerator:  defaultIDGenerator,
		Sampler:          func() bool { return true },
		UseOTEL:          true,
	}
}

// NewTracingMiddlewareWithOTEL creates a TracingMiddleware with OpenTelemetry explicitly
// enabled. This is functionally identical to NewTracingMiddleware (which also defaults
// to OTEL mode) but makes the intent explicit at the call site for clarity.
func NewTracingMiddlewareWithOTEL() *TracingMiddleware {

	return &TracingMiddleware{
		TraceIDGenerator: defaultIDGenerator,
		SpanIDGenerator:  defaultIDGenerator,
		Sampler:          func() bool { return true },
		UseOTEL:          true,
	}
}

// InterceptPublish returns a PublishMiddleware that instruments outgoing messages
// with tracing context. In OTEL mode, it starts a "producer" span with messaging-specific
// attributes (system, destination, operation) following the OpenTelemetry Messaging
// Semantic Conventions. The OTEL span context (trace ID, span ID, sampling flag) is
// extracted and propagated into the NATS message headers so that downstream consumers
// can continue the trace. In manual mode, it either creates a new root trace or a
// child span under an existing trace, then injects the trace context into headers.
func (tracingMiddleware *TracingMiddleware) InterceptPublish() PublishMiddleware {
	return func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, msg *nats.Msg) error {
			if tracingMiddleware.UseOTEL {
				return tracingMiddleware.publishWithOTEL(ctx, msg, next)
			}
			return tracingMiddleware.publishWithManualTrace(ctx, msg, next)
		}
	}
}

// publishWithOTEL instruments a publish operation using OpenTelemetry spans.
func (tracingMiddleware *TracingMiddleware) publishWithOTEL(ctx context.Context, msg *nats.Msg, next PublishHandler) error {
	ctx, span := tracer.StartProducer(ctx, "events.publish",
		tracer.WithAttributes(
			attribute.String("messaging.system", "nats"),
			attribute.String("messaging.destination", msg.Subject),
			attribute.String("messaging.operation", "publish"),
		),
	)
	defer span.End()

	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		ctx = core.WithTraceContext(ctx, &core.TraceContext{
			TraceID: spanCtx.TraceID().String(),
			SpanID:  spanCtx.SpanID().String(),
			Sampled: spanCtx.IsSampled(),
		})
	}

	ensureHeaders(msg)
	core.ExtractHeaders(ctx, msg.Header)

	err := next(ctx, msg)
	if err != nil {
		span.SetError(err)
	} else {
		span.SetOK()
	}
	return err
}

// publishWithManualTrace instruments a publish using manually generated trace context.
func (tracingMiddleware *TracingMiddleware) publishWithManualTrace(ctx context.Context, msg *nats.Msg, next PublishHandler) error {
	tc, ok := core.TraceContextFromContext(ctx)
	if !ok {
		tc = &core.TraceContext{
			TraceID: tracingMiddleware.TraceIDGenerator(),
			SpanID:  tracingMiddleware.SpanIDGenerator(),
			Sampled: tracingMiddleware.Sampler(),
		}
	} else {
		tc = &core.TraceContext{
			TraceID:  tc.TraceID,
			SpanID:   tracingMiddleware.SpanIDGenerator(),
			ParentID: tc.SpanID,
			Sampled:  tc.Sampled,
			State:    tc.State,
		}
	}
	ctx = core.WithTraceContext(ctx, tc)

	ensureHeaders(msg)
	core.ExtractHeaders(ctx, msg.Header)

	return next(ctx, msg)
}

// ensureHeaders initializes the message header map if nil.
func ensureHeaders(msg *nats.Msg) {
	if msg.Header == nil {
		msg.Header = make(nats.Header)
	}
}

// InterceptSubscribe returns a SubscribeMiddleware that extracts tracing context from
// incoming NATS message headers and continues the distributed trace on the consumer side.
// It first extracts any propagated trace context from message headers (supporting both
// W3C Trace Context and B3 formats via core.InjectContext). In OTEL mode, it then
// starts a "consumer" span linked to the propagated context, enabling end-to-end trace
// visualization across producer and consumer services. In manual mode, it ensures a
// trace context exists (creating one if none was propagated) so that downstream
// handlers always have tracing information available.
func (tracingMiddleware *TracingMiddleware) InterceptSubscribe() SubscribeMiddleware {
	return func(next SubscribeHandler) SubscribeHandler {
		return func(ctx context.Context, msg *nats.Msg) error {
			if msg.Header != nil {
				ctx = core.InjectContext(ctx, msg.Header)
			}

			if tracingMiddleware.UseOTEL {
				return tracingMiddleware.subscribeWithOTEL(ctx, msg, next)
			}
			return tracingMiddleware.subscribeWithManualTrace(ctx, msg, next)
		}
	}
}

// subscribeWithOTEL instruments a subscribe operation using OpenTelemetry spans.
func (tracingMiddleware *TracingMiddleware) subscribeWithOTEL(ctx context.Context, msg *nats.Msg, next SubscribeHandler) error {
	ctx, span := tracer.StartConsumer(ctx, "events.receive",
		tracer.WithAttributes(
			attribute.String("messaging.system", "nats"),
			attribute.String("messaging.destination", msg.Subject),
			attribute.String("messaging.operation", "receive"),
		),
	)
	defer span.End()

	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		ctx = core.WithTraceContext(ctx, &core.TraceContext{
			TraceID: spanCtx.TraceID().String(),
			SpanID:  spanCtx.SpanID().String(),
			Sampled: spanCtx.IsSampled(),
		})
	}

	err := next(ctx, msg)
	if err != nil {
		span.SetError(err)
	} else {
		span.SetOK()
	}
	return err
}

// subscribeWithManualTrace ensures a trace context exists for manual tracing mode.
func (tracingMiddleware *TracingMiddleware) subscribeWithManualTrace(ctx context.Context, msg *nats.Msg, next SubscribeHandler) error {
	if _, ok := core.TraceContextFromContext(ctx); !ok {
		ctx = core.WithTraceContext(ctx, &core.TraceContext{
			TraceID: tracingMiddleware.TraceIDGenerator(),
			SpanID:  tracingMiddleware.SpanIDGenerator(),
			Sampled: tracingMiddleware.Sampler(),
		})
	}
	return next(ctx, msg)
}

// W3C Trace Context support
//
// The W3C Trace Context specification (https://www.w3.org/TR/trace-context/) defines
// a standard for propagating distributed trace identity across service boundaries.
// The traceparent header encodes version, trace-id, parent-id (span-id), and
// trace-flags in a single hyphen-delimited string.

// ExtractW3CTraceParent parses a W3C traceparent header string and returns the
// corresponding TraceContext. The traceparent format is:
//
//	{version}-{trace-id}-{parent-id}-{trace-flags}
//
// Example: "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"
//   - version: "00" (current spec version)
//   - trace-id: 32 hex characters (128-bit)
//   - parent-id: 16 hex characters (64-bit span ID)
//   - trace-flags: "01" means sampled, "00" means not sampled
//
// Returns nil if the header is malformed or too short.
func ExtractW3CTraceParent(traceparent string) *core.TraceContext {
	// Minimum length: 2 (version) + 1 (-) + 32 (trace-id) + 1 (-) + 16 (parent-id) + 1 (-) + 2 (flags) = 55
	if len(traceparent) < 55 {
		return nil
	}

	parts := splitTraceparent(traceparent)
	if len(parts) != 4 {
		return nil
	}

	tc := &core.TraceContext{
		TraceID: parts[1],
		SpanID:  parts[2],
		Sampled: parts[3] == "01",
	}

	return tc
}

// FormatW3CTraceParent serializes a TraceContext into the W3C traceparent header
// format: "00-{trace-id}-{span-id}-{trace-flags}". Returns an empty string if the
// TraceContext is nil or missing required fields. The trace-flags byte is set to
// "01" when the trace is sampled, "00" otherwise.
func FormatW3CTraceParent(tc *core.TraceContext) string {
	if tc == nil || tc.TraceID == "" || tc.SpanID == "" {
		return ""
	}

	flags := "00"
	if tc.Sampled {
		flags = "01"
	}

	return "00-" + tc.TraceID + "-" + tc.SpanID + "-" + flags
}

// splitTraceparent splits a traceparent string on '-' delimiters without allocating
// via strings.Split. This avoids importing the strings package for a single use case.
func splitTraceparent(s string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '-' {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		parts = append(parts, s[start:])
	}
	return parts
}

// defaultIDGenerator produces a 128-bit (32 hex character) cryptographically random
// identifier suitable for use as a trace ID or span ID. It reads from crypto/rand
// for high-quality entropy. If crypto/rand fails (which is extremely rare but possible
// on some constrained systems), it falls back to fallbackIDGenerator.
func defaultIDGenerator() string {

	b := make([]byte, 16)

	if _, err := rand.Read(b); err != nil {
		// Fallback to less random but still unique ID
		return fallbackIDGenerator()
	}

	return hex.EncodeToString(b)
}

// fallbackIDGenerator provides a deterministic-but-unique ID when crypto/rand is
// unavailable. It combines a nanosecond timestamp with an atomically incrementing
// counter, formatted as 32 hex characters. While not cryptographically random,
// it guarantees uniqueness within a single process.
func fallbackIDGenerator() string {

	// Use current time nanoseconds combined with counter as entropy source
	counter := atomic.AddUint64(&idCounter, 1)
	timestamp := time.Now().UnixNano()

	return fmt.Sprintf("%016x%016x", timestamp, counter)
}

// Tracing is a convenience constructor that returns a TracingMiddleware with default
// settings (OTEL enabled, always-sample, crypto/rand IDs). It is the recommended
// entry point for adding distributed tracing to publish and subscribe operations.
func Tracing() *TracingMiddleware {
	return NewTracingMiddleware()
}
