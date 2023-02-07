package xos

import "os"

// EnvGet 获取环境变量值，如果未找到，则返回 def 值
func EnvGet(key string, def ...string) string {
	v, ok := os.LookupEnv(key)
	if !ok && len(def) > 0 {
		return def[0]
	}
	return v
}
