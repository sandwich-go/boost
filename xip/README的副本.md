# xip

`ip`、`port` 辅助函数

- 获取除回环地址外的本机所有地址
- 判断某个地址是否是内网地址
- 获取空闲端口

# 例子
```go
port, err := GetFreePort()
if err != nil {
    panic(er)
}
fmt.Println("free port:", port)

localIP := GetLocalIP()
fmt.Printf("%s is intranet? %v\n", localIP, IsIntranet(localIP))
```
Output:
```text
free port: 55304
10.0.48.24 is intranet? true
```