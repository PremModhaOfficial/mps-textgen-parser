package auth

import (
	"context"
	"testing"
)

// ============== Registration Payload Tests ==============

func TestRegistrationPayloadValidate(t *testing.T) {
	// Valid payload
	p := &RegistrationPayload{
		TenantID: "tenant-1",
		JWT:      "jwt-token",
	}
	if err := p.Validate(); err != nil {
		t.Errorf("valid payload should not return error: %v", err)
	}

	// Missing tenant ID
	p = &RegistrationPayload{
		JWT: "jwt-token",
	}
	if err := p.Validate(); err != ErrMissingTenantID {
		t.Errorf("expected ErrMissingTenantID, got %v", err)
	}

	// Missing JWT
	p = &RegistrationPayload{
		TenantID: "tenant-1",
	}
	if err := p.Validate(); err != ErrMissingJWT {
		t.Errorf("expected ErrMissingJWT, got %v", err)
	}
}

func TestParseRegistrationPayload(t *testing.T) {
	// Valid payload
	data := []byte(`{"tenant_id":"tenant-1","jwt":"jwt-token"}`)
	p, err := ParseRegistrationPayload(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.TenantID != "tenant-1" || p.JWT != "jwt-token" {
		t.Errorf("unexpected payload: %+v", p)
	}

	// Empty data
	_, err = ParseRegistrationPayload(nil)
	if err != ErrInvalidPayload {
		t.Errorf("expected ErrInvalidPayload for empty data, got %v", err)
	}

	// Invalid JSON
	_, err = ParseRegistrationPayload([]byte("not json"))
	if err != ErrInvalidPayload {
		t.Errorf("expected ErrInvalidPayload for invalid JSON, got %v", err)
	}

	// Valid JSON but missing fields
	_, err = ParseRegistrationPayload([]byte(`{"tenant_id":"tenant-1"}`))
	if err != ErrMissingJWT {
		t.Errorf("expected ErrMissingJWT, got %v", err)
	}
}

func TestRegistrationRequestToJSON(t *testing.T) {
	r := &RegistrationRequest{
		TenantID:    "tenant-1",
		ServiceName: "my-service",
		Permissions: &PermissionRequest{
			PublishSubjects:   []string{"events.>"},
			SubscribeSubjects: []string{"events.>"},
		},
		Metadata: map[string]string{"key": "value"},
	}

	data, err := r.ToJSON()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Error("expected non-empty JSON")
	}
}

func TestParseRegistrationResponse(t *testing.T) {
	// Valid response
	data := []byte(`{"success":true,"payload":{"tenant_id":"t1","jwt":"j1"}}`)
	r, err := ParseRegistrationResponse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Success {
		t.Error("expected success to be true")
	}
	if r.Payload == nil || r.Payload.TenantID != "t1" {
		t.Errorf("unexpected payload: %+v", r.Payload)
	}

	// Empty data
	_, err = ParseRegistrationResponse(nil)
	if err != ErrInvalidPayload {
		t.Errorf("expected ErrInvalidPayload, got %v", err)
	}

	// Invalid JSON
	_, err = ParseRegistrationResponse([]byte("not json"))
	if err != ErrInvalidPayload {
		t.Errorf("expected ErrInvalidPayload, got %v", err)
	}
}

// ============== Registration Handler Tests ==============

func TestNewRegistrationHandler(t *testing.T) {
	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	h := NewRegistrationHandler(cm)
	if h == nil {
		t.Fatal("NewRegistrationHandler returned nil")
	}
	if h.credManager != cm {
		t.Error("credential manager not set correctly")
	}
}

func TestRegistrationHandlerCallbacks(t *testing.T) {
	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	h := NewRegistrationHandler(cm)

	var registeredID, unregisteredID string
	var errorID string
	var errorErr error

	h.OnRegistered(func(tenantID string) {
		registeredID = tenantID
	})
	h.OnUnregistered(func(tenantID string) {
		unregisteredID = tenantID
	})
	h.OnError(func(tenantID string, err error) {
		errorID = tenantID
		errorErr = err
	})

	// Test registration
	payload := &RegistrationPayload{
		TenantID: "tenant-1",
		JWT:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig",
	}
	if err := h.HandleRegistration(payload); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if registeredID != "tenant-1" {
		t.Errorf("expected registered callback with tenant-1, got %s", registeredID)
	}

	// Test unregistration
	h.HandleUnregistration("tenant-1")
	if unregisteredID != "tenant-1" {
		t.Errorf("expected unregistered callback with tenant-1, got %s", unregisteredID)
	}

	// Test error callback
	invalidPayload := &RegistrationPayload{TenantID: "t1"} // Missing JWT
	h.HandleRegistration(invalidPayload)
	if errorID != "t1" {
		t.Errorf("expected error callback with t1, got %s", errorID)
	}
	if errorErr == nil {
		t.Error("expected error in callback")
	}
}

func TestRegistrationHandlerHandleNil(t *testing.T) {
	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	h := NewRegistrationHandler(cm)

	if err := h.HandleRegistration(nil); err != ErrInvalidPayload {
		t.Errorf("expected ErrInvalidPayload for nil payload, got %v", err)
	}
}

func TestRegistrationHandlerHandleJSON(t *testing.T) {
	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	h := NewRegistrationHandler(cm)

	var registeredID string
	h.OnRegistered(func(tenantID string) {
		registeredID = tenantID
	})

	// Valid JSON
	data := []byte(`{"tenant_id":"tenant-2","jwt":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig"}`)
	if err := h.HandleRegistrationJSON(data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if registeredID != "tenant-2" {
		t.Errorf("expected tenant-2, got %s", registeredID)
	}

	// Invalid JSON
	if err := h.HandleRegistrationJSON([]byte("invalid")); err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestRegistrationHandlerUnregistrationEmpty(t *testing.T) {
	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	h := NewRegistrationHandler(cm)

	// Should not panic
	h.HandleUnregistration("")
}

// ============== Static Provider Tests ==============

func TestNewStaticProvider(t *testing.T) {
	creds := UserPass("user", "pass")
	p := NewStaticProvider(creds)
	if p == nil {
		t.Fatal("NewStaticProvider returned nil")
	}
}

func TestStaticProviderType(t *testing.T) {
	// With credentials
	creds := UserPass("user", "pass")
	p := NewStaticProvider(creds)
	if p.Type() != TypeUserPass {
		t.Errorf("expected TypeUserPass, got %v", p.Type())
	}

	// With nil credentials
	p = NewStaticProvider(nil)
	if p.Type() != TypeNone {
		t.Errorf("expected TypeNone for nil credentials, got %v", p.Type())
	}
}

func TestStaticProviderAuthenticate(t *testing.T) {
	ctx := context.Background()

	// Valid credentials
	creds := Token("my-token")
	p := NewStaticProvider(creds)
	result, err := p.Authenticate(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Token != "my-token" {
		t.Errorf("unexpected token: %s", result.Token)
	}

	// Nil credentials
	p = NewStaticProvider(nil)
	_, err = p.Authenticate(ctx)
	if err != ErrInvalidCredential {
		t.Errorf("expected ErrInvalidCredential, got %v", err)
	}
}

func TestStaticProviderRefresh(t *testing.T) {
	ctx := context.Background()
	creds := Token("my-token")
	p := NewStaticProvider(creds)

	// Refresh returns same credentials
	result, err := p.Refresh(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != creds {
		t.Error("refresh should return same credentials")
	}
}

func TestStaticProviderIsValid(t *testing.T) {
	// Valid credentials
	creds := Token("my-token")
	p := NewStaticProvider(creds)
	if !p.IsValid() {
		t.Error("expected IsValid to return true")
	}

	// Nil credentials
	p = NewStaticProvider(nil)
	if p.IsValid() {
		t.Error("expected IsValid to return false for nil")
	}

	// Empty credentials
	p = NewStaticProvider(&Credentials{})
	if p.IsValid() {
		t.Error("expected IsValid to return false for empty")
	}
}

// ============== Managed Provider Tests ==============

func TestNewManagedProvider(t *testing.T) {
	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	p := NewManagedProvider(cm, "tenant-1")
	if p == nil {
		t.Fatal("NewManagedProvider returned nil")
	}
	if p.tenantID != "tenant-1" {
		t.Errorf("unexpected tenant ID: %s", p.tenantID)
	}
}

func TestManagedProviderType(t *testing.T) {
	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	p := NewManagedProvider(cm, "tenant-1")
	if p.Type() != TypeJWT {
		t.Errorf("expected TypeJWT, got %v", p.Type())
	}
}

func TestManagedProviderAuthenticate(t *testing.T) {
	ctx := context.Background()
	cm := NewCredentialManager(DefaultCredentialManagerConfig())

	// Register credentials
	cm.Register("tenant-1", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig", "seed")

	p := NewManagedProvider(cm, "tenant-1")
	creds, err := p.Authenticate(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds.JWT != "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig" {
		t.Errorf("unexpected JWT: %s", creds.JWT)
	}

	// Non-existent tenant
	p2 := NewManagedProvider(cm, "unknown")
	_, err = p2.Authenticate(ctx)
	if err == nil {
		t.Error("expected error for unknown tenant")
	}
}

func TestManagedProviderRefresh(t *testing.T) {
	ctx := context.Background()
	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	cm.Register("tenant-1", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig", "seed")

	p := NewManagedProvider(cm, "tenant-1")
	creds, err := p.Refresh(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds == nil {
		t.Error("expected credentials from refresh")
	}
}

func TestManagedProviderIsValid(t *testing.T) {
	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	cm.Register("tenant-1", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig", "seed")

	// Valid credentials
	p := NewManagedProvider(cm, "tenant-1")
	if !p.IsValid() {
		t.Error("expected IsValid to return true")
	}

	// Non-existent tenant
	p2 := NewManagedProvider(cm, "unknown")
	if p2.IsValid() {
		t.Error("expected IsValid to return false for unknown tenant")
	}
}

// ============== Benchmarks ==============

func BenchmarkStaticProviderAuthenticate(b *testing.B) {
	ctx := context.Background()
	creds := Token("my-token")
	p := NewStaticProvider(creds)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		p.Authenticate(ctx)
	}
}

func BenchmarkParseRegistrationPayload(b *testing.B) {
	data := []byte(`{"tenant_id":"tenant-1","jwt":"jwt-token","seed":"seed-value"}`)
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ParseRegistrationPayload(data)
	}
}
