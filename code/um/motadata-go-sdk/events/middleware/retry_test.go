package middleware

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
	"github.com/stretchr/testify/assert"
)

/* ========================================================================================================
   RETRY CONFIG TESTS
   ======================================================================================================== */

func TestDefaultRetryConfig(t *testing.T) {

	assertions := assert.New(t)

	config := DefaultRetryConfig()

	assertions.Equal(3, config.MaxAttempts, "Default MaxAttempts should be 3")
	assertions.Equal(100*time.Millisecond, config.InitialInterval, "Default InitialInterval should be 100ms")
	assertions.Equal(5*time.Second, config.MaxInterval, "Default MaxInterval should be 5s")
	assertions.Equal(2.0, config.Multiplier, "Default Multiplier should be 2.0")
	assertions.Equal(0.1, config.Jitter, "Default Jitter should be 0.1")
}

/* ========================================================================================================
   RETRY MIDDLEWARE TESTS
   ======================================================================================================== */

func TestNewRetryMiddleware(t *testing.T) {

	assertions := assert.New(t)

	config := DefaultRetryConfig()
	mw := NewRetryMiddleware(config)

	assertions.NotNil(mw, "RetryMiddleware should not be nil")
	assertions.NotNil(mw.shouldRetry, "shouldRetry should not be nil")
	assertions.NotNil(mw.rng, "rng should not be nil")
}

func TestRetryMiddlewareSuccess(t *testing.T) {

	assertions := assert.New(t)

	config := RetryConfig{
		MaxAttempts:     3,
		InitialInterval: 1 * time.Millisecond,
		MaxInterval:     10 * time.Millisecond,
		Multiplier:      2.0,
		Jitter:          0,
	}

	mw := NewRetryMiddleware(config)

	var attempts int32

	handler := func(ctx context.Context, subject string, msg *core.Message) error {

		atomic.AddInt32(&attempts, 1)

		return nil
	}

	interceptor := mw.InterceptPublish()
	wrappedHandler := interceptor(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	err := wrappedHandler(ctx, "test.subject", msg)

	assertions.NoError(err, "Should succeed on first try")
	assertions.Equal(int32(1), atomic.LoadInt32(&attempts), "Should only attempt once")
}

func TestRetryMiddlewareRetryOnError(t *testing.T) {

	assertions := assert.New(t)

	config := RetryConfig{
		MaxAttempts:     3,
		InitialInterval: 1 * time.Millisecond,
		MaxInterval:     10 * time.Millisecond,
		Multiplier:      2.0,
		Jitter:          0,
	}

	mw := NewRetryMiddleware(config)

	var attempts int32

	handler := func(ctx context.Context, subject string, msg *core.Message) error {

		count := atomic.AddInt32(&attempts, 1)

		if count < 3 {

			return core.ErrNotConnected // Retryable error
		}

		return nil
	}

	interceptor := mw.InterceptPublish()
	wrappedHandler := interceptor(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	err := wrappedHandler(ctx, "test.subject", msg)

	assertions.NoError(err, "Should succeed after retries")
	assertions.Equal(int32(3), atomic.LoadInt32(&attempts), "Should retry until success")
}

func TestRetryMiddlewareExhausted(t *testing.T) {

	assertions := assert.New(t)

	config := RetryConfig{
		MaxAttempts:     2,
		InitialInterval: 1 * time.Millisecond,
		MaxInterval:     10 * time.Millisecond,
		Multiplier:      2.0,
		Jitter:          0,
	}

	mw := NewRetryMiddleware(config)

	var attempts int32
	expectedErr := core.ErrNotConnected

	handler := func(ctx context.Context, subject string, msg *core.Message) error {

		atomic.AddInt32(&attempts, 1)

		return expectedErr
	}

	interceptor := mw.InterceptPublish()
	wrappedHandler := interceptor(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	err := wrappedHandler(ctx, "test.subject", msg)

	assertions.ErrorIs(err, expectedErr, "Should return last error")
	assertions.Equal(int32(3), atomic.LoadInt32(&attempts), "Should attempt MaxAttempts + 1 times")
}

func TestRetryMiddlewareNonRetryableError(t *testing.T) {

	assertions := assert.New(t)

	config := RetryConfig{
		MaxAttempts:     3,
		InitialInterval: 1 * time.Millisecond,
		MaxInterval:     10 * time.Millisecond,
		Multiplier:      2.0,
		Jitter:          0,
	}

	mw := NewRetryMiddleware(config)

	var attempts int32

	handler := func(ctx context.Context, subject string, msg *core.Message) error {

		atomic.AddInt32(&attempts, 1)

		return core.ErrInvalidSubject // Non-retryable
	}

	interceptor := mw.InterceptPublish()
	wrappedHandler := interceptor(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	err := wrappedHandler(ctx, "test.subject", msg)

	assertions.ErrorIs(err, core.ErrInvalidSubject, "Should return non-retryable error immediately")
	assertions.Equal(int32(1), atomic.LoadInt32(&attempts), "Should not retry non-retryable errors")
}

func TestRetryMiddlewareContextCancellation(t *testing.T) {

	assertions := assert.New(t)

	config := RetryConfig{
		MaxAttempts:     10,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     1 * time.Second,
		Multiplier:      2.0,
		Jitter:          0,
	}

	mw := NewRetryMiddleware(config)

	var attempts int32

	handler := func(ctx context.Context, subject string, msg *core.Message) error {

		atomic.AddInt32(&attempts, 1)

		return core.ErrNotConnected
	}

	interceptor := mw.InterceptPublish()
	wrappedHandler := interceptor(handler)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	msg := core.NewMessage([]byte("test"))

	err := wrappedHandler(ctx, "test.subject", msg)

	assertions.Error(err, "Should return error on context cancellation")
	assertions.Less(atomic.LoadInt32(&attempts), int32(10), "Should not complete all retries")
}

/* ========================================================================================================
   CUSTOM SHOULD RETRY TESTS
   ======================================================================================================== */

func TestRetryMiddlewareWithShouldRetry(t *testing.T) {

	assertions := assert.New(t)

	config := RetryConfig{
		MaxAttempts:     3,
		InitialInterval: 1 * time.Millisecond,
		MaxInterval:     10 * time.Millisecond,
		Multiplier:      2.0,
		Jitter:          0,
	}

	customError := errors.New("custom error")

	mw := NewRetryMiddleware(config).WithShouldRetry(func(err error) bool {

		return errors.Is(err, customError)
	})

	var attempts int32

	handler := func(ctx context.Context, subject string, msg *core.Message) error {

		count := atomic.AddInt32(&attempts, 1)

		if count < 3 {

			return customError
		}

		return nil
	}

	interceptor := mw.InterceptPublish()
	wrappedHandler := interceptor(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	err := wrappedHandler(ctx, "test.subject", msg)

	assertions.NoError(err, "Should succeed after retries")
	assertions.Equal(int32(3), atomic.LoadInt32(&attempts), "Should retry custom error")
}

/* ========================================================================================================
   BACKOFF CALCULATION TESTS
   ======================================================================================================== */

func TestCalculateBackoff(t *testing.T) {

	assertions := assert.New(t)

	config := RetryConfig{
		MaxAttempts:     3,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     1 * time.Second,
		Multiplier:      2.0,
		Jitter:          0,
	}

	mw := NewRetryMiddleware(config)

	// Attempt 0: 100ms
	backoff0 := mw.calculateBackoff(0)
	assertions.Equal(100*time.Millisecond, backoff0, "Attempt 0 should use initial interval")

	// Attempt 1: 200ms
	backoff1 := mw.calculateBackoff(1)
	assertions.Equal(200*time.Millisecond, backoff1, "Attempt 1 should double")

	// Attempt 2: 400ms
	backoff2 := mw.calculateBackoff(2)
	assertions.Equal(400*time.Millisecond, backoff2, "Attempt 2 should double again")
}

func TestCalculateBackoffCappedAtMax(t *testing.T) {

	assertions := assert.New(t)

	config := RetryConfig{
		MaxAttempts:     10,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     500 * time.Millisecond,
		Multiplier:      2.0,
		Jitter:          0,
	}

	mw := NewRetryMiddleware(config)

	// High attempt should be capped
	backoff := mw.calculateBackoff(10)
	assertions.Equal(500*time.Millisecond, backoff, "Backoff should be capped at max interval")
}

func TestCalculateBackoffWithJitter(t *testing.T) {

	assertions := assert.New(t)

	config := RetryConfig{
		MaxAttempts:     3,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     1 * time.Second,
		Multiplier:      2.0,
		Jitter:          0.5, // 50% jitter
	}

	mw := NewRetryMiddleware(config)

	// Get multiple samples to verify jitter
	var samples []time.Duration

	for i := 0; i < 10; i++ {

		samples = append(samples, mw.calculateBackoff(0))
	}

	// All samples should be in the jitter range
	for _, sample := range samples {

		assertions.GreaterOrEqual(sample, 100*time.Millisecond, "Backoff should not go below initial")
		assertions.LessOrEqual(sample, 150*time.Millisecond, "Backoff should be within jitter range")
	}

	// Verify there is variation (jitter is working)
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

	config := RetryConfig{
		MaxAttempts:     3,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     1 * time.Second,
		Multiplier:      2.0,
		Jitter:          0,
	}

	mw := NewRetryMiddleware(config)

	backoff := mw.calculateBackoff(-1)
	assertions.Equal(100*time.Millisecond, backoff, "Negative attempt should use initial interval")
}

/* ========================================================================================================
   SUBSCRIBE MIDDLEWARE TESTS
   ======================================================================================================== */

func TestRetryMiddlewareInterceptSubscribe(t *testing.T) {

	assertions := assert.New(t)

	config := DefaultRetryConfig()
	mw := NewRetryMiddleware(config)

	var called bool

	handler := func(ctx context.Context, msg *core.Message) error {

		called = true

		return nil
	}

	interceptor := mw.InterceptSubscribe()
	wrappedHandler := interceptor(handler)

	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	err := wrappedHandler(ctx, msg)

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
		{"retryable error", core.ErrNotConnected, true},
		{"retryable timeout", core.ErrConnectionTimeout, true},
		{"non-retryable", core.ErrInvalidSubject, false},
		{"non-retryable config", core.ErrInvalidConfig, false},
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
