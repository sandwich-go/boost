package xcopy

import (
	proto1 "github.com/golang/protobuf/proto"
	"github.com/sandwich-go/boost/xrand"
	"github.com/sandwich-go/boost/z"
	. "github.com/smartystreets/goconvey/convey"
	proto2 "google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"reflect"
	"strconv"
	"testing"
)

type withDeepCopy struct {
	stringF string
	intF    uint32
	cloneF  int
}

func randWithDeepCopy() *withDeepCopy {
	return &withDeepCopy{intF: z.FastRand(), stringF: xrand.String(20)}
}

func (w withDeepCopy) getCloneFlag() int { return w.cloneF }
func (w *withDeepCopy) DeepCopy() interface{} {
	if w == nil {
		return nil
	}
	return &withDeepCopy{stringF: w.stringF, intF: w.intF, cloneF: cloneFlagDeep}
}

func (w *withDeepCopy) equal(w1 *withDeepCopy) bool {
	if w == nil && w1 == nil {
		return true
	}
	if (w == nil && w1 != nil) || (w1 == nil && w != nil) {
		return false
	}
	return w.intF == w1.intF && w.stringF == w1.stringF
}

func copyMap(in map[string]string) (out map[string]string) {
	if in != nil {
		out = make(map[string]string)
		for k, v := range in {
			out[k] = v
		}
	}
	return
}

func equalMap(a, b map[string]string) bool {
	if a == nil && b == nil {
		return true
	}
	if (a == nil && b != nil) || (b == nil && a != nil) || len(a) != len(b) {
		return false
	}
	for k, v := range a {
		v1, ok := b[k]
		if !ok || v != v1 {
			return false
		}
	}
	return true
}

func randMap() map[string]string {
	n := z.FastRandUint32n(101)
	out := make(map[string]string)
	var i uint32
	for i = 0; i < n; i++ {
		out[strconv.FormatUint(uint64(i), 10)] = xrand.String(20)
	}
	return out
}

func copySlice(in []string) (out []string) {
	if in != nil {
		out = make([]string, len(in))
		copy(out, in)
	}
	return
}

func equalSlice(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}
	if (a == nil && b != nil) || (b == nil && a != nil) || len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if v != b[k] {
			return false
		}
	}
	return true
}

func randSlice() []string {
	n := z.FastRandUint32n(101)
	out := make([]string, n)
	for k := range out {
		out[k] = xrand.String(20)
	}
	return out
}

const (
	cloneFlag = iota
	cloneFlag1
	cloneFlag2
	cloneFlagDeep
)

type withClone1 struct {
	stringF  string
	intF     uint32
	mapF     map[string]string
	sliceF   []string
	pointerF *withDeepCopy
	cloneF   int
}

func randWithClone1() *withClone1 {
	return &withClone1{
		intF: z.FastRand(), stringF: xrand.String(20),
		mapF: randMap(), sliceF: randSlice(),
		pointerF: randWithDeepCopy(),
	}
}

func (w withClone1) getCloneFlag() int { return w.cloneF }
func (*withClone1) Reset()             {}
func (*withClone1) String() string     { return "" }
func (*withClone1) ProtoMessage()      {}
func (w *withClone1) Clone() proto1.Message {
	c := &withClone1{
		stringF: w.stringF,
		intF:    w.intF,
		mapF:    copyMap(w.mapF),
		sliceF:  copySlice(w.sliceF),
		cloneF:  cloneFlag1,
	}
	if w.pointerF != nil {
		c.pointerF = w.pointerF.DeepCopy().(*withDeepCopy)
	}
	return c
}
func (w *withClone1) equal(w1 *withClone1) bool {
	if w == nil && w1 == nil {
		return true
	}
	if (w == nil && w1 != nil) || (w1 == nil && w != nil) {
		return false
	}
	if w.stringF != w1.stringF || w.intF != w1.intF ||
		!equalMap(w.mapF, w1.mapF) || !equalSlice(w.sliceF, w1.sliceF) || !w.pointerF.equal(w1.pointerF) {
		return false
	}
	return true
}

type withClone2 struct {
	stringF  string
	intF     uint32
	mapF     map[string]string
	sliceF   []string
	pointerF *withDeepCopy
	cloneF   int
}

func randWithClone2() *withClone2 {
	return &withClone2{
		intF: z.FastRand(), stringF: xrand.String(20),
		mapF: randMap(), sliceF: randSlice(),
		pointerF: randWithDeepCopy(),
	}
}

func (w withClone2) getCloneFlag() int                { return w.cloneF }
func (withClone2) ProtoReflect() protoreflect.Message { return nil }
func (w *withClone2) Clone() proto2.Message {
	c := &withClone2{
		stringF: w.stringF,
		intF:    w.intF,
		mapF:    copyMap(w.mapF),
		sliceF:  copySlice(w.sliceF),
		cloneF:  cloneFlag2,
	}
	if w.pointerF != nil {
		c.pointerF = w.pointerF.DeepCopy().(*withDeepCopy)
	}
	return c
}
func (w *withClone2) equal(w1 *withClone2) bool {
	if w == nil && w1 == nil {
		return true
	}
	if (w == nil && w1 != nil) || (w1 == nil && w != nil) {
		return false
	}
	if w.stringF != w1.stringF || w.intF != w1.intF ||
		!equalMap(w.mapF, w1.mapF) || !equalSlice(w.sliceF, w1.sliceF) || !w.pointerF.equal(w1.pointerF) {
		return false
	}
	return true
}

func equal(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if (a == nil && b != nil) || (b == nil && a != nil) {
		return false
	}
	if reflect.ValueOf(a).Type().String() != reflect.ValueOf(b).Type().String() {
		return false
	}
	switch aa := a.(type) {
	case *withClone1:
		return aa.equal(b.(*withClone1))
	case *withClone2:
		return aa.equal(b.(*withClone2))
	case *withDeepCopy:
		return aa.equal(b.(*withDeepCopy))
	case *common:
		return aa.equal(b.(*common))
	}
	return false
}

type common struct {
	I uint32
	B *subMap
	C *subSlice
	D subStruct
	E getter
	F *withDeepCopy
	G *withClone1
	H *withClone2
}

func randCommon() *common {
	return &common{
		I: z.FastRand(), B: randSubMap(), C: randSubSlice(),
		D: randSubStruct(), E: randSubInterface(), F: randWithDeepCopy(),
		G: randWithClone1(), H: randWithClone2(),
	}
}

func (c *common) equal(c1 *common) bool {
	if c == nil && c1 == nil {
		return true
	}
	if (c == nil && c1 != nil) || (c1 == nil && c != nil) {
		return false
	}
	if c.I != c1.I || !c.B.equal(c1.B) || !c.C.equal(c1.C) || !c.E.equal(c1.E) || !c.F.equal(c1.F) || !c.G.equal(c1.G) || !c.H.equal(c1.H) {
		return false
	}
	return true
}

type getter interface {
	get() uint32
	equal(getter) bool
}

type subInterface struct {
	I uint32
}

func randSubInterface() *subInterface {
	return &subInterface{I: z.FastRand()}
}

func (s subInterface) get() uint32 { return s.I }

func (s *subInterface) equal(s1 getter) bool {
	if s == nil && s1 == nil {
		return true
	}
	if (s == nil && s1 != nil) || (s1 == nil && s != nil) {
		return false
	}
	return s.get() == s1.get()
}

type subStruct struct {
	I uint32
}

func randSubStruct() subStruct {
	return subStruct{I: z.FastRand()}
}

func (s subStruct) equal(s1 subStruct) bool {
	return s.I == s1.I
}

type subSlice struct {
	B []string
}

func randSubSlice() *subSlice {
	return &subSlice{B: randSlice()}
}

func (s *subSlice) equal(s1 *subSlice) bool {
	if s == nil && s1 == nil {
		return true
	}
	if (s == nil && s1 != nil) || (s1 == nil && s != nil) {
		return false
	}
	return equalSlice(s.B, s1.B)
}

type subMap struct {
	B map[string]string
}

func (s *subMap) equal(s1 *subMap) bool {
	if s == nil && s1 == nil {
		return true
	}
	if (s == nil && s1 != nil) || (s1 == nil && s != nil) {
		return false
	}
	return equalMap(s.B, s1.B)
}

func randSubMap() *subMap {
	return &subMap{B: randMap()}
}

type sub struct {
	i uint32
}

func randSub() *sub {
	return &sub{i: z.FastRand()}
}

func (s *sub) equal(s1 *sub) bool {
	if s == nil && s1 == nil {
		return true
	}
	if (s == nil && s1 != nil) || (s1 == nil && s != nil) {
		return false
	}
	return s.i == s1.i
}

func TestDeepCopy(t *testing.T) {
	Convey("deep copy", t, func() {
		for _, test := range []struct {
			src       interface{}
			cloneFlag int
			notEqual  bool
		}{
			{src: nil, cloneFlag: cloneFlag},
			{src: randWithClone1(), cloneFlag: cloneFlag1},
			{src: randWithClone2(), cloneFlag: cloneFlag2},
			{src: randWithDeepCopy(), cloneFlag: cloneFlagDeep},
			{src: randCommon()},
			{src: randSub(), notEqual: true},
		} {
			dest := DeepCopy(test.src)
			if test.notEqual {
				So(equal(test.src, dest), ShouldBeFalse)
			} else {
				So(equal(test.src, dest), ShouldBeTrue)
			}
			if v, ok := dest.(interface{ getCloneFlag() int }); ok {
				So(v.getCloneFlag(), ShouldEqual, test.cloneFlag)
			}
		}
	})
}
