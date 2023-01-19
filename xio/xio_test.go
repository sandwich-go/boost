package xio

import (
	"bytes"
	"context"
	"io"
	"testing"
)

func TestReader(t *testing.T) {
	buf := []byte("abcdef")
	buf2 := make([]byte, 3)
	r := NewReader(context.Background(), bytes.NewReader(buf))

	// read first half
	n, err := r.Read(buf2)
	if n != 3 {
		t.Error("n should be 3")
	}
	if err != nil {
		t.Error("should have no error")
	}
	if string(buf2) != string(buf[:3]) {
		t.Error("incorrect contents")
	}

	// read second half
	n, err = r.Read(buf2)
	if n != 3 {
		t.Error("n should be 3")
	}
	if err != nil {
		t.Error("should have no error")
	}
	if string(buf2) != string(buf[3:6]) {
		t.Error("incorrect contents")
	}

	// read more.
	n, err = r.Read(buf2)
	if n != 0 {
		t.Error("n should be 0", n)
	}
	if err != io.EOF {
		t.Error("should be EOF", err)
	}
}
