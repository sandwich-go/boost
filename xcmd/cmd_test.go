package xcmd

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCommand(t *testing.T) {
	cleanParsedOptions()
	Convey("add duplicated command flag should panic", t, func() {
		var name = fmt.Sprintf("%stest", FlagPrefix)
		So(func() { AddFlag(name, "test") }, ShouldNotPanic)
		So(func() { AddFlag(name, "test") }, ShouldPanic)
	})
	cleanParsedOptions()
	Convey("command flag format", t, func() {
		Convey("without FlagPrefix should panic", func() {
			So(func() { AddFlag("test", "test") }, ShouldPanic)
		})
		Convey("contain _ should panic", func() {
			So(func() { AddFlag(fmt.Sprintf("%stest", FlagPrefix), "test") }, ShouldPanic)
		})
	})
	cleanParsedOptions()
	Convey("command init", t, func() {
		var name1 = fmt.Sprintf("%sdebug", FlagPrefix)
		var name2 = fmt.Sprintf("%sp1", FlagPrefix)
		var name3 = fmt.Sprintf("%sp2", FlagPrefix)
		Init(fmt.Sprintf("--%s=true", name1), fmt.Sprintf("--%s", name2), "test", fmt.Sprintf("--%s", name3))
		So(IsTrue(GetOptWithEnv(name1)), ShouldBeTrue)
		So(GetOptWithEnv(name2), ShouldEqual, "test")
		So(IsFalse(GetOptWithEnv(name3)), ShouldBeTrue)
	})
	cleanParsedOptions()
	Convey("flag provided but not defined", t, func() {
		var name1 = strings.TrimSuffix(FlagPrefix, "_")
		var name2 = fmt.Sprintf("%sp1", FlagPrefix)
		Convey("should occur", func() {
			flag.CommandLine.Init(name1, flag.PanicOnError)
			os.Args = append(os.Args, fmt.Sprintf("--%s=true", name2))
			flag.CommandLine.SetOutput(ioutil.Discard)
			So(func() { _ = flag.CommandLine.Parse(os.Args[1:]) }, ShouldPanic)
		})
		Convey("should not occur after call DeclareInto", func() {
			AddFlag(name2, "true")
			So(func() { _ = flag.CommandLine.Parse(os.Args[1:]) }, ShouldNotPanic)
		})
	})
	cleanParsedOptions()
	Convey("parameter order,command line flag => os env", t, func() {
		var name = fmt.Sprintf("%s_p2", FlagPrefix)
		err := os.Setenv(name, "v2")
		So(err, ShouldBeNil)
		os.Args = append(os.Args, fmt.Sprintf("--%s=v1", name))
		So(GetOptWithEnv(name), ShouldEqual, "v1")
		So(GetOptWithEnv(name), ShouldEqual, "v1")
	})
}
