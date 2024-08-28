package graph

import (
	"cmp"
	"fmt"
)

// Rectangle wrapper, contain two Point[T].
type Rectangle[T Number] struct {
	min, max Point[T]
}

// R return shorthand for Rectangle[T].{min, max}.
func R[T Number](min, max Point[T]) Rectangle[T] { return Rectangle[T]{min: min, max: max} }

// String returns a string representation of r like "R{(3,4)-(6,5)}".
func (r Rectangle[T]) String() string {
	return fmt.Sprintf("R{%s-%s}", r.min.String(), r.max.String())
}

func (r Rectangle[T]) Min() Point[T] { return r.min }
func (r Rectangle[T]) Max() Point[T] { return r.max }
func (r Rectangle[T]) Center() Point[T] {
	return P(r.min.x+(r.max.x-r.min.x)/2, r.min.y+(r.max.y-r.min.y)/2)
}

// RangePoints range all points in rectangle.
// if with return false, aborted range.
func (r Rectangle[T]) RangePoints(with func(p Point[T]) bool) {
	var ZR Rectangle[T]
	if with == nil || r == ZR {
		return
	}

	for x := r.min.x; cmp.Compare(x, r.max.x) <= 0; x++ {
		for y := r.min.y; cmp.Compare(y, r.max.y) <= 0; y++ {
			if !with(P(x, y)) {
				return
			}
		}
	}
}

// RangePointsMinClosedMaxOpen range all points in rectangle except min x, y.
// if with return false, aborted range.
func (r Rectangle[T]) RangePointsMinClosedMaxOpen(with func(p Point[T]) bool) {
	var ZR Rectangle[T]
	if with == nil || r == ZR {
		return
	}

	for x := r.min.x + 1; cmp.Compare(x, r.max.x) <= 0; x++ {
		for y := r.min.y + 1; cmp.Compare(y, r.max.y) <= 0; y++ {
			if !with(P(x, y)) {
				return
			}
		}
	}
}

// RangePointsMinMaxClosed range all points in rectangle except min/max x, y.
// if with return false, aborted range.
func (r Rectangle[T]) RangePointsMinMaxClosed(with func(p Point[T]) bool) {
	var ZR Rectangle[T]
	if with == nil || r == ZR {
		return
	}

	for x := r.min.x + 1; cmp.Compare(x, r.max.x) < 0; x++ {
		for y := r.min.y + 1; cmp.Compare(y, r.max.y) < 0; y++ {
			if !with(P(x, y)) {
				return
			}
		}
	}
}

// RangePointsMinOpenMaxClosed range all points in rectangle except max x, y.
// if with return false, aborted range.
func (r Rectangle[T]) RangePointsMinOpenMaxClosed(with func(p Point[T]) bool) {
	var ZR Rectangle[T]
	if with == nil || r == ZR {
		return
	}

	for x := r.min.x; cmp.Compare(x, r.max.x) < 0; x++ {
		for y := r.min.y; cmp.Compare(y, r.max.y) < 0; y++ {
			if !with(P(x, y)) {
				return
			}
		}
	}
}

// IntersectionWithLine check line intersection, if intersection, return true
func (r Rectangle[T]) IntersectionWithLine(s Point[T], e Point[T]) bool {
	if cmp.Compare(s.x, r.min.x) <= 0 && cmp.Compare(e.x, r.min.x) <= 0 || cmp.Compare(s.x, r.max.x) >= 0 && cmp.Compare(e.x, r.max.x) >= 0 ||
		cmp.Compare(s.y, r.min.y) <= 0 && cmp.Compare(e.y, r.min.y) <= 0 || cmp.Compare(s.y, r.max.y) >= 0 && cmp.Compare(e.y, r.max.y) >= 0 {
		return false
	}

	a := s.y - e.y
	b := e.x - s.x
	c := e.y*s.x - e.x*s.y
	if cmp.Compare((a*r.min.x+b*r.min.y+c)*(a*r.max.x+b*r.max.y+c), 0) <= 0 || cmp.Compare((a*r.max.x+b*r.min.y+c)*(a*r.min.x+b*r.max.y+c), 0) <= 0 {
		return true
	}

	return false
}

// Dx returns r's width.
func (r Rectangle[T]) Dx() T { return r.max.x - r.min.x }

// Dy returns r's height.
func (r Rectangle[T]) Dy() T { return r.max.y - r.min.y }

// Size returns r's width and height.
func (r Rectangle[T]) Size() Point[T] { return Point[T]{r.max.x - r.min.x, r.max.y - r.min.y} }

// Add returns the rectangle r translated by p.
func (r Rectangle[T]) Add(p Point[T]) Rectangle[T] {
	return Rectangle[T]{Point[T]{r.min.x + p.x, r.min.y + p.y}, Point[T]{r.max.x + p.x, r.max.y + p.y}}
}

// Sub returns the rectangle r translated by -p.
func (r Rectangle[T]) Sub(p Point[T]) Rectangle[T] {
	return Rectangle[T]{Point[T]{r.min.x - p.x, r.min.y - p.y}, Point[T]{r.max.x - p.x, r.max.y - p.y}}
}

// Inset returns the rectangle r inset by n, which may be negative. If either
// of r's dimensions is less than 2*n then an empty rectangle near the center
// of r will be returned.
func (r Rectangle[T]) Inset(n T) Rectangle[T] {
	if cmp.Compare(r.Dx(), 2*n) < 0 {
		r.min.x = (r.min.x + r.max.x) / 2
		r.max.x = r.min.x
	} else {
		r.min.x += n
		r.max.x -= n
	}

	if cmp.Compare(r.Dy(), 2*n) < 0 {
		r.min.y = (r.min.y + r.max.y) / 2
		r.max.y = r.min.y
	} else {
		r.min.y += n
		r.max.y -= n
	}

	return r
}

// IsIntersect reports whether the rectangle intersect other rectangle.
func (r Rectangle[T]) IsIntersect(s Rectangle[T]) bool { return !r.Intersect(s).Empty() }

// Intersect returns the largest rectangle contained by both r and s. If the
// two rectangles do not overlap then the zero rectangle will be returned.
func (r Rectangle[T]) Intersect(s Rectangle[T]) Rectangle[T] {
	if cmp.Compare(r.min.x, s.min.x) < 0 {
		r.min.x = s.min.x
	}
	if cmp.Compare(r.min.y, s.min.y) < 0 {
		r.min.y = s.min.y
	}
	if cmp.Compare(r.max.x, s.max.x) > 0 {
		r.max.x = s.max.x
	}
	if cmp.Compare(r.max.y, s.max.y) > 0 {
		r.max.y = s.max.y
	}
	// Letting r0 and s0 be the values of r and s at the time that the method
	// is called, this next line is equivalent to:
	//
	// if max(r0.Min.X, s0.Min.X) >= min(r0.Max.X, s0.Max.X) || likewiseForY { etc }
	if r.Empty() {
		var ZR Rectangle[T]
		return ZR
	}

	return r
}

// Union returns the smallest rectangle that contains both r and s.
func (r Rectangle[T]) Union(s Rectangle[T]) Rectangle[T] {
	if r.Empty() {
		return s
	}

	if s.Empty() {
		return r
	}

	if cmp.Compare(r.min.x, s.min.x) > 0 {
		r.min.x = s.min.x
	}
	if cmp.Compare(r.min.y, s.min.y) > 0 {
		r.min.y = s.min.y
	}
	if cmp.Compare(r.max.x, s.max.x) < 0 {
		r.max.x = s.max.x
	}
	if cmp.Compare(r.max.y, s.max.y) < 0 {
		r.max.y = s.max.y
	}

	return r
}

// Has reports whether the rectangle contains the point.
func (r Rectangle[T]) Has(p Point[T]) bool {
	return cmp.Compare(p.x, r.min.x) >= 0 && cmp.Compare(p.x, r.max.x) <= 0 && cmp.Compare(p.y, r.min.y) >= 0 && cmp.Compare(p.y, r.max.y) <= 0
}

// Empty reports whether the rectangle contains no points.
func (r Rectangle[T]) Empty() bool {
	return cmp.Compare(r.min.x, r.max.x) >= 0 || cmp.Compare(r.min.y, r.max.y) >= 0
}

// Equals reports whether r and s contain the same set of points. All empty
// rectangles are considered equal.
func (r Rectangle[T]) Equals(s Rectangle[T]) bool { return r == s || r.Empty() && s.Empty() }

// Overlaps reports whether r and s have a non-empty intersection.
func (r Rectangle[T]) Overlaps(s Rectangle[T]) bool {
	return !r.Empty() && !s.Empty() &&
		cmp.Compare(r.min.x, s.max.x) < 0 && cmp.Compare(s.min.x, r.max.x) < 0 &&
		cmp.Compare(r.min.y, s.max.y) < 0 && cmp.Compare(s.min.y, r.max.y) < 0
}

// In reports whether every point in r is in s.
func (r Rectangle[T]) In(s Rectangle[T]) bool {
	if r.Empty() {
		return true
	}

	// Note that r.max is an exclusive bound for r, so that r.In(s)
	// does not require that r.max.In(s).
	return cmp.Compare(s.min.x, r.min.x) <= 0 && cmp.Compare(r.max.x, s.max.x) <= 0 &&
		cmp.Compare(s.min.y, r.min.y) <= 0 && cmp.Compare(r.max.y, s.max.y) <= 0
}

// Bounds returns a rectangle bounds
func (r Rectangle[T]) Bounds() Rectangle[T] { return r }

// Expanded returns a rectangle that has been expanded in the x-direction
// by margin.X, and in y-direction by margin.Y. The resulting rectangle may be empty.
func (r Rectangle[T]) Expanded(margin Point[T]) Rectangle[T] {
	return Rectangle[T]{
		Point[T]{r.min.x - margin.x, r.min.y - margin.y},
		Point[T]{r.max.x + margin.x, r.max.y + margin.y},
	}
}

// ExpandedByMargin returns a rectangle that has been expanded in the x-direction
// by margin, and in y-direction by margin. The resulting rectangle may be empty.
func (r Rectangle[T]) ExpandedByMargin(margin T) Rectangle[T] {
	return r.Expanded(P(margin, margin))
}

// Rect is shorthand for RectangleInt{Pt(x0, y0), Pt(x1, y1)}. The returned
// rectangle has minimum and maximum coordinates swapped if necessary so that
// it is well-formed.
func Rect[T Number](x0, y0, x1, y1 T) Rectangle[T] {
	if cmp.Compare(x0, x1) > 0 {
		x0, x1 = x1, x0
	}
	if cmp.Compare(y0, y1) > 0 {
		y0, y1 = y1, y0
	}

	return Rectangle[T]{Point[T]{x0, y0}, Point[T]{x1, y1}}
}

// RectFromCenterSize constructs a rectangle with the given center and size.
// Both dimensions of size must be non-negative.
func RectFromCenterSize[T Number](center, size Point[T]) Rectangle[T] {
	return Rect(center.x-size.x, center.y-size.y, center.x+size.x, center.y+size.y)
}

// RectFromMinSize constructs a rectangle with the given BottomLeft and size.
// Both dimensions of size must be non-negative.
func RectFromMinSize[T Number](min, size Point[T]) Rectangle[T] {
	return Rect(min.x, min.y, min.x+size.x, min.y+size.y)
}

// HelixRectRangeFromCenterAndMargin 由center节点逆时针螺旋由内向外访问margin区域内的所有节点
// 25   24   23   22   21
// 10    9    8    7   20
// 11    2    1    6   19
// 12    3    4    5   18
// 13   14   15   16   17
func HelixRectRangeFromCenterAndMargin[T Number](center Point[T], margin T, with func(p Point[T]) bool) {
	var x, y, xNow, yNow T
	xLen, yLen := T(1), T(1)
	rectXYLen := margin*2 + 1
	m, maxM := T(0), rectXYLen*rectXYLen
	startX, startY := center.x, center.y

	for cmp.Compare(m, maxM) <= 0 {
		for x = startX; cmp.Compare(x, startX-xLen) >= 0; x-- {
			m++
			if cmp.Compare(m, maxM) > 0 || !with(P(x, startY)) {
				return
			}
		}
		for y = startY - 1; cmp.Compare(y, startY-yLen) >= 0; y-- {
			m++
			if cmp.Compare(m, maxM) > 0 || !with(P(x+1, y)) {
				return
			}
		}
		xLen++
		yLen++
		for xNow = x + 2; cmp.Compare(xNow, x+xLen) <= 0; xNow++ {
			m++
			if cmp.Compare(m, maxM) > 0 || !with(P(xNow, y+1)) {
				return
			}
		}
		for yNow = y + 1; cmp.Compare(yNow, y+yLen) <= 0; yNow++ {
			m++
			if cmp.Compare(m, maxM) > 0 || !with(P(xNow, yNow)) {
				return
			}
		}
		xLen++
		yLen++
		startX = xNow
		startY = yNow
	}
}
