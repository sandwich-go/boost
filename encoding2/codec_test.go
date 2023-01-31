package encoding2

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type MockCodec struct {
	name string
}

func (mc *MockCodec) Marshal(v interface{}) ([]byte, error)      { return nil, nil }
func (mc *MockCodec) Unmarshal(data []byte, v interface{}) error { return nil }
func (mc *MockCodec) Name() string                               { return mc.name }

func TestCodec(t *testing.T) {
	mc := &MockCodec{name: "mock_test"}

	Convey("register codec will panic if codec nil or name empty", t, func() {
		RegisterCodec(mc)
		So(GetCodec(mc.Name()), ShouldEqual, mc)
		So(func() { RegisterCodec(nil) }, ShouldPanic)
		So(func() { RegisterCodec(&MockCodec{}) }, ShouldPanic)
	})
	Convey("get nil codec", t, func() {
		So(FromContext(context.Background()), ShouldEqual, nil)
	})

	Convey("with golang standard context", t, func() {
		ctx := WithContext(context.Background(), mc)
		So(FromContext(ctx), ShouldEqual, mc)
	})
}
