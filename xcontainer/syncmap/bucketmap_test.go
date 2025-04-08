package syncmap

import (
	"strconv"
	"testing"
	"time"
)

func generateTestData(n int) ([]string, []int) {
	keys := make([]string, n)
	values := make([]int, n)
	for i := 0; i < n; i++ {
		keys[i] = "key_" + strconv.Itoa(i)
		values[i] = i
	}
	return keys, values
}

func BenchmarkKeys_Cached(b *testing.B) {
	hashFn := func(s string) int64 {
		return int64(len(s)) // 简单 hash 函数
	}
	m, _, _ := NewBucketMapWithCacheKey[string, int](64, time.Millisecond*500, hashFn)
	keys, values := generateTestData(10000)

	for i := 0; i < len(keys); i++ {
		m.Store(keys[i], values[i])
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Keys()
	}
}

func BenchmarkKeys_Raw(b *testing.B) {
	hashFn := func(s string) int64 {
		return int64(len(s)) // 简单 hash 函数
	}
	m := NewBucketMap[string, int](64, hashFn)
	keys, values := generateTestData(10000)

	for i := 0; i < len(keys); i++ {
		m.Store(keys[i], values[i])
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = m.Keys()
	}
}
