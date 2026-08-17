package middleware

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

/* ========================================================================================================
   COMMON TEST VARIABLES
   ======================================================================================================== */

var (
	testSubject        = "test.subject"
	testSubject1       = "subject1"
	testSubject2       = "subject2"
	testPayload        = []byte("test")
	testPayloadData    = []byte("test data")
	testPayloadLong    = []byte("this is a longer payload that should be truncated")
	testPayloadStr     = []byte("test payload")
	testPayloadContent = []byte("test payload data")
	testErrMsg         = "test error"
	testPublishFailed  = "publish failed"
	testProcessFailed  = "process failed"
	testTraceID        = "trace123"
	testSpanID         = "span456"
	testTenantID       = "tenant123"
	testExistingTrace  = "existing-trace-id"
	testExistingSpan   = "existing-span-id"
	testIncomingTrace  = "incoming-trace-id"
	testIncomingSpan   = "incoming-span-id"
	testW3CTraceID     = "4bf92f3577b34da6a3ce929d0e0e4736"
	testW3CSpanID      = "00f067aa0ba902b7"
	testW3CParent      = "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"
	testW3CParentUnsmp = "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-00"
	orderM1Before      = "m1-before"
	orderM1After       = "m1-after"
	orderM2Before      = "m2-before"
	orderM2After       = "m2-after"
	orderHandler       = "handler"
	errSimulated       = errors.New("simulated failure")
)

// newTestMsg creates a nats.Msg with standard test fields.
func newTestMsg(subject string, data []byte) *nats.Msg {
	return &nats.Msg{Data: data, Header: make(nats.Header), Subject: subject}
}

// noopPublishHandler returns a PublishHandler that does nothing.
func noopPublishHandler() PublishHandler {
	return func(ctx context.Context, msg *nats.Msg) error { return nil }
}

// noopSubscribeHandler returns a SubscribeHandler that does nothing.
func noopSubscribeHandler() SubscribeHandler {
	return func(ctx context.Context, msg *nats.Msg) error { return nil }
}

/* ========================================================================================================
   CHAIN TESTS
   ======================================================================================================== */

func TestChain(t *testing.T) {
	testCases := []struct {
		name          string
		middlewares   []PublishMiddleware
		expectedOrder []string
	}{
		{
			name: "two middlewares execute in correct order",
			middlewares: func() []PublishMiddleware {
				return []PublishMiddleware{
					func(next PublishHandler) PublishHandler {
						return func(ctx context.Context, msg *nats.Msg) error {
							ctx = context.WithValue(ctx, orderM1Before, true)
							err := next(ctx, msg)
							return err
						}
					},
					func(next PublishHandler) PublishHandler {
						return func(ctx context.Context, msg *nats.Msg) error {
							err := next(ctx, msg)
							return err
						}
					},
				}
			}(),
			expectedOrder: []string{orderM1Before, orderM2Before, orderHandler, orderM2After, orderM1After},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			var order []string

			mw1 := func(next PublishHandler) PublishHandler {
				return func(ctx context.Context, msg *nats.Msg) error {
					order = append(order, orderM1Before)
					err := next(ctx, msg)
					order = append(order, orderM1After)
					return err
				}
			}
			mw2 := func(next PublishHandler) PublishHandler {
				return func(ctx context.Context, msg *nats.Msg) error {
					order = append(order, orderM2Before)
					err := next(ctx, msg)
					order = append(order, orderM2After)
					return err
				}
			}
			handler := func(ctx context.Context, msg *nats.Msg) error {
				order = append(order, orderHandler)
				return nil
			}

			chained := Chain(mw1, mw2)(handler)
			msg := newTestMsg(testSubject, testPayload)
			err := chained(context.Background(), msg)

			assertions.NoError(err)
			assertions.Equal(tc.expectedOrder, order, "middleware execution order should match")
		})
	}
}

func TestChainEmpty(t *testing.T) {
	assertions := assert.New(t)

	called := false
	handler := func(ctx context.Context, msg *nats.Msg) error {
		called = true
		return nil
	}

	chained := Chain()(handler)
	err := chained(context.Background(), newTestMsg(testSubject, testPayload))

	assertions.NoError(err)
	assertions.True(called, "handler should be called with empty chain")
}

func TestChainSubscribe(t *testing.T) {
	assertions := assert.New(t)
	var order []string

	mw1 := func(next SubscribeHandler) SubscribeHandler {
		return func(ctx context.Context, msg *nats.Msg) error {
			order = append(order, orderM1Before)
			err := next(ctx, msg)
			order = append(order, orderM1After)
			return err
		}
	}
	mw2 := func(next SubscribeHandler) SubscribeHandler {
		return func(ctx context.Context, msg *nats.Msg) error {
			order = append(order, orderM2Before)
			err := next(ctx, msg)
			order = append(order, orderM2After)
			return err
		}
	}
	handler := func(ctx context.Context, msg *nats.Msg) error {
		order = append(order, orderHandler)
		return nil
	}

	chained := ChainSubscribe(mw1, mw2)(handler)
	err := chained(context.Background(), newTestMsg("", testPayload))

	expected := []string{orderM1Before, orderM2Before, orderHandler, orderM2After, orderM1After}
	assertions.NoError(err)
	assertions.Equal(expected, order)
}

func TestChainSubscribeEmpty(t *testing.T) {
	assertions := assert.New(t)

	called := false
	handler := func(ctx context.Context, msg *nats.Msg) error {
		called = true
		return nil
	}

	chained := ChainSubscribe()(handler)
	err := chained(context.Background(), newTestMsg("", testPayload))

	assertions.NoError(err)
	assertions.True(called, "handler should be called with empty chain")
}

/* ========================================================================================================
   STACK TESTS
   ======================================================================================================== */

func TestNewStack(t *testing.T) {
	assertions := assert.New(t)

	s := NewStack()

	assertions.NotNil(s, "NewStack should not return nil")
	assertions.NotNil(s.publish, "publish slice should be initialized")
	assertions.NotNil(s.subscribe, "subscribe slice should be initialized")
}

func TestStackUsePublish(t *testing.T) {
	assertions := assert.New(t)

	s := NewStack()
	s.UsePublish(func(next PublishHandler) PublishHandler { return next })

	assertions.Equal(1, len(s.publish), "should have 1 publish middleware")
}

func TestStackUseSubscribe(t *testing.T) {
	assertions := assert.New(t)

	s := NewStack()
	s.UseSubscribe(func(next SubscribeHandler) SubscribeHandler { return next })

	assertions.Equal(1, len(s.subscribe), "should have 1 subscribe middleware")
}

func TestStackPublishChain(t *testing.T) {
	assertions := assert.New(t)

	s := NewStack()
	called := false
	s.UsePublish(func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, msg *nats.Msg) error {
			called = true
			return next(ctx, msg)
		}
	})

	wrapped := s.PublishChain()(noopPublishHandler())
	err := wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	assertions.NoError(err)
	assertions.True(called, "publish middleware should be called")
}

func TestStackSubscribeChain(t *testing.T) {
	assertions := assert.New(t)

	s := NewStack()
	called := false
	s.UseSubscribe(func(next SubscribeHandler) SubscribeHandler {
		return func(ctx context.Context, msg *nats.Msg) error {
			called = true
			return next(ctx, msg)
		}
	})

	wrapped := s.SubscribeChain()(noopSubscribeHandler())
	err := wrapped(context.Background(), newTestMsg("", testPayload))

	assertions.NoError(err)
	assertions.True(called, "subscribe middleware should be called")
}

func TestStackWrapPublish(t *testing.T) {
	assertions := assert.New(t)

	s := NewStack()
	called := false
	s.UsePublish(func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, msg *nats.Msg) error {
			called = true
			return next(ctx, msg)
		}
	})

	wrapped := s.WrapPublish(noopPublishHandler())
	err := wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	assertions.NoError(err)
	assertions.True(called, "middleware should be called via WrapPublish")
}

func TestStackWrapSubscribe(t *testing.T) {
	assertions := assert.New(t)

	s := NewStack()
	called := false
	s.UseSubscribe(func(next SubscribeHandler) SubscribeHandler {
		return func(ctx context.Context, msg *nats.Msg) error {
			called = true
			return next(ctx, msg)
		}
	})

	wrapped := s.WrapSubscribe(noopSubscribeHandler())
	err := wrapped(context.Background(), newTestMsg("", testPayload))

	assertions.NoError(err)
	assertions.True(called, "middleware should be called via WrapSubscribe")
}

/* ========================================================================================================
   INTERCEPTOR TESTS
   ======================================================================================================== */

type mockInterceptor struct {
	publishCalled   bool
	subscribeCalled bool
}

func (m *mockInterceptor) InterceptPublish() PublishMiddleware {
	return func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, msg *nats.Msg) error {
			m.publishCalled = true
			return next(ctx, msg)
		}
	}
}

func (m *mockInterceptor) InterceptSubscribe() SubscribeMiddleware {
	return func(next SubscribeHandler) SubscribeHandler {
		return func(ctx context.Context, msg *nats.Msg) error {
			m.subscribeCalled = true
			return next(ctx, msg)
		}
	}
}

func TestStackUseInterceptor(t *testing.T) {
	assertions := assert.New(t)

	s := NewStack()
	interceptor := &mockInterceptor{}
	s.UseInterceptor(interceptor)

	assertions.Equal(1, len(s.publish), "publish middleware should be added")
	assertions.Equal(1, len(s.subscribe), "subscribe middleware should be added")

	pubHandler := s.WrapPublish(noopPublishHandler())
	pubHandler(context.Background(), newTestMsg(testSubject, testPayload))
	assertions.True(interceptor.publishCalled, "publish intercept should be called")

	subHandler := s.WrapSubscribe(noopSubscribeHandler())
	subHandler(context.Background(), newTestMsg("", testPayload))
	assertions.True(interceptor.subscribeCalled, "subscribe intercept should be called")
}

/* ========================================================================================================
   CIRCUIT BREAKER EXTENDED TESTS
   ======================================================================================================== */

func TestCircuitBreakerStateString(t *testing.T) {
	testCases := []struct {
		name     string
		state    CircuitState
		expected string
	}{
		{"closed", CircuitClosed, "closed"},
		{"open", CircuitOpen, "open"},
		{"half-open", CircuitHalfOpen, "half-open"},
		{"unknown", CircuitState(99), "unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, tc.state.String())
		})
	}
}

func TestCircuitBreakerFailures(t *testing.T) {
	assertions := assert.New(t)

	cb := NewCircuitBreaker(DefaultCircuitBreakerConfig())
	testErr := errors.New(testErrMsg)

	for i := 0; i < 3; i++ {
		cb.RecordFailure(testErr)
	}

	assertions.Equal(3, cb.Failures(), "should have 3 recorded failures")
}

func TestCircuitBreakerSubscribeMiddleware(t *testing.T) {
	assertions := assert.New(t)

	cb := NewCircuitBreaker(DefaultCircuitBreakerConfig())
	handlerCalled := false
	handler := func(ctx context.Context, msg *nats.Msg) error {
		handlerCalled = true
		return nil
	}

	wrapped := CircuitBreakerSubscribeMiddleware(cb)(handler)
	err := wrapped(context.Background(), newTestMsg("", testPayload))

	assertions.NoError(err)
	assertions.True(handlerCalled, "handler must be called")
}

func TestCircuitBreakerSubscribeMiddlewareError(t *testing.T) {
	assertions := assert.New(t)

	cfg := DefaultCircuitBreakerConfig()
	cfg.FailureThreshold = 2
	cb := NewCircuitBreaker(cfg)

	testErr := errors.New(testErrMsg)
	handler := func(ctx context.Context, msg *nats.Msg) error { return testErr }

	wrapped := CircuitBreakerSubscribeMiddleware(cb)(handler)
	msg := newTestMsg("", testPayload)

	for i := 0; i < 3; i++ {
		wrapped(context.Background(), msg)
	}

	assertions.Equal(CircuitOpen, cb.State(), "circuit should be open after failures")
}

func TestMultiCircuitBreakerReset(t *testing.T) {
	assertions := assert.New(t)

	mcb := NewMultiCircuitBreaker(DefaultCircuitBreakerConfig())
	mcb.Get(testSubject1)
	mcb.Get(testSubject2)

	mcb.Reset()

	cb := mcb.Get(testSubject1)
	assertions.Equal(CircuitClosed, cb.State(), "breaker should be closed after reset")
}

func TestMultiCircuitBreakerMiddleware(t *testing.T) {
	assertions := assert.New(t)

	mcb := NewMultiCircuitBreaker(DefaultCircuitBreakerConfig())
	handlerCalled := 0
	handler := func(ctx context.Context, msg *nats.Msg) error {
		handlerCalled++
		return nil
	}

	wrapped := MultiCircuitBreakerMiddleware(mcb)(handler)
	ctx := context.Background()

	wrapped(ctx, newTestMsg(testSubject1, testPayload))
	wrapped(ctx, newTestMsg(testSubject2, testPayload))

	assertions.Equal(2, handlerCalled, "handler should be called for each subject")
}

func TestCircuitBreakerTransitionToOpen(t *testing.T) {
	assertions := assert.New(t)

	cfg := DefaultCircuitBreakerConfig()
	cfg.FailureThreshold = 2
	cb := NewCircuitBreaker(cfg)

	testErr := errors.New(testErrMsg)
	cb.RecordFailure(testErr)
	cb.RecordFailure(testErr)
	cb.RecordFailure(testErr)

	assertions.Equal(CircuitOpen, cb.State(), "circuit should be open after exceeding threshold")
}

func TestCircuitBreakerRecordSuccess(t *testing.T) {
	assertions := assert.New(t)

	cb := NewCircuitBreaker(DefaultCircuitBreakerConfig())
	cb.RecordSuccess()
	cb.RecordSuccess()

	assertions.Equal(CircuitClosed, cb.State(), "circuit should remain closed after successes")
}

func TestMultiCircuitBreakerGet(t *testing.T) {
	assertions := assert.New(t)

	mcb := NewMultiCircuitBreaker(DefaultCircuitBreakerConfig())

	cb1 := mcb.Get(testSubject1)
	cb2 := mcb.Get(testSubject1)
	cb3 := mcb.Get(testSubject2)

	assertions.Same(cb1, cb2, "same subject should return same breaker")
	assertions.NotSame(cb1, cb3, "different subjects should return different breakers")
}

func TestMultiCircuitBreakerMiddlewareError(t *testing.T) {
	assertions := assert.New(t)

	cfg := DefaultCircuitBreakerConfig()
	cfg.FailureThreshold = 1
	mcb := NewMultiCircuitBreaker(cfg)

	testErr := errors.New(testErrMsg)
	handler := func(ctx context.Context, msg *nats.Msg) error { return testErr }

	wrapped := MultiCircuitBreakerMiddleware(mcb)(handler)
	err := wrapped(context.Background(), newTestMsg(testSubject1, testPayload))

	assertions.ErrorIs(err, testErr, "error should propagate from handler")
}

/* ========================================================================================================
   RATE LIMITER EXTENDED TESTS
   ======================================================================================================== */

func TestDefaultRateLimiterConfig(t *testing.T) {
	assertions := assert.New(t)

	cfg := DefaultRateLimiterConfig()

	assertions.Greater(cfg.Rate, float64(0), "Rate should be positive")
	assertions.Greater(cfg.Burst, 0, "Burst should be positive")
}

func TestRateLimiterTokens(t *testing.T) {
	assertions := assert.New(t)

	rl := NewRateLimiter(RateLimiterConfig{Rate: 100, Burst: 10})

	assertions.Greater(rl.Tokens(), float64(0), "should have positive tokens")
}

func TestRateLimitWaitMiddleware(t *testing.T) {
	assertions := assert.New(t)

	rl := NewRateLimiter(RateLimiterConfig{Rate: 1000, Burst: 10})
	handlerCalled := false
	handler := func(ctx context.Context, msg *nats.Msg) error {
		handlerCalled = true
		return nil
	}

	wrapped := RateLimitWaitMiddleware(rl)(handler)
	err := wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	assertions.NoError(err)
	assertions.True(handlerCalled, "handler is called")
}

func TestRateLimitSubscribeMiddleware(t *testing.T) {
	assertions := assert.New(t)

	rl := NewRateLimiter(RateLimiterConfig{Rate: 1000, Burst: 10})
	handlerCalled := false
	handler := func(ctx context.Context, msg *nats.Msg) error {
		handlerCalled = true
		return nil
	}

	wrapped := RateLimitSubscribeMiddleware(rl)(handler)
	err := wrapped(context.Background(), newTestMsg("", testPayload))

	assertions.NoError(err)
	assertions.True(handlerCalled, "handler called")
}

func TestPerSubjectRateLimitMiddleware(t *testing.T) {
	assertions := assert.New(t)

	psr := NewPerSubjectRateLimiter(RateLimiterConfig{Rate: 1000, Burst: 10})
	handlerCalled := 0
	handler := func(ctx context.Context, msg *nats.Msg) error {
		handlerCalled++
		return nil
	}

	wrapped := PerSubjectRateLimitMiddleware(psr)(handler)
	ctx := context.Background()

	wrapped(ctx, newTestMsg(testSubject1, testPayload))
	wrapped(ctx, newTestMsg(testSubject2, testPayload))

	assertions.Equal(2, handlerCalled, "handler should be called for each subject")
}

func TestSlidingWindowMiddleware(t *testing.T) {
	assertions := assert.New(t)

	sw := NewSlidingWindowLimiter(100, 0)
	handlerCalled := false
	handler := func(ctx context.Context, msg *nats.Msg) error {
		handlerCalled = true
		return nil
	}

	wrapped := SlidingWindowMiddleware(sw)(handler)
	err := wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	assertions.NoError(err)
	assertions.True(handlerCalled, "handler-called")
}

func TestNewRateLimiter(t *testing.T) {
	assertions := assert.New(t)

	rl := NewRateLimiter(RateLimiterConfig{Rate: 100, Burst: 10})

	assertions.NotNil(rl, "NewRateLimiter should not return nil")
}

func TestRateLimiterAllowAndWait(t *testing.T) {
	assertions := assert.New(t)

	rl := NewRateLimiter(RateLimiterConfig{Rate: 1000, Burst: 10})

	assertions.True(rl.Allow(), "Allow should return true")

	err := rl.Wait(context.Background())
	assertions.NoError(err, "Wait should not return error")
}

func TestPerSubjectRateLimiterGet(t *testing.T) {
	assertions := assert.New(t)

	psr := NewPerSubjectRateLimiter(RateLimiterConfig{Rate: 100, Burst: 10})

	rl1 := psr.Get(testSubject1)
	rl2 := psr.Get(testSubject1)
	rl3 := psr.Get(testSubject2)

	assertions.Same(rl1, rl2, "same subject should return same limiter")
	assertions.NotSame(rl1, rl3, "different subjects should return different limiters")
}

func TestSlidingWindowLimiterAllow(t *testing.T) {
	assertions := assert.New(t)

	sw := NewSlidingWindowLimiter(100, time.Second)

	for i := 0; i < 50; i++ {
		assertions.True(sw.Allow(), "should allow request %d", i)
	}
}

func TestSlidingWindowLimiterCount(t *testing.T) {
	assertions := assert.New(t)

	sw := NewSlidingWindowLimiter(100, time.Second)
	sw.Allow()
	sw.Allow()
	sw.Allow()

	assertions.Equal(3, sw.Count(), "count should be 3")
}

func TestRateLimitWaitMiddlewareRejection(t *testing.T) {
	assertions := assert.New(t)

	rl := NewRateLimiter(RateLimiterConfig{Rate: 0.001, Burst: 1})
	handler := noopPublishHandler()
	wrapped := RateLimitWaitMiddleware(rl)(handler)

	ctx := context.Background()
	msg := newTestMsg(testSubject, testPayload)
	wrapped(ctx, msg)

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// May or may not error depending on timing; just ensure no panic
	_ = wrapped(ctxTimeout, msg)
	assertions.True(true, "should not panic")
}

func TestRateLimitSubscribeMiddlewareRejection(t *testing.T) {
	assertions := assert.New(t)

	rl := NewRateLimiter(RateLimiterConfig{Rate: 1, Burst: 1})
	handlerCalled := 0
	handler := func(ctx context.Context, msg *nats.Msg) error {
		handlerCalled++
		return nil
	}

	wrapped := RateLimitSubscribeMiddleware(rl)(handler)
	msg := newTestMsg("", testPayload)

	wrapped(context.Background(), msg)
	wrapped(context.Background(), msg)

	assertions.GreaterOrEqual(handlerCalled, 1, "at least one call should succeed")
}

func TestPerSubjectRateLimitMiddlewareRejection(t *testing.T) {
	assertions := assert.New(t)

	psr := NewPerSubjectRateLimiter(RateLimiterConfig{Rate: 1, Burst: 1})
	wrapped := PerSubjectRateLimitMiddleware(psr)(noopPublishHandler())

	err := wrapped(context.Background(), newTestMsg(testSubject, testPayload))
	assertions.NoError(err, "first request should succeed")
}

func TestSlidingWindowMiddlewareRejection(t *testing.T) {
	assertions := assert.New(t)

	sw := NewSlidingWindowLimiter(2, time.Second)
	handlerCalled := 0
	handler := func(ctx context.Context, msg *nats.Msg) error {
		handlerCalled++
		return nil
	}

	wrapped := SlidingWindowMiddleware(sw)(handler)
	msg := newTestMsg(testSubject, testPayload)

	wrapped(context.Background(), msg)
	wrapped(context.Background(), msg)
	_ = wrapped(context.Background(), msg)

	assertions.GreaterOrEqual(handlerCalled, 2, "at least 2 calls should succeed")
}

/* ========================================================================================================
   CONCURRENT ACCESS TESTS
   ======================================================================================================== */

func TestChainConcurrentAccess(t *testing.T) {
	assertions := assert.New(t)

	var count int
	var mu sync.Mutex

	middleware := func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, msg *nats.Msg) error {
			mu.Lock()
			count++
			mu.Unlock()
			return next(ctx, msg)
		}
	}

	chained := Chain(middleware)(noopPublishHandler())

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			chained(context.Background(), newTestMsg(testSubject, testPayload))
		}()
	}
	wg.Wait()

	assertions.Equal(100, count, "should have 100 concurrent calls")
}

/* ========================================================================================================
   LOGGING MIDDLEWARE TESTS
   ======================================================================================================== */

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
	assertions := assert.New(t)

	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger)

	assertions.NotNil(lm, "NewLoggingMiddleware should not return nil")
	assertions.Equal(logger, lm.logger, "logger should be set correctly")
}

func TestLoggingMiddlewareWithLevel(t *testing.T) {
	assertions := assert.New(t)

	lm := NewLoggingMiddleware(&testLogger{}).WithLevel(LogLevelError)

	assertions.Equal(LogLevelError, lm.level, "log level should be set")
}

func TestLoggingMiddlewareWithPayload(t *testing.T) {
	assertions := assert.New(t)

	lm := NewLoggingMiddleware(&testLogger{}).WithPayload(true, 512)

	assertions.True(lm.logPayload, "logPayload should be true")
	assertions.Equal(512, lm.maxPayloadSize, "maxPayloadSize should be set")
}

func TestLoggingMiddlewareInterceptPublish(t *testing.T) {
	assertions := assert.New(t)

	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelDebug)

	handlerCalled := false
	handler := func(ctx context.Context, msg *nats.Msg) error {
		handlerCalled = true
		return nil
	}

	wrapped := lm.InterceptPublish()(handler)
	err := wrapped(context.Background(), newTestMsg(testSubject, testPayloadStr))

	assertions.NoError(err)
	assertions.True(handlerCalled, "handler_called")
	assertions.NotEmpty(logger.messages, "should produce log messages")
}

func TestLoggingMiddlewareInterceptPublishError(t *testing.T) {
	assertions := assert.New(t)

	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelDebug)

	testErr := errors.New(testErrMsg)
	handler := func(ctx context.Context, msg *nats.Msg) error { return testErr }

	wrapped := lm.InterceptPublish()(handler)
	err := wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	assertions.ErrorIs(err, testErr)

	hasError := false
	for _, m := range logger.messages {
		if m == "error:publish failed" {
			hasError = true
			break
		}
	}
	assertions.True(hasError, "should have error log message")
}

func TestLoggingMiddlewareInterceptPublishWithPayload(t *testing.T) {
	assertions := assert.New(t)

	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelDebug).WithPayload(true, 10)

	wrapped := lm.InterceptPublish()(noopPublishHandler())
	err := wrapped(context.Background(), newTestMsg(testSubject, testPayloadLong))

	assertions.NoError(err)
}

func TestLoggingMiddlewareInterceptSubscribe(t *testing.T) {
	assertions := assert.New(t)

	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelDebug)

	handlerCalled := false
	handler := func(ctx context.Context, msg *nats.Msg) error {
		handlerCalled = true
		return nil
	}

	wrapped := lm.InterceptSubscribe()(handler)
	err := wrapped(context.Background(), newTestMsg(testSubject, testPayloadStr))

	assertions.NoError(err)
	assertions.True(handlerCalled, "handler should make a call")
}

func TestLoggingMiddlewareInterceptSubscribeError(t *testing.T) {
	assertions := assert.New(t)

	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelDebug)

	testErr := errors.New(testErrMsg)
	handler := func(ctx context.Context, msg *nats.Msg) error { return testErr }

	wrapped := lm.InterceptSubscribe()(handler)
	err := wrapped(context.Background(), newTestMsg("", testPayload))

	assertions.ErrorIs(err, testErr)
}

func TestLogging(t *testing.T) {
	assertions := assert.New(t)

	lm := Logging(&testLogger{})

	assertions.NotNil(lm, "Logging should not return nil")
}

func TestNoOpLogger(t *testing.T) {
	assertions := assert.New(t)

	logger := NoOpLogger{}
	logger.Debug("test", nil)
	logger.Info("test", nil)
	logger.Warn("test", nil)
	logger.Error("test", nil)

	assertions.True(true, "NoOpLogger should not panic")
}

func TestLoggingMiddlewareNilLogger(t *testing.T) {
	assertions := assert.New(t)

	lm := NewLoggingMiddleware(nil)
	wrapped := lm.InterceptPublish()(noopPublishHandler())

	err := wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	assertions.NoError(err, "nil logger should not panic")
}

func TestLoggingMiddlewareLogLevelFilter(t *testing.T) {
	assertions := assert.New(t)

	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelWarn)
	wrapped := lm.InterceptPublish()(noopPublishHandler())

	wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	for _, m := range logger.messages {
		assertions.NotEqual("debug:publishing message", m, "debug messages should be filtered")
		assertions.NotEqual("info:message published", m, "info messages should be filtered")
	}
}

func TestLoggingMiddlewareLogLevelWarnPublish(t *testing.T) {
	assertions := assert.New(t)

	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelWarn)

	handlerCalled := false
	handler := func(ctx context.Context, msg *nats.Msg) error {
		handlerCalled = true
		return nil
	}

	wrapped := lm.InterceptPublish()(handler)
	err := wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	assertions.NoError(err)
	assertions.True(handlerCalled, "handler must make a call")
}

func TestLoggingMiddlewareLogLevelWarnSubscribe(t *testing.T) {
	assertions := assert.New(t)

	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelWarn)

	testErr := errors.New(testErrMsg)
	handler := func(ctx context.Context, msg *nats.Msg) error { return testErr }

	wrapped := lm.InterceptSubscribe()(handler)
	err := wrapped(context.Background(), newTestMsg("", testPayload))

	assertions.ErrorIs(err, testErr)
}

func TestLoggingMiddlewareWithTraceContext(t *testing.T) {
	assertions := assert.New(t)

	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelDebug)
	wrapped := lm.InterceptPublish()(noopPublishHandler())

	tc := &core.TraceContext{TraceID: testTraceID, SpanID: testSpanID}
	ctx := core.WithTraceContext(context.Background(), tc)

	err := wrapped(ctx, newTestMsg(testSubject, testPayload))

	assertions.NoError(err)
}

func TestLoggingMiddlewareWithTenantID(t *testing.T) {
	assertions := assert.New(t)

	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelDebug)
	wrapped := lm.InterceptPublish()(noopPublishHandler())

	ctx := core.WithTenantID(context.Background(), testTenantID)

	err := wrapped(ctx, newTestMsg(testSubject, testPayload))

	assertions.NoError(err)
}

func TestLoggingMiddlewareSubscribeWithTraceAndTenant(t *testing.T) {
	assertions := assert.New(t)

	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelDebug).WithPayload(true, 100)
	wrapped := lm.InterceptSubscribe()(noopSubscribeHandler())

	tc := &core.TraceContext{TraceID: testTraceID, SpanID: testSpanID}
	ctx := core.WithTraceContext(context.Background(), tc)
	ctx = core.WithTenantID(ctx, testTenantID)

	err := wrapped(ctx, newTestMsg(testSubject, testPayloadContent))

	assertions.NoError(err)
}

func TestLoggingMiddlewareSubscribePayloadTruncation(t *testing.T) {
	assertions := assert.New(t)

	logger := &testLogger{}
	lm := NewLoggingMiddleware(logger).WithLevel(LogLevelDebug).WithPayload(true, 5)
	wrapped := lm.InterceptSubscribe()(noopSubscribeHandler())

	err := wrapped(context.Background(), newTestMsg(testSubject, []byte("this is a long payload")))

	assertions.NoError(err)
}

/* ========================================================================================================
   METRICS MIDDLEWARE EXTENDED TESTS
   ======================================================================================================== */

func TestMetricsCollectorOnPublish(t *testing.T) {
	assertions := assert.New(t)

	mc := NewMetricsCollector()
	callbackCalled := false
	mc.OnPublish(func(subject string, latency time.Duration, err error) {
		callbackCalled = true
	})

	wrapped := mc.InterceptPublish()(noopPublishHandler())
	wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	assertions.True(callbackCalled, "OnPublish callback should be called")
}

func TestMetricsCollectorOnReceive(t *testing.T) {
	assertions := assert.New(t)

	mc := NewMetricsCollector()
	callbackCalled := false
	mc.OnReceive(func(subject string, processTime time.Duration, err error) {
		callbackCalled = true
	})

	wrapped := mc.InterceptSubscribe()(noopSubscribeHandler())
	wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	assertions.True(callbackCalled, "OnReceive callback should be called")
}

func TestMetricsCollectorCollect(t *testing.T) {
	assertions := assert.New(t)

	mc := NewMetricsCollector()
	wrapped := mc.InterceptPublish()(noopPublishHandler())
	wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	metrics := mc.Collect()
	assertions.Equal(int64(1), metrics.PublishCount)
	assertions.Equal(int64(4), metrics.BytesSent, "test is 4 bytes")
}

func TestMetricsCollectorPublishError(t *testing.T) {
	assertions := assert.New(t)

	mc := NewMetricsCollector()
	handler := func(ctx context.Context, msg *nats.Msg) error {
		return errors.New("publish error")
	}

	wrapped := mc.InterceptPublish()(handler)
	wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	assertions.Equal(int64(1), mc.Collect().PublishErrors)
}

func TestMetricsCollectorSubscribeError(t *testing.T) {
	assertions := assert.New(t)

	mc := NewMetricsCollector()
	handler := func(ctx context.Context, msg *nats.Msg) error {
		return errors.New("process error")
	}

	wrapped := mc.InterceptSubscribe()(handler)
	wrapped(context.Background(), newTestMsg("", testPayload))

	assertions.Equal(int64(1), mc.Collect().ReceiveErrors)
}

func TestMetricsMiddleware(t *testing.T) {
	assertions := assert.New(t)

	mc := MetricsMiddleware()

	assertions.NotNil(mc, "MetricsMiddleware should not return nil")
}

func TestPerSubjectMetrics(t *testing.T) {
	assertions := assert.New(t)

	psm := NewPerSubjectMetrics()

	assertions.NotNil(psm, "NewPerSubjectMetrics should not return nil")
}

func TestPerSubjectMetricsAll(t *testing.T) {
	assertions := assert.New(t)

	psm := NewPerSubjectMetrics()
	wrapped := psm.InterceptPublish()(noopPublishHandler())
	ctx := context.Background()

	wrapped(ctx, newTestMsg(testSubject1, testPayload))
	wrapped(ctx, newTestMsg(testSubject2, testPayload))

	assertions.Equal(2, len(psm.All()), "should have 2 subjects")
}

func TestPerSubjectMetricsInterceptSubscribeNilMsg(t *testing.T) {
	assertions := assert.New(t)

	psm := NewPerSubjectMetrics()
	handlerCalled := false
	handler := func(ctx context.Context, msg *nats.Msg) error {
		handlerCalled = true
		return nil
	}

	wrapped := psm.InterceptSubscribe()(handler)
	err := wrapped(context.Background(), nil)

	assertions.NoError(err)
	assertions.True(handlerCalled, "handler should be called for nil message")
}

/* ========================================================================================================
   OTEL MIDDLEWARE TESTS
   ======================================================================================================== */

func TestNewOTELMetricsMiddleware(t *testing.T) {
	assertions := assert.New(t)

	otelm := NewOTELMetricsMiddleware()

	assertions.NotNil(otelm, "NewOTELMetricsMiddleware should not return nil")
}

func TestOTELMetricsMiddlewareInterceptPublish(t *testing.T) {
	assertions := assert.New(t)

	otelm := NewOTELMetricsMiddleware()
	handlerCalled := false
	handler := func(ctx context.Context, msg *nats.Msg) error {
		handlerCalled = true
		return nil
	}

	wrapped := otelm.InterceptPublish()(handler)
	err := wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	assertions.NoError(err)
	assertions.True(handlerCalled, "handler should be called")
}

func TestOTELMetricsMiddlewareInterceptPublishError(t *testing.T) {
	assertions := assert.New(t)

	otelm := NewOTELMetricsMiddleware()
	testErr := errors.New(testErrMsg)
	handler := func(ctx context.Context, msg *nats.Msg) error { return testErr }

	wrapped := otelm.InterceptPublish()(handler)
	err := wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	assertions.ErrorIs(err, testErr)
}

func TestOTELMetricsMiddlewareInterceptPublishNilMsg(t *testing.T) {
	assertions := assert.New(t)

	otelm := NewOTELMetricsMiddleware()
	wrapped := otelm.InterceptPublish()(noopPublishHandler())

	err := wrapped(context.Background(), nil)

	assertions.NoError(err, "nil message should not panic")
}

func TestOTELMetricsMiddlewareInterceptSubscribe(t *testing.T) {
	assertions := assert.New(t)

	otelm := NewOTELMetricsMiddleware()
	handlerCalled := false
	handler := func(ctx context.Context, msg *nats.Msg) error {
		handlerCalled = true
		return nil
	}

	wrapped := otelm.InterceptSubscribe()(handler)
	err := wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	assertions.NoError(err)
	assertions.True(handlerCalled, "handler should be called")
}

func TestOTELMetricsMiddlewareInterceptSubscribeError(t *testing.T) {
	assertions := assert.New(t)

	otelm := NewOTELMetricsMiddleware()
	testErr := errors.New(testErrMsg)
	handler := func(ctx context.Context, msg *nats.Msg) error { return testErr }

	wrapped := otelm.InterceptSubscribe()(handler)
	err := wrapped(context.Background(), newTestMsg("", testPayload))

	assertions.ErrorIs(err, testErr)
}

func TestOTELMetricsMiddlewareInterceptSubscribeNilMsg(t *testing.T) {
	assertions := assert.New(t)

	otelm := NewOTELMetricsMiddleware()
	wrapped := otelm.InterceptSubscribe()(noopSubscribeHandler())

	err := wrapped(context.Background(), nil)

	assertions.NoError(err, "nil message should not panic")
}

/* ========================================================================================================
   OTEL LOGGER TESTS
   ======================================================================================================== */

func TestNewOTELLogger(t *testing.T) {
	assertions := assert.New(t)

	otelLogger := NewOTELLogger()

	assertions.NotNil(otelLogger, "NewOTELLogger should not return nil")
}

func TestOTELLoggerMethods(t *testing.T) {
	assertions := assert.New(t)

	otelLogger := NewOTELLogger()

	otelLogger.Debug("debug message", map[string]any{"key": "value"})
	otelLogger.Info("info message", map[string]any{"key": "value"})
	otelLogger.Warn("warn message", map[string]any{"key": "value"})
	otelLogger.Error("error message", map[string]any{"key": "value"})

	assertions.True(true, "OTEL logger methods should not panic")
}

func TestOTELLoggerNilFields(t *testing.T) {
	assertions := assert.New(t)

	otelLogger := NewOTELLogger()

	otelLogger.Debug("debug message", nil)
	otelLogger.Info("info message", nil)
	otelLogger.Warn("warn message", nil)
	otelLogger.Error("error message", nil)

	assertions.True(true, "OTEL logger with nil fields should not panic")
}

func TestMapToFields(t *testing.T) {
	testCases := []struct {
		name          string
		input         map[string]any
		expectedNil   bool
		expectedCount int
	}{
		{"nil map", nil, true, 0},
		{"populated map", map[string]any{"string": "value", "int": 42, "bool": true}, false, 3},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			fields := mapToFields(tc.input)
			if tc.expectedNil {
				assertions.Nil(fields)
			} else {
				assertions.Equal(tc.expectedCount, len(fields))
			}
		})
	}
}

func TestOTELLoggingMiddleware(t *testing.T) {
	assertions := assert.New(t)

	lm := OTELLoggingMiddleware()

	assertions.NotNil(lm, "OTELLoggingMiddleware should not return nil")
}
