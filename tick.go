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
	Exchange Name
	Symbol   string
	ID       int64
	Date     time.Time
	Price    float64
	Volume   float64
	Type     string
}

// DecTick 会把序列化成 []byte 的 Tick 值转换回来
func DecTick(bs []byte) *Tick {
	var tick Tick
	bb := bytes.NewBuffer(bs)
	dec := gob.NewDecoder(bb)
	// Decode 只有在输入不是指针时候，才会报错
	// 显然 &tick 肯定是一个指针
	dec.Decode(&tick)
	return &tick
}
