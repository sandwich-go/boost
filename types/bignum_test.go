package types

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
