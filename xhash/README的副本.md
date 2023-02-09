# xhash

`hash` 算法

- 支持对文件进行 `md5`
- 支持对数据流进行 `md5`
- 支持 `jenkins` `hash` 算法

# 例子
```go
s, err := md5.Buffer(bytes.NewReader([]byte("aaaaaaaa")))
if err != nil {
    panic(er)
}
fmt.Println(s)

hint, _ := jenkins.HashString("aaaaaaaa", 0, 0)
fmt.Println(hint)
```
Output:
```text
3dbe00a167653a1aaee01d93e77e730e
783334759
```