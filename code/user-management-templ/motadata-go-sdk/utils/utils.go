package utils

import (
	"reflect"
	"unsafe"
)

// IsZeroValue checks if a value is the zero value of its type
func IsZeroValue[T any](v T) bool {

	return reflect.ValueOf(&v).Elem().IsZero()
}

func StringToBytes(value string) []byte {

	return []byte(value)
}

func StringToUnsafeBytes(value string) []byte {

	return unsafe.Slice(unsafe.StringData(value), len(value))
}
