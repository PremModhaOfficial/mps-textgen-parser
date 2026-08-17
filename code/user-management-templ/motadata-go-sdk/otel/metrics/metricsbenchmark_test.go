package metrics

import (
	"context"
	"fmt"
	"testing"
	"time"
)

/* ---------------------------------------- Memory Allocation Benchmarks -------------------------------------------------- */

func BenchmarkCounterIncNoLabels(b *testing.B) {
	setupBenchmarkMetrics(b)
	defer Shutdown(context.Background())

	testContext := context.Background()
	counter := NewCounter("bench_counter_no_labels", "Benchmark counter")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		counter.Inc(testContext)
	}
}

func BenchmarkCounterIncWithLabels(b *testing.B) {
	setupBenchmarkMetrics(b)
	defer Shutdown(context.Background())

	testContext := context.Background()
	counter := NewCounter("bench_counter_with_labels", "Benchmark counter")
	labels := Labels{"method": "GET", "status": "200"}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		counter.Inc(testContext, labels)
	}
}

func BenchmarkCounterIncManyLabels(b *testing.B) {
	setupBenchmarkMetrics(b)
	defer Shutdown(context.Background())

	testContext := context.Background()
	counter := NewCounter("bench_counter_many_labels", "Benchmark counter")
	labels := Labels{
		"method":   "GET",
		"status":   "200",
		"path":     "/api/users",
		"service":  "user-service",
		"version":  "v1",
		"region":   "us-east-1",
		"instance": "pod-123",
		"handler":  "GetUser",
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		counter.Inc(testContext, labels)
	}
}

func BenchmarkGaugeSetNoLabels(b *testing.B) {
	setupBenchmarkMetrics(b)
	defer Shutdown(context.Background())

	testContext := context.Background()
	gauge := NewGauge("bench_gauge_no_labels", "Benchmark gauge")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gauge.Set(testContext, float64(i))
	}
}

func BenchmarkGaugeSetWithLabels(b *testing.B) {
	setupBenchmarkMetrics(b)
	defer Shutdown(context.Background())

	testContext := context.Background()
	gauge := NewGauge("bench_gauge_with_labels", "Benchmark gauge")
	labels := Labels{"pool": "main", "type": "active"}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gauge.Set(testContext, float64(i), labels)
	}
}

func BenchmarkHistogramObserve(b *testing.B) {
	setupBenchmarkMetrics(b)
	defer Shutdown(context.Background())

	testContext := context.Background()
	histogram := NewHistogram("bench_histogram", "Benchmark histogram")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		histogram.Observe(testContext, float64(i%100))
	}
}

func BenchmarkHistogramObserveWithLabels(b *testing.B) {
	setupBenchmarkMetrics(b)
	defer Shutdown(context.Background())

	testContext := context.Background()
	histogram := NewHistogram("bench_histogram_labels", "Benchmark histogram")
	labels := Labels{"endpoint": "/api/users", "method": "GET"}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		histogram.Observe(testContext, float64(i%100), labels)
	}
}

func BenchmarkHistogramObserveDuration(b *testing.B) {
	setupBenchmarkMetrics(b)
	defer Shutdown(context.Background())

	testContext := context.Background()
	histogram := NewHistogram("bench_histogram_duration", "Benchmark histogram")
	startTime := time.Now()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		histogram.ObserveDuration(testContext, startTime)
	}
}

/* ---------------------------------------- Timer Benchmarks -------------------------------------------------- */

func BenchmarkTimerStartStop(b *testing.B) {
	setupBenchmarkMetrics(b)
	defer Shutdown(context.Background())

	testContext := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		timer := NewTimer(testContext, "bench_timer", "Benchmark timer")
		timer.Stop()
	}
}

func BenchmarkTimerStartStopWithLabels(b *testing.B) {
	setupBenchmarkMetrics(b)
	defer Shutdown(context.Background())

	testContext := context.Background()
	labels := Labels{"operation": "query"}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		timer := NewTimer(testContext, "bench_timer_labels", "Benchmark timer", labels)
		timer.StopWithLabels(Labels{"status": "success"})
	}
}

/* ---------------------------------------- Concurrent Access Benchmarks -------------------------------------------------- */

func BenchmarkCounterConcurrent(b *testing.B) {
	setupBenchmarkMetrics(b)
	defer Shutdown(context.Background())

	testContext := context.Background()
	counter := NewCounter("bench_counter_concurrent", "Benchmark counter")

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Inc(testContext)
		}
	})
}

func BenchmarkCounterConcurrentWithLabels(b *testing.B) {
	setupBenchmarkMetrics(b)
	defer Shutdown(context.Background())

	testContext := context.Background()
	counter := NewCounter("bench_counter_concurrent_labels", "Benchmark counter")

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		labels := Labels{"method": "GET", "status": "200"}
		for pb.Next() {
			counter.Inc(testContext, labels)
		}
	})
}

func BenchmarkGaugeConcurrent(b *testing.B) {
	setupBenchmarkMetrics(b)
	defer Shutdown(context.Background())

	testContext := context.Background()
	gauge := NewGauge("bench_gauge_concurrent", "Benchmark gauge")

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			gauge.Set(testContext, float64(i))
			i++
		}
	})
}

func BenchmarkHistogramConcurrent(b *testing.B) {
	setupBenchmarkMetrics(b)
	defer Shutdown(context.Background())

	testContext := context.Background()
	histogram := NewHistogram("bench_histogram_concurrent", "Benchmark histogram")

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			histogram.Observe(testContext, float64(i%100))
			i++
		}
	})
}

/* ---------------------------------------- Registry Benchmarks -------------------------------------------------- */

func BenchmarkRegistryMetricCreation(b *testing.B) {
	setupBenchmarkMetrics(b)
	defer Shutdown(context.Background())

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Same name should return cached metric
		_ = NewCounter("cached_counter", "Cached counter")
	}
}

func BenchmarkRegistryNewMetricCreation(b *testing.B) {
	setupBenchmarkMetrics(b)
	defer Shutdown(context.Background())

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Unique name each time
		_ = NewCounter(fmt.Sprintf("unique_counter_%d", i), "Unique counter")
	}
}

func BenchmarkNamespaceCreation(b *testing.B) {
	setupBenchmarkMetrics(b)
	defer Shutdown(context.Background())

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Namespace("service")
	}
}

/* ---------------------------------------- Labels Benchmarks -------------------------------------------------- */

func BenchmarkLabelsMerge(b *testing.B) {
	labels1 := Labels{"a": "1", "b": "2", "c": "3"}
	labels2 := Labels{"d": "4", "e": "5", "f": "6"}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = labels1.Merge(labels2)
	}
}

func BenchmarkLabelsCreation(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Labels{
			"method": "GET",
			"status": "200",
			"path":   "/api/users",
		}
	}
}

/* ---------------------------------------- Service Metrics Benchmarks -------------------------------------------------- */

func BenchmarkServiceMetricsRecord(b *testing.B) {
	setupBenchmarkMetrics(b)
	defer Shutdown(context.Background())

	testContext := context.Background()
	serviceMetrics := NewServiceMetrics("user_service")
	duration := 50 * time.Millisecond

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		serviceMetrics.RecordRequest(testContext, "get_user", duration, nil)
	}
}

func BenchmarkServiceMetricsRecordWithError(b *testing.B) {
	setupBenchmarkMetrics(b)
	defer Shutdown(context.Background())

	testContext := context.Background()
	serviceMetrics := NewServiceMetrics("user_service")
	duration := 50 * time.Millisecond
	testError := context.DeadlineExceeded

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		serviceMetrics.RecordRequest(testContext, "get_user", duration, testError)
	}
}

/* ---------------------------------------- Helper Functions -------------------------------------------------- */

func setupBenchmarkMetrics(b *testing.B) {
	b.Helper()

	Shutdown(context.Background())

	metricsConfig := DefaultConfig()
	metricsConfig.Enabled = true
	metricsConfig.OTELEnabled = false

	if err := Init(metricsConfig); err != nil {
		b.Fatalf("Failed to initialize metrics: %v", err)
	}
}
