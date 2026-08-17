package middleware

import (
	"context"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

/* ========================================================================================================
   RATE LIMITER BASIC TESTS
   ======================================================================================================== */

var (
	rlDenyMsg         = "expected deny after burst exhausted"
	rlAllowAfterRefil = "expected allow after refill"
)

func TestRateLimiterBasic(t *testing.T) {
	assertions := assert.New(t)

	rl := NewRateLimiter(RateLimiterConfig{Rate: 10, Burst: 5})

	for i := 0; i < 5; i++ {
		assertions.True(rl.Allow(), "expected allow on request %d", i)
	}

	assertions.False(rl.Allow(), rlDenyMsg)
}

func TestRateLimiterRefill(t *testing.T) {
	assertions := assert.New(t)

	rl := NewRateLimiter(RateLimiterConfig{Rate: 100, Burst: 2})

	rl.Allow()
	rl.Allow()

	assertions.False(rl.Allow(), rlDenyMsg)

	time.Sleep(25 * time.Millisecond)

	assertions.True(rl.Allow(), rlAllowAfterRefil)
}

/* ========================================================================================================
   ALLOW N TESTS
   ======================================================================================================== */

func TestRateLimiterAllowN(t *testing.T) {
	testCases := []struct {
		name     string
		n        int
		expected bool
	}{
		{"first 5 requests", 5, true},
		{"next 5 requests", 5, true},
		{"exceeds burst", 5, false},
	}

	rl := NewRateLimiter(RateLimiterConfig{Rate: 100, Burst: 10})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, rl.AllowN(tc.n))
		})
	}
}

/* ========================================================================================================
   WAIT TESTS
   ======================================================================================================== */

func TestRateLimiterWait(t *testing.T) {
	assertions := assert.New(t)

	rl := NewRateLimiter(RateLimiterConfig{Rate: 100, Burst: 1})
	rl.Allow()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	err := rl.Wait(ctx)
	elapsed := time.Since(start)

	assertions.NoError(err)
	assertions.Greater(elapsed, 5*time.Millisecond, "should wait before allowing")
}

func TestRateLimiterWaitTimeout(t *testing.T) {
	assertions := assert.New(t)

	rl := NewRateLimiter(RateLimiterConfig{Rate: 1, Burst: 1})
	rl.Allow()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := rl.Wait(ctx)

	assertions.ErrorIs(err, context.DeadlineExceeded, "should return DeadlineExceeded")
}

/* ========================================================================================================
   RATE LIMIT MIDDLEWARE TESTS
   ======================================================================================================== */

func TestRateLimitMiddleware(t *testing.T) {
	assertions := assert.New(t)

	rl := NewRateLimiter(RateLimiterConfig{Rate: 100, Burst: 2})

	callCount := 0
	handler := func(ctx context.Context, msg *nats.Msg) error {
		callCount++
		return nil
	}

	wrapped := RateLimitMiddleware(rl)(handler)
	msg := newTestMsg(testSubject, testPayload)
	ctx := context.Background()

	assertions.NoError(wrapped(ctx, msg))
	assertions.NoError(wrapped(ctx, msg))
	assertions.ErrorIs(wrapped(ctx, msg), ErrRateLimitExceeded, "third call should be rate limited")
	assertions.Equal(2, callCount, "handler should be called twice")
}

/* ========================================================================================================
   PER SUBJECT RATE LIMITER TESTS
   ======================================================================================================== */

func TestPerSubjectRateLimiter(t *testing.T) {
	assertions := assert.New(t)

	psr := NewPerSubjectRateLimiter(RateLimiterConfig{Rate: 100, Burst: 2})

	rl1 := psr.Get(testSubject1)
	rl2 := psr.Get(testSubject2)

	rl1.Allow()
	rl1.Allow()

	assertions.False(rl1.Allow(), "subject1 should be rate limited")
	assertions.True(rl2.Allow(), "subject2 should still allow requests")
}

/* ========================================================================================================
   SLIDING WINDOW LIMITER TESTS
   ======================================================================================================== */

func TestSlidingWindowLimiter(t *testing.T) {
	assertions := assert.New(t)

	sw := NewSlidingWindowLimiter(3, 100*time.Millisecond)

	for i := 0; i < 3; i++ {
		assertions.True(sw.Allow(), "should allow request %d", i)
	}

	assertions.False(sw.Allow(), "should deny after limit reached")
	assertions.Equal(3, sw.Count(), "count should be 3")

	time.Sleep(110 * time.Millisecond)

	assertions.True(sw.Allow(), "should allow after window expired")
}
