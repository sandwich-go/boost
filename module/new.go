package module

import (
	"github.com/rs/xid"
	"github.com/sandwich-go/boost/xpanic"
)

type simpleModule struct {
	cc  *ModuleOptions
	run func(closeChan chan struct{}) (clear func())
}

func (m *simpleModule) OnInit() {
	if m.cc.OnInit != nil {
		m.cc.OnInit()
	}
}
func (m *simpleModule) Run(closeChan chan struct{}) {
	clearFunc := m.run(closeChan)
	defer func() {
		if clearFunc != nil {
			clearFunc()
		}
	}()
	<-closeChan
}
func (m *simpleModule) Name() string { return m.cc.Name }
func (m *simpleModule) OnClose() {
	if m.cc.OnClose != nil {
		m.cc.OnClose()
	}
}

// NewModule 提供一个便捷创建Module的方式
func NewModule(run func(closeChan chan struct{}) (clear func()), opts ...ModuleOption) Module {
	cc := NewModuleOptions(opts...)
	if cc.Name == "" {
		cc.Name = "module-" + xid.New().String()
	}
	xpanic.WhenTrue(run == nil, "run should not be nil")
	return &simpleModule{
		cc:  cc,
		run: run,
	}
}
