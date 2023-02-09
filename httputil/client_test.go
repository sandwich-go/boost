package httputil

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient(t *testing.T) {
	Convey("client", t, func() {
		const raw = `{ "foo": "bar" }`

		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(raw))
		}))
		defer ts.Close()
		res, err := Get(ts.URL)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)

		res, err = Post(ts.URL, "", nil)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)

		bs, err0 := Bytes(ts.URL)
		So(err0, ShouldBeNil)
		So(bs, ShouldNotBeNil)
		So(string(bs), ShouldEqual, raw)

		str, err1 := String(ts.URL)
		So(err1, ShouldBeNil)
		So(str, ShouldEqual, raw)

		var m map[string]string
		err2 := JSON(ts.URL, &m)
		So(err2, ShouldBeNil)
		So(len(m), ShouldEqual, 1)
		So(m["foo"], ShouldEqual, "bar")
	})
}
