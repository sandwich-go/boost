package httputil

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// 本文件目标：补 httputil 包内未覆盖的导出 API + 错误路径：
//   - New（之前 0%，自定义 client 工厂）
//   - SetDefaultTimeout（之前 0%）
//   - Request（之前 0%）
//   - (*Error).Error（之前 0%，String() 返回）
//   - httpClient.err（之前 0%，错误构造）
//   - Bytes / String / JSON 错误路径（之前各 ~75%，缺非 200 status 分支）
//
// 已有 client_test.go 覆盖 happy path（200 OK）。

func TestNew_CustomClient(t *testing.T) {
	Convey("New 用自定义 *http.Client 构造 Client", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("custom"))
		}))
		defer ts.Close()

		c := New(&http.Client{Timeout: 5 * time.Second})
		So(c, ShouldNotBeNil)

		got, err := c.Bytes(ts.URL)
		So(err, ShouldBeNil)
		So(string(got), ShouldEqual, "custom")
	})
}

func TestSetDefaultTimeout(t *testing.T) {
	Convey("SetDefaultTimeout 改 globalClient timeout", t, func() {
		// 保存原值恢复，避免污染其它测试
		orig := globalClient.Client.Timeout
		defer SetDefaultTimeout(orig)

		SetDefaultTimeout(7 * time.Second)
		So(globalClient.Client.Timeout, ShouldEqual, 7*time.Second)
	})
}

func TestRequest_CustomMethod(t *testing.T) {
	Convey("Request 发送任意自定义 method 请求", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Method", r.Method)
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		// PUT method（标准 Get/Post 不支持）
		req, err := http.NewRequest(http.MethodPut, ts.URL, strings.NewReader("body"))
		So(err, ShouldBeNil)

		resp, err := Request(req)
		So(err, ShouldBeNil)
		defer resp.Body.Close()

		So(resp.StatusCode, ShouldEqual, http.StatusOK)
		So(resp.Header.Get("X-Method"), ShouldEqual, "PUT")
	})
}

func TestError_ErrorMethod(t *testing.T) {
	Convey("(*Error).Error 返回 Message 字段", t, func() {
		e := &Error{
			Message:    "boom",
			StatusCode: 500,
			URL:        "http://example.com",
		}
		So(e.Error(), ShouldEqual, "boom")
	})
}

func TestBytes_NonOKStatus(t *testing.T) {
	Convey("Bytes 在非 200 状态下返回 *Error", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer ts.Close()

		_, err := Bytes(ts.URL)
		So(err, ShouldNotBeNil)
		// 验证 err 是 *Error 类型且 StatusCode / URL 透传
		errImpl, ok := err.(*Error)
		So(ok, ShouldBeTrue)
		So(errImpl.StatusCode, ShouldEqual, 500)
		So(errImpl.URL, ShouldContainSubstring, ts.URL)
		// 默认 message 含 status code 信息
		So(errImpl.Message, ShouldContainSubstring, "500")
	})
}

func TestString_NonOKStatus(t *testing.T) {
	Convey("String 在非 200 状态下返回 *Error", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer ts.Close()

		got, err := String(ts.URL)
		So(err, ShouldNotBeNil)
		So(got, ShouldEqual, "")
	})
}

func TestJSON_ErrorPaths(t *testing.T) {
	Convey("JSON 在非 200 状态下返回 *Error", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
		}))
		defer ts.Close()

		var v map[string]string
		err := JSON(ts.URL, &v)
		So(err, ShouldNotBeNil)
	})

	Convey("JSON 在 syntax error 时包成 *Error 返回", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			// 故意写非法 JSON
			_, _ = w.Write([]byte(`{not valid json}`))
		}))
		defer ts.Close()

		var v map[string]string
		err := JSON(ts.URL, &v)
		So(err, ShouldNotBeNil)
		// 实现：json.SyntaxError 被包成 *Error，message="JSON syntax error at <url>"
		errImpl, ok := err.(*Error)
		So(ok, ShouldBeTrue)
		So(errImpl.Message, ShouldContainSubstring, "JSON syntax error")
	})
}

func TestBytes_GetError(t *testing.T) {
	Convey("Bytes 在 Get 失败（无效 URL）时返回非 nil error", t, func() {
		// 不存在的 host，DNS 应该失败
		_, err := Bytes("http://this-domain-should-not-exist-aaaa-bbbb.invalid/")
		So(err, ShouldNotBeNil)
	})
}
