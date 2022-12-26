package gzip

import (
	"compress/gzip"
	"io"
	"io/ioutil"
	"sync"

	c2 "github.com/sandwich-go/boost/encoding2/compressor"
	"github.com/sandwich-go/boost/xerror"
)

// Name is the name registered for the gzip compressor.
const Name = "gzip"

func init() {
	c2.RegisterCompressor(createCompressor())
}

func createCompressor() *compressor {
	c := &compressor{}
	c.poolCompressor.New = func() interface{} {
		return &writer{Writer: gzip.NewWriter(ioutil.Discard), pool: &c.poolCompressor}
	}
	return c
}

type writer struct {
	*gzip.Writer
	pool *sync.Pool
}

// SetLevel updates the registered gzip compressor to use the compression level specified (gzip.HuffmanOnly is not supported).
// NOTE: this function must only be called during initialization time (i.e. in an init() function),
// and is not thread-safe.
//
// The error returned will be nil if the specified level is valid.
func SetLevel(level int) error {
	if level < gzip.DefaultCompression || level > gzip.BestCompression {
		return xerror.NewText("invalid gzip compression level: %d", level)
	}
	c := c2.GetCompressor(Name).(*compressor)
	c.poolCompressor.New = func() interface{} {
		w, err := gzip.NewWriterLevel(ioutil.Discard, level)
		if err != nil {
			panic(err)
		}
		return &writer{Writer: w, pool: &c.poolCompressor}
	}
	return nil
}

func (c *compressor) Compress(w io.Writer) (c2.WriteFlushCloser, error) {
	z := c.poolCompressor.Get().(*writer)
	z.Writer.Reset(w)
	return z, nil
}

func (z *writer) Close() error {
	err := z.Writer.Close()
	z.pool.Put(z)
	return err
}

type reader struct {
	*gzip.Reader
	pool *sync.Pool
}

func (c *compressor) Decompress(r io.Reader) (io.ReadCloser, error) {
	z, inPool := c.poolDecompressor.Get().(*reader)
	if !inPool {
		newZ, err := gzip.NewReader(r)
		if err != nil {
			return nil, err
		}
		return &reader{Reader: newZ, pool: &c.poolDecompressor}, nil
	}
	if err := z.Reset(r); err != nil {
		c.poolDecompressor.Put(z)
		return nil, err
	}
	return z, nil
}

func (z *reader) Close() error {
	err := z.Reader.Close()
	z.pool.Put(z)
	return err
}

func (c *compressor) Name() string {
	return Name
}

type compressor struct {
	poolCompressor   sync.Pool
	poolDecompressor sync.Pool
}
