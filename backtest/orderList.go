package backtest

import "github.com/jujili/exch"

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
	return order.canMatch(price)
}

func (l *orderList) match(tick exch.Tick) []exch.Asset {
	res := make([]exch.Asset, 0, 16)
	var as []exch.Asset
	var order order
	for tick.Volume != 0 && l.canMatch(tick.Price) {
		order, tick, as = l.pop().match(tick)
		res = append(res, as...)
	}
	l.push(&order)
	return res
}
