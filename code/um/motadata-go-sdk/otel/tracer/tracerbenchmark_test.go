package tracer

import (
	"context"
	"fmt"
	"testing"
)

/* ---------------------------------------- Span Creation Benchmarks -------------------------------------------------- */

func BenchmarkSpanCreation(b *testing.B) {
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, span := Start(ctx, "benchmark-span")
		span.End()
	}
}

func BenchmarkSpanCreationWithAttributes(b *testing.B) {
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, span := Start(ctx, "benchmark-span",
			WithAttributes(
				StringAttr("key1", "value1"),
				IntAttr("key2", 42),
				BoolAttr("key3", true),
			),
		)
		span.End()
	}
}

func BenchmarkSpanCreation10Attributes(b *testing.B) {
	ctx := context.Background()

	attrs := make([]Attribute, 10)
	for i := 0; i < 10; i++ {
		attrs[i] = IntAttr(fmt.Sprintf("key%d", i), i)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, span := Start(ctx, "benchmark-span",
			WithAttributes(attrs...),
		)
		span.End()
	}
}

func BenchmarkNestedSpanCreation(b *testing.B) {
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx, parent := Start(ctx, "parent")
		_, child := Start(ctx, "child")
		child.End()
		parent.End()
	}
}

/* ---------------------------------------- Span Operation Benchmarks -------------------------------------------------- */

func BenchmarkSpanSetAttributes(b *testing.B) {
	ctx := context.Background()
	_, span := Start(ctx, "benchmark-span")
	defer span.End()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		span.SetAttributes(IntAttr("key", i))
	}
}

func BenchmarkSpanAddEvent(b *testing.B) {
	ctx := context.Background()
	_, span := Start(ctx, "benchmark-span")
	defer span.End()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		span.AddEvent("benchmark-event")
	}
}

/* ---------------------------------------- Context Benchmarks -------------------------------------------------- */

func BenchmarkSpanFromContext(b *testing.B) {
	ctx := context.Background()
	ctx, span := Start(ctx, "benchmark-span")
	defer span.End()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = SpanFromContext(ctx)
	}
}

func BenchmarkTraceIDFromContext(b *testing.B) {
	ctx := context.Background()
	ctx, span := Start(ctx, "benchmark-span")
	defer span.End()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = TraceIDFromContext(ctx)
	}
}

func BenchmarkSpanIDFromContext(b *testing.B) {
	ctx := context.Background()
	ctx, span := Start(ctx, "benchmark-span")
	defer span.End()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = SpanIDFromContext(ctx)
	}
}

/* ---------------------------------------- Concurrent Benchmarks -------------------------------------------------- */

func BenchmarkConcurrentSpanCreation(b *testing.B) {
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, span := Start(ctx, "concurrent-span")
			span.End()
		}
	})
}

func BenchmarkConcurrentNestedSpanCreation(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		ctx := context.Background()
		for pb.Next() {
			ctx, parent := Start(ctx, "parent")
			_, child := Start(ctx, "child")
			child.End()
			parent.End()
		}
	})
}

/* ---------------------------------------- Helper Function Benchmarks -------------------------------------------------- */

func BenchmarkWithSpan(b *testing.B) {
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = WithSpan(ctx, "benchmark-span", func(ctx context.Context) error {
			return nil
		})
	}
}

func BenchmarkTimeSpan(b *testing.B) {
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = TimeSpan(ctx, "benchmark-span", func(ctx context.Context) {})
	}
}

/* ---------------------------------------- Attribute Creation Benchmarks -------------------------------------------------- */

func BenchmarkAttributeString(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = StringAttr("key", "value")
	}
}

func BenchmarkAttributeInt(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = IntAttr("key", 42)
	}
}

func BenchmarkAttributeInt64(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Int64Attr("key", 9223372036854775807)
	}
}

func BenchmarkAttributeFloat64(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Float64Attr("key", 3.14159)
	}
}

func BenchmarkAttributeBool(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = BoolAttr("key", true)
	}
}

func BenchmarkAttributeStrings(b *testing.B) {
	values := []string{"a", "b", "c"}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = StringsAttr("key", values)
	}
}
