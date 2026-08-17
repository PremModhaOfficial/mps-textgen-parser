# Types Package - Core Data Structures

## Overview

This package provides high-performance, generic data structures for the motadatagosdk:

1. **Slice**: Generic slice wrapper with extended functionality and native Go compatibility
2. **Map**: Generic map wrapper with native Go syntax support for any comparable keys
3. **IntMap**: High-performance generic map for integer keys (int, int64, uint, etc.)
4. **SwissMap**: High-performance generic map for any comparable keys using SwissTable algorithm

## Quick Decision Guide

```
What type of keys do you need?

INTEGER keys + High performance   → Use IntMap[K Integer, V any]
ANY comparable keys + Performance → Use SwissMap[K comparable, V any]
ANY comparable keys + Native syntax → Use Map[K comparable, V any]
GENERIC collection operations     → Use Slice
```

## Slice - Generic Slice Wrapper

### What It Does
- Provides **extended functionality** while preserving native Go slice operations
- Supports **safe access methods** that avoid panics (GetOr, GetSafe, FirstOr, LastOr)
- Offers **in-place mutation** methods (Append, Prepend, Insert, Remove, Reverse)
- **Optimized Prepend and Insert** that avoid allocation when capacity is sufficient

### When to Use
- You need slice operations with safe access patterns
- You want to avoid panics on out-of-bounds access
- You need optimized prepend/insert operations
- You want both method-based and native slice syntax

### Basic Usage

```go
import "motadatagosdk/core/types"

// Create a new slice
s := types.New(1, 2, 3, 4, 5)

// Basic operations
s.Append(6)           // [1, 2, 3, 4, 5, 6]
s.Prepend(0)          // [0, 1, 2, 3, 4, 5, 6]
s.Insert(3, 99)       // [0, 1, 2, 99, 3, 4, 5, 6]
s.Remove(3)           // [0, 1, 2, 3, 4, 5, 6]

// Safe access (no panics)
value := s.Get(2)                    // 2 (panics if out of bounds)
value = s.GetOr(100, -1)             // -1 (default value)
value, ok := s.GetSafe(100)          // 0, false

// First/Last access
first := s.First()                   // 0 (panics if empty)
first = s.FirstOr(-1)                // 0 (or -1 if empty)
first, ok = s.FirstSafe()            // 0, true

last := s.Last()                     // 6 (panics if empty)
last = s.LastOr(-1)                  // 6 (or -1 if empty)
last, ok = s.LastSafe()              // 6, true

// Native Go slice operations still work
sub := s[1:4]                        // native reslicing
for i, v := range s {                // native iteration
    fmt.Println(i, v)
}
```

### Creation Functions

```go
// From elements
s := types.New(1, 2, 3, 4, 5)

// From existing slice
existing := []int{1, 2, 3}
s := types.From(existing)

// With length and capacity
s := types.Make[int](10, 20)         // len=10, cap=20

// With length only (zero values)
s := types.WithLength[int](10)       // len=10, cap=10

// With capacity only (for pre-allocation)
s := types.WithCapacity[int](100)    // len=0, cap=100
```

### Search Operations

```go
s := types.New(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

// Find first element matching predicate
value, found := s.Find(func(v int) bool {
    return v > 5
})
// value=6, found=true

// Count elements matching predicate
count := s.Count(func(v int) bool {
    return v%2 == 0
})
// count=5 (even numbers: 2, 4, 6, 8, 10)
```

### String Operations

```go
s := types.New("apple", "banana", "cherry")

// Join with separator
result := types.Join(s, ", ")
// result = "apple, banana, cherry"
```

## Map - Generic Map with Native Syntax

### What It Does
- **Generic map** wrapping Go's built-in `map[K]V`
- Supports **any comparable key type** and **any value type**
- Provides **helper methods** while supporting **native map syntax**
- Ideal when you need both convenience methods and native operations

### When to Use
- You need a map with native Go syntax support (`m["key"] = value`)
- You want type-safe keys and values
- Simplicity over performance is preferred

### Basic Usage

```go
import "motadatagosdk/core/types"

// Create a new map with string keys and any values
m := types.NewMap[string, any]()

// Using methods
m.Set("name", "John")
m.Set("age", 30)

value, exists := m.Get("name")       // "John", true
value = m.GetValue("name")           // "John" (zero value if not found)
exists = m.Has("name")               // true

m.Delete("age")

// Native map syntax also works
m["email"] = "john@example.com"
email := m["email"]

// Iteration
m.Range(func(key string, value any) bool {
    fmt.Printf("%s: %v\n", key, value)
    return true // continue iteration
})

// Get all keys/values
keys := m.Keys()                     // []string{"name", "email"}
values := m.Values()                 // []any{"John", "john@example.com"}

// Copy and clear
copy := m.Copy()                     // shallow copy
m.Clear()                            // remove all entries
```

### Type-Safe Usage

```go
// With specific value types - no type assertions needed!
m := types.NewMap[string, int]()
m.Set("count", 42)
value, _ := m.Get("count")           // value is int, not any

// With custom key types (any comparable type)
type UserID struct { ID int }
m := types.NewMap[UserID, string]()
m.Set(UserID{ID: 1}, "Alice")

// With int keys
mInt := types.NewMap[int, string]()
mInt.Set(1, "one")
```

### Advanced Operations

```go
// Create with capacity hint
m := types.NewMapWithCapacity[string, any](100)

// Put if not exists
value, inserted := m.PutIfNotExists("key", "value")
// If key exists: returns existing value, false
// If key doesn't exist: inserts and returns new value, true
```

## IntMap - High-Performance Generic Integer Key Map

### What It Does
- **High-performance generic** map optimized for integer keys
- Supports **any integer type** as key (int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64)
- Supports **any value type** through generics
- Uses **open addressing with linear probing** for cache-friendly lookups
- Wraps `github.com/kamstrup/intmap` for optimal performance

### When to Use
- You have integer keys (IDs, indices, counters)
- You want type-safe keys and values
- Performance is critical
- Read-heavy workloads

### Basic Usage

```go
import "motadatagosdk/core/types"

// Create a new IntMap with int keys and any values
m := types.NewIntMap[int, any]()

// Basic operations
m.Set(1, "value1")
m.Set(2, "value2")
m.Set(3, "value3")

value, exists := m.Get(1)            // "value1", true
exists = m.Has(2)                    // true
deleted := m.Delete(3)               // true

// Iteration
m.Range(func(key int, value any) bool {
    fmt.Printf("%d: %v\n", key, value)
    return true // continue iteration
})

// Get all keys/values
keys := m.Keys()                     // []int{1, 2}
values := m.Values()                 // []any{"value1", "value2"}

// Copy and clear
copy := m.Copy()                     // shallow copy
m.Clear()                            // remove all entries

// Length
length := m.Len()                    // 2
```

### Type-Safe Usage

```go
// With specific value types - no type assertions needed!
m := types.NewIntMap[int, string]()
m.Set(1, "hello")
value, _ := m.Get(1)                 // value is string, not any

// With int64 keys
m64 := types.NewIntMap[int64, MyStruct]()
m64.Set(int64(1234567890), MyStruct{...})

// With uint keys
mUint := types.NewIntMap[uint, []byte]()
mUint.Set(uint(1), []byte("data"))
```

### With Capacity

```go
// Create with capacity hint for better performance
m := types.NewIntMapWithCapacity[int, string](1000)

// Put if not exists
value, inserted := m.PutIfNotExists(1, "value")
```

## SwissMap - High-Performance Generic Map

### What It Does
- **High-performance generic** map using Google's SwissTable algorithm
- Supports **any comparable key type** (string, int, struct, etc.)
- Supports **any value type** through generics
- Uses **SIMD instructions** for parallel key matching
- Excellent performance for both small and large maps
- Better memory efficiency than Go's built-in map

### When to Use
- You need a high-performance map with any comparable key type
- Performance is critical
- Large maps with frequent lookups

### Basic Usage

```go
import "motadatagosdk/core/types"

// Create a new SwissMap with string keys and any values
m := types.NewSwissMap[string, any]()

// Basic operations
m.Set("key1", "value1")
m.Set("key2", "value2")
m.Set("key3", "value3")

value, exists := m.Get("key1")       // "value1", true
exists = m.Has("key2")               // true
deleted := m.Delete("key3")          // true

// Iteration
m.Range(func(key string, value any) bool {
    fmt.Printf("%s: %v\n", key, value)
    return true // continue iteration
})

// Get all keys/values
keys := m.Keys()                     // []string{"key1", "key2"}
values := m.Values()                 // []any{"value1", "value2"}

// Copy and clear
copy := m.Copy()                     // shallow copy
m.Clear()                            // remove all entries

// Length and capacity
length := m.Len()                    // 2
capacity := m.Capacity()             // current capacity
```

### Type-Safe Usage

```go
// With specific value types - no type assertions needed!
m := types.NewSwissMap[string, int]()
m.Set("count", 42)
value, _ := m.Get("count")           // value is int, not any

// With custom key types (any comparable type)
type UserID struct { ID int }
m := types.NewSwissMap[UserID, string]()
m.Set(UserID{ID: 1}, "Alice")

// With int keys (alternative to IntMap)
mInt := types.NewSwissMap[int, string]()
mInt.Set(1, "one")
```

### With Capacity

```go
// Create with capacity hint for better performance
m := types.NewSwissMapWithCapacity[string, any](1000)

// Put if not exists
value, inserted := m.PutIfNotExists("key", "value")
```

## Comparison Table

| Feature | Slice | Map | IntMap | SwissMap |
|---------|-------|-----|--------|----------|
| **Key Type** | Index (int) | comparable (generic) | Integer (generic) | comparable (generic) |
| **Value Type** | Generic T | Generic V | Generic V | Generic V |
| **Native Syntax** | Yes (s[i], s[1:3]) | Yes (m["key"]) | No | No |
| **Performance** | Standard | Standard | High | High |
| **Use Case** | Collections | Native syntax maps | ID lookups | General lookups |
| **Memory** | Standard | Standard | Efficient | Efficient |
| **Algorithm** | - | Go map | Linear probing | SwissTable |

## Performance Comparison

### Map Operations (1000 entries)

| Operation | Map | IntMap | SwissMap |
|-----------|-----|--------|----------|
| Set | Baseline | ~2x faster | ~1.5x faster |
| Get | Baseline | ~3x faster | ~2x faster |
| Has | Baseline | ~3x faster | ~2x faster |
| Delete | Baseline | ~2x faster | ~1.5x faster |

*Note: Actual performance varies based on key distribution and hardware.*

### Slice Operations

| Operation | With Capacity | Without Capacity |
|-----------|---------------|------------------|
| Append | O(1) amortized | O(1) amortized |
| Prepend | O(n) or O(1)* | O(n) |
| Insert | O(n) or O(1)* | O(n) |
| Remove | O(n) | O(n) |

*O(1) when capacity is sufficient (no allocation needed)*

## Best Practices

### Slice Best Practices
1. **Pre-allocate capacity**: Use `WithCapacity` when you know the approximate size
2. **Use safe methods**: Prefer `GetOr`, `GetSafe` over `Get` to avoid panics
3. **Return to original**: Prepend/Insert operations modify in place
4. **Clone when sharing**: Use `Clone()` to create independent copies

### Map Best Practices
1. **Choose the right type**: Use IntMap for int keys, SwissMap for string keys
2. **Pre-allocate**: Use `WithCapacity` constructors when size is known
3. **Use Has for existence**: Check with `Has()` before `Get()` if you only need existence
4. **Handle missing keys**: Always check the boolean return from `Get()`

### General Best Practices
1. **Profile first**: Measure before optimizing
2. **Consider memory**: High-performance maps trade memory for speed
3. **Thread safety**: These types are NOT thread-safe; use sync primitives if needed

## Testing

Run tests for all types:

```bash
# Run all tests
go test ./...

# Run with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...

# Run specific benchmarks
go test -bench=BenchmarkSwissMap ./...
go test -bench=BenchmarkIntMap ./...
go test -bench=BenchmarkSlice ./...

# Run with memory allocation stats
go test -bench=. -benchmem ./...
```

## Thread Safety

**None of these types are thread-safe by default.** For concurrent access, wrap operations with appropriate synchronization:

```go
import "sync"

type SafeMap[K comparable, V any] struct {
    mu sync.RWMutex
    m  types.Map[K, V]
}

func (sm *SafeMap[K, V]) Get(key K) (V, bool) {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    return sm.m.Get(key)
}

func (sm *SafeMap[K, V]) Set(key K, value V) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    sm.m.Set(key, value)
}
```

## String Utilities

### Overview

The `string.go` file provides a comprehensive set of string helper functions for common operations including validation, encoding/decoding, hashing, encryption, and password management.

### Empty/Blank Checks

```go
import "motadatagosdk/core/types"

types.IsEmpty("")           // true
types.IsEmpty(" ")          // false
types.IsBlank("   ")        // true (whitespace only)
types.IsNotEmpty("hello")   // true
types.IsNotBlank("  a  ")   // true
```

### Character Type Checks

```go
types.IsASCII("hello")      // true
types.IsASCII("日本語")      // false
types.IsAlpha("Hello")      // true
types.IsNumeric("12345")    // true
types.IsAlphaNumeric("abc123") // true
types.IsLower("hello")      // true
types.IsUpper("HELLO")      // true
types.IsDigit("123")        // true
types.IsHex("1a2b3c")       // true
```

### Format Validation

```go
types.IsEmail("user@example.com")     // true
types.IsURL("https://example.com")    // true
```

### Encoding/Decoding

```go
// Base64
encoded := types.Base64Encode("hello")           // "aGVsbG8="
decoded, err := types.Base64Decode("aGVsbG8=")   // "hello", nil

// URL encoding
encoded := types.URLEncode("hello world")        // "hello+world"
decoded, err := types.URLDecode("hello+world")   // "hello world", nil

// Byte conversion
bytes := types.ToBytes("hello")                  // []byte{104, 101, 108, 108, 111}
```

### Trimming and Whitespace

```go
types.Trim("  hello  ", " ")        // "hello"
types.TrimLeft("xxhello", "x")      // "hello"
types.TrimRight("helloxx", "x")     // "hello"
types.TrimSpace("  hello  ")        // "hello"
types.NormalizeSpace("  a   b  ")   // "a b"
types.CollapseWhitespace("a   b")   // "a b"
types.RemoveWhitespace("a b c")     // "abc"
```

### Case Conversion

```go
types.ToLower("HELLO")              // "hello"
types.ToUpper("hello")              // "HELLO"
types.ToTitle("hello world")        // "Hello World"
```

### Hashing (using CityHash)

CityHash is a fast, non-cryptographic hash function developed by Google.

```go
hash64 := types.Hash64("hello")     // uint64
hash32 := types.Hash32("hello")     // uint32
hash128 := types.Hash128("hello")   // city.U128 (with .Low and .High fields)
checksum := types.Checksum("hello") // uint32 (CRC32)
```

### AES Encryption/Decryption

AES-GCM authenticated encryption with automatic nonce generation.

```go
// Key must be 16, 24, or 32 bytes for AES-128, AES-192, or AES-256
key := "1234567890123456" // 16 bytes for AES-128

// Encrypt
ciphertext, err := types.AESEncrypt("secret message", key)

// Decrypt
plaintext, err := types.AESDecrypt(ciphertext, key)
```

### Password Hashing (using bcrypt)

Secure password hashing with automatic salt generation.

```go
// Hash a password (uses bcrypt with default cost)
// Note: bcrypt only uses the first 72 bytes of the password
hash, err := types.HashPassword("mySecurePassword")

// Verify a password
if types.VerifyPassword(hash, "mySecurePassword") {
    // Password matches
}

// Same password produces different hashes (due to salt)
hash1, _ := types.HashPassword("password")
hash2, _ := types.HashPassword("password")
// hash1 != hash2, but both verify correctly
```

### String Functions Reference

| Category | Functions |
|----------|-----------|
| **Empty/Blank** | `IsEmpty`, `IsBlank`, `IsNotEmpty`, `IsNotBlank` |
| **Character Checks** | `IsASCII`, `IsAlpha`, `IsNumeric`, `IsAlphaNumeric`, `IsLower`, `IsUpper`, `IsDigit`, `IsHex` |
| **Format Validation** | `IsEmail`, `IsURL` |
| **Encoding** | `Base64Encode`, `Base64Decode`, `URLEncode`, `URLDecode`, `ToBytes` |
| **Trimming** | `Trim`, `TrimLeft`, `TrimRight`, `TrimSpace`, `NormalizeSpace`, `CollapseWhitespace`, `RemoveWhitespace` |
| **Case Conversion** | `ToLower`, `ToUpper`, `ToTitle` |
| **Hashing** | `Hash64`, `Hash32`, `Hash128`, `Checksum` |
| **Encryption** | `AESEncrypt`, `AESDecrypt` |
| **Password** | `HashPassword`, `VerifyPassword` |

## Summary

- **Slice**: For generic collections with safe access and optimized operations
- **Map**: For generic maps with native Go syntax support
- **IntMap**: For high-performance generic maps with integer keys
- **SwissMap**: For high-performance generic maps with any comparable keys
- **String Utilities**: For common string operations, validation, encoding, hashing, and encryption

Choose based on your key type and performance requirements. All types are designed for high-performance applications with clean, idiomatic Go APIs.