package utils

import "errors"

// pool errors
var (
	ErrPoolClosed               = errors.New("pool is closed")
	ErrPoolExhausted            = errors.New("pool is exhausted")
	ErrInvalidSize              = errors.New("pool size must be greater than 0")
	ErrFactoryFunctionRequired  = errors.New("factory function is required")
	ErrWorkerPoolNotInitialized = errors.New("pool not initialized - call Init() first")
	ErrQueueFull                = errors.New("task queue is full - increase MaxQueueSize or reduce task submission rate")
)

// circuit breaker errors
var (
	ErrTooManyRequest = errors.New("too many requests in half-open state")
	ErrCircuitOpen    = errors.New("circuit breaker is open")
)
