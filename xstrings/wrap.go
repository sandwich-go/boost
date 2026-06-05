package xstrings

import (
	"bytes"
	"unicode"
)

const nbsp = 0xA0

// Wrap 换行给定字符串
func Wrap(s string, lim uint) string {
	// Initialize a buffer with a slightly larger size to account for breaks
	init := make([]byte, 0, len(s))
	buf := bytes.NewBuffer(init)

	var current uint
	var wordBuf, spaceBuf bytes.Buffer
	var wordBufLen, spaceBufLen uint

	for _, char := range s {
		if char == '\n' {
			// 注意 char=='\n' 分支末尾 current = 0 重置；上面 if/else 中对
			// current 的赋值是 ineffective（被 line 末尾覆盖），仅保留有副
			// 作用的 spaceBuf.WriteTo / Reset 操作。
			if wordBuf.Len() == 0 {
				if current+spaceBufLen <= lim {
					_, _ = spaceBuf.WriteTo(buf)
				}
				spaceBuf.Reset()
				spaceBufLen = 0
			} else {
				_, _ = spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				spaceBufLen = 0
				_, _ = wordBuf.WriteTo(buf)
				wordBuf.Reset()
				wordBufLen = 0
			}
			buf.WriteRune(char)
			current = 0
		} else if unicode.IsSpace(char) && char != nbsp {
			if spaceBuf.Len() == 0 || wordBuf.Len() > 0 {
				current += spaceBufLen + wordBufLen
				_, _ = spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				spaceBufLen = 0
				_, _ = wordBuf.WriteTo(buf)
				wordBuf.Reset()
				wordBufLen = 0
			}

			spaceBuf.WriteRune(char)
			spaceBufLen++
		} else {
			wordBuf.WriteRune(char)
			wordBufLen++

			if current+wordBufLen+spaceBufLen > lim && wordBufLen < lim {
				buf.WriteRune('\n')
				current = 0
				spaceBuf.Reset()
				spaceBufLen = 0
			}
		}
	}

	if wordBuf.Len() == 0 {
		if current+spaceBufLen <= lim {
			_, _ = spaceBuf.WriteTo(buf)
		}
	} else {
		_, _ = spaceBuf.WriteTo(buf)
		_, _ = wordBuf.WriteTo(buf)
	}

	return buf.String()
}
