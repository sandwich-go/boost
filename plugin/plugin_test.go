package plugin

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type mockEntry struct{ num int }

type before interface {
	Name() string
	Before()
	Context() context.Context
}

type filter interface {
	Name() string
	Filter(*mockEntry) *mockEntry
	Context() context.Context
}

type mockBefore struct{ ctx context.Context }

func newMockBefore(ctx context.Context) before { return &mockBefore{ctx: ctx} }
func (b mockBefore) Name() string              { return "Before" }
func (b mockBefore) Context() context.Context  { return b.ctx }
func (b *mockBefore) Before()                  { b.ctx = context.WithValue(b.ctx, b.Name(), struct{}{}) }

type mockFilter struct {
	ctx    context.Context
	maxNum int
}

func newMockFilter(ctx context.Context, maxNum int) filter {
	return &mockFilter{ctx: ctx, maxNum: maxNum}
}
func (b mockFilter) Name() string             { return "Filter" }
func (b mockFilter) Context() context.Context { return b.ctx }
func (b *mockFilter) Filter(in *mockEntry) *mockEntry {
	b.ctx = context.WithValue(b.ctx, b.Name(), struct{}{})
	if in.num > b.maxNum {
		return nil
	}
	return in
}

func TestPlugin(t *testing.T) {
	Convey("Test Plugin", t, func(c C) {
		cc := New(new(before), new(filter))
		So(cc.Add(new(mockEntry)), ShouldNotBeNil)
		So(func() { cc.MustAdd(new(mockEntry)) }, ShouldPanic)

		maxNum := 2
		ctx := context.Background()
		mockBeforeObj := newMockBefore(ctx)
		mockFilterObj := newMockFilter(ctx, maxNum)
		So(cc.Add(mockBeforeObj), ShouldBeNil)
		cc.MustAdd(mockFilterObj)

		var es = make([]*mockEntry, 0)
		for i := 0; i <= maxNum; i++ {
			es = append(es, &mockEntry{num: i + 1})
		}

		cc.Range(func(plugin Plugin) bool {
			if p, ok := plugin.(before); ok {
				p.Before()
			}
			return true
		})

		cc.Range(func(plugin Plugin) bool {
			if p, ok := plugin.(filter); ok {
				var out = make([]*mockEntry, 0, len(es))
				for _, e := range es {
					if e = p.Filter(e); e != nil {
						out = append(out, e)
					}
				}
				So(len(out)+1, ShouldEqual, len(es))
			}
			return true
		})

		So(mockBeforeObj.Context().Value(mockBeforeObj.Name()), ShouldNotBeNil)
		So(mockFilterObj.Context().Value(mockFilterObj.Name()), ShouldNotBeNil)
	})
}
