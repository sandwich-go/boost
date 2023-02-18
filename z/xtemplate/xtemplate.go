package xtemplate

import (
	"bytes"
	"github.com/Masterminds/sprig"
	"github.com/sandwich-go/boost/xos"
	"github.com/sandwich-go/boost/xstrings"
	"github.com/sandwich-go/boost/z/goformat"
	"html/template"
	"path/filepath"
)

var funcMap = template.FuncMap{
	"Unescaped":  func(str string) template.HTML { return template.HTML(str) },
	"CamelCase":  xstrings.CamelCase,
	"SnakeCase":  xstrings.SnakeCase,
	"FirstLower": xstrings.FirstLower,
	"FirstUpper": xstrings.FirstUpper,
}

func init() {
	for k, v := range funcMap {
		funcMap[xstrings.FirstLower(k)] = v
	}
}

// Execute 根据指定的 template 以及 args 生成对应文本内容
func Execute(templateStr string, args interface{}, opts ...Option) ([]byte, error) {
	cfg := NewOptions(opts...)
	t, err := template.New(cfg.GetName()).
		Funcs(funcMap).
		Funcs(sprig.FuncMap()).
		Parse(templateStr)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(nil)
	err = t.Execute(buf, args)
	if err != nil {
		return nil, err
	}
	bytesUsing := buf.Bytes()
	for _, filter := range cfg.GetFilers() {
		if filter == nil {
			continue
		}
		bytesUsing = filter(bytesUsing)
	}
	if len(cfg.GetFileName()) > 0 {
		if filepath.Ext(cfg.GetFileName()) == ".go" {
			bytesUsing, err = goformat.ProcessCode(bytesUsing)
			if err != nil {
				return nil, err
			}
		}
		err = xos.FilePutContents(cfg.GetFileName(), bytesUsing)
		if err != nil {
			return nil, err
		}
	}
	return bytesUsing, nil
}
