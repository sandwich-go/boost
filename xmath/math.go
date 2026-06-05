package xmath

import "math"

const (
	EPSILON64 float64 = 0.00000001
	EPSILON32 float32 = 0.00000001
)

// Integer limit values.
const (
	ConstMaxInt    = math.MaxInt
	ConstMinInt    = math.MinInt
	ConstMaxUint   = math.MaxUint
	ConstMaxInt8   = math.MaxInt8
	ConstMinInt8   = math.MinInt8
	ConstMaxInt16  = math.MaxInt16
	ConstMinInt16  = math.MinInt16
	ConstMaxInt32  = math.MaxInt32
	ConstMinInt32  = math.MinInt32
	ConstMaxInt64  = math.MaxInt64
	ConstMinInt64  = math.MinInt64
	ConstMaxUint8  = math.MaxUint8
	ConstMaxUint16 = math.MaxUint16
	ConstMaxUint32 = math.MaxUint32
	ConstMaxUint64 = math.MaxUint64
)

// Float64Equals 判断 float64 是否相等（EPSILON 近似）。
func Float64Equals(a, b float64) bool {
	return (a-b) < EPSILON64 && (b-a) < EPSILON64
}

// Float32Equals 判断 float32 是否相等（EPSILON 近似）。
func Float32Equals(a, b float32) bool {
	return (a-b) < EPSILON32 && (b-a) < EPSILON32
}

// IsZeroFloat64 判断 float64 是否为零（EPSILON 近似）。
func IsZeroFloat64(v float64) bool {
	return Float64Equals(v, 0)
}

// IsZeroFloat32 判断 float32 是否为零（EPSILON 近似）。
func IsZeroFloat32(v float32) bool {
	return Float32Equals(v, 0)
}

// IsBelowZeroFloat64 判断 float64 是否 <= 0（EPSILON 近似，v == 0 时也返回 true）。
func IsBelowZeroFloat64(v float64) bool {
	return (v - 0) < EPSILON64
}

// IsBelowZeroFloat32 判断 float32 是否 <= 0（EPSILON 近似，v == 0 时也返回 true）。
func IsBelowZeroFloat32(v float32) bool {
	return (v - 0) < EPSILON32
}
