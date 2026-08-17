package endpoint

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

// ============== Status Tests ==============

func TestStatusString(t *testing.T) {
	tests := []struct {
		status   Status
		expected string
	}{
		{StatusUnknown, "unknown"},
		{StatusHealthy, "healthy"},
		{StatusUnhealthy, "unhealthy"},
		{StatusDegraded, "degraded"},
		{StatusMaintenance, "maintenance"},
		{StatusOffline, "offline"},
		{Status(99), "unknown"},
	}

	for _, tt := range tests {
		got := tt.status.String()
		if got != tt.expected {
			t.Errorf("Status(%d).String() = %s, expected %s", tt.status, got, tt.expected)
		}
	}
}

func TestStatusIsAvailable(t *testing.T) {
	tests := []struct {
		status   Status
		expected bool
	}{
		{StatusUnknown, false},
		{StatusHealthy, true},
		{StatusUnhealthy, false},
		{StatusDegraded, true},
		{StatusMaintenance, false},
		{StatusOffline, false},
	}

	for _, tt := range tests {
		got := tt.status.IsAvailable()
		if got != tt.expected {
			t.Errorf("Status(%d).IsAvailable() = %v, expected %v", tt.status, got, tt.expected)
		}
	}
}

// ============== Info Tests ==============

func TestInfoClone(t *testing.T) {
	original := &Info{
		ID:        "test-id",
		ServiceID: "test-service",
		Host:      "localhost",
		Port:      8080,
		Protocol:  ProtocolHTTP,
		Status:    StatusHealthy,
		Methods:   []Method{MethodGET, MethodPOST},
		Tags:      []string{"api", "v1"},
		Metadata:  map[string]string{"key": "value"},
	}

	clone := original.Clone()

	if clone == original {
		t.Error("Clone returned same pointer")
	}
	if clone.ID != original.ID {
		t.Error("ID not cloned correctly")
	}

	// Modify clone and verify original unchanged
	clone.Tags[0] = "modified"
	if original.Tags[0] == "modified" {
		t.Error("Tags not deep copied")
	}

	clone.Metadata["key"] = "modified"
	if original.Metadata["key"] == "modified" {
		t.Error("Metadata not deep copied")
	}
}

func TestInfoCloneNil(t *testing.T) {
	var info *Info
	clone := info.Clone()
	if clone != nil {
		t.Error("Clone of nil should be nil")
	}
}

func TestInfoURL(t *testing.T) {
	tests := []struct {
		info     Info
		expected string
	}{
		{
			info:     Info{Protocol: ProtocolHTTP, Host: "localhost", Port: 8080, Path: "/api"},
			expected: "http://localhost:8080/api",
		},
		{
			info:     Info{Protocol: ProtocolHTTP, Host: "localhost", Port: 80, Path: "/api"},
			expected: "http://localhost/api",
		},
		{
			info:     Info{Protocol: ProtocolHTTPS, Host: "localhost", Port: 443, Path: "/api"},
			expected: "https://localhost/api",
		},
		{
			info:     Info{Protocol: ProtocolHTTPS, Host: "localhost", Port: 8443, Path: "/api"},
			expected: "https://localhost:8443/api",
		},
		{
			info:     Info{Protocol: ProtocolHTTP, Host: "localhost", Port: 0, Path: "/api"},
			expected: "http://localhost/api",
		},
	}

	for _, tt := range tests {
		got := tt.info.URL()
		if got != tt.expected {
			t.Errorf("URL() = %s, expected %s", got, tt.expected)
		}
	}
}

func TestInfoAddress(t *testing.T) {
	tests := []struct {
		info     Info
		expected string
	}{
		{Info{Host: "localhost", Port: 8080}, "localhost:8080"},
		{Info{Host: "localhost", Port: 0}, "localhost"},
	}

	for _, tt := range tests {
		got := tt.info.Address()
		if got != tt.expected {
			t.Errorf("Address() = %s, expected %s", got, tt.expected)
		}
	}
}

func TestInfoIsExpired(t *testing.T) {
	// Not expired - no expiry set
	info := &Info{}
	if info.IsExpired() {
		t.Error("Expected non-expired when ExpiresAt is zero")
	}

	// Not expired - future time
	info.ExpiresAt = time.Now().Add(time.Hour)
	if info.IsExpired() {
		t.Error("Expected non-expired for future time")
	}

	// Expired - past time
	info.ExpiresAt = time.Now().Add(-time.Hour)
	if !info.IsExpired() {
		t.Error("Expected expired for past time")
	}
}

func TestInfoIsHealthy(t *testing.T) {
	info := &Info{Status: StatusHealthy}
	if !info.IsHealthy() {
		t.Error("Expected healthy")
	}

	info.Status = StatusUnhealthy
	if info.IsHealthy() {
		t.Error("Expected not healthy")
	}
}

func TestInfoHasTag(t *testing.T) {
	info := &Info{Tags: []string{"api", "v1", "production"}}

	if !info.HasTag("api") {
		t.Error("Expected to have 'api' tag")
	}
	if info.HasTag("staging") {
		t.Error("Expected not to have 'staging' tag")
	}
}

func TestInfoGetMetadata(t *testing.T) {
	info := &Info{Metadata: map[string]string{"key": "value"}}

	if got := info.GetMetadata("key"); got != "value" {
		t.Errorf("GetMetadata('key') = %s, expected 'value'", got)
	}
	if got := info.GetMetadata("nonexistent"); got != "" {
		t.Errorf("GetMetadata('nonexistent') = %s, expected empty", got)
	}

	// Test with nil metadata
	info2 := &Info{}
	if got := info2.GetMetadata("key"); got != "" {
		t.Errorf("GetMetadata on nil metadata = %s, expected empty", got)
	}
}

func TestInfoSetMetadata(t *testing.T) {
	info := &Info{}
	info.SetMetadata("key", "value")

	if info.Metadata == nil {
		t.Fatal("Metadata should be initialized")
	}
	if info.Metadata["key"] != "value" {
		t.Error("Metadata not set correctly")
	}

	info.SetMetadata("key2", "value2")
	if info.Metadata["key2"] != "value2" {
		t.Error("Second metadata not set correctly")
	}
}

// ============== RegistrationRequest Tests ==============

func TestRegistrationRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     RegistrationRequest
		wantErr bool
	}{
		{
			name:    "nil endpoint",
			req:     RegistrationRequest{Endpoint: nil},
			wantErr: true,
		},
		{
			name:    "missing ID",
			req:     RegistrationRequest{Endpoint: &Info{ServiceID: "svc", Host: "localhost"}},
			wantErr: true,
		},
		{
			name:    "missing service ID",
			req:     RegistrationRequest{Endpoint: &Info{ID: "ep1", Host: "localhost"}},
			wantErr: true,
		},
		{
			name:    "missing host",
			req:     RegistrationRequest{Endpoint: &Info{ID: "ep1", ServiceID: "svc"}},
			wantErr: true,
		},
		{
			name:    "valid",
			req:     RegistrationRequest{Endpoint: &Info{ID: "ep1", ServiceID: "svc", Host: "localhost"}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ============== DeregistrationRequest Tests ==============

func TestDeregistrationRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     DeregistrationRequest
		wantErr bool
	}{
		{
			name:    "missing endpoint ID",
			req:     DeregistrationRequest{},
			wantErr: true,
		},
		{
			name:    "valid",
			req:     DeregistrationRequest{EndpointID: "ep1"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// ============== Query Tests ==============

func TestQueryMatches(t *testing.T) {
	endpoint := &Info{
		ID:         "ep1",
		ServiceID:  "svc1",
		InstanceID: "inst1",
		Host:       "localhost",
		Port:       8080,
		Protocol:   ProtocolHTTP,
		Status:     StatusHealthy,
		Tags:       []string{"api", "v1"},
		Metadata:   map[string]string{"env": "prod"},
		Version:    "1.0.0",
		APIVersion: "v1",
		TenantID:   "tenant1",
		Region:     "us-east-1",
		Zone:       "zone-a",
	}

	tests := []struct {
		name    string
		query   Query
		matches bool
	}{
		{
			name:    "nil endpoint",
			query:   Query{},
			matches: false,
		},
		{
			name:    "empty query matches all",
			query:   Query{},
			matches: true,
		},
		{
			name:    "ID filter match",
			query:   Query{IDs: []string{"ep1"}},
			matches: true,
		},
		{
			name:    "ID filter no match",
			query:   Query{IDs: []string{"ep2"}},
			matches: false,
		},
		{
			name:    "ServiceID filter match",
			query:   Query{ServiceIDs: []string{"svc1"}},
			matches: true,
		},
		{
			name:    "ServiceID filter no match",
			query:   Query{ServiceIDs: []string{"svc2"}},
			matches: false,
		},
		{
			name:    "InstanceID filter match",
			query:   Query{InstanceID: "inst1"},
			matches: true,
		},
		{
			name:    "InstanceID filter no match",
			query:   Query{InstanceID: "inst2"},
			matches: false,
		},
		{
			name:    "Host filter match",
			query:   Query{Host: "localhost"},
			matches: true,
		},
		{
			name:    "Host filter no match",
			query:   Query{Host: "remotehost"},
			matches: false,
		},
		{
			name:    "Port filter match",
			query:   Query{Port: 8080},
			matches: true,
		},
		{
			name:    "Port filter no match",
			query:   Query{Port: 9090},
			matches: false,
		},
		{
			name:    "Protocol filter match",
			query:   Query{Protocol: ProtocolHTTP},
			matches: true,
		},
		{
			name:    "Protocol filter no match",
			query:   Query{Protocol: ProtocolHTTPS},
			matches: false,
		},
		{
			name:    "Statuses filter match",
			query:   Query{Statuses: []Status{StatusHealthy, StatusDegraded}},
			matches: true,
		},
		{
			name:    "Statuses filter no match",
			query:   Query{Statuses: []Status{StatusUnhealthy}},
			matches: false,
		},
		{
			name:    "OnlyHealthy match",
			query:   Query{OnlyHealthy: true},
			matches: true,
		},
		{
			name:    "OnlyAvailable match",
			query:   Query{OnlyAvailable: true},
			matches: true,
		},
		{
			name:    "Tag filter match",
			query:   Query{Tags: []string{"api"}},
			matches: true,
		},
		{
			name:    "Tag filter no match",
			query:   Query{Tags: []string{"grpc"}},
			matches: false,
		},
		{
			name:    "Metadata filter match",
			query:   Query{Metadata: map[string]string{"env": "prod"}},
			matches: true,
		},
		{
			name:    "Metadata filter no match",
			query:   Query{Metadata: map[string]string{"env": "dev"}},
			matches: false,
		},
		{
			name:    "Version filter match",
			query:   Query{Version: "1.0.0"},
			matches: true,
		},
		{
			name:    "Version filter no match",
			query:   Query{Version: "2.0.0"},
			matches: false,
		},
		{
			name:    "APIVersion filter match",
			query:   Query{APIVersion: "v1"},
			matches: true,
		},
		{
			name:    "APIVersion filter no match",
			query:   Query{APIVersion: "v2"},
			matches: false,
		},
		{
			name:    "TenantID filter match",
			query:   Query{TenantID: "tenant1"},
			matches: true,
		},
		{
			name:    "TenantID filter no match",
			query:   Query{TenantID: "tenant2"},
			matches: false,
		},
		{
			name:    "Region filter match",
			query:   Query{Region: "us-east-1"},
			matches: true,
		},
		{
			name:    "Region filter no match",
			query:   Query{Region: "eu-west-1"},
			matches: false,
		},
		{
			name:    "Zone filter match",
			query:   Query{Zone: "zone-a"},
			matches: true,
		},
		{
			name:    "Zone filter no match",
			query:   Query{Zone: "zone-b"},
			matches: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ep *Info
			if tt.name != "nil endpoint" {
				ep = endpoint
			}
			got := tt.query.Matches(ep)
			if got != tt.matches {
				t.Errorf("Matches() = %v, want %v", got, tt.matches)
			}
		})
	}
}

// ============== Event Tests ==============

func TestEventToMessage(t *testing.T) {
	event := &Event{
		Type:       EventRegistered,
		EndpointID: "ep1",
		ServiceID:  "svc1",
		Timestamp:  time.Now(),
	}

	msg, err := event.ToMessage()
	if err != nil {
		t.Fatalf("ToMessage() error = %v", err)
	}
	if msg == nil {
		t.Fatal("ToMessage() returned nil")
	}
	if msg.Headers.Get("event-type") != string(EventRegistered) {
		t.Error("event-type header not set")
	}
	if msg.Headers.Get("endpoint-id") != "ep1" {
		t.Error("endpoint-id header not set")
	}
	if msg.Headers.Get("service-id") != "svc1" {
		t.Error("service-id header not set")
	}
}

func TestEventFromMessage(t *testing.T) {
	originalEvent := &Event{
		Type:       EventRegistered,
		EndpointID: "ep1",
		ServiceID:  "svc1",
		Timestamp:  time.Now().UTC().Truncate(time.Second),
	}

	msg, _ := originalEvent.ToMessage()
	recovered, err := EventFromMessage(msg)
	if err != nil {
		t.Fatalf("EventFromMessage() error = %v", err)
	}
	if recovered.Type != originalEvent.Type {
		t.Errorf("Type = %s, want %s", recovered.Type, originalEvent.Type)
	}
	if recovered.EndpointID != originalEvent.EndpointID {
		t.Errorf("EndpointID = %s, want %s", recovered.EndpointID, originalEvent.EndpointID)
	}
}

func TestEventFromMessageInvalid(t *testing.T) {
	msg := core.NewMessage([]byte("invalid json"))
	_, err := EventFromMessage(msg)
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
}

func TestRegistryRegisterAndGet(t *testing.T) {
	r := NewRegistry()

	ep := &Info{
		ID:        "ep1",
		ServiceID: "svc1",
		Host:      "localhost",
		Port:      8080,
		Protocol:  ProtocolHTTP,
		Status:    StatusHealthy,
	}

	err := r.Register(ep)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	retrieved, found := r.Get("ep1")
	if !found {
		t.Fatal("Endpoint not found")
	}
	if retrieved.ID != ep.ID {
		t.Errorf("Retrieved ID = %s, want %s", retrieved.ID, ep.ID)
	}
}

func TestRegistryRegisterDuplicate(t *testing.T) {
	r := NewRegistry()

	ep := &Info{
		ID:        "ep1",
		ServiceID: "svc1",
		Host:      "localhost",
		Port:      8080,
	}

	err := r.Register(ep)
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}

	// Second registration with same ID should update, not error
	ep2 := &Info{
		ID:        "ep1",
		ServiceID: "svc1",
		Host:      "localhost",
		Port:      9090, // Updated port
	}
	err = r.Register(ep2)
	if err != nil {
		t.Errorf("Second registration should succeed (update): %v", err)
	}

	// Verify it was updated
	got, exists := r.Get("ep1")
	if !exists {
		t.Fatal("Endpoint should exist")
	}
	if got.Port != 9090 {
		t.Errorf("Port should be updated to 9090, got %d", got.Port)
	}
}

func TestRegistryRegisterWithTTL(t *testing.T) {
	r := NewRegistry()

	ep := &Info{
		ID:        "ep1",
		ServiceID: "svc1",
		Host:      "localhost",
	}

	err := r.RegisterWithTTL(ep, 5*time.Minute)
	if err != nil {
		t.Fatalf("RegisterWithTTL() error = %v", err)
	}

	retrieved, _ := r.Get("ep1")
	if retrieved.ExpiresAt.IsZero() {
		t.Error("ExpiresAt should be set")
	}
}

func TestRegistryGetNotFound(t *testing.T) {
	r := NewRegistry()

	_, found := r.Get("nonexistent")
	if found {
		t.Error("Expected not found")
	}
}

func TestRegistryDeregister(t *testing.T) {
	r := NewRegistry()

	ep := &Info{
		ID:        "ep1",
		ServiceID: "svc1",
		Host:      "localhost",
	}

	r.Register(ep)

	err := r.Deregister("ep1")
	if err != nil {
		t.Errorf("Deregister should succeed: %v", err)
	}

	_, found := r.Get("ep1")
	if found {
		t.Error("Endpoint should be removed")
	}
}

func TestRegistryDeregisterNotFound(t *testing.T) {
	r := NewRegistry()

	err := r.Deregister("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent endpoint")
	}
}

func TestRegistryDeregisterWithReason(t *testing.T) {
	r := NewRegistry()

	ep := &Info{
		ID:        "ep1",
		ServiceID: "svc1",
		Host:      "localhost",
	}

	r.Register(ep)

	err := r.DeregisterWithReason("ep1", "test reason")
	if err != nil {
		t.Errorf("DeregisterWithReason should succeed: %v", err)
	}
}

func TestRegistryDeregisterByService(t *testing.T) {
	r := NewRegistry()

	r.Register(&Info{ID: "ep1", ServiceID: "svc1", Host: "host1"})
	r.Register(&Info{ID: "ep2", ServiceID: "svc1", Host: "host2"})
	r.Register(&Info{ID: "ep3", ServiceID: "svc2", Host: "host3"})

	count, err := r.DeregisterByService("svc1")
	if err != nil {
		t.Fatalf("DeregisterByService() error = %v", err)
	}
	if count != 2 {
		t.Errorf("Count = %d, want 2", count)
	}

	if r.Count() != 1 {
		t.Errorf("Registry count = %d, want 1", r.Count())
	}
}

func TestRegistryDeregisterByInstance(t *testing.T) {
	r := NewRegistry()

	r.Register(&Info{ID: "ep1", ServiceID: "svc1", Host: "host1", InstanceID: "inst1"})
	r.Register(&Info{ID: "ep2", ServiceID: "svc1", Host: "host2", InstanceID: "inst1"})
	r.Register(&Info{ID: "ep3", ServiceID: "svc2", Host: "host3", InstanceID: "inst2"})

	count, err := r.DeregisterByInstance("inst1")
	if err != nil {
		t.Fatalf("DeregisterByInstance() error = %v", err)
	}
	if count != 2 {
		t.Errorf("Count = %d, want 2", count)
	}
}

func TestRegistryUpdateStatus(t *testing.T) {
	r := NewRegistry()

	ep := &Info{
		ID:        "ep1",
		ServiceID: "svc1",
		Host:      "localhost",
		Status:    StatusHealthy,
	}

	r.Register(ep)

	err := r.UpdateStatus("ep1", StatusDegraded)
	if err != nil {
		t.Fatalf("UpdateStatus() error = %v", err)
	}

	retrieved, _ := r.Get("ep1")
	if retrieved.Status != StatusDegraded {
		t.Errorf("Status = %v, want %v", retrieved.Status, StatusDegraded)
	}
}

func TestRegistryQuery(t *testing.T) {
	r := NewRegistry()

	// Register multiple endpoints
	endpoints := []*Info{
		{ID: "ep1", ServiceID: "svc1", Host: "host1", Status: StatusHealthy, Tags: []string{"api"}},
		{ID: "ep2", ServiceID: "svc1", Host: "host2", Status: StatusUnhealthy, Tags: []string{"api"}},
		{ID: "ep3", ServiceID: "svc2", Host: "host3", Status: StatusHealthy, Tags: []string{"grpc"}},
	}

	for _, ep := range endpoints {
		r.Register(ep)
	}

	// Query by service
	results := r.Query(&Query{ServiceIDs: []string{"svc1"}})
	if len(results) != 2 {
		t.Errorf("Query by service: got %d, want 2", len(results))
	}

	// Query healthy only
	results = r.Query(&Query{OnlyHealthy: true})
	if len(results) != 2 {
		t.Errorf("Query healthy: got %d, want 2", len(results))
	}

	// Query by tag
	results = r.Query(&Query{Tags: []string{"api"}})
	if len(results) != 2 {
		t.Errorf("Query by tag: got %d, want 2", len(results))
	}

	// Empty query
	results = r.Query(nil)
	if len(results) != 3 {
		t.Errorf("Empty query: got %d, want 3", len(results))
	}
}

func TestRegistryQueryWithPagination(t *testing.T) {
	r := NewRegistry()

	for i := 0; i < 10; i++ {
		r.Register(&Info{
			ID:        "ep" + itoa(i),
			ServiceID: "svc",
			Host:      "host",
		})
	}

	results := r.Query(&Query{Limit: 5})
	if len(results) != 5 {
		t.Errorf("Limit: got %d, want 5", len(results))
	}
}

func TestRegistryAll(t *testing.T) {
	r := NewRegistry()

	r.Register(&Info{ID: "ep1", ServiceID: "svc1", Host: "host1"})
	r.Register(&Info{ID: "ep2", ServiceID: "svc1", Host: "host2"})

	all := r.All()
	if len(all) != 2 {
		t.Errorf("All() returned %d endpoints, want 2", len(all))
	}
}

func TestRegistryCount(t *testing.T) {
	r := NewRegistry()

	if r.Count() != 0 {
		t.Error("Empty registry should have count 0")
	}

	r.Register(&Info{ID: "ep1", ServiceID: "svc1", Host: "host1"})
	r.Register(&Info{ID: "ep2", ServiceID: "svc1", Host: "host2"})

	if r.Count() != 2 {
		t.Errorf("Count() = %d, want 2", r.Count())
	}
}

func TestRegistryServices(t *testing.T) {
	r := NewRegistry()

	r.Register(&Info{ID: "ep1", ServiceID: "svc1", Host: "host1"})
	r.Register(&Info{ID: "ep2", ServiceID: "svc2", Host: "host2"})
	r.Register(&Info{ID: "ep3", ServiceID: "svc1", Host: "host3"})

	services := r.Services()
	if len(services) != 2 {
		t.Errorf("Services() returned %d services, want 2", len(services))
	}
}

func TestRegistryGetByService(t *testing.T) {
	r := NewRegistry()

	r.Register(&Info{ID: "ep1", ServiceID: "svc1", Host: "host1"})
	r.Register(&Info{ID: "ep2", ServiceID: "svc2", Host: "host2"})
	r.Register(&Info{ID: "ep3", ServiceID: "svc1", Host: "host3"})

	endpoints := r.GetByService("svc1")
	if len(endpoints) != 2 {
		t.Errorf("GetByService() returned %d endpoints, want 2", len(endpoints))
	}
}

func TestRegistryCallbacks(t *testing.T) {
	r := NewRegistry()

	registeredCalled := false
	deregisteredCalled := false
	statusChangeCalled := false

	r.OnRegistered(func(ep *Info) {
		registeredCalled = true
	})
	r.OnDeregistered(func(ep *Info, reason string) {
		deregisteredCalled = true
	})
	r.OnStatusChange(func(ep *Info, oldStatus, newStatus Status) {
		statusChangeCalled = true
	})

	ep := &Info{
		ID:        "ep1",
		ServiceID: "svc1",
		Host:      "localhost",
		Status:    StatusHealthy,
	}

	r.Register(ep)
	time.Sleep(10 * time.Millisecond) // Give callback goroutine time
	if !registeredCalled {
		t.Error("OnRegistered callback not called")
	}

	r.UpdateStatus("ep1", StatusDegraded)
	time.Sleep(10 * time.Millisecond)
	if !statusChangeCalled {
		t.Error("OnStatusChange callback not called")
	}

	r.Deregister("ep1")
	time.Sleep(10 * time.Millisecond)
	if !deregisteredCalled {
		t.Error("OnDeregistered callback not called")
	}
}

func TestRegistryConcurrentAccess(t *testing.T) {
	r := NewRegistry()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			ep := &Info{
				ID:        "ep" + itoa(i),
				ServiceID: "svc",
				Host:      "host",
			}
			r.Register(ep)
			r.Get("ep" + itoa(i))
			r.UpdateStatus("ep"+itoa(i), StatusHealthy)
		}(i)
	}
	wg.Wait()

	if r.Count() != 100 {
		t.Errorf("Count = %d, want 100", r.Count())
	}
}

// ============== Error Tests ==============

func TestErrors(t *testing.T) {
	tests := []struct {
		err      error
		message  string
		isNotNil bool
	}{
		{ErrInvalidEndpoint, "invalid endpoint", true},
		{ErrMissingEndpointID, "missing endpoint ID", true},
		{ErrMissingServiceID, "missing service ID", true},
		{ErrMissingHost, "missing host", true},
		{ErrEndpointNotFound, "endpoint not found", true},
		{ErrEndpointExists, "endpoint already exists", true},
		{ErrEndpointExpired, "endpoint has expired", true},
		{ErrHealthCheckFailed, "health check failed", true},
		{ErrHealthCheckTimeout, "health check timeout", true},
		{ErrRegistryShutdown, "registry is shutting down", true},
		{ErrInvalidQuery, "invalid query parameters", true},
		{ErrManagerNotStarted, "manager not started", true},
		{ErrManagerShutdown, "manager is shutting down", true},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			if tt.err == nil && tt.isNotNil {
				t.Error("Error should not be nil")
			}
			if tt.err != nil && tt.err.Error() != tt.message {
				t.Errorf("Error message = %s, want %s", tt.err.Error(), tt.message)
			}
		})
	}
}

func TestEndpointError(t *testing.T) {
	err := NewEndpointError("ep1", "svc1", "register", ErrEndpointExists)

	if err.Error() == "" {
		t.Error("Error string should not be empty")
	}
	if err.Unwrap() != ErrEndpointExists {
		t.Error("Unwrap() should return underlying error")
	}
}

func TestEndpointErrorWithoutServiceID(t *testing.T) {
	err := NewEndpointError("ep1", "", "register", ErrEndpointExists)

	errStr := err.Error()
	if errStr == "" {
		t.Error("Error string should not be empty")
	}
}

func TestEndpointErrorWithoutOp(t *testing.T) {
	err := NewEndpointError("ep1", "", "", ErrEndpointExists)

	errStr := err.Error()
	if errStr == "" {
		t.Error("Error string should not be empty")
	}
}

func TestRegistrationError(t *testing.T) {
	err := RegistrationError{
		EndpointID: "ep1",
		Reason:     "duplicate",
		Err:        ErrEndpointExists,
	}

	if err.Error() == "" {
		t.Error("Error string should not be empty")
	}
	if err.Unwrap() != ErrEndpointExists {
		t.Error("Unwrap() should return underlying error")
	}
}

func TestRegistrationErrorWithoutReason(t *testing.T) {
	err := RegistrationError{
		EndpointID: "ep1",
		Err:        ErrEndpointExists,
	}

	errStr := err.Error()
	if errStr == "" {
		t.Error("Error string should not be empty")
	}
}

func TestHealthCheckError(t *testing.T) {
	err := HealthCheckError{
		EndpointID:       "ep1",
		URL:              "http://localhost:8080/health",
		ConsecutiveFails: 3,
		Err:              ErrHealthCheckTimeout,
	}

	if err.Error() == "" {
		t.Error("Error string should not be empty")
	}
	if err.Unwrap() != ErrHealthCheckTimeout {
		t.Error("Unwrap() should return underlying error")
	}
}

func TestIsEndpointError(t *testing.T) {
	epErr := NewEndpointError("ep1", "svc1", "register", ErrEndpointExists)
	if !IsEndpointError(epErr) {
		t.Error("Should return true for EndpointError")
	}

	regErr := RegistrationError{EndpointID: "ep1", Err: ErrEndpointExists}
	if !IsEndpointError(regErr) {
		t.Error("Should return true for RegistrationError")
	}

	healthErr := HealthCheckError{EndpointID: "ep1", Err: ErrHealthCheckFailed}
	if !IsEndpointError(healthErr) {
		t.Error("Should return true for HealthCheckError")
	}

	if IsEndpointError(ErrInvalidEndpoint) {
		t.Error("Should return false for base error")
	}
}

func TestIsNotFound(t *testing.T) {
	if !IsNotFound(ErrEndpointNotFound) {
		t.Error("Should return true for ErrEndpointNotFound")
	}
	if IsNotFound(ErrEndpointExists) {
		t.Error("Should return false for other errors")
	}
}

func TestIsExists(t *testing.T) {
	if !IsExists(ErrEndpointExists) {
		t.Error("Should return true for ErrEndpointExists")
	}
	if IsExists(ErrEndpointNotFound) {
		t.Error("Should return false for other errors")
	}
}

// ============== Manager Tests ==============

func TestNewManager(t *testing.T) {
	cfg := ManagerConfig{
		HealthCheck: HealthCheckConfig{
			Enabled: false,
		},
	}

	manager, err := NewManager(cfg)
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}
	if manager == nil {
		t.Fatal("NewManager() returned nil")
	}
}

func TestManagerRegisterEndpoint(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	ep := &Info{
		ID:        "ep1",
		ServiceID: "svc1",
		Host:      "localhost",
		Port:      8080,
		Protocol:  ProtocolHTTP,
		Status:    StatusHealthy,
	}

	err := manager.RegisterEndpoint(ctx, ep)
	if err != nil {
		t.Fatalf("RegisterEndpoint() error = %v", err)
	}

	retrieved, found := manager.Get("ep1")
	if !found {
		t.Fatal("Endpoint not found")
	}
	if retrieved.ID != ep.ID {
		t.Error("Retrieved endpoint ID mismatch")
	}
}

func TestManagerDeregisterEndpoint(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	ep := &Info{
		ID:        "ep1",
		ServiceID: "svc1",
		Host:      "localhost",
	}

	manager.RegisterEndpoint(ctx, ep)

	err := manager.DeregisterEndpoint(ctx, "ep1")
	if err != nil {
		t.Fatalf("DeregisterEndpoint() error = %v", err)
	}

	_, found := manager.Get("ep1")
	if found {
		t.Error("Endpoint should be removed")
	}
}

func TestManagerDiscover(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{
		ID:        "ep1",
		ServiceID: "svc1",
		Host:      "host1",
		Status:    StatusHealthy,
	})
	manager.RegisterEndpoint(ctx, &Info{
		ID:        "ep2",
		ServiceID: "svc1",
		Host:      "host2",
		Status:    StatusDegraded, // Degraded is "available" (Healthy + Degraded)
	})
	manager.RegisterEndpoint(ctx, &Info{
		ID:        "ep3",
		ServiceID: "svc1",
		Host:      "host3",
		Status:    StatusUnhealthy, // Unhealthy is NOT available
	})

	// Discover available (healthy + degraded) - default behavior
	endpoints := manager.Discover("svc1")
	if len(endpoints) != 2 {
		t.Errorf("Discover() returned %d endpoints, want 2 (healthy + degraded)", len(endpoints))
	}

	// Discover healthy only (excludes degraded)
	endpoints = manager.Discover("svc1", OnlyHealthy())
	if len(endpoints) != 1 {
		t.Errorf("Discover(OnlyHealthy) returned %d endpoints, want 1", len(endpoints))
	}
}

func TestManagerServices(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: "ep1", ServiceID: "svc1", Host: "host1"})
	manager.RegisterEndpoint(ctx, &Info{ID: "ep2", ServiceID: "svc2", Host: "host2"})

	services := manager.Services()
	if len(services) != 2 {
		t.Errorf("Services() returned %d services, want 2", len(services))
	}
}

func TestManagerShutdown(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	manager.Shutdown(ctx)
	// Should not panic
}

// ============== Discovery Options Tests ==============

func TestDiscoveryOptions(t *testing.T) {
	q := &Query{}

	OnlyHealthy()(q)
	if !q.OnlyHealthy {
		t.Error("OnlyHealthy option not set")
	}

	q2 := &Query{}
	WithTags("tag1", "tag2")(q2)
	if len(q2.Tags) != 2 {
		t.Error("WithTags option not set correctly")
	}

	q3 := &Query{}
	WithVersion("1.0.0")(q3)
	if q3.Version != "1.0.0" {
		t.Error("WithVersion option not set")
	}

	q4 := &Query{}
	WithRegion("us-east-1")(q4)
	if q4.Region != "us-east-1" {
		t.Error("WithRegion option not set")
	}

	q5 := &Query{}
	WithLimit(10)(q5)
	if q5.Limit != 10 {
		t.Error("WithLimit option not set")
	}

	q6 := &Query{}
	WithProtocol(ProtocolHTTP)(q6)
	if q6.Protocol != ProtocolHTTP {
		t.Error("WithProtocol option not set")
	}
}

// ============== Helper Function Tests ==============

func TestItoa(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{123, "123"},
		{-1, "-1"},
		{-123, "-123"},
	}

	for _, tt := range tests {
		got := itoa(tt.input)
		if got != tt.expected {
			t.Errorf("itoa(%d) = %s, expected %s", tt.input, got, tt.expected)
		}
	}
}

func TestContainsString(t *testing.T) {
	slice := []string{"a", "b", "c"}

	if !containsString(slice, "b") {
		t.Error("Expected true for 'b'")
	}
	if containsString(slice, "d") {
		t.Error("Expected false for 'd'")
	}
}

func TestContainsStatus(t *testing.T) {
	slice := []Status{StatusHealthy, StatusDegraded}

	if !containsStatus(slice, StatusHealthy) {
		t.Error("Expected true for StatusHealthy")
	}
	if containsStatus(slice, StatusUnhealthy) {
		t.Error("Expected false for StatusUnhealthy")
	}
}

// ============== Benchmarks ==============

func BenchmarkRegistryRegister(b *testing.B) {
	r := NewRegistry()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ep := &Info{
			ID:        "ep" + itoa(i),
			ServiceID: "svc",
			Host:      "localhost",
		}
		r.Register(ep)
	}
}

func BenchmarkRegistryGet(b *testing.B) {
	r := NewRegistry()

	for i := 0; i < 1000; i++ {
		r.Register(&Info{
			ID:        "ep" + itoa(i),
			ServiceID: "svc",
			Host:      "localhost",
		})
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r.Get("ep" + itoa(i%1000))
	}
}

func BenchmarkRegistryQuery(b *testing.B) {
	r := NewRegistry()

	for i := 0; i < 1000; i++ {
		r.Register(&Info{
			ID:        "ep" + itoa(i),
			ServiceID: "svc" + itoa(i%10),
			Host:      "localhost",
			Status:    StatusHealthy,
			Tags:      []string{"api"},
		})
	}

	query := &Query{
		ServiceIDs:  []string{"svc1", "svc2"},
		OnlyHealthy: true,
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r.Query(query)
	}
}

func BenchmarkInfoClone(b *testing.B) {
	info := &Info{
		ID:        "ep1",
		ServiceID: "svc1",
		Host:      "localhost",
		Tags:      []string{"api", "v1"},
		Metadata:  map[string]string{"key": "value"},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		info.Clone()
	}
}

func BenchmarkQueryMatches(b *testing.B) {
	ep := &Info{
		ID:        "ep1",
		ServiceID: "svc1",
		Host:      "localhost",
		Status:    StatusHealthy,
		Tags:      []string{"api"},
	}

	query := &Query{
		ServiceIDs:  []string{"svc1"},
		OnlyHealthy: true,
		Tags:        []string{"api"},
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		query.Matches(ep)
	}
}

// ============== Additional Registry Tests ==============

func TestRegistryGetByInstance(t *testing.T) {
	r := NewRegistry()

	r.Register(&Info{ID: "ep1", ServiceID: "svc1", Host: "host1", InstanceID: "inst1"})
	r.Register(&Info{ID: "ep2", ServiceID: "svc1", Host: "host2", InstanceID: "inst1"})
	r.Register(&Info{ID: "ep3", ServiceID: "svc2", Host: "host3", InstanceID: "inst2"})

	endpoints := r.GetByInstance("inst1")
	if len(endpoints) != 2 {
		t.Errorf("GetByInstance() returned %d endpoints, want 2", len(endpoints))
	}

	endpoints = r.GetByInstance("nonexistent")
	if len(endpoints) != 0 {
		t.Errorf("GetByInstance() for nonexistent returned %d endpoints, want 0", len(endpoints))
	}
}

func TestRegistryGetByTag(t *testing.T) {
	r := NewRegistry()

	r.Register(&Info{ID: "ep1", ServiceID: "svc1", Host: "host1", Tags: []string{"api", "v1"}})
	r.Register(&Info{ID: "ep2", ServiceID: "svc1", Host: "host2", Tags: []string{"api", "v2"}})
	r.Register(&Info{ID: "ep3", ServiceID: "svc2", Host: "host3", Tags: []string{"grpc"}})

	endpoints := r.GetByTag("api")
	if len(endpoints) != 2 {
		t.Errorf("GetByTag('api') returned %d endpoints, want 2", len(endpoints))
	}

	endpoints = r.GetByTag("v1")
	if len(endpoints) != 1 {
		t.Errorf("GetByTag('v1') returned %d endpoints, want 1", len(endpoints))
	}

	endpoints = r.GetByTag("nonexistent")
	if len(endpoints) != 0 {
		t.Errorf("GetByTag('nonexistent') returned %d endpoints, want 0", len(endpoints))
	}
}

func TestRegistryGetHealthy(t *testing.T) {
	r := NewRegistry()

	r.Register(&Info{ID: "ep1", ServiceID: "svc1", Host: "host1", Status: StatusHealthy})
	r.Register(&Info{ID: "ep2", ServiceID: "svc1", Host: "host2", Status: StatusUnhealthy})
	r.Register(&Info{ID: "ep3", ServiceID: "svc2", Host: "host3", Status: StatusHealthy})

	endpoints := r.GetHealthy()
	if len(endpoints) != 2 {
		t.Errorf("GetHealthy() returned %d endpoints, want 2", len(endpoints))
	}
}

func TestRegistryGetAvailable(t *testing.T) {
	r := NewRegistry()

	r.Register(&Info{ID: "ep1", ServiceID: "svc1", Host: "host1", Status: StatusHealthy})
	r.Register(&Info{ID: "ep2", ServiceID: "svc1", Host: "host2", Status: StatusDegraded})
	r.Register(&Info{ID: "ep3", ServiceID: "svc2", Host: "host3", Status: StatusUnhealthy})

	endpoints := r.GetAvailable()
	if len(endpoints) != 2 {
		t.Errorf("GetAvailable() returned %d endpoints, want 2", len(endpoints))
	}
}

func TestRegistryServiceCount(t *testing.T) {
	r := NewRegistry()

	if r.ServiceCount() != 0 {
		t.Error("Empty registry should have service count 0")
	}

	r.Register(&Info{ID: "ep1", ServiceID: "svc1", Host: "host1"})
	r.Register(&Info{ID: "ep2", ServiceID: "svc2", Host: "host2"})
	r.Register(&Info{ID: "ep3", ServiceID: "svc1", Host: "host3"})

	if r.ServiceCount() != 2 {
		t.Errorf("ServiceCount() = %d, want 2", r.ServiceCount())
	}
}

func TestRegistryExists(t *testing.T) {
	r := NewRegistry()

	r.Register(&Info{ID: "ep1", ServiceID: "svc1", Host: "host1"})

	if !r.Exists("ep1") {
		t.Error("Exists() should return true for existing endpoint")
	}
	if r.Exists("nonexistent") {
		t.Error("Exists() should return false for non-existing endpoint")
	}
}

func TestRegistryUpdateHealth(t *testing.T) {
	r := NewRegistry()

	r.Register(&Info{ID: "ep1", ServiceID: "svc1", Host: "host1", Status: StatusHealthy})

	err := r.UpdateHealth(&HealthUpdate{
		EndpointID: "ep1",
		Status:     StatusHealthy,
		CheckedAt:  time.Now(),
	})
	if err != nil {
		t.Fatalf("UpdateHealth() error = %v", err)
	}

	ep, _ := r.Get("ep1")
	if ep.Status != StatusHealthy {
		t.Errorf("Status = %v, want %v", ep.Status, StatusHealthy)
	}

	err = r.UpdateHealth(&HealthUpdate{
		EndpointID: "ep1",
		Status:     StatusUnhealthy,
		Message:    "health check failed",
		CheckedAt:  time.Now(),
	})
	if err != nil {
		t.Fatalf("UpdateHealth(unhealthy) error = %v", err)
	}

	ep, _ = r.Get("ep1")
	if ep.Status != StatusUnhealthy {
		t.Errorf("Status = %v, want %v", ep.Status, StatusUnhealthy)
	}

	// Test not found
	err = r.UpdateHealth(&HealthUpdate{
		EndpointID: "nonexistent",
		Status:     StatusHealthy,
		CheckedAt:  time.Now(),
	})
	if err == nil {
		t.Error("UpdateHealth() should error for nonexistent endpoint")
	}

	// Test nil update
	err = r.UpdateHealth(nil)
	if err == nil {
		t.Error("UpdateHealth(nil) should error")
	}
}

func TestRegistryRenew(t *testing.T) {
	r := NewRegistry()

	originalExpiry := time.Now().Add(1 * time.Minute)
	r.Register(&Info{ID: "ep1", ServiceID: "svc1", Host: "host1", ExpiresAt: originalExpiry})

	err := r.Renew("ep1", 5*time.Minute)
	if err != nil {
		t.Fatalf("Renew() error = %v", err)
	}

	ep, _ := r.Get("ep1")
	if !ep.ExpiresAt.After(originalExpiry) {
		t.Error("ExpiresAt should be extended")
	}

	// Test not found
	err = r.Renew("nonexistent", 5*time.Minute)
	if err == nil {
		t.Error("Renew() should error for nonexistent endpoint")
	}
}

func TestRegistryTouch(t *testing.T) {
	r := NewRegistry()

	r.Register(&Info{ID: "ep1", ServiceID: "svc1", Host: "host1"})

	ep, _ := r.Get("ep1")
	oldUpdated := ep.UpdatedAt

	time.Sleep(10 * time.Millisecond)

	err := r.Touch("ep1")
	if err != nil {
		t.Fatalf("Touch() error = %v", err)
	}

	ep, _ = r.Get("ep1")
	if !ep.UpdatedAt.After(oldUpdated) {
		t.Error("UpdatedAt should be refreshed")
	}

	// Test not found
	err = r.Touch("nonexistent")
	if err == nil {
		t.Error("Touch() should error for nonexistent endpoint")
	}
}

func TestRegistryClear(t *testing.T) {
	r := NewRegistry()

	r.Register(&Info{ID: "ep1", ServiceID: "svc1", Host: "host1"})
	r.Register(&Info{ID: "ep2", ServiceID: "svc2", Host: "host2"})

	r.Clear()

	if r.Count() != 0 {
		t.Errorf("After Clear(), count = %d, want 0", r.Count())
	}
}

func TestRegistryShutdown(t *testing.T) {
	r := NewRegistry()
	r.Register(&Info{ID: "ep1", ServiceID: "svc1", Host: "host1"})

	ctx := context.Background()
	err := r.Shutdown(ctx)
	if err != nil {
		t.Fatalf("Shutdown() error = %v", err)
	}

	// After shutdown, register should fail
	err = r.Register(&Info{ID: "ep2", ServiceID: "svc2", Host: "host2"})
	if err == nil {
		t.Error("Register() should fail after shutdown")
	}
}

func TestRegistryRegisterValidation(t *testing.T) {
	r := NewRegistry()

	// Nil endpoint
	err := r.Register(nil)
	if err == nil {
		t.Error("Register(nil) should error")
	}

	// Missing ID
	err = r.Register(&Info{ServiceID: "svc", Host: "host"})
	if err == nil {
		t.Error("Register without ID should error")
	}

	// Missing service ID
	err = r.Register(&Info{ID: "ep1", Host: "host"})
	if err == nil {
		t.Error("Register without ServiceID should error")
	}
}

// ============== Additional Manager Tests ==============

func TestManagerRegister(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	req := &RegistrationRequest{
		Endpoint: &Info{
			ID:        "ep1",
			ServiceID: "svc1",
			Host:      "localhost",
		},
	}

	err := manager.Register(ctx, req)
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	_, found := manager.Get("ep1")
	if !found {
		t.Error("Endpoint not found after Register")
	}
}

func TestManagerRegisterWithTTL(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	ep := &Info{
		ID:        "ep1",
		ServiceID: "svc1",
		Host:      "localhost",
	}

	err := manager.RegisterWithTTL(ctx, ep, 5*time.Minute)
	if err != nil {
		t.Fatalf("RegisterWithTTL() error = %v", err)
	}

	retrieved, _ := manager.Get("ep1")
	if retrieved.ExpiresAt.IsZero() {
		t.Error("ExpiresAt should be set")
	}
}

func TestManagerDeregister(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: "ep1", ServiceID: "svc1", Host: "host1"})

	err := manager.Deregister(ctx, &DeregistrationRequest{
		EndpointID: "ep1",
		Reason:     "test deregistration",
	})
	if err != nil {
		t.Fatalf("Deregister() error = %v", err)
	}

	_, found := manager.Get("ep1")
	if found {
		t.Error("Endpoint should be removed")
	}
}

func TestManagerDeregisterService(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: "ep1", ServiceID: "svc1", Host: "host1"})
	manager.RegisterEndpoint(ctx, &Info{ID: "ep2", ServiceID: "svc1", Host: "host2"})
	manager.RegisterEndpoint(ctx, &Info{ID: "ep3", ServiceID: "svc2", Host: "host3"})

	count, err := manager.DeregisterService(ctx, "svc1")
	if err != nil {
		t.Fatalf("DeregisterService() error = %v", err)
	}
	if count != 2 {
		t.Errorf("DeregisterService() returned count = %d, want 2", count)
	}

	if manager.Count() != 1 {
		t.Errorf("Count = %d, want 1", manager.Count())
	}
}

func TestManagerDeregisterInstance(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: "ep1", ServiceID: "svc1", Host: "host1", InstanceID: "inst1"})
	manager.RegisterEndpoint(ctx, &Info{ID: "ep2", ServiceID: "svc1", Host: "host2", InstanceID: "inst1"})
	manager.RegisterEndpoint(ctx, &Info{ID: "ep3", ServiceID: "svc2", Host: "host3", InstanceID: "inst2"})

	count, err := manager.DeregisterInstance(ctx, "inst1")
	if err != nil {
		t.Fatalf("DeregisterInstance() error = %v", err)
	}
	if count != 2 {
		t.Errorf("DeregisterInstance() returned count = %d, want 2", count)
	}

	if manager.Count() != 1 {
		t.Errorf("Count = %d, want 1", manager.Count())
	}
}

func TestManagerQuery(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: "ep1", ServiceID: "svc1", Host: "host1", Status: StatusHealthy})
	manager.RegisterEndpoint(ctx, &Info{ID: "ep2", ServiceID: "svc1", Host: "host2", Status: StatusHealthy})

	results := manager.Query(&Query{ServiceIDs: []string{"svc1"}})
	if len(results) != 2 {
		t.Errorf("Query() returned %d endpoints, want 2", len(results))
	}
}

func TestManagerGetByService(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: "ep1", ServiceID: "svc1", Host: "host1"})
	manager.RegisterEndpoint(ctx, &Info{ID: "ep2", ServiceID: "svc1", Host: "host2"})
	manager.RegisterEndpoint(ctx, &Info{ID: "ep3", ServiceID: "svc2", Host: "host3"})

	endpoints := manager.GetByService("svc1")
	if len(endpoints) != 2 {
		t.Errorf("GetByService() returned %d endpoints, want 2", len(endpoints))
	}
}

func TestManagerGetByInstance(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: "ep1", ServiceID: "svc1", Host: "host1", InstanceID: "inst1"})
	manager.RegisterEndpoint(ctx, &Info{ID: "ep2", ServiceID: "svc1", Host: "host2", InstanceID: "inst1"})

	endpoints := manager.GetByInstance("inst1")
	if len(endpoints) != 2 {
		t.Errorf("GetByInstance() returned %d endpoints, want 2", len(endpoints))
	}
}

func TestManagerGetHealthy(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: "ep1", ServiceID: "svc1", Host: "host1", Status: StatusHealthy})
	manager.RegisterEndpoint(ctx, &Info{ID: "ep2", ServiceID: "svc1", Host: "host2", Status: StatusUnhealthy})

	endpoints := manager.GetHealthy()
	if len(endpoints) != 1 {
		t.Errorf("GetHealthy() returned %d endpoints, want 1", len(endpoints))
	}
}

func TestManagerGetAvailable(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: "ep1", ServiceID: "svc1", Host: "host1", Status: StatusHealthy})
	manager.RegisterEndpoint(ctx, &Info{ID: "ep2", ServiceID: "svc1", Host: "host2", Status: StatusDegraded})
	manager.RegisterEndpoint(ctx, &Info{ID: "ep3", ServiceID: "svc1", Host: "host3", Status: StatusUnhealthy})

	endpoints := manager.GetAvailable()
	if len(endpoints) != 2 {
		t.Errorf("GetAvailable() returned %d endpoints, want 2", len(endpoints))
	}
}

func TestManagerUpdateStatus(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: "ep1", ServiceID: "svc1", Host: "host1", Status: StatusHealthy})

	err := manager.UpdateStatus(ctx, "ep1", StatusDegraded)
	if err != nil {
		t.Fatalf("UpdateStatus() error = %v", err)
	}

	ep, _ := manager.Get("ep1")
	if ep.Status != StatusDegraded {
		t.Errorf("Status = %v, want %v", ep.Status, StatusDegraded)
	}
}

func TestManagerRenew(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: "ep1", ServiceID: "svc1", Host: "host1"})

	err := manager.Renew(ctx, "ep1", 5*time.Minute)
	if err != nil {
		t.Fatalf("Renew() error = %v", err)
	}

	ep, _ := manager.Get("ep1")
	if ep.ExpiresAt.IsZero() {
		t.Error("ExpiresAt should be set after Renew")
	}
}

func TestManagerAll(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: "ep1", ServiceID: "svc1", Host: "host1"})
	manager.RegisterEndpoint(ctx, &Info{ID: "ep2", ServiceID: "svc2", Host: "host2"})

	all := manager.All()
	if len(all) != 2 {
		t.Errorf("All() returned %d endpoints, want 2", len(all))
	}
}

func TestManagerCount(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	if manager.Count() != 0 {
		t.Error("Count should be 0 initially")
	}

	manager.RegisterEndpoint(ctx, &Info{ID: "ep1", ServiceID: "svc1", Host: "host1"})
	manager.RegisterEndpoint(ctx, &Info{ID: "ep2", ServiceID: "svc1", Host: "host2"})

	if manager.Count() != 2 {
		t.Errorf("Count() = %d, want 2", manager.Count())
	}
}

func TestManagerRegistry(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	reg := manager.Registry()
	if reg == nil {
		t.Error("Registry() should not return nil")
	}
}

func TestManagerOnError(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	errorReceived := false
	manager.OnError(func(err error) {
		errorReceived = true
	})

	// Verify callback is set (can't easily trigger error without more setup)
	_ = errorReceived
}

func TestManagerErrors(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	errChan := manager.Errors()
	if errChan == nil {
		t.Error("Errors() should not return nil")
	}
}

// ============== Handler Tests ==============

func TestDefaultSubjects(t *testing.T) {
	subjects := DefaultSubjects()

	if subjects.Register != "endpoints.register" {
		t.Errorf("Register = %s, want endpoints.register", subjects.Register)
	}
	if subjects.Deregister != "endpoints.deregister" {
		t.Errorf("Deregister = %s, want endpoints.deregister", subjects.Deregister)
	}
	if subjects.Renew != "endpoints.renew" {
		t.Errorf("Renew = %s, want endpoints.renew", subjects.Renew)
	}
	if subjects.Query != "endpoints.query" {
		t.Errorf("Query = %s, want endpoints.query", subjects.Query)
	}
	if subjects.Discover != "endpoints.discover" {
		t.Errorf("Discover = %s, want endpoints.discover", subjects.Discover)
	}
	if subjects.Events != "endpoints.events" {
		t.Errorf("Events = %s, want endpoints.events", subjects.Events)
	}
	if subjects.HealthUpdate != "endpoints.health.update" {
		t.Errorf("HealthUpdate = %s, want endpoints.health.update", subjects.HealthUpdate)
	}
}

func TestTenantSubjects(t *testing.T) {
	subjects := TenantSubjects("tenant1")

	if subjects.Register != "tenant1.endpoints.register" {
		t.Errorf("Register = %s, want tenant1.endpoints.register", subjects.Register)
	}
	if subjects.Deregister != "tenant1.endpoints.deregister" {
		t.Errorf("Deregister = %s, want tenant1.endpoints.deregister", subjects.Deregister)
	}
}

func TestNewHandler(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	handler := NewHandler(manager)
	if handler == nil {
		t.Fatal("NewHandler() returned nil")
	}
}

func TestNewHandlerWithOptions(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	customSubjects := Subjects{
		Register:   "custom.register",
		Deregister: "custom.deregister",
	}

	handler := NewHandler(manager, WithSubjects(customSubjects))
	if handler == nil {
		t.Fatal("NewHandler() returned nil")
	}
}

func TestNewHandlerWithPublisher(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	handler := NewHandlerWithPublisher(manager, nil)
	if handler == nil {
		t.Fatal("NewHandlerWithPublisher() returned nil")
	}

	// With custom subjects
	subjects := TenantSubjects("tenant1")
	handler2 := NewHandlerWithPublisher(manager, nil, subjects)
	if handler2 == nil {
		t.Fatal("NewHandlerWithPublisher() with subjects returned nil")
	}
}

// ============== Client Tests ==============

func TestNewClient(t *testing.T) {
	client := NewClient(nil)
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}

	// With custom subjects
	subjects := TenantSubjects("tenant1")
	client2 := NewClient(nil, subjects)
	if client2 == nil {
		t.Fatal("NewClient() with subjects returned nil")
	}
}

func TestClientSetTimeout(t *testing.T) {
	client := NewClient(nil)
	client.SetTimeout(30 * time.Second)
	// Verify it doesn't panic
}

// ============== Health Checker Tests ==============

func TestNewHealthChecker(t *testing.T) {
	cfg := HealthCheckerConfig{
		Timeout:            5 * time.Second,
		HealthyThreshold:   2,
		UnhealthyThreshold: 3,
	}

	checker := NewHealthChecker(cfg)
	if checker == nil {
		t.Fatal("NewHealthChecker() returned nil")
	}
}

func TestNewHealthCheckerDefaults(t *testing.T) {
	// Test with zero values - should use defaults
	cfg := HealthCheckerConfig{}
	checker := NewHealthChecker(cfg)
	if checker == nil {
		t.Fatal("NewHealthChecker() returned nil")
	}
}

func TestHealthCheckConfig(t *testing.T) {
	cfg := HealthCheckConfig{
		Enabled:            true,
		Interval:           30 * time.Second,
		Timeout:            5 * time.Second,
		HealthyThreshold:   2,
		UnhealthyThreshold: 3,
		InitialDelay:       5 * time.Second,
	}

	if !cfg.Enabled {
		t.Error("Enabled should be true")
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("Interval = %v, want 30s", cfg.Interval)
	}
	if cfg.HealthyThreshold != 2 {
		t.Errorf("HealthyThreshold = %d, want 2", cfg.HealthyThreshold)
	}
}

// ============== Manager Configuration Tests ==============

func TestManagerConfigDefaults(t *testing.T) {
	manager, err := NewManager(ManagerConfig{})
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}
	if manager == nil {
		t.Fatal("NewManager() returned nil")
	}
}

func TestManagerWithHealthCheckEnabled(t *testing.T) {
	manager, err := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{
			Enabled:  true,
			Interval: 1 * time.Second,
			Timeout:  500 * time.Millisecond,
		},
	})
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}

	ctx := context.Background()
	manager.Shutdown(ctx)
}

// ============== OnUpdated Callback Test ==============

func TestRegistryOnUpdatedCallback(t *testing.T) {
	r := NewRegistry()

	updatedCalled := false
	r.OnUpdated(func(old, new *Info) {
		updatedCalled = true
	})

	ep := &Info{
		ID:        "ep1",
		ServiceID: "svc1",
		Host:      "localhost",
		Port:      8080,
	}

	r.Register(ep)

	// Update by re-registering with different port
	ep2 := &Info{
		ID:        "ep1",
		ServiceID: "svc1",
		Host:      "localhost",
		Port:      9090,
	}
	r.Register(ep2)

	time.Sleep(10 * time.Millisecond)
	if !updatedCalled {
		t.Error("OnUpdated callback not called")
	}
}

// ============== Query with Offset Test ==============

func TestRegistryQueryWithOffset(t *testing.T) {
	r := NewRegistry()

	for i := 0; i < 10; i++ {
		r.Register(&Info{
			ID:        "ep" + itoa(i),
			ServiceID: "svc",
			Host:      "host",
		})
	}

	results := r.Query(&Query{Limit: 3, Offset: 2})
	if len(results) != 3 {
		t.Errorf("Query with offset: got %d, want 3", len(results))
	}
}

// ============== SetPublisher Test ==============

func TestManagerSetPublisher(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	manager.SetPublisher(nil)
	// Verify it doesn't panic
}

func TestHandlerSetPublisher(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	handler := NewHandler(manager)
	handler.SetPublisher(nil)
	// Verify it doesn't panic
}

// ============== Additional Benchmarks ==============

func BenchmarkRegistryDeregister(b *testing.B) {
	r := NewRegistry()

	for i := 0; i < b.N; i++ {
		ep := &Info{
			ID:        "ep" + itoa(i),
			ServiceID: "svc",
			Host:      "localhost",
		}
		r.Register(ep)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r.Deregister("ep" + itoa(i))
	}
}

func BenchmarkRegistryGetByService(b *testing.B) {
	r := NewRegistry()

	for i := 0; i < 1000; i++ {
		r.Register(&Info{
			ID:        "ep" + itoa(i),
			ServiceID: "svc" + itoa(i%10),
			Host:      "localhost",
		})
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r.GetByService("svc" + itoa(i%10))
	}
}

func BenchmarkRegistryUpdateStatus(b *testing.B) {
	r := NewRegistry()

	for i := 0; i < 1000; i++ {
		r.Register(&Info{
			ID:        "ep" + itoa(i),
			ServiceID: "svc",
			Host:      "localhost",
			Status:    StatusHealthy,
		})
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		r.UpdateStatus("ep"+itoa(i%1000), StatusDegraded)
	}
}

// ============== Mock Publisher for Testing ==============

type mockPublisher struct {
	publishErr   error
	requestReply *core.Message
	requestErr   error
}

func (m *mockPublisher) Publish(ctx context.Context, subject string, msg *core.Message) error {
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

// ============== Handler Tests with Mock Publisher ==============

func TestWithPublisherOption(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	pub := &mockPublisher{}
	handler := NewHandler(manager, WithPublisher(pub))
	if handler == nil {
		t.Fatal("NewHandler() with WithPublisher returned nil")
	}
}

func TestHandlerHandleRegister(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	handler := NewHandler(manager)

	// Create a registration request message
	req := RegistrationRequest{
		Endpoint: &Info{
			ID:        "ep1",
			ServiceID: "svc1",
			Host:      "localhost",
			Port:      8080,
		},
	}
	data, _ := json.Marshal(req)
	msg := core.NewMessage(data)

	ctx := context.Background()
	handler.HandleRegister(ctx, msg)

	// Verify endpoint was registered
	_, found := manager.Get("ep1")
	if !found {
		t.Error("Endpoint should be registered")
	}
}

func TestHandlerHandleRegisterInvalid(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	handler := NewHandler(manager)

	// Invalid JSON
	msg := core.NewMessage([]byte("invalid json"))
	ctx := context.Background()
	handler.HandleRegister(ctx, msg)
	// Should not panic
}

func TestHandlerHandleDeregister(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: "ep1", ServiceID: "svc1", Host: "localhost"})

	handler := NewHandler(manager)

	req := DeregistrationRequest{
		EndpointID: "ep1",
		Reason:     "test",
	}
	data, _ := json.Marshal(req)
	msg := core.NewMessage(data)

	handler.HandleDeregister(ctx, msg)

	// Verify endpoint was deregistered
	_, found := manager.Get("ep1")
	if found {
		t.Error("Endpoint should be deregistered")
	}
}

func TestHandlerHandleDeregisterInvalid(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	handler := NewHandler(manager)
	msg := core.NewMessage([]byte("invalid"))
	ctx := context.Background()
	handler.HandleDeregister(ctx, msg)
	// Should not panic
}

func TestHandlerHandleRenew(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: "ep1", ServiceID: "svc1", Host: "localhost"})

	handler := NewHandler(manager)

	req := RenewRequest{
		EndpointID: "ep1",
		TTL:        5 * time.Minute,
	}
	data, _ := json.Marshal(req)
	msg := core.NewMessage(data)

	handler.HandleRenew(ctx, msg)
	// Should not panic
}

func TestHandlerHandleRenewInvalid(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	handler := NewHandler(manager)
	msg := core.NewMessage([]byte("invalid"))
	ctx := context.Background()
	handler.HandleRenew(ctx, msg)
	// Should not panic
}

func TestHandlerHandleQuery(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: "ep1", ServiceID: "svc1", Host: "localhost", Status: StatusHealthy})

	handler := NewHandler(manager)

	// Query uses Query struct from types.go
	req := Query{
		ServiceIDs: []string{"svc1"},
	}
	data, _ := json.Marshal(req)
	msg := core.NewMessage(data)

	handler.HandleQuery(ctx, msg)
	// Should not panic
}

func TestHandlerHandleQueryInvalid(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	handler := NewHandler(manager)
	msg := core.NewMessage([]byte("invalid"))
	ctx := context.Background()
	handler.HandleQuery(ctx, msg)
	// Should not panic
}

func TestHandlerHandleDiscover(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: "ep1", ServiceID: "svc1", Host: "localhost", Status: StatusHealthy})

	handler := NewHandler(manager)

	req := DiscoverRequest{
		ServiceID: "svc1",
	}
	data, _ := json.Marshal(req)
	msg := core.NewMessage(data)

	handler.HandleDiscover(ctx, msg)
	// Should not panic
}

func TestHandlerHandleDiscoverInvalid(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	handler := NewHandler(manager)
	msg := core.NewMessage([]byte("invalid"))
	ctx := context.Background()
	handler.HandleDiscover(ctx, msg)
	// Should not panic
}

func TestHandlerHandleHealthUpdate(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: "ep1", ServiceID: "svc1", Host: "localhost", Status: StatusHealthy})

	handler := NewHandler(manager)

	update := HealthUpdate{
		EndpointID: "ep1",
		Status:     StatusDegraded,
		Message:    "health degraded",
		CheckedAt:  time.Now(),
	}
	data, _ := json.Marshal(update)
	msg := core.NewMessage(data)

	handler.HandleHealthUpdate(ctx, msg)

	// Verify status updated
	ep, _ := manager.Get("ep1")
	if ep.Status != StatusDegraded {
		t.Errorf("Status = %v, want %v", ep.Status, StatusDegraded)
	}
}

func TestHandlerHandleHealthUpdateInvalid(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	handler := NewHandler(manager)
	msg := core.NewMessage([]byte("invalid"))
	ctx := context.Background()
	handler.HandleHealthUpdate(ctx, msg)
	// Should not panic
}

// ============== Manager HTTP Handler Tests ==============

func TestManagerHTTPHandleRegistration(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	req := RegistrationRequest{
		Endpoint: &Info{
			ID:        "ep1",
			ServiceID: "svc1",
			Host:      "localhost",
		},
		TTL: 5 * time.Minute,
	}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	manager.HandleRegistration(rr, httpReq)

	if rr.Code != http.StatusOK && rr.Code != http.StatusCreated {
		t.Errorf("HandleRegistration returned %d", rr.Code)
	}
}

func TestManagerHTTPHandleRegistrationWrongMethod(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	httpReq := httptest.NewRequest(http.MethodGet, "/register", nil)
	rr := httptest.NewRecorder()

	manager.HandleRegistration(rr, httpReq)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405, got %d", rr.Code)
	}
}

func TestManagerHTTPHandleDeregistration(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: "ep1", ServiceID: "svc1", Host: "localhost"})

	req := DeregistrationRequest{
		EndpointID: "ep1",
		Reason:     "test",
	}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/deregister", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	manager.HandleDeregistration(rr, httpReq)

	if rr.Code != http.StatusOK {
		t.Errorf("HandleDeregistration returned %d", rr.Code)
	}
}

func TestManagerHTTPHandleDeregistrationWrongMethod(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	httpReq := httptest.NewRequest(http.MethodGet, "/deregister", nil)
	rr := httptest.NewRecorder()

	manager.HandleDeregistration(rr, httpReq)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405, got %d", rr.Code)
	}
}

func TestManagerHTTPHandleQuery(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: "ep1", ServiceID: "svc1", Host: "localhost", Status: StatusHealthy})

	httpReq := httptest.NewRequest(http.MethodGet, "/query?service_id=svc1", nil)
	rr := httptest.NewRecorder()

	manager.HandleQuery(rr, httpReq)

	if rr.Code != http.StatusOK {
		t.Errorf("HandleQuery returned %d", rr.Code)
	}
}

func TestManagerHTTPHandleQueryWrongMethod(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	httpReq := httptest.NewRequest(http.MethodDelete, "/query", nil)
	rr := httptest.NewRecorder()

	manager.HandleQuery(rr, httpReq)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405, got %d", rr.Code)
	}
}

func TestManagerHTTPHandleRenew(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: "ep1", ServiceID: "svc1", Host: "localhost"})

	req := RenewRequest{
		EndpointID: "ep1",
		TTL:        10 * time.Minute,
	}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, "/renew", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	manager.HandleRenew(rr, httpReq)

	if rr.Code != http.StatusOK {
		t.Errorf("HandleRenew returned %d", rr.Code)
	}
}

func TestManagerHTTPHandleRenewWrongMethod(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	httpReq := httptest.NewRequest(http.MethodGet, "/renew", nil)
	rr := httptest.NewRecorder()

	manager.HandleRenew(rr, httpReq)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405, got %d", rr.Code)
	}
}

// ============== Default Config Tests ==============

func TestDefaultHealthCheckConfig(t *testing.T) {
	cfg := DefaultHealthCheckConfig()

	if !cfg.Enabled {
		t.Error("Default health check should be enabled")
	}
	if cfg.Interval == 0 {
		t.Error("Default interval should not be zero")
	}
	if cfg.Timeout == 0 {
		t.Error("Default timeout should not be zero")
	}
}

// ============== Registry Deregister with Reason Callback ==============

func TestRegistryDeregisterWithReasonCallback(t *testing.T) {
	r := NewRegistry()

	deregisteredCalled := false
	var gotReason string
	r.OnDeregistered(func(ep *Info, reason string) {
		deregisteredCalled = true
		gotReason = reason
	})

	r.Register(&Info{ID: "ep1", ServiceID: "svc1", Host: "localhost"})
	r.DeregisterWithReason("ep1", "shutdown")

	time.Sleep(10 * time.Millisecond)
	if !deregisteredCalled {
		t.Error("OnDeregistered callback not called")
	}
	if gotReason != "shutdown" {
		t.Errorf("Reason = %s, want shutdown", gotReason)
	}
}

// ============== Additional Discovery Options ==============

func TestDiscoverOptionsChaining(t *testing.T) {
	q := &Query{}

	// Test chaining multiple options
	WithTags("api", "v1")(q)
	WithVersion("1.0.0")(q)
	WithProtocol(ProtocolHTTP)(q)
	WithRegion("us-east-1")(q)
	WithLimit(10)(q)
	OnlyHealthy()(q)

	if len(q.Tags) != 2 {
		t.Errorf("Tags = %v, want 2 tags", q.Tags)
	}
	if q.Version != "1.0.0" {
		t.Errorf("Version = %s, want 1.0.0", q.Version)
	}
	if q.Protocol != ProtocolHTTP {
		t.Errorf("Protocol = %v, want HTTP", q.Protocol)
	}
	if q.Region != "us-east-1" {
		t.Errorf("Region = %s, want us-east-1", q.Region)
	}
	if q.Limit != 10 {
		t.Errorf("Limit = %d, want 10", q.Limit)
	}
	if !q.OnlyHealthy {
		t.Error("OnlyHealthy should be true")
	}
}

// ============== Expired Endpoint Tests ==============

func TestInfoIsExpiredEdgeCases(t *testing.T) {
	// Just at expiry time (edge case)
	info := &Info{ExpiresAt: time.Now()}
	// Should either be expired or not based on exact timing
	_ = info.IsExpired()
}

// ============== Registration Request with TTL ==============

func TestManagerRegisterValidationError(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()

	// Missing endpoint
	req := &RegistrationRequest{}
	err := manager.Register(ctx, req)
	if err == nil {
		t.Error("Register should fail without endpoint")
	}

	// Missing ID
	req = &RegistrationRequest{
		Endpoint: &Info{ServiceID: "svc", Host: "host"},
	}
	err = manager.Register(ctx, req)
	if err == nil {
		t.Error("Register should fail without ID")
	}
}

// ============== Deregister Validation ==============

func TestManagerDeregisterValidation(t *testing.T) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()

	// Missing endpoint ID
	err := manager.Deregister(ctx, &DeregistrationRequest{})
	if err == nil {
		t.Error("Deregister should fail without endpoint ID")
	}
}

// ============== Additional Benchmarks ==============

func BenchmarkManagerRegister(b *testing.B) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		manager.RegisterEndpoint(ctx, &Info{
			ID:        "ep" + itoa(i),
			ServiceID: "svc",
			Host:      "localhost",
		})
	}
}

func BenchmarkManagerDiscover(b *testing.B) {
	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})

	ctx := context.Background()
	for i := 0; i < 100; i++ {
		manager.RegisterEndpoint(ctx, &Info{
			ID:        "ep" + itoa(i),
			ServiceID: "svc",
			Host:      "localhost",
			Status:    StatusHealthy,
		})
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		manager.Discover("svc")
	}
}
