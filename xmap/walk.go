package xmap

import "sort"

// WalkMapDeterministic 有序遍历map
func WalkMapDeterministic(in map[string]string, walkFunc func(k string, v string) bool) {
	var keys []string
	for k := range in {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if walkFunc(k, in[k]) {
			continue
		}
		break
	}
}
