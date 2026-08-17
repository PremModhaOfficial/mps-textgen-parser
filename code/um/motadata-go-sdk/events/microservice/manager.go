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

// Manager handles microservice lifecycle, health checking, and event propagation
type Manager struct {
	registry *Registry

	// Publishing
	publisher       core.Publisher
	publishSubjects PublishSubjects

	// Health checking
	healthChecker     *HealthChecker
	healthConfig      HealthCheckConfig
	healthEnabled     bool
	healthInstances   map[string]struct{}
	healthMu          sync.RWMutex

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Error handling
	errorChan chan error
	errorCb   func(error)
	errorMu   sync.RWMutex
}

// PublishSubjects defines the subjects for publishing events
type PublishSubjects struct {
	ServiceRegistered   string
	ServiceDeregistered string
	ServiceUpdated      string
	ServiceHealth       string
	InstanceRegistered  string
	InstanceDeregistered string
	InstanceHealth      string
}

// DefaultPublishSubjects returns default publish subjects
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

// HealthCheckConfig configures health checking behavior
type HealthCheckConfig struct {
	Enabled            bool
	Interval           time.Duration
	Timeout            time.Duration
	HealthyThreshold   int
	UnhealthyThreshold int
	InitialDelay       time.Duration
}

// DefaultHealthCheckConfig returns default health check configuration
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

// ManagerConfig holds manager configuration
type ManagerConfig struct {
	Registry        *Registry
	Publisher       core.Publisher
	PublishSubjects PublishSubjects
	HealthCheck     HealthCheckConfig
	ErrorBufferSize int
}

// NewManager creates a new microservice manager
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
		ctx:             ctx,
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

// setupCallbacks configures registry callbacks
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

// startHealthCheckLoop starts the background health check goroutine
func (m *Manager) startHealthCheckLoop() {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()

		if m.healthConfig.InitialDelay > 0 {
			select {
			case <-m.ctx.Done():
				return
			case <-time.After(m.healthConfig.InitialDelay):
			}
		}

		ticker := time.NewTicker(m.healthConfig.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-m.ctx.Done():
				return
			case <-ticker.C:
				m.performHealthChecks()
			}
		}
	}()
}

// performHealthChecks runs health checks on all registered instances
func (m *Manager) performHealthChecks() {
	ctx, span := tracer.Start(m.ctx, "microservice.health_check_cycle")
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

func (m *Manager) addToHealthCheck(instanceID string) {
	m.healthMu.Lock()
	m.healthInstances[instanceID] = struct{}{}
	m.healthMu.Unlock()
}

func (m *Manager) removeFromHealthCheck(instanceID string) {
	m.healthMu.Lock()
	delete(m.healthInstances, instanceID)
	m.healthMu.Unlock()
}

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

// RegisterService registers a new service
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
		tracer.StringAttr("service.id", req.Service.ID),
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

// RegisterServiceDirect is a convenience method
func (m *Manager) RegisterServiceDirect(ctx context.Context, service *Info) error {
	ctx, span := tracer.Start(ctx, "microservice.register_service_direct",
		tracer.WithAttributes(
			tracer.StringAttr("service.id", service.ID),
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

// DeregisterService removes a service
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
	span.SetAttributes(tracer.StringAttr("service.id", req.ServiceID))

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

// DeregisterServiceDirect is a convenience method
func (m *Manager) DeregisterServiceDirect(ctx context.Context, serviceID string) error {
	ctx, span := tracer.Start(ctx, "microservice.deregister_service_direct",
		tracer.WithAttributes(
			tracer.StringAttr("service.id", serviceID),
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

// RegisterInstance registers a new instance
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
		tracer.StringAttr("instance.id", req.Instance.ID),
		tracer.StringAttr("service.id", req.Instance.ServiceID),
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

// RegisterInstanceDirect is a convenience method
func (m *Manager) RegisterInstanceDirect(ctx context.Context, instance *Instance) error {
	ctx, span := tracer.Start(ctx, "microservice.register_instance_direct",
		tracer.WithAttributes(
			tracer.StringAttr("instance.id", instance.ID),
			tracer.StringAttr("service.id", instance.ServiceID),
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

// RegisterInstanceWithTTL registers an instance with TTL
func (m *Manager) RegisterInstanceWithTTL(ctx context.Context, instance *Instance, ttl time.Duration) error {
	ctx, span := tracer.Start(ctx, "microservice.register_instance_ttl",
		tracer.WithAttributes(
			tracer.StringAttr("instance.id", instance.ID),
			tracer.StringAttr("service.id", instance.ServiceID),
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

// DeregisterInstance removes an instance
func (m *Manager) DeregisterInstance(ctx context.Context, instanceID string) error {
	ctx, span := tracer.Start(ctx, "microservice.deregister_instance",
		tracer.WithAttributes(
			tracer.StringAttr("instance.id", instanceID),
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

// GetService returns a service by ID
func (m *Manager) GetService(serviceID string) (*Info, bool) {
	return m.registry.GetService(serviceID)
}

// GetServiceByName returns services by name
func (m *Manager) GetServiceByName(name string) []*Info {
	return m.registry.GetServiceByName(name)
}

// GetInstance returns an instance by ID
func (m *Manager) GetInstance(instanceID string) (*Instance, bool) {
	return m.registry.GetInstance(instanceID)
}

// GetInstances returns all instances for a service
func (m *Manager) GetInstances(serviceID string) []*Instance {
	return m.registry.GetInstances(serviceID)
}

// GetHealthyInstances returns healthy instances
func (m *Manager) GetHealthyInstances(serviceID string) []*Instance {
	return m.registry.GetHealthyInstances(serviceID)
}

// GetAvailableInstances returns available instances
func (m *Manager) GetAvailableInstances(serviceID string) []*Instance {
	return m.registry.GetAvailableInstances(serviceID)
}

// QueryServices queries for services
func (m *Manager) QueryServices(q *Query) []*Info {
	return m.registry.QueryServices(q)
}

// Discover finds services matching criteria
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

// DiscoverOption modifies a discover query
type DiscoverOption func(*Query)

// WithType filters by service type
func WithType(t ServiceType) DiscoverOption {
	return func(q *Query) {
		q.Types = []ServiceType{t}
	}
}

// WithCapabilities filters by capabilities
func WithCapabilities(caps ...string) DiscoverOption {
	return func(q *Query) {
		q.Capabilities = caps
	}
}

// WithTags filters by tags
func WithTags(tags ...string) DiscoverOption {
	return func(q *Query) {
		q.Tags = tags
	}
}

// WithVersion filters by version
func WithVersion(version string) DiscoverOption {
	return func(q *Query) {
		q.Version = version
	}
}

// OnlyHealthy filters to only healthy services
func OnlyHealthy() DiscoverOption {
	return func(q *Query) {
		q.OnlyHealthy = true
		q.OnlyAvailable = false
	}
}

// WithMinInstances requires minimum instances
func WithMinInstances(min int) DiscoverOption {
	return func(q *Query) {
		q.MinInstances = min
	}
}

// WithLimit limits results
func WithLimit(limit int) DiscoverOption {
	return func(q *Query) {
		q.Limit = limit
	}
}

// DiscoverInstances finds instances for a service
func (m *Manager) DiscoverInstances(serviceID string, onlyHealthy bool) []*Instance {
	if onlyHealthy {
		return m.registry.GetHealthyInstances(serviceID)
	}
	return m.registry.GetAvailableInstances(serviceID)
}

// RenewInstance extends instance TTL
func (m *Manager) RenewInstance(ctx context.Context, instanceID string, ttl time.Duration) error {
	return m.registry.RenewInstance(instanceID, ttl)
}

// UpdateServiceStatus updates service status
func (m *Manager) UpdateServiceStatus(ctx context.Context, serviceID string, status Status) error {
	return m.registry.UpdateServiceStatus(serviceID, status)
}

// UpdateInstanceStatus updates instance status
func (m *Manager) UpdateInstanceStatus(ctx context.Context, instanceID string, status Status) error {
	return m.registry.UpdateInstanceStatus(instanceID, status)
}

// AllServices returns all services
func (m *Manager) AllServices() []*Info {
	return m.registry.AllServices()
}

// AllInstances returns all instances
func (m *Manager) AllInstances() []*Instance {
	return m.registry.AllInstances()
}

// ServiceCount returns service count
func (m *Manager) ServiceCount() int {
	return m.registry.ServiceCount()
}

// InstanceCount returns instance count
func (m *Manager) InstanceCount() int {
	return m.registry.InstanceCount()
}

// Registry returns the underlying registry
func (m *Manager) Registry() *Registry {
	return m.registry
}

// SetPublisher sets the publisher
func (m *Manager) SetPublisher(publisher core.Publisher) {
	m.publisher = publisher
}

// Errors returns the error channel
func (m *Manager) Errors() <-chan error {
	return m.errorChan
}

// OnError sets an error callback
func (m *Manager) OnError(cb func(error)) {
	m.errorMu.Lock()
	m.errorCb = cb
	m.errorMu.Unlock()
}

// Shutdown gracefully shuts down the manager
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

// HealthChecker performs health checks on instances
type HealthChecker struct {
	client             *http.Client
	timeout            time.Duration
	healthyThreshold   int
	unhealthyThreshold int

	mu         sync.RWMutex
	successCnt map[string]int
	failureCnt map[string]int
}

// HealthCheckerConfig configures the health checker
type HealthCheckerConfig struct {
	Timeout            time.Duration
	HealthyThreshold   int
	UnhealthyThreshold int
}

// HealthCheckResult represents the result of a health check
type HealthCheckResult struct {
	InstanceID string
	ServiceID  string
	Status     Status
	Message    string
	Timestamp  time.Time
	Duration   time.Duration
	Details    map[string]any
}

// NewHealthChecker creates a new health checker
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

// Check performs a health check on an instance
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

func (h *HealthChecker) recordSuccess(instanceID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.successCnt[instanceID]++
	h.failureCnt[instanceID] = 0
}

func (h *HealthChecker) recordFailure(instanceID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.failureCnt[instanceID]++
	h.successCnt[instanceID] = 0
}

func (h *HealthChecker) isHealthy(instanceID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.successCnt[instanceID] >= h.healthyThreshold
}

// HTTP Handlers

// HandleServiceRegistration handles service registration HTTP requests
func (m *Manager) HandleServiceRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

// HandleServiceDeregistration handles service deregistration HTTP requests
func (m *Manager) HandleServiceDeregistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

// HandleInstanceRegistration handles instance registration HTTP requests
func (m *Manager) HandleInstanceRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

// HandleServiceQuery handles service query HTTP requests
func (m *Manager) HandleServiceQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

// HandleInstanceRenew handles instance TTL renewal HTTP requests
func (m *Manager) HandleInstanceRenew(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
