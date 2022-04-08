package annotation

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func validateOk(annot Annotation) bool {
	return true
}

func TestAnnotation(t *testing.T) {
	Convey("different input for Filter", t, func() {
		for _, test := range []struct {
			// input
			panic    bool
			name     string
			in       []string
			validate func(as []Annotation)
		}{
			// 非注解
			{in: []string{`// wvdwadbvb`}, validate: func(as []Annotation) {
				So(len(as), ShouldEqual, 0)
			}},
			// 格式不合法，panic
			{in: []string{`// annotation@X( a = "A" }`}, name: "X", panic: true},
			{in: []string{`// annotation@X( a = A B C D  )`}, name: "X", validate: func(as []Annotation) {
				So(len(as), ShouldEqual, 1)
				So(as[0].Attributes["a"], ShouldEqual, "")
			}},
			{in: []string{`// annotation@X( a = false  )`}, name: "X", validate: func(as []Annotation) {
				So(len(as), ShouldEqual, 1)
				So(as[0].Attributes["a"], ShouldEqual, "")
			}},
			{in: []string{`// annotation@X( a = "A" ) some words`}, name: "X", validate: func(as []Annotation) {
				So(len(as), ShouldEqual, 1)
				So(as[0].Attributes["a"], ShouldEqual, "A")
			}},
			{in: []string{`// annotation@X( a = "A" ) `}, name: "X", validate: func(as []Annotation) {
				So(len(as), ShouldEqual, 1)
				So(as[0].Attributes["a"], ShouldEqual, "A")
			}},
			{in: []string{`// annotation@X( a = "A" )`, `annotation@X( b = "B" )`}, name: "X", validate: func(as []Annotation) {
				So(len(as), ShouldEqual, 2)
				So(as[0].Attributes["a"], ShouldEqual, "A")
				So(as[1].Attributes["b"], ShouldEqual, "B")
			}},
			{in: []string{`// annotation@SomethingElse( aggregate = "@A@")`, `// annotation@Event( aggregate = "@A@")`}, name: "Event", validate: func(as []Annotation) {
				So(len(as), ShouldEqual, 1)
				So(as[0].Attributes["aggregate"], ShouldEqual, "@A@")
			}},
			{in: []string{`// annotation@Doit( a="/A/", b="/B" )`}, name: "Doit", validate: func(as []Annotation) {
				So(len(as), ShouldEqual, 1)
				So(as[0].Attributes["a"], ShouldEqual, "/A/")
				So(as[0].Attributes["b"], ShouldEqual, "/B")
			}},
		} {
			registry := NewRegistry(&Descriptor{
				Name:      test.name,
				Validator: validateOk,
			})
			if test.panic {
				So(func() { registry.ResolveAnnotations(test.in) }, ShouldPanic)
			} else {
				as := registry.ResolveAnnotations(test.in)
				test.validate(as)
			}
		}
	})
}
