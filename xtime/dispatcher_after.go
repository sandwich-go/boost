package xtime

import "time"

// AfterFunc 是对 std time.AfterFunc 的变量级别绑定，便于测试 stub。
// 业务侧调用方应当用 xtime.AfterFunc(...) 而非直接 time.AfterFunc(...)，
// 让全仓有统一的时间抽象点。
var AfterFunc = time.AfterFunc
