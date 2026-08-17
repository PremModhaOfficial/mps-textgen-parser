package types

const (
	TypeSTRING = "STRING"

	TypeFLOAT64 = "FLOAT64"

	TypeINT64 = "INT64"

	TypeByte = "BYTE"

	TypeUnknown = "UNKNOWN"
)

// Constraints for generics
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

type Integer interface {
	Signed | Unsigned
}

type Float interface {
	~float32 | ~float64
}

type Numeric interface {
	Integer | Float
}

type Char interface {
	string
}

func IsInteger[T Integer](v T) T {
	return v // Only compiles if T satisfies Integer constraint
}

func IsFloat[T Float](v T) T {

	return v // Only compiles if T satisfies Float constraint
}

func IsString[T string](v T) T {

	return v // Only compiles if T satisfies String constraint
}

func IsBoolean[T bool](v T) T {

	return v // Only compiles if T satisfies Boolean constraint
}

func IsSigned[T Signed](v T) T {

	return v
}

func IsUnsigned[T Unsigned](v T) T {

	return v
}
