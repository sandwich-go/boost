package xtime

import (
	"fmt"
	"math"
	"sync/atomic"
	"time"

	"github.com/sandwich-go/boost/internal/log"
	"github.com/sandwich-go/boost/xpanic"

	"github.com/sandwich-go/boost/xsync"
)

var (
	// CopToleranceSecond 时间容忍误差，误差超过CopToleranceSecond则后续使用系统时间
	CopToleranceSecond = int64(2)
	// CopToleranceCheckInterval 检测CopToleranceSecond的时间间隔
	CopToleranceCheckInterval = time.Duration(10) * time.Second
)

// providerHolder 包一层让 atomic.Pointer 能装 func（func 自身不能直接做 atomic.Pointer 的 T）。
type providerHolder struct {
	fn func() time.Time
}

// Cop mock time like ruby's time cop
type Cop struct {
	ts        xsync.AtomicInt64
	unhealthy chan struct{}
	closeChan chan struct{}
	running   xsync.AtomicBool
	closeFlag xsync.AtomicInt32
	// nowProvider 用 atomic.Pointer 保护：SetNowProvider 与后台 start/check
	// 以及业务侧 Now() 之间存在并发读写。用 atomic.Pointer 让读侧无锁、写
	// 侧 Store 原子可见。
	nowProvider atomic.Pointer[providerHolder]
}

// NewCop 新建Cop对象
func NewCop(nowProvider func() time.Time) *Cop {
	tc := &Cop{}
	tc.nowProvider.Store(&providerHolder{fn: nowProvider})
	tc.closeChan = make(chan struct{})
	tc.unhealthy = make(chan struct{})
	tc.run()
	return tc
}

// SetNowProvider 原子替换 nowProvider；与后台 goroutine / Now() 并发安全。
func (tc *Cop) SetNowProvider(nowProvider func() time.Time) {
	tc.nowProvider.Store(&providerHolder{fn: nowProvider})
}

// getNowProvider 原子读 nowProvider 当前 fn。
func (tc *Cop) getNowProvider() func() time.Time {
	return tc.nowProvider.Load().fn
}

func (tc *Cop) run() {
	go xpanic.AutoRecover("stime_cop", tc.start)
	go xpanic.AutoRecover("stime_cop_check", tc.check)
}

func (tc *Cop) check() {
	checkTicker := time.NewTicker(CopToleranceCheckInterval)
	defer checkTicker.Stop()
	for {
		select {
		case <-checkTicker.C:
			systemTime := tc.getNowProvider()().Unix()
			copTime := tc.Unix()
			if int64(math.Abs(float64(systemTime-copTime))) > CopToleranceSecond {
				log.Warn(fmt.Sprintf("stime cop tolerance, system: %d, cop: %d, tolerance: %d", systemTime, copTime, CopToleranceSecond))
				tc.unhealthy <- struct{}{}
			}
			return
		case <-tc.closeChan:
			return
		}
	}
}

func (tc *Cop) start() {
	tc.running.Set(true)
	defer func() {
		tc.running.Set(false)
	}()
	tc.ts.Set(tc.getNowProvider()().Unix())
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case now := <-ticker.C:
			tc.ts.Set(now.Unix())
		case <-tc.closeChan:
			return
		case <-tc.unhealthy:
			// 处于不健康状态，暂停使用
			return
		}
	}
}

// Stop 停止模拟时间
func (tc *Cop) Stop() {
	if tc.closeFlag.CompareAndSwap(0, 1) {
		close(tc.closeChan)
		tc.running.Set(false)
	}
}

// Now Cop获取的最小单位为秒,精度低但是效率高, time.Now()最小单位为纳秒
func (tc *Cop) Now() time.Time {
	if tc.running.Get() {
		return time.Unix(tc.ts.Get(), 0)
	}
	return tc.getNowProvider()()
}

// Unix 获取当前Unix时间戳
func (tc *Cop) Unix() int64 { return tc.Now().Unix() }

// UnixMilli 获取当前时间，单位毫秒
func (tc *Cop) UnixMilli() int64 {
	return tc.Now().UnixNano() / int64(time.Millisecond)
}
