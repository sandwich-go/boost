package annotation

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sandwich-go/boost/xstrings"
)

type Register interface {
	ResolveAnnotations(annotationLines []string) []Annotation
	ResolveAnnotationsErrorDuplicate(annotationLines []string) ([]Annotation, error)
	ResolveAnnotationByName(annotationLines []string, name string) Annotation
	ResolveAnnotation(annotationLines string) (Annotation, bool)
}

func ResolveAnnotationByName(methodCommentLines []string, key string) Annotation {
	registry := NewRegistry()
	return registry.ResolveAnnotationByName(methodCommentLines, key)
}

type annotationRegistry struct {
	descriptors []*Descriptor
}

const all = "*"

func NewRegistry(descriptors ...*Descriptor) Register {
	v := &annotationRegistry{
		descriptors: descriptors,
	}
	if len(v.descriptors) == 0 {
		v.descriptors = append(v.descriptors, &Descriptor{Name: all})
	}
	return v
}

type Annotation struct {
	Name       string
	Line       string
	Attributes map[string]string
}

func (a Annotation) Has(key string) bool {
	_, ok := a.Attributes[key]
	return ok
}
func (a Annotation) GetInt(key string, defaultVal ...int) int {
	v, ok := a.Attributes[key]
	if !ok {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return 0
	}
	lowerV := xstrings.Trim(strings.ToLower(v))
	i, err := strconv.Atoi(lowerV)
	if err != nil {
		panic(fmt.Errorf("got err:%s while parse key:%s with val:%s", err.Error(), key, lowerV))
	}
	return i
}
func (a Annotation) GetBool(key string, defaultVal ...bool) bool {
	v, ok := a.Attributes[key]
	if !ok {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return false
	}
	return xstrings.IsTrue(xstrings.Trim(strings.ToLower(v)))
}

func (a Annotation) GetString(key string, defaultVal ...string) string {
	v, ok := a.Attributes[key]
	if !ok {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		}
		return ""
	}
	return xstrings.Trim(v)
}

type validationFunc func(a Annotation) bool

type Descriptor struct {
	Name      string
	Validator validationFunc
}

func (ar *annotationRegistry) ResolveAnnotations(annotationLines []string) []Annotation {
	annotations := make([]Annotation, 0)
	for _, line := range annotationLines {
		if ann, ok := ar.ResolveAnnotation(strings.TrimSpace(line)); ok {
			annotations = append(annotations, ann)
		}
	}
	return annotations
}
func (ar *annotationRegistry) ResolveAnnotationsErrorDuplicate(annotationLines []string) ([]Annotation, error) {
	annotations := make([]Annotation, 0)
	for _, line := range annotationLines {
		if ann, ok := ar.ResolveAnnotation(strings.TrimSpace(line)); ok {
			for _, v := range annotations {
				if v.Name == ann.Name {
					return nil, fmt.Errorf("got duplicate annotation name with name: %s line1: %s line2: %s", v.Name, v.Line, ann.Line)
				}
			}
			annotations = append(annotations, ann)
		}
	}
	return annotations, nil
}

func (ar *annotationRegistry) ResolveAnnotationByName(annotationLines []string, name string) Annotation {
	for _, line := range annotationLines {
		ann, ok := ar.ResolveAnnotation(strings.TrimSpace(line))
		if ok && ann.Name == name {
			return ann
		}
	}
	return Annotation{}
}

func (ar *annotationRegistry) ResolveAnnotation(annotationLines string) (Annotation, bool) {
	for _, descriptor := range ar.descriptors {
		if !strings.Contains(annotationLines, magicPrefix) {
			continue
		}

		ann, err := parseAnnotation(annotationLines)
		if err != nil {
			panic(fmt.Errorf("got error while parse annotation with line: %s err: %s", annotationLines, err.Error()))
		}
		if descriptor.Name != all && ann.Name != descriptor.Name {
			continue
		}

		if descriptor.Validator != nil {
			ok := descriptor.Validator(ann)
			if !ok {
				continue
			}

		}
		return ann, true
	}
	return Annotation{}, false
}
