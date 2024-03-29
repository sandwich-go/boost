// Code generated by tools. DO NOT EDIT.
package xslice

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFloat32s(t *testing.T) {
	Convey("float32 slice", t, func() {
		for _, test := range []struct {
			ss       []float32
			s        float32
			contains bool
		}{
			{ss: nil, s: 1},
			{ss: []float32{1, 2}, s: 3, contains: false},
			{ss: []float32{1, 2}, s: 1, contains: true},
		} {
			So(Float32sContain(test.ss, test.s), ShouldEqual, test.contains)
		}
		dest := Float32sSetAdd(nil, 1)
		So(len(dest), ShouldEqual, 1)
		So(len(Float32sSetAdd(dest, 1)), ShouldEqual, 1)

		src := []float32{1, 2}
		dest = Float32sWalk(src, func(s float32) (float32, bool) {
			return s + 1, true
		})
		So(len(src), ShouldEqual, len(dest))
		for i := 0; i < len(src); i++ {
			So(src[i]+1, ShouldEqual, dest[i])
		}
		for i := 0; i < 2; i++ {
			if i > 0 {
				tooManyElement = 4
			}
			for _, test := range []struct {
				src      []float32
				dest     []float32
				contains bool
			}{
				{src: []float32{1, 2, 2}, dest: []float32{1, 2}},
				{src: []float32{1, 1, 2}, dest: []float32{1, 2}},
				{src: []float32{1, 1, 2, 2, 2}, dest: []float32{1, 2}},
			} {
				So(Float32sRemoveRepeated(test.src), ShouldResemble, test.dest)
			}
		}

		for _, test := range []struct {
			src      []float32
			dest     []float32
			contains bool
		}{
			{src: []float32{1, 0, 2}, dest: []float32{1, 2}},
			{src: []float32{1, 1, 0}, dest: []float32{1, 1}},
			{src: []float32{1, 0, 0, 2, 0}, dest: []float32{1, 2}},
		} {
			v := Float32sRemoveEmpty(test.src)
			So(v, ShouldResemble, test.dest)
			Float32sShuffle(v)
		}
	})
}

func TestFloat64s(t *testing.T) {
	Convey("float64 slice", t, func() {
		for _, test := range []struct {
			ss       []float64
			s        float64
			contains bool
		}{
			{ss: nil, s: 1},
			{ss: []float64{1, 2}, s: 3, contains: false},
			{ss: []float64{1, 2}, s: 1, contains: true},
		} {
			So(Float64sContain(test.ss, test.s), ShouldEqual, test.contains)
		}
		dest := Float64sSetAdd(nil, 1)
		So(len(dest), ShouldEqual, 1)
		So(len(Float64sSetAdd(dest, 1)), ShouldEqual, 1)

		src := []float64{1, 2}
		dest = Float64sWalk(src, func(s float64) (float64, bool) {
			return s + 1, true
		})
		So(len(src), ShouldEqual, len(dest))
		for i := 0; i < len(src); i++ {
			So(src[i]+1, ShouldEqual, dest[i])
		}
		for i := 0; i < 2; i++ {
			if i > 0 {
				tooManyElement = 4
			}
			for _, test := range []struct {
				src      []float64
				dest     []float64
				contains bool
			}{
				{src: []float64{1, 2, 2}, dest: []float64{1, 2}},
				{src: []float64{1, 1, 2}, dest: []float64{1, 2}},
				{src: []float64{1, 1, 2, 2, 2}, dest: []float64{1, 2}},
			} {
				So(Float64sRemoveRepeated(test.src), ShouldResemble, test.dest)
			}
		}

		for _, test := range []struct {
			src      []float64
			dest     []float64
			contains bool
		}{
			{src: []float64{1, 0, 2}, dest: []float64{1, 2}},
			{src: []float64{1, 1, 0}, dest: []float64{1, 1}},
			{src: []float64{1, 0, 0, 2, 0}, dest: []float64{1, 2}},
		} {
			v := Float64sRemoveEmpty(test.src)
			So(v, ShouldResemble, test.dest)
			Float64sShuffle(v)
		}
	})
}

func TestInts(t *testing.T) {
	Convey("int slice", t, func() {
		for _, test := range []struct {
			ss       []int
			s        int
			contains bool
		}{
			{ss: nil, s: 1},
			{ss: []int{1, 2}, s: 3, contains: false},
			{ss: []int{1, 2}, s: 1, contains: true},
		} {
			So(IntsContain(test.ss, test.s), ShouldEqual, test.contains)
		}
		dest := IntsSetAdd(nil, 1)
		So(len(dest), ShouldEqual, 1)
		So(len(IntsSetAdd(dest, 1)), ShouldEqual, 1)

		src := []int{1, 2}
		dest = IntsWalk(src, func(s int) (int, bool) {
			return s + 1, true
		})
		So(len(src), ShouldEqual, len(dest))
		for i := 0; i < len(src); i++ {
			So(src[i]+1, ShouldEqual, dest[i])
		}
		for i := 0; i < 2; i++ {
			if i > 0 {
				tooManyElement = 4
			}
			for _, test := range []struct {
				src      []int
				dest     []int
				contains bool
			}{
				{src: []int{1, 2, 2}, dest: []int{1, 2}},
				{src: []int{1, 1, 2}, dest: []int{1, 2}},
				{src: []int{1, 1, 2, 2, 2}, dest: []int{1, 2}},
			} {
				So(IntsRemoveRepeated(test.src), ShouldResemble, test.dest)
			}
		}

		for _, test := range []struct {
			src      []int
			dest     []int
			contains bool
		}{
			{src: []int{1, 0, 2}, dest: []int{1, 2}},
			{src: []int{1, 1, 0}, dest: []int{1, 1}},
			{src: []int{1, 0, 0, 2, 0}, dest: []int{1, 2}},
		} {
			v := IntsRemoveEmpty(test.src)
			So(v, ShouldResemble, test.dest)
			IntsShuffle(v)
		}
	})
}

func TestInt16s(t *testing.T) {
	Convey("int16 slice", t, func() {
		for _, test := range []struct {
			ss       []int16
			s        int16
			contains bool
		}{
			{ss: nil, s: 1},
			{ss: []int16{1, 2}, s: 3, contains: false},
			{ss: []int16{1, 2}, s: 1, contains: true},
		} {
			So(Int16sContain(test.ss, test.s), ShouldEqual, test.contains)
		}
		dest := Int16sSetAdd(nil, 1)
		So(len(dest), ShouldEqual, 1)
		So(len(Int16sSetAdd(dest, 1)), ShouldEqual, 1)

		src := []int16{1, 2}
		dest = Int16sWalk(src, func(s int16) (int16, bool) {
			return s + 1, true
		})
		So(len(src), ShouldEqual, len(dest))
		for i := 0; i < len(src); i++ {
			So(src[i]+1, ShouldEqual, dest[i])
		}
		for i := 0; i < 2; i++ {
			if i > 0 {
				tooManyElement = 4
			}
			for _, test := range []struct {
				src      []int16
				dest     []int16
				contains bool
			}{
				{src: []int16{1, 2, 2}, dest: []int16{1, 2}},
				{src: []int16{1, 1, 2}, dest: []int16{1, 2}},
				{src: []int16{1, 1, 2, 2, 2}, dest: []int16{1, 2}},
			} {
				So(Int16sRemoveRepeated(test.src), ShouldResemble, test.dest)
			}
		}

		for _, test := range []struct {
			src      []int16
			dest     []int16
			contains bool
		}{
			{src: []int16{1, 0, 2}, dest: []int16{1, 2}},
			{src: []int16{1, 1, 0}, dest: []int16{1, 1}},
			{src: []int16{1, 0, 0, 2, 0}, dest: []int16{1, 2}},
		} {
			v := Int16sRemoveEmpty(test.src)
			So(v, ShouldResemble, test.dest)
			Int16sShuffle(v)
		}
	})
}

func TestInt32s(t *testing.T) {
	Convey("int32 slice", t, func() {
		for _, test := range []struct {
			ss       []int32
			s        int32
			contains bool
		}{
			{ss: nil, s: 1},
			{ss: []int32{1, 2}, s: 3, contains: false},
			{ss: []int32{1, 2}, s: 1, contains: true},
		} {
			So(Int32sContain(test.ss, test.s), ShouldEqual, test.contains)
		}
		dest := Int32sSetAdd(nil, 1)
		So(len(dest), ShouldEqual, 1)
		So(len(Int32sSetAdd(dest, 1)), ShouldEqual, 1)

		src := []int32{1, 2}
		dest = Int32sWalk(src, func(s int32) (int32, bool) {
			return s + 1, true
		})
		So(len(src), ShouldEqual, len(dest))
		for i := 0; i < len(src); i++ {
			So(src[i]+1, ShouldEqual, dest[i])
		}
		for i := 0; i < 2; i++ {
			if i > 0 {
				tooManyElement = 4
			}
			for _, test := range []struct {
				src      []int32
				dest     []int32
				contains bool
			}{
				{src: []int32{1, 2, 2}, dest: []int32{1, 2}},
				{src: []int32{1, 1, 2}, dest: []int32{1, 2}},
				{src: []int32{1, 1, 2, 2, 2}, dest: []int32{1, 2}},
			} {
				So(Int32sRemoveRepeated(test.src), ShouldResemble, test.dest)
			}
		}

		for _, test := range []struct {
			src      []int32
			dest     []int32
			contains bool
		}{
			{src: []int32{1, 0, 2}, dest: []int32{1, 2}},
			{src: []int32{1, 1, 0}, dest: []int32{1, 1}},
			{src: []int32{1, 0, 0, 2, 0}, dest: []int32{1, 2}},
		} {
			v := Int32sRemoveEmpty(test.src)
			So(v, ShouldResemble, test.dest)
			Int32sShuffle(v)
		}
	})
}

func TestInt64s(t *testing.T) {
	Convey("int64 slice", t, func() {
		for _, test := range []struct {
			ss       []int64
			s        int64
			contains bool
		}{
			{ss: nil, s: 1},
			{ss: []int64{1, 2}, s: 3, contains: false},
			{ss: []int64{1, 2}, s: 1, contains: true},
		} {
			So(Int64sContain(test.ss, test.s), ShouldEqual, test.contains)
		}
		dest := Int64sSetAdd(nil, 1)
		So(len(dest), ShouldEqual, 1)
		So(len(Int64sSetAdd(dest, 1)), ShouldEqual, 1)

		src := []int64{1, 2}
		dest = Int64sWalk(src, func(s int64) (int64, bool) {
			return s + 1, true
		})
		So(len(src), ShouldEqual, len(dest))
		for i := 0; i < len(src); i++ {
			So(src[i]+1, ShouldEqual, dest[i])
		}
		for i := 0; i < 2; i++ {
			if i > 0 {
				tooManyElement = 4
			}
			for _, test := range []struct {
				src      []int64
				dest     []int64
				contains bool
			}{
				{src: []int64{1, 2, 2}, dest: []int64{1, 2}},
				{src: []int64{1, 1, 2}, dest: []int64{1, 2}},
				{src: []int64{1, 1, 2, 2, 2}, dest: []int64{1, 2}},
			} {
				So(Int64sRemoveRepeated(test.src), ShouldResemble, test.dest)
			}
		}

		for _, test := range []struct {
			src      []int64
			dest     []int64
			contains bool
		}{
			{src: []int64{1, 0, 2}, dest: []int64{1, 2}},
			{src: []int64{1, 1, 0}, dest: []int64{1, 1}},
			{src: []int64{1, 0, 0, 2, 0}, dest: []int64{1, 2}},
		} {
			v := Int64sRemoveEmpty(test.src)
			So(v, ShouldResemble, test.dest)
			Int64sShuffle(v)
		}
	})
}

func TestInt8s(t *testing.T) {
	Convey("int8 slice", t, func() {
		for _, test := range []struct {
			ss       []int8
			s        int8
			contains bool
		}{
			{ss: nil, s: 1},
			{ss: []int8{1, 2}, s: 3, contains: false},
			{ss: []int8{1, 2}, s: 1, contains: true},
		} {
			So(Int8sContain(test.ss, test.s), ShouldEqual, test.contains)
		}
		dest := Int8sSetAdd(nil, 1)
		So(len(dest), ShouldEqual, 1)
		So(len(Int8sSetAdd(dest, 1)), ShouldEqual, 1)

		src := []int8{1, 2}
		dest = Int8sWalk(src, func(s int8) (int8, bool) {
			return s + 1, true
		})
		So(len(src), ShouldEqual, len(dest))
		for i := 0; i < len(src); i++ {
			So(src[i]+1, ShouldEqual, dest[i])
		}
		for i := 0; i < 2; i++ {
			if i > 0 {
				tooManyElement = 4
			}
			for _, test := range []struct {
				src      []int8
				dest     []int8
				contains bool
			}{
				{src: []int8{1, 2, 2}, dest: []int8{1, 2}},
				{src: []int8{1, 1, 2}, dest: []int8{1, 2}},
				{src: []int8{1, 1, 2, 2, 2}, dest: []int8{1, 2}},
			} {
				So(Int8sRemoveRepeated(test.src), ShouldResemble, test.dest)
			}
		}

		for _, test := range []struct {
			src      []int8
			dest     []int8
			contains bool
		}{
			{src: []int8{1, 0, 2}, dest: []int8{1, 2}},
			{src: []int8{1, 1, 0}, dest: []int8{1, 1}},
			{src: []int8{1, 0, 0, 2, 0}, dest: []int8{1, 2}},
		} {
			v := Int8sRemoveEmpty(test.src)
			So(v, ShouldResemble, test.dest)
			Int8sShuffle(v)
		}
	})
}

func TestStrings(t *testing.T) {
	Convey("string slice", t, func() {
		for _, test := range []struct {
			ss       []string
			s        string
			contains bool
		}{
			{ss: nil, s: "a"},
			{ss: []string{"abc", "b"}, s: "a", contains: false},
			{ss: []string{"abc", "b"}, s: "abc", contains: true},
		} {
			So(StringsContain(test.ss, test.s), ShouldEqual, test.contains)
		}
		for _, test := range []struct {
			ss       []string
			s        string
			contains bool
		}{
			{ss: nil, s: "a"},
			{ss: []string{"abc", "b"}, s: "a", contains: false},
			{ss: []string{"abc", "b"}, s: "abc", contains: true},
			{ss: []string{"ABC", "b"}, s: "abc", contains: true},
		} {
			So(StringsContainEqualFold(test.ss, test.s), ShouldEqual, test.contains)
		}
		dest := StringsSetAdd(nil, "a")
		So(len(dest), ShouldEqual, 1)
		So(len(StringsSetAdd(dest, "a")), ShouldEqual, 1)

		src := []string{"1", "2"}
		dest = StringsWalk(src, func(s string) (string, bool) {
			return s + ",", true
		})
		So(len(src), ShouldEqual, len(dest))
		for i := 0; i < len(src); i++ {
			So(src[i]+",", ShouldEqual, dest[i])
		}

		dest = StringsAddSuffix(src, ",")
		So(len(src), ShouldEqual, len(dest))
		for i := 0; i < len(src); i++ {
			So(src[i]+",", ShouldEqual, dest[i])
		}

		dest = StringsAddPrefix(src, ",")
		So(len(src), ShouldEqual, len(dest))
		for i := 0; i < len(src); i++ {
			So(","+src[i], ShouldEqual, dest[i])
		}
		for i := 0; i < 2; i++ {
			if i > 0 {
				tooManyElement = 4
			}
			for _, test := range []struct {
				src      []string
				dest     []string
				contains bool
			}{
				{src: []string{"abc", "b", "b"}, dest: []string{"abc", "b"}},
				{src: []string{"abc", "abc", "b"}, dest: []string{"abc", "b"}},
				{src: []string{"abc", "abc", "b", "b", "b"}, dest: []string{"abc", "b"}},
			} {
				So(StringsRemoveRepeated(test.src), ShouldResemble, test.dest)
			}
		}

		for _, test := range []struct {
			src      []string
			dest     []string
			contains bool
		}{
			{src: []string{"abc", "", "b"}, dest: []string{"abc", "b"}},
			{src: []string{"abc", "abc", ""}, dest: []string{"abc", "abc"}},
			{src: []string{"abc", "", "", "b", ""}, dest: []string{"abc", "b"}},
		} {
			v := StringsRemoveEmpty(test.src)
			So(v, ShouldResemble, test.dest)
			StringsShuffle(v)
		}
	})
}

func TestUints(t *testing.T) {
	Convey("uint slice", t, func() {
		for _, test := range []struct {
			ss       []uint
			s        uint
			contains bool
		}{
			{ss: nil, s: 1},
			{ss: []uint{1, 2}, s: 3, contains: false},
			{ss: []uint{1, 2}, s: 1, contains: true},
		} {
			So(UintsContain(test.ss, test.s), ShouldEqual, test.contains)
		}
		dest := UintsSetAdd(nil, 1)
		So(len(dest), ShouldEqual, 1)
		So(len(UintsSetAdd(dest, 1)), ShouldEqual, 1)

		src := []uint{1, 2}
		dest = UintsWalk(src, func(s uint) (uint, bool) {
			return s + 1, true
		})
		So(len(src), ShouldEqual, len(dest))
		for i := 0; i < len(src); i++ {
			So(src[i]+1, ShouldEqual, dest[i])
		}
		for i := 0; i < 2; i++ {
			if i > 0 {
				tooManyElement = 4
			}
			for _, test := range []struct {
				src      []uint
				dest     []uint
				contains bool
			}{
				{src: []uint{1, 2, 2}, dest: []uint{1, 2}},
				{src: []uint{1, 1, 2}, dest: []uint{1, 2}},
				{src: []uint{1, 1, 2, 2, 2}, dest: []uint{1, 2}},
			} {
				So(UintsRemoveRepeated(test.src), ShouldResemble, test.dest)
			}
		}

		for _, test := range []struct {
			src      []uint
			dest     []uint
			contains bool
		}{
			{src: []uint{1, 0, 2}, dest: []uint{1, 2}},
			{src: []uint{1, 1, 0}, dest: []uint{1, 1}},
			{src: []uint{1, 0, 0, 2, 0}, dest: []uint{1, 2}},
		} {
			v := UintsRemoveEmpty(test.src)
			So(v, ShouldResemble, test.dest)
			UintsShuffle(v)
		}
	})
}

func TestUint16s(t *testing.T) {
	Convey("uint16 slice", t, func() {
		for _, test := range []struct {
			ss       []uint16
			s        uint16
			contains bool
		}{
			{ss: nil, s: 1},
			{ss: []uint16{1, 2}, s: 3, contains: false},
			{ss: []uint16{1, 2}, s: 1, contains: true},
		} {
			So(Uint16sContain(test.ss, test.s), ShouldEqual, test.contains)
		}
		dest := Uint16sSetAdd(nil, 1)
		So(len(dest), ShouldEqual, 1)
		So(len(Uint16sSetAdd(dest, 1)), ShouldEqual, 1)

		src := []uint16{1, 2}
		dest = Uint16sWalk(src, func(s uint16) (uint16, bool) {
			return s + 1, true
		})
		So(len(src), ShouldEqual, len(dest))
		for i := 0; i < len(src); i++ {
			So(src[i]+1, ShouldEqual, dest[i])
		}
		for i := 0; i < 2; i++ {
			if i > 0 {
				tooManyElement = 4
			}
			for _, test := range []struct {
				src      []uint16
				dest     []uint16
				contains bool
			}{
				{src: []uint16{1, 2, 2}, dest: []uint16{1, 2}},
				{src: []uint16{1, 1, 2}, dest: []uint16{1, 2}},
				{src: []uint16{1, 1, 2, 2, 2}, dest: []uint16{1, 2}},
			} {
				So(Uint16sRemoveRepeated(test.src), ShouldResemble, test.dest)
			}
		}

		for _, test := range []struct {
			src      []uint16
			dest     []uint16
			contains bool
		}{
			{src: []uint16{1, 0, 2}, dest: []uint16{1, 2}},
			{src: []uint16{1, 1, 0}, dest: []uint16{1, 1}},
			{src: []uint16{1, 0, 0, 2, 0}, dest: []uint16{1, 2}},
		} {
			v := Uint16sRemoveEmpty(test.src)
			So(v, ShouldResemble, test.dest)
			Uint16sShuffle(v)
		}
	})
}

func TestUint32s(t *testing.T) {
	Convey("uint32 slice", t, func() {
		for _, test := range []struct {
			ss       []uint32
			s        uint32
			contains bool
		}{
			{ss: nil, s: 1},
			{ss: []uint32{1, 2}, s: 3, contains: false},
			{ss: []uint32{1, 2}, s: 1, contains: true},
		} {
			So(Uint32sContain(test.ss, test.s), ShouldEqual, test.contains)
		}
		dest := Uint32sSetAdd(nil, 1)
		So(len(dest), ShouldEqual, 1)
		So(len(Uint32sSetAdd(dest, 1)), ShouldEqual, 1)

		src := []uint32{1, 2}
		dest = Uint32sWalk(src, func(s uint32) (uint32, bool) {
			return s + 1, true
		})
		So(len(src), ShouldEqual, len(dest))
		for i := 0; i < len(src); i++ {
			So(src[i]+1, ShouldEqual, dest[i])
		}
		for i := 0; i < 2; i++ {
			if i > 0 {
				tooManyElement = 4
			}
			for _, test := range []struct {
				src      []uint32
				dest     []uint32
				contains bool
			}{
				{src: []uint32{1, 2, 2}, dest: []uint32{1, 2}},
				{src: []uint32{1, 1, 2}, dest: []uint32{1, 2}},
				{src: []uint32{1, 1, 2, 2, 2}, dest: []uint32{1, 2}},
			} {
				So(Uint32sRemoveRepeated(test.src), ShouldResemble, test.dest)
			}
		}

		for _, test := range []struct {
			src      []uint32
			dest     []uint32
			contains bool
		}{
			{src: []uint32{1, 0, 2}, dest: []uint32{1, 2}},
			{src: []uint32{1, 1, 0}, dest: []uint32{1, 1}},
			{src: []uint32{1, 0, 0, 2, 0}, dest: []uint32{1, 2}},
		} {
			v := Uint32sRemoveEmpty(test.src)
			So(v, ShouldResemble, test.dest)
			Uint32sShuffle(v)
		}
	})
}

func TestUint64s(t *testing.T) {
	Convey("uint64 slice", t, func() {
		for _, test := range []struct {
			ss       []uint64
			s        uint64
			contains bool
		}{
			{ss: nil, s: 1},
			{ss: []uint64{1, 2}, s: 3, contains: false},
			{ss: []uint64{1, 2}, s: 1, contains: true},
		} {
			So(Uint64sContain(test.ss, test.s), ShouldEqual, test.contains)
		}
		dest := Uint64sSetAdd(nil, 1)
		So(len(dest), ShouldEqual, 1)
		So(len(Uint64sSetAdd(dest, 1)), ShouldEqual, 1)

		src := []uint64{1, 2}
		dest = Uint64sWalk(src, func(s uint64) (uint64, bool) {
			return s + 1, true
		})
		So(len(src), ShouldEqual, len(dest))
		for i := 0; i < len(src); i++ {
			So(src[i]+1, ShouldEqual, dest[i])
		}
		for i := 0; i < 2; i++ {
			if i > 0 {
				tooManyElement = 4
			}
			for _, test := range []struct {
				src      []uint64
				dest     []uint64
				contains bool
			}{
				{src: []uint64{1, 2, 2}, dest: []uint64{1, 2}},
				{src: []uint64{1, 1, 2}, dest: []uint64{1, 2}},
				{src: []uint64{1, 1, 2, 2, 2}, dest: []uint64{1, 2}},
			} {
				So(Uint64sRemoveRepeated(test.src), ShouldResemble, test.dest)
			}
		}

		for _, test := range []struct {
			src      []uint64
			dest     []uint64
			contains bool
		}{
			{src: []uint64{1, 0, 2}, dest: []uint64{1, 2}},
			{src: []uint64{1, 1, 0}, dest: []uint64{1, 1}},
			{src: []uint64{1, 0, 0, 2, 0}, dest: []uint64{1, 2}},
		} {
			v := Uint64sRemoveEmpty(test.src)
			So(v, ShouldResemble, test.dest)
			Uint64sShuffle(v)
		}
	})
}

func TestUint8s(t *testing.T) {
	Convey("uint8 slice", t, func() {
		for _, test := range []struct {
			ss       []uint8
			s        uint8
			contains bool
		}{
			{ss: nil, s: 1},
			{ss: []uint8{1, 2}, s: 3, contains: false},
			{ss: []uint8{1, 2}, s: 1, contains: true},
		} {
			So(Uint8sContain(test.ss, test.s), ShouldEqual, test.contains)
		}
		dest := Uint8sSetAdd(nil, 1)
		So(len(dest), ShouldEqual, 1)
		So(len(Uint8sSetAdd(dest, 1)), ShouldEqual, 1)

		src := []uint8{1, 2}
		dest = Uint8sWalk(src, func(s uint8) (uint8, bool) {
			return s + 1, true
		})
		So(len(src), ShouldEqual, len(dest))
		for i := 0; i < len(src); i++ {
			So(src[i]+1, ShouldEqual, dest[i])
		}
		for i := 0; i < 2; i++ {
			if i > 0 {
				tooManyElement = 4
			}
			for _, test := range []struct {
				src      []uint8
				dest     []uint8
				contains bool
			}{
				{src: []uint8{1, 2, 2}, dest: []uint8{1, 2}},
				{src: []uint8{1, 1, 2}, dest: []uint8{1, 2}},
				{src: []uint8{1, 1, 2, 2, 2}, dest: []uint8{1, 2}},
			} {
				So(Uint8sRemoveRepeated(test.src), ShouldResemble, test.dest)
			}
		}

		for _, test := range []struct {
			src      []uint8
			dest     []uint8
			contains bool
		}{
			{src: []uint8{1, 0, 2}, dest: []uint8{1, 2}},
			{src: []uint8{1, 1, 0}, dest: []uint8{1, 1}},
			{src: []uint8{1, 0, 0, 2, 0}, dest: []uint8{1, 2}},
		} {
			v := Uint8sRemoveEmpty(test.src)
			So(v, ShouldResemble, test.dest)
			Uint8sShuffle(v)
		}
	})
}
