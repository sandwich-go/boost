package xslice

import "strings"

// ContainString 返回slice中是否含有指定字符串
func ContainString(slice []string, v string) bool {
	for _, vv := range slice {
		if vv == v {
			return true
		}
	}
	return false
}

// ContainStringEqualFold 返回slice中是否含有指定字符串，不区分大小写
func ContainStringEqualFold(slice []string, v string) bool {
	for _, vv := range slice {
		if strings.EqualFold(vv, v) {
			return true
		}
	}
	return false
}

// StringSetAdd 如果slice中不存在给定的元素v则添加
func StringSetAdd(slice []string, v ...string) []string {
	for _, vv := range v {
		if !ContainString(slice, vv) {
			slice = append(slice, vv)
		}
	}
	return slice
}

// StringSliceWalk 遍历vs,将f应用到每一个元素，返回更新后的数据
func StringSliceWalk(vs []string, f func(string) (string, bool)) []string {
	vsm := make([]string, 0)
	for _, v := range vs {
		ret, valid := f(v)
		if valid {
			vsm = append(vsm, ret)
		}
	}
	return vsm
}
