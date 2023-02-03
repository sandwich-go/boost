package xip

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFree(t *testing.T) {
	Convey("get free port", t, func() {
		port, err := GetFreePort()
		So(err, ShouldBeNil)
		So(port, ShouldBeGreaterThan, 0)
		t.Log("free port:", port)
	})
}
