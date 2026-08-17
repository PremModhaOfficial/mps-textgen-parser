package common

import (
	"fmt"
	"strings"
)

/* ---------------------------------------- Shutdown error Aggregation -------------------------------------------------- */

// AggregateErrors combines multiple errors into a single error.
// Returns nil if all errors are nil or the slice is empty.
func AggregateErrors(errors []error) error {

	var nonNilErrors []error

	for _, err := range errors {

		if err != nil {
			nonNilErrors = append(nonNilErrors, err)
		}
	}

	if len(nonNilErrors) == 0 {

		return nil
	}

	if len(nonNilErrors) == 1 {

		return nonNilErrors[0]
	}

	var errorMessages []string

	for _, err := range nonNilErrors {
		errorMessages = append(errorMessages, err.Error())
	}

	return fmt.Errorf("multiple errors: [%s]", strings.Join(errorMessages, "; "))
}

/* ---------------------------------------- Type Definitions -------------------------------------------------- */

// ShutdownCollector collects shutdown errors from multiple components.
type ShutdownCollector struct {
	errors []error
}

/* ---------------------------------------- ShutdownCollector Constructor -------------------------------------------------- */

// NewShutdownCollector creates a new shutdown error collector.
func NewShutdownCollector() *ShutdownCollector {

	return &ShutdownCollector{
		errors: make([]error, 0),
	}
}

/* ---------------------------------------- ShutdownCollector Methods -------------------------------------------------- */

// Collect adds an error to the collector if it's not nil.
func (collector *ShutdownCollector) Collect(err error) {

	if err != nil {
		collector.errors = append(collector.errors, err)
	}
}

// CollectFunc executes a shutdown function and collects any error.
func (collector *ShutdownCollector) CollectFunc(shutdownFunc func() error) {

	if shutdownFunc != nil {
		collector.Collect(shutdownFunc())
	}
}

// error returns the aggregated error or nil if no errors were collected.
func (collector *ShutdownCollector) Error() error {

	return AggregateErrors(collector.errors)
}

// HasErrors returns true if any errors were collected.
func (collector *ShutdownCollector) HasErrors() bool {

	return len(collector.errors) > 0
}

// Count returns the number of collected errors.
func (collector *ShutdownCollector) Count() int {

	return len(collector.errors)
}
