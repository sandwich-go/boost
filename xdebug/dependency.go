package xdebug

import (
	"runtime/debug"
	"strings"

	"github.com/coreos/go-semver/semver"

	"github.com/sandwich-go/boost"
)

type dependency interface {
	GetPath() string
	GetRequireVersion() string
	WarnString() string
	GoVersionDisuse() string
}

// dependencies 依赖的包
var dependencies = make([]dependency, 0)

func registerDependency(d dependency) {
	dependencies = append(dependencies, d)
}

func getDependenciesFromBuildInfo() (map[string]semver.Version, *semver.Version, bool) {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		// read fail, don't check require dependencies
		return nil, nil, false
	}
	if bi == nil {
		return nil, nil, true
	}

	ver, _ := semver.NewVersion(strings.TrimPrefix(bi.GoVersion, "go"))

	var out = make(map[string]semver.Version)
	for _, dep := range bi.Deps {
		if v, _ := semver.NewVersion(dep.Version); v != nil {
			out[dep.Path] = *v
		} else {
			out[dep.Path] = semver.Version{}
		}
	}
	return out, ver, true
}

func checkRequireDependency(goVer *semver.Version, deps map[string]semver.Version, requireDependency dependency) bool {
	if goVersionDisuse := requireDependency.GoVersionDisuse(); goVersionDisuse != "" && goVer != nil {
		goDisuseSemVer, _ := semver.NewVersion(goVersionDisuse)
		if goDisuseSemVer != nil && goVer.Compare(*goDisuseSemVer) >= 0 {
			// 如果当前的 go 版本大于等于放弃版本，则不校验
			return true
		}
	}
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
	deps, ver, ok := getDependenciesFromBuildInfo()
	if !ok {
		return
	}
	for _, v := range dependencies {
		if checkRequireDependency(ver, deps, v) {
			continue
		}
		boost.LogWarn(v.WarnString())
	}
}
