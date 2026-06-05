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
