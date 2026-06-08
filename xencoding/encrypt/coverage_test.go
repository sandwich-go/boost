package encrypt

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：补 encrypt 包未覆盖的：
//   - SetKey
//   - Register
//   - codec.Name = "encrypt"
//   - Type.String

func TestSetKey(t *testing.T) {
	Convey("SetKey 设置 AESCodec.key", t, func() {
		// 保存原 key 测完恢复（如果有）
		orig := AESCodec.key
		defer func() { AESCodec.key = orig }()

		key := []byte("0123456789abcdef") // AES-128 16 字节
		SetKey(key)
		So(AESCodec.key, ShouldResemble, key)
	})
}

// stubEncCodec 唯一名 encrypt codec 用于 Register 测试
type stubEncCodec struct{ n string }

func (c stubEncCodec) Name() string                                             { return c.n }
func (c stubEncCodec) Marshal(_ context.Context, _ interface{}) ([]byte, error) { return nil, nil }
func (c stubEncCodec) Unmarshal(_ context.Context, _ []byte, _ interface{}) error {
	return nil
}

func TestRegister_CustomCodec(t *testing.T) {
	Convey("Register 追加自定义 encrypt codec", t, func() {
		const customType Type = 99
		// 用唯一 Name 避免 xencoding.RegisterCodec dup 检查 panic
		Register(customType, stubEncCodec{n: "stub-enc"})
		// 不 panic 即可
	})
}

func TestPackageCodec_Name_AndPaths(t *testing.T) {
	Convey("NewCodec 返回 codec.Name() = \"encrypt\"", t, func() {
		c := NewCodec(NoneType, nil)
		So(c.Name(), ShouldEqual, "encrypt")
	})

	Convey("NewCodec NoneType + Marshal/Unmarshal []byte 透传", t, func() {
		c := NewCodec(NoneType, nil)
		ctx := context.Background()
		input := []byte("plaintext")
		out, err := c.Marshal(ctx, input)
		So(err, ShouldBeNil)
		So(out, ShouldResemble, input)

		var got []byte
		err = c.Unmarshal(ctx, out, &got)
		So(err, ShouldBeNil)
		So(got, ShouldResemble, input)
	})

	Convey("NewCodec AESType + 16 字节 key roundtrip", t, func() {
		key := []byte("0123456789abcdef")
		c := NewCodec(AESType, key)
		ctx := context.Background()
		input := []byte("secret message")
		encrypted, err := c.Marshal(ctx, input)
		So(err, ShouldBeNil)
		So(encrypted, ShouldNotBeEmpty)
		So(encrypted, ShouldNotResemble, input) // 真加密了

		var decrypted []byte
		err = c.Unmarshal(ctx, encrypted, &decrypted)
		So(err, ShouldBeNil)
		So(decrypted, ShouldResemble, input)
	})

	Convey("NewCodec 未知 Type + nil codec → panic", t, func() {
		// codec[unknown] = nil → c.Codec == nil → xpanic.WhenError
		const unknown Type = 200
		So(func() { NewCodec(unknown, nil) }, ShouldPanic)
	})
}

func TestType_String(t *testing.T) {
	Convey("Type.String", t, func() {
		So(NoneType.String(), ShouldNotBeEmpty)
		So(AESType.String(), ShouldNotBeEmpty)
		So(NoneType.String(), ShouldNotEqual, AESType.String())
	})
}

func TestNoneCodec_Unmarshal_NonByteParam(t *testing.T) {
	Convey("NoneCodec.Unmarshal 入参非 *[]byte 返错误", t, func() {
		err := NoneCodec.Unmarshal(context.Background(), []byte("x"), nil)
		So(err, ShouldNotBeNil)
		var s string
		err = NoneCodec.Unmarshal(context.Background(), []byte("x"), &s)
		So(err, ShouldNotBeNil)
	})
}

// TestCodec_Marshal_NonByteParam 覆盖 noneCodec/aesCodec Marshal 的
// 'if !ok → errCodecMarshalParam' 错误路径（line 65 / 86）。
func TestCodec_Marshal_NonByteParam(t *testing.T) {
	Convey("NoneCodec.Marshal 入参非 []byte 返错误", t, func() {
		_, err := NoneCodec.Marshal(context.Background(), "string-not-bytes")
		So(err, ShouldNotBeNil)
		_, err = NoneCodec.Marshal(context.Background(), nil)
		So(err, ShouldNotBeNil)
	})

	Convey("aesCodec.Marshal 入参非 []byte 返错误", t, func() {
		c := NewCodec(AESType, []byte("0123456789abcdef"))
		_, err := c.Marshal(context.Background(), "string-not-bytes")
		So(err, ShouldNotBeNil)
	})
}

// keySetterCodec 实现了 KeySetter 接口的 custom codec，用于覆盖
// NewCodec 的 'KeySetter 分支'（line 124-126）。
type keySetterCodec struct {
	stubEncCodec
	gotKey []byte
}

func (c *keySetterCodec) SetKey(key []byte) { c.gotKey = key }

// TestNewCodec_KeySetter 覆盖 NewCodec 的 'cc.SetKey(key)' 分支
// （line 124-126）。Register 一个实现 KeySetter 的 codec，NewCodec
// 时应调 SetKey 注入 key。
func TestNewCodec_KeySetter(t *testing.T) {
	Convey("NewCodec 对实现 KeySetter 的 codec 注入 key", t, func() {
		const customType Type = 100
		c := &keySetterCodec{stubEncCodec: stubEncCodec{n: "key-setter"}}
		Register(customType, c)

		key := []byte("test-key-123")
		_ = NewCodec(customType, key)
		// SetKey 应被调到，gotKey 设上
		So(c.gotKey, ShouldResemble, key)
	})
}
