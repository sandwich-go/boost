# xio

由于 `Read` 是 `block` 操作，内部为每一次 `Read` 启动了独立协程协助读取，如果超时或退出，则返回错误


# 例子
```go
ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
r = NewReader(ctx, bytes.NewReader(buf))
defer cancel()
_, err := r.Read(buf2)
xpanic.WhenError(err)
```