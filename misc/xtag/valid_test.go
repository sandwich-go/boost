package xtag

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTag(t *testing.T) {
	Convey("valid tag", t, func() {
		t.Log(ValidateStructTag("a:1"))
		So(ValidateStructTag("a:1"), ShouldNotBeNil)
		So(ValidateStructTag(`json:"value"`), ShouldBeNil)
		So(ValidateStructTag(`a:"value"`), ShouldBeNil)
	})
}
