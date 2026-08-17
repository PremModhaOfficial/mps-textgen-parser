package auth

import (
	"encoding/base64"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/* ========================================================================================================
   JWT PARSING TESTS
   ======================================================================================================== */

func createTestJWT(claims *JWTClaims) string {

	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))

	claimsBytes, _ := json.Marshal(claims)
	payload := base64.RawURLEncoding.EncodeToString(claimsBytes)

	signature := base64.RawURLEncoding.EncodeToString([]byte("test-signature"))

	return header + "." + payload + "." + signature
}

func TestParseJWT(t *testing.T) {

	assertions := assert.New(t)

	claims := &JWTClaims{
		Subject:   "test-subject",
		Issuer:    "test-issuer",
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
		Name:      "test-user",
	}

	jwt := createTestJWT(claims)

	parsedClaims, err := ParseJWT(jwt)

	assertions.NoError(err, "Should parse valid JWT")
	assertions.NotNil(parsedClaims, "Claims should not be nil")
	assertions.Equal("test-subject", parsedClaims.Subject, "Subject should match")
	assertions.Equal("test-issuer", parsedClaims.Issuer, "Issuer should match")
	assertions.Equal("test-user", parsedClaims.Name, "Name should match")
}

func TestParseJWTInvalidFormat(t *testing.T) {

	testCases := []struct {
		name string
		jwt  string
	}{
		{"empty", ""},
		{"one part", "header"},
		{"two parts", "header.payload"},
		{"invalid base64", "header.###invalid###.signature"},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			assertions := assert.New(t)
			claims, err := ParseJWT(tc.jwt)

			assertions.Error(err, "Should return error for invalid JWT")
			assertions.ErrorIs(err, ErrInvalidJWTFormat, "Should be format error")
			assertions.Nil(claims, "Claims should be nil")
		})
	}
}

/* ========================================================================================================
   JWT CLAIMS VALIDATION TESTS
   ======================================================================================================== */

func TestValidateClaimsNil(t *testing.T) {

	assertions := assert.New(t)

	err := ValidateClaims(nil)

	assertions.ErrorIs(err, ErrJWTClaimsMissing, "Should return missing claims error")
}

func TestValidateClaimsExpired(t *testing.T) {

	assertions := assert.New(t)

	claims := &JWTClaims{
		Subject:   "test",
		ExpiresAt: time.Now().Add(-time.Hour).Unix(), // Expired 1 hour ago
	}

	err := ValidateClaims(claims)

	assertions.ErrorIs(err, ErrJWTExpired, "Should return expired error")
}

func TestValidateClaimsNotYetValid(t *testing.T) {

	assertions := assert.New(t)

	claims := &JWTClaims{
		Subject:   "test",
		NotBefore: time.Now().Add(time.Hour).Unix(), // Valid in 1 hour
	}

	err := ValidateClaims(claims)

	assertions.ErrorIs(err, ErrJWTNotYetValid, "Should return not yet valid error")
}

func TestValidateClaimsValid(t *testing.T) {

	assertions := assert.New(t)

	claims := &JWTClaims{
		Subject:   "test",
		ExpiresAt: time.Now().Add(time.Hour).Unix(), // Expires in 1 hour
		NotBefore: time.Now().Add(-time.Hour).Unix(), // Valid since 1 hour ago
	}

	err := ValidateClaims(claims)

	assertions.NoError(err, "Should not return error for valid claims")
}

func TestValidateClaimsNoExpiration(t *testing.T) {

	assertions := assert.New(t)

	claims := &JWTClaims{
		Subject: "test",
		// No ExpiresAt or NotBefore
	}

	err := ValidateClaims(claims)

	assertions.NoError(err, "Should not return error for claims without expiration")
}

/* ========================================================================================================
   JWT CREDENTIAL TESTS
   ======================================================================================================== */

func TestNewJWTCredential(t *testing.T) {

	assertions := assert.New(t)

	claims := &JWTClaims{
		Subject:   "test-subject",
		Issuer:    "test-issuer",
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
	}

	jwt := createTestJWT(claims)

	cred, err := NewJWTCredential("tenant-123", jwt, "seed-value")

	assertions.NoError(err, "Should create credential")
	assertions.NotNil(cred, "Credential should not be nil")
	assertions.Equal("tenant-123", cred.TenantID, "TenantID should match")
	assertions.Equal(jwt, cred.JWT, "JWT should match")
	assertions.Equal("seed-value", cred.Seed, "Seed should match")
	assertions.NotNil(cred.Claims, "Claims should be parsed")
	assertions.False(cred.ExpiresAt.IsZero(), "ExpiresAt should be set")
}

func TestNewJWTCredentialInvalidJWT(t *testing.T) {

	assertions := assert.New(t)

	cred, err := NewJWTCredential("tenant-123", "invalid", "seed")

	assertions.Error(err, "Should return error for invalid JWT")
	assertions.Nil(cred, "Credential should be nil")
}

func TestNewJWTCredentialExpiredJWT(t *testing.T) {

	assertions := assert.New(t)

	claims := &JWTClaims{
		Subject:   "test",
		ExpiresAt: time.Now().Add(-time.Hour).Unix(), // Expired
	}

	jwt := createTestJWT(claims)

	cred, err := NewJWTCredential("tenant-123", jwt, "seed")

	assertions.Error(err, "Should return error for expired JWT")
	assertions.ErrorIs(err, ErrJWTExpired, "Should be expired error")
	assertions.Nil(cred, "Credential should be nil")
}

/* ========================================================================================================
   JWT CREDENTIAL EXPIRATION TESTS
   ======================================================================================================== */

func TestJWTCredentialIsExpired(t *testing.T) {

	assertions := assert.New(t)

	// Not expired
	validCred := &JWTCredential{
		ExpiresAt: time.Now().Add(time.Hour),
	}
	assertions.False(validCred.IsExpired(), "Should not be expired")

	// Expired
	expiredCred := &JWTCredential{
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	assertions.True(expiredCred.IsExpired(), "Should be expired")

	// No expiration
	noExpiryCred := &JWTCredential{}
	assertions.False(noExpiryCred.IsExpired(), "Should not be expired when no expiry set")
}

func TestJWTCredentialIsValid(t *testing.T) {

	assertions := assert.New(t)

	// Valid
	validCred := &JWTCredential{
		JWT:       "token",
		ExpiresAt: time.Now().Add(time.Hour),
	}
	assertions.True(validCred.IsValid(), "Should be valid")

	// No JWT
	noJWTCred := &JWTCredential{
		ExpiresAt: time.Now().Add(time.Hour),
	}
	assertions.False(noJWTCred.IsValid(), "Should not be valid without JWT")

	// Expired
	expiredCred := &JWTCredential{
		JWT:       "token",
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	assertions.False(expiredCred.IsValid(), "Should not be valid when expired")
}

func TestJWTCredentialTimeUntilExpiry(t *testing.T) {

	assertions := assert.New(t)

	// Has expiry
	cred := &JWTCredential{
		ExpiresAt: time.Now().Add(time.Hour),
	}
	ttl := cred.TimeUntilExpiry()
	assertions.True(ttl > 59*time.Minute && ttl <= time.Hour, "Time until expiry should be approximately 1 hour")

	// No expiry
	noExpiryCred := &JWTCredential{}
	ttl = noExpiryCred.TimeUntilExpiry()
	assertions.True(ttl > 24*365*100*time.Hour, "Time until expiry should be very large when no expiry")
}

func TestJWTCredentialNeedsRefresh(t *testing.T) {

	assertions := assert.New(t)

	// Needs refresh (expiring soon)
	soonCred := &JWTCredential{
		ExpiresAt: time.Now().Add(2 * time.Minute),
	}
	assertions.True(soonCred.NeedsRefresh(5*time.Minute), "Should need refresh when expiring in 2 min with 5 min threshold")

	// Doesn't need refresh
	laterCred := &JWTCredential{
		ExpiresAt: time.Now().Add(time.Hour),
	}
	assertions.False(laterCred.NeedsRefresh(5*time.Minute), "Should not need refresh when expiring in 1 hour")
}

func TestJWTCredentialToCredentials(t *testing.T) {

	assertions := assert.New(t)

	cred := &JWTCredential{
		JWT:  "jwt-token",
		Seed: "seed-value",
	}

	genericCreds := cred.ToCredentials()

	assertions.Equal(TypeJWT, genericCreds.Type, "Type should be JWT")
	assertions.Equal("jwt-token", genericCreds.JWT, "JWT should match")
	assertions.Equal("seed-value", genericCreds.Seed, "Seed should match")
}
