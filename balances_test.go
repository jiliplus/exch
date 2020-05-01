package exch

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_DecBalancesFunc(t *testing.T) {
	Convey("测试 Balances 的 Decode 函数", t, func() {
		expected := &Balances{
			Asset: Asset{
				Name:   "BTC",
				Free:   100,
				Locked: 0,
				Price:  10000,
			},
			UE: "BTC",
		}
		enc := EncFunc()
		bs := enc(expected)
		dec := DecBalancesFunc()
		actual := dec(bs)
		fmt.Println(actual)
		So(*actual, ShouldResemble, *expected)
	})
}

// decBalances 会把序列化成 []byte 的 Balances 值转换回来
func decBalances(bs []byte) *Balances {
	var balances Balances
	bb := bytes.NewBuffer(bs)
	dec := gob.NewDecoder(bb)
	// Decode 只有在输入不是指针时候，才会报错
	// 显然 &balances 肯定是一个指针
	dec.Decode(&balances)
	return &balances
}

func Benchmark_DecBalance(b *testing.B) {
	bls := &Balances{
		Asset: Asset{
			Name:   "BTC",
			Free:   100,
			Locked: 0,
		},
		UE: "BTC",
	}
	enc := EncFunc()
	bs := enc(bls)
	for i := 1; i < b.N; i++ {
		decBalances(bs)
	}
}

func Benchmark_DecBalanceFunc(b *testing.B) {
	bls := &Balances{
		Asset: Asset{
			Name:   "BTC",
			Free:   100,
			Locked: 0,
		},
		UE: "BTC",
	}
	enc := EncFunc()
	bs := enc(bls)
	dec := DecBalancesFunc()
	for i := 1; i < b.N; i++ {
		dec(bs)
	}
}
