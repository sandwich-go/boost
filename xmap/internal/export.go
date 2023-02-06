package main

import (
	"github.com/sandwich-go/boost/internal/log"
	"github.com/sandwich-go/boost/xmap/internal/template2"
	"github.com/sandwich-go/boost/xtemplate"
)

func main() {
	// 生成 walk.go
	_, err := xtemplate.Execute(template2.WalkTPL, template2.GetWalkTPLArgs(),
		xtemplate.WithOptionName("walk"),
		xtemplate.WithOptionFileName("./../walk.go"),
	)
	if err != nil {
		log.Fatal(err.Error())
	}
	// 生成 walk_test.go
	_, err = xtemplate.Execute(template2.WalkTestTPL, template2.GetWalkTPLArgs(),
		xtemplate.WithOptionName("walk_test"),
		xtemplate.WithOptionFileName("./../walk_test.go"),
	)
	if err != nil {
		log.Fatal(err.Error())
	}

}
