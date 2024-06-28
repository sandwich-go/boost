package geom

import (
	"golang.org/x/exp/constraints"
	"strconv"
)

type Point[T constraints.Integer] struct {
	X, Y T
}

// String returns a string representation of p like "(3,4)".
func (p Point[T]) String() string {
	return "(" + strconv.FormatUint(uint64(p.X), 10) + "," + strconv.FormatUint(uint64(p.Y), 10) + ")"
}

// Add returns the vector p+q.
func (p Point[T]) Add(q Point[T]) Point[T] {
	return Point[T]{p.X + q.X, p.Y + q.Y}
}

// Sub returns the vector p-q.
func (p Point[T]) Sub(q Point[T]) Point[T] {
	return Point[T]{p.X - q.X, p.Y - q.Y}
}

// Mul returns the vector p*k.
func (p Point[T]) Mul(k T) Point[T] {
	return Point[T]{p.X * k, p.Y * k}
}

// Div returns the vector p/k.
func (p Point[T]) Div(k T) Point[T] {
	return Point[T]{p.X / k, p.Y / k}
}

// In reports whether p is in r.
func (p Point[T]) In(r Rectangle[T]) bool {
	return r.Min.X <= p.X && p.X <= r.Max.X && r.Min.Y <= p.Y && p.Y <= r.Max.Y
}

// Eq reports whether p and q are equal.
func (p Point[T]) Eq(q Point[T]) bool {
	return p == q
}

// Pt return shorthand for Point[T].{X, Y}.
func Pt[T constraints.Integer](X, Y T) Point[T] { return Point[T]{X, Y} }

// Rectangle wrapper, contain two Point[T].
type Rectangle[T constraints.Integer] struct {
	Min, Max Point[T]
}

// String returns a string representation of r like "(3,4)-(6,5)".
func (r Rectangle[T]) String() string {
	return r.Min.String() + "-" + r.Max.String()
}

// RangePoints range all points in rectangle.
// if with return false, aborted range.
func (r Rectangle[T]) RangePoints(with func(p Point[T]) bool) {
	var ZR Rectangle[T]
	if with == nil || r == ZR {
		return
	}
	for x := r.Min.X; x <= r.Max.X; x++ {
		for y := r.Min.Y; y <= r.Max.Y; y++ {
			if !with(Pt(x, y)) {
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
	for x := r.Min.X + 1; x <= r.Max.X; x++ {
		for y := r.Min.Y + 1; y <= r.Max.Y; y++ {
			if !with(Pt(x, y)) {
				return
			}
		}
	}
}

// IntersectionWithLine check line intersection, if intersection, return true
func (r Rectangle[T]) IntersectionWithLine(s Point[T], e Point[T]) bool {
	if s.X <= r.Min.X && e.X <= r.Min.X || s.X >= r.Max.X && e.X >= r.Max.X || s.Y <= r.Min.Y && e.Y <= r.Min.Y || s.Y >= r.Max.Y && e.Y >= r.Max.Y {
		return false
	}
	a := s.Y - e.Y
	b := e.X - s.X
	c := e.Y*s.X - e.X*s.Y
	if ((a*r.Min.X+b*r.Min.Y+c)*(a*r.Max.X+b*r.Max.Y+c)) <= 0 || ((a*r.Max.X+b*r.Min.Y+c)*(a*r.Min.X+b*r.Max.Y+c)) <= 0 {
		return true
	}
	return false
}

// RangePointsMinMaxClosed range all points in rectangle except min/max x, y.
// if with return false, aborted range.
func (r Rectangle[T]) RangePointsMinMaxClosed(with func(p Point[T]) bool) {
	var ZR Rectangle[T]
	if with == nil || r == ZR {
		return
	}
	for x := r.Min.X + 1; x < r.Max.X; x++ {
		for y := r.Min.Y + 1; y < r.Max.Y; y++ {
			if !with(Pt(x, y)) {
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
	for x := r.Min.X; x < r.Max.X; x++ {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			if !with(Pt(x, y)) {
				return
			}
		}
	}
}

// Dx returns r's width.
func (r Rectangle[T]) Dx() T {
	return r.Max.X - r.Min.X
}

// Dy returns r's height.
func (r Rectangle[T]) Dy() T {
	return r.Max.Y - r.Min.Y
}

// Size returns r's width and height.
func (r Rectangle[T]) Size() Point[T] {
	return Point[T]{
		r.Max.X - r.Min.X,
		r.Max.Y - r.Min.Y,
	}
}

// Add returns the rectangle r translated by p.
func (r Rectangle[T]) Add(p Point[T]) Rectangle[T] {
	return Rectangle[T]{
		Point[T]{r.Min.X + p.X, r.Min.Y + p.Y},
		Point[T]{r.Max.X + p.X, r.Max.Y + p.Y},
	}
}

// Sub returns the rectangle r translated by -p.
func (r Rectangle[T]) Sub(p Point[T]) Rectangle[T] {
	return Rectangle[T]{
		Point[T]{r.Min.X - p.X, r.Min.Y - p.Y},
		Point[T]{r.Max.X - p.X, r.Max.Y - p.Y},
	}
}

// Inset returns the rectangle r inset by n, which may be negative. If either
// of r's dimensions is less than 2*n then an empty rectangle near the center
// of r will be returned.
func (r Rectangle[T]) Inset(n T) Rectangle[T] {
	if r.Dx() < 2*n {
		r.Min.X = (r.Min.X + r.Max.X) / 2
		r.Max.X = r.Min.X
	} else {
		r.Min.X += n
		r.Max.X -= n
	}
	if r.Dy() < 2*n {
		r.Min.Y = (r.Min.Y + r.Max.Y) / 2
		r.Max.Y = r.Min.Y
	} else {
		r.Min.Y += n
		r.Max.Y -= n
	}
	return r
}

// Intersect returns the largest rectangle contained by both r and s. If the
// two rectangles do not overlap then the zero rectangle will be returned.
func (r Rectangle[T]) Intersect(s Rectangle[T]) Rectangle[T] {
	if r.Min.X < s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y < s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X > s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y > s.Max.Y {
		r.Max.Y = s.Max.Y
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
	if r.Min.X > s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y > s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X < s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y < s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	return r
}

// Empty reports whether the rectangle contains no points.
func (r Rectangle[T]) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

// Eq reports whether r and s contain the same set of points. All empty
// rectangles are considered equal.
func (r Rectangle[T]) Eq(s Rectangle[T]) bool {
	return r == s || r.Empty() && s.Empty()
}

// Overlaps reports whether r and s have a non-empty intersection.
func (r Rectangle[T]) Overlaps(s Rectangle[T]) bool {
	return !r.Empty() && !s.Empty() &&
		r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

// In reports whether every point in r is in s.
func (r Rectangle[T]) In(s Rectangle[T]) bool {
	if r.Empty() {
		return true
	}
	// Note that r.Max is an exclusive bound for r, so that r.In(s)
	// does not require that r.Max.In(s).
	return s.Min.X <= r.Min.X && r.Max.X <= s.Max.X &&
		s.Min.Y <= r.Min.Y && r.Max.Y <= s.Max.Y
}

// Bounds returns a rectangle bounds
func (r Rectangle[T]) Bounds() Rectangle[T] {
	return r
}

// Expanded returns a rectangle that has been expanded in the x-direction
// by margin.X, and in y-direction by margin.Y. The resulting rectangle may be empty.
func (r Rectangle[T]) Expanded(margin Point[T]) Rectangle[T] {
	return Rectangle[T]{
		Point[T]{r.Min.X - margin.X, r.Min.Y - margin.Y},
		Point[T]{r.Max.X + margin.X, r.Max.Y + margin.Y},
	}
}

// ExpandedByMargin returns a rectangle that has been expanded in the x-direction
// by margin, and in y-direction by margin. The resulting rectangle may be empty.
func (r Rectangle[T]) ExpandedByMargin(margin T) Rectangle[T] {
	return r.Expanded(Pt(margin, margin))
}

// Rect is shorthand for RectangleInt{Pt(x0, y0), Pt(x1, y1)}. The returned
// rectangle has minimum and maximum coordinates swapped if necessary so that
// it is well-formed.
func Rect[T constraints.Integer](x0, y0, x1, y1 T) Rectangle[T] {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return Rectangle[T]{Point[T]{x0, y0}, Point[T]{x1, y1}}
}

// RectFromCenterSize constructs a rectangle with the given center and size.
// Both dimensions of size must be non-negative.
func RectFromCenterSize[T constraints.Integer](center, size Point[T]) Rectangle[T] {
	return Rect(center.X-size.X, center.Y-center.Y, center.X+size.X, center.Y+size.Y)
}

// HelixRectRangeFromCenterAndMargin 由center节点逆时针螺旋由内向外访问margin区域内的所有节点
// 25   24   23   22   21
// 10    9    8    7   20
// 11    2    1    6   19
// 12    3    4    5   18
// 13   14   15   16   17
func HelixRectRangeFromCenterAndMargin[T constraints.Integer](center Point[T], margin T, with func(p Point[T]) bool) {
	var x, y, xNow, yNow T
	xLen, yLen := T(1), T(1)
	rectXYLen := margin*2 + 1
	m, maxM := T(0), rectXYLen*rectXYLen
	startX, startY := center.X, center.Y

	for m <= maxM {
		for x = startX; x >= startX-xLen; x-- {
			m++
			if m > maxM || !with(Pt(x, startY)) {
				return
			}
		}
		for y = startY - 1; y >= startY-yLen; y-- {
			m++
			if m > maxM || !with(Pt(x+1, y)) {
				return
			}
		}
		xLen++
		yLen++
		for xNow = x + 2; xNow <= x+xLen; xNow++ {
			m++
			if m > maxM || !with(Pt(xNow, y+1)) {
				return
			}
		}
		for yNow = y + 1; yNow <= y+yLen; yNow++ {
			m++
			if m > maxM || !with(Pt(xNow, yNow)) {
				return
			}
		}
		xLen++
		yLen++
		startX = xNow
		startY = yNow
	}
}
