package tracer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/config"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Link is an alias for trace.Link
type Link = trace.Link

// Status code constants
var (
	CodeOK    = codes.Ok
	CodeError = codes.Error
)

/* ---------------------------------------- Initialization Tests -------------------------------------------------- */

// TestDefaultConfig tests default configuration
func TestDefaultConfig(t *testing.T) {
	assertions := assert.New(t)

	config := DefaultConfig()

	assertions.Equal("app", config.ServiceName)
	assertions.Equal("0.0.0", config.ServiceVersion)
	assertions.Equal("development", config.Environment)
}

// TestInit tests tracer initialization
func TestInit(t *testing.T) {
	assertions := assert.New(t)

	// Cleanup
	Shutdown(context.Background())

	config := DefaultConfig()
	config.Enabled = false // Don't actually connect to OTEL

	tracerInstance, err := Init(config)
	assertions.NoError(err)
	assertions.Nil(tracerInstance) // Returns nil when disabled
}

// TestInitEnabled tests tracer initialization when enabled
func TestInitEnabled(t *testing.T) {
	// Skip if can't connect to OTEL endpoint
	t.Skip("Requires OTEL collector")
}

/* ---------------------------------------- Span Creation Tests -------------------------------------------------- */

// TestSpanStart tests starting a span
func TestSpanStart(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	ctx, span := Start(ctx, "test-span")

	assertions.NotNil(span)
	assertions.NotNil(ctx)

	span.End()
}

// TestSpanWithKind tests creating spans with different kinds
func TestSpanWithKind(t *testing.T) {
	testCases := []struct {
		name    string
		kind    SpanKind
		startFn func(context.Context, string, ...StartSpanOption) (context.Context, Span)
	}{
		{"Internal", SpanKindInternal, Start},
		{"Server", SpanKindServer, StartServer},
		{"Client", SpanKindClient, StartClient},
		{"Producer", SpanKindProducer, StartProducer},
		{"Consumer", SpanKindConsumer, StartConsumer},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)

			ctx := context.Background()
			_, span := tc.startFn(ctx, fmt.Sprintf("test-%s-span", tc.name))

			assertions.NotNil(span)
			span.End()
		})
	}
}

// TestSpanWithAttributes tests adding attributes to spans
func TestSpanWithAttributes(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	_, span := Start(ctx, "test-span",
		WithAttributes(
			StringAttr("key1", "value1"),
			IntAttr("key2", 42),
			BoolAttr("key3", true),
		),
	)

	assertions.NotNil(span)
	span.End()
}

// TestSpanSetAttributes tests setting attributes after span creation
func TestSpanSetAttributes(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	_, span := Start(ctx, "test-span")

	span.SetAttributes(
		StringAttr("key1", "value1"),
		Float64Attr("key2", 3.14),
	)

	assertions.True(true) // Verify no panic
	span.End()
}

/* ---------------------------------------- Span Status Tests -------------------------------------------------- */

// TestSpanSetOK tests setting span status to OK
func TestSpanSetOK(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	_, span := Start(ctx, "test-span")

	span.SetOK()
	assertions.True(true) // No panic
	span.End()
}

// TestSpanSetError tests setting span error status
func TestSpanSetError(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	_, span := Start(ctx, "test-span")

	testError := errors.New("test error")
	span.SetError(testError)
	assertions.True(true) // No panic
	span.End()
}

// TestSpanSetErrorNil tests setting nil error
func TestSpanSetErrorNil(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	_, span := Start(ctx, "test-span")

	// Should not panic with nil error
	assertions.NotPanics(func() {
		span.SetError(nil)
	})
	span.End()
}

// TestSpanRecordError tests recording error as event
func TestSpanRecordError(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	_, span := Start(ctx, "test-span")

	testError := errors.New("recorded error")
	span.RecordError(testError)
	assertions.True(true) // No panic
	span.End()
}

/* ---------------------------------------- Span Events Tests -------------------------------------------------- */

// TestSpanAddEvent tests adding events to spans
func TestSpanAddEvent(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	_, span := Start(ctx, "test-span")

	span.AddEvent("event-1")
	span.AddEvent("event-2")
	assertions.True(true) // No panic
	span.End()
}

// TestSpanSetName tests renaming a span
func TestSpanSetName(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	_, span := Start(ctx, "original-name")

	span.SetName("new-name")
	assertions.True(true) // No panic
	span.End()
}

/* ---------------------------------------- Context Propagation Tests -------------------------------------------------- */

// TestSpanFromContext tests extracting span from context
func TestSpanFromContext(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	ctx, parentSpan := Start(ctx, "parent-span")

	// Extract span from context
	extractedSpan := SpanFromContext(ctx)
	assertions.NotNil(extractedSpan)

	// Span context should match
	assertions.Equal(
		parentSpan.SpanContext().TraceID(),
		extractedSpan.SpanContext().TraceID(),
	)

	parentSpan.End()
}

// TestContextWithSpan tests adding span to context
func TestContextWithSpan(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	_, span := Start(ctx, "test-span")

	newCtx := ContextWithSpan(ctx, span)
	assertions.NotNil(newCtx)

	extractedSpan := SpanFromContext(newCtx)
	assertions.NotNil(extractedSpan)

	span.End()
}

// TestTraceIDFromContext tests extracting trace ID
func TestTraceIDFromContext(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	ctx, span := Start(ctx, "test-span")

	traceID := TraceIDFromContext(ctx)
	// NoopTracer returns empty trace ID
	assertions.NotNil(traceID)

	span.End()
}

// TestSpanIDFromContext tests extracting span ID
func TestSpanIDFromContext(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	ctx, span := Start(ctx, "test-span")

	spanID := SpanIDFromContext(ctx)
	assertions.NotNil(spanID)

	span.End()
}

// TestNestedSpans tests creating nested spans
func TestNestedSpans(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()

	ctx, parentSpan := Start(ctx, "parent")
	assertions.NotNil(parentSpan)

	ctx, childSpan := Start(ctx, "child")
	assertions.NotNil(childSpan)

	_, grandchildSpan := Start(ctx, "grandchild")
	assertions.NotNil(grandchildSpan)

	grandchildSpan.End()
	childSpan.End()
	parentSpan.End()
}

/* ---------------------------------------- WithSpan Helper Tests -------------------------------------------------- */

// TestWithSpan tests the WithSpan helper
func TestWithSpan(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	executed := false

	err := WithSpan(ctx, "test-operation", func(ctx context.Context) error {
		executed = true
		return nil
	})

	assertions.NoError(err)
	assertions.True(executed)
}

// TestWithSpanError tests WithSpan with error
func TestWithSpanError(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	expectedError := errors.New("operation failed")

	err := WithSpan(ctx, "failing-operation", func(ctx context.Context) error {
		return expectedError
	})

	assertions.Error(err)
	assertions.Equal(expectedError, err)
}

// TestTimeSpan tests the TimeSpan helper
func TestTimeSpan(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	executed := false

	duration := TimeSpan(ctx, "timed-operation", func(ctx context.Context) {
		executed = true
	})

	assertions.True(executed)
	assertions.GreaterOrEqual(duration.Nanoseconds(), int64(0))
}

/* ---------------------------------------- Span Wrapper Tests -------------------------------------------------- */

// TestSpanIsRecording tests IsRecording method
func TestSpanIsRecording(t *testing.T) {
	ctx := context.Background()
	_, span := Start(ctx, "test-span")

	// Just verify no panic
	_ = span.IsRecording()
	span.End()
}

// TestSpanUnwrap tests Unwrap method
func TestSpanUnwrap(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	_, span := Start(ctx, "test-span")

	underlyingSpan := span.Unwrap()
	assertions.NotNil(underlyingSpan)
	span.End()
}

// TestSpanTracerProvider tests TracerProvider method
func TestSpanTracerProvider(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	_, span := Start(ctx, "test-span")

	provider := span.TracerProvider()
	assertions.NotNil(provider)
	span.End()
}

/* ---------------------------------------- Concurrent Safety Tests -------------------------------------------------- */

// TestConcurrentSpanCreation tests concurrent span creation
func TestConcurrentSpanCreation(t *testing.T) {
	ctx := context.Background()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_, span := Start(ctx, fmt.Sprintf("span-%d-%d", id, j))
				span.SetAttributes(IntAttr("id", id))
				span.End()
			}
		}(i)
	}

	wg.Wait()
}

// TestConcurrentNestedSpans tests concurrent nested span creation
func TestConcurrentNestedSpans(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ctx := context.Background()

			for j := 0; j < 50; j++ {
				ctx, parent := Start(ctx, fmt.Sprintf("parent-%d-%d", id, j))
				_, child := Start(ctx, fmt.Sprintf("child-%d-%d", id, j))
				child.End()
				parent.End()
			}
		}(i)
	}

	wg.Wait()
}

// TestConcurrentSpanOperations tests concurrent operations on same span
func TestConcurrentSpanOperations(t *testing.T) {
	ctx := context.Background()
	_, span := Start(ctx, "shared-span")

	var wg sync.WaitGroup

	// Multiple goroutines setting attributes
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			span.SetAttributes(IntAttr(fmt.Sprintf("key-%d", id), id))
		}(i)
	}

	// Multiple goroutines adding events
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			span.AddEvent(fmt.Sprintf("event-%d", id))
		}(i)
	}

	wg.Wait()
	span.End()
}

/* ---------------------------------------- Edge Case Tests -------------------------------------------------- */

// TestEmptySpanName tests creating span with empty name
func TestEmptySpanName(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()

	// Should not panic
	assertions.NotPanics(func() {
		_, span := Start(ctx, "")
		span.End()
	})
}

// TestNilContext tests operations with nil context
// Note: nil context is expected to panic as per Go context package behavior
func TestNilContext(t *testing.T) {
	assertions := assert.New(t)

	// Nil context is expected to panic (this is standard Go behavior)
	assertions.Panics(func() {
		_, span := Start(nil, "test-span") //nolint:staticcheck // intentional nil context test
		span.End()
	}, "nil context should panic as per Go context package contract")
}

// TestManyAttributes tests span with many attributes
func TestManyAttributes(t *testing.T) {
	ctx := context.Background()
	_, span := Start(ctx, "many-attributes")

	// Add 200 attributes (above typical SDK limit of 128)
	for i := 0; i < 200; i++ {
		span.SetAttributes(IntAttr(fmt.Sprintf("attr-%d", i), i))
	}

	span.End()
}

// TestLongAttributeValues tests span with very long attribute values
func TestLongAttributeValues(t *testing.T) {
	ctx := context.Background()
	_, span := Start(ctx, "long-attributes")

	// Create a very long string (10KB)
	longValue := make([]byte, 10*1024)
	for i := range longValue {
		longValue[i] = 'a'
	}

	span.SetAttributes(StringAttr("long-key", string(longValue)))
	span.End()
}

// TestSpanDoubleEnd tests ending a span twice
func TestSpanDoubleEnd(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	_, span := Start(ctx, "double-end")

	span.End()

	// Second End should not panic
	assertions.NotPanics(func() {
		span.End()
	})
}

/* ---------------------------------------- Attribute Constructor Tests -------------------------------------------------- */

// TestAttributeConstructors tests all attribute constructors
func TestAttributeConstructors(t *testing.T) {
	assertions := assert.New(t)

	// String
	strAttr := StringAttr("string-key", "string-value")
	assertions.NotNil(strAttr)

	// Int
	intAttr := IntAttr("int-key", 42)
	assertions.NotNil(intAttr)

	// Int64
	int64Attr := Int64Attr("int64-key", 9223372036854775807)
	assertions.NotNil(int64Attr)

	// Float64
	float64Attr := Float64Attr("float64-key", 3.14159)
	assertions.NotNil(float64Attr)

	// Bool
	boolAttr := BoolAttr("bool-key", true)
	assertions.NotNil(boolAttr)

	// String slice
	stringsAttr := StringsAttr("strings-key", []string{"a", "b", "c"})
	assertions.NotNil(stringsAttr)

	// Int slice
	intsAttr := IntsAttr("ints-key", []int{1, 2, 3})
	assertions.NotNil(intsAttr)

	// Int64 slice
	int64sAttr := Int64sAttr("int64s-key", []int64{1, 2, 3})
	assertions.NotNil(int64sAttr)

	// Float64 slice
	float64sAttr := Float64sAttr("float64s-key", []float64{1.1, 2.2, 3.3})
	assertions.NotNil(float64sAttr)

	// Bool slice
	boolsAttr := BoolsAttr("bools-key", []bool{true, false, true})
	assertions.NotNil(boolsAttr)
}

/* ---------------------------------------- Global Tracer Tests -------------------------------------------------- */

// TestT tests the global tracer accessor
func TestT(t *testing.T) {
	assertions := assert.New(t)

	// Clean up first
	_ = Shutdown(context.Background())

	// T() should return a tracer even without initialization
	tracer := T()
	assertions.NotNil(tracer, "T() should return a non-nil tracer")
}

// TestShutdownGlobal tests global shutdown
func TestShutdownGlobal(t *testing.T) {
	assertions := assert.New(t)

	// Clean up first
	_ = Shutdown(context.Background())

	config := DefaultConfig()
	config.Enabled = false

	_, initErr := Init(config)
	assertions.NoError(initErr)

	err := Shutdown(context.Background())
	assertions.NoError(err, "Shutdown should not error")

	// Second shutdown should be idempotent
	err = Shutdown(context.Background())
	assertions.NoError(err, "Second shutdown should not error")
}

// TestShutdownWithoutInit tests shutdown without initialization
func TestShutdownWithoutInit(t *testing.T) {
	assertions := assert.New(t)

	// Clean up any existing state
	_ = Shutdown(context.Background())

	// Shutdown without init should not error
	err := Shutdown(context.Background())
	assertions.NoError(err)
}

/* ---------------------------------------- Tracer Instance Tests -------------------------------------------------- */

// TestTracerInstanceShutdown tests Tracer instance shutdown
func TestTracerInstanceShutdown(t *testing.T) {
	assertions := assert.New(t)

	// Create a tracer with disabled config
	config := DefaultConfig()
	config.Enabled = false

	tracerInstance, err := Init(config)
	assertions.NoError(err)
	assertions.Nil(tracerInstance) // Returns nil when disabled

	// Clean up
	_ = Shutdown(context.Background())
}

/* ---------------------------------------- Config Loading Tests -------------------------------------------------- */

// TestConfigFromEnv tests loading config from environment
func TestConfigFromEnv(t *testing.T) {
	assertions := assert.New(t)

	// Set environment variables
	t.Setenv("TRACER_ENABLED", "false")
	t.Setenv("TRACER_SERVICE_NAME", "test-service")

	config := ConfigFromEnv()
	assertions.NotNil(config)
}

// TestInitFromEnv tests initialization from environment
func TestInitFromEnv(t *testing.T) {
	assertions := assert.New(t)

	// Clean up first
	_ = Shutdown(context.Background())

	// Disable tracer via env
	t.Setenv("TRACER_ENABLED", "false")

	tracerInstance, err := InitFromEnv()
	assertions.NoError(err)
	assertions.Nil(tracerInstance) // Disabled returns nil
}

/* ---------------------------------------- WithLinks Option Tests -------------------------------------------------- */

// TestWithLinks tests adding links to spans
func TestWithLinks(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()

	// Create a parent span to link to
	ctx, parentSpan := Start(ctx, "parent-span")
	parentSpanContext := parentSpan.SpanContext()
	parentSpan.End()

	// Create a new span with link to parent
	_, linkedSpan := Start(ctx, "linked-span",
		WithLinks(Link{SpanContext: parentSpanContext}),
	)

	assertions.NotNil(linkedSpan)
	linkedSpan.End()
}

// TestWithMultipleLinks tests adding multiple links
func TestWithMultipleLinks(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()

	// Create spans to link to
	_, span1 := Start(ctx, "span-1")
	sc1 := span1.SpanContext()
	span1.End()

	_, span2 := Start(ctx, "span-2")
	sc2 := span2.SpanContext()
	span2.End()

	// Create span with multiple links
	_, linkedSpan := Start(ctx, "multi-linked-span",
		WithLinks(
			Link{SpanContext: sc1},
			Link{SpanContext: sc2},
		),
	)

	assertions.NotNil(linkedSpan)
	linkedSpan.End()
}

/* ---------------------------------------- SetStatus Tests -------------------------------------------------- */

// TestSpanSetStatus tests setting span status directly
func TestSpanSetStatus(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	_, span := Start(ctx, "status-span")

	// Test setting status with description
	assertions.NotPanics(func() {
		span.SetStatus(CodeOK, "success")
	})

	assertions.NotPanics(func() {
		span.SetStatus(CodeError, "something went wrong")
	})

	span.End()
}

/* ---------------------------------------- ContextWithSpan Edge Cases -------------------------------------------------- */

// TestContextWithSpanNonWrapper tests ContextWithSpan with non-wrapper span
func TestContextWithSpanNonWrapper(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	_, span := Start(ctx, "test-span")

	// Using a proper span wrapper should work
	newCtx := ContextWithSpan(ctx, span)
	assertions.NotNil(newCtx)

	span.End()
}

/* ---------------------------------------- TraceID and SpanID Edge Cases -------------------------------------------------- */

// TestTraceIDFromEmptyContext tests TraceIDFromContext with no span
func TestTraceIDFromEmptyContext(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	traceID := TraceIDFromContext(ctx)

	// Empty context should return empty string
	assertions.Equal("", traceID)
}

// TestSpanIDFromEmptyContext tests SpanIDFromContext with no span
func TestSpanIDFromEmptyContext(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	spanID := SpanIDFromContext(ctx)

	// Empty context should return empty string
	assertions.Equal("", spanID)
}

/* ---------------------------------------- WithSpanKind Option Tests -------------------------------------------------- */

// TestWithSpanKindOption tests the WithSpanKind option directly
func TestWithSpanKindOption(t *testing.T) {
	assertions := assert.New(t)

	testCases := []struct {
		name string
		kind SpanKind
	}{
		{"Internal", SpanKindInternal},
		{"Server", SpanKindServer},
		{"Client", SpanKindClient},
		{"Producer", SpanKindProducer},
		{"Consumer", SpanKindConsumer},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			_, span := Start(ctx, "test-span", WithSpanKind(tc.kind))
			assertions.NotNil(span)
			span.End()
		})
	}
}

/* ---------------------------------------- Span Options Combination Tests -------------------------------------------------- */

// TestCombinedSpanOptions tests combining multiple span options
func TestCombinedSpanOptions(t *testing.T) {
	assertions := assert.New(t)

	ctx := context.Background()
	_, span := Start(ctx, "combined-options-span",
		WithSpanKind(SpanKindServer),
		WithAttributes(
			StringAttr("key1", "value1"),
			IntAttr("key2", 42),
		),
	)

	assertions.NotNil(span)
	span.End()
}

/* ---------------------------------------- Global Start Functions Tests -------------------------------------------------- */

// TestGlobalStartFunctions tests all global Start* functions
func TestGlobalStartFunctions(t *testing.T) {
	testCases := []struct {
		name    string
		startFn func(context.Context, string, ...StartSpanOption) (context.Context, Span)
	}{
		{"Start", Start},
		{"StartServer", StartServer},
		{"StartClient", StartClient},
		{"StartProducer", StartProducer},
		{"StartConsumer", StartConsumer},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			ctx := context.Background()

			newCtx, span := tc.startFn(ctx, "test-"+tc.name)
			assertions.NotNil(newCtx)
			assertions.NotNil(span)
			span.End()
		})
	}
}

/* ---------------------------------------- createPropagator Tests -------------------------------------------------- */

// TestCreatePropagatorDefault tests default propagator creation
func TestCreatePropagatorDefault(t *testing.T) {
	assertions := assert.New(t)

	// Test with empty names (should use defaults)
	propagator := createPropagator(nil)
	assertions.NotNil(propagator)

	propagator = createPropagator([]string{})
	assertions.NotNil(propagator)
}

// TestCreatePropagatorTypes tests various propagator types
func TestCreatePropagatorTypes(t *testing.T) {
	testCases := []struct {
		name        string
		propagators []string
	}{
		{"tracecontext", []string{"tracecontext"}},
		{"traceparent", []string{"traceparent"}},
		{"baggage", []string{"baggage"}},
		{"b3", []string{"b3"}},
		{"b3multi", []string{"b3multi"}},
		{"jaeger", []string{"jaeger"}},
		{"multiple", []string{"tracecontext", "baggage", "b3"}},
		{"all", []string{"tracecontext", "baggage", "b3", "b3multi", "jaeger"}},
		{"unknown_returns_default", []string{"unknown"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			propagator := createPropagator(tc.propagators)
			assertions.NotNil(propagator)
		})
	}
}

/* ---------------------------------------- InitFromUnifiedConfig Tests -------------------------------------------------- */

// TestInitFromUnifiedConfig tests initialization from unified config
func TestInitFromUnifiedConfig(t *testing.T) {
	assertions := assert.New(t)

	// Clean up first
	_ = Shutdown(context.Background())

	// Create a unified config with tracer disabled
	cfg := config.DefaultConfig()
	cfg.Tracer.Enabled = false

	tracerInstance, err := InitFromUnifiedConfig(cfg)
	assertions.NoError(err)
	assertions.Nil(tracerInstance) // Disabled returns nil
}
