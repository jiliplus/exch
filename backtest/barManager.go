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

// TickBarService 会负责接受从 "tick" 话题中接受 tick 数据，
// 生成 Bar 后，会发送数据到对应的话题中。
// 例如，生成日 bar 线后，发送到 "24h0m0sBar" 话题中
// 例如，生成 30 日 bar 线后，发送到 "720h0m0sBar" 话题中
func TickBarService(ctx context.Context, ps Pubsub, interval time.Duration) {
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
