package core

import (
	"context"
	"encoding/json"
	"sync"
	"time"
)

// ============================================================================
// Request Interface - Production-Level Request-Reply Pattern
// ============================================================================

// Request represents an incoming request in a request-reply pattern.
// It provides methods to access request data and send responses.
type Request interface {
	// Context returns the request context
	Context() context.Context

	// Subject returns the request subject
	Subject() string

	// Reply returns the reply subject
	Reply() string

	// Data returns the raw request payload
	Data() []byte

	// Headers returns the request headers
	Headers() Headers

	// Header returns a specific header value
	Header(key string) string

	// Respond sends a response with raw data
	Respond(data []byte, opts ...RespondOption) error

	// RespondJSON sends a JSON-encoded response
	RespondJSON(v any, opts ...RespondOption) error

	// RespondError sends an error response
	RespondError(code, description string, data []byte, opts ...RespondOption) error

	// Ack acknowledges the request (for JetStream)
	Ack() error

	// Nak negative-acknowledges the request (for JetStream)
	Nak() error

	// NakWithDelay negative-acknowledges with redelivery delay
	NakWithDelay(delay time.Duration) error

	// Term terminates the message (no more redeliveries)
	Term() error

	// InProgress signals the message is being processed
	InProgress() error
}

// RespondOption configures a response
type RespondOption func(*respondConfig)

type respondConfig struct {
	headers Headers
}

// WithResponseHeaders adds headers to the response
func WithResponseHeaders(headers Headers) RespondOption {
	return func(c *respondConfig) {
		c.headers = headers
	}
}

// WithResponseHeader adds a single header to the response
func WithResponseHeader(key, value string) RespondOption {
	return func(c *respondConfig) {
		if c.headers == nil {
			c.headers = make(Headers)
		}
		c.headers.Set(key, value)
	}
}

// ============================================================================
// Service Error Response
// ============================================================================

// ServiceError represents an error response in the NATS micro protocol format
type ServiceError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Data        []byte `json:"data,omitempty"`
}

// Error implements the error interface
func (e *ServiceError) Error() string {
	return e.Description
}

// Standard error header keys (compatible with NATS micro protocol)
const (
	HeaderServiceError     = "Nats-Service-Error"
	HeaderServiceErrorCode = "Nats-Service-Error-Code"
)

// ============================================================================
// Request Handler Types
// ============================================================================

// RequestHandler handles incoming requests
type RequestHandler func(Request)

// RequestHandlerFunc is a function adapter for RequestHandler
type RequestHandlerFunc func(Request)

// Handle implements RequestHandler
func (f RequestHandlerFunc) Handle(r Request) {
	f(r)
}

// ContextRequestHandler handles requests with context support
type ContextRequestHandler func(ctx context.Context, req Request)

// ============================================================================
// Service Definition
// ============================================================================

// ServiceConfig defines a service for request-reply patterns
type ServiceConfig struct {
	// Name is the service name (required)
	Name string

	// Version is the service version (e.g., "1.0.0")
	Version string

	// Description describes the service
	Description string

	// Metadata contains service metadata
	Metadata map[string]string

	// QueueGroup for load balancing (default: "q")
	QueueGroup string

	// Endpoint is the default endpoint configuration
	Endpoint *EndpointConfig

	// DoneHandler is called when the service stops
	DoneHandler func(Service)

	// ErrorHandler is called on service errors
	ErrorHandler func(Service, error)

	// StatsHandler returns custom stats
	StatsHandler func(*Endpoint) any
}

// EndpointConfig defines an endpoint within a service
type EndpointConfig struct {
	// Subject is the endpoint subject (required if no Name)
	Subject string

	// Name is the endpoint name (used to generate subject if Subject is empty)
	Name string

	// Handler handles requests
	Handler RequestHandler

	// Metadata contains endpoint metadata
	Metadata map[string]string

	// QueueGroup overrides service-level queue group
	QueueGroup string

	// QueueGroupDisabled disables queue groups for fanout
	QueueGroupDisabled bool
}

// ============================================================================
// Service Interface
// ============================================================================

// Service represents a running service
type Service interface {
	// Info returns service information
	Info() ServiceInfo

	// Stats returns service statistics
	Stats() ServiceStats

	// Reset resets service statistics
	Reset()

	// Stop stops the service
	Stop() error

	// Stopped returns a channel that closes when the service stops
	Stopped() <-chan struct{}

	// AddEndpoint adds an endpoint to the service
	AddEndpoint(name string, handler RequestHandler, opts ...EndpointOption) error

	// AddGroup creates an endpoint group with a subject prefix
	AddGroup(name string, opts ...GroupOption) Group
}

// Group represents an endpoint group with shared configuration
type Group interface {
	// AddEndpoint adds an endpoint to the group
	AddEndpoint(name string, handler RequestHandler, opts ...EndpointOption) error

	// AddGroup creates a nested group
	AddGroup(name string, opts ...GroupOption) Group
}

// ============================================================================
// Service Information
// ============================================================================

// ServiceInfo contains service metadata
type ServiceInfo struct {
	Type        string            `json:"type"`
	Name        string            `json:"name"`
	ID          string            `json:"id"`
	Version     string            `json:"version"`
	Description string            `json:"description,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Endpoints   []EndpointInfo    `json:"endpoints,omitempty"`
}

// EndpointInfo contains endpoint metadata
type EndpointInfo struct {
	Name       string            `json:"name"`
	Subject    string            `json:"subject"`
	QueueGroup string            `json:"queue_group,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// ============================================================================
// Service Statistics
// ============================================================================

// ServiceStats contains service statistics
type ServiceStats struct {
	Name      string          `json:"name"`
	ID        string          `json:"id"`
	Version   string          `json:"version"`
	Started   time.Time       `json:"started"`
	Endpoints []EndpointStats `json:"endpoints,omitempty"`
}

// EndpointStats contains endpoint statistics
type EndpointStats struct {
	Name                  string        `json:"name"`
	Subject               string        `json:"subject"`
	QueueGroup            string        `json:"queue_group,omitempty"`
	NumRequests           int64         `json:"num_requests"`
	NumErrors             int64         `json:"num_errors"`
	TotalProcessingTime   time.Duration `json:"processing_time"`
	AverageProcessingTime time.Duration `json:"average_processing_time"`
	LastError             string        `json:"last_error,omitempty"`
	LastErrorTime         time.Time     `json:"last_error_time,omitempty"`
	Data                  any           `json:"data,omitempty"`
}

// ============================================================================
// Endpoint
// ============================================================================

// Endpoint represents a service endpoint
type Endpoint struct {
	mu sync.RWMutex

	Name       string
	Subject    string
	QueueGroup string
	Handler    RequestHandler
	Metadata   map[string]string

	// Statistics
	numRequests         int64
	numErrors           int64
	totalProcessingTime time.Duration
	lastError           string
	lastErrorTime       time.Time

	// Custom stats handler
	statsHandler func(*Endpoint) any
}

// RecordRequest records a request
func (e *Endpoint) RecordRequest(duration time.Duration, err error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.numRequests++
	e.totalProcessingTime += duration

	if err != nil {
		e.numErrors++
		e.lastError = err.Error()
		e.lastErrorTime = time.Now()
	}
}

// Stats returns endpoint statistics
func (e *Endpoint) Stats() EndpointStats {
	e.mu.RLock()
	defer e.mu.RUnlock()

	var avg time.Duration
	if e.numRequests > 0 {
		avg = e.totalProcessingTime / time.Duration(e.numRequests)
	}

	stats := EndpointStats{
		Name:                  e.Name,
		Subject:               e.Subject,
		QueueGroup:            e.QueueGroup,
		NumRequests:           e.numRequests,
		NumErrors:             e.numErrors,
		TotalProcessingTime:   e.totalProcessingTime,
		AverageProcessingTime: avg,
		LastError:             e.lastError,
		LastErrorTime:         e.lastErrorTime,
	}

	if e.statsHandler != nil {
		stats.Data = e.statsHandler(e)
	}

	return stats
}

// Reset resets endpoint statistics
func (e *Endpoint) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.numRequests = 0
	e.numErrors = 0
	e.totalProcessingTime = 0
	e.lastError = ""
	e.lastErrorTime = time.Time{}
}

// ============================================================================
// Options
// ============================================================================

// EndpointOption configures an endpoint
type EndpointOption func(*endpointConfig)

type endpointConfig struct {
	subject            string
	metadata           map[string]string
	queueGroup         string
	queueGroupDisabled bool
}

// WithEndpointSubject sets a custom subject for the endpoint
func WithEndpointSubject(subject string) EndpointOption {
	return func(c *endpointConfig) {
		c.subject = subject
	}
}

// WithEndpointMetadata sets endpoint metadata
func WithEndpointMetadata(metadata map[string]string) EndpointOption {
	return func(c *endpointConfig) {
		c.metadata = metadata
	}
}

// WithEndpointQueueGroup sets a custom queue group
func WithEndpointQueueGroup(queueGroup string) EndpointOption {
	return func(c *endpointConfig) {
		c.queueGroup = queueGroup
	}
}

// WithEndpointQueueGroupDisabled disables queue groups
func WithEndpointQueueGroupDisabled() EndpointOption {
	return func(c *endpointConfig) {
		c.queueGroupDisabled = true
	}
}

// GroupOption configures a group
type GroupOption func(*groupConfig)

type groupConfig struct {
	queueGroup         string
	queueGroupDisabled bool
}

// WithGroupQueueGroup sets a custom queue group for the group
func WithGroupQueueGroup(queueGroup string) GroupOption {
	return func(c *groupConfig) {
		c.queueGroup = queueGroup
	}
}

// WithGroupQueueGroupDisabled disables queue groups for the group
func WithGroupQueueGroupDisabled() GroupOption {
	return func(c *groupConfig) {
		c.queueGroupDisabled = true
	}
}

// ============================================================================
// Requester Interface - Client Side
// ============================================================================

// Requester sends requests and receives responses
type Requester interface {
	// Request sends a request and waits for a response
	Request(ctx context.Context, subject string, data []byte, opts ...RequestOption) (*Response, error)

	// RequestJSON sends a JSON request and decodes the JSON response
	RequestJSON(ctx context.Context, subject string, req any, resp any, opts ...RequestOption) error

	// RequestMulti sends a request and collects multiple responses
	RequestMulti(ctx context.Context, subject string, data []byte, opts ...RequestOption) ([]*Response, error)

	// Close closes the requester
	Close() error
}

// RequestOption configures a request
type RequestOption func(*requestConfig)

type requestConfig struct {
	headers     Headers
	timeout     time.Duration
	maxReplies  int
	noMux       bool
	oldStyle    bool
	expectedRTT time.Duration
}

// WithRequestHeaders sets request headers
func WithRequestHeaders(headers Headers) RequestOption {
	return func(c *requestConfig) {
		c.headers = headers
	}
}

// WithRequestTimeout sets request timeout
func WithRequestTimeout(timeout time.Duration) RequestOption {
	return func(c *requestConfig) {
		c.timeout = timeout
	}
}

// WithMaxReplies sets maximum number of replies to collect
func WithMaxReplies(max int) RequestOption {
	return func(c *requestConfig) {
		c.maxReplies = max
	}
}

// WithExpectedRTT sets expected round-trip time for scatter-gather
func WithExpectedRTT(rtt time.Duration) RequestOption {
	return func(c *requestConfig) {
		c.expectedRTT = rtt
	}
}

// Response represents a response to a request
type Response struct {
	Data    []byte
	Headers Headers
	Error   *ServiceError
}

// IsError returns true if the response is an error
func (r *Response) IsError() bool {
	return r.Error != nil
}

// JSON decodes the response data as JSON
func (r *Response) JSON(v any) error {
	if r.Error != nil {
		return r.Error
	}
	return json.Unmarshal(r.Data, v)
}

// ============================================================================
// Control Subjects (NATS Micro Protocol Compatible)
// ============================================================================

// Verb represents a control verb
type Verb int

const (
	// PingVerb for service discovery
	PingVerb Verb = iota
	// StatsVerb for service statistics
	StatsVerb
	// InfoVerb for service information
	InfoVerb
)

// String returns the string representation
func (v Verb) String() string {
	switch v {
	case PingVerb:
		return "PING"
	case StatsVerb:
		return "STATS"
	case InfoVerb:
		return "INFO"
	default:
		return "UNKNOWN"
	}
}

// API constants for NATS micro protocol
const (
	// APIPrefix is the prefix for service API subjects
	APIPrefix = "$SRV"

	// DefaultQueueGroup is the default queue group for services
	DefaultQueueGroup = "q"
)

// ControlSubject generates a control subject for service discovery
func ControlSubject(verb Verb, name, id string) string {
	base := APIPrefix + "." + verb.String()
	if name == "" {
		return base
	}
	if id == "" {
		return base + "." + name
	}
	return base + "." + name + "." + id
}

// ============================================================================
// Ping Response
// ============================================================================

// PingResponse is returned from PING requests
type PingResponse struct {
	Type     string            `json:"type"`
	Name     string            `json:"name"`
	ID       string            `json:"id"`
	Version  string            `json:"version"`
	Metadata map[string]string `json:"metadata,omitempty"`
}
