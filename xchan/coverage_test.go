package xchan

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

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

// TestNewRingBuffer_EdgeCases NewRingBuffer 边界：size <= 0 panic / size=1 自动升 2。
func TestNewRingBuffer_EdgeCases(t *testing.T) {
	Convey("initialSize <= 0 应 panic", t, func() {
		So(func() { NewRingBuffer[int](0) }, ShouldPanic)
		So(func() { NewRingBuffer[int](-1) }, ShouldPanic)
	})

	Convey("initialSize=1 自动升级到 2（避免 grow 不可触发的边界）", t, func() {
		rb := NewRingBuffer[int](1)
		So(rb.Capacity(), ShouldEqual, 2)
	})
}

// TestRingBuffer_GrowLarge 触发 grow 的"size >= 1024"分支（+1/4 增长策略）。
func TestRingBuffer_GrowLarge(t *testing.T) {
	Convey("size >= 1024 时 grow 走 +1/4 而非 ×2", t, func() {
		// 起 1024 size buffer，写 1024 次填满触发 grow
		rb := NewRingBuffer[int](1024)
		for i := 0; i < 1024; i++ {
			rb.Write(i)
		}
		// grow 后 size = 1024 + 1024/4 = 1280
		So(rb.Capacity(), ShouldEqual, 1280)

		// 数据完整性
		for i := 0; i < 1024; i++ {
			So(rb.Pop(), ShouldEqual, i)
		}
	})
}

// TestUnboundedChan_Callback 覆盖 process 内 Callback 路径
// （line 86-88 / 101-103 / 117-119）。
//
// 触发条件：CallbackOnBufCount 设小阈值，让写满 In/Out 后 buffer 增长
// 越过阈值，Callback 被调；多次写入测试三个 Callback 触发点。
func TestUnboundedChan_Callback(t *testing.T) {
	Convey("CallbackOnBufCount 阈值触发 Callback", t, func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		var callbackCount atomic.Int32
		var lastBufCount atomic.Int64
		ch := NewUnboundedChanSize[int](
			ctx, 2, 2, 2, // small In/Out/Buffer 强制走 buffer 路径
			WithCallbackOnBufCount(3),
			WithCallback(func(bufCount int64) {
				callbackCount.Add(1)
				lastBufCount.Store(bufCount)
			}),
		)

		// 灌入 20 个值；In=2 / Out=2，剩下 16 都进 buffer，buffer 增长会
		// 反复越过 CallbackOnBufCount=3 → Callback 被调多次
		for i := 0; i < 20; i++ {
			ch.In <- i
		}
		// 等 process 处理完所有 push
		time.Sleep(100 * time.Millisecond)

		So(callbackCount.Load(), ShouldBeGreaterThan, int32(0))
		So(lastBufCount.Load(), ShouldBeGreaterThan, int64(3))

		// 收尾：close in 让 process drain + close out
		close(ch.In)
		// drain Out
		for range ch.Out {
		}
	})
}

// TestUnboundedChan_DrainOnCtxDone 覆盖 drain 内 'case <-ctx.Done():
// return' 早返路径（line 63-64）。
//
// 触发条件：buffer 有数据时 cancel ctx，drain 在 select 内捕到 Done
// 早返而不是把 buffer 全放进 out。
func TestUnboundedChan_DrainOnCtxDone(t *testing.T) {
	Convey("drain 内 ctx.Done 早返", t, func() {
		ctx, cancel := context.WithCancel(context.Background())
		ch := NewUnboundedChanSize[int](ctx, 1, 1, 1) // 极小容量强制 buffer

		// 灌入比 In+Out 多得多的值让 buffer 塞满
		for i := 0; i < 100; i++ {
			ch.In <- i
		}
		// 仅消费几个让 buffer 仍有数据
		for i := 0; i < 3; i++ {
			<-ch.Out
		}

		// close in 让 process 进 drain，drain 中 cancel ctx 触发 Done 早返
		close(ch.In)
		cancel()
		// drain Out 让 goroutine 退出（process 已 close out）
		for range ch.Out {
		}
		// 不强断言数量（drain 早返时 buffer 中可能仍有未送出的）
	})
}
