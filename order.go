package exch

import (
	"bytes"
	"encoding/gob"
)

// OrderSide is side of order
type OrderSide bool

// OrderType is side of order
const (
	Buy  OrderSide = false
	Sell OrderSide = true
)

// OrderType is type of order
type OrderType uint8

// OrderType is type of order
const (
	Limit OrderType = iota + 1
	Market
	// 以下类型是从 binance 抄过来的
	// https://binance-docs.github.io/apidocs/spot/cn/#trade
	StopLoss
	StopLossLimit
	TakeProfit
	TakeProfitLimit
	LimitMaker
)

// Order 是下单的格式
// TODO: 下单的 order 和挂单的 order 有什么区别吗？
type Order struct {
	Symbol      string
	AssetName   string
	CapitalName string
	ID          int64
	Side        OrderSide
	// 根据 Type 的不同，以下 3 个属性不是全都必须的
	AssetQuantity   float64
	AssetPrice      float64
	CapitalQuantity float64
}

// NewOrder returns a new order
// TODO: 解决 NewOrder 参数过多的问题
// 可以参考：
// https://github.com/xxjwxc/uber_go_guide_cn#%E5%8A%9F%E8%83%BD%E9%80%89%E9%A1%B9
func NewOrder(symbol, asset, capital string, ID int64, side OrderSide, assetQuantity, assetPrice, capitalQuantity float64) *Order {
	return &Order{
		Symbol:          symbol,
		AssetName:       asset,
		CapitalName:     capital,
		ID:              ID,
		Side:            side,
		AssetQuantity:   assetQuantity,
		AssetPrice:      assetPrice,
		CapitalQuantity: capitalQuantity,
	}
}

// DecOrderFunc 返回的函数会把序列化成 []byte 的 Order 值转换回来
func DecOrderFunc() func(bs []byte) *Order {
	var buf bytes.Buffer
	dec := gob.NewDecoder(&buf)
	return func(bs []byte) *Order {
		buf.Reset()
		buf.Write(bs)
		var order Order
		dec.Decode(&order)
		return &order
	}
}
