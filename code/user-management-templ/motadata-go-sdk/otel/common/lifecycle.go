package common

import (
	"fmt"
	"strings"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"
)

/* ---------------------------------------- Shutdown Error Aggregation -------------------------------------------------- */

// AggregateErrors combines multiple errors into a single error.
// This is useful for collecting errors during graceful shutdown of multiple
// OTEL components where we want to report all failures rather than just the first.
//
// Aggregation rules:
//   - Returns nil if all errors are nil or the slice is empty
//   - Returns the single error directly if only one non-nil error exists
//   - Combines multiple errors into a formatted error message
//
// Parameters:
//   - errors: Slice of errors to aggregate (may contain nil values)
//
// Returns:
//   - error: The aggregated error, or nil if no errors
//
// Example:
//
//	errors := []error{nil, fmt.Errorf("tracer shutdown failed"), nil}
//	err := AggregateErrors(errors)  // Returns "tracer shutdown failed"
//
//	errors = []error{
//	    fmt.Errorf("tracer failed"),
//	    fmt.Errorf("metrics failed"),
//	}
//	err := AggregateErrors(errors)  // Returns "multiple errors: [tracer failed; metrics failed]"
func AggregateErrors(errors []error) error {

	// Filter out nil errors to get only actual failures
	// This allows callers to append errors unconditionally
	var nonNilErrors []error

	for _, err := range errors {

		if err != nil {
			nonNilErrors = append(nonNilErrors, err)
		}
	}

	// No errors - return nil to indicate success
	if len(nonNilErrors) == 0 {

		return nil
	}

	// Single error - return it directly without wrapping
	// This preserves error type information for errors.Is() checks
	if len(nonNilErrors) == 1 {

		return nonNilErrors[0]
	}

	// Multiple errors - combine into a formatted message
	// Format: "multiple errors: [error1; error2; error3]"
	var errorMessages []string

	for _, err := range nonNilErrors {
		errorMessages = append(errorMessages, err.Error())
	}

	return fmt.Errorf("%w: [%s]", utils.ErrOTELShutdownFailed, strings.Join(errorMessages, "; "))
}

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// ShutdownCollector collects shutdown errors from multiple components.
// This provides a convenient way to aggregate errors during graceful shutdown
// of OTEL components (logger, tracer, metrics).
//
// The collector is designed for the common shutdown pattern where multiple
// components need to be shut down and all errors should be reported.
//
// Usage:
//
//	collector := NewShutdownCollector()
//	collector.Collect(tracer.Shutdown(ctx))
//	collector.Collect(metrics.Shutdown(ctx))
//	collector.Collect(logger.Close())
//	return collector.Error()
//
// Thread Safety:
//
//	This type is NOT thread-safe. It should only be used from a single goroutine.
//	For concurrent shutdown, use channels or synchronization.
type ShutdownCollector struct {
	// errors holds the collected non-nil errors
	// Pre-allocated with empty slice to avoid nil checks
	errors []error
}

/* ---------------------------------------- ShutdownCollector Constructor -------------------------------------------------- */

// NewShutdownCollector creates a new shutdown error collector.
// The collector is initialized with an empty error slice ready for use.
//
// Returns:
//   - *ShutdownCollector: A new collector instance
//
// Example:
//
//	collector := NewShutdownCollector()
//	defer func() {
//	    if err := collector.Error(); err != nil {
//	        log.Printf("shutdown errors: %v", err)
//	    }
//	}()
func NewShutdownCollector() *ShutdownCollector {

	return &ShutdownCollector{
		// Pre-allocate empty slice to avoid nil map issues
		// Most shutdowns have 0-3 errors, so small initial capacity is fine
		errors: make([]error, 0),
	}
}

/* ---------------------------------------- ShutdownCollector Methods -------------------------------------------------- */

// Collect adds an error to the collector if it's not nil.
// This method is safe to call with nil errors, making it convenient
// to use with functions that may or may not return errors.
//
// Parameters:
//   - err: The error to collect (nil errors are ignored)
//
// Example:
//
//	collector.Collect(provider.Shutdown(ctx))  // Only collected if non-nil
//	collector.Collect(nil)                     // Safely ignored
func (collector *ShutdownCollector) Collect(err error) {

	// Only append non-nil errors to avoid bloating the slice
	// This allows callers to collect errors unconditionally
	if err != nil {
		collector.errors = append(collector.errors, err)
	}
}

// CollectFunc executes a shutdown function and collects any error.
// This is useful when you have a function that needs to be called
// and its error collected in a single operation.
//
// Parameters:
//   - shutdownFunc: The function to execute (can be nil)
//
// Example:
//
//	collector.CollectFunc(func() error {
//	    return provider.Shutdown(ctx)
//	})
func (collector *ShutdownCollector) CollectFunc(shutdownFunc func() error) {

	// Guard against nil function to prevent panic
	if shutdownFunc != nil {
		collector.Collect(shutdownFunc())
	}
}

// Error returns the aggregated error or nil if no errors were collected.
// This method should be called after all shutdown operations are complete.
//
// Returns:
//   - error: The aggregated error, or nil if no errors were collected
//
// Example:
//
//	if err := collector.Error(); err != nil {
//	    return fmt.Errorf("graceful shutdown failed: %w", err)
//	}
func (collector *ShutdownCollector) Error() error {

	// Delegate to AggregateErrors for consistent error formatting
	return AggregateErrors(collector.errors)
}

// HasErrors returns true if any errors were collected.
// This is useful for conditional logic without needing to format the error.
//
// Returns:
//   - bool: True if at least one error was collected
//
// Example:
//
//	if collector.HasErrors() {
//	    log.Warn("shutdown completed with errors")
//	}
func (collector *ShutdownCollector) HasErrors() bool {

	return len(collector.errors) > 0
}

// Count returns the number of collected errors.
// This can be useful for logging or metrics about shutdown failures.
//
// Returns:
//   - int: The number of errors collected
//
// Example:
//
//	fmt.Printf("Shutdown completed with %d errors\n", collector.Count())
func (collector *ShutdownCollector) Count() int {

	return len(collector.errors)
}
