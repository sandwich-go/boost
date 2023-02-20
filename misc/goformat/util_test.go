package goformat

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUtil(t *testing.T) {
	Convey("cut space", t, func() {
		for _, test := range []struct {
			src, before, middle, after []byte
		}{
			{src: []byte{'a', ' ', ' ', 'b'}, before: []byte{}, middle: []byte{'a', ' ', ' ', 'b'}, after: []byte{}},
			{src: []byte{' ', 'a', 'b', ' '}, before: []byte{' '}, middle: []byte{'a', 'b'}, after: []byte{' '}},
		} {
			before, middle, after := cutSpace(test.src)
			So(before, ShouldResemble, test.before)
			So(middle, ShouldResemble, test.middle)
			So(after, ShouldResemble, test.after)
		}
	})

	Convey("match space", t, func() {
		for _, test := range []struct {
			orig, src, ret []byte
		}{
			{
				orig: []byte{' ', ' ', '3', '4', ' ', ' ', ' '},
				src:  []byte{' ', 'a', ' ', ' ', 'b', ' ', ' '},
				ret:  []byte{' ', ' ', 'a', ' ', ' ', 'b', ' ', ' ', ' '},
			},
		} {
			ret := matchSpace(test.orig, test.src)
			So(ret, ShouldResemble, test.ret)
		}
	})
}
