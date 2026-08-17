package types

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"hash/crc32"
	"io"
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-faster/city"
	"golang.org/x/crypto/bcrypt"
)

// String helper functions provide common string operations.
// All functions take a string parameter and return the result.
//
// Usage:
//
//	if types.IsNotEmpty(s) {
//	    upper := types.ToUpper(s)
//	    fmt.Println(upper)
//	}
//
//	hash := types.Hash("hello world")
//	encoded := types.Base64Encode("hello")

// ============================================================================
// Empty/Blank Checks
// ============================================================================

// IsEmpty returns true if the string has zero length.
func IsEmpty(s string) bool {
	return len(s) == 0
}

// IsBlank returns true if the string is empty or contains only whitespace characters.
func IsBlank(s string) bool {
	if len(s) == 0 {
		return true
	}
	for _, r := range s {
		if !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}

// IsNotEmpty returns true if the string has at least one character.
func IsNotEmpty(s string) bool {
	return len(s) > 0
}

// IsNotBlank returns true if the string contains at least one non-whitespace character.
func IsNotBlank(s string) bool {
	return !IsBlank(s)
}

// ============================================================================
// Character Type Checks
// ============================================================================

// IsASCII returns true if all characters in the string are ASCII (0-127).
func IsASCII(s string) bool {
	if len(s) == 0 {
		return true
	}
	for i := 0; i < len(s); i++ {
		if s[i] > 127 {
			return false
		}
	}
	return true
}

// IsAlpha returns true if all characters in the string are alphabetic (a-z, A-Z).
// Returns false for empty strings.
func IsAlpha(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// IsNumeric returns true if all characters in the string are numeric (0-9).
// Returns false for empty strings.
func IsNumeric(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// IsAlphaNumeric returns true if all characters are alphabetic or numeric.
// Returns false for empty strings.
func IsAlphaNumeric(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// IsLower returns true if all cased characters in the string are lowercase.
// Returns false for empty strings or strings with no cased characters.
func IsLower(s string) bool {
	if len(s) == 0 {
		return false
	}
	hasCased := false
	for _, r := range s {
		if unicode.IsUpper(r) {
			return false
		}
		if unicode.IsLower(r) {
			hasCased = true
		}
	}
	return hasCased
}

// IsUpper returns true if all cased characters in the string are uppercase.
// Returns false for empty strings or strings with no cased characters.
func IsUpper(s string) bool {
	if len(s) == 0 {
		return false
	}
	hasCased := false
	for _, r := range s {
		if unicode.IsLower(r) {
			return false
		}
		if unicode.IsUpper(r) {
			hasCased = true
		}
	}
	return hasCased
}

// IsDigit returns true if all characters in the string are digits (0-9).
// This is an alias for IsNumeric.
// Returns false for empty strings.
func IsDigit(s string) bool {
	return IsNumeric(s)
}

// IsHex returns true if all characters in the string are valid hexadecimal digits (0-9, a-f, A-F).
// Returns false for empty strings.
func IsHex(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
			return false
		}
	}
	return true
}

// ============================================================================
// Format Validation
// ============================================================================

// emailRegex is a compiled regex for basic email validation.
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// IsEmail returns true if the string is a valid email address format.
// This performs basic format validation, not deliverability checking.
func IsEmail(s string) bool {
	if len(s) == 0 {
		return false
	}
	return emailRegex.MatchString(s)
}

// IsURL returns true if the string is a valid URL format.
// Validates that the string can be parsed as a URL with a scheme and host.
func IsURL(s string) bool {
	if len(s) == 0 {
		return false
	}
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}

// ============================================================================
// Encoding/Decoding
// ============================================================================

// ToBytes converts the string to a byte slice.
func ToBytes(s string) []byte {
	return []byte(s)
}

// Base64Encode returns the base64 encoded string.
func Base64Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// Base64Decode returns the base64 decoded string and any error.
func Base64Decode(s string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}

// URLEncode returns the URL encoded (percent-encoded) string.
func URLEncode(s string) string {
	return url.QueryEscape(s)
}

// URLDecode returns the URL decoded string and any error.
func URLDecode(s string) (string, error) {
	return url.QueryUnescape(s)
}

// ============================================================================
// Trimming and Whitespace
// ============================================================================

// Trim returns the string with leading and trailing characters removed.
// The cutset string specifies the set of characters to remove.
func Trim(s, cutset string) string {
	return strings.Trim(s, cutset)
}

// TrimLeft returns the string with leading characters removed.
// The cutset string specifies the set of characters to remove.
func TrimLeft(s, cutset string) string {
	return strings.TrimLeft(s, cutset)
}

// TrimRight returns the string with trailing characters removed.
// The cutset string specifies the set of characters to remove.
func TrimRight(s, cutset string) string {
	return strings.TrimRight(s, cutset)
}

// TrimSpace returns the string with leading and trailing whitespace removed.
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// whitespaceRegex matches one or more whitespace characters.
var whitespaceRegex = regexp.MustCompile(`\s+`)

// NormalizeSpace trims the string and replaces all whitespace sequences with a single space.
func NormalizeSpace(s string) string {
	trimmed := strings.TrimSpace(s)
	return whitespaceRegex.ReplaceAllString(trimmed, " ")
}

// CollapseWhitespace replaces all sequences of whitespace with a single space.
// Unlike NormalizeSpace, this does not trim leading/trailing whitespace.
func CollapseWhitespace(s string) string {
	return whitespaceRegex.ReplaceAllString(s, " ")
}

// RemoveWhitespace removes all whitespace characters from the string.
func RemoveWhitespace(s string) string {
	var result strings.Builder
	result.Grow(len(s))
	for _, r := range s {
		if !unicode.IsSpace(r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ============================================================================
// Case Conversion
// ============================================================================

// ToLower returns the string with all characters converted to lowercase.
func ToLower(s string) string {
	return strings.ToLower(s)
}

// ToUpper returns the string with all characters converted to uppercase.
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// ToTitle returns the string with the first letter of each word capitalized.
func ToTitle(s string) string {
	return strings.Title(s)
}

// ============================================================================
// Hashing and Checksums (using CityHash)
// ============================================================================

// Hash64 returns the 64-bit CityHash of the string as a uint64.
// CityHash is a fast, non-cryptographic hash function developed by Google.
func Hash64(s string) uint64 {
	return city.Hash64([]byte(s))
}

// Hash32 returns the 32-bit CityHash of the string as a uint32.
func Hash32(s string) uint32 {
	return city.Hash32([]byte(s))
}

// Hash128 returns the 128-bit CityHash of the string as a city.U128.
func Hash128(s string) city.U128 {
	return city.Hash128([]byte(s))
}

// Checksum returns the CRC32 checksum of the string as a uint32.
func Checksum(s string) uint32 {
	return crc32.ChecksumIEEE([]byte(s))
}

// ============================================================================
// Encryption/Decryption (using AES-GCM)
// ============================================================================

// AESEncrypt encrypts the plaintext using AES-GCM with the provided key.
// The key must be 16, 24, or 32 bytes for AES-128, AES-192, or AES-256.
// Returns the base64-encoded ciphertext (nonce prepended to ciphertext).
func AESEncrypt(plaintext, key string) (string, error) {
	keyBytes := []byte(key)
	if len(keyBytes) != 16 && len(keyBytes) != 24 && len(keyBytes) != 32 {
		return "", errors.New("key must be 16, 24, or 32 bytes")
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return Base64Encode(string(ciphertext)), nil
}

// AESDecrypt decrypts the base64-encoded ciphertext using AES-GCM with the provided key.
// The key must be 16, 24, or 32 bytes for AES-128, AES-192, or AES-256.
// Returns the decrypted plaintext.
func AESDecrypt(ciphertext, key string) (string, error) {
	keyBytes := []byte(key)
	if len(keyBytes) != 16 && len(keyBytes) != 24 && len(keyBytes) != 32 {
		return "", errors.New("key must be 16, 24, or 32 bytes")
	}

	decoded, err := Base64Decode(ciphertext)
	if err != nil {
		return "", err
	}
	data := []byte(decoded)

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// ============================================================================
// Password Hashing (using bcrypt)
// ============================================================================

// HashPassword hashes a password using bcrypt with default cost.
// Note: bcrypt only uses the first 72 bytes of the password.
// Returns the hashed password as a string.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword compares a hashed password with a plaintext password.
// Returns true if the password matches the hash.
func VerifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
