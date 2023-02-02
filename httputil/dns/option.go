package dns

import (
	"context"
	"net"
	"time"
)

type OnLookup func(ctx context.Context, host string, cost time.Duration, ipAddrs []net.IPAddr)

// defaultDialer 默认的拨号器
var defaultDialer = &net.Dialer{
	Timeout:   30 * time.Second,
	KeepAlive: 30 * time.Second,
}

//go:generate optiongen --option_with_struct_name=false --new_func=NewOptions --xconf=true --empty_composite_nil=true --usage_tag_name=usage
func OptionsOptionDeclareWithDefault() interface{} {
	return map[string]interface{}{
		"TTL":           time.Duration(0),                // @MethodComment(Cache过期ttl)
		"Dialer":        (Dialer)(defaultDialer),         // @MethodComment(拨号器)
		"Resolver":      (Resolver)(net.DefaultResolver), // @MethodComment(解析器)
		"Policy":        Policy(PolicyFirst),             // @MethodComment(拨号策略)
		"LookupTimeout": time.Duration(3 * time.Second),  // @MethodComment(搜索超时)
		"OnLookup":      OnLookup(nil),                   // @MethodComment(当成功搜索)
	}
}
