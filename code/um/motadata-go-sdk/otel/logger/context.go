package logger

import (
	"context"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/core/pool/resourcepool"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// Context keys for log field extraction
type ctxKey string

/* ---------------------------------------- Constants -------------------------------------------------- */

const (
	tenantIDKey ctxKey = "log_tenant_id"

	requestIDKey ctxKey = "log_request_id"

	userIDKey ctxKey = "log_user_id"

	correlationIDKey ctxKey = "log_correlation_id"
)

// Maximum number of context fields expected
const maxContextFields = 6

// Maximum pool size for field slices
const fieldPoolMaxSize = 1024

/* ---------------------------------------- Field Pool -------------------------------------------------- */

// FieldSliceWrapper wraps a zap.Field slice for pooling
type FieldSliceWrapper struct {
	Fields []zap.Field
}

// fieldPool provides reusable slices for context field extraction using ResourcePool
var fieldPool *resourcepool.Pool[*FieldSliceWrapper]

// backgroundCtx is a cached background context for pool operations
var backgroundCtx = context.Background()

// init initializes the field pool
func init() {

	pool, err := resourcepool.NewResourcePool(resourcepool.PoolConfig[*FieldSliceWrapper]{
		MaxSize: fieldPoolMaxSize,
		OnCreate: func() (*FieldSliceWrapper, error) {
			return &FieldSliceWrapper{
				Fields: make([]zap.Field, 0, maxContextFields),
			}, nil
		},
		OnReset: func(wrapper *FieldSliceWrapper) error {
			wrapper.Fields = wrapper.Fields[:0]
			return nil
		},
	})

	if err != nil {
		panic("failed to initialize field pool: " + err.Error())
	}

	fieldPool = pool
}

// getFieldSlice acquires a field slice from the pool
func getFieldSlice() *FieldSliceWrapper {
	wrapper, err := fieldPool.TryGet()

	if err != nil {
		// Fallback to creating a new wrapper if pool is exhausted
		return &FieldSliceWrapper{
			Fields: make([]zap.Field, 0, maxContextFields),
		}
	}

	return wrapper
}

// putFieldSlice releases a field slice back to the pool
func putFieldSlice(wrapper *FieldSliceWrapper) {
	if wrapper == nil {
		return
	}

	_ = fieldPool.Put(wrapper)
}

/* ---------------------------------------- Context Setters -------------------------------------------------- */

// WithTenantID adds a tenant ID to the context
func WithTenantID(ctx context.Context, id string) context.Context {

	return context.WithValue(ctx, tenantIDKey, id)
}

// WithRequestID adds a request ID to the context
func WithRequestID(ctx context.Context, id string) context.Context {

	return context.WithValue(ctx, requestIDKey, id)
}

// WithUserID adds a user ID to the context
func WithUserID(ctx context.Context, id string) context.Context {

	return context.WithValue(ctx, userIDKey, id)
}

// WithCorrelationID adds a correlation ID to the context
func WithCorrelationID(ctx context.Context, id string) context.Context {

	return context.WithValue(ctx, correlationIDKey, id)
}

/* ---------------------------------------- Field Extraction Functions -------------------------------------------------- */

// extractFieldsPooled extracts all context fields using a pooled slice
func extractFieldsPooled(ctx context.Context) *FieldSliceWrapper {

	if ctx == nil {

		return nil
	}

	spanCtx := trace.SpanContextFromContext(ctx)
	hasSpan := spanCtx.IsValid()

	tenantID, hasTenant := ctx.Value(tenantIDKey).(string)
	hasTenant = hasTenant && tenantID != ""

	requestID, hasRequest := ctx.Value(requestIDKey).(string)
	hasRequest = hasRequest && requestID != ""

	userID, hasUser := ctx.Value(userIDKey).(string)
	hasUser = hasUser && userID != ""

	correlationID, hasCorrelation := ctx.Value(correlationIDKey).(string)
	hasCorrelation = hasCorrelation && correlationID != ""

	if !hasSpan && !hasTenant && !hasRequest && !hasUser && !hasCorrelation {

		return nil
	}

	wrapper := getFieldSlice()

	if hasSpan {

		wrapper.Fields = append(wrapper.Fields,
			zap.String("trace_id", spanCtx.TraceID().String()),
			zap.String("span_id", spanCtx.SpanID().String()),
		)
	}

	if hasTenant {

		wrapper.Fields = append(wrapper.Fields, zap.String("tenant_id", tenantID))
	}

	if hasRequest {

		wrapper.Fields = append(wrapper.Fields, zap.String("request_id", requestID))
	}

	if hasUser {

		wrapper.Fields = append(wrapper.Fields, zap.String("user_id", userID))
	}

	if hasCorrelation {

		wrapper.Fields = append(wrapper.Fields, zap.String("correlation_id", correlationID))
	}

	return wrapper
}

// extractFields extracts all context fields (allocating version for backward compatibility)
func extractFields(ctx context.Context) []zap.Field {

	wrapper := extractFieldsPooled(ctx)

	if wrapper == nil {

		return nil
	}

	result := make([]zap.Field, len(wrapper.Fields))
	copy(result, wrapper.Fields)

	putFieldSlice(wrapper)

	return result
}
