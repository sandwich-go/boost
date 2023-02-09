# retry

执行函数，若失败，则按配置参数重新执行，直至成功或达到某种指定条件

# 例子

```go
var retrySum uint
err := Do(
    func(uint) error { return fmt.Errorf("%d", retrySum) },
    WithOnRetry(func(n uint, err error) { retrySum += 1 }),
    WithDelay(time.Nanosecond),
    WithLastErrorOnly(true),
)

fmt.Println(err)
```

Output:
```text
9
```