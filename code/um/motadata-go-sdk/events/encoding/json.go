package encoding

import (
	"encoding/json"
	"reflect"

	"dev.azure.com/Motadata/NextGen/motadata-go-sdk/events/core"
)

// JSONEncoder implements JSON encoding
type JSONEncoder struct {
	// Indent enables pretty-printed JSON
	Indent bool
}

// Encode converts a value to JSON bytes
func (e *JSONEncoder) Encode(v any) ([]byte, error) {
	if v == nil {
		return nil, core.SerializationError{
			Operation: "encode",
			Type:      "nil",
			Err:       core.ErrInvalidMessage,
		}
	}

	var data []byte
	var err error

	if e.Indent {
		data, err = json.MarshalIndent(v, "", "  ")
	} else {
		data, err = json.Marshal(v)
	}

	if err != nil {
		return nil, core.SerializationError{
			Operation: "encode",
			Type:      reflect.TypeOf(v).String(),
			Err:       err,
		}
	}

	return data, nil
}

// Decode converts JSON bytes to a value
func (e *JSONEncoder) Decode(data []byte, v any) error {
	if len(data) == 0 {
		return core.SerializationError{
			Operation: "decode",
			Type:      "empty",
			Err:       core.ErrInvalidMessage,
		}
	}

	if err := json.Unmarshal(data, v); err != nil {
		return core.SerializationError{
			Operation: "decode",
			Type:      reflect.TypeOf(v).String(),
			Err:       err,
		}
	}

	return nil
}

// ContentType returns the JSON MIME type
func (e *JSONEncoder) ContentType() string {
	return "application/json"
}

// JSON is the default JSON encoder
var JSON = &JSONEncoder{}

// JSONPretty is a JSON encoder with indentation
var JSONPretty = &JSONEncoder{Indent: true}

// EncodeJSON encodes a value to JSON
func EncodeJSON(v any) ([]byte, error) {
	return JSON.Encode(v)
}

// DecodeJSON decodes JSON to a value
func DecodeJSON(data []byte, v any) error {
	return JSON.Decode(data, v)
}
