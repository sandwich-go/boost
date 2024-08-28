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
	ErrInvalidCellIndex = errors.New("invalid cell index")
)

//go:generate mockgen -source grid.go -package geo -destination grid_mock_test.go

// Cell 网格中的基础单元格
type Cell interface{}

// CellBuilder 基础单元格工厂
type CellBuilder interface {
	Build(x, y int) Cell
}

// Grid 网格，包含二维的 cell
type Grid[T graph.Number] struct {
	cells    [][]Cell
	cellSize T
	width    T
	height   T
}

// G 新建一个网格
// width 地图的宽
// height 地图的高
// cellSize 基础单元格大小
// cellBuilder 基础单元格工厂
func G[T graph.Number](width, height, cellSize T, cellBuilder CellBuilder) *Grid[T] {
	g := &Grid[T]{width: width, height: height, cellSize: cellSize}
	g.mustInitialize(cellBuilder)
	return g
}

func (g *Grid[T]) toSum(t T) int   { return int(math.Ceil(float64(t) / float64(g.cellSize))) }
func (g *Grid[T]) toIndex(t T) int { return int(math.Floor(float64(t) / float64(g.cellSize))) }
func (g *Grid[T]) toPoint(x int) T { return T(x) * g.cellSize }

func (g *Grid[T]) mustInitialize(cellBuilder CellBuilder) {
	xNum := g.toSum(g.width)
	yNum := g.toSum(g.height)
	g.cells = make([][]Cell, xNum)
	for i := 0; i < xNum; i++ {
		g.cells[i] = make([]Cell, yNum)
		for j := 0; j < yNum; j++ {
			g.cells[i][j] = cellBuilder.Build(i, j)
		}
	}
}

func (g *Grid[T]) checkCellIndex(x, y int) bool {
	return x >= 0 && x < len(g.cells) && y >= 0 && y < len(g.cells[0])
}

func (g *Grid[T]) checkPoint(point graph.Point[T]) bool {
	return cmp.Compare(point.X(), 0) >= 0 && cmp.Compare(point.X(), g.width) < 0 && cmp.Compare(point.Y(), 0) >= 0 && cmp.Compare(point.Y(), g.height) < 0
}

// GetCellByIndex 通过基础单元格坐标索引 x,y 获取基础单元格
func (g *Grid[T]) GetCellByIndex(x, y int) (Cell, error) {
	if !g.checkCellIndex(x, y) {
		return nil, xerror.Wrap(ErrInvalidCellIndex, "index: (%d,%d)", x, y)
	}
	return g.cells[x][y], nil
}

// MustGetCellByIndex 通过基础单元格坐标索引 x,y 获取基础单元格
// 找不到则 panic
func (g *Grid[T]) MustGetCellByIndex(x, y int) Cell {
	c, err := g.GetCellByIndex(x, y)
	xpanic.WhenError(err)
	return c
}

// GetCellByPoint 通过坐标， 获取基础单元格
func (g *Grid[T]) GetCellByPoint(point graph.Point[T]) (Cell, error) {
	if !g.checkPoint(point) {
		return nil, xerror.Wrap(ErrInvalidCellIndex, "point: %s", point)
	}
	return g.GetCellByIndex(g.toIndex(point.X()), g.toIndex(point.Y()))
}

// MustGetCellByPoint 通过坐标， 获取基础单元格
// 找不到则 panic
func (g *Grid[T]) MustGetCellByPoint(point graph.Point[T]) Cell {
	c, err := g.GetCellByPoint(point)
	xpanic.WhenError(err)
	return c
}

// RangeCellInRectangle 在一个矩形中遍历基础单元格
func (g *Grid[T]) RangeCellInRectangle(rect graph.Rectangle[T], with func(p Cell) bool) {
	var minP, maxP = rect.Min(), rect.Max()
	var minX, minY, maxX, maxY T = 0, 0, g.width, g.height
	if x := minP.X(); cmp.Compare(x, minX) > 0 {
		minX = x
	}
	if y := minP.Y(); cmp.Compare(y, minY) > 0 {
		minY = y
	}
	// x + g.cellSize - 1 保证右下角的单元格也会被遍历到
	if x := maxP.X() + g.cellSize - 1; cmp.Compare(x, maxX) < 0 {
		maxX = x
	}
	// y + g.cellSize - 1 保证右下角的单元格也会被遍历到
	if y := maxP.Y() + g.cellSize - 1; cmp.Compare(y, maxY) < 0 {
		maxY = y
	}
	for x := g.toIndex(minX); x < g.toSum(maxX); x++ {
		for y := g.toIndex(minY); y < g.toSum(maxY); y++ {
			ok := rect.Has(graph.P(g.toPoint(x), g.toPoint(y)))
			if !ok {
				continue
			}
			c, err := g.GetCellByIndex(x, y)
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

// RangeCellInLine 在一个线中遍历基础单元格
func (g *Grid[T]) RangeCellInLine(l graph.Line[T], with func(p Cell) bool) {
	var gotCells = make(map[int]struct{}, 0)
	l.RangePoints(func(point graph.Point[T]) bool {
		x := g.toIndex(point.X())
		y := g.toIndex(point.Y())
		n := x*len(g.cells) + y
		_, ok := gotCells[n]
		if ok {
			// 已经遍历过，丢弃
			return true
		}
		gotCells[n] = struct{}{}
		c, err := g.GetCellByIndex(x, y)
		if err != nil {
			// 不合法的，丢弃
			return true
		}
		return with(c)
	})
}
