package microservice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

// ============== Error Tests ==============

func TestServiceError(t *testing.T) {
	tests := []struct {
		name     string
		err      ServiceError
		expected string
	}{
		{
			name: "with service name and op",
			err: ServiceError{
				ServiceID:   "svc1",
				ServiceName: "MyService",
				Op:          "register",
				Err:         ErrServiceExists,
			},
			expected: "service MyService (svc1): register: service already exists",
		},
		{
			name: "with op only",
			err: ServiceError{
				ServiceID: "svc1",
				Op:        "register",
				Err:       ErrServiceExists,
			},
			expected: "service svc1: register: service already exists",
		},
		{
			name: "minimal",
			err: ServiceError{
				ServiceID: "svc1",
				Err:       ErrServiceExists,
			},
			expected: "service svc1: service already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %q, want %q", got, tt.expected)
			}
			if unwrapped := tt.err.Unwrap(); unwrapped != tt.err.Err {
				t.Errorf("Unwrap() = %v, want %v", unwrapped, tt.err.Err)
			}
		})
	}
}

func TestNewServiceError(t *testing.T) {
	err := NewServiceError("svc1", "MyService", "register", ErrServiceExists)
	if err.ServiceID != "svc1" || err.ServiceName != "MyService" || err.Op != "register" {
		t.Error("NewServiceError did not set fields correctly")
	}
}

func TestInstanceError(t *testing.T) {
	tests := []struct {
		name     string
		err      InstanceError
		expected string
	}{
		{
			name: "with service ID and op",
			err: InstanceError{
				InstanceID: "inst1",
				ServiceID:  "svc1",
				Op:         "register",
				Err:        ErrInstanceExists,
			},
			expected: "instance inst1 (service svc1): register: instance already exists",
		},
		{
			name: "with op only",
			err: InstanceError{
				InstanceID: "inst1",
				Op:         "register",
				Err:        ErrInstanceExists,
			},
			expected: "instance inst1: register: instance already exists",
		},
		{
			name: "minimal",
			err: InstanceError{
				InstanceID: "inst1",
				Err:        ErrInstanceExists,
			},
			expected: "instance inst1: instance already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %q, want %q", got, tt.expected)
			}
			if unwrapped := tt.err.Unwrap(); unwrapped != tt.err.Err {
				t.Errorf("Unwrap() = %v, want %v", unwrapped, tt.err.Err)
			}
		})
	}
}

func TestNewInstanceError(t *testing.T) {
	err := NewInstanceError("inst1", "svc1", "register", ErrInstanceExists)
	if err.InstanceID != "inst1" || err.ServiceID != "svc1" || err.Op != "register" {
		t.Error("NewInstanceError did not set fields correctly")
	}
}

func TestRegistrationError(t *testing.T) {
	tests := []struct {
		name     string
		err      RegistrationError
		expected string
	}{
		{
			name: "with instance ID and reason",
			err: RegistrationError{
				ServiceID:  "svc1",
				InstanceID: "inst1",
				Reason:     "validation failed",
				Err:        ErrInvalidService,
			},
			expected: "registration failed for inst1@svc1: validation failed: invalid service",
		},
		{
			name: "without instance ID",
			err: RegistrationError{
				ServiceID: "svc1",
				Reason:    "validation failed",
				Err:       ErrInvalidService,
			},
			expected: "registration failed for svc1: validation failed: invalid service",
		},
		{
			name: "without reason",
			err: RegistrationError{
				ServiceID:  "svc1",
				InstanceID: "inst1",
				Err:        ErrInvalidService,
			},
			expected: "registration failed for inst1@svc1: invalid service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("Error() = %q, want %q", got, tt.expected)
			}
			if unwrapped := tt.err.Unwrap(); unwrapped != tt.err.Err {
				t.Errorf("Unwrap() = %v, want %v", unwrapped, tt.err.Err)
			}
		})
	}
}

func TestHealthCheckError(t *testing.T) {
	err := HealthCheckError{
		ServiceID:        "svc1",
		InstanceID:       "inst1",
		URL:              "http://localhost:8080/health",
		Err:              errors.New("connection refused"),
		ConsecutiveFails: 3,
	}

	expected := "health check failed for inst1@svc1 at http://localhost:8080/health (consecutive fails: 3): connection refused"
	if got := err.Error(); got != expected {
		t.Errorf("Error() = %q, want %q", got, expected)
	}

	if unwrapped := err.Unwrap(); unwrapped.Error() != "connection refused" {
		t.Errorf("Unwrap() = %v, want connection refused", unwrapped)
	}

	// Test without instance ID
	err2 := HealthCheckError{
		ServiceID:        "svc1",
		URL:              "http://localhost:8080/health",
		Err:              errors.New("timeout"),
		ConsecutiveFails: 1,
	}
	if !contains(err2.Error(), "svc1") {
		t.Error("Error message should contain service ID")
	}
}

func TestIsServiceError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"ServiceError", ServiceError{ServiceID: "svc1", Err: ErrServiceNotFound}, true},
		{"InstanceError", InstanceError{InstanceID: "inst1", Err: ErrInstanceNotFound}, true},
		{"RegistrationError", RegistrationError{ServiceID: "svc1", Err: ErrInvalidService}, true},
		{"HealthCheckError", HealthCheckError{ServiceID: "svc1", Err: errors.New("fail")}, true},
		{"regular error", errors.New("random error"), false},
		{"nil error", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsServiceError(tt.err); got != tt.expected {
				t.Errorf("IsServiceError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"ErrServiceNotFound", ErrServiceNotFound, true},
		{"ErrInstanceNotFound", ErrInstanceNotFound, true},
		{"wrapped ErrServiceNotFound", ServiceError{Err: ErrServiceNotFound}, true},
		{"other error", ErrServiceExists, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotFound(tt.err); got != tt.expected {
				t.Errorf("IsNotFound() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsExists(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"ErrServiceExists", ErrServiceExists, true},
		{"ErrInstanceExists", ErrInstanceExists, true},
		{"wrapped ErrServiceExists", ServiceError{Err: ErrServiceExists}, true},
		{"other error", ErrServiceNotFound, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsExists(tt.err); got != tt.expected {
				t.Errorf("IsExists() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// ============== Status Tests ==============

func TestStatusString(t *testing.T) {
	tests := []struct {
		status   Status
		expected string
	}{
		{StatusUnknown, "unknown"},
		{StatusStarting, "starting"},
		{StatusRunning, "running"},
		{StatusDegraded, "degraded"},
		{StatusStopping, "stopping"},
		{StatusStopped, "stopped"},
		{StatusFailed, "failed"},
		{StatusMaintenance, "maintenance"},
		{Status(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.status.String(); got != tt.expected {
				t.Errorf("String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestStatusIsAvailable(t *testing.T) {
	tests := []struct {
		status   Status
		expected bool
	}{
		{StatusRunning, true},
		{StatusDegraded, true},
		{StatusStarting, false},
		{StatusStopping, false},
		{StatusStopped, false},
		{StatusFailed, false},
		{StatusUnknown, false},
	}

	for _, tt := range tests {
		t.Run(tt.status.String(), func(t *testing.T) {
			if got := tt.status.IsAvailable(); got != tt.expected {
				t.Errorf("IsAvailable() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStatusIsHealthy(t *testing.T) {
	tests := []struct {
		status   Status
		expected bool
	}{
		{StatusRunning, true},
		{StatusDegraded, false},
		{StatusStarting, false},
		{StatusStopping, false},
		{StatusFailed, false},
	}

	for _, tt := range tests {
		t.Run(tt.status.String(), func(t *testing.T) {
			if got := tt.status.IsHealthy(); got != tt.expected {
				t.Errorf("IsHealthy() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// ============== Info Tests ==============

func TestInfoClone(t *testing.T) {
	// Test nil clone
	var nilInfo *Info
	if nilInfo.Clone() != nil {
		t.Error("Clone of nil should return nil")
	}

	original := &Info{
		ID:           "svc1",
		Name:         "MyService",
		Dependencies: []string{"dep1", "dep2"},
		Capabilities: []string{"cap1"},
		Tags:         []string{"tag1", "tag2"},
		Metadata:     map[string]string{"key": "value"},
		Instances: []*Instance{
			{ID: "inst1", Host: "localhost"},
		},
	}

	clone := original.Clone()

	// Verify deep copy
	if clone.ID != original.ID {
		t.Error("ID not copied correctly")
	}

	// Modify original and check clone is unaffected
	original.Dependencies[0] = "modified"
	if clone.Dependencies[0] == "modified" {
		t.Error("Dependencies not deep copied")
	}

	original.Capabilities[0] = "modified"
	if clone.Capabilities[0] == "modified" {
		t.Error("Capabilities not deep copied")
	}

	original.Tags[0] = "modified"
	if clone.Tags[0] == "modified" {
		t.Error("Tags not deep copied")
	}

	original.Metadata["key"] = "modified"
	if clone.Metadata["key"] == "modified" {
		t.Error("Metadata not deep copied")
	}

	original.Instances[0].Host = "modified"
	if clone.Instances[0].Host == "modified" {
		t.Error("Instances not deep copied")
	}
}

func TestInfoHasCapability(t *testing.T) {
	info := &Info{
		Capabilities: []string{"read", "write", "admin"},
	}

	if !info.HasCapability("read") {
		t.Error("HasCapability should return true for existing capability")
	}
	if info.HasCapability("delete") {
		t.Error("HasCapability should return false for non-existing capability")
	}
}

func TestInfoHasTag(t *testing.T) {
	info := &Info{
		Tags: []string{"production", "critical"},
	}

	if !info.HasTag("production") {
		t.Error("HasTag should return true for existing tag")
	}
	if info.HasTag("staging") {
		t.Error("HasTag should return false for non-existing tag")
	}
}

func TestInfoMetadata(t *testing.T) {
	info := &Info{}

	// Test GetMetadata with nil map
	if info.GetMetadata("key") != "" {
		t.Error("GetMetadata should return empty string for nil map")
	}

	// Test SetMetadata initializes map
	info.SetMetadata("key", "value")
	if info.GetMetadata("key") != "value" {
		t.Error("SetMetadata/GetMetadata not working correctly")
	}

	// Test GetMetadata with missing key
	if info.GetMetadata("missing") != "" {
		t.Error("GetMetadata should return empty string for missing key")
	}
}

func TestInfoHealthyInstanceCount(t *testing.T) {
	info := &Info{
		Instances: []*Instance{
			{ID: "inst1", Status: StatusRunning},
			{ID: "inst2", Status: StatusDegraded},
			{ID: "inst3", Status: StatusRunning},
			{ID: "inst4", Status: StatusFailed},
		},
	}

	if count := info.HealthyInstanceCount(); count != 2 {
		t.Errorf("HealthyInstanceCount() = %d, want 2", count)
	}
}

func TestInfoAvailableInstanceCount(t *testing.T) {
	info := &Info{
		Instances: []*Instance{
			{ID: "inst1", Status: StatusRunning},
			{ID: "inst2", Status: StatusDegraded},
			{ID: "inst3", Status: StatusRunning},
			{ID: "inst4", Status: StatusFailed},
		},
	}

	if count := info.AvailableInstanceCount(); count != 3 {
		t.Errorf("AvailableInstanceCount() = %d, want 3", count)
	}
}

// ============== Instance Tests ==============

func TestInstanceClone(t *testing.T) {
	// Test nil clone
	var nilInst *Instance
	if nilInst.Clone() != nil {
		t.Error("Clone of nil should return nil")
	}

	original := &Instance{
		ID:       "inst1",
		Host:     "localhost",
		Metadata: map[string]string{"key": "value"},
	}

	clone := original.Clone()
	if clone.ID != original.ID {
		t.Error("ID not copied correctly")
	}

	original.Metadata["key"] = "modified"
	if clone.Metadata["key"] == "modified" {
		t.Error("Metadata not deep copied")
	}
}

func TestInstanceAddress(t *testing.T) {
	tests := []struct {
		name     string
		instance *Instance
		expected string
	}{
		{
			name:     "with port",
			instance: &Instance{Host: "localhost", Port: 8080},
			expected: "localhost:8080",
		},
		{
			name:     "without port",
			instance: &Instance{Host: "localhost", Port: 0},
			expected: "localhost",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.instance.Address(); got != tt.expected {
				t.Errorf("Address() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestInstanceIsExpired(t *testing.T) {
	tests := []struct {
		name      string
		expiresAt time.Time
		expected  bool
	}{
		{"zero time", time.Time{}, false},
		{"future time", time.Now().Add(time.Hour), false},
		{"past time", time.Now().Add(-time.Hour), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inst := &Instance{ExpiresAt: tt.expiresAt}
			if got := inst.IsExpired(); got != tt.expected {
				t.Errorf("IsExpired() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestInstanceIsHealthy(t *testing.T) {
	tests := []struct {
		status   Status
		expected bool
	}{
		{StatusRunning, true},
		{StatusDegraded, false},
		{StatusFailed, false},
	}

	for _, tt := range tests {
		t.Run(tt.status.String(), func(t *testing.T) {
			inst := &Instance{Status: tt.status}
			if got := inst.IsHealthy(); got != tt.expected {
				t.Errorf("IsHealthy() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestInstanceMetadata(t *testing.T) {
	inst := &Instance{}

	// Test GetMetadata with nil map
	if inst.GetMetadata("key") != "" {
		t.Error("GetMetadata should return empty string for nil map")
	}

	// Test SetMetadata initializes map
	inst.SetMetadata("key", "value")
	if inst.GetMetadata("key") != "value" {
		t.Error("SetMetadata/GetMetadata not working correctly")
	}
}

// ============== Request Validation Tests ==============

func TestRegistrationRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     *RegistrationRequest
		wantErr error
	}{
		{
			name:    "nil service",
			req:     &RegistrationRequest{Service: nil},
			wantErr: ErrInvalidService,
		},
		{
			name:    "missing service ID",
			req:     &RegistrationRequest{Service: &Info{Name: "MyService"}},
			wantErr: ErrMissingServiceID,
		},
		{
			name:    "missing service name",
			req:     &RegistrationRequest{Service: &Info{ID: "svc1"}},
			wantErr: ErrMissingServiceName,
		},
		{
			name:    "valid request",
			req:     &RegistrationRequest{Service: &Info{ID: "svc1", Name: "MyService"}},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstanceRegistrationRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     *InstanceRegistrationRequest
		wantErr error
	}{
		{
			name:    "nil instance",
			req:     &InstanceRegistrationRequest{Instance: nil},
			wantErr: ErrInvalidInstance,
		},
		{
			name:    "missing instance ID",
			req:     &InstanceRegistrationRequest{Instance: &Instance{ServiceID: "svc1", Host: "localhost"}},
			wantErr: ErrMissingInstanceID,
		},
		{
			name:    "missing service ID",
			req:     &InstanceRegistrationRequest{Instance: &Instance{ID: "inst1", Host: "localhost"}},
			wantErr: ErrMissingServiceID,
		},
		{
			name:    "missing host",
			req:     &InstanceRegistrationRequest{Instance: &Instance{ID: "inst1", ServiceID: "svc1"}},
			wantErr: ErrMissingHost,
		},
		{
			name:    "valid request",
			req:     &InstanceRegistrationRequest{Instance: &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"}},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeregistrationRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     *DeregistrationRequest
		wantErr error
	}{
		{
			name:    "missing service ID",
			req:     &DeregistrationRequest{},
			wantErr: ErrMissingServiceID,
		},
		{
			name:    "valid request",
			req:     &DeregistrationRequest{ServiceID: "svc1"},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

// ============== Query Tests ==============

func TestQueryMatches(t *testing.T) {
	service := &Info{
		ID:           "svc1",
		Name:         "MyService",
		Type:         ServiceTypeAPI,
		Version:      "1.0.0",
		APIVersion:   "v1",
		Status:       StatusRunning,
		TenantID:     "tenant1",
		Capabilities: []string{"read", "write"},
		Tags:         []string{"production"},
		Metadata:     map[string]string{"env": "prod"},
		Instances:    []*Instance{{ID: "inst1"}},
	}

	tests := []struct {
		name     string
		query    *Query
		expected bool
	}{
		{"nil service", &Query{}, false},
		{"empty query", &Query{}, true},
		{"matching IDs", &Query{IDs: []string{"svc1", "svc2"}}, true},
		{"non-matching IDs", &Query{IDs: []string{"other"}}, false},
		{"matching names", &Query{Names: []string{"MyService"}}, true},
		{"non-matching names", &Query{Names: []string{"Other"}}, false},
		{"matching types", &Query{Types: []ServiceType{ServiceTypeAPI}}, true},
		{"non-matching types", &Query{Types: []ServiceType{ServiceTypeWorker}}, false},
		{"matching statuses", &Query{Statuses: []Status{StatusRunning}}, true},
		{"non-matching statuses", &Query{Statuses: []Status{StatusStopped}}, false},
		{"only healthy - healthy", &Query{OnlyHealthy: true}, true},
		{"only available - available", &Query{OnlyAvailable: true}, true},
		{"matching capabilities", &Query{Capabilities: []string{"read"}}, true},
		{"non-matching capabilities", &Query{Capabilities: []string{"admin"}}, false},
		{"matching tags", &Query{Tags: []string{"production"}}, true},
		{"non-matching tags", &Query{Tags: []string{"staging"}}, false},
		{"matching metadata", &Query{Metadata: map[string]string{"env": "prod"}}, true},
		{"non-matching metadata", &Query{Metadata: map[string]string{"env": "dev"}}, false},
		{"matching version", &Query{Version: "1.0.0"}, true},
		{"non-matching version", &Query{Version: "2.0.0"}, false},
		{"matching API version", &Query{APIVersion: "v1"}, true},
		{"non-matching API version", &Query{APIVersion: "v2"}, false},
		{"matching tenant", &Query{TenantID: "tenant1"}, true},
		{"non-matching tenant", &Query{TenantID: "tenant2"}, false},
		{"has instances - has", &Query{HasInstances: true}, true},
		{"min instances - satisfied", &Query{MinInstances: 1}, true},
		{"min instances - not satisfied", &Query{MinInstances: 5}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var svc *Info
			if tt.name != "nil service" {
				svc = service
			}
			if got := tt.query.Matches(svc); got != tt.expected {
				t.Errorf("Matches() = %v, want %v", got, tt.expected)
			}
		})
	}

	// Test unhealthy and unavailable
	unhealthyService := &Info{ID: "svc2", Name: "Unhealthy", Status: StatusFailed}
	if (&Query{OnlyHealthy: true}).Matches(unhealthyService) {
		t.Error("OnlyHealthy should not match unhealthy service")
	}
	if (&Query{OnlyAvailable: true}).Matches(unhealthyService) {
		t.Error("OnlyAvailable should not match unavailable service")
	}

	// Test empty instances
	noInstancesService := &Info{ID: "svc3", Name: "NoInstances", Status: StatusRunning, Instances: []*Instance{}}
	if (&Query{HasInstances: true}).Matches(noInstancesService) {
		t.Error("HasInstances should not match service with no instances")
	}
}

// ============== Event Tests ==============

func TestEventToMessage(t *testing.T) {
	event := &Event{
		Type:       EventServiceRegistered,
		ServiceID:  "svc1",
		InstanceID: "inst1",
		OldStatus:  StatusUnknown,
		NewStatus:  StatusRunning,
		Timestamp:  time.Now(),
	}

	msg, err := event.ToMessage()
	if err != nil {
		t.Fatalf("ToMessage() error: %v", err)
	}

	if msg.Headers.Get("event-type") != string(EventServiceRegistered) {
		t.Error("event-type header not set correctly")
	}
	if msg.Headers.Get("service-id") != "svc1" {
		t.Error("service-id header not set correctly")
	}
	if msg.Headers.Get("instance-id") != "inst1" {
		t.Error("instance-id header not set correctly")
	}

	// Test without instance ID
	event2 := &Event{
		Type:      EventServiceRegistered,
		ServiceID: "svc1",
		Timestamp: time.Now(),
	}
	msg2, err := event2.ToMessage()
	if err != nil {
		t.Fatalf("ToMessage() error: %v", err)
	}
	if msg2.Headers.Get("instance-id") != "" {
		t.Error("instance-id header should not be set")
	}
}

func TestEventFromMessage(t *testing.T) {
	event := &Event{
		Type:       EventServiceRegistered,
		ServiceID:  "svc1",
		InstanceID: "inst1",
		OldStatus:  StatusUnknown,
		NewStatus:  StatusRunning,
		Timestamp:  time.Now(),
	}

	data, _ := json.Marshal(event)
	msg := core.NewMessage(data)

	parsed, err := EventFromMessage(msg)
	if err != nil {
		t.Fatalf("EventFromMessage() error: %v", err)
	}

	if parsed.Type != event.Type || parsed.ServiceID != event.ServiceID {
		t.Error("Event not parsed correctly")
	}

	// Test invalid JSON
	badMsg := core.NewMessage([]byte("invalid json"))
	_, err = EventFromMessage(badMsg)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

// ============== Registry Tests ==============

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	if r == nil {
		t.Fatal("NewRegistry returned nil")
	}
	if r.services == nil {
		t.Error("services map not initialized")
	}
}

func TestDefaultRegistryConfig(t *testing.T) {
	cfg := DefaultRegistryConfig()
	if cfg.CleanupInterval != time.Minute {
		t.Error("Default CleanupInterval incorrect")
	}
	if !cfg.ExpirationEnabled {
		t.Error("Default ExpirationEnabled should be true")
	}
}

func TestRegistryWithConfig(t *testing.T) {
	cfg := RegistryConfig{
		CleanupInterval:   5 * time.Minute,
		ExpirationEnabled: false,
	}
	r := NewRegistry(cfg)
	if r.cleanupInterval != 5*time.Minute {
		t.Error("CleanupInterval not set correctly")
	}
}

func TestRegistryRegisterService(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	// Test nil service
	if err := r.RegisterService(nil); err != ErrInvalidService {
		t.Errorf("Expected ErrInvalidService, got %v", err)
	}

	// Test missing ID
	if err := r.RegisterService(&Info{Name: "Test"}); err != ErrMissingServiceID {
		t.Errorf("Expected ErrMissingServiceID, got %v", err)
	}

	// Test missing name
	if err := r.RegisterService(&Info{ID: "svc1"}); err != ErrMissingServiceName {
		t.Errorf("Expected ErrMissingServiceName, got %v", err)
	}

	// Test successful registration
	service := &Info{ID: "svc1", Name: "Test", Status: StatusRunning}
	if err := r.RegisterService(service); err != nil {
		t.Fatalf("RegisterService() error: %v", err)
	}

	// Test update existing
	service.Status = StatusDegraded
	if err := r.RegisterService(service); err != nil {
		t.Fatalf("RegisterService() update error: %v", err)
	}

	// Verify update
	got, exists := r.GetService("svc1")
	if !exists || got.Status != StatusDegraded {
		t.Error("Service not updated correctly")
	}
}

func TestRegistryRegisterServiceWithCallbacks(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	var (
		registered        bool
		updated           bool
		statusChanged     bool
		mu                sync.Mutex
	)

	r.OnServiceRegistered(func(svc *Info) {
		mu.Lock()
		registered = true
		mu.Unlock()
	})
	r.OnServiceUpdated(func(old, new *Info) {
		mu.Lock()
		updated = true
		mu.Unlock()
	})
	r.OnServiceStatusChange(func(svc *Info, oldStatus, newStatus Status) {
		mu.Lock()
		statusChanged = true
		mu.Unlock()
	})

	// Register new service
	r.RegisterService(&Info{ID: "svc1", Name: "Test", Status: StatusRunning})
	time.Sleep(10 * time.Millisecond)
	mu.Lock()
	if !registered {
		t.Error("OnServiceRegistered callback not called")
	}
	mu.Unlock()

	// Update service
	r.RegisterService(&Info{ID: "svc1", Name: "Test", Status: StatusDegraded})
	time.Sleep(10 * time.Millisecond)
	mu.Lock()
	if !updated {
		t.Error("OnServiceUpdated callback not called")
	}
	if !statusChanged {
		t.Error("OnServiceStatusChange callback not called")
	}
	mu.Unlock()
}

func TestRegistryDeregisterService(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	// Test missing ID
	if err := r.DeregisterService(""); err != ErrMissingServiceID {
		t.Errorf("Expected ErrMissingServiceID, got %v", err)
	}

	// Test not found
	if err := r.DeregisterService("unknown"); err != ErrServiceNotFound {
		t.Errorf("Expected ErrServiceNotFound, got %v", err)
	}

	// Register and deregister
	r.RegisterService(&Info{ID: "svc1", Name: "Test"})
	r.RegisterInstance(&Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"})

	var deregistered bool
	r.OnServiceDeregistered(func(svc *Info, reason string) {
		deregistered = true
	})

	if err := r.DeregisterServiceWithReason("svc1", "shutdown"); err != nil {
		t.Fatalf("DeregisterServiceWithReason() error: %v", err)
	}

	time.Sleep(10 * time.Millisecond)
	if !deregistered {
		t.Error("OnServiceDeregistered callback not called")
	}

	if r.ServiceExists("svc1") {
		t.Error("Service still exists after deregistration")
	}
	if r.InstanceExists("inst1") {
		t.Error("Instance still exists after service deregistration")
	}
}

func TestRegistryRegisterInstance(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	// Register service first
	r.RegisterService(&Info{ID: "svc1", Name: "Test"})

	// Test nil instance
	if err := r.RegisterInstance(nil); err != ErrInvalidInstance {
		t.Errorf("Expected ErrInvalidInstance, got %v", err)
	}

	// Test missing ID
	if err := r.RegisterInstance(&Instance{ServiceID: "svc1", Host: "localhost"}); err != ErrMissingInstanceID {
		t.Errorf("Expected ErrMissingInstanceID, got %v", err)
	}

	// Test missing service ID
	if err := r.RegisterInstance(&Instance{ID: "inst1", Host: "localhost"}); err != ErrMissingServiceID {
		t.Errorf("Expected ErrMissingServiceID, got %v", err)
	}

	// Test missing host
	if err := r.RegisterInstance(&Instance{ID: "inst1", ServiceID: "svc1"}); err != ErrMissingHost {
		t.Errorf("Expected ErrMissingHost, got %v", err)
	}

	// Test service not found
	if err := r.RegisterInstance(&Instance{ID: "inst1", ServiceID: "unknown", Host: "localhost"}); err != ErrServiceNotFound {
		t.Errorf("Expected ErrServiceNotFound, got %v", err)
	}

	// Test successful registration
	var registered bool
	r.OnInstanceRegistered(func(inst *Instance) {
		registered = true
	})

	inst := &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost", Status: StatusRunning}
	if err := r.RegisterInstance(inst); err != nil {
		t.Fatalf("RegisterInstance() error: %v", err)
	}

	time.Sleep(10 * time.Millisecond)
	if !registered {
		t.Error("OnInstanceRegistered callback not called")
	}

	// Test update existing
	var statusChanged bool
	r.OnInstanceStatusChange(func(inst *Instance, oldStatus, newStatus Status) {
		statusChanged = true
	})

	inst.Status = StatusDegraded
	if err := r.RegisterInstance(inst); err != nil {
		t.Fatalf("RegisterInstance() update error: %v", err)
	}

	time.Sleep(10 * time.Millisecond)
	if !statusChanged {
		t.Error("OnInstanceStatusChange callback not called")
	}
}

func TestRegistryRegisterInstanceWithTTL(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	r.RegisterService(&Info{ID: "svc1", Name: "Test"})

	inst := &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"}
	if err := r.RegisterInstanceWithTTL(inst, 5*time.Minute); err != nil {
		t.Fatalf("RegisterInstanceWithTTL() error: %v", err)
	}

	got, exists := r.GetInstance("inst1")
	if !exists {
		t.Fatal("Instance not found")
	}
	if got.ExpiresAt.IsZero() {
		t.Error("ExpiresAt not set")
	}
}

func TestRegistryDeregisterInstance(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	// Test missing ID
	if err := r.DeregisterInstance(""); err != ErrMissingInstanceID {
		t.Errorf("Expected ErrMissingInstanceID, got %v", err)
	}

	// Test not found
	if err := r.DeregisterInstance("unknown"); err != ErrInstanceNotFound {
		t.Errorf("Expected ErrInstanceNotFound, got %v", err)
	}

	// Register and deregister
	r.RegisterService(&Info{ID: "svc1", Name: "Test"})
	r.RegisterInstance(&Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"})

	var deregistered bool
	r.OnInstanceDeregistered(func(inst *Instance, reason string) {
		deregistered = true
	})

	if err := r.DeregisterInstanceWithReason("inst1", "shutdown"); err != nil {
		t.Fatalf("DeregisterInstanceWithReason() error: %v", err)
	}

	time.Sleep(10 * time.Millisecond)
	if !deregistered {
		t.Error("OnInstanceDeregistered callback not called")
	}

	if r.InstanceExists("inst1") {
		t.Error("Instance still exists after deregistration")
	}
}

func TestRegistryGetService(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	// Test not found
	_, exists := r.GetService("unknown")
	if exists {
		t.Error("GetService should return false for unknown service")
	}

	// Test found
	r.RegisterService(&Info{ID: "svc1", Name: "Test"})
	svc, exists := r.GetService("svc1")
	if !exists || svc.ID != "svc1" {
		t.Error("GetService should return service")
	}
}

func TestRegistryGetServiceByName(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	// Test not found
	if services := r.GetServiceByName("unknown"); len(services) > 0 {
		t.Error("GetServiceByName should return empty for unknown name")
	}

	// Test found
	r.RegisterService(&Info{ID: "svc1", Name: "Test"})
	r.RegisterService(&Info{ID: "svc2", Name: "Test"})

	services := r.GetServiceByName("Test")
	if len(services) != 2 {
		t.Errorf("GetServiceByName() returned %d services, want 2", len(services))
	}
}

func TestRegistryGetInstance(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	// Test not found
	_, exists := r.GetInstance("unknown")
	if exists {
		t.Error("GetInstance should return false for unknown instance")
	}

	// Test found
	r.RegisterService(&Info{ID: "svc1", Name: "Test"})
	r.RegisterInstance(&Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"})

	inst, exists := r.GetInstance("inst1")
	if !exists || inst.ID != "inst1" {
		t.Error("GetInstance should return instance")
	}

	// Test expired instance - set ExpiresAt directly since RegisterInstanceWithTTL only sets for TTL > 0
	expiredInst := &Instance{
		ID:        "inst2",
		ServiceID: "svc1",
		Host:      "localhost",
		ExpiresAt: time.Now().Add(-time.Hour), // Already expired
	}
	r.RegisterInstance(expiredInst)
	_, exists = r.GetInstance("inst2")
	if exists {
		t.Error("GetInstance should not return expired instance")
	}
}

func TestRegistryGetInstances(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	// Test empty
	if instances := r.GetInstances("unknown"); len(instances) > 0 {
		t.Error("GetInstances should return empty for unknown service")
	}

	r.RegisterService(&Info{ID: "svc1", Name: "Test"})
	r.RegisterInstance(&Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"})
	r.RegisterInstance(&Instance{ID: "inst2", ServiceID: "svc1", Host: "localhost"})

	instances := r.GetInstances("svc1")
	if len(instances) != 2 {
		t.Errorf("GetInstances() returned %d instances, want 2", len(instances))
	}
}

func TestRegistryGetHealthyInstances(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	r.RegisterService(&Info{ID: "svc1", Name: "Test"})
	r.RegisterInstance(&Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost", Status: StatusRunning})
	r.RegisterInstance(&Instance{ID: "inst2", ServiceID: "svc1", Host: "localhost", Status: StatusDegraded})
	r.RegisterInstance(&Instance{ID: "inst3", ServiceID: "svc1", Host: "localhost", Status: StatusFailed})

	healthy := r.GetHealthyInstances("svc1")
	if len(healthy) != 1 {
		t.Errorf("GetHealthyInstances() returned %d instances, want 1", len(healthy))
	}
}

func TestRegistryGetAvailableInstances(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	r.RegisterService(&Info{ID: "svc1", Name: "Test"})
	r.RegisterInstance(&Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost", Status: StatusRunning})
	r.RegisterInstance(&Instance{ID: "inst2", ServiceID: "svc1", Host: "localhost", Status: StatusDegraded})
	r.RegisterInstance(&Instance{ID: "inst3", ServiceID: "svc1", Host: "localhost", Status: StatusFailed})

	available := r.GetAvailableInstances("svc1")
	if len(available) != 2 {
		t.Errorf("GetAvailableInstances() returned %d instances, want 2", len(available))
	}
}

func TestRegistryQueryServices(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	r.RegisterService(&Info{ID: "svc1", Name: "A-Service", Status: StatusRunning, Tags: []string{"prod"}})
	r.RegisterService(&Info{ID: "svc2", Name: "B-Service", Status: StatusDegraded, Tags: []string{"dev"}})
	r.RegisterService(&Info{ID: "svc3", Name: "C-Service", Status: StatusFailed})

	// Test nil query
	all := r.QueryServices(nil)
	if len(all) != 3 {
		t.Errorf("QueryServices(nil) returned %d services, want 3", len(all))
	}

	// Test only healthy
	healthy := r.QueryServices(&Query{OnlyHealthy: true})
	if len(healthy) != 1 {
		t.Errorf("QueryServices(OnlyHealthy) returned %d services, want 1", len(healthy))
	}

	// Test pagination
	limited := r.QueryServices(&Query{Limit: 1})
	if len(limited) != 1 {
		t.Errorf("QueryServices(Limit:1) returned %d services, want 1", len(limited))
	}

	offset := r.QueryServices(&Query{Offset: 2})
	if len(offset) != 1 {
		t.Errorf("QueryServices(Offset:2) returned %d services, want 1", len(offset))
	}

	// Test offset beyond length
	beyondOffset := r.QueryServices(&Query{Offset: 100})
	if len(beyondOffset) != 0 {
		t.Errorf("QueryServices(Offset:100) returned %d services, want 0", len(beyondOffset))
	}
}

func TestRegistryAllServices(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	r.RegisterService(&Info{ID: "svc1", Name: "Test1"})
	r.RegisterService(&Info{ID: "svc2", Name: "Test2"})

	all := r.AllServices()
	if len(all) != 2 {
		t.Errorf("AllServices() returned %d services, want 2", len(all))
	}
}

func TestRegistryAllInstances(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	r.RegisterService(&Info{ID: "svc1", Name: "Test"})
	r.RegisterInstance(&Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"})
	r.RegisterInstance(&Instance{ID: "inst2", ServiceID: "svc1", Host: "localhost"})

	all := r.AllInstances()
	if len(all) != 2 {
		t.Errorf("AllInstances() returned %d instances, want 2", len(all))
	}
}

func TestRegistryServiceCount(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	if r.ServiceCount() != 0 {
		t.Error("ServiceCount should be 0 initially")
	}

	r.RegisterService(&Info{ID: "svc1", Name: "Test"})
	if r.ServiceCount() != 1 {
		t.Error("ServiceCount should be 1")
	}
}

func TestRegistryInstanceCount(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	if r.InstanceCount() != 0 {
		t.Error("InstanceCount should be 0 initially")
	}

	r.RegisterService(&Info{ID: "svc1", Name: "Test"})
	r.RegisterInstance(&Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"})
	if r.InstanceCount() != 1 {
		t.Error("InstanceCount should be 1")
	}
}

func TestRegistryServiceExists(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	if r.ServiceExists("svc1") {
		t.Error("ServiceExists should return false for unknown service")
	}

	r.RegisterService(&Info{ID: "svc1", Name: "Test"})
	if !r.ServiceExists("svc1") {
		t.Error("ServiceExists should return true for existing service")
	}
}

func TestRegistryInstanceExists(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	if r.InstanceExists("inst1") {
		t.Error("InstanceExists should return false for unknown instance")
	}

	r.RegisterService(&Info{ID: "svc1", Name: "Test"})
	r.RegisterInstance(&Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"})
	if !r.InstanceExists("inst1") {
		t.Error("InstanceExists should return true for existing instance")
	}
}

func TestRegistryUpdateServiceStatus(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	// Test not found
	if err := r.UpdateServiceStatus("unknown", StatusRunning); err != ErrServiceNotFound {
		t.Errorf("Expected ErrServiceNotFound, got %v", err)
	}

	r.RegisterService(&Info{ID: "svc1", Name: "Test", Status: StatusStarting})

	var statusChanged bool
	r.OnServiceStatusChange(func(svc *Info, oldStatus, newStatus Status) {
		statusChanged = true
	})

	if err := r.UpdateServiceStatus("svc1", StatusRunning); err != nil {
		t.Fatalf("UpdateServiceStatus() error: %v", err)
	}

	time.Sleep(10 * time.Millisecond)
	if !statusChanged {
		t.Error("OnServiceStatusChange callback not called")
	}

	svc, _ := r.GetService("svc1")
	if svc.Status != StatusRunning {
		t.Error("Status not updated")
	}
}

func TestRegistryUpdateInstanceStatus(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	// Test not found
	if err := r.UpdateInstanceStatus("unknown", StatusRunning); err != ErrInstanceNotFound {
		t.Errorf("Expected ErrInstanceNotFound, got %v", err)
	}

	r.RegisterService(&Info{ID: "svc1", Name: "Test"})
	r.RegisterInstance(&Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost", Status: StatusStarting})

	// Test status update to running (resets consecutive fails)
	if err := r.UpdateInstanceStatus("inst1", StatusRunning); err != nil {
		t.Fatalf("UpdateInstanceStatus() error: %v", err)
	}

	// Test status update to failed (increments consecutive fails)
	if err := r.UpdateInstanceStatus("inst1", StatusFailed); err != nil {
		t.Fatalf("UpdateInstanceStatus() error: %v", err)
	}

	inst, _ := r.GetInstance("inst1")
	if inst.ConsecutiveFails != 1 {
		t.Errorf("ConsecutiveFails = %d, want 1", inst.ConsecutiveFails)
	}
}

func TestRegistryUpdateInstanceHealth(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	// Test nil update
	if err := r.UpdateInstanceHealth(nil); err != ErrInvalidInstance {
		t.Errorf("Expected ErrInvalidInstance, got %v", err)
	}

	// Test missing instance ID
	if err := r.UpdateInstanceHealth(&HealthUpdate{}); err != ErrInvalidInstance {
		t.Errorf("Expected ErrInvalidInstance, got %v", err)
	}

	// Test not found
	if err := r.UpdateInstanceHealth(&HealthUpdate{InstanceID: "unknown"}); err != ErrInstanceNotFound {
		t.Errorf("Expected ErrInstanceNotFound, got %v", err)
	}

	r.RegisterService(&Info{ID: "svc1", Name: "Test"})
	r.RegisterInstance(&Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"})

	update := &HealthUpdate{
		InstanceID: "inst1",
		Status:     StatusRunning,
		CheckedAt:  time.Now(),
	}
	if err := r.UpdateInstanceHealth(update); err != nil {
		t.Fatalf("UpdateInstanceHealth() error: %v", err)
	}

	// Test failed status
	update.Status = StatusFailed
	if err := r.UpdateInstanceHealth(update); err != nil {
		t.Fatalf("UpdateInstanceHealth() error: %v", err)
	}
}

func TestRegistryRenewInstance(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	// Test not found
	if err := r.RenewInstance("unknown", time.Hour); err != ErrInstanceNotFound {
		t.Errorf("Expected ErrInstanceNotFound, got %v", err)
	}

	r.RegisterService(&Info{ID: "svc1", Name: "Test"})
	r.RegisterInstance(&Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"})

	if err := r.RenewInstance("inst1", time.Hour); err != nil {
		t.Fatalf("RenewInstance() error: %v", err)
	}

	inst, _ := r.GetInstance("inst1")
	if inst.ExpiresAt.IsZero() {
		t.Error("ExpiresAt not set after renewal")
	}
}

func TestRegistryShutdown(t *testing.T) {
	r := NewRegistry() // With expiration enabled

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := r.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown() error: %v", err)
	}
}

func TestRegistryShutdownTimeout(t *testing.T) {
	r := NewRegistry()

	// Create expired context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := r.Shutdown(ctx); err != context.Canceled {
		t.Errorf("Expected context.Canceled, got %v", err)
	}
}

func TestRegistryClear(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	r.RegisterService(&Info{ID: "svc1", Name: "Test"})
	r.RegisterInstance(&Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"})

	r.Clear()

	if r.ServiceCount() != 0 {
		t.Error("ServiceCount should be 0 after Clear")
	}
	if r.InstanceCount() != 0 {
		t.Error("InstanceCount should be 0 after Clear")
	}
}

func TestRegistryShutdownPreventsRegistration(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	r.cancel() // Simulate shutdown

	if err := r.RegisterService(&Info{ID: "svc1", Name: "Test"}); err != ErrRegistryShutdown {
		t.Errorf("Expected ErrRegistryShutdown, got %v", err)
	}

	// Need a valid service for instance registration
	r2 := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	r2.RegisterService(&Info{ID: "svc1", Name: "Test"})
	r2.cancel()

	if err := r2.RegisterInstance(&Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"}); err != ErrRegistryShutdown {
		t.Errorf("Expected ErrRegistryShutdown, got %v", err)
	}
}

// ============== Subjects Tests ==============

func TestDefaultSubjects(t *testing.T) {
	s := DefaultSubjects()
	if s.RegisterService != "services.register" {
		t.Error("RegisterService subject incorrect")
	}
	if s.DeregisterService != "services.deregister" {
		t.Error("DeregisterService subject incorrect")
	}
}

func TestTenantSubjects(t *testing.T) {
	s := TenantSubjects("tenant1")
	if s.RegisterService != "tenant1.services.register" {
		t.Error("RegisterService subject incorrect")
	}
	if s.HealthUpdate != "tenant1.health.update" {
		t.Error("HealthUpdate subject incorrect")
	}
}

// ============== Handler Tests ==============

func TestNewHandler(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer manager.Shutdown(context.Background())

	h := NewHandler(manager)
	if h == nil {
		t.Fatal("NewHandler returned nil")
	}
	if h.manager != manager {
		t.Error("Manager not set correctly")
	}
}

func TestNewHandlerWithOptions(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer manager.Shutdown(context.Background())

	pub := &mockPublisher{}
	subjects := TenantSubjects("tenant1")

	h := NewHandler(manager, WithPublisher(pub), WithSubjects(subjects))
	if h.publisher != pub {
		t.Error("Publisher not set correctly")
	}
	if h.subjects.RegisterService != subjects.RegisterService {
		t.Error("Subjects not set correctly")
	}
}

func TestNewHandlerWithPublisher(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer manager.Shutdown(context.Background())

	pub := &mockPublisher{}
	h := NewHandlerWithPublisher(manager, pub)
	if h.publisher != pub {
		t.Error("Publisher not set correctly")
	}

	// With custom subjects
	subjects := TenantSubjects("tenant1")
	h2 := NewHandlerWithPublisher(manager, pub, subjects)
	if h2.subjects.RegisterService != subjects.RegisterService {
		t.Error("Subjects not set correctly")
	}
}

func TestHandlerSetPublisher(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer manager.Shutdown(context.Background())

	h := NewHandler(manager)
	pub := &mockPublisher{}
	h.SetPublisher(pub)

	if h.publisher != pub {
		t.Error("SetPublisher did not set publisher correctly")
	}
}

// ============== Client Tests ==============

func TestNewClient(t *testing.T) {
	pub := &mockPublisher{}
	c := NewClient(pub)
	if c == nil {
		t.Fatal("NewClient returned nil")
	}
	if c.timeout != 10*time.Second {
		t.Error("Default timeout incorrect")
	}

	// With custom subjects
	subjects := TenantSubjects("tenant1")
	c2 := NewClient(pub, subjects)
	if c2.subjects.RegisterService != subjects.RegisterService {
		t.Error("Subjects not set correctly")
	}
}

func TestClientSetTimeout(t *testing.T) {
	pub := &mockPublisher{}
	c := NewClient(pub)
	c.SetTimeout(30 * time.Second)
	if c.timeout != 30*time.Second {
		t.Error("SetTimeout did not set timeout correctly")
	}
}

// ============== Manager Tests ==============

func TestNewManager(t *testing.T) {
	m, err := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	if err != nil {
		t.Fatalf("NewManager() error: %v", err)
	}
	defer m.Shutdown(context.Background())

	if m.registry == nil {
		t.Error("Registry not initialized")
	}
}

func TestNewManagerWithConfig(t *testing.T) {
	registry := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer registry.Shutdown(context.Background())

	pub := &mockPublisher{}
	cfg := ManagerConfig{
		Registry:        registry,
		Publisher:       pub,
		HealthCheck:     HealthCheckConfig{Enabled: false},
		ErrorBufferSize: 50,
	}

	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("NewManager() error: %v", err)
	}
	defer m.Shutdown(context.Background())

	if m.registry != registry {
		t.Error("Registry not set correctly")
	}
	if m.publisher != pub {
		t.Error("Publisher not set correctly")
	}
}

func TestDefaultHealthCheckConfig(t *testing.T) {
	cfg := DefaultHealthCheckConfig()
	if !cfg.Enabled {
		t.Error("Default Enabled should be true")
	}
	if cfg.Interval != 30*time.Second {
		t.Error("Default Interval incorrect")
	}
	if cfg.Timeout != 5*time.Second {
		t.Error("Default Timeout incorrect")
	}
}

func TestDefaultPublishSubjects(t *testing.T) {
	s := DefaultPublishSubjects()
	if s.ServiceRegistered != "services.registered" {
		t.Error("ServiceRegistered subject incorrect")
	}
}

func TestManagerRegisterService(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	ctx := context.Background()

	// Test invalid request
	if err := m.RegisterService(ctx, &RegistrationRequest{}); err == nil {
		t.Error("Expected error for invalid request")
	}

	// Test valid request
	req := &RegistrationRequest{Service: &Info{ID: "svc1", Name: "Test"}}
	if err := m.RegisterService(ctx, req); err != nil {
		t.Fatalf("RegisterService() error: %v", err)
	}

	// Verify registration
	svc, exists := m.GetService("svc1")
	if !exists || svc.ID != "svc1" {
		t.Error("Service not registered correctly")
	}
}

func TestManagerRegisterServiceDirect(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	if err := m.RegisterServiceDirect(context.Background(), &Info{ID: "svc1", Name: "Test"}); err != nil {
		t.Fatalf("RegisterServiceDirect() error: %v", err)
	}
}

func TestManagerDeregisterService(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: "svc1", Name: "Test"})

	// Test invalid request
	if err := m.DeregisterService(ctx, &DeregistrationRequest{}); err == nil {
		t.Error("Expected error for invalid request")
	}

	// Test valid request
	req := &DeregistrationRequest{ServiceID: "svc1", Reason: "shutdown"}
	if err := m.DeregisterService(ctx, req); err != nil {
		t.Fatalf("DeregisterService() error: %v", err)
	}
}

func TestManagerDeregisterServiceDirect(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	m.RegisterServiceDirect(context.Background(), &Info{ID: "svc1", Name: "Test"})

	if err := m.DeregisterServiceDirect(context.Background(), "svc1"); err != nil {
		t.Fatalf("DeregisterServiceDirect() error: %v", err)
	}
}

func TestManagerRegisterInstance(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: "svc1", Name: "Test"})

	// Test invalid request
	if err := m.RegisterInstance(ctx, &InstanceRegistrationRequest{}); err == nil {
		t.Error("Expected error for invalid request")
	}

	// Test valid request
	req := &InstanceRegistrationRequest{
		Instance: &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"},
	}
	if err := m.RegisterInstance(ctx, req); err != nil {
		t.Fatalf("RegisterInstance() error: %v", err)
	}

	// Test with TTL
	reqWithTTL := &InstanceRegistrationRequest{
		Instance: &Instance{ID: "inst2", ServiceID: "svc1", Host: "localhost"},
		TTL:      5 * time.Minute,
	}
	if err := m.RegisterInstance(ctx, reqWithTTL); err != nil {
		t.Fatalf("RegisterInstance() with TTL error: %v", err)
	}
}

func TestManagerRegisterInstanceDirect(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: "svc1", Name: "Test"})

	if err := m.RegisterInstanceDirect(ctx, &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"}); err != nil {
		t.Fatalf("RegisterInstanceDirect() error: %v", err)
	}
}

func TestManagerRegisterInstanceWithTTL(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: "svc1", Name: "Test"})

	if err := m.RegisterInstanceWithTTL(ctx, &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"}, 5*time.Minute); err != nil {
		t.Fatalf("RegisterInstanceWithTTL() error: %v", err)
	}
}

func TestManagerDeregisterInstance(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: "svc1", Name: "Test"})
	m.RegisterInstanceDirect(ctx, &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"})

	if err := m.DeregisterInstance(ctx, "inst1"); err != nil {
		t.Fatalf("DeregisterInstance() error: %v", err)
	}
}

func TestManagerGetServiceByName(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	m.RegisterServiceDirect(context.Background(), &Info{ID: "svc1", Name: "Test"})
	m.RegisterServiceDirect(context.Background(), &Info{ID: "svc2", Name: "Test"})

	services := m.GetServiceByName("Test")
	if len(services) != 2 {
		t.Errorf("GetServiceByName() returned %d services, want 2", len(services))
	}
}

func TestManagerGetInstance(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: "svc1", Name: "Test"})
	m.RegisterInstanceDirect(ctx, &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"})

	inst, exists := m.GetInstance("inst1")
	if !exists || inst.ID != "inst1" {
		t.Error("GetInstance did not return instance")
	}
}

func TestManagerGetInstances(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: "svc1", Name: "Test"})
	m.RegisterInstanceDirect(ctx, &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost", Status: StatusRunning})
	m.RegisterInstanceDirect(ctx, &Instance{ID: "inst2", ServiceID: "svc1", Host: "localhost", Status: StatusDegraded})

	instances := m.GetInstances("svc1")
	if len(instances) != 2 {
		t.Errorf("GetInstances() returned %d instances, want 2", len(instances))
	}
}

func TestManagerGetHealthyInstances(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: "svc1", Name: "Test"})
	m.RegisterInstanceDirect(ctx, &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost", Status: StatusRunning})
	m.RegisterInstanceDirect(ctx, &Instance{ID: "inst2", ServiceID: "svc1", Host: "localhost", Status: StatusDegraded})

	healthy := m.GetHealthyInstances("svc1")
	if len(healthy) != 1 {
		t.Errorf("GetHealthyInstances() returned %d instances, want 1", len(healthy))
	}
}

func TestManagerGetAvailableInstances(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: "svc1", Name: "Test"})
	m.RegisterInstanceDirect(ctx, &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost", Status: StatusRunning})
	m.RegisterInstanceDirect(ctx, &Instance{ID: "inst2", ServiceID: "svc1", Host: "localhost", Status: StatusDegraded})

	available := m.GetAvailableInstances("svc1")
	if len(available) != 2 {
		t.Errorf("GetAvailableInstances() returned %d instances, want 2", len(available))
	}
}

func TestManagerQueryServices(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	m.RegisterServiceDirect(context.Background(), &Info{ID: "svc1", Name: "Test", Status: StatusRunning})
	m.RegisterServiceDirect(context.Background(), &Info{ID: "svc2", Name: "Other", Status: StatusDegraded})

	services := m.QueryServices(&Query{OnlyHealthy: true})
	if len(services) != 1 {
		t.Errorf("QueryServices() returned %d services, want 1", len(services))
	}
}

func TestManagerDiscover(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{
		ID:     "svc1",
		Name:   "Test",
		Status: StatusRunning,
		Type:   ServiceTypeAPI,
	})
	m.RegisterInstanceDirect(ctx, &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost", Status: StatusRunning})

	// Basic discover
	services := m.Discover("Test")
	if len(services) != 1 {
		t.Errorf("Discover() returned %d services, want 1", len(services))
	}

	// With options
	services = m.Discover("Test", WithType(ServiceTypeAPI), OnlyHealthy())
	if len(services) != 1 {
		t.Errorf("Discover() with options returned %d services, want 1", len(services))
	}
}

func TestDiscoverOptions(t *testing.T) {
	q := &Query{}

	WithType(ServiceTypeAPI)(q)
	if len(q.Types) != 1 || q.Types[0] != ServiceTypeAPI {
		t.Error("WithType not working")
	}

	WithCapabilities("cap1", "cap2")(q)
	if len(q.Capabilities) != 2 {
		t.Error("WithCapabilities not working")
	}

	WithTags("tag1")(q)
	if len(q.Tags) != 1 {
		t.Error("WithTags not working")
	}

	WithVersion("1.0.0")(q)
	if q.Version != "1.0.0" {
		t.Error("WithVersion not working")
	}

	OnlyHealthy()(q)
	if !q.OnlyHealthy || q.OnlyAvailable {
		t.Error("OnlyHealthy not working")
	}

	WithMinInstances(3)(q)
	if q.MinInstances != 3 {
		t.Error("WithMinInstances not working")
	}

	WithLimit(10)(q)
	if q.Limit != 10 {
		t.Error("WithLimit not working")
	}
}

func TestManagerDiscoverInstances(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: "svc1", Name: "Test"})
	m.RegisterInstanceDirect(ctx, &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost", Status: StatusRunning})
	m.RegisterInstanceDirect(ctx, &Instance{ID: "inst2", ServiceID: "svc1", Host: "localhost", Status: StatusDegraded})

	healthy := m.DiscoverInstances("svc1", true)
	if len(healthy) != 1 {
		t.Errorf("DiscoverInstances(onlyHealthy=true) returned %d instances, want 1", len(healthy))
	}

	available := m.DiscoverInstances("svc1", false)
	if len(available) != 2 {
		t.Errorf("DiscoverInstances(onlyHealthy=false) returned %d instances, want 2", len(available))
	}
}

func TestManagerRenewInstance(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: "svc1", Name: "Test"})
	m.RegisterInstanceDirect(ctx, &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"})

	if err := m.RenewInstance(ctx, "inst1", time.Hour); err != nil {
		t.Fatalf("RenewInstance() error: %v", err)
	}
}

func TestManagerUpdateServiceStatus(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: "svc1", Name: "Test", Status: StatusStarting})

	if err := m.UpdateServiceStatus(ctx, "svc1", StatusRunning); err != nil {
		t.Fatalf("UpdateServiceStatus() error: %v", err)
	}
}

func TestManagerUpdateInstanceStatus(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: "svc1", Name: "Test"})
	m.RegisterInstanceDirect(ctx, &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost", Status: StatusStarting})

	if err := m.UpdateInstanceStatus(ctx, "inst1", StatusRunning); err != nil {
		t.Fatalf("UpdateInstanceStatus() error: %v", err)
	}
}

func TestManagerAllServices(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	m.RegisterServiceDirect(context.Background(), &Info{ID: "svc1", Name: "Test1"})
	m.RegisterServiceDirect(context.Background(), &Info{ID: "svc2", Name: "Test2"})

	all := m.AllServices()
	if len(all) != 2 {
		t.Errorf("AllServices() returned %d services, want 2", len(all))
	}
}

func TestManagerAllInstances(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: "svc1", Name: "Test"})
	m.RegisterInstanceDirect(ctx, &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"})
	m.RegisterInstanceDirect(ctx, &Instance{ID: "inst2", ServiceID: "svc1", Host: "localhost"})

	all := m.AllInstances()
	if len(all) != 2 {
		t.Errorf("AllInstances() returned %d instances, want 2", len(all))
	}
}

func TestManagerServiceCount(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	m.RegisterServiceDirect(context.Background(), &Info{ID: "svc1", Name: "Test"})
	if m.ServiceCount() != 1 {
		t.Error("ServiceCount incorrect")
	}
}

func TestManagerInstanceCount(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: "svc1", Name: "Test"})
	m.RegisterInstanceDirect(ctx, &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"})

	if m.InstanceCount() != 1 {
		t.Error("InstanceCount incorrect")
	}
}

func TestManagerRegistry(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	if m.Registry() == nil {
		t.Error("Registry() returned nil")
	}
}

func TestManagerSetPublisher(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	pub := &mockPublisher{}
	m.SetPublisher(pub)
	if m.publisher != pub {
		t.Error("SetPublisher did not set publisher correctly")
	}
}

func TestManagerErrors(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	errChan := m.Errors()
	if errChan == nil {
		t.Error("Errors() returned nil")
	}
}

func TestManagerOnError(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	var called bool
	m.OnError(func(err error) {
		called = true
	})

	// Trigger an error
	m.dispatchError(errors.New("test error"))

	if !called {
		t.Error("OnError callback not called")
	}
}

func TestManagerShutdown(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := m.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown() error: %v", err)
	}
}

// ============== Health Checker Tests ==============

func TestNewHealthChecker(t *testing.T) {
	hc := NewHealthChecker(HealthCheckerConfig{})
	if hc == nil {
		t.Fatal("NewHealthChecker returned nil")
	}
	if hc.timeout != 5*time.Second {
		t.Error("Default timeout incorrect")
	}
	if hc.healthyThreshold != 2 {
		t.Error("Default healthyThreshold incorrect")
	}
	if hc.unhealthyThreshold != 3 {
		t.Error("Default unhealthyThreshold incorrect")
	}
}

func TestHealthCheckerCheck(t *testing.T) {
	hc := NewHealthChecker(HealthCheckerConfig{
		Timeout:            time.Second,
		HealthyThreshold:   1,
		UnhealthyThreshold: 1,
	})

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Parse server URL
	host := server.Listener.Addr().String()

	inst := &Instance{
		ID:              "inst1",
		ServiceID:       "svc1",
		Host:            host,
		Protocol:        "http",
		HealthCheckPath: "/",
	}

	result := hc.Check(context.Background(), inst)
	if result.Status != StatusRunning {
		t.Errorf("Expected StatusRunning, got %v", result.Status)
	}
}

func TestHealthCheckerCheckFailure(t *testing.T) {
	hc := NewHealthChecker(HealthCheckerConfig{
		Timeout:            time.Second,
		HealthyThreshold:   1,
		UnhealthyThreshold: 1,
	})

	inst := &Instance{
		ID:              "inst1",
		ServiceID:       "svc1",
		Host:            "localhost:99999", // Invalid port
		Protocol:        "http",
		HealthCheckPath: "/health",
	}

	result := hc.Check(context.Background(), inst)
	if result.Status != StatusFailed {
		t.Errorf("Expected StatusFailed, got %v", result.Status)
	}
}

func TestHealthCheckerCheck503(t *testing.T) {
	hc := NewHealthChecker(HealthCheckerConfig{
		Timeout:            time.Second,
		HealthyThreshold:   1,
		UnhealthyThreshold: 1,
	})

	// Create a test server that returns 503
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	host := server.Listener.Addr().String()

	inst := &Instance{
		ID:              "inst1",
		ServiceID:       "svc1",
		Host:            host,
		Protocol:        "http",
		HealthCheckPath: "/",
	}

	result := hc.Check(context.Background(), inst)
	if result.Status != StatusMaintenance {
		t.Errorf("Expected StatusMaintenance, got %v", result.Status)
	}
}

func TestHealthCheckerCheckOtherError(t *testing.T) {
	hc := NewHealthChecker(HealthCheckerConfig{
		Timeout:            time.Second,
		HealthyThreshold:   1,
		UnhealthyThreshold: 1,
	})

	// Create a test server that returns 500
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	host := server.Listener.Addr().String()

	inst := &Instance{
		ID:              "inst1",
		ServiceID:       "svc1",
		Host:            host,
		Protocol:        "http",
		HealthCheckPath: "/",
	}

	result := hc.Check(context.Background(), inst)
	if result.Status != StatusFailed {
		t.Errorf("Expected StatusFailed, got %v", result.Status)
	}
}

func TestHealthCheckerDefaultProtocolAndPath(t *testing.T) {
	hc := NewHealthChecker(HealthCheckerConfig{
		Timeout:            time.Second,
		HealthyThreshold:   1,
		UnhealthyThreshold: 1,
	})

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	host := server.Listener.Addr().String()

	inst := &Instance{
		ID:        "inst1",
		ServiceID: "svc1",
		Host:      host,
		// No protocol or HealthCheckPath - should default
	}

	result := hc.Check(context.Background(), inst)
	if result.Status != StatusRunning {
		t.Errorf("Expected StatusRunning (defaults applied), got %v: %s", result.Status, result.Message)
	}
}

// ============== HTTP Handler Tests ==============

func TestManagerHTTPHandleServiceRegistration(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	req := RegistrationRequest{
		Service: &Info{ID: "svc1", Name: "Test"},
	}
	body, _ := json.Marshal(req)

	// Test POST
	httpReq := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	m.HandleServiceRegistration(rr, httpReq)
	if rr.Code != http.StatusCreated {
		t.Errorf("HandleServiceRegistration returned %d, want %d", rr.Code, http.StatusCreated)
	}

	// Test wrong method
	httpReq = httptest.NewRequest(http.MethodGet, "/register", nil)
	rr = httptest.NewRecorder()

	m.HandleServiceRegistration(rr, httpReq)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("HandleServiceRegistration returned %d for GET, want %d", rr.Code, http.StatusMethodNotAllowed)
	}

	// Test invalid JSON
	httpReq = httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte("invalid")))
	rr = httptest.NewRecorder()

	m.HandleServiceRegistration(rr, httpReq)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("HandleServiceRegistration returned %d for invalid JSON, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestManagerHTTPHandleServiceDeregistration(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	m.RegisterServiceDirect(context.Background(), &Info{ID: "svc1", Name: "Test"})

	req := DeregistrationRequest{ServiceID: "svc1"}
	body, _ := json.Marshal(req)

	// Test DELETE
	httpReq := httptest.NewRequest(http.MethodDelete, "/deregister", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	m.HandleServiceDeregistration(rr, httpReq)
	if rr.Code != http.StatusOK {
		t.Errorf("HandleServiceDeregistration returned %d, want %d", rr.Code, http.StatusOK)
	}

	// Test not found
	httpReq = httptest.NewRequest(http.MethodDelete, "/deregister", bytes.NewReader(body))
	rr = httptest.NewRecorder()

	m.HandleServiceDeregistration(rr, httpReq)
	if rr.Code != http.StatusNotFound {
		t.Errorf("HandleServiceDeregistration returned %d for not found, want %d", rr.Code, http.StatusNotFound)
	}

	// Test wrong method
	httpReq = httptest.NewRequest(http.MethodGet, "/deregister", nil)
	rr = httptest.NewRecorder()

	m.HandleServiceDeregistration(rr, httpReq)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("HandleServiceDeregistration returned %d for GET, want %d", rr.Code, http.StatusMethodNotAllowed)
	}
}

func TestManagerHTTPHandleInstanceRegistration(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	m.RegisterServiceDirect(context.Background(), &Info{ID: "svc1", Name: "Test"})

	req := InstanceRegistrationRequest{
		Instance: &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"},
	}
	body, _ := json.Marshal(req)

	// Test POST
	httpReq := httptest.NewRequest(http.MethodPost, "/register-instance", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	m.HandleInstanceRegistration(rr, httpReq)
	if rr.Code != http.StatusCreated {
		t.Errorf("HandleInstanceRegistration returned %d, want %d", rr.Code, http.StatusCreated)
	}

	// Test wrong method
	httpReq = httptest.NewRequest(http.MethodGet, "/register-instance", nil)
	rr = httptest.NewRecorder()

	m.HandleInstanceRegistration(rr, httpReq)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("HandleInstanceRegistration returned %d for GET, want %d", rr.Code, http.StatusMethodNotAllowed)
	}
}

func TestManagerHTTPHandleServiceQuery(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	m.RegisterServiceDirect(context.Background(), &Info{ID: "svc1", Name: "Test", Status: StatusRunning})

	// Test GET with query params
	httpReq := httptest.NewRequest(http.MethodGet, "/services?name=Test&healthy=true&available=true", nil)
	rr := httptest.NewRecorder()

	m.HandleServiceQuery(rr, httpReq)
	if rr.Code != http.StatusOK {
		t.Errorf("HandleServiceQuery returned %d, want %d", rr.Code, http.StatusOK)
	}

	// Test POST with JSON body
	query := Query{OnlyHealthy: true}
	body, _ := json.Marshal(query)
	httpReq = httptest.NewRequest(http.MethodPost, "/services", bytes.NewReader(body))
	rr = httptest.NewRecorder()

	m.HandleServiceQuery(rr, httpReq)
	if rr.Code != http.StatusOK {
		t.Errorf("HandleServiceQuery POST returned %d, want %d", rr.Code, http.StatusOK)
	}

	// Test wrong method
	httpReq = httptest.NewRequest(http.MethodPut, "/services", nil)
	rr = httptest.NewRecorder()

	m.HandleServiceQuery(rr, httpReq)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("HandleServiceQuery returned %d for PUT, want %d", rr.Code, http.StatusMethodNotAllowed)
	}
}

func TestManagerHTTPHandleInstanceRenew(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: "svc1", Name: "Test"})
	m.RegisterInstanceDirect(ctx, &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"})

	req := struct {
		InstanceID string        `json:"instance_id"`
		TTL        time.Duration `json:"ttl"`
	}{
		InstanceID: "inst1",
		TTL:        5 * time.Minute,
	}
	body, _ := json.Marshal(req)

	// Test POST
	httpReq := httptest.NewRequest(http.MethodPost, "/renew", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	m.HandleInstanceRenew(rr, httpReq)
	if rr.Code != http.StatusOK {
		t.Errorf("HandleInstanceRenew returned %d, want %d", rr.Code, http.StatusOK)
	}

	// Test not found
	req.InstanceID = "unknown"
	body, _ = json.Marshal(req)
	httpReq = httptest.NewRequest(http.MethodPost, "/renew", bytes.NewReader(body))
	rr = httptest.NewRecorder()

	m.HandleInstanceRenew(rr, httpReq)
	if rr.Code != http.StatusNotFound {
		t.Errorf("HandleInstanceRenew returned %d for not found, want %d", rr.Code, http.StatusNotFound)
	}

	// Test wrong method
	httpReq = httptest.NewRequest(http.MethodGet, "/renew", nil)
	rr = httptest.NewRecorder()

	m.HandleInstanceRenew(rr, httpReq)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("HandleInstanceRenew returned %d for GET, want %d", rr.Code, http.StatusMethodNotAllowed)
	}
}

// ============== Helper Tests ==============

func TestItoa(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{-1, "-1"},
		{123, "123"},
		{-456, "-456"},
		{1000000, "1000000"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := itoa(tt.input); got != tt.expected {
				t.Errorf("itoa(%d) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// ============== Mock Publisher ==============

type mockPublisher struct {
	mu           sync.Mutex
	publishErr   error
	requestReply *core.Message
	requestErr   error
	published    []struct {
		subject string
		msg     *core.Message
	}
}

func (m *mockPublisher) Publish(ctx context.Context, subject string, msg *core.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.published = append(m.published, struct {
		subject string
		msg     *core.Message
	}{subject, msg})
	return m.publishErr
}

func (m *mockPublisher) PublishAsync(ctx context.Context, subject string, msg *core.Message) core.PubAckFuture {
	return nil
}

func (m *mockPublisher) Request(ctx context.Context, subject string, msg *core.Message) (*core.Message, error) {
	if m.requestErr != nil {
		return nil, m.requestErr
	}
	return m.requestReply, nil
}

func (m *mockPublisher) Close(ctx context.Context) error {
	return nil
}

// ============== Handler Message Tests ==============

func TestHandlerHandleRegisterService(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	pub := &mockPublisher{}
	h := NewHandler(m, WithPublisher(pub))

	// Test valid registration
	req := &RegistrationRequest{Service: &Info{ID: "svc1", Name: "Test"}}
	data, _ := json.Marshal(req)
	msg := core.NewMessage(data)
	msg.Reply = "reply-subject"

	err := h.HandleRegisterService(context.Background(), msg)
	if err != nil {
		t.Errorf("HandleRegisterService() error: %v", err)
	}

	// Verify service was registered
	if !m.registry.ServiceExists("svc1") {
		t.Error("Service not registered")
	}

	// Test invalid JSON
	badMsg := core.NewMessage([]byte("invalid"))
	badMsg.Reply = "reply"
	err = h.HandleRegisterService(context.Background(), badMsg)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}

	// Test without reply subject (no response expected)
	req2 := &RegistrationRequest{Service: &Info{ID: "svc2", Name: "Test2"}}
	data2, _ := json.Marshal(req2)
	msg2 := core.NewMessage(data2)
	// No Reply set
	err = h.HandleRegisterService(context.Background(), msg2)
	if err != nil {
		t.Errorf("HandleRegisterService() without reply error: %v", err)
	}
}

func TestHandlerHandleDeregisterService(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	m.RegisterServiceDirect(context.Background(), &Info{ID: "svc1", Name: "Test"})

	h := NewHandler(m, WithPublisher(&mockPublisher{}))

	req := &DeregistrationRequest{ServiceID: "svc1"}
	data, _ := json.Marshal(req)
	msg := core.NewMessage(data)
	msg.Reply = "reply"

	err := h.HandleDeregisterService(context.Background(), msg)
	if err != nil {
		t.Errorf("HandleDeregisterService() error: %v", err)
	}

	// Test invalid JSON
	badMsg := core.NewMessage([]byte("invalid"))
	badMsg.Reply = "reply"
	err = h.HandleDeregisterService(context.Background(), badMsg)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestHandlerHandleQueryServices(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	m.RegisterServiceDirect(context.Background(), &Info{ID: "svc1", Name: "Test", Status: StatusRunning})

	h := NewHandler(m, WithPublisher(&mockPublisher{}))

	query := &Query{OnlyHealthy: true}
	data, _ := json.Marshal(query)
	msg := core.NewMessage(data)
	msg.Reply = "reply"

	err := h.HandleQueryServices(context.Background(), msg)
	if err != nil {
		t.Errorf("HandleQueryServices() error: %v", err)
	}

	// Test with tenant context
	ctx := core.WithTenantID(context.Background(), "tenant1")
	err = h.HandleQueryServices(ctx, msg)
	if err != nil {
		t.Errorf("HandleQueryServices() with tenant error: %v", err)
	}
}

func TestHandlerHandleDiscoverServices(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	m.RegisterServiceDirect(context.Background(), &Info{
		ID:           "svc1",
		Name:         "Test",
		Status:       StatusRunning,
		Type:         ServiceTypeAPI,
		Version:      "1.0.0",
		Capabilities: []string{"cap1"},
		Tags:         []string{"tag1"},
	})
	m.RegisterInstanceDirect(context.Background(), &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"})

	h := NewHandler(m, WithPublisher(&mockPublisher{}))

	req := &DiscoverRequest{
		Name:         "Test",
		Type:         ServiceTypeAPI,
		Version:      "1.0.0",
		Capabilities: []string{"cap1"},
		Tags:         []string{"tag1"},
		OnlyHealthy:  true,
		TenantID:     "tenant1",
	}
	data, _ := json.Marshal(req)
	msg := core.NewMessage(data)
	msg.Reply = "reply"

	err := h.HandleDiscoverServices(context.Background(), msg)
	if err != nil {
		t.Errorf("HandleDiscoverServices() error: %v", err)
	}

	// Test with tenant context
	ctx := core.WithTenantID(context.Background(), "tenant2")
	err = h.HandleDiscoverServices(ctx, msg)
	if err != nil {
		t.Errorf("HandleDiscoverServices() with tenant error: %v", err)
	}
}

func TestHandlerHandleRegisterInstance(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	m.RegisterServiceDirect(context.Background(), &Info{ID: "svc1", Name: "Test"})

	h := NewHandler(m, WithPublisher(&mockPublisher{}))

	req := &InstanceRegistrationRequest{
		Instance: &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"},
		TTL:      5 * time.Minute,
	}
	data, _ := json.Marshal(req)
	msg := core.NewMessage(data)
	msg.Reply = "reply"

	err := h.HandleRegisterInstance(context.Background(), msg)
	if err != nil {
		t.Errorf("HandleRegisterInstance() error: %v", err)
	}

	// Test invalid JSON
	badMsg := core.NewMessage([]byte("invalid"))
	badMsg.Reply = "reply"
	err = h.HandleRegisterInstance(context.Background(), badMsg)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestHandlerHandleDeregisterInstance(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	m.RegisterServiceDirect(context.Background(), &Info{ID: "svc1", Name: "Test"})
	m.RegisterInstanceDirect(context.Background(), &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"})

	h := NewHandler(m, WithPublisher(&mockPublisher{}))

	req := &InstanceDeregistrationRequest{InstanceID: "inst1"}
	data, _ := json.Marshal(req)
	msg := core.NewMessage(data)
	msg.Reply = "reply"

	err := h.HandleDeregisterInstance(context.Background(), msg)
	if err != nil {
		t.Errorf("HandleDeregisterInstance() error: %v", err)
	}

	// Test invalid JSON
	badMsg := core.NewMessage([]byte("invalid"))
	badMsg.Reply = "reply"
	err = h.HandleDeregisterInstance(context.Background(), badMsg)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestHandlerHandleRenewInstance(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	m.RegisterServiceDirect(context.Background(), &Info{ID: "svc1", Name: "Test"})
	m.RegisterInstanceDirect(context.Background(), &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"})

	h := NewHandler(m, WithPublisher(&mockPublisher{}))

	req := &RenewInstanceRequest{InstanceID: "inst1", TTL: 5 * time.Minute}
	data, _ := json.Marshal(req)
	msg := core.NewMessage(data)
	msg.Reply = "reply"

	err := h.HandleRenewInstance(context.Background(), msg)
	if err != nil {
		t.Errorf("HandleRenewInstance() error: %v", err)
	}

	// Test invalid JSON
	badMsg := core.NewMessage([]byte("invalid"))
	badMsg.Reply = "reply"
	err = h.HandleRenewInstance(context.Background(), badMsg)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestHandlerHandleQueryInstances(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	m.RegisterServiceDirect(context.Background(), &Info{ID: "svc1", Name: "Test"})
	m.RegisterInstanceDirect(context.Background(), &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost", Status: StatusRunning})
	m.RegisterInstanceDirect(context.Background(), &Instance{ID: "inst2", ServiceID: "svc1", Host: "localhost", Status: StatusDegraded})

	h := NewHandler(m, WithPublisher(&mockPublisher{}))

	// Test only healthy
	req := &QueryInstancesRequest{ServiceID: "svc1", OnlyHealthy: true}
	data, _ := json.Marshal(req)
	msg := core.NewMessage(data)
	msg.Reply = "reply"

	err := h.HandleQueryInstances(context.Background(), msg)
	if err != nil {
		t.Errorf("HandleQueryInstances() error: %v", err)
	}

	// Test only available
	req2 := &QueryInstancesRequest{ServiceID: "svc1", OnlyAvailable: true}
	data2, _ := json.Marshal(req2)
	msg2 := core.NewMessage(data2)
	msg2.Reply = "reply"

	err = h.HandleQueryInstances(context.Background(), msg2)
	if err != nil {
		t.Errorf("HandleQueryInstances() only available error: %v", err)
	}

	// Test all
	req3 := &QueryInstancesRequest{ServiceID: "svc1"}
	data3, _ := json.Marshal(req3)
	msg3 := core.NewMessage(data3)
	msg3.Reply = "reply"

	err = h.HandleQueryInstances(context.Background(), msg3)
	if err != nil {
		t.Errorf("HandleQueryInstances() all error: %v", err)
	}

	// Test invalid JSON
	badMsg := core.NewMessage([]byte("invalid"))
	badMsg.Reply = "reply"
	err = h.HandleQueryInstances(context.Background(), badMsg)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestHandlerHandleHealthUpdate(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	m.RegisterServiceDirect(context.Background(), &Info{ID: "svc1", Name: "Test", Status: StatusStarting})
	m.RegisterInstanceDirect(context.Background(), &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"})

	h := NewHandler(m, WithPublisher(&mockPublisher{}))

	// Test instance health update
	update := &HealthUpdate{
		ServiceID:  "svc1",
		InstanceID: "inst1",
		Status:     StatusRunning,
	}
	data, _ := json.Marshal(update)
	msg := core.NewMessage(data)
	msg.Reply = "reply"

	err := h.HandleHealthUpdate(context.Background(), msg)
	if err != nil {
		t.Errorf("HandleHealthUpdate() error: %v", err)
	}

	// Test service health update (without instance ID)
	update2 := &HealthUpdate{
		ServiceID: "svc1",
		Status:    StatusRunning,
	}
	data2, _ := json.Marshal(update2)
	msg2 := core.NewMessage(data2)
	msg2.Reply = "reply"

	err = h.HandleHealthUpdate(context.Background(), msg2)
	if err != nil {
		t.Errorf("HandleHealthUpdate() service error: %v", err)
	}

	// Test missing service ID
	update3 := &HealthUpdate{Status: StatusRunning}
	data3, _ := json.Marshal(update3)
	msg3 := core.NewMessage(data3)
	msg3.Reply = "reply"

	err = h.HandleHealthUpdate(context.Background(), msg3)
	if err == nil {
		t.Error("Expected error for missing service ID")
	}

	// Test invalid JSON
	badMsg := core.NewMessage([]byte("invalid"))
	badMsg.Reply = "reply"
	err = h.HandleHealthUpdate(context.Background(), badMsg)
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
}

func TestHandlerWithoutPublisher(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	h := NewHandler(m) // No publisher set

	req := &RegistrationRequest{Service: &Info{ID: "svc1", Name: "Test"}}
	data, _ := json.Marshal(req)
	msg := core.NewMessage(data)
	msg.Reply = "reply" // Expects reply but no publisher

	err := h.HandleRegisterService(context.Background(), msg)
	if err != nil {
		t.Errorf("HandleRegisterService() should succeed even without publisher: %v", err)
	}
}

func TestHandlerSendResponseErrors(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	// Test with publisher that returns error
	pub := &mockPublisher{publishErr: errors.New("publish failed")}
	h := NewHandler(m, WithPublisher(pub))

	req := &RegistrationRequest{Service: &Info{ID: "svc1", Name: "Test"}}
	data, _ := json.Marshal(req)
	msg := core.NewMessage(data)
	msg.Reply = "reply"

	err := h.HandleRegisterService(context.Background(), msg)
	if err == nil {
		t.Error("Expected error when publisher fails")
	}
}

// ============== Client Method Tests ==============

func TestClientRegisterService(t *testing.T) {
	respData, _ := json.Marshal(&ServiceRegistrationResponse{Success: true, ServiceID: "svc1"})
	pub := &mockPublisher{requestReply: core.NewMessage(respData)}
	c := NewClient(pub)

	service := &Info{ID: "svc1", Name: "Test"}
	resp, err := c.RegisterService(context.Background(), service)
	if err != nil {
		t.Fatalf("RegisterService() error: %v", err)
	}
	if !resp.Success || resp.ServiceID != "svc1" {
		t.Error("Unexpected response")
	}

	// Test error
	pub.requestErr = errors.New("request failed")
	_, err = c.RegisterService(context.Background(), service)
	if err == nil {
		t.Error("Expected error")
	}
}

func TestClientDeregisterService(t *testing.T) {
	respData, _ := json.Marshal(&ServiceDeregistrationResponse{Success: true, ServiceID: "svc1"})
	pub := &mockPublisher{requestReply: core.NewMessage(respData)}
	c := NewClient(pub)

	resp, err := c.DeregisterService(context.Background(), "svc1", "shutdown")
	if err != nil {
		t.Fatalf("DeregisterService() error: %v", err)
	}
	if !resp.Success {
		t.Error("Unexpected response")
	}
}

func TestClientQueryServices(t *testing.T) {
	respData, _ := json.Marshal(&QueryServicesResponse{Services: []*Info{{ID: "svc1"}}, Count: 1})
	pub := &mockPublisher{requestReply: core.NewMessage(respData)}
	c := NewClient(pub)

	resp, err := c.QueryServices(context.Background(), &Query{OnlyHealthy: true})
	if err != nil {
		t.Fatalf("QueryServices() error: %v", err)
	}
	if resp.Count != 1 {
		t.Error("Unexpected response")
	}
}

func TestClientDiscover(t *testing.T) {
	respData, _ := json.Marshal(&DiscoverResponse{Services: []*Info{{ID: "svc1"}}, Count: 1})
	pub := &mockPublisher{requestReply: core.NewMessage(respData)}
	c := NewClient(pub)

	resp, err := c.Discover(context.Background(), "Test")
	if err != nil {
		t.Fatalf("Discover() error: %v", err)
	}
	if resp.Count != 1 {
		t.Error("Unexpected response")
	}

	// Test with options
	resp, err = c.Discover(context.Background(), "Test", WithType(ServiceTypeAPI), WithVersion("1.0.0"))
	if err != nil {
		t.Fatalf("Discover() with options error: %v", err)
	}
	_ = resp // use resp to suppress warning
}

func TestClientRegisterInstance(t *testing.T) {
	respData, _ := json.Marshal(&InstanceRegistrationResponse{Success: true, InstanceID: "inst1", ServiceID: "svc1"})
	pub := &mockPublisher{requestReply: core.NewMessage(respData)}
	c := NewClient(pub)

	inst := &Instance{ID: "inst1", ServiceID: "svc1", Host: "localhost"}
	resp, err := c.RegisterInstance(context.Background(), inst, 5*time.Minute)
	if err != nil {
		t.Fatalf("RegisterInstance() error: %v", err)
	}
	if !resp.Success {
		t.Error("Unexpected response")
	}
}

func TestClientDeregisterInstance(t *testing.T) {
	respData, _ := json.Marshal(&InstanceDeregistrationResponse{Success: true, InstanceID: "inst1"})
	pub := &mockPublisher{requestReply: core.NewMessage(respData)}
	c := NewClient(pub)

	resp, err := c.DeregisterInstance(context.Background(), "inst1", "shutdown")
	if err != nil {
		t.Fatalf("DeregisterInstance() error: %v", err)
	}
	if !resp.Success {
		t.Error("Unexpected response")
	}
}

func TestClientRenewInstance(t *testing.T) {
	respData, _ := json.Marshal(&RenewInstanceResponse{Success: true, InstanceID: "inst1"})
	pub := &mockPublisher{requestReply: core.NewMessage(respData)}
	c := NewClient(pub)

	resp, err := c.RenewInstance(context.Background(), "inst1", 5*time.Minute)
	if err != nil {
		t.Fatalf("RenewInstance() error: %v", err)
	}
	if !resp.Success {
		t.Error("Unexpected response")
	}
}

func TestClientQueryInstances(t *testing.T) {
	respData, _ := json.Marshal(&QueryInstancesResponse{Instances: []*Instance{{ID: "inst1"}}, Count: 1, ServiceID: "svc1"})
	pub := &mockPublisher{requestReply: core.NewMessage(respData)}
	c := NewClient(pub)

	resp, err := c.QueryInstances(context.Background(), "svc1", true)
	if err != nil {
		t.Fatalf("QueryInstances() error: %v", err)
	}
	if resp.Count != 1 {
		t.Error("Unexpected response")
	}
}

func TestClientUpdateHealth(t *testing.T) {
	pub := &mockPublisher{}
	c := NewClient(pub)

	update := &HealthUpdate{ServiceID: "svc1", InstanceID: "inst1", Status: StatusRunning}
	err := c.UpdateHealth(context.Background(), update)
	if err != nil {
		t.Fatalf("UpdateHealth() error: %v", err)
	}

	// Verify publish was called
	if len(pub.published) != 1 {
		t.Error("UpdateHealth did not publish")
	}
}

// ============== Manager Event Publishing Tests ==============

func TestManagerPublishEvents(t *testing.T) {
	pub := &mockPublisher{}
	m, _ := NewManager(ManagerConfig{
		Publisher:   pub,
		HealthCheck: HealthCheckConfig{Enabled: false},
	})
	defer m.Shutdown(context.Background())

	// Register service - triggers ServiceRegistered event
	m.RegisterServiceDirect(context.Background(), &Info{ID: "svc1", Name: "Test", Status: StatusRunning})
	time.Sleep(50 * time.Millisecond)

	pub.mu.Lock()
	initialPubs := len(pub.published)
	pub.mu.Unlock()

	if initialPubs == 0 {
		t.Error("Expected service registered event to be published")
	}

	// Update service - triggers ServiceUpdated event
	m.RegisterServiceDirect(context.Background(), &Info{ID: "svc1", Name: "Test", Status: StatusDegraded})
	time.Sleep(50 * time.Millisecond)

	// Register instance - triggers InstanceRegistered event
	m.RegisterInstanceDirect(context.Background(), &Instance{
		ID:              "inst1",
		ServiceID:       "svc1",
		Host:            "localhost",
		Status:          StatusRunning,
		HealthCheckPath: "/health",
	})
	time.Sleep(50 * time.Millisecond)

	// Deregister instance - triggers InstanceDeregistered event
	m.DeregisterInstance(context.Background(), "inst1")
	time.Sleep(50 * time.Millisecond)

	// Deregister service - triggers ServiceDeregistered event
	m.DeregisterServiceDirect(context.Background(), "svc1")
	time.Sleep(50 * time.Millisecond)

	pub.mu.Lock()
	finalPubs := len(pub.published)
	pub.mu.Unlock()

	if finalPubs <= initialPubs {
		t.Error("Expected events to be published during lifecycle")
	}
}

func TestManagerHealthCheckEnabled(t *testing.T) {
	// Create manager with health check enabled (short interval for testing)
	cfg := ManagerConfig{
		HealthCheck: HealthCheckConfig{
			Enabled:      true,
			Interval:     50 * time.Millisecond,
			Timeout:      time.Second,
			InitialDelay: 10 * time.Millisecond,
		},
	}

	m, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("NewManager() error: %v", err)
	}

	// Register service and instance with health check
	m.RegisterServiceDirect(context.Background(), &Info{ID: "svc1", Name: "Test"})
	m.RegisterInstanceDirect(context.Background(), &Instance{
		ID:              "inst1",
		ServiceID:       "svc1",
		Host:            "localhost:99999", // Invalid port, will fail health check
		HealthCheckPath: "/health",
		Status:          StatusRunning,
	})

	// Wait for health check to run
	time.Sleep(100 * time.Millisecond)

	// Shutdown
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := m.Shutdown(ctx); err != nil {
		t.Fatalf("Shutdown() error: %v", err)
	}
}

// ============== Cleanup Loop Test ==============

func TestRegistryCleanupLoop(t *testing.T) {
	r := NewRegistry(RegistryConfig{
		CleanupInterval:   50 * time.Millisecond,
		ExpirationEnabled: true,
	})

	r.RegisterService(&Info{ID: "svc1", Name: "Test"})

	// Register an instance that will expire
	inst := &Instance{
		ID:        "inst1",
		ServiceID: "svc1",
		Host:      "localhost",
		ExpiresAt: time.Now().Add(20 * time.Millisecond),
	}
	r.RegisterInstance(inst)

	// Wait for instance to expire and cleanup to run
	time.Sleep(100 * time.Millisecond)

	// Check instance was cleaned up
	_, exists := r.GetInstance("inst1")
	if exists {
		t.Error("Expired instance should have been cleaned up")
	}

	r.Shutdown(context.Background())
}

// ============== Index Management Tests ==============

func TestRegistryIndices(t *testing.T) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	// Register service with various indexed fields
	svc := &Info{
		ID:           "svc1",
		Name:         "TestService",
		Type:         ServiceTypeAPI,
		Capabilities: []string{"cap1", "cap2"},
		Tags:         []string{"tag1", "tag2"},
	}
	r.RegisterService(svc)

	// Update service (tests index removal and re-adding)
	svc.Capabilities = []string{"cap3"}
	svc.Tags = []string{"tag3"}
	r.RegisterService(svc)

	// Deregister (tests index cleanup)
	r.DeregisterService("svc1")

	// Verify indices are cleaned up
	if len(r.byName) > 0 {
		t.Error("Name index not cleaned up")
	}
	if len(r.byType) > 0 {
		t.Error("Type index not cleaned up")
	}
}

// ============== Mock Subscriber for Subscription Tests ==============

type mockSubscriber struct {
	subscribeErr error
	subs         []*mockSubscription
}

type mockSubscription struct {
	subject string
	valid   bool
}

func (m *mockSubscription) Subject() string { return m.subject }
func (m *mockSubscription) Unsubscribe() error {
	m.valid = false
	return nil
}
func (m *mockSubscription) Drain() error   { return nil }
func (m *mockSubscription) IsValid() bool { return m.valid }

func (m *mockSubscriber) Subscribe(ctx context.Context, subject string, handler core.MessageHandler) (core.Subscription, error) {
	if m.subscribeErr != nil {
		return nil, m.subscribeErr
	}
	sub := &mockSubscription{subject: subject, valid: true}
	m.subs = append(m.subs, sub)
	return sub, nil
}

func (m *mockSubscriber) QueueSubscribe(ctx context.Context, subject string, queue string, handler core.MessageHandler) (core.Subscription, error) {
	if m.subscribeErr != nil {
		return nil, m.subscribeErr
	}
	sub := &mockSubscription{subject: subject, valid: true}
	m.subs = append(m.subs, sub)
	return sub, nil
}

func (m *mockSubscriber) Close(ctx context.Context) error { return nil }

func TestHandlerRegisterSubscriptions(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	h := NewHandler(m)
	sub := &mockSubscriber{}

	subs, err := h.RegisterSubscriptions(context.Background(), sub)
	if err != nil {
		t.Fatalf("RegisterSubscriptions() error: %v", err)
	}
	if len(subs) != 9 { // 9 handlers
		t.Errorf("Expected 9 subscriptions, got %d", len(subs))
	}

	// Test with subscribe error
	sub2 := &mockSubscriber{subscribeErr: errors.New("subscribe failed")}
	_, err = h.RegisterSubscriptions(context.Background(), sub2)
	if err == nil {
		t.Error("Expected error when subscribe fails")
	}
}

func TestHandlerRegisterQueueSubscriptions(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	h := NewHandler(m)
	sub := &mockSubscriber{}

	subs, err := h.RegisterQueueSubscriptions(context.Background(), sub, "worker-queue")
	if err != nil {
		t.Fatalf("RegisterQueueSubscriptions() error: %v", err)
	}
	if len(subs) != 8 { // 8 queue handlers (health handler not included in queue subscriptions)
		t.Errorf("Expected 8 subscriptions, got %d", len(subs))
	}
}

func TestHandlerUnsubscribeAll(t *testing.T) {
	m, _ := NewManager(ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}})
	defer m.Shutdown(context.Background())

	h := NewHandler(m)

	subs := []core.Subscription{
		&mockSubscription{valid: true},
		nil, // Test nil handling
		&mockSubscription{valid: true},
	}

	h.unsubscribeAll(subs)

	// Verify subscriptions are unsubscribed
	if subs[0].(*mockSubscription).valid {
		t.Error("First subscription not unsubscribed")
	}
	if subs[2].(*mockSubscription).valid {
		t.Error("Third subscription not unsubscribed")
	}
}

// ============== Benchmarks ==============

func BenchmarkRegistryRegisterService(b *testing.B) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r.RegisterService(&Info{ID: "svc1", Name: "Test", Status: StatusRunning})
	}
}

func BenchmarkRegistryGetService(b *testing.B) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	r.RegisterService(&Info{ID: "svc1", Name: "Test"})

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r.GetService("svc1")
	}
}

func BenchmarkRegistryQueryServices(b *testing.B) {
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})
	defer r.Shutdown(context.Background())

	for i := 0; i < 100; i++ {
		r.RegisterService(&Info{
			ID:     "svc" + itoa(i),
			Name:   "Test",
			Status: StatusRunning,
			Tags:   []string{"prod"},
		})
	}

	q := &Query{OnlyHealthy: true, Tags: []string{"prod"}}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r.QueryServices(q)
	}
}

func BenchmarkQueryMatches(b *testing.B) {
	service := &Info{
		ID:           "svc1",
		Name:         "MyService",
		Type:         ServiceTypeAPI,
		Status:       StatusRunning,
		Capabilities: []string{"read", "write"},
		Tags:         []string{"production"},
		Metadata:     map[string]string{"env": "prod"},
		Instances:    []*Instance{{ID: "inst1"}},
	}

	q := &Query{
		OnlyHealthy:  true,
		Capabilities: []string{"read"},
		Tags:         []string{"production"},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		q.Matches(service)
	}
}

func BenchmarkInfoClone(b *testing.B) {
	info := &Info{
		ID:           "svc1",
		Name:         "MyService",
		Dependencies: []string{"dep1", "dep2", "dep3"},
		Capabilities: []string{"cap1", "cap2"},
		Tags:         []string{"tag1", "tag2", "tag3"},
		Metadata:     map[string]string{"key1": "value1", "key2": "value2"},
		Instances: []*Instance{
			{ID: "inst1", Host: "localhost"},
			{ID: "inst2", Host: "localhost"},
		},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		info.Clone()
	}
}

func BenchmarkStatusString(b *testing.B) {
	var s string
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		s = StatusRunning.String()
	}
	_ = s // prevent compiler optimization
}

// ============== Helper functions ==============

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
