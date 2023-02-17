package xdebug

import (
	"github.com/sandwich-go/boost"
	"runtime/debug"
)

// dependencies 依赖的包
var dependencies = map[string]string{
	"go.uber.org/automaxprocs": `for the best performance, please blank import the package 'go.uber.org/automaxprocs'`,
}

// CheckDependencies 检查依赖
func CheckDependencies() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	var nowDeps = make(map[string]struct{})
	if info != nil {
		for _, dep := range info.Deps {
			nowDeps[dep.Path] = struct{}{}
		}
	}
	for k, v := range dependencies {
		if _, ok = nowDeps[k]; !ok {
			boost.LogWarn(v)
		}
	}
}
