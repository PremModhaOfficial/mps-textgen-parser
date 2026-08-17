package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ============================================================================
// Generic Type Constraint Tests
// ============================================================================

func TestIsInteger(t *testing.T) {

	assertions := assert.New(t)

	// Test signed integers
	assertions.Equal(42, IsInteger(42))
	assertions.Equal(int8(42), IsInteger(int8(42)))
	assertions.Equal(int16(42), IsInteger(int16(42)))
	assertions.Equal(int32(42), IsInteger(int32(42)))
	assertions.Equal(int64(42), IsInteger(int64(42)))

	// Test unsigned integers
	assertions.Equal(uint(43), IsInteger(uint(43)))
	assertions.Equal(uint8(43), IsInteger(uint8(43)))
	assertions.Equal(uint16(43), IsInteger(uint16(43)))
	assertions.Equal(uint32(43), IsInteger(uint32(43)))
	assertions.Equal(uint64(43), IsInteger(uint64(43)))
	assertions.Equal(uintptr(43), IsInteger(uintptr(43)))

	// Test negative values for signed types
	assertions.Equal(-42, IsInteger(-42))
	assertions.Equal(int8(-42), IsInteger(int8(-42)))
	assertions.Equal(int16(-42), IsInteger(int16(-42)))
	assertions.Equal(int32(-42), IsInteger(int32(-42)))
	assertions.Equal(int64(-42), IsInteger(int64(-42)))

	// Test zero values
	assertions.Equal(0, IsInteger(0))
	assertions.Equal(uint(0), IsInteger(uint(0)))
}

func TestIsFloat(t *testing.T) {
	assertions := assert.New(t)

	// Test float32
	assertions.Equal(float32(3.14), IsFloat(float32(3.14)))
	assertions.Equal(float32(-3.14), IsFloat(float32(-3.14)))
	assertions.Equal(float32(0.0), IsFloat(float32(0.0)))

	// Test float64
	assertions.Equal(3.14159265, IsFloat(3.14159265))
	assertions.Equal(-3.14159265, IsFloat(-3.14159265))
	assertions.Equal(0.0, IsFloat(0.0))

	// Test special float values
	assertions.Equal(float32(1.23e-10), IsFloat(float32(1.23e-10)))
	assertions.Equal(1.23e308, IsFloat(1.23e308))
}

func TestIsString(t *testing.T) {
	assertions := assert.New(t)

	// Test string values
	assertions.Equal("hello", IsString("hello"))
	assertions.Equal("", IsString(""))
	assertions.Equal("hello world", IsString("hello world"))
	assertions.Equal("123", IsString("123"))
	assertions.Equal("special characters: !@#$%^&*()", IsString("special characters: !@#$%^&*()"))

	// Test strings with unicode
	assertions.Equal("こんにちは", IsString("こんにちは"))
	assertions.Equal("🚀", IsString("🚀"))
}

func TestIsBoolean(t *testing.T) {
	assertions := assert.New(t)

	// Test boolean values
	assertions.Equal(true, IsBoolean(true))
	assertions.Equal(false, IsBoolean(false))
}

func TestIsSigned(t *testing.T) {
	assertions := assert.New(t)

	// Test signed integer types
	assertions.Equal(42, IsSigned(42))
	assertions.Equal(int8(42), IsSigned(int8(42)))
	assertions.Equal(int16(42), IsSigned(int16(42)))
	assertions.Equal(int32(42), IsSigned(int32(42)))
	assertions.Equal(int64(42), IsSigned(int64(42)))

	// Test negative values
	assertions.Equal(-43, IsSigned(-43))
	assertions.Equal(int8(-43), IsSigned(int8(-43)))
	assertions.Equal(int16(-43), IsSigned(int16(-43)))
	assertions.Equal(int32(-43), IsSigned(int32(-43)))
	assertions.Equal(int64(-43), IsSigned(int64(-43)))

	// Test zero values
	assertions.Equal(0, IsSigned(0))
	assertions.Equal(int8(0), IsSigned(int8(0)))
	assertions.Equal(int16(0), IsSigned(int16(0)))
	assertions.Equal(int32(0), IsSigned(int32(0)))
	assertions.Equal(int64(0), IsSigned(int64(0)))
}

func TestIsUnsigned(t *testing.T) {
	assertions := assert.New(t)

	// Test unsigned integer types
	assertions.Equal(uint(42), IsUnsigned(uint(42)))
	assertions.Equal(uint8(42), IsUnsigned(uint8(42)))
	assertions.Equal(uint16(42), IsUnsigned(uint16(42)))
	assertions.Equal(uint32(42), IsUnsigned(uint32(42)))
	assertions.Equal(uint64(42), IsUnsigned(uint64(42)))
	assertions.Equal(uintptr(42), IsUnsigned(uintptr(42)))

	// Test zero values
	assertions.Equal(uint(0), IsUnsigned(uint(0)))
	assertions.Equal(uint8(0), IsUnsigned(uint8(0)))
	assertions.Equal(uint16(0), IsUnsigned(uint16(0)))
	assertions.Equal(uint32(0), IsUnsigned(uint32(0)))
	assertions.Equal(uint64(0), IsUnsigned(uint64(0)))
	assertions.Equal(uintptr(0), IsUnsigned(uintptr(0)))

	// Test maximum values
	assertions.Equal(uint8(255), IsUnsigned(uint8(255)))
	assertions.Equal(uint16(65535), IsUnsigned(uint16(65535)))
}

// ============================================================================
// Type Constraint Interface Tests
// ============================================================================

func TestTypeConstraintsCompileTime(t *testing.T) {
	assertions := assert.New(t)

	// Test that functions can be called with constraint-compatible types
	// This verifies that the constraints work properly at compile time

	// Test Signed constraint through IsSigned function
	assertions.Equal(42, IsSigned(42))
	assertions.Equal(int8(-10), IsSigned(int8(-10)))
	assertions.Equal(int16(1000), IsSigned(int16(1000)))
	assertions.Equal(int32(-50000), IsSigned(int32(-50000)))
	assertions.Equal(int64(9223372036854775807), IsSigned(int64(9223372036854775807)))

	// Test Unsigned constraint through IsUnsigned function
	assertions.Equal(uint(42), IsUnsigned(uint(42)))
	assertions.Equal(uint8(255), IsUnsigned(uint8(255)))
	assertions.Equal(uint16(65535), IsUnsigned(uint16(65535)))
	assertions.Equal(uint32(4294967295), IsUnsigned(uint32(4294967295)))
	assertions.Equal(uint64(18446744073709551615), IsUnsigned(uint64(18446744073709551615)))
	assertions.Equal(uintptr(123456), IsUnsigned(uintptr(123456)))

	// Test Integer constraint through IsInteger function (both signed and unsigned)
	assertions.Equal(42, IsInteger(42))
	assertions.Equal(int8(-10), IsInteger(int8(-10)))
	assertions.Equal(uint(42), IsInteger(uint(42)))
	assertions.Equal(uint8(255), IsInteger(uint8(255)))

	// Test Float constraint through IsFloat function
	assertions.Equal(float32(3.14), IsFloat(float32(3.14)))
	assertions.Equal(3.14159265, IsFloat(3.14159265))
	assertions.Equal(float32(-3.14), IsFloat(float32(-3.14)))
	assertions.Equal(-3.14159265, IsFloat(-3.14159265))

	// Test string through IsString function
	assertions.Equal("hello", IsString("hello"))
	assertions.Equal("", IsString(""))
	assertions.Equal("🚀", IsString("🚀"))
	assertions.Equal("123", IsString("123"))
}

// ============================================================================
// Integration Tests
// ============================================================================

func TestGenericConstraintIntegration(t *testing.T) {
	assertions := assert.New(t)

	// Test combining multiple constraint functions
	intVal := 42
	floatVal := 3.14
	stringVal := "hello"

	// Verify that constraint functions return expected values
	assertions.Equal(intVal, IsInteger(intVal))
	assertions.Equal(floatVal, IsFloat(floatVal))
	assertions.Equal(stringVal, IsString(stringVal))
	assertions.Equal(true, IsBoolean(true))

	// Test that the constraint functions preserve type information
	assertions.IsType(0, IsInteger(intVal))
	assertions.IsType(float64(0), IsFloat(floatVal))
	assertions.IsType("", IsString(stringVal))
	assertions.IsType(false, IsBoolean(true))
}

func TestConstraintEdgeCases(t *testing.T) {
	assertions := assert.New(t)

	// Test minimum and maximum values for different types
	assertions.Equal(int8(-128), IsSigned(int8(-128)))
	assertions.Equal(int8(127), IsSigned(int8(127)))
	assertions.Equal(uint8(0), IsUnsigned(uint8(0)))
	assertions.Equal(uint8(255), IsUnsigned(uint8(255)))

	// Test very small and very large floats
	assertions.Equal(float32(1e-45), IsFloat(float32(1e-45)))
	assertions.Equal(1e-323, IsFloat(1e-323))

	// Test empty and special strings
	assertions.Equal("", IsString(""))
	assertions.Equal("\n\t\r", IsString("\n\t\r"))
	assertions.Equal("null", IsString("null"))
}
