package types

import (
	"testing"
)

// Benchmark tests for IntMap

func BenchmarkIntMapSetOperation(b *testing.B) {
	m := NewIntMap()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Set(i%1000, i)
	}
}

func BenchmarkIntMapGetOperation(b *testing.B) {
	m := NewIntMap()

	// Pre-populate the map
	for i := 0; i < 1000; i++ {
		m.Set(i, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = m.Get(i % 1000)
	}
}

func BenchmarkIntMapSetGetOperation(b *testing.B) {
	m := NewIntMap()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := i % 1000
		m.Set(key, i)
		_, _ = m.Get(key)
	}
}

func BenchmarkIntMapHasOperation(b *testing.B) {
	m := NewIntMap()

	// Pre-populate the map
	for i := 0; i < 1000; i++ {
		m.Set(i, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = m.Has(i % 1000)
	}
}

func BenchmarkIntMapDeleteOperation(b *testing.B) {
	m := NewIntMap()

	// Pre-populate the map with enough entries for all iterations
	for j := 0; j < b.N; j++ {
		m.Set(j, j)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Delete(i)
	}
}

func BenchmarkIntMapDeleteExistingKey(b *testing.B) {
	m := NewIntMap()

	// Pre-populate the map once
	for j := 0; j < 1000; j++ {
		m.Set(j, j)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Delete and re-add to keep the map populated
		key := i % 1000
		m.Delete(key)
		m.Set(key, i)
	}
}

func BenchmarkIntMapDeleteNonExistingKey(b *testing.B) {
	m := NewIntMap()

	// Pre-populate the map once
	for j := 0; j < 1000; j++ {
		m.Set(j, j)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Delete keys that don't exist (keys are 0-999, we use 1000+)
		m.Delete(1000 + i)
	}
}

func BenchmarkIntMapKeysOperation(b *testing.B) {
	m := NewIntMap()

	// Pre-populate the map
	for i := 0; i < 1000; i++ {
		m.Set(i, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = m.Keys()
	}
}

func BenchmarkIntMapValuesOperation(b *testing.B) {
	m := NewIntMap()

	// Pre-populate the map
	for i := 0; i < 1000; i++ {
		m.Set(i, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = m.Values()
	}
}

func BenchmarkIntMapRangeOperation(b *testing.B) {
	m := NewIntMap()

	// Pre-populate the map
	for i := 0; i < 1000; i++ {
		m.Set(i, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Range(func(key int, value any) bool {
			return true
		})
	}
}

// Benchmark with different map sizes

func BenchmarkIntMapSetSmallMap(b *testing.B) {
	m := NewIntMap()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Set(i%10, i)
	}
}

func BenchmarkIntMapSetMediumMap(b *testing.B) {
	m := NewIntMap()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Set(i%100, i)
	}
}

func BenchmarkIntMapSetLargeMap(b *testing.B) {
	m := NewIntMapWithCapacity(10000)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Set(i%10000, i)
	}
}

func BenchmarkIntMapGetSmallMap(b *testing.B) {
	m := NewIntMap()

	for i := 0; i < 10; i++ {
		m.Set(i, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = m.Get(i % 10)
	}
}

func BenchmarkIntMapGetMediumMap(b *testing.B) {
	m := NewIntMap()

	for i := 0; i < 100; i++ {
		m.Set(i, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = m.Get(i % 100)
	}
}

func BenchmarkIntMapGetLargeMap(b *testing.B) {
	m := NewIntMapWithCapacity(10000)

	for i := 0; i < 10000; i++ {
		m.Set(i, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = m.Get(i % 10000)
	}
}

// Comparison benchmarks: IntMap vs standard Map

func BenchmarkComparisonIntMapGet(b *testing.B) {
	m := NewIntMap()

	for i := 0; i < 1000; i++ {
		m.Set(i, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = m.Get(i % 1000)
	}
}

func BenchmarkComparisonStandardMapGet(b *testing.B) {
	m := NewMap()

	for i := 0; i < 1000; i++ {
		m.Set(string(rune(i)), i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = m.Get(string(rune(i % 1000)))
	}
}

func BenchmarkComparisonIntMapSet(b *testing.B) {
	m := NewIntMap()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Set(i%1000, i)
	}
}

func BenchmarkComparisonStandardMapSet(b *testing.B) {
	m := NewMap()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Set(string(rune(i%1000)), i)
	}
}
