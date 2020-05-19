package backtest

import (
	"context"
	"log"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jujili/exch"
)

type pubsub interface {
	Publish(topic string, messages ...*message.Message) error
	Subscribe(ctx context.Context, topic string) (<-chan *message.Message, error)
	Close() error
}

// BackTest 是一个模拟的交易中心
type BackTest struct {
}

// NewBackTest returns a new trade center - bt
// bt subscribe "tick", "bar" and "order" topics from pubsub
// and
// bt publish "balance" topic
//
func NewBackTest(ctx context.Context, ps pubsub, balance exch.Balance) {
	sells := newOrderList()
	buys := newOrderList()

	ticks, err := ps.Subscribe(ctx, "tick")
	if err != nil {
		panic(err)
	}

	bars, err := ps.Subscribe(ctx, "bar")
	if err != nil {
		panic(err)
	}

	orders, err := ps.Subscribe(ctx, "order")
	if err != nil {
		panic(err)
	}

	decOrder := decOrderFunc()
	decTick := exch.DecTickFunc()

	go func() {
		select {
		case <-ctx.Done():
			log.Println("ctx.Done", ctx.Err())
		case msg := <-ticks:
			tick := decTick(msg.Payload)
			buys.match(tick)
			sells.match(tick)
			msg.Ack()
		case <-bars:
			// case msg := <-bars:
			panic("现在还不能处理 bar 数据")
			// msg.Ack()
		case msg := <-orders:
			order := decOrder(msg.Payload)
			if order.Side == exch.BUY {
				buys.push(order)
			} else {
				sells.push(order)
			}
			msg.Ack()
		}
	}()
}
