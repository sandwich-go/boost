package hide

import (
	"strings"
)

// Do 隐藏字符串
func Do(s string, options ...Option) string {
	opts := NewOptions(options...)
	if len(s) <= opts.HideLenMin+opts.PrefixKeep+opts.SuffixKeep {
		return strings.Repeat("*", len(s)) + "@" + opts.Suffix
	}
	replaceLen := opts.HideReplaceLen
	if replaceLen == 0 {
		replaceLen = len(s) - opts.PrefixKeep - opts.SuffixKeep
	}
	return s[:opts.PrefixKeep] + strings.Repeat(string(opts.HideReplaceWith), replaceLen) + s[len(s)-opts.SuffixKeep:] + "@" + opts.Suffix
}
