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

func Test_Asset_Total(t *testing.T) {
	Convey("测试 Asset.Total", t, func() {
		a := NewAsset("BTC", 1, 2)
		Convey("会返回 Free 和 Locked 的总和", func() {
			So(a.Total(), ShouldEqual, 3)
		})
	})
}

func Test_Asset_Add(t *testing.T) {
	Convey("测试 Asset.Add", t, func() {
		b := NewAsset("BTC", 1, 2)
		Convey("Add 不同种类的 Asset", func() {
			l := NewAsset("LTC", 1, 2)
			Convey("会 panic", func() {
				So(func() {
					b.Add(l)
				}, ShouldPanic)
			})
		})
		Convey("Add 相同种类的 Asset 会汇总相同的的属性值", func() {
			a := NewAsset("BTC", 1, 2)
			actual := b.Add(a)
			expected := NewAsset("BTC", 2, 4)
			So(actual, ShouldResemble, expected)
		})
	})
}

func Test_Balance_Total(t *testing.T) {
	Convey("测试 Balance.Total", t, func() {
		b := NewAsset("BTC", 1, 2)
		l := NewAsset("LTC", 10, 20)
		e := NewAsset("ETC", 100, 200)
		u := NewAsset("USDT", 1000, 2000)
		bal := NewBalances(b, l, e, u)
		price := map[string]float64{
			"BTC": 1000,
			"LTC": 100,
			"ETC": 10,
			// lock USDT price
		}
		Convey("缺少价格的话，会 panic", func() {
			So(func() {
				bal.Total(price)
			}, ShouldPanic)
		})
		price["USDT"] = 1
		Convey("提供完整的价目表，则可以计算出来结果", func() {
			actual := bal.Total(price)
			expected := 3000 * 4
			So(actual, ShouldEqual, expected)
		})
		price["DOGE"] = 0.001
		Convey("提供了不必要的价格，不会影响计算出来结果", func() {
			actual := bal.Total(price)
			expected := 3000 * 4
			So(actual, ShouldEqual, expected)
		})
	})
}

func Test_Balance_Add(t *testing.T) {
	Convey("测试 Balance.Add", t, func() {
		b := NewAsset("BTC", 1, 2)
		bal := NewBalances()
		Convey("Add 第 1 种 Asset", func() {
			bal.Add(b)
			Convey("Bal 会发生改变", func() {
				So((*bal)[b.Name], ShouldResemble, b)
			})
			Convey("再次添加会翻倍", func() {
				bal.Add(b)
				So((*bal)[b.Name], ShouldResemble, b.Add(b))
			})
		})
	})
}
