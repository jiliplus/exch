package exch

import (
	"bytes"
	"encoding/gob"
	"time"
)

// Tick 记录了单笔的交易记录
// 各个具体交易所的交易记录
// 要么直接使用 Tick，
// 要么提供转换到 Tick 函数，
type Tick struct {
	// Exchange Name
	// Symbol    string // like "BTCUSDT"
	// Asset  string // like "BTC"
	ID     int64
	Date   time.Time
	Price  float64
	Volume float64
	// Type   string
}

// NewTick returns a new tick
func NewTick(id int64, date time.Time, price, volume float64) Tick {
	return Tick{
		ID:     id,
		Date:   date,
		Price:  price,
		Volume: volume,
	}
}

// DecTickFunc 返回的函数会把序列化成 []byte 的 Balances 值转换回来
// TODO: 把最终返回值修改成 Tick，把 Tick 当作值对象，不需要返回指针
func DecTickFunc() func(bs []byte) *Tick {
	var bb bytes.Buffer
	dec := gob.NewDecoder(&bb)
	return func(bs []byte) *Tick {
		bb.Reset()
		bb.Write(bs)
		var tick Tick
		// dec.Decode 只有在输入不是指针时候，才会报错
		// 显然 &balances 肯定是一个指针
		dec.Decode(&tick)
		return &tick
	}
}
