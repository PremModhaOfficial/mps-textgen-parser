package core

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/* ========================================================================================================
   MESSAGE TESTS - Core message operations
   ======================================================================================================== */

func TestNewMessage(t *testing.T) {

	assertions := assert.New(t)

	data := []byte("test payload")
	msg := NewMessage(data)

	assertions.NotNil(msg, "Message should not be nil")
	assertions.Equal(data, msg.Data, "Message data should match")
	assertions.NotNil(msg.Headers, "Headers should be initialized")
	assertions.False(msg.Timestamp.IsZero(), "Timestamp should be set")
}

func TestMessageWithSubject(t *testing.T) {

	assertions := assert.New(t)

	msg := NewMessage([]byte("data")).WithSubject("test.subject")

	assertions.Equal("test.subject", msg.Subject, "Subject should be set")
}

func TestMessageWithHeader(t *testing.T) {

	assertions := assert.New(t)

	msg := NewMessage([]byte("data")).
		WithHeader("key1", "value1").
		WithHeader("key2", "value2")

	assertions.Equal("value1", msg.Headers.Get("key1"), "Header key1 should be set")
	assertions.Equal("value2", msg.Headers.Get("key2"), "Header key2 should be set")
}

func TestMessageWithHeaders(t *testing.T) {

	assertions := assert.New(t)

	headers := Headers{
		"key1": []string{"value1"},
		"key2": []string{"value2", "value2b"},
	}

	msg := NewMessage([]byte("data")).WithHeaders(headers)

	assertions.Equal("value1", msg.Headers.Get("key1"), "Header key1 should be set")
	assertions.Equal("value2", msg.Headers.Get("key2"), "Header key2 should be first value")
}

func TestMessageWithID(t *testing.T) {

	assertions := assert.New(t)

	msg := NewMessage([]byte("data")).WithID("msg-123")

	assertions.Equal("msg-123", msg.ID, "Message ID should be set")
}

func TestMessageWithReply(t *testing.T) {

	assertions := assert.New(t)

	msg := NewMessage([]byte("data")).WithReply("reply.subject")

	assertions.Equal("reply.subject", msg.Reply, "Reply subject should be set")
}

func TestMessageChaining(t *testing.T) {

	assertions := assert.New(t)

	msg := NewMessage([]byte("data")).
		WithSubject("test.subject").
		WithID("msg-123").
		WithReply("reply.subject").
		WithHeader("Content-Type", "application/json")

	assertions.Equal("test.subject", msg.Subject, "Subject should be set")
	assertions.Equal("msg-123", msg.ID, "Message ID should be set")
	assertions.Equal("reply.subject", msg.Reply, "Reply should be set")
	assertions.Equal("application/json", msg.Headers.Get("Content-Type"), "Content-Type header should be set")
}

/* ========================================================================================================
   HEADERS TESTS - Header operations
   ======================================================================================================== */

func TestHeadersGet(t *testing.T) {

	assertions := assert.New(t)

	headers := Headers{
		"key1": []string{"value1", "value2"},
		"key2": []string{"value3"},
	}

	assertions.Equal("value1", headers.Get("key1"), "Get should return first value")
	assertions.Equal("value3", headers.Get("key2"), "Get should return value")
	assertions.Equal("", headers.Get("nonexistent"), "Get should return empty for missing key")
}

func TestHeadersSet(t *testing.T) {

	assertions := assert.New(t)

	headers := make(Headers)
	headers.Set("key1", "value1")
	headers.Set("key1", "value2")

	assertions.Equal("value2", headers.Get("key1"), "Set should replace existing value")
	assertions.Equal(1, len(headers["key1"]), "Set should only have one value")
}

func TestHeadersAdd(t *testing.T) {

	assertions := assert.New(t)

	headers := make(Headers)
	headers.Add("key1", "value1")
	headers.Add("key1", "value2")

	values := headers.Values("key1")
	assertions.Equal(2, len(values), "Add should append values")
	assertions.Contains(values, "value1", "Should contain value1")
	assertions.Contains(values, "value2", "Should contain value2")
}

func TestHeadersDel(t *testing.T) {

	assertions := assert.New(t)

	headers := Headers{
		"key1": []string{"value1"},
		"key2": []string{"value2"},
	}

	headers.Del("key1")

	assertions.False(headers.Has("key1"), "key1 should be deleted")
	assertions.True(headers.Has("key2"), "key2 should still exist")
}

func TestHeadersHas(t *testing.T) {

	assertions := assert.New(t)

	headers := Headers{
		"key1": []string{"value1"},
	}

	assertions.True(headers.Has("key1"), "Has should return true for existing key")
	assertions.False(headers.Has("key2"), "Has should return false for missing key")
}

func TestHeadersClone(t *testing.T) {

	assertions := assert.New(t)

	original := Headers{
		"key1": []string{"value1", "value2"},
		"key2": []string{"value3"},
	}

	cloned := original.Clone()

	assertions.Equal(original, cloned, "Cloned headers should equal original")

	// Modify clone and verify original unchanged
	cloned.Set("key1", "modified")

	assertions.Equal("value1", original.Get("key1"), "Original should be unchanged")
	assertions.Equal("modified", cloned.Get("key1"), "Clone should be modified")
}

func TestHeadersCloneNil(t *testing.T) {

	assertions := assert.New(t)

	var headers Headers
	cloned := headers.Clone()

	assertions.Nil(cloned, "Clone of nil headers should be nil")
}

/* ========================================================================================================
   METADATA TESTS
   ======================================================================================================== */

func TestTraceContext(t *testing.T) {

	assertions := assert.New(t)

	tc := &TraceContext{
		TraceID:  "trace-123",
		SpanID:   "span-456",
		ParentID: "parent-789",
		Sampled:  true,
		State:    "state-data",
	}

	assertions.Equal("trace-123", tc.TraceID, "TraceID should match")
	assertions.Equal("span-456", tc.SpanID, "SpanID should match")
	assertions.Equal("parent-789", tc.ParentID, "ParentID should match")
	assertions.True(tc.Sampled, "Sampled should be true")
	assertions.Equal("state-data", tc.State, "State should match")
}

func TestMetadata(t *testing.T) {

	assertions := assert.New(t)

	metadata := &Metadata{
		TenantID: "tenant-123",
		TraceContext: &TraceContext{
			TraceID: "trace-123",
		},
		PublishedAt:   time.Now(),
		Sequence:      42,
		Stream:        "test-stream",
		Consumer:      "test-consumer",
		DeliveryCount: 1,
		Custom: map[string]string{
			"custom-key": "custom-value",
		},
	}

	assertions.Equal("tenant-123", metadata.TenantID, "TenantID should match")
	assertions.NotNil(metadata.TraceContext, "TraceContext should not be nil")
	assertions.Equal("trace-123", metadata.TraceContext.TraceID, "TraceID should match")
	assertions.Equal(uint64(42), metadata.Sequence, "Sequence should match")
	assertions.Equal("test-stream", metadata.Stream, "Stream should match")
	assertions.Equal("custom-value", metadata.Custom["custom-key"], "Custom key should match")
}
