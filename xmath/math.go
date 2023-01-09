package xmath

func MaxInt8(a, b int8) int8 {
	if a > b {
		return a
	}
	return b
}

func MinInt8(a, b int8) int8 {
	if a < b {
		return a
	}
	return b
}

func MinInt32(a, b int32) int32 {
	if a > b {
		return b
	}
	return a
}

func MaxInt32(a, b int32) int32 {
	if a < b {
		return b
	}
	return a
}

func MinInt64(a, b int64) int64 {
	if a > b {
		return b
	}
	return a
}

func MaxInt64(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}

func MinUint32(a, b uint32) uint32 {
	if a > b {
		return b
	}
	return a
}

func MaxUint32(a, b uint32) uint32 {
	if a < b {
		return b
	}
	return a
}

func MinUint64(a, b uint64) uint64 {
	if a > b {
		return b
	}
	return a
}

func MaxUint64(a, b uint64) uint64 {
	if a < b {
		return b
	}
	return a
}

func MinInt(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func MaxInt(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func MinFloat32(a, b float32) float32 {
	if a > b {
		return b
	}
	return a
}

func MaxFloat32(a, b float32) float32 {
	if a < b {
		return b
	}
	return a
}

func MinFloat64(a, b float64) float64 {
	if a > b {
		return b
	}
	return a
}

func MaxFloat64(a, b float64) float64 {
	if a < b {
		return b
	}
	return a
}

func AbsInt32(n int32) int32 {
	if n < 0 {
		return -n
	}
	return n
}

func AbsInt64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

func Float64Equals(a, b float64) bool {
	return (a-b) < EPSILON64 && (b-a) < EPSILON64
}

func Float32Equals(a, b float32) bool {
	return (a-b) < float32(EPSILON32) && (b-a) < float32(EPSILON32)
}

func IsZeroFloat64(v float64) bool {
	return Float64Equals(v, 0)
}

func IsZeroFloat32(v float32) bool {
	return Float32Equals(v, 0)
}

// IsBelowZeroFloat64 v == 0 时也返回true
func IsBelowZeroFloat64(v float64) bool {
	return (v - 0) < EPSILON64
}

// IsBelowZeroFloat32 v == 0 时也返回true
func IsBelowZeroFloat32(v float32) bool {
	return (v - 0) < EPSILON32
}

func EffectZeroLimitInt64(v, change int64) int64 {
	v += change
	if v < 0 {
		v = 0
	}
	return v
}

func EffectZeroLimitInt32(v, change int32) int32 {
	v += change
	if v < 0 {
		v = 0
	}
	return v
}
