package xpanic

import "fmt"

// PanicIfErrorAsFmtFirst err不为nil则wrap并panic，将err作为第一个fmt的参数
func PanicIfErrorAsFmtFirst(err error, fmtStr string, args ...interface{}) {
	if err == nil {
		return
	}
	var argList []interface{}
	argList = append(argList, err)
	argList = append(argList, args...)
	panic(fmt.Errorf(fmtStr, argList...))
}

// PanicIfTrue 当condation为true时panic
func PanicIfTrue(condation bool, fmtStr string, args ...interface{}) {
	if !condation {
		return
	}
	panic(fmt.Errorf(fmtStr, args...))
}
