package types

import (
	"strconv"
	"testing"
)

// Benchmark tests for SwissMap

func BenchmarkSwissMapSetOperation(b *testing.B) {
	m := NewSwissMap[string, any]()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i % 1000)
		m.Set(key, i)
	}
}

func BenchmarkSwissMapGetOperation(b *testing.B) {
	m := NewSwissMap[string, any]()

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

func BenchmarkSwissMapSetGetOperation(b *testing.B) {
	m := NewSwissMap[string, any]()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i % 1000)
		m.Set(key, i)
		_, _ = m.Get(key)
	}
}

func BenchmarkSwissMapHasOperation(b *testing.B) {
	m := NewSwissMap[string, any]()

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

func BenchmarkSwissMapDeleteOperation(b *testing.B) {
	m := NewSwissMap[string, any]()

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

func BenchmarkSwissMapDeleteExistingKey(b *testing.B) {
	m := NewSwissMap[string, any]()

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

func BenchmarkSwissMapDeleteNonExistingKey(b *testing.B) {
	m := NewSwissMap[string, any]()

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

func BenchmarkSwissMapKeysOperation(b *testing.B) {
	m := NewSwissMap[string, any]()

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

func BenchmarkSwissMapValuesOperation(b *testing.B) {
	m := NewSwissMap[string, any]()

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

func BenchmarkSwissMapRangeOperation(b *testing.B) {
	m := NewSwissMap[string, any]()

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

func BenchmarkSwissMapSetSmallMap(b *testing.B) {
	m := NewSwissMap[string, any]()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i % 10)
		m.Set(key, i)
	}
}

func BenchmarkSwissMapSetMediumMap(b *testing.B) {
	m := NewSwissMap[string, any]()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i % 100)
		m.Set(key, i)
	}
}

func BenchmarkSwissMapSetLargeMap(b *testing.B) {
	m := NewSwissMapWithCapacity[string, any](10000)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i % 10000)
		m.Set(key, i)
	}
}

func BenchmarkSwissMapGetSmallMap(b *testing.B) {
	m := NewSwissMap[string, any]()

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

func BenchmarkSwissMapGetMediumMap(b *testing.B) {
	m := NewSwissMap[string, any]()

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

func BenchmarkSwissMapGetLargeMap(b *testing.B) {
	m := NewSwissMapWithCapacity[string, any](10000)

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

// Comparison benchmarks: SwissMap vs Map vs IntMap

func BenchmarkComparisonSwissMapGet(b *testing.B) {
	m := NewSwissMap[string, any]()

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

func BenchmarkComparisonMapGet(b *testing.B) {
	m := NewMap[string, any]()

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

func BenchmarkComparisonSwissMapSet(b *testing.B) {
	m := NewSwissMap[string, any]()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i % 1000)
		m.Set(key, i)
	}
}

func BenchmarkComparisonMapSet(b *testing.B) {
	m := NewMap[string, any]()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := strconv.Itoa(i % 1000)
		m.Set(key, i)
	}
}

func BenchmarkComparisonSwissMapHas(b *testing.B) {
	m := NewSwissMap[string, any]()

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

func BenchmarkComparisonMapHas(b *testing.B) {
	m := NewMap[string, any]()

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
