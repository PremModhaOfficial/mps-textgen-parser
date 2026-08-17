package middleware

import (
	"context"
	"testing"
	"time"
)

func BenchmarkRetryMiddlewareSuccess(b *testing.B) {
	mw := NewRetryMiddleware(DefaultRetryConfig())
	wrapped := mw.InterceptPublish()(noopPublishHandler())
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

func BenchmarkCalculateBackoff(b *testing.B) {
	mw := NewRetryMiddleware(RetryConfig{
		MaxAttempts:     3,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     1 * time.Second,
		Multiplier:      retryDefaultMult,
		Jitter:          0.1,
	})

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			mw.calculateBackoff(i % 5)
			i++
		}
	})
}
