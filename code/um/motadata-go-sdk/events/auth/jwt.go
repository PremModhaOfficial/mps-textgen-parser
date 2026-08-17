package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// JWT errors
var (
	ErrInvalidJWT       = errors.New("invalid JWT token")
	ErrJWTExpired       = errors.New("JWT token has expired")
	ErrJWTNotYetValid   = errors.New("JWT token is not yet valid")
	ErrMissingJWT       = errors.New("JWT token is required")
	ErrInvalidJWTFormat = errors.New("invalid JWT format")
	ErrJWTClaimsMissing = errors.New("required JWT claims missing")
)

// JWTClaims represents the claims contained in a JWT token
type JWTClaims struct {
	// Standard JWT claims
	Subject   string `json:"sub"`           // Subject (user/account public key)
	Issuer    string `json:"iss"`           // Issuer (account/operator public key)
	IssuedAt  int64  `json:"iat"`           // Issued at (Unix timestamp)
	ExpiresAt int64  `json:"exp,omitempty"` // Expiration time (Unix timestamp)
	NotBefore int64  `json:"nbf,omitempty"` // Not before (Unix timestamp)
	JWTID     string `json:"jti,omitempty"` // JWT ID

	// NATS-specific claims
	Name string     `json:"name,omitempty"` // Friendly name
	Nats NATSClaims `json:"nats"`           // NATS-specific permissions
}

// NATSClaims contains NATS-specific JWT claims
type NATSClaims struct {
	Type    string `json:"type,omitempty"`    // "user", "account", "operator"
	Version int    `json:"version,omitempty"` // Version

	// Permissions
	Pub  Permission `json:"pub,omitempty"`  // Publish permissions
	Sub  Permission `json:"sub,omitempty"`  // Subscribe permissions
	Resp *Response  `json:"resp,omitempty"` // Response permissions

	// Limits
	Subs    int64 `json:"subs,omitempty"`    // Max subscriptions
	Data    int64 `json:"data,omitempty"`    // Max data in bytes
	Payload int64 `json:"payload,omitempty"` // Max payload size

	// Bearer token (no signature verification)
	BearerToken bool `json:"bearer_token,omitempty"`

	// Tags
	Tags []string `json:"tags,omitempty"`
}

// Permission represents publish or subscribe permissions
type Permission struct {
	Allow []string `json:"allow,omitempty"` // Allowed subjects
	Deny  []string `json:"deny,omitempty"`  // Denied subjects
}

// Response represents response permissions
type Response struct {
	MaxMsgs int           `json:"max,omitempty"` // Max responses
	TTL     time.Duration `json:"ttl,omitempty"` // Response TTL
}

// JWTCredential holds JWT credentials for a tenant
type JWTCredential struct {
	TenantID  string
	JWT       string
	Seed      string // NKey seed (private key) for signing
	CreatedAt time.Time
	ExpiresAt time.Time
	Claims    *JWTClaims
}

// IsExpired checks if the credentials have expired
func (c *JWTCredential) IsExpired() bool {
	if c.ExpiresAt.IsZero() {
		return false // No expiration
	}
	return time.Now().After(c.ExpiresAt)
}

// IsValid checks if the credentials are currently valid
func (c *JWTCredential) IsValid() bool {
	if c.JWT == "" {
		return false
	}
	return !c.IsExpired()
}

// TimeUntilExpiry returns the duration until expiry
func (c *JWTCredential) TimeUntilExpiry() time.Duration {
	if c.ExpiresAt.IsZero() {
		return time.Duration(1<<63 - 1) // Max duration (no expiry)
	}
	return time.Until(c.ExpiresAt)
}

// NeedsRefresh checks if credentials should be refreshed
func (c *JWTCredential) NeedsRefresh(threshold time.Duration) bool {
	return c.TimeUntilExpiry() < threshold
}

// ToCredentials converts to generic credentials
func (c *JWTCredential) ToCredentials() *Credentials {
	return &Credentials{
		Type: TypeJWT,
		JWT:  c.JWT,
		Seed: c.Seed,
	}
}

// ParseJWT parses a JWT token and extracts claims (without signature verification)
// Note: Signature verification is handled by the server
func ParseJWT(token string) (*JWTClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidJWTFormat
	}

	// Decode the payload (second part)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidJWTFormat
	}

	var claims JWTClaims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, ErrInvalidJWTFormat
	}

	return &claims, nil
}

// ValidateClaims validates JWT claims for basic requirements
func ValidateClaims(claims *JWTClaims) error {
	if claims == nil {
		return ErrJWTClaimsMissing
	}

	now := time.Now().Unix()

	// Check expiration
	if claims.ExpiresAt > 0 && now > claims.ExpiresAt {
		return ErrJWTExpired
	}

	// Check not before
	if claims.NotBefore > 0 && now < claims.NotBefore {
		return ErrJWTNotYetValid
	}

	return nil
}

// NewJWTCredential creates a new JWT credential
func NewJWTCredential(tenantID, jwt, seed string) (*JWTCredential, error) {
	claims, err := ParseJWT(jwt)
	if err != nil {
		return nil, err
	}

	if err := ValidateClaims(claims); err != nil {
		return nil, err
	}

	cred := &JWTCredential{
		TenantID:  tenantID,
		JWT:       jwt,
		Seed:      seed,
		CreatedAt: time.Now(),
		Claims:    claims,
	}

	if claims.ExpiresAt > 0 {
		cred.ExpiresAt = time.Unix(claims.ExpiresAt, 0)
	}

	return cred, nil
}
