package cron

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：补 cron 包内未覆盖的 Parse 错误路径 + matchDay 各分支：
//   - Parse 错误：字段数错 / 各 field 值非法
//   - parseCronField 各分支：太多 slash / hyphen / 非数字 / 越界
//   - MustParse 在错误时 panic
//   - matchDay 三个分支：dom 全设 / dow 全设 / 都部分设
//
// 已有 cronexpr_test.go 覆盖 Next 主流程的大量正面 case。

func TestParse_ErrorPaths(t *testing.T) {
	Convey("字段数错（不是 5 或 6）", t, func() {
		_, err := Parse("* * *")
		So(err, ShouldNotBeNil)

		_, err = Parse("* * * * * * *")
		So(err, ShouldNotBeNil)

		_, err = Parse("")
		So(err, ShouldNotBeNil)
	})

	Convey("各字段非法", t, func() {
		// seconds 字段非法
		_, err := Parse("xx * * * * *")
		So(err, ShouldNotBeNil)

		// minutes 字段
		_, err = Parse("* xx * * * *")
		So(err, ShouldNotBeNil)

		// hours
		_, err = Parse("* * xx * * *")
		So(err, ShouldNotBeNil)

		// day-of-month
		_, err = Parse("* * * xx * *")
		So(err, ShouldNotBeNil)

		// month
		_, err = Parse("* * * * xx *")
		So(err, ShouldNotBeNil)

		// day-of-week
		_, err = Parse("* * * * * xx")
		So(err, ShouldNotBeNil)
	})

	Convey("parseCronField 各错误分支", t, func() {
		// 太多 slash：a/b/c
		_, err := Parse("1/2/3 * * * * *")
		So(err, ShouldNotBeNil)

		// 太多 hyphen：a-b-c
		_, err = Parse("1-2-3 * * * * *")
		So(err, ShouldNotBeNil)

		// "*-1" 非法：* 后跟 -
		_, err = Parse("*-1 * * * * *")
		So(err, ShouldNotBeNil)

		// 非数字 start
		_, err = Parse("xx-2 * * * * *")
		So(err, ShouldNotBeNil)

		// 非数字 end
		_, err = Parse("1-xx * * * * *")
		So(err, ShouldNotBeNil)

		// 非数字 incr
		_, err = Parse("1-2/xx * * * * *")
		So(err, ShouldNotBeNil)
	})
}

func TestMustParse_PanicOnInvalid(t *testing.T) {
	Convey("MustParse 在 Parse 出错时 panic", t, func() {
		So(func() { MustParse("invalid") }, ShouldPanic)
		So(func() { MustParse("xx * * * * *") }, ShouldPanic)
	})

	Convey("MustParse 合法表达式正常返回", t, func() {
		expr := MustParse("* * * * * *")
		So(expr, ShouldNotBeNil)
	})
}

func TestMatchDay_AllBranches(t *testing.T) {
	Convey("matchDay：dom=*（全设），仅按 dow 匹配", t, func() {
		// "0 0 * * 1" = 每周一 00:00；dom 留空 = 0xfffffffe
		expr := MustParse("0 0 * * 1")
		// 2024-01-01 是周一
		monday := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		// Next 应该返回当周一（已经是周一 0:00 → 下个周一）或下周一
		next := expr.Next(monday)
		So(next.Weekday(), ShouldEqual, time.Monday)
	})

	Convey("matchDay：dow=*（全设），仅按 dom 匹配", t, func() {
		// "0 0 1 * *" = 每月 1 号 00:00；dow 留空 = 0x7f
		expr := MustParse("0 0 1 * *")
		from := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		next := expr.Next(from)
		So(next.Day(), ShouldEqual, 1)
		So(next.Month(), ShouldEqual, time.February)
	})

	Convey("matchDay：dom 与 dow 都部分指定，OR 匹配", t, func() {
		// "0 0 1 * 1" = 每月 1 号 OR 每周一；OR 语义按 cron 标准
		expr := MustParse("0 0 1 * 1")
		// 2024-01-01 既是 1 号又是周一，肯定命中
		from := time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)
		next := expr.Next(from)
		// 期望命中 2024-01-01 00:00
		So(next.Year(), ShouldEqual, 2024)
		So(next.Month(), ShouldEqual, time.January)
		So(next.Day(), ShouldEqual, 1)
	})
}

func TestParse_5FieldVariant(t *testing.T) {
	Convey("5 字段表达式自动补 0 秒", t, func() {
		// "* * * * *" = 每分钟（5 字段，不含秒）
		expr := MustParse("* * * * *")
		from := time.Date(2024, 1, 1, 0, 0, 30, 0, time.UTC)
		next := expr.Next(from)
		// 下一分钟 0 秒
		So(next, ShouldEqual, time.Date(2024, 1, 1, 0, 1, 0, 0, time.UTC))
	})
}

// TestNext_NoMatchInTwoYears 覆盖 Next 在两年内找不到匹配时返回零 time。
// 难以构造完全不可达的合法 expression（cron 表达式语法限制内不允许"永
// 不命中"），但下面这个 30 号 + 2 月组合永远不可能，超过 year+1 后返零。
func TestNext_NoMatchInTwoYears(t *testing.T) {
	Convey("Next 找不到 2 月 30 号 → 返零 time", t, func() {
		// 2 月 30 号永远不存在。
		// 6 字段：sec min hour dom month dow
		// "0 0 0 30 2 *" = 每年 2 月 30 号 0:0:0；dom=30 month=2，永不可达
		// 注意：dow 必须 *（全设 0x7f）才走 matchDay 的 dow 全设分支，纯按
		// dom 匹配（dom=30 在 2 月永远 false）
		expr := MustParse("0 0 0 30 2 *")
		from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		next := expr.Next(from)
		So(next.IsZero(), ShouldBeTrue)
	})
}

func TestParse_RangeWithStep(t *testing.T) {
	Convey("Range with step：5/10 = 5,15,25,...", t, func() {
		// 0/15 * * * * * = 0,15,30,45 秒触发
		expr := MustParse("0/15 * * * * *")
		from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		next := expr.Next(from)
		So(next.Second(), ShouldEqual, 15)
	})

	Convey("*/5 = 5 每 5 单位", t, func() {
		expr := MustParse("*/10 * * * *")
		from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		next := expr.Next(from)
		So(next.Minute(), ShouldEqual, 10)
	})
}
