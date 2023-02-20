package xgen

import (
	"regexp"
	"strings"

	"github.com/sandwich-go/boost/xslice"
	"github.com/sandwich-go/boost/xstrings"
)

func RemoveCStyleComments(content []byte) []byte {
	// http://blog.ostermiller.org/find-comment
	ccmt := regexp.MustCompile(`\/\*[\s\S]*?\*\/|([^:]|^)\/\/.*$/`)
	return ccmt.ReplaceAll(content, []byte(""))
}

func RemoveCppStyleComments(content []byte) []byte {
	cppcmt := regexp.MustCompile(`//.*`)
	return cppcmt.ReplaceAll(content, []byte(""))
}

func RmoveCAndCppCommentAndBlanklines(src []byte) []byte {
	out := RemoveCppStyleComments(RemoveCStyleComments(src))
	return []byte(strings.Join(xslice.StringsWalk(strings.Split(string(out), "\n"), func(s string) (string, bool) {
		return s, xstrings.Trim(s) != ""
	}), "\n"))
}
