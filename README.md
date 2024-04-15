# boost

`boost` 工具箱

- `geom` [地理位置](https://github.com/sandwich-go/boost/tree/main/geom/README.md)
- `httputil` [HTTP工具](https://github.com/sandwich-go/boost/tree/main/httputil/README.md)
- `middleware` 中间件
- `misc` [杂项工具](https://github.com/sandwich-go/boost/tree/main/misc/README.md)
- `module` [Module工具](https://github.com/sandwich-go/boost/tree/main/module/README.md)
- `plugin` 插件
- `retry` [重试](https://github.com/sandwich-go/boost/tree/main/retry/README.md)
- `singleflight` [SingleFlight](https://github.com/sandwich-go/boost/tree/main/singleflight/README.md)
- `types` [扩展的数据类型](https://github.com/sandwich-go/boost/tree/main/types/README.md)
- `validator` [校验器](https://github.com/sandwich-go/boost/tree/main/validator/README.md)
- `version` [程序版本](https://github.com/sandwich-go/boost/tree/main/version/README.md)
- `xcmd` [命令行/ENV参数](https://github.com/sandwich-go/boost/tree/main/xcmd/README.md)
- `xcompress` [解压缩器](https://github.com/sandwich-go/boost/tree/main/xcompress/README.md)
- `xcontainer` [容器](https://github.com/sandwich-go/boost/tree/main/xcontainer/README.md)
- `xconv` [内置类型转换辅助工具](https://github.com/sandwich-go/boost/tree/main/xconv/README.md)
- `xcopy` [深拷贝](https://github.com/sandwich-go/boost/tree/main/xcopy/README.md)
- `xcrypto` [加密](https://github.com/sandwich-go/boost/tree/main/xcrypto/README.md)
- `xchan` [Unbounded chan with ring buffer](https://github.com/sandwich-go/boost/tree/main/xchan/README.md)
- `xdebug` [Debug](https://github.com/sandwich-go/boost/tree/main/xdebug/README.md)
- `xencoding` [编码解码器](https://github.com/sandwich-go/boost/tree/main/xencoding/README.md)
- `xerror` [error wrapper](https://github.com/sandwich-go/boost/tree/main/xerror/README.md)
- `xhash` [hash算法](https://github.com/sandwich-go/boost/tree/main/xhash/README.md)
- `xio` [异步io](https://github.com/sandwich-go/boost/tree/main/xio/README.md)
- `xip` [ip/port辅助函数](https://github.com/sandwich-go/boost/tree/main/xip/README.md)
- `xmap` [map辅助函数](https://github.com/sandwich-go/boost/tree/main/xmap/README.md)
- `xmath` [math辅助函数](https://github.com/sandwich-go/boost/tree/main/xmath/README.md)
- `xos` [系统辅助函数](https://github.com/sandwich-go/boost/tree/main/xos/README.md)
- `xpanic` [panic辅助函数](https://github.com/sandwich-go/boost/tree/main/xpanic/README.md)
- `xpool` [pool辅助函数](https://github.com/sandwich-go/boost/tree/main/xpool/README.md)
- `xproc` [command辅助函数](https://github.com/sandwich-go/boost/tree/main/xproc/README.md)
- `xrand` [随机辅助函数](https://github.com/sandwich-go/boost/tree/main/xrand/README.md)
- `xslice` [切片辅助函数](https://github.com/sandwich-go/boost/tree/main/xslice/README.md)
- `xstrings` [字符串辅助函数](https://github.com/sandwich-go/boost/tree/main/xstrings/README.md)
- `xsync` [同步辅助函数](https://github.com/sandwich-go/boost/tree/main/xsync/README.md)
- `xtest`
- `xtime` [time辅助函数](https://github.com/sandwich-go/boost/tree/main/xtime/README.md)
- `z` [编译期辅助函数](https://github.com/sandwich-go/boost/tree/main/z/README.md)

# 设置 `logger`
通过 `boost.InstallLogger` 来设置自定义的 `logger`

自定义的 `logger` 需要实现以下接口

```go
type Logger interface {
    Debug(string)
    Info(string)
    Warn(string)
    Error(string)
    Fatal(string)
}
```