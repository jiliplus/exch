package backtest

import (
	"bytes"
	"encoding/gob"

	"github.com/jujili/exch"
)

// Order 是 exch.Order 的复刻
// 利用 gob 两者不必是完全一直的
type order struct {
	ID          int64
	AssetName   string
	CapitalName string
	Side        exch.OrderSide
	Type        exch.OrderType
	// 根据 Type 的不同，以下 3 个属性不是全都必须的
	AssetQuantity   float64
	AssetPrice      float64
	CapitalQuantity float64
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
	if o == nil {
		return false
	}
	if o.Side != a.Side {
		panic("only compare with the same side")
	}
	// MARKET 订单按照先后顺序排列
	if o.Type == exch.MARKET && a.Type == exch.MARKET {
		return o.ID < a.ID
	}
	// MARKET 订单始终排在 LIMIT 订单前面
	if o.Type == exch.MARKET && a.Type == exch.LIMIT {
		return true
	}
	// MARKET 订单始终排在 LIMIT 订单前面
	if o.Type == exch.LIMIT && a.Type == exch.MARKET {
		return false
	}
	// o.Type == LIMIT && a.Type == LIMIT
	// BUY 订单是价高的先成交
	if o.Side == exch.BUY {
		return o.AssetPrice > a.AssetPrice ||
			(o.AssetPrice == a.AssetPrice && o.ID < a.ID)
	}
	// o.Side == SELL
	// SELL 订单是价低的先成交
	return o.AssetPrice < a.AssetPrice ||
		(o.AssetPrice == a.AssetPrice && o.ID < a.ID)
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

// 对于每个 tick 总是认为可以撮合成功，形成交易的。
// 这里没有考虑手续费和滑点。
func (o *order) match(tick *exch.Tick) []exch.Asset {
	if o.Type == exch.MARKET {
		if o.Side == exch.BUY {
			if o.AssetQuantity <= tick.Price*tick.Volume {

				tick = nil
			}

		}
	}
	return nil
}
