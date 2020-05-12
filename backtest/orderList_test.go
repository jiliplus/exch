package backtest

import (
	"testing"
	"time"

	"github.com/jujili/exch"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_DecOrderFunc(t *testing.T) {
	Convey("反向序列化 order", t, func() {
		asset := "BTC"
		capital := "USDT"
		assetQuantity := 100.0
		assetPrice := 10000.0
		order := exch.NewOrder(asset+capital, asset, capital)
		Convey("Limit", func() {
			source := order.With(exch.Limit(exch.BUY, assetQuantity, assetPrice))
			enc := exch.EncFunc()
			dec := decOrderFunc()
			actual := dec(enc(source))
			Convey("具体的值，应该相同", func() {
				// REVIEW: 这里的比较方式太 low 了
				So(actual.ID, ShouldEqual, source.ID)
				So(actual.AssetName, ShouldEqual, source.AssetName)
				So(actual.CapitalName, ShouldEqual, source.CapitalName)
				So(actual.Side, ShouldEqual, source.Side)
				So(actual.Type, ShouldEqual, source.Type)
				So(actual.AssetQuantity, ShouldEqual, source.AssetQuantity)
				So(actual.AssetPrice, ShouldEqual, source.AssetPrice)
				So(actual.CapitalQuantity, ShouldEqual, source.CapitalQuantity)
			})
		})
	})
}

func Test_Order_IsLessThan(t *testing.T) {
	enc := exch.EncFunc()
	dec := decOrderFunc()
	de := func(i interface{}) *order {
		return dec(enc(i))
	}
	Convey("Order less function", t, func() {
		BtcUsdtOrder := exch.NewOrder("BTCUSDT", "BTC", "USDT")
		Convey("nil 的 order 会返回 false", func() {
			var nilOrder *order
			So(nilOrder.isLessThan(nil), ShouldBeFalse)
		})
		Convey("比较不同 side 的 order 会 panic", func() {
			lb := de(BtcUsdtOrder.With(exch.Limit(exch.BUY, 100, 100000)))
			ms := de(BtcUsdtOrder.With(exch.Market(exch.SELL, 100)))
			So(func() { lb.isLessThan(ms) }, ShouldPanic)
		})
		Convey("BUY side 的 order", func() {
			mb0 := de(BtcUsdtOrder.With(exch.Market(exch.BUY, 10000)))
			temp := *mb0
			temp.ID++
			mb1 := &temp
			lb0 := de(BtcUsdtOrder.With(exch.Limit(exch.BUY, 100, 110000)))
			lb1 := de(BtcUsdtOrder.With(exch.Limit(exch.BUY, 100, 100000)))
			Convey("同为 MARKET 类型，则按照 ID 升序排列", func() {
				So(mb0.isLessThan(mb1), ShouldBeTrue)
				So(mb1.isLessThan(mb0), ShouldBeFalse)
			})
			Convey("同为 LIMIT 类型，则按照 AssetPrice 降序排列", func() {
				So(lb0.isLessThan(lb1), ShouldBeTrue)
				So(lb1.isLessThan(lb0), ShouldBeFalse)
			})
			Convey("MARKET 永远排在 LIMIT 前面", func() {
				So(mb0.isLessThan(lb0), ShouldBeTrue)
				So(mb1.isLessThan(lb1), ShouldBeTrue)
				So(lb0.isLessThan(mb0), ShouldBeFalse)
				So(lb1.isLessThan(mb1), ShouldBeFalse)
			})
		})
		Convey("SELL side 的 order", func() {
			ms0 := de(BtcUsdtOrder.With(exch.Market(exch.SELL, 100)))
			temp := *ms0
			temp.ID++
			ms1 := &temp
			ls0 := de(BtcUsdtOrder.With(exch.Limit(exch.SELL, 100, 100000)))
			ls1 := de(BtcUsdtOrder.With(exch.Limit(exch.SELL, 100, 110000)))
			Convey("同为 MARKET 类型，则按照 ID 升序排列", func() {
				So(ms0.isLessThan(ms1), ShouldBeTrue)
				So(ms1.isLessThan(ms0), ShouldBeFalse)
			})
			Convey("同为 LIMIT 类型，则按照 AssetPrice 升序排列", func() {
				So(ls0.isLessThan(ls1), ShouldBeTrue)
				So(ls1.isLessThan(ls0), ShouldBeFalse)
			})
			Convey("MARKET 永远排在 LIMIT 前面", func() {
				So(ms0.isLessThan(ls0), ShouldBeTrue)
				So(ms1.isLessThan(ls1), ShouldBeTrue)
				So(ls0.isLessThan(ms0), ShouldBeFalse)
				So(ls1.isLessThan(ms1), ShouldBeFalse)
			})
		})
	})
}

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

func Test_order_match(t *testing.T) {
	Convey("order.match", t, func() {
		enc := exch.EncFunc()
		dec := decOrderFunc()
		de := func(i interface{}) *order {
			return dec(enc(i))
		}
		BtcUsdtOrder := exch.NewOrder("BTCUSDT", "BTC", "USDT")
		Convey("市价单以 tick 的价格撮合", func() {
			mb := de(BtcUsdtOrder.With(exch.Market(exch.BUY, 100000)))
			Convey("order 的金额小于 tick ", func() {
				price := 10000.
				tick := exch.NewTick(1, time.Now(), price, mb.CapitalQuantity/price*10)
				delta := mb.AssetQuantity / tick.Price
				So(delta, ShouldNotEqual, 0)
				newVolume := tick.Volume - delta
				mb.match(tick)
				// as := mb.match(tick)
				So(tick.Volume, ShouldEqual, newVolume)
				// Convey("Asset 的改变量", func() {
				// So(len(as), ShouldEqual, 2)
				// })
			})
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
		// ol := newOrderList()
		// ol.push(lb1)
		// mb1 := de(BtcUsdtOrder.With(exch.Market(exch.BUY, 100000)))
		// ol.push(mb1)
		// mb2 := de(BtcUsdtOrder.With(exch.Market(exch.BUY, 100000)))
		// mb2.ID++
		// ol.push(mb2)
		// temp := *lb1
		// temp.AssetPrice -= 10000
		// lb2 := &temp
		// ol.push(lb2)
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
