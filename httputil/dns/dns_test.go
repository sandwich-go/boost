package dns

import (
	"context"
	"errors"
	"fmt"
	"github.com/sandwich-go/boost/z"
	. "github.com/smartystreets/goconvey/convey"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

var (
	mockPort      = strconv.FormatInt(int64(z.FastRand()), 10)
	mockIPAddrs   = []net.IPAddr{{IP: []byte("a")}, {IP: []byte("b")}}
	mockConn      = &net.TCPConn{}
	errLookupFail = errors.New("lookup fail")
	errDailFail   = errors.New("dail fail")
	errTimeout    = errors.New("timeout")
)

type mockDialer struct{}

func (mockDialer) DialContext(_ context.Context, _, address string) (net.Conn, error) {
	for _, ipAddr := range mockIPAddrs {
		if strings.HasPrefix(address, ipAddr.String()) && strings.HasSuffix(address, mockPort) {
			return mockConn, nil
		}
	}
	return nil, errDailFail
}

type mockResolver struct{}

func (mockResolver) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	if len(host) == 0 {
		return nil, errLookupFail
	}
	host = strings.TrimSpace(host)
	if len(host) == 0 {
		return nil, nil
	}
	if ss, _ := strconv.ParseInt(host, 10, 64); ss > 0 {
		select {
		case <-ctx.Done():
			return nil, errTimeout
		}
	}
	return mockIPAddrs, nil
}

func TestDNS(t *testing.T) {
	Convey("dns", t, func() {
		var lookupSuccess bool
		d := New(
			WithLookupTimeout(10*time.Millisecond),
			WithDialer(mockDialer{}),
			WithResolver(mockResolver{}),
			WithOnLookup(func(ctx context.Context, host string, cost time.Duration, ipAddrs []net.IPAddr) {
				lookupSuccess = true
			}))
		for _, test := range []struct {
			host string
			err  error
		}{
			{host: "", err: errLookupFail},
			{host: " ", err: ErrNotFound},
			{host: "1", err: errTimeout},
			{host: "0"},
		} {
			ipAddrs, err := d.LookupIPAddr(context.Background(), test.host)
			if test.err != nil {
				So(err, ShouldNotBeNil)
				So(test.err, ShouldEqual, err)
				So(lookupSuccess, ShouldBeFalse)
			} else {
				So(err, ShouldBeNil)
				So(len(ipAddrs), ShouldEqual, len(mockIPAddrs))
				So(ipAddrs, ShouldResemble, mockIPAddrs)
				So(lookupSuccess, ShouldBeTrue)
			}
		}

		lookupSuccess = false

		dail := d.GetDialContext()
		for _, test := range []struct {
			host string
			err  error
		}{
			{host: fmt.Sprintf(":%s", mockPort), err: errLookupFail},
			{host: fmt.Sprintf(" :%s", mockPort), err: ErrNotFound},
			{host: fmt.Sprintf("1:%s", mockPort), err: errTimeout},
			{host: fmt.Sprintf("0:%s", mockPort)},
		} {
			conn, err := dail(context.Background(), "mock", test.host)
			if test.err != nil {
				So(err, ShouldNotBeNil)
				So(conn, ShouldBeNil)
				So(test.err, ShouldEqual, err)
				So(lookupSuccess, ShouldBeFalse)
			} else {
				So(err, ShouldBeNil)
				So(conn, ShouldNotBeNil)
				So(conn, ShouldEqual, mockConn)
				So(lookupSuccess, ShouldBeTrue)
			}
		}
	})
}

func TestCacheDNS(t *testing.T) {
	Convey("cache dns", t, func() {
		var lookupSuccess bool
		d := NewCache(
			WithLookupTimeout(10*time.Millisecond),
			WithDialer(mockDialer{}),
			WithResolver(mockResolver{}),
			WithOnLookup(func(ctx context.Context, host string, cost time.Duration, ipAddrs []net.IPAddr) {
				lookupSuccess = true
			}))
		for _, test := range []struct {
			host string
			err  error
		}{
			{host: "", err: errLookupFail},
			{host: " ", err: ErrNotFound},
			{host: "1", err: errTimeout},
			{host: "0"},
		} {
			ipAddrs, err := d.LookupIPAddr(context.Background(), test.host)
			if test.err != nil {
				So(err, ShouldNotBeNil)
				So(test.err, ShouldEqual, err)
				So(lookupSuccess, ShouldBeFalse)
			} else {
				So(err, ShouldBeNil)
				So(len(ipAddrs), ShouldEqual, len(mockIPAddrs))
				So(ipAddrs, ShouldResemble, mockIPAddrs)
				So(lookupSuccess, ShouldBeTrue)

				_, ok := d.Get(test.host)
				So(ok, ShouldBeTrue)
			}
		}

		lookupSuccess = false

		dail := d.GetDialContext()
		for _, test := range []struct {
			host string
			err  error
		}{
			{host: fmt.Sprintf(":%s", mockPort), err: errLookupFail},
			{host: fmt.Sprintf(" :%s", mockPort), err: ErrNotFound},
			{host: fmt.Sprintf("1:%s", mockPort), err: errTimeout},
			{host: fmt.Sprintf("0:%s", mockPort)},
		} {
			conn, err := dail(context.Background(), "mock", test.host)
			if test.err != nil {
				So(err, ShouldNotBeNil)
				So(conn, ShouldBeNil)
				So(test.err, ShouldEqual, err)
				So(lookupSuccess, ShouldBeFalse)
			} else {
				So(err, ShouldBeNil)
				So(conn, ShouldNotBeNil)
				So(conn, ShouldEqual, mockConn)
				So(lookupSuccess, ShouldBeFalse) // 说明没有真正的去 lookup，而是走的缓存
			}
		}
	})
}
