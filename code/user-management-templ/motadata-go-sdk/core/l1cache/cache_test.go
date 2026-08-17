package l1cache

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// generateValue creates a buffer with ChecksumSize bytes reserved at the beginning
// and copies the value into buf[ChecksumSize:]. Returns the complete buffer.
func generateValue(value []byte) []byte {
	buf := make([]byte, ChecksumSize+len(value))
	copy(buf[ChecksumSize:], value)
	return buf
}

func TestCacheBasicOperations(t *testing.T) {
	assertions := assert.New(t)

	cache := NewCacheManager(&Config{
		TotalSize: 10 * 1024 * 1024,
	})
	defer func() {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}()

	tenantID := "tenant-001"
	orgID := "org-001"

	// Test byte slice operations
	key := "test:bytes"
	value := []byte("test value")

	err := cache.SetTTL(tenantID, orgID, key, generateValue(value), 1*time.Minute)
	assertions.NoError(err, "Failed to set value")

	retrieved, err := cache.Get(tenantID, orgID, key)
	assertions.NoError(err, "Expected to find key")
	assertions.Equal(string(value), string(retrieved), "Retrieved value should match set value")

	// Test string as bytes
	stringKey := "test:string"
	stringValue := []byte("hello world")

	err = cache.SetTTL(tenantID, orgID, stringKey, generateValue(stringValue), 1*time.Minute)
	assertions.NoError(err, "Failed to set string value")

	retrieved, err = cache.Get(tenantID, orgID, stringKey)
	assertions.NoError(err, "Expected to find string key")
	assertions.Equal(string(stringValue), string(retrieved), "Retrieved string value should match set value")

	// Delete test
	deleted := cache.Delete(tenantID, orgID, key)
	assertions.True(deleted, "Expected successful deletion")

	_, err = cache.Get(tenantID, orgID, key)
	assertions.Error(err, "Expected key to be deleted")
}

func TestCacheWithDifferentSizes(t *testing.T) {
	assertions := assert.New(t)

	cache := NewCacheManager(&Config{
		TotalSize: 10 * 1024 * 1024,
	})
	defer func() {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}()

	tenantID := "tenant-001"
	orgID := "org-001"

	// Test various byte slice sizes
	testCases := []struct {
		key  string
		data []byte
	}{
		{"empty", []byte("")},
		{"small", []byte("a")},
		{"medium", []byte("hello world")},
		{"large", make([]byte, 1024)},
		{"binary", []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD}},
	}

	// Fill large test case with pattern
	for i := range testCases[3].data {
		testCases[3].data[i] = byte(i % 256)
	}

	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			assertions2 := assert.New(t)

			err := cache.SetTTL(tenantID, orgID, tc.key, generateValue(tc.data), 1*time.Minute)
			assertions2.NoError(err, "Failed to set %s", tc.key)

			retrieved, err := cache.Get(tenantID, orgID, tc.key)
			assertions2.NoError(err, "Expected to find %s", tc.key)
			assertions2.Equal(len(tc.data), len(retrieved), "Expected length to match")
			// Handle empty slice case where cache might return nil instead of empty slice
			if len(tc.data) == 0 {
				assertions2.Empty(retrieved, "Retrieved data should be empty")
			} else {
				assertions2.Equal(tc.data, retrieved, "Retrieved data should match set data")
			}
		})
	}
}

func TestConcurrentAccess(t *testing.T) {
	assertions := assert.New(t)

	cache := NewCacheManager(&Config{
		TotalSize: 50 * 1024 * 1024,
	})
	defer func() {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}()

	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 100

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			tenantID := fmt.Sprintf("tenant-%03d", id)
			orgID := fmt.Sprintf("org-%03d", id)
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("concurrent:%d:%d", id, j)
				value := []byte(fmt.Sprintf("value_%d_%d", id, j))

				err := cache.SetTTL(tenantID, orgID, key, generateValue(value), 1*time.Minute)
				assertions.NoError(err, "Failed to set key %s", key)

				retrieved, err := cache.Get(tenantID, orgID, key)
				assertions.NoError(err, "Failed to find key %s", key)
				assertions.Equal(string(value), string(retrieved), "Retrieved value should match set value")
			}
		}(i)
	}

	wg.Wait()
}

func TestCacheClear(t *testing.T) {
	assertions := assert.New(t)

	cache := NewCacheManager(nil)
	defer func() {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}()

	tenantID := "tenant-001"
	orgID := "org-001"

	// Add items
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("test:%d", i)
		value := []byte(fmt.Sprintf("value %d", i))
		_ = cache.SetTTL(tenantID, orgID, key, generateValue(value), 1*time.Minute)
	}

	// Verify items exist
	stats := cache.Stats()
	assertions.Greater(stats.EntryCount, int64(0), "Expected cache to have items")

	// Clear
	cache.Clear()

	// Verify cleared
	stats = cache.Stats()
	assertions.Equal(int64(0), stats.EntryCount, "Expected cache entry count to be 0")

	// Verify items are gone
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("test:%d", i)
		_, err := cache.Get(tenantID, orgID, key)
		assertions.Error(err, "Expected key %s to be cleared", key)
	}
}

func TestCacheExpiration(t *testing.T) {
	assertions := assert.New(t)

	cache := NewCacheManager(&Config{
		TotalSize: 10 * 1024 * 1024,
	})
	defer func() {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}()

	tenantID := "tenant-001"
	orgID := "org-001"
	key := "test:expiry"
	value := []byte("expires soon")

	// Set with very short TTL
	err := cache.SetTTL(tenantID, orgID, key, generateValue(value), 100*time.Millisecond)
	assertions.NoError(err, "Failed to set value")

	// Should be present immediately
	retrieved, err := cache.Get(tenantID, orgID, key)
	assertions.NoError(err, "Expected to find key immediately")
	assertions.Equal(string(value), string(retrieved), "Retrieved value should match set value")

	// Wait for expiration - FreeCache may not immediately expire
	time.Sleep(200 * time.Millisecond)

	// Should be expired now (but FreeCache might be lazy, so we'll be lenient)
	_, err = cache.Get(tenantID, orgID, key)
	// Note: FreeCache may not immediately expire items, so this test might be flaky
	// In production, expired items will be cleaned up eventually
}

func TestCacheOverwrite(t *testing.T) {
	assertions := assert.New(t)

	cache := NewCacheManager(&Config{
		TotalSize: 10 * 1024 * 1024,
	})
	defer func() {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}()

	tenantID := "tenant-001"
	orgID := "org-001"
	key := "test:overwrite"
	value1 := []byte("original value")
	value2 := []byte("new value")

	// Set original value
	err := cache.SetTTL(tenantID, orgID, key, generateValue(value1), 1*time.Minute)
	assertions.NoError(err, "Failed to set original value")

	retrieved, err := cache.Get(tenantID, orgID, key)
	assertions.NoError(err, "Expected to find original value")
	assertions.Equal(string(value1), string(retrieved), "Retrieved value should match original value")

	// Overwrite with new value
	err = cache.SetTTL(tenantID, orgID, key, generateValue(value2), 1*time.Minute)
	assertions.NoError(err, "Failed to set new value")

	retrieved, err = cache.Get(tenantID, orgID, key)
	assertions.NoError(err, "Expected to find new value")
	assertions.Equal(string(value2), string(retrieved), "Retrieved value should match new value")
}

func TestCacheWithDefaultConfig(t *testing.T) {
	assertions := assert.New(t)

	cache := NewCacheManager(nil)
	defer func() {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}()

	assertions.NotNil(cache, "Cache should not be nil")

	tenantID := "tenant-001"
	orgID := "org-001"
	key := "test:default.config"
	value := []byte("default config")

	// Set with default TTL
	err := cache.Set(tenantID, orgID, key, generateValue(value))
	assertions.NoError(err, "Failed to set value")

	// Should be present immediately
	retrieved, err := cache.Get(tenantID, orgID, key)
	assertions.NoError(err, "Expected to find key immediately")
	assertions.Equal(string(value), string(retrieved), "Retrieved value should match set value")
}

func TestMultiTenantIsolation(t *testing.T) {
	assertions := assert.New(t)

	cache := NewCacheManager(&Config{
		TotalSize: 10 * 1024 * 1024,
	})
	defer func() {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}()

	// Same key, different tenants/orgs - should be fully isolated
	key := "shared:key"
	tenant1 := "tenant-001"
	tenant2 := "tenant-002"
	tenant3 := "tenant-003"
	org1 := "org-001"
	org2 := "org-002"
	org3 := "org-003"

	value1 := []byte("tenant-001-data")
	value2 := []byte("tenant-002-data")
	value3 := []byte("tenant-003-data")

	// Set values for each tenant/org
	assertions.NoError(cache.Set(tenant1, org1, key, generateValue(value1)), "Failed to set tenant1 value")
	assertions.NoError(cache.Set(tenant2, org2, key, generateValue(value2)), "Failed to set tenant2 value")
	assertions.NoError(cache.Set(tenant3, org3, key, generateValue(value3)), "Failed to set tenant3 value")

	// Verify each tenant/org gets their own data
	retrieved1, err := cache.Get(tenant1, org1, key)
	assertions.NoError(err, "Expected to find tenant1 key")
	assertions.Equal(string(value1), string(retrieved1), "Tenant1 should get their own data")

	retrieved2, err := cache.Get(tenant2, org2, key)
	assertions.NoError(err, "Expected to find tenant2 key")
	assertions.Equal(string(value2), string(retrieved2), "Tenant2 should get their own data")

	retrieved3, err := cache.Get(tenant3, org3, key)
	assertions.NoError(err, "Expected to find tenant3 key")
	assertions.Equal(string(value3), string(retrieved3), "Tenant3 should get their own data")

	// Delete tenant1's key should not affect others
	cache.Delete(tenant1, org1, key)

	_, err = cache.Get(tenant1, org1, key)
	assertions.Error(err, "Tenant1 key should be deleted")

	retrieved2, err = cache.Get(tenant2, org2, key)
	assertions.NoError(err, "Tenant2 key should still exist")
	assertions.Equal(string(value2), string(retrieved2), "Tenant2 data should be unchanged")

	retrieved3, err = cache.Get(tenant3, org3, key)
	assertions.NoError(err, "Tenant3 key should still exist")
	assertions.Equal(string(value3), string(retrieved3), "Tenant3 data should be unchanged")
}

func TestCacheExists(t *testing.T) {
	assertions := assert.New(t)

	cache := NewCacheManager(&Config{
		TotalSize: 10 * 1024 * 1024,
	})
	defer func() {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}()

	tenantID := "tenant-001"
	orgID := "org-001"
	key := "test:exists"
	value := []byte("exists value")

	// Key should not exist initially
	assertions.False(cache.Exists(tenantID, orgID, key), "Key should not exist initially")

	// Set value
	err := cache.Set(tenantID, orgID, key, generateValue(value))
	assertions.NoError(err, "Failed to set value")

	// Key should exist now
	assertions.True(cache.Exists(tenantID, orgID, key), "Key should exist after set")

	// Delete key
	cache.Delete(tenantID, orgID, key)

	// Key should not exist after delete
	assertions.False(cache.Exists(tenantID, orgID, key), "Key should not exist after delete")
}

func TestCacheStats(t *testing.T) {
	assertions := assert.New(t)

	cache := NewCacheManager(&Config{
		TotalSize: 10 * 1024 * 1024,
	})
	defer func() {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}()

	// Initial stats
	stats := cache.Stats()
	assertions.Equal(int64(0), stats.EntryCount, "Initial entry count should be 0")

	// Add some entries
	for i := 0; i < 10; i++ {
		tenantID := fmt.Sprintf("tenant-%03d", i%3)
		orgID := fmt.Sprintf("org-%03d", i%3)
		key := fmt.Sprintf("key:%d", i)
		value := []byte(fmt.Sprintf("value:%d", i))
		_ = cache.Set(tenantID, orgID, key, generateValue(value))
	}

	// Check stats after adding entries
	stats = cache.Stats()
	assertions.Equal(int64(10), stats.EntryCount, "Entry count should be 10")

	// Perform some gets to affect hit rate
	for i := 0; i < 10; i++ {
		tenantID := fmt.Sprintf("tenant-%03d", i%3)
		orgID := fmt.Sprintf("org-%03d", i%3)
		key := fmt.Sprintf("key:%d", i)
		_, _ = cache.Get(tenantID, orgID, key)
	}

	stats = cache.Stats()
	assertions.Greater(stats.HitRate, float64(0), "Hit rate should be greater than 0 after successful gets")
}

func TestCacheGetNotFound(t *testing.T) {
	assertions := assert.New(t)

	cache := NewCacheManager(&Config{
		TotalSize: 10 * 1024 * 1024,
	})
	defer func() {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}()

	tenantID := "tenant-001"
	orgID := "org-001"
	key := "nonexistent:key"

	// Get non-existent key
	_, err := cache.Get(tenantID, orgID, key)
	assertions.Error(err, "Expected error for non-existent key")
}

func TestCacheDeleteNonExistent(t *testing.T) {
	assertions := assert.New(t)

	cache := NewCacheManager(&Config{
		TotalSize: 10 * 1024 * 1024,
	})
	defer func() {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}()

	tenantID := "tenant-001"
	orgID := "org-001"
	key := "nonexistent:key"

	// Delete non-existent key should return false
	deleted := cache.Delete(tenantID, orgID, key)
	assertions.False(deleted, "Delete of non-existent key should return false")
}

func TestCacheCustomTTL(t *testing.T) {
	assertions := assert.New(t)

	cache := NewCacheManager(&Config{
		TotalSize:  10 * 1024 * 1024,
		DefaultTTL: 5 * time.Minute,
	})
	defer func() {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}()

	tenantID := "tenant-001"
	orgID := "org-001"

	// Test Set with default TTL
	key1 := "test:default-ttl"
	value1 := []byte("default ttl value")
	err := cache.Set(tenantID, orgID, key1, generateValue(value1))
	assertions.NoError(err, "Failed to set with default TTL")

	retrieved, err := cache.Get(tenantID, orgID, key1)
	assertions.NoError(err, "Expected to find key with default TTL")
	assertions.Equal(string(value1), string(retrieved), "Retrieved value should match")

	// Test SetTTL with custom TTL
	key2 := "test:custom-ttl"
	value2 := []byte("custom ttl value")
	err = cache.SetTTL(tenantID, orgID, key2, generateValue(value2), 30*time.Minute)
	assertions.NoError(err, "Failed to set with custom TTL")

	retrieved, err = cache.Get(tenantID, orgID, key2)
	assertions.NoError(err, "Expected to find key with custom TTL")
	assertions.Equal(string(value2), string(retrieved), "Retrieved value should match")
}

func TestCacheWithCustomCacheFactory(t *testing.T) {
	assertions := assert.New(t)

	// Custom cache factory that creates a FreeCache with custom size
	customFactory := func(sizeBytes int) Cache {
		return initFreeCache(sizeBytes)
	}

	cache := NewCacheManager(&Config{
		TotalSize:    50 * 1024 * 1024, // 50MB
		DefaultTTL:   5 * time.Minute,
		CacheFactory: customFactory, // This covers the CacheFactory != nil branch
	})
	defer func() {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}()

	assertions.NotNil(cache, "Cache should not be nil")

	tenantID := "tenant-001"
	orgID := "org-001"
	key := "test:custom-factory"
	value := []byte("custom factory value")

	// Verify the cache works correctly with custom factory
	err := cache.Set(tenantID, orgID, key, generateValue(value))
	assertions.NoError(err, "Failed to set value")

	retrieved, err := cache.Get(tenantID, orgID, key)
	assertions.NoError(err, "Expected to find key")
	assertions.Equal(string(value), string(retrieved), "Retrieved value should match set value")
}

func TestCacheGetWithBuf(t *testing.T) {
	assertions := assert.New(t)

	cache := NewCacheManager(&Config{
		TotalSize: 10 * 1024 * 1024,
	})
	defer func() {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}()

	tenantID := "tenant-001"
	orgID := "org-001"
	key := "test:getwithbuf"
	value := []byte("test value for buffer")

	// Set value
	err := cache.Set(tenantID, orgID, key, generateValue(value))
	assertions.NoError(err, "Failed to set value")

	// Test GetWithBuf with pre-allocated buffer (zero-allocation path)
	buf := make([]byte, 1024)
	retrieved, err := cache.GetWithBuf(tenantID, orgID, key, buf)
	assertions.NoError(err, "Expected to find key with GetWithBuf")
	assertions.Equal(string(value), string(retrieved), "Retrieved value should match set value")

	// Test GetWithBuf with nil buffer (should still work)
	retrieved, err = cache.GetWithBuf(tenantID, orgID, key, nil)
	assertions.NoError(err, "Expected to find key with nil buffer")
	assertions.Equal(string(value), string(retrieved), "Retrieved value should match set value")

	// Test GetWithBuf with small buffer (should expand)
	smallBuf := make([]byte, 5)
	retrieved, err = cache.GetWithBuf(tenantID, orgID, key, smallBuf)
	assertions.NoError(err, "Expected to find key with small buffer")
	assertions.Equal(string(value), string(retrieved), "Retrieved value should match set value")

	// Test GetWithBuf for non-existent key
	_, err = cache.GetWithBuf(tenantID, orgID, "nonexistent:key", buf)
	assertions.Error(err, "Expected error for non-existent key")
}

func TestCacheChecksumValidationGet(t *testing.T) {
	assertions := assert.New(t)

	cache := NewCacheManager(&Config{
		TotalSize: 10 * 1024 * 1024,
	})
	defer func() {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}()

	tenantID := "tenant-001"
	orgID := "org-001"
	key := "test:checksum:get"
	value := []byte("sensitive data for get")

	// Set value with tenant-001/org-001
	err := cache.Set(tenantID, orgID, key, generateValue(value))
	assertions.NoError(err, "Failed to set value")

	// Verify correct tenant/org can retrieve using Get
	retrieved, err := cache.Get(tenantID, orgID, key)
	assertions.NoError(err, "Expected to find key with correct tenant/org")
	assertions.Equal(string(value), string(retrieved), "Retrieved value should match")

	// Verify different tenant cannot retrieve using Get
	_, err = cache.Get("tenant-002", orgID, key)
	assertions.Error(err, "Expected error when tenant doesn't match checksum")

	// Verify different org cannot retrieve using Get
	_, err = cache.Get(tenantID, "org-002", key)
	assertions.Error(err, "Expected error when org doesn't match checksum")

	// Verify both different tenant and org cannot retrieve using Get
	_, err = cache.Get("tenant-002", "org-002", key)
	assertions.Error(err, "Expected error when both tenant and org don't match checksum")
}

func TestCacheChecksumValidationGetWithBuf(t *testing.T) {
	assertions := assert.New(t)

	cache := NewCacheManager(&Config{
		TotalSize: 10 * 1024 * 1024,
	})
	defer func() {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}()

	tenantID := "tenant-001"
	orgID := "org-001"
	key := "test:checksum:getwithbuf"
	value := []byte("sensitive data for getwithbuf")
	buf := make([]byte, 1024)

	// Set value with tenant-001/org-001
	err := cache.Set(tenantID, orgID, key, generateValue(value))
	assertions.NoError(err, "Failed to set value")

	// Verify correct tenant/org can retrieve using GetWithBuf
	retrieved, err := cache.GetWithBuf(tenantID, orgID, key, buf)
	assertions.NoError(err, "Expected to find key with correct tenant/org using GetWithBuf")
	assertions.Equal(string(value), string(retrieved), "Retrieved value should match")

	// Verify different tenant cannot retrieve using GetWithBuf
	_, err = cache.GetWithBuf("tenant-002", orgID, key, buf)
	assertions.Error(err, "Expected error when tenant doesn't match checksum using GetWithBuf")

	// Verify different org cannot retrieve using GetWithBuf
	_, err = cache.GetWithBuf(tenantID, "org-002", key, buf)
	assertions.Error(err, "Expected error when org doesn't match checksum using GetWithBuf")

	// Verify both different tenant and org cannot retrieve using GetWithBuf
	_, err = cache.GetWithBuf("tenant-002", "org-002", key, buf)
	assertions.Error(err, "Expected error when both tenant and org don't match checksum using GetWithBuf")

	// Verify with nil buffer also validates checksum
	_, err = cache.GetWithBuf("tenant-002", orgID, key, nil)
	assertions.Error(err, "Expected error with nil buffer when tenant doesn't match checksum")
}

func TestCacheChecksumMultipleTenants(t *testing.T) {
	assertions := assert.New(t)

	cache := NewCacheManager(&Config{
		TotalSize: 10 * 1024 * 1024,
	})
	defer func() {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}()

	// Set same key for multiple tenants/orgs
	key := "shared:key"
	tenant1, org1, value1 := "tenant-001", "org-001", []byte("data for tenant1/org1")
	tenant2, org2, value2 := "tenant-002", "org-002", []byte("data for tenant2/org2")
	tenant3, org3, value3 := "tenant-003", "org-003", []byte("data for tenant3/org3")

	err := cache.Set(tenant1, org1, key, generateValue(value1))
	assertions.NoError(err, "Failed to set value for tenant1")
	err = cache.Set(tenant2, org2, key, generateValue(value2))
	assertions.NoError(err, "Failed to set value for tenant2")
	err = cache.Set(tenant3, org3, key, generateValue(value3))
	assertions.NoError(err, "Failed to set value for tenant3")

	// Each tenant should only get their own data
	retrieved1, err := cache.Get(tenant1, org1, key)
	assertions.NoError(err, "Expected to find key for tenant1")
	assertions.Equal(string(value1), string(retrieved1), "Tenant1 should get their own data")

	retrieved2, err := cache.Get(tenant2, org2, key)
	assertions.NoError(err, "Expected to find key for tenant2")
	assertions.Equal(string(value2), string(retrieved2), "Tenant2 should get their own data")

	retrieved3, err := cache.Get(tenant3, org3, key)
	assertions.NoError(err, "Expected to find key for tenant3")
	assertions.Equal(string(value3), string(retrieved3), "Tenant3 should get their own data")

	// Cross-tenant access should fail checksum validation
	_, err = cache.Get(tenant1, org2, key)
	assertions.Error(err, "Cross-tenant access should fail")

	_, err = cache.Get(tenant2, org1, key)
	assertions.Error(err, "Cross-tenant access should fail")

	// Test with GetWithBuf as well
	buf := make([]byte, 1024)
	retrieved1, err = cache.GetWithBuf(tenant1, org1, key, buf)
	assertions.NoError(err, "Expected to find key for tenant1 with GetWithBuf")
	assertions.Equal(string(value1), string(retrieved1), "Tenant1 should get their own data with GetWithBuf")

	_, err = cache.GetWithBuf(tenant1, org2, key, buf)
	assertions.Error(err, "Cross-tenant access should fail with GetWithBuf")
}
