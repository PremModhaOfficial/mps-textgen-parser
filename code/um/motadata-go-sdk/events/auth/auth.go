package auth

import (
	"context"
)

// Provider is the interface for authentication providers
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

// Type represents the authentication method
type Type int

const (
	TypeNone Type = iota
	TypeUserPass
	TypeToken
	TypeNKey
	TypeCredentialsFile
	TypeJWT
)

// String returns a human-readable type name
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

// Credentials holds authentication credentials
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

// IsEmpty returns true if no credentials are configured
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

// Config holds authentication configuration
type Config struct {
	Type Type

	// UserPass
	Username string
	Password string

	// Token
	Token string

	// NKey
	NKeyFile string

	// Credentials file
	CredentialsFile string

	// JWT (inline or from credential manager)
	JWT  string
	Seed string
}

// ToCredentials converts config to credentials
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

// UserPass creates user/password credentials
func UserPass(username, password string) *Credentials {
	return &Credentials{
		Type:     TypeUserPass,
		Username: username,
		Password: password,
	}
}

// Token creates token credentials
func Token(token string) *Credentials {
	return &Credentials{
		Type:  TypeToken,
		Token: token,
	}
}

// NKey creates NKey credentials from a file
func NKey(nkeyFile string) *Credentials {
	return &Credentials{
		Type:     TypeNKey,
		NKeyFile: nkeyFile,
	}
}

// NKeySeed creates NKey credentials from a seed
func NKeySeed(seed string) *Credentials {
	return &Credentials{
		Type:     TypeNKey,
		NKeySeed: seed,
	}
}

// CredentialsFile creates credentials from a file
func CredentialsFileAuth(file string) *Credentials {
	return &Credentials{
		Type:            TypeCredentialsFile,
		CredentialsFile: file,
	}
}

// JWT creates JWT credentials
func JWT(jwt, seed string) *Credentials {
	return &Credentials{
		Type: TypeJWT,
		JWT:  jwt,
		Seed: seed,
	}
}
