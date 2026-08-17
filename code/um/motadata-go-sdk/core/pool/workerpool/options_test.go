package workerpool

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ═══════════════════════════════════════════════════════════════
// OPTIONS COVERAGE TESTS
// ═══════════════════════════════════════════════════════════════

// Test all option functions comprehensively
func TestTaskOptionsComprehensive(t *testing.T) {
	assertions := assert.New(t)

	// Test default options
	opts := defaultOptions()
	assertions.NotNil(opts)
	assertions.Empty(opts.Name)
	assertions.Equal(30*time.Second, opts.Timeout)
	assertions.Equal(0, opts.Retries)
	assertions.Equal(time.Second, opts.RetryDelay)
	assertions.Nil(opts.OnSuccess)
	assertions.Nil(opts.OnError)
	assertions.Nil(opts.OnRetry)

	// Test WithName
	nameOpt := WithName("test-task")
	nameOpt(opts)
	assertions.Equal("test-task", opts.Name)

	// Test WithTimeout
	timeoutOpt := WithTimeout(5 * time.Second)
	timeoutOpt(opts)
	assertions.Equal(5*time.Second, opts.Timeout)

	// Test WithRetry
	retryOpt := WithRetry(3, 2*time.Second)
	retryOpt(opts)
	assertions.Equal(3, opts.Retries)
	assertions.Equal(2*time.Second, opts.RetryDelay)

	// Test WithOnSuccess
	var successCalled bool
	successOpt := WithOnSuccess(func() {
		successCalled = true
	})
	successOpt(opts)
	assertions.NotNil(opts.OnSuccess)
	opts.OnSuccess() // Call it
	assertions.True(successCalled)

	// Test WithOnError
	var errorCalled bool
	var receivedError error
	errorOpt := WithOnError(func(err error) {
		errorCalled = true
		receivedError = err
	})
	errorOpt(opts)
	assertions.NotNil(opts.OnError)
	testErr := errors.New("test error")
	opts.OnError(testErr) // Call it
	assertions.True(errorCalled)
	assertions.Equal(testErr, receivedError)

	// Test WithOnRetry
	var retryCalled bool
	var retryAttempt int
	var retryError error
	retryOptFunc := WithOnRetry(func(attempt int, err error) {
		retryCalled = true
		retryAttempt = attempt
		retryError = err
	})
	retryOptFunc(opts)
	assertions.NotNil(opts.OnRetry)
	testRetryErr := errors.New("retry error")
	opts.OnRetry(2, testRetryErr) // Call it
	assertions.True(retryCalled)
	assertions.Equal(2, retryAttempt)
	assertions.Equal(testRetryErr, retryError)
}

// Test options with nil callbacks
func TestTaskOptionsNilCallbacks(t *testing.T) {
	assertions := assert.New(t)

	opts := defaultOptions()

	// Test with nil callbacks (should not panic)
	WithOnSuccess(nil)(opts)
	WithOnError(nil)(opts)
	WithOnRetry(nil)(opts)

	assertions.Nil(opts.OnSuccess)
	assertions.Nil(opts.OnError)
	assertions.Nil(opts.OnRetry)
}

// Test options edge values
func TestTaskOptionsEdgeValues(t *testing.T) {
	assertions := assert.New(t)

	opts := defaultOptions()

	// Test zero and negative values
	WithTimeout(0)(opts)
	assertions.Equal(time.Duration(0), opts.Timeout)

	WithTimeout(-time.Second)(opts)
	assertions.Equal(-time.Second, opts.Timeout)

	WithRetry(0, 0)(opts)
	assertions.Equal(0, opts.Retries)
	assertions.Equal(time.Duration(0), opts.RetryDelay)

	WithRetry(-1, -time.Second)(opts)
	assertions.Equal(-1, opts.Retries)
	assertions.Equal(-time.Second, opts.RetryDelay)

	// Test very large values
	WithTimeout(24 * time.Hour)(opts)
	assertions.Equal(24*time.Hour, opts.Timeout)

	WithRetry(1000, time.Hour)(opts)
	assertions.Equal(1000, opts.Retries)
	assertions.Equal(time.Hour, opts.RetryDelay)
}

// Test multiple options applied in sequence
func TestTaskOptionsChaining(t *testing.T) {
	assertions := assert.New(t)

	opts := defaultOptions()

	// Apply multiple options
	options := []TaskOption{
		WithName("chained-task"),
		WithTimeout(10 * time.Second),
		WithRetry(5, 500*time.Millisecond),
		WithOnSuccess(func() {}),
		WithOnError(func(error) {}),
		WithOnRetry(func(int, error) {}),
	}

	for _, opt := range options {
		opt(opts)
	}

	assertions.Equal("chained-task", opts.Name)
	assertions.Equal(10*time.Second, opts.Timeout)
	assertions.Equal(5, opts.Retries)
	assertions.Equal(500*time.Millisecond, opts.RetryDelay)
	assertions.NotNil(opts.OnSuccess)
	assertions.NotNil(opts.OnError)
	assertions.NotNil(opts.OnRetry)
}

// Test option overriding
func TestTaskOptionsOverriding(t *testing.T) {
	assertions := assert.New(t)

	opts := defaultOptions()

	// Set initial values
	WithName("first")(opts)
	WithTimeout(time.Second)(opts)
	WithRetry(1, time.Millisecond)(opts)

	assertions.Equal("first", opts.Name)
	assertions.Equal(time.Second, opts.Timeout)
	assertions.Equal(1, opts.Retries)
	assertions.Equal(time.Millisecond, opts.RetryDelay)

	// Override with new values
	WithName("second")(opts)
	WithTimeout(2 * time.Second)(opts)
	WithRetry(2, 2*time.Millisecond)(opts)

	assertions.Equal("second", opts.Name)
	assertions.Equal(2*time.Second, opts.Timeout)
	assertions.Equal(2, opts.Retries)
	assertions.Equal(2*time.Millisecond, opts.RetryDelay)
}

// Test TaskOptions struct directly
func TestTaskOptionsStruct(t *testing.T) {
	assertions := assert.New(t)

	// Test zero value
	var opts TaskOptions
	assertions.Empty(opts.Name)
	assertions.Equal(time.Duration(0), opts.Timeout)
	assertions.Equal(0, opts.Retries)
	assertions.Equal(time.Duration(0), opts.RetryDelay)
	assertions.Nil(opts.OnSuccess)
	assertions.Nil(opts.OnError)
	assertions.Nil(opts.OnRetry)

	// Test with values
	opts = TaskOptions{
		Name:       "test",
		Timeout:    time.Minute,
		Retries:    10,
		RetryDelay: time.Second,
		OnSuccess:  func() {},
		OnError:    func(error) {},
		OnRetry:    func(int, error) {},
	}

	assertions.Equal("test", opts.Name)
	assertions.Equal(time.Minute, opts.Timeout)
	assertions.Equal(10, opts.Retries)
	assertions.Equal(time.Second, opts.RetryDelay)
	assertions.NotNil(opts.OnSuccess)
	assertions.NotNil(opts.OnError)
	assertions.NotNil(opts.OnRetry)
}

// Test empty string name
func TestTaskOptionsEmptyName(t *testing.T) {
	assertions := assert.New(t)

	opts := defaultOptions()

	// Set empty string explicitly
	WithName("")(opts)
	assertions.Empty(opts.Name)

	// Set name then override with empty
	WithName("test")(opts)
	assertions.Equal("test", opts.Name)
	WithName("")(opts)
	assertions.Empty(opts.Name)
}

// Test callback function execution
func TestTaskOptionsCallbackExecution(t *testing.T) {
	assertions := assert.New(t)

	var (
		successCount int
		errorCount   int
		retryCount   int
		lastError    error
		lastAttempt  int
	)

	opts := &TaskOptions{
		OnSuccess: func() {
			successCount++
		},
		OnError: func(err error) {
			errorCount++
			lastError = err
		},
		OnRetry: func(attempt int, err error) {
			retryCount++
			lastAttempt = attempt
			lastError = err
		},
	}

	// Test multiple calls
	opts.OnSuccess()
	opts.OnSuccess()
	assertions.Equal(2, successCount)

	testErr1 := errors.New("error 1")
	testErr2 := errors.New("error 2")
	opts.OnError(testErr1)
	opts.OnError(testErr2)
	assertions.Equal(2, errorCount)
	assertions.Equal(testErr2, lastError)

	opts.OnRetry(1, testErr1)
	opts.OnRetry(3, testErr2)
	assertions.Equal(2, retryCount)
	assertions.Equal(3, lastAttempt)
	assertions.Equal(testErr2, lastError)
}
