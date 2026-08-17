package types

import (
	"testing"
)

func TestSliceNew(t *testing.T) {
	t.Run("creates empty slice", func(t *testing.T) {
		s := New[int]()

		if s.Len() != 0 {
			t.Errorf("expected length 0, got %d", s.Len())
		}
	})

	t.Run("creates slice with elements", func(t *testing.T) {
		s := New(1, 2, 3)

		if s.Len() != 3 {
			t.Errorf("expected length 3, got %d", s.Len())
		}
		if s.Get(0) != 1 || s.Get(1) != 2 || s.Get(2) != 3 {
			t.Error("elements not correctly initialized")
		}
	})

	t.Run("creates slice with strings", func(t *testing.T) {
		s := New("a", "b", "c")

		if s.Len() != 3 {
			t.Errorf("expected length 3, got %d", s.Len())
		}
		if s.Get(0) != "a" || s.Get(1) != "b" || s.Get(2) != "c" {
			t.Error("elements not correctly initialized")
		}
	})
}

func TestSliceFrom(t *testing.T) {
	t.Run("creates from existing slice", func(t *testing.T) {
		original := []int{1, 2, 3}
		s := From(original)

		if s.Len() != 3 {
			t.Errorf("expected length 3, got %d", s.Len())
		}
	})

	t.Run("creates from nil slice", func(t *testing.T) {
		var original []int
		s := From(original)

		if s.Len() != 0 {
			t.Errorf("expected length 0, got %d", s.Len())
		}
	})
}

func TestSliceMake(t *testing.T) {
	t.Run("creates with length and capacity", func(t *testing.T) {
		s := Make[int](5, 10)

		if s.Len() != 5 {
			t.Errorf("expected length 5, got %d", s.Len())
		}
		if s.Cap() != 10 {
			t.Errorf("expected capacity 10, got %d", s.Cap())
		}
	})
}

func TestSliceWithLength(t *testing.T) {
	t.Run("creates with specified length", func(t *testing.T) {
		s := WithLength[int](5)

		if s.Len() != 5 {
			t.Errorf("expected length 5, got %d", s.Len())
		}
	})

	t.Run("elements are zero values", func(t *testing.T) {
		s := WithLength[int](3)

		for i := 0; i < s.Len(); i++ {
			if s.Get(i) != 0 {
				t.Errorf("expected zero value at index %d", i)
			}
		}
	})
}

func TestSliceWithCapacity(t *testing.T) {
	t.Run("creates with zero length and specified capacity", func(t *testing.T) {
		s := WithCapacity[int](10)

		if s.Len() != 0 {
			t.Errorf("expected length 0, got %d", s.Len())
		}
		if s.Cap() != 10 {
			t.Errorf("expected capacity 10, got %d", s.Cap())
		}
	})
}

func TestSliceLen(t *testing.T) {
	t.Run("returns correct length", func(t *testing.T) {
		s := New(1, 2, 3, 4, 5)

		if s.Len() != 5 {
			t.Errorf("expected length 5, got %d", s.Len())
		}
	})

	t.Run("returns 0 for empty slice", func(t *testing.T) {
		s := New[int]()

		if s.Len() != 0 {
			t.Errorf("expected length 0, got %d", s.Len())
		}
	})
}

func TestSliceCap(t *testing.T) {
	t.Run("returns correct capacity", func(t *testing.T) {
		s := Make[int](3, 10)

		if s.Cap() != 10 {
			t.Errorf("expected capacity 10, got %d", s.Cap())
		}
	})
}

func TestSliceIsEmpty(t *testing.T) {
	t.Run("returns true for empty slice", func(t *testing.T) {
		s := New[int]()

		if !s.IsEmpty() {
			t.Error("expected IsEmpty to return true")
		}
	})

	t.Run("returns false for non-empty slice", func(t *testing.T) {
		s := New(1)

		if s.IsEmpty() {
			t.Error("expected IsEmpty to return false")
		}
	})
}

func TestSliceIsNotEmpty(t *testing.T) {
	t.Run("returns false for empty slice", func(t *testing.T) {
		s := New[int]()

		if s.IsNotEmpty() {
			t.Error("expected IsNotEmpty to return false")
		}
	})

	t.Run("returns true for non-empty slice", func(t *testing.T) {
		s := New(1)

		if !s.IsNotEmpty() {
			t.Error("expected IsNotEmpty to return true")
		}
	})
}

func TestSliceGet(t *testing.T) {
	t.Run("returns element at index", func(t *testing.T) {
		s := New(10, 20, 30)

		if s.Get(0) != 10 {
			t.Errorf("expected 10, got %d", s.Get(0))
		}
		if s.Get(1) != 20 {
			t.Errorf("expected 20, got %d", s.Get(1))
		}
		if s.Get(2) != 30 {
			t.Errorf("expected 30, got %d", s.Get(2))
		}
	})

	t.Run("panics on out of bounds", func(t *testing.T) {
		s := New(1, 2, 3)

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic on out of bounds access")
			}
		}()

		_ = s.Get(10)
	})
}

func TestSliceGetOr(t *testing.T) {
	t.Run("returns element at valid index", func(t *testing.T) {
		s := New(10, 20, 30)

		if s.GetOr(1, 0) != 20 {
			t.Errorf("expected 20, got %d", s.GetOr(1, 0))
		}
	})

	t.Run("returns default for negative index", func(t *testing.T) {
		s := New(10, 20, 30)

		if s.GetOr(-1, 99) != 99 {
			t.Errorf("expected 99, got %d", s.GetOr(-1, 99))
		}
	})

	t.Run("returns default for out of bounds index", func(t *testing.T) {
		s := New(10, 20, 30)

		if s.GetOr(10, 99) != 99 {
			t.Errorf("expected 99, got %d", s.GetOr(10, 99))
		}
	})
}

func TestSliceGetSafe(t *testing.T) {
	t.Run("returns element and true for valid index", func(t *testing.T) {
		s := New(10, 20, 30)

		val, ok := s.GetSafe(1)
		if !ok {
			t.Error("expected ok to be true")
		}
		if val != 20 {
			t.Errorf("expected 20, got %d", val)
		}
	})

	t.Run("returns zero and false for invalid index", func(t *testing.T) {
		s := New(10, 20, 30)

		val, ok := s.GetSafe(10)
		if ok {
			t.Error("expected ok to be false")
		}
		if val != 0 {
			t.Errorf("expected 0, got %d", val)
		}
	})

	t.Run("returns zero and false for negative index", func(t *testing.T) {
		s := New(10, 20, 30)

		val, ok := s.GetSafe(-1)
		if ok {
			t.Error("expected ok to be false")
		}
		if val != 0 {
			t.Errorf("expected 0, got %d", val)
		}
	})
}

func TestSliceFirst(t *testing.T) {
	t.Run("returns first element", func(t *testing.T) {
		s := New(10, 20, 30)

		if s.First() != 10 {
			t.Errorf("expected 10, got %d", s.First())
		}
	})

	t.Run("panics on empty slice", func(t *testing.T) {
		s := New[int]()

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic on empty slice")
			}
		}()

		_ = s.First()
	})
}

func TestSliceFirstOr(t *testing.T) {
	t.Run("returns first element for non-empty slice", func(t *testing.T) {
		s := New(10, 20, 30)

		if s.FirstOr(99) != 10 {
			t.Errorf("expected 10, got %d", s.FirstOr(99))
		}
	})

	t.Run("returns default for empty slice", func(t *testing.T) {
		s := New[int]()

		if s.FirstOr(99) != 99 {
			t.Errorf("expected 99, got %d", s.FirstOr(99))
		}
	})
}

func TestSliceFirstSafe(t *testing.T) {
	t.Run("returns first element and true for non-empty slice", func(t *testing.T) {
		s := New(10, 20, 30)

		val, ok := s.FirstSafe()
		if !ok {
			t.Error("expected ok to be true")
		}
		if val != 10 {
			t.Errorf("expected 10, got %d", val)
		}
	})

	t.Run("returns zero and false for empty slice", func(t *testing.T) {
		s := New[int]()

		val, ok := s.FirstSafe()
		if ok {
			t.Error("expected ok to be false")
		}
		if val != 0 {
			t.Errorf("expected 0, got %d", val)
		}
	})
}

func TestSliceLast(t *testing.T) {
	t.Run("returns last element", func(t *testing.T) {
		s := New(10, 20, 30)

		if s.Last() != 30 {
			t.Errorf("expected 30, got %d", s.Last())
		}
	})

	t.Run("panics on empty slice", func(t *testing.T) {
		s := New[int]()

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected panic on empty slice")
			}
		}()

		_ = s.Last()
	})
}

func TestSliceLastOr(t *testing.T) {
	t.Run("returns last element for non-empty slice", func(t *testing.T) {
		s := New(10, 20, 30)

		if s.LastOr(99) != 30 {
			t.Errorf("expected 30, got %d", s.LastOr(99))
		}
	})

	t.Run("returns default for empty slice", func(t *testing.T) {
		s := New[int]()

		if s.LastOr(99) != 99 {
			t.Errorf("expected 99, got %d", s.LastOr(99))
		}
	})
}

func TestSliceLastSafe(t *testing.T) {
	t.Run("returns last element and true for non-empty slice", func(t *testing.T) {
		s := New(10, 20, 30)

		val, ok := s.LastSafe()
		if !ok {
			t.Error("expected ok to be true")
		}
		if val != 30 {
			t.Errorf("expected 30, got %d", val)
		}
	})

	t.Run("returns zero and false for empty slice", func(t *testing.T) {
		s := New[int]()

		val, ok := s.LastSafe()
		if ok {
			t.Error("expected ok to be false")
		}
		if val != 0 {
			t.Errorf("expected 0, got %d", val)
		}
	})
}

func TestSliceAppend(t *testing.T) {
	t.Run("appends single element", func(t *testing.T) {
		s := New(1, 2, 3)

		s.Append(4)

		if s.Len() != 4 {
			t.Errorf("expected length 4, got %d", s.Len())
		}
		if s.Last() != 4 {
			t.Errorf("expected last element 4, got %d", s.Last())
		}
	})

	t.Run("appends multiple elements", func(t *testing.T) {
		s := New(1, 2, 3)

		s.Append(4, 5, 6)

		if s.Len() != 6 {
			t.Errorf("expected length 6, got %d", s.Len())
		}
	})

	t.Run("appends to empty slice", func(t *testing.T) {
		s := New[int]()

		s.Append(1, 2, 3)

		if s.Len() != 3 {
			t.Errorf("expected length 3, got %d", s.Len())
		}
	})
}

func TestSlicePrepend(t *testing.T) {
	t.Run("prepends single element", func(t *testing.T) {
		s := New(2, 3, 4)

		s.Prepend(1)

		if s.Len() != 4 {
			t.Errorf("expected length 4, got %d", s.Len())
		}
		if s.First() != 1 {
			t.Errorf("expected first element 1, got %d", s.First())
		}
	})

	t.Run("prepends multiple elements", func(t *testing.T) {
		s := New(4, 5, 6)

		s.Prepend(1, 2, 3)

		if s.Len() != 6 {
			t.Errorf("expected length 6, got %d", s.Len())
		}
		if s.First() != 1 {
			t.Errorf("expected first element 1, got %d", s.First())
		}
	})

	t.Run("prepends to empty slice", func(t *testing.T) {
		s := New[int]()

		s.Prepend(1, 2, 3)

		if s.Len() != 3 {
			t.Errorf("expected length 3, got %d", s.Len())
		}
	})
}

func TestSliceConcat(t *testing.T) {
	t.Run("concatenates single slice", func(t *testing.T) {
		s := New(1, 2)
		other := New(3, 4)

		s.Concat(other)

		if s.Len() != 4 {
			t.Errorf("expected length 4, got %d", s.Len())
		}
	})

	t.Run("concatenates multiple slices", func(t *testing.T) {
		s := New(1, 2)
		other1 := New(3, 4)
		other2 := New(5, 6)

		s.Concat(other1, other2)

		if s.Len() != 6 {
			t.Errorf("expected length 6, got %d", s.Len())
		}
	})

	t.Run("concatenates empty slices", func(t *testing.T) {
		s := New(1, 2)
		other := New[int]()

		s.Concat(other)

		if s.Len() != 2 {
			t.Errorf("expected length 2, got %d", s.Len())
		}
	})
}

func TestSliceInsert(t *testing.T) {
	t.Run("inserts at beginning", func(t *testing.T) {
		s := New(2, 3, 4)

		s.Insert(0, 1)

		if s.Len() != 4 {
			t.Errorf("expected length 4, got %d", s.Len())
		}
		if s.First() != 1 {
			t.Errorf("expected first element 1, got %d", s.First())
		}
	})

	t.Run("inserts in middle", func(t *testing.T) {
		s := New(1, 2, 4, 5)

		s.Insert(2, 3)

		if s.Len() != 5 {
			t.Errorf("expected length 5, got %d", s.Len())
		}
		if s.Get(2) != 3 {
			t.Errorf("expected element at index 2 to be 3, got %d", s.Get(2))
		}
	})

	t.Run("inserts at end", func(t *testing.T) {
		s := New(1, 2, 3)

		s.Insert(3, 4)

		if s.Len() != 4 {
			t.Errorf("expected length 4, got %d", s.Len())
		}
		if s.Last() != 4 {
			t.Errorf("expected last element 4, got %d", s.Last())
		}
	})

	t.Run("inserts multiple elements", func(t *testing.T) {
		s := New(1, 5)

		s.Insert(1, 2, 3, 4)

		if s.Len() != 5 {
			t.Errorf("expected length 5, got %d", s.Len())
		}
	})

	t.Run("handles negative index", func(t *testing.T) {
		s := New(2, 3)

		s.Insert(-5, 1)

		if s.First() != 1 {
			t.Errorf("expected first element 1, got %d", s.First())
		}
	})

	t.Run("handles index beyond length", func(t *testing.T) {
		s := New(1, 2)

		s.Insert(100, 3)

		if s.Last() != 3 {
			t.Errorf("expected last element 3, got %d", s.Last())
		}
	})
}

func TestSliceRemove(t *testing.T) {
	t.Run("removes element at index", func(t *testing.T) {
		s := New(1, 2, 3, 4, 5)

		s.Remove(2)

		if s.Len() != 4 {
			t.Errorf("expected length 4, got %d", s.Len())
		}
		if s.Get(2) != 4 {
			t.Errorf("expected element at index 2 to be 4, got %d", s.Get(2))
		}
	})

	t.Run("removes first element", func(t *testing.T) {
		s := New(1, 2, 3)

		s.Remove(0)

		if s.Len() != 2 {
			t.Errorf("expected length 2, got %d", s.Len())
		}
		if s.First() != 2 {
			t.Errorf("expected first element 2, got %d", s.First())
		}
	})

	t.Run("removes last element", func(t *testing.T) {
		s := New(1, 2, 3)

		s.Remove(2)

		if s.Len() != 2 {
			t.Errorf("expected length 2, got %d", s.Len())
		}
		if s.Last() != 2 {
			t.Errorf("expected last element 2, got %d", s.Last())
		}
	})

	t.Run("does nothing for negative index", func(t *testing.T) {
		s := New(1, 2, 3)

		s.Remove(-1)

		if s.Len() != 3 {
			t.Errorf("expected length 3, got %d", s.Len())
		}
	})

	t.Run("does nothing for out of bounds index", func(t *testing.T) {
		s := New(1, 2, 3)

		s.Remove(10)

		if s.Len() != 3 {
			t.Errorf("expected length 3, got %d", s.Len())
		}
	})
}

func TestSliceSet(t *testing.T) {
	t.Run("sets element at valid index", func(t *testing.T) {
		s := New(1, 2, 3)

		s.Set(1, 20)

		if s.Get(1) != 20 {
			t.Errorf("expected 20, got %d", s.Get(1))
		}
	})

	t.Run("does nothing for negative index", func(t *testing.T) {
		s := New(1, 2, 3)

		s.Set(-1, 99)

		if s.Get(0) != 1 && s.Get(1) != 2 && s.Get(2) != 3 {
			t.Error("slice should not be modified")
		}
	})

	t.Run("does nothing for out of bounds index", func(t *testing.T) {
		s := New(1, 2, 3)

		s.Set(10, 99)

		if s.Len() != 3 {
			t.Errorf("expected length 3, got %d", s.Len())
		}
	})
}

func TestSliceClone(t *testing.T) {
	t.Run("creates copy of slice", func(t *testing.T) {
		s := New(1, 2, 3)

		cloned := s.Clone()

		if cloned.Len() != 3 {
			t.Errorf("expected length 3, got %d", cloned.Len())
		}
	})

	t.Run("creates independent copy", func(t *testing.T) {
		s := New(1, 2, 3)

		cloned := s.Clone()
		cloned.Set(0, 100)

		if s.Get(0) != 1 {
			t.Error("modifying clone should not affect original")
		}
		if cloned.Get(0) != 100 {
			t.Error("clone should have modified value")
		}
	})

	t.Run("handles nil slice", func(t *testing.T) {
		var s Slice[int]

		cloned := s.Clone()

		if cloned != nil {
			t.Error("cloning nil slice should return nil")
		}
	})

	t.Run("handles empty slice", func(t *testing.T) {
		s := New[int]()

		cloned := s.Clone()

		if cloned.Len() != 0 {
			t.Errorf("expected length 0, got %d", cloned.Len())
		}
	})
}

func TestSliceReverse(t *testing.T) {
	t.Run("reverses slice", func(t *testing.T) {
		s := New(1, 2, 3, 4, 5)

		s.Reverse()

		expected := []int{5, 4, 3, 2, 1}
		for i, v := range expected {
			if s.Get(i) != v {
				t.Errorf("expected %d at index %d, got %d", v, i, s.Get(i))
			}
		}
	})

	t.Run("reverses single element", func(t *testing.T) {
		s := New(1)

		s.Reverse()

		if s.Get(0) != 1 {
			t.Errorf("expected 1, got %d", s.Get(0))
		}
	})

	t.Run("reverses two elements", func(t *testing.T) {
		s := New(1, 2)

		s.Reverse()

		if s.Get(0) != 2 || s.Get(1) != 1 {
			t.Error("expected reversed order")
		}
	})

	t.Run("handles empty slice", func(t *testing.T) {
		s := New[int]()

		s.Reverse() // should not panic

		if s.Len() != 0 {
			t.Errorf("expected length 0, got %d", s.Len())
		}
	})
}

func TestSliceFind(t *testing.T) {
	t.Run("finds element matching predicate", func(t *testing.T) {
		s := New(1, 2, 3, 4, 5)

		val, found := s.Find(func(v int) bool { return v > 3 })

		if !found {
			t.Error("expected to find element")
		}
		if val != 4 {
			t.Errorf("expected 4, got %d", val)
		}
	})

	t.Run("returns false when not found", func(t *testing.T) {
		s := New(1, 2, 3)

		val, found := s.Find(func(v int) bool { return v > 10 })

		if found {
			t.Error("expected not to find element")
		}
		if val != 0 {
			t.Errorf("expected zero value, got %d", val)
		}
	})

	t.Run("returns first matching element", func(t *testing.T) {
		s := New(1, 2, 3, 4, 5)

		val, found := s.Find(func(v int) bool { return v%2 == 0 })

		if !found {
			t.Error("expected to find element")
		}
		if val != 2 {
			t.Errorf("expected 2 (first even), got %d", val)
		}
	})
}

func TestSliceCount(t *testing.T) {
	t.Run("counts elements matching predicate", func(t *testing.T) {
		s := New(1, 2, 3, 4, 5, 6)

		count := s.Count(func(v int) bool { return v%2 == 0 })

		if count != 3 {
			t.Errorf("expected 3 even numbers, got %d", count)
		}
	})

	t.Run("returns 0 when none match", func(t *testing.T) {
		s := New(1, 3, 5, 7)

		count := s.Count(func(v int) bool { return v%2 == 0 })

		if count != 0 {
			t.Errorf("expected 0, got %d", count)
		}
	})

	t.Run("counts all elements when all match", func(t *testing.T) {
		s := New(2, 4, 6, 8)

		count := s.Count(func(v int) bool { return v%2 == 0 })

		if count != 4 {
			t.Errorf("expected 4, got %d", count)
		}
	})
}

func TestJoin(t *testing.T) {
	t.Run("joins with separator", func(t *testing.T) {
		s := New("a", "b", "c")

		result := Join(s, ", ")

		if result != "a, b, c" {
			t.Errorf("expected 'a, b, c', got '%s'", result)
		}
	})

	t.Run("joins with empty separator", func(t *testing.T) {
		s := New("a", "b", "c")

		result := Join(s, "")

		if result != "abc" {
			t.Errorf("expected 'abc', got '%s'", result)
		}
	})

	t.Run("returns empty string for empty slice", func(t *testing.T) {
		s := New[string]()

		result := Join(s, ", ")

		if result != "" {
			t.Errorf("expected empty string, got '%s'", result)
		}
	})

	t.Run("returns element for single element slice", func(t *testing.T) {
		s := New("only")

		result := Join(s, ", ")

		if result != "only" {
			t.Errorf("expected 'only', got '%s'", result)
		}
	})
}

func TestSliceReslicing(t *testing.T) {
	t.Run("supports standard reslicing", func(t *testing.T) {
		s := New(1, 2, 3, 4, 5)

		sub := s[1:4]

		if len(sub) != 3 {
			t.Errorf("expected length 3, got %d", len(sub))
		}
		if sub[0] != 2 || sub[1] != 3 || sub[2] != 4 {
			t.Error("resliced elements incorrect")
		}
	})

	t.Run("supports range iteration", func(t *testing.T) {
		s := New(1, 2, 3)
		sum := 0

		for _, v := range s {
			sum += v
		}

		if sum != 6 {
			t.Errorf("expected sum 6, got %d", sum)
		}
	})
}

func TestSliceOperationsSequence(t *testing.T) {
	t.Run("complex operations sequence", func(t *testing.T) {
		s := New(1, 2, 3)

		// Append
		s.Append(4, 5)
		if s.Len() != 5 {
			t.Errorf("after append: expected length 5, got %d", s.Len())
		}

		// Prepend
		s.Prepend(0)
		if s.First() != 0 {
			t.Errorf("after prepend: expected first 0, got %d", s.First())
		}

		// Insert
		s.Insert(3, 100)
		if s.Get(3) != 100 {
			t.Errorf("after insert: expected 100 at index 3, got %d", s.Get(3))
		}

		// Remove
		s.Remove(3)
		if s.Get(3) != 3 {
			t.Errorf("after remove: expected 3 at index 3, got %d", s.Get(3))
		}

		// Set
		s.Set(0, 99)
		if s.First() != 99 {
			t.Errorf("after set: expected first 99, got %d", s.First())
		}

		// Reverse
		originalLast := s.Last()
		s.Reverse()
		if s.First() != originalLast {
			t.Error("after reverse: first should be original last")
		}
	})
}

func TestSliceEdgeCases(t *testing.T) {
	t.Run("handles large slice", func(t *testing.T) {
		s := WithCapacity[int](10000)

		for i := 0; i < 10000; i++ {
			s.Append(i)
		}

		if s.Len() != 10000 {
			t.Errorf("expected length 10000, got %d", s.Len())
		}
		if s.Last() != 9999 {
			t.Errorf("expected last 9999, got %d", s.Last())
		}
	})

	t.Run("handles struct elements", func(t *testing.T) {
		type item struct {
			ID   int
			Name string
		}

		s := New(
			item{1, "first"},
			item{2, "second"},
			item{3, "third"},
		)

		if s.Len() != 3 {
			t.Errorf("expected length 3, got %d", s.Len())
		}

		found, ok := s.Find(func(i item) bool { return i.ID == 2 })
		if !ok || found.Name != "second" {
			t.Error("expected to find item with ID 2")
		}
	})

	t.Run("handles pointer elements", func(t *testing.T) {
		a, b, c := 1, 2, 3
		s := New(&a, &b, &c)

		if s.Len() != 3 {
			t.Errorf("expected length 3, got %d", s.Len())
		}
		if *s.Get(0) != 1 {
			t.Errorf("expected 1, got %d", *s.Get(0))
		}
	})
}
