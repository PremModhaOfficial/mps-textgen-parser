package types

import "github.com/dolthub/swiss"

// Package types provides core data structures for the motadata-go-sdk.
//
// SwissMap is a high-performance map implementation based on Google's SwissTable design.
// It wraps github.com/dolthub/swiss which uses SIMD instructions for parallel key matching,
// providing excellent performance for both small and large maps. SwissMap is particularly
// efficient for string keys and offers better memory efficiency than Go's built-in map.
//
// Usage:
//
//	m := types.NewSwissMap()
//	m.Set("key", "value")
//	value, exists := m.Get("key")
//	m.Delete("key")

// SwissMap is a high-performance map with string keys and any values.
// It uses SwissTable algorithm for fast lookups with SIMD optimizations.
type SwissMap swiss.Map[string, any]

// internal returns the underlying swiss.Map for method access.
// This method is inlined by the compiler, so there is no performance overhead.
func (m *SwissMap) internal() *swiss.Map[string, any] {
	return (*swiss.Map[string, any])(m)
}

// NewSwissMap creates and returns a new SwissMap with default capacity (16).
func NewSwissMap() *SwissMap {
	return (*SwissMap)(swiss.NewMap[string, any](16))
}

// NewSwissMapWithCapacity creates and returns a new SwissMap with the specified initial capacity.
// Use this when you know the approximate size to avoid rehashing.
func NewSwissMapWithCapacity(capacity uint32) *SwissMap {
	return (*SwissMap)(swiss.NewMap[string, any](capacity))
}

// Set stores a key-value pair in the map.
// If the key already exists, its value is updated.
func (m *SwissMap) Set(key string, value any) {
	m.internal().Put(key, value)
}

// Get retrieves a value by key.
// Returns the value and true if found, or nil and false if not found.
func (m *SwissMap) Get(key string) (any, bool) {
	return m.internal().Get(key)
}

// Delete removes a key-value pair from the map.
// Returns true if the key was found and deleted, false otherwise.
func (m *SwissMap) Delete(key string) bool {
	return m.internal().Delete(key)
}

// Has checks if a key exists in the map.
// Returns true if the key exists, false otherwise.
func (m *SwissMap) Has(key string) bool {
	return m.internal().Has(key)
}

// Len returns the number of entries in the map.
func (m *SwissMap) Len() int {
	return m.internal().Count()
}

// Keys returns a slice containing all keys in the map.
// The order of keys is not guaranteed.
func (m *SwissMap) Keys() []string {
	keys := make([]string, m.Len())
	index := 0
	m.internal().Iter(func(key string, _ any) bool {
		keys[index] = key
		index++
		return false
	})
	return keys
}

// Values returns a slice containing all values in the map.
// The order of values is not guaranteed.
func (m *SwissMap) Values() []any {
	values := make([]any, m.Len())
	index := 0
	m.internal().Iter(func(_ string, value any) bool {
		values[index] = value
		index++
		return false
	})
	return values
}

// Clear removes all entries from the map.
func (m *SwissMap) Clear() {
	m.internal().Clear()
}

// Range iterates over all key-value pairs in the map.
// The callback function is called for each entry. If it returns false, iteration stops.
func (m *SwissMap) Range(fn func(key string, value any) bool) {
	m.internal().Iter(func(key string, value any) bool {
		return !fn(key, value)
	})
}

// Copy creates and returns a shallow copy of the map.
// The new map is independent of the original.
func (m *SwissMap) Copy() *SwissMap {
	newMap := NewSwissMapWithCapacity(uint32(m.Len()))
	m.Range(func(key string, value any) bool {
		newMap.Set(key, value)
		return true
	})
	return newMap
}

// PutIfNotExists adds a key-value pair only if the key doesn't already exist.
// Returns the value (existing or new) and true if the key was newly inserted.
func (m *SwissMap) PutIfNotExists(key string, value any) (any, bool) {
	if existing, exists := m.Get(key); exists {
		return existing, false
	}
	m.Set(key, value)
	return value, true
}

// Capacity returns the current capacity of the map.
// This is the number of entries that can be stored before a resize is needed.
func (m *SwissMap) Capacity() int {
	return m.internal().Capacity()
}
