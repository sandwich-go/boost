// ringbuf 包实现了一个环形缓冲区，可以对ringbuf进行读写操作
// New 返回一个指定缓冲区大小的环形缓冲区实例
package ringbuf

import (
	"errors"
	"fmt"
	"github.com/sandwich-go/boost/internal/log"
)

// NotEnoughMem 内存不够
var NotEnoughMem = errors.New("not enough memory")

type Ringbuf struct {
	buf         []byte
	start, size int
}

// New 这个 ringBuf 不是线程安全的！！！
func New(size int) *Ringbuf {
	return &Ringbuf{make([]byte, size), 0, 0}
}

// Write 写入数据
func (r *Ringbuf) Write(b []byte) error {
	if cap(r.buf)-r.size < len(b) {
		return NotEnoughMem
	}
	for len(b) > 0 {
		start := (r.start + r.size) % len(r.buf)
		n := copy(r.buf[start:], b)
		b = b[n:]

		r.size += n
	}
	return nil
}

// Read 读数据
func (r *Ringbuf) Read(b []byte, length int) int {
	read := 0
	for length > 0 && r.size > 0 {
		end := r.start + length
		if end > len(r.buf) {
			end = len(r.buf)
		}
		n := copy(b, r.buf[r.start:end])
		read += n
		length -= n
		b = b[n:]

		r.size -= n
		r.start = (r.start + n) % len(r.buf)
	}
	return read
}

// PreUse 预申请缓存
func (r *Ringbuf) PreUse(length int) []byte {
	s := (r.start + r.size) % r.Capacity()
	e := s + length
	if e > r.Capacity() {
		e = r.Capacity()
	}
	if s < r.start && e > r.start {
		e = r.start
	}
	if s == e {
		log.Warn(fmt.Sprintf("not enough, length: %d,s: %d,e: %d,start: %d,size: %d", length, s, e, r.start, r.size))
	}
	return r.buf[s:e]
}

// RealUse 标记实际使用
func (r *Ringbuf) RealUse(length int) {
	if length <= 0 {
		return
	}
	r.size += length
}

// Size buf 大小
func (r *Ringbuf) Size() int { return r.size }

// Capacity buf 容量
func (r *Ringbuf) Capacity() int { return cap(r.buf) }
