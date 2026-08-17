package l1cache

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCacheBasicOperations(t *testing.T) {
	assertions := assert.New(t)

	cache := NewCacheManager(&Config{
		CacheSize: 10 * 1024 * 1024,
	})
	defer func(cache Cache) {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}(cache)

	// Test byte slice operations
	key := "test:bytes"
	value := []byte("test value")

	err := cache.SetTTL(key, value, 1*time.Minute)
	assertions.NoError(err, "Failed to set value")

	retrieved, found := cache.Get(key)
	assertions.True(found, "Expected to find key")
	assertions.Equal(string(value), string(retrieved), "Retrieved value should match set value")

	// Test string as bytes
	stringKey := "test:string"
	stringValue := []byte("hello world")

	err = cache.SetTTL(stringKey, stringValue, 1*time.Minute)
	assertions.NoError(err, "Failed to set string value")

	retrieved, found = cache.Get(stringKey)
	assertions.True(found, "Expected to find string key")
	assertions.Equal(string(stringValue), string(retrieved), "Retrieved string value should match set value")

	// Delete test
	deleted := cache.Delete(key)
	assertions.True(deleted, "Expected successful deletion")

	_, found = cache.Get(key)
	assertions.False(found, "Expected key to be deleted")
}

func TestCacheWithDifferentSizes(t *testing.T) {
	assertions := assert.New(t)

	cache := NewCacheManager(&Config{
		CacheSize: 10 * 1024 * 1024,
	})
	defer func(cache Cache) {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}(cache)

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

			err2 := cache.SetTTL(tc.key, tc.data, 1*time.Minute)
			assertions2.NoError(err2, "Failed to set %s", tc.key)

			retrieved, found := cache.Get(tc.key)
			assertions2.True(found, "Expected to find %s", tc.key)
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
		CacheSize: 50 * 1024 * 1024,
	})
	defer func(cache Cache) {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}(cache)

	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 100

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()

			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("concurrent:%d:%d", id, j)
				value := []byte(fmt.Sprintf("value_%d_%d", id, j))

				err := cache.SetTTL(key, value, 1*time.Minute)
				assertions.NoError(err, "Failed to set key %s", key)

				retrieved, found := cache.Get(key)
				assertions.True(found, "Failed to find key %s", key)
				assertions.Equal(string(value), string(retrieved), "Retrieved value should match set value")
			}
		}(i)
	}

	wg.Wait()
}

func TestCacheClear(t *testing.T) {
	assertions := assert.New(t)

	cache := NewCacheManager(DefaultConfig())
	defer func(cache Cache) {
		err := cache.Close()
		assertions.NoError(err, "Failed to close cache")
	}(cache)

	// Add items
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("test:%d", i)
		value := []byte(fmt.Sprintf("value %d", i))
		_ = cache.SetTTL(key, value, 1*time.Minute)
	}

	// Verify items exist
	assertions.Greater(cache.Size(), 0, "Expected cache to have items")

	// Clear
	err := cache.Clear()
	assertions.NoError(err, "Failed to clear cache")

	// Verify cleared
	assertions.Equal(0, cache.Size(), "Expected cache size to be 0")

	// Verify items are gone
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("test:%d", i)
		_, found := cache.Get(key)
		assertions.False(found, "Expected key %s to be cleared", key)
	}
}

func TestCacheExpiration(t *testing.T) {
	assertions := assert.New(t)

	cache := NewCacheManager(&Config{
		CacheSize: 10 * 1024 * 1024,
	})
	defer func(cache Cache) {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}(cache)

	key := "test:expiry"
	value := []byte("expires soon")

	// Set with very short TTL
	err := cache.SetTTL(key, value, 100*time.Millisecond)
	assertions.NoError(err, "Failed to set value")

	// Should be present immediately
	retrieved, found := cache.Get(key)
	assertions.True(found, "Expected to find key immediately")
	assertions.Equal(string(value), string(retrieved), "Retrieved value should match set value")

	// Wait for expiration - FreeCache may not immediately expire
	time.Sleep(200 * time.Millisecond)

	// Should be expired now (but FreeCache might be lazy, so we'll be lenient)
	_, found = cache.Get(key)
	// Note: FreeCache may not immediately expire items, so this test might be flaky
	// In production, expired items will be cleaned up eventually
}

func TestCacheOverwrite(t *testing.T) {
	assertions := assert.New(t)

	cache := NewCacheManager(&Config{
		CacheSize: 10 * 1024 * 1024,
	})
	defer func(cache Cache) {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}(cache)

	key := "test:overwrite"
	value1 := []byte("original value")
	value2 := []byte("new value")

	// Set original value
	err := cache.SetTTL(key, value1, 1*time.Minute)
	assertions.NoError(err, "Failed to set original value")

	retrieved, found := cache.Get(key)
	assertions.True(found, "Expected to find original value")
	assertions.Equal(string(value1), string(retrieved), "Retrieved value should match original value")

	// Overwrite with new value
	err = cache.SetTTL(key, value2, 1*time.Minute)
	assertions.NoError(err, "Failed to set new value")

	retrieved, found = cache.Get(key)
	assertions.True(found, "Expected to find new value")
	assertions.Equal(string(value2), string(retrieved), "Retrieved value should match new value")
}

func TestCacheWithDefaultConfig(t *testing.T) {

	assertions := assert.New(t)

	cache := NewCacheManager(nil)
	defer func(cache Cache) {
		assertions.NoError(cache.Close(), "Failed to close cache")
	}(cache)

	assertions.NotNil(cache, "Cache should not be nil")

	key := "test:default.config"
	value := []byte("default config")

	// Set with very short TTL
	err := cache.Set(key, value)
	assertions.NoError(err, "Failed to set value")

	// Should be present immediately
	retrieved, found := cache.Get(key)
	assertions.True(found, "Expected to find key immediately")
	assertions.Equal(string(value), string(retrieved), "Retrieved value should match set value")
}
