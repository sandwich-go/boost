package graph

import (
	"cmp"
	"fmt"
	"math"
)

// Line 线
type Line[T Number] struct {
	start, end Point[T]
}

// L return shorthand for Line[T].{start, end}.
func L[T Number](start, end Point[T]) Line[T] { return Line[T]{start: start, end: end} }

// String returns a string representation of r like "L{(3,4)-(6,5)}".
func (x Line[T]) String() string { return fmt.Sprintf("L{%s-%s}", x.start.String(), x.end.String()) }

// Start returns l's start point.
func (x Line[T]) Start() Point[T] { return x.start }

// End returns l's end point.
func (x Line[T]) End() Point[T] { return x.end }

// Equals reports whether r and s contain the same set of points. All empty
// rectangles are considered equal.
func (x Line[T]) Equals(l Line[T]) bool { return x == l }

// RangePoints range all points in rectangle.
// if with return false, aborted range.
func (x Line[T]) RangePoints(with func(p Point[T]) bool) {
	var ZL Line[T]
	if with == nil || x == ZL {
		return
	}

	var ps []Point[T]
	ps = append(ps, x.start)
	if x.start.Equals(x.end) {
		with(x.start)
		return
	}

	if !with(x.start) {
		return
	}
	xy := float64(x.end.y - x.start.y)
	xx := float64(x.end.x - x.start.x)
	xyabs := math.Abs(xy)
	xxabs := math.Abs(xx)
	var positive = false
	if cmp.Compare(xyabs, xxabs) >= 0 {
		positive = cmp.Compare(xy, 0) >= 0
		per := xx / xy
		posY := x.start.y
		var i float64
		for i = 0; cmp.Compare(i, xyabs) < 0; i++ {
			var posX T
			if positive {
				posY = posY + 1
				posX = (T)(math.Round(float64(x.start.x) + per*(i+1)))
			} else {
				posY = posY - 1
				posX = (T)(math.Round(float64(x.start.x) - per*(i+1)))
			}
			if !with(P(posX, posY)) {
				return
			}
		}
	} else {
		positive = cmp.Compare(xx, 0) >= 0
		per := xy / xx
		posX := x.start.x
		var i float64
		for i = 0; cmp.Compare(i, xxabs) < 0; i++ {
			var posY T
			if positive {
				posX = posX + 1
				posY = (T)(math.Round(float64(x.start.y) + per*(i+1)))
			} else {
				posX = posX - 1
				posY = (T)(math.Round(float64(x.start.y) - per*(i+1)))
			}
			if !with(P(posX, posY)) {
				return
			}
		}
	}
}
