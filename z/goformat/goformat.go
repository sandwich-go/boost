package goformat

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"
)

// ProcessFile
// filename 源文件路径
func ProcessFile(filename string, opts ...Option) ([]byte, error) {
	return Process(filename, nil, opts...)
}

// ProcessCode
// src 源码
func ProcessCode(src []byte, opts ...Option) ([]byte, error) {
	return Process("", src, opts...)
}

// Process 格式化文件或内容
// filename 源文件路径
// src 源码
func Process(filename string, src []byte, opts ...Option) ([]byte, error) {
	conf := NewOptions(opts...)

	fileSet := token.NewFileSet()
	file, adjust, err := parse(fileSet, filename, src, conf)
	if err != nil {
		return nil, err
	}

	if err = fillReturnValues(file); err != nil {
		return nil, err
	}

	if conf.RemoveBareReturns {
		if err = removeBareReturns(file); err != nil {
			return nil, err
		}
	}

	var buf bytes.Buffer
	err = printer.Fprint(&buf, fileSet, file)
	if err != nil {
		return nil, err
	}
	out := buf.Bytes()
	if adjust != nil {
		out = adjust(src, out)
	}
	return format.Source(out)
}

func parseMainFragment(fset *token.FileSet, filename string, src []byte, mode parser.Mode) (*ast.File, func(orig, src []byte) []byte, error) {
	psrc := append([]byte("package main;"), src...)
	file, err := parser.ParseFile(fset, filename, psrc, mode)
	if err == nil {
		if containsMainFunc(file) {
			return file, nil, nil
		}
		adjust := func(orig, src []byte) []byte {
			src = src[len("package main\n"):]
			return matchSpace(orig, src)
		}
		return file, adjust, nil
	}
	return nil, nil, err
}

func parseDeclarationFragment(fset *token.FileSet, filename string, src []byte, mode parser.Mode) (*ast.File, func(orig, src []byte) []byte, error) {
	fsrc := append(append([]byte("package p; func _() {"), src...), '}')
	file, err := parser.ParseFile(fset, filename, fsrc, mode)
	if err == nil {
		adjust := func(orig, src []byte) []byte {
			src = src[len("package p\n\nfunc _() {"):]
			src = src[:len(src)-len("}\n")]
			src = bytes.Replace(src, []byte("\n\t"), []byte("\n"), -1)
			return matchSpace(orig, src)
		}
		return file, adjust, nil
	}
	return nil, nil, err
}

func parse(fset *token.FileSet, filename string, src []byte, conf *Options) (*ast.File, func(orig, src []byte) []byte, error) {
	mode := parser.ParseComments
	if conf.AllErrors {
		mode |= parser.AllErrors
	}
	var err error
	var file *ast.File
	if src == nil {
		file, err = parser.ParseFile(fset, filename, nil, mode)
	} else {
		file, err = parser.ParseFile(fset, filename, src, mode)
	}
	if err == nil {
		return file, nil, nil
	}
	if !conf.Fragment || !strings.Contains(err.Error(), "expected 'package'") {
		return nil, nil, err
	}

	var adjust func(orig, src []byte) []byte
	file, adjust, err = parseMainFragment(fset, filename, src, mode)
	if err == nil {
		return file, adjust, nil
	}

	if !strings.Contains(err.Error(), "expected declaration") {
		return nil, nil, err
	}
	file, adjust, err = parseDeclarationFragment(fset, filename, src, mode)
	if err == nil {
		return file, adjust, nil
	}
	return nil, nil, err
}
