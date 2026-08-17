package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// idCounter for fallback ID generation
var idCounter uint64

// TracingMiddleware adds distributed tracing to messages
type TracingMiddleware struct {
	// TraceIDGenerator generates trace IDs
	TraceIDGenerator func() string

	// SpanIDGenerator generates span IDs
	SpanIDGenerator func() string

	// Sampler determines if a trace should be sampled
	Sampler func() bool

	// UseOTEL enables OpenTelemetry integration
	UseOTEL bool
}

// NewTracingMiddleware creates a new tracing middleware
func NewTracingMiddleware() *TracingMiddleware {

	return &TracingMiddleware{
		TraceIDGenerator: defaultIDGenerator,
		SpanIDGenerator:  defaultIDGenerator,
		Sampler:          func() bool { return true },
		UseOTEL:          true,
	}
}

// NewTracingMiddlewareWithOTEL creates a tracing middleware with OpenTelemetry enabled
func NewTracingMiddlewareWithOTEL() *TracingMiddleware {

	return &TracingMiddleware{
		TraceIDGenerator: defaultIDGenerator,
		SpanIDGenerator:  defaultIDGenerator,
		Sampler:          func() bool { return true },
		UseOTEL:          true,
	}
}

// InterceptPublish returns the publish middleware
func (tracingMiddleware *TracingMiddleware) InterceptPublish() PublishMiddleware {

	return func(next PublishHandler) PublishHandler {

		return func(ctx context.Context, subject string, msg *core.Message) error {

			// Use OpenTelemetry if enabled
			if tracingMiddleware.UseOTEL {

				ctx, span := tracer.StartProducer(ctx, "events.publish",
					tracer.WithAttributes(
						attribute.String("messaging.system", "nats"),
						attribute.String("messaging.destination", subject),
						attribute.String("messaging.operation", "publish"),
					),
				)
				defer span.End()

				// Extract trace and span IDs from OpenTelemetry span
				spanCtx := trace.SpanContextFromContext(ctx)

				if spanCtx.IsValid() {

					tc := &core.TraceContext{
						TraceID: spanCtx.TraceID().String(),
						SpanID:  spanCtx.SpanID().String(),
						Sampled: spanCtx.IsSampled(),
					}

					ctx = core.WithTraceContext(ctx, tc)
				}

				// Inject trace context into headers
				if msg.Headers == nil {

					msg.Headers = make(core.Headers)
				}

				core.ExtractHeaders(ctx, msg.Headers)

				err := next(ctx, subject, msg)

				if err != nil {

					span.SetError(err)

				} else {

					span.SetOK()
				}

				return err
			}

			// Fallback to manual trace context
			tc, ok := core.TraceContextFromContext(ctx)

			if !ok {

				tc = &core.TraceContext{
					TraceID: tracingMiddleware.TraceIDGenerator(),
					SpanID:  tracingMiddleware.SpanIDGenerator(),
					Sampled: tracingMiddleware.Sampler(),
				}

				ctx = core.WithTraceContext(ctx, tc)

			} else {
				// Create new span under existing trace
				tc = &core.TraceContext{
					TraceID:  tc.TraceID,
					SpanID:   tracingMiddleware.SpanIDGenerator(),
					ParentID: tc.SpanID,
					Sampled:  tc.Sampled,
					State:    tc.State,
				}

				ctx = core.WithTraceContext(ctx, tc)
			}

			// Inject trace context into headers
			if msg.Headers == nil {

				msg.Headers = make(core.Headers)
			}

			core.ExtractHeaders(ctx, msg.Headers)

			return next(ctx, subject, msg)
		}
	}
}

// InterceptSubscribe returns the subscribe middleware
func (tracingMiddleware *TracingMiddleware) InterceptSubscribe() SubscribeMiddleware {

	return func(next SubscribeHandler) SubscribeHandler {

		return func(ctx context.Context, msg *core.Message) error {

			// Extract trace context from headers first
			if msg.Headers != nil {

				ctx = core.InjectContext(ctx, msg.Headers)
			}

			// Use OpenTelemetry if enabled
			if tracingMiddleware.UseOTEL {

				ctx, span := tracer.StartConsumer(ctx, "events.receive",
					tracer.WithAttributes(
						attribute.String("messaging.system", "nats"),
						attribute.String("messaging.destination", msg.Subject),
						attribute.String("messaging.operation", "receive"),
					),
				)
				defer span.End()

				// Update trace context from OpenTelemetry span
				spanCtx := trace.SpanContextFromContext(ctx)

				if spanCtx.IsValid() {

					tc := &core.TraceContext{
						TraceID: spanCtx.TraceID().String(),
						SpanID:  spanCtx.SpanID().String(),
						Sampled: spanCtx.IsSampled(),
					}

					ctx = core.WithTraceContext(ctx, tc)
				}

				err := next(ctx, msg)

				if err != nil {

					span.SetError(err)

				} else {

					span.SetOK()
				}

				return err
			}

			// Fallback: If no trace context, create one
			if _, ok := core.TraceContextFromContext(ctx); !ok {

				tc := &core.TraceContext{
					TraceID: tracingMiddleware.TraceIDGenerator(),
					SpanID:  tracingMiddleware.SpanIDGenerator(),
					Sampled: tracingMiddleware.Sampler(),
				}

				ctx = core.WithTraceContext(ctx, tc)
			}

			return next(ctx, msg)
		}
	}
}

// W3C Trace Context support

// ExtractW3CTraceParent extracts trace context from W3C traceparent header
func ExtractW3CTraceParent(traceparent string) *core.TraceContext {
	// Format: version-trace_id-parent_id-trace_flags
	// Example: 00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01
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

// FormatW3CTraceParent formats trace context as W3C traceparent
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

// defaultIDGenerator generates cryptographically random trace/span IDs
func defaultIDGenerator() string {

	b := make([]byte, 16)

	if _, err := rand.Read(b); err != nil {
		// Fallback to less random but still unique ID
		return fallbackIDGenerator()
	}

	return hex.EncodeToString(b)
}

// fallbackIDGenerator provides a fallback when crypto/rand fails
func fallbackIDGenerator() string {

	// Use current time nanoseconds combined with counter as entropy source
	counter := atomic.AddUint64(&idCounter, 1)
	timestamp := time.Now().UnixNano()

	return fmt.Sprintf("%016x%016x", timestamp, counter)
}

// Tracing returns a configured tracing middleware
func Tracing() *TracingMiddleware {
	return NewTracingMiddleware()
}
