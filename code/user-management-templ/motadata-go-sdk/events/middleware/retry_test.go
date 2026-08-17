package middleware

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

/* ========================================================================================================
   RETRY CONFIG TESTS
   ======================================================================================================== */

var (
	retryMinimalInterval = 1 * time.Millisecond
	retryMinimalMax      = 10 * time.Millisecond
	retryDefaultMult     = 2.0
	retryZeroJitter      = 0.0
)

func TestDefaultRetryConfig(t *testing.T) {
	assertions := assert.New(t)

	config := DefaultRetryConfig()

	assertions.Equal(3, config.MaxAttempts, "Default MaxAttempts should be 3")
	assertions.Equal(100*time.Millisecond, config.InitialInterval, "Default InitialInterval should be 100ms")
	assertions.Equal(5*time.Second, config.MaxInterval, "Default MaxInterval should be 5s")
	assertions.Equal(retryDefaultMult, config.Multiplier, "Default Multiplier should be 2.0")
	assertions.Equal(0.1, config.Jitter, "Default Jitter should be 0.1")
}

/* ========================================================================================================
   RETRY MIDDLEWARE TESTS
   ======================================================================================================== */

func TestNewRetryMiddleware(t *testing.T) {
	assertions := assert.New(t)

	mw := NewRetryMiddleware(DefaultRetryConfig())

	assertions.NotNil(mw, "RetryMiddleware should not be nil")
	assertions.NotNil(mw.shouldRetry, "shouldRetry should not be nil")
	assertions.NotNil(mw.rng, "rng should not be nil")
}

func TestRetryMiddlewareSuccess(t *testing.T) {
	assertions := assert.New(t)

	mw := NewRetryMiddleware(RetryConfig{
		MaxAttempts:     3,
		InitialInterval: retryMinimalInterval,
		MaxInterval:     retryMinimalMax,
		Multiplier:      retryDefaultMult,
		Jitter:          retryZeroJitter,
	})

	var attempts int32
	handler := func(ctx context.Context, msg *nats.Msg) error {
		atomic.AddInt32(&attempts, 1)
		return nil
	}

	wrapped := mw.InterceptPublish()(handler)
	err := wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	assertions.NoError(err, "Should succeed on first try")
	assertions.Equal(int32(1), atomic.LoadInt32(&attempts), "Should only attempt once")
}

func TestRetryMiddlewareRetryOnError(t *testing.T) {
	assertions := assert.New(t)

	mw := NewRetryMiddleware(RetryConfig{
		MaxAttempts:     3,
		InitialInterval: retryMinimalInterval,
		MaxInterval:     retryMinimalMax,
		Multiplier:      retryDefaultMult,
		Jitter:          retryZeroJitter,
	})

	var attempts int32
	handler := func(ctx context.Context, msg *nats.Msg) error {
		count := atomic.AddInt32(&attempts, 1)
		if count < 3 {
			return utils.ErrNotConnected
		}
		return nil
	}

	wrapped := mw.InterceptPublish()(handler)
	err := wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	assertions.NoError(err, "Should succeed after retries")
	assertions.Equal(int32(3), atomic.LoadInt32(&attempts), "Should retry until success")
}

func TestRetryMiddlewareExhausted(t *testing.T) {
	assertions := assert.New(t)

	mw := NewRetryMiddleware(RetryConfig{
		MaxAttempts:     2,
		InitialInterval: retryMinimalInterval,
		MaxInterval:     retryMinimalMax,
		Multiplier:      retryDefaultMult,
		Jitter:          retryZeroJitter,
	})

	var attempts int32
	expectedErr := utils.ErrNotConnected
	handler := func(ctx context.Context, msg *nats.Msg) error {
		atomic.AddInt32(&attempts, 1)
		return expectedErr
	}

	wrapped := mw.InterceptPublish()(handler)
	err := wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	assertions.ErrorIs(err, expectedErr, "Should return last error")
	assertions.Equal(int32(3), atomic.LoadInt32(&attempts), "Should attempt MaxAttempts + 1 times")
}

func TestRetryMiddlewareNonRetryableError(t *testing.T) {
	assertions := assert.New(t)

	mw := NewRetryMiddleware(RetryConfig{
		MaxAttempts:     3,
		InitialInterval: retryMinimalInterval,
		MaxInterval:     retryMinimalMax,
		Multiplier:      retryDefaultMult,
		Jitter:          retryZeroJitter,
	})

	var attempts int32
	handler := func(ctx context.Context, msg *nats.Msg) error {
		atomic.AddInt32(&attempts, 1)
		return utils.ErrInvalidSubject
	}

	wrapped := mw.InterceptPublish()(handler)
	err := wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	assertions.ErrorIs(err, utils.ErrInvalidSubject, "Should return non-retryable error immediately")
	assertions.Equal(int32(1), atomic.LoadInt32(&attempts), "Should not retry non-retryable errors")
}

func TestRetryMiddlewareContextCancellation(t *testing.T) {
	assertions := assert.New(t)

	mw := NewRetryMiddleware(RetryConfig{
		MaxAttempts:     10,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     1 * time.Second,
		Multiplier:      retryDefaultMult,
		Jitter:          retryZeroJitter,
	})

	var attempts int32
	handler := func(ctx context.Context, msg *nats.Msg) error {
		atomic.AddInt32(&attempts, 1)
		return utils.ErrNotConnected
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	wrapped := mw.InterceptPublish()(handler)
	err := wrapped(ctx, newTestMsg(testSubject, testPayload))

	assertions.Error(err, "Should return error on context cancellation")
	assertions.Less(atomic.LoadInt32(&attempts), int32(10), "Should not complete all retries")
}

/* ========================================================================================================
   CUSTOM SHOULD RETRY TESTS
   ======================================================================================================== */

func TestRetryMiddlewareWithShouldRetry(t *testing.T) {
	assertions := assert.New(t)

	customError := errors.New("custom error")

	mw := NewRetryMiddleware(RetryConfig{
		MaxAttempts:     3,
		InitialInterval: retryMinimalInterval,
		MaxInterval:     retryMinimalMax,
		Multiplier:      retryDefaultMult,
		Jitter:          retryZeroJitter,
	}).WithShouldRetry(func(err error) bool {
		return errors.Is(err, customError)
	})

	var attempts int32
	handler := func(ctx context.Context, msg *nats.Msg) error {
		count := atomic.AddInt32(&attempts, 1)
		if count < 3 {
			return customError
		}
		return nil
	}

	wrapped := mw.InterceptPublish()(handler)
	err := wrapped(context.Background(), newTestMsg(testSubject, testPayload))

	assertions.NoError(err, "Should succeed after retries")
	assertions.Equal(int32(3), atomic.LoadInt32(&attempts), "Should retry custom error")
}

/* ========================================================================================================
   BACKOFF CALCULATION TESTS
   ======================================================================================================== */

func TestCalculateBackoff(t *testing.T) {
	mw := NewRetryMiddleware(RetryConfig{
		MaxAttempts:     3,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     1 * time.Second,
		Multiplier:      retryDefaultMult,
		Jitter:          retryZeroJitter,
	})

	testCases := []struct {
		name     string
		attempt  int
		expected time.Duration
	}{
		{"attempt 0 uses initial interval", 0, 100 * time.Millisecond},
		{"attempt 1 doubles", 1, 200 * time.Millisecond},
		{"attempt 2 doubles again", 2, 400 * time.Millisecond},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, mw.calculateBackoff(tc.attempt))
		})
	}
}

func TestCalculateBackoffCappedAtMax(t *testing.T) {
	assertions := assert.New(t)

	mw := NewRetryMiddleware(RetryConfig{
		MaxAttempts:     10,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     500 * time.Millisecond,
		Multiplier:      retryDefaultMult,
		Jitter:          retryZeroJitter,
	})

	backoff := mw.calculateBackoff(10)

	assertions.Equal(500*time.Millisecond, backoff, "Backoff should be capped at max interval")
}

func TestCalculateBackoffWithJitter(t *testing.T) {
	assertions := assert.New(t)

	mw := NewRetryMiddleware(RetryConfig{
		MaxAttempts:     3,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     1 * time.Second,
		Multiplier:      retryDefaultMult,
		Jitter:          0.5,
	})

	var samples []time.Duration
	for i := 0; i < 10; i++ {
		samples = append(samples, mw.calculateBackoff(0))
	}

	for _, sample := range samples {
		assertions.GreaterOrEqual(sample, 100*time.Millisecond, "Backoff should not go below initial")
		assertions.LessOrEqual(sample, 150*time.Millisecond, "Backoff should be within jitter range")
	}

	allSame := true
	for i := 1; i < len(samples); i++ {
		if samples[i] != samples[0] {
			allSame = false
			break
		}
	}
	assertions.False(allSame, "Jitter should produce variation in backoff times")
}

func TestCalculateBackoffNegativeAttempt(t *testing.T) {
	assertions := assert.New(t)

	mw := NewRetryMiddleware(RetryConfig{
		MaxAttempts:     3,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     1 * time.Second,
		Multiplier:      retryDefaultMult,
		Jitter:          retryZeroJitter,
	})

	backoff := mw.calculateBackoff(-1)

	assertions.Equal(100*time.Millisecond, backoff, "Negative attempt should use initial interval")
}

/* ========================================================================================================
   SUBSCRIBE MIDDLEWARE TESTS
   ======================================================================================================== */

func TestRetryMiddlewareInterceptSubscribe(t *testing.T) {
	assertions := assert.New(t)

	mw := NewRetryMiddleware(DefaultRetryConfig())

	var called bool
	handler := func(ctx context.Context, msg *nats.Msg) error {
		called = true
		return nil
	}

	wrapped := mw.InterceptSubscribe()(handler)
	err := wrapped(context.Background(), newTestMsg("", testPayload))

	assertions.NoError(err, "Subscribe should pass through")
	assertions.True(called, "Handler should be called")
}

/* ========================================================================================================
   DEFAULT SHOULD RETRY TESTS
   ======================================================================================================== */

func TestDefaultShouldRetry(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"retryable error", utils.ErrNotConnected, true},
		{"retryable timeout", utils.ErrConnectionTimeout, true},
		{"non-retryable", utils.ErrInvalidSubject, false},
		{"non-retryable config", utils.ErrInvalidConfig, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			result := DefaultShouldRetry(tc.err)
			assertions.Equal(tc.expected, result, "DefaultShouldRetry should return correct value")
		})
	}
}

/* ========================================================================================================
   HELPER FUNCTION TESTS
   ======================================================================================================== */

func TestRetryFactory(t *testing.T) {
	assertions := assert.New(t)

	mw := Retry()

	assertions.NotNil(mw, "Retry() should return RetryMiddleware")
}

func TestRetryWithConfigFactory(t *testing.T) {
	assertions := assert.New(t)

	config := RetryConfig{
		MaxAttempts:     5,
		InitialInterval: 50 * time.Millisecond,
		MaxInterval:     2 * time.Second,
		Multiplier:      3.0,
		Jitter:          0.2,
	}

	mw := RetryWithConfig(config)

	assertions.NotNil(mw, "RetryWithConfig() should return RetryMiddleware")
	assertions.Equal(5, mw.config.MaxAttempts, "Config should be applied")
}
