package types

// Package types provides core data structures for the motadata-go-sdk.
//
// Slice is a generic slice wrapper that provides extended functionality while
// preserving the ability to use native Go slice operations like reslicing (s[1:3]),
// range iteration, and direct indexing. Unlike wrapper structs, Slice is a type
// definition for []T, making it fully compatible with Go's slice syntax.
//
// Key features:
//   - Full compatibility with native slice operations (s[i], s[1:3], range s)
//   - Safe access methods (GetOr, GetSafe, FirstOr, LastOr) that avoid panics
//   - In-place mutation methods (Append, Prepend, Insert, Remove, Reverse)
//   - Optimized Prepend and Insert that avoid allocation when capacity is sufficient
//   - Search operations (Find, Count) with predicate functions
//
// Usage:
//
//	s := types.New(1, 2, 3, 4, 5)
//	s.Append(6)                    // [1, 2, 3, 4, 5, 6]
//	s.Prepend(0)                   // [0, 1, 2, 3, 4, 5, 6]
//	value := s.Get(2)              // 2
//	value, ok := s.GetSafe(10)     // 0, false
//	sub := s[1:4]                  // native reslicing works

// Slice is a generic slice type that provides extended functionality
// while preserving the ability to use native slice operations like s[1:3].
type Slice[T any] []T

// New creates a new Slice from the given elements
func New[T any](elements ...T) Slice[T] {
	return elements
}

// From creates a Slice from an existing slice
func From[T any](s []T) Slice[T] {
	return s
}

// Make creates a new Slice with the specified length and capacity
func Make[T any](length, capacity int) Slice[T] {
	return make(Slice[T], length, capacity)
}

// WithLength creates a new Slice with the specified length
func WithLength[T any](length int) Slice[T] {
	return make(Slice[T], length)
}

// WithCapacity creates a new Slice with zero length and the specified capacity
func WithCapacity[T any](capacity int) Slice[T] {
	return make(Slice[T], 0, capacity)
}

// ============================================================================
// Basic Operations
// ============================================================================

// Len returns the length of the slice
func (s *Slice[T]) Len() int {
	return len(*s)
}

// Cap returns the capacity of the slice
func (s *Slice[T]) Cap() int {
	return cap(*s)
}

// IsEmpty returns true if the slice has no elements
func (s *Slice[T]) IsEmpty() bool {
	return len(*s) == 0
}

// IsNotEmpty returns true if the slice has at least one element
func (s *Slice[T]) IsNotEmpty() bool {
	return len(*s) > 0
}

// Get returns the element at the given index
// Panics if index is out of bounds
func (s *Slice[T]) Get(index int) T {
	return (*s)[index]
}

// GetOr returns the element at the given index or defaultValue if out of bounds
func (s *Slice[T]) GetOr(index int, defaultValue T) T {
	if index < 0 || index >= len(*s) {
		return defaultValue
	}
	return (*s)[index]
}

// GetSafe returns the element at the given index and a boolean indicating success
func (s *Slice[T]) GetSafe(index int) (T, bool) {
	if index < 0 || index >= len(*s) {
		var zero T
		return zero, false
	}
	return (*s)[index], true
}

// First returns the first element
// Panics if the slice is empty
func (s *Slice[T]) First() T {
	return (*s)[0]
}

// FirstOr returns the first element or defaultValue if empty
func (s *Slice[T]) FirstOr(defaultValue T) T {
	if len(*s) == 0 {
		return defaultValue
	}
	return (*s)[0]
}

// FirstSafe returns the first element and a boolean indicating success
func (s *Slice[T]) FirstSafe() (T, bool) {
	if len(*s) == 0 {
		var zero T
		return zero, false
	}
	return (*s)[0], true
}

// Last returns the last element
// Panics if the slice is empty
func (s *Slice[T]) Last() T {
	return (*s)[len(*s)-1]
}

// LastOr returns the last element or defaultValue if empty
func (s *Slice[T]) LastOr(defaultValue T) T {
	if len(*s) == 0 {
		return defaultValue
	}
	return (*s)[len(*s)-1]
}

// LastSafe returns the last element and a boolean indicating success
func (s *Slice[T]) LastSafe() (T, bool) {
	if len(*s) == 0 {
		var zero T
		return zero, false
	}
	return (*s)[len(*s)-1], true
}

// ============================================================================
// Modification Operations (mutate in place)
// ============================================================================

// Append adds elements to the end of the slice (mutates in place)
func (s *Slice[T]) Append(elements ...T) {
	*s = append(*s, elements...)
}

// Prepend adds elements to the beginning of the slice (mutates in place).
// Optimized to avoid allocation when current capacity is sufficient.
func (s *Slice[T]) Prepend(elements ...T) {
	if len(elements) == 0 {
		return
	}

	oldLen := len(*s)
	newLen := oldLen + len(elements)
	numNew := len(elements)

	if cap(*s) >= newLen {
		// Enough capacity: shift existing elements right, then copy new elements
		*s = (*s)[:newLen]

		// Shift existing elements to the right (copy backwards to handle overlap)
		for i := oldLen - 1; i >= 0; i-- {
			(*s)[i+numNew] = (*s)[i]
		}

		// Copy new elements at the beginning
		copy(*s, elements)
	} else {
		// Not enough capacity: allocate new slice with growth
		newCap := newLen
		if newCap < 2*cap(*s) {
			newCap = 2 * cap(*s)
		}

		newSlice := make([]T, newLen, newCap)
		copy(newSlice, elements)
		copy(newSlice[numNew:], *s)
		*s = newSlice
	}
}

// Concat combines multiple slices into this slice (mutates in place)
func (s *Slice[T]) Concat(others ...Slice[T]) {
	total := len(*s)
	for _, other := range others {
		total += len(other)
	}
	result := make(Slice[T], 0, total)
	result = append(result, *s...)
	for _, other := range others {
		result = append(result, other...)
	}
	*s = result
}

// Insert inserts elements at the given index (mutates in place).
// Optimized to avoid allocation when current capacity is sufficient.
func (s *Slice[T]) Insert(index int, elements ...T) {
	if len(elements) == 0 {
		return
	}

	if index < 0 {
		index = 0
	}
	if index > len(*s) {
		index = len(*s)
	}

	oldLen := len(*s)
	newLen := oldLen + len(elements)
	numNew := len(elements)

	if cap(*s) >= newLen {
		// Enough capacity: shift elements after index to the right, then copy new elements
		*s = (*s)[:newLen]

		// Shift elements from index to the right (copy backwards to handle overlap)
		for i := oldLen - 1; i >= index; i-- {
			(*s)[i+numNew] = (*s)[i]
		}

		// Copy new elements at the index position
		copy((*s)[index:], elements)
	} else {
		// Not enough capacity: allocate new slice with growth
		newCap := newLen
		if newCap < 2*cap(*s) {
			newCap = 2 * cap(*s)
		}

		newSlice := make([]T, newLen, newCap)
		copy(newSlice, (*s)[:index])
		copy(newSlice[index:], elements)
		copy(newSlice[index+numNew:], (*s)[index:])
		*s = newSlice
	}
}

// Remove removes the element at the given index (mutates in place)
func (s *Slice[T]) Remove(index int) {
	if index < 0 || index >= len(*s) {
		return
	}
	*s = append((*s)[:index], (*s)[index+1:]...)
}

// Set sets the element at index to value (mutates in place)
func (s *Slice[T]) Set(index int, value T) {
	if index < 0 || index >= len(*s) {
		return
	}
	(*s)[index] = value
}

// Clone creates a shallow copy of the slice
func (s *Slice[T]) Clone() Slice[T] {
	if *s == nil {
		return nil
	}
	result := make(Slice[T], len(*s))
	copy(result, *s)
	return result
}

// ============================================================================
// Transformation Operations (mutate in place)
// ============================================================================

// Reverse reverses the slice in place
func (s *Slice[T]) Reverse() {
	for i, j := 0, len(*s)-1; i < j; i, j = i+1, j-1 {
		(*s)[i], (*s)[j] = (*s)[j], (*s)[i]
	}
}

// ============================================================================
// Search Operations
// ============================================================================

// Find returns the first element satisfying the predicate
func (s *Slice[T]) Find(predicate func(T) bool) (T, bool) {
	for _, v := range *s {
		if predicate(v) {
			return v, true
		}
	}
	var zero T
	return zero, false
}

// Count returns the number of elements satisfying the predicate
func (s *Slice[T]) Count(predicate func(T) bool) int {
	count := 0
	for _, v := range *s {
		if predicate(v) {
			count++
		}
	}
	return count
}

// ============================================================================
// String Operations (for string slices)
// ============================================================================

// Join concatenates string slice elements with a separator
func Join(s Slice[string], separator string) string {
	if len(s) == 0 {
		return ""
	}
	if len(s) == 1 {
		return s[0]
	}
	n := len(separator) * (len(s) - 1)
	for _, elem := range s {
		n += len(elem)
	}
	result := make([]byte, 0, n)
	result = append(result, s[0]...)
	for _, elem := range s[1:] {
		result = append(result, separator...)
		result = append(result, elem...)
	}
	return string(result)
}
