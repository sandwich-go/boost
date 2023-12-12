package hide

//go:generate optionGen  --option_return_previous=false
func OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Suffix":          "hash",
		"PrefixKeep":      3,
		"SuffixKeep":      3,
		"HideLenMin":      3,
		"HideReplaceWith": rune('*'),
		// 0 等长，否则按照指定的长度替换
		"HideReplaceLen": 0,
	}
}
