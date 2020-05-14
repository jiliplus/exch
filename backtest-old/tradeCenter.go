package backtest

import "github.com/ThreeDotsLabs/watermill/message"

// TradeCenter 是一个模拟的交易中心
type TradeCenter struct {
}

type pubsub interface {
	Publish(topic string, messages ...*message.Message) error
	Subscribe(topic string) (<-chan *message.Message, error)
	Close() error
}
