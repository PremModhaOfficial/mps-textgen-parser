package middleware

import (
	"context"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
	"testing"

	"github.com/nats-io/nats.go"
)

func BenchmarkTracingMiddlewarePublish(b *testing.B) {
	mw := NewTracingMiddleware()
	mw.UseOTEL = false
	wrapped := mw.InterceptPublish()(noopPublishHandler())
	ctx := context.Background()
	msg := &nats.Msg{Data: testPayload, Header: make(nats.Header), Subject: testSubject}

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			wrapped(ctx, msg)
		}
	})
}

func BenchmarkTracingMiddlewareSubscribe(b *testing.B) {
	mw := NewTracingMiddleware()
	mw.UseOTEL = false
	wrapped := mw.InterceptSubscribe()(noopSubscribeHandler())
	ctx := context.Background()
	msg := &nats.Msg{Data: testPayload, Header: make(nats.Header), Subject: testSubject}

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			wrapped(ctx, msg)
		}
	})
}

func BenchmarkDefaultIDGenerator(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			defaultIDGenerator()
		}
	})
}

func BenchmarkExtractW3CTraceParent(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ExtractW3CTraceParent(testW3CParent)
		}
	})
}

func BenchmarkFormatW3CTraceParent(b *testing.B) {
	tc := &core.TraceContext{TraceID: testW3CTraceID, SpanID: testW3CSpanID, Sampled: true}

	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			FormatW3CTraceParent(tc)
		}
	})
}
