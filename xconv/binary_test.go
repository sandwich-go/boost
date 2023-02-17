package xconv

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBinary(t *testing.T) {
	Convey(`test binary`, t, func() {
		var buf = []byte{0x00, 0x00, 0x48, 0x42}
		So(LittleEndianDecodeToFloat32(buf), ShouldEqual, float32(50))
		buf = []byte{24, 45, 68, 84, 251, 33, 9, 64}
		So(LittleEndianDecodeToFloat64(buf), ShouldEqual, 3.141592653589793)
		buf = []byte{0x48}
		So(LittleEndianDecodeToInt64(buf), ShouldEqual, 72)
		So(LittleEndianDecodeToUint64(buf), ShouldEqual, 72)
	})
}
