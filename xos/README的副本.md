# xos

系统跨平台辅助函数

- 支持文件、目录拷贝
- 支持目录创建删除遍历等操作
- 支持判断目录、文件是否存在
- 支持文件内容获取、写入
- 支持文件隐藏、取消隐藏

# 例子
```go
var file, dir = "a.go", "b"
fmt.Printf("file: %s, exists? %v\n", file, ExistsFile(file))
fmt.Printf("dir: %s, exists? %v\n", dir, ExistsDir(dir))
```
Output:
```text
file: a.go, exists? false
dir: b, exists? false
```