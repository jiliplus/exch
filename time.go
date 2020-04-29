package exch

import (
	"time"
)

// BeginFunc 会根据 time 和 interval 计算 time 所在周期的开始时间
type BeginFunc func(time.Time, time.Duration) time.Time

// Begin 会根据 time 和 interval 计算 time 所在周期的开始时间
// Begin 认为 1970-01-01 00:00:00 +0000 UTC 是第一个周期的起点，
// 然后计算当前周期的起点
// Begin 不会改变 t 的 location 信息
func Begin(t time.Time, d time.Duration) time.Time {
	loc := t.Location()
	utc := t.Unix()
	sec := int64(d / time.Second)
	utc = utc / sec * sec
	return time.Unix(utc, 0).In(loc)
}
