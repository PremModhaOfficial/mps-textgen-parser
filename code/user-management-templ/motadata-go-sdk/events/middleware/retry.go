package middleware

import (
	"context"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"
)

// RetryConfig holds the parameters for exponential backoff with jitter.
//
// The backoff duration for attempt N is computed as:
//
//	base = InitialInterval * (Multiplier ^ N)
//	capped = min(base, MaxInterval)
//	jitter = capped * Jitter * random(-1, 1)
//	backoff = max(capped + jitter, InitialInterval)
//
// This produces an exponentially growing delay that is capped at MaxInterval,
// with a random jitter component to decorrelate retries from multiple clients
// (preventing thundering herd problems).
type RetryConfig struct {
	// MaxAttempts is the maximum number of retry attempts after the initial call.
	// For example, MaxAttempts=3 means up to 4 total executions (1 initial + 3 retries).
	MaxAttempts int

	// InitialInterval is the base delay before the first retry. Subsequent retries
	// grow this value exponentially by the Multiplier.
	InitialInterval time.Duration

	// MaxInterval caps the backoff duration. Once the exponential growth exceeds
	// this value, all subsequent retries use MaxInterval as their base delay.
	MaxInterval time.Duration

	// Multiplier is the exponential growth factor applied to InitialInterval on
	// each successive retry. A value of 2.0 doubles the delay each time.
	Multiplier float64

	// Jitter is the fraction of the current backoff to randomize, preventing
	// synchronized retries across clients. A value of 0.1 means +/-10% randomization.
	// Set to 0 to disable jitter.
	Jitter float64
}

// DefaultRetryConfig returns default retry settings
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxAttempts:     3,
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     5 * time.Second,
		Multiplier:      2.0,
		Jitter:          0.1,
	}
}

// RetryMiddleware adds retry logic to publish operations
type RetryMiddleware struct {
	config      RetryConfig
	shouldRetry func(error) bool
	rng         *rand.Rand
	rngMu       sync.Mutex
}

// NewRetryMiddleware creates a new retry middleware
func NewRetryMiddleware(config RetryConfig) *RetryMiddleware {
	return &RetryMiddleware{
		config:      config,
		shouldRetry: DefaultShouldRetry,
		rng:         rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// WithShouldRetry sets a custom retry check function
func (r *RetryMiddleware) WithShouldRetry(fn func(error) bool) *RetryMiddleware {
	r.shouldRetry = fn
	return r
}

// InterceptPublish returns the publish middleware
func (r *RetryMiddleware) InterceptPublish() PublishMiddleware {
	return func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, msg *nats.Msg) error {
			return r.executeWithRetry(ctx, msg, next)
		}
	}
}

// executeWithRetry runs the handler with retry logic and exponential backoff.
func (r *RetryMiddleware) executeWithRetry(ctx context.Context, msg *nats.Msg, next PublishHandler) error {
	var lastErr error

	for attempt := 0; attempt <= r.config.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return firstNonNil(lastErr, err)
		}

		if err := next(ctx, msg); err == nil {
			return nil
		} else {
			lastErr = err
		}

		if !r.shouldRetry(lastErr) || attempt >= r.config.MaxAttempts {
			return lastErr
		}

		if err := r.waitBackoff(ctx, attempt); err != nil {
			return lastErr
		}
	}

	return lastErr
}

// waitBackoff waits for the calculated backoff duration or context cancellation.
func (r *RetryMiddleware) waitBackoff(ctx context.Context, attempt int) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(r.calculateBackoff(attempt)):
		return nil
	}
}

// firstNonNil returns the first non-nil error.
func firstNonNil(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// InterceptSubscribe returns nil (retry not applicable for subscribe)
func (r *RetryMiddleware) InterceptSubscribe() SubscribeMiddleware {
	return func(next SubscribeHandler) SubscribeHandler {
		return next // Pass through
	}
}

// calculateBackoff returns the backoff duration for the given attempt
func (r *RetryMiddleware) calculateBackoff(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}

	// Calculate base backoff with exponential growth
	backoff := float64(r.config.InitialInterval) * math.Pow(r.config.Multiplier, float64(attempt))

	// Cap at maximum interval
	maxBackoff := float64(r.config.MaxInterval)
	if backoff > maxBackoff {
		backoff = maxBackoff
	}

	// Apply jitter (mutex-protected for concurrent safety)
	if r.config.Jitter > 0 {
		jitterRange := backoff * r.config.Jitter
		r.rngMu.Lock()
		jitter := jitterRange * (r.rng.Float64()*2 - 1)
		r.rngMu.Unlock()
		backoff += jitter

		// Ensure backoff doesn't go below minimum
		if backoff < float64(r.config.InitialInterval) {
			backoff = float64(r.config.InitialInterval)
		}
	}

	return time.Duration(backoff)
}

// DefaultShouldRetry is the default retry check function
func DefaultShouldRetry(err error) bool {
	return err != nil && utils.IsRetryable(err)
}

// Retry returns a retry middleware with default config
func Retry() *RetryMiddleware {
	return NewRetryMiddleware(DefaultRetryConfig())
}

// RetryWithConfig returns a retry middleware with custom config
func RetryWithConfig(config RetryConfig) *RetryMiddleware {
	return NewRetryMiddleware(config)
}
