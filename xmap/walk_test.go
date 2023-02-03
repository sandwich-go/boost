package xmap

import (
	. "github.com/smartystreets/goconvey/convey"
	"sort"
	"strconv"
	"testing"
)

func TestMap(t *testing.T) {
	Convey("walk", t, func() {
		var n = 100
		var keys = make([]string, 0, n)
		var tm = make(map[string]string)
		for i := 0; i < n; i++ {
			v := strconv.FormatInt(int64(i), 10)
			tm[v] = v
			keys = append(keys, v)
		}
		sort.Strings(keys)

		var dest = make([]string, 0, n)
		WalkMapDeterministic(tm, func(k string, v string) bool {
			dest = append(dest, k)
			return true
		})

		for k, v := range keys {
			So(v, ShouldEqual, dest[k])
		}
	})
}
