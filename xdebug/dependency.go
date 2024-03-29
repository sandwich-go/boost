package xdebug

import (
	"github.com/coreos/go-semver/semver"
	"github.com/sandwich-go/boost"
	"runtime/debug"
)

type dependency interface {
	GetPath() string
	GetRequireVersion() string
	WarnString() string
}

// dependencies 依赖的包
var dependencies = make([]dependency, 0)

func registerDependency(d dependency) {
	dependencies = append(dependencies, d)
}

func getDependenciesFromBuildInfo() (map[string]semver.Version, bool) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		// read fail, don't check require dependencies
		return nil, false
	}
	if bi == nil {
		return nil, true
	}
	var out = make(map[string]semver.Version)
	for _, dep := range bi.Deps {
		if v, _ := semver.NewVersion(dep.Version); v != nil {
			out[dep.Path] = *v
		} else {
			out[dep.Path] = semver.Version{}
		}
	}
	return out, true
}

func checkRequireDependency(deps map[string]semver.Version, requireDependency dependency) bool {
	// has require dependency?
	depSemVer, ok := deps[requireDependency.GetPath()]
	if !ok {
		return false
	}
	// compare dependency version
	requireVer := requireDependency.GetRequireVersion()
	if len(requireVer) == 0 {
		return true
	}
	requireSemVer, _ := semver.NewVersion(requireVer)
	if requireSemVer == nil {
		return true
	}
	return requireSemVer.LessThan(depSemVer) || requireSemVer.Equal(depSemVer)
}

// CheckRequireDependencies 检查依赖
func CheckRequireDependencies() {
	deps, ok := getDependenciesFromBuildInfo()
	if !ok {
		return
	}
	for _, v := range dependencies {
		if checkRequireDependency(deps, v) {
			continue
		}
		boost.LogWarn(v.WarnString())
	}
}
