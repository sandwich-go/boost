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

// Grid 地块中的基础单元格 地格
type Grid interface{}

// GridBuilder 基础单元格工厂
type GridBuilder interface {
	Build(x, y int) Grid
}

// Block 地块，包含二维的 Grid
type Block[T graph.Number] struct {
	grids    [][]Grid
	gridSize T
	width    T
	height   T
}

// B 新建一个地块
// width 地图的宽
// height 地图的高
// gridSize 基础单元格大小
// gridBuilder 基础单元格工厂
func B[T graph.Number](width, height, gridSize T, gridBuilder GridBuilder) *Block[T] {
	b := &Block[T]{width: width, height: height, gridSize: gridSize}
	b.mustInitialize(gridBuilder)
	return b
}

func (b *Block[T]) toSum(t T) int   { return int(math.Ceil(float64(t) / float64(b.gridSize))) }
func (b *Block[T]) toIndex(t T) int { return int(math.Floor(float64(t) / float64(b.gridSize))) }
func (b *Block[T]) toPoint(x int) T { return T(x) * b.gridSize }

func (b *Block[T]) mustInitialize(gridBuilder GridBuilder) {
	xNum := b.toSum(b.width)
	yNum := b.toSum(b.height)
	b.grids = make([][]Grid, xNum)
	for i := 0; i < xNum; i++ {
		b.grids[i] = make([]Grid, yNum)
		for j := 0; j < yNum; j++ {
			b.grids[i][j] = gridBuilder.Build(i, j)
		}
	}
}

func (b *Block[T]) checkGridIndex(x, y int) bool {
	return x >= 0 && x < len(b.grids) && y >= 0 && y < len(b.grids[0])
}

func (b *Block[T]) checkPoint(point graph.Point[T]) bool {
	return cmp.Compare(point.X(), 0) >= 0 && cmp.Compare(point.X(), b.width) < 0 && cmp.Compare(point.Y(), 0) >= 0 && cmp.Compare(point.Y(), b.height) < 0
}

// GetGridByIndex 通过基础单元格坐标索引 x,y 获取基础单元格
func (b *Block[T]) GetGridByIndex(x, y int) (Grid, error) {
	if !b.checkGridIndex(x, y) {
		return nil, xerror.Wrap(ErrInvalidGridIndex, "index: (%d,%d)", x, y)
	}
	return b.grids[x][y], nil
}

// MustGetGridByIndex 通过基础单元格坐标索引 x,y 获取基础单元格
// 找不到则 panic
func (b *Block[T]) MustGetGridByIndex(x, y int) Grid {
	c, err := b.GetGridByIndex(x, y)
	xpanic.WhenError(err)
	return c
}

// GetGridByPoint 通过坐标， 获取基础单元格
func (b *Block[T]) GetGridByPoint(point graph.Point[T]) (Grid, error) {
	if !b.checkPoint(point) {
		return nil, xerror.Wrap(ErrInvalidGridIndex, "point: %s", point)
	}
	return b.GetGridByIndex(b.toIndex(point.X()), b.toIndex(point.Y()))
}

// MustGetGridByPoint 通过坐标， 获取基础单元格
// 找不到则 panic
func (b *Block[T]) MustGetGridByPoint(point graph.Point[T]) Grid {
	c, err := b.GetGridByPoint(point)
	xpanic.WhenError(err)
	return c
}

// RangeGridInRectangle 在一个矩形中遍历基础单元格
func (b *Block[T]) RangeGridInRectangle(rect graph.Rectangle[T], with func(p Grid) bool) {
	var minP, maxP = rect.Min(), rect.Max()
	var minX, minY, maxX, maxY T = 0, 0, b.width, b.height
	if x := minP.X(); cmp.Compare(x, minX) > 0 {
		minX = x
	}
	if y := minP.Y(); cmp.Compare(y, minY) > 0 {
		minY = y
	}
	// x + b.gridSize - 1 保证右下角的单元格也会被遍历到
	if x := maxP.X() + b.gridSize - 1; cmp.Compare(x, maxX) < 0 {
		maxX = x
	}
	// y + b.gridSize - 1 保证右下角的单元格也会被遍历到
	if y := maxP.Y() + b.gridSize - 1; cmp.Compare(y, maxY) < 0 {
		maxY = y
	}
	for x := b.toIndex(minX); x < b.toSum(maxX); x++ {
		for y := b.toIndex(minY); y < b.toSum(maxY); y++ {
			ok := rect.Has(graph.P(b.toPoint(x), b.toPoint(y)))
			if !ok {
				continue
			}
			c, err := b.GetGridByIndex(x, y)
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
func (b *Block[T]) RangeGridInLine(l graph.Line[T], with func(p Grid) bool) {
	var gotGrids = make(map[int]struct{}, 0)
	l.RangePoints(func(point graph.Point[T]) bool {
		x := b.toIndex(point.X())
		y := b.toIndex(point.Y())
		n := x*len(b.grids) + y
		_, ok := gotGrids[n]
		if ok {
			// 已经遍历过，丢弃
			return true
		}
		gotGrids[n] = struct{}{}
		c, err := b.GetGridByIndex(x, y)
		if err != nil {
			// 不合法的，丢弃
			return true
		}
		return with(c)
	})
}
