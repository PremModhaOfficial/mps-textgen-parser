package l1cache

import (
	"fmt"
	"testing"
	"time"
)

// BenchmarkCacheSet benchmarks cache set operations with []byte
func BenchmarkCacheSet(b *testing.B) {
	cache := NewCacheManager(&Config{
		CacheSize: 100 * 1024 * 1024,
	})
	defer func(cache Cache) {
		_ = cache.Close()
	}(cache)

	data := []byte("benchmark data string")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key:%d", i)
		_ = cache.SetTTL(key, data, 1*time.Minute)
	}
}

// BenchmarkCacheGet benchmarks cache get operations with []byte
func BenchmarkCacheGet(b *testing.B) {
	cache := NewCacheManager(&Config{
		CacheSize: 100 * 1024 * 1024,
	})

	defer func(cache Cache) {
		_ = cache.Close()
	}(cache)

	data := []byte("benchmark data string")

	// Pre-populate cache
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key:%d", i)
		_ = cache.SetTTL(key, data, 10*time.Minute)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key:%d", i%10000)
		cache.Get(key)
	}
}

// BenchmarkConcurrentAccess benchmarks concurrent cache access with []byte
func BenchmarkConcurrentAccess(b *testing.B) {
	cache := NewCacheManager(&Config{
		CacheSize: 100 * 1024 * 1024,
	})

	defer func(cache Cache) {
		_ = cache.Close()
	}(cache)

	// Pre-populate cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key:%d", i)
		value := []byte(fmt.Sprintf("value_%d", i))
		_ = cache.SetTTL(key, value, 10*time.Minute)
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key:%d", i%1000)
			if i%2 == 0 {
				cache.Get(key)
			} else {
				value := []byte(fmt.Sprintf("new_value_%d", i))
				_ = cache.SetTTL(key, value, 10*time.Minute)
			}
			i++
		}
	})
}

// BenchmarkSmallData benchmarks cache with small byte slices
func BenchmarkSmallData(b *testing.B) {
	cache := NewCacheManager(&Config{
		CacheSize: 100 * 1024 * 1024,
	})

	defer func(cache Cache) {
		_ = cache.Close()
	}(cache)

	data := []byte("small")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key:%d", i)
		_ = cache.SetTTL(key, data, 1*time.Minute)
		cache.Get(key)
	}
}

// BenchmarkLargeData benchmarks cache with large byte slices
func BenchmarkLargeData(b *testing.B) {
	cache := NewCacheManager(&Config{
		CacheSize: 500 * 1024 * 1024,
	})

	defer func(cache Cache) {
		_ = cache.Close()
	}(cache)

	// Create 10KB data
	data := make([]byte, 10240)
	for i := range data {
		data[i] = byte(i % 256)
	}

	b.SetBytes(int64(len(data)))
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key:%d", i%100) // Reuse keys to test overwrite
		_ = cache.SetTTL(key, data, 1*time.Minute)
		cache.Get(key)
	}
}

// BenchmarkJSONLikeData benchmarks cache with JSON-like byte data
func BenchmarkJSONLikeData(b *testing.B) {
	cache := NewCacheManager(&Config{
		CacheSize: 100 * 1024 * 1024,
	})

	defer func(cache Cache) {
		_ = cache.Close()
	}(cache)

	jsonData := []byte(`{"id":12345,"name":"John Doe","email":"john@example.com","active":true,"tags":["user","premium"],"metadata":{"created":"2024-01-01","updated":"2024-12-01"}}`)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("user:%d", i)
		_ = cache.SetTTL(key, jsonData, 1*time.Minute)
		cache.Get(key)
	}
}

// BenchmarkBinaryData benchmarks cache with binary data
func BenchmarkBinaryData(b *testing.B) {
	cache := NewCacheManager(&Config{
		CacheSize: 100 * 1024 * 1024,
	})

	defer func(cache Cache) {
		_ = cache.Close()
	}(cache)

	// Create binary data pattern
	binaryData := make([]byte, 256)
	for i := range binaryData {
		binaryData[i] = byte(i)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("binary:%d", i)
		_ = cache.SetTTL(key, binaryData, 1*time.Minute)
		cache.Get(key)
	}
}

// BenchmarkSetOnly benchmarks only Set operations
func BenchmarkSetOnly(b *testing.B) {
	cache := NewCacheManager(&Config{
		CacheSize: 100 * 1024 * 1024,
	})

	defer func(cache Cache) {
		_ = cache.Close()
	}(cache)

	data := []byte("set only benchmark")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key:%d", i)
		_ = cache.SetTTL(key, data, 1*time.Minute)
	}
}

// BenchmarkGetOnly benchmarks only Get operations
func BenchmarkGetOnly(b *testing.B) {
	cache := NewCacheManager(&Config{
		CacheSize: 100 * 1024 * 1024,
	})

	defer func(cache Cache) {
		_ = cache.Close()
	}(cache)

	data := []byte("get only benchmark")

	// Pre-populate cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key:%d", i)
		_ = cache.SetTTL(key, data, 10*time.Minute)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key:%d", i%1000)
		cache.Get(key)
	}
}
