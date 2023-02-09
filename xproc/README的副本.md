# xproc

`command` 辅助函数

- `process` 启动、关闭、`kill`、`Signal` 管理
- `process` 集管理

# 例子
```go
tmpFile := filepath.Join(os.TempDir(), "test.sh")
err := xos.FilePutContents(tmpFile, []byte("echo \"GOT ME $1\""))
if err != nil {
    panic(err)	
}
var stdOut string
stdOut, err = ShellRun(tmpFile, WithArgs("1.2.0"))
if err != nil {
    panic(err)
}
fmt.Println(stdOut)
```
Output:
```text
GOT ME 1.2.0
```