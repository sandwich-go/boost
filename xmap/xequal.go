package xmap

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
