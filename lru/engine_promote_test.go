package lru

import (
	"container/heap"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPromoteDefersExpiry(t *testing.T) {
	Convey("到达旧过期点时按最后访问时间重新排期", t, func() {
		var lock sync.Mutex
		var expired []int
		engine := newEngine(time.Minute, &lock, func(value int) {
			expired = append(expired, value)
		}, true)

		lock.Lock()
		expiredNode := engine.Add(1)
		refreshedNode := engine.Add(2)
		now := engine.now()
		expiredAt := now - 2*time.Minute.Nanoseconds()
		expiredNode.scheduledAt = expiredAt
		expiredNode.lastAccess.Store(expiredAt)
		refreshedNode.scheduledAt = now - 3*time.Minute.Nanoseconds()
		refreshedNode.lastAccess.Store(now)
		heap.Init(&engine.nodes)
		lock.Unlock()

		engine.expire()
		So(expired, ShouldResemble, []int{1})
		So(engine.size(), ShouldEqual, 1)
		So(refreshedNode.engine, ShouldEqual, engine)
	})
}

func TestPromoteUnderReadLock(t *testing.T) {
	var lock sync.RWMutex
	engine := newEngine[int](time.Hour, &lock, nil, true)
	lock.Lock()
	node := engine.Add(1)
	lock.Unlock()

	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				lock.RLock()
				if !engine.Promote(node) {
					t.Error("Promote failed")
				}
				lock.RUnlock()
			}
		}()
	}
	wg.Wait()
}

func TestExpireBatch(t *testing.T) {
	Convey("大量同批到期结点线性整理后保持堆序和滑动 TTL", t, func() {
		var lock sync.Mutex
		expired := make(map[int]bool)
		engine := newEngine(time.Minute, &lock, func(value int) {
			expired[value] = true
		}, true)

		lock.Lock()
		now := engine.now()
		old := now - 2*time.Minute.Nanoseconds()
		const count = expiryBatchThreshold * 4
		for i := 0; i < count; i++ {
			node := engine.Add(i)
			node.scheduledAt = old
			if i%2 == 0 {
				node.lastAccess.Store(old)
			} else {
				node.lastAccess.Store(now)
			}
		}
		future := engine.Add(count)
		heap.Init(&engine.nodes)
		lock.Unlock()

		engine.expire()
		So(len(expired), ShouldEqual, count/2)
		So(engine.size(), ShouldEqual, count/2+1)
		So(future.engine, ShouldEqual, engine)
		for index, node := range engine.nodes {
			if node.heapIndex != index {
				t.Fatalf("node %d heapIndex = %d", index, node.heapIndex)
			}
			left := index*2 + 1
			right := left + 1
			if left < len(engine.nodes) {
				if node.scheduledAt > engine.nodes[left].scheduledAt {
					t.Fatalf("heap order broken at %d -> %d", index, left)
				}
			}
			if right < len(engine.nodes) {
				if node.scheduledAt > engine.nodes[right].scheduledAt {
					t.Fatalf("heap order broken at %d -> %d", index, right)
				}
			}
		}
	})
}
