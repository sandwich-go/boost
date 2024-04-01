package rateshaping

import "time"

//go:generate optionGen  --option_return_previous=false
func OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"per":   time.Duration(time.Second),
		"Slack": int(10),
	}
}
