package main

import (
	"github.com/sandwich-go/boost/internal/log"
	"github.com/sandwich-go/boost/xmath/internal/template2"
	"github.com/sandwich-go/boost/xtemplate"
)

func main() {
	// 生成 walk.go
	_, err := xtemplate.Execute(template2.MathTPL, template2.GetMathTPLArgs(),
		xtemplate.WithOptionName("math"),
		xtemplate.WithOptionFileName("./../math.go"),
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	// 生成 walk_test.go
	_, err = xtemplate.Execute(template2.MathTestTPL, template2.GetMathTPLArgs(),
		xtemplate.WithOptionName("math_test"),
		xtemplate.WithOptionFileName("./../math_test.go"),
	)
	if err != nil {
		log.Fatal(err.Error())
	}
}
