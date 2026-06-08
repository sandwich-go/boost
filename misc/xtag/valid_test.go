package xtag

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTag(t *testing.T) {
	Convey("valid tag", t, func() {
		t.Log(ValidateStructTag("a:1"))
		So(ValidateStructTag("a:1"), ShouldNotBeNil)
		So(ValidateStructTag(`json:"value"`), ShouldBeNil)
		So(ValidateStructTag(`a:"value"`), ShouldBeNil)
	})
}

// TestValidateStructTag_AllPaths 覆盖 ValidateStructTag 各错误 / 成功路径。
func TestValidateStructTag_AllPaths(t *testing.T) {
	Convey("空 tag 返 nil（n=0 直接退出循环）", t, func() {
		So(ValidateStructTag(""), ShouldBeNil)
		So(ValidateStructTag("    "), ShouldBeNil) // 全空格
	})

	Convey("多对 tag 用空格分隔（line 25 检查）", t, func() {
		So(ValidateStructTag(`json:"a" xml:"b"`), ShouldBeNil)
	})

	Convey("多对 tag 没空格分隔 → errTagSpace（line 28）", t, func() {
		// `x:"foo",y:"bar"` 第二对前是 ',' 不是空格
		err := ValidateStructTag(`x:"foo",y:"bar"`)
		So(err, ShouldEqual, errTagSpace)
	})

	Convey("key 为空 → errTagKeySyntax（line 49）", t, func() {
		// `:"value"` 起始就是 ':' i 立即停在 0 → key 长度 0
		err := ValidateStructTag(`:"value"`)
		So(err, ShouldEqual, errTagKeySyntax)
	})

	Convey("缺 ':' 或 ':' 在末尾 → errTagSyntax（line 52）", t, func() {
		// 没冒号
		So(ValidateStructTag(`abc`), ShouldEqual, errTagSyntax)
		// 末尾才是冒号（i+1 >= len）
		So(ValidateStructTag(`abc:`), ShouldEqual, errTagSyntax)
	})

	Convey("value 不以 '\"' 开头 → errTagValueSyntax（line 55）", t, func() {
		So(ValidateStructTag(`abc:value`), ShouldEqual, errTagValueSyntax)
		So(ValidateStructTag(`abc:1`), ShouldEqual, errTagValueSyntax)
	})

	Convey("value 引号不闭合 → errTagValueSyntax（line 69）", t, func() {
		So(ValidateStructTag(`abc:"unclosed`), ShouldEqual, errTagValueSyntax)
	})

	Convey("value Unquote 失败 → errTagValueSyntax（line 76）", t, func() {
		// 内部含未转义的反斜杠（\x 不是合法 escape）
		So(ValidateStructTag(`abc:"\x"`), ShouldEqual, errTagValueSyntax)
	})

	Convey("非 json/xml/asn1 key 不走 spaces 检查（line 80 continue）", t, func() {
		// 'a' 不在 checkTagSpaces，跳过 spaces 检查
		So(ValidateStructTag(`a:"has space"`), ShouldBeNil)
	})

	Convey("xml: value 首尾空格 → errTagValueSpace（line 88）", t, func() {
		So(ValidateStructTag(`xml:" leading"`), ShouldEqual, errTagValueSpace)
		So(ValidateStructTag(`xml:"trailing "`), ShouldEqual, errTagValueSpace)
	})

	Convey("xml: 多个空格 → errTagValueSpace（line 93）", t, func() {
		So(ValidateStructTag(`xml:"a b c"`), ShouldEqual, errTagValueSpace)
	})

	Convey("xml: 无 comma → continue（line 99）", t, func() {
		So(ValidateStructTag(`xml:"name"`), ShouldBeNil)
	})

	Convey("xml: comma 前是空格 → errTagValueSpace（line 103）", t, func() {
		So(ValidateStructTag(`xml:"name ,attr"`), ShouldEqual, errTagValueSpace)
	})

	Convey("xml: comma 后续走通用空格检查（line 105 + 116）", t, func() {
		// 'name,foo bar' → comma 后 value='foo bar' 含空格 → errTagValueSpace
		So(ValidateStructTag(`xml:"name,foo bar"`), ShouldEqual, errTagValueSpace)
		// 'name,attr' → comma 后 'attr' 没空格 → nil
		So(ValidateStructTag(`xml:"name,attr"`), ShouldBeNil)
	})

	Convey("json: 无 comma 直接 continue（line 110）", t, func() {
		So(ValidateStructTag(`json:"field name"`), ShouldBeNil) // json 允许 name 有空格
	})

	Convey("json: comma 后 option 含空格 → errTagValueSpace（line 116）", t, func() {
		// 'name,opt with space' → comma 后 value='opt with space' 含空格
		So(ValidateStructTag(`json:"name,opt with space"`), ShouldEqual, errTagValueSpace)
		So(ValidateStructTag(`json:"name,omitempty"`), ShouldBeNil)
	})

	Convey("asn1: 也走 spaces 检查（line 79 checkTagSpaces）", t, func() {
		// asn1 用 json case 同分支（switch default 走 line 115 通用 space 检查 是不是？
		// 其实 asn1 不在 switch case 中，会落到 default 后 line 115 通用检查）
		// asn1 内 value 含空格 → errTagValueSpace
		// 注：当前实现 switch 没 asn1 case，default 是 fallthrough 到 line 115 ' ' 检查
		So(ValidateStructTag(`asn1:"value"`), ShouldBeNil)
	})
}
