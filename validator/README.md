# validator

通过 `tag` 或者 `rules` 进行数据校验

# 例子

```go
type Test struct {
    Name string `val:"len=4"`
}

func main() {
    s := &Test{
        Name: "TEST",
    }
    err := Default.Struct(context.Background(), s)
    if err == nil {
        fmt.Println("validate ok")
    }   

    s.Name = ""
    err = Default.Struct(context.Background(), s)
    if err != nil {
        fmt.Println(err.Error())
    }
}

```

Output:
```text
validate ok
Key: 'Test.Name' Error:Field validation for 'Name' failed on the 'len' tag
```