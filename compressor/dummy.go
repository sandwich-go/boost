package compressor

type dummyCompressor struct{}

func newDummyCompressor() (Compressor, error) {
	return &dummyCompressor{}, nil
}

func (*dummyCompressor) Flat(data []byte) ([]byte, error)    { return data, nil }
func (*dummyCompressor) Inflate(data []byte) ([]byte, error) { return data, nil }
