package encoding

import (
	"fmt"
	"reflect"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

// BytesEncoder passes through raw bytes
type BytesEncoder struct{}

// Encode returns bytes as-is or converts strings
func (e *BytesEncoder) Encode(v any) ([]byte, error) {
	switch val := v.(type) {
	case []byte:
		return val, nil
	case string:
		return []byte(val), nil
	case *[]byte:
		if val == nil {
			return nil, core.SerializationError{
				Operation: "encode",
				Type:      "*[]byte",
				Err:       core.ErrInvalidMessage,
			}
		}
		return *val, nil
	case *string:
		if val == nil {
			return nil, core.SerializationError{
				Operation: "encode",
				Type:      "*string",
				Err:       core.ErrInvalidMessage,
			}
		}
		return []byte(*val), nil
	default:
		return nil, core.SerializationError{
			Operation: "encode",
			Type:      reflect.TypeOf(v).String(),
			Err:       fmt.Errorf("unsupported type for bytes encoder: %T", v),
		}
	}
}

// Decode copies bytes to the target
func (e *BytesEncoder) Decode(data []byte, v any) error {
	switch val := v.(type) {
	case *[]byte:
		if val == nil {
			return core.SerializationError{
				Operation: "decode",
				Type:      "*[]byte",
				Err:       core.ErrInvalidMessage,
			}
		}
		*val = append([]byte{}, data...)
		return nil
	case *string:
		if val == nil {
			return core.SerializationError{
				Operation: "decode",
				Type:      "*string",
				Err:       core.ErrInvalidMessage,
			}
		}
		*val = string(data)
		return nil
	default:
		return core.SerializationError{
			Operation: "decode",
			Type:      reflect.TypeOf(v).String(),
			Err:       fmt.Errorf("unsupported type for bytes decoder: %T", v),
		}
	}
}

// ContentType returns the octet-stream MIME type
func (e *BytesEncoder) ContentType() string {
	return "application/octet-stream"
}

// Bytes is the default bytes encoder
var Bytes = &BytesEncoder{}

// EncodeBytes encodes to bytes
func EncodeBytes(v any) ([]byte, error) {
	return Bytes.Encode(v)
}

// DecodeBytes decodes bytes
func DecodeBytes(data []byte, v any) error {
	return Bytes.Decode(data, v)
}
