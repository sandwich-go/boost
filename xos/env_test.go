package xos

import (
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestEnv(t *testing.T) {
	Convey("env", t, func() {
		var envName, envValue string
		So(os.Setenv(envName, envValue), ShouldNotBeNil)
		envName, envValue = "a", "b"
		So(os.Setenv(envName, envValue), ShouldBeNil)

		So(EnvGet(envName), ShouldEqual, envValue)
		So(EnvGet(envName+"1", envValue), ShouldEqual, envValue)
	})
}
