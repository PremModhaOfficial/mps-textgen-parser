package types

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSwissMap(t *testing.T) {
	t.Run("creates empty swiss-map", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()

		assertions.NotNil(m, "expected map to be created")
		assertions.Equal(0, m.Len())
	})

	t.Run("creates with capacity", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMapWithCapacity(100)

		assertions.NotNil(m, "expected map to be created")
		assertions.Equal(0, m.Len())
	})

	t.Run("creates independent instances", func(t *testing.T) {
		assertions := assert.New(t)
		m1 := NewSwissMap()
		m2 := NewSwissMap()

		m1.Set("key", "value1")
		m2.Set("key", "value2")

		v1, _ := m1.Get("key")
		v2, _ := m2.Get("key")

		assertions.NotEqual(v1, v2)
	})
}

func TestSwissMapSet(t *testing.T) {
	t.Run("sets single value", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()

		m.Set("key", "value")

		assertions.Equal(1, m.Len())
	})

	t.Run("sets multiple values", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()

		m.Set("key1", "value1")
		m.Set("key2", "value2")
		m.Set("key3", "value3")

		assertions.Equal(3, m.Len())
	})

	t.Run("overwrites existing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()

		m.Set("key", "value1")
		m.Set("key", "value2")

		assertions.Equal(1, m.Len())

		value, exists := m.Get("key")
		assertions.True(exists)
		assertions.Equal("value2", value)
	})

	t.Run("sets various types", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()

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
		m := NewSwissMap()

		m.Set("", "empty key value")

		value, exists := m.Get("")
		assertions.True(exists)
		assertions.Equal("empty key value", value)
	})
}

func TestSwissMapGet(t *testing.T) {
	t.Run("gets existing value", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("key", "value")

		value, exists := m.Get("key")

		assertions.True(exists)
		assertions.Equal("value", value)
	})

	t.Run("returns false for non-existing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()

		_, exists := m.Get("nonexistent")

		assertions.False(exists)
	})

	t.Run("returns nil value correctly", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("nilkey", nil)

		value, exists := m.Get("nilkey")

		assertions.True(exists)
		assertions.Nil(value)
	})

	t.Run("distinguishes nil value from missing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("nilkey", nil)

		_, existsNil := m.Get("nilkey")
		_, existsMissing := m.Get("missing")

		assertions.True(existsNil)
		assertions.False(existsMissing)
	})
}

func TestSwissMapDelete(t *testing.T) {
	t.Run("deletes existing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("key", "value")

		deleted := m.Delete("key")

		assertions.True(deleted)
		assertions.False(m.Has("key"))
		assertions.Equal(0, m.Len())
	})

	t.Run("returns false for non-existing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()

		deleted := m.Delete("nonexistent")

		assertions.False(deleted)
	})

	t.Run("deletes only specified key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("key1", "value1")
		m.Set("key2", "value2")
		m.Set("key3", "value3")

		m.Delete("key2")

		assertions.Equal(2, m.Len())
		assertions.True(m.Has("key1"))
		assertions.False(m.Has("key2"))
		assertions.True(m.Has("key3"))
	})
}

func TestSwissMapHas(t *testing.T) {
	t.Run("returns true for existing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("key", "value")

		assertions.True(m.Has("key"))
	})

	t.Run("returns false for non-existing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()

		assertions.False(m.Has("nonexistent"))
	})

	t.Run("returns true for key with nil value", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("nilkey", nil)

		assertions.True(m.Has("nilkey"))
	})

	t.Run("returns false after deletion", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("key", "value")
		m.Delete("key")

		assertions.False(m.Has("key"))
	})
}

func TestSwissMapLen(t *testing.T) {
	t.Run("returns 0 for empty swiss-map", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()

		assertions.Equal(0, m.Len())
	})

	t.Run("returns correct count after additions", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()

		for i := 0; i < 10; i++ {
			m.Set(string(rune('a'+i)), i)
		}

		assertions.Equal(10, m.Len())
	})

	t.Run("returns correct count after deletions", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)
		m.Delete("b")

		assertions.Equal(2, m.Len())
	})

	t.Run("does not increase on overwrite", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("key", "value1")
		m.Set("key", "value2")
		m.Set("key", "value3")

		assertions.Equal(1, m.Len())
	})
}

func TestSwissMapKeys(t *testing.T) {
	t.Run("returns empty slice for empty map", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()

		keys := m.Keys()

		assertions.Equal(0, len(keys))
	})

	t.Run("returns all keys", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		keys := m.Keys()

		assertions.Equal(3, len(keys))

		keySet := make(map[string]bool)
		for _, k := range keys {
			keySet[k] = true
		}

		assertions.True(keySet["a"])
		assertions.True(keySet["b"])
		assertions.True(keySet["c"])
	})

	t.Run("returns copy of keys", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("key", "value")

		keys := m.Keys()
		keys[0] = "modified"

		originalKeys := m.Keys()
		assertions.NotEqual("modified", originalKeys[0])
	})
}

func TestSwissMapValues(t *testing.T) {
	t.Run("returns empty slice for empty map", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()

		values := m.Values()

		assertions.Equal(0, len(values))
	})

	t.Run("returns all values", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
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
		m := NewSwissMap()
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
		assertions.True(hasNil)
	})
}

func TestSwissMapClear(t *testing.T) {
	t.Run("clears all entries", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		m.Clear()

		assertions.Equal(0, m.Len())
	})

	t.Run("clear on empty map does not panic", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()

		m.Clear()

		assertions.Equal(0, m.Len())
	})

	t.Run("map is usable after clear", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("key", "value")
		m.Clear()

		m.Set("newkey", "newvalue")

		assertions.Equal(1, m.Len())
		value, exists := m.Get("newkey")
		assertions.True(exists)
		assertions.Equal("newvalue", value)
	})

	t.Run("old keys not accessible after clear", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("oldkey", "oldvalue")
		m.Clear()

		assertions.False(m.Has("oldkey"))
	})
}

func TestSwissMapRange(t *testing.T) {
	t.Run("iterates over all swiss-map entries", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
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
		m := NewSwissMap()
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
		m := NewSwissMap()

		count := 0
		m.Range(func(key string, value any) bool {
			count++
			return true
		})

		assertions.Equal(0, count)
	})

	t.Run("provides correct key-value pairs", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
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
		m := NewSwissMap()
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

func TestSwissMapCopy(t *testing.T) {
	t.Run("copies empty map", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()

		copied := m.Copy()

		assertions.NotNil(copied)
		assertions.Equal(0, copied.Len())
	})

	t.Run("copies all entries", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("a", 1)
		m.Set("b", 2)
		m.Set("c", 3)

		copied := m.Copy()

		assertions.Equal(3, copied.Len())

		for _, key := range []string{"a", "b", "c"} {
			origVal, _ := m.Get(key)
			copyVal, exists := copied.Get(key)
			assertions.True(exists)
			assertions.Equal(origVal, copyVal)
		}
	})

	t.Run("creates independent copy", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("key", "original")

		copied := m.Copy()
		copied.Set("key", "modified")

		origVal, _ := m.Get("key")
		assertions.Equal("original", origVal)

		copyVal, _ := copied.Get("key")
		assertions.Equal("modified", copyVal)
	})

	t.Run("adding to copy does not affect original", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("a", 1)

		copied := m.Copy()
		copied.Set("b", 2)

		assertions.Equal(1, m.Len())
		assertions.Equal(2, copied.Len())
	})

	t.Run("deleting from copy does not affect original", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("a", 1)
		m.Set("b", 2)

		copied := m.Copy()
		copied.Delete("a")

		assertions.Equal(2, m.Len())
		assertions.True(m.Has("a"))
		assertions.Equal(1, copied.Len())
	})

	t.Run("clearing copy does not affect original", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("a", 1)
		m.Set("b", 2)

		copied := m.Copy()
		copied.Clear()

		assertions.Equal(2, m.Len())
		assertions.Equal(0, copied.Len())
	})

	t.Run("copies nil values correctly", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("nilkey", nil)

		copied := m.Copy()

		value, exists := copied.Get("nilkey")
		assertions.True(exists)
		assertions.Nil(value)
	})

	t.Run("copies various types", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("string", "hello")
		m.Set("int", 42)
		m.Set("float", 3.14)
		m.Set("bool", true)
		m.Set("slice", []int{1, 2, 3})

		copied := m.Copy()

		assertions.Equal(5, copied.Len())

		strVal, _ := copied.Get("string")
		assertions.Equal("hello", strVal)

		intVal, _ := copied.Get("int")
		assertions.Equal(42, intVal)
	})
}

func TestSwissMapPutIfNotExists(t *testing.T) {
	t.Run("adds new key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()

		val, wasNew := m.PutIfNotExists("key", "value")

		assertions.True(wasNew)
		assertions.Equal("value", val)
		assertions.Equal(1, m.Len())
	})

	t.Run("does not overwrite existing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("key", "original")

		val, wasNew := m.PutIfNotExists("key", "new")

		assertions.False(wasNew)
		assertions.Equal("original", val)

		storedVal, _ := m.Get("key")
		assertions.Equal("original", storedVal)
	})

	t.Run("handles nil value for new key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()

		val, wasNew := m.PutIfNotExists("nilkey", nil)

		assertions.True(wasNew)
		assertions.Nil(val)
		assertions.True(m.Has("nilkey"))
	})

	t.Run("handles nil value for existing key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("key", nil)

		val, wasNew := m.PutIfNotExists("key", "new")

		assertions.False(wasNew)
		assertions.Nil(val)
	})

	t.Run("handles empty string key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()

		val, wasNew := m.PutIfNotExists("", "empty key value")

		assertions.True(wasNew)
		assertions.Equal("empty key value", val)
	})

	t.Run("multiple PutIfNotExists calls", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()

		// First insertion
		val1, wasNew1 := m.PutIfNotExists("key", "first")
		assertions.True(wasNew1)
		assertions.Equal("first", val1)

		// Second insertion (should not overwrite)
		val2, wasNew2 := m.PutIfNotExists("key", "second")
		assertions.False(wasNew2)
		assertions.Equal("first", val2)

		// Third insertion (should not overwrite)
		val3, wasNew3 := m.PutIfNotExists("key", "third")
		assertions.False(wasNew3)
		assertions.Equal("first", val3)

		assertions.Equal(1, m.Len())
	})
}

func TestSwissMapCapacity(t *testing.T) {
	t.Run("returns capacity", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMapWithCapacity(100)

		capacity := m.Capacity()

		assertions.GreaterOrEqual(capacity, 0)
	})
}

func TestSwissMapOperationsSequence(t *testing.T) {
	t.Run("complex operations sequence", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()

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
		assertions.False(m.Has("a"))

		// Check remaining
		assertions.Equal(2, m.Len())

		// Clear and restart
		m.Clear()
		m.Set("new", "entry")

		assertions.Equal(1, m.Len())
	})
}

func TestSwissMapConcurrentExecution(t *testing.T) {
	t.Run("handles concurrent reads on same map", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
		m.Set("key", "value")

		var wg sync.WaitGroup
		const numGoroutines = 100
		const numReads = 100

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
		assertions.Equal(expectedReads, successCount.Load())
	})

	t.Run("concurrent Len calls are safe", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
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
		assertions.Equal(expectedReads, correctCount.Load())
	})

	t.Run("concurrent Has calls are safe", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()
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
		assertions.Equal(expectedChecks, existsCount.Load())
		assertions.Equal(expectedChecks, notExistsCount.Load())
	})
}

func TestSwissMapEdgeCases(t *testing.T) {
	t.Run("handles special string keys", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()

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
			assertions.True(exists)
			assertions.Equal(i, value)
		}
	})

	t.Run("handles large number of entries", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMapWithCapacity(10000)
		const numEntries = 10000

		for i := 0; i < numEntries; i++ {
			key := string(rune(i))
			m.Set(key, i)
		}

		assertions.Equal(numEntries, m.Len())
	})

	t.Run("handles complex value types", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()

		type customStruct struct {
			Name  string
			Value int
		}

		m.Set("struct", customStruct{Name: "test", Value: 42})
		m.Set("pointer", &customStruct{Name: "ptr", Value: 100})
		m.Set("func", func() int { return 1 })
		m.Set("channel", make(chan int))

		assertions.Equal(4, m.Len())

		structVal, _ := m.Get("struct")
		s, ok := structVal.(customStruct)
		assertions.True(ok)
		assertions.Equal("test", s.Name)
	})

	t.Run("delete then re-add same key", func(t *testing.T) {
		assertions := assert.New(t)
		m := NewSwissMap()

		m.Set("key", "first")
		m.Delete("key")
		m.Set("key", "second")

		value, exists := m.Get("key")
		assertions.True(exists)
		assertions.Equal("second", value)
	})
}
