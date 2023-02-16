package module

import (
	"context"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

type testPlugin struct {
}

func (*testPlugin) AfterRunModule(context.Context, Master) {
	fmt.Println("after_run_module")
}
func (*testPlugin) BeforeCloseModule(context.Context, Master) {
	fmt.Println("before_close_module")
}

type testModule struct {
}

func (*testModule) OnInit()           {}
func (*testModule) OnClose()          {}
func (*testModule) Run(chan struct{}) { time.Sleep(time.Millisecond * 10) }
func (*testModule) Name() string      { return "exit_after_10_millisecond" }

type testModuleNotStop struct {
}

func (*testModuleNotStop) OnInit()                     {}
func (*testModuleNotStop) OnClose()                    {}
func (*testModuleNotStop) Run(closeChan chan struct{}) { <-closeChan }
func (*testModuleNotStop) Name() string                { return "testModuleNotStop" }

func TestModuleStop(t *testing.T) {
	Convey("module should stop when no sub module in running state", t, func(c C) {
		m := New()
		m.Register(&testModule{})
		m.Run()
	})
	Convey("module should not stop when still have sub module in running state", t, func(c C) {
		Register(&testModule{}, &testModuleNotStop{})
		AttachPlugin(&testPlugin{})
		closed := make(chan struct{})
		go func() {
			Run()
			close(closed)
		}()
		timeout := false
		select {
		case <-time.After(time.Millisecond * 20):
			timeout = true
		case <-closed:
		}
		So(timeout, ShouldBeTrue)
		Stop()
		<-ShutdownNotify()
	})
}
