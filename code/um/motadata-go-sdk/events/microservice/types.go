// Package microservice provides microservice registration, discovery, and lifecycle management.
//
// This package enables services to register themselves with a central registry,
// discover other services, and manage their lifecycle including health monitoring.
//
// Key concepts:
//   - Microservice: A logical service that can have multiple instances
//   - Instance: A running instance of a microservice
//   - Registry: Thread-safe storage for service and instance information
//   - Manager: Handles service lifecycle, health checking, and event propagation
package microservice

import (
	"encoding/json"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

// Status represents the current state of a microservice or instance
type Status int

const (
	// StatusUnknown indicates the status is not determined
	StatusUnknown Status = iota
	// StatusStarting indicates the service is starting up
	StatusStarting
	// StatusRunning indicates the service is running normally
	StatusRunning
	// StatusDegraded indicates the service is running with reduced capacity
	StatusDegraded
	// StatusStopping indicates the service is shutting down
	StatusStopping
	// StatusStopped indicates the service has stopped
	StatusStopped
	// StatusFailed indicates the service has failed
	StatusFailed
	// StatusMaintenance indicates the service is under maintenance
	StatusMaintenance
)

// String returns a human-readable status name
func (s Status) String() string {
	switch s {
	case StatusUnknown:
		return "unknown"
	case StatusStarting:
		return "starting"
	case StatusRunning:
		return "running"
	case StatusDegraded:
		return "degraded"
	case StatusStopping:
		return "stopping"
	case StatusStopped:
		return "stopped"
	case StatusFailed:
		return "failed"
	case StatusMaintenance:
		return "maintenance"
	default:
		return "unknown"
	}
}

// IsAvailable returns true if the service can accept requests
func (s Status) IsAvailable() bool {
	return s == StatusRunning || s == StatusDegraded
}

// IsHealthy returns true if the service is in a healthy state
func (s Status) IsHealthy() bool {
	return s == StatusRunning
}

// ServiceType represents the type of microservice
type ServiceType string

const (
	ServiceTypeAPI      ServiceType = "api"
	ServiceTypeWorker   ServiceType = "worker"
	ServiceTypeGateway  ServiceType = "gateway"
	ServiceTypeBroker   ServiceType = "broker"
	ServiceTypeDatabase ServiceType = "database"
	ServiceTypeCache    ServiceType = "cache"
	ServiceTypeQueue    ServiceType = "queue"
	ServiceTypeScheduler ServiceType = "scheduler"
	ServiceTypeMonitor   ServiceType = "monitor"
	ServiceTypeCustom   ServiceType = "custom"
)

// Info holds comprehensive microservice information
type Info struct {
	// Core identification
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	DisplayName string      `json:"display_name,omitempty"`
	Description string      `json:"description,omitempty"`
	Type        ServiceType `json:"type"`

	// Versioning
	Version     string `json:"version"`
	APIVersion  string `json:"api_version,omitempty"`
	BuildInfo   string `json:"build_info,omitempty"`
	GitCommit   string `json:"git_commit,omitempty"`

	// Service configuration
	Dependencies []string          `json:"dependencies,omitempty"`
	Capabilities []string          `json:"capabilities,omitempty"`
	Tags         []string          `json:"tags,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`

	// Health and status
	Status          Status    `json:"status"`
	LastHealthCheck time.Time `json:"last_health_check,omitempty"`
	HealthCheckURL  string    `json:"health_check_url,omitempty"`

	// Instance management
	Instances      []*Instance `json:"instances,omitempty"`
	MinInstances   int         `json:"min_instances,omitempty"`
	MaxInstances   int         `json:"max_instances,omitempty"`
	DesiredInstances int       `json:"desired_instances,omitempty"`

	// Lifecycle
	RegisteredAt time.Time `json:"registered_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Multi-tenancy
	TenantID string `json:"tenant_id,omitempty"`
	Owner    string `json:"owner,omitempty"`

	// Documentation
	DocsURL       string `json:"docs_url,omitempty"`
	RepositoryURL string `json:"repository_url,omitempty"`
}

// Clone creates a deep copy of the service info
func (s *Info) Clone() *Info {
	if s == nil {
		return nil
	}

	clone := *s

	// Deep copy slices
	if s.Dependencies != nil {
		clone.Dependencies = make([]string, len(s.Dependencies))
		copy(clone.Dependencies, s.Dependencies)
	}
	if s.Capabilities != nil {
		clone.Capabilities = make([]string, len(s.Capabilities))
		copy(clone.Capabilities, s.Capabilities)
	}
	if s.Tags != nil {
		clone.Tags = make([]string, len(s.Tags))
		copy(clone.Tags, s.Tags)
	}
	if s.Metadata != nil {
		clone.Metadata = make(map[string]string, len(s.Metadata))
		for k, v := range s.Metadata {
			clone.Metadata[k] = v
		}
	}
	if s.Instances != nil {
		clone.Instances = make([]*Instance, len(s.Instances))
		for i, inst := range s.Instances {
			clone.Instances[i] = inst.Clone()
		}
	}

	return &clone
}

// HasCapability checks if the service has a specific capability
func (s *Info) HasCapability(cap string) bool {
	for _, c := range s.Capabilities {
		if c == cap {
			return true
		}
	}
	return false
}

// HasTag checks if the service has a specific tag
func (s *Info) HasTag(tag string) bool {
	for _, t := range s.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// GetMetadata returns a metadata value or empty string if not found
func (s *Info) GetMetadata(key string) string {
	if s.Metadata == nil {
		return ""
	}
	return s.Metadata[key]
}

// SetMetadata sets a metadata value
func (s *Info) SetMetadata(key, value string) {
	if s.Metadata == nil {
		s.Metadata = make(map[string]string)
	}
	s.Metadata[key] = value
}

// HealthyInstanceCount returns the number of healthy instances
func (s *Info) HealthyInstanceCount() int {
	count := 0
	for _, inst := range s.Instances {
		if inst.Status == StatusRunning {
			count++
		}
	}
	return count
}

// AvailableInstanceCount returns the number of available instances
func (s *Info) AvailableInstanceCount() int {
	count := 0
	for _, inst := range s.Instances {
		if inst.Status.IsAvailable() {
			count++
		}
	}
	return count
}

// Instance represents a running instance of a microservice
type Instance struct {
	// Core identification
	ID        string `json:"id"`
	ServiceID string `json:"service_id"`

	// Network information
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Protocol string `json:"protocol,omitempty"`

	// Health and status
	Status           Status    `json:"status"`
	LastHealthCheck  time.Time `json:"last_health_check,omitempty"`
	HealthCheckPath  string    `json:"health_check_path,omitempty"`
	ConsecutiveFails int       `json:"consecutive_fails,omitempty"`

	// Load information
	Weight         int     `json:"weight,omitempty"`
	ActiveRequests int     `json:"active_requests,omitempty"`
	CPUUsage       float64 `json:"cpu_usage,omitempty"`
	MemoryUsage    float64 `json:"memory_usage,omitempty"`

	// Lifecycle
	StartedAt    time.Time `json:"started_at"`
	RegisteredAt time.Time `json:"registered_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`

	// Environment information
	Region     string            `json:"region,omitempty"`
	Zone       string            `json:"zone,omitempty"`
	Datacenter string            `json:"datacenter,omitempty"`
	Hostname   string            `json:"hostname,omitempty"`
	PodName    string            `json:"pod_name,omitempty"`
	ContainerID string           `json:"container_id,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// Clone creates a deep copy of the instance
func (i *Instance) Clone() *Instance {
	if i == nil {
		return nil
	}

	clone := *i
	if i.Metadata != nil {
		clone.Metadata = make(map[string]string, len(i.Metadata))
		for k, v := range i.Metadata {
			clone.Metadata[k] = v
		}
	}

	return &clone
}

// Address returns the host:port address
func (i *Instance) Address() string {
	if i.Port == 0 {
		return i.Host
	}
	return i.Host + ":" + itoa(i.Port)
}

// IsExpired returns true if the instance has expired
func (i *Instance) IsExpired() bool {
	if i.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(i.ExpiresAt)
}

// IsHealthy returns true if the instance is in a healthy state
func (i *Instance) IsHealthy() bool {
	return i.Status == StatusRunning
}

// GetMetadata returns a metadata value
func (i *Instance) GetMetadata(key string) string {
	if i.Metadata == nil {
		return ""
	}
	return i.Metadata[key]
}

// SetMetadata sets a metadata value
func (i *Instance) SetMetadata(key, value string) {
	if i.Metadata == nil {
		i.Metadata = make(map[string]string)
	}
	i.Metadata[key] = value
}

// RegistrationRequest represents a request to register a microservice
type RegistrationRequest struct {
	Service  *Info         `json:"service"`
	TTL      time.Duration `json:"ttl,omitempty"`
	Override bool          `json:"override,omitempty"`
}

// Validate validates the registration request
func (r *RegistrationRequest) Validate() error {
	if r.Service == nil {
		return ErrInvalidService
	}
	if r.Service.ID == "" {
		return ErrMissingServiceID
	}
	if r.Service.Name == "" {
		return ErrMissingServiceName
	}
	return nil
}

// InstanceRegistrationRequest represents a request to register a service instance
type InstanceRegistrationRequest struct {
	Instance *Instance     `json:"instance"`
	TTL      time.Duration `json:"ttl,omitempty"`
	Override bool          `json:"override,omitempty"`
}

// Validate validates the instance registration request
func (r *InstanceRegistrationRequest) Validate() error {
	if r.Instance == nil {
		return ErrInvalidInstance
	}
	if r.Instance.ID == "" {
		return ErrMissingInstanceID
	}
	if r.Instance.ServiceID == "" {
		return ErrMissingServiceID
	}
	if r.Instance.Host == "" {
		return ErrMissingHost
	}
	return nil
}

// DeregistrationRequest represents a request to deregister a service
type DeregistrationRequest struct {
	ServiceID  string `json:"service_id"`
	InstanceID string `json:"instance_id,omitempty"`
	Reason     string `json:"reason,omitempty"`
}

// Validate validates the deregistration request
func (r *DeregistrationRequest) Validate() error {
	if r.ServiceID == "" {
		return ErrMissingServiceID
	}
	return nil
}

// HealthUpdate represents a health status update
type HealthUpdate struct {
	ServiceID  string    `json:"service_id"`
	InstanceID string    `json:"instance_id,omitempty"`
	Status     Status    `json:"status"`
	Message    string    `json:"message,omitempty"`
	CheckedAt  time.Time `json:"checked_at"`
	Details    any       `json:"details,omitempty"`
}

// Query defines filters for searching services
type Query struct {
	// Filter by IDs
	IDs   []string `json:"ids,omitempty"`
	Names []string `json:"names,omitempty"`

	// Filter by type
	Types []ServiceType `json:"types,omitempty"`

	// Filter by status
	Statuses      []Status `json:"statuses,omitempty"`
	OnlyHealthy   bool     `json:"only_healthy,omitempty"`
	OnlyAvailable bool     `json:"only_available,omitempty"`

	// Filter by capabilities and tags
	Capabilities []string          `json:"capabilities,omitempty"`
	Tags         []string          `json:"tags,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`

	// Filter by version
	Version    string `json:"version,omitempty"`
	APIVersion string `json:"api_version,omitempty"`

	// Multi-tenancy filters
	TenantID string `json:"tenant_id,omitempty"`

	// Instance filters
	MinInstances    int  `json:"min_instances,omitempty"`
	HasInstances    bool `json:"has_instances,omitempty"`

	// Pagination
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

// Matches checks if a service matches the query filters
func (q *Query) Matches(s *Info) bool {
	if s == nil {
		return false
	}

	// ID/Name filters
	if len(q.IDs) > 0 && !containsString(q.IDs, s.ID) {
		return false
	}
	if len(q.Names) > 0 && !containsString(q.Names, s.Name) {
		return false
	}

	// Type filter
	if len(q.Types) > 0 && !containsServiceType(q.Types, s.Type) {
		return false
	}

	// Status filters
	if len(q.Statuses) > 0 && !containsStatus(q.Statuses, s.Status) {
		return false
	}
	if q.OnlyHealthy && s.Status != StatusRunning {
		return false
	}
	if q.OnlyAvailable && !s.Status.IsAvailable() {
		return false
	}

	// Capabilities filter (all must match)
	for _, cap := range q.Capabilities {
		if !s.HasCapability(cap) {
			return false
		}
	}

	// Tags filter (all must match)
	for _, tag := range q.Tags {
		if !s.HasTag(tag) {
			return false
		}
	}

	// Metadata filter (all must match)
	for k, v := range q.Metadata {
		if s.GetMetadata(k) != v {
			return false
		}
	}

	// Version filters
	if q.Version != "" && s.Version != q.Version {
		return false
	}
	if q.APIVersion != "" && s.APIVersion != q.APIVersion {
		return false
	}

	// Multi-tenancy filter
	if q.TenantID != "" && s.TenantID != q.TenantID {
		return false
	}

	// Instance filters
	if q.HasInstances && len(s.Instances) == 0 {
		return false
	}
	if q.MinInstances > 0 && len(s.Instances) < q.MinInstances {
		return false
	}

	return true
}

// EventType represents the type of lifecycle event
type EventType string

const (
	EventServiceRegistered   EventType = "service.registered"
	EventServiceDeregistered EventType = "service.deregistered"
	EventServiceUpdated      EventType = "service.updated"
	EventServiceHealthy      EventType = "service.healthy"
	EventServiceUnhealthy    EventType = "service.unhealthy"
	EventInstanceRegistered  EventType = "instance.registered"
	EventInstanceDeregistered EventType = "instance.deregistered"
	EventInstanceHealthy     EventType = "instance.healthy"
	EventInstanceUnhealthy   EventType = "instance.unhealthy"
)

// Event represents a service lifecycle event
type Event struct {
	Type       EventType `json:"type"`
	ServiceID  string    `json:"service_id"`
	InstanceID string    `json:"instance_id,omitempty"`
	Service    *Info     `json:"service,omitempty"`
	Instance   *Instance `json:"instance,omitempty"`
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
	msg.Headers.Set("service-id", e.ServiceID)
	if e.InstanceID != "" {
		msg.Headers.Set("instance-id", e.InstanceID)
	}
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

func containsServiceType(slice []ServiceType, t ServiceType) bool {
	for _, item := range slice {
		if item == t {
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
