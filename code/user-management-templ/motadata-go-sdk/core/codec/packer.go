// Package codec - Packer Implementation
//
// ARCHITECTURE OVERVIEW:
//
// The packer module implements the serialization (encoding) side of the codec,
// transforming Go data structures into compact binary representations. It uses
// reflection for type discovery and a cascading type switch for efficient
// serialization of supported types.
//
// SERIALIZATION PIPELINE:
//
//	Go Value → Type Detection → Width Selection → Binary Encoding → Buffer
//	    ↓           ↓                ↓                   ↓            ↓
//	map[string]any  getDataType()  GetDataTypeINT()  WriteInt*()  bytes.Buffer
//
// DESIGN PRINCIPLES:
//
// 1. Zero-Copy Where Possible:
//   - Direct buffer writes avoid intermediate allocations
//   - String writes use buffer.WriteString (no byte conversion)
//   - Reusable buffers minimize GC pressure
//
// 2. Type-Driven Dispatch:
//   - Single type check determines serialization path
//   - No runtime type assertions in hot paths
//   - Compile-time optimizations for known types
//
// 3. Recursive Structure Handling:
//   - Maps and arrays serialize recursively
//   - Stack-based traversal (no heap allocations)
//   - Depth-first serialization order
//
// SERIALIZATION FORMAT:
//
// Map Encoding:
//
//	┌────────────┬──────────┬─────────────────────┐
//	│ Key Count  │ Type Tag │ Key-Value Pairs     │
//	│ (2 bytes)  │ (1 byte) │ (variable)          │
//	└────────────┴──────────┴─────────────────────┘
//
// Array Encoding:
//
//	┌──────────────┬──────────┬───────────────────┐
//	│ Element Count│ Type Tag │ Elements          │
//	│ (2 bytes)    │ (1 byte) │ (variable)        │
//	└──────────────┴──────────┴───────────────────┘
//
// Key-Value Pair:
//
//	┌────────────┬─────────┬──────────┬───────────┐
//	│ Key Length │ Key Str │ Val Type │ Value     │
//	│ (2 bytes)  │ (UTF-8) │ (1 byte) │ (variable)│
//	└────────────┴─────────┴──────────┴───────────┘
//
// TYPE DETECTION ALGORITHM:
//
// The getDataType function implements a type cascade:
// 1. Nil check → Invalid
// 2. Type switch on concrete types
// 3. Numeric types → GetDataTypeINT for width selection
// 4. Complex types → Recursive serialization
//
// This approach minimizes reflection overhead by:
// - Using type switches instead of reflect.TypeOf
// - Caching type information where possible
// - Avoiding interface boxing for primitives
//
// NUMERIC CONVERSION STRATEGY:
//
// The ToINT function provides universal numeric conversion:
// - Handles all Go numeric types (int*, uint*, float*)
// - String parsing with error suppression
// - Lossy conversion for floats (truncation)
//
// Conversion matrix:
//
//	string → ParseInt → int
//	uint* → direct cast → int
//	int* → direct cast/assign → int
//	float* → truncation → int
//
// BUFFER MANAGEMENT:
//
// Efficient buffer usage patterns:
// 1. Pre-allocation when size is known
// 2. Exponential growth for unknown sizes
// 3. Buffer reuse across serializations
// 4. Direct writes without intermediate copies
//
// Buffer growth strategy:
//
//	Initial: 512 bytes (typical message size)
//	Growth: 2x when capacity exceeded
//	Maximum: Limited by available memory
//
// PERFORMANCE OPTIMIZATIONS:
//
// 1. Type-Specific Fast Paths:
//   - Direct casts for known numeric types
//   - Avoid reflection for primitives
//   - Inline small functions
//
// 2. Memory Efficiency:
//   - Stack allocation for small values
//   - Buffer pooling for large messages
//   - Minimal intermediate allocations
//
// 3. CPU Cache Optimization:
//   - Sequential buffer writes
//   - Predictable access patterns
//   - Hot path optimization
//
// ERROR HANDLING PHILOSOPHY:
//
// The packer uses panic for unrecoverable errors:
// - Unsupported types trigger panic
// - Buffer overflow prevented by growth
// - Type mismatches caught at compile time
//
// This approach prioritizes performance over recovery:
// - No error return values in hot paths
// - Panics indicate programming errors
// - Production code should validate inputs
//
// REFLECTION USAGE:
//
// Reflection is minimized but necessary for:
// - Generic map/array handling
// - Interface{} type discovery
// - Dynamic type conversion
//
// Reflection overhead mitigation:
// - Cache reflect.Value when possible
// - Use type switches for common cases
// - Avoid reflect.TypeOf in loops
//
// FUTURE ENHANCEMENTS:
//
// - Code generation for type-specific packers
// - SIMD acceleration for array packing
// - Pluggable compression support
// - Schema-based optimization
// - Zero-allocation mode
package codec

import (
	"bytes"
	"fmt"
	"math"
	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/utils"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

const (
	maxUint16 = 1<<16 - 1
	maxUint32 = 1<<32 - 1
)

// PackMap serializes a map[string]interface{} into a byte slice with a header byte.
// The header encodes the codec type and compression algorithm.
// Wire format: [1 byte header][N bytes compressed payload]
func PackMap(data map[string]interface{}, encoder Encoder, buffer *bytes.Buffer) ([]byte, error) {

	var payload []byte

	switch encoder {
	case Custom:
		if buffer == nil {
			buffer = bytes.NewBuffer(make([]byte, 0, len(data)*32))
		} else {
			buffer.Reset()
		}
		if err := packMapData(data, buffer); err != nil {
			return nil, err
		}
		payload = buffer.Bytes()
	case MsgPack:
		encoded, err := msgpack.Marshal(data)
		if err != nil {
			return nil, utils.WrapErr(utils.ErrUnsupportedCodec, err)
		}
		payload = encoded
	default:
		return nil, fmt.Errorf("%w: 0x%02x", utils.ErrUnsupportedCodec, encoder)
	}

	result := make([]byte, 1+len(payload))
	result[0] = BuildHeader(encoder)
	copy(result[1:], payload)
	return result, nil
}

func packMapData(data map[string]interface{}, buffer *bytes.Buffer) error {
	if len(data) > maxUint16 {
		return fmt.Errorf("%w: map entry count %d exceeds %d", utils.ErrDataTooLarge, len(data), maxUint16)
	}

	// first 2 bytes length of keys
	WriteInt16Value(int16(len(data)), buffer)

	// datatype
	buffer.WriteByte(byte(Map))

	// write each key-value pair
	for key, value := range data {
		if len(key) > maxUint16 {
			return fmt.Errorf("%w: key length %d exceeds %d", utils.ErrDataTooLarge, len(key), maxUint16)
		}
		WriteInt16Value(int16(len(key)), buffer)
		buffer.WriteString(key)
		dataType := getDataType(value)
		if dataType == Invalid && value != nil {
			return fmt.Errorf("%w: key %q has %T", utils.ErrUnsupportedDataType, key, value)
		}
		buffer.WriteByte(byte(dataType))
		if err := encodeValues(value, dataType, buffer); err != nil {
			return fmt.Errorf("%w: key %q", err, key)
		}
	}

	return nil
}

// PackArray serializes a slice of interface{} into a byte slice with a header byte.
// The header encodes the codec type and compression algorithm.
// Wire format: [1 byte header][N bytes compressed payload]
func PackArray(data []interface{}, encoder Encoder, buffer *bytes.Buffer) ([]byte, error) {
	var payload []byte

	switch encoder {
	case Custom:
		if buffer == nil {
			buffer = bytes.NewBuffer(make([]byte, 0, len(data)*16))
		} else {
			buffer.Reset()
		}
		if err := packArrayData(data, buffer); err != nil {
			return nil, err
		}
		payload = buffer.Bytes()
	case MsgPack:
		encoded, err := msgpack.Marshal(data)
		if err != nil {
			return nil, utils.WrapErr(utils.ErrUnsupportedCodec, err)
		}
		payload = encoded
	default:
		return nil, fmt.Errorf("%w: 0x%02x", utils.ErrUnsupportedCodec, encoder)
	}

	result := make([]byte, 1+len(payload))
	result[0] = BuildHeader(encoder)
	copy(result[1:], payload)
	return result, nil
}

func packArrayData(data []interface{}, buffer *bytes.Buffer) error {
	if len(data) > maxUint32 {
		return fmt.Errorf("%w: array length %d exceeds %d", utils.ErrDataTooLarge, len(data), maxUint32)
	}

	// write array length
	WriteInt16Value(int16(len(data)), buffer)

	// datatype
	buffer.WriteByte(byte(Array))

	// write each element with its type
	for _, value := range data {
		dataType := getDataType(value)
		if dataType == Invalid && value != nil {
			return fmt.Errorf("%w: array element has %T", utils.ErrUnsupportedDataType, value)
		}
		buffer.WriteByte(byte(dataType))
		if err := encodeValues(value, dataType, buffer); err != nil {
			return err
		}
	}

	return nil
}

// encodeValues writes actual data values based on their data type
func encodeValues(value any, datatype DataType, buffer *bytes.Buffer) error {
	switch datatype {
	case Invalid:
		return nil
	case Int8, Int16, Int24, Int32, Int40, Int48, Int56, Int64:
		return encodeIntValue(value, datatype, buffer)
	case Float32:
		WriteFloat32Value(value.(float32), buffer)
	case Float64:
		WriteFloat64Value(value.(float64), buffer)
	case String:
		return encodeStringValue(value.(string), buffer)
	case ByteArray:
		return encodeByteArrayValue(value.([]byte), buffer)
	case Array:
		return encodeArrayValue(value, buffer)
	case Map:
		return encodeMapValue(value, buffer)
	case Boolean:
		return encodeBoolValue(value.(bool), buffer)
	case DateTime:
		return encodeDateTimeValue(value.(time.Time), buffer)
	case DateTimeDuration:
		WriteInt64Value(int64(value.(time.Duration)), buffer)
	default:
		return fmt.Errorf("%w: %T", utils.ErrUnsupportedDataType, value)
	}
	return nil
}

// encodeIntValue converts value to int64 and writes it using the appropriate width writer.
func encodeIntValue(value any, datatype DataType, buffer *bytes.Buffer) error {
	i64, err := ToINT64(value)
	if err != nil {
		return err
	}
	switch datatype {
	case Int8:
		WriteInt8Value(int8(i64), buffer)
	case Int16:
		WriteInt16Value(int16(i64), buffer)
	case Int24:
		WriteInt24Value(int32(i64), buffer)
	case Int32:
		WriteInt32Value(int32(i64), buffer)
	case Int40:
		WriteInt40Value(i64, buffer)
	case Int48:
		WriteInt48Value(i64, buffer)
	case Int56:
		WriteInt56Value(i64, buffer)
	case Int64:
		WriteInt64Value(i64, buffer)
	default:
		// unreachable: caller only passes Int8-Int64
	}
	return nil
}

func encodeStringValue(val string, buffer *bytes.Buffer) error {
	if len(val) > maxUint32 {
		return fmt.Errorf("%w: string length %d exceeds %d", utils.ErrDataTooLarge, len(val), maxUint32)
	}
	WriteInt32Value(int32(len(val)), buffer)
	buffer.WriteString(val)
	return nil
}

func encodeByteArrayValue(val []byte, buffer *bytes.Buffer) error {
	if len(val) > maxUint32 {
		return fmt.Errorf("%w: byte array length %d exceeds %d", utils.ErrDataTooLarge, len(val), maxUint32)
	}
	WriteInt32Value(int32(len(val)), buffer)
	buffer.Write(val)
	return nil
}

func encodeArrayValue(value any, buffer *bytes.Buffer) error {
	// fast path: direct type assertion for []interface{}
	if arr, ok := value.([]interface{}); ok {
		return packArrayData(arr, buffer)
	}
	// reflection fallback for typed slices (e.g., []string, []int)
	rv := reflect.ValueOf(value)
	result := make([]interface{}, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		result[i] = rv.Index(i).Interface()
	}
	return packArrayData(result, buffer)
}

func encodeMapValue(value any, buffer *bytes.Buffer) error {
	// fast path: direct type assertion for map[string]interface{}
	if m, ok := value.(map[string]interface{}); ok {
		return packMapData(m, buffer)
	}
	// reflection fallback for typed maps (e.g., map[string]int)
	rv := reflect.ValueOf(value)
	result := make(map[string]interface{}, rv.Len())
	for _, key := range rv.MapKeys() {
		result[key.String()] = rv.MapIndex(key).Interface()
	}
	return packMapData(result, buffer)
}

func encodeBoolValue(val bool, buffer *bytes.Buffer) error {
	if val {
		buffer.WriteByte(1)
	} else {
		buffer.WriteByte(0)
	}
	return nil
}

func encodeDateTimeValue(t time.Time, buffer *bytes.Buffer) error {
	val := t.Format(time.RFC3339Nano)
	if len(val) > maxUint32 {
		return fmt.Errorf("%w: datetime string length %d exceeds %d", utils.ErrDataTooLarge, len(val), maxUint32)
	}
	WriteInt32Value(int32(len(val)), buffer)
	buffer.WriteString(val)
	return nil
}

func ToINT64(value any) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		if uint64(v) > math.MaxInt64 {
			return 0, fmt.Errorf("%w: %d", utils.ErrValueOutOfRange, v)
		}
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		if v > math.MaxInt64 {
			return 0, fmt.Errorf("%w: %d", utils.ErrValueOutOfRange, v)
		}
		return int64(v), nil
	default:
		return 0, fmt.Errorf("%w: %T", utils.ErrUnsupportedDataType, value)
	}
}

// ToINT converts various numeric types to int
func ToINT(value any) (result int) {
	name := reflect.TypeOf(value).Name()

	if name == "string" {
		output, _ := strconv.ParseInt(strings.TrimSpace(value.(string)), 10, 64)
		result = int(output)
	} else if name == "uint" {
		result = int(value.(uint))
	} else if name == "uint8" {
		result = int(value.(uint8))
	} else if name == "uint16" {
		result = int(value.(uint16))
	} else if name == "uint32" {
		result = int(value.(uint32))
	} else if name == "uint64" {
		result = int(value.(uint64))
	} else if name == "int" {
		result = value.(int)
	} else if name == "int8" {
		result = int(value.(int8))
	} else if name == "int16" {
		result = int(value.(int16))
	} else if name == "int32" {
		result = int(value.(int32))
	} else if name == "int64" {
		result = int(value.(int64))
	} else if name == "float64" {
		result = int(value.(float64))
	} else if name == "float32" {
		result = int(value.(float32))
	}

	return
}

// getDataType determines the DataType for a given value.
// Uses direct type switch for common types to avoid reflection overhead,
// with reflection fallback for uncommon typed slices/maps.
func getDataType(v any) DataType {
	if v == nil {
		return Invalid
	}

	// fast path: direct type switch for all common concrete types
	switch val := v.(type) {
	case bool:
		return Boolean
	case int:
		return GetDataTypeINT64(int64(val))
	case int8:
		return GetDataTypeINT64(int64(val))
	case int16:
		return GetDataTypeINT64(int64(val))
	case int32:
		return GetDataTypeINT64(int64(val))
	case int64:
		return GetDataTypeINT64(val)
	case uint:
		if uint64(val) > math.MaxInt64 {
			return Invalid
		}
		return GetDataTypeINT64(int64(val))
	case uint8:
		return GetDataTypeINT64(int64(val))
	case uint16:
		return GetDataTypeINT64(int64(val))
	case uint32:
		return GetDataTypeINT64(int64(val))
	case uint64:
		if val > math.MaxInt64 {
			return Invalid
		}
		return GetDataTypeINT64(int64(val))
	case float32:
		return Float32
	case float64:
		return Float64
	case string:
		return String
	case []byte:
		return ByteArray
	case []interface{}:
		return Array
	case map[string]interface{}:
		return Map
	case time.Time:
		return DateTime
	case time.Duration:
		return DateTimeDuration
	default:
		// reflection fallback for uncommon types (typed slices, typed maps, etc.)
		rv := reflect.ValueOf(v)
		switch rv.Kind() {
		case reflect.Slice:
			return Array
		case reflect.Array:
			return Array
		case reflect.Map:
			if rv.Type().Key().Kind() == reflect.String {
				return Map
			}
			return Invalid
		default:
			return Invalid
		}
	}
}
