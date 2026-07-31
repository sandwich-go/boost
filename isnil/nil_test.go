package isnil

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/modern-go/reflect2"
)

type foo struct{}

// truthTable 钉住 Check 的绝对返回值，覆盖 pointer-shaped 与非 pointer-shaped 两类。
// 差分对拍（TestCheckMatchesReflect2）只保证两个实现一致，抓不到「对 nil map/chan/func
// 到底返回什么」这类理解错误，故绝对值必须单独断言。
var truthTable = []struct {
	name string
	in   interface{}
	want bool
}{
	// pointer-shaped：值直接存在数据字，nil 值进去数据字就是 0
	{"nil", nil, true},
	{"(*foo)(nil)", (*foo)(nil), true},
	{"(map[string]int)(nil)", (map[string]int)(nil), true},
	{"(chan int)(nil)", (chan int)(nil), true},
	{"(func())(nil)", (func())(nil), true},
	{"unsafe.Pointer(nil)", unsafe.Pointer(nil), true},
	{"(interface{})(nil)", (interface{})(nil), true},
	{"fmt.Stringer(nil)", fmt.Stringer(nil), true},
	// slice 是三字结构，nil slice 装箱后数据字指向 zeroVal，非 nil
	{"([]int)(nil)", ([]int)(nil), false},
	// 非 nil 的 pointer-shaped 与其他类型
	{"&foo{}", &foo{}, false},
	{"map[string]int{}", map[string]int{}, false},
	{"make(chan int)", make(chan int), false},
	{"[]int{}", []int{}, false},
	{"1", 1, false},
	{"\"nope\"", "nope", false},
	{"foo{}", foo{}, false},
	{"big struct", struct{ a, b, c, d int64 }{1, 2, 3, 4}, false},
}

func TestCheckTruthTable(t *testing.T) {
	for _, tc := range truthTable {
		if got := Check(tc.in); got != tc.want {
			t.Errorf("Check(%s) = %v, want %v", tc.name, got, tc.want)
		}
	}
	// func(){} 不能进 truthTable：非 nil func 值不可比较，放进结构体切片没问题，
	// 但保持与其他 case 同构更清晰，单独断言。
	if got := Check(func() {}); got {
		t.Errorf("Check(func(){}) = true, want false")
	}
}

// TestCheckMatchesReflect2 把 Check 的语义钉在 reflect2.IsNil 上：
// Check 用手写 eface 替换 reflect2.unpackEFace 以消除逃逸，行为必须逐 case 等价。
func TestCheckMatchesReflect2(t *testing.T) {
	for _, tc := range truthTable {
		if got, want := Check(tc.in), reflect2.IsNil(tc.in); got != want {
			t.Errorf("%s: Check=%v, reflect2.IsNil=%v", tc.name, got, want)
		}
	}
	if got, want := Check(func() {}), reflect2.IsNil(func() {}); got != want {
		t.Errorf("func(){}: Check=%v, reflect2.IsNil=%v", got, want)
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
