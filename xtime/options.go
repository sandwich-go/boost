package xtime

import "time"

//go:generate optionGen  --option_return_previous=false
func OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		// annotation@TickDuration(comment="tick的duration，大于0是自动开启ticker")
		"TickDuration": time.Duration(0),
		// annotation@TickHostingMode(comment="全托管模式，内部起一个协程执行tick的func")
		"TickHostingMode": true,
		// annotation@CountGauge(comment="统计tick数量监控")
		"TickCount": CountGauge(&noopGauge{}),
	}
}

type CountGauge interface {
	Dec()
	Inc()
}

var _ CountGauge = (*noopGauge)(nil)

type noopGauge struct{}

func (n noopGauge) Dec() {}

func (n noopGauge) Inc() {}
