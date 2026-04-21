package xmap

func Diff[K comparable, V comparable](a, b map[K]V) map[K]V {
	if a == nil {
		return nil
	}

	if b == nil {
		return a
	}

	diff := make(map[K]V)
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			diff[k] = v
		}
	}
	return diff
}
