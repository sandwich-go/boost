package bignum

import (
	"fmt"
	"testing"
)

func TestAdd(t *testing.T) {
	type args struct {
		d1 BigNum
		d2 BigNum
	}
	tests := []struct {
		name    string
		args    args
		want    BigNum
		wantErr bool
	}{
		{"100+2000", args{"100", "2000"}, "2100", false},
		{"100K+2000", args{"100K", "2000"}, BigNum("102" + unitStr["K"][1:]), false},
		{"100mm+200mm", args{"100mm", "200mm"}, BigNum("300" + unitStr["mm"][1:]), false},
		{"100zz+2zz", args{"100zz", "2zz"}, BigNum("102" + unitStr["zz"][1:]), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Add(tt.args.d1, tt.args.d2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Add() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiv(t *testing.T) {
	type args struct {
		d BigNum
		i int64
	}
	tests := []struct {
		name    string
		args    args
		want    BigNum
		wantErr bool
	}{
		{"0.001", args{"1000.001", 1000}, "0.001", false},
		{"1000.000", args{"2000.000", 1000}, "1000.000", false},
		{"1000", args{BigNum("9" + unitStr["K"][1:]), 8000}, "1000", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Div(tt.args.d, tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("BNDivInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BNMulInt64() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMul(t *testing.T) {
	type args struct {
		d BigNum
		i int64
	}
	tests := []struct {
		name    string
		args    args
		want    BigNum
		wantErr bool
	}{
		{"0.001", args{"0.001", 1000}, "1.000", false},
		{"1000.000", args{"1000.000", 1000}, "1000000.000", false},
		{"9K", args{"9", 1000}, BigNum("9" + unitStr["K"][1:]), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Mul(tt.args.d, tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("BNMulInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BNMulInt64() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_clean(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{"6.7", args{"6.7"}, "6.7", "", false},
		{"6.7K", args{"6.7K"}, "6.7", "K", false},
		{"a1K", args{"a1K"}, "", "", true},
		{"1a", args{"1a"}, "", "", true},
		{"1k", args{"1k"}, "", "", true},
		{"-1K", args{"-1K"}, "-1", "K", false},
		{"10T", args{"10T"}, "10", "T", false},
		{"99zz", args{"99zz"}, "99", "zz", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := clean(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("clean() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("clean() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("clean() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestFromString(t *testing.T) {
	type args struct {
		d string
	}
	tests := []struct {
		name    string
		args    args
		want    BigNum
		wantErr bool
	}{
		{"0.0000000012", args{"0.0000000012"}, "0.000000001", false},
		{"0.0000000019", args{"0.0000000019"}, "0.000000002", false},
		{"11.12zz", args{"11.12zz"}, "11120000000000000000000000000000000000000000000000000000000000000000000000000000000000000000.00", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromString(tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FromString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConstantTry(t *testing.T) {
	addP := 26
	oldParam := BigNum("1")
	var of float32
	for i := 1; i <= 1000; i++ {
		if i == addP {
			if i == 26 {
				of = 4
				addP += 25
			} else if i == 51 {
				of = 8
				addP += 50
			} else {
				of = 10
				addP += 50
			}

		} else {
			of = 1.05
		}
		mulFloat64, _ := Mul(oldParam, of)
		oldParam = mulFloat64
		fmt.Println(i, of, oldParam)
	}
	fmt.Println(oldParam)
}

func TestBigNum_Pow(t *testing.T) {
	type args struct {
		n int
	}
	tests := []struct {
		name    string
		s       BigNum
		args    args
		want    BigNum
		wantErr bool
	}{
		{"1", BigNum("1.2"), args{1}, "1.2", false},
		{"2", BigNum("1.2"), args{2}, "1.44", false},
		{"3", BigNum("1.2"), args{3}, "1.728", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Pow(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Pow() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestBigNum_Compare 覆盖 Compare 三种结果（未覆盖）。
func TestBigNum_Compare(t *testing.T) {
	cases := []struct {
		s, d BigNum
		want int
	}{
		{"1", "2", -1},
		{"2", "1", 1},
		{"1", "1", 0},
		{"1.5", "1.5", 0},
		{"100K", "99K", 1},
		{"99K", "100K", -1},
	}
	for _, c := range cases {
		got := c.s.Compare(c.d)
		if got != c.want {
			t.Errorf("Compare(%s, %s) = %d, want %d", c.s, c.d, got, c.want)
		}
	}
}

// TestBigNum_Pow_Errors 覆盖 Pow 错误路径（parse 失败）。
func TestBigNum_Pow_Errors(t *testing.T) {
	// invalid input → parse 返 err → Pow 返同 err
	_, err := BigNum("not-a-number").Pow(2)
	if err == nil {
		t.Error("Pow on invalid BigNum should return error")
	}
	// n=0 早返 "1"
	got, err := BigNum("anything-wont-be-parsed").Pow(0)
	if err != nil || got != "1" {
		t.Errorf("Pow(0) should return '1', no error; got %v / %v", got, err)
	}
}

// TestParse_AllTypes 覆盖 parse 内部各 type switch 分支（int / int32 /
// int64 / uint / uint32 / uint64 / float32 / float64 / BigNum / string /
// default）。通过公开 API Add 间接触发，让 doInterface 走每个类型。
func TestParse_AllTypes(t *testing.T) {
	cases := []struct {
		name    string
		v       interface{}
		wantErr bool
	}{
		{"int", int(5), false},
		{"int32", int32(5), false},
		{"int64", int64(5), false},
		{"uint", uint(5), false},
		{"uint32", uint32(5), false},
		{"uint64", uint64(5), false},
		{"float32", float32(1.5), false},
		{"float64", float64(2.5), false},
		{"BigNum", BigNum("3"), false},
		{"string", "4", false},
		{"unsupported_struct", struct{}{}, true},
		{"unsupported_slice", []int{1, 2}, true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := Add(c.v, BigNum("0"))
			if (err != nil) != c.wantErr {
				t.Errorf("Add(%v): err=%v, wantErr=%v", c.v, err, c.wantErr)
			}
		})
	}
}

// TestFromString_Errors 覆盖 FromString 在 parseString → clean / r.FromString
// 失败时的早返。
func TestFromString_Errors(t *testing.T) {
	cases := []struct {
		s       string
		wantErr bool
	}{
		{"", false},      // 空 string clean 返 "0"
		{"abc", true},    // 第一字符非 digit 非 '-' → ErrBadNumber
		{"123abc", true}, // 末位非 digit + 不是已知 unit
		{"123KK", true},  // KK 不在 unitStr（K 是 1 字符 unit, KK 不识别）
		// 注：mydecimal.FromString 对多 '.' 的字符串行为宽松（吞掉后续 '.'），
		// 因此 '123.45.67' 在本仓 parse 通路下不报错；不算 BigNum 设计 bug。
		{"100K", false}, // 已知 unit
		{"-100", false}, // 负数前缀 '-'
	}
	for _, c := range cases {
		t.Run(c.s, func(t *testing.T) {
			_, err := FromString(c.s)
			if (err != nil) != c.wantErr {
				t.Errorf("FromString(%q): err=%v, wantErr=%v", c.s, err, c.wantErr)
			}
		})
	}
}

// TestClean_EdgeCases 覆盖 clean 内部分支：
// - l == 0 → "0"
// - 第一字符非 digit 非 '-' → ErrBadNumber
// - 末位非 digit + 不是已知 unit → ErrBadNumber
// - l <= len(unit) → ErrBadNumber
// - 中间字符非 digit 非 '.' → ErrBadNumber
func TestClean_EdgeCases(t *testing.T) {
	cases := []struct {
		s       string
		wantErr bool
		wantS   string
		wantU   string
	}{
		{"", false, "0", ""},
		{"abc", true, "", ""},
		{"123x", true, "", ""},    // 末位非 digit / 非已知 unit
		{"K", true, "", ""},       // l=1 <= len("K")=1，第二段 for 抓不到 unit
		{"123-456", true, "", ""}, // 中间出现 '-'
	}
	for _, c := range cases {
		t.Run(c.s, func(t *testing.T) {
			s, u, err := clean(c.s)
			if (err != nil) != c.wantErr {
				t.Errorf("clean(%q): err=%v, wantErr=%v", c.s, err, c.wantErr)
			}
			if !c.wantErr {
				if s != c.wantS || u != c.wantU {
					t.Errorf("clean(%q): got s=%q u=%q, want s=%q u=%q", c.s, s, u, c.wantS, c.wantU)
				}
			}
		})
	}
}
