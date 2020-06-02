package backtest

import (
	"github.com/jujili/exch"
)

type orderList struct {
	head *order
}

func (ol orderList) String() string {
	if ol.head.next == nil {
		return "[[EMPTY orderList]]"
	}
	return ol.head.next.String()
}

func newOrderList() *orderList {
	return &orderList{
		// 因为根本不会查看 head 内部的数据
		// head 完全可以是一个空的
		head: &order{},
	}
}

func (l *orderList) push(a *order) exch.Asset {
	curr, next := l.head, l.head.next
	for next.isLessThan(a) {
		curr, next = next, next.next
	}
	curr.next, a.next = a, next
	return a.pend2Lock()
}

func (l *orderList) remove(a *order) exch.Asset {
	curr, next := l.head, l.head.next
	for next != nil && next.ID != a.ID {
		curr, next = next, next.next
	}
	if next == nil {
		return exch.NewAsset(a.AssetName, 0, 0)
	}
	curr.next = next.next
	return next.cancel2Free()
}

func (l *orderList) pop() *order {
	if l.head.next == nil {
		return nil
	}
	curr, res := l.head, l.head.next
	curr.next, res.next = res.next, nil
	return res
}

func (l *orderList) isEmpty() bool {
	return l.head.next == nil
}

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
	// 防止把 for 循环前的 order 添加进来了
	if order.Type != 0 {
		l.push(&order) // order 此时有可能是空订单
	}
	return res
}
