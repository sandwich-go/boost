# AGENTS.md

给在本仓库做改动的 agent（人或 AI）读的"上岗手册"。主文档是
[`README.md`](./README.md)；本文件只补齐"做事前必须知道的约束 / 陷阱 /
工作流"——读完能避免 80% 的返工。

本文风格参考姊妹仓 myproxy / dataserver 的 AGENTS.md，但 boost 是**通用
工具库**，性质不同：内容已按工具库视角重写，**不要照搬业务系统的约束**
（如 SQL 注入、Close 协议、TLS 透传等本仓不存在）。

---

## 1. 这个仓是什么

**boost = sandwich-go 系列的通用工具库**。47 个顶级子包 / ~550 个 Go 文件，
被 `myproxy` / `dataserver` / `cmd/*` 等业务仓 import 消费。

核心定位：

- **稳定 API 优于性能优化**——下游业务版本号 pin 到 boost minor，公开 API 改名
  / 删除 / 行为变更要走 minor 版本（1.4 vs 1.3），且发版前必检"还有没有下游
  仍在调老 API"
- **每个子包独立可用**，互相依赖尽量少；`xmath`/`xslice`/`xmap` 这种 leaf 工具
  包不能反向 import `httputil` 之类的 heavy 包
- **泛型已可用**（go.mod `go 1.24.0`）；新代码优先泛型实现，老的多类型生成
  代码（`internal/template2/`）逐步退役（已在 1.4/develop 完成 xmath /
  xslice / xmap 三个包的迁移，删除 ~16,300 行生成代码）

下游业务最常用的子包（需要特别保稳定）：

- `xmath / xmap / xslice / xstrings / xconv / xio / xrand` — leaf 工具
- `xerror` — 带栈错误（`xerror.Errorf` / `xerror.Array`），下游依赖错误码语义
- `xtime` — `Periodic / NowFunc / dispatcher`，2026-05 增加 `Periodic` 原语后
  下游 lru / myproxy refresh 路径已切过去
- `retry` — `retry.Do(fn, opts...)`，下游 download / dataserver / pmt 都在用
- `module` — module 级生命周期管理（master/agent），cmd 入口的标配
- `xcontainer/{smap,syncmap,sset,redblacktree,sortedmap}` — 高频并发容器
- `lru` — 分片 LRU，业务侧热点缓存
- `xpool` — bytes pool / sync.Pool 封装

---

## 2. 开工前先跑

```bash
# 工具一次性安装（optiongen / stringer / mockgen / gotemplate / golangci-lint / govulncheck / benchstat）
make tools

# CI 等价的本地准入：vet + build + test + race（lint 当前未钉，见 §4.5）
make ci
```

**Go 版本**：`go.mod` 写 `go 1.24.0`。本地 toolchain 比这版本新没事，但**不要在
go.mod 里手动把 `go 1.24.0` 往下回退**，会让 1.21+ 才有的 `cmp.Ordered` /
`maps` / `slices` 标准库不可用。

**CI 镜像**：`.github/workflows/ci.yml` 跑 `vet → build → test → race`，全绿才
能合入主干。详细命令见 §2.1。

### 2.1 日常命令清单

| 命令 | 什么时候跑 |
|---|---|
| `make ci` | 提交前 / PR 准入。等价于 CI 上 `vet + build + test + race + vuln` |
| `make test` | 局部改动后，比 `make ci` 快（无 race） |
| `make test_race` | 改了并发原语 / 容器 / 池子时，必跑 |
| `make lint` | **当前未钉为硬约束**（见 §4.5），但本地建议跑过；新增告警零容忍 |
| `make cover` | 测试覆盖率变动时，产物 `coverage.out` |
| `make bench_diff` | 动了 hot path（`xrand` / `xpool` / `xcontainer/*` / `xencoding/*`）时本地跑；等价 `bench` + `benchstat`，对比入仓基线 `benchdata/main.txt` |
| `make bench_refresh` | 主干 perf 改动 land 后刷新入仓基线，**独立 chore commit**，不混到其他改动 |
| `make gen` | 改了 `optiongen` / `stringer` / `gotemplate` / `mockgen` 注解的源文件后跑；diff 必须随同源改动一起提交 |
| `make vuln` | 升级依赖前 / 主干同步上游 stdlib CVE 时跑 |

---

## 3. 目录速查

47 个顶级子包，按"使用频度 / 关键性"分组：

```
# === leaf 工具（被大量 import；改动需 PoC + bench） ===
xmath/        Max/Min/Abs/EffectZeroLimit (泛型) + Float* EPSILON 比较
              1.4 已完成 template2 → 泛型迁移，删 44 个具名函数
xmap/         Equal / WalkMapDeterministic / Diff / ToMap (泛型)
              1.4 已完成 template2 → 泛型迁移，删 351 个具名函数
xslice/       Contain / SetAdd / Walk / RemoveRepeated / RemoveEmpty / Shuffle
              + ContainEqualFold / AddPrefix / AddSuffix (string 专属)
              1.4 已完成 template2 → 泛型迁移，删 107+26 个具名函数
xstrings/     字符串处理（Camel / Snake / Trim 等）
xconv/        类型转换（int↔string / bool↔int / 时间格式等）
xrand/        FastRand + String / StringWithTimestamp（v1.4 已优化 -3 alloc）
xio/          异步 IO（Reader / Writer wrappers）
xhash/        fnv / nhash / hash14v / md5

# === 错误处理（下游依赖错误码 / 栈语义） ===
xerror/       Errorf / NewText / NewCode / Wrap + Array (错误聚合) + Stack 抓栈
              改这里要看 §11.4：栈抓取是错误产生时的真大头（>fmt.Sprintf）

# === 时间 / 调度 ===
xtime/        dispatcher (timer 调度) + cron + Periodic (1.3 新增) + NowFunc 抽象
              hot path 慎改，refresh 类业务依赖；race 必跑
xtime/cron/   crontab 表达式

# === 并发 / 容器 ===
xsync/        wg / 超时锁 (1.3 新增) / 各种同步原语
xchan/        ringbuf-based unbounded chan
xcontainer/   sortedmap / syncmap / smap / sset / sarray / slist /
              redblacktree / ringbuf / fanin / 各种 generated 容器
xparallel/    并发执行原语（错误聚合）

# === 池子 / 限流 ===
xpool/        SyncBytesPool (slab) / sync.Pool 封装 / goroutine pool / reference
              改 Alloc/Free 性能必 PoC + bench，见 xpool/buffer_bench_test.go 头部
              "已 PoC 不优"决策记录
ratelimiter/  令牌桶 / 漏桶
singleflight/ 函数级请求去重

# === 缓存 ===
lru/          分片 LRU；workerPerEngine 1.3 切到 xtime.PeriodicWithShutdown

# === 网络 / IO / 编解码 ===
httputil/     HTTP client / dns 策略
xencoding/    msgpack / json / pbjson / protobuf / encrypt / compressor
              这几个 codec 是 hot path，bench 全配齐
xcompress/    gzip / snappy 封装（用 sync.Pool）
xcrypto/      算法封装（aes / mask / hash）

# === 系统级 / 入口 ===
xos/          文件 / 目录 / 进程辅助
xproc/        子进程封装
xcmd/         命令行 / ENV 参数
xip/          IP / port 辅助
module/       module 生命周期 (master/agent)，cmd 入口标配
boost.go      LogXxx 简易 logger 门面（见 §11.6）

# === 调试 / 错误恢复 ===
xpanic/       panic 恢复 + Try / Catch
xdebug/       依赖检查 (registerDependency) — 唯一未泛型化的 template2 用例

# === Misc / 业务子域 ===
misc/         xtemplate / annotation / goformat / cloud / xgen
geo/          地理 / 网格
graph/        图算法
plugin/       插件加载
validator/    校验器 (基于 go-playground/validator)
xvalidator/   protovalidate-go 封装
humanize/     human-readable 格式化
isnil/        nil 判断
types/        bignum / 扩展类型
version/      程序版本元数据
xetl/         ETL / scan helpers
xexp/         实验性子包
xtest/        quantile (P50/P99 等)
xcopy/        深拷贝
z/            纯 stdlib helper（hash / mono / wall / rand / conv），无依赖

# === 其它 ===
internal/log/ 占位 Logger 接口（被 boost.go 用），下游通过 InstallLogger 注入
internal/template2/   见 §6.1（已废弃模式，xmath/xmap/xslice 已删除；xdebug 保留）
```

**1.4 阶段重要变化**（参考 1.3/release..1.4/develop diff）：

- `xmath / xmap / xslice` 的 `internal/template2/` + `internal/export.go` 已物
  理删除；老的 502 个具名函数（`MaxInt32` / `EqualStringStringMap` /
  `IntsContain` 等）全部删除——下游需要改名到泛型版（迁移表见 §6.1）
- `xrand.StringWithTimestamp` 实现重写，砍 3 alloc / -47% ns（PoC + bench
  数据见 commit message + `xrand/string_bench_test.go` 头部基线）

---

## 4. 仓库级硬约束（改动必须满足）

### 4.1 `make ci` 全绿

`make ci = vet + build + test_race + vuln`。其中 **build 是 CI 唯一硬
门槛**（`.github/workflows/ci.yml`）；其余都是软门槛（continue-on-error）。
任何 PR 合入主干前本地都应该跑过一次 `make ci` 看完整信号。

### 4.1.1 当前软门槛清单与未来硬钉路径

boost 仓 1.4/develop HEAD 上有以下 pre-existing 噪声 / 失败，需要单独 PR
系统性清理后再硬钉：

| CI job | 当前状态 | pre-existing 问题 |
|---|---|---|
| `vet` | 软门槛 | ~20 处 noise（`misc/hrff/hrff_test.go` Example 命名 / `xsync/cond_test.go` non-test goroutine 调 `t.Fatal` / `z/conv.go` `reflect.StringHeader` misuse / `xerror/xerror_x_test.go` cancel func discarded） |
| `lint` | 软门槛 | 没 `.golangci.yml`，默认规则集会出 200-300 告警（仓库历史长，从未钉过 lint） |
| `test` | 软门槛 | 当前无 pre-existing FAIL（misc/cloud `TestCloud` 仅在本地有 RELEASE_CLOUD_KEY/SECRET env 时才连真 AWS，CI 无 env 自动跳过；xpanic 一组测试在本轮已修复）。下一步钉硬门槛前需先确认 lint / vet 噪声清完 |
| `race` | 软门槛 | 与 test 同；race 模式下没有额外的 pre-existing 问题需先解决 |
| `vuln` | 软门槛 | 没 0 affecting 承诺基线 |
| `build` | **硬门槛** | 全仓 `go build ./...` 通过 |

**渐进式硬钉路径**（每步独立 PR）：

1. 系统性修 vet noise → CI vet 改 hard gate
2. 修 `misc/cloud` / `xpanic` 测试 → test/race 改 hard gate
3. 加 `.golangci.yml` v2 schema → 修 lint 告警 → lint hard gate
4. 升 deps 把 vuln 清零 → vuln hard gate

每步都是独立 PR，不混进功能改动。

#### §9.2 关联：xpanic %w/%v 历史教训

xpanic/panic_when.go 历史曾发生 `fmt.Errorf` → `fmt.Sprintf` 改写时漏改
`%w` → `%v`（commit f7dd56a），vet 一直报 build error 让 xpanic 测试无
法编译。修复同时补 `WhenHereNotNil` 端到端断言（panic value 字面格式
精确比较），确保未来再回归这种"verb 不匹配"问题时测试层面也能捕获，
不再仅依赖 vet（万一未来 lint/vet 配置改动就漏）。详见 §9.2。

### 4.2 Go 版本钉 1.24.0

`go.mod` 第 3 行 `go 1.24.0`。**不要降到 1.20 / 1.21**——1.21 引入的
`cmp.Ordered` / 标准库 `slices` / `maps` 已被本仓泛型代码 import（看
`xmath/xmath.go` / `xmap/xwalk.go` 头部 import）。

升 Go 版本时单独走 PR：升 `go.mod` → 跑 `make ci` → 跑 `make bench_diff`
确认主要 hot path 性能不退化 → 提交。

### 4.3 公开 API 只增不减（minor 内）

下游业务侧（myproxy / dataserver 等）的 `go.mod` 通常 pin 到 boost 的
**minor 版本**（如 `v1.3.x`），minor 内 API 增加可以但**删除 / 重命名 /
签名变更不允许**——这些是 break change，必须升 minor（1.3 → 1.4）。

**1.4 已经发生的 break change**（见 §6.1 迁移表）：
- `xmath`：`MaxInt32` / `MinFloat64` / `AbsInt` / `EffectZeroLimitInt32` 等 44
  个删除，统一改名到 `Max[Ordered]` / `Min[Ordered]` / `Abs[Signed]` /
  `MaxFloat[~float]` / `EffectZeroLimit[Signed]` 等泛型版
- `xslice`：`StringsContain` / `IntsRemoveRepeated` / `Float32sShuffle` 等
  107 个删除 + 26 个未导出 helper 删除，改名到 `Contain` / `SetAdd` /
  `Walk` / `RemoveRepeated` / `RemoveEmpty` / `Shuffle` 泛型版
- `xmap`：`EqualStringStringMap` / `WalkInt32Float64MapDeterministic` 等
  351 个删除，改名到 `Equal` / `WalkMapDeterministic` 泛型版

升 1.5 之前**不要再加新的 break change**——一次升一个台阶，下游迁移成本可控。

### 4.4 `go test -race` 全绿

并发包（`xtime` / `xsync` / `xchan` / `xcontainer/syncmap` / `xpool` / `lru` /
`ratelimiter` / `singleflight`）任何改动必须 `make test_race` 通过。

某些已知抖动（Apple M2 Pro 上 `xtime.TestTick` 在 `go test ./...` 高并发时偶
发"close of closed channel"，但单独跑 `go test ./xtime/` 总绿）记到 §11，
不堵塞主干。

### 4.5 lint 暂未钉为硬约束

boost 仓**没有 `.golangci.yml`**，`make lint` target 存在但 CI 不跑。原因：
仓库历史长，全量打开 lint 会出几百个告警，集中清理需要单独 PR。

**短期态度**：

- 新代码本地 `make lint` 应零告警（自觉）
- CI 不阻断 lint 失败（`continue-on-error: true`），允许"先合入再清"
- 后续若决定钉 lint，单独 PR：加 `.golangci.yml` v2 schema → 修全部告警 → CI
  把 lint 从 soft gate 改 hard gate

**中期路线图**（待安排）：参考 myproxy 的 `.golangci.yml` v2 schema 接入 +
`internal/template2/` 类已删除路径不再有豁免压力。

### 4.6 `interface{}` 替换为 `any`

非生成代码、非 fork 代码，新 API 用 `any`。旧代码用 `interface{}` 的逐步替换
（见 commit `ebde337` "refactor"），混用不影响行为但风格不一致。

### 4.7 不要回写已删除的 template2 路径

`xmath / xmap / xslice / internal/template2/` 已物理删除，**generator 也删了**
（`internal/export.go`）。新增多类型工具用泛型，**不要回写 template-based
代码生成**。

唯一保留：`xdebug/internal/template2/dependency_template.go`——它生成的不是
"多类型工具"而是 dependency 注册表，不是泛型可表达的场景。

---

## 5. 测试 / bench / pprof 写法

### 5.1 goconvey 优先

仓内 `_test.go` 习惯用 `goconvey` 嵌套 `Convey("desc", t, func(){ ... })`，写
新测试照搬。`testify/require` 偶有出现但不是主流。

### 5.2 bench 入仓作为基线

性能敏感包要写 `BenchmarkXxx`，并把当前数值作为**回归监测基线**。已有的
基线点：

| 包 | bench 文件 | 关键 bench |
|---|---|---|
| `xrand` | `xrand/string_bench_test.go` | `String_*` / `StringWithTimestamp_*`；头部注释带 baseline 数值 |
| `xpool` | `xpool/buffer_bench_test.go` | `SyncBytesPool_Alloc_*`；头部含"已 PoC 不优"决策记录 |
| `xpool` | `xpool/goroutine_benchmark_test.go` | goroutine pool 性能基线 |
| `xtime` | `xtime/cop_test.go` / `periodic_test.go` | dispatcher / Periodic 性能基线 |
| `xencoding/protobuf` | `codec_benchmark_test.go` | proto codec 性能基线 |
| `xencoding/msgpack` | `msgpack_benchmark_test.go` | msgpack 性能基线 |
| `xencoding/json` | `json_benchmark_test.go` | json 性能基线 |
| `xcrypto/algorithm/aes` | `aes_test.go` | AES 性能基线 |

**改 hot path 必须 bench**，规约见 [`AGENTS.md` 全局版 §3.4](~/.config/opencode/AGENTS.md)：

> 任何"性能优化值不值得做"的决策都需要 PoC + bench 实测，**禁止只凭推断给
> 结论**。

`xrand/string_bench_test.go` 头部 + `xpool/buffer_bench_test.go` 头部都有完整
的"PoC + bench → 决策"记录范例，新加 bench 时参考这种 doc 化方式。

### 5.3 pprof 抓 GC 热点

GC 嫌疑场景下抓 alloc profile：

```bash
go test -run=^$ -bench=BenchmarkXxx -benchmem -memprofile=mem.out ./pkg/
go tool pprof -alloc_objects mem.out
# (pprof) top
# (pprof) list <FuncName>
```

仓内 `xrand.StringWithTimestamp` 优化（commit 待查）就是这流程：bench 看到
4 alloc/op，pprof 定位是外层 fmt.Sprintf + 内层 strings.Join，重写后 1 alloc/op。

### 5.4 race 是硬门槛

任何并发原语改动必须 `make test_race` 全绿。CI 也跑同一份，runner 抖动允许
retry 1 次。

### 5.5 `go test ./...` vs 单包跑

仓内某些抖动只在 `go test ./...`（高并发跨包）时出现，单包跑稳定（如
`xtime.TestTick`）。如果 CI 偶发挂某个包但本地单跑稳，**先重跑一次**，再
判断是抖动还是真问题。

### 5.6 不要给生产代码加 testing seam

除非业务封装太死必须暴露，这种修改单独提一个 commit 标注用途。一般白盒测试
放同包 `_test.go`，直接访问包私有字段。

---

## 6. 泛型迁移与技术金

### 6.1 `internal/template2/` → 泛型迁移（1.4 已完成 3 个包）

**背景**：1.3 之前，`xmath` / `xmap` / `xslice` 用 `internal/template2/`
+ `internal/export.go` 跑 generator，从 Go template 字符串生成多类型具名函数
（`MaxInt32` / `EqualStringStringMap` / `IntsContain` 等共 502 个）。Go 1.18+
有泛型后这模式过时。

**1.4 已完成**：删除 502 个具名函数 + generator + template2 目录，全部走
泛型 API。共删除 ~16,300 行代码，新增 ~220 行（含手写测试）。

**迁移表**（破坏性变更，下游需要改）：

```
旧 API → 新 API（boost 泛型版 / Go 标准库）

xmath:
  MaxInt32(a, b)        → xmath.Max(a, b)            // [T cmp.Ordered]
  MaxFloat64(a, b)      → xmath.MaxFloat(a, b)       // [T ~float]，IEEE 754 NaN 传染
  AbsInt(v)             → xmath.Abs(v)               // [T constraints.Signed]
  AbsFloat32(v)         → xmath.AbsFloat(v)          // [T ~float]，保 EPSILON 近似判负
  EffectZeroLimitInt(v, c)    → xmath.EffectZeroLimit(v, c)
  EffectZeroLimitFloat32(v, c) → xmath.EffectZeroLimitFloat(v, c)
  保留：Float32Equals / Float64Equals / IsZeroFloat* / IsBelowZeroFloat* / EPSILON*

xslice:
  IntsContain(s, v)         → xslice.Contain(s, v)        // 或标准库 slices.Contains
  StringsSetAdd(s, v...)    → xslice.SetAdd(s, v...)
  IntsWalk(s, fn)           → xslice.Walk(s, fn)
  IntsRemoveRepeated(s)     → xslice.RemoveRepeated(s)    // ⚠ 与 slices.Compact 不同：保序全局去重
  IntsRemoveEmpty(s)        → xslice.RemoveEmpty(s)
  Float64sShuffle(s)        → xslice.Shuffle(s)
  StringsContainEqualFold   → xslice.ContainEqualFold     // string 专属
  StringsAddPrefix          → xslice.AddPrefix
  StringsAddSuffix          → xslice.AddSuffix

xmap:
  EqualStringStringMap(a, b)            → xmap.Equal(a, b)
  WalkInt32Float64MapDeterministic(...) → xmap.WalkMapDeterministic(...)
```

**语义微差异**（下游迁移时注意）：

- `xmath.Max[Ordered]` 在 NaN 上**不再传染**（用 `if a > b` 比较，NaN 比较都
  false 走 else 返回 b）；要 IEEE 754 NaN 传染语义改用 `xmath.MaxFloat`
- `xmath.AbsFloat` 的 EPSILON 近似判负保留（`v=0` / `v=-0.5*EPSILON` 都返回 -v）
- `xmap.Equal` 区分 nil 与空 map（`nil != map[K]V{}`），与 Go 1.21+ 标准库
  `maps.Equal`（视 nil == empty）**不一致**，是历史 boost API 行为
- `xmap.WalkMapDeterministic` 的 V 约束放宽到 `any`（原老 API 允许 V=interface{}）

### 6.2 标准库 vs boost 泛型版

Go 1.21+ 标准库 `slices` / `maps` 已覆盖大量 boost 历史功能：

| 标准库 | 等价 boost API | 何时用谁 |
|---|---|---|
| `slices.Contains` | `xslice.Contain` | 标准库优先（轻量、无 import 成本） |
| `slices.ContainsFunc` | `xslice.ContainEqualFold` 是它的特化 | 标准库优先 |
| `slices.Compact` | **不等价**于 `xslice.RemoveRepeated` | Compact 只去相邻重复，RemoveRepeated 保序全局去重 |
| `slices.Equal` | — | 标准库 |
| `maps.Equal` | `xmap.Equal` | **语义不同**：标准库视 nil == empty，boost 区分 |
| `maps.Keys` | — | 标准库（返回 iter.Seq[K]） |
| `slices.SortFunc` | — | 标准库 |
| `cmp.Compare` / `cmp.Ordered` | — | 标准库（boost 自己也基于这些） |

**新代码鼓励用标准库**，除非语义/性能确有差异（如 `xmap.Equal` 的 nil/empty
区分是 boost 既定行为，下游已依赖）。

---

## 7. commit / push / tag 规约

### 7.1 分支策略

- **`main`** — 长期分支，大版本发布的来源
- **`X.Y/develop`**（如 `1.4/develop`）— 当前 minor 的开发主干，新 commit 推
  这里；功能分支从它切出、合回
- **`X.Y/release`**（如 `1.3/release`）— 已发版分支，只接受从 develop
  fast-forward / merge
- **`X.Y/feat/<topic>`**（如 `1.3/feat/lru-op`）— 大特性长期分支

切新版本（参考切 1.4/develop 经验）：

```bash
git fetch origin
git checkout -b 1.4/develop origin/1.3/release
git push -u origin 1.4/develop
```

只切 develop，不预先切 release——release 在 minor 真正发版前不存在。

### 7.2 commit 规约

走 **conventional commits 中文**风格（看 `git log` 前 20 条）：

- `feat(xtime): 新增 Periodic 周期任务原语`
- `fix(misc): 修复 pointerShift 在 -race 下触发 checkptr 报错`
- `perf(xrand): StringWithTimestamp 避免 fmt.Sprintf 减 3 alloc`
- `refactor(xmap): template2 → 泛型，删 351 个具名函数`
- `chore(sem): make changelog`
- `test(xrand): 增加 StringWithTimestamp baseline bench`
- `docs(AGENTS): 新增上岗手册`

commit 拆分原则：

- 改动逻辑上独立的就**分多个 commit**，一个主题一条，便于 revert
- 不要在一个 commit 里混 "功能修改" + "格式化"
- `make gen` 产生的 diff 单独一个 commit（和对应 conf.go 改动同 commit 也行）
- bench 基线刷新（`make bench_refresh`）独立 chore commit
- **perf 优化和功能迁移分两个主题 commit**，不混进同一 PR

### 7.3 打 tag / 升版本是显式指令，不是隐含动作

`make version` (= `protokitgo sem release`) 会推下一个 tag 并自动 push 到
origin。tag 是对外语义承诺（CI artifact / 部署系统 pin / changelog 起点），
**只增不减**——push 后回滚成本极高。

**约束**：agent / AI 协作者**只在用户明确说出下列指令时才跑 `make version`**：

- "打 tag" / "打个 tag"
- "升版本" / "打个版本" / "打个新的版本"
- "发布 vX.Y.Z" / "release"

**不触发 tag 的动词**（哪怕推完代码也不要顺手打）：

- "修复" / "提交" / "push" / "同步" / "记录" / "整理" / "重构"
- "更新文档" / "改 X" / "加 Y"

不确定时**问一句不丢人**："本次修复要打 tag 吗？"

详细反例与判断准则见 [`~/.config/opencode/AGENTS.md` §1.2](~/.config/opencode/AGENTS.md)
（个人全局规约）。

### 7.4 changelog 由 protokitgo 生成

仓内 `CHANGELOG-0.1.md` / `CHANGELOG-1.3.md` 由 `protokitgo sem changelog`
生成，**不要手改**。`.sembumprc.yml` 控制生成参数。

`chore(sem): make changelog` 是 `protokitgo` 自动产生的 commit message，
不需要人工写。

### 7.5 force push 守则

- `--force-push` 到 main / X.Y/release **必须警告用户**且仅在用户明确同意时执行
- 一律用 `--force-with-lease`，不裸 `--force`
- 已 push 到 remote 的 commit 不要 amend（除非用户明确要求 force push）

---

## 8. 高风险改动清单

按"动它必须非常小心"排序：

1. **`xtime/dispatcher.go`** — 核心 timer 调度。`tickHosting` / `TriggerTickFuncs`
   改顺序 / 改锁就会让下游 lru / refresh 路径漏 tick 或 panic。改前必跑
   `make test_race ./xtime/` + 多次 `go test ./...` 看 `TestTick` 是否抖动加重
2. **`xerror/xerror_stack.go` / `xerror/array.go`** — 错误栈抓取。下游错误码
   语义靠 `*Error` struct + `Code()` 函数透视，结构体字段重排 / 重命名都
   是 break change
3. **`retry/retry.go`** — `BackOffDelay = Delay << n` 是字面公式。
   下游业务依赖 sleep 序列总长（如 download 的 Limit=5 总 sleep 30s）。
   改公式 = 改下游超时配置语义；公式改动必须 PoC + bench 验证以及 doc 更新
4. **`xpool/buffer.go` SyncBytesPool**`Alloc` / `Free`** — slab pool 的 size
   匹配算法。线性扫已 PoC 不优（见 `xpool/buffer_bench_test.go` 头部决策
   记录），但若动 chunk size 计算或 sync.Pool 改 sync.Map 等本质改动，要
   `make test_race` + bench
5. **`xcontainer/syncmap/*.go`** — 桶分片并发 map。键空间 hash 分布改了会
   让 race 测试假阴。改前看 `xcontainer/syncmap/bucketmap.go` 的 keySlicePool
   实现
6. **`lru/clean_worker.go` workerPerEngine** — 1.3 已切到
   `xtime.PeriodicWithShutdown`（commit `aefae5f`）；改这块必读那 commit
   message 的 rationale，单实例稳态 0 timer 分配是关键不变量
7. **`xsync/*.go`** — 各种同步原语。1.3 新增"超时锁"（commit `06fc89f`）
   后多个下游业务用，改 API 是 break change
8. **`xpanic/*.go`** — `panic` recovery 路径。1.3 修过 `Errorf` → `Sprintf`
   bug（commit `f7dd56a` / `4bf8c4e`），改 panic_when 系列函数前看那两个 commit
9. **`xrand/string.go`** — 1.4 刚优化的 `StringWithTimestamp`（PoC + bench
   见 `string_bench_test.go`）。再动这个文件必须不让 baseline 倒退
10. **`internal/log/log.go` Logger 接口** — 下游通过 `boost.InstallLogger`
    注入。接口加方法是弱 break（要求下游补实现），改方法签名是 hard break
11. **`module/master.go`** — module 生命周期。`master.Start` / `Close` 顺序敏
    感，下游 cmd 入口都依赖。改前必跑下游 cmd 的端到端

这几处改动建议在 PR 描述里显式说 "trigger full regression: race + bench"。

---

## 9. 常见陷阱（踩过的坑）

> boost 没有 myproxy 的 `docs/bug-history.md` 体系；下面是已知踩坑的索引式
> 短条目，未来如有大量积累再单独建文档。

### 9.1 `xtime.TestTick` 在 `go test ./...` 偶发挂

症状：`go test ./...` 高并发跨包时挂 `panic: close of closed channel` 在
`xtime/dispatcher.go:100`，但单独跑 `go test ./xtime/` 总绿。

判断：测试代码自身的 panic（test fixture 双重 close），非生产代码 bug。改其
他包后跑 `./...` 命中这个抖动**不是你引入的**——单包重跑确认。

### 9.2 `xpanic.WhenErrorAsFmtFirst` 用 `Sprintf` 而非 `Errorf`

`fmt.Sprintf` 不支持 `%w`（error wrapping），用了等于把错误转成纯字符串。
1.3 修过（commit `4bf8c4e` / `f7dd56a`）。新增 `panic(...)` 路径用
`fmt.Sprintf` 时确认没用 `%w`。

### 9.3 `boost.LogXxxf` eager fmt.Sprintf

```go
func LogInfof(format string, a ...interface{}) { LogInfo(fmt.Sprintf(format, a...)) }
```

调用方即使 logger 决定丢弃 Info 级别，**fmt.Sprintf 已经 alloc + 格式化**。
hot path 慎用 `LogXxxf`，改用注入的 zerolog / zap 直接调（它们有 lazy 评估）。
仓内只在 `module/master.go` 启动/关闭时调用，不是 hot path。

### 9.4 `internal/template2/` 是已废弃模式

不要回写 generator 生成代码。多类型工具用泛型；唯一保留 `xdebug/internal/
template2/dependency_template.go` 因为不是泛型可表达场景。

### 9.5 `xmap.Equal` 与 `maps.Equal` 语义不同

- `xmap.Equal(nil, map[K]V{}) == false`（区分 nil 与空 map）
- Go 标准库 `maps.Equal(nil, map[K]V{}) == true`（视为相等）

下游迁移时注意取舍：要 nil/empty 区分继续用 `xmap.Equal`，要标准语义用
`maps.Equal`。

### 9.6 `xslice.RemoveRepeated` 与 `slices.Compact` 语义不同

- `xslice.RemoveRepeated([1,2,1]) == [1,2]`（保序全局去重）
- `slices.Compact([1,2,1]) == [1,2,1]`（只去**相邻**重复）

下游迁移时**别 `sed s/RemoveRepeated/slices.Compact/`**——结果不一样。

### 9.7 `make.sh` 是 hook 不是入口

`make.sh` 被 `.sembumprc.yml` 引用为 `before_bump` hook。当前内容几乎为空
（注释掉的 `go test ./...`），不要把它当成项目入口脚本——开发者实际入口
是 `make ci`。

### 9.8 generated 文件清单（不要手改）

- `gen_*_optiongen.go` — optiongen 生成
- `gen_*.go`（在 xcontainer 系列）— gotemplate 生成
- `*_string.go` — stringer 生成
- `*_mock_test.go` — mockgen 生成

修这些文件等于和 `make gen` 对赌——下次重跑 `make gen` 就被覆盖。改源
注解（`//go:generate ...` 标注的 `option.go` / `type.go` / `gen.go` 等）后
跑 `make gen` 重新生成。

---

_最后更新：本文件随仓库演进；提交影响到约束 / 工作流时同步更新此文件。_

_关联文档_：

- 个人全局规约 `~/.config/opencode/AGENTS.md`（精度优先 / commit·push 约束 /
  PoC + bench / 考古词黑名单等通用规则，本文兜底引用）
- 姊妹仓 `myproxy/AGENTS.md` / `dataserver/AGENTS.md`（风格参考；boost 和它们
  视角不同，约束差异见本文 §1）
