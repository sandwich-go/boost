package main

import (
	"github.com/sandwich-go/boost/internal/log"
	"github.com/sandwich-go/boost/misc/xtemplate"
	"github.com/sandwich-go/boost/xdebug/internal/template2"
)

func main() {
	for _, src := range []struct {
		Name     string
		FileName string
		TPL      string
		Args     interface{}
	}{
		{Name: "dependency", FileName: "./../gen_dependency.go", TPL: template2.DependencyTPL, Args: template2.GetDependencyTPLArgs()},
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
