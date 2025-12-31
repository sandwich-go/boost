package protobuf

import (
	"context"
	"fmt"
	"testing"

	"github.com/sandwich-go/boost/xencoding"
	"github.com/sandwich-go/boost/xencoding/protobuf/test_perf"

	"google.golang.org/protobuf/proto"
)

func setupBenchmarkProtoCodecInputs(payloadBaseSize uint32) []proto.Message {
	payloadBase := make([]byte, payloadBaseSize)
	// arbitrary byte slices
	payloadSuffixes := [][]byte{
		[]byte("one"),
		[]byte("two"),
		[]byte("three"),
		[]byte("four"),
		[]byte("five"),
	}
	protoStructs := make([]proto.Message, 0)

	for _, p := range payloadSuffixes {
		ps := &test_perf.Buffer{}
		ps.Body = append(payloadBase, p...)
		protoStructs = append(protoStructs, ps)
	}

	return protoStructs
}

// The possible use of certain protobuf APIs like the proto.Buffer API potentially involves caching
// on our side. This can add checks around memory allocations and possible contention.
// Example run: go test -v -run=^$ -bench=BenchmarkProtoCodec -benchmem
func BenchmarkProtoCodec(b *testing.B) {
	// range of message sizes
	payloadBaseSizes := make([]uint32, 0)
	for i := uint32(0); i <= 12; i += 4 {
		payloadBaseSizes = append(payloadBaseSizes, 1<<i)
	}
	// range of SetParallelism
	parallelisms := make([]int, 0)
	for i := uint32(0); i <= 16; i += 4 {
		parallelisms = append(parallelisms, int(1<<i))
	}
	for _, s := range payloadBaseSizes {
		for _, p := range parallelisms {
			protoStructs := setupBenchmarkProtoCodecInputs(s)
			name := fmt.Sprintf("MinPayloadSize:%v/SetParallelism(%v)", s, p)
			b.Run(name, func(b *testing.B) {
				codec := &codec{}
				b.SetParallelism(p)
				b.RunParallel(func(pb *testing.PB) {
					benchmarkProtoCodec(codec, protoStructs, pb, b)
				})
			})
		}
	}
}

func benchmarkProtoCodec(codec *codec, protoStructs []proto.Message, pb *testing.PB, b *testing.B) {
	counter := 0
	for pb.Next() {
		counter++
		ps := protoStructs[counter%len(protoStructs)]
		fastMarshalAndUnmarshal(codec, ps, b)
	}
}

func fastMarshalAndUnmarshal(codec xencoding.Codec, protoStruct proto.Message, b *testing.B) {
	marshaledBytes, err := codec.Marshal(context.Background(), protoStruct)
	if err != nil {
		b.Errorf("codec.Marshal(_) returned an error")
	}
	res := test_perf.Buffer{}
	if err := codec.Unmarshal(context.Background(), marshaledBytes, &res); err != nil {
		b.Errorf("codec.Unmarshal(_) returned an error")
	}
}

// BenchmarkVTProtoCodec 测试 vtproto 的性能
func BenchmarkVTProtoCodec(b *testing.B) {
	payloadBaseSizes := []uint32{1, 16, 256, 4096}

	for _, s := range payloadBaseSizes {
		protoStructs := setupBenchmarkProtoCodecInputs(s)

		b.Run(fmt.Sprintf("Standard/Size:%v", s), func(b *testing.B) {
			codec := &codec{usingPool: false, name: CodecName}
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				ps := protoStructs[i%len(protoStructs)]
				fastMarshalAndUnmarshal(codec, ps, b)
			}
		})

		b.Run(fmt.Sprintf("WithPool/Size:%v", s), func(b *testing.B) {
			codec := &codec{usingPool: true, name: UsingPoolCodecName}
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				ps := protoStructs[i%len(protoStructs)]
				fastMarshalAndUnmarshal(codec, ps, b)
			}
		})
	}
}

// BenchmarkVTProtoMarshal 单独测试序列化性能
func BenchmarkVTProtoMarshal(b *testing.B) {
	msg := &test_perf.Buffer{Body: make([]byte, 1024)}

	b.Run("VT Standard", func(b *testing.B) {
		codec := &codec{usingPool: false, name: CodecName}
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := codec.Marshal(context.Background(), msg)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("VT WithPool", func(b *testing.B) {
		codec := &codec{usingPool: true, name: UsingPoolCodecName}
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := codec.Marshal(context.Background(), msg)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("DirectVT", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := msg.MarshalVT()
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("StandardProto", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := proto.Marshal(msg)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("StandardProtoWithPool", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := standardMarshalWithPool(msg)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkVTProtoUnmarshal 单独测试反序列化性能
func BenchmarkVTProtoUnmarshal(b *testing.B) {
	msg := &test_perf.Buffer{Body: make([]byte, 1024)}
	data, _ := msg.MarshalVT()

	b.Run("Standard", func(b *testing.B) {
		codec := &codec{usingPool: false, name: CodecName}
		result := &test_perf.Buffer{}
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			err := codec.Unmarshal(context.Background(), data, result)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("DirectVT", func(b *testing.B) {
		result := &test_perf.Buffer{}
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			err := result.UnmarshalVT(data)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("StandardProto", func(b *testing.B) {
		result := &test_perf.Buffer{}
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			err := proto.Unmarshal(data, result)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkMarshalParallel 并发场景下的性能对比
func BenchmarkMarshalParallel(b *testing.B) {
	msg := &test_perf.Buffer{Body: make([]byte, 1024)}

	b.Run("StandardProto", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := proto.Marshal(msg)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	})

	b.Run("StandardProtoWithPool", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := standardMarshalWithPool(msg)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}
