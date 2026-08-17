// Package types - Slice Implementation
//
// ARCHITECTURE OVERVIEW:
//
// Slice provides a type-safe, generic wrapper around Go's native slice type that
// extends functionality while maintaining full compatibility with built-in slice
// operations. Unlike struct-based wrappers, Slice is a type definition that
// preserves native syntax (indexing, slicing, ranging) while adding safety and
// convenience methods through receiver functions.
//
// DESIGN PHILOSOPHY:
//
// Zero-Cost Abstraction:
// Slice[T] is defined as []T, not struct{data []T}, which means:
// 1. No indirection - direct memory access
// 2. No wrapper allocation - slice header only
// 3. Native syntax preserved - s[i], s[1:3], range s all work
// 4. Compiler optimizations apply - bounds check elimination, etc.
//
// This design trades encapsulation for performance and ergonomics, making Slice
// ideal for performance-critical code that still needs safety guarantees.
//
// MEMORY MODEL:
//
// Go Slice Internal Structure:
//
//	┌──────────────────────┐
//	│  SliceHeader         │
//	│  ├── Data *T         │ ← Pointer to backing array
//	│  ├── Len  int        │ ← Current length
//	│  └── Cap  int        │ ← Total capacity
//	└──────────────────────┘
//	         ↓
//	┌──────────────────────────────────┐
//	│  Backing Array                   │
//	│  [T][T][T][T][T]...[T][T][ ][ ]  │
//	│  └─────Len─────┘   └──Cap────┘   │
//	└──────────────────────────────────┘
//
// ARCHITECTURAL DECISIONS:
//
// 1. Type Definition vs Struct:
//   - Type definition: type Slice[T] []T
//   - Preserves all native operations
//   - Zero memory overhead
//   - Direct compatibility with []T
//
// 2. Pointer Receivers:
//   - Mutation methods use *Slice[T]
//   - Allows in-place modifications
//   - Avoids slice header copying
//   - Consistent with append semantics
//
// 3. Safe Access Pattern:
//   - Three-tier safety: Get (panic), GetOr (default), GetSafe (bool)
//   - Follows Go idioms (value, ok pattern)
//   - Predictable error handling
//
// PERFORMANCE OPTIMIZATIONS:
//
//  1. Prepend Optimization:
//     Traditional prepend: O(n) allocation + copy
//     Optimized prepend: O(n) shift when capacity sufficient
//
//     When cap(s) >= newLen:
//     - Extend slice length
//     - Shift existing elements right
//     - Copy new elements to front
//     - Avoids allocation entirely
//
//  2. Insert Optimization:
//     Similar to prepend but for arbitrary index
//     - In-place shift when capacity allows
//     - Single allocation when growth needed
//     - Minimizes memory fragmentation
//
//  3. Growth Strategy:
//     Follows Go's slice growth algorithm:
//     - < 1024: double capacity
//     - >= 1024: grow by 25%
//     - Amortized O(1) append
//
// CAPACITY MANAGEMENT:
//
// Slice capacity behavior:
//
//	s := WithCapacity[int](10)  // len=0, cap=10
//	s.Append(1)                  // len=1, cap=10 (no alloc)
//	s.Prepend(0)                 // len=2, cap=10 (no alloc)
//	s.Insert(1, 5)               // len=3, cap=10 (no alloc)
//
// Growth triggers:
//
//	s := New(1, 2, 3)            // len=3, cap=3
//	s.Append(4)                  // len=4, cap=6 (growth)
//	s.Prepend(0)                 // len=5, cap=6 (no growth)
//	s.Insert(2, 9, 8, 7)         // len=8, cap=12 (growth)
//
// ALGORITHMIC COMPLEXITY:
//
// Time Complexity:
// - Get/Set: O(1) with bounds check
// - Append: O(1) amortized
// - Prepend: O(n) shift, O(1) amortized allocation
// - Insert: O(n) shift from index, O(1) amortized allocation
// - Remove: O(n) shift after index
// - Clone: O(n) copy
// - Reverse: O(n/2) swaps
// - Find: O(n) linear search
//
// Space Complexity:
// - Operations: O(1) additional space
// - Clone: O(n) new allocation
// - Growth: 2x capacity when < 1024 elements
//
// CONCURRENCY CONSIDERATIONS:
//
// Slice is NOT thread-safe. Concurrent access requires synchronization:
//
// Option 1 - Mutex protection:
//
//	type SafeSlice[T any] struct {
//	    mu sync.RWMutex
//	    s  Slice[T]
//	}
//
// Option 2 - Channel-based access:
//
//	type ChanSlice[T any] struct {
//	    ops chan func(*Slice[T])
//	}
//
// Option 3 - Copy-on-write:
//
//	type COWSlice[T any] struct {
//	    s atomic.Value // holds Slice[T]
//	}
//
// USE CASES:
//
// Ideal for:
// - Dynamic arrays with frequent appends
// - Collections requiring safe access
// - Buffer management and pooling
// - String manipulation (Slice[string])
// - Algorithm implementations
//
// Not ideal for:
// - Fixed-size arrays (use [N]T)
// - Concurrent access (needs sync)
// - Very large slices (consider chunking)
// - Sorted data (consider container/heap)
//
// USAGE PATTERNS:
//
// Safe iteration with modification:
//
//	s := New(1, 2, 3, 4, 5)
//	for i := s.Len() - 1; i >= 0; i-- {
//	    if s.Get(i)%2 == 0 {
//	        s.Remove(i) // Safe backward iteration
//	    }
//	}
//
// Builder pattern:
//
//	s := WithCapacity[string](100)
//	s.Append("header")
//	for _, item := range items {
//	    s.Append(process(item))
//	}
//	s.Append("footer")
//
// Functional operations:
//
//	s := New(1, 2, 3, 4, 5)
//	if val, ok := s.Find(func(x int) bool {
//	    return x > 3
//	}); ok {
//	    fmt.Printf("Found: %d\n", val)
//	}
//
// COMPARISON WITH ALTERNATIVES:
//
// vs []T (native slice):
// - Slice[T]: Safe access methods, convenience functions
// - []T: Minimal overhead, standard library compatible
// - Use Slice[T] when safety and convenience matter
//
// vs container/list:
// - Slice[T]: O(1) random access, cache-friendly
// - list.List: O(1) insertion/deletion at any position
// - Use list for frequent mid-sequence modifications
//
// vs sync.Pool:
// - Slice[T]: Owned memory, predictable lifecycle
// - sync.Pool: Shared buffers, GC-managed
// - Use Pool for temporary buffers
//
// BEST PRACTICES:
//
//  1. Pre-allocate capacity:
//     // Bad: Multiple growths
//     s := New[string]()
//     for _, v := range largeData {
//     s.Append(v) // Many allocations
//     }
//
//     // Good: Single allocation
//     s := WithCapacity[string](len(largeData))
//     for _, v := range largeData {
//     s.Append(v) // No allocations
//     }
//
//  2. Use appropriate access method:
//     // Panic on invalid index
//     val := s.Get(unknownIndex)
//
//     // Safe with default
//     val := s.GetOr(unknownIndex, defaultVal)
//
//     // Check existence
//     if val, ok := s.GetSafe(unknownIndex); ok {
//     use(val)
//     }
//
//  3. Avoid slice aliasing:
//     // Bad: Shared backing array
//     s1 := New(1, 2, 3, 4, 5)
//     s2 := s1[1:3] // Shares backing array
//
//     // Good: Independent copy
//     s2 := s1.Clone()[1:3]
//
// PITFALLS:
//
//  1. Append during iteration:
//     // Bad: Infinite loop possible
//     for i := 0; i < s.Len(); i++ {
//     s.Append(s.Get(i)) // Length changes!
//     }
//
//     // Good: Fixed iteration count
//     originalLen := s.Len()
//     for i := 0; i < originalLen; i++ {
//     s.Append(s.Get(i))
//     }
//
//  2. Capacity assumptions:
//     s := WithCapacity[int](10)
//     s = s[:10] // PANIC! Length exceeds current length
//
//     // Correct: Set length explicitly
//     s := Make[int](10, 10) // len=10, cap=10
//
// FUTURE ENHANCEMENTS:
//
// - Parallel operations (Map, Filter, Reduce)
// - Sort with custom comparators
// - Binary search for sorted slices
// - Batch operations (RemoveAll, ReplaceAll)
// - Memory pool integration
package types

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
