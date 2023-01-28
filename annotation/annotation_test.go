package annotation

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAnnotation(t *testing.T) {
	Convey("different input for Filter", t, func() {
		for _, test := range []struct {
			err      bool
			name     string
			in       []string
			validate func(as []Annotation)
		}{
			// 非注解
			{in: []string{`// wvdwadbvb`}, validate: func(as []Annotation) {
				So(len(as), ShouldEqual, 0)
			}},
			{in: []string{`// annotation@X( a = "A" }`}, name: "X", err: true},
			{in: []string{`// annotation@X( a = A B C D  )`}, name: "X", validate: func(as []Annotation) {
				So(len(as), ShouldEqual, 1)
				So(as[0].String("a"), ShouldEqual, "")
			}},
			{in: []string{`// annotation@X( a = false  )`}, name: "X", validate: func(as []Annotation) {
				So(len(as), ShouldEqual, 1)
				So(as[0].String("a"), ShouldEqual, "")
			}},
			{in: []string{`// annotation@X( a = "A" ) some words`}, name: "X", validate: func(as []Annotation) {
				So(len(as), ShouldEqual, 1)
				So(as[0].String("a"), ShouldEqual, "A")
			}},
			{in: []string{`// annotation@X( a = "A" ) `}, name: "X", validate: func(as []Annotation) {
				So(len(as), ShouldEqual, 1)
				So(as[0].String("a"), ShouldEqual, "A")
			}},
			{in: []string{`// annotation@X( a = "A" )`, `annotation@X( b = "B" )`}, name: "X", validate: func(as []Annotation) {
				So(len(as), ShouldEqual, 2)
				So(as[0].String("a"), ShouldEqual, "A")
				So(as[1].String("b"), ShouldEqual, "B")
			}},
			{in: []string{`// annotation@X( a = "A" )`, `annotation@X( b = "B" )`}, name: "*", validate: func(as []Annotation) {
				So(len(as), ShouldEqual, 2)
				So(as[0].String("a"), ShouldEqual, "A")
				So(as[1].String("b"), ShouldEqual, "B")
			}},
			{in: []string{`// annotation@SomethingElse( aggregate = "@A@")`, `// annotation@Event( aggregate = "@A@")`}, name: "Event", validate: func(as []Annotation) {
				So(len(as), ShouldEqual, 1)
				So(as[0].String("aggregate"), ShouldEqual, "@A@")
			}},
			{in: []string{`// annotation@Doit( a="/A/", b="/B" )`}, name: "Doit", validate: func(as []Annotation) {
				So(len(as), ShouldEqual, 1)
				So(as[0].String("a"), ShouldEqual, "/A/")
				So(as[0].String("b"), ShouldEqual, "/B")
			}},
			{in: []string{`// annotation@SomethingElse( aggregate = "A")`, `// annotation@Event( aggregate = "@A@")`}, name: "Event", validate: func(as []Annotation) {
				So(len(as), ShouldEqual, 1)
				So(as[0].String("aggregate"), ShouldEqual, "@A@")
			}},
		} {
			r := New(WithDescriptors(Descriptor{Name: test.name}))
			ann, err := r.ResolveMany(test.in...)
			if test.err {
				So(err, ShouldNotBeNil)
			} else {
				test.validate(ann)
			}
		}
	})

	Convey("annotation", t, func() {
		for _, test := range []struct {
			// input
			name      string
			in        string
			validator func(ann Annotation) bool
			validate  func(ann Annotation)
			err       error
		}{
			{name: "a", in: `// annotation@a( AK="av" )`, err: ErrNoAnnotation},
			{name: "a", in: `// ann@a( ak="av" )`, validate: func(ann Annotation) {
				So(ann.Name(), ShouldEqual, "a")
				So(ann.Line(), ShouldEqual, `// ann@a( ak="av" )`)
				So(ann.Contains("ak"), ShouldBeTrue)
				So(ann.Contains("AK"), ShouldBeFalse)
				So(ann.String("ak"), ShouldEqual, "av")
				_, err := ann.Int64("ak")
				So(err, ShouldNotBeNil)
			}},
			{name: "a", in: `// ann@a( AK="av" )`, validate: func(ann Annotation) {
				So(ann.String("ak"), ShouldEqual, "")
				So(ann.String("AK"), ShouldEqual, "av")
			}},
			{name: "a", in: `// ann@a( AK="av" )`, validator: func(ann Annotation) bool { return ann.String("AK") != "av" }, err: ErrNoAnnotation},
			{name: "a", in: `// ann@a( AK=127, AV=128, AF=0.061, AB=1, INVALID="AAAAA" )`, validate: func(ann Annotation) {
				var check = func(index int, expected, defaultVal interface{}, f func() (interface{}, error)) {
					val, err := f()
					switch index {
					case 0:
						So(err, ShouldBeNil)
						So(val, ShouldEqual, expected)
					case 1:
						So(err, ShouldBeNil)
						So(val, ShouldEqual, defaultVal)
					case 2:
						So(err, ShouldNotBeNil)
					}
				}
				for i, key := range []string{"AK", "no_exists", "INVALID"} {
					for _, defaultVal := range []interface{}{int8(126), int16(128), int32(126), int64(126), 126, uint8(126), uint16(128), uint32(126), uint64(126)} {
						switch vv := defaultVal.(type) {
						case int8:
							check(i, int8(127), vv, func() (interface{}, error) { return ann.Int8(key, vv) })
						case int16:
							check(i, int16(127), vv, func() (interface{}, error) { return ann.Int16(key, vv) })
						case int32:
							check(i, int32(127), vv, func() (interface{}, error) { return ann.Int32(key, vv) })
						case int64:
							check(i, int64(127), vv, func() (interface{}, error) { return ann.Int64(key, vv) })
						case int:
							check(i, 127, vv, func() (interface{}, error) { return ann.Int(key, vv) })
						case uint8:
							check(i, uint8(127), vv, func() (interface{}, error) { return ann.Uint8(key, vv) })
						case uint16:
							check(i, uint16(127), vv, func() (interface{}, error) { return ann.Uint16(key, vv) })
						case uint32:
							check(i, uint32(127), vv, func() (interface{}, error) { return ann.Uint32(key, vv) })
						case uint64:
							check(i, uint64(127), vv, func() (interface{}, error) { return ann.Uint64(key, vv) })
						}
					}
				}

				for i, key := range []string{"AF", "no_exists", "INVALID"} {
					for _, defaultVal := range []interface{}{float32(0.051), 0.051} {
						switch vv := defaultVal.(type) {
						case float32:
							check(i, float32(0.061), vv, func() (interface{}, error) { return ann.Float32(key, vv) })
						case float64:
							check(i, 0.061, vv, func() (interface{}, error) { return ann.Float64(key, vv) })
						}
					}
				}

				for i, key := range []string{"AB", "no_exists"} {
					for _, defaultVal := range []interface{}{true} {
						switch vv := defaultVal.(type) {
						case bool:
							check(i, true, vv, func() (interface{}, error) { return ann.Bool(key, vv) })
						}
					}
				}
				_, err := ann.Int8("AV")
				So(err, ShouldNotBeNil)
			}},
		} {
			r := New(WithMagicPrefix("ann@"), WithDescriptors(Descriptor{Name: test.name, Validator: test.validator}), WithLowerKey(false))
			ann, err := r.Resolve(test.in)
			if test.err != nil {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, test.err.Error())
			} else {
				So(err, ShouldBeNil)
				test.validate(ann)
			}
		}
	})

	Convey("resolve", t, func() {
		ann, err := Resolve(`// wvdwadbvb`)
		So(ann, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, ErrNoAnnotation)

		ann, err = Resolve(`// annotation@X( A = "A" )`)
		So(err, ShouldBeNil)
		So(ann, ShouldNotBeNil)
		So(ann.String("a"), ShouldEqual, "A")
		So(ann.String("A"), ShouldEqual, "")
	})

	Convey("resolve with name", t, func() {
		ann, err := ResolveWithName("Y", `// annotation@X( A = "A" )`)
		So(ann, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, ErrNoAnnotation)

		ann, err = ResolveWithName("X", `// annotation@X( A = "A" }`)
		So(err, ShouldNotBeNil)

		ann, err = ResolveWithName("X", `// annotation@X( A = "A" )`, `// annotation@X( B = "B" )`)
		So(err, ShouldBeNil)
		So(ann, ShouldNotBeNil)
		So(ann.String("a"), ShouldEqual, "A")
		So(ann.String("b"), ShouldBeEmpty)
	})

	Convey("resolve many", t, func() {
		ann, err := ResolveMany(`// annotation@X( A = "A" }`, `// annotation@X( A = "A" )`)
		So(err, ShouldNotBeNil)
		So(len(ann), ShouldEqual, 0)

		ann, err = ResolveMany(`// wvdwadbvb`, `// annotation@X( A = "A" )`)
		So(err, ShouldBeNil)
		So(ann, ShouldNotBeNil)
		So(len(ann), ShouldEqual, 1)
		So(ann[0].String("a"), ShouldEqual, "A")
	})

	Convey("resolve with name", t, func() {
		ann, err := ResolveWithName("Y", `// annotation@X( A = "A" )`)
		So(ann, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err, ShouldEqual, ErrNoAnnotation)

		ann, err = ResolveWithName("X", `// annotation@X( A = "A" )`)
		So(err, ShouldBeNil)
		So(ann, ShouldNotBeNil)
		So(ann.String("a"), ShouldEqual, "A")
	})

	Convey("resolve no duplicate", t, func() {
		ann, err := ResolveNoDuplicate(`// annotation@X( A = "A") `, `// annotation@X( a = "a" )`)
		So(ann, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "got duplicate annotation")

		ann, err = ResolveNoDuplicate(`// annotation@X( A = "A") `, `// wvdwadbvb`)
		So(err, ShouldBeNil)
		So(ann, ShouldNotBeNil)
		So(len(ann), ShouldEqual, 1)
		So(ann[0].Name(), ShouldEqual, "X")
		So(ann[0].String("a"), ShouldEqual, "A")

		ann, err = ResolveNoDuplicate(`// annotation@X( A = "A") `, `// annotation@Y( a = "a" )`, `// wvdwadbvb`)
		So(err, ShouldBeNil)
		So(ann, ShouldNotBeNil)
		So(len(ann), ShouldEqual, 2)
		So(ann[0].Name(), ShouldEqual, "X")
		So(ann[0].String("a"), ShouldEqual, "A")
		So(ann[1].Name(), ShouldEqual, "Y")
		So(ann[1].String("a"), ShouldEqual, "a")
	})
}
