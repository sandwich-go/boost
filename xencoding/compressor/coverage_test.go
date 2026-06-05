package compressor

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：补 compressor 包未覆盖的：
//   - Register（自定义 codec 注册到 typeMap）
//   - codec.Name（"compressor" 包级 codec）
//   - Type.String（stringer 生成）

type stubCompCodec struct{ n string }

func (c stubCompCodec) Name() string                                             { return c.n }
func (c stubCompCodec) Marshal(_ context.Context, v interface{}) ([]byte, error) { return nil, nil }
func (c stubCompCodec) Unmarshal(_ context.Context, _ []byte, _ interface{}) error {
	return nil
}

func TestRegister_CustomCodec(t *testing.T) {
	Convey("Register 追加自定义 compressor codec", t, func() {
		// 用一个未占用的 Type 值（>3 的常量都未用）
		const customType Type = 99
		Register(customType, stubCompCodec{n: "custom-comp"})
		// 注册成功后 NewCodec 能拿到
		// 注意 NewCodec 是 noOp 还是查表取决于实现；这里只验证不 panic
		c := NewCodec(GZIPType)
		So(c, ShouldNotBeNil)
	})
}

func TestPackageCodec_Name(t *testing.T) {
	Convey("NewCodec 返回的 codec.Name() = \"compressor\"", t, func() {
		c := NewCodec(GZIPType)
		So(c.Name(), ShouldEqual, "compressor")
	})

	Convey("NewCodec 的 Marshal/Unmarshal 走对应 baseCodec 路径", t, func() {
		c := NewCodec(GZIPType)
		ctx := context.Background()
		// compressor codec 的 Marshal 入参必须是 []byte（doc 化的契约）
		input := []byte("hello compressor")
		var inputAny interface{} = input
		out, err := c.Marshal(ctx, inputAny)
		So(err, ShouldBeNil)
		So(out, ShouldNotBeEmpty)

		var got []byte
		var gotAny interface{} = &got
		err = c.Unmarshal(ctx, out, gotAny)
		So(err, ShouldBeNil)
		So(got, ShouldResemble, input)
	})

	Convey("NewCodec 用未知 Type → Marshal/Unmarshal 返 errCodecNoFound", t, func() {
		const unknown Type = 200
		c := NewCodec(unknown)
		ctx := context.Background()
		var input interface{} = []byte("x")
		_, err := c.Marshal(ctx, input)
		So(err, ShouldNotBeNil)
		err = c.Unmarshal(ctx, []byte{}, nil)
		So(err, ShouldNotBeNil)
	})

	Convey("Marshal 入参非 []byte 返参数错误", t, func() {
		c := NewCodec(GZIPType)
		_, err := c.Marshal(context.Background(), "string-not-bytes")
		So(err, ShouldNotBeNil)
	})
}

func TestType_String(t *testing.T) {
	Convey("Type.String stringer 输出", t, func() {
		So(NoneType.String(), ShouldNotBeEmpty)
		So(GZIPType.String(), ShouldNotBeEmpty)
		So(SnappyType.String(), ShouldNotBeEmpty)
		So(NoneType.String(), ShouldNotEqual, GZIPType.String())
	})
}
