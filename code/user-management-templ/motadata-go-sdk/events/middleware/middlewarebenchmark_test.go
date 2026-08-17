package middleware

import (
	"context"
	"testing"

	"github.com/nats-io/nats.go"
)

func BenchmarkChain(b *testing.B) {
	middleware := func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, msg *nats.Msg) error {
			return next(ctx, msg)
		}
	}

	chained := Chain(middleware, middleware, middleware)(noopPublishHandler())
	ctx := context.Background()
	msg := newTestMsg(testSubject, testPayload)

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			chained(ctx, msg)
		}
	})
}

func BenchmarkStack(b *testing.B) {
	s := NewStack()
	s.UsePublish(func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, msg *nats.Msg) error {
			return next(ctx, msg)
		}
	})

	wrapped := s.WrapPublish(noopPublishHandler())
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

func BenchmarkRetryMiddleware(b *testing.B) {
	rm := NewRetryMiddleware(DefaultRetryConfig())
	wrapped := rm.InterceptPublish()(noopPublishHandler())
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

func BenchmarkMetricsCollector(b *testing.B) {
	mc := NewMetricsCollector()
	wrapped := mc.InterceptPublish()(noopPublishHandler())
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

func BenchmarkPerSubjectMetrics(b *testing.B) {
	psm := NewPerSubjectMetrics()
	wrapped := psm.InterceptPublish()(noopPublishHandler())
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

func BenchmarkCircuitBreaker(b *testing.B) {
	cb := NewCircuitBreaker(DefaultCircuitBreakerConfig())
	wrapped := CircuitBreakerMiddleware(cb)(noopPublishHandler())
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

func BenchmarkRateLimiter(b *testing.B) {
	rl := NewRateLimiter(RateLimiterConfig{Rate: 1000000, Burst: 1000})
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

func BenchmarkLoggingMiddleware(b *testing.B) {
	lm := NewLoggingMiddleware(NoOpLogger{})
	wrapped := lm.InterceptPublish()(noopPublishHandler())
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
