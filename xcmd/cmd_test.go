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

func cleanParsedOptions() { parsedOptions = make(map[string]string) }

func TestCommand(t *testing.T) {
	SetFlagPrefix("test_")

	cleanParsedOptions()
	Convey("add duplicated command flag should panic", t, func() {
		var name = fmt.Sprintf("%stest", GetFlagPrefix())
		So(func() { MustAddFlag(name, "test") }, ShouldNotPanic)
		So(func() { MustAddFlag(name, "test") }, ShouldPanic)
	})

	cleanParsedOptions()
	Convey("command flag format", t, func() {
		Convey("without FlagPrefix should panic", func() {
			So(func() { MustAddFlag("test", "test") }, ShouldPanic)
		})
		Convey("contain invalid char should panic", func() {
			for _, invalidChar := range invalidChars {
				So(func() { MustAddFlag(fmt.Sprintf("%s%stest", GetFlagPrefix(), string(invalidChar)), "test") }, ShouldPanic)
			}
		})
	})

	cleanParsedOptions()
	Convey("command init", t, func() {
		var name1 = fmt.Sprintf("%sdebug", GetFlagPrefix())
		var name2 = fmt.Sprintf("%sp1", GetFlagPrefix())
		var name3 = fmt.Sprintf("%sp2", GetFlagPrefix())
		Init(fmt.Sprintf("--%s=true", name1), fmt.Sprintf("--%s", name2), "test", fmt.Sprintf("--%s", name3))
		So(IsTrue(GetOptWithEnv(name1)), ShouldBeTrue)
		So(GetOptWithEnv(name2), ShouldEqual, "test")
		So(IsFalse(GetOptWithEnv(name3)), ShouldBeTrue)
	})

	cleanParsedOptions()
	Convey("flag provided but not defined", t, func() {
		var name1 = strings.TrimSuffix(GetFlagPrefix(), "_")
		var name2 = fmt.Sprintf("%sp1", GetFlagPrefix())
		Convey("should occur", func() {
			flag.CommandLine.Init(name1, flag.PanicOnError)
			os.Args = append(os.Args, fmt.Sprintf("--%s=true", name2))
			flag.CommandLine.SetOutput(ioutil.Discard)
			So(func() { _ = flag.CommandLine.Parse(os.Args[1:]) }, ShouldPanic)
		})
		Convey("should not occur after call DeclareInto", func() {
			MustAddFlag(name2, "true")
			DeclareInto(flag.CommandLine)
			So(func() { _ = flag.CommandLine.Parse(os.Args[1:]) }, ShouldNotPanic)
		})
	})

	cleanParsedOptions()
	Convey("parameter order,command line flag => os env", t, func() {
		var name = fmt.Sprintf("%s_p2", GetFlagPrefix())
		err := os.Setenv(name, "v2")
		So(err, ShouldBeNil)
		os.Args = append(os.Args, fmt.Sprintf("--%s=v1", name))
		So(ContainsOpt(name), ShouldBeTrue)
		So(GetOptWithEnv(name), ShouldEqual, "v1")
		So(GetOptWithEnv(name), ShouldEqual, "v1")
	})
}
