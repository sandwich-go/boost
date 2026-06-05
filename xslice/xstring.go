package xslice

import "strings"

// ContainEqualFold 判断 s 中是否含有指定字符串 v（不区分大小写）。
//
// 标准库无直接对应：可用 slices.ContainsFunc(s, func(e string) bool { return strings.EqualFold(e, v) }) 等价表达。
func ContainEqualFold(s []string, v string) bool {
	for _, ele := range s {
		if strings.EqualFold(ele, v) {
			return true
		}
	}
	return false
}

// AddPrefix 给每个元素添加前缀，返回新切片。
func AddPrefix(s []string, prefix string) []string {
	out := make([]string, 0, len(s))
	for _, ele := range s {
		out = append(out, prefix+ele)
	}
	return out
}

// AddSuffix 给每个元素添加后缀，返回新切片。
func AddSuffix(s []string, suffix string) []string {
	out := make([]string, 0, len(s))
	for _, ele := range s {
		out = append(out, ele+suffix)
	}
	return out
}
