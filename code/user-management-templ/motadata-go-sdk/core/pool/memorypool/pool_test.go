package memorypool

import (
	"strings"
	"testing"

	. "dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"

	"github.com/stretchr/testify/assert"
)

func TestNewMemoryPoolString(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[string](5, 10, true)

	assertions.NotNil(memoryPool, "Pool should not be nil")
	assertions.Equal(5, memoryPool.poolSize, "Expected pool size 5")
	assertions.Equal(10, memoryPool.poolLength, "Expected pool length 10")
	assertions.True(memoryPool.expandable, "Expected pool to be expandable")
	assertions.Equal(int64(0), memoryPool.usedPools, "Expected used pools to be 0")
}

func TestNewMemoryPoolInt64(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[int64](3, 8, false)

	assertions.NotNil(memoryPool, "Pool should not be nil")
	assertions.Equal(3, memoryPool.poolSize, "Expected pool size 3")
	assertions.Equal(8, memoryPool.poolLength, "Expected pool length 8")
	assertions.False(memoryPool.expandable, "Expected pool to be non-expandable")
}

func TestNewMemoryPoolFloat64(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[float64](7, 15, true)

	assertions.NotNil(memoryPool, "Pool should not be nil")
	assertions.Equal(7, memoryPool.poolSize, "Expected pool size 7")
	assertions.Equal(15, memoryPool.poolLength, "Expected pool length 15")
	assertions.True(memoryPool.expandable, "Expected pool to be expandable")
}

func TestNewMemoryPoolByte(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[byte](4, 12, false)

	assertions.NotNil(memoryPool, "Pool should not be nil")
	assertions.Equal(4, memoryPool.poolSize, "Expected pool size 4")
	assertions.Equal(12, memoryPool.poolLength, "Expected pool length 12")
	assertions.False(memoryPool.expandable, "Expected pool to be non-expandable")
}

func TestMemoryPoolAcquirePoolString(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[string](3, 5, false)

	poolIndex, buffer := memoryPool.acquirePool(3)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire pool")
	assertions.Equal(3, len(buffer), "Expected buffer length 3")

	memoryPool.releasePool(poolIndex)
}

func TestMemoryPoolAcquirePoolInt64(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[int64](3, 5, false)

	poolIndex, buffer := memoryPool.acquirePool(4)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire pool")
	assertions.Equal(4, len(buffer), "Expected buffer length 4")

	memoryPool.releasePool(poolIndex)
}

func TestMemoryPoolAcquirePoolFloat64(t *testing.T) {
	assertions := assert.New(t)

	pool := newMemoryPool[float64](3, 5, false)

	poolIndex, buffer := pool.acquirePool(2)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire pool")
	assertions.Equal(2, len(buffer), "Expected buffer length 2")

	pool.releasePool(poolIndex)
}

func TestMemoryPoolAcquirePoolByte(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[byte](3, 5, false)

	poolIndex, buffer := memoryPool.acquirePool(5)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire pool")
	assertions.Equal(5, len(buffer), "Expected buffer length 5")

	memoryPool.releasePool(poolIndex)
}

func TestMemoryPoolReleasePool(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[string](3, 5, false)

	poolIndex, _ := memoryPool.acquirePool(3)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire pool")

	memoryPool.releasePool(poolIndex)

	poolIndex2, _ := memoryPool.acquirePool(2)
	assertions.NotEqual(NotAvailable, poolIndex2, "Should be able to acquire pool after release")

	memoryPool.releasePool(poolIndex2)
}

func TestMemoryPoolGetPool(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[string](3, 5, false)

	poolIndex, originalBuffer := memoryPool.acquirePool(4)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire pool")

	retrievedBuffer := memoryPool.getPool(poolIndex)
	assertions.Equal(len(originalBuffer), len(retrievedBuffer), "Retrieved buffer length should match original")

	memoryPool.releasePool(poolIndex)
}

func TestMemoryPoolExpandPool(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[string](3, 5, true)

	poolIndex, _ := memoryPool.acquirePool(3)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire pool")

	expandedBuffer := memoryPool.expandPool(poolIndex, 8)
	assertions.Equal(8, len(expandedBuffer), "Expected expanded buffer length 8")

	memoryPool.releasePool(poolIndex)
}

func TestMemoryPoolGetUsedPools(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[string](5, 10, false)

	assertions.Equal(0, memoryPool.getUsedPools(), "Initially used pools should be 0")

	poolIndex1, _ := memoryPool.acquirePool(5)
	poolIndex2, _ := memoryPool.acquirePool(3)

	assertions.NotEqual(NotAvailable, poolIndex1, "Should be able to acquire first pool")
	assertions.NotEqual(NotAvailable, poolIndex2, "Should be able to acquire second pool")

	usedPools := memoryPool.getUsedPools()
	// usedPools is a bitmask, not a count. With indices 0 and 1 acquired, it should be 0b11 = 3
	assertions.Equal(3, usedPools, "Used pools bitmask should be 3 (0b11) after acquiring two pools at indices 0 and 1")

	memoryPool.releasePool(poolIndex1)
	memoryPool.releasePool(poolIndex2)

	assertions.Equal(0, memoryPool.getUsedPools(), "Used pools should be 0 after releasing all")
}

func TestMemoryPoolGetPoolLength(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[string](3, 15, false)

	assertions.Equal(15, memoryPool.getPoolLength(), "Expected pool length 15")
}

func TestMemoryPoolGetPoolSize(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[string](7, 10, false)

	assertions.Equal(7, memoryPool.getPoolSize(), "Expected pool size 7")
}

func TestMemoryPoolIsExpandableTrue(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[string](3, 5, true)

	assertions.True(memoryPool.expandable, "Pool should be expandable")
}

func TestMemoryPoolIsExpandableFalse(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[string](3, 5, false)

	assertions.False(memoryPool.expandable, "Pool should not be expandable")
}

func TestMemoryPoolShrinkPool(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[string](3, 5, true)

	poolIndex, buffer := memoryPool.acquirePool(10)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire pool")
	assertions.Equal(10, len(buffer), "Expected buffer length 10")

	memoryPool.releasePool(poolIndex)

	memoryPool.shrinkPool()
	// Verify pool is still functional after shrink
	poolIndex2, buffer2 := memoryPool.acquirePool(3)
	assertions.NotEqual(NotAvailable, poolIndex2, "Pool should be functional after shrink")
	assertions.Equal(3, len(buffer2), "Expected buffer length 3 after shrink")
	memoryPool.releasePool(poolIndex2)
}

func TestMemoryPoolShrinkPoolNonExpandable(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[string](3, 5, false)

	poolIndex, buffer := memoryPool.acquirePool(3)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire pool")
	assertions.Equal(3, len(buffer), "Expected buffer length 3")

	memoryPool.releasePool(poolIndex)

	memoryPool.shrinkPool()
	// Verify pool is still functional after shrink
	poolIndex2, buffer2 := memoryPool.acquirePool(4)
	assertions.NotEqual(NotAvailable, poolIndex2, "Pool should be functional after shrink")
	assertions.Equal(4, len(buffer2), "Expected buffer length 4 after shrink")
	memoryPool.releasePool(poolIndex2)
}

func TestMemoryPoolTestPoolLeak(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[string](3, 5, false)

	poolIndex1, _ := memoryPool.acquirePool(3)
	poolIndex2, _ := memoryPool.acquirePool(2)

	assertions.NotEqual(NotAvailable, poolIndex1, "Should be able to acquire first pool")
	assertions.NotEqual(NotAvailable, poolIndex2, "Should be able to acquire second pool")

	usedBefore := memoryPool.getUsedPools()
	// usedPools is a bitmask, not a count. With indices 0 and 1 acquired, it should be 0b11 = 3
	assertions.Equal(3, usedBefore, "Used pools bitmask should be 3 (0b11) before leak test")

	memoryPool.testPoolLeak()

	assertions.Equal(0, memoryPool.getUsedPools(), "Pool leak test should reset used pools")
}

func TestMemoryPoolPoolExhaustion(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[string](2, 5, false)

	poolIndex1, _ := memoryPool.acquirePool(3)
	poolIndex2, _ := memoryPool.acquirePool(4)

	assertions.NotEqual(NotAvailable, poolIndex1, "Should be able to acquire first pool")
	assertions.NotEqual(NotAvailable, poolIndex2, "Should be able to acquire second pool")

	poolIndex3, buffer3 := memoryPool.acquirePool(2)
	assertions.Equal(NotAvailable, poolIndex3, "Should not be able to acquire pool when exhausted")
	assertions.Nil(buffer3, "Buffer should be nil when pool exhausted")

	memoryPool.releasePool(poolIndex1)
	poolIndex4, buffer4 := memoryPool.acquirePool(2)
	assertions.NotEqual(NotAvailable, poolIndex4, "Should be able to acquire pool after release")
	assertions.NotNil(buffer4, "Buffer should not be nil after release")

	memoryPool.releasePool(poolIndex2)
	memoryPool.releasePool(poolIndex4)
}

func TestMemoryPoolExpandableVsNonExpandable(t *testing.T) {
	assertions := assert.New(t)

	expandablePool := newMemoryPool[string](3, 5, true)
	nonExpandablePool := newMemoryPool[string](3, 5, false)

	expandablePoolIndex, expandableBuffer := expandablePool.acquirePool(8)
	nonExpandablePoolIndex, nonExpandableBuffer := nonExpandablePool.acquirePool(8)

	assertions.NotEqual(NotAvailable, expandablePoolIndex, "Expandable pool should handle large sizes")
	assertions.Equal(8, len(expandableBuffer), "Expandable pool buffer should be 8")

	assertions.NotEqual(NotAvailable, nonExpandablePoolIndex, "Non-expandable pool should still acquire but limit size")
	assertions.Equal(5, len(nonExpandableBuffer), "Non-expandable pool buffer should be limited to 5")

	expandablePool.releasePool(expandablePoolIndex)
	nonExpandablePool.releasePool(nonExpandablePoolIndex)
}

func TestMemoryPoolDefaultSizeRequest(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[string](3, 7, false)

	poolIndex, buffer := memoryPool.acquirePool(NotAvailable)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire pool with default size")
	assertions.Equal(7, len(buffer), "Default size should be pool length 7")

	memoryPool.releasePool(poolIndex)
}

func TestMemoryPoolZeroValues(t *testing.T) {
	assertions := assert.New(t)

	stringPool := newMemoryPool[string](3, 5, false)
	int64Pool := newMemoryPool[int64](3, 5, false)
	float64Pool := newMemoryPool[float64](3, 5, false)
	bytePool := newMemoryPool[byte](3, 5, false)

	stringPoolIndex, stringBuffer := stringPool.acquirePool(3)
	int64PoolIndex, int64Buffer := int64Pool.acquirePool(3)
	float64PoolIndex, float64Buffer := float64Pool.acquirePool(3)
	bytePoolIndex, byteBuffer := bytePool.acquirePool(3)

	for i := 0; i < 3; i++ {
		assertions.Equal(Empty, stringBuffer[i], "string buffer should be initialized to zero values")
		assertions.Equal(int64(0), int64Buffer[i], "int64 buffer should be initialized to zero values")
		assertions.Equal(0.0, float64Buffer[i], "float64 buffer should be initialized to zero values")
		assertions.Equal(byte(0), byteBuffer[i], "byte buffer should be initialized to zero values")
	}

	stringPool.releasePool(stringPoolIndex)
	int64Pool.releasePool(int64PoolIndex)
	float64Pool.releasePool(float64PoolIndex)
	bytePool.releasePool(bytePoolIndex)
}

func TestMemoryPoolConcurrentOperationsStringPool(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[string](10, 5, false)

	// Test sequential operations (Pool is not thread-safe by design)
	for i := 0; i < 300; i++ {
		poolIndex, _ := memoryPool.acquirePool(3)
		if poolIndex != NotAvailable {
			memoryPool.releasePool(poolIndex)
		}
	}

	// Verify pool is still functional
	poolIndex, buffer := memoryPool.acquirePool(5)
	assertions.NotEqual(NotAvailable, poolIndex, "Pool should still be functional after sequential operations")
	assertions.Equal(5, len(buffer), "Buffer length should be 3")
	memoryPool.releasePool(poolIndex)
}

func TestMemoryPoolConcurrentOperationsInt64Pool(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[int64](10, 5, false)

	// Test sequential operations (Pool is not thread-safe by design)
	for i := 0; i < 200; i++ {
		poolIndex, _ := memoryPool.acquirePool(3)
		if poolIndex != NotAvailable {
			memoryPool.releasePool(poolIndex)
		}
	}

	// Verify pool is still functional
	poolIndex, buffer := memoryPool.acquirePool(3)
	assertions.NotEqual(NotAvailable, poolIndex, "Pool should still be functional after sequential operations")
	assertions.Equal(3, len(buffer), "Buffer length should be 3")
	memoryPool.releasePool(poolIndex)
}

func TestMemoryPoolConcurrentOperationsFloat64Pool(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[float64](10, 5, false)

	// Test sequential operations (Pool is not thread-safe by design)
	for i := 0; i < 234; i++ {
		poolIndex, _ := memoryPool.acquirePool(3)
		if poolIndex != NotAvailable {
			memoryPool.releasePool(poolIndex)
		}
	}

	// Verify pool is still functional
	poolIndex, buffer := memoryPool.acquirePool(3)
	assertions.NotEqual(NotAvailable, poolIndex, "Pool should still be functional after sequential operations")
	assertions.Equal(3, len(buffer), "Buffer length should be 3")
	memoryPool.releasePool(poolIndex)
}

func TestMemoryPoolConcurrentOperationsBytePool(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[byte](10, 5, false)

	// Test sequential operations (Pool is not thread-safe by design)
	for i := 0; i < 250; i++ {
		poolIndex, _ := memoryPool.acquirePool(3)
		if poolIndex != NotAvailable {
			memoryPool.releasePool(poolIndex)
		}
	}

	// Verify pool is still functional
	poolIndex, buffer := memoryPool.acquirePool(3)
	assertions.NotEqual(NotAvailable, poolIndex, "Pool should still be functional after sequential operations")
	assertions.Equal(3, len(buffer), "Buffer length should be 3")
	memoryPool.releasePool(poolIndex)
}

func TestMemoryPoolEdgeCases(t *testing.T) {
	t.Run("MinimumSizes", func(t *testing.T) {
		assertions := assert.New(t)

		memoryPool := newMemoryPool[string](1, 1, false)

		poolIndex, buffer := memoryPool.acquirePool(1)
		assertions.NotEqual(NotAvailable, poolIndex, "Should handle minimum pool size")
		assertions.Equal(1, len(buffer), "Buffer length should be 1")

		memoryPool.releasePool(poolIndex)
	})

	t.Run("ZeroLengthRequest", func(t *testing.T) {
		assertions := assert.New(t)

		memoryPool := newMemoryPool[string](3, 5, false)

		poolIndex, buffer := memoryPool.acquirePool(0)
		assertions.NotEqual(NotAvailable, poolIndex, "Should handle zero length request")
		assertions.Equal(0, len(buffer), "Zero length request should return zero length buffer")

		memoryPool.releasePool(poolIndex)
	})

	t.Run("LargePoolSize", func(t *testing.T) {
		assertions := assert.New(t)

		memoryPool := newMemoryPool[string](100, 10, false)

		assertions.Equal(100, memoryPool.poolSize, "Should handle large pool size")
	})
}

// TestMemoryPoolUnknownDataType tests the unknown type scenario
func TestMemoryPoolUnknownDataType(t *testing.T) {
	assertions := assert.New(t)

	// Test with an unsupported/unknown type - using interface{} which doesn't match predefined types
	type UnknownType int

	memoryPool := newMemoryPool[UnknownType](3, 5, false)
	assertions.NotNil(memoryPool, "Pool should be created even with unknown type")
	assertions.Equal("unknown", strings.ToLower(memoryPool.dataType), "Data type should be marked as unknown")

	// Verify pool still functions normally
	poolIndex, buffer := memoryPool.acquirePool(3)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire pool with unknown type")
	assertions.Equal(3, len(buffer), "Buffer length should be correct")

	memoryPool.releasePool(poolIndex)
}

// TestMemoryPoolIsAvailableEdgeCases tests edge cases in isAvailable function
func TestMemoryPoolIsAvailableEdgeCases(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[string](3, 5, false)

	// Test isAvailable with NotAvailable index
	isAvailableForInvalidIndex := memoryPool.isAvailable(NotAvailable)
	assertions.True(isAvailableForInvalidIndex, "NotAvailable index should be considered available")

	// Test normal available index
	isAvailableForValidIndex := memoryPool.isAvailable(0)
	assertions.True(isAvailableForValidIndex, "Index 0 should be available initially")

	// Acquire a pool and test unavailable
	poolIndex, _ := memoryPool.acquirePool(3)
	isAvailableAfterAcquire := memoryPool.isAvailable(poolIndex)
	assertions.False(isAvailableAfterAcquire, "Index should not be available after acquire")

	memoryPool.releasePool(poolIndex)
}

// TestMemoryPoolReleasePoolLeakDetection tests leak detection in releasePool
func TestMemoryPoolReleasePoolLeakDetection(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[string](3, 5, false)

	// Test releasing an already available pool (should trigger leak detection)
	// This should print a leak detection message but not cause a test failure
	poolIndex := 0
	memoryPool.releasePool(poolIndex)

	// Verify pool is still functional
	newPoolIndex, buffer := memoryPool.acquirePool(3)
	assertions.NotEqual(NotAvailable, newPoolIndex, "Pool should still be functional after leak detection")
	assertions.Equal(3, len(buffer), "Buffer length should be correct")

	memoryPool.releasePool(newPoolIndex)
}

// TestMemoryPoolShrinkPoolWhileInUse tests shrinking when pools are in use
func TestMemoryPoolShrinkPoolWhileInUse(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[string](3, 5, true)

	// Acquire pool and expand it
	poolIndex, _ := memoryPool.acquirePool(10)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to acquire expandable pool")

	// Try to shrink while pool is in use - should not affect in-use pools
	memoryPool.shrinkPool()

	// Pool should still be in use
	assertions.False(memoryPool.isAvailable(poolIndex), "Pool should still be in use after shrink attempt")

	// Release and then shrink
	memoryPool.releasePool(poolIndex)
	memoryPool.shrinkPool()

	// Verify pool works normally after shrink
	poolIndex2, buffer := memoryPool.acquirePool(3)
	assertions.NotEqual(NotAvailable, poolIndex2, "Pool should work after shrink")
	assertions.Equal(3, len(buffer), "Buffer should have correct length after shrink")

	memoryPool.releasePool(poolIndex2)
}

// TestMemoryPoolGetSizeEdgeCases tests edge cases in getSize function
func TestMemoryPoolGetSizeEdgeCases(t *testing.T) {
	assertions := assert.New(t)

	// Test non-expandable pool with size larger than poolLength
	nonExpandablePool := newMemoryPool[string](3, 5, false)

	// This should trigger the warning and return poolLength
	size := nonExpandablePool.getSize(10)
	assertions.Equal(5, size, "Size should be clamped to poolLength for non-expandable pool")

	// Test with NotAvailable size
	defaultSize := nonExpandablePool.getSize(NotAvailable)
	assertions.Equal(5, defaultSize, "NotAvailable size should return poolLength")

	// Test expandable pool with large size
	expandablePool := newMemoryPool[string](3, 5, true)
	largeSize := expandablePool.getSize(15)
	assertions.Equal(15, largeSize, "Expandable pool should allow large sizes")
}

// TestMemoryPoolAllocatePoolEdgeCases tests edge cases in allocatePool function
func TestMemoryPoolAllocatePoolEdgeCases(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[string](3, 5, true)

	// Test initial allocation larger than poolLength
	poolIndex, buffer := memoryPool.acquirePool(10)
	assertions.NotEqual(NotAvailable, poolIndex, "Should be able to allocate large initial buffer")
	assertions.Equal(10, len(buffer), "Initial buffer should be expanded to requested size")

	memoryPool.releasePool(poolIndex)

	// Test expanding existing buffer
	poolIndex2, _ := memoryPool.acquirePool(3)
	assertions.NotEqual(NotAvailable, poolIndex2, "Should be able to acquire pool")

	// Expand the buffer
	expandedBuffer := memoryPool.expandPool(poolIndex2, 8)
	assertions.Equal(8, len(expandedBuffer), "Buffer should be expanded to requested size")

	memoryPool.releasePool(poolIndex2)

	// Re-acquire and test expansion beyond poolLength
	poolIndex3, _ := memoryPool.acquirePool(12)
	assertions.NotEqual(NotAvailable, poolIndex3, "Should be able to acquire expanded pool")

	memoryPool.releasePool(poolIndex3)

	// Test shrinking buffer (requesting smaller size after expansion)
	poolIndex4, _ := memoryPool.acquirePool(2)
	retrievedBuffer := memoryPool.getPool(poolIndex4)
	assertions.Equal(2, len(retrievedBuffer), "Buffer should be shrunk to requested size")

	memoryPool.releasePool(poolIndex4)
}

// TestMemoryPoolBufferReuse tests buffer reuse and zero-value initialization
func TestMemoryPoolBufferReuse(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[string](3, 5, false)

	// Acquire pool and modify data
	poolIndex, buffer := memoryPool.acquirePool(3)
	buffer[0] = "test data"
	buffer[1] = "more data"

	memoryPool.releasePool(poolIndex)

	// Re-acquire same pool and verify it's been zeroed
	poolIndex2, buffer2 := memoryPool.acquirePool(3)
	assertions.Equal("", buffer2[0], "Buffer should be zeroed on reuse")
	assertions.Equal("", buffer2[1], "Buffer should be zeroed on reuse")

	memoryPool.releasePool(poolIndex2)
}

// TestMemoryPoolSequentialAcquireRelease tests sequential operations without leaks
func TestMemoryPoolSequentialAcquireRelease(t *testing.T) {
	assertions := assert.New(t)

	memoryPool := newMemoryPool[string](5, 10, false)

	// Perform many sequential acquire/release operations
	for i := 0; i < 1000; i++ {
		poolIndex, buffer := memoryPool.acquirePool(5)
		if poolIndex != NotAvailable {
			assertions.Equal(5, len(buffer), "Buffer length should be consistent")
			buffer[0] = "test" // Use the buffer
			memoryPool.releasePool(poolIndex)
		}
	}

	// Verify pool is still functional
	finalIndex, finalBuffer := memoryPool.acquirePool(3)
	assertions.NotEqual(NotAvailable, finalIndex, "Pool should still work after many operations")
	assertions.Equal(3, len(finalBuffer), "Buffer length should be correct")
	assertions.Equal("", finalBuffer[0], "Buffer should be zeroed")

	memoryPool.releasePool(finalIndex)
}

// TestMemoryPoolTypeSpecificBehavior tests type-specific behavior
func TestMemoryPoolTypeSpecificBehavior(t *testing.T) {
	assertions := assert.New(t)

	// Test different data types have correct type identifiers
	stringPool := newMemoryPool[string](3, 5, false)
	assertions.Equal("string", strings.ToLower(stringPool.dataType), "string pool should have correct type")

	int64Pool := newMemoryPool[int64](3, 5, false)
	assertions.Equal("int64", strings.ToLower(int64Pool.dataType), "int64 pool should have correct type")

	float64Pool := newMemoryPool[float64](3, 5, false)
	assertions.Equal("float64", strings.ToLower(float64Pool.dataType), "float64 pool should have correct type")

	bytePool := newMemoryPool[byte](3, 5, false)
	assertions.Equal("byte", strings.ToLower(bytePool.dataType), "byte pool should have correct type")

	// Test that each type initializes to its zero value correctly
	stringIndex, stringBuffer := stringPool.acquirePool(2)
	int64Index, int64Buffer := int64Pool.acquirePool(2)
	float64Index, float64Buffer := float64Pool.acquirePool(2)
	byteIndex, byteBuffer := bytePool.acquirePool(2)

	assertions.Equal("", stringBuffer[0], "string should initialize to empty string")
	assertions.Equal(int64(0), int64Buffer[0], "int64 should initialize to 0")
	assertions.Equal(float64(0), float64Buffer[0], "float64 should initialize to 0")
	assertions.Equal(byte(0), byteBuffer[0], "byte should initialize to 0")

	stringPool.releasePool(stringIndex)
	int64Pool.releasePool(int64Index)
	float64Pool.releasePool(float64Index)
	bytePool.releasePool(byteIndex)
}
