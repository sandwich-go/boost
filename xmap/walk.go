package xmap

import "sort"

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
