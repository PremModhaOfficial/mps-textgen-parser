package encoding

import (
	"fmt"
	"reflect"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

// ProtobufEncoder implements Protocol Buffers encoding
// Requires messages to implement Marshal/Unmarshal methods
type ProtobufEncoder struct{}

// ProtoMarshaler interface for protobuf messages
type ProtoMarshaler interface {
	Marshal() ([]byte, error)
}

// ProtoUnmarshaler interface for protobuf messages
type ProtoUnmarshaler interface {
	Unmarshal([]byte) error
}

// Encode converts a protobuf message to bytes
func (e *ProtobufEncoder) Encode(v any) ([]byte, error) {
	if v == nil {
		return nil, core.SerializationError{
			Operation: "encode",
			Type:      "nil",
			Err:       core.ErrInvalidMessage,
		}
	}

	// Check if value implements Marshal method
	if marshaler, ok := v.(ProtoMarshaler); ok {
		data, err := marshaler.Marshal()
		if err != nil {
			return nil, core.SerializationError{
				Operation: "encode",
				Type:      reflect.TypeOf(v).String(),
				Err:       err,
			}
		}
		return data, nil
	}

	// Check for MarshalBinary (alternative interface)
	if marshaler, ok := v.(interface{ MarshalBinary() ([]byte, error) }); ok {
		data, err := marshaler.MarshalBinary()
		if err != nil {
			return nil, core.SerializationError{
				Operation: "encode",
				Type:      reflect.TypeOf(v).String(),
				Err:       err,
			}
		}
		return data, nil
	}

	return nil, core.SerializationError{
		Operation: "encode",
		Type:      reflect.TypeOf(v).String(),
		Err:       fmt.Errorf("type does not implement Marshal method"),
	}
}

// Decode converts bytes to a protobuf message
func (e *ProtobufEncoder) Decode(data []byte, v any) error {
	if len(data) == 0 {
		return core.SerializationError{
			Operation: "decode",
			Type:      "empty",
			Err:       core.ErrInvalidMessage,
		}
	}

	// Check if value implements Unmarshal method
	if unmarshaler, ok := v.(ProtoUnmarshaler); ok {
		if err := unmarshaler.Unmarshal(data); err != nil {
			return core.SerializationError{
				Operation: "decode",
				Type:      reflect.TypeOf(v).String(),
				Err:       err,
			}
		}
		return nil
	}

	// Check for UnmarshalBinary (alternative interface)
	if unmarshaler, ok := v.(interface{ UnmarshalBinary([]byte) error }); ok {
		if err := unmarshaler.UnmarshalBinary(data); err != nil {
			return core.SerializationError{
				Operation: "decode",
				Type:      reflect.TypeOf(v).String(),
				Err:       err,
			}
		}
		return nil
	}

	return core.SerializationError{
		Operation: "decode",
		Type:      reflect.TypeOf(v).String(),
		Err:       fmt.Errorf("type does not implement Unmarshal method"),
	}
}

// ContentType returns the protobuf MIME type
func (e *ProtobufEncoder) ContentType() string {
	return "application/x-protobuf"
}

// Protobuf is the default protobuf encoder
var Protobuf = &ProtobufEncoder{}

// EncodeProtobuf encodes a protobuf message
func EncodeProtobuf(v any) ([]byte, error) {
	return Protobuf.Encode(v)
}

// DecodeProtobuf decodes a protobuf message
func DecodeProtobuf(data []byte, v any) error {
	return Protobuf.Decode(data, v)
}
