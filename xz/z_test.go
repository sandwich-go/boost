package xz_test

import (
	"hash/fnv"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sandwich-go/boost/xtime"
	"github.com/sandwich-go/boost/xz"
)

// goos: darwin
// goarch: amd64
// pkg: github.com/sandwich-go/boost/xz
// cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
// BenchmarkFnv-12                 21295071                49.35 ns/op     1296.93 MB/s           0 B/op          0 allocs/op
// BenchmarkKeyToHash-12           28054087                39.32 ns/op     1627.77 MB/s          24 B/op          1 allocs/op
// BenchmarkMemHash-12             216110384                5.548 ns/op    11536.34 MB/s          0 B/op          0 allocs/op
// BenchmarkMemHashString-12       219418490                5.494 ns/op    11649.26 MB/s          0 B/op          0 allocs/op
func BenchmarkFnv(b *testing.B) {
	buf := make([]byte, 64)
	rand.Read(buf)
	f := fnv.New64a()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = f.Write(buf)
		f.Sum64()
		f.Reset()
	}
	b.SetBytes(int64(len(buf)))
}
func BenchmarkKeyToHash(b *testing.B) {
	buf := make([]byte, 64)
	rand.Read(buf)
	hashVal := uint64(0)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newVal := xz.KeyToHash(buf)
		if hashVal == 0 {
			hashVal = newVal
		}
		if newVal != hashVal {
			b.Fatal("got different hash val using same string")
		}
	}
	b.SetBytes(int64(len(buf)))
}

func BenchmarkMemHash(b *testing.B) {
	buf := make([]byte, 64)
	rand.Read(buf)
	hashVal := uint64(0)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newVal := xz.MemHash(buf)
		if hashVal == 0 {
			hashVal = newVal
		}
		if newVal != hashVal {
			b.Fatal("got different hash val using same string")
		}
	}
	b.SetBytes(int64(len(buf)))
}

func BenchmarkMemHashString(b *testing.B) {
	buf := make([]byte, 64)
	rand.Read(buf)
	s := string(buf)

	hashVal := uint64(0)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newVal := xz.MemHashString(s)
		if hashVal == 0 {
			hashVal = newVal
		}
		if newVal != hashVal {
			b.Fatal("got different hash val using same string")
		}
	}
	b.SetBytes(int64(len(s)))
}

// goos: darwin
// goarch: amd64
// pkg: github.com/sandwich-go/boost/xz
// cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
//BenchmarkFastRand-12         	1000000000	         0.3866 ns/op
//BenchmarkRandSource-12       	1000000000	         0.8944 ns/op
//BenchmarkRandGlobal-12       	22552677	        55.36 ns/op
//BenchmarkRandAtomic-12       	62376811	        19.32 ns/op
func benchmarkRand(b *testing.B, fab func() func() uint32) {
	b.RunParallel(func(pb *testing.PB) {
		gen := fab()
		for pb.Next() {
			gen()
		}
	})
}

func BenchmarkFastRand(b *testing.B) {
	benchmarkRand(b, func() func() uint32 {
		return xz.FastRand
	})
}

func BenchmarkRandSource(b *testing.B) {
	benchmarkRand(b, func() func() uint32 {
		s := rand.New(rand.NewSource(time.Now().Unix()))
		return func() uint32 { return s.Uint32() }
	})
}

func BenchmarkRandGlobal(b *testing.B) {
	benchmarkRand(b, func() func() uint32 {
		return func() uint32 { return rand.Uint32() }
	})
}

func BenchmarkRandAtomic(b *testing.B) {
	var x uint32
	benchmarkRand(b, func() func() uint32 {
		return func() uint32 { return uint32(atomic.AddUint32(&x, 1)) }
	})
}

// goos: darwin
// goarch: amd64
// pkg: github.com/sandwich-go/boost/xz
// cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
// BenchmarkMono-12                39144883                29.95 ns/op
// BenchmarkSTime-12               366019147                3.333 ns/op
// BenchmarkWall-12                28415515                41.18 ns/op
// BenchmarkGoTime-12              16638219                72.13 ns/op
func BenchmarkMono(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xz.MonoOffset()
	}
}

func BenchmarkSTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xtime.UnixMilli()
	}
}

func BenchmarkWall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		xz.Now()
	}
}

func BenchmarkGoTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Now().UnixNano()
	}
}
