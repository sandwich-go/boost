package types

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/sandwich-go/boost/types/internal"
)

// ErrBadNumber if string not BigNum，return ErrBadNumber
var ErrBadNumber = errors.New("bad number")

func init() {
	t := "1"
	for i := 0; i < len(kUnit); i++ {
		t = t + "000"
		unitStr[kUnit[i]] = t
	}
}

// BigNum big decimal
type BigNum string

// Compare compares one decimal to another, returns -1/0/1.
// s<d: -1; s==d: 0; s>d: 1
func (s BigNum) Compare(d BigNum) int {
	s1, _ := parse(s)
	d1, _ := parse(d)
	return s1.Compare(d1)
}

// Pow returns s**n, the base-s exponential of n.
func (s BigNum) Pow(n int) (BigNum, error) {
	if n == 0 {
		return "1", nil
	}
	b, err := parse(s)
	if err != nil {
		return s, err
	}
	r := new(internal.MyDecimal)
	r.Copy(b)
	for i := 1; i < n; i++ {
		var t internal.MyDecimal
		if e := internal.DecimalMul(r, b, &t); e != nil {
			return "", err
		}
		r.Copy(&t)
	}
	return BigNum(r.String()), nil
}

// FromString returns big decimal.
func FromString(s string) (BigNum, error) {
	d, err := parse(s)
	if err != nil {
		return "", err
	}
	return BigNum(d.String()), nil
}

// Add returns d1 + d2.
func Add(d1, d2 interface{}) (BigNum, error) {
	return doInterface(d1, d2, addFunc)
}

// Div returns d1 - d2.
func Div(d1, d2 interface{}) (BigNum, error) {
	return doInterface(d1, d2, subFunc)
}

// Mul returns d1 * d2.
func Mul(d1, d2 interface{}) (BigNum, error) {
	return doInterface(d1, d2, mulFunc)
}

func doInterface(s, t interface{}, f calcFun) (BigNum, error) {
	id1, err := parse(s)
	if err != nil {
		return "", err
	}
	id2, err := parse(t)
	if err != nil {
		return "", err
	}
	var result internal.MyDecimal
	err = f(id1, id2, &result)
	if err != nil {
		return "", err
	}
	return BigNum(result.String()), nil
}

func parse(i interface{}) (*internal.MyDecimal, error) {
	r := new(internal.MyDecimal)
	var err error
	switch v := i.(type) {
	case int:
		r.FromInt(int64(v))
	case int32:
		r.FromInt(int64(v))
	case int64:
		r.FromInt(v)
	case uint:
		r.FromUint(uint64(v))
	case uint32:
		r.FromUint(uint64(v))
	case uint64:
		r.FromUint(v)
	case float32:
		r, err = parseString(fmt.Sprintf("%f", v))
	case float64:
		var t internal.MyDecimal
		if err = t.FromFloat64(v); err == nil {
			r, err = fixFrac(&t)
		}
	case BigNum:
		r, err = parseString(string(v))
	case string:
		r, err = parseString(v)
	default:
		err = errors.New(fmt.Sprintf("not support type:%s of:%v", reflect.TypeOf(i), i))
	}
	return r, err
}

func parseString(s string) (*internal.MyDecimal, error) {
	s, u, err := clean(s)
	if err != nil {
		return nil, err
	}
	r := new(internal.MyDecimal)
	err = r.FromString([]byte(s))
	if err != nil {
		return nil, err
	}
	if uu, ok := unitStr[u]; ok {
		ud := new(internal.MyDecimal)
		_ = ud.FromString([]byte(uu))
		r2 := new(internal.MyDecimal)
		err2 := internal.DecimalMul(r, ud, r2)
		if err2 != nil {
			return nil, err2
		}
		return r2, nil
	}
	return fixFrac(r)
}

func fixFrac(i *internal.MyDecimal) (*internal.MyDecimal, error) {
	r := new(internal.MyDecimal)
	err := i.Round(r, internal.MaxFraction)
	return r, err
}

func clean(s string) (string, string, error) {
	l := len(s)
	if l == 0 {
		return "0", "", nil
	}
	if !isDigit(s[0]) && s[0] != '-' {
		return "", "", ErrBadNumber
	}
	checkNumber := func(i int) error {
		for j := 1; j < i; j++ {
			if !isDigit(s[j]) && s[j] != '.' {
				return ErrBadNumber
			}
		}
		return nil
	}
	if isDigit(s[l-1]) {
		if e := checkNumber(l); e != nil {
			return "", "", e
		}
		return s, "", nil
	}
	for i := 0; i < len(kUnit); i++ {
		u := kUnit[i]
		if l <= len(u) {
			return "", "", ErrBadNumber
		}
		if s[l-len(u):] == u {
			// 末尾去除单位之后，前面必须全是数字（0位已经判断过跳过）
			if e := checkNumber(l - len(u)); e != nil {
				return "", "", e
			}
			return s[0 : l-len(u)], u, nil
		}
	}
	return "", "", ErrBadNumber
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

type calcFun func(from1 *internal.MyDecimal, from2 *internal.MyDecimal, to *internal.MyDecimal) error

var (
	addFunc = internal.DecimalAdd
	subFunc = internal.DecimalSub
	mulFunc = internal.DecimalMul

	unitStr = map[string]string{}
	kUnit   = []string{"K", "M", "B", "T", "aa", "bb", "cc", "dd", "ee", "ff", "gg", "hh", "ii", "jj",
		"kk", "ll", "mm", "nn", "oo", "pp", "qq", "rr", "ss", "tt", "uu", "vv", "ww", "xx", "yy", "zz"}
)
