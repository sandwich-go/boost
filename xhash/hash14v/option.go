package hash14v

import "github.com/sandwich-go/boost/z"

//go:generate optiongen --option_with_struct_name=false --new_func=NewOptions --xconf=true --empty_composite_nil=true --usage_tag_name=usage
func OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"HashKey":           []byte(z.StringToBytes("nlCwbUUd")), // @MethodComment(hash使用的key)
		"HashOffset":        []byte(z.StringToBytes("FAAAAAA")),  // @MethodComment(hash的偏移值)
		"UsingReservedBuff": false,                               // @MethodComment(解压缩等级)
	}
}
