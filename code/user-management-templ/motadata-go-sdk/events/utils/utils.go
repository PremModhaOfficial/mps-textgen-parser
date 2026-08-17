// Package utils provides utility functions for the events package.
//
// This package contains common helper functions used across the events
// subsystem, including subject validation, ID generation, and string utilities.
package utils

import (
	"crypto/rand"
	"encoding/hex"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// ============================================================================
// Subject Validation
// ============================================================================

// subjectPattern validates NATS subject names
// Allowed characters: alphanumeric, dots, asterisks, greater-than
var subjectPattern = regexp.MustCompile(`^[a-zA-Z0-9._*>-]+$`)

// ValidateSubject checks if a subject name is valid.
// Returns true if the subject follows NATS naming conventions.
func ValidateSubject(subject string) bool {
	if subject == "" {
		return false
	}
	if len(subject) > 256 {
		return false
	}
	return subjectPattern.MatchString(subject)
}

// SanitizeSubject removes invalid characters from a subject name.
// Returns a sanitized subject that conforms to NATS naming conventions.
func SanitizeSubject(subject string) string {
	if subject == "" {
		return ""
	}

	var builder strings.Builder
	builder.Grow(len(subject))

	for _, r := range subject {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '.' || r == '_' || r == '-' || r == '*' || r == '>' {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

// BuildSubject joins subject tokens with dots.
// Example: BuildSubject("events", "user", "created") returns "events.user.created"
func BuildSubject(tokens ...string) string {
	if len(tokens) == 0 {
		return ""
	}

	var nonEmpty []string
	for _, token := range tokens {
		if token != "" {
			nonEmpty = append(nonEmpty, token)
		}
	}

	return strings.Join(nonEmpty, ".")
}

// ParseSubject splits a subject into its component tokens.
// Example: ParseSubject("events.user.created") returns ["events", "user", "created"]
func ParseSubject(subject string) []string {
	if subject == "" {
		return nil
	}
	return strings.Split(subject, ".")
}

// ============================================================================
// ID Generation
// ============================================================================

// GenerateID generates a random ID with the specified byte length.
// The returned string is hex-encoded, so it will be twice the byte length.
func GenerateID(byteLength int) string {
	if byteLength <= 0 {
		byteLength = 16
	}

	bytes := make([]byte, byteLength)
	_, _ = rand.Read(bytes)

	return hex.EncodeToString(bytes)
}

// GenerateMessageID generates a unique message ID suitable for deduplication.
// Returns a 32-character hex string (16 bytes).
func GenerateMessageID() string {
	return GenerateID(16)
}

// GenerateCorrelationID generates a correlation ID for request-reply patterns.
// Returns a 24-character hex string (12 bytes).
func GenerateCorrelationID() string {
	return GenerateID(12)
}

// GenerateTraceID generates a trace ID for distributed tracing.
// Returns a 32-character hex string (16 bytes) conforming to W3C trace context.
func GenerateTraceID() string {
	return GenerateID(16)
}

// GenerateSpanID generates a span ID for distributed tracing.
// Returns a 16-character hex string (8 bytes) conforming to W3C trace context.
func GenerateSpanID() string {
	return GenerateID(8)
}

// ============================================================================
// String Utilities
// ============================================================================

// TruncateString truncates a string to the specified maximum length.
// If the string is longer, it adds "..." suffix.
func TruncateString(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

// CoalesceString returns the first non-empty string from the arguments.
func CoalesceString(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

// StringOrDefault returns the value if non-empty, otherwise returns the default.
func StringOrDefault(value, defaultValue string) string {
	if value != "" {
		return value
	}
	return defaultValue
}

// ============================================================================
// Reflection Utilities
// ============================================================================

// IsZeroValue checks if a value is its zero value.
// Works with any type including structs, pointers, and interfaces.
func IsZeroValue[T any](v T) bool {
	return reflect.ValueOf(&v).Elem().IsZero()
}

// IsNil checks if an interface value is nil.
// Handles the case where the interface contains a nil pointer.
func IsNil(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Chan, reflect.Func, reflect.Interface:
		return rv.IsNil()
	}

	return false
}

// ============================================================================
// Time Utilities
// ============================================================================

// UnixMillis returns the current Unix timestamp in milliseconds.
func UnixMillis() int64 {
	return time.Now().UnixMilli()
}

// UnixMicros returns the current Unix timestamp in microseconds.
func UnixMicros() int64 {
	return time.Now().UnixMicro()
}

// UnixNanos returns the current Unix timestamp in nanoseconds.
func UnixNanos() int64 {
	return time.Now().UnixNano()
}

// TimestampString returns the current timestamp as an RFC3339 string.
func TimestampString() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// DurationString formats a duration in a human-readable format.
func DurationString(d time.Duration) string {
	if d < time.Microsecond {
		return d.String()
	}
	if d < time.Millisecond {
		return d.Round(time.Microsecond).String()
	}
	if d < time.Second {
		return d.Round(time.Millisecond).String()
	}
	return d.Round(time.Second).String()
}

// ============================================================================
// Slice Utilities
// ============================================================================

// Contains checks if a slice contains a specific element.
func Contains[T comparable](slice []T, elem T) bool {
	for _, v := range slice {
		if v == elem {
			return true
		}
	}
	return false
}

// Unique returns a new slice with duplicate elements removed.
// Preserves the order of first occurrences.
func Unique[T comparable](slice []T) []T {
	if len(slice) == 0 {
		return nil
	}

	seen := make(map[T]struct{}, len(slice))
	result := make([]T, 0, len(slice))

	for _, v := range slice {
		if _, exists := seen[v]; !exists {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}

	return result
}

// Filter returns a new slice containing only elements that satisfy the predicate.
func Filter[T any](slice []T, predicate func(T) bool) []T {
	if len(slice) == 0 {
		return nil
	}

	result := make([]T, 0)
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}

	return result
}

// Map applies a transformation function to each element and returns a new slice.
func Map[T, U any](slice []T, transform func(T) U) []U {
	if len(slice) == 0 {
		return nil
	}

	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = transform(v)
	}

	return result
}

// ============================================================================
// Map Utilities
// ============================================================================

// MergeMaps merges multiple maps into a new map.
// Later maps take precedence over earlier ones for duplicate keys.
func MergeMaps[K comparable, V any](maps ...map[K]V) map[K]V {
	result := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// CopyMap creates a shallow copy of a map.
func CopyMap[K comparable, V any](m map[K]V) map[K]V {
	if m == nil {
		return nil
	}

	result := make(map[K]V, len(m))
	for k, v := range m {
		result[k] = v
	}

	return result
}

// Keys returns all keys from a map as a slice.
// The order of keys is not guaranteed.
func Keys[K comparable, V any](m map[K]V) []K {
	if len(m) == 0 {
		return nil
	}

	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

// Values returns all values from a map as a slice.
// The order of values is not guaranteed.
func Values[K comparable, V any](m map[K]V) []V {
	if len(m) == 0 {
		return nil
	}

	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}

	return values
}
