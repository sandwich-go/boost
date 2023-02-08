package z

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type hasher struct{}

func (h hasher) HashUint64() uint64 { return 1 }

type stringer struct{}

func (h stringer) String() string { return "1" }

func TestHash(t *testing.T) {
	Convey("hash", t, func() {
		s := `a {{ .val1 }} {{ .val2 }}`
		s1 := MemHash([]byte(s))
		s2 := MemHashString(s)
		So(s1, ShouldEqual, s2)

		for _, test := range []struct {
			v     interface{}
			hv    uint64
			panic bool
		}{
			{v: nil, hv: 0}, {v: int8(4), hv: 4}, {v: int16(5), hv: 5}, {v: int32(6), hv: 6}, {v: int64(8), hv: 8}, {v: 3, hv: 3},
			{v: byte(2), hv: 2}, {v: uint16(9), hv: 9}, {v: uint64(1), hv: 1}, {v: uint32(7), hv: 7}, {v: uint(10), hv: 10}, {v: uint8(11), hv: 11},
			{v: hasher{}, hv: 1}, {v: stringer{}, hv: MemHashString("1")}, {v: []byte(s), hv: s1}, {v: s, hv: s1},
			{v: float32(1), panic: true},
		} {
			if test.panic {
				So(func() {
					KeyToHash(test.v)
				}, ShouldPanic)
			} else {
				So(KeyToHash(test.v), ShouldEqual, test.hv)
			}
		}
	})
}
