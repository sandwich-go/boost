package xconv

// iInt64 is used for type assert api for Int64().
type iInt64 interface {
	Int64() int64
}

// iUint64 is used for type assert api for Int64().
type iUint64 interface {
	Uint64() uint64
}

// iFloat64 is used for type assert api for Float64().
type iFloat64 interface {
	Float64() float64
}

// iFloat32 is used for type assert api for Float32().
type iFloat32 interface {
	Float32() float32
}

// iString is used for type assert api for String().
type iString interface {
	String() string
}

// iBool is used for type assert api for Bool().
type iBool interface {
	Bool() bool
}
