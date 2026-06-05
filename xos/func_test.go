package xos

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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
