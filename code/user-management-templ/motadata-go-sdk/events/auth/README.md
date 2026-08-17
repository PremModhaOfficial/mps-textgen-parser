# Auth Package

The `auth` package provides authentication and credential management for the events messaging system with support for multiple authentication methods including JWT, NKey, and user/password.

## Overview

This package contains:

- **Credentials**: Authentication credential types and helpers
- **CredentialManager**: Thread-safe credential storage and retrieval
- **JWT Support**: JWT-based authentication with seed signing
- **Provider Interface**: Extensible authentication provider pattern
- **Registration Handler**: HTTP handler for credential registration

## Package Structure

```
auth/
├── auth.go            # Credentials and auth types
├── credentials.go     # CredentialManager implementation
├── jwt.go             # JWT credential handling
├── provider.go        # Auth provider interface
├── doc.go             # Package documentation
├── auth_test.go       # Auth tests
├── credentials_test.go # Credentials tests
├── jwt_test.go        # JWT tests
├── provider_test.go   # Provider tests
└── README.md          # This file
```

## Authentication Types

```go
type Type int

const (
    TypeNone            Type = iota  // No authentication
    TypeUserPass                     // Username/password
    TypeToken                        // Static token
    TypeNKey                         // NATS NKey
    TypeCredentialsFile              // Credentials from file
    TypeJWT                          // JWT with seed (recommended)
)
```

| Type | Use Case | Security Level |
|------|----------|----------------|
| `TypeNone` | Development only | None |
| `TypeUserPass` | Simple setups | Basic |
| `TypeToken` | Static token auth | Basic |
| `TypeNKey` | NATS-native auth | High |
| `TypeCredentialsFile` | File-based creds | High |
| `TypeJWT` | Multi-tenant production | Highest |

## Credentials

### Credentials Structure

```go
type Credentials struct {
    Type Type

    // UserPass authentication
    Username string
    Password string

    // Token authentication
    Token string

    // NKey authentication
    NKeyFile string
    NKeySeed string

    // Credentials file
    CredentialsFile string

    // JWT authentication
    JWT  string
    Seed string
}
```

### Creating Credentials

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"

// Username/Password
creds := auth.UserPass("username", "password")

// Static Token
creds := auth.Token("my-secret-token")

// NKey from file
creds := auth.NKey("/path/to/user.nkey")

// NKey from seed
creds := auth.NKeySeed("SUAM...")

// Credentials file
creds := auth.CredentialsFileAuth("/path/to/user.creds")

// JWT with seed (recommended for production)
creds := auth.JWT(jwtToken, nkeySeed)
```

### Checking Credentials

```go
creds := auth.JWT(jwt, seed)

// Check if credentials are empty
if creds.IsEmpty() {
    log.Fatal("No credentials provided")
}

// Check credential type
switch creds.Type {
case auth.TypeJWT:
    fmt.Println("Using JWT authentication")
case auth.TypeUserPass:
    fmt.Println("Using username/password")
}
```

## Credential Manager

Thread-safe credential storage for multi-tenant applications:

```go
type CredentialManager struct {
    // Thread-safe storage
}

// Methods
func NewCredentialManager(cfg CredentialManagerConfig) *CredentialManager
func (cm *CredentialManager) Register(tenantID, jwt, seed string) error
func (cm *CredentialManager) GetCredentials(tenantID string) (*Credentials, error)
func (cm *CredentialManager) Revoke(tenantID string)
func (cm *CredentialManager) RegisterFromPayload(payload RegistrationPayload) error
func (cm *CredentialManager) Start() error
func (cm *CredentialManager) Stop()
```

### Usage

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"

// Create credential manager
cm := auth.NewCredentialManager(auth.DefaultCredentialManagerConfig())

// Register tenant credentials
err := cm.Register("tenant-1", jwtToken1, seed1)
err := cm.Register("tenant-2", jwtToken2, seed2)

// Get credentials for tenant
creds, err := cm.GetCredentials("tenant-1")
if err != nil {
    if errors.Is(err, auth.ErrCredentialNotFound) {
        // Tenant not registered
    }
}

// Revoke tenant credentials
cm.Revoke("tenant-1")

// Cleanup
cm.Stop()
```

### Configuration

```go
type CredentialManagerConfig struct {
    // Storage backend (memory, redis, etc.)
    StorageType string

    // Credential refresh interval
    RefreshInterval time.Duration

    // Enable automatic refresh
    AutoRefresh bool
}

cfg := auth.DefaultCredentialManagerConfig()
cfg.RefreshInterval = 5 * time.Minute
cfg.AutoRefresh = true

cm := auth.NewCredentialManager(cfg)
```

## JWT Credential

Extended JWT credential with metadata:

```go
type JWTCredential struct {
    TenantID    string
    JWT         string
    Seed        string
    CreatedAt   time.Time
    ExpiresAt   time.Time
    RefreshedAt time.Time
}

// Methods
func (j *JWTCredential) IsExpired() bool
func (j *JWTCredential) NeedsRefresh(threshold time.Duration) bool
func (j *JWTCredential) ToCredentials() *Credentials
```

## Registration Payload

For HTTP-based credential registration:

```go
type RegistrationPayload struct {
    TenantID string `json:"tenant_id"`
    JWT      string `json:"jwt"`
    Seed     string `json:"seed"`
}

// Validation
func (p RegistrationPayload) Validate() error
```

### Registration Handler

HTTP handler for credential registration:

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"

cm := auth.NewCredentialManager(auth.DefaultCredentialManagerConfig())
handler := auth.NewRegistrationHandler(cm)

// Register HTTP endpoints
http.HandleFunc("/credentials/register", handler.HandleRegister)
http.HandleFunc("/credentials/revoke", handler.HandleRevoke)
```

**Register Endpoint:**

```bash
POST /credentials/register
Content-Type: application/json

{
    "tenant_id": "tenant-123",
    "jwt": "eyJ...",
    "seed": "SUAM..."
}
```

**Revoke Endpoint:**

```bash
DELETE /credentials/revoke?tenant_id=tenant-123
```

## Auth Provider Interface

Extensible authentication provider:

```go
type Provider interface {
    // Type returns the authentication type
    Type() Type

    // Authenticate returns credentials for authentication
    Authenticate(ctx context.Context) (*Credentials, error)

    // Refresh refreshes the credentials if needed
    Refresh(ctx context.Context) (*Credentials, error)

    // IsValid checks if current credentials are valid
    IsValid() bool
}
```

### Custom Provider Example

```go
type VaultProvider struct {
    vaultClient *vault.Client
    tenantID    string
}

func (p *VaultProvider) Type() auth.Type {
    return auth.TypeJWT
}

func (p *VaultProvider) Authenticate(ctx context.Context) (*auth.Credentials, error) {
    // Fetch JWT from Vault
    secret, err := p.vaultClient.Read("nats/creds/" + p.tenantID)
    if err != nil {
        return nil, err
    }

    return auth.JWT(
        secret.Data["jwt"].(string),
        secret.Data["seed"].(string),
    ), nil
}

func (p *VaultProvider) Refresh(ctx context.Context) (*auth.Credentials, error) {
    return p.Authenticate(ctx)
}

func (p *VaultProvider) IsValid() bool {
    return p.vaultClient != nil
}
```

## Error Definitions

| Error | Description |
|-------|-------------|
| `ErrInvalidCredential` | Credentials are invalid or empty |
| `ErrCredentialNotFound` | Tenant credentials not registered |
| `ErrCredentialExpired` | JWT has expired |
| `ErrInvalidJWT` | JWT format is invalid |
| `ErrInvalidSeed` | NKey seed is invalid |

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/auth"

creds, err := cm.GetCredentials("tenant-1")
if err != nil {
    switch {
    case errors.Is(err, auth.ErrCredentialNotFound):
        // Register credentials first
    case errors.Is(err, auth.ErrCredentialExpired):
        // Refresh credentials
    case errors.Is(err, auth.ErrInvalidCredential):
        // Invalid credential format
    }
}
```

## JWT Authentication Flow

```
┌─────────────┐     ┌───────────────────┐     ┌─────────────┐
│   Client    │     │ CredentialManager │     │ NATS Server │
└─────────────┘     └───────────────────┘     └─────────────┘
       │                      │                       │
       │ Register(tenantID,   │                       │
       │   jwt, seed)         │                       │
       │─────────────────────>│                       │
       │                      │                       │
       │ Connect(tenantID)    │                       │
       │─────────────────────>│                       │
       │                      │                       │
       │                      │ GetCredentials()      │
       │                      │───────────┐           │
       │                      │           │           │
       │                      │<──────────┘           │
       │                      │                       │
       │      *Credentials    │ Connect with JWT     │
       │<─────────────────────│──────────────────────>│
       │                      │                       │
       │                      │     Authenticated     │
       │                      │<──────────────────────│
       │                      │                       │
```

## Best Practices

### 1. Use JWT for Production

```go
// Recommended for multi-tenant production
creds := auth.JWT(jwtToken, nkeySeed)
```

### 2. Secure Credential Storage

```go
// Use CredentialManager for thread-safe access
cm := auth.NewCredentialManager(auth.DefaultCredentialManagerConfig())

// Register credentials securely
err := cm.Register(tenantID, jwt, seed)

// Never log credentials
// BAD: log.Printf("JWT: %s", jwt)
// GOOD: log.Printf("Registered credentials for tenant: %s", tenantID)
```

### 3. Handle Credential Expiration

```go
creds, err := cm.GetCredentials(tenantID)
if errors.Is(err, auth.ErrCredentialExpired) {
    // Trigger credential refresh
    newJWT, newSeed, err := refreshFromAuthServer(tenantID)
    if err == nil {
        cm.Register(tenantID, newJWT, newSeed)
    }
}
```

### 4. Validate Early

```go
payload := auth.RegistrationPayload{
    TenantID: tenantID,
    JWT:      jwt,
    Seed:     seed,
}

// Validate before processing
if err := payload.Validate(); err != nil {
    return fmt.Errorf("invalid registration: %w", err)
}
```

### 5. Clean Up on Shutdown

```go
func main() {
    cm := auth.NewCredentialManager(cfg)
    defer cm.Stop() // Stop refresh loops, cleanup

    // Application code...
}
```

## Testing

```bash
# Run tests
go test ./events/auth/...

# Run with coverage
go test -cover ./events/auth/...

# Run specific tests
go test -run TestJWT ./events/auth/...
```

## Step-by-Step Usage Guide

### Step 1: Choose Authentication Method

```go
// Development: No auth
creds := auth.Credentials{Type: auth.TypeNone}

// Simple setup: Token
creds = auth.Token("my-secret-token")

// Production (recommended): JWT with seed
creds = auth.JWT(jwtToken, nkeySeed)
```

### Step 2: Create Credential Manager (Multi-Tenant)

```go
cm := auth.NewCredentialManager(auth.DefaultCredentialManagerConfig())
cm.Start()
defer cm.Stop()
```

### Step 3: Register Tenant Credentials

```go
cm.Register("tenant-1", jwt1, seed1)
cm.Register("tenant-2", jwt2, seed2)

// Or from registration payload
payload := auth.RegistrationPayload{
    TenantID: "tenant-3",
    JWT:      jwt3,
    Seed:     seed3,
}
cm.RegisterFromPayload(payload)
```

### Step 4: Retrieve Credentials When Connecting

```go
creds, err := cm.GetCredentials("tenant-1")
if err != nil {
    if errors.Is(err, auth.ErrCredentialNotFound) {
        // Register credentials first
    }
}
```

### Step 5: Handle Credential Refresh

```go
cm.OnRefreshNeeded(func(tenantID string) {
    // Fetch new JWT from auth server
    newJWT, newSeed := refreshFromAuthServer(tenantID)
    cm.Register(tenantID, newJWT, newSeed)
})
```

### Step 6: Set Up HTTP Registration (Optional)

```go
handler := auth.NewRegistrationHandler(cm)
http.HandleFunc("/credentials/register", handler.HandleRegister)
http.HandleFunc("/credentials/revoke", handler.HandleRevoke)
```

## Usage in Other Packages

This package is used by:
- `events/tenant` - Tenant connection authentication
- `events/connection.go` - NATS connection authentication
- `events` (main) - Re-exports for convenience
