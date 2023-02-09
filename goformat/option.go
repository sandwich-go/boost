package goformat

//go:generate optiongen --option_with_struct_name=false --new_func=NewOptions --xconf=true --empty_composite_nil=true --usage_tag_name=usage
func OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Fragment":          false, // @MethodComment(允许解析源文件片段代码)
		"AllErrors":         false, // @MethodComment(打印所有的语法错误到标准输出。如果不使用此标记，则只会打印不同行的前10个错误)
		"RemoveBareReturns": false, // @MethodComment(移除无效的return)
	}
}
