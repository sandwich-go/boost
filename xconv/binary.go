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
