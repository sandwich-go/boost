package graph

import (
	"cmp"
	"fmt"
	"github.com/sandwich-go/boost/xconv"
	"golang.org/x/exp/constraints"
	"math"
)

// Number 支持的坐标类型
type Number interface {
	constraints.Integer | constraints.Float
}

// Point 坐标点
type Point[T Number] struct {
	x, y T
}

// P return shorthand for Point[T].{x, y}.
func P[T Number](x, y T) Point[T] { return Point[T]{x: x, y: y} }

// X return shorthand for Point[T].x
func (p Point[T]) X() T { return p.x }

// Y return shorthand for Point[T].y
func (p Point[T]) Y() T { return p.y }

// String returns a string representation of p like "(3,4)".
func (p Point[T]) String() string {
	return fmt.Sprintf("(%s,%s)", xconv.String(p.x), xconv.String(p.y))
}

// Add returns the vector p+q.
func (p Point[T]) Add(q Point[T]) Point[T] {
	return Point[T]{p.x + q.x, p.y + q.y}
}

// Sub returns the vector p-q.
func (p Point[T]) Sub(q Point[T]) Point[T] {
	return Point[T]{p.x - q.x, p.y - q.y}
}

// Mul returns the vector p*k.
func (p Point[T]) Mul(k T) Point[T] {
	return Point[T]{p.x * k, p.y * k}
}

// Div returns the vector p/k.
func (p Point[T]) Div(k T) Point[T] {
	return Point[T]{p.x / k, p.y / k}
}

// In reports whether p is in r.
func (p Point[T]) In(r Rectangle[T]) bool {
	return cmp.Compare(r.min.x, p.x) <= 0 && cmp.Compare(p.x, r.max.x) <= 0 && cmp.Compare(r.min.y, p.y) <= 0 && cmp.Compare(p.y, r.max.y) <= 0
}

// Equals reports whether p and q are equal.
func (p Point[T]) Equals(q Point[T]) bool { return p == q }

// Distance returns distance p and q.
func (p Point[T]) Distance(q Point[T]) float64 {
	if cmp.Compare(p.x, q.x) == 0 {
		return math.Abs(float64(p.y - q.y))
	}

	if cmp.Compare(p.y, q.y) == 0 {
		return math.Abs(float64(p.x - q.x))
	}

	return math.Sqrt(math.Pow(math.Abs(float64(p.x-q.x)), 2) + math.Pow(math.Abs(float64(p.y-q.y)), 2))
}
