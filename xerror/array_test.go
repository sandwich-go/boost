package xerror

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestArray(t *testing.T) {
	Convey("array", t, func() {
		var arr Array
		So(arr.Err(), ShouldBeNil)
		So(arr.LastErr(), ShouldBeNil)

		var e1 = errors.New("error 1")
		var e2 = errors.New("error 2")
		arr.Push(e1)
		So(arr.Err(), ShouldNotBeNil)
		So(arr.LastErr(), ShouldNotBeNil)
		So(arr.LastErr(), ShouldEqual, e1)

		arr.Push(e2)
		So(arr.Err(), ShouldNotBeNil)
		So(arr.LastErr(), ShouldNotBeNil)
		So(arr.LastErr(), ShouldEqual, e2)

		So(arr.WrappedErrors(), ShouldResemble, []error{e1, e2})

		t.Log(arr.String())
		t.Log(arr.Error())

		arr.SetFormatFunc(DotFormatFunc)
		t.Log(arr.Error())
	})
	Convey("errors.Is", t, func() {
		errArray := &Array{}
		errArray.Push(NewText("1"))
		e1 := NewText("2")
		e2 := Wrap(e1, "wrap with xerror")
		So(errors.Is(errArray, e1), ShouldBeFalse)
		errArray.Push(e2)
		So(errors.Is(errArray, e1), ShouldBeTrue)
		So(errors.Is(errArray, e2), ShouldBeTrue)
		So(errors.Is(errArray.Err(), e1), ShouldBeTrue)
	})
}
