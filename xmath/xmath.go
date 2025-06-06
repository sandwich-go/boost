package xmath

import (
	"cmp"
	"golang.org/x/exp/constraints"
	"math/rand/v2"
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

// Disturb 随机值，值的范围 [n*(100-percent)/100, n*(100+percent)/100]
func Disturb[N constraints.Integer](n N, percent N) N {
	w := rand.N(n * percent / 100)
	if w%2 == 0 {
		return n + w
	}
	return n - w
}
