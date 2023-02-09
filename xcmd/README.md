# xcmd

命令行参数或Env参数

# 例子
```go
var name1 = fmt.Sprintf("%sdebug", GetFlagPrefix())
var name2 = fmt.Sprintf("%sp1", GetFlagPrefix())
var name3 = fmt.Sprintf("%sp2", GetFlagPrefix())
Init(fmt.Sprintf("--%s=true", name1), fmt.Sprintf("--%s", name2), "test", fmt.Sprintf("--%s", name3))

fmt.Println(IsTrue(GetOptWithEnv(name1)))
fmt.Println(GetOptWithEnv(name2))
fmt.Println(IsTrue(GetOptWithEnv(name3)))
```

Output:
```text
true
test
false
```