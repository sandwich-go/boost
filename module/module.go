package module

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	processShutdownNotifyChan = make(chan struct{})
	processShutdownNotifyOnce sync.Once
	processShutdownSignal     os.Signal
)

// ProcessShutdownNotify 进程退出信号通知
func ProcessShutdownNotify() chan struct{} {
	processShutdownNotifyOnce.Do(func() {
		c := make(chan os.Signal, 1)
		c <- syscall.SIGHUP
		waitIncoming := make(chan struct{})
		go func() {
			<-c //读取syscall.SIGHUP
			// os.Kill cannot be trapped
			signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM)
			close(waitIncoming) // 通知已经监听，防止过早退出ProcessShutdownNotify，逻辑执行在监听之前
			processShutdownSignal = <-c
			close(processShutdownNotifyChan)
		}()
		<-waitIncoming
	})
	return processShutdownNotifyChan
}

// Plugin 插件
type Plugin interface {
	// AfterRunModule 运行所有的 Module 之后被调用
	AfterRunModule(context.Context, Master)
	// BeforeCloseModule 关闭所有的 Module 之前被调用
	BeforeCloseModule(context.Context, Master)
}

// Master Module 管理器
type Master interface {
	// AttachPlugin 挂接 Plugin
	AttachPlugin(plugins ...Plugin)
	// Register 注册多个 Module
	Register(ms ...Module)
	// Stop 主动停止 Master
	Stop(reason ...string)
	// Run 运行入口，进程会堵塞在这里直到收到停止信号，可以指定要运行的 Module 列表
	Run(ms ...Module)
	// RunWithCloseTimeout 参考 Run,扩充了关闭超时支持，防止逻辑层堵塞导致进程关闭失败
	RunWithCloseTimeout(closeTimeout time.Duration, ms ...Module)
	// ShutdownNotify Master 停止的通知信号
	ShutdownNotify() chan struct{}
	// RunModule 运行一个单独的 Module
	RunModule(md Module)
}

// Module 进程中的所有 Module 都需注册到 module 中并由其管理
// 模块启动时回调 OnInit，关闭时通过 closeChan 通知 Module 业务逻辑,然后回调 OnClose
type Module interface {
	// OnInit 模块启动时回调
	OnInit()
	// OnClose Run 监听到 closeChan 返回后回调
	OnClose()
	// Run 需要在 Run 逻辑中监听 closeChan 通知，收到通知再返回，否则 master 会认为该 Module 已退出
	Run(closeChan chan struct{})
	// Name 模块名称
	Name() string
}

// defaultMaster 默认的 Master
var defaultMaster = New()

// AttachPlugin 挂接 Plugin
func AttachPlugin(plugins ...Plugin) { defaultMaster.AttachPlugin(plugins...) }

// Register 注册多个 Module
func Register(ms ...Module) { defaultMaster.Register(ms...) }

// Stop 主动停止 Master
func Stop(reason ...string) { defaultMaster.Stop(reason...) }

// Run 运行入口，进程会堵塞在这里直到收到停止信号，可以指定要运行的 Module 列表
func Run(ms ...Module) { defaultMaster.Run(ms...) }

// RunWithCloseTimeout 参考 Run,扩充了关闭超时支持，防止逻辑层堵塞导致进程关闭失败
func RunWithCloseTimeout(closeTimeout time.Duration, ms ...Module) {
	defaultMaster.RunWithCloseTimeout(closeTimeout, ms...)
}

// ShutdownNotify Master 停止的通知信号
func ShutdownNotify() chan struct{} { return defaultMaster.ShutdownNotify() }

// RunModule 运行一个单独的 Module
func RunModule(md Module) { defaultMaster.RunModule(md) }
