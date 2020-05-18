package backtest

import (
	"testing"

	"github.com/jujili/exch"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_orderList_push(t *testing.T) {
	Convey("orderList.push", t, func() {
		enc := exch.EncFunc()
		dec := decOrderFunc()
		de := func(i interface{}) *order {
			return dec(enc(i))
		}
		BtcUsdtOrder := exch.NewOrder("BTCUSDT", "BTC", "USDT")
		lb1 := de(BtcUsdtOrder.With(exch.Limit(exch.BUY, 100, 100000)))
		ol := newOrderList()
		ol.push(lb1)
		Convey("ol 的 head 后面就是 lb1", func() {
			So(ol.head.next, ShouldResemble, lb1)
		})
		mb1 := de(BtcUsdtOrder.With(exch.Market(exch.BUY, 100000)))
		ol.push(mb1)
		Convey("插入市价单后，ol 的 head 后面就是 mb1", func() {
			So(ol.head.next, ShouldResemble, mb1)
		})
		mb2 := de(BtcUsdtOrder.With(exch.Market(exch.BUY, 100000)))
		mb2.ID++
		ol.push(mb2)
		Convey("插入新的市价单后，mb1 的后面是 mb2", func() {
			So(mb1.next, ShouldResemble, mb2)
		})
		temp := *lb1
		temp.AssetPrice -= 10000
		lb2 := &temp
		ol.push(lb2)
		Convey("插入更低的限价买入单后，lb2 应该在最后", func() {
			So(lb1.next, ShouldEqual, lb2)
			So(lb2.next, ShouldBeNil)
		})
		Convey("整个 list 的顺序是", func() {
			So(ol.head.next, ShouldResemble, mb1)
			So(mb1.next, ShouldResemble, mb2)
			So(mb2.next, ShouldResemble, lb1)
			So(lb1.next, ShouldResemble, lb2)
			So(lb2.next, ShouldBeNil)
		})
	})
}

func Test_orderList_pop(t *testing.T) {
	Convey("orderList.pop", t, func() {
		enc := exch.EncFunc()
		dec := decOrderFunc()
		de := func(i interface{}) *order {
			return dec(enc(i))
		}
		BtcUsdtOrder := exch.NewOrder("BTCUSDT", "BTC", "USDT")
		lb1 := de(BtcUsdtOrder.With(exch.Limit(exch.BUY, 100, 100000)))
		ol := newOrderList()
		ol.push(lb1)
		mb1 := de(BtcUsdtOrder.With(exch.Market(exch.BUY, 100000)))
		ol.push(mb1)
		mb2 := de(BtcUsdtOrder.With(exch.Market(exch.BUY, 100000)))
		mb2.ID++
		ol.push(mb2)
		temp := *lb1
		temp.AssetPrice -= 10000
		lb2 := &temp
		ol.push(lb2)
		Convey("整个 list 的顺序是", func() {
			So(ol.head.next, ShouldResemble, mb1)
			So(mb1.next, ShouldResemble, mb2)
			So(mb2.next, ShouldResemble, lb1)
			So(lb1.next, ShouldResemble, lb2)
			So(lb2.next, ShouldBeNil)
		})
		order := ol.pop()
		Convey("应该是 mb1", func() {
			So(order, ShouldResemble, mb1)
		})
		order = ol.pop()
		Convey("应该是 mb2", func() {
			So(order, ShouldResemble, mb2)
		})
		order = ol.pop()
		Convey("应该是 lb1", func() {
			So(order, ShouldResemble, lb1)
		})
		order = ol.pop()
		Convey("应该是 lb2", func() {
			So(order, ShouldResemble, lb2)
		})
		order = ol.pop()
		Convey("应该是 nil", func() {
			So(order, ShouldBeNil)
		})
	})
}

func Test_orderList_canMatch(t *testing.T) {
	Convey("orderList.canMatch", t, func() {
		enc := exch.EncFunc()
		dec := decOrderFunc()
		de := func(i interface{}) *order {
			return dec(enc(i))
		}
		BtcUsdtOrder := exch.NewOrder("BTCUSDT", "BTC", "USDT")
		ol := newOrderList()
		Convey("空的 orderList 不会匹配", func() {
			So(ol.canMatch(0), ShouldBeFalse)
		})
		Convey("市价 BUY 单，总是能够匹配成功", func() {
			mb := de(BtcUsdtOrder.With(exch.Market(exch.BUY, 100000)))
			ol.push(mb)
			So(ol.canMatch(0), ShouldBeTrue)
		})
		Convey("限价 BUY 单", func() {
			lb := de(BtcUsdtOrder.With(exch.Limit(exch.BUY, 100, 100000)))
			ol.push(lb)
			price := lb.AssetPrice
			Convey("对相等或更低的价格**能够**匹配", func() {
				So(ol.canMatch(price), ShouldBeTrue)
				So(ol.canMatch(price*0.99), ShouldBeTrue)
			})
			Convey("对更高的价格**不能够**匹配", func() {
				So(ol.canMatch(price*1.01), ShouldBeFalse)
			})
		})
		Convey("限价 SELL 单", func() {
			ls := de(BtcUsdtOrder.With(exch.Limit(exch.SELL, 100, 100000)))
			ol.push(ls)
			price := ls.AssetPrice
			Convey("对相等或更高的价格**能够**匹配", func() {
				So(ol.canMatch(price), ShouldBeTrue)
				So(ol.canMatch(price*1.01), ShouldBeTrue)
			})
			Convey("对更低的价格**不能够**匹配", func() {
				So(ol.canMatch(price*0.99), ShouldBeFalse)
			})
		})
	})
}
