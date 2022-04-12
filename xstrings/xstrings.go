package xstrings

import (
	"strings"
	"unicode"
)

// FirstUpper 首字符大写
func FirstUpper(s string) string {
	var result string
	for i, word := range s {
		w := word
		if i == 0 {
			w = unicode.ToUpper(w)
		}
		result += string(w)
	}
	return result
}

// FirstLower 首字符小写
func FirstLower(s string) string {
	var result string
	for i, word := range s {
		w := word
		if i == 0 {
			w = unicode.ToLower(w)
		}
		result += string(w)
	}
	return result
}

// HasPrefixIgnoreCase 前缀匹配
func HasPrefixIgnoreCase(str, prefix string) bool {
	return strings.HasPrefix(strings.ToLower(str), strings.ToLower(prefix))
}

// TrimPrefixIgnoreCase 移除前缀,不区分大小写
func TrimPrefixIgnoreCase(str, prefix string) string {
	if !HasPrefixIgnoreCase(str, prefix) {
		return str
	}
	return strings.TrimSpace(str[len(prefix):])
}

// CompareIgnoreCase 比较字符串相等，不区分大小写
func CompareIgnoreCase(s1, s2 string) bool {
	return strings.EqualFold(s1, s2)
}
