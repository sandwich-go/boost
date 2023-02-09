package main

import (
	"github.com/sandwich-go/boost/internal/log"
	"github.com/sandwich-go/boost/xmath/internal/template2"
	"github.com/sandwich-go/boost/xtemplate"
)

func main() {
	for _, src := range []struct {
		Name     string
		FileName string
		TPL      string
		Args     interface{}
	}{
		{Name: "math", FileName: "./../math.go", TPL: template2.MathTPL, Args: template2.GetMathTPLArgs()},
		{Name: "math_test", FileName: "./../math_test.go", TPL: template2.MathTestTPL, Args: template2.GetMathTPLArgs()},
	} {
		_, err := xtemplate.Execute(src.TPL, src.Args,
			xtemplate.WithOptionName(src.Name),
			xtemplate.WithOptionFileName(src.FileName),
		)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
