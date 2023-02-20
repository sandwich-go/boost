package main

import (
	"github.com/sandwich-go/boost/internal/log"
	"github.com/sandwich-go/boost/misc/xtemplate"
	"github.com/sandwich-go/boost/xmap/internal/template2"
)

func main() {
	for _, src := range []struct {
		Name     string
		FileName string
		TPL      string
		Args     interface{}
	}{
		{Name: "walk", FileName: "./../walk.go", TPL: template2.WalkTPL, Args: template2.GetWalkTPLArgs()},
		{Name: "walk_test", FileName: "./../walk_test.go", TPL: template2.WalkTestTPL, Args: template2.GetWalkTPLArgs()},
		{Name: "equal", FileName: "./../equal.go", TPL: template2.EqualTPL, Args: template2.GetEqualTPLArgs()},
		{Name: "equal_test", FileName: "./../equal_test.go", TPL: template2.EqualTestTPL, Args: template2.GetEqualTPLArgs()},
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
