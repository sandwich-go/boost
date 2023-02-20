package goformat

import (
	. "github.com/smartystreets/goconvey/convey"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestFunc(t *testing.T) {
	var tests = []struct {
		src                           string
		hasSingleReturnValueFuncCount int
	}{
		{
			src: `package c
		func a() (error) { return nil }`,
		},
		{
			src: `package c
		func a() (int, error) { return b() }
		func b() (int, error) { return 0, nil }`,
		},
		{
			src: `package c
		func a() (error) { return b() }
		func b() (error) { return nil }`,
			hasSingleReturnValueFuncCount: 1,
		},
		{
			src: `package c
		func d() (error) { return b() }
		func a() (error) { return b() }
		func b() (error) { return nil }`,
			hasSingleReturnValueFuncCount: 2,
		},
		{
			src: `package c
func a() (int, error) { return b() }
func b() (error) { return nil }`,
			hasSingleReturnValueFuncCount: 1,
		},
	}
	Convey("func return value count", t, func() {
		for _, test := range tests {
			f, err := parser.ParseFile(token.NewFileSet(), "", test.src, 0)
			So(err, ShouldBeNil)
			returns := map[*ast.ReturnStmt]*ast.FuncType{}
			ast.Walk(visitor{returns: returns}, f)

			var count int
			for ret := range returns {
				if len(ret.Results) == 0 {
					continue
				}
				if e, ok := ret.Results[0].(*ast.CallExpr); ok && hasSingleCallReturnVal(e) {
					count++
				}
			}
			So(count, ShouldEqual, test.hasSingleReturnValueFuncCount)
		}
	})

	Convey("fill return values", t, func() {
		for _, test := range tests {
			f, err := parser.ParseFile(token.NewFileSet(), "", test.src, 0)
			So(err, ShouldBeNil)
			err = fillReturnValues(f)
			So(err, ShouldBeNil)

			returns := map[*ast.ReturnStmt]*ast.FuncType{}
			ast.Walk(visitor{returns: returns}, f)

			for ret, ftRet := range returns {
				if len(ret.Results) == 0 || ftRet.Results == nil {
					continue
				}
				if e, ok := ret.Results[0].(*ast.CallExpr); ok {
					if !hasSingleCallReturnVal(e) {
						continue
					}
				}
				So(len(ret.Results), ShouldEqual, len(ftRet.Results.List))
			}
		}
	})

	Convey("remove bare returns", t, func() {
		for _, test := range []struct {
			src string
		}{
			{
				src: `package c
func b() (err error) { return  }`,
			},
		} {
			f, err := parser.ParseFile(token.NewFileSet(), "", test.src, 0)
			So(err, ShouldBeNil)
			err = removeBareReturns(f)
			So(err, ShouldBeNil)

			returns := map[*ast.ReturnStmt]*ast.FuncType{}
			ast.Walk(visitor{returns: returns}, f)

			for ret, ftRet := range returns {
				if ftRet.Results == nil {
					continue
				}
				So(len(ret.Results), ShouldEqual, len(ftRet.Results.List))
			}
		}
	})

	Convey("contains main func", t, func() {
		for _, test := range []struct {
			src      string
			contains bool
		}{
			{
				src: `package c
func b() (err error) { return  }`,
			},
			{
				src: `package c
func main() (err error) { return  }`,
			},
			{
				src: `package c
func main() { return  }`, contains: true,
			},
		} {
			f, err := parser.ParseFile(token.NewFileSet(), "", test.src, 0)
			So(err, ShouldBeNil)
			So(containsMainFunc(f), ShouldEqual, test.contains)
		}
	})
}
