# xhash

`hash` 算法

- 支持对文件进行 `md5`
- 支持对数据流进行 `md5`
- 支持 `jenkins` `hash` 算法
- 支持 `hash14v`，对 `uint64` 与 `string` 相互转换

# 例子
```go
s, err := md5.Buffer(bytes.NewReader([]byte("aaaaaaaa")))
xpanic.WhenError(err)
fmt.Println(s)

hint, _ := jenkins.HashString("aaaaaaaa", 0, 0)
fmt.Println(hint)

v := hash14v.ToV(123456789)
fmt.Println(v)

id := hash14v.ToId(v)
fmt.Println(id)
```
Output:
```text
3dbe00a167653a1aaee01d93e77e730e
783334759
SVMAFEV
123456789
```