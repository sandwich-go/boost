package xslice

import "math/rand"

// Contain s 中是否含有指定元素 v
func Contain[T comparable](s []T, v T) bool {
	for _, ele := range s {
		if ele == v {
			return true
		}
	}
	return false
}

// Deprecated: 使用 SetAdd
// Add 如果 s 中不存在给定的元素 v 则添加
func Add[T comparable](s []T, v ...T) []T {
	for _, ele := range v {
		if !Contain(s, ele) {
			s = append(s, ele)
		}
	}
	return s
}

// SetAdd 如果 s 中不存在给定的元素 v 则添加
func SetAdd[T comparable](s []T, v ...T) []T {
	for _, ele := range v {
		if !Contain(s, ele) {
			s = append(s, ele)
		}
	}
	return s
}

// Walk 遍历 s 将 f 应用到每一个元素，返回更新后的数据
func Walk[T comparable](s []T, f func(T) (T, bool)) []T {
	out := make([]T, 0, len(s))
	for _, ele := range s {
		if ret, valid := f(ele); valid {
			out = append(out, ret)
		}
	}
	return out
}

// RemoveRepeated 移除重复元素
func RemoveRepeated[T comparable](s []T) []T {
	if len(s) == 0 {
		return s
	}
	if len(s) < tooManyElement {
		return removeRepeatByLoop(s)
	} else {
		return removeRepeatByMap(s)
	}
}

// RemoveEmpty 移除空元素
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

func removeRepeatByMap[T comparable](s []T) []T {
	out := make([]T, 0, len(s))
	tmp := make(map[T]struct{})
	for _, ele := range s {
		l := len(tmp)
		tmp[ele] = struct{}{}
		if len(tmp) != l {
			out = append(out, ele)
		}
	}
	return out
}

func removeRepeatByLoop[T comparable](s []T) []T {
	out := make([]T, 0, len(s))
	flag := true
	for i := range s {
		flag = true
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

// Shuffle 数组打乱
func Shuffle[T comparable](s []T) {
	for i := range s {
		j := rand.Intn(i + 1)
		s[i], s[j] = s[j], s[i]
	}
}

// ToAny 转换为[]any
func ToAny[T any](s []T) []any {
	result := make([]any, len(s))
	for i, v := range s {
		result[i] = v
	}
	return result
}

// Last 最后一个元素
func Last[T any](s []T) T {
	if len(s) == 0 {
		var t T
		return t
	}
	return s[len(s)-1]
}

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
