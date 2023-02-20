// Copyright © 2012-2013 Lawrence E. Bakst. All rights reserved.

package hrff

import (
	"flag"
	"fmt"
	"testing"
)

var _i1 = Int64{V: 2}
var i1 int64

const ci1 = 1 * 1024 * 1024 * 1024
const ci2 = -2

var _i2 = Int64{V: 3}
var i2 int64

var _f1 = Float64{V: 4.5}
var f1 float64
var _f2 = Float64{V: 5.5}
var f2 float64

func TestHRFF(t *testing.T) {

	_i1.Set("1k")
	i1 = int64(_i1.V)
	if i1 != 1000 {
		t.Errorf("wanted %d got %d\n", ci1, 1000)
	}

	_i1.Set("1M")
	i1 = int64(_i1.V)
	if i1 != 1000000 {
		t.Errorf("wanted %d got %d\n", ci1, 1000000)
	}

	_i1.Set("1000k")
	i1 = int64(_i1.V)
	if i1 != 1000000 {
		t.Errorf("wanted %d got %d\n", ci1, 1000000)
	}

	_i1.Set("1Gi")
	i1 = int64(_i1.V)
	if i1 != ci1 {
		t.Errorf("wanted %d got %d\n", ci1, i1)
	}
	_f1.Set("123.k")
	f1 = 123. * 1000
	if float64(_f1.V) != f1 {
		t.Errorf("wanted %f got %f\n", f1, _f1)
	}
	_f1.Set("-1.M")
	f1 = -1 * 1000000
	if float64(_f1.V) != f1 {
		t.Errorf("wanted %f got %f\n", f1, _f1)
	}
	_f1.Set(".001m")
	f1 = .000001
	if float64(_f1.V) != f1 {
		t.Errorf("wanted %f got %f\n", f1, _f1)
	}

}

func Example001() {
	var i1 = Int64{2000000, "bps"}
	var i2 = Int64{V: 3000}
	var f1 = Float64{-7020000000.1, "B"}
	var f2 = Float64{.000001, "s"}

	var imm1 = Int64{-40000, ""}
	var imm2 = Int64{1 * 1024 * 1024 * 1024, ""}

	flag.Var(&i1, "10", "i1")
	flag.Var(&i2, "11", "i2")
	flag.Var(&f1, "12.1", "f1")
	flag.Var(&f2, "12.2", "f2")

	flag.Parse()
	fmt.Printf("i1=%d, i2=%d, f1=%f, f2=%f, i1=%h, i2=%h, f1=%0.4h, f2=%h, imm1=%h, imm2=%H\n",
		i1.V, i2.V, f1.V, f2.V, i1, i2, f1, f2, imm1, imm2)
	fmt.Printf("%10.3h\n", Int64{V: 0, U: "foobars"})
	// Output: i1=2000000, i2=3000, f1=-7020000000.100000, f2=0.000001, i1=2 Mbps, i2=3 k, f1=-7.0200 GB, f2=1 µs, imm1=-40 k, imm2=1 Gi
	//          0 foobars
}

func Example002() {
	var size = Int64{3 * 1024 * 1024 * 1024, "B"}
	var speed = Float64{2100000, "bps"}

	fmt.Printf("size=%H, speed=%0.2h\n", size, speed)
	// Output: size=3 GiB, speed=2.10 Mbps
}
func Example003() {
	var v = Int64{0, "B"}

	fmt.Printf("v=%h\n", v)
	// Output: v=0 B
}
func Example004() {
	var v = Int64{1, "B"}

	fmt.Printf("v=%h\n", v)
	// Output: v=1 B
}

func Example005() {
	var v = Float64{0, "B"}

	fmt.Printf("v=%h\n", v)
	// Output: v=0 B
}

func Example006() {
	var v = Float64{1, "B"}

	fmt.Printf("v=%h\n", v)
	// Output: v=1 B
}

func Example007() {
	var v = Float64{1024 * 1024 * 1024, "B"}

	fmt.Printf("v=%h\n", v)
	// Output: v=1 GB
}

func Example008() {
	var v = Int64{1000, "B"}
	fmt.Printf("v=%D\n", v)
	// Output: v=1,000
}

func Example009() {
	var v = Int{1000, "B"}
	fmt.Printf("v=%D\n", v)
	// Output: v=1,000
}

func Example010() {
	var v = Int{-100, "B"}
	fmt.Printf("v=%D\n", v)
	// Output: v=-100
}

func Example011() {
	var v = Int{-1000, "B"}
	fmt.Printf("v=%D\n", v)
	// Output: v=-1,000
}

func Example012() {
	var v = Int{1234567, "B"}
	fmt.Printf("v=%D\n", v)
	// Output: v=1,234,567
}

func Example013() {
	var v = Int{-1234567, "B"}
	fmt.Printf("v=%D\n", v)
	// Output: v=-1,234,567
}

func Example015a() {
	var v = Float64{0, "B"}
	fmt.Printf("v=%h\n", v)
	// Output: v=0 B
}

func Example014() {
	fmt.Printf("%h\n", Int64{V: 0, U: "foobars"})
	// Output: 0 foobars
}

func Example015() {
	fmt.Printf("%h\n", Float64{V: 0, U: "foobars"})
	// Output: 0 foobars
}

func Example016() {
	fmt.Printf("%h\n", Int64{V: 11, U: "foobars"})
	// Output: 11 foobars
}

func Example017() {
	fmt.Printf("%h\n", Float64{V: 11, U: "foobars"})
	// Output: 11 foobars
}

func Example018() {
	fmt.Printf("%h\n", Int64{V: 999, U: "foobars"})
	// Output: 999 foobars
}

func Example019() {
	fmt.Printf("%h\n", Float64{V: 999, U: "foobars"})
	// Output: 999 foobars
}

func Example020() {
	fmt.Printf("%h\n", Int64{V: 1000, U: "foobars"})
	// Output: 1 kfoobars
}

func Example021() {
	fmt.Printf("%h\n", Float64{V: 1000, U: "foobars"})
	// Output: 1 kfoobars
}

func Example022() {
	fmt.Printf("%D\n", Int64{V: 1000000, U: "foobars"})
	// Output: 1,000,000
}
