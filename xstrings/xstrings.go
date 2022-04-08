package xstrings

import "unicode"

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
