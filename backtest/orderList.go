package backtest

import (
	"bytes"
	"encoding/gob"

	"github.com/jujili/exch"
)

// Order 是 exch.Order 的复刻
// 利用 gob 两者不必是完全一直的
type order struct {
	exch.Order
	// 指向下一个挂单
	next *order
}

type orderList struct {
	head *order
}

func newOrderList() *orderList {
	return &orderList{
		// 因为根本不会查看 head 内部的数据
		// head 完全可以是一个空的
		head: &order{},
	}
}

// DecOrderFunc 返回的函数会把序列化成 []byte 的 Order 值转换回来
func decOrderFunc() func(bs []byte) *order {
	var buf bytes.Buffer
	dec := gob.NewDecoder(&buf)
	return func(bs []byte) *order {
		buf.Reset()
		buf.Write(bs)
		var order order
		dec.Decode(&order)
		return &order
	}
}

// isLessThan return true if o < a
// REVIEW: 当 Order 的 Type 增加以后，这个方法会爆炸。
func (o *order) isLessThan(a *order) bool {
	if o == nil || a == nil {
		return false
	}
	if o.Side != a.Side {
		panic("only compare with the same side")
	}
	if o.Type != a.Type {
		return o.Type < a.Type
	}
	switch o.Type {
	case exch.MARKET:
		return o.ID < a.ID
	case exch.LIMIT:
		return (o.AssetPrice == a.AssetPrice && o.ID < a.ID) ||
			o.sidePrice() < a.sidePrice()
	default:
		panic("现在只能处理 limit 和 market 类型。")
	}
}

func (o *order) sidePrice() float64 {
	return float64(o.Side) * o.AssetPrice
}

func (l *orderList) push(a *order) {
	curr, next := l.head, l.head.next
	for next.isLessThan(a) {
		curr, next = next, next.next
	}
	curr.next, a.next = a, next
}

func (l *orderList) pop() *order {
	if l.head.next == nil {
		return nil
	}
	curr, res := l.head, l.head.next
	curr.next, res.next = res.next, nil
	return res
}

// TODO: sellOrderList 和 buyOrderList 都存在 MARKET 订单，要如何处理。

func (l *orderList) canMatch(price float64) bool {
	if l.head.next == nil {
		return false
	}
	order := l.head.next
	if order.Type == exch.MARKET {
		return true
	}
	// order.Type == exch.LIMIT
	if order.Side == exch.BUY {
		return order.AssetPrice >= price
	}
	// order.Side == exch.BUY
	return order.AssetPrice <= price
}

// // 对于每个 tick 总是认为可以撮合成功，形成交易的。
// // 这里没有考虑手续费和滑点。
// func (o *order) match(tick *exch.Tick) []exch.Asset {
// 	var add, lost exch.Asset
// 	if o.Side == exch.BUY {
// 		add.Name = o.AssetName
// 		lost.Name = o.CapitalName
// 	} else {
// 		add.Name = o.CapitalName
// 		lost.Name = o.AssetName
// 	}
// 	if o.Type == exch.MARKET {
// 		// 市价单以 tick 的价格成交
// 		add.Free = math.Min(o.CapitalQuantity/tick.Price, tick.Volume)
// 		lost.Locked = -math.Min(o.CapitalQuantity, tick.Price*tick.Volume)
// 		tick.Volume -= add.Free
// 		o.CapitalQuantity += lost.Locked
// 		return []exch.Asset{add, lost}
// 	}
// 	if o.Side == exch.BUY {
// 		if o.Type == exch.LIMIT {
// 			// 限价单以 order 的价格成交
// 			if tick.Price <= o.AssetPrice {
// 				if tick.Volume >= o.AssetQuantity {
// 					add.Free = o.AssetQuantity
// 					lost.Locked = -o.AssetPrice * add.Free
// 					tick.Volume -= add.Free
// 					// o 会被丢弃，无需对其进行修改
// 				} else {
// 					add.Free = tick.Volume
// 					lost.Locked = -o.AssetPrice * tick.Volume
// 					tick.Volume = 0
// 					// o 还要放回 orderList，所以需要对其进行修改
// 					o.AssetQuantity -= add.Free
// 				}
// 			}
// 		}
// 		return []exch.Asset{add, lost}
// 	}
// 	// o.Side == exch.SELL
// 	if o.Type == exch.LIMIT {
// 		// 限价单以 order 的价格成交
// 		if tick.Price >= o.AssetPrice {
// 			if tick.Volume >= o.AssetQuantity {
// 				add.Free = o.AssetQuantity
// 				lost.Locked = -o.AssetPrice * add.Free
// 				tick.Volume -= o.AssetQuantity
// 				// o 会被丢弃，无需对其进行修改
// 			} else {
// 				add.Free = tick.Volume
// 				lost.Locked = -o.AssetPrice * tick.Volume
// 				tick.Volume = 0
// 				// o 还要放回 orderList，所以需要对其进行修改
// 				o.AssetQuantity -= add.Free
// 			}
// 		}
// 	}
// 	return []exch.Asset{add, lost}
// }
