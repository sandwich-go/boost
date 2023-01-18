package annotation

//go:generate optiongen --option_with_struct_name=false --new_func=NewOptions --xconf=true --empty_composite_nil=true --usage_tag_name=usage
func OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"MagicPrefix": "annotation@",  // @MethodComment(只有包含 MagicPrefix 的行，才能萃取到注释)
		"LowerKey":    true,           // @MethodComment(key是否为转化为小写)
		"Descriptors": []Descriptor{}, // @MethodComment(描述数组)
	}
}
