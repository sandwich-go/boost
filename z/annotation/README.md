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
xpanic.WhenError(err0)
fmt.Println("Name:", ann.Name())
fmt.Println("Line:", ann.Line())
akVal, err1 := ann.Int("AK")
xpanic.WhenError(err1)
fmt.Println("AK:", akVal)
fmt.Println("AV:", ann.String("AV"))
```

Output:
```text
Name: a
Line: // annotation@A( AK=127, AV="AAAAA" )
AK: 127
AV: AAAAA
```