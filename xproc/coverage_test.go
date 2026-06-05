package xproc

import (
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：补 xproc 包未覆盖的：
//   - parseCommand 各分支（unix 直接返；windows 引号解析略 skip）
//   - Manager 全部 13 个方法（之前 0%）
//   - Process Pid / Run / Start / NewProcess 各分支
//   - ShellRun / Run 顶层 helper

func TestParseCommand_Unix(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("unix path only")
	}
	Convey("unix 下 parseCommand 直接把命令塞 args[0]", t, func() {
		args := parseCommand("ls -la /tmp")
		So(args, ShouldResemble, []string{"ls -la /tmp"})
	})
}

func TestManager_BasicAPI(t *testing.T) {
	Convey("Manager 空状态：Size=0 / Pids=nil / Processes 空", t, func() {
		m := NewManager()
		So(m.Size(), ShouldEqual, 0)
		So(m.Pids(), ShouldBeNil)
		So(m.Processes(), ShouldBeEmpty)
		So(m.GetProcess(0), ShouldBeNil)
	})

	Convey("Manager.NewProcess 关联 Manager 自身", t, func() {
		m := NewManager()
		p := m.NewProcess("/bin/sh")
		So(p, ShouldNotBeNil)
		So(p.Manager, ShouldEqual, m)
	})

	Convey("Manager.AddProcess 把已存在的 pid 登记进 Manager", t, func() {
		// 历史 bug：AddProcess 调 m.NewProcess("", nil, nil) 传 nil opt
		// 让 NewProcessOptions for-range 时 opt.Apply 在 nil interface 上 panic。
		// 修复：改成 m.NewProcess("")（variadic 空 slice 不迭代）。
		m := NewManager()
		pid := os.Getpid() // 当前测试进程一定存在
		So(func() { m.AddProcess(pid) }, ShouldNotPanic)

		// 端到端断言：登记成功后能查到
		So(m.Size(), ShouldEqual, 1)
		So(m.GetProcess(pid), ShouldNotBeNil)
		So(m.GetProcess(pid).Process, ShouldNotBeNil)
		So(m.GetProcess(pid).Process.Pid, ShouldEqual, pid)
		So(m.Pids(), ShouldResemble, []int{pid})

		// 重复 AddProcess 同一 pid 不会重复登记
		m.AddProcess(pid)
		So(m.Size(), ShouldEqual, 1)

		// 清理
		m.RemoveProcess(pid)
		So(m.Size(), ShouldEqual, 0)
	})

	Convey("Manager.RemoveProcess + Clear（不通过 AddProcess）", t, func() {
		m := NewManager()
		// 直接通过 NewProcess + Manager 关联（绕开 AddProcess bug）
		// 用真启动一个 process 让 m.processes 有元素
		// 这里只验证 RemoveProcess / Clear 在空集合下安全
		m.RemoveProcess(99999) // 不存在的 pid
		So(m.Size(), ShouldEqual, 0)

		m.Clear() // 空集合 Clear
		So(m.Size(), ShouldEqual, 0)
	})

	Convey("Manager 空集合 SignalAll/KillAll 不报错", t, func() {
		m := NewManager()
		So(m.KillAll(), ShouldBeNil)
		So(m.SignalAll(syscall.SIGTERM), ShouldBeNil)
		// WaitAll 在空集合下立即返回
		m.WaitAll()
	})
}

func TestProcess_Run(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses /bin/echo")
	}
	Convey("ShellRun 跑简单 echo", t, func() {
		out, err := ShellRun("echo hello")
		So(err, ShouldBeNil)
		So(out, ShouldContainSubstring, "hello")
	})

	Convey("Run 跑可执行文件（/bin/echo）", t, func() {
		out, err := Run("/bin/echo", WithArgs("world"))
		So(err, ShouldBeNil)
		So(out, ShouldContainSubstring, "world")
	})

	Convey("Run 在不存在的可执行文件上返错误", t, func() {
		_, err := Run("/nonexistent/path/that/should/not/exist")
		So(err, ShouldNotBeNil)
	})
}

func TestProcess_PidBeforeStart(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses /bin/sh")
	}
	Convey("Pid 在 Start 前返 0", t, func() {
		p := NewProcess("/bin/sh")
		So(p.Pid(), ShouldEqual, 0)
	})
}

func TestProcess_StartTwice(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses /bin/echo")
	}
	Convey("Process.Start 二次调用返已有 pid", t, func() {
		// 跑一个超快退出的 echo
		p := NewProcess("/bin/echo", WithArgs("hi"))
		pid1, err := p.Start()
		So(err, ShouldBeNil)
		So(pid1, ShouldBeGreaterThan, 0)

		// Wait 让进程退出
		_ = p.Wait()

		// 再 Start：p.Process 不为 nil 走早返路径（返回原 pid）
		pid2, err := p.Start()
		So(err, ShouldBeNil)
		So(pid2, ShouldEqual, pid1)
	})
}

func TestNewProcessWithOptions_DupPath(t *testing.T) {
	Convey("Args[0] 与 path 相同时跳过重复添加", t, func() {
		path := "/bin/echo"
		// 故意把 path 也当作 Args[0]，看是否被跳过
		p := NewProcessWithOptions(path, NewProcessOptions(
			WithArgs(path, "extra"),
		))
		// 期望 Args = [path, "extra"]，不是 [path, path, "extra"]
		So(p.Args, ShouldResemble, []string{path, "extra"})
	})

	Convey("Args[0] 与 path 不同时全部追加", t, func() {
		p := NewProcessWithOptions("/bin/echo", NewProcessOptions(
			WithArgs("hello", "world"),
		))
		So(p.Args, ShouldResemble, []string{"/bin/echo", "hello", "world"})
	})
}

// 启动一个临时脚本 + 真 Kill / Signal / Release 路径
func TestProcess_KillSignalRelease(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses bash sleep")
	}
	Convey("Kill 停止运行中进程", t, func() {
		// 写一个 sleep 脚本
		tmp := filepath.Join(os.TempDir(), "xproc_kill_test.sh")
		err := os.WriteFile(tmp, []byte("#!/bin/sh\nsleep 30\n"), 0755)
		So(err, ShouldBeNil)
		defer os.Remove(tmp)

		m := NewManager()
		p := m.NewProcess(tmp)
		_, err = p.Start()
		So(err, ShouldBeNil)
		So(p.Pid(), ShouldBeGreaterThan, 0)

		// Kill
		So(p.Kill(), ShouldBeNil)
		// Manager 应自动移除
		So(m.GetProcess(p.Pid()), ShouldBeNil)
	})

	Convey("Signal 发送信号到运行中进程", t, func() {
		tmp := filepath.Join(os.TempDir(), "xproc_signal_test.sh")
		err := os.WriteFile(tmp, []byte("#!/bin/sh\nsleep 30\n"), 0755)
		So(err, ShouldBeNil)
		defer os.Remove(tmp)

		p := NewProcess(tmp)
		_, err = p.Start()
		So(err, ShouldBeNil)

		// 发 SIGTERM
		So(p.Signal(syscall.SIGTERM), ShouldBeNil)
		// 等进程退出
		_ = p.Wait()
	})

	Convey("Release 释放 process resource（在 Wait 之后）", t, func() {
		// 用 echo 让进程立即退出
		p := NewProcess("/bin/echo", WithArgs("done"))
		_, err := p.Start()
		So(err, ShouldBeNil)
		_ = p.Wait()
		// Release 在 Wait 之后调用：unix 上通常返 nil 或"already released"
		// 我们只验证不 panic
		So(func() { _ = p.Release() }, ShouldNotPanic)
	})
}

// TestManager_WithRunningProcesses 让 Manager 真有进程，覆盖 Processes /
// Pids / WaitAll / Size 的循环分支。
func TestManager_WithRunningProcesses(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("uses /bin/echo")
	}
	Convey("Manager 持有 1 个 echo 进程，遍历 API 全跑", t, func() {
		m := NewManager()
		p := m.NewProcess("/bin/echo", WithArgs("hi"))
		_, err := p.Start()
		So(err, ShouldBeNil)

		So(m.Size(), ShouldEqual, 1)
		So(m.Pids(), ShouldHaveLength, 1)
		So(m.Processes(), ShouldHaveLength, 1)

		m.WaitAll() // 走 len>0 分支
		So(m.GetProcess(p.Pid()), ShouldNotBeNil)
	})

	Convey("Manager.SignalAll / KillAll 在 has-process 上工作", t, func() {
		tmp := filepath.Join(os.TempDir(), "xproc_mgr_kill.sh")
		err := os.WriteFile(tmp, []byte("#!/bin/sh\nsleep 30\n"), 0755)
		So(err, ShouldBeNil)
		defer os.Remove(tmp)

		m := NewManager()
		p1 := m.NewProcess(tmp)
		p2 := m.NewProcess(tmp)
		_, err = p1.Start()
		So(err, ShouldBeNil)
		_, err = p2.Start()
		So(err, ShouldBeNil)
		So(m.Size(), ShouldEqual, 2)

		// SignalAll SIGTERM 走 for 循环
		_ = m.SignalAll(syscall.SIGTERM)
		// 等进程退出
		_ = p1.Wait()
		_ = p2.Wait()

		// 用第三个 process 测 KillAll
		p3 := m.NewProcess(tmp)
		_, err = p3.Start()
		So(err, ShouldBeNil)
		// KillAll 让 m 中所有 process 都被 Kill（含已 wait 的，可能返 err）
		_ = m.KillAll()

		// Clear 走 has-process 分支（KillAll 后 Manager 已 Delete 大部分）
		m.Clear()
		So(m.Size(), ShouldEqual, 0)
	})
}
