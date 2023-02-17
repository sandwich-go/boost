package xtime

import (
	"fmt"
	"github.com/sandwich-go/boost/internal/log"
	"github.com/sandwich-go/boost/xpanic"
	"math"
	"time"

	"github.com/sandwich-go/boost/xsync"
)

var (
	// CopToleranceSecond 时间容忍误差，误差超过CopToleranceSecond则后续使用系统时间
	CopToleranceSecond = int64(2)
	// CopToleranceCheckInterval 检测CopToleranceSecond的时间间隔
	CopToleranceCheckInterval = time.Duration(10) * time.Second
)

// Cop mock time like ruby's time cop
type Cop struct {
	ts          xsync.AtomicInt64
	unhealthy   chan struct{}
	closeChan   chan struct{}
	running     xsync.AtomicBool
	closeFlag   xsync.AtomicInt32
	nowProvider func() time.Time
}

// NewCop 新建Cop对象
func NewCop(nowProvider func() time.Time) *Cop {
	tc := &Cop{nowProvider: nowProvider}
	tc.closeChan = make(chan struct{})
	tc.unhealthy = make(chan struct{})
	tc.run()
	return tc
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
			systemTime := tc.nowProvider().Unix()
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
	tc.ts.Set(tc.nowProvider().Unix())
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
	return tc.nowProvider()
}

// Unix 获取当前Unix时间戳
func (tc *Cop) Unix() int64 { return tc.Now().Unix() }

// UnixMilli 获取当前时间，单位毫秒
func (tc *Cop) UnixMilli() int64 {
	return tc.Now().UnixNano() / int64(time.Millisecond)
}
