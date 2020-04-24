package backtest

import (
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
)

// TODO: 把 orderType 改成枚举类型
type orderType int

// 以下订单类型来自于:
// https://binance-docs.github.io/apidocs/spot/cn/#trade-2
const (
	limit orderType = iota
	market
	stopLoss
	stopLossLimit
	takeProfit
	takeProfitLimit
	limitMaker
)

func (t orderType) String() string {
	switch t {
	case limit:
		return "限价单"
	case market:
		return "市价单"
	case stopLoss:
		return "止损单"
	case stopLossLimit:
		return "现价止损单"
	case takeProfit:
		return "TAKE_PROFIT"
	case takeProfitLimit:
		return "TAKE_PROFIT_LIMIT"
	case limitMaker:
		return "LIMIT_MAKER"
	default:
		return "Unknown"
	}
}

type order struct {
	orderType orderType
	price     float64
	amount    float64
	total     float64
}

type tickTrader struct {
	sync.Mutex
	orders []order
}

// TODO: tickTrader 的初始化，需要包装一下 watermill 的 router
func (t *tickTrader) init() {
	return
}

// noPublishHandle
func (t *tickTrader) handleOder(msg *message.Message) {
	t.Lock()
	defer t.Unlock()

	return
}
