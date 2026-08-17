package utils

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

/* ========================================================================================================
   SHARED TEST VARIABLES
   ======================================================================================================== */

// --- Strings shared across both test files ---
var (
	// Subject / token literals
	testSubjectEvents       = "events"
	testSubjectDotted       = "events.user.created"
	testSubjectWildcard     = "events.*.created"
	testSubjectMultiWild    = "events.>"
	testSubjectUnderscore   = "events_v2.user"
	testSubjectDash         = "events-v2.user"
	testSubjectSpace        = "events user"
	testSubjectSpecial      = "events@user"
	testSubjectSpecialHash  = "events@user#created"
	testSubjectPreserved    = "events.user_created-v2"
	testSubjectWildcardsAll = "events.*.>"
	testTokenUser           = "user"
	testTokenCreated        = "created"

	// String values
	testStrHello      = "hello"
	testStrWorld      = "world"
	testStrDefault    = "default"
	testStrHelloWorld = "hello world"
	testStrEmpty      = ""

	// Truncated values
	testStrHelloTrunc = "hello..."

	// Sanitized results
	testSanitizedSpace   = "eventsuser"
	testSanitizedSpecial = "eventsusercreated"

	// Error-related strings (shared with errors_core_test.go)
	testKindConnection = "connection"
	testOpConnect      = "connect"
	testOpDisconnect   = "disconnect"
	testTenantID       = "tenant-123"
	testFieldServers   = "Servers"
	testTypeMyStruct   = "MyStruct"
	testOpDeserialize  = "deserialize"
	testOpSerialize    = "serialize"
	testMsgSingleErr   = "single error"
	testMsgErr1        = "error 1"
	testMsgErr2        = "error 2"
	testMsgErr3        = "error 3"
	testMsgNoErrors    = "no errors"
	testMsgNotConn     = "not connected"
	testNATSServer     = "nats://localhost:4222"
	testMsgAtLeastOne  = "at least one server required"
	testMsgRequired    = "required"
	testMsgJsonErr     = "json unmarshal error"
	testMsgTestErr     = "test"
	testFieldEmail     = "email"
	testMsgInvalidFmt  = "invalid format"
	testExpectedValErr = "validation error: email: invalid format"

	// Error wrapping
	testMsgOriginalErr = "original error"
	testMsgContext     = "context"
	testMsgRootErr     = "root error"
	testMsgLevel1      = "level 1"
	testMsgLevel2      = "level 2"

	// Format error
	testMsgSuccess   = "publish: success"
	testMsgFmtErr    = "publish: original error"
	testOpPublish    = "publish"
	testMsgPanicStr  = "something went wrong"
	testMsgPanicInt  = "panic: 42"
	testMsgFmtWrap   = "operation %s failed"
	testMsgFmtResult = "operation publish failed: original error"
)

/* ======== SUBJECT VALIDATION TESTS ======== */

func TestValidateSubject(t *testing.T) {
	testCases := []struct {
		name    string
		subject string
		want    bool
	}{
		{"valid simple", testSubjectEvents, true},
		{"valid dotted", testSubjectDotted, true},
		{"valid with wildcard", testSubjectWildcard, true},
		{"valid with multi-wildcard", testSubjectMultiWild, true},
		{"valid with underscore", testSubjectUnderscore, true},
		{"valid with dash", testSubjectDash, true},
		{"empty", testStrEmpty, false},
		{"too long", string(make([]byte, 300)), false},
		{"invalid space", testSubjectSpace, false},
		{"invalid special char", testSubjectSpecial, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.want, ValidateSubject(tc.subject), "ValidateSubject(%q)", tc.subject)
		})
	}
}

func TestSanitizeSubject(t *testing.T) {
	testCases := []struct {
		name    string
		subject string
		want    string
	}{
		{"simple", testSubjectEvents, testSubjectEvents},
		{"with spaces", testSubjectSpace, testSanitizedSpace},
		{"with special chars", testSubjectSpecialHash, testSanitizedSpecial},
		{"valid chars preserved", testSubjectPreserved, testSubjectPreserved},
		{"empty", testStrEmpty, testStrEmpty},
		{"wildcards preserved", testSubjectWildcardsAll, testSubjectWildcardsAll},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.want, SanitizeSubject(tc.subject), "SanitizeSubject(%q)", tc.subject)
		})
	}
}

func TestBuildSubject(t *testing.T) {
	testCases := []struct {
		name   string
		tokens []string
		want   string
	}{
		{"single", []string{testSubjectEvents}, testSubjectEvents},
		{"multiple", []string{testSubjectEvents, testTokenUser, testTokenCreated}, testSubjectDotted},
		{"with empty", []string{testSubjectEvents, testStrEmpty, testTokenCreated}, "events.created"},
		{"all empty", []string{testStrEmpty, testStrEmpty}, testStrEmpty},
		{"no tokens", []string{}, testStrEmpty},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.want, BuildSubject(tc.tokens...), "BuildSubject(%v)", tc.tokens)
		})
	}
}

func TestParseSubject(t *testing.T) {
	testCases := []struct {
		name    string
		subject string
		want    []string
	}{
		{"simple", testSubjectEvents, []string{testSubjectEvents}},
		{"dotted", testSubjectDotted, []string{testSubjectEvents, testTokenUser, testTokenCreated}},
		{"empty", testStrEmpty, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.want, ParseSubject(tc.subject), "ParseSubject(%q)", tc.subject)
		})
	}
}

/* ======== ID GENERATION TESTS ======== */

func TestGenerateID(t *testing.T) {
	testCases := []struct {
		name       string
		byteLength int
		wantLen    int
	}{
		{"default", 0, 32},
		{"negative defaults to 16 bytes", -5, 32},
		{"16 bytes", 16, 32},
		{"8 bytes", 8, 16},
		{"4 bytes", 4, 8},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			id := GenerateID(tc.byteLength)
			assertions.Len(id, tc.wantLen, "GenerateID(%d) length", tc.byteLength)
		})
	}
}

func TestGenerateMessageID(t *testing.T) {
	assertions := assert.New(t)

	id := GenerateMessageID()
	assertions.Len(id, 32, "GenerateMessageID() length")

	id2 := GenerateMessageID()
	assertions.NotEqual(id, id2, "GenerateMessageID() should generate unique IDs")
}

func TestGenerateCorrelationID(t *testing.T) {
	assertions := assert.New(t)
	id := GenerateCorrelationID()
	assertions.Len(id, 24, "GenerateCorrelationID() length")
}

func TestGenerateTraceID(t *testing.T) {
	assertions := assert.New(t)
	id := GenerateTraceID()
	assertions.Len(id, 32, "GenerateTraceID() length")
}

func TestGenerateSpanID(t *testing.T) {
	assertions := assert.New(t)
	id := GenerateSpanID()
	assertions.Len(id, 16, "GenerateSpanID() length")
}

/* ======== STRING UTILITIES TESTS ======== */

func TestTruncateString(t *testing.T) {
	testCases := []struct {
		name   string
		s      string
		maxLen int
		want   string
	}{
		{"shorter", testStrHello, 10, testStrHello},
		{"exact", testStrHello, 5, testStrHello},
		{"truncate", testStrHelloWorld, 8, testStrHelloTrunc},
		{"very short max", testStrHello, 2, "he"},
		{"zero max", testStrHello, 0, testStrEmpty},
		{"negative max", testStrHello, -1, testStrEmpty},
		{"max equals 3", testStrHelloWorld, 3, "hel"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.want, TruncateString(tc.s, tc.maxLen), "TruncateString(%q, %d)", tc.s, tc.maxLen)
		})
	}
}

func TestCoalesceString(t *testing.T) {
	testCases := []struct {
		name   string
		values []string
		want   string
	}{
		{"first non-empty", []string{testStrEmpty, testStrHello, testStrWorld}, testStrHello},
		{"first value", []string{testStrHello, testStrWorld}, testStrHello},
		{"all empty", []string{testStrEmpty, testStrEmpty}, testStrEmpty},
		{"no values", []string{}, testStrEmpty},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.want, CoalesceString(tc.values...), "CoalesceString(%v)", tc.values)
		})
	}
}

func TestStringOrDefault(t *testing.T) {
	testCases := []struct {
		name         string
		value        string
		defaultValue string
		want         string
	}{
		{"non-empty", testStrHello, testStrDefault, testStrHello},
		{"empty", testStrEmpty, testStrDefault, testStrDefault},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.want, StringOrDefault(tc.value, tc.defaultValue), "StringOrDefault(%q, %q)", tc.value, tc.defaultValue)
		})
	}
}

/* ======== REFLECTION UTILITIES TESTS ======== */

func TestIsZeroValue(t *testing.T) {
	// NOTE: IsZeroValue uses reflect; boxing concrete zero values into `any`
	// changes reflect.Kind, so each case must call the generic function
	// directly with its concrete type to avoid interface-boxing artifacts.
	type testStruct struct{ X int }

	t.Run("zero int", func(t *testing.T) {
		assertions := assert.New(t)
		assertions.True(IsZeroValue(0), "IsZeroValue(0)")
	})

	t.Run("non-zero int", func(t *testing.T) {
		assertions := assert.New(t)
		assertions.False(IsZeroValue(42), "IsZeroValue(42)")
	})

	t.Run("empty string", func(t *testing.T) {
		assertions := assert.New(t)
		assertions.True(IsZeroValue(testStrEmpty), "IsZeroValue(empty string)")
	})

	t.Run("non-empty string", func(t *testing.T) {
		assertions := assert.New(t)
		assertions.False(IsZeroValue(testStrHello), "IsZeroValue(hello)")
	})

	t.Run("nil pointer", func(t *testing.T) {
		assertions := assert.New(t)
		var p *int
		assertions.True(IsZeroValue(p), "IsZeroValue(nil pointer)")
	})

	t.Run("non-nil pointer", func(t *testing.T) {
		assertions := assert.New(t)
		val := 42
		assertions.False(IsZeroValue(&val), "IsZeroValue(non-nil pointer)")
	})

	t.Run("zero struct", func(t *testing.T) {
		assertions := assert.New(t)
		assertions.True(IsZeroValue(testStruct{}), "IsZeroValue(zero struct)")
	})

	t.Run("non-zero struct", func(t *testing.T) {
		assertions := assert.New(t)
		assertions.False(IsZeroValue(testStruct{X: 1}), "IsZeroValue(non-zero struct)")
	})

	t.Run("zero float64", func(t *testing.T) {
		assertions := assert.New(t)
		assertions.True(IsZeroValue(0.0), "IsZeroValue(0.0)")
	})

	t.Run("non-zero float64", func(t *testing.T) {
		assertions := assert.New(t)
		assertions.False(IsZeroValue(3.14), "IsZeroValue(3.14)")
	})

	t.Run("zero bool", func(t *testing.T) {
		assertions := assert.New(t)
		assertions.True(IsZeroValue(false), "IsZeroValue(false)")
	})

	t.Run("non-zero bool", func(t *testing.T) {
		assertions := assert.New(t)
		assertions.False(IsZeroValue(true), "IsZeroValue(true)")
	})
}

func TestIsNil(t *testing.T) {
	var nilPtr *int
	var nilMap map[string]int
	var nilSlice []int
	var nilChan chan int
	var nilFunc func()
	var nilInterface error

	testCases := []struct {
		name  string
		value any
		want  bool
	}{
		{"nil", nil, true},
		{"nil pointer", nilPtr, true},
		{"nil map", nilMap, true},
		{"nil slice", nilSlice, true},
		{"nil chan", nilChan, true},
		{"nil func", nilFunc, true},
		{"nil interface", nilInterface, true},
		{"non-nil int", 42, false},
		{"non-nil string", testStrHello, false},
		{"non-nil map", map[string]int{"a": 1}, false},
		{"non-nil slice", []int{1}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.want, IsNil(tc.value), "IsNil(%v)", tc.value)
		})
	}
}

/* ======== TIME UTILITIES TESTS ======== */

func TestUnixMillis(t *testing.T) {
	assertions := assert.New(t)

	before := time.Now().UnixMilli()
	got := UnixMillis()
	after := time.Now().UnixMilli()

	assertions.GreaterOrEqual(got, before, "UnixMillis() should be >= before")
	assertions.LessOrEqual(got, after, "UnixMillis() should be <= after")
}

func TestUnixMicros(t *testing.T) {
	assertions := assert.New(t)

	before := time.Now().UnixMicro()
	got := UnixMicros()
	after := time.Now().UnixMicro()

	assertions.GreaterOrEqual(got, before, "UnixMicros() should be >= before")
	assertions.LessOrEqual(got, after, "UnixMicros() should be <= after")
}

func TestUnixNanos(t *testing.T) {
	assertions := assert.New(t)

	before := time.Now().UnixNano()
	got := UnixNanos()
	after := time.Now().UnixNano()

	assertions.GreaterOrEqual(got, before, "UnixNanos() should be >= before")
	assertions.LessOrEqual(got, after, "UnixNanos() should be <= after")
}

func TestTimestampString(t *testing.T) {
	assertions := assert.New(t)

	ts := TimestampString()
	_, err := time.Parse(time.RFC3339, ts)
	assertions.NoError(err, "TimestampString() should produce valid RFC3339")
}

func TestDurationString(t *testing.T) {
	testCases := []struct {
		name string
		d    time.Duration
	}{
		{"nanoseconds", 500 * time.Nanosecond},
		{"microseconds", 500 * time.Microsecond},
		{"milliseconds", 500 * time.Millisecond},
		{"seconds", 5 * time.Second},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			got := DurationString(tc.d)
			assertions.NotEmpty(got, "DurationString(%v) should not be empty", tc.d)
		})
	}
}

/* ======== SLICE UTILITIES TESTS ======== */

func TestContains(t *testing.T) {
	testCases := []struct {
		name  string
		slice []int
		elem  int
		want  bool
	}{
		{"found", []int{1, 2, 3, 4, 5}, 3, true},
		{"not found", []int{1, 2, 3, 4, 5}, 10, false},
		{"empty slice", []int{}, 1, false},
		{"single element found", []int{42}, 42, true},
		{"single element not found", []int{42}, 7, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.want, Contains(tc.slice, tc.elem))
		})
	}
}

func TestContainsStrings(t *testing.T) {
	assertions := assert.New(t)
	slice := []string{testStrHello, testStrWorld}

	assertions.True(Contains(slice, testStrHello))
	assertions.False(Contains(slice, testStrDefault))
}

func TestUnique(t *testing.T) {
	testCases := []struct {
		name  string
		slice []int
		want  []int
	}{
		{"with duplicates", []int{1, 2, 2, 3, 3, 3}, []int{1, 2, 3}},
		{"no duplicates", []int{1, 2, 3}, []int{1, 2, 3}},
		{"empty", []int{}, nil},
		{"nil-slice", nil, nil},
		{"single element", []int{42}, []int{42}},
		{"all same", []int{5, 5, 5, 5}, []int{5}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			got := Unique(tc.slice)
			assertions.Equal(tc.want, got, "Unique(%v)", tc.slice)
		})
	}
}

func TestFilter(t *testing.T) {
	testCases := []struct {
		name      string
		slice     []int
		predicate func(int) bool
		want      []int
	}{
		{"even numbers", []int{1, 2, 3, 4, 5}, func(n int) bool { return n%2 == 0 }, []int{2, 4}},
		{"empty slice", []int{}, func(n int) bool { return true }, nil},
		{"nil_slice", nil, func(n int) bool { return true }, nil},
		{"none match", []int{1, 3, 5}, func(n int) bool { return n%2 == 0 }, []int{}},
		{"all match", []int{2, 4, 6}, func(n int) bool { return n%2 == 0 }, []int{2, 4, 6}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			got := Filter(tc.slice, tc.predicate)
			assertions.Equal(tc.want, got)
		})
	}
}

func TestMap(t *testing.T) {
	testCases := []struct {
		name  string
		slice []int
		want  []int
	}{
		{"double", []int{1, 2, 3}, []int{2, 4, 6}},
		{"empty", []int{}, nil},
		{"nil", nil, nil},
		{"single", []int{5}, []int{10}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			got := Map(tc.slice, func(n int) int { return n * 2 })
			assertions.Equal(tc.want, got)
		})
	}
}

func TestMapTypeConversion(t *testing.T) {
	assertions := assert.New(t)
	slice := []int{1, 2, 3}
	got := Map(slice, func(n int) string {
		return string(rune('a' + n - 1))
	})
	assertions.Equal([]string{"a", "b", "c"}, got)
}

/* ======== MAP UTILITIES TESTS ======== */

func TestMergeMaps(t *testing.T) {
	testCases := []struct {
		name string
		maps []map[string]int
		want map[string]int
	}{
		{
			"two maps with overlap",
			[]map[string]int{{"a": 1, "b": 2}, {"b": 3, "c": 4}},
			map[string]int{"a": 1, "b": 3, "c": 4},
		},
		{
			"empty input",
			[]map[string]int{},
			map[string]int{},
		},
		{
			"single map",
			[]map[string]int{{"a": 1}},
			map[string]int{"a": 1},
		},
		{
			"nil map in list",
			[]map[string]int{{"a": 1}, nil, {"b": 2}},
			map[string]int{"a": 1, "b": 2},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			got := MergeMaps(tc.maps...)
			assertions.Equal(tc.want, got)
		})
	}
}

func TestCopyMap(t *testing.T) {
	testCases := []struct {
		name  string
		input map[string]int
		isNil bool
	}{
		{"normal map", map[string]int{"a": 1, "b": 2}, false},
		{"nil-map", nil, true},
		{"empty map", map[string]int{}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			got := CopyMap(tc.input)
			if tc.isNil {
				assertions.Nil(got)
			} else {
				assertions.Equal(tc.input, got)
				// Verify independence
				if len(got) > 0 {
					got["zzz"] = 999
					_, exists := tc.input["zzz"]
					assertions.False(exists, "CopyMap should create independent copy")
				}
			}
		})
	}
}

func TestKeys(t *testing.T) {
	testCases := []struct {
		name    string
		input   map[string]int
		wantLen int
		isNil   bool
	}{
		{"normal-map", map[string]int{"a": 1, "b": 2}, 2, false},
		{"empty-map", map[string]int{}, 0, true},
		{"nil_map", nil, 0, true},
		{"single entry", map[string]int{"x": 42}, 1, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			got := Keys(tc.input)
			if tc.isNil {
				assertions.Nil(got)
			} else {
				assertions.Len(got, tc.wantLen)
			}
		})
	}
}

func TestValues(t *testing.T) {
	testCases := []struct {
		name    string
		input   map[string]int
		wantLen int
		isNil   bool
	}{
		{"normal_map", map[string]int{"a": 1, "b": 2}, 2, false},
		{"empty_map", map[string]int{}, 0, true},
		{"null map", nil, 0, true},
		{"single entry", map[string]int{"x": 42}, 1, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			got := Values(tc.input)
			if tc.isNil {
				assertions.Nil(got)
			} else {
				assertions.Len(got, tc.wantLen)
			}
		})
	}
}

/* ======== TYPES TESTS — Optional ======== */

func TestOptional(t *testing.T) {
	t.Run("Some", func(t *testing.T) {
		assertions := assert.New(t)
		opt := Some(42)
		assertions.True(opt.IsPresent(), "Some should be present")
		assertions.False(opt.IsEmpty(), "Some should not be empty")
		assertions.Equal(42, opt.Value(), "Value() should return 42")

		v, ok := opt.Get()
		assertions.True(ok, "Get() should return true")
		assertions.Equal(42, v, "Get() value should be 42")

		assertions.Equal(42, opt.ValueOr(99), "ValueOr should return actual value when present")
	})

	t.Run("None", func(t *testing.T) {
		assertions := assert.New(t)
		opt := None[int]()
		assertions.False(opt.IsPresent(), "None should not be present")
		assertions.True(opt.IsEmpty(), "None should be empty")
		assertions.Equal(0, opt.Value(), "Value() on None should return zero")
		assertions.Equal(10, opt.ValueOr(10), "ValueOr(10) should return 10")

		v, ok := opt.Get()
		assertions.False(ok, "Get() should return false")
		assertions.Equal(0, v, "Get() value on None should be zero")
	})

	t.Run("Some with string", func(t *testing.T) {
		assertions := assert.New(t)
		opt := Some(testStrHello)
		assertions.True(opt.IsPresent())
		assertions.Equal(testStrHello, opt.Value())
	})

	t.Run("None with string", func(t *testing.T) {
		assertions := assert.New(t)
		opt := None[string]()
		assertions.Equal(testStrDefault, opt.ValueOr(testStrDefault))
	})
}

/* ======== TYPES TESTS — Result ======== */

func TestResult(t *testing.T) {
	t.Run("Ok", func(t *testing.T) {
		assertions := assert.New(t)
		result := Ok(42)
		assertions.True(result.IsOk(), "Ok should be successful")
		assertions.False(result.IsErr(), "Ok should not have error")
		assertions.Equal(42, result.Value(), "Value() should return 42")
		assertions.Nil(result.Error(), "Error() should be nil on Ok")

		v, err := result.Get()
		assertions.NoError(err, "Get() error should be nil")
		assertions.Equal(42, v, "Get() value should be 42")

		assertions.Equal(42, result.ValueOr(99), "ValueOr on Ok should return actual value")
	})

	t.Run("Err", func(t *testing.T) {
		assertions := assert.New(t)
		testErr := errors.New(testMsgTestErr)
		result := Err[int](testErr)
		assertions.True(result.IsErr(), "Err should have error")
		assertions.False(result.IsOk(), "Err should not be ok")
		assertions.Equal(testErr, result.Error(), "Error() should return the error")
		assertions.Equal(10, result.ValueOr(10), "ValueOr(10) should return 10")

		v, err := result.Get()
		assertions.Error(err, "Get() should return error")
		assertions.Equal(0, v, "Get() value on Err should be zero")
	})
}

/* ======== TYPES TESTS — Pair ======== */

func TestNewPair(t *testing.T) {
	assertions := assert.New(t)

	p := NewPair("key", 42)
	assertions.Equal("key", p.Key)
	assertions.Equal(42, p.Value)

	p2 := NewPair(1, testStrHello)
	assertions.Equal(1, p2.Key)
	assertions.Equal(testStrHello, p2.Value)
}

/* ======== TYPES TESTS — SafeMap ======== */

func TestSafeMap(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSafeMap[string, int]()

		m.Set("a", 1)
		m.Set("b", 2)

		v, ok := m.Get("a")
		assertions.True(ok, "Get(a) should exist")
		assertions.Equal(1, v, "Get(a) value")
		assertions.Equal(2, m.Len(), "Len()")

		m.Delete("a")
		assertions.False(m.Has("a"), "Has(a) should be false after delete")
	})

	t.Run("with capacity", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSafeMapWithCapacity[string, int](10)
		m.Set("x", 99)
		v, ok := m.Get("x")
		assertions.True(ok)
		assertions.Equal(99, v)
	})

	t.Run("GetOrSet existing", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSafeMap[string, int]()
		m.Set("key", 10)

		v, existed := m.GetOrSet("key", 20)
		assertions.True(existed, "key should already exist")
		assertions.Equal(10, v, "should return existing value")
	})

	t.Run("GetOrSet new", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSafeMap[string, int]()

		v, existed := m.GetOrSet("key", 20)
		assertions.False(existed, "key should not exist yet")
		assertions.Equal(20, v, "should return newly set value")
	})

	t.Run("Clear", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSafeMap[string, int]()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Clear()
		assertions.Equal(0, m.Len())
		assertions.False(m.Has("a"))
	})

	t.Run("Keys", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSafeMap[string, int]()
		m.Set("a", 1)
		m.Set("b", 2)
		keys := m.Keys()
		assertions.Len(keys, 2)
		assertions.True(Contains(keys, "a"))
		assertions.True(Contains(keys, "b"))
	})

	t.Run("Values", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSafeMap[string, int]()
		m.Set("a", 1)
		m.Set("b", 2)
		vals := m.Values()
		assertions.Len(vals, 2)
		assertions.True(Contains(vals, 1))
		assertions.True(Contains(vals, 2))
	})

	t.Run("Range", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSafeMap[string, int]()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		visited := 0
		m.Range(func(key string, value int) bool {
			visited++
			return true
		})
		assertions.Equal(3, visited)
	})

	t.Run("Range early exit", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSafeMap[string, int]()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		visited := 0
		m.Range(func(key string, value int) bool {
			visited++
			return false // stop after first
		})
		assertions.Equal(1, visited)
	})

	t.Run("Copy", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSafeMap[string, int]()
		m.Set("a", 1)
		m.Set("b", 2)

		copied := m.Copy()
		assertions.Equal(2, copied.Len())
		v, ok := copied.Get("a")
		assertions.True(ok)
		assertions.Equal(1, v)

		// Verify independence
		copied.Set("c", 3)
		assertions.False(m.Has("c"))
	})

	t.Run("Get non-existent key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSafeMap[string, int]()
		v, ok := m.Get("missing")
		assertions.False(ok)
		assertions.Equal(0, v)
	})
}

/* ======== TYPES TESTS — SafeSlice ======== */

func TestSafeSlice(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		assertions := assert.New(t)
		s := NewSafeSlice[int]()

		s.Append(1, 2, 3)
		assertions.Equal(3, s.Len(), "Len()")

		v, ok := s.Get(1)
		assertions.True(ok, "Get(1) should exist")
		assertions.Equal(2, v, "Get(1) value")

		v, ok = s.Pop()
		assertions.True(ok, "Pop() should succeed")
		assertions.Equal(3, v, "Pop() value")
		assertions.Equal(2, s.Len(), "Len() after Pop()")
	})

	t.Run("with capacity", func(t *testing.T) {
		assertions := assert.New(t)
		s := NewSafeSliceWithCapacity[string](10)
		s.Append(testStrHello)
		assertions.Equal(1, s.Len())
		v, ok := s.Get(0)
		assertions.True(ok)
		assertions.Equal(testStrHello, v)
	})

	t.Run("Get out of bounds", func(t *testing.T) {
		assertions := assert.New(t)
		s := NewSafeSlice[int]()
		s.Append(1, 2)

		_, ok := s.Get(-1)
		assertions.False(ok, "Get(-1) should fail")

		_, ok = s.Get(5)
		assertions.False(ok, "Get(5) should fail")
	})

	t.Run("Set valid index", func(t *testing.T) {
		assertions := assert.New(t)
		s := NewSafeSlice[int]()
		s.Append(1, 2, 3)

		ok := s.Set(1, 99)
		assertions.True(ok, "Set(1, 99) should succeed")
		v, _ := s.Get(1)
		assertions.Equal(99, v)
	})

	t.Run("Set out of bounds", func(t *testing.T) {
		assertions := assert.New(t)
		s := NewSafeSlice[int]()
		s.Append(1, 2)

		ok := s.Set(-1, 99)
		assertions.False(ok, "Set(-1) should fail")

		ok = s.Set(5, 99)
		assertions.False(ok, "Set(5) should fail")
	})

	t.Run("Clear", func(t *testing.T) {
		assertions := assert.New(t)
		s := NewSafeSlice[int]()
		s.Append(1, 2, 3)
		s.Clear()
		assertions.Equal(0, s.Len())
	})

	t.Run("All", func(t *testing.T) {
		assertions := assert.New(t)
		s := NewSafeSlice[int]()
		s.Append(1, 2, 3)

		all := s.All()
		assertions.Equal([]int{1, 2, 3}, all)

		// Verify it is a copy
		all[0] = 999
		v, _ := s.Get(0)
		assertions.Equal(1, v, "All() should return a copy")
	})

	t.Run("Range", func(t *testing.T) {
		assertions := assert.New(t)
		s := NewSafeSlice[int]()
		s.Append(10, 20, 30)

		sum := 0
		s.Range(func(index int, value int) bool {
			sum += value
			return true
		})
		assertions.Equal(60, sum)
	})

	t.Run("Range early exit", func(t *testing.T) {
		assertions := assert.New(t)
		s := NewSafeSlice[int]()
		s.Append(10, 20, 30)

		visited := 0
		s.Range(func(index int, value int) bool {
			visited++
			return false
		})
		assertions.Equal(1, visited)
	})

	t.Run("Pop empty", func(t *testing.T) {
		assertions := assert.New(t)
		s := NewSafeSlice[int]()
		_, ok := s.Pop()
		assertions.False(ok, "Pop() on empty should return false")
	})
}

/* ======== ERROR TESTS ======== */

func TestWrapError(t *testing.T) {
	testCases := []struct {
		name    string
		err     error
		msg     string
		isNil   bool
		wantMsg string
	}{
		{"non-nil error", errors.New(testMsgOriginalErr), testMsgContext, false, "context: original error"},
		{"nil error", nil, testMsgContext, true, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			wrapped := WrapError(tc.err, tc.msg)
			if tc.isNil {
				assertions.Nil(wrapped)
			} else {
				assertions.NotNil(wrapped)
				assertions.Equal(tc.wantMsg, wrapped.Error())
				assertions.True(errors.Is(wrapped, tc.err))
			}
		})
	}
}

func TestWrapErrorf(t *testing.T) {
	testCases := []struct {
		name    string
		err     error
		format  string
		args    []any
		isNil   bool
		wantMsg string
	}{
		{"non-nil error", errors.New(testMsgOriginalErr), testMsgFmtWrap, []any{testOpPublish}, false, testMsgFmtResult},
		{"nil-error", nil, testMsgFmtWrap, []any{testOpPublish}, true, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			wrapped := WrapErrorf(tc.err, tc.format, tc.args...)
			if tc.isNil {
				assertions.Nil(wrapped)
			} else {
				assertions.NotNil(wrapped)
				assertions.Equal(tc.wantMsg, wrapped.Error())
			}
		})
	}
}

func TestFormatError(t *testing.T) {
	testCases := []struct {
		name string
		op   string
		err  error
		want string
	}{
		{"with error", testOpPublish, errors.New(testMsgOriginalErr), testMsgFmtErr},
		{"nil_error", testOpPublish, nil, testMsgSuccess},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.want, FormatError(tc.op, tc.err))
		})
	}
}

func TestIs(t *testing.T) {
	assertions := assert.New(t)

	wrapped := WrapError(ErrNotConnected, testMsgContext)
	assertions.True(Is(wrapped, ErrNotConnected))
	assertions.False(Is(wrapped, ErrConnectionClosed))
	assertions.False(Is(nil, ErrNotConnected))
}

func TestAs(t *testing.T) {
	assertions := assert.New(t)

	serErr := SerializationError{Operation: testOpSerialize, Type: testTypeMyStruct, Err: errors.New(testMsgTestErr)}
	var target SerializationError
	assertions.True(As(serErr, &target))
	assertions.Equal(testOpSerialize, target.Operation)

	assertions.False(As(errors.New(testMsgTestErr), &target))
}

func TestIsAnyOf(t *testing.T) {
	testCases := []struct {
		name    string
		err     error
		targets []error
		want    bool
	}{
		{"matches first", ErrNotConnected, []error{ErrNotConnected, ErrConnectionClosed}, true},
		{"matches second", ErrConnectionClosed, []error{ErrNotConnected, ErrConnectionClosed}, true},
		{"no match", ErrPublishFailed, []error{ErrNotConnected, ErrConnectionClosed}, false},
		{"a nil error", nil, []error{ErrNotConnected}, false},
		{"empty targets", ErrNotConnected, []error{}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tc.want, IsAnyOf(tc.err, tc.targets...))
		})
	}
}

func TestErrorChain(t *testing.T) {
	testCases := []struct {
		name    string
		err     error
		wantLen int
	}{
		{"nil", nil, 0},
		{"single", errors.New(testMsgRootErr), 1},
		{"wrapped twice", WrapError(WrapError(errors.New(testMsgRootErr), testMsgLevel1), testMsgLevel2), 3},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			chain := ErrorChain(tc.err)
			assertions.Len(chain, tc.wantLen)
		})
	}
}

func TestRootCause(t *testing.T) {
	testCases := []struct {
		name string
		err  error
		want error
	}{
		{"nil", nil, nil},
		{"single", errors.New(testMsgRootErr), errors.New(testMsgRootErr)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			cause := RootCause(tc.err)
			if tc.want == nil {
				assertions.Nil(cause)
			} else {
				assertions.Equal(tc.want.Error(), cause.Error())
			}
		})
	}

	t.Run("wrapped chain", func(t *testing.T) {
		assertions := assert.New(t)
		root := errors.New(testMsgRootErr)
		wrapped := WrapError(WrapError(root, testMsgLevel1), testMsgLevel2)
		assertions.Equal(root, RootCause(wrapped))
	})
}

func TestErrorCollector(t *testing.T) {
	t.Run("basic usage", func(t *testing.T) {
		assertions := assert.New(t)
		ec := NewErrorCollector()

		assertions.False(ec.HasErrors())
		assertions.Equal(0, ec.Count())
		assertions.Nil(ec.Error())

		ec.Add(errors.New(testMsgErr1))
		ec.Add(nil)
		ec.Add(errors.New(testMsgErr2))

		assertions.True(ec.HasErrors())
		assertions.Equal(2, ec.Count())
		assertions.NotNil(ec.Error())
		assertions.Len(ec.All(), 2)
	})

	t.Run("single error returns error directly", func(t *testing.T) {
		assertions := assert.New(t)
		ec := NewErrorCollector()
		err1 := errors.New(testMsgErr1)
		ec.Add(err1)
		assertions.Equal(err1, ec.Error())
	})

	t.Run("multiple errors returns combined", func(t *testing.T) {
		assertions := assert.New(t)
		ec := NewErrorCollector()
		ec.Add(errors.New(testMsgErr1))
		ec.Add(errors.New(testMsgErr2))
		ec.Add(errors.New(testMsgErr3))

		combined := ec.Error()
		assertions.NotNil(combined)
		assertions.Contains(combined.Error(), "3 errors occurred")
	})

	t.Run("AddAll", func(t *testing.T) {
		assertions := assert.New(t)
		ec := NewErrorCollector()
		ec.AddAll(errors.New(testMsgErr1), nil, errors.New(testMsgErr2), nil)
		assertions.Equal(2, ec.Count())
	})
}

func TestValidationError(t *testing.T) {
	t.Run("with field", func(t *testing.T) {
		assertions := assert.New(t)
		err := NewValidationError(testFieldEmail, testMsgInvalidFmt)
		assertions.Equal(testExpectedValErr, err.Error())
	})

	t.Run("without field", func(t *testing.T) {
		assertions := assert.New(t)
		err := NewValidationError("", testMsgRequired)
		assertions.Equal("validation error: required", err.Error())
	})

	t.Run("with value", func(t *testing.T) {
		assertions := assert.New(t)
		err := NewValidationErrorWithValue(testFieldEmail, testMsgInvalidFmt, "bad@")
		assertions.Equal(testExpectedValErr, err.Error())
		assertions.Equal("bad@", err.Value)
	})
}

func TestRecoverError(t *testing.T) {
	testCases := []struct {
		name    string
		input   any
		isNil   bool
		wantMsg string
	}{
		{"nil", nil, true, ""},
		{"error type", errors.New(testMsgOriginalErr), false, testMsgOriginalErr},
		{"string type", testMsgPanicStr, false, testMsgPanicStr},
		{"int type", 42, false, testMsgPanicInt},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assertions := assert.New(t)
			err := RecoverError(tc.input)
			if tc.isNil {
				assertions.Nil(err)
			} else {
				assertions.Equal(tc.wantMsg, err.Error())
			}
		})
	}
}

func TestSafeGo(t *testing.T) {
	t.Run("successful function", func(t *testing.T) {
		assertions := assert.New(t)
		errChan := make(chan error, 1)
		SafeGo(func() error { return nil }, errChan)
		// Give goroutine time to complete
		select {
		case err := <-errChan:
			assertions.Fail("should not receive error", err.Error())
		case <-time.After(100 * time.Millisecond):
			// Expected: no error
		}
	})

	t.Run("function returns error", func(t *testing.T) {
		assertions := assert.New(t)
		errChan := make(chan error, 1)
		testErr := errors.New(testMsgTestErr)
		SafeGo(func() error { return testErr }, errChan)

		select {
		case err := <-errChan:
			assertions.Equal(testErr, err)
		case <-time.After(1 * time.Second):
			assertions.Fail("should have received error")
		}
	})

	t.Run("function panics with error", func(t *testing.T) {
		assertions := assert.New(t)
		errChan := make(chan error, 1)
		SafeGo(func() error { panic(errors.New(testMsgPanicStr)) }, errChan)

		select {
		case err := <-errChan:
			assertions.Equal(testMsgPanicStr, err.Error())
		case <-time.After(1 * time.Second):
			assertions.Fail("should have received panic error")
		}
	})

	t.Run("function panics with string", func(t *testing.T) {
		assertions := assert.New(t)
		errChan := make(chan error, 1)
		SafeGo(func() error { panic(testMsgPanicStr) }, errChan)

		select {
		case err := <-errChan:
			assertions.Equal(testMsgPanicStr, err.Error())
		case <-time.After(1 * time.Second):
			assertions.Fail("should have received panic error")
		}
	})

	t.Run("nil errChan with error", func(t *testing.T) {
		// Should not panic
		SafeGo(func() error { return errors.New(testMsgTestErr) }, nil)
		time.Sleep(50 * time.Millisecond)
	})

	t.Run("nil errChan with panic", func(t *testing.T) {
		// Should not panic
		SafeGo(func() error { panic(testMsgPanicStr) }, nil)
		time.Sleep(50 * time.Millisecond)
	})
}

func TestCombinedErrorString(t *testing.T) {
	assertions := assert.New(t)
	ec := NewErrorCollector()
	ec.Add(errors.New(testMsgErr1))
	ec.Add(errors.New(testMsgErr2))

	combined := ec.Error()
	assertions.Contains(combined.Error(), "2 errors occurred")
	assertions.Contains(combined.Error(), testMsgErr1)
	assertions.Contains(combined.Error(), testMsgErr2)
}

/* ======== BENCHMARKS ======== */

func BenchmarkValidateSubject(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ValidateSubject(testSubjectDotted)
		}
	})
}

func BenchmarkBuildSubject(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			BuildSubject(testSubjectEvents, testTokenUser, testTokenCreated)
		}
	})
}

func BenchmarkGenerateID(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			GenerateID(16)
		}
	})
}

func BenchmarkGenerateMessageID(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			GenerateMessageID()
		}
	})
}

func BenchmarkTruncateString(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			TruncateString(testStrHelloWorld, 8)
		}
	})
}

func BenchmarkContains(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Contains(slice, 7)
		}
	})
}

func BenchmarkUnique(b *testing.B) {
	slice := []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Unique(slice)
		}
	})
}

func BenchmarkMergeMaps(b *testing.B) {
	m1 := map[string]int{"a": 1, "b": 2, "c": 3}
	m2 := map[string]int{"b": 20, "d": 4, "e": 5}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			MergeMaps(m1, m2)
		}
	})
}

func BenchmarkIsNil(b *testing.B) {
	var nilPtr *int
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			IsNil(nilPtr)
		}
	})
}

func BenchmarkSafeMapOperations(b *testing.B) {
	m := NewSafeMap[string, int]()
	m.Set("key", 42)
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m.Set("key", 42)
			m.Get("key")
			m.Has("key")
			m.Len()
		}
	})
}

func BenchmarkSafeSliceOperations(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s := NewSafeSlice[int]()
			s.Append(1, 2, 3)
			s.Get(1)
			s.Len()
		}
	})
}

func BenchmarkFilter(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Filter(slice, func(n int) bool { return n%2 == 0 })
		}
	})
}

func BenchmarkRecoverError(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			RecoverError(testMsgPanicStr)
		}
	})
}
