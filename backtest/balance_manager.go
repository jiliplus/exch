package backtest

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/jujili/exch"
)

// balanceManager
// 有一个 balance 帐户和一个 publisher
// 当 balance 的值发生变动时，
// 会利用 pulisher 把变动后的值，发送到 "balance" 话题
// TODO: 完成这个功能
type balanceManager struct {
	Balance exch.Balance
	pub     Publisher
	enc     func(interface{}) []byte
}

func newBalanceManager(pub Publisher, bal exch.Balance) *balanceManager {
	return &balanceManager{
		Balance: bal,
		pub:     pub,
		enc:     exch.EncFunc(),
	}
}

// NOTICE: 并没有核查 bm 内资产的 total，有可能 total 是负值
func (bm *balanceManager) update(as ...exch.Asset) {
	bm.Balance = bm.Balance.Add(as...)
	payload := bm.enc(bm.Balance)
	msg := message.NewMessage(watermill.NewUUID(), payload)
	go bm.pub.Publish("balance", msg)
	// TODO: 为什么这里总是空的
	// log.Println("balance:", bm.Balance)
}
