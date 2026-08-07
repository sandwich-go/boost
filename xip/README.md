# xip

`ip`、`port` 辅助函数

- 获取除回环地址外的本机所有地址
- 判断某个地址是否是内网地址
- 获取空闲端口

# 例子
```go
port, err := GetFreePort()
xpanic.WhenError(err)
fmt.Println("free port:", port)

localIP := GetLocalIP()
fmt.Printf("%s is intranet? %v\n", localIP, IsIntranet(localIP))
```
Output:
```text
free port: 55304
10.0.48.24 is intranet? true
```

# 热路径取本机 IP

`GetLocalIP` 每次调用都要枚举网卡并两次全量扫描环境变量，darwin/arm64 实测单次
79µs、42216 B、581 allocs，开销随环境变量条数增长。本机 IP 在进程生命周期内不变，
热路径（如按请求上报埋点）请用只解析一次的 `GetLocalIPCached`：

```go
// 每条埋点都取一次本机 IP
fields = append(fields, logbus.String("#ip", xip.GetLocalIPCached()))
```

| | ns/op | B/op | allocs/op |
|---|---|---|---|
| `GetLocalIP` | 80514 | 42216 | 581 |
| `GetLocalIPCached`（缓存命中） | 2.4 | 0 | 0 |

```text
go test -run '^$' -bench 'BenchmarkGetLocalIP' -benchmem -benchtime=2s -count=3 ./xip/
goos: darwin  goarch: arm64  cpu: Apple M2 Pro
```

`GetLocalIPCached` 的缓存在**首次调用**时填充，且解析失败（空串）不写缓存、下次重新解析。
依赖 `boost_ip_prefer_vpn` / `x_sandwich_service_host` 的调用方需在首次调用前设好环境变量。