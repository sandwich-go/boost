package xos

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFunc(t *testing.T) {
	Convey("func name", t, func() {
		s, err := FuncFullNameUsingReflect(FuncFullNameUsingReflect)
		So(err, ShouldBeNil)
		s, err = FuncBaseNameUsingReflect(FuncFullNameUsingReflect)
		So(err, ShouldBeNil)
		So(s, ShouldEqual, "FuncFullNameUsingReflect")
		s, err = FuncBaseNameUsingReflect(FuncBaseNameUsingReflect)
		So(err, ShouldBeNil)
		So(s, ShouldEqual, "FuncBaseNameUsingReflect")
	})
}
