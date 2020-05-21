package exch

import (
	"bytes"
	"encoding/gob"
	"time"
)

// Bar 实现了 k 线的相关方法
type Bar struct {
	Begin                  time.Time
	Interval               time.Duration
	Open, High, Low, Close float64 // Price
	Volume                 float64
}

// DecBarFunc 返回的函数会把序列化成 []byte 的 Bar 值转换回来
func DecBarFunc() func(bs []byte) Bar {
	var bb bytes.Buffer
	dec := gob.NewDecoder(&bb)
	return func(bs []byte) Bar {
		bb.Reset()
		bb.Write(bs)
		var bar Bar
		dec.Decode(&bar)
		return bar
	}
}

// newTickBar make the first bar from a tick
func newTickBar(tick *Tick, begin time.Time, interval time.Duration) *Bar {
	tU := tick.Date.Unix()
	beginU := begin.Unix()
	endU := begin.Add(interval).Unix()
	if !(beginU <= tU && tU < endU) {
		panic("newTickBar: tick should in begin,interval")
	}
	return &Bar{
		Begin:    begin,
		Interval: interval,
		Open:     tick.Price,
		High:     tick.Price,
		Low:      tick.Price,
		Close:    tick.Price,
		Volume:   tick.Volume,
	}
}

// TODO: uncomment this block, 添加 bar 生成 bar 的相关方法
// // newBarBar make the first bar from another kind bar
// func newBarBar(bar *Bar, begin time.Time, interval time.Duration) *Bar {
// 	bInterval := bar.Interval
// 	if interval <= bInterval {
// 		panic("newBarBar: 新 bar 应该比旧 bar 宽")
// 	}
// 	bBegin, bEnd := bar.Begin, bar.Begin.Add(bInterval)
// 	end := begin.Add(interval)
// 	if !(begin.Unix() < bBegin.Unix() && bEnd.Unix() < end.Unix()) {
// 		panic("newBarBar: 新 bar 应该完全包住旧 bar")
// 	}
// 	if interval%bInterval != 0 {
// 		panic("newBarBar: 新 bar 的宽度，应该是旧 bar 的整数倍")
// 	}
// 	if bBegin.Sub(begin)%bInterval != 0 {
// 		panic("newBarBar: 新旧 bar 要能够对齐")
// 	}
// 	res := *bar
// 	res.Begin = begin
// 	res.Interval = interval
// 	return &res
// }

// GenTickBarFunc 会返回一个接收 tick 并生成 bar 的闭包函数
// 有以下情况需要处理
// 1. 接收第一个 tick,
//    不返回 bar
// 2. 接收到当前的 interval 的 tick
//    不返回 bar
// 3. 接收到下一个 interval 的 tick
//    返回上一个 bar
// 4. 接收到下一个 interval 后面的 interval 的 tick，市场冷清，长时间没有交易
//    返回多个 bar
func GenTickBarFunc(begin BeginFunc, interval time.Duration) func(*Tick) []*Bar {
	isInited := false
	var bar *Bar
	var lastTickDate time.Time
	return func(tick *Tick) []*Bar {
		tickBegin := begin(tick.Date, interval)
		if !isInited {
			bar = newTickBar(tick, tickBegin, interval)
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
		bar = newTickBar(tick, tickBegin, interval)
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
