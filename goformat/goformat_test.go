package goformat

import (
	. "github.com/smartystreets/goconvey/convey"
	"go/token"
	"testing"
)

func TestParseMainFragment(t *testing.T) {
	Convey("parse main fragment", t, func() {
		for _, test := range []struct {
			src         string
			adjust, err bool
		}{
			{
				src: `
func a() error {return nil}
`, adjust: true,
			},
			{
				src: `
func main() error {return nil}
`, adjust: true,
			},
			{
				src: `
func main() {}
`,
			},
			{
				src: `
err := db.Connect()
`, err: true,
			},
		} {
			_, adjust, err := parseMainFragment(token.NewFileSet(), "", []byte(test.src), 0)
			if test.err {
				So(err, ShouldNotBeNil)
			} else {
				So(err, ShouldBeNil)
			}
			if test.adjust {
				So(adjust, ShouldNotBeNil)
			} else {
				So(adjust, ShouldBeNil)
			}
		}
	})
}

func TestParseDeclarationFragment(t *testing.T) {
	Convey("parse declaration fragment", t, func() {
		for _, test := range []struct {
			src         string
			adjust, err bool
		}{
			{
				src: `
func a() error {return nil}
`, err: true,
			},
			{
				src: `
func main() {}
`, err: true,
			},
			{
				src: `
err := db.Connect()
`, adjust: true,
			},
		} {
			_, adjust, err := parseDeclarationFragment(token.NewFileSet(), "", []byte(test.src), 0)
			if test.err {
				So(err, ShouldNotBeNil)
			} else {
				So(err, ShouldBeNil)
			}
			if test.adjust {
				So(adjust, ShouldNotBeNil)
			} else {
				So(adjust, ShouldBeNil)
			}
		}
	})
}

func TestParse(t *testing.T) {
	Convey("parse", t, func() {
		for _, test := range []struct {
			src         string
			fragment    bool
			adjust, err bool
		}{
			{
				src: `
package main

func main() {}
`,
			},
			{
				src: `
package a

func b() error { return nil }
`,
			},
			{
				src: `
func b() error { return nil }
`, err: true,
			},
			{
				src: `
func a() error {return nil}
`, fragment: true, adjust: true,
			},
			{
				src: `
func main() error {return nil}
`, fragment: true, adjust: true,
			},
			{
				src: `
func main() {}
`, fragment: true,
			},
			{
				src: `
err := db.Connect()
`, fragment: true, adjust: true,
			},
		} {
			_, adjust, err := parse(token.NewFileSet(), "", []byte(test.src), NewOptions(WithFragment(test.fragment)))
			if test.err {
				So(err, ShouldNotBeNil)
			} else {
				So(err, ShouldBeNil)
			}
			if test.adjust {
				So(adjust, ShouldNotBeNil)
			} else {
				So(adjust, ShouldBeNil)
			}
		}
	})
}

func TestProcess(t *testing.T) {
	Convey("process", t, func() {
		for _, test := range []struct {
			fileName string
			code     string
			err      bool
			expected string
			fragment bool
		}{
			{
				fileName: "testgo",
				expected: `package test_file

func a() { return }
`,
			},
			{
				code:     "func a(     ) {return}",
				fragment: true,
				expected: `func a() { return }`,
			},
			{
				code: "func a(     ) {return}",
				err:  true,
			},
		} {
			var err error
			var out []byte
			if len(test.fileName) > 0 {
				out, err = ProcessFile(test.fileName, WithFragment(test.fragment))
			} else {
				out, err = ProcessCode([]byte(test.code), WithFragment(test.fragment))
			}
			if test.err {
				So(err, ShouldNotBeNil)
			} else {
				So(err, ShouldBeNil)
				So(out, ShouldResemble, []byte(test.expected))
			}
		}
	})
}
