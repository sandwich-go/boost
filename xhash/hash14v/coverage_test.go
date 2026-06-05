package hash14v

import (
	"math"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：补 hash14v 包未覆盖的边界 + 默认路径：
//   - ToV id 溢出（math.MaxUint64 - id < idOffset）
//   - ToId 长度越界（< minLength / > maxLength）
//   - encode UsingReservedBuff=false 路径（默认走 stack-allocated rb [16]byte）
//   - Offset 返回 idOffset
//   - 包级 ToV / ToId / Offset（global converter）

func TestToV_OverflowReturnsNil(t *testing.T) {
	Convey("ToV id 溢出（id+offset > MaxUint64）返 nil", t, func() {
		oo := New(WithHashOffset([]byte("FAAAAAA")))
		// id = MaxUint64 - 1，offset 至少 1，加起来溢出
		v := oo.ToV(math.MaxUint64)
		So(v, ShouldBeNil)
	})
}

func TestToId_LengthBoundaries(t *testing.T) {
	oo := New(WithHashOffset([]byte("FAAAAAA")))

	Convey("ToId v 短于 minLength 返 0", t, func() {
		// minLength == len("FAAAAAA") == 7
		So(oo.ToId([]byte("AB")), ShouldEqual, Id(0))
	})

	Convey("ToId v 长于 maxLength=14 返 0", t, func() {
		long := []byte("AAAAAAAAAAAAAAA") // 15 字符
		So(oo.ToId(long), ShouldEqual, Id(0))
	})

	Convey("ToId v 在范围内可转回原 id", t, func() {
		id := Id(1234)
		v := oo.ToV(id)
		So(v, ShouldNotBeNil)
		So(oo.ToId(v), ShouldEqual, id)
	})
}

func TestEncode_DefaultPath_NoReservedBuff(t *testing.T) {
	Convey("encode UsingReservedBuff=false 走 stack alloc 路径", t, func() {
		// 默认 UsingReservedBuff 为 false
		oo := New(WithHashOffset([]byte("FAAAAAA")))
		conv := oo.(*converter)
		So(conv.spec.GetUsingReservedBuff(), ShouldBeFalse)

		// roundtrip 验证
		id := Id(99)
		v := oo.ToV(id)
		So(v, ShouldNotBeEmpty)
		So(oo.ToId(v), ShouldEqual, id)
	})
}

func TestOffset(t *testing.T) {
	Convey("Offset 返回 idOffset（与 hashOffset 解码一致）", t, func() {
		// 默认 hashOffset = "FAAAAAA"，decode 为某个非 0 值
		oo := New()
		So(oo.Offset(), ShouldNotEqual, Id(0))
		// 自定义 hashOffset = "A"，decode 为 1
		oo2 := New(WithHashOffset([]byte("A")))
		So(oo2.Offset(), ShouldEqual, Id(1))
	})
}

func TestPackageLevel_ToVToIdOffset(t *testing.T) {
	Convey("包级 ToV/ToId/Offset 走 global 默认 converter", t, func() {
		// global 默认 hashOffset="FAAAAAA"，offset 非 0
		So(Offset(), ShouldNotEqual, Id(0))

		// roundtrip 验证 ToV / ToId 互逆（global 默认带 hashKey 加密）
		v := ToV(Id(42))
		So(v, ShouldNotBeEmpty)
		So(ToId(v), ShouldEqual, Id(42))
	})
}

// V.String 是 z.BytesToString 包装，覆盖一下
func TestV_String(t *testing.T) {
	Convey("V.String 把 []byte 转 string", t, func() {
		v := V("HELLO")
		So(v.String(), ShouldEqual, "HELLO")
	})
}
