package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

/* ========================================================================================================
   TEST CONSTANTS
   ======================================================================================================== */

var (
	testTenantTwo      = "tenant-2"
	testServiceName    = "my-service"
	testTenantBadJWT   = "tenant-bad-jwt"
	testInvalidJWTFmt  = "invalid-jwt-format"
	testResponseTenant = "t1"
	testResponseJWT    = "j1"
)

/* ========================================================================================================
   REGISTRATION PAYLOAD TESTS
   ======================================================================================================== */

func TestRegistrationPayloadValidate(t *testing.T) {

	testCases := []struct {
		name     string
		payload  *RegistrationPayload
		expected error
	}{
		{
			"valid payload",
			&RegistrationPayload{TenantID: testTenantID, JWT: testJWTTokenLiteral},
			nil,
		},
		{
			"missing tenant ID",
			&RegistrationPayload{JWT: testJWTTokenLiteral},
			ErrMissingTenantID,
		},
		{
			"missing JWT",
			&RegistrationPayload{TenantID: testTenantID},
			ErrMissingJWT,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			assertions := assert.New(t)

			err := tc.payload.Validate()

			if tc.expected == nil {
				assertions.NoError(err, "Should not return error")
			} else {
				assertions.ErrorIs(err, tc.expected, "Error should match")
			}
		})
	}
}

func TestParseRegistrationPayload(t *testing.T) {

	t.Run("valid payload", func(t *testing.T) {

		assertions := assert.New(t)

		data := []byte(`{"tenant_id":"` + testTenantID + `","jwt":"` + testJWTTokenLiteral + `"}`)
		p, err := ParseRegistrationPayload(data)

		assertions.NoError(err, "Should parse successfully")
		assertions.Equal(testTenantID, p.TenantID, "TenantID should match")
		assertions.Equal(testJWTTokenLiteral, p.JWT, "JWT should match")
	})

	t.Run("empty data", func(t *testing.T) {

		assertions := assert.New(t)

		_, err := ParseRegistrationPayload(nil)
		assertions.ErrorIs(err, ErrInvalidPayload, "Should return ErrInvalidPayload")
	})

	t.Run("invalid JSON", func(t *testing.T) {

		assertions := assert.New(t)

		_, err := ParseRegistrationPayload([]byte("not json"))
		assertions.ErrorIs(err, ErrInvalidPayload, "Should return an ErrInvalidPayload")
	})

	t.Run("missing fields", func(t *testing.T) {

		assertions := assert.New(t)

		_, err := ParseRegistrationPayload([]byte(`{"tenant_id":"` + testTenantID + `"}`))
		assertions.ErrorIs(err, ErrMissingJWT, "Should return ErrMissingJWT")
	})
}

func TestRegistrationRequestToJSON(t *testing.T) {

	assertions := assert.New(t)

	r := &RegistrationRequest{
		TenantID:    testTenantID,
		ServiceName: testServiceName,
		Permissions: &PermissionRequest{
			PublishSubjects:   []string{"events.>"},
			SubscribeSubjects: []string{"events.>"},
		},
		Metadata: map[string]string{"key": "value"},
	}

	data, err := r.ToJSON()

	assertions.NoError(err, "Should marshal to JSON")
	assertions.NotEmpty(data, "JSON should not be empty")
}

func TestParseRegistrationResponse(t *testing.T) {

	t.Run("valid response", func(t *testing.T) {

		assertions := assert.New(t)

		data := []byte(`{"success":true,"payload":{"tenant_id":"` + testResponseTenant + `","jwt":"` + testResponseJWT + `"}}`)
		r, err := ParseRegistrationResponse(data)

		assertions.NoError(err, "Should parse successfully")
		assertions.True(r.Success, "Success should be true")
		assertions.NotNil(r.Payload, "Payload should not be nil")
		assertions.Equal(testResponseTenant, r.Payload.TenantID, "TenantID must match")
	})

	t.Run("empty data", func(t *testing.T) {

		assertions := assert.New(t)

		_, err := ParseRegistrationResponse(nil)
		assertions.ErrorIs(err, ErrInvalidPayload, "Must return ErrInvalidPayload")
	})

	t.Run("invalid JSON", func(t *testing.T) {

		assertions := assert.New(t)

		_, err := ParseRegistrationResponse([]byte("not json"))
		assertions.ErrorIs(err, ErrInvalidPayload, "Must return an ErrInvalidPayload")
	})
}

/* ========================================================================================================
   REGISTRATION HANDLER TESTS
   ======================================================================================================== */

func TestNewRegistrationHandler(t *testing.T) {

	assertions := assert.New(t)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	h := NewRegistrationHandler(cm)

	assertions.NotNil(h, "Handler should not be nil")
	assertions.Equal(cm, h.credManager, "Credential manager should be set")
}

func TestRegistrationHandlerCallbacks(t *testing.T) {

	assertions := assert.New(t)

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
		TenantID: testTenantID,
		JWT:      testJWTToken,
	}
	err := h.HandleRegistration(payload)

	assertions.NoError(err, "Should register successfully")
	assertions.Equal(testTenantID, registeredID, "Registered callback should fire")

	// Test unregistration
	h.HandleUnregistration(testTenantID)
	assertions.Equal(testTenantID, unregisteredID, "Unregistered callback should fire")

	// Test error callback
	invalidPayload := &RegistrationPayload{TenantID: testResponseTenant}
	h.HandleRegistration(invalidPayload)
	assertions.Equal(testResponseTenant, errorID, "Error callback should fire")
	assertions.NotNil(errorErr, "Error should not be nil")
}

func TestRegistrationHandlerHandleNil(t *testing.T) {

	assertions := assert.New(t)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	h := NewRegistrationHandler(cm)

	err := h.HandleRegistration(nil)
	assertions.ErrorIs(err, ErrInvalidPayload, "Should be returned an ErrInvalidPayload for nil")
}

func TestRegistrationHandlerHandleJSON(t *testing.T) {

	assertions := assert.New(t)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	h := NewRegistrationHandler(cm)

	var registeredID string
	h.OnRegistered(func(tenantID string) {
		registeredID = tenantID
	})

	data := []byte(`{"tenant_id":"` + testTenantTwo + `","jwt":"` + testJWTToken + `"}`)
	err := h.HandleRegistrationJSON(data)

	assertions.NoError(err, "Should handle valid JSON")
	assertions.Equal(testTenantTwo, registeredID, "Should register tenant-2")

	err = h.HandleRegistrationJSON([]byte("invalid"))
	assertions.Error(err, "Should fail for invalid JSON")
}

func TestRegistrationHandlerUnregistrationEmpty(t *testing.T) {

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	h := NewRegistrationHandler(cm)

	// Should not panic
	h.HandleUnregistration("")
}

func TestHandleRegistrationInvalidJWT(t *testing.T) {

	assertions := assert.New(t)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	handler := NewRegistrationHandler(cm)

	var errorTenantID string
	var errorErr error
	handler.OnError(func(tenantID string, err error) {
		errorTenantID = tenantID
		errorErr = err
	})

	payload := &RegistrationPayload{
		TenantID: testTenantBadJWT,
		JWT:      testInvalidJWTFmt,
	}

	err := handler.HandleRegistration(payload)

	assertions.Error(err, "Should fail for invalid JWT")
	assertions.Equal(testTenantBadJWT, errorTenantID, "Error callback should receive tenant ID")
	assertions.NotNil(errorErr, "Error callback should receive error")
}

/* ========================================================================================================
   STATIC PROVIDER TESTS
   ======================================================================================================== */

func TestNewStaticProvider(t *testing.T) {

	assertions := assert.New(t)

	creds := UserPass(testUser, testPass)
	p := NewStaticProvider(creds)

	assertions.NotNil(p, "Provider should not be nil")
}

func TestStaticProviderType(t *testing.T) {

	assertions := assert.New(t)

	creds := UserPass(testUser, testPass)
	p := NewStaticProvider(creds)
	assertions.Equal(TypeUserPass, p.Type(), "Type should be UserPass")

	p = NewStaticProvider(nil)
	assertions.Equal(TypeNone, p.Type(), "Type should be None for nil credentials")
}

func TestStaticProviderAuthenticate(t *testing.T) {

	assertions := assert.New(t)
	ctx := context.Background()

	creds := Token(testTokenValue)
	p := NewStaticProvider(creds)
	result, err := p.Authenticate(ctx)
	assertions.NoError(err, "Should authenticate successfully")
	assertions.Equal(testTokenValue, result.Token, "Token should match")

	p = NewStaticProvider(nil)
	_, err = p.Authenticate(ctx)
	assertions.ErrorIs(err, ErrInvalidCredential, "Should return ErrInvalidCredential for nil")
}

func TestStaticProviderRefresh(t *testing.T) {

	assertions := assert.New(t)

	creds := Token(testTokenValue)
	p := NewStaticProvider(creds)

	result, err := p.Refresh(context.Background())
	assertions.NoError(err, "Should refresh successfully")
	assertions.Equal(creds, result, "Refresh should return same credentials")
}

func TestStaticProviderIsValid(t *testing.T) {

	assertions := assert.New(t)

	creds := Token(testTokenValue)
	p := NewStaticProvider(creds)
	assertions.True(p.IsValid(), "Should be valid with credentials")

	p = NewStaticProvider(nil)
	assertions.False(p.IsValid(), "Should be invalid for nil")

	p = NewStaticProvider(&Credentials{})
	assertions.False(p.IsValid(), "Should be invalid for empty")
}

/* ========================================================================================================
   MANAGED PROVIDER TESTS
   ======================================================================================================== */

func TestNewManagedProvider(t *testing.T) {

	assertions := assert.New(t)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	p := NewManagedProvider(cm, testTenantID)

	assertions.NotNil(p, "Provider should not be nil")
	assertions.Equal(testTenantID, p.tenantID, "TenantID should be matched")
}

func TestManagedProviderType(t *testing.T) {

	assertions := assert.New(t)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	p := NewManagedProvider(cm, testTenantID)

	assertions.Equal(TypeJWT, p.Type(), "Type should be JWT")
}

func TestManagedProviderAuthenticate(t *testing.T) {

	assertions := assert.New(t)
	ctx := context.Background()

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	cm.Register(testTenantID, testJWTToken, testSeedLiteral)

	p := NewManagedProvider(cm, testTenantID)
	creds, err := p.Authenticate(ctx)
	assertions.NoError(err, "Should authenticate successfully")
	assertions.Equal(testJWTToken, creds.JWT, "JWT should match")

	p2 := NewManagedProvider(cm, "unknown")
	_, err = p2.Authenticate(ctx)
	assertions.Error(err, "Should fail for unknown tenant")
}

func TestManagedProviderRefresh(t *testing.T) {

	assertions := assert.New(t)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	cm.Register(testTenantID, testJWTToken, testSeedLiteral)

	p := NewManagedProvider(cm, testTenantID)
	creds, err := p.Refresh(context.Background())

	assertions.NoError(err, "Should refresh successfully")
	assertions.NotNil(creds, "Credentials should not be nil")
}

func TestManagedProviderIsValid(t *testing.T) {

	assertions := assert.New(t)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	cm.Register(testTenantID, testJWTToken, testSeedLiteral)

	p := NewManagedProvider(cm, testTenantID)
	assertions.True(p.IsValid(), "Should be valid for registered tenant")

	p2 := NewManagedProvider(cm, "unknown")
	assertions.False(p2.IsValid(), "Should be invalid for unknown tenant")
}
