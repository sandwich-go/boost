package xencoding

import (
	"context"
	"sort"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：补 xencoding root 包未覆盖的 Codecs 函数。

// stubCodec 仅用于注册测试
type stubCodec struct{ name string }

func (c *stubCodec) Name() string                                         { return c.name }
func (c *stubCodec) Marshal(context.Context, interface{}) ([]byte, error) { return nil, nil }
func (c *stubCodec) Unmarshal(context.Context, []byte, interface{}) error { return nil }

func TestCodecs_Listing(t *testing.T) {
	Convey("Codecs 返回已注册的 codec 名字（排序）", t, func() {
		// 注册两个 codec（用唯一名避免与现有冲突）
		RegisterCodec(&stubCodec{name: "stub-zzz"})
		RegisterCodec(&stubCodec{name: "stub-aaa"})

		names := Codecs()
		// 验证已排序
		sortedNames := make([]string, len(names))
		copy(sortedNames, names)
		sort.Strings(sortedNames)
		So(names, ShouldResemble, sortedNames)

		// 我们注册的两个应该在里面
		So(names, ShouldContain, "stub-aaa")
		So(names, ShouldContain, "stub-zzz")
	})
}
