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

// ExpireBatch 最多处理 limit 个到期排期，返回是否仍有到期排期需要检查。
// limit 为 0 时不限量，沿用 Expire；负数会 panic。堆的续期重排也占用预算，
// 避免大量存活节点的旧排期校正长时间占锁。每次调用自行获取 locker，调用方
// 应在返回后释放外层锁，再安排下一批。
func (e *Engine[V]) ExpireBatch(limit int) bool {
	if limit < 0 {
		panic("lru: expiry limit must not be negative")
	}
	if limit == 0 {
		e.expire()
		return false
	}
	return e.expireLimited(limit)
}
