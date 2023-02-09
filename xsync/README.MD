# xsync

同步辅助函数

- `bool`、`string`、`time.Time`、`int32`、`int64`、`uint32`、`uint64` 原子性操作
- `sync.Cond` 封装扩展
- 支持设置指定 `timeout` 的 `sync.WaitGroup`

# 例子
```go
var b AtomicBool
fmt.Println(b.Get())
b.Set(true)
fmt.Println(b.Get())
```
Output:
```text
false
true
```