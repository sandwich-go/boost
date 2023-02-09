package version

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestVersion(t *testing.T) {
	Convey(`version`, t, func() {
		for _, ver := range []struct {
			version  string
			userData string
			valid    bool
			expected string
		}{
			{version: "unknown", valid: false, expected: "unknown_unknown_unknown_unknown_unknown"},
			{version: "1.3.2", valid: true, expected: "1.3.2_unknown_unknown_unknown_unknown"},
			{version: "1.3.2", userData: "__data__", valid: true, expected: "1.3.2_unknown_unknown_unknown_unknown___data__"},
		} {
			Version = ver.version
			UserData = ver.userData
			So(ver.valid, ShouldEqual, Valid())
			So(ver.expected, ShouldEqual, String())
		}
	})
}
