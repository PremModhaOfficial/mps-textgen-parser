package workerpool

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// ═══════════════════════════════════════════════════════════════
// MANAGER API BENCHMARKS
// ═══════════════════════════════════════════════════════════════

func BenchmarkManagerGo(b *testing.B) {
	resetAndInit(b)
	defer Shutdown()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = Async(func() {
				_ = cpuBoundWork()
			})
		}
	})

	time.Sleep(100 * time.Millisecond)
}

func BenchmarkManagerGoCtx(b *testing.B) {
	resetAndInit(b)
	defer Shutdown()

	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = AsyncWithCtx(ctx, func(c context.Context) {
				_ = cpuBoundWork()
			})
		}
	})

	time.Sleep(100 * time.Millisecond)
}

func BenchmarkManagerRun(b *testing.B) {
	resetAndInit(b)
	defer Shutdown()

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SyncWithCtx(ctx, func() error {
			_ = cpuBoundWork()
			return nil
		})
	}
}

func BenchmarkManagerGoWithResult(b *testing.B) {
	resetAndInit(b)
	defer Shutdown()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultChan := AsyncWithError(func() error {
			_ = cpuBoundWork()
			return nil
		})
		<-resultChan
	}
}

func BenchmarkManagerGoWithValue(b *testing.B) {
	resetAndInit(b)
	defer Shutdown()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultChan := AsyncResult(func() (int, error) {
			_ = cpuBoundWork()
			return 42, nil
		})
		<-resultChan
	}
}

// ═══════════════════════════════════════════════════════════════
// CONCURRENT WORKLOAD BENCHMARKS
// ═══════════════════════════════════════════════════════════════

func BenchmarkConcurrentTasks(b *testing.B) {
	testCases := []struct {
		name       string
		workers    int
		maxWorkers int
	}{
		{"Small-2x4", 2, 4},
		{"Medium-4x8", 4, 8},
		{"Large-8x16", 8, 16},
		{"XLarge-16x32", 16, 32},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			pool := setupBenchmarkPool(tc.workers, tc.maxWorkers)
			defer pool.close()

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_ = pool.async(func() {
						_ = cpuBoundWork()
					})
				}
			})

			time.Sleep(50 * time.Millisecond)
		})
	}
}

func BenchmarkHighFrequencySubmission(b *testing.B) {
	pool := setupBenchmarkPool(8, 32)
	defer pool.close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			_ = pool.async(func() {
				// Very light work to test submission overhead
				_ = j * 2
			})
		}
	}

	time.Sleep(100 * time.Millisecond)
}

func BenchmarkMixedWorkload(b *testing.B) {
	pool := setupBenchmarkPool(4, 16)
	defer pool.close()

	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			switch i % 4 {
			case 0:
				_ = pool.async(func() { _ = cpuBoundWork() })
			case 1:
				_ = pool.asyncWithCtx(ctx, func(c context.Context) { ioBoundWork() })
			case 2:
				_ = pool.sync(ctx, func() error { memoryWork(); return nil })
			case 3:
				_ = pool.async(func() { time.Sleep(time.Microsecond) })
			}
			i++
		}
	})

	time.Sleep(100 * time.Millisecond)
}

// ═══════════════════════════════════════════════════════════════
// SCALING BENCHMARKS
// ═══════════════════════════════════════════════════════════════

func BenchmarkAutoScaling(b *testing.B) {
	pool := setupBenchmarkPool(1, 20)
	defer pool.close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Submit burst of tasks to trigger scaling
		var wg sync.WaitGroup
		wg.Add(50)

		for j := 0; j < 50; j++ {
			_ = pool.async(func() {
				defer wg.Done()
				_ = cpuBoundWork()
			})
		}

		wg.Wait()
	}
}

func BenchmarkPoolUtilization(b *testing.B) {
	pool := setupBenchmarkPool(4, 16)
	defer pool.close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Test different utilization levels
		for load := 1; load <= 16; load += 3 {
			var wg sync.WaitGroup
			wg.Add(load)

			for j := 0; j < load; j++ {
				_ = pool.async(func() {
					defer wg.Done()
					_ = cpuBoundWork()
				})
			}

			wg.Wait()
		}
	}
}

// ═══════════════════════════════════════════════════════════════
// BATCH OPERATION BENCHMARKS
// ═══════════════════════════════════════════════════════════════

func BenchmarkBatchOperations(b *testing.B) {
	batchSizes := []int{10, 50, 100, 500}

	for _, size := range batchSizes {
		b.Run(fmt.Sprintf("Batch-%d", size), func(b *testing.B) {
			resetAndInit(b)
			defer Shutdown()

			ctx := context.Background()

			// Create tasks
			tasks := make([]func() error, size)
			for i := 0; i < size; i++ {
				tasks[i] = func() error {
					_ = cpuBoundWork()
					return nil
				}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				Batch(ctx, tasks...)
			}
		})
	}
}

func BenchmarkBatchWithContext(b *testing.B) {
	resetAndInit(b)
	defer Shutdown()

	ctx := context.Background()

	tasks := make([]func(context.Context) error, 50)
	for i := 0; i < 50; i++ {
		tasks[i] = func(c context.Context) error {
			_ = cpuBoundWork()
			return nil
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BatchWithContext(ctx, tasks...)
	}
}

func BenchmarkParallelVsSequential(b *testing.B) {
	resetAndInit(b)
	defer Shutdown()

	ctx := context.Background()

	tasks := make([]func() error, 20)
	for i := 0; i < 20; i++ {
		tasks[i] = func() error {
			_ = cpuBoundWork()
			return nil
		}
	}

	b.Run("Parallel", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Parallel(ctx, tasks...)
		}
	})

	b.Run("Sequential", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Sequential(ctx, tasks...)
		}
	})
}

// ═══════════════════════════════════════════════════════════════
// MEMORY ALLOCATION BENCHMARKS
// ═══════════════════════════════════════════════════════════════

func BenchmarkMemoryAllocation(b *testing.B) {
	pool := setupBenchmarkPool(4, 8)
	defer pool.close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = pool.async(func() {
				data := memoryWork()
				_ = data
			})
		}
	})

	time.Sleep(100 * time.Millisecond)
}

func BenchmarkChannelAllocation(b *testing.B) {
	resetAndInit(b)
	defer Shutdown()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		resultChan := AsyncResult(func() ([]int, error) {
			return memoryWork(), nil
		})
		result := <-resultChan
		_ = result
	}
}

// ═══════════════════════════════════════════════════════════════
// COMPARATIVE BENCHMARKS
// ═══════════════════════════════════════════════════════════════

func BenchmarkVsDirectGoroutines(b *testing.B) {
	pool := setupBenchmarkPool(8, 16)
	defer pool.close()

	b.Run("WorkerPool", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ = pool.async(func() {
					_ = cpuBoundWork()
				})
			}
		})
		time.Sleep(50 * time.Millisecond)
	})

	b.Run("DirectGoroutines", func(b *testing.B) {
		var wg sync.WaitGroup
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				wg.Add(1)
				go func() {
					defer wg.Done()
					_ = cpuBoundWork()
				}()
			}
		})
		wg.Wait()
	})
}

func BenchmarkPoolSizes(b *testing.B) {
	poolSizes := []struct {
		min, max int
		name     string
	}{
		{1, 2, "Tiny"},
		{2, 4, "Small"},
		{4, 8, "Medium"},
		{8, 16, "Large"},
		{16, 32, "XLarge"},
	}

	for _, ps := range poolSizes {
		b.Run(ps.name, func(b *testing.B) {
			pool := setupBenchmarkPool(ps.min, ps.max)
			defer pool.close()

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_ = pool.async(func() {
						_ = cpuBoundWork()
					})
				}
			})

			time.Sleep(50 * time.Millisecond)
		})
	}
}

// ═══════════════════════════════════════════════════════════════
// SCHEDULING BENCHMARKS
// ═══════════════════════════════════════════════════════════════

func BenchmarkSchedulingAfter(b *testing.B) {
	resetAndInit(b)
	defer Shutdown()

	var counter atomic.Int64

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task := After(time.Nanosecond, func() {
			counter.Add(1)
		})
		_ = task
	}

	time.Sleep(10 * time.Millisecond)
}

func BenchmarkSchedulingEvery(b *testing.B) {
	resetAndInit(b)
	defer Shutdown()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	var counter atomic.Int64

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task := Every(ctx, time.Microsecond, func() {
			counter.Add(1)
		})
		time.Sleep(10 * time.Microsecond)
		task.Cancel()
	}
}

// ═══════════════════════════════════════════════════════════════
// STRESS TEST BENCHMARKS
// ═══════════════════════════════════════════════════════════════

func BenchmarkStressTest(b *testing.B) {
	pool := setupBenchmarkPool(runtime.NumCPU(), runtime.NumCPU()*4)
	defer pool.close()

	var completed atomic.Int64
	ctx := context.Background()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Mix of different operations
			switch completed.Add(1) % 5 {
			case 0:
				_ = pool.async(func() { _ = cpuBoundWork() })
			case 1:
				_ = pool.asyncWithCtx(ctx, func(c context.Context) { ioBoundWork() })
			case 2:
				_ = pool.sync(ctx, func() error { memoryWork(); return nil })
			case 3:
				resultChan := AsyncWithError(func() error {
					_ = cpuBoundWork()
					return nil
				})
				<-resultChan
			case 4:
				valueChan := AsyncResult(func() (int, error) {
					return cpuBoundWork(), nil
				})
				<-valueChan
			}
		}
	})

	time.Sleep(100 * time.Millisecond)
}

func BenchmarkLongRunningTasks(b *testing.B) {
	pool := setupBenchmarkPool(2, 8)
	defer pool.close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(5)

		for j := 0; j < 5; j++ {
			_ = pool.async(func() {
				defer wg.Done()
				// Simulate longer running task
				for k := 0; k < 10000; k++ {
					_ = cpuBoundWork()
				}
			})
		}

		wg.Wait()
	}
}

// ═══════════════════════════════════════════════════════════════
// BENCHMARK HELPER UTILITIES
// ═══════════════════════════════════════════════════════════════

func resetAndInit(b *testing.B) {
	b.Helper()

	// Reset global state
	instance.Store(nil)
	once = sync.Once{}
	initErr = nil
	initialized.Store(false)

	err := Init(PoolConfig{
		MinWorkers: 4,
		MaxWorkers: 16,
		Timeout:    time.Minute,
	})
	if err != nil {
		b.Fatal("Failed to initialize Pool:", err)
	}
}
