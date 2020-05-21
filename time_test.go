package exch

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Begin(t *testing.T) {
	layout := "2006-01-02 15:04:05.99"
	Convey("测试 Begin 函数", t, func() {
		input, err := time.Parse(layout, "2006-01-02 15:59:58.99")
		So(err, ShouldBeNil)
		title := fmt.Sprintf("输入时间为 %s，%s", input, input.Weekday())
		Convey(title, func() {
			Convey("间隔为 1 分钟", func() {
				d := time.Minute
				expected, err := time.Parse(layout, "2006-01-02 15:59:00")
				So(err, ShouldBeNil)
				actual := Begin(input, d)
				So(actual, ShouldEqual, expected)
			})
			Convey("间隔为 2 分钟", func() {
				d := time.Minute * 2
				expected, err := time.Parse(layout, "2006-01-02 15:58:00")
				So(err, ShouldBeNil)
				actual := Begin(input, d)
				So(actual, ShouldEqual, expected)
			})
			Convey("间隔为 3 分钟", func() {
				d := time.Minute * 3
				expected, err := time.Parse(layout, "2006-01-02 15:57:00")
				So(err, ShouldBeNil)
				actual := Begin(input, d)
				So(actual, ShouldEqual, expected)
			})
			Convey("间隔为 4 分钟", func() {
				d := time.Minute * 4
				expected, err := time.Parse(layout, "2006-01-02 15:56:00")
				So(err, ShouldBeNil)
				actual := Begin(input, d)
				So(actual, ShouldEqual, expected)
			})
			Convey("间隔为 5 分钟", func() {
				d := time.Minute * 5
				expected, err := time.Parse(layout, "2006-01-02 15:55:00")
				So(err, ShouldBeNil)
				actual := Begin(input, d)
				So(actual, ShouldEqual, expected)
			})
			Convey("间隔为 6 分钟", func() {
				d := time.Minute * 6
				expected, err := time.Parse(layout, "2006-01-02 15:54:00")
				So(err, ShouldBeNil)
				actual := Begin(input, d)
				So(actual, ShouldEqual, expected)
			})
			Convey("间隔为 10 分钟", func() {
				d := time.Minute * 10
				expected, err := time.Parse(layout, "2006-01-02 15:50:00")
				So(err, ShouldBeNil)
				actual := Begin(input, d)
				So(actual, ShouldEqual, expected)
			})
			Convey("间隔为 12 分钟", func() {
				d := time.Minute * 12
				expected, err := time.Parse(layout, "2006-01-02 15:48:00")
				So(err, ShouldBeNil)
				actual := Begin(input, d)
				So(actual, ShouldEqual, expected)
			})
			Convey("间隔为 15 分钟", func() {
				d := time.Minute * 15
				expected, err := time.Parse(layout, "2006-01-02 15:45:00")
				So(err, ShouldBeNil)
				actual := Begin(input, d)
				So(actual, ShouldEqual, expected)
			})
			Convey("间隔为 20 分钟", func() {
				d := time.Minute * 20
				expected, err := time.Parse(layout, "2006-01-02 15:40:00")
				So(err, ShouldBeNil)
				actual := Begin(input, d)
				So(actual, ShouldEqual, expected)
			})
			Convey("间隔为 30 分钟", func() {
				d := time.Minute * 30
				expected, err := time.Parse(layout, "2006-01-02 15:30:00")
				So(err, ShouldBeNil)
				actual := Begin(input, d)
				So(actual, ShouldEqual, expected)
			})
			Convey("间隔为 60 分钟", func() {
				d := time.Minute * 60
				expected, err := time.Parse(layout, "2006-01-02 15:00:00")
				So(err, ShouldBeNil)
				actual := Begin(input, d)
				So(actual, ShouldEqual, expected)
			})
			Convey("间隔为 14 分钟", func() {
				d := time.Minute * 14
				expected, err := time.Parse(layout, "2006-01-02 15:52:00")
				So(err, ShouldBeNil)
				actual := Begin(input, d)
				So(actual, ShouldEqual, expected)
			})
			Convey("间隔为 2 小时", func() {
				d := 2 * time.Hour
				expected, err := time.Parse(layout, "2006-01-02 14:00:00")
				So(err, ShouldBeNil)
				actual := Begin(input, d)
				So(actual, ShouldEqual, expected)
			})
			Convey("间隔为 4 小时", func() {
				d := 4 * time.Hour
				expected, err := time.Parse(layout, "2006-01-02 12:00:00")
				So(err, ShouldBeNil)
				actual := Begin(input, d)
				So(actual, ShouldEqual, expected)
			})
			Convey("间隔为 7 天", func() {
				d := 7 * 24 * time.Hour
				expected, err := time.Parse(layout, "2006-01-02 00:00:00")
				So(err, ShouldBeNil)
				actual := Begin(input, d)
				So(actual, ShouldEqual, expected)
			})
		})
	})
	Convey("Begin 不会改变时区信息", t, func() {
		shanghai, err := time.LoadLocation("Asia/Shanghai")
		So(err, ShouldBeNil)
		input, err := time.Parse(layout, "2006-01-02 15:59:58.99")
		So(err, ShouldBeNil)
		input = input.In(shanghai)
		expected := input.Location()
		title := fmt.Sprintf("输入时间的时区为 %s", expected)
		Convey(title, func() {
			actual := Begin(input, time.Minute).Location()
			title := fmt.Sprintf("输出时间的时区为 %s", actual)
			Convey(title, func() {
				So(actual, ShouldEqual, expected)
			})
		})
	})
}
