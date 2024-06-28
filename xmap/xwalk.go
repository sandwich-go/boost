package xmap

import (
	"cmp"
	"sort"
)

// WalkMapDeterministic 有序遍历map
// walkFunc 函数返回 false，停止遍历
func WalkMapDeterministic[K, V cmp.Ordered](in map[K]V, walkFunc func(k K, v V) bool) {
	var keys = make([]K, 0, len(in))
	for k := range in {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return cmp.Compare(keys[i], keys[j]) < 0 })
	for _, k := range keys {
		if walkFunc(k, in[k]) {
			continue
		}
		break
	}
}
