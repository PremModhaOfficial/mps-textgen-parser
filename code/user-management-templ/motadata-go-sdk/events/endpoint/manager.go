package endpoint

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
	trackerEndpointKey = "endpoint.id"
	tracerServiceKey   = "service.id"
	errMethodNotAllowed = "Method not allowed"
)

// Manager handles endpoint lifecycle, health checking, and event propagation
type Manager struct {
	registry *Registry

	// Publishing
	publisher       core.Publisher
	publishSubjects PublishSubjects

	// Health checking
	healthChecker   *HealthChecker
	healthConfig    HealthCheckConfig
	healthEnabled   bool
	healthEndpoints map[string]struct{}
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

// PublishSubjects defines the subjects for publishing endpoint events
type PublishSubjects struct {
	Registered   string
	Deregistered string
	Updated      string
	HealthChange string
}

// DefaultPublishSubjects returns default publish subjects
func DefaultPublishSubjects() PublishSubjects {
	return PublishSubjects{
		Registered:   "endpoints.registered",
		Deregistered: "endpoints.deregistered",
		Updated:      "endpoints.updated",
		HealthChange: "endpoints.health",
	}
}

// HealthCheckConfig configures health checking behavior
type HealthCheckConfig struct {
	// Enabled enables automatic health checking
	Enabled bool

	// Interval is how often to check endpoint health
	Interval time.Duration

	// Timeout is the timeout for each health check
	Timeout time.Duration

	// HealthyThreshold is consecutive successes to mark healthy
	HealthyThreshold int

	// UnhealthyThreshold is consecutive failures to mark unhealthy
	UnhealthyThreshold int

	// InitialDelay before starting health checks
	InitialDelay time.Duration
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
	// Registry is the endpoint registry to use (required)
	Registry *Registry

	// Publisher for publishing events (optional)
	Publisher core.Publisher

	// PublishSubjects configures event subjects
	PublishSubjects PublishSubjects

	// HealthCheck configuration
	HealthCheck HealthCheckConfig

	// ErrorBufferSize is the size of the error channel buffer
	ErrorBufferSize int
}

// NewManager creates a new endpoint manager
func NewManager(cfg ManagerConfig) (*Manager, error) {
	if cfg.Registry == nil {
		cfg.Registry = NewRegistry()
	}

	if cfg.PublishSubjects.Registered == "" {
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
		healthEndpoints: make(map[string]struct{}),
		done:            ctx.Done(),
		cancel:          cancel,
		errorChan:       make(chan error, cfg.ErrorBufferSize),
	}

	// Create health checker
	m.healthChecker = NewHealthChecker(HealthCheckerConfig{
		Timeout:            cfg.HealthCheck.Timeout,
		HealthyThreshold:   cfg.HealthCheck.HealthyThreshold,
		UnhealthyThreshold: cfg.HealthCheck.UnhealthyThreshold,
	})

	// Set up registry callbacks
	m.setupCallbacks()

	// Start health check loop if enabled
	if cfg.HealthCheck.Enabled {
		m.startHealthCheckLoop()
	}

	logger.Info(ctx, "endpoint manager initialized",
		logger.Bool("health_check_enabled", cfg.HealthCheck.Enabled),
		logger.Duration("health_check_interval", cfg.HealthCheck.Interval),
	)

	return m, nil
}

// setupCallbacks configures registry callbacks
func (m *Manager) setupCallbacks() {
	m.registry.OnRegistered(func(ep *Info) {
		m.publishEvent(EventRegistered, ep, StatusUnknown, ep.Status, "")
		if m.healthEnabled {
			m.addToHealthCheck(ep.ID)
		}
	})

	m.registry.OnDeregistered(func(ep *Info, reason string) {
		m.publishEvent(EventDeregistered, ep, ep.Status, StatusOffline, reason)
		if m.healthEnabled {
			m.removeFromHealthCheck(ep.ID)
		}
	})

	m.registry.OnUpdated(func(old, new *Info) {
		m.publishEvent(EventUpdated, new, old.Status, new.Status, "")
	})

	m.registry.OnStatusChange(func(ep *Info, oldStatus, newStatus Status) {
		eventType := EventHealthy
		if newStatus == StatusUnhealthy {
			eventType = EventUnhealthy
		}
		m.publishEvent(eventType, ep, oldStatus, newStatus, "")
	})
}

// startHealthCheckLoop starts the background health check goroutine
func (m *Manager) startHealthCheckLoop() {
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()

		// Initial delay
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

// performHealthChecks runs health checks on all registered endpoints
func (m *Manager) performHealthChecks() {
	ctx, span := tracer.Start(context.Background(), "endpoint.health_check_cycle")
	defer span.End()

	m.healthMu.RLock()
	endpointIDs := make([]string, 0, len(m.healthEndpoints))
	for id := range m.healthEndpoints {
		endpointIDs = append(endpointIDs, id)
	}
	m.healthMu.RUnlock()

	span.SetAttributes(tracer.IntAttr("endpoints.count", len(endpointIDs)))

	for _, id := range endpointIDs {
		ep, exists := m.registry.Get(id)
		if !exists {
			continue
		}

		// Skip if no health check path configured
		if ep.HealthCheckPath == "" {
			continue
		}

		// Perform health check
		result := m.healthChecker.Check(ctx, ep)

		// Update registry
		if err := m.registry.UpdateHealth(&HealthUpdate{
			EndpointID: id,
			Status:     result.Status,
			Message:    result.Message,
			CheckedAt:  result.Timestamp,
			Details:    result.Details,
		}); err != nil {
			logger.Error(ctx, "failed to update endpoint health",
				logger.String("endpoint_id", id),
				logger.Err(err),
			)
			m.dispatchError(err)
		}
	}

	span.SetOK()
	logger.Debug(ctx, "endpoint health check cycle completed",
		logger.Int("endpoints_checked", len(endpointIDs)),
	)
}

// addToHealthCheck adds an endpoint to health checking
func (m *Manager) addToHealthCheck(endpointID string) {
	m.healthMu.Lock()
	m.healthEndpoints[endpointID] = struct{}{}
	m.healthMu.Unlock()
}

// removeFromHealthCheck removes an endpoint from health checking
func (m *Manager) removeFromHealthCheck(endpointID string) {
	m.healthMu.Lock()
	delete(m.healthEndpoints, endpointID)
	m.healthMu.Unlock()
}

// publishEvent publishes an endpoint event
func (m *Manager) publishEvent(eventType EventType, ep *Info, oldStatus, newStatus Status, reason string) {
	if m.publisher == nil {
		return
	}

	event := &Event{
		Type:       eventType,
		EndpointID: ep.ID,
		ServiceID:  ep.ServiceID,
		Endpoint:   ep,
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
	case EventRegistered:
		subject = m.publishSubjects.Registered
	case EventDeregistered:
		subject = m.publishSubjects.Deregistered
	case EventUpdated:
		subject = m.publishSubjects.Updated
	case EventHealthy, EventUnhealthy:
		subject = m.publishSubjects.HealthChange
	default:
		subject = m.publishSubjects.Updated
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

// dispatchError dispatches an error
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
		// Channel full, drop error
	}
}

// Register registers a new endpoint
func (m *Manager) Register(ctx context.Context, req *RegistrationRequest) error {
	ctx, span := tracer.Start(ctx, "endpoint.register")
	defer span.End()

	// Validate first to avoid nil pointer access
	if err := req.Validate(); err != nil {
		span.SetError(err)
		logger.Error(ctx, "endpoint registration validation failed",
			logger.Err(err),
		)
		return err
	}

	// Now safe to access endpoint fields
	span.SetAttributes(
		tracer.StringAttr(trackerEndpointKey, req.Endpoint.ID),
		tracer.StringAttr(tracerServiceKey, req.Endpoint.ServiceID),
	)

	var err error
	// Set expiration if TTL provided
	if req.TTL > 0 {
		err = m.registry.RegisterWithTTL(req.Endpoint, req.TTL)
		span.SetAttributes(tracer.StringAttr("ttl", req.TTL.String()))
	} else {
		err = m.registry.Register(req.Endpoint)
	}

	if err != nil {
		span.SetError(err)
		logger.Error(ctx, "failed to register endpoint",
			logger.String("endpoint_id", req.Endpoint.ID),
			logger.Err(err),
		)
		return err
	}

	span.SetOK()
	logger.Info(ctx, "endpoint registered",
		logger.String("endpoint_id", req.Endpoint.ID),
		logger.String("service_id", req.Endpoint.ServiceID),
	)
	return nil
}

// RegisterEndpoint is a convenience method for registering an endpoint directly
func (m *Manager) RegisterEndpoint(ctx context.Context, endpoint *Info) error {
	ctx, span := tracer.Start(ctx, "endpoint.register_direct",
		tracer.WithAttributes(
			tracer.StringAttr(trackerEndpointKey, endpoint.ID),
			tracer.StringAttr(tracerServiceKey, endpoint.ServiceID),
		),
	)
	defer span.End()

	if err := m.registry.Register(endpoint); err != nil {
		span.SetError(err)
		logger.Error(ctx, "failed to register endpoint",
			logger.String("endpoint_id", endpoint.ID),
			logger.Err(err),
		)
		return err
	}

	span.SetOK()
	logger.Info(ctx, "endpoint registered",
		logger.String("endpoint_id", endpoint.ID),
		logger.String("service_id", endpoint.ServiceID),
	)
	return nil
}

// RegisterWithTTL registers an endpoint with TTL
func (m *Manager) RegisterWithTTL(ctx context.Context, endpoint *Info, ttl time.Duration) error {
	ctx, span := tracer.Start(ctx, "endpoint.register_ttl",
		tracer.WithAttributes(
			tracer.StringAttr(trackerEndpointKey, endpoint.ID),
			tracer.StringAttr(tracerServiceKey, endpoint.ServiceID),
			tracer.StringAttr("ttl", ttl.String()),
		),
	)
	defer span.End()

	if err := m.registry.RegisterWithTTL(endpoint, ttl); err != nil {
		span.SetError(err)
		logger.Error(ctx, "failed to register endpoint with TTL",
			logger.String("endpoint_id", endpoint.ID),
			logger.Err(err),
		)
		return err
	}

	span.SetOK()
	logger.Info(ctx, "endpoint registered with TTL",
		logger.String("endpoint_id", endpoint.ID),
		logger.String("service_id", endpoint.ServiceID),
		logger.Duration("ttl", ttl),
	)
	return nil
}

// Deregister removes an endpoint
func (m *Manager) Deregister(ctx context.Context, req *DeregistrationRequest) error {
	ctx, span := tracer.Start(ctx, "endpoint.deregister")
	defer span.End()

	// Validate first
	if err := req.Validate(); err != nil {
		span.SetError(err)
		logger.Error(ctx, "endpoint deregistration validation failed",
			logger.Err(err),
		)
		return err
	}

	// Now safe to access fields
	span.SetAttributes(tracer.StringAttr(trackerEndpointKey, req.EndpointID))

	if err := m.registry.DeregisterWithReason(req.EndpointID, req.Reason); err != nil {
		span.SetError(err)
		logger.Error(ctx, "failed to deregister endpoint",
			logger.String("endpoint_id", req.EndpointID),
			logger.Err(err),
		)
		return err
	}

	span.SetOK()
	logger.Info(ctx, "endpoint deregistered",
		logger.String("endpoint_id", req.EndpointID),
		logger.String("reason", req.Reason),
	)
	return nil
}

// DeregisterEndpoint is a convenience method for deregistering an endpoint
func (m *Manager) DeregisterEndpoint(ctx context.Context, endpointID string) error {
	ctx, span := tracer.Start(ctx, "endpoint.deregister_direct",
		tracer.WithAttributes(
			tracer.StringAttr(trackerEndpointKey, endpointID),
		),
	)
	defer span.End()

	if err := m.registry.Deregister(endpointID); err != nil {
		span.SetError(err)
		logger.Error(ctx, "failed to deregister endpoint",
			logger.String("endpoint_id", endpointID),
			logger.Err(err),
		)
		return err
	}

	span.SetOK()
	logger.Info(ctx, "endpoint deregistered",
		logger.String("endpoint_id", endpointID),
	)
	return nil
}

// DeregisterService removes all endpoints for a service
func (m *Manager) DeregisterService(ctx context.Context, serviceID string) (int, error) {
	ctx, span := tracer.Start(ctx, "endpoint.deregister_service",
		tracer.WithAttributes(
			tracer.StringAttr(tracerServiceKey, serviceID),
		),
	)
	defer span.End()

	count, err := m.registry.DeregisterByService(serviceID)
	if err != nil {
		span.SetError(err)
		logger.Error(ctx, "failed to deregister endpoints for service",
			logger.String("service_id", serviceID),
			logger.Err(err),
		)
		return count, err
	}

	span.SetAttributes(tracer.IntAttr("endpoints.removed", count))
	span.SetOK()
	logger.Info(ctx, "endpoints deregistered for service",
		logger.String("service_id", serviceID),
		logger.Int("count", count),
	)
	return count, nil
}

// DeregisterInstance removes all endpoints for an instance
func (m *Manager) DeregisterInstance(ctx context.Context, instanceID string) (int, error) {
	ctx, span := tracer.Start(ctx, "endpoint.deregister_instance",
		tracer.WithAttributes(
			tracer.StringAttr("instance.id", instanceID),
		),
	)
	defer span.End()

	count, err := m.registry.DeregisterByInstance(instanceID)
	if err != nil {
		span.SetError(err)
		logger.Error(ctx, "failed to deregister endpoints for instance",
			logger.String("instance_id", instanceID),
			logger.Err(err),
		)
		return count, err
	}

	span.SetAttributes(tracer.IntAttr("endpoints.removed", count))
	span.SetOK()
	logger.Info(ctx, "endpoints deregistered for instance",
		logger.String("instance_id", instanceID),
		logger.Int("count", count),
	)
	return count, nil
}

// Get returns an endpoint by ID
func (m *Manager) Get(endpointID string) (*Info, bool) {
	return m.registry.Get(endpointID)
}

// Query returns endpoints matching the query
func (m *Manager) Query(q *Query) []*Info {
	return m.registry.Query(q)
}

// GetByService returns all endpoints for a service
func (m *Manager) GetByService(serviceID string) []*Info {
	return m.registry.GetByService(serviceID)
}

// GetByInstance returns all endpoints for an instance
func (m *Manager) GetByInstance(instanceID string) []*Info {
	return m.registry.GetByInstance(instanceID)
}

// GetHealthy returns all healthy endpoints
func (m *Manager) GetHealthy() []*Info {
	return m.registry.GetHealthy()
}

// GetAvailable returns all available endpoints
func (m *Manager) GetAvailable() []*Info {
	return m.registry.GetAvailable()
}

// Discover finds endpoints matching criteria (alias for Query with common patterns)
func (m *Manager) Discover(serviceID string, opts ...DiscoverOption) []*Info {
	q := &Query{
		ServiceIDs:    []string{serviceID},
		OnlyAvailable: true,
	}

	for _, opt := range opts {
		opt(q)
	}

	return m.registry.Query(q)
}

// DiscoverOption modifies a discover query
type DiscoverOption func(*Query)

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

// WithProtocol filters by protocol
func WithProtocol(protocol Protocol) DiscoverOption {
	return func(q *Query) {
		q.Protocol = protocol
	}
}

// WithRegion filters by region
func WithRegion(region string) DiscoverOption {
	return func(q *Query) {
		q.Region = region
	}
}

// OnlyHealthy filters to only healthy endpoints
func OnlyHealthy() DiscoverOption {
	return func(q *Query) {
		q.OnlyHealthy = true
		q.OnlyAvailable = false
	}
}

// WithLimit limits results
func WithLimit(limit int) DiscoverOption {
	return func(q *Query) {
		q.Limit = limit
	}
}

// Renew extends endpoint TTL
func (m *Manager) Renew(ctx context.Context, endpointID string, ttl time.Duration) error {
	return m.registry.Renew(endpointID, ttl)
}

// UpdateStatus updates endpoint status
func (m *Manager) UpdateStatus(ctx context.Context, endpointID string, status Status) error {
	return m.registry.UpdateStatus(endpointID, status)
}

// All returns all endpoints
func (m *Manager) All() []*Info {
	return m.registry.All()
}

// Count returns endpoint count
func (m *Manager) Count() int {
	return m.registry.Count()
}

// Services returns all service IDs
func (m *Manager) Services() []string {
	return m.registry.Services()
}

// Registry returns the underlying registry
func (m *Manager) Registry() *Registry {
	return m.registry
}

// SetPublisher sets the publisher for events
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
	ctx, span := tracer.Start(ctx, "endpoint.manager.shutdown")
	defer span.End()

	logger.Info(ctx, "initiating endpoint manager shutdown",
		logger.Int("endpoint_count", m.Count()),
	)

	m.cancel()

	// Wait for goroutines
	done := make(chan struct{})
	go func() {
		m.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		span.SetError(ctx.Err())
		logger.Warn(ctx, "endpoint manager shutdown timeout")
		return ctx.Err()
	case <-done:
	}

	// Shutdown registry
	if err := m.registry.Shutdown(ctx); err != nil {
		span.SetError(err)
		logger.Error(ctx, "failed to shutdown endpoint registry",
			logger.Err(err),
		)
		return err
	}

	close(m.errorChan)

	span.SetOK()
	logger.Info(ctx, "endpoint manager shutdown complete")
	return nil
}

// HealthChecker performs health checks on endpoints
type HealthChecker struct {
	client             *http.Client
	timeout            time.Duration
	healthyThreshold   int
	unhealthyThreshold int

	// Track consecutive results per endpoint
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
	EndpointID string
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

// Check performs a health check on an endpoint
func (h *HealthChecker) Check(ctx context.Context, ep *Info) *HealthCheckResult {
	start := time.Now()
	result := &HealthCheckResult{
		EndpointID: ep.ID,
		Timestamp:  start,
	}

	// Build health check URL
	port := ep.HealthCheckPort
	if port == 0 {
		port = ep.Port
	}

	protocol := ep.Protocol
	if protocol == "" {
		protocol = ProtocolHTTP
	}

	path := ep.HealthCheckPath
	if path == "" {
		path = "/health"
	}

	url := string(protocol) + "://" + ep.Host
	if port != 0 {
		url += ":" + itoa(port)
	}
	url += path

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		result.Status = StatusUnhealthy
		result.Message = err.Error()
		result.Duration = time.Since(start)
		h.recordFailure(ep.ID)
		return result
	}

	// Perform request
	resp, err := h.client.Do(req)
	result.Duration = time.Since(start)

	if err != nil {
		result.Status = StatusUnhealthy
		result.Message = err.Error()
		h.recordFailure(ep.ID)
		return result
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		h.recordSuccess(ep.ID)
		if h.isHealthy(ep.ID) {
			result.Status = StatusHealthy
			result.Message = "OK"
		} else {
			result.Status = StatusDegraded
			result.Message = "Recovering"
		}
	} else if resp.StatusCode == 503 {
		result.Status = StatusMaintenance
		result.Message = "Service unavailable"
		h.recordFailure(ep.ID)
	} else {
		result.Status = StatusUnhealthy
		result.Message = "HTTP " + resp.Status
		h.recordFailure(ep.ID)
	}

	result.Details = map[string]any{
		"status_code": resp.StatusCode,
		"duration_ms": result.Duration.Milliseconds(),
	}

	return result
}

func (h *HealthChecker) recordSuccess(endpointID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.successCnt[endpointID]++
	h.failureCnt[endpointID] = 0
}

func (h *HealthChecker) recordFailure(endpointID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.failureCnt[endpointID]++
	h.successCnt[endpointID] = 0
}

func (h *HealthChecker) isHealthy(endpointID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return h.successCnt[endpointID] >= h.healthyThreshold
}

// HandleRegistration is an HTTP handler for endpoint registration
func (m *Manager) HandleRegistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, errMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var req RegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := m.Register(r.Context(), &req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":      "registered",
		"endpoint_id": req.Endpoint.ID,
	})
}

// HandleDeregistration is an HTTP handler for endpoint deregistration
func (m *Manager) HandleDeregistration(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		http.Error(w, errMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var req DeregistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := m.Deregister(r.Context(), &req); err != nil {
		if err == ErrEndpointNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":      "deregistered",
		"endpoint_id": req.EndpointID,
	})
}

// HandleQuery is an HTTP handler for querying endpoints
func (m *Manager) HandleQuery(w http.ResponseWriter, r *http.Request) {
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
		// Parse query parameters
		params := r.URL.Query()
		if serviceID := params.Get("service_id"); serviceID != "" {
			q.ServiceIDs = []string{serviceID}
		}
		if params.Get("healthy") == "true" {
			q.OnlyHealthy = true
		}
		if params.Get("available") == "true" {
			q.OnlyAvailable = true
		}
	}

	endpoints := m.Query(&q)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(endpoints)
}

// HandleRenew is an HTTP handler for renewing endpoint TTL
func (m *Manager) HandleRenew(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		http.Error(w, errMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		EndpointID string        `json:"endpoint_id"`
		TTL        time.Duration `json:"ttl"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := m.Renew(r.Context(), req.EndpointID, req.TTL); err != nil {
		if err == ErrEndpointNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"status":      "renewed",
		"endpoint_id": req.EndpointID,
	})
}
