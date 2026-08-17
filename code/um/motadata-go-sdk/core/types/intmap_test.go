package types

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestNewIntMap(t *testing.T) {
	t.Run("creates empty map", func(t *testing.T) {
		m := NewIntMap[int, any]()

		if m == nil {
			t.Error("expected map to be created")
		}
		if m.Len() != 0 {
			t.Errorf("expected length 0, got %d", m.Len())
		}
	})

	t.Run("creates with capacity", func(t *testing.T) {
		m := NewIntMapWithCapacity[int, any](100)

		if m == nil {
			t.Error("expected map to be created")
		}
		if m.Len() != 0 {
			t.Errorf("expected length 0, got %d", m.Len())
		}
	})

	t.Run("creates independent instances", func(t *testing.T) {
		m1 := NewIntMap[int, any]()
		m2 := NewIntMap[int, any]()

		m1.Set(1, "value1")
		m2.Set(1, "value2")

		v1, _ := m1.Get(1)
		v2, _ := m2.Get(1)

		if v1 == v2 {
			t.Error("expected independent maps")
		}
	})
}

func TestIntMapSet(t *testing.T) {
	t.Run("sets single value", func(t *testing.T) {
		m := NewIntMap[int, any]()

		m.Set(1, "value")

		if m.Len() != 1 {
			t.Errorf("expected length 1, got %d", m.Len())
		}
	})

	t.Run("sets multiple values", func(t *testing.T) {
		m := NewIntMap[int, any]()

		m.Set(1, "value1")
		m.Set(2, "value2")
		m.Set(3, "value3")

		if m.Len() != 3 {
			t.Errorf("expected length 3, got %d", m.Len())
		}
	})

	t.Run("overwrites existing key", func(t *testing.T) {
		m := NewIntMap[int, any]()

		m.Set(1, "value1")
		m.Set(1, "value2")

		if m.Len() != 1 {
			t.Errorf("expected length 1, got %d", m.Len())
		}

		value, exists := m.Get(1)
		if !exists {
			t.Error("expected key to exist")
		}
		if value != "value2" {
			t.Errorf("expected 'value2', got %v", value)
		}
	})

	t.Run("sets various types", func(t *testing.T) {
		m := NewIntMap[int, any]()

		m.Set(1, "hello")
		m.Set(2, 42)
		m.Set(3, 3.14)
		m.Set(4, true)
		m.Set(5, nil)
		m.Set(6, []int{1, 2, 3})
		m.Set(7, map[string]int{"a": 1})

		if m.Len() != 7 {
			t.Errorf("expected length 7, got %d", m.Len())
		}
	})

	t.Run("sets zero key", func(t *testing.T) {
		m := NewIntMap[int, any]()

		m.Set(0, "zero key value")

		value, exists := m.Get(0)
		if !exists {
			t.Error("expected zero key to exist")
		}
		if value != "zero key value" {
			t.Errorf("expected 'zero key value', got %v", value)
		}
	})

	t.Run("sets negative keys", func(t *testing.T) {
		m := NewIntMap[int, any]()

		m.Set(-1, "negative one")
		m.Set(-100, "negative hundred")

		v1, exists1 := m.Get(-1)
		v2, exists2 := m.Get(-100)

		if !exists1 || v1 != "negative one" {
			t.Error("expected negative key -1 to exist with correct value")
		}
		if !exists2 || v2 != "negative hundred" {
			t.Error("expected negative key -100 to exist with correct value")
		}
	})
}

func TestIntMapGet(t *testing.T) {
	t.Run("gets existing value", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, "value")

		value, exists := m.Get(1)

		if !exists {
			t.Error("expected key to exist")
		}
		if value != "value" {
			t.Errorf("expected 'value', got %v", value)
		}
	})

	t.Run("returns false for non-existing key", func(t *testing.T) {
		m := NewIntMap[int, any]()

		_, exists := m.Get(999)

		if exists {
			t.Error("expected key to not exist")
		}
	})

	t.Run("returns nil value correctly", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, nil)

		value, exists := m.Get(1)

		if !exists {
			t.Error("expected key to exist")
		}
		if value != nil {
			t.Errorf("expected nil, got %v", value)
		}
	})

	t.Run("distinguishes nil value from missing key", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, nil)

		_, existsNil := m.Get(1)
		_, existsMissing := m.Get(999)

		if !existsNil {
			t.Error("expected nil key to exist")
		}
		if existsMissing {
			t.Error("expected missing key to not exist")
		}
	})
}

func TestIntMapDelete(t *testing.T) {
	t.Run("deletes existing key", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, "value")

		deleted := m.Delete(1)

		if !deleted {
			t.Error("expected Delete to return true for existing key")
		}
		if m.Has(1) {
			t.Error("expected key to be deleted")
		}
		if m.Len() != 0 {
			t.Errorf("expected length 0, got %d", m.Len())
		}
	})

	t.Run("returns false for non-existing key", func(t *testing.T) {
		m := NewIntMap[int, any]()

		deleted := m.Delete(999)

		if deleted {
			t.Error("expected Delete to return false for non-existing key")
		}
		if m.Len() != 0 {
			t.Errorf("expected length 0, got %d", m.Len())
		}
	})

	t.Run("deletes only specified key", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, "value1")
		m.Set(2, "value2")
		m.Set(3, "value3")

		m.Delete(2)

		if m.Len() != 2 {
			t.Errorf("expected length 2, got %d", m.Len())
		}
		if !m.Has(1) {
			t.Error("expected key 1 to exist")
		}
		if m.Has(2) {
			t.Error("expected key 2 to be deleted")
		}
		if !m.Has(3) {
			t.Error("expected key 3 to exist")
		}
	})
}

func TestIntMapHas(t *testing.T) {
	t.Run("returns true for existing key", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, "value")

		if !m.Has(1) {
			t.Error("expected Has to return true")
		}
	})

	t.Run("returns false for non-existing key", func(t *testing.T) {
		m := NewIntMap[int, any]()

		if m.Has(999) {
			t.Error("expected Has to return false")
		}
	})

	t.Run("returns true for key with nil value", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, nil)

		if !m.Has(1) {
			t.Error("expected Has to return true for nil value")
		}
	})

	t.Run("returns false after deletion", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, "value")
		m.Delete(1)

		if m.Has(1) {
			t.Error("expected Has to return false after deletion")
		}
	})
}

func TestIntMapLen(t *testing.T) {
	t.Run("returns 0 for empty map", func(t *testing.T) {
		m := NewIntMap[int, any]()

		if m.Len() != 0 {
			t.Errorf("expected length 0, got %d", m.Len())
		}
	})

	t.Run("returns correct count after additions", func(t *testing.T) {
		m := NewIntMap[int, any]()

		for i := 0; i < 10; i++ {
			m.Set(i, i)
		}

		if m.Len() != 10 {
			t.Errorf("expected length 10, got %d", m.Len())
		}
	})

	t.Run("returns correct count after deletions", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, 1)
		m.Set(2, 2)
		m.Set(3, 3)
		m.Delete(2)

		if m.Len() != 2 {
			t.Errorf("expected length 2, got %d", m.Len())
		}
	})

	t.Run("does not increase on overwrite", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, "value1")
		m.Set(1, "value2")
		m.Set(1, "value3")

		if m.Len() != 1 {
			t.Errorf("expected length 1, got %d", m.Len())
		}
	})
}

func TestIntMapKeys(t *testing.T) {
	t.Run("returns empty slice for empty map", func(t *testing.T) {
		m := NewIntMap[int, any]()

		keys := m.Keys()

		if len(keys) != 0 {
			t.Errorf("expected 0 keys, got %d", len(keys))
		}
	})

	t.Run("returns all keys", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, "a")
		m.Set(2, "b")
		m.Set(3, "c")

		keys := m.Keys()

		if len(keys) != 3 {
			t.Errorf("expected 3 keys, got %d", len(keys))
		}

		keySet := make(map[int]bool)
		for _, k := range keys {
			keySet[k] = true
		}

		if !keySet[1] || !keySet[2] || !keySet[3] {
			t.Error("expected all keys to be present")
		}
	})

	t.Run("returns copy of keys", func(t *testing.T) {
		m := NewIntMap[int, any]()
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
		if !hasOriginal {
			t.Error("modifying returned keys should not affect map")
		}
	})
}

func TestIntMapValues(t *testing.T) {
	t.Run("returns empty slice for empty map", func(t *testing.T) {
		m := NewIntMap[int, any]()

		values := m.Values()

		if len(values) != 0 {
			t.Errorf("expected 0 values, got %d", len(values))
		}
	})

	t.Run("returns all values", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, 10)
		m.Set(2, 20)
		m.Set(3, 30)

		values := m.Values()

		if len(values) != 3 {
			t.Errorf("expected 3 values, got %d", len(values))
		}

		sum := 0
		for _, v := range values {
			sum += v.(int)
		}
		if sum != 60 {
			t.Errorf("expected sum 60, got %d", sum)
		}
	})

	t.Run("includes nil values", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, nil)
		m.Set(2, "value")

		values := m.Values()

		if len(values) != 2 {
			t.Errorf("expected 2 values, got %d", len(values))
		}

		hasNil := false
		for _, v := range values {
			if v == nil {
				hasNil = true
				break
			}
		}
		if !hasNil {
			t.Error("expected nil value to be included")
		}
	})
}

func TestIntMapClear(t *testing.T) {
	t.Run("clears all entries", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, 1)
		m.Set(2, 2)
		m.Set(3, 3)

		m.Clear()

		if m.Len() != 0 {
			t.Errorf("expected length 0 after clear, got %d", m.Len())
		}
	})

	t.Run("clear on empty map does not panic", func(t *testing.T) {
		m := NewIntMap[int, any]()

		m.Clear()

		if m.Len() != 0 {
			t.Errorf("expected length 0, got %d", m.Len())
		}
	})

	t.Run("map is usable after clear", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, "value")
		m.Clear()

		m.Set(2, "newvalue")

		if m.Len() != 1 {
			t.Errorf("expected length 1, got %d", m.Len())
		}
		value, exists := m.Get(2)
		if !exists || value != "newvalue" {
			t.Error("expected to get 'newvalue'")
		}
	})

	t.Run("old keys not accessible after clear", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, "oldvalue")
		m.Clear()

		if m.Has(1) {
			t.Error("expected old key to not exist after clear")
		}
	})
}

func TestIntMapRange(t *testing.T) {
	t.Run("iterates over all entries", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, 1)
		m.Set(2, 2)
		m.Set(3, 3)

		count := 0
		m.Range(func(key int, value any) bool {
			count++
			return true
		})

		if count != 3 {
			t.Errorf("expected 3 iterations, got %d", count)
		}
	})

	t.Run("stops iteration when function returns false", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, 1)
		m.Set(2, 2)
		m.Set(3, 3)

		count := 0
		m.Range(func(key int, value any) bool {
			count++
			return count < 2
		})

		if count != 2 {
			t.Errorf("expected 2 iterations, got %d", count)
		}
	})

	t.Run("handles empty map", func(t *testing.T) {
		m := NewIntMap[int, any]()

		count := 0
		m.Range(func(key int, value any) bool {
			count++
			return true
		})

		if count != 0 {
			t.Errorf("expected 0 iterations, got %d", count)
		}
	})

	t.Run("provides correct key-value pairs", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, 10)
		m.Set(2, 20)

		collected := make(map[int]int)
		m.Range(func(key int, value any) bool {
			collected[key] = value.(int)
			return true
		})

		if collected[1] != 10 || collected[2] != 20 {
			t.Error("expected correct key-value pairs")
		}
	})

	t.Run("stops immediately on first false return", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, 1)
		m.Set(2, 2)
		m.Set(3, 3)

		count := 0
		m.Range(func(key int, value any) bool {
			count++
			return false
		})

		if count != 1 {
			t.Errorf("expected 1 iteration, got %d", count)
		}
	})
}

func TestIntMapCopy(t *testing.T) {
	t.Run("copies empty map", func(t *testing.T) {
		m := NewIntMap[int, any]()

		copied := m.Copy()

		if copied == nil {
			t.Error("expected copied map to be created")
		}
		if copied.Len() != 0 {
			t.Errorf("expected length 0, got %d", copied.Len())
		}
	})

	t.Run("copies all entries", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, 10)
		m.Set(2, 20)
		m.Set(3, 30)

		copied := m.Copy()

		if copied.Len() != 3 {
			t.Errorf("expected length 3, got %d", copied.Len())
		}

		for _, key := range []int{1, 2, 3} {
			origVal, _ := m.Get(key)
			copyVal, exists := copied.Get(key)
			if !exists {
				t.Errorf("expected key %d to exist in copy", key)
			}
			if origVal != copyVal {
				t.Errorf("expected value %v, got %v", origVal, copyVal)
			}
		}
	})

	t.Run("creates independent copy", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, "original")

		copied := m.Copy()
		copied.Set(1, "modified")

		origVal, _ := m.Get(1)
		if origVal != "original" {
			t.Error("modifying copy should not affect original")
		}

		copyVal, _ := copied.Get(1)
		if copyVal != "modified" {
			t.Error("copy should have modified value")
		}
	})

	t.Run("adding to copy does not affect original", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, 1)

		copied := m.Copy()
		copied.Set(2, 2)

		if m.Len() != 1 {
			t.Errorf("original length should be 1, got %d", m.Len())
		}
		if copied.Len() != 2 {
			t.Errorf("copy length should be 2, got %d", copied.Len())
		}
	})

	t.Run("deleting from copy does not affect original", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, 1)
		m.Set(2, 2)

		copied := m.Copy()
		copied.Delete(1)

		if m.Len() != 2 {
			t.Errorf("original length should be 2, got %d", m.Len())
		}
		if !m.Has(1) {
			t.Error("original should still have key 1")
		}
		if copied.Len() != 1 {
			t.Errorf("copy length should be 1, got %d", copied.Len())
		}
	})

	t.Run("clearing copy does not affect original", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, 1)
		m.Set(2, 2)

		copied := m.Copy()
		copied.Clear()

		if m.Len() != 2 {
			t.Errorf("original length should be 2, got %d", m.Len())
		}
		if copied.Len() != 0 {
			t.Errorf("copy length should be 0, got %d", copied.Len())
		}
	})

	t.Run("copies nil values correctly", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, nil)

		copied := m.Copy()

		value, exists := copied.Get(1)
		if !exists {
			t.Error("expected key 1 to exist in copy")
		}
		if value != nil {
			t.Errorf("expected nil value, got %v", value)
		}
	})
}

func TestIntMapPutIfNotExists(t *testing.T) {
	t.Run("adds new key", func(t *testing.T) {
		m := NewIntMap[int, any]()

		val, wasNew := m.PutIfNotExists(1, "value")

		// wasNew is true if the key was newly inserted
		if !wasNew {
			t.Error("expected key to be newly inserted")
		}
		if val != "value" {
			t.Errorf("expected 'value', got %v", val)
		}
		if m.Len() != 1 {
			t.Errorf("expected length 1, got %d", m.Len())
		}
	})

	t.Run("does not overwrite existing key", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, "original")

		val, wasNew := m.PutIfNotExists(1, "new")

		// wasNew is false if key already existed
		if wasNew {
			t.Error("expected key to already exist")
		}
		if val != "original" {
			t.Errorf("expected 'original', got %v", val)
		}

		storedVal, _ := m.Get(1)
		if storedVal != "original" {
			t.Errorf("expected stored value 'original', got %v", storedVal)
		}
	})
}

func TestIntMapOperationsSequence(t *testing.T) {
	t.Run("complex operations sequence", func(t *testing.T) {
		m := NewIntMap[int, any]()

		// Add entries
		m.Set(1, 10)
		m.Set(2, 20)
		m.Set(3, 30)

		if m.Len() != 3 {
			t.Errorf("expected length 3, got %d", m.Len())
		}

		// Update entry
		m.Set(2, 200)
		value, _ := m.Get(2)
		if value != 200 {
			t.Errorf("expected 200, got %v", value)
		}

		// Delete entry
		m.Delete(1)
		if m.Has(1) {
			t.Error("expected key 1 to be deleted")
		}

		// Check remaining
		if m.Len() != 2 {
			t.Errorf("expected length 2, got %d", m.Len())
		}

		// Clear and restart
		m.Clear()
		m.Set(100, "new entry")

		if m.Len() != 1 {
			t.Errorf("expected length 1, got %d", m.Len())
		}
	})
}

func TestIntMapConcurrentExecution(t *testing.T) {
	t.Run("handles concurrent reads on same map", func(t *testing.T) {
		m := NewIntMap[int, any]()
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

		expectedReads := int64(numGoroutines * numReads)
		if successCount.Load() != expectedReads {
			t.Errorf("expected %d successful reads, got %d", expectedReads, successCount.Load())
		}
	})

	t.Run("concurrent Len calls are safe", func(t *testing.T) {
		m := NewIntMap[int, any]()
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

		expectedReads := int64(numGoroutines * numReads)
		if correctCount.Load() != expectedReads {
			t.Errorf("expected %d correct Len reads, got %d", expectedReads, correctCount.Load())
		}
	})

	t.Run("concurrent Has calls are safe", func(t *testing.T) {
		m := NewIntMap[int, any]()
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

		expectedChecks := int64(numGoroutines * numChecks)
		if existsCount.Load() != expectedChecks {
			t.Errorf("expected %d exists checks, got %d", expectedChecks, existsCount.Load())
		}
		if notExistsCount.Load() != expectedChecks {
			t.Errorf("expected %d notexists checks, got %d", expectedChecks, notExistsCount.Load())
		}
	})

	t.Run("concurrent Range calls are safe", func(t *testing.T) {
		m := NewIntMap[int, any]()
		m.Set(1, 1)
		m.Set(2, 2)
		m.Set(3, 3)

		var wg sync.WaitGroup
		const numGoroutines = 50
		const numRanges = 50

		var correctRanges atomic.Int64

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
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

		wg.Wait()

		expectedRanges := int64(numGoroutines * numRanges)
		if correctRanges.Load() != expectedRanges {
			t.Errorf("expected %d correct Range calls, got %d", expectedRanges, correctRanges.Load())
		}
	})
}

func TestIntMapEdgeCases(t *testing.T) {
	t.Run("handles large keys", func(t *testing.T) {
		m := NewIntMap[int, any]()

		largeKeys := []int{
			1000000,
			-1000000,
			2147483647,  // max int32
			-2147483648, // min int32
		}

		for i, key := range largeKeys {
			m.Set(key, i)
		}

		for i, key := range largeKeys {
			value, exists := m.Get(key)
			if !exists {
				t.Errorf("expected key %d to exist", key)
			}
			if value != i {
				t.Errorf("expected value %d for key %d, got %v", i, key, value)
			}
		}
	})

	t.Run("handles large number of entries", func(t *testing.T) {
		m := NewIntMapWithCapacity[int, any](10000)
		const numEntries = 10000

		for i := 0; i < numEntries; i++ {
			m.Set(i, i)
		}

		if m.Len() != numEntries {
			t.Errorf("expected length %d, got %d", numEntries, m.Len())
		}
	})

	t.Run("handles complex value types", func(t *testing.T) {
		m := NewIntMap[int, any]()

		type customStruct struct {
			Name  string
			Value int
		}

		m.Set(1, customStruct{Name: "test", Value: 42})
		m.Set(2, &customStruct{Name: "ptr", Value: 100})
		m.Set(3, func() int { return 1 })
		m.Set(4, make(chan int))

		if m.Len() != 4 {
			t.Errorf("expected length 4, got %d", m.Len())
		}

		structVal, _ := m.Get(1)
		if s, ok := structVal.(customStruct); !ok || s.Name != "test" {
			t.Error("expected to get struct value")
		}
	})

	t.Run("delete then re-add same key", func(t *testing.T) {
		m := NewIntMap[int, any]()

		m.Set(1, "first")
		m.Delete(1)
		m.Set(1, "second")

		value, exists := m.Get(1)
		if !exists {
			t.Error("expected key to exist after re-add")
		}
		if value != "second" {
			t.Errorf("expected 'second', got %v", value)
		}
	})
}
