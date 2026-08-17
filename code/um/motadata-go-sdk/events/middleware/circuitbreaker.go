package middleware

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

// CircuitState represents the circuit breaker state
type CircuitState int32

const (
	// CircuitClosed allows requests to pass through
	CircuitClosed CircuitState = iota
	// CircuitOpen blocks all requests
	CircuitOpen
	// CircuitHalfOpen allows limited requests to test recovery
	CircuitHalfOpen
)

func (s CircuitState) String() string {
	switch s {
	case CircuitClosed:
		return "closed"
	case CircuitOpen:
		return "open"
	case CircuitHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// CircuitBreakerConfig configures the circuit breaker
type CircuitBreakerConfig struct {
	// FailureThreshold is the number of failures before opening the circuit
	FailureThreshold int

	// SuccessThreshold is the number of successes in half-open state before closing
	SuccessThreshold int

	// Timeout is how long the circuit stays open before transitioning to half-open
	Timeout time.Duration

	// OnStateChange is called when the circuit state changes
	OnStateChange func(from, to CircuitState)

	// ShouldTrip determines if an error should count as a failure
	ShouldTrip func(error) bool
}

// DefaultCircuitBreakerConfig returns default configuration
func DefaultCircuitBreakerConfig() CircuitBreakerConfig {
	return CircuitBreakerConfig{
		FailureThreshold: 5,
		SuccessThreshold: 2,
		Timeout:          30 * time.Second,
		ShouldTrip: func(err error) bool {
			// All errors trip the circuit by default
			return err != nil
		},
	}
}

// CircuitBreaker implements the circuit breaker pattern
type CircuitBreaker struct {
	config CircuitBreakerConfig

	state           atomic.Int32
	failures        atomic.Int64
	successes       atomic.Int64
	lastFailureTime atomic.Int64

	mu sync.Mutex
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(cfg CircuitBreakerConfig) *CircuitBreaker {
	if cfg.FailureThreshold <= 0 {
		cfg.FailureThreshold = 5
	}
	if cfg.SuccessThreshold <= 0 {
		cfg.SuccessThreshold = 2
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.ShouldTrip == nil {
		cfg.ShouldTrip = func(err error) bool { return err != nil }
	}

	cb := &CircuitBreaker{
		config: cfg,
	}
	cb.state.Store(int32(CircuitClosed))
	return cb
}

// State returns the current circuit state
func (cb *CircuitBreaker) State() CircuitState {
	return CircuitState(cb.state.Load())
}

// Allow checks if a request is allowed to proceed
func (cb *CircuitBreaker) Allow() error {
	state := cb.State()

	switch state {
	case CircuitClosed:
		return nil

	case CircuitOpen:
		// Check if timeout has passed
		lastFailure := time.Unix(0, cb.lastFailureTime.Load())
		if time.Since(lastFailure) > cb.config.Timeout {
			// Transition to half-open
			cb.transitionTo(CircuitHalfOpen)
			return nil
		}
		return ErrCircuitOpen

	case CircuitHalfOpen:
		return nil
	}

	return nil
}

// RecordSuccess records a successful operation
func (cb *CircuitBreaker) RecordSuccess() {
	state := cb.State()

	switch state {
	case CircuitClosed:
		cb.failures.Store(0) // Reset failures on success

	case CircuitHalfOpen:
		successes := cb.successes.Add(1)
		if int(successes) >= cb.config.SuccessThreshold {
			cb.transitionTo(CircuitClosed)
		}
	}
}

// RecordFailure records a failed operation
func (cb *CircuitBreaker) RecordFailure(err error) {
	if !cb.config.ShouldTrip(err) {
		return
	}

	cb.lastFailureTime.Store(time.Now().UnixNano())
	state := cb.State()

	switch state {
	case CircuitClosed:
		failures := cb.failures.Add(1)
		if int(failures) >= cb.config.FailureThreshold {
			cb.transitionTo(CircuitOpen)
		}

	case CircuitHalfOpen:
		// Any failure in half-open state re-opens the circuit
		cb.transitionTo(CircuitOpen)
	}
}

// transitionTo transitions to a new state
func (cb *CircuitBreaker) transitionTo(newState CircuitState) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	oldState := CircuitState(cb.state.Load())
	if oldState == newState {
		return
	}

	// Reset counters on state change
	cb.failures.Store(0)
	cb.successes.Store(0)

	cb.state.Store(int32(newState))

	if cb.config.OnStateChange != nil {
		cb.config.OnStateChange(oldState, newState)
	}
}

// Reset resets the circuit breaker to closed state
func (cb *CircuitBreaker) Reset() {
	cb.transitionTo(CircuitClosed)
}

// Failures returns the current failure count
func (cb *CircuitBreaker) Failures() int {
	return int(cb.failures.Load())
}

// ErrCircuitOpen is returned when the circuit breaker is open
var ErrCircuitOpen = &core.Error{
	Kind: "circuit_breaker",
	Op:   "allow",
	Err:  core.ErrPublishFailed,
}

// CircuitBreakerMiddleware returns publish middleware that applies circuit breaker logic
func CircuitBreakerMiddleware(cb *CircuitBreaker) PublishMiddleware {
	return func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, subject string, msg *core.Message) error {
			if err := cb.Allow(); err != nil {
				return err
			}

			err := next(ctx, subject, msg)
			if err != nil {
				cb.RecordFailure(err)
			} else {
				cb.RecordSuccess()
			}

			return err
		}
	}
}

// CircuitBreakerSubscribeMiddleware returns subscribe middleware with circuit breaker
func CircuitBreakerSubscribeMiddleware(cb *CircuitBreaker) SubscribeMiddleware {
	return func(next SubscribeHandler) SubscribeHandler {
		return func(ctx context.Context, msg *core.Message) error {
			if err := cb.Allow(); err != nil {
				return err
			}

			err := next(ctx, msg)
			if err != nil {
				cb.RecordFailure(err)
			} else {
				cb.RecordSuccess()
			}

			return err
		}
	}
}

// MultiCircuitBreaker manages circuit breakers per subject/destination
type MultiCircuitBreaker struct {
	config   CircuitBreakerConfig
	breakers sync.Map // map[string]*CircuitBreaker
}

// NewMultiCircuitBreaker creates a multi-circuit breaker
func NewMultiCircuitBreaker(cfg CircuitBreakerConfig) *MultiCircuitBreaker {
	return &MultiCircuitBreaker{
		config: cfg,
	}
}

// Get returns the circuit breaker for a specific key
func (mcb *MultiCircuitBreaker) Get(key string) *CircuitBreaker {
	if v, ok := mcb.breakers.Load(key); ok {
		return v.(*CircuitBreaker)
	}

	cb := NewCircuitBreaker(mcb.config)
	actual, _ := mcb.breakers.LoadOrStore(key, cb)
	return actual.(*CircuitBreaker)
}

// Reset resets all circuit breakers
func (mcb *MultiCircuitBreaker) Reset() {
	mcb.breakers.Range(func(key, value any) bool {
		value.(*CircuitBreaker).Reset()
		return true
	})
}

// States returns the state of all circuit breakers
func (mcb *MultiCircuitBreaker) States() map[string]CircuitState {
	states := make(map[string]CircuitState)
	mcb.breakers.Range(func(key, value any) bool {
		states[key.(string)] = value.(*CircuitBreaker).State()
		return true
	})
	return states
}

// MultiCircuitBreakerMiddleware creates middleware that uses per-subject circuit breakers
func MultiCircuitBreakerMiddleware(mcb *MultiCircuitBreaker) PublishMiddleware {
	return func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, subject string, msg *core.Message) error {
			cb := mcb.Get(subject)

			if err := cb.Allow(); err != nil {
				return err
			}

			err := next(ctx, subject, msg)
			if err != nil {
				cb.RecordFailure(err)
			} else {
				cb.RecordSuccess()
			}

			return err
		}
	}
}
