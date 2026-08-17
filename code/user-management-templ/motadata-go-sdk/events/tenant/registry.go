package tenant

import (
	"sync"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"
)

// Info holds tenant information
type Info struct {
	ID       string
	Name     string
	Metadata map[string]string
	Active   bool
}

// Registry manages tenant registration
type Registry struct {
	mu      sync.RWMutex
	tenants map[string]*Info

	// Callbacks
	onRegistered   func(tenantID string)
	onUnregistered func(tenantID string)
}

// NewRegistry creates a new tenant registry
func NewRegistry() *Registry {
	return &Registry{
		tenants: make(map[string]*Info),
	}
}

// Register registers a new tenant
func (r *Registry) Register(info *Info) error {
	if info == nil || info.ID == "" {
		return utils.ErrTenantNotFound
	}

	r.mu.Lock()
	if _, exists := r.tenants[info.ID]; exists {
		r.mu.Unlock()
		return utils.ErrTenantExists
	}
	r.tenants[info.ID] = info
	r.mu.Unlock()

	if r.onRegistered != nil {
		r.onRegistered(info.ID)
	}

	return nil
}

// Unregister removes a tenant
func (r *Registry) Unregister(tenantID string) error {
	r.mu.Lock()
	if _, exists := r.tenants[tenantID]; !exists {
		r.mu.Unlock()
		return utils.ErrTenantNotFound
	}
	delete(r.tenants, tenantID)
	r.mu.Unlock()

	if r.onUnregistered != nil {
		r.onUnregistered(tenantID)
	}

	return nil
}

// Get returns tenant info
func (r *Registry) Get(tenantID string) (*Info, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	info, ok := r.tenants[tenantID]
	return info, ok
}

// Exists checks if a tenant exists
func (r *Registry) Exists(tenantID string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.tenants[tenantID]
	return ok
}

// List returns all tenant IDs
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]string, 0, len(r.tenants))
	for id := range r.tenants {
		ids = append(ids, id)
	}
	return ids
}

// All returns all tenant info
func (r *Registry) All() []*Info {
	r.mu.RLock()
	defer r.mu.RUnlock()

	infos := make([]*Info, 0, len(r.tenants))
	for _, info := range r.tenants {
		infos = append(infos, info)
	}
	return infos
}

// Count returns the number of registered tenants
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.tenants)
}

// Update updates tenant info
func (r *Registry) Update(tenantID string, fn func(*Info)) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	info, ok := r.tenants[tenantID]
	if !ok {
		return utils.ErrTenantNotFound
	}

	fn(info)
	return nil
}

// SetActive sets tenant active status
func (r *Registry) SetActive(tenantID string, active bool) error {
	return r.Update(tenantID, func(info *Info) {
		info.Active = active
	})
}

// OnRegistered sets callback for registration
func (r *Registry) OnRegistered(cb func(tenantID string)) {
	r.onRegistered = cb
}

// OnUnregistered sets callback for unregistration
func (r *Registry) OnUnregistered(cb func(tenantID string)) {
	r.onUnregistered = cb
}

// ActiveTenants returns IDs of active tenants
func (r *Registry) ActiveTenants() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]string, 0)
	for id, info := range r.tenants {
		if info.Active {
			ids = append(ids, id)
		}
	}
	return ids
}
