package xpanic

import (
	"testing"
)

func Test_NormalFlow(T *testing.T) {
	called := false

	Try(func() {
		called = true

	}).Catch(func(_ E) {
		T.Error("Catch must not be called")
	})

	// if try was not called
	if !called {
		T.Error("Try do not called")
	}
}

func Test_NormalFlowFinally(T *testing.T) {
	calledTry := false
	calledFinally := false

	Try(func() {
		calledTry = true

	}).Finally(func() {
		calledFinally = true

	}).Catch(func(_ E) {
		T.Error("Catch must not be called")

	})

	// if try was not called
	if !calledTry {
		T.Error("Try do not called")
	}

	// if finally was not called
	if !calledFinally {
		T.Error("Finally do not called")
	}
}

// catch handler 收到的 e 是 exception 结构体（commit d885113 设计选择），
// 不是 panic 抛出的原值。需要 e.(exception).Error 取 panic value。同包测试
// 直接用未导出的 exception 类型；外部包目前没有便捷断言路径（潜在 API
// 设计问题，留给未来 PR）。
func panicValueFrom(e E) any {
	if exc, ok := e.(exception); ok {
		return exc.Error
	}
	return e
}

func Test_CrashInTry(T *testing.T) {
	calledFinally := false
	calledCatch := false

	Try(func() {
		panic("testing panic")

	}).Finally(func() {
		calledFinally = true

	}).Catch(func(e E) {
		calledCatch = true
		if panicValueFrom(e) != "testing panic" {
			T.Error("error is not 'testing panic'")
		}
	})

	// if catch was not called
	if !calledCatch {
		T.Error("Catch do not called")
	}

	// if finally was not called
	if !calledFinally {
		T.Error("Finally do not called")
	}
}

func Test_CrashInTry2(T *testing.T) {
	calledCatch := false

	Try(func() {
		panic("testing panic")

	}).Catch(func(e E) {
		calledCatch = true
		if panicValueFrom(e) != "testing panic" {
			T.Error("error is not 'testing panic'")
		}
	})

	// if catch was not called
	if !calledCatch {
		T.Error("Catch do not called")
	}
}

func Test_CrashInCatch(T *testing.T) {
	calledFinally := false

	defer func() {
		err := recover()
		if err != "another panic" {
			T.Error("error is not 'another panic'")
		}
		// if finally was not called
		if !calledFinally {
			T.Error("Finally do not called")
		}
	}()
	Try(func() {
		panic("testing panic")

	}).Finally(func() {
		calledFinally = true

	}).Catch(func(e E) {
		if panicValueFrom(e) != "testing panic" {
			T.Error("error is not 'testing panic'")
		}
		panic("another panic")

	})
}

func Test_CrashInCatch2(T *testing.T) {
	defer func() {
		err := recover()
		if err != "another panic" {
			T.Error("error is not 'another panic'")
		}
	}()
	Try(func() {
		panic("testing panic")

	}).Catch(func(e E) {
		if panicValueFrom(e) != "testing panic" {
			T.Error("error is not 'testing panic'")
		}
		panic("another panic")
	})
}

func Test_CrashInThrow(T *testing.T) {
	calledFinally := false

	defer func() {
		err := recover()
		if err != "testing panic" {
			T.Error("error is not 'testing panic'")
		}
		// if finally was not called
		if !calledFinally {
			T.Error("Finally do not called")
		}
	}()

	Try(func() {
		panic("testing panic")

	}).Finally(func() {
		calledFinally = true

	}).Catch(func(e E) {
		if panicValueFrom(e) != "testing panic" {
			T.Error("error is not 'testing panic'")
		}
		Throw()
	})
}

func Test_CrashInThrow2(T *testing.T) {
	defer func() {
		err := recover()
		if err != "testing panic" {
			T.Error("error is not 'testing panic'")
		}
	}()

	Try(func() {
		panic("testing panic")

	}).Catch(func(e E) {
		if panicValueFrom(e) != "testing panic" {
			T.Error("error is not 'testing panic'")
		}
		Throw()
	})
}

func Test_CrashInFinally1(T *testing.T) {
	calledTry := false

	defer func() {
		err := recover()
		if err != "finally panic" {
			T.Error("error is not 'finally panic'")
		}

		// if try was not called
		if !calledTry {
			T.Error("Try do not called")
		}
	}()

	Try(func() {
		calledTry = true

	}).Finally(func() {
		panic("finally panic")

	}).Catch(func(e E) {
		T.Error("Catch must not be called")
	})
}

func Test_CrashInFinally2(T *testing.T) {

	defer func() {
		err := recover()
		if err != "finally panic" {
			T.Error("error is not 'finally panic'")
		}
	}()

	Try(func() {
		panic("testing panic")

	}).Finally(func() {
		panic("finally panic")

	}).Catch(func(e E) {
		if panicValueFrom(e) != "testing panic" {
			T.Error("error is not 'testing panic'")
		}
		panic("another panic")

	})
}

// Test_ExceptionString 验证 exception.String() 输出 "<panic value>\n<stack>"
// 格式（commit d885113 的 catch handler 拿到 exception 结构体后用 %v 即触发
// 此 Stringer，仓内 xparallel/xtime/xcompress 三处真实调用方都依赖该格式）。
func Test_ExceptionString(T *testing.T) {
	var snapshot string
	Try(func() {
		panic("boom")
	}).Catch(func(e E) {
		exc := e.(exception)
		snapshot = exc.String()
	})
	if snapshot == "" {
		T.Fatal("exception.String returned empty")
	}
	// 第一行应该是 panic value
	if !startsWith(snapshot, "boom\n") {
		T.Errorf("exception.String first line not panic value: %q", firstLine(snapshot))
	}
	// 后续应该含栈帧（debug.Stack 输出）
	if !contains(snapshot, "goroutine ") {
		T.Errorf("exception.String missing stack frames: %q", snapshot)
	}
}

// Test_FinallyOnly_PanicInFinally 验证只有 Finally 的路径（无 Catch）panic 透传。
// 覆盖原 try.go Finally 分支里 e.finally != nil 的纯调用路径。
func Test_FinallyOnly_PanicInFinally(T *testing.T) {
	calledFinally := false
	defer func() {
		if r := recover(); r != "finally only panic" {
			T.Errorf("unexpected recover: %v", r)
		}
		if !calledFinally {
			T.Error("Finally not called")
		}
	}()
	Try(func() {
		// 没有 panic，Try 完成
	}).Finally(func() {
		calledFinally = true
		panic("finally only panic")
	}).Catch(func(_ E) {
		T.Error("Catch must not be called")
	})
}

// helper
func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
func firstLine(s string) string {
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			return s[:i]
		}
	}
	return s
}

// Test_FinallyTwice 验证连续两次 Finally panic 防御断言（同一 exception 不
// 允许多次注册 Finally；这是 try.go API 契约）。
func Test_FinallyTwice(T *testing.T) {
	defer func() {
		if r := recover(); r != "finally was only set" {
			T.Errorf("expected 'finally was only set' panic, got %v", r)
		}
	}()
	Try(func() {}).Finally(func() {}).Finally(func() {})
}
