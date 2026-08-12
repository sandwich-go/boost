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

func TestFiveSecondTTLExpiryAndHeapReschedule(t *testing.T) {
	const ttl = 5 * time.Second

	var lock sync.RWMutex
	store := make(map[string]bool)
	expired := make(chan string, 2)
	engine := NewEngine(ttl, &lock, func(key string) {
		delete(store, key)
		expired <- key
	})

	lock.Lock()
	store["cold"] = true
	cold := engine.Add("cold")
	store["hot"] = true
	hot := engine.Add("hot")
	lock.Unlock()

	// 在初始 TTL 到达前访问 hot。它仍留在旧堆位置，只更新时间戳。
	time.Sleep(3 * time.Second)
	lock.RLock()
	oldSchedule := hot.scheduledAt
	if !engine.Promote(hot) {
		t.Fatal("hot promote failed")
	}
	promotedAt := hot.lastAccess.Load()
	if hot.scheduledAt != oldSchedule {
		t.Fatal("Promote 不应在读路径调整堆")
	}
	lock.RUnlock()

	// cold 应按首次写入时间淘汰；hot 到达旧过期点后应按最后访问时间重新排堆。
	select {
	case key := <-expired:
		if key != "cold" {
			t.Fatalf("先淘汰了 %q，期望 cold", key)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("cold 未在 5 秒 TTL 后被淘汰")
	}

	rescheduledDeadline := time.Now().Add(2 * time.Second)
	for {
		lock.RLock()
		rescheduled := hot.scheduledAt == promotedAt
		lock.RUnlock()
		if rescheduled {
			break
		}
		if time.Now().After(rescheduledDeadline) {
			t.Fatal("hot 到达旧过期点后未重新排堆")
		}
		time.Sleep(10 * time.Millisecond)
	}

	lock.RLock()
	if store["cold"] || cold.engine != nil {
		t.Error("cold key 或结点未被完整淘汰")
	}
	if !store["hot"] || hot.engine != engine {
		t.Error("hot key 在续期后不应被淘汰")
	}
	if engine.nodes.Len() != 1 || engine.nodes[0] != hot || hot.heapIndex != 0 {
		t.Error("hot 重新排堆后的根结点或 heapIndex 不正确")
	}
	lock.RUnlock()

	// hot 必须从最后一次访问起再存活完整 TTL，随后才淘汰。
	select {
	case key := <-expired:
		if key != "hot" {
			t.Fatalf("最终淘汰了 %q，期望 hot", key)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("hot 未在最后访问时间起 5 秒后被淘汰")
	}

	lock.RLock()
	defer lock.RUnlock()
	if store["hot"] || hot.engine != nil || engine.nodes.Len() != 0 || hot.heapIndex != -1 {
		t.Fatal("hot 淘汰后 store 或堆中仍有残留")
	}
}
