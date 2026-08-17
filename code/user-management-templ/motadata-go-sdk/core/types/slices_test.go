package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceNew(t *testing.T) {
	t.Run("creates empty slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New[int]()

		assertions.Equal(0, s.Len())
	})

	t.Run("creates slice with elements", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3)

		assertions.Equal(3, s.Len())
		assertions.Equal(1, s.Get(0))
		assertions.Equal(2, s.Get(1))
		assertions.Equal(3, s.Get(2))
	})

	t.Run("creates slice with strings", func(t *testing.T) {
		assertions := assert.New(t)
		s := New("a", "b", "c")

		assertions.Equal(3, s.Len())
		assertions.Equal("a", s.Get(0))
		assertions.Equal("b", s.Get(1))
		assertions.Equal("c", s.Get(2))
	})
}

func TestSliceFrom(t *testing.T) {
	t.Run("creates from existing slice", func(t *testing.T) {
		assertions := assert.New(t)
		original := []int{1, 2, 3}
		s := From(original)

		assertions.Equal(3, s.Len())
	})

	t.Run("creates from nil slice", func(t *testing.T) {
		assertions := assert.New(t)
		var original []int
		s := From(original)

		assertions.Equal(0, s.Len())
	})
}

func TestSliceMake(t *testing.T) {
	t.Run("creates with length and capacity", func(t *testing.T) {
		assertions := assert.New(t)
		s := Make[int](5, 10)

		assertions.Equal(5, s.Len())
		assertions.Equal(10, s.Cap())
	})
}

func TestSliceWithLength(t *testing.T) {
	t.Run("creates with specified length", func(t *testing.T) {
		assertions := assert.New(t)
		s := WithLength[int](5)

		assertions.Equal(5, s.Len())
	})

	t.Run("elements are zero values", func(t *testing.T) {
		assertions := assert.New(t)
		s := WithLength[int](3)

		for i := 0; i < s.Len(); i++ {
			assertions.Equal(0, s.Get(i))
		}
	})
}

func TestSliceWithCapacity(t *testing.T) {
	t.Run("creates with zero length and specified capacity", func(t *testing.T) {
		assertions := assert.New(t)
		s := WithCapacity[int](10)

		assertions.Equal(0, s.Len())
		assertions.Equal(10, s.Cap())
	})
}

func TestSliceLen(t *testing.T) {
	t.Run("returns correct length", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3, 4, 5)

		assertions.Equal(5, s.Len())
	})

	t.Run("returns 0 for empty slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New[int]()

		assertions.Equal(0, s.Len())
	})
}

func TestSliceCap(t *testing.T) {
	t.Run("returns correct capacity", func(t *testing.T) {
		assertions := assert.New(t)
		s := Make[int](3, 10)

		assertions.Equal(10, s.Cap())
	})
}

func TestSliceIsEmpty(t *testing.T) {
	t.Run("returns true for empty slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New[int]()

		assertions.True(s.IsEmpty())
	})

	t.Run("returns false for non-empty slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1)

		assertions.False(s.IsEmpty())
	})
}

func TestSliceIsNotEmpty(t *testing.T) {
	t.Run("returns false for empty slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New[int]()

		assertions.False(s.IsNotEmpty())
	})

	t.Run("returns true for non-empty slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1)

		assertions.True(s.IsNotEmpty())
	})
}

func TestSliceGet(t *testing.T) {
	t.Run("returns element at index", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(10, 20, 30)

		assertions.Equal(10, s.Get(0))
		assertions.Equal(20, s.Get(1))
		assertions.Equal(30, s.Get(2))
	})

	t.Run("panics on out of bounds", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3)

		assertions.Panics(func() {
			_ = s.Get(10)
		})
	})
}

func TestSliceGetOr(t *testing.T) {
	t.Run("returns element at valid index", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(10, 20, 30)

		assertions.Equal(20, s.GetOr(1, 0))
	})

	t.Run("returns default for negative index", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(10, 20, 30)

		assertions.Equal(99, s.GetOr(-1, 99))
	})

	t.Run("returns default for out of bounds index", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(10, 20, 30)

		assertions.Equal(99, s.GetOr(10, 99))
	})
}

func TestSliceGetSafe(t *testing.T) {
	t.Run("returns element and true for valid index", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(10, 20, 30)

		val, ok := s.GetSafe(1)
		assertions.True(ok)
		assertions.Equal(20, val)
	})

	t.Run("returns zero and false for invalid index", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(10, 20, 30)

		val, ok := s.GetSafe(10)
		assertions.False(ok)
		assertions.Equal(0, val)
	})

	t.Run("returns zero and false for negative index", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(10, 20, 30)

		val, ok := s.GetSafe(-1)
		assertions.False(ok)
		assertions.Equal(0, val)
	})
}

func TestSliceFirst(t *testing.T) {
	t.Run("returns first element", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(10, 20, 30)

		assertions.Equal(10, s.First())
	})

	t.Run("panics on empty slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New[int]()

		assertions.Panics(func() {
			_ = s.First()
		})
	})
}

func TestSliceFirstOr(t *testing.T) {
	t.Run("returns first element for non-empty slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(10, 20, 30)

		assertions.Equal(10, s.FirstOr(99))
	})

	t.Run("returns default for empty slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New[int]()

		assertions.Equal(99, s.FirstOr(99))
	})
}

func TestSliceFirstSafe(t *testing.T) {
	t.Run("returns first element and true for non-empty slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(10, 20, 30)

		val, ok := s.FirstSafe()
		assertions.True(ok)
		assertions.Equal(10, val)
	})

	t.Run("returns zero and false for empty slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New[int]()

		val, ok := s.FirstSafe()
		assertions.False(ok)
		assertions.Equal(0, val)
	})
}

func TestSliceLast(t *testing.T) {
	t.Run("returns last element", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(10, 20, 30)

		assertions.Equal(30, s.Last())
	})

	t.Run("panics on empty slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New[int]()

		assertions.Panics(func() {
			_ = s.Last()
		})
	})
}

func TestSliceLastOr(t *testing.T) {
	t.Run("returns last element for non-empty slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(10, 20, 30)

		assertions.Equal(30, s.LastOr(99))
	})

	t.Run("returns default for empty slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New[int]()

		assertions.Equal(99, s.LastOr(99))
	})
}

func TestSliceLastSafe(t *testing.T) {
	t.Run("returns last element and true for non-empty slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(10, 20, 30)

		val, ok := s.LastSafe()
		assertions.True(ok)
		assertions.Equal(30, val)
	})

	t.Run("returns zero and false for empty slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New[int]()

		val, ok := s.LastSafe()
		assertions.False(ok)
		assertions.Equal(0, val)
	})
}

func TestSliceAppend(t *testing.T) {
	t.Run("appends single element", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3)

		s.Append(4)

		assertions.Equal(4, s.Len())
		assertions.Equal(4, s.Last())
	})

	t.Run("appends multiple elements", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3)

		s.Append(4, 5, 6)

		assertions.Equal(6, s.Len())
	})

	t.Run("appends to empty slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New[int]()

		s.Append(1, 2, 3)

		assertions.Equal(3, s.Len())
	})
}

func TestSlicePrepend(t *testing.T) {
	t.Run("prepends single element", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(2, 3, 4)

		s.Prepend(1)

		assertions.Equal(4, s.Len())
		assertions.Equal(1, s.First())
	})

	t.Run("prepends multiple elements", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(4, 5, 6)

		s.Prepend(1, 2, 3)

		assertions.Equal(6, s.Len())
		assertions.Equal(1, s.First())
	})

	t.Run("prepends to empty slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New[int]()

		s.Prepend(1, 2, 3)

		assertions.Equal(3, s.Len())
	})
}

func TestSliceConcat(t *testing.T) {
	t.Run("concatenates single slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2)
		other := New(3, 4)

		s.Concat(other)

		assertions.Equal(4, s.Len())
	})

	t.Run("concatenates multiple slices", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2)
		other1 := New(3, 4)
		other2 := New(5, 6)

		s.Concat(other1, other2)

		assertions.Equal(6, s.Len())
	})

	t.Run("concatenates empty slices", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2)
		other := New[int]()

		s.Concat(other)

		assertions.Equal(2, s.Len())
	})
}

func TestSliceInsert(t *testing.T) {
	t.Run("inserts at beginning", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(2, 3, 4)

		s.Insert(0, 1)

		assertions.Equal(4, s.Len())
		assertions.Equal(1, s.First())
	})

	t.Run("inserts in middle", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 4, 5)

		s.Insert(2, 3)

		assertions.Equal(5, s.Len())
		assertions.Equal(3, s.Get(2))
	})

	t.Run("inserts at end", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3)

		s.Insert(3, 4)

		assertions.Equal(4, s.Len())
		assertions.Equal(4, s.Last())
	})

	t.Run("inserts multiple elements", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 5)

		s.Insert(1, 2, 3, 4)

		assertions.Equal(5, s.Len())
	})

	t.Run("handles negative index", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(2, 3)

		s.Insert(-5, 1)

		assertions.Equal(1, s.First())
	})

	t.Run("handles index beyond length", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2)

		s.Insert(100, 3)

		assertions.Equal(3, s.Last())
	})
}

func TestSliceRemove(t *testing.T) {
	t.Run("removes element at index", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3, 4, 5)

		s.Remove(2)

		assertions.Equal(4, s.Len())
		assertions.Equal(4, s.Get(2))
	})

	t.Run("removes first element", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3)

		s.Remove(0)

		assertions.Equal(2, s.Len())
		assertions.Equal(2, s.First())
	})

	t.Run("removes last element", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3)

		s.Remove(2)

		assertions.Equal(2, s.Len())
		assertions.Equal(2, s.Last())
	})

	t.Run("does nothing for negative index", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3)

		s.Remove(-1)

		assertions.Equal(3, s.Len())
	})

	t.Run("does nothing for out of bounds index", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3)

		s.Remove(10)

		assertions.Equal(3, s.Len())
	})
}

func TestSliceSet(t *testing.T) {
	t.Run("sets element at valid index", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3)

		s.Set(1, 20)

		assertions.Equal(20, s.Get(1))
	})

	t.Run("does nothing for negative index", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3)

		s.Set(-1, 99)

		assertions.Equal(1, s.Get(0))
		assertions.Equal(2, s.Get(1))
		assertions.Equal(3, s.Get(2))
	})

	t.Run("does nothing for out of bounds index", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3)

		s.Set(10, 99)

		assertions.Equal(3, s.Len())
	})
}

func TestSliceClone(t *testing.T) {
	t.Run("creates copy of slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3)

		cloned := s.Clone()

		assertions.Equal(3, cloned.Len())
	})

	t.Run("creates independent copy", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3)

		cloned := s.Clone()
		cloned.Set(0, 100)

		assertions.Equal(1, s.Get(0))
		assertions.Equal(100, cloned.Get(0))
	})

	t.Run("handles nil slice", func(t *testing.T) {
		assertions := assert.New(t)
		var s Slice[int]

		cloned := s.Clone()

		assertions.Nil(cloned)
	})

	t.Run("handles empty slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New[int]()

		cloned := s.Clone()

		assertions.Equal(0, cloned.Len())
	})
}

func TestSliceReverse(t *testing.T) {
	t.Run("reverses slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3, 4, 5)

		s.Reverse()

		expected := []int{5, 4, 3, 2, 1}
		for i, v := range expected {
			assertions.Equal(v, s.Get(i))
		}
	})

	t.Run("reverses single element", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1)

		s.Reverse()

		assertions.Equal(1, s.Get(0))
	})

	t.Run("reverses two elements", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2)

		s.Reverse()

		assertions.Equal(2, s.Get(0))
		assertions.Equal(1, s.Get(1))
	})

	t.Run("handles empty slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New[int]()

		s.Reverse() // should not panic

		assertions.Equal(0, s.Len())
	})
}

func TestSliceFind(t *testing.T) {
	t.Run("finds element matching predicate", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3, 4, 5)

		val, found := s.Find(func(v int) bool { return v > 3 })

		assertions.True(found)
		assertions.Equal(4, val)
	})

	t.Run("returns false when not found", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3)

		val, found := s.Find(func(v int) bool { return v > 10 })

		assertions.False(found)
		assertions.Equal(0, val)
	})

	t.Run("returns first matching element", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3, 4, 5)

		val, found := s.Find(func(v int) bool { return v%2 == 0 })

		assertions.True(found)
		assertions.Equal(2, val)
	})
}

func TestSliceCount(t *testing.T) {
	t.Run("counts elements matching predicate", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3, 4, 5, 6)

		count := s.Count(func(v int) bool { return v%2 == 0 })

		assertions.Equal(3, count)
	})

	t.Run("returns 0 when none match", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 3, 5, 7)

		count := s.Count(func(v int) bool { return v%2 == 0 })

		assertions.Equal(0, count)
	})

	t.Run("counts all elements when all match", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(2, 4, 6, 8)

		count := s.Count(func(v int) bool { return v%2 == 0 })

		assertions.Equal(4, count)
	})
}

func TestJoin(t *testing.T) {
	t.Run("joins with separator", func(t *testing.T) {
		assertions := assert.New(t)
		s := New("a", "b", "c")

		result := Join(s, ", ")

		assertions.Equal("a, b, c", result)
	})

	t.Run("joins with empty separator", func(t *testing.T) {
		assertions := assert.New(t)
		s := New("a", "b", "c")

		result := Join(s, "")

		assertions.Equal("abc", result)
	})

	t.Run("returns empty string for empty slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New[string]()

		result := Join(s, ", ")

		assertions.Equal("", result)
	})

	t.Run("returns element for single element slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := New("only")

		result := Join(s, ", ")

		assertions.Equal("only", result)
	})
}

func TestSliceReslicing(t *testing.T) {
	t.Run("supports standard reslicing", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3, 4, 5)

		sub := s[1:4]

		assertions.Equal(3, len(sub))
		assertions.Equal(2, sub[0])
		assertions.Equal(3, sub[1])
		assertions.Equal(4, sub[2])
	})

	t.Run("supports range iteration", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3)
		sum := 0

		for _, v := range s {
			sum += v
		}

		assertions.Equal(6, sum)
	})
}

func TestSliceOperationsSequence(t *testing.T) {
	t.Run("complex operations sequence", func(t *testing.T) {
		assertions := assert.New(t)
		s := New(1, 2, 3)

		// Append
		s.Append(4, 5)
		assertions.Equal(5, s.Len())

		// Prepend
		s.Prepend(0)
		assertions.Equal(0, s.First())

		// Insert
		s.Insert(3, 100)
		assertions.Equal(100, s.Get(3))

		// Remove
		s.Remove(3)
		assertions.Equal(3, s.Get(3))

		// Set
		s.Set(0, 99)
		assertions.Equal(99, s.First())

		// Reverse
		originalLast := s.Last()
		s.Reverse()
		assertions.Equal(originalLast, s.First())
	})
}

func TestSliceEdgeCases(t *testing.T) {
	t.Run("handles large slice", func(t *testing.T) {
		assertions := assert.New(t)
		s := WithCapacity[int](10000)

		for i := 0; i < 10000; i++ {
			s.Append(i)
		}

		assertions.Equal(10000, s.Len())
		assertions.Equal(9999, s.Last())
	})

	t.Run("handles struct elements", func(t *testing.T) {
		assertions := assert.New(t)
		type item struct {
			ID   int
			Name string
		}

		s := New(
			item{1, "first"},
			item{2, "second"},
			item{3, "third"},
		)

		assertions.Equal(3, s.Len())

		found, ok := s.Find(func(i item) bool { return i.ID == 2 })
		assertions.True(ok)
		assertions.Equal("second", found.Name)
	})

	t.Run("handles pointer elements", func(t *testing.T) {
		assertions := assert.New(t)
		a, b, c := 1, 2, 3
		s := New(&a, &b, &c)

		assertions.Equal(3, s.Len())
		assertions.Equal(1, *s.Get(0))
	})
}
