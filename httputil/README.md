# httputil

`HTTP` 工具

- 便捷返回 `json`、`string`、`bytes` 类型的 `GET` 数据
- 支持带缓存的域名解析

# 例子
```go
str, err := String("https://www.baidu.com/")
xpanic.WhenError(err)
fmt.Println(str)

dc := dns.NewCache()
client := &http.Client{
    Transport: &http.Transport{
        DialContext: dc.GetDialContext(),
    },
}
c := New(client)
str, err = c.String("https://www.baidu.com/")
xpanic.WhenError(err)
fmt.Println(str)
```