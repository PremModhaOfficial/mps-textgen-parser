package middleware

import (
	"context"
	"sync"
	"testing"
)

func BenchmarkMetricsCollectorPublish(b *testing.B) {
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

func BenchmarkMetricsCollectorSubscribe(b *testing.B) {
	mc := NewMetricsCollector()
	wrapped := mc.InterceptSubscribe()(noopSubscribeHandler())
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

func BenchmarkPerSubjectMetricsPublish(b *testing.B) {
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

func BenchmarkMetricsCollectorConcurrent(b *testing.B) {
	mc := NewMetricsCollector()
	wrapped := mc.InterceptPublish()(noopPublishHandler())
	ctx := context.Background()
	msg := newTestMsg(testSubject, testPayload)

	var wg sync.WaitGroup

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			wg.Add(1)
			go func() {
				defer wg.Done()
				wrapped(ctx, msg)
			}()
		}
	})
	wg.Wait()
}
