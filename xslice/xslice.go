package xslice

import "math/rand"

// tooManyElement 是 RemoveRepeated 选择 loop / map 实现的元素数量阈值。
// 阈值以下走 O(n^2) 但常数小的双层 loop；以上走 O(n) 但带 map 分配的实现。
var tooManyElement = 1024

// Contain 判断 s 中是否含有指定元素 v。
// 与 Go 1.21+ 的 slices.Contains 等价，保留是为了与 boost 历史 API 命名习惯一致。
func Contain[T comparable](s []T, v T) bool {
	for _, ele := range s {
		if ele == v {
			return true
		}
	}
	return false
}

// SetAdd 集合添加（幂等）：如果 s 中不存在给定元素 v 则追加，重复元素不会被插入。
// 名字 "Set" 提示集合语义，与单纯的 append 区分。
func SetAdd[T comparable](s []T, v ...T) []T {
	for _, ele := range v {
		if !Contain(s, ele) {
			s = append(s, ele)
		}
	}
	return s
}

// Walk 遍历 s 将 f 应用到每一个元素，f 返回 (newValue, keep)；keep == false 时该元素不进入结果。
// T 约束 any，因为遍历本身不需要 comparable。
func Walk[T any](s []T, f func(T) (T, bool)) []T {
	out := make([]T, 0, len(s))
	for _, ele := range s {
		if ret, valid := f(ele); valid {
			out = append(out, ret)
		}
	}
	return out
}

// RemoveRepeated 移除重复元素，保序：相同元素只保留首次出现。
// 与 Go 1.21+ 的 slices.Compact 不同——slices.Compact 只去相邻重复，本函数全局去重。
//
// 实现策略：元素数 < tooManyElement 走双层 loop（无内存分配，常数小）；
// 大于阈值走 map 实现（O(n) 但有 map 分配）。
func RemoveRepeated[T comparable](s []T) []T {
	if len(s) == 0 {
		return s
	}
	if len(s) < tooManyElement {
		return removeRepeatByLoop(s)
	}
	return removeRepeatByMap(s)
}

// RemoveEmpty 移除空元素（== 类型零值的元素：0 / "" / nil 等）。
func RemoveEmpty[T comparable](s []T) []T {
	var zero T
	out := make([]T, 0, len(s))
	for _, ele := range s {
		if ele != zero {
			out = append(out, ele)
		}
	}
	return out
}

// Shuffle 原地打乱切片。
// 用全局 math/rand，可由调用方先 rand.Seed 控制确定性。
func Shuffle[T any](s []T) {
	for i := range s {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}
}

// ToAny 转换为 []any。
func ToAny[T any](s []T) []any {
	result := make([]any, len(s))
	for i, v := range s {
		result[i] = v
	}
	return result
}

// Last 返回切片最后一个元素；空切片返回 T 的零值。
func Last[T any](s []T) T {
	if len(s) == 0 {
		var t T
		return t
	}
	return s[len(s)-1]
}

// To 把 []F 通过 cb 映射成 []T。from 为 nil 时返回 nil。
func To[F any, T any](from []F, cb func(F) T) []T {
	if from == nil {
		return nil
	}
	to := make([]T, len(from))
	for i, f := range from {
		to[i] = cb(f)
	}
	return to
}

// removeRepeatByMap O(n) map 实现，适合 len(s) >= tooManyElement 的场景。
func removeRepeatByMap[T comparable](s []T) []T {
	out := make([]T, 0, len(s))
	tmp := make(map[T]struct{}, len(s))
	for _, ele := range s {
		l := len(tmp)
		tmp[ele] = struct{}{}
		if len(tmp) != l {
			out = append(out, ele)
		}
	}
	return out
}

// removeRepeatByLoop O(n^2) 双层 loop 实现，适合元素少 + 无 map 分配压力的场景。
func removeRepeatByLoop[T comparable](s []T) []T {
	out := make([]T, 0, len(s))
	for i := range s {
		flag := true
		for j := range out {
			if s[i] == out[j] {
				flag = false
				break
			}
		}
		if flag {
			out = append(out, s[i])
		}
	}
	return out
}
