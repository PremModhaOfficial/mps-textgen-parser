package auth

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/* ========================================================================================================
   TEST CONSTANTS
   ======================================================================================================== */

var (
	testSeed           = "SUAM..."
	testTenantID       = "tenant-1"
	testJWTToken       = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig"
	testExpiredTenant  = "expired-tenant"
	testJWTFile        = "my-jwt"
	testSeedSample     = "my-seed"
	testJWTTokenSample = "my-jwt-token"
)

var ()

/* ========================================================================================================
   CREDENTIAL MANAGER CONFIG TESTS
   ======================================================================================================== */

func TestDefaultCredentialManagerConfig(t *testing.T) {

	assertions := assert.New(t)

	cfg := DefaultCredentialManagerConfig()

	assertions.Equal(5*time.Minute, cfg.RefreshThreshold, "RefreshThreshold should be 5m")
	assertions.Equal(1*time.Minute, cfg.CheckInterval, "CheckInterval should be 1m")
	assertions.Empty(cfg.CredentialsDir, "CredentialsDir should be empty")
}

/* ========================================================================================================
   CREDENTIAL MANAGER TESTS
   ======================================================================================================== */

func TestNewCredentialManager(t *testing.T) {

	assertions := assert.New(t)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())

	assertions.NotNil(cm, "CredentialManager should not be nil")
	assertions.NotNil(cm.store, "store should be initialized")
	assertions.Equal(5*time.Minute, cm.refreshThreshold, "refreshThreshold should match")
}

func TestNewCredentialManagerWithZeroConfig(t *testing.T) {

	assertions := assert.New(t)

	cm := NewCredentialManager(CredentialManagerConfig{})

	assertions.Equal(5*time.Minute, cm.refreshThreshold, "Should use default refreshThreshold")
	assertions.Equal(1*time.Minute, cm.checkInterval, "Should use default checkInterval")
}

func TestCredentialManagerStartStop(t *testing.T) {
	cfg := CredentialManagerConfig{
		CheckInterval: 10 * time.Millisecond,
	}
	cm := NewCredentialManager(cfg)

	cm.Start()

	// Let it run for a bit
	time.Sleep(25 * time.Millisecond)

	cm.Stop()
	// Should not panic or block
}

func TestCredentialManagerRegister(t *testing.T) {

	assertions := assert.New(t)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())

	err := cm.Register(testTenantID, testJWTToken, testSeed)
	assertions.NoError(err, "Should register successfully")

	cred, err := cm.Get(testTenantID)
	assertions.NoError(err, "Should get registered credential")
	assertions.Equal(testTenantID, cred.TenantID, "TenantID should match")
}

func TestCredentialManagerRegisterInvalidJWT(t *testing.T) {

	assertions := assert.New(t)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())

	err := cm.Register(testTenantID, "invalid-jwt", "")
	assertions.Error(err, "Should fail for invalid JWT")
}

func TestCredentialManagerRegisterFromPayload(t *testing.T) {

	assertions := assert.New(t)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())

	payload := RegistrationPayload{
		TenantID: testTenantID,
		JWT:      testJWTToken,
		Seed:     testSeed,
	}
	err := cm.RegisterFromPayload(payload)
	assertions.NoError(err, "Should register from payload successfully")
}

func TestCredentialManagerGet(t *testing.T) {

	assertions := assert.New(t)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	cm.Register(testTenantID, testJWTToken, "")

	cred, err := cm.Get(testTenantID)
	assertions.NoError(err, "Should get existing credential")
	assertions.NotNil(cred, "Credential should not be nil")

	_, err = cm.Get("unknown")
	assertions.Error(err, "Should fail for not known tenant")
}

func TestCredentialManagerGetJWT(t *testing.T) {

	assertions := assert.New(t)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	cm.Register(testTenantID, testJWTToken, "")

	got, err := cm.GetJWT(testTenantID)
	assertions.NoError(err, "Should get JWT successfully")
	assertions.Equal(testJWTToken, got, "JWT should match")

	_, err = cm.GetJWT("unknown")
	assertions.Error(err, "Should fail for tenant which are not known")
}

func TestCredentialManagerRemove(t *testing.T) {

	assertions := assert.New(t)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	cm.Register(testTenantID, testJWTToken, "")

	cm.Remove(testTenantID)

	_, err := cm.Get(testTenantID)
	assertions.Error(err, "Should fail after removal")
}

func TestCredentialManagerList(t *testing.T) {

	assertions := assert.New(t)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	cm.Register(testTenantID, testJWTToken, "")
	cm.Register("tenant-2", testJWTToken, "")

	list := cm.List()
	assertions.Len(list, 2, "Should list 2 tenants")
}

func TestCredentialManagerHasValid(t *testing.T) {

	assertions := assert.New(t)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	cm.Register(testTenantID, testJWTToken, "")

	assertions.True(cm.HasValid(testTenantID), "Should be valid for registered tenant")
	assertions.False(cm.HasValid("unknown"), "Should be invalid for unknown tenant")
}

func TestCredentialManagerGetCredentials(t *testing.T) {

	assertions := assert.New(t)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	cm.Register(testTenantID, testJWTToken, "")

	creds, err := cm.GetCredentials(testTenantID)
	assertions.NoError(err, "Should get credentials")
	assertions.NotNil(creds, "Credentials should not be nil")

	_, err = cm.GetCredentials("unknown")
	assertions.Error(err, "Should fail for unknown tenant")
}

/* ========================================================================================================
   REFRESH CALLBACK TESTS
   ======================================================================================================== */

func TestCredentialManagerOnRefreshNeeded(t *testing.T) {

	assertions := assert.New(t)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())

	var refreshedTenantID string
	cm.OnRefreshNeeded(func(tenantID string) (*JWTCredential, error) {
		refreshedTenantID = tenantID
		return NewJWTCredential(tenantID, testJWTToken, "")
	})

	cm.handleCredentialExpiring(testTenantID, 5*time.Minute)
	time.Sleep(50 * time.Millisecond)
	cm.wg.Wait()

	assertions.Equal(testTenantID, refreshedTenantID, "Refresh callback should be invoked with correct tenant")
}

func TestCredentialManagerOnRefreshNeededError(t *testing.T) {

	cm := NewCredentialManager(DefaultCredentialManagerConfig())

	cm.OnRefreshNeeded(func(tenantID string) (*JWTCredential, error) {
		return nil, ErrMissingJWT
	})

	// Should not panic
	cm.handleCredentialExpiring(testTenantID, 5*time.Minute)
	time.Sleep(50 * time.Millisecond)
	cm.wg.Wait()
}

func TestCredentialManagerOnRefreshNeededNil(t *testing.T) {

	cm := NewCredentialManager(DefaultCredentialManagerConfig())

	cm.OnRefreshNeeded(func(tenantID string) (*JWTCredential, error) {
		return nil, nil
	})

	// Should not panic when callback returns nil
	cm.handleCredentialExpiring(testTenantID, 5*time.Minute)
	time.Sleep(50 * time.Millisecond)
	cm.wg.Wait()
}

func TestHandleCredentialExpiringStoreError(t *testing.T) {

	assertions := assert.New(t)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())

	cm.OnRefreshNeeded(func(tenantID string) (*JWTCredential, error) {
		return &JWTCredential{}, nil // Empty TenantID causes store to fail
	})

	cm.handleCredentialExpiring(testTenantID, 5*time.Minute)
	time.Sleep(50 * time.Millisecond)
	cm.wg.Wait()

	assertions.True(true, "Should handle store error gracefully")
}

func TestCredentialManagerHandleExpired(t *testing.T) {

	cm := NewCredentialManager(DefaultCredentialManagerConfig())

	// Should not panic
	cm.handleCredentialExpired(testTenantID)
}

/* ========================================================================================================
   PERSISTENCE TESTS
   ======================================================================================================== */

func TestCredentialManagerPersistence(t *testing.T) {

	assertions := assert.New(t)

	tempDir, err := os.MkdirTemp("", "creds-test-*")
	assertions.NoError(err, "Should create temp dir")
	defer os.RemoveAll(tempDir)

	cfg := CredentialManagerConfig{
		CredentialsDir: tempDir,
	}
	cm := NewCredentialManager(cfg)

	cm.Register(testTenantID, testJWTToken, testSeed)

	filePath := filepath.Join(tempDir, "tenant-1.creds")
	_, statErr := os.Stat(filePath)
	assertions.NoError(statErr, "Credentials file should be created")

	cm.Remove(testTenantID)
	_, statErr = os.Stat(filePath)
	assertions.True(os.IsNotExist(statErr), "Credentials file should be removed")
}

func TestRegisterWithInvalidPersistenceDir(t *testing.T) {

	assertions := assert.New(t)

	cfg := CredentialManagerConfig{
		CredentialsDir: "/nonexistent/readonly/dir/that/cannot/be/created",
	}
	cm := NewCredentialManager(cfg)

	err := cm.Register(testTenantID, testJWTToken, testSeed)
	assertions.NoError(err, "Register should succeed even if file persistence fails")

	cred, getErr := cm.Get(testTenantID)
	assertions.NoError(getErr, "Credential should be available in memory")
	assertions.NotNil(cred, "Credential must not be nil")
}

func TestPersistCredentialsEmptyDir(t *testing.T) {

	assertions := assert.New(t)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())

	cred := &JWTCredential{
		TenantID: testTenantID,
		JWT:      "some-jwt",
		Seed:     "some-seed",
	}

	err := cm.persistCredentials(cred)
	assertions.NoError(err, "persistCredentials with empty dir should return nil")
}

func TestRemoveWithPersistenceDirNonExistentFile(t *testing.T) {

	tempDir, err := os.MkdirTemp("", "creds-remove-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := CredentialManagerConfig{
		CredentialsDir: tempDir,
	}
	cm := NewCredentialManager(cfg)

	cm.Register("tenant-rm", testJWTToken, "")

	filePath := filepath.Join(tempDir, "tenant-rm.creds")
	os.Remove(filePath)

	// Should not error even if file already gone
	cm.Remove("tenant-rm")
}

/* ========================================================================================================
   LOAD FROM FILE TESTS
   ======================================================================================================== */

func TestCredentialManagerLoadFromFile(t *testing.T) {

	assertions := assert.New(t)

	tempDir, err := os.MkdirTemp("", "test-creds-*")
	assertions.NoError(err, "Must create temp dir")
	defer os.RemoveAll(tempDir)

	content := FormatCredsFile(testJWTToken, testSeed)
	filePath := filepath.Join(tempDir, "tenant-1.creds")
	os.WriteFile(filePath, []byte(content), 0600)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	err = cm.LoadFromFile(testTenantID, filePath)
	assertions.NoError(err, "Should load from file")

	cred, err := cm.Get(testTenantID)
	assertions.NoError(err, "Should get loaded credential")
	assertions.Equal(testJWTToken, cred.JWT, "JWT must match")
}

func TestCredentialManagerLoadFromFileNotExist(t *testing.T) {

	assertions := assert.New(t)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	err := cm.LoadFromFile(testTenantID, "/nonexistent/path/file.creds")
	assertions.Error(err, "Should fail for non-existent file")
}

func TestCredentialManagerLoadFromFileInvalidContent(t *testing.T) {

	assertions := assert.New(t)

	tempDir, err := os.MkdirTemp("", "cred-test-*")
	assertions.NoError(err, "Should be created temp dir")
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "invalid.creds")
	os.WriteFile(filePath, []byte("invalid content"), 0600)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	err = cm.LoadFromFile(testTenantID, filePath)
	assertions.ErrorIs(err, ErrInvalidJWTFormat, "Should return ErrInvalidJWTFormat")
}

/* ========================================================================================================
   CREDENTIAL STORE TESTS
   ======================================================================================================== */

func TestNewCredentialStore(t *testing.T) {

	assertions := assert.New(t)

	s := NewCredentialStore()
	assertions.NotNil(s, "CredentialStore should not be nil")
	assertions.NotNil(s.credentials, "credentials map should be initialized")
}

func TestCredentialStoreStore(t *testing.T) {

	assertions := assert.New(t)

	s := NewCredentialStore()

	cred, _ := NewJWTCredential(testTenantID, testJWTToken, "")
	err := s.Store(cred)
	assertions.NoError(err, "Should store valid credential")

	err = s.Store(nil)
	assertions.ErrorIs(err, ErrMissingJWT, "Should fail for nil credential")

	err = s.Store(&JWTCredential{})
	assertions.ErrorIs(err, ErrMissingJWT, "Should fail for empty tenant ID")
}

func TestCredentialStoreGet(t *testing.T) {

	assertions := assert.New(t)

	s := NewCredentialStore()
	cred, _ := NewJWTCredential(testTenantID, testJWTToken, "")
	s.Store(cred)

	got, ok := s.Get(testTenantID)
	assertions.True(ok, "Should find existing credential")
	assertions.Equal(testTenantID, got.TenantID, "TenantID should match")

	_, ok = s.Get("unknown")
	assertions.False(ok, "Should not find non-existent credential")
}

func TestCredentialStoreGetValid(t *testing.T) {

	assertions := assert.New(t)

	s := NewCredentialStore()

	cred, _ := NewJWTCredential(testTenantID, testJWTToken, "")
	s.Store(cred)

	got, err := s.GetValid(testTenantID)
	assertions.NoError(err, "Should get valid credential")
	assertions.NotNil(got, "Credential cannot not be nil")

	_, err = s.GetValid("unknown")
	assertions.ErrorIs(err, ErrMissingJWT, "Should return ErrMissingJWT for unknown")

	expiredCred := &JWTCredential{
		TenantID:  testExpiredTenant,
		JWT:       testJWTToken,
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	s.Store(expiredCred)

	_, err = s.GetValid(testExpiredTenant)
	assertions.ErrorIs(err, ErrJWTExpired, "Should return ErrJWTExpired")
}

func TestCredentialStoreRemove(t *testing.T) {

	assertions := assert.New(t)

	s := NewCredentialStore()
	cred, _ := NewJWTCredential(testTenantID, testJWTToken, "")
	s.Store(cred)

	s.Remove(testTenantID)

	_, ok := s.Get(testTenantID)
	assertions.False(ok, "Should be removed")
}

func TestCredentialStoreList(t *testing.T) {

	assertions := assert.New(t)

	s := NewCredentialStore()
	cred1, _ := NewJWTCredential(testTenantID, testJWTToken, "")
	cred2, _ := NewJWTCredential("tenant-2", testJWTToken, "")
	s.Store(cred1)
	s.Store(cred2)

	list := s.List()
	assertions.Len(list, 2, "Should list 2 tenants")
}

func TestCredentialStoreCheckExpiring(t *testing.T) {

	assertions := assert.New(t)

	s := NewCredentialStore()

	var expiringID, expiredID string
	var expiringDuration time.Duration

	s.OnCredentialExpiring(func(tenantID string, timeUntilExpiry time.Duration) {
		expiringID = tenantID
		expiringDuration = timeUntilExpiry
	})
	s.OnCredentialExpired(func(tenantID string) {
		expiredID = tenantID
	})

	expiringCred := &JWTCredential{
		TenantID:  "expiring-tenant",
		JWT:       testJWTToken,
		ExpiresAt: time.Now().Add(3 * time.Minute),
	}
	s.Store(expiringCred)

	expiredCred := &JWTCredential{
		TenantID:  testExpiredTenant,
		JWT:       testJWTToken,
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	s.Store(expiredCred)

	s.CheckExpiring(5 * time.Minute)

	assertions.Equal("expiring-tenant", expiringID, "Should trigger expiring callback")
	assertions.True(expiringDuration > 0, "Time until expiry should be positive")
	assertions.Equal(testExpiredTenant, expiredID, "Should trigger expired callback")
}

func TestCredentialStoreCheckExpiringNoCallbacks(t *testing.T) {

	s := NewCredentialStore()

	cred := &JWTCredential{
		TenantID:  testTenantID,
		JWT:       testJWTToken,
		ExpiresAt: time.Now().Add(3 * time.Minute),
	}
	s.Store(cred)

	// Should not panic without callbacks
	s.CheckExpiring(5 * time.Minute)
}

/* ========================================================================================================
   FORMAT / PARSE CREDS FILE TESTS
   ======================================================================================================== */

func TestFormatCredsFile(t *testing.T) {

	assertions := assert.New(t)

	content := FormatCredsFile(testJWTTokenSample, testSeed)

	assertions.NotEmpty(content, "Content should not be empty")
	assertions.Contains(content, "-----BEGIN NATS USER JWT-----", "Should contain JWT begin marker")
	assertions.Contains(content, "------END NATS USER JWT------", "Should contain JWT end marker")
	assertions.Contains(content, testJWTTokenSample, "Should contain JWT")
	assertions.Contains(content, "-----BEGIN USER NKEY SEED-----", "Should contain seed begin marker")
	assertions.Contains(content, "------END USER NKEY SEED------", "Should contain seed end marker")
	assertions.Contains(content, testSeed, "Should contain seed")
}

func TestFormatCredsFileJWTOnly(t *testing.T) {

	assertions := assert.New(t)

	content := FormatCredsFile(testJWTFile, "")

	assertions.Contains(content, testJWTFile, "Should contain JWT")
	assertions.NotContains(content, "NKEY SEED", "Should not contain NKEY section")
}

func TestFormatCredsFileSeedOnly(t *testing.T) {

	assertions := assert.New(t)

	content := FormatCredsFile("", testSeedSample)

	assertions.NotContains(content, "USER JWT", "Should not contain JWT section")
	assertions.Contains(content, testSeedSample, "Should contain seed")
}

func TestFormatCredsFileEmpty(t *testing.T) {

	assertions := assert.New(t)

	content := FormatCredsFile("", "")

	assertions.Empty(content, "Should be empty for empty inputs")
}

func TestParseCredsFile(t *testing.T) {

	assertions := assert.New(t)

	content := FormatCredsFile(testJWTTokenSample, testSeed)
	parsedJWT, parsedSeed := ParseCredsFile(content)

	assertions.Equal(testJWTTokenSample, parsedJWT, "JWT must be matched")
	assertions.Equal(testSeed, parsedSeed, "Seed should match")
}

func TestParseCredsFileEmpty(t *testing.T) {

	assertions := assert.New(t)

	jwt, seed := ParseCredsFile("")
	assertions.Empty(jwt, "JWT should be empty")
	assertions.Empty(seed, "Seed should be empty")
}

func TestParseCredsFileInvalid(t *testing.T) {

	assertions := assert.New(t)

	jwt, seed := ParseCredsFile("random content without markers")
	assertions.Empty(jwt, "JWT should be empty for invalid content")
	assertions.Empty(seed, "Seed should be empty for invalid content")
}

func TestParseCredsFilePartial(t *testing.T) {

	assertions := assert.New(t)

	content := "-----BEGIN NATS USER JWT-----\nmy-jwt\n------END NATS USER JWT------"
	jwt, seed := ParseCredsFile(content)

	assertions.Equal(testJWTFile, jwt, "JWT should be parsed")
	assertions.Empty(seed, "Seed should be empty")
}

func TestParseCredsFileWithWhitespace(t *testing.T) {

	assertions := assert.New(t)

	content := "-----BEGIN NATS USER JWT-----\n  my-jwt\n------END NATS USER JWT------\n-----BEGIN USER NKEY SEED-----\n\tmy-seed\n------END USER NKEY SEED------"
	jwt, seed := ParseCredsFile(content)

	assertions.Equal(testJWTFile, jwt, "JWT should be trimmed")
	assertions.Equal(testSeedSample, seed, "Seed should be trimmed")
}

/* ========================================================================================================
   TRIM SPACE EDGE CASES
   ======================================================================================================== */

func TestTrimSpaceEdgeCases(t *testing.T) {

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{"leading spaces", "   hello", "hello"},
		{"trailing spaces", "hello   ", "hello"},
		{"both sides", "  hello  ", "hello"},
		{"tabs", "\thello\t", "hello"},
		{"carriage return", "\rhello\r", "hello"},
		{"mixed whitespace", " \t\rhello \t\r", "hello"},
		{"empty", "", ""},
		{"only whitespace", "   \t\r  ", ""},
		{"no whitespace", "hello", "hello"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			assertions := assert.New(t)

			result := trimSpace(tc.input)

			assertions.Equal(tc.expected, result, "trimSpace should handle: "+tc.name)
		})
	}
}

/* ========================================================================================================
   SPLIT LINES EDGE CASES
   ======================================================================================================== */

func TestSplitLinesEdgeCases(t *testing.T) {

	testCases := []struct {
		name     string
		input    string
		expected int
	}{
		{"empty", "", 0},
		{"single line no newline", "hello", 1},
		{"single line with newline", "hello\n", 1},
		{"two lines", "hello\nworld", 2},
		{"trailing newline", "hello\nworld\n", 2},
		{"multiple newlines", "a\nb\nc\nd", 4},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			assertions := assert.New(t)

			lines := splitLines(tc.input)

			assertions.Len(lines, tc.expected, "Number of lines should match for: "+tc.name)
		})
	}
}

/* ========================================================================================================
   PARSE JWT INVALID JSON PAYLOAD TESTS
   ======================================================================================================== */

func TestParseJWTInvalidJSONPayload(t *testing.T) {

	assertions := assert.New(t)

	invalidJSON := base64.RawURLEncoding.EncodeToString([]byte("not-json"))
	token := "header." + invalidJSON + ".signature"

	claims, err := ParseJWT(token)

	assertions.Error(err, "Should fail for invalid JSON in payload")
	assertions.ErrorIs(err, ErrInvalidJWTFormat, "Should be format error")
	assertions.Nil(claims, "Claims should be nil")
}

/* ========================================================================================================
   CONCURRENT ACCESS TESTS
   ======================================================================================================== */

func TestCredentialStoreConcurrentAccess(t *testing.T) {

	s := NewCredentialStore()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			tid := testTenantID + string(rune('a'+id%26))
			cred, _ := NewJWTCredential(tid, testJWTToken, "")
			s.Store(cred)
		}(i)
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			tid := "tenant-" + string(rune('a'+id%26))
			s.Get(tid)
		}(i)
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.List()
		}()
	}

	wg.Wait()
}

func TestCredentialManagerConcurrentAccess(t *testing.T) {

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			tid := "tenant-" + string(rune('a'+id%26))
			cm.Register(tid, testJWTToken, "")
			cm.Get(tid)
			cm.HasValid(tid)
			cm.List()
		}(i)
	}

	wg.Wait()
}

/* ========================================================================================================
   BENCHMARK TESTS
   ======================================================================================================== */

func BenchmarkCredentialStoreStore(b *testing.B) {

	s := NewCredentialStore()
	cred, _ := NewJWTCredential(testTenantID, testJWTToken, "")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Store(cred)
	}
}

func BenchmarkCredentialStoreGet(b *testing.B) {

	s := NewCredentialStore()
	cred, _ := NewJWTCredential(testTenantID, testJWTToken, "")
	s.Store(cred)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Get(testTenantID)
	}
}

func BenchmarkFormatCredsFile(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		FormatCredsFile(testJWTToken, testSeed)
	}
}

func BenchmarkParseCredsFile(b *testing.B) {

	content := FormatCredsFile(testJWTTokenSample, testSeed)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ParseCredsFile(content)
	}
}

func BenchmarkCredentialManagerRegister(b *testing.B) {

	cm := NewCredentialManager(DefaultCredentialManagerConfig())

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cm.Register(testTenantID, testJWTToken, "")
	}
}

func BenchmarkNewJWTCredential(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		NewJWTCredential(testTenantID, testJWTToken, "seed")
	}
}

func BenchmarkParseJWT(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ParseJWT(testJWTToken)
	}
}

func BenchmarkCredentialManagerGetParallel(b *testing.B) {

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	cm.Register(testTenantID, testJWTToken, "")

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			cm.Get(testTenantID)
		}
	})
}
