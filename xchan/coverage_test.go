package xchan

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：补 xchan 包未覆盖的 RingBuffer API + 边界。

func TestRingBuffer_BasicAPI(t *testing.T) {
	Convey("Capacity 返回当前 buf 大小", t, func() {
		rb := NewRingBuffer[int](8)
		So(rb.Capacity(), ShouldEqual, 8)
	})

	Convey("Len 计算逻辑：r==w 为 0；w>r 为 w-r；w<r 为 size-r+w（环绕）", t, func() {
		rb := NewRingBuffer[int](4)
		So(rb.Len(), ShouldEqual, 0)
		So(rb.IsEmpty(), ShouldBeTrue)

		// 写 2 个：r=0 w=2 → Len=2
		rb.Write(1)
		rb.Write(2)
		So(rb.Len(), ShouldEqual, 2)

		// Pop 1 个：r=1 w=2 → Len=1
		So(rb.Pop(), ShouldEqual, 1)
		So(rb.Len(), ShouldEqual, 1)

		// Reset 重置
		rb.Reset()
		So(rb.Len(), ShouldEqual, 0)
		So(rb.Capacity(), ShouldEqual, 4) // initialSize
	})
}

func TestRingBuffer_PopPeekOnEmpty(t *testing.T) {
	Convey("Pop 在 empty 上 panic", t, func() {
		rb := NewRingBuffer[int](4)
		So(func() { rb.Pop() }, ShouldPanic)
	})

	Convey("Peek 在 empty 上 panic", t, func() {
		rb := NewRingBuffer[int](4)
		So(func() { rb.Peek() }, ShouldPanic)
	})

	Convey("Peek 不消耗 item，多次 Peek 返同值", t, func() {
		rb := NewRingBuffer[int](4)
		rb.Write(42)
		So(rb.Peek(), ShouldEqual, 42)
		So(rb.Peek(), ShouldEqual, 42)
		// Peek 不改 r，Pop 还是 42
		So(rb.Pop(), ShouldEqual, 42)
	})
}

func TestRingBuffer_Wrap(t *testing.T) {
	Convey("Write/Pop 后 Len 计算正确（含环绕场景）", t, func() {
		rb := NewRingBuffer[int](4)
		rb.Write(1)
		rb.Write(2)
		rb.Write(3)
		So(rb.Pop(), ShouldEqual, 1)
		So(rb.Pop(), ShouldEqual, 2)
		// 现在 r=2 w=3 → Len=1
		So(rb.Len(), ShouldEqual, 1)
		// 再写 2 个让 w 环绕（grow 也可能发生）
		rb.Write(4)
		rb.Write(5)
		So(rb.Len(), ShouldEqual, 3)
	})
}

func TestRingBuffer_Grow(t *testing.T) {
	Convey("Write 触发 grow（写入数 > size）", t, func() {
		rb := NewRingBuffer[int](2)
		// 连续写 10 个，肯定触发 grow
		for i := 0; i < 10; i++ {
			rb.Write(i)
		}
		So(rb.Len(), ShouldEqual, 10)
		So(rb.Capacity(), ShouldBeGreaterThan, 2)

		// 读出顺序应 0..9（FIFO）
		for i := 0; i < 10; i++ {
			So(rb.Pop(), ShouldEqual, i)
		}
	})
}

func TestRingBuffer_Read_OneAtATime(t *testing.T) {
	Convey("Read 单次读出一个 item（FIFO）", t, func() {
		rb := NewRingBuffer[int](8)
		for i := 1; i <= 5; i++ {
			rb.Write(i)
		}

		for i := 1; i <= 5; i++ {
			v, err := rb.Read()
			So(err, ShouldBeNil)
			So(v, ShouldEqual, i)
		}

		// Read 完应该 empty
		So(rb.IsEmpty(), ShouldBeTrue)

		// 再 Read 返 ErrIsEmpty
		_, err := rb.Read()
		So(err, ShouldNotBeNil)
	})
}
