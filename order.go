package exch

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
	Side        OrderSide
	// 根据 Type 的不同，以下 3 个属性不是全都必须的
	AssetQuantity   float64
	AssetPrice      float64
	CapitalQuantity float64
}
