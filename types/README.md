# types

扩展的数据类型

- 支持大数 `BigNum`

# 例子

```go
b, err := bignum.Add("100zz", "2zz")
if err != nil {
    panic(err)
}
fmt.Println(b)
```

Output:
```text
102000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
```