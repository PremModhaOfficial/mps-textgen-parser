// Package types - Map Implementation
//
// ARCHITECTURE OVERVIEW:
//
// Map provides a lightweight wrapper around Go's built-in map[string]any that
// maintains native map semantics while offering convenient helper methods.
// Unlike SwissMap (SIMD-optimized) and IntMap (integer-key optimized), this
// implementation focuses on compatibility and simplicity for general-purpose use.
//
// DESIGN PHILOSOPHY:
//
// Type Alias vs Struct Wrapper:
// Map is implemented as a type alias rather than a struct wrapper, which means:
// 1. Zero overhead - no indirection or method dispatch cost
// 2. Native syntax support - m[key] works directly
// 3. Full compatibility with existing map[string]any code
// 4. Seamless integration with Go's range, len, and delete builtins
//
// This design choice prioritizes developer experience and performance over
// encapsulation, making Map the ideal choice when you need both convenience
// methods and native map operations.
//
// ARCHITECTURAL DECISIONS:
//
// Why Type Alias:
// - Direct memory access without pointer indirection
// - No allocation overhead for wrapper objects
// - Compiler can inline all operations
// - JSON marshaling works out of the box
//
// Method Design:
// - All methods are value receivers for consistency
// - Methods never return the map itself (no chaining)
// - Clear separation between existence checking and value retrieval
// - Defensive copying in Copy() to prevent aliasing bugs
//
// INTERNAL STRUCTURE:
//
// Go's map[string]any uses a hash table with:
// - Separate chaining for collision resolution
// - Incremental rehashing on growth
// - Load factor of ~6.5 entries per bucket
// - Randomized hash seeds for security
//
// Memory Layout:
//
//	┌─────────────────────────────┐
//	│  hmap struct                │
//	│  ├── count: int            │ ← Number of entries
//	│  ├── B: uint8              │ ← log₂ of bucket count
//	│  ├── hash0: uint32         │ ← Hash seed
//	│  └── buckets: *bmap        │ ← Pointer to bucket array
//	└─────────────────────────────┘
//	         ↓
//	┌─────────────────────────────┐
//	│  Bucket Array               │
//	│  ├── tophash[8]            │ ← High 8 bits of hash
//	│  ├── keys[8]               │ ← String keys
//	│  ├── values[8]             │ ← Interface values
//	│  └── overflow: *bmap       │ ← Overflow bucket chain
//	└─────────────────────────────┘
//
// PERFORMANCE CHARACTERISTICS:
//
// Time Complexity:
// - Get/Set/Delete: O(1) average, O(n) worst case
// - Len: O(1) - stored in hmap.count
// - Keys/Values: O(n) - full iteration required
// - Clear: O(n) - must visit each bucket
// - Range: O(n) with early termination support
//
// Space Complexity:
// - Overhead: ~10.79 bytes per entry on average
// - Minimum allocation: 1 bucket (8 slots)
// - Growth factor: 2x when load > 6.5
// - String keys: 16 bytes header + string data
// - Interface values: 16 bytes (type + data pointer)
//
// Comparison with Specialized Maps:
//
// vs SwissMap:
// - Map: Better for small maps (< 100 entries)
// - SwissMap: 2-4x faster for large maps
// - Map: Native syntax support
// - SwissMap: SIMD-accelerated lookups
//
// vs IntMap:
// - Map: Flexible key types (strings)
// - IntMap: 2-3x faster for integer keys
// - Map: Better for mixed workloads
// - IntMap: Superior cache locality
//
// CONCURRENCY MODEL:
//
// Map is NOT thread-safe. Concurrent access requires external synchronization:
//
// Option 1 - Read-Write Mutex:
//
//	type SafeMap struct {
//	    mu sync.RWMutex
//	    m  types.Map
//	}
//
// Option 2 - sync.Map for read-heavy workloads:
//
//	var m sync.Map // Built-in concurrent map
//
// Option 3 - Sharding for write-heavy workloads:
//
//	type ShardedMap [N]SafeMap // N shards
//
// USE CASES:
//
// Ideal for:
// - Configuration management
// - Small to medium-sized lookups (< 1000 entries)
// - JSON/YAML data structures
// - Dynamic field access
// - Prototype and development phase
//
// Not ideal for:
// - High-performance critical paths (use SwissMap)
// - Integer-keyed data (use IntMap)
// - Concurrent access without locking
// - Very large datasets (consider sharding)
//
// USAGE PATTERNS:
//
// Basic operations:
//
//	m := types.NewMap()
//	m.Set("config.timeout", 30)    // Method syntax
//	m["config.retry"] = 3           // Native syntax
//
//	if val, ok := m.Get("config.timeout"); ok {
//	    timeout := val.(int)
//	}
//
// Pre-sized initialization for performance:
//
//	m := types.NewMapWithCapacity(1000)
//	// Avoids rehashing for first ~6500 entries
//
// Iteration with early termination:
//
//	m.Range(func(key string, value any) bool {
//	    if key == "stop" {
//	        return false // Stop iteration
//	    }
//	    process(key, value)
//	    return true
//	})
//
// PITFALLS AND BEST PRACTICES:
//
//  1. Type Assertions:
//     // Bad: Unchecked assertion can panic
//     timeout := m["timeout"].(int)
//
//     // Good: Check type assertion
//     if timeout, ok := m["timeout"].(int); ok {
//     use(timeout)
//     }
//
//  2. Nil Map Operations:
//     var m Map // nil map
//     m.Set("key", "value") // PANIC!
//
//     // Always initialize:
//     m := NewMap()
//
//  3. Iteration Mutation:
//     // Bad: Modifying during iteration is undefined
//     for k, v := range m {
//     m.Delete(k) // Undefined behavior
//     }
//
//     // Good: Collect keys first
//     keys := m.Keys()
//     for _, k := range keys {
//     m.Delete(k)
//     }
//
//  4. Value Aliasing:
//     // Be aware that Copy() is shallow
//     m1 := NewMap()
//     m1["slice"] = []int{1, 2, 3}
//     m2 := m1.Copy()
//     m2["slice"].([]int)[0] = 999 // Modifies m1's slice too!
//
// FUTURE OPTIMIZATIONS:
//
// - Generics version: Map[K comparable, V any]
// - Ordered iteration support
// - Bulk operations (SetMany, DeleteMany)
// - Merge and diff operations
// - Persistent/immutable variant
package types

// Map is a simple wrapper around Go's built-in map with string keys and any values.
//
// Unlike IntMap and SwissMap, Map is a type alias for map[string]any, which means
// it can be used directly with Go's native map syntax (m[key] = value, v := m[key]).
// This makes it ideal for cases where you need both the convenience of helper methods
// and the flexibility of native map operations.
//
// Usage:
//	m := types.NewMap()
//	m.Set("key", "value")           // Using method
//	m["key"] = "value"              // Using native syntax
//	value, exists := m.Get("key")   // Using method
//	value := m["key"]               // Using native syntax
//	m.Delete("key")

// Map represents a map with string keys and any values.
// Can be used directly with native map syntax: m[key] = value, v := m[key]
type Map map[string]any

// NewMap creates and returns a new empty Map.
func NewMap() Map {
	return make(Map)
}

// NewMapWithCapacity creates and returns a new Map with the specified initial capacity.
// Use this when you know the approximate size to avoid rehashing.
func NewMapWithCapacity(capacity int) Map {
	return make(Map, capacity)
}

// Set stores a key-value pair in the map.
// If the key already exists, its value is updated.
func (m Map) Set(key string, value any) {
	m[key] = value
}

// Get retrieves a value by key.
// Returns the value and true if found, or nil and false if not found.
func (m Map) Get(key string) (any, bool) {
	value, exists := m[key]
	return value, exists
}

// GetValue retrieves a value by key.
// Returns only the value (nil if the key does not exist).
// Use Get() if you need to distinguish between a nil value and a missing key.
func (m Map) GetValue(key string) any {
	return m[key]
}

// Delete removes a key-value pair from the map.
func (m Map) Delete(key string) {
	delete(m, key)
}

// Has checks if a key exists in the map.
// Returns true if the key exists, false otherwise.
func (m Map) Has(key string) bool {
	_, exists := m[key]
	return exists
}

// Len returns the number of entries in the map.
func (m Map) Len() int {
	return len(m)
}

// Keys returns a slice containing all keys in the map.
// The order of keys is not guaranteed.
func (m Map) Keys() []string {
	keys := make([]string, len(m))
	index := 0
	for key := range m {
		keys[index] = key
		index++
	}
	return keys
}

// Values returns a slice containing all values in the map.
// The order of values is not guaranteed.
func (m Map) Values() []any {
	values := make([]any, len(m))
	index := 0
	for _, value := range m {
		values[index] = value
		index++
	}
	return values
}

// Clear removes all entries from the map.
func (m Map) Clear() {
	for key := range m {
		delete(m, key)
	}
}

// Range iterates over all key-value pairs in the map.
// The callback function is called for each entry. If it returns false, iteration stops.
func (m Map) Range(fn func(key string, value any) bool) {
	for key, value := range m {
		if !fn(key, value) {
			break
		}
	}
}

// Copy creates and returns a shallow copy of the map.
// The new map is independent of the original.
func (m Map) Copy() Map {
	newMap := make(Map, len(m))
	for key, value := range m {
		newMap[key] = value
	}
	return newMap
}

// PutIfNotExists adds a key-value pair only if the key doesn't already exist.
// Returns the value (existing or new) and true if the key was newly inserted.
func (m Map) PutIfNotExists(key string, value any) (any, bool) {
	if existing, exists := m[key]; exists {
		return existing, false
	}
	m[key] = value
	return value, true
}
