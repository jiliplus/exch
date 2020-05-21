package exch

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_GenBarFunc(t *testing.T) {
	Convey("把 GenBar 赋值给 gb", t, func() {
		date := time.Now()
		interval := time.Minute
		gb := GenBarFunc(Begin, interval)
		tick := &Tick{
			Date:   date,
			Price:  1.0,
			Volume: 1.0,
		}
		Convey("输入第一个 tick", func() {
			actual := gb(tick)
			Convey("返回值应该是 nil", func() {
				So(actual, ShouldBeNil)
			})
			Convey("输入更早的 tick 会 panic", func() {
				So(func() {
					tick := &Tick{Date: date.Add(-time.Second)}
					gb(tick)
				}, ShouldPanic)
			})
			Convey("输入同周期的 tick 会返回 nil", func() {
				actual := gb(tick)
				So(actual, ShouldBeNil)
				Convey("输入下个周期的 tick，会返回一个 bar", func() {
					tick.Date = tick.Date.Add(interval)
					bars := gb(tick)
					So(len(bars), ShouldEqual, 1)
					bar := bars[0]
					Convey("由于同一个 tick 输入了两遍", func() {
						So(bar.Open, ShouldEqual, bar.Close)
						So(bar.High, ShouldEqual, bar.Close)
						So(bar.Low, ShouldEqual, bar.Close)
						So(bar.Volume, ShouldEqual, tick.Volume*2)
					})
				})
				Convey("输入下两个周期的 tick，会返回 2 个 bar", func() {
					tick.Date = tick.Date.Add(2 * interval)
					bars := gb(tick)
					So(len(bars), ShouldEqual, 2)
					bar := bars[0]
					Convey("由于同一个 tick 输入了两遍", func() {
						So(bar.Open, ShouldEqual, bar.Close)
						So(bar.High, ShouldEqual, bar.Close)
						So(bar.Low, ShouldEqual, bar.Close)
						So(bar.Volume, ShouldEqual, tick.Volume*2)
					})
					emptyBar := bars[1]
					Convey("第二个 bar 应该是 empty 的", func() {
						So(emptyBar.Volume, ShouldEqual, 0)
						So(emptyBar.Open, ShouldEqual, bar.Close)
						So(emptyBar.High, ShouldEqual, bar.Close)
						So(emptyBar.Low, ShouldEqual, bar.Close)
						So(emptyBar.Close, ShouldEqual, bar.Close)
					})
				})
			})
		})
	})
}

func Test_DecBarFunc(t *testing.T) {
	Convey("反向序列化 Tick", t, func() {
		expected := Bar{
			Begin:    time.Now(),
			Interval: time.Minute,
			Open:     1,
			High:     2,
			Low:      3,
			Close:    4,
			Volume:   5,
		}
		enc := EncFunc()
		dec := DecBarFunc()
		actual := dec(enc(expected))
		Convey("指针指向的对象应该不同", func() {
			So(actual, ShouldNotEqual, expected)
			Convey("具体的值，应该相同", func() {
				So(actual.Begin.Equal(expected.Begin), ShouldBeTrue)
				actual.Begin = expected.Begin
				// 没有上面两行，直接使用下面的判断语句会报错，
				So(actual, ShouldResemble, expected)
			})
		})
	})
}
