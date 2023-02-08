package main

import (
	"github.com/sandwich-go/boost/internal/log"
	"github.com/sandwich-go/boost/xslice/internal/template2"
	"github.com/sandwich-go/boost/xtemplate"
)

func main() {
	for _, src := range []struct {
		Name     string
		FileName string
		TPL      string
		Args     interface{}
	}{
		{Name: "slice", FileName: "./../slice.go", TPL: template2.SliceTPL, Args: template2.GetSliceTPLArgs()},
		{Name: "slice_test", FileName: "./../slice_test.go", TPL: template2.SliceTestTPL, Args: template2.GetSliceTPLArgs()},
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
