package utils

import (
	"errors"
	"fmt"
	"strings"
)

// ============================================================================
// Error Constructors
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

// NewError creates a new error with the given message.
func NewError(message string) error {
	return errors.New(message)
}

// NewErrorf creates a new error with a formatted message.
func NewErrorf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

// ============================================================================
// Error Checking
// ============================================================================

// Is reports whether any error in err's chain matches target.
// This is a convenient wrapper around errors.Is.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target.
// This is a convenient wrapper around errors.As.
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
// Sentinel Errors for Utils Package
// ============================================================================

var (
	// ErrInvalidArgument indicates an invalid argument was provided
	ErrInvalidArgument = errors.New("invalid argument")

	// ErrNilValue indicates a nil value was provided where non-nil was expected
	ErrNilValue = errors.New("nil value")

	// ErrEmptyValue indicates an empty value was provided
	ErrEmptyValue = errors.New("empty value")

	// ErrOutOfRange indicates a value is out of valid range
	ErrOutOfRange = errors.New("out of range")

	// ErrNotFound indicates a requested item was not found
	ErrNotFound = errors.New("not found")

	// ErrAlreadyExists indicates the item already exists
	ErrAlreadyExists = errors.New("already exists")

	// ErrOperationCanceled indicates the operation was canceled
	ErrOperationCanceled = errors.New("operation canceled")

	// ErrTimeout indicates the operation timed out
	ErrTimeout = errors.New("timeout")

	// ErrClosed indicates the resource is closed
	ErrClosed = errors.New("closed")
)

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
// Error Formatting
// ============================================================================

// FormatError formats an error with additional context.
func FormatError(operation string, err error) string {
	if err == nil {
		return fmt.Sprintf("%s: success", operation)
	}
	return fmt.Sprintf("%s: %v", operation, err)
}

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
// Panic Recovery
// ============================================================================

// RecoverError recovers from a panic and returns it as an error.
// Use in defer: defer func() { err = RecoverError(recover()) }()
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
// If the function panics, the error is sent to the error channel.
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
