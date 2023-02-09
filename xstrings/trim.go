package xstrings

import "strings"

// DefaultTrimChars are the characters which are stripped by Trim* functions in default.
var DefaultTrimChars = string([]byte{
	'\t', // Tab.
	'\v', // Vertical tab.
	'\n', // New line (line feed).
	'\r', // Carriage return.
	'\f', // New page.
	' ',  // Ordinary space.
	0x00, // NUL-byte.
	0x85, // Delete.
	0xA0, // Non-breaking space.
})

// Trim 扩展strings.Trim功能，提供默认的mask字符
func Trim(str string, characterMask ...string) string {
	trimChars := DefaultTrimChars
	if len(characterMask) > 0 {
		trimChars += characterMask[0]
	}
	return strings.Trim(str, trimChars)
}

// SplitAndTrim 将给定字符串分割并trim每一个子元素
func SplitAndTrim(str, delimiter string, characterMask ...string) []string {
	ss := strings.Split(str, delimiter)
	array := make([]string, 0, len(ss))
	for _, v := range ss {
		if len(characterMask) > 0 {
			v = Trim(v, characterMask...)
		}
		if len(v) > 0 {
			array = append(array, v)
		}
	}
	return array
}
