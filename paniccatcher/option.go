package paniccatcher

import (
	"time"
)

type OnRecover = func(tag string, reason interface{})

//go:generate optiongen  --option_with_struct_name=true --option_return_previous=false
func AutoRecoverOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// annotation@DelayTime(comment="每次panic后重启delay的时间 Note: 这里应该可以直接对接到retry package，复用重试逻辑")
		"DelayTime": time.Duration(0),
		"OnRecover": OnRecover(nil),
	}
}
