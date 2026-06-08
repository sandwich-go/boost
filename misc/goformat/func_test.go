package goformat

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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

	Convey("remove bare returns（命名返回值）", t, func() {
		// removeBareReturns 用 ast.Ident 直接复用命名返回值的 name，**不调
		// newZeroValueNode**。这里只验证基础 path（其他类型见下面 fill
		// returns 用例）。
		for _, test := range []struct {
			src string
		}{
			{src: `package c
func b() (err error) { return  }`},
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

	Convey("fill return values 触发 newZeroValueNode 各 type case", t, func() {
		// fillReturnValues 触发条件：numRVs > 0 && numRVs < ftyp.Results 且
		// ret.Results[0] 是 single-return CallExpr。这里 caller 返多值，
		// callee 返单值，让 fillReturnValues 给前面补零值，正好走
		// newZeroValueNode 各 case。
		for _, src := range []string{
			// int 系列（line 81-82）
			`package c
func a() (int, error) { return b() }
func b() (error) { return nil }`,
			// uint8 (byte)
			`package c
func a() (uint8, error) { return b() }
func b() (error) { return nil }`,
			// float64（line 83-84）
			`package c
func a() (float64, error) { return b() }
func b() (error) { return nil }`,
			// complex128（line 85-86）
			`package c
func a() (complex128, error) { return b() }
func b() (error) { return nil }`,
			// bool（line 87-88）
			`package c
func a() (bool, error) { return b() }
func b() (error) { return nil }`,
			// string（line 89-90）
			`package c
func a() (string, error) { return b() }
func b() (error) { return nil }`,
			// error（line 91-92）
			`package c
func a() (error, error) { return b() }
func b() (error) { return nil }`,
			// slice（line 94-97 ArrayType v.Len==nil）
			`package c
func a() ([]int, error) { return b() }
func b() (error) { return nil }`,
			// array（line 99 ArrayType v.Len!=nil）
			`package c
func a() ([3]int, error) { return b() }
func b() (error) { return nil }`,
			// pointer（line 100-101 StarExpr）
			`package c
func a() (*int, error) { return b() }
func b() (error) { return nil }`,
			// 自定义类型 → 走 line 79 default → newZeroValueNode 返 nil →
			// fillReturnValues line 67 continue（不修改），覆盖该路径
			`package c
type T int
func a() (T, error) { return b() }
func b() (error) { return nil }`,
		} {
			f, err := parser.ParseFile(token.NewFileSet(), "", src, 0)
			So(err, ShouldBeNil)
			err = fillReturnValues(f)
			So(err, ShouldBeNil)
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
