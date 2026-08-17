package core

import (
	"context"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
)

/* ========================================================================================================
   TEST CONSTANTS
   ======================================================================================================== */

var (
	// context_test.go literals
	testTenantID123     = "tenant-123"
	testMsgID456        = "msg-456"
	testCorrID789       = "corr-789"
	testTraceABC        = "trace-abc"
	testSpanDEF         = "span-def"
	testStateData       = "state-data"
	testTraceID123      = "trace-123"
	testSpanID456       = "span-456"
	testParentID789     = "parent-789"
	testMsgID123        = "msg-123"
	testCorrID123       = "corr-123"
	testTenantMeta      = "tenant-meta"
	testCorrMeta        = "corr-meta"
	testMsgMeta         = "msg-meta"
	testEnvProduction   = "production"
	testRegionUSEast    = "us-east"
	testCustomValue     = "custom-value"
	testCustomHeader    = "Custom-Header"
	testTraceOnly       = "trace-only"
	testTraceState2V    = "vendor1=value1,vendor2=value2"
	testTenantOTEL      = "tenant-otel"
	testOTELTraceIDHex  = "4bf92f3577b34da6a3ce929d0e0e4736"
	testOTELSpanIDHex   = "00f067aa0ba902b7"
	testB3Trace123      = "b3-trace-123"
	testB3Span456       = "b3-span-456"
	testCustomTrace     = "custom-trace"
	testB3Trace         = "b3-trace"
	testCustomSpan      = "custom-span"
	testB3Span          = "b3-span"
	testCorrOnly        = "corr-only"
	testMsgOnly         = "msg-only"
	testTenantA         = "tenant-A"
	testTenantB         = "tenant-B"
	testCorrA           = "corr-A"
	testCorrB           = "corr-B"
	testTenantAlpha     = "tenant-alpha"
	testCorrAlpha       = "corr-alpha"
	testTraceAlpha      = "trace-alpha"
	testSpanAlpha       = "span-alpha"
	testTenantBeta      = "tenant-beta"
	testCorrBeta        = "corr-beta"
	testTraceBeta       = "trace-beta"
	testSpanBeta        = "span-beta"
	testTenantGamma     = "tenant-gamma"
	testCorrGamma       = "corr-gamma"
	testTraceGamma      = "trace-gamma"
	testSpanGamma       = "span-gamma"
	testStreamValue     = "test-stream"
)

/* ========================================================================================================
   TENANT ID CONTEXT TESTS
   ======================================================================================================== */

func TestWithTenantID(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	ctx = WithTenantID(ctx, testTenantID123)

	tenantID, ok := TenantIDFromContext(ctx)

	assertions.True(ok, "TenantID should exist in context")
	assertions.Equal(testTenantID123, tenantID, "TenantID should match")
}

func TestTenantIDFromContextMissing(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()

	tenantID, ok := TenantIDFromContext(ctx)

	assertions.False(ok, "TenantID should not exist in context")
	assertions.Empty(tenantID, "TenantID should be empty")
}

/* ========================================================================================================
   TRACE CONTEXT TESTS
   ======================================================================================================== */

func TestWithTraceContext(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	tc := &TraceContext{
		TraceID:  testTraceID123,
		SpanID:   testSpanID456,
		ParentID: testParentID789,
		Sampled:  true,
		State:    testStateData,
	}

	ctx = WithTraceContext(ctx, tc)

	retrieved, ok := TraceContextFromContext(ctx)

	assertions.True(ok, "TraceContext should must exist in context")
	assertions.Equal(tc, retrieved, "TraceContext should match")
}

func TestTraceContextFromContextMissing(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()

	tc, ok := TraceContextFromContext(ctx)

	assertions.False(ok, "TraceContext should not exist in context")
	assertions.Nil(tc, "TraceContext should be nil")
}

/* ========================================================================================================
   MESSAGE ID CONTEXT TESTS
   ======================================================================================================== */

func TestWithMessageID(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	ctx = WithMessageID(ctx, testMsgID123)

	messageID, ok := MessageIDFromContext(ctx)

	assertions.True(ok, "MessageID should exist in context")
	assertions.Equal(testMsgID123, messageID, "MessageID should match")
}

func TestMessageIDFromContextMissing(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()

	messageID, ok := MessageIDFromContext(ctx)

	assertions.False(ok, "MessageID should not exist in context")
	assertions.Empty(messageID, "MessageID should be empty")
}

/* ========================================================================================================
   CORRELATION ID CONTEXT TESTS
   ======================================================================================================== */

func TestWithCorrelationID(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	ctx = WithCorrelationID(ctx, testCorrID123)

	correlationID, ok := CorrelationIDFromContext(ctx)

	assertions.True(ok, "CorrelationID should exist in context")
	assertions.Equal(testCorrID123, correlationID, "CorrelationID should match")
}

func TestCorrelationIDFromContextMissing(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()

	correlationID, ok := CorrelationIDFromContext(ctx)

	assertions.False(ok, "CorrelationID should not exist in context")
	assertions.Empty(correlationID, "CorrelationID should be empty")
}

/* ========================================================================================================
   METADATA CONTEXT TESTS
   ======================================================================================================== */

func TestWithMetadata(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	metadata := &Metadata{
		TenantID: testTenantID123,
		Custom:   map[string]string{"stream": testStreamValue},
	}

	ctx = WithMetadata(ctx, metadata)

	retrieved, ok := MetadataFromContext(ctx)

	assertions.True(ok, "Metadata should exist in context")
	assertions.Equal(metadata, retrieved, "Metadata should match")
}

func TestMetadataFromContextMissing(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()

	metadata, ok := MetadataFromContext(ctx)

	assertions.False(ok, "Metadata should not exist in context")
	assertions.Nil(metadata, "Metadata should be nil")
}

func TestMetadataWithTimestamp(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	now := time.Now()
	metadata := &Metadata{
		TenantID:      testTenantMeta,
		CorrelationID: testCorrMeta,
		MessageID:     testMsgMeta,
		Timestamp:     now,
		Custom:        map[string]string{"env": testEnvProduction, "region": testRegionUSEast},
	}

	ctx = WithMetadata(ctx, metadata)
	retrieved, ok := MetadataFromContext(ctx)

	assertions.True(ok, "Metadata should exist in context")
	assertions.Equal(testTenantMeta, retrieved.TenantID, "TenantID must match")
	assertions.Equal(testCorrMeta, retrieved.CorrelationID, "CorrelationID should be matched")
	assertions.Equal(testMsgMeta, retrieved.MessageID, "MessageID should be matched")
	assertions.Equal(now, retrieved.Timestamp, "Timestamp should match")
	assertions.Equal(testEnvProduction, retrieved.Custom["env"], "Custom env should match")
	assertions.Equal(testRegionUSEast, retrieved.Custom["region"], "Custom region should match")
}

/* ========================================================================================================
   EXTRACT HEADERS TESTS
   ======================================================================================================== */

func TestExtractHeaders(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	ctx = WithTenantID(ctx, testTenantID123)
	ctx = WithMessageID(ctx, testMsgID456)
	ctx = WithCorrelationID(ctx, testCorrID789)
	ctx = WithTraceContext(ctx, &TraceContext{
		TraceID: testTraceABC,
		SpanID:  testSpanDEF,
		Sampled: true,
		State:   testStateData,
	})

	headers := ExtractHeaders(ctx, nil)

	assertions.Equal(testTenantID123, headers.Get(HeaderTenantID), "TenantID header should be set")
	assertions.Equal(testMsgID456, headers.Get(HeaderMessageID), "MessageID header should be set")
	assertions.Equal(testCorrID789, headers.Get(HeaderCorrelationID), "CorrelationID header should be set")
	assertions.Equal(testTraceABC, headers.Get(HeaderTraceID), "TraceID header should be set")
	assertions.Equal(testSpanDEF, headers.Get(HeaderSpanID), "SpanID header should be set")
	assertions.Equal(testTraceABC, headers.Get(HeaderB3TraceID), "B3 TraceID header should be set")
	assertions.Equal(testSpanDEF, headers.Get(HeaderB3SpanID), "B3 SpanID header should be set")
	assertions.Equal("1", headers.Get(HeaderB3Sampled), "B3 Sampled header should be set")
	assertions.Equal(testStateData, headers.Get(HeaderTraceState), "TraceState header should be set")
	assertions.Contains(headers.Get(HeaderTraceParent), testTraceABC, "TraceParent should contain TraceID")
	assertions.Contains(headers.Get(HeaderTraceParent), testSpanDEF, "TraceParent should contain SpanID")
	assertions.Contains(headers.Get(HeaderTraceParent), "01", "TraceParent should contain sampled flag")
}

func TestExtractHeadersWithExistingHeaders(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	ctx = WithTenantID(ctx, testTenantID123)

	existingHeaders := nats.Header{
		testCustomHeader: []string{testCustomValue},
	}

	headers := ExtractHeaders(ctx, existingHeaders)

	assertions.Equal(testTenantID123, headers.Get(HeaderTenantID), "TenantID header should be set")
	assertions.Equal(testCustomValue, headers.Get(testCustomHeader), "Custom header should be preserved")
}

func TestExtractHeadersEmptyContext(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	headers := ExtractHeaders(ctx, nil)

	assertions.NotNil(headers, "Headers should not be nil")
	assertions.Empty(headers, "Headers should be empty")
}

func TestExtractHeadersTraceContextNotSampled(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	ctx = WithTraceContext(ctx, &TraceContext{
		TraceID: testTraceID123,
		SpanID:  testSpanID456,
		Sampled: false,
	})

	headers := ExtractHeaders(ctx, nil)

	assertions.Equal(testTraceID123, headers.Get(HeaderTraceID), "TraceID should be set")
	assertions.Contains(headers.Get(HeaderTraceParent), "00-"+testTraceID123+"-"+testSpanID456+"-00", "TraceParent should have sampled=00")
	assertions.Empty(headers.Get(HeaderB3Sampled), "B3 Sampled should not be set when not sampled")
}

func TestExtractHeadersTraceContextPartialIDs(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	ctx = WithTraceContext(ctx, &TraceContext{
		TraceID: testTraceOnly,
	})

	headers := ExtractHeaders(ctx, nil)

	assertions.Equal(testTraceOnly, headers.Get(HeaderTraceID), "TraceID should be set")
	assertions.Empty(headers.Get(HeaderTraceParent), "TraceParent should not be set without both IDs")
}

func TestExtractHeadersTraceContextWithState(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	ctx = WithTraceContext(ctx, &TraceContext{
		TraceID: testTraceABC,
		SpanID:  testSpanDEF,
		State:   testTraceState2V,
		Sampled: true,
	})

	headers := ExtractHeaders(ctx, nil)

	assertions.Equal(testTraceState2V, headers.Get(HeaderTraceState), "TraceState should be set")
}

/* ========================================================================================================
   EXTRACT OTEL TRACE CONTEXT TESTS
   ======================================================================================================== */

func TestExtractOTELTraceContextNoSpan(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	tc := ExtractOTELTraceContext(ctx)

	assertions.Nil(tc, "Should return nil when no OTEL span is active")
}

func TestExtractOTELTraceContextWithSpan(t *testing.T) {

	assertions := assert.New(t)

	traceID, _ := trace.TraceIDFromHex(testOTELTraceIDHex)
	spanID, _ := trace.SpanIDFromHex(testOTELSpanIDHex)

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
		Remote:     false,
	})

	ctx := trace.ContextWithSpanContext(context.Background(), spanCtx)
	tc := ExtractOTELTraceContext(ctx)

	assertions.NotNil(tc, "Should return TraceContext when OTEL span is active")
	assertions.Equal(testOTELTraceIDHex, tc.TraceID, "TraceID should match")
	assertions.Equal(testOTELSpanIDHex, tc.SpanID, "SpanID should match")
	assertions.True(tc.Sampled, "Sampled should be true for active OTEL span")
}

func TestExtractHeadersWithOTELSpan(t *testing.T) {

	assertions := assert.New(t)

	traceID, _ := trace.TraceIDFromHex(testOTELTraceIDHex)
	spanID, _ := trace.SpanIDFromHex(testOTELSpanIDHex)

	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    traceID,
		SpanID:     spanID,
		TraceFlags: trace.FlagsSampled,
		Remote:     false,
	})

	ctx := trace.ContextWithSpanContext(context.Background(), spanCtx)
	ctx = WithTenantID(ctx, testTenantOTEL)

	headers := ExtractHeaders(ctx, nil)

	assertions.Equal(testTenantOTEL, headers.Get(HeaderTenantID), "TenantID should be set")
	assertions.Equal(testOTELTraceIDHex, headers.Get(HeaderTraceID), "TraceID from OTEL should be set")
	assertions.Equal(testOTELSpanIDHex, headers.Get(HeaderSpanID), "SpanID from OTEL should be set")
	assertions.Equal("1", headers.Get(HeaderB3Sampled), "B3 Sampled should be 1")
	assertions.Contains(headers.Get(HeaderTraceParent), testOTELTraceIDHex, "TraceParent should contain trace ID")
}

/* ========================================================================================================
   INJECT CONTEXT TESTS
   ======================================================================================================== */

func TestInjectContext(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	headers := nats.Header{
		HeaderTenantID:      []string{testTenantID123},
		HeaderMessageID:     []string{testMsgID456},
		HeaderCorrelationID: []string{testCorrID789},
		HeaderTraceID:       []string{testTraceABC},
		HeaderSpanID:        []string{testSpanDEF},
		HeaderB3Sampled:     []string{"1"},
		HeaderTraceState:    []string{testStateData},
	}

	ctx = InjectContext(ctx, headers)

	tenantID, ok := TenantIDFromContext(ctx)
	assertions.True(ok, "TenantID should exist")
	assertions.Equal(testTenantID123, tenantID, "TenantID should be matched")

	messageID, ok := MessageIDFromContext(ctx)
	assertions.True(ok, "MessageID should exist")
	assertions.Equal(testMsgID456, messageID, "MessageID must match")

	correlationID, ok := CorrelationIDFromContext(ctx)
	assertions.True(ok, "CorrelationID should exist")
	assertions.Equal(testCorrID789, correlationID, "CorrelationID must match")

	tc, ok := TraceContextFromContext(ctx)
	assertions.True(ok, "TraceContext should exist")
	assertions.Equal(testTraceABC, tc.TraceID, "TraceID should match")
	assertions.Equal(testSpanDEF, tc.SpanID, "SpanID should match")
	assertions.True(tc.Sampled, "Sampled should be true")
	assertions.Equal(testStateData, tc.State, "State should match")
}

func TestInjectContextB3Headers(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	headers := nats.Header{
		HeaderB3TraceID: []string{testB3Trace123},
		HeaderB3SpanID:  []string{testB3Span456},
		HeaderB3Sampled: []string{"true"},
	}

	ctx = InjectContext(ctx, headers)

	tc, ok := TraceContextFromContext(ctx)

	assertions.True(ok, "TraceContext must exist")
	assertions.Equal(testB3Trace123, tc.TraceID, "B3 TraceID should be used")
	assertions.Equal(testB3Span456, tc.SpanID, "B3 SpanID should be used")
	assertions.True(tc.Sampled, "Sampled should be true")
}

func TestInjectContextNilHeaders(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	resultCtx := InjectContext(ctx, nil)

	assertions.Equal(ctx, resultCtx, "Context should be unchanged with nil headers")
}

func TestInjectContextEmptyHeaders(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	headers := make(nats.Header)

	resultCtx := InjectContext(ctx, headers)

	_, ok := TenantIDFromContext(resultCtx)
	assertions.False(ok, "TenantID should not exist with empty headers")
}

func TestInjectContextB3SampledVariants(t *testing.T) {

	testCases := []struct {
		name            string
		sampledValue    string
		expectedSampled bool
	}{
		{"sampled with 1", "1", true},
		{"sampled with true", "true", true},
		{"not sampled with 0", "0", false},
		{"not sampled with false", "false", false},
		{"not sampled empty", "", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			assertions := assert.New(t)

			ctx := context.Background()
			headers := nats.Header{
				HeaderB3TraceID: []string{testTraceID123},
				HeaderB3SpanID:  []string{testSpanID456},
			}
			if tc.sampledValue != "" {
				headers[HeaderB3Sampled] = []string{tc.sampledValue}
			}

			ctx = InjectContext(ctx, headers)

			traceCtx, ok := TraceContextFromContext(ctx)
			assertions.True(ok, "TraceContext should not be empty")
			assertions.Equal(tc.expectedSampled, traceCtx.Sampled, "Sampled flag should match")
		})
	}
}

func TestInjectContextCustomHeaderPrecedence(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	headers := nats.Header{
		HeaderTraceID:   []string{testCustomTrace},
		HeaderB3TraceID: []string{testB3Trace},
		HeaderSpanID:    []string{testCustomSpan},
		HeaderB3SpanID:  []string{testB3Span},
	}

	ctx = InjectContext(ctx, headers)

	tc, ok := TraceContextFromContext(ctx)
	assertions.True(ok, "TraceContext must not be empty")
	assertions.Equal(testCustomTrace, tc.TraceID, "Custom header should take precedence over B3")
	assertions.Equal(testCustomSpan, tc.SpanID, "Custom header should take precedence over B3")
}

func TestInjectContextOnlyCorrelationID(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	headers := nats.Header{
		HeaderCorrelationID: []string{testCorrOnly},
	}

	ctx = InjectContext(ctx, headers)

	correlationID, ok := CorrelationIDFromContext(ctx)
	assertions.True(ok, "CorrelationID should exist")
	assertions.Equal(testCorrOnly, correlationID, "CorrelationID must be matched")

	_, hasTrace := TraceContextFromContext(ctx)
	assertions.False(hasTrace, "No trace context should exist")
}

func TestInjectContextOnlyMessageID(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	headers := nats.Header{
		HeaderMessageID: []string{testMsgOnly},
	}

	ctx = InjectContext(ctx, headers)

	messageID, ok := MessageIDFromContext(ctx)
	assertions.True(ok, "MessageID should exist")
	assertions.Equal(testMsgOnly, messageID, "MessageID must be matched")
}

/* ========================================================================================================
   ROUNDTRIP TESTS
   ======================================================================================================== */

func TestContextRoundtrip(t *testing.T) {

	assertions := assert.New(t)

	originalCtx := context.Background()
	originalCtx = WithTenantID(originalCtx, testTenantID123)
	originalCtx = WithMessageID(originalCtx, testMsgID456)
	originalCtx = WithCorrelationID(originalCtx, testCorrID789)
	originalCtx = WithTraceContext(originalCtx, &TraceContext{
		TraceID: testTraceABC,
		SpanID:  testSpanDEF,
		Sampled: true,
		State:   testStateData,
	})

	headers := ExtractHeaders(originalCtx, nil)
	newCtx := InjectContext(context.Background(), headers)

	tenantID, _ := TenantIDFromContext(newCtx)
	assertions.Equal(testTenantID123, tenantID, "TenantID must be matched after roundtrip")

	messageID, _ := MessageIDFromContext(newCtx)
	assertions.Equal(testMsgID456, messageID, "MessageID should must match after roundtrip")

	correlationID, _ := CorrelationIDFromContext(newCtx)
	assertions.Equal(testCorrID789, correlationID, "CorrelationID should must match after roundtrip")

	tc, _ := TraceContextFromContext(newCtx)
	assertions.Equal(testTraceABC, tc.TraceID, "TraceID should match after roundtrip")
	assertions.Equal(testSpanDEF, tc.SpanID, "SpanID should match after roundtrip")
	assertions.True(tc.Sampled, "Sampled should match after roundtrip")
	assertions.Equal(testStateData, tc.State, "State should match after roundtrip")
}

/* ========================================================================================================
   MULTI-TENANT ISOLATION TESTS
   ======================================================================================================== */

func TestMultiTenantContextIsolation(t *testing.T) {

	assertions := assert.New(t)

	ctx1 := context.Background()
	ctx1 = WithTenantID(ctx1, testTenantA)
	ctx1 = WithCorrelationID(ctx1, testCorrA)

	ctx2 := context.Background()
	ctx2 = WithTenantID(ctx2, testTenantB)
	ctx2 = WithCorrelationID(ctx2, testCorrB)

	tenantA, _ := TenantIDFromContext(ctx1)
	tenantB, _ := TenantIDFromContext(ctx2)
	corrA, _ := CorrelationIDFromContext(ctx1)
	corrB, _ := CorrelationIDFromContext(ctx2)

	assertions.Equal(testTenantA, tenantA, "Tenant A should be isolated")
	assertions.Equal(testTenantB, tenantB, "Tenant B should be isolated")
	assertions.Equal(testCorrA, corrA, "Correlation A should be isolated")
	assertions.Equal(testCorrB, corrB, "Correlation B should be isolated")
}

func TestMultiTenantHeaderRoundtrip(t *testing.T) {

	tenants := []struct {
		tenantID      string
		correlationID string
		traceID       string
		spanID        string
	}{
		{testTenantAlpha, testCorrAlpha, testTraceAlpha, testSpanAlpha},
		{testTenantBeta, testCorrBeta, testTraceBeta, testSpanBeta},
		{testTenantGamma, testCorrGamma, testTraceGamma, testSpanGamma},
	}

	for _, tt := range tenants {
		t.Run(tt.tenantID, func(t *testing.T) {

			assertions := assert.New(t)

			publishCtx := context.Background()
			publishCtx = WithTenantID(publishCtx, tt.tenantID)
			publishCtx = WithCorrelationID(publishCtx, tt.correlationID)
			publishCtx = WithTraceContext(publishCtx, &TraceContext{
				TraceID: tt.traceID,
				SpanID:  tt.spanID,
				Sampled: true,
			})

			headers := ExtractHeaders(publishCtx, nil)
			subscribeCtx := InjectContext(context.Background(), headers)

			tenantID, _ := TenantIDFromContext(subscribeCtx)
			correlationID, _ := CorrelationIDFromContext(subscribeCtx)
			tc, _ := TraceContextFromContext(subscribeCtx)

			assertions.Equal(tt.tenantID, tenantID, "TenantID roundtrip should match")
			assertions.Equal(tt.correlationID, correlationID, "CorrelationID roundtrip should match")
			assertions.Equal(tt.traceID, tc.TraceID, "TraceID roundtrip should match")
			assertions.Equal(tt.spanID, tc.SpanID, "SpanID roundtrip should match")
		})
	}
}

/* ========================================================================================================
   BENCHMARK TESTS
   ======================================================================================================== */

func BenchmarkExtractHeaders(b *testing.B) {

	ctx := context.Background()
	ctx = WithTenantID(ctx, testTenantID123)
	ctx = WithMessageID(ctx, testMsgID456)
	ctx = WithCorrelationID(ctx, testCorrID789)
	ctx = WithTraceContext(ctx, &TraceContext{
		TraceID: testTraceABC,
		SpanID:  testSpanDEF,
		Sampled: true,
		State:   testStateData,
	})

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ExtractHeaders(ctx, nil)
	}
}

func BenchmarkExtractHeadersParallel(b *testing.B) {

	ctx := context.Background()
	ctx = WithTenantID(ctx, testTenantID123)
	ctx = WithMessageID(ctx, testMsgID456)
	ctx = WithCorrelationID(ctx, testCorrID789)
	ctx = WithTraceContext(ctx, &TraceContext{
		TraceID: testTraceABC,
		SpanID:  testSpanDEF,
		Sampled: true,
	})

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ExtractHeaders(ctx, nil)
		}
	})
}

func BenchmarkInjectContext(b *testing.B) {

	headers := nats.Header{
		HeaderTenantID:      []string{testTenantID123},
		HeaderMessageID:     []string{testMsgID456},
		HeaderCorrelationID: []string{testCorrID789},
		HeaderTraceID:       []string{testTraceABC},
		HeaderSpanID:        []string{testSpanDEF},
		HeaderB3Sampled:     []string{"1"},
		HeaderTraceState:    []string{testStateData},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		InjectContext(context.Background(), headers)
	}
}

func BenchmarkInjectContextParallel(b *testing.B) {

	headers := nats.Header{
		HeaderTenantID:      []string{testTenantID123},
		HeaderMessageID:     []string{testMsgID456},
		HeaderCorrelationID: []string{testCorrID789},
		HeaderTraceID:       []string{testTraceABC},
		HeaderSpanID:        []string{testSpanDEF},
		HeaderB3Sampled:     []string{"1"},
	}

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			InjectContext(context.Background(), headers)
		}
	})
}

func BenchmarkWithTenantID(b *testing.B) {

	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		WithTenantID(ctx, testTenantID123)
	}
}

func BenchmarkTenantIDFromContext(b *testing.B) {

	ctx := WithTenantID(context.Background(), testTenantID123)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		TenantIDFromContext(ctx)
	}
}

func BenchmarkContextRoundtrip(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		ctx = WithTenantID(ctx, testTenantID123)
		ctx = WithMessageID(ctx, testMsgID456)
		ctx = WithCorrelationID(ctx, testCorrID789)
		ctx = WithTraceContext(ctx, &TraceContext{
			TraceID: testTraceABC,
			SpanID:  testSpanDEF,
			Sampled: true,
		})

		headers := ExtractHeaders(ctx, nil)
		InjectContext(context.Background(), headers)
	}
}
