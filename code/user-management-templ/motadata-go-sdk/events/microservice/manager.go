package microservice

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/tracer"
)

var (
	attrServiceID       = "service.id"
	attrInstanceID      = "instance.id"
	errMethodNotAllowed = "Method not allowed"
)

// Manager orchestrates the full microservice lifecycle: registration, deregistration,
// health monitoring, discovery, and event propagation. It sits above the Registry
// layer, adding observability (tracing, logging), asynchronous event publishing
// over NATS, and periodic HTTP health checks for registered instances.
//
// Service vs Instance lifecycle:
//   - A "service" is a logical entity registered once. It holds metadata, capabilities,
//     and version info. It exists independently of running instances.
//   - An "instance" is a physical/running replica of a service. Multiple instances
//     share the same ServiceID. Instances can be registered with TTLs, requiring
//     periodic heartbeats to stay alive. When an instance is deregistered (manually
//     or via TTL expiration), the parent service remains.
//   - Deregistering a service cascades: all its instances are removed as well.
//
// The Manager publishes lifecycle events (registration, deregistration, health
// changes) to configured NATS subjects so that external systems can react to
// state changes in real time.
type Manager struct {
	registry *Registry

	// Publishing
	publisher       core.Publisher
	publishSubjects PublishSubjects

	// Health checking
	healthChecker   *HealthChecker
	healthConfig    HealthCheckConfig
	healthEnabled   bool
	healthInstances map[string]struct{}
	healthMu        sync.RWMutex

	// Lifecycle
	done   <-chan struct{}
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Error handling
	errorChan chan error
	errorCb   func(error)
	errorMu   sync.RWMutex
}

// PublishSubjects defines the NATS subject strings to which lifecycle events are
// published. Each event type maps to a distinct subject, allowing subscribers to
// selectively listen for specific lifecycle transitions (e.g., only health changes).
type PublishSubjects struct {
	ServiceRegistered    string
	ServiceDeregistered  string
	ServiceUpdated       string
	ServiceHealth        string
	InstanceRegistered   string
	InstanceDeregistered string
	InstanceHealth       string
}

// DefaultPublishSubjects returns the standard subject naming convention using
// dot-separated hierarchies (e.g., "services.registered", "instances.health").
func DefaultPublishSubjects() PublishSubjects {
	return PublishSubjects{
		ServiceRegistered:    "services.registered",
		ServiceDeregistered:  "services.deregistered",
		ServiceUpdated:       "services.updated",
		ServiceHealth:        "services.health",
		InstanceRegistered:   "instances.registered",
		InstanceDeregistered: "instances.deregistered",
		InstanceHealth:       "instances.health",
	}
}

// HealthCheckConfig configures the periodic health checking behavior of the Manager.
// The health check loop runs every Interval, sending HTTP GET requests to each
// instance's HealthCheckPath. An instance transitions to healthy only after
// HealthyThreshold consecutive successes, and to unhealthy after UnhealthyThreshold
// consecutive failures, preventing status flapping from transient errors.
type HealthCheckConfig struct {
	Enabled            bool
	Interval           time.Duration
	Timeout            time.Duration
	HealthyThreshold   int
	UnhealthyThreshold int
	InitialDelay       time.Duration
}

// DefaultHealthCheckConfig returns sensible defaults: checks every 30s with a 5s
// timeout, 2 consecutive successes to mark healthy, 3 consecutive failures to
// mark unhealthy, and a 10s initial delay to allow services to finish starting.
func DefaultHealthCheckConfig() HealthCheckConfig {
	return HealthCheckConfig{
		Enabled:            true,
		Interval:           30 * time.Second,
		Timeout:            5 * time.Second,
		HealthyThreshold:   2,
		UnhealthyThreshold: 3,
		InitialDelay:       10 * time.Second,
	}
}

// ManagerConfig holds the configuration required to construct a Manager. When
// Registry is nil, a default registry is created. When Publisher is nil, the
// Manager operates without event publishing (useful for testing).
type ManagerConfig struct {
	Registry        *Registry
	Publisher       core.Publisher
	PublishSubjects PublishSubjects
	HealthCheck     HealthCheckConfig
	ErrorBufferSize int
}

// NewManager creates a new microservice manager with the given configuration.
// It initializes the registry (creating a default if none is provided), sets up
// lifecycle callbacks that bridge registry events to NATS publishing, starts the
// background health check loop (if enabled), and returns a ready-to-use Manager.
func NewManager(cfg ManagerConfig) (*Manager, error) {
	if cfg.Registry == nil {
		cfg.Registry = NewRegistry()
	}

	if cfg.PublishSubjects.ServiceRegistered == "" {
		cfg.PublishSubjects = DefaultPublishSubjects()
	}

	if cfg.ErrorBufferSize == 0 {
		cfg.ErrorBufferSize = 100
	}

	ctx, cancel := context.WithCancel(context.Background())

	m := &Manager{
		registry:        cfg.Registry,
		publisher:       cfg.Publisher,
		publishSubjects: cfg.PublishSubjects,
		healthConfig:    cfg.HealthCheck,
		healthEnabled:   cfg.HealthCheck.Enabled,
		healthInstances: make(map[string]struct{}),
		done:            ctx.Done(),
		cancel:          cancel,
		errorChan:       make(chan error, cfg.ErrorBufferSize),
	}

	m.healthChecker = NewHealthChecker(HealthCheckerConfig{
		Timeout:            cfg.HealthCheck.Timeout,
		HealthyThreshold:   cfg.HealthCheck.HealthyThreshold,
		UnhealthyThreshold: cfg.HealthCheck.UnhealthyThreshold,
	})

	m.setupCallbacks()

	if cfg.HealthCheck.Enabled {
		m.startHealthCheckLoop()
	}

	logger.Info(ctx, "microservice manager initialized",
		logger.Bool("health_check_enabled", cfg.HealthCheck.Enabled),
		logger.Duration("health_check_interval", cfg.HealthCheck.Interval),
	)

	return m, nil
}

// setupCallbacks wires registry lifecycle callbacks to the Manager's event publishing
// pipeline. Each registry event (service registered/deregistered/updated, instance
// registered/deregistered, status changes) triggers an asynchronous NATS publish.
// Instance registration callbacks also manage the health check tracking set.
func (m *Manager) setupCallbacks() {
	m.registry.OnServiceRegistered(func(svc *Info) {
		m.publishServiceEvent(EventServiceRegistered, svc, StatusUnknown, svc.Status, "")
	})

	m.registry.OnServiceDeregistered(func(svc *Info, reason string) {
		m.publishServiceEvent(EventServiceDeregistered, svc, svc.Status, StatusStopped, reason)
	})

	m.registry.OnServiceUpdated(func(old, new *Info) {
		m.publishServiceEvent(EventServiceUpdated, new, old.Status, new.Status, "")
	})

	m.registry.OnServiceStatusChange(func(svc *Info, oldStatus, newStatus Status) {
		eventType := EventServiceHealthy
		if !newStatus.IsHealthy() {
			eventType = EventServiceUnhealthy
		}
		m.publishServiceEvent(eventType, svc, oldStatus, newStatus, "")
	})

	m.registry.OnInstanceRegistered(func(inst *Instance) {
		m.publishInstanceEvent(EventInstanceRegistered, inst, StatusUnknown, inst.Status, "")
		if m.healthEnabled && inst.HealthCheckPath != "" {
			m.addToHealthCheck(inst.ID)
		}
	})

	m.registry.OnInstanceDeregistered(func(inst *Instance, reason string) {
		m.publishInstanceEvent(EventInstanceDeregistered, inst, inst.Status, StatusStopped, reason)
		if m.healthEnabled {
			m.removeFromHealthCheck(inst.ID)
		}
	})

	m.registry.OnInstanceStatusChange(func(inst *Instance, oldStatus, newStatus Status) {
		eventType := EventInstanceHealthy
		if !newStatus.IsHealthy() {
			eventType = EventInstanceUnhealthy
		}
		m.publishInstanceEvent(eventType, inst, oldStatus, newStatus, "")
	})
}

// startHealthCheckLoop launches a background goroutine that periodically checks
// the health of all registered instances. The loop first waits for InitialDelay
// (to allow services to finish starting), then ticks at Interval, calling
// performHealthChecks on each tick. The goroutine exits when the Manager's
// context is cancelled during shutdown.
func (m *Manager) startHealthCheckLoop() {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()

		if m.healthConfig.InitialDelay > 0 {
			select {
			case <-m.done:
				return
			case <-time.After(m.healthConfig.InitialDelay):
			}
		}

		ticker := time.NewTicker(m.healthConfig.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-m.done:
				return
			case <-ticker.C:
				m.performHealthChecks()
			}
		}
	}()
}

// performHealthChecks iterates over all instances in the health check tracking set,
// performs an HTTP health check for each one that has a HealthCheckPath configured,
// and updates the registry with the result. Instances that no longer exist in the
// registry are silently skipped. The health check results flow through the
// HealthChecker's threshold logic to prevent status flapping.
func (m *Manager) performHealthChecks() {
	ctx, span := tracer.Start(context.Background(), "microservice.health_check_cycle")
	defer span.End()

	m.healthMu.RLock()
	instanceIDs := make([]string, 0, len(m.healthInstances))
	for id := range m.healthInstances {
		instanceIDs = append(instanceIDs, id)
	}
	m.healthMu.RUnlock()

	span.SetAttributes(tracer.IntAttr("instances.count", len(instanceIDs)))

	for _, id := range instanceIDs {
		inst, exists := m.registry.GetInstance(id)
		if !exists {
			continue
		}

		if inst.HealthCheckPath == "" {
			continue
		}

		result := m.healthChecker.Check(ctx, inst)

		if err := m.registry.UpdateInstanceHealth(&HealthUpdate{
			ServiceID:  inst.ServiceID,
			InstanceID: id,
			Status:     result.Status,
			Message:    result.Message,
			CheckedAt:  result.Timestamp,
			Details:    result.Details,
		}); err != nil {
			logger.Error(ctx, "failed to update instance health",
				logger.String("instance_id", id),
				logger.String("service_id", inst.ServiceID),
				logger.Err(err),
			)
			m.dispatchError(err)
		}
	}

	span.SetOK()
	logger.Debug(ctx, "health check cycle completed",
		logger.Int("instances_checked", len(instanceIDs)),
	)
}

// addToHealthCheck adds an instance to the set of instances that the health
// check loop will probe on each cycle.
func (m *Manager) addToHealthCheck(instanceID string) {
	m.healthMu.Lock()
	m.healthInstances[instanceID] = struct{}{}
	m.healthMu.Unlock()
}

// removeFromHealthCheck removes an instance from the health check tracking set,
// so it will no longer be probed on subsequent health check cycles.
func (m *Manager) removeFromHealthCheck(instanceID string) {
	m.healthMu.Lock()
	delete(m.healthInstances, instanceID)
	m.healthMu.Unlock()
}

// publishServiceEvent constructs a lifecycle Event for a service and publishes it
// asynchronously to the appropriate NATS subject. If no publisher is configured,
// the method returns immediately. Publishing is done in a separate goroutine to
// avoid blocking the caller.
func (m *Manager) publishServiceEvent(eventType EventType, svc *Info, oldStatus, newStatus Status, reason string) {
	if m.publisher == nil {
		return
	}

	event := &Event{
		Type:      eventType,
		ServiceID: svc.ID,
		Service:   svc,
		OldStatus: oldStatus,
		NewStatus: newStatus,
		Reason:    reason,
		Timestamp: time.Now(),
	}

	msg, err := event.ToMessage()
	if err != nil {
		m.dispatchError(err)
		return
	}

	var subject string
	switch eventType {
	case EventServiceRegistered:
		subject = m.publishSubjects.ServiceRegistered
	case EventServiceDeregistered:
		subject = m.publishSubjects.ServiceDeregistered
	case EventServiceUpdated:
		subject = m.publishSubjects.ServiceUpdated
	case EventServiceHealthy, EventServiceUnhealthy:
		subject = m.publishSubjects.ServiceHealth
	default:
		subject = m.publishSubjects.ServiceUpdated
	}

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := m.publisher.Publish(ctx, subject, msg); err != nil {
			m.dispatchError(err)
		}
	}()
}

// publishInstanceEvent constructs a lifecycle Event for an instance and publishes
// it asynchronously to the appropriate NATS subject. Like publishServiceEvent,
// it no-ops when no publisher is configured.
func (m *Manager) publishInstanceEvent(eventType EventType, inst *Instance, oldStatus, newStatus Status, reason string) {
	if m.publisher == nil {
		return
	}

	event := &Event{
		Type:       eventType,
		ServiceID:  inst.ServiceID,
		InstanceID: inst.ID,
		Instance:   inst,
		OldStatus:  oldStatus,
		NewStatus:  newStatus,
		Reason:     reason,
		Timestamp:  time.Now(),
	}

	msg, err := event.ToMessage()
	if err != nil {
		m.dispatchError(err)
		return
	}

	var subject string
	switch eventType {
	case EventInstanceRegistered:
		subject = m.publishSubjects.InstanceRegistered
	case EventInstanceDeregistered:
		subject = m.publishSubjects.InstanceDeregistered
	case EventInstanceHealthy, EventInstanceUnhealthy:
		subject = m.publishSubjects.InstanceHealth
	default:
		subject = m.publishSubjects.InstanceRegistered
	}

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := m.publisher.Publish(ctx, subject, msg); err != nil {
			m.dispatchError(err)
		}
	}()
}

// dispatchError routes an error to both the error callback (if set) and the
// buffered error channel. If the channel is full, the error is silently dropped
// to prevent blocking the caller.
func (m *Manager) dispatchError(err error) {
	m.errorMu.RLock()
	cb := m.errorCb
	m.errorMu.RUnlock()

	if cb != nil {
		cb(err)
	}

	select {
	case m.errorChan <- err:
	default:
	}
}

// RegisterService validates and registers a new service through the RegistrationRequest
// envelope. It creates an OpenTelemetry span for tracing, validates the request,
// delegates to the registry, and emits structured log messages on success or failure.
func (m *Manager) RegisterService(ctx context.Context, req *RegistrationRequest) error {
	ctx, span := tracer.Start(ctx, "microservice.register_service")
	defer span.End()

	// Validate first to avoid nil pointer access
	if err := req.Validate(); err != nil {
		span.SetError(err)
		logger.Error(ctx, "service registration validation failed",
			logger.Err(err),
		)
		return err
	}

	// Now safe to access service fields
	span.SetAttributes(
		tracer.StringAttr(attrServiceID, req.Service.ID),
		tracer.StringAttr("service.name", req.Service.Name),
	)

	if err := m.registry.RegisterService(req.Service); err != nil {
		span.SetError(err)
		logger.Error(ctx, "failed to register service",
			logger.String("service_id", req.Service.ID),
			logger.Err(err),
		)
		return err
	}

	span.SetOK()
	logger.Info(ctx, "service registered",
		logger.String("service_id", req.Service.ID),
		logger.String("service_name", req.Service.Name),
	)
	return nil
}

// RegisterServiceDirect is a convenience method that registers a service directly
// from an Info struct, bypassing the RegistrationRequest envelope. This is useful
// when callers already have a validated Info and do not need TTL or override options.
func (m *Manager) RegisterServiceDirect(ctx context.Context, service *Info) error {
	ctx, span := tracer.Start(ctx, "microservice.register_service_direct",
		tracer.WithAttributes(
			tracer.StringAttr(attrServiceID, service.ID),
			tracer.StringAttr("service.name", service.Name),
		),
	)
	defer span.End()

	if err := m.registry.RegisterService(service); err != nil {
		span.SetError(err)
		logger.Error(ctx, "failed to register service",
			logger.String("service_id", service.ID),
			logger.Err(err),
		)
		return err
	}

	span.SetOK()
	logger.Info(ctx, "service registered",
		logger.String("service_id", service.ID),
		logger.String("service_name", service.Name),
	)
	return nil
}

// DeregisterService removes a service and all its instances from the registry.
// The deregistration reason from the request is forwarded to lifecycle callbacks
// so that subscribers can distinguish between planned shutdowns and forced removals.
func (m *Manager) DeregisterService(ctx context.Context, req *DeregistrationRequest) error {
	ctx, span := tracer.Start(ctx, "microservice.deregister_service")
	defer span.End()

	// Validate first
	if err := req.Validate(); err != nil {
		span.SetError(err)
		logger.Error(ctx, "service deregistration validation failed",
			logger.Err(err),
		)
		return err
	}

	// Now safe to access fields
	span.SetAttributes(tracer.StringAttr(attrServiceID, req.ServiceID))

	if err := m.registry.DeregisterServiceWithReason(req.ServiceID, req.Reason); err != nil {
		span.SetError(err)
		logger.Error(ctx, "failed to deregister service",
			logger.String("service_id", req.ServiceID),
			logger.Err(err),
		)
		return err
	}

	span.SetOK()
	logger.Info(ctx, "service deregistered",
		logger.String("service_id", req.ServiceID),
		logger.String("reason", req.Reason),
	)
	return nil
}

// DeregisterServiceDirect is a convenience method that deregisters a service by ID
// without requiring a DeregistrationRequest envelope. No reason is recorded.
func (m *Manager) DeregisterServiceDirect(ctx context.Context, serviceID string) error {
	ctx, span := tracer.Start(ctx, "microservice.deregister_service_direct",
		tracer.WithAttributes(
			tracer.StringAttr(attrServiceID, serviceID),
		),
	)
	defer span.End()

	if err := m.registry.DeregisterService(serviceID); err != nil {
		span.SetError(err)
		logger.Error(ctx, "failed to deregister service",
			logger.String("service_id", serviceID),
			logger.Err(err),
		)
		return err
	}

	span.SetOK()
	logger.Info(ctx, "service deregistered",
		logger.String("service_id", serviceID),
	)
	return nil
}

// RegisterInstance validates and registers a new instance through the
// InstanceRegistrationRequest envelope. When the request includes a positive TTL,
// the instance is registered with automatic expiration. The parent service must
// already exist in the registry.
func (m *Manager) RegisterInstance(ctx context.Context, req *InstanceRegistrationRequest) error {
	ctx, span := tracer.Start(ctx, "microservice.register_instance")
	defer span.End()

	// Validate first to avoid nil pointer access
	if err := req.Validate(); err != nil {
		span.SetError(err)
		logger.Error(ctx, "instance registration validation failed",
			logger.Err(err),
		)
		return err
	}

	// Now safe to access instance fields
	span.SetAttributes(
		tracer.StringAttr(attrInstanceID, req.Instance.ID),
		tracer.StringAttr(attrServiceID, req.Instance.ServiceID),
	)

	var err error
	if req.TTL > 0 {
		err = m.registry.RegisterInstanceWithTTL(req.Instance, req.TTL)
		span.SetAttributes(tracer.StringAttr("ttl", req.TTL.String()))
	} else {
		err = m.registry.RegisterInstance(req.Instance)
	}

	if err != nil {
		span.SetError(err)
		logger.Error(ctx, "failed to register instance",
			logger.String("instance_id", req.Instance.ID),
			logger.String("service_id", req.Instance.ServiceID),
			logger.Err(err),
		)
		return err
	}

	span.SetOK()
	logger.Info(ctx, "instance registered",
		logger.String("instance_id", req.Instance.ID),
		logger.String("service_id", req.Instance.ServiceID),
	)
	return nil
}

// RegisterInstanceDirect is a convenience method that registers an instance directly
// from an Instance struct, bypassing the InstanceRegistrationRequest envelope.
// The instance is registered without a TTL (it never auto-expires).
func (m *Manager) RegisterInstanceDirect(ctx context.Context, instance *Instance) error {
	ctx, span := tracer.Start(ctx, "microservice.register_instance_direct",
		tracer.WithAttributes(
			tracer.StringAttr(attrInstanceID, instance.ID),
			tracer.StringAttr(attrServiceID, instance.ServiceID),
		),
	)
	defer span.End()

	if err := m.registry.RegisterInstance(instance); err != nil {
		span.SetError(err)
		logger.Error(ctx, "failed to register instance",
			logger.String("instance_id", instance.ID),
			logger.String("service_id", instance.ServiceID),
			logger.Err(err),
		)
		return err
	}

	span.SetOK()
	logger.Info(ctx, "instance registered",
		logger.String("instance_id", instance.ID),
		logger.String("service_id", instance.ServiceID),
	)
	return nil
}

// RegisterInstanceWithTTL registers an instance that will automatically expire
// after the given TTL duration unless renewed via RenewInstance. This is the
// primary mechanism for ephemeral instance registration where instances send
// periodic heartbeats to extend their TTL.
func (m *Manager) RegisterInstanceWithTTL(ctx context.Context, instance *Instance, ttl time.Duration) error {
	ctx, span := tracer.Start(ctx, "microservice.register_instance_ttl",
		tracer.WithAttributes(
			tracer.StringAttr(attrInstanceID, instance.ID),
			tracer.StringAttr(attrServiceID, instance.ServiceID),
			tracer.StringAttr("ttl", ttl.String()),
		),
	)
	defer span.End()

	if err := m.registry.RegisterInstanceWithTTL(instance, ttl); err != nil {
		span.SetError(err)
		logger.Error(ctx, "failed to register instance with TTL",
			logger.String("instance_id", instance.ID),
			logger.String("service_id", instance.ServiceID),
			logger.Err(err),
		)
		return err
	}

	span.SetOK()
	logger.Info(ctx, "instance registered with TTL",
		logger.String("instance_id", instance.ID),
		logger.String("service_id", instance.ServiceID),
		logger.Duration("ttl", ttl),
	)
	return nil
}

// DeregisterInstance removes a specific instance from the registry by its ID.
// The parent service is not affected; only the single instance is removed.
func (m *Manager) DeregisterInstance(ctx context.Context, instanceID string) error {
	ctx, span := tracer.Start(ctx, "microservice.deregister_instance",
		tracer.WithAttributes(
			tracer.StringAttr(attrInstanceID, instanceID),
		),
	)
	defer span.End()

	if err := m.registry.DeregisterInstance(instanceID); err != nil {
		span.SetError(err)
		logger.Error(ctx, "failed to deregister instance",
			logger.String("instance_id", instanceID),
			logger.Err(err),
		)
		return err
	}

	span.SetOK()
	logger.Info(ctx, "instance deregistered",
		logger.String("instance_id", instanceID),
	)
	return nil
}

// GetService returns a deep copy of the service identified by serviceID, or
// (nil, false) if no such service exists. The returned copy is safe to modify.
func (m *Manager) GetService(serviceID string) (*Info, bool) {
	return m.registry.GetService(serviceID)
}

// GetServiceByName returns all services registered under the given name. Multiple
// services can share a name (e.g., different versions), so a slice is returned.
func (m *Manager) GetServiceByName(name string) []*Info {
	return m.registry.GetServiceByName(name)
}

// GetInstance returns a deep copy of the instance identified by instanceID, or
// (nil, false) if it does not exist or has expired.
func (m *Manager) GetInstance(instanceID string) (*Instance, bool) {
	return m.registry.GetInstance(instanceID)
}

// GetInstances returns deep copies of all non-expired instances belonging to the
// given service, regardless of their health status.
func (m *Manager) GetInstances(serviceID string) []*Instance {
	return m.registry.GetInstances(serviceID)
}

// GetHealthyInstances returns only the instances in StatusRunning for the given
// service. This is suitable for strict health-based load balancing.
func (m *Manager) GetHealthyInstances(serviceID string) []*Instance {
	return m.registry.GetHealthyInstances(serviceID)
}

// GetAvailableInstances returns instances that can accept requests (StatusRunning
// or StatusDegraded) for the given service. This provides a broader pool than
// GetHealthyInstances by including degraded but still functional replicas.
func (m *Manager) GetAvailableInstances(serviceID string) []*Instance {
	return m.registry.GetAvailableInstances(serviceID)
}

// QueryServices executes a structured query against the registry, returning all
// services matching the filter criteria. Results are sorted by name and support
// pagination via Query.Offset and Query.Limit.
func (m *Manager) QueryServices(q *Query) []*Info {
	return m.registry.QueryServices(q)
}

// Discover provides a high-level service discovery API. It searches for services
// by name with sensible defaults (only available services with at least one instance).
// The algorithm:
//  1. Builds a base Query with the given service name, OnlyAvailable=true, HasInstances=true.
//  2. Applies each DiscoverOption to override or augment the query (type, capabilities, tags, version, etc.).
//  3. Executes the query against the registry and returns matching services.
//
// This is the recommended entry point for service-to-service discovery.
func (m *Manager) Discover(name string, opts ...DiscoverOption) []*Info {
	q := &Query{
		Names:         []string{name},
		OnlyAvailable: true,
		HasInstances:  true,
	}

	for _, opt := range opts {
		opt(q)
	}

	return m.registry.QueryServices(q)
}

// DiscoverOption is a functional option that modifies the discovery Query built
// by Manager.Discover. Options are applied in order and can override each other.
type DiscoverOption func(*Query)

// WithType narrows discovery results to services of the specified ServiceType
// (e.g., ServiceTypeAPI, ServiceTypeWorker).
func WithType(t ServiceType) DiscoverOption {
	return func(q *Query) {
		q.Types = []ServiceType{t}
	}
}

// WithCapabilities narrows discovery results to services that declare ALL of the
// specified capability strings (AND logic).
func WithCapabilities(caps ...string) DiscoverOption {
	return func(q *Query) {
		q.Capabilities = caps
	}
}

// WithTags narrows discovery results to services that have ALL of the specified
// tags (AND logic).
func WithTags(tags ...string) DiscoverOption {
	return func(q *Query) {
		q.Tags = tags
	}
}

// WithVersion narrows discovery results to services matching the exact version string.
func WithVersion(version string) DiscoverOption {
	return func(q *Query) {
		q.Version = version
	}
}

// OnlyHealthy restricts discovery to services in StatusRunning only (excludes
// StatusDegraded). It overrides the default OnlyAvailable filter.
func OnlyHealthy() DiscoverOption {
	return func(q *Query) {
		q.OnlyHealthy = true
		q.OnlyAvailable = false
	}
}

// WithMinInstances requires that discovered services have at least the specified
// number of registered instances. Useful for ensuring capacity before routing.
func WithMinInstances(min int) DiscoverOption {
	return func(q *Query) {
		q.MinInstances = min
	}
}

// WithLimit caps the maximum number of services returned by discovery.
func WithLimit(limit int) DiscoverOption {
	return func(q *Query) {
		q.Limit = limit
	}
}

// DiscoverInstances returns instances for a service, filtered by health status.
// When onlyHealthy is true, only StatusRunning instances are returned; otherwise,
// both Running and Degraded (available) instances are included.
func (m *Manager) DiscoverInstances(serviceID string, onlyHealthy bool) []*Instance {
	if onlyHealthy {
		return m.registry.GetHealthyInstances(serviceID)
	}
	return m.registry.GetAvailableInstances(serviceID)
}

// RenewInstance extends the TTL of a registered instance, effectively acting as
// a heartbeat. Instances must call this periodically (before their TTL expires)
// to prevent automatic expiration and removal by the registry cleanup loop.
func (m *Manager) RenewInstance(ctx context.Context, instanceID string, ttl time.Duration) error {
	return m.registry.RenewInstance(instanceID, ttl)
}

// UpdateServiceStatus explicitly sets the status of a service. If the status differs
// from the current value, a status change callback and lifecycle event are triggered.
func (m *Manager) UpdateServiceStatus(ctx context.Context, serviceID string, status Status) error {
	return m.registry.UpdateServiceStatus(serviceID, status)
}

// UpdateInstanceStatus explicitly sets the status of an instance. Status transitions
// trigger callbacks and lifecycle events, and affect the instance's ConsecutiveFails counter.
func (m *Manager) UpdateInstanceStatus(ctx context.Context, instanceID string, status Status) error {
	return m.registry.UpdateInstanceStatus(instanceID, status)
}

// AllServices returns deep copies of all registered services in the registry.
func (m *Manager) AllServices() []*Info {
	return m.registry.AllServices()
}

// AllInstances returns deep copies of all non-expired instances across all services.
func (m *Manager) AllInstances() []*Instance {
	return m.registry.AllInstances()
}

// ServiceCount returns the total number of services currently registered.
func (m *Manager) ServiceCount() int {
	return m.registry.ServiceCount()
}

// InstanceCount returns the total number of instances currently registered
// (including expired ones that have not yet been cleaned up).
func (m *Manager) InstanceCount() int {
	return m.registry.InstanceCount()
}

// Registry returns the underlying Registry for direct access to low-level
// registration and query operations that are not exposed through the Manager API.
func (m *Manager) Registry() *Registry {
	return m.registry
}

// SetPublisher sets or replaces the NATS publisher used for lifecycle event
// broadcasting. Passing nil disables event publishing.
func (m *Manager) SetPublisher(publisher core.Publisher) {
	m.publisher = publisher
}

// Errors returns a read-only channel that receives errors from asynchronous
// operations (event publishing failures, health check errors). The channel is
// buffered (default size 100) and drops errors when full to avoid blocking.
func (m *Manager) Errors() <-chan error {
	return m.errorChan
}

// OnError registers a synchronous error callback that is invoked whenever an
// asynchronous error occurs. This is an alternative to consuming the Errors()
// channel and is useful for simple logging or alerting integrations.
func (m *Manager) OnError(cb func(error)) {
	m.errorMu.Lock()
	m.errorCb = cb
	m.errorMu.Unlock()
}

// Shutdown gracefully shuts down the manager by cancelling the background context
// (stopping the health check loop and in-flight publishes), waiting for all
// goroutines to complete (with ctx-based timeout), shutting down the registry,
// and closing the error channel. The provided context controls the maximum time
// allowed for shutdown; if it expires, the method returns the context error.
func (m *Manager) Shutdown(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "microservice.manager.shutdown")
	defer span.End()

	logger.Info(ctx, "initiating microservice manager shutdown",
		logger.Int("service_count", m.ServiceCount()),
		logger.Int("instance_count", m.InstanceCount()),
	)

	m.cancel()

	done := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		span.SetError(ctx.Err())
		logger.Warn(ctx, "microservice manager shutdown timeout")
		return ctx.Err()
	case <-done:
	}

	if err := m.registry.Shutdown(ctx); err != nil {
		span.SetError(err)
		logger.Error(ctx, "failed to shutdown registry",
			logger.Err(err),
		)
		return err
	}

	close(m.errorChan)

	span.SetOK()
	logger.Info(ctx, "microservice manager shutdown complete")
	return nil
}

// HealthChecker performs HTTP-based health checks on service instances and applies
// threshold-based status transitions to prevent flapping. It maintains per-instance
// counters of consecutive successes and failures. An instance is only marked as
// healthy (StatusRunning) after reaching healthyThreshold consecutive successful
// checks, and only marked as unhealthy after unhealthyThreshold consecutive failures.
// During the transition period, the instance is reported as StatusDegraded ("Recovering").
type HealthChecker struct {
	client             *http.Client
	timeout            time.Duration
	healthyThreshold   int
	unhealthyThreshold int

	mu         sync.RWMutex
	successCnt map[string]int
	failureCnt map[string]int
}

// HealthCheckerConfig configures the health checker's HTTP timeout and the
// threshold counters that determine when status transitions actually take effect.
type HealthCheckerConfig struct {
	Timeout            time.Duration
	HealthyThreshold   int
	UnhealthyThreshold int
}

// HealthCheckResult represents the outcome of a single health check probe
// against an instance, including the determined status, timing information,
// and optional diagnostic details.
type HealthCheckResult struct {
	InstanceID string
	ServiceID  string
	Status     Status
	Message    string
	Timestamp  time.Time
	Duration   time.Duration
	Details    map[string]any
}

// NewHealthChecker creates a new health checker with the given configuration.
// Zero values for Timeout, HealthyThreshold, and UnhealthyThreshold are
// replaced with sensible defaults (5s, 2, and 3 respectively).
func NewHealthChecker(cfg HealthCheckerConfig) *HealthChecker {
	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second
	}
	if cfg.HealthyThreshold == 0 {
		cfg.HealthyThreshold = 2
	}
	if cfg.UnhealthyThreshold == 0 {
		cfg.UnhealthyThreshold = 3
	}

	return &HealthChecker{
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
		timeout:            cfg.Timeout,
		healthyThreshold:   cfg.HealthyThreshold,
		unhealthyThreshold: cfg.UnhealthyThreshold,
		successCnt:         make(map[string]int),
		failureCnt:         make(map[string]int),
	}
}

// Check performs an HTTP GET health check against the instance's health endpoint
// (protocol://host:port/HealthCheckPath). The check result maps HTTP responses to
// statuses: 2xx -> healthy (or "Recovering" if below threshold), 503 -> maintenance,
// other codes or errors -> failed. Success/failure counters are updated to support
// threshold-based transitions that prevent status flapping from transient errors.
func (h *HealthChecker) Check(ctx context.Context, inst *Instance) *HealthCheckResult {
	start := time.Now()
	result := &HealthCheckResult{
		InstanceID: inst.ID,
		ServiceID:  inst.ServiceID,
		Timestamp:  start,
	}

	protocol := inst.Protocol
	if protocol == "" {
		protocol = "http"
	}

	path := inst.HealthCheckPath
	if path == "" {
		path = "/health"
	}

	url := protocol + "://" + inst.Host
	if inst.Port != 0 {
		url += ":" + itoa(inst.Port)
	}
	url += path

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		result.Status = StatusFailed
		result.Message = err.Error()
		result.Duration = time.Since(start)
		h.recordFailure(inst.ID)
		return result
	}

	resp, err := h.client.Do(req)
	result.Duration = time.Since(start)

	if err != nil {
		result.Status = StatusFailed
		result.Message = err.Error()
		h.recordFailure(inst.ID)
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		h.recordSuccess(inst.ID)
		if h.isHealthy(inst.ID) {
			result.Status = StatusRunning
			result.Message = "OK"
		} else {
			result.Status = StatusDegraded
			result.Message = "Recovering"
		}
	} else if resp.StatusCode == 503 {
		result.Status = StatusMaintenance
		result.Message = "Service unavailable"
		h.recordFailure(inst.ID)
	} else {
		result.Status = StatusFailed
		result.Message = "HTTP " + resp.Status
		h.recordFailure(inst.ID)
	}

	result.Details = map[string]any{
		"status_code": resp.StatusCode,
		"duration_ms": result.Duration.Milliseconds(),
	}

	return result
}

// recordSuccess increments the success counter and resets the failure counter
// for an instance. This implements one side of the threshold-based transition logic.
func (h *HealthChecker) recordSuccess(instanceID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.successCnt[instanceID]++
	h.failureCnt[instanceID] = 0
}

// recordFailure increments the failure counter and resets the success counter
// for an instance. A streak of failures triggers the unhealthy threshold.
func (h *HealthChecker) recordFailure(instanceID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.failureCnt[instanceID]++
	h.successCnt[instanceID] = 0
}

// isHealthy returns true if the instance has reached the healthyThreshold of
// consecutive successful health checks.
func (h *HealthChecker) isHealthy(instanceID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.successCnt[instanceID] >= h.healthyThreshold
}

// HTTP Handlers -- These methods expose Manager operations as REST endpoints.
// They decode JSON request bodies, delegate to the corresponding Manager methods,
// and return JSON responses. They are designed to be registered directly with
// http.HandleFunc or any compatible HTTP router.

// HandleServiceRegistration handles HTTP POST requests to register a new service.
// It expects a JSON-encoded RegistrationRequest body and returns a 201 Created
// response with the service ID on success.
func (m *Manager) HandleServiceRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, errMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var req RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := m.RegisterService(r.Context(), &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":     "registered",
		"service_id": req.Service.ID,
	})
}

// HandleServiceDeregistration handles HTTP DELETE or POST requests to remove a
// service. It expects a JSON-encoded DeregistrationRequest body and returns 404
// if the service does not exist.
func (m *Manager) HandleServiceDeregistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		http.Error(w, errMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var req DeregistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := m.DeregisterService(r.Context(), &req); err != nil {
		if err == ErrServiceNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":     "deregistered",
		"service_id": req.ServiceID,
	})
}

// HandleInstanceRegistration handles HTTP POST requests to register a new instance.
// It expects a JSON-encoded InstanceRegistrationRequest body and returns 201 Created
// with the instance and service IDs on success.
func (m *Manager) HandleInstanceRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, errMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var req InstanceRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := m.RegisterInstance(r.Context(), &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":      "registered",
		"instance_id": req.Instance.ID,
		"service_id":  req.Instance.ServiceID,
	})
}

// HandleServiceQuery handles HTTP GET or POST requests to query services. GET
// requests support query parameters (name, healthy, available); POST requests
// accept a JSON-encoded Query body for full filter control. Results are returned
// as a JSON array.
func (m *Manager) HandleServiceQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, errMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var q Query
	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		params := r.URL.Query()
		if name := params.Get("name"); name != "" {
			q.Names = []string{name}
		}
		if params.Get("healthy") == "true" {
			q.OnlyHealthy = true
		}
		if params.Get("available") == "true" {
			q.OnlyAvailable = true
		}
	}

	services := m.QueryServices(&q)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(services)
}

// HandleInstanceRenew handles HTTP POST or PUT requests to extend an instance's TTL.
// It expects a JSON body with instance_id and ttl fields. Returns 404 if the
// instance does not exist, enabling clients to detect stale heartbeat targets.
func (m *Manager) HandleInstanceRenew(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		http.Error(w, errMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		InstanceID string        `json:"instance_id"`
		TTL        time.Duration `json:"ttl"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := m.RenewInstance(r.Context(), req.InstanceID, req.TTL); err != nil {
		if err == ErrInstanceNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":      "renewed",
		"instance_id": req.InstanceID,
	})
}
