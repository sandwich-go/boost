# xcopy

深拷贝

- 若实现 Clone() [proto1](https://github.com/golang/protobuf/tree/master/proto).Message 接口，可快速进行深拷贝
- 若实现 Clone() [proto2](https://github.com/protocolbuffers/protobuf-go/tree/master/proto).Message 接口，可快速进行深拷贝
- 若实现 `DeepCopy() interface{}` 接口，可快速进行深拷贝
- 其他形式则通过反射进行深拷贝


**注意** :
- 通过反射进行深拷贝，注意未导出的属性，未导出的属性无法深拷贝

# 例子
```go
type sub0 struct {
    i uint32    // 含有未导出的属性，无法深拷贝
}

type sub1 struct {
    i uint32    // 含有未导出的属性，但实现了 Clone() proto1.Message、Clone() proto2.Message、DeepCopyInterface，可以深拷贝
}

func (s *sub1) DeepCopy() interface{} {
    if s == nil {
        return s
    }
    return &sub1{i: s.i}
}

type sub2 struct{
    I uint32    // 没有未导出的属性，可以深拷贝
}

func main() {
    var s0 = &sub0{i: 1}            // 不会进行拷贝
    var s1 = &sub1{i: 2}
    var s2 = &sub2{I: 3}
	
    fmt.Println(DeepCopy(s0).i)
    fmt.Println(DeepCopy(s1).i)
    fmt.Println(DeepCopy(s2).i)
}
```

Output:
```text
0
2
3
```