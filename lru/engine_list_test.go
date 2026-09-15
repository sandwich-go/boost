//go:build boost_lru_list

package lru

import (
	"sync"
	"testing"
	"time"
)

// ringValues 从 head 沿 next 走一圈，返回 MRU -> LRU 顺序的值。
// 环断裂时会因为走不回 head 而在超过 size 上限时报错，因此它同时充当链表完整性检查。
func ringValues[V any](t *testing.T, e *Engine[V], limit int) []V {
	t.Helper()
	var out []V
	if e.head == nil {
		return out
	}
	node := e.head
	for {
		out = append(out, node.value)
		node = node.next
		if node == e.head {
			return out
		}
		if node == nil {
			t.Fatalf("链表环断裂：next 走到 nil，已收集 %v", out)
		}
		if len(out) > limit {
			t.Fatalf("链表环断裂：走了 %d 步仍未回到 head，已收集 %v", len(out), out)
		}
	}
}

func TestEngine_AddPromoteKeepsMRUOrder(t *testing.T) {
	var lock sync.Mutex
	e := newEngine[int](time.Minute, &lock, nil, true)

	lock.Lock()
	n1 := e.Add(1)
	e.Add(2)
	e.Add(3)
	lock.Unlock()

	if got, want := ringValues(t, e, 8), []int{3, 2, 1}; !equalInts(got, want) {
		t.Fatalf("Add 之后 MRU->LRU 顺序 = %v, want %v", got, want)
	}
	if e.head.prev.value != 1 {
		t.Fatalf("tail 应是最久未访问的 1, got %v", e.head.prev.value)
	}

	lock.Lock()
	ok := e.Promote(n1)
	lock.Unlock()
	if !ok {
		t.Fatal("Promote 在引擎内的结点上应返回 true")
	}
	if got, want := ringValues(t, e, 8), []int{1, 3, 2}; !equalInts(got, want) {
		t.Fatalf("Promote(1) 之后 MRU->LRU 顺序 = %v, want %v", got, want)
	}
	if e.head.prev.value != 2 {
		t.Fatalf("Promote(1) 之后 tail 应是 2, got %v", e.head.prev.value)
	}
	if e.size() != 3 {
		t.Fatalf("size = %d, want 3", e.size())
	}
}

func TestEngine_ExpireRemovesIdleAndStopsAtLive(t *testing.T) {
	const interval = time.Minute
	var lock sync.Mutex
	var expired []int
	e := newEngine(interval, &lock, func(v int) { expired = append(expired, v) }, true)

	lock.Lock()
	n1 := e.Add(1)
	n2 := e.Add(2)
	n3 := e.Add(3)
	// 1 与 2 已空闲超过 interval，3 仍新鲜。链表顺序 head=3,2,1 与访问时间保持一致。
	old := time.Now().Add(-2 * interval)
	n1.accessTime = old
	n2.accessTime = old
	lock.Unlock()

	e.expire()

	// expire 从 tail 往前走，因此先 1 后 2。
	if got, want := expired, []int{1, 2}; !equalInts(got, want) {
		t.Fatalf("expireHandler 收到 %v, want %v（顺序应为 tail 优先）", got, want)
	}
	if n1.engine != nil || n2.engine != nil {
		t.Fatal("过期结点的 engine 应被置 nil")
	}
	if n3.engine != e {
		t.Fatal("未过期结点的 engine 不应被改动")
	}
	if got, want := ringValues(t, e, 8), []int{3}; !equalInts(got, want) {
		t.Fatalf("expire 后剩余 %v, want %v", got, want)
	}
	if e.head.next != e.head || e.head.prev != e.head {
		t.Fatal("只剩一个结点时应自成环")
	}
	if e.size() != 1 {
		t.Fatalf("size = %d, want 1", e.size())
	}
}

func TestEngine_ExpireAllEmptiesEngine(t *testing.T) {
	const interval = time.Minute
	var lock sync.Mutex
	var expired []int
	e := newEngine(interval, &lock, func(v int) { expired = append(expired, v) }, true)

	lock.Lock()
	nodes := []*Node[int]{e.Add(1), e.Add(2), e.Add(3)}
	old := time.Now().Add(-2 * interval)
	for _, n := range nodes {
		n.accessTime = old
	}
	lock.Unlock()

	e.expire()

	if got, want := expired, []int{1, 2, 3}; !equalInts(got, want) {
		t.Fatalf("expireHandler 收到 %v, want %v", got, want)
	}
	if e.head != nil {
		t.Fatal("全部过期后 head 应为 nil")
	}
	if e.size() != 0 {
		t.Fatalf("size = %d, want 0", e.size())
	}
	for _, n := range nodes {
		if n.engine != nil {
			t.Fatal("全部过期后每个结点的 engine 都应为 nil")
		}
	}
}

func TestEngine_PromoteDefersExpiry(t *testing.T) {
	const interval = time.Minute
	var lock sync.Mutex
	var expired []int
	e := newEngine(interval, &lock, func(v int) { expired = append(expired, v) }, true)

	lock.Lock()
	n1 := e.Add(1)
	n2 := e.Add(2)
	old := time.Now().Add(-2 * interval)
	n1.accessTime = old
	n2.accessTime = old
	// 访问 1：刷新访问时间并移到 head，它不该在本轮过期。
	e.Promote(n1)
	lock.Unlock()

	e.expire()

	if got, want := expired, []int{2}; !equalInts(got, want) {
		t.Fatalf("expireHandler 收到 %v, want %v（被 Promote 的 1 不该过期）", got, want)
	}
	if n1.engine != e {
		t.Fatal("被 Promote 的结点不该被过期")
	}
	if got, want := ringValues(t, e, 8), []int{1}; !equalInts(got, want) {
		t.Fatalf("expire 后剩余 %v, want %v", got, want)
	}
}

func TestEngine_RemoveSingleNode(t *testing.T) {
	var lock sync.Mutex
	e := newEngine[int](time.Minute, &lock, nil, true)

	lock.Lock()
	n1 := e.Add(1)
	e.Remove(n1)
	lock.Unlock()

	if e.head != nil {
		t.Fatal("删掉唯一结点后 head 应为 nil")
	}
	if e.size() != 0 {
		t.Fatalf("size = %d, want 0", e.size())
	}
	if n1.engine != nil {
		t.Fatal("被删结点的 engine 应为 nil")
	}
}

func TestEngine_RemoveMiddleKeepsRingIntact(t *testing.T) {
	var lock sync.Mutex
	e := newEngine[int](time.Minute, &lock, nil, true)

	lock.Lock()
	e.Add(1)
	n2 := e.Add(2)
	e.Add(3)
	e.Remove(n2)
	lock.Unlock()

	if got, want := ringValues(t, e, 8), []int{3, 1}; !equalInts(got, want) {
		t.Fatalf("删掉中间结点后剩余 %v, want %v", got, want)
	}
	if e.size() != 2 {
		t.Fatalf("size = %d, want 2", e.size())
	}
}

func TestEngine_RemoveHeadMovesHeadToNext(t *testing.T) {
	var lock sync.Mutex
	e := newEngine[int](time.Minute, &lock, nil, true)

	lock.Lock()
	e.Add(1)
	e.Add(2)
	n3 := e.Add(3)
	e.Remove(n3) // n3 是 head
	lock.Unlock()

	if e.head.value != 2 {
		t.Fatalf("删掉 head 后新 head = %v, want 2", e.head.value)
	}
	if got, want := ringValues(t, e, 8), []int{2, 1}; !equalInts(got, want) {
		t.Fatalf("删掉 head 后剩余 %v, want %v", got, want)
	}
}

func TestEngine_RemoveAndPromoteAreNoopForDetachedNode(t *testing.T) {
	var lock sync.Mutex
	e := newEngine[int](time.Minute, &lock, nil, true)

	lock.Lock()
	e.Add(1)
	n2 := e.Add(2)
	e.Add(3)
	e.Remove(n2)
	before := ringValues(t, e, 8)

	e.Remove(n2) // 重复删除
	if ok := e.Promote(n2); ok {
		t.Fatal("Promote 已脱离引擎的结点应返回 false")
	}
	lock.Unlock()

	if got := ringValues(t, e, 8); !equalInts(got, before) {
		t.Fatalf("对已脱离结点的重复操作改动了链表：%v -> %v", before, got)
	}
	if e.size() != 2 {
		t.Fatalf("size = %d, want 2", e.size())
	}
}

func TestEngine_ForeignNodeIsRejected(t *testing.T) {
	var lock sync.Mutex
	e1 := newEngine[int](time.Minute, &lock, nil, true)
	e2 := newEngine[int](time.Minute, &lock, nil, true)

	lock.Lock()
	foreign := e2.Add(99)
	e1.Add(1)
	before := ringValues(t, e1, 8)

	e1.Remove(foreign)
	if ok := e1.Promote(foreign); ok {
		t.Fatal("Promote 别的引擎的结点应返回 false")
	}
	lock.Unlock()

	if got := ringValues(t, e1, 8); !equalInts(got, before) {
		t.Fatalf("外来结点改动了链表：%v -> %v", before, got)
	}
	if foreign.engine != e2 {
		t.Fatal("外来结点的 engine 不该被改动")
	}
}

func TestEngine_LazyStartsCleanerOnFirstAdd(t *testing.T) {
	var lock sync.Mutex
	// interval 取足够大，后台 cleaner 在用例期间不会真正跑到。
	e := newEngine[int](time.Hour, &lock, nil, false)

	if e.started.Get() != 0 {
		t.Fatal("lazy 引擎在首次 Add 之前不应启动 cleaner")
	}
	lock.Lock()
	e.Add(1)
	lock.Unlock()
	if e.started.Get() != 1 {
		t.Fatal("lazy 引擎首次 Add 之后应启动 cleaner")
	}
}

func TestEngine_ActiveEngineDoesNotStartCleanerOnAdd(t *testing.T) {
	var lock sync.Mutex
	e := newEngine[int](time.Hour, &lock, nil, true)

	lock.Lock()
	e.Add(1)
	lock.Unlock()
	if e.started.Get() != 0 {
		t.Fatal("active 引擎的 cleaner 由构造函数启动，Add 不应重复启动")
	}
}

func TestEngine_NilLockerStillWorks(t *testing.T) {
	var expired []int
	const interval = time.Minute
	e := newEngine(interval, nil, func(v int) { expired = append(expired, v) }, true)

	n1 := e.Add(1)
	n1.accessTime = time.Now().Add(-2 * interval)
	e.Add(2)

	e.expire()

	if got, want := expired, []int{1}; !equalInts(got, want) {
		t.Fatalf("locker 为 nil 时 expireHandler 收到 %v, want %v", got, want)
	}
}

func equalInts[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestImplementationName_LinkedList(t *testing.T) {
	if got := ImplementationName(); got != "linkedlist" {
		t.Fatalf("ImplementationName() = %q, want \"linkedlist\"", got)
	}
}
