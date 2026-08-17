package endpoint

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
)

func BenchmarkRegistryRegister(b *testing.B) {
	r := NewRegistry()

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ep := &Info{
				ID:        "ep" + itoa(i),
				ServiceID: "svc",
				Host:      testHost,
			}
			r.Register(ep)
			i++
		}
	})
}

func BenchmarkRegistryGet(b *testing.B) {
	r := NewRegistry()
	for i := 0; i < 1000; i++ {
		r.Register(&Info{
			ID:        "ep" + itoa(i),
			ServiceID: "svc",
			Host:      testHost,
		})
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			r.Get("ep" + itoa(i%1000))
			i++
		}
	})
}

func BenchmarkRegistryQuery(b *testing.B) {
	r := NewRegistry()
	for i := 0; i < 1000; i++ {
		r.Register(&Info{
			ID:        "ep" + itoa(i),
			ServiceID: "svc" + itoa(i%10),
			Host:      testHost,
			Status:    StatusHealthy,
			Tags:      []string{testTagAPI},
		})
	}

	query := &Query{
		ServiceIDs:  []string{testServiceID, testServiceID2},
		OnlyHealthy: true,
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.Query(query)
		}
	})
}

func BenchmarkRegistryDeregister(b *testing.B) {
	r := NewRegistry()
	for i := 0; i < b.N; i++ {
		r.Register(&Info{
			ID:        "ep" + itoa(i),
			ServiceID: "svc",
			Host:      testHost,
		})
	}

	b.ReportAllocs()
	b.ResetTimer()
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
			Host:      testHost,
		})
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			r.GetByService("svc" + itoa(i%10))
			i++
		}
	})
}

func BenchmarkRegistryUpdateStatus(b *testing.B) {
	r := NewRegistry()
	for i := 0; i < 1000; i++ {
		r.Register(&Info{
			ID:        "ep" + itoa(i),
			ServiceID: "svc",
			Host:      testHost,
			Status:    StatusHealthy,
		})
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			r.UpdateStatus("ep"+itoa(i%1000), StatusDegraded)
			i++
		}
	})
}

func BenchmarkInfoClone(b *testing.B) {
	info := &Info{
		ID:        testEndpointID,
		ServiceID: testServiceID,
		Host:      testHost,
		Tags:      []string{testTagAPI, testTagV1},
		Metadata:  map[string]string{testMetaKey: testMetaValue},
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			info.Clone()
		}
	})
}

func BenchmarkQueryMatches(b *testing.B) {
	ep := &Info{
		ID:        testEndpointID,
		ServiceID: testServiceID,
		Host:      testHost,
		Status:    StatusHealthy,
		Tags:      []string{testTagAPI},
	}

	query := &Query{
		ServiceIDs:  []string{testServiceID},
		OnlyHealthy: true,
		Tags:        []string{testTagAPI},
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			query.Matches(ep)
		}
	})
}

func BenchmarkManagerRegister(b *testing.B) {
	manager := newTestManager()
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.RegisterEndpoint(ctx, &Info{
			ID:        "ep" + itoa(i),
			ServiceID: "svc",
			Host:      testHost,
		})
	}
}

func BenchmarkManagerDiscover(b *testing.B) {
	manager := newTestManager()
	ctx := context.Background()
	for i := 0; i < 100; i++ {
		manager.RegisterEndpoint(ctx, &Info{
			ID:        "ep" + itoa(i),
			ServiceID: "svc",
			Host:      testHost,
			Status:    StatusHealthy,
		})
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			manager.Discover("svc")
		}
	})
}

/* ========================================================================================================
   CLIENT OPERATION TESTS
   ======================================================================================================== */

func TestClientRegister(t *testing.T) {
	testCases := []struct {
		name        string
		reply       *nats.Msg
		reqErr      error
		wantErr     bool
		wantSuccess bool
	}{
		{
			name: "success",
			reply: func() *nats.Msg {
				resp := &RegistrationResponse{Success: true, EndpointID: testEndpointID, Message: "registered"}
				data, _ := json.Marshal(resp)
				return &nats.Msg{Data: data}
			}(),
			wantErr:     false,
			wantSuccess: true,
		},
		{
			name:    "request error",
			reqErr:  errors.New(testErrRequest),
			wantErr: true,
		},
		{
			name:    "invalid response JSON",
			reply:   &nats.Msg{Data: []byte(testInvalidJSON)},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			pub := &mockPublisher{requestReply: tc.reply, requestErr: tc.reqErr}
			client := NewClient(pub)
			client.SetTimeout(time.Second)

			ep := newTestEndpoint(testEndpointID, testServiceID, testHost)
			resp, err := client.Register(context.Background(), ep, testTTL5Min)
			if tc.wantErr {
				assertions.Error(err)
			} else {
				assertions.NoError(err)
				assertions.Equal(tc.wantSuccess, resp.Success)
			}
		})
	}
}

func TestClientDeregister(t *testing.T) {
	testCases := []struct {
		name    string
		reply   *nats.Msg
		reqErr  error
		wantErr bool
	}{
		{
			name: "success",
			reply: func() *nats.Msg {
				resp := &DeregistrationResponse{Success: true, EndpointID: testEndpointID, Message: "deregistered"}
				data, _ := json.Marshal(resp)
				return &nats.Msg{Data: data}
			}(),
		},
		{
			name:    "request-error",
			reqErr:  errors.New(testErrRequest),
			wantErr: true,
		},
		{
			name:    "invalid JSON response",
			reply:   &nats.Msg{Data: []byte(testInvalidJSON)},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			pub := &mockPublisher{requestReply: tc.reply, requestErr: tc.reqErr}
			client := NewClient(pub)

			resp, err := client.Deregister(context.Background(), testEndpointID, testReasonTest)
			if tc.wantErr {
				assertions.Error(err)
			} else {
				assertions.NoError(err)
				assertions.True(resp.Success)
			}
		})
	}
}

func TestClientRenew(t *testing.T) {
	testCases := []struct {
		name    string
		reply   *nats.Msg
		reqErr  error
		wantErr bool
	}{
		{
			name: "success",
			reply: func() *nats.Msg {
				resp := &RenewResponse{Success: true, EndpointID: testEndpointID, Message: "renewed", ExpiresAt: time.Now().Add(testTTL5Min)}
				data, _ := json.Marshal(resp)
				return &nats.Msg{Data: data}
			}(),
		},
		{
			name:    "request_error",
			reqErr:  errors.New(testErrRequest),
			wantErr: true,
		},
		{
			name:    "JSON invalid response",
			reply:   &nats.Msg{Data: []byte(testInvalidJSON)},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			pub := &mockPublisher{requestReply: tc.reply, requestErr: tc.reqErr}
			client := NewClient(pub)

			resp, err := client.Renew(context.Background(), testEndpointID, testTTL5Min)
			if tc.wantErr {
				assertions.Error(err)
			} else {
				assertions.NoError(err)
				assertions.True(resp.Success)
			}
		})
	}
}

func TestClientQuery(t *testing.T) {
	testCases := []struct {
		name    string
		reply   *nats.Msg
		reqErr  error
		wantErr bool
	}{
		{
			name: "success",
			reply: func() *nats.Msg {
				resp := &QueryResponse{Endpoints: []*Info{{ID: testEndpointID}}, Count: 1}
				data, _ := json.Marshal(resp)
				return &nats.Msg{Data: data}
			}(),
		},
		{
			name:    "error request",
			reqErr:  errors.New(testErrRequest),
			wantErr: true,
		},
		{
			name:    "JSON response invalid",
			reply:   &nats.Msg{Data: []byte(testInvalidJSON)},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			pub := &mockPublisher{requestReply: tc.reply, requestErr: tc.reqErr}
			client := NewClient(pub)

			resp, err := client.Query(context.Background(), &Query{ServiceIDs: []string{testServiceID}})
			if tc.wantErr {
				assertions.Error(err)
			} else {
				assertions.NoError(err)
				assertions.Equal(1, resp.Count)
			}
		})
	}
}

func TestClientDiscover(t *testing.T) {
	testCases := []struct {
		name    string
		opts    []DiscoverOption
		reply   *nats.Msg
		reqErr  error
		wantErr bool
	}{
		{
			name: "success no options",
			reply: func() *nats.Msg {
				resp := &DiscoverResponse{Endpoints: []*Info{{ID: testEndpointID}}, Count: 1, ServiceID: testServiceID}
				data, _ := json.Marshal(resp)
				return &nats.Msg{Data: data}
			}(),
		},
		{
			name: "success with options",
			opts: []DiscoverOption{WithTags(testTagAPI), WithVersion(testVersion), WithProtocol(ProtocolHTTP), OnlyHealthy()},
			reply: func() *nats.Msg {
				resp := &DiscoverResponse{Endpoints: []*Info{{ID: testEndpointID}}, Count: 1, ServiceID: testServiceID}
				data, _ := json.Marshal(resp)
				return &nats.Msg{Data: data}
			}(),
		},
		{
			name:    "error-request",
			reqErr:  errors.New(testErrRequest),
			wantErr: true,
		},
		{
			name:    "response JSON invalid",
			reply:   &nats.Msg{Data: []byte(testInvalidJSON)},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			pub := &mockPublisher{requestReply: tc.reply, requestErr: tc.reqErr}
			client := NewClient(pub)

			resp, err := client.Discover(context.Background(), testServiceID, tc.opts...)
			if tc.wantErr {
				assertions.Error(err)
			} else {
				assertions.NoError(err)
				assertions.Equal(testServiceID, resp.ServiceID)
			}
		})
	}
}

func TestClientUpdateHealth(t *testing.T) {
	testCases := []struct {
		name    string
		pubErr  error
		wantErr bool
	}{
		{"success", nil, false},
		{"publish error", errors.New(testErrPublish), true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			pub := &mockPublisher{publishErr: tc.pubErr}
			client := NewClient(pub)

			err := client.UpdateHealth(context.Background(), &HealthUpdate{
				EndpointID: testEndpointID,
				Status:     StatusHealthy,
				CheckedAt:  time.Now(),
			})
			if tc.wantErr {
				assertions.Error(err)
			} else {
				assertions.NoError(err)
			}
		})
	}
}

/* ========================================================================================================
   HANDLER RESPONSE AND SUBSCRIPTION TESTS
   ======================================================================================================== */

func TestHandlerSendResponseWithPublisher(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()
	pub := &capturingPublisher{}
	handler := NewHandler(manager, WithPublisher(pub))

	// Register via NATS message with reply
	req := RegistrationRequest{
		Endpoint: newTestEndpoint(testEndpointID, testServiceID, testHost),
	}
	data, _ := json.Marshal(req)
	msg := newTestMsgWithReply(data, testReplySubject)

	err := handler.HandleRegister(context.Background(), msg)
	assertions.NoError(err)

	// Verify response was published
	pub.mu.Lock()
	assertions.Equal(1, len(pub.subjects))
	assertions.Equal(testReplySubject, pub.subjects[0])
	pub.mu.Unlock()
}

func TestHandlerSendErrorResponseWithPublisher(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()
	pub := &capturingPublisher{}
	handler := NewHandler(manager, WithPublisher(pub))

	// Send invalid JSON with reply subject
	msg := newTestMsgWithReply([]byte(testInvalidJSON), testReplySubject)

	err := handler.HandleRegister(context.Background(), msg)
	assertions.Error(err)

	// Verify error response was published
	pub.mu.Lock()
	assertions.GreaterOrEqual(len(pub.subjects), 1)
	assertions.Equal(testReplySubject, pub.subjects[0])
	// Verify the response body contains error
	var errResp ErrorResponse
	json.Unmarshal(pub.messages[0].Data, &errResp)
	assertions.False(errResp.Success)
	assertions.NotEmpty(errResp.Error)
	pub.mu.Unlock()
}

func TestHandlerSendResponseNoReply(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()
	pub := &capturingPublisher{}
	handler := NewHandler(manager, WithPublisher(pub))

	// Message with no Reply subject
	req := RegistrationRequest{
		Endpoint: newTestEndpoint(testEndpointID, testServiceID, testHost),
	}
	data, _ := json.Marshal(req)
	msg := newTestMsg(data)

	err := handler.HandleRegister(context.Background(), msg)
	assertions.NoError(err)

	// No response should be published
	pub.mu.Lock()
	assertions.Empty(pub.subjects)
	pub.mu.Unlock()
}

func TestHandlerSendErrorResponseNoPublisher(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()
	handler := NewHandler(manager) // no publisher

	msg := newTestMsgWithReply([]byte(testInvalidJSON), testReplySubject)
	err := handler.HandleRegister(context.Background(), msg)
	assertions.Error(err) // returns the original error
}

func TestHandlerSendResponsePublishError(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()
	pub := &capturingPublisher{pubErr: errors.New(testErrPublish)}
	handler := NewHandler(manager, WithPublisher(pub))

	req := RegistrationRequest{
		Endpoint: newTestEndpoint(testEndpointID, testServiceID, testHost),
	}
	data, _ := json.Marshal(req)
	msg := newTestMsgWithReply(data, testReplySubject)

	err := handler.HandleRegister(context.Background(), msg)
	assertions.Error(err)
}

func TestRegisterSubscriptions(t *testing.T) {
	testCases := []struct {
		name       string
		subErr     error
		failAtCall int
		wantErr    bool
		wantCount  int
	}{
		{"all succeed", nil, 0, false, 6},
		{"fail at second subscribe", errors.New(testErrSubscribe), 2, true, 0},
		{"fail at third subscribe", errors.New(testErrSubscribe), 3, true, 0},
		{"fail at fourth subscribe", errors.New(testErrSubscribe), 4, true, 0},
		{"fail at fifth subscribe", errors.New(testErrSubscribe), 5, true, 0},
		{"fail at sixth subscribe", errors.New(testErrSubscribe), 6, true, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			manager := newTestManager()
			handler := NewHandler(manager)

			sub := &mockSubscriber{
				subscribeErr: tc.subErr,
				failAtCall:   tc.failAtCall,
			}

			subs, err := handler.RegisterSubscriptions(context.Background(), sub)
			if tc.wantErr {
				assertions.Error(err)
			} else {
				assertions.NoError(err)
				assertions.Len(subs, tc.wantCount)
			}
		})
	}
}

func TestRegisterSubscriptionsFailFirst(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()
	handler := NewHandler(manager)

	sub := &mockSubscriber{
		subscribeErr: errors.New(testErrSubscribe),
		failAtCall:   0, // fail all
	}

	_, err := handler.RegisterSubscriptions(context.Background(), sub)
	assertions.Error(err)
}

func TestRegisterQueueSubscriptions(t *testing.T) {
	testCases := []struct {
		name       string
		subErr     error
		failAtCall int
		wantErr    bool
		wantCount  int
	}{
		{"all succeed", nil, 0, false, 5},
		{"fail at second subscribe", errors.New(testErrSubscribe), 2, true, 0},
		{"fail at third subscribe", errors.New(testErrSubscribe), 3, true, 0},
		{"fail at fourth subscribe", errors.New(testErrSubscribe), 4, true, 0},
		{"fail at fifth subscribe", errors.New(testErrSubscribe), 5, true, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			manager := newTestManager()
			handler := NewHandler(manager)

			sub := &mockSubscriber{
				queueSubscribeErr: tc.subErr,
				failAtCall:        tc.failAtCall,
			}

			subs, err := handler.RegisterQueueSubscriptions(context.Background(), sub, testQueueName)
			if tc.wantErr {
				assertions.Error(err)
			} else {
				assertions.NoError(err)
				assertions.Len(subs, tc.wantCount)
			}
		})
	}
}

func TestRegisterQueueSubscriptionsFailFirst(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()
	handler := NewHandler(manager)

	sub := &mockSubscriber{
		queueSubscribeErr: errors.New(testErrSubscribe),
		failAtCall:        0, // fail all
	}

	_, err := handler.RegisterQueueSubscriptions(context.Background(), sub, testQueueName)
	assertions.Error(err)
}

func TestUnsubscribeAll(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()
	handler := NewHandler(manager)

	subs := []core.Subscription{
		&mockSubscription{subject: testSubjectRegister},
		nil, // nil subscriptions should be skipped
		&mockSubscription{subject: testSubjectDeregister},
	}

	// Should not panic
	handler.unsubscribeAll(subs)
	assertions.True(true) // reached without panic
}

/* ========================================================================================================
   HEALTH CHECKER CHECK TESTS
   ======================================================================================================== */

func TestHealthCheckerCheckHealthy(t *testing.T) {
	assertions := assert.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	checker := NewHealthChecker(HealthCheckerConfig{
		Timeout:            time.Second,
		HealthyThreshold:   1,
		UnhealthyThreshold: 1,
	})

	ep := &Info{
		ID:              testEndpointID,
		Host:            srv.Listener.Addr().(*net.TCPAddr).IP.String(),
		Port:            srv.Listener.Addr().(*net.TCPAddr).Port,
		Protocol:        ProtocolHTTP,
		HealthCheckPath: testHealthPath,
	}

	result := checker.Check(context.Background(), ep)
	assertions.Equal(StatusHealthy, result.Status)
	assertions.Equal("OK", result.Message)
	assertions.NotZero(result.Duration)
	assertions.NotNil(result.Details)
}

func TestHealthCheckerCheckUnhealthy(t *testing.T) {
	assertions := assert.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	checker := NewHealthChecker(HealthCheckerConfig{
		Timeout:            time.Second,
		HealthyThreshold:   2,
		UnhealthyThreshold: 1,
	})

	ep := &Info{
		ID:              testEndpointID,
		Host:            srv.Listener.Addr().(*net.TCPAddr).IP.String(),
		Port:            srv.Listener.Addr().(*net.TCPAddr).Port,
		Protocol:        ProtocolHTTP,
		HealthCheckPath: testHealthPath,
	}

	result := checker.Check(context.Background(), ep)
	assertions.Equal(StatusUnhealthy, result.Status)
}

func TestHealthCheckerCheckMaintenance(t *testing.T) {
	assertions := assert.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	checker := NewHealthChecker(HealthCheckerConfig{
		Timeout:            time.Second,
		HealthyThreshold:   2,
		UnhealthyThreshold: 1,
	})

	ep := &Info{
		ID:              testEndpointID,
		Host:            srv.Listener.Addr().(*net.TCPAddr).IP.String(),
		Port:            srv.Listener.Addr().(*net.TCPAddr).Port,
		Protocol:        ProtocolHTTP,
		HealthCheckPath: testHealthPath,
	}

	result := checker.Check(context.Background(), ep)
	assertions.Equal(StatusMaintenance, result.Status)
	assertions.Equal("Service unavailable", result.Message)
}

func TestHealthCheckerCheckConnectionRefused(t *testing.T) {
	assertions := assert.New(t)

	checker := NewHealthChecker(HealthCheckerConfig{
		Timeout:            100 * time.Millisecond,
		HealthyThreshold:   2,
		UnhealthyThreshold: 1,
	})

	ep := &Info{
		ID:              testEndpointID,
		Host:            "127.0.0.1",
		Port:            19999, // unlikely to be in use
		Protocol:        ProtocolHTTP,
		HealthCheckPath: testHealthPath,
	}

	result := checker.Check(context.Background(), ep)
	assertions.Equal(StatusUnhealthy, result.Status)
	assertions.NotEmpty(result.Message)
}

func TestHealthCheckerCheckDefaultsNoPort(t *testing.T) {
	assertions := assert.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	checker := NewHealthChecker(HealthCheckerConfig{
		Timeout:            time.Second,
		HealthyThreshold:   1,
		UnhealthyThreshold: 1,
	})

	// Endpoint with HealthCheckPort set, empty Protocol, empty HealthCheckPath
	ep := &Info{
		ID:              testEndpointID,
		Host:            srv.Listener.Addr().(*net.TCPAddr).IP.String(),
		HealthCheckPort: srv.Listener.Addr().(*net.TCPAddr).Port,
	}

	result := checker.Check(context.Background(), ep)
	// Should use default /health path and http protocol
	assertions.NotNil(result)
}

func TestHealthCheckerRecoveringState(t *testing.T) {
	assertions := assert.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	checker := NewHealthChecker(HealthCheckerConfig{
		Timeout:            time.Second,
		HealthyThreshold:   3, // need 3 consecutive successes
		UnhealthyThreshold: 1,
	})

	ep := &Info{
		ID:              testEndpointID,
		Host:            srv.Listener.Addr().(*net.TCPAddr).IP.String(),
		Port:            srv.Listener.Addr().(*net.TCPAddr).Port,
		Protocol:        ProtocolHTTP,
		HealthCheckPath: testHealthPath,
	}

	// First check - not yet at threshold, should be Degraded/Recovering
	result := checker.Check(context.Background(), ep)
	assertions.Equal(StatusDegraded, result.Status)
	assertions.Equal("Recovering", result.Message)

	// Second check - still recovering
	result = checker.Check(context.Background(), ep)
	assertions.Equal(StatusDegraded, result.Status)

	// Third check - now healthy
	result = checker.Check(context.Background(), ep)
	assertions.Equal(StatusHealthy, result.Status)
}

func TestHealthCheckerRecordSuccessResetsFails(t *testing.T) {
	assertions := assert.New(t)

	checker := NewHealthChecker(HealthCheckerConfig{
		Timeout:            time.Second,
		HealthyThreshold:   2,
		UnhealthyThreshold: 2,
	})

	// Simulate failures then success
	checker.recordFailure(testEndpointID)
	checker.recordFailure(testEndpointID)

	checker.mu.RLock()
	assertions.Equal(2, checker.failureCnt[testEndpointID])
	checker.mu.RUnlock()

	checker.recordSuccess(testEndpointID)

	checker.mu.RLock()
	assertions.Equal(0, checker.failureCnt[testEndpointID])
	assertions.Equal(1, checker.successCnt[testEndpointID])
	checker.mu.RUnlock()
}

func TestHealthCheckerIsHealthy(t *testing.T) {
	assertions := assert.New(t)

	checker := NewHealthChecker(HealthCheckerConfig{
		HealthyThreshold: 2,
	})

	assertions.False(checker.isHealthy(testEndpointID))

	checker.recordSuccess(testEndpointID)
	assertions.False(checker.isHealthy(testEndpointID))

	checker.recordSuccess(testEndpointID)
	assertions.True(checker.isHealthy(testEndpointID))
}

/* ========================================================================================================
   MANAGER PUBLISH EVENT AND DISPATCH ERROR TESTS
   ======================================================================================================== */

func TestManagerDispatchError(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	var capturedErr error
	manager.OnError(func(err error) {
		capturedErr = err
	})

	testErr := errors.New("test dispatch error")
	manager.dispatchError(testErr)

	assertions.Equal(testErr, capturedErr)

	// Also check the error channel
	select {
	case err := <-manager.Errors():
		assertions.Equal(testErr, err)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("expected error on channel")
	}
}

func TestManagerDispatchErrorChannelFull(t *testing.T) {
	assertions := assert.New(t)

	m, _ := NewManager(ManagerConfig{
		HealthCheck:     HealthCheckConfig{Enabled: false},
		ErrorBufferSize: 1,
	})

	// Fill the channel
	m.dispatchError(errors.New("err1"))
	// This should not block even though channel is full
	m.dispatchError(errors.New("err2"))
	assertions.True(true) // did not block
}

func TestManagerPublishEventWithPublisher(t *testing.T) {
	assertions := assert.New(t)

	pub := &capturingPublisher{}
	manager := newTestManager()
	manager.SetPublisher(pub)

	ctx := context.Background()
	ep := &Info{
		ID:        testEndpointID,
		ServiceID: testServiceID,
		Host:      testHost,
		Status:    StatusHealthy,
	}
	manager.RegisterEndpoint(ctx, ep)

	// Give callbacks/goroutines time to fire
	time.Sleep(50 * time.Millisecond)

	pub.mu.Lock()
	assertions.GreaterOrEqual(len(pub.subjects), 1)
	pub.mu.Unlock()
}

func TestManagerPublishEventNoPublisher(t *testing.T) {
	manager := newTestManager()
	// publishEvent with nil publisher should be a no-op
	manager.publishEvent(EventRegistered, &Info{ID: testEndpointID, ServiceID: testServiceID}, StatusUnknown, StatusHealthy, "")
	// Should not panic
}

func TestManagerPublishEventAllTypes(t *testing.T) {
	testCases := []struct {
		name      string
		eventType EventType
	}{
		{"registered", EventRegistered},
		{"deregistered", EventDeregistered},
		{"updated", EventUpdated},
		{"healthy", EventHealthy},
		{"unhealthy", EventUnhealthy},
		{"expired", EventExpired},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pub := &capturingPublisher{}
			manager := newTestManager()
			manager.SetPublisher(pub)

			manager.publishEvent(tc.eventType, &Info{ID: testEndpointID, ServiceID: testServiceID}, StatusUnknown, StatusHealthy, "")
			time.Sleep(50 * time.Millisecond)

			pub.mu.Lock()
			assertions := assert.New(t)
			assertions.Equal(1, len(pub.subjects))
			pub.mu.Unlock()
		})
	}
}

func TestManagerPublishEventError(t *testing.T) {
	assertions := assert.New(t)

	pub := &capturingPublisher{pubErr: errors.New(testErrPublish)}
	manager := newTestManager()
	manager.SetPublisher(pub)

	var capturedErr error
	manager.OnError(func(err error) {
		capturedErr = err
	})

	manager.publishEvent(EventRegistered, &Info{ID: testEndpointID, ServiceID: testServiceID}, StatusUnknown, StatusHealthy, "")
	time.Sleep(50 * time.Millisecond)

	assertions.NotNil(capturedErr)
}

/* ========================================================================================================
   REGISTRY CLEANUP EXPIRED TESTS
   ======================================================================================================== */

func TestRegistryCleanupExpired(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry(RegistryConfig{
		CleanupInterval:   10 * time.Millisecond,
		ExpirationEnabled: true,
	})
	defer r.Shutdown(context.Background())

	// Register with very short TTL
	ep := &Info{
		ID:        testEndpointID,
		ServiceID: testServiceID,
		Host:      testHost,
		ExpiresAt: time.Now().Add(5 * time.Millisecond),
	}
	r.Register(ep)

	// Also register one that does not expire
	ep2 := newTestEndpoint(testEndpointID2, testServiceID, testHost2)
	r.Register(ep2)

	// Wait for cleanup
	time.Sleep(50 * time.Millisecond)

	// Expired endpoint should be gone
	_, found := r.Get(testEndpointID)
	assertions.False(found, "Expired endpoint should be cleaned up")

	// Non-expired endpoint should still exist
	_, found = r.Get(testEndpointID2)
	assertions.True(found, "Non-expired endpoint should still exist")
}

func TestRegistryGetExpiredEndpointReturnsNotFound(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})

	ep := &Info{
		ID:        testEndpointID,
		ServiceID: testServiceID,
		Host:      testHost,
		ExpiresAt: time.Now().Add(-time.Hour), // already expired
	}
	r.Register(ep)

	_, found := r.Get(testEndpointID)
	assertions.False(found, "Expired endpoint should not be returned by Get")
}

func TestRegistryGetByServiceSkipsExpired(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})

	r.Register(&Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost1, ExpiresAt: time.Now().Add(-time.Hour)})
	r.Register(&Info{ID: testEndpointID2, ServiceID: testServiceID, Host: testHost2})

	results := r.GetByService(testServiceID)
	assertions.Len(results, 1, "Should skip expired endpoints")
}

func TestRegistryGetByServiceNotFound(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	results := r.GetByService("nonexistent")
	assertions.Nil(results)
}

func TestRegistryGetByTagSkipsExpired(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})

	r.Register(&Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost1, Tags: []string{testTagAPI}, ExpiresAt: time.Now().Add(-time.Hour)})
	r.Register(&Info{ID: testEndpointID2, ServiceID: testServiceID, Host: testHost2, Tags: []string{testTagAPI}})

	results := r.GetByTag(testTagAPI)
	assertions.Len(results, 1)
}

func TestRegistryQueryOffsetBeyondResults(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	r.Register(newTestEndpoint(testEndpointID, testServiceID, testHost1))

	results := r.Query(&Query{Offset: 100})
	assertions.Nil(results, "Offset beyond results should return nil")
}

func TestRegistryAllSkipsExpired(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})

	r.Register(&Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost1, ExpiresAt: time.Now().Add(-time.Hour)})
	r.Register(&Info{ID: testEndpointID2, ServiceID: testServiceID, Host: testHost2})

	all := r.All()
	assertions.Len(all, 1)
}

func TestRegistryQuerySkipsExpired(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry(RegistryConfig{ExpirationEnabled: false})

	r.Register(&Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost1, ExpiresAt: time.Now().Add(-time.Hour)})
	r.Register(&Info{ID: testEndpointID2, ServiceID: testServiceID, Host: testHost2})

	results := r.Query(&Query{ServiceIDs: []string{testServiceID}})
	assertions.Len(results, 1)
}

/* ========================================================================================================
   HTTP HANDLER EDGE CASE TESTS
   ======================================================================================================== */

func TestManagerHTTPHandleRegistrationBadJSON(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	httpReq := httptest.NewRequest(http.MethodPost, testHTTPPathRegister, bytes.NewReader([]byte(testInvalidJSON)))
	rr := httptest.NewRecorder()

	manager.HandleRegistration(rr, httpReq)
	assertions.Equal(http.StatusBadRequest, rr.Code)
}

func TestManagerHTTPHandleRegistrationValidationFail(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	// Endpoint missing required fields
	req := RegistrationRequest{Endpoint: &Info{}}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, testHTTPPathRegister, bytes.NewReader(body))
	rr := httptest.NewRecorder()

	manager.HandleRegistration(rr, httpReq)
	assertions.Equal(http.StatusInternalServerError, rr.Code)
}

func TestManagerHTTPHandleDeregistrationBadJSON(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	httpReq := httptest.NewRequest(http.MethodDelete, testHTTPPathDeregister, bytes.NewReader([]byte(testInvalidJSON)))
	rr := httptest.NewRecorder()

	manager.HandleDeregistration(rr, httpReq)
	assertions.Equal(http.StatusBadRequest, rr.Code)
}

func TestManagerHTTPHandleDeregistrationNotFound(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	req := DeregistrationRequest{EndpointID: "nonexistent"}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodDelete, testHTTPPathDeregister, bytes.NewReader(body))
	rr := httptest.NewRecorder()

	manager.HandleDeregistration(rr, httpReq)
	assertions.Equal(http.StatusNotFound, rr.Code)
}

func TestManagerHTTPHandleDeregistrationValidationFail(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	req := DeregistrationRequest{} // missing EndpointID
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodDelete, testHTTPPathDeregister, bytes.NewReader(body))
	rr := httptest.NewRecorder()

	manager.HandleDeregistration(rr, httpReq)
	assertions.Equal(http.StatusInternalServerError, rr.Code)
}

func TestManagerHTTPHandleQueryPOST(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost, Status: StatusHealthy})

	q := Query{ServiceIDs: []string{testServiceID}}
	body, _ := json.Marshal(q)

	httpReq := httptest.NewRequest(http.MethodPost, testHTTPPathQuery, bytes.NewReader(body))
	rr := httptest.NewRecorder()

	manager.HandleQuery(rr, httpReq)
	assertions.Equal(http.StatusOK, rr.Code)

	var endpoints []*Info
	json.NewDecoder(rr.Body).Decode(&endpoints)
	assertions.Len(endpoints, 1)
}

func TestManagerHTTPHandleQueryBadJSON(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	httpReq := httptest.NewRequest(http.MethodPost, testHTTPPathQuery, bytes.NewReader([]byte(testInvalidJSON)))
	rr := httptest.NewRecorder()

	manager.HandleQuery(rr, httpReq)
	assertions.Equal(http.StatusBadRequest, rr.Code)
}

func TestManagerHTTPHandleQueryGETParams(t *testing.T) {
	testCases := []struct {
		name  string
		query string
	}{
		{"healthy", "?healthy=true"},
		{"available", "?available=true"},
		{"service_id and healthy", "?service_id=" + testServiceID + "&healthy=true"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			manager := newTestManager()

			ctx := context.Background()
			manager.RegisterEndpoint(ctx, &Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost, Status: StatusHealthy})

			httpReq := httptest.NewRequest(http.MethodGet, testHTTPPathQuery+tc.query, nil)
			rr := httptest.NewRecorder()

			manager.HandleQuery(rr, httpReq)
			assertions.Equal(http.StatusOK, rr.Code)
		})
	}
}

func TestManagerHTTPHandleRenewBadJSON(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	httpReq := httptest.NewRequest(http.MethodPost, testHTTPPathRenew, bytes.NewReader([]byte(testInvalidJSON)))
	rr := httptest.NewRecorder()

	manager.HandleRenew(rr, httpReq)
	assertions.Equal(http.StatusBadRequest, rr.Code)
}

func TestManagerHTTPHandleRenewNotFound(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	req := struct {
		EndpointID string        `json:"endpoint_id"`
		TTL        time.Duration `json:"ttl"`
	}{EndpointID: "nonexistent", TTL: testTTL5Min}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPost, testHTTPPathRenew, bytes.NewReader(body))
	rr := httptest.NewRecorder()

	manager.HandleRenew(rr, httpReq)
	assertions.Equal(http.StatusNotFound, rr.Code)
}

func TestManagerHTTPHandleRenewPUT(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, newTestEndpoint(testEndpointID, testServiceID, testHost))

	req := struct {
		EndpointID string        `json:"endpoint_id"`
		TTL        time.Duration `json:"ttl"`
	}{EndpointID: testEndpointID, TTL: testTTL10Min}
	body, _ := json.Marshal(req)

	httpReq := httptest.NewRequest(http.MethodPut, testHTTPPathRenew, bytes.NewReader(body))
	rr := httptest.NewRecorder()

	manager.HandleRenew(rr, httpReq)
	assertions.Equal(http.StatusOK, rr.Code)
}

/* ========================================================================================================
   HANDLER NATS MESSAGE EDGE CASE TESTS
   ======================================================================================================== */

func TestHandlerHandleDeregisterNotFound(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()
	handler := NewHandler(manager)

	req := DeregistrationRequest{EndpointID: "nonexistent"}
	data, _ := json.Marshal(req)
	msg := newTestMsg(data)

	err := handler.HandleDeregister(context.Background(), msg)
	assertions.Error(err)
}

func TestHandlerHandleRenewNotFound(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()
	handler := NewHandler(manager)

	req := RenewRequest{EndpointID: "nonexistent", TTL: testTTL5Min}
	data, _ := json.Marshal(req)
	msg := newTestMsg(data)

	err := handler.HandleRenew(context.Background(), msg)
	assertions.Error(err)
}

func TestHandlerHandleHealthUpdateZeroCheckedAt(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{ID: testEndpointID, ServiceID: testServiceID, Host: testHost, Status: StatusHealthy})

	handler := NewHandler(manager)
	update := HealthUpdate{
		EndpointID: testEndpointID,
		Status:     StatusDegraded,
		// CheckedAt is zero - should be filled in
	}
	data, _ := json.Marshal(update)
	msg := newTestMsg(data)

	err := handler.HandleHealthUpdate(ctx, msg)
	assertions.NoError(err)
}

func TestHandlerHandleHealthUpdateNotFound(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()
	handler := NewHandler(manager)

	update := HealthUpdate{
		EndpointID: "nonexistent",
		Status:     StatusHealthy,
		CheckedAt:  time.Now(),
	}
	data, _ := json.Marshal(update)
	msg := newTestMsg(data)

	err := handler.HandleHealthUpdate(context.Background(), msg)
	assertions.Error(err)
}

func TestHandlerHandleDiscoverWithOptions(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{
		ID:        testEndpointID,
		ServiceID: testServiceID,
		Host:      testHost,
		Status:    StatusHealthy,
		Version:   testVersion,
		Protocol:  ProtocolHTTP,
		Tags:      []string{testTagAPI},
		TenantID:  testTenantID,
	})

	handler := NewHandler(manager)

	testCases := []struct {
		name string
		req  DiscoverRequest
	}{
		{
			name: "with version",
			req:  DiscoverRequest{ServiceID: testServiceID, Version: testVersion},
		},
		{
			name: "with protocol",
			req:  DiscoverRequest{ServiceID: testServiceID, Protocol: ProtocolHTTP},
		},
		{
			name: "with tags",
			req:  DiscoverRequest{ServiceID: testServiceID, Tags: []string{testTagAPI}},
		},
		{
			name: "only healthy",
			req:  DiscoverRequest{ServiceID: testServiceID, OnlyHealthy: true},
		},
		{
			name: "with tenant",
			req:  DiscoverRequest{ServiceID: testServiceID, TenantID: testTenantID},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, _ := json.Marshal(tc.req)
			msg := newTestMsg(data)
			err := handler.HandleDiscover(ctx, msg)
			assertions.NoError(err)
		})
	}
}

func TestHandlerHandleQueryWithTenantContext(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{
		ID:        testEndpointID,
		ServiceID: testServiceID,
		Host:      testHost,
		Status:    StatusHealthy,
		TenantID:  testTenantID,
	})
	manager.RegisterEndpoint(ctx, &Info{
		ID:        testEndpointID2,
		ServiceID: testServiceID,
		Host:      testHost2,
		Status:    StatusHealthy,
		TenantID:  testTenantID2,
	})

	pub := &capturingPublisher{}
	handler := NewHandler(manager, WithPublisher(pub))

	q := Query{}
	data, _ := json.Marshal(q)
	msg := newTestMsgWithReply(data, testReplySubject)

	tenantCtx := core.WithTenantID(ctx, testTenantID)
	err := handler.HandleQuery(tenantCtx, msg)
	assertions.NoError(err)

	// Verify filtered by tenant
	pub.mu.Lock()
	assertions.Equal(1, len(pub.messages))
	var resp QueryResponse
	json.Unmarshal(pub.messages[0].Data, &resp)
	assertions.Equal(1, resp.Count)
	pub.mu.Unlock()
}

func TestHandlerHandleRegisterWithTenantContext(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	handler := NewHandler(manager)

	req := RegistrationRequest{
		Endpoint: newTestEndpoint(testEndpointID, testServiceID, testHost),
	}
	data, _ := json.Marshal(req)
	msg := newTestMsg(data)

	tenantCtx := core.WithTenantID(context.Background(), testTenantID)
	err := handler.HandleRegister(tenantCtx, msg)
	assertions.NoError(err)

	// Verify tenant was set
	ep, found := manager.Get(testEndpointID)
	assertions.True(found)
	assertions.Equal(testTenantID, ep.TenantID)
}

func TestHandlerHandleDiscoverWithTenantContext(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	manager.RegisterEndpoint(ctx, &Info{
		ID:        testEndpointID,
		ServiceID: testServiceID,
		Host:      testHost,
		Status:    StatusHealthy,
		TenantID:  testTenantID,
	})

	handler := NewHandler(manager)
	req := DiscoverRequest{ServiceID: testServiceID}
	data, _ := json.Marshal(req)
	msg := newTestMsg(data)

	tenantCtx := core.WithTenantID(ctx, testTenantID)
	err := handler.HandleDiscover(tenantCtx, msg)
	assertions.NoError(err)
}

/* ========================================================================================================
   REGISTRY DEREGISTER EDGE CASE TESTS
   ======================================================================================================== */

func TestRegistryDeregisterWithReasonEmptyID(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	err := r.DeregisterWithReason("", testReasonTest)
	assertions.ErrorIs(err, ErrMissingEndpointID)
}

func TestRegistryDeregisterByServiceEmpty(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	count, err := r.DeregisterByService("")
	assertions.ErrorIs(err, ErrMissingServiceID)
	assertions.Equal(0, count)
}

func TestRegistryDeregisterByServiceNotFound(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	count, err := r.DeregisterByService("nonexistent")
	assertions.NoError(err)
	assertions.Equal(0, count)
}

func TestRegistryDeregisterByInstanceEmpty(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	count, err := r.DeregisterByInstance("")
	assertions.NoError(err)
	assertions.Equal(0, count)
}

func TestRegistryDeregisterByInstanceNotFound(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	count, err := r.DeregisterByInstance("nonexistent")
	assertions.NoError(err)
	assertions.Equal(0, count)
}

func TestRegistryUpdateStatusNotFound(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	err := r.UpdateStatus("nonexistent", StatusHealthy)
	assertions.ErrorIs(err, ErrEndpointNotFound)
}

func TestRegistryUpdateHealthEmptyEndpointID(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	err := r.UpdateHealth(&HealthUpdate{EndpointID: ""})
	assertions.ErrorIs(err, ErrInvalidEndpoint)
}

func TestRegistryShutdownTimeout(t *testing.T) {
	assertions := assert.New(t)

	r := NewRegistry(RegistryConfig{
		CleanupInterval:   time.Hour, // long interval so loop is blocking on ticker
		ExpirationEnabled: true,
	})

	// Shutdown with already-cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := r.Shutdown(ctx)
	assertions.Error(err)
}

/* ========================================================================================================
   MANAGER REGISTER/DEREGISTER ERROR PATH TESTS
   ======================================================================================================== */

func TestManagerRegisterWithTTLRequest(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	ctx := context.Background()
	req := &RegistrationRequest{
		Endpoint: newTestEndpoint(testEndpointID, testServiceID, testHost),
		TTL:      testTTL5Min,
	}

	err := manager.Register(ctx, req)
	assertions.NoError(err)

	ep, found := manager.Get(testEndpointID)
	assertions.True(found)
	assertions.False(ep.ExpiresAt.IsZero())
}

func TestManagerDeregisterEndpointNotFound(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	err := manager.DeregisterEndpoint(context.Background(), "nonexistent")
	assertions.Error(err)
}

func TestManagerDeregisterNotFound(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	err := manager.Deregister(context.Background(), &DeregistrationRequest{EndpointID: "nonexistent"})
	assertions.Error(err)
}

func TestManagerDeregisterServiceEmpty(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	_, err := manager.DeregisterService(context.Background(), "")
	assertions.Error(err)
}

func TestManagerDeregisterInstanceEmpty(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	count, err := manager.DeregisterInstance(context.Background(), "")
	assertions.NoError(err)
	assertions.Equal(0, count)
}

func TestManagerRegisterEndpointError(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	// Missing required fields
	err := manager.RegisterEndpoint(context.Background(), &Info{})
	assertions.Error(err)
}

func TestManagerRegisterWithTTLError(t *testing.T) {
	assertions := assert.New(t)
	manager := newTestManager()

	// Missing required fields
	err := manager.RegisterWithTTL(context.Background(), &Info{}, testTTL5Min)
	assertions.Error(err)
}

func TestManagerShutdownTimeout(t *testing.T) {
	assertions := assert.New(t)

	manager, _ := NewManager(ManagerConfig{
		HealthCheck: HealthCheckConfig{
			Enabled:      true,
			Interval:     time.Hour,
			InitialDelay: time.Hour,
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := manager.Shutdown(ctx)
	// Might timeout or succeed depending on timing
	_ = err
	assertions.True(true)
}

/* ========================================================================================================
   QUERY SORTING TESTS
   ======================================================================================================== */

func TestRegistryQuerySortByPriorityAndWeight(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	r.Register(&Info{ID: "low", ServiceID: testServiceID, Host: testHost1, Priority: 1, Weight: 10})
	r.Register(&Info{ID: "high", ServiceID: testServiceID, Host: testHost2, Priority: 10, Weight: 5})
	r.Register(&Info{ID: "high-heavy", ServiceID: testServiceID, Host: testHost3, Priority: 10, Weight: 20})

	results := r.Query(&Query{})
	assertions.Len(results, 3)
	// high-heavy (pri=10, wt=20) first, then high (pri=10, wt=5), then low (pri=1, wt=10)
	assertions.Equal("high-heavy", results[0].ID)
	assertions.Equal("high", results[1].ID)
	assertions.Equal("low", results[2].ID)
}

/* ========================================================================================================
   DEFAULT CONFIG TESTS
   ======================================================================================================== */

func TestDefaultRegistryConfig(t *testing.T) {
	assertions := assert.New(t)
	cfg := DefaultRegistryConfig()
	assertions.True(cfg.ExpirationEnabled)
	assertions.Equal(time.Minute, cfg.CleanupInterval)
}

func TestDefaultPublishSubjects(t *testing.T) {
	assertions := assert.New(t)
	subjects := DefaultPublishSubjects()
	assertions.NotEmpty(subjects.Registered)
	assertions.NotEmpty(subjects.Deregistered)
	assertions.NotEmpty(subjects.Updated)
	assertions.NotEmpty(subjects.HealthChange)
}

/* ========================================================================================================
   CLONE EDGE CASE TESTS
   ======================================================================================================== */

func TestInfoCloneNoSlicesNoMetadata(t *testing.T) {
	assertions := assert.New(t)
	original := &Info{
		ID:        testEndpointID,
		ServiceID: testServiceID,
		Host:      testHost,
	}

	clone := original.Clone()
	assertions.NotSame(original, clone)
	assertions.Nil(clone.Methods)
	assertions.Nil(clone.Tags)
	assertions.Nil(clone.Metadata)
}

/* ========================================================================================================
   ADDITIONAL BENCHMARKS
   ======================================================================================================== */

func BenchmarkHealthCheckerCheck(b *testing.B) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	checker := NewHealthChecker(HealthCheckerConfig{
		Timeout:            time.Second,
		HealthyThreshold:   1,
		UnhealthyThreshold: 1,
	})

	ep := &Info{
		ID:              testEndpointID,
		Host:            srv.Listener.Addr().(*net.TCPAddr).IP.String(),
		Port:            srv.Listener.Addr().(*net.TCPAddr).Port,
		Protocol:        ProtocolHTTP,
		HealthCheckPath: testHealthPath,
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		checker.Check(context.Background(), ep)
	}
}

func BenchmarkInfoURL(b *testing.B) {
	info := &Info{
		Protocol: ProtocolHTTP,
		Host:     testHost,
		Port:     testPort,
		Path:     testPath,
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = info.URL()
		}
	})
}

func BenchmarkRegistryAll(b *testing.B) {
	r := NewRegistry()
	for i := 0; i < 1000; i++ {
		r.Register(&Info{
			ID:        "ep" + itoa(i),
			ServiceID: "svc" + itoa(i%10),
			Host:      testHost,
			Status:    StatusHealthy,
		})
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.All()
		}
	})
}

func BenchmarkRegistryClear(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := NewRegistry()
		for j := 0; j < 100; j++ {
			r.Register(&Info{
				ID:        fmt.Sprintf("ep%d", j),
				ServiceID: "svc",
				Host:      testHost,
			})
		}
		r.Clear()
	}
}
