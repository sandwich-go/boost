package geo

import (
	"cmp"
	"errors"
	"github.com/sandwich-go/boost/graph"
	"github.com/sandwich-go/boost/xerror"
	"github.com/sandwich-go/boost/xpanic"
	"math"
)

var (
	ErrInvalidGridIndex = errors.New("invalid grid index")
)

//go:generate mockgen -source grid.go -package geo -destination grid_mock_test.go

// Grid 网格中的基础单元格 地格
type Grid interface{}

// GridBuilder 基础单元格工厂
type GridBuilder interface {
	Build(x, y int) Grid
}

// Block 网格，包含二维的 Grid
type Block[T graph.Number] struct {
	grids    [][]Grid
	gridSize T
	width    T
	height   T
}

// G 新建一个网格
// width 地图的宽
// height 地图的高
// gridSize 基础单元格大小
// gridBuilder 基础单元格工厂
func G[T graph.Number](width, height, gridSize T, gridBuilder GridBuilder) *Block[T] {
	g := &Block[T]{width: width, height: height, gridSize: gridSize}
	g.mustInitialize(gridBuilder)
	return g
}

func (g *Block[T]) toSum(t T) int   { return int(math.Ceil(float64(t) / float64(g.gridSize))) }
func (g *Block[T]) toIndex(t T) int { return int(math.Floor(float64(t) / float64(g.gridSize))) }
func (g *Block[T]) toPoint(x int) T { return T(x) * g.gridSize }

func (g *Block[T]) mustInitialize(gridBuilder GridBuilder) {
	xNum := g.toSum(g.width)
	yNum := g.toSum(g.height)
	g.grids = make([][]Grid, xNum)
	for i := 0; i < xNum; i++ {
		g.grids[i] = make([]Grid, yNum)
		for j := 0; j < yNum; j++ {
			g.grids[i][j] = gridBuilder.Build(i, j)
		}
	}
}

func (g *Block[T]) checkGridIndex(x, y int) bool {
	return x >= 0 && x < len(g.grids) && y >= 0 && y < len(g.grids[0])
}

func (g *Block[T]) checkPoint(point graph.Point[T]) bool {
	return cmp.Compare(point.X(), 0) >= 0 && cmp.Compare(point.X(), g.width) < 0 && cmp.Compare(point.Y(), 0) >= 0 && cmp.Compare(point.Y(), g.height) < 0
}

// GetGridByIndex 通过基础单元格坐标索引 x,y 获取基础单元格
func (g *Block[T]) GetGridByIndex(x, y int) (Grid, error) {
	if !g.checkGridIndex(x, y) {
		return nil, xerror.Wrap(ErrInvalidGridIndex, "index: (%d,%d)", x, y)
	}
	return g.grids[x][y], nil
}

// MustGetGridByIndex 通过基础单元格坐标索引 x,y 获取基础单元格
// 找不到则 panic
func (g *Block[T]) MustGetGridByIndex(x, y int) Grid {
	c, err := g.GetGridByIndex(x, y)
	xpanic.WhenError(err)
	return c
}

// GetGridByPoint 通过坐标， 获取基础单元格
func (g *Block[T]) GetGridByPoint(point graph.Point[T]) (Grid, error) {
	if !g.checkPoint(point) {
		return nil, xerror.Wrap(ErrInvalidGridIndex, "point: %s", point)
	}
	return g.GetGridByIndex(g.toIndex(point.X()), g.toIndex(point.Y()))
}

// MustGetGridByPoint 通过坐标， 获取基础单元格
// 找不到则 panic
func (g *Block[T]) MustGetGridByPoint(point graph.Point[T]) Grid {
	c, err := g.GetGridByPoint(point)
	xpanic.WhenError(err)
	return c
}

// RangeGridInRectangle 在一个矩形中遍历基础单元格
func (g *Block[T]) RangeGridInRectangle(rect graph.Rectangle[T], with func(p Grid) bool) {
	var minP, maxP = rect.Min(), rect.Max()
	var minX, minY, maxX, maxY T = 0, 0, g.width, g.height
	if x := minP.X(); cmp.Compare(x, minX) > 0 {
		minX = x
	}
	if y := minP.Y(); cmp.Compare(y, minY) > 0 {
		minY = y
	}
	// x + g.gridSize - 1 保证右下角的单元格也会被遍历到
	if x := maxP.X() + g.gridSize - 1; cmp.Compare(x, maxX) < 0 {
		maxX = x
	}
	// y + g.gridSize - 1 保证右下角的单元格也会被遍历到
	if y := maxP.Y() + g.gridSize - 1; cmp.Compare(y, maxY) < 0 {
		maxY = y
	}
	for x := g.toIndex(minX); x < g.toSum(maxX); x++ {
		for y := g.toIndex(minY); y < g.toSum(maxY); y++ {
			ok := rect.Has(graph.P(g.toPoint(x), g.toPoint(y)))
			if !ok {
				continue
			}
			c, err := g.GetGridByIndex(x, y)
			if err != nil {
				// 不合法的，丢弃
				continue
			}
			if !with(c) {
				return
			}
		}
	}
}

// RangeGridInLine 在一个线中遍历基础单元格
func (g *Block[T]) RangeGridInLine(l graph.Line[T], with func(p Grid) bool) {
	var gotGrids = make(map[int]struct{}, 0)
	l.RangePoints(func(point graph.Point[T]) bool {
		x := g.toIndex(point.X())
		y := g.toIndex(point.Y())
		n := x*len(g.grids) + y
		_, ok := gotGrids[n]
		if ok {
			// 已经遍历过，丢弃
			return true
		}
		gotGrids[n] = struct{}{}
		c, err := g.GetGridByIndex(x, y)
		if err != nil {
			// 不合法的，丢弃
			return true
		}
		return with(c)
	})
}
