package middleware

import (
	"context"
	"errors"
	"testing"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

func TestCircuitBreakerInitialState(t *testing.T) {
	cb := NewCircuitBreaker(DefaultCircuitBreakerConfig())

	if cb.State() != CircuitClosed {
		t.Errorf("expected initial state to be closed, got %v", cb.State())
	}
}

func TestCircuitBreakerOpensAfterFailures(t *testing.T) {
	cfg := CircuitBreakerConfig{
		FailureThreshold: 3,
		SuccessThreshold: 1,
		Timeout:          100 * time.Millisecond,
	}
	cb := NewCircuitBreaker(cfg)

	testErr := errors.New("test error")

	// Record failures
	for i := 0; i < 3; i++ {
		cb.RecordFailure(testErr)
	}

	if cb.State() != CircuitOpen {
		t.Errorf("expected circuit to be open after %d failures, got %v", cfg.FailureThreshold, cb.State())
	}
}

func TestCircuitBreakerBlocksWhenOpen(t *testing.T) {
	cfg := CircuitBreakerConfig{
		FailureThreshold: 1,
		SuccessThreshold: 1,
		Timeout:          1 * time.Hour, // Long timeout to ensure it stays open
	}
	cb := NewCircuitBreaker(cfg)

	// Open the circuit
	cb.RecordFailure(errors.New("test"))

	// Should block
	err := cb.Allow()
	if err != ErrCircuitOpen {
		t.Errorf("expected ErrCircuitOpen, got %v", err)
	}
}

func TestCircuitBreakerTransitionsToHalfOpen(t *testing.T) {
	cfg := CircuitBreakerConfig{
		FailureThreshold: 1,
		SuccessThreshold: 1,
		Timeout:          10 * time.Millisecond,
	}
	cb := NewCircuitBreaker(cfg)

	// Open the circuit
	cb.RecordFailure(errors.New("test"))

	if cb.State() != CircuitOpen {
		t.Errorf("expected open state, got %v", cb.State())
	}

	// Wait for timeout
	time.Sleep(20 * time.Millisecond)

	// Should allow and transition to half-open
	err := cb.Allow()
	if err != nil {
		t.Errorf("expected nil error after timeout, got %v", err)
	}

	if cb.State() != CircuitHalfOpen {
		t.Errorf("expected half-open state, got %v", cb.State())
	}
}

func TestCircuitBreakerClosesAfterSuccess(t *testing.T) {
	cfg := CircuitBreakerConfig{
		FailureThreshold: 1,
		SuccessThreshold: 2,
		Timeout:          10 * time.Millisecond,
	}
	cb := NewCircuitBreaker(cfg)

	// Open the circuit
	cb.RecordFailure(errors.New("test"))

	// Wait for timeout and transition to half-open
	time.Sleep(20 * time.Millisecond)
	cb.Allow()

	// Record successes
	cb.RecordSuccess()
	cb.RecordSuccess()

	if cb.State() != CircuitClosed {
		t.Errorf("expected closed state after successes, got %v", cb.State())
	}
}

func TestCircuitBreakerReopensOnFailureInHalfOpen(t *testing.T) {
	cfg := CircuitBreakerConfig{
		FailureThreshold: 1,
		SuccessThreshold: 2,
		Timeout:          10 * time.Millisecond,
	}
	cb := NewCircuitBreaker(cfg)

	// Open the circuit
	cb.RecordFailure(errors.New("test1"))

	// Wait for timeout and transition to half-open
	time.Sleep(20 * time.Millisecond)
	cb.Allow()

	if cb.State() != CircuitHalfOpen {
		t.Errorf("expected half-open state, got %v", cb.State())
	}

	// Record a failure
	cb.RecordFailure(errors.New("test2"))

	if cb.State() != CircuitOpen {
		t.Errorf("expected open state after failure in half-open, got %v", cb.State())
	}
}

func TestCircuitBreakerReset(t *testing.T) {
	cfg := CircuitBreakerConfig{
		FailureThreshold: 1,
		SuccessThreshold: 1,
		Timeout:          1 * time.Hour,
	}
	cb := NewCircuitBreaker(cfg)

	// Open the circuit
	cb.RecordFailure(errors.New("test"))

	if cb.State() != CircuitOpen {
		t.Errorf("expected open state, got %v", cb.State())
	}

	// Reset
	cb.Reset()

	if cb.State() != CircuitClosed {
		t.Errorf("expected closed state after reset, got %v", cb.State())
	}
}

func TestCircuitBreakerStateCallback(t *testing.T) {
	var transitions []CircuitState
	cfg := CircuitBreakerConfig{
		FailureThreshold: 1,
		SuccessThreshold: 1,
		Timeout:          10 * time.Millisecond,
		OnStateChange: func(from, to CircuitState) {
			transitions = append(transitions, to)
		},
	}
	cb := NewCircuitBreaker(cfg)

	// Open
	cb.RecordFailure(errors.New("test"))

	// Wait and go to half-open
	time.Sleep(20 * time.Millisecond)
	cb.Allow()

	// Close
	cb.RecordSuccess()

	expected := []CircuitState{CircuitOpen, CircuitHalfOpen, CircuitClosed}
	if len(transitions) != len(expected) {
		t.Errorf("expected %d transitions, got %d", len(expected), len(transitions))
	}

	for i, s := range expected {
		if i < len(transitions) && transitions[i] != s {
			t.Errorf("transition %d: expected %v, got %v", i, s, transitions[i])
		}
	}
}

func TestCircuitBreakerMiddleware(t *testing.T) {
	cfg := CircuitBreakerConfig{
		FailureThreshold: 2,
		SuccessThreshold: 1,
		Timeout:          100 * time.Millisecond,
	}
	cb := NewCircuitBreaker(cfg)

	callCount := 0
	failingHandler := func(ctx context.Context, subject string, msg *core.Message) error {
		callCount++
		return errors.New("simulated failure")
	}

	wrapped := CircuitBreakerMiddleware(cb)(failingHandler)
	ctx := context.Background()
	msg := core.NewMessage([]byte("test"))

	// First two calls should go through (and fail)
	wrapped(ctx, "test", msg)
	wrapped(ctx, "test", msg)

	// Circuit should now be open
	if cb.State() != CircuitOpen {
		t.Errorf("expected open state, got %v", cb.State())
	}

	// Third call should be blocked
	err := wrapped(ctx, "test", msg)
	if err != ErrCircuitOpen {
		t.Errorf("expected ErrCircuitOpen, got %v", err)
	}

	// Handler should only have been called twice
	if callCount != 2 {
		t.Errorf("expected handler to be called 2 times, got %d", callCount)
	}
}

func TestMultiCircuitBreaker(t *testing.T) {
	cfg := CircuitBreakerConfig{
		FailureThreshold: 1,
		SuccessThreshold: 1,
		Timeout:          100 * time.Millisecond,
	}
	mcb := NewMultiCircuitBreaker(cfg)

	// Get circuit breakers for different subjects
	cb1 := mcb.Get("subject1")
	cb2 := mcb.Get("subject2")

	// Open circuit for subject1
	cb1.RecordFailure(errors.New("test"))

	if cb1.State() != CircuitOpen {
		t.Errorf("expected subject1 circuit to be open, got %v", cb1.State())
	}

	if cb2.State() != CircuitClosed {
		t.Errorf("expected subject2 circuit to be closed, got %v", cb2.State())
	}

	// Check states
	states := mcb.States()
	if states["subject1"] != CircuitOpen {
		t.Errorf("expected subject1 state to be open in map")
	}
	if states["subject2"] != CircuitClosed {
		t.Errorf("expected subject2 state to be closed in map")
	}
}

func TestCircuitBreakerShouldTrip(t *testing.T) {
	cfg := CircuitBreakerConfig{
		FailureThreshold: 1,
		SuccessThreshold: 1,
		Timeout:          100 * time.Millisecond,
		ShouldTrip: func(err error) bool {
			// Only trip on specific errors
			return err != nil && err.Error() == "critical"
		},
	}
	cb := NewCircuitBreaker(cfg)

	// Non-critical error should not trip
	cb.RecordFailure(errors.New("non-critical"))
	if cb.State() != CircuitClosed {
		t.Errorf("expected circuit to remain closed for non-critical error")
	}

	// Critical error should trip
	cb.RecordFailure(errors.New("critical"))
	if cb.State() != CircuitOpen {
		t.Errorf("expected circuit to open for critical error")
	}
}
