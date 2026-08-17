package types

import "github.com/kamstrup/intmap"

// Package types provides core data structures for the motadata-go-sdk.
//
// IntMap is a high-performance map implementation optimized for integer keys.
// It wraps github.com/kamstrup/intmap which uses open addressing with linear
// probing for cache-friendly lookups. This makes it significantly faster than
// Go's built-in map for integer keys, especially for read-heavy workloads.
//
// Usage:
//
//	m := types.NewIntMap()
//	m.Set(1, "value")
//	value, exists := m.Get(1)
//	m.Delete(1)

// IntMap is a high-performance map with integer keys and any values.
// It provides O(1) average-case lookups with better cache locality than standard maps.
type IntMap intmap.Map[int, any]

// internal returns the underlying intmap.Map for method access.
// This method is inlined by the compiler, so there is no performance overhead.
func (m *IntMap) internal() *intmap.Map[int, any] {
	return (*intmap.Map[int, any])(m)
}

// NewIntMap creates and returns a new IntMap with default capacity (16).
func NewIntMap() *IntMap {
	return (*IntMap)(intmap.New[int, any](16))
}

// NewIntMapWithCapacity creates and returns a new IntMap with the specified initial capacity.
// Use this when you know the approximate size to avoid rehashing.
func NewIntMapWithCapacity(capacity int) *IntMap {
	return (*IntMap)(intmap.New[int, any](capacity))
}

// Set stores a key-value pair in the map.
// If the key already exists, its value is updated.
func (m *IntMap) Set(key int, value any) {
	m.internal().Put(key, value)
}

// Get retrieves a value by key.
// Returns the value and true if found, or nil and false if not found.
func (m *IntMap) Get(key int) (any, bool) {
	return m.internal().Get(key)
}

// Delete removes a key-value pair from the map.
// Returns true if the key was found and deleted, false otherwise.
func (m *IntMap) Delete(key int) bool {
	return m.internal().Del(key)
}

// Has checks if a key exists in the map.
// Returns true if the key exists, false otherwise.
func (m *IntMap) Has(key int) bool {
	return m.internal().Has(key)
}

// Len returns the number of entries in the map.
func (m *IntMap) Len() int {
	return m.internal().Len()
}

// Keys returns a slice containing all keys in the map.
// The order of keys is not guaranteed.
func (m *IntMap) Keys() []int {
	keys := make([]int, m.Len())
	index := 0
	for key := range m.internal().Keys() {
		keys[index] = key
		index++
	}
	return keys
}

// Values returns a slice containing all values in the map.
// The order of values is not guaranteed.
func (m *IntMap) Values() []any {
	values := make([]any, m.Len())
	index := 0
	for value := range m.internal().Values() {
		values[index] = value
		index++
	}
	return values
}

// Clear removes all entries from the map.
func (m *IntMap) Clear() {
	m.internal().Clear()
}

// Range iterates over all key-value pairs in the map.
// The callback function is called for each entry. If it returns false, iteration stops.
func (m *IntMap) Range(fn func(key int, value any) bool) {
	m.internal().ForEach(fn)
}

// Copy creates and returns a shallow copy of the map.
// The new map is independent of the original.
func (m *IntMap) Copy() *IntMap {
	newMap := NewIntMapWithCapacity(m.Len())
	m.Range(func(key int, value any) bool {
		newMap.Set(key, value)
		return true
	})
	return newMap
}

// PutIfNotExists adds a key-value pair only if the key doesn't already exist.
// Returns the value (existing or new) and true if the key was newly inserted.
func (m *IntMap) PutIfNotExists(key int, value any) (any, bool) {
	return m.internal().PutIfNotExists(key, value)
}
