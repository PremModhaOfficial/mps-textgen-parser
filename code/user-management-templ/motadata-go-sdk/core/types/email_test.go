package types

import (
	"testing"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Email
		wantErr bool
	}{
		{"valid email", "user@example.com", Email("user@example.com"), false},
		{"valid with plus", "user+tag@example.com", Email("user+tag@example.com"), false},
		{"valid with dots", "first.last@example.com", Email("first.last@example.com"), false},
		{"trims whitespace", "  user@example.com  ", Email("user@example.com"), false},
		{"empty string", "", "", true},
		{"no at sign", "userexample.com", "", true},
		{"no domain", "user@", "", true},
		{"no local", "@example.com", "", true},
		{"no tld", "user@example", "", true},
		{"double at", "user@@example.com", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEmail(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmail(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NewEmail(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestMustEmail(t *testing.T) {
	// Valid email should not panic.
	e := MustEmail("user@example.com")
	if e != Email("user@example.com") {
		t.Errorf("MustEmail returned %v, want user@example.com", e)
	}

	// Invalid email should panic.
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustEmail with invalid email should panic")
		}
	}()
	MustEmail("invalid")
}

func TestEmailIsValid(t *testing.T) {
	tests := []struct {
		email Email
		want  bool
	}{
		{Email("user@example.com"), true},
		{Email("a@b.co"), true},
		{Email("user+tag@example.com"), true},
		{Email("user.name@example.com"), true},
		{Email(""), false},
		{Email("noatsign"), false},
		{Email("user@"), false},
		{Email("@domain.com"), false},
		{Email("user@domain"), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.email), func(t *testing.T) {
			if got := tt.email.IsValid(); got != tt.want {
				t.Errorf("Email(%q).IsValid() = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestEmailLocal(t *testing.T) {
	tests := []struct {
		email Email
		want  string
	}{
		{Email("user@example.com"), "user"},
		{Email("first.last@example.com"), "first.last"},
		{Email("user+tag@example.com"), "user+tag"},
		{Email("noatsign"), "noatsign"},
	}

	for _, tt := range tests {
		t.Run(string(tt.email), func(t *testing.T) {
			if got := tt.email.Local(); got != tt.want {
				t.Errorf("Email(%q).Local() = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestEmailDomain(t *testing.T) {
	tests := []struct {
		email Email
		want  string
	}{
		{Email("user@example.com"), "example.com"},
		{Email("user@sub.example.com"), "sub.example.com"},
		{Email("noatsign"), ""},
	}

	for _, tt := range tests {
		t.Run(string(tt.email), func(t *testing.T) {
			if got := tt.email.Domain(); got != tt.want {
				t.Errorf("Email(%q).Domain() = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestEmailNormalize(t *testing.T) {
	tests := []struct {
		email Email
		want  Email
	}{
		{Email("User@Example.COM"), Email("user@example.com")},
		{Email("  User@Example.COM  "), Email("user@example.com")},
		{Email("already@lower.com"), Email("already@lower.com")},
	}

	for _, tt := range tests {
		t.Run(string(tt.email), func(t *testing.T) {
			if got := tt.email.Normalize(); got != tt.want {
				t.Errorf("Email(%q).Normalize() = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}

func TestEmailString(t *testing.T) {
	e := Email("user@example.com")
	if got := e.String(); got != "user@example.com" {
		t.Errorf("Email.String() = %v, want user@example.com", got)
	}
}

func TestEmailIsEmpty(t *testing.T) {
	if !Email("").IsEmpty() {
		t.Error("Email(\"\").IsEmpty() should be true")
	}
	if Email("user@example.com").IsEmpty() {
		t.Error("Email(\"user@example.com\").IsEmpty() should be false")
	}
}

func TestEmailEquals(t *testing.T) {
	tests := []struct {
		a, b Email
		want bool
	}{
		{Email("user@example.com"), Email("user@example.com"), true},
		{Email("User@Example.COM"), Email("user@example.com"), true},
		{Email("user@example.com"), Email("other@example.com"), false},
		{Email(""), Email(""), true},
	}

	for _, tt := range tests {
		t.Run(string(tt.a)+"_vs_"+string(tt.b), func(t *testing.T) {
			if got := tt.a.Equals(tt.b); got != tt.want {
				t.Errorf("Email(%q).Equals(%q) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
