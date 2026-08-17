package circuitbreaker

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"
)

func TestNewCircuitBreaker(t *testing.T) {
	t.Run("creates with default config", func(t *testing.T) {
		cfg := DefaultConfig("test")
		cb := NewCircuitBreaker(cfg)

		if cb.Name() != "test" {
			t.Errorf("expected name 'test', got %s", cb.Name())
		}
		if cb.State() != StateClosed {
			t.Errorf("expected initial state closed, got %s", cb.State())
		}
	})

	t.Run("creates with custom config", func(t *testing.T) {
		cfg := Config{
			Name:             "custom",
			MaxRequests:      3,
			Timeout:          30 * time.Second,
			FailureThreshold: 3,
		}
		cb := NewCircuitBreaker(cfg)

		if cb.Name() != "custom" {
			t.Errorf("expected name 'custom', got %s", cb.Name())
		}
	})

	t.Run("applies defaults for zero values", func(t *testing.T) {
		cfg := Config{Name: "zero-test"}
		cb := NewCircuitBreaker(cfg)

		if cb == nil {
			t.Error("expected circuit breaker to be created")
		}
	})
}

func TestCircuitBreakerExecute(t *testing.T) {
	t.Run("executes successful function", func(t *testing.T) {
		cb := NewCircuitBreaker(DefaultConfig("test"))

		result, err := cb.Execute(func() (any, error) {
			return "success", nil
		})

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != "success" {
			t.Errorf("expected 'success', got %v", result)
		}
	})

	t.Run("returns error from function", func(t *testing.T) {
		cb := NewCircuitBreaker(DefaultConfig("test"))
		expectedErr := errors.New("test error")

		_, err := cb.Execute(func() (any, error) {
			return nil, expectedErr
		})

		if !errors.Is(err, expectedErr) {
			t.Errorf("expected %v, got %v", expectedErr, err)
		}
	})

	t.Run("opens circuit after consecutive failures", func(t *testing.T) {
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

		if cb.State() != StateOpen {
			t.Errorf("expected state open, got %s", cb.State())
		}

		// Next call should fail immediately with ErrCircuitOpen
		_, err := cb.Execute(func() (any, error) {
			return "should not execute", nil
		})

		if !errors.Is(err, utils.ErrCircuitOpen) {
			t.Errorf("expected ErrCircuitOpen, got %v", err)
		}
	})

	t.Run("circuit recovers after timeout", func(t *testing.T) {
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

		if cb.State() != StateOpen {
			t.Errorf("expected state open, got %s", cb.State())
		}

		// Wait for timeout
		time.Sleep(150 * time.Millisecond)

		// Circuit should be half-open now, try a successful request
		result, err := cb.Execute(func() (any, error) {
			return "recovered", nil
		})

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if result != "recovered" {
			t.Errorf("expected 'recovered', got %v", result)
		}

		// Circuit should be closed again
		if cb.State() != StateClosed {
			t.Errorf("expected state closed after recovery, got %s", cb.State())
		}
	})
}

func TestCircuitBreakerState(t *testing.T) {
	t.Run("initial state is closed", func(t *testing.T) {
		cb := NewCircuitBreaker(DefaultConfig("test"))

		if cb.State() != StateClosed {
			t.Errorf("expected closed, got %s", cb.State())
		}
		if !cb.IsClosed() {
			t.Error("expected IsClosed() to be true")
		}
		if cb.IsOpen() {
			t.Error("expected IsOpen() to be false")
		}
	})

	t.Run("IsOpen returns true when circuit is open", func(t *testing.T) {
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

		if !cb.IsOpen() {
			t.Error("expected IsOpen() to be true")
		}
		if cb.IsClosed() {
			t.Error("expected IsClosed() to be false")
		}
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
			if tt.state.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, tt.state.String())
			}
		})
	}
}

func TestCircuitBreakerCounts(t *testing.T) {
	t.Run("tracks successful requests", func(t *testing.T) {
		cb := NewCircuitBreaker(DefaultConfig("test"))

		for i := 0; i < 3; i++ {
			_, _ = cb.Execute(func() (any, error) {
				return nil, nil
			})
		}

		counts := cb.Counts()
		if counts.Requests != 3 {
			t.Errorf("expected 3 requests, got %d", counts.Requests)
		}
		if counts.TotalSuccesses != 3 {
			t.Errorf("expected 3 successes, got %d", counts.TotalSuccesses)
		}
		if counts.TotalFailures != 0 {
			t.Errorf("expected 0 failures, got %d", counts.TotalFailures)
		}
		if counts.ConsecutiveSuccesses != 3 {
			t.Errorf("expected 3 consecutive successes, got %d", counts.ConsecutiveSuccesses)
		}
	})

	t.Run("tracks failed requests", func(t *testing.T) {
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
		if counts.TotalFailures != 3 {
			t.Errorf("expected 3 failures, got %d", counts.TotalFailures)
		}
		if counts.ConsecutiveFailures != 3 {
			t.Errorf("expected 3 consecutive failures, got %d", counts.ConsecutiveFailures)
		}
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

	if len(stateChanges) == 0 {
		t.Fatal("expected state change callback to be called")
	}

	// First change should be closed -> open
	if stateChanges[0].from != StateClosed {
		t.Errorf("expected from state closed, got %s", stateChanges[0].from)
	}
	if stateChanges[0].to != StateOpen {
		t.Errorf("expected to state open, got %s", stateChanges[0].to)
	}
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
	if cb.State() != StateClosed {
		t.Errorf("expected state closed (ignored errors), got %s", cb.State())
	}

	counts := cb.Counts()
	if counts.TotalFailures != 0 {
		t.Errorf("expected 0 failures (ignored), got %d", counts.TotalFailures)
	}
}

func TestCircuitBreakerDefaultConfig(t *testing.T) {
	cfg := DefaultConfig("my-breaker")

	if cfg.Name != "my-breaker" {
		t.Errorf("expected name 'my-breaker', got %s", cfg.Name)
	}
	if cfg.MaxRequests != 1 {
		t.Errorf("expected MaxRequests 1, got %d", cfg.MaxRequests)
	}
	if cfg.Timeout != 60*time.Second {
		t.Errorf("expected Timeout 60s, got %v", cfg.Timeout)
	}
	if cfg.FailureThreshold != 5 {
		t.Errorf("expected FailureThreshold 5, got %d", cfg.FailureThreshold)
	}
	if cfg.SuccessThreshold != 1 {
		t.Errorf("expected SuccessThreshold 1, got %d", cfg.SuccessThreshold)
	}
	if cfg.Interval != 0 {
		t.Errorf("expected Interval 0, got %v", cfg.Interval)
	}
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

		if !errors.Is(err, utils.ErrCircuitOpen) {
			t.Errorf("expected ErrCircuitOpen, got %v", err)
		}
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
		if cb.State() != StateClosed {
			t.Fatalf("Step 1: expected initial state closed, got %s", cb.State())
		}

		// Step 2: Cause failures to transition to OPEN
		for i := 0; i < 3; i++ {
			_, err := cb.Execute(func() (any, error) {
				return nil, errors.New("failure")
			})
			if i < 2 && err == nil {
				t.Errorf("Step 2: expected error on failure %d", i)
			}
		}

		if cb.State() != StateOpen {
			t.Fatalf("Step 2: expected state open after 3 failures, got %s", cb.State())
		}

		// Verify first state transition: Closed -> Open
		if len(stateTransitions) < 1 {
			t.Fatal("Step 2: expected at least one state transition")
		}
		if stateTransitions[0].from != StateClosed || stateTransitions[0].to != StateOpen {
			t.Errorf("Step 2: expected transition closed->open, got %s->%s",
				stateTransitions[0].from, stateTransitions[0].to)
		}

		// Step 3: Verify requests are rejected while OPEN
		_, err := cb.Execute(func() (any, error) {
			return "should not execute", nil
		})
		if !errors.Is(err, utils.ErrCircuitOpen) {
			t.Errorf("Step 3: expected ErrCircuitOpen while open, got %v", err)
		}

		// Step 4: Wait for timeout to transition to HALF-OPEN
		time.Sleep(150 * time.Millisecond)

		// The state transitions to half-open on the next request attempt
		// Make a successful request to trigger the transition and close the circuit
		result, err := cb.Execute(func() (any, error) {
			return "recovery success", nil
		})
		if err != nil {
			t.Errorf("Step 4: unexpected error in half-open state: %v", err)
		}
		if result != "recovery success" {
			t.Errorf("Step 4: expected 'recovery success', got %v", result)
		}

		// Verify state transitions: Open -> Half-Open -> Closed
		if len(stateTransitions) < 3 {
			t.Fatalf("Step 4: expected 3 state transitions, got %d", len(stateTransitions))
		}
		if stateTransitions[1].from != StateOpen || stateTransitions[1].to != StateHalfOpen {
			t.Errorf("Step 4: expected transition open->half-open, got %s->%s",
				stateTransitions[1].from, stateTransitions[1].to)
		}
		if stateTransitions[2].from != StateHalfOpen || stateTransitions[2].to != StateClosed {
			t.Errorf("Step 4: expected transition half-open->closed, got %s->%s",
				stateTransitions[2].from, stateTransitions[2].to)
		}

		// Verify final state is CLOSED
		if cb.State() != StateClosed {
			t.Fatalf("Step 4: expected state closed after recovery, got %s", cb.State())
		}

		// Step 5: Verify circuit is fully operational again
		for i := 0; i < 5; i++ {
			result, err = cb.Execute(func() (any, error) {
				return "operational", nil
			})
			if err != nil {
				t.Errorf("Step 5: unexpected error on request %d: %v", i, err)
			}
		}

		counts := cb.Counts()
		if counts.TotalSuccesses < 5 {
			t.Errorf("Step 5: expected at least 5 successes, got %d", counts.TotalSuccesses)
		}
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
		for i := 0; i < 2; i++ {
			_, _ = cb.Execute(func() (any, error) {
				return nil, errors.New("failure")
			})
		}

		if cb.State() != StateOpen {
			t.Fatalf("Step 1: expected state open, got %s", cb.State())
		}

		// Step 2: Wait for timeout
		time.Sleep(150 * time.Millisecond)

		// Step 3: Fail in half-open state - should go back to OPEN
		_, err := cb.Execute(func() (any, error) {
			return nil, errors.New("failure in half-open")
		})
		if err == nil {
			t.Error("Step 3: expected error from failed execution")
		}

		// Circuit should be back to OPEN
		if cb.State() != StateOpen {
			t.Errorf("Step 3: expected state open after half-open failure, got %s", cb.State())
		}

		// Verify transitions: Closed -> Open -> Half-Open -> Open
		if len(stateTransitions) < 3 {
			t.Fatalf("Step 3: expected at least 3 state transitions, got %d", len(stateTransitions))
		}

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
				t.Errorf("Step 3: missing transition %d", i)
				continue
			}
			if stateTransitions[i].from != expected.from || stateTransitions[i].to != expected.to {
				t.Errorf("Step 3: transition %d: expected %s->%s, got %s->%s",
					i, expected.from, expected.to, stateTransitions[i].from, stateTransitions[i].to)
			}
		}
	})

	t.Run("multiple cycles of open and recovery", func(t *testing.T) {
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
			if cb.State() != StateClosed {
				t.Fatalf("Cycle %d: expected initial state closed, got %s", cycle, cb.State())
			}

			// Trip the circuit
			for i := 0; i < 2; i++ {
				_, _ = cb.Execute(func() (any, error) {
					return nil, errors.New("failure")
				})
			}

			if cb.State() != StateOpen {
				t.Fatalf("Cycle %d: expected state open, got %s", cycle, cb.State())
			}

			// Wait for half-open
			time.Sleep(60 * time.Millisecond)

			// Recover with success
			_, err := cb.Execute(func() (any, error) {
				return "recovered", nil
			})
			if err != nil {
				t.Errorf("Cycle %d: unexpected error during recovery: %v", cycle, err)
			}

			// Should be closed again
			if cb.State() != StateClosed {
				t.Fatalf("Cycle %d: expected state closed after recovery, got %s", cycle, cb.State())
			}
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
		if counts.ConsecutiveFailures != 3 {
			t.Errorf("expected 3 consecutive failures, got %d", counts.ConsecutiveFailures)
		}
		if counts.TotalFailures != 3 {
			t.Errorf("expected 3 total failures, got %d", counts.TotalFailures)
		}

		// Circuit should still be closed
		if cb.State() != StateClosed {
			t.Errorf("expected state closed, got %s", cb.State())
		}

		// Wait for the interval to pass
		time.Sleep(150 * time.Millisecond)

		// Make a successful request to trigger the count reset check
		_, _ = cb.Execute(func() (any, error) {
			return "success", nil
		})

		// Counts should be reset (only the new successful request should be counted)
		counts = cb.Counts()
		if counts.ConsecutiveFailures != 0 {
			t.Errorf("expected 0 consecutive failures after interval reset, got %d", counts.ConsecutiveFailures)
		}
		if counts.TotalFailures != 0 {
			t.Errorf("expected 0 total failures after interval reset, got %d", counts.TotalFailures)
		}
		if counts.TotalSuccesses != 1 {
			t.Errorf("expected 1 total success after interval reset, got %d", counts.TotalSuccesses)
		}
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

		if cb.State() != StateClosed {
			t.Errorf("expected state closed after 2 failures, got %s", cb.State())
		}

		// Wait for interval to pass and counts to reset
		time.Sleep(150 * time.Millisecond)

		// Cause 2 more failures in second interval (would be 4 total without reset)
		for i := 0; i < 2; i++ {
			_, _ = cb.Execute(func() (any, error) {
				return nil, errors.New("failure")
			})
		}

		// Circuit should still be closed because counts were reset
		if cb.State() != StateClosed {
			t.Errorf("expected state closed (counts should have reset), got %s", cb.State())
		}

		counts := cb.Counts()
		if counts.ConsecutiveFailures != 2 {
			t.Errorf("expected 2 consecutive failures in new interval, got %d", counts.ConsecutiveFailures)
		}
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
		if cb.State() != StateOpen {
			t.Errorf("expected state open, got %s", cb.State())
		}
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
		counts := cb.Counts()
		if counts.TotalFailures != 4 {
			t.Errorf("expected 4 total failures (no reset), got %d", counts.TotalFailures)
		}
		if counts.ConsecutiveFailures != 4 {
			t.Errorf("expected 4 consecutive failures (no reset), got %d", counts.ConsecutiveFailures)
		}
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
		if successCount.Load() != numGoroutines {
			t.Errorf("successCount: expected %d, got %d", numGoroutines, successCount.Load())
		}
		if errorCount.Load() != 0 {
			t.Errorf("errorCount: expected 0, got %d", errorCount.Load())
		}

		counts := cb.Counts()
		if counts.Requests != numGoroutines {
			t.Errorf("counts.Requests: expected %d, got %d", numGoroutines, counts.Requests)
		}
		if counts.TotalSuccesses != numGoroutines {
			t.Errorf("counts.TotalSuccesses: expected %d, got %d", numGoroutines, counts.TotalSuccesses)
		}
		if counts.TotalFailures != 0 {
			t.Errorf("counts.TotalFailures: expected 0, got %d", counts.TotalFailures)
		}
		if counts.ConsecutiveSuccesses != numGoroutines {
			t.Errorf("counts.ConsecutiveSuccesses: expected %d, got %d", numGoroutines, counts.ConsecutiveSuccesses)
		}
		if counts.ConsecutiveFailures != 0 {
			t.Errorf("counts.ConsecutiveFailures: expected 0, got %d", counts.ConsecutiveFailures)
		}
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
			wg.Add(1)
			go func() {
				defer wg.Done()
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
			}()
		}

		wg.Wait()

		// Circuit should still be closed (threshold not reached)
		if cb.State() != StateClosed {
			t.Errorf("expected state closed, got %s", cb.State())
		}

		// Assert precise counts
		if failureCount.Load() != numGoroutines {
			t.Errorf("failureCount: expected %d, got %d", numGoroutines, failureCount.Load())
		}
		if successCount.Load() != 0 {
			t.Errorf("successCount: expected 0, got %d", successCount.Load())
		}
		if circuitOpenCount.Load() != 0 {
			t.Errorf("circuitOpenCount: expected 0, got %d", circuitOpenCount.Load())
		}

		counts := cb.Counts()
		if counts.Requests != numGoroutines {
			t.Errorf("counts.Requests: expected %d, got %d", numGoroutines, counts.Requests)
		}
		if counts.TotalFailures != numGoroutines {
			t.Errorf("counts.TotalFailures: expected %d, got %d", numGoroutines, counts.TotalFailures)
		}
		if counts.TotalSuccesses != 0 {
			t.Errorf("counts.TotalSuccesses: expected 0, got %d", counts.TotalSuccesses)
		}
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

		if cb.State() != StateOpen {
			t.Fatalf("expected state open after 3 failures, got %s", cb.State())
		}

		// Phase 2: All concurrent requests should be rejected
		var wg sync.WaitGroup
		var circuitOpenCount atomic.Int64
		var otherErrorCount atomic.Int64
		var successCount atomic.Int64
		const numGoroutines = 50

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
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
			}()
		}

		wg.Wait()

		// Assert precise counts - ALL requests should be rejected
		if circuitOpenCount.Load() != numGoroutines {
			t.Errorf("circuitOpenCount: expected %d, got %d", numGoroutines, circuitOpenCount.Load())
		}
		if successCount.Load() != 0 {
			t.Errorf("successCount: expected 0, got %d", successCount.Load())
		}
		if otherErrorCount.Load() != 0 {
			t.Errorf("otherErrorCount: expected 0, got %d", otherErrorCount.Load())
		}

		// Note: Counts are reset when circuit transitions to open state
		// Rejected requests (ErrCircuitOpen) are not counted in the circuit breaker stats
		counts := cb.Counts()
		if counts.Requests != 0 {
			t.Errorf("counts.Requests: expected 0 (reset on open), got %d", counts.Requests)
		}
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
		if cb.State() != StateClosed {
			t.Errorf("expected state closed, got %s", cb.State())
		}

		// Assert precise counts
		if successCount.Load() != 50 {
			t.Errorf("successCount: expected 50, got %d", successCount.Load())
		}
		if failureCount.Load() != 50 {
			t.Errorf("failureCount: expected 50, got %d", failureCount.Load())
		}
		if circuitOpenCount.Load() != 0 {
			t.Errorf("circuitOpenCount: expected 0, got %d", circuitOpenCount.Load())
		}

		counts := cb.Counts()
		if counts.Requests != numGoroutines {
			t.Errorf("counts.Requests: expected %d, got %d", numGoroutines, counts.Requests)
		}
		if counts.TotalSuccesses != 50 {
			t.Errorf("counts.TotalSuccesses: expected 50, got %d", counts.TotalSuccesses)
		}
		if counts.TotalFailures != 50 {
			t.Errorf("counts.TotalFailures: expected 50, got %d", counts.TotalFailures)
		}
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
		counts := cb.Counts()
		expectedRequests := uint32(numWriters * writesPerGoroutine)
		if counts.Requests != expectedRequests {
			t.Errorf("counts.Requests: expected %d, got %d", expectedRequests, counts.Requests)
		}
		if counts.TotalSuccesses != expectedRequests {
			t.Errorf("counts.TotalSuccesses: expected %d, got %d", expectedRequests, counts.TotalSuccesses)
		}
		if counts.TotalFailures != 0 {
			t.Errorf("counts.TotalFailures: expected 0, got %d", counts.TotalFailures)
		}
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

		if cb.State() != StateOpen {
			t.Fatalf("expected state open, got %s", cb.State())
		}

		// Phase 2: Wait for circuit to transition to half-open
		time.Sleep(60 * time.Millisecond)

		// Phase 3: Single successful request to close the circuit
		_, err := cb.Execute(func() (any, error) {
			return "success", nil
		})
		if err != nil {
			t.Errorf("expected successful execution in half-open, got error: %v", err)
		}

		if cb.State() != StateClosed {
			t.Fatalf("expected state closed after recovery, got %s", cb.State())
		}

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
		if successCount.Load() != numGoroutines {
			t.Errorf("successCount: expected %d, got %d", numGoroutines, successCount.Load())
		}
		if errorCount.Load() != 0 {
			t.Errorf("errorCount: expected 0, got %d", errorCount.Load())
		}
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
		if cb.State() != StateClosed {
			t.Fatalf("Phase 1: expected initial state closed, got %s", cb.State())
		}

		// Phase 2: Trip the circuit sequentially (predictable 5 failures)
		for i := 0; i < 5; i++ {
			_, _ = cb.Execute(func() (any, error) {
				return nil, errors.New("failure")
			})
		}

		if cb.State() != StateOpen {
			t.Fatalf("Phase 2: expected state open after 5 failures, got %s", cb.State())
		}

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
		if rejectedCount.Load() != numRejectedGoroutines {
			t.Errorf("Phase 3 rejectedCount: expected %d, got %d", numRejectedGoroutines, rejectedCount.Load())
		}
		if otherCount.Load() != 0 {
			t.Errorf("Phase 3 otherCount: expected 0, got %d", otherCount.Load())
		}

		// Phase 4: Wait for timeout to transition to HALF-OPEN
		time.Sleep(150 * time.Millisecond)

		// Phase 5: Single request to recover (MaxRequests=1)
		_, err := cb.Execute(func() (any, error) {
			return "recovery", nil
		})
		if err != nil {
			t.Errorf("Phase 5: expected successful recovery, got error: %v", err)
		}

		if cb.State() != StateClosed {
			t.Errorf("Phase 5: expected state closed after recovery, got %s", cb.State())
		}

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
		if finalSuccessCount.Load() != numFinalGoroutines {
			t.Errorf("Phase 6 finalSuccessCount: expected %d, got %d", numFinalGoroutines, finalSuccessCount.Load())
		}
		if finalErrorCount.Load() != 0 {
			t.Errorf("Phase 6 finalErrorCount: expected 0, got %d", finalErrorCount.Load())
		}

		// Verify state transitions
		transitionMu.Lock()
		defer transitionMu.Unlock()

		if len(stateTransitions) != 3 {
			t.Errorf("Expected exactly 3 state transitions, got %d", len(stateTransitions))
		}
		if len(stateTransitions) >= 1 && (stateTransitions[0].from != StateClosed || stateTransitions[0].to != StateOpen) {
			t.Errorf("Transition 1: expected closed->open, got %s->%s", stateTransitions[0].from, stateTransitions[0].to)
		}
		if len(stateTransitions) >= 2 && (stateTransitions[1].from != StateOpen || stateTransitions[1].to != StateHalfOpen) {
			t.Errorf("Transition 2: expected open->half-open, got %s->%s", stateTransitions[1].from, stateTransitions[1].to)
		}
		if len(stateTransitions) >= 3 && (stateTransitions[2].from != StateHalfOpen || stateTransitions[2].to != StateClosed) {
			t.Errorf("Transition 3: expected half-open->closed, got %s->%s", stateTransitions[2].from, stateTransitions[2].to)
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

		if cb.State() != StateOpen {
			t.Fatalf("Phase 1: expected state open, got %s", cb.State())
		}

		// Note: Counts are reset when circuit transitions to open state
		counts := cb.Counts()
		if counts.Requests != 0 {
			t.Errorf("Phase 1 counts.Requests: expected 0 (reset on open), got %d", counts.Requests)
		}

		// Phase 2: Wait for half-open
		time.Sleep(150 * time.Millisecond)

		// Phase 3: Single failing request in half-open state
		_, err := cb.Execute(func() (any, error) {
			return nil, errors.New("still failing")
		})
		if err == nil {
			t.Error("Phase 3: expected error from failed execution")
		}

		if cb.State() != StateOpen {
			t.Errorf("Phase 3: expected state open after half-open failure, got %s", cb.State())
		}

		// Phase 4: Wait again and recover sequentially
		time.Sleep(150 * time.Millisecond)

		_, err = cb.Execute(func() (any, error) {
			return "recovered", nil
		})
		if err != nil {
			t.Errorf("Phase 4: expected successful recovery, got error: %v", err)
		}

		if cb.State() != StateClosed {
			t.Errorf("Phase 4: expected state closed after recovery, got %s", cb.State())
		}

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
		if successCount.Load() != numGoroutines {
			t.Errorf("Phase 5 successCount: expected %d, got %d", numGoroutines, successCount.Load())
		}
		if errorCount.Load() != 0 {
			t.Errorf("Phase 5 errorCount: expected 0, got %d", errorCount.Load())
		}
	})

	t.Run("multiple sequential recovery cycles with concurrent verification", func(t *testing.T) {
		cfg := Config{
			Name:             "multi-cycle-test",
			FailureThreshold: 3,
			MaxRequests:      1,
			Timeout:          50 * time.Millisecond,
		}
		cb := NewCircuitBreaker(cfg)

		for cycle := 1; cycle <= 3; cycle++ {
			// Verify starting in closed state
			if cb.State() != StateClosed {
				t.Fatalf("Cycle %d: expected initial state closed, got %s", cycle, cb.State())
			}

			// Sequential failures to trip the circuit (predictable)
			for i := 0; i < 3; i++ {
				_, _ = cb.Execute(func() (any, error) {
					return nil, errors.New("failure")
				})
			}

			if cb.State() != StateOpen {
				t.Fatalf("Cycle %d: expected state open, got %s", cycle, cb.State())
			}

			// Wait for half-open
			time.Sleep(60 * time.Millisecond)

			// Single recovery request
			_, err := cb.Execute(func() (any, error) {
				return "recovered", nil
			})
			if err != nil {
				t.Errorf("Cycle %d: expected successful recovery, got error: %v", cycle, err)
			}

			if cb.State() != StateClosed {
				t.Fatalf("Cycle %d: expected state closed after recovery, got %s", cycle, cb.State())
			}

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
			if successCount.Load() != numVerifyGoroutines {
				t.Errorf("Cycle %d verify successCount: expected %d, got %d", cycle, numVerifyGoroutines, successCount.Load())
			}
			if errorCount.Load() != 0 {
				t.Errorf("Cycle %d verify errorCount: expected 0, got %d", cycle, errorCount.Load())
			}
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
		if cb.State() != StateClosed {
			t.Errorf("expected state closed, got %s", cb.State())
		}

		// Assert precise counts
		const expectedSuccesses = 8000 // 800 goroutines * 10 requests
		const expectedFailures = 2000  // 200 goroutines * 10 requests
		const expectedTotal = numGoroutines * requestsPerGoroutine

		if successCount.Load() != expectedSuccesses {
			t.Errorf("successCount: expected %d, got %d", expectedSuccesses, successCount.Load())
		}
		if failureCount.Load() != expectedFailures {
			t.Errorf("failureCount: expected %d, got %d", expectedFailures, failureCount.Load())
		}

		counts := cb.Counts()
		if counts.Requests != expectedTotal {
			t.Errorf("counts.Requests: expected %d, got %d", expectedTotal, counts.Requests)
		}
		if counts.TotalSuccesses != expectedSuccesses {
			t.Errorf("counts.TotalSuccesses: expected %d, got %d", expectedSuccesses, counts.TotalSuccesses)
		}
		if counts.TotalFailures != expectedFailures {
			t.Errorf("counts.TotalFailures: expected %d, got %d", expectedFailures, counts.TotalFailures)
		}
	})
}
