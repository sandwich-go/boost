package misc

import (
	"testing"
)

type InterfaceSample interface {
	doHaha() int
}

type AbstractSample[T any, PT interface {
	*T
	InterfaceSample
}] struct {
	Extend[T, PT]
}

func (a *AbstractSample[T, PT]) Haha() int {
	return a.Self().doHaha()
}

func (a *AbstractSample[T, PT]) doHaha() int { return 999 }

type Sample1 struct {
	AbstractSample[Sample1, *Sample1]
}

func (s *Sample1) doHaha() int {
	return 0
}

type Sample2 struct {
	AbstractSample[Sample2, *Sample2]
}

func (s *Sample2) doHaha() int {
	return 1
}

type Sample3 struct {
	AbstractSample[Sample3, *Sample3]
}

func TestExtend(t *testing.T) {
	sample1 := &Sample1{}
	if sample1.Haha() != 0 {
		t.Fail()
	}
	sample2 := &Sample2{}
	if sample2.Haha() != 1 {
		t.Fail()
	}
	sample3 := &Sample3{}
	if sample3.Haha() != 999 {
		t.Fail()
	}
}
