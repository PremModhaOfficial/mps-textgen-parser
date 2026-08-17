package auth

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/otel/logger"
)

// CredentialManager manages JWT credentials for multi-tenant NATS connections.
//
// It owns the full credential lifecycle: registration, storage, periodic expiry
// checks, proactive refresh, optional file persistence, and removal. A background
// goroutine (started by Start) periodically scans all stored credentials and
// triggers the registered onRefreshNeeded callback when a credential is within
// the configured refresh threshold of expiry. Expired credentials trigger a
// separate onCredentialExpired callback for alerting.
//
// Thread-safety: all public methods are safe for concurrent use. The manager
// protects its callback with mu (RWMutex), and delegates credential storage
// to the thread-safe CredentialStore. The background loop is coordinated via
// ctx/cancel and wg to ensure clean shutdown.
type CredentialManager struct {
	// store is the thread-safe, in-memory credential store keyed by tenant ID.
	store *CredentialStore

	// refreshThreshold is how far before expiry to begin refresh attempts.
	refreshThreshold time.Duration
	// checkInterval is the polling period for the background expiry-check loop.
	checkInterval time.Duration
	// credentialsDir, when non-empty, enables file-based persistence of credentials.
	credentialsDir string

	// onRefreshNeeded is the user-supplied callback invoked asynchronously when
	// a credential is about to expire. It should fetch new credentials from an
	// external auth server. Protected by mu.
	onRefreshNeeded func(tenantID string) (*JWTCredential, error)

	// done is closed when cancel is called, signalling background goroutines to stop.
	done <-chan struct{}
	// cancel stops the background expiry-check goroutine.
	cancel context.CancelFunc
	// wg tracks all background goroutines (expiry loop + async refresh calls)
	// so Stop can wait for graceful shutdown.
	wg sync.WaitGroup

	// mu protects the onRefreshNeeded callback from concurrent read/write.
	mu sync.RWMutex
}

// CredentialManagerConfig holds the tuning parameters for CredentialManager.
// Zero values for durations are replaced with sensible defaults by
// NewCredentialManager.
type CredentialManagerConfig struct {
	// RefreshThreshold determines how far before expiry a credential is
	// considered "expiring" and the refresh callback is invoked.
	// Default: 5 minutes.
	RefreshThreshold time.Duration
	// CheckInterval is how often the background loop scans for expiring
	// credentials. Default: 1 minute.
	CheckInterval time.Duration
	// CredentialsDir is the filesystem directory where credential files are
	// persisted. When empty, file persistence is disabled.
	CredentialsDir string
}

// DefaultCredentialManagerConfig returns a CredentialManagerConfig with
// production-ready defaults: 5-minute refresh threshold, 1-minute check
// interval, and no file persistence.
func DefaultCredentialManagerConfig() CredentialManagerConfig {
	return CredentialManagerConfig{
		RefreshThreshold: 5 * time.Minute,
		CheckInterval:    1 * time.Minute,
		CredentialsDir:   "",
	}
}

// NewCredentialManager creates a new CredentialManager with the given config.
// Zero-value durations in config are replaced with defaults. The manager is
// not started until Start is called; call Stop to release background resources.
// It wires the internal CredentialStore's expiring/expired callbacks to the
// manager's own handlers so that refresh and alerting happen automatically.
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
		done:             ctx.Done(),
		cancel:           cancel,
	}

	// Setup store callbacks
	cm.store.OnCredentialExpiring(cm.handleCredentialExpiring)
	cm.store.OnCredentialExpired(cm.handleCredentialExpired)

	return cm
}

// Start launches the background goroutine that periodically checks all stored
// credentials for impending or past expiry. Must be called after NewCredentialManager
// and before any credentials are expected to auto-refresh.
func (cm *CredentialManager) Start() {
	cm.wg.Add(1)
	go cm.expiryCheckLoop()
}

// Stop cancels the background expiry-check goroutine and waits for all
// in-flight refresh goroutines to complete, ensuring a clean shutdown.
func (cm *CredentialManager) Stop() {
	cm.cancel()
	cm.wg.Wait()
}

// expiryCheckLoop runs in a background goroutine, ticking at checkInterval and
// delegating to the store's CheckExpiring to fire expiring/expired callbacks.
// It exits cleanly when the manager's context is cancelled.
func (cm *CredentialManager) expiryCheckLoop() {
	defer cm.wg.Done()

	ticker := time.NewTicker(cm.checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-cm.done:
			return
		case <-ticker.C:
			cm.store.CheckExpiring(cm.refreshThreshold)
		}
	}
}

// handleCredentialExpiring is the callback invoked by CredentialStore when a
// tenant's credential is within the refresh threshold. It reads the user-supplied
// onRefreshNeeded callback under a read lock and, if set, spawns a goroutine to
// fetch new credentials and store them. Refresh failures are logged but do not
// remove the existing (still-valid) credential, allowing retry on the next tick.
func (cm *CredentialManager) handleCredentialExpiring(tenantID string, timeUntilExpiry time.Duration) {
	ctx := context.Background()

	logger.Warn(ctx, "credentials expiring soon",
		logger.String("tenant_id", tenantID),
		logger.Duration("time_until_expiry", timeUntilExpiry),
	)

	cm.mu.RLock()
	refreshFn := cm.onRefreshNeeded
	cm.mu.RUnlock()

	if refreshFn != nil {
		cm.wg.Add(1)
		go func() {
			defer cm.wg.Done()

			newCreds, err := refreshFn(tenantID)
			if err != nil {
				logger.Error(ctx, "failed to refresh credentials",
					logger.String("tenant_id", tenantID),
					logger.Err(err),
				)
				return
			}

			if newCreds != nil {
				if storeErr := cm.store.Store(newCreds); storeErr != nil {
					logger.Error(ctx, "failed to store refreshed credentials",
						logger.String("tenant_id", tenantID),
						logger.Err(storeErr),
					)
					return
				}

				logger.Info(ctx, "credentials refreshed successfully",
					logger.String("tenant_id", tenantID),
				)
			}
		}()
	}
}

// handleCredentialExpired is called when credentials have expired
func (cm *CredentialManager) handleCredentialExpired(tenantID string) {
	logger.Error(context.Background(), "credentials have expired",
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
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			logger.Warn(context.Background(), "failed to remove credential file",
				logger.String("tenant_id", tenantID),
				logger.String("path", filePath),
				logger.Err(err),
			)
		}
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
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onCredentialExpiring = cb
}

// OnCredentialExpired sets a callback for expired credentials
func (s *CredentialStore) OnCredentialExpired(cb func(tenantID string)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onCredentialExpired = cb
}

// CheckExpiring checks all credentials and triggers callbacks
func (s *CredentialStore) CheckExpiring(threshold time.Duration) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for tenantID, cred := range s.credentials {
		s.checkCredentialStatus(tenantID, cred, threshold)
	}
}

// checkCredentialStatus checks a single credential and fires the appropriate callback.
func (s *CredentialStore) checkCredentialStatus(tenantID string, cred *JWTCredential, threshold time.Duration) {
	if cred.IsExpired() && s.onCredentialExpired != nil {
		s.onCredentialExpired(tenantID)
		return
	}

	if cred.NeedsRefresh(threshold) && s.onCredentialExpiring != nil {
		s.onCredentialExpiring(tenantID, cred.TimeUntilExpiry())
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

// credsSection tracks which section of a creds file is being parsed.
type credsSection int

const (
	credsSectionNone credsSection = iota
	credsSectionJWT
	credsSectionSeed
)

// ParseCredsFile parses a NATS creds file
func ParseCredsFile(content string) (jwt, seed string) {
	lines := splitLines(content)
	section := credsSectionNone

	for _, line := range lines {
		line = trimSpace(line)
		section, jwt, seed = parseCredsLine(line, section, jwt, seed)
	}

	return jwt, seed
}

// parseCredsLine processes a single line of a creds file, returning updated state.
func parseCredsLine(line string, section credsSection, jwt, seed string) (credsSection, string, string) {
	switch line {
	case "-----BEGIN NATS USER JWT-----":
		return credsSectionJWT, jwt, seed
	case "------END NATS USER JWT------":
		return credsSectionNone, jwt, seed
	case "-----BEGIN USER NKEY SEED-----":
		return credsSectionSeed, jwt, seed
	case "------END USER NKEY SEED------":
		return credsSectionNone, jwt, seed
	}

	if line == "" {
		return section, jwt, seed
	}

	switch section {
	case credsSectionJWT:
		jwt = line
	case credsSectionSeed:
		seed = line
	}

	return section, jwt, seed
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
