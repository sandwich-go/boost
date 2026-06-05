package pbjson

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：补 pbjson 包未覆盖的 UseEnumNumbers setter。

func TestUseEnumNumbers(t *testing.T) {
	Convey("UseEnumNumbers 切换 marshaler.UseEnumNumbers 字段", t, func() {
		// 保存原值恢复
		orig := marshaler.UseEnumNumbers
		defer func() { marshaler.UseEnumNumbers = orig }()

		UseEnumNumbers(false)
		So(marshaler.UseEnumNumbers, ShouldBeFalse)

		UseEnumNumbers(true)
		So(marshaler.UseEnumNumbers, ShouldBeTrue)
	})
}

func TestCodec_NonProtoParam(t *testing.T) {
	Convey("Marshal 入参非 proto.Message 返 errCodecParam", t, func() {
		_, err := Codec.Marshal(context.Background(), "not proto")
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, errCodecParam)
	})

	Convey("Unmarshal 入参非 proto.Message 返 errCodecParam", t, func() {
		var s string
		err := Codec.Unmarshal(context.Background(), []byte("{}"), &s)
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, errCodecParam)
	})
}
