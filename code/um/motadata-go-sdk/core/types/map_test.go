package types

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestNewMap(t *testing.T) {
	t.Run("creates empty map", func(t *testing.T) {
		m := NewMap[string, any]()

		if m == nil {
			t.Error("expected map to be created")
		}
		if m.Len() != 0 {
			t.Errorf("expected length 0, got %d", m.Len())
		}
	})

	t.Run("creates independent instances", func(t *testing.T) {
		m1 := NewMap[string, any]()
		m2 := NewMap[string, any]()

		m1.Set("key", "value1")
		m2.Set("key", "value2")

		v1, _ := m1.Get("key")
		v2, _ := m2.Get("key")

		if v1 == v2 {
			t.Error("expected independent maps")
		}
	})
}

func TestMapSet(t *testing.T) {
	t.Run("sets single value", func(t *testing.T) {
		m := NewMapWithCapacity[string, any](1)

		m.Set("key", "value")

		if m.Len() != 1 {
			t.Errorf("expected length 1, got %d", m.Len())
		}
	})

	t.Run("sets multiple values", func(t *testing.T) {
		m := NewMap[string, any]()

		m.Set("key1", "value1")
		m.Set("key2", "value2")
		m.Set("key3", "value3")

		if m.Len() != 3 {
			t.Errorf("expected length 3, got %d", m.Len())
		}
	})

	t.Run("overwrites existing key", func(t *testing.T) {
		m := NewMap[string, any]()

		m.Set("key", "value1")
		m.Set("key", "value2")

		if m.Len() != 1 {
			t.Errorf("expected length 1, got %d", m.Len())
		}

		value, exists := m.Get("key")
		if !exists {
			t.Error("expected key to exist")
		}
		if value != "value2" {
			t.Errorf("expected 'value2', got %v", value)
		}
	})

	t.Run("sets various types", func(t *testing.T) {
		m := NewMap[string, any]()

		m.Set("string", "hello")
		m.Set("int", 42)
		m.Set("float", 3.14)
		m.Set("bool", true)
		m.Set("nil", nil)
		m.Set("slice", []int{1, 2, 3})
		m.Set("map", map[string]int{"a": 1})

		if m.Len() != 7 {
			t.Errorf("expected length 7, got %d", m.Len())
		}
	})

	t.Run("sets empty string key", func(t *testing.T) {
		m := NewMap[string, any]()

		m.Set("", "empty key value")

		value, exists := m.Get("")
		if !exists {
			t.Error("expected empty string key to exist")
		}
		if value != "empty key value" {
			t.Errorf("expected 'empty key value', got %v", value)
		}
	})
}

func TestMapGet(t *testing.T) {
	t.Run("gets existing value", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("key", "value")

		value, exists := m.Get("key")

		if !exists {
			t.Error("expected key to exist")
		}
		if value != "value" {
			t.Errorf("expected 'value', got %v", value)
		}
	})

	t.Run("returns false for non-existing key", func(t *testing.T) {
		m := NewMap[string, any]()

		_, exists := m.Get("nonexistent")

		if exists {
			t.Error("expected key to not exist")
		}
	})

	t.Run("returns nil value correctly", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("nilkey", nil)

		value, exists := m.Get("nilkey")

		if !exists {
			t.Error("expected key to exist")
		}
		if value != nil {
			t.Errorf("expected nil, got %v", value)
		}
	})

	t.Run("distinguishes nil value from missing key", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("nilkey", nil)

		_, existsNil := m.Get("nilkey")
		_, existsMissing := m.Get("missing")

		if !existsNil {
			t.Error("expected nil key to exist")
		}
		if existsMissing {
			t.Error("expected missing key to not exist")
		}
	})
}

func TestMapDelete(t *testing.T) {
	t.Run("deletes existing key", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("key", "value")

		m.Delete("key")

		if m.Has("key") {
			t.Error("expected key to be deleted")
		}
		if m.Len() != 0 {
			t.Errorf("expected length 0, got %d", m.Len())
		}
	})

	t.Run("does not panic on non-existing key", func(t *testing.T) {
		m := NewMap[string, any]()

		m.Delete("nonexistent")

		if m.Len() != 0 {
			t.Errorf("expected length 0, got %d", m.Len())
		}
	})

	t.Run("deletes only specified key", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("key1", "value1")
		m.Set("key2", "value2")
		m.Set("key3", "value3")

		m.Delete("key2")

		if m.Len() != 2 {
			t.Errorf("expected length 2, got %d", m.Len())
		}
		if !m.Has("key1") {
			t.Error("expected key1 to exist")
		}
		if m.Has("key2") {
			t.Error("expected key2 to be deleted")
		}
		if !m.Has("key3") {
			t.Error("expected key3 to exist")
		}
	})
}

func TestMapHas(t *testing.T) {
	t.Run("returns true for existing key", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("key", "value")

		if !m.Has("key") {
			t.Error("expected Has to return true")
		}
	})

	t.Run("returns false for non-existing key", func(t *testing.T) {
		m := NewMap[string, any]()

		if m.Has("nonexistent") {
			t.Error("expected Has to return false")
		}
	})

	t.Run("returns true for key with nil value", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("nilkey", nil)

		if !m.Has("nilkey") {
			t.Error("expected Has to return true for nil value")
		}
	})

	t.Run("returns false after deletion", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("key", "value")
		m.Delete("key")

		if m.Has("key") {
			t.Error("expected Has to return false after deletion")
		}
	})
}

func TestMapLen(t *testing.T) {
	t.Run("returns 0 for empty map", func(t *testing.T) {
		m := NewMap[string, any]()

		if m.Len() != 0 {
			t.Errorf("expected length 0, got %d", m.Len())
		}
	})

	t.Run("returns correct count after additions", func(t *testing.T) {
		m := NewMap[string, any]()

		for i := 0; i < 10; i++ {
			m.Set(string(rune('a'+i)), i)
		}

		if m.Len() != 10 {
			t.Errorf("expected length 10, got %d", m.Len())
		}
	})

	t.Run("returns correct count after deletions", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)
		m.Delete("b")

		if m.Len() != 2 {
			t.Errorf("expected length 2, got %d", m.Len())
		}
	})

	t.Run("does not increase on overwrite", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("key", "value1")
		m.Set("key", "value2")
		m.Set("key", "value3")

		if m.Len() != 1 {
			t.Errorf("expected length 1, got %d", m.Len())
		}
	})
}

func TestMapKeys(t *testing.T) {
	t.Run("returns empty slice for empty map", func(t *testing.T) {
		m := NewMap[string, any]()

		keys := m.Keys()

		if len(keys) != 0 {
			t.Errorf("expected 0 keys, got %d", len(keys))
		}
	})

	t.Run("returns all keys", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		keys := m.Keys()

		if len(keys) != 3 {
			t.Errorf("expected 3 keys, got %d", len(keys))
		}

		keySet := make(map[string]bool)
		for _, k := range keys {
			keySet[k] = true
		}

		if !keySet["a"] || !keySet["b"] || !keySet["c"] {
			t.Error("expected all keys to be present")
		}
	})

	t.Run("returns copy of keys", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("key", "value")

		keys := m.Keys()
		keys[0] = "modified"

		originalKeys := m.Keys()
		if originalKeys[0] == "modified" {
			t.Error("modifying returned keys should not affect map")
		}
	})
}

func TestMapValues(t *testing.T) {
	t.Run("returns empty slice for empty map", func(t *testing.T) {
		m := NewMap[string, any]()

		values := m.Values()

		if len(values) != 0 {
			t.Errorf("expected 0 values, got %d", len(values))
		}
	})

	t.Run("returns all values", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		values := m.Values()

		if len(values) != 3 {
			t.Errorf("expected 3 values, got %d", len(values))
		}

		sum := 0
		for _, v := range values {
			sum += v.(int)
		}
		if sum != 6 {
			t.Errorf("expected sum 6, got %d", sum)
		}
	})

	t.Run("includes nil values", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("a", nil)
		m.Set("b", "value")

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

func TestMapClear(t *testing.T) {
	t.Run("clears all entries", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		m.Clear()

		if m.Len() != 0 {
			t.Errorf("expected length 0 after clear, got %d", m.Len())
		}
	})

	t.Run("clear on empty map does not panic", func(t *testing.T) {
		m := NewMap[string, any]()

		m.Clear()

		if m.Len() != 0 {
			t.Errorf("expected length 0, got %d", m.Len())
		}
	})

	t.Run("map is usable after clear", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("key", "value")
		m.Clear()

		m.Set("newkey", "newvalue")

		if m.Len() != 1 {
			t.Errorf("expected length 1, got %d", m.Len())
		}
		value, exists := m.Get("newkey")
		if !exists || value != "newvalue" {
			t.Error("expected to get 'newvalue'")
		}
	})

	t.Run("old keys not accessible after clear", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("oldkey", "oldvalue")
		m.Clear()

		if m.Has("oldkey") {
			t.Error("expected old key to not exist after clear")
		}
	})
}

func TestMapRange(t *testing.T) {
	t.Run("iterates over all entries", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		count := 0
		m.Range(func(key string, value any) bool {
			count++
			return true
		})

		if count != 3 {
			t.Errorf("expected 3 iterations, got %d", count)
		}
	})

	t.Run("stops iteration when function returns false", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		count := 0
		m.Range(func(key string, value any) bool {
			count++
			return count < 2
		})

		if count != 2 {
			t.Errorf("expected 2 iterations, got %d", count)
		}
	})

	t.Run("handles empty map", func(t *testing.T) {
		m := NewMap[string, any]()

		count := 0
		m.Range(func(key string, value any) bool {
			count++
			return true
		})

		if count != 0 {
			t.Errorf("expected 0 iterations, got %d", count)
		}
	})

	t.Run("provides correct key-value pairs", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("a", 1)
		m.Set("b", 2)

		collected := make(map[string]int)
		m.Range(func(key string, value any) bool {
			collected[key] = value.(int)
			return true
		})

		if collected["a"] != 1 || collected["b"] != 2 {
			t.Error("expected correct key-value pairs")
		}
	})

	t.Run("stops immediately on first false return", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		count := 0
		m.Range(func(key string, value any) bool {
			count++
			return false
		})

		if count != 1 {
			t.Errorf("expected 1 iteration, got %d", count)
		}
	})
}

func TestMapOperationsSequence(t *testing.T) {
	t.Run("complex operations sequence", func(t *testing.T) {
		m := NewMap[string, any]()

		// Add entries
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		if m.Len() != 3 {
			t.Errorf("expected length 3, got %d", m.Len())
		}

		// Update entry
		m.Set("b", 20)
		value, _ := m.Get("b")
		if value != 20 {
			t.Errorf("expected 20, got %v", value)
		}

		// Delete entry
		m.Delete("a")
		if m.Has("a") {
			t.Error("expected 'a' to be deleted")
		}

		// Check remaining
		if m.Len() != 2 {
			t.Errorf("expected length 2, got %d", m.Len())
		}

		// Clear and restart
		m.Clear()
		m.Set("new", "entry")

		if m.Len() != 1 {
			t.Errorf("expected length 1, got %d", m.Len())
		}
	})
}

func TestMapConcurrentExecution(t *testing.T) {
	t.Run("handles concurrent reads on separate maps", func(t *testing.T) {
		var wg sync.WaitGroup
		const numGoroutines = 100
		const numReads = 100

		m := NewMap[string, any]()
		m.Set("key", "value")

		var successCount atomic.Int64

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < numReads; j++ {
					value, exists := m.Get("key")
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
		m := NewMap[string, any]()
		m.Set("a", 1)
		m.Set("b", 2)

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
		m := NewMap[string, any]()
		m.Set("exists", true)

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
					if m.Has("exists") {
						existsCount.Add(1)
					}
					if !m.Has("notexists") {
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
		m := NewMap[string, any]()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

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
					m.Range(func(key string, value any) bool {
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

func TestMapCopy(t *testing.T) {
	t.Run("copies empty map", func(t *testing.T) {
		m := NewMap[string, any]()

		copied := m.Copy()

		if copied == nil {
			t.Error("expected copied map to be created")
		}
		if copied.Len() != 0 {
			t.Errorf("expected length 0, got %d", copied.Len())
		}
	})

	t.Run("copies all entries", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		copied := m.Copy()

		if copied.Len() != 3 {
			t.Errorf("expected length 3, got %d", copied.Len())
		}

		for _, key := range []string{"a", "b", "c"} {
			origVal, _ := m.Get(key)
			copyVal, exists := copied.Get(key)
			if !exists {
				t.Errorf("expected key %s to exist in copy", key)
			}
			if origVal != copyVal {
				t.Errorf("expected value %v, got %v", origVal, copyVal)
			}
		}
	})

	t.Run("creates independent copy", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("key", "original")

		copied := m.Copy()
		copied.Set("key", "modified")

		origVal, _ := m.Get("key")
		if origVal != "original" {
			t.Error("modifying copy should not affect original")
		}

		copyVal, _ := copied.Get("key")
		if copyVal != "modified" {
			t.Error("copy should have modified value")
		}
	})
}

func TestMapPutIfNotExists(t *testing.T) {
	t.Run("adds new key", func(t *testing.T) {
		m := NewMap[string, any]()

		val, wasNew := m.PutIfNotExists("key", "value")

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
		m := NewMap[string, any]()
		m.Set("key", "original")

		val, wasNew := m.PutIfNotExists("key", "new")

		if wasNew {
			t.Error("expected key to already exist")
		}
		if val != "original" {
			t.Errorf("expected 'original', got %v", val)
		}

		storedVal, _ := m.Get("key")
		if storedVal != "original" {
			t.Errorf("expected stored value 'original', got %v", storedVal)
		}
	})
}

func TestMapEdgeCases(t *testing.T) {
	t.Run("handles special string keys", func(t *testing.T) {
		m := NewMap[string, any]()

		specialKeys := []string{
			"",
			" ",
			"\t",
			"\n",
			"key with spaces",
			"key\twith\ttabs",
			"unicode: 日本語",
			"emoji: 🚀",
		}

		for i, key := range specialKeys {
			m.Set(key, i)
		}

		for i, key := range specialKeys {
			value, exists := m.Get(key)
			if !exists {
				t.Errorf("expected key %q to exist", key)
			}
			if value != i {
				t.Errorf("expected value %d for key %q, got %v", i, key, value)
			}
		}
	})

	t.Run("handles large number of entries", func(t *testing.T) {
		m := NewMap[string, any]()
		const numEntries = 10000

		for i := 0; i < numEntries; i++ {
			key := string(rune(i))
			m.Set(key, i)
		}

		if m.Len() != numEntries {
			t.Errorf("expected length %d, got %d", numEntries, m.Len())
		}
	})

	t.Run("handles complex value types", func(t *testing.T) {
		m := NewMap[string, any]()

		type customStruct struct {
			Name  string
			Value int
		}

		m.Set("struct", customStruct{Name: "test", Value: 42})
		m.Set("pointer", &customStruct{Name: "ptr", Value: 100})
		m.Set("func", func() int { return 1 })
		m.Set("channel", make(chan int))

		if m.Len() != 4 {
			t.Errorf("expected length 4, got %d", m.Len())
		}

		structVal, _ := m.Get("struct")
		if s, ok := structVal.(customStruct); !ok || s.Name != "test" {
			t.Error("expected to get struct value")
		}
	})

	t.Run("delete then re-add same key", func(t *testing.T) {
		m := NewMap[string, any]()

		m.Set("key", "first")
		m.Delete("key")
		m.Set("key", "second")

		value, exists := m.Get("key")
		if !exists {
			t.Error("expected key to exist after re-add")
		}
		if value != "second" {
			t.Errorf("expected 'second', got %v", value)
		}
	})

	t.Run("range modification during iteration does not panic", func(t *testing.T) {
		m := NewMap[string, any]()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		var keysToDelete []string
		m.Range(func(key string, value any) bool {
			keysToDelete = append(keysToDelete, key)
			return true
		})

		for _, key := range keysToDelete {
			m.Delete(key)
		}

		if m.Len() != 0 {
			t.Errorf("expected length 0, got %d", m.Len())
		}
	})
}
