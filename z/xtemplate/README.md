# xtemplate

`tempalte` 辅助函数

- 模板中支持 `github.com/Masterminds/sprig.FuncMap` 返回的函数
- 支持生成指定文件
- 支持 `go` 文件格式化

# 例子
```go
s := `a {{ .val1 }} {{ .val2 }}`
s1, err := Execute(s, map[string]interface{}{"val1": "b", "val2": 2})
xpanic.WhenError(err)
fmt.Println(string(s1))
```
Output:
```text
a b 2
```