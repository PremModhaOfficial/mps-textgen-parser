package circuitbreaker

import (
	"errors"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"
	"github.com/sony/gobreaker/v2"
)

// Package circuitbreaker provides a circuit breaker implementation for fault tolerance.
//
// A circuit breaker is a design pattern used to detect failures and prevent cascading
// failures in distributed systems. It wraps calls to external services and monitors
// for failures. When failures exceed a threshold, the circuit "opens" and subsequent
// calls fail immediately without attempting the operation, giving the failing service
// time to recover.
//
// Circuit Breaker States:
//   - Closed: Normal operation. Requests pass through and failures are counted.
//   - Open: Circuit is tripped. Requests fail immediately with ErrCircuitOpen.
//   - Half-Open: Testing recovery. Limited requests are allowed through to test
//     if the service has recovered.
//
// Key Features:
//   - Configurable failure threshold to trip the circuit
//   - Configurable timeout for automatic recovery attempts
//   - Configurable interval to reset failure counts in closed state
//   - Custom success evaluation for error handling
//   - State change callbacks for monitoring
//   - Thread-safe for concurrent use
//
// Usage:
//
//	// Create with default configuration
//	cb := circuitbreaker.NewCircuitBreaker(circuitbreaker.DefaultConfig("my-service"))
//
//	// Execute a function through the circuit breaker
//	result, err := cb.Execute(func() (any, error) {
//	    return callExternalService()
//	})
//
//	if errors.Is(err, circuitbreaker.ErrCircuitOpen) {
//	    // Circuit is open, handle gracefully
//	}
//
// Custom Configuration:
//
//	cfg := circuitbreaker.Config{
//	    Name:             "payment-service",
//	    FailureThreshold: 5,              // Open after 5 consecutive failures
//	    Timeout:          30 * time.Second, // Try recovery after 30s
//	    MaxRequests:      3,              // Allow 3 requests in half-open state
//	    OnStateChange: func(name string, from, to circuitbreaker.State) {
//	        log.Printf("Circuit %s: %s -> %s", name, from, to)
//	    },
//	}
//	cb := circuitbreaker.NewCircuitBreaker(cfg)

// State represents the circuit breaker state.
// The circuit breaker transitions between these states based on the success
// or failure of executed operations.
type State int

const (
	StateClosed   State = iota // Normal operation, requests pass through
	StateHalfOpen              // Testing if service recovered
	StateOpen                  // Circuit is open, requests fail fast
)

// String returns the string representation of the state
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateHalfOpen:
		return "half-open"
	case StateOpen:
		return "open"
	default:
		return "unknown"
	}
}

// Config holds the circuit breaker configuration
type Config struct {
	// Name identifies the circuit breaker
	Name string

	// MaxRequests is the maximum number of requests allowed in half-open state
	// Default: 1
	MaxRequests uint32

	// Interval is the cyclic period of the closed state to clear the internal counts
	// If 0, the internal counts are never cleared
	// Default: 0
	Interval time.Duration

	// Timeout is the period of the open state, after which the state becomes half-open
	// Default: 60 seconds
	Timeout time.Duration

	// FailureThreshold is the number of consecutive failures needed to open the circuit
	// Default: 5
	FailureThreshold uint32

	// SuccessThreshold is the number of consecutive successes needed to close the circuit from half-open
	// Default: 1
	SuccessThreshold uint32

	// OnStateChange is called when the circuit breaker state changes
	OnStateChange func(name string, from, to State)

	// IsSuccessful determines if the error should be counted as a failure
	// Return true if the error should NOT be counted as failure
	// Default: nil (all errors are failures)
	IsSuccessful func(err error) bool
}

// DefaultConfig returns a Config with sensible defaults
func DefaultConfig(name string) Config {
	return Config{
		Name:             name,
		MaxRequests:      1,
		Interval:         0,
		Timeout:          60 * time.Second,
		FailureThreshold: 5,
		SuccessThreshold: 1,
	}
}

// CircuitBreaker wraps sony/gobreaker for simplified usage
type CircuitBreaker struct {
	breaker *gobreaker.CircuitBreaker[any]
	config  Config
}

// NewCircuitBreaker creates a new CircuitBreaker with the given configuration
func NewCircuitBreaker(cfg Config) *CircuitBreaker {
	// Apply defaults for zero values
	if cfg.MaxRequests == 0 {
		cfg.MaxRequests = 1
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 60 * time.Second
	}
	if cfg.FailureThreshold == 0 {
		cfg.FailureThreshold = 5
	}
	if cfg.SuccessThreshold == 0 {
		cfg.SuccessThreshold = 1
	}

	settings := gobreaker.Settings{
		Name:        cfg.Name,
		MaxRequests: cfg.MaxRequests,
		Interval:    cfg.Interval,
		Timeout:     cfg.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= cfg.FailureThreshold
		},
		IsSuccessful: cfg.IsSuccessful,
	}

	if cfg.OnStateChange != nil {
		settings.OnStateChange = func(name string, from, to gobreaker.State) {
			cfg.OnStateChange(name, convertState(from), convertState(to))
		}
	}

	return &CircuitBreaker{
		breaker: gobreaker.NewCircuitBreaker[any](settings),
		config:  cfg,
	}
}

// Execute runs the given function through the circuit breaker
func (cb *CircuitBreaker) Execute(fn func() (any, error)) (any, error) {
	result, err := cb.breaker.Execute(fn)
	return result, convertError(err)
}

// State returns the current state of the circuit breaker
func (cb *CircuitBreaker) State() State {
	return convertState(cb.breaker.State())
}

// Name returns the name of the circuit breaker
func (cb *CircuitBreaker) Name() string {
	return cb.breaker.Name()
}

// Counts returns the current counts of the circuit breaker
func (cb *CircuitBreaker) Counts() Counts {
	counts := cb.breaker.Counts()
	return Counts{
		Requests:             counts.Requests,
		TotalSuccesses:       counts.TotalSuccesses,
		TotalFailures:        counts.TotalFailures,
		ConsecutiveSuccesses: counts.ConsecutiveSuccesses,
		ConsecutiveFailures:  counts.ConsecutiveFailures,
	}
}

// IsOpen returns true if the circuit is open
func (cb *CircuitBreaker) IsOpen() bool {
	return cb.State() == StateOpen
}

// IsClosed returns true if the circuit is closed
func (cb *CircuitBreaker) IsClosed() bool {
	return cb.State() == StateClosed
}

// Counts holds the statistics of a circuit breaker
type Counts struct {
	Requests             uint32
	TotalSuccesses       uint32
	TotalFailures        uint32
	ConsecutiveSuccesses uint32
	ConsecutiveFailures  uint32
}

// convertState converts gobreaker.State to our State type
func convertState(state gobreaker.State) State {
	switch state {
	case gobreaker.StateClosed:
		return StateClosed
	case gobreaker.StateHalfOpen:
		return StateHalfOpen
	case gobreaker.StateOpen:
		return StateOpen
	default:
		return StateClosed
	}
}

// convertError converts gobreaker errors to our error types
func convertError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gobreaker.ErrOpenState) {
		return utils.ErrCircuitOpen
	}
	if errors.Is(err, gobreaker.ErrTooManyRequests) {
		return utils.ErrTooManyRequest
	}
	return err
}
