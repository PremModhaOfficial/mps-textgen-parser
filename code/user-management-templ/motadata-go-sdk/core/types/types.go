// Package types provides comprehensive type constraints and validation utilities
// for Go generics, enabling type-safe operations across the SDK.
//
// ARCHITECTURE OVERVIEW:
//
// This file defines the type constraint hierarchy used throughout the SDK for
// generic programming. The constraints form a taxonomy of numeric and primitive
// types that enable compile-time type safety while maintaining flexibility.
//
// TYPE CONSTRAINT HIERARCHY:
//
//	               Numeric
//	              /        \
//	         Integer        Float
//	        /      \         |
//	   Signed    Unsigned   float32/64
//	     |          |
//	int/8/16/32/64 uint/8/16/32/64/ptr
//
// DESIGN PRINCIPLES:
//
// 1. Type Safety First:
//   - Compile-time constraint validation
//   - No runtime type assertions needed
//   - Clear constraint boundaries
//
// 2. Tilde Operator Usage (~):
//   - Allows underlying type matching
//   - Supports type aliases and custom types
//   - More flexible than exact type matching
//
// 3. Composable Constraints:
//   - Build complex from simple constraints
//   - Reusable across packages
//   - Clear semantic meaning
//
// ARCHITECTURAL DECISIONS:
//
// Why Separate Signed/Unsigned:
// - Different overflow behaviors
// - Different valid operations (e.g., bit shifting)
// - Clearer API intentions
//
// Why Include uintptr in Unsigned:
// - Pointer arithmetic operations
// - Memory address calculations
// - Unsafe pointer conversions
//
// Why Validation Functions:
// - Compile-time type checking
// - Zero runtime overhead
// - Documentation through types
//
// USAGE PATTERNS:
//
//  1. Generic Function Constraints:
//     func Sum[T Numeric](values []T) T
//
//  2. Type-Safe Collections:
//     type NumericMap[K comparable, V Numeric] map[K]V
//
//  3. Algorithm Implementations:
//     func BinarySearch[T Integer](arr []T, target T) int
//
// PERFORMANCE CONSIDERATIONS:
//
// - Zero runtime overhead (resolved at compile time)
// - No interface boxing or type assertions
// - Inlined by compiler in most cases
// - No reflection required
//
// FUTURE EXTENSIBILITY:
//
// Additional constraints can be added for:
// - Complex numbers
// - Fixed-point arithmetic
// - Vector/Matrix types
// - Custom numeric types
package types

// ============================================================================
// SECTION: Numeric Type Constraints
//
// These constraints define the numeric type hierarchy used throughout the SDK.
// They enable generic algorithms to work with appropriate numeric types while
// maintaining type safety and preventing invalid operations.
//
// Architecture Notes:
// - Tilde (~) prefix allows derived types and aliases
// - Union types (|) provide flexibility
// - Composition creates higher-level constraints
// ============================================================================

// Signed represents all signed integer types in Go.
// These types can represent negative values and use two's complement arithmetic.
// Range examples:
// - int8: -128 to 127
// - int16: -32,768 to 32,767
// - int32: -2,147,483,648 to 2,147,483,647
// - int64: -9,223,372,036,854,775,808 to 9,223,372,036,854,775,807
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned represents all unsigned integer types in Go.
// These types can only represent non-negative values and support bitwise operations.
// The inclusion of uintptr enables pointer arithmetic in unsafe operations.
// Range examples:
// - uint8: 0 to 255
// - uint16: 0 to 65,535
// - uint32: 0 to 4,294,967,295
// - uint64: 0 to 18,446,744,073,709,551,615
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer combines all signed and unsigned integer types.
// This constraint is used when an algorithm works with any integer type
// regardless of signedness, such as counting, indexing, or basic arithmetic.
type Integer interface {
	Signed | Unsigned
}

// Float represents all floating-point types in Go.
// These types follow IEEE 754 standard for floating-point arithmetic.
// - float32: 32-bit IEEE 754 (7 decimal digits precision)
// - float64: 64-bit IEEE 754 (15 decimal digits precision)
type Float interface {
	~float32 | ~float64
}

// Numeric combines all numeric types (integers and floats).
// This is the most general numeric constraint, used for algorithms that
// work with any numeric type, such as min/max, sum, or statistical operations.
type Numeric interface {
	Integer | Float
}

// Char represents string types for character operations.
// This constraint is primarily for string manipulation functions
// and text processing algorithms.
// Note: Unlike C/C++, Go doesn't have a char type; rune and byte are integers.
type Char interface {
	string
}

// ============================================================================
// SECTION: Type Validation Functions
//
// These functions provide compile-time type validation through Go's generics.
// They serve as both documentation and compile-time checks, ensuring that
// values passed to generic functions satisfy the required constraints.
//
// Architecture Purpose:
// - Compile-time type checking without runtime cost
// - Self-documenting code through explicit constraints
// - Building blocks for more complex generic operations
// ============================================================================

// IsInteger validates that a value satisfies the Integer constraint at compile time.
// Returns the value unchanged if valid, compilation error otherwise.
// This is a zero-cost abstraction - the function is typically inlined.
func IsInteger[T Integer](v T) T {
	return v // Only compiles if T satisfies Integer constraint
}

// IsFloat validates that a value satisfies the Float constraint at compile time.
// Used to ensure floating-point operations are only performed on appropriate types.
func IsFloat[T Float](v T) T {
	return v // Only compiles if T satisfies Float constraint
}

// IsString validates that a value is of string type at compile time.
// Useful for generic text processing functions that require string inputs.
func IsString[T string](v T) T {
	return v // Only compiles if T satisfies String constraint
}

// IsBoolean validates that a value is of bool type at compile time.
// Used in generic conditional logic and boolean algebra operations.
func IsBoolean[T bool](v T) T {
	return v // Only compiles if T satisfies Boolean constraint
}

// IsSigned validates that a value satisfies the Signed constraint at compile time.
// Important for operations that require signed arithmetic, such as
// absolute value calculations or sign-dependent algorithms.
func IsSigned[T Signed](v T) T {
	return v
}

func IsUnsigned[T Unsigned](v T) T {

	return v
}
