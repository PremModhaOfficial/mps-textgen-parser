package middleware

import (
	"context"
	"math"
	"math/rand"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxAttempts     int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
	Jitter          float64
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
		return func(ctx context.Context, subject string, msg *core.Message) error {
			var lastErr error

			for attempt := 0; attempt <= r.config.MaxAttempts; attempt++ {
				// Check context before each attempt
				select {
				case <-ctx.Done():
					if lastErr != nil {
						return lastErr
					}
					return ctx.Err()
				default:
				}

				// Execute the handler
				err := next(ctx, subject, msg)
				if err == nil {
					return nil // Success
				}

				lastErr = err

				// Check if we should retry
				if !r.shouldRetry(err) {
					return err // Non-retryable error
				}

				// Check if we have more attempts
				if attempt >= r.config.MaxAttempts {
					return err // No more retries
				}

				// Calculate backoff and wait
				backoff := r.calculateBackoff(attempt)

				select {
				case <-ctx.Done():
					return lastErr
				case <-time.After(backoff):
					// Continue to next attempt
				}
			}

			return lastErr
		}
	}
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

	// Apply jitter
	if r.config.Jitter > 0 {
		jitterRange := backoff * r.config.Jitter
		jitter := jitterRange * (r.rng.Float64()*2 - 1)
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
	if err == nil {
		return false
	}

	// Don't retry on these errors
	if core.IsRetryable(err) {
		return true
	}

	return false
}

// Retry returns a retry middleware with default config
func Retry() *RetryMiddleware {
	return NewRetryMiddleware(DefaultRetryConfig())
}

// RetryWithConfig returns a retry middleware with custom config
func RetryWithConfig(config RetryConfig) *RetryMiddleware {
	return NewRetryMiddleware(config)
}
