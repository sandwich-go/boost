package graph

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLine(t *testing.T) {
	Convey("test line", t, func() {
		s := P[int8](1, 1)
		e := P[int8](4, 4)
		l0, l1 := L[int8](s, e), L[int8](s, e)
		So(l0.Equals(l1), ShouldBeTrue)

		So(l0.Start().Equals(s), ShouldBeTrue)
		So(l0.End().Equals(e), ShouldBeTrue)

		t.Log(l0)

		var l2 Line[int64]
		var l2PointNum int
		l2.RangePoints(func(p Point[int64]) bool {
			l2PointNum++
			return true
		})
		So(l2PointNum, ShouldBeZeroValue)

		var l3 = L[int8](s, s)
		var l3Points []Point[int8]
		l3.RangePoints(func(p Point[int8]) bool {
			l3Points = append(l3Points, p)
			return true
		})
		So(len(l3Points), ShouldEqual, 1)
		for _, p := range l3Points {
			So(p.Equals(s), ShouldBeTrue)
		}

		var l0Points []Point[int8]
		l0.RangePoints(func(p Point[int8]) bool {
			l0Points = append(l0Points, p)
			return true
		})
		So(len(l0Points), ShouldEqual, 4)

		l4 := L[int8](e, s)
		var l4Points []Point[int8]
		l4.RangePoints(func(p Point[int8]) bool {
			l4Points = append(l4Points, p)
			return true
		})
		So(len(l4Points), ShouldEqual, 4)
	})

	// 覆盖 RangePoints 各分支：xx>xy 横线 / 反向 / with 早退
	Convey("RangePoints 横线 / 反向 / with 早退", t, func() {
		// 横线 (1,1)→(5,1)：xx=4 >> xy=0，走 'xx >= xy' 的 else 分支
		// （line 70-86）+ positive=true 子分支
		var ptsHor []Point[int8]
		L[int8](P[int8](1, 1), P[int8](5, 1)).RangePoints(func(p Point[int8]) bool {
			ptsHor = append(ptsHor, p)
			return true
		})
		So(len(ptsHor), ShouldEqual, 5)

		// 反向横线 (5,1)→(1,1)：走 line 70-86 的 positive=false 子分支
		var ptsHorRev []Point[int8]
		L[int8](P[int8](5, 1), P[int8](1, 1)).RangePoints(func(p Point[int8]) bool {
			ptsHorRev = append(ptsHorRev, p)
			return true
		})
		So(len(ptsHorRev), ShouldEqual, 5)

		// 反向竖线 (1,5)→(1,1)：走 line 51-68 的 positive=false 子分支
		var ptsVerRev []Point[int8]
		L[int8](P[int8](1, 5), P[int8](1, 1)).RangePoints(func(p Point[int8]) bool {
			ptsVerRev = append(ptsVerRev, p)
			return true
		})
		So(len(ptsVerRev), ShouldEqual, 5)

		// with 返 false 早退（line 43-44 + line 65 / 83）
		var earlyCount int
		L[int8](P[int8](1, 1), P[int8](10, 1)).RangePoints(func(p Point[int8]) bool {
			earlyCount++
			return earlyCount < 3 // 第 3 次返 false
		})
		So(earlyCount, ShouldEqual, 3)
	})
}
