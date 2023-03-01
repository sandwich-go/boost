package protobuf

import (
	"context"
	"fmt"
	"github.com/sandwich-go/boost/xencoding"
	"github.com/sandwich-go/boost/xencoding/protobuf/test_perf"
	"testing"

	"github.com/golang/protobuf/proto"
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
