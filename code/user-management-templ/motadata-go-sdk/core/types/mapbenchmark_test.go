package types

import (
	"strconv"
	"testing"
)

// Benchmark tests

func BenchmarkMapSetOperation(b *testing.B) {
	m := NewMap()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i % 1000)
		m.Set(key, i)
	}
}

func BenchmarkMapGetOperation(b *testing.B) {
	m := NewMap()

	// Pre-populate the map
	for i := 0; i < 1000; i++ {
		key := strconv.Itoa(i)
		m.Set(key, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i % 1000)
		_, _ = m.Get(key)
	}
}

func BenchmarkMapSetGetOperation(b *testing.B) {
	m := NewMap()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i % 1000)
		m.Set(key, i)
		_, _ = m.Get(key)
	}
}

func BenchmarkMapHasOperation(b *testing.B) {
	m := NewMap()

	// Pre-populate the map
	for i := 0; i < 1000; i++ {
		key := strconv.Itoa(i)
		m.Set(key, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i % 1000)
		_ = m.Has(key)
	}
}

func BenchmarkMapDeleteOperation(b *testing.B) {
	m := NewMap()

	// Pre-populate the map with enough entries for all iterations
	for j := 0; j < b.N; j++ {
		key := strconv.Itoa(j)
		m.Set(key, j)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i)
		m.Delete(key)
	}
}

func BenchmarkMapDeleteExistingKey(b *testing.B) {
	m := NewMap()

	// Pre-populate the map once
	for j := 0; j < 1000; j++ {
		key := strconv.Itoa(j)
		m.Set(key, j)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Delete and re-add to keep the map populated
		key := strconv.Itoa(i % 1000)
		m.Delete(key)
		m.Set(key, i)
	}
}

func BenchmarkMapDeleteNonExistingKey(b *testing.B) {
	m := NewMap()

	// Pre-populate the map once
	for j := 0; j < 1000; j++ {
		key := strconv.Itoa(j)
		m.Set(key, j)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Delete keys that don't exist (keys are 0-999, we use 1000+)
		key := strconv.Itoa(1000 + i)
		m.Delete(key)
	}
}

func BenchmarkMapKeysOperation(b *testing.B) {
	m := NewMap()

	// Pre-populate the map
	for i := 0; i < 1000; i++ {
		key := strconv.Itoa(i)
		m.Set(key, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = m.Keys()
	}
}

func BenchmarkMapValuesOperation(b *testing.B) {
	m := NewMap()

	// Pre-populate the map
	for i := 0; i < 1000; i++ {
		key := strconv.Itoa(i)
		m.Set(key, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = m.Values()
	}
}

func BenchmarkMapRangeOperation(b *testing.B) {
	m := NewMap()

	// Pre-populate the map
	for i := 0; i < 1000; i++ {
		key := strconv.Itoa(i)
		m.Set(key, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Range(func(key string, value any) bool {
			return true
		})
	}
}

// Benchmark with different map sizes

func BenchmarkMapSetSmallMap(b *testing.B) {
	m := NewMap()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i % 10)
		m.Set(key, i)
	}
}

func BenchmarkMapSetMediumMap(b *testing.B) {
	m := NewMap()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i % 100)
		m.Set(key, i)
	}
}

func BenchmarkMapSetLargeMap(b *testing.B) {
	m := NewMap()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i % 10000)
		m.Set(key, i)
	}
}

func BenchmarkMapGetSmallMap(b *testing.B) {
	m := NewMap()

	for i := 0; i < 10; i++ {
		key := strconv.Itoa(i)
		m.Set(key, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i % 10)
		_, _ = m.Get(key)
	}
}

func BenchmarkMapGetMediumMap(b *testing.B) {
	m := NewMap()

	for i := 0; i < 100; i++ {
		key := strconv.Itoa(i)
		m.Set(key, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i % 100)
		_, _ = m.Get(key)
	}
}

func BenchmarkMapGetLargeMap(b *testing.B) {
	m := NewMap()

	for i := 0; i < 10000; i++ {
		key := strconv.Itoa(i)
		m.Set(key, i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i % 10000)
		_, _ = m.Get(key)
	}
}
