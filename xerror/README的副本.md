# xerror

封装`error`

- 支持多个 `error` 封装成单一 `error` 进行函数参数传递
- 支持携带错误码信息 `error` 对象
- 支持包含调用信息 `error` 对象
- 支持 `Logic` 层异常 `error` 对象
- 支持返回最底层的错误信息
- 支持返回堆栈信息

# 例子
```go
var arr Array
arr.Push(errors.New("error 1"))
arr.Push(errors.New("error 2"))

if arr.Err() != nil {
    fmt.Println(arr.Error())
}

err := New(WithText("io error"), WithCode(500), WithStack())
errW := Wrap(err, "link error")
errW = Wrap(errW, "session error")
fmt.Println(errW.Error())
fmt.Println(Caller(err.Cause(), 0))
```
Output:
```text
2 errors occurred:
    #1: error 1
    #2: error 2
    
session error: link error: io error
array_test.go github.com/sandwich-go/boost/xerror.TestArray.func1 35    
```