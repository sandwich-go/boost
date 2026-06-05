package xpanic

import (
	"fmt"
	"strings"

	"github.com/sandwich-go/boost/isnil"
)

// WhenErrorAsFmtFirst err 不为 nil 则 wrap 并 panic，将 err 作为第一个 fmt 的参数
// xpanic.WhenErrorAsFmtFirst(err, "got error: %w while reading file: %s", filePath)
func WhenErrorAsFmtFirst(err error, fmtStr string, args ...interface{}) {
	if err == nil {
		return
	}
	var argList = make([]interface{}, 0, len(args)+1)
	argList = append(argList, err)
	argList = append(argList, args...)
	panic(fmt.Sprintf(fmtStr, argList...))
}

// WhenError err 不为 nil 则 panic
func WhenError(err error, reason ...string) {
	if err == nil {
		return
	}
	if len(reason) > 0 {
		panic(strings.Join(reason, "\n"))
	} else {
		panic(err)
	}
}

// WhenTrue 当 condition 为 true 时 panic
func WhenTrue(condition bool, fmtStr string, args ...interface{}) {
	if !condition {
		return
	}
	panic(fmt.Sprintf(fmtStr, args...))
}

// WhenFalse 当 condition 为 false 时 panic
func WhenFalse(condition bool, fmtStr string, args ...interface{}) {
	WhenTrue(!condition, fmtStr, args...)
}

// WhenHereNotNil 提供运行到此处返回的error应为nil的语义，避免在框架层吃掉error
// 功能逻辑等同WhenError，但是语义上调用者确定这里不会返回错误
//
// panic value 是 fmt.Sprintf 格式化后的 string；用 %v 而非 %w —— Sprintf 不支持
// %w，且这里不需要 errors.Is/As 透视（panic 用 string value 让 recover 处理更简单）。
// 历史教训：commit f7dd56a 把 fmt.Errorf 改为 fmt.Sprintf 时漏改 %w → %v，
// vet 一直报 build error 让 xpanic 包的所有测试无法编译，直到本次修复才暴露。
func WhenHereNotNil(err error) {
	if err == nil {
		return
	}
	panic(fmt.Sprintf("err should be nil when here, got:%v", err))
}

// WhenNil 如果v为nil则panic
func WhenNil(v any, fmtStr string, args ...interface{}) {
	WhenTrue(isnil.Check(v), fmtStr, args...)
}

// WhenNotNil 如果v不为nil则panic
func WhenNotNil(v any, fmtStr string, args ...interface{}) {
	WhenTrue(!isnil.Check(v), fmtStr, args...)
}
