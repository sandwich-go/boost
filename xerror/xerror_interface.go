package xerror

// APICode Code feature.
type APICode interface {
	error
	Code() int32
}

// apiStack Stack feature.
type apiStack interface {
	error
	Stack() string
}

// apiStack Stack feature.
type apiCaller interface {
	error
	Caller(skip int) (file, funcName string, line int)
}

// apiCause Cause feature.
type apiCause interface {
	Error() string
	Cause() error
}

// apiLogic Logic Exception feature.
type apiLogic interface {
	Error() string
	Logic() bool
}

// Code 返回错误码数据，如果没有实现APICode则根据CodeHandlerForNotAPICode逻辑返回
func Code(err error) int32 {
	if err != nil {
		if e, ok := err.(APICode); ok {
			return e.Code()
		}
	}
	return CodeHandlerForNotAPICode(err)
}

func Caller(err error, skip int) (file, funcName string, line int) {
	if err == nil {
		return
	}
	if e, ok := err.(apiCaller); ok {
		return e.Caller(skip)
	}
	return
}

// Logic 返回是否是Logic层异常，默认为false
func Logic(err error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(apiLogic); ok {
		return e.Logic()
	}
	// 兼容err2接口
	if e, ok := err.(interface{ IsLogicException() bool }); ok {
		return e.IsLogicException()
	}
	return false
}

// Cause 返回最底层的错误信息，如果没有实现apiCause，则返回当前错误信息
func Cause(err error) error {
	if err != nil {
		if e, ok := err.(apiCause); ok {
			return e.Cause()
		}
	}
	return err
}

// Stack 返回堆栈信息，如果err不支持apiStack，则返回错误本身数据
func Stack(err error) string {
	if err == nil {
		return ""
	}
	if e, ok := err.(apiStack); ok {
		return e.Stack()
	}
	return err.Error()
}
