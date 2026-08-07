package xip

import (
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/sandwich-go/boost/xos"
	"github.com/sandwich-go/boost/xslice"
)

const boostIPPreferVPN = "boost_ip_prefer_vpn"
const cmdEnvKeyForServiceHost = "x_sandwich_service_host"

// LocalIpv4Addrs scan all ip addresses with loopback excluded.
// If VPN is connected, prioritize returning the VPN IP address.
func LocalIpv4Addrs() (ips []string, err error) {
	ips = make([]string, 0)

	ifaces, e := net.Interfaces()
	if e != nil {
		return ips, e
	}

	// 优先检查 VPN 接口
	if xos.EnvGetCaseInsensitive(boostIPPreferVPN) != "" {
		for _, iface := range ifaces {
			if isVPNInterface(iface) {
				addrs, e := iface.Addrs()
				if e != nil {
					continue
				}

				for _, addr := range addrs {
					ip := getIPFromAddr(addr)
					if ip != nil && ip.To4() != nil {
						ips = append(ips, ip.String())
					}
				}
			}
		}
	}
	// 检查 K8S 环境变量注入的
	if k8sIP := xos.EnvGetCaseInsensitive(cmdEnvKeyForServiceHost); k8sIP != "" {
		if IsValidIP4(k8sIP) {
			ips = append(ips, k8sIP)
		}
	}

	// 如果没有找到 VPN 接口的 IP，继续检查其他接口
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}

		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}

		// ignore docker and warden bridge
		if strings.HasPrefix(iface.Name, "docker") || strings.HasPrefix(iface.Name, "w-") {
			continue
		}

		addrs, e := iface.Addrs()
		if e != nil {
			return ips, e
		}

		for _, addr := range addrs {
			ip := getIPFromAddr(addr)
			if ip != nil && ip.To4() != nil && IsIntranet(ip.String()) {
				ips = append(ips, ip.String())
			}
		}
	}
	return xslice.StringsRemoveRepeated(ips), nil
}

// isVPNInterface 判断是否是 VPN 接口
func isVPNInterface(iface net.Interface) bool {
	// 常见的 VPN 接口名称前缀
	vpnPrefixes := []string{"tun", "utun", "ppp", "tap"}
	for _, prefix := range vpnPrefixes {
		if strings.HasPrefix(iface.Name, prefix) {
			return true
		}
	}

	// 检查接口标志（某些 VPN 接口可能有特定标志）
	if iface.Flags&net.FlagPointToPoint != 0 {
		return true
	}

	return false
}

// getIPFromAddr 从 net.Addr 中提取 IP 地址
func getIPFromAddr(addr net.Addr) net.IP {
	switch v := addr.(type) {
	case *net.IPNet:
		return v.IP
	case *net.IPAddr:
		return v.IP
	}
	return nil
}

// IsIntranet 是否是内网地址
func IsIntranet(ipStr string) bool {
	if strings.HasPrefix(ipStr, "10.") || strings.HasPrefix(ipStr, "192.168.") {
		return true
	}
	if strings.HasPrefix(ipStr, "172.") {
		// 172.16.0.0-172.31.255.255
		arr := strings.Split(ipStr, ".")
		if len(arr) != 4 {
			return false
		}
		second, err := strconv.ParseInt(arr[1], 10, 64)
		if err != nil {
			return false
		}
		if second >= 16 && second <= 31 {
			return true
		}
	}
	return false
}

// GetLocalIP returns the non loopback local IP of the host
// If VPN is connected, prioritize returning the VPN IP address.
func GetLocalIP() string {
	addrs, err := LocalIpv4Addrs()
	if err != nil || len(addrs) == 0 {
		return ""
	}
	return addrs[0]
}

// localIPCache 缓存 GetLocalIPCached 首次成功解析到的地址。
// 只在 nil -> 非 nil 方向写一次，读路径无锁。
var localIPCache atomic.Pointer[string]

// localIPResolver 是 GetLocalIPCached 的解析入口，供包内测试替换以覆盖解析失败分支。
// 仅在缓存未命中时调用，不在读路径上。
var localIPResolver = GetLocalIP

// GetLocalIPCached 返回与 GetLocalIP 相同的地址，但整个进程只解析一次，之后的调用直接读缓存。
//
// GetLocalIP 每次调用都会枚举一遍网卡（net.Interfaces 走系统调用）并两次全量扫描环境变量
// （xos.EnvGetCaseInsensitive 里 os.Environ 拷贝全量环境 + 逐项 strings.SplitN），开销随
// 环境变量条数增长：darwin/arm64 实测单次 79µs、42216 B、581 allocs。本机 IP 在进程生命
// 周期内不会变，按请求频率调用会把这份开销直接变成 GC 压力，所以热路径上取本机 IP 用本函数，
// 只在启动期取一次的场景用 GetLocalIP 即可。
//
// 解析失败（返回空串）时不写缓存，下次调用重新解析，避免把一次瞬时失败固化成进程终身返回空串；
// 这条退化路径的开销与直接调 GetLocalIP 相同，不会比原来更差。
//
// 缓存在首次调用时填充而非包初始化时填充，因此依赖 boost_ip_prefer_vpn /
// x_sandwich_service_host 的调用方必须在首次调用前把环境变量设置好。
func GetLocalIPCached() string {
	if cached := localIPCache.Load(); cached != nil {
		return *cached
	}
	ip := localIPResolver()
	if ip == "" {
		return ""
	}
	localIPCache.Store(&ip)
	return ip
}

var ip4Reg = regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

// IsValidIP4 是否是合法的 ip4 地址
func IsValidIP4(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")
	if i := strings.LastIndex(ipAddress, ":"); i >= 0 {
		ipAddress = ipAddress[:i] //remove port
	}
	return ip4Reg.MatchString(ipAddress)
}
