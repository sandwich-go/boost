package xmap

import (
	"cmp"
	"sort"
)

// WalkMapDeterministic 按 key 升序遍历 map（确定性遍历）。
// walkFunc 返回 false 时停止遍历。
//
// K 约束 cmp.Ordered（有序：可用 < 比较）；V 约束 any（任意类型，
// 包含 interface{} 与函数等不可比较类型）。
func WalkMapDeterministic[K cmp.Ordered, V any](in map[K]V, walkFunc func(k K, v V) bool) {
	keys := make([]K, 0, len(in))
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
