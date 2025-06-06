package misc

import (
	. "github.com/smartystreets/goconvey/convey"
	"reflect"
	"testing"
	"unsafe"
)

var offset uintptr

func init() {
	b := &B{
		A: A{1},
	}
	a := &b.A

	AMetadata.ModelType = reflect.TypeOf(b)
	AMetadata.offset = uintptr(unsafe.Pointer(b)) - uintptr(unsafe.Pointer(a))
}

type A struct {
	haha uint8
}

func (a *A) Locker() string {
	v := UnsafeShift1[A](a, -offset, 1, AMetadata.ModelType)
	return v.(*B).Locker()
}

func (a *A) Metadata() *Metadata {
	return AMetadata
}

var AMetadata = &Metadata{}

type Metadata struct {
	ModelType reflect.Type
	offset    uintptr
}

type B struct {
	flag uint8
	A
	data interface{}
}

func (b *B) Locker() string {
	return "B"
}

func TestCalcOffset(t *testing.T) {
	b := &B{
		A: A{1},
	}
	a := &b.A
	t.Log(a.Locker())
}

type TestStruct struct {
	privateField int
}

func (ts *TestStruct) GetPrivateField() int {
	return ts.privateField
}

func TestExportField(t *testing.T) {
	testStruct := TestStruct{}
	typ := reflect.TypeOf(&testStruct).Elem()
	val := reflect.ValueOf(&testStruct).Elem()
	f := ExportField(val, typ.Field(0))
	f.Set(reflect.ValueOf(1))

	if testStruct.GetPrivateField() != 1 {
		t.Fail()
	}
}

type MockAnonymousStruct struct {
	privateField int
	PublicField  int
}

func (x MockAnonymousStruct) Equal(x1 MockStruct) bool {
	return x.PublicField == x1.PublicField && x.privateField == x1.privateField
}

type MockStruct struct {
	MockAnonymousStruct
}

func (x MockStruct) Equal(x1 MockAnonymousStruct) bool {
	return x.PublicField == x1.PublicField && x.privateField == x1.privateField
}

func getAnonymousOffset[E any]() uintptr {
	ts := reflect.TypeFor[E]().Elem()
	for i := 0; i < ts.NumField(); i++ {
		if f := ts.Field(i); f.Anonymous {
			return f.Offset
		}
	}
	panic("not found anonymous field")
}

func TestTypes(t *testing.T) {
	Convey("use anonymous field, get parent message", t, func() {
		var m MockAnonymousStruct
		m.PublicField = 1
		m.privateField = 2
		et := reflect.TypeFor[*MockStruct]()
		o := getAnonymousOffset[*MockStruct]()

		a, ok := UnsafeShift0(m, o, -1, et).(*MockStruct)
		So(ok, ShouldBeTrue)
		So(a, ShouldNotBeNil)
		So(a.Equal(m), ShouldBeTrue)

		a, ok = UnsafeShift1[MockAnonymousStruct](&m, o, -1, et).(*MockStruct)
		So(ok, ShouldBeTrue)
		So(a, ShouldNotBeNil)
		So(a.Equal(m), ShouldBeTrue)

		a = UnsafeShift2[MockAnonymousStruct, MockStruct](&m, o, -1)
		So(a, ShouldNotBeNil)
		So(a.Equal(m), ShouldBeTrue)

		a = UnsafeShift2[MockAnonymousStruct, MockStruct](&m, o, 0)
		So(a, ShouldNotBeNil)
		So(a.Equal(m), ShouldBeTrue)
	})

	Convey("use parent message, get anonymous message", t, func() {
		var m MockStruct
		m.PublicField = 3
		m.privateField = 4
		et := reflect.TypeFor[*MockAnonymousStruct]()
		o := getAnonymousOffset[*MockStruct]()

		a, ok := UnsafeShift0(m, o, 1, et).(*MockAnonymousStruct)
		So(ok, ShouldBeTrue)
		So(a, ShouldNotBeNil)
		So(a.Equal(m), ShouldBeTrue)

		a, ok = UnsafeShift1[MockStruct](&m, o, 1, et).(*MockAnonymousStruct)
		So(ok, ShouldBeTrue)
		So(a, ShouldNotBeNil)
		So(a.Equal(m), ShouldBeTrue)

		a = UnsafeShift2[MockStruct, MockAnonymousStruct](&m, o, 1)
		So(a, ShouldNotBeNil)
		So(a.Equal(m), ShouldBeTrue)
	})

	Convey("SliceCast should work ok", t, func() {
		var from []interface{}
		var out = SliceCast[interface{}, int](from)
		So(len(out), ShouldEqual, 0)
		So(out, ShouldBeNil)

		for i := 1; i < 100; i++ {
			from = append(from, i)
		}
		out = SliceCast[interface{}, int](from)
		So(len(out), ShouldEqual, len(from))
		So(func() {
			_ = SliceCast[interface{}, int16](from)
		}, ShouldPanic)

	})

	Convey("FuncName should work ok", t, func() {
		So(func() {
			_ = FuncName(nil)
		}, ShouldPanic)
		So(FuncName(FuncName), ShouldEqual, "github.com/sandwich-go/boost/misc.FuncName")
	})

	Convey("Zero should work ok", t, func() {
		So(Zero[map[string]string](), ShouldBeZeroValue)
		So(Zero[[]int](), ShouldBeZeroValue)
		So(Zero[string](), ShouldBeZeroValue)
		So(Zero[int](), ShouldBeZeroValue)
		So(Zero[MockStruct](), ShouldBeZeroValue)
		So(Zero[*MockStruct](), ShouldBeZeroValue)
	})
}
