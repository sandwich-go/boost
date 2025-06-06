package xtemplate

import "html/template"

// Filter 过滤器
// 在生成指定文件或返回模板填充后数据前，可以指定过滤器对数据进行二次加工
type Filter func([]byte) []byte

//go:generate optiongen  --option_with_struct_name=true --option_return_previous=false
func OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// annotation@Name(comment="指定模板的名称")
		"Name": "xtemplate",
		// annotation@FileName(comment="文件名称，若不为空，则会生成对应的文件，若文件为 .go 文件，则会格式化")
		"FileName": "",
		// annotation@Filers(comment="在生成指定文件或返回模板填充后数据前，指定过滤器，可以对数据进行二次加工")
		"Filers": []Filter(nil),
		// annotation@FuncMap(comment="自定义的函数集")
		"FuncMap": template.FuncMap{},
	}
}
