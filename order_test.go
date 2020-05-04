package exch

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_DecOrderFunc(t *testing.T) {
	Convey("反向序列化 Tick", t, func() {
		expected := NewOrder(
			"BTCUSDT",
			"BTC",
			"USDT",
			-1,
			Buy,
			100,
			10000,
			-1,
		)
		enc := EncFunc()
		dec := DecOrderFunc()
		actual := dec(enc(expected))
		Convey("指针指向的对象应该不同", func() {
			So(actual, ShouldNotEqual, expected)
			Convey("具体的值，应该相同", func() {
				So(*actual, ShouldResemble, *expected)
			})
		})
	})
}
