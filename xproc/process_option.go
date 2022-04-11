package xproc

import (
	"io"
	"os"
)

//go:generate optiongen --option_return_previous=false
func ProcessOptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"Args":       []string{},
		"Stdin":      io.Reader(os.Stdin),
		"Stdout":     io.Writer(os.Stdout),
		"Stderr":     io.Writer(os.Stderr),
		"WorkingDir": "",
		"Env":        []string{},
	}
}
