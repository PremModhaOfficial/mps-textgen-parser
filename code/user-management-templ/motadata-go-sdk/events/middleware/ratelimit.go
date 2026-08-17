package middleware

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nats-io/nats.go"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	rate       float64 // Tokens per second
	burst      int64   // Maximum bucket size
	tokens     atomic.Int64
	lastUpdate atomic.Int64
}

// RateLimiterConfig configures the rate limiter
type RateLimiterConfig struct {
	// Rate is the number of requests allowed per second
	Rate float64

	// Burst is the maximum number of requests that can be made at once
	Burst int
}

// DefaultRateLimiterConfig returns default configuration
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		Rate:  100, // 100 requests per second
		Burst: 50,  // Allow bursts of 50
	}
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(cfg RateLimiterConfig) *RateLimiter {
	if cfg.Rate <= 0 {
		cfg.Rate = 100
	}
	if cfg.Burst <= 0 {
		cfg.Burst = int(cfg.Rate)
	}

	rl := &RateLimiter{
		rate:  cfg.Rate,
		burst: int64(cfg.Burst),
	}
	rl.tokens.Store(int64(cfg.Burst * 1000)) // Store as milliTokens
	rl.lastUpdate.Store(time.Now().UnixNano())
	return rl
}

// Allow checks if a request is allowed
func (rl *RateLimiter) Allow() bool {
	return rl.AllowN(1)
}

// AllowN checks if n requests are allowed
func (rl *RateLimiter) AllowN(n int) bool {
	rl.refill()

	needed := int64(n * 1000) // Convert to milliTokens
	for {
		current := rl.tokens.Load()
		if current < needed {
			return false
		}
		if rl.tokens.CompareAndSwap(current, current-needed) {
			return true
		}
	}
}

// Wait blocks until a request is allowed or context is cancelled
func (rl *RateLimiter) Wait(ctx context.Context) error {
	return rl.WaitN(ctx, 1)
}

// WaitN blocks until n requests are allowed or context is cancelled
func (rl *RateLimiter) WaitN(ctx context.Context, n int) error {
	for {
		if rl.AllowN(n) {
			return nil
		}

		// Calculate wait time
		needed := int64(n * 1000)
		current := rl.tokens.Load()
		deficit := needed - current
		waitTime := time.Duration(float64(deficit) / rl.rate * float64(time.Millisecond))

		if waitTime < time.Millisecond {
			waitTime = time.Millisecond
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitTime):
			// Try again
		}
	}
}

// refill adds tokens based on elapsed time
func (rl *RateLimiter) refill() {
	now := time.Now().UnixNano()
	last := rl.lastUpdate.Load()

	elapsed := float64(now-last) / float64(time.Second)
	if elapsed <= 0 {
		return
	}

	// Calculate new tokens
	newTokens := int64(elapsed * rl.rate * 1000) // milliTokens
	if newTokens <= 0 {
		return
	}

	// Try to update atomically
	for {
		current := rl.tokens.Load()
		updated := current + newTokens
		max := rl.burst * 1000
		if updated > max {
			updated = max
		}

		if rl.tokens.CompareAndSwap(current, updated) {
			rl.lastUpdate.Store(now)
			return
		}
	}
}

// Tokens returns the current number of available tokens
func (rl *RateLimiter) Tokens() float64 {
	rl.refill()
	return float64(rl.tokens.Load()) / 1000
}

// ErrRateLimitExceeded is returned when rate limit is exceeded
var ErrRateLimitExceeded = &utils.Error{
	Kind: "rate_limiter",
	Op:   "allow",
	Err:  utils.ErrPublishFailed,
}

// RateLimitMiddleware returns publish middleware that applies rate limiting
func RateLimitMiddleware(rl *RateLimiter) PublishMiddleware {
	return func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, msg *nats.Msg) error {
			if !rl.Allow() {
				return ErrRateLimitExceeded
			}
			return next(ctx, msg)
		}
	}
}

// RateLimitWaitMiddleware returns middleware that waits for rate limit
func RateLimitWaitMiddleware(rl *RateLimiter) PublishMiddleware {
	return func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, msg *nats.Msg) error {
			if err := rl.Wait(ctx); err != nil {
				return err
			}
			return next(ctx, msg)
		}
	}
}

// RateLimitSubscribeMiddleware returns subscribe middleware with rate limiting
func RateLimitSubscribeMiddleware(rl *RateLimiter) SubscribeMiddleware {
	return func(next SubscribeHandler) SubscribeHandler {
		return func(ctx context.Context, msg *nats.Msg) error {
			if !rl.Allow() {
				return ErrRateLimitExceeded
			}
			return next(ctx, msg)
		}
	}
}

// PerSubjectRateLimiter manages rate limiters per subject
type PerSubjectRateLimiter struct {
	config   RateLimiterConfig
	limiters sync.Map // map[string]*RateLimiter
}

// NewPerSubjectRateLimiter creates a per-subject rate limiter
func NewPerSubjectRateLimiter(cfg RateLimiterConfig) *PerSubjectRateLimiter {
	return &PerSubjectRateLimiter{
		config: cfg,
	}
}

// Get returns the rate limiter for a specific subject
func (psr *PerSubjectRateLimiter) Get(subject string) *RateLimiter {
	if v, ok := psr.limiters.Load(subject); ok {
		return v.(*RateLimiter)
	}

	rl := NewRateLimiter(psr.config)
	actual, _ := psr.limiters.LoadOrStore(subject, rl)
	return actual.(*RateLimiter)
}

// PerSubjectRateLimitMiddleware creates middleware with per-subject rate limiting
func PerSubjectRateLimitMiddleware(psr *PerSubjectRateLimiter) PublishMiddleware {
	return func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, msg *nats.Msg) error {
			rl := psr.Get(msg.Subject)
			if !rl.Allow() {
				return ErrRateLimitExceeded
			}
			return next(ctx, msg)
		}
	}
}

// SlidingWindowLimiter implements a sliding window rate limiter
type SlidingWindowLimiter struct {
	window   time.Duration
	limit    int64
	mu       sync.Mutex
	requests []int64 // Unix nano timestamps
}

// NewSlidingWindowLimiter creates a sliding window limiter
func NewSlidingWindowLimiter(limit int, window time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		window:   window,
		limit:    int64(limit),
		requests: make([]int64, 0, limit),
	}
}

// Allow checks if a request is allowed
func (sw *SlidingWindowLimiter) Allow() bool {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now().UnixNano()
	windowStart := now - int64(sw.window)

	// Remove old requests
	validIdx := 0
	for i, ts := range sw.requests {
		if ts > windowStart {
			validIdx = i
			break
		}
		validIdx = i + 1
	}
	if validIdx > 0 && validIdx <= len(sw.requests) {
		sw.requests = sw.requests[validIdx:]
	}

	// Check limit
	if int64(len(sw.requests)) >= sw.limit {
		return false
	}

	// Add new request
	sw.requests = append(sw.requests, now)
	return true
}

// Count returns the number of requests in the current window
func (sw *SlidingWindowLimiter) Count() int {
	sw.mu.Lock()
	defer sw.mu.Unlock()

	now := time.Now().UnixNano()
	windowStart := now - int64(sw.window)

	count := 0
	for _, ts := range sw.requests {
		if ts > windowStart {
			count++
		}
	}
	return count
}

// SlidingWindowMiddleware returns middleware using sliding window limiter
func SlidingWindowMiddleware(sw *SlidingWindowLimiter) PublishMiddleware {
	return func(next PublishHandler) PublishHandler {
		return func(ctx context.Context, msg *nats.Msg) error {
			if !sw.Allow() {
				return ErrRateLimitExceeded
			}
			return next(ctx, msg)
		}
	}
}
