package types

// Map is a simple wrapper around Go's built-in map with string keys and any values.
//
// Unlike IntMap and SwissMap, Map is a type alias for map[string]any, which means
// it can be used directly with Go's native map syntax (m[key] = value, v := m[key]).
// This makes it ideal for cases where you need both the convenience of helper methods
// and the flexibility of native map operations.
//
// Usage:
//
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
