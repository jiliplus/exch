package exch

import (
	"time"
)

// BeginFunc 会根据 time 和 interval 计算 time 所在周期的开始时间
type BeginFunc func(time.Time, time.Duration) time.Time

// Begin 会根据 time 和 interval 计算 time 所在周期的开始时间
// 当 interval 的单位
//   为分钟或秒时，推荐值为 1,2,3,4,5,6,10,12,15,20,30,60
//   为小时时，   推荐值为 1,2,3,4,6,8,12,24
// NOTICE: 由于每个月的时间长度不一致，无法计算月线的起始日期。年线同理。
func Begin(date time.Time, interval time.Duration) time.Time {
	return date.Add(-interval / 2).Round(interval)
}
