package xexec

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

func (m *Manager) AddProcess(pid int) {
	if _, ok := m.processes.Load(pid); !ok {
		if process, err := os.FindProcess(pid); err == nil {
			p := m.NewProcess("", nil, nil)
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
			p.Wait()
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
