package backtest

import (
	"bytes"
	"encoding/gob"
	"math"

	"github.com/jujili/exch"
)

// Order 是 exch.Order 的复刻
// 利用 gob 两者不必是完全一直的
type order struct {
	exch.Order
	// 指向下一个挂单
	next *order
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
// otherwise return false
func (o *order) isLessThan(a *order) bool {
	if o == nil {
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

// canMatch return true if o < a
// otherwise return false
func (o *order) canMatch(price float64) bool {
	if o == nil {
		return false
	}
	switch o.Type {
	case exch.MARKET:
		// MARKET 总是可以撮合上
		return true
	case exch.LIMIT:
		return o.sidePrice() <= float64(o.Side)*price
	default:
		panic("现在只能处理 limit 和 market 类型。")
	}
}

func (o *order) sidePrice() float64 {
	return float64(o.Side) * o.AssetPrice
}

// 对于每个 tick 总是认为可以撮合成功，形成交易的。
// 这里没有考虑手续费和滑点。
// match 前需要使用 canMatch 进行检查， match 内就不再检查了
// TODO: 测试空订单
func (o order) match(tick exch.Tick) (order, exch.Tick, []exch.Asset) {
	switch o.Type {
	case exch.MARKET:
		return matchMarket(o, tick)
	case exch.LIMIT:
		return matchLimit(o, tick)
	default:
		panic("现在只能处理 limit 和 market 类型")
	}
}

var matchMarket = func(o order, t exch.Tick) (order, exch.Tick, []exch.Asset) {
	var asset, capital exch.Asset
	asset.Name = o.AssetName
	capital.Name = o.CapitalName
	if o.Type != exch.MARKET {
		panic("order.Type should be exch.MARKET")
	}
	if o.Side == exch.SELL {
		diff := math.Min(o.AssetQuantity, t.Volume)
		asset.Locked = -diff
		capital.Free = t.Price * diff
		t.Volume -= diff
		o.AssetQuantity -= diff
	} else {
		diff := math.Min(o.CapitalQuantity, t.Volume*t.Price)
		asset.Free = diff / t.Price
		capital.Locked = -diff
		t.Volume -= diff / t.Price
		o.CapitalQuantity -= diff
	}
	return o, t, []exch.Asset{asset, capital}
}

var matchLimit = func(o order, t exch.Tick) (order, exch.Tick, []exch.Asset) {
	var asset, capital exch.Asset
	asset.Name = o.AssetName
	capital.Name = o.CapitalName
	if o.Type != exch.LIMIT {
		panic("order.Type should be exch.LIMIT")
	}
	if float64(o.Side)*t.Price < o.sidePrice() {
		return o, t, []exch.Asset{asset, capital}
	}
	// 处于谨慎的态度，以 o.AssetPrice 的价格成交
	var diff float64
	diff = math.Min(o.AssetQuantity, t.Volume)
	if o.Side == exch.SELL {
		asset.Locked = -diff
		capital.Free = diff * o.AssetPrice
	} else {
		asset.Free = diff
		capital.Locked = -diff * o.AssetPrice
	}
	t.Volume -= diff
	o.AssetQuantity -= diff
	return o, t, []exch.Asset{asset, capital}
}

func (o *order) pend2Lock() exch.Asset {
	switch o.Type {
	case exch.MARKET:
		return pendMarket(*o)
	case exch.LIMIT:
		return pendLimit(*o)
	default:
		panic("现在只能处理 limit 和 market 类型")
	}
}

var pendMarket = func(o order) exch.Asset {
	var res exch.Asset
	if o.Type != exch.MARKET {
		panic("pendMarket 应该输入 MARKET 类型的 order")
	}
	if o.Side == exch.BUY {
		res.Name = o.CapitalName
		res.Free = -o.CapitalQuantity
		res.Locked = o.CapitalQuantity
	} else { // o.Side == exch.SELL
		res.Name = o.AssetName
		res.Free = -o.AssetQuantity
		res.Locked = o.AssetQuantity
	}
	return res
}

var pendLimit = func(o order) exch.Asset {
	var res exch.Asset
	if o.Type != exch.LIMIT {
		panic("pendLimit 应该输入 LIMIT 类型的 order")
	}
	if o.Side == exch.BUY {
		res.Name = o.CapitalName
		total := o.AssetQuantity * o.AssetPrice
		res.Free = -total
		res.Locked = total
	} else { // o.Side == exch.SELL
		res.Name = o.AssetName
		res.Free = -o.AssetQuantity
		res.Locked = o.AssetQuantity
	}
	return res
}

func (o *order) cancel2Free() exch.Asset {
	switch o.Type {
	case exch.MARKET:
		return cancelMarket(*o)
	case exch.LIMIT:
		return cancelLimit(*o)
	default:
		panic("现在只能处理 limit 和 market 类型")
	}
}

var cancelMarket = func(o order) exch.Asset {
	var res exch.Asset
	if o.Type != exch.MARKET {
		panic("cancelMarket 应该输入 MARKET 类型的 order")
	}
	if o.Side == exch.BUY {
		res.Name = o.CapitalName
		res.Free = o.CapitalQuantity
		res.Locked = -o.CapitalQuantity
	} else { // o.Side == exch.SELL
		res.Name = o.AssetName
		res.Free = o.AssetQuantity
		res.Locked = -o.AssetQuantity
	}
	return res
}

var cancelLimit = func(o order) exch.Asset {
	var res exch.Asset
	if o.Type != exch.LIMIT {
		panic("cancelLimit 应该输入 LIMIT 类型的 order")
	}
	if o.Side == exch.BUY {
		res.Name = o.CapitalName
		total := o.AssetQuantity * o.AssetPrice
		res.Free = total
		res.Locked = -total
	} else { // o.Side == exch.SELL
		res.Name = o.AssetName
		res.Free = o.AssetQuantity
		res.Locked = -o.AssetQuantity
	}
	return res
}
