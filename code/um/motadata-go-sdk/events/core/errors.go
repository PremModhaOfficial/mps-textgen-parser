package core

import (
	"errors"
	"fmt"
)

// Base errors
var (
	// Connection errors
	ErrNotConnected      = errors.New("not connected")
	ErrAlreadyConnected  = errors.New("already connected")
	ErrConnectionClosed  = errors.New("connection closed")
	ErrConnectionTimeout = errors.New("connection timeout")
	ErrReconnectFailed   = errors.New("reconnect failed")

	// Publishing errors
	ErrPublishFailed   = errors.New("publish failed")
	ErrPublishTimeout  = errors.New("publish timeout")
	ErrNoAck           = errors.New("no acknowledgment received")
	ErrDuplicateMsg    = errors.New("duplicate message")
	ErrStreamNotFound  = errors.New("stream not found")
	ErrInvalidSubject  = errors.New("invalid subject")
	ErrInvalidMessage  = errors.New("invalid message")
	ErrMessageTooLarge = errors.New("message too large")
	ErrRequestTimeout  = errors.New("request timeout")
	ErrNoReply         = errors.New("no reply received")

	// Subscription errors
	ErrSubscriptionClosed  = errors.New("subscription closed")
	ErrSubscriptionInvalid = errors.New("subscription invalid")
	ErrMaxMsgsExceeded     = errors.New("max messages exceeded")

	// Authentication errors
	ErrAuthFailed        = errors.New("authentication failed")
	ErrAuthExpired       = errors.New("authentication expired")
	ErrInvalidCredential = errors.New("invalid credentials")
	ErrPermissionDenied  = errors.New("permission denied")

	// Tenant errors
	ErrTenantNotFound     = errors.New("tenant not found")
	ErrTenantExists       = errors.New("tenant already exists")
	ErrTenantDisconnected = errors.New("tenant disconnected")

	// Configuration errors
	ErrInvalidConfig = errors.New("invalid configuration")
	ErrMissingConfig = errors.New("missing required configuration")

	// Serialization errors
	ErrSerializationFailed   = errors.New("serialization failed")
	ErrDeserializationFailed = errors.New("deserialization failed")

	// Lifecycle errors
	ErrShutdownInProgress = errors.New("shutdown in progress")
	ErrAlreadyStarted     = errors.New("already started")
	ErrNotStarted         = errors.New("not started")
)

// Error wraps an error with additional context
type Error struct {
	Op      string // Operation that failed
	Kind    string // Category of error
	Err     error  // Underlying error
	Details map[string]any
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Op != "" {
		return fmt.Sprintf("%s: %s: %v", e.Kind, e.Op, e.Err)
	}
	return fmt.Sprintf("%s: %v", e.Kind, e.Err)
}

// Unwrap returns the underlying error
func (e *Error) Unwrap() error {
	return e.Err
}

// NewError creates a new Error
func NewError(kind, op string, err error) *Error {
	return &Error{
		Op:   op,
		Kind: kind,
		Err:  err,
	}
}

// WithDetails adds details to the error
func (e *Error) WithDetails(details map[string]any) *Error {
	e.Details = details
	return e
}

// TenantError represents an error specific to a tenant
type TenantError struct {
	TenantID string
	Op       string
	Err      error
}

// Error implements the error interface
func (e TenantError) Error() string {
	if e.Op != "" {
		return fmt.Sprintf("tenant %s: %s: %v", e.TenantID, e.Op, e.Err)
	}
	return fmt.Sprintf("tenant %s: %v", e.TenantID, e.Err)
}

// Unwrap returns the underlying error
func (e TenantError) Unwrap() error {
	return e.Err
}

// NewTenantError creates a new TenantError
func NewTenantError(tenantID, op string, err error) TenantError {
	return TenantError{
		TenantID: tenantID,
		Op:       op,
		Err:      err,
	}
}

// ConfigError represents a configuration error
type ConfigError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (e ConfigError) Error() string {
	return fmt.Sprintf("config error: %s: %s", e.Field, e.Message)
}

// SerializationError represents a serialization/deserialization error
type SerializationError struct {
	Operation string // "serialize" or "deserialize"
	Type      string // Type being processed
	Err       error
}

// Error implements the error interface
func (e SerializationError) Error() string {
	return fmt.Sprintf("%s %s: %v", e.Operation, e.Type, e.Err)
}

// Unwrap returns the underlying error
func (e SerializationError) Unwrap() error {
	return e.Err
}

// IsRetryable returns true if the error is potentially retryable
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	// Non-retryable errors
	switch {
	case errors.Is(err, ErrInvalidSubject),
		errors.Is(err, ErrInvalidMessage),
		errors.Is(err, ErrInvalidConfig),
		errors.Is(err, ErrMissingConfig),
		errors.Is(err, ErrSerializationFailed),
		errors.Is(err, ErrPermissionDenied),
		errors.Is(err, ErrDuplicateMsg),
		errors.Is(err, ErrShutdownInProgress),
		errors.Is(err, ErrConnectionClosed):
		return false
	}

	// Check for SerializationError
	var serErr SerializationError
	if errors.As(err, &serErr) {
		return false
	}

	// Check for ConfigError
	var cfgErr ConfigError
	if errors.As(err, &cfgErr) {
		return false
	}

	// Everything else is potentially retryable
	return true
}

// IsTemporary returns true if the error is temporary
func IsTemporary(err error) bool {

	if err == nil {

		return false
	}

	switch {
	case errors.Is(err, ErrConnectionTimeout),
		errors.Is(err, ErrPublishTimeout),
		errors.Is(err, ErrReconnectFailed),
		errors.Is(err, ErrNoAck):

		return true
	}

	return false
}

// MultiError represents multiple errors that occurred during an operation
type MultiError struct {
	Errors []error
}

// Error implements the error interface
func (multiError *MultiError) Error() string {

	if len(multiError.Errors) == 0 {

		return "no errors"
	}

	if len(multiError.Errors) == 1 {

		return multiError.Errors[0].Error()
	}

	return fmt.Sprintf("%d errors occurred: %v", len(multiError.Errors), multiError.Errors)
}

// Unwrap returns the first error for errors.Is/As support
func (multiError *MultiError) Unwrap() error {

	if len(multiError.Errors) > 0 {

		return multiError.Errors[0]
	}

	return nil
}

// All returns all the errors
func (multiError *MultiError) All() []error {

	return multiError.Errors
}

// NewMultiError creates a new MultiError from a slice of errors
func NewMultiError(errs []error) *MultiError {

	// Filter out nil errors
	var filtered []error

	for _, err := range errs {

		if err != nil {

			filtered = append(filtered, err)
		}
	}

	if len(filtered) == 0 {

		return nil
	}

	return &MultiError{Errors: filtered}
}
