package xip

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestXIP(t *testing.T) {
	Convey("xip", t, func() {
		ips, err := LocalIpv4Addrs()
		So(err, ShouldBeNil)
		So(len(ips), ShouldBeGreaterThan, 0)
		t.Log("ip list:", ips)

		localIP := GetLocalIP()
		So(localIP, ShouldNotBeEmpty)
		So(localIP, ShouldNotEqual, "127.0.0.1")
		t.Logf("%s is intranet? %v", localIP, IsIntranet(localIP))

		So(IsValidIP4("127.0"), ShouldBeFalse)
		So(IsValidIP4("127.0.0.1"), ShouldBeTrue)
		So(IsValidIP4("127.0.0.1:0"), ShouldBeTrue)
	})
}
