package backtest

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jujili/clock"
	"github.com/jujili/exch"
)

// BalanceService 会在每天的凌晨零点零分零秒记录 balance 的总价值
// prices 里面需要放好各种资产的价格，不要忘记 capital 的价格是 1
func BalanceService(ctx context.Context, ps Pubsub, prices map[string]float64, asset string) {
	ticks, err := ps.Subscribe(ctx, "tick")
	if err != nil {
		panic(err)
	}
	decTick := exch.DecTickFunc()
	//
	balances, err := ps.Subscribe(ctx, "balance")
	if err != nil {
		panic(err)
	}
	decBal := exch.DecBalanceFunc()
	// 创建模拟 clock
	msg := <-ticks
	tick := decTick(msg.Payload)
	prices[asset] = tick.Price
	clock := clock.NewSimulator(tick.Date)
	everyNewDay := clock.EveryDay(0, 0, 0)
	// 另起一个 goroutine，更新 clock
	go func() {
		tks, _ := ps.Subscribe(ctx, "tick")
		for msg := range tks {
			tick := decTick(msg.Payload)
			clock.SetOrPanic(tick.Date)
			// 设置好了才确认
			msg.Ack()
		}
	}()
	//
	go func() {
		var msg *message.Message
		var bal *exch.Balance
		bs := make([]balanceSnap, 0, 2048)
		ok := true
		for ok {
			select {
			case <-ctx.Done():
				log.Fatalln("BalanceService Down: ", ctx.Err())
			case msg, ok = <-ticks:
				tick := decTick(msg.Payload)
				prices[asset] = tick.Price
			case msg = <-balances:
				bal = decBal(msg.Payload)
				msg.Ack()
			case date := <-everyNewDay:
				newBal := newBalanceSnap(date, bal, prices)
				bs = append(bs, newBal)
				log.Println(date, newBal)
			}
		}
		fmt.Println(bs)
	}()
}

type balanceSnap struct {
	date   time.Time
	amount float64
}

func newBalanceSnap(date time.Time, balance *exch.Balance, prices map[string]float64) balanceSnap {
	return balanceSnap{
		date:   date,
		amount: balance.Total(prices),
	}
}
