package template2

import "sort"

type DependencyConfig struct {
	Name           string
	Path           string
	RequireVersion string
	WarnString     string
}

func GetDependencyTPLArgs() interface{} {
	var dependencyConfigs = make([]DependencyConfig, 0)
	dependencyConfigs = append(dependencyConfigs, DependencyConfig{
		Name:           "automaxprocs",
		Path:           "go.uber.org/automaxprocs",
		RequireVersion: "v1.5.1",
		WarnString:     "fmt.Sprintf(`for the best performance, please blank import the package '%s@%s'`, d.GetPath(), d.GetRequireVersion())",
	})
	sort.Slice(dependencyConfigs, func(i, j int) bool {
		return dependencyConfigs[i].Name < dependencyConfigs[j].Name
	})
	return map[string]interface{}{"DependencyConfigs": dependencyConfigs}
}

const DependencyTPL = `// Code generated by tools. DO NOT EDIT.
package xdebug

import "fmt"

func init() {
{{- range $dependencyConfig := .DependencyConfigs }}
	registerDependency({{ $dependencyConfig.Name }}{})
{{- end }}
}

{{ range $dependencyConfig := .DependencyConfigs }}
type {{ $dependencyConfig.Name }} struct{}

func (d {{ $dependencyConfig.Name }}) GetPath() string           { return "{{ $dependencyConfig.Path }}" }
func (d {{ $dependencyConfig.Name }}) GetRequireVersion() string { return "{{ $dependencyConfig.RequireVersion }}" }
func (d {{ $dependencyConfig.Name }}) WarnString() string { return {{ $dependencyConfig.WarnString | Unescaped }}}
{{ end }}
`
