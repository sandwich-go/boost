# xmath

`math` 辅助函数

- 返回两值最大值
- 返回两值最小值
- 获取绝对值
- 增加值，返回不小于0值的数值
- 判断 `float32`、`float64` 值是否相等，是否为0值
- 将十进制、十六进制字符串转化为 `uint64`

# 例子
```go
fmt.Println("max:", MaxInt32(3, 2))
fmt.Println("min:", MinInt32(3, 2))
fmt.Println("val:", MustParseUint64("0x16"))
```
Output:
```text
max: 3
min: 2
val: 22
```