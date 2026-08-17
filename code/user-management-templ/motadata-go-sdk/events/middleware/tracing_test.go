package middleware

import (
	"context"
	"errors"
	"testing"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

/* ========================================================================================================
   TRACING MIDDLEWARE TESTS
   ======================================================================================================== */

var (
	tracingPublishFailed = "publish failed"
	tracingProcessFailed = "process failed"
)

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
	mw.UseOTEL = false

	var capturedCtx context.Context
	var capturedMsg *nats.Msg

	handler := func(ctx context.Context, msg *nats.Msg) error {
		capturedCtx = ctx
		capturedMsg = msg
		return nil
	}

	wrapped := mw.InterceptPublish()(handler)
	msg := newTestMsg(testSubject, testPayloadData)

	err := wrapped(context.Background(), msg)

	assertions.NoError(err, "Handler should not return error")
	assertions.Equal(msg, capturedMsg, "Message should be passed through")

	tc, ok := core.TraceContextFromContext(capturedCtx)
	assertions.True(ok, "TraceContext should be added to context")
	assertions.NotEmpty(tc.TraceID, "TraceID should be generated")
	assertions.NotEmpty(tc.SpanID, "SpanID should be generated")

	assertions.NotEmpty(msg.Header.Get(core.HeaderTraceID), "TraceID header should be set")
	assertions.NotEmpty(msg.Header.Get(core.HeaderSpanID), "SpanID header should be set")
}

func TestTracingMiddlewareInterceptPublishWithExistingTrace(t *testing.T) {
	assertions := assert.New(t)

	mw := NewTracingMiddleware()
	mw.UseOTEL = false

	var capturedCtx context.Context

	handler := func(ctx context.Context, msg *nats.Msg) error {
		capturedCtx = ctx
		return nil
	}

	wrapped := mw.InterceptPublish()(handler)

	ctx := core.WithTraceContext(context.Background(), &core.TraceContext{
		TraceID: testExistingTrace,
		SpanID:  testExistingSpan,
		Sampled: true,
	})

	err := wrapped(ctx, newTestMsg(testSubject, testPayloadData))

	assertions.NoError(err, "Handler should not return error")

	tc, ok := core.TraceContextFromContext(capturedCtx)
	assertions.True(ok, "TraceContext should exist")
	assertions.Equal(testExistingTrace, tc.TraceID, "TraceID should be preserved")
	assertions.NotEqual(testExistingSpan, tc.SpanID, "SpanID should be new")
	assertions.Equal(testExistingSpan, tc.ParentID, "ParentID should be the old span")
}

func TestTracingMiddlewareInterceptSubscribe(t *testing.T) {
	assertions := assert.New(t)

	mw := NewTracingMiddleware()
	mw.UseOTEL = false

	var capturedCtx context.Context
	var capturedMsg *nats.Msg

	handler := func(ctx context.Context, msg *nats.Msg) error {
		capturedCtx = ctx
		capturedMsg = msg
		return nil
	}

	wrapped := mw.InterceptSubscribe()(handler)

	msg := newTestMsg(testSubject, testPayloadData)
	msg.Header.Set(core.HeaderTraceID, testIncomingTrace)
	msg.Header.Set(core.HeaderSpanID, testIncomingSpan)
	msg.Header.Set(core.HeaderB3Sampled, "1")

	err := wrapped(context.Background(), msg)

	assertions.NoError(err, "Handler should not return error")
	assertions.Equal(msg, capturedMsg, "Message should be passed through")

	tc, ok := core.TraceContextFromContext(capturedCtx)
	assertions.True(ok, "TraceContext should exist")
	assertions.Equal(testIncomingTrace, tc.TraceID, "TraceID should be from headers")
	assertions.Equal(testIncomingSpan, tc.SpanID, "SpanID should be from headers")
	assertions.True(tc.Sampled, "Sampled should be true from headers")
}

func TestTracingMiddlewareInterceptSubscribeNoHeaders(t *testing.T) {
	assertions := assert.New(t)

	mw := NewTracingMiddleware()
	mw.UseOTEL = false

	var capturedCtx context.Context

	handler := func(ctx context.Context, msg *nats.Msg) error {
		capturedCtx = ctx
		return nil
	}

	wrapped := mw.InterceptSubscribe()(handler)
	err := wrapped(context.Background(), newTestMsg("", testPayloadData))

	assertions.NoError(err, "Handler should not return error")

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

	tc := ExtractW3CTraceParent(testW3CParent)

	assertions.NotNil(tc, "TraceContext should not be nil for valid traceparent")
	assertions.Equal(testW3CTraceID, tc.TraceID, "TraceID should match")
	assertions.Equal(testW3CSpanID, tc.SpanID, "SpanID should match")
	assertions.True(tc.Sampled, "Sampled should be true for flag 01")
}

func TestExtractW3CTraceParentNotSampled(t *testing.T) {
	assertions := assert.New(t)

	tc := ExtractW3CTraceParent(testW3CParentUnsmp)

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
	testCases := []struct {
		name     string
		tc       *core.TraceContext
		expected string
	}{
		{
			"sampled",
			&core.TraceContext{TraceID: testW3CTraceID, SpanID: testW3CSpanID, Sampled: true},
			testW3CParent,
		},
		{
			"not sampled",
			&core.TraceContext{TraceID: testW3CTraceID, SpanID: testW3CSpanID, Sampled: false},
			testW3CParentUnsmp,
		},
		{
			"nil context",
			nil,
			"",
		},
		{
			"missing TraceID",
			&core.TraceContext{TraceID: "", SpanID: testW3CSpanID},
			"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			result := FormatW3CTraceParent(tc.tc)
			assertions.Equal(tc.expected, result)
		})
	}
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

	expectedErr := errors.New(tracingPublishFailed)
	handler := func(ctx context.Context, msg *nats.Msg) error { return expectedErr }

	wrapped := mw.InterceptPublish()(handler)
	err := wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	assertions.ErrorIs(err, expectedErr, "Error should be propagated")
}

func TestTracingMiddlewareSubscribeError(t *testing.T) {
	assertions := assert.New(t)

	mw := NewTracingMiddleware()
	mw.UseOTEL = false

	expectedErr := errors.New(tracingProcessFailed)
	handler := func(ctx context.Context, msg *nats.Msg) error { return expectedErr }

	wrapped := mw.InterceptSubscribe()(handler)
	err := wrapped(context.Background(), newTestMsg("", testPayload))

	assertions.ErrorIs(err, expectedErr, "Error should be propagated")
}
