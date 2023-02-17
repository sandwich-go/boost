package sdebug

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

const shoudlAddSpace = 9

// PrintStack fmt.Print打印堆栈信息
func PrintStack(skip ...int) { fmt.Print(Stack(skip...)) }

// Stack 返回调用堆栈
func Stack(skip ...int) string { return StackWithFilter("", skip...) }

// StackWithFilter 根据指定的filter返回调用堆栈数据
func StackWithFilter(filter string, skip ...int) string {
	return stackWithFilters([]string{filter}, skip...)
}

func stackWithFilters(filters []string, skip ...int) string {
	number := 0
	if len(skip) > 0 {
		number = skip[0]
	}
	var name string
	var space string
	var filtered bool
	var index = 1
	var buffer = bytes.NewBuffer(nil)
	var ok = true
	var pc, file, line, start = callerFromIndex(filters)

	for i := start + number; i < maxCallerDepth; i++ {
		if i != start {
			pc, file, line, ok = runtime.Caller(i)
		}
		if !ok {
			break
		}
		// Filter empty file.
		if file == "" {
			continue
		}
		// GOROOT filter.
		if goRootForFilter != "" &&
			len(file) >= len(goRootForFilter) &&
			file[0:len(goRootForFilter)] == goRootForFilter {
			continue
		}
		// Custom filtering.
		filtered = false
		for _, filter := range filters {
			if filter != "" && strings.Contains(file, filter) {
				filtered = true
				break
			}
		}
		if filtered {
			continue
		}
		if fn := runtime.FuncForPC(pc); fn == nil {
			name = "unknown"
		} else {
			name = fn.Name()
		}
		if index > shoudlAddSpace {
			space = " "
		}
		buffer.WriteString(fmt.Sprintf("%d.%s%s\n    %s:%d\n", index, space, name, file, line))
		index++
	}
	return buffer.String()
}
