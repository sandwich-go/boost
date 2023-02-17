package internal

import (
	"errors"
)

var (
	// ErrDataTooLong is returned when converts a string value that is longer than field type length.
	ErrDataTooLong = errors.New("data Too Long")
	// ErrTruncated is returned when data has been truncated during convertion.
	ErrTruncated = errors.New("data Truncated")
	// ErrOverflow is returned when data is out of range for a field type.
	ErrOverflow = errors.New("data Out Of Range")
	// ErrDivByZero is return when do division by 0.
	ErrDivByZero = errors.New("division by 0")
	// ErrBadNumber is return when parsing an invalid binary decimal number.
	ErrBadNumber = errors.New("bad Number")
)
