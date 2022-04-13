package xgen

import (
	"regexp"
	"strings"

	"github.com/sandwich-go/boost/xslice"
)

func RemoveCAndCppComments(src []byte) []byte {
	ccmt := regexp.MustCompile(`\/\*[\s\S]*?\*\/|([^:]|^)\/\/.*$/`)
	out := ccmt.ReplaceAll(src, []byte(""))
	return []byte(strings.Join(xslice.StringsRemoveEmpty(strings.Split(string(out), "\n")), "\n"))
}
