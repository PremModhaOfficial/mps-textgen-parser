package types

import (
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================================================
// Empty/Blank Checks Tests
// ============================================================================

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", true},
		{" ", false},
		{"a", false},
		{"hello", false},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := IsEmpty(tt.input)
		assertions.Equal(tt.expected, result, "IsEmpty(%q) failed", tt.input)
	}
}

func TestIsBlank(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", true},
		{" ", true},
		{"  ", true},
		{"\t", true},
		{"\n", true},
		{" \t\n ", true},
		{"a", false},
		{" a ", false},
		{"hello", false},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := IsBlank(tt.input)
		assertions.Equal(tt.expected, result, "IsBlank(%q) failed", tt.input)
	}
}

func TestIsNotEmpty(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{" ", true},
		{"a", true},
		{"hello", true},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := IsNotEmpty(tt.input)
		assertions.Equal(tt.expected, result, "IsNotEmpty(%q) failed", tt.input)
	}
}

func TestIsNotBlank(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{" ", false},
		{"  ", false},
		{"a", true},
		{" a ", true},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := IsNotBlank(tt.input)
		assertions.Equal(tt.expected, result, "IsNotBlank(%q) failed", tt.input)
	}
}

// ============================================================================
// Character Type Checks Tests
// ============================================================================

func TestIsASCII(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", true},
		{"hello", true},
		{"Hello123", true},
		{"!@#$%", true},
		{"日本語", false},
		{"hello世界", false},
		{"émoji", false},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := IsASCII(tt.input)
		assertions.Equal(tt.expected, result, "IsASCII(%q) failed", tt.input)
	}
}

func TestIsAlpha(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"abc", true},
		{"ABC", true},
		{"abcDEF", true},
		{"abc123", false},
		{"abc ", false},
		{"abc!", false},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := IsAlpha(tt.input)
		assertions.Equal(tt.expected, result, "IsAlpha(%q) failed", tt.input)
	}
}

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"123", true},
		{"0", true},
		{"123abc", false},
		{"12.34", false},
		{"-123", false},
		{"123 ", false},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := IsNumeric(tt.input)
		assertions.Equal(tt.expected, result, "IsNumeric(%q) failed", tt.input)
	}
}

func TestIsAlphaNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"abc", true},
		{"123", true},
		{"abc123", true},
		{"ABC123def", true},
		{"abc 123", false},
		{"abc-123", false},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := IsAlphaNumeric(tt.input)
		assertions.Equal(tt.expected, result, "IsAlphaNumeric(%q) failed", tt.input)
	}
}

func TestIsLower(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"abc", true},
		{"ABC", false},
		{"abcDEF", false},
		{"abc123", true},
		{"123", false},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := IsLower(tt.input)
		assertions.Equal(tt.expected, result, "IsLower(%q) failed", tt.input)
	}
}

func TestIsUpper(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"ABC", true},
		{"abc", false},
		{"ABCdef", false},
		{"ABC123", true},
		{"123", false},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := IsUpper(tt.input)
		assertions.Equal(tt.expected, result, "IsUpper(%q) failed", tt.input)
	}
}

func TestIsDigit(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"123", true},
		{"0", true},
		{"abc", false},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := IsDigit(tt.input)
		assertions.Equal(tt.expected, result, "IsDigit(%q) failed", tt.input)
	}
}

func TestIsHex(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"0123456789", true},
		{"abcdef", true},
		{"ABCDEF", true},
		{"0123456789abcdefABCDEF", true},
		{"ghij", false},
		{"0x123", false},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := IsHex(tt.input)
		assertions.Equal(tt.expected, result, "IsHex(%q) failed", tt.input)
	}
}

// ============================================================================
// Format Validation Tests
// ============================================================================

func TestIsEmail(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"test@example.com", true},
		{"test.name@example.com", true},
		{"test+tag@example.com", true},
		{"test@sub.example.com", true},
		{"invalid", false},
		{"invalid@", false},
		{"@example.com", false},
		{"test@.com", false},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := IsEmail(tt.input)
		assertions.Equal(tt.expected, result, "IsEmail(%q) failed", tt.input)
	}
}

func TestIsURL(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"https://example.com", true},
		{"https://example.com/path", true},
		{"ftp://example.com", true},
		{"example.com", false},
		{"invalid", false},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := IsURL(tt.input)
		assertions.Equal(tt.expected, result, "IsURL(%q) failed", tt.input)
	}
}

// ============================================================================
// Encoding/Decoding Tests
// ============================================================================

func TestToBytes(t *testing.T) {
	assertions := assert.New(t)
	s := "hello"
	expected := []byte{'h', 'e', 'l', 'l', 'o'}
	result := utils.StringToBytes(s)

	assertions.Equal(expected, result, "ToBytes(%q) failed", s)
}

func TestBase64Encode(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"hello", "aGVsbG8="},
		{"hello world", "aGVsbG8gd29ybGQ="},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := Base64Encode(tt.input)
		assertions.Equal(tt.expected, result, "Base64Encode(%q) failed", tt.input)
	}
}

func TestBase64Decode(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		expectError bool
	}{
		{"", "", false},
		{"aGVsbG8=", "hello", false},
		{"aGVsbG8gd29ybGQ=", "hello world", false},
		{"invalid!!!", "", true},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result, err := Base64Decode(tt.input)
		if tt.expectError {
			assertions.Error(err, "Base64Decode(%q) expected error", tt.input)
		} else {
			assertions.NoError(err, "Base64Decode(%q) unexpected error", tt.input)
			assertions.Equal(tt.expected, result, "Base64Decode(%q) failed", tt.input)
		}
	}
}

func TestURLEncode(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "hello"},
		{"hello world", "hello+world"},
		{"hello=world&foo=bar", "hello%3Dworld%26foo%3Dbar"},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := URLEncode(tt.input)
		assertions.Equal(tt.expected, result, "URLEncode(%q) failed", tt.input)
	}
}

func TestURLDecode(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		expectError bool
	}{
		{"hello", "hello", false},
		{"hello+world", "hello world", false},
		{"hello%20world", "hello world", false},
		{"hello%3Dworld", "hello=world", false},
		{"%invalid", "", true},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result, err := URLDecode(tt.input)
		if tt.expectError {
			assertions.Error(err, "URLDecode(%q) expected error", tt.input)
		} else {
			assertions.NoError(err, "URLDecode(%q) unexpected error", tt.input)
			assertions.Equal(tt.expected, result, "URLDecode(%q) failed", tt.input)
		}
	}
}

// ============================================================================
// Trimming and Whitespace Tests
// ============================================================================

func TestTrim(t *testing.T) {
	tests := []struct {
		input    string
		cutset   string
		expected string
	}{
		{"  hello  ", " ", "hello"},
		{"xxhelloxx", "x", "hello"},
		{"hello", "x", "hello"},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := Trim(tt.input, tt.cutset)
		assertions.Equal(tt.expected, result, "Trim(%q, %q) failed", tt.input, tt.cutset)
	}
}

func TestTrimLeft(t *testing.T) {
	tests := []struct {
		input    string
		cutset   string
		expected string
	}{
		{"  hello  ", " ", "hello  "},
		{"xxhelloxx", "x", "helloxx"},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := TrimLeft(tt.input, tt.cutset)
		assertions.Equal(tt.expected, result, "TrimLeft(%q, %q) failed", tt.input, tt.cutset)
	}
}

func TestTrimRight(t *testing.T) {
	tests := []struct {
		input    string
		cutset   string
		expected string
	}{
		{"  hello  ", " ", "  hello"},
		{"xxhelloxx", "x", "xxhello"},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := TrimRight(tt.input, tt.cutset)
		assertions.Equal(tt.expected, result, "TrimRight(%q, %q) failed", tt.input, tt.cutset)
	}
}

func TestTrimSpace(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  hello  ", "hello"},
		{"\thello\n", "hello"},
		{"hello", "hello"},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := TrimSpace(tt.input)
		assertions.Equal(tt.expected, result, "TrimSpace(%q) failed", tt.input)
	}
}

func TestNormalizeSpace(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  hello  world  ", "hello world"},
		{"hello\t\nworld", "hello world"},
		{"hello   world", "hello world"},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := NormalizeSpace(tt.input)
		assertions.Equal(tt.expected, result, "NormalizeSpace(%q) failed", tt.input)
	}
}

func TestCollapseWhitespace(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  hello  world  ", " hello world "},
		{"hello\t\nworld", "hello world"},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := CollapseWhitespace(tt.input)
		assertions.Equal(tt.expected, result, "CollapseWhitespace(%q) failed", tt.input)
	}
}

func TestRemoveWhitespace(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  hello  world  ", "helloworld"},
		{"hello\t\nworld", "helloworld"},
		{"hello", "hello"},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := RemoveWhitespace(tt.input)
		assertions.Equal(tt.expected, result, "RemoveWhitespace(%q) failed", tt.input)
	}
}

// ============================================================================
// Case Conversion Tests
// ============================================================================

func TestToLower(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HELLO", "hello"},
		{"Hello World", "hello world"},
		{"hello", "hello"},
		{"123", "123"},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := ToLower(tt.input)
		assertions.Equal(tt.expected, result, "ToLower(%q) failed", tt.input)
	}
}

func TestToUpper(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "HELLO"},
		{"Hello World", "HELLO WORLD"},
		{"HELLO", "HELLO"},
		{"123", "123"},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := ToUpper(tt.input)
		assertions.Equal(tt.expected, result, "ToUpper(%q) failed", tt.input)
	}
}

func TestToTitle(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello world", "Hello World"},
		{"HELLO WORLD", "HELLO WORLD"},
		{"hello", "Hello"},
	}

	assertions := assert.New(t)
	for _, tt := range tests {
		result := ToTitle(tt.input)
		assertions.Equal(tt.expected, result, "ToTitle(%q) failed", tt.input)
	}
}

// ============================================================================
// Hashing and Checksums Tests
// ============================================================================

func TestHash64(t *testing.T) {
	assertions := assert.New(t)
	// Test that hash produces consistent output
	s := "hello world"
	h1 := Hash64(s)
	h2 := Hash64(s)

	assertions.Equal(h1, h2, "Hash64(%q) produced inconsistent results", s)

	// Test that different strings produce different hashes
	h3 := Hash64("hello world!")
	assertions.NotEqual(h1, h3, "Hash64 produced same result for different inputs")

	// Test that hash returns a non-zero value
	assertions.NotEqual(uint64(0), h1, "Hash64(%q) returned zero", s)
}

func TestHash32(t *testing.T) {
	assertions := assert.New(t)
	s := "hello world"
	h1 := Hash32(s)
	h2 := Hash32(s)

	assertions.Equal(h1, h2, "Hash32(%q) produced inconsistent results", s)

	// Test that hash returns a non-zero value
	assertions.NotEqual(uint32(0), h1, "Hash32(%q) returned zero", s)
}

func TestHash128(t *testing.T) {
	assertions := assert.New(t)
	s := "hello world"
	h1 := Hash128(s)
	h2 := Hash128(s)

	assertions.Equal(h1, h2, "Hash128(%q) produced inconsistent results", s)

	// Test that hash returns a non-zero value
	assertions.False(h1.Low == 0 && h1.High == 0, "Hash128(%q) returned zero", s)
}

func TestChecksum(t *testing.T) {
	assertions := assert.New(t)
	s := "hello world"
	c1 := Checksum(s)
	c2 := Checksum(s)

	assertions.Equal(c1, c2, "Checksum(%q) produced inconsistent results", s)

	// Test that checksum returns a non-zero value
	assertions.NotEqual(uint32(0), c1, "Checksum(%q) returned zero", s)
}

// ============================================================================
// Encryption/Decryption Tests
// ============================================================================

func TestAESEncryptDecrypt(t *testing.T) {
	// Test with AES-128 (16 byte key)
	key16 := "1234567890123456"
	plaintext := "hello world"

	assertions := assert.New(t)

	encrypted, err := AESEncrypt(plaintext, key16)
	assertions.NoError(err, "AESEncrypt failed")

	assertions.NotEqual(plaintext, encrypted, "Encrypted text should not equal plaintext")

	decrypted, err := AESDecrypt(encrypted, key16)
	assertions.NoError(err, "AESDecrypt failed")

	assertions.Equal(plaintext, decrypted, "AESDecrypt failed")

	// Test with AES-256 (32 byte key)
	key32 := "12345678901234567890123456789012"
	encrypted32, err := AESEncrypt(plaintext, key32)
	assertions.NoError(err, "AESEncrypt with 32-byte key failed")

	decrypted32, err := AESDecrypt(encrypted32, key32)
	assertions.NoError(err, "AESDecrypt with 32-byte key failed")

	assertions.Equal(plaintext, decrypted32, "AESDecrypt with 32-byte key failed")
}

func TestAESEncryptInvalidKey(t *testing.T) {
	assertions := assert.New(t)
	_, err := AESEncrypt("hello", "shortkey")
	assertions.Error(err, "AESEncrypt should fail with invalid key length")
}

func TestAESDecryptInvalidKey(t *testing.T) {
	assertions := assert.New(t)
	// First encrypt with valid key
	key := "1234567890123456"
	encrypted, _ := AESEncrypt("hello", key)

	// Try to decrypt with wrong key length
	_, err := AESDecrypt(encrypted, "shortkey")
	assertions.Error(err, "AESDecrypt should fail with invalid key length")

	// Try to decrypt with wrong key (correct length)
	_, err = AESDecrypt(encrypted, "6543210987654321")
	assertions.Error(err, "AESDecrypt should fail with wrong key")
}

func TestAESDecryptInvalidCiphertext(t *testing.T) {
	assertions := assert.New(t)
	key := "1234567890123456"

	// Invalid base64
	_, err := AESDecrypt("not-valid-base64!!!", key)
	assertions.Error(err, "AESDecrypt should fail with invalid base64")

	// Too short ciphertext
	_, err = AESDecrypt("YWJj", key) // "abc" in base64
	assertions.Error(err, "AESDecrypt should fail with ciphertext too short")
}

// ============================================================================
// Password Hashing Tests
// ============================================================================

func TestHashPassword(t *testing.T) {
	assertions := assert.New(t)
	password := "mySecurePassword123"

	hash, err := HashPassword(password)
	assertions.NoError(err, "HashPassword failed")

	assertions.NotEqual(password, hash, "Hash should not equal plaintext password")
	assertions.NotEmpty(hash, "Hash should not be empty")

	// Hash should start with bcrypt identifier
	assertions.Equal(byte('$'), hash[0], "Hash should start with bcrypt identifier")
}

func TestVerifyPassword(t *testing.T) {
	assertions := assert.New(t)
	password := "mySecurePassword123"

	hash, err := HashPassword(password)
	assertions.NoError(err, "HashPassword failed")

	// Correct password should verify
	assertions.True(VerifyPassword(hash, password), "VerifyPassword should return true for correct password")

	// Wrong password should not verify
	assertions.False(VerifyPassword(hash, "wrongPassword"), "VerifyPassword should return false for wrong password")
}

func TestHashPasswordUniqueness(t *testing.T) {
	assertions := assert.New(t)
	password := "samePassword"

	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)

	// Same password should produce different hashes (due to salt)
	assertions.NotEqual(hash1, hash2, "Same password should produce different hashes due to salt")

	// Both hashes should verify correctly
	assertions.True(VerifyPassword(hash1, password), "First hash should verify correctly")
	assertions.True(VerifyPassword(hash2, password), "Second hash should verify correctly")
}
