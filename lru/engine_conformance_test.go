package lru

import (
	"sort"
	"sync"
	"testing"
	"time"
)

// 本文件刻意不带 build tag：lru 有两份实现（惰性堆 / 链表），靠 boost_lru_list 二选一，
// 而两者对外可见的语义必须一致。同一份用例在两种配置下都跑，就是这份一致性的守卫——
// 比写一个「同时引用两种实现」的对比用例更实在，因为 build tag 下两者本就互斥。
//
// 这里只用两份实现都提供的入口：newEngine / Add / Remove / Promote / expire / size /
// Value / SetValue / ImplementationName。任何依赖内部字段的断言都属于某一份实现，应放到
// 带 tag 的测试文件里。

// sameIntSet 比较两个 int 切片的多重集合，不比较顺序。
//
// 过期回调的顺序不是对外承诺：链表实现按 tail->head 走，堆实现按 scheduledAt 出堆，
// 两者在「同时到期」时的先后没有保证，因此一致性用例只钉住集合。
func sameIntSet(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	x := append([]int(nil), a...)
	y := append([]int(nil), b...)
	sort.Ints(x)
	sort.Ints(y)
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func TestConformance_ImplementationNameIsKnown(t *testing.T) {
	switch got := ImplementationName(); got {
	case "heap", "linkedlist":
		t.Logf("当前编译进来的实现 = %s", got)
	default:
		t.Fatalf("ImplementationName() = %q，不是已知实现名", got)
	}
}

func TestConformance_AddRemoveSize(t *testing.T) {
	var lock sync.Mutex
	e := newEngine[int](time.Hour, &lock, nil, true)

	lock.Lock()
	defer lock.Unlock()

	if e.size() != 0 {
		t.Fatalf("空引擎 size = %d, want 0", e.size())
	}
	n1 := e.Add(1)
	n2 := e.Add(2)
	if e.size() != 2 {
		t.Fatalf("Add 两个后 size = %d, want 2", e.size())
	}
	e.Remove(n1)
	if e.size() != 1 {
		t.Fatalf("Remove 一个后 size = %d, want 1", e.size())
	}
	e.Remove(n1) // 重复删除必须是 no-op
	if e.size() != 1 {
		t.Fatalf("重复 Remove 后 size = %d, want 1", e.size())
	}
	e.Remove(n2)
	if e.size() != 0 {
		t.Fatalf("全部删除后 size = %d, want 0", e.size())
	}
}

func TestConformance_PromoteRejectsDetachedNode(t *testing.T) {
	var lock sync.Mutex
	e := newEngine[int](time.Hour, &lock, nil, true)

	lock.Lock()
	defer lock.Unlock()

	n := e.Add(1)
	if !e.Promote(n) {
		t.Fatal("Promote 引擎内的结点应返回 true")
	}
	e.Remove(n)
	if e.Promote(n) {
		t.Fatal("Promote 已脱离引擎的结点应返回 false")
	}
}

func TestConformance_ValueRoundTrip(t *testing.T) {
	var lock sync.Mutex
	e := newEngine[string](time.Hour, &lock, nil, true)

	lock.Lock()
	defer lock.Unlock()

	n := e.Add("a")
	if got := n.Value(); got != "a" {
		t.Fatalf("Value() = %q, want \"a\"", got)
	}
	n.SetValue("b")
	if got := n.Value(); got != "b" {
		t.Fatalf("SetValue 之后 Value() = %q, want \"b\"", got)
	}
}

// TestConformance_ExpireRemovesIdleEntries 空闲超过 interval 的元素必须被清掉。
//
// 时间余量：interval 200ms，等 400ms 才清理，需要 sleep 少于 200ms 才会误判，不可能。
func TestConformance_ExpireRemovesIdleEntries(t *testing.T) {
	const interval = 200 * time.Millisecond

	var lock sync.Mutex
	var expired []int
	e := newEngine(interval, &lock, func(v int) { expired = append(expired, v) }, true)

	lock.Lock()
	e.Add(1)
	e.Add(2)
	e.Add(3)
	lock.Unlock()

	time.Sleep(2 * interval)
	e.expire()

	if want := []int{1, 2, 3}; !sameIntSet(expired, want) {
		t.Fatalf("expireHandler 收到 %v, want 集合 %v", expired, want)
	}
	lock.Lock()
	defer lock.Unlock()
	if e.size() != 0 {
		t.Fatalf("全部过期后 size = %d, want 0", e.size())
	}
}

// TestConformance_PromoteDefersExpiry 访问过的元素必须重新获得完整的空闲期。
//
// 这是滑动 TTL 的核心语义，两份实现走的路径完全不同：链表把结点移到 head 并刷新
// accessTime；堆只原子刷新 lastAccess，等 cleaner 走到旧排期点时再校正。
//
// 时间余量：interval 500ms；t=400ms 访问，t=700ms 清理。t=700 > 500 保证堆实现真的
// 走到「取出结点、按 lastAccess 重新排期」这条路（否则它连堆顶都不会碰，等于没测到）；
// 空闲期 700-400=300ms，要让它被误判成过期，第二段 sleep 需要多睡 200ms 以上。
func TestConformance_PromoteDefersExpiry(t *testing.T) {
	const interval = 500 * time.Millisecond

	var lock sync.Mutex
	var expired []int
	e := newEngine(interval, &lock, func(v int) { expired = append(expired, v) }, true)

	lock.Lock()
	kept := e.Add(1)
	e.Add(2)
	lock.Unlock()

	time.Sleep(400 * time.Millisecond)
	lock.Lock()
	e.Promote(kept)
	lock.Unlock()

	time.Sleep(300 * time.Millisecond)
	e.expire()

	if want := []int{2}; !sameIntSet(expired, want) {
		t.Fatalf("expireHandler 收到 %v, want 集合 %v（被 Promote 的 1 不该过期）", expired, want)
	}
	lock.Lock()
	defer lock.Unlock()
	if e.size() != 1 {
		t.Fatalf("expire 后 size = %d, want 1", e.size())
	}
	if got := kept.Value(); got != 1 {
		t.Fatalf("存活结点的值 = %v, want 1", got)
	}
}
