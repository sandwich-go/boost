package xbit

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestField(t *testing.T) {
	Convey("set invalid field bit, should panic", t, func() {
		So(func() {
			mustValidField(FieldMax + 1)
		}, ShouldPanic)

		So(func() {
			var fs FieldSet
			fs.Set(FieldMax + 1)
		}, ShouldPanic)
	})

	Convey("field bit option", t, func() {
		var fs0 FieldSet
		var fs1 FieldSet
		var f0, f1, f2 Field = 1, 2, 3
		So(fs0.Set(f0), ShouldBeTrue)
		So(fs0.Set(f0), ShouldBeFalse)
		So(fs0.IsSet(f0), ShouldBeTrue)

		So(fs0.IsSet(f1), ShouldBeFalse)
		So(fs0.Clear(f1), ShouldBeFalse)

		So(fs0.IsSet(f0), ShouldBeTrue)

		So(fs0.Set(f1), ShouldBeTrue)
		So(fs0.IsSet(f1), ShouldBeTrue)
		So(fs0.Clear(f1), ShouldBeTrue)
		So(fs0.IsSet(f1), ShouldBeFalse)
		So(fs0.Set(f1), ShouldBeTrue)
		So(fs0.IsSet(f1), ShouldBeTrue)

		So(fs0.IsSet(f0), ShouldBeTrue)

		So(fs1.Clear(f2), ShouldBeFalse)
		So(fs1.Set(f2), ShouldBeTrue)
		So(fs1.IsSet(f2), ShouldBeTrue)

		So(fs0.IsSet(f2), ShouldBeFalse)

		fs0 = fs0.Union(fs1)
		So(fs0.IsSet(f0), ShouldBeTrue)
		So(fs0.IsSet(f1), ShouldBeTrue)
		So(fs0.IsSet(f2), ShouldBeTrue)

		fs0.ClearAll()
		So(fs0.IsSet(f0), ShouldBeFalse)
		So(fs0.IsSet(f1), ShouldBeFalse)
		So(fs0.IsSet(f2), ShouldBeFalse)
	})
}
