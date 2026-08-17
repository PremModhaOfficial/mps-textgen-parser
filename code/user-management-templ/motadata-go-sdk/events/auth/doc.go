// Package auth provides authentication and credential management for event-driven messaging.
//
// This package supports multiple authentication methods for connecting to messaging systems
// like NATS, including:
//   - Username/Password authentication
//   - Token-based authentication
//   - NKey authentication (Ed25519 keys)
//   - Credentials file authentication
//   - JWT-based authentication with automatic refresh
//
// # Authentication Types
//
// The package supports the following authentication types:
//
//	TypeNone            - No authentication
//	TypeUserPass        - Username and password
//	TypeToken           - Bearer token
//	TypeNKey            - NKey seed or file
//	TypeCredentialsFile - NATS credentials file
//	TypeJWT             - JWT token with optional NKey seed
//
// # Creating Credentials
//
// Use the helper functions to create credentials:
//
//	// Username/password
//	creds := auth.UserPass("user", "password")
//
//	// Token
//	creds := auth.Token("my-secret-token")
//
//	// NKey from file
//	creds := auth.NKey("/path/to/nkey.nk")
//
//	// NKey from seed
//	creds := auth.NKeySeed("SUAM...")
//
//	// JWT with seed for signing
//	creds := auth.JWT(jwtToken, nkeySeed)
//
// # Credential Manager
//
// The CredentialManager handles JWT credentials for multi-tenant scenarios:
//
//	// Create manager
//	manager := auth.NewCredentialManager(auth.DefaultCredentialManagerConfig())
//	manager.Start()
//	defer manager.Stop()
//
//	// Register tenant credentials
//	err := manager.Register("tenant-1", jwtToken, nkeySeed)
//
//	// Set refresh callback for expiring credentials
//	manager.OnRefreshNeeded(func(tenantID string) (*auth.JWTCredential, error) {
//	    // Fetch new credentials from your auth server
//	    return refreshCredentialsFromServer(tenantID)
//	})
//
//	// Get credentials for a tenant
//	creds, err := manager.GetCredentials("tenant-1")
//
// # JWT Parsing and Validation
//
// The package can parse and validate JWT tokens without signature verification
// (signature verification is done by the NATS server):
//
//	claims, err := auth.ParseJWT(jwtToken)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Check permissions
//	fmt.Println("Publish permissions:", claims.Nats.Pub.Allow)
//	fmt.Println("Subscribe permissions:", claims.Nats.Sub.Allow)
//
// # Registration Handling
//
// For services that receive credentials from a registration server:
//
//	handler := auth.NewRegistrationHandler(credManager)
//
//	handler.OnRegistered(func(tenantID string) {
//	    log.Printf("Tenant %s registered", tenantID)
//	})
//
//	// Process registration payload
//	err := handler.HandleRegistrationJSON(payloadBytes)
//
// # Provider Interface
//
// Implement the Provider interface for custom authentication:
//
//	type Provider interface {
//	    Type() Type
//	    Authenticate(ctx context.Context) (*Credentials, error)
//	    Refresh(ctx context.Context) (*Credentials, error)
//	    IsValid() bool
//	}
//
// Built-in providers include:
//   - StaticProvider: Returns static credentials
//   - ManagedProvider: Returns credentials from CredentialManager
package auth
