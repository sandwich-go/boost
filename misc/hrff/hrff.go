// Copyright © 2012-2019 Lawrence E. Bakst. All rights reserved.

// Package hrff (Human Readbale Flags and Formatting)
// Allows command line arguments like % dd bs=1Mi
// Provides SI unit formatting via %h and %H format characters
// Defines two news types, Int64 and Float64 which provide methods for flags to accept these kind of args
package hrff // import "leb.io/hrff"

import (
	"fmt"
	"strconv"
)

// SIsufixes is public so you can add a prefix if you want to
var SIsufixes = map[string]float64{

	"geop":   10000000000000000000000000000000, // geop 10^30
	"bronto": 10000000000000000000000000000,    // bronto 10^27
	"Y":      1000000000000000000000000,        // yota
	"Z":      1000000000000000000000,           // zetta
	"E":      1000000000000000000,              // exa
	"P":      1000000000000000,                 // peta
	"T":      1000000000000,                    // tera
	"G":      1000000000,                       // giga
	"M":      1000000,                          // mega
	//"my": 10000,				     // prefix myria- (my-) was formerly used for 10^4 but now depricated
	"k":  1000,                      // kilo, chilo in Italian
	"h":  100,                       // hecto, etto in Italian
	"da": 10,                        // SI: deca, NIST: deka
	"":   1,                         // not real dummy stopper
	"d":  .1,                        // deci
	"c":  .01,                       // centi
	"m":  .001,                      // milli
	"µ":  .000001,                   // micro (unicode char see below)
	"n":  .000000001,                // nano
	"p":  .00000000001,              // pico
	"f":  .000000000000001,          // femto
	"a":  .000000000000000001,       // atto
	"z":  .000000000000000000001,    // zepto
	"y":  .000000000000000000000001, // yocto

	"u": .000001, // micro (with u)

	"Ki": 1024,                                                  // kibi
	"Mi": 1024 * 1024,                                           // mebi
	"Gi": 1024 * 1024 * 1024,                                    // gibi
	"Ti": 1024 * 1024 * 1024 * 1024,                             // tebi
	"Pi": 1024 * 1024 * 1024 * 1024 * 1024,                      // pebi
	"Ei": 1024 * 1024 * 1024 * 1024 * 1024 * 1024,               // exbi
	"Zi": 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024,        // zebi
	"Yi": 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024, // yobi
}

var order = []string{"geop", "bronto", "Y", "Z", "E", "P", "T", "G", "M", "k", "h", "da", "", "d", "c", "m", "µ", "n", "p", "f", "a", "z", "y"}
var order2 = []string{"Yi", "Zi", "Ei", "Pi", "Ti", "Gi", "Mi", "Ki", "", "d", "c", "m", "µ", "n", "p", "f", "a", "z", "y"}
var skips = map[string]bool{"h": true, "da": true, "d": true, "c": true} // The sufixes h, da, d, c aren't much used for scienttific work, skip them

// Classic is what life was like before SI
func Classic() {
	SIsufixes["K"] = SIsufixes["Ki"]
	SIsufixes["M"] = SIsufixes["Mi"]
	SIsufixes["G"] = SIsufixes["Gi"]
	SIsufixes["T"] = SIsufixes["Ti"]
	SIsufixes["P"] = SIsufixes["Pi"]
	SIsufixes["E"] = SIsufixes["Ei"]
	SIsufixes["Z"] = SIsufixes["Zi"]
	SIsufixes["Y"] = SIsufixes["Yi"]
}

// RemoveNominal removes unoffice suffixes
func RemoveNominal() {
	delete(SIsufixes, "h")
	delete(SIsufixes, "da")
	delete(SIsufixes, "d")
	delete(SIsufixes, "c")
}

// UseHella allows use of H over bronto
func UseHella() {
	delete(SIsufixes, "bronto")
	SIsufixes["H"] = 10000000000000000000000000000 // hella (one for the team)
}

// Int is a version of int with a unit
type Int struct {
	V int
	U string
}

// Float64 is a version of float64 with a unit
type Float64 struct {
	V float64
	U string
}

// Int64 is a version of int64 with a unit
type Int64 struct {
	V int64
	U string
}

// AddSkip adds a skip
func AddSkip(sip string, b bool) {
	skips[sip] = b
}

// NoSkips gets rid of all the skips
func NoSkips() {
	for k := range skips {
		skips[k] = false
	}
}

// thanks to my mentor
func knot(c rune, chars string) bool {
	for _, v := range chars {
		if c == v {
			return false
		}
	}
	return true
}

func getPrefix(s string) (float64, int, bool) {
	var m float64 = 1
	var o int = 0

	//	fmt.Printf("getPrefix: s=%q\n", s)
	_, ok := SIsufixes["xxx"] // better way?
	l := len(s)
	if l > 1 {
		if knot(rune(s[l-1]), "0123456789.") {
			if l > 2 {
				if knot(rune(s[l-2]), "0123456789.+-eE") {
					o = 2
				} else {
					o = 1
				}
			} else {
				o = 1
			}
		}
		m, ok = SIsufixes[s[l-o:]]
		//		fmt.Printf("getPrefix: %q, m=%f, l=%d, o=%d, ok=%v\n", s[l-o:], m, l, o, ok)
	}
	return m, l - o, ok
}

// print integer format
func pif(val int64, units string, w, p int, okw, okp bool, order []string) string {
	var sip string

	//fmt.Printf("pif: %d\n", val)
	sgn := ""
	if val < 0 {
		sgn = "-"
		val = -val
	}
	if val == 0 {
		p = 1
		okp = true
	}

	//fs := fmt.Sprintf("%%s%%%d.%dd %%s%%s", w, p)
	fs := ""
	switch {
	case okw == false && okp == false:
		fs = fmt.Sprintf("%%s%%%d.%dd %%s%%s", 0, 1)
	case okw == false && okp == true:
		fs = fmt.Sprintf("%%s%%.%dd %%s%%s", p)
	case okw == true && okp == false:
		fs = fmt.Sprintf("%%s%%%d.d %%s%%s", w)
	case okw == true && okp == true:
		fs = fmt.Sprintf("%%s%%%d.%dd %%s%%s", w, p)
	}

	//fmt.Printf("sgn=%q, fs=%q\n", sgn, fs)

	for _, sip = range order {
		//		fmt.Printf("Format: try %q, ", sip)
		if skips[sip] {
			continue
		}
		if (SIsufixes[sip] <= float64(val)) || (sip == "" && val == 0) {
			break
		}
	}
	//fmt.Printf("pif: sip=%q\n", sip)
	val = val / int64(SIsufixes[sip])
	//fmt.Printf("pif: val=%d\n", val)
	str := fmt.Sprintf(fs, sgn, val, sip, units)
	if str[len(str)-1] == ' ' {
		str = str[:len(str)-1]
	}
	return str
}

// print floating format
func pff(val float64, units string, w, p int, okw, okp bool, order []string) string {
	var sip string

	// fmt.Printf("pff: %f\n", val)
	sgn := ""
	if val < 0 {
		sgn = "-"
		val = -val
	}
	if val == 0 {
		w = 1
	}

	fs := ""
	switch {
	case okw == false && okp == false:
		fs = fmt.Sprintf("%%s%%%d.%df %%s%%s", 0, 0)
	case okw == false && okp == true:
		fs = fmt.Sprintf("%%s%%.%df %%s%%s", p)
	case okw == true && okp == false:
		fs = fmt.Sprintf("%%s%%%d.f %%s%%s", w)
	case okw == true && okp == true:
		fs = fmt.Sprintf("%%s%%%d.%df %%s%%s", w, p)
	}
	//fmt.Printf("sgn=%q, fs=%q\n", sgn, fs)

	for _, sip = range order {
		if skips[sip] {
			continue
		}
		//fmt.Printf("pff: %q, %f <= %f\n", sip, SIsufixes[sip], val)
		if SIsufixes[sip] == 1 {
			if val == 0.0 || val == 1.0 {
				break
			}
			//continue
		}
		if SIsufixes[sip] <= val {
			break
		}
	}
	//fmt.Printf("pff: val=%f, sip=%q\n", val, sip)
	val = val / SIsufixes[sip]
	str := fmt.Sprintf(fs, sgn, val, sip, units)
	if str[len(str)-1] == ' ' {
		str = str[:len(str)-1]
	}
	return str
}

// called to parse format descriptor
func i(v *Int64, s fmt.State, c rune) {
	var val = int64(v.V)
	var str string
	var w, p int
	var okw, okp bool

	// not checking is OK because 0 is the default behavior
	w, okw = s.Width()
	p, okp = s.Precision()
	//fmt.Printf("i: w=%d, okw=%v, p=%d, okp=%v\n", w, okw, p, okp)
	//mi, pl, sh, sp, ze := s.Flag('-'), s.Flag('+'), s.Flag('#'), s.Flag(' '), s.Flag('0')

	switch c {
	case 'h':
		str = pif(val, v.U, w, p, okw, okp, order)
	case 'H':
		str = pif(val, v.U, w, p, okw, okp, order2)
	case 'd':
		str = fmt.Sprintf("%d", val)
	case 'D':
		fs := fmt.Sprintf("%%%d.%dd", w, p)
		tmp := fmt.Sprintf(fs, val)
		str = ""
		for k := range tmp {
			c := string(tmp[len(tmp)-k-1])
			if c < `0` || c > `9` {
				str = tmp[0:len(tmp)-k] + str
				break
			}
			if k > 0 && k%3 == 0 {
				str = "," + str
			}
			str = c + str
		}
	case 'v':
		str = fmt.Sprintf("%v", val)
	default:
		// fmt.Printf("default\n")
		str = fmt.Sprintf("%d", val)
	}
	b := []byte(str)
	s.Write(b)
}

// called to parse format descriptor
func f(v *Float64, s fmt.State, c rune) {
	var val = float64(v.V)
	var str string
	var w, p int
	var okw, okp bool

	w, okw = s.Width()
	p, okp = s.Precision()

	switch c {
	case 'h':
		str = pff(val, v.U, w, p, okw, okp, order)
	case 'H':
		str = pff(val, v.U, w, p, okw, okp, order2)
	case 'd':
		str = fmt.Sprintf("%d", int(val))
	case 'v':
		str = fmt.Sprintf("%v", val)
	default:
		// fmt.Printf("default\n")
		str = fmt.Sprintf("%d", int(val))
	}
	b := []byte(str)
	s.Write(b)
}

// FIX FIX FIX check ok or err? if no prefix we must convert anyway not err
func (r *Int64) Set(s string) error {

	m, l, _ := getPrefix(s)
	v, err := strconv.ParseInt(s[:l], 10, 64)
	if err != nil {
		return err
	}
	// fmt.Printf("Set: v=%d, m=%f, v*m=%v\n", v, m, v*int64(m))
	r.V = int64(v * int64(m))
	return err
}

func (r *Int) Set(s string) error {
	m, l, _ := getPrefix(s)
	v, err := strconv.ParseInt(s[:l], 10, 64)
	if err != nil {
		return err
	}
	// fmt.Printf("Set: v=%d, m=%f, v*m=%v\n", v, m, v*int64(m))
	r.V = int(v * int64(m))
	return err
}

func (r *Float64) Set(s string) error {

	m, l, ok := getPrefix(s)
	v, err := strconv.ParseFloat(s[:l], 64)
	if !ok {
		return err
	}
	// fmt.Printf("Set: v=%f, m=%f, v*m=%v\n", v, m, v*m)
	r.V = float64(v * m)
	return err
}

func (v Int64) String() string {
	//	fmt.Printf("String: I\n")
	return fmt.Sprintf("%s", pif(v.V, v.U, 0, 0, true, true, order))
}

func (v Int) String() string {
	//	fmt.Printf("String: I\n")
	return fmt.Sprintf("%s", pif(int64(v.V), v.U, 0, 0, true, true, order))
}

func (v Float64) String() string {
	//	fmt.Printf("String: F\n")
	return fmt.Sprintf("%s", pff(v.V, v.U, 0, 0, true, true, order))
}

func (v Int64) Format(s fmt.State, c rune) {
	i(&v, s, c)
}

func (v Int) Format(s fmt.State, c rune) {
	v2 := Int64{int64(v.V), v.U}
	i(&v2, s, c)
}

func (v Float64) Format(s fmt.State, c rune) {
	f(&v, s, c)
}
