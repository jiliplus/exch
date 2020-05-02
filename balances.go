package exch

import (
	"bytes"
	"encoding/gob"
)

// Balances 记录了交易所中的资产
type Balances map[string]*Asset

// NewBalances returns a new Balances
func NewBalances(assets ...*Asset) *Balances {
	b := make(Balances, len(assets))
	for _, a := range assets {
		b[a.Name] = a
	}
	return &b
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

// Asset 代表了交易所中，某一项资产的状态和数目
type Asset struct {
	Name         string
	Free, Locked float64
}

// NewAsset return new asset pointer
func NewAsset(name string, free, locked float64) *Asset {
	return &Asset{
		Name:   name,
		Free:   free,
		Locked: locked,
	}
}

// Lock return true if lock f successfully
// otherwise return false and do NOT change asset
func (a *Asset) Lock(f float64) bool {
	if a.Free < f {
		return false
	}
	a.Free -= f
	a.Locked += f
	return true
}

// Unlock return true if unlock f successfully
// otherwise return false and do NOT change asset
func (a *Asset) Unlock(f float64) bool {
	if a.Locked < f {
		return false
	}
	a.Free += f
	a.Locked -= f
	return true
}

// UnlockAll unlock all locked asset
func (a *Asset) UnlockAll() {
	a.Free += a.Locked
	a.Locked = 0
}

// LockAll lock all free asset
func (a *Asset) LockAll() {
	a.Locked += a.Free
	a.Free = 0
}

// Total returns total asset of this asset
func (a *Asset) Total() float64 {
	return a.Free + a.Locked
}
