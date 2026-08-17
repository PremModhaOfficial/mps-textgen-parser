package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/* ========================================================================================================
   TEST CONSTANTS
   ======================================================================================================== */

var (
	testUser            = "user"
	testPass            = "pass"
	testTokenValue      = "my-token"
	testSeedValue       = "seed-value"
	testNKeyPath        = "/path/to/nkey"
	testCredsFilePath   = "/path/to/creds"
	testNKeyFileName    = "nkey.file"
	testCredsFileName   = "creds.file"
	testJWTTokenLiteral = "jwt-token"
	testSeedLiteral     = "seed"
)

/* ========================================================================================================
   AUTH TYPE TESTS
   ======================================================================================================== */

func TestTypeString(t *testing.T) {

	testCases := []struct {
		authType Type
		expected string
	}{
		{TypeNone, "none"},
		{TypeUserPass, "userpass"},
		{TypeToken, "token"},
		{TypeNKey, "nkey"},
		{TypeCredentialsFile, "credentials"},
		{TypeJWT, "jwt"},
		{Type(99), "unknown"},
	}

	for _, tc := range testCases {

		t.Run(tc.expected, func(t *testing.T) {

			assertions := assert.New(t)
			assertions.Equal(tc.expected, tc.authType.String(), "Type string should match")
		})
	}
}

/* ========================================================================================================
   CREDENTIALS TESTS
   ======================================================================================================== */

func TestCredentialsIsEmpty(t *testing.T) {

	assertions := assert.New(t)

	// Nil credentials
	var nilCreds *Credentials
	assertions.True(nilCreds.IsEmpty(), "Nil credentials should be empty")

	// Empty credentials
	emptyCreds := &Credentials{}
	assertions.True(emptyCreds.IsEmpty(), "Empty credentials should be empty")

	// With type only
	typedCreds := &Credentials{Type: TypeUserPass}
	assertions.False(typedCreds.IsEmpty(), "Credentials with type should not be empty")

	// With username
	userCreds := &Credentials{Username: testUser}
	assertions.False(userCreds.IsEmpty(), "Credentials with username should not be empty")

	// With token
	tokenCreds := &Credentials{Token: "token"}
	assertions.False(tokenCreds.IsEmpty(), "Credentials with token should not be empty")

	// With JWT
	jwtCreds := &Credentials{JWT: "jwt"}
	assertions.False(jwtCreds.IsEmpty(), "Credentials with JWT should not be empty")
}

func TestConfigToCredentials(t *testing.T) {

	assertions := assert.New(t)

	config := &Config{
		Type:            TypeUserPass,
		Username:        testUser,
		Password:        testPass,
		Token:           "token",
		NKeyFile:        testNKeyFileName,
		CredentialsFile: testCredsFileName,
		JWT:             "jwt",
		Seed:            testSeedLiteral,
	}

	creds := config.ToCredentials()

	assertions.Equal(TypeUserPass, creds.Type, "Type should match")
	assertions.Equal(testUser, creds.Username, "Username should match")
	assertions.Equal(testPass, creds.Password, "Password should match")
	assertions.Equal("token", creds.Token, "Token should match")
	assertions.Equal(testNKeyFileName, creds.NKeyFile, "NKeyFile should match")
	assertions.Equal(testCredsFileName, creds.CredentialsFile, "CredentialsFile should match")
	assertions.Equal("jwt", creds.JWT, "JWT should match")
	assertions.Equal(testSeedLiteral, creds.Seed, "Seed should match")
}

/* ========================================================================================================
   CREDENTIAL FACTORY TESTS
   ======================================================================================================== */

func TestUserPass(t *testing.T) {

	assertions := assert.New(t)

	creds := UserPass(testUser, testPass)

	assertions.Equal(TypeUserPass, creds.Type, "Type should be UserPass")
	assertions.Equal(testUser, creds.Username, "Username should match")
	assertions.Equal(testPass, creds.Password, "Password should match")
}

func TestToken(t *testing.T) {

	assertions := assert.New(t)

	creds := Token(testTokenValue)

	assertions.Equal(TypeToken, creds.Type, "Type should be Token")
	assertions.Equal(testTokenValue, creds.Token, "Token should match")
}

func TestNKey(t *testing.T) {

	assertions := assert.New(t)

	creds := NKey(testNKeyPath)

	assertions.Equal(TypeNKey, creds.Type, "Type should be NKey")
	assertions.Equal(testNKeyPath, creds.NKeyFile, "NKeyFile should match")
}

func TestNKeySeed(t *testing.T) {

	assertions := assert.New(t)

	creds := NKeySeed(testSeedValue)

	assertions.Equal(TypeNKey, creds.Type, "Type should be NKey")
	assertions.Equal(testSeedValue, creds.NKeySeed, "NKeySeed should match")
}

func TestCredentialsFileAuth(t *testing.T) {

	assertions := assert.New(t)

	creds := CredentialsFileAuth(testCredsFilePath)

	assertions.Equal(TypeCredentialsFile, creds.Type, "Type should be CredentialsFile")
	assertions.Equal(testCredsFilePath, creds.CredentialsFile, "CredentialsFile should match")
}

func TestJWT(t *testing.T) {

	assertions := assert.New(t)

	creds := JWT(testJWTTokenLiteral, testSeedValue)

	assertions.Equal(TypeJWT, creds.Type, "Type should be JWT")
	assertions.Equal(testJWTTokenLiteral, creds.JWT, "JWT should match")
	assertions.Equal(testSeedValue, creds.Seed, "Seed should match")
}
