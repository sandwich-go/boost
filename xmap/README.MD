# xmap

`map`、辅助函数

- 有序遍历`map`（根据 `map` 中 `KEY` 进行排序，依次遍历）

# 例子
```go
var tm = make(map[string]string)
tm["b"] = "c"
tm["a"] = "b"
tm["1"] = "2"
tm["c"] = "d"

WalkStringStringMapDeterministic(tm, func(k string, v string) bool {
    fmt.Println("key:", k, "value:", v)
    return true
})
```
Output:
```text
key: 1 value: 2
key: a value: b
key: b value: c
key: c value: d
```