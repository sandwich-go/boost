package xerror

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// 内部测试：覆盖未导出 protoEnum 接口路径 + containsCode + Error.Format
// 各 verb / Cause 嵌套路径 / *Error nil receiver 路径。
//
// 这些路径要么访问未导出标识符（protoEnum），要么涉及 *Error 自身的 Code/
// Cause/Format 方法的全部分支，xerror_test 外部包没法或不便构造，写在内
// 部测试更直接。

// fakeProtoEnum 实现 protoEnum 接口（fmt.Stringer + NumberInt32 +
// EnumDescriptor），让 NewProtoEnum / WrapProtoEnum 能拿到 code 和 text。
type fakeProtoEnum struct {
	name string
	num  int32
}

func (f fakeProtoEnum) String() string                  { return f.name }
func (f fakeProtoEnum) NumberInt32() int32              { return f.num }
func (f fakeProtoEnum) EnumDescriptor() ([]byte, []int) { return nil, nil }

func TestNewProtoEnum_WrapProtoEnum(t *testing.T) {
	Convey("NewProtoEnum 把 enum.NumberInt32 当 code, enum.String 当 text", t, func() {
		pe := fakeProtoEnum{name: "ERR_NOT_FOUND", num: 404}
		err := NewProtoEnum(pe)
		So(err, ShouldNotBeNil)
		So(err.Code(), ShouldEqual, int32(404))
		So(err.Error(), ShouldEqual, "ERR_NOT_FOUND")
	})

	Convey("WrapProtoEnum nil 输入透传 nil", t, func() {
		pe := fakeProtoEnum{name: "ERR_X", num: 1}
		So(WrapProtoEnum(pe, nil), ShouldBeNil)
	})

	Convey("WrapProtoEnum 包装非 nil err，Code 跟随 enum.NumberInt32", t, func() {
		pe := fakeProtoEnum{name: "ERR_INTERNAL", num: 500}
		base := errors.New("disk fault")
		wrapped := WrapProtoEnum(pe, base)
		So(wrapped, ShouldNotBeNil)
		// wrapped 是 *Error，Code() 走 enum
		xerr, ok := wrapped.(*Error)
		So(ok, ShouldBeTrue)
		So(xerr.Code(), ShouldEqual, int32(500))
		So(errors.Is(wrapped, base), ShouldBeTrue)
	})
}

func TestContainsCode(t *testing.T) {
	Convey("ContainsCode nil err 返回 false", t, func() {
		So(ContainsCode(nil, 100), ShouldBeFalse)
	})

	Convey("ContainsCode 单层 *Error", t, func() {
		err := NewCode(123, "boom")
		So(ContainsCode(err, 123), ShouldBeTrue)
		So(ContainsCode(err, 999), ShouldBeFalse)
	})

	Convey("ContainsCode 经 Wrap 链向下找", t, func() {
		inner := NewCode(456, "inner")
		outer := WrapCode(789, inner, "outer")
		So(ContainsCode(outer, 789), ShouldBeTrue)
		So(ContainsCode(outer, 456), ShouldBeTrue)
		So(ContainsCode(outer, 0), ShouldBeFalse)
	})

	Convey("ContainsCode 在普通 error 链上 Unwrap 失败时返回 false", t, func() {
		// fmt.Errorf("...%w") wrap 一个非 APICode 的普通 error，最终 Unwrap 走完
		// 没找到 code 返回 false
		base := errors.New("plain")
		wrapped := fmt.Errorf("wrap: %w", base)
		So(ContainsCode(wrapped, 1), ShouldBeFalse)
	})

	Convey("ContainsCode 处理 Unwrap() []error 多分支", t, func() {
		// errors.Join 实现 Unwrap() []error
		e1 := NewCode(11, "a")
		e2 := NewCode(22, "b")
		joined := errors.Join(e1, e2)
		So(ContainsCode(joined, 11), ShouldBeTrue)
		So(ContainsCode(joined, 22), ShouldBeTrue)
		So(ContainsCode(joined, 33), ShouldBeFalse)

		// Join 中含 nil 元素：containsCode 内部 continue 跳过
		e3 := NewCode(44, "c")
		joinedWithNil := errors.Join(nil, e3, nil)
		So(ContainsCode(joinedWithNil, 44), ShouldBeTrue)
	})
}

func TestErrorMethod_NilReceiver(t *testing.T) {
	Convey("*Error nil receiver 各方法返回零值不 panic", t, func() {
		var nilE *Error
		So(nilE.Error(), ShouldEqual, "")
		So(nilE.Code(), ShouldEqual, ErrorCodeOk)
		So(nilE.Cause(), ShouldBeNil)
	})

	Convey("Error.Error 仅 text 无 wrapped err", t, func() {
		err := NewText("only text")
		So(err.Error(), ShouldEqual, "only text")
	})

	Convey("Error.Error 仅 wrapped err 无 text", t, func() {
		base := errors.New("base")
		// 通过 New + WithErr 构造 text="" 但 err!=nil 的 *Error
		err := New(WithErr(base))
		So(err.Error(), ShouldEqual, "base")
	})

	Convey("Error.Error text 与 wrapped 同时存在用 ': ' 拼接", t, func() {
		err := Wrap(errors.New("inner"), "outer")
		So(err.Error(), ShouldEqual, "outer: inner")
	})
}

func TestErrorFormat_AllVerbs(t *testing.T) {
	Convey("Error.Format 全 verb 路径", t, func() {
		old := IsErrorWithStack
		IsErrorWithStack = true
		defer func() { IsErrorWithStack = old }()

		err := New(WithText("io error"), WithStack())
		errW := Wrap(err, "session")

		// %s = 全错误链
		s := fmt.Sprintf("%s", errW)
		So(s, ShouldContainSubstring, "io error")
		So(s, ShouldContainSubstring, "session")

		// %-s = 当前错误的 text（不递归）
		neg := fmt.Sprintf("%-s", errW.(*Error))
		So(neg, ShouldEqual, "session")

		// %+s = 错误链 + stack
		plus := fmt.Sprintf("%+s", errW.(*Error))
		So(plus, ShouldContainSubstring, "session")
		So(plus, ShouldContainSubstring, "\n") // 栈帧之间有换行

		// %v 与 %s 等价
		v := fmt.Sprintf("%v", errW)
		So(v, ShouldEqual, s)

	})

	Convey("Error.Format %-s 在 text 空时回退 Error()", t, func() {
		// 直接构造 text="" 但 err 非空
		base := errors.New("inner")
		e := New(WithErr(base)) // *Error，text=""
		got := fmt.Sprintf("%-s", e)
		So(got, ShouldEqual, "inner") // 走 cc.Error() 回退分支
	})
}

func TestCause_NestedXerror(t *testing.T) {
	Convey("Cause 沿 *Error 链向下找原 error", t, func() {
		base := errors.New("plain")
		l1 := Wrap(base, "l1")
		l2 := Wrap(l1, "l2")
		l3 := Wrap(l2, "l3")
		So(Cause(l3), ShouldEqual, base)
	})

	Convey("Cause 遇到非 *Error 但实现 apiCause 的 err 返回其 Cause", t, func() {
		// xerror.NewText 是 *Error，Wrap 它得到 *Error 也是；底层无 plain err 时
		// 链尾返回最初的 *Error 自己
		e1 := NewText("root")
		e2 := Wrap(e1, "wrapped")
		// Cause 应该返回 e1（最底层的 *Error）
		root := Cause(e2)
		So(root.Error(), ShouldContainSubstring, "root")
	})
}

func TestStack_ApiPath(t *testing.T) {
	Convey("Stack(nil) 返回空", t, func() {
		So(Stack(nil), ShouldEqual, "")
	})

	Convey("Stack 在普通 error 上回退到 err.Error()", t, func() {
		base := errors.New("plain")
		So(Stack(base), ShouldEqual, "plain")
	})

	Convey("Stack 在 *Error 上调 apiStack.Stack()", t, func() {
		old := IsErrorWithStack
		IsErrorWithStack = true
		defer func() { IsErrorWithStack = old }()

		err := NewText("with stack")
		out := Stack(err)
		So(out, ShouldContainSubstring, "with stack")
		// 栈帧体现：函数名 / 文件路径出现
		So(out, ShouldContainSubstring, "xerror")
	})
}

func TestCaller_ApiPath(t *testing.T) {
	Convey("Caller(nil) 返回零值", t, func() {
		f, fn, ln := Caller(nil, 0)
		So(f, ShouldEqual, "")
		So(fn, ShouldEqual, "")
		So(ln, ShouldEqual, 0)
	})

	Convey("Caller 在普通 error 上返回零值", t, func() {
		f, fn, ln := Caller(errors.New("plain"), 0)
		So(f, ShouldEqual, "")
		So(fn, ShouldEqual, "")
		So(ln, ShouldEqual, 0)
	})

	Convey("Caller 在 *Error 上调 apiCaller.Caller()", t, func() {
		old := IsErrorWithStack
		IsErrorWithStack = true
		defer func() { IsErrorWithStack = old }()

		err := NewText("with caller")
		f, fn, ln := Caller(err, 0)
		So(f, ShouldNotEqual, "")
		So(fn, ShouldNotEqual, "")
		So(ln, ShouldBeGreaterThan, 0)
	})
}

func TestLogic_ApiPath(t *testing.T) {
	Convey("Logic(nil) 返回 false", t, func() {
		So(Logic(nil), ShouldBeFalse)
	})

	Convey("Logic 普通 error 返回 false", t, func() {
		So(Logic(errors.New("plain")), ShouldBeFalse)
	})

	Convey("Logic 实现 apiLogic 接口的 *Error", t, func() {
		err := New(WithLogic(true))
		So(Logic(err), ShouldBeTrue)
		err2 := New(WithLogic(false))
		So(Logic(err2), ShouldBeFalse)
	})

	Convey("Logic 兼容 err2 IsLogicException 接口", t, func() {
		// 自定义实现 IsLogicException 但不实现 apiLogic
		fake := &fakeLogicErr{msg: "logic", logic: true}
		So(Logic(fake), ShouldBeTrue)
		fake2 := &fakeLogicErr{msg: "not-logic", logic: false}
		So(Logic(fake2), ShouldBeFalse)
	})
}

type fakeLogicErr struct {
	msg   string
	logic bool
}

func (f *fakeLogicErr) Error() string          { return f.msg }
func (f *fakeLogicErr) IsLogicException() bool { return f.logic }

// 防 -unused 报错：让 strings 在测试编译时被引用
var _ = strings.Builder{}

// TestNewWrap_FormatBranches 覆盖 NewText/NewCode/Wrap/WrapCode 的两条分支：
// len(args) == 0 直接用 format 当 text；len(args) > 0 走 fmt.Sprintf。
//
// 包级 bench 也跑两条分支（_Format 后缀），但 bench 不记入 coverage，所以
// 单测得显式覆盖。
func TestNewWrap_FormatBranches(t *testing.T) {
	Convey("NewText 无 args 直接用 format 字符串", t, func() {
		So(NewText("plain msg").Error(), ShouldEqual, "plain msg")
	})
	Convey("NewText 有 args 走 fmt.Sprintf", t, func() {
		So(NewText("hello %s, %d", "world", 42).Error(), ShouldEqual, "hello world, 42")
	})

	Convey("NewCode 无 args / 有 args 两条路径", t, func() {
		So(NewCode(100, "no args").Error(), ShouldEqual, "no args")
		So(NewCode(101, "code=%d", 999).Error(), ShouldEqual, "code=999")
	})

	Convey("Wrap 无 args / 有 args 两条路径", t, func() {
		base := errors.New("inner")
		So(Wrap(base, "wrap-no-args").Error(), ShouldEqual, "wrap-no-args: inner")
		So(Wrap(base, "wrap-%s", "fmt").Error(), ShouldEqual, "wrap-fmt: inner")
	})

	Convey("WrapCode nil 透传 nil", t, func() {
		So(WrapCode(500, nil, "no-op"), ShouldBeNil)
	})

	Convey("WrapCode 无 args / 有 args 两条路径", t, func() {
		base := errors.New("inner")
		w1 := WrapCode(500, base, "wrapped").(*Error)
		So(w1.Error(), ShouldEqual, "wrapped: inner")
		So(w1.Code(), ShouldEqual, int32(500))

		w2 := WrapCode(501, base, "wrapped-%d", 42).(*Error)
		So(w2.Error(), ShouldEqual, "wrapped-42: inner")
		So(w2.Code(), ShouldEqual, int32(501))
	})
}

// TestErrorStack_NilAndChain 覆盖 *Error.Caller / *Error.Stack 的 nil receiver
// 和混合 chain（中间塞个非 *Error）路径。
func TestErrorStack_NilAndChain(t *testing.T) {
	Convey("(*Error)(nil).Caller 返回零值", t, func() {
		var nilE *Error
		f, fn, ln := nilE.Caller(0)
		So(f, ShouldEqual, "")
		So(fn, ShouldEqual, "")
		So(ln, ShouldEqual, 0)
	})

	Convey("(*Error)(nil).Stack 返回空字符串", t, func() {
		var nilE *Error
		So(nilE.Stack(), ShouldEqual, "")
	})

	Convey("Caller(skip>=stack length) 走完循环不命中返回零值", t, func() {
		old := IsErrorWithStack
		IsErrorWithStack = true
		defer func() { IsErrorWithStack = old }()

		err := NewText("with stack")
		// skip=10000 远超栈深度
		f, fn, ln := err.Caller(10000)
		So(f, ShouldEqual, "")
		So(fn, ShouldEqual, "")
		So(ln, ShouldEqual, 0)
	})

	Convey("Stack 在非 *Error chain tail 处 fallback 输出 err.Error()", t, func() {
		old := IsErrorWithStack
		IsErrorWithStack = true
		defer func() { IsErrorWithStack = old }()

		// 链：*Error → *Error → plain error
		// Stack 走到最末非 *Error 时进入 line 47-50 分支：format 错误信息后 break
		base := errors.New("plain-tail")
		mid := Wrap(base, "mid")
		top := Wrap(mid, "top").(*Error)
		out := top.Stack()
		So(out, ShouldContainSubstring, "top")
		So(out, ShouldContainSubstring, "mid")
		So(out, ShouldContainSubstring, "plain-tail")
	})
}

// TestOptionsAndSetters 覆盖剩余的 option / setter：WithCode / WithSkip /
// WithTimeout / SetLogic / UnsetLogic。这些是 fluent API 的简单 setter，
// 测试目的是确保字段被正确写入（防止下游 sed 替换字段名时跑路）。
func TestOptionsAndSetters(t *testing.T) {
	Convey("WithCode 写入 code 字段", t, func() {
		e := New(WithCode(42))
		So(e.Code(), ShouldEqual, int32(42))
	})

	Convey("WithSkip 影响 callers 帧偏移（与 WithStack 协同）", t, func() {
		// IsErrorWithStack=true 时 New 用 e.skip 抓栈；WithSkip 改 skip
		old := IsErrorWithStack
		IsErrorWithStack = true
		defer func() { IsErrorWithStack = old }()
		e := New(WithText("skipped"), WithSkip(1))
		// 不严格断言栈帧位置（依赖 runtime），只确保栈被抓到
		So(e.Stack(), ShouldContainSubstring, "skipped")
	})

	Convey("WithTimeout 写 setTimeout 标志 + timeout 值", t, func() {
		e := New(WithTimeout(true))
		So(e.Timeout(), ShouldBeTrue)
		e2 := New(WithTimeout(false))
		So(e2.Timeout(), ShouldBeFalse)
	})

	Convey("SetLogic / UnsetLogic 切换 logic 字段", t, func() {
		e := New()
		So(e.Logic(), ShouldBeFalse)
		e.SetLogic()
		So(e.Logic(), ShouldBeTrue)
		e.UnsetLogic()
		So(e.Logic(), ShouldBeFalse)
	})
}
