package retry

import (
	"context"
	"time"
)

//go:generate optionGen  --option_return_previous=false
func OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// annotation@Limit(comment="最大尝试次数")
		"Limit": uint(10),
		// annotation@Delay(comment="固定延迟")
		"Delay": time.Duration(100 * time.Millisecond),
		// annotation@MaxJitter(comment="延迟最大抖动")
		"MaxJitter": time.Duration(100 * time.Millisecond),
		// annotation@OnRetry(comment="每次重试会先调用此方法")
		"OnRetry": func(n uint, err error) { /*do nothing now*/ },
		// annotation@RetryIf(comment="何种error进行重试")
		"RetryIf": (func(err error) bool)(IsRecoverable),
		// annotation@DelayType(comment="何种error进行重试")
		"DelayType": DelayTypeFunc(CombineDelay(BackOffDelay, RandomDelay)),
		// annotation@LastErrorOnly(comment="是否只返回最后遇到的error")
		"LastErrorOnly": false,
		// annotation@Context(comment="context，可以设定超时等")
		"Context": context.Context(context.Background()),
		// annotation@MaxDelay(comment="最大延迟时间")
		"MaxDelay":         time.Duration(0),
		"MaxBackOffNInner": uint(0),
	}
}
