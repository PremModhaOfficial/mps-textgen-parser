package endpoint

import (
	"errors"
	"fmt"
)

// Base errors for endpoint operations
var (
	// Registration errors
	ErrInvalidEndpoint   = errors.New("invalid endpoint")
	ErrMissingEndpointID = errors.New("missing endpoint ID")
	ErrMissingServiceID  = errors.New("missing service ID")
	ErrMissingHost       = errors.New("missing host")
	ErrEndpointExists    = errors.New("endpoint already exists")
	ErrEndpointNotFound  = errors.New("endpoint not found")
	ErrEndpointExpired   = errors.New("endpoint has expired")

	// Health check errors
	ErrHealthCheckFailed  = errors.New("health check failed")
	ErrHealthCheckTimeout = errors.New("health check timeout")

	// Operation errors
	ErrRegistryShutdown = errors.New("registry is shutting down")
	ErrInvalidQuery     = errors.New("invalid query parameters")

	// Manager errors
	ErrManagerNotStarted = errors.New("manager not started")
	ErrManagerShutdown   = errors.New("manager is shutting down")
)

// EndpointError represents an error specific to an endpoint
type EndpointError struct {
	EndpointID string
	ServiceID  string
	Op         string
	Err        error
}

// Error implements the error interface
func (e EndpointError) Error() string {
	if e.ServiceID != "" {
		return fmt.Sprintf("endpoint %s (service %s): %s: %v", e.EndpointID, e.ServiceID, e.Op, e.Err)
	}
	if e.Op != "" {
		return fmt.Sprintf("endpoint %s: %s: %v", e.EndpointID, e.Op, e.Err)
	}
	return fmt.Sprintf("endpoint %s: %v", e.EndpointID, e.Err)
}

// Unwrap returns the underlying error
func (e EndpointError) Unwrap() error {
	return e.Err
}

// NewEndpointError creates a new EndpointError
func NewEndpointError(endpointID, serviceID, op string, err error) EndpointError {
	return EndpointError{
		EndpointID: endpointID,
		ServiceID:  serviceID,
		Op:         op,
		Err:        err,
	}
}

// RegistrationError represents an error during endpoint registration
type RegistrationError struct {
	EndpointID string
	Reason     string
	Err        error
}

// Error implements the error interface
func (e RegistrationError) Error() string {
	if e.Reason != "" {
		return fmt.Sprintf("registration failed for endpoint %s: %s: %v", e.EndpointID, e.Reason, e.Err)
	}
	return fmt.Sprintf("registration failed for endpoint %s: %v", e.EndpointID, e.Err)
}

// Unwrap returns the underlying error
func (e RegistrationError) Unwrap() error {
	return e.Err
}

// HealthCheckError represents an error during health checking
type HealthCheckError struct {
	EndpointID       string
	URL              string
	Err              error
	ConsecutiveFails int
}

// Error implements the error interface
func (e HealthCheckError) Error() string {
	return fmt.Sprintf("health check failed for endpoint %s at %s (consecutive fails: %d): %v",
		e.EndpointID, e.URL, e.ConsecutiveFails, e.Err)
}

// Unwrap returns the underlying error
func (e HealthCheckError) Unwrap() error {
	return e.Err
}

// IsEndpointError checks if the error is an endpoint-related error
func IsEndpointError(err error) bool {
	var endpointErr EndpointError
	var regErr RegistrationError
	var healthErr HealthCheckError

	return errors.As(err, &endpointErr) ||
		errors.As(err, &regErr) ||
		errors.As(err, &healthErr)
}

// IsNotFound checks if the error indicates an endpoint was not found
func IsNotFound(err error) bool {
	return errors.Is(err, ErrEndpointNotFound)
}

// IsExists checks if the error indicates an endpoint already exists
func IsExists(err error) bool {
	return errors.Is(err, ErrEndpointExists)
}
