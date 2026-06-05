package xmath

import (
	"cmp"
	"math"
	"math/rand/v2"
	"unsafe"

	"golang.org/x/exp/constraints"
)

// Max 返回较大值。NaN 处理沿用 Go 内建 max 语义（NaN 不传染，依赖比较顺序）。
// 如需 IEEE 754 NaN 传染语义，请用 [MaxFloat]。
func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Min 返回较小值。NaN 处理沿用 Go 内建 min 语义（NaN 不传染，依赖比较顺序）。
// 如需 IEEE 754 NaN 传染语义，请用 [MinFloat]。
func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// MaxFloat 返回较大浮点值，IEEE 754 语义（NaN 传染、Inf 处理）。
// 底层走 math.Max；与 [Max] 不同，[Max] 是 if a>b 比较，NaN 不会被传染。
func MaxFloat[T ~float32 | ~float64](a, b T) T {
	return T(math.Max(float64(a), float64(b)))
}

// MinFloat 返回较小浮点值，IEEE 754 语义（NaN 传染、Inf 处理）。
// 底层走 math.Min；与 [Min] 不同，[Min] 是 if a<b 比较，NaN 不会被传染。
func MinFloat[T ~float32 | ~float64](a, b T) T {
	return T(math.Min(float64(a), float64(b)))
}

// Abs 返回有符号/无符号整数的绝对值。
// 注意：math.MinInt / MinInt8 等最小值取反会溢出（不做溢出保护，与 if v<0 直接 -v 等价）。
// 对无符号类型 v 始终 >= 0，分支永远不命中。
func Abs[T constraints.Unsigned | constraints.Signed](v T) T {
	if v < 0 {
		return -v
	}
	return v
}

// AbsFloat 返回浮点绝对值，使用 EPSILON 近似判负：
// 当 v < EPSILON 时取反（IsBelowZeroFloat* 为 true 即取反）。
//
// 这与 if v < 0 不等价：
//   - v = -0.5*EPSILON 仍会返回 -v
//   - v = 0 时 IsBelowZeroFloat*(0) == true，会返回 -0
//
// 因 EPSILON32 / EPSILON64 精度不同，按 T 的实际字节大小路由到对应 EPSILON
// （unsafe.Sizeof 在编译期确定，零运行时开销）。
func AbsFloat[T ~float32 | ~float64](v T) T {
	if isBelowZeroFloat(v) {
		return -v
	}
	return v
}

// EffectZeroLimit 加 change 后返回整数结果，结果不会小于 0。
// 对无符号类型 v + change 不会变负（前提 change 也是无符号），分支永远不命中。
func EffectZeroLimit[T constraints.Unsigned | constraints.Signed](v, change T) T {
	v += change
	if v < 0 {
		v = 0
	}
	return v
}

// EffectZeroLimitFloat 加 change 后返回浮点结果，结果不会小于 0；
// 使用 EPSILON 近似判负（与 [AbsFloat] 一致）。
func EffectZeroLimitFloat[T ~float32 | ~float64](v, change T) T {
	v += change
	if isBelowZeroFloat(v) {
		v = 0
	}
	return v
}

// Disturb 随机值，值的范围 [n*(100-percent)/100, n*(100+percent)/100]。
func Disturb[N constraints.Integer](n N, percent N) N {
	w := rand.N(n * percent / 100)
	if w%2 == 0 {
		return n + w
	}
	return n - w
}

// isBelowZeroFloat 按 T 的真实精度路由到对应 EPSILON 比较。
// 用 unsafe.Sizeof 在编译期决定走 EPSILON32 还是 EPSILON64，零运行时开销。
// 语义与导出版 [IsBelowZeroFloat32] / [IsBelowZeroFloat64] 字面一致。
func isBelowZeroFloat[T ~float32 | ~float64](v T) bool {
	if unsafe.Sizeof(v) == 4 {
		return float32(v)-0 < EPSILON32
	}
	return float64(v)-0 < EPSILON64
}
