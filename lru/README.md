# LRU

`Engine[V]` 提供按最后访问时间续期的滑动 TTL。默认使用双向环形链表，构建时指定
`-tags boost_lru_heap` 使用惰性最小堆。两套实现共用 `NewEngine`、`NewLazyEngine`、
`Add`、`Remove`、`Promote`、`Node.Value` 和 `Node.SetValue` 接口。

默认链表的 `Add`、`Remove`、`Promote` 都修改节点链接，调用方必须持构造时传入的
locker 对应的写锁。过期扫描从最久未访问的节点开始，遇到存活节点即停止。

堆模式的 `Add`、`Remove` 仍需写锁；`Promote` 只原子更新访问时间，可在读锁下调用。
cleaner 在旧排期到达时校正节点位置，大量到期节点采用批量整理。
`Node.Value` 的读取和 `Node.SetValue` 的修改也须由调用方同步。
locker 为 nil 时，调用方须自行保证清理与节点操作互斥。

`PromoteNeedsWriteLock` 是当前实现的锁要求，调用方可用这个编译期常量选择锁策略。
`ImplementationName()` 返回 `linkedlist` 或 `heap`，可用于启动日志确认构建模式。
服务、测试与 Go plugin 等构建入口应使用一致的 tag。

`NewEngine(..., WithAutoClean(false))` 关闭每个 Engine 自带的后台清理，
由外层共用的调度任务调用 `Expire()`。`Add` 不会启动清理任务；不传选项时仍自动清理。
`Expire()` 使用构造时传入的 locker；外层已经持锁时传 nil locker，避免重复加锁。

`ExpireBatch(limit)` 在单次调用内最多处理 `limit` 个到期排期，返回是否仍有排期需要检查。
`limit=0` 沿用不限量的 `Expire()`，负数会 panic。正数预算同时限制堆模式的续期重排，
不会触发全堆线性整理；返回 true 不一定意味着有同样数量的实际过期对象。
外层收到 true 后，应释放本批持有的锁，再继续下一批；单次返回时节点结构和索引回调已更新。


```sh
go test ./lru ./xcontainer/redblacktree
go test -tags boost_lru_heap ./lru ./xcontainer/redblacktree
go test -race ./lru ./xcontainer/redblacktree
go test -race -tags boost_lru_heap ./lru ./xcontainer/redblacktree
```
