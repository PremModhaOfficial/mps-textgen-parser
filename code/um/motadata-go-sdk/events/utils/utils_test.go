package utils

import (
	"errors"
	"testing"
	"time"
)

// ============================================================================
// Subject Validation Tests
// ============================================================================

func TestValidateSubject(t *testing.T) {
	tests := []struct {
		name    string
		subject string
		want    bool
	}{
		{"valid simple", "events", true},
		{"valid dotted", "events.user.created", true},
		{"valid with wildcard", "events.*.created", true},
		{"valid with multi-wildcard", "events.>", true},
		{"valid with underscore", "events_v2.user", true},
		{"valid with dash", "events-v2.user", true},
		{"empty", "", false},
		{"too long", string(make([]byte, 300)), false},
		{"invalid space", "events user", false},
		{"invalid special char", "events@user", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateSubject(tt.subject); got != tt.want {
				t.Errorf("ValidateSubject(%q) = %v, want %v", tt.subject, got, tt.want)
			}
		})
	}
}

func TestSanitizeSubject(t *testing.T) {
	tests := []struct {
		name    string
		subject string
		want    string
	}{
		{"simple", "events", "events"},
		{"with spaces", "events user", "eventsuser"},
		{"with special chars", "events@user#created", "eventsusercreated"},
		{"valid chars preserved", "events.user_created-v2", "events.user_created-v2"},
		{"empty", "", ""},
		{"wildcards preserved", "events.*.>", "events.*.>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizeSubject(tt.subject); got != tt.want {
				t.Errorf("SanitizeSubject(%q) = %q, want %q", tt.subject, got, tt.want)
			}
		})
	}
}

func TestBuildSubject(t *testing.T) {
	tests := []struct {
		name   string
		tokens []string
		want   string
	}{
		{"single", []string{"events"}, "events"},
		{"multiple", []string{"events", "user", "created"}, "events.user.created"},
		{"with empty", []string{"events", "", "created"}, "events.created"},
		{"all empty", []string{"", ""}, ""},
		{"no tokens", []string{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildSubject(tt.tokens...); got != tt.want {
				t.Errorf("BuildSubject(%v) = %q, want %q", tt.tokens, got, tt.want)
			}
		})
	}
}

func TestParseSubject(t *testing.T) {
	tests := []struct {
		name    string
		subject string
		want    []string
	}{
		{"simple", "events", []string{"events"}},
		{"dotted", "events.user.created", []string{"events", "user", "created"}},
		{"empty", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseSubject(tt.subject)
			if len(got) != len(tt.want) {
				t.Errorf("ParseSubject(%q) = %v, want %v", tt.subject, got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("ParseSubject(%q)[%d] = %q, want %q", tt.subject, i, got[i], tt.want[i])
				}
			}
		})
	}
}

// ============================================================================
// ID Generation Tests
// ============================================================================

func TestGenerateID(t *testing.T) {
	tests := []struct {
		name       string
		byteLength int
		wantLen    int
	}{
		{"default", 0, 32},   // 16 bytes * 2 (hex)
		{"16 bytes", 16, 32}, // 16 bytes * 2 (hex)
		{"8 bytes", 8, 16},   // 8 bytes * 2 (hex)
		{"4 bytes", 4, 8},    // 4 bytes * 2 (hex)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := GenerateID(tt.byteLength)
			if len(id) != tt.wantLen {
				t.Errorf("GenerateID(%d) length = %d, want %d", tt.byteLength, len(id), tt.wantLen)
			}
		})
	}
}

func TestGenerateMessageID(t *testing.T) {
	id := GenerateMessageID()
	if len(id) != 32 {
		t.Errorf("GenerateMessageID() length = %d, want 32", len(id))
	}

	// Ensure uniqueness
	id2 := GenerateMessageID()
	if id == id2 {
		t.Error("GenerateMessageID() should generate unique IDs")
	}
}

func TestGenerateCorrelationID(t *testing.T) {
	id := GenerateCorrelationID()
	if len(id) != 24 {
		t.Errorf("GenerateCorrelationID() length = %d, want 24", len(id))
	}
}

func TestGenerateTraceID(t *testing.T) {
	id := GenerateTraceID()
	if len(id) != 32 {
		t.Errorf("GenerateTraceID() length = %d, want 32", len(id))
	}
}

func TestGenerateSpanID(t *testing.T) {
	id := GenerateSpanID()
	if len(id) != 16 {
		t.Errorf("GenerateSpanID() length = %d, want 16", len(id))
	}
}

// ============================================================================
// String Utilities Tests
// ============================================================================

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		maxLen int
		want   string
	}{
		{"shorter", "hello", 10, "hello"},
		{"exact", "hello", 5, "hello"},
		{"truncate", "hello world", 8, "hello..."},
		{"very short max", "hello", 2, "he"},
		{"zero max", "hello", 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TruncateString(tt.s, tt.maxLen); got != tt.want {
				t.Errorf("TruncateString(%q, %d) = %q, want %q", tt.s, tt.maxLen, got, tt.want)
			}
		})
	}
}

func TestCoalesceString(t *testing.T) {
	tests := []struct {
		name   string
		values []string
		want   string
	}{
		{"first non-empty", []string{"", "hello", "world"}, "hello"},
		{"first value", []string{"hello", "world"}, "hello"},
		{"all empty", []string{"", ""}, ""},
		{"no values", []string{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CoalesceString(tt.values...); got != tt.want {
				t.Errorf("CoalesceString(%v) = %q, want %q", tt.values, got, tt.want)
			}
		})
	}
}

func TestStringOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		value        string
		defaultValue string
		want         string
	}{
		{"non-empty", "hello", "default", "hello"},
		{"empty", "", "default", "default"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringOrDefault(tt.value, tt.defaultValue); got != tt.want {
				t.Errorf("StringOrDefault(%q, %q) = %q, want %q", tt.value, tt.defaultValue, got, tt.want)
			}
		})
	}
}

// ============================================================================
// Reflection Utilities Tests
// ============================================================================

func TestIsZeroValue(t *testing.T) {
	// Test with concrete types to avoid interface boxing issues
	t.Run("zero int", func(t *testing.T) {
		if !IsZeroValue(0) {
			t.Error("IsZeroValue(0) should be true")
		}
	})

	t.Run("non-zero int", func(t *testing.T) {
		if IsZeroValue(42) {
			t.Error("IsZeroValue(42) should be false")
		}
	})

	t.Run("empty string", func(t *testing.T) {
		if !IsZeroValue("") {
			t.Error("IsZeroValue(\"\") should be true")
		}
	})

	t.Run("non-empty string", func(t *testing.T) {
		if IsZeroValue("hello") {
			t.Error("IsZeroValue(\"hello\") should be false")
		}
	})

	t.Run("nil pointer", func(t *testing.T) {
		var p *int
		if !IsZeroValue(p) {
			t.Error("IsZeroValue(nil pointer) should be true")
		}
	})

	t.Run("non-nil pointer", func(t *testing.T) {
		val := 42
		if IsZeroValue(&val) {
			t.Error("IsZeroValue(non-nil pointer) should be false")
		}
	})

	t.Run("zero struct", func(t *testing.T) {
		type testStruct struct{ X int }
		if !IsZeroValue(testStruct{}) {
			t.Error("IsZeroValue(zero struct) should be true")
		}
	})

	t.Run("non-zero struct", func(t *testing.T) {
		type testStruct struct{ X int }
		if IsZeroValue(testStruct{X: 1}) {
			t.Error("IsZeroValue(non-zero struct) should be false")
		}
	})
}

func TestIsNil(t *testing.T) {
	var nilPtr *int
	var nilMap map[string]int
	var nilSlice []int
	var nilChan chan int

	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"nil", nil, true},
		{"nil pointer", nilPtr, true},
		{"nil map", nilMap, true},
		{"nil slice", nilSlice, true},
		{"nil chan", nilChan, true},
		{"non-nil int", 42, false},
		{"non-nil string", "hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNil(tt.value); got != tt.want {
				t.Errorf("IsNil(%v) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}
}

// ============================================================================
// Time Utilities Tests
// ============================================================================

func TestUnixMillis(t *testing.T) {
	before := time.Now().UnixMilli()
	got := UnixMillis()
	after := time.Now().UnixMilli()

	if got < before || got > after {
		t.Errorf("UnixMillis() = %d, want between %d and %d", got, before, after)
	}
}

func TestTimestampString(t *testing.T) {
	ts := TimestampString()
	_, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		t.Errorf("TimestampString() = %q is not valid RFC3339: %v", ts, err)
	}
}

func TestDurationString(t *testing.T) {
	tests := []struct {
		name string
		d    time.Duration
	}{
		{"nanoseconds", 500 * time.Nanosecond},
		{"microseconds", 500 * time.Microsecond},
		{"milliseconds", 500 * time.Millisecond},
		{"seconds", 5 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DurationString(tt.d)
			if got == "" {
				t.Errorf("DurationString(%v) returned empty string", tt.d)
			}
		})
	}
}

// ============================================================================
// Slice Utilities Tests
// ============================================================================

func TestContains(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	if !Contains(slice, 3) {
		t.Error("Contains should return true for existing element")
	}

	if Contains(slice, 10) {
		t.Error("Contains should return false for non-existing element")
	}
}

func TestUnique(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		want  []int
	}{
		{"with duplicates", []int{1, 2, 2, 3, 3, 3}, []int{1, 2, 3}},
		{"no duplicates", []int{1, 2, 3}, []int{1, 2, 3}},
		{"empty", []int{}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Unique(tt.slice)
			if len(got) != len(tt.want) {
				t.Errorf("Unique(%v) = %v, want %v", tt.slice, got, tt.want)
			}
		})
	}
}

func TestFilter(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}
	got := Filter(slice, func(n int) bool { return n%2 == 0 })

	if len(got) != 2 || got[0] != 2 || got[1] != 4 {
		t.Errorf("Filter() = %v, want [2, 4]", got)
	}
}

func TestMap(t *testing.T) {
	slice := []int{1, 2, 3}
	got := Map(slice, func(n int) int { return n * 2 })

	if len(got) != 3 || got[0] != 2 || got[1] != 4 || got[2] != 6 {
		t.Errorf("Map() = %v, want [2, 4, 6]", got)
	}
}

// ============================================================================
// Map Utilities Tests
// ============================================================================

func TestMergeMaps(t *testing.T) {
	m1 := map[string]int{"a": 1, "b": 2}
	m2 := map[string]int{"b": 3, "c": 4}

	got := MergeMaps(m1, m2)

	if got["a"] != 1 || got["b"] != 3 || got["c"] != 4 {
		t.Errorf("MergeMaps() = %v, want map[a:1 b:3 c:4]", got)
	}
}

func TestCopyMap(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	got := CopyMap(m)

	got["c"] = 3
	if _, exists := m["c"]; exists {
		t.Error("CopyMap should create independent copy")
	}
}

func TestKeys(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	keys := Keys(m)

	if len(keys) != 2 {
		t.Errorf("Keys() length = %d, want 2", len(keys))
	}
}

func TestValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	values := Values(m)

	if len(values) != 2 {
		t.Errorf("Values() length = %d, want 2", len(values))
	}
}

// ============================================================================
// Types Tests
// ============================================================================

func TestOptional(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		opt := Some(42)
		if !opt.IsPresent() {
			t.Error("Some should be present")
		}
		if opt.Value() != 42 {
			t.Errorf("Value() = %d, want 42", opt.Value())
		}
	})

	t.Run("None", func(t *testing.T) {
		opt := None[int]()
		if opt.IsPresent() {
			t.Error("None should not be present")
		}
		if opt.ValueOr(10) != 10 {
			t.Errorf("ValueOr(10) = %d, want 10", opt.ValueOr(10))
		}
	})
}

func TestResult(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		result := Ok(42)
		if !result.IsOk() {
			t.Error("Ok should be successful")
		}
		if result.Value() != 42 {
			t.Errorf("Value() = %d, want 42", result.Value())
		}
	})

	t.Run("Err", func(t *testing.T) {
		result := Err[int](errors.New("test error"))
		if !result.IsErr() {
			t.Error("Err should have error")
		}
		if result.ValueOr(10) != 10 {
			t.Errorf("ValueOr(10) = %d, want 10", result.ValueOr(10))
		}
	})
}

func TestSafeMap(t *testing.T) {
	m := NewSafeMap[string, int]()

	m.Set("a", 1)
	m.Set("b", 2)

	if v, ok := m.Get("a"); !ok || v != 1 {
		t.Errorf("Get(a) = %d, %v, want 1, true", v, ok)
	}

	if m.Len() != 2 {
		t.Errorf("Len() = %d, want 2", m.Len())
	}

	m.Delete("a")
	if m.Has("a") {
		t.Error("Has(a) should be false after delete")
	}
}

func TestSafeSlice(t *testing.T) {
	s := NewSafeSlice[int]()

	s.Append(1, 2, 3)

	if s.Len() != 3 {
		t.Errorf("Len() = %d, want 3", s.Len())
	}

	if v, ok := s.Get(1); !ok || v != 2 {
		t.Errorf("Get(1) = %d, %v, want 2, true", v, ok)
	}

	if v, ok := s.Pop(); !ok || v != 3 {
		t.Errorf("Pop() = %d, %v, want 3, true", v, ok)
	}

	if s.Len() != 2 {
		t.Errorf("Len() after Pop() = %d, want 2", s.Len())
	}
}

// ============================================================================
// Error Tests
// ============================================================================

func TestWrapError(t *testing.T) {
	original := errors.New("original error")
	wrapped := WrapError(original, "context")

	if wrapped == nil {
		t.Error("WrapError should not return nil")
	}

	if !errors.Is(wrapped, original) {
		t.Error("wrapped error should contain original")
	}
}

func TestErrorCollector(t *testing.T) {
	ec := NewErrorCollector()

	ec.Add(errors.New("error 1"))
	ec.Add(nil) // should be ignored
	ec.Add(errors.New("error 2"))

	if ec.Count() != 2 {
		t.Errorf("Count() = %d, want 2", ec.Count())
	}

	err := ec.Error()
	if err == nil {
		t.Error("Error() should not return nil")
	}
}

func TestValidationError(t *testing.T) {
	err := NewValidationError("email", "invalid format")
	expected := "validation error: email: invalid format"

	if err.Error() != expected {
		t.Errorf("Error() = %q, want %q", err.Error(), expected)
	}
}

func TestRootCause(t *testing.T) {
	root := errors.New("root error")
	wrapped := WrapError(WrapError(root, "level 1"), "level 2")

	cause := RootCause(wrapped)
	if cause != root {
		t.Errorf("RootCause() = %v, want %v", cause, root)
	}
}
