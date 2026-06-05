package lru

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：补 lru 包内未覆盖的导出 API + 双向链表关键路径：
//   - Node.SetValue / Node.Value（之前 0%）
//   - NewLazyEngine（之前 0%）
//   - moveNodeToHead 各分支（之前 46.7%，LRU 核心移动逻辑）
//   - Add / Remove / Promote 边界（误注入其他 engine 的 node）
//   - expire 过期淘汰
//
// 已有 engine_test.go 覆盖主流程；本文件补漏 + 边界。

// TestNode_SetValueAndValue 覆盖 Node.SetValue / Value setter/getter。
func TestNode_SetValueAndValue(t *testing.T) {
	Convey("Node.SetValue / Value", t, func() {
		mx := sync.Mutex{}
		e := NewEngine[string](100*time.Millisecond, &mx, nil)

		mx.Lock()
		n := e.Add("initial")
		mx.Unlock()

		So(n.Value(), ShouldEqual, "initial")

		n.SetValue("updated")
		So(n.Value(), ShouldEqual, "updated")
	})
}

// TestNewLazyEngine 覆盖 NewLazyEngine：第一次 Add 才启动 cleaner。
func TestNewLazyEngine(t *testing.T) {
	Convey("NewLazyEngine 在第一次 Add 前 cleaner 不启动", t, func() {
		mx := sync.Mutex{}
		var expired atomic.Int32
		handler := func(string) { expired.Add(1) }

		ttl := 30 * time.Millisecond
		e := NewLazyEngine[string](ttl, &mx, handler)

		// 还没 Add：started 应该是 0
		So(e.started.Get(), ShouldEqual, int32(0))

		mx.Lock()
		n := e.Add("v1")
		mx.Unlock()

		// Add 后 cleaner 启动
		So(e.started.Get(), ShouldEqual, int32(1))
		So(n.Value(), ShouldEqual, "v1")

		// 等过期触发：interval=30ms + sleep 100ms，确保 expire 跑过
		time.Sleep(100 * time.Millisecond)
		So(expired.Load(), ShouldBeGreaterThanOrEqualTo, int32(1))
	})

	Convey("NewLazyEngine 第二次 startCleaner 走 CompareAndSwap 失败早返", t, func() {
		mx := sync.Mutex{}
		e := NewLazyEngine[string](100*time.Millisecond, &mx, nil)

		// 第一次手动调 startCleaner
		e.startCleaner()
		So(e.started.Get(), ShouldEqual, int32(1))

		// 再调一次：CompareAndSwap(0, 1) 失败，早返不重复注册
		e.startCleaner()
		So(e.started.Get(), ShouldEqual, int32(1)) // 仍然是 1，没变成 2
	})
}

// TestEngine_MoveNodeToHead_AllBranches 覆盖 moveNodeToHead 4 个分支：
//   - node == e.head 早返
//   - e.head == nil 初始化
//   - node 已有 next/prev 时摘链 + 插头
//   - 多次 Promote 验证链表完整性
func TestEngine_MoveNodeToHead_AllBranches(t *testing.T) {
	Convey("Promote 同一 node 多次：第二次 node==head 早返", t, func() {
		mx := sync.Mutex{}
		e := NewEngine[string](100*time.Millisecond, &mx, nil)

		mx.Lock()
		defer mx.Unlock()

		n := e.Add("a")
		// 此时 n == e.head；再 Promote 同一个 node 触发 head 早返
		ok := e.Promote(n)
		So(ok, ShouldBeTrue)
		So(e.size(), ShouldEqual, 1)
	})

	Convey("Promote 中间节点：摘链后插到头", t, func() {
		mx := sync.Mutex{}
		e := NewEngine[string](100*time.Millisecond, &mx, nil)

		mx.Lock()
		defer mx.Unlock()

		n1 := e.Add("a")
		n2 := e.Add("b")
		n3 := e.Add("c") // head 顺序：c -> b -> a -> c (循环)

		So(e.head, ShouldEqual, n3)
		So(e.size(), ShouldEqual, 3)

		// 把 n1（最旧）promote 到 head
		So(e.Promote(n1), ShouldBeTrue)
		So(e.head, ShouldEqual, n1)
		So(e.size(), ShouldEqual, 3)

		// 验证链表完整：n1.next 应为 n3（原 head），n1.prev 应为 n2（原 tail）
		So(n1.next, ShouldEqual, n3)
		So(n1.prev, ShouldEqual, n2)
		_ = n2
	})

	Convey("Promote 来自其它 engine 的 node 返回 false 不动作", t, func() {
		mx := sync.Mutex{}
		e1 := NewEngine[string](100*time.Millisecond, &mx, nil)
		e2 := NewEngine[string](100*time.Millisecond, &mx, nil)

		mx.Lock()
		defer mx.Unlock()

		n := e1.Add("from-e1")
		// 拿 e1 的 node 去 Promote 到 e2
		ok := e2.Promote(n)
		So(ok, ShouldBeFalse)
		// e2 仍为空
		So(e2.size(), ShouldEqual, 0)
		// e1 不受影响
		So(e1.size(), ShouldEqual, 1)
	})
}

// TestEngine_Remove_AllBranches 覆盖 Remove 各分支：
//   - 来自其它 engine 的 node 早返
//   - 删 head 且只剩 1 个：head=nil
//   - 删 head 且还有其它：head=node.next
//   - 删非 head：摘链
func TestEngine_Remove_AllBranches(t *testing.T) {
	mx := sync.Mutex{}
	Convey("Remove 来自其它 engine 的 node 早返不动作", t, func() {
		e1 := NewEngine[string](100*time.Millisecond, &mx, nil)
		e2 := NewEngine[string](100*time.Millisecond, &mx, nil)

		mx.Lock()
		defer mx.Unlock()

		n := e1.Add("e1-node")
		e2.Remove(n) // node.engine != e2，早返
		So(e1.size(), ShouldEqual, 1)
	})

	Convey("Remove 唯一 head：head 变 nil", t, func() {
		e := NewEngine[int](100*time.Millisecond, &mx, nil)
		mx.Lock()
		defer mx.Unlock()

		n := e.Add(1)
		So(e.head, ShouldEqual, n)
		e.Remove(n)
		So(e.head, ShouldBeNil)
		So(e.size(), ShouldEqual, 0)
	})

	Convey("Remove head 但还有其它：head=next", t, func() {
		e := NewEngine[int](100*time.Millisecond, &mx, nil)
		mx.Lock()
		defer mx.Unlock()

		n1 := e.Add(1)
		n2 := e.Add(2) // head=n2; n2.next=n1
		So(e.head, ShouldEqual, n2)

		e.Remove(n2)
		// n2 是 head；删除后 head 应该变成 n2.next（在循环链表中是 n1）
		So(e.head, ShouldEqual, n1)
		So(e.size(), ShouldEqual, 1)
	})

	Convey("Remove 非 head 节点：摘链不影响 head", t, func() {
		e := NewEngine[int](100*time.Millisecond, &mx, nil)
		mx.Lock()
		defer mx.Unlock()

		n1 := e.Add(1) // 最旧
		_ = e.Add(2)
		n3 := e.Add(3) // 最新（head）

		So(e.head, ShouldEqual, n3)
		e.Remove(n1) // 删最旧的（非 head）
		So(e.head, ShouldEqual, n3)
		So(e.size(), ShouldEqual, 2)
	})
}

// TestEngine_Expire_HandlerNil 覆盖 expire 中 expireHandler == nil 的早返
// 分支（之前 88.9%）。
func TestEngine_Expire_HandlerNil(t *testing.T) {
	Convey("expire handler==nil 仍然清理过期 entry", t, func() {
		mx := sync.Mutex{}
		// 不传 expireHandler
		e := NewEngine[int](20*time.Millisecond, &mx, nil)

		mx.Lock()
		_ = e.Add(1)
		_ = e.Add(2)
		mx.Unlock()

		// 等过期
		time.Sleep(200 * time.Millisecond)

		mx.Lock()
		So(e.size(), ShouldEqual, 0)
		mx.Unlock()
	})

	Convey("expire 在空链表上立即返回", t, func() {
		mx := sync.Mutex{}
		e := NewEngine[int](20*time.Millisecond, &mx, nil)

		// 不 Add 任何元素，直接调 expire（expire 内部会 Lock，所以这里不
		// 能持锁调用——会 self-deadlock）
		e.expire()

		mx.Lock()
		So(e.head, ShouldBeNil)
		mx.Unlock()
	})

	Convey("expire 不带 locker 也不 panic", t, func() {
		// locker = nil 路径：if e.locker != nil 早返不加锁
		// NewEngine 不允许 nil locker（startCleaner 后调用 expire 会 Lock）；
		// 但实现里 expire 有 `if e.locker != nil` 守卫。我们用 NewLazyEngine
		// 避开 startCleaner 启动后台 goroutine 立即 expire 的死锁场景。
		e := NewLazyEngine[int](20*time.Millisecond, nil, nil)
		// 直接调 expire，验证 nil locker 不 panic
		e.expire()
		// 不 Add 任何元素，head 应该 nil
		So(e.head, ShouldBeNil)
	})
}
