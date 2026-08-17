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

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

/* ========================================================================================================
   TEST CONSTANTS
   ======================================================================================================== */

var (
	testServiceID1   = "svc1"
	testServiceID2   = "svc2"
	testServiceID3   = "svc3"
	testServiceName  = "MyService"
	testServiceTest  = "Test"
	testServiceTest1 = "Test1"
	testServiceTest2 = "Test2"
	testServiceOther = "Other"
	testInstanceID1  = "inst1"
	testInstanceID2  = "inst2"
	testInstanceID3  = "inst3"
	testInstanceID4  = "inst4"
	testHost         = "localhost"
	testTenantID1    = "tenant1"
	testTenantID2    = "tenant2"
	testVersion      = "1.0.0"
	testAPIVersion   = "v1"
	testRegisterOp   = "register"
	testReplySubject = "reply"
	testHealthPath   = "/health"
	testHealthURL    = "http://localhost:8080/health"
	testInvalidJSON  = "invalid"
	testProtocolHTTP = "http"
	testShutdown     = "shutdown"

	testRegistryCfg = RegistryConfig{ExpirationEnabled: false}
	testManagerCfg  = ManagerConfig{HealthCheck: HealthCheckConfig{Enabled: false}}

	// endpoints
	registerEndpoint         = "/register"
	deregisterEndpoint       = "/deregister"
	servicesEndpoint         = "/services"
	renewEndpoint            = "/renew"
	registerInstanceEndpoint = "/register-instance"
)

/* ========================================================================================================
   MOCK PUBLISHER
   ======================================================================================================== */

type mockPublisher struct {
	mu           sync.Mutex
	publishErr   error
	requestReply *nats.Msg
	requestErr   error
	published    []struct {
		subject string
		msg     *nats.Msg
	}
}

func (m *mockPublisher) Publish(ctx context.Context, subject string, msg *nats.Msg) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.published = append(m.published, struct {
		subject string
		msg     *nats.Msg
	}{subject, msg})
	return m.publishErr
}

func (m *mockPublisher) Request(ctx context.Context, subject string, msg *nats.Msg) (*nats.Msg, error) {
	if m.requestErr != nil {
		return nil, m.requestErr
	}
	return m.requestReply, nil
}

func (m *mockPublisher) Close(ctx context.Context) error {
	return nil
}

/* ========================================================================================================
   MOCK SUBSCRIBER
   ======================================================================================================== */

type mockSubscriber struct {
	subscribeErr error
	subs         []*mockSubscription
}

type mockSubscription struct {
	subject string
	valid   bool
}

func (m *mockSubscription) Subject() string    { return m.subject }
func (m *mockSubscription) Unsubscribe() error { m.valid = false; return nil }
func (m *mockSubscription) Drain() error       { return nil }
func (m *mockSubscription) IsValid() bool      { return m.valid }

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

/* ========================================================================================================
   HELPER FUNCTIONS
   ======================================================================================================== */

func newTestRegistry(t *testing.T) *Registry {
	t.Helper()
	r := NewRegistry(testRegistryCfg)
	t.Cleanup(func() { r.Shutdown(context.Background()) })
	return r
}

func newTestManager(t *testing.T) *Manager {
	t.Helper()
	m, err := NewManager(testManagerCfg)
	assert.New(t).NoError(err)
	t.Cleanup(func() { m.Shutdown(context.Background()) })
	return m
}

func newTestManagerWithPub(t *testing.T, pub core.Publisher) *Manager {
	t.Helper()
	m, err := NewManager(ManagerConfig{Publisher: pub, HealthCheck: HealthCheckConfig{Enabled: false}})
	assert.New(t).NoError(err)
	t.Cleanup(func() { m.Shutdown(context.Background()) })
	return m
}

func makeNATSMsg(t *testing.T, v any, reply string) *nats.Msg {
	t.Helper()
	data, err := json.Marshal(v)
	assert.New(t).NoError(err)
	msg := &nats.Msg{Data: data, Header: make(nats.Header)}
	msg.Reply = reply
	return msg
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

/* ========================================================================================================
   ERROR TESTS
   ======================================================================================================== */

func TestServiceError(t *testing.T) {
	testCases := []struct {
		name     string
		err      ServiceError
		expected string
	}{
		{
			name:     "with service name and op",
			err:      ServiceError{ServiceID: testServiceID1, ServiceName: testServiceName, Op: testRegisterOp, Err: ErrServiceExists},
			expected: "service MyService (svc1): register: service already exists",
		},
		{
			name:     "with op only",
			err:      ServiceError{ServiceID: testServiceID1, Op: testRegisterOp, Err: ErrServiceExists},
			expected: "service svc1: register: service already exists",
		},
		{
			name:     "minimal",
			err:      ServiceError{ServiceID: testServiceID1, Err: ErrServiceExists},
			expected: "service svc1: service already exists",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, tc.err.Error())
			assertions.Equal(tc.err.Err, tc.err.Unwrap())
		})
	}
}

func TestNewServiceError(t *testing.T) {
	assertions := assert.New(t)
	err := NewServiceError(testServiceID1, testServiceName, testRegisterOp, ErrServiceExists)
	assertions.Equal(testServiceID1, err.ServiceID)
	assertions.Equal(testServiceName, err.ServiceName)
	assertions.Equal(testRegisterOp, err.Op)
}

func TestInstanceError(t *testing.T) {
	testCases := []struct {
		name     string
		err      InstanceError
		expected string
	}{
		{
			name:     "with service ID and op",
			err:      InstanceError{InstanceID: testInstanceID1, ServiceID: testServiceID1, Op: testRegisterOp, Err: ErrInstanceExists},
			expected: "instance inst1 (service svc1): register: instance already exists",
		},
		{
			name:     "with op only",
			err:      InstanceError{InstanceID: testInstanceID1, Op: testRegisterOp, Err: ErrInstanceExists},
			expected: "instance inst1: register: instance already exists",
		},
		{
			name:     "minimal",
			err:      InstanceError{InstanceID: testInstanceID1, Err: ErrInstanceExists},
			expected: "instance inst1: instance already exists",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, tc.err.Error())
			assertions.Equal(tc.err.Err, tc.err.Unwrap())
		})
	}
}

func TestNewInstanceError(t *testing.T) {
	assertions := assert.New(t)
	err := NewInstanceError(testInstanceID1, testServiceID1, testRegisterOp, ErrInstanceExists)
	assertions.Equal(testInstanceID1, err.InstanceID)
	assertions.Equal(testServiceID1, err.ServiceID)
	assertions.Equal(testRegisterOp, err.Op)
}

func TestRegistrationError(t *testing.T) {
	testCases := []struct {
		name     string
		err      RegistrationError
		expected string
	}{
		{
			name:     "with instance ID and reason",
			err:      RegistrationError{ServiceID: testServiceID1, InstanceID: testInstanceID1, Reason: "validation failed", Err: ErrInvalidService},
			expected: "registration failed for inst1@svc1: validation failed: invalid service",
		},
		{
			name:     "without instance ID",
			err:      RegistrationError{ServiceID: testServiceID1, Reason: "validation failed", Err: ErrInvalidService},
			expected: "registration failed for svc1: validation failed: invalid service",
		},
		{
			name:     "without reason",
			err:      RegistrationError{ServiceID: testServiceID1, InstanceID: testInstanceID1, Err: ErrInvalidService},
			expected: "registration failed for inst1@svc1: invalid service",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, tc.err.Error())
			assertions.Equal(tc.err.Err, tc.err.Unwrap())
		})
	}
}

func TestHealthCheckError(t *testing.T) {
	assertions := assert.New(t)

	err := HealthCheckError{
		ServiceID:        testServiceID1,
		InstanceID:       testInstanceID1,
		URL:              testHealthURL,
		Err:              errors.New("connection refused"),
		ConsecutiveFails: 3,
	}

	expected := "health check failed for inst1@svc1 at http://localhost:8080/health (consecutive fails: 3): connection refused"
	assertions.Equal(expected, err.Error())
	assertions.Equal("connection refused", err.Unwrap().Error())

	err2 := HealthCheckError{
		ServiceID:        testServiceID1,
		URL:              testHealthURL,
		Err:              errors.New("timeout"),
		ConsecutiveFails: 1,
	}
	assertions.True(contains(err2.Error(), testServiceID1))
}

func TestIsServiceError(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{"ServiceError", ServiceError{ServiceID: testServiceID1, Err: ErrServiceNotFound}, true},
		{"InstanceError", InstanceError{InstanceID: testInstanceID1, Err: ErrInstanceNotFound}, true},
		{"RegistrationError", RegistrationError{ServiceID: testServiceID1, Err: ErrInvalidService}, true},
		{"HealthCheckError", HealthCheckError{ServiceID: testServiceID1, Err: errors.New("fail")}, true},
		{"regular error", errors.New("random error"), false},
		{"nil error", nil, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, IsServiceError(tc.err))
		})
	}
}

func TestIsNotFound(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{"ErrServiceNotFound", ErrServiceNotFound, true},
		{"ErrInstanceNotFound", ErrInstanceNotFound, true},
		{"wrapped ErrServiceNotFound", ServiceError{Err: ErrServiceNotFound}, true},
		{"other error", ErrServiceExists, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, IsNotFound(tc.err))
		})
	}
}

func TestIsExists(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{"ErrServiceExists", ErrServiceExists, true},
		{"ErrInstanceExists", ErrInstanceExists, true},
		{"wrapped ErrServiceExists", ServiceError{Err: ErrServiceExists}, true},
		{"other error", ErrServiceNotFound, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, IsExists(tc.err))
		})
	}
}

/* ========================================================================================================
   STATUS TESTS
   ======================================================================================================== */

func TestStatusString(t *testing.T) {
	testCases := []struct {
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

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, tc.status.String())
		})
	}
}

func TestStatusIsAvailable(t *testing.T) {
	testCases := []struct {
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

	for _, tc := range testCases {
		t.Run(tc.status.String(), func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, tc.status.IsAvailable())
		})
	}
}

func TestStatusIsHealthy(t *testing.T) {
	testCases := []struct {
		status   Status
		expected bool
	}{
		{StatusRunning, true},
		{StatusDegraded, false},
		{StatusStarting, false},
		{StatusStopping, false},
		{StatusFailed, false},
	}

	for _, tc := range testCases {
		t.Run(tc.status.String(), func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, tc.status.IsHealthy())
		})
	}
}

/* ========================================================================================================
   INFO TESTS
   ======================================================================================================== */

func TestInfoClone(t *testing.T) {
	assertions := assert.New(t)

	var nilInfo *Info
	assertions.Nil(nilInfo.Clone())

	original := &Info{
		ID:           testServiceID1,
		Name:         testServiceName,
		Dependencies: []string{"dep1", "dep2"},
		Capabilities: []string{"cap1"},
		Tags:         []string{"tag1", "tag2"},
		Metadata:     map[string]string{"key": "value"},
		Instances:    []*Instance{{ID: testInstanceID1, Host: testHost}},
	}

	clone := original.Clone()
	assertions.Equal(original.ID, clone.ID)

	original.Dependencies[0] = "modified"
	assertions.NotEqual("modified", clone.Dependencies[0])

	original.Capabilities[0] = "modified"
	assertions.NotEqual("modified", clone.Capabilities[0])

	original.Tags[0] = "modified"
	assertions.NotEqual("modified", clone.Tags[0])

	original.Metadata["key"] = "modified"
	assertions.NotEqual("modified", clone.Metadata["key"])

	original.Instances[0].Host = "modified"
	assertions.NotEqual("modified", clone.Instances[0].Host)
}

func TestInfoHasCapability(t *testing.T) {
	assertions := assert.New(t)
	info := &Info{Capabilities: []string{"read", "write", "admin"}}

	assertions.True(info.HasCapability("read"))
	assertions.False(info.HasCapability("delete"))
}

func TestInfoHasTag(t *testing.T) {
	assertions := assert.New(t)
	info := &Info{Tags: []string{"production", "critical"}}

	assertions.True(info.HasTag("production"))
	assertions.False(info.HasTag("staging"))
}

func TestInfoMetadata(t *testing.T) {
	assertions := assert.New(t)
	info := &Info{}

	assertions.Empty(info.GetMetadata("key"))
	info.SetMetadata("key", "value")
	assertions.Equal("value", info.GetMetadata("key"))
	assertions.Empty(info.GetMetadata("missing"))
}

func TestInfoHealthyInstanceCount(t *testing.T) {
	assertions := assert.New(t)
	info := &Info{
		Instances: []*Instance{
			{ID: testInstanceID1, Status: StatusRunning},
			{ID: testInstanceID2, Status: StatusDegraded},
			{ID: testInstanceID3, Status: StatusRunning},
			{ID: testInstanceID4, Status: StatusFailed},
		},
	}
	assertions.Equal(2, info.HealthyInstanceCount())
}

func TestInfoAvailableInstanceCount(t *testing.T) {
	assertions := assert.New(t)
	info := &Info{
		Instances: []*Instance{
			{ID: testInstanceID1, Status: StatusRunning},
			{ID: testInstanceID2, Status: StatusDegraded},
			{ID: testInstanceID3, Status: StatusRunning},
			{ID: testInstanceID4, Status: StatusFailed},
		},
	}
	assertions.Equal(3, info.AvailableInstanceCount())
}

/* ========================================================================================================
   INSTANCE TESTS
   ======================================================================================================== */

func TestInstanceClone(t *testing.T) {
	assertions := assert.New(t)

	var nilInst *Instance
	assertions.Nil(nilInst.Clone())

	original := &Instance{
		ID:       testInstanceID1,
		Host:     testHost,
		Metadata: map[string]string{"key": "value"},
	}
	clone := original.Clone()
	assertions.Equal(original.ID, clone.ID)

	original.Metadata["key"] = "modified"
	assertions.NotEqual("modified", clone.Metadata["key"])
}

func TestInstanceAddress(t *testing.T) {
	testCases := []struct {
		name     string
		instance *Instance
		expected string
	}{
		{name: "with port", instance: &Instance{Host: testHost, Port: 8080}, expected: "localhost:8080"},
		{name: "without port", instance: &Instance{Host: testHost, Port: 0}, expected: testHost},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, tc.instance.Address())
		})
	}
}

func TestInstanceIsExpired(t *testing.T) {
	testCases := []struct {
		name      string
		expiresAt time.Time
		expected  bool
	}{
		{"zero time", time.Time{}, false},
		{"future time", time.Now().Add(time.Hour), false},
		{"past time", time.Now().Add(-time.Hour), true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			inst := &Instance{ExpiresAt: tc.expiresAt}
			assertions.Equal(tc.expected, inst.IsExpired())
		})
	}
}

func TestInstanceIsHealthy(t *testing.T) {
	testCases := []struct {
		status   Status
		expected bool
	}{
		{StatusRunning, true},
		{StatusDegraded, false},
		{StatusFailed, false},
	}

	for _, tc := range testCases {
		t.Run(tc.status.String(), func(t *testing.T) {
			assertions := assert.New(t)
			inst := &Instance{Status: tc.status}
			assertions.Equal(tc.expected, inst.IsHealthy())
		})
	}
}

func TestInstanceMetadata(t *testing.T) {
	assertions := assert.New(t)
	inst := &Instance{}

	assertions.Empty(inst.GetMetadata("key"))
	inst.SetMetadata("key", "value")
	assertions.Equal("value", inst.GetMetadata("key"))
}

/* ========================================================================================================
   REQUEST VALIDATION TESTS
   ======================================================================================================== */

func TestRegistrationRequestValidate(t *testing.T) {
	testCases := []struct {
		name    string
		req     *RegistrationRequest
		wantErr error
	}{
		{"nil service", &RegistrationRequest{Service: nil}, ErrInvalidService},
		{"missing-service-ID", &RegistrationRequest{Service: &Info{Name: testServiceName}}, ErrMissingServiceID},
		{"missing service name", &RegistrationRequest{Service: &Info{ID: testServiceID1}}, ErrMissingServiceName},
		{"valid request", &RegistrationRequest{Service: &Info{ID: testServiceID1, Name: testServiceName}}, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			err := tc.req.Validate()
			if tc.wantErr != nil {
				assertions.ErrorIs(err, tc.wantErr)
			} else {
				assertions.NoError(err)
			}
		})
	}
}

func TestInstanceRegistrationRequestValidate(t *testing.T) {
	testCases := []struct {
		name    string
		req     *InstanceRegistrationRequest
		wantErr error
	}{
		{"nil instance", &InstanceRegistrationRequest{Instance: nil}, ErrInvalidInstance},
		{"missing instance ID", &InstanceRegistrationRequest{Instance: &Instance{ServiceID: testServiceID1, Host: testHost}}, ErrMissingInstanceID},
		{"missing service Id", &InstanceRegistrationRequest{Instance: &Instance{ID: testInstanceID1, Host: testHost}}, ErrMissingServiceID},
		{"missing host", &InstanceRegistrationRequest{Instance: &Instance{ID: testInstanceID1, ServiceID: testServiceID1}}, ErrMissingHost},
		{"valid-request", &InstanceRegistrationRequest{Instance: &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost}}, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			err := tc.req.Validate()
			if tc.wantErr != nil {
				assertions.ErrorIs(err, tc.wantErr)
			} else {
				assertions.NoError(err)
			}
		})
	}
}

func TestDeregistrationRequestValidate(t *testing.T) {
	testCases := []struct {
		name    string
		req     *DeregistrationRequest
		wantErr error
	}{
		{"missing service id", &DeregistrationRequest{}, ErrMissingServiceID},
		{"valid_request", &DeregistrationRequest{ServiceID: testServiceID1}, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			err := tc.req.Validate()
			if tc.wantErr != nil {
				assertions.ErrorIs(err, tc.wantErr)
			} else {
				assertions.NoError(err)
			}
		})
	}
}

/* ========================================================================================================
   QUERY TESTS
   ======================================================================================================== */

func TestQueryMatches(t *testing.T) {
	service := &Info{
		ID:           testServiceID1,
		Name:         testServiceName,
		Type:         ServiceTypeAPI,
		Version:      testVersion,
		APIVersion:   testAPIVersion,
		Status:       StatusRunning,
		TenantID:     testTenantID1,
		Capabilities: []string{"read", "write"},
		Tags:         []string{"production"},
		Metadata:     map[string]string{"env": "prod"},
		Instances:    []*Instance{{ID: testInstanceID1}},
	}

	testCases := []struct {
		name     string
		svc      *Info
		query    *Query
		expected bool
	}{
		{"nil-service", nil, &Query{}, false},
		{"empty query", service, &Query{}, true},
		{"matching IDs", service, &Query{IDs: []string{testServiceID1, testServiceID2}}, true},
		{"non-matching IDs", service, &Query{IDs: []string{"other"}}, false},
		{"matching names", service, &Query{Names: []string{testServiceName}}, true},
		{"non-matching names", service, &Query{Names: []string{testServiceOther}}, false},
		{"matching types", service, &Query{Types: []ServiceType{ServiceTypeAPI}}, true},
		{"non-matching types", service, &Query{Types: []ServiceType{ServiceTypeWorker}}, false},
		{"matching statuses", service, &Query{Statuses: []Status{StatusRunning}}, true},
		{"non-matching statuses", service, &Query{Statuses: []Status{StatusStopped}}, false},
		{"only healthy - healthy", service, &Query{OnlyHealthy: true}, true},
		{"only available - available", service, &Query{OnlyAvailable: true}, true},
		{"matching capabilities", service, &Query{Capabilities: []string{"read"}}, true},
		{"non-matching capabilities", service, &Query{Capabilities: []string{"admin"}}, false},
		{"matching tags", service, &Query{Tags: []string{"production"}}, true},
		{"non-matching tags", service, &Query{Tags: []string{"staging"}}, false},
		{"matching metadata", service, &Query{Metadata: map[string]string{"env": "prod"}}, true},
		{"non-matching metadata", service, &Query{Metadata: map[string]string{"env": "dev"}}, false},
		{"matching version", service, &Query{Version: testVersion}, true},
		{"non-matching version", service, &Query{Version: "2.0.0"}, false},
		{"matching API version", service, &Query{APIVersion: testAPIVersion}, true},
		{"non-matching API version", service, &Query{APIVersion: "v2"}, false},
		{"matching tenant", service, &Query{TenantID: testTenantID1}, true},
		{"non-matching tenant", service, &Query{TenantID: testTenantID2}, false},
		{"has instances - has", service, &Query{HasInstances: true}, true},
		{"min instances - satisfied", service, &Query{MinInstances: 1}, true},
		{"min instances - not satisfied", service, &Query{MinInstances: 5}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, tc.query.Matches(tc.svc))
		})
	}

	t.Run("unhealthy and unavailable service", func(t *testing.T) {
		assertions := assert.New(t)
		unhealthy := &Info{ID: testServiceID2, Name: "Unhealthy", Status: StatusFailed}
		assertions.False((&Query{OnlyHealthy: true}).Matches(unhealthy))
		assertions.False((&Query{OnlyAvailable: true}).Matches(unhealthy))
	})

	t.Run("empty instances", func(t *testing.T) {
		assertions := assert.New(t)
		noInstances := &Info{ID: testServiceID3, Name: "NoInstances", Status: StatusRunning, Instances: []*Instance{}}
		assertions.False((&Query{HasInstances: true}).Matches(noInstances))
	})
}

/* ========================================================================================================
   EVENT TESTS
   ======================================================================================================== */

func TestEventToMessage(t *testing.T) {
	assertions := assert.New(t)

	event := &Event{
		Type:       EventServiceRegistered,
		ServiceID:  testServiceID1,
		InstanceID: testInstanceID1,
		OldStatus:  StatusUnknown,
		NewStatus:  StatusRunning,
		Timestamp:  time.Now(),
	}

	msg, err := event.ToMessage()
	assertions.NoError(err)
	assertions.Equal(string(EventServiceRegistered), msg.Header.Get("event-type"))
	assertions.Equal(testServiceID1, msg.Header.Get("service-id"))
	assertions.Equal(testInstanceID1, msg.Header.Get("instance-id"))

	event2 := &Event{Type: EventServiceRegistered, ServiceID: testServiceID1, Timestamp: time.Now()}
	msg2, err := event2.ToMessage()
	assertions.NoError(err)
	assertions.Empty(msg2.Header.Get("instance-id"))
}

func TestEventFromMessage(t *testing.T) {
	assertions := assert.New(t)

	event := &Event{
		Type:       EventServiceRegistered,
		ServiceID:  testServiceID1,
		InstanceID: testInstanceID1,
		OldStatus:  StatusUnknown,
		NewStatus:  StatusRunning,
		Timestamp:  time.Now(),
	}

	data, _ := json.Marshal(event)
	msg := &nats.Msg{Data: data, Header: make(nats.Header)}

	parsed, err := EventFromMessage(msg)
	assertions.NoError(err)
	assertions.Equal(event.Type, parsed.Type)
	assertions.Equal(event.ServiceID, parsed.ServiceID)

	badMsg := &nats.Msg{Data: []byte(testInvalidJSON), Header: make(nats.Header)}
	_, err = EventFromMessage(badMsg)
	assertions.Error(err)
}

/* ========================================================================================================
   REGISTRY TESTS
   ======================================================================================================== */

func TestNewRegistry(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()
	assertions.NotNil(r)
	assertions.NotNil(r.services)
	r.Shutdown(context.Background())
}

func TestDefaultRegistryConfig(t *testing.T) {
	assertions := assert.New(t)
	cfg := DefaultRegistryConfig()
	assertions.Equal(time.Minute, cfg.CleanupInterval)
	assertions.True(cfg.ExpirationEnabled)
}

func TestRegistryWithConfig(t *testing.T) {
	assertions := assert.New(t)
	cfg := RegistryConfig{CleanupInterval: 5 * time.Minute, ExpirationEnabled: false}
	r := NewRegistry(cfg)
	assertions.Equal(5*time.Minute, r.cleanupInterval)
	r.Shutdown(context.Background())
}

func TestRegistryRegisterService(t *testing.T) {
	testCases := []struct {
		name    string
		service *Info
		wantErr error
	}{
		{"nil_service", nil, ErrInvalidService},
		{"missing-ID", &Info{Name: testServiceTest}, ErrMissingServiceID},
		{"missing name", &Info{ID: testServiceID1}, ErrMissingServiceName},
		{"valid registration", &Info{ID: testServiceID1, Name: testServiceTest, Status: StatusRunning}, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			r := newTestRegistry(t)

			err := r.RegisterService(tc.service)
			if tc.wantErr != nil {
				assertions.ErrorIs(err, tc.wantErr)
			} else {
				assertions.NoError(err)
			}
		})
	}

	t.Run("update existing service", func(t *testing.T) {
		assertions := assert.New(t)
		r := newTestRegistry(t)

		svc := &Info{ID: testServiceID1, Name: testServiceTest, Status: StatusRunning}
		assertions.NoError(r.RegisterService(svc))

		svc.Status = StatusDegraded
		assertions.NoError(r.RegisterService(svc))

		got, exists := r.GetService(testServiceID1)
		assertions.True(exists)
		assertions.Equal(StatusDegraded, got.Status)
	})
}

func TestRegistryRegisterServiceWithCallbacks(t *testing.T) {
	assertions := assert.New(t)
	r := newTestRegistry(t)

	var (
		registered    bool
		updated       bool
		statusChanged bool
		mu            sync.Mutex
	)

	r.OnServiceRegistered(func(svc *Info) { mu.Lock(); registered = true; mu.Unlock() })
	r.OnServiceUpdated(func(old, new *Info) { mu.Lock(); updated = true; mu.Unlock() })
	r.OnServiceStatusChange(func(svc *Info, oldStatus, newStatus Status) { mu.Lock(); statusChanged = true; mu.Unlock() })

	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest, Status: StatusRunning})
	time.Sleep(10 * time.Millisecond)
	mu.Lock()
	assertions.True(registered)
	mu.Unlock()

	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest, Status: StatusDegraded})
	time.Sleep(10 * time.Millisecond)
	mu.Lock()
	assertions.True(updated)
	assertions.True(statusChanged)
	mu.Unlock()
}

func TestRegistryDeregisterService(t *testing.T) {
	testCases := []struct {
		name    string
		id      string
		wantErr error
	}{
		{"missing_ID", "", ErrMissingServiceID},
		{"not found", "unknown", ErrServiceNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			r := newTestRegistry(t)
			assertions.ErrorIs(r.DeregisterService(tc.id), tc.wantErr)
		})
	}

	t.Run("successful deregistration with callback", func(t *testing.T) {
		assertions := assert.New(t)
		r := newTestRegistry(t)
		r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})
		r.RegisterInstance(&Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})

		var deregistered bool
		r.OnServiceDeregistered(func(svc *Info, reason string) { deregistered = true })

		assertions.NoError(r.DeregisterServiceWithReason(testServiceID1, testShutdown))
		time.Sleep(10 * time.Millisecond)
		assertions.True(deregistered)
		assertions.False(r.ServiceExists(testServiceID1))
		assertions.False(r.InstanceExists(testInstanceID1))
	})
}

func TestRegistryRegisterInstance(t *testing.T) {
	testCases := []struct {
		name    string
		inst    *Instance
		wantErr error
	}{
		{"nil instance", nil, ErrInvalidInstance},
		{"missing ID", &Instance{ServiceID: testServiceID1, Host: testHost}, ErrMissingInstanceID},
		{"missing_service_ID", &Instance{ID: testInstanceID1, Host: testHost}, ErrMissingServiceID},
		{"missing host", &Instance{ID: testInstanceID1, ServiceID: testServiceID1}, ErrMissingHost},
		{"service not found", &Instance{ID: testInstanceID1, ServiceID: "unknown", Host: testHost}, ErrServiceNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			r := newTestRegistry(t)
			r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})
			assertions.ErrorIs(r.RegisterInstance(tc.inst), tc.wantErr)
		})
	}

	t.Run("successful registration with callback", func(t *testing.T) {
		assertions := assert.New(t)
		r := newTestRegistry(t)
		r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})

		var registered bool
		r.OnInstanceRegistered(func(inst *Instance) { registered = true })

		inst := &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost, Status: StatusRunning}
		assertions.NoError(r.RegisterInstance(inst))

		time.Sleep(10 * time.Millisecond)
		assertions.True(registered)
	})

	t.Run("update existing instance triggers status change", func(t *testing.T) {
		assertions := assert.New(t)
		r := newTestRegistry(t)
		r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})

		inst := &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost, Status: StatusRunning}
		assertions.NoError(r.RegisterInstance(inst))

		var statusChanged bool
		r.OnInstanceStatusChange(func(inst *Instance, oldStatus, newStatus Status) { statusChanged = true })

		inst.Status = StatusDegraded
		assertions.NoError(r.RegisterInstance(inst))
		time.Sleep(10 * time.Millisecond)
		assertions.True(statusChanged)
	})
}

func TestRegistryRegisterInstanceWithTTL(t *testing.T) {
	assertions := assert.New(t)
	r := newTestRegistry(t)
	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})

	inst := &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost}
	assertions.NoError(r.RegisterInstanceWithTTL(inst, 5*time.Minute))

	got, exists := r.GetInstance(testInstanceID1)
	assertions.True(exists)
	assertions.False(got.ExpiresAt.IsZero())
}

func TestRegistryDeregisterInstance(t *testing.T) {
	testCases := []struct {
		name    string
		id      string
		wantErr error
	}{
		{"missing ID", "", ErrMissingInstanceID},
		{"not found", "unknown", ErrInstanceNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			r := newTestRegistry(t)
			assertions.ErrorIs(r.DeregisterInstance(tc.id), tc.wantErr)
		})
	}

	t.Run("successful deregistration with callback", func(t *testing.T) {
		assertions := assert.New(t)
		r := newTestRegistry(t)
		r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})
		r.RegisterInstance(&Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})

		var deregistered bool
		r.OnInstanceDeregistered(func(inst *Instance, reason string) { deregistered = true })

		assertions.NoError(r.DeregisterInstanceWithReason(testInstanceID1, testShutdown))
		time.Sleep(10 * time.Millisecond)
		assertions.True(deregistered)
		assertions.False(r.InstanceExists(testInstanceID1))
	})
}

func TestRegistryGetService(t *testing.T) {
	assertions := assert.New(t)
	r := newTestRegistry(t)

	_, exists := r.GetService("unknown")
	assertions.False(exists)

	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})
	svc, exists := r.GetService(testServiceID1)
	assertions.True(exists)
	assertions.Equal(testServiceID1, svc.ID)
}

func TestRegistryGetServiceByName(t *testing.T) {
	assertions := assert.New(t)
	r := newTestRegistry(t)

	assertions.Empty(r.GetServiceByName("unknown"))

	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})
	r.RegisterService(&Info{ID: testServiceID2, Name: testServiceTest})
	assertions.Len(r.GetServiceByName(testServiceTest), 2)
}

func TestRegistryGetInstance(t *testing.T) {
	assertions := assert.New(t)
	r := newTestRegistry(t)

	_, exists := r.GetInstance("unknown")
	assertions.False(exists)

	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})
	r.RegisterInstance(&Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})

	inst, exists := r.GetInstance(testInstanceID1)
	assertions.True(exists)
	assertions.Equal(testInstanceID1, inst.ID)

	expiredInst := &Instance{
		ID: testInstanceID2, ServiceID: testServiceID1, Host: testHost,
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	r.RegisterInstance(expiredInst)
	_, exists = r.GetInstance(testInstanceID2)
	assertions.False(exists, "expired instance should not be returned")
}

func TestRegistryGetInstances(t *testing.T) {
	assertions := assert.New(t)
	r := newTestRegistry(t)

	assertions.Empty(r.GetInstances("unknown"))

	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})
	r.RegisterInstance(&Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})
	r.RegisterInstance(&Instance{ID: testInstanceID2, ServiceID: testServiceID1, Host: testHost})
	assertions.Len(r.GetInstances(testServiceID1), 2)
}

func TestRegistryGetHealthyInstances(t *testing.T) {
	assertions := assert.New(t)
	r := newTestRegistry(t)

	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})
	r.RegisterInstance(&Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost, Status: StatusRunning})
	r.RegisterInstance(&Instance{ID: testInstanceID2, ServiceID: testServiceID1, Host: testHost, Status: StatusDegraded})
	r.RegisterInstance(&Instance{ID: testInstanceID3, ServiceID: testServiceID1, Host: testHost, Status: StatusFailed})

	assertions.Len(r.GetHealthyInstances(testServiceID1), 1)
}

func TestRegistryGetAvailableInstances(t *testing.T) {
	assertions := assert.New(t)
	r := newTestRegistry(t)

	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})
	r.RegisterInstance(&Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost, Status: StatusRunning})
	r.RegisterInstance(&Instance{ID: testInstanceID2, ServiceID: testServiceID1, Host: testHost, Status: StatusDegraded})
	r.RegisterInstance(&Instance{ID: testInstanceID3, ServiceID: testServiceID1, Host: testHost, Status: StatusFailed})

	assertions.Len(r.GetAvailableInstances(testServiceID1), 2)
}

func TestRegistryQueryServices(t *testing.T) {
	assertions := assert.New(t)
	r := newTestRegistry(t)

	r.RegisterService(&Info{ID: testServiceID1, Name: "A-Service", Status: StatusRunning, Tags: []string{"prod"}})
	r.RegisterService(&Info{ID: testServiceID2, Name: "B-Service", Status: StatusDegraded, Tags: []string{"dev"}})
	r.RegisterService(&Info{ID: testServiceID3, Name: "C-Service", Status: StatusFailed})

	assertions.Len(r.QueryServices(nil), 3)
	assertions.Len(r.QueryServices(&Query{OnlyHealthy: true}), 1)
	assertions.Len(r.QueryServices(&Query{Limit: 1}), 1)
	assertions.Len(r.QueryServices(&Query{Offset: 2}), 1)
	assertions.Empty(r.QueryServices(&Query{Offset: 100}))
}

func TestRegistryAllServices(t *testing.T) {
	assertions := assert.New(t)
	r := newTestRegistry(t)

	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest1})
	r.RegisterService(&Info{ID: testServiceID2, Name: testServiceTest2})
	assertions.Len(r.AllServices(), 2)
}

func TestRegistryAllInstances(t *testing.T) {
	assertions := assert.New(t)
	r := newTestRegistry(t)

	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})
	r.RegisterInstance(&Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})
	r.RegisterInstance(&Instance{ID: testInstanceID2, ServiceID: testServiceID1, Host: testHost})
	assertions.Len(r.AllInstances(), 2)
}

func TestRegistryServiceCount(t *testing.T) {
	assertions := assert.New(t)
	r := newTestRegistry(t)

	assertions.Equal(0, r.ServiceCount())
	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})
	assertions.Equal(1, r.ServiceCount())
}

func TestRegistryInstanceCount(t *testing.T) {
	assertions := assert.New(t)
	r := newTestRegistry(t)

	assertions.Equal(0, r.InstanceCount())
	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})
	r.RegisterInstance(&Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})
	assertions.Equal(1, r.InstanceCount())
}

func TestRegistryServiceExists(t *testing.T) {
	assertions := assert.New(t)
	r := newTestRegistry(t)

	assertions.False(r.ServiceExists(testServiceID1))
	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})
	assertions.True(r.ServiceExists(testServiceID1))
}

func TestRegistryInstanceExists(t *testing.T) {
	assertions := assert.New(t)
	r := newTestRegistry(t)

	assertions.False(r.InstanceExists(testInstanceID1))
	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})
	r.RegisterInstance(&Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})
	assertions.True(r.InstanceExists(testInstanceID1))
}

func TestRegistryUpdateServiceStatus(t *testing.T) {
	assertions := assert.New(t)
	r := newTestRegistry(t)

	assertions.ErrorIs(r.UpdateServiceStatus("unknown", StatusRunning), ErrServiceNotFound)

	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest, Status: StatusStarting})

	var statusChanged bool
	r.OnServiceStatusChange(func(svc *Info, oldStatus, newStatus Status) { statusChanged = true })

	assertions.NoError(r.UpdateServiceStatus(testServiceID1, StatusRunning))
	time.Sleep(10 * time.Millisecond)
	assertions.True(statusChanged)

	svc, _ := r.GetService(testServiceID1)
	assertions.Equal(StatusRunning, svc.Status)
}

func TestRegistryUpdateInstanceStatus(t *testing.T) {
	assertions := assert.New(t)
	r := newTestRegistry(t)

	assertions.ErrorIs(r.UpdateInstanceStatus("unknown", StatusRunning), ErrInstanceNotFound)

	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})
	r.RegisterInstance(&Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost, Status: StatusStarting})

	assertions.NoError(r.UpdateInstanceStatus(testInstanceID1, StatusRunning))
	assertions.NoError(r.UpdateInstanceStatus(testInstanceID1, StatusFailed))

	inst, _ := r.GetInstance(testInstanceID1)
	assertions.Equal(1, inst.ConsecutiveFails)
}

func TestRegistryUpdateInstanceHealth(t *testing.T) {
	assertions := assert.New(t)
	r := newTestRegistry(t)

	assertions.ErrorIs(r.UpdateInstanceHealth(nil), ErrInvalidInstance)
	assertions.ErrorIs(r.UpdateInstanceHealth(&HealthUpdate{}), ErrInvalidInstance)
	assertions.ErrorIs(r.UpdateInstanceHealth(&HealthUpdate{InstanceID: "unknown"}), ErrInstanceNotFound)

	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})
	r.RegisterInstance(&Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})

	assertions.NoError(r.UpdateInstanceHealth(&HealthUpdate{InstanceID: testInstanceID1, Status: StatusRunning, CheckedAt: time.Now()}))
	assertions.NoError(r.UpdateInstanceHealth(&HealthUpdate{InstanceID: testInstanceID1, Status: StatusFailed, CheckedAt: time.Now()}))
}

func TestRegistryRenewInstance(t *testing.T) {
	assertions := assert.New(t)
	r := newTestRegistry(t)

	assertions.ErrorIs(r.RenewInstance("unknown", time.Hour), ErrInstanceNotFound)

	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})
	r.RegisterInstance(&Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})
	assertions.NoError(r.RenewInstance(testInstanceID1, time.Hour))

	inst, _ := r.GetInstance(testInstanceID1)
	assertions.False(inst.ExpiresAt.IsZero())
}

func TestRegistryShutdown(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	assertions.NoError(r.Shutdown(ctx))
}

func TestRegistryShutdownTimeout(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	assertions.ErrorIs(r.Shutdown(ctx), context.Canceled)
}

func TestRegistryClear(t *testing.T) {
	assertions := assert.New(t)
	r := newTestRegistry(t)

	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})
	r.RegisterInstance(&Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})

	r.Clear()
	assertions.Equal(0, r.ServiceCount())
	assertions.Equal(0, r.InstanceCount())
}

func TestRegistryShutdownPreventsRegistration(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry(testRegistryCfg)
	r.cancel()

	assertions.ErrorIs(r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest}), ErrRegistryShutdown)

	r2 := NewRegistry(testRegistryCfg)
	r2.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})
	r2.cancel()

	assertions.ErrorIs(r2.RegisterInstance(&Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost}), ErrRegistryShutdown)
}

/* ========================================================================================================
   SUBJECTS TESTS
   ======================================================================================================== */

func TestDefaultSubjects(t *testing.T) {
	assertions := assert.New(t)
	s := DefaultSubjects()
	assertions.Equal("services.register", s.RegisterService)
	assertions.Equal("services.deregister", s.DeregisterService)
}

func TestTenantSubjects(t *testing.T) {
	assertions := assert.New(t)
	s := TenantSubjects(testTenantID1)
	assertions.Equal("tenant1.services.register", s.RegisterService)
	assertions.Equal("tenant1.health.update", s.HealthUpdate)
}

/* ========================================================================================================
   HANDLER TESTS
   ======================================================================================================== */

func TestNewHandler(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)

	h := NewHandler(m)
	assertions.NotNil(h)
	assertions.Equal(m, h.manager)
}

func TestNewHandlerWithOptions(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	pub := &mockPublisher{}
	subjects := TenantSubjects(testTenantID1)

	h := NewHandler(m, WithPublisher(pub), WithSubjects(subjects))
	assertions.Equal(pub, h.publisher)
	assertions.Equal(subjects.RegisterService, h.subjects.RegisterService)
}

func TestNewHandlerWithPublisher(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	pub := &mockPublisher{}

	h := NewHandlerWithPublisher(m, pub)
	assertions.Equal(pub, h.publisher)

	subjects := TenantSubjects(testTenantID1)
	h2 := NewHandlerWithPublisher(m, pub, subjects)
	assertions.Equal(subjects.RegisterService, h2.subjects.RegisterService)
}

func TestHandlerSetPublisher(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	h := NewHandler(m)
	pub := &mockPublisher{}
	h.SetPublisher(pub)
	assertions.Equal(pub, h.publisher)
}

/* ========================================================================================================
   CLIENT TESTS
   ======================================================================================================== */

func TestNewClient(t *testing.T) {
	assertions := assert.New(t)
	pub := &mockPublisher{}
	c := NewClient(pub)
	assertions.NotNil(c)
	assertions.Equal(10*time.Second, c.timeout)

	subjects := TenantSubjects(testTenantID1)
	c2 := NewClient(pub, subjects)
	assertions.Equal(subjects.RegisterService, c2.subjects.RegisterService)
}

func TestClientSetTimeout(t *testing.T) {
	assertions := assert.New(t)
	c := NewClient(&mockPublisher{})
	c.SetTimeout(30 * time.Second)
	assertions.Equal(30*time.Second, c.timeout)
}

/* ========================================================================================================
   MANAGER TESTS
   ======================================================================================================== */

func TestNewManager(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	assertions.NotNil(m.registry)
}

func TestNewManagerWithConfig(t *testing.T) {
	assertions := assert.New(t)
	registry := newTestRegistry(t)
	pub := &mockPublisher{}

	m, err := NewManager(ManagerConfig{
		Registry:        registry,
		Publisher:       pub,
		HealthCheck:     HealthCheckConfig{Enabled: false},
		ErrorBufferSize: 50,
	})
	assertions.NoError(err)
	defer m.Shutdown(context.Background())

	assertions.Equal(registry, m.registry)
	assertions.Equal(pub, m.publisher)
}

func TestDefaultHealthCheckConfig(t *testing.T) {
	assertions := assert.New(t)
	cfg := DefaultHealthCheckConfig()
	assertions.True(cfg.Enabled)
	assertions.Equal(30*time.Second, cfg.Interval)
	assertions.Equal(5*time.Second, cfg.Timeout)
}

func TestDefaultPublishSubjects(t *testing.T) {
	assertions := assert.New(t)
	s := DefaultPublishSubjects()
	assertions.Equal("services.registered", s.ServiceRegistered)
}

func TestManagerRegisterService(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()

	assertions.Error(m.RegisterService(ctx, &RegistrationRequest{}))

	req := &RegistrationRequest{Service: &Info{ID: testServiceID1, Name: testServiceTest}}
	assertions.NoError(m.RegisterService(ctx, req))

	svc, exists := m.GetService(testServiceID1)
	assertions.True(exists)
	assertions.Equal(testServiceID1, svc.ID)
}

func TestManagerRegisterServiceDirect(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	assertions.NoError(m.RegisterServiceDirect(context.Background(), &Info{ID: testServiceID1, Name: testServiceTest}))
}

func TestManagerDeregisterService(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})

	assertions.Error(m.DeregisterService(ctx, &DeregistrationRequest{}))
	assertions.NoError(m.DeregisterService(ctx, &DeregistrationRequest{ServiceID: testServiceID1, Reason: testShutdown}))
}

func TestManagerDeregisterServiceDirect(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	m.RegisterServiceDirect(context.Background(), &Info{ID: testServiceID1, Name: testServiceTest})
	assertions.NoError(m.DeregisterServiceDirect(context.Background(), testServiceID1))
}

func TestManagerRegisterInstance(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})

	assertions.Error(m.RegisterInstance(ctx, &InstanceRegistrationRequest{}))

	assertions.NoError(m.RegisterInstance(ctx, &InstanceRegistrationRequest{
		Instance: &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost},
	}))

	assertions.NoError(m.RegisterInstance(ctx, &InstanceRegistrationRequest{
		Instance: &Instance{ID: testInstanceID2, ServiceID: testServiceID1, Host: testHost},
		TTL:      5 * time.Minute,
	}))
}

func TestManagerRegisterInstanceDirect(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})
	assertions.NoError(m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost}))
}

func TestManagerRegisterInstanceWithTTL(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})
	assertions.NoError(m.RegisterInstanceWithTTL(ctx, &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost}, 5*time.Minute))
}

func TestManagerDeregisterInstance(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})
	assertions.NoError(m.DeregisterInstance(ctx, testInstanceID1))
}

func TestManagerGetServiceByName(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID2, Name: testServiceTest})
	assertions.Len(m.GetServiceByName(testServiceTest), 2)
}

func TestManagerGetInstance(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})

	inst, exists := m.GetInstance(testInstanceID1)
	assertions.True(exists)
	assertions.Equal(testInstanceID1, inst.ID)
}

func TestManagerGetInstances(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost, Status: StatusRunning})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID2, ServiceID: testServiceID1, Host: testHost, Status: StatusDegraded})
	assertions.Len(m.GetInstances(testServiceID1), 2)
}

func TestManagerGetHealthyInstances(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost, Status: StatusRunning})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID2, ServiceID: testServiceID1, Host: testHost, Status: StatusDegraded})
	assertions.Len(m.GetHealthyInstances(testServiceID1), 1)
}

func TestManagerGetAvailableInstances(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost, Status: StatusRunning})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID2, ServiceID: testServiceID1, Host: testHost, Status: StatusDegraded})
	assertions.Len(m.GetAvailableInstances(testServiceID1), 2)
}

func TestManagerQueryServices(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest, Status: StatusRunning})
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID2, Name: testServiceOther, Status: StatusDegraded})
	assertions.Len(m.QueryServices(&Query{OnlyHealthy: true}), 1)
}

func TestManagerDiscover(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest, Status: StatusRunning, Type: ServiceTypeAPI})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost, Status: StatusRunning})

	assertions.Len(m.Discover(testServiceTest), 1)
	assertions.Len(m.Discover(testServiceTest, WithType(ServiceTypeAPI), OnlyHealthy()), 1)
}

func TestDiscoverOptions(t *testing.T) {
	assertions := assert.New(t)
	q := &Query{}

	WithType(ServiceTypeAPI)(q)
	assertions.Equal([]ServiceType{ServiceTypeAPI}, q.Types)

	WithCapabilities("cap1", "cap2")(q)
	assertions.Len(q.Capabilities, 2)

	WithTags("tag1")(q)
	assertions.Len(q.Tags, 1)

	WithVersion(testVersion)(q)
	assertions.Equal(testVersion, q.Version)

	OnlyHealthy()(q)
	assertions.True(q.OnlyHealthy)
	assertions.False(q.OnlyAvailable)

	WithMinInstances(3)(q)
	assertions.Equal(3, q.MinInstances)

	WithLimit(10)(q)
	assertions.Equal(10, q.Limit)
}

func TestManagerDiscoverInstances(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost, Status: StatusRunning})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID2, ServiceID: testServiceID1, Host: testHost, Status: StatusDegraded})

	assertions.Len(m.DiscoverInstances(testServiceID1, true), 1)
	assertions.Len(m.DiscoverInstances(testServiceID1, false), 2)
}

func TestManagerRenewInstance(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})
	assertions.NoError(m.RenewInstance(ctx, testInstanceID1, time.Hour))
}

func TestManagerUpdateServiceStatus(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest, Status: StatusStarting})
	assertions.NoError(m.UpdateServiceStatus(ctx, testServiceID1, StatusRunning))
}

func TestManagerUpdateInstanceStatus(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost, Status: StatusStarting})
	assertions.NoError(m.UpdateInstanceStatus(ctx, testInstanceID1, StatusRunning))
}

func TestManagerAllServices(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest1})
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID2, Name: testServiceTest2})
	assertions.Len(m.AllServices(), 2)
}

func TestManagerAllInstances(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID2, ServiceID: testServiceID1, Host: testHost})
	assertions.Len(m.AllInstances(), 2)
}

func TestManagerServiceCount(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	m.RegisterServiceDirect(context.Background(), &Info{ID: testServiceID1, Name: testServiceTest})
	assertions.Equal(1, m.ServiceCount())
}

func TestManagerInstanceCount(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})
	assertions.Equal(1, m.InstanceCount())
}

func TestManagerRegistry(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	assertions.NotNil(m.Registry())
}

func TestManagerSetPublisher(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	pub := &mockPublisher{}
	m.SetPublisher(pub)
	assertions.Equal(pub, m.publisher)
}

func TestManagerErrors(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	assertions.NotNil(m.Errors())
}

func TestManagerOnError(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)

	var called bool
	m.OnError(func(err error) { called = true })
	m.dispatchError(errors.New("test error"))
	assertions.True(called)
}

func TestManagerShutdown(t *testing.T) {
	assertions := assert.New(t)
	m, _ := NewManager(testManagerCfg)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	assertions.NoError(m.Shutdown(ctx))
}

/* ========================================================================================================
   HEALTH CHECKER TESTS
   ======================================================================================================== */

func TestNewHealthChecker(t *testing.T) {
	assertions := assert.New(t)
	hc := NewHealthChecker(HealthCheckerConfig{})
	assertions.NotNil(hc)
	assertions.Equal(5*time.Second, hc.timeout)
	assertions.Equal(2, hc.healthyThreshold)
	assertions.Equal(3, hc.unhealthyThreshold)
}

func TestHealthCheckerCheck(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		wantStatus Status
	}{
		{"200 OK", http.StatusOK, StatusRunning},
		{"503 unavailable", http.StatusServiceUnavailable, StatusMaintenance},
		{"500 error", http.StatusInternalServerError, StatusFailed},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			hc := NewHealthChecker(HealthCheckerConfig{Timeout: time.Second, HealthyThreshold: 1, UnhealthyThreshold: 1})

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
			}))
			defer server.Close()

			inst := &Instance{
				ID: testInstanceID1, ServiceID: testServiceID1,
				Host: server.Listener.Addr().String(), Protocol: testProtocolHTTP, HealthCheckPath: "/",
			}

			result := hc.Check(context.Background(), inst)
			assertions.Equal(tc.wantStatus, result.Status)
		})
	}

	t.Run("connection failure", func(t *testing.T) {
		assertions := assert.New(t)
		hc := NewHealthChecker(HealthCheckerConfig{Timeout: time.Second, HealthyThreshold: 1, UnhealthyThreshold: 1})
		inst := &Instance{
			ID: testInstanceID1, ServiceID: testServiceID1,
			Host: "localhost:99999", Protocol: testProtocolHTTP, HealthCheckPath: testHealthPath,
		}
		result := hc.Check(context.Background(), inst)
		assertions.Equal(StatusFailed, result.Status)
	})
}

func TestHealthCheckerDefaultProtocolAndPath(t *testing.T) {
	assertions := assert.New(t)
	hc := NewHealthChecker(HealthCheckerConfig{Timeout: time.Second, HealthyThreshold: 1, UnhealthyThreshold: 1})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != testHealthPath {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	inst := &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: server.Listener.Addr().String()}

	result := hc.Check(context.Background(), inst)
	assertions.Equal(StatusRunning, result.Status)
}

/* ========================================================================================================
   HTTP HANDLER TESTS
   ======================================================================================================== */

func TestManagerHTTPHandleServiceRegistration(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)

	req := RegistrationRequest{Service: &Info{ID: testServiceID1, Name: testServiceTest}}
	body, _ := json.Marshal(req)

	rr := httptest.NewRecorder()
	m.HandleServiceRegistration(rr, httptest.NewRequest(http.MethodPost, registerEndpoint, bytes.NewReader(body)))
	assertions.Equal(http.StatusCreated, rr.Code)

	rr = httptest.NewRecorder()
	m.HandleServiceRegistration(rr, httptest.NewRequest(http.MethodGet, registerEndpoint, nil))
	assertions.Equal(http.StatusMethodNotAllowed, rr.Code)

	rr = httptest.NewRecorder()
	m.HandleServiceRegistration(rr, httptest.NewRequest(http.MethodPost, registerEndpoint, bytes.NewReader([]byte(testInvalidJSON))))
	assertions.Equal(http.StatusBadRequest, rr.Code)
}

func TestManagerHTTPHandleServiceDeregistration(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	m.RegisterServiceDirect(context.Background(), &Info{ID: testServiceID1, Name: testServiceTest})

	req := DeregistrationRequest{ServiceID: testServiceID1}
	body, _ := json.Marshal(req)

	rr := httptest.NewRecorder()
	m.HandleServiceDeregistration(rr, httptest.NewRequest(http.MethodDelete, deregisterEndpoint, bytes.NewReader(body)))
	assertions.Equal(http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	m.HandleServiceDeregistration(rr, httptest.NewRequest(http.MethodDelete, deregisterEndpoint, bytes.NewReader(body)))
	assertions.Equal(http.StatusNotFound, rr.Code)

	rr = httptest.NewRecorder()
	m.HandleServiceDeregistration(rr, httptest.NewRequest(http.MethodGet, deregisterEndpoint, nil))
	assertions.Equal(http.StatusMethodNotAllowed, rr.Code)
}

func TestManagerHTTPHandleInstanceRegistration(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	m.RegisterServiceDirect(context.Background(), &Info{ID: testServiceID1, Name: testServiceTest})

	req := InstanceRegistrationRequest{Instance: &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost}}
	body, _ := json.Marshal(req)

	rr := httptest.NewRecorder()
	m.HandleInstanceRegistration(rr, httptest.NewRequest(http.MethodPost, registerInstanceEndpoint, bytes.NewReader(body)))
	assertions.Equal(http.StatusCreated, rr.Code)

	rr = httptest.NewRecorder()
	m.HandleInstanceRegistration(rr, httptest.NewRequest(http.MethodGet, registerInstanceEndpoint, nil))
	assertions.Equal(http.StatusMethodNotAllowed, rr.Code)
}

func TestManagerHTTPHandleServiceQuery(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	m.RegisterServiceDirect(context.Background(), &Info{ID: testServiceID1, Name: testServiceTest, Status: StatusRunning})

	rr := httptest.NewRecorder()
	m.HandleServiceQuery(rr, httptest.NewRequest(http.MethodGet, "/services?name=Test&healthy=true&available=true", nil))
	assertions.Equal(http.StatusOK, rr.Code)

	query := Query{OnlyHealthy: true}
	body, _ := json.Marshal(query)
	rr = httptest.NewRecorder()
	m.HandleServiceQuery(rr, httptest.NewRequest(http.MethodPost, servicesEndpoint, bytes.NewReader(body)))
	assertions.Equal(http.StatusOK, rr.Code)

	rr = httptest.NewRecorder()
	m.HandleServiceQuery(rr, httptest.NewRequest(http.MethodPut, servicesEndpoint, nil))
	assertions.Equal(http.StatusMethodNotAllowed, rr.Code)
}

func TestManagerHTTPHandleInstanceRenew(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})

	reqBody := struct {
		InstanceID string        `json:"instance_id"`
		TTL        time.Duration `json:"ttl"`
	}{InstanceID: testInstanceID1, TTL: 5 * time.Minute}
	body, _ := json.Marshal(reqBody)

	rr := httptest.NewRecorder()
	m.HandleInstanceRenew(rr, httptest.NewRequest(http.MethodPost, renewEndpoint, bytes.NewReader(body)))
	assertions.Equal(http.StatusOK, rr.Code)

	reqBody.InstanceID = "unknown"
	body, _ = json.Marshal(reqBody)
	rr = httptest.NewRecorder()
	m.HandleInstanceRenew(rr, httptest.NewRequest(http.MethodPost, renewEndpoint, bytes.NewReader(body)))
	assertions.Equal(http.StatusNotFound, rr.Code)

	rr = httptest.NewRecorder()
	m.HandleInstanceRenew(rr, httptest.NewRequest(http.MethodGet, renewEndpoint, nil))
	assertions.Equal(http.StatusMethodNotAllowed, rr.Code)
}

/* ========================================================================================================
   HELPER TESTS
   ======================================================================================================== */

func TestItoa(t *testing.T) {
	testCases := []struct {
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

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, itoa(tc.input))
		})
	}
}

/* ========================================================================================================
   HANDLER MESSAGE TESTS
   ======================================================================================================== */

func TestHandlerHandleRegisterService(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	pub := &mockPublisher{}
	h := NewHandler(m, WithPublisher(pub))

	req := &RegistrationRequest{Service: &Info{ID: testServiceID1, Name: testServiceTest}}
	msg := makeNATSMsg(t, req, testReplySubject)
	assertions.NoError(h.HandleRegisterService(context.Background(), msg))
	assertions.True(m.registry.ServiceExists(testServiceID1))

	badMsg := &nats.Msg{Data: []byte(testInvalidJSON), Header: make(nats.Header), Reply: testReplySubject}
	assertions.Error(h.HandleRegisterService(context.Background(), badMsg))

	req2 := &RegistrationRequest{Service: &Info{ID: testServiceID2, Name: testServiceTest2}}
	msg2 := makeNATSMsg(t, req2, "")
	assertions.NoError(h.HandleRegisterService(context.Background(), msg2))
}

func TestHandlerHandleDeregisterService(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	m.RegisterServiceDirect(context.Background(), &Info{ID: testServiceID1, Name: testServiceTest})

	h := NewHandler(m, WithPublisher(&mockPublisher{}))

	req := &DeregistrationRequest{ServiceID: testServiceID1}
	msg := makeNATSMsg(t, req, testReplySubject)
	assertions.NoError(h.HandleDeregisterService(context.Background(), msg))

	badMsg := &nats.Msg{Data: []byte(testInvalidJSON), Header: make(nats.Header), Reply: testReplySubject}
	assertions.Error(h.HandleDeregisterService(context.Background(), badMsg))
}

func TestHandlerHandleQueryServices(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	m.RegisterServiceDirect(context.Background(), &Info{ID: testServiceID1, Name: testServiceTest, Status: StatusRunning})

	h := NewHandler(m, WithPublisher(&mockPublisher{}))

	query := &Query{OnlyHealthy: true}
	msg := makeNATSMsg(t, query, testReplySubject)
	assertions.NoError(h.HandleQueryServices(context.Background(), msg))

	ctx := core.WithTenantID(context.Background(), testTenantID1)
	assertions.NoError(h.HandleQueryServices(ctx, msg))
}

func TestHandlerHandleDiscoverServices(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{
		ID: testServiceID1, Name: testServiceTest, Status: StatusRunning,
		Type: ServiceTypeAPI, Version: testVersion, Capabilities: []string{"cap1"}, Tags: []string{"tag1"},
	})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})

	h := NewHandler(m, WithPublisher(&mockPublisher{}))

	req := &DiscoverRequest{
		Name: testServiceTest, Type: ServiceTypeAPI, Version: testVersion,
		Capabilities: []string{"cap1"}, Tags: []string{"tag1"}, OnlyHealthy: true, TenantID: testTenantID1,
	}
	msg := makeNATSMsg(t, req, testReplySubject)
	assertions.NoError(h.HandleDiscoverServices(ctx, msg))

	tenantCtx := core.WithTenantID(ctx, testTenantID2)
	assertions.NoError(h.HandleDiscoverServices(tenantCtx, msg))
}

func TestHandlerHandleRegisterInstance(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	m.RegisterServiceDirect(context.Background(), &Info{ID: testServiceID1, Name: testServiceTest})

	h := NewHandler(m, WithPublisher(&mockPublisher{}))

	req := &InstanceRegistrationRequest{
		Instance: &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost},
		TTL:      5 * time.Minute,
	}
	msg := makeNATSMsg(t, req, testReplySubject)
	assertions.NoError(h.HandleRegisterInstance(context.Background(), msg))

	badMsg := &nats.Msg{Data: []byte(testInvalidJSON), Header: make(nats.Header), Reply: testReplySubject}
	assertions.Error(h.HandleRegisterInstance(context.Background(), badMsg))
}

func TestHandlerHandleDeregisterInstance(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})

	h := NewHandler(m, WithPublisher(&mockPublisher{}))

	req := &InstanceDeregistrationRequest{InstanceID: testInstanceID1}
	msg := makeNATSMsg(t, req, testReplySubject)
	assertions.NoError(h.HandleDeregisterInstance(ctx, msg))

	badMsg := &nats.Msg{Data: []byte(testInvalidJSON), Header: make(nats.Header), Reply: testReplySubject}
	assertions.Error(h.HandleDeregisterInstance(ctx, badMsg))
}

func TestHandlerHandleRenewInstance(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})

	h := NewHandler(m, WithPublisher(&mockPublisher{}))

	req := &RenewInstanceRequest{InstanceID: testInstanceID1, TTL: 5 * time.Minute}
	msg := makeNATSMsg(t, req, testReplySubject)
	assertions.NoError(h.HandleRenewInstance(ctx, msg))

	badMsg := &nats.Msg{Data: []byte(testInvalidJSON), Header: make(nats.Header), Reply: testReplySubject}
	assertions.Error(h.HandleRenewInstance(ctx, badMsg))
}

func TestHandlerHandleQueryInstances(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost, Status: StatusRunning})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID2, ServiceID: testServiceID1, Host: testHost, Status: StatusDegraded})

	h := NewHandler(m, WithPublisher(&mockPublisher{}))

	testCases := []struct {
		name string
		req  *QueryInstancesRequest
	}{
		{"only healthy", &QueryInstancesRequest{ServiceID: testServiceID1, OnlyHealthy: true}},
		{"only available", &QueryInstancesRequest{ServiceID: testServiceID1, OnlyAvailable: true}},
		{"all", &QueryInstancesRequest{ServiceID: testServiceID1}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			msg := makeNATSMsg(t, tc.req, testReplySubject)
			assertions.NoError(h.HandleQueryInstances(ctx, msg))
		})
	}

	badMsg := &nats.Msg{Data: []byte(testInvalidJSON), Header: make(nats.Header), Reply: testReplySubject}
	assertions.Error(h.HandleQueryInstances(ctx, badMsg))
}

func TestHandlerHandleHealthUpdate(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest, Status: StatusStarting})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})

	h := NewHandler(m, WithPublisher(&mockPublisher{}))

	// instance health update
	msg := makeNATSMsg(t, &HealthUpdate{ServiceID: testServiceID1, InstanceID: testInstanceID1, Status: StatusRunning}, testReplySubject)
	assertions.NoError(h.HandleHealthUpdate(ctx, msg))

	// service health update
	msg2 := makeNATSMsg(t, &HealthUpdate{ServiceID: testServiceID1, Status: StatusRunning}, testReplySubject)
	assertions.NoError(h.HandleHealthUpdate(ctx, msg2))

	// missing service ID
	msg3 := makeNATSMsg(t, &HealthUpdate{Status: StatusRunning}, testReplySubject)
	assertions.Error(h.HandleHealthUpdate(ctx, msg3))

	// invalid JSON
	badMsg := &nats.Msg{Data: []byte(testInvalidJSON), Header: make(nats.Header), Reply: testReplySubject}
	assertions.Error(h.HandleHealthUpdate(ctx, badMsg))
}

func TestHandlerWithoutPublisher(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	h := NewHandler(m)

	req := &RegistrationRequest{Service: &Info{ID: testServiceID1, Name: testServiceTest}}
	msg := makeNATSMsg(t, req, testReplySubject)
	assertions.NoError(h.HandleRegisterService(context.Background(), msg))
}

func TestHandlerSendResponseErrors(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	pub := &mockPublisher{publishErr: errors.New("publish failed")}
	h := NewHandler(m, WithPublisher(pub))

	req := &RegistrationRequest{Service: &Info{ID: testServiceID1, Name: testServiceTest}}
	msg := makeNATSMsg(t, req, testReplySubject)
	assertions.Error(h.HandleRegisterService(context.Background(), msg))
}

/* ========================================================================================================
   CLIENT METHOD TESTS
   ======================================================================================================== */

func TestClientRegisterService(t *testing.T) {
	assertions := assert.New(t)
	respData, _ := json.Marshal(&ServiceRegistrationResponse{Success: true, ServiceID: testServiceID1})
	pub := &mockPublisher{requestReply: &nats.Msg{Data: respData, Header: make(nats.Header)}}
	c := NewClient(pub)

	resp, err := c.RegisterService(context.Background(), &Info{ID: testServiceID1, Name: testServiceTest})
	assertions.NoError(err)
	assertions.True(resp.Success)
	assertions.Equal(testServiceID1, resp.ServiceID)

	pub.requestErr = errors.New("request failed")
	_, err = c.RegisterService(context.Background(), &Info{ID: testServiceID1, Name: testServiceTest})
	assertions.Error(err)
}

func TestClientDeregisterService(t *testing.T) {
	assertions := assert.New(t)
	respData, _ := json.Marshal(&ServiceDeregistrationResponse{Success: true, ServiceID: testServiceID1})
	pub := &mockPublisher{requestReply: &nats.Msg{Data: respData, Header: make(nats.Header)}}
	c := NewClient(pub)

	resp, err := c.DeregisterService(context.Background(), testServiceID1, testShutdown)
	assertions.NoError(err)
	assertions.True(resp.Success)
}

func TestClientQueryServices(t *testing.T) {
	assertions := assert.New(t)
	respData, _ := json.Marshal(&QueryServicesResponse{Services: []*Info{{ID: testServiceID1}}, Count: 1})
	pub := &mockPublisher{requestReply: &nats.Msg{Data: respData, Header: make(nats.Header)}}
	c := NewClient(pub)

	resp, err := c.QueryServices(context.Background(), &Query{OnlyHealthy: true})
	assertions.NoError(err)
	assertions.Equal(1, resp.Count)
}

func TestClientDiscover(t *testing.T) {
	assertions := assert.New(t)
	respData, _ := json.Marshal(&DiscoverResponse{Services: []*Info{{ID: testServiceID1}}, Count: 1})
	pub := &mockPublisher{requestReply: &nats.Msg{Data: respData, Header: make(nats.Header)}}
	c := NewClient(pub)

	resp, err := c.Discover(context.Background(), testServiceTest)
	assertions.NoError(err)
	assertions.Equal(1, resp.Count)

	resp, err = c.Discover(context.Background(), testServiceTest, WithType(ServiceTypeAPI), WithVersion(testVersion))
	assertions.NoError(err)
	_ = resp
}

func TestClientRegisterInstance(t *testing.T) {
	assertions := assert.New(t)
	respData, _ := json.Marshal(&InstanceRegistrationResponse{Success: true, InstanceID: testInstanceID1, ServiceID: testServiceID1})
	pub := &mockPublisher{requestReply: &nats.Msg{Data: respData, Header: make(nats.Header)}}
	c := NewClient(pub)

	resp, err := c.RegisterInstance(context.Background(), &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost}, 5*time.Minute)
	assertions.NoError(err)
	assertions.True(resp.Success)
}

func TestClientDeregisterInstance(t *testing.T) {
	assertions := assert.New(t)
	respData, _ := json.Marshal(&InstanceDeregistrationResponse{Success: true, InstanceID: testInstanceID1})
	pub := &mockPublisher{requestReply: &nats.Msg{Data: respData, Header: make(nats.Header)}}
	c := NewClient(pub)

	resp, err := c.DeregisterInstance(context.Background(), testInstanceID1, testShutdown)
	assertions.NoError(err)
	assertions.True(resp.Success)
}

func TestClientRenewInstance(t *testing.T) {
	assertions := assert.New(t)
	respData, _ := json.Marshal(&RenewInstanceResponse{Success: true, InstanceID: testInstanceID1})
	pub := &mockPublisher{requestReply: &nats.Msg{Data: respData, Header: make(nats.Header)}}
	c := NewClient(pub)

	resp, err := c.RenewInstance(context.Background(), testInstanceID1, 5*time.Minute)
	assertions.NoError(err)
	assertions.True(resp.Success)
}

func TestClientQueryInstances(t *testing.T) {
	assertions := assert.New(t)
	respData, _ := json.Marshal(&QueryInstancesResponse{Instances: []*Instance{{ID: testInstanceID1}}, Count: 1, ServiceID: testServiceID1})
	pub := &mockPublisher{requestReply: &nats.Msg{Data: respData, Header: make(nats.Header)}}
	c := NewClient(pub)

	resp, err := c.QueryInstances(context.Background(), testServiceID1, true)
	assertions.NoError(err)
	assertions.Equal(1, resp.Count)
}

func TestClientUpdateHealth(t *testing.T) {
	assertions := assert.New(t)
	pub := &mockPublisher{}
	c := NewClient(pub)

	update := &HealthUpdate{ServiceID: testServiceID1, InstanceID: testInstanceID1, Status: StatusRunning}
	assertions.NoError(c.UpdateHealth(context.Background(), update))
	assertions.Len(pub.published, 1)
}

/* ========================================================================================================
   MANAGER EVENT PUBLISHING TESTS
   ======================================================================================================== */

func TestManagerPublishEvents(t *testing.T) {
	assertions := assert.New(t)
	pub := &mockPublisher{}
	m := newTestManagerWithPub(t, pub)
	ctx := context.Background()

	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest, Status: StatusRunning})
	time.Sleep(50 * time.Millisecond)

	pub.mu.Lock()
	initialPubs := len(pub.published)
	pub.mu.Unlock()
	assertions.Greater(initialPubs, 0)

	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest, Status: StatusDegraded})
	time.Sleep(50 * time.Millisecond)

	m.RegisterInstanceDirect(ctx, &Instance{
		ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost, Status: StatusRunning, HealthCheckPath: testHealthPath,
	})
	time.Sleep(50 * time.Millisecond)

	m.DeregisterInstance(ctx, testInstanceID1)
	time.Sleep(50 * time.Millisecond)

	m.DeregisterServiceDirect(ctx, testServiceID1)
	time.Sleep(50 * time.Millisecond)

	pub.mu.Lock()
	finalPubs := len(pub.published)
	pub.mu.Unlock()
	assertions.Greater(finalPubs, initialPubs)
}

func TestManagerHealthCheckEnabled(t *testing.T) {
	m, err := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{
			Enabled: true, Interval: 50 * time.Millisecond,
			Timeout: time.Second, InitialDelay: 10 * time.Millisecond,
		},
	})
	assert.New(t).NoError(err)

	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})
	m.RegisterInstanceDirect(ctx, &Instance{
		ID: testInstanceID1, ServiceID: testServiceID1, Host: "localhost:99999",
		HealthCheckPath: testHealthPath, Status: StatusRunning,
	})

	time.Sleep(100 * time.Millisecond)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	assert.New(t).NoError(m.Shutdown(shutdownCtx))
}

/* ========================================================================================================
   CLEANUP LOOP TESTS
   ======================================================================================================== */

func TestRegistryCleanupLoop(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry(RegistryConfig{CleanupInterval: 50 * time.Millisecond, ExpirationEnabled: true})

	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})
	r.RegisterInstance(&Instance{
		ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost,
		ExpiresAt: time.Now().Add(20 * time.Millisecond),
	})

	time.Sleep(100 * time.Millisecond)

	_, exists := r.GetInstance(testInstanceID1)
	assertions.False(exists, "expired instance should have been cleaned up")
	r.Shutdown(context.Background())
}

/* ========================================================================================================
   INDEX MANAGEMENT TESTS
   ======================================================================================================== */

func TestRegistryIndices(t *testing.T) {
	assertions := assert.New(t)
	r := newTestRegistry(t)

	svc := &Info{
		ID: testServiceID1, Name: "TestService", Type: ServiceTypeAPI,
		Capabilities: []string{"cap1", "cap2"}, Tags: []string{"tag1", "tag2"},
	}
	r.RegisterService(svc)

	svc.Capabilities = []string{"cap3"}
	svc.Tags = []string{"tag3"}
	r.RegisterService(svc)

	r.DeregisterService(testServiceID1)

	assertions.Empty(r.byName)
	assertions.Empty(r.byType)
}

/* ========================================================================================================
   SUBSCRIPTION TESTS
   ======================================================================================================== */

func TestHandlerRegisterSubscriptions(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	h := NewHandler(m)
	sub := &mockSubscriber{}

	subs, err := h.RegisterSubscriptions(context.Background(), sub)
	assertions.NoError(err)
	assertions.Len(subs, 9)

	sub2 := &mockSubscriber{subscribeErr: errors.New("subscribe failed")}
	_, err = h.RegisterSubscriptions(context.Background(), sub2)
	assertions.Error(err)
}

func TestHandlerRegisterQueueSubscriptions(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	h := NewHandler(m)
	sub := &mockSubscriber{}

	subs, err := h.RegisterQueueSubscriptions(context.Background(), sub, "worker-queue")
	assertions.NoError(err)
	assertions.Len(subs, 8)
}

func TestHandlerUnsubscribeAll(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	h := NewHandler(m)

	subs := []core.Subscription{
		&mockSubscription{valid: true},
		nil,
		&mockSubscription{valid: true},
	}

	h.unsubscribeAll(subs)
	assertions.False(subs[0].(*mockSubscription).valid)
	assertions.False(subs[2].(*mockSubscription).valid)
}

/* ========================================================================================================
   BENCHMARKS
   ======================================================================================================== */

/* ========================================================================================================
   ADDITIONAL COVERAGE TESTS
   ======================================================================================================== */

// failingAfterNSubscriber fails after N successful subscriptions.
type failingAfterNSubscriber struct {
	callCount    int
	failAfter    int
	subscribeErr error
	subs         []*mockSubscription
}

func (f *failingAfterNSubscriber) Subscribe(_ context.Context, subject string, _ core.MessageHandler) (core.Subscription, error) {
	f.callCount++
	if f.callCount > f.failAfter {
		return nil, f.subscribeErr
	}
	sub := &mockSubscription{subject: subject, valid: true}
	f.subs = append(f.subs, sub)
	return sub, nil
}

func (f *failingAfterNSubscriber) QueueSubscribe(_ context.Context, subject string, _ string, _ core.MessageHandler) (core.Subscription, error) {
	f.callCount++
	if f.callCount > f.failAfter {
		return nil, f.subscribeErr
	}
	sub := &mockSubscription{subject: subject, valid: true}
	f.subs = append(f.subs, sub)
	return sub, nil
}

func (f *failingAfterNSubscriber) Close(_ context.Context) error { return nil }

var (
	testSubscribeFailed = "subscribe failed"
	testQueueName       = "worker-queue"
)

func TestRegisterSubscriptionsFailAtEachStep(t *testing.T) {
	testCases := []struct {
		name      string
		failAfter int
	}{
		{"fail at 2nd subscribe", 1},
		{"fail at 3rd subscribe", 2},
		{"fail at 4th subscribe", 3},
		{"fail at 5th subscribe", 4},
		{"fail at 6th subscribe", 5},
		{"fail at 7th subscribe", 6},
		{"fail at 8th subscribe", 7},
		{"fail at 9th subscribe", 8},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			m := newTestManager(t)
			h := NewHandler(m)
			sub := &failingAfterNSubscriber{
				failAfter:    tc.failAfter,
				subscribeErr: errors.New(testSubscribeFailed),
			}

			_, err := h.RegisterSubscriptions(context.Background(), sub)
			assertions.Error(err)
			// Verify earlier subscriptions were cleaned up
			for _, s := range sub.subs {
				assertions.False(s.valid)
			}
		})
	}
}

func TestRegisterQueueSubscriptionsFailAtEachStep(t *testing.T) {
	testCases := []struct {
		name      string
		failAfter int
	}{
		{"fail at 2nd queue subscribe", 1},
		{"fail at 3rd queue subscribe", 2},
		{"fail at 4th queue subscribe", 3},
		{"fail at 5th queue subscribe", 4},
		{"fail at 6th queue subscribe", 5},
		{"fail at 7th queue subscribe", 6},
		{"fail at 8th queue subscribe", 7},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			m := newTestManager(t)
			h := NewHandler(m)
			sub := &failingAfterNSubscriber{
				failAfter:    tc.failAfter,
				subscribeErr: errors.New(testSubscribeFailed),
			}

			_, err := h.RegisterQueueSubscriptions(context.Background(), sub, testQueueName)
			assertions.Error(err)
			for _, s := range sub.subs {
				assertions.False(s.valid)
			}
		})
	}
}

func TestSendErrorResponseWithPublisher(t *testing.T) {
	var testCaseMsg = "original error"

	testCases := []struct {
		name       string
		reply      string
		publisher  *mockPublisher
		wantErrMsg string
	}{
		{
			name:       "no reply subject returns original error",
			reply:      "",
			publisher:  &mockPublisher{},
			wantErrMsg: testCaseMsg,
		},
		{
			name:       "no publisher returns original error",
			reply:      testReplySubject,
			publisher:  nil,
			wantErrMsg: testCaseMsg,
		},
		{
			name:       "publish fails returns publish error",
			reply:      testReplySubject,
			publisher:  &mockPublisher{publishErr: errors.New("publish has been failed")},
			wantErrMsg: "publish has been failed",
		},
		{
			name:       "publish succeeds returns original error",
			reply:      testReplySubject,
			publisher:  &mockPublisher{},
			wantErrMsg: testCaseMsg,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			m := newTestManager(t)
			h := NewHandler(m)
			if tc.publisher != nil {
				h.SetPublisher(tc.publisher)
			}

			msg := &nats.Msg{
				Data:   []byte("{}"),
				Header: make(nats.Header),
				Reply:  tc.reply,
			}
			err := h.sendErrorResponse(context.Background(), msg, errors.New(testCaseMsg))
			assertions.Error(err)
			assertions.True(contains(err.Error(), tc.wantErrMsg))
		})
	}
}

func TestManagerRemoveFromHealthCheck(t *testing.T) {
	assertions := assert.New(t)

	m, err := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{
			Enabled: true, Interval: 5 * time.Second,
			Timeout: time.Second, InitialDelay: time.Hour,
		},
	})
	assertions.NoError(err)

	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})
	m.RegisterInstanceDirect(ctx, &Instance{
		ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost,
		HealthCheckPath: testHealthPath, Status: StatusRunning,
	})

	// Wait for async callback to fire
	time.Sleep(30 * time.Millisecond)

	// Verify instance was added to health check
	m.healthMu.RLock()
	_, inHealthCheck := m.healthInstances[testInstanceID1]
	m.healthMu.RUnlock()
	assertions.True(inHealthCheck)

	// Deregister instance triggers removeFromHealthCheck
	assertions.NoError(m.DeregisterInstance(ctx, testInstanceID1))
	time.Sleep(30 * time.Millisecond)

	m.healthMu.RLock()
	_, inHealthCheck = m.healthInstances[testInstanceID1]
	m.healthMu.RUnlock()
	assertions.False(inHealthCheck)

	shutdownCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	assertions.NoError(m.Shutdown(shutdownCtx))
}

func TestManagerDirectMethodErrorPaths(t *testing.T) {
	t.Run("RegisterServiceDirect duplicate after shutdown", func(t *testing.T) {
		assertions := assert.New(t)
		m := newTestManager(t)
		// Cancel the registry to trigger errors
		m.registry.cancel()
		err := m.RegisterServiceDirect(context.Background(), &Info{ID: testServiceID1, Name: testServiceTest})
		assertions.Error(err)
	})

	t.Run("DeregisterServiceDirect not found", func(t *testing.T) {
		assertions := assert.New(t)
		m := newTestManager(t)
		err := m.DeregisterServiceDirect(context.Background(), "nonexistent")
		assertions.Error(err)
	})

	t.Run("RegisterInstanceDirect service not found", func(t *testing.T) {
		assertions := assert.New(t)
		m := newTestManager(t)
		err := m.RegisterInstanceDirect(context.Background(), &Instance{ID: testInstanceID1, ServiceID: "nonexistent", Host: testHost})
		assertions.Error(err)
	})

	t.Run("RegisterInstanceWithTTL service not found", func(t *testing.T) {
		assertions := assert.New(t)
		m := newTestManager(t)
		err := m.RegisterInstanceWithTTL(context.Background(), &Instance{ID: testInstanceID1, ServiceID: "nonexistent", Host: testHost}, 5*time.Minute)
		assertions.Error(err)
	})

	t.Run("DeregisterInstance not found", func(t *testing.T) {
		assertions := assert.New(t)
		m := newTestManager(t)
		err := m.DeregisterInstance(context.Background(), "nonexistent")
		assertions.Error(err)
	})
}

func TestManagerShutdownTimeout(t *testing.T) {
	assertions := assert.New(t)
	pub := &mockPublisher{}
	m, err := NewManager(ManagerConfig{Publisher: pub, HealthCheck: HealthCheckConfig{Enabled: false}})
	assertions.NoError(err)

	// Add a goroutine to the WaitGroup that blocks until we release it
	blocker := make(chan struct{})
	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		<-blocker
	}()

	// Shutdown with already-cancelled context should timeout
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err = m.Shutdown(ctx)
	assertions.ErrorIs(err, context.Canceled)

	close(blocker)
}

func TestHandlerRegisterServiceWithTenantContext(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	pub := &mockPublisher{}
	h := NewHandler(m, WithPublisher(pub))

	req := &RegistrationRequest{Service: &Info{ID: testServiceID1, Name: testServiceTest}}
	msg := makeNATSMsg(t, req, testReplySubject)

	ctx := core.WithTenantID(context.Background(), testTenantID1)
	assertions.NoError(h.HandleRegisterService(ctx, msg))

	svc, exists := m.GetService(testServiceID1)
	assertions.True(exists)
	assertions.Equal(testTenantID1, svc.TenantID)
}

func TestHandlerRegisterServiceManagerError(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	pub := &mockPublisher{}
	h := NewHandler(m, WithPublisher(pub))

	// Register invalid service (missing name) to trigger manager error
	req := &RegistrationRequest{Service: &Info{ID: testServiceID1}}
	msg := makeNATSMsg(t, req, testReplySubject)
	err := h.HandleRegisterService(context.Background(), msg)
	assertions.Error(err)
}

func TestHandlerDeregisterServiceManagerError(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	pub := &mockPublisher{}
	h := NewHandler(m, WithPublisher(pub))

	// Deregister non-existent service
	req := &DeregistrationRequest{ServiceID: "nonexistent"}
	msg := makeNATSMsg(t, req, testReplySubject)
	err := h.HandleDeregisterService(context.Background(), msg)
	assertions.Error(err)
}

func TestHandlerRegisterInstanceManagerError(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	pub := &mockPublisher{}
	h := NewHandler(m, WithPublisher(pub))

	// Register instance for non-existent service
	req := &InstanceRegistrationRequest{
		Instance: &Instance{ID: testInstanceID1, ServiceID: "nonexistent", Host: testHost},
	}
	msg := makeNATSMsg(t, req, testReplySubject)
	err := h.HandleRegisterInstance(context.Background(), msg)
	assertions.Error(err)
}

func TestHandlerDeregisterInstanceManagerError(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	pub := &mockPublisher{}
	h := NewHandler(m, WithPublisher(pub))

	// Deregister non-existent instance
	req := &InstanceDeregistrationRequest{InstanceID: "nonexistent"}
	msg := makeNATSMsg(t, req, testReplySubject)
	err := h.HandleDeregisterInstance(context.Background(), msg)
	assertions.Error(err)
}

func TestHandlerRenewInstanceManagerError(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	pub := &mockPublisher{}
	h := NewHandler(m, WithPublisher(pub))

	// Renew non-existent instance
	req := &RenewInstanceRequest{InstanceID: "nonexistent", TTL: 5 * time.Minute}
	msg := makeNATSMsg(t, req, testReplySubject)
	err := h.HandleRenewInstance(context.Background(), msg)
	assertions.Error(err)
}

func TestHandlerQueryServicesInvalidJSON(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	pub := &mockPublisher{}
	h := NewHandler(m, WithPublisher(pub))

	badMsg := &nats.Msg{Data: []byte(testInvalidJSON), Header: make(nats.Header), Reply: testReplySubject}
	assertions.Error(h.HandleQueryServices(context.Background(), badMsg))
}

func TestHandlerDiscoverServicesInvalidJSON(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	pub := &mockPublisher{}
	h := NewHandler(m, WithPublisher(pub))

	badMsg := &nats.Msg{Data: []byte(testInvalidJSON), Header: make(nats.Header), Reply: testReplySubject}
	assertions.Error(h.HandleDiscoverServices(context.Background(), badMsg))
}

func TestHandlerHealthUpdateWithZeroCheckedAt(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})

	pub := &mockPublisher{}
	h := NewHandler(m, WithPublisher(pub))

	// Health update with zero CheckedAt should auto-set to now
	update := &HealthUpdate{InstanceID: testInstanceID1, ServiceID: testServiceID1, Status: StatusRunning}
	msg := makeNATSMsg(t, update, testReplySubject)
	assertions.NoError(h.HandleHealthUpdate(ctx, msg))
}

func TestHTTPHandleServiceDeregistrationBadJSON(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)

	rr := httptest.NewRecorder()
	m.HandleServiceDeregistration(rr, httptest.NewRequest(http.MethodDelete, deregisterEndpoint, bytes.NewReader([]byte(testInvalidJSON))))
	assertions.Equal(http.StatusBadRequest, rr.Code)
}

func TestHTTPHandleServiceDeregistrationPOST(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	m.RegisterServiceDirect(context.Background(), &Info{ID: testServiceID1, Name: testServiceTest})

	req := DeregistrationRequest{ServiceID: testServiceID1}
	body, _ := json.Marshal(req)

	rr := httptest.NewRecorder()
	m.HandleServiceDeregistration(rr, httptest.NewRequest(http.MethodPost, deregisterEndpoint, bytes.NewReader(body)))
	assertions.Equal(http.StatusOK, rr.Code)
}

func TestHTTPHandleServiceDeregistrationInternalError(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	// Cancel registry to force internal error
	m.registry.cancel()

	req := DeregistrationRequest{ServiceID: testServiceID1}
	body, _ := json.Marshal(req)

	rr := httptest.NewRecorder()
	m.HandleServiceDeregistration(rr, httptest.NewRequest(http.MethodDelete, deregisterEndpoint, bytes.NewReader(body)))
	// not-found or internal error
	assertions.True(rr.Code == http.StatusNotFound || rr.Code == http.StatusInternalServerError)
}

func TestHTTPHandleInstanceRegistrationBadJSON(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)

	rr := httptest.NewRecorder()
	m.HandleInstanceRegistration(rr, httptest.NewRequest(http.MethodPost, registerInstanceEndpoint, bytes.NewReader([]byte(testInvalidJSON))))
	assertions.Equal(http.StatusBadRequest, rr.Code)
}

func TestHTTPHandleInstanceRegistrationError(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)

	// Register for non-existent service to trigger internal error
	req := InstanceRegistrationRequest{Instance: &Instance{ID: testInstanceID1, ServiceID: "nonexistent", Host: testHost}}
	body, _ := json.Marshal(req)

	rr := httptest.NewRecorder()
	m.HandleInstanceRegistration(rr, httptest.NewRequest(http.MethodPost, registerInstanceEndpoint, bytes.NewReader(body)))
	assertions.Equal(http.StatusInternalServerError, rr.Code)
}

func TestHTTPHandleInstanceRenewBadJSON(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)

	rr := httptest.NewRecorder()
	m.HandleInstanceRenew(rr, httptest.NewRequest(http.MethodPost, renewEndpoint, bytes.NewReader([]byte(testInvalidJSON))))
	assertions.Equal(http.StatusBadRequest, rr.Code)
}

func TestHTTPHandleInstanceRenewPUT(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})
	m.RegisterInstanceDirect(ctx, &Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})

	reqBody := struct {
		InstanceID string        `json:"instance_id"`
		TTL        time.Duration `json:"ttl"`
	}{InstanceID: testInstanceID1, TTL: 5 * time.Minute}
	body, _ := json.Marshal(reqBody)

	rr := httptest.NewRecorder()
	m.HandleInstanceRenew(rr, httptest.NewRequest(http.MethodPut, renewEndpoint, bytes.NewReader(body)))
	assertions.Equal(http.StatusOK, rr.Code)
}

func TestHTTPHandleServiceQueryBadJSON(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)

	rr := httptest.NewRecorder()
	m.HandleServiceQuery(rr, httptest.NewRequest(http.MethodPost, servicesEndpoint, bytes.NewReader([]byte(testInvalidJSON))))
	assertions.Equal(http.StatusBadRequest, rr.Code)
}

func TestHTTPHandleServiceRegistrationError(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)

	// Register with invalid service to trigger internal error
	req := RegistrationRequest{Service: &Info{ID: testServiceID1}}
	body, _ := json.Marshal(req)

	rr := httptest.NewRecorder()
	m.HandleServiceRegistration(rr, httptest.NewRequest(http.MethodPost, registerEndpoint, bytes.NewReader(body)))
	assertions.Equal(http.StatusInternalServerError, rr.Code)
}

func TestClientRequestErrors(t *testing.T) {
	var testCaseMsg = "request has been failed"
	t.Run("DeregisterService request error", func(t *testing.T) {
		assertions := assert.New(t)
		pub := &mockPublisher{requestErr: errors.New(testCaseMsg)}
		c := NewClient(pub)
		_, err := c.DeregisterService(context.Background(), testServiceID1, testShutdown)
		assertions.Error(err)
	})

	t.Run("QueryServices request error", func(t *testing.T) {
		assertions := assert.New(t)
		pub := &mockPublisher{requestErr: errors.New(testCaseMsg)}
		c := NewClient(pub)
		_, err := c.QueryServices(context.Background(), &Query{})
		assertions.Error(err)
	})

	t.Run("Discover request error", func(t *testing.T) {
		assertions := assert.New(t)
		pub := &mockPublisher{requestErr: errors.New(testCaseMsg)}
		c := NewClient(pub)
		_, err := c.Discover(context.Background(), testServiceTest)
		assertions.Error(err)
	})

	t.Run("RegisterInstance request error", func(t *testing.T) {
		assertions := assert.New(t)
		pub := &mockPublisher{requestErr: errors.New(testCaseMsg)}
		c := NewClient(pub)
		_, err := c.RegisterInstance(context.Background(), &Instance{ID: testInstanceID1}, 5*time.Minute)
		assertions.Error(err)
	})

	t.Run("DeregisterInstance request error", func(t *testing.T) {
		assertions := assert.New(t)
		pub := &mockPublisher{requestErr: errors.New(testCaseMsg)}
		c := NewClient(pub)
		_, err := c.DeregisterInstance(context.Background(), testInstanceID1, testShutdown)
		assertions.Error(err)
	})

	t.Run("RenewInstance request error", func(t *testing.T) {
		assertions := assert.New(t)
		pub := &mockPublisher{requestErr: errors.New(testCaseMsg)}
		c := NewClient(pub)
		_, err := c.RenewInstance(context.Background(), testInstanceID1, 5*time.Minute)
		assertions.Error(err)
	})

	t.Run("QueryInstances request error", func(t *testing.T) {
		assertions := assert.New(t)
		pub := &mockPublisher{requestErr: errors.New(testCaseMsg)}
		c := NewClient(pub)
		_, err := c.QueryInstances(context.Background(), testServiceID1, true)
		assertions.Error(err)
	})

	t.Run("UpdateHealth publish error", func(t *testing.T) {
		assertions := assert.New(t)
		pub := &mockPublisher{publishErr: errors.New("publish is failed")}
		c := NewClient(pub)
		err := c.UpdateHealth(context.Background(), &HealthUpdate{ServiceID: testServiceID1, Status: StatusRunning})
		assertions.Error(err)
	})
}

func TestClientUnmarshalErrors(t *testing.T) {
	badReply := &nats.Msg{Data: []byte(testInvalidJSON), Header: make(nats.Header)}

	t.Run("RegisterService bad response", func(t *testing.T) {
		assertions := assert.New(t)
		pub := &mockPublisher{requestReply: badReply}
		c := NewClient(pub)
		_, err := c.RegisterService(context.Background(), &Info{ID: testServiceID1, Name: testServiceTest})
		assertions.Error(err)
	})

	t.Run("DeregisterService bad response", func(t *testing.T) {
		assertions := assert.New(t)
		pub := &mockPublisher{requestReply: badReply}
		c := NewClient(pub)
		_, err := c.DeregisterService(context.Background(), testServiceID1, testShutdown)
		assertions.Error(err)
	})

	t.Run("QueryServices bad response", func(t *testing.T) {
		assertions := assert.New(t)
		pub := &mockPublisher{requestReply: badReply}
		c := NewClient(pub)
		_, err := c.QueryServices(context.Background(), &Query{})
		assertions.Error(err)
	})

	t.Run("Discover bad response", func(t *testing.T) {
		assertions := assert.New(t)
		pub := &mockPublisher{requestReply: badReply}
		c := NewClient(pub)
		_, err := c.Discover(context.Background(), testServiceTest)
		assertions.Error(err)
	})

	t.Run("RegisterInstance bad response", func(t *testing.T) {
		assertions := assert.New(t)
		pub := &mockPublisher{requestReply: badReply}
		c := NewClient(pub)
		_, err := c.RegisterInstance(context.Background(), &Instance{ID: testInstanceID1}, 5*time.Minute)
		assertions.Error(err)
	})

	t.Run("DeregisterInstance bad response", func(t *testing.T) {
		assertions := assert.New(t)
		pub := &mockPublisher{requestReply: badReply}
		c := NewClient(pub)
		_, err := c.DeregisterInstance(context.Background(), testInstanceID1, testShutdown)
		assertions.Error(err)
	})

	t.Run("RenewInstance bad response", func(t *testing.T) {
		assertions := assert.New(t)
		pub := &mockPublisher{requestReply: badReply}
		c := NewClient(pub)
		_, err := c.RenewInstance(context.Background(), testInstanceID1, 5*time.Minute)
		assertions.Error(err)
	})

	t.Run("QueryInstances bad response", func(t *testing.T) {
		assertions := assert.New(t)
		pub := &mockPublisher{requestReply: badReply}
		c := NewClient(pub)
		_, err := c.QueryInstances(context.Background(), testServiceID1, true)
		assertions.Error(err)
	})
}

func TestManagerPublishServiceEventErrors(t *testing.T) {
	assertions := assert.New(t)
	pub := &mockPublisher{publishErr: errors.New("publish is being failed")}
	m := newTestManagerWithPub(t, pub)

	var errorReceived bool
	m.OnError(func(err error) { errorReceived = true })

	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest, Status: StatusRunning})
	time.Sleep(50 * time.Millisecond)

	assertions.True(errorReceived)
}

func TestManagerPublishInstanceEventErrors(t *testing.T) {
	assertions := assert.New(t)
	pub := &mockPublisher{publishErr: errors.New("publish has got failed")}
	m := newTestManagerWithPub(t, pub)

	var errorReceived bool
	m.OnError(func(err error) { errorReceived = true })

	ctx := context.Background()
	m.RegisterServiceDirect(ctx, &Info{ID: testServiceID1, Name: testServiceTest})
	time.Sleep(20 * time.Millisecond)

	m.RegisterInstanceDirect(ctx, &Instance{
		ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost, Status: StatusRunning,
	})
	time.Sleep(50 * time.Millisecond)

	assertions.True(errorReceived)
}

func TestRegistryAllInstancesWithExpired(t *testing.T) {
	assertions := assert.New(t)
	r := newTestRegistry(t)

	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})
	r.RegisterInstance(&Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})
	r.RegisterInstance(&Instance{
		ID: testInstanceID2, ServiceID: testServiceID1, Host: testHost,
		ExpiresAt: time.Now().Add(-time.Hour),
	})

	instances := r.AllInstances()
	assertions.Len(instances, 1)
	assertions.Equal(testInstanceID1, instances[0].ID)
}

func TestRegistryGetInstancesWithExpired(t *testing.T) {
	assertions := assert.New(t)
	r := newTestRegistry(t)

	r.RegisterService(&Info{ID: testServiceID1, Name: testServiceTest})
	r.RegisterInstance(&Instance{ID: testInstanceID1, ServiceID: testServiceID1, Host: testHost})
	r.RegisterInstance(&Instance{
		ID: testInstanceID2, ServiceID: testServiceID1, Host: testHost,
		ExpiresAt: time.Now().Add(-time.Hour),
	})

	instances := r.GetInstances(testServiceID1)
	assertions.Len(instances, 1)
}

func TestManagerRegisterServiceRegistryError(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()

	// Register service, then re-register with different data to trigger an update (not error).
	// But first cancel registry context to force an error.
	m.registry.cancel()

	req := &RegistrationRequest{Service: &Info{ID: testServiceID1, Name: testServiceTest}}
	err := m.RegisterService(ctx, req)
	assertions.Error(err)
}

func TestHandlerHealthUpdateWithErrors(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t)
	ctx := context.Background()
	pub := &mockPublisher{}
	h := NewHandler(m, WithPublisher(pub))

	// Health update for non-existent instance
	update := &HealthUpdate{InstanceID: "nonexistent", ServiceID: testServiceID1, Status: StatusRunning}
	msg := makeNATSMsg(t, update, testReplySubject)
	err := h.HandleHealthUpdate(ctx, msg)
	assertions.Error(err)

	// Health update for non-existent service (service-only update)
	update2 := &HealthUpdate{ServiceID: "nonexistent", Status: StatusRunning}
	msg2 := makeNATSMsg(t, update2, testReplySubject)
	err = h.HandleHealthUpdate(ctx, msg2)
	assertions.Error(err)
}
