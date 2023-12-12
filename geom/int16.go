// Code generated by tools. DO NOT EDIT.
package geom

import "strconv"

// PointInt16 point wrapper
type PointInt16 struct {
	X, Y int16
}

// String returns a string representation of p like "(3,4)".
func (p PointInt16) String() string {
	return "(" + strconv.FormatInt(int64(p.X), 10) + "," + strconv.FormatInt(int64(p.Y), 10) + ")"
}

// Add returns the vector p+q.
func (p PointInt16) Add(q PointInt16) PointInt16 {
	return PointInt16{p.X + q.X, p.Y + q.Y}
}

// Sub returns the vector p-q.
func (p PointInt16) Sub(q PointInt16) PointInt16 {
	return PointInt16{p.X - q.X, p.Y - q.Y}
}

// Mul returns the vector p*k.
func (p PointInt16) Mul(k int16) PointInt16 {
	return PointInt16{p.X * k, p.Y * k}
}

// Div returns the vector p/k.
func (p PointInt16) Div(k int16) PointInt16 {
	return PointInt16{p.X / k, p.Y / k}
}

// In reports whether p is in r.
func (p PointInt16) In(r RectangleInt16) bool {
	return r.Min.X <= p.X && p.X <= r.Max.X && r.Min.Y <= p.Y && p.Y <= r.Max.Y
}

// Eq reports whether p and q are equal.
func (p PointInt16) Eq(q PointInt16) bool {
	return p == q
}

// ZPInt16 is the zero PointInt16.
var ZPInt16 PointInt16

// PtInt16 is shorthand for PointInt16.{X, Y}.
func PtInt16(X, Y int16) PointInt16 {
	return PointInt16{X, Y}
}

// RectangleInt16 rectangle wrapper, contain two point.
type RectangleInt16 struct {
	Min, Max PointInt16
}

// String returns a string representation of r like "(3,4)-(6,5)".
func (r RectangleInt16) String() string {
	return r.Min.String() + "-" + r.Max.String()
}

// RangePoints range all points in rectangle.
// if with return false, aborted range.
func (r RectangleInt16) RangePoints(with func(p PointInt16) bool) {
	if with == nil || r == ZRInt16 {
		return
	}
	for x := r.Min.X; x <= r.Max.X; x++ {
		for y := r.Min.Y; y <= r.Max.Y; y++ {
			if !with(PtInt16(x, y)) {
				return
			}
		}
	}
}

// RangePointsMinClosedMaxOpen range all points in rectangle except min x, y.
// if with return false, aborted range.
func (r RectangleInt16) RangePointsMinClosedMaxOpen(with func(p PointInt16) bool) {
	if with == nil || r == ZRInt16 {
		return
	}
	for x := r.Min.X + 1; x <= r.Max.X; x++ {
		for y := r.Min.Y + 1; y <= r.Max.Y; y++ {
			if !with(PtInt16(x, y)) {
				return
			}
		}
	}
}

// IntersectionWithLine check line intersection, if intersection, return true
func (r RectangleInt16) IntersectionWithLine(s PointInt16, e PointInt16) bool {
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
func (r RectangleInt16) RangePointsMinMaxClosed(with func(p PointInt16) bool) {
	if with == nil || r == ZRInt16 {
		return
	}
	for x := r.Min.X + 1; x < r.Max.X; x++ {
		for y := r.Min.Y + 1; y < r.Max.Y; y++ {
			if !with(PtInt16(x, y)) {
				return
			}
		}
	}
}

// RangePointsMinOpenMaxClosed range all points in rectangle except max x, y.
// if with return false, aborted range.
func (r RectangleInt16) RangePointsMinOpenMaxClosed(with func(p PointInt16) bool) {
	if with == nil || r == ZRInt16 {
		return
	}
	for x := r.Min.X; x < r.Max.X; x++ {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			if !with(PtInt16(x, y)) {
				return
			}
		}
	}
}

// Dx returns r's width.
func (r RectangleInt16) Dx() int16 {
	return r.Max.X - r.Min.X
}

// Dy returns r's height.
func (r RectangleInt16) Dy() int16 {
	return r.Max.Y - r.Min.Y
}

// Size returns r's width and height.
func (r RectangleInt16) Size() PointInt16 {
	return PointInt16{
		r.Max.X - r.Min.X,
		r.Max.Y - r.Min.Y,
	}
}

// Add returns the rectangle r translated by p.
func (r RectangleInt16) Add(p PointInt16) RectangleInt16 {
	return RectangleInt16{
		PointInt16{r.Min.X + p.X, r.Min.Y + p.Y},
		PointInt16{r.Max.X + p.X, r.Max.Y + p.Y},
	}
}

// Sub returns the rectangle r translated by -p.
func (r RectangleInt16) Sub(p PointInt16) RectangleInt16 {
	return RectangleInt16{
		PointInt16{r.Min.X - p.X, r.Min.Y - p.Y},
		PointInt16{r.Max.X - p.X, r.Max.Y - p.Y},
	}
}

// Inset returns the rectangle r inset by n, which may be negative. If either
// of r's dimensions is less than 2*n then an empty rectangle near the center
// of r will be returned.
func (r RectangleInt16) Inset(n int16) RectangleInt16 {
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
func (r RectangleInt16) Intersect(s RectangleInt16) RectangleInt16 {
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
		return ZRInt16
	}
	return r
}

// Union returns the smallest rectangle that contains both r and s.
func (r RectangleInt16) Union(s RectangleInt16) RectangleInt16 {
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
func (r RectangleInt16) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

// Eq reports whether r and s contain the same set of points. All empty
// rectangles are considered equal.
func (r RectangleInt16) Eq(s RectangleInt16) bool {
	return r == s || r.Empty() && s.Empty()
}

// Overlaps reports whether r and s have a non-empty intersection.
func (r RectangleInt16) Overlaps(s RectangleInt16) bool {
	return !r.Empty() && !s.Empty() &&
		r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

// In reports whether every point in r is in s.
func (r RectangleInt16) In(s RectangleInt16) bool {
	if r.Empty() {
		return true
	}
	// Note that r.Max is an exclusive bound for r, so that r.In(s)
	// does not require that r.Max.In(s).
	return s.Min.X <= r.Min.X && r.Max.X <= s.Max.X &&
		s.Min.Y <= r.Min.Y && r.Max.Y <= s.Max.Y
}

// Bounds returns a rectangle bounds
func (r RectangleInt16) Bounds() RectangleInt16 {
	return r
}

// Expanded returns a rectangle that has been expanded in the x-direction
// by margin.X, and in y-direction by margin.Y. The resulting rectangle may be empty.
func (r RectangleInt16) Expanded(margin PointInt16) RectangleInt16 {
	return RectangleInt16{
		PointInt16{r.Min.X - margin.X, r.Min.Y - margin.Y},
		PointInt16{r.Max.X + margin.X, r.Max.Y + margin.Y},
	}
}

// ExpandedByMargin returns a rectangle that has been expanded in the x-direction
// by margin, and in y-direction by margin. The resulting rectangle may be empty.
func (r RectangleInt16) ExpandedByMargin(margin int16) RectangleInt16 {
	return r.Expanded(PtInt16(margin, margin))
}

// ZRInt16 is the zero RectangleInt16.
var ZRInt16 RectangleInt16

// RectInt16 is shorthand for RectangleInt16{Pt(x0, y0), Pt(x1, y1)}. The returned
// rectangle has minimum and maximum coordinates swapped if necessary so that
// it is well-formed.
func RectInt16(x0, y0, x1, y1 int16) RectangleInt16 {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return RectangleInt16{PointInt16{x0, y0}, PointInt16{x1, y1}}
}

// RectInt16FromCenterSize constructs a rectangle with the given center and size.
// Both dimensions of size must be non-negative.
func RectInt16FromCenterSize(center, size PointInt16) RectangleInt16 {
	return RectInt16(center.X-size.X, center.Y-center.Y, center.X+size.X, center.Y+size.Y)
}

// HelixRectRangeFromCenterAndMarginInt16 由center节点逆时针螺旋由内向外访问margin区域内的所有节点
// 25   24   23   22   21
// 10    9    8    7   20
// 11    2    1    6   19
// 12    3    4    5   18
// 13   14   15   16   17
func HelixRectRangeFromCenterAndMarginInt16(center PointInt16, margin int16, with func(p PointInt16) bool) {
	var x, y, xNow, yNow int16
	xLen, yLen := int16(1), int16(1)
	rectXYLen := margin*2 + 1
	m, max := int16(0), rectXYLen*rectXYLen
	startX, startY := center.X, center.Y

	for m <= max {
		for x = startX; x >= startX-xLen; x-- {
			m++
			if m > max || !with(PtInt16(x, startY)) {
				return
			}
		}
		for y = startY - 1; y >= startY-yLen; y-- {
			m++
			if m > max || !with(PtInt16(x+1, y)) {
				return
			}
		}
		xLen++
		yLen++
		for xNow = x + 2; xNow <= x+xLen; xNow++ {
			m++
			if m > max || !with(PtInt16(xNow, y+1)) {
				return
			}
		}
		for yNow = y + 1; yNow <= y+yLen; yNow++ {
			m++
			if m > max || !with(PtInt16(xNow, yNow)) {
				return
			}
		}
		xLen++
		yLen++
		startX = xNow
		startY = yNow
	}
}