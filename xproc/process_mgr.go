package xproc

import (
	"os"
	"sync"
)

type Manager struct {
	processes sync.Map
}

func NewManager() *Manager { return &Manager{} }

func (m *Manager) NewProcess(path string, opt ...ProcessOption) *Process {
	p := NewProcess(path, opt...)
	p.Manager = m
	return p
}

func (m *Manager) GetProcess(pid int) *Process {
	if v, ok := m.processes.Load(pid); ok {
		return v.(*Process)
	}
	return nil
}

// AddProcess 通过 pid 把已存在的 OS 进程登记到 Manager。
// 仅做"接管 / 登记"，不会启动新进程；若 pid 已在 Manager 中则忽略。
//
// 注意：Manager.NewProcess 的 opt 是 variadic，不传等于空 slice，
// NewProcessOptions for-range 不会迭代 nil opt（详见 §历史 bug：传
// nil ProcessOption 会让 opt.Apply 在 nil interface 上 panic）。
func (m *Manager) AddProcess(pid int) {
	if _, ok := m.processes.Load(pid); !ok {
		if process, err := os.FindProcess(pid); err == nil {
			p := m.NewProcess("")
			p.Process = process
			m.processes.Store(pid, p)
		}
	}
}

func (m *Manager) RemoveProcess(pid int) {
	m.processes.Delete(pid)
}

func (m *Manager) Processes() []*Process {
	processes := make([]*Process, 0)
	m.processes.Range(func(key, value interface{}) bool {
		processes = append(processes, value.(*Process))
		return true
	})
	return processes
}

func (m *Manager) Pids() (ret []int) {
	m.processes.Range(func(key, value interface{}) bool {
		ret = append(ret, key.(int))
		return true
	})
	return
}

func (m *Manager) WaitAll() {
	processes := m.Processes()
	if len(processes) > 0 {
		for _, p := range processes {
			// WaitAll 语义是"等所有进程结束"，单个 Wait 错误不阻断其他进程
			// 等待，显式忽略。调用方需要单进程错误的话用 Manager.Processes
			// 拿到 *Process 后自行 Wait。
			_ = p.Wait()
		}
	}
}

func (m *Manager) KillAll() error {
	for _, p := range m.Processes() {
		if err := p.Kill(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Manager) SignalAll(sig os.Signal) error {
	for _, p := range m.Processes() {
		if err := p.Signal(sig); err != nil {
			return err
		}
	}
	return nil
}
func (m *Manager) Clear() {
	m.processes.Range(func(key, value interface{}) bool {
		m.processes.Delete(key)
		return true
	})
}

func (m *Manager) Size() (c int) {
	m.processes.Range(func(key, value interface{}) bool {
		c++
		return true
	})
	return
}
