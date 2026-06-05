package xerror

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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

	// NewArraySafe 工厂构造的 Array.goroutineSafe = true，所有方法走加锁分支。
	// 之前的测试只覆盖默认 goroutineSafe=false 路径，这里专门测 true 路径让
	// 各方法的两条分支都覆盖到。
	Convey("NewArraySafe goroutineSafe path", t, func() {
		arr := NewArraySafe()
		So(arr, ShouldNotBeNil)

		// 走加锁分支的所有方法
		So(arr.Err(), ShouldBeNil)
		So(arr.LastErr(), ShouldBeNil)
		So(arr.WrappedErrors(), ShouldBeEmpty)
		// nil push 在 goroutineSafe 下也能被吞掉
		arr.Push(nil)
		So(arr.Err(), ShouldBeNil)

		e1 := errors.New("safe-1")
		arr.Push(e1)
		So(arr.LastErr(), ShouldEqual, e1)
		So(arr.Err(), ShouldNotBeNil)
		So(arr.Is(e1), ShouldBeTrue)
		So(arr.Is(errors.New("nope")), ShouldBeFalse)
		So(arr.Error(), ShouldContainSubstring, "safe-1")
		So(arr.String(), ShouldNotBeEmpty)

		// SetFormatFunc 走加锁写分支
		arr.SetFormatFunc(DotFormatFunc)
		arr.Push(errors.New("safe-2"))
		So(arr.Error(), ShouldContainSubstring, "safe-1,safe-2")
	})

	// 错误格式化函数 DotFormatFunc 默认行为：用 ',' 拼接 error.Error()
	Convey("DotFormatFunc / ListFormatFunc", t, func() {
		es := []error{errors.New("a"), errors.New("b"), errors.New("c")}
		So(DotFormatFunc(es), ShouldEqual, "a,b,c")
		out := ListFormatFunc(es)
		So(out, ShouldContainSubstring, "3 errors occurred")
		So(out, ShouldContainSubstring, "#1: a")
		So(out, ShouldContainSubstring, "#2: b")
		So(out, ShouldContainSubstring, "#3: c")
	})
}
