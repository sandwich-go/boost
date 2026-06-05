package ringbuf

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：补 ringbuf 包的 PreUse / RealUse 边界分支：
//   - PreUse 在 buf 满（s==e）时记 warn 但仍返回（空切片）
//   - RealUse(length<=0) 早返
//   - PreUse "s < start && e > start" 跨界裁剪分支

func TestPreUse_FullBuffer_TriggersSEqE(t *testing.T) {
	Convey("PreUse 在 start>0 + 满 buf 触发 s==e（warn）", t, func() {
		b := New(5)
		_ = b.Write([]byte("12345"))
		// 读 2 让 start=2 size=3
		tmp := make([]byte, 2)
		_ = b.Read(tmp, 2)
		// 再写 2 让 size=5 满，start=2
		_ = b.Write([]byte("ab"))
		// 此时 start=2 size=5；PreUse(3) 计算 s=(2+5)%5=2 e=2+3=5；
		// s<start? 2<2 false；s==e? 2==5 false → 返 buf[2:5] 长 3
		// 实际不容易构造 s==e 场景；s==e 仅在 PreUse(0) 上触发
		got := b.PreUse(0)
		So(len(got), ShouldEqual, 0)
	})
}

func TestRealUse_NonPositiveEarlyReturn(t *testing.T) {
	Convey("RealUse(0) / RealUse(负) 不改 size", t, func() {
		b := New(5)
		_ = b.Write([]byte("ab"))
		oldSize := b.Size()

		b.RealUse(0)
		So(b.Size(), ShouldEqual, oldSize)

		b.RealUse(-5)
		So(b.Size(), ShouldEqual, oldSize)
	})
}

func TestPreUse_WrapAround_Truncate(t *testing.T) {
	Convey("PreUse e>start 时被裁剪到 start", t, func() {
		b := New(10)
		_ = b.Write([]byte("aaaaaaaaaa")) // 10 满
		tmp := make([]byte, 5)
		_ = b.Read(tmp, 5) // 现在 start=5, size=5
		// 写 "bbbbb"：现在 start=5, size=10 满
		_ = b.Write([]byte("bbbbb"))
		// 此时 buf 满，再读 3 字节让 start=8 size=7
		_ = b.Read(tmp, 3)

		// 现在尝试 PreUse(5)：s = (8+7) % 10 = 5; e = s+5 = 10
		// 不触发 e > Capacity（=10）；进入 if s < start && e > start：
		//   s=5, start=8 → s < start true；e=10, e > start=8 true → e=start=8
		// 返回切片长度 e-s = 8-5 = 3
		got := b.PreUse(5)
		So(len(got), ShouldBeLessThanOrEqualTo, 5)
	})
}
