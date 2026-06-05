package fnv

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：xhash/fnv 之前 0% 覆盖。
//
// 实现注意（行为契约）：Hash 只对 string 走 fnv1aHash；对所有数值类型
// 直接返回 uint64(v)（负数取反）。这与 fnv1a 名字暗示的"统一 hash"不
// 一致，但是当前实现行为，测试按此契约固定下来。

func TestHash_NumericTypes(t *testing.T) {
	Convey("Hash 数值类型直接返 uint64（不走 fnv1a）", t, func() {
		So(Hash(42), ShouldEqual, uint64(42))
		So(Hash(int(-42)), ShouldEqual, uint64(42)) // 负数取反
		So(Hash(uint(99)), ShouldEqual, uint64(99))

		So(Hash(int32(100)), ShouldEqual, uint64(100))
		So(Hash(int32(-100)), ShouldEqual, uint64(100))
		So(Hash(uint32(200)), ShouldEqual, uint64(200))

		So(Hash(int64(1<<40)), ShouldEqual, uint64(1<<40))
		So(Hash(int64(-(1 << 40))), ShouldEqual, uint64(1<<40))
		So(Hash(uint64(1<<50)), ShouldEqual, uint64(1<<50))

		// float 截断为 uint64（行为契约：浮点小数部分丢失）
		So(Hash(float32(3.0)), ShouldEqual, uint64(3))
		So(Hash(float32(-3.0)), ShouldEqual, uint64(3))
		So(Hash(float64(7.0)), ShouldEqual, uint64(7))
		So(Hash(float64(-7.0)), ShouldEqual, uint64(7))
	})
}

func TestHash_String(t *testing.T) {
	Convey("Hash string 走 fnv1a", t, func() {
		// 同字符串得到相同 hash
		h1 := Hash("hello")
		h2 := Hash("hello")
		So(h1, ShouldEqual, h2)

		// 不同字符串大概率不同 hash（fnv-1a 哈希）
		h3 := Hash("world")
		So(h1, ShouldNotEqual, h3)

		// 空字符串：buf 为空，fnv1a 返回 fnvOffsetBasis
		So(Hash(""), ShouldEqual, fnvOffsetBasis)
	})

	Convey("Hash 同字符串不同 case 不同 hash", t, func() {
		So(Hash("Hello"), ShouldNotEqual, Hash("hello"))
	})
}

func TestHash_UnsupportedTypePanic(t *testing.T) {
	Convey("Hash 不支持的类型 panic", t, func() {
		So(func() { Hash([]int{1, 2, 3}) }, ShouldPanic)
		So(func() { Hash(struct{}{}) }, ShouldPanic)
		So(func() { Hash(true) }, ShouldPanic) // bool 不在 switch case 里
		So(func() { Hash(nil) }, ShouldPanic)
	})
}

func TestFnv1aHash_Direct(t *testing.T) {
	Convey("fnv1aHash 内部函数：空 input 返 offset basis", t, func() {
		So(fnv1aHash(nil), ShouldEqual, fnvOffsetBasis)
		So(fnv1aHash([]byte{}), ShouldEqual, fnvOffsetBasis)
	})

	Convey("fnv1aHash 单字节正确性：基础值 ^ byte * prime（uint64 自然溢出）", t, func() {
		// 算 1 字节 0x01 的 fnv1a:
		// hash = offset ^ 0x01
		// hash *= prime（uint64 mod 2^64 溢出环绕）
		var h uint64 = fnvOffsetBasis
		h ^= 0x01
		h *= fnvPrime
		So(fnv1aHash([]byte{0x01}), ShouldEqual, h)
	})
}

// BenchmarkHash hot path：业务侧 hash 的常见用法是 string key。
func BenchmarkHash_String(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Hash("benchmark-key-string")
	}
}

func BenchmarkHash_Int(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Hash(i)
	}
}
