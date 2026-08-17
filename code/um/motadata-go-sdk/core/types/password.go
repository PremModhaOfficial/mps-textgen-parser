package types

import (
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// Password is a type-safe plaintext password wrapper.
// String() returns a masked value to prevent accidental logging.
type Password string

// NewPassword creates a Password from a plaintext string.
func NewPassword(s string) Password {
	return Password(s)
}

// Len returns the number of runes in the password.
func (p Password) Len() int {
	return len([]rune(string(p)))
}

// HasUpper returns true if the password contains at least one uppercase letter.
func (p Password) HasUpper() bool {
	for _, r := range string(p) {
		if unicode.IsUpper(r) {
			return true
		}
	}
	return false
}

// HasLower returns true if the password contains at least one lowercase letter.
func (p Password) HasLower() bool {
	for _, r := range string(p) {
		if unicode.IsLower(r) {
			return true
		}
	}
	return false
}

// HasDigit returns true if the password contains at least one digit.
func (p Password) HasDigit() bool {
	for _, r := range string(p) {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

// HasSpecial returns true if the password contains at least one non-alphanumeric character.
func (p Password) HasSpecial() bool {
	for _, r := range string(p) {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

// IsEmpty returns true if the password is the zero value.
func (p Password) IsEmpty() bool {
	return len(p) == 0
}

// MeetsPolicy checks the password against the given policy requirements.
func (p Password) MeetsPolicy(minLen int, requireUpper, requireLower, requireDigit, requireSpecial bool) bool {
	if p.Len() < minLen {
		return false
	}
	if requireUpper && !p.HasUpper() {
		return false
	}
	if requireLower && !p.HasLower() {
		return false
	}
	if requireDigit && !p.HasDigit() {
		return false
	}
	if requireSpecial && !p.HasSpecial() {
		return false
	}
	return true
}

// Hash hashes the password using bcrypt with the default cost.
func (p Password) Hash() (HashedPassword, error) {
	return p.HashWithCost(bcrypt.DefaultCost)
}

// HashWithCost hashes the password using bcrypt with the specified cost.
func (p Password) HashWithCost(cost int) (HashedPassword, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p), cost)
	if err != nil {
		return "", err
	}
	return HashedPassword(hash), nil
}

// String returns a masked value to prevent accidental logging.
// Use string(p) to access the plaintext when needed.
func (p Password) String() string {
	return "****"
}

// HashedPassword is a bcrypt-hashed password.
type HashedPassword string

// Verify checks whether the plaintext password matches this hash.
func (h HashedPassword) Verify(plain Password) bool {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(plain)) == nil
}

// String returns the hash string.
func (h HashedPassword) String() string {
	return string(h)
}

// IsEmpty returns true if the hashed password is the zero value.
func (h HashedPassword) IsEmpty() bool {
	return len(h) == 0
}
