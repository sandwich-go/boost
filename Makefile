export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on

# boost 是通用工具库，47 个顶级子包；上游业务（myproxy / dataserver）通过
# minor 版本号 pin。Makefile 设计目标：
#   - make ci：本地准入，等价 .github/workflows/ci.yml 的硬门槛
#   - make bench_diff：动 hot path 时跑（xrand / xpool / xcontainer / xencoding），
#     与入仓基线 benchdata/main.txt 对比
#   - make gen：仓内有 optiongen / gotemplate / stringer / mockgen 生成代码
#
# 不在本 Makefile 里：fuzz_smoke / fuzz_long / docker e2e / chaos
# 原因：boost 是库不是服务进程，没有 docker 部署；仓内目前没有 Fuzz* 函数
# （rg '^func Fuzz' --glob '*_test.go' 为空）。需要时再补，避免占位空骨架。

TEST_FILES?=$$(go list ./... | grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor | grep -v '/gen/')
COVER_PROFILE?=coverage.out

default: ci

# ---------------------------------------------------------------------------
# tools：一次性安装代码生成 / lint / vuln / benchstat 工具
# ---------------------------------------------------------------------------

tools: tools_cover tools_gen tools_lint tools_vuln tools_benchstat

tools_cover:
	@go tool cover 2>/dev/null >/dev/null; if [ $$? -eq 3 ]; then \
		go install golang.org/x/tools/cmd/cover@latest; \
	fi

# tools_gen：boost 仓内 //go:generate 指令实际用到的 4 个工具
#   optiongen   gen_*_optiongen.go（约 20 个包）
#   gotemplate  xtime/wheel_map.go / xcontainer 系列
#   stringer    httputil/dns/policy_string.go / xencoding/{encrypt,compressor}/type.go
#   mockgen     geo/grid_mock_test.go
#
# stringer 版本陷阱（与 myproxy 一致）：2024-03 之前的版本在 Go 1.25 toolchain
# 下会报 "package context without types was imported"（types 包 ABI 变更）。
# 若本机 stringer 长期没更新，make gen 会挂在 *_string.go 这一步。
# 解决：rm "$(go env GOPATH)/bin/stringer" && make tools_gen
# 这里不做强制 go install @latest，避免每次 make gen 都走网络。
tools_gen:
	@command -v optiongen   >/dev/null 2>&1 || go install github.com/timestee/optiongen/cmd/optiongen@latest
	@command -v gotemplate  >/dev/null 2>&1 || go install github.com/ncw/gotemplate@latest
	@command -v stringer    >/dev/null 2>&1 || go install golang.org/x/tools/cmd/stringer@latest
	@command -v mockgen     >/dev/null 2>&1 || go install go.uber.org/mock/mockgen@latest

GOLANGCI_LINT_VERSION ?= v2.12.2
GOLANGCI_LINT         := $(shell go env GOPATH)/bin/golangci-lint

# tools_lint：钉版本安装，原因与 myproxy 一致：
#   - 旧版 v1.x 与 v2 schema 不兼容
#   - go install 编译需要本地 go 版本 >= lint 源码要求；用预编译 binary 解耦
#
# 注意 boost 仓 §4.5：当前 lint 未钉为硬约束，CI 软门槛跑。本 target 仍按
# 钉版安装，让 `make lint` 本地行为可预期。
tools_lint:
	@gopath="$$(go env GOPATH)"; \
	if [ -z "$$gopath" ]; then echo "tools_lint: go env GOPATH 为空，先装 Go 或修 GOPATH"; exit 1; fi; \
	if $(GOLANGCI_LINT) version 2>/dev/null | grep -q "$(GOLANGCI_LINT_VERSION:v%=%)"; then \
		exit 0; \
	fi; \
	echo "==> installing golangci-lint $(GOLANGCI_LINT_VERSION) → $$gopath/bin"; \
	tmpdir="$$(mktemp -d)"; \
	trap 'rm -rf "$$tmpdir"' EXIT; \
	if ! curl -fsSL --max-time 60 -o "$$tmpdir/install.sh" https://raw.githubusercontent.com/golangci/golangci-lint/main/install.sh; then \
		echo "tools_lint: 拉 install.sh 失败（网络 / DNS / 防火墙？）"; \
		exit 1; \
	fi; \
	if command -v timeout >/dev/null 2>&1; then \
		timeout 180 sh "$$tmpdir/install.sh" -b "$$gopath/bin" $(GOLANGCI_LINT_VERSION); \
	else \
		sh "$$tmpdir/install.sh" -b "$$gopath/bin" $(GOLANGCI_LINT_VERSION); \
	fi; \
	rc=$$?; \
	if [ $$rc -ne 0 ]; then \
		echo "tools_lint: install.sh 退出码 $$rc"; \
		exit $$rc; \
	fi

GOVULNCHECK := $(shell go env GOPATH)/bin/govulncheck
BENCHSTAT   := $(shell go env GOPATH)/bin/benchstat

tools_vuln:
	@command -v govulncheck >/dev/null 2>&1 || go install golang.org/x/vuln/cmd/govulncheck@latest

tools_benchstat:
	@command -v benchstat >/dev/null 2>&1 || go install golang.org/x/perf/cmd/benchstat@latest

# ---------------------------------------------------------------------------
# format / static check
# ---------------------------------------------------------------------------

fmt:
	gofmt -w $(GOFMT_FILES)

# goimports 的 -local 设成 boost 自己的 module path，让 boost 内部 import
# 排在第三方包之后第一组。
imports import goimports:
	goimports -local=github.com/sandwich-go/boost -w $(GOFMT_FILES)

vet:
	go vet ./...

# lint：当前未钉为硬约束（详见 AGENTS.md §4.5）。本 target 仍提供，让本地
# 自查 / 后续接入硬门槛时不需要改 Makefile。.golangci.yml 暂未入仓：lint
# 会用默认规则集，已知会出 200-300 个告警（boost 是老仓），先以"本地参考"
# 看告警密度趋势。
lint: tools_lint
	$(GOLANGCI_LINT) run ./...

# vuln：govulncheck call-graph 驱动，只报"代码路径真正会触发"的 CVE。
# 升级依赖前 / 主干同步上游 stdlib CVE 时跑。
vuln: tools_vuln
	$(GOVULNCHECK) ./...

# ---------------------------------------------------------------------------
# test / race / cover
# ---------------------------------------------------------------------------

# 主测试目标：覆盖 + 5min 超时（boost 单测应远快于此，超时多半是 deadlock）。
test:
	go test $(TEST_FILES) -cover -timeout=5m -parallel=4 $(TESTARGS)

# test_short：跑带 -short flag 的快速子集，用于 IDE save-on-test。
test_short:
	go test $(TEST_FILES) -short -cover -timeout=2m -parallel=4 $(TESTARGS)

# test_race：race 是并发包硬约束（xtime / xsync / xchan / xcontainer/syncmap /
# xpool / lru / ratelimiter / singleflight）。10min 超时给 race 检测器额外预算
# （-race 会让测试慢 2-10x）。
test_race race:
	go test -race $(TEST_FILES) -timeout=10m $(TESTARGS)

# cover：保留 $(COVER_PROFILE) 给后续 cover_html / sonar 用。
# 过滤掉 generated 文件（gen_*.go / *_string.go / option.go）和 internal/
# template2 等，避免污染覆盖率统计。
cover: tools_cover
	go test $(TEST_FILES) -coverprofile=$(COVER_PROFILE).tmp -timeout=5m
	@cat $(COVER_PROFILE).tmp \
		| grep -v "gen.go" \
		| grep -v "option.go" \
		| grep -v "_test.go" \
		| grep -v "gen_" \
		| grep -v vendor \
		| grep -v test_perf \
		| grep -v internal/template2 \
		> $(COVER_PROFILE)
	@rm -f $(COVER_PROFILE).tmp
	@go tool cover -func=$(COVER_PROFILE) | tail -n 30
	@echo "==> coverage profile written to $(COVER_PROFILE)"

cover_html: cover
	go tool cover -html=$(COVER_PROFILE)

# 简洁文本：每个包一行覆盖率（不写 profile 文件）
cover_pkg:
	@go test $(TEST_FILES) -cover -timeout=5m | grep -E "coverage:|^FAIL|^ok " || true

# ---------------------------------------------------------------------------
# benchmark：关键路径基线 + 与入仓基线对比
# ---------------------------------------------------------------------------
#
# AGENTS.md §5.2 列出已入仓 bench 文件。改 hot path 必须 PoC + bench：
#   1. main 上跑 `make bench` → cp bench.txt benchdata/main.txt（如果还没 baseline）
#   2. 切分支改代码 → `make bench_diff`，看 ±X% + p-value
#   3. 主干 perf 改动 land 后 `make bench_refresh` 单独 chore commit
#
# BENCH_PKGS 列出当前真有 Benchmark* 的包；动新包前要在这里加：
BENCH_PKGS    ?= ./xrand/... ./xpool/... ./xtime/... ./xencoding/... \
                 ./xcrypto/algorithm/aes/... ./xhash/...
BENCH_RUN      ?= -run=^$$ -bench=. -benchmem -benchtime=1s -count=5
BENCH_OUT      ?= bench.txt
BENCH_BASELINE ?= benchdata/main.txt

bench:
	@echo "==> running benchmarks: $(BENCH_PKGS)"
	@go test $(BENCH_PKGS) $(BENCH_RUN) | tee $(BENCH_OUT)
	@echo "==> wrote $(BENCH_OUT)"

# benchstat：用 benchdata/main.txt 入仓基线对比当前 bench.txt。
# 机器抖动会让 ±X% 产生假阳性，看 p-value / geomean 综合判断。
# AGENTS.md §11 提到的"单次 benchstat 不能直接当退化结论"对 boost 同样适用。
benchstat: tools_benchstat
	@if [ ! -f $(BENCH_BASELINE) ]; then \
		echo "缺少 $(BENCH_BASELINE)（入仓基线），请先在主干跑 'make bench' 再 cp bench.txt benchdata/main.txt"; \
		exit 1; \
	fi
	@if [ ! -f $(BENCH_OUT) ]; then \
		echo "缺少 $(BENCH_OUT)，请先 'make bench'"; \
		exit 1; \
	fi
	$(BENCHSTAT) $(BENCH_BASELINE) $(BENCH_OUT)

# bench_diff = bench + benchstat 一键跑：对照入仓基线看改动侧影响
bench_diff: bench benchstat
	@echo "==> bench_diff against $(BENCH_BASELINE) done"

# bench_refresh：把当前 bench.txt 作为新基线写入 benchdata/main.txt。
# 只在主干 perf 真实变动（修复 / 大改动 landed）后跑一次，单独 chore commit。
bench_refresh: bench
	@mkdir -p benchdata
	@cp $(BENCH_OUT) $(BENCH_BASELINE)
	@echo "==> refreshed baseline: $(BENCH_BASELINE)"
	@echo "    请在 chore commit 里独立提交，不与功能改动混合（AGENTS.md §7.2）"

# ---------------------------------------------------------------------------
# code generation：optiongen / gotemplate / stringer / mockgen
# ---------------------------------------------------------------------------
#
# 仓内 //go:generate 指令分布：
#   - 20+ 个 gen_*_optiongen.go（optiongen 从 option.go 等生成 setter/getter/visitor）
#   - xtime/wheel_map.go（gotemplate 从 base/container/templates/syncmap 生成）
#   - xcontainer/syncmap/* 等（gotemplate）
#   - httputil/dns/policy_string.go / xencoding/*/type.go（stringer）
#   - geo/grid_mock_test.go（mockgen）
#
# 改源文件后跑 make gen，diff 必须随同源改动一起提交（AGENTS.md §7.2）。
gen generate: tools_gen
	go generate ./...

# ---------------------------------------------------------------------------
# composite
# ---------------------------------------------------------------------------

# 本地准入：跑测试 + 自动 fmt/imports + go mod tidy
commit: test style

# CI 准入：不修改源码，只校验。等价 .github/workflows/ci.yml 的硬门槛。
# - vet：编译期错误 + 简单静态检查
# - build：所有包编译通过（含 ./...）
# - test_race：并发包 race 检测；race 是仓库硬约束（AGENTS.md §4.4）
# - vuln：govulncheck call-graph 驱动，0 affecting CVE 期望（AGENTS.md §4.1）
#
# 不含 lint：当前未钉为硬约束（AGENTS.md §4.5）。
ci: vet build test_race vuln
	@echo "ci OK"

# 本地 commit 前的自动整理（修改源码）：errcheck 这里没接，靠 lint 兜底
style: fmt imports
	go mod tidy

# build：编译所有包但不产出 binary（boost 是库）。
# 用 ./... 一次性确保所有子包都能 import；avoid silent breakage on某个边
# 缘子包没 import path 错误。
build:
	go build ./...

# version：发版工具（protokitgo sem release）会推下一个 tag 并 push origin。
# 仅在用户明确要求"打 tag / 升版本"时才跑（AGENTS.md §7.3）。
version:
	protokitgo sem release

.NOTPARALLEL:
.PHONY: default tools tools_cover tools_gen tools_lint tools_vuln tools_benchstat \
	fmt imports import goimports vet lint vuln \
	test test_short test_race race \
	cover cover_pkg cover_html \
	bench benchstat bench_diff bench_refresh \
	gen generate build \
	commit ci style version
