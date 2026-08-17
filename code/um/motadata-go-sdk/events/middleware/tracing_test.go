package middleware

import (
	"context"
	"errors"
	"testing"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
	"github.com/stretchr/testify/assert"
)

/* ========================================================================================================
   TRACING MIDDLEWARE TESTS
   ======================================================================================================== */

func TestNewTracingMiddleware(t *testing.T) {

	assertions := assert.New(t)

	mw := NewTracingMiddleware()

	assertions.NotNil(mw, "TracingMiddleware should not be nil")
	assertions.NotNil(mw.TraceIDGenerator, "TraceIDGenerator should not be nil")
	assertions.NotNil(mw.SpanIDGenerator, "SpanIDGenerator should not be nil")
	assertions.NotNil(mw.Sampler, "Sampler should not be nil")
	assertions.True(mw.UseOTEL, "UseOTEL should be true by default")
}

func TestTracingMiddlewareInterceptPublish(t *testing.T) {

	assertions := assert.New(t)

	mw := NewTracingMiddleware()
	mw.UseOTEL = false // Disable OTEL for testing

	var capturedCtx context.Context
	var capturedSubject string
	var capturedMsg *core.Message

	handler := func(ctx context.Context, subject string, msg *core.Message) error {

		capturedCtx = ctx
		capturedSubject = subject
		capturedMsg = msg

		return nil
	}

	interceptor := mw.InterceptPublish()
	wrappedHandler := interceptor(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test data"))

	err := wrappedHandler(ctx, "test.subject", msg)

	assertions.NoError(err, "Handler should not return error")
	assertions.Equal("test.subject", capturedSubject, "Subject should be passed through")
	assertions.Equal(msg, capturedMsg, "Message should be passed through")

	// Verify trace context was added
	tc, ok := core.TraceContextFromContext(capturedCtx)
	assertions.True(ok, "TraceContext should be added to context")
	assertions.NotEmpty(tc.TraceID, "TraceID should be generated")
	assertions.NotEmpty(tc.SpanID, "SpanID should be generated")

	// Verify headers were injected
	assertions.NotEmpty(msg.Headers.Get(core.HeaderTraceID), "TraceID header should be set")
	assertions.NotEmpty(msg.Headers.Get(core.HeaderSpanID), "SpanID header should be set")
}

func TestTracingMiddlewareInterceptPublishWithExistingTrace(t *testing.T) {

	assertions := assert.New(t)

	mw := NewTracingMiddleware()
	mw.UseOTEL = false

	var capturedCtx context.Context

	handler := func(ctx context.Context, subject string, msg *core.Message) error {

		capturedCtx = ctx

		return nil
	}

	interceptor := mw.InterceptPublish()
	wrappedHandler := interceptor(handler)

	// Create context with existing trace
	ctx := context.Background()
	ctx = core.WithTraceContext(ctx, &core.TraceContext{
		TraceID: "existing-trace-id",
		SpanID:  "existing-span-id",
		Sampled: true,
	})

	msg := core.NewMessage([]byte("test data"))
	err := wrappedHandler(ctx, "test.subject", msg)

	assertions.NoError(err, "Handler should not return error")

	// Verify trace context was preserved/updated
	tc, ok := core.TraceContextFromContext(capturedCtx)
	assertions.True(ok, "TraceContext should exist")
	assertions.Equal("existing-trace-id", tc.TraceID, "TraceID should be preserved")
	assertions.NotEqual("existing-span-id", tc.SpanID, "SpanID should be new")
	assertions.Equal("existing-span-id", tc.ParentID, "ParentID should be the old span")
}

func TestTracingMiddlewareInterceptSubscribe(t *testing.T) {

	assertions := assert.New(t)

	mw := NewTracingMiddleware()
	mw.UseOTEL = false

	var capturedCtx context.Context
	var capturedMsg *core.Message

	handler := func(ctx context.Context, msg *core.Message) error {

		capturedCtx = ctx
		capturedMsg = msg

		return nil
	}

	interceptor := mw.InterceptSubscribe()
	wrappedHandler := interceptor(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test data")).
		WithSubject("test.subject").
		WithHeader(core.HeaderTraceID, "incoming-trace-id").
		WithHeader(core.HeaderSpanID, "incoming-span-id").
		WithHeader(core.HeaderB3Sampled, "1")

	err := wrappedHandler(ctx, msg)

	assertions.NoError(err, "Handler should not return error")
	assertions.Equal(msg, capturedMsg, "Message should be passed through")

	// Verify trace context was extracted from headers
	tc, ok := core.TraceContextFromContext(capturedCtx)
	assertions.True(ok, "TraceContext should exist")
	assertions.Equal("incoming-trace-id", tc.TraceID, "TraceID should be from headers")
	assertions.Equal("incoming-span-id", tc.SpanID, "SpanID should be from headers")
	assertions.True(tc.Sampled, "Sampled should be true from headers")
}

func TestTracingMiddlewareInterceptSubscribeNoHeaders(t *testing.T) {

	assertions := assert.New(t)

	mw := NewTracingMiddleware()
	mw.UseOTEL = false

	var capturedCtx context.Context

	handler := func(ctx context.Context, msg *core.Message) error {

		capturedCtx = ctx

		return nil
	}

	interceptor := mw.InterceptSubscribe()
	wrappedHandler := interceptor(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test data"))

	err := wrappedHandler(ctx, msg)

	assertions.NoError(err, "Handler should not return error")

	// Verify trace context was created
	tc, ok := core.TraceContextFromContext(capturedCtx)
	assertions.True(ok, "TraceContext should be created")
	assertions.NotEmpty(tc.TraceID, "TraceID should be generated")
	assertions.NotEmpty(tc.SpanID, "SpanID should be generated")
}

/* ========================================================================================================
   DEFAULT ID GENERATOR TESTS
   ======================================================================================================== */

func TestDefaultIDGenerator(t *testing.T) {

	assertions := assert.New(t)

	id1 := defaultIDGenerator()
	id2 := defaultIDGenerator()

	assertions.NotEmpty(id1, "ID should not be empty")
	assertions.NotEmpty(id2, "ID should not be empty")
	assertions.NotEqual(id1, id2, "IDs should be unique")
	assertions.Equal(32, len(id1), "ID should be 32 hex characters (16 bytes)")
}

func TestFallbackIDGenerator(t *testing.T) {

	assertions := assert.New(t)

	id1 := fallbackIDGenerator()
	id2 := fallbackIDGenerator()

	assertions.NotEmpty(id1, "Fallback ID should not be empty")
	assertions.NotEmpty(id2, "Fallback ID should not be empty")
	assertions.NotEqual(id1, id2, "Fallback IDs should be unique")
}

/* ========================================================================================================
   W3C TRACE CONTEXT TESTS
   ======================================================================================================== */

func TestExtractW3CTraceParent(t *testing.T) {

	assertions := assert.New(t)

	// Valid traceparent
	tc := ExtractW3CTraceParent("00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01")

	assertions.NotNil(tc, "TraceContext should not be nil for valid traceparent")
	assertions.Equal("4bf92f3577b34da6a3ce929d0e0e4736", tc.TraceID, "TraceID should match")
	assertions.Equal("00f067aa0ba902b7", tc.SpanID, "SpanID should match")
	assertions.True(tc.Sampled, "Sampled should be true for flag 01")
}

func TestExtractW3CTraceParentNotSampled(t *testing.T) {

	assertions := assert.New(t)

	tc := ExtractW3CTraceParent("00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00")

	assertions.NotNil(tc, "TraceContext should not be nil")
	assertions.False(tc.Sampled, "Sampled should be false for flag 00")
}

func TestExtractW3CTraceParentInvalid(t *testing.T) {

	testCases := []struct {
		name        string
		traceparent string
	}{
		{"too short", "00-abc-def-01"},
		{"empty", ""},
		{"missing parts", "00-abcd"},
		{"wrong format", "invalid-format-here-now"},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			assertions := assert.New(t)
			result := ExtractW3CTraceParent(tc.traceparent)
			assertions.Nil(result, "TraceContext should be nil for invalid traceparent")
		})
	}
}

func TestFormatW3CTraceParent(t *testing.T) {

	assertions := assert.New(t)

	tc := &core.TraceContext{
		TraceID: "4bf92f3577b34da6a3ce929d0e0e4736",
		SpanID:  "00f067aa0ba902b7",
		Sampled: true,
	}

	result := FormatW3CTraceParent(tc)

	assertions.Equal("00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01", result, "Traceparent should be correctly formatted")
}

func TestFormatW3CTraceParentNotSampled(t *testing.T) {

	assertions := assert.New(t)

	tc := &core.TraceContext{
		TraceID: "4bf92f3577b34da6a3ce929d0e0e4736",
		SpanID:  "00f067aa0ba902b7",
		Sampled: false,
	}

	result := FormatW3CTraceParent(tc)

	assertions.Equal("00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00", result, "Traceparent should have 00 flag")
}

func TestFormatW3CTraceParentNil(t *testing.T) {

	assertions := assert.New(t)

	result := FormatW3CTraceParent(nil)

	assertions.Empty(result, "Should return empty string for nil TraceContext")
}

func TestFormatW3CTraceParentMissingFields(t *testing.T) {

	assertions := assert.New(t)

	tc := &core.TraceContext{
		TraceID: "",
		SpanID:  "00f067aa0ba902b7",
	}

	result := FormatW3CTraceParent(tc)

	assertions.Empty(result, "Should return empty string when TraceID is missing")
}

/* ========================================================================================================
   TRACING HELPER FUNCTION TESTS
   ======================================================================================================== */

func TestTracing(t *testing.T) {

	assertions := assert.New(t)

	mw := Tracing()

	assertions.NotNil(mw, "Tracing() should return TracingMiddleware")
	assertions.NotNil(mw.TraceIDGenerator, "TraceIDGenerator should not be nil")
}

/* ========================================================================================================
   ERROR HANDLING TESTS
   ======================================================================================================== */

func TestTracingMiddlewarePublishError(t *testing.T) {

	assertions := assert.New(t)

	mw := NewTracingMiddleware()
	mw.UseOTEL = false

	expectedErr := errors.New("publish failed")

	handler := func(ctx context.Context, subject string, msg *core.Message) error {

		return expectedErr
	}

	interceptor := mw.InterceptPublish()
	wrappedHandler := interceptor(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	err := wrappedHandler(ctx, "test.subject", msg)

	assertions.ErrorIs(err, expectedErr, "Error should be propagated")
}

func TestTracingMiddlewareSubscribeError(t *testing.T) {

	assertions := assert.New(t)

	mw := NewTracingMiddleware()
	mw.UseOTEL = false

	expectedErr := errors.New("process failed")

	handler := func(ctx context.Context, msg *core.Message) error {

		return expectedErr
	}

	interceptor := mw.InterceptSubscribe()
	wrappedHandler := interceptor(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	err := wrappedHandler(ctx, msg)

	assertions.ErrorIs(err, expectedErr, "Error should be propagated")
}
