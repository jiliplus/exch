package backtest

import (
	"testing"
	"time"

	"github.com/jujili/exch"
	"github.com/prashantv/gostub"
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

func Test_order_isLessThan(t *testing.T) {
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
		Convey("现在只能比较 LIMIT 和 MARKET", func() {
			ms := de(BtcUsdtOrder.With(exch.Market(exch.SELL, 100)))
			ms.Type = exch.STOPloss
			ls := de(BtcUsdtOrder.With(exch.Limit(exch.SELL, 100, 100000)))
			ls.Type = exch.STOPloss
			Convey("强行比较会 panic", func() {
				So(func() {
					ms.isLessThan(ls)
				}, ShouldPanic)
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

func Test_order_canMatch(t *testing.T) {
	Convey("检测 order.canMatch", t, func() {
		Convey("nil.canMatch 会返回 false", func() {
			var nilOrder *order
			So(nilOrder.canMatch(1), ShouldBeFalse)
		})
		Convey("现在只能检测 LIMIT 和 MARKET 类型的 order", func() {
			order := &order{}
			order.Type = exch.OrderType(3)
			So(func() {
				order.canMatch(1)
			}, ShouldPanic)
		})
	})
}

func checkMatch(
	matchFunc func(order, exch.Tick) (order, exch.Tick, []exch.Asset),
	ao, eo order,
	at, et exch.Tick,
	ea, ec exch.Asset,
) {
	ao, at, as := matchFunc(ao, at)
	Convey("order 应该与预期相等", func() {
		So(ao, ShouldResemble, eo)
	})
	Convey("tick 应该与预期相等", func() {
		So(at, ShouldResemble, et)
	})
	aa, ac := as[0], as[1]
	Convey("asset 应该与预期相等", func() {
		So(aa, ShouldResemble, ea)
	})
	Convey("capital 应该与预期相等", func() {
		So(ac, ShouldResemble, ec)
	})
}

// de 方便生成 order
func de(i interface{}) *order {
	enc := exch.EncFunc()
	dec := decOrderFunc()
	return dec(enc(i))
}

func Test_matchMarket(t *testing.T) {
	Convey("matchMarket 撮合市价单", t, func() {
		BtcUsdtOrder := exch.NewOrder("BTCUSDT", "BTC", "USDT")
		//
		Convey("输入别的类型的 order 会 panic", func() {
			lb := de(BtcUsdtOrder.With(exch.Limit(exch.BUY, 100, 100000)))
			So(lb.Type, ShouldNotEqual, exch.MARKET)
			So(func() {
				var tk exch.Tick
				matchMarket(*lb, tk)
			}, ShouldPanic)
		})
		//
		eAsset := exch.NewAsset(BtcUsdtOrder.AssetName, 0, 0)
		eCapital := exch.NewAsset(BtcUsdtOrder.CapitalName, 0, 0)
		//
		Convey("BUY 时", func() {
			capitalQuantity := 10000.
			mb := de(BtcUsdtOrder.With(exch.Market(exch.BUY, capitalQuantity)))
			tk := exch.NewTick(0, time.Now(), 1000, 100)
			Convey("如果 tick.Volume*tick.Price < mb.CapitalQuantity", func() {
				tk.Volume = mb.CapitalQuantity / tk.Price * 0.5
				//
				et := tk
				et.Volume = 0
				//
				eo := *mb
				eo.CapitalQuantity = mb.CapitalQuantity / 2
				//
				eAsset.Free = tk.Volume
				eCapital.Locked = -mb.CapitalQuantity / 2
				checkMatch(matchMarket, *mb, eo, tk, et, eAsset, eCapital)
			})
			Convey("如果 tick.Volume*tick.Price = mb.CapitalQuantity", func() {
				tk.Volume = mb.CapitalQuantity / tk.Price
				//
				et := tk
				et.Volume = 0
				//
				eo := *mb
				eo.CapitalQuantity = 0
				//
				eAsset.Free = tk.Volume
				eCapital.Locked = -mb.CapitalQuantity
				checkMatch(matchMarket, *mb, eo, tk, et, eAsset, eCapital)
			})
			Convey("如果 tick.Volume*tick.Price > mb.CapitalQuantity", func() {
				tk.Volume = mb.CapitalQuantity / tk.Price * 2
				//
				et := tk
				et.Volume = tk.Volume / 2
				//
				eo := *mb
				eo.CapitalQuantity = 0
				//
				eAsset.Free = tk.Volume / 2
				eCapital.Locked = -mb.CapitalQuantity
				checkMatch(matchMarket, *mb, eo, tk, et, eAsset, eCapital)
			})
		})
		Convey("SELL 时", func() {
			assetQuantity := 100.
			ms := de(BtcUsdtOrder.With(exch.Market(exch.SELL, assetQuantity)))
			tk := exch.NewTick(0, time.Now(), 1000, 100)
			Convey("如果 tick.Volume < ms.AssetQuantity", func() {
				tk.Volume = ms.AssetQuantity * 0.75
				//
				et := tk
				et.Volume = 0
				//
				eo := *ms
				eo.AssetQuantity = ms.AssetQuantity - tk.Volume
				//
				eAsset.Locked = -tk.Volume
				eCapital.Free = tk.Volume * tk.Price
				checkMatch(matchMarket, *ms, eo, tk, et, eAsset, eCapital)
			})
			Convey("如果 tick.Volume = ms.AssetQuantity", func() {
				tk.Volume = ms.AssetQuantity
				//
				et := tk
				et.Volume = 0
				//
				eo := *ms
				eo.AssetQuantity = ms.AssetQuantity - tk.Volume
				//
				eAsset.Locked = -tk.Volume
				eCapital.Free = tk.Volume * tk.Price
				checkMatch(matchMarket, *ms, eo, tk, et, eAsset, eCapital)
			})
			Convey("如果 tick.Volume > ms.AssetQuantity", func() {
				tk.Volume = ms.AssetQuantity * 1.25
				//
				et := tk
				et.Volume = tk.Volume - ms.AssetQuantity
				//
				eo := *ms
				eo.AssetQuantity = 0
				//
				eAsset.Locked = -ms.AssetQuantity
				eCapital.Free = ms.AssetQuantity * tk.Price
				checkMatch(matchMarket, *ms, eo, tk, et, eAsset, eCapital)
			})
		})
	})
}

func Test_matchLimit(t *testing.T) {
	Convey("matchLimit 撮合限价单", t, func() {
		BtcUsdtOrder := exch.NewOrder("BTCUSDT", "BTC", "USDT")
		//
		Convey("输入别的类型的 order 会 panic", func() {
			lb := de(BtcUsdtOrder.With(exch.Market(exch.BUY, 100000)))
			So(lb.Type, ShouldNotEqual, exch.LIMIT)
			So(func() {
				var tk exch.Tick
				matchLimit(*lb, tk)
			}, ShouldPanic)
		})
		//
		eAsset := exch.NewAsset(BtcUsdtOrder.AssetName, 0, 0)
		eCapital := exch.NewAsset(BtcUsdtOrder.CapitalName, 0, 0)
		//
		Convey("SELL 时", func() {
			quantity, price := 10000., 100.
			ls := de(BtcUsdtOrder.With(exch.Limit(exch.SELL, quantity, price)))
			tk := exch.NewTick(0, time.Now(), 1000, 100)
			Convey("如果 tick.Price < ls.AssetPrice，则无法成交", func() {
				tk.Price = ls.AssetPrice / 2
				//
				et := tk
				//
				eo := *ls
				checkMatch(matchLimit, *ls, eo, tk, et, eAsset, eCapital)
			})
			Convey("如果 tick.Price = ls.AssetPrice，则可以成交", func() {
				tk.Price = ls.AssetPrice
				//
				Convey("如果 tick.Volume < ls.AssetQuantity", func() {
					tk.Volume = ls.AssetQuantity / 2
					//
					et := tk
					et.Volume = 0
					//
					eo := *ls
					eo.AssetQuantity = ls.AssetQuantity - tk.Volume
					//
					eAsset.Locked = -tk.Volume
					eCapital.Free = ls.AssetPrice * tk.Volume
					checkMatch(matchLimit, *ls, eo, tk, et, eAsset, eCapital)
				})
				Convey("如果 tick.Volume = ls.AssetQuantity", func() {
					tk.Volume = ls.AssetQuantity
					//
					et := tk
					et.Volume = 0
					//
					eo := *ls
					eo.AssetQuantity = 0
					//
					eAsset.Locked = -tk.Volume
					eCapital.Free = ls.AssetPrice * tk.Volume
					checkMatch(matchLimit, *ls, eo, tk, et, eAsset, eCapital)
				})
				Convey("如果 tick.Volume > ls.AssetQuantity", func() {
					tk.Volume = ls.AssetQuantity * 2
					//
					et := tk
					et.Volume = tk.Volume - ls.AssetQuantity
					//
					eo := *ls
					eo.AssetQuantity = 0
					//
					eAsset.Locked = -ls.AssetQuantity
					eCapital.Free = ls.AssetPrice * ls.AssetQuantity
					checkMatch(matchLimit, *ls, eo, tk, et, eAsset, eCapital)
				})
			})
			Convey("如果 tick.Price > ls.AssetPrice，则可以成交", func() {
				tk.Price = ls.AssetPrice * 2
				//
				Convey("如果 tick.Volume < ls.AssetQuantity", func() {
					tk.Volume = ls.AssetQuantity / 2
					//
					et := tk
					et.Volume = 0
					//
					eo := *ls
					eo.AssetQuantity = ls.AssetQuantity - tk.Volume
					//
					eAsset.Locked = -tk.Volume
					eCapital.Free = ls.AssetPrice * tk.Volume
					checkMatch(matchLimit, *ls, eo, tk, et, eAsset, eCapital)
				})
				Convey("如果 tick.Volume = ls.AssetQuantity", func() {
					tk.Volume = ls.AssetQuantity
					//
					et := tk
					et.Volume = 0
					//
					eo := *ls
					eo.AssetQuantity = 0
					//
					eAsset.Locked = -tk.Volume
					eCapital.Free = ls.AssetPrice * tk.Volume
					checkMatch(matchLimit, *ls, eo, tk, et, eAsset, eCapital)
				})
				Convey("如果 tick.Volume > ls.AssetQuantity", func() {
					tk.Volume = ls.AssetQuantity * 2
					//
					et := tk
					et.Volume = tk.Volume - ls.AssetQuantity
					//
					eo := *ls
					eo.AssetQuantity = 0
					//
					eAsset.Locked = -ls.AssetQuantity
					eCapital.Free = ls.AssetPrice * ls.AssetQuantity
					checkMatch(matchLimit, *ls, eo, tk, et, eAsset, eCapital)
				})
			})
		})
		Convey("BUY 时", func() {
			quantity, price := 10000., 100.
			lb := de(BtcUsdtOrder.With(exch.Limit(exch.BUY, quantity, price)))
			tk := exch.NewTick(0, time.Now(), 1000, 100)
			Convey("如果 tick.Price > lb.AssetPrice，则无法成交", func() {
				tk.Price = lb.AssetPrice * 2
				//
				et := tk
				//
				eo := *lb
				checkMatch(matchLimit, *lb, eo, tk, et, eAsset, eCapital)
			})
			Convey("如果 tick.Price = lb.AssetPrice，则可以成交", func() {
				tk.Price = lb.AssetPrice
				//
				Convey("如果 tick.Volume < lb.AssetQuantity", func() {
					tk.Volume = lb.AssetQuantity / 2
					//
					et := tk
					et.Volume = 0
					//
					eo := *lb
					eo.AssetQuantity = lb.AssetQuantity - tk.Volume
					//
					eAsset.Free = tk.Volume
					eCapital.Locked = -lb.AssetPrice * tk.Volume
					checkMatch(matchLimit, *lb, eo, tk, et, eAsset, eCapital)
				})
				Convey("如果 tick.Volume = lb.AssetQuantity", func() {
					tk.Volume = lb.AssetQuantity
					//
					et := tk
					et.Volume = 0
					//
					eo := *lb
					eo.AssetQuantity = 0
					//
					eAsset.Free = tk.Volume
					eCapital.Locked = -lb.AssetPrice * tk.Volume
					checkMatch(matchLimit, *lb, eo, tk, et, eAsset, eCapital)
				})
				Convey("如果 tick.Volume > lb.AssetQuantity", func() {
					tk.Volume = lb.AssetQuantity * 2
					//
					et := tk
					et.Volume = tk.Volume - lb.AssetQuantity
					//
					eo := *lb
					eo.AssetQuantity = 0
					//
					eAsset.Free = lb.AssetQuantity
					eCapital.Locked = -lb.AssetPrice * lb.AssetQuantity
					checkMatch(matchLimit, *lb, eo, tk, et, eAsset, eCapital)
				})
			})
			Convey("如果 tick.Price < lb.AssetPrice，则可以成交", func() {
				tk.Price = lb.AssetPrice / 2
				//
				Convey("如果 tick.Volume < lb.AssetQuantity", func() {
					tk.Volume = lb.AssetQuantity / 2
					//
					et := tk
					et.Volume = 0
					//
					eo := *lb
					eo.AssetQuantity = lb.AssetQuantity - tk.Volume
					//
					eAsset.Free = tk.Volume
					eCapital.Locked = -lb.AssetPrice * tk.Volume
					checkMatch(matchLimit, *lb, eo, tk, et, eAsset, eCapital)
				})
				Convey("如果 tick.Volume = lb.AssetQuantity", func() {
					tk.Volume = lb.AssetQuantity
					//
					et := tk
					et.Volume = 0
					//
					eo := *lb
					eo.AssetQuantity = 0
					//
					eAsset.Free = tk.Volume
					eCapital.Locked = -lb.AssetPrice * tk.Volume
					checkMatch(matchLimit, *lb, eo, tk, et, eAsset, eCapital)
				})
				Convey("如果 tick.Volume > lb.AssetQuantity", func() {
					tk.Volume = lb.AssetQuantity * 2
					//
					et := tk
					et.Volume = tk.Volume - lb.AssetQuantity
					//
					eo := *lb
					eo.AssetQuantity = 0
					//
					eAsset.Free = lb.AssetQuantity
					eCapital.Locked = -lb.AssetPrice * lb.AssetQuantity
					checkMatch(matchLimit, *lb, eo, tk, et, eAsset, eCapital)
				})
			})
		})
	})
}

func Test_order_match(t *testing.T) {
	Convey("测试 order.match", t, func() {
		BtcUsdtOrder := exch.NewOrder("BTCUSDT", "BTC", "USDT")
		//
		Convey("输入别的类型的 order 会 panic", func() {
			lb := de(BtcUsdtOrder.With(exch.Market(exch.BUY, 100000)))
			lb.Type = 3
			// NOTICE: order.match 匹配的返回扩大以后，需要修改这个断言。
			So(lb.Type, ShouldNotBeBetweenOrEqual, 1, 2)
			So(func() {
				var tk exch.Tick
				lb.match(tk)
			}, ShouldPanicWith, "现在只能处理 limit 和 market 类型")
		})
		Convey("MARKET 的 order.match 会调用 matchMarket", func() {
			hasCalled := false
			stubs := gostub.Stub(&matchMarket, func(o order, t exch.Tick) (order, exch.Tick, []exch.Asset) {
				hasCalled = true
				return o, t, nil
			})
			defer stubs.Reset()
			mb := de(BtcUsdtOrder.With(exch.Market(exch.BUY, 1000)))
			tk := exch.NewTick(0, time.Now(), 1000, 100)
			mb.match(tk)
			Convey("matchMarket 应该被调用了", func() {
				So(hasCalled, ShouldBeTrue)
			})
		})
		//
		Convey("LIMIT 的 order.match 会调用 matchLimit", func() {
			hasCalled := false
			stubs := gostub.Stub(&matchLimit, func(o order, t exch.Tick) (order, exch.Tick, []exch.Asset) {
				hasCalled = true
				return o, t, nil
			})
			defer stubs.Reset()
			ls := de(BtcUsdtOrder.With(exch.Limit(exch.SELL, 10000, 100)))
			tk := exch.NewTick(0, time.Now(), 1000, 100)
			ls.match(tk)
			Convey("matchLimit 应该被调用了", func() {
				So(hasCalled, ShouldBeTrue)
			})
		})
	})
}
