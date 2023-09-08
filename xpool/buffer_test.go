package xpool

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBufferPool(t *testing.T) {
	Convey("bytes pool", t, func() {
		debug = true
		So(func() { NewSyncBytesPool(0, 100, 2) }, ShouldPanic)
		p := NewSyncBytesPool(1, 100, 2)
		buff0 := p.Alloc(6)
		So(p.(*SyncBytesPool).allocTimesFromPool, ShouldEqual, 1)
		buff1 := p.Alloc(101)
		So(p.(*SyncBytesPool).allocTimesFromPool, ShouldEqual, 1)
		p.Free(buff0)
		So(p.(*SyncBytesPool).freeTimesToPool, ShouldEqual, 1)
		p.Free(buff1)
		So(p.(*SyncBytesPool).freeTimesToPool, ShouldEqual, 1)
		var frame = make([]byte, 7)
		p.Free(frame)
		So(func() { _ = p.Alloc(8) }, ShouldNotPanic)

		frame = make([]byte, 2, 2)
		p = NewSyncBytesPool(4, 100, 2)
		p.Free(frame)
		So(func() { _ = p.Alloc(4) }, ShouldNotPanic)
	})
}
