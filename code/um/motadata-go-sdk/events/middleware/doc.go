// Package middleware provides cross-cutting concerns for the events system.
//
// This package implements middleware for:
//   - Retry logic with exponential backoff
//   - Distributed tracing propagation
//   - Metrics collection
//   - Structured logging
//   - Circuit breaker pattern
//   - Rate limiting
//
// # Middleware Pattern
//
// Middleware wraps publish and subscribe operations:
//
//	type PublishMiddleware func(next PublishFunc) PublishFunc
//	type SubscribeMiddleware func(next core.MessageHandler) core.MessageHandler
//
// # Middleware Stack
//
// Use the middleware stack to compose multiple middlewares:
//
//	stack := middleware.NewStack()
//
//	// Add publish middleware (applied in order)
//	stack.AddPublish(
//	    middleware.TracingPublish(),
//	    middleware.MetricsMiddleware(metricsConfig),
//	    middleware.RetryPublish(retryConfig),
//	)
//
//	// Add subscribe middleware
//	stack.AddSubscribe(
//	    middleware.TracingSubscribe(),
//	    middleware.LoggingSubscribe(logConfig),
//	)
//
// # Retry Middleware
//
// Automatically retry failed publish operations:
//
//	retryConfig := middleware.RetryConfig{
//	    MaxAttempts:     3,
//	    InitialInterval: 100 * time.Millisecond,
//	    MaxInterval:     5 * time.Second,
//	    Multiplier:      2.0,
//	    Jitter:          0.1,
//	}
//
//	middleware := middleware.Retry(retryConfig)
//
// # Tracing Middleware
//
// Propagate distributed tracing context:
//
//	// For publish operations
//	tracingPublish := middleware.TracingPublish()
//
//	// For subscribe operations
//	tracingSubscribe := middleware.TracingSubscribe()
//
// The middleware extracts trace context from the Go context and adds it to
// message headers, supporting both W3C Trace Context and B3 formats.
//
// # Metrics Middleware
//
// Collect publishing metrics:
//
//	metricsConfig := middleware.MetricsConfig{
//	    Namespace: "myapp",
//	    Subsystem: "events",
//	    Buckets:   []float64{.001, .005, .01, .05, .1, .5, 1},
//	}
//
//	metricsMiddleware := middleware.MetricsMiddleware(metricsConfig)
//
// Collected metrics:
//   - Publish count (success/failure)
//   - Publish latency histogram
//   - Message size histogram
//
// # Logging Middleware
//
// Add structured logging to operations:
//
//	loggingConfig := middleware.LoggingConfig{
//	    LogPayload:    false,  // Don't log message data
//	    LogHeaders:    true,   // Log message headers
//	    SlowThreshold: 100 * time.Millisecond,
//	}
//
//	loggingMiddleware := middleware.Logging(loggingConfig)
//
// # Circuit Breaker
//
// Prevent cascading failures with the circuit breaker pattern:
//
//	cb := middleware.NewCircuitBreaker(middleware.CircuitBreakerConfig{
//	    FailureThreshold:   5,                    // Open after 5 failures
//	    SuccessThreshold:   2,                    // Close after 2 successes
//	    Timeout:            30 * time.Second,     // Stay open for 30s
//	    HalfOpenMaxRequests: 3,                   // Allow 3 requests in half-open
//	})
//
//	// Use as middleware
//	cbMiddleware := middleware.CircuitBreakerMiddleware(cb)
//
//	// Check circuit state
//	if cb.State() == middleware.CircuitOpen {
//	    log.Println("Circuit is open, requests will be rejected")
//	}
//
// Circuit states:
//
//	middleware.CircuitClosed   - Normal operation
//	middleware.CircuitOpen     - Failing, rejecting requests
//	middleware.CircuitHalfOpen - Testing if service recovered
//
// # Rate Limiting
//
// Control message throughput:
//
//	// Basic rate limiter
//	limiter := middleware.NewRateLimiter(middleware.RateLimiterConfig{
//	    Rate:       1000,           // 1000 messages per second
//	    BurstSize:  100,            // Allow burst of 100
//	})
//
//	// Reject if rate exceeded
//	rateLimitMiddleware := middleware.RateLimitMiddleware(limiter)
//
//	// Wait for available capacity (blocking)
//	waitMiddleware := middleware.RateLimitWaitMiddleware(limiter)
//
//	// Per-subject rate limiting
//	perSubject := middleware.NewPerSubjectRateLimiter(100, 10) // 100/s, burst 10
//	perSubjectMiddleware := middleware.PerSubjectRateLimitMiddleware(perSubject)
//
// # Multi-Circuit Breaker
//
// Separate circuit breakers per subject:
//
//	mcb := middleware.NewMultiCircuitBreaker(middleware.CircuitBreakerConfig{
//	    FailureThreshold: 5,
//	    Timeout:          30 * time.Second,
//	})
//
//	// Get circuit breaker for a specific subject
//	cb := mcb.GetBreaker("orders.created")
//
// # Error Handling
//
// Middleware errors:
//
//	middleware.ErrCircuitOpen       - Circuit breaker is open
//	middleware.ErrRateLimitExceeded - Rate limit exceeded
//
// Example:
//
//	err := publisher.Publish(ctx, subject, msg)
//	if errors.Is(err, middleware.ErrCircuitOpen) {
//	    log.Println("Service unavailable, circuit open")
//	}
package middleware
