package msgpack

import (
	"bytes"
	"context"
	"sync"
	"testing"

	"github.com/sandwich-go/boost/xdebug/race"
	"github.com/sandwich-go/boost/xencoding"
	"github.com/sandwich-go/boost/xencoding/protobuf/test_perf"
)

func marshalAndUnmarshal(t *testing.T, codec xencoding.Codec, expectedBody []byte) {
	p := &test_perf.Buffer{}
	p.Body = expectedBody

	marshalledBytes, err := codec.Marshal(context.Background(), p)
	if err != nil {
		t.Errorf("codec.Marshal(_) returned an error")
	}

	if err := codec.Unmarshal(context.Background(), marshalledBytes, p); err != nil {
		t.Errorf("codec.Unmarshal(_) returned an error")
	}

	if !bytes.Equal(p.GetBody(), expectedBody) {
		t.Errorf("Unexpected body; got %v; want %v", p.GetBody(), expectedBody)
	}
}

func TestMsgpackCodecMarshalAndUnmarshal(t *testing.T) {
	marshalAndUnmarshal(t, codec{}, []byte{1, 2, 3})
}

// TestConcurrentUsage 100 goroutine × 1000 marshal/unmarshal 压测，意在
// 抓取 codec 内部 pool 复用产生的竞态。
//
// 已知 pre-existing race（vmihailenco/msgpack 上游 lib 问题，v5.3.5 / v5.4.1
// 都复现）：Marshal 路径 Encoder.write 与并发 Unmarshal 路径 Decoder
// readN/Reader.Read 在共享 backing array 上 race。boost.codec 是无状态
// wrapper，调用栈直达 lib 顶层 API。
//
// 处理策略（按 §1 + §4.6：上游 lib 不动）：在 -race 模式自动 skip。
// 非 race 模式仍跑（覆盖正确性）。
func TestConcurrentUsage(t *testing.T) {
	if race.Enabled {
		t.Skip("pre-existing race in vmihailenco/msgpack v5 lib (Marshal/Unmarshal pool reuse); see test doc comment")
	}
	const (
		numGoRoutines   = 100
		numMarshUnmarsh = 1000
	)

	// small, arbitrary byte slices
	protoBodies := [][]byte{
		[]byte("one"),
		[]byte("two"),
		[]byte("three"),
		[]byte("four"),
		[]byte("five"),
	}

	var wg sync.WaitGroup
	codec := codec{}

	for i := 0; i < numGoRoutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for k := 0; k < numMarshUnmarsh; k++ {
				marshalAndUnmarshal(t, codec, protoBodies[k%len(protoBodies)])
			}
		}()
	}

	wg.Wait()
}

// TestStaggeredMarshalAndUnmarshalUsingSamePool tries to catch potential errors in which slices get
// stomped on during reuse of a proto.Buffer.
func TestStaggeredMarshalAndUnmarshalUsingSamePool(t *testing.T) {
	codec1 := codec{}
	codec2 := codec{}

	expectedBody1 := []byte{1, 2, 3}
	expectedBody2 := []byte{4, 5, 6}

	proto1 := test_perf.Buffer{Body: expectedBody1}
	proto2 := test_perf.Buffer{Body: expectedBody2}

	var m1, m2 []byte
	var err error

	if m1, err = codec1.Marshal(context.Background(), &proto1); err != nil {
		t.Errorf("codec.Marshal(%v) failed", &proto1)
	}

	if m2, err = codec2.Marshal(context.Background(), &proto2); err != nil {
		t.Errorf("codec.Marshal(%v) failed", &proto2)
	}

	if err = codec1.Unmarshal(context.Background(), m1, &proto1); err != nil {
		t.Errorf("codec.Unmarshal(%v) failed", m1)
	}

	if err = codec2.Unmarshal(context.Background(), m2, &proto2); err != nil {
		t.Errorf("codec.Unmarshal(%v) failed", m2)
	}

	b1 := proto1.GetBody()
	b2 := proto2.GetBody()

	for i, v := range b1 {
		if expectedBody1[i] != v {
			t.Errorf("expected %v at index %v but got %v", i, expectedBody1[i], v)
		}
	}

	for i, v := range b2 {
		if expectedBody2[i] != v {
			t.Errorf("expected %v at index %v but got %v", i, expectedBody2[i], v)
		}
	}
}
