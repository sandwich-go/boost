package xerror_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/sandwich-go/boost/xerror"

	. "github.com/smartystreets/goconvey/convey"
)

// TestSafeArray_BasicAPI 覆盖 SafeArray 8 个方法的非并发路径（与 Array
// 行为对齐）。SafeArray 内部包一层 sync.RWMutex 保护一个嵌套 Array，所以
// 行为应等价。
func TestSafeArray_BasicAPI(t *testing.T) {
	Convey("SafeArray basic API", t, func() {
		var arr xerror.SafeArray
		So(arr.Err(), ShouldBeNil)
		So(arr.LastErr(), ShouldBeNil)
		So(arr.WrappedErrors(), ShouldBeEmpty)

		// nil 被 Push 丢弃
		arr.Push(nil)
		So(arr.Err(), ShouldBeNil)
		So(arr.LastErr(), ShouldBeNil)

		e1 := errors.New("error 1")
		e2 := errors.New("error 2")
		arr.Push(e1)
		So(arr.Err(), ShouldNotBeNil)
		So(arr.LastErr(), ShouldEqual, e1)

		arr.Push(e2)
		So(arr.LastErr(), ShouldEqual, e2)
		So(arr.WrappedErrors(), ShouldResemble, []error{e1, e2})

		// String / Error 都能跑出文本，具体格式由 ListFormatFunc 决定（默认）
		So(arr.String(), ShouldContainSubstring, "errors")
		So(arr.Error(), ShouldContainSubstring, "error 1")
		So(arr.Error(), ShouldContainSubstring, "error 2")

		// errors.Is 透视嵌套 *Error
		xerr := xerror.NewText("nested")
		wrapped := xerror.Wrap(xerr, "wrapped")
		arr.Push(wrapped)
		So(arr.Is(xerr), ShouldBeTrue)
		So(arr.Is(errors.New("nope")), ShouldBeFalse)

		// 自定义 format func 后 Error 输出走新格式
		arr.SetFormatFunc(xerror.DotFormatFunc)
		So(arr.Error(), ShouldContainSubstring, "error 1,error 2,")
	})
}

// TestSafeArray_Concurrent 真正测 SafeArray 的并发安全（这是 SafeArray 与
// Array 的唯一差别，也是 §4.4 race 硬约束的覆盖点）。
//
// 1000 个 goroutine 并发 Push + Read，最终条数应一致；race detector 不报。
// run 时配合 -race，CI 上 race job 会跑这条。
func TestSafeArray_Concurrent(t *testing.T) {
	const goroutines = 100
	const perGoroutine = 50

	var arr xerror.SafeArray
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for g := 0; g < goroutines; g++ {
		go func(gid int) {
			defer wg.Done()
			for i := 0; i < perGoroutine; i++ {
				arr.Push(errors.New("oops"))
				// 同时读，触发 RLock/Lock 竞争
				_ = arr.LastErr()
				_ = arr.Err()
			}
		}(g)
	}
	wg.Wait()

	if got, want := len(arr.WrappedErrors()), goroutines*perGoroutine; got != want {
		t.Fatalf("SafeArray.Push lost data under concurrency: got=%d want=%d", got, want)
	}

	// Is/Error/String 在并发完成后也应稳定可调
	if !arr.Is(errors.New("oops-no-match")) && arr.Err() == nil {
		t.Error("Err() unexpectedly nil after concurrent Push")
	}
}
