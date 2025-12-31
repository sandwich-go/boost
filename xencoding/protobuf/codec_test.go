package protobuf

import (
	"bytes"
	"context"
	"sync"
	"testing"

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

func TestBasicProtoCodecMarshalAndUnmarshal(t *testing.T) {
	marshalAndUnmarshal(t, &codec{}, []byte{1, 2, 3})
}
func TestBasicJsonCodecMarshalAndUnmarshal(t *testing.T) {
	marshalAndUnmarshal(t, &codec{}, []byte{1, 2, 3})
}

// Try to catch possible race conditions around use of pools
func TestConcurrentUsage(t *testing.T) {
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
	codec := &codec{}

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
		t.Errorf("codec.Marshal failed: %v", err)
	}

	if m2, err = codec2.Marshal(context.Background(), &proto2); err != nil {
		t.Errorf("codec.Marshal failed: %v", err)
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

// TestVTProtoMarshal 测试 vtproto 的 MarshalVT 方法
func TestVTProtoMarshal(t *testing.T) {
	expectedBody := []byte("test vtproto marshal")
	msg := &test_perf.Buffer{Body: expectedBody}

	// 使用不带 pool 的 codec
	c := codec{usingPool: false, name: CodecName}
	data, err := c.Marshal(context.Background(), msg)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 验证可以正确反序列化
	result := &test_perf.Buffer{}
	if err := c.Unmarshal(context.Background(), data, result); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if !bytes.Equal(result.Body, expectedBody) {
		t.Errorf("Expected body %v, got %v", expectedBody, result.Body)
	}
}

// TestVTProtoMarshalWithPool 测试使用 pool 的 vtproto marshal
func TestVTProtoMarshalWithPool(t *testing.T) {
	expectedBody := []byte("test vtproto marshal with pool")
	msg := &test_perf.Buffer{Body: expectedBody}

	// 使用带 pool 的 codec
	c := codec{usingPool: true, name: UsingPoolCodecName}
	data, err := c.Marshal(context.Background(), msg)
	if err != nil {
		t.Fatalf("Marshal with pool failed: %v", err)
	}

	// 验证可以正确反序列化
	result := &test_perf.Buffer{}
	if err := c.Unmarshal(context.Background(), data, result); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if !bytes.Equal(result.Body, expectedBody) {
		t.Errorf("Expected body %v, got %v", expectedBody, result.Body)
	}
}

// TestVTProtoUnmarshal 测试 vtproto 的 UnmarshalVT 方法
func TestVTProtoUnmarshal(t *testing.T) {
	expectedBody := []byte("test vtproto unmarshal")
	msg := &test_perf.Buffer{Body: expectedBody}

	// 先用 VT 方法序列化
	data, err := msg.MarshalVT()
	if err != nil {
		t.Fatalf("MarshalVT failed: %v", err)
	}

	// 使用 codec 的 Unmarshal（应该会调用 UnmarshalVT）
	c := codec{usingPool: false, name: CodecName}
	result := &test_perf.Buffer{}
	if err := c.Unmarshal(context.Background(), data, result); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if !bytes.Equal(result.Body, expectedBody) {
		t.Errorf("Expected body %v, got %v", expectedBody, result.Body)
	}
}

// TestVTProtoCompatibility 测试 vtproto 与标准 proto 的兼容性
func TestVTProtoCompatibility(t *testing.T) {
	expectedBody := []byte("test compatibility between vtproto and proto")

	t.Run("VT_marshal_proto_unmarshal", func(t *testing.T) {
		// VT 序列化
		msg := &test_perf.Buffer{Body: expectedBody}
		data, err := msg.MarshalVT()
		if err != nil {
			t.Fatalf("MarshalVT failed: %v", err)
		}

		// 标准 proto 反序列化
		result := &test_perf.Buffer{}
		if err := result.UnmarshalVT(data); err != nil {
			t.Fatalf("UnmarshalVT failed: %v", err)
		}

		if !bytes.Equal(result.Body, expectedBody) {
			t.Errorf("Expected body %v, got %v", expectedBody, result.Body)
		}
	})

	t.Run("codec_marshal_unmarshal_consistency", func(t *testing.T) {
		msg := &test_perf.Buffer{Body: expectedBody}

		// 使用不带 pool 的 codec
		c1 := codec{usingPool: false, name: CodecName}
		data1, err := c1.Marshal(context.Background(), msg)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		// 使用带 pool 的 codec
		c2 := codec{usingPool: true, name: UsingPoolCodecName}
		data2, err := c2.Marshal(context.Background(), msg)
		if err != nil {
			t.Fatalf("Marshal with pool failed: %v", err)
		}

		// 两种方式序列化的结果应该一致
		if !bytes.Equal(data1, data2) {
			t.Errorf("Marshal results differ between pool and non-pool codecs")
		}

		// 两种数据都应该能正确反序列化
		result1 := &test_perf.Buffer{}
		if err := c1.Unmarshal(context.Background(), data1, result1); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		result2 := &test_perf.Buffer{}
		if err := c2.Unmarshal(context.Background(), data2, result2); err != nil {
			t.Fatalf("Unmarshal with pool failed: %v", err)
		}

		if !bytes.Equal(result1.Body, expectedBody) || !bytes.Equal(result2.Body, expectedBody) {
			t.Errorf("Unmarshal results incorrect")
		}
	})
}

// TestUnmarshalReset 测试 Unmarshal 时的 Reset 行为
func TestUnmarshalReset(t *testing.T) {
	// 第一次设置数据
	msg := &test_perf.Buffer{Body: []byte("old data")}

	// 序列化新数据（body 为空）
	newMsg := &test_perf.Buffer{}
	c := codec{usingPool: false, name: CodecName}
	data, err := c.Marshal(context.Background(), newMsg)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// 反序列化到已有数据的 msg 中，应该被重置
	if err := c.Unmarshal(context.Background(), data, msg); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// 由于 Reset，旧数据应该被清除
	if len(msg.Body) != 0 {
		t.Errorf("Expected body to be empty after unmarshal, got %v", msg.Body)
	}
}

// TestStandardMarshalWithPool 测试 standardMarshalWithPool 的正确性
func TestStandardMarshalWithPool(t *testing.T) {
	expectedBody := []byte("test standard marshal with pool")
	msg := &test_perf.Buffer{Body: expectedBody}

	t.Run("basic_marshal_unmarshal", func(t *testing.T) {
		// 使用带 pool 的 codec（即使消息支持 vtproto，也测试标准方法）
		c := codec{usingPool: true, name: UsingPoolCodecName}
		data, err := c.Marshal(context.Background(), msg)
		if err != nil {
			t.Fatalf("Marshal with pool failed: %v", err)
		}

		// 验证可以正确反序列化
		result := &test_perf.Buffer{}
		if err := c.Unmarshal(context.Background(), data, result); err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if !bytes.Equal(result.Body, expectedBody) {
			t.Errorf("Expected body %v, got %v", expectedBody, result.Body)
		}
	})

	t.Run("consistency_with_standard_marshal", func(t *testing.T) {
		// 使用带 pool 的方式序列化
		cPool := codec{usingPool: true, name: UsingPoolCodecName}
		dataPool, err := cPool.Marshal(context.Background(), msg)
		if err != nil {
			t.Fatalf("Marshal with pool failed: %v", err)
		}

		// 使用不带 pool 的方式序列化
		cNoPool := codec{usingPool: false, name: CodecName}
		dataNoPool, err := cNoPool.Marshal(context.Background(), msg)
		if err != nil {
			t.Fatalf("Marshal without pool failed: %v", err)
		}

		// 两种方式的结果应该一致
		if !bytes.Equal(dataPool, dataNoPool) {
			t.Errorf("Marshal results differ between pool and non-pool")
		}
	})

	t.Run("concurrent_usage", func(t *testing.T) {
		const numGoroutines = 50
		const numIterations = 100

		testBodies := [][]byte{
			[]byte("concurrent test 1"),
			[]byte("concurrent test 2"),
			[]byte("concurrent test 3"),
			[]byte("concurrent test 4"),
			[]byte("concurrent test 5"),
		}

		var wg sync.WaitGroup
		c := codec{usingPool: true, name: UsingPoolCodecName}

		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				for j := 0; j < numIterations; j++ {
					body := testBodies[j%len(testBodies)]
					msg := &test_perf.Buffer{Body: body}

					// 序列化
					data, err := c.Marshal(context.Background(), msg)
					if err != nil {
						t.Errorf("Marshal failed: %v", err)
						return
					}

					// 反序列化
					result := &test_perf.Buffer{}
					if err := c.Unmarshal(context.Background(), data, result); err != nil {
						t.Errorf("Unmarshal failed: %v", err)
						return
					}

					// 验证数据正确性
					if !bytes.Equal(result.Body, body) {
						t.Errorf("Body mismatch: expected %v, got %v", body, result.Body)
					}
				}
			}(i)
		}

		wg.Wait()
	})

	t.Run("large_message", func(t *testing.T) {
		// 测试大消息（超过初始 pool 容量 512 字节）
		largeBody := make([]byte, 10000)
		for i := range largeBody {
			largeBody[i] = byte(i % 256)
		}

		largeMsg := &test_perf.Buffer{Body: largeBody}
		c := codec{usingPool: true, name: UsingPoolCodecName}

		data, err := c.Marshal(context.Background(), largeMsg)
		if err != nil {
			t.Fatalf("Marshal large message failed: %v", err)
		}

		result := &test_perf.Buffer{}
		if err := c.Unmarshal(context.Background(), data, result); err != nil {
			t.Fatalf("Unmarshal large message failed: %v", err)
		}

		if !bytes.Equal(result.Body, largeBody) {
			t.Errorf("Large body mismatch")
		}
	})

	t.Run("empty_message", func(t *testing.T) {
		// 测试空消息
		emptyMsg := &test_perf.Buffer{}
		c := codec{usingPool: true, name: UsingPoolCodecName}

		data, err := c.Marshal(context.Background(), emptyMsg)
		if err != nil {
			t.Fatalf("Marshal empty message failed: %v", err)
		}

		result := &test_perf.Buffer{}
		if err := c.Unmarshal(context.Background(), data, result); err != nil {
			t.Fatalf("Unmarshal empty message failed: %v", err)
		}

		if len(result.Body) != 0 {
			t.Errorf("Expected empty body, got %v", result.Body)
		}
	})
}
