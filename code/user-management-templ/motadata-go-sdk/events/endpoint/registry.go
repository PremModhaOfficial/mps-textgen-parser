package endpoint

import (
	"context"
	"sort"
	"sync"
	"time"
)

// Registry manages endpoint registration and discovery
type Registry struct {
	mu        sync.RWMutex
	endpoints map[string]*Info

	// Indices for efficient lookups
	byService  map[string]map[string]struct{} // serviceID -> endpointIDs
	byInstance map[string]map[string]struct{} // instanceID -> endpointIDs
	byHost     map[string]map[string]struct{} // host -> endpointIDs
	byTag      map[string]map[string]struct{} // tag -> endpointIDs

	// Lifecycle
	done   <-chan struct{}
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Callbacks
	onRegistered   func(endpoint *Info)
	onDeregistered func(endpoint *Info, reason string)
	onUpdated      func(old, new *Info)
	onStatusChange func(endpoint *Info, oldStatus, newStatus Status)

	// Configuration
	cleanupInterval   time.Duration
	expirationEnabled bool
}

// RegistryConfig holds registry configuration
type RegistryConfig struct {
	// CleanupInterval is how often to check for expired endpoints
	CleanupInterval time.Duration

	// ExpirationEnabled enables automatic expiration of endpoints
	ExpirationEnabled bool
}

// DefaultRegistryConfig returns default registry configuration
func DefaultRegistryConfig() RegistryConfig {
	return RegistryConfig{
		CleanupInterval:   time.Minute,
		ExpirationEnabled: true,
	}
}

// NewRegistry creates a new endpoint registry
func NewRegistry(cfg ...RegistryConfig) *Registry {
	var config RegistryConfig
	if len(cfg) > 0 {
		config = cfg[0]
	} else {
		config = DefaultRegistryConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	r := &Registry{
		endpoints:         make(map[string]*Info),
		byService:         make(map[string]map[string]struct{}),
		byInstance:        make(map[string]map[string]struct{}),
		byHost:            make(map[string]map[string]struct{}),
		byTag:             make(map[string]map[string]struct{}),
		done:              ctx.Done(),
		cancel:            cancel,
		cleanupInterval:   config.CleanupInterval,
		expirationEnabled: config.ExpirationEnabled,
	}

	if config.ExpirationEnabled && config.CleanupInterval > 0 {
		r.startCleanupLoop()
	}

	return r
}

// startCleanupLoop starts the background cleanup goroutine
func (r *Registry) startCleanupLoop() {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()

		ticker := time.NewTicker(r.cleanupInterval)
		defer ticker.Stop()

		for {
			select {
			case <-r.done:
				return
			case <-ticker.C:
				r.cleanupExpired()
			}
		}
	}()
}

// cleanupExpired removes expired endpoints
func (r *Registry) cleanupExpired() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for id, ep := range r.endpoints {
		if !ep.ExpiresAt.IsZero() && now.After(ep.ExpiresAt) {
			_ = r.removeEndpointLocked(id, "expired")
		}
	}
}

// Register registers a new endpoint or updates an existing one
func (r *Registry) Register(endpoint *Info) error {
	if endpoint == nil {
		return ErrInvalidEndpoint
	}
	if endpoint.ID == "" {
		return ErrMissingEndpointID
	}
	if endpoint.ServiceID == "" {
		return ErrMissingServiceID
	}

	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	// Check for shutdown
	select {
	case <-r.done:
		return ErrRegistryShutdown
	default:
	}

	existing, exists := r.endpoints[endpoint.ID]
	if exists {
		// Update existing
		oldEndpoint := existing.Clone()

		// Preserve registration time
		endpoint.RegisteredAt = existing.RegisteredAt
		endpoint.UpdatedAt = now

		// Remove from old indices if service/instance/host/tags changed
		r.removeFromIndicesLocked(existing)

		// Store updated endpoint
		r.endpoints[endpoint.ID] = endpoint.Clone()

		// Add to indices
		r.addToIndicesLocked(endpoint)

		// Callback
		if r.onUpdated != nil {
			go r.onUpdated(oldEndpoint, endpoint)
		}

		// Status change callback
		if r.onStatusChange != nil && oldEndpoint.Status != endpoint.Status {
			go r.onStatusChange(endpoint, oldEndpoint.Status, endpoint.Status)
		}
	} else {
		// New registration
		endpoint.RegisteredAt = now
		endpoint.UpdatedAt = now

		// Clone and store
		stored := endpoint.Clone()
		r.endpoints[endpoint.ID] = stored

		// Add to indices
		r.addToIndicesLocked(endpoint)

		// Callback
		if r.onRegistered != nil {
			go r.onRegistered(stored)
		}
	}

	return nil
}

// RegisterWithTTL registers an endpoint with automatic expiration
func (r *Registry) RegisterWithTTL(endpoint *Info, ttl time.Duration) error {
	if ttl > 0 {
		endpoint.ExpiresAt = time.Now().Add(ttl)
	}
	return r.Register(endpoint)
}

// Deregister removes an endpoint from the registry
func (r *Registry) Deregister(endpointID string) error {
	return r.DeregisterWithReason(endpointID, "")
}

// DeregisterWithReason removes an endpoint with a reason
func (r *Registry) DeregisterWithReason(endpointID, reason string) error {
	if endpointID == "" {
		return ErrMissingEndpointID
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	return r.removeEndpointLocked(endpointID, reason)
}

// removeEndpointLocked removes an endpoint (must be called with lock held)
func (r *Registry) removeEndpointLocked(endpointID, reason string) error {
	endpoint, exists := r.endpoints[endpointID]
	if !exists {
		return ErrEndpointNotFound
	}

	// Remove from indices
	r.removeFromIndicesLocked(endpoint)

	// Remove from main map
	delete(r.endpoints, endpointID)

	// Callback
	if r.onDeregistered != nil {
		go r.onDeregistered(endpoint, reason)
	}

	return nil
}

// DeregisterByService removes all endpoints for a service
func (r *Registry) DeregisterByService(serviceID string) (int, error) {
	if serviceID == "" {
		return 0, ErrMissingServiceID
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	endpointIDs, exists := r.byService[serviceID]
	if !exists {
		return 0, nil
	}

	count := 0
	for id := range endpointIDs {
		if err := r.removeEndpointLocked(id, "service deregistered"); err == nil {
			count++
		}
	}

	return count, nil
}

// DeregisterByInstance removes all endpoints for an instance
func (r *Registry) DeregisterByInstance(instanceID string) (int, error) {
	if instanceID == "" {
		return 0, nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	endpointIDs, exists := r.byInstance[instanceID]
	if !exists {
		return 0, nil
	}

	count := 0
	for id := range endpointIDs {
		if err := r.removeEndpointLocked(id, "instance deregistered"); err == nil {
			count++
		}
	}

	return count, nil
}

// Get returns an endpoint by ID
func (r *Registry) Get(endpointID string) (*Info, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	endpoint, exists := r.endpoints[endpointID]
	if !exists {
		return nil, false
	}

	// Skip expired endpoints
	if !endpoint.ExpiresAt.IsZero() && time.Now().After(endpoint.ExpiresAt) {
		return nil, false
	}

	return endpoint.Clone(), true
}

// GetByService returns all endpoints for a service
func (r *Registry) GetByService(serviceID string) []*Info {
	r.mu.RLock()
	defer r.mu.RUnlock()

	endpointIDs, exists := r.byService[serviceID]
	if !exists {
		return nil
	}

	now := time.Now()
	result := make([]*Info, 0, len(endpointIDs))
	for id := range endpointIDs {
		if ep, ok := r.endpoints[id]; ok {
			// Skip expired
			if !ep.ExpiresAt.IsZero() && now.After(ep.ExpiresAt) {
				continue
			}
			result = append(result, ep.Clone())
		}
	}

	return result
}

// GetByInstance returns all endpoints for an instance
func (r *Registry) GetByInstance(instanceID string) []*Info {
	r.mu.RLock()
	defer r.mu.RUnlock()

	endpointIDs, exists := r.byInstance[instanceID]
	if !exists {
		return nil
	}

	now := time.Now()
	result := make([]*Info, 0, len(endpointIDs))
	for id := range endpointIDs {
		if ep, ok := r.endpoints[id]; ok {
			if !ep.ExpiresAt.IsZero() && now.After(ep.ExpiresAt) {
				continue
			}
			result = append(result, ep.Clone())
		}
	}

	return result
}

// GetByTag returns all endpoints with a specific tag
func (r *Registry) GetByTag(tag string) []*Info {
	r.mu.RLock()
	defer r.mu.RUnlock()

	endpointIDs, exists := r.byTag[tag]
	if !exists {
		return nil
	}

	now := time.Now()
	result := make([]*Info, 0, len(endpointIDs))
	for id := range endpointIDs {
		if ep, ok := r.endpoints[id]; ok {
			if !ep.ExpiresAt.IsZero() && now.After(ep.ExpiresAt) {
				continue
			}
			result = append(result, ep.Clone())
		}
	}

	return result
}

// GetHealthy returns all healthy endpoints
func (r *Registry) GetHealthy() []*Info {
	return r.Query(&Query{OnlyHealthy: true})
}

// GetAvailable returns all available endpoints (healthy or degraded)
func (r *Registry) GetAvailable() []*Info {
	return r.Query(&Query{OnlyAvailable: true})
}

// Query returns endpoints matching the query
func (r *Registry) Query(q *Query) []*Info {
	if q == nil {
		q = &Query{}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	now := time.Now()
	result := make([]*Info, 0)

	for _, ep := range r.endpoints {
		// Skip expired
		if !ep.ExpiresAt.IsZero() && now.After(ep.ExpiresAt) {
			continue
		}

		if q.Matches(ep) {
			result = append(result, ep.Clone())
		}
	}

	// Sort by priority (higher first), then weight (higher first)
	sort.Slice(result, func(i, j int) bool {
		if result[i].Priority != result[j].Priority {
			return result[i].Priority > result[j].Priority
		}
		return result[i].Weight > result[j].Weight
	})

	// Apply pagination
	if q.Offset > 0 {
		if q.Offset >= len(result) {
			return nil
		}
		result = result[q.Offset:]
	}
	if q.Limit > 0 && len(result) > q.Limit {
		result = result[:q.Limit]
	}

	return result
}

// All returns all registered endpoints
func (r *Registry) All() []*Info {
	r.mu.RLock()
	defer r.mu.RUnlock()

	now := time.Now()
	result := make([]*Info, 0, len(r.endpoints))
	for _, ep := range r.endpoints {
		if !ep.ExpiresAt.IsZero() && now.After(ep.ExpiresAt) {
			continue
		}
		result = append(result, ep.Clone())
	}

	return result
}

// Count returns the number of registered endpoints
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.endpoints)
}

// ServiceCount returns the number of unique services
func (r *Registry) ServiceCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.byService)
}

// Services returns all unique service IDs
func (r *Registry) Services() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	services := make([]string, 0, len(r.byService))
	for serviceID := range r.byService {
		services = append(services, serviceID)
	}
	return services
}

// Exists checks if an endpoint exists
func (r *Registry) Exists(endpointID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.endpoints[endpointID]
	return exists
}

// UpdateStatus updates the status of an endpoint
func (r *Registry) UpdateStatus(endpointID string, status Status) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	endpoint, exists := r.endpoints[endpointID]
	if !exists {
		return ErrEndpointNotFound
	}

	oldStatus := endpoint.Status
	endpoint.Status = status
	endpoint.LastHealthCheck = time.Now()
	endpoint.UpdatedAt = time.Now()

	if r.onStatusChange != nil && oldStatus != status {
		go r.onStatusChange(endpoint, oldStatus, status)
	}

	return nil
}

// UpdateHealth updates the health information of an endpoint
func (r *Registry) UpdateHealth(update *HealthUpdate) error {
	if update == nil || update.EndpointID == "" {
		return ErrInvalidEndpoint
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	endpoint, exists := r.endpoints[update.EndpointID]
	if !exists {
		return ErrEndpointNotFound
	}

	oldStatus := endpoint.Status
	endpoint.Status = update.Status
	endpoint.LastHealthCheck = update.CheckedAt
	endpoint.UpdatedAt = time.Now()

	if update.Status == StatusHealthy {
		endpoint.ConsecutiveFails = 0
	} else if update.Status == StatusUnhealthy {
		endpoint.ConsecutiveFails++
	}

	if r.onStatusChange != nil && oldStatus != update.Status {
		go r.onStatusChange(endpoint, oldStatus, update.Status)
	}

	return nil
}

// Renew extends the expiration time of an endpoint
func (r *Registry) Renew(endpointID string, ttl time.Duration) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	endpoint, exists := r.endpoints[endpointID]
	if !exists {
		return ErrEndpointNotFound
	}

	endpoint.ExpiresAt = time.Now().Add(ttl)
	endpoint.UpdatedAt = time.Now()

	return nil
}

// Touch updates the last activity time without changing expiration
func (r *Registry) Touch(endpointID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	endpoint, exists := r.endpoints[endpointID]
	if !exists {
		return ErrEndpointNotFound
	}

	endpoint.UpdatedAt = time.Now()
	return nil
}

// Shutdown gracefully shuts down the registry
func (r *Registry) Shutdown(ctx context.Context) error {
	r.cancel()

	// Wait for cleanup loop with timeout
	done := make(chan struct{})
	go func() {
		r.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

// Clear removes all endpoints
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.endpoints = make(map[string]*Info)
	r.byService = make(map[string]map[string]struct{})
	r.byInstance = make(map[string]map[string]struct{})
	r.byHost = make(map[string]map[string]struct{})
	r.byTag = make(map[string]map[string]struct{})
}

// Callbacks

// OnRegistered sets callback for endpoint registration
func (r *Registry) OnRegistered(cb func(endpoint *Info)) {
	r.onRegistered = cb
}

// OnDeregistered sets callback for endpoint deregistration
func (r *Registry) OnDeregistered(cb func(endpoint *Info, reason string)) {
	r.onDeregistered = cb
}

// OnUpdated sets callback for endpoint updates
func (r *Registry) OnUpdated(cb func(old, new *Info)) {
	r.onUpdated = cb
}

// OnStatusChange sets callback for status changes
func (r *Registry) OnStatusChange(cb func(endpoint *Info, oldStatus, newStatus Status)) {
	r.onStatusChange = cb
}

// Index management helpers

func (r *Registry) addToIndicesLocked(ep *Info) {
	// Add to service index
	if ep.ServiceID != "" {
		if r.byService[ep.ServiceID] == nil {
			r.byService[ep.ServiceID] = make(map[string]struct{})
		}
		r.byService[ep.ServiceID][ep.ID] = struct{}{}
	}

	// Add to instance index
	if ep.InstanceID != "" {
		if r.byInstance[ep.InstanceID] == nil {
			r.byInstance[ep.InstanceID] = make(map[string]struct{})
		}
		r.byInstance[ep.InstanceID][ep.ID] = struct{}{}
	}

	// Add to host index
	if ep.Host != "" {
		if r.byHost[ep.Host] == nil {
			r.byHost[ep.Host] = make(map[string]struct{})
		}
		r.byHost[ep.Host][ep.ID] = struct{}{}
	}

	// Add to tag indices
	for _, tag := range ep.Tags {
		if r.byTag[tag] == nil {
			r.byTag[tag] = make(map[string]struct{})
		}
		r.byTag[tag][ep.ID] = struct{}{}
	}
}

// removeFromIndex removes endpointID from an index map and cleans up the key if empty.
func removeFromIndex(index map[string]map[string]struct{}, key, endpointID string) {
	if key == "" {
		return
	}
	entries, ok := index[key]
	if !ok {
		return
	}
	delete(entries, endpointID)
	if len(entries) == 0 {
		delete(index, key)
	}
}

func (r *Registry) removeFromIndicesLocked(ep *Info) {
	removeFromIndex(r.byService, ep.ServiceID, ep.ID)
	removeFromIndex(r.byInstance, ep.InstanceID, ep.ID)
	removeFromIndex(r.byHost, ep.Host, ep.ID)

	for _, tag := range ep.Tags {
		removeFromIndex(r.byTag, tag, ep.ID)
	}
}
