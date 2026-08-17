package codec

import (
	"bytes"
	"testing"
)

// ============== Benchmark Tests ==============

// ============== Size Optimization Benchmarks ==============

func BenchmarkPackMapInt8Values(b *testing.B) {
	// All values fit in Int8
	data := map[string]interface{}{
		"a": 0,
		"b": 100,
		"c": -100,
		"d": 127,
		"e": -128,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = PackMap(data, Custom, nil)
	}
}

func BenchmarkPackMapInt16Values(b *testing.B) {
	// All values require Int16
	data := map[string]interface{}{
		"a": 200,
		"b": 1000,
		"c": -1000,
		"d": 32767,
		"e": -32768,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = PackMap(data, Custom, nil)
	}
}

func BenchmarkPackMapInt24Values(b *testing.B) {
	// All values require Int24
	data := map[string]interface{}{
		"a": 40000,
		"b": 100000,
		"c": -100000,
		"d": MaxInt24,
		"e": MinInt24,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = PackMap(data, Custom, nil)
	}
}

func BenchmarkPackMapInt32Values(b *testing.B) {
	// All values require Int32
	data := map[string]interface{}{
		"a": 10000000,
		"b": 100000000,
		"c": -100000000,
		"d": 2147483647,
		"e": -2147483648,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = PackMap(data, Custom, nil)
	}
}

func BenchmarkPackMapMixedIntSizes(b *testing.B) {
	// Mixed sizes to test optimization
	data := map[string]interface{}{
		"int8":  100,                 // Int8 (1 byte)
		"int16": 1000,                // Int16 (2 bytes)
		"int24": 100000,              // Int24 (3 bytes)
		"int32": 100000000,           // Int32 (4 bytes)
		"int40": 500000000000,        // Int40 (5 bytes)
		"int48": 50000000000000,      // Int48 (6 bytes)
		"int56": 5000000000000000,    // Int56 (7 bytes)
		"int64": 5000000000000000000, // Int64 (8 bytes)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = PackMap(data, Custom, nil)
	}
}

func BenchmarkPackMapSmall(b *testing.B) {
	data := map[string]interface{}{
		"key1": "value1",
		"key2": 42, // Will use Int8 now instead of Int64
		"key3": true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = PackMap(data, Custom, nil)
	}
}

func BenchmarkPackMapLarge(b *testing.B) {
	data := make(map[string]interface{})
	for i := 0; i < 100; i++ {
		key := "key_" + string(rune(i))
		data[key] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = PackMap(data, Custom, nil)
	}
}

func BenchmarkPackArraySmall(b *testing.B) {
	data := []interface{}{"string", int64(42), true, 3.14}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = PackArray(data, Custom, nil)
	}
}

func BenchmarkPackArrayLarge(b *testing.B) {
	data := make([]interface{}, 100)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = PackArray(data, Custom, nil)
	}
}

func BenchmarkUnpackMap(b *testing.B) {
	data := map[string]interface{}{
		"key1": "value1",
		"key2": int64(42),
		"key3": true,
		"nested": map[string]interface{}{
			"inner": "value",
		},
	}
	packed, _ := PackMap(data, Custom, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UnpackMap(packed)
	}
}

func BenchmarkUnpackArray(b *testing.B) {
	data := []interface{}{
		"string",
		int64(42),
		true,
		3.14,
		[]interface{}{1, 2, 3},
	}
	packed, _ := PackArray(data, Custom, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UnpackArray(packed)
	}
}

// ============== Round-Trip Benchmarks ==============

func BenchmarkRoundTripMapInt8(b *testing.B) {
	data := map[string]interface{}{
		"a": 10,
		"b": 50,
		"c": -50,
		"d": 127,
		"e": -128,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		packed, _ := PackMap(data, Custom, nil)
		_, _ = UnpackMap(packed)
	}
}

func BenchmarkRoundTripMapInt24(b *testing.B) {
	data := map[string]interface{}{
		"a": 40000,
		"b": 100000,
		"c": -100000,
		"d": MaxInt24,
		"e": MinInt24,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		packed, _ := PackMap(data, Custom, nil)
		_, _ = UnpackMap(packed)
	}
}

func BenchmarkRoundTripMapInt64(b *testing.B) {
	data := map[string]interface{}{
		"a": 1000000000000000000,
		"b": 5000000000000000000,
		"c": -1000000000000000000,
		"d": 9223372036854775807,
		"e": -9223372036854775808,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		packed, _ := PackMap(data, Custom, nil)
		_, _ = UnpackMap(packed)
	}
}

func BenchmarkRoundTripMapMixed(b *testing.B) {
	data := map[string]interface{}{
		"int8":   100,                 // 1 byte
		"int16":  1000,                // 2 bytes
		"int24":  100000,              // 3 bytes
		"int32":  100000000,           // 4 bytes
		"int40":  500000000000,        // 5 bytes
		"int48":  50000000000000,      // 6 bytes
		"int56":  5000000000000000,    // 7 bytes
		"int64":  5000000000000000000, // 8 bytes
		"string": "test",
		"bool":   true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		packed, _ := PackMap(data, Custom, nil)
		_, _ = UnpackMap(packed)
	}
}

func BenchmarkRoundTripMap(b *testing.B) {
	data := map[string]interface{}{
		"name":   "test",
		"count":  42, // Will use Int8 now
		"active": true,
		"nested": map[string]interface{}{
			"key": "value",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		packed, _ := PackMap(data, Custom, nil)
		_, _ = UnpackMap(packed)
	}
}

func BenchmarkRoundTripArray(b *testing.B) {
	data := []interface{}{
		"string",
		int64(42),
		true,
		3.14,
		[]interface{}{1, 2, 3},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		packed, _ := PackArray(data, Custom, nil)
		_, _ = UnpackArray(packed)
	}
}

func BenchmarkGetDataType(b *testing.B) {
	values := []interface{}{
		nil,
		int8(1),
		int16(1),
		int32(1),
		int64(1),
		uint8(1),
		uint16(1),
		uint32(1),
		uint64(1),
		float32(1.0),
		1.0,
		"string",
		true,
		[]byte{1, 2, 3},
		[]interface{}{1, 2},
		map[string]interface{}{"k": "v"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			getDataType(v)
		}
	}
}

// ============== Integer Write Benchmarks ==============

func BenchmarkWriteInt8(b *testing.B) {
	buffer := &bytes.Buffer{}
	value := int8(127)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Reset()
		WriteInt8Value(value, buffer)
	}
}

func BenchmarkWriteInt16(b *testing.B) {
	buffer := &bytes.Buffer{}
	value := int16(32767)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Reset()
		WriteInt16Value(value, buffer)
	}
}

func BenchmarkWriteInt24(b *testing.B) {
	buffer := &bytes.Buffer{}
	value := int32(8388607) // MaxInt24

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Reset()
		WriteInt24Value(value, buffer)
	}
}

func BenchmarkWriteInt32(b *testing.B) {
	buffer := &bytes.Buffer{}
	value := int32(2147483647)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Reset()
		WriteInt32Value(value, buffer)
	}
}

func BenchmarkWriteInt40(b *testing.B) {
	buffer := &bytes.Buffer{}
	value := int64(549755813887) // MaxInt40

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Reset()
		WriteInt40Value(value, buffer)
	}
}

func BenchmarkWriteInt48(b *testing.B) {
	buffer := &bytes.Buffer{}
	value := int64(140737488355327) // MaxInt48

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Reset()
		WriteInt48Value(value, buffer)
	}
}

func BenchmarkWriteInt56(b *testing.B) {
	buffer := &bytes.Buffer{}
	value := int64(36028797018963967) // MaxInt56

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Reset()
		WriteInt56Value(value, buffer)
	}
}

func BenchmarkWriteInt64(b *testing.B) {
	buffer := &bytes.Buffer{}
	value := int64(9223372036854775807)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Reset()
		WriteInt64Value(value, buffer)
	}
}

// ============== Integer Read Benchmarks ==============

func BenchmarkReadINT16(b *testing.B) {
	data := []byte{0xFF, 0x7F} // MaxInt16

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ReadINT16Value(data)
	}
}

func BenchmarkReadINT24(b *testing.B) {
	data := []byte{0xFF, 0xFF, 0x7F} // MaxInt24

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ReadINT24Value(data)
	}
}

func BenchmarkReadINT32(b *testing.B) {
	data := []byte{0xFF, 0xFF, 0xFF, 0x7F} // MaxInt32

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ReadINT32Value(data)
	}
}

func BenchmarkReadINT40(b *testing.B) {
	data := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x7F} // MaxInt40

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ReadINT40Value(data)
	}
}

func BenchmarkReadINT48(b *testing.B) {
	data := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x7F} // MaxInt48

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ReadINT48Value(data)
	}
}

func BenchmarkReadINT56(b *testing.B) {
	data := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x7F} // MaxInt56

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ReadINT56Value(data)
	}
}

func BenchmarkReadINT64(b *testing.B) {
	data := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x7F} // MaxInt64

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ReadINT64Value(data)
	}
}

// ============== Helper Function Benchmarks ==============

func BenchmarkGetDataTypeINT(b *testing.B) {
	// Test various values that require different integer sizes
	values := []int{
		0,                   // Int8
		100,                 // Int8
		-100,                // Int8
		1000,                // Int16
		-30000,              // Int16
		100000,              // Int24
		-8000000,            // Int24
		100000000,           // Int32
		-2000000000,         // Int32
		500000000000,        // Int40
		-300000000000,       // Int40
		50000000000000,      // Int48
		-30000000000000,     // Int48
		5000000000000000,    // Int56
		-3000000000000000,   // Int56
		5000000000000000000, // Int64
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			_ = GetDataTypeINT(v)
		}
	}
}

func BenchmarkToINT(b *testing.B) {
	// Test various input types
	values := []interface{}{
		"123",
		"-456",
		uint(123),
		uint8(123),
		uint16(123),
		uint32(123),
		uint64(123),
		123,
		int8(123),
		int16(123),
		int32(123),
		int64(123),
		float32(123.7),
		123.7,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			_ = ToINT(v)
		}
	}
}

func BenchmarkWriteFloat64(b *testing.B) {
	buffer := &bytes.Buffer{}
	value := 3.141592653589793

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Reset()
		WriteFloat64Value(value, buffer)
	}
}

// ============== MsgPack Benchmarks ==============

func BenchmarkPackMapMsgPack(b *testing.B) {
	data := map[string]interface{}{
		"key1": "value1",
		"key2": int64(42),
		"key3": true,
		"nested": map[string]interface{}{
			"inner": "value",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = PackMap(data, MsgPack, nil)
	}
}

func BenchmarkPackArrayMsgPack(b *testing.B) {
	data := []interface{}{"string", int64(42), true, 3.14}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = PackArray(data, MsgPack, nil)
	}
}

func BenchmarkRoundTripMapMsgPack(b *testing.B) {
	data := map[string]interface{}{
		"name":   "test",
		"count":  int64(42),
		"active": true,
		"nested": map[string]interface{}{
			"key": "value",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		packed, _ := PackMap(data, MsgPack, nil)
		_, _ = UnpackMap(packed)
	}
}

func BenchmarkUnpackMapMsgPack(b *testing.B) {
	data := map[string]interface{}{
		"key1": "value1",
		"key2": int64(42),
		"key3": true,
		"nested": map[string]interface{}{
			"inner": "value",
		},
	}
	packed, _ := PackMap(data, MsgPack, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = UnpackMap(packed)
	}
}
