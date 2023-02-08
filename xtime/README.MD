# xtime

`time` 辅助函数

- 时间函数
- 时间模拟
- `timewheel`

# 例子
```go
w := NewWheel(time.Second, 3)
defer w.Stop()

time.Sleep(500 * time.Millisecond)
t1 := time.Now()

go func() {
    select {
    case <-w.After(1 * time.Second):
        fmt.Printf("expected 1s, got %s\n", time.Since(t1))
    }
}()
```
Output:
```text
expected 1s, got 1.498680893s
```