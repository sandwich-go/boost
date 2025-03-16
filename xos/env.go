package xos

import (
	"os"
	"strings"
)

// EnvGet 获取环境变量值，如果未找到，则返回 def 值
func EnvGet(key string, def ...string) string {
	v, ok := os.LookupEnv(key)
	if !ok && len(def) > 0 {
		return def[0]
	}
	return v
}

// EnvGetCaseInsensitive 查找环境变量，大小写不敏感
// 如果未找到，则返回 def 值，如未提供def值则返回空
func EnvGetCaseInsensitive(key string, def ...string) string {
	upperKey := strings.ToUpper(key)
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) == 2 && strings.ToUpper(strings.TrimSpace(pair[0])) == upperKey {
			return pair[1]
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}
