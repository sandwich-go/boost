package main

import (
	"github.com/sandwich-go/boost/geom/internal/template2"
	"github.com/sandwich-go/boost/internal/log"
	"github.com/sandwich-go/boost/xtemplate"
)

func main() {
	for _, src := range []struct {
		Name     string
		FileName string
		TPL      string
		Args     interface{}
	}{
		{Name: "int", FileName: "./../int.go", TPL: template2.GeomTPL, Args: map[string]interface{}{"Key": "int"}},
		{Name: "int_test", FileName: "./../int_test.go", TPL: template2.GeomTestTPL, Args: map[string]interface{}{"Key": "int"}},

		{Name: "int8", FileName: "./../int8.go", TPL: template2.GeomTPL, Args: map[string]interface{}{"Key": "int8"}},
		{Name: "int8_test", FileName: "./../int8_test.go", TPL: template2.GeomTestTPL, Args: map[string]interface{}{"Key": "int8"}},

		{Name: "int16", FileName: "./../int16.go", TPL: template2.GeomTPL, Args: map[string]interface{}{"Key": "int16"}},
		{Name: "int16_test", FileName: "./../int16_test.go", TPL: template2.GeomTestTPL, Args: map[string]interface{}{"Key": "int16"}},

		{Name: "int32", FileName: "./../int32.go", TPL: template2.GeomTPL, Args: map[string]interface{}{"Key": "int32"}},
		{Name: "int32_test", FileName: "./../int32_test.go", TPL: template2.GeomTestTPL, Args: map[string]interface{}{"Key": "int32"}},

		{Name: "int64", FileName: "./../int64.go", TPL: template2.GeomTPL, Args: map[string]interface{}{"Key": "int64"}},
		{Name: "int64_test", FileName: "./../int64_test.go", TPL: template2.GeomTestTPL, Args: map[string]interface{}{"Key": "int64"}},
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
