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

// StringsAddPrefix 每一个元素添加前缀
func StringsAddPrefix(s []string, prefix string) []string {
	var out []string
	for _, v := range s {
		out = append(out, prefix+v)
	}
	return out
}

// StringsAddPrefix 每一个元素添加后缀
func StringsAddSuffix(s []string, suffix string) []string {
	var out []string
	for _, v := range s {
		out = append(out, v+suffix)
	}
	return out
}

// StringsRemoveRepeated 移除重复元素
func StringsRemoveRepeated(slc []string) []string {
	if len(slc) == 0 {
		return slc
	}
	if len(slc) < 1024 {
		// 切片长度小于1024的时候，循环来过滤
		return removeRepeatByLoop(slc)
	} else {
		// 大于的时候，通过map来过滤
		return removeRepeatByMap(slc)
	}
}

// StringsRemoveEmpty 移除空元素
func StringsRemoveEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func removeRepeatByMap(slc []string) []string {
	if len(slc) == 0 {
		return slc
	}
	result := []string{}
	tempMap := map[string]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}
func removeRepeatByLoop(slc []string) []string {
	if len(slc) == 0 {
		return slc
	}
	result := []string{} // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}
