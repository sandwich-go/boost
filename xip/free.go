package xip

import (
	"fmt"
	"github.com/sandwich-go/boost/internal/log"
	"net"
)

// GetFreePort asks the kernel for a free open port that is ready to use.
func GetFreePort() (int, error) {
	addr, err0 := net.ResolveTCPAddr("tcp", "localhost:0")
	if err0 != nil {
		return 0, err0
	}
	l, err1 := net.ListenTCP("tcp", addr)
	if err1 != nil {
		return 0, err1
	}
	freePort := l.Addr().(*net.TCPAddr).Port
	if err2 := l.Close(); err2 != nil {
		log.Error(fmt.Sprintf("close free port: %d, err: %s", freePort, err2.Error()))
	}
	return freePort, nil
}
