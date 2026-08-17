package auth

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// ============== CredentialManagerConfig Tests ==============

func TestDefaultCredentialManagerConfig(t *testing.T) {
	cfg := DefaultCredentialManagerConfig()
	if cfg.RefreshThreshold != 5*time.Minute {
		t.Errorf("expected RefreshThreshold 5m, got %v", cfg.RefreshThreshold)
	}
	if cfg.CheckInterval != 1*time.Minute {
		t.Errorf("expected CheckInterval 1m, got %v", cfg.CheckInterval)
	}
	if cfg.CredentialsDir != "" {
		t.Errorf("expected empty CredentialsDir, got %s", cfg.CredentialsDir)
	}
}

// ============== CredentialManager Tests ==============

func TestNewCredentialManager(t *testing.T) {
	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	if cm == nil {
		t.Fatal("NewCredentialManager returned nil")
	}
	if cm.store == nil {
		t.Error("store not initialized")
	}
	if cm.refreshThreshold != 5*time.Minute {
		t.Errorf("unexpected refreshThreshold: %v", cm.refreshThreshold)
	}
}

func TestNewCredentialManagerWithZeroConfig(t *testing.T) {
	cm := NewCredentialManager(CredentialManagerConfig{})
	if cm.refreshThreshold != 5*time.Minute {
		t.Errorf("expected default refreshThreshold, got %v", cm.refreshThreshold)
	}
	if cm.checkInterval != 1*time.Minute {
		t.Errorf("expected default checkInterval, got %v", cm.checkInterval)
	}
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
	cm := NewCredentialManager(DefaultCredentialManagerConfig())

	// Valid registration
	err := cm.Register("tenant-1", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig", "SUAM...")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should be able to get it
	cred, err := cm.Get("tenant-1")
	if err != nil {
		t.Fatalf("unexpected error getting credential: %v", err)
	}
	if cred.TenantID != "tenant-1" {
		t.Errorf("unexpected tenant ID: %s", cred.TenantID)
	}
}

func TestCredentialManagerRegisterInvalidJWT(t *testing.T) {
	cm := NewCredentialManager(DefaultCredentialManagerConfig())

	// Invalid JWT format
	err := cm.Register("tenant-1", "invalid-jwt", "")
	if err == nil {
		t.Error("expected error for invalid JWT")
	}
}

func TestCredentialManagerRegisterFromPayload(t *testing.T) {
	cm := NewCredentialManager(DefaultCredentialManagerConfig())

	payload := RegistrationPayload{
		TenantID: "tenant-1",
		JWT:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig",
		Seed:     "SUAM...",
	}
	err := cm.RegisterFromPayload(payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCredentialManagerGet(t *testing.T) {
	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	cm.Register("tenant-1", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig", "")

	// Get existing
	cred, err := cm.Get("tenant-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cred == nil {
		t.Error("expected credential")
	}

	// Get non-existent
	_, err = cm.Get("unknown")
	if err == nil {
		t.Error("expected error for unknown tenant")
	}
}

func TestCredentialManagerGetJWT(t *testing.T) {
	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig"
	cm.Register("tenant-1", jwt, "")

	// Get existing
	got, err := cm.GetJWT("tenant-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != jwt {
		t.Errorf("unexpected JWT: %s", got)
	}

	// Get non-existent
	_, err = cm.GetJWT("unknown")
	if err == nil {
		t.Error("expected error for unknown tenant")
	}
}

func TestCredentialManagerRemove(t *testing.T) {
	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	cm.Register("tenant-1", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig", "")

	cm.Remove("tenant-1")

	_, err := cm.Get("tenant-1")
	if err == nil {
		t.Error("expected error after removal")
	}
}

func TestCredentialManagerList(t *testing.T) {
	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	cm.Register("tenant-1", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig", "")
	cm.Register("tenant-2", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig", "")

	list := cm.List()
	if len(list) != 2 {
		t.Errorf("expected 2 tenants, got %d", len(list))
	}
}

func TestCredentialManagerHasValid(t *testing.T) {
	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	cm.Register("tenant-1", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig", "")

	if !cm.HasValid("tenant-1") {
		t.Error("expected HasValid to return true")
	}

	if cm.HasValid("unknown") {
		t.Error("expected HasValid to return false for unknown")
	}
}

func TestCredentialManagerGetCredentials(t *testing.T) {
	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	cm.Register("tenant-1", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig", "")

	creds, err := cm.GetCredentials("tenant-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if creds == nil {
		t.Error("expected credentials")
	}

	_, err = cm.GetCredentials("unknown")
	if err == nil {
		t.Error("expected error for unknown tenant")
	}
}

func TestCredentialManagerOnRefreshNeeded(t *testing.T) {
	cm := NewCredentialManager(DefaultCredentialManagerConfig())

	var refreshedTenantID string
	cm.OnRefreshNeeded(func(tenantID string) (*JWTCredential, error) {
		refreshedTenantID = tenantID
		return NewJWTCredential(tenantID, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.new", "")
	})

	// Simulate expiring credential callback
	cm.handleCredentialExpiring("tenant-1", 5*time.Minute)

	// Wait for goroutine to complete
	time.Sleep(50 * time.Millisecond)
	cm.wg.Wait()

	if refreshedTenantID != "tenant-1" {
		t.Errorf("expected refresh callback for tenant-1, got %s", refreshedTenantID)
	}
}

func TestCredentialManagerOnRefreshNeededError(t *testing.T) {
	cm := NewCredentialManager(DefaultCredentialManagerConfig())

	cm.OnRefreshNeeded(func(tenantID string) (*JWTCredential, error) {
		return nil, ErrMissingJWT
	})

	// Should not panic
	cm.handleCredentialExpiring("tenant-1", 5*time.Minute)
	time.Sleep(50 * time.Millisecond)
	cm.wg.Wait()
}

func TestCredentialManagerOnRefreshNeededNil(t *testing.T) {
	cm := NewCredentialManager(DefaultCredentialManagerConfig())

	cm.OnRefreshNeeded(func(tenantID string) (*JWTCredential, error) {
		return nil, nil
	})

	// Should not panic when callback returns nil
	cm.handleCredentialExpiring("tenant-1", 5*time.Minute)
	time.Sleep(50 * time.Millisecond)
	cm.wg.Wait()
}

func TestCredentialManagerHandleExpired(t *testing.T) {
	cm := NewCredentialManager(DefaultCredentialManagerConfig())

	// Should not panic
	cm.handleCredentialExpired("tenant-1")
}

func TestCredentialManagerPersistence(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "creds-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cfg := CredentialManagerConfig{
		CredentialsDir: tempDir,
	}
	cm := NewCredentialManager(cfg)

	// Register should persist
	cm.Register("tenant-1", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig", "SUAM...")

	// Check file exists
	filePath := filepath.Join(tempDir, "tenant-1.creds")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("credentials file not created")
	}

	// Remove should delete file
	cm.Remove("tenant-1")
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("credentials file not removed")
	}
}

func TestCredentialManagerLoadFromFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "creds-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create credentials file
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig"
	content := FormatCredsFile(jwt, "SUAM...")
	filePath := filepath.Join(tempDir, "tenant-1.creds")
	os.WriteFile(filePath, []byte(content), 0600)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	err = cm.LoadFromFile("tenant-1", filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should be able to get credentials
	cred, err := cm.Get("tenant-1")
	if err != nil {
		t.Fatalf("unexpected error getting credential: %v", err)
	}
	if cred.JWT != jwt {
		t.Errorf("unexpected JWT: %s", cred.JWT)
	}
}

func TestCredentialManagerLoadFromFileNotExist(t *testing.T) {
	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	err := cm.LoadFromFile("tenant-1", "/nonexistent/path/file.creds")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestCredentialManagerLoadFromFileInvalidContent(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "creds-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create file with invalid content
	filePath := filepath.Join(tempDir, "invalid.creds")
	os.WriteFile(filePath, []byte("invalid content"), 0600)

	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	err = cm.LoadFromFile("tenant-1", filePath)
	if err != ErrInvalidJWTFormat {
		t.Errorf("expected ErrInvalidJWTFormat, got %v", err)
	}
}

// ============== CredentialStore Tests ==============

func TestNewCredentialStore(t *testing.T) {
	s := NewCredentialStore()
	if s == nil {
		t.Fatal("NewCredentialStore returned nil")
	}
	if s.credentials == nil {
		t.Error("credentials map not initialized")
	}
}

func TestCredentialStoreStore(t *testing.T) {
	s := NewCredentialStore()

	// Valid credential
	cred, _ := NewJWTCredential("tenant-1", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig", "")
	err := s.Store(cred)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Nil credential
	err = s.Store(nil)
	if err != ErrMissingJWT {
		t.Errorf("expected ErrMissingJWT, got %v", err)
	}

	// Empty tenant ID
	err = s.Store(&JWTCredential{})
	if err != ErrMissingJWT {
		t.Errorf("expected ErrMissingJWT for empty tenant, got %v", err)
	}
}

func TestCredentialStoreGet(t *testing.T) {
	s := NewCredentialStore()
	cred, _ := NewJWTCredential("tenant-1", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig", "")
	s.Store(cred)

	// Get existing
	got, ok := s.Get("tenant-1")
	if !ok {
		t.Error("expected credential to exist")
	}
	if got.TenantID != "tenant-1" {
		t.Errorf("unexpected tenant ID: %s", got.TenantID)
	}

	// Get non-existent
	_, ok = s.Get("unknown")
	if ok {
		t.Error("expected credential to not exist")
	}
}

func TestCredentialStoreGetValid(t *testing.T) {
	s := NewCredentialStore()

	// Store valid credential
	cred, _ := NewJWTCredential("tenant-1", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig", "")
	s.Store(cred)

	got, err := s.GetValid("tenant-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Error("expected credential")
	}

	// Get non-existent
	_, err = s.GetValid("unknown")
	if err != ErrMissingJWT {
		t.Errorf("expected ErrMissingJWT, got %v", err)
	}

	// Store expired credential
	expiredCred := &JWTCredential{
		TenantID:  "expired-tenant",
		JWT:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjF9.sig",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	s.Store(expiredCred)

	_, err = s.GetValid("expired-tenant")
	if err != ErrJWTExpired {
		t.Errorf("expected ErrJWTExpired, got %v", err)
	}
}

func TestCredentialStoreRemove(t *testing.T) {
	s := NewCredentialStore()
	cred, _ := NewJWTCredential("tenant-1", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig", "")
	s.Store(cred)

	s.Remove("tenant-1")

	_, ok := s.Get("tenant-1")
	if ok {
		t.Error("expected credential to be removed")
	}
}

func TestCredentialStoreList(t *testing.T) {
	s := NewCredentialStore()
	cred1, _ := NewJWTCredential("tenant-1", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig", "")
	cred2, _ := NewJWTCredential("tenant-2", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig", "")
	s.Store(cred1)
	s.Store(cred2)

	list := s.List()
	if len(list) != 2 {
		t.Errorf("expected 2 tenants, got %d", len(list))
	}
}

func TestCredentialStoreCheckExpiring(t *testing.T) {
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

	// Store expiring credential (expires in 3 minutes, threshold is 5 minutes)
	expiringCred := &JWTCredential{
		TenantID:  "expiring-tenant",
		JWT:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig",
		ExpiresAt: time.Now().Add(3 * time.Minute),
	}
	s.Store(expiringCred)

	// Store expired credential
	expiredCred := &JWTCredential{
		TenantID:  "expired-tenant",
		JWT:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjF9.sig",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	s.Store(expiredCred)

	// Check with 5 minute threshold
	s.CheckExpiring(5 * time.Minute)

	if expiringID != "expiring-tenant" {
		t.Errorf("expected expiring callback for expiring-tenant, got %s", expiringID)
	}
	if expiringDuration <= 0 {
		t.Error("expected positive time until expiry")
	}
	if expiredID != "expired-tenant" {
		t.Errorf("expected expired callback for expired-tenant, got %s", expiredID)
	}
}

func TestCredentialStoreCheckExpiringNoCallbacks(t *testing.T) {
	s := NewCredentialStore()

	// Store credentials without setting callbacks
	cred := &JWTCredential{
		TenantID:  "tenant-1",
		JWT:       "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig",
		ExpiresAt: time.Now().Add(3 * time.Minute),
	}
	s.Store(cred)

	// Should not panic
	s.CheckExpiring(5 * time.Minute)
}

// ============== Format/Parse Creds File Tests ==============

func TestFormatCredsFile(t *testing.T) {
	jwt := "my-jwt-token"
	seed := "SUAM..."

	content := FormatCredsFile(jwt, seed)

	if content == "" {
		t.Error("expected non-empty content")
	}

	// Should contain JWT markers
	if !containsString(content, "-----BEGIN NATS USER JWT-----") {
		t.Error("missing JWT begin marker")
	}
	if !containsString(content, "------END NATS USER JWT------") {
		t.Error("missing JWT end marker")
	}
	if !containsString(content, jwt) {
		t.Error("JWT not found in content")
	}

	// Should contain NKEY markers
	if !containsString(content, "-----BEGIN USER NKEY SEED-----") {
		t.Error("missing seed begin marker")
	}
	if !containsString(content, "------END USER NKEY SEED------") {
		t.Error("missing seed end marker")
	}
	if !containsString(content, seed) {
		t.Error("seed not found in content")
	}
}

func TestFormatCredsFileJWTOnly(t *testing.T) {
	content := FormatCredsFile("my-jwt", "")
	if !containsString(content, "my-jwt") {
		t.Error("JWT not found")
	}
	if containsString(content, "NKEY SEED") {
		t.Error("should not contain NKEY section for empty seed")
	}
}

func TestFormatCredsFileSeedOnly(t *testing.T) {
	content := FormatCredsFile("", "my-seed")
	if containsString(content, "USER JWT") {
		t.Error("should not contain JWT section for empty JWT")
	}
	if !containsString(content, "my-seed") {
		t.Error("seed not found")
	}
}

func TestFormatCredsFileEmpty(t *testing.T) {
	content := FormatCredsFile("", "")
	if content != "" {
		t.Error("expected empty content for empty inputs")
	}
}

func TestParseCredsFile(t *testing.T) {
	jwt := "my-jwt-token"
	seed := "SUAM..."

	content := FormatCredsFile(jwt, seed)
	parsedJWT, parsedSeed := ParseCredsFile(content)

	if parsedJWT != jwt {
		t.Errorf("expected JWT %s, got %s", jwt, parsedJWT)
	}
	if parsedSeed != seed {
		t.Errorf("expected seed %s, got %s", seed, parsedSeed)
	}
}

func TestParseCredsFileEmpty(t *testing.T) {
	jwt, seed := ParseCredsFile("")
	if jwt != "" || seed != "" {
		t.Error("expected empty results for empty input")
	}
}

func TestParseCredsFileInvalid(t *testing.T) {
	jwt, seed := ParseCredsFile("random content without markers")
	if jwt != "" || seed != "" {
		t.Error("expected empty results for invalid input")
	}
}

func TestParseCredsFilePartial(t *testing.T) {
	// Only JWT
	content := `-----BEGIN NATS USER JWT-----
my-jwt
------END NATS USER JWT------`
	jwt, seed := ParseCredsFile(content)
	if jwt != "my-jwt" {
		t.Errorf("expected jwt 'my-jwt', got %s", jwt)
	}
	if seed != "" {
		t.Error("expected empty seed")
	}
}

func TestParseCredsFileWithWhitespace(t *testing.T) {
	content := `-----BEGIN NATS USER JWT-----
  my-jwt
------END NATS USER JWT------
-----BEGIN USER NKEY SEED-----
	my-seed
------END USER NKEY SEED------`
	jwt, seed := ParseCredsFile(content)
	if jwt != "my-jwt" {
		t.Errorf("expected jwt 'my-jwt', got '%s'", jwt)
	}
	if seed != "my-seed" {
		t.Errorf("expected seed 'my-seed', got '%s'", seed)
	}
}

// ============== Concurrent Access Tests ==============

func TestCredentialStoreConcurrentAccess(t *testing.T) {
	s := NewCredentialStore()
	var wg sync.WaitGroup

	// Concurrent stores
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			tenantID := "tenant-" + string(rune('a'+id%26))
			cred, _ := NewJWTCredential(tenantID, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig", "")
			s.Store(cred)
		}(i)
	}

	// Concurrent gets
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			tenantID := "tenant-" + string(rune('a'+id%26))
			s.Get(tenantID)
		}(i)
	}

	// Concurrent lists
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

	// Concurrent operations
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			tenantID := "tenant-" + string(rune('a'+id%26))
			cm.Register(tenantID, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig", "")
			cm.Get(tenantID)
			cm.HasValid(tenantID)
			cm.List()
		}(i)
	}

	wg.Wait()
}

// ============== Helper Functions ==============

func containsString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ============== Benchmarks ==============

func BenchmarkCredentialStoreStore(b *testing.B) {
	s := NewCredentialStore()
	cred, _ := NewJWTCredential("tenant-1", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig", "")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		s.Store(cred)
	}
}

func BenchmarkCredentialStoreGet(b *testing.B) {
	s := NewCredentialStore()
	cred, _ := NewJWTCredential("tenant-1", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig", "")
	s.Store(cred)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		s.Get("tenant-1")
	}
}

func BenchmarkFormatCredsFile(b *testing.B) {
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig"
	seed := "SUAM..."

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		FormatCredsFile(jwt, seed)
	}
}

func BenchmarkParseCredsFile(b *testing.B) {
	content := FormatCredsFile("my-jwt-token", "SUAM...")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		ParseCredsFile(content)
	}
}

func BenchmarkCredentialManagerRegister(b *testing.B) {
	cm := NewCredentialManager(DefaultCredentialManagerConfig())
	jwt := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk5OTk5OTk5OTl9.sig"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		cm.Register("tenant-1", jwt, "")
	}
}
