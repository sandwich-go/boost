package xchan

//go:generate optionGen  --option_return_previous=false
func OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"CallbackOnBufCount": int64(0),
		"Callback":           (func(bufCount int64))(nil),
	}
}
