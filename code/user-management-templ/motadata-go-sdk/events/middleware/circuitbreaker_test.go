package middleware

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

/* ========================================================================================================
   CIRCUIT BREAKER STATE TESTS
   ======================================================================================================== */

var (
	cbTestErr         = errors.New(testErrMsg)
	cbCriticalErr     = errors.New("critical")
	cbNonCriticalErr  = errors.New("non-critical")
	cbShortTimeout    = 10 * time.Millisecond
	cbLongTimeout     = 1 * time.Hour
	cbSleepForTimeout = 20 * time.Millisecond
)

func TestCircuitBreakerInitialState(t *testing.T) {
	assertions := assert.New(t)

	cb := NewCircuitBreaker(DefaultCircuitBreakerConfig())

	assertions.Equal(CircuitClosed, cb.State(), "initial state should be closed")
}

func TestCircuitBreakerOpensAfterFailures(t *testing.T) {
	assertions := assert.New(t)

	cfg := CircuitBreakerConfig{
		FailureThreshold: 3,
		SuccessThreshold: 1,
		Timeout:          100 * time.Millisecond,
	}
	cb := NewCircuitBreaker(cfg)

	for i := 0; i < 3; i++ {
		cb.RecordFailure(cbTestErr)
	}

	assertions.Equal(CircuitOpen, cb.State(), "circuit should be open after threshold failures")
}

func TestCircuitBreakerBlocksWhenOpen(t *testing.T) {
	assertions := assert.New(t)

	cfg := CircuitBreakerConfig{
		FailureThreshold: 1,
		SuccessThreshold: 1,
		Timeout:          cbLongTimeout,
	}
	cb := NewCircuitBreaker(cfg)
	cb.RecordFailure(cbTestErr)

	err := cb.Allow()

	assertions.ErrorIs(err, ErrCircuitOpen, "should return ErrCircuitOpen when open")
}

func TestCircuitBreakerTransitionsToHalfOpen(t *testing.T) {
	assertions := assert.New(t)

	cfg := CircuitBreakerConfig{
		FailureThreshold: 1,
		SuccessThreshold: 1,
		Timeout:          cbShortTimeout,
	}
	cb := NewCircuitBreaker(cfg)
	cb.RecordFailure(cbTestErr)

	assertions.Equal(CircuitOpen, cb.State())

	time.Sleep(cbSleepForTimeout)
	err := cb.Allow()

	assertions.NoError(err, "should allow after timeout")
	assertions.Equal(CircuitHalfOpen, cb.State(), "should transition to half-open")
}

func TestCircuitBreakerClosesAfterSuccess(t *testing.T) {
	assertions := assert.New(t)

	cfg := CircuitBreakerConfig{
		FailureThreshold: 1,
		SuccessThreshold: 2,
		Timeout:          cbShortTimeout,
	}
	cb := NewCircuitBreaker(cfg)
	cb.RecordFailure(cbTestErr)

	time.Sleep(cbSleepForTimeout)
	cb.Allow()

	cb.RecordSuccess()
	cb.RecordSuccess()

	assertions.Equal(CircuitClosed, cb.State(), "should close after enough successes in half-open")
}

func TestCircuitBreakerReopensOnFailureInHalfOpen(t *testing.T) {
	assertions := assert.New(t)

	cfg := CircuitBreakerConfig{
		FailureThreshold: 1,
		SuccessThreshold: 2,
		Timeout:          cbShortTimeout,
	}
	cb := NewCircuitBreaker(cfg)
	cb.RecordFailure(errors.New("test1"))

	time.Sleep(cbSleepForTimeout)
	cb.Allow()

	assertions.Equal(CircuitHalfOpen, cb.State())

	cb.RecordFailure(errors.New("test2"))

	assertions.Equal(CircuitOpen, cb.State(), "should reopen on failure in half-open")
}

func TestCircuitBreakerReset(t *testing.T) {
	assertions := assert.New(t)

	cfg := CircuitBreakerConfig{
		FailureThreshold: 1,
		SuccessThreshold: 1,
		Timeout:          cbLongTimeout,
	}
	cb := NewCircuitBreaker(cfg)
	cb.RecordFailure(cbTestErr)

	assertions.Equal(CircuitOpen, cb.State())

	cb.Reset()

	assertions.Equal(CircuitClosed, cb.State(), "should be closed after reset")
}

/* ========================================================================================================
   CIRCUIT BREAKER CALLBACK TESTS
   ======================================================================================================== */

func TestCircuitBreakerStateCallback(t *testing.T) {
	assertions := assert.New(t)

	var transitions []CircuitState
	cfg := CircuitBreakerConfig{
		FailureThreshold: 1,
		SuccessThreshold: 1,
		Timeout:          cbShortTimeout,
		OnStateChange: func(from, to CircuitState) {
			transitions = append(transitions, to)
		},
	}
	cb := NewCircuitBreaker(cfg)

	cb.RecordFailure(cbTestErr)

	time.Sleep(cbSleepForTimeout)
	cb.Allow()

	cb.RecordSuccess()

	expected := []CircuitState{CircuitOpen, CircuitHalfOpen, CircuitClosed}
	assertions.Equal(expected, transitions, "state transitions should match")
}

/* ========================================================================================================
   CIRCUIT BREAKER MIDDLEWARE TESTS
   ======================================================================================================== */

func TestCircuitBreakerMiddleware(t *testing.T) {
	assertions := assert.New(t)

	cfg := CircuitBreakerConfig{
		FailureThreshold: 2,
		SuccessThreshold: 1,
		Timeout:          100 * time.Millisecond,
	}
	cb := NewCircuitBreaker(cfg)

	callCount := 0
	failingHandler := func(ctx context.Context, msg *nats.Msg) error {
		callCount++
		return errSimulated
	}

	wrapped := CircuitBreakerMiddleware(cb)(failingHandler)
	msg := newTestMsg(testSubject, testPayload)
	ctx := context.Background()

	wrapped(ctx, msg)
	wrapped(ctx, msg)

	assertions.Equal(CircuitOpen, cb.State(), "circuit should be open after failures")

	err := wrapped(ctx, msg)

	assertions.ErrorIs(err, ErrCircuitOpen, "third call should be blocked")
	assertions.Equal(2, callCount, "handler should only be called twice")
}

/* ========================================================================================================
   MULTI CIRCUIT BREAKER TESTS
   ======================================================================================================== */

func TestMultiCircuitBreaker(t *testing.T) {
	assertions := assert.New(t)

	cfg := CircuitBreakerConfig{
		FailureThreshold: 1,
		SuccessThreshold: 1,
		Timeout:          100 * time.Millisecond,
	}
	mcb := NewMultiCircuitBreaker(cfg)

	cb1 := mcb.Get(testSubject1)
	cb2 := mcb.Get(testSubject2)

	cb1.RecordFailure(cbTestErr)

	assertions.Equal(CircuitOpen, cb1.State(), "subject1 circuit should be open")
	assertions.Equal(CircuitClosed, cb2.State(), "subject2 circuit should be closed")

	states := mcb.States()
	assertions.Equal(CircuitOpen, states[testSubject1])
	assertions.Equal(CircuitClosed, states[testSubject2])
}

/* ========================================================================================================
   SHOULD TRIP TESTS
   ======================================================================================================== */

func TestCircuitBreakerShouldTrip(t *testing.T) {
	testCases := []struct {
		name          string
		err           error
		expectedState CircuitState
	}{
		{"non-critical error stays closed", cbNonCriticalErr, CircuitClosed},
		{"critical error opens circuit", cbCriticalErr, CircuitOpen},
	}

	cfg := CircuitBreakerConfig{
		FailureThreshold: 1,
		SuccessThreshold: 1,
		Timeout:          100 * time.Millisecond,
		ShouldTrip: func(err error) bool {
			return err != nil && err.Error() == "critical"
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)

			cb := NewCircuitBreaker(cfg)
			cb.RecordFailure(tc.err)

			assertions.Equal(tc.expectedState, cb.State())
		})
	}
}
