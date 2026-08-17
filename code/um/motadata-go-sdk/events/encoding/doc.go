// Package encoding provides message serialization and deserialization for the events system.
//
// This package supports multiple encoding formats for message payloads:
//   - JSON: Human-readable format with broad compatibility
//   - Protocol Buffers: Efficient binary format for high-performance scenarios
//   - Raw bytes: Pass-through encoding for pre-serialized data
//
// # Basic Usage
//
// Encode and decode data using the global functions:
//
//	// Encode to JSON
//	data, err := encoding.EncodeJSON(user)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Decode from JSON
//	var user User
//	err = encoding.DecodeJSON(data, &user)
//
// # Available Encoders
//
// The package provides three built-in encoders:
//
//	encoding.JSON      - JSON encoder (application/json)
//	encoding.Bytes     - Raw bytes encoder (application/octet-stream)
//	encoding.Protobuf  - Protocol Buffers encoder (application/protobuf)
//
// # JSON Encoding
//
// For standard JSON encoding:
//
//	type Event struct {
//	    Type    string    `json:"type"`
//	    Payload any       `json:"payload"`
//	    Time    time.Time `json:"time"`
//	}
//
//	event := Event{Type: "user.created", Payload: user, Time: time.Now()}
//
//	// Encode
//	data, err := encoding.EncodeJSON(event)
//
//	// Decode
//	var decoded Event
//	err = encoding.DecodeJSON(data, &decoded)
//
// # Protocol Buffers
//
// For Protocol Buffer encoding (requires proto.Message implementation):
//
//	// Your protobuf message
//	protoMsg := &pb.UserEvent{
//	    UserId: "123",
//	    Action: pb.Action_CREATED,
//	}
//
//	// Encode
//	data, err := encoding.EncodeProtobuf(protoMsg)
//
//	// Decode
//	var decoded pb.UserEvent
//	err = encoding.DecodeProtobuf(data, &decoded)
//
// # Raw Bytes
//
// For pre-serialized or binary data:
//
//	// Encode (pass-through)
//	data, err := encoding.EncodeBytes(rawData)
//
//	// Decode (copies to destination)
//	var result []byte
//	err = encoding.DecodeBytes(data, &result)
//
// # Encoder Interface
//
// Implement the Encoder interface for custom formats:
//
//	type Encoder interface {
//	    Encode(v any) ([]byte, error)
//	    Decode(data []byte, v any) error
//	    ContentType() string
//	}
//
// Example custom encoder:
//
//	type MsgPackEncoder struct{}
//
//	func (e MsgPackEncoder) Encode(v any) ([]byte, error) {
//	    return msgpack.Marshal(v)
//	}
//
//	func (e MsgPackEncoder) Decode(data []byte, v any) error {
//	    return msgpack.Unmarshal(data, v)
//	}
//
//	func (e MsgPackEncoder) ContentType() string {
//	    return "application/msgpack"
//	}
//
// # Encoder Registry
//
// Use the registry to manage multiple encoders:
//
//	// Get the default registry
//	registry := encoding.DefaultRegistry
//
//	// Register a custom encoder
//	registry.Register(myCustomEncoder)
//
//	// Encode using content type
//	data, err := registry.Encode("application/msgpack", payload)
//
//	// Decode using content type
//	err = registry.Decode("application/msgpack", data, &result)
//
// # Message Integration
//
// Set content type in message headers:
//
//	msg := core.NewMessage(data)
//	msg.WithHeader(core.HeaderContentType, encoding.JSON.ContentType())
//
// Decode based on message content type:
//
//	contentType := msg.Headers.Get(core.HeaderContentType)
//	err := encoding.Decode(contentType, msg.Data, &payload)
//
// # Error Handling
//
// Encoding errors are wrapped in core.SerializationError:
//
//	data, err := encoding.EncodeJSON(invalidData)
//	var serErr core.SerializationError
//	if errors.As(err, &serErr) {
//	    log.Printf("Serialization error: %s %s: %v",
//	        serErr.Operation, serErr.Type, serErr.Err)
//	}
package encoding
