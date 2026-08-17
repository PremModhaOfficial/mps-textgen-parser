# Codec Package - Binary Serialization

## Overview

The codec package provides compact binary serialization and deserialization for Go data structures. It offers efficient encoding and decoding of maps and arrays with support for all primitive Go types and nested structures.

### Key Features

- **Compact Binary Format**: Minimal overhead for efficient storage and transmission
- **Type Safety**: Strong type preservation during serialization/deserialization
- **Nested Structure Support**: Handles deeply nested maps and arrays
- **Error Recovery**: Graceful handling of corrupted or invalid data
- **Little-Endian Format**: Consistent byte ordering across platforms
- **Buffer Reuse**: Supports custom buffers for reduced allocations
- **Codec Selection**: Custom binary or MsgPack serialization formats
- **Compression Support**: None, Snappy, Zstd, Gzip, LZ4

## Installation

```go
import "motadatagosdk/core/codec"
```

## Usage Guide

### Quick Start

The codec package provides two main functions for serialization:
- `PackMap()` - Serialize map data structures
- `PackArray()` - Serialize arrays/slices

And two for deserialization:
- `UnpackMap()` - Deserialize to map
- `UnpackArray()` - Deserialize to array

### API

```go
func PackMap(data map[string]interface{}, encoder Encoder, compression Compression, buffer *bytes.Buffer) ([]byte, error)
func PackArray(data []interface{}, encoder Encoder, compression Compression, buffer *bytes.Buffer) ([]byte, error)
func UnpackMap(data []byte) (map[string]interface{}, error)
func UnpackArray(data []byte) ([]interface{}, error)
```

### Wire Format

- `Pack*` produces: `[1-byte header][compressed payload]`
- Header: high 4 bits = compression, low 4 bits = codec type

## Basic Usage

### Packing and Unpacking Maps

```go
import "motadatagosdk/core/codec"

// Create data to pack
data := map[string]interface{}{
    "name":   "John Doe",
    "age":    30,
    "active": true,
    "scores": []interface{}{95, 87, 92},
}

// Pack the map
packed, err := codec.PackMap(data, codec.CodecCustom, codec.CompressNone, nil)
if err != nil {
    // Handle error
}

// Unpack the data
unpacked, err := codec.UnpackMap(packed)
if err != nil {
    // Handle error
}

// Use unpacked data
name := unpacked["name"].(string)  // "John Doe"
age := unpacked["age"].(int8)      // 30 (smallest type that fits)
```

### Packing and Unpacking Arrays

```go
// Create array data
data := []interface{}{
    "item1",
    42,
    true,
    3.14,
    map[string]interface{}{"key": "value"},
}

// Pack the array
packed, err := codec.PackArray(data, codec.CodecCustom, codec.CompressNone, nil)
if err != nil {
    // Handle error
}

// Unpack the array
unpacked, err := codec.UnpackArray(packed)
if err != nil {
    // Handle error
}

// Use unpacked data
first := unpacked[0].(string)  // "item1"
number := unpacked[1].(int8)   // 42
```

### Using Custom Buffers

```go
import "bytes"

// Reuse buffers for better performance
buffer := &bytes.Buffer{}

// Pack with custom buffer (buffer is reset automatically)
data := map[string]interface{}{"key": "value"}
packed, err := codec.PackMap(data, codec.CodecCustom, codec.CompressNone, buffer)
```

### Using MsgPack Codec

```go
// Use MsgPack for cross-service communication
packed, err := codec.PackMap(data, codec.CodecMsgPack, codec.CompressNone, nil)
if err != nil {
    // Handle error
}

unpacked, err := codec.UnpackMap(packed)
```

### Using Compression

```go
// Pack with Snappy compression
packed, err := codec.PackMap(data, codec.CodecCustom, codec.CompressSnappy, nil)

// Pack with Zstd compression
packed, err = codec.PackMap(data, codec.CodecCustom, codec.CompressZstd, nil)

// Pack with Gzip compression
packed, err = codec.PackMap(data, codec.CodecCustom, codec.CompressGzip, nil)

// Pack with LZ4 compression
packed, err = codec.PackMap(data, codec.CodecCustom, codec.CompressLZ4, nil)

// Unpack auto-detects codec and compression from the header byte
unpacked, err := codec.UnpackMap(packed)
```

## Codec/Compression Constants

```go
// Codec types
codec.CodecCustom   // 0x00 - Inter-service communication (custom binary)
codec.CodecMsgPack  // 0x01 - Cross-service communication (standard format)

// Compression algorithms
codec.CompressNone   // 0x00
codec.CompressSnappy // 0x01
codec.CompressZstd   // 0x02
codec.CompressGzip   // 0x03
codec.CompressLZ4    // 0x04
```

## Supported Data Types

The codec supports the following Go types:

| Go Type | Packed As | Unpacked As | Notes |
|---------|-----------|-------------|-------|
| `nil` | Invalid | `nil` | Preserved as nil |
| `bool` | Boolean | `bool` | True/false |
| `int8` | Int8 | `int8` | 8-bit signed |
| `int16` | Int16 | `int16` | 16-bit signed |
| `int32` | Int32 | `int32` | 32-bit signed |
| `int64` | Int64 | `int64` | 64-bit signed |
| `int` | Int8-Int64 | `int8`-`int64` | Smallest type that fits |
| `uint8` | Int8-Int16 | `int8`-`int16` | Converted to signed |
| `uint16` | Int16-Int32 | `int16`-`int32` | Converted to signed |
| `uint32` | Int32-Int64 | `int32`-`int64` | Converted to signed |
| `uint64` | Int64 | `int64` | Max `math.MaxInt64` |
| `uint` | Int8-Int64 | `int8`-`int64` | Converted to signed |
| `float32` | Float32 | `float32` | 32-bit float |
| `float64` | Float64 | `float64` | 64-bit float |
| `string` | String | `string` | UTF-8 encoded |
| `[]byte` | ByteArray | `[]byte` | Raw byte data |
| `[]interface{}` | Array | `[]interface{}` | Generic array |
| Any slice/array | Array | `[]interface{}` | Via reflection |
| `map[string]interface{}` | Map | `map[string]interface{}` | String-keyed map |
| Any `map[string]T` | Map | `map[string]interface{}` | Via reflection |

### Notes

- Integer values are unpacked using the smallest signed type that fits (`int8`..`int64`).
- `uint64` values above `math.MaxInt64` return an error in custom codec mode.
- Custom codec length fields for maps/arrays/map keys are 16-bit; oversized values return an error.

## Binary Format Specification

### Map Format
```
[2 bytes] - Number of key-value pairs (uint16, little-endian)
[1 byte]  - Data type marker (Map)
For each key-value pair:
    [2 bytes]    - Key length (uint16, little-endian)
    [N bytes]    - Key string (UTF-8)
    [1 byte]     - Value type
    [Variable]   - Value data (format depends on type)
```

### Array Format
```
[2 bytes] - Number of elements (uint16, little-endian)
[1 byte]  - Data type marker (Array)
For each element:
    [1 byte]     - Element type
    [Variable]   - Element data (format depends on type)
```

### Numeric Encoding (Little-Endian)
- **Int8/UInt8**: 1 byte
- **Int16/UInt16**: 2 bytes
- **Int24**: 3 bytes
- **Int32/UInt32**: 4 bytes
- **Int40**: 5 bytes
- **Int48**: 6 bytes
- **Int56**: 7 bytes
- **Int64/UInt64**: 8 bytes
- **Float32**: 4 bytes (IEEE 754)
- **Float64**: 8 bytes (IEEE 754)

### String Encoding
```
[4 bytes] - String length (uint32, little-endian)
[N bytes] - String data (UTF-8)
```

### ByteArray Encoding
```
[4 bytes] - Byte array length (uint32, little-endian)
[N bytes] - Raw bytes
```

### Boolean Encoding
- `true`: 0x01
- `false`: 0x00

## Advanced Usage

### Nested Structures

```go
// Complex nested structure
data := map[string]interface{}{
    "user": map[string]interface{}{
        "profile": map[string]interface{}{
            "name": "Alice",
            "tags": []interface{}{"admin", "developer"},
        },
        "settings": map[string]interface{}{
            "theme": "dark",
            "notifications": true,
        },
    },
    "timestamp": int64(1234567890),
}

// Pack and unpack preserves structure
packed, err := codec.PackMap(data, codec.CodecCustom, codec.CompressNone, nil)
unpacked, err := codec.UnpackMap(packed)

// Navigate nested structure
user := unpacked["user"].(map[string]interface{})
profile := user["profile"].(map[string]interface{})
name := profile["name"].(string)  // "Alice"
```

### Typed Maps and Slices

```go
// Typed maps are supported via reflection
data := map[string]interface{}{
    "typed": map[string]int{"a": 1, "b": 2},
    "slice": []int{10, 20, 30},
}

packed, err := codec.PackMap(data, codec.CodecCustom, codec.CompressNone, nil)
```

### Error Handling

```go
// Handle corrupted data
corruptedData := []byte{255, 255, 255, 255}
result, err := codec.UnpackMap(corruptedData)
if err != nil {
    if errors.Is(err, utils.ErrUnpackFailed) {
        // Data is corrupted or invalid
        fmt.Println("Failed to unpack: corrupted data")
    }
}
```

## Error Handling

Codec errors are defined in `motadatagosdk/utils`:

```go
var (
    ErrUnpackFailed           = errors.New("unpack failed: invalid or corrupted data")
    ErrUnsupportedCodec       = errors.New("unsupported codec type")
    ErrUnsupportedCompression = errors.New("unsupported compression algorithm")
    ErrUnsupportedDataType    = errors.New("unsupported data type")
    ErrValueOutOfRange        = errors.New("value out of supported range")
    ErrDataTooLarge           = errors.New("data too large for codec length fields")
    ErrCompressionFailed      = errors.New("compression failed")
    ErrDecompressionFailed    = errors.New("decompression failed")
)
```

## Implementation Details

### Space Characteristics
- Minimal overhead: 3 bytes for empty map/array
- Compact numeric encoding with little-endian format
- Variable-width integer encoding (Int8 through Int64)
- No field names stored (uses map keys directly)

### Time Complexity
| Operation | Complexity |
|-----------|------------|
| PackMap | O(n) where n is total elements |
| UnpackMap | O(n) where n is total elements |
| PackArray | O(n) where n is array length |
| UnpackArray | O(n) where n is array length |

### Memory Usage
- Buffer reuse supported for reduced allocations
- In-place encoding with provided buffers
- Minimal temporary allocations during unpacking

### Size Limits
```go
// Be aware of size limits
// - Max map/array size: 65535 elements (uint16)
// - Max string length: 4GB (uint32)
// - Max key length: 65535 bytes (uint16)
// - Max uint64 value: math.MaxInt64 (custom codec)
```

## Best Practices

### 1. Buffer Reuse
```go
// Reuse buffers in hot paths
buffer := &bytes.Buffer{}
for _, data := range dataSlice {
    packed, err := codec.PackMap(data, codec.CodecCustom, codec.CompressNone, buffer)
    // Process packed data
}
```

### 2. Type Assertions
```go
// Always check types after unpacking
unpacked, err := codec.UnpackMap(data)
if err != nil {
    return err
}

// Safe type assertion
if name, ok := unpacked["name"].(string); ok {
    // Use name
}
```

### 3. Integer Types
```go
// Integers are unpacked as the smallest type that fits
packed, _ := codec.PackMap(map[string]interface{}{
    "count": 42,  // int
}, codec.CodecCustom, codec.CompressNone, nil)

unpacked, _ := codec.UnpackMap(packed)
count := unpacked["count"].(int8)  // Fits in int8
```

## Use Cases

### Ideal For:
- **Internal service communication** within Go applications
- **Cache serialization** for memory or disk storage
- **Message queue payloads** where size matters
- **Binary protocols** with known Go endpoints
- **Configuration storage** in binary format

### Considerations:
- Best suited for Go-to-Go communication
- Not human-readable (binary format)
- Schema evolution requires careful planning
- Use MsgPack codec for cross-language interoperability

## Testing

Run tests:
```bash
# Run all tests
go test ./...

# Run with race detection
go test -race ./...

# Run benchmarks
go test -bench=. ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestPackMap ./...
```

## Thread Safety

**The codec functions are stateless and thread-safe for concurrent use.** Multiple goroutines can safely call pack/unpack functions simultaneously with different data.

However, if using shared buffers, ensure proper synchronization:

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return &bytes.Buffer{}
    },
}

func packWithPool(data map[string]interface{}) ([]byte, error) {
    buffer := bufferPool.Get().(*bytes.Buffer)
    defer bufferPool.Put(buffer)
    return codec.PackMap(data, codec.CodecCustom, codec.CompressNone, buffer)
}
```

## Complete Round-Trip Example

```go
package main

import (
    "fmt"
    "motadatagosdk/core/codec"
)

func main() {
    // Original data
    original := map[string]interface{}{
        "id": int64(123),
        "user": map[string]interface{}{
            "name": "Alice",
            "email": "alice@example.com",
            "active": true,
        },
        "scores": []interface{}{95.5, 87.0, 92.5},
        "metadata": nil,
    }

    // Pack the data with Snappy compression
    packed, err := codec.PackMap(original, codec.CodecCustom, codec.CompressSnappy, nil)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Packed size: %d bytes\n", len(packed))

    // Unpack the data (auto-detects codec and compression)
    unpacked, err := codec.UnpackMap(packed)
    if err != nil {
        panic(err)
    }

    // Access unpacked data
    user := unpacked["user"].(map[string]interface{})
    fmt.Printf("User name: %s\n", user["name"])

    scores := unpacked["scores"].([]interface{})
    fmt.Printf("First score: %v\n", scores[0])
}
```

## Summary

The codec package provides efficient, type-safe binary serialization for Go applications. It offers compact binary encoding for maps and arrays with multiple codec formats (custom binary, MsgPack) and compression algorithms (Snappy, Zstd, Gzip, LZ4), making it suitable for internal data serialization, caching, and message passing within Go applications.
