# encoding2

编码解码器

- 提供压缩/解压缩编码解码方式
- 提供加密/解密编码解码方式
- 提供 `json` 编码解码方式
- 提供 `pbjson` 编码解码方式
- 提供 `protobuf` 编码解码方式
- 提供 `msgpack` 编码解码方式


# 例子

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

frame := []byte("time.Duration,[]time.Duration,map[string]*Redis此类的非基础类型的slice或者map都需要辅助指明类型")

for _, n := Codecs() {
    codec := GetCodec(nn)
    bs, err := codec.Marshal(frame)
    xpanic.WhenError(err)
    var raw []byte
    err = codec.Unmarshal(bs, &raw)
    xpanic.WhenError(err)
    fmt.Println(string(raw))
}

ctx = WithContext(ctx, GetCodec(encrypt.AESCodecName))
_, err0 := FromContext(ctx).Marshal(frame)
xpanic.WhenError(err0)

var raw1 []byte
err0 = FromContext(ctx).Unmarshal(bs, &raw1)
xpanic.WhenError(err0)
fmt.Println(string(raw1))
```