package microservice

import (
	"context"
	"sort"
	"sync"
	"time"
)

// Registry manages microservice and instance registration
type Registry struct {
	mu       sync.RWMutex
	services map[string]*Info

	// Instance storage (separate for efficient instance operations)
	instances        map[string]*Instance          // instanceID -> Instance
	instancesByService map[string]map[string]struct{} // serviceID -> instanceIDs

	// Indices for efficient lookups
	byName       map[string]map[string]struct{} // name -> serviceIDs
	byType       map[ServiceType]map[string]struct{}
	byCapability map[string]map[string]struct{}
	byTag        map[string]map[string]struct{}

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Callbacks
	onServiceRegistered   func(service *Info)
	onServiceDeregistered func(service *Info, reason string)
	onServiceUpdated      func(old, new *Info)
	onServiceStatusChange func(service *Info, oldStatus, newStatus Status)
	onInstanceRegistered  func(instance *Instance)
	onInstanceDeregistered func(instance *Instance, reason string)
	onInstanceStatusChange func(instance *Instance, oldStatus, newStatus Status)

	// Configuration
	cleanupInterval   time.Duration
	expirationEnabled bool
}

// RegistryConfig holds registry configuration
type RegistryConfig struct {
	CleanupInterval   time.Duration
	ExpirationEnabled bool
}

// DefaultRegistryConfig returns default registry configuration
func DefaultRegistryConfig() RegistryConfig {
	return RegistryConfig{
		CleanupInterval:   time.Minute,
		ExpirationEnabled: true,
	}
}

// NewRegistry creates a new microservice registry
func NewRegistry(cfg ...RegistryConfig) *Registry {
	var config RegistryConfig
	if len(cfg) > 0 {
		config = cfg[0]
	} else {
		config = DefaultRegistryConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	r := &Registry{
		services:           make(map[string]*Info),
		instances:          make(map[string]*Instance),
		instancesByService: make(map[string]map[string]struct{}),
		byName:             make(map[string]map[string]struct{}),
		byType:             make(map[ServiceType]map[string]struct{}),
		byCapability:       make(map[string]map[string]struct{}),
		byTag:              make(map[string]map[string]struct{}),
		ctx:                ctx,
		cancel:             cancel,
		cleanupInterval:    config.CleanupInterval,
		expirationEnabled:  config.ExpirationEnabled,
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
			case <-r.ctx.Done():
				return
			case <-ticker.C:
				r.cleanupExpired()
			}
		}
	}()
}

// cleanupExpired removes expired instances
func (r *Registry) cleanupExpired() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for id, inst := range r.instances {
		if !inst.ExpiresAt.IsZero() && now.After(inst.ExpiresAt) {
			_ = r.removeInstanceLocked(id, "expired")
		}
	}
}

// RegisterService registers a new service or updates an existing one
func (r *Registry) RegisterService(service *Info) error {
	if service == nil {
		return ErrInvalidService
	}
	if service.ID == "" {
		return ErrMissingServiceID
	}
	if service.Name == "" {
		return ErrMissingServiceName
	}

	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-r.ctx.Done():
		return ErrRegistryShutdown
	default:
	}

	existing, exists := r.services[service.ID]
	if exists {
		// Update existing
		oldService := existing.Clone()

		// Preserve registration time and instances
		service.RegisteredAt = existing.RegisteredAt
		service.UpdatedAt = now
		service.Instances = existing.Instances

		// Remove from old indices
		r.removeFromIndicesLocked(existing)

		// Store updated service
		r.services[service.ID] = service.Clone()

		// Add to indices
		r.addToIndicesLocked(service)

		// Callbacks
		if r.onServiceUpdated != nil {
			go r.onServiceUpdated(oldService, service)
		}
		if r.onServiceStatusChange != nil && oldService.Status != service.Status {
			go r.onServiceStatusChange(service, oldService.Status, service.Status)
		}
	} else {
		// New registration
		service.RegisteredAt = now
		service.UpdatedAt = now
		if service.Instances == nil {
			service.Instances = make([]*Instance, 0)
		}

		stored := service.Clone()
		r.services[service.ID] = stored
		r.instancesByService[service.ID] = make(map[string]struct{})

		// Add to indices
		r.addToIndicesLocked(service)

		if r.onServiceRegistered != nil {
			go r.onServiceRegistered(stored)
		}
	}

	return nil
}

// DeregisterService removes a service and all its instances
func (r *Registry) DeregisterService(serviceID string) error {
	return r.DeregisterServiceWithReason(serviceID, "")
}

// DeregisterServiceWithReason removes a service with a reason
func (r *Registry) DeregisterServiceWithReason(serviceID, reason string) error {
	if serviceID == "" {
		return ErrMissingServiceID
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	service, exists := r.services[serviceID]
	if !exists {
		return ErrServiceNotFound
	}

	// Remove all instances first
	if instanceIDs, ok := r.instancesByService[serviceID]; ok {
		for instID := range instanceIDs {
			if inst, ok := r.instances[instID]; ok {
				delete(r.instances, instID)
				if r.onInstanceDeregistered != nil {
					go r.onInstanceDeregistered(inst, "service deregistered")
				}
			}
		}
		delete(r.instancesByService, serviceID)
	}

	// Remove from indices
	r.removeFromIndicesLocked(service)

	// Remove service
	delete(r.services, serviceID)

	if r.onServiceDeregistered != nil {
		go r.onServiceDeregistered(service, reason)
	}

	return nil
}

// RegisterInstance registers a new instance for a service
func (r *Registry) RegisterInstance(instance *Instance) error {
	if instance == nil {
		return ErrInvalidInstance
	}
	if instance.ID == "" {
		return ErrMissingInstanceID
	}
	if instance.ServiceID == "" {
		return ErrMissingServiceID
	}
	if instance.Host == "" {
		return ErrMissingHost
	}

	now := time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	select {
	case <-r.ctx.Done():
		return ErrRegistryShutdown
	default:
	}

	// Verify service exists
	service, exists := r.services[instance.ServiceID]
	if !exists {
		return ErrServiceNotFound
	}

	existingInst, instExists := r.instances[instance.ID]
	if instExists {
		// Update existing instance
		oldStatus := existingInst.Status

		instance.RegisteredAt = existingInst.RegisteredAt
		instance.UpdatedAt = now

		stored := instance.Clone()
		r.instances[instance.ID] = stored

		// Update in service's instance list
		for i, inst := range service.Instances {
			if inst.ID == instance.ID {
				service.Instances[i] = stored
				break
			}
		}

		if r.onInstanceStatusChange != nil && oldStatus != instance.Status {
			go r.onInstanceStatusChange(stored, oldStatus, instance.Status)
		}
	} else {
		// New instance
		instance.RegisteredAt = now
		instance.UpdatedAt = now
		if instance.StartedAt.IsZero() {
			instance.StartedAt = now
		}

		stored := instance.Clone()
		r.instances[instance.ID] = stored

		// Add to service instance map
		if r.instancesByService[instance.ServiceID] == nil {
			r.instancesByService[instance.ServiceID] = make(map[string]struct{})
		}
		r.instancesByService[instance.ServiceID][instance.ID] = struct{}{}

		// Add to service's instance list
		service.Instances = append(service.Instances, stored)

		if r.onInstanceRegistered != nil {
			go r.onInstanceRegistered(stored)
		}
	}

	return nil
}

// RegisterInstanceWithTTL registers an instance with TTL
func (r *Registry) RegisterInstanceWithTTL(instance *Instance, ttl time.Duration) error {
	if ttl > 0 {
		instance.ExpiresAt = time.Now().Add(ttl)
	}
	return r.RegisterInstance(instance)
}

// DeregisterInstance removes an instance
func (r *Registry) DeregisterInstance(instanceID string) error {
	return r.DeregisterInstanceWithReason(instanceID, "")
}

// DeregisterInstanceWithReason removes an instance with a reason
func (r *Registry) DeregisterInstanceWithReason(instanceID, reason string) error {
	if instanceID == "" {
		return ErrMissingInstanceID
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	return r.removeInstanceLocked(instanceID, reason)
}

// removeInstanceLocked removes an instance (must be called with lock held)
func (r *Registry) removeInstanceLocked(instanceID, reason string) error {
	instance, exists := r.instances[instanceID]
	if !exists {
		return ErrInstanceNotFound
	}

	// Remove from service's instance map
	if svcInstances, ok := r.instancesByService[instance.ServiceID]; ok {
		delete(svcInstances, instanceID)
	}

	// Remove from service's instance list
	if service, ok := r.services[instance.ServiceID]; ok {
		for i, inst := range service.Instances {
			if inst.ID == instanceID {
				service.Instances = append(service.Instances[:i], service.Instances[i+1:]...)
				break
			}
		}
	}

	// Remove from main map
	delete(r.instances, instanceID)

	if r.onInstanceDeregistered != nil {
		go r.onInstanceDeregistered(instance, reason)
	}

	return nil
}

// GetService returns a service by ID
func (r *Registry) GetService(serviceID string) (*Info, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	service, exists := r.services[serviceID]
	if !exists {
		return nil, false
	}

	return service.Clone(), true
}

// GetServiceByName returns services by name
func (r *Registry) GetServiceByName(name string) []*Info {
	r.mu.RLock()
	defer r.mu.RUnlock()

	serviceIDs, exists := r.byName[name]
	if !exists {
		return nil
	}

	result := make([]*Info, 0, len(serviceIDs))
	for id := range serviceIDs {
		if svc, ok := r.services[id]; ok {
			result = append(result, svc.Clone())
		}
	}

	return result
}

// GetInstance returns an instance by ID
func (r *Registry) GetInstance(instanceID string) (*Instance, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	instance, exists := r.instances[instanceID]
	if !exists {
		return nil, false
	}

	// Skip expired
	if !instance.ExpiresAt.IsZero() && time.Now().After(instance.ExpiresAt) {
		return nil, false
	}

	return instance.Clone(), true
}

// GetInstances returns all instances for a service
func (r *Registry) GetInstances(serviceID string) []*Instance {
	r.mu.RLock()
	defer r.mu.RUnlock()

	instanceIDs, exists := r.instancesByService[serviceID]
	if !exists {
		return nil
	}

	now := time.Now()
	result := make([]*Instance, 0, len(instanceIDs))
	for id := range instanceIDs {
		if inst, ok := r.instances[id]; ok {
			if !inst.ExpiresAt.IsZero() && now.After(inst.ExpiresAt) {
				continue
			}
			result = append(result, inst.Clone())
		}
	}

	return result
}

// GetHealthyInstances returns healthy instances for a service
func (r *Registry) GetHealthyInstances(serviceID string) []*Instance {
	instances := r.GetInstances(serviceID)
	result := make([]*Instance, 0)
	for _, inst := range instances {
		if inst.Status == StatusRunning {
			result = append(result, inst)
		}
	}
	return result
}

// GetAvailableInstances returns available instances for a service
func (r *Registry) GetAvailableInstances(serviceID string) []*Instance {
	instances := r.GetInstances(serviceID)
	result := make([]*Instance, 0)
	for _, inst := range instances {
		if inst.Status.IsAvailable() {
			result = append(result, inst)
		}
	}
	return result
}

// QueryServices returns services matching the query
func (r *Registry) QueryServices(q *Query) []*Info {
	if q == nil {
		q = &Query{}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*Info, 0)
	for _, svc := range r.services {
		if q.Matches(svc) {
			result = append(result, svc.Clone())
		}
	}

	// Sort by name
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
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

// AllServices returns all registered services
func (r *Registry) AllServices() []*Info {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*Info, 0, len(r.services))
	for _, svc := range r.services {
		result = append(result, svc.Clone())
	}

	return result
}

// AllInstances returns all registered instances
func (r *Registry) AllInstances() []*Instance {
	r.mu.RLock()
	defer r.mu.RUnlock()

	now := time.Now()
	result := make([]*Instance, 0, len(r.instances))
	for _, inst := range r.instances {
		if !inst.ExpiresAt.IsZero() && now.After(inst.ExpiresAt) {
			continue
		}
		result = append(result, inst.Clone())
	}

	return result
}

// ServiceCount returns the number of registered services
func (r *Registry) ServiceCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.services)
}

// InstanceCount returns the number of registered instances
func (r *Registry) InstanceCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.instances)
}

// ServiceExists checks if a service exists
func (r *Registry) ServiceExists(serviceID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.services[serviceID]
	return exists
}

// InstanceExists checks if an instance exists
func (r *Registry) InstanceExists(instanceID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.instances[instanceID]
	return exists
}

// UpdateServiceStatus updates the status of a service
func (r *Registry) UpdateServiceStatus(serviceID string, status Status) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	service, exists := r.services[serviceID]
	if !exists {
		return ErrServiceNotFound
	}

	oldStatus := service.Status
	service.Status = status
	service.LastHealthCheck = time.Now()
	service.UpdatedAt = time.Now()

	if r.onServiceStatusChange != nil && oldStatus != status {
		go r.onServiceStatusChange(service, oldStatus, status)
	}

	return nil
}

// UpdateInstanceStatus updates the status of an instance
func (r *Registry) UpdateInstanceStatus(instanceID string, status Status) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	instance, exists := r.instances[instanceID]
	if !exists {
		return ErrInstanceNotFound
	}

	oldStatus := instance.Status
	instance.Status = status
	instance.LastHealthCheck = time.Now()
	instance.UpdatedAt = time.Now()

	if status == StatusRunning {
		instance.ConsecutiveFails = 0
	} else if status == StatusFailed {
		instance.ConsecutiveFails++
	}

	if r.onInstanceStatusChange != nil && oldStatus != status {
		go r.onInstanceStatusChange(instance, oldStatus, status)
	}

	return nil
}

// UpdateInstanceHealth updates instance health information
func (r *Registry) UpdateInstanceHealth(update *HealthUpdate) error {
	if update == nil || update.InstanceID == "" {
		return ErrInvalidInstance
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	instance, exists := r.instances[update.InstanceID]
	if !exists {
		return ErrInstanceNotFound
	}

	oldStatus := instance.Status
	instance.Status = update.Status
	instance.LastHealthCheck = update.CheckedAt
	instance.UpdatedAt = time.Now()

	if update.Status == StatusRunning {
		instance.ConsecutiveFails = 0
	} else if update.Status == StatusFailed {
		instance.ConsecutiveFails++
	}

	if r.onInstanceStatusChange != nil && oldStatus != update.Status {
		go r.onInstanceStatusChange(instance, oldStatus, update.Status)
	}

	return nil
}

// RenewInstance extends the TTL of an instance
func (r *Registry) RenewInstance(instanceID string, ttl time.Duration) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	instance, exists := r.instances[instanceID]
	if !exists {
		return ErrInstanceNotFound
	}

	instance.ExpiresAt = time.Now().Add(ttl)
	instance.UpdatedAt = time.Now()

	return nil
}

// Shutdown gracefully shuts down the registry
func (r *Registry) Shutdown(ctx context.Context) error {
	r.cancel()

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

// Clear removes all services and instances
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.services = make(map[string]*Info)
	r.instances = make(map[string]*Instance)
	r.instancesByService = make(map[string]map[string]struct{})
	r.byName = make(map[string]map[string]struct{})
	r.byType = make(map[ServiceType]map[string]struct{})
	r.byCapability = make(map[string]map[string]struct{})
	r.byTag = make(map[string]map[string]struct{})
}

// Callbacks

func (r *Registry) OnServiceRegistered(cb func(service *Info)) {
	r.onServiceRegistered = cb
}

func (r *Registry) OnServiceDeregistered(cb func(service *Info, reason string)) {
	r.onServiceDeregistered = cb
}

func (r *Registry) OnServiceUpdated(cb func(old, new *Info)) {
	r.onServiceUpdated = cb
}

func (r *Registry) OnServiceStatusChange(cb func(service *Info, oldStatus, newStatus Status)) {
	r.onServiceStatusChange = cb
}

func (r *Registry) OnInstanceRegistered(cb func(instance *Instance)) {
	r.onInstanceRegistered = cb
}

func (r *Registry) OnInstanceDeregistered(cb func(instance *Instance, reason string)) {
	r.onInstanceDeregistered = cb
}

func (r *Registry) OnInstanceStatusChange(cb func(instance *Instance, oldStatus, newStatus Status)) {
	r.onInstanceStatusChange = cb
}

// Index management helpers

func (r *Registry) addToIndicesLocked(svc *Info) {
	// Name index
	if svc.Name != "" {
		if r.byName[svc.Name] == nil {
			r.byName[svc.Name] = make(map[string]struct{})
		}
		r.byName[svc.Name][svc.ID] = struct{}{}
	}

	// Type index
	if svc.Type != "" {
		if r.byType[svc.Type] == nil {
			r.byType[svc.Type] = make(map[string]struct{})
		}
		r.byType[svc.Type][svc.ID] = struct{}{}
	}

	// Capability index
	for _, cap := range svc.Capabilities {
		if r.byCapability[cap] == nil {
			r.byCapability[cap] = make(map[string]struct{})
		}
		r.byCapability[cap][svc.ID] = struct{}{}
	}

	// Tag index
	for _, tag := range svc.Tags {
		if r.byTag[tag] == nil {
			r.byTag[tag] = make(map[string]struct{})
		}
		r.byTag[tag][svc.ID] = struct{}{}
	}
}

func (r *Registry) removeFromIndicesLocked(svc *Info) {
	// Name index
	if svc.Name != "" {
		if names, ok := r.byName[svc.Name]; ok {
			delete(names, svc.ID)
			if len(names) == 0 {
				delete(r.byName, svc.Name)
			}
		}
	}

	// Type index
	if svc.Type != "" {
		if types, ok := r.byType[svc.Type]; ok {
			delete(types, svc.ID)
			if len(types) == 0 {
				delete(r.byType, svc.Type)
			}
		}
	}

	// Capability index
	for _, cap := range svc.Capabilities {
		if caps, ok := r.byCapability[cap]; ok {
			delete(caps, svc.ID)
			if len(caps) == 0 {
				delete(r.byCapability, cap)
			}
		}
	}

	// Tag index
	for _, tag := range svc.Tags {
		if tags, ok := r.byTag[tag]; ok {
			delete(tags, svc.ID)
			if len(tags) == 0 {
				delete(r.byTag, tag)
			}
		}
	}
}
