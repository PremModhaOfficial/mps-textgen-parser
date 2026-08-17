package utils

import "sync"

// ============================================================================
// Type Constraints (Generic Type Definitions)
// ============================================================================

// Signed is a constraint for signed integer types.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is a constraint for unsigned integer types.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer is a constraint for all integer types.
type Integer interface {
	Signed | Unsigned
}

// Float is a constraint for floating-point types.
type Float interface {
	~float32 | ~float64
}

// Numeric is a constraint for all numeric types.
type Numeric interface {
	Integer | Float
}

// Ordered is a constraint for types that support ordering (<, >, <=, >=).
type Ordered interface {
	Integer | Float | ~string
}

// ============================================================================
// Optional Value Type
// ============================================================================

// Optional represents a value that may or may not be present.
// This is useful for distinguishing between "value not set" and "value set to zero".
type Optional[T any] struct {
	value T
	valid bool
}

// Some creates an Optional containing a value.
func Some[T any](value T) Optional[T] {
	return Optional[T]{value: value, valid: true}
}

// None creates an Optional with no value.
func None[T any]() Optional[T] {
	return Optional[T]{}
}

// Get returns the value and whether it's present.
func (o Optional[T]) Get() (T, bool) {
	return o.value, o.valid
}

// Value returns the value or panics if not present.
func (o Optional[T]) Value() T {
	if !o.valid {
		panic("optional value not present")
	}
	return o.value
}

// ValueOr returns the value or the provided default if not present.
func (o Optional[T]) ValueOr(defaultValue T) T {
	if o.valid {
		return o.value
	}
	return defaultValue
}

// IsPresent returns true if a value is present.
func (o Optional[T]) IsPresent() bool {
	return o.valid
}

// IsEmpty returns true if no value is present.
func (o Optional[T]) IsEmpty() bool {
	return !o.valid
}

// ============================================================================
// Result Type
// ============================================================================

// Result represents the outcome of an operation that may fail.
// It either contains a value of type T or an error.
type Result[T any] struct {
	value T
	err   error
}

// Ok creates a successful Result with a value.
func Ok[T any](value T) Result[T] {
	return Result[T]{value: value}
}

// Err creates a failed Result with an error.
func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

// Get returns the value and error.
func (r Result[T]) Get() (T, error) {
	return r.value, r.err
}

// Value returns the value or panics if there's an error.
func (r Result[T]) Value() T {
	if r.err != nil {
		panic(r.err)
	}
	return r.value
}

// ValueOr returns the value or the provided default if there's an error.
func (r Result[T]) ValueOr(defaultValue T) T {
	if r.err != nil {
		return defaultValue
	}
	return r.value
}

// Error returns the error (nil if successful).
func (r Result[T]) Error() error {
	return r.err
}

// IsOk returns true if the result is successful.
func (r Result[T]) IsOk() bool {
	return r.err == nil
}

// IsErr returns true if the result contains an error.
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// ============================================================================
// Pair Type
// ============================================================================

// Pair represents a key-value pair.
type Pair[K, V any] struct {
	Key   K
	Value V
}

// NewPair creates a new Pair.
func NewPair[K, V any](key K, value V) Pair[K, V] {
	return Pair[K, V]{Key: key, Value: value}
}

// ============================================================================
// SafeMap Type (Thread-Safe Map)
// ============================================================================

// SafeMap is a thread-safe map implementation using sync.RWMutex.
type SafeMap[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]V
}

// NewSafeMap creates a new SafeMap.
func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		items: make(map[K]V),
	}
}

// NewSafeMapWithCapacity creates a new SafeMap with the specified capacity.
func NewSafeMapWithCapacity[K comparable, V any](capacity int) *SafeMap[K, V] {
	return &SafeMap[K, V]{
		items: make(map[K]V, capacity),
	}
}

// Set stores a key-value pair in the map.
func (m *SafeMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	m.items[key] = value
	m.mu.Unlock()
}

// Get retrieves a value by key.
func (m *SafeMap[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	value, exists := m.items[key]
	m.mu.RUnlock()
	return value, exists
}

// GetOrSet returns the existing value for a key, or sets and returns a new value.
func (m *SafeMap[K, V]) GetOrSet(key K, value V) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if existing, exists := m.items[key]; exists {
		return existing, true
	}

	m.items[key] = value
	return value, false
}

// Delete removes a key-value pair from the map.
func (m *SafeMap[K, V]) Delete(key K) {
	m.mu.Lock()
	delete(m.items, key)
	m.mu.Unlock()
}

// Has checks if a key exists in the map.
func (m *SafeMap[K, V]) Has(key K) bool {
	m.mu.RLock()
	_, exists := m.items[key]
	m.mu.RUnlock()
	return exists
}

// Len returns the number of entries in the map.
func (m *SafeMap[K, V]) Len() int {
	m.mu.RLock()
	length := len(m.items)
	m.mu.RUnlock()
	return length
}

// Clear removes all entries from the map.
func (m *SafeMap[K, V]) Clear() {
	m.mu.Lock()
	m.items = make(map[K]V)
	m.mu.Unlock()
}

// Keys returns all keys in the map.
func (m *SafeMap[K, V]) Keys() []K {
	m.mu.RLock()
	keys := make([]K, 0, len(m.items))
	for k := range m.items {
		keys = append(keys, k)
	}
	m.mu.RUnlock()
	return keys
}

// Values returns all values in the map.
func (m *SafeMap[K, V]) Values() []V {
	m.mu.RLock()
	values := make([]V, 0, len(m.items))
	for _, v := range m.items {
		values = append(values, v)
	}
	m.mu.RUnlock()
	return values
}

// Range iterates over all key-value pairs.
// The callback function should return true to continue iteration.
func (m *SafeMap[K, V]) Range(fn func(key K, value V) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for k, v := range m.items {
		if !fn(k, v) {
			break
		}
	}
}

// Copy creates a shallow copy of the map.
func (m *SafeMap[K, V]) Copy() *SafeMap[K, V] {
	m.mu.RLock()
	defer m.mu.RUnlock()

	newMap := NewSafeMapWithCapacity[K, V](len(m.items))
	for k, v := range m.items {
		newMap.items[k] = v
	}
	return newMap
}

// ============================================================================
// SafeSlice Type (Thread-Safe Slice)
// ============================================================================

// SafeSlice is a thread-safe slice implementation.
type SafeSlice[T any] struct {
	mu    sync.RWMutex
	items []T
}

// NewSafeSlice creates a new SafeSlice.
func NewSafeSlice[T any]() *SafeSlice[T] {
	return &SafeSlice[T]{
		items: make([]T, 0),
	}
}

// NewSafeSliceWithCapacity creates a new SafeSlice with the specified capacity.
func NewSafeSliceWithCapacity[T any](capacity int) *SafeSlice[T] {
	return &SafeSlice[T]{
		items: make([]T, 0, capacity),
	}
}

// Append adds elements to the end of the slice.
func (s *SafeSlice[T]) Append(elements ...T) {
	s.mu.Lock()
	s.items = append(s.items, elements...)
	s.mu.Unlock()
}

// Get returns the element at the given index.
func (s *SafeSlice[T]) Get(index int) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if index < 0 || index >= len(s.items) {
		var zero T
		return zero, false
	}
	return s.items[index], true
}

// Set sets the element at the given index.
func (s *SafeSlice[T]) Set(index int, value T) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if index < 0 || index >= len(s.items) {
		return false
	}
	s.items[index] = value
	return true
}

// Len returns the length of the slice.
func (s *SafeSlice[T]) Len() int {
	s.mu.RLock()
	length := len(s.items)
	s.mu.RUnlock()
	return length
}

// Clear removes all elements from the slice.
func (s *SafeSlice[T]) Clear() {
	s.mu.Lock()
	s.items = s.items[:0]
	s.mu.Unlock()
}

// All returns a copy of all elements.
func (s *SafeSlice[T]) All() []T {
	s.mu.RLock()
	result := make([]T, len(s.items))
	copy(result, s.items)
	s.mu.RUnlock()
	return result
}

// Range iterates over all elements.
// The callback function should return true to continue iteration.
func (s *SafeSlice[T]) Range(fn func(index int, value T) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i, v := range s.items {
		if !fn(i, v) {
			break
		}
	}
}

// Pop removes and returns the last element.
func (s *SafeSlice[T]) Pop() (T, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.items) == 0 {
		var zero T
		return zero, false
	}

	index := len(s.items) - 1
	value := s.items[index]
	s.items = s.items[:index]
	return value, true
}
