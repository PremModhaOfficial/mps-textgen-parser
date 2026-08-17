package types

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIntMap(t *testing.T) {
	t.Run("creates empty map", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()

		assertions.NotNil(m, "expected map to be created")
		assertions.Equal(0, m.Len())
	})

	t.Run("creates with capacity", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMapWithCapacity(100)

		assertions.NotNil(m, "expected map to be created")
		assertions.Equal(0, m.Len())
	})

	t.Run("creates independent instances", func(t *testing.T) {
		assertions := assert.New(t)
		m1 := NewIntMap()
		m2 := NewIntMap()

		m1.Set(1, "value1")
		m2.Set(1, "value2")

		v1, _ := m1.Get(1)
		v2, _ := m2.Get(1)

		assertions.NotEqual(v1, v2, "expected independent maps")
	})
}

func TestIntMapSet(t *testing.T) {
	t.Run("sets single value", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()

		m.Set(1, "value")

		assertions.Equal(1, m.Len())
	})

	t.Run("sets multiple values", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()

		m.Set(1, "value1")
		m.Set(2, "value2")
		m.Set(3, "value3")

		assertions.Equal(3, m.Len())
	})

	t.Run("overwrites existing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()

		m.Set(1, "value1")
		m.Set(1, "value2")

		assertions.Equal(1, m.Len())

		value, exists := m.Get(1)
		assertions.True(exists, "expected key to exist")
		assertions.Equal("value2", value)
	})

	t.Run("sets various types", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()

		m.Set(1, "hello")
		m.Set(2, 42)
		m.Set(3, 3.14)
		m.Set(4, true)
		m.Set(5, nil)
		m.Set(6, []int{1, 2, 3})
		m.Set(7, map[string]int{"a": 1})

		assertions.Equal(7, m.Len())
	})

	t.Run("sets zero key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()

		m.Set(0, "zero key value")

		value, exists := m.Get(0)
		assertions.True(exists, "expected zero key to exist")
		assertions.Equal("zero key value", value)
	})

	t.Run("sets negative keys", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()

		m.Set(-1, "negative one")
		m.Set(-100, "negative hundred")

		v1, exists1 := m.Get(-1)
		v2, exists2 := m.Get(-100)

		assertions.True(exists1 && v1 == "negative one", "expected negative key -1 to exist with correct value")
		assertions.True(exists2 && v2 == "negative hundred", "expected negative key -100 to exist with correct value")
	})
}

func TestIntMapGet(t *testing.T) {
	t.Run("gets existing value", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, "value")

		value, exists := m.Get(1)

		assertions.True(exists, "expected key to exist")
		assertions.Equal("value", value)
	})

	t.Run("returns false for non-existing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()

		_, exists := m.Get(999)

		assertions.False(exists, "expected key to not exist")
	})

	t.Run("returns nil value correctly", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, nil)

		value, exists := m.Get(1)

		assertions.True(exists, "expected key to exist")
		assertions.Nil(value)
	})

	t.Run("distinguishes nil value from missing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, nil)

		_, existsNil := m.Get(1)
		_, existsMissing := m.Get(999)

		assertions.True(existsNil, "expected nil key to exist")
		assertions.False(existsMissing, "expected missing key to not exist")
	})
}

func TestIntMapDelete(t *testing.T) {
	t.Run("deletes existing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, "value")

		deleted := m.Delete(1)

		assertions.True(deleted, "expected Delete to return true for existing key")
		assertions.False(m.Has(1), "expected key to be deleted")
		assertions.Equal(0, m.Len())
	})

	t.Run("returns false for non-existing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()

		deleted := m.Delete(999)

		assertions.False(deleted, "expected Delete to return false for non-existing key")
		assertions.Equal(0, m.Len())
	})

	t.Run("deletes only specified key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, "value1")
		m.Set(2, "value2")
		m.Set(3, "value3")

		m.Delete(2)

		assertions.Equal(2, m.Len())
		assertions.True(m.Has(1), "expected key 1 to exist")
		assertions.False(m.Has(2), "expected key 2 to be deleted")
		assertions.True(m.Has(3), "expected key 3 to exist")
	})
}

func TestIntMapHas(t *testing.T) {
	t.Run("returns true for existing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, "value")

		assertions.True(m.Has(1), "expected Has to return true")
	})

	t.Run("returns false for non-existing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()

		assertions.False(m.Has(999), "expected Has to return false")
	})

	t.Run("returns true for key with nil value", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, nil)

		assertions.True(m.Has(1), "expected Has to return true for nil value")
	})

	t.Run("returns false after deletion", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, "value")
		m.Delete(1)

		assertions.False(m.Has(1), "expected Has to return false after deletion")
	})
}

func TestIntMapLen(t *testing.T) {
	t.Run("returns 0 for empty map", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()

		assertions.Equal(0, m.Len())
	})

	t.Run("returns correct count after additions", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()

		for i := 0; i < 10; i++ {
			m.Set(i, i)
		}

		assertions.Equal(10, m.Len())
	})

	t.Run("returns correct count after deletions", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, 1)
		m.Set(2, 2)
		m.Set(3, 3)
		m.Delete(2)

		assertions.Equal(2, m.Len())
	})

	t.Run("does not increase on overwrite", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, "value1")
		m.Set(1, "value2")
		m.Set(1, "value3")

		assertions.Equal(1, m.Len())
	})
}

func TestIntMapKeys(t *testing.T) {
	t.Run("returns empty slice for empty map", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()

		keys := m.Keys()

		assertions.Equal(0, len(keys))
	})

	t.Run("returns all keys", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, "a")
		m.Set(2, "b")
		m.Set(3, "c")

		keys := m.Keys()

		assertions.Equal(3, len(keys))

		keySet := make(map[int]bool)
		for _, k := range keys {
			keySet[k] = true
		}

		assertions.True(keySet[1] && keySet[2] && keySet[3], "expected all keys to be present")
	})

	t.Run("returns copy of keys", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, "value")

		keys := m.Keys()
		keys[0] = 999

		originalKeys := m.Keys()
		hasOriginal := false
		for _, k := range originalKeys {
			if k == 1 {
				hasOriginal = true
				break
			}
		}
		assertions.True(hasOriginal, "modifying returned keys should not affect map")
	})
}

func TestIntMapValues(t *testing.T) {
	t.Run("returns empty slice for empty map", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()

		values := m.Values()

		assertions.Equal(0, len(values))
	})

	t.Run("returns all values", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, 10)
		m.Set(2, 20)
		m.Set(3, 30)

		values := m.Values()

		assertions.Equal(3, len(values))

		sum := 0
		for _, v := range values {
			sum += v.(int)
		}
		assertions.Equal(60, sum)
	})

	t.Run("includes nil values", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, nil)
		m.Set(2, "value")

		values := m.Values()

		assertions.Equal(2, len(values))

		hasNil := false
		for _, v := range values {
			if v == nil {
				hasNil = true
				break
			}
		}
		assertions.True(hasNil, "expected nil value to be included")
	})
}

func TestIntMapClear(t *testing.T) {
	t.Run("clears all entries", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, 1)
		m.Set(2, 2)
		m.Set(3, 3)

		m.Clear()

		assertions.Equal(0, m.Len())
	})

	t.Run("clear on empty map does not panic", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()

		m.Clear()

		assertions.Equal(0, m.Len())
	})

	t.Run("map is usable after clear", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, "value")
		m.Clear()

		m.Set(2, "newvalue")

		assertions.Equal(1, m.Len())
		value, exists := m.Get(2)
		assertions.True(exists && value == "newvalue", "expected to get 'newvalue'")
	})

	t.Run("old keys not accessible after clear", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, "oldvalue")
		m.Clear()

		assertions.False(m.Has(1), "expected old key to not exist after clear")
	})
}

func TestIntMapRange(t *testing.T) {
	t.Run("iterates over all entries", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, 1)
		m.Set(2, 2)
		m.Set(3, 3)

		count := 0
		m.Range(func(key int, value any) bool {
			count++
			return true
		})

		assertions.Equal(3, count)
	})

	t.Run("stops iteration when function returns false", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, 1)
		m.Set(2, 2)
		m.Set(3, 3)

		count := 0
		m.Range(func(key int, value any) bool {
			count++
			return count < 2
		})

		assertions.Equal(2, count)
	})

	t.Run("handles empty map", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()

		count := 0
		m.Range(func(key int, value any) bool {
			count++
			return true
		})

		assertions.Equal(0, count)
	})

	t.Run("provides correct key-value pairs", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, 10)
		m.Set(2, 20)

		collected := make(map[int]int)
		m.Range(func(key int, value any) bool {
			collected[key] = value.(int)
			return true
		})

		assertions.Equal(10, collected[1])
		assertions.Equal(20, collected[2])
	})

	t.Run("stops immediately on first false return", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, 1)
		m.Set(2, 2)
		m.Set(3, 3)

		count := 0
		m.Range(func(key int, value any) bool {
			count++
			return false
		})

		assertions.Equal(1, count)
	})
}

func TestIntMapCopy(t *testing.T) {
	t.Run("copies empty map", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()

		copied := m.Copy()

		assertions.NotNil(copied, "expected copied map to be created")
		assertions.Equal(0, copied.Len())
	})

	t.Run("copies all entries", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, 10)
		m.Set(2, 20)
		m.Set(3, 30)

		copied := m.Copy()

		assertions.Equal(3, copied.Len())

		for _, key := range []int{1, 2, 3} {
			origVal, _ := m.Get(key)
			copyVal, exists := copied.Get(key)
			assertions.True(exists, "expected key %d to exist in copy", key)
			assertions.Equal(origVal, copyVal)
		}
	})

	t.Run("creates independent copy", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, "original")

		copied := m.Copy()
		copied.Set(1, "modified")

		origVal, _ := m.Get(1)
		assertions.Equal("original", origVal, "modifying copy should not affect original")

		copyVal, _ := copied.Get(1)
		assertions.Equal("modified", copyVal, "copy should have modified value")
	})

	t.Run("adding to copy does not affect original", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, 1)

		copied := m.Copy()
		copied.Set(2, 2)

		assertions.Equal(1, m.Len(), "original length should be 1")
		assertions.Equal(2, copied.Len(), "copy length should be 2")
	})

	t.Run("deleting from copy does not affect original", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, 1)
		m.Set(2, 2)

		copied := m.Copy()
		copied.Delete(1)

		assertions.Equal(2, m.Len(), "original length should be 2")
		assertions.True(m.Has(1), "original should still have key 1")
		assertions.Equal(1, copied.Len(), "copy length should be 1")
	})

	t.Run("clearing copy does not affect original", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, 1)
		m.Set(2, 2)

		copied := m.Copy()
		copied.Clear()

		assertions.Equal(2, m.Len(), "original length should be 2")
		assertions.Equal(0, copied.Len(), "copy length should be 0")
	})

	t.Run("copies nil values correctly", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, nil)

		copied := m.Copy()

		value, exists := copied.Get(1)
		assertions.True(exists, "expected key 1 to exist in copy")
		assertions.Nil(value)
	})
}

func TestIntMapPutIfNotExists(t *testing.T) {
	t.Run("adds new key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()

		val, wasNew := m.PutIfNotExists(1, "value")

		// wasNew is true if the key was newly inserted
		assertions.True(wasNew, "expected key to be newly inserted")
		assertions.Equal("value", val)
		assertions.Equal(1, m.Len())
	})

	t.Run("does not overwrite existing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewIntMap()
		m.Set(1, "original")

		val, wasNew := m.PutIfNotExists(1, "new")

		// wasNew is false if key already existed
		assertions.False(wasNew, "expected key to already exist")
		assertions.Equal("original", val)

		storedVal, _ := m.Get(1)
		assertions.Equal("original", storedVal)
	})
}

func TestIntMapOperationsSequence(t *testing.T) {
	t.Run("complex operations sequence", func(t *testing.T) {
		m := NewIntMap()

		// Add entries
		m.Set(1, 10)
		m.Set(2, 20)
		m.Set(3, 30)

		assertions := assert.New(t)
		assertions.Equal(3, m.Len())

		// Update entry
		m.Set(2, 200)
		value, _ := m.Get(2)
		assertions.Equal(200, value)

		// Delete entry
		m.Delete(1)
		assertions.False(m.Has(1), "expected key 1 to be deleted")

		// Check remaining
		assertions.Equal(2, m.Len())

		// Clear and restart
		m.Clear()
		m.Set(100, "new entry")

		assertions.Equal(1, m.Len())
	})
}

func TestIntMapConcurrentExecution(t *testing.T) {
	t.Run("handles concurrent reads on same map", func(t *testing.T) {
		m := NewIntMap()
		m.Set(1, "value")

		var wg sync.WaitGroup
		const numGoroutines = 100
		const numReads = 100

		var successCount atomic.Int64

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < numReads; j++ {
					value, exists := m.Get(1)
					if exists && value == "value" {
						successCount.Add(1)
					}
				}
			}()
		}

		wg.Wait()

		assertions := assert.New(t)
		expectedReads := int64(numGoroutines * numReads)
		assertions.Equal(expectedReads, successCount.Load())
	})

	t.Run("concurrent Len calls are safe", func(t *testing.T) {
		m := NewIntMap()
		m.Set(1, 1)
		m.Set(2, 2)

		var wg sync.WaitGroup
		const numGoroutines = 50
		const numReads = 100

		var correctCount atomic.Int64

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < numReads; j++ {
					if m.Len() == 2 {
						correctCount.Add(1)
					}
				}
			}()
		}

		wg.Wait()

		assertions := assert.New(t)
		expectedReads := int64(numGoroutines * numReads)
		assertions.Equal(expectedReads, correctCount.Load())
	})

	t.Run("concurrent Has calls are safe", func(t *testing.T) {
		m := NewIntMap()
		m.Set(1, true)

		var wg sync.WaitGroup
		const numGoroutines = 50
		const numChecks = 100

		var existsCount atomic.Int64
		var notExistsCount atomic.Int64

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < numChecks; j++ {
					if m.Has(1) {
						existsCount.Add(1)
					}
					if !m.Has(999) {
						notExistsCount.Add(1)
					}
				}
			}()
		}

		wg.Wait()

		assertions := assert.New(t)
		expectedChecks := int64(numGoroutines * numChecks)
		assertions.Equal(expectedChecks, existsCount.Load())
		assertions.Equal(expectedChecks, notExistsCount.Load())
	})

	t.Run("concurrent Range calls are safe", func(t *testing.T) {
		m := NewIntMap()
		m.Set(1, 1)
		m.Set(2, 2)
		m.Set(3, 3)

		var waitGroup sync.WaitGroup
		const numGoroutines = 50
		const numRanges = 50

		var correctRanges atomic.Int64

		for i := 0; i < numGoroutines; i++ {
			waitGroup.Add(1)
			go func() {
				defer waitGroup.Done()
				for j := 0; j < numRanges; j++ {
					count := 0
					m.Range(func(key int, value any) bool {
						count++
						return true
					})
					if count == 3 {
						correctRanges.Add(1)
					}
				}
			}()
		}

		waitGroup.Wait()

		assertions := assert.New(t)
		expectedRanges := int64(numGoroutines * numRanges)
		assertions.Equal(expectedRanges, correctRanges.Load())
	})
}

func TestIntMapEdgeCases(t *testing.T) {
	t.Run("handles large keys", func(t *testing.T) {
		m := NewIntMap()

		largeKeys := []int{
			1000000,
			-1000000,
			2147483647,  // max int32
			-2147483648, // min int32
		}

		for i, key := range largeKeys {
			m.Set(key, i)
		}

		assertions := assert.New(t)
		for i, key := range largeKeys {
			value, exists := m.Get(key)
			assertions.True(exists, "expected key %d to exist", key)
			assertions.Equal(i, value, "expected value %d for key %d", i, key)
		}
	})

	t.Run("handles large number of entries", func(t *testing.T) {
		m := NewIntMapWithCapacity(10000)
		const numEntries = 10000

		for i := 0; i < numEntries; i++ {
			m.Set(i, i)
		}

		assertions := assert.New(t)
		assertions.Equal(numEntries, m.Len())
	})

	t.Run("handles complex value types", func(t *testing.T) {
		m := NewIntMap()

		type customStruct struct {
			Name  string
			Value int
		}

		m.Set(1, customStruct{Name: "test", Value: 42})
		m.Set(2, &customStruct{Name: "ptr", Value: 100})
		m.Set(3, func() int { return 1 })
		m.Set(4, make(chan int))

		assertions := assert.New(t)
		assertions.Equal(4, m.Len())

		structVal, _ := m.Get(1)
		if s, ok := structVal.(customStruct); ok {
			assertions.Equal("test", s.Name)
		} else {
			assertions.Fail("expected to get struct value")
		}
	})

	t.Run("delete then re-add same key", func(t *testing.T) {
		m := NewIntMap()

		m.Set(1, "first")
		m.Delete(1)
		m.Set(1, "second")

		assertions := assert.New(t)
		value, exists := m.Get(1)
		assertions.True(exists, "expected key to exist after re-add")
		assertions.Equal("second", value)
	})
}
