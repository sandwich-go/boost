# module

`module` 管理工具

- 管理自定义的 `module`
- 按注册顺序启动 `module`
  - 依次调用 `module` 的 `OnInit` 函数
  - 依次调用 `module` 的 `Run` 函数
- 按注册反序关闭 `module`
  - 依次调用 `module` 的 `OnClose` 函数

# 例子
```go
type infrastructure struct {
    handler   Handler
    startedCh chan struct{}
}

func newInfrastructure(handler Handler) *infrastructure {
    return &infrastructure{startedCh: make(chan struct{}), handler: handler}
}

func (s *infrastructure) OnInit() {
    s.handler.infrastructureInitialize()
    close(s.startedCh)
}

func (s *infrastructure) OnClose() {
    s.handler.infrastructureDestroy()
}

func (s *infrastructure) Run(closeChan chan struct{}) { <-closeChan }
func (s *infrastructure) Name() string                { return "module-infrastructure" }

func main() {
    handler := newHandler()
    module.Register(newInfrastructure(handler))
    module.Run()
}
```