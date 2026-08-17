package circuitbreaker

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

// BenchmarkExecuteSuccess measures the performance of successful executions
func BenchmarkExecuteSuccess(b *testing.B) {
	cb := NewCircuitBreaker(DefaultConfig("benchmark"))
	fn := func() (any, error) {
		return "success", nil
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cb.Execute(fn)
	}
}

// BenchmarkExecuteFailure measures the performance of failed executions
// Uses high threshold to prevent circuit from opening
func BenchmarkExecuteFailure(b *testing.B) {
	cfg := Config{
		Name:             "benchmark",
		FailureThreshold: uint32(b.N) + 1000, // Prevent circuit from opening
	}
	cb := NewCircuitBreaker(cfg)
	testErr := errors.New("benchmark error")
	fn := func() (any, error) {
		return nil, testErr
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cb.Execute(fn)
	}
}

// BenchmarkExecuteCircuitOpen measures the fast-fail performance when circuit is open
func BenchmarkExecuteCircuitOpen(b *testing.B) {
	cfg := Config{
		Name:             "benchmark",
		FailureThreshold: 1,
		Timeout:          1 * time.Hour, // Keep circuit open for entire benchmark
	}
	cb := NewCircuitBreaker(cfg)

	// Trip the circuit
	_, _ = cb.Execute(func() (any, error) {
		return nil, errors.New("trip")
	})

	fn := func() (any, error) {
		return "should not execute", nil
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cb.Execute(fn)
	}
}

// BenchmarkExecuteParallel measures concurrent execution performance
func BenchmarkExecuteParallel(b *testing.B) {
	cfg := Config{
		Name:             "benchmark-parallel",
		FailureThreshold: 1000000, // High threshold to prevent tripping
	}
	cb := NewCircuitBreaker(cfg)
	fn := func() (any, error) {
		return "success", nil
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = cb.Execute(fn)
		}
	})
}

// BenchmarkExecuteParallelMixed measures concurrent mixed success/failure performance
func BenchmarkExecuteParallelMixed(b *testing.B) {
	cfg := Config{
		Name:             "benchmark-parallel-mixed",
		FailureThreshold: 1000000, // High threshold to prevent tripping
	}
	cb := NewCircuitBreaker(cfg)
	testErr := errors.New("benchmark error")
	var counter atomic.Int64

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c := counter.Add(1)
			if c%2 == 0 {
				_, _ = cb.Execute(func() (any, error) {
					return "success", nil
				})
			} else {
				_, _ = cb.Execute(func() (any, error) {
					return nil, testErr
				})
			}
		}
	})
}

// BenchmarkExecuteIsSuccessful measures overhead of custom success evaluation
func BenchmarkExecuteIsSuccessful(b *testing.B) {
	b.Run("without_IsSuccessful", func(b *testing.B) {
		cfg := Config{
			Name:             "benchmark-no-issuccessful",
			FailureThreshold: uint32(b.N) + 1000,
		}
		cb := NewCircuitBreaker(cfg)
		testErr := errors.New("expected error")
		fn := func() (any, error) {
			return nil, testErr
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = cb.Execute(fn)
		}
	})

	b.Run("with_IsSuccessful", func(b *testing.B) {
		expectedErr := errors.New("expected error")
		cfg := Config{
			Name:             "benchmark-with-issuccessful",
			FailureThreshold: uint32(b.N) + 1000,
			IsSuccessful: func(err error) bool {
				return errors.Is(err, expectedErr)
			},
		}
		cb := NewCircuitBreaker(cfg)
		fn := func() (any, error) {
			return nil, expectedErr
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = cb.Execute(fn)
		}
	})
}
