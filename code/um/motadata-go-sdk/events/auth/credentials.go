package auth

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"
)

// CredentialManager manages JWT credentials for tenant connections
type CredentialManager struct {
	store *CredentialStore

	// Configuration
	refreshThreshold time.Duration
	checkInterval    time.Duration
	credentialsDir   string

	// Refresh callback
	onRefreshNeeded func(tenantID string) (*JWTCredential, error)

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	mu sync.RWMutex
}

// CredentialManagerConfig holds configuration
type CredentialManagerConfig struct {
	RefreshThreshold time.Duration // Default: 5 minutes before expiry
	CheckInterval    time.Duration // Default: 1 minute
	CredentialsDir   string        // Directory for credential files
}

// DefaultCredentialManagerConfig returns default configuration
func DefaultCredentialManagerConfig() CredentialManagerConfig {
	return CredentialManagerConfig{
		RefreshThreshold: 5 * time.Minute,
		CheckInterval:    1 * time.Minute,
		CredentialsDir:   "",
	}
}

// NewCredentialManager creates a new credential manager
func NewCredentialManager(config CredentialManagerConfig) *CredentialManager {
	ctx, cancel := context.WithCancel(context.Background())

	if config.RefreshThreshold == 0 {
		config.RefreshThreshold = 5 * time.Minute
	}
	if config.CheckInterval == 0 {
		config.CheckInterval = 1 * time.Minute
	}

	cm := &CredentialManager{
		store:            NewCredentialStore(),
		refreshThreshold: config.RefreshThreshold,
		checkInterval:    config.CheckInterval,
		credentialsDir:   config.CredentialsDir,
		ctx:              ctx,
		cancel:           cancel,
	}

	// Setup store callbacks
	cm.store.OnCredentialExpiring(cm.handleCredentialExpiring)
	cm.store.OnCredentialExpired(cm.handleCredentialExpired)

	return cm
}

// Start starts the credential manager background processes
func (cm *CredentialManager) Start() {
	cm.wg.Add(1)
	go cm.expiryCheckLoop()
}

// Stop stops the credential manager
func (cm *CredentialManager) Stop() {
	cm.cancel()
	cm.wg.Wait()
}

// expiryCheckLoop periodically checks for expiring credentials
func (cm *CredentialManager) expiryCheckLoop() {
	defer cm.wg.Done()

	ticker := time.NewTicker(cm.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-cm.ctx.Done():
			return
		case <-ticker.C:
			cm.store.CheckExpiring(cm.refreshThreshold)
		}
	}
}

// handleCredentialExpiring is called when credentials are about to expire
func (credentialManager *CredentialManager) handleCredentialExpiring(tenantID string, timeUntilExpiry time.Duration) {

	ctx := context.Background()

	logger.Warn(ctx, "credentials expiring soon",
		logger.String("tenant_id", tenantID),
		logger.Duration("time_until_expiry", timeUntilExpiry),
	)

	credentialManager.mu.RLock()
	refreshFn := credentialManager.onRefreshNeeded
	credentialManager.mu.RUnlock()

	if refreshFn != nil {
		credentialManager.wg.Add(1)
		go func() {
			defer credentialManager.wg.Done()

			newCreds, err := refreshFn(tenantID)

			if err != nil {

				logger.Error(ctx, "failed to refresh credentials",
					logger.String("tenant_id", tenantID),
					logger.Err(err),
				)

				return
			}

			if newCreds != nil {

				_ = credentialManager.store.Store(newCreds)

				logger.Info(ctx, "credentials refreshed successfully",
					logger.String("tenant_id", tenantID),
				)
			}
		}()
	}
}

// handleCredentialExpired is called when credentials have expired
func (credentialManager *CredentialManager) handleCredentialExpired(tenantID string) {

	ctx := context.Background()

	logger.Error(ctx, "credentials have expired",
		logger.String("tenant_id", tenantID),
	)
}

// OnRefreshNeeded sets a callback to refresh credentials
func (cm *CredentialManager) OnRefreshNeeded(cb func(tenantID string) (*JWTCredential, error)) {
	cm.mu.Lock()
	cm.onRefreshNeeded = cb
	cm.mu.Unlock()
}

// Register registers JWT credentials for a tenant
func (cm *CredentialManager) Register(tenantID, jwt, seed string) error {
	cred, err := NewJWTCredential(tenantID, jwt, seed)
	if err != nil {
		return err
	}

	if err := cm.store.Store(cred); err != nil {
		return err
	}

	// Optionally persist to file
	if cm.credentialsDir != "" {
		if err := cm.persistCredentials(cred); err != nil {
			// Non-fatal error - log but don't fail registration
			logger.Warn(context.Background(), "failed to persist credentials to file",
				logger.String("tenant_id", tenantID),
				logger.String("dir", cm.credentialsDir),
				logger.Err(err),
			)
		}
	}

	return nil
}

// RegisterFromPayload registers credentials from a registration payload
func (cm *CredentialManager) RegisterFromPayload(payload RegistrationPayload) error {
	return cm.Register(payload.TenantID, payload.JWT, payload.Seed)
}

// Get retrieves valid credentials for a tenant
func (cm *CredentialManager) Get(tenantID string) (*JWTCredential, error) {
	return cm.store.GetValid(tenantID)
}

// GetJWT retrieves just the JWT token for a tenant
func (cm *CredentialManager) GetJWT(tenantID string) (string, error) {
	cred, err := cm.store.GetValid(tenantID)
	if err != nil {
		return "", err
	}
	return cred.JWT, nil
}

// Remove removes credentials for a tenant
func (cm *CredentialManager) Remove(tenantID string) {
	cm.store.Remove(tenantID)

	// Remove persisted file if exists
	if cm.credentialsDir != "" {
		filePath := cm.credentialFilePath(tenantID)
		_ = os.Remove(filePath)
	}
}

// List returns all tenant IDs with registered credentials
func (cm *CredentialManager) List() []string {
	return cm.store.List()
}

// HasValid checks if a tenant has valid credentials
func (cm *CredentialManager) HasValid(tenantID string) bool {
	cred, ok := cm.store.Get(tenantID)
	if !ok {
		return false
	}
	return cred.IsValid()
}

// GetCredentials returns generic credentials for a tenant
func (cm *CredentialManager) GetCredentials(tenantID string) (*Credentials, error) {
	cred, err := cm.Get(tenantID)
	if err != nil {
		return nil, err
	}
	return cred.ToCredentials(), nil
}

// persistCredentials writes credentials to a file
func (cm *CredentialManager) persistCredentials(cred *JWTCredential) error {
	if cm.credentialsDir == "" {
		return nil
	}

	// Ensure directory exists
	if err := os.MkdirAll(cm.credentialsDir, 0700); err != nil {
		return err
	}

	filePath := cm.credentialFilePath(cred.TenantID)
	content := FormatCredsFile(cred.JWT, cred.Seed)

	return os.WriteFile(filePath, []byte(content), 0600)
}

// LoadFromFile loads credentials from a file
func (cm *CredentialManager) LoadFromFile(tenantID, filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	jwt, seed := ParseCredsFile(string(data))
	if jwt == "" {
		return ErrInvalidJWTFormat
	}

	return cm.Register(tenantID, jwt, seed)
}

// credentialFilePath returns the file path for a tenant's credentials
func (cm *CredentialManager) credentialFilePath(tenantID string) string {
	return filepath.Join(cm.credentialsDir, tenantID+".creds")
}

// CredentialStore manages credentials for multiple tenants
type CredentialStore struct {
	mu          sync.RWMutex
	credentials map[string]*JWTCredential

	// Callbacks
	onCredentialExpiring func(tenantID string, timeUntilExpiry time.Duration)
	onCredentialExpired  func(tenantID string)
}

// NewCredentialStore creates a new credential store
func NewCredentialStore() *CredentialStore {
	return &CredentialStore{
		credentials: make(map[string]*JWTCredential),
	}
}

// Store stores credentials for a tenant
func (s *CredentialStore) Store(cred *JWTCredential) error {
	if cred == nil || cred.TenantID == "" {
		return ErrMissingJWT
	}

	s.mu.Lock()
	s.credentials[cred.TenantID] = cred
	s.mu.Unlock()

	return nil
}

// Get retrieves credentials for a tenant
func (s *CredentialStore) Get(tenantID string) (*JWTCredential, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cred, ok := s.credentials[tenantID]
	return cred, ok
}

// GetValid retrieves valid (non-expired) credentials
func (s *CredentialStore) GetValid(tenantID string) (*JWTCredential, error) {
	cred, ok := s.Get(tenantID)
	if !ok {
		return nil, ErrMissingJWT
	}

	if cred.IsExpired() {
		return nil, ErrJWTExpired
	}

	return cred, nil
}

// Remove removes credentials for a tenant
func (s *CredentialStore) Remove(tenantID string) {
	s.mu.Lock()
	delete(s.credentials, tenantID)
	s.mu.Unlock()
}

// List returns all tenant IDs
func (s *CredentialStore) List() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ids := make([]string, 0, len(s.credentials))
	for id := range s.credentials {
		ids = append(ids, id)
	}
	return ids
}

// OnCredentialExpiring sets a callback for expiring credentials
func (s *CredentialStore) OnCredentialExpiring(cb func(tenantID string, timeUntilExpiry time.Duration)) {
	s.onCredentialExpiring = cb
}

// OnCredentialExpired sets a callback for expired credentials
func (s *CredentialStore) OnCredentialExpired(cb func(tenantID string)) {
	s.onCredentialExpired = cb
}

// CheckExpiring checks all credentials and triggers callbacks
func (s *CredentialStore) CheckExpiring(threshold time.Duration) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for tenantID, cred := range s.credentials {
		if cred.IsExpired() {
			if s.onCredentialExpired != nil {
				s.onCredentialExpired(tenantID)
			}
		} else if cred.NeedsRefresh(threshold) {
			if s.onCredentialExpiring != nil {
				s.onCredentialExpiring(tenantID, cred.TimeUntilExpiry())
			}
		}
	}
}

// FormatCredsFile formats credentials in NATS creds file format
func FormatCredsFile(jwt, seed string) string {
	var content string

	if jwt != "" {
		content += "-----BEGIN NATS USER JWT-----\n"
		content += jwt + "\n"
		content += "------END NATS USER JWT------\n"
	}

	if seed != "" {
		content += "\n"
		content += "************************* IMPORTANT *************************\n"
		content += "NKEY Seed printed below can be used to sign and prove identity.\n"
		content += "NKEYs are sensitive and should be treated as secrets.\n"
		content += "\n"
		content += "-----BEGIN USER NKEY SEED-----\n"
		content += seed + "\n"
		content += "------END USER NKEY SEED------\n"
	}

	return content
}

// ParseCredsFile parses a NATS creds file
func ParseCredsFile(content string) (jwt, seed string) {
	lines := splitLines(content)

	inJWT := false
	inSeed := false

	for _, line := range lines {
		line = trimSpace(line)

		if line == "-----BEGIN NATS USER JWT-----" {
			inJWT = true
			continue
		}
		if line == "------END NATS USER JWT------" {
			inJWT = false
			continue
		}
		if line == "-----BEGIN USER NKEY SEED-----" {
			inSeed = true
			continue
		}
		if line == "------END USER NKEY SEED------" {
			inSeed = false
			continue
		}

		if inJWT && line != "" {
			jwt = line
		}
		if inSeed && line != "" {
			seed = line
		}
	}

	return jwt, seed
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}
