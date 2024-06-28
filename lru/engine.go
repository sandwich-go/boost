package lru

import (
	"github.com/sandwich-go/boost/xmath"
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
}

// NewEngine 创建 lru 类型的引擎
// interval ticker 间隔时间，过期时间需要 [interval*0.75, interval*1.25] 之前才准确
// locker 节点锁，保证添加/删除/过期元素并发安全
// expireHandler 过期元素处理器
func NewEngine[V any](interval time.Duration, locker sync.Locker, expireHandler func(V)) *Engine[V] {
	var e = &Engine[V]{
		cleanInterval: interval,
		expireHandler: expireHandler,
	}
	go func() {
		slept := min(time.Second, interval)
		for {
			time.Sleep(xmath.Disturb(slept, 25))
			e.expire(locker)
		}
	}()
	return e
}

// Add 添加元素
func (e *Engine[V]) Add(value V) *Node[V] {
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

func (e *Engine[V]) expire(locker sync.Locker) {
	if locker != nil {
		locker.Lock()
		defer locker.Unlock()
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
