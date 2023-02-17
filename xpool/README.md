# xpool


- 支持 `bytes` 池
- 支持工作协程池

# 例子
```go
p := NewSyncBytesPool(1, 100, 2)
buff := p.Alloc(6)
defer func() {
	p.Free(buff)
}
```