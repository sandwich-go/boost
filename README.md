# boost

`boost` 工具箱

- `annotation` [注释解析器](https://github.com/sandwich-go/boost/tree/main/annotation/README.md)
- `cloud` [对象存储](https://github.com/sandwich-go/boost/tree/main/cloud/README.md)
- `compressor` [解压缩器](https://github.com/sandwich-go/boost/tree/main/compressor/README.md)
- `encoding2` [编码解码器](https://github.com/sandwich-go/boost/tree/main/encoding2/README.md)
- `encryption` [加密](https://github.com/sandwich-go/boost/tree/main/encryption/README.md)
- `geom` [地理位置](https://github.com/sandwich-go/boost/tree/main/geom/README.md)
- `goformat` [Golang格式化](https://github.com/sandwich-go/boost/tree/main/goformat/README.md)
- `hrff` [Human Readable Flags and Formatting](https://github.com/sandwich-go/boost/tree/main/hrff/README.md)
- `httputil` [HTTP工具](https://github.com/sandwich-go/boost/tree/main/httputil/README.md)
- `paniccatcher` [panic工具](https://github.com/sandwich-go/boost/tree/main/paniccatcher/README.md)
- `retry` [重试](https://github.com/sandwich-go/boost/tree/main/retry/README.md)
- `version` [程序版本](https://github.com/sandwich-go/boost/tree/main/version/README.md)
- `xcmd` [命令行/ENV参数](https://github.com/sandwich-go/boost/tree/main/xcmd/README.md)
- `xcopy` [深拷贝](https://github.com/sandwich-go/boost/tree/main/xcopy/README.md)
- `xerror` [error wrapper](https://github.com/sandwich-go/boost/tree/main/xerror/README.md)
- `xhash` [hash算法](https://github.com/sandwich-go/boost/tree/main/xhash/README.md)
- `xio` [异步io](https://github.com/sandwich-go/boost/tree/main/xio/README.md)
- `xip` [ip/port辅助函数](https://github.com/sandwich-go/boost/tree/main/xip/README.md)
- `xmap` [map辅助函数](https://github.com/sandwich-go/boost/tree/main/xmap/README.md)
- `xmath` [math辅助函数](https://github.com/sandwich-go/boost/tree/main/xmath/README.md)
- `xpanic` [panic辅助函数](https://github.com/sandwich-go/boost/tree/main/xpanic/README.md)
- `xproc` [command辅助函数](https://github.com/sandwich-go/boost/tree/main/xproc/README.md)
- `xrand` [随机辅助函数](https://github.com/sandwich-go/boost/tree/main/xrand/README.md)
- `xslice` [切片辅助函数](https://github.com/sandwich-go/boost/tree/main/xslice/README.md)
- `xstrings` [字符串辅助函数](https://github.com/sandwich-go/boost/tree/main/xstrings/README.md)
- `xsync` [同步辅助函数](https://github.com/sandwich-go/boost/tree/main/xsync/README.md)
- `xtag` [sturct tag辅助函数](https://github.com/sandwich-go/boost/tree/main/xtag/README.md)
- `xtemplate` [tempalte辅助函数](https://github.com/sandwich-go/boost/tree/main/xtemplate/README.md)
- `xtime` [time辅助函数](https://github.com/sandwich-go/boost/tree/main/xtime/README.md)
- `z` [辅助函数](https://github.com/sandwich-go/boost/tree/main/z/README.md)

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