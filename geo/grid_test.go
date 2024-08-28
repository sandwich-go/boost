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
		b := NewMockCellBuilder(ctrl)
		b.EXPECT().Build(gomock.Any(), gomock.Any()).DoAndReturn(func(x, y int) Cell {
			c := NewMockCell(ctrl)
			return c
		}).AnyTimes()

		var width, height, cellSize int64 = 1200, 1200, 3
		g := G[int64](width, height, cellSize, b)

		c, err := g.GetCellByIndex(int(width/cellSize+1), int(height/cellSize))
		So(err, ShouldNotBeNil)
		c, err = g.GetCellByPoint(graph.P[int64](width+1, height+1))
		So(err, ShouldNotBeNil)

		So(func() {
			_ = g.MustGetCellByIndex(int(width/cellSize+1), int(height/cellSize))
		}, ShouldPanic)

		So(func() {
			_ = g.MustGetCellByPoint(graph.P[int64](width+1, height+1))
		}, ShouldPanic)

		c, err = g.GetCellByIndex(0, 1)
		So(err, ShouldBeNil)
		So(c, ShouldNotBeNil)

		var count int
		g.RangeCellInRectangle(graph.Rect[int64](0, 0, 2, 2), func(p Cell) bool {
			count++
			return true
		})
		So(count, ShouldEqual, 1)

		count = 0
		g.RangeCellInRectangle(graph.Rect[int64](0, 0, 3, 3), func(p Cell) bool {
			count++
			return true
		})
		So(count, ShouldEqual, 4)

		count = 0
		g.RangeCellInRectangle(graph.Rect[int64](0, 0, 4, 4), func(p Cell) bool {
			count++
			return true
		})
		So(count, ShouldEqual, 4)

		count = 0
		g.RangeCellInRectangle(graph.Rect[int64](0, 0, 6, 6), func(p Cell) bool {
			count++
			return true
		})
		So(count, ShouldEqual, 9)

		count = 0
		g.RangeCellInRectangle(graph.Rect[int64](0, 0, 7, 7), func(p Cell) bool {
			count++
			return true
		})
		So(count, ShouldEqual, 9)

		count = 0
		g.RangeCellInLine(graph.L[int64](graph.P[int64](0, 0), graph.P[int64](2, 2)), func(p Cell) bool {
			count++
			return true
		})
		So(count, ShouldEqual, 1)

		count = 0
		g.RangeCellInLine(graph.L[int64](graph.P[int64](0, 0), graph.P[int64](4, 4)), func(p Cell) bool {
			count++
			return true
		})
		So(count, ShouldEqual, 2)
	})

}
