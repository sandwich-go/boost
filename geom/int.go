// Code generated by tools. DO NOT EDIT.
package geom

import "strconv"

type PointInt struct {
	X, Y int
}

// String returns a string representation of p like "(3,4)".
func (p PointInt) String() string {
	return "(" + strconv.FormatInt(int64(p.X), 10) + "," + strconv.FormatInt(int64(p.Y), 10) + ")"
}

// Add returns the vector p+q.
func (p PointInt) Add(q PointInt) PointInt {
	return PointInt{p.X + q.X, p.Y + q.Y}
}

// Sub returns the vector p-q.
func (p PointInt) Sub(q PointInt) PointInt {
	return PointInt{p.X - q.X, p.Y - q.Y}
}

// Mul returns the vector p*k.
func (p PointInt) Mul(k int) PointInt {
	return PointInt{p.X * k, p.Y * k}
}

// Div returns the vector p/k.
func (p PointInt) Div(k int) PointInt {
	return PointInt{p.X / k, p.Y / k}
}

// In reports whether p is in r.
func (p PointInt) In(r RectangleInt) bool {
	return r.Min.X <= p.X && p.X <= r.Max.X && r.Min.Y <= p.Y && p.Y <= r.Max.Y
}

// Eq reports whether p and q are equal.
func (p PointInt) Eq(q PointInt) bool {
	return p == q
}

// ZPInt is the zero PointInt.
var ZPInt PointInt

// PtInt is shorthand for PointInt.{X, Y}.
func PtInt(X, Y int) PointInt {
	return PointInt{X, Y}
}

type RectangleInt struct {
	Min, Max PointInt
}

// String returns a string representation of r like "(3,4)-(6,5)".
func (r RectangleInt) String() string {
	return r.Min.String() + "-" + r.Max.String()
}

func (r RectangleInt) RangePoints(with func(p PointInt) bool) {
	if with == nil || r == ZRInt {
		return
	}
	for x := r.Min.X; x <= r.Max.X; x++ {
		for y := r.Min.Y; y <= r.Max.Y; y++ {
			if !with(PtInt(x, y)) {
				return
			}
		}
	}
}

func (r RectangleInt) RangePointsMinClosedMaxOpen(with func(p PointInt) bool) {
	if with == nil || r == ZRInt {
		return
	}
	for x := r.Min.X + 1; x <= r.Max.X; x++ {
		for y := r.Min.Y + 1; y <= r.Max.Y; y++ {
			if !with(PtInt(x, y)) {
				return
			}
		}
	}
}

func (r RectangleInt) IntersectionWithLine(s PointInt, e PointInt) bool {
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

func (r RectangleInt) RangePointsMinMaxClosed(with func(p PointInt) bool) {
	if with == nil || r == ZRInt {
		return
	}
	for x := r.Min.X + 1; x < r.Max.X; x++ {
		for y := r.Min.Y + 1; y < r.Max.Y; y++ {
			if !with(PtInt(x, y)) {
				return
			}
		}
	}
}

func (r RectangleInt) RangePointsMinOpenMaxClosed(with func(p PointInt) bool) {
	if with == nil || r == ZRInt {
		return
	}
	for x := r.Min.X; x < r.Max.X; x++ {
		for y := r.Min.Y; y < r.Max.Y; y++ {
			if !with(PtInt(x, y)) {
				return
			}
		}
	}
}

// Dx returns r's width.
func (r RectangleInt) Dx() int {
	return r.Max.X - r.Min.X
}

// Dy returns r's height.
func (r RectangleInt) Dy() int {
	return r.Max.Y - r.Min.Y
}

// Size returns r's width and height.
func (r RectangleInt) Size() PointInt {
	return PointInt{
		r.Max.X - r.Min.X,
		r.Max.Y - r.Min.Y,
	}
}

// Add returns the rectangle r translated by p.
func (r RectangleInt) Add(p PointInt) RectangleInt {
	return RectangleInt{
		PointInt{r.Min.X + p.X, r.Min.Y + p.Y},
		PointInt{r.Max.X + p.X, r.Max.Y + p.Y},
	}
}

// Sub returns the rectangle r translated by -p.
func (r RectangleInt) Sub(p PointInt) RectangleInt {
	return RectangleInt{
		PointInt{r.Min.X - p.X, r.Min.Y - p.Y},
		PointInt{r.Max.X - p.X, r.Max.Y - p.Y},
	}
}

// Inset returns the rectangle r inset by n, which may be negative. If either
// of r's dimensions is less than 2*n then an empty rectangle near the center
// of r will be returned.
func (r RectangleInt) Inset(n int) RectangleInt {
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
func (r RectangleInt) Intersect(s RectangleInt) RectangleInt {
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
		return ZRInt
	}
	return r
}

// Union returns the smallest rectangle that contains both r and s.
func (r RectangleInt) Union(s RectangleInt) RectangleInt {
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
func (r RectangleInt) Empty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

// Eq reports whether r and s contain the same set of points. All empty
// rectangles are considered equal.
func (r RectangleInt) Eq(s RectangleInt) bool {
	return r == s || r.Empty() && s.Empty()
}

// Overlaps reports whether r and s have a non-empty intersection.
func (r RectangleInt) Overlaps(s RectangleInt) bool {
	return !r.Empty() && !s.Empty() &&
		r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

// In reports whether every point in r is in s.
func (r RectangleInt) In(s RectangleInt) bool {
	if r.Empty() {
		return true
	}
	// Note that r.Max is an exclusive bound for r, so that r.In(s)
	// does not require that r.Max.In(s).
	return s.Min.X <= r.Min.X && r.Max.X <= s.Max.X &&
		s.Min.Y <= r.Min.Y && r.Max.Y <= s.Max.Y
}

func (r RectangleInt) Bounds() RectangleInt {
	return r
}

// Expanded returns a rectangle that has been expanded in the x-direction
// by margin.X, and in y-direction by margin.Y. The resulting rectangle may be empty.
func (r RectangleInt) Expanded(margin PointInt) RectangleInt {
	return RectangleInt{
		PointInt{r.Min.X - margin.X, r.Min.Y - margin.Y},
		PointInt{r.Max.X + margin.X, r.Max.Y + margin.Y},
	}
}

func (r RectangleInt) ExpandedByMargin(margin int) RectangleInt {
	return r.Expanded(PtInt(margin, margin))
}

// ZRInt is the zero RectangleInt.
var ZRInt RectangleInt

// RectInt is shorthand for RectangleInt{Pt(x0, y0), Pt(x1, y1)}. The returned
// rectangle has minimum and maximum coordinates swapped if necessary so that
// it is well-formed.
func RectInt(x0, y0, x1, y1 int) RectangleInt {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return RectangleInt{PointInt{x0, y0}, PointInt{x1, y1}}
}

// RectIntFromCenterSize constructs a rectangle with the given center and size.
// Both dimensions of size must be non-negative.
func RectIntFromCenterSize(center, size PointInt) RectangleInt {
	return RectInt(center.X-size.X, center.Y-center.Y, center.X+size.X, center.Y+size.Y)
}

// HelixRectRangeFromCenterAndMarginInt 由center节点逆时针螺旋由内向外访问margin区域内的所有节点
// 25   24   23   22   21
// 10    9    8    7   20
// 11    2    1    6   19
// 12    3    4    5   18
// 13   14   15   16   17
func HelixRectRangeFromCenterAndMarginInt(center PointInt, margin int, with func(p PointInt) bool) {
	var x, y, xNow, yNow int
	xLen, yLen := int(1), int(1)
	rectXYLen := margin*2 + 1
	m, max := int(0), rectXYLen*rectXYLen
	startX, startY := center.X, center.Y

	for m <= max {
		for x = startX; x >= startX-xLen; x-- {
			m++
			if m > max || !with(PtInt(x, startY)) {
				return
			}
		}
		for y = startY - 1; y >= startY-yLen; y-- {
			m++
			if m > max || !with(PtInt(x+1, y)) {
				return
			}
		}
		xLen++
		yLen++
		for xNow = x + 2; xNow <= x+xLen; xNow++ {
			m++
			if m > max || !with(PtInt(xNow, y+1)) {
				return
			}
		}
		for yNow = y + 1; yNow <= y+yLen; yNow++ {
			m++
			if m > max || !with(PtInt(xNow, yNow)) {
				return
			}
		}
		xLen++
		yLen++
		startX = xNow
		startY = yNow
	}
}
