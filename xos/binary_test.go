package xos

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBinary(t *testing.T) {
	Convey("search binary file", t, func() {
		file := "binary.go"
		So(SearchBinary(file), ShouldEqual, file)
		So(SearchBinaryPath(file), ShouldBeEmpty)
	})
}
