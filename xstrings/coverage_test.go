package xstrings

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：补 xstrings 包内 0% 或低覆盖的导出函数：
//   - From（之前 0%）：interface{} → string 转换的 type-switch 各分支
//   - CamelCaseSlice（之前 0%）：CamelCase 的 slice 包装
//   - TrimPrefixIgnoreCase（66.7%）：HasPrefix 不匹配早返路径
//   - ValidJSON（80%）：空 / 有效 / 无效 三个分支
//   - Wrap（74.5%）：边界场景

func TestFrom(t *testing.T) {
	Convey("From 处理 nil", t, func() {
		So(From(nil), ShouldEqual, "")
	})

	Convey("From 处理整数家族", t, func() {
		So(From(int(42)), ShouldEqual, "42")
		So(From(int8(-8)), ShouldEqual, "-8")
		So(From(int16(-16)), ShouldEqual, "-16")
		So(From(int32(-32)), ShouldEqual, "-32")
		So(From(int64(-64)), ShouldEqual, "-64")
		So(From(uint(7)), ShouldEqual, "7")
		So(From(uint8(8)), ShouldEqual, "8")
		So(From(uint16(16)), ShouldEqual, "16")
		So(From(uint32(32)), ShouldEqual, "32")
		So(From(uint64(64)), ShouldEqual, "64")
	})

	Convey("From 处理浮点", t, func() {
		So(From(float32(3.14)), ShouldEqual, "3.14")
		So(From(float64(2.718)), ShouldEqual, "2.718")
	})

	Convey("From 处理 bool / string / []byte", t, func() {
		So(From(true), ShouldEqual, "true")
		So(From(false), ShouldEqual, "false")
		So(From("hello"), ShouldEqual, "hello")
		So(From([]byte("bytes")), ShouldEqual, "bytes")
	})

	Convey("From 处理 time.Time", t, func() {
		var zeroT time.Time
		So(From(zeroT), ShouldEqual, "")

		// 非零 time
		t1 := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
		So(From(t1), ShouldNotBeEmpty)

		// *time.Time nil / 非 nil
		var nilT *time.Time
		So(From(nilT), ShouldEqual, "")
		So(From(&t1), ShouldNotBeEmpty)
	})

	Convey("From 处理实现 String() 接口的类型", t, func() {
		So(From(stringer("hi")), ShouldEqual, "stringer:hi")
	})

	Convey("From 处理 error 接口", t, func() {
		So(From(errType("oops")), ShouldEqual, "oops")
	})

	Convey("From 处理 nil chan / map / slice / func / ptr 走 reflect", t, func() {
		var nilChan chan int
		So(From(nilChan), ShouldEqual, "")
		var nilMap map[string]int
		So(From(nilMap), ShouldEqual, "")
		var nilSlice []int
		So(From(nilSlice), ShouldEqual, "")
		var nilFn func()
		So(From(nilFn), ShouldEqual, "")
	})

	Convey("From 处理 derived string type", t, func() {
		type myStr string
		So(From(myStr("custom")), ShouldEqual, "custom")
	})

	Convey("From fallback 走 json.Marshal（struct）", t, func() {
		got := From(struct {
			A int    `json:"a"`
			B string `json:"b"`
		}{A: 1, B: "x"})
		So(got, ShouldEqual, `{"a":1,"b":"x"}`)
	})

	Convey("From json.Marshal 失败 fallback 走 fmt.Sprint", t, func() {
		// 含 chan 的 struct 让 json.Marshal 报错（chan 不可序列化）
		ch := make(chan int)
		got := From(struct {
			Ch chan int
		}{Ch: ch})
		// fmt.Sprint 输出含 chan 地址或 "0xc0xxxxx" 等格式
		So(got, ShouldNotBeEmpty)
		So(got, ShouldContainSubstring, "{")
	})
}

// stringer 实现 fmt.Stringer / iString
type stringer string

func (s stringer) String() string { return "stringer:" + string(s) }

// errType 实现 error
type errType string

func (e errType) Error() string { return string(e) }

func TestCamelCaseSlice(t *testing.T) {
	Convey("CamelCaseSlice 把 slice join 后 CamelCase", t, func() {
		So(CamelCaseSlice([]string{"hello", "world"}), ShouldEqual, "HelloWorld")
		So(CamelCaseSlice([]string{"foo", "bar", "baz"}), ShouldEqual, "FooBarBaz")
		So(CamelCaseSlice([]string{}), ShouldEqual, "")
		So(CamelCaseSlice([]string{"single"}), ShouldEqual, "Single")
	})
}

func TestTrimPrefixIgnoreCase(t *testing.T) {
	Convey("TrimPrefixIgnoreCase 不匹配前缀直接返回原串", t, func() {
		So(TrimPrefixIgnoreCase("hello world", "xyz"), ShouldEqual, "hello world")
	})

	Convey("TrimPrefixIgnoreCase 匹配前缀（不区分大小写）后剥离 + TrimSpace", t, func() {
		So(TrimPrefixIgnoreCase("HELLO world", "hello"), ShouldEqual, "world")
		So(TrimPrefixIgnoreCase("Hello   world", "HELLO"), ShouldEqual, "world")
	})
}

func TestValidJSON_AllBranches(t *testing.T) {
	Convey("ValidJSON 空 bytes 视为有效", t, func() {
		So(ValidJSON([]byte("")), ShouldBeTrue)
		So(ValidJSON(nil), ShouldBeTrue)
	})
	Convey("ValidJSON 合法 JSON", t, func() {
		So(ValidJSON([]byte(`{"a":1}`)), ShouldBeTrue)
		So(ValidJSON([]byte(`[1,2,3]`)), ShouldBeTrue)
		So(ValidJSON([]byte(`"plain string"`)), ShouldBeTrue)
		So(ValidJSON([]byte(`null`)), ShouldBeTrue)
	})
	Convey("ValidJSON 非法 JSON", t, func() {
		So(ValidJSON([]byte(`{not json}`)), ShouldBeFalse)
		So(ValidJSON([]byte(`{"a":}`)), ShouldBeFalse)
	})
}

func TestWrap_EdgeCases(t *testing.T) {
	Convey("Wrap 空字符串", t, func() {
		So(Wrap("", 10), ShouldEqual, "")
	})

	Convey("Wrap 短文本不换行", t, func() {
		So(Wrap("short", 100), ShouldEqual, "short")
	})

	Convey("Wrap 长文本按 word 边界换行", t, func() {
		got := Wrap("hello world foo bar baz", 10)
		// 输出应该含 \n
		So(got, ShouldContainSubstring, "\n")
	})

	Convey("Wrap 含显式 \\n 的多行文本", t, func() {
		got := Wrap("line1\nline2 word3\nline3", 100)
		// 显式换行应保留
		So(got, ShouldContainSubstring, "line1")
		So(got, ShouldContainSubstring, "line2")
		So(got, ShouldContainSubstring, "line3")
	})

	Convey("Wrap 单词超过 lim 不被切（wordBufLen >= lim 早返）", t, func() {
		// 一个超长单词，lim 很小：实现选择不切单词
		got := Wrap("supercalifragilisticexpialidocious word", 5)
		So(got, ShouldContainSubstring, "supercalifragilisticexpialidocious")
	})
}

// CamelCase 还有一处 0.1 边界缺口：含数字开头的字段
func TestCamelCase_WithDigit(t *testing.T) {
	Convey("CamelCase 处理数字开头 / 含数字片段", t, func() {
		// 这种 case 是 protobuf-style identifier 转换中的边界
		So(CamelCase("foo_2bar"), ShouldNotBeEmpty)
		So(CamelCase("2foo"), ShouldNotBeEmpty)
	})
}
