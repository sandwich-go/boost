package dns

import (
	"context"
	"github.com/sandwich-go/boost/z"
	"net"
	"strings"
	"sync"
	"time"
)

type dns struct {
	Resolver
	opts *Options
}

// New 创建 dns
func New(opts ...Option) DNS {
	cfg := NewOptions(opts...)
	d := &dns{opts: cfg}
	d.Resolver = d
	return d
}

// GetDialContext 获取 Dial 函数
func (d *dns) GetDialContext() func(ctx context.Context, network, address string) (net.Conn, error) {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		sepIndex := strings.LastIndex(addr, ":")
		ipAddrs, err := d.Resolver.LookupIPAddr(ctx, addr[:sepIndex])
		if err != nil {
			return nil, err
		}
		if len(ipAddrs) == 0 {
			return nil, ErrNotFound
		}
		index := 0
		if d.opts.GetPolicy() == PolicyRandom {
			index = int(z.FastRand()) % len(ipAddrs)
		}
		addr = ipAddrs[index].String() + addr[sepIndex:]
		return d.opts.GetDialer().DialContext(ctx, network, addr)
	}
}

// LookupIPAddr 通过主机搜索 ip 地址列表
func (d *dns) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	var start = time.Now()
	var cancel context.CancelFunc
	if _, ok := ctx.Deadline(); !ok {
		ctx, cancel = context.WithTimeout(ctx, d.opts.GetLookupTimeout())
	}
	defer func() {
		if cancel != nil {
			cancel()
		}
	}()

	ipAddrs, err := d.opts.GetResolver().LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
	}
	if len(ipAddrs) == 0 {
		return nil, ErrNotFound
	}
	// 成功则回调
	if d.opts.GetOnLookup() != nil {
		d.opts.GetOnLookup()(ctx, host, time.Since(start), ipAddrs)
	}
	return ipAddrs, nil
}

type cacheDns struct {
	*dns
	caches sync.Map
}

// NewCache 创建带缓存的 dns
func NewCache(opts ...Option) CacheDNS {
	dc := &cacheDns{dns: New(opts...).(*dns)}
	dc.dns.Resolver = dc
	return dc
}

// LookupIPAddr 通过主机搜索 ip 地址列表
func (dc *cacheDns) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	ipCache, exists := dc.Get(host)
	if exists {
		// 如果创建时间小于0，表示永久有效
		// 如果在有效期内，直接返回
		if ipCache.CreatedAt.IsZero() || ipCache.CreatedAt.Add(dc.opts.GetTTL()).After(time.Now()) {
			return ipCache.IPAddrs, nil
		}
	}
	ipAddrs, err := dc.dns.LookupIPAddr(ctx, host)
	if err != nil {
		return nil, err
	}
	if dc.opts.GetTTL() > 0 {
		dc.Set(host, IpCache{IPAddrs: ipAddrs, CreatedAt: time.Now()})
	} else {
		dc.Set(host, IpCache{IPAddrs: ipAddrs})
	}
	return ipAddrs, nil
}

func (dc *cacheDns) Set(host string, cache IpCache) {
	dc.caches.Store(host, cache)
}

func (dc *cacheDns) Remove(host string) {
	dc.caches.Delete(host)
}

func (dc *cacheDns) Get(host string) (IpCache, bool) {
	c, ok := dc.caches.Load(host)
	if !ok {
		return EmptyIpCache, false
	}
	return c.(IpCache), true
}
