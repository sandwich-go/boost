package lru

// EngineOption 配置 NewEngine 的清理行为。
type EngineOption func(*engineOptions)

type engineOptions struct {
	autoClean bool
}

// WithAutoClean 控制是否启动后台清理，默认开启。
// 关闭后由调用方定期调用 Expire，Add 不会再启动后台任务。
func WithAutoClean(enabled bool) EngineOption {
	return func(options *engineOptions) { options.autoClean = enabled }
}

func autoCleanEnabled(opts []EngineOption) bool {
	options := engineOptions{autoClean: true}
	for _, option := range opts {
		option(&options)
	}
	return options.autoClean
}

// Expire 立即清理过期项，沿用构造时的 locker 和过期回调。
// 外部已持锁时应在构造时传 nil locker，避免重复加锁。
func (e *Engine[V]) Expire() { e.expire() }
