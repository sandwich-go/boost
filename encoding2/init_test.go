package encoding2

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

var fileData []byte

func BenchmarkCompressGzip(t *testing.B) {
	compressor := GetCompressor("gzip")
	var buf bytes.Buffer
	c, _ := compressor.Compress(&buf)
	c.Write(fileData)
	c.Close()
}

func BenchmarkCompressLz4(t *testing.B) {
	compressor := GetCompressor("lz4")
	var buf bytes.Buffer
	c, _ := compressor.Compress(&buf)
	c.Write(fileData)
	c.Close()
}

func TestGetCompressor(t *testing.T) {
	GetCompressor("lz4")
}

func TestMain(m *testing.M) {
	fileData, _ = ioutil.ReadFile("./pb_uncompress.data")
	os.Exit(m.Run())
}
