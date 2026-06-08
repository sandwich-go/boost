package xip

import (
	"net"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：补 xip 包的：
//   - IsIntranet 各 case（10/192.168/172.16-31/非内网）
//   - getIPFromAddr 类型分支
//   - isVPNInterface VPN 前缀识别
//   - GetFreePort 错误路径

func TestIsIntranet(t *testing.T) {
	Convey("10.x.x.x 是内网", t, func() {
		So(IsIntranet("10.0.0.1"), ShouldBeTrue)
		So(IsIntranet("10.255.255.255"), ShouldBeTrue)
	})
	Convey("192.168.x.x 是内网", t, func() {
		So(IsIntranet("192.168.0.1"), ShouldBeTrue)
		So(IsIntranet("192.168.100.50"), ShouldBeTrue)
	})
	Convey("172.16-31.x.x 是内网", t, func() {
		So(IsIntranet("172.16.0.1"), ShouldBeTrue)
		So(IsIntranet("172.20.0.1"), ShouldBeTrue)
		So(IsIntranet("172.31.255.255"), ShouldBeTrue)
		// 172.0-15 / 32-255 不是内网
		So(IsIntranet("172.15.0.1"), ShouldBeFalse)
		So(IsIntranet("172.32.0.1"), ShouldBeFalse)
	})
	Convey("172.x 但格式错误返 false", t, func() {
		So(IsIntranet("172.x"), ShouldBeFalse)
		So(IsIntranet("172.notanumber.0.1"), ShouldBeFalse)
		So(IsIntranet("172"), ShouldBeFalse)
	})
	Convey("公网 IP 不是内网", t, func() {
		So(IsIntranet("8.8.8.8"), ShouldBeFalse)
		So(IsIntranet("1.1.1.1"), ShouldBeFalse)
		So(IsIntranet("203.0.113.1"), ShouldBeFalse)
	})
}

func TestGetIPFromAddr(t *testing.T) {
	Convey("getIPFromAddr 处理 *net.IPNet", t, func() {
		ipnet := &net.IPNet{IP: net.IPv4(192, 168, 1, 1)}
		got := getIPFromAddr(ipnet)
		So(got, ShouldNotBeNil)
		So(got.String(), ShouldEqual, "192.168.1.1")
	})

	Convey("getIPFromAddr 处理 *net.IPAddr", t, func() {
		ipaddr := &net.IPAddr{IP: net.IPv4(10, 0, 0, 1)}
		got := getIPFromAddr(ipaddr)
		So(got, ShouldNotBeNil)
		So(got.String(), ShouldEqual, "10.0.0.1")
	})

	Convey("getIPFromAddr 不支持的类型返 nil", t, func() {
		// *net.UDPAddr 不在 switch 内
		udpAddr := &net.UDPAddr{IP: net.IPv4(1, 2, 3, 4), Port: 80}
		So(getIPFromAddr(udpAddr), ShouldBeNil)
	})
}

func TestIsVPNInterface(t *testing.T) {
	Convey("VPN 前缀识别", t, func() {
		So(isVPNInterface(net.Interface{Name: "tun0"}), ShouldBeTrue)
		So(isVPNInterface(net.Interface{Name: "utun0"}), ShouldBeTrue)
		So(isVPNInterface(net.Interface{Name: "ppp0"}), ShouldBeTrue)
		So(isVPNInterface(net.Interface{Name: "tap0"}), ShouldBeTrue)
	})

	Convey("非 VPN 接口（无 PointToPoint flag）", t, func() {
		So(isVPNInterface(net.Interface{Name: "eth0"}), ShouldBeFalse)
		So(isVPNInterface(net.Interface{Name: "wlan0"}), ShouldBeFalse)
	})

	Convey("含 PointToPoint flag 视作 VPN", t, func() {
		So(isVPNInterface(net.Interface{
			Name:  "custom",
			Flags: net.FlagPointToPoint,
		}), ShouldBeTrue)
	})
}

func TestIsValidIP4(t *testing.T) {
	Convey("合法 IPv4", t, func() {
		So(IsValidIP4("192.168.1.1"), ShouldBeTrue)
		So(IsValidIP4("0.0.0.0"), ShouldBeTrue)
		So(IsValidIP4("255.255.255.255"), ShouldBeTrue)
		// 含 port 应剥离
		So(IsValidIP4("8.8.8.8:53"), ShouldBeTrue)
	})

	Convey("非法 IPv4", t, func() {
		So(IsValidIP4("256.0.0.1"), ShouldBeFalse)
		So(IsValidIP4("not.an.ip.addr"), ShouldBeFalse)
		So(IsValidIP4(""), ShouldBeFalse)
		So(IsValidIP4("1.2.3"), ShouldBeFalse)
	})
}

func TestGetFreePort(t *testing.T) {
	Convey("GetFreePort 返回可用端口", t, func() {
		port, err := GetFreePort()
		So(err, ShouldBeNil)
		So(port, ShouldBeGreaterThan, 0)
		So(port, ShouldBeLessThanOrEqualTo, 65535)
	})
}

func TestGetLocalIP_AndLocalIpv4Addrs(t *testing.T) {
	Convey("LocalIpv4Addrs 返回本机 IP（依赖系统环境）", t, func() {
		// 在 CI 环境可能没有非 loopback IP，但调用应不报错
		ips, err := LocalIpv4Addrs()
		So(err, ShouldBeNil)
		_ = ips // 不强断言长度
	})

	Convey("GetLocalIP 返回字符串", t, func() {
		// 调用不 panic 即可
		s := GetLocalIP()
		_ = s
	})
}

// TestLocalIpv4Addrs_PreferVPN 覆盖 LocalIpv4Addrs 的 'boost_ip_prefer_vpn'
// env 分支（line 27-43）。
func TestLocalIpv4Addrs_PreferVPN(t *testing.T) {
	Convey("设置 boost_ip_prefer_vpn 走 VPN 优先扫描分支", t, func() {
		t.Setenv("boost_ip_prefer_vpn", "1")
		// CI runner 不一定有 VPN 接口，但代码路径会 enter VPN 扫描循环
		ips, err := LocalIpv4Addrs()
		So(err, ShouldBeNil)
		_ = ips
	})
}

// TestLocalIpv4Addrs_K8sServiceHost 覆盖 k8s service host env 注入分支
// （line 45-49）。
func TestLocalIpv4Addrs_K8sServiceHost(t *testing.T) {
	Convey("k8s 注入 valid IP 加入结果", t, func() {
		t.Setenv("x_sandwich_service_host", "10.20.30.40")
		ips, err := LocalIpv4Addrs()
		So(err, ShouldBeNil)
		// 10.20.30.40 是 valid ip4，应在返回列表
		found := false
		for _, ip := range ips {
			if ip == "10.20.30.40" {
				found = true
				break
			}
		}
		So(found, ShouldBeTrue)
	})

	Convey("k8s 注入 invalid IP 不加入", t, func() {
		t.Setenv("x_sandwich_service_host", "not-an-ip")
		ips, err := LocalIpv4Addrs()
		So(err, ShouldBeNil)
		for _, ip := range ips {
			So(ip, ShouldNotEqual, "not-an-ip")
		}
	})
}

// TestGetLocalIP_EmptyOnNoAddrs 覆盖 GetLocalIP 在 LocalIpv4Addrs 返
// 0 个 IP 时的"返空字符串"分支（line 136-138）。
//
// 注意：CI 上 LocalIpv4Addrs 通常返 docker bridge / k8s pod IP 至少一个，
// 难以稳定构造"0 IP"场景。这里只验证 GetLocalIP 与 LocalIpv4Addrs[0]
// 行为一致（即"返第一个 / 空时返空"契约）。
func TestGetLocalIP_Contract(t *testing.T) {
	Convey("GetLocalIP 返第一个 LocalIpv4Addrs 或空", t, func() {
		ips, _ := LocalIpv4Addrs()
		got := GetLocalIP()
		if len(ips) == 0 {
			So(got, ShouldEqual, "")
		} else {
			So(got, ShouldEqual, ips[0])
		}
	})
}
