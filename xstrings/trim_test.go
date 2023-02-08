package xstrings

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTrim(t *testing.T) {
	Convey("trim", t, func() {
		for _, test := range []struct {
			s        string
			expected string
			mask     string
		}{
			{s: "a\t", expected: "a"},
			{s: "a\v", expected: "a"},
			{s: "a\n", expected: "a"},
			{s: "a\r", expected: "a"},
			{s: "a\f\f", expected: "a"},
			{s: "a  ", expected: "a"},
			{s: string([]byte{'a', 0x00}), expected: "a"},
			{s: string([]byte{'a', 0x85}), expected: "a"},
			{s: string([]byte{'a', 0xA0}), expected: "a"},
			{s: "ab", mask: "b", expected: "a"},
		} {
			So(Trim(test.s, test.mask), ShouldEqual, test.expected)
		}

		for _, test := range []struct {
			s         string
			expected  []string
			delimiter string
			mask      string
		}{
			{s: "a\t,b\t", expected: []string{"a", "b"}, delimiter: ","},
			{s: "a\v,b\v", expected: []string{"a", "b"}, delimiter: ","},
			{s: "a\n,a\n", expected: []string{"a", "a"}, delimiter: ","},
			{s: "a\r,a\n", expected: []string{"a", "a"}, delimiter: ","},
			{s: "a\f\f", expected: []string{"a"}, delimiter: ","},
			{s: "a  ,b  ", expected: []string{"a", "b"}, delimiter: ","},
			{s: string([]byte{'a', 0x00, ',', 'b', 0x00}), expected: []string{"a", "b"}, delimiter: ","},
			{s: string([]byte{'a', 0x85, ',', 'b', 0x00}), expected: []string{"a", "b"}, delimiter: ","},
			{s: string([]byte{'a', 0xA0, ',', 'b', 0x85}), expected: []string{"a", "b"}, delimiter: ","},
			{s: "a,b", mask: "b", expected: []string{"a"}, delimiter: ","},
		} {
			So(SplitAndTrim(test.s, test.delimiter, test.mask), ShouldResemble, test.expected)
		}
	})
}
