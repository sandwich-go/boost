package xos

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestHide(t *testing.T) {
	Convey("hide", t, func() {
		var needCreate bool
		var tmpFile = "/tmp/test_dir"
		t.Cleanup(func() {
			if needCreate {
				_ = os.RemoveAll(tmpFile)
			}
		})
		_, err := os.Stat(tmpFile)
		if err != nil {
			if os.IsNotExist(err) {
				_, err = IsHidden(tmpFile)
				So(err, ShouldNotBeNil)

				err = nil
				needCreate = true
				err = os.MkdirAll(tmpFile, 0775)
			}
		}
		So(err, ShouldBeNil)
		is, err1 := IsHidden(tmpFile)
		So(err1, ShouldBeNil)
		So(is, ShouldBeFalse)

		var newTmpFile string
		newTmpFile, err = Hide(tmpFile)
		So(err, ShouldBeNil)
		is, err = IsHidden(newTmpFile)
		So(err, ShouldBeNil)
		So(is, ShouldBeTrue)
		tmpFile, err = UnHide(newTmpFile)
		So(err, ShouldBeNil)
		is, err = IsHidden(tmpFile)
		So(err, ShouldBeNil)
		So(is, ShouldBeFalse)
	})
}
