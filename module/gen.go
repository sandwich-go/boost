package module

//go:generate optiongen --option_with_struct_name=true
func ModuleOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Name":    string(""),
		"OnInit":  (func())(nil),
		"OnClose": (func())(nil),
	}
}
