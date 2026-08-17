package encoding

import (
	"errors"
	"testing"
)

// ============== JSON Encoder Tests ==============

type testStruct struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestJSONEncode(t *testing.T) {
	data := testStruct{Name: "test", Value: 42}
	encoded, err := JSON.Encode(data)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `{"name":"test","value":42}`
	if string(encoded) != expected {
		t.Errorf("expected %s, got %s", expected, string(encoded))
	}
}

func TestJSONEncodeNil(t *testing.T) {
	_, err := JSON.Encode(nil)
	if err == nil {
		t.Error("expected error for nil value")
	}
}

func TestJSONEncodeInvalid(t *testing.T) {
	// Channels can't be encoded to JSON
	ch := make(chan int)
	_, err := JSON.Encode(ch)
	if err == nil {
		t.Error("expected error for invalid type")
	}
}

func TestJSONEncodeIndent(t *testing.T) {
	data := testStruct{Name: "test", Value: 42}
	encoded, err := JSONPretty.Encode(data)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should contain newlines and spaces
	if len(encoded) <= 20 {
		t.Error("expected indented output to be longer")
	}
}

func TestJSONDecode(t *testing.T) {
	data := []byte(`{"name":"test","value":42}`)
	var result testStruct

	err := JSON.Decode(data, &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "test" || result.Value != 42 {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestJSONDecodeEmpty(t *testing.T) {
	var result testStruct
	err := JSON.Decode([]byte{}, &result)
	if err == nil {
		t.Error("expected error for empty data")
	}
}

func TestJSONDecodeInvalid(t *testing.T) {
	var result testStruct
	err := JSON.Decode([]byte("not json"), &result)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestJSONContentType(t *testing.T) {
	if JSON.ContentType() != "application/json" {
		t.Errorf("unexpected content type: %s", JSON.ContentType())
	}
}

func TestEncodeJSON(t *testing.T) {
	data := map[string]string{"key": "value"}
	encoded, err := EncodeJSON(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(encoded) != `{"key":"value"}` {
		t.Errorf("unexpected result: %s", string(encoded))
	}
}

func TestDecodeJSON(t *testing.T) {
	var result map[string]string
	err := DecodeJSON([]byte(`{"key":"value"}`), &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["key"] != "value" {
		t.Errorf("unexpected result: %v", result)
	}
}

// ============== Bytes Encoder Tests ==============

func TestBytesEncodeBytes(t *testing.T) {
	input := []byte("hello")
	encoded, err := Bytes.Encode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(encoded) != "hello" {
		t.Errorf("unexpected result: %s", string(encoded))
	}
}

func TestBytesEncodeString(t *testing.T) {
	input := "hello"
	encoded, err := Bytes.Encode(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(encoded) != "hello" {
		t.Errorf("unexpected result: %s", string(encoded))
	}
}

func TestBytesEncodePointerBytes(t *testing.T) {
	input := []byte("hello")
	encoded, err := Bytes.Encode(&input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(encoded) != "hello" {
		t.Errorf("unexpected result: %s", string(encoded))
	}
}

func TestBytesEncodePointerString(t *testing.T) {
	input := "hello"
	encoded, err := Bytes.Encode(&input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(encoded) != "hello" {
		t.Errorf("unexpected result: %s", string(encoded))
	}
}

func TestBytesEncodeNilPointerBytes(t *testing.T) {
	var input *[]byte
	_, err := Bytes.Encode(input)
	if err == nil {
		t.Error("expected error for nil pointer")
	}
}

func TestBytesEncodeNilPointerString(t *testing.T) {
	var input *string
	_, err := Bytes.Encode(input)
	if err == nil {
		t.Error("expected error for nil pointer")
	}
}

func TestBytesEncodeUnsupported(t *testing.T) {
	_, err := Bytes.Encode(42)
	if err == nil {
		t.Error("expected error for unsupported type")
	}
}

func TestBytesDecodeBytes(t *testing.T) {
	var result []byte
	err := Bytes.Decode([]byte("hello"), &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "hello" {
		t.Errorf("unexpected result: %s", string(result))
	}
}

func TestBytesDecodeString(t *testing.T) {
	var result string
	err := Bytes.Decode([]byte("hello"), &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "hello" {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestBytesDecodeNilPointerBytes(t *testing.T) {
	var result *[]byte
	err := Bytes.Decode([]byte("hello"), result)
	if err == nil {
		t.Error("expected error for nil pointer")
	}
}

func TestBytesDecodeNilPointerString(t *testing.T) {
	var result *string
	err := Bytes.Decode([]byte("hello"), result)
	if err == nil {
		t.Error("expected error for nil pointer")
	}
}

func TestBytesDecodeUnsupported(t *testing.T) {
	var result int
	err := Bytes.Decode([]byte("42"), &result)
	if err == nil {
		t.Error("expected error for unsupported type")
	}
}

func TestBytesContentType(t *testing.T) {
	if Bytes.ContentType() != "application/octet-stream" {
		t.Errorf("unexpected content type: %s", Bytes.ContentType())
	}
}

func TestEncodeBytes(t *testing.T) {
	encoded, err := EncodeBytes([]byte("test"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(encoded) != "test" {
		t.Errorf("unexpected result: %s", string(encoded))
	}
}

func TestDecodeBytes(t *testing.T) {
	var result []byte
	err := DecodeBytes([]byte("test"), &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(result) != "test" {
		t.Errorf("unexpected result: %s", string(result))
	}
}

// ============== Protobuf Encoder Tests ==============

// Mock protobuf message
type mockProtoMessage struct {
	data      []byte
	marshalErr error
	unmarshalErr error
}

func (m *mockProtoMessage) Marshal() ([]byte, error) {
	if m.marshalErr != nil {
		return nil, m.marshalErr
	}
	return m.data, nil
}

func (m *mockProtoMessage) Unmarshal(data []byte) error {
	if m.unmarshalErr != nil {
		return m.unmarshalErr
	}
	m.data = data
	return nil
}

// Mock binary marshaler
type mockBinaryMarshaler struct {
	data      []byte
	marshalErr error
	unmarshalErr error
}

func (m *mockBinaryMarshaler) MarshalBinary() ([]byte, error) {
	if m.marshalErr != nil {
		return nil, m.marshalErr
	}
	return m.data, nil
}

func (m *mockBinaryMarshaler) UnmarshalBinary(data []byte) error {
	if m.unmarshalErr != nil {
		return m.unmarshalErr
	}
	m.data = data
	return nil
}

func TestProtobufEncodeNil(t *testing.T) {
	_, err := Protobuf.Encode(nil)
	if err == nil {
		t.Error("expected error for nil value")
	}
}

func TestProtobufEncodeProtoMarshaler(t *testing.T) {
	msg := &mockProtoMessage{data: []byte("proto")}
	encoded, err := Protobuf.Encode(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(encoded) != "proto" {
		t.Errorf("unexpected result: %s", string(encoded))
	}
}

func TestProtobufEncodeProtoMarshalerError(t *testing.T) {
	msg := &mockProtoMessage{marshalErr: errors.New("marshal error")}
	_, err := Protobuf.Encode(msg)
	if err == nil {
		t.Error("expected error")
	}
}

func TestProtobufEncodeBinaryMarshaler(t *testing.T) {
	msg := &mockBinaryMarshaler{data: []byte("binary")}
	encoded, err := Protobuf.Encode(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(encoded) != "binary" {
		t.Errorf("unexpected result: %s", string(encoded))
	}
}

func TestProtobufEncodeBinaryMarshalerError(t *testing.T) {
	msg := &mockBinaryMarshaler{marshalErr: errors.New("marshal error")}
	_, err := Protobuf.Encode(msg)
	if err == nil {
		t.Error("expected error")
	}
}

func TestProtobufEncodeUnsupported(t *testing.T) {
	_, err := Protobuf.Encode("string")
	if err == nil {
		t.Error("expected error for unsupported type")
	}
}

func TestProtobufDecodeEmpty(t *testing.T) {
	msg := &mockProtoMessage{}
	err := Protobuf.Decode([]byte{}, msg)
	if err == nil {
		t.Error("expected error for empty data")
	}
}

func TestProtobufDecodeProtoUnmarshaler(t *testing.T) {
	msg := &mockProtoMessage{}
	err := Protobuf.Decode([]byte("proto"), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(msg.data) != "proto" {
		t.Errorf("unexpected result: %s", string(msg.data))
	}
}

func TestProtobufDecodeProtoUnmarshalerError(t *testing.T) {
	msg := &mockProtoMessage{unmarshalErr: errors.New("unmarshal error")}
	err := Protobuf.Decode([]byte("proto"), msg)
	if err == nil {
		t.Error("expected error")
	}
}

func TestProtobufDecodeBinaryUnmarshaler(t *testing.T) {
	msg := &mockBinaryMarshaler{}
	err := Protobuf.Decode([]byte("binary"), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(msg.data) != "binary" {
		t.Errorf("unexpected result: %s", string(msg.data))
	}
}

func TestProtobufDecodeBinaryUnmarshalerError(t *testing.T) {
	msg := &mockBinaryMarshaler{unmarshalErr: errors.New("unmarshal error")}
	err := Protobuf.Decode([]byte("binary"), msg)
	if err == nil {
		t.Error("expected error")
	}
}

func TestProtobufDecodeUnsupported(t *testing.T) {
	var result string
	err := Protobuf.Decode([]byte("data"), &result)
	if err == nil {
		t.Error("expected error for unsupported type")
	}
}

func TestProtobufContentType(t *testing.T) {
	if Protobuf.ContentType() != "application/x-protobuf" {
		t.Errorf("unexpected content type: %s", Protobuf.ContentType())
	}
}

func TestEncodeProtobuf(t *testing.T) {
	msg := &mockProtoMessage{data: []byte("test")}
	encoded, err := EncodeProtobuf(msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(encoded) != "test" {
		t.Errorf("unexpected result: %s", string(encoded))
	}
}

func TestDecodeProtobuf(t *testing.T) {
	msg := &mockProtoMessage{}
	err := DecodeProtobuf([]byte("test"), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(msg.data) != "test" {
		t.Errorf("unexpected result: %s", string(msg.data))
	}
}

// ============== Registry Tests ==============

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	if r == nil {
		t.Fatal("NewRegistry returned nil")
	}

	// Should have JSON and Bytes registered
	if r.Get("application/json") != JSON {
		t.Error("JSON encoder not registered")
	}
	if r.Get("application/octet-stream") != Bytes {
		t.Error("Bytes encoder not registered")
	}
}

func TestRegistryRegister(t *testing.T) {
	r := NewRegistry()
	r.Register(Protobuf)

	if r.Get("application/x-protobuf") != Protobuf {
		t.Error("Protobuf encoder not registered")
	}
}

func TestRegistryGetFallback(t *testing.T) {
	r := NewRegistry()

	// Unknown content type should return fallback (JSON)
	enc := r.Get("unknown/type")
	if enc != JSON {
		t.Error("expected fallback to be JSON")
	}
}

func TestRegistrySetFallback(t *testing.T) {
	r := NewRegistry()
	r.SetFallback(Bytes)

	enc := r.Get("unknown/type")
	if enc != Bytes {
		t.Error("expected fallback to be Bytes")
	}
}

func TestRegistryEncode(t *testing.T) {
	r := NewRegistry()

	data := testStruct{Name: "test", Value: 1}
	encoded, err := r.Encode("application/json", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(encoded) != `{"name":"test","value":1}` {
		t.Errorf("unexpected result: %s", string(encoded))
	}
}

func TestRegistryDecode(t *testing.T) {
	r := NewRegistry()

	var result testStruct
	err := r.Decode("application/json", []byte(`{"name":"test","value":1}`), &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "test" || result.Value != 1 {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestDefaultRegistryEncode(t *testing.T) {
	data := testStruct{Name: "test", Value: 1}
	encoded, err := Encode("application/json", data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(encoded) != `{"name":"test","value":1}` {
		t.Errorf("unexpected result: %s", string(encoded))
	}
}

func TestDefaultRegistryDecode(t *testing.T) {
	var result testStruct
	err := Decode("application/json", []byte(`{"name":"test","value":1}`), &result)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "test" || result.Value != 1 {
		t.Errorf("unexpected result: %+v", result)
	}
}

// ============== Benchmarks ==============

func BenchmarkJSONEncode(b *testing.B) {
	data := testStruct{Name: "test", Value: 42}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = JSON.Encode(data)
	}
}

func BenchmarkJSONDecode(b *testing.B) {
	data := []byte(`{"name":"test","value":42}`)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result testStruct
		_ = JSON.Decode(data, &result)
	}
}

func BenchmarkBytesEncode(b *testing.B) {
	data := []byte("test data for encoding")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Bytes.Encode(data)
	}
}

func BenchmarkBytesDecode(b *testing.B) {
	data := []byte("test data for decoding")
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result []byte
		_ = Bytes.Decode(data, &result)
	}
}

func BenchmarkRegistryGet(b *testing.B) {
	r := NewRegistry()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Get("application/json")
	}
}
