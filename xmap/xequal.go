package xmap

// Equal 判断两个 map 的 key/value 是否完全相等，区分 nil 与空 map：
//   - nil == nil 相等
//   - nil 与 len(b)==0 的非 nil map 不相等
//
// 这与 Go 1.21+ 标准库 maps.Equal 不同——maps.Equal 视 nil == empty。
// 如需 nil == empty 语义请用标准库 maps.Equal。
func Equal[K, V comparable](a, b map[K]V) bool {
	if a == nil && b == nil {
		return true
	}
	if (a != nil && b == nil) || (b != nil && a == nil) || (len(a) != len(b)) {
		return false
	}
	for k, v := range a {
		v1, ok := b[k]
		if !ok || v != v1 {
			return false
		}
	}
	return true
}
