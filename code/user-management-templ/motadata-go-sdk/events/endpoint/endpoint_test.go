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

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

/* ========================================================================================================
   TEST CONSTANTS AND VARIABLES
   ======================================================================================================== */

var (
	// Endpoint identifiers
	testEndpointID  = "ep1"
	testEndpointID2 = "ep2"
	testEndpointID3 = "ep3"
	testServiceID   = "svc1"
	testServiceID2  = "svc2"
	testInstanceID  = "inst1"
	testInstanceID2 = "inst2"

	// Network
	testHost      = "localhost"
	testHost1     = "host1"
	testHost2     = "host2"
	testHost3     = "host3"
	testPort      = 8080
	testPortAlt   = 9090
	testPortHTTPS = 8443
	testPath      = "/api"
	testHealthURL = "http://localhost:8080/health"

	// Multi-tenancy
	testTenantID  = "tenant1"
	testTenantID2 = "tenant2"
	testRegion    = "us-east-1"
	testRegionAlt = "eu-west-1"
	testZone      = "zone-a"
	testZoneAlt   = "zone-b"

	// Versioning
	testVersion       = "1.0.0"
	testVersionAlt    = "2.0.0"
	testAPIVersion    = "v1"
	testAPIVersionAlt = "v2"

	// Tags and metadata
	testTagAPI        = "api"
	testTagV1         = "v1"
	testTagV2         = "v2"
	testTagGRPC       = "grpc"
	testTagProduction = "production"
	testMetaKey       = "key"
	testMetaValue     = "value"
	testMetaKey2      = "key2"
	testMetaValue2    = "value2"
	testMetaEnvKey    = "env"
	testMetaEnvProd   = "prod"
	testMetaEnvDev    = "dev"

	// Subjects
	testSubjectRegister     = "endpoints.register"
	testSubjectDeregister   = "endpoints.deregister"
	testSubjectRenew        = "endpoints.renew"
	testSubjectQuery        = "endpoints.query"
	testSubjectDiscover     = "endpoints.discover"
	testSubjectEvents       = "endpoints.events"
	testSubjectHealthUpdate = "endpoints.health.update"

	// Handler/reason strings
	testReasonShutdown       = "shutdown"
	testReasonTest           = "test"
	testReasonDuplicate      = "duplicate"
	testReasonDeregistration = "test deregistration"
	testInvalidJSON          = "invalid json"
	testInvalidPayload       = "invalid"

	// HTTP paths
	testHTTPPathRegister   = "/register"
	testHTTPPathDeregister = "/deregister"
	testHTTPPathQuery      = "/query"
	testHTTPPathRenew      = "/renew"

	// Durations
	testTTL5Min  = 5 * time.Minute
	testTTL10Min = 10 * time.Minute
	testTTL1Min  = 1 * time.Minute

	// Callback sleep
	testCallbackWait = 10 * time.Millisecond

	// Queue name
	testQueueName = "endpoint-workers"

	// Reply subjects
	testReplySubject = "_INBOX.reply123"

	// Error messages
	testErrPublish   = "publish failed"
	testErrRequest   = "request failed"
	testErrSubscribe = "subscribe failed"

	// Health check
	testHealthPath = "/health"

	// Content type
	testContentType    = "application/json"
	testContentTypeKey = "Content-Type"
)

/* ========================================================================================================
   HELPER FUNCTIONS
   ======================================================================================================== */

// newTestMsg creates a *nats.Msg suitable for testing.
func newTestMsg(data []byte) *nats.Msg {
	return &nats.Msg{Data: data, Header: make(nats.Header)}
}

// newTestManager creates a Manager with health checks disabled.
func newTestManager() *Manager {
	m, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})
	return m
}

// newTestEndpoint creates a standard test endpoint Info.
func newTestEndpoint(id, serviceID, host string) *Info {
	return &Info{
		ID:        id,
		ServiceID: serviceID,
		Host:      host,
	}
}

/* ========================================================================================================
   MOCK PUBLISHER
   ======================================================================================================== */

type mockPublisher struct {
	publishErr   error
	requestReply *nats.Msg
	requestErr   error
}

func (m *mockPublisher) Publish(_ context.Context, _ string, _ *nats.Msg) error {
	return m.publishErr
}

func (m *mockPublisher) Request(_ context.Context, _ string, _ *nats.Msg) (*nats.Msg, error) {
	if m.requestErr != nil {
		return nil, m.requestErr
	}
	return m.requestReply, nil
}

func (m *mockPublisher) Close(_ context.Context) error {
	return nil
}

/* ========================================================================================================
   MOCK SUBSCRIBER AND SUBSCRIPTION
   ======================================================================================================== */

type mockSubscription struct {
	subject string
	err     error
}

func (s *mockSubscription) Subject() string    { return s.subject }
func (s *mockSubscription) Unsubscribe() error { return s.err }
func (s *mockSubscription) Drain() error       { return s.err }
func (s *mockSubscription) IsValid() bool      { return true }

type mockSubscriber struct {
	subscribeErr      error
	queueSubscribeErr error
	callCount         int
	failAtCall        int // fail at this call number (0 = never fail)
}

func (m *mockSubscriber) Subscribe(_ context.Context, subject string, _ core.MessageHandler) (core.Subscription, error) {
	m.callCount++
	if m.failAtCall > 0 && m.callCount >= m.failAtCall {
		return nil, m.subscribeErr
	}
	if m.subscribeErr != nil && m.failAtCall == 0 {
		return nil, m.subscribeErr
	}
	return &mockSubscription{subject: subject}, nil
}

func (m *mockSubscriber) QueueSubscribe(_ context.Context, subject, _ string, _ core.MessageHandler) (core.Subscription, error) {
	m.callCount++
	if m.failAtCall > 0 && m.callCount >= m.failAtCall {
		return nil, m.queueSubscribeErr
	}
	if m.queueSubscribeErr != nil && m.failAtCall == 0 {
		return nil, m.queueSubscribeErr
	}
	return &mockSubscription{subject: subject}, nil
}

func (m *mockSubscriber) Close(_ context.Context) error {
	return nil
}

// capturingPublisher records published messages for verification.
type capturingPublisher struct {
	mu       sync.Mutex
	messages []*nats.Msg
	subjects []string
	pubErr   error
}

func (c *capturingPublisher) Publish(_ context.Context, subject string, msg *nats.Msg) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.messages = append(c.messages, msg)
	c.subjects = append(c.subjects, subject)
	return c.pubErr
}

func (c *capturingPublisher) Request(_ context.Context, _ string, _ *nats.Msg) (*nats.Msg, error) {
	return nil, nil
}

func (c *capturingPublisher) Close(_ context.Context) error {
	return nil
}

// newTestMsgWithReply creates a *nats.Msg with a Reply subject for request-reply testing.
func newTestMsgWithReply(data []byte, reply string) *nats.Msg {
	return &nats.Msg{Data: data, Header: make(nats.Header), Reply: reply}
}

/* ========================================================================================================
   STATUS TESTS
   ======================================================================================================== */

func TestStatusString(t *testing.T) {
	testCases := []struct {
		name     string
		status   Status
		expected string
	}{
		{"unknown", StatusUnknown, "unknown"},
		{"healthy", StatusHealthy, "healthy"},
		{"unhealthy", StatusUnhealthy, "unhealthy"},
		{"degraded", StatusDegraded, "degraded"},
		{"maintenance", StatusMaintenance, "maintenance"},
		{"offline", StatusOffline, "offline"},
		{"invalid value", Status(99), "unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, tc.status.String())
		})
	}
}

func TestStatusIsAvailable(t *testing.T) {
	testCases := []struct {
		name     string
		status   Status
		expected bool
	}{
		{"unknown is not available", StatusUnknown, false},
		{"healthy is available", StatusHealthy, true},
		{"unhealthy is not available", StatusUnhealthy, false},
		{"degraded is available", StatusDegraded, true},
		{"maintenance is not available", StatusMaintenance, false},
		{"offline is not available", StatusOffline, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, tc.status.IsAvailable())
		})
	}
}

/* ========================================================================================================
   INFO TESTS
   ======================================================================================================== */

func TestInfoClone(t *testing.T) {
	assertions := assert.New(t)

	original := &Info{
		ID:        testEndpointID,
		ServiceID: testServiceID,
		Host:      testHost,
		Port:      testPort,
		Protocol:  ProtocolHTTP,
		Status:    StatusHealthy,
		Methods:   []Method{MethodGET, MethodPOST},
		Tags:      []string{testTagAPI, testTagV1},
		Metadata:  map[string]string{testMetaKey: testMetaValue},
	}

	clone := original.Clone()

	assertions.NotSame(original, clone, "Clone returned same pointer")
	assertions.Equal(original.ID, clone.ID, "ID not cloned correctly")

	// Modify clone and verify original unchanged
	clone.Tags[0] = "modified"
	assertions.NotEqual("modified", original.Tags[0], "Tags not deep copied")

	clone.Metadata[testMetaKey] = "modified"
	assertions.NotEqual("modified", original.Metadata[testMetaKey], "Metadata not deep copied")
}

func TestInfoCloneNil(t *testing.T) {
	assertions := assert.New(t)
	var info *Info
	clone := info.Clone()
	assertions.Nil(clone, "Clone of nil should be nil")
}

func TestInfoURL(t *testing.T) {
	testCases := []struct {
		name     string
		info     Info
		expected string
	}{
		{
			name:     "standard HTTP with non-default port",
			info:     Info{Protocol: ProtocolHTTP, Host: testHost, Port: testPort, Path: testPath},
			expected: "http://localhost:8080/api",
		},
		{
			name:     "HTTP with default port 80 omitted",
			info:     Info{Protocol: ProtocolHTTP, Host: testHost, Port: 80, Path: testPath},
			expected: "http://localhost/api",
		},
		{
			name:     "HTTPS with default port 443 omitted",
			info:     Info{Protocol: ProtocolHTTPS, Host: testHost, Port: 443, Path: testPath},
			expected: "https://localhost/api",
		},
		{
			name:     "HTTPS with non-default port",
			info:     Info{Protocol: ProtocolHTTPS, Host: testHost, Port: testPortHTTPS, Path: testPath},
			expected: "https://localhost:8443/api",
		},
		{
			name:     "HTTP with zero port omitted",
			info:     Info{Protocol: ProtocolHTTP, Host: testHost, Port: 0, Path: testPath},
			expected: "http://localhost/api",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, tc.info.URL())
		})
	}
}

func TestInfoAddress(t *testing.T) {
	testCases := []struct {
		name     string
		info     Info
		expected string
	}{
		{"with port", Info{Host: testHost, Port: testPort}, "localhost:8080"},
		{"without port", Info{Host: testHost, Port: 0}, "localhost"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, tc.info.Address())
		})
	}
}

func TestInfoIsExpired(t *testing.T) {
	testCases := []struct {
		name     string
		info     *Info
		expected bool
	}{
		{"zero ExpiresAt never expires", &Info{}, false},
		{"future time not expired", &Info{ExpiresAt: time.Now().Add(time.Hour)}, false},
		{"past time is expired", &Info{ExpiresAt: time.Now().Add(-time.Hour)}, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, tc.info.IsExpired())
		})
	}
}

func TestInfoIsExpiredEdgeCases(t *testing.T) {
	// Just at expiry time (edge case) - just verify no panic
	info := &Info{ExpiresAt: time.Now()}
	_ = info.IsExpired()
}

func TestInfoIsHealthy(t *testing.T) {
	testCases := []struct {
		name     string
		status   Status
		expected bool
	}{
		{"healthy returns true", StatusHealthy, true},
		{"unhealthy returns false", StatusUnhealthy, false},
		{"degraded returns false", StatusDegraded, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			info := &Info{Status: tc.status}
			assertions.Equal(tc.expected, info.IsHealthy())
		})
	}
}

func TestInfoHasTag(t *testing.T) {
	assertions := assert.New(t)
	info := &Info{Tags: []string{testTagAPI, testTagV1, testTagProduction}}

	assertions.True(info.HasTag(testTagAPI), "Expected to have api tag")
	assertions.False(info.HasTag("staging"), "Expected not to have staging tag")
}

func TestInfoGetMetadata(t *testing.T) {
	testCases := []struct {
		name     string
		info     *Info
		key      string
		expected string
	}{
		{"existing key", &Info{Metadata: map[string]string{testMetaKey: testMetaValue}}, testMetaKey, testMetaValue},
		{"nonexistent key", &Info{Metadata: map[string]string{testMetaKey: testMetaValue}}, "nonexistent", ""},
		{"nil metadata", &Info{}, testMetaKey, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, tc.info.GetMetadata(tc.key))
		})
	}
}

func TestInfoSetMetadata(t *testing.T) {
	assertions := assert.New(t)

	info := &Info{}
	info.SetMetadata(testMetaKey, testMetaValue)

	assertions.NotNil(info.Metadata, "Metadata should be initialized")
	assertions.Equal(testMetaValue, info.Metadata[testMetaKey])

	info.SetMetadata(testMetaKey2, testMetaValue2)
	assertions.Equal(testMetaValue2, info.Metadata[testMetaKey2])
}

/* ========================================================================================================
   REGISTRATION REQUEST TESTS
   ======================================================================================================== */

func TestRegistrationRequestValidate(t *testing.T) {
	testCases := []struct {
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
			req:     RegistrationRequest{Endpoint: &Info{ServiceID: testServiceID, Host: testHost}},
			wantErr: true,
		},
		{
			name:    "missing service ID",
			req:     RegistrationRequest{Endpoint: &Info{ID: testEndpointID, Host: testHost}},
			wantErr: true,
		},
		{
			name:    "missing host",
			req:     RegistrationRequest{Endpoint: &Info{ID: testEndpointID, ServiceID: testServiceID}},
			wantErr: true,
		},
		{
			name:    "valid",
			req:     RegistrationRequest{Endpoint: &Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost}},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			err := tc.req.Validate()
			if tc.wantErr {
				assertions.Error(err)
			} else {
				assertions.NoError(err)
			}
		})
	}
}

/* ========================================================================================================
   DEREGISTRATION REQUEST TESTS
   ======================================================================================================== */

func TestDeregistrationRequestValidate(t *testing.T) {
	testCases := []struct {
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
			req:     DeregistrationRequest{EndpointID: testEndpointID},
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			err := tc.req.Validate()
			if tc.wantErr {
				assertions.Error(err)
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
	endpoint := &Info{
		ID:         testEndpointID,
		ServiceID:  testServiceID,
		InstanceID: testInstanceID,
		Host:       testHost,
		Port:       testPort,
		Protocol:   ProtocolHTTP,
		Status:     StatusHealthy,
		Tags:       []string{testTagAPI, testTagV1},
		Metadata:   map[string]string{testMetaEnvKey: testMetaEnvProd},
		Version:    testVersion,
		APIVersion: testAPIVersion,
		TenantID:   testTenantID,
		Region:     testRegion,
		Zone:       testZone,
	}

	testCases := []struct {
		name    string
		query   Query
		ep      *Info
		matches bool
	}{
		{"null endpoint", Query{}, nil, false},
		{"empty query matches all", Query{}, endpoint, true},
		{"ID filter match", Query{IDs: []string{testEndpointID}}, endpoint, true},
		{"ID filter no match", Query{IDs: []string{testEndpointID2}}, endpoint, false},
		{"ServiceID filter match", Query{ServiceIDs: []string{testServiceID}}, endpoint, true},
		{"ServiceID filter no match", Query{ServiceIDs: []string{testServiceID2}}, endpoint, false},
		{"InstanceID filter match", Query{InstanceID: testInstanceID}, endpoint, true},
		{"InstanceID filter no match", Query{InstanceID: testInstanceID2}, endpoint, false},
		{"Host filter match", Query{Host: testHost}, endpoint, true},
		{"Host filter no match", Query{Host: "remotehost"}, endpoint, false},
		{"Port filter match", Query{Port: testPort}, endpoint, true},
		{"Port filter no match", Query{Port: testPortAlt}, endpoint, false},
		{"Protocol filter match", Query{Protocol: ProtocolHTTP}, endpoint, true},
		{"Protocol filter no match", Query{Protocol: ProtocolHTTPS}, endpoint, false},
		{"Statuses filter match", Query{Statuses: []Status{StatusHealthy, StatusDegraded}}, endpoint, true},
		{"Statuses filter no match", Query{Statuses: []Status{StatusUnhealthy}}, endpoint, false},
		{"OnlyHealthy match", Query{OnlyHealthy: true}, endpoint, true},
		{"OnlyAvailable match", Query{OnlyAvailable: true}, endpoint, true},
		{"Tag filter match", Query{Tags: []string{testTagAPI}}, endpoint, true},
		{"Tag filter no match", Query{Tags: []string{testTagGRPC}}, endpoint, false},
		{"Metadata filter match", Query{Metadata: map[string]string{testMetaEnvKey: testMetaEnvProd}}, endpoint, true},
		{"Metadata filter no match", Query{Metadata: map[string]string{testMetaEnvKey: testMetaEnvDev}}, endpoint, false},
		{"Version filter match", Query{Version: testVersion}, endpoint, true},
		{"Version filter no match", Query{Version: testVersionAlt}, endpoint, false},
		{"APIVersion filter match", Query{APIVersion: testAPIVersion}, endpoint, true},
		{"APIVersion filter no match", Query{APIVersion: testAPIVersionAlt}, endpoint, false},
		{"TenantID filter match", Query{TenantID: testTenantID}, endpoint, true},
		{"TenantID filter no match", Query{TenantID: testTenantID2}, endpoint, false},
		{"Region filter match", Query{Region: testRegion}, endpoint, true},
		{"Region filter no match", Query{Region: testRegionAlt}, endpoint, false},
		{"Zone filter match", Query{Zone: testZone}, endpoint, true},
		{"Zone filter no match", Query{Zone: testZoneAlt}, endpoint, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.matches, tc.query.Matches(tc.ep))
		})
	}
}

/* ========================================================================================================
   EVENT TESTS
   ======================================================================================================== */

func TestEventToMessage(t *testing.T) {
	assertions := assert.New(t)

	event := &Event{
		Type:       EventRegistered,
		EndpointID: testEndpointID,
		ServiceID:  testServiceID,
		Timestamp:  time.Now(),
	}

	msg, err := event.ToMessage()
	assertions.NoError(err)
	assertions.NotNil(msg)
	assertions.Equal(string(EventRegistered), msg.Header.Get("event-type"))
	assertions.Equal(testEndpointID, msg.Header.Get("endpoint-id"))
	assertions.Equal(testServiceID, msg.Header.Get("service-id"))
}

func TestEventFromMessage(t *testing.T) {
	assertions := assert.New(t)

	originalEvent := &Event{
		Type:       EventRegistered,
		EndpointID: testEndpointID,
		ServiceID:  testServiceID,
		Timestamp:  time.Now().UTC().Truncate(time.Second),
	}

	msg, _ := originalEvent.ToMessage()
	recovered, err := EventFromMessage(msg)
	assertions.NoError(err)
	assertions.Equal(originalEvent.Type, recovered.Type)
	assertions.Equal(originalEvent.EndpointID, recovered.EndpointID)
}

func TestEventFromMessageInvalid(t *testing.T) {
	assertions := assert.New(t)
	msg := newTestMsg([]byte(testInvalidJSON))
	_, err := EventFromMessage(msg)
	assertions.Error(err, "Expected error for invalid JSON")
}

/* ========================================================================================================
   REGISTRY TESTS
   ======================================================================================================== */

func TestNewRegistry(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()
	assertions.NotNil(r)
}

func TestRegistryRegisterAndGet(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	ep := &Info{
		ID:        testEndpointID,
		ServiceID: testServiceID,
		Host:      testHost,
		Port:      testPort,
		Protocol:  ProtocolHTTP,
		Status:    StatusHealthy,
	}

	err := r.Register(ep)
	assertions.NoError(err)

	retrieved, found := r.Get(testEndpointID)
	assertions.True(found, "Endpoint not found")
	assertions.Equal(ep.ID, retrieved.ID)
}

func TestRegistryRegisterDuplicate(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	ep := &Info{
		ID:        testEndpointID,
		ServiceID: testServiceID,
		Host:      testHost,
		Port:      testPort,
	}

	err := r.Register(ep)
	assertions.NoError(err)

	// Second registration with same ID should update, not error
	ep2 := &Info{
		ID:        testEndpointID,
		ServiceID: testServiceID,
		Host:      testHost,
		Port:      testPortAlt,
	}
	err = r.Register(ep2)
	assertions.NoError(err, "Second registration should succeed (update)")

	got, exists := r.Get(testEndpointID)
	assertions.True(exists)
	assertions.Equal(testPortAlt, got.Port, "Port should be updated")
}

func TestRegistryRegisterWithTTL(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	ep := newTestEndpoint(testEndpointID, testServiceID, testHost)

	err := r.RegisterWithTTL(ep, testTTL5Min)
	assertions.NoError(err)

	retrieved, _ := r.Get(testEndpointID)
	assertions.False(retrieved.ExpiresAt.IsZero(), "ExpiresAt should be set")
}

func TestRegistryGetNotFound(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	_, found := r.Get("nonexistent")
	assertions.False(found)
}

func TestRegistryDeregister(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	ep := newTestEndpoint(testEndpointID, testServiceID, testHost)
	r.Register(ep)

	err := r.Deregister(testEndpointID)
	assertions.NoError(err)

	_, found := r.Get(testEndpointID)
	assertions.False(found, "Endpoint should be removed")
}

func TestRegistryDeregisterNotFound(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	err := r.Deregister("nonexistent")
	assertions.Error(err)
}

func TestRegistryDeregisterWithReason(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	ep := newTestEndpoint(testEndpointID, testServiceID, testHost)
	r.Register(ep)

	err := r.DeregisterWithReason(testEndpointID, testReasonTest)
	assertions.NoError(err)
}

func TestRegistryDeregisterByService(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	r.Register(newTestEndpoint(testEndpointID, testServiceID, testHost1))
	r.Register(newTestEndpoint(testEndpointID2, testServiceID, testHost2))
	r.Register(newTestEndpoint(testEndpointID3, testServiceID2, testHost3))

	count, err := r.DeregisterByService(testServiceID)
	assertions.NoError(err)
	assertions.Equal(2, count)
	assertions.Equal(1, r.Count())
}

func TestRegistryDeregisterByInstance(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	r.Register(&Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost1, InstanceID: testInstanceID})
	r.Register(&Info{ID: testEndpointID2, ServiceID: testServiceID, Host: testHost2, InstanceID: testInstanceID})
	r.Register(&Info{ID: testEndpointID3, ServiceID: testServiceID2, Host: testHost3, InstanceID: testInstanceID2})

	count, err := r.DeregisterByInstance(testInstanceID)
	assertions.NoError(err)
	assertions.Equal(2, count)
}

func TestRegistryUpdateStatus(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	ep := &Info{
		ID:        testEndpointID,
		ServiceID: testServiceID,
		Host:      testHost,
		Status:    StatusHealthy,
	}
	r.Register(ep)

	err := r.UpdateStatus(testEndpointID, StatusDegraded)
	assertions.NoError(err)

	retrieved, _ := r.Get(testEndpointID)
	assertions.Equal(StatusDegraded, retrieved.Status)
}

func TestRegistryQuery(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	endpoints := []*Info{
		{ID: testEndpointID, ServiceID: testServiceID, Host: testHost1, Status: StatusHealthy, Tags: []string{testTagAPI}},
		{ID: testEndpointID2, ServiceID: testServiceID, Host: testHost2, Status: StatusUnhealthy, Tags: []string{testTagAPI}},
		{ID: testEndpointID3, ServiceID: testServiceID2, Host: testHost3, Status: StatusHealthy, Tags: []string{testTagGRPC}},
	}

	for _, ep := range endpoints {
		r.Register(ep)
	}

	// Query by service
	results := r.Query(&Query{ServiceIDs: []string{testServiceID}})
	assertions.Len(results, 2, "Query by service")

	// Query healthy only
	results = r.Query(&Query{OnlyHealthy: true})
	assertions.Len(results, 2, "Query healthy")

	// Query by tag
	results = r.Query(&Query{Tags: []string{testTagAPI}})
	assertions.Len(results, 2, "Query by tag")

	// Empty query
	results = r.Query(nil)
	assertions.Len(results, 3, "Empty query")
}

func TestRegistryQueryWithPagination(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	for i := 0; i < 10; i++ {
		r.Register(&Info{
			ID:        "ep" + itoa(i),
			ServiceID: "svc",
			Host:      "host",
		})
	}

	results := r.Query(&Query{Limit: 5})
	assertions.Len(results, 5)
}

func TestRegistryQueryWithOffset(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	for i := 0; i < 10; i++ {
		r.Register(&Info{
			ID:        "ep" + itoa(i),
			ServiceID: "svc",
			Host:      "host",
		})
	}

	results := r.Query(&Query{Limit: 3, Offset: 2})
	assertions.Len(results, 3)
}

func TestRegistryAll(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	r.Register(newTestEndpoint(testEndpointID, testServiceID, testHost1))
	r.Register(newTestEndpoint(testEndpointID2, testServiceID, testHost2))

	all := r.All()
	assertions.Len(all, 2)
}

func TestRegistryCount(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	assertions.Equal(0, r.Count(), "Empty registry should have count 0")

	r.Register(newTestEndpoint(testEndpointID, testServiceID, testHost1))
	r.Register(newTestEndpoint(testEndpointID2, testServiceID, testHost2))

	assertions.Equal(2, r.Count())
}

func TestRegistryServiceCount(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	assertions.Equal(0, r.ServiceCount())

	r.Register(newTestEndpoint(testEndpointID, testServiceID, testHost1))
	r.Register(newTestEndpoint(testEndpointID2, testServiceID2, testHost2))
	r.Register(newTestEndpoint(testEndpointID3, testServiceID, testHost3))

	assertions.Equal(2, r.ServiceCount())
}

func TestRegistryServices(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	r.Register(newTestEndpoint(testEndpointID, testServiceID, testHost1))
	r.Register(newTestEndpoint(testEndpointID2, testServiceID2, testHost2))
	r.Register(newTestEndpoint(testEndpointID3, testServiceID, testHost3))

	services := r.Services()
	assertions.Len(services, 2)
}

func TestRegistryExists(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	r.Register(newTestEndpoint(testEndpointID, testServiceID, testHost1))

	assertions.True(r.Exists(testEndpointID))
	assertions.False(r.Exists("nonexistent"))
}

func TestRegistryGetByService(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	r.Register(newTestEndpoint(testEndpointID, testServiceID, testHost1))
	r.Register(newTestEndpoint(testEndpointID2, testServiceID2, testHost2))
	r.Register(newTestEndpoint(testEndpointID3, testServiceID, testHost3))

	endpoints := r.GetByService(testServiceID)
	assertions.Len(endpoints, 2)
}

func TestRegistryGetByInstance(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	r.Register(&Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost1, InstanceID: testInstanceID})
	r.Register(&Info{ID: testEndpointID2, ServiceID: testServiceID, Host: testHost2, InstanceID: testInstanceID})
	r.Register(&Info{ID: testEndpointID3, ServiceID: testServiceID2, Host: testHost3, InstanceID: testInstanceID2})

	endpoints := r.GetByInstance(testInstanceID)
	assertions.Len(endpoints, 2)

	endpoints = r.GetByInstance("nonexistent")
	assertions.Empty(endpoints)
}

func TestRegistryGetByTag(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	r.Register(&Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost1, Tags: []string{testTagAPI, testTagV1}})
	r.Register(&Info{ID: testEndpointID2, ServiceID: testServiceID, Host: testHost2, Tags: []string{testTagAPI, testTagV2}})
	r.Register(&Info{ID: testEndpointID3, ServiceID: testServiceID2, Host: testHost3, Tags: []string{testTagGRPC}})

	assertions.Len(r.GetByTag(testTagAPI), 2)
	assertions.Len(r.GetByTag(testTagV1), 1)
	assertions.Empty(r.GetByTag("nonexistent"))
}

func TestRegistryGetHealthy(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	r.Register(&Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost1, Status: StatusHealthy})
	r.Register(&Info{ID: testEndpointID2, ServiceID: testServiceID, Host: testHost2, Status: StatusUnhealthy})
	r.Register(&Info{ID: testEndpointID3, ServiceID: testServiceID2, Host: testHost3, Status: StatusHealthy})

	assertions.Len(r.GetHealthy(), 2)
}

func TestRegistryGetAvailable(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	r.Register(&Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost1, Status: StatusHealthy})
	r.Register(&Info{ID: testEndpointID2, ServiceID: testServiceID, Host: testHost2, Status: StatusDegraded})
	r.Register(&Info{ID: testEndpointID3, ServiceID: testServiceID2, Host: testHost3, Status: StatusUnhealthy})

	assertions.Len(r.GetAvailable(), 2)
}

func TestRegistryUpdateHealth(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	r.Register(&Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost1, Status: StatusHealthy})

	err := r.UpdateHealth(&HealthUpdate{
		EndpointID: testEndpointID,
		Status:     StatusHealthy,
		CheckedAt:  time.Now(),
	})
	assertions.NoError(err)

	ep, _ := r.Get(testEndpointID)
	assertions.Equal(StatusHealthy, ep.Status)

	err = r.UpdateHealth(&HealthUpdate{
		EndpointID: testEndpointID,
		Status:     StatusUnhealthy,
		Message:    "health check failed",
		CheckedAt:  time.Now(),
	})
	assertions.NoError(err)

	ep, _ = r.Get(testEndpointID)
	assertions.Equal(StatusUnhealthy, ep.Status)

	// Test not found
	err = r.UpdateHealth(&HealthUpdate{
		EndpointID: "nonexistent",
		Status:     StatusHealthy,
		CheckedAt:  time.Now(),
	})
	assertions.Error(err)

	// Test nil update
	err = r.UpdateHealth(nil)
	assertions.Error(err)
}

func TestRegistryRenew(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	originalExpiry := time.Now().Add(testTTL1Min)
	r.Register(&Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost1, ExpiresAt: originalExpiry})

	err := r.Renew(testEndpointID, testTTL5Min)
	assertions.NoError(err)

	ep, _ := r.Get(testEndpointID)
	assertions.True(ep.ExpiresAt.After(originalExpiry), "ExpiresAt should be extended")

	// Test not found
	err = r.Renew("nonexistent", testTTL5Min)
	assertions.Error(err)
}

func TestRegistryTouch(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	r.Register(newTestEndpoint(testEndpointID, testServiceID, testHost1))

	ep, _ := r.Get(testEndpointID)
	oldUpdated := ep.UpdatedAt

	time.Sleep(testCallbackWait)

	err := r.Touch(testEndpointID)
	assertions.NoError(err)

	ep, _ = r.Get(testEndpointID)
	assertions.True(ep.UpdatedAt.After(oldUpdated), "UpdatedAt should be refreshed")

	// Test not found
	err = r.Touch("nonexistent")
	assertions.Error(err)
}

func TestRegistryClear(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	r.Register(newTestEndpoint(testEndpointID, testServiceID, testHost1))
	r.Register(newTestEndpoint(testEndpointID2, testServiceID2, testHost2))

	r.Clear()
	assertions.Equal(0, r.Count())
}

func TestRegistryShutdown(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()
	r.Register(newTestEndpoint(testEndpointID, testServiceID, testHost1))

	ctx := context.Background()
	err := r.Shutdown(ctx)
	assertions.NoError(err)

	// After shutdown, register should fail
	err = r.Register(newTestEndpoint(testEndpointID2, testServiceID2, testHost2))
	assertions.Error(err, "Register should fail after shutdown")
}

func TestRegistryRegisterValidation(t *testing.T) {
	testCases := []struct {
		name string
		ep   *Info
	}{
		{"a nil endpoint", nil},
		{"missing Id", &Info{ServiceID: "svc", Host: "host"}},
		{"missing service Id", &Info{ID: testEndpointID, Host: "host"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			r := NewRegistry()
			err := r.Register(tc.ep)
			assertions.Error(err)
		})
	}
}

/* ========================================================================================================
   REGISTRY CALLBACK TESTS
   ======================================================================================================== */

func TestRegistryCallbacks(t *testing.T) {
	assertions := assert.New(t)
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
		ID:        testEndpointID,
		ServiceID: testServiceID,
		Host:      testHost,
		Status:    StatusHealthy,
	}

	r.Register(ep)
	time.Sleep(testCallbackWait)
	assertions.True(registeredCalled, "OnRegistered callback not called")

	r.UpdateStatus(testEndpointID, StatusDegraded)
	time.Sleep(testCallbackWait)
	assertions.True(statusChangeCalled, "OnStatusChange callback not called")

	r.Deregister(testEndpointID)
	time.Sleep(testCallbackWait)
	assertions.True(deregisteredCalled, "OnDeregistered callback not called")
}

func TestRegistryOnUpdatedCallback(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	updatedCalled := false
	r.OnUpdated(func(old, new *Info) {
		updatedCalled = true
	})

	ep := &Info{
		ID:        testEndpointID,
		ServiceID: testServiceID,
		Host:      testHost,
		Port:      testPort,
	}
	r.Register(ep)

	// Update by re-registering with different port
	ep2 := &Info{
		ID:        testEndpointID,
		ServiceID: testServiceID,
		Host:      testHost,
		Port:      testPortAlt,
	}
	r.Register(ep2)

	time.Sleep(testCallbackWait)
	assertions.True(updatedCalled, "OnUpdated callback not called")
}

func TestRegistryDeregisterWithReasonCallback(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	deregisteredCalled := false
	var gotReason string
	r.OnDeregistered(func(ep *Info, reason string) {
		deregisteredCalled = true
		gotReason = reason
	})

	r.Register(newTestEndpoint(testEndpointID, testServiceID, testHost))
	r.DeregisterWithReason(testEndpointID, testReasonShutdown)

	time.Sleep(testCallbackWait)
	assertions.True(deregisteredCalled, "OnDeregistered callback not called")
	assertions.Equal(testReasonShutdown, gotReason)
}

/* ========================================================================================================
   REGISTRY CONCURRENT ACCESS TESTS
   ======================================================================================================== */

func TestRegistryConcurrentAccess(t *testing.T) {
	assertions := assert.New(t)
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

	assertions.Equal(100, r.Count())
}

/* ========================================================================================================
   ERROR TESTS
   ======================================================================================================== */

func TestErrors(t *testing.T) {
	testCases := []struct {
		name    string
		err     error
		message string
	}{
		{"ErrInvalidEndpoint", ErrInvalidEndpoint, "invalid endpoint"},
		{"ErrMissingEndpointID", ErrMissingEndpointID, "missing endpoint ID"},
		{"ErrMissingServiceID", ErrMissingServiceID, "missing service ID"},
		{"ErrMissingHost", ErrMissingHost, "missing host"},
		{"ErrEndpointNotFound", ErrEndpointNotFound, "endpoint not found"},
		{"ErrEndpointExists", ErrEndpointExists, "endpoint already exists"},
		{"ErrEndpointExpired", ErrEndpointExpired, "endpoint has expired"},
		{"ErrHealthCheckFailed", ErrHealthCheckFailed, "health check failed"},
		{"ErrHealthCheckTimeout", ErrHealthCheckTimeout, "health check timeout"},
		{"ErrRegistryShutdown", ErrRegistryShutdown, "registry is shutting down"},
		{"ErrInvalidQuery", ErrInvalidQuery, "invalid query parameters"},
		{"ErrManagerNotStarted", ErrManagerNotStarted, "manager not started"},
		{"ErrManagerShutdown", ErrManagerShutdown, "manager is shutting down"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.NotNil(tc.err)
			assertions.Equal(tc.message, tc.err.Error())
		})
	}
}

func TestEndpointError(t *testing.T) {
	assertions := assert.New(t)

	err := NewEndpointError(testEndpointID, testServiceID, "register", ErrEndpointExists)
	assertions.NotEmpty(err.Error())
	assertions.Equal(ErrEndpointExists, err.Unwrap())
}

func TestEndpointErrorWithoutServiceID(t *testing.T) {
	assertions := assert.New(t)
	err := NewEndpointError(testEndpointID, "", "register", ErrEndpointExists)
	assertions.NotEmpty(err.Error())
}

func TestEndpointErrorWithoutOp(t *testing.T) {
	assertions := assert.New(t)
	err := NewEndpointError(testEndpointID, "", "", ErrEndpointExists)
	assertions.NotEmpty(err.Error())
}

func TestRegistrationError(t *testing.T) {
	assertions := assert.New(t)

	err := RegistrationError{
		EndpointID: testEndpointID,
		Reason:     testReasonDuplicate,
		Err:        ErrEndpointExists,
	}

	assertions.NotEmpty(err.Error())
	assertions.Equal(ErrEndpointExists, err.Unwrap())
}

func TestRegistrationErrorWithoutReason(t *testing.T) {
	assertions := assert.New(t)

	err := RegistrationError{
		EndpointID: testEndpointID,
		Err:        ErrEndpointExists,
	}
	assertions.NotEmpty(err.Error())
}

func TestHealthCheckError(t *testing.T) {
	assertions := assert.New(t)

	err := HealthCheckError{
		EndpointID:       testEndpointID,
		URL:              testHealthURL,
		ConsecutiveFails: 3,
		Err:              ErrHealthCheckTimeout,
	}

	assertions.NotEmpty(err.Error())
	assertions.Equal(ErrHealthCheckTimeout, err.Unwrap())
}

func TestIsEndpointError(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected bool
	}{
		{"EndpointError", NewEndpointError(testEndpointID, testServiceID, "register", ErrEndpointExists), true},
		{"RegistrationError", RegistrationError{EndpointID: testEndpointID, Err: ErrEndpointExists}, true},
		{"HealthCheckError", HealthCheckError{EndpointID: testEndpointID, Err: ErrHealthCheckFailed}, true},
		{"base error", ErrInvalidEndpoint, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, IsEndpointError(tc.err))
		})
	}
}

func TestIsNotFound(t *testing.T) {
	assertions := assert.New(t)
	assertions.True(IsNotFound(ErrEndpointNotFound))
	assertions.False(IsNotFound(ErrEndpointExists))
}

func TestIsExists(t *testing.T) {
	assertions := assert.New(t)
	assertions.True(IsExists(ErrEndpointExists))
	assertions.False(IsExists(ErrEndpointNotFound))
}

/* ========================================================================================================
   MANAGER TESTS
   ======================================================================================================== */

func TestNewManager(t *testing.T) {
	assertions := assert.New(t)

	manager, err := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{Enabled: false},
	})
	assertions.NoError(err)
	assertions.NotNil(manager)
}

func TestManagerConfigDefaults(t *testing.T) {
	assertions := assert.New(t)

	manager, err := NewManager(ManagerConfig{})
	assertions.NoError(err)
	assertions.NotNil(manager)
}

func TestManagerWithHealthCheckEnabled(t *testing.T) {
	assertions := assert.New(t)

	manager, err := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{
			Enabled:  true,
			Interval: 1 * time.Second,
			Timeout:  500 * time.Millisecond,
		},
	})
	assertions.NoError(err)

	ctx := context.Background()
	manager.Shutdown(ctx)
}

func TestManagerRegisterEndpoint(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	ep := &Info{
		ID:        testEndpointID,
		ServiceID: testServiceID,
		Host:      testHost,
		Port:      testPort,
		Protocol:  ProtocolHTTP,
		Status:    StatusHealthy,
	}

	err := manager.RegisterEndpoint(ctx, ep)
	assertions.NoError(err)

	retrieved, found := manager.Get(testEndpointID)
	assertions.True(found)
	assertions.Equal(ep.ID, retrieved.ID)
}

func TestManagerRegister(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	req := &RegistrationRequest{
		Endpoint: newTestEndpoint(testEndpointID, testServiceID, testHost),
	}

	err := manager.Register(ctx, req)
	assertions.NoError(err)

	_, found := manager.Get(testEndpointID)
	assertions.True(found)
}

func TestManagerRegisterWithTTL(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	ep := newTestEndpoint(testEndpointID, testServiceID, testHost)

	err := manager.RegisterWithTTL(ctx, ep, testTTL5Min)
	assertions.NoError(err)

	retrieved, _ := manager.Get(testEndpointID)
	assertions.False(retrieved.ExpiresAt.IsZero(), "ExpiresAt should be set")
}

func TestManagerRegisterValidationError(t *testing.T) {
	testCases := []struct {
		name string
		req  *RegistrationRequest
	}{
		{"missing endpoint", &RegistrationRequest{}},
		{"missing id", &RegistrationRequest{Endpoint: &Info{ServiceID: "svc", Host: "host"}}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			manager := newTestManager()
			err := manager.Register(context.Background(), tc.req)
			assertions.Error(err)
		})
	}
}

func TestManagerDeregisterEndpoint(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, newTestEndpoint(testEndpointID, testServiceID, testHost))

	err := manager.DeregisterEndpoint(ctx, testEndpointID)
	assertions.NoError(err)

	_, found := manager.Get(testEndpointID)
	assertions.False(found, "Endpoint should be removed")
}

func TestManagerDeregister(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, newTestEndpoint(testEndpointID, testServiceID, testHost1))

	err := manager.Deregister(ctx, &DeregistrationRequest{
		EndpointID: testEndpointID,
		Reason:     testReasonDeregistration,
	})
	assertions.NoError(err)

	_, found := manager.Get(testEndpointID)
	assertions.False(found)
}

func TestManagerDeregisterValidation(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	err := manager.Deregister(context.Background(), &DeregistrationRequest{})
	assertions.Error(err, "Deregister should fail without endpoint ID")
}

func TestManagerDeregisterService(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, newTestEndpoint(testEndpointID, testServiceID, testHost1))
	manager.RegisterEndpoint(ctx, newTestEndpoint(testEndpointID2, testServiceID, testHost2))
	manager.RegisterEndpoint(ctx, newTestEndpoint(testEndpointID3, testServiceID2, testHost3))

	count, err := manager.DeregisterService(ctx, testServiceID)
	assertions.NoError(err)
	assertions.Equal(2, count)
	assertions.Equal(1, manager.Count())
}

func TestManagerDeregisterInstance(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost1, InstanceID: testInstanceID})
	manager.RegisterEndpoint(ctx, &Info{ID: testEndpointID2, ServiceID: testServiceID, Host: testHost2, InstanceID: testInstanceID})
	manager.RegisterEndpoint(ctx, &Info{ID: testEndpointID3, ServiceID: testServiceID2, Host: testHost3, InstanceID: testInstanceID2})

	count, err := manager.DeregisterInstance(ctx, testInstanceID)
	assertions.NoError(err)
	assertions.Equal(2, count)
	assertions.Equal(1, manager.Count())
}

func TestManagerDiscover(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost1, Status: StatusHealthy})
	manager.RegisterEndpoint(ctx, &Info{ID: testEndpointID2, ServiceID: testServiceID, Host: testHost2, Status: StatusDegraded})
	manager.RegisterEndpoint(ctx, &Info{ID: testEndpointID3, ServiceID: testServiceID, Host: testHost3, Status: StatusUnhealthy})

	// Discover available (healthy + degraded)
	endpoints := manager.Discover(testServiceID)
	assertions.Len(endpoints, 2, "Discover should return healthy + degraded")

	// Discover healthy only
	endpoints = manager.Discover(testServiceID, OnlyHealthy())
	assertions.Len(endpoints, 1, "Discover(OnlyHealthy) should return only healthy")
}

func TestManagerQuery(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost1, Status: StatusHealthy})
	manager.RegisterEndpoint(ctx, &Info{ID: testEndpointID2, ServiceID: testServiceID, Host: testHost2, Status: StatusHealthy})

	results := manager.Query(&Query{ServiceIDs: []string{testServiceID}})
	assertions.Len(results, 2)
}

func TestManagerGetByService(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, newTestEndpoint(testEndpointID, testServiceID, testHost1))
	manager.RegisterEndpoint(ctx, newTestEndpoint(testEndpointID2, testServiceID, testHost2))
	manager.RegisterEndpoint(ctx, newTestEndpoint(testEndpointID3, testServiceID2, testHost3))

	assertions.Len(manager.GetByService(testServiceID), 2)
}

func TestManagerGetByInstance(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost1, InstanceID: testInstanceID})
	manager.RegisterEndpoint(ctx, &Info{ID: testEndpointID2, ServiceID: testServiceID, Host: testHost2, InstanceID: testInstanceID})

	assertions.Len(manager.GetByInstance(testInstanceID), 2)
}

func TestManagerGetHealthy(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost1, Status: StatusHealthy})
	manager.RegisterEndpoint(ctx, &Info{ID: testEndpointID2, ServiceID: testServiceID, Host: testHost2, Status: StatusUnhealthy})

	assertions.Len(manager.GetHealthy(), 1)
}

func TestManagerGetAvailable(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost1, Status: StatusHealthy})
	manager.RegisterEndpoint(ctx, &Info{ID: testEndpointID2, ServiceID: testServiceID, Host: testHost2, Status: StatusDegraded})
	manager.RegisterEndpoint(ctx, &Info{ID: testEndpointID3, ServiceID: testServiceID, Host: testHost3, Status: StatusUnhealthy})

	assertions.Len(manager.GetAvailable(), 2)
}

func TestManagerUpdateStatus(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost1, Status: StatusHealthy})

	err := manager.UpdateStatus(ctx, testEndpointID, StatusDegraded)
	assertions.NoError(err)

	ep, _ := manager.Get(testEndpointID)
	assertions.Equal(StatusDegraded, ep.Status)
}

func TestManagerRenew(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, newTestEndpoint(testEndpointID, testServiceID, testHost1))

	err := manager.Renew(ctx, testEndpointID, testTTL5Min)
	assertions.NoError(err)

	ep, _ := manager.Get(testEndpointID)
	assertions.False(ep.ExpiresAt.IsZero(), "ExpiresAt should be set after Renew")
}

func TestManagerAll(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, newTestEndpoint(testEndpointID, testServiceID, testHost1))
	manager.RegisterEndpoint(ctx, newTestEndpoint(testEndpointID2, testServiceID2, testHost2))

	assertions.Len(manager.All(), 2)
}

func TestManagerCount(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	assertions.Equal(0, manager.Count())

	manager.RegisterEndpoint(ctx, newTestEndpoint(testEndpointID, testServiceID, testHost1))
	manager.RegisterEndpoint(ctx, newTestEndpoint(testEndpointID2, testServiceID, testHost2))

	assertions.Equal(2, manager.Count())
}

func TestManagerServices(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, newTestEndpoint(testEndpointID, testServiceID, testHost1))
	manager.RegisterEndpoint(ctx, newTestEndpoint(testEndpointID2, testServiceID2, testHost2))

	assertions.Len(manager.Services(), 2)
}

func TestManagerRegistry(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()
	assertions.NotNil(manager.Registry())
}

func TestManagerSetPublisher(t *testing.T) {
	manager := newTestManager()
	manager.SetPublisher(nil) // Should not panic
}

func TestManagerOnError(t *testing.T) {
	manager := newTestManager()
	manager.OnError(func(err error) {
		// Verify callback is set without panic
	})
}

func TestManagerErrors(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()
	assertions.NotNil(manager.Errors())
}

func TestManagerShutdown(t *testing.T) {
	manager := newTestManager()
	ctx := context.Background()
	manager.Shutdown(ctx)
	// Should not panic
}

/* ========================================================================================================
   DISCOVERY OPTIONS TESTS
   ======================================================================================================== */

func TestDiscoveryOptions(t *testing.T) {
	testCases := []struct {
		name   string
		apply  func(q *Query)
		verify func(assertions *assert.Assertions, q *Query)
	}{
		{
			name:  "OnlyHealthy",
			apply: func(q *Query) { OnlyHealthy()(q) },
			verify: func(assertions *assert.Assertions, q *Query) {
				assertions.True(q.OnlyHealthy)
			},
		},
		{
			name:  "WithTags",
			apply: func(q *Query) { WithTags("tag1", "tag2")(q) },
			verify: func(assertions *assert.Assertions, q *Query) {
				assertions.Len(q.Tags, 2)
			},
		},
		{
			name:  "WithVersion",
			apply: func(q *Query) { WithVersion(testVersion)(q) },
			verify: func(assertions *assert.Assertions, q *Query) {
				assertions.Equal(testVersion, q.Version)
			},
		},
		{
			name:  "WithRegion",
			apply: func(q *Query) { WithRegion(testRegion)(q) },
			verify: func(assertions *assert.Assertions, q *Query) {
				assertions.Equal(testRegion, q.Region)
			},
		},
		{
			name:  "WithLimit",
			apply: func(q *Query) { WithLimit(10)(q) },
			verify: func(assertions *assert.Assertions, q *Query) {
				assertions.Equal(10, q.Limit)
			},
		},
		{
			name:  "WithProtocol",
			apply: func(q *Query) { WithProtocol(ProtocolHTTP)(q) },
			verify: func(assertions *assert.Assertions, q *Query) {
				assertions.Equal(ProtocolHTTP, q.Protocol)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			q := &Query{}
			tc.apply(q)
			tc.verify(assertions, q)
		})
	}
}

func TestDiscoverOptionsChaining(t *testing.T) {
	assertions := assert.New(t)
	q := &Query{}

	WithTags(testTagAPI, testTagV1)(q)
	WithVersion(testVersion)(q)
	WithProtocol(ProtocolHTTP)(q)
	WithRegion(testRegion)(q)
	WithLimit(10)(q)
	OnlyHealthy()(q)

	assertions.Len(q.Tags, 2)
	assertions.Equal(testVersion, q.Version)
	assertions.Equal(ProtocolHTTP, q.Protocol)
	assertions.Equal(testRegion, q.Region)
	assertions.Equal(10, q.Limit)
	assertions.True(q.OnlyHealthy)
}

/* ========================================================================================================
   HANDLER SUBJECT TESTS
   ======================================================================================================== */

func TestDefaultSubjects(t *testing.T) {
	assertions := assert.New(t)
	subjects := DefaultSubjects()

	assertions.Equal(testSubjectRegister, subjects.Register)
	assertions.Equal(testSubjectDeregister, subjects.Deregister)
	assertions.Equal(testSubjectRenew, subjects.Renew)
	assertions.Equal(testSubjectQuery, subjects.Query)
	assertions.Equal(testSubjectDiscover, subjects.Discover)
	assertions.Equal(testSubjectEvents, subjects.Events)
	assertions.Equal(testSubjectHealthUpdate, subjects.HealthUpdate)
}

func TestTenantSubjects(t *testing.T) {
	assertions := assert.New(t)
	subjects := TenantSubjects(testTenantID)

	assertions.Equal(testTenantID+"."+testSubjectRegister, subjects.Register)
	assertions.Equal(testTenantID+"."+testSubjectDeregister, subjects.Deregister)
}

/* ========================================================================================================
   HANDLER TESTS
   ======================================================================================================== */

func TestNewHandler(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()
	handler := NewHandler(manager)
	assertions.NotNil(handler)
}

func TestNewHandlerWithOptions(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	customSubjects := Subjects{
		Register:   "custom.register",
		Deregister: "custom.deregister",
	}
	handler := NewHandler(manager, WithSubjects(customSubjects))
	assertions.NotNil(handler)
}

func TestWithPublisherOption(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()
	pub := &mockPublisher{}
	handler := NewHandler(manager, WithPublisher(pub))
	assertions.NotNil(handler)
}

func TestNewHandlerWithPublisher(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	handler := NewHandlerWithPublisher(manager, nil)
	assertions.NotNil(handler)

	// With custom subjects
	subjects := TenantSubjects(testTenantID)
	handler2 := NewHandlerWithPublisher(manager, nil, subjects)
	assertions.NotNil(handler2)
}

func TestHandlerSetPublisher(t *testing.T) {
	manager := newTestManager()
	handler := NewHandler(manager)
	handler.SetPublisher(nil) // Should not panic
}

func TestHandlerHandleRegister(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()
	handler := NewHandler(manager)

	req := RegistrationRequest{
		Endpoint: &Info{
			ID:        testEndpointID,
			ServiceID: testServiceID,
			Host:      testHost,
			Port:      testPort,
		},
	}
	data, _ := json.Marshal(req)
	msg := newTestMsg(data)

	handler.HandleRegister(context.Background(), msg)

	_, found := manager.Get(testEndpointID)
	assertions.True(found, "Endpoint should be registered")
}

func TestHandlerHandleRegisterInvalid(t *testing.T) {
	manager := newTestManager()
	handler := NewHandler(manager)
	msg := newTestMsg([]byte(testInvalidJSON))
	handler.HandleRegister(context.Background(), msg)
	// Should not panic
}

func TestHandlerHandleDeregister(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, newTestEndpoint(testEndpointID, testServiceID, testHost))

	handler := NewHandler(manager)
	req := DeregistrationRequest{EndpointID: testEndpointID, Reason: testReasonTest}
	data, _ := json.Marshal(req)
	msg := newTestMsg(data)

	handler.HandleDeregister(ctx, msg)

	_, found := manager.Get(testEndpointID)
	assertions.False(found, "Endpoint should be deregistered")
}

func TestHandlerHandleDeregisterInvalid(t *testing.T) {
	manager := newTestManager()
	handler := NewHandler(manager)
	msg := newTestMsg([]byte(testInvalidPayload))
	handler.HandleDeregister(context.Background(), msg)
	// Should not panic
}

func TestHandlerHandleRenew(t *testing.T) {
	manager := newTestManager()
	ctx := context.Background()
	manager.RegisterEndpoint(ctx, newTestEndpoint(testEndpointID, testServiceID, testHost))

	handler := NewHandler(manager)
	req := RenewRequest{EndpointID: testEndpointID, TTL: testTTL5Min}
	data, _ := json.Marshal(req)
	msg := newTestMsg(data)

	handler.HandleRenew(ctx, msg)
	// Should not panic
}

func TestHandlerHandleRenewInvalid(t *testing.T) {
	manager := newTestManager()
	handler := NewHandler(manager)
	msg := newTestMsg([]byte(testInvalidPayload))
	handler.HandleRenew(context.Background(), msg)
	// Should not panic
}

func TestHandlerHandleQuery(t *testing.T) {
	manager := newTestManager()
	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost, Status: StatusHealthy})

	handler := NewHandler(manager)
	req := Query{ServiceIDs: []string{testServiceID}}
	data, _ := json.Marshal(req)
	msg := newTestMsg(data)

	handler.HandleQuery(ctx, msg)
	// Should not panic
}

func TestHandlerHandleQueryInvalid(t *testing.T) {
	manager := newTestManager()
	handler := NewHandler(manager)
	msg := newTestMsg([]byte(testInvalidPayload))
	handler.HandleQuery(context.Background(), msg)
	// Should not panic
}

func TestHandlerHandleDiscover(t *testing.T) {
	manager := newTestManager()
	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost, Status: StatusHealthy})

	handler := NewHandler(manager)
	req := DiscoverRequest{ServiceID: testServiceID}
	data, _ := json.Marshal(req)
	msg := newTestMsg(data)

	handler.HandleDiscover(ctx, msg)
	// Should not panic
}

func TestHandlerHandleDiscoverInvalid(t *testing.T) {
	manager := newTestManager()
	handler := NewHandler(manager)
	msg := newTestMsg([]byte(testInvalidPayload))
	handler.HandleDiscover(context.Background(), msg)
	// Should not panic
}

func TestHandlerHandleHealthUpdate(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost, Status: StatusHealthy})

	handler := NewHandler(manager)
	update := HealthUpdate{
		EndpointID: testEndpointID,
		Status:     StatusDegraded,
		Message:    "health degraded",
		CheckedAt:  time.Now(),
	}
	data, _ := json.Marshal(update)
	msg := newTestMsg(data)

	handler.HandleHealthUpdate(ctx, msg)

	ep, _ := manager.Get(testEndpointID)
	assertions.Equal(StatusDegraded, ep.Status)
}

func TestHandlerHandleHealthUpdateInvalid(t *testing.T) {
	manager := newTestManager()
	handler := NewHandler(manager)
	msg := newTestMsg([]byte(testInvalidPayload))
	handler.HandleHealthUpdate(context.Background(), msg)
	// Should not panic
}

/* ========================================================================================================
   CLIENT TESTS
   ======================================================================================================== */

func TestNewClient(t *testing.T) {
	assertions := assert.New(t)

	client := NewClient(nil)
	assertions.NotNil(client)

	// With custom subjects
	subjects := TenantSubjects(testTenantID)
	client2 := NewClient(nil, subjects)
	assertions.NotNil(client2)
}

func TestClientSetTimeout(t *testing.T) {
	client := NewClient(nil)
	client.SetTimeout(30 * time.Second)
	// Should not panic
}

/* ========================================================================================================
   HEALTH CHECKER TESTS
   ======================================================================================================== */

func TestNewHealthChecker(t *testing.T) {
	assertions := assert.New(t)

	cfg := HealthCheckerConfig{
		Timeout:            5 * time.Second,
		HealthyThreshold:   2,
		UnhealthyThreshold: 3,
	}
	checker := NewHealthChecker(cfg)
	assertions.NotNil(checker)
}

func TestNewHealthCheckerDefaults(t *testing.T) {
	assertions := assert.New(t)
	cfg := HealthCheckerConfig{}
	checker := NewHealthChecker(cfg)
	assertions.NotNil(checker)
}

func TestHealthCheckConfig(t *testing.T) {
	assertions := assert.New(t)

	cfg := HealthCheckConfig{
		Enabled:            true,
		Interval:           30 * time.Second,
		Timeout:            5 * time.Second,
		HealthyThreshold:   2,
		UnhealthyThreshold: 3,
		InitialDelay:       5 * time.Second,
	}

	assertions.True(cfg.Enabled)
	assertions.Equal(30*time.Second, cfg.Interval)
	assertions.Equal(2, cfg.HealthyThreshold)
}

func TestDefaultHealthCheckConfig(t *testing.T) {
	assertions := assert.New(t)
	cfg := DefaultHealthCheckConfig()

	assertions.True(cfg.Enabled)
	assertions.NotZero(cfg.Interval)
	assertions.NotZero(cfg.Timeout)
}

/* ========================================================================================================
   MANAGER HTTP HANDLER TESTS
   ======================================================================================================== */

func TestManagerHTTPHandleRegistration(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	req := RegistrationRequest{
		Endpoint: newTestEndpoint(testEndpointID, testServiceID, testHost),
		TTL:      testTTL5Min,
	}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, testHTTPPathRegister, bytes.NewReader(body))
	httpReq.Header.Set(testContentTypeKey, testContentType)
	rr := httptest.NewRecorder()

	manager.HandleRegistration(rr, httpReq)
	assertions.Contains([]int{http.StatusOK, http.StatusCreated}, rr.Code)
}

func TestManagerHTTPHandleRegistrationWrongMethod(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	httpReq := httptest.NewRequest(http.MethodGet, testHTTPPathRegister, nil)
	rr := httptest.NewRecorder()

	manager.HandleRegistration(rr, httpReq)
	assertions.Equal(http.StatusMethodNotAllowed, rr.Code)
}

func TestManagerHTTPHandleDeregistration(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, newTestEndpoint(testEndpointID, testServiceID, testHost))

	req := DeregistrationRequest{EndpointID: testEndpointID, Reason: testReasonTest}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, testHTTPPathDeregister, bytes.NewReader(body))
	httpReq.Header.Set(testContentTypeKey, testContentType)
	rr := httptest.NewRecorder()

	manager.HandleDeregistration(rr, httpReq)
	assertions.Equal(http.StatusOK, rr.Code)
}

func TestManagerHTTPHandleDeregistrationWrongMethod(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	httpReq := httptest.NewRequest(http.MethodGet, testHTTPPathDeregister, nil)
	rr := httptest.NewRecorder()

	manager.HandleDeregistration(rr, httpReq)
	assertions.Equal(http.StatusMethodNotAllowed, rr.Code)
}

func TestManagerHTTPHandleQuery(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost, Status: StatusHealthy})

	httpReq := httptest.NewRequest(http.MethodGet, testHTTPPathQuery+"?service_id="+testServiceID, nil)
	rr := httptest.NewRecorder()

	manager.HandleQuery(rr, httpReq)
	assertions.Equal(http.StatusOK, rr.Code)
}

func TestManagerHTTPHandleQueryWrongMethod(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	httpReq := httptest.NewRequest(http.MethodDelete, testHTTPPathQuery, nil)
	rr := httptest.NewRecorder()

	manager.HandleQuery(rr, httpReq)
	assertions.Equal(http.StatusMethodNotAllowed, rr.Code)
}

func TestManagerHTTPHandleRenew(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, newTestEndpoint(testEndpointID, testServiceID, testHost))

	req := RenewRequest{EndpointID: testEndpointID, TTL: testTTL10Min}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, testHTTPPathRenew, bytes.NewReader(body))
	httpReq.Header.Set(testContentTypeKey, testContentType)
	rr := httptest.NewRecorder()

	manager.HandleRenew(rr, httpReq)
	assertions.Equal(http.StatusOK, rr.Code)
}

func TestManagerHTTPHandleRenewWrongMethod(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	httpReq := httptest.NewRequest(http.MethodGet, testHTTPPathRenew, nil)
	rr := httptest.NewRecorder()

	manager.HandleRenew(rr, httpReq)
	assertions.Equal(http.StatusMethodNotAllowed, rr.Code)
}

/* ========================================================================================================
   MULTI-TENANT ISOLATION TESTS
   ======================================================================================================== */

func TestTenantIsolationRegistryQuery(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	// Register endpoints for different tenants
	r.Register(&Info{ID: "t1-ep1", ServiceID: testServiceID, Host: testHost1, TenantID: testTenantID, Status: StatusHealthy})
	r.Register(&Info{ID: "t1-ep2", ServiceID: testServiceID, Host: testHost2, TenantID: testTenantID, Status: StatusHealthy})
	r.Register(&Info{ID: "tent2-ep1", ServiceID: testServiceID, Host: testHost3, TenantID: testTenantID2, Status: StatusHealthy})

	// Query scoped to tenant1
	results := r.Query(&Query{TenantID: testTenantID})
	assertions.Len(results, 2, "Tenant1 should see only its own endpoints")

	// Query scoped to tenant2
	results = r.Query(&Query{TenantID: testTenantID2})
	assertions.Len(results, 1, "Tenant2 should see only its own endpoints")

	// Unscoped query sees all
	results = r.Query(&Query{})
	assertions.Len(results, 3, "Unscoped query should see all endpoints")
}

func TestTenantIsolationWithTenantSubjectsUniqueness(t *testing.T) {
	var testCaseMsg = "Tenant subjects must differ"
	assertions := assert.New(t)

	s1 := TenantSubjects(testTenantID)
	s2 := TenantSubjects(testTenantID2)

	assertions.NotEqual(s1.Register, s2.Register, testCaseMsg)
	assertions.NotEqual(s1.Deregister, s2.Deregister, testCaseMsg)
	assertions.NotEqual(s1.Query, s2.Query, testCaseMsg)
	assertions.NotEqual(s1.Discover, s2.Discover, testCaseMsg)
	assertions.NotEqual(s1.Events, s2.Events, testCaseMsg)
	assertions.NotEqual(s1.HealthUpdate, s2.HealthUpdate, testCaseMsg)
}

func TestTenantIsolationWithDeregisterDoesNotAffectOtherTenants(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	r.Register(&Info{ID: "te1-ep1", ServiceID: testServiceID, Host: testHost1, TenantID: testTenantID})
	r.Register(&Info{ID: "te2-ep1", ServiceID: testServiceID, Host: testHost2, TenantID: testTenantID2})

	// Deregister tenant1's endpoint
	err := r.Deregister("te1-ep1")
	assertions.NoError(err)

	// Tenant2's endpoint must still exist
	ep, found := r.Get("te2-ep1")
	assertions.True(found, "Tenant2 endpoint should still exist")
	assertions.Equal(testTenantID2, ep.TenantID)
}

func TestTenantIsolationiWithManagerDiscoverByTenant(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: "ten1-ep1", ServiceID: testServiceID, Host: testHost1, TenantID: testTenantID, Status: StatusHealthy})
	manager.RegisterEndpoint(ctx, &Info{ID: "ten2-ep1", ServiceID: testServiceID, Host: testHost2, TenantID: testTenantID2, Status: StatusHealthy})

	// Query filtered by tenant
	results := manager.Query(&Query{TenantID: testTenantID, OnlyAvailable: true})
	assertions.Len(results, 1)
	assertions.Equal(testTenantID, results[0].TenantID)
}

/* ========================================================================================================
   HELPER FUNCTION TESTS
   ======================================================================================================== */

func TestItoa(t *testing.T) {
	testCases := []struct {
		name     string
		input    int
		expected string
	}{
		{"zero", 0, "0"},
		{"positive single", 1, "1"},
		{"positive multi", 123, "123"},
		{"negative single", -1, "-1"},
		{"negative multi", -123, "-123"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.expected, itoa(tc.input))
		})
	}
}

func TestContainsString(t *testing.T) {
	assertions := assert.New(t)
	slice := []string{"a", "b", "c"}

	assertions.True(containsString(slice, "b"))
	assertions.False(containsString(slice, "d"))
}

func TestContainsStatus(t *testing.T) {
	assertions := assert.New(t)
	slice := []Status{StatusHealthy, StatusDegraded}

	assertions.True(containsStatus(slice, StatusHealthy))
	assertions.False(containsStatus(slice, StatusUnhealthy))
}
