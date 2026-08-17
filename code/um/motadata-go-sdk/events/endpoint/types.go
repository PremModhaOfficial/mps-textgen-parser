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

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

// Status represents the current state of an endpoint
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

// String returns a human-readable status name
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

// IsAvailable returns true if the endpoint can accept requests
func (s Status) IsAvailable() bool {
	return s == StatusHealthy || s == StatusDegraded
}

// Protocol represents the communication protocol of an endpoint
type Protocol string

const (
	ProtocolHTTP  Protocol = "http"
	ProtocolHTTPS Protocol = "https"
	ProtocolGRPC  Protocol = "grpc"
	ProtocolWS    Protocol = "ws"
	ProtocolWSS   Protocol = "wss"
	ProtocolTCP   Protocol = "tcp"
	ProtocolUDP   Protocol = "udp"
)

// Method represents HTTP/gRPC method types
type Method string

const (
	MethodGET     Method = "GET"
	MethodPOST    Method = "POST"
	MethodPUT     Method = "PUT"
	MethodDELETE  Method = "DELETE"
	MethodPATCH   Method = "PATCH"
	MethodHEAD    Method = "HEAD"
	MethodOPTIONS Method = "OPTIONS"
	MethodAny     Method = "*"
)

// Info holds comprehensive endpoint information
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

// Clone creates a deep copy of the endpoint info
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

// URL returns the full URL for the endpoint
func (e *Info) URL() string {
	if e.Port == 0 || (e.Protocol == ProtocolHTTP && e.Port == 80) || (e.Protocol == ProtocolHTTPS && e.Port == 443) {
		return string(e.Protocol) + "://" + e.Host + e.Path
	}
	return string(e.Protocol) + "://" + e.Host + ":" + itoa(e.Port) + e.Path
}

// Address returns the host:port address
func (e *Info) Address() string {
	if e.Port == 0 {
		return e.Host
	}
	return e.Host + ":" + itoa(e.Port)
}

// IsExpired returns true if the endpoint has expired
func (e *Info) IsExpired() bool {
	if e.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(e.ExpiresAt)
}

// IsHealthy returns true if the endpoint is in a healthy state
func (e *Info) IsHealthy() bool {
	return e.Status == StatusHealthy
}

// HasTag checks if the endpoint has a specific tag
func (e *Info) HasTag(tag string) bool {
	for _, t := range e.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// GetMetadata returns a metadata value or empty string if not found
func (e *Info) GetMetadata(key string) string {
	if e.Metadata == nil {
		return ""
	}
	return e.Metadata[key]
}

// SetMetadata sets a metadata value
func (e *Info) SetMetadata(key, value string) {
	if e.Metadata == nil {
		e.Metadata = make(map[string]string)
	}
	e.Metadata[key] = value
}

// RegistrationRequest represents a request to register an endpoint
type RegistrationRequest struct {
	// Endpoint information
	Endpoint *Info `json:"endpoint"`

	// TTL for automatic expiration (0 means no expiration)
	TTL time.Duration `json:"ttl,omitempty"`

	// Override existing registration if present
	Override bool `json:"override,omitempty"`
}

// Validate validates the registration request
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

// DeregistrationRequest represents a request to deregister an endpoint
type DeregistrationRequest struct {
	EndpointID string `json:"endpoint_id"`
	ServiceID  string `json:"service_id,omitempty"`
	InstanceID string `json:"instance_id,omitempty"`
	Reason     string `json:"reason,omitempty"`
}

// Validate validates the deregistration request
func (r *DeregistrationRequest) Validate() error {
	if r.EndpointID == "" {
		return ErrMissingEndpointID
	}
	return nil
}

// HealthUpdate represents a health status update for an endpoint
type HealthUpdate struct {
	EndpointID string    `json:"endpoint_id"`
	Status     Status    `json:"status"`
	Message    string    `json:"message,omitempty"`
	CheckedAt  time.Time `json:"checked_at"`
	Details    any       `json:"details,omitempty"`
}

// Query defines filters for searching endpoints
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

// Matches checks if an endpoint matches the query filters
func (q *Query) Matches(e *Info) bool {
	if e == nil {
		return false
	}

	// ID filters
	if len(q.IDs) > 0 && !containsString(q.IDs, e.ID) {
		return false
	}
	if len(q.ServiceIDs) > 0 && !containsString(q.ServiceIDs, e.ServiceID) {
		return false
	}
	if q.InstanceID != "" && e.InstanceID != q.InstanceID {
		return false
	}

	// Network filters
	if q.Host != "" && e.Host != q.Host {
		return false
	}
	if q.Port != 0 && e.Port != q.Port {
		return false
	}
	if q.Protocol != "" && e.Protocol != q.Protocol {
		return false
	}

	// Status filters
	if len(q.Statuses) > 0 && !containsStatus(q.Statuses, e.Status) {
		return false
	}
	if q.OnlyHealthy && e.Status != StatusHealthy {
		return false
	}
	if q.OnlyAvailable && !e.Status.IsAvailable() {
		return false
	}

	// Tag filters (all tags must match)
	for _, tag := range q.Tags {
		if !e.HasTag(tag) {
			return false
		}
	}

	// Metadata filters (all metadata must match)
	for k, v := range q.Metadata {
		if e.GetMetadata(k) != v {
			return false
		}
	}

	// Version filters
	if q.Version != "" && e.Version != q.Version {
		return false
	}
	if q.APIVersion != "" && e.APIVersion != q.APIVersion {
		return false
	}

	// Multi-tenancy filters
	if q.TenantID != "" && e.TenantID != q.TenantID {
		return false
	}
	if q.Region != "" && e.Region != q.Region {
		return false
	}
	if q.Zone != "" && e.Zone != q.Zone {
		return false
	}

	return true
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

// ToMessage converts the event to a core.Message
func (e *Event) ToMessage() (*core.Message, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	msg := core.NewMessage(data)
	msg.Headers.Set("event-type", string(e.Type))
	msg.Headers.Set("endpoint-id", e.EndpointID)
	msg.Headers.Set("service-id", e.ServiceID)
	return msg, nil
}

// EventFromMessage creates an Event from a core.Message
func EventFromMessage(msg *core.Message) (*Event, error) {
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
