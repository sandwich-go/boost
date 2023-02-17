package sdebug

import (
	"runtime"
	"strings"
)

const (
	maxCallerDepth            = 1000
	StackFilterKeyForSandwich = "/funplus/sandwich/"
	stackFilterKey            = "/debug/sdebug/"
)

var (
	goRootForFilter = runtime.GOROOT()
)

func init() {
	if goRootForFilter != "" {
		goRootForFilter = strings.Replace(goRootForFilter, "\\", "/", -1)
	}
}

func shouldContinue(file string) bool {
	if Enabled() {
		if strings.Contains(file, stackFilterKey) {
			return true
		}
	} else {
		if strings.Contains(file, StackFilterKeyForSandwich) {
			return true
		}
	}
	return false
}

func callerFromIndex(filters []string) (pc uintptr, file string, line, index int) {
	var filtered, ok bool
	for index = 0; index < maxCallerDepth; index++ {
		if pc, file, line, ok = runtime.Caller(index); !ok {
			continue
		}
		filtered = false
		for _, filter := range filters {
			if filter != "" && strings.Contains(file, filter) {
				filtered = true
				break
			}
		}
		if filtered || shouldContinue(file) {
			continue
		}

		if index > 0 {
			index--
		}
		return
	}
	return 0, "", -1, -1
}
