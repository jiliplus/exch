package backtest

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_NextIDFunc(t *testing.T) {
	Convey("测试 NextID", t, func() {
		nextID := NextIDFunc()
		for i := 1; i < 10; i++ {
			So(nextID(), ShouldEqual, i)
		}
	})
}
