package xio

import (
	"context"
	"io"
	"time"
)

type ioRet struct {
	n   int
	err error
}

type ctxReader struct {
	r   io.Reader
	ctx context.Context
}

// NewReader 返回context敏感的io reader
func NewReader(ctx context.Context, r io.Reader) io.Reader {
	return &ctxReader{ctx: ctx, r: r}
}

// block for debug
var block = false

// Read 由于Read是block操作，内部为每一次Read启动了独立协程协助读取,如果超时，则会返回(0，ctx.Error)
func (r *ctxReader) Read(buf []byte) (int, error) {
	bufReading := make([]byte, len(buf))

	c := make(chan ioRet, 1)

	go func() {
		if block {
			time.Sleep(10 * time.Second)
		}
		n, err := r.r.Read(bufReading)
		c <- ioRet{n: n, err: err}
		close(c)
	}()

	select {
	case ret := <-c:
		copy(buf, bufReading)
		return ret.n, ret.err
	case <-r.ctx.Done():
		return 0, r.ctx.Err()
	}
}
