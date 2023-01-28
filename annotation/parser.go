package annotation

import (
	"fmt"
	"github.com/sandwich-go/boost/xstrings"
	"strings"
	"text/scanner"
)

func parser(line string, lowerKey bool) (Annotation, error) {
	var (
		s           scanner.Scanner
		token       rune
		currentStep = initialStep
		attrName    string
		ann         = annotation{
			line:       line,
			attributes: make(map[string]string),
		}
	)
	s.Init(strings.NewReader(strings.TrimLeft(strings.TrimSpace(line), "/")))

	for token != scanner.EOF && currentStep < doneStep {
		token = s.Scan()
		switch token {
		case '@':
			currentStep = annotationNameStep
		case '(':
			currentStep = attributeNameStep
		case '=':
			currentStep = attributeValueStep
		case ',':
			currentStep = attributeNameStep
		case ')':
			currentStep = doneStep
		case scanner.Ident:
			switch currentStep {
			case annotationNameStep:
				if n := s.TokenText(); len(n) > 0 {
					ann.name = n
				}
			case attributeNameStep:
				if n := s.TokenText(); len(n) > 0 {
					attrName = n
				}
			}
		default:
			switch currentStep {
			case attributeValueStep:
				var key = attrName
				if lowerKey {
					key = strings.ToLower(key)
				}
				ann.attributes[key] =
					xstrings.Trim(strings.Trim(strings.Trim(strings.Trim(s.TokenText(), "\""), "'"), "`"))
			}
		}
	}
	if currentStep != doneStep {
		return nil, fmt.Errorf("invalid completion-status %v name:%s for annotation:%s", currentStep, attrName, line)
	}
	return ann, nil
}
