# paniccatcher

`panic` 工具

# 例子
```go
Do(func() {
    fmt.Println("Doing something...")
    panic("Something wrong happened!")
}, func(p *Panic) {
    fmt.Println("Caught a panic:", p.Reason)
})

AutoRecover("something", func(){
    // do something
    // if panic, auto execute this function continue.
})
```