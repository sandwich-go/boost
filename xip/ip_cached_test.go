package xip

import (
	"sync"
	"testing"
)

// resetLocalIPCache 把缓存与解析入口恢复成初始状态，并在用例结束后还原，
// 避免用例之间互相污染（包内用例共享 localIPCache / localIPResolver）。
func resetLocalIPCache(t *testing.T, resolver func() string) {
	t.Helper()
	prevResolver := localIPResolver
	localIPCache.Store(nil)
	localIPResolver = resolver
	t.Cleanup(func() {
		localIPCache.Store(nil)
		localIPResolver = prevResolver
	})
}

func TestGetLocalIPCached_SameValueAsGetLocalIP(t *testing.T) {
	resetLocalIPCache(t, GetLocalIP)

	want := GetLocalIP()
	if want == "" {
		t.Skip("本机解析不到非回环内网 IPv4，跳过取值一致性校验")
	}
	if got := GetLocalIPCached(); got != want {
		t.Fatalf("GetLocalIPCached() = %q, want %q（应与 GetLocalIP 同值）", got, want)
	}
}

func TestGetLocalIPCached_ResolveOnce(t *testing.T) {
	var calls int
	resetLocalIPCache(t, func() string {
		calls++
		return "10.0.0.1"
	})

	for i := 0; i < 100; i++ {
		if got := GetLocalIPCached(); got != "10.0.0.1" {
			t.Fatalf("第 %d 次调用返回 %q, want %q", i+1, got, "10.0.0.1")
		}
	}
	if calls != 1 {
		t.Fatalf("解析次数 = %d, want 1（后续调用必须命中缓存）", calls)
	}
}

// TestGetLocalIPCached_EmptyResultNotCached 覆盖解析失败分支：空串不能写进缓存，
// 否则一次瞬时失败会让进程终身返回空串。
func TestGetLocalIPCached_EmptyResultNotCached(t *testing.T) {
	var calls int
	ret := ""
	resetLocalIPCache(t, func() string {
		calls++
		return ret
	})

	if got := GetLocalIPCached(); got != "" {
		t.Fatalf("解析失败时返回 %q, want 空串", got)
	}
	if got := GetLocalIPCached(); got != "" {
		t.Fatalf("解析持续失败时返回 %q, want 空串", got)
	}
	if calls != 2 {
		t.Fatalf("解析次数 = %d, want 2（空结果不得写缓存）", calls)
	}

	// 解析恢复后应立刻生效并落缓存。
	ret = "10.0.0.2"
	if got := GetLocalIPCached(); got != "10.0.0.2" {
		t.Fatalf("解析恢复后返回 %q, want %q", got, "10.0.0.2")
	}
	if got := GetLocalIPCached(); got != "10.0.0.2" {
		t.Fatalf("缓存命中返回 %q, want %q", got, "10.0.0.2")
	}
	if calls != 3 {
		t.Fatalf("解析次数 = %d, want 3（成功后应停止解析）", calls)
	}
}

// TestGetLocalIPCached_Concurrent 在 -race 下校验读写缓存无数据竞争，且所有 goroutine
// 看到同一个值。resolver 在起 goroutine 前装好，运行期间不再改动。
func TestGetLocalIPCached_Concurrent(t *testing.T) {
	resetLocalIPCache(t, func() string { return "10.0.0.3" })

	const n = 64
	var wg sync.WaitGroup
	got := make([]string, n)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(idx int) {
			defer wg.Done()
			got[idx] = GetLocalIPCached()
		}(i)
	}
	wg.Wait()

	for i, v := range got {
		if v != "10.0.0.3" {
			t.Fatalf("goroutine %d 得到 %q, want %q", i, v, "10.0.0.3")
		}
	}
}

func BenchmarkGetLocalIP(b *testing.B) {
	for b.Loop() {
		_ = GetLocalIP()
	}
}

func BenchmarkGetLocalIPCached(b *testing.B) {
	localIPCache.Store(nil)
	b.Cleanup(func() { localIPCache.Store(nil) })
	if GetLocalIPCached() == "" {
		b.Skip("本机解析不到非回环内网 IPv4，缓存无法预热")
	}
	for b.Loop() {
		_ = GetLocalIPCached()
	}
}

func BenchmarkGetLocalIPCachedParallel(b *testing.B) {
	localIPCache.Store(nil)
	b.Cleanup(func() { localIPCache.Store(nil) })
	if GetLocalIPCached() == "" {
		b.Skip("本机解析不到非回环内网 IPv4，缓存无法预热")
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = GetLocalIPCached()
		}
	})
}
