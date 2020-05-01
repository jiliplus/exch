package exch

import (
	"bytes"
	"encoding/gob"
)

// Balances 记录了交易所中的资产
type Balances struct {
	Asset Asset
	// Universal Equivalent
	UE string
}

// Asset 代表了交易所中，某一项资产的状态和数目
type Asset struct {
	Name         string
	Free, Locked float64
	// Price 是相对于 Balances 中 UE 的价格
	Price float64
}

// DecBalancesFunc 返回的函数会把序列化成 []byte 的 Balances 值转换回来
func DecBalancesFunc() func(bs []byte) *Balances {
	var bb bytes.Buffer
	dec := gob.NewDecoder(&bb)
	return func(bs []byte) *Balances {
		bb.Reset()
		bb.Write(bs)
		var balances Balances
		// dec.Decode 只有在输入不是指针时候，才会报错
		// 显然 &balances 肯定是一个指针
		dec.Decode(&balances)
		return &balances
	}
}
