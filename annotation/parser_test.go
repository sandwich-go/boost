package annotation

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestParser(t *testing.T) {
	Convey("parser", t, func() {
		for _, test := range []struct {
			line     string
			validate func(Annotation)
			lowKey   bool
		}{
			{`@X( a = "A" }`, nil, true},
			{`@ X( a = "A" )`, func(ann Annotation) {
				So(ann.Name(), ShouldEqual, "X")
				So(ann.String("a"), ShouldEqual, "A")
			}, true},
			{`@X( a = "A" )`, func(ann Annotation) {
				So(ann.Name(), ShouldEqual, "X")
				So(ann.String("a"), ShouldEqual, "A")
			}, true},
			{`@X( A = "A" )`, func(ann Annotation) {
				So(ann.Name(), ShouldEqual, "X")
				So(ann.String("a"), ShouldEqual, "A")
			}, true},
			{`//aaaaa@X( A = "A" )`, func(ann Annotation) {
				So(ann.Name(), ShouldEqual, "X")
				So(ann.String("a"), ShouldEqual, "A")
			}, true},
			{`@ X ( A = "A" , b = "B", c= "c" )`, func(ann Annotation) {
				So(ann.Name(), ShouldEqual, "X")
				So(ann.String("a"), ShouldEqual, "A")
				So(ann.String("b"), ShouldEqual, "B")
				So(ann.String("c"), ShouldEqual, "c")
			}, true},
			{`@ X ( A = "A" , b = "B", c= "c" )`, func(ann Annotation) {
				So(ann.Name(), ShouldEqual, "X")
				So(ann.String("A"), ShouldEqual, "A")
				So(ann.String("b"), ShouldEqual, "B")
				So(ann.String("c"), ShouldEqual, "c")
			}, false},
		} {
			ann, err := parser(test.line, test.lowKey)
			if test.validate != nil {
				So(err, ShouldBeNil)
				test.validate(ann)
			} else {
				So(err, ShouldNotBeNil)
			}
		}
	})
}
