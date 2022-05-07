package xproc

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/sandwich-go/boost/xos"
	"github.com/sandwich-go/boost/xslice"
)

const envKeyPPid = "GO_PROC_PPID"

type Process struct {
	cc *ProcessOptions
	exec.Cmd
	Manager *Manager
	PPid    int
}

// NewProcess creates and returns a new Process.
func NewProcessWithOptions(path string, cc *ProcessOptions) *Process {
	process := &Process{
		cc:      cc,
		Manager: nil,
		PPid:    os.Getpid(),
		Cmd: exec.Cmd{
			Args:       []string{path},
			Path:       path,
			Stdin:      cc.Stdin,
			Stdout:     cc.Stdout,
			Stderr:     cc.Stderr,
			Env:        append(os.Environ(), cc.Env...),
			Dir:        cc.WorkingDir,
			ExtraFiles: make([]*os.File, 0), // not support on windows
		},
	}
	if len(cc.Args) > 0 {
		start := 0
		if strings.EqualFold(path, cc.Args[0]) {
			start = 1
		}
		process.Args = append(process.Args, cc.Args[start:]...)
	}
	return process
}

// NewProcess creates and returns a new Process.
func NewProcess(path string, opt ...ProcessOption) *Process {
	return NewProcessWithOptions(path, NewProcessOptions(opt...))
}

// NewProcessShellCmdWithOptions creates and returns a process with given command and optional environment variable array.
func NewProcessShellCmdWithOptions(cmd string, cc *ProcessOptions) *Process {
	argsLen := len(cc.Args)
	cc.Args = xslice.StringSetAdd(parseCommand(cmd), cc.Args...)
	if argsLen == 0 {
		cc.Args = xslice.StringSetAdd([]string{xos.GetShellOption()}, cc.Args...)
	}
	fmt.Println(cc.Args)
	return NewProcessWithOptions(xos.GetShell(), cc)
}

// Start starts executing the process in non-blocking way.
// It returns the pid if success, or else it returns an error.
func (p *Process) Start() (int, error) {
	if p.Process != nil {
		return p.Pid(), nil
	}
	p.Env = append(p.Env, fmt.Sprintf("%s=%d", envKeyPPid, p.PPid))
	if err := p.Cmd.Start(); err == nil {
		if p.Manager != nil {
			p.Manager.processes.Store(p.Process.Pid, p)
		}
		return p.Process.Pid, nil
	} else {
		return 0, err
	}
}

// Run executes the process in blocking way.
func (p *Process) Run() error {
	if _, err := p.Start(); err == nil {
		return p.Wait()
	} else {
		return err
	}
}

func (p *Process) Pid() int {
	if p.Process != nil {
		return p.Process.Pid
	}
	return 0
}

func (p *Process) Release() error {
	return p.Process.Release()
}

// Kill causes the Process to exit immediately.
func (p *Process) Kill() (err error) {
	if err = p.Process.Kill(); err != nil {
		return err
	}
	if p.Manager != nil {
		p.Manager.processes.Delete(p.Pid())
	}
	if runtime.GOOS != "windows" {
		if err = p.Process.Release(); err != nil {
			return err
		}
	}
	_, err = p.Process.Wait()
	return err
}

// Signal sends a signal to the Process.
// Sending Interrupt on Windows is not implemented.
func (p *Process) Signal(sig os.Signal) error {
	return p.Process.Signal(sig)
}
