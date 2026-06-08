package xtime

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestTick 历史抖动（§9.1 doc 化）：原代码 'if count >= 2 { close(stop) }'
// 在 dispatcher tick 间隔 5ms / 主线超时 12ms 之间，dispatcher tick goroutine
// 可能在主线退出前再触发第 3 次 tick，count 已 >= 2 重复 close(stop) 触发
// 'close of closed channel' panic。go test ./... 高并发跨包时调度延后让概率
// 显著上升。
//
// 修复双保险：
// 1. count 改 atomic.Int32（防 dispatcher tick goroutine 与主线 race）
// 2. close(stop) 走 sync.Once 让重复 close 安全
// 3. 主线退出前 defer d.Close() 让 tick 停（避免后续 tick 还在跑触发 panic）
func TestTick(t *testing.T) {
	d := NewDispatcher(0, WithTickDuration(time.Millisecond*5))
	defer d.Close()
	var count atomic.Int32
	stop := make(chan struct{})
	var stopOnce sync.Once
	d.TickFunc("123", func(_ context.Context) {
		if count.Add(1) >= 2 {
			stopOnce.Do(func() { close(stop) })
		}
	})
	d.Start()
	select {
	case <-stop:
		t.Log("trigger tick twice")
	case <-time.After(time.Millisecond * 12):
		t.Fatal("ticker not work")
	}
}

func TestTickExternalHost(t *testing.T) {
	d := NewDispatcher(0, WithTickDuration(time.Millisecond*5), WithTickHostingMode(false))
	defer d.Close()
	var count atomic.Int32
	stop := make(chan struct{})
	var stopOnce sync.Once
	d.TickFunc("123", func(_ context.Context) {
		if count.Add(1) >= 2 {
			stopOnce.Do(func() { close(stop) })
		}
	})
	d.Start()
	td := d.(TickerDispatcher)

	for {
		select {
		case <-td.TickerC():
			td.TriggerTickFuncs(context.Background())
		case <-stop:
			t.Log("trigger tick twice")
			return
		case <-time.After(time.Millisecond * 12):
			t.Fatal("ticker not work")
		}
	}

}
