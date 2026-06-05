package xsync

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：补 xsync 包内 0% 覆盖的导出函数：
//   - AtomicInt32 / AtomicUint32 / AtomicUint64 / AtomicDuration 各方法
//   - RWTimeoutLock 5 个方法（1.3 新增"超时锁"，AGENTS.md §7 列出）
//   - WaitContext 缺失分支
//
// 已被 atomic_test / cond_test / wg_test 覆盖的方法不重复测。

// TestAtomicInt32 覆盖 AtomicInt32 全方法（Add/Set/Get/CompareAndSwap）。
func TestAtomicInt32(t *testing.T) {
	Convey("AtomicInt32 单线程契约", t, func() {
		var a AtomicInt32
		So(a.Get(), ShouldEqual, int32(0))

		a.Set(42)
		So(a.Get(), ShouldEqual, int32(42))

		So(a.Add(8), ShouldEqual, int32(50))
		So(a.Get(), ShouldEqual, int32(50))

		// CompareAndSwap 成功
		So(a.CompareAndSwap(50, 100), ShouldBeTrue)
		So(a.Get(), ShouldEqual, int32(100))

		// CompareAndSwap 失败（old 不匹配）
		So(a.CompareAndSwap(50, 999), ShouldBeFalse)
		So(a.Get(), ShouldEqual, int32(100))
	})

	Convey("AtomicInt32 并发 Add 累加正确性", t, func() {
		var a AtomicInt32
		const goroutines = 100
		const perG = 100
		var wg sync.WaitGroup
		wg.Add(goroutines)
		for g := 0; g < goroutines; g++ {
			go func() {
				defer wg.Done()
				for i := 0; i < perG; i++ {
					a.Add(1)
				}
			}()
		}
		wg.Wait()
		So(a.Get(), ShouldEqual, int32(goroutines*perG))
	})
}

// TestAtomicUint32 覆盖 AtomicUint32 全方法。
func TestAtomicUint32(t *testing.T) {
	Convey("AtomicUint32 单线程契约", t, func() {
		var a AtomicUint32
		So(a.Get(), ShouldEqual, uint32(0))
		a.Set(123)
		So(a.Get(), ShouldEqual, uint32(123))
		So(a.Add(7), ShouldEqual, uint32(130))
		So(a.CompareAndSwap(130, 200), ShouldBeTrue)
		So(a.CompareAndSwap(0, 999), ShouldBeFalse)
		So(a.Get(), ShouldEqual, uint32(200))
	})
}

// TestAtomicUint64 覆盖 AtomicUint64 全方法。
func TestAtomicUint64(t *testing.T) {
	Convey("AtomicUint64 单线程契约", t, func() {
		var a AtomicUint64
		So(a.Get(), ShouldEqual, uint64(0))
		a.Set(1<<40 + 7)
		So(a.Get(), ShouldEqual, uint64(1<<40+7))
		So(a.Add(1), ShouldEqual, uint64(1<<40+8))
		So(a.CompareAndSwap(1<<40+8, 999), ShouldBeTrue)
		So(a.CompareAndSwap(0, 1234), ShouldBeFalse)
		So(a.Get(), ShouldEqual, uint64(999))
	})
}

// TestAtomicInt64_Set 补 Int64.Set（其他方法已被现有测试覆盖）。
func TestAtomicInt64_Set(t *testing.T) {
	Convey("AtomicInt64.Set", t, func() {
		var a AtomicInt64
		a.Set(-1234567890)
		So(a.Get(), ShouldEqual, int64(-1234567890))
	})
}

// TestAtomicDuration 覆盖 AtomicDuration 全方法。
func TestAtomicDuration(t *testing.T) {
	Convey("AtomicDuration 单线程契约", t, func() {
		var d AtomicDuration
		So(d.Get(), ShouldEqual, time.Duration(0))

		d.Set(500 * time.Millisecond)
		So(d.Get(), ShouldEqual, 500*time.Millisecond)

		got := d.Add(500 * time.Millisecond)
		So(got, ShouldEqual, time.Second)
		So(d.Get(), ShouldEqual, time.Second)

		So(d.CompareAndSwap(time.Second, 2*time.Second), ShouldBeTrue)
		So(d.CompareAndSwap(0, 99*time.Second), ShouldBeFalse)
		So(d.Get(), ShouldEqual, 2*time.Second)
	})
}

// TestRWTimeoutLock 覆盖 RWTimeoutLock 全部 5 个方法。
//
// RWTimeoutLock 基于 semaphore.Weighted：Lock = 拿全部权重；RLock = 拿 1。
// 关键不变量：
//  1. 多读并发可同时持有
//  2. 写排他：写持有时读必等
//  3. ctx 超时正确传播
func TestRWTimeoutLock(t *testing.T) {
	Convey("NewRWTimeoutLock(0) panic", t, func() {
		So(func() { NewRWTimeoutLock(0) }, ShouldPanic)
		So(func() { NewRWTimeoutLock(-1) }, ShouldPanic)

		// 进一步验证 panic value 是 error 类型且消息正确（防止 panic 实现被
		// 误改为字符串）
		defer func() {
			r := recover()
			err, ok := r.(error)
			So(ok, ShouldBeTrue)
			So(err.Error(), ShouldEqual, "maxReaders must be greater than zero")
		}()
		NewRWTimeoutLock(0)
	})

	Convey("单读 RLock + RUnlock", t, func() {
		l := NewRWTimeoutLock(4)
		ctx := context.Background()
		So(l.RLock(ctx), ShouldBeNil)
		l.RUnlock()
	})

	Convey("多读并发持有（最多 maxReaders 个）", t, func() {
		const max int64 = 4
		l := NewRWTimeoutLock(max)
		ctx := context.Background()
		// 4 个 reader 都能拿到
		for i := int64(0); i < max; i++ {
			So(l.RLock(ctx), ShouldBeNil)
		}
		// 第 5 个在 1ms ctx 超时下应该失败
		ctxShort, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
		defer cancel()
		So(l.RLock(ctxShort), ShouldNotBeNil) // ctx.Err() = DeadlineExceeded

		// 释放一个，第 5 个能拿到了
		l.RUnlock()
		So(l.RLock(ctx), ShouldBeNil)

		// 收尾全部释放（4 个 reader）
		for i := int64(0); i < max; i++ {
			l.RUnlock()
		}
	})

	Convey("Lock 排他：写持有时读必等到超时", t, func() {
		l := NewRWTimeoutLock(4)
		ctx := context.Background()
		So(l.Lock(ctx), ShouldBeNil)

		// 写锁持有时，RLock 立即超时失败
		ctxShort, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
		defer cancel()
		So(l.RLock(ctxShort), ShouldNotBeNil)

		l.Unlock()

		// 释放后 RLock 立即成功
		So(l.RLock(ctx), ShouldBeNil)
		l.RUnlock()
	})

	Convey("Lock 在已有 reader 时 ctx 超时返回错误", t, func() {
		l := NewRWTimeoutLock(2)
		ctx := context.Background()
		So(l.RLock(ctx), ShouldBeNil) // 持有 1 个 reader

		// Lock 需要 2 个权重，但只剩 1 个 → ctx 超时
		ctxShort, cancel := context.WithTimeout(context.Background(), 5*time.Millisecond)
		defer cancel()
		So(l.Lock(ctxShort), ShouldNotBeNil)

		l.RUnlock()
	})
}

// TestWaitContext_Branches 覆盖 WaitContext 两条分支：
//   - 返回 true：ctx 超时/取消（ctx.Done() 先被关闭）
//   - 返回 false：wg.Wait 正常完成
//
// 注意：这与 stdlib WaitGroup 风格相反。boost 的语义是"返回 true 表示
// 出错（超时）"。
func TestWaitContext_Branches(t *testing.T) {
	Convey("WaitContext ctx 已取消时返回 true（超时分支）", t, func() {
		var wg sync.WaitGroup
		wg.Add(1)
		// 永不调 wg.Done

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // 立即取消

		timedOut := WaitContext(&wg, ctx)
		So(timedOut, ShouldBeTrue)
	})

	Convey("WaitContext wg 完成时返回 false（成功分支）", t, func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			time.Sleep(10 * time.Millisecond)
			wg.Done()
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		timedOut := WaitContext(&wg, ctx)
		So(timedOut, ShouldBeFalse)
	})
}

// BenchmarkAtomicInt32_Add hot path bench：原子计数是高频用法。
// 与 sync/atomic.AddInt32 对比，量化 boost AtomicInt32 包装成本。
func BenchmarkAtomicInt32_Add(b *testing.B) {
	var a AtomicInt32
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		a.Add(1)
	}
}

// BenchmarkAtomicInt32_Add_Stdlib 对照：直接 sync/atomic
func BenchmarkAtomicInt32_Add_Stdlib(b *testing.B) {
	var a int32
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		atomic.AddInt32(&a, 1)
	}
}

// BenchmarkRWTimeoutLock_RLock_RUnlock：semaphore-based 读锁与 sync.RWMutex
// 对比，量化超时支持的开销代价。
func BenchmarkRWTimeoutLock_RLock_RUnlock(b *testing.B) {
	l := NewRWTimeoutLock(int64(b.N) + 1)
	ctx := context.Background()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = l.RLock(ctx)
		l.RUnlock()
	}
}

// BenchmarkRWMutex_RLock_RUnlock_Stdlib 对照
func BenchmarkRWMutex_RLock_RUnlock_Stdlib(b *testing.B) {
	var m sync.RWMutex
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		m.RLock()
		m.RUnlock()
	}
}
