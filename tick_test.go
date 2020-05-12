package exch

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_DecTickFunc(t *testing.T) {
	Convey("反向序列化 Tick", t, func() {
		expected := &Tick{
			// Exchange: OKEX,
			// Symbol:   "BTCUSDT",
			ID:     110,
			Date:   time.Now(),
			Price:  1222,
			Volume: 100,
			// Type:   "Type",
		}
		enc := EncFunc()
		dec := DecTickFunc()
		actual := dec(enc(expected))
		Convey("指针指向的对象应该不同", func() {
			So(actual, ShouldNotEqual, expected)
			Convey("具体的值，应该相同", func() {
				So(actual.Date.Equal(expected.Date), ShouldBeTrue)
				actual.Date = expected.Date
				// 没有上面两行，直接使用下面的判断语句会报错，
				So(*actual, ShouldResemble, *expected)
			})
		})
	})
}
