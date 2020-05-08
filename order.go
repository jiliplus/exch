package exch

import (
	"bytes"
	"encoding/gob"
	"time"
)

// OrderSide is side of order
type OrderSide uint8

// OrderType is side of order
const (
	BUY OrderSide = iota + 1
	SELL
)

func (t OrderSide) String() string {
	switch t {
	case BUY:
		return "BUY"
	case SELL:
		return "SELL"
	default:
		panic("meet UNKNOWN Order Side")
	}
}

// OrderType is type of order
type OrderType uint8

// OrderType is type of order
// 类型值从 iota+1 也就是 1 开始
// 是为了避开默认的 0 值
const (
	LIMIT OrderType = iota + 1
	MARKET
	// 以下类型是从 binance 抄过来的
	// https://binance-docs.github.io/apidocs/spot/cn/#trade
	STOPloss
	STOPlossLIMIT
	TAKEprofit
	TAKEprofitLIMIT
	LIMITmaker
)

func (t OrderType) String() string {
	switch t {
	case LIMIT:
		return "LIMIT"
	case MARKET:
		return "MARKET"
	case STOPloss:
		return "STOP_LOSS"
	case STOPlossLIMIT:
		return "STOP_LOSS_LIMIT"
	case TAKEprofit:
		return "TAKE_PROFIT"
	case TAKEprofitLIMIT:
		return "TAKE_PROFIT_LIMIT"
	case LIMITmaker:
		return "LIMIT_MAKER"
	default:
		panic("meet UNKNOWN Order Type")
	}
}

// Order 是下单的格式
// TODO: 下单的 order 和挂单的 order 有什么区别吗？
type Order struct {
	Symbol      string
	AssetName   string
	CapitalName string
	// if ID is negative value, means unset
	// ID is time.Now().Unix()
	ID   int64
	Side OrderSide
	Type OrderType
	// 根据 Type 的不同，以下 3 个属性不是全都必须的
	AssetQuantity   float64
	AssetPrice      float64
	CapitalQuantity float64
}

// NewOrder returns a order with default Symbol, Asset, Capital.
// return value is INCOMPLETE.
// It need run 'With' method to make a complete copy.
func NewOrder(symbol, asset, capital string) Order {
	return Order{
		Symbol:      symbol,
		AssetName:   asset,
		CapitalName: capital,
		ID:          -1,
	}
}

// With 可以生成一个根据 apply 实施的新订单
func (o Order) With(apply func(*Order)) *Order {
	res := o // deep copy
	apply(&res)
	res.ID = time.Now().Unix()
	return &res
}

// Limit 会按照限价单的方式设置订单
func Limit(side OrderSide, quantity, price float64) func(*Order) {
	return func(o *Order) {
		o.Type = LIMIT
		o.Side = side
		o.AssetQuantity = quantity
		o.AssetPrice = price
	}
}

// Market 会按照市价单的方式设置订单
func Market(side OrderSide, quantity float64) func(*Order) {
	return func(o *Order) {
		o.Type = MARKET
		o.Side = side
		switch side {
		case BUY:
			o.CapitalQuantity = quantity
		case SELL:
			o.AssetQuantity = quantity
		}
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

// IsLessThan return true if o < a
func (o *Order) IsLessThan(a *Order) bool {
	if o.Side != a.Side {
		panic("only compare with the same side")
	}
	if o.Type == MARKET && a.Type == MARKET {
		return o.ID < a.ID
	}
	if o.Type == MARKET && a.Type == LIMIT {
		return true
	}
	if o.Type == LIMIT && a.Type == MARKET {
		return false
	}
	// o.Type == LIMIT && a.Type == LIMIT
	if o.Side == BUY {
		// BUY 订单是价高的先成交
		return o.AssetPrice > a.AssetPrice
	}
	// o.Side == SELL
	// SELL 订单是价低的先成交
	return o.AssetPrice < a.AssetPrice
}

type cancelOrder struct {
	ID int64
}
