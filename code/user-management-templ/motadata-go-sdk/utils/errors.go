package utils

import (
	"errors"
	"fmt"
)

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

// L1 cache errors

var (
	ErrKeyNotFound            = errors.New("key not found")
	ErrValueNotFound          = errors.New("value not found")
	ErrTenantChecksumNotMatch = errors.New("checksum not matching")
)

// OTEL errors - errors related to OpenTelemetry logger, tracer, and metrics components

var (
	// ErrOTELProviderNotInitialized indicates OTEL provider was not initialized before use
	ErrOTELProviderNotInitialized = errors.New("OTEL provider not initialized")

	// ErrOTELProviderAlreadyClosed indicates an attempt to use a closed OTEL provider
	ErrOTELProviderAlreadyClosed = errors.New("OTEL provider already closed")

	// ErrOTELUnsupportedProtocol indicates an unsupported OTEL exporter protocol was specified
	ErrOTELUnsupportedProtocol = errors.New("unsupported OTEL protocol - use 'grpc' or 'http'")

	// ErrOTELResourceCreationFailed indicates failure to create OTEL resource
	ErrOTELResourceCreationFailed = errors.New("failed to create OTEL resource")

	// ErrOTELExporterCreationFailed indicates failure to create OTEL exporter
	ErrOTELExporterCreationFailed = errors.New("failed to create OTEL exporter")

	// ErrOTELInvalidLogLevel indicates an invalid log level string was provided
	ErrOTELInvalidLogLevel = errors.New("invalid log level")

	// ErrOTELMetricCreationFailed indicates failure to create a metric instrument
	ErrOTELMetricCreationFailed = errors.New("failed to create metric instrument")

	// ErrOTELShutdownFailed indicates one or more components failed during shutdown
	ErrOTELShutdownFailed = errors.New("OTEL shutdown encountered errors")
)

// codec errors

var (
	ErrUnpackFailed        = errors.New("unpack failed: invalid or corrupted data")
	ErrUnsupportedCodec    = errors.New("unsupported codec type")
	ErrUnsupportedDataType = errors.New("unsupported data type")
	ErrValueOutOfRange     = errors.New("value out of supported range")
	ErrDataTooLarge        = errors.New("data too large for codec length fields")
)

// WrapErr wraps a sentinel error with a cause for consistent error formatting.
func WrapErr(sentinel, cause error) error {
	return fmt.Errorf("%w: %v", sentinel, cause)
}
