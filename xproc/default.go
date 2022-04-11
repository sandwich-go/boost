package xproc

import (
	"bytes"
	"errors"
	"strings"
)

var DefaultManager = NewManager()

// Run 运行指定命令,接管std out std error
func Run(path string, opt ...ProcessOption) (string, error) {
	cc := NewProcessOptions(opt...)
	stdOut := new(bytes.Buffer)
	stdErr := new(bytes.Buffer)
	cc.Stdout = stdOut
	cc.Stderr = stdErr
	return run(NewProcessWithOptions(path, cc), stdOut, stdErr)
}

// ShellRun 启动shell运行指定命令
func ShellRun(cmd string, opt ...ProcessOption) (string, error) {
	cc := NewProcessOptions(opt...)
	stdOut := new(bytes.Buffer)
	stdErr := new(bytes.Buffer)
	cc.Stdout = stdOut
	cc.Stderr = stdErr
	return run(NewProcessShellCmdWithOptions(cmd, cc), stdOut, stdErr)
}

func run(process *Process, stdOut, stdErr *bytes.Buffer) (string, error) {
	err := process.Run()
	if err != nil {
		if stdErr.Len() > 0 {
			return "", errors.New(strings.TrimSuffix(stdErr.String(), "\n"))
		}
		return "", err
	}
	return strings.TrimSuffix(stdOut.String(), "\n"), nil
}
