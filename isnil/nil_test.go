package isnil

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/modern-go/reflect2"
)

type foo struct{}

// diffCases 覆盖 Check 与 reflect2.IsNil 需要保持一致的全部 kind，
// 含 nil slice/map/chan/func —— 它们装箱后数据字非 nil，两者都返回 false。
var diffCases = []interface{}{
	nil,
	1,
	"nope",
	foo{},
	&foo{},
	(*foo)(nil),
	(interface{})(nil),
	fmt.Stringer(nil),
	([]int)(nil),
	[]int{},
	(map[string]int)(nil),
	map[string]int{},
	(chan int)(nil),
	(func())(nil),
	unsafe.Pointer(nil),
	struct{ a, b, c, d int64 }{1, 2, 3, 4},
}

// TestCheckMatchesReflect2 把 Check 的语义钉在 reflect2.IsNil 上：
// Check 用手写 eface 替换 reflect2.unpackEFace 以消除逃逸，行为必须逐 case 等价。
func TestCheckMatchesReflect2(t *testing.T) {
	for i, in := range diffCases {
		if got, want := Check(in), reflect2.IsNil(in); got != want {
			t.Errorf("case %d (%T): Check=%v, reflect2.IsNil=%v", i, in, got, want)
		}
	}
}

// TestCheckNoAlloc 防止 unsafe.Pointer(&any) 派生的指针再次流出函数体：
// 一旦流出，参数会被判为 leaking param，每次调用堆分配 16B eface。
func TestCheckNoAlloc(t *testing.T) {
	var pv = &foo{}
	var iv interface{} = (*foo)(nil)
	var sv = struct{ a, b, c, d int64 }{1, 2, 3, 4}
	for name, fn := range map[string]func(){
		"ptr":    func() { sink = Check(pv) },
		"iface":  func() { sink = Check(iv) },
		"int":    func() { sink = Check(counter) },
		"struct": func() { sink = Check(sv) },
	} {
		if n := testing.AllocsPerRun(100, fn); n != 0 {
			t.Errorf("Check(%s): got %v allocs/op, want 0", name, n)
		}
	}
}

var (
	sink    bool
	counter = 1 << 20 // 大于 256，避免命中 runtime.staticuint64s 免装箱路径
)

func TestIsNil(t *testing.T) {
	cases := []struct {
		in  interface{}
		exp bool
	}{
		{1, false},
		{"nope", false},
		{foo{}, false},
		{&foo{}, false},
		{nil, true},
		{(*foo)(nil), true},
		{(interface{})(nil), true},
		{fmt.Stringer(nil), true},
	}

	for _, tcase := range cases {
		if out := Check(tcase.in); out != tcase.exp {
			if tcase.exp {
				t.Errorf("Expected %++v to be nil", tcase.in)
			} else {
				t.Errorf("Expected %++v to not be nil", tcase.in)
			}
		}
	}
}

func BenchmarkEqNilBasic(b *testing.B) {
	var v *int
	for i := 0; i < b.N; i++ {
		_ = v != nil
	}
}

func BenchmarkEqNilInterface(b *testing.B) {
	var v interface{}
	v = (*foo)(nil)
	for i := 0; i < b.N; i++ {
		_ = v != nil
	}
}

func BenchmarkIsNilBasic(b *testing.B) {
	var v *int
	for i := 0; i < b.N; i++ {
		Check(v)
	}
}

func BenchmarkIsNilInterface(b *testing.B) {
	var v interface{}
	v = (*foo)(nil)
	for i := 0; i < b.N; i++ {
		Check(v)
	}
}

func BenchmarkIsNilNil(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Check(nil)
	}
}
