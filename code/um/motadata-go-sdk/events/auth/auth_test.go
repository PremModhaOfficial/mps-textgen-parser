package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	userCreds := &Credentials{Username: "user"}
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
		Username:        "user",
		Password:        "pass",
		Token:           "token",
		NKeyFile:        "nkey.file",
		CredentialsFile: "creds.file",
		JWT:             "jwt",
		Seed:            "seed",
	}

	creds := config.ToCredentials()

	assertions.Equal(TypeUserPass, creds.Type, "Type should match")
	assertions.Equal("user", creds.Username, "Username should match")
	assertions.Equal("pass", creds.Password, "Password should match")
	assertions.Equal("token", creds.Token, "Token should match")
	assertions.Equal("nkey.file", creds.NKeyFile, "NKeyFile should match")
	assertions.Equal("creds.file", creds.CredentialsFile, "CredentialsFile should match")
	assertions.Equal("jwt", creds.JWT, "JWT should match")
	assertions.Equal("seed", creds.Seed, "Seed should match")
}

/* ========================================================================================================
   CREDENTIAL FACTORY TESTS
   ======================================================================================================== */

func TestUserPass(t *testing.T) {

	assertions := assert.New(t)

	creds := UserPass("user", "pass")

	assertions.Equal(TypeUserPass, creds.Type, "Type should be UserPass")
	assertions.Equal("user", creds.Username, "Username should match")
	assertions.Equal("pass", creds.Password, "Password should match")
}

func TestToken(t *testing.T) {

	assertions := assert.New(t)

	creds := Token("my-token")

	assertions.Equal(TypeToken, creds.Type, "Type should be Token")
	assertions.Equal("my-token", creds.Token, "Token should match")
}

func TestNKey(t *testing.T) {

	assertions := assert.New(t)

	creds := NKey("/path/to/nkey")

	assertions.Equal(TypeNKey, creds.Type, "Type should be NKey")
	assertions.Equal("/path/to/nkey", creds.NKeyFile, "NKeyFile should match")
}

func TestNKeySeed(t *testing.T) {

	assertions := assert.New(t)

	creds := NKeySeed("seed-value")

	assertions.Equal(TypeNKey, creds.Type, "Type should be NKey")
	assertions.Equal("seed-value", creds.NKeySeed, "NKeySeed should match")
}

func TestCredentialsFileAuth(t *testing.T) {

	assertions := assert.New(t)

	creds := CredentialsFileAuth("/path/to/creds")

	assertions.Equal(TypeCredentialsFile, creds.Type, "Type should be CredentialsFile")
	assertions.Equal("/path/to/creds", creds.CredentialsFile, "CredentialsFile should match")
}

func TestJWT(t *testing.T) {

	assertions := assert.New(t)

	creds := JWT("jwt-token", "seed-value")

	assertions.Equal(TypeJWT, creds.Type, "Type should be JWT")
	assertions.Equal("jwt-token", creds.JWT, "JWT should match")
	assertions.Equal("seed-value", creds.Seed, "Seed should match")
}
