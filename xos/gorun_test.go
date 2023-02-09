package xos

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGoRun(t *testing.T) {
	Convey("go run", t, func() {
		So(MustGetBinaryFilePath(), ShouldContainSubstring, ".test")
		So(IsGoRun(), ShouldBeTrue)
	})
}
