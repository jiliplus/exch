package exch

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_DecOrderFunc(t *testing.T) {
	Convey("反向序列化 Tick", t, func() {
		order := NewOrder("BTCUSDT", "BTC", "USDT")
		Convey("Limit", func() {
			expected := order.With(Limit(BUY, 100, 10000))
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
		Convey("Market BUY", func() {
			expected := order.With(Market(BUY, 10000))
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
		Convey("Market SELL", func() {
			expected := order.With(Market(SELL, 100))
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
	})
}

func Test_OrderType_String(t *testing.T) {
	Convey("", t, func() {

	})
}
