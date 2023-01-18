package annotation

import (
	"fmt"
	"github.com/sandwich-go/boost/xstrings"
	"strconv"
	"strings"
)

const all = "*"

type resolver struct {
	opts *Options
}

// New 创建一个解析器，默认使用 annotation@ 作为 MagicPrefix，只有包含 MagicPrefix 的行，才能萃取到注释
// descriptors 可以指定萃取的注释名，以及是否为合法的注释
// 若不指定 descriptors，则任意名的注释，均为合法的注释
func New(opts ...Option) Resolver {
	v := &resolver{opts: NewOptions(opts...)}
	if len(v.opts.GetDescriptors()) == 0 {
		v.opts.ApplyOption(WithDescriptors(Descriptor{Name: all}))
	}
	return v
}

type annotation struct {
	name       string
	line       string
	attributes map[string]string
}

func (a annotation) Name() string { return a.name }
func (a annotation) Line() string { return a.line }
func (a annotation) Contains(key string) bool {
	_, ok := a.attributes[key]
	return ok
}

func (a annotation) String(key string, defaultVal ...string) string {
	v, ok := a.attributes[key]
	if !ok {
		return append(defaultVal, "")[0]
	}
	return xstrings.Trim(v)
}

func (a annotation) Int8(key string, defaultVal ...int8) (int8, error) {
	val := a.String(key)
	if len(val) == 0 {
		return append(defaultVal, 0)[0], nil
	}
	iVal, err := strconv.ParseInt(val, 10, 8)
	if err != nil {
		return 0, err
	}
	return int8(iVal), nil
}

func (a annotation) Int16(key string, defaultVal ...int16) (int16, error) {
	val := a.String(key)
	if len(val) == 0 {
		return append(defaultVal, 0)[0], nil
	}
	iVal, err := strconv.ParseInt(val, 10, 16)
	if err != nil {
		return 0, err
	}
	return int16(iVal), nil
}

func (a annotation) Int32(key string, defaultVal ...int32) (int32, error) {
	val := a.String(key)
	if len(val) == 0 {
		return append(defaultVal, 0)[0], nil
	}
	iVal, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(iVal), nil
}

func (a annotation) Int64(key string, defaultVal ...int64) (int64, error) {
	val := a.String(key)
	if len(val) == 0 {
		return append(defaultVal, 0)[0], nil
	}
	return strconv.ParseInt(val, 10, 64)
}

func (a annotation) Uint8(key string, defaultVal ...uint8) (uint8, error) {
	val := a.String(key)
	if len(val) == 0 {
		return append(defaultVal, 0)[0], nil
	}
	iVal, err := strconv.ParseUint(val, 10, 8)
	if err != nil {
		return 0, err
	}
	return uint8(iVal), nil
}

func (a annotation) Uint16(key string, defaultVal ...uint16) (uint16, error) {
	val := a.String(key)
	if len(val) == 0 {
		return append(defaultVal, 0)[0], nil
	}
	iVal, err := strconv.ParseUint(val, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(iVal), nil
}

func (a annotation) Uint32(key string, defaultVal ...uint32) (uint32, error) {
	val := a.String(key)
	if len(val) == 0 {
		return append(defaultVal, 0)[0], nil
	}
	iVal, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(iVal), nil
}

func (a annotation) Uint64(key string, defaultVal ...uint64) (uint64, error) {
	val := a.String(key)
	if len(val) == 0 {
		return append(defaultVal, 0)[0], nil
	}
	return strconv.ParseUint(val, 10, 64)
}

func (a annotation) Int(key string, defaultVal ...int) (int, error) {
	val := a.String(key)
	if len(val) == 0 {
		return append(defaultVal, 0)[0], nil
	}
	iVal, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(iVal), nil
}

func (a annotation) Float32(key string, defaultVal ...float32) (float32, error) {
	val := a.String(key)
	if len(val) == 0 {
		return append(defaultVal, 0)[0], nil
	}
	iVal, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return 0, err
	}
	return float32(iVal), nil
}

func (a annotation) Float64(key string, defaultVal ...float64) (float64, error) {
	val := a.String(key)
	if len(val) == 0 {
		return append(defaultVal, 0)[0], nil
	}
	return strconv.ParseFloat(val, 64)
}

func (a annotation) Bool(key string, defaultVal ...bool) (bool, error) {
	val := a.String(key)
	if len(val) == 0 {
		return append(defaultVal, false)[0], nil
	}
	return xstrings.IsTrue(xstrings.Trim(strings.ToLower(val))), nil
}

func (desc Descriptor) match(ann Annotation) bool {
	if desc.Name != all && ann.Name() != desc.Name {
		return false
	}
	if desc.Validator != nil {
		return desc.Validator(ann)
	}
	return true
}

func (r *resolver) Resolve(line string) (Annotation, error) {
	line = strings.TrimSpace(line)
	for _, desc := range r.opts.GetDescriptors() {
		if !strings.Contains(line, r.opts.GetMagicPrefix()) {
			continue
		}
		ann, err := parser(line, r.opts.GetLowerKey())
		if err != nil {
			return ann, err
		}
		if !desc.match(ann) {
			continue
		}
		return ann, nil
	}
	return annotation{}, ErrNoAnnotation
}

func (r *resolver) ResolveMultiple(lines []string) ([]Annotation, error) {
	as := make([]Annotation, 0, len(lines))
	for _, line := range lines {
		if ann, err := r.Resolve(line); err == nil {
			as = append(as, ann)
		} else if err != ErrNoAnnotation {
			return nil, err
		}
	}
	return as, nil
}

func (r *resolver) ResolveWithName(lines []string, name string) (Annotation, error) {
	for _, line := range lines {
		if ann, err := r.Resolve(line); err == nil && ann.Name() == name {
			return ann, nil
		} else if err != nil && err != ErrNoAnnotation {
			return annotation{}, err
		}
	}
	return annotation{}, ErrNoAnnotation
}

func (r *resolver) ResolveNoDuplicate(lines []string) ([]Annotation, error) {
	as, err := r.ResolveMultiple(lines)
	if err != nil || len(as) == 0 {
		return nil, err
	}
	mapping := make(map[string]Annotation)
	for _, v := range as {
		if v1, ok := mapping[v.Name()]; ok {
			return nil, fmt.Errorf("got duplicate annotation name with name: %s line1: %s line2: %s", v.Name(), v.Line(), v1.Line())
		}
		mapping[v.Name()] = v
	}
	return as, nil
}
