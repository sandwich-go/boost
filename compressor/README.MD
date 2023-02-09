# compressor

解压缩器

- 支持 `gzip` 方式，以不同的压缩等级进行解压缩
- 支持 `snappy` 方式解压缩

**注意** :
- 默认使用 `gzip` 方式解压缩

# 例子

```go
frame := []byte("time.Duration,[]time.Duration,map[string]*Redis此类的非基础类型的slice或者map都需要辅助指明类型")
c, err0 := New(WithType(GZIP), WithLevel(BestSpeed))
if err0 != nil {
    panic(err0)
}
flatFrame, err1 := c.Flat(frame)
if err1 != nil {
    panic(err1)
}
inflateFrame, err2 := c.Inflate(flatFrame)
if err2 != nil {
    panic(err2)
}
if string(frame) != string(inflateFrame) {
    panic("flat/inflate fail")
}
fmt.Println("OK")
```

Output:
```text
OK
```