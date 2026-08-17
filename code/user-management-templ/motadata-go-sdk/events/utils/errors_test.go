package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

/* ======== STRUCTURED ERROR TESTS ======== */

func TestError(t *testing.T) {
	t.Run("with operation", func(t *testing.T) {
		assertions := assert.New(t)
		err := NewError(testKindConnection, testOpConnect, ErrNotConnected)

		assertions.Contains(err.Error(), testKindConnection, "Error should contain kind")
		assertions.Contains(err.Error(), testOpConnect, "Error should contain operation")
		assertions.Contains(err.Error(), testMsgNotConn, "Error should contain underlying error")
	})

	t.Run("without operation", func(t *testing.T) {
		assertions := assert.New(t)
		err := NewError(testKindConnection, "", ErrNotConnected)

		assertions.Contains(err.Error(), testKindConnection)
		assertions.Contains(err.Error(), testMsgNotConn)
		assertions.NotContains(err.Error(), ": :")
	})
}

func TestErrorUnwrap(t *testing.T) {
	assertions := assert.New(t)

	err := NewError(testKindConnection, testOpConnect, ErrNotConnected)

	assertions.True(errors.Is(err, ErrNotConnected), "Unwrap should return underlying error")
}

func TestErrorWithDetails(t *testing.T) {
	assertions := assert.New(t)

	err := NewError(testKindConnection, testOpConnect, ErrNotConnected).
		WithDetails(map[string]any{
			"server": testNATSServer,
			"retry":  3,
		})

	assertions.NotNil(err.Details, "Details should not be nil")
	assertions.Equal(testNATSServer, err.Details["server"], "Server detail should match")
	assertions.Equal(3, err.Details["retry"], "Retry detail should match")
}

/* ======== TENANT ERROR TESTS ======== */

func TestTenantError(t *testing.T) {
	t.Run("with operation", func(t *testing.T) {
		assertions := assert.New(t)
		err := NewTenantError(testTenantID, testOpDisconnect, ErrTenantDisconnected)

		assertions.Contains(err.Error(), testTenantID, "Error should contain tenant ID")
		assertions.Contains(err.Error(), testOpDisconnect, "Error must contain operation")
		assertions.True(errors.Is(err, ErrTenantDisconnected), "Unwrap must return underlying error")
	})

	t.Run("without operation", func(t *testing.T) {
		assertions := assert.New(t)
		err := NewTenantError(testTenantID, "", ErrTenantDisconnected)

		assertions.Contains(err.Error(), testTenantID)
		assertions.Contains(err.Error(), "tenant disconnected")
		assertions.NotContains(err.Error(), ": :")
	})
}

/* ======== CONFIG ERROR TESTS ======== */

func TestConfigError(t *testing.T) {
	assertions := assert.New(t)

	err := ConfigError{
		Field:   testFieldServers,
		Message: testMsgAtLeastOne,
	}

	assertions.Contains(err.Error(), testFieldServers, "Error should contain field")
	assertions.Contains(err.Error(), testMsgAtLeastOne, "Error should contain message")
}

/* ======== SERIALIZATION ERROR TESTS ======== */

func TestSerializationError(t *testing.T) {
	assertions := assert.New(t)

	underlying := errors.New(testMsgJsonErr)
	err := SerializationError{
		Operation: testOpDeserialize,
		Type:      testTypeMyStruct,
		Err:       underlying,
	}

	assertions.Contains(err.Error(), testOpDeserialize, "Error contains operation")
	assertions.Contains(err.Error(), testTypeMyStruct, "Error should contain type")
	assertions.True(errors.Is(err, underlying), "Unwrap returns underlying error")
}

/* ======== ERROR CLASSIFICATION TESTS ======== */

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
		{"ErrInvalidSubject", ErrInvalidSubject, false},
		{"ErrInvalidMessage", ErrInvalidMessage, false},
		{"ErrInvalidConfig", ErrInvalidConfig, false},
		{"ErrMissingConfig", ErrMissingConfig, false},
		{"ErrSerializationFailed", ErrSerializationFailed, false},
		{"ErrPermissionDenied", ErrPermissionDenied, false},
		{"ErrDuplicateMsg", ErrDuplicateMsg, false},
		{"ErrShutdownInProgress", ErrShutdownInProgress, false},
		{"ErrConnectionClosed", ErrConnectionClosed, false},
		{"ErrPublishFailed (retryable)", ErrPublishFailed, true},
		{"ErrNoAck (retryable)", ErrNoAck, true},
		{"ErrTimeout (retryable)", ErrTimeout, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, IsRetryable(tc.err), "IsRetryable for %s", tc.name)
		})
	}
}

func TestIsRetryableSerializationError(t *testing.T) {
	assertions := assert.New(t)

	err := SerializationError{
		Operation: testOpSerialize,
		Type:      testTypeMyStruct,
		Err:       errors.New(testMsgTestErr),
	}

	assertions.False(IsRetryable(err), "SerializationError should not be retryable")
}

func TestIsRetryableConfigError(t *testing.T) {
	assertions := assert.New(t)

	err := ConfigError{
		Field:   testFieldServers,
		Message: testMsgRequired,
	}

	assertions.False(IsRetryable(err), "ConfigError should not be retryable")
}

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
		{"ErrNotConnected", ErrNotConnected, false},
		{"ErrInvalidSubject", ErrInvalidSubject, false},
		{"ErrPermissionDenied", ErrPermissionDenied, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, IsTemporary(tc.err), "IsTemporary for %s", tc.name)
		})
	}
}

/* ======== MULTI-ERROR TESTS ======== */

func TestMultiError(t *testing.T) {
	assertions := assert.New(t)

	err1 := errors.New(testMsgErr1)
	err2 := errors.New(testMsgErr2)
	err3 := errors.New(testMsgErr3)

	multiErr := &MultiError{
		Errors: []error{err1, err2, err3},
	}

	assertions.Contains(multiErr.Error(), "3 errors occurred", "Should indicate number of errors")
	assertions.Equal(err1, multiErr.Unwrap(), "Unwrap should return first error")
	assertions.Len(multiErr.All(), 3, "All should return all errors")
}

func TestMultiErrorSingle(t *testing.T) {
	assertions := assert.New(t)

	err1 := errors.New(testMsgSingleErr)
	multiErr := &MultiError{
		Errors: []error{err1},
	}

	assertions.Equal(testMsgSingleErr, multiErr.Error(), "Single error should return its message")
}

func TestMultiErrorEmpty(t *testing.T) {
	assertions := assert.New(t)

	multiErr := &MultiError{
		Errors: []error{},
	}

	assertions.Equal(testMsgNoErrors, multiErr.Error(), "Empty MultiError should say no errors")
	assertions.Nil(multiErr.Unwrap(), "Unwrap on empty should return nil")
}

func TestNewMultiError(t *testing.T) {
	assertions := assert.New(t)

	err1 := errors.New(testMsgErr1)
	err2 := errors.New(testMsgErr2)

	multiErr := NewMultiError([]error{err1, nil, err2, nil})
	assertions.NotNil(multiErr, "NewMultiError should not be nil")
	assertions.Len(multiErr.Errors, 2, "Should filter out nil errors")

	nilMultiErr := NewMultiError([]error{nil, nil})
	assertions.Nil(nilMultiErr, "NewMultiError with all nils should return nil")

	emptyMultiErr := NewMultiError([]error{})
	assertions.Nil(emptyMultiErr, "NewMultiError with empty slice should return nil")
}

/* ======== COMBINED ERROR UNWRAP TEST ======== */

func TestCombinedErrorUnwrap(t *testing.T) {
	assertions := assert.New(t)
	ec := NewErrorCollector()
	ec.Add(ErrNotConnected)
	ec.Add(ErrPublishFailed)

	combined := ec.Error()
	assertions.NotNil(combined)

	// The combined error should support errors.Is for both errors
	assertions.True(errors.Is(combined, ErrNotConnected))
	assertions.True(errors.Is(combined, ErrPublishFailed))
}

func TestCombinedErrorEmpty(t *testing.T) {
	assertions := assert.New(t)
	ec := NewErrorCollector()
	ec.Add(errors.New(testMsgErr1))
	ec.Add(errors.New(testMsgErr2))
	ec.Add(errors.New(testMsgErr3))

	combined := ec.Error()
	assertions.Contains(combined.Error(), "[1]")
	assertions.Contains(combined.Error(), "[2]")
	assertions.Contains(combined.Error(), "[3]")
}
