# annotation

注释解析器

# 格式
```text
// {$MagicPrefix}{$Name}({$Key0}={$Value0}, {$Key1}={$Value1}, {$Key2}="{$Value2}" ...)
```

**注意** :
- 字符串类型值需要使用 `""`
- `$MagicPrefix` 默认为 `annotation@`

# 例子
```text
// annotation@A( AK=127, AV="AAAAA" )
```

解析 :
```go
line := `// annotation@A( AK=127, AV="AAAAA" )`
ann, err0 := Default.Resolve(line)
if err0 != nil {
    panic(err0)
}
fmt.Println("Name:", ann.Name())
fmt.Println("Line:", ann.Line())
if akVal, err1 := ann.Int("AK"); err1 != nil {
    fmt.Println("AK:", akVal)
} else {
    panic(err1)
}
fmt.Println("AV:", ann.String("AV"))
```

Output:
```text
Name: a
Line: // annotation@A( AK=127, AV="AAAAA" )
AK: 127
AV: AAAAA
```