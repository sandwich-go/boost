# xpanic

`panic` 辅助函数

- `error` 不为 `nil` 时， `panic`
- `true` 时，`panic`
- `try cache panic`

# 例子
```go
Try(func() {
    WhenErrorAsFmtFirst(err, "%w, %d", 1)
}).Catch(func(err E) {
    fmt.Println(err)
})
```
Output:
```text
error, 1
```