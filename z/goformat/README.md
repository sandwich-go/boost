# goformat

`Golang` 代码/文件格式化

# 例子

```go
var code = `func a(     ) {return}`
out, err := ProcessCode([]byte(code), WithFragment(true))
xpanic.WhenError(err)
fmt.Println(string(out))
```

Output:
```text
func a() { return }
```