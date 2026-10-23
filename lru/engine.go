//go:build boost_lru_heap

package lru

import (
	"container/heap"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/xid"
	"github.com/sandwich-go/boost/xsync"
)

type Node[V any] struct {
	value V

	lastAccess  atomic.Int64 //最后访问时间
	scheduledAt int64        //只有cleaner修改
	heapIndex   int
	engine      *Engine[V]
}

// SetValue 设置节点上的值
func (n *Node[V]) SetValue(value V) { n.value = value }

// Value 获取节点上的值
func (n *Node[V]) Value() V { return n.value }

type nodeHeap[V any] []*Node[V]

func (h *nodeHeap[V]) Len() int { return len(*h) }
func (h *nodeHeap[V]) Less(i, j int) bool {
	return (*h)[i].scheduledAt < (*h)[j].scheduledAt
}
func (h *nodeHeap[V]) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
	(*h)[i].heapIndex = i
	(*h)[j].heapIndex = j
}
func (h *nodeHeap[V]) Push(value any) {
	node := value.(*Node[V])
	node.heapIndex = len(*h)
	*h = append(*h, node)
}
func (h *nodeHeap[V]) Pop() any {
	old := *h
	last := len(old) - 1
	node := old[last]
	old[last] = nil
	node.heapIndex = -1
	*h = old[:last]
	return node
}

// Engine 使用惰性最小堆维护滑动 TTL。Promote 只原子更新访问时间，cleaner 到达
// 结点的旧过期点时再校正堆中的排期，因此读取方可以在共享读锁下调用 Promote。
type Engine[V any] struct {
	cleanInterval      time.Duration
	cleanIntervalNanos int64
	epoch              time.Time
	nodes              nodeHeap[V]
	expireHandler      func(V)

	locker  sync.Locker
	started xsync.AtomicInt32
	active  bool
	id      string
}

func newEngine[V any](interval time.Duration, locker sync.Locker, expireHandler func(V), active bool) *Engine[V] {
	return &Engine[V]{
		cleanInterval:      interval,
		cleanIntervalNanos: interval.Nanoseconds(),
		epoch:              time.Now(),
		locker:             locker,
		expireHandler:      expireHandler,
		active:             active,
		id:                 xid.New().String(),
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
func NewEngine[V any](interval time.Duration, locker sync.Locker, expireHandler func(V), opts ...EngineOption) *Engine[V] {
	var e = newEngine(interval, locker, expireHandler, true)
	if autoCleanEnabled(opts) {
		e.startCleaner()
	}
	return e
}

// PromoteNeedsWriteLock 声明 Promote 是否要求调用方持有写锁。
const PromoteNeedsWriteLock = false

// ImplementationName 返回当前编译的实现名：boost_lru_heap 选择堆，默认使用链表。
func ImplementationName() string { return "heap" }

func (e *Engine[V]) now() int64 {
	return time.Since(e.epoch).Nanoseconds()
}

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
	now := e.now()
	node := &Node[V]{
		value:       value,
		scheduledAt: now,
		heapIndex:   -1,
		engine:      e,
	}
	node.lastAccess.Store(now)
	heap.Push(&e.nodes, node)
	return node
}

// Remove 删除元素。调用方须持有 locker 对应的写锁。
func (e *Engine[V]) Remove(node *Node[V]) {
	if node.engine != e {
		return
	}
	heap.Remove(&e.nodes, node.heapIndex)
	node.engine = nil
}

// Promote 更新元素的最后访问时间。它不调整堆，可在 locker 对应的读锁下调用。
func (e *Engine[V]) Promote(node *Node[V]) bool {
	if node.engine != e {
		return false
	}
	now := e.now()
	e.promoteAt(node, now)
	return true
}

// IsExpired 只判断逻辑过期，不删除节点；调用方须持有保护节点的读锁或写锁。
func (e *Engine[V]) IsExpired(node *Node[V]) bool {
	return node == nil || node.engine != e || e.now()-node.lastAccess.Load() > e.cleanIntervalNanos
}

// TryPromote 只为尚未过期的节点续期，避免分批淘汰期间读操作复活旧数据。
// 与 Promote 一样，可在保护节点的读锁下调用。
func (e *Engine[V]) TryPromote(node *Node[V]) bool {
	if node == nil || node.engine != e {
		return false
	}
	now := e.now()
	if now-node.lastAccess.Load() > e.cleanIntervalNanos {
		return false
	}
	e.promoteAt(node, now)
	return true
}

func (e *Engine[V]) promoteAt(node *Node[V], now int64) {
	for lastAccess := node.lastAccess.Load(); now > lastAccess; lastAccess = node.lastAccess.Load() {
		if node.lastAccess.CompareAndSwap(lastAccess, now) {
			break
		}
	}
}

func (e *Engine[V]) size() int { return e.nodes.Len() }

// expiryBatchThreshold 避免少量到期结点触发全堆扫描；达到该数量后，继续逐个
// heap.Pop/Push 的成本已经足以摊平一次线性整理。
const expiryBatchThreshold = 256

// expireLimited 逐个校正堆顶；正数预算下不能触发全堆线性整理，否则会突破锁内预算。
func (e *Engine[V]) expireLimited(limit int) bool {
	if e.locker != nil {
		e.locker.Lock()
		defer e.locker.Unlock()
	}
	now := e.now()
	for processed := 0; e.nodes.Len() > 0; processed++ {
		node := e.nodes[0]
		if now-node.scheduledAt <= e.cleanIntervalNanos {
			return false
		}
		if processed == limit {
			return true
		}
		lastAccess := node.lastAccess.Load()
		if now-lastAccess <= e.cleanIntervalNanos {
			node.scheduledAt = lastAccess
			heap.Fix(&e.nodes, 0)
			continue
		}
		heap.Pop(&e.nodes)
		e.expireNode(node)
	}
	return false
}

func (e *Engine[V]) expire() {
	if e.locker != nil {
		e.locker.Lock()
		defer e.locker.Unlock()
	}

	now := e.now()
	rescheduled := make([]*Node[V], 0, expiryBatchThreshold)
	processed := 0
	for e.nodes.Len() > 0 {
		node := e.nodes[0]
		if now-node.scheduledAt <= e.cleanIntervalNanos {
			break
		}

		heap.Pop(&e.nodes)
		lastAccess := node.lastAccess.Load()
		if now-lastAccess <= e.cleanIntervalNanos {
			node.scheduledAt = lastAccess
			rescheduled = append(rescheduled, node)
		} else {
			e.expireNode(node)
		}

		processed++
		if processed == expiryBatchThreshold {
			e.expireBatch(now, rescheduled)
			return
		}
	}

	for _, node := range rescheduled {
		heap.Push(&e.nodes, node)
	}
}

// expireBatch 在线性扫描中完成同一批旧排期的校正，最后只建堆一次。
func (e *Engine[V]) expireBatch(now int64, rescheduled []*Node[V]) {
	survivors := e.nodes[:0]
	for _, node := range e.nodes {
		if now-node.scheduledAt > e.cleanIntervalNanos {
			lastAccess := node.lastAccess.Load()
			if now-lastAccess <= e.cleanIntervalNanos {
				node.scheduledAt = lastAccess
			} else {
				e.expireNode(node)
				continue
			}
		}
		survivors = append(survivors, node)
	}

	survivors = append(survivors, rescheduled...)
	e.nodes = survivors
	for index, node := range e.nodes {
		node.heapIndex = index
	}
	heap.Init(&e.nodes)
}

func (e *Engine[V]) expireNode(node *Node[V]) {
	node.heapIndex = -1
	node.engine = nil
	if e.expireHandler != nil {
		e.expireHandler(node.value)
	}
}
