package xip

import (
	"net"
	"regexp"
	"strconv"
	"strings"

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
	return xslice.RemoveRepeated(ips), nil
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

var ip4Reg = regexp.MustCompile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

// IsValidIP4 是否是合法的 ip4 地址
func IsValidIP4(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")
	if i := strings.LastIndex(ipAddress, ":"); i >= 0 {
		ipAddress = ipAddress[:i] //remove port
	}
	return ip4Reg.MatchString(ipAddress)
}
