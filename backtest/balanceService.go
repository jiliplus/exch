package backtest

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jujili/exch"
)

// BalanceService 会在每天的凌晨零点零分零秒记录 balance 的总价值
// TODO: 从这里开始
func BalanceService(ctx context.Context, ps pubsub, interval time.Duration) {
	topic := fmt.Sprintf("%sBar", interval)
	log.Printf(`从 tick 生成的 bar 会发送到 "%s" 话题中`, topic)
	//
	ticks, err := ps.Subscribe(ctx, "tick")
	if err != nil {
		panic(err)
	}
	decTick := exch.DecTickFunc()
	//
	gtb := exch.GenTickBarFunc(exch.Begin, interval)
	//
	enc := exch.EncFunc()
	//
	var bars []exch.Bar
	//
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Fatalln("TickBarService Down: ", ctx.Err())
			case msg, ok := <-ticks:
				msg.Ack()
				if !ok {
					bars = gtb(exch.NilTick)
				} else {
					tick := decTick(msg.Payload)
					bars = gtb(tick)
				}
				for _, bar := range bars {
					msg := message.NewMessage(watermill.NewUUID(), enc(bar))
					ps.Publish(topic, msg)
				}
			}
		}
	}()
}
