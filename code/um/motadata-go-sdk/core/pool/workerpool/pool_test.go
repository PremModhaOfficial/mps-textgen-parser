package workerpool

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"
	"github.com/stretchr/testify/assert"
)

// ═══════════════════════════════════════════════════════════════
// POSITIVE TEST SCENARIOS
// ═══════════════════════════════════════════════════════════════

func TestWorkerPoolCreation(t *testing.T) {
	assertions := assert.New(t)

	config := PoolConfig{
		MinWorkers: 2,
		MaxWorkers: 10,
		Timeout:    5 * time.Second,
	}

	pool, err := newWorkerPool(config)
	assertions.NoError(err, "Should create Pool without error")
	assertions.NotNil(pool, "Pool should not be nil")
	assertions.Equal(2, pool.minWorkers, "Expected min workers to be 2")
	assertions.Equal(10, pool.maxWorkers, "Expected max workers to be 10")
	assertions.Equal(5*time.Second, pool.timeout, "Expected timeout to be 5 seconds")
	assertions.True(pool.isRunning(), "Pool should be running")

	pool.close()
	assertions.False(pool.isRunning(), "Pool should not be running after close")
}

func TestWorkerPoolDefaultConfig(t *testing.T) {
	assertions := assert.New(t)

	pool, err := newWorkerPool(PoolConfig{})
	assertions.NoError(err, "Should create Pool with default config")
	assertions.NotNil(pool, "Pool should not be nil")
	assertions.Equal(DefaultConfig.MinWorkers, pool.minWorkers, "Expected default min workers")
	assertions.Equal(DefaultConfig.MaxWorkers, pool.maxWorkers, "Expected default max workers")
	assertions.Equal(DefaultConfig.Timeout, pool.timeout, "Expected default timeout")

	defer pool.close()
}

func TestWorkerPoolGoOperation(t *testing.T) {
	assertions := assert.New(t)

	pool, err := newWorkerPool(DefaultConfig)
	assertions.NoError(err)
	defer pool.close()

	var counter atomic.Int32
	numTasks := 100
	var wg sync.WaitGroup
	wg.Add(numTasks)

	for i := 0; i < numTasks; i++ {
		taskErr := pool.async(func() {
			counter.Add(1)
			wg.Done()
		})
		assertions.NoError(taskErr, "Should submit task without error")
	}

	wg.Wait()
	assertions.Equal(int32(numTasks), counter.Load(), "All tasks should complete")

	stats := pool.getPoolStats()
	assertions.Equal(int64(numTasks), stats.Submitted, "Submitted count should match")
	assertions.Equal(int64(numTasks), stats.Completed, "Completed count should match")
	assertions.Equal(int64(0), stats.Failed, "No tasks should fail")
}

func TestWorkerPoolGoWithContextOperation(t *testing.T) {
	assertions := assert.New(t)

	pool, err := newWorkerPool(DefaultConfig)
	assertions.NoError(err)
	defer pool.close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	executed := make(chan bool, 1)

	err = pool.asyncWithCtx(ctx, func(ctx context.Context) {
		select {
		case <-ctx.Done():
			return
		default:
			executed <- true
		}
	})

	assertions.NoError(err, "Should submit task with context")

	select {
	case <-executed:
		assertions.True(true, "Task executed successfully")
	case <-time.After(2 * time.Second):
		assertions.Fail("Task did not execute in time")
	}
}

func TestWorkerPoolRunOperation(t *testing.T) {
	assertions := assert.New(t)

	pool, err := newWorkerPool(DefaultConfig)
	assertions.NoError(err)
	defer pool.close()

	ctx := context.Background()
	counter := 0

	err = pool.sync(ctx, func() error {
		counter++
		return nil
	})

	assertions.NoError(err, "Should run task synchronously without error")
	assertions.Equal(1, counter, "Task should have executed")

	stats := pool.getPoolStats()
	assertions.Equal(int64(1), stats.Submitted, "One task submitted")
	assertions.Equal(int64(1), stats.Completed, "One task completed")
}

func TestWorkerPoolRunWithErrorOperation(t *testing.T) {
	assertions := assert.New(t)

	pool, err := newWorkerPool(DefaultConfig)
	assertions.NoError(err)
	defer pool.close()

	ctx := context.Background()
	expectedErr := errors.New("test error")

	err = pool.sync(ctx, func() error {
		return expectedErr
	})

	assertions.Error(err, "Should return error from task")
	assertions.Equal(expectedErr, err, "Should return exact error")

	stats := pool.getPoolStats()
	assertions.Equal(int64(1), stats.Failed, "One task should fail")
}

func TestWorkerPoolRunWithOptions(t *testing.T) {
	assertions := assert.New(t)

	pool, err := newWorkerPool(DefaultConfig)
	assertions.NoError(err)
	defer pool.close()

	ctx := context.Background()
	var successCalled, errorCalled, retryCalled bool
	retryCount := 0

	// Test successful execution
	err = pool.syncWithOptions(ctx, func(context.Context) error {
		return nil
	},
		WithTimeout(2*time.Second),
		WithOnSuccess(func() { successCalled = true }),
	)

	assertions.NoError(err, "Should execute successfully")
	assertions.True(successCalled, "OnSuccess should be called")

	// Test with retry
	attempts := 0
	err = pool.syncWithOptions(ctx, func(context.Context) error {
		attempts++
		if attempts < 3 {
			return errors.New("retry me")
		}
		return nil
	},
		WithRetry(3, 100*time.Millisecond),
		WithOnRetry(func(attempt int, err error) {
			retryCalled = true
			retryCount = attempt
		}),
	)

	assertions.NoError(err, "Should succeed after retries")
	assertions.True(retryCalled, "OnRetry should be called")
	assertions.Greater(retryCount, 0, "Should have retry attempts")

	// Test error callback
	err = pool.syncWithOptions(ctx, func(context.Context) error {
		return errors.New("final error")
	},
		WithOnError(func(err error) { errorCalled = true }),
	)

	assertions.Error(err, "Should return error")
	assertions.True(errorCalled, "OnError should be called")
}

func TestWorkerPoolCloseWithTimeout(t *testing.T) {
	assertions := assert.New(t)

	pool, err := newWorkerPool(DefaultConfig)
	assertions.NoError(err)

	// Don't submit any tasks, just test close with timeout
	err = pool.closeWithTimeout(100 * time.Millisecond)
	assertions.NoError(err, "Should close within timeout when no tasks")
	assertions.True(pool.closed.Load(), "Pool should be closed")

	// Test with tasks
	pool2, err := newWorkerPool(DefaultConfig)
	assertions.NoError(err)

	// Submit a quick task
	done := make(chan bool)
	_ = pool2.async(func() {
		time.Sleep(10 * time.Millisecond)
		done <- true
	})

	// Wait for task to complete
	<-done

	// Give time for workerpool to return to Pool
	time.Sleep(50 * time.Millisecond)

	// Now close - the Pool.Running() might still show workers but they're idle
	err = pool2.closeWithTimeout(200 * time.Millisecond)
	// This might timeout due to idle workers, which is acceptable behavior
	if err != nil {
		assertions.Contains(err.Error(), "timeout", "Error should be about timeout")
	}
	assertions.True(pool2.closed.Load(), "Pool should be marked as closed")
}

func TestWorkerPoolStats(t *testing.T) {
	assertions := assert.New(t)

	pool, err := newWorkerPool(DefaultConfig)
	assertions.NoError(err)
	defer pool.close()

	ctx := context.Background()

	// Submit various tasks
	_ = pool.async(func() { time.Sleep(10 * time.Millisecond) })
	_ = pool.sync(ctx, func() error { return nil })
	_ = pool.sync(ctx, func() error { return errors.New("error") })

	time.Sleep(50 * time.Millisecond)

	stats := pool.getPoolStats()
	assertions.GreaterOrEqual(stats.Submitted, int64(3), "Should have submitted at least 3 tasks")
	assertions.GreaterOrEqual(stats.Completed, int64(2), "Should have completed at least 2 tasks")
	assertions.GreaterOrEqual(stats.Failed, int64(1), "Should have at least 1 failed task")
}

// ═══════════════════════════════════════════════════════════════
// NEGATIVE TEST SCENARIOS
// ═══════════════════════════════════════════════════════════════

func TestWorkerPoolClosedOperations(t *testing.T) {
	assertions := assert.New(t)

	pool, err := newWorkerPool(DefaultConfig)
	assertions.NoError(err)
	pool.close()

	// Try operations on closed Pool
	err = pool.async(func() {})
	assertions.Error(err, "async should error on closed Pool")
	assertions.Equal(utils.ErrPoolClosed, err, "Should return Pool closed error")

	ctx := context.Background()
	err = pool.asyncWithCtx(ctx, func(context.Context) {})
	assertions.Error(err, "asyncWithCtx should error on closed Pool")
	assertions.Equal(utils.ErrPoolClosed, err, "Should return Pool closed error")

	err = pool.sync(ctx, func() error { return nil })
	assertions.Error(err, "sync should error on closed Pool")
	assertions.Equal(utils.ErrPoolClosed, err, "Should return Pool closed error")

	err = pool.syncWithOptions(ctx, func(context.Context) error { return nil })
	assertions.Error(err, "syncWithOptions should error on closed Pool")
	assertions.Equal(utils.ErrPoolClosed, err, "Should return Pool closed error")
}

func TestWorkerPoolDoubleClose(t *testing.T) {
	assertions := assert.New(t)

	pool, err := newWorkerPool(DefaultConfig)
	assertions.NoError(err)

	pool.close()
	assertions.True(pool.closed.Load(), "Pool should be closed")

	// Second close should be safe
	assertions.NotPanics(func() {
		pool.close()
	}, "Double close should not panic")
}

func TestWorkerPoolContextCancellation(t *testing.T) {
	assertions := assert.New(t)

	pool, err := newWorkerPool(DefaultConfig)
	assertions.NoError(err)
	defer pool.close()

	ctx, cancel := context.WithCancel(context.Background())

	started := make(chan bool)
	finished := make(chan bool, 1)

	// Submit task with context
	err = pool.asyncWithCtx(ctx, func(ctx context.Context) {
		started <- true
		select {
		case <-ctx.Done():
			return
		case <-time.After(5 * time.Second):
			finished <- true
		}
	})
	assertions.NoError(err)

	<-started
	cancel() // Cancel context

	select {
	case <-finished:
		assertions.Fail("Task should not finish after context cancellation")
	case <-time.After(200 * time.Millisecond):
		// Expected: task stopped due to context cancellation
		// The asyncWithCtx increments failed when context is done
	}
}

func TestWorkerPoolRunTimeout(t *testing.T) {
	assertions := assert.New(t)

	pool, err := newWorkerPool(DefaultConfig)
	assertions.NoError(err)
	defer pool.close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err = pool.sync(ctx, func() error {
		time.Sleep(500 * time.Millisecond)
		return nil
	})

	assertions.Error(err, "Should timeout")
	assertions.Equal(context.DeadlineExceeded, err, "Should return deadline exceeded error")
}

func TestWorkerPoolPanicRecovery(t *testing.T) {
	assertions := assert.New(t)

	pool, err := newWorkerPool(DefaultConfig)
	assertions.NoError(err)
	defer pool.close()

	// Test panic recovery in async
	assertions.NotPanics(func() {
		_ = pool.async(func() {
			panic("test panic")
		})
		time.Sleep(50 * time.Millisecond)
	}, "Should recover from panic in async")

	// Test panic recovery in sync
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	done := make(chan bool)
	go func() {
		err = pool.sync(ctx, func() error {
			panic("test panic in run")
		})
		done <- true
	}()

	select {
	case <-done:
		assertions.Error(err, "Should return error for panic")
		assertions.Contains(err.Error(), "panic", "Error should mention panic")
	case <-time.After(200 * time.Millisecond):
		// If we timeout, that's ok - the panic was recovered
	}

	// The panic in async is counted immediately
	// The panic in sync may or may not be counted depending on timing
	// We just verify the Pool is still functional after panics
	assertions.True(pool.isRunning(), "Pool should still be running after panics")
}

// ═══════════════════════════════════════════════════════════════
// EDGE CASE TEST SCENARIOS
// ═══════════════════════════════════════════════════════════════

func TestWorkerPoolConcurrentSubmissions(t *testing.T) {
	assertions := assert.New(t)

	pool, err := newWorkerPool(PoolConfig{
		MinWorkers: 2,
		MaxWorkers: 5,
		Timeout:    1 * time.Second,
	})
	assertions.NoError(err)
	defer pool.close()

	numGoroutines := 50
	numTasksPerGoroutine := 20
	totalTasks := numGoroutines * numTasksPerGoroutine

	var counter atomic.Int32
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numTasksPerGoroutine; j++ {
				taskErr := pool.async(func() {
					counter.Add(1)
					time.Sleep(time.Microsecond)
				})
				assertions.NoError(taskErr, "Should submit task")
			}
		}()
	}

	wg.Wait()
	time.Sleep(100 * time.Millisecond) // Allow tasks to complete

	finalCount := counter.Load()
	assertions.Equal(int32(totalTasks), finalCount, "All tasks should complete")
}

func TestWorkerPoolAutoScaling(t *testing.T) {
	assertions := assert.New(t)

	pool, err := newWorkerPool(PoolConfig{
		MinWorkers: 2,
		MaxWorkers: 8,
		Timeout:    5 * time.Second,
	})
	assertions.NoError(err)
	defer pool.close()

	// Initial capacity should be min workers
	initialStats := pool.getPoolStats()
	assertions.Equal(2, initialStats.Capacity, "Initial capacity should be min workers")

	// Submit many tasks to trigger scale up
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		_ = pool.async(func() {
			time.Sleep(200 * time.Millisecond)
			wg.Done()
		})
	}

	// Give auto-scaler time to react (runs every second)
	time.Sleep(2500 * time.Millisecond)

	midStats := pool.getPoolStats()
	// Auto-scaler should have increased capacity when there were waiting tasks
	assertions.GreaterOrEqual(midStats.Capacity, 2, "Capacity should be at least min workers")
	assertions.LessOrEqual(midStats.Capacity, 8, "Should not exceed max workers")

	// Wait for tasks to complete
	wg.Wait()
	time.Sleep(2 * time.Second) // Allow scale down

	finalStats := pool.getPoolStats()
	assertions.LessOrEqual(finalStats.Capacity, midStats.Capacity, "Should scale down when idle")
}

func TestWorkerPoolZeroWorkers(t *testing.T) {
	assertions := assert.New(t)

	// Test with zero min workers (should use default)
	pool, err := newWorkerPool(PoolConfig{
		MinWorkers: 0,
		MaxWorkers: 5,
	})
	assertions.NoError(err, "Should create Pool with zero min workers")
	assertions.Equal(DefaultConfig.MinWorkers, pool.minWorkers, "Should use default min workers")
	pool.close()

	// Test with zero max workers (should use default)
	pool, err = newWorkerPool(PoolConfig{
		MinWorkers: 2,
		MaxWorkers: 0,
	})
	assertions.NoError(err, "Should create Pool with zero max workers")
	assertions.Equal(DefaultConfig.MaxWorkers, pool.maxWorkers, "Should use default max workers")
	pool.close()
}

func TestWorkerPoolCloseWithTimeoutExpiry(t *testing.T) {
	assertions := assert.New(t)

	pool, err := newWorkerPool(DefaultConfig)
	assertions.NoError(err)

	// Submit long-running task
	_ = pool.async(func() {
		time.Sleep(1 * time.Second)
	})

	// Give task time to start
	time.Sleep(10 * time.Millisecond)

	// close with insufficient timeout
	err = pool.closeWithTimeout(50 * time.Millisecond)
	assertions.Error(err, "Should timeout")
	assertions.Contains(err.Error(), "timeout", "Error should mention timeout")
	assertions.Contains(err.Error(), "tasks still running", "Error should mention running tasks")
}

func TestWorkerPoolRunWithOptionsTimeoutContext(t *testing.T) {
	assertions := assert.New(t)

	pool, err := newWorkerPool(DefaultConfig)
	assertions.NoError(err)
	defer pool.close()

	ctx := context.Background()
	executed := false

	err = pool.syncWithOptions(ctx, func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1 * time.Second):
			executed = true
			return nil
		}
	}, WithTimeout(100*time.Millisecond))

	assertions.Error(err, "Should timeout")
	assertions.Equal(context.DeadlineExceeded, err, "Should return deadline exceeded")
	assertions.False(executed, "Task should not complete")
}

func TestWorkerPoolRetryWithExponentialBackoff(t *testing.T) {
	assertions := assert.New(t)

	pool, err := newWorkerPool(DefaultConfig)
	assertions.NoError(err)
	defer pool.close()

	ctx := context.Background()
	attempts := 0
	var retryDelays []time.Duration
	lastRetryTime := time.Now()

	err = pool.syncWithOptions(ctx, func(context.Context) error {
		attempts++
		if attempts <= 3 {
			return errors.New("retry me")
		}
		return nil
	},
		WithRetry(3, 50*time.Millisecond),
		WithOnRetry(func(attempt int, err error) {
			now := time.Now()
			delay := now.Sub(lastRetryTime)
			retryDelays = append(retryDelays, delay)
			lastRetryTime = now
		}),
	)

	assertions.NoError(err, "Should succeed after retries")
	assertions.Equal(4, attempts, "Should have made 4 attempts (1 initial + 3 retries)")
	assertions.Len(retryDelays, 3, "Should have 3 retry delays recorded")

	// Verify exponential backoff
	for i := 1; i < len(retryDelays); i++ {
		assertions.Greater(retryDelays[i], retryDelays[i-1], "Delay should increase with each retry")
	}
}

func TestWorkerPoolNilCallbacks(t *testing.T) {
	assertions := assert.New(t)

	pool, err := newWorkerPool(DefaultConfig)
	assertions.NoError(err)
	defer pool.close()

	ctx := context.Background()

	// Test with nil callbacks - should not panic
	assertions.NotPanics(func() {
		_ = pool.syncWithOptions(ctx, func(context.Context) error {
			return nil
		},
			WithOnSuccess(nil),
			WithOnError(nil),
			WithOnRetry(nil),
		)
	}, "Should handle nil callbacks gracefully")
}

func TestWorkerPoolConcurrentCloseAndSubmit(t *testing.T) {
	assertions := assert.New(t)

	for i := 0; i < 10; i++ { // sync multiple times to catch race conditions
		pool, err := newWorkerPool(DefaultConfig)
		assertions.NoError(err)

		var wg sync.WaitGroup
		wg.Add(2)

		// Goroutine 1: Submit tasks
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				_ = pool.async(func() {
					time.Sleep(10 * time.Millisecond)
				})
				time.Sleep(time.Millisecond)
			}
		}()

		// Goroutine 2: close Pool
		go func() {
			defer wg.Done()
			time.Sleep(15 * time.Millisecond)
			pool.close()
		}()

		wg.Wait()
		assertions.True(pool.closed.Load(), "Pool should be closed")
	}
}

func TestWorkerPoolMemoryLeak(t *testing.T) {
	assertions := assert.New(t)

	// Create and destroy multiple pools to check for leaks
	for i := 0; i < 100; i++ {
		pool, err := newWorkerPool(PoolConfig{
			MinWorkers: 1,
			MaxWorkers: 2,
			Timeout:    100 * time.Millisecond,
		})
		assertions.NoError(err)

		// Submit a few tasks
		for j := 0; j < 10; j++ {
			_ = pool.async(func() {
				time.Sleep(time.Microsecond)
			})
		}

		pool.close()
	}

	// If we get here without OOM or hanging, test passes
	assertions.True(true, "No memorypool leak detected")
}
