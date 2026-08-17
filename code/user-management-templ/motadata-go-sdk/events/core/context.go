package core

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer"
)

// Standard header keys used to propagate metadata through NATS message headers.
// These constants ensure consistent key naming across publishers, subscribers,
// and middleware. Two tracing formats are supported: W3C Trace Context
// (HeaderTraceParent, HeaderTraceState) and Zipkin B3 (HeaderB3*), allowing
// interoperability with different distributed tracing backends.
const (
	// HeaderContentType indicates the MIME type of the message payload (e.g., "application/json").
	HeaderContentType = "Content-Type"

	// HeaderMessageID is a unique message identifier used by NATS JetStream
	// for exactly-once delivery via server-side deduplication.
	HeaderMessageID = "Nats-Msg-Id"

	// HeaderCorrelationID links related messages in a workflow or saga,
	// enabling end-to-end request tracing across multiple services.
	HeaderCorrelationID = "Correlation-ID"

	// HeaderTenantID carries the multi-tenant identifier so that downstream
	// handlers can scope operations (database queries, cache keys) to the
	// correct tenant without additional lookups.
	HeaderTenantID = "X-Tenant-ID"

	// HeaderTraceID is a custom header carrying the distributed trace identifier.
	// Used alongside HeaderSpanID for lightweight tracing when W3C or B3 headers
	// are not present.
	HeaderTraceID = "X-Trace-ID"

	// HeaderSpanID carries the current span identifier within a trace.
	HeaderSpanID = "X-Span-ID"

	// HeaderTraceParent follows the W3C Trace Context specification
	// (format: "00-<traceId>-<spanId>-<flags>") for cross-service trace propagation.
	HeaderTraceParent = "traceparent"

	// HeaderTraceState carries vendor-specific trace data as defined by the
	// W3C Trace Context specification (e.g., "congo=t61rcWkgMzE").
	HeaderTraceState = "tracestate"

	// HeaderB3TraceID is the Zipkin B3 trace identifier, supported for
	// compatibility with systems that use B3 propagation.
	HeaderB3TraceID = "X-B3-TraceId"

	// HeaderB3SpanID is the Zipkin B3 span identifier.
	HeaderB3SpanID = "X-B3-SpanId"

	// HeaderB3ParentSpanID is the Zipkin B3 parent span identifier, used to
	// reconstruct the call tree in B3-based tracing systems.
	HeaderB3ParentSpanID = "X-B3-ParentSpanId"

	// HeaderB3Sampled indicates whether the trace is sampled in B3 format
	// ("1" or "true" means sampled).
	HeaderB3Sampled = "X-B3-Sampled"
)

// TraceContext holds distributed tracing identifiers that travel across service
// boundaries via NATS message headers. It supports both W3C Trace Context and
// Zipkin B3 formats. When an OpenTelemetry span is active, ExtractOTELTraceContext
// can populate this struct from the OTEL span; otherwise it is built from
// incoming message headers by InjectContext.
type TraceContext struct {
	// TraceID is the globally unique identifier for the entire distributed trace.
	TraceID string

	// SpanID identifies the current unit of work within the trace.
	SpanID string

	// ParentID identifies the parent span that initiated this span, used to
	// reconstruct the causal call tree.
	ParentID string

	// Sampled indicates whether this trace should be recorded by the tracing
	// backend. When false, backends may drop the trace to reduce storage costs.
	Sampled bool

	// State carries vendor-specific trace data (W3C tracestate header value).
	State string
}

// Metadata is a convenience aggregate that bundles all per-message context
// values into a single struct. It can be stored in a Go context via
// WithMetadata and retrieved via MetadataFromContext. The Custom map allows
// application-specific key-value pairs to travel alongside the standard fields.
type Metadata struct {
	// TenantID identifies the tenant that owns or originated this message.
	TenantID string

	// CorrelationID links this message to a broader workflow or request chain.
	CorrelationID string

	// MessageID is a unique identifier for this specific message, used by
	// JetStream for deduplication.
	MessageID string

	// Timestamp records when the message was created or published.
	Timestamp time.Time

	// Custom holds arbitrary application-defined metadata key-value pairs
	// that do not fit into the standard fields above.
	Custom map[string]string
}

// contextKey is an unexported type used as the key for context.WithValue to
// prevent collisions with keys defined in other packages.
type contextKey int

// context key constants identify the event-related values stored in a Go context.
// Each key corresponds to a With*/FromContext function pair below.
const (
	// tenantIDKey stores the multi-tenant identifier; used by the tenant package
	// for routing and by middleware for scoped metrics/logging.
	tenantIDKey contextKey = iota

	// traceContextKey stores a *TraceContext; used by the tracing middleware
	// and ExtractHeaders to propagate distributed trace IDs across NATS messages.
	traceContextKey

	// messageIDKey stores the unique message identifier; used by JetStream
	// publishers for deduplication.
	messageIDKey

	// correlationIDKey stores the correlation identifier; used by middleware
	// and handlers to link related messages in a workflow.
	correlationIDKey

	// metadataKey stores a *Metadata aggregate; provides a single-lookup
	// alternative to fetching each field individually.
	metadataKey
)

// WithTenantID returns a new context carrying the given tenant identifier.
// The tenant ID is later extracted by ExtractHeaders to populate the
// X-Tenant-ID NATS header, and by the tenant package for routing decisions.
func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantIDKey, tenantID)
}

// TenantIDFromContext retrieves the tenant identifier previously stored by
// WithTenantID. The second return value is false when no tenant ID is present,
// allowing callers to distinguish between an empty string and a missing value.
func TenantIDFromContext(ctx context.Context) (string, bool) {
	tenantIDValue, exists := ctx.Value(tenantIDKey).(string)
	return tenantIDValue, exists
}

// WithTraceContext returns a new context carrying the given TraceContext.
// This is used by InjectContext when processing inbound messages and by
// application code that needs to manually propagate trace information.
func WithTraceContext(ctx context.Context, traceContext *TraceContext) context.Context {
	return context.WithValue(ctx, traceContextKey, traceContext)
}

// TraceContextFromContext retrieves the TraceContext previously stored by
// WithTraceContext. Returns nil and false when no trace context is present.
// Used by ExtractHeaders as a fallback when no active OTEL span exists.
func TraceContextFromContext(ctx context.Context) (*TraceContext, bool) {
	traceContextValue, exists := ctx.Value(traceContextKey).(*TraceContext)
	return traceContextValue, exists
}

// WithMessageID returns a new context carrying the given message identifier.
// The message ID is propagated into the Nats-Msg-Id header by ExtractHeaders,
// enabling JetStream server-side deduplication for exactly-once publishing.
func WithMessageID(ctx context.Context, messageID string) context.Context {
	return context.WithValue(ctx, messageIDKey, messageID)
}

// MessageIDFromContext retrieves the message identifier previously stored by
// WithMessageID. Returns an empty string and false when absent.
func MessageIDFromContext(ctx context.Context) (string, bool) {
	messageIDValue, exists := ctx.Value(messageIDKey).(string)
	return messageIDValue, exists
}

// WithCorrelationID returns a new context carrying the given correlation identifier.
// Correlation IDs tie together all messages belonging to a single user request or
// workflow, making it possible to trace a complete operation across services in logs
// and monitoring dashboards.
func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, correlationIDKey, correlationID)
}

// CorrelationIDFromContext retrieves the correlation identifier previously stored
// by WithCorrelationID. Returns an empty string and false when absent.
func CorrelationIDFromContext(ctx context.Context) (string, bool) {
	correlationIDValue, exists := ctx.Value(correlationIDKey).(string)
	return correlationIDValue, exists
}

// WithMetadata returns a new context carrying the given Metadata aggregate.
// This is a convenience for storing all per-message context in a single value
// rather than calling WithTenantID, WithCorrelationID, and WithMessageID
// individually. Note that ExtractHeaders reads each field separately, so
// storing Metadata alone is not sufficient for header propagation -- the
// individual With* functions must also be called, or ExtractHeaders must be
// extended to read from Metadata.
func WithMetadata(ctx context.Context, metadata *Metadata) context.Context {
	return context.WithValue(ctx, metadataKey, metadata)
}

// MetadataFromContext retrieves the Metadata aggregate previously stored by
// WithMetadata. Returns nil and false when no metadata is present.
func MetadataFromContext(ctx context.Context) (*Metadata, bool) {
	metadataValue, exists := ctx.Value(metadataKey).(*Metadata)
	return metadataValue, exists
}

// extractOTELTraceHeaders writes OTEL-based trace headers when an active span exists.
// Returns true if OTEL trace info was found and written.
func extractOTELTraceHeaders(ctx context.Context, headers nats.Header) bool {
	traceID := tracer.TraceIDFromContext(ctx)
	spanID := tracer.SpanIDFromContext(ctx)

	if traceID == "" || spanID == "" {
		return false
	}

	headers.Set(HeaderTraceID, traceID)
	headers.Set(HeaderB3TraceID, traceID)
	headers.Set(HeaderSpanID, spanID)
	headers.Set(HeaderB3SpanID, spanID)
	headers.Set(HeaderB3Sampled, "1")
	headers.Set(HeaderTraceParent, "00-"+traceID+"-"+spanID+"-01")

	return true
}

// extractLocalTraceHeaders writes trace headers from a manually stored TraceContext.
func extractLocalTraceHeaders(ctx context.Context, headers nats.Header) {
	traceContext, exists := TraceContextFromContext(ctx)
	if !exists {
		return
	}

	if traceContext.TraceID != "" {
		headers.Set(HeaderTraceID, traceContext.TraceID)
		headers.Set(HeaderB3TraceID, traceContext.TraceID)
	}

	if traceContext.SpanID != "" {
		headers.Set(HeaderSpanID, traceContext.SpanID)
		headers.Set(HeaderB3SpanID, traceContext.SpanID)
	}

	if traceContext.Sampled {
		headers.Set(HeaderB3Sampled, "1")
	}

	if traceContext.TraceID != "" && traceContext.SpanID != "" {
		sampledFlag := "00"
		if traceContext.Sampled {
			sampledFlag = "01"
		}
		headers.Set(HeaderTraceParent, "00-"+traceContext.TraceID+"-"+traceContext.SpanID+"-"+sampledFlag)
	}

	if traceContext.State != "" {
		headers.Set(HeaderTraceState, traceContext.State)
	}
}

// ExtractHeaders reads tenant ID, trace context, correlation ID, and message ID
// from the Go context and writes them into NATS message headers. This is the
// publish-side counterpart to InjectContext.
//
// Trace propagation strategy:
//  1. Prefer an active OpenTelemetry span (via the otel/tracer package).
//  2. Fall back to a manually stored TraceContext from the context.
//
// Both W3C Trace Context (traceparent) and Zipkin B3 headers are emitted so
// that heterogeneous tracing backends can all participate in the same trace.
//
// If headers is nil a new nats.Header map is allocated and returned.
func ExtractHeaders(ctx context.Context, headers nats.Header) nats.Header {
	if headers == nil {
		headers = make(nats.Header)
	}

	if tenantID, exists := TenantIDFromContext(ctx); exists {
		headers.Set(HeaderTenantID, tenantID)
	}

	// Prefer OTEL span context, fall back to local trace context
	if !extractOTELTraceHeaders(ctx, headers) {
		extractLocalTraceHeaders(ctx, headers)
	}

	if correlationID, exists := CorrelationIDFromContext(ctx); exists {
		headers.Set(HeaderCorrelationID, correlationID)
	}

	if messageID, exists := MessageIDFromContext(ctx); exists {
		headers.Set(HeaderMessageID, messageID)
	}

	return headers
}

// ExtractOTELTraceContext reads the trace and span IDs from the active
// OpenTelemetry span in the context (via the otel/tracer package) and returns
// a TraceContext suitable for manual propagation. Returns nil when no active
// OTEL span is found, indicating that the caller should fall back to other
// trace propagation mechanisms.
func ExtractOTELTraceContext(ctx context.Context) *TraceContext {
	traceID := tracer.TraceIDFromContext(ctx)
	spanID := tracer.SpanIDFromContext(ctx)

	if traceID == "" || spanID == "" {
		return nil
	}

	return &TraceContext{
		TraceID: traceID,
		SpanID:  spanID,
		Sampled: true, // OTEL spans are always sampled when active
	}
}

// InjectContext is the subscribe-side counterpart to ExtractHeaders. It reads
// tenant ID, trace context, correlation ID, and message ID from incoming NATS
// message headers and stores them in the Go context so that downstream
// handlers and middleware can access them via the *FromContext functions.
//
// For trace context, both the custom X-Trace-ID/X-Span-ID headers and the
// Zipkin B3 headers are checked, with the custom headers taking precedence.
// If headers is nil the original context is returned unchanged.
func InjectContext(ctx context.Context, headers nats.Header) context.Context {
	if headers == nil {
		return ctx
	}

	if tenantID := headers.Get(HeaderTenantID); tenantID != "" {
		ctx = WithTenantID(ctx, tenantID)
	}

	ctx = injectTraceContext(ctx, headers)

	if correlationID := headers.Get(HeaderCorrelationID); correlationID != "" {
		ctx = WithCorrelationID(ctx, correlationID)
	}

	if messageID := headers.Get(HeaderMessageID); messageID != "" {
		ctx = WithMessageID(ctx, messageID)
	}

	return ctx
}

// injectTraceContext extracts trace context from headers and stores it in the context.
func injectTraceContext(ctx context.Context, headers nats.Header) context.Context {
	tc, ok := parseTraceHeaders(headers)
	if ok {
		ctx = WithTraceContext(ctx, tc)
	}
	return ctx
}

// parseTraceHeaders reads trace IDs from custom and B3 headers.
// Returns the parsed TraceContext and true if any trace info was found.
func parseTraceHeaders(headers nats.Header) (*TraceContext, bool) {
	tc := &TraceContext{}
	hasTrace := false

	tc.TraceID = headerWithFallback(headers, HeaderTraceID, HeaderB3TraceID)
	if tc.TraceID != "" {
		hasTrace = true
	}

	tc.SpanID = headerWithFallback(headers, HeaderSpanID, HeaderB3SpanID)
	if tc.SpanID != "" {
		hasTrace = true
	}

	if sampledFlag := headers.Get(HeaderB3Sampled); sampledFlag == "1" || sampledFlag == "true" {
		tc.Sampled = true
	}

	if traceState := headers.Get(HeaderTraceState); traceState != "" {
		tc.State = traceState
	}

	return tc, hasTrace
}

// headerWithFallback returns the value of the primary header, falling back to the fallback header.
func headerWithFallback(headers nats.Header, primary, fallback string) string {
	if v := headers.Get(primary); v != "" {
		return v
	}
	return headers.Get(fallback)
}
