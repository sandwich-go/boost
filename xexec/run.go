package xexec

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

// Run 运行指定命令，in作为std in输入存在
func Run(arg, dir string, in ...*bytes.Buffer) (string, error) {
	shellPath := getShell()
	cmd := exec.Command(shellPath, "-c", arg)
	if len(dir) > 0 {
		cmd.Dir = dir
	}
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	if len(in) > 0 {
		cmd.Stdin = in[0]
	}
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err != nil {
		if stderr.Len() > 0 {
			return "", errors.New(strings.TrimSuffix(stderr.String(), "\n"))
		}
		return "", err
	}

	return strings.TrimSuffix(stdout.String(), "\n"), nil
}
