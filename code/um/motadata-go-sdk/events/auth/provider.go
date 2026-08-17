package auth

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

// Registration errors
var (
	ErrInvalidPayload      = errors.New("invalid registration payload")
	ErrMissingTenantID     = errors.New("tenant ID is required")
	ErrRegistrationFailed  = errors.New("tenant registration failed")
	ErrTenantAlreadyExists = errors.New("tenant already registered")
)

// RegistrationPayload represents credentials from server registration
type RegistrationPayload struct {
	TenantID         string            `json:"tenant_id"`
	JWT              string            `json:"jwt"`
	Seed             string            `json:"seed,omitempty"`
	AccountPublicKey string            `json:"account_public_key,omitempty"`
	IssuedAt         time.Time         `json:"issued_at,omitempty"`
	ExpiresAt        *time.Time        `json:"expires_at,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

// Validate checks if the payload is valid
func (p *RegistrationPayload) Validate() error {
	if p.TenantID == "" {
		return ErrMissingTenantID
	}
	if p.JWT == "" {
		return ErrMissingJWT
	}
	return nil
}

// ParseRegistrationPayload parses JSON into payload
func ParseRegistrationPayload(data []byte) (*RegistrationPayload, error) {
	if len(data) == 0 {
		return nil, ErrInvalidPayload
	}

	var payload RegistrationPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, ErrInvalidPayload
	}

	if err := payload.Validate(); err != nil {
		return nil, err
	}

	return &payload, nil
}

// RegistrationRequest for requesting credentials from server
type RegistrationRequest struct {
	TenantID    string             `json:"tenant_id"`
	ServiceName string             `json:"service_name"`
	ServiceID   string             `json:"service_id,omitempty"`
	Permissions *PermissionRequest `json:"permissions,omitempty"`
	Metadata    map[string]string  `json:"metadata,omitempty"`
}

// PermissionRequest specifies requested permissions
type PermissionRequest struct {
	PublishSubjects   []string `json:"publish,omitempty"`
	SubscribeSubjects []string `json:"subscribe,omitempty"`
}

// ToJSON converts request to JSON
func (r *RegistrationRequest) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

// RegistrationResponse from server
type RegistrationResponse struct {
	Success bool                 `json:"success"`
	Error   string               `json:"error,omitempty"`
	Payload *RegistrationPayload `json:"payload,omitempty"`
}

// ParseRegistrationResponse parses JSON response
func ParseRegistrationResponse(data []byte) (*RegistrationResponse, error) {
	if len(data) == 0 {
		return nil, ErrInvalidPayload
	}

	var resp RegistrationResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, ErrInvalidPayload
	}

	return &resp, nil
}

// RegistrationHandler handles tenant registration events
type RegistrationHandler struct {
	credManager *CredentialManager

	// Callbacks
	onRegistered   func(tenantID string)
	onUnregistered func(tenantID string)
	onError        func(tenantID string, err error)
}

// NewRegistrationHandler creates a new handler
func NewRegistrationHandler(credManager *CredentialManager) *RegistrationHandler {
	return &RegistrationHandler{
		credManager: credManager,
	}
}

// OnRegistered sets callback for successful registrations
func (h *RegistrationHandler) OnRegistered(cb func(tenantID string)) {
	h.onRegistered = cb
}

// OnUnregistered sets callback for unregistrations
func (h *RegistrationHandler) OnUnregistered(cb func(tenantID string)) {
	h.onUnregistered = cb
}

// OnError sets callback for errors
func (h *RegistrationHandler) OnError(cb func(tenantID string, err error)) {
	h.onError = cb
}

// HandleRegistration processes a registration payload
func (h *RegistrationHandler) HandleRegistration(payload *RegistrationPayload) error {
	if payload == nil {
		return ErrInvalidPayload
	}

	if err := payload.Validate(); err != nil {
		h.notifyError(payload.TenantID, err)
		return err
	}

	if err := h.credManager.RegisterFromPayload(*payload); err != nil {
		h.notifyError(payload.TenantID, err)
		return err
	}

	if h.onRegistered != nil {
		h.onRegistered(payload.TenantID)
	}

	return nil
}

// HandleRegistrationJSON processes JSON payload
func (h *RegistrationHandler) HandleRegistrationJSON(data []byte) error {
	payload, err := ParseRegistrationPayload(data)
	if err != nil {
		h.notifyError("", err)
		return err
	}

	return h.HandleRegistration(payload)
}

// HandleUnregistration processes unregistration
func (h *RegistrationHandler) HandleUnregistration(tenantID string) {
	if tenantID == "" {
		return
	}

	h.credManager.Remove(tenantID)

	if h.onUnregistered != nil {
		h.onUnregistered(tenantID)
	}
}

func (h *RegistrationHandler) notifyError(tenantID string, err error) {
	if h.onError != nil {
		h.onError(tenantID, err)
	}
}

// StaticProvider provides static credentials
type StaticProvider struct {
	credentials *Credentials
}

// NewStaticProvider creates a provider with static credentials
func NewStaticProvider(creds *Credentials) *StaticProvider {
	return &StaticProvider{credentials: creds}
}

// Type returns the auth type
func (p *StaticProvider) Type() Type {
	if p.credentials != nil {
		return p.credentials.Type
	}
	return TypeNone
}

// Authenticate returns the static credentials
func (p *StaticProvider) Authenticate(ctx context.Context) (*Credentials, error) {
	if p.credentials == nil {
		return nil, ErrInvalidCredential
	}
	return p.credentials, nil
}

// Refresh returns the same credentials (no refresh for static)
func (p *StaticProvider) Refresh(ctx context.Context) (*Credentials, error) {
	return p.Authenticate(ctx)
}

// IsValid returns true if credentials exist
func (p *StaticProvider) IsValid() bool {
	return p.credentials != nil && !p.credentials.IsEmpty()
}

// ManagedProvider provides credentials from CredentialManager
type ManagedProvider struct {
	manager  *CredentialManager
	tenantID string
}

// NewManagedProvider creates a provider backed by CredentialManager
func NewManagedProvider(manager *CredentialManager, tenantID string) *ManagedProvider {
	return &ManagedProvider{
		manager:  manager,
		tenantID: tenantID,
	}
}

// Type returns JWT type
func (p *ManagedProvider) Type() Type {
	return TypeJWT
}

// Authenticate returns credentials from manager
func (p *ManagedProvider) Authenticate(ctx context.Context) (*Credentials, error) {
	return p.manager.GetCredentials(p.tenantID)
}

// Refresh triggers refresh and returns new credentials
func (p *ManagedProvider) Refresh(ctx context.Context) (*Credentials, error) {
	// Manager handles refresh automatically via callbacks
	return p.Authenticate(ctx)
}

// IsValid checks if tenant has valid credentials
func (p *ManagedProvider) IsValid() bool {
	return p.manager.HasValid(p.tenantID)
}

// ErrInvalidCredential error
var ErrInvalidCredential = errors.New("invalid credentials")
