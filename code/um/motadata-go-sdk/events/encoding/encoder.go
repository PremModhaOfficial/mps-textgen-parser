package encoding

// Encoder defines the interface for message encoding
type Encoder interface {
	// Encode converts a value to bytes
	Encode(v any) ([]byte, error)

	// Decode converts bytes to a value
	Decode(data []byte, v any) error

	// ContentType returns the MIME type
	ContentType() string
}

// EncoderFunc is a function type that implements Encoder for encoding
type EncoderFunc func(v any) ([]byte, error)

// DecoderFunc is a function type for decoding
type DecoderFunc func(data []byte, v any) error

// Registry holds registered encoders by content type
type Registry struct {
	encoders map[string]Encoder
	fallback Encoder
}

// NewRegistry creates a new encoder registry
func NewRegistry() *Registry {
	r := &Registry{
		encoders: make(map[string]Encoder),
		fallback: JSON, // Default fallback
	}

	// Register default encoders
	r.Register(JSON)
	r.Register(Bytes)

	return r
}

// Register adds an encoder to the registry
func (r *Registry) Register(enc Encoder) {
	r.encoders[enc.ContentType()] = enc
}

// Get retrieves an encoder by content type
func (r *Registry) Get(contentType string) Encoder {
	if enc, ok := r.encoders[contentType]; ok {
		return enc
	}
	return r.fallback
}

// SetFallback sets the fallback encoder
func (r *Registry) SetFallback(enc Encoder) {
	r.fallback = enc
}

// Encode encodes using the appropriate encoder for the content type
func (r *Registry) Encode(contentType string, v any) ([]byte, error) {
	return r.Get(contentType).Encode(v)
}

// Decode decodes using the appropriate encoder for the content type
func (r *Registry) Decode(contentType string, data []byte, v any) error {
	return r.Get(contentType).Decode(data, v)
}

// DefaultRegistry is the global encoder registry
var DefaultRegistry = NewRegistry()

// Encode encodes using the default registry
func Encode(contentType string, v any) ([]byte, error) {
	return DefaultRegistry.Encode(contentType, v)
}

// Decode decodes using the default registry
func Decode(contentType string, data []byte, v any) error {
	return DefaultRegistry.Decode(contentType, data, v)
}
