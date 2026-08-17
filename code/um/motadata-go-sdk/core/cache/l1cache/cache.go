package l1cache

import (
	"time"
)

// Config holds cache configuration for FreeCache
type Config struct {
	DefaultTTL      time.Duration // Default TTL for items
	CleanupInterval time.Duration // How often to clean expired items
	CacheSize       int           // Cache size in bytes
}

// DefaultConfig returns optimized defaults for FreeCache
func DefaultConfig() *Config {
	return &Config{
		DefaultTTL:      5 * time.Minute,
		CleanupInterval: 30 * time.Second,
		CacheSize:       100 * 1024 * 1024, // 100MB
	}
}

// Cache is a high-performance cache interface that works with raw bytes
type Cache interface {
	// Get retrieves a value as bytes
	Get(key string) ([]byte, bool)

	// Set stores a value as bytes
	Set(key string, value []byte) error

	// SetTTL stores a value as bytes with provided TTL
	SetTTL(key string, value []byte, ttl time.Duration) error

	// Delete removes a key from cache
	Delete(key string) bool

	// Clear removes all entries
	Clear() error

	// Size returns current number of entries
	Size() int

	// Close stops the cache and cleans up resources
	Close() error
}

// CacheManager provides a simple cache implementation that works with raw bytes
type CacheManager struct {
	cache  Cache
	config *Config
}

// NewCacheManager creates a new cache manager that works with raw bytes
func NewCacheManager(config *Config) Cache {
	if config == nil {
		config = DefaultConfig()
	}

	return &CacheManager{
		cache:  initFreeCache(config),
		config: config,
	}
}

// Get retrieves a value as raw bytes
func (cacheManager *CacheManager) Get(key string) ([]byte, bool) {
	return cacheManager.cache.Get(key)
}

// Set stores a value as raw bytes
func (cacheManager *CacheManager) Set(key string, value []byte) error {
	return cacheManager.cache.Set(key, value)
}

// SetTTL stores a value as raw bytes with provided TTL
func (cacheManager *CacheManager) SetTTL(key string, value []byte, ttl time.Duration) error {
	return cacheManager.cache.SetTTL(key, value, ttl)
}

// Delete removes a key from cache
func (cacheManager *CacheManager) Delete(key string) bool {
	return cacheManager.cache.Delete(key)
}

// Clear removes all entries
func (cacheManager *CacheManager) Clear() error {
	return cacheManager.cache.Clear()
}

// Size returns current number of entries
func (cacheManager *CacheManager) Size() int {
	return cacheManager.cache.Size()
}

// Close cleans up the cache and its resources
func (cacheManager *CacheManager) Close() error {
	return cacheManager.cache.Close()
}
