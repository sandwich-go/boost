package xpanic

import "fmt"

// WhenErrorAsFmtFirst err不为nil则wrap并panic，将err作为第一个fmt的参数
// xpanic.WhenErrorAsFmtFirst(err,"got error:%w while reading file:%s",filePath)
func WhenErrorAsFmtFirst(err error, fmtStr string, args ...interface{}) {
	if err == nil {
		return
	}
	var argList []interface{}
	argList = append(argList, err)
	argList = append(argList, args...)
	panic(fmt.Errorf(fmtStr, argList...))
}

// WhenError err不为nil则panic
func WhenError(err error) {
	if err == nil {
		return
	}
	panic(err)
}

// WhenTrue 当condation为true时panic
func WhenTrue(condation bool, fmtStr string, args ...interface{}) {
	if !condation {
		return
	}
	panic(fmt.Errorf(fmtStr, args...))
}
