package types

import (
	"errors"
	"strings"
)

// Email is a type-safe email address representation with validation and utility methods.
type Email string

// NewEmail creates a new Email after validating the format.
// Returns an error if the email format is invalid.
func NewEmail(s string) (Email, error) {
	e := Email(strings.TrimSpace(s))
	if !e.IsValid() {
		return "", errors.New("invalid email format")
	}
	return e, nil
}

// MustEmail creates a new Email, panicking if the format is invalid.
// Use for constants and tests only.
func MustEmail(s string) Email {
	e, err := NewEmail(s)
	if err != nil {
		panic(err)
	}
	return e
}

// IsValid returns true if the email matches a valid email format.
func (e Email) IsValid() bool {
	s := string(e)
	if len(s) == 0 || len(s) > 254 {
		return false
	}
	return emailRegex.MatchString(s)
}

// Local returns the local part of the email (before @).
func (e Email) Local() string {
	s := string(e)
	at := strings.LastIndexByte(s, '@')
	if at < 0 {
		return s
	}
	return s[:at]
}

// Domain returns the domain part of the email (after @).
func (e Email) Domain() string {
	s := string(e)
	at := strings.LastIndexByte(s, '@')
	if at < 0 {
		return ""
	}
	return s[at+1:]
}

// Normalize returns the email lowercased and trimmed.
func (e Email) Normalize() Email {
	return Email(strings.ToLower(strings.TrimSpace(string(e))))
}

// String returns the email as a string.
func (e Email) String() string {
	return string(e)
}

// IsEmpty returns true if the email is the zero value.
func (e Email) IsEmpty() bool {
	return len(e) == 0
}

// Equals performs a case-insensitive comparison with another Email.
func (e Email) Equals(other Email) bool {
	return strings.EqualFold(string(e), string(other))
}
