package xtemplate

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTag(t *testing.T) {
	Convey("template", t, func() {
		s := `a {{ .val1 }} {{ .val2 }}`
		s1, err := Execute(s, map[string]interface{}{"val1": "b", "val2": 2})
		So(err, ShouldBeNil)
		So(string(s1), ShouldEqual, "a b 2")
	})
}
