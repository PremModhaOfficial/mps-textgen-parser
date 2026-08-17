package circuitbreaker

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"

	"github.com/stretchr/testify/assert"
)

func TestNewCircuitBreaker(t *testing.T) {
	t.Run("creates with default config", func(t *testing.T) {
		assertions := assert.New(t)
		cfg := DefaultConfig("test")
		cb := NewCircuitBreaker(cfg)

		assertions.Equal("test", cb.Name())
		assertions.Equal(StateClosed, cb.State())
	})

	t.Run("creates with custom config", func(t *testing.T) {
		assertions := assert.New(t)
		cfg := Config{
			Name:             "custom",
			MaxRequests:      3,
			Timeout:          30 * time.Second,
			FailureThreshold: 3,
		}
		cb := NewCircuitBreaker(cfg)

		assertions.Equal("custom", cb.Name())
	})

	t.Run("applies defaults for zero values", func(t *testing.T) {
		assertions := assert.New(t)
		cfg := Config{Name: "zero-test"}
		cb := NewCircuitBreaker(cfg)

		assertions.NotNil(cb, "expected circuit breaker to be created")
	})
}

func TestCircuitBreakerExecute(t *testing.T) {
	t.Run("executes successful function", func(t *testing.T) {
		assertions := assert.New(t)
		cb := NewCircuitBreaker(DefaultConfig("test"))

		result, err := cb.Execute(func() (any, error) {
			return "success", nil
		})

		assertions.NoError(err)
		assertions.Equal("success", result)
	})

	t.Run("returns error from function", func(t *testing.T) {
		assertions := assert.New(t)
		cb := NewCircuitBreaker(DefaultConfig("test"))
		expectedErr := errors.New("test error")

		_, err := cb.Execute(func() (any, error) {
			return nil, expectedErr
		})

		assertions.ErrorIs(err, expectedErr)
	})

	t.Run("opens circuit after consecutive failures", func(t *testing.T) {
		assertions := assert.New(t)
		cfg := Config{
			Name:             "test",
			FailureThreshold: 3,
			Timeout:          1 * time.Second,
		}
		cb := NewCircuitBreaker(cfg)

		// Cause failures to trip the circuit
		for i := 0; i < 3; i++ {
			_, _ = cb.Execute(func() (any, error) {
				return nil, errors.New("failure")
			})
		}

		assertions.Equal(StateOpen, cb.State())

		// Next call should fail immediately with ErrCircuitOpen
		_, err := cb.Execute(func() (any, error) {
			return "should not execute", nil
		})

		assertions.ErrorIs(err, utils.ErrCircuitOpen)
	})

	t.Run("circuit recovers after timeout", func(t *testing.T) {
		assertions := assert.New(t)
		cfg := Config{
			Name:             "test",
			FailureThreshold: 2,
			Timeout:          100 * time.Millisecond,
		}
		cb := NewCircuitBreaker(cfg)

		// Open the circuit
		for i := 0; i < 2; i++ {
			_, _ = cb.Execute(func() (any, error) {
				return nil, errors.New("failure")
			})
		}

		assertions.Equal(StateOpen, cb.State())

		// Wait for timeout
		time.Sleep(150 * time.Millisecond)

		// Circuit should be half-open now, try a successful request
		result, err := cb.Execute(func() (any, error) {
			return "recovered", nil
		})

		assertions.NoError(err)
		assertions.Equal("recovered", result)

		// Circuit should be closed again
		assertions.Equal(StateClosed, cb.State())
	})
}

func TestCircuitBreakerState(t *testing.T) {
	t.Run("initial state is closed", func(t *testing.T) {
		assertions := assert.New(t)
		cb := NewCircuitBreaker(DefaultConfig("test"))

		assertions.Equal(StateClosed, cb.State())
		assertions.True(cb.IsClosed())
		assertions.False(cb.IsOpen())
	})

	t.Run("IsOpen returns true when circuit is open", func(t *testing.T) {
		assertions := assert.New(t)
		cfg := Config{
			Name:             "test",
			FailureThreshold: 1,
			Timeout:          1 * time.Second,
		}
		cb := NewCircuitBreaker(cfg)

		// Trip the circuit
		_, _ = cb.Execute(func() (any, error) {
			return nil, errors.New("failure")
		})

		assertions.True(cb.IsOpen())
		assertions.False(cb.IsClosed())
	})
}

func TestCircuitBreakerStateString(t *testing.T) {
	tests := []struct {
		state    State
		expected string
	}{
		{StateClosed, "closed"},
		{StateHalfOpen, "half-open"},
		{StateOpen, "open"},
		{State(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tt.expected, tt.state.String())
		})
	}
}

func TestCircuitBreakerCounts(t *testing.T) {
	t.Run("tracks successful requests", func(t *testing.T) {
		assertions := assert.New(t)
		cb := NewCircuitBreaker(DefaultConfig("test"))

		for i := 0; i < 3; i++ {
			_, _ = cb.Execute(func() (any, error) {
				return nil, nil
			})
		}

		counts := cb.Counts()
		assertions.Equal(uint32(3), counts.Requests)
		assertions.Equal(uint32(3), counts.TotalSuccesses)
		assertions.Equal(uint32(0), counts.TotalFailures)
		assertions.Equal(uint32(3), counts.ConsecutiveSuccesses)
	})

	t.Run("tracks failed requests", func(t *testing.T) {
		assertions := assert.New(t)
		cfg := Config{
			Name:             "test",
			FailureThreshold: 10, // High threshold so circuit stays closed
		}
		cb := NewCircuitBreaker(cfg)

		for i := 0; i < 3; i++ {
			_, _ = cb.Execute(func() (any, error) {
				return nil, errors.New("failure")
			})
		}

		counts := cb.Counts()
		assertions.Equal(uint32(3), counts.TotalFailures)
		assertions.Equal(uint32(3), counts.ConsecutiveFailures)
	})
}

func TestCircuitBreakerOnStateChange(t *testing.T) {
	stateChanges := make([]struct {
		from State
		to   State
	}, 0)

	cfg := Config{
		Name:             "test",
		FailureThreshold: 2,
		Timeout:          100 * time.Millisecond,
		OnStateChange: func(name string, from, to State) {
			stateChanges = append(stateChanges, struct {
				from State
				to   State
			}{from, to})
		},
	}
	cb := NewCircuitBreaker(cfg)

	// Cause failures to open the circuit
	for i := 0; i < 2; i++ {
		_, _ = cb.Execute(func() (any, error) {
			return nil, errors.New("failure")
		})
	}

	assertions := assert.New(t)
	assertions.True(len(stateChanges) > 0, "expected state change callback to be called")

	// First change should be closed -> open
	assertions.Equal(StateClosed, stateChanges[0].from)
	assertions.Equal(StateOpen, stateChanges[0].to)
}

func TestCircuitBreakerIsSuccessful(t *testing.T) {
	// Custom error that should not count as failure
	ignoredErr := errors.New("ignored error")

	cfg := Config{
		Name:             "test",
		FailureThreshold: 2,
		IsSuccessful: func(err error) bool {
			return errors.Is(err, ignoredErr)
		},
	}
	cb := NewCircuitBreaker(cfg)

	// These should not count as failures
	for i := 0; i < 5; i++ {
		_, _ = cb.Execute(func() (any, error) {
			return nil, ignoredErr
		})
	}

	// Circuit should still be closed because errors were ignored
	assertions := assert.New(t)
	assertions.Equal(StateClosed, cb.State(), "expected state closed (ignored errors)")

	counts := cb.Counts()
	assertions.Equal(uint32(0), counts.TotalFailures, "expected 0 failures (ignored)")
}

func TestCircuitBreakerDefaultConfig(t *testing.T) {
	assertions := assert.New(t)
	cfg := DefaultConfig("my-breaker")

	assertions.Equal("my-breaker", cfg.Name)
	assertions.Equal(uint32(1), cfg.MaxRequests)
	assertions.Equal(60*time.Second, cfg.Timeout)
	assertions.Equal(uint32(5), cfg.FailureThreshold)
	assertions.Equal(uint32(1), cfg.SuccessThreshold)
	assertions.Equal(time.Duration(0), cfg.Interval)
}

func TestCircuitBreakerErrors(t *testing.T) {
	t.Run("ErrCircuitOpen is returned when circuit is open", func(t *testing.T) {
		cfg := Config{
			Name:             "test",
			FailureThreshold: 1,
			Timeout:          1 * time.Second,
		}
		cb := NewCircuitBreaker(cfg)

		// Trip the circuit
		_, _ = cb.Execute(func() (any, error) {
			return nil, errors.New("failure")
		})

		_, err := cb.Execute(func() (any, error) {
			return nil, nil
		})

		assertions := assert.New(t)
		assertions.ErrorIs(err, utils.ErrCircuitOpen)
	})
}

func TestCircuitBreakerFullStateTransitionCycle(t *testing.T) {
	t.Run("closed -> open -> half-open -> closed cycle", func(t *testing.T) {
		stateTransitions := make([]struct {
			from State
			to   State
		}, 0)

		cfg := Config{
			Name:             "lifecycle-test",
			FailureThreshold: 3,
			MaxRequests:      1, // Only 1 successful request needed in half-open to close
			Timeout:          100 * time.Millisecond,
			OnStateChange: func(name string, from, to State) {
				stateTransitions = append(stateTransitions, struct {
					from State
					to   State
				}{from, to})
			},
		}
		cb := NewCircuitBreaker(cfg)

		// Step 1: Verify initial state is CLOSED
		assertions := assert.New(t)
		assertions.Equal(StateClosed, cb.State(), "Step 1: expected initial state closed")

		// Step 2: Cause failures to transition to OPEN
		for i := 0; i < 3; i++ {
			_, err := cb.Execute(func() (any, error) {
				return nil, errors.New("failure")
			})
			if i < 2 && err == nil {
				assertions.Error(err, "Step 2: expected error on failure %d", i)
			}
		}

		assertions.Equal(StateOpen, cb.State(), "Step 2: expected state open after 3 failures")

		// Verify first state transition: Closed -> Open
		assertions.True(len(stateTransitions) >= 1, "Step 2: expected at least one state transition")
		assertions.Equal(StateClosed, stateTransitions[0].from, "Step 2: expected transition from closed")
		assertions.Equal(StateOpen, stateTransitions[0].to, "Step 2: expected transition to open")

		// Step 3: Verify requests are rejected while OPEN
		_, err := cb.Execute(func() (any, error) {
			return "should not execute", nil
		})
		assertions.ErrorIs(err, utils.ErrCircuitOpen, "Step 3: expected ErrCircuitOpen while open")

		// Step 4: Wait for timeout to transition to HALF-OPEN
		time.Sleep(150 * time.Millisecond)

		// The state transitions to half-open on the next request attempt
		// Make a successful request to trigger the transition and close the circuit
		result, err := cb.Execute(func() (any, error) {
			return "recovery success", nil
		})
		assertions.NoError(err, "Step 4: unexpected error in half-open state")
		assertions.Equal("recovery success", result, "Step 4: expected 'recovery success'")

		// Verify state transitions: Open -> Half-Open -> Closed
		assertions.True(len(stateTransitions) >= 3, "Step 4: expected at least 3 state transitions")
		assertions.Equal(StateOpen, stateTransitions[1].from, "Step 4: expected transition from open")
		assertions.Equal(StateHalfOpen, stateTransitions[1].to, "Step 4: expected transition to half-open")
		assertions.Equal(StateHalfOpen, stateTransitions[2].from, "Step 4: expected transition from half-open")
		assertions.Equal(StateClosed, stateTransitions[2].to, "Step 4: expected transition to closed")

		// Verify final state is CLOSED
		assertions.Equal(StateClosed, cb.State(), "Step 4: expected state closed after recovery")

		// Step 5: Verify circuit is fully operational again
		for i := 0; i < 5; i++ {
			result, err = cb.Execute(func() (any, error) {
				return "operational", nil
			})
			assertions.NoError(err, "Step 5: unexpected error on request %d", i)
		}

		counts := cb.Counts()
		assertions.GreaterOrEqual(counts.TotalSuccesses, uint32(5), "Step 5: expected at least 5 successes")
	})

	t.Run("closed -> open -> half-open -> open (failure in half-open)", func(t *testing.T) {
		stateTransitions := make([]struct {
			from State
			to   State
		}, 0)

		cfg := Config{
			Name:             "half-open-failure-test",
			FailureThreshold: 2,
			SuccessThreshold: 2,
			MaxRequests:      1,
			Timeout:          100 * time.Millisecond,
			OnStateChange: func(name string, from, to State) {
				stateTransitions = append(stateTransitions, struct {
					from State
					to   State
				}{from, to})
			},
		}
		cb := NewCircuitBreaker(cfg)

		// Step 1: Trip the circuit to OPEN
		assertions := assert.New(t)
		for i := 0; i < 2; i++ {
			_, _ = cb.Execute(func() (any, error) {
				return nil, errors.New("failure")
			})
		}

		assertions.Equal(StateOpen, cb.State(), "Step 1: expected state open")

		// Step 2: Wait for timeout
		time.Sleep(150 * time.Millisecond)

		// Step 3: Fail in half-open state - should go back to OPEN
		_, err := cb.Execute(func() (any, error) {
			return nil, errors.New("failure in half-open")
		})
		assertions.Error(err, "Step 3: expected error from failed execution")

		// Circuit should be back to OPEN
		assertions.Equal(StateOpen, cb.State(), "Step 3: expected state open after half-open failure")

		// Verify transitions: Closed -> Open -> Half-Open -> Open
		assertions.True(len(stateTransitions) >= 3, "Step 3: expected at least 3 state transitions")

		expectedTransitions := []struct {
			from State
			to   State
		}{
			{StateClosed, StateOpen},
			{StateOpen, StateHalfOpen},
			{StateHalfOpen, StateOpen},
		}

		for i, expected := range expectedTransitions {
			if i >= len(stateTransitions) {
				assertions.Failf("missing transition", "Step 3: missing transition %d", i)
				continue
			}
			assertions.Equal(expected.from, stateTransitions[i].from, "Step 3: transition %d from state", i)
			assertions.Equal(expected.to, stateTransitions[i].to, "Step 3: transition %d to state", i)
		}
	})

	t.Run("multiple cycles of open and recovery", func(t *testing.T) {
		assertions := assert.New(t)
		cfg := Config{
			Name:             "multi-cycle-test",
			FailureThreshold: 2,
			SuccessThreshold: 1,
			MaxRequests:      1,
			Timeout:          50 * time.Millisecond,
		}
		cb := NewCircuitBreaker(cfg)

		for cycle := 1; cycle <= 3; cycle++ {
			// Verify starting in closed state
			assertions.Equal(StateClosed, cb.State(), "Cycle %d: expected initial state closed", cycle)

			// Trip the circuit
			for i := 0; i < 2; i++ {
				_, _ = cb.Execute(func() (any, error) {
					return nil, errors.New("failure")
				})
			}

			assertions.Equal(StateOpen, cb.State(), "Cycle %d: expected state open", cycle)

			// Wait for half-open
			time.Sleep(60 * time.Millisecond)

			// Recover with success
			_, err := cb.Execute(func() (any, error) {
				return "recovered", nil
			})
			assertions.NoError(err, "Cycle %d: unexpected error during recovery", cycle)

			// Should be closed again
			assertions.Equal(StateClosed, cb.State(), "Cycle %d: expected state closed after recovery", cycle)
		}
	})
}

func TestCircuitBreakerIntervalResetsCounts(t *testing.T) {
	t.Run("resets failure count after interval expires", func(t *testing.T) {
		cfg := Config{
			Name:             "interval-test",
			FailureThreshold: 5,
			Interval:         100 * time.Millisecond, // Counts reset every 100ms
			Timeout:          1 * time.Second,
		}
		cb := NewCircuitBreaker(cfg)

		// Cause some failures (but not enough to trip the circuit)
		for i := 0; i < 3; i++ {
			_, _ = cb.Execute(func() (any, error) {
				return nil, errors.New("failure")
			})
		}

		counts := cb.Counts()
		assertions := assert.New(t)
		assertions.Equal(uint32(3), counts.ConsecutiveFailures)
		assertions.Equal(uint32(3), counts.TotalFailures)

		// Circuit should still be closed
		assertions.Equal(StateClosed, cb.State())

		// Wait for the interval to pass
		time.Sleep(150 * time.Millisecond)

		// Make a successful request to trigger the count reset check
		_, _ = cb.Execute(func() (any, error) {
			return "success", nil
		})

		// Counts should be reset (only the new successful request should be counted)
		counts = cb.Counts()
		assertions.Equal(uint32(0), counts.ConsecutiveFailures, "expected 0 consecutive failures after interval reset")
		assertions.Equal(uint32(0), counts.TotalFailures, "expected 0 total failures after interval reset")
		assertions.Equal(uint32(1), counts.TotalSuccesses, "expected 1 total success after interval reset")
	})

	t.Run("does not trip circuit if failures spread across intervals", func(t *testing.T) {
		cfg := Config{
			Name:             "interval-spread-test",
			FailureThreshold: 3,
			Interval:         100 * time.Millisecond,
			Timeout:          1 * time.Second,
		}
		cb := NewCircuitBreaker(cfg)

		// Cause 2 failures in first interval
		for i := 0; i < 2; i++ {
			_, _ = cb.Execute(func() (any, error) {
				return nil, errors.New("failure")
			})
		}

		assertions := assert.New(t)
		assertions.Equal(StateClosed, cb.State(), "expected state closed after 2 failures")

		// Wait for interval to pass and counts to reset
		time.Sleep(150 * time.Millisecond)

		// Cause 2 more failures in second interval (would be 4 total without reset)
		for i := 0; i < 2; i++ {
			_, _ = cb.Execute(func() (any, error) {
				return nil, errors.New("failure")
			})
		}

		// Circuit should still be closed because counts were reset
		assertions.Equal(StateClosed, cb.State(), "expected state closed (counts should have reset)")

		counts := cb.Counts()
		assertions.Equal(uint32(2), counts.ConsecutiveFailures, "expected 2 consecutive failures in new interval")
	})

	t.Run("trips circuit when failures exceed threshold within interval", func(t *testing.T) {
		cfg := Config{
			Name:             "interval-trip-test",
			FailureThreshold: 3,
			Interval:         500 * time.Millisecond, // Long interval
			Timeout:          1 * time.Second,
		}
		cb := NewCircuitBreaker(cfg)

		// Cause failures quickly within the same interval
		for i := 0; i < 3; i++ {
			_, _ = cb.Execute(func() (any, error) {
				return nil, errors.New("failure")
			})
		}

		// Circuit should be open because all failures happened within the interval
		assertions := assert.New(t)
		assertions.Equal(StateOpen, cb.State(), "expected state open")
	})

	t.Run("no interval means counts never reset", func(t *testing.T) {
		cfg := Config{
			Name:             "no-interval-test",
			FailureThreshold: 10,
			Interval:         0, // No interval - counts never reset
			Timeout:          1 * time.Second,
		}
		cb := NewCircuitBreaker(cfg)

		// Cause some failures
		for i := 0; i < 3; i++ {
			_, _ = cb.Execute(func() (any, error) {
				return nil, errors.New("failure")
			})
		}

		// Wait some time
		time.Sleep(100 * time.Millisecond)

		// Make another request
		_, _ = cb.Execute(func() (any, error) {
			return nil, errors.New("failure")
		})

		// Counts should NOT be reset
		assertions := assert.New(t)
		counts := cb.Counts()
		assertions.Equal(uint32(4), counts.TotalFailures, "expected 4 total failures (no reset)")
		assertions.Equal(uint32(4), counts.ConsecutiveFailures, "expected 4 consecutive failures (no reset)")
	})
}

func TestCircuitBreakerConcurrentExecution(t *testing.T) {
	t.Run("handles concurrent successful requests", func(t *testing.T) {
		// Use high failure threshold to ensure circuit stays closed
		cfg := Config{
			Name:             "concurrent-success-test",
			FailureThreshold: 1000,
			Timeout:          1 * time.Second,
		}
		cb := NewCircuitBreaker(cfg)

		var wg sync.WaitGroup
		var successCount atomic.Int64
		var errorCount atomic.Int64
		const numGoroutines = 100

		for i := 0; i < numGoroutines; i++ {

			wg.Go(func() {

				_, err := cb.Execute(func() (any, error) {
					return "success", nil
				})
				if err == nil {
					successCount.Add(1)
				} else {
					errorCount.Add(1)
				}
			})
		}

		wg.Wait()

		// Assert precise counts
		assertions := assert.New(t)
		assertions.Equal(int64(numGoroutines), successCount.Load(), "successCount should match numGoroutines")
		assertions.Equal(int64(0), errorCount.Load(), "errorCount should be 0")

		counts := cb.Counts()
		assertions.Equal(uint32(numGoroutines), counts.Requests, "counts.Requests should match numGoroutines")
		assertions.Equal(uint32(numGoroutines), counts.TotalSuccesses, "counts.TotalSuccesses should match numGoroutines")
		assertions.Equal(uint32(0), counts.TotalFailures, "counts.TotalFailures should be 0")
		assertions.Equal(uint32(numGoroutines), counts.ConsecutiveSuccesses, "counts.ConsecutiveSuccesses should match numGoroutines")
		assertions.Equal(uint32(0), counts.ConsecutiveFailures, "counts.ConsecutiveFailures should be 0")
	})

	t.Run("handles concurrent failures with circuit staying closed", func(t *testing.T) {
		// Use high failure threshold to keep circuit closed
		cfg := Config{
			Name:             "concurrent-failure-closed-test",
			FailureThreshold: 1000,
			Timeout:          1 * time.Second,
		}
		cb := NewCircuitBreaker(cfg)

		var wg sync.WaitGroup
		var failureCount atomic.Int64
		var successCount atomic.Int64
		var circuitOpenCount atomic.Int64
		const numGoroutines = 100

		for i := 0; i < numGoroutines; i++ {

			wg.Go(func() {

				_, err := cb.Execute(func() (any, error) {
					return nil, errors.New("failure")
				})
				if err == nil {
					successCount.Add(1)
				} else if errors.Is(err, utils.ErrCircuitOpen) {
					circuitOpenCount.Add(1)
				} else {
					failureCount.Add(1)
				}
			})
		}

		wg.Wait()

		// Circuit should still be closed (threshold not reached)
		assertions := assert.New(t)
		assertions.Equal(StateClosed, cb.State(), "expected state closed")

		// Assert precise counts
		assertions.Equal(int64(numGoroutines), failureCount.Load(), "failureCount should match numGoroutines")
		assertions.Equal(int64(0), successCount.Load(), "successCount should be 0")
		assertions.Equal(int64(0), circuitOpenCount.Load(), "circuitOpenCount should be 0")

		counts := cb.Counts()
		assertions.Equal(uint32(numGoroutines), counts.Requests, "counts.Requests should match numGoroutines")
		assertions.Equal(uint32(numGoroutines), counts.TotalFailures, "counts.TotalFailures should match numGoroutines")
		assertions.Equal(uint32(0), counts.TotalSuccesses, "counts.TotalSuccesses should be 0")
	})

	t.Run("sequential trip then concurrent rejection", func(t *testing.T) {
		cfg := Config{
			Name:             "trip-then-reject-test",
			FailureThreshold: 3,
			Timeout:          1 * time.Second,
		}
		cb := NewCircuitBreaker(cfg)

		// Phase 1: Trip the circuit sequentially (predictable)
		for i := 0; i < 3; i++ {
			_, _ = cb.Execute(func() (any, error) {
				return nil, errors.New("failure")
			})
		}

		assertions := assert.New(t)
		assertions.Equal(StateOpen, cb.State(), "expected state open after 3 failures")

		// Phase 2: All concurrent requests should be rejected
		var wg sync.WaitGroup
		var circuitOpenCount atomic.Int64
		var otherErrorCount atomic.Int64
		var successCount atomic.Int64
		const numGoroutines = 50

		for i := 0; i < numGoroutines; i++ {

			wg.Go(func() {
				_, err := cb.Execute(func() (any, error) {
					return "should not execute", nil
				})
				if err == nil {
					successCount.Add(1)
				} else if errors.Is(err, utils.ErrCircuitOpen) {
					circuitOpenCount.Add(1)
				} else {
					otherErrorCount.Add(1)
				}
			})
		}

		wg.Wait()

		// Assert precise counts - ALL requests should be rejected
		assertions.Equal(int64(numGoroutines), circuitOpenCount.Load(), "circuitOpenCount should match numGoroutines")
		assertions.Equal(int64(0), successCount.Load(), "successCount should be 0")
		assertions.Equal(int64(0), otherErrorCount.Load(), "otherErrorCount should be 0")

		// Note: Counts are reset when circuit transitions to open state
		// Rejected requests (ErrCircuitOpen) are not counted in the circuit breaker stats
		counts := cb.Counts()
		assertions.Equal(uint32(0), counts.Requests, "counts.Requests should be 0 (reset on open)")
	})

	t.Run("handles mixed concurrent success and failure", func(t *testing.T) {
		// Use very high failure threshold to ensure circuit stays closed
		cfg := Config{
			Name:             "concurrent-mixed-test",
			FailureThreshold: 1000,
			Timeout:          1 * time.Second,
		}
		cb := NewCircuitBreaker(cfg)

		var wg sync.WaitGroup
		var successCount atomic.Int64
		var failureCount atomic.Int64
		var circuitOpenCount atomic.Int64
		const numGoroutines = 100
		// 50 successes (even indices), 50 failures (odd indices)

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				_, err := cb.Execute(func() (any, error) {
					if idx%2 == 0 {
						return "success", nil
					}
					return nil, errors.New("failure")
				})
				if err == nil {
					successCount.Add(1)
				} else if errors.Is(err, utils.ErrCircuitOpen) {
					circuitOpenCount.Add(1)
				} else {
					failureCount.Add(1)
				}
			}(i)
		}

		wg.Wait()

		// Circuit should still be closed
		assertions := assert.New(t)
		assertions.Equal(StateClosed, cb.State(), "expected state closed")

		// Assert precise counts
		assertions.Equal(int64(50), successCount.Load(), "successCount should be 50")
		assertions.Equal(int64(50), failureCount.Load(), "failureCount should be 50")
		assertions.Equal(int64(0), circuitOpenCount.Load(), "circuitOpenCount should be 0")

		counts := cb.Counts()
		assertions.Equal(uint32(numGoroutines), counts.Requests, "counts.Requests should match numGoroutines")
		assertions.Equal(uint32(50), counts.TotalSuccesses, "counts.TotalSuccesses should be 50")
		assertions.Equal(uint32(50), counts.TotalFailures, "counts.TotalFailures should be 50")
	})

	t.Run("concurrent state reads are safe", func(t *testing.T) {
		cfg := Config{
			Name:             "concurrent-state-test",
			FailureThreshold: 1000,
			Timeout:          100 * time.Millisecond,
		}
		cb := NewCircuitBreaker(cfg)

		var wg sync.WaitGroup
		const numReaders = 50
		const numWriters = 20
		const writesPerGoroutine = 10

		// Start readers that continuously check state
		for i := 0; i < numReaders; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					_ = cb.State()
					_ = cb.IsOpen()
					_ = cb.IsClosed()
					_ = cb.Counts()
				}
			}()
		}

		// Start writers - all successful requests
		for i := 0; i < numWriters; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < writesPerGoroutine; j++ {
					_, _ = cb.Execute(func() (any, error) {
						return "success", nil
					})
				}
			}()
		}

		wg.Wait()

		// Assert precise counts
		assertions := assert.New(t)
		counts := cb.Counts()
		expectedRequests := uint32(numWriters * writesPerGoroutine)
		assertions.Equal(expectedRequests, counts.Requests, "counts.Requests should match expected")
		assertions.Equal(expectedRequests, counts.TotalSuccesses, "counts.TotalSuccesses should match expected")
		assertions.Equal(uint32(0), counts.TotalFailures, "counts.TotalFailures should be 0")
	})

	t.Run("sequential half-open then concurrent requests", func(t *testing.T) {
		cfg := Config{
			Name:             "halfopen-concurrent-test",
			FailureThreshold: 2,
			MaxRequests:      1, // Only 1 request allowed in half-open state
			Timeout:          50 * time.Millisecond,
		}
		cb := NewCircuitBreaker(cfg)

		// Phase 1: Trip the circuit sequentially
		for i := 0; i < 2; i++ {
			_, _ = cb.Execute(func() (any, error) {
				return nil, errors.New("failure")
			})
		}

		assertions := assert.New(t)
		assertions.Equal(StateOpen, cb.State(), "expected state open")

		// Phase 2: Wait for circuit to transition to half-open
		time.Sleep(60 * time.Millisecond)

		// Phase 3: Single successful request to close the circuit
		_, err := cb.Execute(func() (any, error) {
			return "success", nil
		})
		assertions.NoError(err, "expected successful execution in half-open")

		assertions.Equal(StateClosed, cb.State(), "expected state closed after recovery")

		// Phase 4: Now all concurrent requests should succeed
		var wg sync.WaitGroup
		var successCount atomic.Int64
		var errorCount atomic.Int64
		const numGoroutines = 50

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := cb.Execute(func() (any, error) {
					return "success", nil
				})
				if err == nil {
					successCount.Add(1)
				} else {
					errorCount.Add(1)
				}
			}()
		}

		wg.Wait()

		// Assert precise counts
		assertions.Equal(int64(numGoroutines), successCount.Load(), "successCount should match numGoroutines")
		assertions.Equal(int64(0), errorCount.Load(), "errorCount should be 0")
	})

	t.Run("concurrent state transition cycle with precise counts", func(t *testing.T) {
		var stateTransitions []struct {
			from State
			to   State
		}
		var transitionMu sync.Mutex

		cfg := Config{
			Name:             "concurrent-lifecycle-test",
			FailureThreshold: 5,
			MaxRequests:      1,
			Timeout:          100 * time.Millisecond,
			OnStateChange: func(name string, from, to State) {
				transitionMu.Lock()
				stateTransitions = append(stateTransitions, struct {
					from State
					to   State
				}{from, to})
				transitionMu.Unlock()
			},
		}
		cb := NewCircuitBreaker(cfg)

		// Phase 1: Verify initial state is CLOSED
		assertions := assert.New(t)
		assertions.Equal(StateClosed, cb.State(), "Phase 1: expected initial state closed")

		// Phase 2: Trip the circuit sequentially (predictable 5 failures)
		for i := 0; i < 5; i++ {
			_, _ = cb.Execute(func() (any, error) {
				return nil, errors.New("failure")
			})
		}

		assertions.Equal(StateOpen, cb.State(), "Phase 2: expected state open after 5 failures")

		// Phase 3: Concurrent requests should ALL be rejected while OPEN
		var wg sync.WaitGroup
		var rejectedCount atomic.Int64
		var otherCount atomic.Int64
		const numRejectedGoroutines = 20

		for i := 0; i < numRejectedGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := cb.Execute(func() (any, error) {
					return "should not execute", nil
				})
				if errors.Is(err, utils.ErrCircuitOpen) {
					rejectedCount.Add(1)
				} else {
					otherCount.Add(1)
				}
			}()
		}
		wg.Wait()

		// Assert precise counts - ALL should be rejected
		assertions.Equal(int64(numRejectedGoroutines), rejectedCount.Load(), "Phase 3 rejectedCount should match expected")
		assertions.Equal(int64(0), otherCount.Load(), "Phase 3 otherCount should be 0")

		// Phase 4: Wait for timeout to transition to HALF-OPEN
		time.Sleep(150 * time.Millisecond)

		// Phase 5: Single request to recover (MaxRequests=1)
		_, err := cb.Execute(func() (any, error) {
			return "recovery", nil
		})
		assertions.NoError(err, "Phase 5: expected successful recovery")

		assertions.Equal(StateClosed, cb.State(), "Phase 5: expected state closed after recovery")

		// Phase 6: Verify circuit is fully operational with concurrent successful requests
		var finalSuccessCount atomic.Int64
		var finalErrorCount atomic.Int64
		const numFinalGoroutines = 50

		for i := 0; i < numFinalGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := cb.Execute(func() (any, error) {
					return "operational", nil
				})
				if err == nil {
					finalSuccessCount.Add(1)
				} else {
					finalErrorCount.Add(1)
				}
			}()
		}
		wg.Wait()

		// Assert precise counts
		assertions.Equal(int64(numFinalGoroutines), finalSuccessCount.Load(), "Phase 6 finalSuccessCount should match expected")
		assertions.Equal(int64(0), finalErrorCount.Load(), "Phase 6 finalErrorCount should be 0")

		// Verify state transitions
		transitionMu.Lock()
		defer transitionMu.Unlock()

		assertions.Equal(3, len(stateTransitions), "Expected exactly 3 state transitions")
		if len(stateTransitions) >= 1 {
			assertions.Equal(StateClosed, stateTransitions[0].from, "Transition 1: expected from closed")
			assertions.Equal(StateOpen, stateTransitions[0].to, "Transition 1: expected to open")
		}
		if len(stateTransitions) >= 2 {
			assertions.Equal(StateOpen, stateTransitions[1].from, "Transition 2: expected from open")
			assertions.Equal(StateHalfOpen, stateTransitions[1].to, "Transition 2: expected to half-open")
		}
		if len(stateTransitions) >= 3 {
			assertions.Equal(StateHalfOpen, stateTransitions[2].from, "Transition 3: expected from half-open")
			assertions.Equal(StateClosed, stateTransitions[2].to, "Transition 3: expected to closed")
		}
	})

	t.Run("sequential recovery failure then concurrent success", func(t *testing.T) {
		cfg := Config{
			Name:             "recovery-failure-test",
			FailureThreshold: 3,
			MaxRequests:      1,
			Timeout:          100 * time.Millisecond,
		}
		cb := NewCircuitBreaker(cfg)

		// Phase 1: Trip the circuit sequentially
		for i := 0; i < 3; i++ {
			_, _ = cb.Execute(func() (any, error) {
				return nil, errors.New("failure")
			})
		}

		assertions := assert.New(t)
		assertions.Equal(StateOpen, cb.State(), "Phase 1: expected state open")

		// Note: Counts are reset when circuit transitions to open state
		counts := cb.Counts()
		assertions.Equal(uint32(0), counts.Requests, "Phase 1 counts.Requests: expected 0 (reset on open)")

		// Phase 2: Wait for half-open
		time.Sleep(150 * time.Millisecond)

		// Phase 3: Single failing request in half-open state
		_, err := cb.Execute(func() (any, error) {
			return nil, errors.New("still failing")
		})
		assertions.Error(err, "Phase 3: expected error from failed execution")

		assertions.Equal(StateOpen, cb.State(), "Phase 3: expected state open after half-open failure")

		// Phase 4: Wait again and recover sequentially
		time.Sleep(150 * time.Millisecond)

		_, err = cb.Execute(func() (any, error) {
			return "recovered", nil
		})
		assertions.NoError(err, "Phase 4: expected successful recovery")

		assertions.Equal(StateClosed, cb.State(), "Phase 4: expected state closed after recovery")

		// Phase 5: Now all concurrent requests should succeed
		var wg sync.WaitGroup
		var successCount atomic.Int64
		var errorCount atomic.Int64
		const numGoroutines = 50

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := cb.Execute(func() (any, error) {
					return "success", nil
				})
				if err == nil {
					successCount.Add(1)
				} else {
					errorCount.Add(1)
				}
			}()
		}
		wg.Wait()

		// Assert precise counts
		assertions.Equal(int64(numGoroutines), successCount.Load(), "Phase 5 successCount should match numGoroutines")
		assertions.Equal(int64(0), errorCount.Load(), "Phase 5 errorCount should be 0")
	})

	t.Run("multiple sequential recovery cycles with concurrent verification", func(t *testing.T) {
		cfg := Config{
			Name:             "multi-cycle-test",
			FailureThreshold: 3,
			MaxRequests:      1,
			Timeout:          50 * time.Millisecond,
		}
		cb := NewCircuitBreaker(cfg)

		assertions := assert.New(t)
		for cycle := 1; cycle <= 3; cycle++ {
			// Verify starting in closed state
			assertions.Equal(StateClosed, cb.State(), "Cycle %d: expected initial state closed", cycle)

			// Sequential failures to trip the circuit (predictable)
			for i := 0; i < 3; i++ {
				_, _ = cb.Execute(func() (any, error) {
					return nil, errors.New("failure")
				})
			}

			assertions.Equal(StateOpen, cb.State(), "Cycle %d: expected state open", cycle)

			// Wait for half-open
			time.Sleep(60 * time.Millisecond)

			// Single recovery request
			_, err := cb.Execute(func() (any, error) {
				return "recovered", nil
			})
			assertions.NoError(err, "Cycle %d: expected successful recovery", cycle)

			assertions.Equal(StateClosed, cb.State(), "Cycle %d: expected state closed after recovery", cycle)

			// Concurrent verification - all should succeed
			var wg sync.WaitGroup
			var successCount atomic.Int64
			var errorCount atomic.Int64
			const numVerifyGoroutines = 20

			for i := 0; i < numVerifyGoroutines; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					_, err := cb.Execute(func() (any, error) {
						return "verify", nil
					})
					if err == nil {
						successCount.Add(1)
					} else {
						errorCount.Add(1)
					}
				}()
			}
			wg.Wait()

			// Assert precise counts
			assertions.Equal(int64(numVerifyGoroutines), successCount.Load(), "Cycle %d verify successCount should match expected", cycle)
			assertions.Equal(int64(0), errorCount.Load(), "Cycle %d verify errorCount should be 0", cycle)
		}
	})

	t.Run("stress test with precise counts", func(t *testing.T) {
		cfg := Config{
			Name:             "stress-test",
			FailureThreshold: 100000, // Very high threshold to keep circuit closed
			Timeout:          100 * time.Millisecond,
		}
		cb := NewCircuitBreaker(cfg)

		var wg sync.WaitGroup
		const numGoroutines = 1000
		const requestsPerGoroutine = 10
		var successCount atomic.Int64
		var failureCount atomic.Int64

		// 200 goroutines will fail (idx % 5 == 0), 800 will succeed
		// Each makes 10 requests
		// Expected: 2000 failures, 8000 successes

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				for j := 0; j < requestsPerGoroutine; j++ {
					_, err := cb.Execute(func() (any, error) {
						if idx%5 == 0 {
							return nil, errors.New("failure")
						}
						return "success", nil
					})
					if err == nil {
						successCount.Add(1)
					} else {
						failureCount.Add(1)
					}
				}
			}(i)
		}

		wg.Wait()

		// Circuit should still be closed
		assertions := assert.New(t)
		assertions.Equal(StateClosed, cb.State(), "expected state closed")

		// Assert precise counts
		const expectedSuccesses = 8000 // 800 goroutines * 10 requests
		const expectedFailures = 2000  // 200 goroutines * 10 requests
		const expectedTotal = numGoroutines * requestsPerGoroutine

		assertions.Equal(int64(expectedSuccesses), successCount.Load(), "successCount should match expected")
		assertions.Equal(int64(expectedFailures), failureCount.Load(), "failureCount should match expected")

		counts := cb.Counts()
		assertions.Equal(uint32(expectedTotal), counts.Requests, "counts.Requests should match expected")
		assertions.Equal(uint32(expectedSuccesses), counts.TotalSuccesses, "counts.TotalSuccesses should match expected")
		assertions.Equal(uint32(expectedFailures), counts.TotalFailures, "counts.TotalFailures should match expected")
	})
}
