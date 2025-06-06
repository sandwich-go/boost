package mask

import (
	"strings"
)

// 判断是否为邮箱地址
func isEmail(s string) bool {
	atIndex := strings.Index(s, "@")
	if atIndex == -1 {
		return false
	}
	domainPart := s[atIndex+1:]
	return len(domainPart) > 0 && strings.Contains(domainPart, ".")
}

func maskString(s string, opts *Options) string {
	if len(s) <= opts.HideLenMin+opts.PrefixKeep+opts.SuffixKeep {
		return strings.Repeat(string(opts.HideReplaceWith), len(s)) + opts.Suffix
	}
	replaceLen := opts.HideReplaceLen
	if replaceLen == 0 {
		replaceLen = len(s) - opts.PrefixKeep - opts.SuffixKeep
	}
	return s[:opts.PrefixKeep] +
		strings.Repeat(string(opts.HideReplaceWith), replaceLen) +
		s[len(s)-opts.SuffixKeep:] + opts.Suffix
}

// Do 主逻辑函数
func Do(s string, options ...Option) string {
	opts := NewOptions(options...)

	if isEmail(s) {
		// 如果是邮箱，仅对本地部分（@ 之前）进行遮掩
		atIndex := strings.Index(s, "@")
		localPart := s[:atIndex]
		domainPart := s[atIndex:]
		maskedLocal := maskString(localPart, opts)
		return maskedLocal + domainPart
	}
	// 非邮箱直接遮掩整个字符串
	return maskString(s, opts)
}
