package l1cache

import (
	"fmt"
	"testing"
	"time"
)

// BenchmarkCacheSet benchmarks cache set operations with []byte
func BenchmarkCacheSet(b *testing.B) {
	cache := NewCacheManager(&Config{
		TotalSize: 100 * 1024 * 1024,
	})
	defer func() {
		_ = cache.Close()
	}()

	tenantID := "tenant-001"
	orgID := "org-001"
	data := []byte("benchmark data string")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key:%d", i)
		_ = cache.SetTTL(tenantID, orgID, key, generateValue(data), 1*time.Minute)
	}
}

// BenchmarkCacheGet benchmarks cache get operations with []byte
func BenchmarkCacheGet(b *testing.B) {
	cache := NewCacheManager(&Config{
		TotalSize: 100 * 1024 * 1024,
	})
	defer func() {
		_ = cache.Close()
	}()

	tenantID := "tenant-001"
	orgID := "org-001"
	data := []byte("benchmark data string")

	// Pre-populate cache
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key:%d", i)
		_ = cache.SetTTL(tenantID, orgID, key, generateValue(data), 10*time.Minute)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key:%d", i%10000)
		_, _ = cache.Get(tenantID, orgID, key)
	}
}

// BenchmarkConcurrentAccess benchmarks concurrent cache access with []byte
func BenchmarkConcurrentAccess(b *testing.B) {
	cache := NewCacheManager(&Config{
		TotalSize: 100 * 1024 * 1024,
	})
	defer func() {
		_ = cache.Close()
	}()

	tenantID := "tenant-001"
	orgID := "org-001"

	// Pre-populate cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key:%d", i)
		value := []byte(fmt.Sprintf("value_%d", i))
		_ = cache.SetTTL(tenantID, orgID, key, generateValue(value), 10*time.Minute)
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("key:%d", i%1000)
			if i%2 == 0 {
				_, _ = cache.Get(tenantID, orgID, key)
			} else {
				value := []byte(fmt.Sprintf("new_value_%d", i))
				_ = cache.SetTTL(tenantID, orgID, key, generateValue(value), 10*time.Minute)
			}
			i++
		}
	})
}

// BenchmarkSmallData benchmarks cache with small byte slices
func BenchmarkSmallData(b *testing.B) {
	cache := NewCacheManager(&Config{
		TotalSize: 100 * 1024 * 1024,
	})
	defer func() {
		_ = cache.Close()
	}()

	tenantID := "tenant-001"
	orgID := "org-001"
	data := []byte("small")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key:%d", i)
		_ = cache.SetTTL(tenantID, orgID, key, generateValue(data), 1*time.Minute)
		_, _ = cache.Get(tenantID, orgID, key)
	}
}

// BenchmarkLargeData benchmarks cache with large byte slices
func BenchmarkLargeData(b *testing.B) {
	cache := NewCacheManager(&Config{
		TotalSize: 500 * 1024 * 1024,
	})
	defer func() {
		_ = cache.Close()
	}()

	tenantID := "tenant-001"
	orgID := "org-001"

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
		_ = cache.SetTTL(tenantID, orgID, key, generateValue(data), 1*time.Minute)
		_, _ = cache.Get(tenantID, orgID, key)
	}
}

// BenchmarkJSONLikeData benchmarks cache with JSON-like byte data
func BenchmarkJSONLikeData(b *testing.B) {
	cache := NewCacheManager(&Config{
		TotalSize: 100 * 1024 * 1024,
	})
	defer func() {
		_ = cache.Close()
	}()

	tenantID := "tenant-001"
	orgID := "org-001"
	jsonData := []byte(`{"id":12345,"name":"John Doe","email":"john@example.com","active":true,"tags":["user","premium"],"metadata":{"created":"2024-01-01","updated":"2024-12-01"}}`)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("user:%d", i)
		_ = cache.SetTTL(tenantID, orgID, key, generateValue(jsonData), 1*time.Minute)
		_, _ = cache.Get(tenantID, orgID, key)
	}
}

// BenchmarkBinaryData benchmarks cache with binary data
func BenchmarkBinaryData(b *testing.B) {
	cache := NewCacheManager(&Config{
		TotalSize: 100 * 1024 * 1024,
	})
	defer func() {
		_ = cache.Close()
	}()

	tenantID := "tenant-001"
	orgID := "org-001"

	// Create binary data pattern
	binaryData := make([]byte, 256)
	for i := range binaryData {
		binaryData[i] = byte(i)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("binary:%d", i)
		_ = cache.SetTTL(tenantID, orgID, key, generateValue(binaryData), 1*time.Minute)
		_, _ = cache.Get(tenantID, orgID, key)
	}
}

// BenchmarkSetOnly benchmarks only Set operations
func BenchmarkSetOnly(b *testing.B) {
	cache := NewCacheManager(&Config{
		TotalSize: 100 * 1024 * 1024,
	})
	defer func() {
		_ = cache.Close()
	}()

	tenantID := "tenant-001"
	orgID := "org-001"
	data := []byte("set only benchmark")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key:%d", i)
		_ = cache.SetTTL(tenantID, orgID, key, generateValue(data), 1*time.Minute)
	}
}

// BenchmarkGetOnly benchmarks only Get operations
func BenchmarkGetOnly(b *testing.B) {
	cache := NewCacheManager(&Config{
		TotalSize: 100 * 1024 * 1024,
	})
	defer func() {
		_ = cache.Close()
	}()

	tenantID := "tenant-001"
	orgID := "org-001"
	data := []byte("get only benchmark")

	// Pre-populate cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key:%d", i)
		_ = cache.SetTTL(tenantID, orgID, key, generateValue(data), 10*time.Minute)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key:%d", i%1000)
		_, _ = cache.Get(tenantID, orgID, key)
	}
}

// BenchmarkMultiTenantSet benchmarks set operations across multiple tenants
func BenchmarkMultiTenantSet(b *testing.B) {
	cache := NewCacheManager(&Config{
		TotalSize: 100 * 1024 * 1024,
	})
	defer func() {
		_ = cache.Close()
	}()

	data := []byte("multi-tenant benchmark data")

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenantID := fmt.Sprintf("tenant-%03d", i%100)
		orgID := fmt.Sprintf("org-%03d", i%100)
		key := fmt.Sprintf("key:%d", i)
		_ = cache.SetTTL(tenantID, orgID, key, generateValue(data), 1*time.Minute)
	}
}

// BenchmarkMultiTenantGet benchmarks get operations across multiple tenants
func BenchmarkMultiTenantGet(b *testing.B) {
	cache := NewCacheManager(&Config{
		TotalSize: 100 * 1024 * 1024,
	})
	defer func() {
		_ = cache.Close()
	}()

	data := []byte("multi-tenant benchmark data")

	// Pre-populate cache for multiple tenants
	for i := 0; i < 10000; i++ {
		tenantID := fmt.Sprintf("tenant-%03d", i%100)
		orgID := fmt.Sprintf("org-%03d", i%100)
		key := fmt.Sprintf("key:%d", i)
		_ = cache.SetTTL(tenantID, orgID, key, generateValue(data), 10*time.Minute)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		tenantID := fmt.Sprintf("tenant-%03d", i%100)
		orgID := fmt.Sprintf("org-%03d", i%100)
		key := fmt.Sprintf("key:%d", i%10000)
		_, _ = cache.Get(tenantID, orgID, key)
	}
}

// BenchmarkMultiTenantConcurrent benchmarks concurrent access across multiple tenants
func BenchmarkMultiTenantConcurrent(b *testing.B) {
	cache := NewCacheManager(&Config{
		TotalSize: 100 * 1024 * 1024,
	})
	defer func() {
		_ = cache.Close()
	}()

	// Pre-populate cache for multiple tenants
	for i := 0; i < 1000; i++ {
		tenantID := fmt.Sprintf("tenant-%03d", i%10)
		orgID := fmt.Sprintf("org-%03d", i%10)
		key := fmt.Sprintf("key:%d", i)
		value := []byte(fmt.Sprintf("value_%d", i))
		_ = cache.SetTTL(tenantID, orgID, key, generateValue(value), 10*time.Minute)
	}

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			tenantID := fmt.Sprintf("tenant-%03d", i%10)
			orgID := fmt.Sprintf("org-%03d", i%10)
			key := fmt.Sprintf("key:%d", i%1000)
			if i%2 == 0 {
				_, _ = cache.Get(tenantID, orgID, key)
			} else {
				value := []byte(fmt.Sprintf("new_value_%d", i))
				_ = cache.SetTTL(tenantID, orgID, key, generateValue(value), 10*time.Minute)
			}
			i++
		}
	})
}

// BenchmarkDeleteOperation benchmarks delete operations
func BenchmarkDeleteOperation(b *testing.B) {
	cache := NewCacheManager(&Config{
		TotalSize: 100 * 1024 * 1024,
	})
	defer func() {
		_ = cache.Close()
	}()

	tenantID := "tenant-001"
	orgID := "org-001"
	data := []byte("delete benchmark data")

	// Pre-populate cache
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key:%d", i)
		_ = cache.SetTTL(tenantID, orgID, key, generateValue(data), 10*time.Minute)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key:%d", i)
		cache.Delete(tenantID, orgID, key)
	}
}

// BenchmarkExistsOperation benchmarks exists check operations
func BenchmarkExistsOperation(b *testing.B) {
	cache := NewCacheManager(&Config{
		TotalSize: 100 * 1024 * 1024,
	})
	defer func() {
		_ = cache.Close()
	}()

	tenantID := "tenant-001"
	orgID := "org-001"
	data := []byte("exists benchmark data")

	// Pre-populate cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key:%d", i)
		_ = cache.SetTTL(tenantID, orgID, key, generateValue(data), 10*time.Minute)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key:%d", i%1000)
		cache.Exists(tenantID, orgID, key)
	}
}

// BenchmarkGetWithBuf benchmarks GetWithBuf with pre-allocated buffer (zero-allocation)
func BenchmarkGetWithBuf(b *testing.B) {
	cache := NewCacheManager(&Config{
		TotalSize: 100 * 1024 * 1024,
	})
	defer func() {
		_ = cache.Close()
	}()

	tenantID := "tenant-001"
	orgID := "org-001"
	data := []byte("benchmark data for getwithbuf")

	// Pre-populate cache
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key:%d", i)
		_ = cache.SetTTL(tenantID, orgID, key, generateValue(data), 10*time.Minute)
	}

	// Pre-allocate buffer for zero-allocation reads
	buf := make([]byte, 1024)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key:%d", i%1000)
		_, _ = cache.GetWithBuf(tenantID, orgID, key, buf)
	}
}
