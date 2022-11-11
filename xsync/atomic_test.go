package xsync

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAtomicString(t *testing.T) {
	s := NewAtomicString("")
	if s.Get() != "" {
		t.Errorf("want empty, got %s", s.Get())
	}
	s.Set("a")
	if s.Get() != "a" {
		t.Errorf("want a, got %s", s.Get())
	}
}

func TestAtomicBool(t *testing.T) {
	Convey("atomic bool", t, func() {
		var b AtomicBool
		So(b.Get(), ShouldBeFalse)
		b.Set(true)
		So(b.Get(), ShouldBeTrue)
		b.Set(false)
		So(b.Get(), ShouldBeFalse)
	})
}
