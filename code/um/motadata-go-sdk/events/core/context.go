package core

import (
	"context"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer"
)

// Context keys for storing event-related values
type contextKey int

const (
	tenantIDKey contextKey = iota
	traceContextKey
	messageIDKey
	correlationIDKey
	metadataKey
)

// WithTenantID adds tenant ID to context
func WithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantIDKey, tenantID)
}

// TenantIDFromContext retrieves tenant ID from context
func TenantIDFromContext(ctx context.Context) (string, bool) {
	tenantIDValue, exists := ctx.Value(tenantIDKey).(string)
	return tenantIDValue, exists
}

// WithTraceContext adds trace context to context
func WithTraceContext(ctx context.Context, traceContext *TraceContext) context.Context {
	return context.WithValue(ctx, traceContextKey, traceContext)
}

// TraceContextFromContext retrieves trace context from context
func TraceContextFromContext(ctx context.Context) (*TraceContext, bool) {
	traceContextValue, exists := ctx.Value(traceContextKey).(*TraceContext)
	return traceContextValue, exists
}

// WithMessageID adds message ID to context
func WithMessageID(ctx context.Context, messageID string) context.Context {
	return context.WithValue(ctx, messageIDKey, messageID)
}

// MessageIDFromContext retrieves message ID from context
func MessageIDFromContext(ctx context.Context) (string, bool) {
	messageIDValue, exists := ctx.Value(messageIDKey).(string)
	return messageIDValue, exists
}

// WithCorrelationID adds correlation ID to context
func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, correlationIDKey, correlationID)
}

// CorrelationIDFromContext retrieves correlation ID from context
func CorrelationIDFromContext(ctx context.Context) (string, bool) {
	correlationIDValue, exists := ctx.Value(correlationIDKey).(string)
	return correlationIDValue, exists
}

// WithMetadata adds metadata to context
func WithMetadata(ctx context.Context, metadata *Metadata) context.Context {
	return context.WithValue(ctx, metadataKey, metadata)
}

// MetadataFromContext retrieves metadata from context
func MetadataFromContext(ctx context.Context) (*Metadata, bool) {
	metadataValue, exists := ctx.Value(metadataKey).(*Metadata)
	return metadataValue, exists
}

// ExtractHeaders extracts context values into headers
func ExtractHeaders(ctx context.Context, headers Headers) Headers {
	if headers == nil {
		headers = make(Headers)
	}

	// Extract tenant ID
	if tenantID, exists := TenantIDFromContext(ctx); exists {
		headers.Set(HeaderTenantID, tenantID)
	}

	// First try to extract from OTEL span context (preferred)
	traceID := tracer.TraceIDFromContext(ctx)
	spanID := tracer.SpanIDFromContext(ctx)

	if traceID != "" && spanID != "" {
		// Got trace context from OTEL
		headers.Set(HeaderTraceID, traceID)
		headers.Set(HeaderB3TraceID, traceID)
		headers.Set(HeaderSpanID, spanID)
		headers.Set(HeaderB3SpanID, spanID)
		headers.Set(HeaderB3Sampled, "1") // OTEL spans are always sampled
		// W3C traceparent format: version-traceId-spanId-flags
		headers.Set(HeaderTraceParent, "00-"+traceID+"-"+spanID+"-01")
	} else if traceContext, exists := TraceContextFromContext(ctx); exists {
		// Fall back to local trace context
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
		// W3C traceparent format
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

	// Extract correlation ID
	if correlationID, exists := CorrelationIDFromContext(ctx); exists {
		headers.Set(HeaderCorrelationID, correlationID)
	}

	// Extract message ID
	if messageID, exists := MessageIDFromContext(ctx); exists {
		headers.Set(HeaderMessageID, messageID)
	}

	return headers
}

// ExtractOTELTraceContext extracts OTEL trace context from the current span
// and returns a TraceContext that can be propagated
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

// InjectContext injects header values into context
func InjectContext(ctx context.Context, headers Headers) context.Context {
	if headers == nil {
		return ctx
	}

	// Inject tenant ID
	if tenantID := headers.Get(HeaderTenantID); tenantID != "" {
		ctx = WithTenantID(ctx, tenantID)
	}

	// Inject trace context
	traceContext := &TraceContext{}
	hasTraceContext := false

	if traceID := headers.Get(HeaderTraceID); traceID != "" {
		traceContext.TraceID = traceID
		hasTraceContext = true
	} else if traceIDFromB3 := headers.Get(HeaderB3TraceID); traceIDFromB3 != "" {
		traceContext.TraceID = traceIDFromB3
		hasTraceContext = true
	}

	if spanID := headers.Get(HeaderSpanID); spanID != "" {
		traceContext.SpanID = spanID
		hasTraceContext = true
	} else if spanIDFromB3 := headers.Get(HeaderB3SpanID); spanIDFromB3 != "" {
		traceContext.SpanID = spanIDFromB3
		hasTraceContext = true
	}

	if sampledFlag := headers.Get(HeaderB3Sampled); sampledFlag == "1" || sampledFlag == "true" {
		traceContext.Sampled = true
	}

	if traceState := headers.Get(HeaderTraceState); traceState != "" {
		traceContext.State = traceState
	}

	if hasTraceContext {
		ctx = WithTraceContext(ctx, traceContext)
	}

	// Inject correlation ID
	if correlationID := headers.Get(HeaderCorrelationID); correlationID != "" {
		ctx = WithCorrelationID(ctx, correlationID)
	}

	// Inject message ID
	if messageID := headers.Get(HeaderMessageID); messageID != "" {
		ctx = WithMessageID(ctx, messageID)
	}

	return ctx
}
