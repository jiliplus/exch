package exch

import (
	"fmt"
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
	Convey("测试 OrderType 的字符化", t, func() {
		tests := []struct {
			name     string
			t        OrderType
			expected string
		}{
			{"LIMIT", LIMIT, "LIMIT"},
			{"MARKET", MARKET, "MARKET"},
			{"STOP_LOSS", STOPloss, "STOP_LOSS"},
			{"STOP_LOSS_LIMIT", STOPlossLIMIT, "STOP_LOSS_LIMIT"},
			{"TAKE_PROFIT", TAKEprofit, "TAKE_PROFIT"},
			{"TAKE_PROFIT_LIMIT", TAKEprofitLIMIT, "TAKE_PROFIT_LIMIT"},
			{"LIMIT_MAKER", LIMITmaker, "LIMIT_MAKER"},
		}
		for _, tt := range tests {
			title := fmt.Sprintf("测试 %s", tt.name)
			Convey(title, func() {
				actual := tt.t.String()
				So(actual, ShouldEqual, tt.expected)
			})
		}
	})
	Convey("遇到未定义的 OrderType 会 panic", t, func() {
		So(func() { _ = OrderType(0).String() }, ShouldPanic)
	})
}

func Test_OrderSide_String(t *testing.T) {
	Convey("测试 OrderSide 的字符化", t, func() {
		tests := []struct {
			name     string
			t        OrderSide
			expected string
		}{
			{"BUY", BUY, "BUY"},
			{"SELL", SELL, "SELL"},
		}
		for _, tt := range tests {
			title := fmt.Sprintf("测试 %s", tt.name)
			Convey(title, func() {
				actual := tt.t.String()
				So(actual, ShouldEqual, tt.expected)
			})
		}
	})
	Convey("遇到未定义的 OrderSide 会 panic", t, func() {
		So(func() { _ = OrderSide(0).String() }, ShouldPanic)
	})
}

func Test_Order_IsLessThan(t *testing.T) {
	Convey("Order less function", t, func() {
		BtcUsdtOrder := NewOrder("BTCUSDT", "BTC", "USDT")
		Convey("比较不同 side 的 order 会 panic", func() {
			lb := BtcUsdtOrder.With(Limit(BUY, 100, 100000))
			ms := BtcUsdtOrder.With(Market(SELL, 100))
			So(func() { lb.IsLessThan(ms) }, ShouldPanic)
		})
		Convey("BUY side 的 order", func() {
			lb0 := BtcUsdtOrder.With(Limit(BUY, 100, 100000))
			var temp Order
			temp := *lb0
			lb1 := BtcUsdtOrder.With(Limit(BUY, 100, 100000))
			Convey("同为 limit 类型，则按照 ID 升序排列", func() {
				So(lb0.IsLessThan(lb1), ShouldBeTrue)
				So(lb1.IsLessThan(lb0), ShouldBeFalse)
			})
		})
	})
}
