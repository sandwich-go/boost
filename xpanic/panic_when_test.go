package xpanic

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPanicWhen(t *testing.T) {
	Convey("panic when", t, func() {
		var err = errors.New("error")

		// WhenErrorAsFmtFirst 用 fmt.Sprintf 把 err 拼成 string 后 panic（commit
		// f7dd56a 把原 fmt.Errorf 改为 fmt.Sprintf，让 panic value 是 string 而
		// 非 *fmt.wrapError）；同时 %w 在 Sprintf 中无效，应改用 %v —— 这里测
		// 试期望按当前实现的 fmt.Sprintf("%v, %d", err, 1) 输出来断言。
		// 注：调用方仍可写 "%w" 但 Sprintf 会把它当无效 verb 输出 "%!w(...)"，
		// 所以测试用 %v 对齐实现的"显式建议"形式。
		So(func() {
			WhenErrorAsFmtFirst(err, "%v, %d", 1)
		}, ShouldPanic)
		So(func() {
			WhenErrorAsFmtFirst(nil, "%v, %d", 1)
		}, ShouldNotPanic)

		Try(func() {
			WhenErrorAsFmtFirst(err, "%v, %d", 1)
		}).Catch(func(e E) {
			// catch 收到的 e 是 xpanic.exception 结构体（commit d885113 设计）；
			// panicValueFrom 取出原 panic value（这里是 fmt.Sprintf 出来的 string）。
			So(panicValueFrom(e), ShouldEqual, "error, 1")
		})

		So(func() { WhenError(err) }, ShouldPanic)
		So(func() { WhenError(nil) }, ShouldNotPanic)
		// WhenError reason 非空分支：panic value 是 reason join
		Try(func() {
			WhenError(err, "step1", "step2")
		}).Catch(func(e E) {
			So(panicValueFrom(e), ShouldEqual, "step1\nstep2")
		})

		So(func() { WhenTrue(true, "%d", 1) }, ShouldPanic)
		So(func() { WhenTrue(false, "%d", 1) }, ShouldNotPanic)
		So(func() { WhenFalse(false, "%d", 1) }, ShouldPanic)
		So(func() { WhenFalse(true, "%d", 1) }, ShouldNotPanic)

		// WhenHereNotNil 走 fmt.Sprintf("err should be nil when here, got:%v", err)
		// （历史 bug：曾用 %w，本次随包修复改为 %v）。验证 panic value 字面格式。
		Try(func() {
			WhenHereNotNil(err)
		}).Catch(func(e E) {
			So(panicValueFrom(e), ShouldEqual, "err should be nil when here, got:error")
		})
		So(func() { WhenHereNotNil(nil) }, ShouldNotPanic)
	})

	Convey("WhenNil / WhenNotNil", t, func() {
		// 区分 typed-nil 与 untyped-nil（isnil.Check 处理 typed-nil）
		var typedNilPtr *int
		var nonNilPtr = new(int)

		// nil → WhenNil 触发；WhenNotNil 不触发
		So(func() { WhenNil(nil, "should panic: %s", "untyped nil") }, ShouldPanic)
		So(func() { WhenNotNil(nil, "should not") }, ShouldNotPanic)

		// typed nil（*int(nil)）—— isnil.Check 也视作 nil
		So(func() { WhenNil(typedNilPtr, "should panic on typed nil") }, ShouldPanic)
		So(func() { WhenNotNil(typedNilPtr, "should not") }, ShouldNotPanic)

		// non-nil → WhenNil 不触发；WhenNotNil 触发
		So(func() { WhenNil(nonNilPtr, "should not") }, ShouldNotPanic)
		So(func() { WhenNotNil(nonNilPtr, "should panic on non-nil: %d", 42) }, ShouldPanic)

		// panic value 是 fmt.Sprintf 后的 string（与 WhenTrue 一致）
		Try(func() {
			WhenNotNil(nonNilPtr, "got %d", 42)
		}).Catch(func(e E) {
			So(panicValueFrom(e), ShouldEqual, "got 42")
		})
	})
}
