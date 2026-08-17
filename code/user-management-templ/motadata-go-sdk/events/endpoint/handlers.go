package endpoint

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nats-io/nats.go"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

// Subjects defines the NATS subjects for endpoint operations
type Subjects struct {
	// Registration subjects
	Register   string
	Deregister string
	Renew      string

	// Query subjects
	Query    string
	Discover string

	// Event subjects (for publishing)
	Events string

	// Health subjects
	HealthUpdate string
}

// DefaultSubjects returns default subject names
func DefaultSubjects() Subjects {
	return Subjects{
		Register:     "endpoints.register",
		Deregister:   "endpoints.deregister",
		Renew:        "endpoints.renew",
		Query:        "endpoints.query",
		Discover:     "endpoints.discover",
		Events:       "endpoints.events",
		HealthUpdate: "endpoints.health.update",
	}
}

// TenantSubjects returns subjects with tenant prefix
func TenantSubjects(tenantID string) Subjects {
	prefix := tenantID + "."
	return Subjects{
		Register:     prefix + "endpoints.register",
		Deregister:   prefix + "endpoints.deregister",
		Renew:        prefix + "endpoints.renew",
		Query:        prefix + "endpoints.query",
		Discover:     prefix + "endpoints.discover",
		Events:       prefix + "endpoints.events",
		HealthUpdate: prefix + "endpoints.health.update",
	}
}

// Handler handles endpoint-related message operations
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

// NewHandler creates a new endpoint handler
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

// NewHandlerWithPublisher creates a new endpoint handler with a publisher for responses
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
	subscriptions := make([]core.Subscription, 0, 5)

	// Registration handler
	sub, err := subscriber.Subscribe(ctx, h.subjects.Register, h.HandleRegister)
	if err != nil {
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	// Deregistration handler
	sub, err = subscriber.Subscribe(ctx, h.subjects.Deregister, h.HandleDeregister)
	if err != nil {
		h.unsubscribeAll(subscriptions)
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	// Renew handler
	sub, err = subscriber.Subscribe(ctx, h.subjects.Renew, h.HandleRenew)
	if err != nil {
		h.unsubscribeAll(subscriptions)
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	// Query handler
	sub, err = subscriber.Subscribe(ctx, h.subjects.Query, h.HandleQuery)
	if err != nil {
		h.unsubscribeAll(subscriptions)
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	// Discover handler
	sub, err = subscriber.Subscribe(ctx, h.subjects.Discover, h.HandleDiscover)
	if err != nil {
		h.unsubscribeAll(subscriptions)
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	// Health update handler
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
	subscriptions := make([]core.Subscription, 0, 5)

	// Registration handler
	sub, err := subscriber.QueueSubscribe(ctx, h.subjects.Register, queue, h.HandleRegister)
	if err != nil {
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	// Deregistration handler
	sub, err = subscriber.QueueSubscribe(ctx, h.subjects.Deregister, queue, h.HandleDeregister)
	if err != nil {
		h.unsubscribeAll(subscriptions)
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	// Renew handler
	sub, err = subscriber.QueueSubscribe(ctx, h.subjects.Renew, queue, h.HandleRenew)
	if err != nil {
		h.unsubscribeAll(subscriptions)
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	// Query handler
	sub, err = subscriber.QueueSubscribe(ctx, h.subjects.Query, queue, h.HandleQuery)
	if err != nil {
		h.unsubscribeAll(subscriptions)
		return nil, err
	}
	subscriptions = append(subscriptions, sub)

	// Discover handler
	sub, err = subscriber.QueueSubscribe(ctx, h.subjects.Discover, queue, h.HandleDiscover)
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

// HandleRegister handles endpoint registration messages
func (h *Handler) HandleRegister(ctx context.Context, msg *nats.Msg) error {
	var req RegistrationRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	// Extract tenant from context if available
	if tenantID, ok := core.TenantIDFromContext(ctx); ok && tenantID != "" && req.Endpoint != nil {
		req.Endpoint.TenantID = tenantID
	}

	if err := h.manager.Register(ctx, &req); err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	return h.sendResponse(ctx, msg, &RegistrationResponse{
		Success:    true,
		EndpointID: req.Endpoint.ID,
		Message:    "registered",
	})
}

// HandleDeregister handles endpoint deregistration messages
func (h *Handler) HandleDeregister(ctx context.Context, msg *nats.Msg) error {
	var req DeregistrationRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	if err := h.manager.Deregister(ctx, &req); err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	return h.sendResponse(ctx, msg, &DeregistrationResponse{
		Success:    true,
		EndpointID: req.EndpointID,
		Message:    "deregistered",
	})
}

// HandleRenew handles TTL renewal messages
func (h *Handler) HandleRenew(ctx context.Context, msg *nats.Msg) error {
	var req RenewRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	if err := h.manager.Renew(ctx, req.EndpointID, req.TTL); err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	return h.sendResponse(ctx, msg, &RenewResponse{
		Success:    true,
		EndpointID: req.EndpointID,
		Message:    "renewed",
		ExpiresAt:  time.Now().Add(req.TTL),
	})
}

// HandleQuery handles endpoint query messages
func (h *Handler) HandleQuery(ctx context.Context, msg *nats.Msg) error {
	var q Query
	if err := json.Unmarshal(msg.Data, &q); err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	// Apply tenant filter from context
	if tenantID, ok := core.TenantIDFromContext(ctx); ok && tenantID != "" && q.TenantID == "" {
		q.TenantID = tenantID
	}

	endpoints := h.manager.Query(&q)

	return h.sendResponse(ctx, msg, &QueryResponse{
		Endpoints: endpoints,
		Count:     len(endpoints),
	})
}

// HandleDiscover handles service discovery messages
func (h *Handler) HandleDiscover(ctx context.Context, msg *nats.Msg) error {
	var req DiscoverRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	// Build query from discover request
	q := &Query{
		OnlyAvailable: true,
	}

	if req.ServiceID != "" {
		q.ServiceIDs = []string{req.ServiceID}
	}
	if req.Version != "" {
		q.Version = req.Version
	}
	if req.Protocol != "" {
		q.Protocol = req.Protocol
	}
	if len(req.Tags) > 0 {
		q.Tags = req.Tags
	}
	if req.OnlyHealthy {
		q.OnlyHealthy = true
		q.OnlyAvailable = false
	}

	// Apply tenant filter from context
	if tenantID, ok := core.TenantIDFromContext(ctx); ok && tenantID != "" {
		q.TenantID = tenantID
	} else if req.TenantID != "" {
		q.TenantID = req.TenantID
	}

	endpoints := h.manager.Query(q)

	return h.sendResponse(ctx, msg, &DiscoverResponse{
		Endpoints: endpoints,
		Count:     len(endpoints),
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

	if err := h.manager.Registry().UpdateHealth(&update); err != nil {
		return h.sendErrorResponse(ctx, msg, err)
	}

	return h.sendResponse(ctx, msg, &HealthUpdateResponse{
		Success:    true,
		EndpointID: update.EndpointID,
		Status:     update.Status.String(),
	})
}

// Helper methods for sending responses

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

// Request/Response types for handlers

// RegistrationResponse is the response to a registration request
type RegistrationResponse struct {
	Success    bool   `json:"success"`
	EndpointID string `json:"endpoint_id"`
	Message    string `json:"message"`
}

// DeregistrationResponse is the response to a deregistration request
type DeregistrationResponse struct {
	Success    bool   `json:"success"`
	EndpointID string `json:"endpoint_id"`
	Message    string `json:"message"`
}

// RenewRequest is a request to renew endpoint TTL
type RenewRequest struct {
	EndpointID string        `json:"endpoint_id"`
	TTL        time.Duration `json:"ttl"`
}

// RenewResponse is the response to a renew request
type RenewResponse struct {
	Success    bool      `json:"success"`
	EndpointID string    `json:"endpoint_id"`
	Message    string    `json:"message"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// QueryResponse is the response to a query request
type QueryResponse struct {
	Endpoints []*Info `json:"endpoints"`
	Count     int     `json:"count"`
}

// DiscoverRequest is a request to discover endpoints
type DiscoverRequest struct {
	ServiceID   string   `json:"service_id"`
	Version     string   `json:"version,omitempty"`
	Protocol    Protocol `json:"protocol,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	OnlyHealthy bool     `json:"only_healthy,omitempty"`
	TenantID    string   `json:"tenant_id,omitempty"`
}

// DiscoverResponse is the response to a discover request
type DiscoverResponse struct {
	Endpoints []*Info `json:"endpoints"`
	Count     int     `json:"count"`
	ServiceID string  `json:"service_id"`
}

// HealthUpdateResponse is the response to a health update
type HealthUpdateResponse struct {
	Success    bool   `json:"success"`
	EndpointID string `json:"endpoint_id"`
	Status     string `json:"status"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// Client provides a client interface for endpoint operations over messaging
type Client struct {
	publisher core.Publisher
	subjects  Subjects
	timeout   time.Duration
}

// NewClient creates a new endpoint client
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

// Register registers an endpoint
func (c *Client) Register(ctx context.Context, endpoint *Info, ttl time.Duration) (*RegistrationResponse, error) {
	req := &RegistrationRequest{
		Endpoint: endpoint,
		TTL:      ttl,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	msg := &nats.Msg{Data: data, Header: make(nats.Header)}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	reply, err := c.publisher.Request(ctx, c.subjects.Register, msg)
	if err != nil {
		return nil, err
	}

	var resp RegistrationResponse
	if err := json.Unmarshal(reply.Data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Deregister deregisters an endpoint
func (c *Client) Deregister(ctx context.Context, endpointID, reason string) (*DeregistrationResponse, error) {
	req := &DeregistrationRequest{
		EndpointID: endpointID,
		Reason:     reason,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	msg := &nats.Msg{Data: data, Header: make(nats.Header)}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	reply, err := c.publisher.Request(ctx, c.subjects.Deregister, msg)
	if err != nil {
		return nil, err
	}

	var resp DeregistrationResponse
	if err := json.Unmarshal(reply.Data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Renew renews an endpoint's TTL
func (c *Client) Renew(ctx context.Context, endpointID string, ttl time.Duration) (*RenewResponse, error) {
	req := &RenewRequest{
		EndpointID: endpointID,
		TTL:        ttl,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	msg := &nats.Msg{Data: data, Header: make(nats.Header)}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	reply, err := c.publisher.Request(ctx, c.subjects.Renew, msg)
	if err != nil {
		return nil, err
	}

	var resp RenewResponse
	if err := json.Unmarshal(reply.Data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Query queries for endpoints
func (c *Client) Query(ctx context.Context, q *Query) (*QueryResponse, error) {
	data, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}

	msg := &nats.Msg{Data: data, Header: make(nats.Header)}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	reply, err := c.publisher.Request(ctx, c.subjects.Query, msg)
	if err != nil {
		return nil, err
	}

	var resp QueryResponse
	if err := json.Unmarshal(reply.Data, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

// Discover discovers endpoints for a service
func (c *Client) Discover(ctx context.Context, serviceID string, opts ...DiscoverOption) (*DiscoverResponse, error) {
	q := &Query{
		ServiceIDs:    []string{serviceID},
		OnlyAvailable: true,
	}

	for _, opt := range opts {
		opt(q)
	}

	req := &DiscoverRequest{
		ServiceID:   serviceID,
		Version:     q.Version,
		Protocol:    q.Protocol,
		Tags:        q.Tags,
		OnlyHealthy: q.OnlyHealthy,
		TenantID:    q.TenantID,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	msg := &nats.Msg{Data: data, Header: make(nats.Header)}

	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	reply, err := c.publisher.Request(ctx, c.subjects.Discover, msg)
	if err != nil {
		return nil, err
	}

	var resp DiscoverResponse
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
