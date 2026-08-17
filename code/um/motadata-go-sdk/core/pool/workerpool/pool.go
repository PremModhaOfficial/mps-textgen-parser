// Package workerpool provides a high-performance, auto-scaling worker Pool implementation
// with comprehensive metrics, retry logic, and graceful shutdown capabilities.
package workerpool

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"
	"github.com/panjf2000/ants/v2"
)

// ═══════════════════════════════════════════════════════════════
// TYPE DEFINITIONS
// ═══════════════════════════════════════════════════════════════

// Pool represents a worker Pool with auto-scaling capabilities and comprehensive metrics.
// It wraps the ants Pool library and provides additional functionality like
// automatic scaling, graceful shutdown, and detailed performance tracking.
type Pool struct {
	// Core Pool configuration
	antsPool   *ants.Pool    // underlying worker Pool implementation
	minWorkers int           // minimum number of workers to maintain
	maxWorkers int           // maximum number of workers allowed
	timeout    time.Duration // default task timeout

	// Lifecycle management
	ctx       context.Context    // Pool context for cancellation
	cancel    context.CancelFunc // cancel function for graceful shutdown
	waitGroup sync.WaitGroup     // tracks running goroutines for clean shutdown

	// State management
	closed atomic.Bool // indicates if Pool is closed
	lock   sync.Mutex  // protects close operations from race conditions

	// Performance metrics (thread-safe atomic counters)
	submitted atomic.Int64 // total tasks submitted to the Pool
	completed atomic.Int64 // total tasks completed successfully
	failed    atomic.Int64 // total tasks that failed
	retried   atomic.Int64 // total retry attempts made
}

// PoolConfig defines configuration parameters for creating a new worker Pool.
type PoolConfig struct {
	// MinWorkers specifies the minimum number of workers to maintain in the Pool
	MinWorkers int
	// MaxWorkers specifies the maximum number of workers the Pool can scale to
	MaxWorkers int
	// Timeout specifies the default timeout for tasks submitted to the Pool
	Timeout time.Duration
}

// Stats provides comprehensive statistics about the worker Pool's performance and state.
// These metrics are useful for monitoring, debugging, and capacity planning.
type Stats struct {
	// Current state metrics
	Running  int // number of currently executing tasks
	Waiting  int // number of tasks waiting for available workers
	Capacity int // current Pool capacity (number of workers)

	// Lifetime performance metrics (atomic counters)
	Submitted int64 // total number of tasks submitted to the Pool
	Completed int64 // total number of tasks completed successfully
	Failed    int64 // total number of tasks that failed
	Retried   int64 // total number of retry attempts made
}

// DefaultConfig provides sensible defaults for worker Pool configuration.
// MinWorkers: 2 (ensures basic concurrency)
// MaxWorkers: 6 (balances performance with resource usage)
// Timeout: 30 seconds (prevents hung tasks from blocking indefinitely)
var DefaultConfig = PoolConfig{
	MinWorkers: 2,
	MaxWorkers: 6,
	Timeout:    30 * time.Second,
}

// ═══════════════════════════════════════════════════════════════
// POOL LIFECYCLE - INITIALIZATION
// ═══════════════════════════════════════════════════════════════

// newWorkerPool creates a new worker Pool with the specified configuration.
// It validates the config parameters and applies defaults where necessary,
// then initializes the underlying ants Pool with appropriate settings.
func newWorkerPool(config PoolConfig) (*Pool, error) {
	// Validate and apply default configuration values
	if config.MinWorkers <= 0 {
		config.MinWorkers = DefaultConfig.MinWorkers
	}

	if config.MaxWorkers <= 0 {
		config.MaxWorkers = DefaultConfig.MaxWorkers
	}

	if config.Timeout <= 0 {
		config.Timeout = DefaultConfig.Timeout
	}

	// Create the underlying ants Pool with optimized settings
	antsPool, err := ants.NewPool(
		config.MinWorkers,
		ants.WithExpiryDuration(time.Minute), // clean up idle workers after 1 minute
		ants.WithPreAlloc(true),              // pre-allocate goroutine stack for better performance
		ants.WithNonblocking(false),          // block when Pool is full (backpressure)
		ants.WithPanicHandler(func(i interface{}) {
			// Gracefully handle panics in worker goroutines
			fmt.Printf("workerpool panic recovered: %v\n%s\n", i, debug.Stack())
		}),
	)

	if err != nil {
		return nil, err
	}

	// Create cancellable context for Pool lifecycle management
	ctx, cancel := context.WithCancel(context.Background())

	// Initialize the Pool with validated configuration
	pool := &Pool{
		antsPool:   antsPool,
		minWorkers: config.MinWorkers,
		maxWorkers: config.MaxWorkers,
		timeout:    config.Timeout,
		ctx:        ctx,
		cancel:     cancel,
	}

	// Start the auto-scaler goroutine to manage dynamic Pool sizing
	pool.waitGroup.Add(1)
	go pool.autoScaler()

	return pool, nil
}

// ═══════════════════════════════════════════════════════════════
// POOL LIFECYCLE - TASK EXECUTION
// ═══════════════════════════════════════════════════════════════

// async submits a task for asynchronous execution (fire and forget).
// The task will be executed when a worker becomes available.
// Returns immediately without waiting for task completion.
func (workerPool *Pool) async(fn func()) error {
	// Check if Pool is still accepting tasks
	if workerPool.closed.Load() {
		return utils.ErrPoolClosed
	}

	// Update metrics: increment submitted counter
	workerPool.submitted.Add(1)

	// Submit task to underlying Pool with error recovery
	return workerPool.antsPool.Submit(func() {
		// Ensure panic recovery for safe execution
		// Deferred function to catch and handle panics
		defer func() {
			if r := recover(); r != nil {
				// Update failure metrics for any panic
				workerPool.failed.Add(1)

				fmt.Printf("task panic: %v\n%s\n", r, debug.Stack())
			}
		}()
		// Execute the user function
		fn()
		// Update metrics: increment completed counter
		workerPool.completed.Add(1)
	})
}

// asyncWithCtx submits a task with context support for cancellation.
// The task will be cancelled if the context is cancelled before execution.
// Useful for tasks that need to respect deadlines or cancellation signals.
func (workerPool *Pool) asyncWithCtx(ctx context.Context, fn func(ctx2 context.Context)) error {
	// Check if Pool is still accepting tasks
	if workerPool.closed.Load() {
		return utils.ErrPoolClosed
	}

	// Update metrics: increment submitted counter
	workerPool.submitted.Add(1)

	// Submit context-aware task to underlying Pool
	return workerPool.antsPool.Submit(func() {
		// Ensure panic recovery for safe execution
		// Deferred function to catch and handle panics
		defer func() {
			if r := recover(); r != nil {
				// Update failure metrics for any panic
				workerPool.failed.Add(1)

				fmt.Printf("task panic: %v\n%s\n", r, debug.Stack())
			}
		}()

		// Check if context was cancelled before execution
		select {
		case <-ctx.Done():
			// Context cancelled, mark as failed and exit
			workerPool.failed.Add(1)
			return
		default:
			// Execute the user function with context
			fn(ctx)
			// Update metrics: increment completed counter
			workerPool.completed.Add(1)
		}
	})
}

// sync executes a task synchronously with context support.
// Blocks until the task completes or the context is cancelled.
// Returns the error from the task execution or context cancellation.
func (workerPool *Pool) sync(ctx context.Context, fn func() error) error {
	// Check if Pool is still accepting tasks
	if workerPool.closed.Load() {
		return utils.ErrPoolClosed
	}

	// Create buffered channel to receive the result
	errs := make(chan error, 1)

	// Update metrics: increment submitted counter
	workerPool.submitted.Add(1)

	// Submit synchronous task to underlying Pool
	poolErr := workerPool.antsPool.Submit(func() {
		// Ensure panic recovery with error channel
		// Deferred function to catch and handle panics
		defer func() {
			if r := recover(); r != nil {
				// Update failure metrics for any panic
				workerPool.failed.Add(1)

				// Send panic as error through channel
				errs <- fmt.Errorf("panic: %v", r)
			}
		}()

		// Execute the user function and handle result
		if err := fn(); err != nil {
			// Task failed, update metrics and send error
			workerPool.failed.Add(1)
			errs <- err
		} else {
			// Task succeeded, update metrics and send success
			workerPool.completed.Add(1)
			errs <- nil
		}
	})

	// Handle Pool submission error
	if poolErr != nil {
		return poolErr
	}

	// Wait for either task completion or context cancellation
	select {
	case err := <-errs:
		// Task completed (successfully or with error)
		return err
	case <-ctx.Done():
		// Context cancelled before task completion
		workerPool.failed.Add(1)
		return ctx.Err()
	}
}

// syncWithOptions executes a task synchronously with comprehensive configuration options.
// Supports timeout, retry logic, and callback handlers for success/failure scenarios.
// This is the most feature-rich execution method, suitable for critical operations.
func (workerPool *Pool) syncWithOptions(ctx context.Context, fn func(ctx2 context.Context) error, opts ...TaskOption) error {
	// Check if Pool is still accepting tasks
	if workerPool.closed.Load() {
		return utils.ErrPoolClosed
	}

	// Apply all provided options to default configuration
	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	// Apply timeout if configured
	if options.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, options.Timeout)
		defer cancel() // ensure timeout cleanup
	}

	// Initialize retry logic
	var lastErr error
	retries := options.Retries + 1 // +1 for initial attempt

	// Execute task with retry logic
	for retry := 1; retry <= retries; retry++ {
		// Create result channel for this attempt
		errs := make(chan error, 1)

		// Update metrics: increment submitted counter
		workerPool.submitted.Add(1)

		// Submit task attempt to underlying Pool
		poolErr := workerPool.antsPool.Submit(func() {
			// Ensure panic recovery with error channel
			// Deferred function to catch and handle panics
			defer func() {
				if r := recover(); r != nil {
					// Update failure metrics for any panic
					workerPool.failed.Add(1)

					// Send panic as error through channel
					errs <- fmt.Errorf("panic: %v", r)
				}
			}()
			// Execute user function and send result
			errs <- fn(ctx)
		})

		// Handle Pool submission error
		if poolErr != nil {
			return poolErr
		}

		// Wait for task completion or context cancellation
		select {
		case err := <-errs:
			// Task completed, check result
			if err == nil {
				// Task succeeded
				workerPool.completed.Add(1)
				if options.OnSuccess != nil {
					options.OnSuccess()
				}
				return nil
			}

			// Task failed, store error for potential retry
			lastErr = err
			workerPool.failed.Add(1)

			// Check if we should retry
			if retry < retries {
				workerPool.retried.Add(1)
				if options.OnRetry != nil {
					options.OnRetry(retry, err)
				}
				// Apply exponential backoff delay
				time.Sleep(options.RetryDelay * time.Duration(retry))
				continue
			}

			// All retries exhausted, execute error callback and return
			if options.OnError != nil {
				options.OnError(err)
			}
			return err

		case <-ctx.Done():
			// Context cancelled during task execution
			workerPool.failed.Add(1)
			if options.OnError != nil {
				options.OnError(ctx.Err())
			}
			return ctx.Err()
		}
	}

	// This should never be reached, but return last error as fallback
	return lastErr
}

// ═══════════════════════════════════════════════════════════════
// POOL LIFECYCLE - MANAGEMENT AND OPTIMIZATION
// ═══════════════════════════════════════════════════════════════

// autoScaler runs in a separate goroutine and automatically adjusts
// the Pool size based on current load and demand patterns.
// It monitors the Pool's state every second and scales up or down as needed.
func (workerPool *Pool) autoScaler() {
	// Ensure proper cleanup and panic recovery
	defer func() {
		workerPool.waitGroup.Done()
		if r := recover(); r != nil {
			fmt.Printf("autoScaler panic recovered: %v\n%s\n", r, debug.Stack())
			// Don't restart - let it die gracefully to avoid infinite panic loops
		}
	}()

	// Use ticker for periodic scaling decisions (1 second interval)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop() // ensure ticker cleanup to prevent goroutine leaks

	// Main auto-scaling loop
	for {
		select {
		case <-workerPool.ctx.Done():
			// Pool is shutting down, exit gracefully
			return
		case <-ticker.C:
			// Check if Pool is closed before attempting to scale
			if workerPool.closed.Load() {
				return
			}
			// Execute scaling logic
			workerPool.runAutoScaler()
		}
	}
}

// runAutoScaler implements the core auto-scaling logic.
// This is separated for better testing and readability.
// Scaling decisions are based on current load patterns and configured limits.
func (workerPool *Pool) runAutoScaler() {
	// Get current Pool metrics for scaling decisions
	waiting := workerPool.antsPool.Waiting() // tasks waiting for workers
	running := workerPool.antsPool.Running() // currently executing tasks
	capacity := workerPool.antsPool.Cap()    // current Pool capacity

	// Scale up logic: increase capacity when demand exceeds current capacity
	if waiting > 0 && capacity < workerPool.maxWorkers {
		// Calculate new capacity (at least current + 2, but respect max limit)
		qualified := min(capacity+max(waiting, 2), workerPool.maxWorkers)
		workerPool.antsPool.Tune(qualified)
	}

	// Scale down logic: reduce capacity when utilization is low
	if waiting == 0 && running < capacity/2 && capacity > workerPool.minWorkers {
		// Calculate new capacity (keep some buffer above running + respect min limit)
		qualified := max(running+2, workerPool.minWorkers)
		workerPool.antsPool.Tune(qualified)
	}
}

// ═══════════════════════════════════════════════════════════════
// POOL LIFECYCLE - MONITORING AND STATISTICS
// ═══════════════════════════════════════════════════════════════

// getPoolStats returns current Pool statistics including real-time metrics
// and cumulative performance counters.
// This provides a comprehensive view of the Pool's health and performance.
func (workerPool *Pool) getPoolStats() Stats {
	return Stats{
		// Real-time state from underlying Pool
		Running:  workerPool.antsPool.Running(),
		Waiting:  workerPool.antsPool.Waiting(),
		Capacity: workerPool.antsPool.Cap(),
		// Atomic counters for thread-safe metrics
		Submitted: workerPool.submitted.Load(),
		Completed: workerPool.completed.Load(),
		Failed:    workerPool.failed.Load(),
		Retried:   workerPool.retried.Load(),
	}
}

// isRunning returns true if the Pool is active and accepting new tasks.
// Returns false if the Pool has been closed or is in the process of shutting down.
func (workerPool *Pool) isRunning() bool {
	return !workerPool.closed.Load()
}

// ═══════════════════════════════════════════════════════════════
// POOL LIFECYCLE - SHUTDOWN
// ═══════════════════════════════════════════════════════════════

// close gracefully shuts down the worker Pool.
// Stops accepting new tasks, cancels the auto-scaler, and releases all resources.
// This method is idempotent - calling it multiple times has no effect.
func (workerPool *Pool) close() {
	// Ensure thread-safe shutdown
	workerPool.lock.Lock()
	defer workerPool.lock.Unlock()

	// Check if already closed (atomic operation for thread safety)
	if workerPool.closed.Swap(true) {
		return // already closed, nothing to do
	}

	// Shutdown sequence: cancel context, wait for goroutines, release Pool
	workerPool.cancel()           // stop auto-scaler and signal shutdown
	workerPool.waitGroup.Wait()   // wait for auto-scaler to exit
	workerPool.antsPool.Release() // release underlying Pool resources
}

// closeWithTimeout gracefully shuts down the Pool with a timeout for pending tasks.
// Waits for currently running tasks to complete before forcing shutdown.
// Returns an error if the timeout is exceeded while waiting for tasks to complete.
func (workerPool *Pool) closeWithTimeout(timeout time.Duration) error {
	// Ensure thread-safe shutdown
	workerPool.lock.Lock()
	defer workerPool.lock.Unlock()

	// Check if already closed (atomic operation for thread safety)
	if workerPool.closed.Swap(true) {
		return nil // already closed
	}

	// Stop accepting new tasks and cancel auto-scaler
	workerPool.cancel()
	workerPool.waitGroup.Wait() // wait for auto-scaler to exit

	// Wait for running tasks to complete with timeout and exponential backoff
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Start with short polling interval, increase gradually
	pollInterval := 10 * time.Millisecond
	maxInterval := 100 * time.Millisecond

	// Poll until all tasks complete or timeout occurs
	for workerPool.antsPool.Running() > 0 {
		select {
		case <-ctx.Done():
			// Timeout reached, force release and return error
			workerPool.antsPool.Release()
			return fmt.Errorf("timeout: %d tasks still running", workerPool.antsPool.Running())
		case <-time.After(pollInterval):
			// Exponential backoff to reduce CPU usage during long waits
			if pollInterval < maxInterval {
				pollInterval *= 2
			}
		}
	}

	// All tasks completed within timeout
	workerPool.antsPool.Release()
	return nil
}
