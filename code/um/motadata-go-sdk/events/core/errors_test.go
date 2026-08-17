package core

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

/* ========================================================================================================
   ERROR TYPES TESTS
   ======================================================================================================== */

func TestError(t *testing.T) {

	assertions := assert.New(t)

	err := NewError("connection", "connect", ErrNotConnected)

	assertions.Contains(err.Error(), "connection", "Error should contain kind")
	assertions.Contains(err.Error(), "connect", "Error should contain operation")
	assertions.Contains(err.Error(), "not connected", "Error should contain underlying error")
}

func TestErrorUnwrap(t *testing.T) {

	assertions := assert.New(t)

	err := NewError("connection", "connect", ErrNotConnected)

	assertions.True(errors.Is(err, ErrNotConnected), "Unwrap should return underlying error")
}

func TestErrorWithDetails(t *testing.T) {

	assertions := assert.New(t)

	err := NewError("connection", "connect", ErrNotConnected).
		WithDetails(map[string]any{
			"server": "nats://localhost:4222",
			"retry":  3,
		})

	assertions.NotNil(err.Details, "Details should not be nil")
	assertions.Equal("nats://localhost:4222", err.Details["server"], "Server detail should match")
	assertions.Equal(3, err.Details["retry"], "Retry detail should match")
}

func TestTenantError(t *testing.T) {

	assertions := assert.New(t)

	err := NewTenantError("tenant-123", "disconnect", ErrTenantDisconnected)

	assertions.Contains(err.Error(), "tenant-123", "Error should contain tenant ID")
	assertions.Contains(err.Error(), "disconnect", "Error should contain operation")
	assertions.True(errors.Is(err, ErrTenantDisconnected), "Unwrap should return underlying error")
}

func TestConfigError(t *testing.T) {

	assertions := assert.New(t)

	err := ConfigError{
		Field:   "Servers",
		Message: "at least one server required",
	}

	assertions.Contains(err.Error(), "Servers", "Error should contain field")
	assertions.Contains(err.Error(), "at least one server required", "Error should contain message")
}

func TestSerializationError(t *testing.T) {

	assertions := assert.New(t)

	underlying := errors.New("json unmarshal error")
	err := SerializationError{
		Operation: "deserialize",
		Type:      "MyStruct",
		Err:       underlying,
	}

	assertions.Contains(err.Error(), "deserialize", "Error should contain operation")
	assertions.Contains(err.Error(), "MyStruct", "Error should contain type")
	assertions.True(errors.Is(err, underlying), "Unwrap should return underlying error")
}

/* ========================================================================================================
   RETRYABLE ERROR TESTS
   ======================================================================================================== */

func TestIsRetryable(t *testing.T) {

	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"ErrNotConnected", ErrNotConnected, true},
		{"ErrConnectionTimeout", ErrConnectionTimeout, true},
		{"ErrPublishTimeout", ErrPublishTimeout, true},
		{"ErrReconnectFailed", ErrReconnectFailed, true},

		// Non-retryable errors
		{"ErrInvalidSubject", ErrInvalidSubject, false},
		{"ErrInvalidMessage", ErrInvalidMessage, false},
		{"ErrInvalidConfig", ErrInvalidConfig, false},
		{"ErrSerializationFailed", ErrSerializationFailed, false},
		{"ErrPermissionDenied", ErrPermissionDenied, false},
		{"ErrDuplicateMsg", ErrDuplicateMsg, false},
		{"ErrShutdownInProgress", ErrShutdownInProgress, false},
		{"ErrConnectionClosed", ErrConnectionClosed, false},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			assertions := assert.New(t)
			result := IsRetryable(tc.err)
			assertions.Equal(tc.expected, result, "IsRetryable for %s should be %v", tc.name, tc.expected)
		})
	}
}

func TestIsRetryableSerializationError(t *testing.T) {

	assertions := assert.New(t)

	err := SerializationError{
		Operation: "serialize",
		Type:      "MyStruct",
		Err:       errors.New("test"),
	}

	assertions.False(IsRetryable(err), "SerializationError should not be retryable")
}

func TestIsRetryableConfigError(t *testing.T) {

	assertions := assert.New(t)

	err := ConfigError{
		Field:   "Servers",
		Message: "required",
	}

	assertions.False(IsRetryable(err), "ConfigError should not be retryable")
}

/* ========================================================================================================
   TEMPORARY ERROR TESTS
   ======================================================================================================== */

func TestIsTemporary(t *testing.T) {

	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"ErrConnectionTimeout", ErrConnectionTimeout, true},
		{"ErrPublishTimeout", ErrPublishTimeout, true},
		{"ErrReconnectFailed", ErrReconnectFailed, true},
		{"ErrNoAck", ErrNoAck, true},

		// Non-temporary errors
		{"ErrNotConnected", ErrNotConnected, false},
		{"ErrInvalidSubject", ErrInvalidSubject, false},
		{"ErrPermissionDenied", ErrPermissionDenied, false},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			assertions := assert.New(t)
			result := IsTemporary(tc.err)
			assertions.Equal(tc.expected, result, "IsTemporary for %s should be %v", tc.name, tc.expected)
		})
	}
}

/* ========================================================================================================
   MULTI ERROR TESTS
   ======================================================================================================== */

func TestMultiError(t *testing.T) {

	assertions := assert.New(t)

	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	err3 := errors.New("error 3")

	multiErr := &MultiError{
		Errors: []error{err1, err2, err3},
	}

	assertions.Contains(multiErr.Error(), "3 errors occurred", "Should indicate number of errors")
	assertions.Equal(err1, multiErr.Unwrap(), "Unwrap should return first error")
	assertions.Equal(3, len(multiErr.All()), "All should return all errors")
}

func TestMultiErrorSingle(t *testing.T) {

	assertions := assert.New(t)

	err1 := errors.New("single error")

	multiErr := &MultiError{
		Errors: []error{err1},
	}

	assertions.Equal("single error", multiErr.Error(), "Single error should return its message")
}

func TestMultiErrorEmpty(t *testing.T) {

	assertions := assert.New(t)

	multiErr := &MultiError{
		Errors: []error{},
	}

	assertions.Equal("no errors", multiErr.Error(), "Empty MultiError should say no errors")
	assertions.Nil(multiErr.Unwrap(), "Unwrap on empty should return nil")
}

func TestNewMultiError(t *testing.T) {

	assertions := assert.New(t)

	err1 := errors.New("error 1")
	err2 := errors.New("error 2")

	// With valid errors
	multiErr := NewMultiError([]error{err1, nil, err2, nil})

	assertions.NotNil(multiErr, "NewMultiError should not be nil")
	assertions.Equal(2, len(multiErr.Errors), "Should filter out nil errors")

	// All nil errors
	nilMultiErr := NewMultiError([]error{nil, nil})

	assertions.Nil(nilMultiErr, "NewMultiError with all nils should return nil")

	// Empty slice
	emptyMultiErr := NewMultiError([]error{})

	assertions.Nil(emptyMultiErr, "NewMultiError with empty slice should return nil")
}
