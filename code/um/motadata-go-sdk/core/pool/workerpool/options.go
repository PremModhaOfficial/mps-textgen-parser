// Package workerpool provides a high-performance worker Pool implementation
// with auto-scaling capabilities, metrics tracking, and advanced task scheduling.
package workerpool

import (
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"
)

// TaskOptions defines configuration options for task execution including
// timeout, retry policies, and callback handlers.
type TaskOptions struct {
	// Name is the optional identifier for the task (useful for debugging)
	Name string
	// Timeout specifies the maximum execution time for the task
	Timeout time.Duration
	// Retries defines the number of retry attempts on task failure
	Retries int
	// RetryDelay specifies the delay between retry attempts
	RetryDelay time.Duration
	// OnSuccess callback executed when task completes successfully
	OnSuccess func()
	// OnError callback executed when task fails (after all retries)
	OnError func(error2 error)
	// OnRetry callback executed before each retry attempt
	OnRetry func(attempt int, err error)
}

// TaskOption is a functional option for configuring TaskOptions.
type TaskOption func(*TaskOptions)

// defaultOptions returns TaskOptions with sensible default values.
func defaultOptions() *TaskOptions {
	return &TaskOptions{
		Name:       utils.Empty,
		Timeout:    30 * time.Second,
		Retries:    0,
		RetryDelay: time.Second,
	}
}

// WithName sets a custom name for the task (useful for debugging and monitoring).
func WithName(name string) TaskOption {
	return func(taskOptions *TaskOptions) { taskOptions.Name = name }
}

// WithTimeout sets the maximum execution time for the task.
// Tasks exceeding this duration will be cancelled.
func WithTimeout(duration time.Duration) TaskOption {
	return func(taskOptions *TaskOptions) { taskOptions.Timeout = duration }
}

// WithRetry configures retry behavior for task failures.
// retries: number of retry attempts (0 = no retries)
// delay: base delay between retries (actual delay may use exponential backoff)
func WithRetry(retries int, delay time.Duration) TaskOption {
	return func(taskOptions *TaskOptions) {
		taskOptions.Retries = retries
		taskOptions.RetryDelay = delay
	}
}

// WithOnSuccess sets a callback function executed when the task completes successfully.
func WithOnSuccess(fn func()) TaskOption {
	return func(taskOptions *TaskOptions) { taskOptions.OnSuccess = fn }
}

// WithOnError sets a callback function executed when the task fails after all retries.
func WithOnError(fn func(error)) TaskOption {
	return func(taskOptions *TaskOptions) { taskOptions.OnError = fn }
}

// WithOnRetry sets a callback function executed before each retry attempt.
// The callback receives the attempt number (1-based) and the error that caused the retry.
func WithOnRetry(fn func(attempt int, err error)) TaskOption {
	return func(taskOptions *TaskOptions) { taskOptions.OnRetry = fn }
}
