package workerpool

import (
	"context"
	"testing"
	"time"
)

// ═══════════════════════════════════════════════════════════════
// BENCHMARK HELPER FUNCTIONS
// ═══════════════════════════════════════════════════════════════

// Simple CPU-bound work simulation
func cpuBoundWork() int {
	sum := 0
	for i := 0; i < 1000; i++ {
		sum += i * i
	}
	return sum
}

// Memory allocation work simulation
func memoryWork() []int {
	data := make([]int, 100)
	for i := range data {
		data[i] = i
	}
	return data
}

// IO-bound work simulation
func ioBoundWork() {
	time.Sleep(time.Microsecond)
}

// Setup helper for benchmarks
func setupBenchmarkPool(minWorkers, maxWorkers int) *Pool {
	pool, _ := newWorkerPool(PoolConfig{
		MinWorkers: minWorkers,
		MaxWorkers: maxWorkers,
		Timeout:    time.Minute,
	})
	return pool
}

// ═══════════════════════════════════════════════════════════════
// BASIC POOL OPERATION BENCHMARKS
// ═══════════════════════════════════════════════════════════════

func BenchmarkWorkerPoolGo(b *testing.B) {
	pool := setupBenchmarkPool(4, 16)
	defer pool.close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = pool.async(func() {
				_ = cpuBoundWork()
			})
		}
	})

	// Wait for all tasks to complete
	time.Sleep(100 * time.Millisecond)
}

func BenchmarkWorkerPoolGoWithContext(b *testing.B) {
	pool := setupBenchmarkPool(4, 16)
	defer pool.close()

	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = pool.asyncWithCtx(ctx, func(c context.Context) {
				_ = cpuBoundWork()
			})
		}
	})

	time.Sleep(100 * time.Millisecond)
}

func BenchmarkWorkerPoolRun(b *testing.B) {
	pool := setupBenchmarkPool(4, 16)
	defer pool.close()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pool.sync(ctx, func() error {
			_ = cpuBoundWork()
			return nil
		})
	}
}

func BenchmarkWorkerPoolRunWithOptions(b *testing.B) {
	pool := setupBenchmarkPool(4, 16)
	defer pool.close()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pool.syncWithOptions(ctx, func(c context.Context) error {
			_ = cpuBoundWork()
			return nil
		}, WithTimeout(time.Second))
	}
}

// BenchmarkInitialization measures Pool initialization overhead
func BenchmarkInitialization(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pool, err := newWorkerPool(DefaultConfig)
		if err != nil {
			b.Fatal(err)
		}
		pool.close()
	}
}

// BenchmarkShutdown measures Pool shutdown time
func BenchmarkShutdown(b *testing.B) {
	pools := make([]*Pool, b.N)

	// Create pools first (not measured)
	for i := 0; i < b.N; i++ {
		pools[i] = setupBenchmarkPool(4, 8)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pools[i].close()
	}
}
