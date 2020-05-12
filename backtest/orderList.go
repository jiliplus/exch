package backtest

import "github.com/jujili/exch"

// Order 是 exch.Order 的复刻
// 利用 gob 两者不必是完全一直的
type order struct {
	Symbol      string
	AssetName   string
	CapitalName string
	// if ID is negative value, means unset
	// ID is time.Now().Unix()
	ID   int64
	Side exch.OrderSide
	Type exch.OrderType
	// 根据 Type 的不同，以下 3 个属性不是全都必须的
	AssetQuantity   float64
	AssetPrice      float64
	CapitalQuantity float64

	// 指向下一个挂单
	next *order
}

type orderList struct {
	head *exch.Order
}

func newOrderList() *orderList {
	return &orderList{
		// 因为根本不会查看 head 内部的数据
		// head 完全可以是一个空的
		head: &exch.Order{},
	}
}

func (l *orderList) push(a *exch.Order) {
	curr, next := o.head, o.head.Next
	return
}
