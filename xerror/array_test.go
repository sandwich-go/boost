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

		err := New(WithText("io error"), WithCode(500), WithStack())
		errW := Wrap(err, "link error")
		errW = Wrap(errW, "session error")

		t.Log(errW.Error())
		t.Log(Caller(err.Cause(), 0))
	})
}
