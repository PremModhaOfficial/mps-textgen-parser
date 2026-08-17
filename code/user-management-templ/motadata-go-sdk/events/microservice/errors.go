package microservice

import (
	"errors"
	"fmt"
)

// Base errors for microservice operations
var (
	// Service errors
	ErrInvalidService    = errors.New("invalid service")
	ErrMissingServiceID  = errors.New("missing service ID")
	ErrMissingServiceName = errors.New("missing service name")
	ErrServiceExists     = errors.New("service already exists")
	ErrServiceNotFound   = errors.New("service not found")

	// Instance errors
	ErrInvalidInstance   = errors.New("invalid instance")
	ErrMissingInstanceID = errors.New("missing instance ID")
	ErrMissingHost       = errors.New("missing host")
	ErrInstanceExists    = errors.New("instance already exists")
	ErrInstanceNotFound  = errors.New("instance not found")
	ErrInstanceExpired   = errors.New("instance has expired")

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

// ServiceError represents an error specific to a service
type ServiceError struct {
	ServiceID   string
	ServiceName string
	Op          string
	Err         error
}

// Error implements the error interface
func (e ServiceError) Error() string {
	if e.ServiceName != "" {
		return fmt.Sprintf("service %s (%s): %s: %v", e.ServiceName, e.ServiceID, e.Op, e.Err)
	}
	if e.Op != "" {
		return fmt.Sprintf("service %s: %s: %v", e.ServiceID, e.Op, e.Err)
	}
	return fmt.Sprintf("service %s: %v", e.ServiceID, e.Err)
}

// Unwrap returns the underlying error
func (e ServiceError) Unwrap() error {
	return e.Err
}

// NewServiceError creates a new ServiceError
func NewServiceError(serviceID, serviceName, op string, err error) ServiceError {
	return ServiceError{
		ServiceID:   serviceID,
		ServiceName: serviceName,
		Op:          op,
		Err:         err,
	}
}

// InstanceError represents an error specific to an instance
type InstanceError struct {
	InstanceID string
	ServiceID  string
	Op         string
	Err        error
}

// Error implements the error interface
func (e InstanceError) Error() string {
	if e.ServiceID != "" {
		return fmt.Sprintf("instance %s (service %s): %s: %v", e.InstanceID, e.ServiceID, e.Op, e.Err)
	}
	if e.Op != "" {
		return fmt.Sprintf("instance %s: %s: %v", e.InstanceID, e.Op, e.Err)
	}
	return fmt.Sprintf("instance %s: %v", e.InstanceID, e.Err)
}

// Unwrap returns the underlying error
func (e InstanceError) Unwrap() error {
	return e.Err
}

// NewInstanceError creates a new InstanceError
func NewInstanceError(instanceID, serviceID, op string, err error) InstanceError {
	return InstanceError{
		InstanceID: instanceID,
		ServiceID:  serviceID,
		Op:         op,
		Err:        err,
	}
}

// RegistrationError represents an error during registration
type RegistrationError struct {
	ServiceID  string
	InstanceID string
	Reason     string
	Err        error
}

// Error implements the error interface
func (e RegistrationError) Error() string {
	id := e.ServiceID
	if e.InstanceID != "" {
		id = e.InstanceID + "@" + e.ServiceID
	}
	if e.Reason != "" {
		return fmt.Sprintf("registration failed for %s: %s: %v", id, e.Reason, e.Err)
	}
	return fmt.Sprintf("registration failed for %s: %v", id, e.Err)
}

// Unwrap returns the underlying error
func (e RegistrationError) Unwrap() error {
	return e.Err
}

// HealthCheckError represents an error during health checking
type HealthCheckError struct {
	ServiceID        string
	InstanceID       string
	URL              string
	Err              error
	ConsecutiveFails int
}

// Error implements the error interface
func (e HealthCheckError) Error() string {
	id := e.ServiceID
	if e.InstanceID != "" {
		id = e.InstanceID + "@" + e.ServiceID
	}
	return fmt.Sprintf("health check failed for %s at %s (consecutive fails: %d): %v",
		id, e.URL, e.ConsecutiveFails, e.Err)
}

// Unwrap returns the underlying error
func (e HealthCheckError) Unwrap() error {
	return e.Err
}

// IsServiceError checks if the error is a service-related error
func IsServiceError(err error) bool {
	var svcErr ServiceError
	var instErr InstanceError
	var regErr RegistrationError
	var healthErr HealthCheckError

	return errors.As(err, &svcErr) ||
		errors.As(err, &instErr) ||
		errors.As(err, &regErr) ||
		errors.As(err, &healthErr)
}

// IsNotFound checks if the error indicates a service or instance was not found
func IsNotFound(err error) bool {
	return errors.Is(err, ErrServiceNotFound) || errors.Is(err, ErrInstanceNotFound)
}

// IsExists checks if the error indicates a service or instance already exists
func IsExists(err error) bool {
	return errors.Is(err, ErrServiceExists) || errors.Is(err, ErrInstanceExists)
}
