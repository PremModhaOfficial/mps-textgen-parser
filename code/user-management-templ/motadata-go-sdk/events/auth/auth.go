package auth

import (
	"context"
)

// Provider defines the contract for authentication providers that supply credentials
// to NATS connections. Implementations manage the full credential lifecycle: initial
// authentication, periodic refresh, and validity checks. The two built-in implementations
// are StaticProvider (fixed credentials) and ManagedProvider (CredentialManager-backed
// credentials with automatic refresh).
type Provider interface {
	// Type returns the authentication method used by this provider (e.g., TypeJWT, TypeToken).
	// Callers use this to determine which NATS connection option to apply.
	Type() Type

	// Authenticate retrieves the current credentials for establishing a connection.
	// Returns an error if credentials are unavailable or the provider is misconfigured.
	Authenticate(ctx context.Context) (*Credentials, error)

	// Refresh obtains fresh credentials, typically by delegating to a credential manager
	// or an external auth server. For static providers, this simply returns the existing
	// credentials unchanged.
	Refresh(ctx context.Context) (*Credentials, error)

	// IsValid reports whether the provider currently holds usable credentials.
	// This is a lightweight check (no network calls) suitable for health probes.
	IsValid() bool
}

// Type represents the authentication method used to connect to NATS.
// Each type maps to a specific NATS connection option (e.g., nats.UserInfo,
// nats.Token, nats.NkeyOptionFromSeed). The zero value TypeNone indicates
// an unauthenticated connection.
type Type int

const (
	// TypeNone indicates no authentication is configured.
	TypeNone Type = iota
	// TypeUserPass uses a username and password pair for NATS basic authentication.
	TypeUserPass
	// TypeToken uses a single bearer token for NATS token authentication.
	TypeToken
	// TypeNKey uses Ed25519-based NKey authentication, either from a key file or an inline seed.
	TypeNKey
	// TypeCredentialsFile uses a NATS credentials file (.creds) that bundles a JWT and NKey seed.
	TypeCredentialsFile
	// TypeJWT uses a JWT token with an optional NKey seed for request signing.
	TypeJWT
)

// String returns a human-readable name for the authentication type.
// Unknown or future type values return "unknown" to ensure safe logging.
func (authType Type) String() string {
	switch authType {
	case TypeNone:
		return "none"
	case TypeUserPass:
		return "userpass"
	case TypeToken:
		return "token"
	case TypeNKey:
		return "nkey"
	case TypeCredentialsFile:
		return "credentials"
	case TypeJWT:
		return "jwt"
	default:
		return "unknown"
	}
}

// Credentials holds the authentication material needed to connect to NATS.
// Only the fields relevant to the specified Type are populated; for example,
// a TypeUserPass credential has Username and Password set while all other fields
// remain empty. This struct is created by factory functions (UserPass, Token, NKey,
// NKeySeed, CredentialsFileAuth, JWT) or by converting from Config/JWTCredential.
type Credentials struct {
	// Type determines which set of fields below is meaningful.
	Type Type

	// Username and Password are used for TypeUserPass authentication.
	Username string
	Password string

	// Token is used for TypeToken bearer-token authentication.
	Token string

	// NKeyFile is the filesystem path to an NKey file (TypeNKey).
	NKeyFile string
	// NKeySeed is the raw Ed25519 NKey seed string (TypeNKey), used as an
	// alternative to NKeyFile when the seed is provided inline.
	NKeySeed string

	// CredentialsFile is the path to a NATS .creds file (TypeCredentialsFile)
	// that bundles both a JWT and an NKey seed.
	CredentialsFile string

	// JWT is the encoded JWT token (TypeJWT).
	JWT string
	// Seed is the NKey seed paired with the JWT for request signing (TypeJWT).
	Seed string
}

// IsEmpty returns true if the Credentials struct is nil or has no meaningful
// authentication data configured. This is used by providers to determine
// whether the credential is usable before attempting a connection.
func (credentials *Credentials) IsEmpty() bool {
	if credentials == nil {
		return true
	}
	return credentials.Type == TypeNone &&
		credentials.Username == "" &&
		credentials.Password == "" &&
		credentials.Token == "" &&
		credentials.NKeyFile == "" &&
		credentials.NKeySeed == "" &&
		credentials.CredentialsFile == "" &&
		credentials.JWT == ""
}

// Config holds authentication configuration typically loaded from a configuration
// file or environment variables. It mirrors the Credentials struct but is intended
// as a deserialization target. Call ToCredentials to convert it into a Credentials
// value suitable for establishing a NATS connection.
type Config struct {
	// Type selects the authentication method.
	Type Type

	// Username and Password for TypeUserPass.
	Username string
	Password string

	// Token for TypeToken.
	Token string

	// NKeyFile path for TypeNKey.
	NKeyFile string

	// CredentialsFile path for TypeCredentialsFile.
	CredentialsFile string

	// JWT token string and NKey Seed for TypeJWT authentication.
	// The JWT may be provided inline or retrieved from a credential manager.
	JWT  string
	Seed string
}

// ToCredentials converts the Config into a Credentials value by copying all
// fields. This enables a clean separation between the configuration layer
// (deserialization, validation) and the runtime credential usage.
func (authConfig *Config) ToCredentials() *Credentials {
	return &Credentials{
		Type:            authConfig.Type,
		Username:        authConfig.Username,
		Password:        authConfig.Password,
		Token:           authConfig.Token,
		NKeyFile:        authConfig.NKeyFile,
		CredentialsFile: authConfig.CredentialsFile,
		JWT:             authConfig.JWT,
		Seed:            authConfig.Seed,
	}
}

// UserPass creates a Credentials value configured for username/password
// authentication against NATS. The returned credentials have Type set to
// TypeUserPass.
func UserPass(username, password string) *Credentials {
	return &Credentials{
		Type:     TypeUserPass,
		Username: username,
		Password: password,
	}
}

// Token creates a Credentials value configured for bearer-token authentication.
// The returned credentials have Type set to TypeToken.
func Token(token string) *Credentials {
	return &Credentials{
		Type:  TypeToken,
		Token: token,
	}
}

// NKey creates a Credentials value for NKey authentication using a key file
// on disk. The file should contain the Ed25519 seed. The returned credentials
// have Type set to TypeNKey with NKeyFile populated.
func NKey(nkeyFile string) *Credentials {
	return &Credentials{
		Type:     TypeNKey,
		NKeyFile: nkeyFile,
	}
}

// NKeySeed creates a Credentials value for NKey authentication using an inline
// seed string (e.g., "SUAM..."). This is an alternative to NKey when the seed
// is available in memory rather than on disk. The returned credentials have
// Type set to TypeNKey with NKeySeed populated.
func NKeySeed(seed string) *Credentials {
	return &Credentials{
		Type:     TypeNKey,
		NKeySeed: seed,
	}
}

// CredentialsFileAuth creates a Credentials value that references a NATS .creds
// file on disk. The creds file bundles a JWT and NKey seed in NATS' standard
// PEM-like format. The returned credentials have Type set to TypeCredentialsFile.
func CredentialsFileAuth(file string) *Credentials {
	return &Credentials{
		Type:            TypeCredentialsFile,
		CredentialsFile: file,
	}
}

// JWT creates a Credentials value for JWT-based authentication. The jwt parameter
// is the encoded token and seed is the NKey seed used for signing requests.
// The returned credentials have Type set to TypeJWT.
func JWT(jwt, seed string) *Credentials {
	return &Credentials{
		Type: TypeJWT,
		JWT:  jwt,
		Seed: seed,
	}
}
