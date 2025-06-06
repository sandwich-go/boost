package geo

import (
	"github.com/sandwich-go/boost/graph"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestGrid(t *testing.T) {
	ctrl := gomock.NewController(t)

	Convey("test grid", t, func() {
		b := NewMockGridBuilder(ctrl)
		b.EXPECT().Build(gomock.Any(), gomock.Any()).DoAndReturn(func(x, y int) Grid {
			c := NewMockGrid(ctrl)
			return c
		}).AnyTimes()

		var width, height, GridSize int64 = 1200, 1200, 3
		g := B[int64](width, height, GridSize, b)

		c, err := g.GetGridByIndex(int(width/GridSize+1), int(height/GridSize))
		So(err, ShouldNotBeNil)
		c, err = g.GetGridByPoint(graph.P[int64](width+1, height+1))
		So(err, ShouldNotBeNil)

		So(func() {
			_ = g.MustGetGridByIndex(int(width/GridSize+1), int(height/GridSize))
		}, ShouldPanic)

		So(func() {
			_ = g.MustGetGridByPoint(graph.P[int64](width+1, height+1))
		}, ShouldPanic)

		c, err = g.GetGridByIndex(0, 1)
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		var count int
		g.RangeGridInRectangle(graph.Rect[int64](0, 0, 2, 2), func(p Grid) bool {
			count++
			return true
		})
		So(count, ShouldEqual, 1)

		count = 0
		g.RangeGridInRectangle(graph.Rect[int64](0, 0, 3, 3), func(p Grid) bool {
			count++
			return true
		})
		So(count, ShouldEqual, 4)

		count = 0
		g.RangeGridInRectangle(graph.Rect[int64](0, 0, 4, 4), func(p Grid) bool {
			count++
			return true
		})
		So(count, ShouldEqual, 4)

		count = 0
		g.RangeGridInRectangle(graph.Rect[int64](0, 0, 6, 6), func(p Grid) bool {
			count++
			return true
		})
		So(count, ShouldEqual, 9)

		count = 0
		g.RangeGridInRectangle(graph.Rect[int64](0, 0, 7, 7), func(p Grid) bool {
			count++
			return true
		})
		So(count, ShouldEqual, 9)

		count = 0
		g.RangeGridInLine(graph.L[int64](graph.P[int64](0, 0), graph.P[int64](2, 2)), func(p Grid) bool {
			count++
			return true
		})
		So(count, ShouldEqual, 1)

		count = 0
		g.RangeGridInLine(graph.L[int64](graph.P[int64](0, 0), graph.P[int64](4, 4)), func(p Grid) bool {
			count++
			return true
		})
		So(count, ShouldEqual, 2)
	})

}
