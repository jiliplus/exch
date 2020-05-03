package exch

import (
	"bytes"
	"encoding/gob"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func getABalances() *Balance {
	a1 := NewAsset("BTC", 100, 200)
	a2 := NewAsset("DOGE", 900000000, 100000000)
	return NewBalances(a1, a2)
}

func Test_DecBalancesFunc(t *testing.T) {

	Convey("测试 Balances 的 Decode 函数", t, func() {
		expected := getABalances()
		enc := EncFunc()
		bs := enc(expected)
		dec := DecBalanceFunc()
		actual := dec(bs)
		So(*actual, ShouldResemble, *expected)
	})
}

// decBalances 会把序列化成 []byte 的 Balances 值转换回来
func decBalances(bs []byte) *Balance {
	var balances Balance
	bb := bytes.NewBuffer(bs)
	dec := gob.NewDecoder(bb)
	// Decode 只有在输入不是指针时候，才会报错
	// 显然 &balances 肯定是一个指针
	dec.Decode(&balances)
	return &balances
}

func Benchmark_DecBalance(b *testing.B) {
	bl := getABalances()
	enc := EncFunc()
	bs := enc(bl)
	for i := 1; i < b.N; i++ {
		decBalances(bs)
	}
}

func Benchmark_DecBalanceFunc(b *testing.B) {
	bls := getABalances()
	enc := EncFunc()
	bs := enc(bls)
	dec := DecBalanceFunc()
	for i := 1; i < b.N; i++ {
		dec(bs)
	}
}
