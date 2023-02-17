package xcompress

//go:generate optiongen --option_with_struct_name=false --new_func=NewOptions --xconf=true --empty_composite_nil=true --usage_tag_name=usage
func OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Type":  Type(GZIP),              // @MethodComment(解压缩类型)
		"Level": int(DefaultCompression), // @MethodComment(解压缩等级)
	}
}
