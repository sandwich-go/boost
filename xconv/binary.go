package xconv

import (
	"encoding/binary"
	"math"
)

// LittleEndianDecodeToFloat32 bytes to float32 use little-endian
func LittleEndianDecodeToFloat32(b []byte) float32 {
	return math.Float32frombits(binary.LittleEndian.Uint32(LeFillUpSize(b, 4)))
}

// LittleEndianDecodeToFloat64 bytes to float64 use little-endian
func LittleEndianDecodeToFloat64(b []byte) float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(LeFillUpSize(b, 8)))
}

// LittleEndianDecodeToInt64 bytes to int64 use little-endian
func LittleEndianDecodeToInt64(b []byte) int64 {
	return int64(binary.LittleEndian.Uint64(LeFillUpSize(b, 8)))
}

// LittleEndianDecodeToUint64 bytes to uint64 use little-endian
func LittleEndianDecodeToUint64(b []byte) uint64 {
	return binary.LittleEndian.Uint64(LeFillUpSize(b, 8))
}

// LittleEndianEncodeFromFloat32 float32 to bytes (little-endian, 4 bytes)
func LittleEndianEncodeFromFloat32(v float32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, math.Float32bits(v))
	return b
}

// LittleEndianEncodeFromFloat64 float64 to bytes (little-endian, 8 bytes)
func LittleEndianEncodeFromFloat64(v float64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, math.Float64bits(v))
	return b
}

// LittleEndianEncodeFromInt64 int64 to bytes (little-endian, 8 bytes)
func LittleEndianEncodeFromInt64(v int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(v))
	return b
}

// LittleEndianEncodeFromUint64 uint64 to bytes (little-endian, 8 bytes)
func LittleEndianEncodeFromUint64(v uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, v)
	return b
}
