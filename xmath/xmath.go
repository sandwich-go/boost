package xmath

import (
	"cmp"
	"golang.org/x/exp/constraints"
)

// Max 返回大值
func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Min 返回小值
func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Abs 返回绝对值
func Abs[T constraints.Unsigned](v T) T {
	if v < 0 {
		return -v
	}
	return v
}

// EffectZeroLimit 加 change 值，返回值，该值不会小于0
func EffectZeroLimit[T constraints.Unsigned](v, change T) T {
	v += change
	if v < 0 {
		v = 0
	}
	return v
}
