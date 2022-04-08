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

func Trim(str string, characterMask ...string) string {
	trimChars := DefaultTrimChars
	if len(characterMask) > 0 {
		trimChars += characterMask[0]
	}
	return strings.Trim(str, trimChars)
}
