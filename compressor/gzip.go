package compressor

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/sandwich-go/boost"
	"io/ioutil"
	"sync"
)

type gzipCompressor struct {
	level                        int
	spWriter, spReader, spBuffer sync.Pool
}

func newGzipCompressor(level int) (Compressor, error) {
	c := &gzipCompressor{level: level}
	err := c.init()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *gzipCompressor) init() error {
	if c.level < HuffmanOnly || c.level > BestCompression {
		return fmt.Errorf("gzip: invalid compression level: %d", c.level)
	}
	c.spWriter = sync.Pool{New: func() interface{} { w, _ := gzip.NewWriterLevel(nil, c.level); return w }}
	c.spReader = sync.Pool{New: func() interface{} { return new(gzip.Reader) }}
	c.spBuffer = sync.Pool{New: func() interface{} { return bytes.NewBuffer(nil) }}
	return nil
}

func (c *gzipCompressor) Flat(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}
	buf := c.spBuffer.Get().(*bytes.Buffer)
	w := c.spWriter.Get().(*gzip.Writer)
	w.Reset(buf)

	defer func() {
		buf.Reset()
		c.spBuffer.Put(buf)
		boost.LogErrorAndEatError(w.Close())
		c.spWriter.Put(w)
	}()
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Flush()
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *gzipCompressor) Inflate(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}
	buf := c.spBuffer.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		c.spBuffer.Put(buf)
	}()

	_, err := buf.Write(data)
	if err != nil {
		return nil, err
	}

	gr := c.spReader.Get().(*gzip.Reader)
	defer func() {
		c.spReader.Put(gr)
	}()
	err = gr.Reset(buf)
	if err != nil {
		return nil, err
	}
	defer func() {
		boost.LogErrorAndEatError(gr.Close())
	}()
	return ioutil.ReadAll(gr)
}
