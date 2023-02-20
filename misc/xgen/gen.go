package xgen

import (
	"bytes"
	"fmt"
	"strings"
)

type Gen struct {
	*bytes.Buffer
	indent      string
	cellIndent  string
	indentCount int
}

// In Indents the output one tab stop.
func (g *Gen) In() {
	g.indent += g.cellIndent
	g.indentCount++
}

// Out unindents the output one tab stop.
func (g *Gen) Out() {
	if g.indentCount > 0 {
		g.indent = g.indent[len(g.cellIndent):]
		g.indentCount--
	}
}

func (g *Gen) Indent() string {
	return g.indent
}
func (g *Gen) IndentCount() int {
	return g.indentCount
}

func (g *Gen) SetPrefixIndent(base int) {
	for i := 0; i < base; i++ {
		g.In()
	}
}

func (g *Gen) OutPrefixIndent(base int) {
	for i := 0; i < base; i++ {
		g.Out()
	}
}

func (g *Gen) TrimRight(s string) {
	ret := bytes.TrimRight(g.Buffer.Bytes(), s)
	g.Buffer.Reset()
	g.Buffer.Write(ret)
}

func (g *Gen) PFormat(fmtStr string, str ...interface{}) {
	g.P(fmt.Sprintf(fmtStr, str...))
}

// P prints the arguments to the generated output.  It handles strings and int32s, plus
// handling indirections because they may be *string, etc.
func (g *Gen) P(str ...interface{}) {
	_, _ = g.WriteString(g.indent)
	for _, v := range str {
		switch s := v.(type) {
		case string:
			_, _ = g.WriteString(s)
		case *string:
			_, _ = g.WriteString(*s)
		case bool:
			_, _ = fmt.Fprintf(g, "%t", s)
		case *bool:
			_, _ = fmt.Fprintf(g, "%t", *s)
		case int:
			_, _ = fmt.Fprintf(g, "%d", s)
		case *int32:
			_, _ = fmt.Fprintf(g, "%d", *s)
		case *int64:
			_, _ = fmt.Fprintf(g, "%d", *s)
		case float64:
			_, _ = fmt.Fprintf(g, "%g", s)
		case *float64:
			_, _ = fmt.Fprintf(g, "%g", *s)
		default:
			g.Fail(fmt.Sprintf("unknown type in printer: %T", v))
		}
	}
	_ = g.WriteByte('\n')
}

// Fail reports a problem and exits the program.
func (g *Gen) Fail(msgs ...string) {
	panic("ProtoKit: error:" + strings.Join(msgs, " "))
}

func (g *Gen) StringTrimNewline() string {
	return strings.TrimRight(g.String(), "\n")
}

func NewGeneratorWithTabIndent() *Gen {
	return &Gen{Buffer: new(bytes.Buffer), cellIndent: "\t"}
}
func NewGeneratorWithSpaceIndent(n int) *Gen {
	return &Gen{Buffer: new(bytes.Buffer), cellIndent: strings.Repeat(" ", n)}
}
