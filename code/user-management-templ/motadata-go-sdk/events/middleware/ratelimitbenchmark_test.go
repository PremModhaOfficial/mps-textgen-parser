package middleware

import (
	"context"
	"testing"
	"time"
)

func BenchmarkRateLimiterAllow(b *testing.B) {
	rl := NewRateLimiter(RateLimiterConfig{Rate: 1000000, Burst: 1000000})

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rl.Allow()
		}
	})
}

func BenchmarkSlidingWindowAllow(b *testing.B) {
	sw := NewSlidingWindowLimiter(1000000, time.Minute)

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			sw.Allow()
		}
	})
}

func BenchmarkRateLimitMiddlewarePublish(b *testing.B) {
	rl := NewRateLimiter(RateLimiterConfig{Rate: 1000000, Burst: 1000000})
	wrapped := RateLimitMiddleware(rl)(noopPublishHandler())
	ctx := context.Background()
	msg := newTestMsg(testSubject, testPayload)

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			wrapped(ctx, msg)
		}
	})
}

func BenchmarkPerSubjectRateLimitMiddlewarePublish(b *testing.B) {
	psr := NewPerSubjectRateLimiter(RateLimiterConfig{Rate: 1000000, Burst: 1000000})
	wrapped := PerSubjectRateLimitMiddleware(psr)(noopPublishHandler())
	ctx := context.Background()
	msg := newTestMsg(testSubject, testPayload)

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			wrapped(ctx, msg)
		}
	})
}
