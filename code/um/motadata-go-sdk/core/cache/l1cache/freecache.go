package l1cache

import (
	"time"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"
	"github.com/coocood/freecache"
)

// FreeCache implements Cache using FreeCache
type FreeCache struct {
	cache  *freecache.Cache
	config *Config
}

// initFreeCache creates a new FreeCache cache
func initFreeCache(config *Config) *FreeCache {
	cache := freecache.NewCache(config.CacheSize)

	return &FreeCache{
		cache:  cache,
		config: config,
	}
}

// Get retrieves a value from FreeCache
func (freeCache *FreeCache) Get(key string) ([]byte, bool) {
	data, err := freeCache.cache.Get(utils.UnsafeStringToBytes(key))
	if err != nil {
		return nil, false
	}
	return data, true
}

// Set stores a value in FreeCache
func (freeCache *FreeCache) Set(key string, value []byte) error {

	return freeCache.cache.Set(utils.UnsafeStringToBytes(key), value, int(freeCache.config.DefaultTTL.Seconds()))
}

func (freeCache *FreeCache) SetTTL(key string, value []byte, ttl time.Duration) error {

	return freeCache.cache.Set(utils.UnsafeStringToBytes(key), value, int(ttl.Seconds()))
}

// Delete removes a key from FreeCache
func (freeCache *FreeCache) Delete(key string) bool {
	return freeCache.cache.Del(utils.UnsafeStringToBytes(key))
}

// Clear removes all entries from FreeCache
func (freeCache *FreeCache) Clear() error {
	freeCache.cache.Clear()
	return nil
}

// Size returns the number of entries in FreeCache
func (freeCache *FreeCache) Size() int {
	return int(freeCache.cache.EntryCount())
}

// Close closes the FreeCache cache
func (freeCache *FreeCache) Close() error {
	return nil // FreeCache doesn't need explicit close
}
