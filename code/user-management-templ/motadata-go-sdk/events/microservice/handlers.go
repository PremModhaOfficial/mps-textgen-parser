package microservice

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

// Subjects defines the NATS subjects for microservice operations
type Subjects struct {
	// Service subjects
	RegisterService   string
	DeregisterService string
	QueryServices     string
	DiscoverServices  string

	// Instance subjects
	RegisterInstance   string
	DeregisterInstance string
	RenewInstance      string
	QueryInstances     string

	// Event subjects
	ServiceEvents  string
	InstanceEvents string

	// Health subjects
	HealthUpdate string
}

// DefaultSubjects returns default subject names
func DefaultSubjects() Subjects {
	return Subjects{
		RegisterService:    "services.register",
		DeregisterService:  "services.deregister",
		QueryServices:      "services.query",
		DiscoverServices:   "services.discover",
		RegisterInstance:   "instances.register",
		DeregisterInstance: "instances.deregister",
		RenewInstance:      "instances.renew",
		QueryInstances:     "instances.query",
		ServiceEvents:      "services.events",
		InstanceEvents:     "instances.events",
		HealthUpdate:       "health.update",
	}
}

// TenantSubjects returns subjects with tenant prefix
func TenantSubjects(tenantID string) Subjects {
	prefix := tenantID + "."
	return Subjects{
		RegisterService:    prefix + "services.register",
		DeregisterService:  prefix + "services.deregister",
		QueryServices:      prefix + "services.query",
		DiscoverServices:   prefix + "services.discover",
		RegisterInstance:   prefix + "instances.register",
		DeregisterInstance: prefix + "instances.deregister",
		RenewInstance:      prefix + "instances.renew",
		QueryInstances:     prefix + "instances.query",
		ServiceEvents:      prefix + "services.events",
		InstanceEvents:     prefix + "instances.events",
		HealthUpdate:       prefix + "health.update",
	}
}

// Handler handles microservice-related message operations
type Handler struct {
	manager   *Manager
	publisher core.Publisher
	subjects  Subjects
}

// HandlerOption is a functional option for configuring Handler
type HandlerOption func(*Handler)

// WithPublisher sets the publisher for sending responses
func WithPublisher(publisher core.Publisher) HandlerOption {
	return func(h *Handler) {
		h.publisher = publisher
	}
}

// WithSubjects sets custom subjects
func WithSubjects(subjects Subjects) HandlerOption {
	return func(h *Handler) {
		h.subjects = subjects
	}
}

// NewHandler creates a new microservice handler
func NewHandler(manager *Manager, opts ...HandlerOption) *Handler {
	h := &Handler{
		manager:  manager,
		subjects: DefaultSubjects(),
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

// NewHandlerWithPublisher creates a new microservice handler with a publisher for responses
// Deprecated: Use NewHandler with WithPublisher option instead
func NewHandlerWithPublisher(manager *Manager, publisher core.Publisher, subjects ...Subjects) *Handler {
	var s Subjects
	if len(subjects) > 0 {
		s = subjects[0]
	} else {
		s = DefaultSubjects()
	}

	return &Handler{
		manager:   manager,
		publisher: publisher,
		subjects:  s,
	}
}

// SetPublisher sets the publisher for sending responses
func (h *Handler) SetPublisher(publisher core.Publisher) {
	h.publisher = publisher
}

// RegisterSubscriptions sets up all subscription handlers
func (h *Handler) RegisterSubscriptions(ctx context.Context, subscriber core.Subscriber) ([]core.Subscription, error) {
	subscriptions := make([]core.Subscription, 0, 8)

	// Service handlers
	sub, err := subscriber.Subscribe(ctx, h.subjects.RegisterService, h.HandleRegisterService)
	if err != nil {
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	sub, err = subscriber.Subscribe(ctx, h.subjects.DeregisterService, h.HandleDeregisterService)
	if err != nil {
		h.unsubscribeAll(subscriptions)
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	sub, err = subscriber.Subscribe(ctx, h.subjects.QueryServices, h.HandleQueryServices)
	if err != nil {
		h.unsubscribeAll(subscriptions)
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	sub, err = subscriber.Subscribe(ctx, h.subjects.DiscoverServices, h.HandleDiscoverServices)
	if err != nil {
		h.unsubscribeAll(subscriptions)
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	// Instance handlers
	sub, err = subscriber.Subscribe(ctx, h.subjects.RegisterInstance, h.HandleRegisterInstance)
	if err != nil {
		h.unsubscribeAll(subscriptions)
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	sub, err = subscriber.Subscribe(ctx, h.subjects.DeregisterInstance, h.HandleDeregisterInstance)
	if err != nil {
		h.unsubscribeAll(subscriptions)
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	sub, err = subscriber.Subscribe(ctx, h.subjects.RenewInstance, h.HandleRenewInstance)
	if err != nil {
		h.unsubscribeAll(subscriptions)
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	sub, err = subscriber.Subscribe(ctx, h.subjects.QueryInstances, h.HandleQueryInstances)
	if err != nil {
		h.unsubscribeAll(subscriptions)
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	// Health handler
	sub, err = subscriber.Subscribe(ctx, h.subjects.HealthUpdate, h.HandleHealthUpdate)
	if err != nil {
		h.unsubscribeAll(subscriptions)
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	return subscriptions, nil
}

// RegisterQueueSubscriptions sets up queue subscription handlers for load balancing
func (h *Handler) RegisterQueueSubscriptions(ctx context.Context, subscriber core.Subscriber, queue string) ([]core.Subscription, error) {
	subscriptions := make([]core.Subscription, 0, 8)

	// Service handlers
	sub, err := subscriber.QueueSubscribe(ctx, h.subjects.RegisterService, queue, h.HandleRegisterService)
	if err != nil {
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	sub, err = subscriber.QueueSubscribe(ctx, h.subjects.DeregisterService, queue, h.HandleDeregisterService)
	if err != nil {
		h.unsubscribeAll(subscriptions)
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	sub, err = subscriber.QueueSubscribe(ctx, h.subjects.QueryServices, queue, h.HandleQueryServices)
	if err != nil {
		h.unsubscribeAll(subscriptions)
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	sub, err = subscriber.QueueSubscribe(ctx, h.subjects.DiscoverServices, queue, h.HandleDiscoverServices)
	if err != nil {
		h.unsubscribeAll(subscriptions)
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	// Instance handlers
	sub, err = subscriber.QueueSubscribe(ctx, h.subjects.RegisterInstance, queue, h.HandleRegisterInstance)
	if err != nil {
		h.unsubscribeAll(subscriptions)
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	sub, err = subscriber.QueueSubscribe(ctx, h.subjects.DeregisterInstance, queue, h.HandleDeregisterInstance)
	if err != nil {
		h.unsubscribeAll(subscriptions)
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	sub, err = subscriber.QueueSubscribe(ctx, h.subjects.RenewInstance, queue, h.HandleRenewInstance)
	if err != nil {
		h.unsubscribeAll(subscriptions)
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	sub, err = subscriber.QueueSubscribe(ctx, h.subjects.QueryInstances, queue, h.HandleQueryInstances)
	if err != nil {
		h.unsubscribeAll(subscriptions)
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	return subscriptions, nil
}

func (h *Handler) unsubscribeAll(subs []core.Subscription) {
	for _, sub := range subs {
		if sub != nil {
			_ = sub.Unsubscribe()
		}
	}
}

// HandleRegisterService handles service registration messages
func (h *Handler) HandleRegisterService(ctx context.Context, msg *nats.Msg) error {
	var req RegistrationRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	// Extract tenant from context
	if tenantID, ok := core.TenantIDFromContext(ctx); ok && tenantID != "" && req.Service != nil {
		req.Service.TenantID = tenantID
	}

	if err := h.manager.RegisterService(ctx, &req); err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	return h.sendResponse(ctx, msg, &ServiceRegistrationResponse{
		Success:   true,
		ServiceID: req.Service.ID,
		Message:   "registered",
	})
}

// HandleDeregisterService handles service deregistration messages
func (h *Handler) HandleDeregisterService(ctx context.Context, msg *nats.Msg) error {
	var req DeregistrationRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	if err := h.manager.DeregisterService(ctx, &req); err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	return h.sendResponse(ctx, msg, &ServiceDeregistrationResponse{
		Success:   true,
		ServiceID: req.ServiceID,
		Message:   "deregistered",
	})
}

// HandleQueryServices handles service query messages
func (h *Handler) HandleQueryServices(ctx context.Context, msg *nats.Msg) error {
	var q Query
	if err := json.Unmarshal(msg.Data, &q); err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	// Apply tenant filter from context
	if tenantID, ok := core.TenantIDFromContext(ctx); ok && tenantID != "" && q.TenantID == "" {
		q.TenantID = tenantID
	}

	services := h.manager.QueryServices(&q)

	return h.sendResponse(ctx, msg, &QueryServicesResponse{
		Services: services,
		Count:    len(services),
	})
}

// HandleDiscoverServices handles service discovery messages
func (h *Handler) HandleDiscoverServices(ctx context.Context, msg *nats.Msg) error {
	var req DiscoverRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	q := &Query{
		OnlyAvailable: true,
		HasInstances:  true,
	}

	if req.Name != "" {
		q.Names = []string{req.Name}
	}
	if req.Type != "" {
		q.Types = []ServiceType{req.Type}
	}
	if req.Version != "" {
		q.Version = req.Version
	}
	if len(req.Capabilities) > 0 {
		q.Capabilities = req.Capabilities
	}
	if len(req.Tags) > 0 {
		q.Tags = req.Tags
	}
	if req.OnlyHealthy {
		q.OnlyHealthy = true
		q.OnlyAvailable = false
	}

	// Apply tenant filter
	if tenantID, ok := core.TenantIDFromContext(ctx); ok && tenantID != "" {
		q.TenantID = tenantID
	} else if req.TenantID != "" {
		q.TenantID = req.TenantID
	}

	services := h.manager.QueryServices(q)

	return h.sendResponse(ctx, msg, &DiscoverResponse{
		Services: services,
		Count:    len(services),
	})
}

// HandleRegisterInstance handles instance registration messages
func (h *Handler) HandleRegisterInstance(ctx context.Context, msg *nats.Msg) error {
	var req InstanceRegistrationRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	if err := h.manager.RegisterInstance(ctx, &req); err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	return h.sendResponse(ctx, msg, &InstanceRegistrationResponse{
		Success:    true,
		InstanceID: req.Instance.ID,
		ServiceID:  req.Instance.ServiceID,
		Message:    "registered",
	})
}

// HandleDeregisterInstance handles instance deregistration messages
func (h *Handler) HandleDeregisterInstance(ctx context.Context, msg *nats.Msg) error {
	var req InstanceDeregistrationRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	if err := h.manager.DeregisterInstance(ctx, req.InstanceID); err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	return h.sendResponse(ctx, msg, &InstanceDeregistrationResponse{
		Success:    true,
		InstanceID: req.InstanceID,
		Message:    "deregistered",
	})
}

// HandleRenewInstance handles instance TTL renewal messages
func (h *Handler) HandleRenewInstance(ctx context.Context, msg *nats.Msg) error {
	var req RenewInstanceRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	if err := h.manager.RenewInstance(ctx, req.InstanceID, req.TTL); err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	return h.sendResponse(ctx, msg, &RenewInstanceResponse{
		Success:    true,
		InstanceID: req.InstanceID,
		Message:    "renewed",
		ExpiresAt:  time.Now().Add(req.TTL),
	})
}

// HandleQueryInstances handles instance query messages
func (h *Handler) HandleQueryInstances(ctx context.Context, msg *nats.Msg) error {
	var req QueryInstancesRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	var instances []*Instance
	if req.OnlyHealthy {
		instances = h.manager.GetHealthyInstances(req.ServiceID)
	} else if req.OnlyAvailable {
		instances = h.manager.GetAvailableInstances(req.ServiceID)
	} else {
		instances = h.manager.GetInstances(req.ServiceID)
	}

	return h.sendResponse(ctx, msg, &QueryInstancesResponse{
		Instances: instances,
		Count:     len(instances),
		ServiceID: req.ServiceID,
	})
}

// HandleHealthUpdate handles health status update messages
func (h *Handler) HandleHealthUpdate(ctx context.Context, msg *nats.Msg) error {
	var update HealthUpdate
	if err := json.Unmarshal(msg.Data, &update); err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	if update.CheckedAt.IsZero() {
		update.CheckedAt = time.Now()
	}

	var err error
	if update.InstanceID != "" {
		err = h.manager.Registry().UpdateInstanceHealth(&update)
	} else if update.ServiceID != "" {
		err = h.manager.UpdateServiceStatus(ctx, update.ServiceID, update.Status)
	} else {
		return h.sendErrorResponse(ctx, msg, ErrMissingServiceID)
	}

	if err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	return h.sendResponse(ctx, msg, &HealthUpdateResponse{
		Success:    true,
		ServiceID:  update.ServiceID,
		InstanceID: update.InstanceID,
		Status:     update.Status.String(),
	})
}

// Helper methods

func (h *Handler) sendResponse(ctx context.Context, msg *nats.Msg, response any) error {
	if msg.Reply == "" {
		return nil // No reply expected
	}

	if h.publisher == nil {
		return nil // No publisher configured, can't send response
	}

	data, err := json.Marshal(response)
	if err != nil {
		return err
	}

	replyMsg := &nats.Msg{Data: data, Header: make(nats.Header)}
	replyMsg.Header.Set("status", "success")

	return h.publisher.Publish(ctx, msg.Reply, replyMsg)
}

func (h *Handler) sendErrorResponse(ctx context.Context, msg *nats.Msg, err error) error {
	if msg.Reply == "" {
		return err // No reply expected, return the error
	}

	if h.publisher == nil {
		return err // No publisher configured, return the error
	}

	response := &ErrorResponse{
		Success: false,
		Error:   err.Error(),
	}

	data, marshalErr := json.Marshal(response)
	if marshalErr != nil {
		return marshalErr
	}

	replyMsg := &nats.Msg{Data: data, Header: make(nats.Header)}
	replyMsg.Header.Set("status", "error")

	if pubErr := h.publisher.Publish(ctx, msg.Reply, replyMsg); pubErr != nil {
		return pubErr
	}
	return err
}

// Request/Response types

// ServiceRegistrationResponse is the response to a service registration
type ServiceRegistrationResponse struct {
	Success   bool   `json:"success"`
	ServiceID string `json:"service_id"`
	Message   string `json:"message"`
}

// ServiceDeregistrationResponse is the response to a service deregistration
type ServiceDeregistrationResponse struct {
	Success   bool   `json:"success"`
	ServiceID string `json:"service_id"`
	Message   string `json:"message"`
}

// QueryServicesResponse is the response to a service query
type QueryServicesResponse struct {
	Services []*Info `json:"services"`
	Count    int     `json:"count"`
}

// DiscoverRequest is a request to discover services
type DiscoverRequest struct {
	Name         string      `json:"name,omitempty"`
	Type         ServiceType `json:"type,omitempty"`
	Version      string      `json:"version,omitempty"`
	Capabilities []string    `json:"capabilities,omitempty"`
	Tags         []string    `json:"tags,omitempty"`
	OnlyHealthy  bool        `json:"only_healthy,omitempty"`
	TenantID     string      `json:"tenant_id,omitempty"`
}

// DiscoverResponse is the response to a discover request
type DiscoverResponse struct {
	Services []*Info `json:"services"`
	Count    int     `json:"count"`
}

// InstanceRegistrationResponse is the response to an instance registration
type InstanceRegistrationResponse struct {
	Success    bool   `json:"success"`
	InstanceID string `json:"instance_id"`
	ServiceID  string `json:"service_id"`
	Message    string `json:"message"`
}

// InstanceDeregistrationRequest is a request to deregister an instance
type InstanceDeregistrationRequest struct {
	InstanceID string `json:"instance_id"`
	ServiceID  string `json:"service_id,omitempty"`
	Reason     string `json:"reason,omitempty"`
}

// InstanceDeregistrationResponse is the response to an instance deregistration
type InstanceDeregistrationResponse struct {
	Success    bool   `json:"success"`
	InstanceID string `json:"instance_id"`
	Message    string `json:"message"`
}

// RenewInstanceRequest is a request to renew instance TTL
type RenewInstanceRequest struct {
	InstanceID string        `json:"instance_id"`
	TTL        time.Duration `json:"ttl"`
}

// RenewInstanceResponse is the response to a renew request
type RenewInstanceResponse struct {
	Success    bool      `json:"success"`
	InstanceID string    `json:"instance_id"`
	Message    string    `json:"message"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// QueryInstancesRequest is a request to query instances
type QueryInstancesRequest struct {
	ServiceID     string `json:"service_id"`
	OnlyHealthy   bool   `json:"only_healthy,omitempty"`
	OnlyAvailable bool   `json:"only_available,omitempty"`
}

// QueryInstancesResponse is the response to an instance query
type QueryInstancesResponse struct {
	Instances []*Instance `json:"instances"`
	Count     int         `json:"count"`
	ServiceID string      `json:"service_id"`
}

// HealthUpdateResponse is the response to a health update
type HealthUpdateResponse struct {
	Success    bool   `json:"success"`
	ServiceID  string `json:"service_id,omitempty"`
	InstanceID string `json:"instance_id,omitempty"`
	Status     string `json:"status"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// Client provides a client interface for microservice operations over messaging
type Client struct {
	publisher core.Publisher
	subjects  Subjects
	timeout   time.Duration
}

// NewClient creates a new microservice client
func NewClient(publisher core.Publisher, subjects ...Subjects) *Client {
	var s Subjects
	if len(subjects) > 0 {
		s = subjects[0]
	} else {
		s = DefaultSubjects()
	}

	return &Client{
		publisher: publisher,
		subjects:  s,
		timeout:   10 * time.Second,
	}
}

// SetTimeout sets the request timeout
func (c *Client) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

// RegisterService registers a service
func (c *Client) RegisterService(ctx context.Context, service *Info) (*ServiceRegistrationResponse, error) {
	req := &RegistrationRequest{
		Service: service,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	msg := &nats.Msg{Data: data, Header: make(nats.Header)}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	reply, err := c.publisher.Request(ctx, c.subjects.RegisterService, msg)
	if err != nil {
		return nil, err
	}

	var resp ServiceRegistrationResponse
	if err := json.Unmarshal(reply.Data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// DeregisterService deregisters a service
func (c *Client) DeregisterService(ctx context.Context, serviceID, reason string) (*ServiceDeregistrationResponse, error) {
	req := &DeregistrationRequest{
		ServiceID: serviceID,
		Reason:    reason,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	msg := &nats.Msg{Data: data, Header: make(nats.Header)}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	reply, err := c.publisher.Request(ctx, c.subjects.DeregisterService, msg)
	if err != nil {
		return nil, err
	}

	var resp ServiceDeregistrationResponse
	if err := json.Unmarshal(reply.Data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// QueryServices queries for services
func (c *Client) QueryServices(ctx context.Context, q *Query) (*QueryServicesResponse, error) {
	data, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}

	msg := &nats.Msg{Data: data, Header: make(nats.Header)}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	reply, err := c.publisher.Request(ctx, c.subjects.QueryServices, msg)
	if err != nil {
		return nil, err
	}

	var resp QueryServicesResponse
	if err := json.Unmarshal(reply.Data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Discover discovers services
func (c *Client) Discover(ctx context.Context, name string, opts ...DiscoverOption) (*DiscoverResponse, error) {
	q := &Query{
		Names:         []string{name},
		OnlyAvailable: true,
		HasInstances:  true,
	}

	for _, opt := range opts {
		opt(q)
	}

	req := &DiscoverRequest{
		Name:         name,
		Version:      q.Version,
		Capabilities: q.Capabilities,
		Tags:         q.Tags,
		OnlyHealthy:  q.OnlyHealthy,
		TenantID:     q.TenantID,
	}

	if len(q.Types) > 0 {
		req.Type = q.Types[0]
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	msg := &nats.Msg{Data: data, Header: make(nats.Header)}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	reply, err := c.publisher.Request(ctx, c.subjects.DiscoverServices, msg)
	if err != nil {
		return nil, err
	}

	var resp DiscoverResponse
	if err := json.Unmarshal(reply.Data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// RegisterInstance registers an instance
func (c *Client) RegisterInstance(ctx context.Context, instance *Instance, ttl time.Duration) (*InstanceRegistrationResponse, error) {
	req := &InstanceRegistrationRequest{
		Instance: instance,
		TTL:      ttl,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	msg := &nats.Msg{Data: data, Header: make(nats.Header)}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	reply, err := c.publisher.Request(ctx, c.subjects.RegisterInstance, msg)
	if err != nil {
		return nil, err
	}

	var resp InstanceRegistrationResponse
	if err := json.Unmarshal(reply.Data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// DeregisterInstance deregisters an instance
func (c *Client) DeregisterInstance(ctx context.Context, instanceID, reason string) (*InstanceDeregistrationResponse, error) {
	req := &InstanceDeregistrationRequest{
		InstanceID: instanceID,
		Reason:     reason,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	msg := &nats.Msg{Data: data, Header: make(nats.Header)}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	reply, err := c.publisher.Request(ctx, c.subjects.DeregisterInstance, msg)
	if err != nil {
		return nil, err
	}

	var resp InstanceDeregistrationResponse
	if err := json.Unmarshal(reply.Data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// RenewInstance renews an instance's TTL
func (c *Client) RenewInstance(ctx context.Context, instanceID string, ttl time.Duration) (*RenewInstanceResponse, error) {
	req := &RenewInstanceRequest{
		InstanceID: instanceID,
		TTL:        ttl,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	msg := &nats.Msg{Data: data, Header: make(nats.Header)}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	reply, err := c.publisher.Request(ctx, c.subjects.RenewInstance, msg)
	if err != nil {
		return nil, err
	}

	var resp RenewInstanceResponse
	if err := json.Unmarshal(reply.Data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// QueryInstances queries for instances
func (c *Client) QueryInstances(ctx context.Context, serviceID string, onlyHealthy bool) (*QueryInstancesResponse, error) {
	req := &QueryInstancesRequest{
		ServiceID:   serviceID,
		OnlyHealthy: onlyHealthy,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	msg := &nats.Msg{Data: data, Header: make(nats.Header)}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	reply, err := c.publisher.Request(ctx, c.subjects.QueryInstances, msg)
	if err != nil {
		return nil, err
	}

	var resp QueryInstancesResponse
	if err := json.Unmarshal(reply.Data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// UpdateHealth sends a health update
func (c *Client) UpdateHealth(ctx context.Context, update *HealthUpdate) error {
	data, err := json.Marshal(update)
	if err != nil {
		return err
	}

	msg := &nats.Msg{Data: data, Header: make(nats.Header)}
	return c.publisher.Publish(ctx, c.subjects.HealthUpdate, msg)
}
