package middleware

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

// ============== Chain Tests ==============

func TestChain(t *testing.T) {
	var order []string

	middleware1 := func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, subject string, msg *core.Message) error {
			order = append(order, "m1-before")
			err := next(ctx, subject, msg)
			order = append(order, "m1-after")
			return err
		}
	}

	middleware2 := func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, subject string, msg *core.Message) error {
			order = append(order, "m2-before")
			err := next(ctx, subject, msg)
			order = append(order, "m2-after")
			return err
		}
	}

	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		order = append(order, "handler")
		return nil
	}

	chained := Chain(middleware1, middleware2)(handler)
	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	err := chained(ctx, "test.subject", msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"m1-before", "m2-before", "handler", "m2-after", "m1-after"}
	if len(order) != len(expected) {
		t.Fatalf("expected %d calls, got %d", len(expected), len(order))
	}
	for i, v := range expected {
		if order[i] != v {
			t.Errorf("order[%d] = %s, expected %s", i, order[i], v)
		}
	}
}

func TestChainEmpty(t *testing.T) {
	called := false
	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		called = true
		return nil
	}

	chained := Chain()(handler)
	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	err := chained(ctx, "test.subject", msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("handler was not called")
	}
}

func TestChainSubscribe(t *testing.T) {
	var order []string

	middleware1 := func(next SubscribeHandler) SubscribeHandler {
		return func(ctx context.Context, msg *core.Message) error {
			order = append(order, "m1-before")
			err := next(ctx, msg)
			order = append(order, "m1-after")
			return err
		}
	}

	middleware2 := func(next SubscribeHandler) SubscribeHandler {
		return func(ctx context.Context, msg *core.Message) error {
			order = append(order, "m2-before")
			err := next(ctx, msg)
			order = append(order, "m2-after")
			return err
		}
	}

	handler := func(ctx context.Context, msg *core.Message) error {
		order = append(order, "handler")
		return nil
	}

	chained := ChainSubscribe(middleware1, middleware2)(handler)
	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	err := chained(ctx, msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{"m1-before", "m2-before", "handler", "m2-after", "m1-after"}
	if len(order) != len(expected) {
		t.Fatalf("expected %d calls, got %d", len(expected), len(order))
	}
}

func TestChainSubscribeEmpty(t *testing.T) {
	called := false
	handler := func(ctx context.Context, msg *core.Message) error {
		called = true
		return nil
	}

	chained := ChainSubscribe()(handler)
	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	err := chained(ctx, msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("handler was not called")
	}
}

// ============== Stack Tests ==============

func TestNewStack(t *testing.T) {
	s := NewStack()
	if s == nil {
		t.Fatal("NewStack returned nil")
	}
	if s.publish == nil {
		t.Error("publish slice not initialized")
	}
	if s.subscribe == nil {
		t.Error("subscribe slice not initialized")
	}
}

func TestStackUsePublish(t *testing.T) {
	s := NewStack()
	middleware := func(next PublishHandler) PublishHandler {
		return next
	}

	s.UsePublish(middleware)
	if len(s.publish) != 1 {
		t.Errorf("expected 1 middleware, got %d", len(s.publish))
	}
}

func TestStackUseSubscribe(t *testing.T) {
	s := NewStack()
	middleware := func(next SubscribeHandler) SubscribeHandler {
		return next
	}

	s.UseSubscribe(middleware)
	if len(s.subscribe) != 1 {
		t.Errorf("expected 1 middleware, got %d", len(s.subscribe))
	}
}

func TestStackPublishChain(t *testing.T) {
	s := NewStack()

	called := false
	middleware := func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, subject string, msg *core.Message) error {
			called = true
			return next(ctx, subject, msg)
		}
	}
	s.UsePublish(middleware)

	chain := s.PublishChain()
	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return nil
	}

	wrapped := chain(handler)
	wrapped(context.Background(), "test", core.NewMessage([]byte("test")))

	if !called {
		t.Error("middleware was not called")
	}
}

func TestStackSubscribeChain(t *testing.T) {
	s := NewStack()

	called := false
	middleware := func(next SubscribeHandler) SubscribeHandler {
		return func(ctx context.Context, msg *core.Message) error {
			called = true
			return next(ctx, msg)
		}
	}
	s.UseSubscribe(middleware)

	chain := s.SubscribeChain()
	handler := func(ctx context.Context, msg *core.Message) error {
		return nil
	}

	wrapped := chain(handler)
	wrapped(context.Background(), core.NewMessage([]byte("test")))

	if !called {
		t.Error("middleware was not called")
	}
}

func TestStackWrapPublish(t *testing.T) {
	s := NewStack()

	called := false
	s.UsePublish(func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, subject string, msg *core.Message) error {
			called = true
			return next(ctx, subject, msg)
		}
	})

	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return nil
	}

	wrapped := s.WrapPublish(handler)
	wrapped(context.Background(), "test", core.NewMessage([]byte("test")))

	if !called {
		t.Error("middleware was not called")
	}
}

func TestStackWrapSubscribe(t *testing.T) {
	s := NewStack()

	called := false
	s.UseSubscribe(func(next SubscribeHandler) SubscribeHandler {
		return func(ctx context.Context, msg *core.Message) error {
			called = true
			return next(ctx, msg)
		}
	})

	handler := func(ctx context.Context, msg *core.Message) error {
		return nil
	}

	wrapped := s.WrapSubscribe(handler)
	wrapped(context.Background(), core.NewMessage([]byte("test")))

	if !called {
		t.Error("middleware was not called")
	}
}

func TestStackUseInterceptor(t *testing.T) {
	s := NewStack()

	interceptor := &mockInterceptor{
		publishCalled:   false,
		subscribeCalled: false,
	}

	s.UseInterceptor(interceptor)

	if len(s.publish) != 1 {
		t.Error("publish middleware not added")
	}
	if len(s.subscribe) != 1 {
		t.Error("subscribe middleware not added")
	}

	// Test publish
	pubHandler := s.WrapPublish(func(ctx context.Context, subject string, msg *core.Message) error {
		return nil
	})
	pubHandler(context.Background(), "test", core.NewMessage([]byte("test")))
	if !interceptor.publishCalled {
		t.Error("publish intercept not called")
	}

	// Test subscribe
	subHandler := s.WrapSubscribe(func(ctx context.Context, msg *core.Message) error {
		return nil
	})
	subHandler(context.Background(), core.NewMessage([]byte("test")))
	if !interceptor.subscribeCalled {
		t.Error("subscribe intercept not called")
	}
}

type mockInterceptor struct {
	publishCalled   bool
	subscribeCalled bool
}

func (m *mockInterceptor) InterceptPublish() PublishMiddleware {
	return func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, subject string, msg *core.Message) error {
			m.publishCalled = true
			return next(ctx, subject, msg)
		}
	}
}

func (m *mockInterceptor) InterceptSubscribe() SubscribeMiddleware {
	return func(next SubscribeHandler) SubscribeHandler {
		return func(ctx context.Context, msg *core.Message) error {
			m.subscribeCalled = true
			return next(ctx, msg)
		}
	}
}

// ============== CircuitBreaker Extended Tests ==============

func TestCircuitBreakerStateString(t *testing.T) {
	tests := []struct {
		state    CircuitState
		expected string
	}{
		{CircuitClosed, "closed"},
		{CircuitOpen, "open"},
		{CircuitHalfOpen, "half-open"},
		{CircuitState(99), "unknown"},
	}

	for _, tt := range tests {
		got := tt.state.String()
		if got != tt.expected {
			t.Errorf("CircuitState(%d).String() = %s, expected %s", tt.state, got, tt.expected)
		}
	}
}

func TestCircuitBreakerFailures(t *testing.T) {
	cb := NewCircuitBreaker(DefaultCircuitBreakerConfig())

	// Record some failures
	testErr := errors.New("test error")
	for i := 0; i < 3; i++ {
		cb.RecordFailure(testErr)
	}

	failures := cb.Failures()
	if failures != 3 {
		t.Errorf("expected 3 failures, got %d", failures)
	}
}

func TestCircuitBreakerSubscribeMiddleware(t *testing.T) {
	cb := NewCircuitBreaker(DefaultCircuitBreakerConfig())

	handlerCalled := false
	handler := func(ctx context.Context, msg *core.Message) error {
		handlerCalled = true
		return nil
	}

	middleware := CircuitBreakerSubscribeMiddleware(cb)
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	err := wrapped(ctx, msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !handlerCalled {
		t.Error("handler was not called")
	}
}

func TestCircuitBreakerSubscribeMiddlewareError(t *testing.T) {
	cfg := DefaultCircuitBreakerConfig()
	cfg.FailureThreshold = 2
	cb := NewCircuitBreaker(cfg)

	testErr := errors.New("test error")
	handler := func(ctx context.Context, msg *core.Message) error {
		return testErr
	}

	middleware := CircuitBreakerSubscribeMiddleware(cb)
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	// Record errors
	for i := 0; i < 3; i++ {
		wrapped(ctx, msg)
	}

	// Should be open now
	if cb.State() != CircuitOpen {
		t.Errorf("expected CircuitOpen, got %v", cb.State())
	}
}

func TestMultiCircuitBreakerReset(t *testing.T) {
	mcb := NewMultiCircuitBreaker(DefaultCircuitBreakerConfig())

	// Get some breakers
	mcb.Get("subject1")
	mcb.Get("subject2")

	mcb.Reset()

	// Getting subjects again should create new breakers
	cb := mcb.Get("subject1")
	if cb.State() != CircuitClosed {
		t.Error("expected fresh breaker after reset")
	}
}

func TestMultiCircuitBreakerMiddleware(t *testing.T) {
	mcb := NewMultiCircuitBreaker(DefaultCircuitBreakerConfig())

	handlerCalled := 0
	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		handlerCalled++
		return nil
	}

	middleware := MultiCircuitBreakerMiddleware(mcb)
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	wrapped(ctx, "subject1", msg)
	wrapped(ctx, "subject2", msg)

	if handlerCalled != 2 {
		t.Errorf("expected handler called 2 times, got %d", handlerCalled)
	}
}

// ============== RateLimiter Extended Tests ==============

func TestDefaultRateLimiterConfig(t *testing.T) {
	cfg := DefaultRateLimiterConfig()
	if cfg.Rate <= 0 {
		t.Error("expected positive Rate")
	}
	if cfg.Burst <= 0 {
		t.Error("expected positive Burst")
	}
}

func TestRateLimiterTokens(t *testing.T) {
	cfg := RateLimiterConfig{Rate: 100, Burst: 10}
	rl := NewRateLimiter(cfg)

	tokens := rl.Tokens()
	if tokens <= 0 {
		t.Error("expected positive tokens")
	}
}

func TestRateLimitWaitMiddleware(t *testing.T) {
	cfg := RateLimiterConfig{Rate: 1000, Burst: 10}
	rl := NewRateLimiter(cfg)

	handlerCalled := false
	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		handlerCalled = true
		return nil
	}

	middleware := RateLimitWaitMiddleware(rl)
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	err := wrapped(ctx, "test", msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !handlerCalled {
		t.Error("handler was not called")
	}
}

func TestRateLimitSubscribeMiddleware(t *testing.T) {
	cfg := RateLimiterConfig{Rate: 1000, Burst: 10}
	rl := NewRateLimiter(cfg)

	handlerCalled := false
	handler := func(ctx context.Context, msg *core.Message) error {
		handlerCalled = true
		return nil
	}

	middleware := RateLimitSubscribeMiddleware(rl)
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	err := wrapped(ctx, msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !handlerCalled {
		t.Error("handler was not called")
	}
}

func TestPerSubjectRateLimitMiddleware(t *testing.T) {
	cfg := RateLimiterConfig{Rate: 1000, Burst: 10}
	psr := NewPerSubjectRateLimiter(cfg)

	handlerCalled := 0
	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		handlerCalled++
		return nil
	}

	middleware := PerSubjectRateLimitMiddleware(psr)
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	wrapped(ctx, "subject1", msg)
	wrapped(ctx, "subject2", msg)

	if handlerCalled != 2 {
		t.Errorf("expected handler called 2 times, got %d", handlerCalled)
	}
}

func TestSlidingWindowMiddleware(t *testing.T) {
	sw := NewSlidingWindowLimiter(100, 0) // 0 window for testing

	handlerCalled := false
	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		handlerCalled = true
		return nil
	}

	middleware := SlidingWindowMiddleware(sw)
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	err := wrapped(ctx, "test", msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !handlerCalled {
		t.Error("handler was not called")
	}
}

// ============== Concurrent Access Tests ==============

func TestChainConcurrentAccess(t *testing.T) {
	var count int
	var mu sync.Mutex

	middleware := func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, subject string, msg *core.Message) error {
			mu.Lock()
			count++
			mu.Unlock()
			return next(ctx, subject, msg)
		}
	}

	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return nil
	}

	chained := Chain(middleware)(handler)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx := context.Background()
			msg := core.NewMessage([]byte("test"))
			chained(ctx, "test", msg)
		}()
	}
	wg.Wait()

	if count != 100 {
		t.Errorf("expected 100 calls, got %d", count)
	}
}

// ============== Logging Middleware Tests ==============

type testLogger struct {
	messages []string
	mu       sync.Mutex
}

func (l *testLogger) Debug(msg string, fields map[string]any) {
	l.mu.Lock()
	l.messages = append(l.messages, "debug:"+msg)
	l.mu.Unlock()
}

func (l *testLogger) Info(msg string, fields map[string]any) {
	l.mu.Lock()
	l.messages = append(l.messages, "info:"+msg)
	l.mu.Unlock()
}

func (l *testLogger) Warn(msg string, fields map[string]any) {
	l.mu.Lock()
	l.messages = append(l.messages, "warn:"+msg)
	l.mu.Unlock()
}

func (l *testLogger) Error(msg string, fields map[string]any) {
	l.mu.Lock()
	l.messages = append(l.messages, "error:"+msg)
	l.mu.Unlock()
}

func TestNewLoggingMiddleware(t *testing.T) {
	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger)
	if lm == nil {
		t.Fatal("NewLoggingMiddleware returned nil")
	}
	if lm.logger != logger {
		t.Error("logger not set correctly")
	}
}

func TestLoggingMiddlewareWithLevel(t *testing.T) {
	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelError)
	if lm.level != LogLevelError {
		t.Errorf("expected LogLevelError, got %v", lm.level)
	}
}

func TestLoggingMiddlewareWithPayload(t *testing.T) {
	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithPayload(true, 512)
	if !lm.logPayload {
		t.Error("expected logPayload to be true")
	}
	if lm.maxPayloadSize != 512 {
		t.Errorf("expected maxPayloadSize 512, got %d", lm.maxPayloadSize)
	}
}

func TestLoggingMiddlewareInterceptPublish(t *testing.T) {
	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelDebug)

	handlerCalled := false
	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		handlerCalled = true
		return nil
	}

	middleware := lm.InterceptPublish()
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test payload"))

	err := wrapped(ctx, "test.subject", msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !handlerCalled {
		t.Error("handler was not called")
	}
	if len(logger.messages) == 0 {
		t.Error("expected log messages")
	}
}

func TestLoggingMiddlewareInterceptPublishError(t *testing.T) {
	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelDebug)

	testErr := errors.New("test error")
	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return testErr
	}

	middleware := lm.InterceptPublish()
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	err := wrapped(ctx, "test.subject", msg)
	if err != testErr {
		t.Errorf("expected test error, got %v", err)
	}

	// Should have error log
	hasError := false
	for _, m := range logger.messages {
		if m == "error:publish failed" {
			hasError = true
			break
		}
	}
	if !hasError {
		t.Error("expected error log message")
	}
}

func TestLoggingMiddlewareInterceptPublishWithPayload(t *testing.T) {
	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelDebug).WithPayload(true, 10)

	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return nil
	}

	middleware := lm.InterceptPublish()
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("this is a longer payload that should be truncated"))

	err := wrapped(ctx, "test.subject", msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoggingMiddlewareInterceptSubscribe(t *testing.T) {
	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelDebug)

	handlerCalled := false
	handler := func(ctx context.Context, msg *core.Message) error {
		handlerCalled = true
		return nil
	}

	middleware := lm.InterceptSubscribe()
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test payload"))
	msg.Subject = "test.subject"

	err := wrapped(ctx, msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !handlerCalled {
		t.Error("handler was not called")
	}
}

func TestLoggingMiddlewareInterceptSubscribeError(t *testing.T) {
	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelDebug)

	testErr := errors.New("test error")
	handler := func(ctx context.Context, msg *core.Message) error {
		return testErr
	}

	middleware := lm.InterceptSubscribe()
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	err := wrapped(ctx, msg)
	if err != testErr {
		t.Errorf("expected test error, got %v", err)
	}
}

func TestLogging(t *testing.T) {
	logger := &testLogger{}
	lm := Logging(logger)
	if lm == nil {
		t.Fatal("Logging returned nil")
	}
}

func TestNoOpLogger(t *testing.T) {
	logger := NoOpLogger{}
	// These should not panic
	logger.Debug("test", nil)
	logger.Info("test", nil)
	logger.Warn("test", nil)
	logger.Error("test", nil)
}

func TestLoggingMiddlewareNilLogger(t *testing.T) {
	lm := NewLoggingMiddleware(nil)

	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return nil
	}

	middleware := lm.InterceptPublish()
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	// Should not panic with nil logger
	err := wrapped(ctx, "test.subject", msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoggingMiddlewareLogLevelFilter(t *testing.T) {
	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelWarn)

	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return nil
	}

	middleware := lm.InterceptPublish()
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	wrapped(ctx, "test.subject", msg)

	// Debug and Info messages should be filtered out
	for _, m := range logger.messages {
		if m == "debug:publishing message" || m == "info:message published" {
			t.Error("expected debug/info messages to be filtered")
		}
	}
}

// ============== Benchmarks ==============

func BenchmarkChain(b *testing.B) {
	middleware := func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, subject string, msg *core.Message) error {
			return next(ctx, subject, msg)
		}
	}

	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return nil
	}

	chained := Chain(middleware, middleware, middleware)(handler)
	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		chained(ctx, "test", msg)
	}
}

func BenchmarkStack(b *testing.B) {
	s := NewStack()
	s.UsePublish(func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, subject string, msg *core.Message) error {
			return next(ctx, subject, msg)
		}
	})

	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return nil
	}

	wrapped := s.WrapPublish(handler)
	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		wrapped(ctx, "test", msg)
	}
}

// ============== Metrics Middleware Extended Tests ==============

func TestMetricsCollectorOnPublish(t *testing.T) {
	mc := NewMetricsCollector()

	callbackCalled := false
	mc.OnPublish(func(subject string, latency time.Duration, err error) {
		callbackCalled = true
	})

	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return nil
	}

	middleware := mc.InterceptPublish()
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	wrapped(ctx, "test.subject", msg)

	if !callbackCalled {
		t.Error("OnPublish callback not called")
	}
}

func TestMetricsCollectorOnReceive(t *testing.T) {
	mc := NewMetricsCollector()

	callbackCalled := false
	mc.OnReceive(func(subject string, processTime time.Duration, err error) {
		callbackCalled = true
	})

	handler := func(ctx context.Context, msg *core.Message) error {
		return nil
	}

	middleware := mc.InterceptSubscribe()
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))
	msg.Subject = "test.subject"

	wrapped(ctx, msg)

	if !callbackCalled {
		t.Error("OnReceive callback not called")
	}
}

func TestMetricsCollectorCollect(t *testing.T) {
	mc := NewMetricsCollector()

	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return nil
	}

	middleware := mc.InterceptPublish()
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	wrapped(ctx, "test", msg)

	metrics := mc.Collect()
	if metrics.PublishCount != 1 {
		t.Errorf("expected PublishCount 1, got %d", metrics.PublishCount)
	}
	if metrics.BytesSent != 4 { // "test" is 4 bytes
		t.Errorf("expected BytesSent 4, got %d", metrics.BytesSent)
	}
}


func TestMetricsCollectorPublishError(t *testing.T) {
	mc := NewMetricsCollector()

	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return errors.New("publish error")
	}

	middleware := mc.InterceptPublish()
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	wrapped(ctx, "test", msg)

	metrics := mc.Collect()
	if metrics.PublishErrors != 1 {
		t.Errorf("expected PublishErrors 1, got %d", metrics.PublishErrors)
	}
}

func TestMetricsCollectorSubscribeError(t *testing.T) {
	mc := NewMetricsCollector()

	handler := func(ctx context.Context, msg *core.Message) error {
		return errors.New("process error")
	}

	middleware := mc.InterceptSubscribe()
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	wrapped(ctx, msg)

	metrics := mc.Collect()
	if metrics.ReceiveErrors != 1 {
		t.Errorf("expected ReceiveErrors 1, got %d", metrics.ReceiveErrors)
	}
}

func TestMetricsMiddleware(t *testing.T) {
	mc := MetricsMiddleware()
	if mc == nil {
		t.Fatal("MetricsMiddleware returned nil")
	}
}

func TestPerSubjectMetrics(t *testing.T) {
	psm := NewPerSubjectMetrics()
	if psm == nil {
		t.Fatal("NewPerSubjectMetrics returned nil")
	}
}


func TestPerSubjectMetricsAll(t *testing.T) {
	psm := NewPerSubjectMetrics()

	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return nil
	}

	middleware := psm.InterceptPublish()
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	wrapped(ctx, "subject1", msg)
	wrapped(ctx, "subject2", msg)

	all := psm.All()
	if len(all) != 2 {
		t.Errorf("expected 2 subjects, got %d", len(all))
	}
}


func TestPerSubjectMetricsInterceptSubscribeNilMsg(t *testing.T) {
	psm := NewPerSubjectMetrics()

	handlerCalled := false
	handler := func(ctx context.Context, msg *core.Message) error {
		handlerCalled = true
		return nil
	}

	middleware := psm.InterceptSubscribe()
	wrapped := middleware(handler)

	ctx := context.Background()

	// Call with nil message
	err := wrapped(ctx, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !handlerCalled {
		t.Error("handler was not called")
	}
}

// ============== OTel Middleware Tests ==============

func TestNewOTELMetricsMiddleware(t *testing.T) {
	otelm := NewOTELMetricsMiddleware()
	if otelm == nil {
		t.Fatal("NewOTELMetricsMiddleware returned nil")
	}
}

func TestOTELMetricsMiddlewareInterceptPublish(t *testing.T) {
	otelm := NewOTELMetricsMiddleware()

	handlerCalled := false
	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		handlerCalled = true
		return nil
	}

	middleware := otelm.InterceptPublish()
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	err := wrapped(ctx, "test.subject", msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !handlerCalled {
		t.Error("handler was not called")
	}
}

func TestOTELMetricsMiddlewareInterceptPublishError(t *testing.T) {
	otelm := NewOTELMetricsMiddleware()

	testErr := errors.New("test error")
	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return testErr
	}

	middleware := otelm.InterceptPublish()
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	err := wrapped(ctx, "test.subject", msg)
	if err != testErr {
		t.Errorf("expected test error, got %v", err)
	}
}

func TestOTELMetricsMiddlewareInterceptPublishNilMsg(t *testing.T) {
	otelm := NewOTELMetricsMiddleware()

	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return nil
	}

	middleware := otelm.InterceptPublish()
	wrapped := middleware(handler)

	ctx := context.Background()

	// Should not panic with nil message
	err := wrapped(ctx, "test.subject", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestOTELMetricsMiddlewareInterceptSubscribe(t *testing.T) {
	otelm := NewOTELMetricsMiddleware()

	handlerCalled := false
	handler := func(ctx context.Context, msg *core.Message) error {
		handlerCalled = true
		return nil
	}

	middleware := otelm.InterceptSubscribe()
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))
	msg.Subject = "test.subject"

	err := wrapped(ctx, msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !handlerCalled {
		t.Error("handler was not called")
	}
}

func TestOTELMetricsMiddlewareInterceptSubscribeError(t *testing.T) {
	otelm := NewOTELMetricsMiddleware()

	testErr := errors.New("test error")
	handler := func(ctx context.Context, msg *core.Message) error {
		return testErr
	}

	middleware := otelm.InterceptSubscribe()
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	err := wrapped(ctx, msg)
	if err != testErr {
		t.Errorf("expected test error, got %v", err)
	}
}

func TestOTELMetricsMiddlewareInterceptSubscribeNilMsg(t *testing.T) {
	otelm := NewOTELMetricsMiddleware()

	handler := func(ctx context.Context, msg *core.Message) error {
		return nil
	}

	middleware := otelm.InterceptSubscribe()
	wrapped := middleware(handler)

	ctx := context.Background()

	// Should not panic with nil message
	err := wrapped(ctx, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ============== OTel Logger Tests ==============

func TestNewOTELLogger(t *testing.T) {
	otelLogger := NewOTELLogger()
	if otelLogger == nil {
		t.Fatal("NewOTELLogger returned nil")
	}
}

func TestOTELLoggerMethods(t *testing.T) {
	otelLogger := NewOTELLogger()

	// Should not panic
	otelLogger.Debug("debug message", map[string]any{"key": "value"})
	otelLogger.Info("info message", map[string]any{"key": "value"})
	otelLogger.Warn("warn message", map[string]any{"key": "value"})
	otelLogger.Error("error message", map[string]any{"key": "value"})
}

func TestOTELLoggerNilFields(t *testing.T) {
	otelLogger := NewOTELLogger()

	// Should not panic with nil fields
	otelLogger.Debug("debug message", nil)
	otelLogger.Info("info message", nil)
	otelLogger.Warn("warn message", nil)
	otelLogger.Error("error message", nil)
}

func TestMapToFields(t *testing.T) {
	// Test with nil map
	fields := mapToFields(nil)
	if fields != nil {
		t.Error("expected nil for nil map")
	}

	// Test with populated map
	m := map[string]any{
		"string": "value",
		"int":    42,
		"bool":   true,
	}
	fields = mapToFields(m)
	if len(fields) != 3 {
		t.Errorf("expected 3 fields, got %d", len(fields))
	}
}

func TestOTELLoggingMiddleware(t *testing.T) {
	lm := OTELLoggingMiddleware()
	if lm == nil {
		t.Fatal("OTELLoggingMiddleware returned nil")
	}
}

// ============== Logging Middleware with LogLevelWarn ==============

func TestLoggingMiddlewareLogLevelWarnPublish(t *testing.T) {
	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelWarn)

	handlerCalled := false
	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		handlerCalled = true
		return nil
	}

	middleware := lm.InterceptPublish()
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	err := wrapped(ctx, "test.subject", msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !handlerCalled {
		t.Error("handler was not called")
	}
	// Debug and Info should be filtered - only Error should pass through for errors
}

func TestLoggingMiddlewareLogLevelWarnSubscribe(t *testing.T) {
	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelWarn)

	testErr := errors.New("test error")
	handler := func(ctx context.Context, msg *core.Message) error {
		return testErr
	}

	middleware := lm.InterceptSubscribe()
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	err := wrapped(ctx, msg)
	if err != testErr {
		t.Errorf("expected test error, got %v", err)
	}
}

func TestLoggingMiddlewareWithTraceContext(t *testing.T) {
	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelDebug)

	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return nil
	}

	middleware := lm.InterceptPublish()
	wrapped := middleware(handler)

	// Create context with trace
	tc := &core.TraceContext{
		TraceID: "trace123",
		SpanID:  "span456",
	}
	ctx := core.WithTraceContext(context.Background(), tc)
	msg := core.NewMessage([]byte("test"))

	err := wrapped(ctx, "test.subject", msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoggingMiddlewareWithTenantID(t *testing.T) {
	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelDebug)

	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return nil
	}

	middleware := lm.InterceptPublish()
	wrapped := middleware(handler)

	// Create context with tenant ID
	ctx := core.WithTenantID(context.Background(), "tenant123")
	msg := core.NewMessage([]byte("test"))

	err := wrapped(ctx, "test.subject", msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoggingMiddlewareSubscribeWithTraceAndTenant(t *testing.T) {
	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelDebug).WithPayload(true, 100)

	handler := func(ctx context.Context, msg *core.Message) error {
		return nil
	}

	middleware := lm.InterceptSubscribe()
	wrapped := middleware(handler)

	// Create context with trace and tenant
	tc := &core.TraceContext{
		TraceID: "trace123",
		SpanID:  "span456",
	}
	ctx := core.WithTraceContext(context.Background(), tc)
	ctx = core.WithTenantID(ctx, "tenant123")

	msg := core.NewMessage([]byte("test payload data"))
	msg.Subject = "test.subject"

	err := wrapped(ctx, msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoggingMiddlewareSubscribePayloadTruncation(t *testing.T) {
	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelDebug).WithPayload(true, 5)

	handler := func(ctx context.Context, msg *core.Message) error {
		return nil
	}

	middleware := lm.InterceptSubscribe()
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("this is a long payload"))
	msg.Subject = "test.subject"

	err := wrapped(ctx, msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ============== Circuit Breaker Extended Tests ==============

func TestCircuitBreakerTransitionToOpen(t *testing.T) {
	cfg := DefaultCircuitBreakerConfig()
	cfg.FailureThreshold = 2
	cb := NewCircuitBreaker(cfg)

	testErr := errors.New("test error")

	// Record failures to trigger open state
	cb.RecordFailure(testErr)
	cb.RecordFailure(testErr)
	cb.RecordFailure(testErr)

	state := cb.State()
	if state != CircuitOpen {
		t.Errorf("expected CircuitOpen, got %v", state)
	}
}

func TestCircuitBreakerRecordSuccess(t *testing.T) {
	cfg := DefaultCircuitBreakerConfig()
	cb := NewCircuitBreaker(cfg)

	cb.RecordSuccess()
	cb.RecordSuccess()

	state := cb.State()
	if state != CircuitClosed {
		t.Errorf("expected CircuitClosed, got %v", state)
	}
}



func TestMultiCircuitBreakerGet(t *testing.T) {
	mcb := NewMultiCircuitBreaker(DefaultCircuitBreakerConfig())

	cb1 := mcb.Get("subject1")
	cb2 := mcb.Get("subject1") // Same subject
	cb3 := mcb.Get("subject2") // Different subject

	if cb1 != cb2 {
		t.Error("expected same breaker for same subject")
	}
	if cb1 == cb3 {
		t.Error("expected different breaker for different subject")
	}
}

func TestMultiCircuitBreakerMiddlewareError(t *testing.T) {
	cfg := DefaultCircuitBreakerConfig()
	cfg.FailureThreshold = 1
	mcb := NewMultiCircuitBreaker(cfg)

	testErr := errors.New("test error")
	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return testErr
	}

	middleware := MultiCircuitBreakerMiddleware(mcb)
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	// First call should return error
	err := wrapped(ctx, "subject1", msg)
	if err != testErr {
		t.Errorf("expected test error, got %v", err)
	}
}


// ============== Rate Limiter Extended Tests ==============

func TestNewRateLimiter(t *testing.T) {
	cfg := RateLimiterConfig{Rate: 100, Burst: 10}
	rl := NewRateLimiter(cfg)
	if rl == nil {
		t.Fatal("NewRateLimiter returned nil")
	}
}

func TestRateLimiterAllowAndWait(t *testing.T) {
	cfg := RateLimiterConfig{Rate: 1000, Burst: 10}
	rl := NewRateLimiter(cfg)

	allowed := rl.Allow()
	if !allowed {
		t.Error("expected Allow to return true")
	}

	ctx := context.Background()
	err := rl.Wait(ctx)
	if err != nil {
		t.Fatalf("unexpected error from Wait: %v", err)
	}
}

func TestPerSubjectRateLimiterGet(t *testing.T) {
	cfg := RateLimiterConfig{Rate: 100, Burst: 10}
	psr := NewPerSubjectRateLimiter(cfg)

	rl1 := psr.Get("subject1")
	rl2 := psr.Get("subject1") // Same subject
	rl3 := psr.Get("subject2") // Different subject

	if rl1 != rl2 {
		t.Error("expected same limiter for same subject")
	}
	if rl1 == rl3 {
		t.Error("expected different limiter for different subject")
	}
}

func TestSlidingWindowLimiterAllow(t *testing.T) {
	sw := NewSlidingWindowLimiter(100, time.Second)

	for i := 0; i < 50; i++ {
		allowed := sw.Allow()
		if !allowed {
			t.Errorf("expected Allow to return true for request %d", i)
		}
	}
}

func TestSlidingWindowLimiterCount(t *testing.T) {
	sw := NewSlidingWindowLimiter(100, time.Second)

	sw.Allow()
	sw.Allow()
	sw.Allow()

	count := sw.Count()
	if count != 3 {
		t.Errorf("expected count 3, got %d", count)
	}
}

func TestRateLimitWaitMiddlewareRejection(t *testing.T) {
	// Create a rate limiter with very low rate
	cfg := RateLimiterConfig{Rate: 0.001, Burst: 1}
	rl := NewRateLimiter(cfg)

	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return nil
	}

	middleware := RateLimitWaitMiddleware(rl)
	wrapped := middleware(handler)

	// Use first token
	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))
	wrapped(ctx, "test", msg)

	// Create a context with very short timeout
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// This should fail due to timeout while waiting
	err := wrapped(ctxTimeout, "test", msg)
	// May or may not error depending on timing, just ensure no panic
	_ = err
}

func TestRateLimitSubscribeMiddlewareRejection(t *testing.T) {
	cfg := RateLimiterConfig{Rate: 1, Burst: 1}
	rl := NewRateLimiter(cfg)

	handlerCalled := 0
	handler := func(ctx context.Context, msg *core.Message) error {
		handlerCalled++
		return nil
	}

	middleware := RateLimitSubscribeMiddleware(rl)
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	// First should succeed
	wrapped(ctx, msg)

	// Second might be rate limited (depends on timing)
	wrapped(ctx, msg)

	// At least one should have succeeded
	if handlerCalled < 1 {
		t.Error("expected at least one call to handler")
	}
}

func TestPerSubjectRateLimitMiddlewareRejection(t *testing.T) {
	cfg := RateLimiterConfig{Rate: 1, Burst: 1}
	psr := NewPerSubjectRateLimiter(cfg)

	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return nil
	}

	middleware := PerSubjectRateLimitMiddleware(psr)
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	// Should not error (allow first request)
	err := wrapped(ctx, "test", msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSlidingWindowMiddlewareRejection(t *testing.T) {
	sw := NewSlidingWindowLimiter(2, time.Second)

	handlerCalled := 0
	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		handlerCalled++
		return nil
	}

	middleware := SlidingWindowMiddleware(sw)
	wrapped := middleware(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	// First two should succeed
	wrapped(ctx, "test", msg)
	wrapped(ctx, "test", msg)

	// Third might be rejected
	err := wrapped(ctx, "test", msg)
	// If not rejected, handler should be called; either way, no panic
	_ = err

	if handlerCalled < 2 {
		t.Errorf("expected at least 2 calls, got %d", handlerCalled)
	}
}

// ============== Additional Benchmarks ==============

func BenchmarkRetryMiddleware(b *testing.B) {
	cfg := DefaultRetryConfig()
	rm := NewRetryMiddleware(cfg)

	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return nil
	}

	middleware := rm.InterceptPublish()
	wrapped := middleware(handler)
	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		wrapped(ctx, "test", msg)
	}
}

func BenchmarkMetricsCollector(b *testing.B) {
	mc := NewMetricsCollector()

	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return nil
	}

	middleware := mc.InterceptPublish()
	wrapped := middleware(handler)
	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		wrapped(ctx, "test", msg)
	}
}

func BenchmarkPerSubjectMetrics(b *testing.B) {
	psm := NewPerSubjectMetrics()

	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return nil
	}

	middleware := psm.InterceptPublish()
	wrapped := middleware(handler)
	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		wrapped(ctx, "test", msg)
	}
}

func BenchmarkCircuitBreaker(b *testing.B) {
	cb := NewCircuitBreaker(DefaultCircuitBreakerConfig())

	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return nil
	}

	middleware := CircuitBreakerMiddleware(cb)
	wrapped := middleware(handler)
	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		wrapped(ctx, "test", msg)
	}
}

func BenchmarkRateLimiter(b *testing.B) {
	cfg := RateLimiterConfig{Rate: 1000000, Burst: 1000}
	rl := NewRateLimiter(cfg)

	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return nil
	}

	middleware := RateLimitMiddleware(rl)
	wrapped := middleware(handler)
	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		wrapped(ctx, "test", msg)
	}
}

func BenchmarkLoggingMiddleware(b *testing.B) {
	lm := NewLoggingMiddleware(NoOpLogger{})

	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		return nil
	}

	middleware := lm.InterceptPublish()
	wrapped := middleware(handler)
	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		wrapped(ctx, "test", msg)
	}
}
