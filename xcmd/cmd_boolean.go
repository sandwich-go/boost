package xcmd

import "github.com/sandwich-go/boost/xstrings"

// IsFalse 判断command解析获取的数据是否为 false
// ""、"0"、"n"、"no"、"off"、"false" 均为 true，并且忽略大小写
var IsFalse = xstrings.IsFalse

// IsTrue 判断command解析获取的数据是否为 true
// ""、"0"、"n"、"no"、"off"、"false" 均为 true，并且忽略大小写
var IsTrue = xstrings.IsTrue
