package xmath

import "strconv"

// Integer limit values.
const (
	ConstMaxInt8   = 1<<7 - 1
	ConstMinInt8   = -1 << 7
	ConstMaxInt16  = 1<<15 - 1
	ConstMinInt16  = -1 << 15
	ConstMaxInt32  = 1<<31 - 1
	ConstMinInt32  = -1 << 31
	ConstMaxInt64  = 1<<63 - 1
	ConstMinInt64  = -1 << 63
	ConstMaxUint8  = 1<<8 - 1
	ConstMaxUint16 = 1<<16 - 1
	ConstMaxUint32 = 1<<32 - 1
	ConstMaxUint64 = 1<<64 - 1
)

// ParseUint64 parses s as an integer in decimal or hexadecimal syntax.
// Leading zeros are accepted. The empty string parses as zero.
func ParseUint64(s string) (uint64, bool) {
	if s == "" {
		return 0, true
	}
	if len(s) >= 2 && (s[:2] == "0x" || s[:2] == "0X") {
		v, err := strconv.ParseUint(s[2:], 16, 64)
		return v, err == nil
	}
	v, err := strconv.ParseUint(s, 10, 64)
	return v, err == nil
}

// MustParseUint64 parses s as an integer and panics if the string is invalid.
func MustParseUint64(s string) uint64 {
	v, ok := ParseUint64(s)
	if !ok {
		panic("invalid unsigned 64 bit integer: " + s)
	}
	return v
}
