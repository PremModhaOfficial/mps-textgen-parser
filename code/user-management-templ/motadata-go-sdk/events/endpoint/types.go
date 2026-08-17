// Package endpoint provides endpoint registration, discovery, and lifecycle management.
//
// This package enables microservices to register their endpoints with a central
// registry and discover other service endpoints for inter-service communication.
//
// Key concepts:
//   - Endpoint: A specific API endpoint exposed by a service (e.g., /api/v1/users)
//   - EndpointInfo: Metadata about an endpoint including health status
//   - Registry: Thread-safe in-memory storage for endpoint information
//   - Manager: Handles endpoint lifecycle, health checking, and event propagation
package endpoint

import (
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
)

// Status represents the current lifecycle state of an endpoint.
// Endpoints transition through these states based on health check results
// and administrative actions. The state determines whether the endpoint
// is eligible to receive traffic (see IsAvailable).
type Status int

const (
	// StatusUnknown indicates the endpoint status is not determined
	StatusUnknown Status = iota
	// StatusHealthy indicates the endpoint is responding correctly
	StatusHealthy
	// StatusUnhealthy indicates the endpoint is not responding correctly
	StatusUnhealthy
	// StatusDegraded indicates the endpoint is partially functional
	StatusDegraded
	// StatusMaintenance indicates the endpoint is under maintenance
	StatusMaintenance
	// StatusOffline indicates the endpoint is intentionally offline
	StatusOffline
)

// String returns a human-readable status name suitable for logging and serialization.
func (s Status) String() string {
	switch s {
	case StatusUnknown:
		return "unknown"
	case StatusHealthy:
		return "healthy"
	case StatusUnhealthy:
		return "unhealthy"
	case StatusDegraded:
		return "degraded"
	case StatusMaintenance:
		return "maintenance"
	case StatusOffline:
		return "offline"
	default:
		return "unknown"
	}
}

// IsAvailable returns true if the endpoint can accept requests.
// Both StatusHealthy and StatusDegraded are considered available because
// a degraded endpoint is still partially functional and can serve traffic,
// unlike unhealthy, maintenance, or offline endpoints.
func (s Status) IsAvailable() bool {
	return s == StatusHealthy || s == StatusDegraded
}

// Protocol represents the communication protocol of an endpoint.
// It determines how clients should connect to the endpoint and is used
// when constructing the endpoint URL. The protocol also influences which
// default port is assumed (e.g., 80 for HTTP, 443 for HTTPS).
type Protocol string

const (
	// ProtocolHTTP represents plain-text HTTP connections.
	ProtocolHTTP Protocol = "http"
	// ProtocolHTTPS represents TLS-encrypted HTTP connections.
	ProtocolHTTPS Protocol = "https"
	// ProtocolGRPC represents gRPC connections (typically over HTTP/2).
	ProtocolGRPC Protocol = "grpc"
	// ProtocolWS represents plain-text WebSocket connections.
	ProtocolWS Protocol = "ws"
	// ProtocolWSS represents TLS-encrypted WebSocket connections.
	ProtocolWSS Protocol = "wss"
	// ProtocolTCP represents raw TCP connections.
	ProtocolTCP Protocol = "tcp"
	// ProtocolUDP represents UDP connections.
	ProtocolUDP Protocol = "udp"
)

// Method represents HTTP/gRPC method types supported by an endpoint.
// It is used to declare which HTTP methods an endpoint accepts, enabling
// clients to validate requests before sending them.
type Method string

const (
	// MethodGET represents the HTTP GET method for retrieving resources.
	MethodGET Method = "GET"
	// MethodPOST represents the HTTP POST method for creating resources.
	MethodPOST Method = "POST"
	// MethodPUT represents the HTTP PUT method for full resource replacement.
	MethodPUT Method = "PUT"
	// MethodDELETE represents the HTTP DELETE method for removing resources.
	MethodDELETE Method = "DELETE"
	// MethodPATCH represents the HTTP PATCH method for partial resource updates.
	MethodPATCH Method = "PATCH"
	// MethodHEAD represents the HTTP HEAD method for retrieving headers only.
	MethodHEAD Method = "HEAD"
	// MethodOPTIONS represents the HTTP OPTIONS method for CORS preflight and capability discovery.
	MethodOPTIONS Method = "OPTIONS"
	// MethodAny is a wildcard indicating the endpoint accepts all HTTP methods.
	MethodAny Method = "*"
)

// Info holds comprehensive endpoint information including its network location,
// health status, routing metadata, versioning, and multi-tenancy context.
// Info is the primary data structure exchanged during endpoint registration,
// discovery, and health checking. It is designed to be cloned for thread-safe
// sharing across goroutines.
type Info struct {
	// Core identification
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	// Service association
	ServiceID   string `json:"service_id"`
	ServiceName string `json:"service_name"`
	InstanceID  string `json:"instance_id,omitempty"`

	// Network information
	Host     string   `json:"host"`
	Port     int      `json:"port"`
	Path     string   `json:"path"`
	Protocol Protocol `json:"protocol"`
	Methods  []Method `json:"methods,omitempty"`

	// Health and status
	Status           Status    `json:"status"`
	LastHealthCheck  time.Time `json:"last_health_check,omitempty"`
	HealthCheckPath  string    `json:"health_check_path,omitempty"`
	HealthCheckPort  int       `json:"health_check_port,omitempty"`
	ConsecutiveFails int       `json:"consecutive_fails,omitempty"`

	// Routing and load balancing
	Weight   int               `json:"weight,omitempty"`
	Priority int               `json:"priority,omitempty"`
	Tags     []string          `json:"tags,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`

	// Versioning
	Version    string `json:"version,omitempty"`
	APIVersion string `json:"api_version,omitempty"`

	// Lifecycle
	RegisteredAt time.Time `json:"registered_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`

	// Multi-tenancy
	TenantID string `json:"tenant_id,omitempty"`
	Region   string `json:"region,omitempty"`
	Zone     string `json:"zone,omitempty"`
}

// Clone creates a deep copy of the endpoint info, duplicating all slices and maps
// so that mutations to the clone do not affect the original. This is essential for
// thread-safe sharing between the registry's internal state and callers.
func (e *Info) Clone() *Info {
	if e == nil {
		return nil
	}

	clone := *e

	// Deep copy slices
	if e.Methods != nil {
		clone.Methods = make([]Method, len(e.Methods))
		copy(clone.Methods, e.Methods)
	}
	if e.Tags != nil {
		clone.Tags = make([]string, len(e.Tags))
		copy(clone.Tags, e.Tags)
	}
	if e.Metadata != nil {
		clone.Metadata = make(map[string]string, len(e.Metadata))
		for k, v := range e.Metadata {
			clone.Metadata[k] = v
		}
	}

	return &clone
}

// URL returns the full URL for the endpoint by combining protocol, host, port, and path.
// Standard ports (80 for HTTP, 443 for HTTPS) and zero ports are omitted from the URL
// to produce clean, canonical URLs.
func (e *Info) URL() string {
	if e.Port == 0 || (e.Protocol == ProtocolHTTP && e.Port == 80) || (e.Protocol == ProtocolHTTPS && e.Port == 443) {
		return string(e.Protocol) + "://" + e.Host + e.Path
	}
	return string(e.Protocol) + "://" + e.Host + ":" + itoa(e.Port) + e.Path
}

// Address returns the host:port address string suitable for net.Dial or similar use.
// If port is zero, only the host is returned.
func (e *Info) Address() string {
	if e.Port == 0 {
		return e.Host
	}
	return e.Host + ":" + itoa(e.Port)
}

// IsExpired returns true if the endpoint has exceeded its TTL-based expiration time.
// Endpoints with a zero ExpiresAt never expire.
func (e *Info) IsExpired() bool {
	if e.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(e.ExpiresAt)
}

// IsHealthy returns true if the endpoint is in a healthy state.
// Unlike IsAvailable on Status, this strictly requires StatusHealthy and
// does not include degraded endpoints.
func (e *Info) IsHealthy() bool {
	return e.Status == StatusHealthy
}

// HasTag checks if the endpoint has a specific tag by performing a linear scan.
// Tags are used for filtering during service discovery queries.
func (e *Info) HasTag(tag string) bool {
	for _, t := range e.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// GetMetadata returns a metadata value by key, or an empty string if the key
// is not present or the metadata map is nil.
func (e *Info) GetMetadata(key string) string {
	if e.Metadata == nil {
		return ""
	}
	return e.Metadata[key]
}

// SetMetadata sets a metadata key-value pair on the endpoint, lazily initializing
// the metadata map if it has not been created yet.
func (e *Info) SetMetadata(key, value string) {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	e.Metadata[key] = value
}

// RegistrationRequest represents a request to register an endpoint with the registry.
// It bundles the endpoint metadata, an optional TTL for automatic expiration, and
// an override flag to control behavior when the endpoint already exists.
type RegistrationRequest struct {
	// Endpoint information
	Endpoint *Info `json:"endpoint"`

	// TTL for automatic expiration (0 means no expiration)
	TTL time.Duration `json:"ttl,omitempty"`

	// Override existing registration if present
	Override bool `json:"override,omitempty"`
}

// Validate validates the registration request by checking that the endpoint is non-nil
// and that required fields (ID, ServiceID, Host) are present. This should be called
// before passing the request to the Manager or Registry to get clear error messages.
func (r *RegistrationRequest) Validate() error {
	if r.Endpoint == nil {
		return ErrInvalidEndpoint
	}
	if r.Endpoint.ID == "" {
		return ErrMissingEndpointID
	}
	if r.Endpoint.ServiceID == "" {
		return ErrMissingServiceID
	}
	if r.Endpoint.Host == "" {
		return ErrMissingHost
	}
	return nil
}

// DeregistrationRequest represents a request to remove an endpoint from the registry.
// It identifies the endpoint by ID and optionally includes a reason for deregistration,
// which is propagated through lifecycle callbacks and events.
type DeregistrationRequest struct {
	EndpointID string `json:"endpoint_id"`
	ServiceID  string `json:"service_id,omitempty"`
	InstanceID string `json:"instance_id,omitempty"`
	Reason     string `json:"reason,omitempty"`
}

// Validate validates the deregistration request by ensuring the endpoint ID is provided.
func (r *DeregistrationRequest) Validate() error {
	if r.EndpointID == "" {
		return ErrMissingEndpointID
	}
	return nil
}

// HealthUpdate represents a health status update for an endpoint, typically
// produced by the HealthChecker or reported by external monitoring systems.
// It carries the new status, a human-readable message, and optional details.
type HealthUpdate struct {
	EndpointID string    `json:"endpoint_id"`
	Status     Status    `json:"status"`
	Message    string    `json:"message,omitempty"`
	CheckedAt  time.Time `json:"checked_at"`
	Details    any       `json:"details,omitempty"`
}

// Query defines filters for searching endpoints in the registry.
// All non-zero fields act as AND filters -- an endpoint must satisfy every
// specified criterion to be included in results. The query also supports
// pagination via Limit and Offset fields.
type Query struct {
	// Filter by IDs
	IDs        []string `json:"ids,omitempty"`
	ServiceIDs []string `json:"service_ids,omitempty"`
	InstanceID string   `json:"instance_id,omitempty"`

	// Filter by network
	Host     string   `json:"host,omitempty"`
	Port     int      `json:"port,omitempty"`
	Protocol Protocol `json:"protocol,omitempty"`

	// Filter by status
	Statuses      []Status `json:"statuses,omitempty"`
	OnlyHealthy   bool     `json:"only_healthy,omitempty"`
	OnlyAvailable bool     `json:"only_available,omitempty"`

	// Filter by tags/metadata
	Tags     []string          `json:"tags,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`

	// Filter by version
	Version    string `json:"version,omitempty"`
	APIVersion string `json:"api_version,omitempty"`

	// Multi-tenancy filters
	TenantID string `json:"tenant_id,omitempty"`
	Region   string `json:"region,omitempty"`
	Zone     string `json:"zone,omitempty"`

	// Pagination
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

// Matches checks if an endpoint satisfies all non-zero query filters.
// Each filter category (ID, network, status, tags, metadata, version, tenancy)
// is evaluated independently; the endpoint must pass all of them.
// A nil endpoint always returns false.
func (q *Query) Matches(e *Info) bool {
	if e == nil {
		return false
	}

	return q.matchesIdentity(e) &&
		q.matchesNetwork(e) &&
		q.matchesStatus(e) &&
		q.matchesTags(e) &&
		q.matchesMetadata(e) &&
		q.matchesVersion(e) &&
		q.matchesTenancy(e)
}

func (q *Query) matchesIdentity(e *Info) bool {
	if len(q.IDs) > 0 && !containsString(q.IDs, e.ID) {
		return false
	}
	if len(q.ServiceIDs) > 0 && !containsString(q.ServiceIDs, e.ServiceID) {
		return false
	}
	return q.InstanceID == "" || e.InstanceID == q.InstanceID
}

func (q *Query) matchesNetwork(e *Info) bool {
	if q.Host != "" && e.Host != q.Host {
		return false
	}
	if q.Port != 0 && e.Port != q.Port {
		return false
	}
	return q.Protocol == "" || e.Protocol == q.Protocol
}

func (q *Query) matchesStatus(e *Info) bool {
	if len(q.Statuses) > 0 && !containsStatus(q.Statuses, e.Status) {
		return false
	}
	if q.OnlyHealthy && e.Status != StatusHealthy {
		return false
	}
	return !q.OnlyAvailable || e.Status.IsAvailable()
}

func (q *Query) matchesTags(e *Info) bool {
	for _, tag := range q.Tags {
		if !e.HasTag(tag) {
			return false
		}
	}
	return true
}

func (q *Query) matchesMetadata(e *Info) bool {
	for k, v := range q.Metadata {
		if e.GetMetadata(k) != v {
			return false
		}
	}
	return true
}

func (q *Query) matchesVersion(e *Info) bool {
	if q.Version != "" && e.Version != q.Version {
		return false
	}
	return q.APIVersion == "" || e.APIVersion == q.APIVersion
}

func (q *Query) matchesTenancy(e *Info) bool {
	if q.TenantID != "" && e.TenantID != q.TenantID {
		return false
	}
	if q.Region != "" && e.Region != q.Region {
		return false
	}
	return q.Zone == "" || e.Zone == q.Zone
}

// Event types for endpoint lifecycle
type EventType string

const (
	EventRegistered   EventType = "endpoint.registered"
	EventDeregistered EventType = "endpoint.deregistered"
	EventUpdated      EventType = "endpoint.updated"
	EventHealthy      EventType = "endpoint.healthy"
	EventUnhealthy    EventType = "endpoint.unhealthy"
	EventExpired      EventType = "endpoint.expired"
)

// Event represents an endpoint lifecycle event
type Event struct {
	Type       EventType `json:"type"`
	EndpointID string    `json:"endpoint_id"`
	ServiceID  string    `json:"service_id"`
	Endpoint   *Info     `json:"endpoint,omitempty"`
	OldStatus  Status    `json:"old_status,omitempty"`
	NewStatus  Status    `json:"new_status,omitempty"`
	Reason     string    `json:"reason,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

// ToMessage converts the event to a *nats.Msg
func (e *Event) ToMessage() (*nats.Msg, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	msg := &nats.Msg{Data: data, Header: make(nats.Header)}
	msg.Header.Set("event-type", string(e.Type))
	msg.Header.Set("endpoint-id", e.EndpointID)
	msg.Header.Set("service-id", e.ServiceID)
	return msg, nil
}

// EventFromMessage creates an Event from a *nats.Msg
func EventFromMessage(msg *nats.Msg) (*Event, error) {
	var event Event
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

// Helper functions

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func containsStatus(slice []Status, s Status) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	if i < 0 {
		return "-" + itoa(-i)
	}
	var b [20]byte
	bp := len(b) - 1
	for i > 0 {
		b[bp] = byte('0' + i%10)
		bp--
		i /= 10
	}
	return string(b[bp+1:])
}
