package xpanic

import "fmt"

// WhenErrorAsFmtFirst err 不为 nil 则 wrap 并 panic，将 err 作为第一个 fmt 的参数
// xpanic.WhenErrorAsFmtFirst(err, "got error: %w while reading file: %s", filePath)
func WhenErrorAsFmtFirst(err error, fmtStr string, args ...interface{}) {
	if err == nil {
		return
	}
	var argList = make([]interface{}, 0, len(args)+1)
	argList = append(argList, err)
	argList = append(argList, args...)
	panic(fmt.Errorf(fmtStr, argList...))
}

// WhenError err 不为 nil 则 panic
func WhenError(err error) {
	if err == nil {
		return
	}
	panic(err)
}

// WhenTrue 当 condition 为 true 时 panic
func WhenTrue(condition bool, fmtStr string, args ...interface{}) {
	if !condition {
		return
	}
	panic(fmt.Errorf(fmtStr, args...))
}

// WhenHereShouldNotError 提供运行到此处返回的error应为nil的语义，避免在框架层吃掉error
func WhenHereShouldNotError(err error) {
	WhenError(err)
}
