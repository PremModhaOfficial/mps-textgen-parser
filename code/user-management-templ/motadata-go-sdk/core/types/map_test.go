package types

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMap(t *testing.T) {
	t.Run("creates empty map", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()

		assertions.NotNil(m, "expected map to be created")
		assertions.Equal(0, m.Len())
	})

	t.Run("creates independent instances", func(t *testing.T) {
		assertions := assert.New(t)
		m1 := NewMap()
		m2 := NewMap()

		m1.Set("key", "value1")
		m2.Set("key", "value2")

		v1, _ := m1.Get("key")
		v2, _ := m2.Get("key")

		assertions.NotEqual(v1, v2, "expected independent maps")
	})
}

func TestMapSet(t *testing.T) {
	t.Run("sets single value", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMapWithCapacity(1)

		m.Set("key", "value")

		assertions.Equal(1, m.Len())
	})

	t.Run("sets multiple values", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()

		m.Set("key1", "value1")
		m.Set("key2", "value2")
		m.Set("key3", "value3")

		assertions.Equal(3, m.Len())
	})

	t.Run("overwrites existing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()

		m.Set("key", "value1")
		m.Set("key", "value2")

		assertions.Equal(1, m.Len())

		value, exists := m.Get("key")
		assertions.True(exists, "expected key to exist")
		assertions.Equal("value2", value)
	})

	t.Run("sets various types", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()

		m.Set("string", "hello")
		m.Set("int", 42)
		m.Set("float", 3.14)
		m.Set("bool", true)
		m.Set("nil", nil)
		m.Set("slice", []int{1, 2, 3})
		m.Set("map", map[string]int{"a": 1})

		assertions.Equal(7, m.Len())
	})

	t.Run("sets empty string key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()

		m.Set("", "empty key value")

		value, exists := m.Get("")
		assertions.True(exists, "expected empty string key to exist")
		assertions.Equal("empty key value", value)
	})
}

func TestMapGet(t *testing.T) {
	t.Run("gets existing value", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("key", "value")

		value, exists := m.Get("key")

		assertions.True(exists, "expected key to exist")
		assertions.Equal("value", value)
	})

	t.Run("returns false for non-existing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()

		_, exists := m.Get("nonexistent")

		assertions.False(exists, "expected key to not exist")
	})

	t.Run("returns nil value correctly", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("nilkey", nil)

		value, exists := m.Get("nilkey")

		assertions.True(exists, "expected key to exist")
		assertions.Nil(value)
	})

	t.Run("distinguishes nil value from missing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("nilkey", nil)

		_, existsNil := m.Get("nilkey")
		_, existsMissing := m.Get("missing")

		assertions.True(existsNil, "expected nil key to exist")
		assertions.False(existsMissing, "expected missing key to not exist")
	})
}

func TestMapDelete(t *testing.T) {
	t.Run("deletes existing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("key", "value")

		m.Delete("key")

		assertions.False(m.Has("key"), "expected key to be deleted")
		assertions.Equal(0, m.Len())
	})

	t.Run("does not panic on non-existing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()

		m.Delete("nonexistent")

		assertions.Equal(0, m.Len())
	})

	t.Run("deletes only specified key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("key1", "value1")
		m.Set("key2", "value2")
		m.Set("key3", "value3")

		m.Delete("key2")

		assertions.Equal(2, m.Len())
		assertions.True(m.Has("key1"), "expected key1 to exist")
		assertions.False(m.Has("key2"), "expected key2 to be deleted")
		assertions.True(m.Has("key3"), "expected key3 to exist")
	})
}

func TestMapHas(t *testing.T) {
	t.Run("returns true for existing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("key", "value")

		assertions.True(m.Has("key"), "expected Has to return true")
	})

	t.Run("returns false for non-existing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()

		assertions.False(m.Has("nonexistent"), "expected Has to return false")
	})

	t.Run("returns true for key with nil value", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("nilkey", nil)

		assertions.True(m.Has("nilkey"), "expected Has to return true for nil value")
	})

	t.Run("returns false after deletion", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("key", "value")
		m.Delete("key")

		assertions.False(m.Has("key"), "expected Has to return false after deletion")
	})
}

func TestMapLen(t *testing.T) {
	t.Run("returns 0 for empty map", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()

		assertions.Equal(0, m.Len())
	})

	t.Run("returns correct count after additions", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()

		for i := 0; i < 10; i++ {
			m.Set(string(rune('a'+i)), i)
		}

		assertions.Equal(10, m.Len())
	})

	t.Run("returns correct count after deletions", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)
		m.Delete("b")

		assertions.Equal(2, m.Len())
	})

	t.Run("does not increase on overwrite", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("key", "value1")
		m.Set("key", "value2")
		m.Set("key", "value3")

		assertions.Equal(1, m.Len())
	})
}

func TestMapKeys(t *testing.T) {
	t.Run("returns empty slice for empty map", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()

		keys := m.Keys()

		assertions.Equal(0, len(keys))
	})

	t.Run("returns all keys", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		keys := m.Keys()

		assertions.Equal(3, len(keys))

		keySet := make(map[string]bool)
		for _, k := range keys {
			keySet[k] = true
		}

		assertions.True(keySet["a"] && keySet["b"] && keySet["c"], "expected all keys to be present")
	})

	t.Run("returns copy of keys", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("key", "value")

		keys := m.Keys()
		keys[0] = "modified"

		originalKeys := m.Keys()
		assertions.NotEqual("modified", originalKeys[0], "modifying returned keys should not affect map")
	})
}

func TestMapValues(t *testing.T) {
	t.Run("returns empty slice value for empty map", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()

		values := m.Values()

		assertions.Equal(0, len(values))
	})

	t.Run("returns all values", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		values := m.Values()

		assertions.Equal(3, len(values))

		sum := 0
		for _, v := range values {
			sum += v.(int)
		}
		assertions.Equal(6, sum)
	})

	t.Run("includes nil values", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("a", nil)
		m.Set("b", "value")

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

func TestMapClear(t *testing.T) {
	t.Run("clears all entries", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		m.Clear()

		assertions.Equal(0, m.Len())
	})

	t.Run("clear on empty map does not panic", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()

		m.Clear()

		assertions.Equal(0, m.Len())
	})

	t.Run("map is usable after clear", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("key", "value")
		m.Clear()

		m.Set("newkey", "newvalue")

		assertions.Equal(1, m.Len())
		value, exists := m.Get("newkey")
		assertions.True(exists && value == "newvalue", "expected to get 'newvalue'")
	})

	t.Run("old keys not accessible after clear", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("oldkey", "oldvalue")
		m.Clear()

		assertions.False(m.Has("oldkey"), "expected old key to not exist after clear")
	})
}

func TestMapRange(t *testing.T) {
	t.Run("iterates over all entries", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		count := 0
		m.Range(func(key string, value any) bool {
			count++
			return true
		})

		assertions.Equal(3, count)
	})

	t.Run("stops iteration when function returns false", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		count := 0
		m.Range(func(key string, value any) bool {
			count++
			return count < 2
		})

		assertions.Equal(2, count)
	})

	t.Run("handles empty map", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()

		count := 0
		m.Range(func(key string, value any) bool {
			count++
			return true
		})

		assertions.Equal(0, count)
	})

	t.Run("provides correct key-value pairs", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("a", 1)
		m.Set("b", 2)

		collected := make(map[string]int)
		m.Range(func(key string, value any) bool {
			collected[key] = value.(int)
			return true
		})

		assertions.Equal(1, collected["a"])
		assertions.Equal(2, collected["b"])
	})

	t.Run("stops immediately on first false return", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		count := 0
		m.Range(func(key string, value any) bool {
			count++
			return false
		})

		assertions.Equal(1, count)
	})
}

func TestMapOperationsSequence(t *testing.T) {
	t.Run("complex operations sequence", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()

		// Add entries
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		assertions.Equal(3, m.Len())

		// Update entry
		m.Set("b", 20)
		value, _ := m.Get("b")
		assertions.Equal(20, value)

		// Delete entry
		m.Delete("a")
		assertions.False(m.Has("a"), "expected 'a' to be deleted")

		// Check remaining
		assertions.Equal(2, m.Len())

		// Clear and restart
		m.Clear()
		m.Set("new", "entry")

		assertions.Equal(1, m.Len())
	})
}

func TestMapConcurrentExecution(t *testing.T) {
	t.Run("handles concurrent reads on separate maps", func(t *testing.T) {
		var wg sync.WaitGroup
		const numGoroutines = 100
		const numReads = 100

		m := NewMap()
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

		assertions := assert.New(t)
		expectedReads := int64(numGoroutines * numReads)
		assertions.Equal(expectedReads, successCount.Load())
	})

	t.Run("concurrent Len calls are safe", func(t *testing.T) {
		m := NewMap()
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

		assertions := assert.New(t)
		expectedReads := int64(numGoroutines * numReads)
		assertions.Equal(expectedReads, correctCount.Load())
	})

	t.Run("concurrent Has calls are safe", func(t *testing.T) {
		m := NewMap()
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

		assertions := assert.New(t)
		expectedChecks := int64(numGoroutines * numChecks)
		assertions.Equal(expectedChecks, existsCount.Load())
		assertions.Equal(expectedChecks, notExistsCount.Load())
	})

	t.Run("concurrent Range calls are safe", func(t *testing.T) {
		m := NewMap()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		const numGoroutines = 50
		const numRanges = 50

		var correctRanges atomic.Int64

		var wg sync.WaitGroup

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

		assertions := assert.New(t)
		expectedRanges := int64(numGoroutines * numRanges)
		assertions.Equal(expectedRanges, correctRanges.Load())
	})
}

func TestMapCopy(t *testing.T) {
	t.Run("copies empty map", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()

		copied := m.Copy()

		assertions.NotNil(copied, "expected copied map to be created")
		assertions.Equal(0, copied.Len())
	})

	t.Run("copies all entries", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		copied := m.Copy()

		assertions.Equal(3, copied.Len())

		for _, key := range []string{"a", "b", "c"} {
			origVal, _ := m.Get(key)
			copyVal, exists := copied.Get(key)
			assertions.True(exists, "expected key %s to exist in copy", key)
			assertions.Equal(origVal, copyVal)
		}
	})

	t.Run("creates independent copy", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("key", "original")

		copied := m.Copy()
		copied.Set("key", "modified")

		origVal, _ := m.Get("key")
		assertions.Equal("original", origVal, "modifying copy should not affect original")

		copyVal, _ := copied.Get("key")
		assertions.Equal("modified", copyVal, "copy should have modified value")
	})
}

func TestMapPutIfNotExists(t *testing.T) {
	t.Run("adds new key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()

		val, wasNew := m.PutIfNotExists("key", "value")

		assertions.True(wasNew, "expected key to be newly inserted")
		assertions.Equal("value", val)
		assertions.Equal(1, m.Len())
	})

	t.Run("does not overwrite existing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewMap()
		m.Set("key", "original")

		val, wasNew := m.PutIfNotExists("key", "new")

		assertions.False(wasNew, "expected key to already exist")
		assertions.Equal("original", val)

		storedVal, _ := m.Get("key")
		assertions.Equal("original", storedVal)
	})
}

func TestMapEdgeCases(t *testing.T) {
	t.Run("handles special string keys", func(t *testing.T) {
		m := NewMap()

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

		assertions := assert.New(t)
		for i, key := range specialKeys {
			value, exists := m.Get(key)
			assertions.True(exists, "expected key %q to exist", key)
			assertions.Equal(i, value, "expected value %d for key %q", i, key)
		}
	})

	t.Run("handles large number of entries", func(t *testing.T) {
		m := NewMap()
		const numEntries = 10000

		for i := 0; i < numEntries; i++ {
			key := string(rune(i))
			m.Set(key, i)
		}

		assertions := assert.New(t)
		assertions.Equal(numEntries, m.Len())
	})

	t.Run("handles complex value types", func(t *testing.T) {
		m := NewMap()

		type customStruct struct {
			Name  string
			Value int
		}

		m.Set("struct", customStruct{Name: "test", Value: 42})
		m.Set("pointer", &customStruct{Name: "ptr", Value: 100})
		m.Set("func", func() int { return 1 })
		m.Set("channel", make(chan int))

		assertions := assert.New(t)
		assertions.Equal(4, m.Len())

		structVal, _ := m.Get("struct")
		if s, ok := structVal.(customStruct); ok {
			assertions.Equal("test", s.Name)
		} else {
			assertions.Fail("expected to get struct value")
		}
	})

	t.Run("delete then re-add same key", func(t *testing.T) {
		m := NewMap()

		m.Set("key", "first")
		m.Delete("key")
		m.Set("key", "second")

		assertions := assert.New(t)
		value, exists := m.Get("key")
		assertions.True(exists, "expected key to exist after re-add")
		assertions.Equal("second", value)
	})

	t.Run("range modification during iteration does not panic", func(t *testing.T) {
		m := NewMap()
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

		assertions := assert.New(t)
		assertions.Equal(0, m.Len())
	})
}
