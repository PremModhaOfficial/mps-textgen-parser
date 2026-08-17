package types

import (
	"testing"
)

// Benchmark tests for Slice

func BenchmarkSliceAppendOperation(b *testing.B) {
	s := New[int]()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Append(i)
	}
}

func BenchmarkSlicePrependOperation(b *testing.B) {
	s := New[int]()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Prepend(i)
	}
}

func BenchmarkSliceGetOperation(b *testing.B) {
	s := WithLength[int](1000)
	for i := 0; i < 1000; i++ {
		s.Set(i, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.Get(i % 1000)
	}
}

func BenchmarkSliceGetOrOperation(b *testing.B) {
	s := WithLength[int](1000)
	for i := 0; i < 1000; i++ {
		s.Set(i, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.GetOr(i%1000, -1)
	}
}

func BenchmarkSliceGetSafeOperation(b *testing.B) {
	s := WithLength[int](1000)
	for i := 0; i < 1000; i++ {
		s.Set(i, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = s.GetSafe(i % 1000)
	}
}

func BenchmarkSliceSetOperation(b *testing.B) {
	s := WithLength[int](1000)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Set(i%1000, i)
	}
}

func BenchmarkSliceInsertOperation(b *testing.B) {
	s := New[int]()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Insert(i%100, i)
	}
}

func BenchmarkSliceRemoveOperation(b *testing.B) {
	s := New(1, 2, 3, 4, 5)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if s.Len() > 1000 {
			s.Remove(i % s.Len())
		} else {
			s = New(1, 2, 3, 4, 5)
		}
	}
}

func BenchmarkSliceCloneOperation(b *testing.B) {
	s := New[int]()
	for i := 0; i < 1000; i++ {
		s.Append(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.Clone()
	}
}

func BenchmarkSliceCloneSmall(b *testing.B) {
	s := New(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.Clone()
	}
}

func BenchmarkSliceCloneMedium(b *testing.B) {
	s := New[int]()
	for i := 0; i < 100; i++ {
		s.Append(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.Clone()
	}
}

func BenchmarkSliceCloneLarge(b *testing.B) {
	s := New[int]()
	for i := 0; i < 10000; i++ {
		s.Append(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.Clone()
	}
}

func BenchmarkSliceReverseOperation(b *testing.B) {
	s := New[int]()
	for i := 0; i < 1000; i++ {
		s.Append(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Reverse()
	}
}

func BenchmarkSliceReverseSmall(b *testing.B) {
	s := New(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Reverse()
	}
}

func BenchmarkSliceReverseLarge(b *testing.B) {
	s := New[int]()
	for i := 0; i < 10000; i++ {
		s.Append(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Reverse()
	}
}

func BenchmarkSliceFindOperation(b *testing.B) {
	s := New[int]()
	for i := 0; i < 1000; i++ {
		s.Append(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = s.Find(func(v int) bool { return v == 500 })
	}
}

func BenchmarkSliceFindFirst(b *testing.B) {
	s := New[int]()
	for i := 0; i < 1000; i++ {
		s.Append(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = s.Find(func(v int) bool { return v == 0 })
	}
}

func BenchmarkSliceFindLast(b *testing.B) {
	s := New[int]()
	for i := 0; i < 1000; i++ {
		s.Append(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = s.Find(func(v int) bool { return v == 999 })
	}
}

func BenchmarkSliceFindNotFound(b *testing.B) {
	s := New[int]()
	for i := 0; i < 1000; i++ {
		s.Append(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = s.Find(func(v int) bool { return v == 9999 })
	}
}

func BenchmarkSliceCountOperation(b *testing.B) {
	s := New[int]()
	for i := 0; i < 1000; i++ {
		s.Append(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.Count(func(v int) bool { return v%2 == 0 })
	}
}

func BenchmarkSliceLenOperation(b *testing.B) {
	s := New[int]()
	for i := 0; i < 1000; i++ {
		s.Append(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.Len()
	}
}

func BenchmarkSliceCapOperation(b *testing.B) {
	s := Make[int](1000, 2000)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.Cap()
	}
}

func BenchmarkSliceIsEmptyOperation(b *testing.B) {
	s := New[int]()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.IsEmpty()
	}
}

func BenchmarkSliceFirstOperation(b *testing.B) {
	s := New[int]()
	for i := 0; i < 1000; i++ {
		s.Append(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.First()
	}
}

func BenchmarkSliceLastOperation(b *testing.B) {
	s := New[int]()
	for i := 0; i < 1000; i++ {
		s.Append(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.Last()
	}
}

func BenchmarkSliceConcatOperation(b *testing.B) {
	s1 := New[int]()
	s2 := New[int]()
	for i := 0; i < 100; i++ {
		s1.Append(i)
		s2.Append(i + 100)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		s := s1.Clone()
		b.StartTimer()

		s.Concat(s2)
	}
}

func BenchmarkJoinOperation(b *testing.B) {
	s := New[string]()
	for i := 0; i < 100; i++ {
		s.Append("element")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Join(s, ", ")
	}
}

func BenchmarkJoinSmall(b *testing.B) {
	s := New("a", "b", "c", "d", "e")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Join(s, ", ")
	}
}

func BenchmarkJoinLarge(b *testing.B) {
	s := New[string]()
	for i := 0; i < 1000; i++ {
		s.Append("element")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = Join(s, ", ")
	}
}

// Parallel benchmarks

func BenchmarkSliceParallelGet(b *testing.B) {
	s := WithLength[int](1000)
	for i := 0; i < 1000; i++ {
		s.Set(i, i)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			_ = s.Get(counter % 1000)
			counter++
		}
	})
}

func BenchmarkSliceParallelLen(b *testing.B) {
	s := WithLength[int](1000)
	for i := 0; i < 1000; i++ {
		s.Set(i, i)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = s.Len()
		}
	})
}

func BenchmarkSliceParallelClone(b *testing.B) {
	s := New[int]()
	for i := 0; i < 100; i++ {
		s.Append(i)
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = s.Clone()
		}
	})
}

// Size comparison benchmarks

func BenchmarkSliceAppendSmall(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s := New[int]()
		for j := 0; j < 10; j++ {
			s.Append(j)
		}
	}
}

func BenchmarkSliceAppendMedium(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s := New[int]()
		for j := 0; j < 100; j++ {
			s.Append(j)
		}
	}
}

func BenchmarkSliceAppendLarge(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s := New[int]()
		for j := 0; j < 1000; j++ {
			s.Append(j)
		}
	}
}

func BenchmarkSliceAppendWithCapacity(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s := WithCapacity[int](1000)
		for j := 0; j < 1000; j++ {
			s.Append(j)
		}
	}
}

// Comparison with standard slice

func BenchmarkComparisonSliceAppend(b *testing.B) {
	s := New[int]()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.Append(i)
	}
}

func BenchmarkComparisonStandardSliceAppend(b *testing.B) {
	var s []int

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s = append(s, i)
	}
}

func BenchmarkComparisonSliceGet(b *testing.B) {
	s := WithLength[int](1000)
	for i := 0; i < 1000; i++ {
		s.Set(i, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.Get(i % 1000)
	}
}

func BenchmarkComparisonStandardSliceGet(b *testing.B) {
	s := make([]int, 1000)
	for i := 0; i < 1000; i++ {
		s[i] = i
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s[i%1000]
	}
}
