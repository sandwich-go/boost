package md5

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：补 md5 包错误路径：
//   - File 不存在路径返错误
//   - Buffer io.Reader 报错时返错误

func TestFile_NotExist(t *testing.T) {
	Convey("File 在文件不存在时返错误", t, func() {
		_, err := File("/nonexistent/path/that/should/not/exist.txt")
		So(err, ShouldNotBeNil)
	})
}

// errReader 故意 Read 时返错误
type errReader struct{}

func (errReader) Read(_ []byte) (int, error) {
	return 0, errors.New("simulated read error")
}

func TestBuffer_ReadError(t *testing.T) {
	Convey("Buffer 在 reader 报错时返错误", t, func() {
		_, err := Buffer(errReader{})
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "simulated read error")
	})
}
