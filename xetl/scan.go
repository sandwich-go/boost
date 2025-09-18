package scan

import (
	"database/sql"
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// Scanner scans a row of []interface{} into a struct pointed to by dst.
// It must be created by NewScanner.
type Scanner func(row []string, dst interface{}) error

// NewScanner takes a header (column names) and a dst pointer to struct,
// and builds a Scanner that maps []interface{} rows into the struct fields.
func NewScanner(header []string, dst interface{}) (Scanner, error) {
	st := reflect.ValueOf(dst)
	if st.Kind() != reflect.Ptr {
		panic("scan: dst must be a pointer to a struct type")
	}
	st = reflect.Indirect(st)
	if !st.IsValid() || st.Type().Kind() != reflect.Struct {
		panic("scan: dst must be a pointer to a struct type")
	}

	var setters []setter
	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)
		if field.PkgPath != "" {
			continue // unexported
		}
		tag := field.Tag.Get("db")
		if tag == "" {
			continue
		}
		colIdx := indexOf(header, tag)
		if colIdx == -1 {
			continue
		}
		sf := setter{
			colIdx:   colIdx,
			fieldIdx: i,
			fn:       buildAssignFunc(field.Type),
		}
		if sf.fn == nil {
			return nil, fmt.Errorf("unsupported field type %s for column %q", field.Type, tag)
		}
		setters = append(setters, sf)
	}

	if len(setters) == 0 {
		return nil, errors.New("no matches found between header and db-tagged struct fields")
	}

	originalType := reflect.ValueOf(dst).Type()
	return func(row []string, dst interface{}) error {
		st := reflect.ValueOf(dst)
		if st.Type() != originalType {
			panic("scan: Scanner called with a different type than NewScanner")
		}
		st = reflect.Indirect(st)
		for _, s := range setters {
			if s.colIdx >= len(row) {
				continue
			}
			if err := s.fn(st.Field(s.fieldIdx), row[s.colIdx]); err != nil {
				return fmt.Errorf("scan column %d: %w", s.colIdx, err)
			}
		}
		return nil
	}, nil
}

type setter struct {
	colIdx   int
	fieldIdx int
	fn       func(field reflect.Value, v interface{}) error
}

func buildAssignFunc(t reflect.Type) func(field reflect.Value, v interface{}) error {
	// TextUnmarshaler
	if reflect.PointerTo(t).Implements(textUnmarshalerType) {
		return func(field reflect.Value, v interface{}) error {
			if v == nil {
				return nil
			}
			bs := []byte(fmt.Sprint(v))
			return field.Addr().Interface().(encoding.TextUnmarshaler).UnmarshalText(bs)
		}
	}

	// SQL Null types
	switch t {
	case reflect.TypeOf(sql.NullString{}):
		return func(field reflect.Value, v interface{}) error {
			ns := sql.NullString{}
			if v != nil && fmt.Sprint(v) != "" {
				ns.String = fmt.Sprint(v)
				ns.Valid = true
			}
			field.Set(reflect.ValueOf(ns))
			return nil
		}
	case reflect.TypeOf(sql.NullInt64{}):
		return func(field reflect.Value, v interface{}) error {
			ni := sql.NullInt64{}
			if v != nil && fmt.Sprint(v) != "" {
				val, err := strconv.ParseInt(fmt.Sprint(v), 10, 64)
				if err != nil {
					return err
				}
				ni.Int64 = val
				ni.Valid = true
			}
			field.Set(reflect.ValueOf(ni))
			return nil
		}
	case reflect.TypeOf(sql.NullFloat64{}):
		return func(field reflect.Value, v interface{}) error {
			nf := sql.NullFloat64{}
			if v != nil && fmt.Sprint(v) != "" {
				val, err := strconv.ParseFloat(fmt.Sprint(v), 64)
				if err != nil {
					return err
				}
				nf.Float64 = val
				nf.Valid = true
			}
			field.Set(reflect.ValueOf(nf))
			return nil
		}
	case reflect.TypeOf(sql.NullBool{}):
		return func(field reflect.Value, v interface{}) error {
			nb := sql.NullBool{}
			if v != nil && fmt.Sprint(v) != "" {
				val, err := strconv.ParseBool(fmt.Sprint(v))
				if err != nil {
					return err
				}
				nb.Bool = val
				nb.Valid = true
			}
			field.Set(reflect.ValueOf(nb))
			return nil
		}
	}

	// fallback to basic parsing
	return func(field reflect.Value, v interface{}) error {
		if v == nil {
			return nil
		}
		val := reflect.ValueOf(v)
		if !val.Type().AssignableTo(field.Type()) {
			if val.Type().ConvertibleTo(field.Type()) {
				val = val.Convert(field.Type())
			} else {
				tmp, err := parseBasic(fmt.Sprint(v), field.Type())
				if err != nil {
					return err
				}
				val = tmp
			}
		}
		field.Set(val)
		return nil
	}
}
func indexOf(s []string, x string) int {
	for i := range s {
		if x == s[i] {
			return i
		}
	}
	return -1
}

// parseBasic 将字符串形式的值解析成目标类型的 reflect.Value。
// 目标类型可以是基本类型、[]byte、interface{}，也支持指针（会返回指向目标类型的指针）。
func parseBasic(s string, t reflect.Type) (reflect.Value, error) {
	// handle pointer target: produce a pointer to parsed element
	if t.Kind() == reflect.Ptr {
		ev, err := parseBasic(s, t.Elem())
		if err != nil {
			return reflect.Value{}, err
		}
		pv := reflect.New(t.Elem())
		pv.Elem().Set(ev)
		return pv, nil
	}

	switch t.Kind() {
	case reflect.String:
		return reflect.ValueOf(s).Convert(t), nil
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 { // []byte
			return reflect.ValueOf([]byte(s)).Convert(t), nil
		}
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("parse bool %q: %w", s, err)
		}
		return reflect.ValueOf(b).Convert(t), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bits := intBits(t)
		v, err := strconv.ParseInt(s, 0, bits)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("parse int %q -> %s: %w", s, t.String(), err)
		}
		return reflect.ValueOf(v).Convert(t), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		bits := intBits(t)
		v, err := strconv.ParseUint(s, 0, bits)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("parse uint %q -> %s: %w", s, t.String(), err)
		}
		return reflect.ValueOf(v).Convert(t), nil
	case reflect.Float32, reflect.Float64:
		bits := 32
		if t.Kind() == reflect.Float64 {
			bits = 64
		}
		f, err := strconv.ParseFloat(s, bits)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("parse float %q -> %s: %w", s, t.String(), err)
		}
		return reflect.ValueOf(f).Convert(t), nil
	case reflect.Interface:
		// if target is interface{}, just return string (caller may handle)
		return reflect.ValueOf(s), nil
	default:
		return reflect.Value{}, fmt.Errorf("unsupported conversion to type %s", t.String())
	}

	return reflect.Value{}, fmt.Errorf("unsupported conversion to type %s", t.String())
}

// intBits 返回目标整型的位宽（用于 strconv）
func intBits(t reflect.Type) int {
	switch t.Kind() {
	case reflect.Int8:
		return 8
	case reflect.Int16:
		return 16
	case reflect.Int32:
		return 32
	case reflect.Int64:
		return 64
	case reflect.Int:
		return strconv.IntSize
	case reflect.Uint8:
		return 8
	case reflect.Uint16:
		return 16
	case reflect.Uint32:
		return 32
	case reflect.Uint64:
		return 64
	case reflect.Uint:
		return strconv.IntSize
	case reflect.Uintptr:
		return strconv.IntSize
	default:
		return 0
	}
}

var textUnmarshalerType = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
