package types

import (
	"testing"
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

	for _, tt := range tests {
		result := IsEmpty(tt.input)
		if result != tt.expected {
			t.Errorf("IsEmpty(%q) = %v, want %v", tt.input, result, tt.expected)
		}
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

	for _, tt := range tests {
		result := IsBlank(tt.input)
		if result != tt.expected {
			t.Errorf("IsBlank(%q) = %v, want %v", tt.input, result, tt.expected)
		}
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

	for _, tt := range tests {
		result := IsNotEmpty(tt.input)
		if result != tt.expected {
			t.Errorf("IsNotEmpty(%q) = %v, want %v", tt.input, result, tt.expected)
		}
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

	for _, tt := range tests {
		result := IsNotBlank(tt.input)
		if result != tt.expected {
			t.Errorf("IsNotBlank(%q) = %v, want %v", tt.input, result, tt.expected)
		}
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

	for _, tt := range tests {
		result := IsASCII(tt.input)
		if result != tt.expected {
			t.Errorf("IsASCII(%q) = %v, want %v", tt.input, result, tt.expected)
		}
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

	for _, tt := range tests {
		result := IsAlpha(tt.input)
		if result != tt.expected {
			t.Errorf("IsAlpha(%q) = %v, want %v", tt.input, result, tt.expected)
		}
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

	for _, tt := range tests {
		result := IsNumeric(tt.input)
		if result != tt.expected {
			t.Errorf("IsNumeric(%q) = %v, want %v", tt.input, result, tt.expected)
		}
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

	for _, tt := range tests {
		result := IsAlphaNumeric(tt.input)
		if result != tt.expected {
			t.Errorf("IsAlphaNumeric(%q) = %v, want %v", tt.input, result, tt.expected)
		}
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

	for _, tt := range tests {
		result := IsLower(tt.input)
		if result != tt.expected {
			t.Errorf("IsLower(%q) = %v, want %v", tt.input, result, tt.expected)
		}
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

	for _, tt := range tests {
		result := IsUpper(tt.input)
		if result != tt.expected {
			t.Errorf("IsUpper(%q) = %v, want %v", tt.input, result, tt.expected)
		}
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

	for _, tt := range tests {
		result := IsDigit(tt.input)
		if result != tt.expected {
			t.Errorf("IsDigit(%q) = %v, want %v", tt.input, result, tt.expected)
		}
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

	for _, tt := range tests {
		result := IsHex(tt.input)
		if result != tt.expected {
			t.Errorf("IsHex(%q) = %v, want %v", tt.input, result, tt.expected)
		}
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

	for _, tt := range tests {
		result := IsEmail(tt.input)
		if result != tt.expected {
			t.Errorf("IsEmail(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

func TestIsURL(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", false},
		{"https://example.com", true},
		{"http://example.com", true},
		{"https://example.com/path", true},
		{"ftp://example.com", true},
		{"example.com", false},
		{"invalid", false},
	}

	for _, tt := range tests {
		result := IsURL(tt.input)
		if result != tt.expected {
			t.Errorf("IsURL(%q) = %v, want %v", tt.input, result, tt.expected)
		}
	}
}

// ============================================================================
// Encoding/Decoding Tests
// ============================================================================

func TestToBytes(t *testing.T) {
	s := "hello"
	expected := []byte{'h', 'e', 'l', 'l', 'o'}
	result := ToBytes(s)

	if len(result) != len(expected) {
		t.Errorf("ToBytes(%q) length = %d, want %d", s, len(result), len(expected))
	}
	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("ToBytes(%q)[%d] = %d, want %d", s, i, result[i], expected[i])
		}
	}
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

	for _, tt := range tests {
		result := Base64Encode(tt.input)
		if result != tt.expected {
			t.Errorf("Base64Encode(%q) = %q, want %q", tt.input, result, tt.expected)
		}
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

	for _, tt := range tests {
		result, err := Base64Decode(tt.input)
		if tt.expectError && err == nil {
			t.Errorf("Base64Decode(%q) expected error, got nil", tt.input)
		}
		if !tt.expectError && err != nil {
			t.Errorf("Base64Decode(%q) unexpected error: %v", tt.input, err)
		}
		if result != tt.expected {
			t.Errorf("Base64Decode(%q) = %q, want %q", tt.input, result, tt.expected)
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

	for _, tt := range tests {
		result := URLEncode(tt.input)
		if result != tt.expected {
			t.Errorf("URLEncode(%q) = %q, want %q", tt.input, result, tt.expected)
		}
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

	for _, tt := range tests {
		result, err := URLDecode(tt.input)
		if tt.expectError && err == nil {
			t.Errorf("URLDecode(%q) expected error, got nil", tt.input)
		}
		if !tt.expectError && err != nil {
			t.Errorf("URLDecode(%q) unexpected error: %v", tt.input, err)
		}
		if !tt.expectError && result != tt.expected {
			t.Errorf("URLDecode(%q) = %q, want %q", tt.input, result, tt.expected)
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

	for _, tt := range tests {
		result := Trim(tt.input, tt.cutset)
		if result != tt.expected {
			t.Errorf("Trim(%q, %q) = %q, want %q", tt.input, tt.cutset, result, tt.expected)
		}
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

	for _, tt := range tests {
		result := TrimLeft(tt.input, tt.cutset)
		if result != tt.expected {
			t.Errorf("TrimLeft(%q, %q) = %q, want %q", tt.input, tt.cutset, result, tt.expected)
		}
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

	for _, tt := range tests {
		result := TrimRight(tt.input, tt.cutset)
		if result != tt.expected {
			t.Errorf("TrimRight(%q, %q) = %q, want %q", tt.input, tt.cutset, result, tt.expected)
		}
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

	for _, tt := range tests {
		result := TrimSpace(tt.input)
		if result != tt.expected {
			t.Errorf("TrimSpace(%q) = %q, want %q", tt.input, result, tt.expected)
		}
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

	for _, tt := range tests {
		result := NormalizeSpace(tt.input)
		if result != tt.expected {
			t.Errorf("NormalizeSpace(%q) = %q, want %q", tt.input, result, tt.expected)
		}
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

	for _, tt := range tests {
		result := CollapseWhitespace(tt.input)
		if result != tt.expected {
			t.Errorf("CollapseWhitespace(%q) = %q, want %q", tt.input, result, tt.expected)
		}
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

	for _, tt := range tests {
		result := RemoveWhitespace(tt.input)
		if result != tt.expected {
			t.Errorf("RemoveWhitespace(%q) = %q, want %q", tt.input, result, tt.expected)
		}
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

	for _, tt := range tests {
		result := ToLower(tt.input)
		if result != tt.expected {
			t.Errorf("ToLower(%q) = %q, want %q", tt.input, result, tt.expected)
		}
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

	for _, tt := range tests {
		result := ToUpper(tt.input)
		if result != tt.expected {
			t.Errorf("ToUpper(%q) = %q, want %q", tt.input, result, tt.expected)
		}
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

	for _, tt := range tests {
		result := ToTitle(tt.input)
		if result != tt.expected {
			t.Errorf("ToTitle(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

// ============================================================================
// Hashing and Checksums Tests
// ============================================================================

func TestHash64(t *testing.T) {
	// Test that hash produces consistent output
	s := "hello world"
	h1 := Hash64(s)
	h2 := Hash64(s)

	if h1 != h2 {
		t.Errorf("Hash64(%q) produced inconsistent results: %d != %d", s, h1, h2)
	}

	// Test that different strings produce different hashes
	h3 := Hash64("hello world!")
	if h1 == h3 {
		t.Errorf("Hash64 produced same result for different inputs")
	}

	// Test that hash returns a non-zero value
	if h1 == 0 {
		t.Errorf("Hash64(%q) returned zero", s)
	}
}

func TestHash32(t *testing.T) {
	s := "hello world"
	h1 := Hash32(s)
	h2 := Hash32(s)

	if h1 != h2 {
		t.Errorf("Hash32(%q) produced inconsistent results: %d != %d", s, h1, h2)
	}

	// Test that hash returns a non-zero value
	if h1 == 0 {
		t.Errorf("Hash32(%q) returned zero", s)
	}
}

func TestHash128(t *testing.T) {
	s := "hello world"
	h1 := Hash128(s)
	h2 := Hash128(s)

	if h1 != h2 {
		t.Errorf("Hash128(%q) produced inconsistent results", s)
	}

	// Test that hash returns a non-zero value
	if h1.Low == 0 && h1.High == 0 {
		t.Errorf("Hash128(%q) returned zero", s)
	}
}

func TestChecksum(t *testing.T) {
	s := "hello world"
	c1 := Checksum(s)
	c2 := Checksum(s)

	if c1 != c2 {
		t.Errorf("Checksum(%q) produced inconsistent results: %d != %d", s, c1, c2)
	}

	// Test that checksum returns a non-zero value
	if c1 == 0 {
		t.Errorf("Checksum(%q) returned zero", s)
	}
}

// ============================================================================
// Encryption/Decryption Tests
// ============================================================================

func TestAESEncryptDecrypt(t *testing.T) {
	// Test with AES-128 (16 byte key)
	key16 := "1234567890123456"
	plaintext := "hello world"

	encrypted, err := AESEncrypt(plaintext, key16)
	if err != nil {
		t.Fatalf("AESEncrypt failed: %v", err)
	}

	if encrypted == plaintext {
		t.Error("Encrypted text should not equal plaintext")
	}

	decrypted, err := AESDecrypt(encrypted, key16)
	if err != nil {
		t.Fatalf("AESDecrypt failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("AESDecrypt returned %q, want %q", decrypted, plaintext)
	}

	// Test with AES-256 (32 byte key)
	key32 := "12345678901234567890123456789012"
	encrypted32, err := AESEncrypt(plaintext, key32)
	if err != nil {
		t.Fatalf("AESEncrypt with 32-byte key failed: %v", err)
	}

	decrypted32, err := AESDecrypt(encrypted32, key32)
	if err != nil {
		t.Fatalf("AESDecrypt with 32-byte key failed: %v", err)
	}

	if decrypted32 != plaintext {
		t.Errorf("AESDecrypt returned %q, want %q", decrypted32, plaintext)
	}
}

func TestAESEncryptInvalidKey(t *testing.T) {
	_, err := AESEncrypt("hello", "shortkey")
	if err == nil {
		t.Error("AESEncrypt should fail with invalid key length")
	}
}

func TestAESDecryptInvalidKey(t *testing.T) {
	// First encrypt with valid key
	key := "1234567890123456"
	encrypted, _ := AESEncrypt("hello", key)

	// Try to decrypt with wrong key length
	_, err := AESDecrypt(encrypted, "shortkey")
	if err == nil {
		t.Error("AESDecrypt should fail with invalid key length")
	}

	// Try to decrypt with wrong key (correct length)
	_, err = AESDecrypt(encrypted, "6543210987654321")
	if err == nil {
		t.Error("AESDecrypt should fail with wrong key")
	}
}

func TestAESDecryptInvalidCiphertext(t *testing.T) {
	key := "1234567890123456"

	// Invalid base64
	_, err := AESDecrypt("not-valid-base64!!!", key)
	if err == nil {
		t.Error("AESDecrypt should fail with invalid base64")
	}

	// Too short ciphertext
	_, err = AESDecrypt("YWJj", key) // "abc" in base64
	if err == nil {
		t.Error("AESDecrypt should fail with ciphertext too short")
	}
}

// ============================================================================
// Password Hashing Tests
// ============================================================================

func TestHashPassword(t *testing.T) {
	password := "mySecurePassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == password {
		t.Error("Hash should not equal plaintext password")
	}

	if len(hash) == 0 {
		t.Error("Hash should not be empty")
	}

	// Hash should start with bcrypt identifier
	if hash[0] != '$' {
		t.Error("Hash should start with bcrypt identifier")
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "mySecurePassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// Correct password should verify
	if !VerifyPassword(hash, password) {
		t.Error("VerifyPassword should return true for correct password")
	}

	// Wrong password should not verify
	if VerifyPassword(hash, "wrongPassword") {
		t.Error("VerifyPassword should return false for wrong password")
	}
}

func TestHashPasswordUniqueness(t *testing.T) {
	password := "samePassword"

	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)

	// Same password should produce different hashes (due to salt)
	if hash1 == hash2 {
		t.Error("Same password should produce different hashes due to salt")
	}

	// Both hashes should verify correctly
	if !VerifyPassword(hash1, password) {
		t.Error("First hash should verify correctly")
	}
	if !VerifyPassword(hash2, password) {
		t.Error("Second hash should verify correctly")
	}
}
