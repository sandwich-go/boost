package buildinfo

import "runtime/debug"

const (
	AutoMaxProcs = "go.uber.org/automaxprocs"
)

// CheckDependency 判断是否依赖指定module
func CheckDependency(pkgName string) bool {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, dep := range info.Deps {
			if dep.Path == pkgName {
				return true
			}
		}
	}
	return false
}
