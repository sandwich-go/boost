package xconv

import (
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：补 xconv 包内部分覆盖函数的剩余分支：
//   - String reflect.Value / iString / error / nil chan / json fail 路径
//   - Int64 hex / octal / float fallback
//   - Uint64 hex / octal / nil checks
//   - Bool 各 case + iBool 接口
//   - Float64 string / []byte 路径
//
// 已有 *_test.go 覆盖主流程；本文件补漏。

// stringerType 实现 iString 接口
type stringerType struct{ s string }

func (st stringerType) String() string { return st.s }

// errType 实现 error
type errType struct{ msg string }

func (e errType) Error() string { return e.msg }

func TestString_Branches(t *testing.T) {
	Convey("String reflect.Value 输入：递归走原始类型", t, func() {
		// 有效 reflect.Value：内部递归调 String(rv.Interface())
		rv := reflect.ValueOf(42)
		So(String(rv), ShouldEqual, "42")

		// 无效 reflect.Value（零值）：直接返空
		var zero reflect.Value
		So(String(zero), ShouldEqual, "")
	})

	Convey("String iString 接口", t, func() {
		So(String(stringerType{s: "hi"}), ShouldEqual, "hi")
	})

	Convey("String error 接口", t, func() {
		So(String(errType{msg: "boom"}), ShouldEqual, "boom")
	})

	Convey("String nil chan / map / slice / func 走 reflect IsNil", t, func() {
		var nilChan chan int
		So(String(nilChan), ShouldEqual, "")
		var nilMap map[string]int
		So(String(nilMap), ShouldEqual, "")
		var nilSlice []int
		So(String(nilSlice), ShouldEqual, "")
		var nilFn func()
		So(String(nilFn), ShouldEqual, "")
	})

	Convey("String derived string type（reflect.String 路径）", t, func() {
		type myStr string
		So(String(myStr("x")), ShouldEqual, "x")
	})

	Convey("String time.Time 零值返空", t, func() {
		var zeroT time.Time
		So(String(zeroT), ShouldEqual, "")
	})

	Convey("String *time.Time nil 返空", t, func() {
		var nilT *time.Time
		So(String(nilT), ShouldEqual, "")
	})

	Convey("String json.Marshal fallback（含 chan 让 marshal 失败）", t, func() {
		// 含 chan field 的 struct，json.Marshal 报错走 fmt.Sprint
		got := String(struct {
			Ch chan int
		}{})
		So(got, ShouldNotBeEmpty)
	})
}

func TestInt64_Branches(t *testing.T) {
	Convey("Int64 hex string", t, func() {
		So(Int64("0x10"), ShouldEqual, int64(16))
		So(Int64("-0xFF"), ShouldEqual, int64(-255))
		So(Int64("+0xA"), ShouldEqual, int64(10))
	})

	Convey("Int64 octal string", t, func() {
		So(Int64("010"), ShouldEqual, int64(8))
		So(Int64("-010"), ShouldEqual, int64(-8))
	})

	Convey("Int64 float string fallback", t, func() {
		// "3.14" 不是合法 int 但能解析为 float
		So(Int64("3.14"), ShouldEqual, int64(3))
	})

	Convey("Int64 实现 iInt64 的类型", t, func() {
		So(Int64(int64Stringer{v: 99}), ShouldEqual, int64(99))
	})
}

// int64Stringer 实现 iInt64 接口
type int64Stringer struct{ v int64 }

func (s int64Stringer) Int64() int64 { return s.v }

func TestUint64_Branches(t *testing.T) {
	Convey("Uint64 hex / octal / 默认", t, func() {
		So(Uint64("0xFF"), ShouldEqual, uint64(255))
		So(Uint64("010"), ShouldEqual, uint64(8))
		So(Uint64("123"), ShouldEqual, uint64(123))
	})

	Convey("Uint64 nil / bool / 各种类型", t, func() {
		So(Uint64(nil), ShouldEqual, uint64(0))
		So(Uint64(true), ShouldEqual, uint64(1))
		So(Uint64(false), ShouldEqual, uint64(0))
	})

	Convey("Uint64 float fallback", t, func() {
		So(Uint64("3.99"), ShouldEqual, uint64(3))
	})
}

func TestBool_Branches(t *testing.T) {
	Convey("Bool 各种 truthy / falsy 输入", t, func() {
		So(Bool(true), ShouldBeTrue)
		So(Bool(false), ShouldBeFalse)
		So(Bool(nil), ShouldBeFalse)
		So(Bool(1), ShouldBeTrue)
		So(Bool(0), ShouldBeFalse)
		So(Bool("true"), ShouldBeTrue)
		So(Bool("false"), ShouldBeFalse)
		// 空字符串视为 false
		So(Bool(""), ShouldBeFalse)
	})
}

func TestFloat64_Branches(t *testing.T) {
	Convey("Float64 string / []byte 路径", t, func() {
		So(Float64("3.14"), ShouldEqual, 3.14)
		So(Float64([]byte("2.718")), ShouldNotEqual, 0.0)
	})

	Convey("Float64 nil 返 0", t, func() {
		So(Float64(nil), ShouldEqual, 0.0)
	})

	Convey("Float64 bool 走 default String 路径返 0（行为契约）", t, func() {
		// Float64 没有 bool case；走 default → String(true)="true" →
		// ParseFloat("true") 失败 → 返 0。这是实现行为契约，doc 化。
		So(Float64(true), ShouldEqual, 0.0)
		So(Float64(false), ShouldEqual, 0.0)
	})
}
