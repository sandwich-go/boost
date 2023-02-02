package dns

import (
	"context"
	"errors"
	"net"
	"time"
)

var ErrNotFound = errors.New("not found ip address list")

// Dialer 拨号器
type Dialer interface {
	// DialContext 连接
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

// Resolver 解析器
type Resolver interface {
	// LookupIPAddr 通过主机搜索 ip 地址列表
	LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error)
}

// DNS 域名解析器
type DNS interface {
	Resolver

	// GetDialContext 获取 Dial 函数
	GetDialContext() func(ctx context.Context, network, address string) (net.Conn, error)
}

// EmptyIpCache 空的 ip 缓存信息
var EmptyIpCache = IpCache{}

// IpCache ip 缓存信息
type IpCache struct {
	IPAddrs   []net.IPAddr // ip 列表
	CreatedAt time.Time    // 创建时间，如果不设置，则表示永久有效
}

// CacheDNS 带缓存的域名解析器
type CacheDNS interface {
	DNS

	// Set 设置 host 对应的 ip 缓存信息
	Set(host string, cache IpCache)
	// Remove 删除 host 对应的 ip 缓存信息
	Remove(host string)
	// Get 获取 host 对应的 ip 缓存信息
	Get(host string) (IpCache, bool)
}
