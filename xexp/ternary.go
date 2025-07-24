package xexp

// Ternary 类似三元表达式 condition ? a : b。
func Ternary[T any](condition bool, a, b T) T {
	if condition {
		return a
	}
	return b
}

// TernaryFunc 使用所选的元素作为参数来执行给定的函数
func TernaryFunc[T any, R any](condition bool, a, b T, f func(T) R) R {
	if condition {
		return f(a)
	}
	return f(b)
}
