package codec

import (
	"bytes"
	"math"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ============== Codec Write Functions Tests ==============

func TestWriteInt8Value(t *testing.T) {
	tests := []struct {
		name     string
		value    int8
		expected []byte
	}{
		{"zero", 0, []byte{0}},
		{"positive", 127, []byte{127}},
		{"negative", -128, []byte{128}},
		{"small positive", 1, []byte{1}},
		{"small negative", -1, []byte{255}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			buffer := &bytes.Buffer{}
			WriteInt8Value(tt.value, buffer)

			assertions.Equal(tt.expected, buffer.Bytes(), "WriteInt8Value(%d) should produce expected bytes", tt.value)
		})
	}
}

func TestWriteInt16Value(t *testing.T) {
	tests := []struct {
		name     string
		value    int16
		expected []byte
	}{
		{"zero", 0, []byte{0, 0}},
		{"positive max", 32767, []byte{255, 127}},
		{"negative min", -32768, []byte{0, 128}},
		{"value 256", 256, []byte{0, 1}},
		{"value -1", -1, []byte{255, 255}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			buffer := &bytes.Buffer{}
			WriteInt16Value(tt.value, buffer)

			assertions.Equal(tt.expected, buffer.Bytes(), "WriteInt16Value(%d) should produce expected bytes", tt.value)
		})
	}
}

func TestWriteInt32Value(t *testing.T) {
	tests := []struct {
		name     string
		value    int32
		expected []byte
	}{
		{"zero", 0, []byte{0, 0, 0, 0}},
		{"positive max", 2147483647, []byte{255, 255, 255, 127}},
		{"negative min", -2147483648, []byte{0, 0, 0, 128}},
		{"value 16777216", 16777216, []byte{0, 0, 0, 1}},
		{"value -1", -1, []byte{255, 255, 255, 255}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			buffer := &bytes.Buffer{}
			WriteInt32Value(tt.value, buffer)

			assertions.Equal(tt.expected, buffer.Bytes(), "WriteInt32Value(%d) should produce expected bytes", tt.value)
		})
	}
}

func TestWriteInt64Value(t *testing.T) {
	tests := []struct {
		name     string
		value    int64
		expected []byte
	}{
		{"zero", 0, []byte{0, 0, 0, 0, 0, 0, 0, 0}},
		{"positive max", 9223372036854775807, []byte{255, 255, 255, 255, 255, 255, 255, 127}},
		{"negative min", -9223372036854775808, []byte{0, 0, 0, 0, 0, 0, 0, 128}},
		{"value -1", -1, []byte{255, 255, 255, 255, 255, 255, 255, 255}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			buffer := &bytes.Buffer{}
			WriteInt64Value(tt.value, buffer)

			assertions.Equal(tt.expected, buffer.Bytes(), "WriteInt64Value(%d) should produce expected bytes", tt.value)
		})
	}
}

func TestWriteInt24Value(t *testing.T) {
	tests := []struct {
		name     string
		value    int32
		expected []byte
	}{
		{"zero", 0, []byte{0x00, 0x00, 0x00}},
		{"positive", 0x123456, []byte{0x56, 0x34, 0x12}},
		{"negative", -1, []byte{0xFF, 0xFF, 0xFF}},
		{"max_int24", MaxInt24, []byte{0xFF, 0xFF, 0x7F}},
		{"min_int24", MinInt24, []byte{0x00, 0x00, 0x80}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			buffer := &bytes.Buffer{}
			WriteInt24Value(tt.value, buffer)

			assertions.Equal(tt.expected, buffer.Bytes(), "WriteInt24Value(%d) should produce expected bytes", tt.value)
		})
	}
}

func TestWriteInt40Value(t *testing.T) {
	tests := []struct {
		name     string
		value    int64
		expected []byte
	}{
		{"zero", 0, []byte{0x00, 0x00, 0x00, 0x00, 0x00}},
		{"positive", 0x123456789A, []byte{0x9A, 0x78, 0x56, 0x34, 0x12}},
		{"negative", -1, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			buffer := &bytes.Buffer{}
			WriteInt40Value(tt.value, buffer)

			assertions.Equal(tt.expected, buffer.Bytes(), "WriteInt40Value(%d) should produce expected bytes", tt.value)
		})
	}
}

func TestWriteInt48Value(t *testing.T) {
	tests := []struct {
		name     string
		value    int64
		expected []byte
	}{
		{"zero", 0, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"positive", 0x123456789ABC, []byte{0xBC, 0x9A, 0x78, 0x56, 0x34, 0x12}},
		{"negative", -1, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			buffer := &bytes.Buffer{}
			WriteInt48Value(tt.value, buffer)

			assertions.Equal(tt.expected, buffer.Bytes(), "WriteInt48Value(%d) should produce expected bytes", tt.value)
		})
	}
}

func TestWriteInt56Value(t *testing.T) {
	tests := []struct {
		name     string
		value    int64
		expected []byte
	}{
		{"zero", 0, []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"positive", 0x123456789ABCDE, []byte{0xDE, 0xBC, 0x9A, 0x78, 0x56, 0x34, 0x12}},
		{"negative", -1, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			buffer := &bytes.Buffer{}
			WriteInt56Value(tt.value, buffer)

			assertions.Equal(tt.expected, buffer.Bytes(), "WriteInt56Value(%d) should produce expected bytes", tt.value)
		})
	}
}

func TestReadINT24Value(t *testing.T) {
	tests := []struct {
		name     string
		bytes    []byte
		expected int32
	}{
		{"zero", []byte{0x00, 0x00, 0x00}, 0},
		{"positive", []byte{0x56, 0x34, 0x12}, 0x123456},
		{"negative", []byte{0xFF, 0xFF, 0xFF}, -1},
		{"max_int24", []byte{0xFF, 0xFF, 0x7F}, MaxInt24},
		{"min_int24", []byte{0x00, 0x00, 0x80}, MinInt24},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			result := ReadINT24Value(tt.bytes)
			assertions.Equal(tt.expected, result, "ReadINT24Value should decode correctly")
		})
	}
}

func TestReadINT40Value(t *testing.T) {
	tests := []struct {
		name     string
		bytes    []byte
		expected int64
	}{
		{"zero", []byte{0x00, 0x00, 0x00, 0x00, 0x00}, 0},
		{"positive", []byte{0x9A, 0x78, 0x56, 0x34, 0x12}, 0x123456789A},
		{"negative", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, -1},
		{"max_int40", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x7F}, MaxInt40},
		{"min_int40", []byte{0x00, 0x00, 0x00, 0x00, 0x80}, MinInt40},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			result := ReadINT40Value(tt.bytes)
			assertions.Equal(tt.expected, result, "ReadINT40Value should decode correctly")
		})
	}
}

func TestReadINT48Value(t *testing.T) {
	tests := []struct {
		name     string
		bytes    []byte
		expected int64
	}{
		{"zero", []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 0},
		{"positive", []byte{0xBC, 0x9A, 0x78, 0x56, 0x34, 0x12}, 0x123456789ABC},
		{"negative", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, -1},
		{"max_int48", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x7F}, MaxInt48},
		{"min_int48", []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x80}, MinInt48},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			result := ReadINT48Value(tt.bytes)
			assertions.Equal(tt.expected, result, "ReadINT48Value should decode correctly")
		})
	}
}

func TestReadINT56Value(t *testing.T) {
	tests := []struct {
		name     string
		bytes    []byte
		expected int64
	}{
		{"zero", []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 0},
		{"positive", []byte{0xDE, 0xBC, 0x9A, 0x78, 0x56, 0x34, 0x12}, 0x123456789ABCDE},
		{"negative", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, -1},
		{"max_int56", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x7F}, MaxInt56},
		{"min_int56", []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80}, MinInt56},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			result := ReadINT56Value(tt.bytes)
			assertions.Equal(tt.expected, result, "ReadINT56Value should decode correctly")
		})
	}
}

func TestReadINT64Value(t *testing.T) {
	tests := []struct {
		name     string
		bytes    []byte
		expected int64
	}{
		{"zero", []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 0},
		{"positive", []byte{0xEF, 0xCD, 0xAB, 0x89, 0x67, 0x45, 0x23, 0x01}, 0x0123456789ABCDEF},
		{"negative", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}, -1},
		{"max_int64", []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x7F}, math.MaxInt64},
		{"min_int64", []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80}, math.MinInt64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			result := ReadINT64Value(tt.bytes)
			assertions.Equal(tt.expected, result, "ReadINT64Value should decode correctly")
		})
	}
}

func TestReadINT16Value(t *testing.T) {
	tests := []struct {
		name     string
		bytes    []byte
		expected int16
	}{
		{"zero", []byte{0x00, 0x00}, 0},
		{"positive", []byte{0x34, 0x12}, 0x1234},
		{"negative", []byte{0xFF, 0xFF}, -1},
		{"max_int16", []byte{0xFF, 0x7F}, math.MaxInt16},
		{"min_int16", []byte{0x00, 0x80}, math.MinInt16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			result := ReadINT16Value(tt.bytes)
			assertions.Equal(tt.expected, result, "ReadINT16Value should decode correctly")
		})
	}
}

func TestReadINT32Value(t *testing.T) {
	tests := []struct {
		name     string
		bytes    []byte
		expected int32
	}{
		{"zero", []byte{0x00, 0x00, 0x00, 0x00}, 0},
		{"positive", []byte{0x78, 0x56, 0x34, 0x12}, 0x12345678},
		{"negative", []byte{0xFF, 0xFF, 0xFF, 0xFF}, -1},
		{"max_int32", []byte{0xFF, 0xFF, 0xFF, 0x7F}, math.MaxInt32},
		{"min_int32", []byte{0x00, 0x00, 0x00, 0x80}, math.MinInt32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			result := ReadINT32Value(tt.bytes)
			assertions.Equal(tt.expected, result, "ReadINT32Value should decode correctly")
		})
	}
}

func TestWriteFloat32Value(t *testing.T) {
	tests := []struct {
		name  string
		value float32
	}{
		{"zero", 0.0},
		{"positive", 3.14},
		{"negative", -3.14},
		{"max", math.MaxFloat32},
		{"smallest positive", math.SmallestNonzeroFloat32},
		{"infinity", float32(math.Inf(1))},
		{"negative infinity", float32(math.Inf(-1))},
		{"NaN", float32(math.NaN())},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			buffer := &bytes.Buffer{}
			WriteFloat32Value(tt.value, buffer)

			assertions.Equal(4, buffer.Len(), "WriteFloat32Value should write exactly 4 bytes")
		})
	}
}

func TestWriteFloat64Value(t *testing.T) {
	tests := []struct {
		name  string
		value float64
	}{
		{"zero", 0.0},
		{"positive", 3.141592653589793},
		{"negative", -3.141592653589793},
		{"max", math.MaxFloat64},
		{"smallest positive", math.SmallestNonzeroFloat64},
		{"infinity", math.Inf(1)},
		{"negative infinity", math.Inf(-1)},
		{"NaN", math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			buffer := &bytes.Buffer{}
			WriteFloat64Value(tt.value, buffer)

			assertions.Equal(8, buffer.Len(), "WriteFloat64Value should write exactly 8 bytes")
		})
	}
}

// ============== Header Tests ==============

func TestBuildAndParseHeader(t *testing.T) {
	tests := []struct {
		name     string
		encoder  Encoder
		expected byte
	}{
		{"custom", Custom, 0x00},
		{"msgpack", MsgPack, 0x01},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)

			header := BuildHeader(tt.encoder)
			assertions.Equal(tt.expected, header, "BuildHeader should produce expected byte")

			parsedCodec := ParseHeader(header)
			assertions.Equal(tt.encoder, parsedCodec, "ParseHeader should recover codec type")
		})
	}
}

// ============== Packer Functions Tests ==============

func TestPackMap(t *testing.T) {
	t.Run("empty map", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{}
		result, err := PackMap(data, Custom, nil)
		assertions.NoError(err)

		assertions.NotNil(result, "Result should not be nil")
		assertions.NotEmpty(result, "Result should not be empty for empty map")

		// Verify can unpack
		unpacked, err := UnpackMap(result)
		assertions.NoError(err, "Should unpack empty map without error")
		assertions.Empty(unpacked, "Unpacked map should be empty")
	})

	t.Run("simple map", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{
			"name":   "test",
			"count":  int64(42),
			"active": true,
		}

		result, err := PackMap(data, Custom, nil)
		assertions.NoError(err)

		unpacked, err := UnpackMap(result)
		assertions.NoError(err, "Unpack should not fail")

		assertions.Equal("test", unpacked["name"], "Name should match")
		assertions.Equal(int8(42), unpacked["count"], "Count should match")
		assertions.Equal(true, unpacked["active"], "Active should match")
	})

	t.Run("all data types", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{
			"int8":    int8(127),
			"int16":   int16(32767),
			"int32":   int32(2147483647),
			"int64":   int64(9223372036854775807),
			"int":     42,
			"uint8":   uint8(255),
			"uint16":  uint16(65535),
			"uint32":  uint32(4294967295),
			"uint64":  uint64(math.MaxInt64),
			"uint":    uint(100),
			"float32": float32(3.14),
			"float64": 3.141592653589793,
			"string":  "hello world",
			"bool":    true,
			"nil":     nil,
		}

		packed, err := PackMap(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err, "Unpack should not fail")

		// Check all keys exist
		for key := range data {
			assertions.Contains(unpacked, key, "Key %s should exist in unpacked map", key)
		}

		assertions.Equal(int16(255), unpacked["uint8"])
		assertions.Equal(int32(65535), unpacked["uint16"])
		assertions.Equal(int64(4294967295), unpacked["uint32"])
		assertions.Equal(int64(math.MaxInt64), unpacked["uint64"])
		assertions.Equal(int8(100), unpacked["uint"])
	})

	t.Run("uint64 out of range", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{
			"too_big": uint64(math.MaxUint64),
		}

		_, err := PackMap(data, Custom, nil)
		assertions.ErrorIs(err, utils.ErrUnsupportedDataType)
	})

	t.Run("nested map", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]interface{}{
					"value": "nested",
				},
			},
		}

		packed, err := PackMap(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err, "Unpack should not fail")

		level1, ok := unpacked["level1"].(map[string]interface{})
		assertions.True(ok, "level1 should be a map")

		level2, ok := level1["level2"].(map[string]interface{})
		assertions.True(ok, "level2 should be a map")

		assertions.Equal("nested", level2["value"], "Nested value should match")
	})

	t.Run("typed map value support", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{
			"typed": map[string]int{
				"a": 1,
				"b": 2,
			},
		}

		packed, err := PackMap(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		typedMap, ok := unpacked["typed"].(map[string]interface{})
		assertions.True(ok)
		assertions.Equal(int8(1), typedMap["a"])
		assertions.Equal(int8(2), typedMap["b"])
	})

	t.Run("map with array", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{
			"items": []interface{}{"item1", "item2", "item3"},
			"count": int64(3),
		}

		packed, err := PackMap(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err, "Unpack should not fail")

		items, ok := unpacked["items"].([]interface{})
		assertions.True(ok, "items should be an array")
		assertions.Len(items, 3, "Should have 3 items")
	})

	t.Run("map with byte array", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{
			"blob": []byte{1, 2, 3, 4, 5},
		}
		packed, err := PackMap(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		assertions.Equal([]byte{1, 2, 3, 4, 5}, unpacked["blob"])
	})

	t.Run("with buffer", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{"key": "value"}
		buffer := &bytes.Buffer{}

		result, err := PackMap(data, Custom, buffer)
		assertions.NoError(err)

		// Result starts with header byte followed by buffer contents
		assertions.Equal(byte(0x00), result[0], "Header should be 0x00 for Custom")
		assertions.True(len(result) > 1, "Result should have header + payload")
	})

	t.Run("reused buffer is reset", func(t *testing.T) {
		assertions := assert.New(t)
		buffer := &bytes.Buffer{}
		first := map[string]interface{}{"k1": "v1"}
		second := map[string]interface{}{"k2": "v2"}

		firstPacked, err := PackMap(first, Custom, buffer)
		assertions.NoError(err)
		secondPacked, err := PackMap(second, Custom, buffer)
		assertions.NoError(err)

		firstUnpacked, err := UnpackMap(firstPacked)
		assertions.NoError(err)
		secondUnpacked, err := UnpackMap(secondPacked)
		assertions.NoError(err)

		assertions.Equal("v1", firstUnpacked["k1"])
		assertions.Equal("v2", secondUnpacked["k2"])
	})

	t.Run("too large key length", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{
			strings.Repeat("k", maxUint16+1): "value",
		}
		_, err := PackMap(data, Custom, nil)
		assertions.ErrorIs(err, utils.ErrDataTooLarge)
	})

	t.Run("unsupported codec type", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{"key": "value"}

		_, err := PackMap(data, Encoder(0x0F), nil)
		assertions.ErrorIs(err, utils.ErrUnsupportedCodec)
	})
}

func TestPackArray(t *testing.T) {
	t.Run("empty array", func(t *testing.T) {
		assertions := assert.New(t)
		var data []interface{}
		result, err := PackArray(data, Custom, nil)
		assertions.NoError(err)

		assertions.NotNil(result, "Result should not be nil")
		assertions.NotEmpty(result, "Result should not be empty for empty array")

		unpacked, err := UnpackArray(result)
		assertions.NoError(err, "Should unpack empty array without error")
		assertions.Empty(unpacked, "Unpacked array should be empty")
	})

	t.Run("simple array", func(t *testing.T) {
		assertions := assert.New(t)
		data := []interface{}{
			"string",
			int64(42),
			true,
			3.14,
		}

		packed, err := PackArray(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackArray(packed)
		assertions.NoError(err, "Unpack should not fail")

		assertions.Len(unpacked, 4, "Should have 4 elements")
		assertions.Equal("string", unpacked[0], "Element 0 should match")
		assertions.Equal(int8(42), unpacked[1], "Element 1 should match")
		assertions.Equal(true, unpacked[2], "Element 2 should match")
	})

	t.Run("nested arrays", func(t *testing.T) {
		assertions := assert.New(t)
		data := []interface{}{
			[]interface{}{1, 2, 3},
			[]interface{}{"a", "b", "c"},
		}

		packed, err := PackArray(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackArray(packed)
		assertions.NoError(err, "Unpack should not fail")

		assertions.Len(unpacked, 2, "Should have 2 elements")

		arr1, ok := unpacked[0].([]interface{})
		assertions.True(ok, "First element should be an array")
		assertions.Len(arr1, 3, "First array should have 3 elements")

		arr2, ok := unpacked[1].([]interface{})
		assertions.True(ok, "Second element should be an array")
		assertions.Len(arr2, 3, "Second array should have 3 elements")
	})

	t.Run("array with maps", func(t *testing.T) {
		assertions := assert.New(t)
		data := []interface{}{
			map[string]interface{}{"key1": "value1"},
			map[string]interface{}{"key2": "value2"},
		}

		packed, err := PackArray(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackArray(packed)
		assertions.NoError(err, "Unpack should not fail")

		assertions.Len(unpacked, 2, "Should have 2 elements")

		map1, ok := unpacked[0].(map[string]interface{})
		assertions.True(ok, "First element should be a map")
		assertions.Equal("value1", map1["key1"], "Map1 value should match")
	})

	t.Run("with buffer", func(t *testing.T) {
		assertions := assert.New(t)
		data := []interface{}{"test"}
		buffer := &bytes.Buffer{}

		result, err := PackArray(data, Custom, buffer)
		assertions.NoError(err)

		assertions.Equal(byte(0x00), result[0], "Header should be 0x00 for Custom")
		assertions.True(len(result) > 1, "Result should have header + payload")
	})

	t.Run("byte array", func(t *testing.T) {
		assertions := assert.New(t)
		data := []interface{}{[]byte{1, 2, 3, 4}}
		packed, err := PackArray(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackArray(packed)
		assertions.NoError(err)

		assertions.Len(unpacked, 1)
		assertions.Equal([]byte{1, 2, 3, 4}, unpacked[0])
	})

	t.Run("typed slice support", func(t *testing.T) {
		assertions := assert.New(t)
		data := []interface{}{[]int{1, 2, 3}}
		packed, err := PackArray(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackArray(packed)
		assertions.NoError(err)

		arr, ok := unpacked[0].([]interface{})
		assertions.True(ok)
		assertions.Equal(int8(1), arr[0])
		assertions.Equal(int8(2), arr[1])
		assertions.Equal(int8(3), arr[2])
	})
}

// ============== Unpacker Functions Tests ==============

func TestUnpackMap(t *testing.T) {
	t.Run("invalid data", func(t *testing.T) {
		assertions := assert.New(t)

		_, err := UnpackMap([]byte{})
		assertions.Error(err, "Should error on empty data")
	})

	t.Run("corrupted data recovery", func(t *testing.T) {
		assertions := assert.New(t)
		// Create corrupted data: header byte 0x00 (Custom) + garbage
		corruptedData := []byte{0x00, 255, 255, 255, 255, 255}

		result, err := UnpackMap(corruptedData)
		assertions.ErrorIs(err, utils.ErrUnpackFailed, "Should return ErrUnpackFailed for corrupted data")
		assertions.Nil(result, "Result should be nil for corrupted data")
	})

	t.Run("complex structure", func(t *testing.T) {
		assertions := assert.New(t)
		original := map[string]interface{}{
			"nested": map[string]interface{}{
				"array": []interface{}{1, 2, 3},
				"value": "test",
			},
			"top": "level",
		}

		packed, err := PackMap(original, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err, "Unpack should not fail")

		assertions.Equal("level", unpacked["top"], "Top level value should match")

		nested, ok := unpacked["nested"].(map[string]interface{})
		assertions.True(ok, "Nested should be a map")
		assertions.Equal("test", nested["value"], "Nested value should match")
	})
}

func TestUnpackArray(t *testing.T) {
	t.Run("invalid array data", func(t *testing.T) {
		assertions := assert.New(t)

		_, err := UnpackArray([]byte{})
		assertions.Error(err, "Should error on empty data")
	})

	t.Run("corrupted array data recovery", func(t *testing.T) {
		assertions := assert.New(t)
		// Create corrupted data: header byte 0x00 (Custom) + garbage
		corruptedData := []byte{0x00, 255, 255, 255, 255, 255}

		result, err := UnpackArray(corruptedData)
		assertions.ErrorIs(err, utils.ErrUnpackFailed, "Should return ErrUnpackFailed for corrupted data")
		assertions.Nil(result, "Result should be nil for corrupted data")
	})

	t.Run("mixed array types", func(t *testing.T) {
		assertions := assert.New(t)
		original := []interface{}{
			nil,
			true,
			int8(8),
			int16(16),
			int32(32),
			int64(64),
			uint8(8),
			uint16(16),
			uint32(32),
			uint64(64),
			float32(32.0),
			64.0,
			"string",
			[]interface{}{1, 2},
			map[string]interface{}{"k": "v"},
		}

		packed, err := PackArray(original, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackArray(packed)
		assertions.NoError(err, "Unpack should not fail")

		assertions.Len(unpacked, len(original), "Length should match")

		// Check specific values
		assertions.Nil(unpacked[0], "Element 0 should be nil")
		assertions.Equal(true, unpacked[1], "Element 1 should be true")
		assertions.Equal("string", unpacked[12], "Element 12 should be 'string'")
	})
}

// ============== Helper Functions Tests ==============

func TestGetDataType(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected DataType
	}{
		{"nil", nil, Invalid},
		{"int8", int8(1), Int8},                            // Small value fits in Int8
		{"int16", int16(1), Int8},                          // Small value fits in Int8
		{"int32", int32(1), Int8},                          // Small value fits in Int8
		{"int64", int64(1), Int8},                          // Small value fits in Int8
		{"int", 1, Int8},                                   // Small value fits in Int8
		{"int16 large", int16(32000), Int16},               // Value requires Int16
		{"int24 range", int32(8000000), Int24},             // Value requires Int24
		{"int32 large", int32(100000000), Int32},           // Value requires Int32
		{"int40 range", int64(500000000000), Int40},        // Value requires Int40
		{"int48 range", int64(50000000000000), Int48},      // Value requires Int48
		{"int56 range", int64(5000000000000000), Int56},    // Value requires Int56
		{"int64 large", int64(5000000000000000000), Int64}, // Value requires Int64
		{"float32", float32(1.0), Float32},
		{"float64", 1.0, Float64},
		{"string", "test", String},
		{"bool true", true, Boolean},
		{"bool false", false, Boolean},
		{"byte array", []byte{1, 2, 3}, ByteArray},
		{"slice", []interface{}{1, 2}, Array},
		{"map", map[string]interface{}{"k": "v"}, Map},
		{"unknown type", struct{}{}, Invalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			result := getDataType(tt.input)
			assertions.Equal(tt.expected, result, "getDataType(%T) should return expected type", tt.input)
		})
	}
}

func TestGetDataTypeINT(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		expected DataType
	}{
		{"zero", 0, Int8},
		{"small_positive", 100, Int8},
		{"small_negative", -100, Int8},
		{"int8_max", 127, Int8},
		{"int8_min", -128, Int8},
		{"int16_range", 1000, Int16},
		{"int16_max", 32767, Int16},
		{"int16_min", -32768, Int16},
		{"int24_range", 100000, Int24},
		{"int24_max", MaxInt24, Int24},
		{"int24_min", MinInt24, Int24},
		{"int32_range", 100000000, Int32},
		{"int32_max", math.MaxInt32, Int32},
		{"int32_min", math.MinInt32, Int32},
		{"int40_range", 500000000000, Int40},
		{"int48_range", 50000000000000, Int48},
		{"int56_range", 5000000000000000, Int56},
		{"int64_large", math.MaxInt64, Int64},
		{"int64_negative", math.MinInt64, Int64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			result := GetDataTypeINT(tt.value)
			assertions.Equal(tt.expected, result, "GetDataTypeINT(%d) should return expected type", tt.value)
		})
	}
}

func TestToINT(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected int
	}{
		{"string_number", "123", 123},
		{"string_negative", "-456", -456},
		{"string_zero", "0", 0},
		{"uint", uint(123), 123},
		{"uint8", uint8(123), 123},
		{"uint16", uint16(123), 123},
		{"uint32", uint32(123), 123},
		{"uint64", uint64(123), 123},
		{"int", 123, 123},
		{"int8", int8(123), 123},
		{"int16", int16(123), 123},
		{"int32", int32(123), 123},
		{"int64", int64(123), 123},
		{"float32", float32(123.7), 123},
		{"float64", 123.7, 123},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			result := ToINT(tt.value)
			assertions.Equal(tt.expected, result, "ToINT(%v) should return expected value", tt.value)
		})
	}
}

func TestEncodeValues(t *testing.T) {
	t.Run("nil value", func(t *testing.T) {
		assertions := assert.New(t)
		buffer := &bytes.Buffer{}
		assertions.NoError(encodeValues(nil, Invalid, buffer))

		assertions.Zero(buffer.Len(), "Buffer should be empty for nil value")
	})

	t.Run("boolean", func(t *testing.T) {
		assertions := assert.New(t)
		buffer := &bytes.Buffer{}
		assertions.NoError(encodeValues(true, Boolean, buffer))

		assertions.Equal(1, buffer.Len(), "Buffer should have 1 byte")
		assertions.Equal(byte(1), buffer.Bytes()[0], "True should encode as 1")

		buffer.Reset()
		assertions.NoError(encodeValues(false, Boolean, buffer))

		assertions.Equal(1, buffer.Len(), "Buffer should have 1 byte")
		assertions.Equal(byte(0), buffer.Bytes()[0], "False should encode as 0")
	})

	t.Run("string", func(t *testing.T) {
		assertions := assert.New(t)
		buffer := &bytes.Buffer{}
		value := "hello"
		assertions.NoError(encodeValues(value, String, buffer))

		// 4 bytes for length + 5 bytes for "hello"
		assertions.Equal(9, buffer.Len(), "Buffer should have 9 bytes for 'hello'")

		// Check empty string
		buffer.Reset()
		assertions.NoError(encodeValues("", String, buffer))

		assertions.Equal(4, buffer.Len(), "Buffer should have 4 bytes for empty string")
	})

	t.Run("error on invalid type", func(t *testing.T) {
		assertions := assert.New(t)
		buffer := &bytes.Buffer{}
		err := encodeValues("test", DataType(255), buffer)
		assertions.ErrorIs(err, utils.ErrUnsupportedDataType)
	})
}

// ============== Round Trip Tests ==============

func TestIntegerSizeOptimization(t *testing.T) {
	t.Run("verifies correct size selection", func(t *testing.T) {
		// Test that values are encoded with minimum required size
		testCases := []struct {
			name     string
			value    int
			expected DataType
			size     int // Expected size in bytes after encoding
		}{
			// Int8 range tests
			{"int8_zero", 0, Int8, 1},
			{"int8_max", 127, Int8, 1},
			{"int8_min", -128, Int8, 1},

			// Int16 range tests
			{"int16_small", 128, Int16, 2},
			{"int16_max", 32767, Int16, 2},
			{"int16_min", -32768, Int16, 2},

			// Int24 range tests
			{"int24_small", 32768, Int24, 3},
			{"int24_max", MaxInt24, Int24, 3},
			{"int24_min", MinInt24, Int24, 3},

			// Int32 range tests
			{"int32_small", MaxInt24 + 1, Int32, 4},
			{"int32_max", math.MaxInt32, Int32, 4},
			{"int32_min", math.MinInt32, Int32, 4},

			// Int40 range tests
			{"int40_small", math.MaxInt32 + 1, Int40, 5},
			{"int40_max", MaxInt40, Int40, 5},
			{"int40_min", MinInt40, Int40, 5},

			// Int48 range tests
			{"int48_small", MaxInt40 + 1, Int48, 6},
			{"int48_max", MaxInt48, Int48, 6},
			{"int48_min", MinInt48, Int48, 6},

			// Int56 range tests
			{"int56_small", MaxInt48 + 1, Int56, 7},
			{"int56_max", MaxInt56, Int56, 7},
			{"int56_min", MinInt56, Int56, 7},

			// Int64 range tests
			{"int64_small", MaxInt56 + 1, Int64, 8},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				assertions := assert.New(t)
				// Check that GetDataTypeINT returns the expected type
				dataType := GetDataTypeINT(tc.value)
				assertions.Equal(tc.expected, dataType, "GetDataTypeINT(%d) should return %v", tc.value, tc.expected)

				// Encode the value and verify size
				data := map[string]interface{}{"test": tc.value}
				packed, err := PackMap(data, Custom, nil)
				assertions.NoError(err)

				// The packed data includes:
				// 1 byte for header
				// 2 bytes for map length
				// 1 byte for Map type marker
				// 2 bytes for key length
				// 4 bytes for "test" string
				// 1 byte for value type marker
				// tc.size bytes for the actual value
				expectedSize := 1 + 2 + 1 + 2 + 4 + 1 + tc.size
				assertions.Equal(expectedSize, len(packed),
					"Packed size for value %d should be %d bytes (value takes %d bytes)",
					tc.value, expectedSize, tc.size)

				// Verify round-trip
				unpacked, err := UnpackMap(packed)
				assertions.NoError(err)

				// Convert the unpacked value to int for comparison
				// The unpacker returns the specific type (int8, int16, int32, int64) based on encoding
				var actualValue int
				switch v := unpacked["test"].(type) {
				case int8:
					actualValue = int(v)
				case int16:
					actualValue = int(v)
				case int32:
					actualValue = int(v)
				case int64:
					actualValue = int(v)
				default:
					assertions.Fail("Unexpected type %T for value", v)
				}

				assertions.Equal(tc.value, actualValue, "Value should be preserved after round-trip")
			})
		}
	})
}

func TestIntegerEncodingRoundTrip(t *testing.T) {
	t.Run("all integer sizes", func(t *testing.T) {
		assertions := assert.New(t)

		// Test data with various integer values that require different sizes
		data := map[string]interface{}{
			"tiny":       int8(100),                  // Fits in Int8
			"small":      int16(1000),                // Fits in Int16
			"medium":     int32(100000),              // Requires Int24
			"large":      int32(100000000),           // Requires Int32
			"larger":     int64(500000000000),        // Requires Int40
			"huge":       int64(50000000000000),      // Requires Int48
			"bigger":     int64(5000000000000000),    // Requires Int56
			"massive":    int64(5000000000000000000), // Requires Int64
			"negative8":  int8(-100),
			"negative16": int16(-30000),
			"negative24": int32(-8000000),
			"negative32": int32(-2000000000),
			"negative40": int64(-300000000000),
			"negative48": int64(-30000000000000),
			"negative56": int64(-3000000000000000),
			"negative64": int64(-5000000000000000000),
		}

		// Pack the data
		packed, err := PackMap(data, Custom, nil)
		assertions.NoError(err)

		// Unpack the data
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err, "Unpack should not fail")

		// Verify all values are preserved correctly
		assertions.Equal(int8(100), unpacked["tiny"])
		assertions.Equal(int16(1000), unpacked["small"])
		assertions.Equal(int32(100000), unpacked["medium"])
		assertions.Equal(int32(100000000), unpacked["large"])
		assertions.Equal(int64(500000000000), unpacked["larger"])
		assertions.Equal(int64(50000000000000), unpacked["huge"])
		assertions.Equal(int64(5000000000000000), unpacked["bigger"])
		assertions.Equal(int64(5000000000000000000), unpacked["massive"])
		assertions.Equal(int8(-100), unpacked["negative8"])
		assertions.Equal(int16(-30000), unpacked["negative16"])
		assertions.Equal(int32(-8000000), unpacked["negative24"])
		assertions.Equal(int32(-2000000000), unpacked["negative32"])
		assertions.Equal(int64(-300000000000), unpacked["negative40"])
		assertions.Equal(int64(-30000000000000), unpacked["negative48"])
		assertions.Equal(int64(-3000000000000000), unpacked["negative56"])
		assertions.Equal(int64(-5000000000000000000), unpacked["negative64"])
	})
}

func TestRoundTrip(t *testing.T) {
	t.Run("preserve data types", func(t *testing.T) {
		assertions := assert.New(t)
		original := map[string]interface{}{
			"int8":    int8(-128),
			"int16":   int16(-32768),
			"int32":   int32(-2147483648),
			"int64":   int64(-9223372036854775808),
			"uint8":   uint8(255),
			"uint16":  uint16(65535),
			"uint32":  uint32(4294967295),
			"uint64":  uint64(math.MaxInt64),
			"float32": float32(-3.14),
			"float64": -3.141592653589793,
			"string":  "test string with special chars: @#$%^&*()",
			"bool":    false,
			"nil":     nil,
		}

		packed, err := PackMap(original, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err, "Unpack should not fail")

		// Verify all keys exist
		for key := range original {
			assertions.Contains(unpacked, key, "Key %s should exist in unpacked map", key)
		}

		// Check nil value
		assertions.Nil(unpacked["nil"], "Nil value should remain nil")
		assertions.Equal(int16(255), unpacked["uint8"])
		assertions.Equal(int32(65535), unpacked["uint16"])
		assertions.Equal(int64(4294967295), unpacked["uint32"])
		assertions.Equal(int64(math.MaxInt64), unpacked["uint64"])
	})

	t.Run("deep nesting", func(t *testing.T) {
		assertions := assert.New(t)
		// Create 10 levels of nesting
		deepest := map[string]interface{}{"value": "deep"}
		current := deepest
		for i := 0; i < 9; i++ {
			current = map[string]interface{}{
				"level": current,
			}
		}

		packed, err := PackMap(current, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err, "Unpack should not fail")

		// Navigate to deep-est level
		nav := unpacked
		for i := 0; i < 9; i++ {
			inner, ok := nav["level"].(map[string]interface{})
			assertions.True(ok, "Level %d should be a map", i)
			nav = inner
		}

		assertions.Equal("deep", nav["value"], "Deepest value should match")
	})

	t.Run("large data", func(t *testing.T) {
		assertions := assert.New(t)
		// Create large map
		large := make(map[string]interface{})
		for i := 0; i < 1000; i++ {
			key := "key_" + string(rune(i))
			large[key] = i
		}

		packed, err := PackMap(large, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err, "Unpack should not fail")

		assertions.Len(unpacked, 1000, "Should have 1000 elements")
	})

	t.Run("special strings", func(t *testing.T) {
		assertions := assert.New(t)
		original := map[string]interface{}{
			"empty":     "",
			"spaces":    "   ",
			"newline":   "line1\nline2",
			"tab":       "col1\tcol2",
			"unicode":   "Hello 世界 🌍",
			"escape":    "\"quoted\"",
			"null byte": "before\x00after",
			"long":      string(make([]byte, 10000)),
		}

		packed, err := PackMap(original, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err, "Unpack should not fail")

		for key, value := range original {
			assertions.Equal(value, unpacked[key], "String value for key %s should match", key)
		}
	})
}

// ============== DateTime Tests ==============

func TestDateTimeRoundTrip(t *testing.T) {
	t.Run("time.Time in map", func(t *testing.T) {
		assertions := assert.New(t)
		now := time.Now()
		original := map[string]interface{}{
			"created_at": now,
			"name":       "test",
		}

		packed, err := PackMap(original, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		decoded, ok := unpacked["created_at"].(time.Time)
		assertions.True(ok, "decoded value should be time.Time")
		assertions.True(now.Equal(decoded), "time values should be equal")
		assertions.Equal("test", unpacked["name"])
	})

	t.Run("time.Time UTC", func(t *testing.T) {
		assertions := assert.New(t)
		utcTime := time.Date(2024, 1, 15, 10, 30, 0, 123456789, time.UTC)
		original := map[string]interface{}{
			"timestamp": utcTime,
		}

		packed, err := PackMap(original, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		decoded := unpacked["timestamp"].(time.Time)
		assertions.True(utcTime.Equal(decoded))
	})

	t.Run("time.Time with timezone", func(t *testing.T) {
		assertions := assert.New(t)
		loc := time.FixedZone("IST", 5*3600+30*60)
		istTime := time.Date(2024, 6, 15, 14, 30, 45, 0, loc)
		original := map[string]interface{}{
			"event_time": istTime,
		}

		packed, err := PackMap(original, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		decoded := unpacked["event_time"].(time.Time)
		assertions.True(istTime.Equal(decoded))
	})

	t.Run("time.Time zero value", func(t *testing.T) {
		assertions := assert.New(t)
		original := map[string]interface{}{
			"zero_time": time.Time{},
		}

		packed, err := PackMap(original, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		decoded := unpacked["zero_time"].(time.Time)
		assertions.True(decoded.IsZero())
	})

	t.Run("time.Time in array", func(t *testing.T) {
		assertions := assert.New(t)
		t1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		t2 := time.Date(2025, 6, 15, 12, 30, 0, 0, time.UTC)
		original := []interface{}{t1, t2}

		packed, err := PackArray(original, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackArray(packed)
		assertions.NoError(err)

		assertions.Len(unpacked, 2)
		assertions.True(t1.Equal(unpacked[0].(time.Time)))
		assertions.True(t2.Equal(unpacked[1].(time.Time)))
	})
}

// ============== DateTimeDuration Tests ==============

func TestDateTimeDurationRoundTrip(t *testing.T) {
	t.Run("time.Duration in map", func(t *testing.T) {
		assertions := assert.New(t)
		original := map[string]interface{}{
			"timeout":  5 * time.Second,
			"interval": 100 * time.Millisecond,
			"name":     "config",
		}

		packed, err := PackMap(original, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		decoded, ok := unpacked["timeout"].(time.Duration)
		assertions.True(ok, "decoded value should be time.Duration")
		assertions.Equal(5*time.Second, decoded)

		decoded2, ok := unpacked["interval"].(time.Duration)
		assertions.True(ok, "decoded value should be time.Duration")
		assertions.Equal(100*time.Millisecond, decoded2)

		assertions.Equal("config", unpacked["name"])
	})

	t.Run("time.Duration zero value", func(t *testing.T) {
		assertions := assert.New(t)
		original := map[string]interface{}{
			"zero_dur": time.Duration(0),
		}

		packed, err := PackMap(original, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		decoded := unpacked["zero_dur"].(time.Duration)
		assertions.Equal(time.Duration(0), decoded)
	})

	t.Run("time.Duration negative value", func(t *testing.T) {
		assertions := assert.New(t)
		original := map[string]interface{}{
			"neg_dur": -30 * time.Second,
		}

		packed, err := PackMap(original, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		decoded := unpacked["neg_dur"].(time.Duration)
		assertions.Equal(-30*time.Second, decoded)
	})

	t.Run("time.Duration large value", func(t *testing.T) {
		assertions := assert.New(t)
		largeDur := 24 * time.Hour * 365 // ~1 year
		original := map[string]interface{}{
			"large_dur": largeDur,
		}

		packed, err := PackMap(original, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		decoded := unpacked["large_dur"].(time.Duration)
		assertions.Equal(largeDur, decoded)
	})

	t.Run("time.Duration nanosecond precision", func(t *testing.T) {
		assertions := assert.New(t)
		original := map[string]interface{}{
			"nano":  time.Duration(1),                                          // 1 nanosecond
			"micro": time.Microsecond,                                          // 1 microsecond
			"mixed": 2*time.Hour + 30*time.Minute + 15*time.Second + 123456789, // complex duration
		}

		packed, err := PackMap(original, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		assertions.Equal(time.Duration(1), unpacked["nano"].(time.Duration))
		assertions.Equal(time.Microsecond, unpacked["micro"].(time.Duration))
		assertions.Equal(2*time.Hour+30*time.Minute+15*time.Second+123456789, unpacked["mixed"].(time.Duration))
	})

	t.Run("time.Duration in array", func(t *testing.T) {
		assertions := assert.New(t)
		d1 := 5 * time.Second
		d2 := 100 * time.Millisecond
		d3 := -1 * time.Minute
		original := []interface{}{d1, d2, d3}

		packed, err := PackArray(original, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackArray(packed)
		assertions.NoError(err)

		assertions.Len(unpacked, 3)
		assertions.Equal(d1, unpacked[0].(time.Duration))
		assertions.Equal(d2, unpacked[1].(time.Duration))
		assertions.Equal(d3, unpacked[2].(time.Duration))
	})

	t.Run("time.Duration mixed with time.Time", func(t *testing.T) {
		assertions := assert.New(t)
		now := time.Now()
		original := map[string]interface{}{
			"created_at": now,
			"ttl":        30 * time.Minute,
		}

		packed, err := PackMap(original, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		decodedTime := unpacked["created_at"].(time.Time)
		assertions.True(now.Equal(decodedTime))

		decodedDur := unpacked["ttl"].(time.Duration)
		assertions.Equal(30*time.Minute, decodedDur)
	})
}

// ============== Edge Cases Tests ==============

func TestEdgeCases(t *testing.T) {
	t.Run("map with empty string keys", func(t *testing.T) {
		assertions := assert.New(t)
		original := map[string]interface{}{
			"":     "empty key",
			"key1": "value1",
		}

		packed, err := PackMap(original, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err, "Unpack should not fail")

		assertions.Equal("empty key", unpacked[""], "Empty key value should match")
	})

	t.Run("array with single nil", func(t *testing.T) {
		assertions := assert.New(t)
		original := []interface{}{nil}

		packed, err := PackArray(original, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackArray(packed)
		assertions.NoError(err, "Unpack should not fail")

		assertions.Len(unpacked, 1, "Should have one element")
		assertions.Nil(unpacked[0], "Element should be nil")
	})

	t.Run("map with all nil values", func(t *testing.T) {
		assertions := assert.New(t)
		original := map[string]interface{}{
			"key1": nil,
			"key2": nil,
			"key3": nil,
		}

		packed, err := PackMap(original, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err, "Unpack should not fail")

		for key := range original {
			assertions.Nil(unpacked[key], "Value for key %s should be nil", key)
		}
	})

	t.Run("cyclic reference prevention", func(t *testing.T) {
		assertions := assert.New(t)
		// Note: Direct cyclic references aren't possible with interface{} maps,
		// but we can test deeply nested structures
		depth := 100
		current := map[string]interface{}{"value": "deep"}

		for i := 0; i < depth; i++ {
			current = map[string]interface{}{"nested": current}
		}

		packed, err := PackMap(current, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err, "Unpack should not fail with deep nesting")

		// Navigate to verify structure
		nav := unpacked
		for i := 0; i < depth; i++ {
			inner, ok := nav["nested"].(map[string]interface{})
			assertions.True(ok, "Navigation should succeed at depth %d", i)
			nav = inner
		}

		assertions.Equal("deep", nav["value"], "Deep value should match")
	})
}

// ============== MsgPack Tests ==============

func TestMsgPackMap(t *testing.T) {
	t.Run("simple map", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{
			"name":   "test",
			"count":  int64(42),
			"active": true,
		}

		packed, err := PackMap(data, MsgPack, nil)
		assertions.NoError(err)

		// Verify header byte
		assertions.Equal(byte(0x01), packed[0], "Header should indicate MsgPack")

		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		assertions.Equal("test", unpacked["name"])
		// MsgPack may decode integers differently, so check value not exact type
		assertions.EqualValues(42, unpacked["count"])
		assertions.Equal(true, unpacked["active"])
	})

	t.Run("nested map", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]interface{}{
					"value": "nested",
				},
			},
		}

		packed, err := PackMap(data, MsgPack, nil)
		assertions.NoError(err)

		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		level1, ok := unpacked["level1"].(map[string]interface{})
		assertions.True(ok)
		level2, ok := level1["level2"].(map[string]interface{})
		assertions.True(ok)
		assertions.Equal("nested", level2["value"])
	})

	t.Run("empty map", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{}

		packed, err := PackMap(data, MsgPack, nil)
		assertions.NoError(err)

		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)
		assertions.Empty(unpacked)
	})
}

func TestMsgPackArray(t *testing.T) {
	t.Run("simple array", func(t *testing.T) {
		assertions := assert.New(t)
		data := []interface{}{"hello", int64(42), true, 3.14}

		packed, err := PackArray(data, MsgPack, nil)
		assertions.NoError(err)

		assertions.Equal(byte(0x01), packed[0], "Header should indicate MsgPack")

		unpacked, err := UnpackArray(packed)
		assertions.NoError(err)

		assertions.Len(unpacked, 4)
		assertions.Equal("hello", unpacked[0])
		assertions.EqualValues(42, unpacked[1])
		assertions.Equal(true, unpacked[2])
	})

	t.Run("empty array", func(t *testing.T) {
		assertions := assert.New(t)
		var data []interface{}

		packed, err := PackArray(data, MsgPack, nil)
		assertions.NoError(err)

		unpacked, err := UnpackArray(packed)
		assertions.NoError(err)
		assertions.Empty(unpacked)
	})
}

// ============== Unsupported Codec Unpack Tests ==============

func TestUnpackUnsupportedCodec(t *testing.T) {
	// Encoder 0x0F is not Custom or MsgPack
	invalidHeader := byte(0x0F)
	payload := []byte{invalidHeader, 0x01, 0x02}

	t.Run("unpack map unsupported codec", func(t *testing.T) {
		assertions := assert.New(t)
		_, err := UnpackMap(payload)
		assertions.ErrorIs(err, utils.ErrUnsupportedCodec)
	})

	t.Run("unpack array unsupported codec", func(t *testing.T) {
		assertions := assert.New(t)
		_, err := UnpackArray(payload)
		assertions.ErrorIs(err, utils.ErrUnsupportedCodec)
	})
}

// ============== PackArray Negative Tests ==============

func TestPackArrayNegative(t *testing.T) {
	t.Run("unsupported codec type", func(t *testing.T) {
		assertions := assert.New(t)
		data := []interface{}{"test"}
		_, err := PackArray(data, Encoder(0x0F), nil)
		assertions.ErrorIs(err, utils.ErrUnsupportedCodec)
	})

	t.Run("unsupported type in array", func(t *testing.T) {
		assertions := assert.New(t)
		data := []interface{}{struct{ Name string }{"bad"}}
		_, err := PackArray(data, Custom, nil)
		assertions.ErrorIs(err, utils.ErrUnsupportedDataType)
	})

	t.Run("reused buffer for array", func(t *testing.T) {
		assertions := assert.New(t)
		buffer := &bytes.Buffer{}
		buffer.WriteString("stale data")

		data := []interface{}{"fresh"}
		packed, err := PackArray(data, Custom, buffer)
		assertions.NoError(err)

		unpacked, err := UnpackArray(packed)
		assertions.NoError(err)
		assertions.Len(unpacked, 1)
		assertions.Equal("fresh", unpacked[0])
	})
}

// ============== PackMap Negative Tests ==============

func TestPackMapNegative(t *testing.T) {
	t.Run("unsupported type in map value", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{
			"bad": struct{ X int }{42},
		}
		_, err := PackMap(data, Custom, nil)
		assertions.ErrorIs(err, utils.ErrUnsupportedDataType)
	})

	t.Run("non-string-key map returns invalid", func(t *testing.T) {
		assertions := assert.New(t)
		dt := getDataType(map[int]string{1: "a"})
		assertions.Equal(Invalid, dt)
	})
}

// ============== ToINT64 Negative Tests ==============

func TestToINT64Negative(t *testing.T) {
	t.Run("unsupported type returns error", func(t *testing.T) {
		assertions := assert.New(t)
		_, err := ToINT64("not a number")
		assertions.ErrorIs(err, utils.ErrUnsupportedDataType)
	})

	t.Run("uint64 max overflows int64", func(t *testing.T) {
		assertions := assert.New(t)
		_, err := ToINT64(uint64(math.MaxUint64))
		assertions.ErrorIs(err, utils.ErrValueOutOfRange)
	})
}

// ============== getDataType Additional Tests ==============

func TestGetDataTypeAdditional(t *testing.T) {
	t.Run("uint64 exceeding MaxInt64 returns Invalid", func(t *testing.T) {
		assertions := assert.New(t)
		dt := getDataType(uint64(math.MaxUint64))
		assertions.Equal(Invalid, dt)
	})

	t.Run("typed array via reflect", func(t *testing.T) {
		assertions := assert.New(t)
		dt := getDataType([3]int{1, 2, 3})
		assertions.Equal(Array, dt)
	})

	t.Run("time.Duration detected", func(t *testing.T) {
		assertions := assert.New(t)
		dt := getDataType(5 * time.Second)
		assertions.Equal(DateTimeDuration, dt)
	})

	t.Run("time.Time detected", func(t *testing.T) {
		assertions := assert.New(t)
		dt := getDataType(time.Now())
		assertions.Equal(DateTime, dt)
	})

	t.Run("unsigned integers", func(t *testing.T) {
		assertions := assert.New(t)
		assertions.Equal(Int8, getDataType(uint(50)))
		assertions.Equal(Int16, getDataType(uint(200)))
		assertions.Equal(Int32, getDataType(uint32(100000000)))
	})
}

// ============== encodeValues Additional Tests ==============

func TestEncodeValuesAdditional(t *testing.T) {
	t.Run("encode float32", func(t *testing.T) {
		assertions := assert.New(t)
		buffer := &bytes.Buffer{}
		err := encodeValues(float32(1.5), Float32, buffer)
		assertions.NoError(err)
		assertions.Equal(4, buffer.Len())
	})

	t.Run("encode float64", func(t *testing.T) {
		assertions := assert.New(t)
		buffer := &bytes.Buffer{}
		err := encodeValues(3.14, Float64, buffer)
		assertions.NoError(err)
		assertions.Equal(8, buffer.Len())
	})

	t.Run("encode byte array", func(t *testing.T) {
		assertions := assert.New(t)
		buffer := &bytes.Buffer{}
		err := encodeValues([]byte{0xDE, 0xAD}, ByteArray, buffer)
		assertions.NoError(err)
		assertions.Equal(6, buffer.Len()) // 4 bytes length + 2 bytes data
	})

	t.Run("encode DateTimeDuration", func(t *testing.T) {
		assertions := assert.New(t)
		buffer := &bytes.Buffer{}
		err := encodeValues(5*time.Second, DateTimeDuration, buffer)
		assertions.NoError(err)
		assertions.Equal(8, buffer.Len())
	})

	t.Run("encode DateTime", func(t *testing.T) {
		assertions := assert.New(t)
		buffer := &bytes.Buffer{}
		now := time.Now()
		err := encodeValues(now, DateTime, buffer)
		assertions.NoError(err)
		assertions.True(buffer.Len() > 4) // 4 bytes length + formatted string
	})

	t.Run("encode nested map via encodeValues", func(t *testing.T) {
		assertions := assert.New(t)
		buffer := &bytes.Buffer{}
		inner := map[string]interface{}{"k": "v"}
		err := encodeValues(inner, Map, buffer)
		assertions.NoError(err)
		assertions.True(buffer.Len() > 0)
	})

	t.Run("encode nested array via encodeValues", func(t *testing.T) {
		assertions := assert.New(t)
		buffer := &bytes.Buffer{}
		arr := []interface{}{1, 2}
		err := encodeValues(arr, Array, buffer)
		assertions.NoError(err)
		assertions.True(buffer.Len() > 0)
	})

	t.Run("encode all int widths directly", func(t *testing.T) {
		assertions := assert.New(t)
		intCases := []struct {
			value    interface{}
			dataType DataType
			size     int
		}{
			{int8(1), Int8, 1},
			{int16(200), Int16, 2},
			{int32(100000), Int24, 3},
			{int32(100000000), Int32, 4},
			{int64(500000000000), Int40, 5},
			{int64(50000000000000), Int48, 6},
			{int64(5000000000000000), Int56, 7},
			{int64(5000000000000000000), Int64, 8},
		}

		for _, tc := range intCases {
			buffer := &bytes.Buffer{}
			err := encodeValues(tc.value, tc.dataType, buffer)
			assertions.NoError(err)
			assertions.Equal(tc.size, buffer.Len(), "DataType %d should produce %d bytes", tc.dataType, tc.size)
		}
	})
}

// ============== readValues Default Case Test ==============

func TestReadValuesUnknownType(t *testing.T) {
	assertions := assert.New(t)
	// Craft minimal data; unknown DataType should return nil
	data := []byte{0x00}
	val, pos := readValues(data, 0, DataType(255))
	assertions.Nil(val)
	assertions.Equal(0, pos)
}

// ============== Unpack Panic Recovery Tests ==============

func TestUnpackMapPanicRecovery(t *testing.T) {
	t.Run("truncated custom payload", func(t *testing.T) {
		assertions := assert.New(t)
		// Header says Custom, but payload is too short for a valid map
		data := []byte{0x00, 0x01, 0x00}
		_, err := UnpackMap(data)
		assertions.ErrorIs(err, utils.ErrUnpackFailed)
	})

	t.Run("wrong type tag in map payload", func(t *testing.T) {
		assertions := assert.New(t)
		// Header: Custom, length=1, type=Array (wrong)
		data := []byte{0x00, 0x01, 0x00, byte(Array)}
		_, err := UnpackMap(data)
		assertions.ErrorIs(err, utils.ErrUnpackFailed)
	})
}

func TestUnpackArrayPanicRecovery(t *testing.T) {
	t.Run("truncated custom payload", func(t *testing.T) {
		assertions := assert.New(t)
		data := []byte{0x00, 0x01, 0x00}
		_, err := UnpackArray(data)
		assertions.ErrorIs(err, utils.ErrUnpackFailed)
	})

	t.Run("wrong type tag in array payload", func(t *testing.T) {
		assertions := assert.New(t)
		// Header: Custom, length=1, type=Map (wrong)
		data := []byte{0x00, 0x01, 0x00, byte(Map)}
		_, err := UnpackArray(data)
		assertions.ErrorIs(err, utils.ErrUnpackFailed)
	})
}

// ============== MsgPack Error Tests ==============

func TestMsgPackUnpackErrors(t *testing.T) {
	t.Run("corrupted msgpack map", func(t *testing.T) {
		assertions := assert.New(t)
		header := BuildHeader(MsgPack)
		data := []byte{header, 0xFF, 0xFF, 0xFF}
		_, err := UnpackMap(data)
		assertions.Error(err)
	})

	t.Run("corrupted msgpack array", func(t *testing.T) {
		assertions := assert.New(t)
		header := BuildHeader(MsgPack)
		data := []byte{header, 0xFF, 0xFF, 0xFF}
		_, err := UnpackArray(data)
		assertions.Error(err)
	})
}

// ============== Corner Case: Empty & Nil Data ==============

func TestEmptyAndNilCornerCases(t *testing.T) {
	t.Run("nil data in unpack map", func(t *testing.T) {
		assertions := assert.New(t)
		_, err := UnpackMap(nil)
		assertions.Error(err)
	})

	t.Run("nil data in unpack array", func(t *testing.T) {
		assertions := assert.New(t)
		_, err := UnpackArray(nil)
		assertions.Error(err)
	})

	t.Run("single header byte only map", func(t *testing.T) {
		assertions := assert.New(t)
		_, err := UnpackMap([]byte{0x00})
		assertions.Error(err)
	})

	t.Run("single header byte only array", func(t *testing.T) {
		assertions := assert.New(t)
		_, err := UnpackArray([]byte{0x00})
		assertions.Error(err)
	})

	t.Run("empty string value", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{"empty": ""}
		packed, err := PackMap(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)
		assertions.Equal("", unpacked["empty"])
	})

	t.Run("empty byte array value", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{"empty": []byte{}}
		packed, err := PackMap(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)
		assertions.Equal([]byte{}, unpacked["empty"])
	})

	t.Run("map with only nil values", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{"a": nil, "b": nil}
		packed, err := PackMap(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)
		assertions.Nil(unpacked["a"])
		assertions.Nil(unpacked["b"])
	})

	t.Run("array with only nils", func(t *testing.T) {
		assertions := assert.New(t)
		data := []interface{}{nil, nil, nil}
		packed, err := PackArray(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackArray(packed)
		assertions.NoError(err)
		assertions.Len(unpacked, 3)
		for _, v := range unpacked {
			assertions.Nil(v)
		}
	})
}

// ============== Corner Case: Boundary Integer Values ==============

func TestBoundaryIntegerValues(t *testing.T) {
	t.Run("int8 boundary crossover", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{
			"max_int8":        127,
			"min_int8":        -128,
			"just_above_int8": 128,
			"just_below_int8": -129,
		}
		packed, err := PackMap(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		assertions.Equal(int8(127), unpacked["max_int8"])
		assertions.Equal(int8(-128), unpacked["min_int8"])
		assertions.Equal(int16(128), unpacked["just_above_int8"])
		assertions.Equal(int16(-129), unpacked["just_below_int8"])
	})

	t.Run("int16 boundary crossover", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{
			"max_int16":        32767,
			"min_int16":        -32768,
			"just_above_int16": 32768,
			"just_below_int16": -32769,
		}
		packed, err := PackMap(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		assertions.Equal(int16(32767), unpacked["max_int16"])
		assertions.Equal(int16(-32768), unpacked["min_int16"])
		assertions.Equal(int32(32768), unpacked["just_above_int16"])
		assertions.Equal(int32(-32769), unpacked["just_below_int16"])
	})

	t.Run("int24 boundary crossover", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{
			"max_int24":        MaxInt24,
			"min_int24":        MinInt24,
			"just_above_int24": MaxInt24 + 1,
			"just_below_int24": MinInt24 - 1,
		}
		packed, err := PackMap(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		assertions.Equal(int32(MaxInt24), unpacked["max_int24"])
		assertions.Equal(int32(MinInt24), unpacked["min_int24"])
		assertions.Equal(int32(MaxInt24+1), unpacked["just_above_int24"])
		assertions.Equal(int32(MinInt24-1), unpacked["just_below_int24"])
	})

	t.Run("negative int64", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{
			"min_int64": int64(math.MinInt64),
			"max_int64": int64(math.MaxInt64),
		}
		packed, err := PackMap(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		assertions.Equal(int64(math.MinInt64), unpacked["min_int64"])
		assertions.Equal(int64(math.MaxInt64), unpacked["max_int64"])
	})
}

// ============== Corner Case: Float Special Values ==============

func TestFloatSpecialValues(t *testing.T) {
	t.Run("float32 special values", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{
			"zero":     float32(0),
			"neg_zero": float32(-0),
			"max":      float32(math.MaxFloat32),
			"min_pos":  float32(math.SmallestNonzeroFloat32),
		}
		packed, err := PackMap(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		assertions.Equal(float32(0), unpacked["zero"])
		assertions.Equal(float32(math.MaxFloat32), unpacked["max"])
		assertions.Equal(float32(math.SmallestNonzeroFloat32), unpacked["min_pos"])
	})

	t.Run("float64 special values", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{
			"zero":    0.0,
			"max":     math.MaxFloat64,
			"min_pos": math.SmallestNonzeroFloat64,
			"neg":     -math.MaxFloat64,
		}
		packed, err := PackMap(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		assertions.Equal(0.0, unpacked["zero"])
		assertions.Equal(math.MaxFloat64, unpacked["max"])
		assertions.Equal(math.SmallestNonzeroFloat64, unpacked["min_pos"])
		assertions.Equal(-math.MaxFloat64, unpacked["neg"])
	})

	t.Run("float NaN and Inf in array", func(t *testing.T) {
		assertions := assert.New(t)
		data := []interface{}{
			math.NaN(),
			math.Inf(1),
			math.Inf(-1),
		}
		packed, err := PackArray(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackArray(packed)
		assertions.NoError(err)

		assertions.True(math.IsNaN(unpacked[0].(float64)))
		assertions.True(math.IsInf(unpacked[1].(float64), 1))
		assertions.True(math.IsInf(unpacked[2].(float64), -1))
	})
}

// ============== Corner Case: Large String / Byte Array ==============

func TestLargeStringAndByteArray(t *testing.T) {
	t.Run("large string round-trip", func(t *testing.T) {
		assertions := assert.New(t)
		large := strings.Repeat("x", 100000)
		data := map[string]interface{}{"big": large}
		packed, err := PackMap(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)
		assertions.Equal(large, unpacked["big"])
	})

	t.Run("large byte array round-trip", func(t *testing.T) {
		assertions := assert.New(t)
		large := make([]byte, 100000)
		for i := range large {
			large[i] = byte(i % 256)
		}
		data := map[string]interface{}{"blob": large}
		packed, err := PackMap(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)
		assertions.Equal(large, unpacked["blob"])
	})
}

// ============== Corner Case: Unicode Strings ==============

func TestUnicodeStrings(t *testing.T) {
	t.Run("unicode key and value", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{
			"emoji":              "\U0001F600\U0001F601\U0001F602",
			"chinese":            "\u4F60\u597D\u4E16\u754C",
			"\u00E9\u00E8\u00EA": "accented key",
		}
		packed, err := PackMap(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		assertions.Equal("\U0001F600\U0001F601\U0001F602", unpacked["emoji"])
		assertions.Equal("\u4F60\u597D\u4E16\u754C", unpacked["chinese"])
		assertions.Equal("accented key", unpacked["\u00E9\u00E8\u00EA"])
	})
}

// ============== Corner Case: Deeply Nested Mixed Structures ==============

func TestDeeplyNestedMixed(t *testing.T) {
	t.Run("array in map in array in map", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{
			"outer": []interface{}{
				map[string]interface{}{
					"inner": []interface{}{1, "two", true},
				},
			},
		}
		packed, err := PackMap(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		outer := unpacked["outer"].([]interface{})
		innerMap := outer[0].(map[string]interface{})
		inner := innerMap["inner"].([]interface{})
		assertions.Equal(int8(1), inner[0])
		assertions.Equal("two", inner[1])
		assertions.Equal(true, inner[2])
	})
}

// ============== Corner Case: Typed Map as Nested Value ==============

func TestTypedMapAndSliceNested(t *testing.T) {
	t.Run("map[string]string nested in map", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{
			"headers": map[string]string{
				"Content-Type": "application/json",
				"Accept":       "text/html",
			},
		}
		packed, err := PackMap(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		headers := unpacked["headers"].(map[string]interface{})
		assertions.Equal("application/json", headers["Content-Type"])
		assertions.Equal("text/html", headers["Accept"])
	})

	t.Run("[]string nested in map", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{
			"tags": []string{"go", "codec", "binary"},
		}
		packed, err := PackMap(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackMap(packed)
		assertions.NoError(err)

		tags := unpacked["tags"].([]interface{})
		assertions.Equal("go", tags[0])
		assertions.Equal("codec", tags[1])
		assertions.Equal("binary", tags[2])
	})

	t.Run("[]float64 nested in array", func(t *testing.T) {
		assertions := assert.New(t)
		data := []interface{}{[]float64{1.1, 2.2, 3.3}}
		packed, err := PackArray(data, Custom, nil)
		assertions.NoError(err)
		unpacked, err := UnpackArray(packed)
		assertions.NoError(err)

		inner := unpacked[0].([]interface{})
		assertions.Equal(1.1, inner[0])
		assertions.Equal(2.2, inner[1])
		assertions.Equal(3.3, inner[2])
	})
}

// ============== encodeValues Error Paths via ToINT64 Failure ==============

func TestEncodeValuesIntErrorPaths(t *testing.T) {
	// Force ToINT64 error by passing a non-numeric type with an integer DataType.
	// This exercises the error returns inside each Int case of encodeValues.
	intTypes := []DataType{Int8, Int16, Int24, Int32, Int40, Int48, Int56, Int64}
	for _, dt := range intTypes {
		t.Run("error_"+string(rune('0'+dt)), func(t *testing.T) {
			assertions := assert.New(t)
			buffer := &bytes.Buffer{}
			err := encodeValues("not_a_number", dt, buffer)
			assertions.ErrorIs(err, utils.ErrUnsupportedDataType)
		})
	}
}

// ============== packArrayData / packMapData Error Propagation ==============

func TestPackArrayDataErrorPropagation(t *testing.T) {
	t.Run("unsupported type error propagates from encodeValues", func(t *testing.T) {
		assertions := assert.New(t)
		// A complex number cannot be encoded
		data := []interface{}{complex(1, 2)}
		_, err := PackArray(data, Custom, nil)
		assertions.ErrorIs(err, utils.ErrUnsupportedDataType)
	})
}

func TestPackMapDataErrorPropagation(t *testing.T) {
	t.Run("encode error propagates with key context", func(t *testing.T) {
		assertions := assert.New(t)
		data := map[string]interface{}{
			"bad_value": complex(1, 2),
		}
		_, err := PackMap(data, Custom, nil)
		assertions.ErrorIs(err, utils.ErrUnsupportedDataType)
	})
}

// ============== readValues: DateTime parse error ==============

func TestReadValuesDateTimeParseError(t *testing.T) {
	t.Run("invalid datetime string triggers panic recovery", func(t *testing.T) {
		assertions := assert.New(t)
		// Manually craft a map payload with a DateTime type tag and invalid date string
		// Format: [header][length_lo][length_hi][Map_tag][key_len_lo][key_len_hi][key][DateTime_tag][str_len(4 bytes)][invalid_date_string]
		buf := &bytes.Buffer{}
		// Map length = 1
		WriteInt16Value(1, buf)
		// Map type tag
		buf.WriteByte(byte(Map))
		// Key "d" (length=1)
		WriteInt16Value(1, buf)
		buf.WriteString("d")
		// Value type: DateTime
		buf.WriteByte(byte(DateTime))
		// String length (4 bytes): 7
		WriteInt32Value(7, buf)
		// Invalid date string
		buf.WriteString("INVALID")

		header := BuildHeader(Custom)
		data := append([]byte{header}, buf.Bytes()...)

		_, err := UnpackMap(data)
		assertions.ErrorIs(err, utils.ErrUnpackFailed)
	})
}

// ============== encodeValues: string/bytearray length validation ==============

func TestEncodeValuesStringEdgeCases(t *testing.T) {
	t.Run("empty string", func(t *testing.T) {
		assertions := assert.New(t)
		buffer := &bytes.Buffer{}
		err := encodeValues("", String, buffer)
		assertions.NoError(err)
		assertions.Equal(4, buffer.Len()) // just the 4-byte length prefix
	})

	t.Run("empty byte array", func(t *testing.T) {
		assertions := assert.New(t)
		buffer := &bytes.Buffer{}
		err := encodeValues([]byte{}, ByteArray, buffer)
		assertions.NoError(err)
		assertions.Equal(4, buffer.Len()) // just the 4-byte length prefix
	})
}

// ============== PackMap / PackArray: MsgPack error path ==============

func TestPackMapMsgPackUnsupportedType(t *testing.T) {
	assertions := assert.New(t)
	// Channels can't be serialized by msgpack
	ch := make(chan int)
	data := map[string]interface{}{"ch": ch}
	_, err := PackMap(data, MsgPack, nil)
	assertions.Error(err)
}

func TestPackArrayMsgPackUnsupportedType(t *testing.T) {
	assertions := assert.New(t)
	ch := make(chan int)
	data := []interface{}{ch}
	_, err := PackArray(data, MsgPack, nil)
	assertions.Error(err)
}

// ============== Size Limit Tests ==============

func TestPackMapTooManyEntries(t *testing.T) {
	assertions := assert.New(t)
	data := make(map[string]interface{}, maxUint16+1)
	for i := 0; i <= maxUint16; i++ {
		data[strings.Repeat("k", 1)+string(rune(i/256))+string(rune(i%256))] = i
	}
	_, err := PackMap(data, Custom, nil)
	assertions.ErrorIs(err, utils.ErrDataTooLarge)
}
