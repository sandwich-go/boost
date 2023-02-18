package goformat

import (
	"go/ast"
	"go/token"
)

// hasSingleCallReturnVal call 函数是否有单个返回值
func hasSingleCallReturnVal(ce *ast.CallExpr) bool {
	if id, ok0 := ce.Fun.(*ast.Ident); ok0 && id.Obj != nil {
		if fn, ok1 := id.Obj.Decl.(*ast.FuncDecl); ok1 {
			return len(fn.Type.Results.List) == 1
		}
	}
	return false
}

type visitor struct {
	enclosing *ast.FuncType                     // innermost enclosing func
	returns   map[*ast.ReturnStmt]*ast.FuncType // potentially incomplete returns
}

func (v visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return v
	}
	switch n := node.(type) {
	case *ast.FuncDecl:
		return visitor{enclosing: n.Type, returns: v.returns}
	case *ast.FuncLit:
		return visitor{enclosing: n.Type, returns: v.returns}
	case *ast.ReturnStmt:
		v.returns[n] = v.enclosing
	}
	return v
}

// fillReturnValues 补充call 函数返回值
func fillReturnValues(f *ast.File) error {
	incReturns := map[*ast.ReturnStmt]*ast.FuncType{}
	ast.Walk(visitor{returns: incReturns}, f)

returnsLoop:
	for ret, ftyp := range incReturns {
		if ftyp.Results == nil {
			continue
		}
		numRVs := len(ret.Results)
		if numRVs == len(ftyp.Results.List) {
			continue
		}
		if numRVs == 0 {
			continue
		}
		if numRVs > len(ftyp.Results.List) {
			continue
		}
		if e, ok := ret.Results[0].(*ast.CallExpr); ok {
			if !hasSingleCallReturnVal(e) {
				continue
			}
		}
		zvs := make([]ast.Expr, len(ftyp.Results.List)-numRVs)
		for i, rt := range ftyp.Results.List[:len(zvs)] {
			zv := newZeroValueNode(rt.Type)
			if zv == nil {
				continue returnsLoop
			}
			zvs[i] = zv
		}
		ret.Results = append(zvs, ret.Results...)
	}
	return nil
}

// newZeroValueNode 新建零值节点
func newZeroValueNode(typ ast.Expr) ast.Expr {
	switch v := typ.(type) {
	case *ast.Ident:
		switch v.Name {
		case "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "byte", "rune", "uint", "int", "uintptr":
			return &ast.BasicLit{Kind: token.INT, Value: "0"}
		case "float32", "float64":
			return &ast.BasicLit{Kind: token.FLOAT, Value: "0"}
		case "complex64", "complex128":
			return &ast.BasicLit{Kind: token.IMAG, Value: "0"}
		case "bool":
			return &ast.Ident{Name: "false"}
		case "string":
			return &ast.BasicLit{Kind: token.STRING, Value: `""`}
		case "error":
			return &ast.Ident{Name: "nil"}
		}
	case *ast.ArrayType:
		if v.Len == nil {
			// slice
			return &ast.Ident{Name: "nil"}
		}
		return &ast.CompositeLit{Type: v}
	case *ast.StarExpr:
		return &ast.Ident{Name: "nil"}
	}
	return nil
}

func removeBareReturns(f *ast.File) error {
	returns := map[*ast.ReturnStmt]*ast.FuncType{}
	ast.Walk(visitor{returns: returns}, f)

returnsLoop:
	for ret, ftyp := range returns {
		if ftyp.Results == nil {
			continue
		}
		numRVs := len(ret.Results)
		if numRVs == len(ftyp.Results.List) {
			continue
		}

		if numRVs == 0 && len(ftyp.Results.List) > 0 {
			zvs := make([]ast.Expr, len(ftyp.Results.List))
			for i, rt := range ftyp.Results.List {
				if len(rt.Names) == 0 {
					continue returnsLoop
				}
				zv := &ast.Ident{Name: rt.Names[0].Name}
				zvs[i] = zv
			}
			ret.Results = append(zvs, ret.Results...)
		}
	}
	return nil
}

func containsMainFunc(file *ast.File) bool {
	for _, decl := range file.Decls {
		if f, ok := decl.(*ast.FuncDecl); ok {
			if f.Name.Name != "main" {
				continue
			}
			if len(f.Type.Params.List) != 0 {
				continue
			}
			if f.Type.Results != nil && len(f.Type.Results.List) != 0 {
				continue
			}
			return true
		}
	}
	return false
}
