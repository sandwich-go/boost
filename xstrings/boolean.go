package xstrings

import "strings"

// Empty strings.
var emptyStringMap = map[string]struct{}{
	"":      {},
	"0":     {},
	"n":     {},
	"no":    {},
	"off":   {},
	"false": {},
}

// IsFalse 判断command解析获取的数据是否为false
func IsFalse(v string) bool {
	_, ok := emptyStringMap[strings.ToLower(string(v))]
	return ok
}

// IsTrue 判断command解析获取的数据是否为true
func IsTrue(v string) bool { return !IsFalse(v) }
