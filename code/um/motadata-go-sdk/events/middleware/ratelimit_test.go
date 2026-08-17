package middleware

import (
	"context"
	"testing"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

func TestRateLimiterBasic(t *testing.T) {
	cfg := RateLimiterConfig{
		Rate:  10,
		Burst: 5,
	}
	rl := NewRateLimiter(cfg)

	// Should allow burst
	for i := 0; i < 5; i++ {
		if !rl.Allow() {
			t.Errorf("expected allow on request %d", i)
		}
	}

	// Should deny after burst
	if rl.Allow() {
		t.Error("expected deny after burst exhausted")
	}
}

func TestRateLimiterRefill(t *testing.T) {
	cfg := RateLimiterConfig{
		Rate:  100, // 100 per second = 1 per 10ms
		Burst: 2,
	}
	rl := NewRateLimiter(cfg)

	// Exhaust burst
	rl.Allow()
	rl.Allow()

	if rl.Allow() {
		t.Error("expected deny after burst exhausted")
	}

	// Wait for refill
	time.Sleep(25 * time.Millisecond)

	// Should allow after refill
	if !rl.Allow() {
		t.Error("expected allow after refill")
	}
}

func TestRateLimiterAllowN(t *testing.T) {
	cfg := RateLimiterConfig{
		Rate:  100,
		Burst: 10,
	}
	rl := NewRateLimiter(cfg)

	// Request 5 at once
	if !rl.AllowN(5) {
		t.Error("expected allow for 5 requests")
	}

	// Request 5 more
	if !rl.AllowN(5) {
		t.Error("expected allow for next 5 requests")
	}

	// Should deny more
	if rl.AllowN(5) {
		t.Error("expected deny after burst exhausted")
	}
}

func TestRateLimiterWait(t *testing.T) {
	cfg := RateLimiterConfig{
		Rate:  100, // 100 per second
		Burst: 1,
	}
	rl := NewRateLimiter(cfg)

	// Use the burst
	rl.Allow()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	err := rl.Wait(ctx)
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	if elapsed < 5*time.Millisecond {
		t.Errorf("expected wait time > 5ms, got %v", elapsed)
	}
}

func TestRateLimiterWaitTimeout(t *testing.T) {
	cfg := RateLimiterConfig{
		Rate:  1,  // 1 per second
		Burst: 1,
	}
	rl := NewRateLimiter(cfg)

	// Use the burst
	rl.Allow()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := rl.Wait(ctx)
	if err != context.DeadlineExceeded {
		t.Errorf("expected DeadlineExceeded, got %v", err)
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	cfg := RateLimiterConfig{
		Rate:  100,
		Burst: 2,
	}
	rl := NewRateLimiter(cfg)

	callCount := 0
	handler := func(ctx context.Context, subject string, msg *core.Message) error {
		callCount++
		return nil
	}

	wrapped := RateLimitMiddleware(rl)(handler)
	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	// First two should succeed
	if err := wrapped(ctx, "test", msg); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if err := wrapped(ctx, "test", msg); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}

	// Third should be rate limited
	if err := wrapped(ctx, "test", msg); err != ErrRateLimitExceeded {
		t.Errorf("expected ErrRateLimitExceeded, got %v", err)
	}

	if callCount != 2 {
		t.Errorf("expected handler called 2 times, got %d", callCount)
	}
}

func TestPerSubjectRateLimiter(t *testing.T) {
	cfg := RateLimiterConfig{
		Rate:  100,
		Burst: 2,
	}
	psr := NewPerSubjectRateLimiter(cfg)

	// Get limiters for different subjects
	rl1 := psr.Get("subject1")
	rl2 := psr.Get("subject2")

	// Exhaust subject1
	rl1.Allow()
	rl1.Allow()

	if rl1.Allow() {
		t.Error("expected subject1 to be rate limited")
	}

	// subject2 should still work
	if !rl2.Allow() {
		t.Error("expected subject2 to allow requests")
	}
}

func TestSlidingWindowLimiter(t *testing.T) {
	sw := NewSlidingWindowLimiter(3, 100*time.Millisecond)

	// Should allow first 3
	for i := 0; i < 3; i++ {
		if !sw.Allow() {
			t.Errorf("expected allow on request %d", i)
		}
	}

	// Should deny 4th
	if sw.Allow() {
		t.Error("expected deny after limit reached")
	}

	// Count should be 3
	if sw.Count() != 3 {
		t.Errorf("expected count 3, got %d", sw.Count())
	}

	// Wait for window to slide
	time.Sleep(110 * time.Millisecond)

	// Should allow again
	if !sw.Allow() {
		t.Error("expected allow after window expired")
	}
}

func BenchmarkRateLimiterAllow(b *testing.B) {
	rl := NewRateLimiter(RateLimiterConfig{
		Rate:  1000000,
		Burst: 1000000,
	})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		rl.Allow()
	}
}

func BenchmarkSlidingWindowAllow(b *testing.B) {
	sw := NewSlidingWindowLimiter(1000000, time.Minute)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		sw.Allow()
	}
}
