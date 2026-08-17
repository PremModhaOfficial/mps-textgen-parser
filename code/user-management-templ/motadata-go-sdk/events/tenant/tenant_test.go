package tenant

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/config"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/utils"
)

/* ========================================================================================================
   TEST CONSTANTS
   ======================================================================================================== */

var (
	testTenantID1  = "tenant-1"
	testTenantID2  = "tenant-2"
	testTenantID3  = "tenant-3"
	testTenantName = "Test Tenant"
	testNameOrig   = "Original"
	testNameUpd    = "Updated"
	testTenantTest = "Test"
	testUnknown    = "unknown"
	testUser       = "user"
	testPass       = "pass"
	testTestTenant = "test-tenant"
	testTestID     = "test"
	testJWT        = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig"
)

/* ========================================================================================================
   MOCK BUILDER
   ======================================================================================================== */

// mockBuilder returns a ConnectionBuilder that creates mock NATS connections.
// Since we cannot create real nats.Conn without a server, tests that call
// Connect will use a builder that returns an error or nil connections as needed.
func mockBuilder(returnErr error) ConnectionBuilder {
	return func(ctx context.Context, tenantID string, cfg *config.EventsConfig, creds *auth.Credentials) (*nats.Conn, jetstream.JetStream, error) {
		if returnErr != nil {
			return nil, nil, returnErr
		}
		return nil, nil, nil
	}
}

// newTestManager creates a Manager with mock builder and defers Shutdown.
func newTestManager(t *testing.T, builder ConnectionBuilder) *Manager {
	t.Helper()
	m, err := NewManager(ManagerConfig{Builder: builder})
	assert.New(t).NoError(err)
	t.Cleanup(func() { m.Shutdown(context.Background()) })
	return m
}

/* ========================================================================================================
   REGISTRY TESTS
   ======================================================================================================== */

func TestNewRegistry(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()
	assertions.NotNil(r, "NewRegistry returned nil")
	assertions.NotNil(r.tenants, "tenants map not initialized")
}

func TestRegistryRegister(t *testing.T) {
	type testCase struct {
		name    string
		info    *Info
		wantErr error
	}

	testCases := []testCase{
		{name: "nil info", info: nil, wantErr: utils.ErrTenantNotFound},
		{name: "empty ID", info: &Info{}, wantErr: utils.ErrTenantNotFound},
		{name: "valid registration", info: &Info{ID: testTenantID1, Name: testTenantName, Active: true}, wantErr: nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			r := NewRegistry()

			err := r.Register(tc.info)
			if tc.wantErr != nil {
				assertions.ErrorIs(err, tc.wantErr)
			} else {
				assertions.NoError(err)
			}
		})
	}

	t.Run("duplicate registration", func(t *testing.T) {
		assertions := assert.New(t)
		r := NewRegistry()
		info := &Info{ID: testTenantID1, Name: testTenantName, Active: true}
		assertions.NoError(r.Register(info))
		assertions.Equal(1, r.Count())
		assertions.ErrorIs(r.Register(info), utils.ErrTenantExists)
	})
}

func TestRegistryRegisterCallback(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()

	var registeredID string
	r.OnRegistered(func(tenantID string) {
		registeredID = tenantID
	})

	assertions.NoError(r.Register(&Info{ID: testTenantID1}))
	assertions.Equal(testTenantID1, registeredID)
}

func TestRegistryUnregister(t *testing.T) {
	type testCase struct {
		name     string
		tenantID string
		setup    bool
		wantErr  error
	}

	testCases := []testCase{
		{name: "existing tenant", tenantID: testTenantID1, setup: true, wantErr: nil},
		{name: "non-existent tenant", tenantID: testUnknown, setup: false, wantErr: utils.ErrTenantNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			r := NewRegistry()

			if tc.setup {
				assertions.NoError(r.Register(&Info{ID: tc.tenantID}))
			}

			err := r.Unregister(tc.tenantID)
			if tc.wantErr != nil {
				assertions.ErrorIs(err, tc.wantErr)
			} else {
				assertions.NoError(err)
				assertions.Equal(0, r.Count())
			}
		})
	}
}

func TestRegistryUnregisterCallback(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()
	assertions.NoError(r.Register(&Info{ID: testTenantID1}))

	var unregisteredID string
	r.OnUnregistered(func(tenantID string) {
		unregisteredID = tenantID
	})

	assertions.NoError(r.Unregister(testTenantID1))
	assertions.Equal(testTenantID1, unregisteredID)
}

func TestRegistryGet(t *testing.T) {
	type testCase struct {
		name     string
		tenantID string
		wantOK   bool
	}

	testCases := []testCase{
		{name: "existing tenant case", tenantID: testTenantID1, wantOK: true},
		{name: "non-existent tenant case", tenantID: testUnknown, wantOK: false},
	}

	r := NewRegistry()
	_ = r.Register(&Info{ID: testTenantID1, Name: testTenantTest})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			info, ok := r.Get(tc.tenantID)
			assertions.Equal(tc.wantOK, ok)
			if tc.wantOK {
				assertions.Equal(testTenantTest, info.Name)
			}
		})
	}
}

func TestRegistryExists(t *testing.T) {
	type testCase struct {
		name     string
		tenantID string
		want     bool
	}

	testCases := []testCase{
		{name: "existing-tenant", tenantID: testTenantID1, want: true},
		{name: "non-existent-tenant", tenantID: testUnknown, want: false},
	}

	r := NewRegistry()
	_ = r.Register(&Info{ID: testTenantID1})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.want, r.Exists(tc.tenantID))
		})
	}
}

func TestRegistryList(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()
	_ = r.Register(&Info{ID: testTenantID1})
	_ = r.Register(&Info{ID: testTenantID2})

	assertions.Len(r.List(), 2)
}

func TestRegistryAll(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()
	_ = r.Register(&Info{ID: testTenantID1, Name: "Test 1"})
	_ = r.Register(&Info{ID: testTenantID2, Name: "Test 2"})

	assertions.Len(r.All(), 2)
}

func TestRegistryCount(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()
	assertions.Equal(0, r.Count())

	_ = r.Register(&Info{ID: testTenantID1})
	assertions.Equal(1, r.Count())
}

func TestRegistryUpdate(t *testing.T) {
	type testCase struct {
		name     string
		tenantID string
		wantErr  error
	}

	testCases := []testCase{
		{name: "existing_tenant", tenantID: testTenantID1, wantErr: nil},
		{name: "non-existent_tenant", tenantID: testUnknown, wantErr: utils.ErrTenantNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			r := NewRegistry()
			_ = r.Register(&Info{ID: testTenantID1, Name: testNameOrig})

			err := r.Update(tc.tenantID, func(info *Info) {
				info.Name = testNameUpd
			})

			if tc.wantErr != nil {
				assertions.ErrorIs(err, tc.wantErr)
			} else {
				assertions.NoError(err)
				info, _ := r.Get(testTenantID1)
				assertions.Equal(testNameUpd, info.Name)
			}
		})
	}
}

func TestRegistrySetActive(t *testing.T) {
	type testCase struct {
		name     string
		tenantID string
		wantErr  error
	}

	testCases := []testCase{
		{name: "a existing tenant case", tenantID: testTenantID1, wantErr: nil},
		{name: "a non-existent tenant case", tenantID: testUnknown, wantErr: utils.ErrTenantNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			r := NewRegistry()
			_ = r.Register(&Info{ID: testTenantID1, Active: false})

			err := r.SetActive(tc.tenantID, true)
			if tc.wantErr != nil {
				assertions.ErrorIs(err, tc.wantErr)
			} else {
				assertions.NoError(err)
				info, _ := r.Get(testTenantID1)
				assertions.True(info.Active)
			}
		})
	}
}

func TestRegistryActiveTenants(t *testing.T) {
	assertions := assert.New(t)
	r := NewRegistry()
	_ = r.Register(&Info{ID: testTenantID1, Active: true})
	_ = r.Register(&Info{ID: testTenantID2, Active: false})
	_ = r.Register(&Info{ID: testTenantID3, Active: true})

	assertions.Len(r.ActiveTenants(), 2)
}

func TestRegistryConcurrentAccess(t *testing.T) {
	r := NewRegistry()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			tenantID := "tenant-" + string(rune('a'+id%26))
			r.Register(&Info{ID: tenantID})
		}(i)
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			tenantID := "tenant-" + string(rune('a'+id%26))
			r.Get(tenantID)
			r.Exists(tenantID)
			r.List()
		}(i)
	}

	wg.Wait()
}

/* ========================================================================================================
   MANAGER TESTS
   ======================================================================================================== */

func TestNewManager(t *testing.T) {
	assertions := assert.New(t)
	m, err := NewManager(ManagerConfig{Builder: mockBuilder(nil)})
	assertions.NoError(err)
	assertions.NotNil(m)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	m.Shutdown(ctx)
}

func TestNewManagerInvalidConfig(t *testing.T) {
	assertions := assert.New(t)
	_, err := NewManager(ManagerConfig{
		Config: &config.EventsConfig{CleanupInterval: 0},
	})
	assertions.Error(err)
}

func TestManagerConnect(t *testing.T) {
	type testCase struct {
		name      string
		builder   ConnectionBuilder
		creds     *auth.Credentials
		wantErr   bool
		wantTCNil bool
	}

	testCases := []testCase{
		{
			name:    "valid connect",
			builder: mockBuilder(nil),
			creds:   auth.UserPass(testUser, testPass),
			wantErr: false,
		},
		{
			name:    "nil credentials",
			builder: mockBuilder(nil),
			creds:   nil,
			wantErr: true,
		},
		{
			name:    "empty credentials",
			builder: mockBuilder(nil),
			creds:   &auth.Credentials{},
			wantErr: true,
		},
		{
			name:    "builder error",
			builder: mockBuilder(utils.ErrNotConnected),
			creds:   auth.UserPass(testUser, testPass),
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			m := newTestManager(t, tc.builder)

			conn, err := m.Connect(context.Background(), testTenantID1, tc.creds)
			if tc.wantErr {
				assertions.Error(err)
			} else {
				assertions.NoError(err)
				assertions.NotNil(conn)
				assertions.Equal(testTenantID1, conn.TenantID())
			}
		})
	}
}

func TestManagerConnectExisting(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t, mockBuilder(nil))
	ctx := context.Background()
	creds := auth.UserPass(testUser, testPass)

	tc1, _ := m.Connect(ctx, testTenantID1, creds)
	tc2, _ := m.Connect(ctx, testTenantID1, creds)
	assertions.Equal(tc1.TenantID(), tc2.TenantID())
}

func TestManagerConnectNoBuilder(t *testing.T) {
	assertions := assert.New(t)
	m, _ := NewManager(ManagerConfig{})
	defer m.Shutdown(context.Background())

	_, err := m.Connect(context.Background(), testTenantID1, auth.UserPass(testUser, testPass))
	assertions.Error(err)
}

func TestManagerConnectWithCredentialManager(t *testing.T) {
	type testCase struct {
		name       string
		setupCM    bool
		registerID string
		connectID  string
		wantErr    bool
	}

	testCases := []testCase{
		{
			name:       "with valid credentials",
			setupCM:    true,
			registerID: testTenantID1,
			connectID:  testTenantID1,
			wantErr:    false,
		},
		{
			name:      "without credential manager",
			setupCM:   false,
			connectID: testTenantID1,
			wantErr:   true,
		},
		{
			name:       "unknown tenant in credential manager",
			setupCM:    true,
			registerID: testTenantID1,
			connectID:  "unknown-tenant",
			wantErr:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			builder := mockBuilder(nil)
			cfg := ManagerConfig{Builder: builder}

			if tc.setupCM {
				cm := auth.NewCredentialManager(auth.DefaultCredentialManagerConfig())
				if tc.registerID != "" {
					cm.Register(tc.registerID, testJWT, "")
				}
				cfg.CredManager = cm
			}

			m, _ := NewManager(cfg)
			defer m.Shutdown(context.Background())

			conn, err := m.ConnectWithCredentialManager(context.Background(), tc.connectID)
			if tc.wantErr {
				assertions.Error(err)
			} else {
				assertions.NoError(err)
				assertions.NotNil(conn)
			}
		})
	}
}

func TestManagerConnectWithPayload(t *testing.T) {
	type testCase struct {
		name    string
		payload auth.RegistrationPayload
		wantErr bool
	}

	testCases := []testCase{
		{
			name:    "valid payload",
			payload: auth.RegistrationPayload{TenantID: testTenantID1, JWT: testJWT},
			wantErr: false,
		},
		{
			name:    "invalid payload",
			payload: auth.RegistrationPayload{},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			m := newTestManager(t, mockBuilder(nil))

			conn, err := m.ConnectWithPayload(context.Background(), tc.payload)
			if tc.wantErr {
				assertions.Error(err)
			} else {
				assertions.NoError(err)
				assertions.NotNil(conn)
			}
		})
	}
}

func TestManagerConnectWithPayloadAndCredentialManager(t *testing.T) {
	assertions := assert.New(t)
	builder := mockBuilder(nil)
	cm := auth.NewCredentialManager(auth.DefaultCredentialManagerConfig())

	m, _ := NewManager(ManagerConfig{Builder: builder, CredManager: cm})
	defer m.Shutdown(context.Background())

	payload := auth.RegistrationPayload{TenantID: testTenantID1, JWT: testJWT}
	tc, err := m.ConnectWithPayload(context.Background(), payload)
	assertions.NoError(err)
	assertions.NotNil(tc)
	assertions.True(cm.HasValid(testTenantID1))
}

func TestManagerConnectWithCredentialManagerFromManager(t *testing.T) {
	assertions := assert.New(t)
	builder := mockBuilder(nil)
	cm := auth.NewCredentialManager(auth.DefaultCredentialManagerConfig())
	cm.Register(testTenantID1, testJWT, "")

	m, _ := NewManager(ManagerConfig{Builder: builder, CredManager: cm})
	defer m.Shutdown(context.Background())

	tc, err := m.Connect(context.Background(), testTenantID1, nil)
	assertions.NoError(err)
	assertions.NotNil(tc)
}

func TestManagerGet(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t, mockBuilder(nil))
	creds := auth.UserPass(testUser, testPass)
	m.Connect(context.Background(), testTenantID1, creds)

	tc, ok := m.Get(testTenantID1)
	assertions.True(ok)
	assertions.NotNil(tc)

	_, ok = m.Get(testUnknown)
	assertions.False(ok)
}

func TestManagerDisconnect(t *testing.T) {
	type testCase struct {
		name     string
		tenantID string
		setup    bool
		wantErr  error
	}

	testCases := []testCase{
		{name: "existing connection", tenantID: testTenantID1, setup: true, wantErr: nil},
		{name: "non-existent connection", tenantID: testUnknown, setup: false, wantErr: utils.ErrTenantNotFound},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			m := newTestManager(t, mockBuilder(nil))
			ctx := context.Background()

			if tc.setup {
				m.Connect(ctx, testTenantID1, auth.UserPass(testUser, testPass))
			}

			err := m.Disconnect(ctx, tc.tenantID)
			if tc.wantErr != nil {
				assertions.ErrorIs(err, tc.wantErr)
			} else {
				assertions.NoError(err)
				_, ok := m.Get(testTenantID1)
				assertions.False(ok)
			}
		})
	}
}

func TestManagerShutdown(t *testing.T) {
	assertions := assert.New(t)
	m, _ := NewManager(ManagerConfig{Builder: mockBuilder(nil)})
	creds := auth.UserPass(testUser, testPass)
	ctx := context.Background()
	m.Connect(ctx, testTenantID1, creds)
	m.Connect(ctx, testTenantID2, creds)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	assertions.NoError(m.Shutdown(shutdownCtx))
	assertions.Equal(0, m.ActiveConnections())
}

func TestManagerShutdownTimeout(t *testing.T) {
	m, _ := NewManager(ManagerConfig{Builder: mockBuilder(nil)})
	creds := auth.UserPass(testUser, testPass)
	m.Connect(context.Background(), testTenantID1, creds)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Should not panic
	m.Shutdown(shutdownCtx)
}

func TestManagerErrors(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t, mockBuilder(nil))
	assertions.NotNil(m.Errors())
}

func TestManagerOnError(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t, mockBuilder(nil))

	var receivedErr utils.TenantError
	m.OnError(func(err utils.TenantError) {
		receivedErr = err
	})

	testErr := utils.TenantError{TenantID: testTenantID1, Err: utils.ErrTenantNotFound}
	m.dispatchError(testErr)
	assertions.Equal(testTenantID1, receivedErr.TenantID)
}

func TestManagerHealth(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t, mockBuilder(nil))
	creds := auth.UserPass(testUser, testPass)
	ctx := context.Background()
	m.Connect(ctx, testTenantID1, creds)
	m.Connect(ctx, testTenantID2, creds)

	health := m.Health()
	assertions.Len(health, 2)

	for tid, info := range health {
		assertions.Equal(tid, info.TenantID)
		assertions.False(info.Connected, "nil conn means IsConnected() returns false")
	}
}

func TestManagerTenantHealth(t *testing.T) {
	type testCase struct {
		name     string
		tenantID string
		wantOK   bool
	}

	testCases := []testCase{
		{name: "existing tenant usecase", tenantID: testTenantID1, wantOK: true},
		{name: "non-existent tenant usecase", tenantID: testUnknown, wantOK: false},
	}

	m := newTestManager(t, mockBuilder(nil))
	m.Connect(context.Background(), testTenantID1, auth.UserPass(testUser, testPass))

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			health, ok := m.TenantHealth(tc.tenantID)
			assertions.Equal(tc.wantOK, ok)
			if tc.wantOK {
				assertions.Equal(testTenantID1, health.TenantID)
			}
		})
	}
}

func TestManagerActiveConnections(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t, mockBuilder(nil))
	assertions.Equal(0, m.ActiveConnections())

	m.Connect(context.Background(), testTenantID1, auth.UserPass(testUser, testPass))
	assertions.Equal(1, m.ActiveConnections())
}

func TestManagerTenantIDs(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t, mockBuilder(nil))
	creds := auth.UserPass(testUser, testPass)
	ctx := context.Background()
	m.Connect(ctx, testTenantID1, creds)
	m.Connect(ctx, testTenantID2, creds)

	assertions.Len(m.TenantIDs(), 2)
}

func TestManagerSetCredentialManager(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t, mockBuilder(nil))

	cm := auth.NewCredentialManager(auth.DefaultCredentialManagerConfig())
	m.SetCredentialManager(cm)
	assertions.Equal(cm, m.CredentialManager())
}

func TestManagerConnectAfterShutdown(t *testing.T) {
	assertions := assert.New(t)
	m, _ := NewManager(ManagerConfig{Builder: mockBuilder(nil)})
	m.Shutdown(context.Background())

	_, err := m.Connect(context.Background(), testTenantID1, auth.UserPass(testUser, testPass))
	assertions.ErrorIs(err, utils.ErrShutdownInProgress)
}

func TestManagerConcurrentConnect(t *testing.T) {
	assertions := assert.New(t)
	m := newTestManager(t, mockBuilder(nil))
	creds := auth.UserPass(testUser, testPass)
	ctx := context.Background()
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.Connect(ctx, testTenantID1, creds)
		}()
	}

	wg.Wait()
	assertions.Equal(1, m.ActiveConnections())
}

func TestManagerDispatchErrorAfterShutdown(t *testing.T) {
	m, _ := NewManager(ManagerConfig{Builder: mockBuilder(nil)})
	m.Shutdown(context.Background())

	// Should not panic
	testErr := utils.TenantError{TenantID: testTenantID1, Err: utils.ErrTenantNotFound}
	m.dispatchError(testErr)
}

func TestManagerCleanupIdleConnections(t *testing.T) {
	assertions := assert.New(t)
	testConfig := config.DefaultEventsConfig()
	testConfig.IdleTimeout = 50 * time.Millisecond
	testConfig.CleanupInterval = 25 * time.Millisecond

	m, err := NewManager(ManagerConfig{
		Builder: mockBuilder(nil),
		Config:  testConfig,
	})
	assertions.NoError(err)
	defer m.Shutdown(context.Background())

	m.Connect(context.Background(), testTenantID1, auth.UserPass(testUser, testPass))
	time.Sleep(100 * time.Millisecond)

	assertions.Equal(0, m.ActiveConnections(), "connection should be cleaned up due to idle timeout")
}

func TestManagerDispatchErrorChannelFull(t *testing.T) {
	m, _ := NewManager(ManagerConfig{
		Builder:         mockBuilder(nil),
		ErrorBufferSize: 1,
	})
	defer m.Shutdown(context.Background())

	for i := 0; i < 10; i++ {
		testErr := utils.TenantError{TenantID: testTenantID1, Err: utils.ErrTenantNotFound}
		m.dispatchError(testErr)
	}
	// Should not panic
}

/* ========================================================================================================
   TENANT CONNECTION TESTS
   ======================================================================================================== */

func TestTenantConnectionTenantID(t *testing.T) {
	assertions := assert.New(t)
	tc := &TenantConnection{tenantID: testTestTenant}
	assertions.Equal(testTestTenant, tc.TenantID())
}

func TestTenantConnectionConn(t *testing.T) {
	assertions := assert.New(t)
	tc := &TenantConnection{}
	assertions.Nil(tc.Conn())
}

func TestTenantConnectionJetStream(t *testing.T) {
	assertions := assert.New(t)
	tc := &TenantConnection{}
	assertions.Nil(tc.JetStream())
}

func TestTenantConnectionTouch(t *testing.T) {
	assertions := assert.New(t)
	tc := &TenantConnection{tenantID: testTestID}
	tc.Touch()
	assertions.NotZero(tc.lastUsed.Load())
}

func TestTenantConnectionIsIdle(t *testing.T) {
	type testCase struct {
		name    string
		touch   bool
		sleep   time.Duration
		timeout time.Duration
		want    bool
	}

	testCases := []testCase{
		{name: "never touched", touch: false, timeout: time.Millisecond, want: false},
		{name: "idle after timeout", touch: true, sleep: 50 * time.Millisecond, timeout: 10 * time.Millisecond, want: true},
		{name: "not idle with long timeout", touch: true, sleep: 0, timeout: time.Hour, want: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			conn := &TenantConnection{tenantID: testTestID}

			if tc.touch {
				conn.Touch()
			}
			if tc.sleep > 0 {
				time.Sleep(tc.sleep)
			}

			assertions.Equal(tc.want, conn.IsIdle(tc.timeout))
		})
	}
}

func TestTenantConnectionIsConnected(t *testing.T) {
	assertions := assert.New(t)
	tc := &TenantConnection{}
	assertions.False(tc.IsConnected())
}

func TestTenantConnectionCloseNilConn(t *testing.T) {
	assertions := assert.New(t)
	tc := &TenantConnection{}
	assertions.NoError(tc.Close(context.Background()))
}
