package xos

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

// FuncFullNameUsingReflect 使用反射获取函数名
func FuncFullNameUsingReflect(f interface{}) (string, error) {
	t := reflect.TypeOf(f).Kind()
	if t != reflect.Func {
		return "", fmt.Errorf("args must be func, now %v", t)
	}
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), nil
}

// FuncBaseNameUsingReflect 使用反射获取基础函数名
func FuncBaseNameUsingReflect(f interface{}) (string, error) {
	n, err := FuncFullNameUsingReflect(f)
	if err != nil {
		return "", err
	}
	ss := strings.Split(filepath.Base(n), ".")
	return ss[len(ss)-1], nil
}
