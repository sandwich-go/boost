package xcmd

import (
	"github.com/sandwich-go/boost/xstrings"
)

const defaultSliceSeparator = ","

// Slice 将 defaultSliceSeparator 分割的字符获取slice
func Slice(v string) []string { return xstrings.SplitAndTrim(v, defaultSliceSeparator) }
