package xstrings

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestJson(t *testing.T) {
	Convey("json", t, func() {
		s := `{a: 1}`
		So(ValidJSON([]byte(s)), ShouldBeFalse)

		s = `{"a": 1}`
		So(ValidJSON([]byte(s)), ShouldBeTrue)

		s = `
{
	"a": 1,
	"b": 2
}
`
		So(string(UglyJSON([]byte(s))), ShouldEqual, `{"a":1,"b":2}`)
	})
}
