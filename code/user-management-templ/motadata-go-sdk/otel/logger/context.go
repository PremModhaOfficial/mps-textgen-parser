package logger

import (
	"context"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/core/pool/resourcepool"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// ctxKey is a private type for context keys to avoid collisions.
// Using a custom type prevents accidental key collisions with other packages
// that might use string keys in the context.
type ctxKey string

/* ---------------------------------------- Constants -------------------------------------------------- */

// Context keys for log field extraction.
// These keys are used to store and retrieve logging-related values from context.
// Each key is prefixed with "log_" to clearly identify them as logger-specific.
const (
	// tenantIDKey stores the tenant identifier for multi-tenant applications
	tenantIDKey ctxKey = "log_tenant_id"

	// requestIDKey stores the unique request identifier for request tracing
	requestIDKey ctxKey = "log_request_id"

	// userIDKey stores the user identifier for user-specific logging
	userIDKey ctxKey = "log_user_id"

	// correlationIDKey stores the correlation ID for distributed tracing
	correlationIDKey ctxKey = "log_correlation_id"
)

/* ---------------------------------------- Field Pool -------------------------------------------------- */

// FieldSliceWrapper wraps a zap.Field slice for pooling.
// This wrapper enables object pooling of field slices, reducing allocations
// during high-throughput logging operations.
//
// The pooling mechanism significantly improves performance when logging
// at high rates by reusing pre-allocated slices instead of allocating new ones.
type FieldSliceWrapper struct {
	// Fields holds the extracted context fields
	// Capacity is pre-allocated to OTELMaxContextFields to avoid reallocation
	Fields []zap.Field
}

// fieldPool provides reusable slices for context field extraction using ResourcePool.
// The pool is initialized in init() and provides zero-allocation field extraction
// for the common case where context contains standard fields.
var fieldPool *resourcepool.Pool[*FieldSliceWrapper]

// init initializes the field pool with pre-configured settings.
// This runs automatically at package initialization and sets up the
// object pool used for efficient field slice management.
func init() {

	// Create a resource pool for field slice wrappers
	// The pool uses ResourcePool from core/pool for thread-safe pooling
	pool, err := resourcepool.NewResourcePool(resourcepool.PoolConfig[*FieldSliceWrapper]{
		// Maximum number of pooled wrappers
		// Higher values use more memory but reduce allocation under high load
		MaxSize: utils.OTELFieldPoolMaxSize,

		// OnCreate allocates a new wrapper with pre-sized field slice
		OnCreate: func() (*FieldSliceWrapper, error) {

			return &FieldSliceWrapper{
				// Pre-allocate capacity for expected number of context fields
				// This avoids slice growth during field extraction
				Fields: make([]zap.Field, 0, utils.OTELMaxContextFields),
			}, nil
		},

		// OnReset clears the slice while preserving capacity
		// This allows reuse without reallocation
		OnReset: func(wrapper *FieldSliceWrapper) error {

			// Reset slice length to 0 while keeping the underlying array
			wrapper.Fields = wrapper.Fields[:0]

			return nil
		},
	})

	// Panic on pool initialization failure as logging is a critical service
	// This should only happen due to programming errors (e.g., invalid config)
	if err != nil {
		panic("failed to initialize field pool: " + err.Error())
	}

	fieldPool = pool
}

// getFieldSlice acquires a field slice wrapper from the pool.
// If the pool is exhausted, it falls back to creating a new wrapper.
// This ensures logging never blocks due to pool exhaustion.
//
// Returns:
//   - *FieldSliceWrapper: A wrapper ready for use
func getFieldSlice() *FieldSliceWrapper {

	// TryGet returns immediately without blocking
	// This is preferred for logging to avoid latency spikes
	wrapper, err := fieldPool.TryGet()

	if err != nil {
		// Fallback to creating a new wrapper if pool is exhausted
		// This ensures logging never fails due to pool capacity
		return &FieldSliceWrapper{
			Fields: make([]zap.Field, 0, utils.OTELMaxContextFields),
		}
	}

	return wrapper
}

// putFieldSlice releases a field slice wrapper back to the pool.
// This should be called after the fields have been used for logging.
//
// Parameters:
//   - wrapper: The wrapper to return (nil is safely ignored)
func putFieldSlice(wrapper *FieldSliceWrapper) {

	// Guard against nil to allow unconditional calls
	if wrapper == nil {
		return
	}

	// Return to pool - errors are ignored as the pool handles overflow gracefully
	_ = fieldPool.Put(wrapper)
}

/* ---------------------------------------- Context Setters -------------------------------------------------- */

// WithTenantID adds a tenant ID to the context for multi-tenant logging.
// This allows automatic inclusion of tenant_id in all log entries
// made with this context.
//
// Parameters:
//   - ctx: The parent context
//   - id: The tenant identifier
//
// Returns:
//   - context.Context: A new context with the tenant ID
//
// Example:
//
//	ctx = logger.WithTenantID(ctx, "tenant-123")
//	logger.Info(ctx, "processing request")  // Includes tenant_id="tenant-123"
func WithTenantID(ctx context.Context, id string) context.Context {

	return context.WithValue(ctx, tenantIDKey, id)
}

// WithRequestID adds a request ID to the context for request tracing.
// This allows automatic inclusion of request_id in all log entries
// made with this context.
//
// Parameters:
//   - ctx: The parent context
//   - id: The request identifier (typically a UUID)
//
// Returns:
//   - context.Context: A new context with the request ID
//
// Example:
//
//	ctx = logger.WithRequestID(ctx, uuid.New().String())
//	logger.Info(ctx, "handling request")  // Includes request_id="..."
func WithRequestID(ctx context.Context, id string) context.Context {

	return context.WithValue(ctx, requestIDKey, id)
}

// WithUserID adds a user ID to the context for user-specific logging.
// This allows automatic inclusion of user_id in all log entries
// made with this context.
//
// Parameters:
//   - ctx: The parent context
//   - id: The user identifier
//
// Returns:
//   - context.Context: A new context with the user ID
//
// Example:
//
//	ctx = logger.WithUserID(ctx, "user-456")
//	logger.Info(ctx, "user action")  // Includes user_id="user-456"
func WithUserID(ctx context.Context, id string) context.Context {

	return context.WithValue(ctx, userIDKey, id)
}

// WithCorrelationID adds a correlation ID to the context for distributed tracing.
// This allows automatic inclusion of correlation_id in all log entries
// made with this context.
//
// Parameters:
//   - ctx: The parent context
//   - id: The correlation identifier
//
// Returns:
//   - context.Context: A new context with the correlation ID
//
// Example:
//
//	ctx = logger.WithCorrelationID(ctx, "corr-789")
//	logger.Info(ctx, "processing step")  // Includes correlation_id="corr-789"
func WithCorrelationID(ctx context.Context, id string) context.Context {

	return context.WithValue(ctx, correlationIDKey, id)
}

/* ---------------------------------------- Field Extraction Functions -------------------------------------------------- */

// extractFieldsPooled extracts all context fields using a pooled slice.
// This is the high-performance version that reuses field slices from a pool.
//
// The function extracts:
//   - trace_id and span_id from OTEL span context (if present)
//   - tenant_id, request_id, user_id, correlation_id from context values
//
// Parameters:
//   - ctx: The context to extract fields from
//
// Returns:
//   - *FieldSliceWrapper: Wrapper containing extracted fields (nil if no fields)
//
// Important: The caller MUST return the wrapper to the pool via putFieldSlice()
// after using the fields.
func extractFieldsPooled(ctx context.Context) *FieldSliceWrapper {

	// Nil context has no fields to extract
	if ctx == nil {

		return nil
	}

	// Extract OTEL span context for trace correlation
	// This enables log-to-trace correlation in observability tools
	spanCtx := trace.SpanContextFromContext(ctx)
	hasSpan := spanCtx.IsValid()

	// Extract custom context values
	// Each extraction includes a non-empty check to avoid empty fields
	tenantID, hasTenant := ctx.Value(tenantIDKey).(string)
	hasTenant = hasTenant && tenantID != ""

	requestID, hasRequest := ctx.Value(requestIDKey).(string)
	hasRequest = hasRequest && requestID != ""

	userID, hasUser := ctx.Value(userIDKey).(string)
	hasUser = hasUser && userID != ""

	correlationID, hasCorrelation := ctx.Value(correlationIDKey).(string)
	hasCorrelation = hasCorrelation && correlationID != ""

	// Fast path: no fields to extract, return nil to avoid pool operations
	if !hasSpan && !hasTenant && !hasRequest && !hasUser && !hasCorrelation {

		return nil
	}

	// Get a wrapper from the pool for field storage
	wrapper := getFieldSlice()

	// Add span context fields for trace correlation
	// These enable correlation between logs and traces in observability backends
	if hasSpan {

		wrapper.Fields = append(wrapper.Fields,
			zap.String("trace_id", spanCtx.TraceID().String()),
			zap.String("span_id", spanCtx.SpanID().String()),
		)
	}

	// Add custom context fields
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

// extractFields extracts all context fields (allocating version).
// This is a convenience function for cases where pooled allocation is not needed.
// It copies the pooled fields to a new slice and returns the pooled wrapper.
//
// Parameters:
//   - ctx: The context to extract fields from
//
// Returns:
//   - []zap.Field: Slice of extracted fields (nil if no fields)
func extractFields(ctx context.Context) []zap.Field {

	// Use the pooled extraction
	wrapper := extractFieldsPooled(ctx)

	if wrapper == nil {

		return nil
	}

	// Copy fields to a new slice to allow returning the wrapper to the pool
	// This is less efficient than the pooled version but simpler to use
	result := make([]zap.Field, len(wrapper.Fields))
	copy(result, wrapper.Fields)

	// Return wrapper to pool for reuse
	putFieldSlice(wrapper)

	return result
}
