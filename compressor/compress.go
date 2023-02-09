package compressor

type compress struct {
	opts *Options
	Compressor
}

// New 初始化 Compressor
// 使用 GZIP 方式时，解压缩等级才有效，其他模式该参数无效
// 默认方式，或者使用非 GZIP / Snappy 方式，均为 Dummy 方式
func New(opts ...Option) (Compressor, error) {
	var err error
	c := &compress{opts: NewOptions(opts...)}
	switch c.opts.GetType() {
	case GZIP:
		c.Compressor, err = newGzipCompressor(c.opts.GetLevel())
	case Snappy:
		c.Compressor, err = newSnappyCompressor()
	default:
		c.Compressor, err = newDummyCompressor()
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

// MustNew 初始化 Compressor，否则 panic
// 当使用 GZIP 方式时，必须传入正确的解压缩等级
func MustNew(opts ...Option) Compressor {
	c, err := New(opts...)
	if err != nil {
		panic(err)
	}
	return c
}
