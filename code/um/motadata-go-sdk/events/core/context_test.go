package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

/* ========================================================================================================
   TENANT ID CONTEXT TESTS
   ======================================================================================================== */

func TestWithTenantID(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	ctx = WithTenantID(ctx, "tenant-123")

	tenantID, ok := TenantIDFromContext(ctx)

	assertions.True(ok, "TenantID should exist in context")
	assertions.Equal("tenant-123", tenantID, "TenantID should match")
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
		TraceID:  "trace-123",
		SpanID:   "span-456",
		ParentID: "parent-789",
		Sampled:  true,
		State:    "state-data",
	}

	ctx = WithTraceContext(ctx, tc)

	retrieved, ok := TraceContextFromContext(ctx)

	assertions.True(ok, "TraceContext should exist in context")
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
	ctx = WithMessageID(ctx, "msg-123")

	messageID, ok := MessageIDFromContext(ctx)

	assertions.True(ok, "MessageID should exist in context")
	assertions.Equal("msg-123", messageID, "MessageID should match")
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
	ctx = WithCorrelationID(ctx, "corr-123")

	correlationID, ok := CorrelationIDFromContext(ctx)

	assertions.True(ok, "CorrelationID should exist in context")
	assertions.Equal("corr-123", correlationID, "CorrelationID should match")
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
		TenantID: "tenant-123",
		Stream:   "test-stream",
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

/* ========================================================================================================
   EXTRACT HEADERS TESTS
   ======================================================================================================== */

func TestExtractHeaders(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	ctx = WithTenantID(ctx, "tenant-123")
	ctx = WithMessageID(ctx, "msg-456")
	ctx = WithCorrelationID(ctx, "corr-789")
	ctx = WithTraceContext(ctx, &TraceContext{
		TraceID: "trace-abc",
		SpanID:  "span-def",
		Sampled: true,
		State:   "state-data",
	})

	headers := ExtractHeaders(ctx, nil)

	assertions.Equal("tenant-123", headers.Get(HeaderTenantID), "TenantID header should be set")
	assertions.Equal("msg-456", headers.Get(HeaderMessageID), "MessageID header should be set")
	assertions.Equal("corr-789", headers.Get(HeaderCorrelationID), "CorrelationID header should be set")
	assertions.Equal("trace-abc", headers.Get(HeaderTraceID), "TraceID header should be set")
	assertions.Equal("span-def", headers.Get(HeaderSpanID), "SpanID header should be set")
	assertions.Equal("trace-abc", headers.Get(HeaderB3TraceID), "B3 TraceID header should be set")
	assertions.Equal("span-def", headers.Get(HeaderB3SpanID), "B3 SpanID header should be set")
	assertions.Equal("1", headers.Get(HeaderB3Sampled), "B3 Sampled header should be set")
	assertions.Equal("state-data", headers.Get(HeaderTraceState), "TraceState header should be set")
	assertions.Contains(headers.Get(HeaderTraceParent), "trace-abc", "TraceParent should contain TraceID")
	assertions.Contains(headers.Get(HeaderTraceParent), "span-def", "TraceParent should contain SpanID")
	assertions.Contains(headers.Get(HeaderTraceParent), "01", "TraceParent should contain sampled flag")
}

func TestExtractHeadersWithExistingHeaders(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	ctx = WithTenantID(ctx, "tenant-123")

	existingHeaders := Headers{
		"Custom-Header": []string{"custom-value"},
	}

	headers := ExtractHeaders(ctx, existingHeaders)

	assertions.Equal("tenant-123", headers.Get(HeaderTenantID), "TenantID header should be set")
	assertions.Equal("custom-value", headers.Get("Custom-Header"), "Custom header should be preserved")
}

func TestExtractHeadersEmptyContext(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	headers := ExtractHeaders(ctx, nil)

	assertions.NotNil(headers, "Headers should not be nil")
	assertions.Empty(headers, "Headers should be empty")
}

/* ========================================================================================================
   INJECT CONTEXT TESTS
   ======================================================================================================== */

func TestInjectContext(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	headers := Headers{
		HeaderTenantID:      []string{"tenant-123"},
		HeaderMessageID:     []string{"msg-456"},
		HeaderCorrelationID: []string{"corr-789"},
		HeaderTraceID:       []string{"trace-abc"},
		HeaderSpanID:        []string{"span-def"},
		HeaderB3Sampled:     []string{"1"},
		HeaderTraceState:    []string{"state-data"},
	}

	ctx = InjectContext(ctx, headers)

	tenantID, ok := TenantIDFromContext(ctx)
	assertions.True(ok, "TenantID should exist")
	assertions.Equal("tenant-123", tenantID, "TenantID should match")

	messageID, ok := MessageIDFromContext(ctx)
	assertions.True(ok, "MessageID should exist")
	assertions.Equal("msg-456", messageID, "MessageID should match")

	correlationID, ok := CorrelationIDFromContext(ctx)
	assertions.True(ok, "CorrelationID should exist")
	assertions.Equal("corr-789", correlationID, "CorrelationID should match")

	tc, ok := TraceContextFromContext(ctx)
	assertions.True(ok, "TraceContext should exist")
	assertions.Equal("trace-abc", tc.TraceID, "TraceID should match")
	assertions.Equal("span-def", tc.SpanID, "SpanID should match")
	assertions.True(tc.Sampled, "Sampled should be true")
	assertions.Equal("state-data", tc.State, "State should match")
}

func TestInjectContextB3Headers(t *testing.T) {

	assertions := assert.New(t)

	ctx := context.Background()
	headers := Headers{
		HeaderB3TraceID: []string{"b3-trace-123"},
		HeaderB3SpanID:  []string{"b3-span-456"},
		HeaderB3Sampled: []string{"true"},
	}

	ctx = InjectContext(ctx, headers)

	tc, ok := TraceContextFromContext(ctx)

	assertions.True(ok, "TraceContext should exist")
	assertions.Equal("b3-trace-123", tc.TraceID, "B3 TraceID should be used")
	assertions.Equal("b3-span-456", tc.SpanID, "B3 SpanID should be used")
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
	headers := make(Headers)

	resultCtx := InjectContext(ctx, headers)

	_, ok := TenantIDFromContext(resultCtx)
	assertions.False(ok, "TenantID should not exist with empty headers")
}

/* ========================================================================================================
   ROUNDTRIP TESTS
   ======================================================================================================== */

func TestContextRoundtrip(t *testing.T) {

	assertions := assert.New(t)

	// Create context with all values
	originalCtx := context.Background()
	originalCtx = WithTenantID(originalCtx, "tenant-123")
	originalCtx = WithMessageID(originalCtx, "msg-456")
	originalCtx = WithCorrelationID(originalCtx, "corr-789")
	originalCtx = WithTraceContext(originalCtx, &TraceContext{
		TraceID: "trace-abc",
		SpanID:  "span-def",
		Sampled: true,
		State:   "state-data",
	})

	// Extract to headers
	headers := ExtractHeaders(originalCtx, nil)

	// Inject back to new context
	newCtx := InjectContext(context.Background(), headers)

	// Verify values match
	tenantID, _ := TenantIDFromContext(newCtx)
	assertions.Equal("tenant-123", tenantID, "TenantID should match after roundtrip")

	messageID, _ := MessageIDFromContext(newCtx)
	assertions.Equal("msg-456", messageID, "MessageID should match after roundtrip")

	correlationID, _ := CorrelationIDFromContext(newCtx)
	assertions.Equal("corr-789", correlationID, "CorrelationID should match after roundtrip")

	tc, _ := TraceContextFromContext(newCtx)
	assertions.Equal("trace-abc", tc.TraceID, "TraceID should match after roundtrip")
	assertions.Equal("span-def", tc.SpanID, "SpanID should match after roundtrip")
	assertions.True(tc.Sampled, "Sampled should match after roundtrip")
	assertions.Equal("state-data", tc.State, "State should match after roundtrip")
}
