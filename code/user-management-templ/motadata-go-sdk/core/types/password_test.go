package types

import (
	"testing"
)

func TestNewPassword(t *testing.T) {
	p := NewPassword("secret123")
	if string(p) != "secret123" {
		t.Errorf("NewPassword() = %v, want secret123", string(p))
	}
}

func TestPasswordLen(t *testing.T) {
	tests := []struct {
		name string
		pass Password
		want int
	}{
		{"empty", Password(""), 0},
		{"ascii", Password("hello"), 5},
		{"unicode", Password("héllo"), 5},
		{"emoji", Password("🔒key"), 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pass.Len(); got != tt.want {
				t.Errorf("Password(%q).Len() = %v, want %v", string(tt.pass), got, tt.want)
			}
		})
	}
}

func TestPasswordHasUpper(t *testing.T) {
	tests := []struct {
		pass Password
		want bool
	}{
		{Password("Hello"), true},
		{Password("HELLO"), true},
		{Password("hello"), false},
		{Password("123!@#"), false},
		{Password(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.pass), func(t *testing.T) {
			if got := tt.pass.HasUpper(); got != tt.want {
				t.Errorf("Password(%q).HasUpper() = %v, want %v", string(tt.pass), got, tt.want)
			}
		})
	}
}

func TestPasswordHasLower(t *testing.T) {
	tests := []struct {
		pass Password
		want bool
	}{
		{Password("Hello"), true},
		{Password("hello"), true},
		{Password("HELLO"), false},
		{Password("123!@#"), false},
		{Password(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.pass), func(t *testing.T) {
			if got := tt.pass.HasLower(); got != tt.want {
				t.Errorf("Password(%q).HasLower() = %v, want %v", string(tt.pass), got, tt.want)
			}
		})
	}
}

func TestPasswordHasDigit(t *testing.T) {
	tests := []struct {
		pass Password
		want bool
	}{
		{Password("hello1"), true},
		{Password("123"), true},
		{Password("hello"), false},
		{Password("!@#$"), false},
		{Password(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.pass), func(t *testing.T) {
			if got := tt.pass.HasDigit(); got != tt.want {
				t.Errorf("Password(%q).HasDigit() = %v, want %v", string(tt.pass), got, tt.want)
			}
		})
	}
}

func TestPasswordHasSpecial(t *testing.T) {
	tests := []struct {
		pass Password
		want bool
	}{
		{Password("hello!"), true},
		{Password("pass@word"), true},
		{Password("hello"), false},
		{Password("Hello123"), false},
		{Password(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.pass), func(t *testing.T) {
			if got := tt.pass.HasSpecial(); got != tt.want {
				t.Errorf("Password(%q).HasSpecial() = %v, want %v", string(tt.pass), got, tt.want)
			}
		})
	}
}

func TestPasswordIsEmpty(t *testing.T) {
	if !Password("").IsEmpty() {
		t.Error("Password(\"\").IsEmpty() should be true")
	}
	if Password("secret").IsEmpty() {
		t.Error("Password(\"secret\").IsEmpty() should be false")
	}
}

func TestPasswordMeetsPolicy(t *testing.T) {
	tests := []struct {
		name           string
		pass           Password
		minLen         int
		upper, lower   bool
		digit, special bool
		want           bool
	}{
		{"all requirements met", Password("Hello1!"), 6, true, true, true, true, true},
		{"too short", Password("Hi1!"), 6, true, true, true, true, false},
		{"missing upper", Password("hello1!"), 6, true, true, true, true, false},
		{"missing lower", Password("HELLO1!"), 6, true, true, true, true, false},
		{"missing digit", Password("Hello!!"), 6, true, true, true, true, false},
		{"missing special", Password("Hello12"), 6, true, true, true, true, false},
		{"no requirements", Password("a"), 1, false, false, false, false, true},
		{"only length", Password("abcdef"), 6, false, false, false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pass.MeetsPolicy(tt.minLen, tt.upper, tt.lower, tt.digit, tt.special)
			if got != tt.want {
				t.Errorf("Password(%q).MeetsPolicy() = %v, want %v", string(tt.pass), got, tt.want)
			}
		})
	}
}

func TestPasswordString(t *testing.T) {
	p := Password("supersecret")
	if got := p.String(); got != "****" {
		t.Errorf("Password.String() = %v, want ****", got)
	}
}

func TestPasswordHashAndVerify(t *testing.T) {
	p := NewPassword("mypassword123")

	hashed, err := p.Hash()
	if err != nil {
		t.Fatalf("Password.Hash() error = %v", err)
	}
	if hashed.IsEmpty() {
		t.Error("Hash should not be empty")
	}

	// Correct password should verify.
	if !hashed.Verify(p) {
		t.Error("HashedPassword.Verify() should return true for correct password")
	}

	// Wrong password should not verify.
	if hashed.Verify(NewPassword("wrongpassword")) {
		t.Error("HashedPassword.Verify() should return false for wrong password")
	}
}

func TestPasswordHashWithCost(t *testing.T) {
	p := NewPassword("testpass")

	// Low cost for fast test.
	hashed, err := p.HashWithCost(4)
	if err != nil {
		t.Fatalf("Password.HashWithCost() error = %v", err)
	}
	if !hashed.Verify(p) {
		t.Error("HashedPassword.Verify() should return true after HashWithCost")
	}
}

func TestHashedPasswordString(t *testing.T) {
	h := HashedPassword("$2a$10$somehashvalue")
	if got := h.String(); got != "$2a$10$somehashvalue" {
		t.Errorf("HashedPassword.String() = %v, want the hash itself", got)
	}
}

func TestHashedPasswordIsEmpty(t *testing.T) {
	if !HashedPassword("").IsEmpty() {
		t.Error("HashedPassword(\"\").IsEmpty() should be true")
	}
	if HashedPassword("somehash").IsEmpty() {
		t.Error("HashedPassword(\"somehash\").IsEmpty() should be false")
	}
}
