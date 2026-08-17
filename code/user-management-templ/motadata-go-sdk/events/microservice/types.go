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
//
// # Status Lifecycle
//
// A service or instance transitions through the following statuses:
//
//	StatusUnknown -> StatusStarting -> StatusRunning -> StatusStopping -> StatusStopped
//	                                    |                                   ^
//	                                    +-> StatusDegraded ----------------+
//	                                    +-> StatusFailed ------------------+
//	                                    +-> StatusMaintenance -------------+
//
// Only StatusRunning and StatusDegraded are considered "available" for accepting
// requests (see Status.IsAvailable). Only StatusRunning is considered "healthy"
// (see Status.IsHealthy).
package microservice

import (
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"
)

// Status represents the current lifecycle state of a microservice or instance.
// It follows a defined transition model where services start as StatusUnknown,
// progress through StatusStarting and StatusRunning, and eventually reach
// StatusStopped. Intermediate states like StatusDegraded and StatusMaintenance
// indicate partial availability, while StatusFailed indicates a critical error.
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

// String returns the lowercase human-readable name of the status.
// Unknown or out-of-range values default to "unknown".
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

// IsAvailable returns true if the service can accept requests. Both
// StatusRunning and StatusDegraded are considered available because a
// degraded service still processes requests, albeit at reduced capacity.
func (s Status) IsAvailable() bool {
	return s == StatusRunning || s == StatusDegraded
}

// IsHealthy returns true only when the service is fully operational (StatusRunning).
// Unlike IsAvailable, degraded services are not considered healthy because
// they may have diminished capacity or reliability.
func (s Status) IsHealthy() bool {
	return s == StatusRunning
}

// ServiceType categorizes microservices by their architectural role. This
// classification enables type-based discovery queries, allowing callers to
// locate all API gateways, all workers, or all caches in the system without
// knowing their individual names.
type ServiceType string

const (
	// ServiceTypeAPI represents a REST or gRPC API service that handles synchronous client requests.
	ServiceTypeAPI ServiceType = "api"
	// ServiceTypeWorker represents a background processing service (e.g., job consumers, data pipelines).
	ServiceTypeWorker ServiceType = "worker"
	// ServiceTypeGateway represents an API gateway that routes and aggregates upstream service calls.
	ServiceTypeGateway ServiceType = "gateway"
	// ServiceTypeBroker represents a message broker service (e.g., Kafka, RabbitMQ frontends).
	ServiceTypeBroker ServiceType = "broker"
	// ServiceTypeDatabase represents a database or data-access service.
	ServiceTypeDatabase ServiceType = "database"
	// ServiceTypeCache represents a caching layer service (e.g., Redis, Memcached).
	ServiceTypeCache ServiceType = "cache"
	// ServiceTypeQueue represents a task/job queue service.
	ServiceTypeQueue ServiceType = "queue"
	// ServiceTypeScheduler represents a job scheduling or cron-like service.
	ServiceTypeScheduler ServiceType = "scheduler"
	// ServiceTypeMonitor represents a monitoring, metrics, or observability service.
	ServiceTypeMonitor ServiceType = "monitor"
	// ServiceTypeCustom represents a user-defined service type not covered by the predefined categories.
	ServiceTypeCustom ServiceType = "custom"
)

// Info holds comprehensive metadata about a logical microservice. A single Info
// record can own multiple Instance records, representing horizontally scaled
// replicas of the same service. The Info struct tracks the service's identity,
// versioning, declared capabilities, health status, multi-tenancy ownership,
// and the desired/actual instance counts for scaling decisions.
type Info struct {
	// Core identification
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	DisplayName string      `json:"display_name,omitempty"`
	Description string      `json:"description,omitempty"`
	Type        ServiceType `json:"type"`

	// Versioning
	Version    string `json:"version"`
	APIVersion string `json:"api_version,omitempty"`
	BuildInfo  string `json:"build_info,omitempty"`
	GitCommit  string `json:"git_commit,omitempty"`

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
	Instances        []*Instance `json:"instances,omitempty"`
	MinInstances     int         `json:"min_instances,omitempty"`
	MaxInstances     int         `json:"max_instances,omitempty"`
	DesiredInstances int         `json:"desired_instances,omitempty"`

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

// Clone creates a deep copy of the service info, including all slices, maps,
// and nested Instance pointers. This ensures that callers receive an isolated
// snapshot that can be read or modified without affecting registry state.
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

// HasCapability checks whether the service declares the given capability string.
// Capabilities are used for discovery: callers can query services that advertise
// specific features (e.g., "auth", "payments") without knowing service names.
func (s *Info) HasCapability(cap string) bool {
	for _, c := range s.Capabilities {
		if c == cap {
			return true
		}
	}
	return false
}

// HasTag checks whether the service has been labeled with the given tag.
// Tags provide a secondary classification mechanism (e.g., "core", "experimental")
// that can be combined with other query filters for fine-grained discovery.
func (s *Info) HasTag(tag string) bool {
	for _, t := range s.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// GetMetadata returns the metadata value associated with the given key, or an
// empty string if the key does not exist. This provides safe access to the
// arbitrary key-value metadata without risk of nil-map panics.
func (s *Info) GetMetadata(key string) string {
	if s.Metadata == nil {
		return ""
	}
	return s.Metadata[key]
}

// SetMetadata stores a key-value pair in the service's metadata map,
// lazily initializing the map if it has not been created yet.
func (s *Info) SetMetadata(key, value string) {
	if s.Metadata == nil {
		s.Metadata = make(map[string]string)
	}
	s.Metadata[key] = value
}

// HealthyInstanceCount returns the number of instances that are in StatusRunning.
// This count excludes degraded instances, providing a measure of fully operational
// capacity for scaling and alerting decisions.
func (s *Info) HealthyInstanceCount() int {
	count := 0
	for _, inst := range s.Instances {
		if inst.Status == StatusRunning {
			count++
		}
	}
	return count
}

// AvailableInstanceCount returns the number of instances that can accept
// requests (StatusRunning or StatusDegraded). This count is useful for load
// balancing decisions where degraded instances are still viable targets.
func (s *Info) AvailableInstanceCount() int {
	count := 0
	for _, inst := range s.Instances {
		if inst.Status.IsAvailable() {
			count++
		}
	}
	return count
}

// Instance represents a single running replica of a microservice. Multiple
// instances can belong to the same service (identified by ServiceID), enabling
// horizontal scaling. Each instance carries its own network endpoint (Host:Port),
// health status, load metrics (Weight, ActiveRequests, CPU/Memory usage), and
// optional TTL-based expiration for automatic cleanup of stale registrations.
//
// The Weight field supports weight-based load balancing: instances with higher
// weights receive proportionally more traffic. A weight of 0 means the instance
// uses the default weight. ActiveRequests, CPUUsage, and MemoryUsage provide
// real-time load information for more sophisticated balancing strategies.
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
	Region      string            `json:"region,omitempty"`
	Zone        string            `json:"zone,omitempty"`
	Datacenter  string            `json:"datacenter,omitempty"`
	Hostname    string            `json:"hostname,omitempty"`
	PodName     string            `json:"pod_name,omitempty"`
	ContainerID string            `json:"container_id,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// Clone creates a deep copy of the instance, including the metadata map.
// This prevents callers from accidentally mutating registry-owned data.
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

// Address returns the instance's network address as "host:port". If the port
// is zero (e.g., for Unix socket or discovery-only registrations), only the
// host is returned without a colon separator.
func (i *Instance) Address() string {
	if i.Port == 0 {
		return i.Host
	}
	return i.Host + ":" + itoa(i.Port)
}

// IsExpired returns true if the instance's TTL has elapsed. Instances without
// an ExpiresAt value (zero time) never expire and must be deregistered explicitly.
func (i *Instance) IsExpired() bool {
	if i.ExpiresAt.IsZero() {
		return false
	}
	return time.Now().After(i.ExpiresAt)
}

// IsHealthy returns true only when the instance is fully operational (StatusRunning).
// Degraded or failed instances are not considered healthy.
func (i *Instance) IsHealthy() bool {
	return i.Status == StatusRunning
}

// GetMetadata returns the metadata value for the given key, or an empty string
// if the key is absent or the metadata map is nil.
func (i *Instance) GetMetadata(key string) string {
	if i.Metadata == nil {
		return ""
	}
	return i.Metadata[key]
}

// SetMetadata stores a key-value pair in the instance's metadata map,
// lazily initializing the map if it has not been created yet.
func (i *Instance) SetMetadata(key, value string) {
	if i.Metadata == nil {
		i.Metadata = make(map[string]string)
	}
	i.Metadata[key] = value
}

// RegistrationRequest represents a request to register a microservice.
// The TTL field is reserved for future use at the service level (instance-level
// TTL is supported via InstanceRegistrationRequest). When Override is true,
// existing service metadata is replaced rather than merged.
type RegistrationRequest struct {
	Service  *Info         `json:"service"`
	TTL      time.Duration `json:"ttl,omitempty"`
	Override bool          `json:"override,omitempty"`
}

// Validate ensures the registration request contains a non-nil service with
// both an ID and a Name. These are the minimum required fields for a valid
// service registration entry.
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

// InstanceRegistrationRequest represents a request to register a single running
// instance of a service. When TTL is positive, the instance automatically
// expires after the specified duration unless renewed via heartbeat. Override
// controls whether an existing instance with the same ID is replaced entirely.
type InstanceRegistrationRequest struct {
	Instance *Instance     `json:"instance"`
	TTL      time.Duration `json:"ttl,omitempty"`
	Override bool          `json:"override,omitempty"`
}

// Validate ensures the instance registration request contains a non-nil instance
// with an ID, a ServiceID (to link it to its parent service), and a Host
// (the network address where the instance is reachable).
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

// DeregistrationRequest represents a request to remove a service (and optionally
// a specific instance) from the registry. When InstanceID is empty, the entire
// service and all its instances are removed. The Reason field is passed to
// lifecycle callbacks to provide context about why the deregistration occurred.
type DeregistrationRequest struct {
	ServiceID  string `json:"service_id"`
	InstanceID string `json:"instance_id,omitempty"`
	Reason     string `json:"reason,omitempty"`
}

// Validate ensures the deregistration request includes a ServiceID, which is
// the minimum identifier needed to locate the target in the registry.
func (r *DeregistrationRequest) Validate() error {
	if r.ServiceID == "" {
		return ErrMissingServiceID
	}
	return nil
}

// HealthUpdate represents a health status update for a service or instance.
// When InstanceID is set, the update targets a specific instance; otherwise
// it targets the service-level status. The Details field can carry structured
// diagnostic information (e.g., latency percentiles, error counts) from the
// health check probe.
type HealthUpdate struct {
	ServiceID  string    `json:"service_id"`
	InstanceID string    `json:"instance_id,omitempty"`
	Status     Status    `json:"status"`
	Message    string    `json:"message,omitempty"`
	CheckedAt  time.Time `json:"checked_at"`
	Details    any       `json:"details,omitempty"`
}

// Query defines a composable set of filters for searching services in the registry.
// All non-empty filter fields are combined with AND logic: a service must satisfy
// every specified criterion to appear in results. Slice filters (IDs, Names, Types,
// etc.) use OR logic within the slice (e.g., Types=["api","worker"] matches either).
// Capabilities and Tags require ALL listed values to be present on the service.
// Pagination is supported via Limit and Offset for large registries.
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
	MinInstances int  `json:"min_instances,omitempty"`
	HasInstances bool `json:"has_instances,omitempty"`

	// Pagination
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

// Matches evaluates all non-empty query filters against the given service Info.
// Returns true only if the service satisfies every criterion. Filters are applied
// in a short-circuit fashion: the first failing check immediately returns false,
// avoiding unnecessary comparisons for early mismatches.
func (q *Query) Matches(s *Info) bool {
	if s == nil {
		return false
	}

	return q.matchesIdentity(s) &&
		q.matchesStatus(s) &&
		q.matchesAttributes(s) &&
		q.matchesVersion(s) &&
		q.matchesInstances(s)
}

func (q *Query) matchesIdentity(s *Info) bool {
	if len(q.IDs) > 0 && !containsString(q.IDs, s.ID) {
		return false
	}
	if len(q.Names) > 0 && !containsString(q.Names, s.Name) {
		return false
	}
	return len(q.Types) == 0 || containsServiceType(q.Types, s.Type)
}

func (q *Query) matchesStatus(s *Info) bool {
	if len(q.Statuses) > 0 && !containsStatus(q.Statuses, s.Status) {
		return false
	}
	if q.OnlyHealthy && s.Status != StatusRunning {
		return false
	}
	return !q.OnlyAvailable || s.Status.IsAvailable()
}

func (q *Query) matchesAttributes(s *Info) bool {
	for _, cap := range q.Capabilities {
		if !s.HasCapability(cap) {
			return false
		}
	}
	for _, tag := range q.Tags {
		if !s.HasTag(tag) {
			return false
		}
	}
	for k, v := range q.Metadata {
		if s.GetMetadata(k) != v {
			return false
		}
	}
	return true
}

func (q *Query) matchesVersion(s *Info) bool {
	if q.Version != "" && s.Version != q.Version {
		return false
	}
	if q.APIVersion != "" && s.APIVersion != q.APIVersion {
		return false
	}
	return q.TenantID == "" || s.TenantID == q.TenantID
}

func (q *Query) matchesInstances(s *Info) bool {
	if q.HasInstances && len(s.Instances) == 0 {
		return false
	}
	return q.MinInstances <= 0 || len(s.Instances) >= q.MinInstances
}

// EventType represents the type of lifecycle event emitted by the Manager
// when service or instance state changes. These events are published over
// NATS to configured subjects, enabling reactive architectures where
// downstream systems can respond to registration, deregistration, and
// health transitions in real time.
type EventType string

const (
	// EventServiceRegistered is emitted when a new service is added to the registry.
	EventServiceRegistered EventType = "service.registered"
	// EventServiceDeregistered is emitted when a service is removed from the registry.
	EventServiceDeregistered EventType = "service.deregistered"
	// EventServiceUpdated is emitted when a service's metadata is modified.
	EventServiceUpdated EventType = "service.updated"
	// EventServiceHealthy is emitted when a service transitions to a healthy status.
	EventServiceHealthy EventType = "service.healthy"
	// EventServiceUnhealthy is emitted when a service transitions to an unhealthy status.
	EventServiceUnhealthy EventType = "service.unhealthy"
	// EventInstanceRegistered is emitted when a new instance is added to a service.
	EventInstanceRegistered EventType = "instance.registered"
	// EventInstanceDeregistered is emitted when an instance is removed from a service.
	EventInstanceDeregistered EventType = "instance.deregistered"
	// EventInstanceHealthy is emitted when an instance transitions to a healthy status.
	EventInstanceHealthy EventType = "instance.healthy"
	// EventInstanceUnhealthy is emitted when an instance transitions to an unhealthy status.
	EventInstanceUnhealthy EventType = "instance.unhealthy"
)

// Event represents a service or instance lifecycle event that is published
// over NATS messaging. It captures the event type, affected service/instance
// identifiers, old and new status values for transition tracking, and the
// full service or instance snapshot at the time of the event.
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

// ToMessage serializes the event to a NATS message with JSON body and
// structured headers (event-type, service-id, instance-id). Headers enable
// NATS consumers to filter or route messages without deserializing the body.
func (e *Event) ToMessage() (*nats.Msg, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	msg := &nats.Msg{Data: data, Header: make(nats.Header)}
	msg.Header.Set("event-type", string(e.Type))
	msg.Header.Set("service-id", e.ServiceID)
	if e.InstanceID != "" {
		msg.Header.Set("instance-id", e.InstanceID)
	}
	return msg, nil
}

// EventFromMessage deserializes a NATS message body into an Event struct.
// This is the inverse of Event.ToMessage and is used by event subscribers
// to reconstruct lifecycle events received from the messaging bus.
func EventFromMessage(msg *nats.Msg) (*Event, error) {
	var event Event
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}

// Helper functions for slice membership checks used by Query.Matches.

// containsString performs a linear search for s within the given string slice.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// containsStatus checks whether a Status value exists in the given slice.
func containsStatus(slice []Status, s Status) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// containsServiceType checks whether a ServiceType value exists in the given slice.
func containsServiceType(slice []ServiceType, t ServiceType) bool {
	for _, item := range slice {
		if item == t {
			return true
		}
	}
	return false
}

// itoa converts an integer to its string representation without importing strconv.
// It uses a fixed-size byte buffer and writes digits in reverse order for efficiency.
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
