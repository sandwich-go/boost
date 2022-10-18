package xerror

import "fmt"

// Errorf returns an error that formats as the given format and args.
var Errorf = NewText

// NewText returns an error that formats as the given format and args.
func NewText(format string, args ...interface{}) *Error {
	text := format
	if len(args) != 0 {
		text = fmt.Sprintf(format, args...)
	}
	return &Error{
		callStack: callersCheckIsErrorWithStack(1),
		text:      text,
		code:      ErrorCodeUnsetAsDefault,
	}
}

// proto-gen-go v2
// https://github.com/golang/protobuf/blob/5b6d0471e5bd463898a8121123ace60c7a0d9ca7/reflect/protoreflect/value.go#L9-L22
type protoEnum interface {
	fmt.Stringer
	NumberInt32() int32
	EnumDescriptor() ([]byte, []int)
}

// NewCode 根据错误码和指定的format, args(可为空)构造错误信息
func NewCode(code int32, format string, args ...interface{}) *Error {
	text := format
	if len(args) != 0 {
		text = fmt.Sprintf(format, args...)
	}
	return &Error{
		callStack: callersCheckIsErrorWithStack(1),
		text:      text,
		code:      code,
	}
}

// NewProtoEnum 根据proto enum错误码和指定的format, args(可为空)构造错误信息
func NewProtoEnum(code protoEnum) *Error {
	return &Error{
		callStack: callersCheckIsErrorWithStack(1),
		text:      code.String(),
		code:      code.NumberInt32(),
	}
}

// Wrap 封装底层错误，提供另一个错误信息,如果error为nil则返回nil
func Wrap(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	text := format
	if len(args) != 0 {
		text = fmt.Sprintf(format, args...)
	}
	return &Error{
		err:       err,
		callStack: callersCheckIsErrorWithStack(1),
		text:      text,
		code:      Code(err),
	}
}

// WrapCode 封装底层错误，指定错误码,如果error为nil则返回nil
func WrapCode(code int32, err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	text := format
	if len(args) != 0 {
		text = fmt.Sprintf(format, args...)
	}
	return &Error{
		err:       err,
		callStack: callersCheckIsErrorWithStack(1),
		text:      text,
		code:      code,
	}
}

// WrapProtoEnum 封装底层错误，指定proto Enum错误码,如果error为nil则返回nil
func WrapProtoEnum(code protoEnum, err error) error {
	if err == nil {
		return nil
	}
	return &Error{
		err:       err,
		callStack: callersCheckIsErrorWithStack(1),
		text:      code.String(),
		code:      code.NumberInt32(),
	}
}
