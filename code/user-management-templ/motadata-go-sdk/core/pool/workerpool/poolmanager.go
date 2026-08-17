package workerpool

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"
)

// ═══════════════════════════════════════════════════════════════
// GLOBAL SINGLETON MANAGEMENT
// ═══════════════════════════════════════════════════════════════

// Global singleton pattern for worker Pool management.
// This ensures a single Pool instance across the application.
var (
	instance    atomic.Pointer[Pool] // thread-safe pointer to singleton Pool instance
	once        sync.Once            // ensures initialization happens only once
	initErr     error                // stores initialization error if any
	initialized atomic.Bool          // atomic flag indicating initialization status
)

// ═══════════════════════════════════════════════════════════════
// INITIALIZATION AND CONFIGURATION
// ═══════════════════════════════════════════════════════════════

// Init initializes the global worker Pool with optional configuration.
// This function is thread-safe and ensures only one initialization occurs.
// If no config is provided, DefaultConfig will be used.
// Returns an error if initialization fails.
func Init(config ...PoolConfig) error {
	// Use sync.Once to guarantee single initialization
	once.Do(func() {
		// Apply provided config or use defaults
		cfg := DefaultConfig
		if len(config) > 0 {
			cfg = PoolConfig{
				MinWorkers:   config[0].MinWorkers,
				MaxWorkers:   config[0].MaxWorkers,
				Timeout:      config[0].Timeout,
				MaxQueueSize: config[0].MaxQueueSize,
				PanicHandler: config[0].PanicHandler,
			}
		}

		// Create the worker Pool instance
		pool, err := newWorkerPool(cfg)
		if err != nil {
			initErr = err
			return
		}

		// Store instance and mark as initialized
		instance.Store(pool)
		initialized.Store(true)
	})
	return initErr
}

// MustInit initializes the global worker Pool and panics if initialization fails.
// This is useful when you want to fail fast during application startup.
// Use this in main() or init() functions where failure should stop the program.
func MustInit(cfg ...PoolConfig) {
	if err := Init(cfg...); err != nil {
		panic(fmt.Sprintf("worker Pool init failed: %v", err))
	}
}

// ═══════════════════════════════════════════════════════════════
// INTERNAL HELPERS
// ═══════════════════════════════════════════════════════════════

// getInstance safely returns the global Pool instance.
// Returns an error if the Pool has not been initialized.
// This is a private helper used by all public API functions.
func getInstance() (*Pool, error) {
	// Check if Pool was initialized
	if !initialized.Load() {
		return nil, utils.ErrWorkerPoolNotInitialized
	}

	// Load the Pool instance
	pool := instance.Load()
	if pool == nil {
		return nil, utils.ErrWorkerPoolNotInitialized
	}

	return pool, nil
}

// GetWorkerPool returns the underlying Pool instance for advanced usage.
// Returns nil if the Pool is not initialized.
// Use this only when you need direct access to Pool internals.
func GetWorkerPool() *Pool {
	pool, err := getInstance()
	if err != nil {
		return nil
	}
	return pool
}

// ═══════════════════════════════════════════════════════════════
// BASIC TASK EXECUTION API
// ═══════════════════════════════════════════════════════════════

// Async runs a function asynchronously (fire and forget).
// Returns immediately without waiting for the function to complete.
// The function will be executed by an available worker from the Pool.
func Async(fn func()) error {
	pool, err := getInstance()
	if err != nil {
		return err
	}
	return pool.async(fn)
}

// AsyncWithCtx runs a function asynchronously with context support.
// The function will be cancelled if the context is cancelled before execution.
// Useful for implementing timeouts and cancellation.
func AsyncWithCtx(ctx context.Context, fn func(context.Context)) error {
	pool, err := getInstance()
	if err != nil {
		return err
	}
	return pool.asyncWithCtx(ctx, fn)
}

// SyncWithCtx executes a function synchronously with context support.
// Blocks until the function completes or the context is cancelled.
// Returns the error from the function or context cancellation.
func SyncWithCtx(ctx context.Context, fn func() error) error {
	pool, err := getInstance()
	if err != nil {
		return err
	}
	return pool.sync(ctx, fn)
}

// SyncWithTimeout executes a function synchronously with a timeout.
// Blocks until the function completes or the timeout is reached.
// Returns the function's error or a timeout error.
func SyncWithTimeout(timeout time.Duration, fn func() error) error {
	pool, err := getInstance()
	if err != nil {
		return err
	}
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return pool.sync(ctx, fn)
}

// ═══════════════════════════════════════════════════════════════
// ASYNC EXECUTION WITH RESULT CHANNELS
// ═══════════════════════════════════════════════════════════════

// AsyncWithError runs a function asynchronously and returns an error channel.
// The channel will receive the function's error (or nil on success).
// The channel is closed after the result is sent.
func AsyncWithError(fn func() error) <-chan error {
	result := make(chan error, 1)
	pool, err := getInstance()
	if err != nil {
		result <- err
		close(result)
		return result
	}

	submitErr := pool.async(func() {
		// Execute function and send result, then close channel
		result <- fn()
		close(result)
	})

	// Only close here if submit failed (goroutine never started)
	// If submit succeeded, the goroutine will close the channel
	if submitErr != nil {
		result <- submitErr
		close(result)
	}
	return result
}

// AsyncWithCtxError runs a function with context and returns an error channel.
// The function receives the context and can respond to cancellation.
// The channel is closed after the result is sent.
func AsyncWithCtxError(ctx context.Context, fn func(context.Context) error) <-chan error {
	result := make(chan error, 1)
	pool, err := getInstance()
	if err != nil {
		result <- err
		close(result)
		return result
	}

	submitErr := pool.asyncWithCtx(ctx, func(c context.Context) {
		// Check context and execute function, then close channel
		select {
		case <-c.Done():
			result <- c.Err()
		default:
			result <- fn(c)
		}
		close(result)
	})

	// Only close here if submit failed (goroutine never started)
	// If submit succeeded, the goroutine will close the channel
	if submitErr != nil {
		result <- submitErr
		close(result)
	}
	return result
}

// Result holds a value and error pair for generic async operations.
// Used by AsyncResult and related functions to return both value and error.
type Result[T any] struct {
	Value T     // The result value (zero value if error occurred)
	Error error // The error (nil if successful)
}

// AsyncResult runs a function asynchronously and returns a value channel.
// The function can return both a value and an error.
// The channel receives a Result containing both the value and error.
func AsyncResult[T any](fn func() (T, error)) <-chan Result[T] {
	result := make(chan Result[T], 1)
	pool, err := getInstance()
	if err != nil {
		var zero T
		result <- Result[T]{Value: zero, Error: err}
		close(result)
		return result
	}

	submitErr := pool.async(func() {
		defer close(result) // Ensure channel is always closed
		var value T
		value, err = fn()
		result <- Result[T]{Value: value, Error: err}
	})

	if submitErr != nil {
		var zero T
		result <- Result[T]{Value: zero, Error: submitErr}
		close(result)
	}
	return result
}

// AsyncWithCtxResult runs a function with context and returns a value channel.
// The function receives context and returns both a value and an error.
// Supports context cancellation during execution.
func AsyncWithCtxResult[T any](ctx context.Context, fn func(context.Context) (T, error)) <-chan Result[T] {
	result := make(chan Result[T], 1)
	pool, err := getInstance()
	if err != nil {
		var zero T
		result <- Result[T]{Value: zero, Error: err}
		close(result)
		return result
	}

	submitErr := pool.asyncWithCtx(ctx, func(c context.Context) {
		defer close(result) // Ensure channel is always closed
		var zero T
		select {
		case <-c.Done():
			result <- Result[T]{Value: zero, Error: c.Err()}
		default:
			var value T
			value, err = fn(c)
			result <- Result[T]{Value: value, Error: err}
		}
	})

	if submitErr != nil {
		var zero T
		result <- Result[T]{Value: zero, Error: submitErr}
		close(result)
	}
	return result
}

// ═══════════════════════════════════════════════════════════════
// SYNCHRONIZATION AND WAITING
// ═══════════════════════════════════════════════════════════════

// WaitForAll waits for multiple result channels and returns all results.
// Blocks until all channels have produced results.
// Results are returned in the same order as the input channels.
func WaitForAll[T any](results ...<-chan Result[T]) []Result[T] {
	all := make([]Result[T], len(results))
	for i, result := range results {
		all[i] = <-result
	}
	return all
}

// WaitForAllErrors waits for multiple error channels and returns all errors.
// Blocks until all channels have produced errors (or nil).
// Errors are returned in the same order as the input channels.
func WaitForAllErrors(results ...<-chan error) []error {
	all := make([]error, len(results))
	for i, result := range results {
		all[i] = <-result
	}
	return all
}

// WaitWithTimeout waits for a result with timeout.
// Returns the result and true if received before timeout.
// Returns zero value and false if timeout occurs.
func WaitWithTimeout[T any](result <-chan Result[T], timeout time.Duration) (Result[T], bool) {
	select {
	case r := <-result:
		return r, true
	case <-time.After(timeout):
		return Result[T]{}, false
	}
}

// WaitErrorWithTimeout waits for an error with timeout.
// Returns the error and true if received before timeout.
// Returns nil and false if timeout occurs.
func WaitErrorWithTimeout(result <-chan error, timeout time.Duration) (error, bool) {
	select {
	case err := <-result:
		return err, true
	case <-time.After(timeout):
		return nil, false
	}
}

// ═══════════════════════════════════════════════════════════════
// ADVANCED TASK SUBMISSION
// ═══════════════════════════════════════════════════════════════

// Task represents a unit of work with comprehensive configuration options.
// Use this for complex tasks that need retry logic, timeouts, and callbacks.
type Task struct {
	Name       string                       // Optional task name for debugging
	Fn         func(context.Context) error  // The function to execute
	Timeout    time.Duration                // Task-specific timeout
	Retries    int                          // Number of retry attempts
	RetryDelay time.Duration                // Delay between retries
	OnSuccess  func()                       // Callback on successful completion
	OnError    func(error)                  // Callback on final failure
	OnRetry    func(attempt int, err error) // Callback before each retry
}

// Submit executes a task with full configuration options.
// Supports retries, timeouts, and event callbacks.
// This is the most feature-rich task submission method.
func Submit(ctx context.Context, task Task) error {
	pool, err := getInstance()
	if err != nil {
		return err
	}

	// Build task options from Task struct
	opts := []TaskOption{
		WithName(task.Name),
	}

	if task.Timeout > 0 {
		opts = append(opts, WithTimeout(task.Timeout))
	}
	if task.Retries > 0 {
		delay := task.RetryDelay
		if delay == 0 {
			delay = time.Second
		}
		opts = append(opts, WithRetry(task.Retries, delay))
	}
	if task.OnSuccess != nil {
		opts = append(opts, WithOnSuccess(task.OnSuccess))
	}
	if task.OnError != nil {
		opts = append(opts, WithOnError(task.OnError))
	}
	if task.OnRetry != nil {
		opts = append(opts, WithOnRetry(task.OnRetry))
	}

	return pool.syncWithOptions(ctx, task.Fn, opts...)
}

// ═══════════════════════════════════════════════════════════════
// BATCH OPERATIONS
// ═══════════════════════════════════════════════════════════════

// Batch runs multiple functions concurrently and waits for all to complete.
// Returns a slice of errors in the same order as the input functions.
// nil entries indicate successful execution.
func Batch(ctx context.Context, fns ...func() error) []error {
	errs := make([]error, len(fns))

	pool, err := getInstance()
	if err != nil {
		for i := range errs {
			errs[i] = err
		}
		return errs
	}

	if len(fns) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(len(fns))

	// Create a context for cancellation propagation
	batchCtx, batchCancel := context.WithCancel(ctx)
	defer batchCancel()

	for i, fn := range fns {
		i2, fn2 := i, fn // Capture loop variables
		submitErr := pool.async(func() {
			defer wg.Done()
			// Check context before execution
			select {
			case <-batchCtx.Done():
				errs[i2] = batchCtx.Err()
				return
			default:
			}

			// Execute with proper error handling
			errs[i2] = fn2() // Store error without cancelling batch
		})

		if submitErr != nil {
			errs[i2] = submitErr
			wg.Done() // Decrease counter for failed submissions
		}
	}

	// Wait for all tasks to complete
	// We must wait for all goroutines to finish writing to errs before returning
	// to avoid race conditions
	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case <-done:
		// All tasks completed normally
		return errs
	case <-ctx.Done():
		// Context cancelled - signal all tasks to stop
		batchCancel()
		// Must wait for all goroutines to finish to avoid race on errs slice
		<-done
		return errs
	}
}

// BatchWithContext runs functions with context passed to each function.
// Each function receives the context and can respond to cancellation.
// Returns errors in the same order as input functions.
func BatchWithContext(ctx context.Context, fns ...func(context.Context) error) []error {
	errs := make([]error, len(fns))
	pool, err := getInstance()
	if err != nil {
		for i := range errs {
			errs[i] = err
		}
		return errs
	}
	if len(fns) == 0 {
		return nil
	}

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(fns))

	// Create a context for batch coordination
	batchCtx, batchCancel := context.WithCancel(ctx)
	defer batchCancel()

	for i, fn := range fns {
		i2, fn2 := i, fn // Capture loop variables
		submitErr := pool.asyncWithCtx(batchCtx, func(c context.Context) {
			defer waitGroup.Done()
			// Immediate context check
			select {
			case <-c.Done():
				errs[i2] = c.Err()
				return
			default:
			}

			// Execute with proper error handling
			errs[i2] = fn2(c) // Store error without cancelling batch
		})

		if submitErr != nil {
			errs[i2] = submitErr
			waitGroup.Done() // Decrease counter for failed submissions
		}
	}

	// Wait for all tasks to complete
	// We must wait for all goroutines to finish writing to errs before returning
	// to avoid race conditions
	done := make(chan struct{})
	go func() {
		defer close(done)
		waitGroup.Wait()
	}()

	select {
	case <-done:
		// All tasks completed normally
		return errs
	case <-ctx.Done():
		// Context cancelled - signal all tasks to stop
		batchCancel()
		// Must wait for all goroutines to finish to avoid race on errs slice
		<-done
		return errs
	}
}

// Parallel runs all functions concurrently and returns the first error.
// Returns nil if all functions complete successfully.
// Useful when you need all operations to succeed.
func Parallel(ctx context.Context, fns ...func() error) error {
	errs := Batch(ctx, fns...)
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// Sequential runs functions in order, stops on first error.
// Each function must complete successfully before the next starts.
// Returns the first error encountered or nil if all succeed.
func Sequential(ctx context.Context, fns ...func() error) error {
	for i, fn := range fns {
		if err := SyncWithCtx(ctx, fn); err != nil {
			return fmt.Errorf("step %d failed: %w", i, err)
		}
	}
	return nil
}

// ═══════════════════════════════════════════════════════════════
// SCHEDULED TASK EXECUTION
// ═══════════════════════════════════════════════════════════════

// ScheduledTask represents a cancellable scheduled task.
// Call Cancel() to stop the task from executing.
type ScheduledTask interface {
	// Cancel stops the scheduled task from executing.
	// Returns true if the task was cancelled, false if it already executed.
	Cancel() bool
}

// timerTask wraps a time.Timer for the ScheduledTask interface.
// Used internally by After() for one-time scheduled tasks.
type timerTask struct {
	*time.Timer
}

// Cancel stops the timer and prevents the task from executing.
func (t *timerTask) Cancel() bool {
	return t.Timer.Stop()
}

// After runs a function after a specified delay.
// Returns a ScheduledTask that can be cancelled before execution.
// More efficient than using goroutines with time.Sleep.
func After(delay time.Duration, fn func()) ScheduledTask {
	timer := time.AfterFunc(delay, func() {
		pool, err := getInstance()
		if err != nil {
			return // Pool not available, skip execution
		}
		if pool.isRunning() {
			_ = pool.async(fn)
		}
	})

	return &timerTask{Timer: timer}
}

// repeatingTask wraps a cancel function for the ScheduledTask interface.
// Used internally by Every() for recurring scheduled tasks.
type repeatingTask struct {
	cancel context.CancelFunc
}

// Cancel stops the repeating task.
func (r *repeatingTask) Cancel() bool {
	r.cancel()
	return true
}

// Every runs a function repeatedly at a specified interval.
// Returns a ScheduledTask that can be cancelled to stop the repetition.
// The function runs in the worker Pool, not in the scheduler goroutine.
func Every(ctx context.Context, interval time.Duration, fn func()) ScheduledTask {
	ctx, cancel := context.WithCancel(ctx)

	// Use a more efficient approach with cleanup tracking
	go func() {
		defer cancel() // Ensure cleanup on exit
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				pool, err := getInstance()
				if err != nil || !pool.isRunning() {
					return // Pool not available or stopped
				}

				// Submit task with error handling
				if submitErr := pool.async(fn); submitErr != nil {
					// If we can't submit, stop the scheduler
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return &repeatingTask{cancel: cancel}
}

// ═══════════════════════════════════════════════════════════════
// MONITORING AND STATISTICS
// ═══════════════════════════════════════════════════════════════

// GetPoolStats returns current Pool statistics.
// Provides real-time metrics about Pool state and performance.
// Returns zero values if Pool is not initialized.
func GetPoolStats() Stats {
	pool, err := getInstance()
	if err != nil {
		return Stats{}
	}
	return pool.getPoolStats()
}

// IsRunning returns true if the Pool is active and accepting tasks.
// Returns false if the Pool is not initialized or has been shut down.
func IsRunning() bool {
	pool, err := getInstance()
	if err != nil {
		return false
	}
	return pool.isRunning()
}

// ═══════════════════════════════════════════════════════════════
// HEALTH CHECK API
// ═══════════════════════════════════════════════════════════════

// HealthStatus represents the health state of the worker pool.
// This is designed for Kubernetes liveness/readiness probes and monitoring systems.
type HealthStatus struct {
	// Healthy indicates whether the pool is operational and can accept tasks
	Healthy bool `json:"healthy"`
	// Reason provides a human-readable explanation of the health status
	Reason string `json:"reason"`
	// Running is the number of currently executing tasks
	Running int `json:"running"`
	// Waiting is the number of tasks waiting for available workers
	Waiting int `json:"waiting"`
	// Capacity is the current pool capacity (number of workers)
	Capacity int `json:"capacity"`
	// Utilization is the percentage of pool capacity in use (0.0 to 1.0)
	Utilization float64 `json:"utilization"`
}

// GetHealthStatus returns the current health status of the worker pool.
// This is designed for Kubernetes liveness/readiness probes.
//
// Health criteria:
//   - Pool must be initialized
//   - Pool must not be closed
//   - Pool must be accepting tasks
//
// Example usage with HTTP health endpoint:
//
//	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
//	    status := workerpool.GetHealthStatus()
//	    if !status.Healthy {
//	        w.WriteHeader(http.StatusServiceUnavailable)
//	    }
//	    json.NewEncoder(w).Encode(status)
//	})
func GetHealthStatus() HealthStatus {
	pool, err := getInstance()
	if err != nil {
		return HealthStatus{
			Healthy: false,
			Reason:  "pool not initialized",
		}
	}

	// Check if pool is closed
	if pool.closed.Load() {
		return HealthStatus{
			Healthy: false,
			Reason:  "pool is closed",
		}
	}

	// Get current stats
	stats := pool.getPoolStats()

	// Calculate utilization
	var utilization float64
	if stats.Capacity > 0 {
		utilization = float64(stats.Running) / float64(stats.Capacity)
	}

	return HealthStatus{
		Healthy:     true,
		Reason:      "pool is healthy and accepting tasks",
		Running:     stats.Running,
		Waiting:     stats.Waiting,
		Capacity:    stats.Capacity,
		Utilization: utilization,
	}
}

// ═══════════════════════════════════════════════════════════════
// SHUTDOWN AND CLEANUP
// ═══════════════════════════════════════════════════════════════

// Shutdown gracefully stops the global worker Pool.
// This will stop accepting new tasks and close all workers.
// Does nothing if the Pool is not initialized.
func Shutdown() {
	if pool := instance.Load(); pool != nil {
		pool.close()
		instance.Store(nil)
		initialized.Store(false)
	}
}

// ShutdownWithTimeout gracefully stops the global worker Pool with a timeout.
// Waits for currently running tasks to complete before forcing shutdown.
// Returns an error if tasks are still running after the timeout.
func ShutdownWithTimeout(timeout time.Duration) error {
	if pool := instance.Load(); pool != nil {
		err := pool.closeWithTimeout(timeout)
		instance.Store(nil)
		initialized.Store(false)
		return err
	}
	return nil
}
