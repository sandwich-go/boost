package md5

import (
	"bytes"
	"fmt"
	"github.com/sandwich-go/boost/xhash/nhash/jenkins"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestMd5(t *testing.T) {
	Convey("md5", t, func() {
		s, err := File("test")
		So(err, ShouldBeNil)

		f, err0 := os.Open("test")
		So(err0, ShouldBeNil)
		defer func() {
			_ = f.Close()
		}()
		s1, err1 := Buffer(f)
		So(err1, ShouldBeNil)
		So(s, ShouldEqual, s1)

		s2, err2 := Buffer(bytes.NewReader([]byte("aaaaaaaa")))
		So(err2, ShouldBeNil)
		t.Log(s2)

		hint, _ := jenkins.HashString("aaaaaaaa", 0, 0)
		fmt.Println(hint)
	})
}
