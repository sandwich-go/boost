package lru

import (
	"github.com/sandwich-go/boost/xmath"
	"github.com/sandwich-go/boost/xsync"
	"sync"
	"time"
)

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

type Engine[V any] struct {
	cleanInterval time.Duration
	head          *Node[V]
	expireHandler func(V)

	locker  sync.Locker
	started xsync.AtomicInt32
	active  bool
}

func newEngine[V any](interval time.Duration, locker sync.Locker, expireHandler func(V), active bool) *Engine[V] {
	return &Engine[V]{
		cleanInterval: interval,
		locker:        locker,
		expireHandler: expireHandler,
		active:        active,
	}
}

// NewLazyEngine 创建 lru 类型的引擎, 当第一个元素添加的时候，才会启动元素清理协程
// interval ticker 间隔时间，过期时间需要 [interval*0.9, interval*1.1] 之前才准确
// locker 节点锁，保证添加/删除/过期元素并发安全
// expireHandler 过期元素处理器
func NewLazyEngine[V any](interval time.Duration, locker sync.Locker, expireHandler func(V)) *Engine[V] {
	return newEngine(interval, locker, expireHandler, false)
}

// NewEngine 创建 lru 类型的引擎
// interval ticker 间隔时间，过期时间需要 [interval*0.9, interval*1.1] 之前才准确
// locker 节点锁，保证添加/删除/过期元素并发安全
// expireHandler 过期元素处理器
func NewEngine[V any](interval time.Duration, locker sync.Locker, expireHandler func(V)) *Engine[V] {
	var e = newEngine(interval, locker, expireHandler, true)
	e.startCleaner()
	return e
}

func (e *Engine[V]) startCleaner() {
	if !e.started.CompareAndSwap(0, 1) {
		return
	}
	go func() {
		slept := min(time.Second, e.cleanInterval)
		for {
			time.Sleep(xmath.Disturb(slept, 10))
			e.expire()
		}
	}()
}

// Add 添加元素
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

// Remove 删除元素
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

// Promote 访问元素后，更新此元素的访问时间
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
