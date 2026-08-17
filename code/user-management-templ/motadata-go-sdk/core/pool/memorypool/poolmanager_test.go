package memorypool

import (
	"testing"

	. "dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"

	"github.com/stretchr/testify/assert"
)

func TestNewPoolManagerGlobal(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   5,
		PoolLength: 10,
		Expandable: true,
		Global:     true,
	}

	pm1 := NewPoolManager(config)
	pm2 := NewPoolManager(config)

	// Use Same to check pointer equality for singleton pattern
	assertions.Same(pm1, pm2, "Global pool managers should return the same instance")
	assertions.NotNil(pm1, "Pool manager should not be nil")
	assertions.True(pm1.global, "Global pool manager should have global field set to true")
}

func TestNewPoolManagerNonGlobal(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   5,
		PoolLength: 10,
		Expandable: true,
		Global:     false,
	}

	pm1 := NewPoolManager(config)
	pm2 := NewPoolManager(config)

	// Use NotSame to check pointer inequality, not deep equality
	assertions.NotSame(pm1, pm2, "Non-global pool managers should be different instances")
	assertions.NotNil(pm1, "First pool manager should not be nil")
	assertions.NotNil(pm2, "Second pool manager should not be nil")
	assertions.False(pm1.global, "Non-global pool manager should have global field set to false")
}

func TestPoolManagerAcquireStringPool(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   3,
		PoolLength: 5,
		Expandable: true,
		Global:     false,
	}

	pm := NewPoolManager(config)

	poolIndex, pool := pm.AcquireStringPool(3)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire string pool")
	assertions.Equal(3, len(pool), "Expected pool length 3")

	pm.ReleaseStringPool(poolIndex)
}

func TestPoolManagerReleaseStringPool(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   3,
		PoolLength: 5,
		Expandable: true,
		Global:     false,
	}

	pm := NewPoolManager(config)

	poolIndex, _ := pm.AcquireStringPool(3)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire string pool")

	pm.ReleaseStringPool(poolIndex)

	poolIndex2, _ := pm.AcquireStringPool(2)
	assertions.NotEqual(NotAvailable, poolIndex2, "Should be able to acquire string pool after release")

	pm.ReleaseStringPool(poolIndex2)
}

func TestPoolManagerGetStringPool(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   3,
		PoolLength: 5,
		Expandable: true,
		Global:     false,
	}

	pm := NewPoolManager(config)

	poolIndex, originalPool := pm.AcquireStringPool(4)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire string pool")

	retrievedPool := pm.GetStringPool(poolIndex)
	assertions.Equal(len(originalPool), len(retrievedPool), "Retrieved pool length should match original")

	pm.ReleaseStringPool(poolIndex)
}

func TestPoolManagerAcquireInt64Pool(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   3,
		PoolLength: 5,
		Expandable: true,
		Global:     false,
	}

	pm := NewPoolManager(config)

	poolIndex, pool := pm.AcquireInt64Pool(4)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire int64 pool")
	assertions.Equal(4, len(pool), "Expected pool length 4")

	pm.ReleaseInt64Pool(poolIndex)
}

func TestPoolManagerReleaseInt64Pool(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   3,
		PoolLength: 5,
		Expandable: true,
		Global:     false,
	}

	pm := NewPoolManager(config)

	poolIndex, _ := pm.AcquireInt64Pool(3)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire int64 pool")

	pm.ReleaseInt64Pool(poolIndex)

	poolIndex2, _ := pm.AcquireInt64Pool(2)
	assertions.NotEqual(NotAvailable, poolIndex2, "Should be able to acquire int64 pool after release")

	pm.ReleaseInt64Pool(poolIndex2)
}

func TestPoolManagerGetInt64Pool(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   3,
		PoolLength: 5,
		Expandable: true,
		Global:     false,
	}

	pm := NewPoolManager(config)

	poolIndex, originalPool := pm.AcquireInt64Pool(4)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire int64 pool")

	retrievedPool := pm.GetInt64Pool(poolIndex)
	assertions.Equal(len(originalPool), len(retrievedPool), "Retrieved pool length should match original")

	pm.ReleaseInt64Pool(poolIndex)
}

func TestPoolManagerAcquireFloat64Pool(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   3,
		PoolLength: 5,
		Expandable: true,
		Global:     false,
	}

	pm := NewPoolManager(config)

	poolIndex, pool := pm.AcquireFloat64Pool(3)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire float64 pool")
	assertions.Equal(3, len(pool), "Expected pool length 3")

	pm.ReleaseFloat64Pool(poolIndex)
}

func TestPoolManagerReleaseFloat64Pool(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   3,
		PoolLength: 5,
		Expandable: true,
		Global:     false,
	}

	pm := NewPoolManager(config)

	poolIndex, _ := pm.AcquireFloat64Pool(3)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire float64 pool")

	pm.ReleaseFloat64Pool(poolIndex)

	poolIndex2, _ := pm.AcquireFloat64Pool(1)
	assertions.NotEqual(NotAvailable, poolIndex2, "Should be able to acquire float64 pool after release")

	pm.ReleaseFloat64Pool(poolIndex2)
}

func TestPoolManagerGetFloat64Pool(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   3,
		PoolLength: 5,
		Expandable: true,
		Global:     false,
	}

	pm := NewPoolManager(config)

	poolIndex, originalPool := pm.AcquireFloat64Pool(3)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire float64 pool")

	retrievedPool := pm.GetFloat64Pool(poolIndex)
	assertions.Equal(len(originalPool), len(retrievedPool), "Retrieved pool length should match original")

	pm.ReleaseFloat64Pool(poolIndex)
}

func TestPoolManagerAcquireBytePool(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   3,
		PoolLength: 5,
		Expandable: true,
		Global:     false,
	}

	pm := NewPoolManager(config)

	poolIndex, pool := pm.AcquireBytePool(5)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire byte pool")
	assertions.Equal(5, len(pool), "Expected pool length 5")

	pm.ReleaseBytePool(poolIndex)
}

func TestPoolManagerReleaseBytePool(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   3,
		PoolLength: 5,
		Expandable: true,
		Global:     false,
	}

	pm := NewPoolManager(config)

	poolIndex, _ := pm.AcquireBytePool(5)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire byte pool")

	pm.ReleaseBytePool(poolIndex)

	poolIndex2, _ := pm.AcquireBytePool(3)
	assertions.NotEqual(NotAvailable, poolIndex2, "Should be able to acquire byte pool after release")

	pm.ReleaseBytePool(poolIndex2)
}

func TestPoolManagerGetBytePool(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   3,
		PoolLength: 5,
		Expandable: true,
		Global:     false,
	}

	pm := NewPoolManager(config)

	poolIndex, originalPool := pm.AcquireBytePool(5)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire byte pool")

	retrievedPool := pm.GetBytePool(poolIndex)
	assertions.Equal(len(originalPool), len(retrievedPool), "Retrieved pool length should match original")

	pm.ReleaseBytePool(poolIndex)
}

func TestPoolManagerShrinkAllPools(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   3,
		PoolLength: 5,
		Expandable: true,
		Global:     false,
	}

	pm := NewPoolManager(config)

	poolIndex1, _ := pm.AcquireStringPool(10)
	poolIndex2, _ := pm.AcquireInt64Pool(8)

	assertions.NotEqual(NotAvailable, poolIndex1, "Should be able to acquire string pool")
	assertions.NotEqual(NotAvailable, poolIndex2, "Should be able to acquire int64 pool")

	pm.ReleaseStringPool(poolIndex1)
	pm.ReleaseInt64Pool(poolIndex2)

	pm.ShrinkAllPools()

	// Verify pools are still functional after shrink
	poolIndex3, buffer3 := pm.AcquireStringPool(3)
	assertions.NotEqual(NotAvailable, poolIndex3, "String pool should be functional after shrink")
	assertions.Equal(3, len(buffer3), "String pool buffer length should be 3")
	pm.ReleaseStringPool(poolIndex3)

	poolIndex4, buffer4 := pm.AcquireInt64Pool(4)
	assertions.NotEqual(NotAvailable, poolIndex4, "Int64 pool should be functional after shrink")
	assertions.Equal(4, len(buffer4), "Int64 pool buffer length should be 4")
	pm.ReleaseInt64Pool(poolIndex4)
}

func TestPoolManagerTestPoolLeak(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   3,
		PoolLength: 5,
		Expandable: false,
		Global:     false,
	}

	pm := NewPoolManager(config)

	poolIndex1, _ := pm.AcquireStringPool(3)
	poolIndex2, _ := pm.AcquireInt64Pool(3)
	poolIndex3, _ := pm.AcquireFloat64Pool(3)
	poolIndex4, _ := pm.AcquireBytePool(3)

	assertions.NotEqual(NotAvailable, poolIndex1, "Should be able to acquire string pool")
	assertions.NotEqual(NotAvailable, poolIndex2, "Should be able to acquire int64 pool")
	assertions.NotEqual(NotAvailable, poolIndex3, "Should be able to acquire float64 pool")
	assertions.NotEqual(NotAvailable, poolIndex4, "Should be able to acquire byte pool")

	statsBeforeLeak := pm.GetPoolStats()
	assertions.NotEqual(0, statsBeforeLeak.UsedStringPools, "String pools should be in use before leak test")
	assertions.NotEqual(0, statsBeforeLeak.UsedInt64Pools, "Int64 pools should be in use before leak test")
	assertions.NotEqual(0, statsBeforeLeak.UsedFloat64Pools, "Float64 pools should be in use before leak test")
	assertions.NotEqual(0, statsBeforeLeak.UsedBytePools, "Byte pools should be in use before leak test")

	pm.TestPoolLeak()

	statsAfterLeak := pm.GetPoolStats()
	assertions.Equal(0, statsAfterLeak.UsedStringPools, "String pools should be reset after leak test")
	assertions.Equal(0, statsAfterLeak.UsedInt64Pools, "Int64 pools should be reset after leak test")
	assertions.Equal(0, statsAfterLeak.UsedFloat64Pools, "Float64 pools should be reset after leak test")
	assertions.Equal(0, statsAfterLeak.UsedBytePools, "Byte pools should be reset after leak test")
}

func TestPoolManagerGetPoolStats(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   5,
		PoolLength: 10,
		Expandable: false,
		Global:     false,
	}

	pm := NewPoolManager(config)

	stats := pm.GetPoolStats()
	assertions.Equal(0, stats.UsedStringPools, "Initially string pool usage should be 0")
	assertions.Equal(0, stats.UsedInt64Pools, "Initially int64 pool usage should be 0")
	assertions.Equal(0, stats.UsedFloat64Pools, "Initially float64 pool usage should be 0")
	assertions.Equal(0, stats.UsedBytePools, "Initially byte pool usage should be 0")

	stringIndex, _ := pm.AcquireStringPool(5)
	int64Index, _ := pm.AcquireInt64Pool(5)

	stats = pm.GetPoolStats()
	assertions.NotEqual(0, stats.UsedStringPools, "String pool usage should be non-zero after acquiring")
	assertions.NotEqual(0, stats.UsedInt64Pools, "Int64 pool usage should be non-zero after acquiring")

	pm.ReleaseStringPool(stringIndex)
	pm.ReleaseInt64Pool(int64Index)

	stats = pm.GetPoolStats()
	assertions.Equal(0, stats.UsedStringPools, "String pool usage should return to 0 after releasing")
	assertions.Equal(0, stats.UsedInt64Pools, "Int64 pool usage should return to 0 after releasing")
}

func TestPoolManagerPoolExhaustion(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   2,
		PoolLength: 5,
		Expandable: false,
		Global:     false,
	}

	pm := NewPoolManager(config)

	poolIndex1, _ := pm.AcquireStringPool(3)
	poolIndex2, _ := pm.AcquireStringPool(3)

	assertions.NotEqual(NotAvailable, poolIndex1, "Should be able to acquire first pool")
	assertions.NotEqual(NotAvailable, poolIndex2, "Should be able to acquire second pool")

	poolIndex3, pool3 := pm.AcquireStringPool(3)
	assertions.Equal(NotAvailable, poolIndex3, "Should not be able to acquire pool when exhausted")
	assertions.Nil(pool3, "Pool should be nil when exhausted")

	pm.ReleaseStringPool(poolIndex1)
	poolIndex4, pool4 := pm.AcquireStringPool(2)
	assertions.NotEqual(NotAvailable, poolIndex4, "Should be able to acquire pool after release")
	assertions.NotNil(pool4, "Pool should not be nil after release")

	pm.ReleaseStringPool(poolIndex2)
	pm.ReleaseStringPool(poolIndex4)
}

func TestPoolManagerConcurrentAccessNonGlobal(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   10,
		PoolLength: 5,
		Expandable: false,
		Global:     false,
	}

	pm := NewPoolManager(config)

	// Test that non-global pools work correctly with sequential access
	// (concurrent access without locking would cause race conditions)
	for i := 0; i < 100; i++ {
		poolIndex, _ := pm.AcquireStringPool(3)
		if poolIndex != NotAvailable {
			pm.ReleaseStringPool(poolIndex)
		}
	}

	// Verify pool is still functional
	poolIndex, pool := pm.AcquireStringPool(3)
	assertions.NotEqual(NotAvailable, poolIndex, "Pool should still be functional after sequential operations")
	assertions.Equal(3, len(pool), "Pool length should be 3")
	pm.ReleaseStringPool(poolIndex)
}

func TestPoolManagerConcurrentAccessGlobal(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   10,
		PoolLength: 5,
		Expandable: false,
		Global:     true,
	}

	pm := NewPoolManager(config)

	done := make(chan bool)
	numGoroutines := 5

	for i := 0; i < numGoroutines; i++ {
		go func() {
			for j := 0; j < 50; j++ {
				poolIndex, _ := pm.AcquireStringPool(3)

				assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire pool in goroutine")

				pm.ReleaseStringPool(poolIndex)
			}
			done <- true
		}()
	}

	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}

func TestPoolManagerEdgeCases(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   1,
		PoolLength: 1,
		Expandable: false,
		Global:     false,
	}

	pm := NewPoolManager(config)

	poolIndex, pool := pm.AcquireStringPool(1)
	assertions.NotEqual(NotAvailable, poolIndex, "Should handle minimum pool size")
	assertions.Equal(1, len(pool), "Pool length should be 1")

	pm.ReleaseStringPool(poolIndex)

	poolIndex2, pool2 := pm.AcquireStringPool(NotAvailable)
	assertions.NotEqual(NotAvailable, poolIndex2, "Should handle default size request")
	assertions.Equal(1, len(pool2), "Default size should be 1")

	pm.ReleaseStringPool(poolIndex2)
}

func TestPoolManagerAllPoolTypesSimultaneous(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   5,
		PoolLength: 10,
		Expandable: true,
		Global:     false,
	}

	pm := NewPoolManager(config)

	stringIndex, stringPool := pm.AcquireStringPool(5)
	int64Index, int64Pool := pm.AcquireInt64Pool(7)
	float64Index, float64Pool := pm.AcquireFloat64Pool(3)
	byteIndex, bytePool := pm.AcquireBytePool(9)

	assertions.NotEqual(NotAvailable, stringIndex, "String pool acquisition should succeed")
	assertions.Equal(5, len(stringPool), "String pool length should be 5")

	assertions.NotEqual(NotAvailable, int64Index, "Int64 pool acquisition should succeed")
	assertions.Equal(7, len(int64Pool), "Int64 pool length should be 7")

	assertions.NotEqual(NotAvailable, float64Index, "Float64 pool acquisition should succeed")
	assertions.Equal(3, len(float64Pool), "Float64 pool length should be 3")

	assertions.NotEqual(NotAvailable, byteIndex, "Byte pool acquisition should succeed")
	assertions.Equal(9, len(bytePool), "Byte pool length should be 9")

	stats := pm.GetPoolStats()
	assertions.NotEqual(0, stats.UsedStringPools, "String pool should show usage")
	assertions.NotEqual(0, stats.UsedInt64Pools, "Int64 pool should show usage")
	assertions.NotEqual(0, stats.UsedFloat64Pools, "Float64 pool should show usage")
	assertions.NotEqual(0, stats.UsedBytePools, "Byte pool should show usage")

	pm.ReleaseStringPool(stringIndex)
	pm.ReleaseInt64Pool(int64Index)
	pm.ReleaseFloat64Pool(float64Index)
	pm.ReleaseBytePool(byteIndex)

	stats = pm.GetPoolStats()
	assertions.Equal(0, stats.UsedStringPools, "String pool should be released")
	assertions.Equal(0, stats.UsedInt64Pools, "Int64 pool should be released")
	assertions.Equal(0, stats.UsedFloat64Pools, "Float64 pool should be released")
	assertions.Equal(0, stats.UsedBytePools, "Byte pool should be released")
}

// TestPoolManagerGlobalSingletonBehavior tests global singleton pattern edge cases
func TestPoolManagerGlobalSingletonBehavior(t *testing.T) {
	assertions := assert.New(t)

	// Test multiple configurations with global=true should return same instance
	config1 := &PoolConfig{
		PoolSize:   5,
		PoolLength: 10,
		Expandable: true,
		Global:     true,
	}

	config2 := &PoolConfig{
		PoolSize:   10, // Different config
		PoolLength: 20,
		Expandable: false,
		Global:     true,
	}

	pm1 := NewPoolManager(config1)
	pm2 := NewPoolManager(config2) // Should ignore config2 and return same instance

	assertions.Same(pm1, pm2, "Global pool managers should return same instance regardless of config")

	// Test that the original config is preserved
	stringIndex, _ := pm2.AcquireStringPool(5)
	assertions.NotEqual(NotAvailable, stringIndex, "Should work with original config")
	pm2.ReleaseStringPool(stringIndex)

	// Test accessing internal pools through both references
	stringIndex1, _ := pm1.AcquireStringPool(3)
	stringIndex2, _ := pm2.AcquireStringPool(4)

	assertions.NotEqual(stringIndex1, stringIndex2, "Should get different pool indices")

	pm1.ReleaseStringPool(stringIndex1)
	pm2.ReleaseStringPool(stringIndex2)
}

// TestPoolManagerNonGlobalInstanceSeparation tests non-global instance separation
func TestPoolManagerNonGlobalInstanceSeparation(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   3,
		PoolLength: 5,
		Expandable: false,
		Global:     false,
	}

	pm1 := NewPoolManager(config)
	pm2 := NewPoolManager(config)

	assertions.NotSame(pm1, pm2, "Non-global instances should be different")

	// Test that they maintain separate state
	stringIndex1, _ := pm1.AcquireStringPool(3)
	stringIndex2, _ := pm2.AcquireStringPool(3)

	// Both should be able to use index 0 since they're separate instances
	assertions.Equal(0, stringIndex1, "First instance should get index 0")
	assertions.Equal(0, stringIndex2, "Second instance should also get index 0")

	stats1 := pm1.GetPoolStats()
	stats2 := pm2.GetPoolStats()

	assertions.NotEqual(0, stats1.UsedStringPools, "First instance should show pool usage")
	assertions.NotEqual(0, stats2.UsedStringPools, "Second instance should show pool usage")

	pm1.ReleaseStringPool(stringIndex1)
	pm2.ReleaseStringPool(stringIndex2)
}

// TestPoolManagerRWLockBehaviorGlobal tests read-write lock behavior in global mode
func TestPoolManagerRWLockBehaviorGlobal(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   5,
		PoolLength: 10,
		Expandable: true,
		Global:     true,
	}

	pm := NewPoolManager(config)

	// Test that GetPool operations use read locks (should not block each other)
	stringIndex, _ := pm.AcquireStringPool(5)
	int64Index, _ := pm.AcquireInt64Pool(6)

	// These should work fine even with pools acquired (read locks)
	retrievedString := pm.GetStringPool(stringIndex)
	retrievedInt64 := pm.GetInt64Pool(int64Index)
	stats := pm.GetPoolStats()

	assertions.Equal(5, len(retrievedString), "Should be able to get string pool info")
	assertions.Equal(6, len(retrievedInt64), "Should be able to get int64 pool info")
	assertions.NotEqual(0, stats.UsedStringPools, "Stats should show usage")

	pm.ReleaseStringPool(stringIndex)
	pm.ReleaseInt64Pool(int64Index)
}

// TestPoolManagerShrinkAllPoolsExtensive tests comprehensive shrink behavior
func TestPoolManagerShrinkAllPoolsExtensive(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   3,
		PoolLength: 5,
		Expandable: true,
		Global:     false,
	}

	pm := NewPoolManager(config)

	// Expand all pool types
	stringIndex, _ := pm.AcquireStringPool(15)
	int64Index, _ := pm.AcquireInt64Pool(12)
	float64Index, _ := pm.AcquireFloat64Pool(18)
	byteIndex, _ := pm.AcquireBytePool(20)

	assertions.NotEqual(NotAvailable, stringIndex, "String pool expansion should work")
	assertions.NotEqual(NotAvailable, int64Index, "Int64 pool expansion should work")
	assertions.NotEqual(NotAvailable, float64Index, "Float64 pool expansion should work")
	assertions.NotEqual(NotAvailable, byteIndex, "Byte pool expansion should work")

	// Release all pools
	pm.ReleaseStringPool(stringIndex)
	pm.ReleaseInt64Pool(int64Index)
	pm.ReleaseFloat64Pool(float64Index)
	pm.ReleaseBytePool(byteIndex)

	// Shrink all pools
	pm.ShrinkAllPools()

	// Test that all pools still work after shrinking
	stringIndex2, stringBuffer := pm.AcquireStringPool(4)
	int64Index2, int64Buffer := pm.AcquireInt64Pool(3)
	float64Index2, float64Buffer := pm.AcquireFloat64Pool(5)
	byteIndex2, byteBuffer := pm.AcquireBytePool(2)

	assertions.NotEqual(NotAvailable, stringIndex2, "String pool should work after shrink")
	assertions.Equal(4, len(stringBuffer), "String buffer should have correct size")

	assertions.NotEqual(NotAvailable, int64Index2, "Int64 pool should work after shrink")
	assertions.Equal(3, len(int64Buffer), "Int64 buffer should have correct size")

	assertions.NotEqual(NotAvailable, float64Index2, "Float64 pool should work after shrink")
	assertions.Equal(5, len(float64Buffer), "Float64 buffer should have correct size")

	assertions.NotEqual(NotAvailable, byteIndex2, "Byte pool should work after shrink")
	assertions.Equal(2, len(byteBuffer), "Byte buffer should have correct size")

	pm.ReleaseStringPool(stringIndex2)
	pm.ReleaseInt64Pool(int64Index2)
	pm.ReleaseFloat64Pool(float64Index2)
	pm.ReleaseBytePool(byteIndex2)
}

// TestPoolManagerTestPoolLeakComprehensive tests comprehensive leak detection
func TestPoolManagerTestPoolLeakComprehensive(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   5,
		PoolLength: 10,
		Expandable: false,
		Global:     false,
	}

	pm := NewPoolManager(config)

	// Acquire multiple pools from each type without releasing
	stringIndex1, _ := pm.AcquireStringPool(5)
	stringIndex2, _ := pm.AcquireStringPool(3)
	int64Index1, _ := pm.AcquireInt64Pool(7)
	float64Index1, _ := pm.AcquireFloat64Pool(4)
	byteIndex1, _ := pm.AcquireBytePool(8)
	byteIndex2, _ := pm.AcquireBytePool(6)

	// Verify all acquisitions succeeded
	assertions.NotEqual(NotAvailable, stringIndex1, "String pool 1 should be acquired")
	assertions.NotEqual(NotAvailable, stringIndex2, "String pool 2 should be acquired")
	assertions.NotEqual(NotAvailable, int64Index1, "Int64 pool should be acquired")
	assertions.NotEqual(NotAvailable, float64Index1, "Float64 pool should be acquired")
	assertions.NotEqual(NotAvailable, byteIndex1, "Byte pool 1 should be acquired")
	assertions.NotEqual(NotAvailable, byteIndex2, "Byte pool 2 should be acquired")

	// Check stats before leak test
	statsBefore := pm.GetPoolStats()
	assertions.NotEqual(0, statsBefore.UsedStringPools, "String pools should be in use")
	assertions.NotEqual(0, statsBefore.UsedInt64Pools, "Int64 pools should be in use")
	assertions.NotEqual(0, statsBefore.UsedFloat64Pools, "Float64 pools should be in use")
	assertions.NotEqual(0, statsBefore.UsedBytePools, "Byte pools should be in use")

	// Test pool leak - should reset all pools
	pm.TestPoolLeak()

	// Check stats after leak test
	statsAfter := pm.GetPoolStats()
	assertions.Equal(0, statsAfter.UsedStringPools, "String pools should be reset")
	assertions.Equal(0, statsAfter.UsedInt64Pools, "Int64 pools should be reset")
	assertions.Equal(0, statsAfter.UsedFloat64Pools, "Float64 pools should be reset")
	assertions.Equal(0, statsAfter.UsedBytePools, "Byte pools should be reset")

	// Verify all pools are functional after leak test
	newStringIndex, newStringBuffer := pm.AcquireStringPool(3)
	newInt64Index, newInt64Buffer := pm.AcquireInt64Pool(4)

	assertions.NotEqual(NotAvailable, newStringIndex, "String pool should work after leak test")
	assertions.Equal(3, len(newStringBuffer), "String buffer should have correct size")

	assertions.NotEqual(NotAvailable, newInt64Index, "Int64 pool should work after leak test")
	assertions.Equal(4, len(newInt64Buffer), "Int64 buffer should have correct size")

	pm.ReleaseStringPool(newStringIndex)
	pm.ReleaseInt64Pool(newInt64Index)
}

// TestPoolManagerHighConcurrencyGlobal tests high concurrency with global pools
func TestPoolManagerHighConcurrencyGlobal(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   10,
		PoolLength: 5,
		Expandable: false,
		Global:     true,
	}

	pm := NewPoolManager(config)

	done := make(chan bool, 20)
	numGoroutines := 20
	operationsPerGoroutine := 100

	// Launch multiple goroutines for stress testing
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < operationsPerGoroutine; j++ {
				// Random pool operations
				switch j % 4 {
				case 0:
					if stringIndex, _ := pm.AcquireStringPool(3); stringIndex != NotAvailable {
						pm.ReleaseStringPool(stringIndex)
					}
				case 1:
					if int64Index, _ := pm.AcquireInt64Pool(4); int64Index != NotAvailable {
						pm.ReleaseInt64Pool(int64Index)
					}
				case 2:
					if float64Index, _ := pm.AcquireFloat64Pool(2); float64Index != NotAvailable {
						pm.ReleaseFloat64Pool(float64Index)
					}
				case 3:
					if byteIndex, _ := pm.AcquireBytePool(5); byteIndex != NotAvailable {
						pm.ReleaseBytePool(byteIndex)
					}
				}

				// Occasionally read stats and pool data
				if j%10 == 0 {
					pm.GetPoolStats()
				}
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify final state
	finalStats := pm.GetPoolStats()
	assertions.Equal(0, finalStats.UsedStringPools, "No string pools should be in use")
	assertions.Equal(0, finalStats.UsedInt64Pools, "No int64 pools should be in use")
	assertions.Equal(0, finalStats.UsedFloat64Pools, "No float64 pools should be in use")
	assertions.Equal(0, finalStats.UsedBytePools, "No byte pools should be in use")

	// Verify pools still work after concurrency test
	stringIndex, stringBuffer := pm.AcquireStringPool(3)
	assertions.NotEqual(NotAvailable, stringIndex, "String pool should still work")
	assertions.Equal(3, len(stringBuffer), "String buffer should have correct size")
	pm.ReleaseStringPool(stringIndex)
}

// TestPoolManagerPoolTypeConsistency tests pool type consistency and behavior
func TestPoolManagerPoolTypeConsistency(t *testing.T) {
	assertions := assert.New(t)

	config := &PoolConfig{
		PoolSize:   3,
		PoolLength: 5,
		Expandable: true,
		Global:     false,
	}

	pm := NewPoolManager(config)

	// Test that each pool type maintains its own separate state
	stringPoolIndex, stringBuffer1 := pm.AcquireStringPool(3)
	stringPoolIndex2, stringBuffer2 := pm.AcquireStringPool(4)
	int64PoolIndex, int64Buffer1 := pm.AcquireInt64Pool(3)
	int64PoolIndex2, int64Buffer2 := pm.AcquireInt64Pool(4)

	// Modify buffers to test independence
	stringBuffer1[0] = "string1"
	stringBuffer2[0] = "string2"
	int64Buffer1[0] = 100
	int64Buffer2[0] = 200

	assertions.Equal("string1", stringBuffer1[0], "String buffer 1 should retain value")
	assertions.Equal("string2", stringBuffer2[0], "String buffer 2 should retain value")
	assertions.Equal(int64(100), int64Buffer1[0], "Int64 buffer 1 should retain value")
	assertions.Equal(int64(200), int64Buffer2[0], "Int64 buffer 2 should retain value")

	// Release and verify independence
	pm.ReleaseStringPool(stringPoolIndex)
	pm.ReleaseInt64Pool(int64PoolIndex)

	// The remaining buffers should still have their values
	assertions.Equal("string2", stringBuffer2[0], "String buffer 2 should still have value")
	assertions.Equal(int64(200), int64Buffer2[0], "Int64 buffer 2 should still have value")

	pm.ReleaseStringPool(stringPoolIndex2)
	pm.ReleaseInt64Pool(int64PoolIndex2)

	// Re-acquire and verify zero initialization
	stringPoolIndex, stringBuffer := pm.AcquireStringPool(3)
	int64PoolIndex, int64Buffer := pm.AcquireInt64Pool(3)

	assertions.Equal("", stringBuffer[0], "New string buffer should be zeroed")
	assertions.Equal(int64(0), int64Buffer[0], "New int64 buffer should be zeroed")

	pm.ReleaseStringPool(stringPoolIndex)
	pm.ReleaseInt64Pool(int64PoolIndex)
}
