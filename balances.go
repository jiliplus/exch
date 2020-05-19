package exch

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// Balance 记录了交易所中的资产
type Balance map[string]Asset

// NewBalances returns a new Balances
func NewBalances(assets ...Asset) *Balance {
	b := make(Balance, len(assets))
	for _, a := range assets {
		b[a.Name] = a
	}
	return &b
}

// DecBalanceFunc 返回的函数会把序列化成 []byte 的 Balances 值转换回来
func DecBalanceFunc() func(bs []byte) *Balance {
	var bb bytes.Buffer
	dec := gob.NewDecoder(&bb)
	return func(bs []byte) *Balance {
		bb.Reset()
		bb.Write(bs)
		var balance Balance
		// dec.Decode 只有在输入不是指针时候，才会报错
		// 显然 &balances 肯定是一个指针
		dec.Decode(&balance)
		return &balance
	}
}

// Add change *Balance with a slice of Asset
func (b *Balance) Add(as ...Asset) {
	for _, a := range as {
		old, ok := (*b)[a.Name]
		if ok {
			(*b)[a.Name] = old.Add(a)
		} else {
			(*b)[a.Name] = a
		}
	}
}

// Total count the total value of balance
func (b *Balance) Total(prices map[string]float64) float64 {
	var total float64
	for name, asset := range *b {
		price, ok := prices[name]
		if !ok {
			msg := fmt.Sprintf("Balance.Total: %s do NOT have a price", name)
			panic(msg)
		}
		total += asset.Total() * price
	}
	return total
}

// Asset 代表了交易所中，某一项资产的状态和数目
// Asset 是一个值对象 value object
type Asset struct {
	Name         string
	Free, Locked float64
}

// NewAsset return new asset
func NewAsset(name string, free, locked float64) Asset {
	return Asset{
		Name:   name,
		Free:   free,
		Locked: locked,
	}
}

// Add return new asset by a change with delta
func (a Asset) Add(delta Asset) Asset {
	if a.Name != delta.Name {
		panic("Asset can NOT change with a different asset")
	}
	return NewAsset(
		a.Name,
		a.Free+delta.Free,
		a.Locked+delta.Locked,
	)
}

// Total returns total asset of this asset
func (a Asset) Total() float64 {
	return a.Free + a.Locked
}
