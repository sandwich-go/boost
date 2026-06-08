# boost

`boost` 工具箱 —— 通用 Go 工具库，sandwich-go 系列下游业务（myproxy /
dataserver / cmd/* 等）共享依赖。

**Go 版本要求**：`go 1.25.0` + `toolchain go1.25.11`（详见 `go.mod`；
下游 myproxy / dataserver 都已 1.25.x，本仓与之同步）。

**1.4 主要变化**（详见 [AGENTS.md §6.1](./AGENTS.md) 迁移表）：

- `xmath` / `xmap` / `xslice` 删除 502 个具名函数（`MaxInt32` / `EqualStringStringMap`
  / `IntsContain` 等），全部走泛型版（`Max` / `Equal` / `Contain` 等）。
- `internal/template2/` 大量退役（仅 `xdebug/internal/template2/` 保留），
  `xrand.StringWithTimestamp` 等 hot path 优化。
- CI 全门槛硬钉：vet / build / test / race / lint / vuln 全部硬约束，
  `benchdata/main.txt` 入仓作 perf 回归基线。

## 子包列表

### Leaf 工具（被大量 import）

- `xmath` [math 辅助函数（泛型）](https://github.com/sandwich-go/boost/tree/main/xmath/README.md)
- `xmap` [map 辅助函数（泛型）](https://github.com/sandwich-go/boost/tree/main/xmap/README.md)
- `xslice` [切片辅助函数（泛型）](https://github.com/sandwich-go/boost/tree/main/xslice/README.md)
- `xstrings` [字符串辅助函数](https://github.com/sandwich-go/boost/tree/main/xstrings/README.md)
- `xconv` [内置类型转换辅助工具](https://github.com/sandwich-go/boost/tree/main/xconv/README.md)
- `xrand` [随机辅助函数](https://github.com/sandwich-go/boost/tree/main/xrand/README.md)
- `xio` [异步 io](https://github.com/sandwich-go/boost/tree/main/xio/README.md)
- `xhash` [hash 算法](https://github.com/sandwich-go/boost/tree/main/xhash/README.md)
- `isnil` 安全 nil 检查（区分 typed-nil / untyped-nil）

### 错误 / 时间 / 调度

- `xerror` [error wrapper（带栈、错误码、Array 聚合）](https://github.com/sandwich-go/boost/tree/main/xerror/README.md)
- `xtime` [time 辅助函数（dispatcher / cron / Periodic / NowFunc）](https://github.com/sandwich-go/boost/tree/main/xtime/README.md)
- `retry` [重试](https://github.com/sandwich-go/boost/tree/main/retry/README.md)

### 并发 / 容器 / 池

- `xsync` [同步辅助函数（atomic / 超时锁 / wg）](https://github.com/sandwich-go/boost/tree/main/xsync/README.md)
- `xchan` [Unbounded chan with ring buffer](https://github.com/sandwich-go/boost/tree/main/xchan/README.md)
- `xcontainer` [容器（sortedmap / syncmap / smap / sset / redblacktree / ringbuf）](https://github.com/sandwich-go/boost/tree/main/xcontainer/README.md)
- `xparallel` 并发执行原语（错误聚合 / SliceV）
- `xpool` [pool 辅助函数（bytes pool / sync.Pool / goroutine pool）](https://github.com/sandwich-go/boost/tree/main/xpool/README.md)
- `ratelimiter` 令牌桶 / 漏桶限流
- `singleflight` [SingleFlight](https://github.com/sandwich-go/boost/tree/main/singleflight/README.md)

### 缓存

- `lru` 分片 LRU（workerPerEngine / workerHashPool 双实现）

### 网络 / IO / 编解码

- `httputil` [HTTP 工具](https://github.com/sandwich-go/boost/tree/main/httputil/README.md)
- `xencoding` [编码解码器（msgpack / json / pbjson / protobuf / encrypt / compressor）](https://github.com/sandwich-go/boost/tree/main/xencoding/README.md)
- `xcompress` [解压缩器（gzip / snappy）](https://github.com/sandwich-go/boost/tree/main/xcompress/README.md)
- `xcrypto` [加密（aes / mask / vigenere / curve25519）](https://github.com/sandwich-go/boost/tree/main/xcrypto/README.md)

### 系统级 / 入口

- `xos` [系统辅助函数](https://github.com/sandwich-go/boost/tree/main/xos/README.md)
- `xproc` [command 辅助函数](https://github.com/sandwich-go/boost/tree/main/xproc/README.md)
- `xcmd` [命令行 / ENV 参数](https://github.com/sandwich-go/boost/tree/main/xcmd/README.md)
- `xip` [ip / port 辅助函数](https://github.com/sandwich-go/boost/tree/main/xip/README.md)
- `module` [Module 工具（master / agent 生命周期）](https://github.com/sandwich-go/boost/tree/main/module/README.md)

### 调试 / 错误恢复

- `xpanic` [panic 辅助函数（recover + Try/Catch + WhenXxx）](https://github.com/sandwich-go/boost/tree/main/xpanic/README.md)
- `xdebug` [Debug（依赖检查 / build info）](https://github.com/sandwich-go/boost/tree/main/xdebug/README.md)
- `xcopy` [深拷贝](https://github.com/sandwich-go/boost/tree/main/xcopy/README.md)

### Misc / 业务子域

- `misc` [杂项工具（cloud / hrff / xtemplate / annotation / xgen）](https://github.com/sandwich-go/boost/tree/main/misc/README.md)
- `geo` 地理 / 网格
- `graph` [图形 / 几何](https://github.com/sandwich-go/boost/tree/main/graph/README.md)
- `plugin` 插件加载
- `validator` [校验器（go-playground/validator）](https://github.com/sandwich-go/boost/tree/main/validator/README.md)
- `xvalidator` protovalidate-go 封装
- `humanize` human-readable 格式化
- `types` [扩展的数据类型（bignum / mydecimal）](https://github.com/sandwich-go/boost/tree/main/types/README.md)
- `xetl` ETL / scan helpers
- `xexp` 实验性子包
- `xtest` quantile (P50/P99 等)
- `version` [程序版本](https://github.com/sandwich-go/boost/tree/main/version/README.md)
- `middleware` 中间件
- `z` [编译期辅助函数（hash / mono / wall / rand / conv）](https://github.com/sandwich-go/boost/tree/main/z/README.md)

## 设置 `logger`

通过 `boost.InstallLogger` 设置自定义的 `logger`。

自定义的 `logger` 需要实现以下接口：

```go
type Logger interface {
    Debug(string)
    Info(string)
    Warn(string)
    Error(string)
    Fatal(string)
}
```

注意：`boost.LogXxxf` 是 eager `fmt.Sprintf`（即使 logger 决定丢弃也会
先 alloc + 格式化）。hot path 慎用，业务侧建议直接调注入的 logger（zerolog
/ zap 有 lazy 评估）。仓内只在 `module/master.go` 启动/关闭时使用。

## 开发

```bash
make tools     # 一次性安装代码生成 / lint / vuln / benchstat 工具
make ci        # 提交前 / PR 准入 = vet + build + test_race + vuln
make lint      # golangci-lint v2（0 告警基线）
make bench_diff   # 改 hot path 时跑：当前 bench 对比 benchdata/main.txt
```

详见 [AGENTS.md](./AGENTS.md)（开工前必读）。

## 文档

- [AGENTS.md](./AGENTS.md) - 上岗手册（约束 / 陷阱 / 工作流）
- [CHANGELOG-1.3.md](./CHANGELOG-1.3.md) - 1.3 版本历史
- [CHANGELOG-0.1.md](./CHANGELOG-0.1.md) - 0.1 版本历史
- 1.4 CHANGELOG 由 `protokitgo sem changelog` 自动生成（每次发版时）
