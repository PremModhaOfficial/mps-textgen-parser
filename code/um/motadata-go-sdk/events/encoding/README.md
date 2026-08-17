# Encoding Package

The `encoding` package provides message serialization and deserialization for the events messaging system with support for JSON, Protobuf, and raw bytes.

## Overview

This package contains:

- **Encoder Interface**: Common interface for encoders
- **JSON Encoder**: JSON serialization with standard library
- **Protobuf Encoder**: Protocol Buffers serialization
- **Bytes Encoder**: Raw bytes passthrough

## Package Structure

```
encoding/
├── encoder.go   # Encoder interface
├── json.go      # JSON encoder implementation
├── protobuf.go  # Protobuf encoder implementation
├── bytes.go     # Bytes encoder implementation
└── README.md    # This file
```

## Encoder Interface

```go
type Encoder interface {
    // Encode serializes a value to bytes
    Encode(v any) ([]byte, error)

    // Decode deserializes bytes to a value
    Decode(data []byte, v any) error

    // ContentType returns the MIME content type
    ContentType() string
}
```

## Available Encoders

### JSON Encoder

Standard JSON encoding using `encoding/json`:

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/encoding"

// Use singleton
encoder := encoding.JSON

// Or use functions directly
data, err := encoding.EncodeJSON(myStruct)
err = encoding.DecodeJSON(data, &myStruct)
```

**Content-Type:** `application/json`

**Example:**

```go
type Order struct {
    ID     string  `json:"id"`
    Amount float64 `json:"amount"`
}

order := Order{ID: "123", Amount: 99.99}

// Encode
data, err := encoding.EncodeJSON(order)
// data: {"id":"123","amount":99.99}

// Decode
var decoded Order
err = encoding.DecodeJSON(data, &decoded)
```

### Protobuf Encoder

Protocol Buffers encoding for high-performance scenarios:

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/encoding"

// Use singleton
encoder := encoding.Protobuf

// Or use functions directly
data, err := encoding.EncodeProtobuf(protoMessage)
err = encoding.DecodeProtobuf(data, protoMessage)
```

**Content-Type:** `application/protobuf`

**Requirements:**
- Values must implement `proto.Message`
- Generate Go code from `.proto` files

**Example:**

```go
// Assuming generated proto message
message := &pb.OrderCreated{
    OrderId: "123",
    Amount:  99.99,
}

// Encode
data, err := encoding.EncodeProtobuf(message)

// Decode
var decoded pb.OrderCreated
err = encoding.DecodeProtobuf(data, &decoded)
```

### Bytes Encoder

Raw bytes passthrough for binary data:

```go
import "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/encoding"

// Use singleton
encoder := encoding.Bytes

// Or use functions directly
data, err := encoding.EncodeBytes(rawData)
err = encoding.DecodeBytes(data, &rawData)
```

**Content-Type:** `application/octet-stream`

**Supported Types:**
- `[]byte` - Direct passthrough
- `string` - Converted to bytes
- `io.Reader` - Read all bytes

**Example:**

```go
// Encode bytes
rawData := []byte{0x01, 0x02, 0x03, 0x04}
data, err := encoding.EncodeBytes(rawData)

// Encode string
strData := "Hello, World!"
data, err := encoding.EncodeBytes(strData)

// Decode
var decoded []byte
err = encoding.DecodeBytes(data, &decoded)
```

## Usage with Messages

### Setting Content-Type

```go
import (
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
    "dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/encoding"
)

// Create message with JSON
order := Order{ID: "123", Amount: 99.99}
data, _ := encoding.EncodeJSON(order)

msg := core.NewMessage(data).
    WithHeader(core.HeaderContentType, encoding.JSON.ContentType())
```

### Auto-Detecting Encoding

```go
func decodeMessage(msg *core.Message, v any) error {
    contentType := msg.Headers.Get(core.HeaderContentType)

    switch contentType {
    case "application/json":
        return encoding.DecodeJSON(msg.Data, v)
    case "application/protobuf":
        return encoding.DecodeProtobuf(msg.Data, v)
    case "application/octet-stream":
        return encoding.DecodeBytes(msg.Data, v)
    default:
        // Default to JSON
        return encoding.DecodeJSON(msg.Data, v)
    }
}
```

## Encoder Comparison

| Encoder | Speed | Size | Human Readable | Schema |
|---------|-------|------|----------------|--------|
| JSON | Medium | Large | Yes | Optional |
| Protobuf | Fast | Small | No | Required |
| Bytes | Fastest | Exact | No | None |

### When to Use Each

**JSON:**
- API responses
- Configuration data
- Debugging/logging
- Cross-language compatibility
- Human-readable requirements

**Protobuf:**
- High-throughput systems
- Large message volumes
- Bandwidth-constrained environments
- Strong schema requirements
- Backward compatibility needs

**Bytes:**
- Binary file transfer
- Pre-serialized data
- Custom serialization
- Raw sensor data

## Helper Functions

### Quick Encoding

```go
// JSON
jsonData, err := encoding.EncodeJSON(myStruct)
err = encoding.DecodeJSON(jsonData, &myStruct)

// Protobuf
protoData, err := encoding.EncodeProtobuf(protoMsg)
err = encoding.DecodeProtobuf(protoData, protoMsg)

// Bytes
bytesData, err := encoding.EncodeBytes(rawBytes)
err = encoding.DecodeBytes(bytesData, &rawBytes)
```

### Using Encoder Interface

```go
func publishWithEncoder(pub core.Publisher, subject string, v any, enc encoding.Encoder) error {
    data, err := enc.Encode(v)
    if err != nil {
        return fmt.Errorf("encode: %w", err)
    }

    msg := core.NewMessage(data).
        WithHeader(core.HeaderContentType, enc.ContentType())

    return pub.Publish(context.Background(), subject, msg)
}

// Usage
publishWithEncoder(pub, "orders", order, encoding.JSON)
publishWithEncoder(pub, "events", event, encoding.Protobuf)
```

## Custom Encoder

Implement the `Encoder` interface for custom serialization:

```go
type MsgPackEncoder struct{}

func (e *MsgPackEncoder) Encode(v any) ([]byte, error) {
    return msgpack.Marshal(v)
}

func (e *MsgPackEncoder) Decode(data []byte, v any) error {
    return msgpack.Unmarshal(data, v)
}

func (e *MsgPackEncoder) ContentType() string {
    return "application/msgpack"
}

// Usage
encoder := &MsgPackEncoder{}
data, _ := encoder.Encode(myStruct)
```

## Error Handling

```go
data, err := encoding.EncodeJSON(value)
if err != nil {
    // Handle encoding error
    // Common causes: unsupported types, circular references
    return fmt.Errorf("failed to encode: %w", err)
}

err = encoding.DecodeJSON(data, &value)
if err != nil {
    // Handle decoding error
    // Common causes: invalid JSON, type mismatch
    return fmt.Errorf("failed to decode: %w", err)
}
```

### Common Errors

| Error | Cause | Solution |
|-------|-------|----------|
| `unsupported type` | Type can't be serialized | Use supported types |
| `json: cannot unmarshal` | Type mismatch | Check target type |
| `proto: wrong type` | Not proto.Message | Use generated types |
| `unexpected EOF` | Incomplete data | Check message integrity |

## Performance Tips

### 1. Reuse Encoders

```go
// Good - use singletons
data, _ := encoding.JSON.Encode(v)

// Avoid - creating new encoders
encoder := &encoding.JSONEncoder{}
data, _ := encoder.Encode(v)
```

### 2. Pre-allocate for Decoding

```go
// Reuse struct instead of creating new
var order Order
for msg := range messages {
    encoding.DecodeJSON(msg.Data, &order)
    process(order)
    order = Order{} // Reset for next iteration
}
```

### 3. Use Protobuf for High Volume

```go
// For high-throughput scenarios, Protobuf is 3-5x faster
data, _ := encoding.EncodeProtobuf(protoMsg)  // Faster
data, _ := encoding.EncodeJSON(struct)        // Slower
```

### 4. Skip Encoding for Bytes

```go
// If data is already bytes, use Bytes encoder
rawData := []byte{...}
msg := core.NewMessage(rawData)  // No encoding overhead
```

## Testing

```bash
# Run tests
go test ./events/encoding/...

# Run with coverage
go test -cover ./events/encoding/...

# Run benchmarks
go test -bench=. ./events/encoding/...
```

## Usage in Other Packages

This package is used by:
- `events/transport` - Message encoding in transport
- `events/middleware` - Content-type handling
- `events` (main) - Re-exports for convenience
- Application code - Message serialization
