package plugin

import (
	"errors"
	"github.com/sandwich-go/boost/xpanic"
	"reflect"
)

var ErrInvaliPlugin = errors.New("invalid plugin")

// Plugin 插件
type Plugin interface{}

// Container 插件容器
type Container interface {
	// Add 添加 Plugin
	Add(Plugin) error
	// MustAdd 添加 Plugin，若为非法的 Plugin，则 panic
	MustAdd(Plugin)
	// Range 遍历所有的 Plugin 并执行函数，若返回 false，则中止遍历
	Range(func(Plugin) bool)
	// Size 插件数量
	Size() int
}

type container struct {
	plugins []Plugin
	types   []interface{}
}

func New(types ...interface{}) Container { return &container{types: types} }
func (p *container) MustAdd(pg Plugin)   { xpanic.WhenError(p.Add(pg)) }
func (p *container) Size() int           { return len(p.plugins) }
func (p *container) Add(pg Plugin) error {
	var pType = reflect.TypeOf(pg)
	for _, v := range p.types {
		if pType.Implements(reflect.TypeOf(v).Elem()) {
			p.plugins = append(p.plugins, pg)
			return nil
		}
	}
	return ErrInvaliPlugin
}

func (p *container) Range(f func(Plugin) bool) {
	for _, plugin := range p.plugins {
		if !f(plugin) {
			return
		}
	}
}
