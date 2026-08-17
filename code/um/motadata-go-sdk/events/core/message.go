package core

import (
	"time"
)

// Message represents a message to be published or received
type Message struct {
	// Subject is the destination subject/topic
	Subject string

	// Data is the message payload
	Data []byte

	// Headers contains message metadata
	Headers Headers

	// Reply is the reply subject for request-reply patterns
	Reply string

	// Timestamp when the message was created
	Timestamp time.Time

	// ID is a unique message identifier (for deduplication)
	ID string
}

// NewMessage creates a new message with the given data
func NewMessage(data []byte) *Message {
	return &Message{
		Data:      data,
		Headers:   make(Headers),
		Timestamp: time.Now(),
	}
}

// WithSubject sets the message subject
func (message *Message) WithSubject(subject string) *Message {
	message.Subject = subject
	return message
}

// WithHeader adds a header to the message
func (message *Message) WithHeader(key, value string) *Message {
	if message.Headers == nil {
		message.Headers = make(Headers)
	}
	message.Headers.Set(key, value)
	return message
}

// WithHeaders sets multiple headers
func (message *Message) WithHeaders(headers Headers) *Message {
	if message.Headers == nil {
		message.Headers = make(Headers)
	}
	for headerKey, headerValues := range headers {
		message.Headers[headerKey] = headerValues
	}
	return message
}

// WithID sets the message ID for deduplication
func (message *Message) WithID(messageID string) *Message {
	message.ID = messageID
	return message
}

// WithReply sets the reply subject
func (message *Message) WithReply(replySubject string) *Message {
	message.Reply = replySubject
	return message
}

// Headers represents message headers/metadata
type Headers map[string][]string

// Get returns the first value for the given key
func (headers Headers) Get(key string) string {
	if headerValues, exists := headers[key]; exists && len(headerValues) > 0 {
		return headerValues[0]
	}
	return ""
}

// Set sets a header value, replacing any existing values
func (headers Headers) Set(key, value string) {
	headers[key] = []string{value}
}

// Add adds a value to the header (supports multiple values)
func (headers Headers) Add(key, value string) {
	headers[key] = append(headers[key], value)
}

// Del removes a header
func (headers Headers) Del(key string) {
	delete(headers, key)
}

// Has returns true if the header exists
func (headers Headers) Has(key string) bool {
	_, exists := headers[key]
	return exists
}

// Values returns all values for a key
func (headers Headers) Values(key string) []string {
	return headers[key]
}

// Clone creates a copy of the headers
func (headers Headers) Clone() Headers {
	if headers == nil {
		return nil
	}
	clonedHeaders := make(Headers, len(headers))
	for headerKey, headerValues := range headers {
		clonedHeaders[headerKey] = append([]string{}, headerValues...)
	}
	return clonedHeaders
}

// Standard header keys
const (
	HeaderContentType   = "Content-Type"
	HeaderMessageID     = "Nats-Msg-Id"
	HeaderCorrelationID = "Correlation-ID"
	HeaderTenantID      = "X-Tenant-ID"
	HeaderTraceID       = "X-Trace-ID"
	HeaderSpanID        = "X-Span-ID"
	HeaderTraceParent   = "traceparent"
	HeaderTraceState    = "tracestate"
	HeaderB3TraceID     = "X-B3-TraceId"
	HeaderB3SpanID      = "X-B3-SpanId"
	HeaderB3Sampled     = "X-B3-Sampled"
)

// Metadata contains additional message context
type Metadata struct {
	// TenantID identifies the tenant
	TenantID string

	// TraceContext for distributed tracing
	TraceContext *TraceContext

	// Timestamp when message was published
	PublishedAt time.Time

	// Sequence number (from stream)
	Sequence uint64

	// Stream name (for JetStream)
	Stream string

	// Consumer name (for JetStream)
	Consumer string

	// DeliveryCount for redelivered messages
	DeliveryCount int

	// Custom metadata
	Custom map[string]string
}

// TraceContext holds distributed tracing information
type TraceContext struct {
	TraceID  string
	SpanID   string
	ParentID string
	Sampled  bool
	State    string
}
