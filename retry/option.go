package retry

import (
	"context"
	"time"
)

//go:generate optionGen  --option_return_previous=false
func OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Limit":            uint(10),                                       // 最大尝试次数
		"Delay":            time.Duration(100 * time.Millisecond),          // 固定延迟
		"MaxJitter":        time.Duration(100 * time.Millisecond),          // 延迟最大抖动
		"OnRetry":          func(n uint, err error) { /*do nothing now*/ }, // 每次重试会先调用此方法
		"RetryIf":          (func(err error) bool)(IsRecoverable),          // 何种error进行重试
		"DelayType":        DelayTypeFunc(CombineDelay(BackOffDelay, RandomDelay)),
		"LastErrorOnly":    false,                                 // 是否只返回最后遇到的error
		"Context":          context.Context(context.Background()), // context，可以设定超时等
		"MaxDelay":         time.Duration(0),
		"MaxBackOffNInner": uint(0),
	}
}
