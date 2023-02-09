# boost

`boost` 工具箱

- `annotation` [注释解析器](https://github.com/sandwich-go/boost/tree/main/annotation)
- `compressor` [解压缩器](https://github.com/sandwich-go/boost/tree/main/compressor)
- `encoding2` [编码解码器](https://github.com/sandwich-go/boost/tree/main/encoding2)
- `encryption` [加密](https://github.com/sandwich-go/boost/tree/main/encryption)
- `goformat` [Golang格式化](https://github.com/sandwich-go/boost/tree/main/goformat)
- `hrff` [Human Readable Flags and Formatting](https://github.com/sandwich-go/boost/tree/main/hrff)
- `httputil` [HTTP工具](https://github.com/sandwich-go/boost/tree/main/httputil)
- `paniccatcher` [panic工具](https://github.com/sandwich-go/boost/tree/main/paniccatcher)
- `retry` [重试](https://github.com/sandwich-go/boost/tree/main/retry)
- `version` [程序版本](https://github.com/sandwich-go/boost/tree/main/version)
- `xcmd` [命令行/ENV参数](https://github.com/sandwich-go/boost/tree/main/xcmd)
- `xcopy` [深拷贝](https://github.com/sandwich-go/boost/tree/main/xcopy)
- `xerror` [error wrapper](https://github.com/sandwich-go/boost/tree/main/xerror)
- `xhash` [hash算法](https://github.com/sandwich-go/boost/tree/main/xhash)
- `xio` [异步io](https://github.com/sandwich-go/boost/tree/main/xio)
- `xip` [ip/port辅助函数](https://github.com/sandwich-go/boost/tree/main/xip)
- `xmap` [map辅助函数](https://github.com/sandwich-go/boost/tree/main/xmap)
- `xmath` [math辅助函数](https://github.com/sandwich-go/boost/tree/main/xmath)
- `xpanic` [panic辅助函数](https://github.com/sandwich-go/boost/tree/main/xpanic)
- `xproc` [command辅助函数](https://github.com/sandwich-go/boost/tree/main/xproc)
- `xrand` [随机辅助函数](https://github.com/sandwich-go/boost/tree/main/xrand)
- `xslice` [切片辅助函数](https://github.com/sandwich-go/boost/tree/main/xslice)
- `xstrings` [字符串辅助函数](https://github.com/sandwich-go/boost/tree/main/xstrings)
- `xsync` [同步辅助函数](https://github.com/sandwich-go/boost/tree/main/xsync)
- `xtag` [sturct tag辅助函数](https://github.com/sandwich-go/boost/tree/main/xtag)
- `xtemplate` [tempalte辅助函数](https://github.com/sandwich-go/boost/tree/main/xtemplate)
- `xtime` [time辅助函数](https://github.com/sandwich-go/boost/tree/main/xtime)
- `z` [辅助函数](https://github.com/sandwich-go/boost/tree/main/z)

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