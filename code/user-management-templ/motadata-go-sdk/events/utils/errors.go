// Package utils provides shared error types, sentinel errors, and error
// classification utilities for the events subsystem.
//
// Errors are organized into categories (connection, publish, subscribe, etc.)
// so callers can use errors.Is to match specific failure modes and decide
// whether to retry, reconnect, or surface the error to the user.
package utils

import (
	"errors"
	"fmt"
	"strings"
)

// ============================================================================
// Sentinel Errors — Connection
//
// These errors indicate failures in the NATS connection lifecycle. They are
// returned by connection managers and health-check routines.
// ============================================================================

var (
	// ErrNotConnected is returned when an operation requires an active
	// connection but no connection has been established yet.
	ErrNotConnected = errors.New("not connected")

	// ErrAlreadyConnected is returned when a caller tries to connect to
	// a broker that is already connected, preventing duplicate connections.
	ErrAlreadyConnected = errors.New("already connected")

	// ErrConnectionClosed is returned when an operation is attempted on a
	// connection that has been explicitly closed. This is non-retryable
	// because the closure was intentional.
	ErrConnectionClosed = errors.New("connection closed")

	// ErrConnectionTimeout is returned when a connection attempt exceeds the
	// configured timeout duration. This is typically temporary and retryable.
	ErrConnectionTimeout = errors.New("connection timeout")

	// ErrReconnectFailed is returned when all automatic reconnection attempts
	// have been exhausted without successfully restoring the connection.
	ErrReconnectFailed = errors.New("reconnect failed")
)

// ============================================================================
// Sentinel Errors — Publishing
//
// These errors occur when publishing messages to NATS subjects or JetStream
// streams. Some are retryable (timeouts) while others indicate permanent
// problems (invalid subject, duplicate message).
// ============================================================================

var (
	// ErrPublishFailed is returned when a publish operation fails for a
	// reason not covered by a more specific sentinel error.
	ErrPublishFailed = errors.New("publish failed")

	// ErrPublishTimeout is returned when a publish operation does not
	// complete within its configured deadline. This is temporary and retryable.
	ErrPublishTimeout = errors.New("publish timeout")

	// ErrNoAck is returned when a JetStream publish does not receive an
	// acknowledgment from the server, indicating the message may not have
	// been persisted. This is temporary and retryable.
	ErrNoAck = errors.New("no acknowledgment received")

	// ErrDuplicateMsg is returned by JetStream when a message with the same
	// deduplication ID has already been accepted. This is non-retryable
	// because the message was already processed successfully.
	ErrDuplicateMsg = errors.New("duplicate message")

	// ErrStreamNotFound is returned when a publish targets a JetStream stream
	// that does not exist on the server.
	ErrStreamNotFound = errors.New("stream not found")

	// ErrInvalidSubject is returned when a subject name fails validation
	// (empty, too long, or contains disallowed characters). Non-retryable.
	ErrInvalidSubject = errors.New("invalid subject")

	// ErrInvalidMessage is returned when a message fails structural
	// validation (nil payload, missing required headers, etc.). Non-retryable.
	ErrInvalidMessage = errors.New("invalid message")

	// ErrMessageTooLarge is returned when a message payload exceeds the
	// configured maximum size (DefaultMaxPayload). Non-retryable without
	// reducing the payload.
	ErrMessageTooLarge = errors.New("message too large")

	// ErrRequestTimeout is returned when a request-reply operation does not
	// receive a response within the configured deadline. Temporary and retryable.
	ErrRequestTimeout = errors.New("request timeout")

	// ErrNoReply is returned when a request-reply operation receives no
	// response from any subscriber, indicating no service is listening.
	ErrNoReply = errors.New("no reply received")
)

// ============================================================================
// Sentinel Errors — Subscription
//
// These errors relate to NATS subscriptions and consumer operations. They are
// typically surfaced by subscriber callbacks or consumer drain routines.
// ============================================================================

var (
	// ErrSubscriptionClosed is returned when a message delivery is attempted
	// on a subscription that has already been unsubscribed or drained.
	ErrSubscriptionClosed = errors.New("subscription closed")

	// ErrSubscriptionInvalid is returned when a subscription handle is
	// malformed or was never properly initialized.
	ErrSubscriptionInvalid = errors.New("subscription invalid")

	// ErrMaxMsgsExceeded is returned when a subscription has received more
	// messages than its configured maximum, causing automatic unsubscribe.
	ErrMaxMsgsExceeded = errors.New("max messages exceeded")
)

// ============================================================================
// Sentinel Errors — Authentication
//
// These errors are returned when NATS authentication or authorization fails.
// Callers should generally avoid retrying auth errors unless credentials have
// been refreshed (e.g., after re-reading a rotated NKey or JWT).
// ============================================================================

var (
	// ErrAuthFailed is returned when the NATS server rejects the supplied
	// credentials during the initial handshake.
	ErrAuthFailed = errors.New("authentication failed")

	// ErrAuthExpired is returned when a previously valid credential (JWT or
	// token) has expired and can no longer be used.
	ErrAuthExpired = errors.New("authentication expired")

	// ErrInvalidCredential is returned when the credential format is wrong
	// or the credential cannot be parsed (e.g., corrupt NKey seed).
	ErrInvalidCredential = errors.New("invalid credentials")

	// ErrPermissionDenied is returned when the authenticated user lacks the
	// required publish or subscribe permission for the target subject.
	// Non-retryable unless server-side permissions are updated.
	ErrPermissionDenied = errors.New("permission denied")
)

// ============================================================================
// Sentinel Errors — Tenant
//
// These errors are specific to multi-tenant connection management. Each tenant
// maintains an isolated NATS connection, and these errors indicate problems
// with tenant lifecycle operations (registration, lookup, disconnect).
// ============================================================================

var (
	// ErrTenantNotFound is returned when a requested tenant ID does not
	// exist in the tenant registry.
	ErrTenantNotFound = errors.New("tenant not found")

	// ErrTenantExists is returned when attempting to register a tenant
	// whose ID is already present in the registry.
	ErrTenantExists = errors.New("tenant already exists")

	// ErrTenantDisconnected is returned when an operation targets a tenant
	// whose NATS connection is no longer active.
	ErrTenantDisconnected = errors.New("tenant disconnected")
)

// ============================================================================
// Sentinel Errors — Configuration
//
// These errors are returned during configuration loading and validation.
// They are always non-retryable; the caller must fix the configuration
// before retrying.
// ============================================================================

var (
	// ErrInvalidConfig is returned when a configuration value fails
	// validation (out of range, wrong type, conflicting settings, etc.).
	ErrInvalidConfig = errors.New("invalid configuration")

	// ErrMissingConfig is returned when a required configuration field
	// (e.g., server URL) is not provided and has no usable default.
	ErrMissingConfig = errors.New("missing required configuration")
)

// ============================================================================
// Sentinel Errors — Serialization
//
// These errors occur during message encoding/decoding. They are non-retryable
// because the same input will produce the same failure.
// ============================================================================

var (
	// ErrSerializationFailed is returned when a message payload cannot be
	// encoded to the target format (JSON, Protobuf, MsgPack, etc.).
	ErrSerializationFailed = errors.New("serialization failed")

	// ErrDeserializationFailed is returned when raw bytes cannot be decoded
	// into the expected Go type, typically due to schema mismatch or
	// corrupted data.
	ErrDeserializationFailed = errors.New("deserialization failed")
)

// ============================================================================
// Sentinel Errors — Lifecycle
//
// These errors relate to the start/stop lifecycle of event system components
// (connections, microservices, publishers, subscribers). They prevent
// double-start, double-stop, and operations on uninitialized components.
// ============================================================================

var (
	// ErrShutdownInProgress is returned when an operation is attempted while
	// a graceful shutdown is already underway. Non-retryable.
	ErrShutdownInProgress = errors.New("shutdown in progress")

	// ErrAlreadyStarted is returned when Start() is called on a component
	// that is already running, preventing duplicate goroutines or listeners.
	ErrAlreadyStarted = errors.New("already started")

	// ErrNotStarted is returned when an operation requires a component to
	// be running but Start() has not been called yet.
	ErrNotStarted = errors.New("not started")

	// ErrJetStreamNotEnabled is returned when a JetStream operation (stream
	// management, KV, object store) is attempted on a connection that was
	// established without JetStream support.
	ErrJetStreamNotEnabled = errors.New("jetstream not enabled")
)

// ============================================================================
// Sentinel Errors — General Utility
//
// Generic sentinel errors that can be used across all subsystems when a
// more specific category does not apply.
// ============================================================================

var (
	// ErrInvalidArgument is returned when a function receives a parameter
	// that does not meet its preconditions (wrong type, out of range, etc.).
	ErrInvalidArgument = errors.New("invalid argument")

	// ErrNilValue is returned when a required pointer, interface, or slice
	// argument is nil.
	ErrNilValue = errors.New("nil value")

	// ErrEmptyValue is returned when a required string or collection
	// argument is empty.
	ErrEmptyValue = errors.New("empty value")

	// ErrOutOfRange is returned when a numeric argument (index, size, count)
	// falls outside the accepted bounds.
	ErrOutOfRange = errors.New("out of range")

	// ErrNotFound is a generic "not found" error for lookups that do not
	// match a more specific sentinel (e.g., ErrTenantNotFound, ErrStreamNotFound).
	ErrNotFound = errors.New("not found")

	// ErrAlreadyExists is a generic "already exists" error for creation
	// operations that detect a duplicate.
	ErrAlreadyExists = errors.New("already exists")

	// ErrOperationCanceled is returned when an operation is aborted because
	// its context was canceled or its deadline was exceeded.
	ErrOperationCanceled = errors.New("operation canceled")

	// ErrTimeout is a generic timeout error for operations that do not have
	// a more specific timeout sentinel (e.g., ErrConnectionTimeout).
	ErrTimeout = errors.New("timeout")

	// ErrClosed is a generic "closed" error for operations attempted on a
	// resource that has been shut down.
	ErrClosed = errors.New("closed")
)

// ============================================================================
// Structured Error Types
// ============================================================================

// Error is a structured error that enriches a root cause with the operation
// name, an error category (kind), and optional key-value details. It
// implements both the error and errors.Unwrap interfaces so callers can
// use errors.Is / errors.As to inspect the underlying cause.
//
// Example:
//
//	NewError("publish", "Publish", ErrPublishTimeout).WithDetails(map[string]any{"subject": "orders.new"})
type Error struct {
	Op      string         // Op is the name of the operation that failed (e.g., "Connect", "Publish").
	Kind    string         // Kind is the error category (e.g., "connection", "publish", "auth").
	Err     error          // Err is the underlying root cause.
	Details map[string]any // Details holds optional structured metadata about the failure.
}

// Error returns a human-readable representation of the structured error,
// including the kind, operation (if set), and the underlying cause.
func (e *Error) Error() string {
	if e.Op != "" {
		return fmt.Sprintf("%s: %s: %v", e.Kind, e.Op, e.Err)
	}
	return fmt.Sprintf("%s: %v", e.Kind, e.Err)
}

// Unwrap returns the underlying error so that errors.Is and errors.As
// can traverse the error chain to match sentinel errors.
func (e *Error) Unwrap() error {
	return e.Err
}

// NewError creates a new structured Error with the given category (kind),
// operation name, and root cause. Use WithDetails to attach additional context.
func NewError(kind, op string, err error) *Error {
	return &Error{
		Op:   op,
		Kind: kind,
		Err:  err,
	}
}

// WithDetails attaches arbitrary key-value metadata to the error for
// structured logging or diagnostics. Returns the same pointer to allow
// method chaining.
//
// NOTE: This mutates the receiver. In practice this is safe because
// WithDetails is called immediately after NewError in a builder pattern.
func (e *Error) WithDetails(details map[string]any) *Error {
	e.Details = details
	return e
}

// TenantError represents a failure scoped to a specific tenant. It includes
// the tenant identifier so log aggregation and alerting systems can correlate
// failures per tenant. Implements error and errors.Unwrap interfaces.
type TenantError struct {
	TenantID string // TenantID identifies the tenant that experienced the failure.
	Op       string // Op is the operation that failed (e.g., "Connect", "Subscribe").
	Err      error  // Err is the underlying root cause.
}

// Error returns a human-readable string prefixed with the tenant ID
// and operation, making it easy to filter logs by tenant.
func (e TenantError) Error() string {
	if e.Op != "" {
		return fmt.Sprintf("tenant %s: %s: %v", e.TenantID, e.Op, e.Err)
	}
	return fmt.Sprintf("tenant %s: %v", e.TenantID, e.Err)
}

// Unwrap returns the underlying error to support errors.Is / errors.As
// chain traversal through the tenant error wrapper.
func (e TenantError) Unwrap() error {
	return e.Err
}

// NewTenantError creates a TenantError binding the failure to a specific
// tenant, operation, and root cause.
func NewTenantError(tenantID, op string, err error) TenantError {
	return TenantError{
		TenantID: tenantID,
		Op:       op,
		Err:      err,
	}
}

// ConfigError represents a configuration validation failure tied to a
// specific configuration field. It is returned during config loading and
// is always non-retryable (see IsRetryable).
type ConfigError struct {
	Field   string // Field is the configuration key that failed validation.
	Message string // Message describes why the value is invalid.
}

// Error returns a formatted string identifying the invalid config field
// and the reason for rejection.
func (e ConfigError) Error() string {
	return fmt.Sprintf("config error: %s: %s", e.Field, e.Message)
}

// SerializationError captures details about a codec failure during message
// encoding or decoding. It records the direction (serialize/deserialize),
// the target Go type, and the underlying codec error. Non-retryable
// because the same payload will produce the same failure.
type SerializationError struct {
	Operation string // Operation is either "serialize" or "deserialize".
	Type      string // Type is the Go type name being encoded/decoded.
	Err       error  // Err is the underlying codec error.
}

// Error returns a human-readable description including the operation
// direction, the affected type, and the root cause.
func (e SerializationError) Error() string {
	return fmt.Sprintf("%s %s: %v", e.Operation, e.Type, e.Err)
}

// Unwrap returns the underlying codec error to enable errors.Is / errors.As
// chain traversal.
func (e SerializationError) Unwrap() error {
	return e.Err
}

// MultiError aggregates multiple errors that occurred during a single
// logical operation (e.g., shutting down several subscribers). It
// implements the standard error interface and supports errors.Is / errors.As
// matching against the first collected error via Unwrap.
type MultiError struct {
	Errors []error // Errors is the list of non-nil errors collected during the operation.
}

// Error returns a summary string. When there is exactly one error it
// delegates to that error's message; otherwise it reports the count and
// lists all errors.
func (me *MultiError) Error() string {
	if len(me.Errors) == 0 {
		return "no errors"
	}
	if len(me.Errors) == 1 {
		return me.Errors[0].Error()
	}
	return fmt.Sprintf("%d errors occurred: %v", len(me.Errors), me.Errors)
}

// Unwrap returns the first error so that errors.Is and errors.As can
// match against the most significant failure in the collection.
func (me *MultiError) Unwrap() error {
	if len(me.Errors) > 0 {
		return me.Errors[0]
	}
	return nil
}

// All returns the complete slice of collected errors for callers that
// need to inspect every failure individually.
func (me *MultiError) All() []error {
	return me.Errors
}

// NewMultiError creates a MultiError from a slice of errors, automatically
// filtering out nil entries. Returns nil when no non-nil errors remain,
// so callers can safely use: if err := NewMultiError(errs); err != nil { ... }
func NewMultiError(errs []error) *MultiError {
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

// ============================================================================
// Error Classification
//
// Classification functions help callers decide how to react to errors
// without hard-coding sentinel checks throughout the codebase.
// ============================================================================

// IsRetryable returns true if the error is potentially recoverable by
// retrying the same operation. Errors that are deterministic failures
// (invalid input, config errors, serialization errors, permission denied,
// duplicate messages, intentional closure, shutdown) return false because
// retrying would produce the same result.
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

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

	var serErr SerializationError
	if errors.As(err, &serErr) {
		return false
	}

	var cfgErr ConfigError
	if errors.As(err, &cfgErr) {
		return false
	}

	return true
}

// IsTemporary returns true if the error represents a transient condition
// that is expected to resolve on its own (timeouts, temporary network
// issues). Unlike IsRetryable, which is permissive (defaults to true),
// IsTemporary is conservative and only returns true for explicitly
// recognized transient conditions.
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

// ============================================================================
// Validation Errors
// ============================================================================

// ValidationError represents a validation failure.
type ValidationError struct {
	Field   string
	Message string
	Value   any
}

// Error implements the error interface.
func (e ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error: %s: %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// NewValidationError creates a new ValidationError.
func NewValidationError(field, message string) ValidationError {
	return ValidationError{Field: field, Message: message}
}

// NewValidationErrorWithValue creates a new ValidationError with the invalid value.
func NewValidationErrorWithValue(field, message string, value any) ValidationError {
	return ValidationError{Field: field, Message: message, Value: value}
}

// ============================================================================
// Error Wrapping & Formatting
// ============================================================================

// WrapError wraps an error with additional context.
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// WrapErrorf wraps an error with formatted context.
func WrapErrorf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(format+": %w", append(args, err)...)
}

// FormatError formats an error with additional context.
func FormatError(operation string, err error) string {
	if err == nil {
		return fmt.Sprintf("%s: success", operation)
	}
	return fmt.Sprintf("%s: %v", operation, err)
}

// ============================================================================
// Error Checking
// ============================================================================

// Is reports whether any error in err's chain matches target.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target.
func As[T error](err error, target *T) bool {
	return errors.As(err, target)
}

// IsAnyOf checks if the error matches any of the target errors.
func IsAnyOf(err error, targets ...error) bool {
	for _, target := range targets {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

// ============================================================================
// Error Chain Utilities
// ============================================================================

// ErrorChain returns all errors in the error chain.
func ErrorChain(err error) []error {
	if err == nil {
		return nil
	}

	var chain []error
	for err != nil {
		chain = append(chain, err)
		err = errors.Unwrap(err)
	}
	return chain
}

// RootCause returns the root cause of an error (the last error in the chain).
func RootCause(err error) error {
	if err == nil {
		return nil
	}

	for {
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			return err
		}
		err = unwrapped
	}
}

// ============================================================================
// Error Collection
// ============================================================================

// ErrorCollector collects multiple errors and combines them.
type ErrorCollector struct {
	errors []error
}

// NewErrorCollector creates a new ErrorCollector.
func NewErrorCollector() *ErrorCollector {
	return &ErrorCollector{}
}

// Add adds an error to the collector (ignores nil errors).
func (ec *ErrorCollector) Add(err error) {
	if err != nil {
		ec.errors = append(ec.errors, err)
	}
}

// AddAll adds multiple errors to the collector.
func (ec *ErrorCollector) AddAll(errs ...error) {
	for _, err := range errs {
		ec.Add(err)
	}
}

// HasErrors returns true if any errors were collected.
func (ec *ErrorCollector) HasErrors() bool {
	return len(ec.errors) > 0
}

// Count returns the number of collected errors.
func (ec *ErrorCollector) Count() int {
	return len(ec.errors)
}

// All returns all collected errors.
func (ec *ErrorCollector) All() []error {
	return ec.errors
}

// Error returns a combined error, or nil if no errors were collected.
func (ec *ErrorCollector) Error() error {
	switch len(ec.errors) {
	case 0:
		return nil
	case 1:
		return ec.errors[0]
	default:
		return &combinedError{errors: ec.errors}
	}
}

// combinedError represents multiple errors combined.
type combinedError struct {
	errors []error
}

func (ce *combinedError) Error() string {
	if len(ce.errors) == 0 {
		return "no errors"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d errors occurred:\n", len(ce.errors)))
	for i, err := range ce.errors {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Sprintf("  [%d] %s", i+1, err.Error()))
	}
	return sb.String()
}

func (ce *combinedError) Unwrap() []error {
	return ce.errors
}

// ============================================================================
// Panic Recovery
// ============================================================================

// RecoverError recovers from a panic and returns it as an error.
func RecoverError(r any) error {
	if r == nil {
		return nil
	}

	switch v := r.(type) {
	case error:
		return v
	case string:
		return errors.New(v)
	default:
		return fmt.Errorf("panic: %v", v)
	}
}

// SafeGo runs a function in a goroutine with panic recovery.
func SafeGo(fn func() error, errChan chan<- error) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				if errChan != nil {
					select {
					case errChan <- RecoverError(r):
					default:
					}
				}
			}
		}()

		if err := fn(); err != nil && errChan != nil {
			select {
			case errChan <- err:
			default:
			}
		}
	}()
}
