package exch

import (
	"bytes"
	"encoding/gob"
	"time"
)

// Bar 实现了 k 线的相关方法
type Bar struct {
	Begin                  time.Time
	Open, High, Low, Close float64
	Volume                 float64
	// 以下属性可以通过其他方式获取
	// 为了节约空间，不要保存在数据库中
	Symbol   string
	Exchange Name
	// TODO: Interval 有必要吗
	Interval time.Duration
}

// DecBarFunc 返回的函数会把序列化成 []byte 的 Bar 值转换回来
func DecBarFunc() func(bs []byte) *Bar {
	var bb bytes.Buffer
	dec := gob.NewDecoder(&bb)
	return func(bs []byte) *Bar {
		bb.Reset()
		bb.Write(bs)
		var bar Bar
		dec.Decode(&bar)
		return &bar
	}
}

func newBar(tick *Tick, date time.Time) *Bar {
	return &Bar{
		Begin:  date,
		Open:   tick.Price,
		High:   tick.Price,
		Low:    tick.Price,
		Close:  tick.Price,
		Volume: tick.Volume,
	}
}

// GenBarFunc 会返回一个接收 tick 并生成 bar 的闭包函数
// TODO: 完成这个闭包函数
// 有以下情况需要处理
// 1. 接收第一个 tick,
//    不返回 bar
// 2. 接收到当前的 interval 的 tick
//    不返回 bar
// 3. 接收到下一个 interval 的 tick
//    返回上一个 bar
// 4. 接收到下一个 interval 后面的 interval 的 tick，市场冷清，长时间没有交易
//    返回多个 bar
func GenBarFunc(begin BeginFunc, interval time.Duration) func(*Tick) []*Bar {
	isInited := false
	var bar *Bar
	var lastTickDate time.Time
	return func(tick *Tick) []*Bar {
		tickBegin := begin(tick.Date, interval)
		if !isInited {
			bar = newBar(tick, tickBegin)
			lastTickDate = tick.Date
			isInited = true
			return nil
		}
		// GenBar 不接受乱序的 ticks
		if tick.Date.Before(lastTickDate) {
			panic("GenBar: Ticks should be sorted in date")
		}
		lastTickDate = tick.Date
		//
		if tickBegin.Equal(bar.Begin) {
			bar.High = maxFloat64(bar.High, tick.Price)
			bar.Low = minFloat64(bar.Low, tick.Price)
			bar.Close = tick.Price
			bar.Volume += tick.Volume
			return nil
		}
		res := make([]*Bar, 0, 256)
		for bar.Begin.Before(tickBegin) {
			res = append(res, bar)
			bar = nextEmptyBar(bar, interval)
		}
		bar = newBar(tick, tickBegin)
		return res
	}
}

func nextEmptyBar(bar *Bar, interval time.Duration) *Bar {
	return &Bar{
		Begin: bar.Begin.Add(interval),
		Open:  bar.Close,
		High:  bar.Close,
		Low:   bar.Close,
		Close: bar.Close,
	}
}
