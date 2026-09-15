//go:build boost_lru_list

package lru

import (
	"sync"
	"time"

	"github.com/rs/xid"
	"github.com/sandwich-go/boost/xsync"
)

// Node 是 Engine 的元素，同时是双向环形链表的结点，因此 Engine 的任何
// 变更操作都会改写它的 next/prev。
type Node[V any] struct {
	value      V
	accessTime time.Time
	next       *Node[V]
	prev       *Node[V]
	engine     *Engine[V]
}

// SetValue 设置节点上的值
func (n *Node[V]) SetValue(value V) { n.value = value }

// Value 获取节点上的值
func (n *Node[V]) Value() V {
	return n.value
}

// Engine 用双向环形链表维护 LRU 顺序，head 指向最近访问的元素，head.prev 是最久
// 未访问的元素。expire 从 head.prev 往前走，遇到第一个未超时的元素即停，因此单次清理
// 的成本与本轮实际过期的元素数成正比，与元素总数无关。
//
// 加锁契约：Add、Remove、Promote 三者都会改写链表结构，调用方必须全部在 locker 对应
// 的**写锁**下调用。Promote 尤其注意——它会把结点移到 head，不是只更新时间戳，在读锁
// 下并发调用会破坏链表，且没有任何运行时提示。需要在读锁下 Promote 时请改用 Engine。
type Engine[V any] struct {
	cleanInterval time.Duration
	head          *Node[V]
	expireHandler func(V)

	locker  sync.Locker
	started xsync.AtomicInt32
	active  bool
	id      string
}

func newEngine[V any](interval time.Duration, locker sync.Locker, expireHandler func(V), active bool) *Engine[V] {
	return &Engine[V]{
		cleanInterval: interval,
		locker:        locker,
		expireHandler: expireHandler,
		active:        active,
		id:            xid.New().String(),
	}
}

// NewLazyEngine 创建链表实现的 lru 引擎, 当第一个元素添加的时候，才会启动元素清理协程
// interval ticker 间隔时间，过期时间需要 [interval*0.9, interval*1.1] 之前才准确
// locker 节点锁，保证添加/删除/过期元素并发安全
// expireHandler 过期元素处理器
//
// Add、Remove、Promote 均需在 locker 对应的写锁下调用，详见 Engine 的说明。
func NewLazyEngine[V any](interval time.Duration, locker sync.Locker, expireHandler func(V)) *Engine[V] {
	return newEngine(interval, locker, expireHandler, false)
}

// NewEngine 创建链表实现的 lru 引擎
// interval ticker 间隔时间，过期时间需要 [interval*0.9, interval*1.1] 之前才准确
// locker 节点锁，保证添加/删除/过期元素并发安全
// expireHandler 过期元素处理器
//
// Add、Remove、Promote 均需在 locker 对应的写锁下调用，详见 Engine 的说明。
func NewEngine[V any](interval time.Duration, locker sync.Locker, expireHandler func(V)) *Engine[V] {
	var e = newEngine(interval, locker, expireHandler, true)
	e.startCleaner()
	return e
}

// ImplementationName 返回当前编译进二进制的 lru 实现名。lru 有两份实现，靠
// boost_lru_list build tag 二选一，且必须在所有构建路径上一致；服务启动时打印它可以
// 让 tag 配错在日志里一眼可见，而不是等行为异常了才发现。
func ImplementationName() string { return "linkedlist" }

func (e *Engine[V]) startCleaner() {
	if !e.started.CompareAndSwap(0, 1) {
		return
	}
	slept := min(time.Second, e.cleanInterval)
	DefaultCleanWorker.Clean(slept, e.id, e.expire)
}

// Add 添加元素。调用方须持有 locker 对应的写锁。
func (e *Engine[V]) Add(value V) *Node[V] {
	if !e.active && e.started.Get() == 0 {
		e.startCleaner()
	}
	n := &Node[V]{
		value:  value,
		engine: e,
	}
	e.Promote(n)
	return n
}

// Remove 删除元素。调用方须持有 locker 对应的写锁。
func (e *Engine[V]) Remove(node *Node[V]) {
	if node.engine != e {
		return
	}

	node.engine = nil

	if node == e.head {
		if node.next == node {
			e.head = nil
			return
		}

		e.head = node.next
	}

	node.next.prev = node.prev
	node.prev.next = node.next
}

// Promote 访问元素后，更新此元素的访问时间，并把它移到链表头部。
//
// 调用方须持有 locker 对应的**写锁**：本方法改写链表结构，不是只更新时间戳。
func (e *Engine[V]) Promote(node *Node[V]) bool {
	if node.engine != e {
		return false
	}

	node.accessTime = time.Now()
	e.moveNodeToHead(node)
	return true
}

func (e *Engine[V]) size() int {
	if e.head == nil {
		return 0
	}
	var i int
	tail := e.head.prev
	for tail != nil {
		i++
		if tail == e.head {
			tail = nil
		} else {
			tail = tail.prev
		}
	}
	return i
}

func (e *Engine[V]) moveNodeToHead(node *Node[V]) {
	if node == e.head || node.engine != e {
		return
	}

	if e.head == nil {
		e.head = node
		e.head.next = node
		e.head.prev = node
		return
	}

	if node.next != nil {
		node.next.prev = node.prev
		node.prev.next = node.next
	}

	node.next = e.head
	node.prev = e.head.prev
	e.head.prev.next = node
	e.head.prev = node
	e.head = node
}

func (e *Engine[V]) expire() {
	if e.locker != nil {
		e.locker.Lock()
		defer e.locker.Unlock()
	}
	if e.head == nil {
		return
	}

	now := time.Now()

	tail := e.head.prev
	for tail != nil && now.Sub(tail.accessTime) > e.cleanInterval {
		tail.engine = nil

		if e.expireHandler != nil {
			e.expireHandler(tail.value)
		}

		if tail == e.head {
			e.head = nil
			tail = nil
		} else {
			tail = tail.prev
		}
	}

	if tail != nil && tail.engine != nil {
		tail.next = e.head
		e.head.prev = tail
	}
}
